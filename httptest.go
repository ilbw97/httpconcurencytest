package main

import (
	"crypto/tls"
	"regexp"

	"flag"
	"net/http"
	"strings"
	"sync"

	"github.com/ilbw97/debuglog"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger = &logrus.Logger{}

func curl(host string, method string) {
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

func main() {

	fprotocol := flag.String("protocol", "", "EX) http/https")
	fhost := flag.String("host", "", "EX) wordpress.jam10000bo.com")
	fmethod := flag.String("method", "", "EX) get/put/post/update")
	fpath := flag.String("path", "", "EX) /cloud2team")
	fport := flag.String("port", "", "EX) 8099")
	fcount := flag.Int("count", 0, "EX) 100")

	flag.Parse()
	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	if *fprotocol == "" || *fhost == "" || *fmethod == "" || *fpath == "" || *fport == "" || *fcount == 0 {
		flag.Usage()
		return
	}

	method := strings.ToUpper(*fmethod)
	switch method {
	case "GET":
	case "PUT":
	case "POST":
	case "UPDATE":
	default:
		flag.Usage()
	}

	wg.Add(*fcount)

	reqpath := *fprotocol + "://" + *fhost + ":" + *fport + *fpath
	logregex := regexp.MustCompile(`[\{\}\[\]\/?.,;:|\)*~!^\_+<>@\#$%&\\\=\(\'\"\n\r]+`)
	logname := logregex.ReplaceAllString("httptest_"+reqpath, "_")
	log = debuglog.DebugLogInit(logname)

	for i := 0; i <= *fcount; i++ {
		log.Infof("Trying to %s %s, No.%d", method, reqpath, i)
		// go func() {
		go curl(reqpath, method)
		// defer wg.Done()
		// }()
	}
	wg.Wait()
}

var wg sync.WaitGroup
