package helper

import (
	"regexp"
	"strings"
)

func ReplaceNumberToWord(wordWithNumber string, wordToReplace string) string {
	re := regexp.MustCompile(`\d+`)
	return re.ReplaceAllString(wordWithNumber, wordToReplace)
}

func NameToGmail(name string) string {
	withoutSpace := strings.ReplaceAll(name, " ", "")
	return withoutSpace + "gmail.com"
}
