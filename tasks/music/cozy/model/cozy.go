package model

type Cozy struct {
	CozyId      string `json:"cozyId"`
	Headline    string `json:"headline"`
	Subtitle    string `json:"subtitle"`
	Category    string `json:"category"`
	ImageUrl    string `json:"imageUrl"`
	MinImageUrl string `json:"minImageUrl"`
	IconUrl     string `json:"iconUrl"`
	AudioUrl    string `json:"audioUrl"`
	Free        bool   `json:"free"`
	Recommend   bool   `json:"recommend"`
	CoverHash   string `json:"cover_hash"`
	ResHash     string `json:"res_hash"`
	CoverLink   string `json:"cover_link"`
	ResLink     string `json:"res_link"`
}
