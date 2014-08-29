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

var C__batchSize = 100
var C__DEBUG = true

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

func send_over_time(num_msgs int, dur time.Duration,text []string) {

    total_in_sec := int(dur.Seconds())   // total num of seconds for the duration
    interval_period := 1            // in seconds
    num_msgs_interval := num_msgs   // number of messages per interval -- init
    
    // performing the dance to figure out what the "ideal" interval should be to send at least one message, starting with a 1 sec interval.
    for i:=1;num_msgs*interval_period/total_in_sec < 1;i++ {
        interval_period = i*5       // we increase by 30s every time we find the interval is too short.
    }
    num_msgs_interval = num_msgs*interval_period/total_in_sec
    
    // now we have a rate per period, we can now start sending!
    fmt.Printf("rate = %d msgs/%d seconds\n",num_msgs_interval,interval_period)
    
    // start the tick now!
    start_time := time.Now()
    in_between := time.Now()
    progress := 0
    for count := 0; count < num_msgs; count += num_msgs_interval {
        in_between = time.Now()
        for i := 0; i< num_msgs_interval; i++ {
            log.Notice(gen_rubbish(10,text))
        }
        // pretty progress bar
        progress = int(20*count/num_msgs)
        fmt.Printf("\r[%s%s] -- %s                 ",strings.Repeat("*",progress),strings.Repeat(" ",20-progress),time.Since(start_time).String())
        // now we just wait
        for ;int(time.Since(in_between).Seconds()) < interval_period; {
            // .1s should be enough so we don't lose too much time waiting and limit the processing overhead
            time.Sleep(time.Millisecond*100)
        }
    }
    fmt.Printf("Sent %d messages in %s\n",num_msgs,time.Since(start_time).String())
    
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

    /* now := time.Now()
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
    } */
    
    tempo,_ := time.ParseDuration("2m")
    send_over_time(10000,tempo,ipsum)
    
    

}
