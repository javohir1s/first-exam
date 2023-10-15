package vazifalar

import (
	"encoding/json"
	"log"
	"os"
	"playground/newProject/config"
	"playground/newProject/handler"
	"playground/newProject/models"
	"playground/newProject/pkg"
	"playground/newProject/storage"
	"time"

	"github.com/spf13/cast"

	"sort"
)

type DailyProductStats struct {
	Date   string
	Amount int
}
type UserPrTransaction struct {
	UserID      string
	ProductsIn  int
	ProductsOut int
}

type UserDailyTransaction struct {
	UserID           string
	Date             string
	TransactionCount int
	TotalSum         int
}
type UserTransactionSummary struct {
	UserID   string
	UserName string
	TotalSum int
}

type ProductByDay struct {
	CreatedAt time.Time
	Quantity  int
}
type Vazifalar struct {
	strg storage.StorageI
	cfg  config.Config
}
type top struct {
	Count int
	Name  string
}
type branch struct {
	Sum        int
	BranchName string
}
type eachBranch struct {
	Categories []category
	BranchName string
}
type category struct {
	CategoryName            string
	ProductTransactionCount int
}
type Product struct {
	ID       string
	Quantity int
	Name     string
}
type plus_minus struct {
	BranchName string
	Plus       int
	Minus      int
	Sum        struct {
		Plus  int
		Minus int
	}
}
type Branch_product_transaction struct {
	Plus  int
	Minus int
}
type Sum struct {
	Plus  int
	Minus int
}

type BranchTranz struct {
	ProductId       string
	Quantity        int
	TransactionType string
}
type BranchProductValue struct {
	BranchName        string
	TotalProductValue int
}

var fileNames = models.FileNames{
	BranchFile:                   "../data/branches.json",
	UserFile:                     "../data/users.json",
	CategoryFile:                 "../data/categories.json",
	ProductFile:                  "../data/products.json",
	BranchProductFile:            "../data/branch_products.json",
	BranchProductTransactionFile: "../data/branch_pr_transaction.json",
}

func NewVazifalarController(strg storage.StorageI, cfg config.Config) *Vazifalar {
	return &Vazifalar{
		strg: strg,
		cfg:  cfg,
	}
}

// twelve task

func (v *Vazifalar) GetUserProductTransactions() map[string]UserPrTransaction {
	dataTransactions, err := pkg.Read(fileNames.BranchProductTransactionFile)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	userTransactions := make(map[string]UserPrTransaction)

	for _, val := range dataTransactions {
		transaction := cast.ToStringMap(val)
		userID := cast.ToString(transaction["user_id"])
		quantity := cast.ToInt(transaction["quantity"])
		transactionType := cast.ToString(transaction["transaction_type"])

		userTransaction, exists := userTransactions[userID]
		if !exists {
			userTransaction = UserPrTransaction{
				UserID: userID,
			}
		}

		if transactionType == "plus" {
			userTransaction.ProductsIn += quantity
		} else {
			userTransaction.ProductsOut += quantity
		}

		userTransactions[userID] = userTransaction
	}

	return userTransactions
}

// eleventh task
func (v *Vazifalar) GetUserDailyTransactionSummaries() []UserDailyTransaction {
	dataTransactions, err := pkg.Read(fileNames.BranchProductTransactionFile)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	h := handler.NewHandler(v.strg, v.cfg)

	userSummaries := make(map[string]map[string]UserDailyTransaction)

	for _, val := range dataTransactions {
		transaction := cast.ToStringMap(val)
		userID := cast.ToString(transaction["user_id"])
		transactionDateStr := cast.ToString(transaction["created_at"])
		quantity := cast.ToInt(transaction["quantity"])
		productID := cast.ToString(transaction["product_id"])

		product := h.GetProduct(productID)
		if product.Price != 0 {
			price := product.Price * quantity
			transactionDate, dateErr := time.Parse("2006-01-02 15:04:05", transactionDateStr)
			if dateErr != nil {
				log.Println("Error parsing date:", dateErr)
				continue
			}
			dayKey := transactionDate.Format("2006-01-02")

			if userSummaries[userID] == nil {
				userSummaries[userID] = make(map[string]UserDailyTransaction)
			}

			dailyTransaction, exists := userSummaries[userID][dayKey]
			if !exists {
				dailyTransaction = UserDailyTransaction{
					UserID:           userID,
					Date:             dayKey,
					TransactionCount: 0,
					TotalSum:         0,
				}
			}

			dailyTransaction.TransactionCount++
			dailyTransaction.TotalSum += price
			userSummaries[userID][dayKey] = dailyTransaction
		}
	}

	var result []UserDailyTransaction
	for _, dailyTransactions := range userSummaries {
		for _, dailyTransaction := range dailyTransactions {
			result = append(result, dailyTransaction)
		}
	}

	return result
}

