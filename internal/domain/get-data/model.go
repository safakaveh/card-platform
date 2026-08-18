package getdata

type PendingCard struct {
	CardUUID  string `json:"card_uuid"`
	OrderUUID string `json:"order_uuid"`
	OrderName string `json:"order_name"`
	BlockNo   int    `json:"block_no"`
	CreatedAt int64  `json:"created_at"`
}

type PendingResponse struct {
	Count int           `json:"count"`
	Items []PendingCard `json:"items"`
}
