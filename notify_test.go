package notify

import (
	"os"
	"syscall"
	"testing"
	"time"
)

var pid = os.Getpid()

func TestNotify_On(t *testing.T) {
	callSize := 0
	off1 := On(func() { callSize++ }, syscall.SIGUSR1)
	off2 := On(func() { callSize++ }, syscall.SIGUSR2)
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
	Once(func() { callSize++ }, syscall.SIGUSR1)
	syscall.Kill(pid, syscall.SIGUSR1)
	time.Sleep(time.Millisecond)
	syscall.Kill(pid, syscall.SIGUSR1)
	time.Sleep(time.Millisecond)
	if callSize != 1 {
		t.Fail()
	}
}

func TestNotify_OnSlice(t *testing.T) {
	callSize := 0
	off := On(func() { callSize++ }, syscall.SIGUSR1, syscall.SIGUSR2)

	syscall.Kill(pid, syscall.SIGUSR1)
	time.Sleep(time.Millisecond)
	syscall.Kill(pid, syscall.SIGUSR2)
	time.Sleep(time.Millisecond)
	if callSize != 2 {
		t.Fail()
	}
	off()

	syscall.Kill(pid, syscall.SIGUSR1)
	time.Sleep(time.Millisecond)
	syscall.Kill(pid, syscall.SIGUSR2)
	time.Sleep(time.Millisecond)
	if callSize != 2 {
		t.Fail()
	}
}

func TestNotify_OnceSlice(t *testing.T) {
	callSize := 0
	Once(func() { callSize++ }, syscall.SIGUSR1, syscall.SIGUSR2)
	syscall.Kill(pid, syscall.SIGUSR1)
	time.Sleep(time.Millisecond)
	syscall.Kill(pid, syscall.SIGUSR2)
	time.Sleep(time.Millisecond)
	syscall.Kill(pid, syscall.SIGUSR1)
	time.Sleep(time.Millisecond)
	syscall.Kill(pid, syscall.SIGUSR2)
	time.Sleep(time.Millisecond)
	if callSize != 1 {
		t.Log(callSize)
		t.Fail()
	}
}
