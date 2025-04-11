import type { SettlersIncomingMessage } from "./messages";

export function safeParse(text: string):
  | {
      [K in keyof SettlersIncomingMessage]: SettlersIncomingMessage;
    }[keyof SettlersIncomingMessage]
  | null {
  try {
    const parsed = JSON.parse(text);
    if (parsed && typeof parsed.type === "string") {
      return parsed;
    }
  } catch (err) {
    console.error("Error while safe parsing the following text blob:", text);
    return null;
  }
  return null;
}
