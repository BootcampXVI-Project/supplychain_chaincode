package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type CounterNO struct {
	counter int `json:"counter"`
}

type User struct {
	userId   string `json:"userId"`
	email    string `json:"email"`
	password string `json:"password"`
	userName string `json:"userName"`
	address  string `json:"address"`
	userType string `json:"userType"`
	role     string `json:"role"`
	status   string `json:"status"`
	identify string `json:"identify"`
}


type ProductDates struct {
	cultivated     string `json:"cultivated"` // supplier
	harvested      string `json:"harvested"`
	imported       string `json:"imported"` // manufacturer
	manufacturered string `json:"manufacturered"`
	exported       string `json:"exported"`
	distributed    string `json:"distributed"` // distributor
	sold           string `json:"sold"`        // retailer
}

type ProductActors struct {
	supplierId     string `json:"supplierId"`
	manufacturerId string `json:"manufacturerId"`
	distributorId  string `json:"distributorId"`
	retailerId     string `json:"retailerId"`
}

// Unit: kg, box/boxes, bottle, bottles

// Supplier: id, cultivate, harvest => cultivating, harvested
// Manufacturer: id, import, manufacture, export => imported, manufacturing, exported
// Distributor: id, distribute => distributed/distributing
// Retailer: id, sell => sold
type Product struct {
	productId   	string        `json:"productId"`
	image 			[]string	  `json:"image"`
	productName 	string        `json:"productName"`
	dates       	ProductDates  `json:"dates"`
	actors      	ProductActors `json:"actors"`
	price       	string        `json:"price"`
	status      	string        `json:"status"`
	description 	string        `json:"description"`
	certificateURL 	string 		  `json:"certificate"`
	cooperationId 	string 		  `json:"cooperationId"`
}

type ProductHistory struct {
	record    *Product  `json:"record"`
	txId      string    `json:"txId"`
	timestamp time.Time `json:"timestamp"`
	isDelete  bool      `json:"isDelete"`
}

// order

type Signature struct {
	distributorSignature  	string 	`json:"distributorSignature"`
	retailerSignature 		string  `json:"retailerSignature"`
}


type ProductItem struct {
	product  Product `json:"product"`
	quantity string  `json:"quantity"`
}

type DeliveryStatus struct {
	distributedId 	string 		`json:"distributedId"`
	DeliveryDate 	string		`json:"deliveryDate"`
	Status       	string    	`json:"status"`
}

type Order struct {
	orderID 		string      	`json:"orderID"`
	productItemList []ProductItem 	`json:"productItemList"`
	signature 		Signature 		`json:"signature"`
	// DateCreate 		string 			`json:"dateCreate"`
	// DateFinish      string      	`json:"dateFinish"`
	deliveryStatus 	[]DeliveryStatus `json:"deliveryStatus"`
	status          string     	 	`json:"status"`
	distributorId  	string 			`json:"distributorId"`
	retailerId     	string 			`json:"retailerId"`
}

// Init initializes chaincode
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	error := initCounter(ctx)

	if error != nil {
		return fmt.Errorf("error init counter: %s", error.Error())
	}

	return nil
}

func initCounter(ctx contractapi.TransactionContextInterface) error {
	// Initializing Product Counter
	ProductCounterBytes, _ := ctx.GetStub().GetState("ProductCounterNO")
	if ProductCounterBytes == nil {
		var ProductCounter = CounterNO{counter: 0}
		ProductCounterBytes, _ := json.Marshal(ProductCounter)
		err := ctx.GetStub().PutState("ProductCounterNO", ProductCounterBytes)
		if err != nil {
			return fmt.Errorf("failed to Intitate Product Counter: %s", err.Error())
		}
	}
	OrderCounterBytes, _ := ctx.GetStub().GetState("OrdertCounterNO")
	if OrderCounterBytes == nil {
		var OrderCounter = CounterNO{counter: 0}
		OrderCounterBytes, _ := json.Marshal(OrderCounter)
		err := ctx.GetStub().PutState("OrderCounterNO", OrderCounterBytes)
		if err != nil {
			return fmt.Errorf("failed to Intitate Order Counter: %s", err.Error())
		}
	}

	return nil
}

