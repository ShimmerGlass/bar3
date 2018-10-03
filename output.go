package main

import (
	"encoding/json"
	"io"
)

const witerPreamble = `{"version":1,"click_events":true}
[[]
`

const MarkupPango = "pango"

type Block struct {
	FullText  string `json:"full_text"`
	Color     string `json:"color"`
	Name      string `json:"name"`
	Separator bool   `json:"separator"`
	Markup    string `json:"markup"`
}

type Writer struct {
	out     io.Writer
	started bool
}

func NewWriter(out io.Writer) *Writer {
	return &Writer{out: out}
}

func (w *Writer) Write(blocks ...Block) error {
	if !w.started {
		_, err := w.out.Write([]byte(witerPreamble))
		if err != nil {
			return err
		}
		w.started = true
	}

	_, err := w.out.Write([]byte(","))
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w.out)
	enc.SetEscapeHTML(false)
	return enc.Encode(blocks)
}
