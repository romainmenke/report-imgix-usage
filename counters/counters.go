package counters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func Get(client *http.Client, id string, from time.Time, to time.Time) (*Counters, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.imgix.com/v4/sources/"+id+"/counters?filter%5Bend%5D="+fmt.Sprint(to.Unix())+"&filter%5Bstart%5D="+fmt.Sprint(from.Unix()), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
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

func (x MultipleCounters) Sum() *Counters {
	if len(x) == 0 {
		return &Counters{}
	}

	sum := &Counters{
		Included: x[0].Included,
		Jsonapi:  x[0].Jsonapi,
		Meta:     x[0].Meta,
	}

	for _, e := range x {
		sum.Data = append(sum.Data, e.Data...)
	}

	sum.SetSum()

	return sum
}

type Counters struct {
	Data     []Data        `json:"data"`
	Included []interface{} `json:"included"`
	Jsonapi  Jsonapi       `json:"jsonapi"`
	Meta     Meta          `json:"meta"`
	Sum      *Attributes   `json:"-"`
}

func (x *Counters) SetSum() {
	if x.Sum != nil {
		return
	}

	sum := Attributes{}

	for _, data := range x.Data {
		sum.Bandwidth += data.Attributes.Bandwidth
		sum.CdnBandwidth += data.Attributes.CdnBandwidth
		sum.CdnErrors += data.Attributes.CdnErrors
		sum.CdnErrorsMild += data.Attributes.CdnErrorsMild
		sum.CdnErrorsSevere += data.Attributes.CdnErrorsSevere
		sum.CdnMisses += data.Attributes.CdnMisses
		sum.CdnRequests += data.Attributes.CdnRequests
		sum.ConvertBandwidth += data.Attributes.ConvertBandwidth
		sum.ConvertImages += data.Attributes.ConvertImages
		sum.ConvertTime += data.Attributes.ConvertTime
		sum.Errors += data.Attributes.Errors
		sum.ErrorsMild += data.Attributes.ErrorsMild
		sum.ErrorsSevere += data.Attributes.ErrorsSevere
		sum.FetchBandwidth += data.Attributes.FetchBandwidth
		sum.FetchRequests += data.Attributes.FetchRequests
		sum.FetchTime += data.Attributes.FetchTime
		sum.Images += data.Attributes.Images
		sum.QueueTime += data.Attributes.QueueTime
		sum.RenderBandwidth += data.Attributes.RenderBandwidth
		sum.RenderErrors += data.Attributes.RenderErrors
		sum.RenderErrorsMild += data.Attributes.RenderErrorsMild
		sum.RenderErrorsSevere += data.Attributes.RenderErrorsSevere
		sum.RenderRequests += data.Attributes.RenderRequests
		sum.RenderTime += data.Attributes.RenderTime
		sum.Renders += data.Attributes.Renders
		sum.Requests += data.Attributes.Requests
		sum.ResponseTime += data.Attributes.ResponseTime
		sum.Successes += data.Attributes.Successes
	}

	x.Sum = &sum
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
