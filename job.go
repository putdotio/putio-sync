package main

import (
	"fmt"
)

type Job interface {
	Run() error
	String() string
}

type DeleteState struct {
	state State
}

func (j *DeleteState) String() string {
	return fmt.Sprintf("Deleting state %q", j.state.relpath)
}

func (j *DeleteState) Run() error {
	// TODO Delete temp download file
	// TODO Delete temp upload file
	return j.state.Delete()
}
