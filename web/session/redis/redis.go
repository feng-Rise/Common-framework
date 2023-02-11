package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"goStudy/web/session"
	"time"
)

var errSessionNotExist = errors.New("redis-session: session 不存在")

type StoreOption func(store *Store)

type Store struct {
	prefix     string
	client     redis.Cmdable
	expiration time.Duration
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	const lua = `redis.call("hset", KEYS[1], ARGV[1], ARGV[2])
		return redis.call("pexpire", KEYS[1], ARGV[3])`
	key := s.key(id)
	_, err := s.client.Eval(ctx, lua, []string{key}, "_sess_id", id, s.expiration.Milliseconds()).Result()
	if err != nil {
		return nil, err
	}
	return &Session{
		key:    key,
		id:     id,
		client: s.client,
	}, nil
}
func (s *Store) Refresh(ctx context.Context, id string) error {
	key := s.key(id)
	affected, err := s.client.Expire(ctx, key, s.expiration).Result()
	if err != nil {
		return err
	}
	if !affected {
		return errSessionNotExist
	}
	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	_, err := s.client.Del(ctx, s.key(id)).Result()
	return err
}

func (s *Store) GET(ctx context.Context, id string) (session.Session, error) {
	key := s.key(id)
	// 这里不需要考虑并发的问题，因为在你检测的当下，没有就是没有
	i, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if i < 0 {
		return nil, errors.New("redis-session: session 不存在")
	}
	return &Session{
		id:     id,
		key:    key,
		client: s.client,
	}, nil
}

func NewStore(client redis.Cmdable, opts ...StoreOption) *Store {
	res := &Store{
		client: client,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (s *Store) key(id string) string {
	res := fmt.Sprint("%s_%s", s.prefix, id)
	return res
}

type Session struct {
	key    string
	id     string
	client redis.Cmdable
}

func (s *Session) GET(ctx context.Context, key string) (string, error) {
	return s.client.HGet(ctx, s.key, key).Result()
}

func (s *Session) SET(ctx context.Context, key, value string) error {
	const lua = `
		if redis.call("exists", KEYS[1])
		then
			return redis.call("hset", KEYS[1], ARGV[1], ARGV[2])
		else
			return -1
		end
		`
	res, err := s.client.Eval(ctx, lua, []string{s.key}, key, value).Int()
	if err != nil {
		return err
	}
	if res < 0 {
		return errSessionNotExist
	}
	return nil
}

func (s *Session) ID() string {
	return s.id
}
