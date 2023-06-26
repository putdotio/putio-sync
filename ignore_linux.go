package putiosync

import "strings"

func shouldIgnoreName(name string) bool {
	return strings.ContainsAny(name, `\x00/`)
}
