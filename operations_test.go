package goMXP

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_PreapplyOperation(t *testing.T) {
	type input struct {
		handler                 http.Handler
		preapplyOperationsInput PreapplyOperationsInput
	}

	type want struct {
		err         bool
		errContains string
		operations  []Operations
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles invalid input",
			input{
				gtGoldenHTTPMock(
					preapplyOperationsHandlerMock(
						readResponse(rpcerrors),
						blankHandler,
					),
				),
				PreapplyOperationsInput{},
			},
			want{
				true,
				"invalid input: Key: 'PreapplyOperationsInput.Blockhash'",
				nil,
			},
		},
		{
			"handles rpc error",
			input{
				gtGoldenHTTPMock(
					preapplyOperationsHandlerMock(
						readResponse(rpcerrors),
						blankHandler,
					),
				),
				PreapplyOperationsInput{
					Blockhash: "some_hash",
					Protocol:  "some_protocol",
					Signature: "some_sig",
					Contents:  []Contents{},
				},
			},
			want{
				true,
				"failed to preapply operation",
				nil,
			},
		},
		{
			"handles failure to unmarshal",
			input{
				gtGoldenHTTPMock(
					preapplyOperationsHandlerMock(
						[]byte("junk"),
						blankHandler,
					),
				),
				PreapplyOperationsInput{
					Blockhash: "some_hash",
					Protocol:  "some_protocol",
					Signature: "some_sig",
					Contents:  []Contents{},
				},
			},
			want{
				true,
				"failed to unmarshal operation",
				nil,
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(
					preapplyOperationsHandlerMock(
						readResponse(preapplyOperations),
						blankHandler,
					),
				),
				PreapplyOperationsInput{
					Blockhash: "some_hash",
					Protocol:  "some_protocol",
					Signature: "some_sig",
					Contents:  []Contents{},
				},
			},
			want{
				false,
				"",
				[]Operations{
					{
						Contents: []Contents{
							{
								Kind:         "transaction",
								Source:       "tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH",
								Fee:          NewInt(3000),
								Counter:      NewInt(1263232),
								GasLimit:     NewInt(20000),
								StorageLimit: NewInt(0),
								Amount:       NewInt(50),
								Destination:  "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
								Metadata: &ContentsMetadata{
									BalanceUpdates: []BalanceUpdates{
										{
											Kind:     "contract",
											Contract: "tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH",
											Change:   NewInt(-3000),
										},
										{
											Kind:     "freezer",
											Category: "fees",
											Delegate: "tz1Ke2h7sDdakHJQh8WX4Z372du1KChsksyU",
											Cycle:    229,
											Change:   NewInt(3000),
										},
									},
									OperationResult: &OperationResult{
										Status: "applied",
										BalanceUpdates: []BalanceUpdates{
											{
												Kind:     "contract",
												Contract: "tz1W3HW533csCBLor4NPtU79R2TT2sbKfJDH",
												Change:   NewInt(-50),
											},
											{
												Kind:     "contract",
												Contract: "tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc",
												Change:   NewInt(50),
											},
										},
										ConsumedGas: NewInt(10207),
									},
								},
							},
						},
						Signature: "edsig...."},
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

			operations, err := gt.PreapplyOperations(tt.input.preapplyOperationsInput)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operations, operations)
		})
	}
}

func Test_InjectOperation(t *testing.T) {
	goldenOp := "a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb0300b1a803000008ba0cb2fad622697145cf1665124096d25bc31e00"
	goldenHash := []byte(`"oopfasdfadjkfalksj"`)

	type input struct {
		handler http.Handler
	}

	type want struct {
		err         bool
		errContains string
		result      string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"returns rpc error",
			input{
				gtGoldenHTTPMock(
					injectionOperationHandlerMock(
						readResponse(rpcerrors),
						blankHandler,
					),
				),
			},
			want{
				true,
				"failed to inject operation",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			input{
				gtGoldenHTTPMock(
					injectionOperationHandlerMock(
						[]byte("junk"),
						blankHandler,
					),
				),
			},
			want{
				true,
				"failed to unmarshal operation",
				"",
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(
					injectionOperationHandlerMock(
						goldenHash,
						blankHandler,
					),
				),
			},
			want{
				false,
				"",
				"oopfasdfadjkfalksj",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			result, err := gt.InjectionOperation(InjectionOperationInput{
				Operation: &goldenOp,
			})
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.result, result)
		})
	}
}

func Test_InjectBlock(t *testing.T) {
	goldenRPCError := readResponse(rpcerrors)
	goldenHash := []byte("some_hash")
	type input struct {
		handler http.Handler
	}

	type want struct {
		err         bool
		errContains string
		result      []byte
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"returns rpc error",
			input{
				gtGoldenHTTPMock(
					injectionBlockHandlerMock(
						readResponse(rpcerrors),
						blankHandler,
					),
				),
			},
			want{
				true,
				"failed to inject block",
				goldenRPCError,
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(
					injectionBlockHandlerMock(
						goldenHash,
						blankHandler,
					),
				),
			},
			want{
				false,
				"",
				goldenHash,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			result, err := gt.InjectionBlock(InjectionBlockInput{
				Block: &Block{},
			})
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.result, result)
		})
	}
}

func Test_Counter(t *testing.T) {
	goldenCounter := 10
	goldenRPCError := readResponse(rpcerrors)
	type input struct {
		handler http.Handler
	}

	type want struct {
		err         bool
		errContains string
		counter     int
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"failed to unmarshal counter",
			input{
				gtGoldenHTTPMock(
					counterHandlerMock(
						[]byte(`bad_counter_data`),
						blankHandler,
					),
				),
			},
			want{
				true,
				"failed to unmarshal counter",
				0,
			},
		},
		{
			"returns rpc error",
			input{
				gtGoldenHTTPMock(
					counterHandlerMock(
						goldenRPCError,
						blankHandler,
					),
				),
			},
			want{
				true,
				"failed to get counter",
				0,
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(
					counterHandlerMock(
						readResponse(counter),
						blankHandler,
					),
				),
			},
			want{
				false,
				"",
				goldenCounter,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.handler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			counter, err := gt.Counter(mockBlockHash, mockAddressTz1)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.counter, counter)
		})
	}
}

