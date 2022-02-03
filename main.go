package main

import (
	//"fyne.io/fyne/app"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	fyne2 "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

const (
	recordStartText = "Record"
	recordEndText   = "Stop recording"
)

var (
	recording          bool
	window             fyne2.Window
	record             *widget.Button
	chooseBindKey      *widget.Select
	binded             *widget.Label
	processHookChan    chan hook.Event
	mappedKeys         sync.Map
	mappedKeysTrack    []string
	rwlock             sync.RWMutex
	selectedValue      string
	selectValueOptions = []string{
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"10",
		"space",
		"w",
		"a",
		"s",
		"d",
		"q",
		"e",
		"z",
		"x",
	}
)

func selectValueHooksInit(start bool) {
	if !start {
		hook.End()
	}

	hook.Register(hook.KeyDown, []string{"shift", "r"}, func(ev hook.Event) {
		if record.Text == recordStartText {
			record.Text = recordEndText
			recordStartEvent()
			window.Content().Refresh()
		} else {
			record.Text = recordStartText
			recordEndEvent()
			window.Content().Refresh()
		}
	})

	hook.Register(hook.KeyDown, []string{"t"}, func(e hook.Event) {
		if !recording {
			return
		}
		x, y := robotgo.GetMousePos()
		log.Printf("Mouse click at x %d : y %d mapped to key %s", x, y, selectedValue)
		mappedKeys.Store(selectedValue, []int{x, y})
		mappedKeysTrack = append(mappedKeysTrack, selectedValue)
		binded.Text += fmt.Sprintf("Key [%s]: click at x %d : y %d \n", selectedValue, x, y)
		chooseBindKey.ClearSelected()
		selectValueHooksInit(false)
	})

	for _, v := range mappedKeysTrack {
		hook.Register(hook.KeyDown, []string{v}, func(e hook.Event) {
			var key string
			// since this isn't a single char, keychar comes empty
			log.Printf("keycode %d", e.Keycode)
			if e.Keychar == 32 {
				key = "space"
			} else {
				key = string(e.Keychar)
			}
			log.Printf("v %+v\n", e.Keychar)
			if v, ok := mappedKeys.Load(key); ok {
				var coords []int
				if coords, ok = v.([]int); !ok {
					log.Fatal("Failed to assert coords from store")
				}
				robotgo.MoveClick(coords[0], coords[1])
				log.Printf("%+v\n", mappedKeys)
			} else {
				log.Printf("Failed to bind, mappedKeys %+v keycode %d", mappedKeys, e.Kind)
			}
		})
	}
	processHookChan = hook.Start()
	<-hook.Process(processHookChan)
}

func recordStartEvent() {
	fmt.Println("Recording")
	recording = true
}

func recordEndEvent() {
	fmt.Println("Stopped recording")
	recording = false
}

func init() {

}

func main() {
	m := make(map[string][]int)
	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, os.Interrupt)
	if b, err := os.ReadFile("key_bindings.json"); err == nil {
		if err = json.Unmarshal(b, &m); err == nil {
			for k, v := range m {
				mappedKeys.Store(k, v)
				mappedKeysTrack = append(mappedKeysTrack, k)
			}
		} else {
			log.Println(err.Error())
		}
	} else {
		log.Println(err.Error())
	}

	// UI init
	a := app.New()
	window = a.NewWindow("Key to click")
	binded = widget.NewLabel("")
	text := widget.NewLabel(`
This app maps mouse clicks at position to specified keys.
[shift] + [r] keys to start/stop recording.
[t] key while keeping the mouse in the right position to bind to the selected key from the menu.
	`)

	record = widget.NewButton(recordStartText, func() {
		if record.Text == recordStartText {
			record.Text = recordEndText
			recordStartEvent()
			window.Content().Refresh()
		} else {
			record.Text = recordStartText
			recordEndEvent()
			window.Content().Refresh()
		}
	})

	chooseBindKey = widget.NewSelect(selectValueOptions,
		func(value string) {
			log.Println("Select set to", value)
			rwlock.Lock()
			selectedValue = value
			rwlock.Unlock()
		})

	container := container.NewVBox(text, record, chooseBindKey, binded)
	window.SetContent(container)

	// Key listen
	go func() {
		selectValueHooksInit(true)
	}()
	// UI loop start
	window.Content().Refresh()
	window.ShowAndRun()
	<-sigInt
	for _, v := range mappedKeysTrack {
		i, ok := mappedKeys.Load(v)
		if ok {
			if coords, ok := i.([]int); ok {
				m[v] = coords
			}
		}
	}
	if b, err := json.Marshal(m); err == nil {
		os.WriteFile("key_bindings.json", b, os.ModePerm)
	}
}
