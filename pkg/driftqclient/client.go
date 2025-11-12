package driftqclient

import (
	"errors"
	"fmt"
)

type Client struct {
	addr string
}

func New(addr string) *Client { return &Client{addr: addr} }

type Topic struct {
	Name       string
	Partitions int
	Compacted  bool
}

func (c *Client) Health() (bool, error) {
	_ = c.addr
	return true, nil
}

func (c *Client) TopicCreate(name string, partitions int, compacted bool) error {
	if name == "" || partitions <= 0 { return errors.New("invalid topic params") }
	fmt.Printf("(stub) create topic name=%s partitions=%d compacted=%v\n", name, partitions, compacted)
	return nil
}

func (c *Client) TopicList() ([]Topic, error) {
	return []Topic{
		{Name: "orders", Partitions: 3, Compacted: false},
		{Name: "payments", Partitions: 1, Compacted: true},
	}, nil
}

func (c *Client) Produce(topic, key string, value []byte) error {
	if topic == "" { return errors.New("missing topic") }
	fmt.Printf("(stub) produced to %s key=%q value_bytes=%d\n", topic, key, len(value))
	return nil
}

func (c *Client) Consume(topic, group, from string, handler func(key string, value []byte) error) error {
	if topic == "" || group == "" { return errors.New("missing topic/group") }
	_ = handler("k1", []byte("hello"))
	_ = handler("k2", []byte("world"))
	return nil
}

func (c *Client) Lag(topic, group string) (int64, error) {
	return 0, nil
}
