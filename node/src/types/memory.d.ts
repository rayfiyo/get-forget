// src/types/memory.d.ts
export interface Memory {
  content: string;
  timestamp: number;
  initialImportance: number;
  useCount: number;
}

export interface Message {
  id: number;
  type: "user" | "ai" | "system";
  content: string;
  timestamp: string;
}
