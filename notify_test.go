package notify

import (
	"os"
	"syscall"
	"testing"
	"time"
)

var pid = os.Getpid()

func TestNotify(t *testing.T) {
	callSize := 0
	off1 := On(syscall.SIGUSR1, func() { callSize++ })
	off2 := On(syscall.SIGUSR2, func() { callSize++ })
	syscall.Kill(pid, syscall.SIGUSR1)
	time.Sleep(time.Millisecond)
	syscall.Kill(pid, syscall.SIGUSR2)
	time.Sleep(time.Millisecond)
	if callSize != 2 {
		t.Fail()
	}
	off1()
	off2()
	syscall.Kill(pid, syscall.SIGUSR1)
	time.Sleep(time.Millisecond)
	syscall.Kill(pid, syscall.SIGUSR2)
	time.Sleep(time.Millisecond)
	if callSize != 2 {
		t.Fail()
	}
}

func TestNotify_Once(t *testing.T) {
	callSize := 0
	Once(syscall.SIGUSR1, func() { callSize++ })
	syscall.Kill(pid, syscall.SIGUSR1)
	time.Sleep(time.Millisecond)
	syscall.Kill(pid, syscall.SIGUSR1)
	time.Sleep(time.Millisecond)
	if callSize != 1 {
		t.Fail()
	}
}
