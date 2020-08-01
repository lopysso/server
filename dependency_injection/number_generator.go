//+build wireinject

package dependency_injection

import (
	"github.com/bwmarrin/snowflake"
	"github.com/google/wire"
	"github.com/lopysso/server/service/number_generator"
)


func InjectSnowflakeNode() *snowflake.Node {
	panic(wire.Build(InjectConfig,number_generator.InstanceSnowflakeNode))
}
