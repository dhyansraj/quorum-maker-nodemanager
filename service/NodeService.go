package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/contractcompiler"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/magiconair/properties"
	log "github.com/sirupsen/logrus"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/client"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/contractclient"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/contracthandler"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/env"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/util"
	"gopkg.in/gomail.v2"
)

type ConnectionInfo struct {
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	Enode string `json:"enode"`
}

type PendingRequests struct {
	NodeName string `json:"nodeName"`
	Enode    string `json:"enode"`
	Message  string `json:"message"`
	EnodeID  string `json:"enodeid"`
	IP       string `json:"ip"`
}

type NodeInfo struct {
	NodeName       string           `json:"nodeName"`
	NodeCount      int              `json:"nodeCount"`
	TotalNodeCount int              `json:"totalNodeCount"`
	Active         string           `json:"active"`
	ConnectionInfo ConnectionInfo   `json:"connectionInfo"`
	RaftRole       string           `json:"raftRole"`
	RaftID         int              `json:"raftID"`
	BlockNumber    int64            `json:"blockNumber"`
	PendingTxCount int              `json:"pendingTxCount"`
	Genesis        string           `json:"genesis"`
	AdminInfo      client.AdminInfo `json:"adminInfo"`
}

type JoinNetworkRequest struct {
	EnodeID   string `json:"enode-id,omitempty"`
	IPAddress string `json:"ip-address,omitempty"`
	Nodename  string `json:"nodename,omitempty"`
}

type GetGenesisResponse struct {
	ContstellationPort string `json:"contstellation-port"`
	NetID              string `json:"netID"`
	Genesis            string `json:"genesis"`
}

type BlockDetailsResponse struct {
	Number           int64                        `json:"number"`
	Hash             string                       `json:"hash"`
	ParentHash       string                       `json:"parentHash"`
	Nonce            string                       `json:"nonce"`
	Sha3Uncles       string                       `json:"sha3Uncles"`
	LogsBloom        string                       `json:"logsBloom"`
	TransactionsRoot string                       `json:"transactionsRoot"`
	StateRoot        string                       `json:"stateRoot"`
	Miner            string                       `json:"miner"`
	Difficulty       int64                        `json:"difficulty"`
	TotalDifficulty  int64                        `json:"totalDifficulty"`
	ExtraData        string                       `json:"extraData"`
	Size             int64                        `json:"size"`
	GasLimit         int64                        `json:"gasLimit"`
	GasUsed          int64                        `json:"gasUsed"`
	Timestamp        int64                        `json:"timestamp"`
	Transactions     []TransactionDetailsResponse `json:"transactions"`
	Uncles           []string                     `json:"uncles"`
	TimeElapsed      int64                        `json:"TimeElapsed"`
}

type TransactionDetailsResponse struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      int64  `json:"blockNumber"`
	From             string `json:"from"`
	Gas              int64  `json:"gas"`
	GasPrice         int64  `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            int64  `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex int64  `json:"transactionIndex"`
	Value            int64  `json:"value"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
	TransactionType  string `json:"transactionType"`
	TimeElapsed      int64  `json:"TimeElapsed"`
}

type TransactionReceiptResponse struct {
	BlockHash         string                         `json:"blockHash"`
	BlockNumber       int64                          `json:"blockNumber"`
	ContractAddress   string                         `json:"contractAddress"`
	CumulativeGasUsed int64                          `json:"cumulativeGasUsed"`
	From              string                         `json:"from"`
	Gas               int64                          `json:"gas"`
	GasPrice          int64                          `json:"gasPrice"`
	GasUsed           int64                          `json:"gasUsed"`
	Input             string                         `json:"input"`
	Logs              []Logs                         `json:"logs"`
	LogsBloom         string                         `json:"logsBloom"`
	Nonce             int64                          `json:"nonce"`
	Root              string                         `json:"root"`
	To                string                         `json:"to"`
	TransactionHash   string                         `json:"transactionHash"`
	TransactionIndex  int64                          `json:"transactionIndex"`
	Value             int64                          `json:"value"`
	V                 string                         `json:"v"`
	R                 string                         `json:"r"`
	S                 string                         `json:"s"`
	TransactionType   string                         `json:"transactionType"`
	TimeElapsed       int64                          `json:"TimeElapsed"`
	DecodedInputs     []contractclient.ParamTableRow `json:"decodedInputs,omitempty"`
	FunctionDetails   string                         `json:"functionDetails,omitempty"`
	DecodeFailed      DecodeFailure                  `json:"decodeFailed,omitempty"`
}

type DecodeFailure struct {
	Label string `json:"label"`
	Type  string `json:"type"`
}

type Logs struct {
	Address          string   `json:"address"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      int64    `json:"blockNumber"`
	Data             string   `json:"data"`
	LogIndex         int64    `json:"logIndex"`
	Topics           []string `json:"topics"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex int64    `json:"transactionIndex"`
}

type JoinNetworkResponse struct {
	EnodeID string `json:"enode-id"`
	Status  string `json:"status"`
}



type CreateNetworkScriptArgs struct {
	Nodename          string `json:"nodename,omitempty"`
	CurrentIP         string `json:"currentIP,omitempty"`
	RPCPort           string `json:"rpcPort,omitempty"`
	WhisperPort       string `json:"whisperPort,omitempty"`
	TesseraPort string `json:"tesseraPort,omitempty"`
	RaftPort          string `json:"raftPort,omitempty"`
	NodeManagerPort   string `json:"nodeManagerPort,omitempty"`
}

type JoinNetworkScriptArgs struct {
	Nodename              string `json:"nodename,omitempty"`
	CurrentIP             string `json:"currentIP,omitempty"`
	RPCPort               string `json:"rpcPort,omitempty"`
	WhisperPort           string `json:"whisperPort,omitempty"`
	TesseraPort     string `json:"tesseraPort,omitempty"`
	RaftPort              string `json:"raftPort,omitempty"`
	NodeManagerPort       string `json:"nodeManagerPort,omitempty"`
	MasterNodeManagerPort string `json:"masterNodeManagerPort,omitempty"`
	MasterIP              string `json:"masterIP,omitempty"`
}

type SuccessResponse struct {
	Status string `json:"statusMessage"`
}

type SuccessResponseBool struct {
	Status bool `json:"statusMessage"`
}

type LatestBlockResponse struct {
	LatestBlockNumber int64 `json:"latestBlockNumber"`
	TimeElapsed       int64 `json:"TimeElapsed"`
}

type NodeList struct {
	NodeName  string `json:"nodeName"`
	Role      string `json:"role,omitempty"`
	PublicKey string `json:"publicKey"`
	IP        string `json:"ip,omitempty"`
	Enode     string `json:"enode,omitempty"`
}

type MailServerConfig struct {
	Host          string `json:"smtpServerHost"`
	Port          string `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	RecipientList string `json:"recipientList"`
}

type LatencyResponse struct {
	EnodeID string `json:"enode-id"`
	Latency string `json:"latency"`
}

type NodeServiceImpl struct {
	Url string
}

type ChartInfo struct {
	TimeStamp        int `json:"timeStamp"`
	BlockCount       int `json:"blockCount"`
	TransactionCount int `json:"transactionCount"`
}

type ContractTableRow struct {
	ContractAdd  string `json:"contractAddress"`
	ContractName string `json:"contractName"`
	ABIContent   string `json:"abi"`
	Sender       string `json:"sender"`
	ContractType string `json:"contractType"`
	Description  string `json:"description"`
	Timestamp    string `json:"timestamp"`
}

type ContractCounter struct {
	TotalContracts int `json:"totalContracts"`
	ABIavailable   int `json:"abis"`
}

