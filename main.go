package main

import (
  "fmt"
  "math/rand"
  "net/http"
  "time"
)

/*
#cgo CFLAGS: -I/home/vagrant/dynatrace-6.1.0/adk/include
#cgo LDFLAGS: -L/home/vagrant/dynatrace-6.1.0/adk/lib64 -l dtadk
#include <stdlib.h>
#include <dynatrace_adk.h>

struct DTContext {
  int methodId;
  int serialNo;
};

static void init_dt()
{
//TODO - figure out how to use a golang slice in C land

int argCount = 2;
char *strs[2] = {"dummy","--dt_debugadk=true"};
char **argv = &strs;
DYNATRACE_INITIALIZE(&argCount,&argv);

}

int getDynatraceMethodId(char *methodName) {
  return dynatrace_get_method_id(methodName, "XAVI II", 666, "Gopher3", 0);
}

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

var handleCallId C.int
var childCallId C.int

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

func childCall(w http.ResponseWriter) {
  ctx := C.enter_dt(childCallId)
  fmt.Printf("method id: %v sequence: %v\n", ctx.methodId, ctx.serialNo)


  delayBase := 1 + rand.Intn(5)
  time.Sleep(time.Duration(delayBase*300) * time.Millisecond)
	w.Write([]byte("Here's some content, dog\n"))

  C.exit_dt(ctx)
}


func handleCall(w http.ResponseWriter, r *http.Request) {

  ctx := C.enter_dt(handleCallId)
  fmt.Printf("method id: %v sequence: %v\n", ctx.methodId, ctx.serialNo)

  delayBase := 1 + rand.Intn(5)
  time.Sleep(time.Duration(delayBase*100) * time.Millisecond)
	w.Write([]byte("Here's some content, dog\n"))

  childCall(w);

  C.exit_dt(ctx)

}

func main() {
  handleCallHandler := http.HandlerFunc(handleCall)
  http.ListenAndServe(":8080", handleCallHandler)
}
