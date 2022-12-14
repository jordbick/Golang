package product

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/jordbick/Golang/inventory-service/database"
)

// keep all of our data access separate from the web service code
// Allow us to easily replace the implementations of these methods once we start working with a DB

func getProduct(productID int) (*Product, error) {
	// if query takes longer than 15s will cancel and return
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Get single rwo from DB so can use QueryRow
	row := database.DbConn.QueryRowContext(ctx, `SELECT productId,
	manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products
	WHERE productId = ?`, productID)

	product := &Product{}
	err := row.Scan(&product.ProductID,
		&product.Manufacturer,
		&product.Sku,
		&product.Upc,
		&product.PricePerUnit,
		&product.QuantityOnHand,
		&product.ProductName)

	// If no rows, return nil as record doesn't exist
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return product, nil
}

// DELETE
func removeProduct(productID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM products where productId = ?`, productID)
	if err != nil {
		return err
	}
	return nil
}

// GET ALL
// Convert into SELECT statements to query the DB rather than static data
// Change function to return an error as when we're working with a DB there could be a connection problem
func getProductList() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// Use DB Query method as we are returning a list
	// use the DB connection that we created in our main function
	results, err := database.DbConn.QueryContext(ctx, `SELECT productId,
	manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products`)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	// Previously we were using a struct with a mutex and a map to manage our products,
	// But now can just return a slice of products (bevause we are getting our products straight from the DB instead of from memory)
	// To do this we grab the Rows object that comes back from the Query method and using a for loop we can use the Next method to move the cursor to the next method
	// Within the loop call the Scan method and pass in the field names for the struct we want to map the column names to
	// Have to be in same order as the SELECT statement
	products := make([]Product, 0)
	for results.Next() {
		var product Product
		results.Scan(&product.ProductID,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.PricePerUnit,
			&product.QuantityOnHand,
			&product.ProductName)
		// append this object product to our slice of products,
		products = append(products, product)
	}
	return products, nil
}

// Similar to getProductList but restricting number of products back
func GetTopTenProducts() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT 
	productId, 
	manufacturer, 
	sku, 
	upc, 
	pricePerUnit, 
	quantityOnHand, 
	productName 
	FROM products ORDER BY quantityOnHand DESC LIMIT 10
	`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	products := make([]Product, 0)
	for results.Next() {
		var product Product
		results.Scan(&product.ProductID,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.PricePerUnit,
			&product.QuantityOnHand,
			&product.ProductName)

		products = append(products, product)
	}
	return products, nil
}

func updateProduct(product Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// Call to Exec method
	_, err := database.DbConn.ExecContext(ctx, `UPDATE products SET 
	manufacturer=?,
	sku=?,
	upc=?,
	pricePerUnit=CAST(? AS DECIMAL(13,2)),
	quantityOnHand=?,
	productName=?
	WHERE productId=?`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName,
		product.ProductID,
	)
	if err != nil {
		return err
	}
	return nil
}

func insertProduct(product Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result, err := database.DbConn.ExecContext(ctx, `INSERT into products (
	manufacturer,
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName) VALUES (?, ?, ?, ?, ?, ?)`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName)

	if err != nil {
		return 0, nil
	}
	// want to know the last product ID of the record that was inserted
	insertID, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(insertID), nil
}