type CrawledABI struct {
	Filename         string `json:"filename"`
	ModificationTime int64  `json:"modificationTime"`
	Processed        bool   `json:"processed"`
	Contractname     string `json:"contractname"`
}

type contractJSONTruffle struct {
	Abi             []interface{} `json:"abi"`
	Interface       []interface{} `json:"interface"`
	Bytecode        string        `json:"bytecode"`
	ContractName    string        `json:"contractName"`
	Address         string        `json:"address"`

}

type ethAccount struct {
	AccountAddress string `json:"accountAddress"`
	Coinbase       bool   `json:"coinbase"`
	Balance        string `json:"balance"`
}

var txnMap = map[string]TransactionReceiptResponse{}
var abiMap = map[string]string{}
var contractCrawlerMutex = 0
var crawledABIs []CrawledABI
var abiCrawlerMutex = 0

var contDescriptionMap = map[string]string{}
var contTypeMap = map[string]string{}
var contTimeMap = map[string]string{}
var contSenderMap = map[string]string{}
var contNameMap = map[string]string{}
var chartSize = 10
var warning = 0
var lastCrawledBlock = 0
var mailServerConfig MailServerConfig

func (nsi *NodeServiceImpl) getGenesis(url string) (response GetGenesisResponse) {

	b, err := ioutil.ReadFile(env.GetAppConfig().NodeDir + "/genesis.json")
	if err != nil {
		//log.Println(err)
	}
	genesis := string(b)
	genesis = strings.Replace(genesis, "\n", "", -1)

	response = GetGenesisResponse{env.GetSetupConf().TesseraPort, env.GetSetupConf().NetworkId, genesis}
	return response
}

func (nsi *NodeServiceImpl) joinNetwork(enode string, url string) string {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	raftId := ethClient.RaftAddPeer(enode)

	collatedInfo := fmt.Sprint(raftId, ":", env.GetSetupConf().ContractAdd)
	return collatedInfo
}

//@TODO: If this function is repeatedly called from UI, please cache the static informations.
func (nsi *NodeServiceImpl) getCurrentNode(url string) NodeInfo {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	fromAddress := ethClient.Coinbase()

	nms := contractclient.NetworkMapContractClient{client.EthClient{url}, contracthandler.ContractParam{fromAddress, env.GetSetupConf().ContractAdd, "", nil}}

	totalCount := len(nms.GetNodeDetailsList())
	var activeStatus string
	active := ethClient.NetListening()
	if active == true {
		activeStatus = "active"
	} else {
		activeStatus = "inactive"
	}
	otherPeersResponse := ethClient.AdminPeers()
	count := len(otherPeersResponse)
	count = count + 1

	raftIdInt, err := strconv.Atoi(env.GetSetupConf().RaftId)
	if err != nil {
		log.Println(err)
	}

	rpcPortInt, err := strconv.Atoi(env.GetSetupConf().RpcPort)
	if err != nil {
		log.Println(err)
	}

	thisAdminInfo := ethClient.AdminNodeInfo()
	enode := thisAdminInfo.Enode

	pendingTxResponse := ethClient.PendingTransactions()
	pendingTxCount := len(pendingTxResponse)

	blockNumber := ethClient.BlockNumber()
	blockNumberInt := util.HexStringtoInt64(blockNumber)

	raftRole := ethClient.RaftRole()

	raftRole = strings.TrimSuffix(raftRole, "\n")

	b, err := ioutil.ReadFile(env.GetAppConfig().NodeDir + "/genesis.json")

	if err != nil {
		//log.Println(err)
	}

	genesis := string(b)
	genesis = strings.Replace(genesis, "\n", "", -1)
	conn := ConnectionInfo{env.GetSetupConf().CurrentIp, rpcPortInt, enode}
	responseObj := NodeInfo{env.GetSetupConf().Nodename, count, totalCount, activeStatus, conn, raftRole, raftIdInt, blockNumberInt, pendingTxCount, genesis, thisAdminInfo}
	return responseObj
}

func (nsi *NodeServiceImpl) getOtherPeer(peerId string, url string) client.AdminPeers {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	otherPeersResponse := ethClient.AdminPeers()
	for _, item := range otherPeersResponse {
		if item.ID == peerId {
			peerResponse := item
			return peerResponse
		}
	}
	return client.AdminPeers{}
}

func (nsi *NodeServiceImpl) getOtherPeers(url string) []client.AdminPeers {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	otherPeersResponse := ethClient.AdminPeers()
	return otherPeersResponse
}

func (nsi *NodeServiceImpl) getPendingTransactions(url string) []TransactionDetailsResponse {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	clientPendingTxResponses := ethClient.PendingTransactions()

	pendingTxResponse := make([]TransactionDetailsResponse, len(clientPendingTxResponses))

	for i, clientPendingTxResponse := range clientPendingTxResponses {

		pendingTxResponse[i] = ConvertToReadable(clientPendingTxResponse, true, true)
	}

	return pendingTxResponse
}

func (nsi *NodeServiceImpl) getBlockInfo(blockno int64, url string) BlockDetailsResponse {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	blockNoHex := strconv.FormatInt(blockno, 16)
	bNoHex := fmt.Sprint("0x", blockNoHex)
	var blockResponse BlockDetailsResponse
	blockResponseClient := ethClient.GetBlockByNumber(bNoHex)
	currentTime := time.Now().Unix()
	creationTime := util.HexStringtoInt64(blockResponseClient.Timestamp)
	creationTimeUnix := creationTime / 1000000000
	elapsedTime := currentTime - creationTimeUnix
	blockResponse.TimeElapsed = elapsedTime

	//@TODO: Create a utility function to convert block object to readable object.
	blockResponse.Number = util.HexStringtoInt64(blockResponseClient.Number)
	blockResponse.Difficulty = util.HexStringtoInt64(blockResponseClient.Difficulty)
	blockResponse.TotalDifficulty = util.HexStringtoInt64(blockResponseClient.TotalDifficulty)
	blockResponse.Size = util.HexStringtoInt64(blockResponseClient.Size)
	blockResponse.GasLimit = util.HexStringtoInt64(blockResponseClient.GasLimit)
	blockResponse.GasUsed = util.HexStringtoInt64(blockResponseClient.GasUsed)
	blockResponse.Timestamp = util.HexStringtoInt64(blockResponseClient.Timestamp)
	blockResponse.Hash = blockResponseClient.Hash
	blockResponse.ParentHash = blockResponseClient.ParentHash
	blockResponse.Nonce = blockResponseClient.Nonce
	blockResponse.Sha3Uncles = blockResponseClient.Sha3Uncles
	blockResponse.LogsBloom = blockResponseClient.LogsBloom
	blockResponse.TransactionsRoot = blockResponseClient.TransactionsRoot
	blockResponse.StateRoot = blockResponseClient.StateRoot
	blockResponse.Miner = blockResponseClient.Miner
	blockResponse.ExtraData = blockResponseClient.ExtraData
	blockResponse.Uncles = blockResponseClient.Uncles
	txnNo := len(blockResponseClient.Transactions)
	txResponse := make([]TransactionDetailsResponse, txnNo)
	for i, clientTransactions := range blockResponseClient.Transactions {

		txGetClient := ethClient.GetTransactionByHash(clientTransactions.Hash)
		private := ethClient.GetQuorumPayload(txGetClient.Input)
		txResponse[i] = ConvertToReadable(clientTransactions, false, (private == "0x"))

	}
	blockResponse.Transactions = txResponse
	return blockResponse
}