func Test_ForgeOperation(t *testing.T) {
	type input struct {
		contents []Contents
		branch   string
	}

	type want struct {
		err         bool
		errContains string
		operation   string
	}
	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful transaction",
			input{
				[]Contents{
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      NewInt(10),
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(12345),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         TRANSACTIONOP,
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(34567123),
						Counter:      NewInt(8),
						GasLimit:     NewInt(56787),
						StorageLimit: NewInt(0),
						Amount:       NewInt(54321),
						Destination:  "KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj",
						Kind:         TRANSACTIONOP,
					},
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				false,
				"",
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb0300b1a803018b88e99e66c1c2587f87118449f781cb7d44c9c40000",
			},
		},
		{
			"is successful reveal",
			input{
				[]Contents{
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      NewInt(10),
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
						Kind:         REVEALOP,
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(34567123),
						Counter:      NewInt(8),
						GasLimit:     NewInt(56787),
						StorageLimit: NewInt(0),
						Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
						Kind:         REVEALOP,
					},
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				false,
				"",
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36b0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0000136083897bc97879c53e3e7855838fbbc87303ddd376080fc3d3e136b55d028b6b0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb030000136083897bc97879c53e3e7855838fbbc87303ddd376080fc3d3e136b55d028b",
			},
		},
		{
			"is successful origination",
			input{
				[]Contents{
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      NewInt(10),
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Kind:         ORIGINATIONOP,
						Balance:      NewInt(328763282),
						Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				false,
				"",
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36d0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00928fe29c01ff0008ba0cb2fad622697145cf1665124096d25bc31e000000c602000000c105000764085e036c055f036d0000000325646f046c000000082564656661756c740501035d050202000000950200000012020000000d03210316051f02000000020317072e020000006a0743036a00000313020000001e020000000403190325072c020000000002000000090200000004034f0327020000000b051f02000000020321034c031e03540348020000001e020000000403190325072c020000000002000000090200000004034f0327034f0326034202000000080320053d036d03420000001a0a000000150008ba0cb2fad622697145cf1665124096d25bc31e",
			},
		},
		{
			"is successful delegation",
			input{
				[]Contents{
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      NewInt(10),
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Kind:         DELEGATIONOP,
						Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				false,
				"",
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36e0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00ff0008ba0cb2fad622697145cf1665124096d25bc31e",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			operation, err := ForgeOperation(tt.input.branch, tt.input.contents...)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, operation)
		})
	}
}

func Test_ForgeTransactionOperation(t *testing.T) {
	type input struct {
		contents []ForgeTransactionOperationInput
		branch   string
	}

	type want struct {
		err         bool
		errContains string
		operation   string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				[]ForgeTransactionOperationInput{
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(400000),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(400000),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(400000),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(400000),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(400000),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(400000),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(400000),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(400000),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(400000),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(400000),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(400000),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(400000),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(400000),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(400000),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(400000),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				false,
				"",
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0080b518000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0080b518000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0080b518000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0080b518000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0080b518000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0080b518000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0080b518000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0080b518000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0080b518000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0080b518000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0080b518000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0080b518000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0080b518000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0080b518000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0080b518000008ba0cb2fad622697145cf1665124096d25bc31e00",
			},
		},
		{
			"handles bad branch",
			input{
				[]ForgeTransactionOperationInput{
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(30),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(10000),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
				},
				"junk",
			},
			want{
				true,
				"failed to forge operation: failed to clean branch: failed to decode payload: junk",
				"",
			},
		},
		{
			"handles missing fields",
			input{
				[]ForgeTransactionOperationInput{
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      10,
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(10000),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				true,
				"failed to forge operation: failed to forge transaction: missing amount",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			operation, err := ForgeTransactionOperation(tt.input.branch, tt.input.contents...)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, operation)
		})
	}
}

func Test_ForgeRevealOperation(t *testing.T) {
	type input struct {
		contents ForgeRevealOperationInput
		branch   string
	}

	type want struct {
		err         bool
		errContains string
		operation   string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				ForgeRevealOperationInput{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      10,
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				false,
				"",
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36b0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0000136083897bc97879c53e3e7855838fbbc87303ddd376080fc3d3e136b55d028b",
			},
		},
		{
			"handles bad branch",
			input{
				ForgeRevealOperationInput{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2Ypc",
					Fee:          NewInt(10100),
					Counter:      10,
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
				},
				"junk",
			},
			want{
				true,
				"failed to forge operation: failed to clean branch: failed to decode payload: junk",
				"",
			},
		},
		{
			"handles missing fields",
			input{
				ForgeRevealOperationInput{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2Ypc",
					Fee:          NewInt(10100),
					Counter:      10,
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				true,
				"failed to forge operation: failed to forge reveal operation: missing phk",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			operation, err := ForgeRevealOperation(tt.input.branch, tt.input.contents)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, operation)
		})
	}
}

