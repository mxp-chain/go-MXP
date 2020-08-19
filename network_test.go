package goMXP

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Version(t *testing.T) {
	goldenVersion := getResponse(version).(Version)

	type want struct {
		wantErr     bool
		containsErr string
		wantVersion Version
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(versionsHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"could not get network version",
				Version{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(versionsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal network version",
				Version{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(versionsHandlerMock(readResponse(version), blankHandler)),
			want{
				false,
				"",
				goldenVersion,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			version, err := gt.Version()
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantVersion, version)
		})
	}
}

func Test_Constants(t *testing.T) {
	goldenConstants := getResponse(constants).(Constants)

	type want struct {
		wantErr       bool
		containsErr   string
		wantConstants Constants
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(newConstantsMock().handler(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"could not get network constants",
				Constants{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(newConstantsMock().handler([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal network constants",
				Constants{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(newConstantsMock().handler(readResponse(constants), blankHandler)),
			want{
				false,
				"",
				goldenConstants,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			constants, err := gt.Constants("BLzGD63HA4RP8Fh5xEtvdQSMKa2WzJMZjQPNVUc4Rqy8Lh5BEY1")
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantConstants, constants)
		})
	}
}

func Test_Connections(t *testing.T) {

	var goldenConnections Connections
	json.Unmarshal(readResponse(connections), &goldenConnections)

	type want struct {
		wantErr         bool
		containsErr     string
		wantConnections Connections
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(connectionsHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"could not get network connections",
				Connections{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(connectionsHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal network connections",
				Connections{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(connectionsHandlerMock(readResponse(connections), blankHandler)),
			want{
				false,
				"",
				goldenConnections,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			connections, err := gt.Connections()
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantConnections, connections)
		})
	}
}

func Test_Bootsrap(t *testing.T) {
	var goldenBootstrap Bootstrap
	json.Unmarshal(readResponse(bootstrap), &goldenBootstrap)

	type want struct {
		wantErr       bool
		containsErr   string
		wantBootstrap Bootstrap
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(bootstrapHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"could not get bootstrap",
				Bootstrap{},
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(bootstrapHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could not unmarshal bootstrap",
				Bootstrap{},
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(bootstrapHandlerMock(readResponse(bootstrap), blankHandler)),
			want{
				false,
				"",
				goldenBootstrap,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			bootstrap, err := gt.Bootstrap()
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantBootstrap, bootstrap)
		})
	}
}

func Test_Commit(t *testing.T) {
	var goldenCommit string
	json.Unmarshal(readResponse(commit), &goldenCommit)

	type want struct {
		wantErr     bool
		containsErr string
		wantCommit  string
	}

	cases := []struct {
		name        string
		inputHanler http.Handler
		want
	}{
		{
			"returns rpc error",
			gtGoldenHTTPMock(commitHandlerMock(readResponse(rpcerrors), blankHandler)),
			want{
				true,
				"could not get commit",
				"",
			},
		},
		{
			"fails to unmarshal",
			gtGoldenHTTPMock(commitHandlerMock([]byte(`junk`), blankHandler)),
			want{
				true,
				"could unmarshal commit",
				"",
			},
		},
		{
			"is successful",
			gtGoldenHTTPMock(commitHandlerMock(readResponse(commit), blankHandler)),
			want{
				false,
				"",
				goldenCommit,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.inputHanler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			commit, err := gt.Commit()
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.want.containsErr)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want.wantCommit, commit)
		})
	}
}

func Test_Cycle(t *testing.T) {
	type input struct {
		handler http.Handler
		cycle   int
	}

	type want struct {
		err         bool
		errContains string
		cycle       Cycle
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"failed to get head block",
			input{
				gtGoldenHTTPMock(newBlockMock().handler([]byte(`not_block_data`), blankHandler)),
				10,
			},
			want{
				true,
				"could not get cycle '10': could not get head block",
				Cycle{},
			},
		},
		{
			"failed to get cycle because cycle is in the future",
			input{
				gtGoldenHTTPMock(newBlockMock().handler(readResponse(block), blankHandler)),
				300,
			},
			want{
				true,
				"request is in the future",
				Cycle{},
			},
		},
		{
			"failed to get block less than cycle",
			input{
				gtGoldenHTTPMock(
					newBlockMock().handler(
						readResponse(block),
						newBlockMock().handler(
							[]byte(`not_block_data`),
							blankHandler,
						),
					),
				),
				2,
			},
			want{
				true,
				"could not get block",
				Cycle{},
			},
		},
		{
			"failed to unmarshal cycle",
			input{
				gtGoldenHTTPMock(
					cycleHandlerMock(
						[]byte(`bad_cycle_data`),
						newBlockMock().handler(
							readResponse(block),
							newBlockMock().handler(
								readResponse(block),
								blankHandler,
							),
						),
					),
				),
				2,
			},
			want{
				true,
				"could not unmarshal at cycle hash",
				Cycle{},
			},
		},
		{
			"failed to get cycle block level",
			input{
				gtGoldenHTTPMock(
					cycleHandlerMock(
						readResponse(cycle),
						newBlockMock().handler(
							readResponse(block),
							newBlockMock().handler(
								readResponse(block),
								newBlockMock().handler(
									[]byte(`not_block_data`),
									blankHandler,
								),
							),
						),
					),
				),
				2,
			},
			want{
				true,
				"could not get block",
				Cycle{
					RandomSeed:   "04dca5c197fc2e18309b60844148c55fc7ccdbcb498bd57acd4ac29f16e22846",
					RollSnapshot: 4,
				},
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(mockCycleSuccessful(blankHandler)),
				2,
			},
			want{
				false,
				"",
				Cycle{
					RandomSeed:   "04dca5c197fc2e18309b60844148c55fc7ccdbcb498bd57acd4ac29f16e22846",
					RollSnapshot: 4,
					BlockHash:    "BLfEWKVudXH15N8nwHZehyLNjRuNLoJavJDjSZ7nq8ggfzbZ18p",
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			cycle, err := gt.Cycle(tt.input.cycle)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.cycle, cycle)
		})
	}
}

func Test_ActiveChains(t *testing.T) {
	type input struct {
		handler http.Handler
	}

	type want struct {
		err          bool
		errContains  string
		activeChains ActiveChains
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"returns rpc error",
			input{
				gtGoldenHTTPMock(activeChainsHandlerMock(readResponse(rpcerrors), blankHandler)),
			},
			want{
				true,
				"failed to get active chains",
				nil,
			},
		},
		{
			"fails to unmarshal",
			input{
				gtGoldenHTTPMock(activeChainsHandlerMock([]byte(`junk`), blankHandler)),
			},
			want{
				true,
				"failed to unmarshal active chains",
				nil,
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(activeChainsHandlerMock(readResponse(activechains), blankHandler)),
			},
			want{
				false,
				"",
				ActiveChains{
					{
						ChainID: "NetXdQprcVkpaWU",
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			activeChains, err := gt.ActiveChains()
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.activeChains, activeChains)
		})
	}
}

func mockCycleSuccessful(next http.Handler) http.Handler {
	var blockmock blockHandlerMock
	var oldHTTPBlock blockHandlerMock
	var blockAtLevel blockHandlerMock
	return cycleHandlerMock(
		readResponse(cycle),
		blockmock.handler(
			readResponse(block),
			oldHTTPBlock.handler(
				readResponse(block),
				blockAtLevel.handler(
					readResponse(block),
					next,
				),
			),
		),
	)
}

func mockCycleFailed(next http.Handler) http.Handler {
	var blockmock blockHandlerMock
	var oldHTTPBlock blockHandlerMock
	return cycleHandlerMock(
		[]byte(`bad_cycle_data`),
		blockmock.handler(
			readResponse(block),
			oldHTTPBlock.handler(
				readResponse(block),
				blankHandler,
			),
		),
	)
}
