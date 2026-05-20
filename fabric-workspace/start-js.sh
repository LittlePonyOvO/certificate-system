#!/bin/bash

echo "========== 启动学业证书存证系统 (JavaScript链码) =========="

# 进入测试网络目录
cd ~/certificate-system/fabric-workspace/fabric-samples/test-network

# 停止旧网络
echo "0. 停止旧网络..."
./network.sh down

# 启动网络
echo "1. 启动Fabric网络..."
./network.sh up createChannel -c mychannel

# 部署 JavaScript 链码
echo "2. 部署JavaScript链码..."
./network.sh deployCC -c mychannel \
  -ccn certcode \
  -ccp ../../certificate-chaincode \
  -ccl javascript

# 检查部署是否成功
if [ $? -ne 0 ]; then
    echo "❌ 链码部署失败！"
    exit 1
fi

# 设置环境变量
echo "3. 设置环境变量..."
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=${PWD}/../config
source scripts/envVar.sh
setGlobals 1

# 等待链码生效（关键：链码部署后需要时间启动容器）
echo "4. 等待链码容器启动..."
sleep 15

# 校验链码是否已部署成功
echo "5. 校验链码部署状态..."
peer lifecycle chaincode querycommitted --channelID mychannel --name certcode > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "❌ 链码未成功部署到通道！"
    exit 1
fi
echo "✅ 链码已成功部署到通道"

# 测试上传证书
echo "6. 测试上传证书..."
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${ORDERER_CA}" \
  -C mychannel -n certcode \
  --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} \
  --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} \
  -c '{"function":"uploadCertificate","Args":["CERT001","张三","2024001","计算机科学与技术","2025-06-20","abc123456789xyz"]}'

# 检查上传是否成功
if [ $? -ne 0 ]; then
    echo "❌ 证书上传失败！"
    exit 1
fi

sleep 2

# 测试查询证书
echo "7. 测试查询证书..."
peer chaincode query -C mychannel -n certcode -c '{"function":"queryCertificate","Args":["CERT001"]}'

# 检查查询是否成功
if [ $? -ne 0 ]; then
    echo "❌ 证书查询失败！"
    exit 1
fi

echo ""
echo "========== ✅ 学业证书存证系统(JavaScript链码)启动完成！=========="
echo ""
echo "常用命令："
echo "  查询证书: peer chaincode query -C mychannel -n certcode -c '{\"function\":\"queryCertificate\",\"Args\":[\"证书ID\"]}'"
echo "  上传证书: peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile \"\${ORDERER_CA}\" -C mychannel -n certcode --peerAddresses localhost:7051 --tlsRootCertFiles \${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles \${PEER0_ORG2_CA} -c '{\"function\":\"uploadCertificate\",\"Args\":[\"证书ID\",\"姓名\",\"学号\",\"专业\",\"颁发时间\",\"哈希值\"]}'"
echo "  停止网络: cd ~/certificate-system/fabric-workspace/fabric-samples/test-network && ./network.sh down"