package vlcclient

type VlcStatusResponse struct {
	Fullscreen   int  `json:"fullscreen"`
	SeekSec      int  `json:"seek_sec"`
	Apiversion   int  `json:"apiversion"`
	Currentplid  int  `json:"currentplid"`
	Time         int  `json:"time"`
	Volume       int  `json:"volume"`
	Length       int  `json:"length"`
	Random       bool `json:"random"`
	Audiofilters struct {
		Filter0 string `json:"filter_0"`
	} `json:"audiofilters"`
	Rate         int `json:"rate"`
	Videoeffects struct {
		Hue        int `json:"hue"`
		Saturation int `json:"saturation"`
		Contrast   int `json:"contrast"`
		Brightness int `json:"brightness"`
		Gamma      int `json:"gamma"`
	} `json:"videoeffects"`
	State         string `json:"state"`
	Loop          bool   `json:"loop"`
	Version       string `json:"version"`
	Position      int    `json:"position"`
	Audiodelay    int    `json:"audiodelay"`
	Repeat        bool   `json:"repeat"`
	Subtitledelay int    `json:"subtitledelay"`
	Equalizer     []any  `json:"equalizer"`
}

//func (c Client) GetStatus() (VlcStatusResponse, error) {
//
//}
