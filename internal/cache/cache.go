package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var C *cache.Cache

func Init(defaultTTL time.Duration) {
	C = cache.New(defaultTTL, 2*defaultTTL)
}
