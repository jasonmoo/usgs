package earthquake

import (
	"errors"
	"testing"
	"time"
)

var client = NewClient()

func TestGetApplicationInfo(t *testing.T) {

	resp, err := client.GetApplicationInfo()
	if err != nil {
		t.Error(err)
	}
	if len(resp.Catalogs) == 0 {
		t.Errorf("expected resp.Catalogs, got none")
	}
	if len(resp.Contributors) == 0 {
		t.Errorf("expected resp.Contributors, got none")
	}
	if len(resp.EventTypes) == 0 {
		t.Errorf("expected resp.EventTypes, got none")
	}
	if len(resp.MagnitudeTypes) == 0 {
		t.Errorf("expected resp.MagnitudeTypes, got none")
	}
	if len(resp.ProductTypes) == 0 {
		t.Errorf("expected resp.ProductTypes, got none")
	}

}

func TestGetApplicationWADL(t *testing.T) {
	resp, err := client.GetApplicationWADL()
	if err != nil {
		t.Error(err)
	}

	const expectedBase = "https://earthquake.usgs.gov/fdsnws/event/1/"
	if resp.Resources.Base != expectedBase {
		t.Errorf("expected %q, got %q", expectedBase, resp.Resources.Base)
	}

	if len(resp.Resources.Resource) == 0 {
		t.Errorf("expected resp.Resources.Resources, got none")
	}
}

func TestGetCatalogs(t *testing.T) {
	resp, err := client.GetCatalogs()
	if err != nil {
		t.Error(err)
	}
	if len(resp.Catalogs) == 0 {
		t.Errorf("expected resp.Catalogs, got none")
	}
}

func TestGetContributors(t *testing.T) {
	resp, err := client.GetContributors()
	if err != nil {
		t.Error(err)
	}
	if len(resp.Contributors) == 0 {
		t.Errorf("expected resp.Contributors, got none")
	}
}

func TestGetCount(t *testing.T) {

	qp := NewQueryParameters()

	qp.StartTime = time.Date(2018, 1, 2, 3, 4, 5, 0, time.UTC)
	qp.EndTime = time.Date(2018, 1, 3, 3, 4, 5, 0, time.UTC)

	resp, err := client.GetCount(qp)
	if err != nil {
		t.Error(err)
	}

	const (
		Count      = 317
		MaxAllowed = 20000
	)
	if resp.Count != Count {
		t.Errorf("expected %d, got %d", Count, resp.Count)
	}
	if resp.MaxAllowed != MaxAllowed {
		t.Errorf("expected %d, got %d", MaxAllowed, resp.MaxAllowed)
	}

}

func TestGetQuery(t *testing.T) {

	qp := NewQueryParameters()

	qp.StartTime = time.Date(2018, 1, 2, 3, 4, 5, 0, time.UTC)
	qp.EndTime = time.Date(2018, 1, 2, 4, 4, 5, 0, time.UTC)
	qp.Limit = 1

	resp, err := client.GetQuery(qp)
	if err != nil {
		t.Error(err)
	}

	const expected = "uw61362166"

	if len(resp.Features) == 0 {
		t.Errorf("expected resp.Features, got none")
	}
	if resp.Features[0].ID != expected {
		t.Errorf("expected %q, got %q", expected, resp.Features[0].ID)
	}

}

func TestGetQueryPaged(t *testing.T) {

	qp := NewQueryParameters()

	qp.StartTime = time.Date(2018, 1, 2, 3, 4, 5, 0, time.UTC)
	qp.EndTime = time.Date(2018, 1, 2, 4, 4, 5, 0, time.UTC)
	qp.Limit = 1
	qp.TotalResults = 3

	var (
		expected = [...]string{
			"uw61362166",
			"us2000ck9q",
			"ci37845255",
		}

		i int
	)
	err := client.GetQueryPaged(qp, func(resp *GetQueryResponse) error {

		if i > len(expected)-1 {
			return errors.New("exceeded expected total results")
		}

		if len(resp.Features) == 0 {
			t.Errorf("expected resp.Features, got none")
		}
		if resp.Features[0].ID != expected[i] {
			t.Errorf("expected %q, got %q", expected[i], resp.Features[0].ID)
		}
		i++

		return nil

	})
	if err != nil {
		t.Error(err)
	}

}

func TestGetVersion(t *testing.T) {
	resp, err := client.GetVersion()
	if err != nil {
		t.Error(err)
	}
	if resp.Version == "" {
		t.Errorf("expected version, got none")
	}
}
