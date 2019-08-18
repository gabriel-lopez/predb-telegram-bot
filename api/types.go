package api

type apiResponse struct {
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Data    apiRowData `json:"data"`
}

type apiRowData struct {
	RowCount int         `json:"rowCount"`
	Rows     []sphinxRow `json:"rows"`
	Offset   int         `json:"offset"`
	ReqCount int         `json:"reqCount"`
	Total    int         `json:"total"`
	Time     float64     `json:"time"`
}

type sphinxRow struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Team  string  `json:"team"`
	Cat   string  `json:"cat"`
	Genre string  `json:"genre"`
	URL   string  `json:"url"`
	Size  float64 `json:"size"`
	Files int     `json:"files"`
	PreAt int64   `json:"preAt"`
}