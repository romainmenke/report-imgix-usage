package sources

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/romainmenke/report-imgix-usage/counters"
)

func Get(auth string, page int) (*Sources, error) {
	req, err := http.NewRequest("GET", "https://api.imgix.com/v4/sources?filter%5Benabled%5D=true&page%5Bnumber%5D="+fmt.Sprint(page)+"&page%5Bsize%5D=40", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	sources := &Sources{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(sources)
	if err != nil {
		return nil, err
	}

	if sources.Meta.Pagination.HasNextPage {
		extraSources, err := Get(auth, page+1)
		if err != nil {
			return nil, err
		}

		sources.Data = append(sources.Data, extraSources.Data...)
	}

	return sources, nil
}

type Sources struct {
	Data     []*Data       `json:"data"`
	Included []interface{} `json:"included"`
	Jsonapi  Jsonapi       `json:"jsonapi"`
	Meta     Meta          `json:"meta"`
}

type Production struct {
	AwsAccessKey           interface{}   `json:"aws_access_key"`
	AwsBucket              interface{}   `json:"aws_bucket"`
	AwsPrefix              interface{}   `json:"aws_prefix"`
	AwsSecretKeyEncrypted  interface{}   `json:"aws_secret_key_encrypted"`
	AzureHost              interface{}   `json:"azure_host"`
	AzurePrefix            interface{}   `json:"azure_prefix"`
	AzureSasString         interface{}   `json:"azure_sas_string"`
	CacheTTLBehavior       string        `json:"cache_ttl_behavior"`
	CacheTTLError          int           `json:"cache_ttl_error"`
	CacheTTLValue          int           `json:"cache_ttl_value"`
	CrossdomainCorsEnabled bool          `json:"crossdomain_cors_enabled"`
	CrossdomainXMLEnabled  bool          `json:"crossdomain_xml_enabled"`
	CustomDomains          []interface{} `json:"custom_domains"`
	DefaultParams          interface{}   `json:"default_params"`
	GcsAccessKey           interface{}   `json:"gcs_access_key"`
	GcsBucket              interface{}   `json:"gcs_bucket"`
	GcsPrefix              interface{}   `json:"gcs_prefix"`
	GcsSecretKeyEncrypted  interface{}   `json:"gcs_secret_key_encrypted"`
	ImageError             interface{}   `json:"image_error"`
	ImageErrorAppendQs     bool          `json:"image_error_append_qs"`
	ImageMissing           interface{}   `json:"image_missing"`
	ImageMissingAppendQs   bool          `json:"image_missing_append_qs"`
	ImgixSubdomains        []string      `json:"imgix_subdomains"`
	SecureURLEnabled       bool          `json:"secure_url_enabled"`
	Type                   string        `json:"type"`
	WebfolderPrefix        string        `json:"webfolder_prefix"`
}

type Attributes struct {
	DateCreated      int        `json:"date_created"`
	DateDeployed     int        `json:"date_deployed"`
	DateModified     int        `json:"date_modified"`
	DeploymentStatus string     `json:"deployment_status"`
	Enabled          bool       `json:"enabled"`
	GlobalVersion    int        `json:"global_version"`
	Name             string     `json:"name"`
	Production       Production `json:"production"`
	SecureURLToken   string     `json:"secure_url_token"`
}

type DataContainer struct {
	Data Data `json:"data"`
}

type Relationships struct {
	Account              DataContainer `json:"account"`
	DeployingUser        DataContainer `json:"deploying_user"`
	StagingConfiguration DataContainer `json:"staging_configuration"`
}

type Data struct {
	sync.Mutex

	Attributes    Attributes                    `json:"attributes"`
	ID            string                        `json:"id"`
	Relationships *Relationships                `json:"relationships,omitempty"`
	Type          string                        `json:"type"`
	Counters      map[string]*counters.Counters `json:"-"`
}

func (x *Data) GetCounters(auth string, from time.Time, to time.Time) error {
	c, err := counters.Get(auth, x.ID, from, to)
	if err != nil {
		return err
	}

	x.Lock()
	defer x.Unlock()

	if x.Counters == nil {
		x.Counters = make(map[string]*counters.Counters)
	}

	x.Counters[from.Format("2006-01")] = c

	return nil
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

type Pagination struct {
	CurrentPage     int  `json:"currentPage"`
	HasNextPage     bool `json:"hasNextPage"`
	HasPreviousPage bool `json:"hasPreviousPage"`
	NextPage        int  `json:"nextPage"`
	PageSize        int  `json:"pageSize"`
	PreviousPage    int  `json:"previousPage"`
	TotalPages      int  `json:"totalPages"`
	TotalRecords    int  `json:"totalRecords"`
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
	Pagination     Pagination     `json:"pagination"`
	Server         Server         `json:"server"`
}
