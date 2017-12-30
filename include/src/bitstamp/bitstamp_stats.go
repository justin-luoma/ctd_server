package bitstamp

import (
	"restful_query"
	json2 "encoding/json"
	"github.com/golang/glog"
	"strconv"
)

type BitstampStatsQuery struct {
	Product   string
	High      string `json:"high"`
	Last      string `json:"last"`
	Timestamp int64 `json:"timestamp"`
	Bid       string `json:"bid"`
	Vwap      string `json:"vwap"`
	Volume    string `json:"volume"`
	Low       string `json:"low"`
	Ask       string `json:"ask"`
	Open      string `json:"open"`
}

type BitstampStats struct {
	Product   string
	High      float64 `json:"high"`
	Last      float64 `json:"last"`
	Timestamp int64   `json:"timestamp"`
	Bid       float64 `json:"bid"`
	Vwap      float64 `json:"vwap"`
	Volume    float64 `json:"volume"`
	Low       float64 `json:"low"`
	Ask       float64 `json:"ask"`
	Open      float64 `json:"open"`
}

//noinspection ALL
func get_stats() ([]BitstampStats, error) {
	onlineProducts, _, err := get_online_products()
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}
	var stats []BitstampStats

	for _, product := range onlineProducts {
		bodyBytes, err := restful_query.Get(apiUrl + "v2/ticker/" + product.URLSymbol)
		if err != nil {
			glog.Errorln(err)
			return nil, err
		}
		var queryStat = BitstampStatsQuery{
			Product:        product.Name,
		}
		json2.Unmarshal(bodyBytes, &queryStat)

		var floatStat BitstampStats
		convert_stats(&queryStat, &floatStat)
		stats = append(stats, floatStat)
	}

	return stats, nil
}

//noinspection ALL
func get_product_stats(productId string) (BitstampStats, error) {
	var productQStat = BitstampStatsQuery{
		Product:        productId,
	}

	var productStats BitstampStats

	bodyBytes, err := restful_query.Get(apiUrl + "v2/ticker/" + productId)
	if err != nil {
		glog.Errorln(err)
		return productStats, err
	}
	json2.Unmarshal(bodyBytes, &productQStat)

	convert_stats(&productQStat, &productStats)

	return productStats, nil
}

//noinspection ALL
func convert_stats(stringStatsStruct *BitstampStatsQuery, floatStatsStruct *BitstampStats) {
	floatStatsStruct.Product = stringStatsStruct.Product
	floatStatsStruct.Timestamp = stringStatsStruct.Timestamp
	fOpen, errO := strconv.ParseFloat(stringStatsStruct.Open, 64)
	fLast, errL := strconv.ParseFloat(stringStatsStruct.Last, 64)
	fHigh, errH := strconv.ParseFloat(stringStatsStruct.High, 64)
	fVolume, errV := strconv.ParseFloat(stringStatsStruct.Volume, 64)
	fLow, errLow := strconv.ParseFloat(stringStatsStruct.Low, 64)
	errs := []error{errO, errL, errH, errV, errLow}
	for _, err := range errs {
		if err != nil {
			glog.Errorln("Unable to convert Bitstamp strings to floats: ", err)
		}
	}
	floatStatsStruct.Open = fOpen
	floatStatsStruct.Last = fLast
	floatStatsStruct.High = fHigh
	floatStatsStruct.Volume = fVolume
	floatStatsStruct.Low = fLow
}
