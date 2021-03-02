package internal

import (
	"github.com/DoNewsCode/core"
	"github.com/DoNewsCode/core/observability"
	"github.com/DoNewsCode/core/otgorm"
	"github.com/DoNewsCode/core/otredis"
	"github.com/DoNewsCode/core/ots3"
)

func provide(c *core.C) {
	c.Provide(observability.Providers())
	c.Provide(otgorm.Providers())
	c.Provide(otredis.Providers())
	c.Provide(ots3.Providers())
}
