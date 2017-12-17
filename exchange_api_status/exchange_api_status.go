package exchange_api_status

import (
	"sync"
	"time"
	"github.com/golang/glog"
)

var exchanges = []string{"gdax", "coincap"}

/*
status holder map, values:
	0 : offline
	1 : online
	2 : hasn't checked in
*/
var statusLock = struct {
	sync.RWMutex
	exchangeStatusHolder map[string]*status
}{exchangeStatusHolder: make(map[string]*status)}

type status struct {
	Status      int
	LastUpdated int64
}

type Status struct {
	ExchangeName string
	Status       int
	LastUpdated  int64
}

func init() {
	statusLock.Lock()
	for _, exchange := range exchanges {
		exchangeHolder := statusLock.exchangeStatusHolder
		exchangeHolder[exchange] = &status{Status: 2,
			LastUpdated: time.Now().Unix()}
	}
	statusLock.Unlock()
}

func Update_Status(exchange string, newStatus int) {
	statusLock.Lock()
	exchangeHolder := statusLock.exchangeStatusHolder
	exchangeHolder[exchange] = &status{Status: newStatus,
		LastUpdated: time.Now().Unix()}
}

func check_status() []Status {
	statusLock.RLock()

	var exchangeStatuses []Status
	exchangeHolder := statusLock.exchangeStatusHolder
	statusLock.RUnlock()
	for exchange, value := range exchangeHolder {
		exchangeStatus := Status{ExchangeName: exchange,
			Status:      value.Status,
			LastUpdated: value.LastUpdated,
		}
		exchangeStatuses = append(exchangeStatuses, exchangeStatus)
	}
	return exchangeStatuses
}

func watch_status() {
	for {
		time.Sleep(5 * time.Second)
		exchangeStatuses := check_status()
		glog.Infoln(exchangeStatuses)
	}
}

func Start_Exchange_Monitoring()  {
	go watch_status()
}
