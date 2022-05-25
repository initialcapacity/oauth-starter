package websupport_test

import (
  "bytes"
  "github.com/initialcapacity/oauth-starter/pkg/websupport"
  "github.com/initialcapacity/oauth-starter/pkg/websupport/test"

  "github.com/stretchr/testify/assert"
  "io"
  "net/http/httptest"
  "testing"
)

func TestModelAndView(t *testing.T) {
  w := &httptest.ResponseRecorder{Body: new(bytes.Buffer)}
  data := websupport.Model{Map: map[string]any{"name": "aName"}}
  _ = websupport.ModelAndView(w, &websupport_test.Resources, "test", data)
  body, _ := io.ReadAll(w.Body)
  assert.Equal(t, `
    <!DOCTYPE html>
    <html lang="en">
    <body>aName
    </body>
    </html>`, string(body))
}
