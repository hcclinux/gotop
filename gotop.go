package gotop


import (
	"reflect"
	"encoding/json"
	"os"
	"fmt"

	"gotop/utils"
	"gotop/brokers"
	"gotop/strategy"
)
/*
Cerebro 大脑是一个主要入口
整个回测过程都是在此结构下完成
*/
type Cerebro struct {
	Broker 		brokers.Broker
	strategy 	strategy.Strategy
}

// StatsPrint 输出统计结果
func (c *Cerebro) StatsPrint() {

}

// MakeData 生成数据给图表调用
func (c Cerebro) MakeData() {
	k := c.Broker.GetKLine()
	filePtr, err := os.Create("./datas.json")
    if err != nil {
        fmt.Println("创建文件失败，err=",err)
        return
    }
    defer filePtr.Close()
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(k)
    if err != nil {
        fmt.Println("编码失败，err=",err)
    } else {
        fmt.Println("编码成功")
    }

}


// AddBroker 添加经纪人
func (c *Cerebro) AddBroker(b brokers.Broker) {
	c.Broker = b
}


// AddStrategy 添加策略
func (c *Cerebro) AddStrategy(s strategy.Strategy) {
	c.strategy = s
}


// Run 启动运行回测
func (c *Cerebro) Run() {
	c.strategy.Init()
	c.strategy.AddBroker(c.Broker)
	for {
		candle := c.Broker.Next()
		if reflect.DeepEqual(candle, utils.Candle{}) {
			break
		}
		c.strategy.Set(candle)
		c.strategy.Next()
		c.Broker.SetIndex()
	}
	c.Broker.Print()
}