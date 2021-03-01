package main

import(
	"fmt"
	"strconv"
	"net/http"
	"bytes"
	"io/ioutil"
	"time"
	"math/rand"
	"benchmark/internal/distuv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"sort"
	// "gonum.org/v1/gonum/stat/distuv"
)

const (
	TestTime = 5
)

type Function struct {
	url string
	ips int  // invocations per second (ips) 
	executionsTime []int64
    requestTime []int
	index int
	inputs []string
}

var functions []Function

func main() {
   
    // executionTime := sendPostRequest(functions[0], "google.com")
	functions = initialize()
	
	coutner := 0
	for {
		for i, f := range functions {
             go sendPostRequest(i, f.url, f.inputs[coutner%len(f.inputs)])
		}
		time.Sleep(500*time.Millisecond)
		coutner++
		if coutner > 18 {
			break
		}
	}
	time.Sleep(5000*time.Millisecond)
	for _, f := range functions {
		fmt.Println("Execution time:", f.executionsTime)
    }
	plotInfo()
	// fmt.Println("executionTime:", executionTime.Milliseconds())

}

func sendPostRequest(index int, url string, data string) time.Duration  {
    fmt.Println("sendPostRequest URL:>", url)

    var jsonStr = []byte(data)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
	t1 := time.Now()
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
	duration := time.Since(t1)
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    // fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", len(body), ", execution time:", duration.Milliseconds())
	functions[index].executionsTime = append(functions[index].executionsTime, duration.Milliseconds())
	return duration
}

func initialize() []Function {
	var images []string
	var websites []string
	var functions []Function
    var requestTime []int
	

	counter := 1
	for {
		images = append(images, "http://mvatandoosts.ir/assets/images/I"+strconv.Itoa(counter)+".jpg" )
		counter++
		if counter > 10 {
			break
		}
	}

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

    // requestTime = append(requestTime, 1)
	// requestTime = append(requestTime, 2)
	// requestTime = append(requestTime, 3)
	// requestTime = append(requestTime, 4)
	// requestTime = append(requestTime, 5)
	// requestTime = append(requestTime, 6)
	// requestTime = append(requestTime, 7)
	// requestTime = append(requestTime, 8)
	// requestTime = append(requestTime, 9)
	// requestTime = append(requestTime, 10)

	requestTime = makeRequestTime(10)
	sort.Ints(requestTime)
	fmt.Println("requestTime: ", requestTime)
	functions = append(functions, Function{url: "http://localhost:8080/function/nslookup", ips: 5, index: 0, inputs: websites, requestTime: requestTime} )
	requestTime = makeRequestTime(10)
	sort.Ints(requestTime)
	fmt.Println("requestTime: ", requestTime)
	functions = append(functions, Function{url: "http://localhost:8080/function/face-detect-pigo", ips: 5, index: 0, inputs: images, requestTime: requestTime} )
	requestTime = makeRequestTime(10)
	sort.Ints(requestTime)
	fmt.Println("requestTime: ", requestTime)
	functions = append(functions, Function{url: "http://localhost:8080/function/qrcode-go", ips: 5, index: 0, inputs: websites, requestTime: requestTime} )
	requestTime = makeRequestTime(10)
	sort.Ints(requestTime)
	fmt.Println("requestTime: ", requestTime)
	functions = append(functions, Function{url: "http://localhost:8080/function/face-blur", ips: 5, index: 0, inputs: images, requestTime: requestTime} )
	requestTime = makeRequestTime(10)
	sort.Ints(requestTime)
	fmt.Println("requestTime: ", requestTime)
	functions = append(functions, Function{url: "http://localhost:8080/function/business-strategy-generator", ips: 5, index: 0, inputs: websites, requestTime: requestTime} )

    
	return functions
}


func makeRequestTime(lambda float64) []int  {
	var requestTime []int
	distribution := distuv.Poisson{Lambda: lambda}
	counter := 1
	for {
		requestTime = append(requestTime, int(distribution.Rand()))
		counter++
		if counter > 24 {
			break
		}
	} 
	return requestTime
}



func plotInfo()  {
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

	// Save the plot to a PNG file.
	if err := p.Save(18*vg.Inch, 10*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
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
    