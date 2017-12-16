package restful_query

import (
	"github.com/valyala/fasthttp"
	"log"
	"errors"
)

func Get(url string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.Add("User-Agent", "CTD-Query")

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	if err != nil {
		log.Fatal(err)
	}
	if resp.Header.StatusCode() != 200 {
		functionError := string(resp.Header.StatusCode()) + string(req.URI().FullURI())
		return nil, errors.New(functionError)
	}

	fasthttp.ReleaseRequest(req)

	bodyBytes := resp.Body()

	fasthttp.ReleaseResponse(resp)

	return bodyBytes, nil
}
