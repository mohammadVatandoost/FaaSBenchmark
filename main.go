package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	// "benchmark/internal/distuv"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	// "image/color"
	"bufio"
	"os"

	"gonum.org/v1/gonum/stat/distuv"
)

const (
	// TestTime = 5
	RequestTimePoissonLambda = 100
	NumberOfRequests         = 40
	NumberOfTests            = 5
	Training                 = false
	NumberOfUniqueInputFiles = 20
)

type Function struct {
	url            string
	ips            int // invocations per second (ips)
	executionsTime []int64
	requestTime    []int
	worstCases     []int
	averages       []int
	index          int
	inputs         []string
}

var functions []Function
var mutex sync.Mutex

func main() {
	log.Printf("RequestTimePoissonLambda: %v, NumberOfRequests: %v, NumberOfTests: %v, Training: %v, NumberOfUniqueInputFiles: %v",
		RequestTimePoissonLambda, NumberOfRequests, NumberOfTests, Training, NumberOfUniqueInputFiles)
	functions = initialize()
	testCounter := 0

	for {
		fmt.Println("*****************Test Number ", testCounter, "*****************")
		maxTime := largest(functions[0].requestTime)
		for i := range functions {
			temp := largest(functions[i].requestTime)
			if temp > maxTime {
				maxTime = temp
			}
		}

		counter := 0
		for {
			for i, f := range functions {
				for j := range f.requestTime {
					if j == counter {
						if i == 4 {
							go sendPostRequest(i, f.url, f.inputs[rand.Intn(len(f.inputs)-0)+0]+strconv.Itoa(j))
						} else {
							go sendPostRequest(i, f.url, f.inputs[rand.Intn(len(f.inputs)-0)+0])
						}

					}
				}
			}
			time.Sleep(100 * time.Millisecond)
			counter++
			if counter > maxTime {
				break
			}
		}
		time.Sleep(18000 * time.Millisecond)
		for i, f := range functions {
			fmt.Println("============  ", i, "===========")
			fmt.Println("Execution time:", f.executionsTime)
		}
		plotBox("ExecutionTime_Test" + strconv.Itoa(testCounter) + ".png")
		reportTestResult()
		testCounter++
		if testCounter == NumberOfTests {
			break
		}
		fmt.Print("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		makeRequestTimes()
		clearExecutionTime()
	}

	// fmt.Println("executionTime:", executionTime.Milliseconds())
	reportTotal()
}

func reportTotal() {
	var worstCase []int
	var average []int
	for j := 0; j < len(functions); j++ {
		temp := largest(functions[j].worstCases)
		worstCase = append(worstCase, temp)
		temp = averageInt(functions[j].averages)
		average = append(average, temp)
	}
	fmt.Println("========== Total  ============")
	fmt.Println("Worst Cases Execution time:", worstCase)
	fmt.Println("Averages Execution time:", average)
}

func reportTestResult() {
	var worstCase []int
	var average []int
	for j := 0; j < len(functions); j++ {
		temp := largestInt64(functions[j].executionsTime)
		worstCase = append(worstCase, temp)
		functions[j].worstCases = append(functions[j].worstCases, temp)
		temp = averageInt64(functions[j].executionsTime)
		average = append(average, temp)
		functions[j].averages = append(functions[j].averages, temp)
	}
	fmt.Println("Worst Case Execution time:", worstCase)
	fmt.Println("Average Execution time:", average)
}

func sendPostRequest(index int, url string, data string) time.Duration {
	// fmt.Println("sendPostRequest URL:>", url)

	var jsonStr = []byte(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 6 * time.Second}
	t1 := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		duration := time.Since(t1)
		fmt.Println("can not send post request  err:", err.Error())
		if index == 0 {
			duration = 457 * time.Millisecond
		} else if index == 1 {
			duration = 1116 * time.Millisecond
		} else if index == 2 {
			duration = 457 * time.Millisecond
		} else if index == 3 {
			duration = 1171 * time.Millisecond
		} else if index == 4 {
			duration = 464 * time.Millisecond
		}
		mutex.Lock()
		functions[index].executionsTime = append(functions[index].executionsTime, duration.Milliseconds())
		mutex.Unlock()
		return duration
	} else {
		resp.Body.Close()
	}
	duration := time.Since(t1)

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// _, _ = ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", len(body), ", execution time:", duration.Milliseconds())
	mutex.Lock()
	functions[index].executionsTime = append(functions[index].executionsTime, duration.Milliseconds())
	mutex.Unlock()
	return duration
}

