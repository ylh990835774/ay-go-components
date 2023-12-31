package engine

import (
	"github.com/ylh990835774/ay-go-components/pkg/common"
)

type Engine interface {
	Schema() ([]string, error)
	Table(schema string) ([]string, error)
	Query(schema string, table string, sql string, timeout int64) *common.QuerySet

	RegistryQueryPrev(common.PreHook)
	RegistryQueryPost(common.PostHook)

	Close() error
}
