package stream

import (
	"github.com/satori/go.uuid"
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

var PlayerList []*Player

// Run Running the cloud gaming server
func (server *StreamServer) Run() {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.Port))
	if err != nil {
		log.Fatal("Error listening", err.Error())
	}
	log.Println("Listen port:", server.Port)

	// Set the first player to None

	NonePlayer := new(Player)
	NonePlayer.ID = "None"
	PlayerList = append(PlayerList, NonePlayer)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting", err.Error())
			return
		}

		// Generate unique ID for each player
		PlayerID := uuid.NewV4()
		player := &Player{
			Conn:     conn,
			ID:       PlayerID.String(),
			GameList: &server.GameList,
		}

		PlayerList = append(PlayerList, player)

		if player.InitTelnet() {
			go player.Serve()
		}
	}
}
