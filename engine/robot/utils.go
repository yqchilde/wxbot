package robot

import (
	"bytes"
	"strconv"
	"unicode/utf16"
	"unicode/utf8"
)

func formatTextMessage(msg string) string {
	buff := bytes.NewBuffer(make([]byte, 0, len(msg)*2))
	for _, r := range msg {
		switch utf8.RuneLen(r) {
		case 1, 2, 3:
			buff.WriteString(string(r))
		case 4:
			r1, r2 := utf16.EncodeRune(r)
			buff.WriteString(`[emoji=\u`)
			buff.WriteString(strconv.FormatInt(int64(r1), 16))
			buff.WriteString(`\u`)
			buff.WriteString(strconv.FormatInt(int64(r2), 16) + `]`)
		}
	}
	return buff.String()
}
