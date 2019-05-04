package server

import (
	"fmt"
	"github.com/HFO4/gbc-in-cloud/driver"
	"github.com/HFO4/gbc-in-cloud/gb"
	"github.com/logrusorgru/aurora"
	"log"
	"net"
	"strconv"
)

// Player Single player model
type Player struct {
	Conn     net.Conn
	Emulator *gb.Core
	ID       string
	Selected int
	GameList *[]GameInfo

	SelectedPlayer   int
	SelectedPlayerID string
}

// Send TELNET options
func (player *Player) InitTelnet() bool {
	// Send telnet options
	_, err := player.Conn.Write([]byte{255, 253, 34})
	_, err = player.Conn.Write([]byte{255, 250, 34, 1, 0, 255, 240})
	// NOT ECHO
	_, err = player.Conn.Write([]byte{0xFF, 0xFB, 0x01})
	if err != nil {
		return false
	}
	return true
}

func (player *Player) Init() bool {

	if player.Emulator == nil {
		Driver := &driver.ASCII{
			Conn: player.Conn,
		}

		core := &gb.Core{
			// Terminal gaming dose not require high FPS,
			// 10 FPS is a decent choice in most situation.
			FPS:           10,
			Clock:         4194304,
			Debug:         false,
			DisplayDriver: Driver,
			Controller:    new(driver.TelnetController),
			DrawSignal:    make(chan bool),
			SpeedMultiple: 0,
			ToggleSound:   false,
		}

		player.Emulator = core

		log.Println("New Player:", player.ID)

	}
	return true

}

// Generate welcome and game selection screen
func (player *Player) RenderWelcomeScreen() []byte {
	res := "\033[H"
	res += "Welcome to " + fmt.Stringer(aurora.Bold(aurora.Green("Gameboy.Live"))).String() + ", you can enjoy GAMEBOY games in your terminal with \"cloud gaming\" experience.\r\n"
	res += "Use " + fmt.Stringer(aurora.Gray(1-1, "Direction keys").BgGray(24-1)).String() + " in your keyboard to select a game, " + fmt.Stringer(aurora.Gray(1-1, " Enter ").BgGray(24-1)).String() + " key to confirm, " + fmt.Stringer(aurora.Gray(1-1, " M ").BgGray(24-1)).String() + " key to enter multi-player mode and select a partner.\r\n"
	res += "\r\n\r\n"

	for k, v := range *player.GameList {
		if player.Selected == k {
			res += "    " + fmt.Stringer(aurora.Gray(1-1, strconv.Itoa(k+1)+".  "+v.Title+"\r\n").BgGray(24-1)).String()
		} else {
			res += "    " + strconv.Itoa(k+1) + ".  " + v.Title + "\r\n"
		}

	}

	res += "\r\n\r\n" + fmt.Stringer(aurora.Yellow("This service is only playable in terminals with ANSI standard and UTF-8 charset support.")).String() + "\r\n"
	res += "Source code of this project is available at: " + fmt.Stringer(aurora.Underline("https://github.com/HFO4/gameboy.live")).String() + " \r\n"
	return []byte(res)
}

/*
	Show the welcome and game select screen, return
	selected game ID.
*/
func (player *Player) Welcome() int {

	//Clean screen
	_, err := player.Conn.Write([]byte("\033[2J\033[H"))
	if err != nil {
		return -1
	}

	player.Init()

	for {
		var n int
		_, err = player.Conn.Write(player.RenderWelcomeScreen())
		buf := make([]byte, 512)
		n, err = player.Conn.Read(buf)
		inputKey := buf[:n]
		if err != nil {
			return -1
		}

		switch inputKey[len(inputKey)-1] {
		// Up key pressed
		case 65:
			if player.Selected == 0 {
				player.Selected = len(*player.GameList) - 1
			} else {
				player.Selected--
			}
		// Down key pressed
		case 66:
			if player.Selected == len(*player.GameList)-1 {
				player.Selected = 0
			} else {
				player.Selected++
			}
		// Enter key pressed
		case 10, 0:
			return player.Selected
		case 109:
			player.SelectPlayer()
			_, err = player.Conn.Write([]byte("\033[2J\033[H"))

			// If choose each other, connect their serial driver
			if player.SelectedPlayerID != "" && PlayerList[player.SelectedPlayer].SelectedPlayerID == player.ID {
				PlayerList[player.SelectedPlayer].Emulator.Serial.SetTarget(&player.Emulator.Serial)
				player.Emulator.Serial.SetTarget(&PlayerList[player.SelectedPlayer].Emulator.Serial)
				log.Printf("[Serial] Player %s connect with Player %s", player.SelectedPlayerID, PlayerList[player.SelectedPlayer].SelectedPlayerID)
			}
		}

	}

}

/*
	Render select multiplayer screen
*/
func (player *Player) RenderSelectPlayer() []byte {
	res := "\033[2J\033[H"
	res += "You can play multiplayer game with your friend or strangers. The list below lists players who are currently online. Both of you need to choose each other, so that the connection can be established.\r\n"
	res += "Your player ID: " + fmt.Stringer(aurora.Gray(1-1, player.ID).BgGray(24-1)).String() + "\r\n"
	res += "Player list (Press R to refresh):\r\n\r\n"

	for k, v := range PlayerList {

		if player.SelectedPlayer == k {
			res += "    " + fmt.Stringer(aurora.Gray(1-1, v.ID+"\r\n").BgGray(24-1)).String()
		} else {
			res += "    " + v.ID + "\r\n"
		}

	}
	return []byte(res)
}

