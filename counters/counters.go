package counters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func Get(auth string, id string, from time.Time, to time.Time) (*Counters, error) {
	req, err := http.NewRequest("GET", "https://api.imgix.com/v4/sources/"+id+"/counters?filter%5Bend%5D="+fmt.Sprint(to.Unix())+"&filter%5Bstart%5D="+fmt.Sprint(from.Unix()), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	counters := &Counters{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(counters)
	if err != nil {
		return nil, err
	}

	return counters, nil
}

type MultipleCounters []*Counters

func (x MultipleCounters) Cumulative() *Counters {
	if len(x) == 0 {
		return &Counters{}
	}

	cumulative := &Counters{
		Included: x[0].Included,
		Jsonapi:  x[0].Jsonapi,
		Meta:     x[0].Meta,
	}

	for _, e := range x {
		cumulative.Data = append(cumulative.Data, e.Data...)
	}

	cumulative.SetCumulative()

	return cumulative
}

type Counters struct {
	Data       []Data        `json:"data"`
	Included   []interface{} `json:"included"`
	Jsonapi    Jsonapi       `json:"jsonapi"`
	Meta       Meta          `json:"meta"`
	Cumulative *Attributes   `json:"-"`
}

func (x *Counters) SetCumulative() {
	if x.Cumulative != nil {
		return
	}

	cumulative := Attributes{}

	for _, data := range x.Data {
		cumulative.Bandwidth += data.Attributes.Bandwidth
		cumulative.CdnBandwidth += data.Attributes.CdnBandwidth
		cumulative.CdnErrors += data.Attributes.CdnErrors
		cumulative.CdnErrorsMild += data.Attributes.CdnErrorsMild
		cumulative.CdnErrorsSevere += data.Attributes.CdnErrorsSevere
		cumulative.CdnMisses += data.Attributes.CdnMisses
		cumulative.CdnRequests += data.Attributes.CdnRequests
		cumulative.ConvertBandwidth += data.Attributes.ConvertBandwidth
		cumulative.ConvertImages += data.Attributes.ConvertImages
		cumulative.ConvertTime += data.Attributes.ConvertTime
		cumulative.Errors += data.Attributes.Errors
		cumulative.ErrorsMild += data.Attributes.ErrorsMild
		cumulative.ErrorsSevere += data.Attributes.ErrorsSevere
		cumulative.FetchBandwidth += data.Attributes.FetchBandwidth
		cumulative.FetchRequests += data.Attributes.FetchRequests
		cumulative.FetchTime += data.Attributes.FetchTime
		cumulative.Images += data.Attributes.Images
		cumulative.QueueTime += data.Attributes.QueueTime
		cumulative.RenderBandwidth += data.Attributes.RenderBandwidth
		cumulative.RenderErrors += data.Attributes.RenderErrors
		cumulative.RenderErrorsMild += data.Attributes.RenderErrorsMild
		cumulative.RenderErrorsSevere += data.Attributes.RenderErrorsSevere
		cumulative.RenderRequests += data.Attributes.RenderRequests
		cumulative.RenderTime += data.Attributes.RenderTime
		cumulative.Renders += data.Attributes.Renders
		cumulative.Requests += data.Attributes.Requests
		cumulative.ResponseTime += data.Attributes.ResponseTime
		cumulative.Successes += data.Attributes.Successes
	}

	x.Cumulative = &cumulative
}

func (x *Counters) CsvRow() []string {
	if x.Cumulative == nil {
		x.SetCumulative()
	}

	row := []string{
		fmt.Sprintf("%d", x.Cumulative.Bandwidth/(1024*1024)),
		fmt.Sprintf("%d", x.Cumulative.CdnBandwidth/(1024*1024)),
		fmt.Sprintf("%d", x.Cumulative.CdnErrors),
		fmt.Sprintf("%d", x.Cumulative.CdnErrorsMild),
		fmt.Sprintf("%d", x.Cumulative.CdnErrorsSevere),
		fmt.Sprintf("%d", x.Cumulative.CdnMisses),
		fmt.Sprintf("%d", x.Cumulative.CdnRequests),
		fmt.Sprintf("%d", x.Cumulative.ConvertBandwidth/(1024*1024)),
		fmt.Sprintf("%d", x.Cumulative.ConvertImages),
		fmt.Sprintf("%d", x.Cumulative.ConvertTime),
		fmt.Sprintf("%d", x.Cumulative.Errors),
		fmt.Sprintf("%d", x.Cumulative.ErrorsMild),
		fmt.Sprintf("%d", x.Cumulative.ErrorsSevere),
		fmt.Sprintf("%d", x.Cumulative.FetchBandwidth/(1024*1024)),
		fmt.Sprintf("%d", x.Cumulative.FetchRequests),
		fmt.Sprintf("%d", x.Cumulative.FetchTime),
		fmt.Sprintf("%d", x.Cumulative.Images),
		fmt.Sprintf("%d", x.Cumulative.QueueTime),
		fmt.Sprintf("%d", x.Cumulative.RenderBandwidth/(1024*1024)),
		fmt.Sprintf("%d", x.Cumulative.RenderErrors),
		fmt.Sprintf("%d", x.Cumulative.RenderErrorsMild),
		fmt.Sprintf("%d", x.Cumulative.RenderErrorsSevere),
		fmt.Sprintf("%d", x.Cumulative.RenderRequests),
		fmt.Sprintf("%d", x.Cumulative.RenderTime),
		fmt.Sprintf("%d", x.Cumulative.Renders),
		fmt.Sprintf("%d", x.Cumulative.Requests),
		fmt.Sprintf("%d", x.Cumulative.ResponseTime),
		fmt.Sprintf("%d", x.Cumulative.Successes),
	}

	return row
}

func CsvHeaders() []string {
	headers := []string{
		"Bandwidth",
		"CdnBandwidth",
		"CdnErrors",
		"CdnErrorsMild",
		"CdnErrorsSevere",
		"CdnMisses",
		"CdnRequests",
		"ConvertBandwidth",
		"ConvertImages",
		"ConvertTime",
		"Errors",
		"ErrorsMild",
		"ErrorsSevere",
		"FetchBandwidth",
		"FetchRequests",
		"FetchTime",
		"Images",
		"QueueTime",
		"RenderBandwidth",
		"RenderErrors",
		"RenderErrorsMild",
		"RenderErrorsSevere",
		"RenderRequests",
		"RenderTime",
		"Renders",
		"Requests",
		"ResponseTime",
		"Successes",
	}

	return headers
}

type Attributes struct {
	Bandwidth          int   `json:"bandwidth"`
	CdnBandwidth       int   `json:"cdn_bandwidth"`
	CdnErrors          int   `json:"cdn_errors"`
	CdnErrorsMild      int   `json:"cdn_errors_mild"`
	CdnErrorsSevere    int   `json:"cdn_errors_severe"`
	CdnMisses          int   `json:"cdn_misses"`
	CdnRequests        int   `json:"cdn_requests"`
	ConvertBandwidth   int64 `json:"convert_bandwidth"`
	ConvertImages      int   `json:"convert_images"`
	ConvertTime        int   `json:"convert_time"`
	Errors             int   `json:"errors"`
	ErrorsMild         int   `json:"errors_mild"`
	ErrorsSevere       int   `json:"errors_severe"`
	FetchBandwidth     int   `json:"fetch_bandwidth"`
	FetchRequests      int   `json:"fetch_requests"`
	FetchTime          int   `json:"fetch_time"`
	Images             int   `json:"images"`
	QueueTime          int   `json:"queue_time"`
	RenderBandwidth    int   `json:"render_bandwidth"`
	RenderErrors       int   `json:"render_errors"`
	RenderErrorsMild   int   `json:"render_errors_mild"`
	RenderErrorsSevere int   `json:"render_errors_severe"`
	RenderRequests     int   `json:"render_requests"`
	RenderTime         int   `json:"render_time"`
	Renders            int   `json:"renders"`
	Requests           int   `json:"requests"`
	ResponseTime       int   `json:"response_time"`
	Successes          int   `json:"successes"`
	Timestamp          int   `json:"timestamp"`
}

type Account struct {
	Data Data `json:"data"`
}

type Source struct {
	Data Data `json:"data"`
}

type Relationships struct {
	Account Account `json:"account"`
	Source  Source  `json:"source"`
}

type Data struct {
	Attributes    Attributes     `json:"attributes"`
	ID            string         `json:"id"`
	Relationships *Relationships `json:"relationships"`
	Type          string         `json:"type"`
}

type Jsonapi struct {
	Version string `json:"version"`
}

type Authentication struct {
	Authorized bool        `json:"authorized"`
	ClientID   interface{} `json:"clientId"`
	Mode       string      `json:"mode"`
	ModeTitle  string      `json:"modeTitle"`
	Tag        string      `json:"tag"`
	User       string      `json:"user"`
}

type Status struct {
	Healthy   bool `json:"healthy"`
	ReadOnly  bool `json:"read_only"`
	Tombstone bool `json:"tombstone"`
}

type Server struct {
	Commit  string `json:"commit"`
	Status  Status `json:"status"`
	Version string `json:"version"`
}

type Meta struct {
	Authentication Authentication `json:"authentication"`
	Server         Server         `json:"server"`
}
