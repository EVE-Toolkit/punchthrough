package main

import (
	"embed"
	"encoding/json"
	"errors"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	appMenu := menu.NewMenu()
	FileMenu := appMenu.AddSubmenu("File")
	FileMenu.AddText("&Save", keys.CmdOrCtrl("s"), func(cd *menu.CallbackData) {
		fileName, err := runtime.SaveFileDialog(app.ctx, runtime.SaveDialogOptions{
			Filters: []runtime.FileFilter{
				{
					DisplayName: "Punchthrough JSON",
					Pattern:     "*.json",
				},
			},
		})

		if err != nil {
			runtime.LogError(app.ctx, err.Error())

			return
		}

		if fileName == "" {
			runtime.LogError(app.ctx, errors.New("It appears the dialog was cancelled. Please try again.").Error())

			return
		}

		stringified, err := json.Marshal(app.Chain)

		if err != nil {
			runtime.LogError(app.ctx, err.Error())

			return
		}

		err = os.WriteFile(fileName, stringified, os.ModePerm)

		if err != nil {
			runtime.LogError(app.ctx, err.Error())

			return
		}
	})

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "punchthrough",
		Width:  500,
		Height: 250,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Menu:             appMenu,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		AlwaysOnTop: true, // remember to turn on
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
