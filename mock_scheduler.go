package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.comcast.com/viper-cog/rio-dashboard/comdef"
	"github.comcast.com/viper-cog/rio-dashboard/rio-mockutil/httputil/httpctxserver"
	"github.comcast.com/viper-cog/rio-dashboard/rio-mockutil/httputil/httputilfuncs"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	myAddr      = "localhost:9003"
	rmURL       = ""
	timeOut     = 10 * time.Second
	recordURL   = "http://localhost:9003" + comdef.URLRecordBridge
	verifyURL   = "http://localhost:9003" + comdef.URLVerify
	updateURL   = "http://localhost:9003" + comdef.URLUpdate
	upCount     = 0
	updtDataGlb = make([]comdef.UpdateReq, 200)
	recDataGlb  = make([]comdef.RecordData, 200)
	wrCount     = 0
	recordingID = ""
	numRec      = 0
	nowTime     = time.Now()
	recEntries  = 0
	rm_host     = ""
)

func writeRecordData(r comdef.RecordReq) {

	var datas []comdef.RecordData
	datas = r.Entries

	for k := range datas {
		recDataGlb[wrCount].RecordingID = datas[k].RecordingID
		recDataGlb[wrCount].StartTime = datas[k].StartTime
		recDataGlb[wrCount].EndTime = datas[k].EndTime
		recDataGlb[wrCount].Action = datas[k].Action
		recDataGlb[wrCount].AccountID = datas[k].AccountID
		recDataGlb[wrCount].StreamID = datas[k].StreamID
		recDataGlb[wrCount].AdZone = datas[k].AdZone
		recDataGlb[wrCount].StationID = datas[k].StationID
		recDataGlb[wrCount].ListingID = datas[k].ListingID
		recDataGlb[wrCount].ScheduledTime = datas[k].ScheduledTime
		wrCount = wrCount + 1
	}
}

//writeUpdateContent to write the update request to a
func writeUpdateContent(reqPut *comdef.UpdateReq) {

	var datas []comdef.StatusObject
	datas = reqPut.Entries
	slice := make([]comdef.StatusObject, len(datas))
	entFlag := false
	fmt.Println(reqPut)
	for k := range datas {
		for i := 0; i < upCount; i++ {
			updtEntries := updtDataGlb[i].Entries

			for j := range updtEntries {
				if updtEntries[j].RecordingID == datas[k].RecordingID && updtEntries[j].Status == "STARTED" {
					updtEntries[j].Status = datas[k].Status
					updtEntries[k].SegmentCount = datas[k].SegmentCount
					entFlag = true
				} else {
					entFlag = false
				}
			}
		}
		if datas[k].Status == "STARTED" && entFlag == false {
			slice[k].RecordingID = datas[k].RecordingID
			slice[k].Status = datas[k].Status
			slice[k].ScheduledTime = datas[k].ScheduledTime //hard coded the time since there is no ScheduledTime property in request
			slice[k].SegmentCount = datas[k].SegmentCount

		}
	}
	if entFlag == false {
		updtDataGlb[upCount].RequestID = reqPut.RequestID
		updtDataGlb[upCount].Entries = slice
		upCount = upCount + 1

	}

}

//Process subscribe request received from Recorder Manager
func processSubscribe(input httpctxserver.HandlerProcessInput) (bool, int, error) {
	var (
		req   *comdef.SubscribeReq
		resp  *comdef.SubscribeResp
		found bool
	)

	//Sanity check to make sure req and resp data type are passed in correctly
	if req, found = input.Req.(*comdef.SubscribeReq); found != true || req == nil {
		err := errors.New("Unsupported req type for Subscribe")
		log.Fatal(err)
		return false, http.StatusInternalServerError, err
	}

	if resp, found = input.Resp.(*comdef.SubscribeResp); found != true || resp == nil {
		err := errors.New("Unsupported resp type for Subscribe")
		log.Fatal(err)
		return false, http.StatusInternalServerError, err
	}

	//set the record URL for application
	rmURL = req.RecordURL

	//generate fake response data
	genSubResp(req, resp)

	return true, http.StatusOK, nil
}

