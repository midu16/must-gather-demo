// must-gather/pkg/gather/constants.go
package gather

const (
	gatherContainerName      = "gather"
	unreachableTaintKey      = "node.kubernetes.io/unreachable"
	volumeUsageCheckerScript = `set -euo pipefail

AVAILABLE=$(df -k %s | tail -1 | awk '{print $4}')
TOTAL=$(df -k %s | tail -1 | awk '{print $2}')
PERCENTAGE=$(( (100 * (TOTAL - AVAILABLE)) / TOTAL ))

if [ $PERCENTAGE -ge %d ]; then
    echo "Error: Volume usage is $PERCENTAGE%%, exceeds %d%% limit" >&2
    exit 1
fi

%s`
)