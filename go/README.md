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

## after deploy chain code 


### predeclare for authorize
```bash
ORDERER_ADDRESS=localhost:7050
ORDERER_TLS_CERT="${PWD}/organizations/ordererOrganizations/supplychain.com/orderers/orderer.supplychain.com/msp/tlscacerts/tlsca.supplychain.com-cert.pem"

SUPPLIER_PEER_ADDRESS=localhost:7051
SUPPLIER_PEER_TLS_CERT="${PWD}/organizations/peerOrganizations/supplier.supplychain.com/peers/peer0.supplier.supplychain.com/tls/ca.crt"

MANUFACTURER_PEER_ADDRESS=localhost:7061
MANUFACTURER_PEER_TLS_CERT="${PWD}/organizations/peerOrganizations/manufacturer.supplychain.com/peers/peer0.manufacturer.supplychain.com/tls/ca.crt"

DISTRIBUTOR_PEER_ADDRESS=localhost:7071
DISTRIBUTOR_PEER_TLS_CERT="${PWD}/organizations/peerOrganizations/distributor.supplychain.com/peers/peer0.distributor.supplychain.com/tls/ca.crt"

RETAILER_PEER_ADDRESS=localhost:7081
RETAILER_PEER_TLS_CERT="${PWD}/organizations/peerOrganizations/retailer.supplychain.com/peers/peer0.retailer.supplychain.com/tls/ca.crt"

CONSUMER_PEER_ADDRESS=localhost:7091
CONSUMER_PEER_TLS_CERT="${PWD}/organizations/peerOrganizations/consumer.supplychain.com/peers/peer0.consumer.supplychain.com/tls/ca.crt"
```
### defined name channel and database
```bash
CHANNEL_NAME=supplychain-channel
CHAINCODE_NAME=basic
```

### defined function want to call
```bash
INVOKE_PARAMS='{"function":"InitLedger","Args":[]}'
```

### create user
```bash
'{"function":"CreateUser","Args":["giahung@gmail.com","giahung","giahung","DaNang","supplier","supplier"]}'

'{"function":"CreateUser","Args":["nero@gmail.com","nero","nero","Hue","manufacturer","manufacturer"]}'

'{"function":"CreateUser","Args":["nora@gmail.com","nora","nora","Hue","distributor","distributor"]}'

'{"function":"CreateUser","Args":["eden@gmail.com","eden","eden","Quang Nam","retailer","retailer"]}'

```
### get all user
```bash
'{"Args":["GetAllUsers"]}'
```

### get all product
```bash
'{"Args":["GetAllProducts"]}'
```

### supplier's functions
### cultivate product
```bash
'{"function":"CultivateProduct","Args":["User1","FirstProduct","109000.00","first product"]}'
```

### harvert product
```bash
'{"function":"HarvertProduct","Args":["User1","Product1"]}'
```

### manufacturer's functions
### import product
```bash
'{"function":"ImportProduct","Args":["User2","Product1","219000.00"]}'
```

### manufacture product
```bash
'{"function":"ManufactureProduct","Args":["User2","Product1"]}'
```

### export product
```bash
'{"function":"ExportProduct","Args":["User2","Product1","219000.00"]}'
```

### distributor's functions
### distribute product
```bash
'{"function":"DistributeProduct","Args":["User3","Product1"]}'
```

### retailer's functions
### sell product
```bash
'{"function":"SellProduct","Args":["User4","Product1","299000.00"]}'
```


### call function
```bash
peer chaincode invoke \
  -o $ORDERER_ADDRESS \
  --ordererTLSHostnameOverride orderer.supplychain.com \
  --tls \
  --cafile $ORDERER_TLS_CERT \
  -C $CHANNEL_NAME \
  -n $CHAINCODE_NAME \
  --peerAddresses $SUPPLIER_PEER_ADDRESS \
  --tlsRootCertFiles $SUPPLIER_PEER_TLS_CERT \
  --peerAddresses $MANUFACTURER_PEER_ADDRESS \
  --tlsRootCertFiles $MANUFACTURER_PEER_TLS_CERT \
  --peerAddresses $DISTRIBUTOR_PEER_ADDRESS \
  --tlsRootCertFiles $DISTRIBUTOR_PEER_TLS_CERT \
  --peerAddresses $RETAILER_PEER_ADDRESS \
  --tlsRootCertFiles $RETAILER_PEER_TLS_CERT \
  --peerAddresses $CONSUMER_PEER_ADDRESS \
  --tlsRootCertFiles $CONSUMER_PEER_TLS_CERT \
  -c $INVOKE_PARAMS
```