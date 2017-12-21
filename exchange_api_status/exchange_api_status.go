package exchange_api_status

import (
	"github.com/golang/glog"
	"sync"
	"time"
)

var exchanges = []string{"gdax", "coincap"}

/*
Status holder map, values:
	0 : offline
	1 : online
	2 : hasn't checked in
*/
var statusLock = struct {
	sync.RWMutex
	exchangeStatusHolder map[string]*Status
}{exchangeStatusHolder: make(map[string]*Status)}

type Status struct {
	ExchangeName string
	Status       int
	LastUpdated  int64
}

/*
type Status struct {
	Status       int
	LastUpdated  int64
}
*/

func init() {
	statusLock.Lock()
	defer statusLock.Unlock()
	for _, exchange := range exchanges {
		exchangeHolder := statusLock.exchangeStatusHolder
		exchangeHolder[exchange] = &Status{Status: 2,
			LastUpdated: time.Now().Unix()}
	}
}

func Update_Status(exchange string, newStatus int) {
	statusLock.Lock()
	defer statusLock.Unlock()
	exchangeHolder := statusLock.exchangeStatusHolder
	exchangeHolder[exchange] = &Status{Status: newStatus,
		LastUpdated: time.Now().Unix()}
}

func check_status() []Status {
	statusLock.RLock()
	defer statusLock.RUnlock()

	var exchangeStatuses []Status
	exchangeHolder := statusLock.exchangeStatusHolder
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

func Start_Exchange_Monitoring() {
	go watch_status()
}
