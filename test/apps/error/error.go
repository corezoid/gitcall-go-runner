package main

import (
	"context"
	"fmt"

	"github.com/corezoid/gitcall-go-runner/gitcall"
)

func usercode(_ context.Context, _ map[string]interface{}) error {

	return fmt.Errorf("error-happened")
}

func main() {
	gitcall.Handle(usercode)
}
