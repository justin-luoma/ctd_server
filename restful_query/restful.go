package restful_query

import (
	"github.com/valyala/fasthttp"
	"errors"
	"strconv"
	"time"
	"github.com/golang/glog"
	"flag"
)

func init() {
	flag.Parse()
}

func Get(url string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.Add("User-Agent", "CTD-Query")

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	if err != nil {
		glog.Error ("Source: restful->Get->client", err)
		return nil, err
	}
	if resp.Header.StatusCode() != 200 {
		if resp.Header.StatusCode() == 429 {
			glog.Warningln("api limit exceeded, trying again")
			time.Sleep(time.Second)
			return Get(url)
		} else {
			glog.Errorln("url: " + url + "status code: " + strconv.Itoa(resp.Header.StatusCode()))
			return nil, errors.New("url: " + string(req.URI().FullURI()) + " status: " + strconv.Itoa(resp.Header.StatusCode()))
		}

	}

	glog.V(2).Infoln("queried: " + string(req.URI().FullURI()))

	fasthttp.ReleaseRequest(req)

	bodyBytes := resp.Body()

	fasthttp.ReleaseResponse(resp)

	return bodyBytes, nil
}
