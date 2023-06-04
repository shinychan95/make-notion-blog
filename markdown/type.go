package markdown

import (
	"fmt"
	"strings"
)

func Header(indent, text string) string {
	return fmt.Sprintf("%s# %s\n\n", indent, text)
}

func SubHeader(indent, text string) string {
	return fmt.Sprintf("%s## %s\n\n", indent, text)
}

func SubSubHeader(indent, text string) string {
	return fmt.Sprintf("%s### %s\n", indent, text)
}

func Text(indent, text string) string {
	if text == "" {
		text = "<br/>\n\n"
	}

	return fmt.Sprintf("%s%s\n\n", indent, text)
}

func Code(indent, lang, text string) string {
	return fmt.Sprintf("%s```%s\n%s%s\n%s```\n\n", indent, lang, indent, text, indent)
}

func Divider(indent string) string {
	return fmt.Sprintf("%s\n---\n\n", indent)
}

func BulletedList(indent, text string) string {
	return fmt.Sprintf("%s- %s\n\n", indent, text)
}

func NumberedList(indent string, number uint8, text string) string {
	return fmt.Sprintf("%s%d. %s\n\n", indent, number, text)
}

func Toggle(indent, text, content string) string {
	summary := fmt.Sprintf("%s<summary>%s</summary>\n", indent, text)
	return fmt.Sprintf("%s<details>\n%s%s%s</details>\n\n", indent, summary, content, indent)
}

func Quote(indent, text string) string {
	text = strings.Replace(text, "\n", "<br/>", -1)
	return fmt.Sprintf("%s> %s\n\n", indent, text)
}

func Callout(indent, text string) string {
	return fmt.Sprintf("%s> 🦖 %s\n\n", indent, text)
}

func Image(indent, imagePath string) string {
	return fmt.Sprintf("%s![](%s)\n", indent, imagePath)
}

func ToDo(indent, text string, checked bool) string {
	if checked {
		return fmt.Sprintf("%s- [x] %s\n", indent, text)
	}
	return fmt.Sprintf("%s- [ ] %s\n", indent, text)
}
