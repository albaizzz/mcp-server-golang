import fetch from "node-fetch";
import type { McpResource } from "./types";

const HOST = process.env.MCP_HOST || "http://localhost:4000";

export async function listProducts(serverName = "demo-products"): Promise<McpResource> {
  const r = await fetch(`${HOST}/mcp/${serverName}/resources/products`);
  if (!r.ok) throw new Error(`Host error ${r.status}`);
  return (await r.json()) as McpResource;
}

export async function getProductById(id: number, serverName = "demo-products") {
  const r = await fetch(`${HOST}/mcp/${serverName}/resources/products/${id}`);
  if (!r.ok) throw new Error(`Host error ${r.status}`);
  return await r.json();
}