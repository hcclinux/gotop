package strategy

import (
	"gotop/utils"
	"gotop/brokers"
)
// Strategy 策略基础接口
type Strategy interface {
	Next()
	Init()
	Set(utils.Candle)
	AddBroker(brokers.Broker)
}

const (
	// TOP .
	TOP = "top" 
	// BOTTOM .
	BOTTOM = "bottom"
	// UP .
	UP = "up"
	// DOWN .
	DOWN = "down"
)