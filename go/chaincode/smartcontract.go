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
	UserId      string `json:"userId"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FullName    string `json:"fullName"`
	UserName    string `json:"userName"`
	Address     string `json:"address"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	Signature   string `json:"signature"`
}

type ProductDates struct {
	Cultivated     string `json:"cultivated"` // Supplier
	Harvested      string `json:"harvested"`
	Imported       string `json:"imported"` // Manufacturer
	Manufacturered string `json:"manufacturered"`
	Exported       string `json:"exported"`
	Distributed    string `json:"distributed"` // Distributor
	Selling        string `json:"selling"` // Retailer
	Sold           string `json:"sold"` 
}

type ProductActors struct {
	SupplierId     string `json:"supplierId"`
	ManufacturerId string `json:"manufacturerId"`
	DistributorId  string `json:"distributorId"`
	RetailerId     string `json:"retailerId"`
}

type Product struct {
	ProductId      string        `json:"productId"`
	Image          []string      `json:"image" metadata:",optional"`
	ProductName    string        `json:"productName"`
	Dates          ProductDates  `json:"dates"`
	Actors         ProductActors `json:"actors"`
	Expired        string        `json:"expireTime"`
	Price          string        `json:"price"`
	Status         string        `json:"status"`
	Description    string        `json:"description"`
	CertificateUrl string        `json:"certificateUrl"`
	CooperationId  string        `json:"cooperationId"`
	QRCode		   string		 `json:"qrCode"`
}

type ProductHistory struct {
	Record    *Product  `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

type OrderHistory struct {
	Record    *Order    `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

type Signature struct {
	DistributorSignature string `json:"distributorSignature"`
	RetailerSignature    string `json:"retailerSignature"`
}

type ProductItem struct {
	Product  Product `json:"product"`
	Quantity string  `json:"quantity"`
}

type DeliveryStatus struct {
	DistributorId 	string 		`json:"distributorId"`
	DeliveryDate 	string		`json:"deliveryDate"`
	Status       	string    	`json:"status"`
	Longitude		string    	`json:"longitude"`
	Latitude		string    	`json:"latitude"`
}

type Order struct {
	OrderId 		string      	`json:"orderId"`
	ProductItemList []ProductItem 	`json:"productItemList" metadata:",optional"`
	DeliveryStatus 	[]DeliveryStatus `json:"deliveryStatus" metadata:",optional"`
	Signature 		Signature 		`json:"signature"`
	Status          string     	 	`json:"status"`
	DistributorId  	string 			`json:"distributorId"`
	RetailerId     	string 			`json:"retailerId"`
	QRCode		   	string		 	`json:"qrCode"`
	CreateDate 		string 			`json:"createDate"`
	UpdateDate 		string 			`json:"updateDate"`
	FinishDate   	string      	`json:"finishDate"`
}

// Initialize chaincode
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	error := initCounter(ctx)
	if error != nil {
		return fmt.Errorf("error init counter: %s", error.Error())
	}
	return nil
}

func initCounter(ctx contractapi.TransactionContextInterface) error {
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

func getCounter(ctx contractapi.TransactionContextInterface, assetType string) (int, error) {
	counterAsBytes, _ := ctx.GetStub().GetState(assetType)
	counterAsset := CounterNO{}
	json.Unmarshal(counterAsBytes, &counterAsset)
	return counterAsset.Counter, nil
}

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
	return counterAsset.Counter, nil
}

func (s *SmartContract) GetTxTimestampChannel(ctx contractapi.TransactionContextInterface) (string, error) {
	txTimeAsPtr, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		fmt.Printf("Returning error in TimeStamp \n")
		return "Error", err
	}
	timeStr := time.Unix(txTimeAsPtr.Seconds, int64(txTimeAsPtr.Nanos)).String()
	return timeStr, nil
}

