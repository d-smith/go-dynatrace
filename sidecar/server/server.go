package server

import (
	"fmt"
	"github.com/d-smith/go-dynatrace/sidecar/timing"
	"math/rand"
	"net/http"
	"time"
)

func childCall(w http.ResponseWriter) {

	fmt.Printf("childCall\n")

	delayBase := 1 + rand.Intn(5)
	time.Sleep(time.Duration(delayBase*300) * time.Millisecond)
	w.Write([]byte("Here's some child content, dog\n"))
}

//HandleCall handles http service requests
func HandleCall(w http.ResponseWriter, r *http.Request) {

	//Start API timing
	var svcTiming timing.ServiceTimer
	svcTiming.Name = "HandleCall"
	now := time.Now()

	//Simulate processing in the API layer
	fmt.Println("HandleCall")
	delayBase := 1 + rand.Intn(5)
	time.Sleep(time.Duration(delayBase*100) * time.Millisecond)
	w.Write([]byte("Here's some content, dog\n"))

	//Start proxied service timing
	childStart := time.Now()

	//Simulate proxied service call
	childCall(w)

	//Capture service call timing
	childElapsed := time.Since(childStart)
	childCallTime := timing.CallTime{"childCall", int(childElapsed / 1000000)}

	//Capture api timing
	svcTiming.Millis = int(time.Since(now) / 1000000)
	svcTiming.ChildCalls = append(svcTiming.ChildCalls, childCallTime)

	//Send timing data to go routine for processing
	go timing.WriteTimingData(svcTiming)
}