//Process Verify request from Recorder Manager
func processVerify(input httpctxserver.HandlerProcessInput) (bool, int, error) {
	var (
		req   *comdef.VerifyReq
		resp  *comdef.VerifyResp
		found bool
	)

	//Sanity check to make sure req and resp data type are passed in correctly
	if req, found = input.Req.(*comdef.VerifyReq); found != true || req == nil {
		err := errors.New("Unsupported req type for Verify")
		log.Fatal(err)
		return false, http.StatusInternalServerError, err
	}

	if resp, found = input.Resp.(*comdef.VerifyResp); found != true || resp == nil {
		err := errors.New("Unsupported resp type for Verify")
		log.Fatal(err)
		return false, http.StatusInternalServerError, err
	}

	//generate fake response data
	genVerifyResp(req, resp)

	return true, http.StatusOK, nil
}

//Process Update request from Recorder Manager
func processUpdate(input httpctxserver.HandlerProcessInput) (bool, int, error) {
	var (
		req   *comdef.UpdateReq
		resp  *comdef.UpdateResp
		found bool
	)

	//Sanity check to make sure req and resp data type are passed in correctly
	if req, found = input.Req.(*comdef.UpdateReq); found != true || req == nil {
		err := errors.New("Unsupported req type for Update")
		log.Fatal(err)
		return false, http.StatusInternalServerError, err
	}
	if resp, found = input.Resp.(*comdef.UpdateResp); found != true || resp == nil {
		err := errors.New("Unsupported resp type for Update")
		log.Fatal(err)
		return false, http.StatusInternalServerError, err
	}
	genUpdateResp(req, resp)
	writeUpdateContent(req)

	//For now we simply return everything is successfully processed
	return false, http.StatusOK, nil
}

//Process Record Request from external source.  It will forward the request to Recorder Manager
//and forward the response received from RM to the external source.
func processBridgeRecord(input httpctxserver.HandlerProcessInput) (bool, int, error) {

	var (
		req   *comdef.RecordReq
		resp  *comdef.RecordResp
		found bool
	)

	//Sanity check to make sure req and resp data type are passed in correctly
	if req, found = input.Req.(*comdef.RecordReq); found != true || req == nil {
		err := errors.New("Unsupported req type for Bridge Record")
		log.Fatal(err)
		return false, http.StatusInternalServerError, err
	}

	if resp, found = input.Resp.(*comdef.RecordResp); found != true || resp == nil {
		err := errors.New("Unsupported resp type for Bridge Record")
		log.Fatal(err)
		return false, http.StatusInternalServerError, err
	}

	//sanity check if we have a URL for RM
	if len(rmURL) == 0 {
		log.Printf("Cannot Forward request to Recorder Manager, no RM URL available")
		err := errors.New("No RM URL is aviable")
		return false, http.StatusBadGateway, err
	}

	//fake request id
	req.RequestID = fmt.Sprintf("xxxxxxxx-aaaa-bbbb-cccc-%v", time.Now().UnixNano)

	//Send request to RM
	httpStatus, err := httputilfuncs.PostHTTPJSONString(rmURL, *req, resp, timeOut)
	if err != nil {
		log.Printf("Error when sending Update Request, err: %v", err)
		return false, http.StatusBadGateway, nil

	}

	//forward the RM response back to external source
	log.Printf("HttpStatus: %d, Resp: %v", httpStatus, *resp)
	if httpStatus == http.StatusOK {
		return true, httpStatus, nil
	}

	return false, httpStatus, nil

}

//GetCurrenttime
func genCurentTime() string {

	t := nowTime.Format("2006-01-02T15:04:05.000Z07:00")

	return t
}

//genSubResp generate subscribe response
func genSubResp(reqPut *comdef.SubscribeReq, respPut *comdef.SubscribeResp) {

	respPut.RequestID = reqPut.RequestID
	respPut.MaxEntries = 1

	verifyURL := "http://" + myAddr + comdef.URLVerify
	respPut.VerifyURL = verifyURL

	//Assigning the value to RM host
	rm_host = strings.Replace(rmURL, "http://", "", 1)
	rm_host = strings.Replace(rm_host, ":9004/a8/record", "", 1)

}

