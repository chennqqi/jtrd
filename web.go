package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	mutils "github.com/malice-plugins/go-plugin-utils/utils"
)

var (
	resultExp = regexp.MustCompile(`(.*)\b\s*\(([^\)]+)\)`)
)

type Web struct {
	fileto   time.Duration
	zipto    time.Duration
	callback string
}

func (s *Web) version(c *gin.Context) {
	txt, _ := ioutil.ReadFile("/opt/jtrd/VERSION")
	c.Data(200, "", txt)
}

func (s *Web) simple(c *gin.Context) {
	var err error
	to := s.zipto
	timeout, ok := c.GetQuery("timeout")
	if ok {
		to, err = time.ParseDuration(timeout)
		if err != nil {
			to = s.fileto
		}
	}

	upf, err := c.FormFile("filename")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	src, err := upf.Open()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("open form err: %s", err.Error()))
		return
	}
	defer src.Close()
	f, err := ioutil.TempFile("/dev/shm", "shadow_")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	defer os.Remove(f.Name())
	io.Copy(f, src)
	f.Close()

	r, _ := jtrSimple(f.Name(), to)
	c.Header("Content-type", "application/json")
	s.doCallback(c, r)
	c.String(200, r)
}

func (s *Web) Run(port int) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/version", s.version)
	r.POST("/simple", s.simple)

	//r.POST("/wordlist", s.wordlist)
	//r.POST("/single", s.single)
	r.Run(fmt.Sprintf(":%d", port))
}

func (s *Web) doCallback(c *gin.Context, r string) {
	callback := c.Query("callback")
	if callback == "" {
		callback = s.callback
	}
	if callback != "" {
		go func(r string) {
			body := strings.NewReader(r)
			http.Post(callback, "application/json", body)
		}(r)
	}
}

func jtrSimple(file string, to time.Duration) (string, error) {
	fmt.Println("start simple crack")
	ctx, cancel := context.WithTimeout(context.TODO(), to)
	defer cancel()

	/*
		docker run -it -v `pwd`/yourfiletocrack:/crackme.txt adamoss/john-the-ripper /crackme.txt
	*/

	type RcdList struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	}
	var ret struct {
		Status  int       `json:"status"`
		List    []RcdList `json:"list"`
		Message string    `json:"message"`
	}

	r, err := mutils.RunCommand(ctx, "/usr/sbin/jhon", file)
	if err != nil {
		ret.Status = 500
		ret.Message = err.Error()
	} else {
		ret.Status = 200
		ret.Message = "OK"
		results := resultExp.FindAllStringSubmatch(r, -1)
		ret.List = make([]RcdList, len(results))
		for idx := 0; idx < len(results); idx++ {
			ret.List[idx].User = results[idx][1]
			ret.List[idx].Pass = results[idx][2]
		}
	}
	txt, _ := json.Marshal(ret)
	return string(txt), nil
}

func jtrWordList(dir string, to time.Duration) (string, error) {
	fmt.Println("start scan ", dir)
	ctx, cancel := context.WithTimeout(context.TODO(), to)
	defer cancel()
	return mutils.RunCommand(ctx, "/usr/sbin/jhon", "call", dir)
}
