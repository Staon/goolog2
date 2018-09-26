package goolog2_test

import (
	"bytes"
	. "goolog2"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

type mockApacheObject struct {
	remoteHost  string
	identity    string
	user        string
	requestTime time.Time
	method      string
	resource    string
	protocol    string
	status      int
	length      uint64
	referer     string
	agent       string
}

func (this *mockApacheObject) GetValues() (
	remoteHost string,
	identity string,
	user string,
	requestTime time.Time,
	method string,
	resource string,
	protocol string,
	status int,
	length uint64,
	referer string,
	agent string) {
	return this.remoteHost,
		this.identity,
		this.user,
		this.requestTime,
		this.method,
		this.resource,
		this.protocol,
		this.status,
		this.length,
		this.referer,
		this.agent
}

func TestApacheLogger(t *testing.T) {
	logfile := "apache.log"
	os.Remove(logfile)

	now, _ := time.Parse("2006-01-02T15:04:05 -0700 MST", "2018-08-29T22:16:26 +0200 CEST")
	timesrc := &mockTimeSource{
		now: now,
	}

	InitWithTimeSource("testlog", timesrc)
	AddApacheLogger("apache", "", MaskAll, 5, logfile, false)

	/* -- empty object */
	object := &mockApacheObject{}
	LogObject("", Error, 1, object)

	/* -- filled object */
	object = &mockApacheObject{
		remoteHost:  "127.0.0.1",
		identity:    "fooUser",
		user:        "me",
		requestTime: now.Add(10000000),
		method:      "GET",
		resource:    "/sws/my_resource.json",
		protocol:    "HTTP/1.0",
		status:      200,
		length:      56,
		referer:     "http://www.google.com/",
		agent:       "Chrome 1.0",
	}
	LogObject("", Error, 1, object)

	Destroy()

	/* -- check the log */
	expected, err1 := ioutil.ReadFile("apachelogger.log")
	current, err2 := ioutil.ReadFile(logfile)
	if err1 != nil || err2 != nil || !bytes.Equal(expected, current) {
		t.Errorf("generated log file is different!")
	}
}
