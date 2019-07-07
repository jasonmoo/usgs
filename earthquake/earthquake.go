package earthquake

//go:generate go run generate_constants.go

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

type (
	// https://earthquake.usgs.gov/fdsnws/event/1/
	Client struct {
		c *http.Client
	}

	GetApplicationInfoResponse struct {
		Catalogs       []Catalog       `json:"catalogs"`
		Contributors   []Contributor   `json:"contributors"`
		EventTypes     []EventType     `json:"eventtypes"`
		MagnitudeTypes []MagnitudeType `json:"magnitudetypes"`
		ProductTypes   []ProductType   `json:"producttypes"`
	}

	GetApplicationWADLResponse struct {
		XMLName   xml.Name `xml:"application"`
		Xmlns     string   `xml:"xmlns,attr"`
		Q         string   `xml:"q,attr"`
		Xs        string   `xml:"xs,attr"`
		Resources struct {
			Base     string `xml:"base,attr"`
			Resource []struct {
				Path   string `xml:"path,attr"`
				Method struct {
					ID       string `xml:"id,attr"`
					Name     string `xml:"name,attr"`
					Response struct {
						Status         string `xml:"status,attr"`
						Representation []struct {
							MediaType string `xml:"mediaType,attr"`
							Element   string `xml:"element,attr"`
						} `xml:"representation"`
					} `xml:"response"`
					Request struct {
						Param []struct {
							Name      string `xml:"name,attr"`
							Style     string `xml:"style,attr"`
							Type      string `xml:"type,attr"`
							Default   string `xml:"default,attr"`
							MediaType string `xml:"mediaType,attr"`
							Option    []struct {
								Value     string `xml:"value,attr"`
								MediaType string `xml:"mediaType,attr"`
							} `xml:"option"`
						} `xml:"param"`
					} `xml:"request"`
				} `xml:"method"`
			} `xml:"resource"`
		} `xml:"resources"`
	}

	GetCatalogsResponse struct {
		XMLName  xml.Name  `xml:"Catalogs"`
		Catalogs []Catalog `xml:"Catalog"`
	}

	GetContributorsResponse struct {
		XMLName      xml.Name      `xml:"Contributors"`
		Contributors []Contributor `xml:"Contributor"`
	}

	GetCountResponse struct {
		Count      int `json:"count"`
		MaxAllowed int `json:"maxAllowed"`
	}

	GetQueryResponse struct {
		Bbox     []float64 `json:"bbox"`
		Features []struct {
			Geometry struct {
				Coordinates []float64 `json:"coordinates"`
				Type        string    `json:"type"`
			} `json:"geometry"`
			ID         string `json:"id"`
			Properties struct {
				Alert   interface{} `json:"alert"`
				Cdi     interface{} `json:"cdi"`
				Code    string      `json:"code"`
				Detail  string      `json:"detail"`
				Dmin    float64     `json:"dmin"`
				Felt    interface{} `json:"felt"`
				Gap     int         `json:"gap"`
				Ids     string      `json:"ids"`
				Mag     float64     `json:"mag"`
				MagType string      `json:"magType"`
				Mmi     interface{} `json:"mmi"`
				Net     string      `json:"net"`
				Nst     int         `json:"nst"`
				Place   string      `json:"place"`
				Rms     float64     `json:"rms"`
				Sig     int         `json:"sig"`
				Sources string      `json:"sources"`
				Status  string      `json:"status"`
				Time    UnixEpoch   `json:"time"`
				Title   string      `json:"title"`
				Tsunami int         `json:"tsunami"`
				Type    string      `json:"type"`
				Types   string      `json:"types"`
				Tz      int         `json:"tz"`
				Updated UnixEpoch   `json:"updated"`
				URL     string      `json:"url"`
			} `json:"properties"`
			Type string `json:"type"`
		} `json:"features"`
		Metadata struct {
			API       string    `json:"api"`
			Count     int       `json:"count"`
			Generated UnixEpoch `json:"generated"`
			Status    int       `json:"status"`
			Title     string    `json:"title"`
			URL       string    `json:"url"`
		} `json:"metadata"`
		Type string `json:"type"`
	}

	GetVersionResponse struct {
		Version string
	}

	UnixEpoch struct{ time.Time }
)

func (e *UnixEpoch) UnmarshalJSON(b []byte) error {
	v, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	(*e).Time = time.Unix(0, v*int64(time.Millisecond))
	return nil
}

type transportFunc func(req *http.Request) (*http.Response, error)

