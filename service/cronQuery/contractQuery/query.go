package contractQuery

import (
	"time"
)

func Init() {
	for {
		query()
		time.Sleep(10 * time.Second)
	}
}
