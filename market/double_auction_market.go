package market

import (
	"container/heap"

	"github.com/robbrit/econerra/goods"
)

type doubleAuction struct {
	bids       orderMaxHeap
	offers     orderMinHeap
	lastHigh   Price
	lastLow    Price
	high       Price
	low        Price
	lastVolume Size
	volume     Size
	bid        Price
	ask        Price
	good       goods.Good
}

// NewDoubleAuction constructs a new market for a given good.
func NewDoubleAuction(good goods.Good) Market {
	m := &doubleAuction{
		good: good,
	}
	m.Reset()
	return m
}

func (m *doubleAuction) Bid() Price       { return m.bid }
func (m *doubleAuction) Ask() Price       { return m.ask }
func (m *doubleAuction) High() Price      { return m.lastHigh }
func (m *doubleAuction) Low() Price       { return m.lastLow }
func (m *doubleAuction) Volume() Size     { return m.lastVolume }
func (m *doubleAuction) Good() goods.Good { return m.good }

// Post sends an order to the market. If this order results in a fill,
// the owner(s) will be notified. If not, the order will remain open in
// the market.
func (m *doubleAuction) Post(o *Order) {
	if o.Size == 0 {
		return
	}
	if o.Price <= 0 {
		return
	}

	switch o.Side {
	case Buy:
		if len(m.offers) == 0 || o.Price < m.offers[0].Price {
			heap.Push(&m.bids, o)
			return
		}

		// Pop sell orders off the heap until we have filled the entire amount.
		size := o.Size
		for len(m.offers) > 0 && o.Price >= m.offers[0].Price && size > 0 {
			if m.offers[0].Size <= size {
				sell := heap.Pop(&m.offers).(*Order)
				m.handleFill(o, sell, sell.Price, sell.Size)
				size -= sell.Size
			} else {
				sell := m.offers[0]
				m.handleFill(o, sell, sell.Price, size)
				m.offers[0].Size -= size
				size = 0
			}
		}

		if size > 0 {
			o.Size = size
			heap.Push(&m.bids, o)
		}
	case Sell:
		if len(m.bids) == 0 || o.Price > m.bids[0].Price {
			heap.Push(&m.offers, o)
			return
		}

		// Pop buy orders off the heap until we have filled the entire amount.
		size := o.Size
		for len(m.bids) > 0 && o.Price <= m.bids[0].Price && size > 0 {
			if m.bids[0].Size <= size {
				buy := heap.Pop(&m.bids).(*Order)
				m.handleFill(buy, o, buy.Price, buy.Size)
				size -= buy.Size
			} else {
				buy := m.bids[0]
				m.handleFill(buy, o, buy.Price, size)
				m.bids[0].Size -= size
				size = 0
			}
		}

		if size > 0 {
			o.Size = size
			heap.Push(&m.offers, o)
		}
	}
}

func (m *doubleAuction) handleFill(buy, sell *Order, price Price, size Size) {
	buy.Owner.OnFill(m.good, Buy, price, size)
	sell.Owner.OnFill(m.good, Sell, price, size)

	if price > m.high {
		m.high = price
	}
	if m.low == 0 || price < m.low {
		m.low = price
	}
	m.volume += size
}

func (m *doubleAuction) Reset() {
	m.lastLow = m.low
	m.lastHigh = m.high
	m.lastVolume = m.volume
	m.high = 0
	m.low = 0
	m.volume = 0

	// Clear out all the orders, sending unfilled notifications as needed.
	for _, order := range m.bids {
		order.Owner.OnUnfilled(m.good, Buy, order.Size)
	}
	for _, order := range m.offers {
		order.Owner.OnUnfilled(m.good, Sell, order.Size)
	}

	if len(m.bids) > 0 {
		m.bid = m.bids[0].Price
	} else {
		m.bid = 0
	}
	if len(m.offers) > 0 {
		m.ask = m.offers[0].Price
	} else {
		m.ask = 0
	}

	m.bids = orderMaxHeap{}
	m.offers = orderMinHeap{}
	heap.Init(&m.bids)
	heap.Init(&m.offers)
}