func (t transportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return t(req)
}

func NewClient() *Client {
	return &Client{
		c: &http.Client{
			Transport: transportFunc(func(req *http.Request) (*http.Response, error) {

				// https://earthquake.usgs.gov/fdsnws/event/1/[METHOD[?PARAMETERS]]
				req.URL.Scheme = "https"
				req.URL.Host = "earthquake.usgs.gov"
				req.URL.Path = path.Join("/fdsnws/event/1", req.URL.Path)

				req.Header.Set("User-Agent", "github.com/jasonmoo/usgs/earthquake v1.0")

				resp, err := http.DefaultTransport.RoundTrip(req)
				if err != nil {
					return nil, err
				}
				if resp.StatusCode != 200 {
					var body []byte
					if !strings.Contains(resp.Header.Get("Content-Type"), "html") {
						body, _ = ioutil.ReadAll(resp.Body)
					}
					resp.Body.Close()
					return nil, fmt.Errorf("%s (%d): %q", resp.Status, resp.StatusCode, string(body))
				}

				return resp, nil

			}),
		},
	}
}

// request known enumerated parameter values for the interface.
func (c *Client) GetApplicationInfo() (*GetApplicationInfoResponse, error) {
	// https://earthquake.usgs.gov/fdsnws/event/1/application.json
	resp, err := c.c.Get("/application.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var v GetApplicationInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

// request WADL for the interface.
func (c *Client) GetApplicationWADL() (*GetApplicationWADLResponse, error) {
	// https://earthquake.usgs.gov/fdsnws/event/1/application.wadl
	resp, err := c.c.Get("/application.wadl")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var v GetApplicationWADLResponse
	if err := xml.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

// request available catalogs.
func (c *Client) GetCatalogs() (*GetCatalogsResponse, error) {
	// https://earthquake.usgs.gov/fdsnws/event/1/catalogs
	resp, err := c.c.Get("/catalogs")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var v GetCatalogsResponse
	if err := xml.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

// request available contributors
func (c *Client) GetContributors() (*GetContributorsResponse, error) {
	// https://earthquake.usgs.gov/fdsnws/event/1/contributors
	resp, err := c.c.Get("/contributors")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var v GetContributorsResponse
	if err := xml.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

// to perform a count on a data request. Count uses the same parameters as the query method, and is availablein these formats: plain text (default), geojson, and xml.
func (c *Client) GetCount(qp *queryParameters) (*GetCountResponse, error) {
	// https://earthquake.usgs.gov/fdsnws/event/1/count?format=geojson
	// https://earthquake.usgs.gov/fdsnws/event/1/count?starttime=2014-01-01&endtime=2014-01-02
	resp, err := c.c.Get("/count?" + qp.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var v GetCountResponse
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

// to submit a data request. See the parameters section for supported url parameters.
func (c *Client) GetQuery(qp *queryParameters) (*GetQueryResponse, error) {
	// https://earthquake.usgs.gov/fdsnws/event/1/query?format=geojson&starttime=2014-01-01&endtime=2014-01-02
	// https://earthquake.usgs.gov/fdsnws/event/1/query?format=xml&starttime=2014-01-01&endtime=2014-01-02&minmagnitude=5
	resp, err := c.c.Get("/query?" + qp.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var v GetQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return &v, nil
}

// run a query and retrieve the full dataset over multiple requests
func (c *Client) GetQueryPaged(qp *queryParameters, f func(*GetQueryResponse) error) error {

	limit := qp.Limit
	qp.Limit = 0

	cresp, err := c.GetCount(qp)
	if err != nil {
		return err
	}

	if limit == 0 {
		qp.Limit = cresp.MaxAllowed
	} else {
		qp.Limit = limit
	}

	if cresp.Count < qp.TotalResults {
		qp.TotalResults = cresp.Count
	}

	for qp.Offset = 1; qp.Offset <= qp.TotalResults; qp.Offset += qp.Limit {
		if err := func() error {
			resp, err := c.c.Get("/query?" + qp.Encode())
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			var v GetQueryResponse
			if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
				return err
			}
			if err := f(&v); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			return err
		}
	}

	return nil

}

// request full service version number
func (c *Client) GetVersion() (*GetVersionResponse, error) {
	// https://earthquake.usgs.gov/fdsnws/event/1/version
	resp, err := c.c.Get("/version")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &GetVersionResponse{Version: strings.TrimSpace(string(data))}, nil
}
