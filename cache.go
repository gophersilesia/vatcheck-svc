package vatcheck

import (
	. "github.com/jgautheron/workshop/vat/config"
	"github.com/wunderlist/ttlcache"
)

var ch *ttlcache.Cache

// cache returns the ttlcache instance.
func cache() *ttlcache.Cache {
	if ch != nil {
		return ch
	}
	ch = ttlcache.NewCache(Config.CacheDuration)
	return ch
}
