package vlw

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/imroc/req/v3"

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
			buff.WriteString(`[emoji=\u`)
			buff.WriteString(fmt.Sprintf("%04x", r) + `]`)
		case 4:
			r1, r2 := utf16.EncodeRune(r)
			buff.WriteString(`[emoji=\u`)
			buff.WriteString(strconv.FormatInt(int64(r1), 16))
			buff.WriteString(`\u`)
			buff.WriteString(strconv.FormatInt(int64(r2), 16) + `]`)
		default:
			buff.WriteString(string(r))
		}
	}
	return buff.String()
}

func (f *Framework) GetRobotInfo() (*robot.User, error) {
	info, err := f.GetObjectInfo(f.BotWxId)
	if err != nil {
		log.Errorf("[VLW] GetRobotInfo error: %v", err.Error())
		return nil, err
	}
	return info, nil
}

func (f *Framework) GetMemePictures(msg *robot.Message) string {
	// 获取图片base64
	path := msg.Content[5 : len(msg.Content)-1]
	payload := map[string]interface{}{
		"api":   "GetFileFoBase64",
		"token": f.ApiToken,
		"path":  path,
	}

	var dataResp MessageResp
	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).SetResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[VLW] GetMemePictures error: %v", err.Error())
		return ""
	}
	return dataResp.ReturnStr
}

