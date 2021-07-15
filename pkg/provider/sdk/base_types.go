package sdk

type Pagination struct {
	TotalResults int `json:"totalResults"`
	TotalPages   int `json:"totalPages"`
	PerPage      int `json:"perPage"`
	Page         int `json:"page"`
}
