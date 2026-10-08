package t

import (
	"testing"
)

// Helper in different file
func crossFileHelper(t *testing.T) {
	t.Helper()
	t.Parallel()
}
