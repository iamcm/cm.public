package main

import (
	"html/template"
	_ "log"
	"net/http"
)

func main() {
	//	http.HandleFunc("/do", doHandler)
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	t, _ := template.New("home").Parse(htmltemplate)
	t.Execute(rw, nil)
}

var htmltemplate string = `
<html>
<body>
<script type="text/javascript">
function removeDuplicates(arr) {
	_arr = [];

    for (var i = 0, len = arr.length; i < len; ++i) {
        if (i == 0 || _arr.indexOf(arr[i]) == -1) {
            _arr.push(arr[i]);
        }
    }
    return _arr;
}

function go(){
	var recognition = new webkitSpeechRecognition();
recognition.lang = "en-GB";
	recognition.continuous = true;
	recognition.interimResults = true;
	words = [];
	recognition.onresult = function(event) { 
		words = words.concat(event.results[0][0].transcript.split(' '))
		words = removeDuplicates(words)
	  	console.log(words) 
	  	if (words.indexOf('play') != -1 && words.indexOf('music') != -1) {
	  		recognition.stop();
	  		console.log('FOUND COMMAND play music')
	  	}
	}
	recognition.start();
}
  </script>

    <input type="button" value="Click to Speak" onclick="go()">
</body>
</html>
`
