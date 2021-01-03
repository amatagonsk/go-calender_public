package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	container2 "fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"myModule"
	"path/filepath"
	"strconv"
	"time"
)

const minCellHeight = 70
const minCellWidth = 80
const bgColorR, bgColorG, bgColorB, bgColorA uint8 = 60, 120, 60, 120
const holidayColorR, holidayColorG, holidayColorB, holidayColorA uint8 = 120, 60, 60, 120

var weeks = [7]string{"Sun", "Mun", "Tue", "Wed", "Thu", "Fri", "Sat"}
var paddingSize float32
var emptyR = canvas.NewRectangle(color.RGBA{})

type calenderStruct struct {
	top    *fyne.Container
	middle *fyne.Container
	window fyne.Window
}

func createSimpleText(str string) *canvas.Text {
	center := fyne.TextAlignCenter
	style := fyne.TextStyle{Monospace: true}
	return &canvas.Text{Text: str, Color: theme.ForegroundColor(), Alignment: center, TextStyle: style, TextSize: theme.TextSize()}
}

func createColorText(str string, color color.RGBA) *canvas.Text {
	result := createSimpleText(str)
	result.Color = color
	return result
}

func createSizedText(str string, size int) *canvas.Text {
	result := createSimpleText(str)
	result.TextSize = float32(size)
	return result
}

func createCells(monthDay int, offset int, showTime time.Time) *fyne.Container {
	cellContainer := container2.New(layout.NewGridLayoutWithColumns(7))

	for i := 1; i <= offset; i++ {
		cellContainer.Add(widget.NewLabel(""))
	}

	holidayRect := canvas.NewRectangle(color.RGBA{R: holidayColorR, G: holidayColorG, B: holidayColorB, A: holidayColorA})
	todayRect := canvas.NewRectangle(color.RGBA{R: bgColorR, G: bgColorG, B: bgColorB, A: bgColorA})
	emptyRect := canvas.NewRectangle(color.RGBA{})
	emptyRect.SetMinSize(fyne.Size{Height: minCellHeight, Width: minCellWidth})

	path := getFileAbsFilePath()

	var mds []myModule.CalenderMdInt
	if err := checkExist(path); err != nil {
		mds = []myModule.CalenderMdInt{}
	} else {
		mds = convertMdToInt(path)
	}

	holidaysOfMonth := myModule.GetJapaneseHolidays(showTime)
	for i := 1; i <= monthDay; i++ {
		item := createSimpleText(strconv.Itoa(i))
		item.TextSize = theme.TextSize() * 7 / 5

		// 日が同じで現時刻と表示月が同じかつ休日があり、日付が休日
		if i == showTime.Day() && myModule.IsNowTargetYearMonth(showTime) && myModule.IndexOfIntElement(holidaysOfMonth, i) != -1 {
			emptyR.SetMinSize(fyne.Size{Width: paddingSize, Height: paddingSize})
			wrappedToday := container2.New(layout.NewBorderLayout(emptyR, emptyR, emptyR, emptyR), todayRect, item)
			cellContainer.Add(container2.New(layout.NewMaxLayout(), emptyRect, holidayRect, wrappedToday, notifyCircle(i, mds)))
			// 休日がある
		} else if myModule.IndexOfIntElement(holidaysOfMonth, i) != -1 {
			cellContainer.Add(container2.New(layout.NewMaxLayout(), emptyRect, holidayRect, notifyCircle(i, mds), item))
			// 日付が現日時とおなじ
		} else if i == showTime.Day() && myModule.IsNowTargetYearMonth(showTime) {
			cellContainer.Add(container2.New(layout.NewMaxLayout(), emptyRect, todayRect, notifyCircle(i, mds), item))
		} else {
			cellContainer.Add(container2.New(layout.NewMaxLayout(), emptyRect, notifyCircle(i, mds), item))
		}
	}
	return cellContainer
}

func getFileAbsFilePath() string {
	return filepath.ToSlash(filepath.Clean(
		fmt.Sprintf("%v/%04d/%02d.md", OutputPath, TargetTime.Year(), int(TargetTime.Month())),
	))
}

func notifyCircle(day int, mds []myModule.CalenderMdInt) *fyne.Container {
	redCircle := &canvas.Circle{FillColor: color.RGBA{R: 255, G: 74, B: 74, A: 255}, Position1: fyne.Position{X: 0, Y: 0}}
	notifyContainer := container2.New(layout.NewGridLayout(3))

	if count := getNotifyCountFromMd(day, mds); count > 0 && len(mds) > 0 {
		var input *fyne.Container
		var circleText string
		if count > 99 {
			circleText = "99+"
		} else {
			circleText = strconv.Itoa(count)
		}

		// fixme: want top left, so fill empty
		for i := 0; i < 6; i++ {
			if i == 0 {
				sizedText := createSizedText(circleText, int(theme.TextSize()*4/5))
				sizedText.Color = color.White
				sizedText.TextStyle = fyne.TextStyle{Bold: true}
				input = container2.New(layout.NewMaxLayout(), redCircle, sizedText)
			} else {
				input = container2.New(layout.NewMaxLayout(), createSimpleText(""))
			}
			notifyContainer.Add(input)
		}
	}
	return notifyContainer
}

