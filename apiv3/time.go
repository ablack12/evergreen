package apiv3

import "time"

const (
	// This string defines ISO-8601 UTC with 3 fractional seconds behind a dot
	// specified by the API spec document.
	APITimeFormat = "\"2006-01-02T15:04:05.000Z\""
)

type APITime time.Time

// NewTime creates a new APITime from an existing time.Time. It handles changing
// converting from the times time zone to UTC.
func NewTime(t time.Time) APITime {
	utcT := t.In(time.UTC)
	return APITime(time.Date(
		utcT.Year(),
		utcT.Month(),
		utcT.Day(),
		utcT.Hour(),
		utcT.Minute(),
		utcT.Second(),
		utcT.Nanosecond(),
		time.FixedZone("", 0),
	))
}

// UnmarshalJSON implements the custom unmarshalling of this type so that it can
// be correctly parsed from an API request.
func (at *APITime) UnmarshalJSON(b []byte) error {
	t, err := time.ParseInLocation(APITimeFormat, string(b), time.FixedZone("", 0))
	if err != nil {
		return err
	}
	(*at) = APITime(t)
	return nil

}

// MarshalJSON implements the custom marshalling of this type so that it can
// be correctly written out in an API response.
func (at APITime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(at).Format(APITimeFormat)), nil
}

func (at APITime) String() string {
	return time.Time(at).String()
}
