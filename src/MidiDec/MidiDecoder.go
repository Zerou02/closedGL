package mididec

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

type Track struct {
	events                 []Event
	timeElapsedSinceLastEv float32
}
type Song struct {
	header HeaderChunk
	tracks []Track
}

type MidiDecoder struct {
	currentPtr         int
	data               []byte
	endOfTrack         bool
	runningStatus      byte
	song               Song
	tempoMSPerQn       uint32
	divisionTicksPerQn uint32
	timeElapsed        float32
	CurrNoteEvents     []Event
	test               []NoteOnEv
}

func NewMidiDecoder() MidiDecoder {
	return MidiDecoder{currentPtr: 0, song: Song{}, tempoMSPerQn: 120}
}

func (this *MidiDecoder) Deserialize(path string) {
	this.currentPtr = 0
	var bytes, _ = os.ReadFile(path)
	//format 0
	var data = []byte{
		//MThd
		0x4D, 0x54, 0x68, 0x64,
		0x00, 0x00, 0x00, 0x06,
		0x00, 0x00,
		0x00, 0x01,
		0x00, 0x60,
		//MTrk
		0x4D, 0x54, 0x72, 0x6B,
		0x00, 0x00, 0x00, 0x3B,
		//Events
		0x00, 0xFF, 0x58, 0x04, 0x04, 0x02, 0x18, 0x08,
		0x00, 0xFF, 0x51, 0x03, 0x07, 0xA1, 0x20,
		0x00, 0xC0, 0x05,
		0x00, 0xC1, 0x2E,
		0x00, 0xC2, 0x46,
		0x00, 0x92, 0x30, 0x60,
		0x00, 0x3C, 0x60,
		0x60, 0x91, 0x43, 0x40,
		0x60, 0x90, 0x4C, 0x20,
		0x81, 0x40, 0x82, 0x30, 0x40,
		0x00, 0x3C, 0x40,
		0x00, 0x81, 0x43, 0x40,
		0x00, 0x80, 0x4C, 0x40,
		0x00, 0xFF, 0x2F, 0x00,
	}
	//format 1
	var dataFormat1 = []byte{
		0x4D, 0x54, 0x68, 0x64,
		0x00, 0x00, 0x00, 0x06,
		0x00, 0x01,
		0x00, 0x04,
		0x00, 0x60,

		0x4D, 0x54, 0x72, 0x6B,
		0x00, 0x00, 0x00, 0x14,
		0x00, 0xFF, 0x58, 0x04, 0x04, 0x02, 0x18, 0x08,
		0x00, 0xFF, 0x51, 0x03, 0x07, 0xA1, 0x20,
		0x83, 0x00, 0xFF, 0x2F, 0x00,

		0x4D, 0x54, 0x72, 0x6B,
		0x00, 0x00, 0x00, 0x10,
		0x00, 0xC0, 0x05,
		0x81, 0x40, 0x90, 0x4C, 0x20,
		0x81, 0x40, 0x4C, 0x00,
		0x00, 0xFF, 0x2F, 0x00,

		0x4D, 0x54, 0x72, 0x6B,
		0x00, 0x00, 0x00, 0x0F,
		0x00, 0xC1, 0x2E,
		0x60, 0x91, 0x43, 0x40,
		0x82, 0x20, 0x43, 0x00,
		0x00, 0xFF, 0x2F, 0x00,

		0x4D, 0x54, 0x72, 0x6B,
		0x00, 0x00, 0x00, 0x15,

		0x00, 0xC2, 0x46,
		0x00, 0x92, 0x30, 0x60,
		0x00, 0x3C, 0x60,
		0x83, 0x00, 0x30, 0x00,
		0x00, 0x3C, 0x00,
		0x00, 0xFF, 0x2F, 0x00,
	}

	_ = data
	_ = dataFormat1
	_ = bytes
	this.data = bytes
	this.parse()
	//this.debugSave("test")
	this.convertNotes()
}

