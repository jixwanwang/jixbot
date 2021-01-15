package irc

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Client struct {
	channel  string
	oauth    string
	username string

	socket       net.Conn
	br           *bufio.Reader
	sentMessages []message
	// Max messages per 15 seconds
	messageRate int
	server      string
	closed      bool
}

type message struct {
	message string
	t       time.Time
}

type Event struct {
	Tags    map[string]string
	From    string
	Kind    string
	Message string
	Err     error
}

func New(server, channel, oauth, username string, messageRate int) *Client {
	client := &Client{
		channel:     channel,
		oauth:       oauth,
		username:    username,
		messageRate: messageRate,
		server:      server,
	}

	client.Reload()

	return client
}

func (C *Client) Reload() error {
	conn, err := net.Dial("tcp", C.server)
	br := bufio.NewReaderSize(conn, 4096)
	C.br = br
	C.socket = conn
	C.closed = false

	C.Send(fmt.Sprintf("PASS %s", C.oauth))
	C.Send(fmt.Sprintf("NICK %s", C.username))
	C.Send(fmt.Sprintf("JOIN #%s", C.channel))
	C.Send("CAP REQ :twitch.tv/commands")
	C.Send("CAP REQ :twitch.tv/membership")
	C.Send("CAP REQ :twitch.tv/tags")
	C.Send(fmt.Sprintf("JOIN #%s", C.username))

	return err
}

func (C *Client) Send(msg string) {
	C.socket.Write([]byte(msg + "\r\n"))
}

func (C *Client) Say(channel, msg string) {
	// Prune message list
	i := 0
	for i = range C.sentMessages {
		if time.Since(C.sentMessages[i].t) < 15*time.Second {
			break
		}
	}
	if i == len(C.sentMessages) {
		C.sentMessages = []message{}
	} else {
		C.sentMessages = C.sentMessages[i:]
	}

	if len(C.sentMessages) < C.messageRate {
		log.Printf(">PRIVMSG %s :%s", channel, msg)
		C.Send(fmt.Sprintf("PRIVMSG %s :%s", channel, msg))
		C.sentMessages = append(C.sentMessages, message{
			message: msg,
			t:       time.Now(),
		})
	}
}

func (C *Client) ReadLoop() chan Event {
	events := make(chan Event, 10)

	go func() {
		for {
			t := time.Now()
			t = t.Add(2 * time.Minute)

			C.socket.SetReadDeadline(t)

			msg, err := C.br.ReadString('\n')
			if err != nil {
				events <- Event{Err: err}
				continue
			}

			tags := map[string]string{}
			// Parse tags
			if msg[0:1] == "@" {
				space1 := strings.Index(msg, " ")
				tagString := msg[:space1]
				tagParts := strings.Split(tagString, ";")
				for _, s := range tagParts {
					i := strings.Index(s, "=")
					tags[s[:i]] = s[i+1:]
				}
				msg = msg[space1+1:]
			}

			space := strings.Index(msg, " ")
			if space == -1 {
				events <- Event{}
				continue
			}
			from := msg[:space]

			msg = msg[space+1:]
			space = strings.Index(msg, " ")
			if space == -1 {
				events <- Event{
					Tags:    tags,
					From:    from,
					Kind:    msg,
					Message: "",
				}
				continue
			}
			kind := msg[:space]

			message := strings.TrimSpace(msg[space:])

			events <- Event{
				Tags:    tags,
				From:    from,
				Kind:    kind,
				Message: message,
			}
		}
	}()

	return events
}
