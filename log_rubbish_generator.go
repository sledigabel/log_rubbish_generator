package main

import (
    // stdlog "log"
    // "os"
    "fmt"
    "strings"
    "math/rand"
    "time"
    "github.com/op/go-logging"
)

var batchSize = 100

var log = logging.MustGetLogger("test")
// setting rand
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// for dummy text

var ipsum = strings.Split("Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."," ")

func gen_rubbish (length int,blah []string) string {

    var str = ""
    for i := 0;i<length;i++ {
        str = str+" "+blah[r.Int()%len(blah)]
    }
    
    return str

}


func main() {
    // Customize the output format
    //logging.SetFormatter(logging.MustStringFormatter("▶ %{level:.1s} 0x%{id:x} %{message}"))

    // Setup one stdout and one syslog backend.
    //logBackend := logging.NewLogBackend(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile)
    //logBackend.Color = true

    syslogBackend, err := logging.NewSyslogBackend("")
    if err != nil {
        log.Fatal(err)
    }

    // Combine them both into one logging backend.
    //logging.SetBackend(logBackend, syslogBackend)
    logging.SetBackend(syslogBackend)

    now := time.Now()
    count := 0
    for ;; {
        for i := 0; i< batchSize; i++ {
            log.Notice(gen_rubbish(10,ipsum)) 
        }
        count += 100
        elapsed := time.Since(now).Seconds()
        if elapsed > 5 {
            fmt.Printf("generated %d in %d seconds\n",count,int(elapsed))
            now = time.Now()
            count = 0
        }
    //time.Sleep(5*time.Second)
    //elapsed = time.Since(now)
    //fmt.Println("%d",elapsed)
    }
    
    

}