func Test_ForgeOriginationOperation(t *testing.T) {
	type input struct {
		contents ForgeOriginationOperationInput
		branch   string
	}

	type want struct {
		err         bool
		errContains string
		operation   string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				ForgeOriginationOperationInput{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      10,
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Balance:      NewInt(328763282),
					Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				false,
				"",
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36d0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00928fe29c01ff0008ba0cb2fad622697145cf1665124096d25bc31e000000c602000000c105000764085e036c055f036d0000000325646f046c000000082564656661756c740501035d050202000000950200000012020000000d03210316051f02000000020317072e020000006a0743036a00000313020000001e020000000403190325072c020000000002000000090200000004034f0327020000000b051f02000000020321034c031e03540348020000001e020000000403190325072c020000000002000000090200000004034f0327034f0326034202000000080320053d036d03420000001a0a000000150008ba0cb2fad622697145cf1665124096d25bc31e",
			},
		},
		{
			"handles bad branch",
			input{
				ForgeOriginationOperationInput{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      10,
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Balance:      NewInt(328763282),
					Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
				"junk",
			},
			want{
				true,
				"failed to forge operation: failed to clean branch: failed to decode payload: junk",
				"",
			},
		},
		{
			"handles missing fields",
			input{
				ForgeOriginationOperationInput{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      10,
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				true,
				"failed to forge operation: failed to forge transaction: missing balance",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			operation, err := ForgeOriginationOperation(tt.input.branch, tt.input.contents)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, operation)
		})
	}
}

func Test_ForgeDelegationOperation(t *testing.T) {
	type input struct {
		contents ForgeDelegationOperationInput
		branch   string
	}

	type want struct {
		err         bool
		errContains string
		operation   string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				ForgeDelegationOperationInput{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      10,
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				false,
				"",
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36e0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00ff0008ba0cb2fad622697145cf1665124096d25bc31e",
			},
		},
		{
			"handles bad branch",
			input{
				ForgeDelegationOperationInput{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      10,
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
				"junk",
			},
			want{
				true,
				"failed to forge operation: failed to clean branch: failed to decode payload: junk",
				"",
			},
		},
		{
			"handles missing fields",
			input{
				ForgeDelegationOperationInput{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      10,
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
				},
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			},
			want{
				true,
				"failed to forge operation: failed to forge delegation operation: missing delegate",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			operation, err := ForgeDelegationOperation(tt.input.branch, tt.input.contents)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, operation)
		})
	}
}
func Test_forgeTransactionOperation(t *testing.T) {
	type input struct {
		contents Contents
	}

	type want struct {
		err         bool
		errContains string
		operation   string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"works with tz1 addresses",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Amount:       NewInt(30),
					Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Kind:         TRANSACTIONOP,
				},
			},
			want{
				false,
				"",
				"6c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e001e000008ba0cb2fad622697145cf1665124096d25bc31e00",
			},
		},
		{
			"works with tz1 to kt",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Amount:       NewInt(30),
					Destination:  "KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj",
					Kind:         TRANSACTIONOP,
				},
			},
			want{
				false,
				"",
				"6c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e001e018b88e99e66c1c2587f87118449f781cb7d44c9c40000",
			},
		},
		{
			"handles common forge error",
			input{
				Contents{
					Source:       "LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Amount:       NewInt(30),
					Destination:  "KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj",
					Kind:         TRANSACTIONOP,
				},
			},
			want{
				true,
				"failed to remove tz1 from source prefix",
				"",
			},
		},
		{
			"handles failed to remove kt prefix from destination",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Amount:       NewInt(30),
					Destination:  "KTJUNK",
					Kind:         TRANSACTIONOP,
				},
			},
			want{
				true,
				"failed to forge transaction: provided destination is not a valid KT1 address: failed to decode payload: KTJUNK",
				"",
			},
		},
		{
			"handles failed to remove tz1 prefix from destination",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Amount:       NewInt(30),
					Destination:  "tz1JUNK",
					Kind:         TRANSACTIONOP,
				},
			},
			want{
				true,
				"failed to forge transaction: provided destination is not a valid tz1 address: failed to decode payload: tz1JUNK",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			operation, err := forgeTransactionOperation(tt.input.contents)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, operation)
		})
	}
}

func Test_forgeRevealOperation(t *testing.T) {
	type input struct {
		contents Contents
	}

	type want struct {
		err         bool
		errContains string
		operation   string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
					Kind:         REVEALOP,
				},
			},
			want{
				false,
				"",
				"6b0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0000136083897bc97879c53e3e7855838fbbc87303ddd376080fc3d3e136b55d028b",
			},
		},
		{
			"handles failure to forge common",
			input{
				Contents{
					Source:       "tAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
					Kind:         REVEALOP,
				},
			},
			want{
				true,
				"failed to forge reveal operation: failed to remove tz1 from source prefix",
				"",
			},
		},
		{
			"handles failure to clean pub key",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Phk:          "tnktxAzm32--9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
					Kind:         REVEALOP,
				},
			},
			want{
				true,
				"failed to forge reveal operation: failed to decode payload",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			operation, err := forgeRevealOperation(tt.input.contents)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, operation)
		})
	}
}

func Test_forgeOriginationOperation(t *testing.T) {
	type input struct {
		contents Contents
	}

	type want struct {
		err         bool
		errContains string
		operation   string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Kind:         ORIGINATIONOP,
					Balance:      NewInt(328763282),
					Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
			},
			want{
				false,
				"",
				"6d0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00928fe29c01ff0008ba0cb2fad622697145cf1665124096d25bc31e000000c602000000c105000764085e036c055f036d0000000325646f046c000000082564656661756c740501035d050202000000950200000012020000000d03210316051f02000000020317072e020000006a0743036a00000313020000001e020000000403190325072c020000000002000000090200000004034f0327020000000b051f02000000020321034c031e03540348020000001e020000000403190325072c020000000002000000090200000004034f0327034f0326034202000000080320053d036d03420000001a0a000000150008ba0cb2fad622697145cf1665124096d25bc31e",
			},
		},
		{
			"handles failure to forge common",
			input{
				Contents{
					Source:       "tAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Kind:         ORIGINATIONOP,
					Balance:      NewInt(328763282),
					Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
			},
			want{
				true,
				"failed to forge origination operation: failed to remove tz1 from source prefix",
				"",
			},
		},
		{
			"handles failure to clean delegate",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Kind:         ORIGINATIONOP,
					Balance:      NewInt(328763282),
					Delegate:     "tz1LSAy890--cAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
			},
			want{
				true,
				"failed to forge origination operation: failed to decode payload",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			operation, err := forgeOriginationOperation(tt.input.contents)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, operation)
		})
	}
}

