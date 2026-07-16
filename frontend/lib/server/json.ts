// src/lib/server/json.ts
"use server";
import fs from "fs/promises";
import path from "path";
import { NextResponse } from "next/server";

type JsonRecord = Record<string, string>;

async function readJsonObjectSafe(fileLoc: string): Promise<JsonRecord> {
  const baseDir = path.join(process.cwd());
  const filePath = path.join(baseDir, fileLoc);

  if (!filePath.startsWith(baseDir)) {
    throw new Error("Invalid file path");
  }

  const file = await fs.readFile(filePath, "utf-8");
  const data: unknown = JSON.parse(file);

  if (data === null || typeof data !== "object" || Array.isArray(data)) {
    throw new Error("JSON data must be a plain object");
  }

  return data as JsonRecord;
}

export async function readJsonFileSafe(fileLoc: string): Promise<JsonRecord> {
  return readJsonObjectSafe(fileLoc);
}

export async function getJsonMap(fileLoc: string): Promise<Map<string, string>> {
  const data = await readJsonObjectSafe(fileLoc);
  return new Map(Object.entries(data));
}

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const fileLoc = searchParams.get("fileLoc");

  if (!fileLoc) {
    return NextResponse.json({ message: "fileLoc is required" }, { status: 400 });
  }

  try {
    const map = await getJsonMap(fileLoc);

    return NextResponse.json({
      entries: Array.from(map.entries()),
    });
  } catch {
    return NextResponse.json({ message: "Failed to read JSON file" }, { status: 500 });
  }
}
