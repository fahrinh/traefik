package plugin

import (
	"net/http"
)

type IMatcher interface {
	MatcherFunc (req *http.Request) bool
}


