package robot

import (
	"regexp"

	"github.com/eatmoreapple/openwechat"
)

type Message struct {
	*openwechat.Message
}

type User struct {
	*openwechat.User
}

func (m *Message) Sender() (*User, error) {
	sender, err := m.Message.Sender()
	return &User{sender}, err
}

func (m *Message) SenderInGroup() (*User, error) {
	group, err := m.Message.SenderInGroup()
	return &User{group}, err
}

func (m *Message) Receiver() (*User, error) {
	receiver, err := m.Message.Receiver()
	return &User{receiver}, err
}

func (m *Message) MatchTextCommand(commands []string) bool {
	if m.IsText() {
		for i := range commands {
			if commands[i] == m.Content {
				return true
			}
		}
	}
	return false
}

func (m *Message) MatchRegexCommand(commands []string) bool {
	if m.IsText() {
		for i := range commands {
			re := regexp.MustCompile(commands[i])
			return re.MatchString(m.Content)
		}
	}
	return false
}
