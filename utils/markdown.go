package utils

import (
	"fmt"
	"strings"
)

func inlineCode(text string) string {
	return fmt.Sprintf("`%s`", text)
}

func bold(text string) string {
	return fmt.Sprintf("**%s**", text)
}

func italic(text string) string {
	return fmt.Sprintf("_%s_", text)
}

func strikethrough(text string) string {
	return fmt.Sprintf("~~%s~~", text)
}

func equation(text string) string {
	return fmt.Sprintf("$%s$", text)
}

func underline(text string) string {
	return fmt.Sprintf("<u>%s</u>", text)
}

func link(text, href string) string {
	return fmt.Sprintf("[%s](%s)", text, href)
}

func codeBlock(text, language string) string {
	if language == "plain text" {
		language = "text"
	}

	return fmt.Sprintf("```%s\n%s\n```", language, text)
}

func heading1(text string) string {
	return fmt.Sprintf("# %s", text)
}

func heading2(text string) string {
	return fmt.Sprintf("## %s", text)
}

func heading3(text string) string {
	return fmt.Sprintf("### %s", text)
}

func quote(text string) string {
	return fmt.Sprintf("> %s", strings.ReplaceAll(text, "\n", "  \n> "))
}

func callout(text string, icon string) string {
	var emoji string
	if icon != "" {
		emoji = icon
	}

	return fmt.Sprintf("> %s%s", emoji, strings.ReplaceAll(text, "\n", "  \n> "))
}

func bullet(text string, count int) string {
	renderText := strings.TrimSpace(text)
	if count > 0 {
		return fmt.Sprintf("%d. %s", count, renderText)
	}
	return fmt.Sprintf("- %s", renderText)
}

func todo(text string, checked bool) string {
	if checked {
		return fmt.Sprintf("- [x] %s", text)
	}
	return fmt.Sprintf("- [ ] %s", text)
}

func image(alt, href string) string {
	return fmt.Sprintf("![%s](%s)", alt, href)
}

func addTabSpace(text string, n int) string {
	tab := "	"
	for i := 0; i < n; i++ {
		if strings.Contains(text, "\n") {
			multiLineText := strings.Join(strings.SplitAfter(text, "\n"), tab)
			text = tab + multiLineText
		} else {
			text = tab + text
		}
	}
	return text
}

func divider() string {
	return "---"
}

func toggle(summary, children string) string {
	if summary == "" {
		return children
	}
	return fmt.Sprintf("<details>\n  <summary>%s</summary>\n%s\n</details>", summary, children)
}

// TODO - Table ([][]string 입력)
