import fetch from "node-fetch";

export type LlmRequest = { system?: string; prompt: string; max_tokens?: number };
export type LlmResponse = { text: string };

const MODE = process.env.LLM_MODE || "mock"; // mock | openai | llama3
const LLM_URL = process.env.LLM_URL || "http://localhost:8080";
const LLM_MODEL = process.env.LLM_MODEL || "llama-3-8b";
const API_KEY = process.env.LLM_API_KEY || "";
const OPENAI_MODEL = process.env.OPENAI_MODEL || "gpt-4o-mini";

async function callOpenAI(req: LlmRequest): Promise<LlmResponse> {
  const apiKey = API_KEY;
  if (!apiKey) throw new Error("OPENAI_API_KEY missing");
  const res = await fetch("https://api.openai.com/v1/chat/completions", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${apiKey}`,
    },
    body: JSON.stringify({
      model: OPENAI_MODEL,
      messages: [
        ...(req.system ? [{ role: "system", content: req.system }] : []),
        { role: "user", content: req.prompt },
      ],
      max_tokens: req.max_tokens ?? 512,
      temperature: 0.2,
    }),
  });
  if (!res.ok) throw new Error(`OpenAI error ${res.status}`);
  const j: any = await res.json();
  const text = j.choices?.[0]?.message?.content ?? "";
  return { text };
}

async function callLlama3Generic(req: LlmRequest): Promise<LlmResponse> {
  const headers: any = { "Content-Type": "application/json" };
  if (API_KEY) headers["Authorization"] = `Bearer ${API_KEY}`;

  // Try OpenAI-compatible first
  const openaiCompat = await fetch(`${LLM_URL}/v1/completions`, {
    method: "POST",
    headers,
    body: JSON.stringify({
      model: LLM_MODEL,
      input: req.prompt,
      max_tokens: req.max_tokens ?? 512,
      temperature: 0.2,
    }),
  }).catch(() => null);

  if (openaiCompat && openaiCompat.ok) {
    const j: any = await openaiCompat.json();
    if (j.choices?.[0]?.text) return { text: j.choices[0].text };
    if (j.choices?.[0]?.message?.content) return { text: j.choices[0].message.content };
  }

  // Try TGI style
  const tgi = await fetch(`${LLM_URL}/v1/models/${LLM_MODEL}/generate`, {
    method: "POST",
    headers,
    body: JSON.stringify({
      inputs: req.prompt,
      parameters: { max_new_tokens: req.max_tokens ?? 512, temperature: 0.2 }
    }),
  }).catch(() => null);

  if (tgi && tgi.ok) {
    const j: any = await tgi.json();
    if (j.generated_text) return { text: j.generated_text };
    if (j.output && Array.isArray(j.output)) {
      const blocks = j.output[0]?.content;
      if (Array.isArray(blocks)) return { text: blocks.map((b: any) => b.text ?? b).join("") };
      if (typeof blocks === "string") return { text: blocks };
    }
  }

  // Fallback raw
  const raw = await fetch(LLM_URL, {
    method: "POST",
    headers,
    body: JSON.stringify({ model: LLM_MODEL, prompt: req.prompt, max_new_tokens: req.max_tokens ?? 512 }),
  }).catch(() => null);

  if (raw && raw.ok) {
    const j: any = await raw.json();
    if (j.generated_text) return { text: j.generated_text };
    if (j.data?.[0]?.generated_text) return { text: j.data[0].generated_text };
    return { text: JSON.stringify(j).slice(0, 1200) };
  }

  throw new Error("Llama3 adapter: unable to get a response (check LLM_URL/model).");
}

export async function llmGenerate(req: LlmRequest): Promise<LlmResponse> {
  if (MODE === "mock") {
    return { text: `MOCK: ${req.prompt.substring(0, 200)}...` };
  }
  if (MODE === "openai") {
    return await callOpenAI(req);
  }
  if (MODE === "llama3") {
    return await callLlama3Generic(req);
  }
  throw new Error(`Unknown LLM_MODE=${MODE}`);
}