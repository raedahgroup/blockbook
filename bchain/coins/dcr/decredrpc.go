package dcr

import (
	"blockbook/bchain"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"blockbook/bchain/coins/btc"

	"github.com/golang/glog"
	"github.com/juju/errors"
)

type DecredRPC struct {
	*btc.BitcoinRPC
	client      http.Client
	rpcURL      string
	rpcUser     string
	rpcPassword string
}

// NewDecredRPC returns new DecredRPC instance.
func NewDecredRPC(config json.RawMessage, pushHandler func(bchain.NotificationType)) (bchain.BlockChain, error) {
	b, err := btc.NewBitcoinRPC(config, pushHandler)
	if err != nil {
		return nil, err
	}

	var c btc.Configuration
	err = json.Unmarshal(config, &c)
	if err != nil {
		return nil, errors.Annotate(err, "Invalid configuration file")
	}

	transport := &http.Transport{
		Dial:                (&net.Dialer{KeepAlive: 600 * time.Second}).Dial,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100, // necessary to not to deplete ports
	}

	d := &DecredRPC{
		BitcoinRPC:  b.(*btc.BitcoinRPC),
		client:      http.Client{Timeout: time.Duration(c.RPCTimeout) * time.Second, Transport: transport},
		rpcURL:      c.RPCURL,
		rpcUser:     c.RPCUser,
		rpcPassword: c.RPCPass,
	}

	d.BitcoinRPC.RPCMarshaler = btc.JSONMarshalerV1{}
	d.BitcoinRPC.ChainConfig.SupportsEstimateSmartFee = false

	return d, nil
}

// Initialize initializes DecredRPC instance.
func (d *DecredRPC) Initialize() error {
	chainInfo, err := d.GetChainInfo()
	if err != nil {
		return err
	}

	chainName := chainInfo.Chain
	glog.Info("Chain name ", chainName)

	params := GetChainParams(chainName)

	// always create parser
	d.BitcoinRPC.Parser = NewDecredParser(params, d.BitcoinRPC.ChainConfig)

	// parameters for getInfo request
	if params.Net == MainnetMagic {
		d.BitcoinRPC.Testnet = false
		d.BitcoinRPC.Network = "livenet"
	} else {
		d.BitcoinRPC.Testnet = true
		d.BitcoinRPC.Network = "testnet"
	}

	glog.Info("rpc: block chain ", params.Name)

	return nil
}

