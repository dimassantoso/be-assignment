package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"runtime/debug"

	"github.com/golangid/candi/codebase/app"
	"github.com/golangid/candi/config"

	service "billing-engine/internal"
)

func main() {
	const serviceName = "billing-engine"

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\x1b[31;1mFailed to start %s service: %v\x1b[0m\n", serviceName, r)
			fmt.Printf("Stack trace: \n%s\n", debug.Stack())
		}
	}()

	cfg := config.Init(serviceName)
	defer cfg.Exit()

	app.New(service.NewService(cfg)).Run()
}
