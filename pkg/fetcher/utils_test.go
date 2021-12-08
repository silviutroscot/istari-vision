package fetcher

import (
	"encoding/json"
	"math/big"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockHandler struct {
	responseFunc func(r *http.Request) ([]byte, int, error)
}

func (mh *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if errorCode := query.Get("errorCode"); errorCode != "" {
		errorCodeInt, _ := strconv.Atoi(errorCode)
		w.WriteHeader(errorCodeInt)
		w.Write([]byte("error on purpose"))
		return
	}

	payload, status, err := mh.responseFunc(r)
	if err != nil {
		w.WriteHeader(status)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(status)
	w.Write(payload)
}

func Test_JSON_big(t *testing.T) {
	t.Parallel()

	f := big.NewFloat(3.1415)
	data, err := json.Marshal(f)
	require.NoError(t, err)
	t.Logf("%s", string(data))

	var fBack *big.Float
	err = json.Unmarshal(data, &fBack)
	require.NoError(t, err)
	require.Equal(t, f.String(), fBack.String())

	var f2 struct {
		Value *big.Float `json:"float"`
	}

	f2.Value = big.NewFloat(3.14156)

	data, err = json.Marshal(f2)
	require.NoError(t, err)
	t.Logf("%s", string(data))
}
