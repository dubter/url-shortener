package request

import (
	"net"
	"net/url"
	"regexp"
	"url-shortner/internal/domain/validation"
)

// URLInput defines structure for create short code url request
type URLInput struct {
	URL  string `json:"url" binding:"required"`
	Host string `json:"-"`
}

// URLFilter defines structure for short code list and search request
type URLFilter struct {
	ShortCode string `json:"short_code"`
	Keyword   string `json:"keyword"`
	Page      string `json:"page"`
}

// @see https://github.com/asaskevich/govalidator/blob/master/patterns.go
var (
	IP           = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	URLSchema    = `((ftp|https?):\/\/)`
	URLUsername  = `(\S+(:\S*)?@)`
	URLPath      = `((\/|\?|#)[^\s]*)`
	URLPort      = `(:(\d{1,5}))`
	URLIP        = `([1-9]\d?|1\d\d|2[01]\d|22[0-3]|24\d|25[0-5])(\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-5]))`
	URLSubdomain = `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`

	URLMinLength   = 15
	URLMaxLength   = 2048
	URLRegex       = `^` + URLSchema + `?` + URLUsername + `?` + `((` + URLIP + `|(\[` + IP + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + URLSubdomain + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + URLPort + `?` + URLPath + `?$`
	URLFilterRegex = `(xxx|localhost|127\.0\.0\.1|\.local)`
)

var (
	urlRe    = regexp.MustCompile(URLRegex)
	filterRe = regexp.MustCompile(URLFilterRegex)
)

// Validate validates the url input before saving to db
// It returns error if something is not valid.
func (input *URLInput) Validate() error {
	if l := len(input.URL); l < URLMinLength || l > URLMaxLength {
		return validation.ErrInvalidURLLen
	}

	if filterRe.MatchString(input.URL) {
		return validation.ErrFilteredURL
	}

	uri, err := url.ParseRequestURI(input.URL)
	if err != nil {
		return validation.ErrInvalidURL
	}

	input.Host = uri.Host
	if host, _, _ := net.SplitHostPort(uri.Host); host != "" {
		input.Host = host
	}

	if !urlRe.MatchString(input.URL) {
		return validation.ErrInvalidURL
	}

	return nil
}
