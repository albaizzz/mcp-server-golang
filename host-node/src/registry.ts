import type { McpServerMeta } from "./types";

export const registry: McpServerMeta[] = [
  { name: "demo-products", baseUrl: process.env.MCP_SERVER1 || "http://localhost:3001" },
];