func (nsi *NodeServiceImpl) getLatestBlockInfo(count string, reference string, url string) []BlockDetailsResponse {
	countValInt, _ := strconv.Atoi(count)
	countVal := int64(countValInt)
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	var blockNumber int64
	if reference == "" {
		blockNumber = util.HexStringtoInt64(ethClient.BlockNumber())
	} else {
		var err error
		blockNumber, err = strconv.ParseInt(reference, 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		blockNumber = blockNumber - 1
	}
	start := blockNumber - countVal + 1
	blockResponse := make([]BlockDetailsResponse, countVal)

	//@TODO: call the utility function to convert to readable block object
	for i := start; i <= blockNumber; i++ {
		blockNoHex := strconv.FormatInt(i, 16)
		bNoHex := fmt.Sprint("0x", blockNoHex)
		blockResponseClient := ethClient.GetBlockByNumber(bNoHex)
		blockResponse[blockNumber-i].Number = util.HexStringtoInt64(blockResponseClient.Number)
		blockResponse[blockNumber-i].Hash = blockResponseClient.Hash
		currentTime := time.Now().Unix()
		creationTime := util.HexStringtoInt64(blockResponseClient.Timestamp)
		creationTimeUnix := creationTime / 1000000000
		elapsedTime := currentTime - creationTimeUnix
		blockResponse[blockNumber-i].TimeElapsed = elapsedTime
		txnNo := len(blockResponseClient.Transactions)
		txResponse := make([]TransactionDetailsResponse, txnNo)

		for i, clientTransactions := range blockResponseClient.Transactions {

			txGetClient := ethClient.GetTransactionByHash(clientTransactions.Hash)
			private := ethClient.GetQuorumPayload(txGetClient.Input)
			txResponse[i] = ConvertToReadable(clientTransactions, false, (private == "0x"))

		}

		blockResponse[blockNumber-i].Transactions = txResponse
	}
	return blockResponse
}

func (nsi *NodeServiceImpl) getLatestTransactionInfo(count string, url string) []BlockDetailsResponse {
	countValInt, _ := strconv.Atoi(count)
	countVal := int64(countValInt)
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	blockNumber := util.HexStringtoInt64(ethClient.BlockNumber())
	start := blockNumber - countVal + 1
	blockResponse := make([]BlockDetailsResponse, countVal)
	for i := start; i <= blockNumber; i++ {
		blockNoHex := strconv.FormatInt(i, 16)
		bNoHex := fmt.Sprint("0x", blockNoHex)
		blockResponseClient := ethClient.GetBlockByNumber(bNoHex)
		currentTime := time.Now().Unix()
		creationTime := util.HexStringtoInt64(blockResponseClient.Timestamp)
		creationTimeUnix := creationTime / 1000000000
		elapsedTime := currentTime - creationTimeUnix
		blockResponse[blockNumber-i].TimeElapsed = elapsedTime
		blockResponse[blockNumber-i].Number = util.HexStringtoInt64(blockResponseClient.Number)
		txnNo := len(blockResponseClient.Transactions)
		txResponse := make([]TransactionDetailsResponse, txnNo)

		for i, clientTransactions := range blockResponseClient.Transactions {

			txGetClient := ethClient.GetTransactionByHash(clientTransactions.Hash)
			private := ethClient.GetQuorumPayload(txGetClient.Input)
			txResponse[i] = ConvertToReadable(clientTransactions, false, (private == "0x"))

		}

		blockResponse[blockNumber-i].Transactions = txResponse
	}
	return blockResponse
}

func (nsi *NodeServiceImpl) getTransactionInfo(txno string, url string) TransactionDetailsResponse {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	var txResponse TransactionDetailsResponse
	txResponseClient := ethClient.GetTransactionByHash(txno)

	private := ethClient.GetQuorumPayload(txResponseClient.Input)
	txResponse = ConvertToReadable(txResponseClient, false, (private == "0x"))

	blockResponseClient := ethClient.GetBlockByNumber(txResponseClient.BlockNumber)
	currentTime := time.Now().Unix()
	creationTime := util.HexStringtoInt64(blockResponseClient.Timestamp)
	creationTimeUnix := creationTime / 1000000000
	elapsedTime := currentTime - creationTimeUnix
	txResponse.TimeElapsed = elapsedTime
	return txResponse
}

func (nsi *NodeServiceImpl) getTransactionReceipt(txno string, url string) TransactionReceiptResponse {
	if txnMap[txno].TransactionHash == "" {
		txResponse := populateTransactionObject(txno, url)
		decodeTransactionObject(&txResponse, url)
		return txResponse
	} else {
		txnDetails := txnMap[txno]
		calculateTimeElapsed(&txnDetails, url)
		txnMap[txno] = txnDetails
		return txnDetails
	}
}

func populateTransactionObject(txno string, url string) TransactionReceiptResponse {
	ethClient := client.EthClient{url}
	getTransaction := ethClient.GetTransactionByHash(txno)
	var txResponse TransactionReceiptResponse
	getTransactionReceipt := ethClient.GetTransactionReceipt(txno)
	txResponse.BlockNumber = util.HexStringtoInt64(getTransactionReceipt.BlockNumber)
	txResponse.CumulativeGasUsed = util.HexStringtoInt64(getTransactionReceipt.CumulativeGasUsed)
	txResponse.GasUsed = util.HexStringtoInt64(getTransactionReceipt.GasUsed)
	txResponse.TransactionIndex = util.HexStringtoInt64(getTransactionReceipt.TransactionIndex)
	txResponse.BlockHash = getTransactionReceipt.BlockHash
	txResponse.From = getTransactionReceipt.From
	txResponse.ContractAddress = getTransactionReceipt.ContractAddress
	txResponse.LogsBloom = getTransactionReceipt.LogsBloom
	txResponse.Root = getTransactionReceipt.Root
	txResponse.To = getTransactionReceipt.To
	txResponse.TransactionHash = getTransactionReceipt.TransactionHash
	txResponse.Gas = util.HexStringtoInt64(getTransaction.Gas)
	txResponse.GasPrice = util.HexStringtoInt64(getTransaction.GasPrice)
	txResponse.Input = getTransaction.Input
	txResponse.Nonce = util.HexStringtoInt64(getTransaction.Nonce)
	txResponse.Value = util.HexStringtoInt64(getTransaction.Value)
	txResponse.V = getTransaction.V
	txResponse.R = getTransaction.R
	txResponse.S = getTransaction.S
	eventNo := len(getTransactionReceipt.Logs)
	txResponseBuffer := make([]Logs, eventNo)
	for i := 0; i < eventNo; i++ {
		txResponseBuffer[i].BlockNumber = util.HexStringtoInt64(getTransactionReceipt.Logs[i].BlockNumber)
		txResponseBuffer[i].LogIndex = util.HexStringtoInt64(getTransactionReceipt.Logs[i].LogIndex)
		txResponseBuffer[i].TransactionIndex = util.HexStringtoInt64(getTransactionReceipt.Logs[i].TransactionIndex)
		txResponseBuffer[i].Address = getTransactionReceipt.Logs[i].Address
		txResponseBuffer[i].BlockHash = getTransactionReceipt.Logs[i].BlockHash
		txResponseBuffer[i].Data = getTransactionReceipt.Logs[i].Data
		txResponseBuffer[i].TransactionHash = getTransactionReceipt.Logs[i].TransactionHash
		txResponseBuffer[i].Topics = getTransactionReceipt.Logs[i].Topics
	}
	txResponse.Logs = txResponseBuffer
	blockResponseClient := ethClient.GetBlockByNumber(getTransactionReceipt.BlockNumber)
	currentTime := time.Now().Unix()
	creationTime := util.HexStringtoInt64(blockResponseClient.Timestamp)
	creationTimeUnix := creationTime / 1000000000
	elapsedTime := currentTime - creationTimeUnix
	txResponse.TimeElapsed = elapsedTime
	return txResponse
}

func decodeTransactionObject(txnDetails *TransactionReceiptResponse, url string) {
	var quorumPayload string
	var decoded bool

	ethClient := client.EthClient{url}

	if util.HexStringtoInt64(txnDetails.V) == 37 || util.HexStringtoInt64(txnDetails.V) == 38 {
		quorumPayload = ethClient.GetQuorumPayload(txnDetails.Input)
		if quorumPayload == "0x" {
			txnDetails.TransactionType = "Hash Only"
		} else {
			txnDetails.TransactionType = "Private"
		}
	} else {
		txnDetails.TransactionType = "Public"

	}
	if txnDetails.ContractAddress == "" {
		if txnDetails.TransactionType == "Private" && abiMap[txnDetails.To] != "" && abiMap[txnDetails.To] != "missing" {
			txnDetails.Input = quorumPayload
			decodedData, functionDetails := contractclient.ABIParser(txnDetails.To, abiMap[txnDetails.To], quorumPayload)
			if decodedData[0].Key == "decodeFailed" {
				var decodeFail DecodeFailure
				decodeFail.Label = decodedData[0].Value
				decodeFail.Type = "red"
				txnDetails.DecodeFailed = decodeFail
			} else {
				txnDetails.DecodedInputs = decodedData
			}
			if functionDetails != "" {
				txnDetails.FunctionDetails = functionDetails
				decoded = true
			}
		} else if txnDetails.TransactionType == "Public" && abiMap[txnDetails.To] != "" && abiMap[txnDetails.To] != "missing" {
			decodedData, functionDetails := contractclient.ABIParser(txnDetails.To, abiMap[txnDetails.To], txnDetails.Input)
			if decodedData[0].Key == "decodeFailed" {
				var decodeFail DecodeFailure
				decodeFail.Label = decodedData[0].Value
				decodeFail.Type = "red"
				txnDetails.DecodeFailed = decodeFail
			} else {
				txnDetails.DecodedInputs = decodedData
			}
			if functionDetails != "" {
				txnDetails.FunctionDetails = functionDetails
				decoded = true
			}
		} else if txnDetails.TransactionType == "Hash Only" {
			var decodeFail DecodeFailure
			decodeFail.Label = "Hash Only Transaction"
			decodeFail.Type = "red"
			txnDetails.DecodeFailed = decodeFail
			decoded = true
		} else if abiMap[txnDetails.To] == "" {
			if txnDetails.Input == "0x" && txnDetails.Value != 0 {
				var decodeFail DecodeFailure
				decodeFail.Label = "Ether Transfer"
				decodeFail.Type = "yellow"
				txnDetails.DecodeFailed = decodeFail
				decoded = true
			} else {
				var decodeFail DecodeFailure
				decodeFail.Label = "Decode in Progress"
				decodeFail.Type = "yellow"
				txnDetails.DecodeFailed = decodeFail
			}
		} else if abiMap[txnDetails.To] == "missing" {
			var decodeFail DecodeFailure
			decodeFail.Label = "ABI Missing"
			decodeFail.Type = "red"
			txnDetails.DecodeFailed = decodeFail
		}
	}

	if decoded {
		txnMap[txnDetails.TransactionHash] = *txnDetails
	}
}

func calculateTimeElapsed(txnDetails *TransactionReceiptResponse, url string) {
	ethClient := client.EthClient{url}
	getTransactionReceipt := ethClient.GetTransactionReceipt(txnDetails.TransactionHash)
	blockResponseClient := ethClient.GetBlockByNumber(getTransactionReceipt.BlockNumber)
	currentTime := time.Now().Unix()
	creationTime := util.HexStringtoInt64(blockResponseClient.Timestamp)
	creationTimeUnix := creationTime / 1000000000
	elapsedTime := currentTime - creationTimeUnix
	txnDetails.TimeElapsed = elapsedTime
}

func (nsi *NodeServiceImpl) joinRequestResponse(enode string, status string) SuccessResponse {
	var successResponse SuccessResponse
	peerMap[enode] = status
	var enodeString []string
	var ipString []string

	//@TODO: Use regex grouping to extract parts
	enodeVal := strings.TrimPrefix(enode, "enode://")
	enodeString = strings.Split(enodeVal, "@")
	ipString = strings.Split(enodeString[1], ":")
	ip := ipString[0]
	successResponse.Status = fmt.Sprintf("Successfully updated status of %s node with IP: %s to %s", nameMap[enode], ip, status)
	return successResponse
}

func (nsi *NodeServiceImpl) deployContract(pubKeys []string, fileName []string, private bool, url string) []contractcompiler.ContractJson {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	fromAddress := ethClient.Coinbase()


	contractJsons := make([]contractcompiler.ContractJson, 0)

	nms := contractclient.NetworkMapContractClient{client.EthClient{url}, contracthandler.ContractParam{fromAddress, env.GetSetupConf().ContractAdd, "", nil}}
	if private == true && pubKeys[0] == "" {
		enode := ethClient.AdminNodeInfo().ID
		peerNo := len(nms.GetNodeDetailsList())
		publicKeys := make([]string, peerNo-1)
		for i := 0; i < peerNo; i++ {
			if enode != nms.GetNodeDetails(i).Enode {
				publicKeys[i-1] = nms.GetNodeDetails(i).PublicKey
			}
		}
		//pubKeys = []string{"R1fOFUfzBbSVaXEYecrlo9rENW0dam0kmaA2pasGM14=", "Er5J8G+jXQA9O2eu7YdhkraYM+j+O5ArnMSZ24PpLQY="}
		pubKeys = publicKeys
	}

	for _, file := range fileName {
		contractJsons = append(contractJsons, contractcompiler.Compile(file, pubKeys, ethClient, private, contNameMap, contTimeMap, abiMap)...)
	}
	return contractJsons
}

func (nsi *NodeServiceImpl) createNetworkScriptCall(nodename string, currentIP string, rpcPort string, whisperPort string, tesseraPort string, raftPort string, nodeManagerPort string) SuccessResponse {
	var successResponse SuccessResponse
	cmd := exec.Command("./setup.sh", "1", nodename)
	cmd.Dir = "./Setup"
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}

	var setupConf string
	setupConf = "CURRENT_IP=" + currentIP + "\n" + "RPC_PORT=" + rpcPort + "\n" + "WHISPER_PORT=" + whisperPort + "\n" + "TESSERA_PORT=" + tesseraPort + "\n" + "RAFT_PORT=" + raftPort + "\n" + "NODEMANAGER_PORT=" + nodeManagerPort + "\n"
	setupConfByte := []byte(setupConf)
	err = ioutil.WriteFile("./Setup/"+nodename+"/setup.conf", setupConfByte, 0775)
	if err != nil {
		fmt.Println(err)
	}
	successResponse.Status = "success"
	return successResponse
}

