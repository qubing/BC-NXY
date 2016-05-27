package main

import (
	"errors"
	"fmt"
	"encoding/json"
	//"strconv"
	//"reflect"
	//"unsafe"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Chaincode example simple Chaincode implementation
type Chaincode struct {
}

// args[0]	number		票号
func (t *Chaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var err error
	var Number string

	if len(args) != 1 {return nil, errors.New("sign Init Expecting 1 number of arguments.")}

	Number = args[0]

	err = stub.PutState(Number, []byte("custodianPublish"))
	if err != nil {return nil, err}

	return nil, nil
}

// args[0]	number		票号
// args[1]	stepName	票号名称
// args[2]	salerParty
// args[3]	buyinParty
// args[4]	dealAmount


// custodianPublish 	(票据托管)
// custodianAccept	(接受托管)
// onsaleApplication 	(卖方提出申请)
// salerPartyReview 	(卖方复核)
// salerPartyApproval 	(卖方审批)
// buyinPartyCheck 	(买方审核)
// buyinPartyApproval 	(买方审批)
// settlement 		(清算处理，所有权变更)
func (t *Chaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var err error
	var Number string
	var stepName string

	if len(args) != 5 {return nil, errors.New("sign Invoke Expecting 5 number of arguments.")}

	Number = args[0]
	stepName = args[1]

	err = stub.PutState(Number, []byte(stepName))
	if err != nil {return nil, err}

	fmt.Printf("after sign Invoke, stepName = %s \n", stepName)

	if stepName == "onsaleApplication" {
		fmt.Printf("onsaleApplication approval is done, jump into contract invoke. \n")

		var f string
		var err error

		dealAmount := args[4]
		buyinParty := args[3]

		// query in contract
		chaincodename := "contract"
		f = "query"
		queryArgs := []string{Number}

		fmt.Printf("stub.QueryChaincode(chaincodename, f, queryArgs) start \n")
		queryResponse_Byte, err := stub.QueryChaincode(chaincodename, f, queryArgs)
		fmt.Printf("stub.QueryChaincode(chaincodename, f, queryArgs) completed \n")

		queryResponse_String := string(queryResponse_Byte)
		if err != nil {return nil, err}

		//json str 转map
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(queryResponse_String), &data);
		err == nil {
			fmt.Println("===============json str 转 map===============\n")
			fmt.Println(data)
			fmt.Println("\n")
		}

		if dealAmount == "" { return nil, errors.New("===============dealAmount is blank! It is incorrect!===============\n")}
		if buyinParty == "" { return nil, errors.New("===============buyinParty is blank! It is incorrect!===============\n")}

		// invoke in contract, add value of buyinParty
		f = "invoke"

		party := data["party"].(string)
		number := data["number"].(string)
		attribute := data["attribute"].(string)
		billType := data["billType"].(string)
		issuerName := data["issuerName"].(string)
		issuerAccountID := data["issuerAccountID"].(string)
		issuerAccountBankID := data["issuerAccountBankID"].(string)
		custodianName := data["custodianName"].(string)
		custodianAccountID := data["custodianAccountID"].(string)
		custodianAccountBankID := data["custodianAccountBankID"].(string)
		faceAmount := data["faceAmount"].(string)
		acceptorName := data["acceptorName"].(string)
		acceptorAccountID := data["acceptorAccountID"].(string)
		acceptorBankID := data["acceptorBankID"].(string)
		issueDate := data["issueDate"].(string)
		dueDate := data["dueDate"].(string)
		acceptDate := data["acceptDate"].(string)
		payBankID := data["payBankID"].(string)
		transferableFlag := data["transferableFlag"].(string)

		salerParty := data["salerParty"].(string)

		invokeArgs := []string{ party, number ,attribute ,billType ,issuerName ,issuerAccountID ,issuerAccountBankID ,
			custodianName ,custodianAccountID ,custodianAccountBankID ,faceAmount ,acceptorName ,acceptorAccountID ,acceptorBankID ,
			issueDate, dueDate ,acceptDate ,payBankID ,transferableFlag ,salerParty ,buyinParty, dealAmount }

		invokeResponse, err := stub.InvokeChaincode(chaincodename, f, invokeArgs)
		if err != nil {return nil, err}

		fmt.Printf("contract Invoke successfully. Got response %s \n", string(invokeResponse))



	}

	if stepName == "settlement" {

		fmt.Printf("All approval done, jump into contract invoke. \n")

		var f string

		chaincodename := "contract"
		f = "query"
		queryArgs := []string{Number}
		queryResponse_Byte, err := stub.QueryChaincode(chaincodename, f, queryArgs)
		queryResponse_String := string(queryResponse_Byte)
		if err != nil {return nil, err}

		//json str 转map
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(queryResponse_String), &data); err == nil {
			fmt.Println("==============json str 转 map=======================\n")
			fmt.Println(data)
			fmt.Println(data["party"])
		}

		if data["party"].(string) != data["salerParty"].(string) {
			return nil, errors.New("bill's belongings party is not the same as salerParty. It can not on sale.\n")
		}

		// invoke to buyinParty
		chaincodename = "contract"
		f = "invoke"

		//party := data["party"]
		number := data["number"].(string)
		attribute := data["attribute"].(string)
		billType := data["billType"].(string)
		issuerName := data["issuerName"].(string)
		issuerAccountID := data["issuerAccountID"].(string)
		issuerAccountBankID := data["issuerAccountBankID"].(string)
		custodianName := data["custodianName"].(string)
		custodianAccountID := data["custodianAccountID"].(string)
		custodianAccountBankID := data["custodianAccountBankID"].(string)
		faceAmount := data["faceAmount"].(string)
		acceptorName := data["acceptorName"].(string)
		acceptorAccountID := data["acceptorAccountID"].(string)
		acceptorBankID := data["acceptorBankID"].(string)
		issueDate := data["issueDate"].(string)
		dueDate := data["dueDate"].(string)
		acceptDate := data["acceptDate"].(string)
		payBankID := data["payBankID"].(string)
		transferableFlag := data["transferableFlag"].(string)

		salerParty := data["salerParty"].(string)
		buyinParty := data["buyinParty"].(string)
		dealAmount := data["dealAmount"].(string)

		fmt.Printf("contact invoke:: salerParty:%s, buyinParty:%s, dealAmount:%s \n")

		invokeArgs := []string{ buyinParty, number ,attribute ,billType ,issuerName ,issuerAccountID ,issuerAccountBankID ,
			custodianName ,custodianAccountID ,custodianAccountBankID ,faceAmount ,acceptorName ,acceptorAccountID ,acceptorBankID ,
			issueDate, dueDate ,acceptDate ,payBankID ,transferableFlag ,"" ,"" ,dealAmount }

		invokeResponse, err := stub.InvokeChaincode(chaincodename, f, invokeArgs)
		if err != nil {return nil, err}

		fmt.Printf("contract Invoke successfully. Got response %s \n", string(invokeResponse))

		// execute cash transfer
		chaincodename = "cash"
		f = "invoke"
		invokeArgs = []string{ buyinParty, salerParty, dealAmount }
		invokeResponse, err = stub.InvokeChaincode(chaincodename, f, invokeArgs)

		fmt.Printf("cash transfer finished. \n")

	}

	return nil, nil

}

// args[0]	number		票号
func (t *Chaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var err error
	var Number string
	var stepName string

	if len(args) != 1 {return nil, errors.New("sign Query Expecting 1 number of arguments.")}
	Number = args[0]

	stepName_Byte, err := stub.GetState(Number)
	stepName = string(stepName_Byte)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + Number + "\"}"
		return nil, errors.New(jsonResp)
	}
	if stepName == "" {
		jsonResp := "{\"Error\":\"stepName is Nil for " + Number + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Number\":\"" + Number + "\",\"stepName\":\"" + string(stepName) + "\"}"

	fmt.Printf("Query Response:%s\n", jsonResp)
	return []byte(jsonResp), nil

}


func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {fmt.Printf("Error starting sign chaincode: %s", err)}
}