package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"gopkg.in/ini.v1"
	"log"
	"myModule"
	"os"
	"time"
)

var TargetTimeLocal = time.Local
var TargetTime = time.Now().In(TargetTimeLocal)

var OutputPath string
var EditorPath string

var iniSectionKey = "MD Calender"
var iniDirectoryKey = "outputPath"
var iniEditorKey = "editorPath"

//const location = "Asia/Tokyo"
//var TargetTime = time.Now().In(time.FixedZone(location, 9*60*60))
var myApp = app.New()

func main() {
	iconPath := "./Icon.png"
	mdIcon, _ := fyne.LoadResourceFromPath(iconPath)
	myApp.SetIcon(mdIcon)

	if checkIniSetting() != nil {
		msg := "load ini setting file failed"
		fyne.LogError(msg, checkIniSetting())
		ShowErr(msg, checkIniSetting())
	} else {
		log.Print(fmt.Sprintf("output path:%v, editor path:%v", OutputPath, EditorPath))
		myModule.TargetTimeLocal = TargetTimeLocal
		paddingSize = theme.TextSize() / 4
		c := NewCalender()
		c.LoadUI(myApp)
	}

	myApp.Run()
}

func checkIniSetting() error {
	iniConf, iniErr := ini.Load(iniFileName)
	if iniErr != nil {
		return iniErr
	}

	OutputPath = iniConf.Section(iniSectionKey).Key(iniDirectoryKey).String()
	EditorPath = iniConf.Section(iniSectionKey).Key(iniEditorKey).String()

	if err := checkExist(OutputPath); err != nil {
		return err
	}
	if err := checkExist(EditorPath); err != nil {
		return err
	}

	return nil
}

func checkExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	return nil
}
