package ghmonitor

import "encoding/xml"

type SubscriptionMsgModel struct {
	XMLName xml.Name `xml:"msg"`
	Text    string   `xml:",chardata"`
	Appmsg  struct {
		Text        string `xml:",chardata"`
		Appid       string `xml:"appid,attr"`
		Sdkver      string `xml:"sdkver,attr"`
		Title       string `xml:"title"`
		Des         string `xml:"des"`
		Action      string `xml:"action"`
		Type        string `xml:"type"`
		Showtype    string `xml:"showtype"`
		Content     string `xml:"content"`
		Contentattr string `xml:"contentattr"`
		URL         string `xml:"url"`
		Lowurl      string `xml:"lowurl"`
		Appattach   struct {
			Text     string `xml:",chardata"`
			Totallen string `xml:"totallen"`
			Attachid string `xml:"attachid"`
			Fileext  string `xml:"fileext"`
		} `xml:"appattach"`
		Extinfo  string `xml:"extinfo"`
		Mmreader struct {
			Text     string `xml:",chardata"`
			Category struct {
				Text   string `xml:",chardata"`
				Type   string `xml:"type,attr"`
				Count  string `xml:"count,attr"`
				Name   string `xml:"name"`
				Topnew struct {
					Text   string `xml:",chardata"`
					Cover  string `xml:"cover"`
					Width  string `xml:"width"`
					Height string `xml:"height"`
					Digest string `xml:"digest"`
				} `xml:"topnew"`
				Item struct {
					Text         string `xml:",chardata"`
					Itemshowtype string `xml:"itemshowtype"`
					Title        string `xml:"title"`
					URL          string `xml:"url"`
					Shorturl     string `xml:"shorturl"`
					Longurl      string `xml:"longurl"`
					PubTime      string `xml:"pub_time"`
					Summary      string `xml:"summary"`
					Cover        string `xml:"cover"`
					Tweetid      string `xml:"tweetid"`
					Digest       string `xml:"digest"`
					Fileid       string `xml:"fileid"`
					Sources      struct {
						Text   string `xml:",chardata"`
						Source struct {
							Text string `xml:",chardata"`
							Name string `xml:"name"`
						} `xml:"source"`
					} `xml:"sources"`
					Styles         string `xml:"styles"`
					NativeURL      string `xml:"native_url"`
					DelFlag        string `xml:"del_flag"`
					Contentattr    string `xml:"contentattr"`
					PlayLength     string `xml:"play_length"`
					PlayURL        string `xml:"play_url"`
					VoiceID        string `xml:"voice_id"`
					Player         string `xml:"player"`
					MusicSource    string `xml:"music_source"`
					PicNum         string `xml:"pic_num"`
					Vid            string `xml:"vid"`
					Author         string `xml:"author"`
					Recommendation string `xml:"recommendation"`
					PicUrls        string `xml:"pic_urls"`
					CommentTopicID string `xml:"comment_topic_id"`
					Cover2351      string `xml:"cover_235_1"`
					Cover11        string `xml:"cover_1_1"`
					Cover169       string `xml:"cover_16_9"`
					AppmsgLikeType string `xml:"appmsg_like_type"`
					VideoWidth     string `xml:"video_width"`
					VideoHeight    string `xml:"video_height"`
					IsPaySubscribe string `xml:"is_pay_subscribe"`
					FinderFeed     struct {
						Text          string `xml:",chardata"`
						ObjectID      string `xml:"object_id"`
						ObjectNonceID string `xml:"object_nonce_id"`
						FeedType      string `xml:"feed_type"`
						Nickname      string `xml:"nickname"`
						Avatar        string `xml:"avatar"`
						Desc          string `xml:"desc"`
						MediaCount    string `xml:"media_count"`
						MediaList     string `xml:"media_list"`
						MegaVideo     struct {
							Text          string `xml:",chardata"`
							ObjectID      string `xml:"object_id"`
							ObjectNonceID string `xml:"object_nonce_id"`
						} `xml:"mega_video"`
					} `xml:"finder_feed"`
					FinderLive struct {
						Text              string `xml:",chardata"`
						FinderUsername    string `xml:"finder_username"`
						Category          string `xml:"category"`
						FinderNonceID     string `xml:"finder_nonce_id"`
						ExportID          string `xml:"export_id"`
						Nickname          string `xml:"nickname"`
						HeadURL           string `xml:"head_url"`
						Desc              string `xml:"desc"`
						LiveStatus        string `xml:"live_status"`
						LiveSourceTypeStr string `xml:"live_source_type_str"`
						ExtFlag           string `xml:"ext_flag"`
						AuthIconURL       string `xml:"auth_icon_url"`
						AuthIconTypeStr   string `xml:"auth_icon_type_str"`
						Media             struct {
							Text     string `xml:",chardata"`
							CoverURL string `xml:"cover_url"`
							Height   string `xml:"height"`
							Width    string `xml:"width"`
						} `xml:"media"`
					} `xml:"finder_live"`
				} `xml:"item"`
			} `xml:"category"`
			Publisher struct {
				Text     string `xml:",chardata"`
				Username string `xml:"username"`
				Nickname string `xml:"nickname"`
			} `xml:"publisher"`
			TemplateHeader string `xml:"template_header"`
			TemplateDetail string `xml:"template_detail"`
			ForbidForward  string `xml:"forbid_forward"`
		} `xml:"mmreader"`
		Thumburl string `xml:"thumburl"`
	} `xml:"appmsg"`
	Fromusername string `xml:"fromusername"`
	Appinfo      struct {
		Text          string `xml:",chardata"`
		Version       string `xml:"version"`
		Appname       string `xml:"appname"`
		Isforceupdate string `xml:"isforceupdate"`
	} `xml:"appinfo"`
}

// SliceUnion 返回切片并集
func SliceUnion(a, b []string) []string {
	m := make(map[string]bool)
	for _, v := range a {
		m[v] = true
	}
	for _, v := range b {
		m[v] = true
	}
	var r []string
	for k := range m {
		r = append(r, k)
	}
	return r
}
