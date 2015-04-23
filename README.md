This is a simple sample that shows how to integrate a golang program
with Dynatrace using the Dynatrace ADK.

Running the sample requires several environment variables to be set:

<pre>
DT_AGENTLIBRARY=/home/vagrant/agent/dynatrace-6.1.0/agent/lib64/libdtagent.so
DT_AGENTNAME=NativeApplication_Monitoring
DT_DEBUGADK=true
DT_SERVER=localhost:9998
</pre>

The ADK also has a library that needs to be in the LD_LIBRARY_PATH... apparently it is not as special as the agent
library which requires a special environment variable to load it from.

<pre>
LD_LIBRARY_PATH=/home/vagrant/dynatrace-6.1.0/adk/lib64
</pre>

