// +build live insecure

// This file allows for quick testing by force setting the crashKey in a live/development environment.
// This may need its tag, but that remains to be seen.

package authentication

func init() {
	crashKey = "12345"
}
