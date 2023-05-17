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
	Counter int `json:"counter"`
}

type User struct {
	UserId   string `json:"userId"`
	Email    string `json:"email"`
	Password string `json:"password"`
	UserName string `json:"userName"`
	Address  string `json:"address"`
	UserType string `json:"userType"`
	Role     string `json:"role"`
	Status   string `json:"status"`
	Identify string `json:"identify"`
}


type ProductDates struct {
	Cultivated     string `json:"cultivated"` // supplier
	Harvested      string `json:"harvested"`
	Imported       string `json:"imported"` // manufacturer
	Manufacturered string `json:"manufacturered"`
	Exported       string `json:"exported"`
	Distributed    string `json:"distributed"` // distributor
	Sold           string `json:"sold"`        // retailer
}

type ProductActors struct {
	SupplierId     string `json:"supplierId"`
	ManufacturerId string `json:"manufacturerId"`
	DistributorId  string `json:"distributorId"`
	RetailerId     string `json:"retailerId"`
}

// Unit: kg, box/boxes, bottle, bottles

// Supplier: id, cultivate, harvest => cultivating, harvested
// Manufacturer: id, import, manufacture, export => imported, manufacturing, exported
// Distributor: id, distribute => distributed/distributing
// Retailer: id, sell => sold
type Product struct {
	ProductId   	string        `json:"productId"`
	Image 			[]string	  `json:"image"`
	ProductName 	string        `json:"productName"`
	Dates       	ProductDates  `json:"dates"`
	Actors      	ProductActors `json:"actors"`
	Price       	string        `json:"price"`
	Status      	string        `json:"status"`
	Description 	string        `json:"description"`
	CertificateURL 	string 		  `json:"certificate"`
	CooperativeId 	string 		  `json:"cooperativeId"`
}

