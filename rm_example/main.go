package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.comcast.com/viper-cog/rio-mock-a8-scheduler/comdef"
	"github.comcast.com/viper-cog/rio-mockutil/httputil/httpctxserver"
	"github.comcast.com/viper-cog/rio-mockutil/httputil/httputilfuncs"
)

var (
	a8VerifyURL string
	timeOut     = 10 * time.Second
	dataChan    = make(chan comdef.RecordData)
	doneChan    = make(chan byte)
)

type updateData struct {
	StatusObj comdef.StatusObject
	RecordObj comdef.RecordData
}

func sendUpdateAndVerify(dataChan <-chan comdef.RecordData, doneChan chan<- byte) {

	quit := false

	pendingUpdateList := []updateData{}
	pendingVerifyList := []updateData{}

	for !quit {

		select {

		case data := <-dataChan:
			statusObject := comdef.StatusObject{Status: "UNKNOWN",
				RecordingID:   data.RecordingID,
				ScheduledTime: data.ScheduledTime,
				SegmentCount:  0,
			}
			updataData := updateData{StatusObj: statusObject, RecordObj: data}
			pendingUpdateList = append(pendingUpdateList, updataData)

			break
		case <-time.After(time.Second * 1):
			if len(pendingVerifyList) > 0 {
				verifyReq := comdef.VerifyReq{RequestID: fmt.Sprintf("r12345678-1234-aaaa-bbbb-%d", time.Now().Unix()),
					Entries: []comdef.VerifyReqData{}}

				for _, obj := range pendingVerifyList {
					verifyData := comdef.VerifyReqData{RecordingID: obj.RecordObj.RecordingID}
					verifyReq.Entries = append(verifyReq.Entries, verifyData)
				}

				//send verify req, we assume all status Obj has the same update URL for now
				log.Printf("Sending Verify Request: req:%v", verifyReq)
				verifyResp := comdef.VerifyResp{}
				httpStatus, err := httputilfuncs.PostHTTPJSONString(a8VerifyURL, verifyReq, &verifyResp, timeOut)
				if err != nil {
					log.Printf("Error when sending Update Request, err: %v", err)
				} else {
					log.Printf("HttpStatus: %d, Resp: %v", httpStatus, verifyResp)
				}

				//clear list for gc
				pendingVerifyList = []updateData{}
			}

			if len(pendingUpdateList) > 0 {
				updateReq := comdef.UpdateReq{RequestID: fmt.Sprintf("r12345678-1234-cccc-dddd-%d", time.Now().Unix()),
					Entries: []comdef.StatusObject{}}
				for _, obj := range pendingUpdateList {
					updateReq.Entries = append(updateReq.Entries, obj.StatusObj)
				}
				//send update req, we assume all status Obj has the same update URL for now
				log.Printf("1 -- Sending Update Request: req:%v", updateReq)
				updateResp := comdef.UpdateResp{}
				fmt.Println(pendingUpdateList[0].RecordObj.UpdateURL)
				httpStatus, err := httputilfuncs.PostHTTPJSONString(pendingUpdateList[0].RecordObj.UpdateURL, updateReq, &updateResp, timeOut)
				if err != nil {
					log.Printf("Error when sending Update Request, err: %v", err)
				} else {
					log.Printf("HttpStatus: %d, Resp: %v", httpStatus, updateResp)
				}

				//clear list for gc
				pendingVerifyList = pendingUpdateList
				pendingUpdateList = []updateData{}
			}
			break
		} //end switch
	} //end for

	doneChan <- 1
}

func proessRecord(input httpctxserver.HandlerProcessInput) (bool, int, error) {

	var (
		req   *comdef.RecordReq
		resp  *comdef.RecordResp
		found bool
	)

	if req, found = input.Req.(*comdef.RecordReq); found != true || req == nil {
		err := errors.New("Unsupported req type for Record Request")
		log.Fatal(err)
		return false, http.StatusInternalServerError, err
	}

	if resp, found = input.Resp.(*comdef.RecordResp); found != true || resp == nil {
		err := errors.New("Unsupported resp type for Verify Request")
		log.Fatal(err)
		return false, http.StatusInternalServerError, err
	}

	//we will trigger to send update and verify request
	for _, data := range req.Entries {
		dataChan <- data
	}

	//for now we always consume all content
	return false, http.StatusNoContent, nil
}

