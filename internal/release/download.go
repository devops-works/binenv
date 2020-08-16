package release

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// Download handles direct binary releases
type Download struct {
	url string
}

// Fetch gets the package and returns location of downloaded file
func (d Download) Fetch(ctx context.Context, v string) (string, error) {
	type fullVersion struct {
		OS           string
		Arch         string
		Version      string
		NakedVersion string
	}

	fv := fullVersion{
		Arch:         runtime.GOARCH,
		OS:           runtime.GOOS,
		Version:      v,
		NakedVersion: strings.TrimLeft(v, "vV"),
	}

	tmpl, err := template.New("test").Parse(d.url)
	if err != nil {
		panic(err)
	}

	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, fv)
	if err != nil {
		panic(err)
	}

	url := buf.String()

	fmt.Printf("fetching version %q for arch %q and OS %q at %s\n", v, runtime.GOARCH, runtime.GOOS, url)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpfile, err := ioutil.TempFile("", v)
	if err != nil {
		log.Fatal(err)
	}

	defer tmpfile.Close()

	// Write the body to file
	_, err = io.Copy(tmpfile, resp.Body)

	fmt.Printf("file saved at %q\n", tmpfile.Name())

	return tmpfile.Name(), err
}
