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

type ReadItem struct {
	CardUUID  string        `json:"card_uuid"`
	OrderUUID string        `json:"order_uuid"`
	OrderName string        `json:"order_name"`
	Sequence  int64         `json:"sequence"`
	ReadAt    *int64        `json:"read_at,omitempty"`
	Laser     []LaserValue  `json:"laser"`
	Magnet    []MagnetValue `json:"magnet"`
	Mifare    []MifareValue `json:"mifare"`
}

type LaserValue struct {
	Side        string `json:"side"`
	Row         int64  `json:"row"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
}

type MagnetValue struct {
	Track   int    `json:"track"`
	Content string `json:"content"`
}

type MifareValue struct {
	Block   int    `json:"block"`
	Content string `json:"content"`
}

type ReadResponse struct {
	Count int        `json:"count"`
	Items []ReadItem `json:"items"`
}

type ReadReportItem struct {
	ReadItem
	CreatedAt int64 `json:"created_at"`
	Read      bool  `json:"read"`
}

type ReadReportResponse struct {
	Count int              `json:"count"`
	Items []ReadReportItem `json:"items"`
}
