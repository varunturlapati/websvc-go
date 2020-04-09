package redisclient

import (
	"fmt"
	"time"
	
	"github.com/go-redis/redis"
	
	"github.com/varunturlapati/simpleWebSvc/pkg/article"
	db "github.com/varunturlapati/simpleWebSvc/pkg/db"
)

func New() (*RedisClient, error) {
	return &RedisClient{redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})}, nil
}

type RedisClient struct {
	*redis.Client
}

func (r *RedisClient) PingPong() (string, error) {
	return r.Ping().Result()
}

func (r *RedisClient) AddEntry(e *db.Entry) error {
	if e == nil {
		return fmt.Errorf("entry is nil")
	}
	if e.Id == "" {
		return fmt.Errorf("can't add an entry with a blank id")
	}
	err := r.Set(e.Id, &e.Value, 0).Err()
	if err != nil {
		return fmt.Errorf("error adding entry with Id %s. Err: %v\n", e.Id, err)
	}
	return nil
}

func (r *RedisClient) RemoveEntry(e string) error {
	if e == "" {
		return fmt.Errorf("can't remove entry with a blank id")
	}
	err := r.Del(e).Err()
	return err
}

func (r *RedisClient) ChangeEntry(e1 string, e2 *db.Entry) error {
	if e2 == nil {
		return fmt.Errorf("entry is nil")
	}
	if e1 == "" {
		return fmt.Errorf("can't find an entry with a blank id")
	}
	err := r.Set(e1, e2, 0).Err()
	return err
}

func (r *RedisClient) GetEntry(e string) (*db.Entry, error) {
	if e == "" {
		return nil, fmt.Errorf("entry id is blank")
	}
	res, err := r.Get(e).Result()
	if err != nil {
		return nil, fmt.Errorf("didn't find an entry with id %s\n", e)
	}
	var art article.Article
	err = art.UnmarshalBinary([]byte(res))
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal the result into an article")
	}
	resEntry := &db.Entry{
		Id:    e,
		Value: art,
		Ts:    time.Time{},
	}
	return resEntry, nil
}

func (r *RedisClient) GetAllEntries() ([]*db.Entry, error) {
	var res []*db.Entry
	ks, err := r.Keys(`*`).Result()
	if err != nil {
		return res, fmt.Errorf("error when fetching all keys. Err: %v\n", err)
	}
	for i, k := range ks {
		resN, errN := r.GetEntry(k)
		if errN != nil {
			fmt.Printf("err fetching item with key %s. At index %d, Err: %v\n", k, i, errN)
		} else {
			res = append(res, resN)
		}
	}
	return res, nil
}
