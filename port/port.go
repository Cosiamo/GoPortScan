package port

import (
	"net"
	"strconv"
	"sync"
	"time"

	color "github.com/fatih/color"
)

type ScanResult struct {
	Port 	string
	State 	string
}

// checks of a port is open
func ScanPort(protocol string, hostname string, port int) ScanResult {
	portStr := strconv.Itoa(port)
	closed := color.RedString("Closed")
	open := color.GreenString("Open")

	result := ScanResult{Port: protocol + "/" + color.YellowString(portStr)}
	address := hostname + ":" + portStr
	conn, err := net.DialTimeout(protocol, address, 60*time.Second)

	// if the connection isn't accepting any requests
	// unless you have a hostname, then return false
	if err != nil {
		result.State = closed
		return result
	}
	// tried using ScanPort as a goroutine but couldn't
	// because net.Conn kept giving me the error:
	// invalid memory address or nil pointer dereference goroutines
	// no matter what I tried, I couldn't get this to work concurrently
	defer conn.Close()

	result.State = open

	// open ports return much faster than closed ports
	// sleeping every time a result state is open
	time.Sleep(2350 * time.Millisecond)
	return result
}

// scans lower end of the port range (1 - 1024)
// these are well-known ports that are preallocated
// HTTP:80 SSH:22 FTP:21
func Scan(hostname string, tcp bool, udp bool) (<-chan ScanResult) {
	res := make(chan ScanResult)
	odds := make(chan int)
	evens := make(chan int)
	threes := make(chan int)
	fives := make(chan int)

	wg := &sync.WaitGroup{}
	wg.Add(5)
	
	endPORT := 1024
	startPORT := 1

	// ===Fan In / Fan Out goroutine pattern===
	go portNums(startPORT, endPORT, odds, evens, threes, fives, wg)
	go result(odds, res, tcp, udp, hostname, wg)
	go result(evens, res, tcp, udp, hostname, wg)
	go result(threes, res, tcp, udp, hostname, wg)
	go result(fives, res, tcp, udp, hostname, wg)
	go func() {
		wg.Wait()
		close(res)
	}()

	return res
}

func portNums(startPORT int, endPORT int, odds chan int, evens chan int, threes chan int, fives chan int, wg *sync.WaitGroup) {
	// Mutex prevents race conditions
	// race conditions happen due to repeated thread access
	var m sync.Mutex
	go func() {
		defer wg.Done()
		for startPORT <= endPORT {
			time.Sleep(16 * time.Millisecond)
			if startPORT%5 == 0 {
				fives <- startPORT
			} else if startPORT%3 == 0 {
				threes <- startPORT
			} else if startPORT%2 == 0 {
				evens <- startPORT
			} else {
				odds <- startPORT
			}
			// mutext lock
			m.Lock()
			// modify num
			startPORT++
			// mutext unlock
			m.Unlock()
		}
		close(odds)
		close(evens)
		close(threes)
		close(fives)
	}()
}

func result(ch chan int, res chan ScanResult, tcp bool, udp bool, hostname string, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range ch {
		if tcp {
			res <- ScanPort("tcp", hostname, v)
		}
		if udp {
			res <- ScanPort("udp", hostname, v)
		}
	}
}