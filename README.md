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

- 18/09/2014
Added options: number of messages, duration, input file
Added feature to read lines from a file.

- 29/08/2014
Added a feature to spread log generation over a period of time
Added a rate display
Added a progress bar with time elapsed and number of messages


