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
	UserId   string `json:"UserId"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
	UserName string `json:"UserName"`
	Address  string `json:"Address"`
	UserType string `json:"UserType"`
	Role     string `json:"Role"`
	Status   string `json:"Status"`
}

type ProductDates struct {
	Cultivated     string `json:"Cultivated"` // supplier
	Harvested      string `json:"Harvested"`
	Imported       string `json:"Imported"` // manufacturer
	Manufacturered string `json:"Manufacturered"`
	Exported       string `json:"Exported"`
	Distributed    string `json:"Distributed"` // distributor
	Sold           string `json:"Sold"`        // retailer
}

type ProductActors struct {
	SupplierId     string `json:"SupplierId"`
	ManufacturerId string `json:"ManufacturerId"`
	DistributorId  string `json:"DistributorId"`
	RetailerId     string `json:"RetailerId"`
}

// Unit: kg, box/boxes, bottle, bottles

// Supplier: id, cultivate, harvest => cultivating, harvested
// Manufacturer: id, import, manufacture, export => imported, manufacturing, exported
// Distributor: id, distribute => distributed/distributing
// Retailer: id, sell => sold
type Product struct {
	ProductId   string        `json:"ProductId"`
	ProductName string        `json:"ProductName"`
	Dates       ProductDates  `json:"Dates"`
	Actors      ProductActors `json:"Actors"`
	Price       float64       `json:"Price"`
	Status      string        `json:"Status"`
	Description string        `json:"Description"`
}

type ProductHistory struct {
	Record    *Product  `json:"Record"`
	TxId      string    `json:"TxId"`
	Timestamp time.Time `json:"Timestamp"`
	IsDelete  bool      `json:"IsDelete"`
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

	// Initializing User Counter
	UserCounterBytes, _ := ctx.GetStub().GetState("UserCounterNO")
	if UserCounterBytes == nil {
		var UserCounter = CounterNO{Counter: 0}
		UserCounterBytes, _ := json.Marshal(UserCounter)
		err := ctx.GetStub().PutState("UserCounterNO", UserCounterBytes)
		if err != nil {
			return fmt.Errorf("failed to Intitate User Counter: %s", err.Error())
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

// sign in
func (s *SmartContract) SignIn(ctx contractapi.TransactionContextInterface, email string, password string) (*User, error) {

	results, err := s.GetAllUsers(ctx)

	if err != nil {
		return nil, err
	}

	for _, user := range results {
		_email := (*user).Email
		_password := (*user).Password
		_userId := (*user).UserId
		if _email == email && _password == password {
			userBytes, _ := ctx.GetStub().GetState(_userId)
			_user := new(User)
			_ = json.Unmarshal(userBytes, _user)
			return _user, nil
		}

	}
	return nil, fmt.Errorf("user is not exists")
}

// create user
func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, email string, password string, username string, address string, userType string, role string) error {

	results, err := s.GetAllUsers(ctx)

	if err != nil {
		return err
	}

	for _, user := range results {
		_email := (*user).Email
		if _email == email {
			return fmt.Errorf("this user is exists")
		}
	}

	userCounter, _ := getCounter(ctx, "UserCounterNO")
	userCounter++

	user := User{
		UserId:   "User" + strconv.Itoa(userCounter),
		Email:    email,
		Password: password,
		UserName: username,
		Address:  address,
		UserType: userType,
		Role:     role,
	}

	userAsBytes, errMarshal := json.Marshal(user)
	if errMarshal != nil {
		return fmt.Errorf("marshal Error in Product: %s", errMarshal)
	}

	incrementCounter(ctx, "UserCounterNO")

	return ctx.GetStub().PutState(user.UserId, userAsBytes)
}

// SUPPLIER FUNCTION
// cultivate product // gieo trồng sảm phẩm
func (s *SmartContract) CultivateProduct(ctx contractapi.TransactionContextInterface, userId string, productName string, price float64, description string) error {

	// get user details from the stub ie. Chaincode stub in network using the user id passed
	userBytes, _ := ctx.GetStub().GetState(userId)
	if userBytes == nil {
		return fmt.Errorf("cannot find User")
	}

	user := new(User)
	_ = json.Unmarshal(userBytes, user)

	// User type check for the function
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
	actors.SupplierId = userId
	var product = Product{
		ProductId:   "Product" + strconv.Itoa(productCounter),
		ProductName: productName,
		Dates:       dates,
		Actors:      actors,
		Price:       price,
		Status:      "CULTIVATING",
		Description: description,
	}

	productAsBytes, _ := json.Marshal(product)

	incrementCounter(ctx, "ProductCounterNO")

	return ctx.GetStub().PutState(product.ProductId, productAsBytes)
}

// havert product // thu hoạch
func (s *SmartContract) HarvertProduct(ctx contractapi.TransactionContextInterface, userId string, productId string) error {

	// get user details from the stub ie. Chaincode stub in network using the user id passed
	userBytes, _ := ctx.GetStub().GetState(userId)
	if userBytes == nil {
		return fmt.Errorf("cannot find this user")
	}

	user := new(User)
	_ = json.Unmarshal(userBytes, user)

	if user.UserType != "supplier" {
		return fmt.Errorf("User must be a supplier")
	}

	// get product details from the stub ie. Chaincode stub in network using the product id passed
	productBytes, _ := ctx.GetStub().GetState(productId)
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
func (s *SmartContract) SupplierUpdateProduct(ctx contractapi.TransactionContextInterface, userId string, productId string, productName string, price float64, description string) error {

	// get user details from the stub ie. Chaincode stub in network using the user id passed
	userBytes, _ := ctx.GetStub().GetState(userId)
	if userBytes == nil {
		return fmt.Errorf("cannot find this user")
	}

	user := new(User)
	_ = json.Unmarshal(userBytes, user)

	if user.UserType != "supplier" {
		return fmt.Errorf("User must be a supplier")
	}

	// get product details from the stub ie. Chaincode stub in network using the product id passed
	productBytes, _ := ctx.GetStub().GetState(productId)
	if productBytes == nil {
		return fmt.Errorf("cannot find this product")
	}

	product := new(Product)
	_ = json.Unmarshal(productBytes, product)

	// Updating the product values withe the new values
	product.ProductName = productName
	product.Price = price
	product.Description = description

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)
}

// MANUFACTURER
// import product
func (s *SmartContract) ImportProduct(ctx contractapi.TransactionContextInterface, userId string, productId string, price float64) error {

	userBytes, _ := ctx.GetStub().GetState(userId)
	if userBytes == nil {
		return fmt.Errorf("cannot find User")
	}

	user := new(User)
	_ = json.Unmarshal(userBytes, user)

	// User type check for the function
	if user.UserType != "manufacturer" {
		return fmt.Errorf("User must be a manufacturer")
	}

	productBytes, _ := ctx.GetStub().GetState(productId)
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
	product.Price = price
	product.Status = "IMPORTED"
	product.Actors.ManufacturerId = userId

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)
}

// manufacture product
func (s *SmartContract) ManufactureProduct(ctx contractapi.TransactionContextInterface, userId string, productId string) error {

	userBytes, _ := ctx.GetStub().GetState(userId)
	if userBytes == nil {
		return fmt.Errorf("cannot find this user")
	}

	user := new(User)
	_ = json.Unmarshal(userBytes, user)

	// User type check for the function
	if user.UserType != "manufacturer" {
		return fmt.Errorf("User must be a manufacturer")
	}

	productBytes, _ := ctx.GetStub().GetState(productId)
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

	if product.Actors.ManufacturerId != userId {
		return fmt.Errorf("Permission denied!")
	}

	// Updating the product values withe the new values
	product.Dates.Manufacturered = txTimeAsPtr
	product.Status = "MANUFACTURED"

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)
}

// export product
func (s *SmartContract) ExportProduct(ctx contractapi.TransactionContextInterface, userId string, productId string, price float64) error {

	userBytes, _ := ctx.GetStub().GetState(userId)
	if userBytes == nil {
		return fmt.Errorf("cannot find this user")
	}

	user := new(User)
	_ = json.Unmarshal(userBytes, user)

	// User type check for the function
	if user.UserType != "manufacturer" {
		return fmt.Errorf("User must be a manufacturer")
	}

	productBytes, _ := ctx.GetStub().GetState(productId)
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

	if product.Actors.ManufacturerId != userId {
		return fmt.Errorf("Permission denied!")
	}

	// Updating the product values withe the new values
	product.Dates.Exported = txTimeAsPtr
	product.Price = price
	product.Status = "EXPORTED"

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)
}

// DISTRIBUTOR
// distribute product
func (s *SmartContract) DistributeProduct(ctx contractapi.TransactionContextInterface, userId string, productId string) error {

	userBytes, _ := ctx.GetStub().GetState(userId)
	if userBytes == nil {
		return fmt.Errorf("cannot find this user")
	}

	user := new(User)
	_ = json.Unmarshal(userBytes, user)

	// User type check for the function
	if user.UserType != "distributor" {
		return fmt.Errorf("User must be a distributor")
	}

	productBytes, _ := ctx.GetStub().GetState(productId)
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
	product.Dates.Distributed = txTimeAsPtr
	product.Status = "DISTRIBUTED"
	product.Actors.DistributorId = userId

	updatedProductAsBytes, _ := json.Marshal(product)

	return ctx.GetStub().PutState(product.ProductId, updatedProductAsBytes)
}

// RETAILER
// sell product
func (s *SmartContract) SellProduct(ctx contractapi.TransactionContextInterface, userId string, productId string, price float64) error {

	userBytes, _ := ctx.GetStub().GetState(userId)
	if userBytes == nil {
		return fmt.Errorf("cannot find this user")
	}

	user := new(User)
	_ = json.Unmarshal(userBytes, user)

	// User type check for the function
	if user.UserType != "retailer" {
		return fmt.Errorf("user must be a retailer")
	}

	// get product details from the stub ie. Chaincode stub in network using the product id passed
	productBytes, _ := ctx.GetStub().GetState(productId)
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
	product.Price = price
	product.Actors.RetailerId = userId

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

// get all asset
func (s *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*User, error) {
	assetCounter, _ := getCounter(ctx, "UserCounterNO")
	startKey := "User1"
	endKey := "User" + strconv.Itoa(assetCounter+1)
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	var users []*User

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var user User
		_ = json.Unmarshal(response.Value, &user)

		users = append(users, &user)
	}
	return users, nil
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

	return products, nil
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
