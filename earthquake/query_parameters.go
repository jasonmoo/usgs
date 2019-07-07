package earthquake

import (
	"math"
	"net/url"
	"strconv"
	"time"
)

type (
	// Query method Parameters
	// These parameters should be submitted as key=value pairs using the HTTP GET method and may not be specified more than once; if a parameter is submitted multiple times the result is undefined.
	queryParameters struct {
		// Formats
		// If no format is specified quakeml will be returned by default.
		// parameter	type	default	description
		// format	String	quakeml	Specify the output format.
		// format=csv      Response format is CSV. Mime-type is “text/csv”.
		// format=geojson  Response format is GeoJSON. Mime-type is “application/json”.
		// format=kml      Response format is KML. Mime-type is “vnd.google-earth.kml+xml”.
		// format=quakeml  Alias for "xml" format.
		// format=text     Response format is plain text. Mime-type is “text/plain”.
		// format=xml      The xml format is dependent upon the request method used.
		//
		// Note(jasonmoo): currently only geoJSON is supported by this library
		format string

		// Time
		// All times use ISO8601 Date/Time format. Unless a timezone is specified, UTC is assumed.
		// Examples:
		// 2019-07-06, Implicit UTC timezone, and time at start of the day (00:00:00)
		// 2019-07-06T18:44:38, Implicit UTC timezone.
		// 2019-07-06T18:44:38+00:00, Explicit timezone.
		// parameter	type	default	description
		// endtime      String  present time  Limit to events on or before the specified end time.
		// starttime    String  NOW - 30 days  Limit to events on or after the specified start time.
		// updatedafter String  null  Limit to events updated after the specified time.
		StartTime    time.Time
		EndTime      time.Time
		UpdatedAfter time.Time

		// Location
		// Requests that use both rectangle and circle will return the intersection,
		// which may be empty, use with caution.

		// Rectangle
		// Requests may use any combination of these parameters.
		// parameter	type	default	description
		// minlatitude	Decimal [-90,90] degrees	-90	Limit to events with a latitude larger than the specified minimum. NOTE: min values must be less than max values.
		// minlongitude	Decimal [-360,360] degrees	-180	Limit to events with a longitude larger than the specified minimum. NOTE: rectangles may cross the date line by using a minlongitude < -180 or maxlongitude > 180. NOTE: min values must be less than max values.
		// maxlatitude	Decimal [-90,90] degrees	90	Limit to events with a latitude smaller than the specified maximum. NOTE: min values must be less than max values.
		// maxlongitude	Decimal [-360,360] degrees	180	Limit to events with a longitude smaller than the specified maximum. NOTE: rectangles may cross the date line by using a minlongitude < -180 or maxlongitude > 180. NOTE: min values must be less than max values.
		MinLatitude  float64
		MinLongitude float64
		MaxLatitude  float64
		MaxLongitude float64

		// Circle
		// Requests must include all of latitude, longitude, and maxradius to perform a circle search.
		// parameter	type	default	description
		// latitude	Decimal [-90,90] degrees	null	Specify the latitude to be used for a radius search.
		// longitude	Decimal [-180,180] degrees	null	Specify the longitude to be used for a radius search.
		// maxradius	Decimal [0, 180] degrees	180	Limit to events within the specified maximum number of degrees from the geographic point defined by the latitude and longitude parameters. NOTE: This option is mutually exclusive with maxradiuskm and specifying both will result in an error.
		// maxradiuskm	Decimal [0, 20001.6] km	20001.6	Limit to events within the specified maximum number of kilometers from the geographic point defined by the latitude and longitude parameters. NOTE: This option is mutually exclusive with maxradius and specifying both will result in an error.
		Latitude    float64
		Longitude   float64
		MaxRadius   float64
		MaxRadiusKM float64

		// Other
		// parameter            type    default description
		// catalog              String  null    Limit to events from a specified catalog. Use the Catalogs Method to find available catalogs. NOTE: when catalog and contributor are omitted, the most preferred information from any catalog or contributor for the event is returned.
		// contributor          String  null    Limit to events contributed by a specified contributor. Use the Contributors Method to find available contributors. NOTE: when catalog and contributor are omitted, the most preferred information from any catalog or contributor for the event is returned.
		// eventid              String  null    Select a specific event by ID; event identifiers are data center specific. NOTE: Selecting a specific event implies includeallorigins, includeallmagnitudes, and, additionally, associated moment tensor and focal-mechanisms are included.
		// includeallmagnitudes Boolean false   Specify if all magnitudes for the event should be included, default is data center dependent but is suggested to be the preferred magnitude only. NOTE: because magnitudes and origins are strongly associated, this parameter is interchangeable with includeallmagnitudes
		// includeallorigins    Boolean false   Specify if all origins for the event should be included, default is data center dependent but is suggested to be the preferred origin only. NOTE: because magnitudes and origins are strongly associated, this parameter is interchangable with includeallmagnitudes
		// *includearrivals     Boolean false   Specify if phase arrivals should be included. NOTE: NOT CURRENTLY IMPLEMENTED
		// includedeleted       Boolean false   Specify if deleted products and events should be included. Deleted events otherwise return the HTTP status 409 Conflict.
		Catalog              Catalog
		Contributor          Contributor
		EventID              string
		IncludeAllMagnitudes bool
		IncludeAllOrigins    bool
		IncludeDeleted       bool

		// NOTE: Only supported by the csv and geojson formats, which include status.

		// includesuperseded Boolean  false	              Specify if superseded products should be included. This also includes all deleted products, and is mutually exclusive to the includedeleted parameter. NOTE: Only works when specifying eventid parameter.
		// limit             Integer  [1,20000]	null      Limit the results to the specified number of events. NOTE: The service limits queries to 20000, and any that exceed this limit will generate a HTTP response code “400 Bad Request”.
		// maxdepth          Decimal  [-100, 1000]km 1000 Limit to events with depth less than the specified maximum.
		// maxmagnitude      Decimal  null	              Limit to events with a magnitude smaller than the specified maximum.
		// mindepth          Decimal  [-100, 1000]km -100 Limit to events with depth more than the specified minimum.
		// minmagnitude      Decimal  null	              Limit to events with a magnitude larger than the specified minimum.
		// offset            Integer  [1,∞]	1	          Return results starting at the event count specified, starting at 1.
		// orderby           String	  [time,time-asc, magnitude, magnitude-asc] time Order the results.
		IncludeSuperseded bool
		Limit             int
		MaxDepth          float64
		MaxMagnitude      float64
		MinDepth          float64
		MinMagnitude      float64
		Offset            int
		OrderBy           Order

		// NOTE(jasonmoo): adding TotalResults as a way to handle paged querying, it is ignored in non-paged query
		TotalResults int

		// Extensions
		// parameter	type	default	description
		// alertlevel   String  null    [green, yellow, orange, red] Limit to events with a specific PAGER alert level.
		// eventtype    String  null    Limit to events of a specific type. NOTE: “earthquake” will filter non-earthquake events.
		// maxcdi       Decimal [0,12]  null	Maximum value for Maximum Community Determined Intensity reported by DYFI.
		// maxgap       Decimal [0,360] degrees	null	Limit to events with no more than this azimuthal gap.
		// maxmmi       Decimal [0,12]  null	Maximum value for Maximum Modified Mercalli Intensity reported by ShakeMap.
		// maxsig       Integer null    Limit to events with no more than this significance.
		// mincdi       Decimal null    Minimum value for Maximum Community Determined Intensity reported by DYFI.
		// minfelt      Integer [1,∞]   null	Limit to events with this many DYFI responses.
		// mingap       Decimal [0,360] degrees	null	Limit to events with no less than this azimuthal gap.
		// minsig       Integer null    Limit to events with no less than this significance.
		// nodata       Integer (204|404)	204	Define the error code that will be returned when no data is found.
		// producttype  String  null    Limit to events that have this type of product associated.
		//                              [moment-tensor, focal-mechanism, shakemap, losspager, dyfi]
		// productcode  String  null    Return the event that is associated with the productcode. The event will be returned even if the productcode is not the preferred code for the event. Example productcodes: nn00458749, at00ndf1fr
		// reviewstatus String  all     Limit to events with a specific review status. The different review statuses are:
		//                              [automatic,reviewed]
		AlertLevel   AlertLevel
		EventType    EventType
		MaxCdi       float64
		MaxGap       float64
		MaxMmi       float64
		MaxSig       int
		MinCdi       float64
		MinFelt      int
		MinGap       float64
		MinSig       int
		ProductType  ProductType
		ProductCode  string
		ReviewStatus ReviewStatus
	}

	Catalog       string
	Contributor   string
	EventType     string
	MagnitudeType string
	ProductType   string

	Order        string
	AlertLevel   string
	ReviewStatus string
)

