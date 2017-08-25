package feedgetter

import (
	"net/http"
	"testing"
)

func TestEncode(t *testing.T) {
	h := &http.Header{}
	h.Set("content-type", "text/html; charset=euc-jp")
	enc := encode(h, []byte{})
	if enc != "euc-jp" {
		t.Error("encode from content-type is error.")
	}
}
