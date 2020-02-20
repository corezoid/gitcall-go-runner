package main

import (
	"context"

	"github.com/corezoid/gitcall-go-runner/runner"
)

func usercode(ctx context.Context, data map[string]interface{}) error {

	data["foo"] = 123

	return nil
}

func main() {
	runner.Run(usercode)
}
