// fui
package main

import (
	"fmt"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	fmt.Println("Hello World!")
	a := app.New()
	w := a.NewWindow("hello")
	btOpenA := widget.NewButton("pathA", func() {
		file := ui.OpenFile(mainwin)
		if file != "" {
			filesA += file + "\r\n"
			mEntryA.SetText(filesA)
		}
	})

	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("quit", func() { a.Quit() }),
	))
	w.ShowAndRun()
}
