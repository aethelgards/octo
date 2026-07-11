package main

import (
	"context"

	"github.com/aethelgards/octo/setup"
)

func main() {
	ctx := context.Background()
	setup.Init(ctx)
}
