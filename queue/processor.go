package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/108356037/torn-root-fetcher/redis"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

func LogProcessor(processChan <-chan []byte) {
	for proc := range processChan {
		// if len(proc) == 0 {
		// 	continue
		// }
		fmt.Printf("we are in LogProcessor: %v, counter: %d\n", proc, Counter)
	}
}

func TornRootRedisProcessor(processChan <-chan []byte) {
	for proc := range processChan {
		if len(proc) != 52 {
			log.Warn("Got byte length unexpected", "data", proc, "length", len(proc))
			return
		}

		concatedValue := common.Bytes2Hex(proc)
		err := redis.RedisClient.Set(context.Background(), concatedValue, true, time.Duration(time.Minute*30)).Err()
		if err != nil {
			log.Warn("Error when writing to redis", "data", concatedValue)
		}
		log.Debug("Successfully write to redis", "data", proc)

		if log.GetLevel() == log.DebugLevel {
			res := redis.RedisClient.Get(context.Background(), concatedValue)
			if res != nil {
				exist, _ := res.Bool()
				log.Debug("Redis get", "key", concatedValue, "value", exist)
			}
		}
	}
}
