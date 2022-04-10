package sitemap

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"

	"github.com/Shamanskiy/gophercises/testutils"
)

var htmlTemplate string = `
<html>
  <body>
    {{ range .}}
      <a href="{{.}}">Some link text</a>
    {{ end }}
  </body>
</html>
`

var responseWithLinks *template.Template = template.Must(template.New("").Parse(htmlTemplate))

func TestMapBuilder(t *testing.T) {
	links := []string{"/home", "/about"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Server: GET request %v\n", r.URL)
		responseWithLinks.Execute(w, links)
	}))
	defer server.Close()

	want := []string{
		server.URL + "/home",
		server.URL + "/about",
	}
	got, err := ParseSite(server.URL)

	testutils.CheckError(err, t)
	compareSiteMaps(got, want, t)
}

func compareSiteMaps(got, want []string, t *testing.T) {
	t.Helper()
	if !testutils.SameElements(got, want) {
		testutils.ReportDifferentSlices(got, want, "Different sitemaps!", t)
	}
}