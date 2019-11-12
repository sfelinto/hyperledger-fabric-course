package main

import (
	"fmt"

	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"github.com/hyperledger/fabric/protos/peer"
)

type TokenChaincode struct {
}

func (token *TokenChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {

	fmt.Println("Init executed")

	return shim.Success(nil)
}

func (token *TokenChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	fmt.Println("Invoke executed")

	fmt.Printf("GetTxID() => %s\n", stub.GetTxID())

	fmt.Printf("GetChannelID() => %s\n", stub.GetChannelID())

	TxTimestamp, _ := stub.GetTxTimestamp()
	timeStr := time.Unix(TxTimestamp.GetSeconds(), 0)
	fmt.Printf("GetTxTimestamp() => %s\n", timeStr)

	return shim.Success(nil)
}

func main() {
	fmt.Println("Started chaincode.")

	err := shim.Start(new(TokenChaincode))

	if err != nil {
		fmt.Printf("Error starting chaincode: %s ", err)
	}
}