func genUpdateResp(reqPut *comdef.UpdateReq, respPut *comdef.UpdateResp) {

	var datas []comdef.StatusObject
	datas = reqPut.Entries
	slice := make([]comdef.UpdateRespEntries, len(datas))
	for i := range datas {
		//respPut.Entries = append(respPut.Entries)
		slice[i].RecordingID = datas[i].RecordingID
		slice[i].Status = 0
	}
	
	respPut.RequestID = reqPut.RequestID
	respPut.Entries = slice

}

//genSubResp generate verify response
func genVerifyResp(reqPut *comdef.VerifyReq, respPut *comdef.VerifyResp) {

	var datas []comdef.VerifyReqData
	datas = reqPut.Entries
	slice := make([]comdef.StatusObject, len(datas))
	for i := range datas {
		respPut.Entries = append(respPut.Entries)
		slice[i].RecordingID = datas[i].RecordingID
		slice[i].Status = "UNKNOWN"
		slice[i].ScheduledTime = "2015-04-14T17:14:45.920181262Z" //hard coded the time since there is no ScheduledTime property in request
		slice[i].SegmentCount = 0
	}
	respPut.RequestID = reqPut.RequestID
	respPut.Entries = slice

}

//handleSubscribe to handle subscribe request from RM
func handleSubscribe(w http.ResponseWriter, r *http.Request) error {

	param := httpctxserver.HandlerProcessInput{
		Req:     &comdef.SubscribeReq{},
		Resp:    &comdef.SubscribeResp{},
		UsrData: nil}

	return httpctxserver.DoJSONReq(w, r, processSubscribe, param)

}

//handleVerify to handle verify request form RM and client
func handleVerify(w http.ResponseWriter, r *http.Request) error {

	param := httpctxserver.HandlerProcessInput{
		Req:     &comdef.VerifyReq{},
		Resp:    &comdef.VerifyResp{},
		UsrData: nil}

	return httpctxserver.DoJSONReq(w, r, processVerify, param)

}

//handleUpadte to handle Update request from RM
func handleUpdate(w http.ResponseWriter, r *http.Request) error {

	param := httpctxserver.HandlerProcessInput{
		Req:     &comdef.UpdateReq{},
		Resp:    &comdef.UpdateResp{},
		UsrData: nil}

	return httpctxserver.DoJSONReq(w, r, processUpdate, param)

}

//handleRecordBridge to handle Record request from client
func handleRecordBridge(w http.ResponseWriter, r *http.Request) error {

	param := httpctxserver.HandlerProcessInput{
		Req:     &comdef.RecordReq{},
		Resp:    &comdef.RecordResp{},
		UsrData: nil}

	return httpctxserver.DoJSONReq(w, r, processBridgeRecord, param)
}

//handleRecord to handle record request from client and forward it to the mock RM
func handleRecord(w http.ResponseWriter, r *http.Request) error {

	var msgString []byte

	nowTime = time.Now()

	_, numRec := getrecPayload(r)
	if numRec != 0 {
		for k := 0; k < numRec; k++ {
			recEntries = recEntries + 1
			req, _ := getrecPayload(r)
			httpStatus, err, resp := postRecord(req, r)
			//go runtime.GC()

			if err != nil {
				log.Printf("Failed on send Record request! err:%v\n", err)

			}
			if httpStatus == 204 {
				msgString = []byte(strconv.Itoa(numRec) + " Recordings Submitted sucessfully to RM ")

			} else if httpStatus == 200 {
				msgString = []byte("Error code " + strconv.Itoa(resp.Entries[0].Status) + ":" + resp.Entries[0].Error + " in Record Request")

			} else if httpStatus != 200 {
				//p := &Page{Title: "RESPONSE from RM", Body: []byte("Recording submitted sucessfully to RM")}
				msgString = []byte("RM not available please retry after Some Time, ERROR: " + strconv.Itoa(httpStatus))
			}

		}
	} else {
		msgString = []byte("Enter a valid number of copies to be created")
	}

	p := &comdef.Page{Body: msgString, RMUrl: rm_host}
	pageTmpl := "./html/resp_record.html"
	renderTemplate(w, pageTmpl, p)

	return nil

}

