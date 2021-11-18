package fe-measure

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/com-gft-tsbo-source/go-common/ms-framework/microservice"
)

// ###########################################################################
// ###########################################################################
// FeMeasure
// ###########################################################################
// ###########################################################################

type header struct {
	Key   string
	Value string
}

type deviceReply struct {
	Data   *interface{}
	Header *http.Header
}

type templateData struct {
	Frontend           *FeMeasure
	Count              int
	BaseURL            string
	Reply              []deviceReply
	Referer            string
	RequestEnvironment string
	RequestSource      string
	RequestTarget      string
	RequestURI         string
	RequestURL         string
	RequestReferer     string
	RequestHeader      *http.Header
}

// FeMeasure Encapsulates the fe-measure data
type FeMeasure struct {
	microservice.MicroService
	*FrontendConfiguration
	DeviceUrlConfiguration

	starttime       time.Time
	lastRequest     time.Time
	template        *template.Template
	templateModTime time.Time
}

// ###########################################################################

// InitFromArgs ...
func InitFromArgs(ms *FeMeasure, args []string, flagset *flag.FlagSet) *FeMeasure {
	var cfg Configuration
	var err error

	if flagset == nil {
		flagset = flag.NewFlagSet("fe-measure", flag.PanicOnError)
	}

	InitConfigurationFromArgs(&cfg, args, flagset)
	ms.FrontendConfiguration = &cfg.FrontendConfiguration
	ms.DeviceURLs = cfg.DeviceURLs
	contentHandler := ms.DefaultHandler()
	contentHandler.Get = ms.httpGetIndex
	microservice.Init(&ms.MicroService, &cfg.Configuration, contentHandler)

	err = ms.loadTemplate()
	if err != nil {
		panic(fmt.Sprintf("Error: Templatefile '%s' not loaded. Error was %s!\n", ms.GetTemplateFile(), err))
	}

	ms.starttime = time.Now()
	ms.lastRequest = ms.starttime
	staticHandler := ms.DefaultHandler()
	staticHandler.Get = ms.httpGetStatic
	ms.AddHandler(ms.GetTemplateURL(), contentHandler)
	ms.AddHandler(ms.GetStaticURL(), staticHandler)
	return ms
}

// ---------------------------------------------------------------------------

func (ms *FeMeasure) httpGetIndex(w http.ResponseWriter, r *http.Request) (status int, contentLen int, msg string) {
	ms.lastRequest = time.Now()
	if r.URL.Path != "/" && r.URL.Path != ms.GetTemplateURL() {
		return ms.PageNotFound(w, r)
	}

	var environment string

	// ------------------------------------------------------------------------
	// Parse parameters

	r.ParseForm()

	if len(r.Form) > 0 {
		var newUrls []string

		// for url, _ := range r.Form {
		for key, values := range r.Form {
			if key == "environment" {
				for _, value := range values {
					environment = value
				}
				continue
			}

			if key == "device" {
				for _, value := range values {
					if (!strings.HasPrefix(value, "http:")) && (!strings.HasPrefix(value, "https:")) {
						value = "http://" + value + "/measure"
					}
					newUrls = append(newUrls, value)
				}
				ms.DeviceURLs = newUrls
				continue
			}
		}
	}

	if len(environment) == 0 {
		environment = r.Header.Get("X-Environment")
	}

	// ------------------------------------------------------------------------
	// Determine base URL

	var baseUrl string

	baseUrl = r.Header.Get("X-Forwarded-Host")
	if len(baseUrl) == 0 {
		baseUrl = r.Host
	}
	if len(baseUrl) == 0 {
		baseUrl = ms.GetHost()
		if ms.GetPort() > 0 {
			baseUrl = fmt.Sprintf("%s:%d", baseUrl, ms.GetPort())
		}
	}

	if len(baseUrl) >= 0 {
		if ms.HasTLS() {
			baseUrl = "https://" + baseUrl
		} else {
			baseUrl = "http://" + baseUrl
		}
		baseUrl = baseUrl + "/"
		// if len(environment) > 0 {
		// 	baseUrl = fmt.Sprintf("%s%s/", baseUrl, environment)
		// }
	}

	// Query the devices
	var data = &templateData{
		Frontend:           ms,
		Count:              0,
		BaseURL:            baseUrl,
		RequestEnvironment: environment,
		RequestReferer:     r.Referer(),
		RequestURL:         r.URL.String(),
		RequestURI:         r.RequestURI,
		RequestSource:      r.RemoteAddr,
		RequestTarget:      r.Host,
	}

	// ------------------------------------------------------------------------
	// Check for devices in the URL and use them from that point on

	for _, url := range ms.DeviceURLs {
		value, replyHeader, err := ms.loadUrl(url, r, environment)
		if err != nil {
			ms.GetLogger().Println(err.Error())
			continue
		}
		data.Reply = append(data.Reply, deviceReply{Data: value, Header: replyHeader})
	}

	data.Count = len(data.Reply)

	data.RequestHeader = &r.Header
	// for key, values := range r.Header {
	// 	for _, value := range values {
	// 		data.RequestHeader = append(data.RequestHeader, header{Key: key, Value: value})
	// 	}
	// }
	// for key, values := range r.Trailer {
	// 	for _, value := range values {
	// 		data.RequestHeader = append(data.RequestHeader, header{Key: key, Value: value})
	// 	}
	// }

	// Parse the template and prepare the reply
	status = http.StatusOK
	err := ms.loadTemplate()
	if err != nil {
		panic(fmt.Sprintf("Error: Templatefile '%s' not loaded. Error was %s!\n", ms.GetTemplateFile(), err))
	}
	msg = fmt.Sprintf("'v%s' in '%s' serving '@T/%s'.", ms.GetVersion(), environment, ms.GetTemplateFile())
	var buffer bytes.Buffer
	err = ms.template.Execute(&buffer, data)
	if err != nil {
		panic(fmt.Sprintf("Error: Templatefile '%s' not loaded. Error was %s!\n", ms.GetTemplateFile(), err))
	}

	// Reply ...
	ms.SetResponseHeaders("text/html; charset=utf-8", w, r)
	w.WriteHeader(status)
	w.Write(buffer.Bytes())

	return status, buffer.Len(), msg
}

