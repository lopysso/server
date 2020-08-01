//+build wireinject

package dependency_injection

import (
	"github.com/google/wire"
	"github.com/lopysso/server/config"
	configService "github.com/lopysso/server/service/config"
)

func InjectConfig() config.Config {
	panic(wire.Build(configService.ConfigInstance))

}
