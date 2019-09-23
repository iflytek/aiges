package finder

// JSONResult defined for companion
type JSONResult struct {
	Ret  int                    `json:"ret"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}
