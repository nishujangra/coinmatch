package engine

import (
	"container/heap"

	"github.com/nishujangra/coinmatch/lib/models"
)

// BuyOrderQueue implements a max-heap for buy orders (highest price first)
type BuyOrderPQ []*models.Order

func (pq BuyOrderPQ) Len() int {
	return len(pq)
}

func (pq BuyOrderPQ) Less(i, j int) bool {
	// Max heap by price, then FIFO by time for same price
	if pq[i].Price == pq[j].Price {
		return pq[i].CreatedAt.Before(pq[j].CreatedAt)
	}
	return pq[i].Price > pq[j].Price
}

func (pq BuyOrderPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *BuyOrderPQ) Push(x interface{}) {
	order := x.(*models.Order)
	*pq = append(*pq, order)
}

func (pq *BuyOrderPQ) Pop() interface{} {
	old := *pq
	n := len(old)

	// last order at index n-1
	order := old[n-1]
	old[n-1] = nil

	*pq = old[0 : n-1]
	return order
}

type SellOrderPQ []*models.Order

func (pq SellOrderPQ) Len() int {
	return len(pq)
}

func (pq SellOrderPQ) Less(i, j int) bool {
	// Max heap by price, then FIFO by time for same price
	if pq[i].Price == pq[j].Price {
		return pq[i].CreatedAt.Before(pq[j].CreatedAt)
	}
	return pq[i].Price < pq[j].Price
}

func (pq SellOrderPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *SellOrderPQ) Push(x interface{}) {
	order := x.(*models.Order)
	*pq = append(*pq, order)
}

func (pq *SellOrderPQ) Pop() interface{} {
	old := *pq
	n := len(old)

	// last order at index n-1
	order := old[n-1]
	old[n-1] = nil

	*pq = old[0 : n-1]
	return order
}

type OrderBook struct {
	BuyPQ  BuyOrderPQ
	SellPQ SellOrderPQ
}

var Books = map[string]*OrderBook{} // pair is the key

func MatchOrder(order *models.Order, book *OrderBook) {
	if order.Side == "buy" {
		for order.Quantity > 0 && book.SellPQ.Len() > 0 {
			bestSell := book.SellPQ[0] // heap top
			if order.Price < bestSell.Price {
				break // no match
			}
			matchQty := min(order.Quantity, bestSell.Quantity)
			order.Quantity -= matchQty
			bestSell.Quantity -= matchQty

			if bestSell.Quantity == 0 {
				heap.Pop(&book.SellPQ)
			}
		}
		if order.Quantity > 0 {
			heap.Push(&book.BuyPQ, order)
		}
	} else if order.Side == "sell" {
		for order.Quantity > 0 && book.BuyPQ.Len() > 0 {
			bestBuy := book.BuyPQ[0]
			if order.Price > bestBuy.Price {
				break
			}
			matchQty := min(order.Quantity, bestBuy.Quantity)
			order.Quantity -= matchQty
			bestBuy.Quantity -= matchQty

			if bestBuy.Quantity == 0 {
				heap.Pop(&book.BuyPQ)
			}
		}
		if order.Quantity > 0 {
			heap.Push(&book.SellPQ, order)
		}
	}
}