type ProductHistory struct {
	Record    *Product  `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

// order

type Signature struct {
	DistributorSignature  	string 	`json:"distributorSignature"`
	RetailerSignature 		string  `json:"retailerSignature"`
}


type ProductItem struct {
	Product  Product `json:"product"`
	Quantity string  `json:"quantity"`
}

type DeliveryStatus struct {
	DistributedId 	string 		`json:"distributedId"`
	DeliveryDate 	string		`json:"deliveryDate"`
	Status       	string    	`json:"status"`
}

type Order struct {
	OrderID 		string      	`json:"orderID"`
	ProductItemList []ProductItem 	`json:"productItemList"`
	Signature 		Signature 		`json:"signature"`
	// DateCreate 		string 			`json:"dateCreate"`
	// DateFinish      string      	`json:"dateFinish"`
	DeliveryStatus 	[]DeliveryStatus `json:"deliveryStatus"`
	Status          string     	 	`json:"status"`
	DistributorId  	string 			`json:"distributorId"`
	RetailerId     	string 			`json:"retailerId"`
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
		var ProductCounter = CounterNO{Counter: 0}
		ProductCounterBytes, _ := json.Marshal(ProductCounter)
		err := ctx.GetStub().PutState("ProductCounterNO", ProductCounterBytes)
		if err != nil {
			return fmt.Errorf("failed to Intitate Product Counter: %s", err.Error())
		}
	}
	OrderCounterBytes, _ := ctx.GetStub().GetState("OrdertCounterNO")
	if OrderCounterBytes == nil {
		var OrderCounter = CounterNO{Counter: 0}
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
	return counterAsset.Counter, nil
}

// incrementCounter to the increase value of the counter based on the Asset Type provided as input parameter by 1
func incrementCounter(ctx contractapi.TransactionContextInterface, assetType string) (int, error) {
	counterAsBytes, _ := ctx.GetStub().GetState(assetType)
	counterAsset := CounterNO{}

	json.Unmarshal(counterAsBytes, &counterAsset)
	counterAsset.Counter++
	counterAsBytes, _ = json.Marshal(counterAsset)

	err := ctx.GetStub().PutState(assetType, counterAsBytes)
	if err != nil {
		return -1, fmt.Errorf("failed to Increment Counter: %s", err.Error())
	}
	fmt.Printf("Printf in incrementing counter  %v", counterAsset)
	return counterAsset.Counter, nil
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

	if user.UserType != "supplier" {
		return fmt.Errorf("User must be a supplier")
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
	dates.Cultivated = txTimeAsPtr
	actors := ProductActors{}
	actors.SupplierId = user.UserId
	var product = Product{
		ProductId:   "Product" + strconv.Itoa(productCounter),
		ProductName: productObj.ProductName,
		Dates:       dates,
		Actors:      actors,
		Price:       productObj.Price,
		Status:      "CULTIVATING",
		Description: productObj.Description,
		CooperativeId : productObj.CooperativeId,
	}

	productAsBytes, _ := json.Marshal(product)

	incrementCounter(ctx, "ProductCounterNO")

	return ctx.GetStub().PutState(product.ProductId, productAsBytes)
}

// havert product // thu hoạch
func (s *SmartContract) HarvertProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {

	if user.UserType != "supplier" {
		return fmt.Errorf("User must be a supplier")
	}

	// get product details from the stub ie. Chaincode stub in network using the product id passed
	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
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
	product.Dates.Harvested = txTimeAsPtr
	product.Status = "HAVERTED"

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)
}

// supplier update
func (s *SmartContract) SupplierUpdateProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {

	if user.UserType != "supplier" {
		return fmt.Errorf("User must be a supplier")
	}

	// get product details from the stub ie. Chaincode stub in network using the product id passed
	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	// Updating the product values withe the new values
	product.ProductName = productObj.ProductName
	product.Price = productObj.Price
	product.Description = productObj.Description

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)
}
func (s *SmartContract) AddCertificate(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {

	if user.UserType != "supplier" {
		return fmt.Errorf("User must be a supplier")
	}

	// get product details from the stub ie. Chaincode stub in network using the product id passed
	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}
	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	product.CertificateURL = productObj.CertificateURL
	
	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)

}
// MANUFACTURER
// import product
func (s *SmartContract) ImportProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {

	if user.UserType != "manufacturer" {
		return fmt.Errorf("User must be a manufacturer")
	}

	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
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
	product.Dates.Imported = txTimeAsPtr
	product.Price = productObj.Price
	product.Status = "IMPORTED"
	product.Actors.ManufacturerId = user.UserId

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)
}

// manufacture product
func (s *SmartContract) ManufactureProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {

	if user.UserType != "manufacturer" {
		return fmt.Errorf("User must be a manufacturer")
	}

	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
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

	if product.Actors.ManufacturerId != user.UserId {
		return fmt.Errorf("Permission denied!")
	}

	// Updating the product values withe the new values
	product.Dates.Manufacturered = txTimeAsPtr
	product.Status = "MANUFACTURED"

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)
}

// export product
func (s *SmartContract) ExportProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {

	if user.UserType != "manufacturer" {
		return fmt.Errorf("User must be a manufacturer")
	}

	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
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

	if product.Actors.ManufacturerId != user.UserId {
		return fmt.Errorf("Permission denied!")
	}

	// Updating the product values withe the new values
	product.Dates.Exported = txTimeAsPtr
	product.Price = productObj.Price
	product.Status = "EXPORTED"

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)
}

// DISTRIBUTOR
// distribute product
func (s *SmartContract) DistributeProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {

	if user.UserType != "distributor" {
		return fmt.Errorf("User must be a distributor")
	}

	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
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
	// product.Dates.Distributed[0].DistributedId = user.UserId
	product.Dates.Distributed = txTimeAsPtr
	// product.Dates.Distributed[0].Status = "Start delivery"

	product.Status = "DISTRIBUTED"
	product.Actors.DistributorId = user.UserId

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)
}

// RETAILER
// sell product
func (s *SmartContract) SellProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {

	if user.UserType != "retailer" {
		return fmt.Errorf("user must be a retailer")
	}

	// get product details from the stub ie. Chaincode stub in network using the product id passed
	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
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
	product.Dates.Sold = txTimeAsPtr
	product.Status = "SOLD"
	product.Price = productObj.Price
	product.Actors.RetailerId = user.UserId

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)
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
	if user.UserType != "distributor" {
		return fmt.Errorf("User must be a distributor")
	}

	orderCounter, _ := getCounter(ctx, "OrderCounterNO")
	orderCounter++

	//To Get the transaction TimeStamp from the Channel Header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in Transaction TimeStamp") 
	}

	firstdelivery := DeliveryStatus{
		DistributedId: user.UserId,
        Status:     "Start delivery",
        DeliveryDate:  txTimeAsPtr,
	}
	var deliveryStatus []DeliveryStatus

	deliveryStatus = append(deliveryStatus, firstdelivery)

	// DATES
	var order = Order{
		OrderID:   			"Order" + strconv.Itoa(orderCounter),
		ProductItemList: 	orderObj.ProductItemList,
		Signature:       	orderObj.Signature,
		DeliveryStatus:     deliveryStatus,
		Status:     		orderObj.Status,
		DistributorId: 		user.UserId,
		RetailerId: 		orderObj.RetailerId,
	}

	orderAsBytes, _ := json.Marshal(order)

	incrementCounter(ctx, "OrderCounterNO")

	return ctx.GetStub().PutState(order.OrderID, orderAsBytes)
}

func (s *SmartContract) updateOrder(ctx contractapi.TransactionContextInterface,user User,orderObj Order ) error {

	if user.UserType != "distributor" {
		return fmt.Errorf("User must be a distributor")
	}

	//To Get the transaction TimeStamp from the Channel Header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in Transaction TimeStamp") 
	}

	orderBytes, _ := ctx.GetStub().GetState(orderObj.OrderID)
	if orderBytes == nil {
		return fmt.Errorf("cannot find this order")
	}

	order := new(Order)
	_ = json.Unmarshal(orderBytes, order)

	if order.DistributorId != user.UserId {
		return fmt.Errorf("Permission denied!")
	}
	delivery := DeliveryStatus{
		DistributedId: user.UserId,
        Status:     "Delivering "+ user.Address,
        DeliveryDate:  txTimeAsPtr,
	}
	order.DeliveryStatus = append(order.DeliveryStatus, delivery)
	order.Status = orderObj.Status
	// order.Signature = orderObj.Signature

	updateOrderAsBytes, _ := json.Marshal(order)

	return ctx.GetStub().PutState(order.OrderID, updateOrderAsBytes)
}

func (s *SmartContract) finishOrder(ctx contractapi.TransactionContextInterface,user User,orderObj Order ) error {

	if user.UserType != "retailer" {
		return fmt.Errorf("User must be a retailer")
	}

	//To Get the transaction TimeStamp from the Channel Header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in Transaction TimeStamp") 
	}

	orderBytes, _ := ctx.GetStub().GetState(orderObj.OrderID)
	if orderBytes == nil {
		return fmt.Errorf("cannot find this order")
	}

	order := new(Order)
	_ = json.Unmarshal(orderBytes, order)

	if order.DistributorId != user.UserId {
		return fmt.Errorf("Permission denied!")
	}
	delivery := DeliveryStatus{
		DistributedId: user.UserId,
        Status:     "Done delivery to "+ user.Address,
        DeliveryDate:  txTimeAsPtr,
	}
	order.DeliveryStatus = append(order.DeliveryStatus, delivery)
	order.Status = orderObj.Status
	order.Signature = orderObj.Signature

	updateOrderAsBytes, _ := json.Marshal(order)

	return ctx.GetStub().PutState(order.OrderID, updateOrderAsBytes)
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
				ProductId: productId,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		productHistory := ProductHistory{
			Record:    &product,
			TxId:      response.TxId,
			Timestamp: timestamp,
			IsDelete:  response.IsDelete,
		}
		histories = append(histories, productHistory)
	}

	return histories, nil
}
