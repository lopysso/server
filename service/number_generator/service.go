package number_generator

import (
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/lopysso/server/config"
)

func NewSnowflakeNode(conf config.Config) *snowflake.Node {
	a, err := snowflake.NewNode(conf.NumberGenerator.Snowflake.Node)
	if err != nil {
		panic("number generator init error")
	}

	return a
}

var snowflakeNode *snowflake.Node
var snowflakeNodeOnce sync.Once

func InstanceSnowflakeNode(conf config.Config) *snowflake.Node {

	snowflakeNodeOnce.Do(func() {
		snowflakeNode = NewSnowflakeNode(conf)
	})

	return snowflakeNode
}
