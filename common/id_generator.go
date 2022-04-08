package common

import (
	"github.com/bwmarrin/snowflake"
	"github.com/resyon/jincai-im/log"
)

const (
	nodeNumber = 1
)

var (
	_node *snowflake.Node = initSnowFlakeNode()
)

func initSnowFlakeNode() *snowflake.Node {

	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(nodeNumber)
	if err != nil {
		log.LOG.Panicf("Fail to init snowflake node, Err=%+v", err)
	}
	return node
}

func GenerateID() int64 {
	id := _node.Generate()
	return id.Int64()
}
