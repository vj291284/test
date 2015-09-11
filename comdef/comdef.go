package comdef

//URLSubscribe - A8 subscribe end point
//URLUpdate - A8 update end point
//URLVerify - A8 verify end point
const (
	URLSubscribe    = "/a8/subscribe"
	URLUpdate       = "/a8/update"
	URLVerify       = "/a8/verify"
	URLRecordBridge = "/bridge/record"
	TimeFrmtStr     = "2006-01-02T15:04:05.000Z07:00"
	FileHandle      = "/html/"
	RecordHandle    = "/record"
	VerifyHandle    = "/verify"
	DelHandle       = "/delete"
	StopHandle      = "/stop"
	IndexHandle     = "/"
)

//SubscribeReq request payload /a8/subscribe:request
type SubscribeReq struct {
	RequestID         string //`json:"RequestID"`         //identifier for request
	RecordingLocality string //`json:"RecordingLocality"` //identifier for grouping of recordings
	MaxEntries        int    //`json:"MaxEntries"`        //Max no of recordings that can be included in any record or verify message
	RecordURL         string //`json:"RecordURL"`         //Url used by scheduler to send record request to RM
	VerifyURL         string //`json:"VerifyURL"`         //Url used by sccheduler to send verify request to RM
}

//SubscribeResp response payload /a8/subscribe:response
type SubscribeResp struct {
	RequestID  string //`json:"RequestID"`  //identifier for request
	MaxEntries int    //`json:"MaxEntries"` //Max no of recordings that can be included in any record or verify message
	VerifyURL  string //`json:"VerifyURL"`  //Url used by RM to send verify request to scheduler
}

//RecordReq JSON payload to send request for recording /rio/record:request to RM
type RecordReq struct {
	RequestID string       //`json:"RequestID"` //identifier for request
	Entries   []RecordData //`json:"Entries"`   //Struct var for record data entries
}

//RecordData Entries for recording payload
type RecordData struct {
	RecordingID   string `json:"RecordingId"` //Identifier for recording
	Action        string //Action to be performed on recording
	ScheduledTime string //Last update time of recording by scheduler
	AccountID     string //Identifier for owner of this recording
	StartTime     string `json:",omitempty"` //The time recording should be started
	EndTime       string `json:",omitempty"` //The time recording should be ended
	StreamID      string `json:",omitempty"` //StreamID for the stream to be recorded
	AdZone        string `json:",omitempty"` //AdZone identifier for associated AccountID
	StationID     string `json:",omitempty"` //Identifier for station id for the stream
	UpdateURL     string `json:",omitempty"` //URL used by RM to send update requests to scheduler
	ListingID     string `json:",omitempty"`
	Status        string `json:",omitempty"` //Optional field only used for displaying the reocrdings
}

//RecordResp response data to RecordReq
type RecordResp struct {
	RequestID string //identifier for request
	Entries   []RespObject
}

//RespObject response data to RecordReq, UpdateReq
type RespObject struct {
	Status      int
	RecordingID string `json:"RecordingId"`
	Error       string `json:",omitempty"`
}

//UpdateReq request payload received from RM  /a8/update:request
type UpdateReqClient struct {
	RecordingID string //identifier for request

}

//UpdateReq request payload received from RM  /a8/update:request
type UpdateReq struct {
	RequestID string         //identifier for request
	Entries   []StatusObject `json:",omitempty"`
}

//UpdateResp response to UpdateReq
type UpdateResp struct {
	RequestID string              //identifier for request
	Entries   []UpdateRespEntries `json:",omitempty"`
}
type UpdateRespEntries struct {
	Status      int    // Status code
	RecordingID string `json:"RecordingId"` // Identifier for recording ( XRID )
	Error       string `json:",omitempty"`  // Error Description
}

//StatusObject Update entries data for request payload from RM
type StatusObject struct {
	Status        string    `json:",omitempty"`  //Current status of recording
	RecordingID   string    `json:"RecordingId"` //Identifier for ecording
	ScheduledTime string    `json:",omitempty"`  //Last update time of recording by scheduler
	SegmentCount  int       `json:",omitempty"`  //Number of Segments in recording
	Segments      []Segment `json:",omitempty"`
	FailureCode   int       `json:",omitempty"` //Failure code to identify why recording failed
	FailureString string    `json:",omitempty"` //Failure string gives details of why recording failed
}

//Segment object, used in StatusObject
type Segment struct {
	Ordinal     int
	ActualStart string
	ActualEnd   string
}

//VerifyReq request payload received form RM /a8/verify/:request
type VerifyReq struct {
	RequestID string          //identifier for request
	Entries   []VerifyReqData //Array of recordingId
}

//VerifyReqData request data
type VerifyReqData struct {
	RecordingID string `json:"RecordingId"`
}

//VerifyResp response payload to be sent to RM  /a8/verify/:response
type VerifyResp struct {
	RequestID string //identifier for request
	Entries   []StatusObject
}

//Page Struct for page template
type Page struct {
	Title      string         `json:",omitempty"`
	Body       []byte         `json:",omitempty"`
	DataUpdate []UpdateReq    `json:",omitempty"`
	DataStatus []StatusObject `json:",omitempty"`
	DataRecord []RecordData   `json:",omitempty"`
	RMUrl      string         `json:",omitempty"`
}
