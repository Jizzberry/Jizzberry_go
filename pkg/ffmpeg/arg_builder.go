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

var processMutex sync.Mutex

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
	UID     string
	process *ProcessHolder
}

type ProcessHolder struct {
	Type      ProcessType
	Cmd       []*exec.Cmd
	Timeout   int
	forceKill chan bool
}

func NewEncoder() *Encoder {
	uid := uuid.New().String()
	e := Encoder{
		UID: uid,
		process: &ProcessHolder{
			Cmd: make([]*exec.Cmd, 0),
		},
	}
	return &e
}

func (e *Encoder) Run(args []string, Type ProcessType, timeout int, wait bool) (out bytes.Buffer, stderr bytes.Buffer) {
	processMutex.Lock()

	cmd := buildArg(args, Type)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if wait {
		err := cmd.Run()
		if err != nil {
			helpers.LogError(err.Error())
		}
		return
	}

	err := cmd.Start()
	if err != nil {
		helpers.LogError(err.Error())
	}

	e.process = &ProcessHolder{
		Type:    Type,
		Cmd:     append(e.process.Cmd, cmd),
		Timeout: timeout,
	}

	e.determineHold(len(e.process.Cmd) - 1)
	processMutex.Unlock()

	return
}

func (e *Encoder) Pipe(args []string, Type ProcessType, timeout int) io.ReadCloser {
	processMutex.Lock()

	cmd := buildArg(args, Type)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		helpers.LogError(err.Error())
		return nil
	}

	err = cmd.Start()
	if err != nil {
		helpers.LogError(err.Error())
		return nil
	}

	e.process = &ProcessHolder{
		Type:    Type,
		Cmd:     append(e.process.Cmd, cmd),
		Timeout: timeout,
	}

	e.determineHold(len(e.process.Cmd) - 1)
	processMutex.Unlock()

	return stdout
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

func (e *Encoder) determineHold(index int) {
	switch e.process.Timeout {
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

func (e *Encoder) blockTillEnd(index int) {
	err := e.process.Cmd[index].Wait()
	if err != nil {
		helpers.LogError(err.Error())
	}
}

func (e *Encoder) killOnTimeout(index int) {
	done := make(chan error)
	go func() { done <- e.process.Cmd[index].Wait() }()
	select {
	case err := <-done:
		if err != nil {
			helpers.LogError(err.Error())
		}
		e.deregisterProcess(index)
		return
	case <-time.After(time.Duration(e.process.Timeout) * time.Second):
		e.KillProcess(index)
		return
	}
}

func (e *Encoder) KillProcess(index int) {
	processMutex.Lock()
	if e.process.Cmd != nil && e.process.Cmd[index] != nil {
		err := e.process.Cmd[index].Process.Kill()
		if err != nil {
			helpers.LogError("Unable to kill process: ", err.Error(), " The process probably exited on its own")
			e.deregisterProcess(index)
			processMutex.Unlock()
			return
		}
		helpers.LogInfo(fmt.Sprintf("Killed process [%v]: %d", e.process.Type, e.process.Cmd[index].Process.Pid))
	}
	e.deregisterProcess(index)
	processMutex.Unlock()
}

func (e *Encoder) deregisterProcess(index int) {
	e.process.Cmd = append(e.process.Cmd[:index], e.process.Cmd[index+1:]...)
}

func (e *Encoder) KillPrev() {
	e.KillProcess(0)
}