// tenth task
func (v *Vazifalar) GetUserTransactionSummaries() []UserTransactionSummary {
	userTransactionSummaries := make(map[string]UserTransactionSummary)

	dataTransactions, err := pkg.Read(fileNames.BranchProductTransactionFile)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	h := handler.NewHandler(v.strg, v.cfg)

	for _, val := range dataTransactions {
		transaction := cast.ToStringMap(val)
		userName := cast.ToString(transaction["user_name"])
		userID := cast.ToString(transaction["user_id"])
		quantity := cast.ToInt(transaction["quantity"])
		productID := cast.ToString(transaction["product_id"])

		product := h.GetProduct(productID)
		if product.Price != 0 {
			price := product.Price * quantity
			userTransactionSummaries[userID] = UserTransactionSummary{
				UserID:   userID,
				UserName: userName,
				TotalSum: userTransactionSummaries[userID].TotalSum + price,
			}
		}
	}

	var result []UserTransactionSummary
	for _, summary := range userTransactionSummaries {
		result = append(result, summary)
	}
	sort.Slice(result, func(a, b int) bool {
		return result[a].TotalSum > result[b].TotalSum
	})

	return result
}

// ninth task
func (v *Vazifalar) GetBranchProductValues() ([]BranchProductValue, error) {
	dataTransactions, err := pkg.Read(fileNames.BranchProductTransactionFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	branchProductSum := make(map[string]int)

	h := handler.NewHandler(v.strg, v.cfg)

	for _, val := range dataTransactions {
		transaction := cast.ToStringMap(val)
		branchID := cast.ToString(transaction["branch_id"])
		quantity := cast.ToInt(transaction["quantity"])
		productID := cast.ToString(transaction["product_id"])

		product := h.GetProduct(productID)
		if product.Price != 0 {
			price := product.Price * quantity
			branchProductSum[branchID] += price
		}
	}

	var result []BranchProductValue
	for branchID, sum := range branchProductSum {
		branch := h.GetBranch(branchID)
		result = append(result, BranchProductValue{
			BranchName:        branch.Name,
			TotalProductValue: sum,
		})
	}
	sort.Slice(result, func(a, b int) bool {
		return result[a].TotalProductValue > result[b].TotalProductValue
	})

	return result, nil
}

// eights task
func (v *Vazifalar) GetProductTransactionSummary() (map[string]map[string]int, error) {
	branchProductTransactions, err := v.read()
	if err != nil {
		return nil, err
	}

	productTransactionSummary := make(map[string]map[string]int)

	productData, err := v.readProductData("../data/products.json")
	if err != nil {
		return nil, err
	}

	for _, product := range productData {
		productTransactionSummary[product.Name] = map[string]int{
			"Kiritilgan":  0,
			"Chiqarilgan": 0,
		}
	}

	for _, transaction := range branchProductTransactions {
		productData, err := v.readProductData("../data/products.json")
		if err != nil {
			return nil, err
		}

		for _, product := range productData {
			if product.Id == transaction.ProductID {
				productTransactionSummary[product.Name]["Kiritilgan"] += transaction.Quantity

			} else {
				productTransactionSummary[product.Name]["Chiqarilgan"] += transaction.Quantity
			}
		}
	}

	return productTransactionSummary, nil
}

func (v *Vazifalar) readProductData(fileName string) ([]models.Product, error) {
	var products []models.Product

	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Printf("Error while reading data: %+v", err)
		return nil, err
	}

	err = json.Unmarshal(data, &products)
	if err != nil {
		log.Printf("Error while unmarshaling data: %+v", err)
		return nil, err
	}

	return products, nil
}

