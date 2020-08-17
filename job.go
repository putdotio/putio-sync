package putiosync

import (
	"context"
)

type iJob interface {
	Run(context.Context) error
	String() string
}
