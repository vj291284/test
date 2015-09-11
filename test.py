import json
import urllib
data = urllib.urlencode({"RecordingID":"buddy","ContentEncoding":"MPEG4"})
u = urllib.urlopen("http://localhost:9004/rm/RecordingInfo?%s" % data)
