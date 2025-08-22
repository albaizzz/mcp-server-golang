package models

// JSON-RPC 2.0 types
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Method  string      `json:"method"`
	Params  any         `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      interface{}      `json:"id,omitempty"`
	Result  any              `json:"result,omitempty"`
	Error   *JSONRPCErrorObj `json:"error,omitempty"`
}

type JSONRPCErrorObj struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Domain
type ZipRecord struct {
	Kodepos   string   `json:"kodepos"`
	Kelurahan string   `json:"kelurahan"`
	Kecamatan string   `json:"kecamatan"`
	KotaKab   string   `json:"kota_kab"`
	Provinsi  string   `json:"provinsi"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

// Method params
type InitializeParams struct {
	Client string `json:"client,omitempty"`
}

type InitializeResult struct {
	ServerName   string   `json:"serverName"`
	Capabilities []string `json:"capabilities"`
	Protocol     string   `json:"protocol"`
	Methods      []string `json:"methods"`
	Version      string   `json:"version"`
}

type ZipLookupParams struct {
	Kodepos string `json:"kodepos"`
}

type ZipFindParams struct {
	Query     string `json:"query,omitempty"` // bebas: kelurahan/kecamatan/kota/provinsi/kodepos
	Provinsi  string `json:"provinsi,omitempty"`
	KotaKab   string `json:"kota_kab,omitempty"`
	Kecamatan string `json:"kecamatan,omitempty"`
	Kelurahan string `json:"kelurahan,omitempty"`
	Limit     int    `json:"limit,omitempty"` // default 20
	Offset    int    `json:"offset,omitempty"`
}

type ZipSuggestParams struct {
	Prefix string `json:"prefix"` // prefix nama wilayah atau kodepos
	Limit  int    `json:"limit,omitempty"`
}
