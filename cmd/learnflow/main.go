package main

import (
	"context"
	"learningflow/internal/app"
)

func main() {
	ctx := context.Background()
	a := app.New(ctx)
	a.Run(ctx)
}
