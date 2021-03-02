package main

import(
	"fmt"
	"strconv"
	"net/http"
	"bytes"
	"io/ioutil"
	"time"
	"math/rand"
	// "benchmark/internal/distuv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"sort"
	// "image/color"


	"gonum.org/v1/gonum/stat/distuv"
)

const (
	TestTime = 5
	RequestTimePoissonLambda = 30
	NumberOfRequests = 40
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
					go sendPostRequest(i, f.url, f.inputs[rand.Intn(len(f.inputs) - 0) + 0])
				}
			}
		}
		time.Sleep(100*time.Millisecond)
		counter++
		if counter > maxTime {
			break
		}
	}
	time.Sleep(15000*time.Millisecond)
	for _, f := range functions {
		fmt.Println("Execution time:", f.executionsTime)
    }
	plotBox();
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

    // fmt.Println("response Status:", resp.Status)
    // fmt.Println("response Headers:", resp.Header)
    _, _ = ioutil.ReadAll(resp.Body)
    // fmt.Println("response Body:", len(body), ", execution time:", duration.Milliseconds())
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
		if counter > 30 {
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



	requestTime = makeRequestTime(RequestTimePoissonLambda)
	sort.Ints(requestTime)
	fmt.Println("Number of Requests: ", len(requestTime))
	fmt.Println("requestTime: ", requestTime)
	functions = append(functions, Function{url: "http://localhost:8080/function/nslookup", ips: 5, index: 0, inputs: websites, requestTime: requestTime} )
	requestTime = makeRequestTime(RequestTimePoissonLambda)
	sort.Ints(requestTime)
	fmt.Println("requestTime: ", requestTime)
	functions = append(functions, Function{url: "http://localhost:8080/function/face-detect-pigo", ips: 5, index: 0, inputs: images, requestTime: requestTime} )
	requestTime = makeRequestTime(RequestTimePoissonLambda)
	sort.Ints(requestTime)
	fmt.Println("requestTime: ", requestTime)
	functions = append(functions, Function{url: "http://localhost:8080/function/qrcode-go", ips: 5, index: 0, inputs: websites, requestTime: requestTime} )
	requestTime = makeRequestTime(RequestTimePoissonLambda)
	sort.Ints(requestTime)
	fmt.Println("requestTime: ", requestTime)
	functions = append(functions, Function{url: "http://localhost:8080/function/face-blur", ips: 5, index: 0, inputs: images, requestTime: requestTime} )
	requestTime = makeRequestTime(RequestTimePoissonLambda)
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
		if counter > NumberOfRequests {
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

func plotBox()  {
	// Get some data to display in our plot.
	rand.Seed(int64(0))
	n := len(functions[0].executionsTime)
	var results []plotter.Values
	for i := 0; i < len(functions); i++ {
		results = append(results, make(plotter.Values, n))
	}
	// f1 := make(plotter.Values, n)
	// f2 := make(plotter.Values, n)
	// f3 := make(plotter.Values, n)
	// f4 := make(plotter.Values, n)
	// f5 := make(plotter.Values, n)
	// fmt.Println(len(functions[0].requestTime))
	// fmt.Println(len(functions[1].requestTime))
	// fmt.Println(len(functions[2].requestTime))
	// fmt.Println(n)
	for i := 0; i < n; i++ {
		for j := 0; j < len(functions); j++ {
			results[j][i] = float64(functions[j].executionsTime[i])
		}
	}
	// fmt.Println(f1[n-1])
	// fmt.Println(f2[n-1])
	// fmt.Println(f3[n-1])

	

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
	// b0, err := plotter.NewBoxPlot(w, 0, f1)
    //     // b0.FillColor = color.RGBA{127, 188, 165, 1}
	// if err != nil {
	// 	panic(err)
	// }
	// b1, err := plotter.NewBoxPlot(w, 1, f2)
    //     // b1.FillColor = color.RGBA{127, 188, 165, 1}
	// if err != nil {
	// 	panic(err)
	// }
	// b2, err := plotter.NewBoxPlot(w, 2, f3)
    //     // b2.FillColor = color.RGBA{127, 188, 165, 1}
	// if err != nil {
	// 	panic(err)
	// }
	// p.Add(b0, b1, b2)

	// Set the X axis of the plot to nominal with
	// the given names for x=0, x=1 and x=2.
	p.NominalX("nslookup", "face-detect\npigo",
		"qrcode-go", "face-blur", "business-strategy\ngenerator")

	if err := p.Save(10*vg.Inch, 8*vg.Inch, "boxplot.png"); err != nil {
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
 
    