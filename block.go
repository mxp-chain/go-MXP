package goMXP

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

/*
Int Wrapper
Description: Int wraps go's big.Int.
*/
type Int struct {
	Big *big.Int
}

// NewInt returns a pointer GoMXP's wrapper Int
func NewInt(i int) *Int {
	return &Int{Big: big.NewInt(int64(i))}
}

func newInt(bigintstring []byte) (*Int, error) {
	i := &Int{}
	err := i.UnmarshalJSON(bigintstring)
	return i, err
}

/*
UnmarshalJSON implements the json.Marshaler interface for BigInt

Parameters:

	b:
		The byte representation of a BigInt.
*/
func (i *Int) UnmarshalJSON(b []byte) error {
	var val string
	err := json.Unmarshal(b, &val)
	if err != nil {
		return err
	}
	i.Big = big.NewInt(0)
	i.Big.SetString(val, 10)

	return nil
}

/*
MarshalJSON implements the json.Marshaler interface for BigInt
*/
func (i *Int) MarshalJSON() ([]byte, error) {
	val, err := i.Big.MarshalText()
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("\"%s\"", val)), nil
}

/*
Block represents a MXP block.

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://MXP.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Block struct {
	Protocol   string         `json:"protocol"`
	ChainID    string         `json:"chain_id"`
	Hash       string         `json:"hash"`
	Header     Header         `json:"header"`
	Metadata   Metadata       `json:"metadata"`
	Operations [][]Operations `json:"operations"`
}

/*
Header represents the header in a MXP block

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://MXP.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Header struct {
	Level            int       `json:"level"`
	Proto            int       `json:"proto"`
	Predecessor      string    `json:"Predecessor"`
	Timestamp        time.Time `json:"timestamp"`
	ValidationPass   int       `json:"validation_pass"`
	OperationsHash   string    `json:"operations_hash"`
	Fitness          []string  `json:"fitness"`
	Context          string    `json:"context"`
	Priority         int       `json:"priority"`
	ProofOfWorkNonce string    `json:"proof_of_work_nonce"`
	Signature        string    `json:"signature"`
}

/*
Metadata represents the metadata in a MXP block

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://MXP.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Metadata struct {
	Protocol               string                   `json:"protocol"`
	NextProtocol           string                   `json:"next_protocol"`
	TestChainStatus        TestChainStatus          `json:"test_chain_status"`
	MaxOperationsTTL       int                      `json:"max_operations_ttl"`
	MaxOperationDataLength int                      `json:"max_operation_data_length"`
	MaxBlockHeaderLength   int                      `json:"max_block_header_length"`
	MaxOperationListLength []MaxOperationListLength `json:"max_operation_list_length"`
	Baker                  string                   `json:"baker"`
	Level                  Level                    `json:"level"`
	VotingPeriodKind       string                   `json:"voting_period_kind"`
	NonceHash              interface{}              `json:"nonce_hash"`
	ConsumedGas            string                   `json:"consumed_gas"`
	Deactivated            []string                 `json:"deactivated"`
	BalanceUpdates         []BalanceUpdates         `json:"balance_updates"`
}

/*
TestChainStatus represents the testchainstatus in a MXP block

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://MXP.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type TestChainStatus struct {
	Status string `json:"status"`
}

/*
MaxOperationListLength represents the maxoperationlistlength in a MXP block

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://MXP.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type MaxOperationListLength struct {
	MaxSize int `json:"max_size"`
	MaxOp   int `json:"max_op,omitempty"`
}

/*
Level represents the level in a MXP block

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://MXP.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Level struct {
	Level                int  `json:"level"`
	LevelPosition        int  `json:"level_position"`
	Cycle                int  `json:"cycle"`
	CyclePosition        int  `json:"cycle_position"`
	VotingPeriod         int  `json:"voting_period"`
	VotingPeriodPosition int  `json:"voting_period_position"`
	ExpectedCommitment   bool `json:"expected_commitment"`
}

/*
BalanceUpdates represents the balance updates in a MXP block

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://MXP.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type BalanceUpdates struct {
	Kind     string `json:"kind"`
	Contract string `json:"contract,omitempty"`
	Change   *Int   `json:"change"`
	Category string `json:"category,omitempty"`
	Delegate string `json:"delegate,omitempty"`
	Cycle    int    `json:"cycle,omitempty"`
	Level    int    `json:"level,omitempty"`
}

/*
OperationResult represents the operation result in a MXP block

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://MXP.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type OperationResult struct {
	BalanceUpdates      []BalanceUpdates `json:"balance_updates"`
	OriginatedContracts []string         `json:"originated_contracts"`
	Status              string           `json:"status"`
	ConsumedGas         *Int             `json:"consumed_gas,omitempty"`
	Errors              []Error          `json:"errors,omitempty"`
}

/*
Operations represents the operations in a MXP block

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://MXP.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Operations struct {
	Protocol  string     `json:"protocol,omitempty"`
	ChainID   string     `json:"chain_id,omitempty"`
	Hash      string     `json:"hash,omitempty"`
	Branch    string     `json:"branch"`
	Contents  []Contents `json:"contents"`
	Signature string     `json:"signature,omitempty"`
}

/*
Contents represents the contents in a MXP operations

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://MXP.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Contents struct {
	Kind             string            `json:"kind,omitempty"`
	Source           string            `json:"source,omitempty"`
	Fee              *Int              `json:"fee,omitempty"`
	Counter          *Int              `json:"counter,omitempty"`
	GasLimit         *Int              `json:"gas_limit,omitempty"`
	StorageLimit     *Int              `json:"storage_limit,omitempty"`
	Amount           *Int              `json:"amount,omitempty"`
	Destination      string            `json:"destination,omitempty"`
	Delegate         string            `json:"delegate,omitempty"`
	Phk              string            `json:"phk,omitempty"`
	Secret           string            `json:"secret,omitempty"`
	Level            int               `json:"level,omitempty"`
	ManagerPublicKey string            `json:"managerPubkey,omitempty"`
	Balance          *Int              `json:"balance,omitempty"`
	Period           int               `json:"period,omitempty"`
	Proposal         string            `json:"proposal,omitempty"`
	Proposals        []string          `json:"proposals,omitempty"`
	Ballot           string            `json:"ballot,omitempty"`
	Metadata         *ContentsMetadata `json:"metadata,omitempty"`
}

func (c *Contents) equal(contents Contents) (bool, error) {
	x, err := json.Marshal(c)
	if err != nil {
		return false, errors.New("failed to compare")
	}

	y, err := json.Marshal(contents)
	if err != nil {
		return false, errors.New("failed to compare")
	}

	if string(x) == string(y) {
		return true, nil
	}

	return false, nil
}

/*
ContentsMetadata represents the contents metadata in a MXP operations

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://MXP.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type ContentsMetadata struct {
	BalanceUpdates           []BalanceUpdates            `json:"balance_updates"`
	OperationResult          *OperationResult            `json:"operation_result,omitempty"`
	Slots                    []int                       `json:"slots"`
	InternalOperationResults []*InternalOperationResults `json:"internal_operation_results,omitempty"`
}

/*
InternalOperationResults represents a field in contents metadata in a MXP operations

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://MXP.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type InternalOperationResults struct {
	Kind        string           `json:"kind"`
	Source      string           `json:"source"`
	Nonce       uint64           `json:"nonce"`
	Amount      string           `json:"amount"`
	Destination string           `json:"destination"`
	Result      *OperationResult `json:"result"`
}

/*
Error respresents an error for operation results

RPC:
	/chains/<chain_id>/blocks/<block_id> (<dyn>)

Link:
	https://MXP.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance
*/
type Error struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

