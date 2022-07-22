package main

import (
	"crypto/tls"
	"regexp"
	"strconv"
	"time"

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

	log.Infof("SUCCESS TO POST %v, STATUS IS : %v", host, res.Status)

	return

}

func checkFlag() *info {

	option := new(info)

	option.protocol = flag.String("protocol", "http", "YOU CAN ENTER http / https. DEFAULT IS http")
	option.host = flag.String("host", "", "EX) wordpress.jam10000bo.com")
	option.method = flag.String("method", "get", "YOU CAN ENTER 'get / put / post / update'. DEFAULT IS GET")
	option.path = flag.String("path", "/", "EX) /cloud2team. DEFAULT IS /")
	option.port = flag.String("port", "80", "YOU CAN ENTER ONLY POSITIVE NUMBER. DEFAULT IS 80")
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

	port, err := strconv.Atoi(*option.port)
	if err != nil || port <= 0 {
		flag.Usage()
		return nil
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

	reqpath := *option.protocol + "://" + *option.host + ":" + *option.port + *option.path
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
}
