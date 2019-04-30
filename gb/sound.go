package gb

import (
	"github.com/HFO4/gbc-in-cloud/util"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"log"
	"math"
	"math/rand"
	"time"
)

type Sound struct {
	Channel4 Channel
	Channel3 Channel
	Channel2 Channel
	Channel1 Channel

	enable bool

	leftVolume  uint8
	rightVolume uint8

	VRAMCache   []byte
	SampleCache [32]float64
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

	sweepIncrease bool
	sweepNumber   byte
	sweepTime     int
	lastSweep     float64
	sweepTick     float64
	freqInitial   int
	freqLast      int

	stopWhileTimeout bool

	Freq     int
	freqLow  uint16
	freqHigh uint16
	waveDuty byte

	volume float64

	duration   float64
	sampleTick float64
	tickUnit   float64

	// For noise wave
	lastGenerate     float64
	lastGenerateTick float64

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

var sweepTime = [8]float64{
	0: 0.0,
	1: 1.0 / 128,
	2: 2.0 / 128,
	3: 3.0 / 128,
	4: 4.0 / 128,
	5: 5.0 / 128,
	6: 6.0 / 128,
	7: 7.0 / 128,
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

	sound.Channel3.enable = false
	sound.Channel3.self = &sound.Channel3
	sound.Channel3.parent = sound
	sound.Channel3.wave = 3

	sound.Channel4.enable = false
	sound.Channel4.self = &sound.Channel4
	sound.Channel4.parent = sound
	sound.Channel4.wave = 2

	go sound.Play()
}

func (sound *Sound) Play() {
	sr := beep.SampleRate(44100)
	err := speaker.Init(sr, sr.N(time.Second/30))
	if err != nil {
		log.Println("[Warning] Failed to init sound speaker")
	}

	stream := &beep.Mixer{}
	stream.Add(sound.Channel1)
	stream.Add(sound.Channel2)
	stream.Add(sound.Channel3)
	stream.Add(sound.Channel4)

	//done := make(chan bool)
	//rawStream := beep.Seq(&stream, beep.Callback(func() {
	//	done <- true
	//}))
	volume := &effects.Volume{
		Streamer: stream,
		Base:     2,
		Volume:   -3,
	}
	speaker.Play(volume)
	//<-done
}

/*
	When sound related memory is writen, this function will be
	called to update sound props.
*/
func (sound *Sound) Trigger(address uint16, val byte, vram []byte) {
	sound.VRAMCache = vram
	//log.Printf("new Sound:%X : %X tick:%d\n", address, val,sound.Channel2.sampleTick)
	if address >= 0xFF30 {
		count := 0
		for i := 0; i < 0xF; i++ {
			sound.SampleCache[count] = float64(vram[0x20+i]>>4) / float64(0xf)
			count++
			sound.SampleCache[count] = float64(vram[0x20+i]&0xF) / float64(0xf)
		}
	}
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
	case 0xFF14:
		if util.TestBit(val, 7) {
			sound.Channel1.reSet()
			//Envelope options
			sound.Channel1.envelopeIncrease = util.TestBit(sound.VRAMCache[0x02], 3)
			sound.Channel1.envelopeInitial = (sound.VRAMCache[0x02] & 0xF0) >> 4
			sound.Channel1.envelopeSweepNum = sound.VRAMCache[0x02] & 0x7
			sound.Channel1.volume = float64(sound.Channel1.envelopeInitial) / float64(0xf)

			//Sweep options
			sound.Channel1.sweepIncrease = !util.TestBit(sound.VRAMCache[0x00], 3)
			sound.Channel1.sweepNumber = sound.VRAMCache[0x00] & 0x7
			sound.Channel1.sweepTime = int((sound.VRAMCache[0x00] & 0x70) >> 4)

			sound.Channel1.freqInitial = int(sound.Channel1.freqHigh + sound.Channel1.freqLow)
			sound.Channel1.freqLast = int(sound.Channel1.freqHigh + sound.Channel1.freqLow)

			sound.Channel1.stopWhileTimeout = util.TestBit(val, 6)
			sound.Channel1.freqHigh = uint16((val & 0x7)) << 8
			sound.Channel1.Freq = 131072 / (2048 - int(sound.Channel1.freqHigh+uint16(sound.VRAMCache[0x03])))
			sound.Channel1.tickUnit = 44100.0 / float64(sound.Channel1.Freq)

			sound.Channel1.waveDuty = (sound.VRAMCache[0x01] >> 6) + 1
			sound.Channel1.duration = (64.0 - float64(sound.VRAMCache[0x01]&0x3F)) * (1.0 / 256.0)

			sound.Channel1.enable = true

			//log.Println(sound.Channel1)
		}
		//log.Println(sound.Channel1)
	case 0xFF11:
		sound.Channel1.waveDuty = (sound.VRAMCache[0x01] >> 6) + 1
	case 0xFF13:
		//sound.Channel1.freqLow = uint16(val)
		//sound.Channel1.Freq = 131072 / (2048 - int(sound.Channel1.freqHigh+sound.Channel1.freqLow))
		//sound.Channel1.tickUnit = 44100.0 / float64(sound.Channel1.Freq)
		//log.Println(sound.Channel1)

	// Channel 2
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
			sound.Channel2.envelopeInitial = (sound.VRAMCache[0x07] & 0xF0) >> 4
			sound.Channel2.envelopeSweepNum = sound.VRAMCache[0x07] & 0x7
			sound.Channel2.volume = float64(sound.Channel2.envelopeInitial) / float64(0xf)

			sound.Channel2.stopWhileTimeout = util.TestBit(val, 6)
			sound.Channel2.freqHigh = uint16((val & 0x7)) << 8
			sound.Channel2.Freq = 131072 / (2048 - int(sound.Channel2.freqHigh+sound.Channel2.freqLow))
			sound.Channel2.tickUnit = 44100.0 / float64(sound.Channel2.Freq)

			sound.Channel2.waveDuty = (sound.VRAMCache[0x06] >> 6) + 1
			sound.Channel2.duration = (64.0 - float64(sound.VRAMCache[0x06]&0x3F)) * (1.0 / 256.0)

			sound.Channel2.enable = true
		}
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

	//Channel 3
	case 0xFF1A:
		/*
			FF1A - NR30 - Channel 3 Sound on/off (R/W)
			  Bit 7 - Sound Channel 3 Off  (0=Stop, 1=Playback)  (Read/Write)
		*/
		if util.TestBit(val, 7) {
			sound.Channel3.enable = true
		} else {
			sound.Channel3.enable = false
		}
	case 0xFF1C:
		/*
			FF1C - NR32 - Channel 3 Select output level (R/W)
			  Bit 6-5 - Select output level (Read/Write)
			Possible Output levels are:
			  0: Mute (No sound)
			  1: 100% Volume (Produce Wave Pattern RAM Data as it is)
			  2:  50% Volume (Produce Wave Pattern RAM data shifted once to the right)
			  3:  25% Volume (Produce Wave Pattern RAM data shifted twice to the right)
		*/
		switch (val & 0x60) >> 5 {
		case 0x00:
			sound.Channel3.volume = 0
		case 0x01:
			sound.Channel3.volume = 1
		case 0x02:
			sound.Channel3.volume = 0.5
		case 0x03:
			sound.Channel3.volume = 0.25
		}

	case 0xFF1E:
		if util.TestBit(val, 7) {
			sound.Channel3.reSet()
			sound.Channel3.Freq = 65536 / (2048 - int(uint16((val&0x7))<<8+uint16(sound.VRAMCache[0x0D])))
			sound.Channel3.stopWhileTimeout = util.TestBit(val, 6)
			sound.Channel3.waveDuty = 1
			sound.Channel3.duration = (256.0 - float64(int(sound.VRAMCache[0x0B]))) * (1.0 / 256.0)
			sound.Channel3.enable = true
		}

	// Channel 4
	case 0xFF22:
		/*
			FF22 - NR43 - Channel 4 Polynomial Counter (R/W)
				The amplitude is randomly switched between high and low at the given frequency.
				A higher frequency will make the noise to appear 'softer'.
				When Bit 3 is set, the output will become more regular,
				and some frequencies will sound more like Tone than Noise.
				  Bit 7-4 - Shift Clock Frequency (s)
				  Bit 3   - Counter Step/Width (0=15 bits, 1=7 bits)
				  Bit 2-0 - Dividing Ratio of Frequencies (r)
				Frequency = 524288 Hz / r / 2^(s+1) ;For r=0 assume r=0.5 instead
		*/
		shiftClockFrequency := float64((val & 0xF0) >> 4)
		dividingRatio := float64(val & 0x7)
		if dividingRatio == 0 {
			dividingRatio = 0.5
		}
		sound.Channel4.Freq = int(524288 / dividingRatio / math.Pow(2, shiftClockFrequency+1))
	case 0xFF23:
		if util.TestBit(val, 7) {
			sound.Channel4.enable = false
			sound.Channel4.envelopeTick = 0
			sound.Channel4.lastEnvelope = 0
			sound.Channel4.sampleTick = 0
			sound.Channel4.lastEnvelope = 0

			//Envelope options
			sound.Channel4.envelopeIncrease = util.TestBit(sound.VRAMCache[0x11], 3)
			sound.Channel4.envelopeInitial = (sound.VRAMCache[0x11] & 0xF0) >> 4
			sound.Channel4.envelopeSweepNum = sound.VRAMCache[0x11] & 0x7
			sound.Channel4.volume = float64(sound.Channel4.envelopeInitial) / float64(0xf)
			sound.Channel4.waveDuty = 1

			sound.Channel4.duration = 64.0 - float64(int(sound.VRAMCache[0x10]))*(1.0/256.0)
			sound.Channel4.stopWhileTimeout = util.TestBit(val, 6)

			sound.Channel4.enable = true
		}

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
			case 3:
				sampleID := math.Floor(math.Mod(channel.self.sampleTick, 1.0) * 32)
				samples[i][0] = channel.parent.SampleCache[int(sampleID)] * channel.self.volume
				samples[i][1] = channel.parent.SampleCache[int(sampleID)] * channel.self.volume
			case 2:
				if channel.self.sampleTick-channel.self.lastGenerateTick > 1 || channel.self.sampleTick-channel.self.lastGenerateTick < -1 {
					sample := rand.Float64()*2 - 1
					samples[i][0] = sample * channel.self.volume
					samples[i][1] = sample * channel.self.volume
					channel.self.lastGenerate = sample
					channel.self.lastGenerateTick = channel.self.sampleTick
				} else {
					samples[i][0] = channel.self.lastGenerate * channel.self.volume
					samples[i][1] = channel.self.lastGenerate * channel.self.volume
				}
			}

			channel.self.duration -= secondPerTick
		} else {
			samples[i][0] = 0
			samples[i][1] = 0
		}

		channel.self.Envelope()
		channel.self.Sweep()
	}
	return len(samples), true
}

