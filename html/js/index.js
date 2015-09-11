
function newPopup(url) {
        popupWindow = window.open(
        url,'popUpWindow','height=500,width=1003,left=10,top=10,resizable=no,scrollbars=yes,toolbar=yes,menubar=no,location=no,directories=no,status=yes,fullscreen=yes');
    }
    
function setRowValue(el){
		var recID=$(el).text()
        sessionStorage.setItem("RecordingID", recID);
    }

    
$('#selectAll').click(function (e) {
    $(this).closest('table').find('td input:checkbox').prop('checked', this.checked);
});

function funcSort(){
	var TableData = new Array();
	var sortedData = new Array();
	var i=0;
	$('#recordingtbl tr').each(function(row, tr){
		var row = $(this);
    		TableData[i]={
        		"RecordingId" :$(tr).find('td:eq(1)').text()
			, "ScheduledTime" : $(tr).find('td:eq(2)').text()
        		, "SegmentCount" : $(tr).find('td:eq(3)').text()
			,"Status":$(tr).find('td:eq(4)').text()
				
		}
		sortedData[i] =$(tr).find('td:eq(1)').text();
		i=i+1
	}); 
	
	sortedData.shift();
	sortedData = sortedData.sort()
	TableData.shift();
	
	j=0
	$('#recordingtbl tr').each(function(row, tr){
		
		if(($(tr).find('td:eq(1)').text())!=""){
			for(k=0;k<TableData.length;k++){
			
				if (sortedData[j] == TableData[k].RecordingId){
					$(tr).find('a[id="RecordingID"]').text(TableData[k].RecordingId)
					$(tr).find('td:eq(2)').text(TableData[k].ScheduledTime)
					$(tr).find('td:eq(3)').text(TableData[k].SegmentCount)
					$(tr).find('td:eq(4)').text(TableData[k].Status)
				
				}
		}
		j=j+1;
		}
		
	});
	 $("#recordingtbl tr").each(function(index,cell) {

    var status = $(this).find("td").eq(4).text();
	
    		if (status == "DELETED"){
			
			$(this).find('button[id="delicon"]').attr("disabled","disabled")
			$(this).find('button[id="stpcon"]').attr("disabled","disabled")
			$(this).find('input[id="checkRow"]').attr("disabled",true)

    		}
		if (status == "COMPLETE" || status == "STARTED"){
			
			$(this).find('button[id="delicon"]').attr("disabled",false)
			$(this).find('button[id="stpcon"]').attr("disabled","disabled")
			$(this).find('input[id="checkRow"]').attr("disabled",false)

    		}
		if (status == "STARTED"){
			
			$(this).find('button[id="delicon"]').attr("disabled",false)
			$(this).find('button[id="stpcon"]').attr("disabled",false)
			$(this).find('input[id="checkRow"]').attr("disabled",false)

    		}
			
    });
}
function checkDel(){
	var i=0
	$('#recordingtbl tr').each(function(row, tr){
		var row = $(this);
		if (row.find('input[type="checkbox"]').is(':checked')){
		i=i+1
		}
		
	}); 
	if (i == 0){
		document.getElementById("errormsg").innerHTML = "*Please select atleast one row to DELETE."
	}else if (i>1){
		document.getElementById("errormsg").innerHTML = "*Please select only one row to DELETE."
	
	}else{
		onDelete()
	}
}
function checkStop(){
	var i=0
	$('#recordingtbl tr').each(function(row, tr){
		var row = $(this);
		if (row.find('input[type="checkbox"]').is(':checked')){
		i=i+1
		}
		
	}); 
	
	
	if (i == 0){
		document.getElementById("errormsg").innerHTML = "*Please select atleast one row to STOP."
	}else if (i>1){
		document.getElementById("errormsg").innerHTML = "*Please select only one row to STOP."
	
	}else{
		onStop()
	}
}

function onDelete(){
	var TableData = new Array();
	var DelData = new Array();
	var StatusData = [];
	var recid 
	var statusUpdt
	var deleteBtn
	var i=0
	$('#recordingtbl tr').each(function(row, tr){
		var row = $(this);
		if (row.find('input[type="checkbox"]').is(':checked')){
			recid = $(tr).find('td:eq(1)').text()
			statusUpdt = $(tr).find('td:eq(4)')
			deleteBtn = $(this).find('button[id="delicon"]')
			chckBox = $(this).find('input[type="checkbox"]')
		
    		TableData[i]={
        		"RecordingId" :$(tr).find('td:eq(1)').text()
        		, "ScheduledTime" : $(tr).find('td:eq(2)').text()
        		, "SegmentCount" : $(tr).find('td:eq(3)').text()
		
    		}
		DelData[i]=recid
		StatusData[i]=statusUpdt
		i=i+1
		}else{
			document.getElementById("errormsg").innerHTML = "*Please select atleast one row to DELETE."
		}
		
	}); 
	
	TableData.shift();  // first row is the table header - so remove
	for(k=0;k<DelData.length;k++){
	
	$.ajax({
   		type: "POST",
    		url: "/delete",
    		data: "RecordingID="+DelData[k],
    		success: function(msg){
			document.getElementById("errormsg").innerHTML="Deleting "+DelData.length+" recordings."
    		}
	});
	}
}

function onStop(){
	var TableData = new Array();
	var recid 
	var statusUpdt
	var deleteBtn
	var StopData = new Array();
	var i=0
	$('#recordingtbl tr').each(function(row, tr){
		var row = $(this);
		if (row.find('input[type="checkbox"]').is(':checked')){
			recid = $(tr).find('td:eq(1)').text()
			statusUpdt = $(tr).find('td:eq(4)')
			deleteBtn = $(this).find('button[id="delicon"]')
			chckBox = $(this).find('input[type="checkbox"]')
		
    		TableData[row]={
        		"RecordingId" :$(tr).find('td:eq(1)').text()
        		, "ScheduledTime" : $(tr).find('td:eq(2)').text()
        		, "SegmentCount" : $(tr).find('td:eq(3)').text()
		
    		}
		StopData[i]=recid
		i=i+1
		}
	}); 
	TableData.shift();  // first row is the table header - so remove
	StopData[i]=recid
	for(k=0;k<StopData.length;k++){
 	$.ajax({
   		type: "POST",
    		url: "/stop",
    		data: "RecordingID="+StopData[k],
    		success: function(msg){
		document.getElementById("errormsg").innerHTML="Stopping "+(StopData.length-1)+" recordings."
    		}
	});
	}
}


