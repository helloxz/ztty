package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 设置路由
	http.HandleFunc("/ws", handleWebSocket)

	// 启动 HTTP 服务器
	log.Println("Starting server on :65524...")
	go http.ListenAndServe("127.0.0.1:65524", nil)
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Ztty",
		Width:  900,
		Height: 560,

		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		// Windows: &windows.Options{
		// 	WebviewIsTransparent: true,
		// 	WindowIsTranslucent:  true,
		// },
		// Mac: &mac.Options{
		// 	WindowIsTranslucent:  true,
		// 	WebviewIsTransparent: true,
		// },
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
		},

		// CSSDragValue:    "drag",
		// BackgroundColour: &options.RGBA{R: 15, G: 15, B: 20, A: 1},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
