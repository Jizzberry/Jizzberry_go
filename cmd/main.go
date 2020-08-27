package main

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers/initializer"
	"os"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	helpers.SetWorkingDirectory(wd)
	err = initializer.Init()
	if err != nil {
		panic(err)
	}
}