func (nsi *NodeServiceImpl) joinRequestResponseCall(nodename string, currentIP string, rpcPort string, whisperPort string, tesseraPort string, raftPort string, nodeManagerPort string, masterNodeManagerPort string, masterIP string) SuccessResponse {
	var successResponse SuccessResponse
	cmd := exec.Command("./setup.sh", "2", nodename, masterIP, masterNodeManagerPort, currentIP, rpcPort, whisperPort, tesseraPort, raftPort, nodeManagerPort)
	cmd.Dir = "./Setup"
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}

	var setupConf string
	setupConf = "CURRENT_IP=" + currentIP + "\n" + "RPC_PORT=" + rpcPort + "\n" + "WHISPER_PORT=" + whisperPort + "\n" + "TESSERA_PORT=" + tesseraPort + "\n" + "RAFT_PORT=" + raftPort + "\n" + "THIS_NODEMANAGER_PORT=" + nodeManagerPort + "\n" + "MASTER_IP=" + masterIP + "\n" + "NODEMANAGER_PORT=" + masterNodeManagerPort + "\n"
	setupConfByte := []byte(setupConf)
	err = ioutil.WriteFile("./Setup/"+nodename+"/setup.conf", setupConfByte, 0775)
	if err != nil {
		fmt.Println(err)
	}
	successResponse.Status = "success"
	return successResponse
}

