package static

import (
	"fmt"
	"github.com/HFO4/gbc-in-cloud/driver"
	"github.com/HFO4/gbc-in-cloud/gb"
	"image/png"
	"net/http"
	"strconv"
	"time"
)

type StaticServer struct {
	Port     int
	GamePath string

	driver *driver.StaticImage
}

// Run Running the static-image gaming server
func (server *StaticServer) Run() {
	// startup the emulator
	server.driver = &driver.StaticImage{}
	core := &gb.Core{
		FPS:           60,
		Clock:         4194304,
		Debug:         false,
		DisplayDriver: server.driver,
		Controller:    server.driver,
		DrawSignal:    make(chan bool),
		SpeedMultiple: 0,
		ToggleSound:   false,
	}
	go core.DisplayDriver.Run(core.DrawSignal, func() {})
	core.Init(server.GamePath)
	go core.Run()

	// image and control server
	http.HandleFunc("/image", showImage(server))
	http.HandleFunc("/control", newInput(server))
	http.ListenAndServe(fmt.Sprintf(":%d", server.Port), nil)
}

func showImage(server *StaticServer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Cache-control", "no-cache")
		w.Header().Set("Content-type", "image/png")
		img := server.driver.Render()
		png.Encode(w, img)
	}
}

func newInput(server *StaticServer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		keys, ok := req.URL.Query()["button"]
		callback, _ := req.URL.Query()["callback"]

		if !ok || len(keys[0]) < 1 {
			return
		}

		buttonByte, err := strconv.ParseUint(keys[0], 10, 32)
		if err != nil || buttonByte > 7 {
			return
		}

		server.driver.EnqueueInput(byte(buttonByte))
		time.Sleep(time.Duration(500) * time.Millisecond)
		http.Redirect(w, req, callback[0], http.StatusSeeOther)
	}
}
