package thread

import (
	"fmt"
)

func keyForStream(id string, runId uint64) string {
	return fmt.Sprintf("z:ts:%s:%d", id, runId)
}
