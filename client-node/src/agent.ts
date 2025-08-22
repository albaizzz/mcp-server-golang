import dotenv from "dotenv";
import { listProducts, getProductById } from "./mcpClient";
import { llmGenerate } from "./llm";

dotenv.config();

async function main() {
  const resource = await listProducts();
  const items = resource.items;
  const detail = await getProductById(2);

  const prompt = `Kamu adalah analis. Buat ringkasan inventori berikut (harga & stok), dan rekomendasikan 1 produk unggulan.
Daftar Produk:
${items.map((p:any)=>`- ${p.name} (ID:${p.id}) Rp${p.price}, stok ${p.stock}`).join("\n")}

Sertakan ringkas detail item ID=2: ${JSON.stringify(detail.item)}.
`;

  const res = await llmGenerate({
    system: "Jawab singkat, jelas, bullet points bila perlu.",
    prompt,
    max_tokens: 400,
  });

  console.log("=== HASIL AGENT ===\n");
  console.log(res.text);
}

main().catch(err => {
  console.error(err);
  process.exit(1);
});