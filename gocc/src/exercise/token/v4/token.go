package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"github.com/hyperledger/fabric/protos/peer"
)

const ChaincodeName = "tokenv2"

var logger = shim.NewLogger(ChaincodeName)

func (token *TokenChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {

	//fmt.Println("Init executed")

	logger.Debug("Init executed v2 - DEBUG")

	return shim.Success(nil)
}

func (token *TokenChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	//fmt.Println("Invoke executed")
	logger.Debug("DEBUG=Invoke executed v2")
	logger.Info("INFO=Invoke executed v2")
	logger.Noticef("NOTICE format string  Value=%s", "Notice executed v2")
	logger.Warning("WARNING=Invoke executed v2", " [any number of parameters of different types]", 123)

	return shim.Success(nil)
}

func main() {
	fmt.Println("Started chaincode.")

	err := shim.Start(new(TokenChaincode))

	if err != nil {
		fmt.Printf("Error starting chaincode: %s ", err)
	}
}
