package redis

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string `mapstructure:"address"`
	Password string `mapstructure:"password"`
}

type Connector struct {
	clientPool map[int]*redis.Client
	cfg        Config
	m          sync.Mutex
}

func NewConnector(cfg Config) *Connector {
	return &Connector{
		clientPool: make(map[int]*redis.Client),
		cfg:        cfg,
		m:          sync.Mutex{},
	}
}

func (c *Connector) initClient(db int) {
	client := redis.NewClient(&redis.Options{
		Addr:     c.cfg.Addr,
		Password: c.cfg.Password,
		DB:       db,
	})

	c.clientPool[db] = client
}

func (c *Connector) GetClient(db int) *redis.Client {
	c.m.Lock()
	defer c.m.Unlock()

	if client, exist := c.clientPool[db]; exist {
		return client
	}

	c.initClient(db)

	return c.clientPool[db]
}

func (c *Connector) Close() error {
	c.m.Lock()
	defer c.m.Unlock()

	for db, client := range c.clientPool {
		if err := client.Close(); err != nil {
			return err
		}
		delete(c.clientPool, db)
	}

	return nil
}
