package static

import (
	"fmt"
	"github.com/HFO4/gbc-in-cloud/driver"
	"github.com/HFO4/gbc-in-cloud/gb"
	"net/http"
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
	http.ListenAndServe(fmt.Sprintf(":%d", server.Port), nil)
}

func showImage(server *StaticServer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		server.driver.Render()
	}
}
