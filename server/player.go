package server

import (
	"fmt"
	"github.com/HFO4/gbc-in-cloud/driver"
	"github.com/HFO4/gbc-in-cloud/gb"
	"log"
	"net"
)

type Player struct {
	Conn     net.Conn
	Emulator *gb.Core
	ID       int
}

var i int
var playerList []*Player

func Run(port string) {

	playerList = make([]*Player, 20)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Error listening", err.Error())
	}

	i = 0

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return // 终止程序
		}
		if i == 1 {
			playerList[20].Emulator.DebugControl = 1
		}
		player := &Player{
			Conn: conn,
			ID:   i,
		}
		playerList = append(playerList, player)
		i++
		if player.Init() {
			go player.Serve()
		}
	}
}

func (player *Player) Init() bool {

	//Send telnet options
	_, err := player.Conn.Write([]byte{255, 253, 34})
	_, err = player.Conn.Write([]byte{255, 250, 34, 1, 0, 255, 240})
	//NOT ECHO
	_, err = player.Conn.Write([]byte{0xFF, 0xFB, 0x01})

	Driver := &driver.ASCII{
		Conn: player.Conn,
	}

	controller := &driver.TelnetController{}
	core := &gb.Core{
		FPS:           60,
		Clock:         4194304,
		Debug:         false,
		DebugControl:  uint16(10 - player.ID),
		DisplayDriver: Driver,
		Controller:    controller,
		DrawSignal:    make(chan bool),
		SpeedMultiple: 0,
		ToggleSound:   true,
		PlayerList:    playerList,
	}

	player.Emulator = core

	if err != nil {
		return false
	}
	log.Println("New Player:", player.ID)
	return true
}

func (player *Player) Serve() {
	go player.Emulator.DisplayDriver.Run(player.Emulator.DrawSignal)
	player.Emulator.Init("G:\\LearnGo\\gb\\Legend of Zelda, The - Link's Awakening (U) (V1.2) [!].gb")
	go player.Emulator.Run()

	for {
		buf := make([]byte, 512)
		n, err := player.Conn.Read(buf)
		if err != nil {
			log.Println("Error reading", err.Error())
			player.Emulator.Exit = true
			log.Println(playerList)
			return
		}
		log.Println(buf[:n])
		player.Emulator.Controller.NewInput(buf[:n])

	}
}