type GenericCmd struct {
	ID     int           `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

type GetBlockChainInfoResult struct {
	Chain                string  `json:"chain"`
	Blocks               int64   `json:"blocks"`
	Headers              int64   `json:"headers"`
	SyncHeight           int64   `json:"syncheight"`
	BestBlockHash        string  `json:"bestblockhash"`
	Difficulty           uint32  `json:"difficulty"`
	VerificationProgress float64 `json:"verificationprogress"`
	ChainWork            string  `json:"chainwork"`
	InitialBlockDownload bool    `json:"initialblockdownload"`
	MaxBlockSize         int64   `json:"maxblocksize"`
}

type GetNetworkInfoResult struct {
	Version         int32   `json:"version"`
	ProtocolVersion int32   `json:"protocolversion"`
	TimeOffset      int64   `json:"timeoffset"`
	Connections     int32   `json:"connections"`
	RelayFee        float64 `json:"relayfee"`
}

type GetBestBlockResult struct {
	Hash   string `json:"hash"`
	Height int64  `json:"height"`
}

type GetBlockHashCmd struct {
	ID     int           `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

type GetBlockHashResult struct {
}

type GetBlockHeaderCmd struct {
	ID     int           `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

type GetBlockInfoResult struct {
	Hash          string  `json:"hash"`
	Confirmations int64   `json:"confirmations"`
	Version       int32   `json:"version"`
	MerkleRoot    string  `json:"merkleroot"`
	StakeRoot     string  `json:"stakeroot"`
	VoteBits      uint16  `json:"votebits"`
	FinalState    string  `json:"finalstate"`
	Voters        uint16  `json:"voters"`
	FreshStake    uint8   `json:"freshstake"`
	Revocations   uint8   `json:"revocations"`
	PoolSize      uint32  `json:"poolsize"`
	Bits          string  `json:"bits"`
	SBits         float64 `json:"sbits"`
	Height        uint32  `json:"height"`
	Size          uint32  `json:"size"`
	Time          int64   `json:"time"`
	Nonce         uint32  `json:"nonce"`
	ExtraData     string  `json:"extradata"`
	StakeVersion  uint32  `json:"stakeversion"`
	Difficulty    float64 `json:"difficulty"`
	ChainWork     string  `json:"chainwork"`
	PreviousHash  string  `json:"previousblockhash,omitempty"`
	NextHash      string  `json:"nextblockhash,omitempty"`
}

type GetBlockCmd struct {
	ID     int           `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

type GetBlockResult struct {
}

type GetTransactionCmd struct {
	ID     int           `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

type Vin struct {
	Coinbase    string     `json:"coinbase"`
	Stakebase   string     `json:"stakebase"`
	Txid        string     `json:"txid"`
	Vout        uint32     `json:"vout"`
	Tree        int8       `json:"tree"`
	Sequence    uint32     `json:"sequence"`
	AmountIn    float64    `json:"amountin"`
	BlockHeight uint32     `json:"blockheight"`
	BlockIndex  uint32     `json:"blockindex"`
	ScriptSig   *ScriptSig `json:"scriptSig"`
}

type ScriptPubKeyResult struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex,omitempty"`
	ReqSigs   int32    `json:"reqSigs,omitempty"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses,omitempty"`
	CommitAmt *float64 `json:"commitamt,omitempty"`
}

type Vout struct {
	Value        float64            `json:"value"`
	N            uint32             `json:"n"`
	Version      uint16             `json:"version"`
	ScriptPubKey ScriptPubKeyResult `json:"scriptPubKey"`
}

type GetTransactionResult struct {
	Hex           string `json:"hex"`
	Txid          string `json:"txid"`
	Version       int32  `json:"version"`
	LockTime      uint32 `json:"locktime"`
	Vin           []Vin  `json:"vin"`
	Vout          []Vout `json:"vout"`
	Expiry        uint32 `json:"expiry"`
	BlockHash     string `json:"blockhash,omitempty"`
	BlockHeight   int64  `json:"blockheight,omitempty"`
	BlockIndex    uint32 `json:"blockindex,omitempty"`
	Confirmations int64  `json:"confirmations,omitempty"`
	Time          int64  `json:"time,omitempty"`
	Blocktime     int64  `json:"blocktime,omitempty"`
}

type EstimateSmartFeeResult struct {
	FeeRate float64  `json:"feerate"`
	Errors  []string `json:"errors"`
	Blocks  int64    `json:"blocks"`
}

type SendRawTransactionResult struct {
}

func (d *DecredRPC) GetChainInfo() (*bchain.ChainInfo, error) {
	blockchainInfoRequest := GenericCmd{
		ID:     1,
		Method: "getblockchaininfo",
	}
	blockchainInfoResult := GetBlockChainInfoResult{}
	err := d.Call(blockchainInfoRequest, &blockchainInfoResult)
	if err != nil {
		return nil, err
	}

	networkInfoRequest := GenericCmd{
		ID:     2,
		Method: "getnetworkinfo",
	}
	networkInfoResult := GetNetworkInfoResult{}
	err = d.Call(networkInfoRequest, &networkInfoResult)
	if err != nil {
		return nil, err
	}

	chainInfo := &bchain.ChainInfo{
		Chain:           blockchainInfoResult.Chain,
		Blocks:          int(blockchainInfoResult.Blocks),
		Headers:         int(blockchainInfoResult.Headers),
		Bestblockhash:   blockchainInfoResult.BestBlockHash,
		Difficulty:      strconv.Itoa(int(blockchainInfoResult.Difficulty)),
		SizeOnDisk:      blockchainInfoResult.SyncHeight,
		Version:         strconv.Itoa(int(networkInfoResult.Version)),
		Subversion:      "",
		ProtocolVersion: strconv.Itoa(int(networkInfoResult.ProtocolVersion)),
		Timeoffset:      float64(networkInfoResult.TimeOffset),
		Warnings:        "",
	}
	return chainInfo, nil
}

func (d *DecredRPC) getBestBlock() (GetBestBlockResult, error) {
	bestBlockRequest := GenericCmd{
		ID:     1,
		Method: "getbestblock",
	}
	bestBlockResult := GetBestBlockResult{}
	err := d.Call(bestBlockRequest, &bestBlockResult)
	return bestBlockResult, err
}

func (d *DecredRPC) GetBestBlockHash() (string, error) {
	bestBlock, err := d.getBestBlock()
	if err != nil {
		return "", err
	}

	return bestBlock.Hash, nil
}

func (d *DecredRPC) GetBestBlockHeight() (uint32, error) {
	bestBlock, err := d.getBestBlock()
	if err != nil {
		return 0, err
	}

	return uint32(bestBlock.Height), err
}

func (d *DecredRPC) GetBlockHash(height uint32) (string, error) {
	blockHashRequest := GetBlockHashCmd{
		ID:     1,
		Method: "getblockhash",
		Params: []interface{}{height},
	}
	blockHashResponse := ""
	err := d.Call(blockHashRequest, &blockHashResponse)
	return blockHashResponse, err
}

func (d *DecredRPC) GetBlockHeader(hash string) (*bchain.BlockHeader, error) {
	blockInfo, err := d.getBlockInfo(hash)
	if err != nil {
		return nil, err
	}

	header := &bchain.BlockHeader{
		Hash:          blockInfo.Hash,
		Prev:          blockInfo.PreviousHash,
		Next:          blockInfo.NextHash,
		Height:        blockInfo.Height,
		Confirmations: int(blockInfo.Confirmations),
		Size:          int(blockInfo.Size),
		Time:          int64(blockInfo.Size),
	}

	return header, nil
}

func (d *DecredRPC) getBlockInfo(hash string) (GetBlockInfoResult, error) {
	blockHeaderRequest := GetBlockHeaderCmd{
		ID:     1,
		Method: "getblockheader",
		Params: []interface{}{hash},
	}
	blockInfoResult := GetBlockInfoResult{}
	err := d.Call(blockHeaderRequest, &blockInfoResult)

	return blockInfoResult, err
}

func (d *DecredRPC) GetBlockHeaderByHeight(height uint32) (*bchain.BlockHeader, error) {
	return nil, nil
}

func (d *DecredRPC) GetBlock(hash string, height uint32) (*bchain.Block, error) {
	return nil, nil
}

func (d *DecredRPC) GetBlockInfo(hash string) (*bchain.BlockInfo, error) {
	if hash == "" {
		return nil, bchain.ErrBlockNotFound
	}

	blockInfo, err := d.getBlockInfo(hash)
	if err != nil {
		return nil, err
	}

	header := bchain.BlockHeader{
		Hash:          blockInfo.Hash,
		Prev:          blockInfo.PreviousHash,
		Next:          blockInfo.NextHash,
		Height:        blockInfo.Height,
		Confirmations: int(blockInfo.Confirmations),
		Size:          int(blockInfo.Size),
		Time:          int64(blockInfo.Size),
	}

	bInfo := &bchain.BlockInfo{
		BlockHeader: header,
		MerkleRoot:  blockInfo.MerkleRoot,
		Version:     json.Number(blockInfo.Version),
		Txids:       []string{},
	}

	return bInfo, nil
}

func (d *DecredRPC) GetMempoolTransactions() ([]string, error) {
	return nil, nil
}

func (d *DecredRPC) GetTransaction(txid string) (*bchain.Tx, error) {
	if txid == "" {
		return nil, bchain.ErrTxidMissing
	}

	getTxRequest := GetTransactionCmd{
		ID:     1,
		Method: "getrawtransaction",
		Params: []interface{}{txid},
	}
	getTxResponse := GetTransactionResult{}
	err := d.Call(getTxRequest, &getTxResponse)
	if err != nil {
		return nil, err
	}

	var Vin []bchain.Vin
	var Vout []bchain.Vout

	for _, vin := range getTxResponse.Vin {
		item := bchain.Vin{
			Coinbase: vin.Coinbase,
			Txid:     vin.Txid,
			Vout:     vin.Vout,
			ScriptSig: bchain.ScriptSig{
				Hex: vin.ScriptSig.Hex,
			},
			Sequence: vin.Sequence,
		}

		Vin = append(Vin, item)
	}

	for _, vout := range getTxResponse.Vout {
		item := bchain.Vout{
			//ValueSat: big.Int(vout.Value),
			JsonValue: json.Number(int64(vout.Value)),
			N:         vout.N,
			ScriptPubKey: bchain.ScriptPubKey{
				Hex:       vout.ScriptPubKey.Hex,
				Addresses: vout.ScriptPubKey.Addresses,
			},
		}

		Vout = append(Vout, item)
	}

	tx := &bchain.Tx{
		Hex:           getTxResponse.Hex,
		Txid:          getTxResponse.Txid,
		Version:       getTxResponse.Version,
		LockTime:      getTxResponse.LockTime,
		Vin:           Vin,
		Vout:          Vout,
		Confirmations: uint32(getTxResponse.Confirmations),
		Time:          getTxResponse.Time,
		Blocktime:     getTxResponse.Blocktime,
	}

	return tx, nil
}

func (d *DecredRPC) GetTransactionForMempool(txid string) (*bchain.Tx, error) {
	return nil, nil
}

func (d *DecredRPC) GetTransactionSpecific(tx *bchain.Tx) (json.RawMessage, error) {
	return nil, nil
}

func (d *DecredRPC) EstimateSmartFee(blocks int, conservative bool) (big.Int, error) {
	estimateSmartFeeRequest := GenericCmd{
		ID:     1,
		Method: "estimatesmartfee",
		Params: []interface{}{blocks},
	}
	estimateSmartFeeResult := EstimateSmartFeeResult{}

	err := d.Call(estimateSmartFeeRequest, &estimateSmartFeeResult)
	if err != nil {
		return *big.NewInt(0), nil
	}

	return *big.NewInt(int64(estimateSmartFeeResult.FeeRate)), nil
}

func (d *DecredRPC) SendRawTransaction(tx string) (string, error) {
	sendRawTxRequest := &GenericCmd{
		ID:     1,
		Method: "sendrawtransaction",
		Params: []interface{}{tx},
	}

	var res string
	err := d.Call(sendRawTxRequest, res)
	if err != nil {
		return "", err
	}

	return res, nil
}

// Call calls Backend RPC interface, using RPCMarshaler interface to marshall the request
func (d *DecredRPC) Call(req interface{}, res interface{}) error {
	httpData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest("POST", d.rpcURL, bytes.NewBuffer(httpData))
	if err != nil {
		return err
	}
	httpReq.SetBasicAuth(d.rpcUser, d.rpcPassword)
	httpRes, err := d.client.Do(httpReq)
	// in some cases the httpRes can contain data even if it returns error
	// see http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/
	if httpRes != nil {
		defer httpRes.Body.Close()
	}
	if err != nil {
		return err
	}

	// if server returns HTTP error code it might not return json with response
	// handle both cases
	if httpRes.StatusCode != 200 {
		err = safeDecodeResponse(httpRes.Body, &res)
		if err != nil {
			return errors.Errorf("%v %v", httpRes.Status, err)
		}
		return nil
	}
	return safeDecodeResponse(httpRes.Body, &res)
}

func safeDecodeResponse(body io.ReadCloser, res *interface{}) (err error) {
	var data []byte
	defer func() {
		if r := recover(); r != nil {
			glog.Error("unmarshal json recovered from panic: ", r, "; data: ", string(data))
			debug.PrintStack()
			if len(data) > 0 && len(data) < 2048 {
				err = errors.Errorf("Error: %v", string(data))
			} else {
				err = errors.New("Internal error")
			}
		}
	}()
	data, err = ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	error := json.Unmarshal(data, res)
	return error
}
