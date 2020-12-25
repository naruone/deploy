package request

type BasePageInfo struct {
    CurrentPage int    `json:"currentPage"`
    PageSize    int    `json:"pageSize"`
    Condition   string `json:"sCondition"`
    SearchValue string `json:"sValue"`
}

type ComPageInfo struct {
    BasePageInfo
    Status int `json:"status"`
    Type   int `json:"type"`
}
