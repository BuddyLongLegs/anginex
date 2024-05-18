package analytics

import (
	"fmt"
	"regexp"
	"strings"
)

func ReplaceAuthInfo(input string) string {
	re := regexp.MustCompile(`:\/\/[^@]*@`)
	return re.ReplaceAllStringFunc(input, func(m string) string {
		return "://" + "analytics:password" + "@"
	})
}

func ReplaceValsInSQLQuery(query string, index int, replacement string) string {
	// replace $index with replacement
	key := "$" + fmt.Sprintf("%d", index)
	query = strings.ReplaceAll(query, key, replacement)
	return query
}
