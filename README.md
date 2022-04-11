# DEM-Chain
Digital Experience Monitoring on the Blockchain.

# Instructions
Check out this code on a client machine. 

```
git clone https://github.com/tonylepage/demchain.git
```

The Client must have access to a node. You will need your MEMBER_ID and NODE_ENDPOINT_URL.

Create a Docker Compose configuration file named docker-compose-cli.yaml in the /home/ec2-user directory, which you use to run the Hyperledger Fabric CLI. You use this CLI to interact with peer nodes that your member owns. 

```yaml
version: '2'
services:
  cli:
    container_name: cli
    image: hyperledger/fabric-tools:2.2
    tty: true
    environment:
      - GOPATH=/opt/gopath
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - FABRIC_LOGGING_SPEC=info # Set logging level to debug for more verbose logging
      - CORE_PEER_ID=cli
      - CORE_CHAINCODE_KEEPALIVE=10
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_TLS_ROOTCERT_FILE=/opt/home/managedblockchain-tls-chain.pem
      - CORE_PEER_LOCALMSPID=MEMBER_ID
      - CORE_PEER_MSPCONFIGPATH=/opt/home/admin-msp
      - CORE_PEER_ADDRESS=NODE_ENDPOINT_URL
    working_dir: /opt/home
    command: /bin/bash
    volumes:
        - /var/run/:/host/var/run/
        - /home/ec2-user/demchain:/opt/gopath/src/github.com/
        - /home/ec2-user:/opt/home
```

Run the following command to start the Hyperledger Fabric peer CLI container:

```
docker-compose -f docker-compose-cli.yaml up -d
```

Setup an environment variable called ORDERER (use your own node URi):
```
export ORDERER=orderer.n-MWY63ZJZU5HGNCMBQER7IN6OIU.managedblockchain.amazonaws.com:30001
```

Create a channel:
```
docker exec cli peer channel create -c demchannel \
-f /opt/home/demchannel.pb -o $ORDERER \
--cafile /opt/home/managedblockchain-tls-chain.pem --tls
```

Join the channel:
```
docker exec cli peer channel join -b demchannel.block \
-o $ORDERER --cafile /opt/home/managedblockchain-tls-chain.pem --tls
```

Install the dependencies:
```
sudo chown -R ec2-user:ec2-user demchain/
cd demchain/
GO111MODULE=on go mod vendor
cd -
```
