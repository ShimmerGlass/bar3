package main

import (
	"fmt"
	"html"
	"strings"
)

type Pango interface {
	Pango() string
}

type Color struct {
	R, G, B uint8
}

func (c Color) String() string {
	return fmt.Sprintf("#%.2x%.2x%.2x", c.R, c.G, c.B)
}
func (c Color) Pango() string {
	return fmt.Sprintf("color=\"%s\"", c.String())
}

var (
	ColorText       = Color{0xd4, 0xe5, 0xf7}
	ColorInactive   = Color{0x2f, 0x34, 0x3f}
	ColorActive     = Color{0x62, 0x6b, 0x82}
	ColorHighlight  = Color{0xff, 0x1a, 0x67}
	ColorHighlight2 = Color{0x25, 0xB2, 0x84}
	ColorHighlight3 = Color{0x9c, 0x00, 0xff}
	ColorError      = Color{0xff, 0, 0}
	ColorLove       = Color{0xef, 0x09, 0x46}
	ColorWarning    = Color{0xff, 0xd1, 0x35}
	ColorDanger     = Color{0xff, 0x79, 0x35}
)

type Font string

func (f Font) Pango() string {
	return fmt.Sprintf("font_family=\"%s\"", f)
}

const (
	FontIcon Font = "NotoSansDisplay Nerd Font Mono"
	FontMono Font = "Noto Mono"
)

type FontSize string

func (f FontSize) Pango() string {
	return fmt.Sprintf("font_size=\"%s\"", f)
}

const (
	FontSizeSmall  FontSize = "small"
	FontSizeMedium FontSize = "medium"
	FontSizeLarge  FontSize = "large"
)

func Style(txt string, opts ...Pango) string {
	p := "<span"
	for _, o := range opts {
		p += " " + o.Pango()
	}
	return p + ">" + html.EscapeString(txt) + "</span>"
}

func Icon(code string, opts ...Pango) string {
	return Style(code, append(opts, FontIcon)...)
}

func Error(s string) string {
	return Style(fmt.Sprintf("<%s>", s), ColorError)
}

func Comb(s ...string) string {
	return strings.Join(s, "")
}

func Elipsis(s string, l int) string {
	sr := []rune(s)
	if len(sr) <= l {
		return s
	}

	return string(sr[:l]) + "…"
}

type GradStop struct {
	At    float64
	Color Color
}

func Grad(v float64, stops ...GradStop) Color {
	if len(stops) == 0 {
		return ColorText
	}

	if v <= stops[0].At {
		return stops[0].Color
	}
	if v >= stops[len(stops)-1].At {
		return stops[len(stops)-1].Color
	}

	var left, right GradStop

	for i, s := range stops {
		if v >= s.At {
			left = s
			right = stops[i+1]
		}
	}

	r := (v - left.At) / (right.At - left.At)

	return Color{
		uint8(float64(left.Color.R) + (float64(right.Color.R)-float64(left.Color.R))*r),
		uint8(float64(left.Color.G) + (float64(right.Color.G)-float64(left.Color.G))*r),
		uint8(float64(left.Color.B) + (float64(right.Color.B)-float64(left.Color.B))*r),
	}
}
