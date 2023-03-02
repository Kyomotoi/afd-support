package db

type AfdianOrders struct {
	OrderNo  string
	Time     int64
	UserID   string
	Consumed int
}

type AfdianUsers struct {
	UserID   string
	UserName string
}
