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
	Avatar     	string `json:"avatar"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	Signature   string `json:"signature"`
	Cart		[]ProductIdItem `json:"cart" metadata:",optional"`
}

type UserPayload struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	Role   		string `json:"role"`
	Email       string `json:"email"`
	FullName    string `json:"fullName"`
	Address     string `json:"address"`
	Avatar     	string `json:"avatar"`
	Signature   string `json:"signature"`
}

type Actor struct {
	UserId      string `json:"userId"`
	PhoneNumber string `json:"phoneNumber"`
	FullName    string `json:"fullName"`
	Address     string `json:"address"`
	Avatar     	string `json:"avatar"`
	Role        string `json:"role"`
}

type ProductActors struct {
	SupplierId     string `json:"supplierId"`
	ManufacturerId string `json:"manufacturerId"`
	DistributorId  string `json:"distributorId"`
	RetailerId     string `json:"retailerId"`
}

type ProductDate struct {
	Status     	string 	 `json:"status"`
	Time 		string	 `json:"time"`
	Actor  		Actor 	 `json:"actor"`
}

type Product struct {
	ProductId      string         `json:"productId"`
	ProductName    string         `json:"productName"`
	Supplier 	   Actor          `json:"supplier"`
	Dates          []ProductDate  `json:"dates" metadata:",optional"`
	Image          []string       `json:"image" metadata:",optional"`
	Expired        string         `json:"expireTime"`
	Price          string         `json:"price"`
	Amount         string         `json:"amount"`
	Unit           string         `json:"unit"`
	Status         string         `json:"status"`
	Description    string         `json:"description"`
	CertificateUrl string         `json:"certificateUrl"`
	QRCode		   string		  `json:"qrCode"`
}

type ProductPayload struct {
	ProductName    string        `json:"productName"`
	Image          []string      `json:"image" metadata:",optional"`
	Price          string        `json:"price"`
	Amount         string        `json:"amount"`
	Unit           string        `json:"unit"`
	Description    string        `json:"description"`
	CertificateUrl string        `json:"certificateUrl"`
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

type ProductItem struct {
	Product  Product `json:"product"`
	Quantity string  `json:"quantity"`
}

type ProductIdItem struct {
	ProductId  	string 	`json:"productId"`
	Quantity 	string  `json:"quantity"`
}

type ProductItemPayload struct {
	ProductId  	string 	`json:"productId"`
	Quantity 	string  `json:"quantity"`
}

type DeliveryStatus struct {
	Status       	string    	`json:"status"`
	DeliveryDate 	string		`json:"deliveryDate"`
	Address			string    	`json:"address"`
	Actor 			Actor 		`json:"actor"`
}

type DeliveryStatusCreateOrder struct {
	Address			string    	`json:"address"`
}

type Order struct {
	OrderId 		string      	 `json:"orderId"`
	ProductItemList []ProductItem 	 `json:"productItemList" metadata:",optional"`
	DeliveryStatuses[]DeliveryStatus `json:"deliveryStatuses" metadata:",optional"`
	Signatures 		[]string 		 `json:"signatures"`
	Status          string     	 	 `json:"status"`
	CreateDate 		string 			 `json:"createDate"`
	UpdateDate 		string 			 `json:"updateDate"`
	FinishDate   	string      	 `json:"finishDate"`
	QRCode		   	string		 	 `json:"qrCode"`
	Retailer     	Actor 			 `json:"retailer"`
	Manufacturer  	Actor 			 `json:"manufacturer"`
	Distributor  	Actor 			 `json:"distributor"`
}

type OrderPayload struct {
	ProductIdItems 	[]ProductIdItem 			`json:"productIdItems" metadata:",optional"`
	DeliveryStatus 	DeliveryStatusCreateOrder 	`json:"deliveryStatus"`
	Signatures 		[]string 					`json:"signatures"`
	QRCode		   	string		 				`json:"qrCode"`
}

func parseUserToActor(user User) Actor {
	actor := Actor{
		UserId:user.UserId,
		PhoneNumber:user.PhoneNumber,
		FullName:user.FullName,
		Address:user.Address,
		Avatar:user.Avatar,
		Role:user.Role,
	}
	return actor
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

func (s *SmartContract) GetCounterOfType(ctx contractapi.TransactionContextInterface, assetType string) (int, error) {
	counterAsBytes, _ := ctx.GetStub().GetState(assetType)
	counterAsset := CounterNO{}
	json.Unmarshal(counterAsBytes, &counterAsset)
	return counterAsset.Counter, nil
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

// supplier
func (s *SmartContract) CultivateProduct(ctx contractapi.TransactionContextInterface, user User, productObj ProductPayload) (*Product, error) {
	if user.Role != "supplier" {
		return nil, fmt.Errorf("user must be a supplier")
	}

	productCounter, _ := getCounter(ctx, "ProductCounterNO")
	productCounter++

	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return nil, fmt.Errorf("transaction timeStamp error")
	}

	actor := parseUserToActor(user)
	var datesArray []ProductDate
	date := ProductDate{
		Status: "CULTIVATED",
		Time: txTimeAsPtr,
		Actor: actor,
	}
	// date1 := ProductDate{
	// 	Status: "HARVESTED",
	// 	Time: txTimeAsPtr,
	// 	Actor: actor,
	// }
	// date2 := ProductDate{
	// 	Status: "IMPORTED",
	// 	Time: txTimeAsPtr,
	// 	Actor: actor,
	// }
	// date3 := ProductDate{
	// 	Status: "MANUFACTURED",
	// 	Time: txTimeAsPtr,
	// 	Actor: actor,
	// }
	// dates := append(datesArray, date, date1, date2, date3)

	dates := append(datesArray, date)
	
	var product = Product{
		ProductId:      "Product" + strconv.Itoa(productCounter),
		ProductName:    productObj.ProductName,
		Image:          productObj.Image,
		Dates:          dates,
		Price:          productObj.Price,
		Amount:         productObj.Amount,
		Unit:         	productObj.Unit,
		Status:         "CULTIVATING",
		Description:    productObj.Description,
		CertificateUrl: productObj.CertificateUrl,
		Supplier:  		actor,
	}
	productAsBytes, _ := json.Marshal(product)
	incrementCounter(ctx, "ProductCounterNO")

	ctx.GetStub().PutState(product.ProductId, productAsBytes)

	return &product, nil
}

func (s *SmartContract) InventoryProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) (*Product, error) {
	if user.Role != "manufacturer" {
		return nil, fmt.Errorf("user must be a manufacturer")
	}

	productCounter, _ := getCounter(ctx, "ProductCounterNO")
	productCounter++

	actor := parseUserToActor(user)
	
	var product = Product{
		ProductId:      "Product" + strconv.Itoa(productCounter),
		ProductName:    productObj.ProductName,
		Image:          productObj.Image,
		Dates:          productObj.Dates,
		Expired:        productObj.Expired,
		Price:          productObj.Price,
		Amount:         productObj.Amount,
		Unit:         	productObj.Unit,
		Status:         "MANUFACTURED",
		Description:    productObj.Description,
		CertificateUrl: productObj.CertificateUrl,
		QRCode:  		productObj.QRCode,
		Supplier:  		actor,
	}
	productAsBytes, _ := json.Marshal(product)
	incrementCounter(ctx, "ProductCounterNO")

	ctx.GetStub().PutState(product.ProductId, productAsBytes)

	return &product, nil
}

// Need refactor
func (s *SmartContract) HarvestProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) (*Product, error) {
	if user.Role != "supplier" {
		return nil, fmt.Errorf("user must be a supplier")
	}

	// get product details
	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return nil, fmt.Errorf("product not found")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return nil, fmt.Errorf("transaction timeStamp error")
	}

	actor := parseUserToActor(user)
	date := ProductDate{
		Status: "HARVESTED",
		Time: txTimeAsPtr,
		Actor: actor,
	}
	dates := append(product.Dates, date)

	// update product
	product.Dates = dates
	product.Status = "HARVESTED"
	product.Amount = productObj.Amount

	updatedProductAsBytes, _ := json.Marshal(product)
	ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)

	return product, nil
}

// supplier
func (s *SmartContract) SupplierUpdateProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) (*Product, error) {
	if user.Role != "supplier" {
		return nil, fmt.Errorf("user must be a supplier")
	}

	// get product
	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return nil, fmt.Errorf("product not found")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	// update product
	product.ProductName = productObj.ProductName
	product.Description = productObj.Description
	product.Price = productObj.Price
	product.Amount = productObj.Amount

	updatedProductAsBytes, _ := json.Marshal(product)
	ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)

	return product, nil
}

func (s *SmartContract) UpdateProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) (*Product, error) {
	if user.Role != "supplier" && user.Role != "manufacturer" && user.Role != "distributor" && user.Role != "retailer" {
		return nil, fmt.Errorf("user must be a supplier, manufacturer, distributor, or retailer")
	}

	// get product
	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return nil, fmt.Errorf("product not found")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	// update product
	product = &productObj
	updatedProductAsBytes, _ := json.Marshal(product)
	ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)

	return product, nil
}

func (s *SmartContract) AddCertificate(ctx contractapi.TransactionContextInterface, user User, productObj Product) (*Product, error) {
	if user.Role != "supplier" {
		return nil, fmt.Errorf("user must be a supplier")
	}

	// get product
	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return nil, fmt.Errorf("product not found")
	}
	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	product.CertificateUrl = productObj.CertificateUrl
	updatedProductAsBytes, _ := json.Marshal(product)
	ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)

	return product, nil
}

// manufacturer
func (s *SmartContract) ImportProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) (*Product, error) {
	if user.Role != "manufacturer" {
		return nil, fmt.Errorf("user must be a manufacturer")
	}

	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return nil, fmt.Errorf("product not found")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return nil, fmt.Errorf("transaction timeStamp error")
	}

	actor := parseUserToActor(user)
	date := ProductDate{
		Status: "IMPORTED",
		Time: txTimeAsPtr,
		Actor: actor,
	}
	dates := append(productObj.Dates, date)

	// update product
	product.Dates = dates
	product.Image = productObj.Image
	product.Price = productObj.Price
	product.Status = "IMPORTED"

	updatedProductAsBytes, _ := json.Marshal(product)
	ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)

	return product, nil
}

func (s *SmartContract) ManufactureProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) (*Product, error) {
	if user.Role != "manufacturer" {
		return nil, fmt.Errorf("user must be a manufacturer")
	}

	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return nil, fmt.Errorf("product not found")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return nil, fmt.Errorf("transaction timeStamp error")
	}

	if product.Dates[2].Actor.UserId != user.UserId {
		return nil, fmt.Errorf("Permission denied!")
	}

	actor := parseUserToActor(user)
	date := ProductDate{
		Status: "MANUFACTURED",
		Time: txTimeAsPtr,
		Actor: actor,
	}
	dates := append(product.Dates, date)

	// update product
	product.Dates = dates
	product.Image = productObj.Image
	product.QRCode = productObj.QRCode
	product.Expired = productObj.Expired
	product.Status = "MANUFACTURED"

	updatedProductAsBytes, _ := json.Marshal(product)
	ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)

	return product, nil
}

func (s *SmartContract) ExportProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) (*Product, error) {
	if user.Role != "manufacturer" {
		return nil, fmt.Errorf("user must be a manufacturer")
	}

	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return nil, fmt.Errorf("product not found")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return nil, fmt.Errorf("transaction timeStamp error")
	}

	if product.Dates[3].Actor.UserId != user.UserId {
		return nil, fmt.Errorf("Permission denied!")
	}

	actor := parseUserToActor(user)
	date := ProductDate{
		Status: "EXPORTED",
		Time: txTimeAsPtr,
		Actor: actor,
	}
	dates := append(product.Dates, date)

	// update product
	product.Dates = dates
	product.Price = productObj.Price
	product.Status = "EXPORTED"

	updatedProductAsBytes, _ := json.Marshal(product)
	ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)

	return product, nil
}

// distributor
func (s *SmartContract) DistributeProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) (*Product, error) {
	if user.Role != "distributor" {
		return nil, fmt.Errorf("user must be a distributor")
	}

	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return nil, fmt.Errorf("product not found")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return nil, fmt.Errorf("transaction timeStamp error")
	}

	actor := parseUserToActor(user)
	date := ProductDate{
		Status: "DISTRIBUTED",
		Time: txTimeAsPtr,
		Actor: actor,
	}
	dates := append(product.Dates, date)

	// update product
	product.Dates = dates
	product.Status = "DISTRIBUTED"

	updatedProductAsBytes, _ := json.Marshal(product)
	ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)

	return product, nil
}

// retailer
func (s *SmartContract) ImportRetailerProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) (*Product, error) {
	if user.Role != "retailer" {
		return nil, fmt.Errorf("user must be a retailer")
	}

	// get product details from the stub ie. Chaincode stub in network using the product id passed
	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return nil, fmt.Errorf("product not found")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return nil, fmt.Errorf("transaction timeStamp error")
	}

	actor := parseUserToActor(user)
	date := ProductDate{
		Status: "SELLING",
		Time: txTimeAsPtr,
		Actor: actor,
	}
	dates := append(product.Dates, date)

	// update product
	product.Dates = dates
	product.Price = productObj.Price
	product.Status = "SELLING"

	updatedProductAsBytes, _ := json.Marshal(product)
	ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)

	return product, nil
}

func (s *SmartContract) SellProduct(ctx contractapi.TransactionContextInterface, user User, productObj Product) (*Product, error) {
	if user.Role != "retailer" {
		return nil, fmt.Errorf("user must be a retailer")
	}

	// get product details from the stub ie. Chaincode stub in network using the product id passed
	productBytes, _ := ctx.GetStub().GetState(productObj.ProductId)
	if productBytes == nil {
		return nil, fmt.Errorf("product not found")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return nil, fmt.Errorf("transaction timeStamp error")
	}

	actor := parseUserToActor(user)
	date := ProductDate{
		Status: "SOLD",
		Time: txTimeAsPtr,
		Actor: actor,
	}
	dates := append(product.Dates, date)

	// update product
	product.Dates = dates
	product.Price = productObj.Price
	product.Status = "SOLD"

	updatedProductAsBytes, _ := json.Marshal(product)
	ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)

	return product, nil
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

func (s *SmartContract) GetAllProducts(ctx contractapi.TransactionContextInterface) ([]*Product, error) {
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

	if len(products) == 0 {
		return []*Product{}, nil
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

func (s *SmartContract) GetAllOrders(ctx contractapi.TransactionContextInterface, status string) ([]*Order, error) {
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

		if status == "" || order.Status == status {
			orders = append(orders, &order)
		}
	}

	if len(orders) == 0 {
		return []*Order{}, nil
	}

	return orders, nil
}

// func (s *SmartContract) GetAllOrdersByAddress(ctx contractapi.TransactionContextInterface, longitude string, latitude string, shippingStatus string) ([]*Order, error) {
    // assetCounter, _ := getCounter(ctx, "OrderCounterNO")
	// startKey := "Order1"
	// endKey := "Order" + strconv.Itoa(assetCounter+1)
	// resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	// if err != nil {
	// 	return nil, err
	// }

	// defer resultsIterator.Close()
	// var orders []*Order

	// for resultsIterator.HasNext() {
	// 	response, err := resultsIterator.Next()

	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	var order Order
	// 	_ = json.Unmarshal(response.Value, &order)

	// 	if shippingStatus == "" {
	// 		for _, status := range order.DeliveryStatus {
	// 			if status.Longitude == longitude && status.Latitude == latitude {
	// 				orders = append(orders, &order)
	// 				break 
	// 			}
	// 		}
	// 	} else {
	// 		for _, status := range order.DeliveryStatus {
	// 			if status.Longitude == longitude && status.Latitude == latitude && order.Status == shippingStatus {
	// 				orders = append(orders, &order)
	// 				break 
	// 			}
	// 		}
	// 	}
	// }

	// if len(orders) == 0 {
	// 	return []*Order{}, nil
	// }

	// return orders, nil
// }

func (s *SmartContract) GetAllOrdersOfManufacturer(ctx contractapi.TransactionContextInterface, userId string, status string) ([]*Order, error) {
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

		if order.Manufacturer.UserId == userId && status == "" || order.Status == status {
			orders = append(orders, &order)
		}
	}

	if len(orders) == 0 {
		return []*Order{}, nil
	}

	return orders, nil
}

func (s *SmartContract) GetAllOrdersOfDistributor(ctx contractapi.TransactionContextInterface, userId string, status string) ([]*Order, error) {
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

		if order.Distributor.UserId == userId && status == "" || order.Status == status {
			orders = append(orders, &order)
		}
	}

	if len(orders) == 0 {
		return []*Order{}, nil
	}

	return orders, nil
}

func (s *SmartContract) GetAllOrdersOfRetailer(ctx contractapi.TransactionContextInterface, userId string, status string) ([]*Order, error) {
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

		if order.Retailer.UserId == userId && status == "" || order.Status == status {
			orders = append(orders, &order)
		}
	}

	if len(orders) == 0 {
		return []*Order{}, nil
	}

	return orders, nil
}

// retailer
func (s *SmartContract) CreateOrder(ctx contractapi.TransactionContextInterface, user User, orderObj OrderPayload) (*Order, error) {
	if user.Role != "retailer" {
		return nil, fmt.Errorf("user must be a retailer")
	}

	orderCounter, _ := getCounter(ctx, "OrderCounterNO")
	orderCounter++

	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return nil, fmt.Errorf("transaction timeStamp error")
	}

	actor := parseUserToActor(user)

	delivery := DeliveryStatus{
		Status:        	"PENDING",
		DeliveryDate:  	txTimeAsPtr,
		Address: 		orderObj.DeliveryStatus.Address,
		Actor: 			actor,
	}
	var deliveryStatuses []DeliveryStatus
	deliveryStatuses = append(deliveryStatuses, delivery)

	var productItemList []ProductItem
	for _, item := range orderObj.ProductIdItems {
		productAsBytes, err := ctx.GetStub().GetState(item.ProductId)
		if err != nil {
			return nil, fmt.Errorf("product not found")
		}

		product := new(Product)
		_ = json.Unmarshal(productAsBytes, product)

		productItem := ProductItem{ 
			Product: *product, 
			Quantity: item.Quantity, 
		}

		productItemList = append(productItemList, productItem)
	}

	var order = Order{
		OrderId:   			"Order" + strconv.Itoa(orderCounter),
		ProductItemList: 	productItemList,
		Signatures:       	orderObj.Signatures,
		DeliveryStatuses:   deliveryStatuses,
		Status:     		"PENDING",
		Manufacturer:		actor,
		Distributor: 		actor,
		Retailer: 			actor,
		QRCode:				orderObj.QRCode,
		CreateDate: 		txTimeAsPtr,
		UpdateDate: 		"",
		FinishDate: 		"",
	}

	orderAsBytes, _ := json.Marshal(order)
	incrementCounter(ctx, "OrderCounterNO")
	ctx.GetStub().PutState(order.OrderId, orderAsBytes)

	return &order, nil
}

// manufacturer
func (s *SmartContract) ApproveOrder(ctx contractapi.TransactionContextInterface, user User, orderId string) (*Order, error) {
	if user.Role != "manufacturer" {
		return nil, fmt.Errorf("user must be a manufacturer")
	}

	orderAsBytes, err := ctx.GetStub().GetState(orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	}
	if orderAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", orderId)
	}

	order := new(Order)
	_ = json.Unmarshal(orderAsBytes, order)

	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return nil, fmt.Errorf("transaction timeStamp error")
	}

	// export products in order
	var productItemList []ProductItem
	for _, item := range order.ProductItemList {
		actor := parseUserToActor(user)
		date := ProductDate{
			Status: "EXPORTED",
			Time: txTimeAsPtr,
			Actor: actor,
		}
		dates := append(item.Product.Dates, date)

		// update product in chaincode
		item.Product.Dates = dates
		item.Product.Price = item.Product.Price
		item.Product.Status = "EXPORTED"

		updatedProductAsBytes, _ := json.Marshal(item.Product)
		ctx.GetStub().PutState(item.Product.ProductId, updatedProductAsBytes)

		// update updated products into order
		productItem := ProductItem{
			Product: item.Product,
			Quantity: item.Quantity,
		}
		productItemList = append(productItemList, productItem)
	}

	// if order.Distributor.UserId != user.UserId {
	// 	return nil, fmt.Errorf("Permission denied!")
	// }

	actor := parseUserToActor(user)
	delivery := DeliveryStatus{
		Status:        	"APPROVED",
		DeliveryDate:  	txTimeAsPtr,
		Address: 		"",
		Actor: 			actor,
	}
	deliveryStatuses := append(order.DeliveryStatuses, delivery)

	order.ProductItemList = productItemList
	order.DeliveryStatuses = deliveryStatuses
	order.UpdateDate = txTimeAsPtr
	order.Status = "APPROVED"

	updateOrderAsBytes, _ := json.Marshal(order)
	ctx.GetStub().PutState(order.OrderId, updateOrderAsBytes)

	return order, nil
}

// distributor
func (s *SmartContract) UpdateOrder(ctx contractapi.TransactionContextInterface, user User, orderObj Order) (*Order, error) {
	if user.Role != "distributor" {
		return nil, fmt.Errorf("user must be a distributor")
	}

	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return nil, fmt.Errorf("transaction timeStamp error")
	}

	orderBytes, _ := ctx.GetStub().GetState(orderObj.OrderId)
	if orderBytes == nil {
		return nil, fmt.Errorf("cannot find this order")
	}

	order := new(Order)
	_ = json.Unmarshal(orderBytes, order)

	// if order.Distributor.UserId != user.UserId {
	// 	return nil, fmt.Errorf("Permission denied!")
	// }

	actor := parseUserToActor(user)
	delivery := DeliveryStatus{
		Status:        	"SHIPPING",
		DeliveryDate:  	txTimeAsPtr,
		Address: 		orderObj.DeliveryStatuses[0].Address,
		Actor: 			actor,
	}
	deliveryStatuses := append(order.DeliveryStatuses, delivery)

	order.DeliveryStatuses = deliveryStatuses
	order.UpdateDate = txTimeAsPtr
	order.Status = "SHIPPING"

	updateOrderAsBytes, _ := json.Marshal(order)
	ctx.GetStub().PutState(order.OrderId, updateOrderAsBytes)

	return order, nil
}

// distributor
func (s *SmartContract) FinishOrder(ctx contractapi.TransactionContextInterface, user User, orderObj Order) (*Order, error) {
	if user.Role != "distributor" {
		return nil, fmt.Errorf("user must be a distributor")
	}

	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return nil, fmt.Errorf("transaction timeStamp error")
	}

	orderBytes, _ := ctx.GetStub().GetState(orderObj.OrderId)
	if orderBytes == nil {
		return nil, fmt.Errorf("cannot find this order")
	}

	order := new(Order)
	_ = json.Unmarshal(orderBytes, order)

	// if order.Distributor.UserId != user.UserId {
	// 	return nil, fmt.Errorf("Permission denied!")
	// }
	
	actor := parseUserToActor(user)
	delivery := DeliveryStatus{
		Status:        	"SHIPPED",
		DeliveryDate:  	txTimeAsPtr,
		Address: 		orderObj.DeliveryStatuses[0].Address,
		Actor: 			actor,
	}
	deliveryStatuses := append(order.DeliveryStatuses, delivery)

	order.DeliveryStatuses = deliveryStatuses
	order.FinishDate = txTimeAsPtr
	order.Status = "SHIPPED"
	order.Signatures = orderObj.Signatures

	finishOrderAsBytes, _ := json.Marshal(order)
	ctx.GetStub().PutState(order.OrderId, finishOrderAsBytes)

	return order, nil
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

	if len(histories) == 0 {
		return []ProductHistory{}, nil
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

	if len(histories) == 0 {
		return []OrderHistory{}, nil
	}

	return histories, nil
}
