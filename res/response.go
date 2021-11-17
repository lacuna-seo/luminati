package res

import "time"

type Response struct {
	General struct {
		SearchEngine string    `json:"search_engine"`
		Query        string    `json:"query"`
		ResultsCnt   int       `json:"results_cnt"`
		SearchTime   float64   `json:"search_time"`
		Language     string    `json:"language"`
		Location     string    `json:"location"`
		Mobile       bool      `json:"mobile"`
		BasicView    bool      `json:"basic_view"`
		SearchType   string    `json:"search_type"`
		PageTitle    string    `json:"page_title"`
		CodeVersion  string    `json:"code_version"`
		Timestamp    time.Time `json:"timestamp"`
	} `json:"general"`
	Input struct {
		OriginalURL string `json:"original_url"`
	} `json:"input"`
	Organic []struct {
		Link        string `json:"link"`
		DisplayLink string `json:"display_link"`
		Title       string `json:"title"`
		Description string `json:"description"`
		CachedLink  string `json:"cached_link,omitempty"`
		Extensions  []struct {
			Inline bool   `json:"inline"`
			Type   string `json:"type"`
			Text   string `json:"text"`
			Rank   int    `json:"rank"`
			Link   string `json:"link,omitempty"`
		} `json:"extensions,omitempty"`
		InfoDescription string `json:"info_description,omitempty"`
		InfoLogo        string `json:"info_logo,omitempty"`
		Rank            int    `json:"rank"`
		GlobalRank      int    `json:"global_rank"`
		SimilarLink     string `json:"similar_link,omitempty"`
		InfoSource      string `json:"info_source,omitempty"`
		InfoLink        string `json:"info_link,omitempty"`
		Image           string `json:"image,omitempty"`
		ImageAlt        string `json:"image_alt,omitempty"`
		ImageURL        string `json:"image_url,omitempty"`
		Duration        string `json:"duration,omitempty"`
		DurationSec     int    `json:"duration_sec,omitempty"`
	} `json:"organic"`
	Images []struct {
		Image      string `json:"image"`
		ImageAlt   string `json:"image_alt"`
		ImageURL   string `json:"image_url"`
		Rank       int    `json:"rank"`
		GlobalRank int    `json:"global_rank"`
	} `json:"images"`
	SnackPackMap struct {
		Image       string  `json:"image"`
		ImageAlt    string  `json:"image_alt"`
		ImageBase64 string  `json:"image_base64"`
		Link        string  `json:"link"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		Altitude    int     `json:"altitude"`
	} `json:"snack_pack_map"`
	SnackPack []struct {
		Cid               string   `json:"cid"`
		Name              string   `json:"name"`
		Image             string   `json:"image"`
		ImageBase64       string   `json:"image_base64"`
		Rating            float64  `json:"rating"`
		ReviewsCnt        int      `json:"reviews_cnt"`
		Type              string   `json:"type"`
		Price             string   `json:"price"`
		WorkStatus        string   `json:"work_status"`
		WorkStatusDetails string   `json:"work_status_details"`
		Address           string   `json:"address"`
		Tags              []string `json:"tags"`
		Rank              int      `json:"rank"`
		GlobalRank        int      `json:"global_rank"`
	} `json:"snack_pack"`
	Knowledge struct {
		Name      string `json:"name"`
		Summary   string `json:"summary"`
		Address   string `json:"address"`
		OpenHours []struct {
			Day   string `json:"day"`
			Hours string `json:"hours"`
		} `json:"open_hours"`
		Phone               string  `json:"phone"`
		Site                string  `json:"site"`
		Fid                 string  `json:"fid"`
		ReviewsCnt          int     `json:"reviews_cnt"`
		MapsLink            string  `json:"maps_link"`
		Latitude            float64 `json:"latitude"`
		Longitude           float64 `json:"longitude"`
		Zoom                int     `json:"zoom"`
		MerchantDescription string  `json:"merchant_description"`
		Facts               []struct {
			Key       string `json:"key"`
			KeyLink   string `json:"key_link,omitempty"`
			Predicate string `json:"predicate"`
			Value     []struct {
				Text string `json:"text"`
			} `json:"value"`
		} `json:"facts"`
		Widgets []struct {
			Type      string `json:"type"`
			Key       string `json:"key"`
			Predicate string `json:"predicate"`
			Items     []struct {
				Link string `json:"link"`
				Name string `json:"name"`
				Rank int    `json:"rank"`
			} `json:"items"`
			Rank       int `json:"rank"`
			GlobalRank int `json:"global_rank"`
		} `json:"widgets"`
	} `json:"knowledge"`
	BottomAds []struct {
		Link         string `json:"link"`
		ReferralLink string `json:"referral_link"`
		Title        string `json:"title"`
		Phone        string `json:"phone,omitempty"`
		Description  string `json:"description"`
		Extensions   []struct {
			Type string `json:"type"`
			Link string `json:"link"`
			Text string `json:"text"`
		} `json:"extensions,omitempty"`
		Rank       int `json:"rank"`
		GlobalRank int `json:"global_rank"`
	} `json:"bottom_ads"`
	Pagination struct {
		CurrentPage   int    `json:"current_page"`
		NextPageLink  string `json:"next_page_link"`
		NextPageStart int    `json:"next_page_start"`
		NextPage      int    `json:"next_page"`
		Pages         []struct {
			Page  int    `json:"page"`
			Link  string `json:"link"`
			Start int    `json:"start"`
		} `json:"pages"`
	} `json:"pagination"`
	Related []struct {
		Text       string `json:"text"`
		ListGroup  bool   `json:"list_group"`
		Expanded   bool   `json:"expanded,omitempty"`
		Image      string `json:"image,omitempty"`
		ImageAlt   string `json:"image_alt,omitempty"`
		ImageURL   string `json:"image_url,omitempty"`
		Rank       int    `json:"rank"`
		GlobalRank int    `json:"global_rank"`
		Link       string `json:"link,omitempty"`
	} `json:"related"`
	PeopleAlsoAsk []struct {
		Question          string `json:"question"`
		QuestionLink      string `json:"question_link"`
		AnswerSource      string `json:"answer_source"`
		AnswerLink        string `json:"answer_link"`
		AnswerDisplayLink string `json:"answer_display_link"`
		AnswerHTML        string `json:"answer_html"`
		Answers           []struct {
			Type  string `json:"type"`
			Title string `json:"title"`
			Items []struct {
				Value string `json:"value"`
				Rank  int    `json:"rank"`
			} `json:"items"`
			Rank int `json:"rank"`
		} `json:"answers"`
		Rank       int `json:"rank"`
		GlobalRank int `json:"global_rank"`
	} `json:"people_also_ask"`
}
