package main

import (
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers/initializer"
)

func main() {
	helpers.SetWorkingDirectory(".")
	err := initializer.Init()
	if err != nil {
		panic(err)
		return
	}
}
