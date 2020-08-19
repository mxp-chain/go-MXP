package goMXP

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ContractStorage(t *testing.T) {
	goldenStorage := []byte(`"Hello MXP!"`)
	goldenRPCErrors := readResponse(rpcerrors)
	type want struct {
		err         bool
		containsErr string
		rpcerr      []byte
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(storageHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"could not get storage",
				goldenRPCErrors,
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(storageHandlerMock(goldenStorage, blankHandler)),
			want{
				false,
				"",
				goldenStorage,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			rpcerr, err := gt.ContractStorage("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1", "KT1LfoE9EbpdsfUzowRckGUfikGcd5PyVKg")
			checkErr(t, tt.want.err, tt.containsErr, err)
			assert.Equal(t, tt.want.rpcerr, rpcerr)
		})
	}
}
