package alert

import "time"

// now is a package-level variable so tests can override it.
var now = func() time.Time {
	return time.Now()
}
