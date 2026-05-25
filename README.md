# 学业证书区块链存证系统

基于 Hyperledger Fabric 2.4 的极简版学业证书存证系统，支持证书上传和查询功能。

## 📁 项目结构

```
certificate-system/
├── fabric-workspace/
│   ├── certificate-chaincode/    # Go链码实现
│   │   ├── main.go               # 核心链码逻辑
│   │   ├── go.mod
│   │   ├── go.sum
│   │   └── vendor/               # Go依赖
│   ├── fabric-samples/           # Fabric官方测试网络
│   ├── client/                   # 极简客户端
│   │   ├── index.html            # 网页前端
│   │   ├── main.go               # Go后端服务
│   │   └── client-server         # 编译后的服务
│   ├── start-go.sh               # Go链码启动脚本
│   ├── start-js.sh               # JavaScript链码启动脚本
│   └── start.sh                  # 旧版启动脚本
└── README.md                     # 部署说明
```

## 🚀 快速开始

### 环境要求

| 依赖 | 版本 | 说明 |
|------|------|------|
| Ubuntu | 20.04/22.04 | 推荐使用 |
| Docker | 20.10+ | 必需 |
| Docker Compose | 1.29+ | 必需 |
| Git | 任意 | 必需 |
| Go | 1.17+ | 链码开发用 |

### 步骤1：安装依赖

```bash
# 更新软件源
sudo apt update && sudo apt upgrade -y

# 安装 Git、curl
sudo apt install git curl -y

# 安装 Docker
sudo apt install docker.io -y

# 启动 Docker 并设置开机自启
sudo systemctl start docker
sudo systemctl enable docker

# 安装 Docker Compose
sudo apt install docker-compose -y

# 给 Docker 加权限（避免每次输 sudo）
sudo usermod -aG docker $USER

# 重启终端生效
```

### 步骤2：下载项目和 Fabric

```bash
# 克隆项目（或复制项目文件）
git clone https://github.com/LittlePonyOvO/certificate-system.git

# 下载 Hyperledger Fabric 2.4
 cd ~/certificate-system/fabric-workspace && curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.4.9 1.5.3
# 等待下载完成（约5-10分钟）
```

### 步骤3：启动区块链网络

```bash
# 进入项目目录
cd ~/certificate-system/fabric-workspace

# 使用 Go 链码启动（推荐）
./start-go.sh

# 或使用 JavaScript 链码启动（需要 Node.js 环境）
# ./start-js.sh
```

**启动成功标志**：
```
========== 学业证书存证系统(Go链码)启动完成！==========
```

### 步骤4：启动客户端

```bash
# 进入客户端目录
cd ~/certificate-system/fabric-workspace/client

# 启动客户端服务
./client-server
```

**客户端启动成功标志**：
```
🚀 学业证书存证客户端启动
📡 服务地址: http://localhost:8080
```

### 步骤5：访问网页

打开浏览器访问：**http://localhost:8080**

## 📋 功能说明

### 1. 证书存证

填写以下信息，点击"提交存证"：
- 证书ID：唯一标识（如 CERT001）
- 学生姓名：证书持有者姓名
- 学号：学生学号
- 专业：所学专业
- 颁发时间：证书颁发日期
- 证书哈希：证书文件的哈希值

### 2. 证书查询

输入证书ID，点击"查询证书"，显示证书详细信息。

## 🔧 命令行操作

### 查询证书

```bash
# 设置环境变量
cd ~/certificate-system/fabric-workspace/fabric-samples/test-network
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=${PWD}/../config
source scripts/envVar.sh
setGlobals 1

# 查询证书
peer chaincode query -C mychannel -n certcode -c '{"function":"queryCertificate","Args":["CERT001"]}'
```

### 上传证书

```bash
peer chaincode invoke -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile "${ORDERER_CA}" \
  -C mychannel -n certcode \
  --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} \
  --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} \
  -c '{"function":"uploadCertificate","Args":["证书ID","姓名","学号","专业","颁发时间","哈希值"]}'
```

### 停止网络

```bash
cd ~/certificate-system/fabric-workspace/fabric-samples/test-network
./network.sh down
```

### 停止客户端

```bash
# 方法1：如果在前台运行，按 Ctrl + C

# 方法2：查找并杀死进程
pkill client-server
```

## 🏗️ 架构说明

### 区块链网络

- **组织数量**：2 个（Org1、Org2）
- **通道数量**：1 个（mychannel）
- **排序节点**：1 个（orderer.example.com）
- **Peer节点**：2 个（peer0.org1、peer0.org2）

### 链码功能

| 方法 | 功能 | 参数 |
|------|------|------|
| `uploadCertificate` | 上传证书存证 | certId, name, studentId, major, issueTime, hash |
| `queryCertificate` | 查询证书信息 | certId |

### API接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/upload` | POST | 上传证书存证 |
| `/query` | GET | 查询证书信息 |
| `/` | GET | 网页首页 |

## 🐛 常见问题

### 1. Docker 权限错误

```
Got permission denied while trying to connect to the Docker daemon socket
```
**解决方案**：
```bash
sudo usermod -aG docker $USER
# 重启终端
```

### 2. peer 命令找不到

```
Command 'peer' not found
```
**解决方案**：
```bash
export PATH=~/certificate-system/fabric-workspace/fabric-samples/bin:$PATH
```

### 3. 链码部署失败

**解决方案**：
```bash
# 先停止网络
cd ~/certificate-system/fabric-workspace/fabric-samples/test-network
./network.sh down

# 重新启动
cd ~/certificate-system/fabric-workspace
./start-go.sh
```

### 4. 查询证书不存在

**解决方案**：
- 检查证书ID是否正确
- 确认网络已启动
- 确认链码已部署成功

## 📝 课设报告参考

### 系统功能
基于 Hyperledger Fabric 实现学业证书去中心化存证，防止伪造，支持存证上传和按ID查询。

### 技术栈
- Hyperledger Fabric 2.4（区块链框架）
- Docker（容器化部署）
- Go（链码开发）
- HTML/CSS/JS（前端界面）

### 核心流程
1. 启动区块链网络 → 2. 部署存证链码 → 3. 上传证书上链 → 4. 区块链不可篡改存储 → 5. 按ID查询

### 优势
- 去中心化：数据存储在多个节点
- 不可篡改：区块链特性保证数据完整性
- 安全可信：使用数字证书进行身份认证

## 📄 许可证

MIT License

---

**项目作者**：吴敏晟 陈凯鑫 卢学城
**版本**：v1.0  
**日期**：2026年5月21日
