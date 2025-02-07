package jpush

import (
	"github.com/gin-gonic/gin"
	"time"
)

const (
	largeIcon = "https://static-us.innolabs.work/fly-logo.png"
	smallIcon = "https://static-us.innolabs.work/fly-logo.png\""
)

// 极光推送推送给单个用户
func JPushSingleNotify(jpushKey, jpushSecret string, platform PlatformType, id, title, content string, largeUrl string, smallUrl string) (err error) {

	//if largeUrl == "" {
	//	largeUrl = largeIcon
	//}
	//if smallUrl == "" {
	//	smallUrl = smallIcon
	//}
	cid := NewCidRequest(1, "")
	cidList, err := cid.GetCidList(jpushKey, jpushSecret)
	if err != nil {
		return
	}
	// 构建推送平台
	var pf Platform
	_ = pf.Add(platform)
	// 构建接收目标
	var at Audience
	//at.All()
	at.SetID([]string{id})
	// 构建通知
	payload := NewPayLoad()
	var n Notification
	n.SetAlternateAlert("")
	if platform == IOS {
		n.SetIos(&IosNotification{Alert: gin.H{"body": content, "title": title}, Badge: "+1", InterruptionLevel: "active", MutableContent: 1, Sound: "default", ThreadId: ""})
	} else if platform == ANDROID {
		n.Alert = content
		n.SetAndroid(&AndroidNotification{Alert: content, Title: title, AlertType: 7, Intent: map[string]string{"url": ""}, Priority: 0, Category: "", BadgeAddNum: 1, LargeIcon: largeUrl, SmallIconUri: smallUrl})
	}
	// 构建消息
	//var m jpush.Message
	//m.MsgContent = content
	//m.Title = title
	// 构建消息负载
	payload.SetPlatform(&pf)
	payload.SetAudience(&at)
	payload.SetNotification(&n)
	//payload.SetMessage(&m)
	payload.SetOptions(&Options{ApnsProduction: true, TimeToLive: 86400, SendNo: int(Unixtime()), Classification: 0, AlternateSet: false})
	//payload.SetInappMessage(&jpush.InappMessage{InAppMessage: true})
	payload.Cid = cidList.CidList[0]
	data, err := payload.Bytes()
	if err != nil {
		return
	}
	c := NewJPushClient(jpushKey, jpushSecret)
	_, err = c.Push(data)
	return
}

func Unixtime() uint {
	return uint(time.Now().Unix())
}
