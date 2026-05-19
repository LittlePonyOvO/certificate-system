package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type CertificateChaincode struct {
}

type Certificate struct {
	CertId     string `json:"certId"`
	Name       string `json:"name"`
	StudentId  string `json:"studentId"`
	Major      string `json:"major"`
	IssueTime  string `json:"issueTime"`
	Hash       string `json:"hash"`
	UploadTime string `json:"uploadTime"`
}

func (t *CertificateChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("学业证书存证系统初始化成功！")
	return shim.Success(nil)
}

func (t *CertificateChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Printf("调用函数: %s\n", function)

	switch function {
	case "uploadCertificate":
		return t.uploadCertificate(stub, args)
	case "queryCertificate":
		return t.queryCertificate(stub, args)
	default:
		return shim.Error("未知函数: " + function)
	}
}

func (t *CertificateChaincode) uploadCertificate(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 6 {
		return shim.Error("参数错误：需要6个参数（certId, name, studentId, major, issueTime, hash）")
	}

	certId := args[0]
	name := args[1]
	studentId := args[2]
	major := args[3]
	issueTime := args[4]
	hash := args[5]

	certificate := &Certificate{
		CertId:     certId,
		Name:       name,
		StudentId:  studentId,
		Major:      major,
		IssueTime:  issueTime,
		Hash:       hash,
		UploadTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	certificateJSON, err := json.Marshal(certificate)
	if err != nil {
		return shim.Error("序列化证书失败: " + err.Error())
	}

	err = stub.PutState(certId, certificateJSON)
	if err != nil {
		return shim.Error("写入区块链失败: " + err.Error())
	}

	return shim.Success([]byte("证书存证成功！ID：" + certId))
}

func (t *CertificateChaincode) queryCertificate(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("参数错误：需要1个参数（certId）")
	}

	certId := args[0]

	certificateJSON, err := stub.GetState(certId)
	if err != nil {
		return shim.Error("查询证书失败: " + err.Error())
	}
	if certificateJSON == nil {
		return shim.Error("证书ID " + certId + " 不存在")
	}

	return shim.Success(certificateJSON)
}

func main() {
	err := shim.Start(new(CertificateChaincode))
	if err != nil {
		fmt.Printf("启动链码失败: %s", err)
	}
}