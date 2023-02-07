package qianxun

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/antchfx/xmlquery"

	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
)

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

func (f *Framework) GetMemePictures(msg *robot.Message) string {
	doc, err := xmlquery.Parse(strings.NewReader(msg.Content))
	if err != nil {
		return ""
	}
	node, err := xmlquery.Query(doc, "//emoji")
	if err != nil {
		return ""
	}
	return node.SelectAttr("cdnurl")
}

func (f *Framework) SendText(toWxId, text string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0001",
		"data": map[string]interface{}{
			"wxid": toWxId,
			"msg":  f.msgFormat(text),
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[千寻] SendText error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendTextAndAt(toGroupWxId, toWxId, toWxName, text string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0001",
		"data": map[string]interface{}{
			"wxid": toGroupWxId,
			"msg":  fmt.Sprintf("[@,wxid=%s,nick=%s,isAuto=true] %s", toWxId, toWxName, f.msgFormat(text)),
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[千寻] SendTextAndAt error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendImage(toWxId, path string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0010",
		"data": map[string]interface{}{
			"wxid": toWxId,
			"path": path,
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[千寻] SendImage error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendShareLink(toWxId, title, desc, imageUrl, jumpUrl string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
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

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[千寻] SendShareLink error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendFile(toWxId, path string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0011",
		"data": map[string]interface{}{
			"wxid": toWxId,
			"path": path,
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[千寻] SendFile error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendVideo(toWxId, path string) error {
	log.Errorf("[千寻] SendVideo not support")
	return errors.New("SendVideo not support")
}

func (f *Framework) SendEmoji(toWxId, path string) error {
	log.Errorf("[千寻] SendEmoji not support")
	return errors.New("SendEmoji not support")
}

func (f *Framework) SendMusic(toWxId, name, author, app, jumpUrl, musicUrl, coverUrl string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0014",
		"data": map[string]interface{}{
			"wxid":     toWxId,
			"name":     name,
			"author":   author,
			"app":      app,
			"jumpUrl":  jumpUrl,
			"musicUrl": musicUrl,
			"imageUrl": coverUrl,
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[千寻] SendMusic error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendMiniProgram(toWxId, ghId, title, content, imagePath, jumpPath string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0013",
		"data": map[string]interface{}{
			"wxid":     toWxId,
			"title":    title,
			"content":  content,
			"jumpPath": jumpPath,
			"gh":       ghId,
			"path":     imagePath,
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[千寻] SendMiniProgram error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendMessageRecord(toWxId, title string, dataList []map[string]interface{}) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0009",
		"data": map[string]interface{}{
			"wxid":     toWxId,
			"title":    title,
			"dataList": dataList,
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[千寻] SendMessageRecord error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendMessageRecordXML(toWxId, xmlStr string) error {
	log.Errorf("[千寻] SendMessageRecordXML not support")
	return errors.New("SendMessageRecordXML not support, please use SendMessageRecord")
}

func (f *Framework) SendFavorites(toWxId, favoritesId string) error {
	log.Errorf("[千寻] SendFavorites not support")
	return errors.New("SendFavorites not support")
}

func (f *Framework) SendXML(toWxId, xmlStr string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0015",
		"data": map[string]interface{}{
			"wxid": toWxId,
			"xml":  xmlStr,
		},
	}

	buf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(payload)
	if err := NewRequest().Post(apiUrl).SetBodyJsonString(buf.String()).Do().Err; err != nil {
		log.Errorf("[千寻] SendXML error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendBusinessCard(toWxId, targetWxId string) error {
	info, err := f.GetObjectInfo(targetWxId)
	if err != nil {
		log.Errorf("[千寻] SendBusinessCard error: %v", err)
		return err
	}

	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0025",
		"data": map[string]interface{}{
			"wxid": toWxId,
			"xml":  fmt.Sprintf(`<msg username="%s" nickname="%s" />`, targetWxId, info.Nick),
		},
	}

	buf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(payload)
	if err := NewRequest().Post(apiUrl).SetBodyJsonString(buf.String()).Do().Err; err != nil {
		log.Errorf("[千寻] SendBusinessCard error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) AgreeFriendVerify(v3, v4, scene string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0017",
		"data": map[string]interface{}{
			"scene": scene,
			"v3":    v3,
			"v4":    v4,
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[千寻] AgreeFriendVerify error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) InviteIntoGroup(groupWxId, wxId string, typ int) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0021",
		"data": map[string]interface{}{
			"wxid":    groupWxId,
			"objWxid": wxId,
			"type":    typ,
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[千寻] InviteIntoGroup error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) GetObjectInfo(wxId string) (*robot.ObjectInfo, error) {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0004",
		"data": map[string]interface{}{
			"wxid": wxId,
		},
	}

	var dataResp ObjectInfoResp
	if err := NewRequest().Post(apiUrl).SetBody(payload).SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[千寻] GetObjectInfo error: %v", err.Error())
		return nil, err
	}
	return &robot.ObjectInfo{
		WxId:                    dataResp.Result.Wxid,
		WxNum:                   dataResp.Result.WxNum,
		Nick:                    dataResp.Result.Nick,
		Remark:                  dataResp.Result.Remark,
		NickBrief:               dataResp.Result.NickBrief,
		NickWhole:               dataResp.Result.NickWhole,
		RemarkBrief:             dataResp.Result.RemarkBrief,
		RemarkWhole:             dataResp.Result.RemarkWhole,
		EnBrief:                 dataResp.Result.EnBrief,
		EnWhole:                 dataResp.Result.EnWhole,
		V3:                      dataResp.Result.V3,
		V4:                      dataResp.Result.V4,
		Sign:                    dataResp.Result.Sign,
		Country:                 dataResp.Result.Country,
		Province:                dataResp.Result.Province,
		City:                    dataResp.Result.City,
		MomentsBackgroundImgUrl: dataResp.Result.MomentsBackgroudImgUrl,
		AvatarMinUrl:            dataResp.Result.AvatarMinUrl,
		AvatarMaxUrl:            dataResp.Result.AvatarMaxUrl,
		Sex:                     dataResp.Result.Sex,
		MemberNum:               dataResp.Result.MemberNum,
	}, nil
}

func (f *Framework) GetFriendsList(isRefresh bool) ([]*robot.FriendInfo, error) {
	dataType := 1
	if isRefresh {
		dataType = 2
	}

	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0005",
		"data": map[string]interface{}{
			"type": dataType,
		},
	}

	var dataResp FriendsListResp
	if err := NewRequest().Post(apiUrl).SetBody(payload).SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[千寻] GetFriendsList error: %v", err.Error())
		return nil, err
	}
	var friendsInfoList []*robot.FriendInfo
	for i := range dataResp.Result {
		friendsInfoList = append(friendsInfoList, &robot.FriendInfo{
			WxId:                    dataResp.Result[i].Wxid,
			WxNum:                   dataResp.Result[i].WxNum,
			Nick:                    dataResp.Result[i].Nick,
			Remark:                  dataResp.Result[i].Remark,
			NickBrief:               dataResp.Result[i].NickBrief,
			NickWhole:               dataResp.Result[i].NickWhole,
			RemarkBrief:             dataResp.Result[i].RemarkBrief,
			RemarkWhole:             dataResp.Result[i].RemarkWhole,
			EnBrief:                 dataResp.Result[i].EnBrief,
			EnWhole:                 dataResp.Result[i].EnWhole,
			V3:                      dataResp.Result[i].V3,
			Sign:                    dataResp.Result[i].Sign,
			Country:                 dataResp.Result[i].Country,
			Province:                dataResp.Result[i].Province,
			City:                    dataResp.Result[i].City,
			MomentsBackgroundImgUrl: dataResp.Result[i].MomentsBackgroudImgUrl,
			AvatarMinUrl:            dataResp.Result[i].AvatarMinUrl,
			AvatarMaxUrl:            dataResp.Result[i].AvatarMaxUrl,
			Sex:                     dataResp.Result[i].Sex,
			MemberNum:               dataResp.Result[i].MemberNum,
		})
	}

	// 过滤系统用户
	var SystemUserWxId = map[string]struct{}{"medianote": {}, "newsapp": {}, "fmessage": {}, "floatbottle": {}}
	var filteredFriendInfo []*robot.FriendInfo
	for i := range friendsInfoList {
		if _, ok := SystemUserWxId[friendsInfoList[i].WxId]; !ok {
			filteredFriendInfo = append(filteredFriendInfo, friendsInfoList[i])
		}
	}
	return filteredFriendInfo, nil
}

func (f *Framework) GetGroupList(isRefresh bool) ([]*robot.GroupInfo, error) {
	dataType := 1
	if isRefresh {
		dataType = 2
	}

	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0006",
		"data": map[string]interface{}{
			"type": dataType,
		},
	}

	var dataResp GroupListResp
	if err := NewRequest().Post(apiUrl).SetBody(payload).SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[千寻] GetGroupList error: %v", err.Error())
		return nil, err
	}
	var groupInfoList []*robot.GroupInfo
	for i := range dataResp.Result {
		groupInfoList = append(groupInfoList, &robot.GroupInfo{
			WxId:        dataResp.Result[i].Wxid,
			WxNum:       dataResp.Result[i].WxNum,
			Nick:        dataResp.Result[i].Nick,
			Remark:      dataResp.Result[i].Remark,
			NickBrief:   dataResp.Result[i].NickBrief,
			NickWhole:   dataResp.Result[i].NickWhole,
			RemarkBrief: dataResp.Result[i].RemarkBrief,
			RemarkWhole: dataResp.Result[i].RemarkWhole,
			EnBrief:     dataResp.Result[i].EnBrief,
			EnWhole:     dataResp.Result[i].EnWhole,
		})
	}
	return groupInfoList, nil
}

func (f *Framework) GetSubscriptionList(isRefresh bool) ([]*robot.SubscriptionInfo, error) {
	dataType := 1
	if isRefresh {
		dataType = 2
	}

	apiUrl := fmt.Sprintf("%s/DaenWxHook/httpapi/?wxid=%s", f.ApiUrl, f.BotWxId)
	payload := map[string]interface{}{
		"type": "Q0007",
		"data": map[string]interface{}{
			"type": dataType,
		},
	}

	var dataResp SubscriptionListResp
	if err := NewRequest().Post(apiUrl).SetBody(payload).SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[千寻] GetSubscriptionList error: %v", err.Error())
		return nil, err
	}
	var subscriptionInfoList []*robot.SubscriptionInfo
	for i := range dataResp.Result {
		subscriptionInfoList = append(subscriptionInfoList, &robot.SubscriptionInfo{
			WxId:                    dataResp.Result[i].Wxid,
			WxNum:                   dataResp.Result[i].WxNum,
			Nick:                    dataResp.Result[i].Nick,
			Remark:                  dataResp.Result[i].Remark,
			NickBrief:               dataResp.Result[i].NickBrief,
			NickWhole:               dataResp.Result[i].NickWhole,
			RemarkBrief:             dataResp.Result[i].RemarkBrief,
			RemarkWhole:             dataResp.Result[i].RemarkWhole,
			EnBrief:                 dataResp.Result[i].EnBrief,
			EnWhole:                 dataResp.Result[i].EnWhole,
			V3:                      dataResp.Result[i].V3,
			Sign:                    dataResp.Result[i].Sign,
			Country:                 dataResp.Result[i].Country,
			Province:                dataResp.Result[i].Province,
			City:                    dataResp.Result[i].City,
			MomentsBackgroundImgUrl: dataResp.Result[i].MomentsBackgroudImgUrl,
			AvatarMinUrl:            dataResp.Result[i].AvatarMinUrl,
			AvatarMaxUrl:            dataResp.Result[i].AvatarMaxUrl,
		})
	}
	return subscriptionInfoList, nil
}
