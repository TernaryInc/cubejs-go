package cube

import (
	"fmt"
	"time"
)

// QualifiedDimension prepends the specified dimension string with the requested schema
// TODO: Is there a way we can still have a typed schema here?  embedding??
func QualifiedDimension(schema string, dimension string) string {
	return fmt.Sprintf("%s.%s", schema, dimension)
}

// An interface that allows us to mock time results freely but use real time in production.
type Nower interface {
	Now() time.Time
}

// An implementation of Nower which uses the real time.
type TimeNower struct{}

// Return the real time.
func (n TimeNower) Now() time.Time {
	return time.Now()
}

// An implementation of Nower which returns a specified time.
// Generally for use in tests.
type TestNower struct {
	T time.Time
}

// Return the specified time.
func (n TestNower) Now() time.Time {
	return n.T
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
