package qianxun

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/imroc/req/v3"
	"github.com/yqchilde/pkgs/log"
)

type MessageResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (f *Framework) msgFormat(msg string) string {
	buff := bytes.NewBuffer(make([]byte, 0, len(msg)*2))
	for _, r := range msg {
		if unicode.Is(unicode.Han, r) || unicode.IsLetter(r) {
			buff.WriteString(string(r))
			continue
		}
		switch utf8.RuneLen(r) {
		case 2, 3:
			buff.WriteString(`[emoji=`)
			buff.WriteString(fmt.Sprintf("%04x", r) + `]`)
		case 4:
			r1, r2 := utf16.EncodeRune(r)
			buff.WriteString(`[emoji=`)
			buff.WriteString(fmt.Sprintf("%04x]", r1))
			buff.WriteString(`[emoji=`)
			buff.WriteString(fmt.Sprintf("%04x]", r2))
		default:
			buff.WriteString(string(r))
		}
	}
	return strings.ReplaceAll(strings.ReplaceAll(buff.String(), "\r\n", "\r"), "\n", "\r")
}

func (f *Framework) SendText(toWxId, text string) error {
	payload := map[string]interface{}{
		"type": "Q0001",
		"data": map[string]interface{}{
			"wxid": toWxId,
			"msg":  f.msgFormat(text),
		},
	}

	var msgResp MessageResp
	resp := req.C().Post(f.ApiUrl).SetBody(payload).Do()
	if err := resp.Into(&msgResp); err != nil {
		log.Errorf("[千寻] SendText error: %v", err)
		return err
	}
	if msgResp.Code != 200 {
		log.Errorf("[千寻] SendText error: %s", resp.String())
		return errors.New(msgResp.Msg)
	}
	return nil
}

func (f *Framework) SendTextAndAt(toGroupWxId, toWxId, toWxName, text string) error {
	payload := map[string]interface{}{
		"type": "Q0001",
		"data": map[string]interface{}{
			"wxid": toGroupWxId,
			"msg":  fmt.Sprintf("[@,wxid=%s,nick=%s,isAuto=true] %s", toWxId, toWxName, f.msgFormat(text)),
		},
	}

	var msgResp MessageResp
	resp := req.C().Post(f.ApiUrl).SetBody(payload).Do()
	if err := resp.Into(&msgResp); err != nil {
		log.Errorf("[千寻] SendTextAndAt error: %v", err)
		return err
	}
	if msgResp.Code != 200 {
		log.Errorf("[千寻] SendTextAndAt error: %s", resp.String())
		return errors.New(msgResp.Msg)
	}
	return nil
}

func (f *Framework) SendImage(toWxId, path string) error {
	payload := map[string]interface{}{
		"type": "Q0010",
		"data": map[string]interface{}{
			"wxid": toWxId,
			"path": path,
		},
	}

	var msgResp MessageResp
	resp := req.C().Post(f.ApiUrl).SetBody(payload).Do()
	if err := resp.Into(&msgResp); err != nil {
		log.Errorf("[千寻] SendImage error: %v", err)
		return err
	}
	if msgResp.Code != 200 {
		log.Errorf("[千寻] SendImage error: %s", resp.String())
		return errors.New(msgResp.Msg)
	}
	return nil
}

func (f *Framework) SendShareLink(toWxId, title, desc, imageUrl, jumpUrl string) error {
	payload := map[string]interface{}{
		"type": "Q0012",
		"data": map[string]interface{}{
			"wxid":    toWxId,
			"title":   title,
			"content": desc,
			"jumpUrl": jumpUrl,
			"path":    imageUrl,
		},
	}

	var msgResp MessageResp
	resp := req.C().Post(f.ApiUrl).SetBody(payload).Do()
	if err := resp.Into(&msgResp); err != nil {
		log.Errorf("[千寻] SendShareLink error: %v", err)
		return err
	}
	if msgResp.Code != 200 {
		log.Errorf("[千寻] SendShareLink error: %s", resp.String())
		return errors.New(msgResp.Msg)
	}
	return nil
}