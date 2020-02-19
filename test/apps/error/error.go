package main

import (
	"context"
	"fmt"

	"git.corezoid.com/gitcall/go-runner/runner"
)

func usercode(ctx context.Context, data map[string]interface{}) error {

	return fmt.Errorf("error-happened")
}

func main() {
	runner.Run(usercode)
}
