package feeds

import (
	"os"
	"fmt"
	"time"
	"bytes"
	"context"
	"strings"
	"strconv"

	"gotop/utils"
	"gotop/brokers"
	goframe "github.com/rocketlaunchr/dataframe-go"
	"github.com/rocketlaunchr/dataframe-go/imports"
)




/*
CSVFeed 用于读取csv文件后转换为可操作数据
Compression: 压缩比，例如TimeFrame=Minutes,那么Compression=5表示数据是5分钟k线数据。
TimeFrame: 时间帧类型，Minutes、Days等类型
FromDate: 起始时期
ToDate: 结束日期
Format("2006-01-02 15:04:05")
*/
type CSVFeed struct {
	Compression		uint8
	TimeFrame		uint8
	FromDate		time.Time
	ToDate			time.Time
}



/* 
New 读取csv文件,把数据处理好后返回给cerebro使用
path: 文件路径
*/
func (c *CSVFeed) New(path, name string) (brokers.DBroker, error) {
	var err error
	broker := brokers.DBroker{Name: name}
	broker.Init()
	file, err := os.Open(path)
	if err != nil {
  		return brokers.DBroker{}, err
	}
	defer file.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	s := buf.String()
	ctx := context.TODO()
	df, err := imports.LoadFromCSV(ctx, strings.NewReader(s))
	if err != nil {
		return brokers.DBroker{}, err
	}
	kline, err := c.handleDF(df)
	if err != nil {
		return brokers.DBroker{}, err
	}
	broker.SetKLine(kline)
	return broker, nil
}

func (c *CSVFeed) isDiscard(k *utils.KLine) {
	var begin, end int
	if c.FromDate.IsZero() && c.ToDate.IsZero() {
		return
	} else if !c.FromDate.IsZero() && c.ToDate.IsZero() {
		for i, j := range k.OHLC {
			if c.FromDate.Equal(j.Date) {
				begin = i
				break
			}
		}
		k.OHLC = k.OHLC[begin:]
		return
	} else if c.FromDate.IsZero() && !c.ToDate.IsZero() {
		for i, j := range k.OHLC {
			if c.ToDate.Equal(j.Date) {
				end = i
				break
			}
		}
		k.OHLC = k.OHLC[:end+1]
		return
	} else if !c.FromDate.IsZero() && !c.ToDate.IsZero() {
		for i, j := range k.OHLC {
			if c.FromDate.Equal(j.Date) {
				begin = i
			}
			if c.ToDate.Equal(j.Date) {
				end = i
				break
			}
		}
		k.OHLC = k.OHLC[begin:end+1]
		return
	}
}

/*
HandleDF 处理读取到csv数据
*/
func (c *CSVFeed) handleDF(df *goframe.DataFrame) (*utils.KLine, error) {
	iterator := df.ValuesIterator(
		goframe.ValuesOptions{
		InitialRow:0,
		Step:1,
		DontReadLock:true},
	)
	var kline utils.KLine
	df.Lock()
	for {
		row, vals, _ := iterator()
 		if row == nil {
  			break
		}
		data := c.convert(vals)
		ca, err := c.handleCandle(data)
		if err != nil {
			return &kline, err
		}
		kline.Append(ca)
	}
	df.Unlock()
	c.isDiscard(&kline)
	return &kline, nil
}


func (c *CSVFeed) handleCandle(cdata map[string]interface{}) (utils.Candle, error) {
	var date time.Time
	var err error
	switch c.TimeFrame {
	case Minute:
		date, err = time.Parse("2006-01-02 15:04:05", cdata["date"].(string))
	default:
		date, err = time.Parse("2006-01-02", cdata["date"].(string))
	}
	if err != nil {
		fmt.Println("Date conversion failed")
		return utils.Candle{}, err
	}
	data := utils.Candle{
		Date: date,
		Open: cdata["open"].(float64),
		High: cdata["high"].(float64),
		Low: cdata["low"].(float64),
		Close: cdata["close"].(float64),
		Volume: cdata["volume"].(float64),
	}
	return data, nil
}


func (c *CSVFeed) convert(vals map[interface{}]interface{}) (map[string]interface{}) {
	data := map[string]interface{}{}
	for key, val := range vals {
		switch key := key.(type) {
		case string:
			switch val := val.(type) {
			case string:
				if key != "date" {
					p, _ := strconv.ParseFloat(val, 64)
					data[key] = p
				} else {
					data[key] = val
				}
			}
		}
	}
	return data
}
