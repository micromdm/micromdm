package mdm

import (
	"fmt"
	"github.com/groob/plist"
	"github.com/micromdm/mdm/test"
	"io/ioutil"
	"testing"
)

func TestQueryResponseMac(t *testing.T) {
	response := &Response{}
	plistBuf := []byte(test.OSXElCapQueryResponses)
	err := plist.Unmarshal(plistBuf, response)

	if err != nil {
		t.Fatal(err)
	}
}

func TestQueryResponseIpadIOS8(t *testing.T) {
	response := &Response{}
	plistBuf := []byte(test.IOS8IpadQueryResponses)
	err := plist.Unmarshal(plistBuf, response)

	if err != nil {
		t.Fatal(err)
	}
}

func TestQueryResponseIphoneIOS8(t *testing.T) {
	response := &Response{}
	plistBuf := []byte(test.IOS8IphoneQueryResponses)
	err := plist.Unmarshal(plistBuf, response)

	if err != nil {
		t.Fatal(err)
	}
}

func TestSecurityInfoMac(t *testing.T) {
	response := &Response{}
	plistBuf := []byte(test.OSXElCapSecurityInfoNoFDE)
	err := plist.Unmarshal(plistBuf, response)

	if err != nil {
		t.Fatal(err)
	}
}

func TestSecurityInfoIpadIOS8(t *testing.T) {
	response := &Response{}
	plistBuf := []byte(test.IOS8IpadSecurityInfo)
	err := plist.Unmarshal(plistBuf, response)

	if err != nil {
		t.Fatal(err)
	}
}

func TestErrorChain(t *testing.T) {
	errorResponseBody, err := ioutil.ReadFile("./test/responses/error_invalidreq.plist")

	if err != nil {
		t.Fatal(err)
	}

	response := &Response{}
	if err := plist.Unmarshal(errorResponseBody, response); err != nil {
		t.Fatal(err)
	}

	if response.Status != "Error" {
		t.Fatal("Response status was not `Error`.")
	}

	if response.ErrorChain == nil {
		t.Fatal("Response did not contain expected error chain struct")
	}

	for _, v := range response.ErrorChain {
		if v.ErrorCode == 0 {
			t.Fatal("Error response did not contain ErrorCode")
		}

		if v.ErrorDomain == "" {
			t.Fatal("Error response did not contain ErrorDomain")
		}

		if v.LocalizedDescription == "" {
			t.Fatal("Error response did not contain localized description")
		}

		if v.USEnglishDescription == "" {
			t.Fatal("Error response did not contain english description")
		}
	}
}

func TestInstalledApplicationListResponse(t *testing.T) {
	appListResponseBody, err := ioutil.ReadFile("./test/responses/installed_application_list.plist")

	if err != nil {
		t.Fatal(err)
	}

	response := &Response{}
	if err := plist.Unmarshal(appListResponseBody, response); err != nil {
		t.Fatal(err)
	}

	if response.Status != "Acknowledged" {
		t.Fatal("Response status was not `Acknowledged`.")
	}

	if response.InstalledApplicationList == nil {
		t.Fatal("No installed application list in response")
	}

	fmt.Printf("%v\n", response)
}
