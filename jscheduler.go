package main

import (
	"flag"
	"fmt"
	"github.com/georgepar/jscheduler-go/jscheduler"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"
)

func main() {
	// Get command line args
	var pid string
	var interval int
	var help bool
	policies := jscheduler.NewThreadPolicyArgList()

	policiesUsage := `The threads which need to be rescheduled with the
    respective scheduling policies.
    Must be given in the format:
        "threadNameRegex1;threadPriority1;cpuPool1::threadNameRegex2;threadPriority2;cpuPool2::..."
    The above configuration will pin the threads that match with threadNameRegex1 to cpuPool1
    with priority threadPriority1 e.t.c.
    Priorities and cpu pools may be left unspecified (but the semicolons must exist),
    in which case the default values given by the OS will be left untouched.
    For example:
        "threadNameRegex1;cpuPool1::threadNameRegex2;threadPriority2;::..."
    `

	flag.BoolVar(&help, "help", false, "Display usage information")
	flag.StringVar(&pid, "pid", "-1", "The pid of the monitored java process. This argument is required.")
	flag.IntVar(&interval, "interval", 3000, "Time to wait between polling jstack in milliseconds. Default value is 3s.")
	flag.Var(&policies, "policies", policiesUsage)

	flag.Parse()

	if pid == "-1" {
		flag.Usage()
		os.Exit(1)
	}

	threadCount := make(map[string]int)
	modifiedThreads := make(map[string]struct{})

	// Print thread occurrence count on CTRL-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		printThreadCount(threadCount)
		os.Exit(1)
	}()

	for {
		// Get thread dump
		threadDump, err := jscheduler.GetJstackThreadDump(os.Getenv("JAVA_HOME"), pid)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Parse thread dump
		threads, err := jscheduler.ParseThreadDump(threadDump, modifiedThreads)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Filter and adjust thread policies
		jscheduler.AdjustThreadPolicies(threads, policies.Get())

		// Set Thread affinities and priorities
		jscheduler.RescheduleThreadGroup(threads)

		for _, t := range *threads {
			if t.HasPolicy {
				modifiedThreads[t.Name] = struct{}{}
			}
            threadCount[t.Name]++
		}

		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}

func printThreadCount(threadCount map[string]int) {
	if len(threadCount) == 0 {
		fmt.Println("No threads found")
	}
	keys := make([]string, 0, len(threadCount))
	for k := range threadCount {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s: %d\n", k, threadCount[k])
	}
}