// getCounter to the latest value of the counter based on the Asset Type provided as input parameter
func getCounter(ctx contractapi.TransactionContextInterface, assetType string) (int, error) {
	counterAsBytes, _ := ctx.GetStub().GetState(assetType)
	counterAsset := CounterNO{}

	json.Unmarshal(counterAsBytes, &counterAsset)
	// fmt.Sprintf("Counter Current Value %d of Asset Type %s", counterAsset.Counter, assetType)
	return counterAsset.counter, nil
}

// incrementCounter to the increase value of the counter based on the Asset Type provided as input parameter by 1
func incrementCounter(ctx contractapi.TransactionContextInterface, assetType string) (int, error) {
	counterAsBytes, _ := ctx.GetStub().GetState(assetType)
	counterAsset := CounterNO{}

	json.Unmarshal(counterAsBytes, &counterAsset)
	counterAsset.counter++
	counterAsBytes, _ = json.Marshal(counterAsset)

	err := ctx.GetStub().PutState(assetType, counterAsBytes)
	if err != nil {
		return -1, fmt.Errorf("failed to Increment Counter: %s", err.Error())
	}
	fmt.Printf("Printf in incrementing counter  %v", counterAsset)
	return counterAsset.counter, nil
}

// GetTxTimestampChannel Function gets the Transaction time when the chain code was executed it remains same on all the peers where chaincode executes
func (s *SmartContract) GetTxTimestampChannel(ctx contractapi.TransactionContextInterface) (string, error) {
	txTimeAsPtr, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		fmt.Printf("Returning error in TimeStamp \n")
		return "Error", err
	}
	fmt.Printf("\t returned value from ctx.GetStub(): %v\n", txTimeAsPtr)
	timeStr := time.Unix(txTimeAsPtr.Seconds, int64(txTimeAsPtr.Nanos)).String()
	return timeStr, nil
}


// SUPPLIER FUNCTION
// cultivate product // gieo trồng sảm phẩm
func (s *SmartContract) CultivateProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {

	if user.userType != "supplier" {
		return fmt.Errorf("user must be a supplier")
	}

	productCounter, _ := getCounter(ctx, "ProductCounterNO")
	productCounter++

	//To Get the transaction TimeStamp from the Channel Header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in Transaction TimeStamp")
	}

	// DATES
	dates := ProductDates{}
	dates.cultivated = txTimeAsPtr
	actors := ProductActors{}
	actors.supplierId = user.userId
	var product = Product{
		productId:   "Product" + strconv.Itoa(productCounter),
		productName: productObj.productName,
		dates:       dates,
		actors:      actors,
		price:       productObj.price,
		status:      "CULTIVATING",
		description: productObj.description,
		cooperationId : productObj.cooperationId,
	}

	productAsBytes, _ := json.Marshal(product)

	incrementCounter(ctx, "ProductCounterNO")

	return ctx.GetStub().PutState(product.productId, productAsBytes)
}

// havert product // thu hoạch
func (s *SmartContract) HarvertProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {

	if user.userType != "supplier" {
		return fmt.Errorf("user must be a supplier")
	}

	// get product details from the stub ie. Chaincode stub in network using the product id passed
	productBytes, _ := ctx.GetStub().GetState(productObj.productId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	//To Get the transaction TimeStamp from the Channel Header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in Transaction TimeStamp")
	}

	// Updating the product values withe the new values
	product.dates.harvested = txTimeAsPtr
	product.status = "HAVERTED"

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.productId, updatedProductAsBytes)
}

// supplier update
func (s *SmartContract) SupplierUpdateProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {

	if user.userType != "supplier" {
		return fmt.Errorf("user must be a supplier")
	}

	// get product details from the stub ie. Chaincode stub in network using the product id passed
	productBytes, _ := ctx.GetStub().GetState(productObj.productId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	// Updating the product values withe the new values
	product.productName = productObj.productName
	product.price = productObj.price
	product.description = productObj.description

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.productId, updatedProductAsBytes)
}
func (s *SmartContract) AddCertificate(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {

	if user.userType != "supplier" {
		return fmt.Errorf("user must be a supplier")
	}

	// get product details from the stub ie. Chaincode stub in network using the product id passed
	productBytes, _ := ctx.GetStub().GetState(productObj.productId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}
	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	product.certificateURL = productObj.certificateURL
	
	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.productId, updatedProductAsBytes)

}
// MANUFACTURER
// import product
func (s *SmartContract) ImportProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {

	if user.userType != "manufacturer" {
		return fmt.Errorf("user must be a manufacturer")
	}

	productBytes, _ := ctx.GetStub().GetState(productObj.productId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	//To Get the transaction TimeStamp from the Channel Header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in Transaction TimeStamp")
	}

	// Updating the product values withe the new values
	product.image = productObj.image
	product.dates.imported = txTimeAsPtr
	product.price = productObj.price
	product.status = "IMPORTED"
	product.actors.manufacturerId = user.userId

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.productId, updatedProductAsBytes)
}

