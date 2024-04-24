package websocket

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/iris-contrib/go.uuid"
)

const (
	// DefaultWebsocketWriteTimeout 0, no timeout
	DefaultWebsocketWriteTimeout = 0
	// DefaultWebsocketReadTimeout 0, no timeout
	DefaultWebsocketReadTimeout = 0
	// DefaultWebsocketPongTimeout 60 * time.Second
	DefaultWebsocketPongTimeout = 60 * time.Second
	// DefaultWebsocketPingPeriod (DefaultPongTimeout * 9) / 10
	DefaultWebsocketPingPeriod = (DefaultWebsocketPongTimeout * 9) / 10
	// DefaultWebsocketMaxMessageSize 1024
	DefaultWebsocketMaxMessageSize = 1024
	// DefaultWebsocketReadBufferSize 4096
	DefaultWebsocketReadBufferSize = 4096
	// DefaultWebsocketWriterBufferSize 4096
	DefaultWebsocketWriterBufferSize = 4096
	// DefaultEvtMessageKey is the default prefix of the underline websocket events
	// that are being established under the hoods.
	//
	// Defaults to "gin-websocket-message:".
	// Last character of the prefix should be ':'.
	DefaultEvtMessageKey = "message:"
)

var (
	// DefaultIDGenerator returns a random unique for a new connection.
	// Used when config.IDGenerator is nil.
	DefaultIDGenerator = func(*gin.Context) string {
		id, err := uuid.NewV4()
		if err != nil {
			return randomString(64)
		}
		return id.String()
	}
)

// Config the websocket server configuration
// all of these are optional.
type Config struct {
	// IDGenerator used to create (and later on, set)
	// an ID for each incoming websocket connections (clients).
	// The request is an input parameter which you can use to generate the ID (from headers for example).
	// If empty then the ID is generated by DefaultIDGenerator: randomString(64)
	IDGenerator func(ctx *gin.Context) string
	// EvtMessagePrefix is the prefix of the underline websocket events that are being established under the hoods.
	// This prefix is visible only to the javascript side (code) and it has nothing to do
	// with the message that the end-user receives.
	// Do not change it unless it is absolutely necessary.
	//
	// If empty then defaults to []byte("iris-websocket-message:").
	EvtMessagePrefix []byte
	// Error is the function that will be fired if any client couldn't upgrade the HTTP connection
	// to a websocket connection, a handshake error.
	Error func(w http.ResponseWriter, r *http.Request, status int, reason error)
	// CheckOrigin a function that is called right before the handshake,
	// if returns false then that client is not allowed to connect with the websocket server.
	CheckOrigin func(r *http.Request) bool
	// HandshakeTimeout specifies the duration for the handshake to complete.
	HandshakeTimeout time.Duration
	// WriteTimeout time allowed to write a message to the connection.
	// 0 means no timeout.
	// Default value is 0
	WriteTimeout time.Duration
	// ReadTimeout time allowed to read a message from the connection.
	// 0 means no timeout.
	// Default value is 0
	ReadTimeout time.Duration
	// PongTimeout allowed to read the next pong message from the connection.
	// Default value is 60 * time.Second
	PongTimeout time.Duration
	// PingPeriod send ping messages to the connection within this period. Must be less than PongTimeout.
	// Default value is 60 *time.Second
	PingPeriod time.Duration
	// MaxMessageSize max message size allowed from connection.
	// Default value is 1024
	MaxMessageSize int64
	// BinaryMessages set it to true in order to denotes binary data messages instead of utf-8 text
	// compatible if you wanna use the Connection's EmitMessage to send a custom binary data to the client, like a native server-client communication.
	// Default value is false
	BinaryMessages bool
	// ReadBufferSize is the buffer size for the connection reader.
	// Default value is 4096
	ReadBufferSize int
	// WriteBufferSize is the buffer size for the connection writer.
	// Default value is 4096
	WriteBufferSize int
	// EnableCompression specify if the server should attempt to negotiate per
	// message compression (RFC 7692). Setting this value to true does not
	// guarantee that compression will be supported. Currently only "no context
	// takeover" modes are supported.
	//
	// Defaults to false and it should be remain as it is, unless special requirements.
	EnableCompression bool

	// Subprotocols specifies the server's supported protocols in order of
	// preference. If this field is set, then the Upgrade method negotiates a
	// subprotocol by selecting the first match in this list with a protocol
	// requested by the client.
	Subprotocols []string
}

// Validate validates the configuration
func (c Config) Validate() Config {
	// 0 means no timeout.
	if c.WriteTimeout < 0 {
		c.WriteTimeout = DefaultWebsocketWriteTimeout
	}

	if c.ReadTimeout < 0 {
		c.ReadTimeout = DefaultWebsocketReadTimeout
	}

	if c.PongTimeout < 0 {
		c.PongTimeout = DefaultWebsocketPongTimeout
	}

	if c.PingPeriod <= 0 {
		c.PingPeriod = DefaultWebsocketPingPeriod
	}

	if c.MaxMessageSize <= 0 {
		c.MaxMessageSize = DefaultWebsocketMaxMessageSize
	}

	if c.ReadBufferSize <= 0 {
		c.ReadBufferSize = DefaultWebsocketReadBufferSize
	}

	if c.WriteBufferSize <= 0 {
		c.WriteBufferSize = DefaultWebsocketWriterBufferSize
	}

	if c.Error == nil {
		c.Error = func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			//empty
		}
	}

	if c.CheckOrigin == nil {
		c.CheckOrigin = func(r *http.Request) bool {
			// allow all connections by default
			return true
		}
	}

	if len(c.EvtMessagePrefix) == 0 {
		c.EvtMessagePrefix = []byte(DefaultEvtMessageKey)
	}

	if c.IDGenerator == nil {
		c.IDGenerator = DefaultIDGenerator
	}

	return c
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// random takes a parameter (int) and returns random slice of byte
// ex: var randomstrbytes []byte; randomstrbytes = utils.Random(32)
func random(n int) []byte {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

// randomString accepts a number(10 for example) and returns a random string using simple but fairly safe random algorithm
func randomString(n int) string {
	return string(random(n))
}
