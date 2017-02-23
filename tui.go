package main

import (
	"github.com/jroimartin/gocui"
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"log"
)

func layout(sync todoist.Sync, c *cli.Context, g *gocui.Gui) error {

	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		writer = NewTSVWriter(v)
		itemList := makeList(sync, c)

		for _, strings := range itemList {
			writer.Write(strings)
		}

		writer.Flush()
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func TUI(sync todoist.Sync, c *cli.Context) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(func(g *gocui.Gui) error {
		return layout(sync, c, g)
	})

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