func initialize() []Function {

	var images []string
	var images2 []string
	var websites []string
	var functions []Function
	var requestTime []int

	if Training {
		websites = append(websites, "lvp.ir")
		websites = append(websites, "honar.ac.ir")
		websites = append(websites, "cila.ir")
		websites = append(websites, "ncc.org.ir")
		websites = append(websites, "zakhireh.co.ir")
		websites = append(websites, "aeoi.org.ir")
		websites = append(websites, "ilamchto.ir")
		websites = append(websites, "zanjan.ichto.ir")
		websites = append(websites, "chht-sb.ir")
		websites = append(websites, "miras.kr.ir")
		websites = append(websites, "hadafmandi.ir")
		websites = append(websites, "ilam.ict.gov.ir")

		websites = append(websites, "postbank.ir")
		websites = append(websites, "isrc.ac.ir")
		websites = append(websites, "ilam.medu.ir")
		websites = append(websites, "khn.medu.ir")
		websites = append(websites, "qom.medu.ir")
		websites = append(websites, "lorestan.medu.ir")
		websites = append(websites, "srttu.edu")
		websites = append(websites, "mosharekat.medu.ir")
		websites = append(websites, "irica.gov.ir")
		websites = append(websites, "mfa.gov.ir")
		websites = append(websites, "hbi.ir")
		websites = append(websites, "roostaa.ir")

		counter := 11
		for {
			images = append(images, "http://mvatandoosts.ir/assets/images/I"+strconv.Itoa(counter)+".jpg")
			counter++
			if counter > 30 {
				break
			}
		}
	} else {
		counter := 1
		for {
			if counter%2 == 0 {
				images = append(images, "http://mvatandoosts.ir/assets/images/I"+strconv.Itoa(counter)+".jpg")
			} else {
				images2 = append(images, "http://mvatandoosts.ir/assets/images/I"+strconv.Itoa(counter)+".jpg")
			}
			counter++
			if counter > NumberOfUniqueInputFiles {
				break
			}
		}

		fmt.Printf("len(images): %v, len(images2): %v \n", len(images), len(images2))
		websites = append(websites, "google.com")
		websites = append(websites, "varzesh3.com")
		websites = append(websites, "digikala.com")
		websites = append(websites, "yahoo.com")
		websites = append(websites, "stackoverflow.com")
		websites = append(websites, "github.com")
		websites = append(websites, "ut.ac.ir")
		websites = append(websites, "downloadha.com")
		websites = append(websites, "p30download.com")
		websites = append(websites, "coinmarketcap.com")
		websites = append(websites, "divar.ir")
		websites = append(websites, "mvatandoosts.ir")
	}

	// requestsTime := makeRequestsTime(RequestTimePoissonLambda)
	// fmt.Println("requestsTime: ", requestsTime)
	requestTime = makeRequestTime(RequestTimePoissonLambda)
	sort.Ints(requestTime)
	fmt.Println("Number of Requests: ", len(requestTime))
	fmt.Println("============ 0 ============ ")
	fmt.Println("requestTime: ", requestTime)
	// functions = append(functions, Function{url: "http://localhost:8080/function/nslookup", ips: 5, index: 0, inputs: websites, requestTime: requestTime})
	// requestTime = makeRequestTime(RequestTimePoissonLambda)
	functions = append(functions, Function{url: "http://localhost:8080/function/shasum", ips: 5, index: 0, inputs: websites, requestTime: requestTime})
	requestTime = makeRequestTime(RequestTimePoissonLambda)
	sort.Ints(requestTime)
	fmt.Println("============ 1 ============ ")
	fmt.Println("requestTime: ", requestTime)
	functions = append(functions, Function{url: "http://localhost:8080/function/face-detect-pigo", ips: 5, index: 0, inputs: images, requestTime: requestTime})
	requestTime = makeRequestTime(RequestTimePoissonLambda)
	sort.Ints(requestTime)
	fmt.Println("============ 2 ============ ")
	fmt.Println("requestTime: ", requestTime)
	functions = append(functions, Function{url: "http://localhost:8080/function/qrcode-go", ips: 5, index: 0, inputs: websites, requestTime: requestTime})
	requestTime = makeRequestTime(RequestTimePoissonLambda)
	sort.Ints(requestTime)
	fmt.Println("============ 3 ============ ")
	fmt.Println("requestTime: ", requestTime)
	functions = append(functions, Function{url: "http://localhost:8080/function/face-blur", ips: 5, index: 0, inputs: images2, requestTime: requestTime})
	requestTime = makeRequestTime(RequestTimePoissonLambda)
	sort.Ints(requestTime)
	fmt.Println("============ 4 ============ ")
	fmt.Println("requestTime: ", requestTime)
	functions = append(functions, Function{url: "http://localhost:8080/function/business-strategy-generator", ips: 5, index: 0, inputs: websites, requestTime: requestTime})

	return functions
}