func Test_forgeDelegationOperation(t *testing.T) {
	type input struct {
		contents Contents
	}

	type want struct {
		err         bool
		errContains string
		operation   string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Kind:         DELEGATIONOP,
					Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
			},
			want{
				false,
				"",
				"6e0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00ff0008ba0cb2fad622697145cf1665124096d25bc31e",
			},
		},
		{
			"handles failure to forge common",
			input{
				Contents{
					Source:       "tAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Kind:         DELEGATIONOP,
					Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
			},
			want{
				true,
				"failed to forge delegation operation: failed to remove tz1 from source prefix",
				"",
			},
		},
		{
			"handles failure to clean delegate tz1",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Kind:         DELEGATIONOP,
					Balance:      NewInt(328763282),
					Delegate:     "tz1LSAy890--cAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
			},
			want{
				true,
				"failed to forge delegation operation: failed to decode payload",
				"",
			},
		},
		{
			"handles failure to clean delegate KT1",
			input{
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Kind:         DELEGATIONOP,
					Balance:      NewInt(328763282),
					Delegate:     "KT1LSAy890--cAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
			},
			want{
				true,
				"failed to forge delegation operation: failed to decode payload",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			operation, err := forgeDelegationOperation(tt.input.contents)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, operation)
		})
	}
}

func Test_UnforgeOperation(t *testing.T) {
	mockHash := "BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk"
	type input struct {
		handler   http.Handler
		hexString string
		signed    bool
	}

	type want struct {
		err         bool
		errContains string
		contents    *[]Contents
		branch      *string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful transaction",
			input{
				gtGoldenHTTPMock(blankHandler),
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36c0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e006c0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb0300b1a803018b88e99e66c1c2587f87118449f781cb7d44c9c40000",
				false,
			},
			want{
				false,
				"",
				&[]Contents{
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      NewInt(10),
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Amount:       NewInt(12345),
						Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Kind:         TRANSACTIONOP,
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(34567123),
						Counter:      NewInt(8),
						GasLimit:     NewInt(56787),
						StorageLimit: NewInt(0),
						Amount:       NewInt(54321),
						Destination:  "KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj",
						Kind:         TRANSACTIONOP,
					},
				},
				&mockHash,
			},
		},
		{
			"is successful reveal",
			input{
				gtGoldenHTTPMock(blankHandler),
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36b0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0000136083897bc97879c53e3e7855838fbbc87303ddd376080fc3d3e136b55d028b6b0008ba0cb2fad622697145cf1665124096d25bc31ed3e7bd1008d3bb030000136083897bc97879c53e3e7855838fbbc87303ddd376080fc3d3e136b55d028ba",
				false,
			},
			want{
				false,
				"",
				&[]Contents{
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      NewInt(10),
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
						Kind:         REVEALOP,
					},
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(34567123),
						Counter:      NewInt(8),
						GasLimit:     NewInt(56787),
						StorageLimit: NewInt(0),
						Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
						Kind:         REVEALOP,
					},
				},
				&mockHash,
			},
		},
		{
			"is successful origination",
			input{
				gtGoldenHTTPMock(blankHandler),
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36d0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00928fe29c01ff0008ba0cb2fad622697145cf1665124096d25bc31e000000c602000000c105000764085e036c055f036d0000000325646f046c000000082564656661756c740501035d050202000000950200000012020000000d03210316051f02000000020317072e020000006a0743036a00000313020000001e020000000403190325072c020000000002000000090200000004034f0327020000000b051f02000000020321034c031e03540348020000001e020000000403190325072c020000000002000000090200000004034f0327034f0326034202000000080320053d036d03420000001a0a000000150008ba0cb2fad622697145cf1665124096d25bc31e",
				false,
			},
			want{
				false,
				"",
				&[]Contents{
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      NewInt(10),
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Kind:         ORIGINATIONOP,
						Balance:      NewInt(328763282),
						Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
				},
				&mockHash,
			},
		},
		{
			"is successful delegation",
			input{
				gtGoldenHTTPMock(blankHandler),
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36e0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00ff0008ba0cb2fad622697145cf1665124096d25bc31e",
				false,
			},
			want{
				false,
				"",
				&[]Contents{
					{
						Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
						Fee:          NewInt(10100),
						Counter:      NewInt(10),
						GasLimit:     NewInt(10100),
						StorageLimit: NewInt(0),
						Kind:         DELEGATIONOP,
						Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					},
				},
				&mockHash,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			branch, contents, err := UnforgeOperation(tt.input.hexString, tt.input.signed)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.branch, branch)
			assert.Equal(t, tt.want.contents, contents)
		})
	}
}

func Test_unforgeTransactionOperation(t *testing.T) {
	type input struct {
		operation string
	}

	type want struct {
		err         bool
		errContains string
		contents    Contents
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"works with tz1 addresses",
			input{
				"0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e001e000008ba0cb2fad622697145cf1665124096d25bc31e00",
			},
			want{
				false,
				"",
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Amount:       NewInt(30),
					Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Kind:         TRANSACTIONOP,
				},
			},
		},
		{
			"works with tz1 to kt",
			input{
				"0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e001e018b88e99e66c1c2587f87118449f781cb7d44c9c40000",
			},
			want{
				false,
				"",
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Amount:       NewInt(30),
					Destination:  "KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj",
					Kind:         TRANSACTIONOP,
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			contents, _, err := unforgeTransactionOperation(tt.input.operation)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.contents, contents)
		})
	}
}

func Test_unforgeRevealOperation(t *testing.T) {
	type input struct {
		operation string
	}

	type want struct {
		err         bool
		errContains string
		contents    Contents
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				"0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e0000136083897bc97879c53e3e7855838fbbc87303ddd376080fc3d3e136b55d028b",
			},
			want{
				false,
				"",
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
					Kind:         REVEALOP,
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			contents, _, err := unforgeRevealOperation(tt.input.operation)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.contents, contents)
		})
	}
}

