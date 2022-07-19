package main

import (
	"crypto/tls"
	"regexp"
	"strconv"

	"flag"
	"net/http"
	"strings"
	"sync"

	"github.com/ilbw97/debuglog"

	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger = &logrus.Logger{}
	wg  sync.WaitGroup
)

func initlog(reqpath string) {
	logregex := regexp.MustCompile(`[\{\}\[\]\/?.,;:|\)*~!^\_+<>@\#$%&\\\=\(\'\"\n\r]+`)
	logname := logregex.ReplaceAllString("httptest_"+reqpath, "_")
	log = debuglog.DebugLogInit(logname)
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

	log.Infof("SUCCESS TO POST %v", host)

	return

}

type info struct {
	protocol *string
	host     *string
	method   *string
	path     *string
	port     *string
	count    *int
}

func checkFlag() *info {
	n := new(info)

	n.protocol = flag.String("protocol", "http", "YOU CAN ENTER http / https. DEFAULT IS http")
	n.host = flag.String("host", "", "EX) wordpress.jam10000bo.com")
	n.method = flag.String("method", "get", "YOU CAN ENTER 'get / put / post / update'. DEFAULT IS GET")
	n.path = flag.String("path", "/", "EX) /cloud2team. DEFAULT IS /")
	n.port = flag.String("port", "80", "YOU CAN ENTER ONLY POSITIVE NUMBER. DEFAULT IS 80")
	n.count = flag.Int("count", 1, "YOU CAN ENTER ONLY POSITIVE NUMBER. DEFUALT IS 1.")

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		return nil
	}

	if *n.host == "" {
		flag.Usage()
		return nil
	}

	port, err := strconv.Atoi(*n.port)
	if err != nil || port <= 0 {
		flag.Usage()
		return nil
	}

	if *n.count <= 0 {
		flag.Usage()
		return nil
	}

	switch *n.protocol {
	case "http":
	case "https":
	default:
		flag.Usage()
		return nil
	}

	method := strings.ToUpper(*n.method)
	switch method {
	case "GET":
	case "PUT":
	case "POST":
	case "UPDATE":
	default:
		flag.Usage()
		return nil
	}

	return n
}

func main() {

	n := checkFlag()
	if n == nil {
		return
	}

	reqpath := *n.protocol + "://" + *n.host + ":" + *n.port + *n.path
	initlog(reqpath)

	wg.Add(*n.count)

	for i := 0; i < *n.count; i++ {
		log.Infof("Trying to %s %s, No.%d", *n.method, reqpath, i)
		go makeRequest(reqpath, *n.method)
	}

	wg.Wait()
}