func (nsi *NodeServiceImpl) resetCurrentNode() SuccessResponse {
	var successResponse SuccessResponse
	cmd := exec.Command("./reset_chain.sh")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		successResponse.Status = "failure"
		return successResponse
	}
	successResponse.Status = "success"
	return successResponse
}

func (nsi *NodeServiceImpl) restartCurrentNode() SuccessResponse {
	var successResponse SuccessResponse
	r, _ := regexp.Compile("[s][t][a][r][t][_][A-Za-z0-9]*[.][s][h]")
	files, err := ioutil.ReadDir(env.GetAppConfig().NodeDir)
	if err != nil {
		log.Println(err)
	}
	var filename string
	for _, f := range files {
		match, _ := regexp.MatchString("[s][t][a][r][t][_][A-Za-z0-9]*[.][s][h]", f.Name())
		if match {
			filename = r.FindString(f.Name())
		}
	}
	filepath := fmt.Sprint("./", filename)
	cmd := exec.Command(filepath)
	cmd.Dir = env.GetAppConfig().NodeDir
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Start()
	if err != nil {
		log.Println(err)
		successResponse.Status = "failure"
		return successResponse
	}
	successResponse.Status = "success"
	return successResponse
}

func (nsi *NodeServiceImpl) latestBlockDetails(url string) LatestBlockResponse {
	var latestBlockResponse LatestBlockResponse
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	currentTime := time.Now().Unix()
	blockNumber := ethClient.BlockNumber()
	blockResponseClient := ethClient.GetBlockByNumber(blockNumber)
	blockNumberInt := util.HexStringtoInt64(blockNumber)
	creationTime := blockResponseClient.Timestamp
	creationTimeInt := util.HexStringtoInt64(creationTime)
	creationTimeUnix := creationTimeInt / 1000000000
	elapsedTime := currentTime - creationTimeUnix
	latestBlockResponse.LatestBlockNumber = blockNumberInt
	latestBlockResponse.TimeElapsed = elapsedTime
	return latestBlockResponse
}

//func (nsi *NodeServiceImpl) latency(url string) ([]LatencyResponse) {
//	var nodeUrl = url
//	ethClient := client.EthClient{nodeUrl}
//	otherPeersResponse := ethClient.AdminPeers()
//	peerCount := len(otherPeersResponse)
//	latencyResponse := make([]LatencyResponse, peerCount+1)
//	for i := 0; i < peerCount+1; i++ {
//		var latOut bytes.Buffer
//		var ip string
//		if i == peerCount {
//			ip = "localhost"
//			thisAdminInfo := ethClient.AdminNodeInfo()
//			latencyResponse[i].EnodeID = thisAdminInfo.ID
//		} else {
//			ip = otherPeersResponse[i].Network.LocalAddress
//			ipString := strings.Split(ip, ":")
//			ip = ipString[0]
//			latencyResponse[i].EnodeID = otherPeersResponse[i].ID
//		}
//		command := fmt.Sprint("ping -c 4 ", ip, " |  awk -F'/' '{ print $5 }' | tail -1")
//		cmd := exec.Command("bash", "-c", command)
//		cmd.Stdout = &latOut
//		err := cmd.Run()
//		if err != nil {
//			fmt.Println(err)
//		}
//		latencyString := strings.TrimSuffix(latOut.String(), "\n")
//		latency, err := strconv.ParseFloat(latencyString, 10)
//		latency = latency * 1000
//		latencyStr := strconv.FormatFloat(latency, 'f', 0, 64)
//		latencyResponse[i].Latency = latencyStr
//	}
//	return latencyResponse
//}

func (nsi *NodeServiceImpl) latency(url string) []LatencyResponse {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	fromAddress := ethClient.Coinbase()

	nms := contractclient.NetworkMapContractClient{client.EthClient{url}, contracthandler.ContractParam{fromAddress, env.GetSetupConf().ContractAdd, "", nil}}

	peerNo := len(nms.GetNodeDetailsList())

	latencyResponse := make([]LatencyResponse, peerNo)
	for i := 0; i < peerNo; i++ {
		var latOut bytes.Buffer
		ip := nms.GetNodeDetails(i).IP
		latencyResponse[i].EnodeID = nms.GetNodeDetails(i).Enode
		command := fmt.Sprint("ping -c 4 ", ip, " |  awk -F'/' '{ print $5 }' | tail -1")
		cmd := exec.Command("bash", "-c", command)
		cmd.Stdout = &latOut
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
		latencyString := strings.TrimSuffix(latOut.String(), "\n")
		latency, err := strconv.ParseFloat(latencyString, 10)
		latency = latency * 1000
		latencyStr := strconv.FormatFloat(latency, 'f', 0, 64)
		latencyResponse[i].Latency = latencyStr
	}
	return latencyResponse
}

func (nsi *NodeServiceImpl) transactionSearchDetails(txno string, url string) BlockDetailsResponse {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	txGetClient := ethClient.GetTransactionReceipt(txno)
	blockNumber := util.HexStringtoInt64(txGetClient.BlockNumber)
	blockDetailsResponse := nsi.getBlockInfo(blockNumber, url)
	return blockDetailsResponse
}

func (nsi *NodeServiceImpl) emailServerConfig(host string, port string, username string, password string, recipientList string, url string) SuccessResponse {
	var successResponse SuccessResponse

	mailServerConfig.Host = host
	mailServerConfig.Port = port
	mailServerConfig.Username = username
	mailServerConfig.Password = password
	mailServerConfig.RecipientList = recipientList

	registered := fmt.Sprint("RECIPIENTLIST=", recipientList, "\n")
	util.AppendStringToFile(env.GetAppConfig().HomeDir+"/setup.conf", registered)

	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			//fmt.Println("Healthcheck done at: ", t)
			if warning > 0 {
				//fmt.Println("Ticker stopped")
				ticker.Stop()
			}
			nsi.healthCheck(url)

		}
	}()
	go func() {
		nsi.sendTestMail()
	}()
	successResponse.Status = "success"
	return successResponse
}

func (nsi *NodeServiceImpl) healthCheck(url string) {
	ethClient := client.EthClient{url}
	blockNumber := ethClient.BlockNumber()
	if blockNumber == "" {
		if warning > 0 {
			exists := util.PropertyExists("RECIPIENTLIST", env.GetAppConfig().HomeDir+"/setup.conf")
			if exists != "" {
				p := properties.MustLoadFile(env.GetAppConfig().HomeDir+"/setup.conf", properties.UTF8)
				recipientList := util.MustGetString("RECIPIENTLIST", p)
				recipients := strings.Split(recipientList, ",")

				b, err := ioutil.ReadFile(env.GetAppConfig().RootDir +"/NodeUnavailableTemplate.txt")

				if err != nil {
					log.Println(err)
				}

				mailCont := string(b)
				mailCont = strings.Replace(mailCont, "\n", "", -1)
				for i := 0; i < len(recipients); i++ {
					nsi.sendMail(mailServerConfig.Host, mailServerConfig.Port, mailServerConfig.Username, mailServerConfig.Password, "Node is not responding", mailCont, recipients[i])
				}
			}
		}
		warning++
	}
}

