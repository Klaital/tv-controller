package vlcclient

import "fmt"

type VlcStatusResponse struct {
	Fullscreen   bool `json:"fullscreen"`
	SeekSec      int  `json:"seek_sec"`
	ApiVersion   int  `json:"apiversion"`
	CurrentPlId  int  `json:"currentplid"`
	Time         int  `json:"time"`
	Volume       int  `json:"volume"`
	Length       int  `json:"length"`
	Random       bool `json:"random"`
	AudioFilters struct {
		Filter0 string `json:"filter_0"`
	} `json:"audiofilters"`
	Rate         int `json:"rate"`
	VideoEffects struct {
		Hue        int `json:"hue"`
		Saturation int `json:"saturation"`
		Contrast   int `json:"contrast"`
		Brightness int `json:"brightness"`
		Gamma      int `json:"gamma"`
	} `json:"videoeffects"`
	State         string  `json:"state"`
	Loop          bool    `json:"loop"`
	Version       string  `json:"version"`
	Position      float64 `json:"position"`
	AudioDelay    int     `json:"audiodelay"`
	Repeat        bool    `json:"repeat"`
	SubtitleDelay int     `json:"subtitledelay"`
	Equalizer     []any   `json:"equalizer"`
}

func (c Client) GetStatus() (*VlcStatusResponse, error) {
	var resp VlcStatusResponse
	err := c.Do("/requests/status.json", nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("GetStatus: %w", err)
	}
	return &resp, nil
}
