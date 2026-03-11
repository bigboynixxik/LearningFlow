package main

import (
	"context"
	"learningflow/internal/app"
)

func main() {
	ctx := context.Background()
	a, err := app.New(ctx)
	if err != nil {
		panic(err)
	}
	a.Run()
}
