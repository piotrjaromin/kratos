package attack

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestCreateTargeterShouldCreateRequestFromJsScript(t *testing.T) {
	url := "http://localhost/test"
	method := "GET"
	body := "string body"

	testFile := fmt.Sprintf(`
		function getRequestOptions() {
			return {
				url: "%s",
				method: "%s",
				body: "%s",
				headers: {
					"header-1": "val-header-1",
					"header-2": "val-header-2",
				}
			}
		}
	`, url, method, body)

	testProvider := func() []byte { return []byte(testFile) }

	targeter := CreateTargeter(testProvider)
	require.NotNil(t, targeter)

	target := &vegeta.Target{}

	err := targeter(target)
	require.Nil(t, err)

	require.Equal(t, url, target.URL)
	require.Equal(t, method, target.Method)
	require.Equal(t, []byte(body), target.Body)

	expectedHeaders := make(http.Header)
	expectedHeaders.Set("header-1", "val-header-1")
	expectedHeaders.Set("header-2", "val-header-2")

	require.Equal(t, expectedHeaders, target.Header)
}

func TestCreateTargeterShouldCallGetRequestOptions(t *testing.T) {
	url := "http://localhost/test/"

	testFile := fmt.Sprintf(`
		var count = 0;
		function getRequestOptions() {
			var opts = {
				url: "%s" + count,
				method: "GET",
			}

			count++;
			return opts;
		}
	`, url)

	testProvider := func() []byte { return []byte(testFile) }

	targeter := CreateTargeter(testProvider)
	require.NotNil(t, targeter)

	assertTarget := func(expectedUrl string) {
		target := &vegeta.Target{}
		err := targeter(target)
		require.Nil(t, err)

		require.Equal(t, expectedUrl, target.URL)
		require.Equal(t, []byte(nil), target.Body)

		require.Equal(t, make(http.Header), target.Header)
	}

	assertTarget(fmt.Sprintf("%s0", url))
	assertTarget(fmt.Sprintf("%s1", url))
	assertTarget(fmt.Sprintf("%s2", url))
}
