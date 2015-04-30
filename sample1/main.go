package main

import (
  "fmt"
  "math/rand"
  "net/http"
  "sync"
  "time"
)

/*
#cgo CFLAGS: -I/home/vagrant/dynatrace-6.1.0/adk/include
#cgo LDFLAGS: -L/home/vagrant/dynatrace-6.1.0/adk/lib64 -l dtadk
#include <stdlib.h>
#include <dynatrace_adk.h>

// Structure for passing methodId and serialNo context back to
// the dynatrace exit function
struct DTContext {
  int methodId;
  int serialNo;
};

//Initialize dynatrace
static void init_dt()
{
  //TODO - look at the initialize macro and determine the
  //correct set up. Might be able to handle everything via environment
  //variables.

  int argCount = 2;
  char *strs[2] = {"dummy","--dt_debugadk=true"};
  char **argv = &strs;
  DYNATRACE_INITIALIZE(&argCount,&argv);

}

//TODO - install a shutdown hook to cleanly disconnect from dynatrace

static void shutdown_dt() {
  DYNATRACE_UNINITIALIZE();
}

struct DTContext api_call()
{
  DYNATRACE_API("Gopher");
  DYNATRACE_START_PUREPATH();
  struct DTContext ctx = {__dynatrace_method_id__, __dynatrace_serial_no__};
  return ctx;
}

struct DTContext svc_call()
{
  DYNATRACE_ENTER();
  struct DTContext ctx = {__dynatrace_method_id__, __dynatrace_serial_no__};
  return ctx;
}

// End the pure path contribution
void exit_dt(struct DTContext ctx)
{
  int __dynatrace_method_id__ = ctx.methodId;
  int __dynatrace_serial_no__ = ctx.serialNo;
  DYNATRACE_EXIT();
}




*/
import "C"

// Use a mutex to sync access to Dynatrace
var mutex sync.Mutex

//Initialize the dynatrace method ids
func init() {
  println("init here")
  C.init_dt()
}

//Sample child method called from handleCall - shown as descendant in pure path
func childCall(w http.ResponseWriter) {

  mutex.Lock()
  ctx := C.svc_call()
  mutex.Unlock()
  fmt.Printf("method id: %v sequence: %v\n", ctx.methodId, ctx.serialNo)


  delayBase := 1 + rand.Intn(5)
  time.Sleep(time.Duration(delayBase*300) * time.Millisecond)
	w.Write([]byte("Here's some content, dog\n"))

  mutex.Lock()
  C.exit_dt(ctx)
  mutex.Unlock()
}

//Top level pure path call
func handleCall(w http.ResponseWriter, r *http.Request) {


  mutex.Lock()
  ctx := C.api_call()
  mutex.Unlock()
  fmt.Printf("method id: %v sequence: %v\n", ctx.methodId, ctx.serialNo)

  delayBase := 1 + rand.Intn(5)
  time.Sleep(time.Duration(delayBase*100) * time.Millisecond)
	w.Write([]byte("Here's some content, dog\n"))

  childCall(w);

  mutex.Lock()
  C.exit_dt(ctx)
  mutex.Unlock()

}

//Serve up some http goodness
func serveHttp() {
  handleCallHandler := http.HandlerFunc(handleCall)
  http.ListenAndServe(":8080", handleCallHandler)
}


func main() {
  serveHttp()
}
