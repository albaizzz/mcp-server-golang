import axios from "axios";
import dotenv from "dotenv";

dotenv.config();

const HOST_URL = process.env.HOST_URL || "http://localhost:3000";
const LLM_MODE = process.env.LLM_MODE || "mock";

// Fungsi panggil Host
async function fetchDrivers() {
  try {
    const res = await axios.get(`${HOST_URL}/drivers`);
    return res.data;
  } catch (err) {
    console.error("Error fetch drivers:", err.message);
    return [];
  }
}

// Mock integrasi LLM
async function askLLM(prompt) {
  if (LLM_MODE === "mock") {
    return `Mock LLM response untuk prompt: "${prompt}"`;
  }

  // misal kalau pakai llama3 (via OpenAI API compat)
  try {
    const res = await axios.post(
      process.env.LLM_ENDPOINT || "http://localhost:8000/v1/chat/completions",
      {
        model: "llama3",
        messages: [{ role: "user", content: prompt }]
      },
      {
        headers: {
          Authorization: `Bearer ${process.env.LLM_API_KEY || "mock-key"}`
        }
      }
    );
    return res.data.choices[0].message.content;
  } catch (err) {
    return `LLM error: ${err.message}`;
  }
}

async function main() {
  const drivers = await fetchDrivers();
  console.log("Drivers:", drivers);

  const llmResp = await askLLM("Siapa driver terdekat?");
  console.log("LLM Response:", llmResp);
}

main();
