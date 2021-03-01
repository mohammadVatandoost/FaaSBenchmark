package main

import(
	"fmt"
	"strconv"
	"net/http"
	"bytes"
	"io/ioutil"
	"time"
)

const (
	TestTime = 5
)

type Function struct {
	url string
	ips int  // invocations per second (ips) 
	executionsTime []int
    requestTime []int
	index int
}

func main() {
   
    executionTime := sendPostRequest(functions[0], "google.com")
	fmt.Println("executionTime:", executionTime.Milliseconds())

}

func sendPostRequest(url string, data string) time.Duration  {
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
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
	return duration
}

func init()  {
	var images []string
	var websites []string
	var functions []Function


	functions = append(functions, Function{url: "http://localhost:8080/function/nslookup", ips: 5, index: 0} )
	functions = append(functions, Function{url: "http://localhost:8080/function/face-detect-pigo", ips: 5, index: 0} )
	functions = append(functions, Function{url: "http://localhost:8080/function/qrcode-go", ips: 5, index: 0} )
	functions = append(functions, Function{url: "http://localhost:8080/function/face-blur", ips: 5, index: 0} )
	functions = append(functions, Function{url: "http://localhost:8080/function/business-strategy-generator", ips: 5, index: 0} )

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
    
}