// manufacture product
func (s *SmartContract) ManufactureProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {

	if user.userType != "manufacturer" {
		return fmt.Errorf("user must be a manufacturer")
	}

	productBytes, _ := ctx.GetStub().GetState(productObj.productId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	//To Get the transaction TimeStamp from the Channel Header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in Transaction TimeStamp")
	}

	if product.actors.manufacturerId != user.userId {
		return fmt.Errorf("Permission denied!")
	}

	// Updating the product values withe the new values
	product.dates.manufacturered = txTimeAsPtr
	product.status = "MANUFACTURED"

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.productId, updatedProductAsBytes)
}

// export product
func (s *SmartContract) ExportProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {

	if user.userType != "manufacturer" {
		return fmt.Errorf("user must be a manufacturer")
	}

	productBytes, _ := ctx.GetStub().GetState(productObj.productId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	//To Get the transaction TimeStamp from the Channel Header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in Transaction TimeStamp")
	}

	if product.actors.manufacturerId != user.userId {
		return fmt.Errorf("Permission denied!")
	}

	// Updating the product values withe the new values
	product.dates.exported = txTimeAsPtr
	product.price = productObj.price
	product.status = "EXPORTED"

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.productId, updatedProductAsBytes)
}

// DISTRIBUTOR
// distribute product
func (s *SmartContract) DistributeProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {

	if user.userType != "distributor" {
		return fmt.Errorf("user must be a distributor")
	}

	productBytes, _ := ctx.GetStub().GetState(productObj.productId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	//To Get the transaction TimeStamp from the Channel Header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in Transaction TimeStamp")
	}

	// Updating the product values withe the new values
	// product.Dates.distributed[0].distributedId = user.userId
	product.dates.distributed = txTimeAsPtr
	// product.Dates.distributed[0].Status = "Start delivery"

	product.status = "DISTRIBUTED"
	product.actors.distributorId = user.userId

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.productId, updatedProductAsBytes)
}

// RETAILER
// sell product
func (s *SmartContract) SellProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {

	if user.userType != "retailer" {
		return fmt.Errorf("user must be a retailer")
	}

	// get product details from the stub ie. Chaincode stub in network using the product id passed
	productBytes, _ := ctx.GetStub().GetState(productObj.productId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	//To Get the transaction TimeStamp from the Channel Header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in Transaction TimeStamp")
	}

	// Updating the product values to be updated after the function
	product.dates.sold = txTimeAsPtr
	product.status = "SOLD"
	product.price = productObj.price
	product.actors.retailerId = user.userId

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.productId, updatedProductAsBytes)
}

// get a asset
func (s *SmartContract) GetProduct(ctx contractapi.TransactionContextInterface, productId string) (*Product, error) {
	productAsBytes, err := ctx.GetStub().GetState(productId)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	}

	if productAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", productId)
	}

	product := new(Product)
	_ = json.Unmarshal(productAsBytes, product)

	return product, nil
}

func (s *SmartContract) GetAllProducts(ctx contractapi.TransactionContextInterface, productObj Product) ([]*Product, error) {

	assetCounter, _ := getCounter(ctx, "ProductCounterNO")
	startKey := "Product1"
	endKey := "Product" + strconv.Itoa(assetCounter+1)
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	var products []*Product

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var product Product
		_ = json.Unmarshal(response.Value, &product)

		products = append(products, &product)
	}

	return products, nil
}

