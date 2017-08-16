log_rubbish_generator
=====================

A project in go to spit random logs at syslog with various options

# INSTALL

1. Retrieve the git code
```
git clone https://github.com/sledigabel/log_rubbish_generator.git
```

2. Install golang

3. Build log_rubbish_generator
```
go build log_rubbish_generator.go
```

4. Options

```
$ ./log_rubbish_generator --help
Usage: ./log_rubbish_generator [--time=<duration>] [--num=<number>] [--ifile=<file>] [--ofile=<file>] [--ml] [--ltrim=<number>]
  -help
    	prints usage
  -ifile string
    	Input file with logs to send
  -ltrim int
    	Number of fields to trim on the left
  -ml
    	Reads as multiline
  -num int
    	Number of messages to send (default 100000)
  -ofile string
    	Output file to send logs to
  -time string
    	Time lapse to send the logs (default "5m")
```

5. Examples

Random rubbish in syslog, default values (100000 in 5m)
```
./log_rubbish_generator
```

Reading from file, output to a new file
```
./log_rubbish_generator -ifile application_sample -ofile test
```

Reading from file, output to a new file, multiline, 1M messages in 10m
```
./log_rubbish_generator -ifile application_sample -ofile test -ml -num 1000000 -time 10m
```

Reading from file, output to a new file, multiline, 1M messages in 10m, remove the 3 first fields
```
./log_rubbish_generator -ifile application_sample -ofile test -ml -num 1000000 -time 10m -ltrim 3
```

- 16/08/2017
Added options: output file, multiline (naive), ltrim

- 18/09/2014
Added options: number of messages, duration, input file
Added feature to read lines from a file.

- 29/08/2014
Added a feature to spread log generation over a period of time
Added a rate display
Added a progress bar with time elapsed and number of messages
