package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/isucon/isucandar/agent"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/lo"
)

type AccessLogLine struct {
	//Time      time.Time `json:"time"`
	Level     string   `json:"level"`
	Msg       string   `json:"msg"`
	RequestID string   `json:"request_id"`
	TraceID   string   `json:"trace_id"`
	Request   Request  `json:"request"`
	Response  Response `json:"response"`
}

type Request struct {
	Method     string          `json:"method"`
	URL        *url.URL        `json:"url"`
	RequestURI string          `json:"request_uri"`
	Header     http.Header     `json:"header"`
	Cookies    []*http.Cookie  `json:"cookies"`
	Body       json.RawMessage `json:"body"`
}

type Response struct {
	Header http.Header     `json:"header"`
	Status int             `json:"status"`
	Body   json.RawMessage `json:"body"`
}

var (
	accessLog     string
	targetURLRaw  string
	ignoreHeaders string

	// trace id: agent
	agentPool = make(map[string]*agent.Agent)
)

func main() {
	flag.StringVar(&accessLog, "accesslog", "./path/to/access.log", "access log file")
	flag.StringVar(&targetURLRaw, "target-url", "http://localhost:1323/", "target url")
	flag.StringVar(&ignoreHeaders, "ignore-headers", "Date,Content-Length,Transfer-Encoding", "Comma-separated list of headers to ignoreHeaders")
	flag.Parse()
	if strings.HasSuffix(targetURLRaw, "/") {
		targetURLRaw = targetURLRaw[:len(targetURLRaw)-1]
	}

	f, err := os.Open(accessLog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v", err)
		return
	}
	defer f.Close()

	targetURL, err := url.Parse(targetURLRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse url: %v", err)
		return
	}
	if targetURL.Scheme != "http" && targetURL.Scheme != "https" {
		fmt.Fprintf(os.Stderr, "scheme must be http or https, got: %v", targetURL.Scheme)
		return
	}
	if targetURL.Host == "" {
		fmt.Fprintf(os.Stderr, "target url host must not be empty")
		return
	}
	if targetURL.Path != "" {
		fmt.Fprintf(os.Stderr, "target url path must be / or empty, got: %v", targetURL.Path)
		return
	}

	fstat, err := f.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get file stat: %v", err)
		return
	}
	fsize := fstat.Size()
	fmt.Println("file size: ", fsize)

	excludeHeaders := make(map[string]struct{})
	if ignoreHeaders != "" {
		h := strings.Split(ignoreHeaders, ",")

		for i := 0; i < len(h); i++ {
			excludeHeaders[http.CanonicalHeaderKey(strings.Trim(h[i], " "))] = struct{}{}
		}
	}

	//p := mpb.New(
	//	mpb.WithWidth(60),
	//	mpb.WithRefreshRate(180*time.Millisecond),
	//	mpb.WithOutput(os.Stderr),
	//)
	//
	//bar := p.New(fsize,
	//	mpb.BarStyle().Rbound("|"),
	//	mpb.PrependDecorators(
	//		decor.Counters(decor.SizeB1024(0), "% .2f / % .2f"),
	//	),
	//	mpb.AppendDecorators(
	//		decor.EwmaETA(decor.ET_STYLE_GO, 30),
	//		decor.Name(" ] "),
	//		decor.EwmaSpeed(decor.SizeB1024(0), "% .2f", 30),
	//	),
	//)

	//proxyReader := bar.ProxyReader(f)
	//defer proxyReader.Close()
	//scanner := bufio.NewScanner(proxyReader)
	//scanner := bufio.NewScanner(f)
	b, err := io.ReadAll(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read file: %v", err)
		return
	}
	buf := bytes.NewBuffer(b)
	lines := strings.Split(buf.String(), "\n")
	lines = lo.Filter(lines, func(line string, idx int) bool {
		return line != ""
	})
	for i, line := range lines {
		var accessLogLine AccessLogLine
		err = json.Unmarshal([]byte(line), &accessLogLine)
		if err != nil {
			fmt.Fprintf(os.Stderr, "line number: %v, trace_id: %v, request_id: %v\n", i+1, accessLogLine.TraceID, accessLogLine.RequestID)
			fmt.Fprintf(os.Stderr, "failed to unmarshal: %v", err)
			return
		}

		ag, err := GetOrNewAgent(accessLogLine.TraceID, agent.WithDefaultTransport(), agent.WithTimeout(120*time.Second))
		if err != nil {
			fmt.Fprintf(os.Stderr, "line number: %v, trace_id: %v, request_id: %v\n", i+1, accessLogLine.TraceID, accessLogLine.RequestID)
			fmt.Fprintf(os.Stderr, "failed to GetOrNewAgent: %v", err)
			return
		}
		cookieURL := &url.URL{
			Scheme: targetURL.Scheme,
			Host:   targetURL.Host + ":" + targetURL.Port(),
			Path:   "/",
		}
		ag.HttpClient.Jar.SetCookies(cookieURL, accessLogLine.Request.Cookies)
		requestBody, err := strconv.Unquote(string(accessLogLine.Request.Body))
		if err != nil {
			fmt.Fprintf(os.Stderr, "line number: %v, trace_id: %v, request_id: %v\n", i+1, accessLogLine.TraceID, accessLogLine.RequestID)
			fmt.Fprintf(os.Stderr, "failed to strconv.Unquote: %v", err)
			return
		}

		middleware.AddTrailingSlash()
		bodyBuf := bytes.NewBufferString(requestBody)
		req, err := ag.NewRequest(accessLogLine.Request.Method, targetURLRaw+accessLogLine.Request.RequestURI, bodyBuf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "line number: %v, trace_id: %v, request_id: %v\n", i+1, accessLogLine.TraceID, accessLogLine.RequestID)
			fmt.Fprintf(os.Stderr, "failed to ag.NewRequest: %v", err)
			return
		}
		res, err := ag.Do(context.Background(), req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "line number: %v, trace_id: %v, request_id: %v\n", i+1, accessLogLine.TraceID, accessLogLine.RequestID)
			fmt.Fprintf(os.Stderr, "failed to ag.Do: %v", err)
			return
		}
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "line number: %v, trace_id: %v, request_id: %v\n", i+1, accessLogLine.TraceID, accessLogLine.RequestID)
			fmt.Fprintf(os.Stderr, "failed to ReadAll response body: %v", err)
			return
		}
		res.Body.Close()

		// --------------------------------- validation ---------------------------------
		var notSameStatusCode bool
		var notSameHeaders bool
		var notSameResponseBody bool

		if res.StatusCode != accessLogLine.Response.Status {
			notSameStatusCode = true
			//fmt.Fprintf(os.Stderr, "status code is not match: %v, %v\n", res.StatusCode, accessLogLine.Response.Status)
		}
		unequalHeaders := compareHeaders(res.Header, accessLogLine.Response.Header, excludeHeaders)
		if len(unequalHeaders) != 0 {
			notSameHeaders = true
			//fmt.Fprintf(os.Stderr, "header is not match: \n%v, \n%v\n", res.Header, accessLogLine.Response.Header)
		}
		unquotedExpectedResBody, err := strconv.Unquote(string(accessLogLine.Response.Body))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to strconv.Unquote: %v", err)
		}
		notSameResponseBody, err = jsonEqual(unquotedExpectedResBody, string(resBody))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to jsonEqual: %v", err)
		}

		// --------------------------------- output ---------------------------------
		if notSameStatusCode || notSameHeaders || notSameResponseBody {
			reqInfo(i+1, accessLogLine)
		}
		if notSameStatusCode {
			vsi(accessLogLine.Response.Status, res.StatusCode, "Different status code: ")
		}
		if notSameHeaders {
			for _, h := range unequalHeaders {
				vs(accessLogLine.Response.Header.Get(h), res.Header.Get(h), "%s header different: ", green(h))
			}
		}
		if notSameResponseBody {
			if len(resBody) != len(accessLogLine.Response.Body) {
				vsi(len(accessLogLine.Response.Body), len(resBody), "Different response body length: ")
			}
			filenameExpected, err := dumpBodyToTempFile("resbody_expected", string(accessLogLine.Response.Body))
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to dump response body to temp file: %v", err)
			}
			filenameActual, err := dumpBodyToTempFile("resbody_actual", string(resBody))
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to dump response body to temp file: %v", err)
			}
			vs(filenameExpected, filenameActual, "Response body files: ")
		}
		if notSameStatusCode || notSameHeaders || notSameResponseBody {
			// TODO os.exit(1)
			return
		}
	}

	fmt.Println("ok")
	//p.Wait()
}

