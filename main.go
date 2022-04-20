package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/Cosiamo/GoPortScan/port"
)

var hostname string
var tcp, udp, help bool

func main() {
	flag.StringVar(&hostname, "host", "localhost", "Sets hostname")
	flag.BoolVar(&tcp, "tcp", false, "Returns TCP packets")
	flag.BoolVar(&udp, "udp", false, "Returns UDP packets")
	flag.BoolVar(&help, "help", false, "Lists flags")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		return
	}

	if !tcp && !udp {
		fmt.Println("Please use the flag:")
		fmt.Println("	-tcp")
		fmt.Println("	and/or")
		fmt.Println("	-udp")
		fmt.Println(" ")
		fmt.Println("For more info use -help")
		return
	}

	start := time.Now()

	for results := range port.InitialScan(hostname, tcp, udp) {
		fmt.Print(results)
	}

	res := TimeMeasurement(start)
	fmt.Println(res)
}

func TimeMeasurement(start time.Time) string {
	elapsed := time.Since(start)
	res := fmt.Sprintf("Execution time: %s", elapsed)
	return res
}