func postRecord(reqData comdef.RecordReq, r *http.Request) (httpStatus int, err error, resp *comdef.RecordResp) {

	//resp := comdef.RecordResp{}
	httpStatus, err = httputilfuncs.PostHTTPJSONString(rmURL,
		reqData, &resp, 10*time.Second)

	if err != nil {

		log.Printf("Failed on send Record request! err:%v\n", err)

	}
	if httpStatus == 204 {
		log.Printf("Recording Submitted sucessfully to RM")

	}
	if httpStatus == 200 {
		log.Printf("One of the fields is missing in the request")
	}

	return httpStatus, err, resp

}

//handleVerResp to handle verify request from client and forward it to the mock RM
func handleVerResp(w http.ResponseWriter, r *http.Request) error {

	req := getverifyPayload(r)

	resp := comdef.VerifyResp{}
	httpStatus, err := httputilfuncs.PostHTTPJSONString(verifyURL,
		req, &resp, 10*time.Second)

	if err != nil {
		log.Printf("Failed on send Record request! err:%v\n", err)
		return err
	}

	log.Printf("Http Status: %d, resp: %v", httpStatus, resp)
	if httpStatus == 200 {
		w.Write([]byte("Recording verified sucessfully:\n"))
		b, _ := json.Marshal(resp)
		w.Write(b)

	} else {
		w.Write([]byte("Recording not verified sucessfully please resubmit"))
	}

	return err

}

func handleUpdateClient(w http.ResponseWriter, r *http.Request) error {
	var retError error

	//var flag = false
	var recData []comdef.UpdateReq
	r.ParseForm()
	reqID := r.PostFormValue("UpdateReqID")
	updateOption := r.PostFormValue("SelectAll")

	if updateOption == "SelectAll" {
		for j := range updtDataGlb {
			if j <= (upCount - 1) {
				if updtDataGlb[j].RequestID != "" {

					recData = append(recData, updtDataGlb[j])
				} else {
					break
				}

			}
		}
	} else {
		for j := range updtDataGlb {
			if j <= (upCount - 1) {
				if reqID == updtDataGlb[j].RequestID {
					recData = append(recData, updtDataGlb[j])
				}
			}
		}

	}

	p := &comdef.Page{Title: "UPDATE STATUS", DataUpdate: recData, RMUrl: rm_host} //[]byte(byteData)}
	pageTmpl := "./html/resp_update.html"
	renderTemplate(w, pageTmpl, p)

	return retError
}

func handleDelete(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()
	delRecid := r.PostFormValue("RecordingID")
	reqData := getdelPayload(delRecid)
	httpStatus, err, _ := postRecord(reqData, r)
	

	if err != nil {
		log.Printf("Failed on send Record request! err:%v\n", err)

	}
	if httpStatus == 204 {
		log.Printf("Delete request Submitted Sucessfully to RM")
		for i := 0; i < upCount; i++ {
			updtEntries := updtDataGlb[i].Entries

			for j := range updtEntries {
				if updtEntries[j].RecordingID == delRecid && updtEntries[j].Status == "COMPLETE" {
					updtEntries[j].Status = "DELETED"

				}
			}
		}

	} else if httpStatus == 500 {
		log.Printf("RM not available please retry, ERROR: " + strconv.Itoa(httpStatus))

	}

	return nil

}
func handleStop(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()

	stpRecid := r.PostFormValue("RecordingID")
	
	reqData := getstopPayload(stpRecid)
	httpStatus, err, _ := postRecord(reqData, r)
	

	if err != nil {
		log.Printf("Failed on send Record request! err:%v\n", err)

	}
	if httpStatus == 204 {
		log.Printf("Stop request Submitted Sucessfully to RM")
		for i := 0; i < upCount; i++ {
			updtEntries := updtDataGlb[i].Entries

			for j := range updtEntries {
				if updtEntries[j].RecordingID == stpRecid && updtEntries[j].Status == "COMPLETE" {
					updtEntries[j].Status = "STOPPED"

				}
			}
		}

	} else if httpStatus == 500 {
		log.Printf("RM not available please retry, ERROR: " + strconv.Itoa(httpStatus))

	}

	return nil

}