// ---------------------------------------------------------------------------

func (ms *FeMeasure) httpGetStatic(w http.ResponseWriter, r *http.Request) (status int, contentLen int, msg string) {
	ms.lastRequest = time.Now()
	if len(ms.GetStaticDir()) == 0 {
		return ms.PageNotFound(w, r)
	}

	subPath := r.URL.Path[7:]
	realPath := path.Join(ms.GetStaticDir(), subPath)

	f, err := os.Open(realPath)
	if err != nil {
		return ms.PageNotFound(w, r)
	}

	stat, err := os.Stat(realPath)
	if err != nil {
		return ms.PageNotFound(w, r)
	}

	environment := r.Header.Get("X-Environment")
	status = http.StatusOK
	msg = fmt.Sprintf("'v%s' in '%s' serving '@S%s'.", ms.GetVersion(), environment, subPath)
	http.ServeContent(w, r, "."+r.URL.Path, stat.ModTime(), f)
	return status, int(stat.Size()), msg
}

// ---------------------------------------------------------------------------

// Run ...
func (ms *FeMeasure) Run() {

	if len(ms.GetStaticDir()) != 0 {
		if len(ms.GetStaticURL()) != 0 {
			ms.GetLogger().Println(fmt.Sprintf("Serving static files in '%s' at '%s'.", ms.GetStaticDir(), ms.GetStaticURL()))
		}
	}

	ms.GetLogger().Println(fmt.Sprintf("Serving template file '%s' at '%s'.", ms.GetTemplateFile(), ms.GetTemplateURL()))
	ms.MicroService.Run()
}

// ---------------------------------------------------------------------------

// loadTemplate ...
func (ms *FeMeasure) loadTemplate() error {
	var err error
	var info os.FileInfo

	info, err = os.Stat(ms.GetTemplateFile())

	if err != nil {
		return err
	}

	if ms.template == nil {
		ms.template, err = template.ParseFiles(ms.GetTemplateFile())
		ms.GetLogger().Println(fmt.Sprintf("Loading template '%s'.", ms.GetTemplateFile()))
	} else if ms.templateModTime.Before(info.ModTime()) {
		ms.GetLogger().Println(fmt.Sprintf("Refreshing template '%s'.", ms.GetTemplateFile()))
		ms.template, err = template.ParseFiles(ms.GetTemplateFile())
	}

	if err == nil {
		ms.templateModTime = info.ModTime()
	}

	return err
}

// ---------------------------------------------------------------------------

// Update reads the new value
func (ms *FeMeasure) loadUrl(url string, in *http.Request, environment string) (*interface{}, *http.Header, error) {

	req, err := http.NewRequest(http.MethodGet, url, strings.NewReader(url))

	if err != nil {
		return nil, nil, fmt.Errorf("Could not query device, message was '%s'!", err.Error())
	}

	ms.SetRequestHeaders("", req, in)

	if len(environment) > 0 {
		req.Header.Set("X-Environment", environment)
	}

	rep, err := ms.HTTPClient.Do(req)

	if err != nil {
		return nil, nil, err
	}

	defer rep.Body.Close()
	body, err := ioutil.ReadAll(rep.Body)
	ms.HTTPClient.CloseIdleConnections()

	if err != nil {
		return nil, nil, err
	}

	if rep.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("Could not update device, status was '%d', message was '%s'!", rep.StatusCode, body)
	}

	var value interface{}
	err = json.Unmarshal(body, &value)

	if err != nil {
		return nil, nil, err
	}

	return &value, &rep.Header, nil
}
