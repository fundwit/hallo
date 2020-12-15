package auth

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var TokenCache = cache.New(24*time.Hour, 1*time.Minute)
var RegisterTokenCache = cache.New(30*time.Minute, 1*time.Minute)
