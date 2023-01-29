package utils

import (
	"errors"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

// check envariable and set value if exist, or get value from user input

type Env struct {
	IsProxy   bool
	TrustMode int // 1. proxy hop count, 2. Trust local, cloudflare, 3. more turst address, if not set remoteAddr

	// if select option 2
	TrustAddress []string

	// if select option 3
	ProxyHopCount int
}

func SetEnv() (Env, error) {
	if os.Getenv("PROXY_HOP_COUNT") == "" {
		return Env{}, errors.New("PROXY_HOP_COUNT is not set")
	}
	proxyHopCount, err := strconv.Atoi(os.Getenv("PROXY_HOP_COUNT"))
	if err != nil {
		return Env{}, err
	}
	return Env{
		ProxyHopCount: proxyHopCount,
	}, nil
}
