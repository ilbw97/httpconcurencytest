package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"

	"flag"
	"net/http"
	"strings"
	"sync"

	"github.com/google/brotli/go/cbrotli"
	"github.com/ilbw97/debuglog"

	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger = &logrus.Logger{}
	wg  sync.WaitGroup
	re  result = result{}
)

type info struct {
	protocol  *string
	host      *string
	method    *string
	path      *string
	port      *string
	count     *int
	loop      *int
	interval  *int
	directory *bool
}

type result struct {
	success int
	fail    int
	blocked int
}

func initlog(reqpath string, directory bool) {
	logregex := regexp.MustCompile(`[\{\}\[\]\/?.,;:|\)*~!^\_+<>@\#$%&\\\=\(\'\"\n\r]+`)
	logname := logregex.ReplaceAllString("http_concurencytest_"+reqpath, "_")
	log = debuglog.DebugLogInit(logname, directory)
}

func makeRequest(host string, method string) {

	defer wg.Done()

	req, err := http.NewRequest(method, host, nil)
	if err != nil {
		log.Errorf("httpRequest Error : %v", err)
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{Transport: tr}

	res, err := client.Do(req)
	if err != nil {
		log.Errorf("client.Do Error : %v", err)
		return
	}

	defer res.Body.Close()
	blocked := bodyCheck(res, maxsize)

	if res.StatusCode == http.StatusOK {
		if blocked {
			re.blocked += 1
			log.Infof("SUCCESS TO MAKE POST REQUEST TO %v, STATUS IS : %v BUT BLOCKED", host, res.Status)
		}
		log.Infof("SUCCESS TO POST %v, STATUS IS : %v", host, res.Status)
		re.success += 1
	} else {
		if blocked {
			re.blocked += 1
		}
		log.Infof("SUCCESS TO MAKE POST REQUEST TO %v, BUT STATUS IS %v AND BLOCKED", host, res.Status)
		re.fail += 1
	}
	return

}

var maxsize int64 = 1000 * 1000 * 1

func bodyCheck(res *http.Response, maxsize int64) bool {
	reader := io.LimitReader(ioutil.NopCloser(res.Body), maxsize)
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Error(err)
		return false
	}

	newreader := io.MultiReader(bytes.NewBuffer(body), res.Body)
	res.Body = ioutil.NopCloser(newreader)

	if len(body) == int(maxsize) {
		log.Infof("max size over %d", maxsize)
		return false
	}

	contentEncoding := res.Header.Get("Content-encoding")

	var decompresser io.ReadCloser

	if contentEncoding == "gzip" {
		decompresser, err = gzip.NewReader(bytes.NewBuffer(body))
		if err != nil {
			log.Error(err)
			return false
		}

	} else if contentEncoding == "br" {
		decompresser = cbrotli.NewReader(bytes.NewBuffer(body))
	} else if contentEncoding == "deflate" {
		decompresser = flate.NewReader(bytes.NewBuffer(body))
	} else {
		decompresser = nil
	}

	var plainBody []byte

	if decompresser != nil {
		plainBody, err = ioutil.ReadAll(decompresser)
		decompresser.Close()
		if err != nil {
			log.Error(err)
			return false
		}
	} else {
		plainBody = body
	}
	if bytes.ContainsAny(plainBody, "Block") {
		log.Infof("blocked!!!")
		return true
	}
	return false
}

func checkFlag() *info {

	option := new(info)

	option.protocol = flag.String("protocol", "http", "YOU CAN ENTER http / https. DEFAULT IS http")
	option.host = flag.String("host", "", "EX) wordpress.jam10000bo.com")
	option.method = flag.String("method", "get", "YOU CAN ENTER 'get / put / post / update'. DEFAULT IS GET")
	option.path = flag.String("path", "/", "EX) /cloud2team. DEFAULT IS /")
	option.port = flag.String("port", "", "YOU CAN ENTER ONLY POSITIVE NUMBER. DEFAULT IS nothing")
	option.count = flag.Int("count", 1, "YOU CAN ENTER ONLY POSITIVE NUMBER. DEFUALT IS 1.")
	option.loop = flag.Int("loop", 1, "YOU CAN ENTER ONLY POSITIVE NUMBER. DEFUALT IS 1.")
	option.interval = flag.Int("interval", 0, "YOU CAN ENTER ONLY POSITIVE NUMBER. DEFUALT IS 1.")
	option.directory = flag.Bool("directory", true, "YOU CAN ENTRER ONLY TRUE OR FALSE. DEFAULT IS TRUE")

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		return nil
	}

	if *option.host == "" {
		flag.Usage()
		return nil
	}

	if *option.port != "" {
		port, err := strconv.Atoi(*option.port)
		if err != nil || port <= 0 {
			flag.Usage()
			return nil
		}
	}

	if *option.count <= 0 {
		flag.Usage()
		return nil
	}

	if *option.loop < 1 {
		flag.Usage()
		return nil
	}

	if *option.interval < 1 {
		flag.Usage()
		return nil
	}

	switch *option.protocol {
	case "http":
	case "https":
	default:
		flag.Usage()
		return nil
	}

	method := strings.ToUpper(*option.method)
	switch method {
	case "GET":
	case "PUT":
	case "POST":
	case "UPDATE":
	default:
		flag.Usage()
		return nil
	}

	return option
}

func main() {

	option := checkFlag()
	if option == nil {
		return
	}

	var reqpath string
	if *option.port == "" {
		reqpath = *option.protocol + "://" + *option.host + *option.path
	} else {
		reqpath = *option.protocol + "://" + *option.host + ":" + *option.port + *option.path
	}

	initlog(reqpath, *option.directory)

	wg.Add(*option.count * *option.loop)

	var cnt int
	for i := 0; i < *option.loop; i++ {
		for j := 0; j < *option.count; j++ {
			cnt += 1
			log.Infof("Trying to %s %s, No.%d", *option.method, reqpath, cnt)
			go makeRequest(reqpath, *option.method)
		}
		log.Infof("WILL SEND %d REQUEST TO %s %d Seconds Later", *option.count, reqpath, *option.interval)
		time.Sleep(time.Second * time.Duration(*option.interval))
	}

	wg.Wait()
	log.Infof("SUCCESS TO SEND %d REQUEST TO %s", cnt, reqpath)
	log.Infof("STATUS OK : %d, NOT OK : %d, BLOCKED COUNT : %d", re.success, re.fail, re.blocked)

}
