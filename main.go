package main

import (
	"context"

	"github.com/aethelgards/octo/setup"
	"github.com/aethelgards/octo/tui"
)

func main() {
	ctx := context.Background()
	setup.Init(ctx)
	err := tui.Init(ctx)
	if err != nil {
		panic(err)
	}
}
