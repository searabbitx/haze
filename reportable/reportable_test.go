package reportable

import (
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/testutils"
	"testing"
)

func TestShouldNotReport200(t *testing.T) {
	res := http.Response{200, []byte{}}

	got := IsReportable(res)

	testutils.AssertFalse(t, got)
}

func TestShouldReport500(t *testing.T) {
	res := http.Response{500, []byte{}}

	got := IsReportable(res)

	testutils.AssertTrue(t, got)
}
