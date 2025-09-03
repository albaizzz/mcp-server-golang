import requests
import json

OLLAMA_URL = "http://localhost:11434/api/chat"
MCP_SERVER_URL = "http://localhost:8080/mcp"

class MCPClient:
    def __init__(self, url=MCP_SERVER_URL):
        self.url = url
        self.req_id = 0

    def call(self, method, params=None):
        self.req_id += 1
        payload = {
            "jsonrpc": "2.0",
            "id": self.req_id,
            "method": method,
            "params": params or {}
        }
        resp = requests.post(self.url, json=payload)
        return resp.json()

mcp_client = MCPClient()

def call_llama3(user_prompt):
    # Define tools
    tools = [
        {
            "type": "function",
            "function": {
                "name": "zip_lookup",
                "description": "Cari detail lokasi berdasarkan kodepos",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "kodepos": {"type": "string"}
                    },
                    "required": ["kodepos"]
                }
            }
        },
        {
            "type": "function",
            "function": {
                "name": "zip_find",
                "description": "Cari kodepos berdasarkan kecamatan/kota/provinsi",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "kecamatan": {"type": "string"},
                        "kota": {"type": "string"},
                        "provinsi": {"type": "string"}
                    }
                }
            }
        }
    ]

    payload = {
        "model": "llama3:8b-instruct-q4_K_M",
        "messages": [{"role": "user", "content": user_prompt}],
        "tools": tools,
        "stream": False
    }

    resp = requests.post(OLLAMA_URL, json=payload).json()
    msg = resp["message"]

    # Kalau Llama3 minta pakai tool
    if "tool_calls" in msg:
        for tool_call in msg["tool_calls"]:
            name = tool_call["function"]["name"]
            args = json.loads(tool_call["function"]["arguments"])
            print(f"🔧 Llama3 memilih tool: {name} dengan args {args}")

            if name == "zip_lookup":
                result = mcp_client.call("zip.lookup", args)
            elif name == "zip_find":
                result = mcp_client.call("zip.find", args)
            else:
                result = {"error": "Unknown tool"}

            # balikin hasil ke Llama3
            followup_payload = {
                "model": "llama3",
                "messages": [
                    {"role": "user", "content": user_prompt},
                    msg,
                    {"role": "tool", "content": json.dumps(result), "name": name}
                ],
                "stream": False
            }
            followup_resp = requests.post(OLLAMA_URL, json=followup_payload).json()
            return followup_resp["message"]["content"]

    return msg["content"]

# ==== DEMO ====
print(call_llama3("Kodepos untuk Gambir berapa?"))