// SUPPLIER FUNCTION
func (s *SmartContract) CultivateProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {
	if user.Role != "supplier" {
		return fmt.Errorf("user must be a supplier")
	}
	productCounter, _ := getCounter(ctx, "ProductCounterNO")
	productCounter++

	// get transaction timestamp from channel header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in transaction timeStamp")
	}

	// DATES
	dates := ProductDates{}
	dates.Cultivated = txTimeAsPtr
	actors := ProductActors{}
	actors.SupplierId = user.UserId
	var product = Product{
		ProductId:      "Product" + strconv.Itoa(productCounter),
		ProductName:    productObj.ProductName,
		Image:          productObj.Image,
		Dates:          dates,
		Actors:         actors,
		Expired:        "",
		Price:          productObj.Price,
		Status:         "CULTIVATING",
		Description:    productObj.Description,
		CertificateUrl: productObj.CertificateUrl,
		CooperationId:  productObj.CooperationId,
		QRCode:  		productObj.QRCode,
	}
	productAsBytes, _ := json.Marshal(product)
	incrementCounter(ctx, "ProductCounterNO")

	return ctx.GetStub().PutState(product.ProductId, productAsBytes)
}

func (s *SmartContract) HarvestProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {
	if user.Role != "supplier" {
		return fmt.Errorf("user must be a supplier")
	}

	// get product details from the stub ie. Chaincode stub in network using the product id passed
	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	// get transaction timestamp from channel header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in transaction timeStamp")
	}

	// Updating the product values withe the new values
	product.Dates.Harvested = txTimeAsPtr
	product.Status = "HARVESTED"

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)
}

func (s *SmartContract) SupplierUpdateProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {
	if user.Role != "supplier" {
		return fmt.Errorf("user must be a supplier")
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
	if user.Role != "supplier" {
		return fmt.Errorf("user must be a supplier")
	}

	// get product details from the stub ie. Chaincode stub in network using the product id passed
	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}
	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	product.CertificateUrl = productObj.CertificateUrl

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)

}

// MANUFACTURER
func (s *SmartContract) ImportProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {
	if user.Role != "manufacturer" {
		return fmt.Errorf("user must be a manufacturer")
	}

	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	// get transaction timestamp from channel header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in transaction timeStamp")
	}

	// Updating the product values withe the new values
	// product.image = productObj.image
	product.Dates.Imported = txTimeAsPtr
	product.Price = productObj.Price
	product.Status = "IMPORTED"
	product.Actors.ManufacturerId = user.UserId

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)
}

func (s *SmartContract) ManufactureProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {
	if user.Role != "manufacturer" {
		return fmt.Errorf("user must be a manufacturer")
	}

	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	// get transaction timestamp from channel header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in transaction timeStamp")
	}

	if product.Actors.ManufacturerId != user.UserId {
		return fmt.Errorf("Permission denied!")
	}

	// Updating the product values withe the new values
	product.Image = productObj.Image
	product.Dates.Manufacturered = txTimeAsPtr
	product.Status = "MANUFACTURED"
	product.QRCode = productObj.QRCode
	product.Expired = productObj.Expired

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)
}

func (s *SmartContract) ExportProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {
	if user.Role != "manufacturer" {
		return fmt.Errorf("user must be a manufacturer")
	}

	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	// get transaction timestamp from channel header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in transaction timeStamp")
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
func (s *SmartContract) DistributeProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {
	if user.Role != "distributor" {
		return fmt.Errorf("user must be a distributor")
	}

	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	// get transaction timestamp from channel header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in transaction timeStamp")
	}

	// Updating the product values withe the new values
	// product.Dates.distributed[0].distributorId = user.UserId
	product.Dates.Distributed = txTimeAsPtr
	// product.Dates.distributed[0].Status = "Start delivery"

	product.Status = "DISTRIBUTED"
	product.Actors.DistributorId = user.UserId

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)
}