func processVerify(input httpctxserver.HandlerProcessInput) (bool, int, error) {

	var (
		req   *comdef.VerifyReq
		resp  *comdef.VerifyResp
		found bool
	)

	if req, found = input.Req.(*comdef.VerifyReq); found != true || req == nil {
		err := errors.New("Unsupported req type for Verify Request")
		log.Fatal(err)
		return false, http.StatusInternalServerError, err
	}

	if resp, found = input.Resp.(*comdef.VerifyResp); found != true || resp == nil {
		err := errors.New("Unsupported resp type for Verify Request")
		log.Fatal(err)
		return false, http.StatusInternalServerError, err
	}

	resp.RequestID = req.RequestID
	for _, data := range req.Entries {
		statusObject := comdef.StatusObject{
			Status:        "UNKNOWN",
			RecordingID:   data.RecordingID,
			ScheduledTime: time.Now().UTC().Format(time.RFC3339Nano),
			SegmentCount:  0,
		}
		resp.Entries = append(resp.Entries, statusObject)
	}

	return true, http.StatusOK, nil
}

func handleRecord(w http.ResponseWriter, r *http.Request) error {

	param := httpctxserver.HandlerProcessInput{
		Req:     &comdef.RecordReq{},
		Resp:    &comdef.RecordResp{},
		UsrData: nil}

	return httpctxserver.DoJSONReq(w, r, proessRecord, param)

}

func handleVerify(w http.ResponseWriter, r *http.Request) error {

	param := httpctxserver.HandlerProcessInput{
		Req:     &comdef.VerifyReq{},
		Resp:    &comdef.VerifyResp{},
		UsrData: nil}

	return httpctxserver.DoJSONReq(w, r, processVerify, param)

}

func main() {

	const (
		URLRecord = "/a8/record"
		URLVerify = "/a8/verify"
	)

	myAddrStr := flag.String("myaddr", "localhost:9002", "Address to listen on")
	a8AddrStr := flag.String("a8addr", "localhost:9003", "Address of A8 scheduler")
	flag.Parse()

	versionStr := "Mock A8 Scheduler Test"
	serverSetup := httpctxserver.HTTPServerSetup{AddrStr: *myAddrStr, VersionStr: versionStr}

	//Subscribe to Mock A8 scheduler
	var subscribeReq = comdef.SubscribeReq{
		RequestID:         "c4646160-6448-11df-9c5e-0050569339e3",
		RecordingLocality: "mile-high",
		MaxEntries:        10,
		RecordURL:         "http://" + *myAddrStr + URLRecord,
		VerifyURL:         "http://" + *myAddrStr + URLVerify,
	}
	var subscribeResp comdef.SubscribeResp

	httpStatus, err := httputilfuncs.PostHTTPJSONString("http://"+*a8AddrStr+comdef.URLSubscribe,
		subscribeReq, &subscribeResp, timeOut)

	if err != nil {
		log.Printf("Failed on Subscribe request! err:%v\n", err)
		return
	}

	log.Printf("HttpStatus: %d, resp: %v", httpStatus, subscribeResp)

	if httpStatus != http.StatusOK {
		log.Printf("HttpStatus is not ok, exit")
		return
	}

	if subscribeResp.RequestID != subscribeReq.RequestID {
		log.Printf("Resp.RequstID: %s is not same as Req.RequestID: %s", subscribeResp.RequestID, subscribeReq.RequestID)
	}

	a8VerifyURL = subscribeResp.VerifyURL

	mockScheduler := httpctxserver.NewHTTPServer(serverSetup)

	//Add handler
	mockScheduler.AddHandler(URLRecord, handleRecord)
	mockScheduler.AddHandler(URLVerify, handleVerify)

	//Start go routine to send update and verify request
	go sendUpdateAndVerify(dataChan, doneChan)

	//Start Server
	mockScheduler.Start()

}
