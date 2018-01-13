package scheduler

import (
	"gdax"
	"time"
)

const loopInterval = 3 * time.Second

func Start_Scheduler()  {
	go data_update_loop()
}

func data_update_loop() {
	for {
		time.Sleep(loopInterval)
		go gdax.Update_Data(false)
	}
}