func Test_unforgeOriginationOperation(t *testing.T) {
	type input struct {
		operation string
	}

	type want struct {
		err         bool
		errContains string
		contents    Contents
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				"0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00928fe29c01ff0008ba0cb2fad622697145cf1665124096d25bc31e000000c602000000c105000764085e036c055f036d0000000325646f046c000000082564656661756c740501035d050202000000950200000012020000000d03210316051f02000000020317072e020000006a0743036a00000313020000001e020000000403190325072c020000000002000000090200000004034f0327020000000b051f02000000020321034c031e03540348020000001e020000000403190325072c020000000002000000090200000004034f0327034f0326034202000000080320053d036d03420000001a0a000000150008ba0cb2fad622697145cf1665124096d25bc31e",
			},
			want{
				false,
				"",
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Kind:         ORIGINATIONOP,
					Balance:      NewInt(328763282),
					Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			contents, _, err := unforgeOriginationOperation(tt.input.operation)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.contents, contents)
		})
	}
}

func Test_unforgeDelegationOperation(t *testing.T) {
	type input struct {
		operation string
	}

	type want struct {
		err         bool
		errContains string
		contents    Contents
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				"0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00ff0008ba0cb2fad622697145cf1665124096d25bc31e",
			},
			want{
				false,
				"",
				Contents{
					Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
					Fee:          NewInt(10100),
					Counter:      NewInt(10),
					GasLimit:     NewInt(10100),
					StorageLimit: NewInt(0),
					Kind:         DELEGATIONOP,
					Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			contents, _, err := unforgeDelegationOperation(tt.input.operation)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.contents, contents)
		})
	}
}

func Test_StripBranchFromForgedOperation(t *testing.T) {
	op := "a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add36d0008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00928fe29c01ff0008ba0cb2fad622697145cf1665124096d25bc31e000000c602000000c105000764085e036c055f036d0000000325646f046c000000082564656661756c740501035d050202000000950200000012020000000d03210316051f02000000020317072e020000006a0743036a00000313020000001e020000000403190325072c020000000002000000090200000004034f0327020000000b051f02000000020321034c031e03540348020000001e020000000403190325072c020000000002000000090200000004034f0327034f0326034202000000080320053d036d03420000001a0a000000150008ba0cb2fad622697145cf1665124096d25bc31e"
	branch, _, err := StripBranchFromForgedOperation(op, false)
	assert.Nil(t, err)
	assert.Equal(t, "BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk", branch)
}

func Test_checkBoolean(t *testing.T) {
	type input struct {
		hexString string
	}

	type want struct {
		err         bool
		errContains string
		res         bool
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is boolean",
			input{
				hexString: "ff",
			},
			want{
				false,
				"",
				true,
			},
		},
		{
			"is not boolean",
			input{
				hexString: "00",
			},
			want{
				false,
				"",
				false,
			},
		},
		{
			"is unkown",
			input{
				hexString: "dssdf",
			},
			want{
				true,
				"boolean value is invalid",
				false,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := checkBoolean(tt.input.hexString)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.res, res)
		})
	}
}

func Test_parseAddress(t *testing.T) {
	type input struct {
		hexString string
	}

	type want struct {
		err         bool
		errContains string
		res         string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful tz1",
			input{
				hexString: "000008ba0cb2fad622697145cf1665124096d25bc31e",
			},
			want{
				false,
				"",
				"tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
			},
		},
		{
			"is successful KT1",
			input{
				hexString: "018b88e99e66c1c2587f87118449f781cb7d44c9c400",
			},
			want{
				false,
				"",
				"KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj",
			},
		},
		{
			"handles junk",
			input{
				hexString: "e66c1c2587f87118449f781cb7d44c9c400",
			},
			want{
				true,
				"address format not supported",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := parseAddress(tt.input.hexString)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.res, res)
		})
	}
}

func Test_removeHexPrefix(t *testing.T) {
	type input struct {
		payload string
		prefix  prefix
	}

	type want struct {
		err         bool
		errContains string
		res         string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful tz1",
			input{
				"tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				tz1prefix,
			},
			want{
				false,
				"",
				"08ba0cb2fad622697145cf1665124096d25bc31e",
			},
		},
		{
			"is successful KT1",
			input{
				"KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj",
				ktprefix,
			},
			want{
				false,
				"",
				"8b88e99e66c1c2587f87118449f781cb7d44c9c4",
			},
		},
		{
			"is successful KT1",
			input{
				"KT1MJZWHKZU7ViybRLsphP3ppiiTc7myP2aj",
				ktprefix,
			},
			want{
				false,
				"",
				"8b88e99e66c1c2587f87118449f781cb7d44c9c4",
			},
		},
		{
			"is successful branch",
			input{
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
				branchprefix,
			},
			want{
				false,
				"",
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add3",
			},
		},
		{
			"handles payload not matching prefix",
			input{
				"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
				edpkprefix,
			},
			want{
				true,
				"payload did not match prefix",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := removeHexPrefix(tt.input.payload, tt.input.prefix)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.res, res)
		})
	}
}

func Test_bigNumberToZarith(t *testing.T) {
	type input struct {
		num *Int
	}

	type want struct {
		res string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful positive number",
			input{
				NewInt(302393),
			},
			want{
				"b9ba12",
			},
		},
		{
			"is successful negative number",
			input{
				NewInt(-302393),
			},
			want{
				"b9ba00",
			},
		},
		{
			"is successful zero",
			input{
				NewInt(0),
			},
			want{
				"00",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			res := bigNumberToZarith(*tt.input.num)
			assert.Equal(t, tt.want.res, res)
		})
	}
}

