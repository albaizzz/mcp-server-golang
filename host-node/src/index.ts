import express from "express";
import fetch from "node-fetch";
import dotenv from "dotenv";
import cors from "cors";
import { registry } from "./registry";
import type { McpNormalizedResource, McpResourceRef } from "./types";

dotenv.config();

const app = express();
app.use(express.json());
app.use(cors());

app.get("/mcp/servers", (_req, res) => {
  res.json({ servers: registry.map(s => ({ name: s.name, baseUrl: s.baseUrl })) });
});

app.get("/mcp/resources", async (_req, res) => {
  const results: { server: string; resources: McpResourceRef[] }[] = [];
  for (const s of registry) {
    try {
      const r = await fetch(`${s.baseUrl}/mcp/resources`);
      const j = (await r.json()) as { resources: McpResourceRef[] };
      results.push({ server: s.name, resources: j.resources || [] });
    } catch (e) {
      results.push({ server: s.name, resources: [] });
    }
  }
  res.json({ aggregated: results });
});

app.get("/mcp/:server/resources/:resource", async (req, res) => {
  const { server, resource } = req.params;
  const s = registry.find(x => x.name === server);
  if (!s) return res.status(404).json({ error: "server not found" });
  try {
    const r = await fetch(`${s.baseUrl}/mcp/resources/${resource}`);
    const j = await r.json();
    const normalized: McpNormalizedResource = {
      server: s.name,
      name: resource,
      items: (j as any).data ?? j,
    };
    res.json(normalized);
  } catch (e) {
    res.status(502).json({ error: "upstream error" });
  }
});

app.get("/mcp/:server/resources/:resource/:id", async (req, res) => {
  const { server, resource, id } = req.params;
  const s = registry.find(x => x.name === server);
  if (!s) return res.status(404).json({ error: "server not found" });
  try {
    const r = await fetch(`${s.baseUrl}/mcp/resources/${resource}/${id}`);
    const j = await r.json();
    res.json({ server: s.name, name: resource, item: (j as any).data ?? j });
  } catch (e) {
    res.status(502).json({ error: "upstream error" });
  }
});

const port = Number(process.env.PORT || 4000);
app.listen(port, () => {
  console.log(`MCP Host running on http://localhost:${port}`);
});