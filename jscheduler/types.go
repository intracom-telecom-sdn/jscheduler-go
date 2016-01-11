package jscheduler

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type CpuPool []int

func NewEmptyCpuPool() CpuPool {
	return make([]int, 0)
}

func NewCpuPool(numCpus int) CpuPool {
	pool := make([]int, numCpus)

	for i := 0; i < numCpus; i++ {
		pool[i] = i
	}

	return pool
}

func ParseCpuPool(pool string) CpuPool {
	elements := strings.Split(pool, ",")
	cpus := NewEmptyCpuPool()
	cpuSet := make(map[int]struct{})
	for _, el := range elements {
		step := 1
		cpuRange := strings.Split(el, "-")
		if len(cpuRange) > 1 {
			rangeSplit := strings.Split(cpuRange[1], ":")
			if len(rangeSplit) > 1 {
				step, _ = strconv.Atoi(rangeSplit[1])
			}
			c1, _ := strconv.Atoi(cpuRange[0])
			c2, _ := strconv.Atoi(rangeSplit[0])
			for c := c1; c <= c2; c += step {
				if _, cpuExists := cpuSet[c]; !cpuExists {
					cpus = append(cpus, c)
					cpuSet[c] = struct{}{}
				}
			}
		} else {
			c, _ := strconv.Atoi(cpuRange[0])
			if _, cpuExists := cpuSet[c]; !cpuExists {
				cpus = append(cpus, c)
				cpuSet[c] = struct{}{}
			}
		}
	}
	fmt.Println("Parsed CPU pool:", cpus)
	return cpus
}

type ThreadPolicy struct {
	Filter string
	Prio   int
	Cpus   CpuPool
}

func NewThreadPolicy() ThreadPolicy {
	return ThreadPolicy{
		Filter: "",
		Prio:   0,
		Cpus:   NewEmptyCpuPool(),
	}
}

type Thread struct {
	Name    string
	Tid     int
	Prio    int
	Cpus    CpuPool
	HasPolicy bool
}

func NewThread(name string, tid int) Thread {
	return Thread{
		Name:    name,
		Tid:     tid,
		Prio:    0,
		Cpus:    NewCpuPool(runtime.NumCPU()),
		HasPolicy: false,
	}
}

func (t *Thread) FilterAndSetPolicy(policy ThreadPolicy) {
	if regexp.MustCompile(policy.Filter).MatchString(t.Name) {
		t.SetPolicy(policy)
	}
}

func (t *Thread) SetPolicy(policy ThreadPolicy) {
	t.Prio = policy.Prio
	t.Cpus = policy.Cpus
	t.HasPolicy = true
}

type ThreadList []Thread

func NewThreadList() ThreadList {
	return make([]Thread, 0)
}

type ThreadPolicyArgList struct {
	Value []ThreadPolicy
}

func NewThreadPolicyArgList() ThreadPolicyArgList {
	return ThreadPolicyArgList{
		Value: make([]ThreadPolicy, 0),
	}
}

func (lst *ThreadPolicyArgList) String() string {
	strLst := make([]string, 0)

	for _, el := range lst.Value {
		filter := el.Filter
		prio := el.Prio
		cpus := strconv.Itoa(el.Cpus[0])
		for _, c := range el.Cpus[1:] {
			cpus += fmt.Sprintf(",%s", strconv.Itoa(c))
		}
		strLst = append(strLst, fmt.Sprintf("\"%s\";%d;%s", filter, prio, cpus))
	}

	return strings.Join(strLst[:], "::")
}

func (lst *ThreadPolicyArgList) Set(s string) error {
	strLst := strings.Split(s, "::")
	lst.Value = make([]ThreadPolicy, 0)
	fmt.Println("Thread Schedule Configuration")
	for _, el := range strLst {
		ts := NewThreadPolicy()
		tsEl := strings.Split(el, ";")
		if tsEl[0] != "" {
			fmt.Printf("    Filter: %s\n", tsEl[0])
			ts.Filter = tsEl[0]
		}
		if tsEl[1] != "" {
			fmt.Printf("    Priority: %s\n", tsEl[1])
			ts.Prio, _ = strconv.Atoi(tsEl[1])
		}
		if tsEl[2] != "" {
			fmt.Printf("    Cpu Pool: %s\n", tsEl[2])
			ts.Cpus = ParseCpuPool(tsEl[2])
		}
		lst.Value = append(lst.Value, ts)
	}

	fmt.Println(lst.Value)
	return nil
}

func (lst *ThreadPolicyArgList) Get() []ThreadPolicy {
	return lst.Value
}

func (lst *ThreadPolicyArgList) IsSet() bool {
	return len(lst.Value) > 0
}
