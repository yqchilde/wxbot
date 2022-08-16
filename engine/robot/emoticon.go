package robot

import "encoding/xml"

// Emoticon 表情包类型，并不是emoji的类型
type Emoticon struct {
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
}

func UnMarshalForEmoticon(data string) (emoticon *Emoticon, err error) {
	if err := xml.Unmarshal([]byte(data), &emoticon); err != nil {
		return nil, err
	}
	return emoticon, nil
}