func (this *MidiDecoder) convertNotes() {
	//dt
	for i := 0; i < len(this.song.tracks); i++ {
		var x = &this.song.tracks[i]
		var newEvents = []Event{}
		for _, y := range x.events {
			if y.EvType == NoteOn || y.EvType == NoteOff || y.EvType == SetTempo {
				newEvents = append(newEvents, y)
			}
		}
		x.events = newEvents
	}
	var newTracks = []Track{}
	for _, x := range this.song.tracks {
		if len(x.events) != 0 {
			x.timeElapsedSinceLastEv = 0
			newTracks = append(newTracks, x)
		}
	}
	this.song.tracks = newTracks
}

func (this *MidiDecoder) Process(delta float32) {
	this.timeElapsed += delta
	this.CurrNoteEvents = []Event{}
	for i := 0; i < len(this.song.tracks); i++ {
		var currTrack = &this.song.tracks[i]
		if len(currTrack.events) == 0 {
			continue
		}
		currTrack.timeElapsedSinceLastEv += delta
		var nextDT = float32(this.dtToMs(currTrack.events[0].deltaTime)) / 1000.0
		if currTrack.timeElapsedSinceLastEv > float32(nextDT) {
			var evToHandle = currTrack.events[0]
			currTrack.events = currTrack.events[1:]
			currTrack.timeElapsedSinceLastEv = 0

			if evToHandle.EvType == SetTempo {
				this.parseSetTempoEv(evToHandle.Event)
			} else if evToHandle.EvType == NoteOn {
				this.CurrNoteEvents = append(this.CurrNoteEvents, evToHandle)
			} else if evToHandle.EvType == NoteOff {
				this.CurrNoteEvents = append(this.CurrNoteEvents, evToHandle)
			}
		}
	}

}

func (this *MidiDecoder) parseSetTempoEv(ev SingleEvent) {
	var bytes = ev.GetBytes()
	var val uint32 = 0
	val |= uint32(bytes[0])
	val <<= 8
	val |= uint32(bytes[1])
	val <<= 8
	val |= uint32(bytes[2])
	val <<= 8
	val |= uint32(bytes[3])
	this.tempoMSPerQn = uint32(val)
}

func (this *MidiDecoder) debugPrintEvents() {
	for i, x := range this.song.tracks {
		println("now parsing", i)
		for j, y := range x.events {
			print("event Nr", j)
			println(" ", y.EvType)
		}
	}
}

func (this *MidiDecoder) debugSave(path string) {
	var file, _ = os.Create(path + ".txt")
	for i, x := range this.song.tracks {
		file.WriteString("Chunk:" + strconv.FormatInt(int64(i), 10) + "\n\n")
		for j, y := range x.events {
			var s = "eventNr" + strconv.FormatInt(int64(j), 10) + ": " + strconv.FormatInt(int64(y.deltaTime), 10) + ":" + strconv.FormatInt(int64(y.EvType), 10) + "\n"
			if y.EvType == NoteOn {
				var bytes = y.Event.GetBytes()
				s += strconv.FormatInt(int64(bytes[1]), 10) + ", "
				s += strconv.FormatInt(int64(bytes[1]), 10) + "\n"
			} else if y.EvType == NoteOff {
				var bytes = y.Event.GetBytes()
				s += strconv.FormatInt(int64(bytes[1]), 10) + ", "
				s += strconv.FormatInt(int64(bytes[1]), 10) + "\n"
			}
			file.WriteString(s)
		}
	}
}

func (this *MidiDecoder) dtToMs(dt uint32) uint32 {
	return uint32((float32(dt) * (float32(this.tempoMSPerQn) / float32(this.divisionTicksPerQn))) / 1000.0)
}

func (this *MidiDecoder) parse() {
	this.parseHeader()
}

func (this *MidiDecoder) parseHeader() {
	var chunkType = this.readStr(4)
	var len = this.readUint32()
	_, _ = chunkType, len
	var format = this.readUint16()
	var nrTracks = this.readUint16()
	var division = this.readUint16()
	var deltaTicksPerQn uint16 = 0
	if division>>7 == 0x00 {
		deltaTicksPerQn = division
	} else {
		division = 24
		//panic("Division is type 1")
	}
	this.divisionTicksPerQn = uint32(deltaTicksPerQn)
	this.song.header = HeaderChunk{format: format, nrTracks: nrTracks, division: division, deltaTicksPerQn: deltaTicksPerQn}
	for i := 0; i < int(this.song.header.nrTracks); i++ {
		this.endOfTrack = false
		this.parseTrackChunk()
	}
}