//IndexHandler for handling the / request
func IndexHandler(w http.ResponseWriter, r *http.Request) error {

	var retData []comdef.StatusObject
	var rData []comdef.RecordData

	for itrCount := 0; itrCount <= upCount; itrCount++ {
		datas := updtDataGlb[itrCount].Entries
		for j := range datas {
			if updtDataGlb[itrCount].Entries[j].RecordingID != "" {
				if len(retData) != 0 {
					var retFlag bool
					for k := range retData {
						if retData[k].RecordingID == updtDataGlb[itrCount].Entries[j].RecordingID {
							retFlag = false
						} else {
							retFlag = true
						}
					}
					if retFlag == true {
						retData = append(retData, updtDataGlb[itrCount].Entries[j])

					}
				} else {
					retData = append(retData, updtDataGlb[itrCount].Entries[j])

				}

			} else {
				break
			}
		}
	}
	for j := range recDataGlb {

		if j <= wrCount {
			if recDataGlb[j].RecordingID != "" {
				rData = append(rData, recDataGlb[j])
			} else {

				break
			}

		}
	}

	p := &comdef.Page{Title: "RIO RECORDINGS", DataStatus: retData, DataRecord: rData, RMUrl: rm_host}

	pageTmpl := "./html/index.html"
	renderTemplate(w, pageTmpl, p)

	return nil
}

//FileHandler for handling the file(image/html) request
func FileHandler(w http.ResponseWriter, r *http.Request) error {

	http.ServeFile(w, r, r.URL.Path[1:])

	return nil

}
func getstopPayload(stpRecid string) comdef.RecordReq {

	var stpData comdef.RecordReq
	slice := make([]comdef.RecordData, 1)
	nwTime     := time.Now()
	for i := range recDataGlb {
		if stpRecid == recDataGlb[i].RecordingID {
			slice[0].Action = "SET"
			slice[0].RecordingID = stpRecid
			slice[0].ScheduledTime = recDataGlb[i].ScheduledTime
			slice[0].AccountID = recDataGlb[i].AccountID
			slice[0].StartTime = recDataGlb[i].StartTime //nowTime.Add(time.Duration(1*30) * time.Second).Format(comdef.TimeFrmtStr) //r.PostFormValue("StartTime")
			slice[0].EndTime = nwTime.Format(comdef.TimeFrmtStr)
			slice[0].StreamID = recDataGlb[i].StreamID
			slice[0].AdZone = recDataGlb[i].AdZone
			slice[0].StationID = recDataGlb[i].StationID
			slice[0].UpdateURL = updateURL //hard coded the scheduler to localhost
			slice[0].ListingID = recDataGlb[i].ListingID

		}
	}
	stpData.RequestID = fmt.Sprintf("xxxxxxxx-aaaa-bbbb-cccc-%v", time.Now().UnixNano)
	stpData.Entries = slice
	return stpData
}

//getrecPayload to generate the recording payload in JSON using the post value from the form
func getrecPayload(r *http.Request) (comdef.RecordReq, int) {

	var recData comdef.RecordReq

	r.ParseForm()
	recordingID = fmt.Sprintf("XRID-%d-%d", nowTime.Unix(), recEntries)
	ScheduledTime := genCurentTime()

	slice := make([]comdef.RecordData, 1) //Created a slice only for one entry

	slice[0].RecordingID = recordingID
	slice[0].Action = r.PostFormValue("Action")
	slice[0].ScheduledTime = ScheduledTime
	slice[0].AccountID = r.PostFormValue("AccountID")
	slice[0].StartTime = r.PostFormValue("StartTime") //nowTime.Add(time.Duration(1*30) * time.Second).Format(comdef.TimeFrmtStr) //r.PostFormValue("StartTime")
	slice[0].EndTime = r.PostFormValue("EndTime")     //nowTime.Add(time.Duration((2)*30) * time.Second).Format(comdef.TimeFrmtStr)
	slice[0].StreamID = r.PostFormValue("StreamID")
	slice[0].AdZone = r.PostFormValue("azone")
	slice[0].StationID = r.PostFormValue("StationID")
	slice[0].UpdateURL = updateURL //hard coded the scheduler to localhost
	slice[0].ListingID = r.PostFormValue("ListingID")

	numRec, _ = strconv.Atoi(r.PostFormValue("NoOfRec")) //To get the no of recordings to be scheduled


	recData.RequestID = fmt.Sprintf("xxxxxxxx-aaaa-bbbb-cccc-%v", time.Now().UnixNano)
	recData.Entries = slice
	writeRecordData(recData)
	return recData, numRec
}

