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

type ItemView struct {
	W         int
	H         int
	BaseY     int
	Ptr       int
	PtrY      int
	ItemCount int
	PageSize  int
}

func createItemView() (*ItemView, error) {
	err := termbox.Init()
	if err != nil {
		return nil, err
	}
	w, h := termbox.Size()
	iv := &ItemView{
		W:         w,
		H:         h,
		BaseY:     0,
		Ptr:       0,
		PtrY:      0,
		ItemCount: 0,
		PageSize:  h / 2,
	}
	return iv, nil
}

func (iv *ItemView) SetItemCount(count int) *ItemView {
	iv.ItemCount = count
	return iv
}

func (iv *ItemView) MovePtrTo(count int) *ItemView {
	iv.Ptr = count
	if iv.Ptr < 0 {
		iv.Ptr = 0
	}
	if iv.Ptr >= iv.ItemCount {
		iv.Ptr = iv.ItemCount - 1
	}
	return iv
}

func (iv *ItemView) MovePtr(count int) *ItemView {
	iv.Ptr = iv.Ptr + count
	if iv.Ptr < 0 {
		iv.Ptr = 0
	}
	if iv.Ptr >= iv.ItemCount {
		iv.Ptr = iv.ItemCount - 1
	}
	return iv
}

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

func (iv *ItemView) adjust() {
	iv.W, iv.H = termbox.Size()
	if iv.PtrY >= iv.H {
		iv.BaseY = iv.BaseY + iv.PageSize
	} else if iv.PtrY < 0 {
		iv.BaseY = iv.BaseY - iv.PageSize
	}
}

func (iv *ItemView) draw(sync todoist.Sync, c *cli.Context) {
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
			drawLine(0, y-iv.BaseY, currentProject.Name+" Tasks")
			drawLine(0, y+1-iv.BaseY, "---")
			y = y + 2
		}
		if iv.Ptr == i {
			iv.PtrY = y - iv.BaseY
		}
		y = y + drawItem(iv.Ptr, i, y-iv.BaseY, paddings, sync, c, itemList[i], order)
	}

	termbox.Flush()
}

func TUI(sync todoist.Sync, c *cli.Context) {
	iv, err := createItemView()
	if err != nil {
		panic(err)
	}
	iv.SetItemCount(len(sync.ItemOrders))
	// log.Printf("start pageSize: %d", iv.PageSize)

	defer termbox.Close()

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	iv.draw(sync, c)
loop:
	for {
		select {
		case ev := <-event_queue:
			if ev.Type == termbox.EventKey {
				if ev.Key == termbox.KeyEsc || ev.Ch == 'q' {
					break loop
				}
				if ev.Ch == 'j' {
					iv.MovePtr(1)
				}
				if ev.Ch == 'k' {
					iv.MovePtr(-1)
				}
				if ev.Key == termbox.KeyCtrlD {
					iv.MovePtr(iv.PageSize)
				}
				if ev.Key == termbox.KeyCtrlU {
					iv.MovePtr(-iv.PageSize)
				}
			}
		default:
			iv.draw(sync, c)
			iv.adjust()

			log.Printf("h: %d, ptrY: %d, baseY: %d", iv.H, iv.PtrY, iv.BaseY)

			time.Sleep(10 * time.Millisecond)
		}
	}
}