func inputWindow(c calenderStruct) {
	response := widget.NewEntry()
	response.Validator = validation.NewRegexp(`^\d{1,4}-[0]?[1-9]$|^\d{1,4}-[0-1][1-2]$|^\d{1,4}-10$`, "response can only numbers, and '-'")

	items := []*widget.FormItem{
		widget.NewFormItem("", response),
	}

	response.Text = fmt.Sprintf("%04d-%02d", TargetTime.Year(), TargetTime.Month())

	dialog.ShowForm(
		"Jump to YYYY-MM", "ok", "cancel", items, func(b bool) {
			if !b {
				return
			}

			inputTime, _ := myModule.ParseYearMonth(response.Text)
			if myModule.IsNowTargetYearMonth(*inputTime) {
				TargetTime = time.Now()
			} else {
				TargetTime = *inputTime
			}
			updateByInput(TargetTime, c)
		},
		c.window)
}

func createTop(date time.Time, c calenderStruct) *fyne.Container {
	vBox := container2.New(layout.NewVBoxLayout())

	leftContainer := container2.New(layout.NewHBoxLayout())
	rightContainer := container2.New(layout.NewHBoxLayout())
	strYearMonth := fmt.Sprintf("%04d-%02d", date.Year(), int(date.Month()))

	weekContainer := container2.New(layout.NewGridLayoutWithColumns(7))

	redSunday := createColorText(weeks[0], color.RGBA{R: 255, A: 255})
	weekContainer.Add(redSunday)

	for _, day := range weeks[1:] {
		weekContainer.Add(createSimpleText(day))
	}

	current := widget.NewButton(strYearMonth, func() {
		inputWindow(c)
	})

	prev := widget.NewButton("<", func() {
		updateByArrow(false, c)
	})

	next := widget.NewButton(">", func() {
		updateByArrow(true, c)
	})

	leftContainer.Add(prev)
	leftContainer.Add(current)
	leftContainer.Add(next)

	// pencil icon
	pencilButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		launchEditor()
	})
	rightContainer.Add(pencilButton)
	// gear icon
	rightContainer.Add(widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		launchFyneSetting()
	}))

	container := container2.New(layout.NewBorderLayout(nil, nil, leftContainer, rightContainer), leftContainer, rightContainer)

	vBox.Add(container)
	vBox.Add(weekContainer)

	return container2.New(layout.NewMaxLayout(), vBox)
}

func updateByInput(TargetTime time.Time, c calenderStruct) {
	c.top = createTop(TargetTime, c)
	c.middle = createCells(
		myModule.DayOfMonth(TargetTime.Year(), TargetTime.Month()),
		myModule.StartOffsetOfMonth(TargetTime.Year(), TargetTime.Month()),
		TargetTime,
	)

	content := container2.New(layout.NewBorderLayout(c.top, nil, nil, nil), c.top, c.middle)
	c.window.SetContent(content)
}

func updateByArrow(isToFuture bool, c calenderStruct) {
	var tempTime time.Time
	if isToFuture {
		tempTime = TargetTime.AddDate(0, 1, 0)
	} else {
		tempTime = TargetTime.AddDate(0, -1, 0)
	}

	showMinYear := time.Date(0, 1, 1, 0, 0, 0, 0, TargetTimeLocal)
	showMaxYear := time.Date(10000, 1, 1, 0, 0, 0, -1, TargetTimeLocal)
	if tempTime.Before(showMinYear) || tempTime.After(showMaxYear) {
		return
	} else if myModule.IsNowTargetYearMonth(tempTime) {
		TargetTime = time.Now()
	} else {
		TargetTime = tempTime
	}

	c.top = createTop(TargetTime, c)
	c.middle = createCells(
		myModule.DayOfMonth(TargetTime.Year(), TargetTime.Month()),
		myModule.StartOffsetOfMonth(TargetTime.Year(), TargetTime.Month()),
		TargetTime,
	)

	content := container2.New(layout.NewBorderLayout(c.top, nil, nil, nil), c.top, c.middle)
	c.window.SetContent(content)
}

func NewCalender() *calenderStruct {
	c := &calenderStruct{}
	c.top = createTop(TargetTime, *c)
	c.middle = createCells(
		myModule.DayOfMonth(TargetTime.Year(), TargetTime.Month()),
		myModule.StartOffsetOfMonth(TargetTime.Year(), TargetTime.Month()),
		TargetTime,
	)
	return c
}

func (c calenderStruct) LoadUI(app fyne.App) {
	c.window = app.NewWindow("MD calender")
	c.top = createTop(TargetTime, c)
	c.middle = createCells(
		myModule.DayOfMonth(TargetTime.Year(), TargetTime.Month()),
		myModule.StartOffsetOfMonth(TargetTime.Year(), TargetTime.Month()),
		TargetTime,
	)
	content := container2.New(layout.NewBorderLayout(c.top, nil, nil, nil), c.top, c.middle)

	c.window.SetContent(content)
	c.window.Resize(fyne.Size{Width: 500, Height: 500})
	c.window.Show()
}

func ShowErr(msg string, err error) {
	w := fyne.CurrentApp().NewWindow(iniSectionKey + " error")
	vBox := container2.New(layout.NewVBoxLayout(), widget.NewLabel(msg), widget.NewLabel(err.Error()))

	w.SetContent(vBox)
	w.SetFixedSize(true)
	w.Show()
}
