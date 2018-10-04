package main

import (
	"fmt"
)

const (
	ColorInactive   = "#2f343f"
	ColorActive     = "#626b82"
	ColorHighlight  = "#2f6eaf"
	ColorHighlight2 = "#25B284"
	ColorHighlight3 = "#5e57ba"
)

func iconC(c string, color string) string {
	return fmt.Sprintf(`<span color="%s" font_family="FontAwesome">%s</span>`, color, string(c))
}

func icon(c string) string {
	return fmt.Sprintf(`<span font_family="FontAwesome">%s</span>`, string(c))
}

func font(font string, msg string) string {
	return fmt.Sprintf(`<span font_family="%s">%s</span>`, font, msg)
}

func bold(msg string) string {
	return fmt.Sprintf("<b>%s</b>", msg)
}

func color(c string, msg string) string {
	return fmt.Sprintf(`<span color="%s">%s</span>`, c, msg)
}

func errMsg(s string) string {
	return color("#ff0000", fmt.Sprintf("<%s>", s))
}

func elipsis(s string, l int) string {
	sr := []rune(s)
	if len(sr) <= l {
		return s
	}

	return string(sr[:l]) + "…"
}
