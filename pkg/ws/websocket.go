package ws

import (
	"encoding/json"
	"fmt"
	"gochatapp/model"
	"gochatapp/pkg/redisrepo"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

/*
Client 구조체는 사용자의 WebSocket 연결 주소와 사용자명을 저장합니다.
Message 구조체는 메시지의 유형, 사용자명 및 채팅 모델을 저장합니다.
이는 ReactJS 애플리케이션으로부터 수신된 메시지입니다. 누군가가 누군가에게 메시지를 보낼 때와 같은 경우입니다.
2개의 메시지 유형이 있을 수 있습니다. 1. bootup과 2. 그 외의 모든 메시지. 클라이언트가 WebSocket 서버에 처음 연결될 때 클라이언트는 bootup 유형과 사용자명을 함께 보내어 연결을 사용자명과 매핑합니다.
clients: 모든 애플리케이션에 연결된 클라이언트를 저장하는 맵입니다.
broadcast: 클라이언트로부터 메시지를 수신하면 해당 메시지를 연결된 모든 클라이언트에게 브로드캐스트하는 채널입니다.
upgrader: WebSocket 업그레이드 구성입니다.
serveWs: TCP 연결로의 연결 업그레이드를 수행하는 HTTP 핸들러입니다. 성공적인 업그레이드 후 새로운 클라이언트를 clients 맵에 추가합니다. 그런 다음 receiver 함수를 사용하여 연결을 수신 대기합니다. receiver 함수에서 빠져나오면 연결이 끊어지고 클라이언트가 더 이상 활동하지 않는 것을 의미합니다. delete 함수는 비활성 클라이언트를 clients 맵에서 제거합니다.
receiver: 클라이언트의 연결을 계속 수신 대기하는 무한 for 루프입니다. 메시지를 수신하면 부팅 메시지인지 확인합니다. 부팅 메시지인 경우 클라이언트의 연결을 사용자명과 매핑합니다.
일반 메시지인 경우 Redis에 채팅을 생성하고 브로드캐스트합니다. 이렇게 하면 수신 대상 사용자가 연결되어 있다면 메시지를 즉시 수신할 수 있습니다.
broadcaster: 무한 for 루프에서 브로드캐스트 채널을 계속 수신 대기합니다. 여기서 broadcaster는 메시지와 관련된 클라이언트를 필터링한 다음 해당 메시지를 보냅니다. broadcaster가 모든 클라이언트에게 이 메시지를 보내면 clients는 그룹으로 동작합니다. 다음 부분에서는 그룹 채팅을 구현할 것입니다.
StartWebsocketServer: 먼저 Redis 인스턴스를 초기화합니다. 별도의 go 루틴에서 broadcaster를 시작한 다음 라우트를 설정하고 WebSocket 서버를 시작합니다.
*/

type Client struct {
	Conn *websocket.Conn
	Username string
}
 
type Message struct {
	Type string     `json:"type"`
	User string     `json:"user,omitempty"`
	Chat model.Chat `json:"chat,omitempty"`
}

//TODO: make map 예시

var clients = make(map[*Client]bool)
var broadcast = make(chan *model.Chat)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	
	CheckOrigin: func(r *http.Request) bool { return true },
}

func receiver(client *Client) {
	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		m := &Message{}

		err = json.Unmarshal(p, m)
		if err != nil {
			log.Println("error while unmarshaling chat", err)
			continue
		}

		fmt.Println("host", client.Conn.RemoteAddr())
		if m.Type == "bootup" {
			client.Username = m.User
			fmt.Println("client successfully mapped", &client, client, client.Username)
		} else {
			fmt.Println("received message", m.Type, m.Chat)
			c := m.Chat
			c.Timestamp = time.Now().Unix()

			id, err := redisrepo.CreateChat(&c)
			if err != nil {
				log.Println("errors while saving chat in redis", err)
				return
			}

			c.ID = id
			broadcast <- &c
		}
	}
}

func broadcaster() {
	for {
		message := <-broadcast
		// send to every client that is currently connected
		fmt.Println("new message", message)

		for client := range clients {
			// send message only to involved users
			fmt.Println(
				"username:", client.Username,
				"from:", message.From,
				"to:", message.To,
			)

			if client.Username == message.From || client.Username == message.To {
				err := client.Conn.WriteJSON(message)
				if err != nil {
					log.Println("Websocket error:", err)
					client.Conn.Close()
					delete(clients, client)
				}
			}
		}
	}
}

// define our WebSocket endpoint
func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host, r.URL.Query())

	// upgrade this connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := &Client{Conn: ws}
	//register client
	clients[client] = true
	fmt.Println("clients", len(clients), clients, ws.RemoteAddr())

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	receiver(client)

	fmt.Println("exiting", ws.RemoteAddr().String())
	delete(clients, client)
}

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Server")
	})
	// map our `/ws` endpoint to the `serveWs` function
	http.HandleFunc("/ws", serveWs)
}

func StartWebsocketServer() {
	redisClient := redisrepo.InitializeRedis()
	defer redisClient.Close()

	go broadcaster()
	setupRoutes()
	http.ListenAndServe(":8081", nil)
}