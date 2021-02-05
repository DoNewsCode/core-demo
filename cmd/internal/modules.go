package internal

import (
	"github.com/DoNewsCode/skeleton/app"
	"github.com/DoNewsCode/std/pkg/container"
	"github.com/DoNewsCode/std/pkg/core"
	"github.com/DoNewsCode/std/pkg/ots3"
	"github.com/DoNewsCode/std/pkg/srvhttp"
)

func register(c *core.C) {
	c.AddModuleViaFunc(app.InjectAppModule)
	c.AddModule(ots3.New(c.ConfigAccessor, c.LevelLogger))
	c.AddModule(container.HttpFunc(srvhttp.Doc))
	c.AddModule(container.HttpFunc(srvhttp.HealthCheck))
	c.AddModule(container.HttpFunc(srvhttp.Metrics))
	c.AddModule(container.HttpFunc(srvhttp.Debug))
}
