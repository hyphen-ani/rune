package namespace

import "rune/internal/constants"

func Normalize(value string) string {
	if value == "" {
		return constants.DefaultNamespace
	}

	return value
}
