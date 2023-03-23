package generator

import (
	"strings"

	"github.com/Felixoid/braxpansion"
)

func expandString(expandableName string, variables map[string]string) []string {
	if len(variables) > 0 {
		s, next, ok := strings.Cut(expandableName, "{{")
		if ok {
			var buf strings.Builder
			buf.WriteString(s)
			for next != "" {
				var sV string
				if sV, next, ok = strings.Cut(next, "}}"); ok {
					sV = strings.TrimSpace(sV)
					if v, ok := variables[sV]; ok {
						s = v
					}
				}
				buf.WriteString(s)
				if s, next, ok = strings.Cut(next, "{{"); !ok {
					buf.WriteString(s)
				}
			}

			return braxpansion.ExpandString(buf.String())
		}
	}

	return braxpansion.ExpandString(expandableName)
}
