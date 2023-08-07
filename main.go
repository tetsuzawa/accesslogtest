package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	ignoreBodies  string

	// trace id: agent
	agentPool = make(map[string]*agent.Agent)
)

func main() {
	flag.StringVar(&accessLog, "accesslog", "./path/to/access.log", "access log file")
	flag.StringVar(&targetURLRaw, "target-url", "http://localhost:1323/", "target url")
	flag.StringVar(&ignoreHeaders, "ignore-headers", "Date,Content-Length,Transfer-Encoding,Connection,Set-Cookie,Server,Vary,Content-Encoding", "Comma-separated list of headers to ignore")
	flag.StringVar(&ignoreBodies, "ignore-bodies", "sessionId", "Comma-separated list of bodies to ignore")
	flag.Parse()
	if strings.HasSuffix(targetURLRaw, "/") {
		targetURLRaw = targetURLRaw[:len(targetURLRaw)-1]
	}

	f, err := os.Open(accessLog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v\n", err)
		return
	}
	defer f.Close()

	targetURL, err := url.Parse(targetURLRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse url: %v\n", err)
		return
	}
	if targetURL.Scheme != "http" && targetURL.Scheme != "https" {
		fmt.Fprintf(os.Stderr, "scheme must be http or https, got: %v\n", targetURL.Scheme)
		return
	}
	if targetURL.Host == "" {
		fmt.Fprintf(os.Stderr, "target url host must not be empty")
		return
	}
	if targetURL.Path != "" {
		fmt.Fprintf(os.Stderr, "target url path must be / or empty, got: %v\n", targetURL.Path)
		return
	}

	fstat, err := f.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get file stat: %v\n", err)
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
	excludeBodies := make(map[string]struct{})
	if ignoreBodies != "" {
		h := strings.Split(ignoreBodies, ",")

		for i := 0; i < len(h); i++ {
			excludeBodies[strings.Trim(h[i], " ")] = struct{}{}
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
		fmt.Fprintf(os.Stderr, "failed to read file: %v\n", err)
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
			fmt.Fprintf(os.Stderr, "failed to unmarshal: %v\n", err)
			return
		}

		ag, err := GetOrNewAgent(accessLogLine.TraceID, agent.WithDefaultTransport(), agent.WithTimeout(120*time.Second))
		if err != nil {
			fmt.Fprintf(os.Stderr, "line number: %v, trace_id: %v, request_id: %v\n", i+1, accessLogLine.TraceID, accessLogLine.RequestID)
			fmt.Fprintf(os.Stderr, "failed to GetOrNewAgent: %v\n", err)
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
			fmt.Fprintf(os.Stderr, "failed to strconv.Unquote: %v\n", err)
			return
		}

		bodyBuf := bytes.NewBufferString(requestBody)
		req, err := ag.NewRequest(accessLogLine.Request.Method, targetURLRaw+accessLogLine.Request.RequestURI, bodyBuf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "line number: %v, trace_id: %v, request_id: %v\n", i+1, accessLogLine.TraceID, accessLogLine.RequestID)
			fmt.Fprintf(os.Stderr, "failed to ag.NewRequest: %v\n", err)
			return
		}
		for k, vs := range accessLogLine.Request.Header {
			defaultValues := req.Header.Values(k)
			for _, v := range vs {
				if !lo.Contains(defaultValues, v) {
					req.Header.Add(k, v)
				}
			}
		}
		res, err := ag.Do(context.Background(), req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "line number: %v, trace_id: %v, request_id: %v\n", i+1, accessLogLine.TraceID, accessLogLine.RequestID)
			fmt.Fprintf(os.Stderr, "failed to ag.Do: %v\n", err)
			return
		}
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "line number: %v, trace_id: %v, request_id: %v\n", i+1, accessLogLine.TraceID, accessLogLine.RequestID)
			fmt.Fprintf(os.Stderr, "failed to ReadAll response body: %v\n", err)
			return
		}
		res.Body.Close()

		// --------------------------------- validation ---------------------------------
		var isNotSameStatusCode bool
		var isNotSameHeaders bool
		var isNotSameResponseBody bool

		if res.StatusCode != accessLogLine.Response.Status {
			isNotSameStatusCode = true
			//fmt.Fprintf(os.Stderr, "status code is not match: %v, %v\n", res.StatusCode, accessLogLine.Response.Status)
		}
		unequalHeaders := compareHeaders(res.Header, accessLogLine.Response.Header, excludeHeaders)
		if len(unequalHeaders) != 0 {
			isNotSameHeaders = true
			//fmt.Fprintf(os.Stderr, "header is not match: \n%v, \n%v\n", res.Header, accessLogLine.Response.Header)
		}
		unquotedExpectedResBody, err := strconv.Unquote(string(accessLogLine.Response.Body))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to strconv.Unquote: %v\n", err)
		}
		fmt.Printf("unquotedExpectedResBody: \n`%v`\nstring(resBody): \n`%v`\n", unquotedExpectedResBody, string(resBody))
		isSameResponseBody, err := jsonEqual(unquotedExpectedResBody, string(resBody), excludeBodies)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to jsonEqual: %v\n", err)
		}
		isNotSameResponseBody = !isSameResponseBody

		// --------------------------------- output ---------------------------------
		if isNotSameStatusCode || isNotSameHeaders || isNotSameResponseBody {
			reqInfo(i+1, accessLogLine)
		}
		if isNotSameStatusCode {
			printDiffInt(accessLogLine.Response.Status, res.StatusCode, "Different status code: ")
		}
		if isNotSameHeaders {
			for _, h := range unequalHeaders {
				printDiff(accessLogLine.Response.Header.Get(h), res.Header.Get(h), "%s header different: ", green(h))
			}
		}
		if isNotSameResponseBody {
			fmt.Printf("%v\n", green("Different Response body: "))
			filenameExpected, err := dumpBodyToTempFile("response_body_expected_", []byte(unquotedExpectedResBody))
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to dump response body to temp file: %v\n", err)
			}
			filenameActual, err := dumpBodyToTempFile("response_body_actual_", resBody)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to dump response body to temp file: %v\n", err)
			}
			printDiff(filenameExpected, filenameActual, "Response body files: ")
		}
		if isNotSameStatusCode || isNotSameHeaders || isNotSameResponseBody {
			// TODO os.exit(1)
			return
		}
	}

	fmt.Println("ok")
	//p.Wait()
}

