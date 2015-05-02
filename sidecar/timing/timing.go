package timing

/*
#cgo CFLAGS: -I/home/vagrant/dynatrace-6.1.0/adk/include
#cgo LDFLAGS: -L/home/vagrant/dynatrace-6.1.0/adk/lib64 -l dtadk
#include <stdlib.h>
#include <dynatrace_adk.h>
#include <stdio.h>
#include <time.h>

//C code for inserting golang information into the Dynatrace pure path.

//Initialize dynatrace
static void init_dt()
{
  int argCount = 2;
  char *strs[2] = {"dummy","--dt_debugadk=true"};
  char **argv = &strs;
  DYNATRACE_INITIALIZE(&argCount,&argv);

}

typedef struct {
	char* name;
	int timing;
	char *subname;
	int subtiming;
} TimingData;

void subcall(int subcallTime) {
	char buf[256];
	sprintf(buf, "proxied service call time: %d ms", subcallTime);

	DYNATRACE_ENTER();
	DYNATRACE_LOG(0, buf, "xavi");
	DYNATRACE_EXIT();
}

void api(int apiTime, int subcallTime) {
	char buf[256];
	sprintf(buf, "api call time: %dms", apiTime);

	DYNATRACE_ENTER();
	DYNATRACE_LOG(0, buf, "xavi");
	subcall(subcallTime);
	DYNATRACE_EXIT();
}

static void handleTiming(TimingData td) {
	DYNATRACE_API("Son of Sidecar");
	DYNATRACE_START_PUREPATH();
	api(td.timing, td.subtiming);
	DYNATRACE_EXIT();
}


//TODO: Add dt shutdown hook.


*/
import "C"
import "unsafe"

func init() {
	C.init_dt()
}

type CallTime struct {
	Name   string
	Millis int
}

type ServiceTime struct {
	CallTime
	ChildCalls []CallTime
}

type ServiceTimer struct {
	ServiceTime
}

//WriteTimingData composes a C structure to pass to the C code above that uses the
//Dynatrace ADK C macros. Since the C macros rely on thread local storage, the work they do
//needs to happen on the same thread: if general calls are made using the Dynatrace ADK
//in a golang program, the pure path will become corrupted anytime the go scheduler suspends
//a thread, as the thread context gets overwritten. This is because the golang scheduler multiplexs
//logical threads onto physical threads.
func WriteTimingData(svcTiming ServiceTimer) {

	var td C.TimingData

	//Timing context for outer measurement point (Xavi entry/exit timing)
	td.name = C.CString(svcTiming.Name)
	defer C.free(unsafe.Pointer(td.name))
	td.timing = C.int(svcTiming.Millis)

	//Timing context for inner context (proxied service call). TODO: handle multiple inner timing points.
	td.subname = C.CString(svcTiming.ChildCalls[0].Name)
	defer C.free(unsafe.Pointer(td.subname))
	td.subtiming = C.int(svcTiming.ChildCalls[0].Millis)

	C.handleTiming(td)
}
