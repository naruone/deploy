package response

type PageResult struct {
    List        interface{} `json:"list"`
    Total       int         `json:"total"`
    CurrentPage int         `json:"currentPage"`
    PageSize    int         `json:"pageSize"`
}
