package main

import (
	"context"
	"fmt"
)

func sync() error {
	ai, err := client.Account.Info(context.TODO())
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", ai)
	return nil
}
