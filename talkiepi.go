package talkiepi

import (
	"crypto/tls"

	"github.com/stianeikeland/go-rpio"
	"github.com/dchote/gumble/gumble"
	"github.com/dchote/gumble/gumbleopenal"
)

// Raspberry Pi GPIO pin assignments (CPU pin definitions)
const (
	OnlineLEDPin       uint = 26  // real GPIO.25
	ParticipantsLEDPin uint = 0   // real GPIO.27
	TransmitLEDPin     uint = 16  // not used now
	ButtonPin          uint = 22  // real GPIO.3
)

type Talkiepi struct {
	Config *gumble.Config
	Client *gumble.Client

	Address   string
	TLSConfig tls.Config

	ConnectAttempts uint

	Stream *gumbleopenal.Stream

	ChannelName    string
	IsConnected    bool
	IsTransmitting bool

	GPIOEnabled     bool
	OnlineLED       rpio.Pin
	ParticipantsLED rpio.Pin
	TransmitLED     rpio.Pin
	Button		rpio.Pin
	ButtonState     rpio.State
}
