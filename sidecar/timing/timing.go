package timing

/*
#cgo CFLAGS: -I/home/vagrant/dynatrace-6.1.0/adk/include
#cgo LDFLAGS: -L/home/vagrant/dynatrace-6.1.0/adk/lib64 -l dtadk
#include <stdlib.h>
#include <dynatrace_adk.h>
#include <stdio.h>
#include <time.h>

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

typedef struct {
	char* name;
	int timing;
	char *subname;
	int subtiming;
} TimingData;

void subcall(int subcallTime) {
	DYNATRACE_ENTER();
	char buf[256];
	sprintf(buf, "proxied service call time: %d ms", subcallTime);
	DYNATRACE_LOG(0, buf, "xavi");
	DYNATRACE_EXIT();
}

void api(int apiTime, int subcallTime) {
	DYNATRACE_ENTER();
	char buf[256];
	sprintf(buf, "api call time: %dms", apiTime);
	DYNATRACE_LOG(0, buf, "xavi");
	subcall(subcallTime);
	DYNATRACE_EXIT();
}

static void handleTiming(TimingData td) {
	
	DYNATRACE_API("Son of Sidecar");
	DYNATRACE_START_PUREPATH();
	//printf("Name: %s, timing: %d\n", td.name, td.timing);
	//printf("Subname: %s, subtiming: %d\n", td.subname, td.subtiming);
	api(td.timing, td.subtiming);
	DYNATRACE_EXIT();
}





*/
import "C"
import "unsafe"

func init() {
	C.init_dt()
}

type CallTime struct {
	Name string
	Millis int
}

type ServiceTime struct {
	CallTime
	ChildCalls []CallTime
}

type ServiceTimer struct {
	ServiceTime
}

func WriteTimingData(svcTiming ServiceTimer) {
	var td C.TimingData
	td.name = C.CString(svcTiming.Name)
	defer C.free(unsafe.Pointer(td.name))
	td.timing = C.int(svcTiming.Millis)
	td.subname = C.CString(svcTiming.ChildCalls[0].Name)
	defer C.free(unsafe.Pointer(td.subname))
	td.subtiming = C.int(svcTiming.ChildCalls[0].Millis)

	C.handleTiming(td)
}

