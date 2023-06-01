package dean

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

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

func (f *Framework) GetRobotInfo() (*robot.User, error) {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0003",
		"data": map[string]interface{}{},
	}

	var dataResp RobotInfoResp
	if err := NewRequest().Post(apiUrl).SetBody(payload).SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[Dean] GetRobotInfo error: %v", err.Error())
		return nil, err
	}
	return &robot.User{
		WxId:         dataResp.Result.Wxid,
		WxNum:        dataResp.Result.WxNum,
		Nick:         dataResp.Result.Nick,
		Country:      dataResp.Result.Country,
		Province:     dataResp.Result.Province,
		City:         dataResp.Result.City,
		AvatarMinUrl: dataResp.Result.AvatarUrl,
		AvatarMaxUrl: dataResp.Result.AvatarUrl,
	}, nil
}

func (f *Framework) GetMemePictures(msg *robot.Message) string {
	var emoji EmojiXml
	if err := xml.Unmarshal([]byte(msg.Content), &emoji); err != nil {
		return ""
	}
	return emoji.Emoji.Cdnurl
}

func (f *Framework) SendText(toWxId, text string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0001",
		"data": map[string]interface{}{
			"wxid": toWxId,
			"msg":  f.msgFormat(text),
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[Dean] SendText error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendTextAndAt(toGroupWxId, toWxId, toWxName, text string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0001",
		"data": map[string]interface{}{
			"wxid": toGroupWxId,
			"msg":  fmt.Sprintf("[@,wxid=%s,nick=%s,isAuto=true] %s", toWxId, toWxName, f.msgFormat(text)),
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[Dean] SendTextAndAt error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendImage(toWxId, path string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0010",
		"data": map[string]interface{}{
			"wxid": toWxId,
			"path": path,
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[Dean] SendImage error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendShareLink(toWxId, title, desc, imageUrl, jumpUrl string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
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
		log.Errorf("[Dean] SendShareLink error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendFile(toWxId, path string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0011",
		"data": map[string]interface{}{
			"wxid": toWxId,
			"path": path,
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[Dean] SendFile error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendVideo(toWxId, path string) error {
	log.Errorf("[Dean] SendVideo not support")
	return errors.New("SendVideo not support")
}

func (f *Framework) SendEmoji(toWxId, path string) error {
	log.Errorf("[Dean] SendEmoji not support")
	return errors.New("SendEmoji not support")
}

func (f *Framework) SendMusic(toWxId, name, author, app, jumpUrl, musicUrl, coverUrl string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
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
		log.Errorf("[Dean] SendMusic error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendMiniProgram(toWxId, ghId, title, content, imagePath, jumpPath string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
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
		log.Errorf("[Dean] SendMiniProgram error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendMessageRecord(toWxId, title string, dataList []map[string]interface{}) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0009",
		"data": map[string]interface{}{
			"wxid":     toWxId,
			"title":    title,
			"dataList": dataList,
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[Dean] SendMessageRecord error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendMessageRecordXML(toWxId, xmlStr string) error {
	log.Errorf("[Dean] SendMessageRecordXML not support")
	return errors.New("SendMessageRecordXML not support, please use SendMessageRecord")
}

func (f *Framework) SendFavorites(toWxId, favoritesId string) error {
	log.Errorf("[Dean] SendFavorites not support")
	return errors.New("SendFavorites not support")
}

func (f *Framework) SendXML(toWxId, xmlStr string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
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
		log.Errorf("[Dean] SendXML error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendBusinessCard(toWxId, targetWxId string) error {
	info, err := f.GetObjectInfo(targetWxId)
	if err != nil {
		log.Errorf("[Dean] SendBusinessCard error: %v", err)
		return err
	}

	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
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
		log.Errorf("[Dean] SendBusinessCard error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) AgreeFriendVerify(v3, v4, scene string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0017",
		"data": map[string]interface{}{
			"scene": scene,
			"v3":    v3,
			"v4":    v4,
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[Dean] AgreeFriendVerify error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) InviteIntoGroup(groupWxId, wxId string, typ int) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0021",
		"data": map[string]interface{}{
			"wxid":    groupWxId,
			"objWxid": wxId,
			"type":    typ,
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[Dean] InviteIntoGroup error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) GetObjectInfo(wxId string) (*robot.User, error) {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0004",
		"data": map[string]interface{}{
			"wxid": wxId,
		},
	}

	var dataResp ObjectInfoResp
	if err := NewRequest().Post(apiUrl).SetBody(payload).SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[Dean] GetObjectInfo error: %v", err.Error())
		return nil, err
	}
	return &robot.User{
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

func (f *Framework) GetFriends(isRefresh bool) ([]*robot.User, error) {
	dataType := 1
	if isRefresh {
		dataType = 2
	}

	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0005",
		"data": map[string]interface{}{
			"type": dataType,
		},
	}

	var dataResp FriendsListResp
	if err := NewRequest().Post(apiUrl).SetBody(payload).SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[Dean] GetFriends error: %v", err.Error())
		return nil, err
	}
	var friendsInfoList []*robot.User
	for _, res := range dataResp.Result {
		friendsInfoList = append(friendsInfoList, &robot.User{
			WxId:                    res.Wxid,
			WxNum:                   res.WxNum,
			Nick:                    res.Nick,
			Remark:                  res.Remark,
			NickBrief:               res.NickBrief,
			NickWhole:               res.NickWhole,
			RemarkBrief:             res.RemarkBrief,
			RemarkWhole:             res.RemarkWhole,
			EnBrief:                 res.EnBrief,
			EnWhole:                 res.EnWhole,
			V3:                      res.V3,
			Sign:                    res.Sign,
			Country:                 res.Country,
			Province:                res.Province,
			City:                    res.City,
			MomentsBackgroundImgUrl: res.MomentsBackgroudImgUrl,
			AvatarMinUrl:            res.AvatarMinUrl,
			AvatarMaxUrl:            res.AvatarMaxUrl,
			Sex:                     res.Sex,
			MemberNum:               res.MemberNum,
		})
	}

	// 过滤系统用户
	var SystemUserWxId = map[string]struct{}{"medianote": {}, "newsapp": {}, "fmessage": {}, "floatbottle": {}}
	var filteredFriendInfo []*robot.User
	for i := range friendsInfoList {
		if _, ok := SystemUserWxId[friendsInfoList[i].WxId]; !ok {
			filteredFriendInfo = append(filteredFriendInfo, friendsInfoList[i])
		}
	}
	return filteredFriendInfo, nil
}

func (f *Framework) GetGroups(isRefresh bool) ([]*robot.User, error) {
	dataType := 1
	if isRefresh {
		dataType = 2
	}

	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0006",
		"data": map[string]interface{}{
			"type": dataType,
		},
	}

	var dataResp GroupListResp
	if err := NewRequest().Post(apiUrl).SetBody(payload).SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[Dean] GetGroups error: %v", err.Error())
		return nil, err
	}
	var groupInfoList []*robot.User
	for _, res := range dataResp.Result {
		groupInfoList = append(groupInfoList, &robot.User{
			WxId:         res.Wxid,
			WxNum:        res.WxNum,
			Nick:         res.Nick,
			Remark:       res.Remark,
			NickBrief:    res.NickBrief,
			NickWhole:    res.NickWhole,
			RemarkBrief:  res.RemarkBrief,
			RemarkWhole:  res.RemarkWhole,
			EnBrief:      res.EnBrief,
			EnWhole:      res.EnWhole,
			MemberNum:    res.MemberNum,
			AvatarMinUrl: res.AvatarMinUrl,
			AvatarMaxUrl: res.AvatarMaxUrl,
		})
	}
	return groupInfoList, nil
}

func (f *Framework) GetGroupMembers(groupWxId string, isRefresh bool) ([]*robot.User, error) {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0008",
		"data": map[string]interface{}{
			"wxid": groupWxId,
		},
	}

	var dataResp GroupMemberListResp
	if err := NewRequest().Post(apiUrl).SetBody(payload).SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[Dean] GetGroupMembers error: %v", err.Error())
		return nil, err
	}
	var groupMemberInfoList []*robot.User
	for _, res := range dataResp.Result {
		groupMemberInfoList = append(groupMemberInfoList, &robot.User{
			WxId: res.Wxid,
			Nick: res.GroupNick,
		})
	}
	return groupMemberInfoList, nil
}

func (f *Framework) GetMPs(isRefresh bool) ([]*robot.User, error) {
	dataType := 1
	if isRefresh {
		dataType = 2
	}

	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0007",
		"data": map[string]interface{}{
			"type": dataType,
		},
	}

	var dataResp SubscriptionListResp
	if err := NewRequest().Post(apiUrl).SetBody(payload).SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[Dean] GetMPs error: %v", err.Error())
		return nil, err
	}
	var subscriptionInfoList []*robot.User
	for _, res := range dataResp.Result {
		subscriptionInfoList = append(subscriptionInfoList, &robot.User{
			WxId:                    res.Wxid,
			WxNum:                   res.WxNum,
			Nick:                    res.Nick,
			Remark:                  res.Remark,
			NickBrief:               res.NickBrief,
			NickWhole:               res.NickWhole,
			RemarkBrief:             res.RemarkBrief,
			RemarkWhole:             res.RemarkWhole,
			EnBrief:                 res.EnBrief,
			EnWhole:                 res.EnWhole,
			V3:                      res.V3,
			Sign:                    res.Sign,
			Country:                 res.Country,
			Province:                res.Province,
			City:                    res.City,
			MomentsBackgroundImgUrl: res.MomentsBackgroudImgUrl,
			AvatarMinUrl:            res.AvatarMinUrl,
			AvatarMaxUrl:            res.AvatarMaxUrl,
		})
	}
	return subscriptionInfoList, nil
}
