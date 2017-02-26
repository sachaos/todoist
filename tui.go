package main

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"strconv"
	"strings"
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

func drawItem(ptr int, num int, y int, paddings []int, sync todoist.Sync, c *cli.Context, strings []string, order todoist.ItemOrder) (height int) {
	x := 0
	// item := order.Data.(todoist.Item)
	fg := termbox.ColorDefault
	bg := termbox.ColorDefault
	if ptr == num {
		bg = termbox.ColorCyan
	}
	for i, string := range strings {
		for _, c := range string {
			termbox.SetCell(x, y, c, fg, bg)
			x = x + charWidth(c)
		}
		for ; x < paddings[i+1]; x++ {
			termbox.SetCell(x, y, ' ', fg, bg)
		}
	}
	return 1
}

func TUIMakeList(sync todoist.Sync, c *cli.Context) [][]string {
	itemList := [][]string{}
	for _, itemOrder := range sync.ItemOrders {
		item := itemOrder.Data.(todoist.Item)
		itemList = append(itemList, []string{
			strconv.Itoa(item.ID),
			"p" + strconv.Itoa(item.Priority),
			strings.Repeat("    ", itemOrder.Indent-1) + ContentFormat(item),
		})
	}
	return itemList
}

func draw(sync todoist.Sync, c *cli.Context, ptr int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	y := 0

	itemList := TUIMakeList(sync, c)
	maxWidth := make([]int, len(itemList[0]))
	for i := 0; i < len(itemList[0]); i++ {
		for _, strings := range itemList {
			sw := stringWidth(strings[i])
			if maxWidth[i] < sw {
				maxWidth[i] = sw
			}
		}
	}

	paddings := make([]int, len(itemList[0])+1)
	for i := 0; i < len(maxWidth); i++ {
		paddings[i+1] = paddings[i] + maxWidth[i] + 1
	}

	var currentProject todoist.Project
	for i, order := range sync.ItemOrders {
		item := order.Data.(todoist.Item)
		if currentProject.ID != item.ProjectID {
			project, err := todoist.SearchByID(sync.Projects, item.ProjectID)
			if err != nil {
				panic(err)
			}
			currentProject = project.(todoist.Project)
			y = y + 1
			drawLine(0, y, currentProject.Name+" Tasks")
			drawLine(0, y+1, "---")
			y = y + 2
		}
		y = y + drawItem(ptr, i, y, paddings, sync, c, itemList[i], order)
	}

	termbox.Flush()
}

func TUI(sync todoist.Sync, c *cli.Context) {
	ptr := 0
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

	draw(sync, c, ptr)
loop:
	for {
		select {
		case ev := <-event_queue:
			if ev.Type == termbox.EventKey {
				if ev.Key == termbox.KeyEsc || ev.Ch == 'q' {
					break loop
				}
				if ev.Ch == 'j' {
					ptr += 1
				}
				if ev.Ch == 'k' {
					ptr -= 1
					if ptr < 0 {
						ptr = 0
					}
				}
			}
		default:
			draw(sync, c, ptr)
			time.Sleep(10 * time.Millisecond)
		}
	}
}
