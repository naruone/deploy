package request

type BasePageInfo struct {
    CurrentPage int    `json:"currentPage" form:"currentPage"`
    PageSize    int    `json:"pageSize" form:"pageSize"`
    Condition   string `json:"sCondition" form:"sCondition"`
    SearchValue string `json:"sValue" form:"sValue"`
}

type ComPageInfo struct {
    BasePageInfo
    Status int `json:"status" form:"status"`
    Type   int `json:"type" form:"type"`
}
