package mritd

import (
	"bytes"
	"fmt"
	"github.com/mritd/promptx/list"
	"github.com/mritd/promptx/utils"
	"github.com/mritd/readline"
	"os"
	"strings"
	"text/template"
)

const (
	DefaultActiveTpl       = "{{ . | cyan }}"
	DefaultInactiveTpl     = "{{ . | white }}"
	DefaultSelectedTpl     = "{{ . | cyan }}"
	DefaultSelectHeaderTpl = "{{ \"Use the arrow keys to navigate: ↓ ↑ → ←\" | faint }}"
	DefaultSelectPromptTpl = "{{ \"Select\" | faint }} {{ . | faint}}:"
	DefaultDisPlaySize     = 5
	DefaultShowWidth       = 20
	NewLine                = "\n"
)

type Select struct {
	Config *SelectConfig
	Items  interface{}
	buf    bytes.Buffer
	height int

	selectPrompt *template.Template
	selectHeader *template.Template
	selected     *template.Template
	active       *template.Template
	inactive     *template.Template
}

type SelectConfig struct {
	ActiveTpl       string
	InactiveTpl     string
	SelectedTpl     string
	DisPlaySize     int
	SelectPrompt    string
	SelectHeaderTpl string
	SelectPromptTpl string
	ShowWidth       int
	ShowBorder      bool
}

const (
	CENTER      = "═"
	LeftTop     = "╔"
	RightTop    = "╗"
	LeftCenter  = "╠"
	RightCenter = "╣"
	LeftBottom  = "╚"
	RightBottom = "╝"
)

func (s *Select) prepareTemplates() {

	var err error

	// set default value
	if s.Config.SelectHeaderTpl == "" {
		s.Config.SelectHeaderTpl = DefaultSelectHeaderTpl
	}
	if s.Config.SelectPromptTpl == "" {
		s.Config.SelectPromptTpl = DefaultSelectPromptTpl
	}
	if s.Config.SelectedTpl == "" {
		s.Config.SelectedTpl = DefaultSelectedTpl
	}
	if s.Config.ActiveTpl == "" {
		s.Config.ActiveTpl = DefaultActiveTpl
	}
	if s.Config.InactiveTpl == "" {
		s.Config.InactiveTpl = DefaultInactiveTpl
	}
	if s.Config.DisPlaySize < 1 {
		s.Config.DisPlaySize = DefaultDisPlaySize
	}
	if s.Config.ShowWidth < 1 {
		s.Config.ShowWidth = DefaultShowWidth
	}

	// Select prepare
	s.selectHeader, err = template.New("").Funcs(FuncMap).Parse(s.Config.SelectHeaderTpl + NewLine)
	utils.CheckAndExit(err)
	s.selectPrompt, err = template.New("").Funcs(FuncMap).Parse(s.Config.SelectPromptTpl + NewLine)
	utils.CheckAndExit(err)
	s.selected, err = template.New("").Funcs(FuncMap).Parse(s.Config.SelectedTpl)
	utils.CheckAndExit(err)
	s.active, err = template.New("").Funcs(FuncMap).Parse(s.Config.ActiveTpl + NewLine)
	utils.CheckAndExit(err)
	s.inactive, err = template.New("").Funcs(FuncMap).Parse(s.Config.InactiveTpl + NewLine)
	utils.CheckAndExit(err)
}

func (s *Select) writeData(l *list.List) {
	// clean buffer
	s.buf.Reset()

	// clean terminal
	for i := 0; i < s.height; i++ {
		s.buf.WriteString(moveUp)
		s.buf.WriteString(clearLine)
	}

	if s.Config.ShowBorder {
		s.buf.WriteString(fmt.Sprintln(LeftTop + strings.Repeat(CENTER, s.Config.ShowWidth+2) + RightTop))
	}
	// select header
	s.buf.Write(utils.Render(s.selectHeader, nil))

	// select prompt
	s.buf.Write(utils.Render(s.selectPrompt, s.Config.SelectPrompt))
	// print empty line
	s.buf.WriteString(fmt.Sprintln(LeftCenter + strings.Repeat(CENTER, s.Config.ShowWidth+2) + RightCenter))

	items, idx := l.Items()
	for i, item := range items {
		if i == idx {
			s.buf.Write(utils.Render(s.active, item))
		} else {
			s.buf.Write(utils.Render(s.inactive, item))
		}
	}
	if s.Config.ShowBorder {
		s.buf.WriteString(fmt.Sprintln(LeftBottom + strings.Repeat(CENTER, s.Config.ShowWidth+2) + RightBottom))
	}

	// set high
	s.height = len(strings.Split(s.buf.String(), "\n")) - 1
}

func (s *Select) Run() int {
	s.prepareTemplates()

	dataList, err := list.New(s.Items, s.Config.DisPlaySize)
	utils.CheckAndExit(err)

	l, err := readline.NewEx(&readline.Config{
		Prompt:                 "",
		DisableAutoSaveHistory: true,
		HistoryLimit:           -1,
		InterruptPrompt:        "^C",
		UniqueEditLine:         true,
		DisableBell:            true,
		Stdin:                  readline.NewCancelableStdin(os.Stdin),
	})
	utils.CheckAndExit(err)
	defer func() { _ = l.Close() }()

	filterInput := func(r rune) (rune, bool) {
		ok := false
		switch r {
		case readline.CharInterrupt:
			// show cursor
			_, _ = l.Write([]byte(showCursor))
			l.Refresh()
			return r, true
		case readline.CharEnter:
			return r, true
		case readline.CharReadLineExit:
			return r, true
		case readline.CharNext:
			dataList.Next()
			ok = true
		case readline.CharPrev:
			dataList.Prev()
			ok = true
		case readline.CharForward:
			dataList.PageDown()
			ok = true
		case readline.CharBackward:
			dataList.PageUp()
			ok = true
		case readline.CharZero:
			dataList.Go(0)
		case readline.CharOne:
			dataList.Go(1)
		case readline.CharTwo:
			dataList.Go(2)
		case readline.CharThree:
			dataList.Go(3)
		case readline.CharFour:
			dataList.Go(4)
		case readline.CharFive:
			dataList.Go(5)
		case readline.CharSix:
			dataList.Go(6)
		case readline.CharSeven:
			dataList.Go(7)
		case readline.CharEight:
			dataList.Go(8)
		case readline.CharNine:
			dataList.Go(9)
		// block other key
		default:
			return r, false
		}
		s.writeData(dataList)
		_, _ = l.Write(s.buf.Bytes())
		l.Refresh()
		return r, ok
	}

	l.Config.FuncFilterInputRune = filterInput

	// hide cursor
	_, _ = l.Write([]byte(hideCursor))

	// write data
	s.writeData(dataList)

	// write to terminal
	_, err = l.Write(s.buf.Bytes())
	utils.CheckAndExit(err)

	// read
	_, err = l.Readline()
	utils.CheckAndExit(err)

	// get select option
	items, idx := dataList.Items()
	result := items[idx]

	// clean terminal
	s.buf.Reset()
	for i := 0; i < s.height; i++ {
		s.buf.WriteString(moveUp)
		s.buf.WriteString(clearLine)
	}

	_, err = l.Write(s.buf.Bytes())
	utils.CheckAndExit(err)

	// show cursor
	_, err = l.Write([]byte(showCursor))
	utils.CheckAndExit(err)
	l.Refresh()

	fmt.Println(string(utils.Render(s.selected, result)))

	return dataList.Index()
}
