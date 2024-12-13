package types

type OnChainRequest struct {
	Business  []string `json:"business,optional"`
	Source    []string `json:"source,optional"`
	Data      string   `json:"data"`
	Extension []string `json:"extension,optional"`
	RequestId string   `json:"requestId,optional"`
}
