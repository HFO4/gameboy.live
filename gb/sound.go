package gb

import (
	"github.com/HFO4/gbc-in-cloud/util"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"log"
	"math"
	"time"
)

type Sound struct {
	Channel2 Channel
	Channel1 Channel

	enable bool

	leftVolume  uint8
	rightVolume uint8

	VRAMCache []byte
}

type Channel struct {
	self   *Channel
	parent *Sound

	enable bool

	envelopeIncrease bool
	envelopeInitial  byte
	envelopeSweepNum byte
	lastEnvelope     float64
	envelopeTick     float64

	stopWhileTimeout bool

	Freq     int
	freqLow  uint16
	freqHigh uint16
	waveDuty byte

	volume float64

	duration   float64
	sampleTick float64
	tickUnit   float64

	// Wave type
	//		0 - Square
	//		1 - Noise
	//		2 - Sample
	wave int
}

/*
	Wave Duty:
	  00: 12.5% ( _-------_-------_------- )
	  01: 25%   ( __------__------__------ )
	  10: 50%   ( ____----____----____---- ) (normal)
	  11: 75%   ( ______--______--______-- )
*/
var waveDutyMap = [5]float64{
	1: -0.25,
	2: -0.5,
	3: 0,
	4: 0.5,
}

const secondPerTick = 1 / 44100.0

func (sound *Sound) Init() {
	log.Println("[Sound] Initialize Sound process unit")
	sound.enable = true
	sound.Channel2.enable = false
	sound.Channel2.self = &sound.Channel2
	sound.Channel2.parent = sound

	sound.Channel1.enable = false
	sound.Channel1.self = &sound.Channel1
	sound.Channel1.parent = sound
	go sound.Play()
}

func (sound *Sound) Play() {
	sr := beep.SampleRate(44100)
	err := speaker.Init(sr, sr.N(time.Second/30))
	if err != nil {
		log.Println("[Warning] Failed to init sound speaker")
	}

	stream := beep.Mixer{}
	stream.Add(sound.Channel1)
	stream.Add(sound.Channel2)

	done := make(chan bool)
	rawStream := beep.Seq(&stream, beep.Callback(func() {
		done <- true
	}))
	volume := &effects.Volume{
		Streamer: rawStream,
		Base:     2,
		Volume:   -3,
	}
	speaker.Play(volume)
	<-done
}

