package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.comcast.com/viper-cog/rio-mock-a8-scheduler/comdef"
	"github.comcast.com/viper-cog/rio-mockutil/httputil/httpctxserver"
)

func doTestPost(t *testing.T, url string, handlerToBeTested httpctxserver.HandleFuncWithReturn,
	reqStr string, expRespStr string, exphttpStatus int, expError bool) bool {

	//add request:
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(reqStr)))
	if req == nil || err != nil {
		t.Errorf("Test failed on NewRequest...err: %v", err)
		return false
	}

	w := httptest.NewRecorder()
	err = handlerToBeTested(w, req)
	if w.Code != exphttpStatus {
		t.Errorf("Test failed on handler function, http status code:%d", w.Code)
		return false
	}

	if expError == true && err == nil {
		t.Errorf("Test expect error occured, but there is no error")
		return false
	}

	if expRespStr != "" {
		body, err := ioutil.ReadAll(w.Body)
		if err != nil || string(body) != expRespStr {
			t.Logf("\nGot:    %s\nExpect: %s", string(body), expRespStr)
			t.Errorf("Test failed on handler function, json string does not match")
			return false
		}
	}

	return true

}

func TestPostVerify(t *testing.T) {

	reqJSONStr := `{"RequestID" : "c4646160-6448-11df-9c5e-0059339e9056","Entries" :[{"RecordingId" : "4616064481:?1df9cZXE50569339e4"},{"RecordingId" : "4616064481:?1df9cZXE50569339e5"}]}`
	expectedRespJSONStr := `{"RequestID":"c4646160-6448-11df-9c5e-0059339e9056","Entries":[{"Status":"UNKNOWN","RecordingId":"4616064481:?1df9cZXE50569339e4","ScheduledTime":"2015-04-14T17:14:45.920181262Z","SegmentCount":0},{"Status":"UNKNOWN","RecordingId":"4616064481:?1df9cZXE50569339e5","ScheduledTime":"2015-04-14T17:14:45.920181262Z","SegmentCount":0}]}`

	if doTestPost(t, comdef.URLVerify, handleVerify, reqJSONStr, expectedRespJSONStr, http.StatusOK, false) != true {
		t.Errorf("Test PostVerify failed")
		return
	}

}

func TestPostSubscribe(t *testing.T) {

	reqJSONStr := `{"RequestID" : "c4646160-6448-11df-9c5e-0059339e9056","RecordingLocality" : "mile-high","MaxEntries" : 10,"RecordUrl" : "http://rm.cdvr.comcast.net/a8/record","VerifyUrl" : "http://rm.cdvr.comcast.net/a8/verify"}`
	expectedRespJSONStr := `{"RequestID":"c4646160-6448-11df-9c5e-0059339e9056","MaxEntries":10,"VerifyURL":"http://localhost:9001/a8/verify"}`

	if doTestPost(t, comdef.URLVerify, handleSubscribe, reqJSONStr, expectedRespJSONStr, http.StatusOK, false) != true {
		t.Errorf("Test PostSubscribe failed")
		return
	}

}

func TestPostUpdate(t *testing.T) {

	reqJSONStr := `{
 	"RequestId" : "c4646160-6448-11df-9c5e-0050569339e6",
 	"Entries" : [{
 	"Status" : "FAILED",
 	"RecordingId" : "4616064481:?1df9cZXE50569339e4",
 	"ScheduledTime" : "2014-01-18T18:33:01.324Z",
 	"SegmentCount" : 0,
 	"FailureCode" : 528451,
 	"FailureString" : "Insufficient Bandwidth"},

 	{"Status" : "COMPLETE",
 	"RecordingId" : "4616064481:?1df9cZXE50569339e5",
 	"ScheduledTime" : "2014-01-18T18:33:01.324Z",
 	"SegmentCount" : 1,
 	"Segments" : [{
 	"Ordinal" : 712,
 	"ActualStart" : "2014-01-18T18:30:00.527Z",
 	"ActualEnd" : "2014-01-18T19:30:00.140Z"}]}]

	}`

	expectedRespJSONStr := ""

	if doTestPost(t, comdef.URLVerify, handleUpdate, reqJSONStr, expectedRespJSONStr, http.StatusNoContent, false) != true {
		t.Errorf("Test PostUpdate failed")
		return
	}

}
