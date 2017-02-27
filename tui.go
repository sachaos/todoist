package main

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"github.com/sachaos/todoist/lib"
	"github.com/urfave/cli"
	"log"
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

func draw(sync todoist.Sync, c *cli.Context, baseY, ptr int) int {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

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

	ptrY := 0
	y := 0
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
			drawLine(0, y-baseY, currentProject.Name+" Tasks")
			drawLine(0, y+1-baseY, "---")
			y = y + 2
		}
		if ptr == i {
			ptrY = y - baseY
		}
		y = y + drawItem(ptr, i, y-baseY, paddings, sync, c, itemList[i], order)
	}

	termbox.Flush()

	return ptrY
}

func TUI(sync todoist.Sync, c *cli.Context) {
	var ptrY int
	ptr := 0
	baseY := 0
	err := termbox.Init()
	_, iH := termbox.Size()
	pageSize := iH / 2
	itemCount := len(sync.ItemOrders)
	log.Printf("start pageSize: %d", pageSize)

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

	ptrY = draw(sync, c, baseY, ptr)
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
				}
				if ev.Key == termbox.KeyCtrlD {
					ptr += pageSize
				}
				if ev.Key == termbox.KeyCtrlU {
					ptr -= pageSize
				}
				if ptr < 0 {
					ptr = 0
				}
				if ptr >= itemCount {
					ptr = itemCount - 1
				}
			}
		default:
			ptrY = draw(sync, c, baseY, ptr)
			_, h := termbox.Size()

			if ptrY >= h {
				baseY = baseY + pageSize
			} else if ptrY < 0 {
				baseY = baseY - pageSize
			}

			log.Printf("h: %d, ptrY: %d, baseY: %d", h, ptrY, baseY)

			time.Sleep(10 * time.Millisecond)
		}
	}
}
