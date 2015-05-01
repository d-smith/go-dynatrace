package server

import (
	"fmt"
	"net/http"
	"math/rand"
	"time"
	"github.com/d-smith/go-dynatrace/sidecar/timing"
)

func childCall(w http.ResponseWriter) {

  fmt.Printf("childCall\n")

  delayBase := 1 + rand.Intn(5)
  time.Sleep(time.Duration(delayBase*300) * time.Millisecond)
  w.Write([]byte("Here's some child content, dog\n"))
}

//HandleCall handles http service requests
func HandleCall(w http.ResponseWriter, r *http.Request) {
	var svcTiming timing.ServiceTimer
	svcTiming.Name = "HandleCall"
	now := time.Now()
 
  fmt.Printf("HandleCall")

  delayBase := 1 + rand.Intn(5)
  time.Sleep(time.Duration(delayBase*100) * time.Millisecond)
	w.Write([]byte("Here's some content, dog\n"))
	childStart := time.Now()
    childCall(w);
    childElapsed := time.Since(childStart)
    childCallTime := timing.CallTime{"childCall", int(childElapsed/1000000)}

    svcTiming.Millis = int(time.Since(now)/1000000)
    svcTiming.ChildCalls = append(svcTiming.ChildCalls, childCallTime)

    
    go timing.WriteTimingData(svcTiming)
}