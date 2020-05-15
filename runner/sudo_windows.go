package runner

import (
	"strconv"
	"strings"
)

func withSudo(args []string) []string {
	var escaped []string

	for _, arg := range args {
		escaped = append(escaped, strconv.Quote(arg))
	}

	return []string{"runas", "/user:administrator", strconv.Quote(strings.Join(escaped, " "))}
}
