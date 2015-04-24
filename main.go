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

// Call this in go init method to get the method id to use for calls.
// This keeps us  from calling the dynatrace_get_method redundantly,
// which the macro from in the ADK avoids and reduces the amount
// of string allocation significantly.
int getDynatraceMethodId(char *methodName) {
  return dynatrace_get_method_id(methodName, "XAVI II", 666, "Gopher3", 0);
}

// Start a pure path contribution
struct DTContext enter_dt(int methodId)
{
  //DYNATRACE_API("Gopher");

  //DYNATRACE_START_PUREPATH();
  int __dynatrace_serial_no__;
	__dynatrace_serial_no__ = dynatrace_get_serial_no(methodId, 1);
  __dynatrace_serial_no__ = dynatrace_enter(methodId, __dynatrace_serial_no__);

  struct DTContext ctx = {methodId, __dynatrace_serial_no__};
  return ctx;
}

// End the pure path contribution
void exit_dt(struct DTContext ctx)
{
  int __dynatrace_method_id__ = ctx.methodId;
  int __dynatrace_serial_no__ = ctx.serialNo;

  //DYNATRACE_EXIT();

  dynatrace_exit(__dynatrace_method_id__, __dynatrace_serial_no__);
}


*/
import "C"
import "unsafe"

// Variables used to hold dynatrace method ids for instrumented methods
var (
  handleCallId C.int
  childCallId C.int
)

// Mutex for synchronizing Dynatrace ADK calls
var mutex sync.Mutex

//Initialize the dynatrace method ids
func init() {
  println("init here")
  C.init_dt()

  handleCall := C.CString("handleCall")
  defer C.free(unsafe.Pointer(handleCall))
  handleCallId = C.getDynatraceMethodId(handleCall)

  childCall := C.CString("childCall")
  defer C.free(unsafe.Pointer(childCall))
  childCallId = C.getDynatraceMethodId(childCall)

}

//Sample child method called from handleCall - shown as descendant in pure path
func childCall(w http.ResponseWriter) {
  mutex.Lock()
  ctx := C.enter_dt(childCallId)
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
  ctx := C.enter_dt(handleCallId)
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
func main() {
  handleCallHandler := http.HandlerFunc(handleCall)
  http.ListenAndServe(":8080", handleCallHandler)
}
