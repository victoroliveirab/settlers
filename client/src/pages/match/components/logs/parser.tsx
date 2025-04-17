import { Die } from "@/components/custom/die";
import { GameCard } from "@/components/custom/game-card";
import { renderToStaticMarkup } from "react-dom/server";

// NOTE: Complexity needed? Probably not. Fun to build? Absolutely!

type ParsedChunk =
  | { type: "text"; content: string }
  | { type: "res"; q: number; v: string }
  | { type: "dev"; v: string }
  | { type: "dice"; v: number };

function tokenizer(input: string): ParsedChunk[] {
  const chunks: ParsedChunk[] = [];
  let i = 0;

  while (i < input.length) {
    if (input[i] === "[") {
      // Start of a tag
      const tagStart = i;
      const tagEnd = input.indexOf("]", tagStart);
      if (tagEnd === -1) throw new Error("Unclosed tag at position " + tagStart);

      const tagContent = input.slice(tagStart + 1, tagEnd).trim();
      const parts = tagContent.split(/\s+/); // we can still split by whitespace

      const tagName = parts[0];
      const params: Record<string, string> = {};

      for (let j = 1; j < parts.length; j++) {
        const part = parts[j];
        const eqIdx = part.indexOf("=");
        if (eqIdx === -1) throw new Error(`Invalid parameter format in tag: "${part}"`);

        const key = part.slice(0, eqIdx);
        const value = part.slice(eqIdx + 1);
        if (!key || !value) throw new Error(`Missing key or value in parameter: "${part}"`);
        params[key] = value;
      }

      // Parse tag types
      if (tagName === "res") {
        const q = parseInt(params["q"], 10);
        const v = params["v"];
        if (isNaN(q) || !v) throw new Error(`Invalid [res] tag at position ${tagStart}`);
        chunks.push({ type: "res", q, v });
      } else if (tagName === "dev") {
        const v = params["v"];
        if (!v) throw new Error(`Missing 'v' in [dev] tag at position ${tagStart}`);
        chunks.push({ type: "dev", v });
      } else if (tagName === "dice") {
        const v = parseInt(params["v"], 10);
        if (isNaN(v)) throw new Error(`Invalid 'v' in [dice] tag at position ${tagStart}`);
        chunks.push({ type: "dice", v });
      } else {
        throw new Error(`Unknown tag: "${tagName}" at position ${tagStart}`);
      }

      i = tagEnd + 1; // move past the closing ]
    } else {
      // Regular text
      const textStart = i;
      while (i < input.length && input[i] !== "[") {
        i++;
      }
      chunks.push({ type: "text", content: input.slice(textStart, i) });
    }
  }

  return chunks;
}

function transformToken(element: HTMLElement, token: ParsedChunk) {
  switch (token.type) {
    case "text": {
      const text = token.content.trim();
      if (!text) return;
      element.innerHTML += text;
      break;
    }
    case "res": {
      const resourceCard = renderToStaticMarkup(
        <GameCard size="xs" value={token.v as SettlersCore.Resource} />,
      );
      for (let i = 0; i < token.q; ++i) {
        element.innerHTML += resourceCard;
      }
      break;
    }
    case "dice": {
      const die = renderToStaticMarkup(<Die className="rounded-xs h-6" value={token.v} />);
      const div = document.createElement("div");
      div.innerHTML = die;
      element.appendChild(div);
      break;
    }
  }
}

export function parseLog(root: HTMLElement, content: string): boolean {
  let tokens: ParsedChunk[];
  try {
    tokens = tokenizer(content);
  } catch (err) {
    console.error(`Error tokenizing input "${content}": ${err}`);
    return false;
  }
  for (const token of tokens) {
    transformToken(root, token);
  }
  return true;
}
