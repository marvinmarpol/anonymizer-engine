package tuple

import "strings"

func HasPrefixInList(s *string, prefixes *[]string, removePrefix bool) bool {
	for _, prefix := range *prefixes {
		if strings.HasPrefix(*s, prefix) {
			if removePrefix {
				// Remove the prefix from the string
				*s = strings.TrimPrefix(*s, prefix)
			}
			return true
		}
	}
	return false
}
