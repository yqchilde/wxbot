package qianxun

import "encoding/xml"

// EmojiXml emoji xml
type EmojiXml struct {
	XMLName xml.Name `xml:"msg"`
	Text    string   `xml:",chardata"`
	Emoji   struct {
		Text              string `xml:",chardata"`
		Fromusername      string `xml:"fromusername,attr"`
		Tousername        string `xml:"tousername,attr"`
		Type              string `xml:"type,attr"`
		Idbuffer          string `xml:"idbuffer,attr"`
		Md5               string `xml:"md5,attr"`
		Len               string `xml:"len,attr"`
		Productid         string `xml:"productid,attr"`
		Androidmd5        string `xml:"androidmd5,attr"`
		Androidlen        string `xml:"androidlen,attr"`
		S60v3md5          string `xml:"s60v3md5,attr"`
		S60v3len          string `xml:"s60v3len,attr"`
		S60v5md5          string `xml:"s60v5md5,attr"`
		S60v5len          string `xml:"s60v5len,attr"`
		Cdnurl            string `xml:"cdnurl,attr"`
		Designerid        string `xml:"designerid,attr"`
		Thumburl          string `xml:"thumburl,attr"`
		Encrypturl        string `xml:"encrypturl,attr"`
		Aeskey            string `xml:"aeskey,attr"`
		Externurl         string `xml:"externurl,attr"`
		Externmd5         string `xml:"externmd5,attr"`
		Width             string `xml:"width,attr"`
		Height            string `xml:"height,attr"`
		Tpurl             string `xml:"tpurl,attr"`
		Tpauthkey         string `xml:"tpauthkey,attr"`
		Attachedtext      string `xml:"attachedtext,attr"`
		Attachedtextcolor string `xml:"attachedtextcolor,attr"`
		Lensid            string `xml:"lensid,attr"`
		Emojiattr         string `xml:"emojiattr,attr"`
		Linkid            string `xml:"linkid,attr"`
		Desc              string `xml:"desc,attr"`
	} `xml:"emoji"`
	Gameext struct {
		Text    string `xml:",chardata"`
		Type    string `xml:"type,attr"`
		Content string `xml:"content,attr"`
	} `xml:"gameext"`
}

// ReferenceXml 引用消息xml
type ReferenceXml struct {
	XMLName xml.Name `xml:"msg"`
	Text    string   `xml:",chardata"`
	Appmsg  struct {
		Text          string `xml:",chardata"`
		Appid         string `xml:"appid,attr"`
		Sdkver        string `xml:"sdkver,attr"`
		Title         string `xml:"title"`
		Des           string `xml:"des"`
		Action        string `xml:"action"`
		Type          string `xml:"type"`
		Showtype      string `xml:"showtype"`
		Soundtype     string `xml:"soundtype"`
		Mediatagname  string `xml:"mediatagname"`
		Messageext    string `xml:"messageext"`
		Messageaction string `xml:"messageaction"`
		Content       string `xml:"content"`
		Contentattr   string `xml:"contentattr"`
		URL           string `xml:"url"`
		Lowurl        string `xml:"lowurl"`
		Dataurl       string `xml:"dataurl"`
		Lowdataurl    string `xml:"lowdataurl"`
		Appattach     struct {
			Text        string `xml:",chardata"`
			Totallen    string `xml:"totallen"`
			Attachid    string `xml:"attachid"`
			Emoticonmd5 string `xml:"emoticonmd5"`
			Fileext     string `xml:"fileext"`
			Aeskey      string `xml:"aeskey"`
		} `xml:"appattach"`
		Extinfo           string `xml:"extinfo"`
		Sourceusername    string `xml:"sourceusername"`
		Sourcedisplayname string `xml:"sourcedisplayname"`
		Thumburl          string `xml:"thumburl"`
		Md5               string `xml:"md5"`
		Statextstr        string `xml:"statextstr"`
		Refermsg          *struct {
			Text        string `xml:",chardata"`
			Type        string `xml:"type"`
			Svrid       string `xml:"svrid"`
			Fromusr     string `xml:"fromusr"`
			Chatusr     string `xml:"chatusr"`
			Displayname string `xml:"displayname"`
			Content     string `xml:"content"`
		} `xml:"refermsg"`
	} `xml:"appmsg"`
	Fromusername string `xml:"fromusername"`
	Scene        string `xml:"scene"`
	Appinfo      struct {
		Text    string `xml:",chardata"`
		Version string `xml:"version"`
		Appname string `xml:"appname"`
	} `xml:"appinfo"`
	Commenturl string `xml:"commenturl"`
}