/*
	When sound related memory is writen, this function will be
	called to update sound props.
*/
func (sound *Sound) Trigger(address uint16, val byte, vram []byte) {
	sound.VRAMCache = vram
	//log.Printf("new Sound:%X : %X tick:%d\n", address, val,sound.Channel2.sampleTick)
	switch address {
	case 0xFF26:
		/*
			FF26 - NR52 - Sound on/off
			  Bit 7 - All sound on/off  (0: stop all sound circuits) (Read/Write)
			  Bit 3 - Sound 4 ON flag (Read Only)
			  Bit 2 - Sound 3 ON flag (Read Only)
			  Bit 1 - Sound 2 ON flag (Read Only)
			  Bit 0 - Sound 1 ON flag (Read Only)
		*/
		sound.enable = util.TestBit(val, 7)
	case 0xFF25:
		/*
			FF25 - NR51 - Selection of Sound output terminal (R/W)
			  Bit 7 - Output sound 4 to SO2 terminal
			  Bit 6 - Output sound 3 to SO2 terminal
			  Bit 5 - Output sound 2 to SO2 terminal
			  Bit 4 - Output sound 1 to SO2 terminal
			  Bit 3 - Output sound 4 to SO1 terminal
			  Bit 2 - Output sound 3 to SO1 terminal
			  Bit 1 - Output sound 2 to SO1 terminal
			  Bit 0 - Output sound 1 to SO1 terminal

			TODO: Separate Right/Left channel
		*/
		if !(util.TestBit(val, 5) || util.TestBit(val, 1)) {
			sound.Channel2.enable = false
		}
	case 0xFF24:
		/*
			FF24 - NR50 - Channel control / ON-OFF / Volume (R/W)
				The volume bits specify the "Master Volume" for Left/Right sound output.
				  Bit 7   - Output Vin to SO2 terminal (1=Enable)
				  Bit 6-4 - SO2 output level (volume)  (0-7)
				  Bit 3   - Output Vin to SO1 terminal (1=Enable)
				  Bit 2-0 - SO1 output level (volume)  (0-7)
		*/
		sound.leftVolume = uint8(val & 0x7)
		sound.rightVolume = uint8((val >> 4) & 0x7)
	// Channel 1
	case 0xFF11:
		/*
			FF16 - NR21 - Channel 2 Sound Length/Wave Pattern Duty (R/W)
			  Bit 7-6 - Wave Pattern Duty (Read/Write)
			  Bit 5-0 - Sound length data (Write Only) (t1: 0-63)
		*/
		sound.Channel1.waveDuty = (val >> 6) + 1
		sound.Channel1.duration = 64.0 - float64(val&0x3F)*(1.0/256.0)
	case 0xFF14:
		/*
			FF19 - NR24 - Channel 2 Frequency hi data (R/W)
			  Bit 7   - Initial (1=Restart Sound)     (Write Only)
			  Bit 6   - Counter/consecutive selection (Read/Write)
						(1=Stop output when length in NR21 expires)
			  Bit 2-0 - Frequency's higher 3 bits (x) (Write Only)
			Frequency = 131072/(2048-x) Hz
		*/
		if util.TestBit(val, 7) {
			sound.Channel1.reSet()
			sound.Channel1.envelopeIncrease = util.TestBit(sound.VRAMCache[0x02], 3)
			sound.Channel1.envelopeInitial = sound.VRAMCache[0x02] >> 4
			sound.Channel1.envelopeSweepNum = sound.VRAMCache[0x02] & 0x7
			sound.Channel1.volume = float64(sound.Channel1.envelopeInitial) / float64(0xf)
			sound.Channel1.enable = true
		}
		sound.Channel1.stopWhileTimeout = util.TestBit(val, 6)
		sound.Channel1.freqHigh = uint16((val & 0x7)) << 8
		sound.Channel1.Freq = 131072 / (2048 - int(sound.Channel1.freqHigh+sound.Channel1.freqLow))
		sound.Channel1.tickUnit = 44100.0 / float64(sound.Channel1.Freq)
		//log.Println(sound.Channel1)
	case 0xFF13:
		/*
			FF18 - NR23 - Channel 2 Frequency lo data (W)
				Frequency's lower 8 bits of 11 bit data (x).
				Next 3 bits are in NR24 ($FF19).
		*/
		sound.Channel1.freqLow = uint16(val)
		sound.Channel1.Freq = 131072 / (2048 - int(sound.Channel1.freqHigh+sound.Channel1.freqLow))
		sound.Channel1.tickUnit = 44100.0 / float64(sound.Channel1.Freq)
		//log.Println(sound.Channel1)

	// Channel 2
	case 0xFF16:
		/*
			FF16 - NR21 - Channel 2 Sound Length/Wave Pattern Duty (R/W)
			  Bit 7-6 - Wave Pattern Duty (Read/Write)
			  Bit 5-0 - Sound length data (Write Only) (t1: 0-63)
		*/
		sound.Channel2.waveDuty = (val >> 6) + 1
		sound.Channel2.duration = 64.0 - float64(val&0x3F)*(1.0/256.0)
	case 0xFF19:
		/*
			FF19 - NR24 - Channel 2 Frequency hi data (R/W)
			  Bit 7   - Initial (1=Restart Sound)     (Write Only)
			  Bit 6   - Counter/consecutive selection (Read/Write)
						(1=Stop output when length in NR21 expires)
			  Bit 2-0 - Frequency's higher 3 bits (x) (Write Only)
			Frequency = 131072/(2048-x) Hz
		*/
		if util.TestBit(val, 7) {
			sound.Channel2.reSet()
			sound.Channel2.envelopeIncrease = util.TestBit(sound.VRAMCache[0x07], 3)
			sound.Channel2.envelopeInitial = sound.VRAMCache[0x07] >> 4
			sound.Channel2.envelopeSweepNum = sound.VRAMCache[0x07] & 0x7
			sound.Channel2.volume = float64(sound.Channel2.envelopeInitial) / float64(0xf)
			sound.Channel2.enable = true
		}
		sound.Channel2.stopWhileTimeout = util.TestBit(val, 6)
		sound.Channel2.freqHigh = uint16((val & 0x7)) << 8
		sound.Channel2.Freq = 131072 / (2048 - int(sound.Channel2.freqHigh+sound.Channel2.freqLow))
		sound.Channel2.tickUnit = 44100.0 / float64(sound.Channel2.Freq)
		//log.Println(sound.Channel2)
	case 0xFF18:
		/*
			FF18 - NR23 - Channel 2 Frequency lo data (W)
				Frequency's lower 8 bits of 11 bit data (x).
				Next 3 bits are in NR24 ($FF19).
		*/
		sound.Channel2.freqLow = uint16(val)
		sound.Channel2.Freq = 131072 / (2048 - int(sound.Channel2.freqHigh+sound.Channel2.freqLow))
		sound.Channel2.tickUnit = 44100.0 / float64(sound.Channel2.Freq)
	}
}

func (channel *Channel) shouldPlay() bool {

	// Is the master sound control enabled?
	if !channel.parent.enable {
		return false
	}

	if channel.Freq <= 0 {
		return false
	}

	// Is this channel enabled?
	if !channel.enable {
		return false
	}

	if channel.waveDuty == 0 {
		return false
	}

	// Is this channel reach time limit?
	if channel.stopWhileTimeout && channel.duration <= 0 {
		return false
	}
	return true
}

func (channel Channel) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {

		channel.self.sampleTick += float64(channel.self.Freq) / 44100.0

		if channel.self.shouldPlay() {
			tickInCycle := channel.self.sampleTick * 2 * 3.1415926
			switch channel.wave {
			case 0:
				if math.Sin(tickInCycle) <= waveDutyMap[int(channel.self.waveDuty)] {
					samples[i][0] = 1 * channel.self.volume
					samples[i][1] = 1 * channel.self.volume
				} else {
					samples[i][0] = 0
					samples[i][1] = 0
				}

			}
			channel.self.duration -= secondPerTick
		} else {
			samples[i][0] = 0
			samples[i][1] = 0
		}

		channel.self.Envelope()
	}
	return len(samples), true
}

func (channel *Channel) Envelope() {
	channel.envelopeTick += 1 / 44100.0
	if channel.envelopeSweepNum > 0 {
		step := float64(channel.envelopeSweepNum) * (1.0 / 64)
		if channel.envelopeTick-channel.lastEnvelope >= step {
			if channel.envelopeInitial > 0 {
				channel.envelopeInitial--
				channel.volume = float64(channel.envelopeInitial) / float64(0xf)
				channel.lastEnvelope = channel.envelopeTick
			}
		}
	}
}

func (channel *Channel) reSet() {
	channel.enable = false
	channel.Freq = 0
	channel.envelopeTick = 0
	channel.lastEnvelope = 0
}

func (channel Channel) Err() error {
	return nil
}
