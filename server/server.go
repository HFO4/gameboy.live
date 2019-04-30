package server

import (
	"log"
	"net"
	"strconv"
)

type StreamServer struct {
	Port     int
	GameList []GameInfo
}

type GameInfo struct {
	Title string
	Path  string
}

// Run Running the cloud gaming server
func (server *StreamServer) Run() {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.Port))
	if err != nil {
		log.Fatal("Error listening", err.Error())
	}
	log.Println("Listen port:", server.Port)

	i := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting", err.Error())
			return
		}
		player := &Player{
			Conn:     conn,
			ID:       i,
			GameList: &server.GameList,
		}
		i++
		if player.InitTelnet() {
			go player.Serve()
		}
	}
}
