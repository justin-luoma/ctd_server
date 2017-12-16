package restful_query

import (
	"github.com/valyala/fasthttp"
	"log"
	"errors"
	"strconv"
	"time"
)

func Get(url string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.Add("User-Agent", "CTD-Query")

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	if err != nil {
		log.Fatalln("Source: restful->Get->client")
		log.Fatal(err)
		return nil, err
	}
	if resp.Header.StatusCode() != 200 {
		if resp.Header.StatusCode() == 429 {
			log.Fatalln("api limit exceeded, trying again")
			time.Sleep(time.Second)
			return Get(url)
		} else {
			log.Fatalln("url: " + url + "status code: " + strconv.Itoa(resp.Header.StatusCode()))
			return nil, errors.New("url: " + string(req.URI().FullURI()) + " status: " + strconv.Itoa(resp.Header.StatusCode()))
		}

	}

	log.Println("queried: " + string(req.URI().FullURI()))

	fasthttp.ReleaseRequest(req)

	bodyBytes := resp.Body()

	fasthttp.ReleaseResponse(resp)

	return bodyBytes, nil
}
