package tools

import (
	"strings"
	"regexp"
)

func RemoveHtmlTags(content string) string {
	content = strings.TrimSpace(content)
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	content = re.ReplaceAllStringFunc(content, strings.ToLower)
	reg := regexp.MustCompile(`<!--[^>]+>|<iframe[\S\s]+?</iframe>|<a[^>]+>|</a>|<script[\S\s]+?</script>`)
	return reg.ReplaceAllString(content, "[invalid]")
}
