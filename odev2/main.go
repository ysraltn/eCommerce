package main

import (
	"fmt"
	"odev2/models"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	pChannel := make(chan models.PurchaseResult, 10)
	iSChannel := make(chan *models.Item, 5)
	market := models.MarketInit()

	account1 := models.AccountInit(1, "ahmet", 10000, market)
	fmt.Println(account1)
	account2 := models.AccountInit(2, "mahmut", 1000, market)
	fmt.Println(account2)
	account3 := models.AccountInit(3, "veli", 1000, market)
	fmt.Println(account3)

	item1 := models.ItemInit(1, "masa", 10, 220, market)
	fmt.Println(item1)

	item2 := models.ItemInit(2, "sandalye", 12, 95, market)
	fmt.Println(item2)

	account1.AddToCart(item1, 3)
	account1.AddToCart(item2, 5)
	fmt.Println(account1)
	account1.ListTheCart()

	account2.AddToCart(item2, 6)
	account2.ListTheCart()

	account3.AddToCart(item2, 2)
	account3.ListTheCart()

	market.ListAccounts()
	market.ListItems()

	//mutex ile itemin korunduğunu ve stogun hiç eksiye
	//dusmedigini daha net görmek icin i<0 yapabilirsiniz.
	//account.completePurchase fonksiyonu içindeki mutex lock/unlocklari kaldirirsaniz
	//bazen race condition olustugunu ve sandalye stogunun -1'e düstügünü görebilirsiniz.
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(item *models.Item) {
			defer wg.Done()
			item2.IncrementStock(2, iSChannel)
		}(item2)
	}

	for _, account := range market.Accounts {
		wg.Add(1)
		//var input2 string
		//fmt.Scanln(&input2)
		go func(account *models.Account) {
			defer wg.Done()
			account.CompletePurchase(pChannel)
		}(account)
	}

	/* 	//wg.Add(1)
	   	go func() {
	   		//defer wg.Done()
	   		for purschaseResult := range pChannel {
	   			fmt.Println("purchase details:")
	   			fmt.Printf("	success : %t\n", purschaseResult.Success)
	   			fmt.Printf("	account id: %d\n", purschaseResult.AccountID)
	   			fmt.Printf("	account name: %s\n", purschaseResult.AccountName)
	   			fmt.Printf("	items:\n")
	   			for item, quantity := range purschaseResult.ItemNameAndQuantity {
	   				fmt.Printf("		%s x %d\n", item, quantity)
	   			}
	   			fmt.Printf("	total cost: %.2f\n", purschaseResult.TotalCost)
	   			fmt.Printf("	message: %s\n", purschaseResult.PurchaseText)
	   		}
	   	}() */

	for i := 0; i < 10; i++ {
		select {
		case purschaseResult := <-pChannel:
			fmt.Println("purchase details:")
			fmt.Printf("	success : %t\n", purschaseResult.Success)
			fmt.Printf("	account id: %d\n", purschaseResult.AccountID)
			fmt.Printf("	account name: %s\n", purschaseResult.AccountName)
			fmt.Printf("	items:\n")
			for item, quantity := range purschaseResult.ItemNameAndQuantity {
				fmt.Printf("		%s x %d\n", item, quantity)
			}
			fmt.Printf("	total cost: %.2f\n", purschaseResult.TotalCost)
			fmt.Printf("	message: %s\n", purschaseResult.PurchaseText)
		case updatedItem := <-iSChannel:
			fmt.Println("updated item details:")
			fmt.Printf("	ID: %d\n", updatedItem.ID)
			fmt.Printf("	name: %s\n", updatedItem.Name)
			fmt.Printf("	stock: %d\n", updatedItem.Stock)
			fmt.Printf("	price: %.2f\n", updatedItem.Price)
		default:
			fmt.Println("default")
		}
	}

	wg.Wait()
	close(pChannel)
	close(iSChannel)
	//var input string
	//fmt.Scanln(&input)
	fmt.Println(account2)
	market.ListItems()
}
