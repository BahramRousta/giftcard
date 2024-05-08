package requester

import "net/url"

type Request struct {
	ID          string
	RequestBody string
	UserIP      string
	Uri         string
	Method      string
	Host        string
	Header      map[string][]string
	Params      url.Values
}