func (s *SmartContract) CreateOrder(ctx contractapi.TransactionContextInterface,user User,orderObj Order ) error {
	// newOrder := Order{}
	// for _, o := range orders {
	// 	newOrder = append(newOrder, struct {
	// 		Product  Product
	// 		Quantity string
	// 	}{
	// 		Product:  o.Product,
	// 		Quantity: o.Quantity,
	// 	})
	// }
	// return newOrder
	if user.userType != "distributor" {
		return fmt.Errorf("user must be a distributor")
	}

	orderCounter, _ := getCounter(ctx, "OrderCounterNO")
	orderCounter++

	//To Get the transaction TimeStamp from the Channel Header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in Transaction TimeStamp") 
	}

	firstdelivery := DeliveryStatus{
		distributedId: user.userId,
        Status:     "Start delivery",
        DeliveryDate:  txTimeAsPtr,
	}
	var deliveryStatus []DeliveryStatus

	deliveryStatus = append(deliveryStatus, firstdelivery)

	// DATES
	var order = Order{
		orderID:   			"Order" + strconv.Itoa(orderCounter),
		productItemList: 	orderObj.productItemList,
		signature:       	orderObj.signature,
		deliveryStatus:     deliveryStatus,
		status:     		orderObj.status,
		distributorId: 		user.userId,
		retailerId: 		orderObj.retailerId,
	}

	orderAsBytes, _ := json.Marshal(order)

	incrementCounter(ctx, "OrderCounterNO")

	return ctx.GetStub().PutState(order.orderID, orderAsBytes)
}

func (s *SmartContract) updateOrder(ctx contractapi.TransactionContextInterface,user User,orderObj Order ) error {

	if user.userType != "distributor" {
		return fmt.Errorf("user must be a distributor")
	}

	//To Get the transaction TimeStamp from the Channel Header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in Transaction TimeStamp") 
	}

	orderBytes, _ := ctx.GetStub().GetState(orderObj.orderID)
	if orderBytes == nil {
		return fmt.Errorf("cannot find this order")
	}

	order := new(Order)
	_ = json.Unmarshal(orderBytes, order)

	if order.distributorId != user.userId {
		return fmt.Errorf("Permission denied!")
	}
	delivery := DeliveryStatus{
		distributedId: user.userId,
        Status:     "Delivering "+ user.address,
        DeliveryDate:  txTimeAsPtr,
	}
	order.deliveryStatus = append(order.deliveryStatus, delivery)
	order.status = orderObj.status
	// order.Signature = orderObj.Signature

	updateOrderAsBytes, _ := json.Marshal(order)

	return ctx.GetStub().PutState(order.orderID, updateOrderAsBytes)
}

func (s *SmartContract) finishOrder(ctx contractapi.TransactionContextInterface,user User,orderObj Order ) error {

	if user.userType != "retailer" {
		return fmt.Errorf("user must be a retailer")
	}

	//To Get the transaction TimeStamp from the Channel Header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in Transaction TimeStamp") 
	}

	orderBytes, _ := ctx.GetStub().GetState(orderObj.orderID)
	if orderBytes == nil {
		return fmt.Errorf("cannot find this order")
	}

	order := new(Order)
	_ = json.Unmarshal(orderBytes, order)

	if order.distributorId != user.userId {
		return fmt.Errorf("Permission denied!")
	}
	delivery := DeliveryStatus{
		distributedId: user.userId,
        Status:     "Done delivery to "+ user.address,
        DeliveryDate:  txTimeAsPtr,
	}
	order.deliveryStatus = append(order.deliveryStatus, delivery)
	order.status = orderObj.status
	order.signature = orderObj.signature

	updateOrderAsBytes, _ := json.Marshal(order)

	return ctx.GetStub().PutState(order.orderID, updateOrderAsBytes)
}


// get the history transaction of product
func (s *SmartContract) GetHistory(ctx contractapi.TransactionContextInterface, productId string) ([]ProductHistory, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(productId)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	defer resultsIterator.Close()

	var histories []ProductHistory

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var product Product
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &product)
			if err != nil {
				return nil, err
			}
		} else {
			product = Product{
				productId: productId,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		productHistory := ProductHistory{
			record:    &product,
			txId:      response.TxId,
			timestamp: timestamp,
			isDelete:  response.IsDelete,
		}
		histories = append(histories, productHistory)
	}

	return histories, nil
}
