package main

/*
#cgo CFLAGS: -I/home/vagrant/dynatrace-6.1.0/adk/include
#cgo LDFLAGS: -L/home/vagrant/dynatrace-6.1.0/adk/lib64 -l dtadk
#include <stdlib.h>
#include <dynatrace_adk.h>
#include <stdio.h>



static void B() {
  DYNATRACE_ENTER(); 
  printf("\tB called\n");
  DYNATRACE_EXIT(); 
}

static void C() {
  DYNATRACE_ENTER(); 
  printf("\tC called\n");
  DYNATRACE_EXIT(); 
}

static void A() {
  DYNATRACE_API("Golang Sample II");
  DYNATRACE_START_PUREPATH();

  printf("A called\n"); 
  B();
  C();
  DYNATRACE_EXIT(); 
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

//TODO - install a shutdown hook to cleanly disconnect from dynatrace

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
)

func init() {
  C.init_dt()
}

var finishedChannel = make(chan bool)
var stopped = false

func doStuff2() {
  for ; stopped == false; {
    C.A()
    time.Sleep(1 * time.Second)
  }
  C.shutdown_dt()
  fmt.Println("shutdown complete")
  finishedChannel <- true
}

func doStuff1() {
  for ; stopped == false; {
    C.A()
    time.Sleep(2 * time.Second)
  }
  C.shutdown_dt()
  fmt.Println("shutdown complete")
  finishedChannel <- true
}

func main() {
  go doStuff1()
  go doStuff2()
  
  signalChannel := make(chan os.Signal, 1)
  
  signal.Notify(signalChannel, os.Interrupt)
  go func() {
    for _ = range signalChannel {
      fmt.Println("Received interrupt - shutting down")
      stopped = true
    }
  }()
  <- finishedChannel
  <- finishedChannel
  
}
