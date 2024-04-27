package stock

import (
	"fmt"
	"math"
	"strconv"
	"time"

	pd "github.com/SiverPineValley/parseduration"
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/sirupsen/logrus"
    c "github.com/wcharczuk/go-chart/v2"
)

type MinMaxGraphForDurationResponse struct {
	Min   float64
	Max   float64
	Graph c.Chart
}

type dailyData struct {
    close float64
    timestamp time.Time
  }
  
type totalData struct {
    min float64
    max float64
    data []dailyData
}

type IStockService interface {
	GetMinMaxGraphForDuration(symbol string, duration string) (MinMaxGraphForDurationResponse, error)
}

type stockService struct{
    logger *logrus.Logger
}

func (s *stockService) GetMinMaxGraphForDuration(symbol string, duration string) (MinMaxGraphForDurationResponse, error) {
    s.logger.Infof("symbol: %s, duration: %s", symbol, duration)
    
    now := time.Now()
    durationInTime, err := pd.ParseDuration(duration)
    if err != nil {
        return MinMaxGraphForDurationResponse{}, err
    }

    startTime := now.Add(-durationInTime)

    resp := chart.Get(&chart.Params{
        Symbol: symbol,
        Start: datetime.New(&startTime),
        End: datetime.New(&now),
        Interval: datetime.OneDay,
    })

    if err := resp.Err(); err != nil {
        return MinMaxGraphForDurationResponse{}, err
    }

    allData := totalData{min: math.MaxFloat64, data: make([]dailyData, 0)}
    for resp.Next() {
        bar := resp.Bar()
        if innerErr := s.processData(bar, &allData); innerErr != nil {
            return MinMaxGraphForDurationResponse{}, innerErr
        }
    }

    s.logger.Debugf("%#v", allData)

    xv := make([]time.Time,0)
    yv := make([]float64,0)

    for _, v := range allData.data {
        yv = append(yv, v.close)
        xv = append(xv, v.timestamp)
    }

    priceSeries := c.TimeSeries{
		Style: c.Style{
			StrokeColor: c.GetDefaultColor(0),
		},
		XValues: xv,
		YValues: yv,
	}

    graph := c.Chart{
        Title: fmt.Sprintf("%s - Max: %.2f | Min: %.2f", symbol, allData.max, allData.min),
		XAxis: c.XAxis{
			TickPosition: c.TickPositionBetweenTicks,
		},
		YAxis: c.YAxis{
			Range: &c.ContinuousRange{
				Max: allData.max,
				Min: allData.min,
			},
		},
		Series: []c.Series{
			priceSeries,
		},
	}

    return MinMaxGraphForDurationResponse{
        Min: allData.min,
        Max: allData.max,
        Graph: graph,
    }, nil
}

func New(logger *logrus.Logger) IStockService {
	return &stockService{logger: logger}
}

func addToSlice[T any](slice []T, val ...T) []T {
    return append(slice, val...)
}

func (s *stockService) processData(bar *finance.ChartBar, allData *totalData) error {
    s.logger.Debugf("working on: #v", bar)
    
    low, err := strconv.ParseFloat(bar.Low.StringFixed(2), 64)
    if err != nil {
       return err
    }

    if low < allData.min {
        if low == 0 {
            return nil
        }
        allData.min = low
    }

    high, err := strconv.ParseFloat(bar.High.StringFixed(2), 64)
    if err != nil {
       return err
    }

    if high > allData.max {
        allData.max = high
    }

    closeVal, err := strconv.ParseFloat(bar.Open.StringFixed(2), 64)
    if err != nil {
        return err
    }

    allData.data = addToSlice(allData.data, dailyData{close: closeVal, timestamp: time.Unix(int64(bar.Timestamp), 0)})
    return nil
}