func (nsi *NodeServiceImpl) sendTestMail() {
	existsA := util.PropertyExists("RECIPIENTLIST", env.GetAppConfig().HomeDir+"/setup.conf")
	existsB := util.PropertyExists("NODENAME", env.GetAppConfig().HomeDir+"/setup.conf")

	if existsA != "" && existsB != "" {
		p := properties.MustLoadFile(env.GetAppConfig().HomeDir+"/setup.conf", properties.UTF8)
		nodename := util.MustGetString("NODENAME", p)
		recipientList := util.MustGetString("RECIPIENTLIST", p)
		recipients := strings.Split(recipientList, ",")
		b, err := ioutil.ReadFile(env.GetAppConfig().RootDir + "/TestMailTemplate.txt")
		if err != nil {
			log.Println(err)
		}

		mailCont := string(b)
		message := fmt.Sprintf(mailCont, nodename)
		for i := 0; i < len(recipients); i++ {
			nsi.sendMail(mailServerConfig.Host, mailServerConfig.Port, mailServerConfig.Username, mailServerConfig.Password, "Quorum Maker Notification Service configured", message, recipients[i])
		}
	}
}

func (nsi *NodeServiceImpl) sendMail(host string, port string, username string, password string, subject string, mailContent string, to string) {
	portNo, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	m := gomail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text", mailContent)

	d := gomail.NewDialer(host, int(portNo), username, password)

	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
	}
}

//@TODO: Implement logrotate command to do this.
func (nsi *NodeServiceImpl) LogRotaterGeth() {
	command := "cat $(ls | grep log | grep -v _) > Geth_$(date| sed -e 's/ /_/g')"

	command1 := "echo -en '' > $(ls | grep log | grep -v _)"

	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = env.GetAppConfig().NodeDir + "/qdata/gethLogs"
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	cmd1 := exec.Command("bash", "-c", command1)
	cmd1.Dir = env.GetAppConfig().NodeDir + "/qdata/gethLogs"
	err1 := cmd1.Run()
	if err1 != nil {
		fmt.Println(err)
	}

}

func (nsi *NodeServiceImpl) LogRotaterConst() {

	command := "cat $(ls | grep log | grep _) > Tessera_$(date| sed -e 's/ /_/g')"

	command1 := "echo -en '' > $(ls | grep log | grep _)"

	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = env.GetAppConfig().NodeDir + "/qdata/tesseraLogs"
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	cmd1 := exec.Command("bash", "-c", command1)
	cmd1.Dir = env.GetAppConfig().NodeDir + "/qdata/tesseraLogs"
	err1 := cmd1.Run()
	if err1 != nil {
		fmt.Println(err)
	}

}

func (nsi *NodeServiceImpl) RegisterNodeDetails(url string) {
	mode := currentMode()
	if mode == "PASSIVE" || mode == "ACTIVENI" {
		return
	}
	var nodeUrl = url

	if env.GetSetupConf().Registered != "TRUE" {
		ethClient := client.EthClient{nodeUrl}

		enode := ethClient.AdminNodeInfo().ID
		fromAddress := ethClient.Coinbase()

		registered := fmt.Sprint("REGISTERED=TRUE", "\n")
		util.AppendStringToFile(env.GetAppConfig().HomeDir+"/setup.conf", registered)
		util.DeleteProperty("REGISTERED=", env.GetAppConfig().HomeDir+"/setup.conf")
		util.DeleteProperty("ROLE=Unassigned", env.GetAppConfig().HomeDir+"/setup.conf")
		nms := contractclient.NetworkMapContractClient{client.EthClient{url}, contracthandler.ContractParam{fromAddress, env.GetSetupConf().ContractAdd, "", nil}}
		nms.RegisterNode(env.GetSetupConf().Nodename, env.GetSetupConf().Role, env.GetSetupConf().PubKey, enode, env.GetSetupConf().CurrentIp, env.GetSetupConf().RaftId)
	}
}

func (nsi *NodeServiceImpl) NetworkManagerContractDeployer(url string) {
	mode := currentMode()
	if mode == "PASSIVE" || mode == "ACTIVENI" {
		return
	}

	if env.GetSetupConf().ContractAdd == "" {
		log.Info("Deploying Network Manager Contract")
		filename := []string{"NetworkManagerContract.sol"}
		deployedContract := nsi.deployContract(nil, filename, false, url)
		contAdd := deployedContract[0].ContractAddress
		contAddAppend := fmt.Sprint("CONTRACT_ADD=", contAdd, "\n")
		util.AppendStringToFile(env.GetAppConfig().HomeDir+"/setup.conf", contAddAppend)
		util.DeleteProperty("CONTRACT_ADD=", env.GetAppConfig().HomeDir+"/setup.conf")

		env.GetSetupConf(true)
	}
}

func ConvertToReadable(p client.TransactionDetailsResponse, pending bool, hash bool) TransactionDetailsResponse {
	var readableTransactionDetailsResponse TransactionDetailsResponse

	readableTransactionDetailsResponse.BlockNumber = util.HexStringtoInt64(p.BlockNumber)
	readableTransactionDetailsResponse.Gas = util.HexStringtoInt64(p.Gas)
	readableTransactionDetailsResponse.GasPrice = util.HexStringtoInt64(p.GasPrice)
	readableTransactionDetailsResponse.TransactionIndex = util.HexStringtoInt64(p.TransactionIndex)
	readableTransactionDetailsResponse.Value = util.HexStringtoInt64(p.Value)
	readableTransactionDetailsResponse.Nonce = util.HexStringtoInt64(p.Nonce)
	readableTransactionDetailsResponse.BlockHash = p.BlockHash
	readableTransactionDetailsResponse.From = p.From
	readableTransactionDetailsResponse.Hash = p.Hash
	readableTransactionDetailsResponse.Input = p.Input
	readableTransactionDetailsResponse.To = p.To
	readableTransactionDetailsResponse.V = p.V
	readableTransactionDetailsResponse.R = p.R
	readableTransactionDetailsResponse.S = p.S
	if util.HexStringtoInt64(p.V) == 37 || util.HexStringtoInt64(p.V) == 38 {
		if pending {
			readableTransactionDetailsResponse.TransactionType = "Private or Hash Only"
		} else if hash {
			readableTransactionDetailsResponse.TransactionType = "Hash Only"
		} else {
			readableTransactionDetailsResponse.TransactionType = "Private"
		}

	} else {
		readableTransactionDetailsResponse.TransactionType = "Public"
	}

	return readableTransactionDetailsResponse
}

func (nsi *NodeServiceImpl) CheckGethStatus(url string) bool {
	ethClient := client.EthClient{url}
	var coinbase string
	for coinbase == "" {
		time.Sleep(1 * time.Second)
		coinbase = ethClient.Coinbase()
	}
	return true
}

