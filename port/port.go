package port

import (
	"net"
	"strconv"
	"time"

	color "github.com/fatih/color"
)

type ScanResult struct {
	Port string
	State string
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
	return result
}

// scans lower end of the port range (1 - 1024)
// these are well-known ports that are preallocated
// HTTP:80 SSH:22 FTP:21
func InitialScan(hostname string, tcp bool, udp bool) (chan []ScanResult) {
	res := make(chan []ScanResult)
	var results []ScanResult

	go func() {
		// i is the port num
		// iterates over the first 1024 ports
		for i := 1; i <= 1024; i++ {
			if tcp {
				res <- append(results, ScanPort("tcp", hostname, i))
			}
			if udp {
				res <- append(results, ScanPort("udp", hostname, i))
			}
		}
		close(res)
	}()

	return res
}