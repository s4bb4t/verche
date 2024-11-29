package liner

import (
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`^([a-zA-Z0-9\.\-_/]+)(?:\s+)(v?[0-9]+\.[0-9]+\.[0-9]+.*)$`)

func TakeALook(line string) (string, string, bool) {
	if strings.HasPrefix(line, "go 1.") {
		return "", "", false
	}
	line = strings.Split(line, "//")[0]
	line = strings.TrimSpace(line)

	matches := re.FindStringSubmatch(line)
	if len(matches) == 3 {
		return matches[1], matches[2], true
	}
	return "", "", false
}
