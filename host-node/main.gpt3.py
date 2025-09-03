import requests
import json

OPENAI_API_KEY = "OPENAI_API_KEY"
OPENAI_URL = "https://api.openai.com/v1/chat/completions"
MCP_SERVER_URL = "http://localhost:8111/mcp"

headers = {
    "Authorization": f"Bearer {OPENAI_API_KEY}",
    "Content-Type": "application/json"
}

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

def call_gpt35(user_prompt):
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

    # Step 1: kirim user prompt
    payload = {
        "model": "gpt-3.5-turbo-1106",
        "messages": [{"role": "user", "content": user_prompt}],
        "tools": tools,
        "tool_choice": "auto"
    }
    response = requests.post(OPENAI_URL, headers=headers, json=payload)
    resp = response.json()

    msg = resp["choices"][0]["message"]

    # Step 2: cek apakah model memanggil tool
    if "tool_calls" in msg:
        for tool_call in msg["tool_calls"]:
            name = tool_call["function"]["name"]
            args = json.loads(tool_call["function"]["arguments"])
            print(f"🔧 GPT memilih tool: {name} dengan args {args}")

            if name == "zip_lookup":
                result = mcp_client.call("zip.lookup", args)
            elif name == "zip_find":
                result = mcp_client.call("zip.find", args)
            else:
                result = {"error": "Unknown tool"}

            # Step 3: kirim balik hasil tool ke GPT
            followup_payload = {
                "model": "gpt-3.5-turbo-1106",
                "messages": [
                    {"role": "user", "content": user_prompt},
                    msg,
                    {
                        "role": "tool",
                        "tool_call_id": tool_call["id"],
                        "content": json.dumps(result)
                    }
                ]
            }
            followup_resp = requests.post(
                OPENAI_URL, headers=headers, json=followup_payload
            ).json()
            return followup_resp["choices"][0]["message"]["content"]

    return msg["content"]


# ==== DEMO ====
print(call_gpt35("Kodepos untuk Gambir berapa?"))