// RETAILER
func (s *SmartContract) ImportRetailerProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {
	if user.Role != "retailer" {
		return fmt.Errorf("user must be a retailer")
	}

	// get product details from the stub ie. Chaincode stub in network using the product id passed
	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	// get transaction timestamp from channel header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in transaction timeStamp")
	}

	// Updating the product values to be updated after the function
	product.Dates.Selling = txTimeAsPtr
	product.Status = "SELLING"
	product.Price = productObj.Price
	product.Actors.RetailerId = user.UserId

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)
}

func (s *SmartContract) SellProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) error {
	if user.Role != "retailer" {
		return fmt.Errorf("user must be a retailer")
	}

	// get product details from the stub ie. Chaincode stub in network using the product id passed
	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	// get transaction timestamp from channel header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in transaction timeStamp")
	}

	// Updating the product values to be updated after the function
	product.Dates.Sold = txTimeAsPtr
	product.Status = "SOLD"
	product.Price = productObj.Price
	product.Actors.RetailerId = user.UserId

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)
}

func (s *SmartContract) GetProduct(ctx contractapi.TransactionContextInterface, ProductId string) (*Product, error) {
	productAsBytes, err := ctx.GetStub().GetState(ProductId)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	}
	if productAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", ProductId)
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

func (s *SmartContract) GetOrder(ctx contractapi.TransactionContextInterface, OrderId string) (*Order, error) {
	orderAsBytes, err := ctx.GetStub().GetState(OrderId)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	}
	if orderAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", OrderId)
	}

	order := new(Order)
	_ = json.Unmarshal(orderAsBytes, order)

	return order, nil
}

func (s *SmartContract) GetAllOrders(ctx contractapi.TransactionContextInterface) ([]*Order, error) {
	assetCounter, _ := getCounter(ctx, "OrderCounterNO")
	startKey := "Order1"
	endKey := "Order" + strconv.Itoa(assetCounter+1)
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()
	var orders []*Order

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var order Order
		_ = json.Unmarshal(response.Value, &order)
		orders = append(orders, &order)
	}

	return orders, nil
}

func (s *SmartContract) GetAllOrdersByAddress(ctx contractapi.TransactionContextInterface, longitude string, latitude string, shippingStatus string) ([]*Order, error) {
    assetCounter, _ := getCounter(ctx, "OrderCounterNO")
	startKey := "Order1"
	endKey := "Order" + strconv.Itoa(assetCounter+1)
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()
	var orders []*Order

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var order Order
		_ = json.Unmarshal(response.Value, &order)

		if shippingStatus == "" {
			for _, status := range order.DeliveryStatus {
				if status.Longitude == longitude && status.Latitude == latitude {
					orders = append(orders, &order)
					break 
				}
			}
		} else {
			for _, status := range order.DeliveryStatus {
				if status.Longitude == longitude && status.Latitude == latitude && order.Status == shippingStatus {
					orders = append(orders, &order)
					break 
				}
			}
		}
	}

	return orders, nil
}

// manufacturer
func (s *SmartContract) CreateOrder(ctx contractapi.TransactionContextInterface, user User, orderObj Order) error {
	if user.Role != "manufacturer" {
		return fmt.Errorf("user must be a manufacturer")
	}

	orderCounter, _ := getCounter(ctx, "OrderCounterNO")
	orderCounter++

	// get transaction timestamp from channel header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in transaction timeStamp")
	}

	firstdelivery := DeliveryStatus{
		DistributorId: orderObj.DistributorId,
		Status:        "NOT-SHIPPED-YET",
		DeliveryDate:  txTimeAsPtr,
		Longitude: orderObj.DeliveryStatus[0].Longitude,
		Latitude: orderObj.DeliveryStatus[0].Latitude,
	}
	var deliveryStatus []DeliveryStatus

	deliveryStatus = append(deliveryStatus, firstdelivery)

	// DATES
	var order = Order{
		OrderId:   			"Order" + strconv.Itoa(orderCounter),
		ProductItemList: 	orderObj.ProductItemList,
		Signature:       	orderObj.Signature,
		DeliveryStatus:     deliveryStatus,
		Status:     		"NOT-SHIPPED-YET",
		DistributorId: 		orderObj.DistributorId,
		RetailerId: 		orderObj.RetailerId,
		QRCode:				orderObj.QRCode,
		CreateDate: 		txTimeAsPtr,
		UpdateDate: 		"",
		FinishDate: 		"",
	}

	orderAsBytes, _ := json.Marshal(order)
	incrementCounter(ctx, "OrderCounterNO")

	return ctx.GetStub().PutState(order.OrderId, orderAsBytes)
}

