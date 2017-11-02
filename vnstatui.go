package main

import (
	"fmt"
	"net/http"
	"os/exec"
  "os"
  "strings"
)

type view struct {
	arg, img string
}

var pagecontent = `<!DOCTYPE html>
<html>
<head>
  <title>vnstat-webdaemon</title>
  <style>
  img {
    vertical-align: top;
    margin-bottom: 0.5em;
    margin-right: 0.5em;
  }
  </style>
</head>
<body>
  <img src="/img/vnstat-summary.png"/>
  <img src="/img/vnstat-hourly.png"/>
  <img src="/img/vnstat-daily.png"/>
  <img src="/img/vnstat-monthly.png"/>
</body>
</html>`

var views = []view{
  view{arg:"-s",img:"vnstat-summary.png"},
  view{arg:"-h",img:"vnstat-hourly.png"},
  view{arg:"-d",img:"vnstat-daily.png"},
  view{arg:"-m",img:"vnstat-monthly.png"},
}

var iface = "eth0"
var command = "vnstati"
var commandArgsList = make([][]string,4)

func handler(w http.ResponseWriter, r *http.Request) {
  os.MkdirAll("/tmp/vnstat-webdaemon/img", os.ModePerm);
  for _, args := range commandArgsList {
    cmd := exec.Command(command, args...)
    if err := cmd.Run(); err != nil {
  		fmt.Fprintln(os.Stderr, "Error when running : "+command+" "+strings.Join(args," "), err)
  	}
  }
	fmt.Fprintf(w, pagecontent)
}

func main() {
  // get network interface from command line parameter
  if len(os.Args) > 1 {
    iface = os.Args[1]
  }

  // init commandArgsList
  for _,view := range views {
    commandArgsList = append(commandArgsList,[]string{view.arg, "-i", iface, "-o", "/tmp/vnstat-webdaemon/img/"+view.img})
  }

  fmt.Println("Running...")

  http.Handle("/img/", http.FileServer(http.Dir("/tmp/vnstat-webdaemon")))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":7991", nil)
}