/*
Head gets all the information about the head block.

Path:
	/chains/<chain_id>/blocks/head (GET)

Link:
	https://MXP.gitlab.io/api/rpc.html#get-chains-chain-id-blocks
*/
func (t *GoMXP) Head() (*Block, error) {
	resp, err := t.get("/chains/main/blocks/head")
	if err != nil {
		return &Block{}, errors.Wrapf(err, "could not get head block")
	}

	var block Block
	err = json.Unmarshal(resp, &block)
	if err != nil {
		return &block, errors.Wrapf(err, "could not get head block")
	}

	return &block, nil
}

/*
Block gets all the information about block.RPC

Path
	/chains/<chain_id>/blocks/<block_id> (GET)
Link
	https://MXP.gitlab.io/api/rpc.html#get-chains-chain-id-blocks

Parameters:

	id:
		hash = <string> : The block hash.
		level = <int> : The block level.
*/
func (t *GoMXP) Block(id interface{}) (*Block, error) {
	blockID, err := idToString(id)
	if err != nil {
		return &Block{}, errors.Wrapf(err, "could not get block '%s'", blockID)
	}

	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s", blockID))
	if err != nil {
		return &Block{}, errors.Wrapf(err, "could not get block '%s'", blockID)
	}

	var block Block
	err = json.Unmarshal(resp, &block)
	if err != nil {
		return &block, errors.Wrapf(err, "could not get block '%s'", blockID)
	}

	return &block, nil
}

/*
OperationHashes is the hashes of all the operations included in the block.

Path:
	../<block_id>/operation_hashes (GET)
Link:
	https://MXP.gitlab.io/api/rpc.html#get-block-id-context-contracts-contract-id-balance

Parameters:

	blockhash:
		The hash of block (height) of which you want to make the query.
*/
func (t *GoMXP) OperationHashes(blockhash string) ([][]string, error) {
	resp, err := t.get(fmt.Sprintf("/chains/main/blocks/%s/operation_hashes", blockhash))
	if err != nil {
		return [][]string{}, errors.Wrapf(err, "could not get operation hashes")
	}

	var operations [][]string
	err = json.Unmarshal(resp, &operations)
	if err != nil {
		return [][]string{}, errors.Wrapf(err, "could not unmarshal operation hashes")
	}

	return operations, nil
}

func idToString(id interface{}) (string, error) {
	switch v := id.(type) {
	case int:
		return strconv.Itoa(v), nil
	case string:
		return v, nil
	default:
		return "", errors.Errorf("id must be block level (int) or block hash (string)")
	}
}