func (f *Framework) SendText(toWxId, text string) error {
	payload := map[string]interface{}{
		"api":        "SendTextMsg",
		"token":      f.ApiToken,
		"msg":        f.msgFormat(text),
		"to_wxid":    toWxId,
		"robot_wxid": f.BotWxId,
	}

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendText error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendTextAndAt(toGroupWxId, toWxId, toWxName, text string) error {
	payload := map[string]interface{}{
		"api":         "SendGroupMsgAndAt",
		"token":       f.ApiToken,
		"msg":         f.msgFormat(text),
		"robot_wxid":  f.BotWxId,
		"group_wxid":  toGroupWxId,
		"member_wxid": toWxId,
		"member_name": toWxName,
	}

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendTextAndAt error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendImage(toWxId, path string) error {
	payload := map[string]interface{}{
		"api":        "SendImageMsg",
		"token":      f.ApiToken,
		"path":       path,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
	}

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendImage error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendShareLink(toWxId, title, desc, imageUrl, jumpUrl string) error {
	payload := map[string]interface{}{
		"api":        "SendShareLinkMsg",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
		"title":      title,
		"desc":       desc,
		"image_url":  imageUrl,
		"url":        jumpUrl,
	}

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendShareLink error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendFile(toWxId, path string) error {
	payload := map[string]interface{}{
		"api":        "SendFileMsg",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
		"path":       path,
	}

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendFile error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendVideo(toWxId, path string) error {
	payload := map[string]interface{}{
		"api":        "SendVideoMsg",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
		"path":       path,
	}

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendVideo error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendEmoji(toWxId, path string) error {
	payload := map[string]interface{}{
		"api":        "SendEmojiMsg",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
		"path":       path,
	}

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendEmoji error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendMusic(toWxId, name, author, app, jumpUrl, musicUrl, coverUrl string) error {
	payload := map[string]interface{}{
		"api":        "SendMusicLinkMsg",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
		"title":      name,
		"desc":       author,
		"url":        jumpUrl,
		"dataurl":    musicUrl,
		"thumburl":   coverUrl,
	}

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendMusic error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendMiniProgram(toWxId, ghId, title, content, imagePath, jumpPath string) error {
	log.Errorf("[VLW] SendMiniProgram not support")
	return errors.New("SendMiniProgram not support，please use SendXML")
}

func (f *Framework) SendMessageRecord(toWxId, title string, dataList []map[string]interface{}) error {
	log.Errorf("[VLW] SendMessageRecord not support")
	return errors.New("SendMessageRecord not support, please use SendMessageRecordXML")
}

func (f *Framework) SendMessageRecordXML(toWxId, xmlStr string) error {
	payload := map[string]interface{}{
		"api":        "SendMessageRecord",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
		"content":    xmlStr,
	}

	buf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(payload)
	if err := req.C().Post(f.ApiUrl).SetBodyJsonString(buf.String()).Do().Err; err != nil {
		log.Errorf("[VLW] SendMessageRecordXML error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendFavorites(toWxId, favoritesId string) error {
	payload := map[string]interface{}{
		"api":        "SendFavorites",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
		"local_id":   favoritesId,
	}

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendFavorites error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendXML(toWxId, xmlStr string) error {
	payload := map[string]interface{}{
		"api":        "SendXmlMsg",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
		"xml":        xmlStr,
	}

	buf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(payload)
	if err := req.C().Post(f.ApiUrl).SetBodyJsonString(buf.String()).Do().Err; err != nil {
		log.Errorf("[VLW] SendXML error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendBusinessCard(toWxId, targetWxId string) error {
	payload := map[string]interface{}{
		"api":        "SendCardMsg",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    toWxId,
		"content":    targetWxId,
	}

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] SendBusinessCard error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) AgreeFriendVerify(v1, v2, scene string) error {
	sceneInt, err := strconv.Atoi(scene)
	if err != nil {
		log.Errorf("[VLW] AgreeFriendVerify error: %v", err)
		return err
	}

	payload := map[string]interface{}{
		"api":        "AgreeFriendVerify",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"v1":         v1,
		"v2":         v2,
		"type":       sceneInt,
	}

	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] AgreeFriendVerify error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) InviteIntoGroup(groupWxId, wxId string, typ int) error {
	payload := make(map[string]interface{}, 5)
	switch typ {
	case 1:
		payload = map[string]interface{}{
			"api":         "InviteInGroup",
			"token":       f.ApiToken,
			"robot_wxid":  f.BotWxId,
			"group_wxid":  groupWxId,
			"friend_wxid": wxId,
		}
	case 2:
		payload = map[string]interface{}{
			"api":         "InviteInGroupByLink",
			"token":       f.ApiToken,
			"robot_wxid":  f.BotWxId,
			"group_wxid":  groupWxId,
			"friend_wxid": wxId,
		}
	}
	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[VLW] InviteIntoGroup error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) GetObjectInfo(wxId string) (*robot.User, error) {
	payload := map[string]interface{}{
		"api":        "GetInfoByWxid",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"to_wxid":    wxId,
	}

	var dataResp ObjectInfoResp
	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[VLW] GetObjectInfo error: %v", err.Error())
		return nil, err
	}
	return &robot.User{
		WxId:         dataResp.ReturnJson.Data.Wxid,
		WxNum:        dataResp.ReturnJson.Data.Account,
		Nick:         dataResp.ReturnJson.Data.Nickname,
		Remark:       dataResp.ReturnJson.Data.Remark,
		V3:           dataResp.ReturnJson.Data.V1,
		V4:           dataResp.ReturnJson.Data.V2,
		Sign:         dataResp.ReturnJson.Data.Signature,
		Country:      dataResp.ReturnJson.Data.Country,
		Province:     dataResp.ReturnJson.Data.Province,
		City:         dataResp.ReturnJson.Data.City,
		AvatarMinUrl: dataResp.ReturnJson.Data.SmallAvatar,
		AvatarMaxUrl: dataResp.ReturnJson.Data.Avatar,
		Sex:          strconv.Itoa(dataResp.ReturnJson.Data.Sex),
	}, nil
}

func (f *Framework) GetFriends(isRefresh bool) ([]*robot.User, error) {
	dataType := 0
	if isRefresh {
		dataType = 1
	}
	payload := map[string]interface{}{
		"api":        "GetFriendlist",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"is_refresh": dataType,
	}

	var dataResp FriendsListResp
	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[VLW] GetFriends error: %v", err.Error())
		return nil, err
	}
	var friendsInfoList []*robot.User
	for _, res := range dataResp.ReturnJson {
		friendsInfoList = append(friendsInfoList, &robot.User{
			WxId:         res.Wxid,
			WxNum:        res.WxNum,
			Nick:         res.Nickname,
			Remark:       res.Note,
			Country:      res.Country,
			Province:     res.Province,
			City:         res.City,
			AvatarMinUrl: res.Avatar,
			AvatarMaxUrl: res.Avatar,
			Sex:          strconv.Itoa(res.Sex),
		})
	}
	return friendsInfoList, nil
}

func (f *Framework) GetGroups(isRefresh bool) ([]*robot.User, error) {
	dataType := 0
	if isRefresh {
		dataType = 1
	}
	payload := map[string]interface{}{
		"api":        "GetGrouplist",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"is_refresh": dataType,
	}

	var dataResp GroupListResp
	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[VLW] GetGroups error: %v", err.Error())
		return nil, err
	}
	var groupInfoList []*robot.User
	for _, res := range dataResp.ReturnJson {
		groupInfoList = append(groupInfoList, &robot.User{
			WxId:         res.Wxid,
			Nick:         res.Nickname,
			MemberNum:    res.TotalMember,
			AvatarMinUrl: res.Avatar,
			AvatarMaxUrl: res.Avatar,
		})
	}
	return groupInfoList, nil
}

func (f *Framework) GetGroupMembers(groupWxId string, isRefresh bool) ([]*robot.User, error) {
	dataType := 0
	if isRefresh {
		dataType = 1
	}
	payload := map[string]interface{}{
		"api":        "GetGroupMember",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"group_wxid": groupWxId,
		"is_refresh": dataType,
	}

	var dataResp GroupMemberListResp
	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[VLW] GetGroupMembers error: %v", err.Error())
		return nil, err
	}
	var groupMemberInfoList []*robot.User
	for _, res := range dataResp.ReturnJson.MemberList {
		groupMemberInfoList = append(groupMemberInfoList, &robot.User{
			WxId:         res.Wxid,
			WxNum:        res.WxNum,
			Nick:         res.Nickname,
			Remark:       res.Remark,
			Country:      res.Country,
			Province:     res.Province,
			City:         res.City,
			AvatarMinUrl: res.Avatar,
			AvatarMaxUrl: res.Avatar,
			Sex:          strconv.Itoa(res.Sex),
		})
	}
	return groupMemberInfoList, nil
}

func (f *Framework) GetMPs(isRefresh bool) ([]*robot.User, error) {
	dataType := 0
	if isRefresh {
		dataType = 1
	}
	payload := map[string]interface{}{
		"api":        "GetSubscriptionlist",
		"token":      f.ApiToken,
		"robot_wxid": f.BotWxId,
		"is_refresh": dataType,
	}

	var dataResp SubscriptionListResp
	if err := NewRequest().Post(f.ApiUrl).SetBody(payload).SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[VLW] GetMPs error: %v", err.Error())
		return nil, err
	}
	var subscriptionInfoList []*robot.User
	for _, res := range dataResp.ReturnJson {
		subscriptionInfoList = append(subscriptionInfoList, &robot.User{
			WxId:         res.Wxid,
			Nick:         res.Nickname,
			AvatarMinUrl: res.Avatar,
			AvatarMaxUrl: res.Avatar,
		})
	}
	return subscriptionInfoList, nil
}