const (
	OrderTimeDesc         Order        = "time"
	OrderTimeAsc          Order        = "time-asc"
	OrderMagnitudeDesc    Order        = "magnitude"
	OrderMagnitudeAsc     Order        = "magnitude-asc"
	AlertLevelGreen       AlertLevel   = "green"
	AlertLevelYellow      AlertLevel   = "yellow"
	AlertLevelOrange      AlertLevel   = "orange"
	AlertLevelRed         AlertLevel   = "red"
	ReviewStatusAll       ReviewStatus = "all"
	ReviewStatusAutomatic ReviewStatus = "automatic"
	ReviewStatusReviewed  ReviewStatus = "reviewed"
)

func NewQueryParameters() *queryParameters {
	return &queryParameters{
		format:       "geojson",
		MinLatitude:  math.NaN(),
		MinLongitude: math.NaN(),
		MaxLatitude:  math.NaN(),
		MaxLongitude: math.NaN(),
		Latitude:     math.NaN(),
		Longitude:    math.NaN(),
		MaxRadius:    math.NaN(),
		MaxRadiusKM:  math.NaN(),
		MaxDepth:     math.NaN(),
		MaxMagnitude: math.NaN(),
		MinDepth:     math.NaN(),
		MinMagnitude: math.NaN(),
		MaxCdi:       math.NaN(),
		MaxGap:       math.NaN(),
		MaxMmi:       math.NaN(),
		MinCdi:       math.NaN(),
		MinGap:       math.NaN(),
	}
}

