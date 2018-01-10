package gdax

import (
	json2 "encoding/json"
	"errors"
	"exchange_api_status"
	"github.com/golang/glog"
	"restful_query"
	"time"
)

const productDataOldDuration = 1 * time.Minute

type gdaxProductsApi struct {
	Id            string `json:"id,omitempty"`
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
	Status        string `json:"status,omitempty"`
	StatusMessage string `json:"status_message,omitempty"`
}

//noinspection ALL
func get_products() (*[]gdaxProductsApi, error) {
	bodyBytes, err := restful_query.Get(apiUrl + "products")
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	var product []gdaxProductsApi
	json2.Unmarshal(bodyBytes, &product)

	return &product, nil
}

//noinspection ALL
func get_online_products() ([]gdaxProductsApi, []gdaxProductsApi, error) {
	products, err := get_products()
	if err != nil {
		glog.Errorln(err)
		return nil, nil, err
	}

	var onlineProducts []gdaxProductsApi
	var offlineProducts []gdaxProductsApi

	for _, product := range *products {
		if product.Status == "online" {
			onlineProducts = append(onlineProducts, product)
			onlineProductIds = append(onlineProductIds, product.Id)
		} else {
			offlineProducts = append(offlineProducts, product)
			glog.Warningln("skipped GDAX product: " + product.Id + " status: " + product.Status)
		}
	}

	return onlineProducts, offlineProducts, nil
}

type gdaxProduct struct {
	Id            string `json:"id,omitempty"`
	BaseCurrency  string `json:"base_currency,omitempty"`
	QuoteCurrency string `json:"quote_currency,omitempty"`
	IsActive      bool   `json:"is_active,omitempty"`
	StatusMessage string `json:"status_message,omitempty"`
}

func update_products() *[]gdaxProduct {

	var tmpProduct gdaxProduct
	var tmpProducts []gdaxProduct

	apiProducts, err := get_products()
	if err != nil || len(*apiProducts) < 1 {
		tmpProducts = []gdaxProduct{}

		exchange_api_status.Update_Status("gdax", 0)
		return &tmpProducts
	}

	for _, apiProduct := range *apiProducts {

		var isActive bool
		if apiProduct.Status == "online" {
			isActive = true
		} else {
			isActive = false
		}

		var statusMessage = ""
		if apiProduct.StatusMessage != "" {
			statusMessage = apiProduct.StatusMessage
		}

		tmpProduct = gdaxProduct{
			Id:            apiProduct.Id,
			BaseCurrency:  apiProduct.BaseCurrency,
			QuoteCurrency: apiProduct.QuoteCurrency,
			IsActive:      isActive,
			StatusMessage: statusMessage,
		}

		tmpProducts = append(tmpProducts, tmpProduct)
	}

	exchange_api_status.Update_Status("gdax", 1)

	return &tmpProducts
}

type gdaxProducts struct {
	Products       []gdaxProduct
	queryTimestamp int64
}

func init_products() *gdaxProducts {
	var gP = gdaxProducts{
		queryTimestamp: time.Now().Unix(),
	}
	gP.Products = *update_products()

	return &gP
}

func (gP *gdaxProducts) Test() *[]gdaxProduct {
	return &gP.Products
}

func (gP *gdaxProducts) update_data(force bool) {
	if force {
		gP.queryTimestamp = time.Now().Unix()
		gP.Products = *update_products()

		return
	} else {
		if dataAge := time.Since(time.Unix(gP.queryTimestamp, 0)); dataAge > productDataOldDuration {
			gP.queryTimestamp = time.Now().Unix()
			gP.Products = *update_products()
		}
	}
}

func (gP *gdaxProducts) get_online_products() *[]gdaxProduct {
	var tmpProduct gdaxProduct
	var tmpProducts []gdaxProduct
	for _, product := range gP.Products {
		if product.IsActive {
			tmpProduct = product
			tmpProducts = append(tmpProducts, tmpProduct)
		}
	}

	return &tmpProducts
}

func (gP *gdaxProducts) get_products() *[]gdaxProduct {
	return &gP.Products
}

func (gP *gdaxProducts) get_product(id string) (*gdaxProduct, error) {
	for _, product := range gP.Products {
		if id == product.Id {
			return &product, nil
		}
	}
	err := errors.New("invalid product id: " + id)

	return nil, err
}
