package jpush

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"testing"
)

func setJPushCfg() {
	global.GVA_CONFIG.JPush.Key = "9fa824898bb6dc18bbc759cb"
	global.GVA_CONFIG.JPush.Secret = "75ca35675c466c47f1d44342"
}

func TestJPush(t *testing.T) {
	setJPushCfg()
	//err := JPushSingleNotify(jpush.ANDROID, "18071adc021b0f0d32b", "title", "content")
	//err := JPushSingleNotify(jpush.IOS, "1a1018970ba04bd9aa6", "Low Battery Alarm", fmt.Sprintf("Your eBike \"%s\" is less than 20 percent charged", "STSX2450002"))
	err := utils.JPushSingleNotify(utils.jpush.IOS, "1a1018970ba04bd9aa6", "Position Alarm", fmt.Sprintf("Your eBike \"%s\" was moved while it was parked. Please pay attention to the positioning of your eBike.", "STSX2450002"))
	fmt.Println(err)
}
