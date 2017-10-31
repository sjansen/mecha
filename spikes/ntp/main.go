package main

import (
	"fmt"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	options := ntp.QueryOptions{Timeout: 30 * time.Second}
	if response, err := ntp.QueryWithOptions("0.beevik-ntp.pool.ntp.org", options); err != nil {
		fmt.Println(err)
	} else {
		err = response.Validate()
		fmt.Println("offset:", response.ClockOffset)
		fmt.Println("poll:", response.Poll)
		fmt.Println("stratum:", response.Stratum)
		fmt.Println("valid:", err == nil)
	}
}
