<br />

<div align="center">
  <a href="https://barreleyescan.com/dashboard">
    <img src="https://github.com/k930503/k930503/assets/48827393/15d2445b-b46f-4056-92c8-6ec18115f29e" alt="Logo"  height="200">
  </a>

  <br />

  <h2>@Barreleye Chain &middot; <img src="https://img.shields.io/badge/Go-1.22-success" alt="go version" height="18"/></h2>

  <p align="center">
   Official open source of <b>Barreleye Blockchain. </b>

 
  with initial developer [@Youngmin Kim](https://github.com/k930503), [@Nayoung Kim](https://github.com/usiyoung)

  
</a></h6>
  </p>
</div>
<br/>
..
## Barreleye Blockchain Overview.

&nbsp;배럴아이 블록체인은 다양한 핵심 기능들을 포함하고 있습니다. 노드들은 서로 연결되어 블록과 트랜잭션을 주고받으며, 각 노드는 수신한 트랜잭션의 서명을 검증하고 이를 멤풀에 저장한 후, 다른 노드들에게 전파합니다. 블록 생성이 완료되면, 생성한 노드는 보상을 받고 멤풀에 담긴 트랜잭션들을 처리해 상태를 업데이트한 후 블록을 체인에 연결하고, 이를 다시 다른 노드들과 공유합니다. 연결된 노드들은 블록 서명 검증 후, 블록 내 트랜잭션을 처리해 상태를 동기화하며 체인에 반영하고 전파합니다.

&nbsp;블록체인은 데이터 무결성을 유지하기 위해 블록 분기와 같은 다양한 상황에서도 배럴아이 코어만의 메커니즘을 통해 문제없이 동작합니다. 또한, 각 노드는 블록체인 상태 정보를 조회할 수 있는 API를 제공하며, 실시간으로 배럴아이스캔(https://barreleyescan.com) 을 통해 블록체인 현황을 확인할 수 있습니다. 누구나 GitHub의 매뉴얼을 참고해 배럴아이 노드를 구축하여 메인 네트워크에 참여하거나 개인 네트워크를 설정할 수도 있습니다.

<br/>

## Barreleye Usage.

### Prerequisites.

Docker download here [docker.com](https://www.docker.com/products/docker-desktop/).

## **1. Pull Docker Image.**
Pull the Barreleye Docker image.

```shell
$ docker pull kym6772/barreleye:1.0.0
```

 

## **2. Write a shell script.**
Fill in the variables needed to run the node.
```text
# example
name="my-node"
role="normal"
port="4100"
peers="172.30.1.5:4101"
httpPort="9000"
key="a2288db63c7016b815c55c1084c2491b8599834500408ba863ec379895373ae9"
docker run -it -p ${port}:${port} -d kym6772/barreleye:1.0.0 /barreleye/bin/barreleye -name=${name} -role=${role} -port=${port} -peer=${peer} -http.port=${httpPort} -key=${key}
```

* `name` - the node name you want.
* `role` - If it is the first node running in a private network, the role is `genesis`, otherwise it is `normal`.
* `port` - Port number for communication between nodes based on TCP/IP.
* `peers` - Peer's port number. If role is genesis, fill in `none`. also, it can be an array. For example, "x.x.x.x:3000,y.y.y.y:4000,..."
* `httpPort` - Port number for REST API.
* `key` - Node’s private key for signing and verifying blocks.

 
## **3. Run a shell script.**
```shell
$ ./{file_name}.sh
```


Result of executing the command.

<img width="1210" alt="tutorial1" src="https://github.com/barreleye-labs/barreleye/assets/48827393/abc5a149-024a-449e-afb2-675822b3c7e2" alt="Result of executing the command.">
<br>
<br>
If this is the first node in your private network, it will stop at a line like the one above. This is because mining begins only when two or more nodes participate. Run two or more nodes.<br>
<br>

![MergedImages](https://github.com/barreleye-labs/barreleye/assets/48827393/e84562af-6f64-4e72-ab41-16a18031fa68)

You can connect infinite nodes as shown above. As you can see from the log, nodes verify and process transactions. Nodes then broadcast blocks and transactions to synchronize data with each other. In this way, nodes earn rewards through mining in return for maintaining the Barreleye blockchain network. Let’s participate as a node in the main network. Or let's build your own private network!

<br/>

## **REST API Documentation.**

|        path        | method | request                                                                                                                                                                                                                                                                                                                                                                                                                                    | response                                                                                                                                 |
|:------------------:|:------:|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------|
|      /blocks       | `GET`  | `query`<br/>page<br/>size                                                                                                                                                                                                                                                                                                                                                                                                                  | blocks                                                                                                                                   |
|    /blocks/:id     | `GET`  | `param`<br/>id - hash or height                                                                                                                                                                                                                                                                                                                                                                                                            | hash<br/>version<br/>dataHash<br/>prevBlockHash<br/>height<br/>timestamp<br/>signer<br/>extra<br/>signature<br/>txCount<br/>transactions |
|    /last-block     | `GET`  | none                                                                                                                                                                                                                                                                                                                                                                                                                                       | block                                                                                                                                    |
|        /txs        | `GET`  | `query`<br/>page<br/>size                                                                                                                                                                                                                                                                                                                                                                                                 | transactions                                                                                                                             |
|      /txs/:id      | `GET`  | `param`<br/>id - hash or number                                                                                                                                                                                                                                                                                                                                                                                                            | hash<br/>nonce<br/>blockHeight<br/>timestamp<br/>from<br/>to<br/>value<br/>data<br/>signer<br/>signature                                 |
|        /txs        | `POST` | `body`<br/>from - <span style="color:gray">*hex string*</span><br/>to - <span style="color:gray">*hex string*</span><br/>value - <span style="color:gray">*hex string*</span><br/>data - <span style="color:gray">*hex string*</span><br/>signerX - <span style="color:gray">*hex string*</span><br/>signerY - <span style="color:gray">*hex string*</span><br/>signatureR - <span style="color:gray">*hex string*</span><br/>signatureS - <span style="color:gray">*hex string*</span> | transaction                                                                                                                              |
|      /faucet       | `POST` | `body`<br/>accountAddress - <span style="color:gray">*hex string*</span>                                                                                                                                                                                                                                                                                                                                                                   | transaction                                                                                                                              |
| /accounts/:address &nbsp; | `GET`  | `param`<br/>address                                                                                                                                                                                                                                                                                                                                                                                                                        | address<br/>nonce<br/>balance                                                                                                |                                                                                                          |

<br/>

## **Specification.**
* `Block time` - 10 seconds on average.<br>
* `Block reward` - 10 barrel per block.<br>
* `Hash algorithm` - SHA256.<br>
* `Cryptography algorithm` - ECDSA secp256k1.<br>
* `Consensus algorithm` - Proof of random

<br/>

## Explorer & Wallet.
https://barreleyescan.com

<br/>

## Our projects.
![barreleye-fish-black-24](https://github.com/barreleye-labs/barreleye/assets/48827393/698b04c7-454a-4cb9-8680-ac5647b558fc)&nbsp;&nbsp;&nbsp;[Barreleye](https://github.com/barreleye-labs/barreleye)

![barreleye-fish-black-24](https://github.com/barreleye-labs/barreleye/assets/48827393/698b04c7-454a-4cb9-8680-ac5647b558fc)&nbsp;&nbsp;&nbsp;[Barreleyescan](https://github.com/barreleye-labs/barreleye-explorer-react)

<br/>

## Please inquire about participating in the main network.
* k930503@gmail.com<br>
* usiyoung7@gmail.com
