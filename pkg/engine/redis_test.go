package engine

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ylh990835774/ay-go-components/pkg/common"
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
		},
		{
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
		},
		{
			"hset key1 key2 val2",
			false,
		},
		{
			"  get key2",
			true,
		},
		{
			`hset
			get key2 val2`,
			false,
		},
		{
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