func Test_splitAndReturnRest(t *testing.T) {
	type input struct {
		payload string
		length  int
	}

	type want struct {
		first  string
		second string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				"08ba0cb2fad622697145cf1665124096d25bc31e",
				15,
			},
			want{
				"08ba0cb2fad6226",
				"97145cf1665124096d25bc31e",
			},
		},
		{
			"is successful when payload is too short",
			input{
				"08ba0cb2fad622697145cf1665124096d25bc31e",
				300,
			},
			want{
				"08ba0cb2fad622697145cf1665124096d25bc31e",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			first, second := splitAndReturnRest(tt.input.payload, tt.input.length)
			assert.Equal(t, tt.want.first, first)
			assert.Equal(t, tt.want.second, second)
		})
	}
}

func Test_prefixAndBase58Encode(t *testing.T) {
	type input struct {
		payload string
		prefix  prefix
	}

	type want struct {
		err         bool
		errContains string
		res         string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				"08ba0cb2fad622697145cf1665124096d25bc31e",
				tz1prefix,
			},
			want{
				false,
				"",
				"tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
			},
		},
		{
			"handles failed encode",
			input{
				"08ba0cb----***20()2fad622697145cf1665124096d25bc31e",
				tz1prefix,
			},
			want{
				true,
				"failed to encode to base58",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := prefixAndBase58Encode(tt.input.payload, tt.input.prefix)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.res, res)
		})
	}
}

func Test_zarithToBigNumber(t *testing.T) {
	type input struct {
		hexString string
	}

	type want struct {
		err         bool
		errContains string
		res         *Int
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful positive number",
			input{
				"b9ba12",
			},
			want{
				false,
				"",
				NewInt(302393),
			},
		},
		{
			"is successful negative number",
			input{
				"b9ba00",
			},
			want{
				false,
				"",
				NewInt(7481),
			},
		},
		{
			"is successful zero",
			input{
				"00",
			},
			want{
				false,
				"",
				NewInt(0),
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := zarithToBigNumber(tt.input.hexString)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.res, res)
		})
	}
}

func Test_findZarithEndIndex(t *testing.T) {
	type input struct {
		hexString string
	}

	type want struct {
		err         bool
		errContains string
		res         int
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				"08ba0cb2fad622697145cf1665124096d25bc31e",
			},
			want{
				false,
				"",
				2,
			},
		},
		{
			"handles failed to find Zarith end index",
			input{
				"^^^^^^---()*97145cf1665124096d25bc31e",
			},
			want{
				true,
				"failed to find Zarith end index",
				0,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := findZarithEndIndex(tt.input.hexString)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.res, res)
		})
	}
}

func Test_parsePublicKey(t *testing.T) {
	type input struct {
		hexString string
	}

	type want struct {
		err         bool
		errContains string
		res         string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				"00136083897bc97879c53e3e7855838fbbc87303ddd376080fc3d3e136b55d028b",
			},
			want{
				false,
				"",
				"edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
			},
		},
		{
			"handles public key format not supported",
			input{
				"136083897bc97879c53e3e7855838fbbc87303ddd376080fc3d3e136b55d028b",
			},
			want{
				true,
				"public key format not supported",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := parsePublicKey(tt.input.hexString)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.res, res)
		})
	}
}

func Test_parseTzAddress(t *testing.T) {
	type input struct {
		hexString string
	}

	type want struct {
		err         bool
		errContains string
		res         string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				"0008ba0cb2fad622697145cf1665124096d25bc31e",
			},
			want{
				false,
				"",
				"tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
			},
		},
		{
			"handles address format not supported",
			input{
				"136083897bc97879c53e3e7855838fbbc87303ddd376080fc3d3e136b55d028b",
			},
			want{
				true,
				"address format not supported",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := parseTzAddress(tt.input.hexString)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.res, res)
		})
	}
}

func Test_cleanBranch(t *testing.T) {
	type want struct {
		err         bool
		errContains string
		branch      string
	}

	cases := []struct {
		name  string
		input string
		want  want
	}{
		{
			"is successful",
			"BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk",
			want{
				false,
				"",
				"a732d3520eeaa3de98d78e5e5cb6c85f72204fd46feb9f76853841d4a701add3",
			},
		},
		{
			"handles error",
			"junk",
			want{
				true,
				"failed to clean branch: failed to decode payload: junk",
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			branch, err := cleanBranch(tt.input)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.branch, branch)
		})
	}

}

func Test_validateTransaction(t *testing.T) {
	type want struct {
		err         bool
		errContains string
	}

	cases := []struct {
		name  string
		input Contents
		want  want
	}{
		{
			"is successful",
			Contents{
				Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				Fee:          NewInt(10100),
				Counter:      NewInt(10),
				GasLimit:     NewInt(10100),
				StorageLimit: NewInt(0),
				Amount:       NewInt(10000),
				Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2Ypc",
				Kind:         TRANSACTIONOP,
			},
			want{
				false,
				"",
			},
		},
		{
			"handles invalid",
			Contents{
				Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				Counter:      NewInt(10),
				StorageLimit: NewInt(0),
				Kind:         REVEALOP,
			},
			want{
				true,
				"wrong kind for transaction: missing amount: missing destination: missing fee: missing gas limit",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTransaction(tt.input)
			checkErr(t, tt.want.err, tt.want.errContains, err)
		})
	}
}
func Test_validateOrigination(t *testing.T) {
	type want struct {
		err         bool
		errContains string
	}

	cases := []struct {
		name  string
		input Contents
		want  want
	}{
		{
			"is successful",
			Contents{
				Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				Fee:          NewInt(10100),
				Counter:      NewInt(10),
				GasLimit:     NewInt(10100),
				StorageLimit: NewInt(0),
				Balance:      NewInt(10000),
				Kind:         ORIGINATIONOP,
			},
			want{
				false,
				"",
			},
		},
		{
			"handles invalid",
			Contents{
				Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				Counter:      NewInt(10),
				StorageLimit: NewInt(0),
				Kind:         REVEALOP,
			},
			want{
				true,
				"wrong kind for origination: missing balance: missing fee: missing gas limit",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOrigination(tt.input)
			checkErr(t, tt.want.err, tt.want.errContains, err)
		})
	}
}
func Test_validateDelegation(t *testing.T) {
	type want struct {
		err         bool
		errContains string
	}

	cases := []struct {
		name  string
		input Contents
		want  want
	}{
		{
			"is successful",
			Contents{
				Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				Fee:          NewInt(10100),
				Counter:      NewInt(10),
				GasLimit:     NewInt(10100),
				StorageLimit: NewInt(0),
				Delegate:     "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				Kind:         DELEGATIONOP,
			},
			want{
				false,
				"",
			},
		},
		{
			"handles invalid",
			Contents{
				Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				Counter:      NewInt(10),
				StorageLimit: NewInt(0),
				Kind:         REVEALOP,
			},
			want{
				true,
				"wrong kind for delegation: missing delegate: missing fee: missing gas limit",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDelegation(tt.input)
			checkErr(t, tt.want.err, tt.want.errContains, err)
		})
	}
}

