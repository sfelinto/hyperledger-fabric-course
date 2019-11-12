package main

/**
 * Demonstrates the use of CID
 **/
import (
	// For printing messages on console
	"fmt"

	// The shim package
	"github.com/hyperledger/fabric/core/chaincode/shim"

	// peer.Response is in the peer package
	"github.com/hyperledger/fabric/protos/peer"

	// Client Identity Library
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"

	// Standard go crypto package
	"crypto/x509"

	"strconv"
)

// CidChaincode Represents our chaincode object
type CidChaincode struct {
}

// Invoke method
func (clientdid *CidChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the function name and parameters
	funcName, args := stub.GetFunctionAndParameters()

	// Just to satisfy the compiler - otherwise it will complain that args declared but not used
	fmt.Println(len(args))

	if funcName == "ReadAttributesOfCaller" {

		// Return JSON
		return clientdid.ReadAttributesOfCaller(stub)

	} else if funcName == "AsssertOnCallersDepartment" {

		// Returns the Trade Rejecton/Approval result
		return clientdid.AssertOnCallersDepartment(stub)

	} else if funcName == "ApproveTrade" {
		// To be coded in an exercise
		return clientdid.ApproveTrade(stub, args)
	}

	return shim.Error("Bad Func Name!!!")
}

// ReadAttributesOfCaller reads the attributes of the callers cert and return it as JSON
func (clientdid *CidChaincode) ReadAttributesOfCaller(stub shim.ChaincodeStubInterface) peer.Response {

	// Variable to hold the result
	jsonResult := "{"

	// 1. Get the unique ID of the user
	id, err := cid.GetID(stub)

	if err != nil {
		fmt.Println("Error GetID() =" + err.Error())
		return shim.Error(err.Error())
	}
	// Format and add the attribute to JSON
	jsonResult += SetJSONNV("id", id)

	// 2. Get the MSP ID of the user
	var mspid string
	mspid, err = cid.GetMSPID(stub)

	if err != nil {
		fmt.Println("Error GetMSPID() =" + err.Error())
		return shim.Error(err.Error())
	}
	// Format and add the attribute to JSON
	jsonResult += "," + SetJSONNV("MSPID", mspid)

	// 3. Get the standard attributes added by default
	// "hf.Affiliation" ,"hf.EnrollmentID", "hf.Type"
	affiliation, _, _ := cid.GetAttributeValue(stub, "hf.Affiliation")
	enrollID, _, _ := cid.GetAttributeValue(stub, "hf.EnrollmentID")
	userType, _, _ := cid.GetAttributeValue(stub, "hf.Type")
	// Format and add the attribute to JSON
	jsonResult += "," + SetJSONNV("affiliation", affiliation)
	jsonResult += "," + SetJSONNV("enrollID", enrollID)
	jsonResult += "," + SetJSONNV("userType", userType)

	// 4. Get the attr value for "app.accounting.role"
	attrValue, flag, _ := cid.GetAttributeValue(stub, "app.accounting.role")
	if !flag {
		attrValue = "NOT SET"
	}
	// Format and add the attribute to JSON
	jsonResult += "," + SetJSONNV("app.accounting.role", attrValue)

	// 5. Get the attr value for "department"
	attrValue, flag, _ = cid.GetAttributeValue(stub, "department")
	if !flag {
		attrValue = "NOT SET"
	}
	// Format and add the attribute to JSON
	jsonResult += "," + SetJSONNV("department", attrValue)

	// 6. Get the Certificate of the caller - not sending it back in JSON
	var cert *x509.Certificate
	cert, err = cid.GetX509Certificate(stub)
	if err != nil {
		fmt.Println("Error GetX509Certificate() =" + err.Error())
		return shim.Error(err.Error())
	}
	fmt.Println("GetX509Certificate() = " + string(cert.RawSubject))

	// Close the JSON and send it as response
	jsonResult += "}"
	return shim.Success([]byte(jsonResult))
}

// AssertOnCallersDepartment uses the cid AsssertAttributeValue
// Rule = Only a caller with department=accounting can invoke this function
func (clientdid *CidChaincode) AssertOnCallersDepartment(stub shim.ChaincodeStubInterface) peer.Response {

	// Get the enrollID and dept
	enrollID, _, _ := cid.GetAttributeValue(stub, "hf.EnrollmentID")
	dept, _, _ := cid.GetAttributeValue(stub, "department")

	// We can use if statement or Assert call to check the rule
	// Check if the department attribute is set to "accounting)"
	err := cid.AssertAttributeValue(stub, "department", "accounting")

	// Check if valid err returned
	if err != nil {
		return shim.Error("Access Denied to " + enrollID + " from " + dept + " !!!")
	}

	// Return success
	return shim.Success([]byte("Access Granted to " + enrollID + " from " + dept))
}

// ApproveTrade - checks the amount of trade passed in args[0] - applies the business rules
// Rule#1  Caller MUST be from accounting dept
// Rule#2 If trade < 100K it can be approved by anyone from accounting
// Rule#3 If trade >= 100K the caller MUST have a role = manager
func (clientdid *CidChaincode) ApproveTrade(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Must provide Trade Value in args[0] !!!")
	}

	enrollID, _, _ := cid.GetAttributeValue(stub, "hf.EnrollmentID")

	// Rule#1  Caller MUST be from accounting dept
	attrValue, flag, _ := cid.GetAttributeValue(stub, "department")
	if !flag || attrValue != "accounting" {
		return shim.Error("REJECTED - Caller MUST be from accounting department to Approve the trade !!!")
	}

	// Convert trade value to number - ignoring error
	tradeValue, _ := strconv.ParseUint(string(args[0]), 10, 64)

	// Rule#2 If trade < 100K it can be approved by anyone from accounting
	if tradeValue < 100000 {

		return clientdid.ProcessTheTrade(stub, args[0], enrollID)
	}

	// Rule#3 If trade >= 100K the caller MUST have a role = manager
	attrValue, flag, _ = cid.GetAttributeValue(stub, "app.accounting.role")

	if !flag || attrValue != "manager" {
		return shim.Error("REJECTED - Caller has role='" + attrValue + "' but since Tradevalue=" + args[0] + " it requires role='manager'")
	}

	// All rules fulfilled
	return clientdid.ProcessTheTrade(stub, args[0], enrollID)
}

// ProcessTheTrade - dummy function - in real sceanrio the state will change for asset in the trade
func (clientdid *CidChaincode) ProcessTheTrade(stub shim.ChaincodeStubInterface, tradeValue, enrollID string) peer.Response {

	// Result string sent to caller as part of the endorsement
	result := "APPROVED - Trade value=" + tradeValue + " by " + enrollID

	return shim.Success([]byte(result))
}

// SetJSONNV returns a name value pair in JSON format
func SetJSONNV(attr, value string) string {
	return " \"" + attr + "\":\"" + value + "\""
}

// Init Implements the Init method
func (clientdid *CidChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Simply print a message
	fmt.Println("Init executed in history")

	// Return success
	return shim.Success(nil)
}

// Chaincode registers with the Shim on startup
func main() {
	fmt.Printf("Started Chaincode. token/cid\n")
	err := shim.Start(new(CidChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
