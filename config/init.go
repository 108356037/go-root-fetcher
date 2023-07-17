package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
)

var (
	ETH_URL   string
	ETH_WSS   string
	REDIS_URL string

	TORN_01_ETH  common.Address
	TORN_1_ETH   common.Address
	TORN_10_ETH  common.Address
	TORN_100_ETH common.Address

	TORN_100_100_DAI common.Address
	TORN_10_000_DAI2 common.Address

	UPDATE_INTERVAL int
	LOG_LEVEL       int
)

// loadFromDotenvFile is helper func that loads env from .env in project root.
// This is handy when doing local test, however in production
// we will pass env through K8s configuration.
func loadFromDotenvFile() {
	var (
		_, b, _, _       = runtime.Caller(0)
		configModulePath = filepath.Dir(b)
	)
	envfilePath := filepath.Join(configModulePath, "..", ".env")
	err := godotenv.Load(envfilePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("No .env file at project root")
		} else {
			log.Fatalf("Error in loading .env file: %s", err.Error())
		}
	}
}

func mustGetEnv(env string) string {
	result, exist := os.LookupEnv(env)
	if !exist {
		log.Fatalf("Environment variable %s must be set", env)
	}
	return result
}

func init() {
	loadFromDotenvFile()

	ETH_URL = mustGetEnv("ETH_URL")
	ETH_WSS = mustGetEnv("ETH_WSS")
	REDIS_URL = mustGetEnv("REDIS_URL")

	TORN_01_ETH = convert2CommonAddr(mustGetEnv("TORN_01_ETH"))
	TORN_1_ETH = convert2CommonAddr(mustGetEnv("TORN_1_ETH"))
	TORN_10_ETH = convert2CommonAddr(mustGetEnv("TORN_10_ETH"))
	TORN_100_ETH = convert2CommonAddr(mustGetEnv("TORN_100_ETH"))

	TORN_100_100_DAI = convert2CommonAddr(mustGetEnv("TORN_100_100_DAI"))
	TORN_10_000_DAI2 = convert2CommonAddr(mustGetEnv("TORN_10_000_DAI2"))

	UPDATE_INTERVAL, _ = strconv.Atoi(mustGetEnv("UPDATE_INTERVAL"))
	LOG_LEVEL, _ = strconv.Atoi(mustGetEnv("LOG_LEVEL"))
}

func convert2CommonAddr(hexString string) common.Address {
	ensureType := func(hexAddress string) bool {
		if !common.IsHexAddress(hexAddress) {
			return false
		}
		return true
	}

	if !(ensureType(hexString)) {
		log.Fatal(fmt.Errorf("Cannot convert to common address: %s", hexString))
	}

	return common.HexToAddress(hexString)
}
