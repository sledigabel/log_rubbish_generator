package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/op/go-logging"
	// "io/ioutil"
	"log/syslog"
	"math/rand"
	"os"
	"strings"
	"time"
)

var log = logging.MustGetLogger("test")

// setting rand
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// file string header == YYYY
var tHeader = time.Now().Format("2006")

// for dummy text
var ipsum = strings.Split("Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.", " ")

func gen_rubbish(length int, blah []string) string {

	var str = ""
	for i := 0; i < length; i++ {
		str = str + " " + blah[r.Int()%len(blah)]
	}
	return str
}

func pick_log(blah []string) string {
	return blah[r.Int()%len(blah)]
}

func send_over_time(num_msgs int, dur time.Duration, text []string, debug bool) {

	total_in_sec := int(dur.Seconds()) // total num of seconds for the duration
	interval_period := 1               // in seconds
	num_msgs_interval := num_msgs      // number of messages per interval -- init

	need_rubbish := len(text) == 0

	// performing the dance to figure out what the "ideal" interval should be to send at least one message, starting with a 1 sec interval.
	for i := 1; num_msgs*interval_period/total_in_sec < 1; i++ {
		interval_period = i * 5 // we increase by 30s every time we find the interval is too short.
	}

	num_msgs_interval = num_msgs * interval_period / total_in_sec
	if debug {
		fmt.Printf("tempo: %s\nnum: %d\n", dur, num_msgs)
		// now we have a rate per period, we can now start sending!
		if interval_period != 1 {
			fmt.Printf("rate = %d msgs/%d seconds\n", num_msgs_interval, interval_period, debug)
		} else {
			fmt.Printf("rate = %d msgs/s\n", num_msgs_interval)
		}
	}

	// start the tick now!
	start_time := time.Now()
	in_between := time.Now()
	progress := 0

	for count := 0; count < num_msgs; count += num_msgs_interval {
		in_between = time.Now()
		for i := 0; i < num_msgs_interval; i++ {
			if need_rubbish {
				log.Notice(gen_rubbish(10, ipsum))
			} else {
				// traditional logging
				log.Notice(pick_log(text))
			}
		}
		// pretty progress bar
		// TODO: make it optional
		if debug{
			progress = int(20 * count / num_msgs)
			fmt.Printf("\r[%s%s] --\t%s\t-- %d", strings.Repeat("*", progress), strings.Repeat(" ", 20-progress), strings.Split(time.Since(start_time).String(), ".")[0]+"s", count)
		}
		// now we just wait
		// TODO: Make it better, this is ugly
		for int(time.Since(in_between).Seconds()) < interval_period {
			// .1s should be enough so we don't lose too much time waiting and limit the processing overhead
			time.Sleep(time.Millisecond * 100)
		}
	}
	if debug {
		fmt.Printf("Sent %d messages in %s\n", num_msgs, time.Since(start_time).String())
	}

}

func usage() {
	fmt.Printf("Usage: %s [--time=<duration>] [--num=<number>] [--ifile=<file>] [--ofile=<file>] [--ml] [--ltrim=<number>] [--]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func multiLineSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {

    // Return nothing if at end of file and no data passed
    if atEOF && len(data) == 0 {
        return 0, nil, nil
    }

    // Find the index of the input of a newline followed by a
    // pound sign.
    if i := strings.Index(string(data), fmt.Sprintf("\n%s",tHeader)); i >= 0 {
        return i + 1, data[0:i], nil
    }

    // If at end of file with data return the data
    if atEOF {
        return len(data), data, nil
    }

    return
}


func main() {

	param_help := flag.Bool("help", false, "prints usage")
	param_tempo := flag.String("time", "5m", "Time lapse to send the logs")
	param_num := flag.Int("num", 100000, "Number of messages to send")
	param_ifile := flag.String("ifile", "", "Input file with logs to send")
	param_ofile := flag.String("ofile", "", "Output file to send logs to")
	param_multiline := flag.Bool("ml", false, "Reads as multiline")
	param_ltrim := flag.Int("ltrim", 0, "Number of fields to trim on the left")
	param_progress := flag.Bool("progress", true, "Display progress bar and summary")

	flag.Usage = usage
	flag.Parse()

	var logs []string
	var is_debug bool = *param_progress

	if *param_help {
		usage()
	}
	if len(*param_ifile) != 0 {
		// reading the file and random logs
		inFile, err := os.Open(*param_ifile)
		if err != nil {
			//Do something
			fmt.Println(err)
			os.Exit(5)
		}
		// TODO: improve here to handle
		scanner := bufio.NewScanner(inFile)
		if *param_multiline {
			scanner.Split(multiLineSplit)
		}
		for scanner.Scan(){
			l := scanner.Text()
			if *param_ltrim > 0 {
				spl := strings.SplitN(l, string(' '), *param_ltrim+1)
				logs = append(logs, spl[len(spl)-1])
			} else {
				logs = append(logs, scanner.Text())
			}
		}
	}

	if len(*param_ofile) != 0 {
		// logging out to a file
		var format = logging.MustStringFormatter(`%{time:2006-01-02T15:04:05.999:00} %{message}`)
		outputfile, err := os.OpenFile(*param_ofile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(6)
		}
		fileBackend := logging.NewLogBackend(outputfile, "", 0)
		fileFormatter := logging.NewBackendFormatter(fileBackend, format)
		fileLeveled := logging.AddModuleLevel(fileFormatter)
		logging.SetBackend(fileLeveled)
	} else {
		// no output file provided... assuming syslog
		syslogBackend, err := logging.NewSyslogBackendPriority("", syslog.LOG_LOCAL5)
		if err != nil {
			log.Fatal(err)
		}
		// Combine them both into one logging backend.
		//logging.SetBackend(logBackend, syslogBackend)
		logging.SetBackend(syslogBackend)
	}

	tempo, _ := time.ParseDuration(*param_tempo)
	send_over_time(*param_num, tempo, logs, is_debug)

}
