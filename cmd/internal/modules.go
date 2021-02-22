package internal

import (
	"github.com/DoNewsCode/core"
	"github.com/DoNewsCode/core/config"
	"github.com/DoNewsCode/core/otgorm"
	"github.com/DoNewsCode/core/ots3"
	"github.com/DoNewsCode/core/srvhttp"
	"github.com/DoNewsCode/skeleton/app"
)

func register(c *core.C) {
	c.AddModuleFunc(app.InjectAppModule)
	c.AddModuleFunc(ots3.New)
	c.AddModuleFunc(config.New)
	c.AddModuleFunc(core.NewServeModule)
	c.AddModuleFunc(otgorm.New)
	c.AddModule(srvhttp.DebugModule{})
	c.AddModule(srvhttp.DocsModule{})
	c.AddModule(srvhttp.HealthCheckModule{})
	c.AddModule(srvhttp.MetricsModule{})
}