// seventh task
func (v *Vazifalar) GetDailyProductStats() ([]DailyProductStats, error) {
	dataTransactions, err := pkg.Read(fileNames.BranchProductTransactionFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	dailyStats := make(map[string]int)

	for _, val := range dataTransactions {
		transaction := cast.ToStringMap(val)
		createdAt := cast.ToString(transaction["created_at"])
		quantity := cast.ToInt(transaction["quantity"])

		// Распарсим дату в формат "2006-01-02"
		transactionDate, dateErr := time.Parse("2006-01-02 15:04:05", createdAt)
		if dateErr != nil {
			log.Println("Error parsing date:", dateErr)
			continue
		}
		dayKey := transactionDate.Format("2006-01-02")

		dailyStats[dayKey] += quantity
	}

	// Преобразуем карту в срез
	var result []DailyProductStats
	for date, amount := range dailyStats {
		result = append(result, DailyProductStats{Date: date, Amount: amount})
	}

	// Сортируем по дате (как строки)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date < result[j].Date
	})

	return result, nil
}

// sixth task
func (v *Vazifalar) CalculateBranchTransactionsStats() []plus_minus {
	h := handler.NewHandler(v.strg, v.cfg)
	dataTransactions, err := pkg.Read(fileNames.BranchProductTransactionFile)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	branchStats := make(map[string]plus_minus)

	for _, val := range dataTransactions {
		transaction := cast.ToStringMap(val)
		branchID := cast.ToString(transaction["branch_id"])
		transactionType := cast.ToString(transaction["transaction_type"])
		quantity := cast.ToInt(transaction["quantity"])
		productID := cast.ToString(transaction["product_id"])

		stats, exists := branchStats[branchID]
		if !exists {
			stats = plus_minus{
				BranchName: h.GetBranch(branchID).Name,
			}
		}

		product := h.GetProduct(productID)
		price := product.Price
		value := price * quantity

		if transactionType == "plus" {
			stats.Plus += quantity
			stats.Sum.Plus += value
		} else {
			stats.Minus += quantity
			stats.Sum.Minus += value
		}
		branchStats[branchID] = stats
	}

	var result []plus_minus
	for _, stats := range branchStats {
		result = append(result, stats)
	}

	return result
}

// fifth task
func (v *Vazifalar) FindCounterCategoriesInTransactions() []eachBranch {
	h := handler.NewHandler(v.strg, v.cfg)

	dataTransactions, err := pkg.Read(fileNames.BranchProductTransactionFile)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	branchProductMap := make(map[string]map[string]int)

	for _, transaction := range dataTransactions {
		transactionMap, ok := transaction.(map[string]interface{})
		if !ok {
			continue
		}

		branchID, ok := transactionMap["branch_id"].(string)
		if !ok {
			continue
		}

		productID, ok := transactionMap["product_id"].(string)
		if !ok {
			continue
		}

		product := h.GetProduct(productID)
		categoryID := product.CategoryId

		if _, exists := branchProductMap[branchID]; !exists {
			branchProductMap[branchID] = make(map[string]int)
		}

		branchProductMap[branchID][categoryID]++
	}

	var result []eachBranch

	for branchID, categoryCounts := range branchProductMap {
		categories := make([]category, 0, len(categoryCounts))

		for categoryID, count := range categoryCounts {
			categoryName := h.GetCategory(categoryID).Name
			categories = append(categories, category{
				CategoryName:            categoryName,
				ProductTransactionCount: count,
			})
		}

		branchName := h.GetBranch(branchID).Name
		result = append(result, eachBranch{
			BranchName: branchName,
			Categories: categories,
		})
	}

	return result
}

