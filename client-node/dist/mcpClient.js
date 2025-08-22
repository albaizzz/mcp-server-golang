import fetch from "node-fetch";
const HOST = process.env.MCP_HOST || "http://localhost:4000";
export async function listProducts(serverName = "demo-products") {
    const r = await fetch(`${HOST}/mcp/${serverName}/resources/products`);
    if (!r.ok)
        throw new Error(`Host error ${r.status}`);
    return (await r.json());
}
export async function getProductById(id, serverName = "demo-products") {
    const r = await fetch(`${HOST}/mcp/${serverName}/resources/products/${id}`);
    if (!r.ok)
        throw new Error(`Host error ${r.status}`);
    return await r.json();
}
