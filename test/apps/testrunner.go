package main

import (
	"context"
	"fmt"

	"github.com/corezoid/gitcall-go-runner/gitcall"
)

func usercode(_ context.Context, data map[string]interface{}) error {
	switch data["case"] {
	case "success":
		data["foo"] = 123

		return nil
	case "error":
		return fmt.Errorf("error-happened")
	case "panic":
		panic("something went wrong")
	}

	return nil
}

func main() {
	gitcall.Handle(usercode)
}
