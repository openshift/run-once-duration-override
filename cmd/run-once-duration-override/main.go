package main

import (
	"github.com/openshift/run-once-duration-override/metrics"

	"github.com/openshift/generic-admission-server/pkg/cmd"
)

func main() {
	metrics.Register()
	cmd.RunAdmissionServer(&runOnceDurationOverrideHook{})
}
