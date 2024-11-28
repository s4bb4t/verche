package liner

import (
	"regexp"
	"strings"
)

func TakeALook(line string) (string, string, bool) {
	line = strings.Split(line, "//")[0]
	line = strings.TrimSpace(line)

	re := regexp.MustCompile(`^([a-zA-Z0-9\.\-_/]+)\s([vV]?[0-9]+\.[0-9]+\.[0-9]+.*)$`)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 3 {
		return matches[1], matches[2], true
	}
	return "", "", false
}