func (this *MidiDecoder) parseTrackChunk() {
	this.song.tracks = append(this.song.tracks, Track{})
	var chunkType = this.readStr(4)
	var len = this.readUint32()
	_, _ = chunkType, len
	for !this.endOfTrack {
		this.parseEvent()
	}
}

func (this *MidiDecoder) parseEvent() {
	var currTrack = &this.song.tracks[len(this.song.tracks)-1]
	currTrack.events = append(currTrack.events, Event{deltaTime: this.readVLQ()})

	var evType = this.readByte()
	if evType == 0xF0 {
		this.parseNormalSysexEvent()
	} else if evType == 0xF7 {
		this.parseEscapeSysexEvent()
	} else if evType == 0xFF {
		this.parseMetaEvent()
	} else {
		this.parseMidiEvent()
	}
}

func (this *MidiDecoder) goBack(n int) {
	this.currentPtr -= n
}

// Nur hier running status
func (this *MidiDecoder) parseMidiEvent() {
	var currTrack = &this.song.tracks[len(this.song.tracks)-1]
	var currEv = &currTrack.events[len(currTrack.events)-1]
	this.goBack(1)
	var evType = this.readByte()
	//msb = 0
	if evType>>7&0x01 == 0 {
		evType = this.runningStatus
		this.goBack(1)
	}
	this.runningStatus = evType
	if evType >= 0xC0 && evType <= 0xCF {
		var channel = evType - 0xC0 + 1
		var ev = ProgramChangeEv{channel: channel, changeNr: this.readByte()}
		currEv.EvType = ProgramChange
		currEv.Event = ev
	} else if evType >= 0x90 && evType <= 0x9F {
		var channel = evType - 0x90 + 1
		var note = this.readByte()
		var vel = this.readByte()
		var noteOnEv = NoteOnEv{channel: channel, note: note, vel: vel}
		currEv.EvType = NoteOn
		currEv.Event = noteOnEv
		this.test = append(this.test, noteOnEv)
	} else if evType >= 0x80 && evType <= 0x8F {
		var channel = evType - 0x80 + 1
		var note = this.readByte()
		var vel = this.readByte()
		var ev = NoteOffEv{channel: channel, note: note, vel: vel}
		currEv.EvType = NoteOff
		currEv.Event = ev
	} else if evType >= 0xA0 && evType <= 0xAF {
		var channel = evType - 0xA0 + 1
		var ev = PolyphonicKeyPressureEv{channel: channel, note: this.readByte(), pressure: this.readByte()}
		currEv.EvType = PolyphonicKeyPressure
		currEv.Event = ev
	} else if evType >= 0xB0 && evType <= 0xBF {
		var channel = evType - 0xB0 + 1
		var ev = ControllerChangeEv{channel: channel, controller: this.readByte(), value: this.readByte()}
		currEv.EvType = ControllerChange
		currEv.Event = ev
	} else if evType >= 0xE0 && evType <= 0xEF {
		var channel = evType - 0xE0 + 1
		var val uint16 = 0
		var lower = this.readByte()
		val |= uint16(this.readByte())
		val <<= 7
		val |= uint16(lower)
		var ev = PitchBendChangeEv{channel: channel, value: val}
		currEv.EvType = PitchBendChange
		currEv.Event = ev
	} else {
		printHex(evType)
		panic("not implemented midi event")
	}
}

func (this *MidiDecoder) parseNormalSysexEvent() {
	printHex(this.readByte())
	panic("not implemented normal sysex event")
}

