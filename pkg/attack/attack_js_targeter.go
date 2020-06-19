package attack

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/robertkrimen/otto"
	vegeta "github.com/tsenart/vegeta/v12/lib"

	"github.com/thoas/go-funk"
)

type TestProvider func() []byte

func TestFileProvider(fileName string) (TestProvider, error) {
	contents, err := ioutil.ReadFile(fileName)
	return func() []byte {
		return contents
	}, err
}

func CreateTargeter(testFile TestProvider) vegeta.Targeter {
	content := testFile()
	script := string(content)
	vm := otto.New()

	firstRun := true
	return func(t *vegeta.Target) (err error) {
		if firstRun {
			if _, err := vm.Run(script); err != nil {
				return fmt.Errorf("Unable to parse test file. Details: %s", err.Error())
			}

			firstRun = false
		}

		value, err := vm.Run(`getRequestOptions();`)
		if err != nil {
			return fmt.Errorf("'getRequestOptions()' run fail. Details: %s", err.Error())
		}

		reqOpts := value.Object()
		if err != nil {
			return fmt.Errorf("'getRequestOptions()' did not return object. Details: %s", err.Error())
		}

		if err := setStringField(t, "Method", reqOpts, "method"); err != nil {
			return fmt.Errorf("Invalid 'method' returned by getRequestOptions(), expected string value. Details: %s", err.Error())
		}

		if err := setStringField(t, "URL", reqOpts, "url"); err != nil {
			return fmt.Errorf("Invalid 'url' returned by getRequestOptions(), expected string value. Details: %s", err.Error())
		}

		if funk.ContainsString(reqOpts.Keys(), "body") {
			body, err := getStringValue(reqOpts, "body")

			if err != nil {
				return fmt.Errorf("Invalid 'body' returned by getRequestOptions(), expected string value. Details: %s", err.Error())
			}

			t.Body = []byte(body)
		}

		t.Header = make(http.Header)
		if funk.ContainsString(reqOpts.Keys(), "headers") {
			headersValue, err := reqOpts.Get("headers")
			if err != nil {
				return fmt.Errorf("Error wile geting headers object. Details: %s", err.Error())
			}

			if !headersValue.IsObject() {
				return fmt.Errorf("headers field should be an object")
			}

			headersObj := headersValue.Object()
			for _, headerName := range headersObj.Keys() {
				headerValue, err := getStringValue(headersObj, headerName)
				if err != nil {
					return fmt.Errorf("Invalid 'header.%s' returned by getRequestOptions(), expected string value. Details: %s", headerName, err.Error())
				}

				t.Header.Set(headerName, headerValue)
			}
		}

		return err
	}
}

func setStringField(dst *vegeta.Target, dstFieldName string, src *otto.Object, srcFieldName string) error {
	strValue, err := getStringValue(src, srcFieldName)
	if err != nil {
		return err
	}

	reflect.ValueOf(dst).Elem().FieldByName(dstFieldName).SetString(strValue)
	return nil
}

func getStringValue(src *otto.Object, srcFieldName string) (string, error) {
	field, err := src.Get(srcFieldName)
	if err != nil {
		return "", err
	}

	if !field.IsString() {
		return "", fmt.Errorf("'%s' field is not an string", srcFieldName)
	}

	return field.String(), nil
}
