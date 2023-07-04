package redisrepo

import (
	"encoding/json"
	"gochatapp/model"
	"log"

	"github.com/go-redis/redis/v8"
)

type Document struct {
	ID      string `json:"_id"`
	Payload []byte `json:"payload"`
	Total   int64  `json:"total"`
}

const (
	resLengthThreshold = 1
	halfValue =2
)

//res 인터페이스를 역직렬화하여 []Document로 변환 하는 함수
func Deserialize(res interface{}) []Document {
	//TODO: switch, interface 타입 변환 예시
	switch v := res.(type) {
	case []interface{}:
		if len(v) > resLengthThreshold {
			total := len(v) - 1
			//TODO: capacity 설정 예시
			docs := make([]Document, 0, total/halfValue)

			for i := 1; i <= total; i = i + 2 {
				arrOfValues := v[i+1].([]interface{})
				value := arrOfValues[len(arrOfValues)-1].(string)

				doc := Document{
					ID:      v[i].(string),
					Payload: []byte(value),
					Total:   v[0].(int64),
				}
				docs = append(docs, doc)
			}
			return docs
		}
	default:
		log.Println("Check response type. $T", res)
		return nil

	}
	return nil
}

func DeserializeChat(docs []Document) []model.Chat {
	chats := []model.Chat{}
	for _, doc := range docs {
		var c model.Chat
		json.Unmarshal(doc.Payload, &c)

		c.ID = doc.ID
		chats = append(chats, c)
	}
	return chats
}

func DeserializeContactList(contacts []redis.Z) []model.ContactList {
	contactList := make([]model.ContactList, 0, len(contacts))

	for _, contact := range contacts {
		contactList = append(contactList, model.ContactList{
			Username: contact.Member.(string),
			LastActivity: int64(contact.Score),
		})
	}
	return contactList
}
