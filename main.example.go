package main

import (
	"context"

	"github.com/corezoid/gitcall-go-runner/runner"
)

func usercode(ctx context.Context, data map[string]interface{}) error {

	data["foo"] = "bar"

	return nil
}

func main() {
	runner.Run(usercode)
}
