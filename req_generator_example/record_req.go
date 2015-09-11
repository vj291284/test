package main

import (
	//"encoding/json"
	//"fmt"
	"math/rand"
	//"net/http"
	"log"
	"time"
	//"bytes"

	"github.comcast.com/viper-cog/rio-mock-a8-scheduler/comdef"
	"github.comcast.com/viper-cog/rio-mockutil/httputil/httputilfuncs"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
var recordURL = "http://localhost:9001" + comdef.URLRecordBridge

func genCurentTime() (string, string) {
	startTime := time.Now().UTC().Format(time.RFC3339Nano)
	endTime := time.Now().Add(3600000000000).UTC().Format(time.RFC3339Nano) //Adding 1 hour to the current time

	return startTime, endTime
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getrecPayload() comdef.RecordReq {

	var recData comdef.RecordReq
	RecordingID := randSeq(10) + ":?" + randSeq(17)
	StartTime, EndTime := genCurentTime()
	ScheduledTime := StartTime

	slice := make([]comdef.RecordData, 1) //Created a slice only for one entry
	slice[0].RecordingID = RecordingID
	slice[0].ScheduledTime = ScheduledTime
	slice[0].AccountID = "4614481:?1df9cZXE505693U9e4"
	slice[0].StartTime = StartTime
	slice[0].EndTime = EndTime
	slice[0].StreamID = "6064481:?1df9cZXE50569339R4"
	slice[0].AdZone = "4616061:?1df9cZXE505693W9e4"
	slice[0].StationID = "4616064481df9cZXE505Q9339e4"
	slice[0].UpdateURL = "http://localhost:9001/a8/update" //hard coded the scheduler to localhost
	slice[0].ListingID = "1646061:?1cZdf9XE5056W993e4"

	recData.RequestID = ""
	recData.Entries = slice

	return recData
}

func recordReq() {

	req := getrecPayload()
	resp := comdef.RecordResp{}
	httpStatus, err := httputilfuncs.PostHTTPJSONString(recordURL,
		req, &resp, 10*time.Second)

	if err != nil {
		log.Printf("Failed on send Record request! err:%v\n", err)
		return
	}

	log.Printf("Http Status: %d, resp: %v", httpStatus, resp)

}

func main() {
	recordReq()
}
