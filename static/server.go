package static

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/HFO4/gbc-in-cloud/driver"
	"github.com/HFO4/gbc-in-cloud/gb"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	http.HandleFunc("/svg", showSVG(server))
	http.HandleFunc("/control", newInput(server))
	http.ListenAndServe(fmt.Sprintf(":%d", server.Port), nil)
}

func showSVG(server *StaticServer) func(http.ResponseWriter, *http.Request) {
	svg, _ := ioutil.ReadFile("gb.svg")

	return func(w http.ResponseWriter, req *http.Request) {
		callback, _ := req.URL.Query()["callback"]

		w.Header().Set("Cache-control", "no-cache,max-age=0")
		w.Header().Set("Content-type", "image/svg+xml")
		w.Header().Set("Expires", time.Now().Add(time.Duration(-1)*time.Hour).UTC().Format(http.TimeFormat))

		// Encode image to Base64
		img := server.driver.Render()
		var imageBuf bytes.Buffer
		png.Encode(&imageBuf, img)
		encoded := base64.StdEncoding.EncodeToString(imageBuf.Bytes())

		// Embaded image into svg template
		res := strings.ReplaceAll(string(svg), "{image}", "data:image/png;base64,"+encoded)

		// Replace callback url
		res = strings.ReplaceAll(res, "{callback}", callback[0])

		w.Write([]byte(res))
	}
}

func showImage(server *StaticServer) func(http.ResponseWriter, *http.Request) {
	lastSave := time.Now().Add(time.Duration(-1) * time.Hour)
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Cache-control", "no-cache,max-age=0")
		w.Header().Set("Content-type", "image/png")
		w.Header().Set("Expires", time.Now().Add(time.Duration(-1)*time.Hour).UTC().Format(http.TimeFormat))
		img := server.driver.Render()
		png.Encode(w, img)

		// Save snapshot every 10 minutes
		if time.Now().Sub(lastSave).Minutes() > 10 {
			lastSave = time.Now()
			if snapshot, err := os.Create("snapshots/" + strconv.FormatInt(time.Now().Unix(), 10) + ".png"); err == nil {
				png.Encode(snapshot, img)
				snapshot.Close()
			} else {
				fmt.Println(err)
			}
		}
	}
}

func newInput(server *StaticServer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		keys, ok := req.URL.Query()["button"]
		callback, _ := req.URL.Query()["callback"]

		if !ok || len(keys) < 1 || len(callback) < 1 {
			return
		}

		buttonByte, err := strconv.ParseUint(keys[0], 10, 32)
		if err != nil || buttonByte > 7 {
			return
		}

		if callback[0] == "v2" {
			callback[0] = "https://www.v2ex.com/t/721249"
		}

		server.driver.EnqueueInput(byte(buttonByte))
		time.Sleep(time.Duration(500) * time.Millisecond)
		http.Redirect(w, req, callback[0], http.StatusSeeOther)

		// record input log
		go func(ip, refer, button string) {
			body := map[string]string{
				"refer":  refer,
				"button": button,
				"ip":     ip,
			}
			bodyJson, _ := json.Marshal(body)
			http.Post("http://playground.aoaoao.me/Api/NewGBCommand", "application/json", bytes.NewReader(bodyJson))
		}(ReadUserIP(req), req.Header.Get("referer"), keys[0])
	}
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
