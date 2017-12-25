package gdax

import (
	"restful_query"
	json2 "encoding/json"
	"errors"
	"github.com/golang/glog"
	"strconv"
	"time"
)

type GdaxStatsQuery struct {
	Product        string `json:"id"`
	Open           string `json:"open"`
	High           string `json:"high"`
	Low            string `json:"low"`
	Volume         string `json:"volume"`
	Last           string `json:"last"`
	Volume30Day    string `json:"volume_30day"`
	QueryTimeStamp int64  `json:"query_timestamp"`
}

type GdaxStats struct {
	Product        string  `json:"id"`
	Open           float64 `json:"open"`
	High           float64 `json:"high"`
	Low            float64 `json:"low"`
	Volume         float64 `json:"volume"`
	Last           float64 `json:"last"`
	Volume30Day    float64 `json:"volume_30day"`
	QueryTimeStamp int64   `json:"query_timestamp"`
}

//noinspection ALL
func get_stats() ([]GdaxStats, error) {
	onlineProducts, _, err := get_online_products()
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}
	var stats []GdaxStats

	for _, product := range onlineProducts {
		bodyBytes, err := restful_query.Get(apiUrl + "products/" + product.Id + "/stats")
		if err != nil {
			glog.Errorln(err)
			return nil, err
		}
		var queryStat = GdaxStatsQuery{
			Product:        product.Id,
			QueryTimeStamp: time.Now().Unix(),
		}
		json2.Unmarshal(bodyBytes, &queryStat)

		var floatStat GdaxStats
		convert_stats(&queryStat, &floatStat)
		stats = append(stats, floatStat)
	}

	return stats, nil
}

//noinspection ALL
func get_product_stats(productId string) (GdaxStats, error) {
	var productQStat = GdaxStatsQuery{
		Product:        productId,
		QueryTimeStamp: time.Now().Unix()}

	var productStats GdaxStats

	bodyBytes, err := restful_query.Get(apiUrl + "products/" + productId + "/stats")
	if err != nil {
		glog.Errorln(err)
		return productStats, err
	}
	json2.Unmarshal(bodyBytes, &productQStat)

	/*
		some of GDAX's products return 0 because they are not setup yet when they are
		converted to strings they end up as empty strings so here we disreguard them
	*/
	if productQStat.Last == "" && productQStat.Open == "" && productQStat.Low == "" &&
		productQStat.Volume30Day == "" && productQStat.Volume == "" && productQStat.High == "" {
		glog.Warningln("GDAX stats for product: " + productId + " are not populated, skipping")
		err := errors.New("GDAX stats for product: " + productId + " are not populated, skipping")
		return productStats, err
	}

	convert_stats(&productQStat, &productStats)

	return productStats, nil
}

//noinspection ALL
func convert_stats(stringStatsStruct *GdaxStatsQuery, floatStatsStruct *GdaxStats) {
	floatStatsStruct.Product = stringStatsStruct.Product
	floatStatsStruct.QueryTimeStamp = stringStatsStruct.QueryTimeStamp
	fOpen, errO := strconv.ParseFloat(stringStatsStruct.Open, 64)
	fLast, errL := strconv.ParseFloat(stringStatsStruct.Last, 64)
	fHigh, errH := strconv.ParseFloat(stringStatsStruct.High, 64)
	fVolume, errV := strconv.ParseFloat(stringStatsStruct.Volume, 64)
	f30Day, err3 := strconv.ParseFloat(stringStatsStruct.Volume30Day, 64)
	fLow, errLow := strconv.ParseFloat(stringStatsStruct.Low, 64)
	errs := []error{errO, errL, errH, errV, err3, errLow}
	for _, err := range errs {
		if err != nil {
			glog.Errorln("Unable to convert GDAX strings to floats: ", err)
		}
	}
	floatStatsStruct.Open = fOpen
	floatStatsStruct.Last = fLast
	floatStatsStruct.High = fHigh
	floatStatsStruct.Volume = fVolume
	floatStatsStruct.Volume30Day = f30Day
	floatStatsStruct.Low = fLow
}
