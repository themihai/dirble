package dirble

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Station struct {
	ID             int    `json:"id"`       //e.g. 15
	Name           string `json:"name"`     // Neradio House  Trance
	Accepted       int    `json:"accepted"` // 1,
	Country        string `json:"country"`  // "SE",
	Description    string `json:"description"`
	TotalListeners int    `json:"total_listeners"` // 20,
	Image          Image  `json:"image"`

	Slug    string `json:"slug"`    // neradio-house-trance",
	Website string `json:"website"` //
	//created_at": "2012-01-15T07:49:54.000+01:00",
	//"updated_at": "2015-04-11T14:10:45.000+02:00",
	Streams    []Stream   `json:"streams"`
	Categories []Category `json:"categories"`
}

type Category struct {
	ID          int    `json:"id"`          //: 1,
	Title       string `json:"title"`       //: "Trance",
	Description string `json:"description"` //: "stations that plays commercial and other things in trance-music genre.",
	Slug        string `json:"slug"`        //: "trance",
	Ancestry    string `json:"ancestry"`    //: "14",
}

type Image struct {
	URL   string `json:"url"`
	Thumb Thumb  `json:"thumb"`
}

type Thumb struct {
	URL string `json:"url"`
}

type Stream struct {
	Stream      string `json:"stream"`       // "http://fire1.neradio.com",
	Bitrate     int    `json:"bitrate"`      // 128,
	ContentType string `json:"content_type"` // "audio/mpeg\r\n",
	Listeners   int    `json:"listeners"`    // 10,
	Status      int    `json:"status"`
}

func StationByID(cx context.Context, token string, stationID int) (*Station, error) {
	URL := "http://api.dirble.com/v2/station/" + strconv.Itoa(stationID) +
		"?token=" + token
	b, err := get(cx, URL)
	if err != nil {
		return nil, err
	}
	v := new(Station)
	if err = json.Unmarshal(b, v); err != nil {
		return nil, fmt.Errorf("err %v, \ndata %s", err, b)
	}
	return v, nil
}

// @limit - How many stations per page to show
// @page - show which per_page stations to show(page number?)
// @offset
func StationsByCountry(cx context.Context, token, country string, page, offset, limit int) ([]Station, error) {
	URL := "http://api.dirble.com/v2/countries/" + country +
		"/stations?token=" + token
	b, err := get(cx, URL)
	if err != nil {
		return nil, err
	}
	var v []Station
	if err = json.Unmarshal(b, &v); err != nil {
		return nil, fmt.Errorf("err %v, \ndata %s", err, b)
	}
	return v, nil
}

func get(cx context.Context, URL string) ([]byte, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(cx)
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("Status %v, body %s", rsp.StatusCode, b)
	}
	return b, nil
}