func (qp *queryParameters) Encode() string {
	v := make(url.Values)
	v.Set("format", qp.format)
	if !qp.StartTime.IsZero() {
		v.Set("starttime", qp.StartTime.UTC().Format(time.RFC3339))
	}
	if !qp.EndTime.IsZero() {
		v.Set("endtime", qp.EndTime.UTC().Format(time.RFC3339))
	}
	if !qp.UpdatedAfter.IsZero() {
		v.Set("updatedafter", qp.UpdatedAfter.UTC().Format(time.RFC3339))
	}
	if !math.IsNaN(qp.MinLatitude) {
		v.Set("minlatitude", strconv.FormatFloat(qp.MinLatitude, 'f', -1, 64))
	}
	if !math.IsNaN(qp.MinLongitude) {
		v.Set("minlongitude", strconv.FormatFloat(qp.MinLongitude, 'f', -1, 64))
	}
	if !math.IsNaN(qp.MaxLatitude) {
		v.Set("maxlatitude", strconv.FormatFloat(qp.MaxLatitude, 'f', -1, 64))
	}
	if !math.IsNaN(qp.MaxLongitude) {
		v.Set("maxlongitude", strconv.FormatFloat(qp.MaxLongitude, 'f', -1, 64))
	}
	if !math.IsNaN(qp.Latitude) {
		v.Set("latitude", strconv.FormatFloat(qp.Latitude, 'f', -1, 64))
	}
	if !math.IsNaN(qp.Longitude) {
		v.Set("longitude", strconv.FormatFloat(qp.Longitude, 'f', -1, 64))
	}
	if !math.IsNaN(qp.MaxRadius) {
		v.Set("maxradius", strconv.FormatFloat(qp.MaxRadius, 'f', -1, 64))
	}
	if !math.IsNaN(qp.MaxRadiusKM) {
		v.Set("maxradiuskm", strconv.FormatFloat(qp.MaxRadiusKM, 'f', -1, 64))
	}
	if qp.Catalog != "" {
		v.Set("catalog", string(qp.Catalog))
	}
	if qp.Contributor != "" {
		v.Set("contributor", string(qp.Contributor))
	}
	if qp.EventID != "" {
		v.Set("eventid", string(qp.EventID))
	}
	if qp.IncludeAllMagnitudes {
		v.Set("includeallmagnitudes", "true")
	}
	if qp.IncludeAllOrigins {
		v.Set("includeallorigins", "true")
	}
	if qp.IncludeDeleted {
		v.Set("includedeleted", "true")
	}
	if qp.IncludeSuperseded {
		v.Set("includesuperseded", "true")
	}
	if qp.Limit != 0 {
		v.Set("limit", strconv.Itoa(qp.Limit))
	}
	if !math.IsNaN(qp.MaxDepth) {
		v.Set("maxdepth", strconv.FormatFloat(qp.MaxDepth, 'f', -1, 64))
	}
	if !math.IsNaN(qp.MaxMagnitude) {
		v.Set("maxmagnitude", strconv.FormatFloat(qp.MaxMagnitude, 'f', -1, 64))
	}
	if !math.IsNaN(qp.MinDepth) {
		v.Set("mindepth", strconv.FormatFloat(qp.MinDepth, 'f', -1, 64))
	}
	if !math.IsNaN(qp.MinMagnitude) {
		v.Set("minmagnitude", strconv.FormatFloat(qp.MinMagnitude, 'f', -1, 64))
	}
	if qp.Offset != 0 {
		v.Set("offset", strconv.Itoa(qp.Offset))
	}
	if qp.OrderBy != "" {
		v.Set("orderby", string(qp.OrderBy))
	}
	if qp.AlertLevel != "" {
		v.Set("alertlevel", string(qp.AlertLevel))
	}
	if qp.EventType != "" {
		v.Set("eventtype", string(qp.EventType))
	}
	if !math.IsNaN(qp.MaxCdi) {
		v.Set("maxcdi", strconv.FormatFloat(qp.MaxCdi, 'f', -1, 64))
	}
	if !math.IsNaN(qp.MaxGap) {
		v.Set("maxgap", strconv.FormatFloat(qp.MaxGap, 'f', -1, 64))
	}
	if !math.IsNaN(qp.MaxMmi) {
		v.Set("maxmmi", strconv.FormatFloat(qp.MaxMmi, 'f', -1, 64))
	}
	if qp.MaxSig != 0 {
		v.Set("maxsig", strconv.Itoa(qp.MaxSig))
	}
	if !math.IsNaN(qp.MinCdi) {
		v.Set("mincdi", strconv.FormatFloat(qp.MinCdi, 'f', -1, 64))
	}
	if qp.MinFelt != 0 {
		v.Set("minfelt", strconv.Itoa(qp.MinFelt))
	}
	if !math.IsNaN(qp.MinGap) {
		v.Set("mingap", strconv.FormatFloat(qp.MinGap, 'f', -1, 64))
	}
	if qp.MinSig != 0 {
		v.Set("minsig", strconv.Itoa(qp.MinSig))
	}
	if qp.ProductType != "" {
		v.Set("producttype", string(qp.ProductType))
	}
	if qp.ProductCode != "" {
		v.Set("productcode", string(qp.ProductCode))
	}
	if qp.ReviewStatus != "" {
		v.Set("reviewstatus", string(qp.ReviewStatus))
	}
	return v.Encode()
}
