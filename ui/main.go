// main.go
package main

import (
	"github.com/catsalt/wordstat/stat"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

var mainwin *ui.Window

func fileRead() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	vbox.Append(ui.NewHorizontalSeparator(), false)

	btOpenA := ui.NewButton("目录A")
	btOpenB := ui.NewButton("目录B")
	btClearA := ui.NewButton("清空A")
	btClearB := ui.NewButton("清空B")
	btStat := ui.NewButton("统计A")
	btCompare := ui.NewButton("B对比A")
	btGrade := ui.NewButton("B分级A")
	btSave := ui.NewButton("目录输出:")

	cbTidyA := ui.NewCheckbox("整理A")
	cbTidyB := ui.NewCheckbox("整理B")
	cbGroup := ui.NewCheckbox("分组")
	entryC := ui.NewEntry()

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	hbox.Append(btOpenA, false)
	hbox.Append(btClearA, false)
	hbox.Append(cbTidyA, false)
	hbox.Append(cbGroup, false)
	hbox.Append(btStat, false)
	hbox.Append(btSave, false)
	hbox.Append(entryC, false)

	vbox.Append(hbox, false)
	vbox.Append(ui.NewHorizontalSeparator(), false)

	hbox = ui.NewHorizontalBox()
	hbox.SetPadded(true)
	hbox.Append(btOpenB, false)
	hbox.Append(btClearB, false)
	hbox.Append(cbTidyB, false)
	hbox.Append(btCompare, false)
	hbox.Append(btGrade, false)

	vbox.Append(hbox, false)
	vbox.Append(ui.NewHorizontalSeparator(), false)

	mEntryA, mEntryB := ui.NewMultilineEntry(), ui.NewMultilineEntry()
	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	entryForm.Append("A-", mEntryA, true)
	entryForm.Append("B-", mEntryB, true)
	group := ui.NewGroup("目录")
	group.SetMargined(true)
	group.SetChild(entryForm)

	vbox.Append(group, true)
	vbox.Append(ui.NewHorizontalSeparator(), false)

	pbar := ui.NewProgressBar()
	hbox = ui.NewVerticalBox()
	hbox.SetPadded(true)
	hbox.Append(pbar, true)

	vbox.Append(hbox, false)
	// below for caculating

	var filesA, filesB string
	cbTidyA.SetChecked(true)
	filesC := "F:/outX Stat"

	btOpenA.OnClicked(func(*ui.Button) {
		file := ui.OpenFile(mainwin)
		if file != "" {
			filesA += file + "\r\n"
			mEntryA.SetText(filesA)
		}
	})
	btOpenB.OnClicked(func(*ui.Button) {
		file := ui.OpenFile(mainwin)
		if file != "" {
			filesB += file + "\r\n"
			mEntryB.SetText(filesB)
		}
	})
	btClearA.OnClicked(func(*ui.Button) {
		filesA = ""
		pbar.SetValue(0)
		mEntryA.SetText(filesA)
	})
	btClearB.OnClicked(func(*ui.Button) {
		filesB = ""
		pbar.SetValue(0)
		mEntryB.SetText(filesB)
	})
	btStat.OnClicked(func(*ui.Button) {
		if filesA != "" {
			pbar.SetValue(0)
			pbar.SetValue(-1)
			stat.ZsTidy(filesA, filesC, cbGroup.Checked(), cbTidyA.Checked())
			pbar.SetValue(100)
		}
	})
	btCompare.OnClicked(func(*ui.Button) {
		switch {
		case filesA == filesB:
		case filesA == "":
		case filesB == "":
		default:
			pbar.SetValue(0)
			pbar.SetValue(-1)
			stat.ZsCompare(filesA, filesB, filesC, cbTidyA.Checked(), cbTidyB.Checked())
			pbar.SetValue(100)
		}
	})
	btGrade.OnClicked(func(*ui.Button) {
		switch {
		case filesA == filesB:
		case filesA == "":
		case filesB == "":
		default:
			pbar.SetValue(0)
			pbar.SetValue(-1)
			stat.ZsGrade(filesA, filesB, filesC, cbTidyA.Checked(), cbTidyB.Checked())
			pbar.SetValue(100)
		}
	})
	// entryC.SetText(filesC)
	btSave.OnClicked(func(*ui.Button) {
		entryC.SetText(filesC)
		entryC.OnChanged(func(*ui.Entry) {
			filesC = entryC.Text()
		})
		// fmt.Println(filesC)
	})
	return vbox
}

func setupUI() {
	mainwin = ui.NewWindow("英语单词统计工具", 720, 540, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})
	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)
	tab.Append("问题反馈:", fileRead())
	tab.SetMargined(0, true)
	mainwin.Show()
}

func main() {
	ui.Main(setupUI)
}
