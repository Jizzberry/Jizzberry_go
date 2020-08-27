package ffmpeg

import (
	"bytes"
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/google/uuid"
	"io"
	"os/exec"
	"sync"
	"time"
)

type ProcessType int

const (
	FFProbe         ProcessType = 0
	TranscodeStream ProcessType = 1
	TranscodeStore  ProcessType = 2
	FastStartConv   ProcessType = 3
	ThumbnailGen    ProcessType = 4

	TimeoutForget = -1
	TimeoutBlock  = -2
)

type Encoder struct {
	UID          string
	process      map[string]*ProcessHolder
	processMutex sync.Mutex
}

type ProcessHolder struct {
	Type      ProcessType
	Cmd       *exec.Cmd
	Timeout   int
	forceKill chan bool
}

func NewEncoder() *Encoder {
	uid := uuid.New().String()
	e := Encoder{
		UID:     uid,
		process: make(map[string]*ProcessHolder),
	}
	return &e
}

func (e *Encoder) Run(args []string, Type ProcessType, timeout int, wait bool) (out bytes.Buffer, stderr bytes.Buffer, uid string) {
	e.processMutex.Lock()

	cmd := buildArg(args, Type)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if wait {
		err := cmd.Run()
		if err != nil {
			helpers.LogError(err.Error())
		}
		e.processMutex.Unlock()
		return
	}

	err := cmd.Start()
	if err != nil {
		helpers.LogError(err.Error())
	}

	uid = uuid.New().String()
	e.process[uid] = &ProcessHolder{
		Type:    Type,
		Timeout: timeout,
		Cmd:     cmd,
	}

	e.determineHold(uid)
	e.processMutex.Unlock()

	return
}

func (e *Encoder) Pipe(args []string, Type ProcessType, timeout int) (io.ReadCloser, string) {
	e.processMutex.Lock()

	cmd := buildArg(args, Type)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		helpers.LogError(err.Error())
		e.processMutex.Unlock()
		return nil, ""
	}

	err = cmd.Start()
	if err != nil {
		helpers.LogError(err.Error())
		e.processMutex.Unlock()
		return nil, ""
	}

	helpers.LogInfo("new stream", args)

	uid := uuid.New().String()
	e.process[uid] = &ProcessHolder{
		Type:    Type,
		Timeout: timeout,
		Cmd:     cmd,
	}

	e.determineHold(uid)
	e.processMutex.Unlock()

	return stdout, uid
}

func buildArg(args []string, Type ProcessType) *exec.Cmd {
	return exec.Command(
		func() string {
			if Type == FFProbe {
				return helpers.GetConfig().FFPROBE
			} else {
				return helpers.GetConfig().FFMEPG
			}
		}(),
		args...)
}

func (e *Encoder) determineHold(index string) {
	switch e.process[index].Timeout {
	case TimeoutBlock:
		e.blockTillEnd(index)
		return
	case TimeoutForget:
		return
	default:
		e.killOnTimeout(index)
		return
	}
}

func (e *Encoder) blockTillEnd(index string) {
	err := e.process[index].Cmd.Wait()
	if err != nil {
		helpers.LogError(err.Error())
	}
}

func (e *Encoder) killOnTimeout(index string) {
	done := make(chan error)
	go func() { done <- e.process[index].Cmd.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			helpers.LogError(err.Error())
		}
		e.deregisterProcess(index)
		return
	case <-time.After(time.Duration(e.process[index].Timeout) * time.Second):
		e.KillProcess(index)
		return
	}
}

func (e *Encoder) KillProcess(index string) {
	e.processMutex.Lock()
	if e.process[index].Cmd != nil {
		err := e.process[index].Cmd.Process.Kill()
		if err != nil {
			helpers.LogError("Unable to kill process: ", err.Error(), " The process probably exited on its own")
			e.deregisterProcess(index)
			e.processMutex.Unlock()
			return
		}
		helpers.LogInfo(fmt.Sprintf("Killed process [%v]: %d", e.process[index].Type, e.process[index].Cmd.Process.Pid))
	}
	e.deregisterProcess(index)
	e.processMutex.Unlock()
}

func (e *Encoder) deregisterProcess(index string) {
	e.process[index] = nil
}