func Test_validateReveal(t *testing.T) {
	type want struct {
		err         bool
		errContains string
	}

	cases := []struct {
		name  string
		input Contents
		want  want
	}{
		{
			"is successful",
			Contents{
				Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				Fee:          NewInt(10100),
				Counter:      NewInt(10),
				GasLimit:     NewInt(10100),
				StorageLimit: NewInt(0),
				Phk:          "edpktnktxAzmXPD9XVNqAvdCFb76vxzQtkbVkSEtXcTz33QZQdb4JQ",
				Kind:         REVEALOP,
			},
			want{
				false,
				"",
			},
		},
		{
			"handles invalid",
			Contents{
				Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				Counter:      NewInt(10),
				StorageLimit: NewInt(0),
				Kind:         DELEGATIONOP,
			},
			want{
				true,
				"wrong kind for reveal: missing phk: missing fee: missing gas limit",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := validateReveal(tt.input)
			checkErr(t, tt.want.err, tt.want.errContains, err)
		})
	}
}

func Test_validateCommon(t *testing.T) {
	type want struct {
		err         bool
		errContains string
	}

	cases := []struct {
		name  string
		input Contents
		want  want
	}{
		{
			"is successful",
			Contents{
				Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				Fee:          NewInt(10100),
				Counter:      NewInt(10),
				GasLimit:     NewInt(10100),
				StorageLimit: NewInt(0),
				Kind:         DELEGATIONOP,
			},
			want{
				false,
				"",
			},
		},
		{
			"handles invalid",
			Contents{
				Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
				Counter:      NewInt(10),
				StorageLimit: NewInt(0),
				Kind:         DELEGATIONOP,
			},
			want{
				true,
				"missing fee: missing gas limit",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCommon(tt.input)
			checkErr(t, tt.want.err, tt.want.errContains, err)
		})
	}
}

func Test_shrinkMultiError(t *testing.T) {
	type want struct {
		err         bool
		errContains string
	}

	cases := []struct {
		name  string
		input []error
		want  want
	}{
		{
			"is successful",
			[]error{
				errors.New("some error"),
				errors.New("another error"),
			},
			want{
				true,
				"some error: another error",
			},
		},
		{
			"handles empty",
			[]error{},
			want{
				false,
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := shrinkMultiError(tt.input)
			checkErr(t, tt.want.err, tt.want.errContains, err)
		})
	}
}

func Test_UnForgeOperationWithRPC(t *testing.T) {
	type input struct {
		inputHandler http.Handler
		operation    UnforgeOperationWithRPCInput
	}

	type want struct {
		err         bool
		errContains string
		operations  []Operations
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles invalid input",
			input{
				gtGoldenHTTPMock(unforgeOperationWithRPCMock(readResponse(rpcerrors), blankHandler)),
				UnforgeOperationWithRPCInput{},
			},
			want{
				true,
				"invalid input: Key: 'UnforgeOperationWithRPCInput.Operations'",
				[]Operations{},
			},
		},
		{
			"handles rpc error",
			input{
				gtGoldenHTTPMock(unforgeOperationWithRPCMock(readResponse(rpcerrors), blankHandler)),
				UnforgeOperationWithRPCInput{
					Operations: []UnforgeOperationWithRPCOperation{
						{
							Data:   "some_data",
							Branch: "some_branch",
						},
					},
				},
			},
			want{
				true,
				"failed to unforge forge operations with RPC",
				[]Operations{},
			},
		},
		{
			"handles failure to unmarshal",
			input{
				gtGoldenHTTPMock(unforgeOperationWithRPCMock([]byte(`junk`), blankHandler)),
				UnforgeOperationWithRPCInput{
					Operations: []UnforgeOperationWithRPCOperation{
						{
							Data:   "some_data",
							Branch: "some_branch",
						},
					},
				},
			},
			want{
				true,
				"failed to unforge forge operations with RPC: invalid character",
				[]Operations{},
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(unforgeOperationWithRPCMock(readResponse(parseOperations), blankHandler)),
				UnforgeOperationWithRPCInput{
					Operations: []UnforgeOperationWithRPCOperation{
						{
							Data:   "some_data",
							Branch: "some_branch",
						},
					},
				},
			},
			want{
				false,
				"",
				[]Operations{
					{
						Protocol: "",
						ChainID:  "",
						Hash:     "",
						Branch:   "BLz6yCE4BUL4ppo1zsEWdK9FRCt15WAY7ECQcuK9RtWg4xeEVL7",
						Contents: []Contents{
							{
								Kind:             "transaction",
								Source:           "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
								Fee:              NewInt(10100),
								Counter:          NewInt(10),
								GasLimit:         NewInt(10100),
								StorageLimit:     NewInt(0),
								Amount:           NewInt(12345),
								Destination:      "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
								Delegate:         "",
								Phk:              "",
								Secret:           "",
								Level:            0,
								ManagerPublicKey: "",
								Balance:          nil,
								Period:           0,
								Proposal:         "",
								Proposals:        []string(nil),
								Ballot:           "",
								Metadata:         nil,
							},
						},
						Signature: "edsigtXomBKi5CTRf5cjATJWSyaRvhfYNHqSUGrn4SdbYRcGwQrUGjzEfQDTuqHhuA8b2d8NarZjz8TRf65WkpQmo423BtomS8Q",
					},
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.inputHandler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			op, err := gt.UnforgeOperationWithRPC("BLyvCRkxuTXkx1KeGvrcEXiPYj4p1tFxzvFDhoHE7SFKtmP1rbk", tt.input.operation)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operations, op)
		})
	}
}

