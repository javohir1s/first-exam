package main

import (
	"fmt"
	"playground/newProject/config"
	"playground/newProject/models"
	"playground/newProject/storage/memory"
	"playground/newProject/vazifalar"
)

func main() {

	cfg := config.Load()
	strg := memory.NewStorage(models.FileNames{
		BranchFile:                   "../data/branches.json",
		UserFile:                     "../data/users.json",
		CategoryFile:                 "../data/categories.json",
		ProductFile:                  "../data/products.json",
		BranchProductFile:            "../data/branch_products.json",
		BranchProductTransactionFile: "../data/branch_pr_transaction.json",
	})
	// h := handler.NewHandler(strg, *cfg)
	// h.CreateCategory("FastFood")
	// h.GetAllBranchProductTransaction(1, 10, )
	// h.CreateBranchProductTransaction("f7497e5f-ce02-452b-8ea4-d5aab1c69a1c", "b9e54f3a-7e11-4b0a-94c3-b7e27e601044", "904d377a-75e4-4dd0-a80d-9449e7b91677", "plus", 2)
	// Tasks
	vazifalar := vazifalar.NewVazifalarController(strg, *cfg)
	// firsr task
	topResults := vazifalar.FindBranchesByTransactionCount()
	fmt.Println("First Task - Top Branches by Transaction Count:")
	for i, result := range topResults {
		fmt.Printf("%d. Branch: %s, Transaction Count: %d\n", i+1, result.Name, result.Count)
	}
	// second task
	branchResults := vazifalar.FindBranchesByTransactionSum()
	fmt.Println("Second Task - Branches by Transaction Sum:")
	for i, result := range branchResults {
		fmt.Printf("%d. Branch: %s, Transaction Sum: %d\n", i+1, result.BranchName, result.Sum)
	}
	// third task
	topProductsResults := vazifalar.FindTopProductsInTransactions()
	fmt.Println("Third Task - Top Products by Transaction Count:")
	for i, result := range topProductsResults {
		fmt.Printf("%d. Product: %s, Transaction Count: %d\n", i+1, result.Name, result.Count)
	}

	// fourth task
	topCategoriesResults := vazifalar.FindTopCategoriesInTransactions()
	fmt.Println("Fourth Task - Top Categories by Transaction Count:")
	for i, result := range topCategoriesResults {
		fmt.Printf("%d. Category: %s, Transaction Count: %d\n", i+1, result.Name, result.Count)
	}
	// fifth task
	counterCategoriesResults := vazifalar.FindCounterCategoriesInTransactions()
	fmt.Println("Fifth Task - Categories Count in Transactions by Branch:")
	for i, result := range counterCategoriesResults {
		fmt.Printf("%d. Branch: %s\n", i+1, result.BranchName)
		for j, categoryResult := range result.Categories {
			fmt.Printf("   %d. Category: %s, Transaction Count: %d\n", j+1, categoryResult.CategoryName, categoryResult.ProductTransactionCount)
		}
	}

	// sixth task

	sixthTaskResults := vazifalar.CalculateBranchTransactionsStats()
	fmt.Println("Sixth Task - Branch Transaction Stats:")
	for i, result := range sixthTaskResults {
		fmt.Printf("%d. Branch: %s, Plus: %d, Minus: %d, Sum Plus: %d, Sum Minus: %d\n",
			i+1, result.BranchName, result.Plus, result.Minus, result.Sum.Plus, result.Sum.Minus)
	}

	// seventh task
	dailyProductStats, err := vazifalar.GetDailyProductStats()
    if err != nil {
        fmt.Printf("Error while getting daily product stats: %v\n", err)
    } else {
        fmt.Println("Seventh Task - Daily Product Stats:")
        for i, stat := range dailyProductStats {
            fmt.Printf("%d. Date: %s, Amount: %d\n", i+1, stat.Date, stat.Amount)
        }
    }
	// eighth task
	productTransactionSummary, err := vazifalar.GetProductTransactionSummary()
	if err != nil {
		fmt.Printf("Error while getting product transaction summary: %v\n", err)
	} else {
		fmt.Println("Eighth Task - Product Transaction Summary:")
		for productName, transactionSummary := range productTransactionSummary {
			fmt.Printf("Product: %s\n", productName)
			fmt.Printf("Kiritilgan: %d\n", transactionSummary["Kiritilgan"])
			fmt.Printf("Chiqarilgan: %d\n", transactionSummary["Chiqarilgan"])
			fmt.Println()
		}
	}

	// ninth task
	ninthTaskResult, err := vazifalar.GetBranchProductValues()
	if err != nil {
		fmt.Printf("Error while getting branch product values: %v\n", err)
	} else {
		fmt.Println("Ninth Task - Branch Product Values:")
		for i, result := range ninthTaskResult {
			fmt.Printf("%d. Branch: %s, Total Product Value: %d\n", i+1, result.BranchName, result.TotalProductValue)
		}
	}
	// tenth task
	userTransactionSummaries := vazifalar.GetUserTransactionSummaries()
	fmt.Println("tenth task - User Transaction Summaries:")
	for i, summary := range userTransactionSummaries {
		fmt.Printf("%d. UserID: %s, UserName: %s, TotalSum: %d\n", i+1, summary.UserID, summary.UserName, summary.TotalSum)

	}
	// // eleventh task
	summaries := vazifalar.GetUserDailyTransactionSummaries()
	fmt.Println("eleventh task GetUserDailyTransactionSummaries")
	for _, summary := range summaries {
		fmt.Printf("User: %s, Date: %s, Transaction Count: %d, Total Sum: %d\n",
			summary.UserID, summary.Date, summary.TransactionCount, summary.TotalSum)
	}
	// twelfth task
	userProductTransactions := vazifalar.GetUserProductTransactions()
	fmt.Println("twelfth task - User Product Transactions:")
	for _, transaction := range userProductTransactions {
		fmt.Printf("User: %s, Products In: %d, Products Out: %d\n", transaction.UserID, transaction.ProductsIn, transaction.ProductsOut)
	}

	/*===================*/
	// Category
	// h.CreateCategory("fruits")
	// h.UpdateCategory("eee3d471-78c1-4aaf-8c62-4df8acbb5c1f", "imbir")
	// h.DeleteCategory("eee3d471-78c1-4aaf-8c62-4df8acbb5c1f")
	// h.GetAllCategory(1, 10)
	// h.GetCategory("ff4cb227-92d7-4a7f-83fb-4d7e3e52d4aa")
	// Product
	// h.CreateProduct("52675e16-f992-41ae-9f8d-ec2fd8c77e44", "Sok", 10500)
	// h.UpdateProduct("e4b2c0f4-b0b5-43ec-a48b-f66513e73d66", "b8301a79-a369-40b5-b518-3d68d2b43c38", "olma", 100_000)
	// h.GetAllProduct(1, 10)
	// h.GetProduct("e4b2c0f4-b0b5-43ec-a48b-f66513e73d66")
	// h.DeleteProduct("e4b2c0f4-b0b5-43ec-a48b-f66513e73d66")
	// }
	// }
}
