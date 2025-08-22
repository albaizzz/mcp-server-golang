package mcp

import (
	"encoding/json"
	"net/http"

	"mcp-kodepos/internal/models"
	"mcp-kodepos/internal/repo"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repo *repo.KodeposRepo
}

func NewHandler(r *repo.KodeposRepo) *Handler {
	return &Handler{Repo: r}
}

func (h *Handler) Post(c *gin.Context) {
	var req models.JSONRPCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.JSONRPCResponse{
			JSONRPC: "2.0",
			Error:   &models.JSONRPCErrorObj{Code: -32700, Message: "Parse error"},
		})
		return
	}

	// Ensure jsonrpc version
	if req.JSONRPC != "2.0" && req.JSONRPC != "" {
		c.JSON(http.StatusOK, h.err(req.ID, -32600, "Invalid Request: jsonrpc must be '2.0'", nil))
		return
	}

	switch req.Method {
	case "initialize":
		result := models.InitializeResult{
			ServerName:   "mcp-kodepos",
			Capabilities: []string{"zip.lookup", "zip.find", "zip.suggest"},
			Protocol:     "jsonrpc-2.0",
			Methods:      []string{"initialize", "zip.lookup", "zip.find", "zip.suggest"},
			Version:      "1.0.0",
		}
		c.JSON(http.StatusOK, models.JSONRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: result})

	case "zip.lookup":
		var p models.ZipLookupParams
		if err := bindParams(req.Params, &p); err != nil || p.Kodepos == "" {
			c.JSON(http.StatusOK, h.err(req.ID, -32602, "Invalid params: kodepos required", nil))
			return
		}
		rows, err := h.Repo.Lookup(c, p.Kodepos)
		if err != nil {
			c.JSON(http.StatusOK, h.err(req.ID, -32000, "DB error", err.Error()))
			return
		}
		c.JSON(http.StatusOK, models.JSONRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: rows})

	case "zip.find":
		var p models.ZipFindParams
		if err := bindParams(req.Params, &p); err != nil {
			c.JSON(http.StatusOK, h.err(req.ID, -32602, "Invalid params", nil))
			return
		}
		rows, err := h.Repo.Find(c, p)
		if err != nil {
			c.JSON(http.StatusOK, h.err(req.ID, -32000, "DB error", err.Error()))
			return
		}
		c.JSON(http.StatusOK, models.JSONRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: rows})

	case "zip.suggest":
		var p models.ZipSuggestParams
		if err := bindParams(req.Params, &p); err != nil || p.Prefix == "" {
			c.JSON(http.StatusOK, h.err(req.ID, -32602, "Invalid params: prefix required", nil))
			return
		}
		rows, err := h.Repo.Suggest(c, p.Prefix, p.Limit)
		if err != nil {
			c.JSON(http.StatusOK, h.err(req.ID, -32000, "DB error", err.Error()))
			return
		}
		c.JSON(http.StatusOK, models.JSONRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: rows})

	default:
		c.JSON(http.StatusOK, h.err(req.ID, -32601, "Method not found", req.Method))
	}
}

func (h *Handler) err(id any, code int, msg string, data any) models.JSONRPCResponse {
	return models.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &models.JSONRPCErrorObj{Code: code, Message: msg, Data: data},
	}
}

// --- helpers
func bindParams(src any, dst any) error {
	// gin already unmarshals into interface{} with map[string]any; re-marshal & unmarshal is simplest
	b, err := jsonMarshal(src)
	if err != nil {
		return err
	}
	return jsonUnmarshal(b, dst)
}

// tiny wrappers to avoid importing in other files
func jsonMarshal(v any) ([]byte, error)   { return json.Marshal(v) }
func jsonUnmarshal(b []byte, v any) error { return json.Unmarshal(b, v) }
