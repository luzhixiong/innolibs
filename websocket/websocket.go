package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Websocket struct {
	ReadBufferSize  int  `mapstructure:"read-buffer-size" json:"read-buffer-size" yaml:"read-buffer-size"`    // io读大小
	WriteBufferSize int  `mapstructure:"write-buffer-size" json:"write-buffer-size" yaml:"write-buffer-size"` // io写大小
	MultiLogin      bool `mapstructure:"multi-login" json:"multi-login" yaml:"multi-login"`                   // 多点登陆
}

// websocket 连接建立
func WebsocketStart(socketConfig []byte, sockFun ConnectionFunc, checkFun func(r *http.Request) bool) (ws *SocketServer) {
	var thisConfig Websocket
	err := json.Unmarshal(socketConfig, &thisConfig)
	if err != nil {
		fmt.Println(err)
	}

	ws = New(Config{
		ReadBufferSize:  thisConfig.ReadBufferSize,
		WriteBufferSize: thisConfig.WriteBufferSize,
		CheckOrigin:     checkFun,
	})
	ws.OnConnection(sockFun)

	return ws
}
