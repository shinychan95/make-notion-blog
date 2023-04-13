package markdown

import (
	"fmt"
)

func InlineCode(text string) string {
	return fmt.Sprintf("`%s`", text)
}

func Bold(text string) string {
	return fmt.Sprintf("**%s**", text)
}

func Italic(text string) string {
	return fmt.Sprintf("_%s_", text)
}

func Strikethrough(text string) string {
	return fmt.Sprintf("~~%s~~", text)
}

func Equation(text string) string {
	return fmt.Sprintf("$%s$", text)
}

func Underline(text string) string {
	return fmt.Sprintf("<u>%s</u>", text)
}

func Link(text, href string) string {
	return fmt.Sprintf("[%s](%s)", text, href)
}
