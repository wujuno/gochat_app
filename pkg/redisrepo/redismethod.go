package redisrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"gochatapp/model"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// 새로운 유저를 등록 후, Set 자료구조에 추가합니다.
// 등록과정이나 추가과정에 오류가 있을 시 에러를 반환합니다.
func RegisterNewUser(username, password string) error {
	err := redisClient.Set(context.Background(), username, password, 0).Err()
	if err != nil {
		log.Println("error while adding new user", err)
		return err
	}

	err = redisClient.SAdd(context.Background(), userSetKey(), username).Err()
	if err != nil {
		log.Println("error while adding in set", err)

		redisClient.Del(context.Background(),username)
		return err
	}

	return nil
}

func IsUserExist(username string) bool {
	return redisClient.SIsMember(context.Background(), userSetKey(), username).Val()
}

// 사용자이름과 비밀번호를 비교합니다.
func IsUserAuthentic(username, password string) error {
	p := redisClient.Get(context.Background(), username).Val()

	if !strings.EqualFold(p, password) {
		return fmt.Errorf("invalid username or password")
	}

	return nil
}

func UpdateContactList(username, contact string) error {
	zs := &redis.Z{Score: float64(time.Now().Unix()), Member: contact}

	err := redisClient.ZAdd(context.Background(),contactListZkey(username),zs).Err()

	if err != nil {
		log.Println("error while updating contact list. username:", username, "contact:", contact, err)
		return err
	}

	return nil
}


// 채팅 데이터를 JSON형식으로 직렬화하여 Redis에 저장합니다.
// 저장 후 연락처 목록을 업데이트 합니다.
// chat key, error를 반환합니다.
func CreateChat(c *model.Chat) (string, error) {
	chatKey := chatKey()
	fmt.Println("chat key:", chatKey)

	by, _ := json.Marshal(c)

	// JSON.SET chat#1661360942123 $ '{"from": "sun", "to":"earth","message":"good morning!"}'
	res, err := redisClient.Do(
		context.Background(),
		"JSON.SET",
		chatKey,
		"$",
		string(by),
	).Result()
	
	if err != nil {
		log.Println("error while setting chat json", err)
		return "", err
	}

	log.Println("chat successfully set", res)

	err = UpdateContactList(c.From, c.To)
	if err != nil {
		log.Println("error while updating contact list of", c.From)
	}

	err = UpdateContactList(c.To, c.From)
	if err != nil {
		log.Println("error while updating contact list of", c.To)
	}

	return chatKey, nil
}

// 전문 검색 인덱스를 생성합니다.
func CreateFetchChatBetweenIndex() {
	res, err := redisClient.Do(
		context.Background(),
		"FT.CREATE",
		chatIndex(),
		"ON", "JSON",
		"PREFIX", "1", "chat#",
		"SCHEMA", "$.from", "AS", "from", "TAG",
		"$.to", "AS", "to", "TAG",
		"$.timestamp", "AS", "timestamp", "NUMERIC", "SORTABLE", 
	).Result()

	fmt.Println(res, err)
}

// 두 사용자 간의 특정 시간 범위 내에 있는 채팅을 검색하는 함수
func FetchChatBetween(username1, username2, fromTS, toTS string) ([]model.Chat, error) {
	
	query := fmt.Sprintf("@from:{%s|%s} @to:{%s|%s} @timestamp:[%s %s]",
	username1, username2, username1, username2, fromTS, toTS)

	res, err := redisClient.Do(
		context.Background(),
		"FT.SEARCH",
		chatIndex(),
		query,
		"SORTBY", "timestamp", "DESC",
	).Result()

	if err != nil {
		return nil, err
	}

	data := Deserialize(res)

	chats := DeserializeChat(data)
	return chats, nil
}

// 특정 사용자의 연락처 목록을 가져오는 기능을 수행하는 함수
func FetchContactList(username string) ([]model.ContactList, error) {
	zRangeArg := redis.ZRangeArgs{
		Key: contactListZkey(username),
		Start: 0,
		Stop: -1,
		Rev: true,
	}

	res, err := redisClient.ZRangeArgsWithScores(context.Background(), zRangeArg).Result()

	if err != nil {
		log.Println("error while fetching contact list. username:", username, err)
		return nil, err
	}

	contactList := DeserializeContactList(res)

	return contactList, nil
}
