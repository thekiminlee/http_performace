package main

import (
	"fmt"
	// "net/http"
	"os"
	"sort"
	"math"
	"time"

	"net/url"
	"net"
	"bufio"
	"io"
	"log"
	"strings"

	flag "github.com/ogier/pflag"
)

var (
	workerURL   string
	profile     int

	total float64 = 0.00
	performanceArr []float64
	fastest float64 
	slowest float64
	success int = 0
	failure int = 0
	errorCodes []string
	smallest = math.Inf(1)
	largest = math.Inf(-1)
)

func init() {
	flag.StringVar(&workerURL, "url", "", "(required) URL of the target for analysis")
	flag.IntVar(&profile, "profile", 0, "Number of requests to make to target URL for analysis. Outputs performance summary after executing requests")
}

func main() {
	flag.Parse()

	if flag.NFlag() == 0 || workerURL == "" {
		fmt.Printf("Usage: %s [options]\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	} else if flag.NFlag() == 2 {
		analyze(true)
	} else {
		analyze(false)
	}
}

func analyze(checkPeformance bool) {
	response, _, err := httpGet()
	if err {		
		fmt.Println("Error: " + response)
	} else {
		fmt.Println("Response:")
		fmt.Println(response)
	}
	
	if checkPeformance {
		performance := execGet()
		fmt.Println()
		fmt.Println(performance)
	}
}

func execGet() string {
	for i := 0; i < profile; i++ {
		start := time.Now()
		res, size, err := httpGet()
		end := time.Now()
		elapsed := float64(end.Sub(start))
		total += elapsed
		performanceArr = append(performanceArr, elapsed)

		if smallest > float64(size) {
			smallest = float64(size)
		} 
		if largest < float64(size) {
			largest = float64(size)
		}

		if err {
			failure++
			errorCodes = append(errorCodes, res)
		} else {
			success++
		}
	}

	sort.Float64s(performanceArr)
	fastest = performanceArr[0]
	slowest = performanceArr[profile-1]

	summary := stringResult()

	return summary
}

func stringResult() string {
	result := ""
	result += fmt.Sprintf("Number of requests: %v\n", profile)
	result += fmt.Sprintf("Fastest time: %v\n", time.Duration(fastest))
	result += fmt.Sprintf("Slowest time: %v\n", time.Duration(slowest))
	result += fmt.Sprintf("Average time: %v\n", time.Duration(total/float64(profile)))
	result += fmt.Sprintf("Median time: %v\n", time.Duration(getMedian(performanceArr, len(performanceArr))))
	result += fmt.Sprintf("Success rate: %v", (success/(success+failure)) * 100)
	result += fmt.Sprint("%\n")
	result += fmt.Sprintf("Error codes: %v\n", errorCodes)
	result += fmt.Sprintf("Smallest response size: %v\n", smallest)
	result += fmt.Sprintf("Largest response size: %v\n", largest)

	return result
}

func getMedian(arr []float64, length int) float64{
	var median float64
	mid := length/2
	if length % 2 == 0 {
		median = (arr[mid] + arr[mid-1])/2
	} else {
		median = arr[mid]
	}

	return median
} 

func httpGet() (string, int, bool) {
	var resp string
	var errorCode string
	var fail bool = false
	var size int = 0

	u, err := url.Parse(workerURL)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.Dial("tcp", u.Host+":80")
	if err != nil {
		log.Fatal(err)
	}

	req := fmt.Sprintf("GET %v HTTP/1.1\r\n", u.Path)
	req += fmt.Sprintf("Host: %v\r\n", u.Host)
	req += fmt.Sprintf("Connection: close\r\n")
	req += fmt.Sprintf("\r\n")

	_, err = conn.Write([]byte(req))
	if err != nil {
		log.Fatal(err)
	}
	
	respReader := bufio.NewReader(conn)
	for {
		lineByte, err := respReader.ReadBytes('\n')
		line := string(lineByte)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		
		header := strings.Split(line, " ")
		if header[0] == "HTTP/1.1" {
			if header[1] != "200" {
				errorCode = header[1]
				fail = true
			}
		}
		resp += fmt.Sprint(line)
		size += len(lineByte)
	}
	conn.Close()
	
	if fail {
		return errorCode, size, fail
	}
	return resp, size, fail

}

// func HttpGet() (string, int64, bool) {
// 	resp, err := http.Get(workerURL)
// 	size := resp.ContentLength
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)

// 	if err != nil {
// 		statusCode := fmt.Sprint(resp.StatusCode)
// 		return statusCode, size, true
// 	}

// 	if size == -1 {
// 		size = int64(len(body))
// 	}

// 	return string(body), size, false
// }
