package main

import (
  "math/rand"
  "net/http"
  "time"
)

/*
#cgo CFLAGS: -I/home/vagrant/dynatrace-6.1.0/adk/include
#cgo LDFLAGS: -L/home/vagrant/dynatrace-6.1.0/adk/lib64 -l dtadk
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

struct DTContext enter_dt()
{
  //TODO - figure out how to specify method name, not pick it up via __FUNCTION__
  //DYNATRACE_API("Go Lang");
  static const char *__dynatrace_api__ = "Go Lang";

  //DYNATRACE_START_PUREPATH();
  static int __dynatrace_method_id__ = 0;
  int __dynatrace_serial_no__;
  if (!__dynatrace_method_id__) {
		__dynatrace_method_id__ = dynatrace_get_method_id(__DYNATRACE_CURRENT_FUNCTION__, __DYNATRACE_CURRENT_SOURCE_FILE__, __DYNATRACE_CURRENT_SOURCE_LINENO__, __dynatrace_api__, 0);
	}
	__dynatrace_serial_no__ = dynatrace_get_serial_no(__dynatrace_method_id__, 1);
  __dynatrace_serial_no__ = dynatrace_enter(__dynatrace_method_id__, __dynatrace_serial_no__);




  struct DTContext ctx = {__dynatrace_method_id__, __dynatrace_serial_no__};
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


func init() {
  println("init here")
  C.init_dt()
}

func clientCall(w http.ResponseWriter) {
  ctx := C.enter_dt()

  delayBase := 1 + rand.Intn(5)
  time.Sleep(time.Duration(delayBase*300) * time.Millisecond)
	w.Write([]byte("Here's some content, dog\n"))

  C.exit_dt(ctx)
}


func handleCall(w http.ResponseWriter, r *http.Request) {
  ctx := C.enter_dt()

  delayBase := 1 + rand.Intn(5)
  time.Sleep(time.Duration(delayBase*100) * time.Millisecond)
	w.Write([]byte("Here's some content, dog\n"))

  clientCall(w);

  C.exit_dt(ctx)

}

func main() {
  handleCallHandler := http.HandlerFunc(handleCall)
  http.ListenAndServe(":8080", handleCallHandler)
}
