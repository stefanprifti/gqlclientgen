package countries

type Continent struct {
	Code      string    `json:"code"`
	Countries []Country `json:"countries"`
	Name      string    `json:"name"`
}
type ContinentFilterInput struct {
	Code StringQueryOperatorInput `json:"code,omitempty"`
}
type Country struct {
	AwsRegion  string     `json:"awsRegion"`
	Capital    string     `json:"capital,omitempty"`
	Code       string     `json:"code"`
	Continent  Continent  `json:"continent"`
	Currencies []string   `json:"currencies"`
	Currency   string     `json:"currency,omitempty"`
	Emoji      string     `json:"emoji"`
	EmojiU     string     `json:"emojiU"`
	Languages  []Language `json:"languages"`
	Name       string     `json:"name"`
	Native     string     `json:"native"`
	Phone      string     `json:"phone"`
	Phones     []string   `json:"phones"`
	States     []State    `json:"states"`
}
type CountryFilterInput struct {
	Code      StringQueryOperatorInput `json:"code,omitempty"`
	Continent StringQueryOperatorInput `json:"continent,omitempty"`
	Currency  StringQueryOperatorInput `json:"currency,omitempty"`
}
type Language struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Native string `json:"native"`
	Rtl    bool   `json:"rtl"`
}
type LanguageFilterInput struct {
	Code StringQueryOperatorInput `json:"code,omitempty"`
}
type State struct {
	Code    string  `json:"code,omitempty"`
	Country Country `json:"country"`
	Name    string  `json:"name"`
}
type StringQueryOperatorInput struct {
	Eq    string   `json:"eq,omitempty"`
	In    []string `json:"in,omitempty"`
	Ne    string   `json:"ne,omitempty"`
	Nin   []string `json:"nin,omitempty"`
	Regex string   `json:"regex,omitempty"`
}
type CountryRequest struct {
	Code string `json:"code"`
}
type CountryResponse struct {
	Country struct {
		Name      string `json:"name"`
		Native    string `json:"native"`
		Languages struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"languages"`
		Emoji     string `json:"emoji"`
		Currency  string `json:"currency,omitempty"`
		Languages struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"languages"`
	} `json:"country,omitempty"`
}
