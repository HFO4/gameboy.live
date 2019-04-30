# Gameboy.Live
🕹️ `Gameboy.Live` is a Gameboy emulator written in go for learning purpose. You can simply play Gameboy games on your desktop:

![https://github.com/HFO4/gameboy.live/raw/master/doc/screenshot.png](https://github.com/HFO4/gameboy.live/raw/master/doc/screenshot.png)

Or, play "Cloud gaming" in your terminal with one single command:

```
telnet gameboy.live 1989
```

![https://github.com/HFO4/gameboy.live/raw/master/doc/cloud-gaming.gif](https://github.com/HFO4/gameboy.live/raw/master/doc/cloud-gaming.gif)

## Installation

You can directly download executable file from [Release](https://github.com/HFO4/gameboy.live/releases) page, or build from the source:

```
git clone https://github.com/HFO4/gameboy.live.git
cd gameboy.live
go build -o gbdotlive main.go
```

## Usage

```
Usage of gbdotlive:
  -c config
        Set the game option list config file path
  -d    Use Debugger in GUI mode
  -f FPS
        Set the FPS in GUI mode (default 60)
  -g    Play specific game in GUI mode (default true)
  -h    This help
  -m    Turn on sound in GUI mode (default true)
  -p port
        Set the port for the cloud-gaming server (default 1989)
  -r ROM
        Set ROM file path to be played in GUI mode
  -s    Start a cloud-gaming server
```

### GUI mode

Play specified ROM file in GUI mode:

```
gbdotlive -r "Tetris.gb" 
```

### Set up a cloud gaming server

You can use `Gameboy.Live` as a "cloud gaming" server, where players use telnet to play Gameboy games in terminal without additional software installation required. (Except telnet itself XD)

First of all, a `gamelist.json` config file is required to specify game options, this is a typical example:

```json
[{
	"Title": "Tetris",
	"Path": "test.gb"
}, {
	"Title": "Dr. Mario",
	"Path": "Dr. Mario (JU) (V1.1).gb"
}, {
	"Title": "Legend of Zelda - Link's Awakening",
	"Path": "Legend of Zelda, The - Link's Awakening (U) (V1.2) [!].gb"
}]

```

It is recommended to test every ROMs before put them in the config file.

Next, start `Gameboy.Live` server with config file from the previous step:

```
gbdotlive -s -c "gamelist.json"
```

You will see an output like this, which means your server has started successfuly:

```
2019/04/30 21:27:56 Listen port: 1989 
```

Now, you can play games anywhere you want. The simulation and rendering process is done entirely on the server.

```
telnet <ip of your server>:<port>
```

To be careful of, this "cloud gaming" is only playable in terminals support [ANSI](https://en.wikipedia.org/wiki/ANSI_escape_code) standard and UTF-8 charset. You can use `WSL` instead of `CMD` on Windows.

## Keyboard instruction

| Keyboard | Gameboy |
| -------- | ------- |
| <kbd>Enter</kbd>     | Start   |
|<kbd>Backspace</kbd>  | Select  |
| <kbd>↑</kbd>  | Up      |
|  <kbd>↓</kbd> | Down    |
|   <kbd>←</kbd> | Left    |
|   <kbd>→</kbd>  | Right   |
|    <kbd>X</kbd>  | B       |
|     <kbd>Z</kbd>     | A       |

## Features & TODOs

- [x] CPU instruction emulation
- [x] Timer and interrupt
- [x] Support for ROM-only、MBC1、MBC2、MBC3 cartridge
- [x] Sound emulation
- [x] Graphics emulation
- [x] Cloud gaming
- [x] ROM debugger

There are still many places to be perfected：

- [ ] Support Gameboy Color emulation
- [ ] Support for MBC4、MBC5、HuC1 cartridge
- [ ] Sound simulation is incomplete, still got differences compared to the Gameboy real machine
- [ ] Sprite priority issue (see `Wario Land II` and `Metroid II: Return of Samus`)
- [ ] Failed to pass Blargg's instruction timing test
- [ ] Game saving & restore in cartridge level
- [ ] Game saving & restore in emulator level

## Testing

![Testing result](https://github.com/HFO4/gameboy.live/raw/master/doc/Testing.jpg)

## Contribution

This emulator is just for learning and entertainment purpose. There are still many places to be perfected. Any suggestions or contribution is welcomed!

## Reference

* [Pan Docs](http://bgb.bircd.org/pandocs.htm)
* [http://www.codeslinger.co.uk/pages/projects/gameboy/beginning.html](http://www.codeslinger.co.uk/pages/projects/gameboy/beginning.html)
* [http://www.devrs.com/gb/files/GBCPU_Instr.html](http://www.devrs.com/gb/files/GBCPU_Instr.html)
* [https://github.com/Humpheh/goboy](https://github.com/Humpheh/goboy)
* [The Ultimate Game Boy Talk (33c3)](https://www.youtube.com/watch?v=HyzD8pNlpwI)
* [http://gameboy.mongenel.com/dmg/asmmemmap.html](http://gameboy.mongenel.com/dmg/asmmemmap.html)
* ......
