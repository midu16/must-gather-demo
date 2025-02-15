// must-gather/pkg/gather/options.go
package gather

import (
	"time"
)

type MustGatherOptions struct {
	Command          []string
	SourceDir        string
	VolumePercentage int
	HostNetwork      bool
	Since            time.Duration
	SinceTime        string
}