export type McpServerMeta = {
  name: string;
  baseUrl: string; // e.g. http://localhost:3001
};

export type McpResourceRef = { name: string; path: string };

export type McpNormalizedResource = {
  server: string;
  name: string;
  items: unknown;
};