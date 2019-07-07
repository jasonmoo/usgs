package earthquake

import (
	"testing"
	"time"
)

func TestQueryParameters(t *testing.T) {

	qp := NewQueryParameters()

	qp.StartTime = time.Date(2019, 1, 2, 3, 4, 1, 0, time.UTC)
	qp.EndTime = time.Date(2019, 1, 2, 3, 4, 2, 0, time.UTC)
	qp.UpdatedAfter = time.Date(2019, 1, 2, 3, 4, 3, 0, time.UTC)
	qp.MinLatitude = 1.1
	qp.MinLongitude = 1.2
	qp.MaxLatitude = 1.3
	qp.MaxLongitude = 1.4
	qp.Latitude = 1.5
	qp.Longitude = 1.6
	qp.MaxRadius = 1.7
	qp.MaxRadiusKM = 1.8
	qp.Catalog = "catalog"
	qp.Contributor = "contributor"
	qp.EventID = "string"
	qp.IncludeAllMagnitudes = true
	qp.IncludeAllOrigins = true
	qp.IncludeDeleted = true
	qp.IncludeSuperseded = true
	qp.Limit = 1000
	qp.MaxDepth = 1.9
	qp.MaxMagnitude = 1.11
	qp.MinDepth = 1.12
	qp.MinMagnitude = 1.13
	qp.Offset = 1001
	qp.OrderBy = "order"
	qp.AlertLevel = "alertlevel"
	qp.EventType = "eventtype"
	qp.MaxCdi = 1.14
	qp.MaxGap = 1.15
	qp.MaxMmi = 1.16
	qp.MaxSig = 1002
	qp.MinCdi = 1.17
	qp.MinFelt = 1003
	qp.MinGap = 1.18
	qp.MinSig = 1004
	qp.ProductType = "producttype"
	qp.ProductCode = "string"
	qp.ReviewStatus = "reviewstatus"

	const expected = `alertlevel=alertlevel&catalog=catalog&contributor=contributor&endtime=2019-01-02T03%3A04%3A02Z&eventid=string&eventtype=eventtype&format=geojson&includeallmagnitudes=true&includeallorigins=true&includedeleted=true&includesuperseded=true&latitude=1.5&limit=1000&longitude=1.6&maxcdi=1.14&maxdepth=1.9&maxgap=1.15&maxlatitude=1.3&maxlongitude=1.4&maxmagnitude=1.11&maxmmi=1.16&maxradius=1.7&maxradiuskm=1.8&maxsig=1002&mincdi=1.17&mindepth=1.12&minfelt=1003&mingap=1.18&minlatitude=1.1&minlongitude=1.2&minmagnitude=1.13&minsig=1004&offset=1001&orderby=order&productcode=string&producttype=producttype&reviewstatus=reviewstatus&starttime=2019-01-02T03%3A04%3A01Z&updatedafter=2019-01-02T03%3A04%3A03Z`

	if out := qp.Encode(); out != expected {
		t.Errorf("expected: %q\n got: %q", expected, out)
	}

}
