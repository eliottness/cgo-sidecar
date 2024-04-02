package main

import (
	"bytes"
	"os"
	"os/exec"
	"reflect"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"
)

func main() {
	shmid, err := os.OpenFile("/dev/shm/cgo-sidecar", os.O_RDWR|os.O_CREATE|syscall.O_NOFOLLOW|syscall.O_CLOEXEC, 0600)
	if err != nil {
		panic(err)
	}

	err = syscall.Ftruncate(int(shmid.Fd()), int64(reflect.TypeOf(int32(0)).Size()))
	if err != nil {
		panic(err)
	}

	// Close shared memory
	defer func ()  {
		shmid.Close()
		os.Remove("/dev/shm/cgo-sidecar")
	}()

	data, err := syscall.Mmap(int(shmid.Fd()), 0, int(reflect.TypeOf(int32(0)).Size()), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}

	init := (*int32)(unsafe.Pointer(&data[0]))

	var buffer bytes.Buffer
	cmd := exec.Command("../target/debug/cgo-sidecar","/dev/shm/cgo-sidecar")
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 30; i++ {
		atomic.AddInt32(init, 1)
		time.Sleep(1 * time.Millisecond)
	}

	cmd.Process.Signal(syscall.SIGTERM)

	cmd.Wait()
	println(string(buffer.String()))

}