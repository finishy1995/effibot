# effibot

[English](./README.md) | [中文](./README_CN.md)

A ChatGPT server that stores and processes data using a tree-based data structure, providing users with a mind map-like 
Q&A experience with ChatGPT. The tree structure greatly optimizes the transmission of context (tokens) and provides 
a better experience when used within a company.

![Demo](docs/demo.png)

> The image shows a Demo client; the UI is for reference only.

In work scenarios, the need to deeply ask the same question is relatively rare, so in most cases, the token count can be 
controlled within 2000. Therefore, GPT 3.5's token limit (4096) is sufficient (no need to consider GPT4 for accuracy).

![Data](docs/data.png)

> The time between the two screenshots does not exceed 5 minutes. Due to multiple users, detailed logs need to be viewed 
> to distinguish the token consumption of the five questions mentioned, but the overall token consumption can be seen as 
> controllable.

## Demo

[43.206.107.75:4000](http://43.206.107.75:4000)

Demo environment is deployed on a cloud server. And **NOT** set the OpenAI token, so it will start mock mode.

## 📢Update Plan

- [ ] User login
- [ ] Data persistence storage
- [x] Quick deployment script, uploaded to Docker Hub
- [ ] Add support for writing long-form fiction scenarios

**Updates will be made as needed.** More updates will be provided if the project is widely used, and updates will be 
made based on interest if the project has fewer users.

Feel free to develop a Web UI based on this project! The UI in the Demo is written by me, a beginner in UI design. 
PRs are welcome!

## Approach

Organize user input into a multi-branch tree and only pass the content of the current branch as context information to 
GPT. The amount of content we transmit each time is equal to the depth of the current node. Optimize the selection and 
transmission of context through the multi-branch tree.

A binary tree with n nodes has a depth of logn. The depth here refers to the context information we need to pass to the 
GPT API. If we don't process the context, it can be regarded as a one-dimensional tree, which degenerates into a line 
segment, naturally the most complex case. By organizing the session into a tree structure, we can create a mind map.

## Environment Notice
It is recommended to choose a server location in a country or region supported by OpenAI. Data centers and cloud hosts
are both acceptable, and the following clouds have been tested:
1. Azure
2. AWS

If you insist on testing in an unsupported country or region, this project fully supports proxies, but the proxy itself
may **affect the experience and pose risks**. See the configuration file Spec.GPT.TransportUrl for the proxy
configuration details.

The use of proxies is not recommended. Use at your own risk.

## Clone and enter project directory

```bash
git clone https://github.com/finishy1995/effibot.git
cd effibot
```

## Configuration

The default configuration is Mock mode, which means it will not actually call the GPT API but return the user's input
as a response. The default REST API port is `4001`, and all configurations can be modified in the
`http/etc/http-api.yaml` file.

```bash
vi http/etc/http-api.yaml
```

```yaml
Name: http-api
Host: 0.0.0.0
Port: 4001 # Port of http server, default 4001
Timeout: 30000 # Timeout of http request, default 30000(ms)
Log:
  Level: debug
  Mode: file # Log mode, default console 日志模式，可选 console（命令行输出） 或 file
  Path: ../logs # Log file path, default ../logs
Spec:
  GPT:
#    Token: "sk-" # Token of OpenAI, will start mock mode if not set. OpenAI 密钥，如果不设置则启用 mock 模式
#    TransportUrl: "http://localhost:4002" # Transport url of OpenAI, default "http://localhost:4002 代理地址，如果不设置则不启用代理
    Timeout: 20s # Timeout of OpenAI request, default 20s
    MaxToken: 1000 # Max token of OpenAI response, default 1000
```

After modifying the file, if you need `One-click deployment` or `container deployment`, please execute the following command
```bash
mkdir -p ./effibot_config
cp http/etc/http-api.yaml ./effibot_config
```

## One-click deployment

Please ensure that `docker` and `docker-compose` are correctly installed and enabled.

```bash
docker-compose up -d
```

Demo client will run on port `4000`, and REST API will run on both ports `4000` and `4001`.

If you don't have `docker-compose`, you can use the following command:

```bash
docker network create effibot
docker run -p 4001:4001 -v ./effibot_config:/app/etc --network effibot --name effibot -d finishy/effibot:latest
docker run -p 4000:4000 --network effibot --name effibot-demo -d finishy/effibot-demo:latest
```

## Custom Development/Deployment

### Server Development

Ensure that golang 1.18+ is installed and configure.

```bash
cd http
go run http.go # go build http.go && ./http
```

Exit directory
```bash
cd ..
```

### Container Packaging

```bash
docker build -t effibot:latest -f http/Dockerfile .
```

### Container Network Bridge Create

```bash
docker network create effibot
```

### Container Deployment

```bash
# Modify the configuration file as needed, such as adding the OpenAI token and change the log mode to console
docker run -p 4001:4001 -v ./effibot_config:/app/etc --network effibot --name effibot -d effibot:latest
```

### Demo Client Container Packaging

```bash
docker build -t effibot-demo:latest -f demo/Dockerfile .
```

### Demo Client Container Deployment

```bash
docker run -p 4000:4000 --network effibot --name effibot-demo -d effibot-demo:latest
```

### Demo Client Development

Demo client is developed by Vue.js + Vite + TypeScript, and requires Node.js 14+ environment.

```bash
cd demo
yarn && yarn dev
```

Demo client will automatically open at [http://localhost:5173](http://localhost:5173).