func (channel *Channel) Sweep() {
	channel.sweepTick += 1 / 44100.0
	if channel.sweepNumber > 0 {
		if channel.sweepTick-channel.lastSweep >= sweepTime[channel.sweepTime] {
			if channel.Freq > 0 {
				newFreq := 0
				if channel.sweepIncrease {
					newFreq = channel.freqLast + channel.freqLast/2 ^ int(channel.sweepNumber)
				} else {
					newFreq = channel.freqLast - channel.freqLast/2 ^ int(channel.sweepNumber)
				}
				channel.freqLast = newFreq
				channel.Freq = 131072 / (2048 - int(newFreq))
				channel.lastSweep = channel.sweepTick
			}
		}
	}
}

func (channel *Channel) Envelope() {
	channel.envelopeTick += 1 / 44100.0
	if channel.envelopeSweepNum > 0 {
		step := float64(channel.envelopeSweepNum) * (1.0 / 64)
		if channel.envelopeTick-channel.lastEnvelope >= step {
			if channel.envelopeInitial > 0 {
				channel.envelopeInitial--
				if channel.envelopeIncrease {
					channel.volume = 1 - float64(channel.envelopeInitial)/float64(0xf)
				} else {
					channel.volume = float64(channel.envelopeInitial) / float64(0xf)
				}

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
	channel.sampleTick = 0
	channel.lastEnvelope = 0
}

func (channel Channel) Err() error {
	return nil
}