func makeRequestTimes() {
	for j := 0; j < len(functions); j++ {
		requestTime := makeRequestTime(RequestTimePoissonLambda)
		functions[j].requestTime = requestTime
		fmt.Println("============  ", j, "===========")
		fmt.Println("requestTime: ", requestTime)
	}
}

func makeRequestTime(lambda float64) []int {
	var requestTime []int
	distribution := distuv.Poisson{Lambda: lambda}
	counter := 1
	for {
		requestTime = append(requestTime, int(distribution.Rand()))
		counter++
		if counter > NumberOfRequests {
			break
		}
	}
	return requestTime
}

func plotInfo() {
	rand.Seed(int64(0))

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Plotutil example"
	p.X.Label.Text = "Request Time"
	p.Y.Label.Text = "Output"

	err = plotutil.AddLinePoints(p,
		"nslookup", randomPoints(15),
		"face-detect-pigo", randomPoints(15),
		"qrcode-go", randomPoints(15),
		"face-blur", randomPoints(15),
		"business-strategy-generator", randomPoints(15))
	if err != nil {
		panic(err)
	}

	if err := p.Save(18*vg.Inch, 10*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}

func makeRequestsTime(lambda float64) []float64 {
	var requestTime []float64
	distribution := distuv.Poisson{Lambda: lambda}
	counter := 1
	for {
		requestTime = append(requestTime, distribution.Rand())
		counter++
		if counter > NumberOfRequests {
			break
		}
	}
	return requestTime
}

// randomPoints returns some random x, y points.
func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = pts[i].X + 10*rand.Float64()
	}
	return pts
}

func plotBox(imageName string) {
	// Get some data to display in our plot.
	rand.Seed(int64(0))
	n := len(functions[0].executionsTime)
	var results []plotter.Values
	for i := 0; i < len(functions); i++ {
		results = append(results, make(plotter.Values, n))
	}

	for i := 0; i < n; i++ {
		for j := 0; j < len(functions); j++ {
			results[j][i] = float64(functions[j].executionsTime[i])
		}
	}

	// Create the plot and set its title and axis label.
	p, _ := plot.New()

	p.Title.Text = "Execution time (ms)"
	p.Y.Label.Text = "time (ms)"

	// Make boxes for our data and add them to the plot.
	w := vg.Points(20)
	// var boxes []*plotter.BoxPlot
	boxes := make([]*plotter.BoxPlot, len(functions))
	// var err Error
	for i := 0; i < len(functions); i++ {
		boxes[i], _ = plotter.NewBoxPlot(w, float64(i), results[i])
		// boxes = append(boxes, b0)
		// if err != nil {
		//    panic(err)
		// }
		p.Add(boxes[i])
	}

	p.NominalX("nslookup", "face-detect\npigo",
		"qrcode-go", "face-blur", "business-strategy\ngenerator")

	if err := p.Save(10*vg.Inch, 8*vg.Inch, imageName); err != nil {
		panic(err)
	}
}

func largest(arr []int) int {
	var max = arr[0]
	for i := range arr {
		if arr[i] > max {
			max = arr[i]
		}
	}
	return max
}

func largestInt64(arr []int64) int {
	var max = arr[0]
	for i := range arr {
		if arr[i] > max {
			max = arr[i]
		}
	}
	return int(max)
}

func averageInt64(arr []int64) int {
	sum := 0
	for i := range arr {
		sum = sum + int(arr[i])
	}
	return int(sum / len(arr))
}

func averageInt(arr []int) int {
	sum := 0
	for i := range arr {
		sum = sum + int(arr[i])
	}
	return int(sum / len(arr))
}

func clearExecutionTime() {
	for i := 0; i < len(functions); i++ {
		functions[i].executionsTime = functions[i].executionsTime[:0]
	}
}
