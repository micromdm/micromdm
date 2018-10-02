package appstore

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"testing"
)

func TestDecodeUploadRequest(t *testing.T) {
	var aFile, bFile bytes.Buffer
	aFile.Write([]byte("hello"))
	bFile.Write([]byte("world"))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("app_manifest_filedata", "manifest.plist")
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, &aFile)
	if err != nil {
		t.Fatal(err)
	}
	err = writer.WriteField("app_manifest_filename", "manifest.plist")
	if err != nil {
		t.Fatal(err)
	}

	partPkg, err := writer.CreateFormFile("pkg_filedata", "foo.pkg")
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(partPkg, &bFile)
	if err != nil {
		t.Fatal(err)
	}
	err = writer.WriteField("pkg_name", "hello.pkg")
	if err != nil {
		t.Fatal(err)
	}
	writer.Close()

	req := httptest.NewRequest("POST", "https://mdm.acme.co/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	request, err := decodeAppUploadRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	decoded := request.(appUploadRequest)

	if have, want := decoded.ManifestName, "manifest.plist"; have != want {
		t.Errorf("have %s, want %s", have, want)
	}
	if have, want := decoded.PKGFilename, "hello.pkg"; have != want {
		t.Errorf("have %s, want %s", have, want)
	}

	var a, b bytes.Buffer
	_, err = io.Copy(&a, decoded.ManifestFile)
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(&b, decoded.PKGFile)
	if err != nil {
		t.Fatal(err)
	}
	if have, want := a.String(), "hello"; have != want {
		t.Errorf("have %s, want %s", have, want)
	}
	if have, want := b.String(), "world"; have != want {
		t.Errorf("have %s, want %s", have, want)
	}
}
