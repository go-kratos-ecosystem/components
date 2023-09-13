package redis

import (
	"context"
	"testing"
	"time"

	"github.com/go-packagist/go-kratos-components/contracts/cache"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var (
	ctx  = context.Background()
	addr = "127.0.0.1:6379"
)

type person struct {
	Name string
	Age  int
	Sex  string
}

func createStore() cache.Store {
	return New(
		Prefix("test"),
		Redis(
			redis.NewClient(&redis.Options{
				Addr: addr,
			}),
		),
	)
}

func TestRedis_Put(t *testing.T) {
	c := createStore()
	defer c.Flush(ctx)
	c.Flush(ctx)

	// put
	assert.NoError(t, c.Put(ctx, "int", 1, time.Second*10))
	assert.NoError(t, c.Put(ctx, "int32", int32(1), time.Second*10))
	assert.NoError(t, c.Put(ctx, "string", "test", time.Second*10))
	assert.NoError(t, c.Put(ctx, "bool", true, time.Second*10))
	var now = time.Now()
	assert.NoError(t, c.Put(ctx, "time", time.Now(), time.Second*10))
	assert.NoError(t, c.Put(ctx, "person", &person{
		Name: "test",
		Age:  18,
		Sex:  "male",
	}, time.Second*10))
	assert.NoError(t, c.Put(ctx, "person", person{
		Name: "test",
		Age:  18,
		Sex:  "male",
	}, time.Second*10))
	assert.NoError(t, c.Put(ctx, "nil", nil, time.Second*10))

	// check
	var (
		intVal     int
		int32Val   int32
		stringVal  string
		boolVal    bool
		timeVal    time.Time
		personVal  *person
		personVal2 person
		nilVal     interface{}
	)
	assert.NoError(t, c.Get(ctx, "int", &intVal))
	assert.NoError(t, c.Get(ctx, "int32", &int32Val))
	assert.NoError(t, c.Get(ctx, "string", &stringVal))
	assert.NoError(t, c.Get(ctx, "bool", &boolVal))
	assert.NoError(t, c.Get(ctx, "time", &timeVal))
	assert.NoError(t, c.Get(ctx, "person", &personVal))
	assert.NoError(t, c.Get(ctx, "person", &personVal2))
	assert.NoError(t, c.Get(ctx, "nil", &nilVal))

	assert.Equal(t, 1, intVal)
	assert.Equal(t, int32(1), int32Val)
	assert.Equal(t, "test", stringVal)
	assert.Equal(t, true, boolVal)
	assert.Equal(t, now.Format("2006-01-02 15:04:05"), timeVal.Format("2006-01-02 15:04:05"))
	assert.Equal(t, &person{
		Name: "test",
		Age:  18,
		Sex:  "male",
	}, personVal)
	assert.Equal(t, person{
		Name: "test",
		Age:  18,
		Sex:  "male",
	}, personVal2)
	assert.Nil(t, nilVal)
}

func TestRedis_Get(t *testing.T) {
	c := createStore()
	defer c.Flush(ctx)
	c.Flush(ctx)

	assert.NoError(t, c.Put(ctx, "int", 1, time.Second*10))
	var intVal int
	assert.NoError(t, c.Get(ctx, "int", &intVal))
	assert.Equal(t, 1, intVal)

	// err
	var notExistVal int
	assert.Error(t, c.Get(ctx, "not_exist", &notExistVal))
}

func TestRedis_GetPrefix(t *testing.T) {
	c := createStore()

	assert.Equal(t, "test:", c.GetPrefix())
}

func TestRedis_Forget(t *testing.T) {
	c := createStore()
	defer c.Flush(ctx)
	c.Flush(ctx)

	assert.NoError(t, c.Put(ctx, "int", 1, time.Second*10))
	var intVal int
	assert.NoError(t, c.Get(ctx, "int", &intVal))
	assert.Equal(t, 1, intVal)

	assert.NoError(t, c.Forget(ctx, "int"))

	var notExistVal int
	assert.Error(t, c.Get(ctx, "int", &notExistVal))
}
