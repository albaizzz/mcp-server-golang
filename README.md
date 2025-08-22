
# MCP Suite (Server-Go, Host-Node, Client-Node + LLM Integration)

Semua dalam satu proyek:
- **server-go**: MCP Server (Go + Gin) expose `products`.
- **host-node**: MCP Host (Node+TS) untuk discovery/proxy.
- **client-node**: MCP Client/Agent (Node+TS) + integrasi LLM (mock, llama3 generic, openai-compat).

## Prasyarat
- Go 1.20+
- Node 18+
- (Opsional) Llama3/TGI server kalau mau mode `llama3`

## Cara Jalan (lokal)
### 1) MCP Server (Go)
```bash
cd server-go
go mod tidy
go run main.go
# http://localhost:3001
```

### 2) MCP Host (Node)
```bash
cd host-node
npm install
npm run dev
# http://localhost:4000
```

### 3) Client/Agent (Node)
Mock (tanpa LLM):
```bash
cd client-node
npm install
export LLM_MODE=mock
npm run dev
```

Llama3 (generic):
```bash
export LLM_MODE=llama3
export LLM_URL=http://localhost:8080
export LLM_MODEL=llama-3-8b
npm run dev
```

OpenAI-compat:
```bash
export LLM_MODE=openai
export OPENAI_MODEL=gpt-4o-mini
export LLM_API_KEY=sk-... # API key
npm run dev
```

## Endpoint
- Server: `/mcp/resources`, `/mcp/resources/products`, `/mcp/resources/products/:id`
- Host: `/mcp/servers`, `/mcp/resources`, `/mcp/demo-products/resources/products`

## Custom ke Ride-Hailing
Ganti resource `products` → `drivers` / `rides` dan sesuaikan prompt di `client-node/src/agent.ts`.
