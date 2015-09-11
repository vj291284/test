rio-mock-segment-recorder
===============================


##RIO Mock A8 Scheduler

To run the server: rio-mock-a8-scheduler -addr=<addr to listen on, i.e.:localhost:9001>  

If -addr is not specified, "localhost:9001" is the default.  

Server supports "/bridge/record" end point.  It will take a record request from external source,
fill in the request id and update URL, then forward the request to Recorder Manager.  After
receiving response back from the Recorder Manager, it will forward the response back to external source.  

The URL of Recorder Manager is set by Subscribe Request sent by Recorder Manager. 

##comdef  

comdef/ has all the structures defined with JSON annotation for interface protocol.  

##rm_example

This is a Mock Recorder Manager Example.

To run rm-example: rm-example -myaddr=<addr to listen on, i.e. "http://localhost:9002"> -a8addr==<addr of a8, i.e. "http://localhost:9001">

If -myaddr is not specified, "localhost:9002" is the default.
If -a8addr is not specified, "localhost:9001" is the default.  

When Mock Recorder Manager starts, it will send Subscribe Request to A8 scheduler.  After it receives
a Record Request, it will respond with http status code 204 and start a timer of 1 second.
 After timer expires, it will send Update and Verify message to A8 scheduler.  
  
##req_generator_example
This is a sample external source to send a record request to Mock A8 Scheduler's end point: "/bridge/record".

To run req_generator_example: req_generator_example -addr=<addr to send request to. i.e. "http://localhost:9001">

If -addr is not specified, "localhost:9001" is the default.  
  
##To run the entire workflow
1.  Start Mock A8 Scheduler.  
2.  Start rm-example.  
3.  Upon starting up, Mock RM sends a Subscribe Request to Mock A8 Scheduler.  Mock A8 Scheduler will remember the URL of Mock RM.  
4.  run req_generator_example.  A randomly generated record request is sent to Mock A8 Scheduler.  
5.  Upon receiving the record request, Mock A8 Scheduler fills in request id and update URL.  It then forwards the request to Mock RM.  It will not respond to external source until a response is received from Mock RM or error occurs.  
6.  Upon receiving the record request from Mock A8 Scheduler,  Mock RM responds with http status code 204.  
7.  Mock RM will also start a timer.  When timer expires, it will send Update and Verify request to Mock A8 Scheduler. 
  

