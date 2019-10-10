package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	//"os"
	"regexp"
	"strconv"
	"strings"
	//"sync"

	"../config/db"
	"../model"

	log "github.com/albrow/prtty"
	"github.com/elazarl/goproxy"
	//"github.com/davecgh/go-spew/spew"
)

const (
	outdir = `../public/system_new/src/netent`
)

var (
	//id      Id
	exclude = regexp.MustCompile(`(?iU)(^redirect$|^localhost$|icloud|google|gstatic|yandex|apple)`)
	iframe  = regexp.MustCompile(`(?iUm)(?:width="(?P<w>[0-9]+)")|(?:height="(?P<h>[0-9]+)")`)
)

func main() {
	//id = Id{I: 1}
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8081", "proxy listen address")
	flag.Parse()

	setCA(caCert, caKey)

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *verbose

	// request handlers
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		//spew.Dump(req.Header)
		//ctx.UserData
		req1 := OnRequest(req)
		if req1 != nil {
			return nil, req1
		}
		return req, nil
	})

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)


	// responce handlers
	/*
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		u := *ctx.Req.URL

		/*
		if u.Hostname() == "www.askgamblers.com" {
			for k, v := range ctx.Req.Header {
				for _, v2 := range v {
					if k == "Cookie" {
						cookies := strings.Split(v2, "; ")
						for _, cookie := range cookies {
							if strings.HasPrefix(cookie, "GameId=") {
								log.Default.Println(cookie)
							}
						}
					}
				}
			}
			//spew.Dump(ctx.Req.Header)
		}


			if !exclude.MatchString(u.Hostname()) {
				//if u.Hostname() != "redirect" && u.Hostname() != "localhost" && !strings.Contains(u.Hostname(), "icloud") && !strings.Contains(u.Hostname(), "google") && !strings.Contains(u.Hostname(), "yandex") {
				parts := strings.Split(u.RequestURI(), "/")
				dir := fmt.Sprintf(`%s/%s`, outdir, u.Hostname())
				file := ""
				for k, v := range parts {
					if v != "" && k < len(parts)-1 {
						dir = dir + `/` + v
					}
					if k == len(parts)-1 {
						file = v
					}
				}

				parts = strings.Split(file, "?")
				file = parts[0]

				if file == "" {
					file = "index.html"
				}

				os.MkdirAll(dir, 0755)

				f, err := os.Create(dir + `/` + file + `.header`)
				if err != nil {
					log.Error.Println(err)
				}

				//log.Printf("Saving %s/%s.header", dir, file)
				for k, v := range resp.Header {
					for _, v2 := range v {
						f.WriteString(fmt.Sprintf("%s: %s\n", k, v2))
					}
				}

				//log.Printf("Saving %s/%s", dir, file)
				resp.Body = NewTeeReadCloser(resp.Body, NewFileStream(dir+`/`+file))

			}

		return resp
	})
	*/

	err := http.ListenAndServe(*addr, proxy); if err != nil {
		log.Error.Fatal(err)
	}
}

func OnRequest(req *http.Request) *http.Response {
	u := *req.URL

	if u.Hostname() == "redirect" && strings.HasPrefix(u.RequestURI(), `/?id=`) {
		parts := strings.Split(u.RequestURI(), "=")

		n, _ := strconv.Atoi(parts[len(parts)-1])
		/*
		if uint(n) != id.Get() {
			id.Set(uint(n))
		}
		*/
		game := model.Game{}
		db.DB.First(&game, n)

		return GetLoader(req, game.Object)
	}

	return nil
}

func GetLoader(req *http.Request, content string) *http.Response {
	content = regexp.MustCompile(`(?iUm)(?:width="([0-9]+)")`).ReplaceAllStringFunc(content, func(s string) string {
		return `width="1024"`
	})
	content = regexp.MustCompile(`(?iUm)(?:height="([0-9]+)")`).ReplaceAllStringFunc(content, func(s string) string {
		return `height="768"`
	})
	//match := iframe.ReplaceAllString(content, `${w}`)
	//spew.Dump(content)

	html := fmt.Sprintf(
		`<!DOCTYPE html>
<html>
<head>
  <title></title>
</head>
<body>
  <div id="content">%s</div>
</body>
</html>`, content)

	resp := http.Response{
		StatusCode:    200,
		ProtoMajor:    1,
		ProtoMinor:    1,
		Request:       req,
		Body:          ioutil.NopCloser(bytes.NewBuffer([]byte(html))),
		ContentLength: int64(len(html)),
	}

	resp.TransferEncoding = req.TransferEncoding
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Type", "text/html; charset=utf-8")

	return &resp
}

/*
type Id struct {
	sync.RWMutex
	I uint
}

func (c *Id) Set(i uint) {
	c.Lock()
	c.I = i
	c.Unlock()
}

func (c *Id) Get() uint {
	c.RLock()
	defer c.RUnlock()
	return c.I
}
*/
