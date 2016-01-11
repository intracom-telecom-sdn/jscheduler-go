package jscheduler

import (
	"fmt"
	"golang.org/x/sys/unix"
	"runtime"
	"syscall"
	"unsafe"
    "regexp"
)

func AdjustThreadPolicies(threads *ThreadList, policies []ThreadPolicy) {
	for _, policy := range policies {
		for i, _ := range *threads {
			if regexp.MustCompile(policy.Filter).MatchString((*threads)[i].Name) {
				if !(*threads)[i].HasPolicy {
					(*threads)[i].SetPolicy(policy)
				}
			}
		}
	}
}

// SetAffinity attend the cpu list to pid,
// note: SetAffinity apply to thread ID only,
// to fully control one process, call SetAffinity for all thread of the process.
// use os.GetThreadIDs() to get all thread of the process
// check ghttps://github.com/golang/go/issues/11243
func SetAffinity(pid int, cpus []int) error {
	var mask [1024 / 64]uintptr
	if pid <= 0 {
		pidget, _, _ := syscall.RawSyscall(unix.SYS_GETPID, 0, 0, 0)
		pid = int(pidget)
	}
	for _, cpuIdx := range cpus {
		cpuIndex := uint(cpuIdx)
		mask[cpuIndex/64] |= 1 << (cpuIndex % 64)
	}
	_, _, err := syscall.RawSyscall(unix.SYS_SCHED_SETAFFINITY, uintptr(pid), uintptr(len(mask)*8), uintptr(unsafe.Pointer(&mask[0])))
	if err != 0 {
		return err
	}
	return nil
}

func SetAffinityThreadGroup(threads *ThreadList) error {
	for _, t := range *threads {
		if !t.HasPolicy {
			continue
		}
		if err := SetAffinity(t.Tid, t.Cpus); err != nil {
			return err
		}
	}
	return nil
}

func SetPriorityThreadGroup(threads *ThreadList) error {
	for _, t := range *threads {
		if !t.HasPolicy {
			continue
		}
		if err := unix.Setpriority(unix.PRIO_PROCESS, t.Tid, t.Prio); err != nil {
			return err
		}
	}
	return nil
}

func RescheduleThreadGroup(threads *ThreadList) error {
	for _, t := range *threads {
		if !t.HasPolicy {
			continue
		}
		if len(t.Cpus) == 0 {
			t.Cpus = NewCpuPool(runtime.NumCPU())
		}
		fmt.Println("Pinning thread", t.Name, "to CPU set", t.Cpus)
		if err := SetAffinity(t.Tid, t.Cpus); err != nil {
			return err
		}
		fmt.Println("Setting thread", t.Name, "priority to", t.Prio)
		if err1 := unix.Setpriority(unix.PRIO_PROCESS, t.Tid, t.Prio); err1 != nil {
			return err1
		}
	}
	return nil
}
