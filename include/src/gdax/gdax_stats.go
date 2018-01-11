package gdax

import (
	json2 "encoding/json"
	"errors"
	"github.com/golang/glog"
	"restful_query"
	"strconv"
	"time"
)

const statsDataOldDuration = 10 * time.Second

type gdaxStatsApiResponse struct {
	Product        string `json:"id"`
	Open           string `json:"open"`
	High           string `json:"high"`
	Low            string `json:"low"`
	Volume         string `json:"volume"`
	Last           string `json:"last"`
	Volume30Day    string `json:"volume_30day"`
	QueryTimestamp int64  `json:"query_timestamp"`
}

type gdaxStatApi struct {
	Product        string  `json:"id"`
	Open           float64 `json:"open"`
	High           float64 `json:"high"`
	Low            float64 `json:"low"`
	Volume         float64 `json:"volume"`
	Last           float64 `json:"last"`
	Volume30Day    float64 `json:"volume_30day"`
	QueryTimestamp int64   `json:"query_timestamp"`
}

//noinspection ALL
func get_stats() ([]gdaxStatApi, error) {
	onlineProducts, _, err := get_online_products()
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}
	var stats []gdaxStatApi

	for _, product := range onlineProducts {
		bodyBytes, err := restful_query.Get(apiUrl + "products/" + product.Id + "/stats")
		if err != nil {
			glog.Errorln(err)
			return nil, err
		}
		var queryStat = gdaxStatsApiResponse{
			Product:        product.Id,
			QueryTimestamp: time.Now().Unix(),
		}
		json2.Unmarshal(bodyBytes, &queryStat)

		var floatStat gdaxStatApi
		convert_stats(&queryStat, &floatStat)
		stats = append(stats, floatStat)
	}

	return stats, nil
}

//noinspection ALL
func get_product_stats(productId string) (gdaxStatApi, error) {
	var productStatApi = gdaxStatsApiResponse{
		Product:        productId,
		QueryTimestamp: time.Now().Unix()}

	var productStats gdaxStatApi

	bodyBytes, err := restful_query.Get(apiUrl + "products/" + productId + "/stats")
	if err != nil {
		glog.Errorln(err)
		return productStats, err
	}
	json2.Unmarshal(bodyBytes, &productStatApi)

	/*
		some of GDAX's products return 0 because they are not setup yet when they are
		converted to strings they end up as empty strings so here we disreguard them
	*/
	if productStatApi.Last == "" && productStatApi.Open == "" && productStatApi.Low == "" &&
		productStatApi.Volume30Day == "" && productStatApi.Volume == "" && productStatApi.High == "" {
		glog.Warningln("GDAX stats for product: " + productId + " are not populated, skipping")
		err := errors.New("GDAX stats for product: " + productId + " are not populated, skipping")
		return productStats, err
	}

	convert_stats(&productStatApi, &productStats)

	return productStats, nil
}

//noinspection ALL
func convert_stats(stringStatsStruct *gdaxStatsApiResponse, floatStatsStruct *gdaxStatApi) {
	//floatStatsStruct.Product = stringStatsStruct.Product
	//floatStatsStruct.QueryTimestamp = stringStatsStruct.QueryTimestamp
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
	*floatStatsStruct = gdaxStatApi{
		Product:        stringStatsStruct.Product,
		QueryTimestamp: stringStatsStruct.QueryTimestamp,
		Open:           fOpen,
		Last:           fLast,
		High:           fHigh,
		Volume:         fVolume,
		Volume30Day:    f30Day,
		Low:            fLow,
	}
	//floatStatsStruct.Open = fOpen
	//floatStatsStruct.Last = fLast
	//floatStatsStruct.High = fHigh
	//floatStatsStruct.Volume = fVolume
	//floatStatsStruct.Volume30Day = f30Day
	//floatStatsStruct.Low = fLow
}

type gdaxStats struct {
	Stats          []gdaxStatApi
	queryTimestamp int64
}

func init_stats() *gdaxStats {
	var gS = gdaxStats{}
	gS.update_stats(gP)

	return &gS
}

func (gS *gdaxStats) update_stats(gP *gdaxProducts) {
	var timestamp = time.Now().Unix()
	var tmpStat gdaxStatApi
	var tmpStats []gdaxStatApi
	var err error

	for _, product := range *gP.get_online_products() {
		tmpStat, err = get_product_stats(product.Id)
		if err != nil {
			continue
		} else {
			tmpStats = append(tmpStats, tmpStat)
		}
	}

	*gS = gdaxStats{
		Stats:          tmpStats,
		queryTimestamp: timestamp,
	}
}

func (gS *gdaxStats) update_data(force bool) {
	if force {
		gS.update_stats(gP)

		return
	} else {
		if dataAge := time.Since(time.Unix(gS.queryTimestamp, 0)); dataAge > statsDataOldDuration {
			gS.update_stats(gP)
		}
	}
}

func (gS *gdaxStats) get_product_stats(productId string) *gdaxStatApi {
	for _, stat := range gS.Stats {
		if stat.Product == productId {
			return &stat
		}
	}

	return nil
}

func (gS *gdaxStats) Test() *gdaxStatApi {
	//return &gS.Stats
	return gS.get_product_stats("BTC-USD")
}