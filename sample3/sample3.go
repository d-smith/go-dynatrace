package main

/*
#cgo CFLAGS: -I/home/vagrant/dynatrace-6.1.0/adk/include
#cgo LDFLAGS: -L/home/vagrant/dynatrace-6.1.0/adk/lib64 -l dtadk
#include <stdlib.h>
#include <dynatrace_adk.h>
#include <stdio.h>

extern void go_sleepy();
extern sleepy();

struct dt_context {
  int methodId;
  int serialNo;
};

struct dt_context dt_start_pp() {
  DYNATRACE_API("Golang Sample II");
  DYNATRACE_START_PUREPATH();
  struct dt_context ctx = { __dynatrace_method_id__, __dynatrace_serial_no__};
  return ctx;
}

void dt_exit(struct dt_context ctx) {
  int __dynatrace_method_id__ = ctx.methodId;
  int __dynatrace_serial_no__ = ctx.serialNo;
  DYNATRACE_EXIT();
}

struct dt_context dt_enter() {
  DYNATRACE_ENTER(); 
  struct dt_context ctx = { __dynatrace_method_id__, __dynatrace_serial_no__};
  return ctx;
}



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


static void shutdown_dt() {
  DYNATRACE_UNINITIALIZE();
}




*/
import "C"

import (
  "fmt"
  "os"
  "os/signal"
  "time"
  "strconv"
)

func init() {
//  C.init_dt()
}

func B() {
  ctx := C.dt_enter()
  fmt.Printf("\tB called ctx is %d - %d\n", ctx.methodId, ctx.serialNo)
  time.Sleep(400 * time.Millisecond)
  C.dt_exit(ctx)
}

func C() {
  ctx := C.dt_enter()
  time.Sleep(400 * time.Millisecond)
  fmt.Printf("\tC called ctx is %d - %d\n", ctx.methodId, ctx.serialNo);
  C.dt_exit(ctx)
}

func A() {
  
  ctx := C.dt_start_pp()
  fmt.Printf("A called ctx is %d - %d\n", ctx.methodId, ctx.serialNo)
  B()
  C()
  C.dt_exit(ctx);
}

var finishedChannel = make(chan bool)
var stopped = false

func doStuff2() {
  for ; stopped == false; {
    A()
    time.Sleep(400 * time.Millisecond)
  }
  C.shutdown_dt()
  fmt.Println("shutdown complete")
  finishedChannel <- true
}

func doStuff1() {
  for ; stopped == false; {
    time.Sleep(1000 * time.Millisecond)
    A()
  }
  C.shutdown_dt()
  fmt.Println("shutdown complete")
  finishedChannel <- true
}

func main() {
  argsWithoutProg := os.Args[1:]
  if len(argsWithoutProg) != 1 {
    println("need a single args that's the number of go routines to spawn")
    return
  }

  println(argsWithoutProg[0])
  goroutines, err := strconv.Atoi(argsWithoutProg[0])
  if err != nil {
    println(err.Error())
    return
  }

  C.init_dt()

  for i := 0; i < goroutines; i++ {
    go doStuff1()
  }
  
  signalChannel := make(chan os.Signal, 1)
  
  signal.Notify(signalChannel, os.Interrupt)
  go func() {
    for _ = range signalChannel {
      fmt.Println("Received interrupt - shutting down")
      stopped = true
    }
  }()

  for i := 0; i < goroutines; i++ {
    <- finishedChannel
  }
  
}
