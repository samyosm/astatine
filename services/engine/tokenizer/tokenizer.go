package tokenizer

import "regexp"

func Wordpunkt(text string) []string {
	r, _ := regexp.Compile(`\w+|[^\w\s]+`)
	return r.FindAllString(text, -1)
}
