package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/fusidic/Eagle-Scheduler/pkg/register"
	"k8s.io/component-base/logs"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	defer logs.FlushLogs()
	command := register.Register()
	if err := command.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
