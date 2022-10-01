package helpers

import "fmt"

func CreateSocketAddress(processID string) string {
	return fmt.Sprintf("/tmp/%s.sock", processID)
}
