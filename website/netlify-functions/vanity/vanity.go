// Binary vanity serves go-import and go-source URLs for Go module imports.
package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// See "go help importpath" or https://golang.org/cmd/go/#hdr-Remote_import_paths
var vanityTemplate = template.Must(template.New("vanity.html").Parse(`
<head>
  <meta name="go-import" content="{{ .Permalink }} git {{ .Repo }}">
  <meta name="go-source" content="{{ .Permalink }} {{ .Repo }} {{ .Repo }}/tree/{{ .Tree }}{/dir} {{ .Repo }}/blob/{{ .Tree }}{/dir}/{file}#L{line}">
</head>
`))

type server struct {
	logger *log.Logger
	myURL  *url.URL
}

func (srv server) handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	srv.logger.Printf("request_path=%s", request.Path)

	pURL, err := srv.myURL.Parse(request.Path)
	if err != nil {
		srv.logger.Printf("parsing relative to base url: %v", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "oops",
		}, nil
	}

	var vanity = struct {
		Permalink string
		Repo      string
		Tree      string
	}{
		Permalink: strings.TrimPrefix(pURL.String(), "https://"),
		Repo:      "https://github.com/micromdm/micromdm",
		Tree:      "v2dev",
	}

	srv.logger.Printf("request_path=%s permalink=%s", request.Path, vanity.Permalink)

	buf := new(bytes.Buffer)
	if err := vanityTemplate.Execute(buf, vanity); err != nil {
		srv.logger.Printf("execute vanityTemplate: %v", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "oops",
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "text/html; charset=UTF-8"},
		Body:       buf.String(),
	}, nil
}

func main() {
	myURL, _ := url.Parse("https://micromdm.io/")
	srv := server{
		logger: log.New(os.Stderr, "", log.LstdFlags),
		myURL:  myURL,
	}
	lambda.Start(srv.handler)
}