// distributor
func (s *SmartContract) UpdateOrder(ctx contractapi.TransactionContextInterface, user User, orderObj Order, longitude string, latitude string) error {
	if user.Role != "distributor" {
		return fmt.Errorf("user must be a distributor")
	}

	// get transaction timestamp from channel header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in transaction timeStamp")
	}

	orderBytes, _ := ctx.GetStub().GetState(orderObj.OrderId)
	if orderBytes == nil {
		return fmt.Errorf("cannot find this order")
	}

	order := new(Order)
	_ = json.Unmarshal(orderBytes, order)

	if order.DistributorId != user.UserId {
		return fmt.Errorf("Permission denied!")
	}

	delivery := DeliveryStatus{
		DistributorId: 	user.UserId,
		Status:        	"SHIPPING",
		DeliveryDate:  	txTimeAsPtr,
		Longitude: 		longitude,
		Latitude: 		latitude,
	}
	order.DeliveryStatus = append(order.DeliveryStatus, delivery)
	order.Status = "SHIPPING"
	order.UpdateDate = txTimeAsPtr
	// order.Signature = orderObj.Signature
	// for i := range order.ProductItemList {
	// 	order.ProductItemList[i].Quantity = orderObj.ProductItemList[i].Quantity
	// }

	updateOrderAsBytes, _ := json.Marshal(order)

	return ctx.GetStub().PutState(order.OrderId, updateOrderAsBytes)
}

// distributor
func (s *SmartContract) FinishOrder(ctx contractapi.TransactionContextInterface, user User, orderObj Order, longitude string, latitude string) error {
	if user.Role != "distributor" {
		return fmt.Errorf("user must be a distributor")
	}

	// get transaction timestamp from channel header
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("returning error in transaction timeStamp")
	}

	orderBytes, _ := ctx.GetStub().GetState(orderObj.OrderId)
	if orderBytes == nil {
		return fmt.Errorf("cannot find this order")
	}

	order := new(Order)
	_ = json.Unmarshal(orderBytes, order)

	if order.DistributorId != user.UserId {
		return fmt.Errorf("Permission denied!")
	}
	delivery := DeliveryStatus{
		DistributorId: 	user.UserId,
		Status:        	"SHIPPED",
		DeliveryDate:  	txTimeAsPtr,		
		Longitude: 		longitude,
		Latitude: 		latitude,
	}

	order.DeliveryStatus = append(order.DeliveryStatus, delivery)
	order.Status = "SHIPPED"
	order.FinishDate = txTimeAsPtr
	// order.Signature = orderObj.Signature

	finishOrderAsBytes, _ := json.Marshal(order)

	return ctx.GetStub().PutState(order.OrderId, finishOrderAsBytes)
}

func (s *SmartContract) GetProductTransactionHistory(ctx contractapi.TransactionContextInterface, productId string) ([]ProductHistory, error) {
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

func (s *SmartContract) GetOrderTransactionHistory(ctx contractapi.TransactionContextInterface, orderId string) ([]OrderHistory, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(orderId)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	defer resultsIterator.Close()
	var histories []OrderHistory

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var order Order
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &order)
			if err != nil {
				return nil, err
			}
		} else {
			order = Order{
				OrderId: orderId,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		orderHistory := OrderHistory{
			Record:    &order,
			TxId:      response.TxId,
			Timestamp: timestamp,
			IsDelete:  response.IsDelete,
		}
		histories = append(histories, orderHistory)
	}

	return histories, nil
}
