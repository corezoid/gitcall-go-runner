package main

import (
	"context"

	"git.corezoid.com/gitcall/go-runner/runner"
)

func usercode(ctx context.Context, data map[string]interface{}) error {

	data["foo"] = 123

	return nil
}

func main() {
	runner.Run(usercode)
}
