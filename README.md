# Jscheduler


## Description

Jscheduler is a utility tool that allows the user to monitor a process that is running on the JVM and enforce a scheduling policy on the threads created by that process at runtime. 

The motivation behind the Jscheduler is the lack of an easy way to modify the priority and the reserved resources of a specific Java thread without getting down to JNI code or setting up and using 3rd party libraries. 

The main purpose of Jscheduler is to be used as a testing tool to measure experimentally the impact of certain threads on the overall system performance for an arbitrary large project. 

By being able to easily change the throughput and/or the latency of processing elements in predetermined spots inside the codebase one can locate performance bottlenecks or zero in an optimal scheduling policy more efficiently. 


## Features

- Live monitoring of JVM processes using JStack
- Dynamic thread dump parsing
- Dynamic name based thread matching
- Dynamic CPU affinity enforcement
- Dynamic thread priority enforcement
- Low footprint execution

## Getting Started

### Dependencies
- A recent version of Go (go1.4 and go1.5 are supported)
- [golang.org/x/sys/unix](https://godoc.org/golang.org/x/sys/unix) library
- Java1.4+
- Properly set `$GOPATH` and `$JAVA_HOME` environmental variables
- Superuser permissions (if you need to increase thread priorities)

### Setup

- Get the `golang.org/x/sys/unix` library
  ```bash
  go get golang.org/x/sys/unix
  ```
- Get Jscheduler
  ```bash
  go get github.com/intracom-telecom-sdn/jscheduler-go
  ```
- Compile and install Jscheduler
  ```bash
  go build -o $GOPATH/bin/jscheduler $GOPATH/src/github.com/intracom-telecom-sdn/jscheduler-go/jscheduler.go
  ```

## Usage

For the purpose of clarity we must define that
> A scheduling policy or simply policy from now on will refer to a tuple of (thread name filter, thread priority, thread CPU affinity)

You can invoke the Jscheduler from the command line with the following options

```bash
jscheduler-go --pid <java_process_pid> --interval <monitoring_interval> --policies <thread_policy_list>
```

_Parameters_

- `pid` is the monitored Java process pid
- `interval` is the monitoring interval. The default value is 3 seconds. _Note_ that the capture of a thread dump is a relatively expensive operation, so the interval should be relatively large if you care about the Jscheduler execution footprint.
- `policies` is a list of the scheduling policies

The general syntax for the `policies` list is given below

```
"threadNameRegex1;threadPriority1;cpuPool1::threadNameRegex2;threadPriority2;cpuPool2::..."
```

Here the consecutive list elements are separated with the operator `::`, while the fields of each element are separated with a semicolon `;`.

You can leave priority and cpu pool fields unspecified, in which case they will take the default value given by the OS. The semicolons must exists. An example is given here

```
"threadNameRegex1;;cpuPool1::threadNameRegex2;threadPriority2;::..."
```

The `cpuPool` field follows the `taskset` command syntax, i.e. it is a numerical list of processors separated by commas and may include ranges. For example `1,3,10-16:2` stands for the CPUS 1,3,10,12,14,16. 


We give an example use case in the following gif
![zero2hero.gif: Image not found](https://raw.githubusercontent.com/intracom-telecom-sdn/jscheduler-go/master/figs/zero2hero.gif) 
*<p align="center">Zero to Hero with Jscheduler</p>*

**Test characteristics**
- We use a mock Java benchmark that creates 2 threadpools with 5 threads each
- All the threads perform the same unit of CPU intensive work repeatedly
- The benchmark runs on a VM with 4 VCPUs and 4GB of RAM
At first we demonstrate how to run a successful setup of the Jscheduler. Then we compile and run the benchmark. Lastly we execute 
the Jscheduler and enforce the following policies:
- `pool-1-thread-2` thread: Highest priority, isolate in CPU 0
- `pool-1-thread-*` and `pool-2-thread-*` threads: Lowest priority, run in CPUs 1,2,3

*NOTES* 
1. The `pool-1-thread-2` thread throughput has a _3x_ increase, i.e. from ~10 jobs/sec to ~30 jobs/sec. This is expected because in the default configuration each cpu is 