func Test_ForgeOperationWithRPC(t *testing.T) {
	type input struct {
		inputHandler               http.Handler
		forgeOperationWithRPCInput ForgeOperationWithRPCInput
	}

	type want struct {
		err         bool
		errContains string
		operation   string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles invalid input",
			input{
				gtGoldenHTTPMock(forgeOperationWithRPCMock(readResponse(rpcerrors), blankHandler)),
				ForgeOperationWithRPCInput{},
			},
			want{
				true,
				"invalid input: Key: 'ForgeOperationWithRPCInput.Blockhash'",
				"",
			},
		},
		{
			"handles rpc error",
			input{
				gtGoldenHTTPMock(forgeOperationWithRPCMock(readResponse(rpcerrors), blankHandler)),
				ForgeOperationWithRPCInput{
					Blockhash: "some_hash",
					Branch:    "some_branch",
					Contents:  []Contents{},
				},
			},
			want{
				true,
				"failed to forge operation: rpc error (somekind)",
				"",
			},
		},
		{
			"handles failure to unmarshal",
			input{
				gtGoldenHTTPMock(forgeOperationWithRPCMock([]byte(`junk`), blankHandler)),
				ForgeOperationWithRPCInput{
					Blockhash: "some_hash",
					Branch:    "some_branch",
					Contents:  []Contents{},
				},
			},
			want{
				true,
				"failed to forge operation: invalid character",
				"",
			},
		},
		{
			"handles failure to strip operation branch",
			input{
				gtGoldenHTTPMock(forgeOperationWithRPCMock([]byte(`"some_junk_op_string"`), unforgeOperationWithRPCMock(readResponse(rpcerrors), blankHandler))),
				ForgeOperationWithRPCInput{
					Blockhash: "some_hash",
					Branch:    "some_branch",
					Contents:  []Contents{},
				},
			},
			want{
				true,
				"failed to forge operation: unable to verify rpc returned a valid contents",
				"some_junk_op_string",
			},
		},
		{
			"handles failure to parse forged operation",
			input{
				gtGoldenHTTPMock(forgeOperationWithRPCMock([]byte(`"some_operation_string"`), unforgeOperationWithRPCMock(readResponse(rpcerrors), blankHandler))),
				ForgeOperationWithRPCInput{
					Blockhash: "some_hash",
					Branch:    "some_branch",
					Contents:  []Contents{},
				},
			},
			want{
				true,
				"failed to forge operation: unable to verify rpc returned a valid contents",
				"some_operation_string",
			},
		},
		{
			"handles failure to match forge with expected contents",
			input{
				gtGoldenHTTPMock(forgeOperationWithRPCMock([]byte(`"a79ec80dba1f8ddb2cde90b8f12f7c62fdc36556030281ff8904a3d0df82cddc08000008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e00"`), unforgeOperationWithRPCMock(readResponse(parseOperations), blankHandler))),
				ForgeOperationWithRPCInput{
					Blockhash: "some_hash",
					Branch:    "some_branch",
					Contents: []Contents{
						{
							Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
							Fee:          NewInt(100),
							Counter:      NewInt(10),
							GasLimit:     NewInt(10100),
							StorageLimit: NewInt(0),
							Amount:       NewInt(12345),
							Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
							Kind:         TRANSACTIONOP,
						},
					},
				},
			},
			want{
				true,
				"failed to forge operation: alert rpc returned invalid contents",
				"a79ec80dba1f8ddb2cde90b8f12f7c62fdc36556030281ff8904a3d0df82cddc08000008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e00",
			},
		},
		{
			"is successful",
			input{
				gtGoldenHTTPMock(forgeOperationWithRPCMock([]byte(`"a79ec80dba1f8ddb2cde90b8f12f7c62fdc36556030281ff8904a3d0df82cddc08000008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e00"`), unforgeOperationWithRPCMock(readResponse(parseOperations), blankHandler))),
				ForgeOperationWithRPCInput{
					Blockhash: "some_hash",
					Branch:    "some_branch",
					Contents: []Contents{
						{
							Source:       "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
							Fee:          NewInt(10100),
							Counter:      NewInt(10),
							GasLimit:     NewInt(10100),
							StorageLimit: NewInt(0),
							Amount:       NewInt(12345),
							Destination:  "tz1LSAycAVcNdYnXCy18bwVksXci8gUC2YpA",
							Kind:         TRANSACTIONOP,
						},
					},
				},
			},
			want{
				false,
				"",
				"a79ec80dba1f8ddb2cde90b8f12f7c62fdc36556030281ff8904a3d0df82cddc08000008ba0cb2fad622697145cf1665124096d25bc31ef44e0af44e00b960000008ba0cb2fad622697145cf1665124096d25bc31e00",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.input.inputHandler)
			defer server.Close()

			gt, err := New(server.URL)
			assert.Nil(t, err)

			op, err := gt.ForgeOperationWithRPC(tt.input.forgeOperationWithRPCInput)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.operation, op)
		})
	}
}