func reqInfo(lineNum int, line AccessLogLine) {
	filename, err := dumpBodyToTempFile("reqbody", string(line.Request.Body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to dump request body to temp file: %v", err)
	}

	fmt.Printf("line number: %v\n    trace_id: %v\n    request_id: %v\n    request url: %v %v\n    request body file: %v\n",
		lineNum,
		line.TraceID,
		line.RequestID,
		line.Request.Method,
		line.Request.RequestURI,
		filename,
	)
}
func dumpBodyToTempFile(filePrefix string, body string) (tmpFilename string, err error) {
	unquotedBody, err := strconv.Unquote(body)
	if err != nil {
		return "", err
	}
	f, err := os.CreateTemp("", filePrefix)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := &bytes.Buffer{}
	if err := json.Indent(buf, []byte(unquotedBody), "", "  "); err != nil {
		return "", err
	}
	_, err = f.Write(buf.Bytes())
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}

func GetOrNewAgent(traceId string, opts ...agent.AgentOption) (*agent.Agent, error) {
	if traceId == "" {
		return agent.NewAgent(opts...)
	}
	if a, ok := agentPool[traceId]; ok {
		return a, nil
	}
	ag, err := agent.NewAgent(opts...)
	if err != nil {
		return nil, err
	}
	agentPool[traceId] = ag
	return ag, nil
}

var mono = false
var notsame = false

// ANSI escape functions and print helpers
func on(i int, s string) string {
	if mono {
		return fmt.Sprintf("%d: %s", i+1, s)
	}
	return fmt.Sprintf("\x1b[3%dm%s\x1b[0m", i*3+1, s)
}
func oni(i, d int) string {
	return on(i, fmt.Sprintf("%d", d))
}
func green(s string) string {
	if mono {
		return fmt.Sprintf("'%s'", s)
	}
	return fmt.Sprintf("\x1b[32m%s\x1b[0m", s)
}
func vs(a, b string, f string, v ...interface{}) bool {
	notsame = a != b
	if notsame {
		s := fmt.Sprintf(f, v...)
		fmt.Printf("%s\n    %s\n    %s\n", s, on(0, a), on(1, b))
	}
	return notsame
}
func vsi(a, b int, f string, v ...interface{}) bool {
	notsame = a != b
	if notsame {
		s := fmt.Sprintf(f, v...)
		fmt.Printf("%s\n    %s\n    %s\n", s, oni(0, a), oni(1, b))
	}
	return notsame
}

func compareHeaders(a, b http.Header, excludeHeaders map[string]struct{}) (unequalHeaders []string) {
	for key, valA := range a {
		if _, exclude := excludeHeaders[key]; exclude {
			continue
		}
		valB, exists := b[key]
		if !exists || !slicesEqual(valA, valB) {
			unequalHeaders = append(unequalHeaders, key)
			continue
		}
	}

	for key := range b {
		if _, exclude := excludeHeaders[key]; exclude {
			continue
		}
		if _, exists := a[key]; !exists {
			unequalHeaders = append(unequalHeaders, key)
		}
	}

	return unequalHeaders
}

// 順序を無視してスライスの要素を比較
func slicesEqual(sliceA, sliceB []string) bool {
	if len(sliceA) != len(sliceB) {
		return false
	}

	copyA := append([]string{}, sliceA...)
	copyB := append([]string{}, sliceB...)
	sort.Strings(copyA)
	sort.Strings(copyB)

	for i := range copyA {
		if copyA[i] != copyB[i] {
			return false
		}
	}

	return true
}

// 順序を無視してJSONを比較
func jsonEqual(a, b string) (bool, error) {
	var objA interface{}
	var objB interface{}

	if err := json.Unmarshal([]byte(a), &objA); err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(b), &objB); err != nil {
		return false, err
	}

	return deepEqual(objA, objB), nil
}

// 順序を無視して要素を比較
func deepEqual(a, b interface{}) bool {
	switch aVal := a.(type) {
	case []interface{}:
		bVal, ok := b.([]interface{})
		if !ok || len(aVal) != len(bVal) {
			return false
		}

		// Create a map to count occurrences
		counts := make(map[interface{}]int)
		for _, val := range aVal {
			counts[val]++
		}
		for _, val := range bVal {
			counts[val]--
		}

		// If any count is non-zero, slices are not equal
		for _, count := range counts {
			if count != 0 {
				return false
			}
		}
		return true
	case map[string]interface{}:
		bVal, ok := b.(map[string]interface{})
		if !ok || len(aVal) != len(bVal) {
			return false
		}

		for key, valA := range aVal {
			valB, exists := bVal[key]
			if !exists || !deepEqual(valA, valB) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}
