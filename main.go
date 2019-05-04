package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"github.com/HFO4/gbc-in-cloud/driver"
	"github.com/HFO4/gbc-in-cloud/gb"
	"github.com/HFO4/gbc-in-cloud/server"
	"github.com/faiface/pixel/pixelgl"
	"log"
	"os"
)

var (
	h bool

	GUIMode    bool
	ServerMode bool

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
	flag.BoolVar(&ServerMode, "s", false, "Start a cloud-gaming server")
	flag.BoolVar(&SoundOn, "m", true, "Turn on sound in GUI mode")
	flag.BoolVar(&Debug, "d", false, "Use Debugger in GUI mode")
	flag.IntVar(&ListenPort, "p", 1989, "Set the `port` for the cloud-gaming server")
	flag.IntVar(&FPS, "f", 60, "Set the `FPS` in GUI mode")
	flag.StringVar(&ConfigPath, "c", "gamelist.json", "Set the game option list `config` file path")
	flag.StringVar(&ROMPath, "r", "test.gb", "Set `ROM` file path to be played in GUI mode")
}

func startGUI() {
	loop := make(chan bool)

	Driver := new(driver.LCD)
	core := new(gb.Core)
	core.FPS = 60
	core.Clock = 4194304
	core.Debug = Debug
	core.DisplayDriver = Driver
	core.Controller = Driver
	core.DrawSignal = make(chan bool)
	core.SpeedMultiple = 0
	core.ToggleSound = SoundOn
	go core.DisplayDriver.Run(core.DrawSignal)
	core.Init(ROMPath)
	go core.Run()

	//t:=0
	//fmt.Scanf("%d",&t)

	Driver2 := new(driver.LCD)
	core2 := new(gb.Core)
	core2.FPS = 60
	core2.Clock = 4194304
	core2.Debug = Debug
	core2.DisplayDriver = Driver2
	core2.Controller = Driver2
	core2.DrawSignal = make(chan bool)
	core2.SpeedMultiple = 0
	core2.ToggleSound = SoundOn
	go core2.DisplayDriver.Run(core2.DrawSignal)
	core2.Init(ROMPath)
	go core2.Run()

	<-loop
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

	streamServer := new(server.StreamServer)
	streamServer.Port = ListenPort
	var gameList []server.GameInfo
	err = json.Unmarshal(gameListStr, &gameList)
	if err != nil {
		log.Fatal("Unable to decode game list config file.")
	}
	streamServer.GameList = gameList
	streamServer.Run()
}

func run() {
	flag.Parse()
	if h {
		flag.Usage()
		return
	}
	runServer()

	if ServerMode {
		runServer()
		return
	}

	//if GUIMode {
	//	startGUI()
	//	return
	//}
}

func main() {
	pixelgl.Run(run)
}
