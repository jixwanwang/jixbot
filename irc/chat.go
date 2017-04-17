package irc

import "fmt"

type TwitchChatClient struct {
	channel string
	main    *Client
	fancy   *Client
	whisper *Client
}

func NewTwitchChat(channel string, irc, ircFancy, ircWhisper *Client) *TwitchChatClient {
	return &TwitchChatClient{
		channel: channel,
		main:    irc,
		fancy:   ircFancy,
		whisper: ircWhisper,
	}
}

func (C *TwitchChatClient) ReadLoops() (chan Event, chan Event) {
	return C.main.ReadLoop(), C.whisper.ReadLoop()
}

func (C *TwitchChatClient) Say(msg string) {
	C.main.Say("#"+C.channel, msg)
}

func (C *TwitchChatClient) FancySay(msg string) {
	C.fancy.Say("#"+C.channel, msg)
}

func (C *TwitchChatClient) Whisper(to, msg string) {
	C.whisper.Say("#"+C.channel, fmt.Sprintf("/w %s %s", to, msg))
}

func (C *TwitchChatClient) Reload() error {
	err := C.main.Reload()
	if err != nil {
		return err
	}

	err = C.fancy.Reload()
	if err != nil {
		return err
	}

	err = C.whisper.Reload()
	if err != nil {
		return err
	}

	return nil
}
