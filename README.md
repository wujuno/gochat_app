## 프로젝트 구조

```
📦gochatapp
┣ 📂client
┣ 📂model
┃ ┗ 📜chat.go
┣ 📂pkg
┃ ┣ 📂httpserver
┃ ┃ ┣ 📜chathandler.go
┃ ┃ ┗ 📜httpserver.go
┃ ┣ 📂redisrepo
┃ ┃ ┣ 📜client.go
┃ ┃ ┣ 📜deserialize.go
┃ ┃ ┣ 📜key.go
┃ ┃ ┗ 📜redismethod.go
┃ ┗ 📂ws
┃   ┗ 📜websocket.go
┣ 📜go.mod
┣ 📜go.sum
┗ 📜main.go
```