// Building out the WHERE clause using the fields in our productFilter
func searchForProductData(productFilter ProductReportFilter) ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var queryArgs = make([]interface{}, 0)
	var queryBuilder strings.Builder
	// appending to our query string programmatically
	queryBuilder.WriteString(`SELECT 
		productId, 
		LOWER(manufacturer), 
		LOWER(sku), 
		upc, 
		pricePerUnit, 
		quantityOnHand, 
		LOWER(productName) 
		FROM products WHERE `)
	if productFilter.NameFilter != "" {
		queryBuilder.WriteString(`productName LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(productFilter.NameFilter)+"%")
	}
	if productFilter.ManufacturerFilter != "" {
		if len(queryArgs) > 0 {
			queryBuilder.WriteString(" AND ")
		}
		queryBuilder.WriteString(`manufacturer LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(productFilter.ManufacturerFilter)+"%")
	}
	if productFilter.SKUFilter != "" {
		if len(queryArgs) > 0 {
			queryBuilder.WriteString(" AND ")
		}
		queryBuilder.WriteString(`sku LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(productFilter.SKUFilter)+"%")
	}

	results, err := database.DbConn.QueryContext(ctx, queryBuilder.String(), queryArgs...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	products := make([]Product, 0)
	// scan over results we get back and return those to a slice, and return it to the caller
	for results.Next() {
		var product Product
		results.Scan(&product.ProductID,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.PricePerUnit,
			&product.QuantityOnHand,
			&product.ProductName)

		products = append(products, product)
	}
	return products, nil
}

// func addOrUpdateProduct(product Product) (int, error) {
// 	// if the product id is set, update, otherwise add
// 	addOrUpdateID := -1
// 	if product.ProductID > 0 {
// 		oldProduct, err := getProduct(product.ProductID)
// 		if err != nil {
// 			return addOrUpdateID, err
// 		}
// 		// if it exists, replace it, otherwise return error
// 		if oldProduct == nil {
// 			return 0, fmt.Errorf("product id [%d] doesn't exist", product.ProductID)
// 		}
// 		addOrUpdateID = product.ProductID
// 	} else {
// 		addOrUpdateID = getNextProductID()
// 		product.ProductID = addOrUpdateID
// 	}
// 	productMap.Lock()
// 	productMap.m[addOrUpdateID] = product
// 	productMap.Unlock()
// 	return addOrUpdateID, nil
// }

// // Function to iterate over products and get highest product ID value
// // Function to get list of IDs sorted in ascending order and then a helper function to wrap that, which will give us the next highest ID value
// func getProductIds() []int {
// 	productMap.RLock()
// 	productIds := []int{}
// 	for key := range productMap.m {
// 		productIds = append(productIds, key)
// 	}
// 	productMap.RUnlock()
// 	sort.Ints(productIds)
// 	return productIds
// }

// func getNextProductID() int {
// 	productIDs := getProductIds()
// 	return productIDs[len(productIDs)-1] + 1
// }

// GET ALL
// func getProductList() []Product {
// 	productMap.RLock()
// 	products := make([]Product, 0, len(productMap.m))
// 	for _, value := range productMap.m {
// 		products = append(products, value)
// 	}
// 	productMap.RUnlock()
// 	return products
// }

// // DELETE
// func removeProduct(productID int) {
// 	productMap.Lock()
// 	defer productMap.Unlock()
// 	delete(productMap.m, productID)
// }

// replace the slice of products with a new type which includes a map of ints to products and a read/write mutex
// By using a map we can use the ProductID as the key for the map and the Product as the value
// Allows us to acces products without having to iterate over all the items in our slice
// Need to use the mutex bevause our web services are multi-threaded, and maps in Go are inherently not thread safe
// Need to wrap mour map using a mutex to avoid 2 threads from writing and reading at the same time

// var productMap = struct {
// 	sync.RWMutex
// 	m map[int]Product
// }{m: make(map[int]Product)}

// // To use this new struct we need to modify our functions that interact with the list of products
// // Call our function to load the JSON file (containing our data)
// func init() {
// 	fmt.Println("loading products...")
// 	prodMap, err := loadProductMap()
// 	productMap.m = prodMap
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("%d products loaded...\n", len(productMap.m))
// }

// // returns a map or error
// func loadProductMap() (map[int]Product, error) {
// 	fileName := "products.json"
// 	// checks that the file exists
// 	_, err := os.Stat(fileName)
// 	if os.IsNotExist(err) {
// 		return nil, fmt.Errorf("file [%s] does not exist", fileName)
// 	}

// 	// ReadFile to read all the data in the file into a byte slice
// 	file, _ := ioutil.ReadFile(fileName)
// 	productList := make([]Product, 0)
// 	// json.Unmarshal to derserialise the bytes into a slice of products
// 	err = json.Unmarshal([]byte(file), &productList)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	prodMap := make(map[int]Product)
// 	// iterate over the slice to initialise the items in our map
// 	for i := 0; i < len(productList); i++ {
// 		prodMap[productList[i].ProductID] = productList[i]
// 	}
// 	return prodMap, nil
// }

// // Functions to work with our product map type
// // GET BY ID
// func getProduct(productID int) *Product {
// 	// Calling the RLock function to get a read lock from the mutex
// 	// Prevents another thread from getting a write lock on the struct while we're reading from the map
// 	productMap.RLock()
// 	// Then calling defer RUnlock to release the read lock from the mutex
// 	defer productMap.RUnlock()
// 	if product, ok := productMap.m[productID]; ok {
// 		return &(product)
// 	}
// 	return nil
// }
