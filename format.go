package main

import (
	"fmt"
	"image/color"
)

var (
	ColorText       = color.RGBA{0xd4, 0xe5, 0xf7, 0xff}
	ColorInactive   = color.RGBA{0x2f, 0x34, 0x3f, 0xff}
	ColorActive     = color.RGBA{0x62, 0x6b, 0x82, 0xff}
	ColorHighlight  = color.RGBA{0x2f, 0x6e, 0xaf, 0xff}
	ColorHighlight2 = color.RGBA{0x25, 0xB2, 0x84, 0xff}
	ColorHighlight3 = color.RGBA{0x5e, 0x57, 0xba, 0xff}
	ColorError      = color.RGBA{0xff, 0, 0, 0xff}
	ColorLove       = color.RGBA{0xef, 0x09, 0x46, 0xff}
	ColorWarning    = color.RGBA{0xff, 0xd1, 0x35, 0xff}
	ColorDanger     = color.RGBA{0xff, 0x79, 0x35, 0xff}
)

func iconC(c string, color color.Color) string {
	return fmt.Sprintf(`<span color="%s" font_family="FontAwesome">%s</span>`, colorFmt(color), string(c))
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

func colorize(c color.Color, msg string) string {
	return fmt.Sprintf(`<span color="%s">%s</span>`, colorFmt(c), msg)
}

func errMsg(s string) string {
	return colorize(ColorError, fmt.Sprintf("<%s>", s))
}

func colorFmt(c color.Color) string {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return fmt.Sprintf("#%.2x%.2x%.2x", rgba.R, rgba.G, rgba.B)
}

func elipsis(s string, l int) string {
	sr := []rune(s)
	if len(sr) <= l {
		return s
	}

	return string(sr[:l]) + "…"
}

type ColorSlideStop struct {
	At    float64
	Color color.Color
}

func colorSlide(v float64, stops ...ColorSlideStop) color.Color {
	if len(stops) == 0 {
		return ColorText
	}

	if v < stops[0].At {
		return stops[0].Color
	}
	if v > stops[len(stops)-1].At {
		return stops[len(stops)-1].Color
	}

	var left, right ColorSlideStop

	for i, s := range stops {
		if v >= s.At {
			left = s
			right = stops[i+1]
		}
	}

	r := (v - left.At) / (right.At - left.At)
	leftrgb := color.RGBAModel.Convert(left.Color).(color.RGBA)
	rightrgb := color.RGBAModel.Convert(right.Color).(color.RGBA)

	return color.RGBA{
		uint8(float64(leftrgb.R) + (float64(rightrgb.R)-float64(leftrgb.R))*r),
		uint8(float64(leftrgb.G) + (float64(rightrgb.G)-float64(leftrgb.G))*r),
		uint8(float64(leftrgb.B) + (float64(rightrgb.B)-float64(leftrgb.B))*r),
		0xff,
	}
}
