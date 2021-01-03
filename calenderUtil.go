package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"gopkg.in/ini.v1"
	"log"
	"myModule"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

var iniFileName = "setting.ini"
var iniConfig, _ = ini.Load(iniFileName)

func getNotifyCountFromMd(day int, mds []myModule.CalenderMdInt) int {
	for _, item := range mds {
		if item.Day == day {
			return len(item.Events)
		}
	}
	return 0
}

func convertMdToInt(path string) []myModule.CalenderMdInt {
	strMDs := myModule.SplitFromMarkdown(path)

	var result []myModule.CalenderMdInt
	for _, item := range strMDs {
		cnv, err := myModule.CalenderConverter(item)
		if err != nil {
			continue
		} else {
			result = append(result, cnv)
		}
	}
	return result
}

func fileOpened(f fyne.URIReadCloser) {
	if f == nil {
		log.Println("Cancelled")
		return
	}

	ext := f.URI().Extension()
	if ext == ".exe" {
		uri, _ := url.ParseRequestURI(f.URI().String())
		//fmt.Printf("url:%v\nscheme:%v\n host:%v\n Path:%v\n", uri, uri.Scheme, uri.Host, uri.Path)
		iniConfig.Section(iniSectionKey).Key(iniEditorKey).SetValue(uri.Host + uri.Path)
		iniConfig.SaveTo(iniFileName)
	}
	err := f.Close()
	if err != nil {
		fyne.LogError("Failed to close stream", err)
	}
}

func launchEditor() {
	fileInfo := getFilePathOrMake()
	cmd := exec.Command(EditorPath, fileInfo)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func getFilePathOrMake() string {
	absFilePath := getFileAbsFilePath()
	_, existErr := os.Stat(absFilePath)
	if existErr != nil {
		dir := filepath.Dir(filepath.Clean(absFilePath))
		//file := filepath.Base(absFilePath)

		_ = os.Mkdir(dir, 0777)
		//_, _ = os.OpenFile(file, os.O_RDONLY|os.O_CREATE, 0777)
	}
	return absFilePath
}

func launchEditorSetting(c calenderStruct) {
	//todo: impl
	win := c.window

	////cfg, err := ini.Load("R:\\project\\goCalender\\setting.ini")
	//cfg, err := ini.Load("setting.ini")
	//if err != nil {
	//	fmt.Printf("Fail to read file: %v", err)
	//	os.Exit(1)
	//}

	//// Now, make some changes and save it
	//cfg.Section("").Key("app_mode").SetValue("production")
	//cfg.SaveTo("my.ini.local")

	//launchEditor(editorPath)

	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err == nil && reader == nil {
			return
		}
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		fileOpened(reader)
	}, win)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".exe"}))

	fd.Show()
}

func launchFyneSetting() {
	w := myApp.NewWindow("Fyne Settings")
	w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
	w.Resize(fyne.NewSize(480, 480))
	w.Show()
}
