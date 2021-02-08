package internal

import (
	"github.com/DoNewsCode/std/pkg/core"
	"github.com/DoNewsCode/std/pkg/observability"
	"github.com/DoNewsCode/std/pkg/otgorm"
	"github.com/DoNewsCode/std/pkg/otredis"
)

func provide(c *core.C) {
	c.ProvideItself()
	c.Provide(observability.Observability)
	c.Provide(otgorm.ProvideDefaultDatabase)
	c.Provide(otredis.ProvideDefaultRedis)
}
