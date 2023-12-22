package engine

import (
	"testing"

	"git.qpaas.com/go-components/webconsole/pkg/common"
	"github.com/stretchr/testify/require"
)

func TestIsRedisCMDSafe(t *testing.T) {
	testCase := []struct {
		redisCMD string
		isSafe   bool
	}{
		{
			"keys *",
			false,
		},
		{
			"ttl key01",
			true,
		}, {
			"type key02",
			true,
		},
		{
			"get key01",
			true,
		},
		{
			"hget key1 key2",
			true,
		}, {
			"hset key1 key2 val2",
			false,
		}, {
			"  get key2",
			true,
		}, {
			`hset 
			get key2 val2`,
			false,
		}, {
			`hget
			 key1 key2`,
			true,
		},
	}

	for _, item := range testCase {
		sql, isSafe, _ := IsRedisCMDSafe(item.redisCMD, common.DefaultRedisWhiteCMD)
		require.Equal(t, item.isSafe, isSafe, sql)
	}
}