func (this *MidiDecoder) parseEscapeSysexEvent() {
	printHex(this.readByte())
	panic("not implemented escape sysex event")
}
func (this *MidiDecoder) parseMetaEvent() {
	var currTrack = &this.song.tracks[len(this.song.tracks)-1]
	var currEv = &currTrack.events[len(currTrack.events)-1]

	var evType = this.readByte()
	//Time Signature
	if evType == 0x58 {
		var len = this.readByte()
		_ = len
		var numerator = this.readByte()
		var denom = byte(math.Pow(2, float64(this.readByte())))
		//24 clocks ^= 1 Viertelnote
		var clocksPerDottedQuarter = this.readByte()
		//  1 Viertelnote ^= x 32stel
		var thirtytwothNotesPerQuarter = this.readByte()
		var ev = TimeSignatureEv{num: numerator, denom: denom, midiClocksPerTick: clocksPerDottedQuarter, nrOf32thPerQn: thirtytwothNotesPerQuarter}
		currEv.EvType = TimeSignature
		currEv.Event = ev
		//Set Tempo: Kompabilitätsevent
	} else if evType == 0x51 {
		var len = this.readByte()
		_ = len
		var msPerQuarter = this.readUint24()

		var ev = SetTempoEv{msPerQn: msPerQuarter}
		this.tempoMSPerQn = msPerQuarter
		currEv.EvType = SetTempo
		currEv.Event = ev
	} else if evType == 0x2F {
		var _ = this.readByte()
		this.endOfTrack = true
		currEv.EvType = EndOfTrack
		currEv.Event = nil
	} else if evType == 0x03 {
		var len = this.readByte()
		var ev = SeqTrkNameEv{text: this.readStr(int(len))}
		currEv.Event = ev
		currEv.EvType = SeqTrkName
	} else if evType == 0x04 {
		var len = this.readByte()
		var ev = InstrumentNameEv{text: this.readStr(int(len))}
		currEv.EvType = InstrumentName
		currEv.Event = ev
	} else {
		printHex(evType)
		panic("not implemented meta event")
	}
}

func parseDivision(division uint16) {
	//TODO:implement
}

func (this *MidiDecoder) readStr(len int) string {
	var retStr = ""
	for i := 0; i < len; i++ {
		retStr += string(this.data[this.currentPtr])
		this.currentPtr++
	}
	return retStr
}

func (this *MidiDecoder) readUint24() uint32 {
	var retInt uint32 = 0
	for i := 0; i < 3; i++ {
		retInt |= uint32(this.data[this.currentPtr])
		if i < 2 {
			retInt <<= 8
		}
		this.currentPtr++
	}
	return retInt
}

func (this *MidiDecoder) readUint16() uint16 {
	var retInt uint16 = 0
	for i := 0; i < 2; i++ {
		retInt |= uint16(this.data[this.currentPtr])
		if i < 1 {
			retInt <<= 8
		}
		this.currentPtr++
	}
	return retInt
}

func (this *MidiDecoder) readByte() byte {
	var retByte = this.data[this.currentPtr]
	this.currentPtr++
	return retByte
}

func (this *MidiDecoder) readUint32() uint32 {
	var retInt uint32 = 0
	for i := 0; i < 4; i++ {
		retInt |= uint32(this.data[this.currentPtr])
		if i < 3 {
			retInt <<= 8
		}
		this.currentPtr++
	}
	return retInt
}

func (this *MidiDecoder) readVLQ() uint32 {
	var bytes = []byte{}
	for this.data[this.currentPtr]>>7 == 1 {
		bytes = append(bytes, this.data[this.currentPtr])
		this.currentPtr++
	}
	bytes = append(bytes, this.data[this.currentPtr])
	this.currentPtr++
	return readVLQ(bytes)
}

func printHex(b byte) {
	fmt.Printf("0x%X \n", b)
}

func writeVLQ(value int64) int64 {
	var buffer int64 = 0
	buffer &= 0x7f
	for (value >> 7) > 0 {
		value = value >> 7
		buffer <<= 8
		buffer |= 0x80
		buffer += (value & 0x7f)
	}
	return buffer
}

func readVLQ(bytes []byte) uint32 {
	var retVal uint32 = 0
	for _, x := range bytes {
		if (x >> 7) == 1 {
			retVal |= uint32(x & 0x7F)
			retVal <<= 7
		} else {
			break
		}
	}
	if retVal == 0 && len(bytes) == 1 {
		retVal = uint32(bytes[0])
	}
	return retVal
}

func (this *MidiDecoder) assert(val, equals int64) {
	if val != equals {
		fmt.Printf("Error: assert failed. Should be %d, but got %d", equals, val)
	}
}