func (nsi *NodeServiceImpl) GetChartData(url string) []ChartInfo {
	ethClient := client.EthClient{url}
	chartResponse := make([]ChartInfo, chartSize)
	currentTimeRaw := time.Now().Unix()
	currentTime := currentTimeRaw - (currentTimeRaw % 60)
	startTime := currentTime
	currentBlockNumber := util.HexStringtoInt64(ethClient.BlockNumber())
	bucketTime := currentTime - 60
	stopTime := currentTime - int64(60*chartSize)
	i := 0
	lastBlockNoHex := strconv.FormatInt(currentBlockNumber, 16)
	lastBNoHex := fmt.Sprint("0x", lastBlockNoHex)
	blockResponseClient := ethClient.GetBlockByNumber(lastBNoHex)
	lastCreationTimeRaw := util.HexStringtoInt64(blockResponseClient.Timestamp)
	lastCreationTime := lastCreationTimeRaw / 1000000000
	lastCreationTimeSec := lastCreationTime - (lastCreationTime % 60)
	if lastCreationTimeSec > stopTime {
		for currentTime > stopTime {
			blockCount := 0
			txnCount := 0
			for currentTime > bucketTime {
				blockNoHex := strconv.FormatInt(currentBlockNumber, 16)
				bNoHex := fmt.Sprint("0x", blockNoHex)
				blockResponseClient := ethClient.GetBlockByNumber(bNoHex)
				creationTimeRaw := util.HexStringtoInt64(blockResponseClient.Timestamp)
				creationTimeRaw = creationTimeRaw / 1000000000
				currentTime = creationTimeRaw - (creationTimeRaw % 60)
				if currentTime > bucketTime {
					currentBlockNumber = currentBlockNumber - 1
					txnCount = txnCount + len(blockResponseClient.Transactions)
					blockCount++
				}
			}
			chartResponse[i].BlockCount = blockCount
			chartResponse[i].TransactionCount = txnCount
			chartResponse[i].TimeStamp = (int(bucketTime) + 60) * 1000
			bucketTime = bucketTime - 60
			i++
		}
	}
	for i := 0; i < chartSize; i++ {

		if chartResponse[i].TimeStamp == 0 {
			chartResponse[i].TimeStamp = int(startTime) * 1000
		}
		startTime = startTime - 60
	}
	for i := 0; i < len(chartResponse)/2; i++ {
		j := len(chartResponse) - i - 1
		chartResponse[i], chartResponse[j] = chartResponse[j], chartResponse[i]
	}

	return chartResponse
}

func (nsi *NodeServiceImpl) ContractCrawler(url string) {
	ticker := time.NewTicker(15 * time.Second)
	go func() {
		for range ticker.C {
			if contractCrawlerMutex == 0 {
				nsi.getContracts(url)
			}
		}
	}()
}

func (nsi *NodeServiceImpl) getContracts(url string) {
	contractCrawlerMutex = 1
	ethClient := client.EthClient{url}
	blockNumber := int(util.HexStringtoInt64(ethClient.BlockNumber()))
	for i := lastCrawledBlock + 1; i <= blockNumber; i++ {
		blockNoHex := strconv.FormatInt(int64(i), 16)
		bNoHex := fmt.Sprint("0x", blockNoHex)
		blockResponseClient := ethClient.GetBlockByNumber(bNoHex)
		for _, clientTransactions := range blockResponseClient.Transactions {
			txGetClient := ethClient.GetTransactionReceipt(clientTransactions.Hash)
			if txGetClient.ContractAddress != "" {
				if abiMap[txGetClient.ContractAddress] == "" {
					abiMap[txGetClient.ContractAddress] = "missing"
				}

				if util.HexStringtoInt64(clientTransactions.V) == 37 || util.HexStringtoInt64(clientTransactions.V) == 38 {
					private := ethClient.GetQuorumPayload(clientTransactions.Input)
					if private == "0x" {
						contTypeMap[txGetClient.ContractAddress] = "Hash Only"
					} else {
						contTypeMap[txGetClient.ContractAddress] = "Private"
					}
				} else {
					contTypeMap[txGetClient.ContractAddress] = "Public"
					mode := currentMode()
					if mode == "ACTIVENI" {
						nsi.attachModeRegisterDetails(url, txGetClient.ContractAddress)
					}
				}
				contSenderMap[txGetClient.ContractAddress] = clientTransactions.From
				contTimeMap[txGetClient.ContractAddress] = strconv.Itoa(int(util.HexStringtoInt64(blockResponseClient.Timestamp) / 1000000000))
			}
		}
	}
	mode := currentMode()
	if mode == "ACTIVENI" {
		util.DeleteProperty("MODE=ACTIVENI", env.GetAppConfig().HomeDir+"/setup.conf")
		modeActive := fmt.Sprint("MODE=ACTIVE\n")
		util.AppendStringToFile(env.GetAppConfig().HomeDir+"/setup.conf", modeActive)
		nsi.NetworkManagerContractDeployer(url)
		nsi.RegisterNodeDetails(url)
	}
	lastCrawledBlock = blockNumber
	contractCrawlerMutex = 0
}

func (nsi *NodeServiceImpl) attachModeRegisterDetails(url string, contractAdd string) {
	nmcBytecode, err := ioutil.ReadFile(env.GetAppConfig().RootDir + "/nmcBytecode")
	if err != nil {
		log.Println(err)
	}
	nmcBytecodeString := string(nmcBytecode)
	nmcBytecodeString = strings.Replace(nmcBytecodeString, "\n", "", -1)
	ethClient := client.EthClient{url}
	bytecode := ethClient.GetCode(contractAdd)
	hashIndex := len(bytecode) - 68
	bytecode = bytecode[:hashIndex]
	if bytecode == nmcBytecodeString {
		util.DeleteProperty("MODE=ACTIVENI", env.GetAppConfig().HomeDir+"/setup.conf")
		modeActive := fmt.Sprint("MODE=ACTIVE\n")
		util.AppendStringToFile(env.GetAppConfig().HomeDir+"/setup.conf", modeActive)
		contAddAppend := fmt.Sprint("CONTRACT_ADD=", contractAdd, "\n")
		util.AppendStringToFile(env.GetAppConfig().HomeDir+"/setup.conf", contAddAppend)
		util.DeleteProperty("CONTRACT_ADD=", env.GetAppConfig().HomeDir+"/setup.conf")
		nsi.RegisterNodeDetails(url)
	}
}

func (nsi *NodeServiceImpl) ContractList() []ContractTableRow {
	contractList := make([]ContractTableRow, len(abiMap))
	i := 0
	for key := range abiMap {
		contractList[i].ContractAdd = key
		contractList[i].ABIContent = abiMap[key]
		if abiMap[key] == "missing" {
			contractList[i].ABIContent = ""
		}
		contractList[i].ContractName = contNameMap[key]
		contractList[i].ContractType = contTypeMap[key]
		contractList[i].Sender = contSenderMap[key]
		contractList[i].Timestamp = contTimeMap[key]
		contractList[i].Description = contDescriptionMap[key]
		i++
	}

	return contractList
}

func (nsi *NodeServiceImpl) ContractCount() ContractCounter {
	availableABIs := 0
	totalContracts := 0
	for key := range abiMap {
		if abiMap[key] != "missing" {
			availableABIs++
		}
		totalContracts++
	}
	var contractCount ContractCounter
	contractCount.TotalContracts = totalContracts
	contractCount.ABIavailable = availableABIs
	return contractCount
}

func (nsi *NodeServiceImpl) updateContractDetails(contractAddress string, contractName string, abi string, description string) SuccessResponse {
	var successResponse SuccessResponse
	contNameMap[contractAddress] = contractName
	abiMap[contractAddress] = abi
	contDescriptionMap[contractAddress] = description
	successResponse.Status = "Successfully updated contract details"
	return successResponse
}

func (nsi *NodeServiceImpl) returnCurrentInitializationState() SuccessResponseBool {
	var successResponse SuccessResponseBool
	state := currentState()
	if state == "I" {
		successResponse.Status = true
	}
	return successResponse
}