func reqInfo(lineNum int, line AccessLogLine) {
	unquotedReqBody, err := strconv.Unquote(string(line.Request.Body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to strconv.Unquote: %v\n", err)
	}

	var isBodyEmpty bool
	filename, err := dumpBodyToTempFile("request_body_", []byte(unquotedReqBody))
	if errors.Is(err, ErrEmptyBody) {
		isBodyEmpty = true
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "failed to dump request body to temp file: %v\n", err)
	}
	if isBodyEmpty {
		fmt.Printf("--------- Request Info ---------\naccess log line number: %v\ntrace_id: %v\nrequest_id: %v\nrequest url: %v %v\n--------------------------------\n\n",
			lineNum,
			line.TraceID,
			line.RequestID,
			line.Request.Method,
			line.Request.RequestURI,
		)
	} else {
		fmt.Printf("--------- Request Info ---------\naccess log line number: %v\ntrace_id: %v\nrequest_id: %v\nrequest url: %v %v\nrequest body file: %v\n--------------------------------\n\n",
			lineNum,
			line.TraceID,
			line.RequestID,
			line.Request.Method,
			line.Request.RequestURI,
			filename,
		)
	}

}

var ErrEmptyBody = errors.New("empty body")

func dumpBodyToTempFile(filePrefix string, body []byte) (tmpFilename string, err error) {
	if string(body) == "" {
		return "", ErrEmptyBody
	}
	f, err := os.CreateTemp("", filePrefix)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer f.Close()

	buf := &bytes.Buffer{}
	if err := json.Indent(buf, body, "", "  "); err != nil {
		fmt.Printf("body:%v\n", string(body))
		return "", fmt.Errorf("failed to json.Indent: %w", err)
	}
	_, err = f.Write(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to write to temp file: %w", err)
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

//func vs(a, b string, f string, v ...interface{}) bool {
//	notsame = a != b
//	if notsame {
//		s := fmt.Sprintf(f, v...)
//		fmt.Printf("%s\n    %s\n    %s\n", s, on(0, a), on(1, b))
//	}
//	return notsame
//}
//
//func vsi(a, b int, f string, v ...interface{}) bool {
//	notsame = a != b
//	if notsame {
//		s := fmt.Sprintf(f, v...)
//		fmt.Printf("%s\n    %s\n    %s\n", s, oni(0, a), oni(1, b))
//	}
//	return notsame
//}

func printDiff(expected, actual string, f string, v ...interface{}) {
	s := fmt.Sprintf(f, v...)
	fmt.Printf("%s\n    expected: %s\n    actual  : %s\n", s, on(0, expected), on(1, actual))
}
func printDiffInt(expected, actual int, f string, v ...interface{}) {
	printDiff(fmt.Sprintf("%d", expected), fmt.Sprintf("%d", actual), f, v...)
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
func jsonEqual(a, b string, ignoreBodies map[string]struct{}) (bool, error) {
	var objA interface{}
	var objB interface{}

	if err := json.Unmarshal([]byte(a), &objA); err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(b), &objB); err != nil {
		return false, err
	}

	return deepEqual(objA, objB, ignoreBodies, ""), nil
}

// 順序を無視して要素を比較
func deepEqual(a, b interface{}, ignoreBodies map[string]struct{}, keyHierarchy string) bool {
	switch aVal := a.(type) {
	case []interface{}:
		bVal, ok := b.([]interface{})
		if !ok || len(aVal) != len(bVal) {
			fmt.Println("here 1", keyHierarchy, aVal, bVal)
			return false
		}

		// sortするとNlog(N)になるが、uuidなど都度変わる値があるとsortが安定でなくなるため全探索する
		visited := make([]bool, len(bVal))
		for _, vA := range aVal {
			matchFound := false
			for j, vB := range bVal {
				if !visited[j] && deepEqual(vA, vB, ignoreBodies, keyHierarchy) {
					visited[j] = true
					matchFound = true
					break
				}
			}
			if !matchFound {
				fmt.Println("here 7", keyHierarchy, aVal, bVal)
				return false
			}
		}
		return true
	case map[string]interface{}:
		bVal, ok := b.(map[string]interface{})
		if !ok || len(aVal) != len(bVal) {
			fmt.Println("here 3", keyHierarchy, aVal, bVal)
			return false
		}

		for key, valA := range aVal {
			valB, exists := bVal[key]
			if !exists {
				fmt.Println("here 4", keyHierarchy, aVal, bVal)
				return false
			}

			presentKeyHierarchy := ""
			if keyHierarchy == "" {
				presentKeyHierarchy = key
			} else {
				presentKeyHierarchy = keyHierarchy + "." + key
			}
			fmt.Printf("key: %s, presentKeyHiralchy: %s\n", key, presentKeyHierarchy)
			if _, ok := ignoreBodies[presentKeyHierarchy]; ok {
				continue
			}
			if !deepEqual(valA, valB, ignoreBodies, presentKeyHierarchy) {
				fmt.Println("here 5", keyHierarchy, aVal, bVal)
				return false
			}
		}
		return true
	default:
		if _, ok := ignoreBodies[keyHierarchy]; ok {
			return true
		}
		if a != b {
			fmt.Println("here 6", keyHierarchy, a, b)
		}
		return a == b
	}
}
