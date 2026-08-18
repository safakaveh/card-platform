package uploadcsv

type ImportResult struct {
	UUID         string   `json:"uuid"`
	OrderName    string   `json:"order_name"`
	FileName     string   `json:"file_name"`
	RowsImported int64    `json:"rows_imported"`
	FrontColumns []string `json:"front_columns"`
	BackColumns  []string `json:"back_columns"`
	HasUID       bool     `json:"has_uid"`
}

type ImportSummary struct {
	UUID      string `json:"uuid"`
	OrderName string `json:"order_name"`
	Status    string `json:"status"`
	FileName  string `json:"file_name"`
	CardCount int64  `json:"card_count"`
	CreatedAt int64  `json:"created_at"`
}

type ImportDetails struct {
	ImportSummary
	FrontColumns int64 `json:"front_columns"`
	BackColumns  int64 `json:"back_columns"`
	UIDFields    int64 `json:"uid_fields"`
}

type columnMapping struct {
	index       int
	side        string
	rowNumber   int64
	contentType string
	header      string
	isUID       bool
	isImage     bool
	trackNo     int
}