/*
	Select multiplayer
*/
func (player *Player) SelectPlayer() int {

	for {
		var n int
		_, err := player.Conn.Write(player.RenderSelectPlayer())
		if err != nil {
			return -1
		}
		buf := make([]byte, 512)
		n, err = player.Conn.Read(buf)
		inputKey := buf[:n]
		if err != nil {
			return -1
		}

		switch inputKey[len(inputKey)-1] {
		// Up key pressed
		case 65:
			if player.SelectedPlayer == 0 {
				player.SelectedPlayer = len(PlayerList) - 1
			} else {
				player.SelectedPlayer--
			}
		// Down key pressed
		case 66:
			if player.SelectedPlayer == len(PlayerList)-1 {
				player.SelectedPlayer = 0
			} else {
				player.SelectedPlayer++
			}
		// Enter key pressed
		case 10, 0:
			// Cannot choose yourself
			if PlayerList[player.SelectedPlayer].ID == player.ID {
				continue
			}

			// Choose none
			if player.SelectedPlayer == 0 {
				player.SelectedPlayerID = ""
				return 0
			}

			player.SelectedPlayerID = PlayerList[player.SelectedPlayer].ID
			return 0
		// R key pressed
		case 114:
			continue
		}

		log.Println(inputKey)
	}
	return 0
}

/*
	Generate the control instruction screen,
	ascii art by Joan Stark.
*/

func (player *Player) Instruction() int {
	ret := "Here's the key instruction, press " + fmt.Stringer(aurora.Gray(1-1, "Enter").BgGray(24-1)).String() + " key to enter the game, " + fmt.Stringer(aurora.Gray(1-1, " Q ").BgGray(24-1)).String() + " to quit the game.\r\n\r\n"
	ret += "                      __________________________\r\n" + "                     |OFFo oON                  |\r\n" + "                     | .----------------------. |\r\n" + "                     | |  .----------------.  | |\r\n" + "                     | |  |                |  | |\r\n" + "                     | |))|                |  | |\r\n" + "                     | |  |                |  | |\r\n" + "                     | |  |                |  | |\r\n" + "                     | |  |                |  | |\r\n" + "                     | |  |                |  | |\r\n" + "                     | |  |                |  | |\r\n" + "                     | |  '----------------'  | |\r\n" + "                     | |__GAME BOY____________/ |\r\n" + "    Keyboard:Up↑ <--------+     ________        |\r\n" + "                     |    +    (Nintendo)       |\r\n" + "                     |  _| |_   \"\"\"\"\"\"\"\"   .-.  |\r\n" + "  Keyboard:Left← <----+[_   _]---+    .-. ( +---------> Keyboard:Z\r\n" + "                     |   |_|     |   (   ) '-'  |\r\n" + "                     |    +      |    '-+   A   |\r\n" + "  Keyboard:Down↓ <--------+ +----+     B+-------------> Keyboard:X\r\n" + "                     |      |   ___   ___       |\r\n" + "                     |      |  (___) (___)  ,., |\r\n" + "Keyboard:Right→ <-----------+ select st+rt ;:;: |\r\n" + "                     |           +     |  ,;:;' /\r\n" + "                  jgs|           |     | ,:;:'.'\r\n" + "                     '-----------------------`\r\n" + "                                 |     |\r\n" + "           Keyboard:Backspace <--+     +-> Keyboard:Enter\r\n"
	// Clean screen
	_, err := player.Conn.Write([]byte("\033[2J\033[H" + ret))
	if err != nil {
		return -1
	}
	for {
		buf := make([]byte, 512)
		n, err := player.Conn.Read(buf)
		inputKey := buf[:n]
		if err != nil {
			return -1
		}

		// Enter key pressed
		if inputKey[len(inputKey)-1] == 0 || inputKey[len(inputKey)-1] == 10 {
			return 1
		}
	}
}

func (player *Player) Logout() {
	// Disconnect serial port
	if player.Emulator.Serial.Target != nil {
		player.Emulator.Serial.Target.Target = nil
	}

	playerIndex := 0
	for k, v := range PlayerList {
		if v.ID == player.ID {
			playerIndex = k
			break
		}
	}

	if playerIndex != 0 {
		PlayerList = append(PlayerList[:playerIndex], PlayerList[playerIndex+1:]...)
	}
}

func (player *Player) Serve() {

	game := player.Welcome()

	if game < 0 {
		log.Println("User quit")
		player.Logout()
		return
	}

	if player.Instruction() < 0 {
		log.Println("User quit")
		player.Logout()
		return
	}

	// Set the display driver to TELNET
	go player.Emulator.DisplayDriver.Run(player.Emulator.DrawSignal)
	player.Emulator.Init((*player.GameList)[player.Selected].Path)
	go player.Emulator.Run()

	for {
		buf := make([]byte, 512)
		n, err := player.Conn.Read(buf)
		if err != nil {
			log.Println("Error reading", err.Error())
			player.Emulator.Exit = true
			player.Logout()
			return
		}
		// If "Q" was pressed ,close the connection
		if buf[n-1] == 113 {
			log.Println("User quit")
			player.Emulator.Exit = true
			err := player.Conn.Close()
			if err != nil {
				log.Println("Failed to close connection")
			}
			player.Logout()
			return
		}
		// Handle user input
		player.Emulator.Controller.NewInput(buf[:n])
	}
}
