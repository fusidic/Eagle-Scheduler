package register

import (
	"github.com/fusidic/Eagle-Scheduler/pkg/eagle"
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
)

// Register custom plugins to scheduler framework.
func Register() *cobra.Command {
	return app.NewSchedulerCommand(
		app.WithPlugin(eagle.Name, eagle.New),
	)
}
