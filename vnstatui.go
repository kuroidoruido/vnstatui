package main

import (
	"fmt"
	"net/http"
	"os/exec"
  "path/filepath"
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
  <img src="/img/vnstat-daily.png"/>
  <img src="/img/vnstat-hourly.png"/>
  <img src="/img/vnstat-five-minutes.png"/>
  <img src="/img/vnstat-hour-graph.png"/>
</body>
</html>`

var views = []view{
  view{arg:"-s",img:"vnstat-summary.png"},
  view{arg:"-h",img:"vnstat-hourly.png"},
  view{arg:"-d",img:"vnstat-daily.png"},
  view{arg:"-5",img:"vnstat-five-minutes.png"},
  view{arg:"-hg",img:"vnstat-hour-graph.png"},
}

var iface = "enp3s0+wlp2s0"
var command = "vnstati"
var commandArgsList = make([][]string,4)

func handler(w http.ResponseWriter, r *http.Request) {
  for _, args := range commandArgsList {
    cmd := exec.Command(command, args...)
    fmt.Fprintln(os.Stdout, cmd)
    if err := cmd.Run(); err != nil {
  		fmt.Fprintln(os.Stderr, "Error when running : "+command+" "+strings.Join(args," "), err)
  	}
  }
	fmt.Fprintf(w, pagecontent)
}

func main() {
  path := filepath.Join("/tmp", "vnstat-webdaemon")
  os.MkdirAll(path, os.ModePerm)
  fmt.Fprintln(os.Stdout, path)

  // get network interface from command line parameter
  if len(os.Args) > 1 {
    iface = os.Args[1]
  }

  // init commandArgsList
  for _,view := range views {
    commandArgsList = append(commandArgsList,[]string{view.arg, "-i", iface, "-o", filepath.Join(path,view.img)})
  }

  fmt.Println("Running...")

  http.Handle("/img/", http.FileServer(http.Dir(path)))
	http.HandleFunc("/", handler)
	http.ListenAndServe(":7991", nil)
}
