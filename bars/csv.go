package bars

import (
	"os"
	"time"
	"bytes"
	"context"
	"strings"
	"strconv"
	pkgerr "github.com/pkg/errors"
	goframe "github.com/rocketlaunchr/dataframe-go"
	"github.com/rocketlaunchr/dataframe-go/imports"
)



// CSVBars .
type CSVBars struct {
	opts 	Options
	k 		[]*Bar
}


// NewCSVBars .
func NewCSVBars(opt ...Option) *CSVBars {
	opts := newOptions(opt...)
	return &CSVBars{
		opts: opts,
	}
}

// Init .
func (cb *CSVBars) Init(opts ...Option) (err error) {
	var (
		file *os.File
		df *goframe.DataFrame
	)

	for _, o := range opts {
		o(&cb.opts)
	}
	if cb.opts.Path == "" {
		return pkgerr.New("The path field required")
	}
	if file, err = os.Open(cb.opts.Path); err != nil {
		return
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	if df, err = imports.LoadFromCSV(context.TODO(), strings.NewReader(buf.String())); err != nil {
		return
	}
	cb.k = make([]*Bar, 0)
	if err = cb.handleDataFrame(df); err != nil {
		return
	}
	return nil
}

func (cb *CSVBars) isDiscard() {
	var begin, end int
	if cb.opts.From.IsZero() && cb.opts.To.IsZero() {
		return
	} else if !cb.opts.From.IsZero() && cb.opts.To.IsZero() {
		for i, j := range cb.k {
			if cb.opts.From.Equal(time.Unix(j.Timestamp, 0)) {
				begin = i
				break
			}
		}
		cb.k = cb.k[begin:]
		return
	} else if cb.opts.From.IsZero() && !cb.opts.To.IsZero() {
		for i, j := range cb.k {
			if cb.opts.To.Equal(time.Unix(j.Timestamp, 0)) {
				end = i
				break
			}
		}
		cb.k = cb.k[:end+1]
		return
	} else if !cb.opts.From.IsZero() && !cb.opts.To.IsZero() {
		for i, j := range cb.k {
			if cb.opts.From.Equal(time.Unix(j.Timestamp, 0)) {
				begin = i
			}
			if cb.opts.To.Equal(time.Unix(j.Timestamp, 0)) {
				end = i
				break
			}
		}
		cb.k = cb.k[begin:end+1]
		return
	}
}

func (cb *CSVBars) handleDataFrame(df *goframe.DataFrame) error {

	iterator := df.ValuesIterator(
		goframe.ValuesOptions{
		InitialRow:0,
		Step:1,
		DontReadLock:true},
	)

	for {
		row, vals, _ := iterator()
 		if row == nil {
  			break
		}
		data := cb.convert(vals)
		bar, err := cb.convertToBar(data)
		if err != nil {
			return err
		}
		cb.k = append(cb.k, bar)
	}
	cb.isDiscard()
	return nil
}


func (cb *CSVBars) convertToBar(cdata map[string]interface{}) (bar *Bar, err error) {
	var date time.Time
	switch cb.opts.TimeFrame {
	case Minute:
		date, err = time.Parse("2006-01-02 15:04:05", cdata["date"].(string))
	default:
		date, err = time.Parse("2006-01-02", cdata["date"].(string))
	}
	if err != nil {
		return &Bar{}, err
	}
	bar = &Bar{
		Timestamp: date.Unix(),
		Open: cdata["open"].(float64),
		High: cdata["high"].(float64),
		Low: cdata["low"].(float64),
		Close: cdata["close"].(float64),
		Volume: cdata["volume"].(float64),
	}
	return
}


func (cb *CSVBars) convert(vals map[interface{}]interface{}) (map[string]interface{}) {
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
