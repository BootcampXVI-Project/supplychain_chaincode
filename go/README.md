## Running the sample

The Fabric test network is used to deploy and run this sample. Follow these steps in order:

1. Create the network and a channel (from the `SupplychainNetwork` folder).
```bash
   ./network.sh up createChannel  -ca -s couchdb
```

1. Deploy one of the smart contract implementations (from the `SupplychainNetwork` folder).

# To deploy the Go chaincode implementation
```bash
   ./network.sh deployCC -ccn basic -ccp ../supplychain_chaincode/go/ -ccl go
```

## config peer command

```bash
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
```

You can then set up the environment variables for each organization. The `./scripts/envVar.sh` command is designed to be run as follows.

```bash
source ./scripts/envVar.sh && setGlobals $ORG
```

## after deploy chain code 
### Create Ledger
```bash
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.supplychain.com --tls --cafile "${PWD}/organizations/ordererOrganizations/supplychain.com/orderers/orderer.supplychain.com/msp/tlscacerts/tlsca.supplychain.com-cert.pem" -C supplychain-channel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/supplier.supplychain.com/peers/peer0.supplier.supplychain.com/tls/ca.crt" --peerAddresses localhost:7061 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/manufacturer.supplychain.com/peers/peer0.manufacturer.supplychain.com/tls/ca.crt" --peerAddresses localhost:7071 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/distributor.supplychain.com/peers/peer0.distributor.supplychain.com/tls/ca.crt"  --peerAddresses localhost:7081 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/retailer.supplychain.com/peers/peer0.retailer.supplychain.com/tls/ca.crt" --peerAddresses localhost:7091 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/consumer.supplychain.com/peers/peer0.consumer.supplychain.com/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}'
```
#### You can check results by couchDB at ......_basic

### ADD User
```bash
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.supplychain.com --tls --cafile "${PWD}/organizations/ordererOrganizations/supplychain.com/orderers/orderer.supplychain.com/msp/tlscacerts/tlsca.supplychain.com-cert.pem" -C supplychain-channel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/supplier.supplychain.com/peers/peer0.supplier.supplychain.com/tls/ca.crt" --peerAddresses localhost:7061 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/manufacturer.supplychain.com/peers/peer0.manufacturer.supplychain.com/tls/ca.crt" --peerAddresses localhost:7071 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/distributor.supplychain.com/peers/peer0.distributor.supplychain.com/tls/ca.crt"  --peerAddresses localhost:7081 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/retailer.supplychain.com/peers/peer0.retailer.supplychain.com/tls/ca.crt" --peerAddresses localhost:7091 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/consumer.supplychain.com/peers/peer0.consumer.supplychain.com/tls/ca.crt" -c '{"function":"CreateUser","Args":["giahung@gmail.com","giahung","giahung","DaNang","supplier","supplier"]}'
```
### Then view it
```bash
peer chaincode query -C supplychain-channel -n basic -c '{"Args":["GetAllUsers"]}'
```
