package mididec

type HeaderChunk struct {
	format, nrTracks, division, deltaTicksPerQn uint16
}
type TimeSignatureEv struct {
	num, denom, midiClocksPerTick, nrOf32thPerQn byte
}

func (this TimeSignatureEv) GetBytes() []byte {
	return []byte{this.num, this.denom, this.midiClocksPerTick, this.nrOf32thPerQn}
}

type SetTempoEv struct {
	msPerQn uint32
}

func (this SetTempoEv) GetBytes() []byte {
	var lowest = byte(this.msPerQn & 0xFF)
	var lowest2 = byte((this.msPerQn >> 8) & 0xFF)
	var lowest3 = byte((this.msPerQn >> 16) & 0xFF)
	var lowest4 = byte((this.msPerQn >> 24) & 0xFF)

	return []byte{lowest4, lowest3, lowest2, lowest}
}

type ProgramChangeEv struct {
	channel, changeNr byte
}

func (this ProgramChangeEv) GetBytes() []byte {
	return []byte{}
}

type NoteOffEv struct {
	channel, note, vel byte
}

func (this NoteOffEv) GetBytes() []byte {
	return []byte{this.channel, this.note, this.vel}
}

type NoteOnEv struct {
	channel, note, vel byte
}

func (this NoteOnEv) GetBytes() []byte {
	return []byte{this.channel, this.note, this.vel}
}

type InstrumentNameEv struct {
	text string
}

func (this InstrumentNameEv) GetBytes() []byte {
	return []byte{}
}

type SeqTrkNameEv struct {
	text string
}

func (this SeqTrkNameEv) GetBytes() []byte {
	return []byte{}
}

type PitchBendChangeEv struct {
	channel byte
	value   uint16
}

func (this PitchBendChangeEv) GetBytes() []byte {
	return []byte{}
}

type EventType uint64

const (
	InstrumentName = iota
	NoteOff
	NoteOn
	SeqTrkName
	PitchBendChange
	SetTempo
	ProgramChange
	TimeSignature
	PolyphonicKeyPressure
	EndOfTrack
	ControllerChange
)

type Event struct {
	EvType    EventType
	Event     SingleEvent
	deltaTime uint32
}

type PolyphonicKeyPressureEv struct {
	channel, note, pressure byte
}

func (this PolyphonicKeyPressureEv) GetBytes() []byte {
	return []byte{}
}

type ControllerChangeEv struct {
	channel, controller, value byte
}

func (this ControllerChangeEv) GetBytes() []byte {
	return []byte{}
}

type SingleEvent interface {
	GetBytes() []byte
}
