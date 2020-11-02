package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"github.com/HFO4/gbc-in-cloud/driver"
	"github.com/HFO4/gbc-in-cloud/fyne"
	"github.com/HFO4/gbc-in-cloud/gb"
	"github.com/HFO4/gbc-in-cloud/static"
	"github.com/HFO4/gbc-in-cloud/stream"
	"log"
	"os"
)

var (
	h bool

	GUIMode          bool
	FyneMode         bool
	StreamServerMode bool
	StaticServerMode bool

	ConfigPath string
	ListenPort int
	ROMPath    string
	SoundOn    bool
	FPS        int
	Debug      bool
)

func init() {
	flag.BoolVar(&h, "h", false, "This help")
	flag.BoolVar(&GUIMode, "g", true, "Play specific game in GUI mode")
	flag.BoolVar(&FyneMode, "G", false, "Play specific game in Fyne GUI mode")
	flag.BoolVar(&StreamServerMode, "s", false, "Start a cloud-gaming server")
	flag.BoolVar(&StaticServerMode, "S", false, "Start a static image cloud-gaming server")
	flag.BoolVar(&SoundOn, "m", true, "Turn on sound in GUI mode")
	flag.BoolVar(&Debug, "d", false, "Use Debugger in GUI mode")
	flag.IntVar(&ListenPort, "p", 1989, "Set the `port` for the cloud-gaming server")
	flag.IntVar(&FPS, "f", 60, "Set the `FPS` in GUI mode")
	flag.StringVar(&ConfigPath, "c", "", "Set the game option list `config` file path")
	flag.StringVar(&ROMPath, "r", "", "Set `ROM` file path to be played in GUI mode")
}

func startGUI(screen driver.DisplayDriver, control driver.ControllerDriver) {
	core := new(gb.Core)
	core.FPS = FPS
	core.Clock = 4194304
	core.Debug = Debug
	core.DisplayDriver = screen
	core.Controller = control
	core.DrawSignal = make(chan bool)
	core.SpeedMultiple = 0
	core.ToggleSound = SoundOn
	core.Init(ROMPath)

	go core.Run()
	screen.Run(core.DrawSignal, func() {
		core.SaveRAM()
	})
}

func runStaticServer() {
	server := static.StaticServer{
		Port:     ListenPort,
		GamePath: ROMPath,
	}
	server.Run()
}

func runServer() {
	if ConfigPath == "" {
		log.Fatal("[Error] Game list not specified")
	}

	// Read config file
	configFile, err := os.Open(ConfigPath)
	defer configFile.Close()
	if err != nil {
		log.Fatal("[Error] Failed to read game list config file,", err)
	}
	stats, statsErr := configFile.Stat()
	if statsErr != nil {
		log.Fatal(statsErr)
	}
	var size = stats.Size()
	gameListStr := make([]byte, size)
	bufReader := bufio.NewReader(configFile)
	_, err = bufReader.Read(gameListStr)

	streamServer := new(stream.StreamServer)
	streamServer.Port = ListenPort
	var gameList []stream.GameInfo
	err = json.Unmarshal(gameListStr, &gameList)
	if err != nil {
		log.Fatal("Unable to decode game list config file.")
	}
	streamServer.GameList = gameList
	streamServer.Run()
}

func main() {
	flag.Parse()
	if h {
		flag.Usage()
		return
	}

	if StreamServerMode {
		runServer()
		return
	}

	if StaticServerMode {
		runStaticServer()
		return
	}

	if FyneMode {
		driver := new(fyne.LCD)
		startGUI(driver, driver)
		return
	} else if GUIMode {
		driver := new(driver.LCD)
		startGUI(driver, driver)
		return
	}
}
