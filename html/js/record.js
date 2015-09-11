var mins = 0
var hrs = 0
var secs = 0
var now = new Date()
var currtHrs = now.getHours()
var currtMins = now.getMinutes()
var currtSecs = now.getSeconds()

function submitForm(value) {
	

    if (value == "recordForm"){
        setDateFields();
    }        
}

function setFields(){
	var rmURL = sessionStorage.getItem("rmURL");
	console.log(rmURL)
    document.getElementById("StartTime").value=formatLocalDate("StartTime");
    document.getElementById("EndTime").value=formatLocalDate("EndTime");
	document.getElementById("rmurl").textContent=rmURL;

}


function changeTime(){
	var strtdiff = document.getElementById("StartTmRand").value
	

	if (strtdiff >= 60){
		mins=Math.floor(strtdiff/60)
		secs=strtdiff%60
	}
	
	
	secs = currtSecs+secs
	mins = currtMins+mins
	hrs = currtHrs+hrs
	
	
	console.log("time")
	console.log(currtHrs+":"+currtMins+":"+currtSecs)
	console.log(hrs+":"+mins+":"+secs)
	
}

function formatLocalDate(value) {
    var now = new Date(),
        tzo = -now.getTimezoneOffset(),
        dif = tzo >= 0 ? '+' : '-',
        pad = function(num) {
            var norm = Math.abs(Math.floor(num));
            return (norm < 10 ? '0' : '') + norm;
        };
	if (value == "EndTime"){
		console.log(value)
		 return now.getFullYear() 
        + '-' + pad(now.getMonth()+1)
        + '-' + pad(now.getDate())
        + 'T' + pad(now.getHours())
        + ':' + pad(now.getMinutes()+2) 
        + ':' + pad(now.getSeconds()) 
		+ '.' + pad(now.getMilliseconds())
        + dif + pad(tzo / 60) 
        + ':' + pad(tzo % 60);
	}else{
		 return now.getFullYear() 
        + '-' + pad(now.getMonth()+1)
        + '-' + pad(now.getDate())
        + 'T' + pad(now.getHours())
        + ':' + pad(now.getMinutes()+1) 
        + ':' + pad(now.getSeconds()) 
		+ '.' + pad(now.getMilliseconds())
        + dif + pad(tzo / 60) 
        + ':' + pad(tzo % 60);
	};
   
}

var retrievedData = sessionStorage.getItem("record");
	
	var array = JSON.parse(retrievedData);
	
	var recID = sessionStorage.getItem("RecordingID"); 
	
	var table = document.getElementById("rstTable");
	for (var i in array){
		
		if (recID == array[i].RecordingId){
			
		var row = table.insertRow(1);
    	var cell1 = row.insertCell(0);
    	var cell2 = row.insertCell(1);
		var cell3 = row.insertCell(2);
		var cell4 = row.insertCell(3);
		var cell5 = row.insertCell(4);
		var cell6 = row.insertCell(5);
			
		cell1.innerHTML = array[i].RecordingId;
		cell2.innerHTML = array[i].StartTime;
		cell3.innerHTML = array[i].EndTime;
		cell4.innerHTML = array[i].AccountID;
		cell5.innerHTML = array[i].StreamID;
		cell6.innerHTML = array[i].StationID;	
		}
	}