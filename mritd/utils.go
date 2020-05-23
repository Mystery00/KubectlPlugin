package mritd

import (
	"fmt"
	"github.com/mritd/promptx"
	"text/template"
)

func ToString(i int) string {
	return fmt.Sprintf("%d", i)
}

func FormatLeftLength(len int, s string) string {
	return fmt.Sprintf("%-"+ToString(len)+"s", s)
}

func FormatRightLength(len int, s string) string {
	return fmt.Sprintf("%"+ToString(len)+"s", s)
}

var FuncMap = template.FuncMap{
	"black":     promptx.Styler(promptx.FGBlack),
	"red":       promptx.Styler(promptx.FGRed),
	"green":     promptx.Styler(promptx.FGGreen),
	"yellow":    promptx.Styler(promptx.FGYellow),
	"blue":      promptx.Styler(promptx.FGBlue),
	"magenta":   promptx.Styler(promptx.FGMagenta),
	"cyan":      promptx.Styler(promptx.FGCyan),
	"white":     promptx.Styler(promptx.FGWhite),
	"bgBlack":   promptx.Styler(promptx.BGBlack),
	"bgRed":     promptx.Styler(promptx.BGRed),
	"bgGreen":   promptx.Styler(promptx.BGGreen),
	"bgYellow":  promptx.Styler(promptx.BGYellow),
	"bgBlue":    promptx.Styler(promptx.BGBlue),
	"bgMagenta": promptx.Styler(promptx.BGMagenta),
	"bgCyan":    promptx.Styler(promptx.BGCyan),
	"bgWhite":   promptx.Styler(promptx.BGWhite),
	"bold":      promptx.Styler(promptx.FGBold),
	"faint":     promptx.Styler(promptx.FGFaint),
	"italic":    promptx.Styler(promptx.FGItalic),
	"underline": promptx.Styler(promptx.FGUnderline),
	"lLength":   FormatLeftLength,
	"rLength":   FormatRightLength,
}
