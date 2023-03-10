package websocket

import (
	"encoding/json"
	"github.com/sunmfei/mfus/common/MFei"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// ClientManager is a websocket manager
type ClientManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

// Client is a websocket client
type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan []byte
}

// Message is return msg
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

// Manager define a ws server manager
var Manager = ClientManager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[*Client]bool),
}

// Start is  项目运行前, 协程开启start -> go Manager.Start()
// 启动websocket服务
func (manager *ClientManager) Start() {
	for {
		MFei.LOGGER.Info("<---管道通信--->")

		select {
		case conn := <-Manager.Register:
			MFei.LOGGER.Info("新用户加入:%v", conn.ID)
			Manager.Clients[conn] = true
			jsonMessage, _ := json.Marshal(&Message{Content: "Successful connection to socket service"})
			Manager.Send(jsonMessage, conn)
		case conn := <-Manager.Unregister:
			MFei.LOGGER.Info("用户离开:%v", conn.ID)
			if _, ok := Manager.Clients[conn]; ok {
				close(conn.Send)
				delete(Manager.Clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Content: "A socket has disconnected"})
				Manager.Send(jsonMessage, conn)
			}

		case message := <-Manager.Broadcast:

			jsonMessage, _ := json.Marshal(&Message{Content: string(message)})
			for conn := range Manager.Clients {
				select {
				case conn.Send <- jsonMessage:
				default:
					close(conn.Send)
					delete(Manager.Clients, conn)
				}
			}
		}
	}
}

// Send is to send ws message to ws client
// 向连接websocket的管道chan写入数据
func (manager *ClientManager) Send(message []byte, ignore *Client) {

	for conn := range manager.Clients {
		// if conn != ignore { //向除了自己的socket 用户发送
		conn.Send <- message
		// }
	}
}

// 读取在websocket管道中的数据
func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()

	for {

		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		MFei.LOGGER.Info("读取到客户端的信息:%s", string(message))
		Manager.Broadcast <- message
	}
}

// 通过websocket协议向连接到ws的客户端发送数据
func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				err := c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					return
				}
				return
			}
			MFei.LOGGER.Info("发送到到客户端的信息:%s", string(message))

			err := c.Socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}

}

// TestHandler socket 连接 中间件 作用:升级协议,用户验证,自定义信息等
// ws链接交互的中间件, 用于协议升级, 用户信息验证等
func TestHandler(c *gin.Context, userID string) {
	c.Request.Header.Add("proxy_http_version", "1.1")
	c.Request.Header.Add("Upgrade", "websocket")
	c.Request.Header.Add("Connection", "Upgrade")
	c.Request.Header.Add("proxy_read_timeout", "600s")
	c.Request.Header.Add("Sec-Websocket-Version", "13")
	c.Request.Method = "GET"
	conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)

	// 升级协议
	/*conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		fmt.Println("升级协议", "ua:", r.Header["User-Agent"], "referer:", r.Header["Referer"])

		return true
	}}).Upgrade(c.Writer, c.Request, nil)*/

	if err != nil {
		http.NotFound(c.Writer, c.Request)
		MFei.LOGGER.Error("sw连接失败!", err)
		return
	}

	//可以添加用户信息验证
	client := &Client{
		ID:     userID,
		Socket: conn,
		Send:   make(chan []byte),
	}
	Manager.Register <- client
	go client.Read()
	go client.Write()
}
