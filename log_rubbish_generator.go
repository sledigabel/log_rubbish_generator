package main

import (
    "fmt"
    "flag"
    "io/ioutil"
    "os"
    "strings"
    "math/rand"
    "time"
    "log/syslog"
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

func pick_log (blah []string) string {
    return blah[r.Int()%len(blah)]
}


func send_over_time(num_msgs int, dur time.Duration,text []string) {

    total_in_sec := int(dur.Seconds())   // total num of seconds for the duration
    interval_period := 1            // in seconds
    num_msgs_interval := num_msgs   // number of messages per interval -- init

    need_rubbish := ( len(text) == 0 )
    
    // performing the dance to figure out what the "ideal" interval should be to send at least one message, starting with a 1 sec interval.
    for i:=1;num_msgs*interval_period/total_in_sec < 1;i++ {
        interval_period = i*5       // we increase by 30s every time we find the interval is too short.
    }
    num_msgs_interval = num_msgs*interval_period/total_in_sec
    
    // now we have a rate per period, we can now start sending!
    if interval_period != 1 {
        fmt.Printf("rate = %d msgs/%d seconds\n",num_msgs_interval,interval_period)
    } else {
        fmt.Printf("rate = %d msgs/s\n",num_msgs_interval)
    }
    
    // start the tick now!
    start_time := time.Now()
    in_between := time.Now()
    progress := 0
    
    for count := 0; count < num_msgs; count += num_msgs_interval {
        in_between = time.Now()
        for i := 0; i< num_msgs_interval; i++ {
            if need_rubbish {
                log.Notice(gen_rubbish(10,ipsum))
            } else{
                // traditional logging
                log.Notice(pick_log(text))
            }
        }
        // pretty progress bar
        progress = int(20*count/num_msgs)
        fmt.Printf("\r[%s%s] --\t%s\t-- %d",strings.Repeat("*",progress),strings.Repeat(" ",20-progress),strings.Split(time.Since(start_time).String(),".")[0]+"s",count)
        // now we just wait
        for ;int(time.Since(in_between).Seconds()) < interval_period; {
            // .1s should be enough so we don't lose too much time waiting and limit the processing overhead
            time.Sleep(time.Millisecond*100)
        }
    }
    fmt.Printf("Sent %d messages in %s\n",num_msgs,time.Since(start_time).String())
    
}

func usage() {
    
    fmt.Printf("Usage: %s [--time=<duration>] [--num=<number>] [--file=<input_file>]\n",os.Args[0])
    flag.PrintDefaults()
    os.Exit(1)
}

func main() {
    // Customize the output format
    //logging.SetFormatter(logging.MustStringFormatter("▶ %{level:.1s} 0x%{id:x} %{message}"))

    // Setup one stdout and one syslog backend.
    //logBackend := logging.NewLogBackend(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile)
    //logBackend.Color = true

    syslogBackend, err := logging.NewSyslogBackendPriority("",syslog.LOG_LOCAL5)
    if err != nil {
        log.Fatal(err)
    }

    // Combine them both into one logging backend.
    //logging.SetBackend(logBackend, syslogBackend)
    logging.SetBackend(syslogBackend)

    param_help := flag.Bool("help",false,"prints usage")
    param_tempo := flag.String("time","5m","Time lapse to send the logs")
    param_num := flag.Int("num",100000,"Number of messages to send")
    param_ifile := flag.String("file","","Input file with logs to send")
    param_groot := flag.Bool("groot",false,"I AM GROOT")
    
    flag.Usage = usage
    flag.Parse()
    
    var logs []string

    if *param_groot {
        fmt.Println("I AM GROOT !!!")
    }
    if *param_help {
        usage()
    }
    if ( *param_ifile != "" ){
        // reading the file and random logs
        content, err := ioutil.ReadFile(*param_ifile)
        if err != nil {
            //Do something
            fmt.Println(err)
            os.Exit(5)
        }
        logs = strings.Split(string(content), "\n")
    }
    
    tempo,_ := time.ParseDuration(*param_tempo)
    fmt.Printf("tempo: %s\n",tempo)
    send_over_time(*param_num,tempo,logs)
    
}
