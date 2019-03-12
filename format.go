package main

import (
	"unicode/utf8"

	"github.com/lucasb-eyer/go-colorful"
)

type Part struct {
	Text   string
	Sat    float64
	Lum    float64
	Styles []Pango
}

func IconPart(icon string) Part {
	return Part{
		Text:   icon,
		Sat:    .7,
		Lum:    .6,
		Styles: []Pango{FontIcon},
	}
}

func TextPart(text string, styles ...Pango) Part {
	return Part{
		Text:   text,
		Sat:    .4,
		Lum:    .85,
		Styles: styles,
	}
}

func Render(hstart, hend float64, parts []Part) string {
	res := ""
	length := 0

	for _, part := range parts {
		length += utf8.RuneCountInString(part.Text)
	}

	i := 0
	for _, part := range parts {
		for _, c := range part.Text {
			color := rainbow(length, part.Lum, part.Sat, hstart, hend, i)
			res += Style(string(c), append(part.Styles, color)...)
			i++
		}
	}

	return res
}

func rainbow(length int, lum, sat, hstart, hend float64, pos int) Color {
	span := 0.0
	if hend > hstart {
		span = hend - hstart
	} else {
		span = 360 - hstart + hend
	}
	step := span / float64(length)
	hue := (step*float64(pos) + hstart)
	if hue > 360 {
		hue -= 360
	}

	color := colorful.Hsl(hue, sat, lum)
	r, g, b := color.RGB255()

	return Color{r, g, b}
}
