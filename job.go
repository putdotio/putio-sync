package main

import (
	"context"
)

type Job interface {
	Run(context.Context) error
	String() string
}
