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
	if os.Getenv("TRUST_MODE") == "1" {
		if os.Getenv("PROXY_HOP_COUNT") == "" {
			return Env{}, errors.New("PROXY_HOP_COUNT is not set")
		}
		proxyHopCount, err := strconv.Atoi(os.Getenv("PROXY_HOP_COUNT"))
		if err != nil {
			return Env{}, err
		}
		return Env{
			IsProxy:       true,
			TrustMode:     1,
			ProxyHopCount: proxyHopCount,
		}, nil
		// } else if os.Getenv("TRUST_MODE") == "2" {
		// 	return Env{
		// 		IsProxy:   true,
		// 		TrustMode: 1,
		// 	}, nil
		// } else if os.Getenv("TRUST_MODE") == "3" {
		// 	if os.Getenv("TRUST_ADDRESS") == "" {
		// 		return Env{}, errors.New("TRUST_ADDRESS is not set")
		// 	}
		// 	return Env{
		// 		IsProxy:      true,
		// 		TrustMode:    3,
		// 		TrustAddress: strings.Split(os.Getenv("TRUST_ADDRESS"), ","),
		// 	}, nil
	} else {
		return Env{
			IsProxy: false,
		}, nil
	}
}
