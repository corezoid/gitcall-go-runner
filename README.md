# Gitcall usercode runner


# Usage
To build your Go usercode runner you need to include 
`github.com/corezoid/gitcall-go-runner/gitcall` package,
 implement `usercode` handler function and `main` function.
See example code:

```go
package main

import (
	"context"

	"github.com/corezoid/gitcall-go-runner/gitcall"
)

func main() {
	gitcall.Handle(func (_ context.Context, data map[string]interface{}) error {
        data["foo"] = "bar"
       
        return nil
   })
}
```

# Test

```bash
> make install
> make build-test
> make test
```