func currentMode() string {
	return env.GetSetupConf().Mode
}

func currentState() string {

	return env.GetSetupConf().State
}

func (nsi *NodeServiceImpl) ABICrawler(url string) {
	updateLastCheckedTime("0")
	ticker := time.NewTicker(15 * time.Second)
	go func() {
		for range ticker.C {
			if abiCrawlerMutex == 0 {
				nsi.DirectoryCrawl()
			}
		}
	}()
}

func (nsi *NodeServiceImpl) DirectoryCrawl() {
	abiCrawlerMutex = 1
	ABIList := getFilesFromDirectory(env.GetAppConfig().ContractsDir)
	nsi.populateABIMap(ABIList)
	abiCrawlerMutex = 0
	updateLastCheckedTime(strconv.Itoa(int(time.Now().Unix())))
}

func getFilesFromDirectory(searchDir string) []CrawledABI {
	fileList := []string{}
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	if err != nil {
		fmt.Println()
	}
	for _, file := range fileList {
		crawledABIs = append(crawledABIs, getABIsFromDirectory(file)...)
	}
	return crawledABIs
}

func getABIsFromDirectory(searchDir string) []CrawledABI {
	r := regexp.MustCompile(`.json$`)
	var crawledABIs []CrawledABI
	files, _ := ioutil.ReadDir(searchDir)
	for _, file := range files {
		var crawledABI CrawledABI
		if r.MatchString(file.Name()) && !file.IsDir() {
			crawledABI.Filename = searchDir + "/" + file.Name()
			crawledABI.Contractname = strings.Split(file.Name(), ".")[0]
			crawledABI.ModificationTime = file.ModTime().Unix()
			if crawledABI.ModificationTime < getLastCheckedTime() {
				continue
			}
			crawledABIs = append(crawledABIs, crawledABI)
		}
	}
	return crawledABIs
}

func (nsi *NodeServiceImpl) populateABIMap(jsons []CrawledABI) {
	for _, file := range jsons {
		if !file.Processed {
			nsi.parseABIJson(file)
			file.Processed = true
		}
	}
}

func (nsi *NodeServiceImpl) parseABIJson(file CrawledABI) {
	var contractJSONContent contractJSONTruffle
	fileBytes, err := ioutil.ReadFile(file.Filename)
	if err != nil {
		log.Println(err)
	}

	jsonContent := string(fileBytes)
	jsonContent = strings.Replace(jsonContent, "\n", "", -1)
	json.Unmarshal([]byte(jsonContent), &contractJSONContent)
	abiContent, _ := json.Marshal(contractJSONContent.Abi)
	abiString := make([]string, len(abiContent))
	for i := 0; i < len(abiContent); i++ {
		abiString[i] = string(abiContent[i])
	}
	abiData := fmt.Sprint(strings.Join(abiString, ""))

	interfaceContent, _ := json.Marshal(contractJSONContent.Interface)
	interfaceString := make([]string, len(interfaceContent))
	for i := 0; i < len(interfaceContent); i++ {
		interfaceString[i] = string(interfaceContent[i])
	}
	interfaceData := fmt.Sprint(strings.Join(interfaceString, ""))
	bytecodeData := contractJSONContent.Bytecode
	contractName := contractJSONContent.ContractName
	contractAddress := contractJSONContent.Address

	var data string
	if len(abiData) != 4 {
		data = abiData
	} else if len(interfaceData) != 4 {
		data = interfaceData
	} else {
		data = jsonContent
		data = strings.Replace(data, "\n", "", -1)
	}

	if contractAddress != "" && contractName != "" {
		nsi.writeContractDetailsToDisk(data, bytecodeData, contractAddress, contractName)
		nsi.updateContractDetails(contractAddress, contractName, data, "default")
	} else if contractAddress != "" && contractName == "" {
		contNameMap[contractAddress] = file.Contractname
		abiMap[contractAddress] = data
		contDescriptionMap[contractAddress] = "default"
	}
}

func (nsi *NodeServiceImpl) writeContractDetailsToDisk(data string, bytecodeData string, contractAddress string, contractName string) {
	jsonString := util.ComposeJSON(data, bytecodeData, contractAddress)

	path := env.GetAppConfig().ContractsDir + "/" + contractAddress + "_" + contractName

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0775)
	}

	filePath := path + "/" + contractName + ".json"
	jsByte := []byte(jsonString)
	err := ioutil.WriteFile(filePath, jsByte, 0775)
	if err != nil {
		fmt.Println(err)
	}
}

func getLastCheckedTime() int64 {
	fileBytes, err := ioutil.ReadFile(env.GetAppConfig().ContractsDir + "/.lastCheckedTime")
	if err != nil {
		log.Println(err)
	}
	jsonContent := string(fileBytes)
	jsonContent = strings.Replace(jsonContent, "\n", "", -1)
	lastChecked, err := strconv.Atoi(jsonContent)
	return int64(lastChecked)
}

func updateLastCheckedTime(timeVal string) {
	util.DeleteFile(env.GetAppConfig().ContractsDir + "/.lastCheckedTime")
	util.CreateFile(env.GetAppConfig().ContractsDir + "/.lastCheckedTime")
	util.WriteFile(env.GetAppConfig().ContractsDir + "/.lastCheckedTime", timeVal)
}

func (nsi *NodeServiceImpl) createAccount(password string, url string) SuccessResponse {
	var accountDetail SuccessResponse
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	accountAddress := ethClient.CreateAccount(password)
	accountDetail.Status = fmt.Sprint("Account ", accountAddress, " has been created successfully")
	return accountDetail
}

func (nsi *NodeServiceImpl) getAccounts(url string) []ethAccount {
	var accountList []ethAccount
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	coinbase := ethClient.Coinbase()
	accounts := ethClient.GetAccounts()
	for _, accountID := range accounts {
		var account ethAccount
		account.AccountAddress = accountID
		if accountID == coinbase {
			account.Coinbase = true
		}
		account.Balance = util.HexStringtoLargeInt64(ethClient.GetBalance(accountID))
		accountList = append(accountList, account)
	}
	return accountList
}

func (nsi *NodeServiceImpl) getNodeIPs(url string) []connectedIP {
	var nodeUrl = url
	var ipList []connectedIP
	var connectedIPs = map[string]int{}
	ethClient := client.EthClient{nodeUrl}
	fromAddress := ethClient.Coinbase()
	enode := ethClient.AdminNodeInfo().ID

	nms := contractclient.NetworkMapContractClient{client.EthClient{url}, contracthandler.ContractParam{fromAddress, env.GetSetupConf().ContractAdd, "", nil}}
	nodeList := nms.GetNodeDetailsList()
	for _, node := range nodeList {
		if node.Enode != enode {
			count := connectedIPs[node.IP]
			connectedIPs[node.IP] = count + 1
		}
	}
	for k := range connectedIPs {
		var connected connectedIP
		connected.IP = k
		connected.Count = connectedIPs[k]
		ipList = append(ipList, connected)
	}
	return ipList
}

func (nsi *NodeServiceImpl) updateWhitelist(ipList []string) SuccessResponse {
	var update SuccessResponse
	util.DeleteFile(env.GetAppConfig().ContractsDir + "/.whiteList")
	util.CreateFile(env.GetAppConfig().ContractsDir + "/.whiteList")
	for _, ip := range ipList {
		util.AppendStringToFile(env.GetAppConfig().ContractsDir + "/.whiteList", fmt.Sprint(ip, "\n"))
	}
	update.Status = "IP Whitelist has been updated successfully"
	return update
}
