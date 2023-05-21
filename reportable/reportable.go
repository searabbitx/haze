package reportable

import (
	"github.com/kamil-s-solecki/haze/http"
)

func IsReportable(res http.Response) bool {
	return res.Code == 500
}
