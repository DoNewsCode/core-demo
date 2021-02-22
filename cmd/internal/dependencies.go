package internal

import (
	"github.com/DoNewsCode/core"
	"github.com/DoNewsCode/core/observability"
	"github.com/DoNewsCode/core/otgorm"
	"github.com/DoNewsCode/core/otredis"
	"github.com/DoNewsCode/core/ots3"
)

func provide(c *core.C) {
	c.Provide(observability.Provide)
	c.Provide(otgorm.Provide)
	c.Provide(otredis.Provide)
	c.Provide(ots3.Provide)
}
