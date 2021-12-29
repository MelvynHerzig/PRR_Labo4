package wave

import (
	"server/network"
)

var messageChan = make(chan message)

func Handle() {
	for {
		select {
			case data := <- network.DataChan:

				if data.Message == "START" && !running {
					running = true
					go searchSP()
				} else {
					messageChan <- deserialize(data.Message)
				}
		}
	}
}
