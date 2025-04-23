package params

import "github.com/sviatilnik/url-shortener/internal/app/util"

type URLConfigParam struct {
	value string
}

func NewURLConfigParam(URL string) URLConfigParam {
	return URLConfigParam{
		value: URL,
	}
}

func (p URLConfigParam) Validate() bool {
	return util.IsURL(p.value)
}

func (p URLConfigParam) Value() interface{} {
	return p.value
}
