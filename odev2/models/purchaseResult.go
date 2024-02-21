package models

type PurchaseResult struct{
	AccountID int
	AccountName string
	Success bool
	ItemNameAndQuantity map[string]int
	TotalCost float64
	PurchaseText string
}