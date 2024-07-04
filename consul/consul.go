package consul

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"bitbucket.org/junglee_games/getsetgo/logger"
	"github.com/hashicorp/consul/api"
)

type Agent struct {
	kv   *api.KV
	path string
}

type Config struct {
	Address string
	Name    string
	Token   string
}

// NewAgent returns an initalized Agent.
func New(path string, cfg Config) (Agent, error) {
	log.SetPrefix(cfg.Name + ": CONSUL  :  ")
	// Get a new client
	client, err := api.NewClient(&api.Config{
		Address: cfg.Address,
		Token:   cfg.Token,
	})
	if err != nil {
		return Agent{}, err
	}

	// Get a handle to the KV API
	kv := client.KV()
	return Agent{kv, path}, nil
}

type ConfigSetter interface {
	Set(map[string]string) error
}

func (a Agent) InitAndGetConfig(cs ConfigSetter) error {
	pairs, _, err := a.kv.List(a.path, nil)
	if err != nil {
		return err
	}
	ch := make(chan *api.KVPair)
	input := make(map[string]string)
	for _, kv := range pairs {
		_key := strings.ReplaceAll(kv.Key, a.path+"/", "")
		if _key != "" {
			input[_key] = strings.TrimSpace(string(kv.Value))
		}
	}
	go subscribeToChanges(a.kv, a.path, pairs, ch)
	go func() {
		for c := range ch {
			val := strings.TrimSpace(string(c.Value))
			_key := strings.ReplaceAll(c.Key, a.path+"/", "")
			if input[_key] != val {
				logger.Info(context.Background(), fmt.Sprintln("Change detected : ", c.Key, val))
				input[_key] = val
				err := cs.Set(input)
				if err != nil {
					logger.Info(context.Background(), fmt.Sprintln("Unable to map Change, invalid type : ", c.Key, string(c.Value)))
				}
			}
		}
	}()
	if err := cs.Set(input); err != nil {
		return err
	}
	return nil
}

func subscribeToChanges(kv *api.KV, path string, pairs api.KVPairs, ch chan *api.KVPair) {
	var watchers = make(map[string]uint64)
	for {
		for _, pair := range pairs {
			_key := strings.ReplaceAll(pair.Key, path+"/", "")
			if _key == "" {
				continue
			}
			currentIndex, ok := watchers[_key]
			if !ok {
				currentIndex = pair.CreateIndex
			}
			pair, meta, err := kv.Get(pair.Key, &api.QueryOptions{
				WaitIndex: currentIndex,
				WaitTime:  time.Second,
			})
			if err != nil {
				logger.Error(context.Background(), fmt.Sprint("Error read from KV", err.Error(), err))
				continue
			}
			if pair != nil {
				ch <- pair
				watchers[_key] = meta.LastIndex
			}
		}
		// Query won’t be blocked if key not found
		time.Sleep(30 * time.Second)
	}
}
