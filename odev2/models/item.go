package models

import "sync"

type Item struct {
	ID    int
	Name  string
	Stock int
	Price float64
	mu    sync.Mutex
}

func ItemInit(id int, name string, stock int, price float64, market *Market) *Item {
	item := Item{ID: id, Name: name, Stock: stock, Price: price}
	market.addItem(&item)
	return &item
}

func (item *Item) IncrementStock(amount int, incStockChan chan *Item) {
	item.mu.Lock()
	item.Stock += amount
	item.mu.Unlock()
	incStockChan <- item
}
