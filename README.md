# DEM-Chain
Digital Experience Monitoring on the Blockchain.

# Installation Instructions
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

To check which chaincode has been installed on the node:
```
$ docker exec cli peer lifecycle chaincode queryinstalled
```

In this example, 2 packages have been installed, so the command returns:
```
Installed chaincodes on peer:
Package ID: abstore_1:3918d0438fd2ebe48ed1bde01533513a14f788846fd2d72ef054482760e73409, Label: abstore_1
Package ID: demchain_1:e60917100fc9af5d6bca17592d78711b077af972861d382970107fef2d0e9cdc, Label: demchain_1
```

Now that the chaincode is deployed on a node, it should be approved. First we need to check commit readiness:
```
export CC_PACKAGE_ID=demchain_1:084cd76392ebd046fafe36d75ba91467c0fc313bbccad83489b11d941ad42fe9
docker exec cli peer lifecycle chaincode approveformyorg \
--orderer $ORDERER --tls --cafile /opt/home/managedblockchain-tls-chain.pem \
--channelID demchannel --name demcc --version v0 --sequence 1 --package-id $CC_PACKAGE_ID
```

And returns:
```
2022-04-12 02:28:31.169 UTC [chaincodeCmd] ClientWait -> INFO 001 txid [769799fac6bdf25e7d8e090276587e7a85cc72bbd1e12129256edc65277d97a9] committed with status (VALID) at nd-t22kaurwnvcbbdqntsxrtsqse4.m-ehbcgjxqgbakdjjyl3uewcdhzu.n-4yfjzomgsvadtlyb63niexxspe.managedblockchain.us-east-1.amazonaws.com:30003
```

Then we can commit the chaincode:
```
docker exec cli peer lifecycle chaincode commit \
--orderer $ORDERER --tls --cafile /opt/home/managedblockchain-tls-chain.pem \
--channelID demchannel --name demcc --version v0 --sequence 1
```

And finally, verify it:
```
docker exec cli peer lifecycle chaincode querycommitted \
--channelID demchannel
```

# Running the Chaincode

To execute methods on the chaincode, continue to use the cli.

Initialise the chaincode using this command:
```
docker exec cli peer chaincode invoke --tls --cafile /opt/home/managedblockchain-tls-chain.pem --channelID demchannel --name demcc -c '{"Function":"InitLedger","Args":[""]}'
```

To retrieve all the measurement values from the chaincode:
```
docker exec cli peer chaincode invoke --tls --cafile /opt/home/managedblockchain-tls-chain.pem --channelID demchannel --name demcc -c '{"Function":"GetAllMeasurements","Args":[""]}'
```

To add a new measurement:
```
docker exec cli peer chaincode invoke --tls --cafile /opt/home/managedblockchain-tls-chain.pem --channelID demchannel --name demcc -c '{"Function":"CreateMeasurement","Args":["Taipei, Taiwan", 1649602278, 217, "Stackpath", "console-tester"]}'
```
