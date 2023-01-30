package redis

import (
	"testing"
	"time"
)

func TestRedisSingleton(t *testing.T) {
	instance := GetRedisInstance()
	instance.Connect("localhost", "", 6379)
	rdb := instance.Client
	ctx := instance.Ctx
	key := "test-key"

	// rdb.Set(ctx, key, "val")

	rdb.HSet(ctx, key, "status", "not-verified")
	rdb.HSet(ctx, key, "chance-left", "2")
	rdb.Expire(ctx, key, 10*time.Minute)

	state := rdb.HGet(ctx, key, "status").Val()
	if state != "not-verified" {
		t.Errorf("expect state to be 'not-verified', got: %s", state)
	}

	chanceLeft := rdb.HGet(ctx, key, "chance-left").Val()
	if chanceLeft != "2" {
		t.Errorf("expect state to be '2', got: %s", chanceLeft)
	}
}