// fouoth task
func (v *Vazifalar) FindTopCategoriesInTransactions() []top {
	h := handler.NewHandler(v.strg, v.cfg)

	dataTransactions, err := pkg.Read(fileNames.BranchProductTransactionFile)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	categoryTransactionCount := make(map[string]int)

	for _, transaction := range dataTransactions {
		transactionMap, ok := transaction.(map[string]interface{})
		if !ok {
			continue
		}

		productID, ok := transactionMap["product_id"].(string)
		if !ok {
			continue
		}

		categoryID := h.GetProduct(productID).CategoryId
		categoryTransactionCount[categoryID]++
	}

	topCategories := make([]top, 0, len(categoryTransactionCount))

	for categoryID, count := range categoryTransactionCount {
		categoryName := h.GetCategory(categoryID).Name
		topCategories = append(topCategories, top{
			Count: count,
			Name:  categoryName,
		})
	}

	sort.Slice(topCategories, func(a, b int) bool {
		return topCategories[a].Count > topCategories[b].Count
	})

	return topCategories
}

// third task
func (v *Vazifalar) FindTopProductsInTransactions() []top {
	h := handler.NewHandler(v.strg, v.cfg)

	dataTransactions, err := pkg.Read(fileNames.BranchProductTransactionFile)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	productTransactionCount := make(map[string]int)
	topProducts := []top{}

	for _, transaction := range dataTransactions {
		transactionMap, ok := transaction.(map[string]interface{})
		if !ok {
			continue
		}

		productID, ok := transactionMap["product_id"].(string)
		if !ok {
			continue
		}

		productTransactionCount[productID]++
	}

	for productID, count := range productTransactionCount {
		productName := h.GetProduct(productID).Name
		topProducts = append(topProducts, top{
			Count: count,
			Name:  productName,
		})
	}

	sort.Slice(topProducts, func(a, b int) bool {
		return topProducts[a].Count > topProducts[b].Count
	})

	return topProducts
}

// second task
func (v *Vazifalar) FindBranchesByTransactionSum() []branch {
	h := handler.NewHandler(v.strg, v.cfg)

	dataTransactions, err := pkg.Read(fileNames.BranchProductTransactionFile)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	branchTransactionSum := make(map[string]int)
	branches := map[string]branch{}

	for _, transaction := range dataTransactions {
		transactionMap, ok := transaction.(map[string]interface{})
		if !ok {
			continue
		}

		branchID := cast.ToString(transactionMap["branch_id"])
		productID := cast.ToString(transactionMap["product_id"])
		quantity := cast.ToInt(transactionMap["quantity"])
		product := h.GetProduct(productID)
		if product.Price != 0 {
			price := int(quantity) * product.Price
			branchTransactionSum[branchID] += price
		}
	}

	for branchID, sum := range branchTransactionSum {
		branchs := h.GetBranch(branchID)
		if branchs.Name != "" {
			branches[branchID] = branch{
				Sum:        sum,
				BranchName: branchs.Name,
			}
		}
	}

	var result []branch
	for _, b := range branches {
		result = append(result, b)
	}

	sort.Slice(result, func(a, b int) bool {
		return result[a].Sum > result[b].Sum
	})

	return result
}

// frist task
func (v *Vazifalar) FindBranchesByTransactionCount() []top {
	h := handler.NewHandler(v.strg, v.cfg)

	dataTransaction, err := pkg.Read(fileNames.BranchProductTransactionFile)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	branchTransactionCount := make(map[string]int)

	for _, val := range dataTransaction {
		Map, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		branchID, ok := Map["branch_id"].(string)
		if !ok {
			continue
		}
		branchTransactionCount[branchID]++
	}

	tranzaction := make([]top, 0, len(branchTransactionCount))

	for branchID, count := range branchTransactionCount {
		branch := h.GetBranch(branchID)
		tranzaction = append(tranzaction, top{
			Count: count,
			Name:  branch.Name,
		})
	}

	sort.Slice(tranzaction, func(a, b int) bool {
		return tranzaction[a].Count > tranzaction[b].Count
	})

	return tranzaction
}

func (v *Vazifalar) read() ([]models.BranchProductTransaction, error) {
	var branchProductTransactions []models.BranchProductTransaction

	data, err := os.ReadFile(fileNames.BranchProductTransactionFile)
	if err != nil {
		log.Printf("Error while reading data: %+v", err)
		return nil, err
	}

	err = json.Unmarshal(data, &branchProductTransactions)
	if err != nil {
		log.Printf("Error while unmarshaling data: %+v", err)
		return nil, err
	}

	return branchProductTransactions, nil
}
