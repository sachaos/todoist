package main

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"time"
)

func charWidth(c rune) int {
	w := runewidth.RuneWidth(c)
	if w == 0 || w == 2 && runewidth.IsAmbiguousWidth(c) {
		w = 1
	}
	return w
}

func stringWidth(s string) int {
	w := 0
	for _, c := range s {
		w += charWidth(c)
	}
	return w
}

func drawLine(x, y int, str string) {
	color := termbox.ColorDefault
	backgroundColor := termbox.ColorDefault

	w := x
	for _, s := range str {
		termbox.SetCell(w, y, s, color, backgroundColor)
		w = w + charWidth(s)
	}
}

func draw(sync todoist.Sync, c *cli.Context) {
	// w, h := termbox.Size()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	itemList := makeList(sync, c)
	for i, strings := range itemList {
		var string string
		for _, str := range strings {
			string = string + " " + str
		}
		drawLine(0, i, string)
	}

	termbox.Flush()
}

func TUI(sync todoist.Sync, c *cli.Context) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	draw(sync, c)
loop:
	for {
		select {
		case ev := <-event_queue:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				break loop
			}
		default:
			draw(sync, c)
			time.Sleep(10 * time.Millisecond)
		}
	}
}
