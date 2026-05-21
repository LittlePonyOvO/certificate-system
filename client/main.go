package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

// 跨域处理
func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

// 上传证书处理
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "解析表单失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	certId := r.FormValue("certId")
	name := r.FormValue("name")
	studentId := r.FormValue("studentId")
	major := r.FormValue("major")
	issueTime := r.FormValue("issueTime")
	hash := r.FormValue("hash")

	if certId == "" || name == "" || studentId == "" || major == "" || issueTime == "" || hash == "" {
		http.Error(w, "请填写所有字段", http.StatusBadRequest)
		return
	}

	// 构建 peer chaincode invoke 命令
	cmd := exec.Command("/bin/bash", "-c", `
		cd ~/certificate-system/fabric-workspace/fabric-samples/test-network && \
		export PATH=${PWD}/../bin:$PATH && \
		export FABRIC_CFG_PATH=${PWD}/../config && \
		source scripts/envVar.sh && \
		setGlobals 1 && \
		peer chaincode invoke -o localhost:7050 \
			--ordererTLSHostnameOverride orderer.example.com \
			--tls --cafile "${ORDERER_CA}" \
			-C mychannel -n certcode \
			--peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} \
			--peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} \
			-c '{"function":"uploadCertificate","Args":["`+certId+`","`+name+`","`+studentId+`","`+major+`","`+issueTime+`","`+hash+`"]}'
	`)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("存证失败: %v, output: %s", err, string(output))
		http.Error(w, "存证失败: "+string(output), http.StatusInternalServerError)
		return
	}

	// 检查是否成功
	if strings.Contains(string(output), "Chaincode invoke successful") {
		w.Write([]byte("存证成功！证书ID：" + certId))
	} else {
		http.Error(w, "存证失败: "+string(output), http.StatusInternalServerError)
	}
}

// 查询证书处理
func queryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	certId := r.URL.Query().Get("certId")
	if certId == "" {
		http.Error(w, "请提供certId参数", http.StatusBadRequest)
		return
	}

	// 构建 peer chaincode query 命令
	cmd := exec.Command("/bin/bash", "-c", `
		cd ~/certificate-system/fabric-workspace/fabric-samples/test-network && \
		export PATH=${PWD}/../bin:$PATH && \
		export FABRIC_CFG_PATH=${PWD}/../config && \
		source scripts/envVar.sh && \
		setGlobals 1 && \
		peer chaincode query -C mychannel -n certcode -c '{"function":"queryCertificate","Args":["`+certId+`"]}'
	`)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("查询失败: %v, output: %s", err, string(output))
		http.Error(w, "查询失败: "+string(output), http.StatusInternalServerError)
		return
	}

	// 检查是否包含错误信息
	outputStr := strings.TrimSpace(string(output))
	if strings.Contains(outputStr, "Error:") || strings.Contains(outputStr, "error") {
		http.Error(w, "查询失败: "+outputStr, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(outputStr))
}

// 首页处理
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", cors(uploadHandler))
	http.HandleFunc("/query", cors(queryHandler))

	fmt.Println("🚀 学业证书存证客户端启动")
	fmt.Println("📡 服务地址: http://localhost:8080")
	fmt.Println("🔗 使用前请确保已启动Fabric网络")
	log.Fatal(http.ListenAndServe(":8080", nil))
}