func getdelPayload(delRecid string) comdef.RecordReq {
	var delData comdef.RecordReq
	slice := make([]comdef.RecordData, 1)

	for i := range recDataGlb {
		if delRecid == recDataGlb[i].RecordingID {
			slice[0].Action = "REMOVE"
			slice[0].RecordingID = delRecid
			slice[0].ScheduledTime = recDataGlb[i].ScheduledTime
			slice[0].AccountID = recDataGlb[i].AccountID

		}
	}
	delData.RequestID = fmt.Sprintf("xxxxxxxx-aaaa-bbbb-cccc-%v", time.Now().UnixNano)
	delData.Entries = slice
	return delData
}

//getverifyPayload to generate the verify payload in JSON using the post value from the form
func getverifyPayload(r *http.Request) comdef.VerifyReq {

	var verifyData comdef.VerifyReq
	r.ParseForm()

	slice := make([]comdef.VerifyReqData, 1) //Created a slice only for one entry
	slice[0].RecordingID = r.PostFormValue("VerifyReqID")
	verifyData.RequestID = r.PostFormValue("VerifyRecID")
	verifyData.Entries = slice

	return verifyData
}

//getupdatePayload to generate the update payload in JSON using the post value from the form
func getupdatePayload(r *http.Request) comdef.UpdateReqClient {

	var updateData comdef.UpdateReqClient
	r.ParseForm()

	updateData.RecordingID = r.PostFormValue("UpdateRecID")

	return updateData
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *comdef.Page) {
	t, err := template.ParseFiles(tmpl)

	t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func main() {

	var versionStr string

	addrStr := flag.String("addr", ":9003", "Address to listen on")
	flag.Parse()

	myAddr = *addrStr

	out, err := json.MarshalIndent(Version, "", "  ")
	if err != nil {
		log.Printf("Failed to automatically create version string, use default one.  Err:%v", err)

		versionStr = "{\n" +
			"\"BuildTime\": \"" + Version.BuildTime + "\",\n" +
			"\"BuildUser\": \"" + Version.BuildUser + "\",\n" +
			"\"GitTag\": \"" + Version.GitTag + "\",\n" +
			"\"GitBranch\": \"" + Version.GitBranch + "\"\n}"
	} else {
		versionStr = string(out)
	}

	serverSetup := httpctxserver.HTTPServerSetup{AddrStr: myAddr, VersionStr: versionStr}

	mockScheduler := httpctxserver.NewHTTPServer(serverSetup)
	//Add handler
	mockScheduler.AddHandler(comdef.URLSubscribe, handleSubscribe)
	mockScheduler.AddHandler(comdef.URLVerify, handleVerify)
	mockScheduler.AddHandler(comdef.URLUpdate, handleUpdate)
	mockScheduler.AddHandler(comdef.URLRecordBridge, handleRecordBridge)
	mockScheduler.AddHandler(comdef.FileHandle, FileHandler)
	mockScheduler.AddHandler(comdef.IndexHandle, IndexHandler)
	mockScheduler.AddHandler(comdef.RecordHandle, handleRecord)
	mockScheduler.AddHandler(comdef.VerifyHandle, handleVerResp)
	//mockScheduler.AddHandler("/updateClient", handleUpdateClient)
	mockScheduler.AddHandler(comdef.DelHandle, handleDelete)
	mockScheduler.AddHandler(comdef.StopHandle, handleStop)

	mockScheduler.Start()

}
