"use client";

import { useState } from "react";

type Employee = {
  name: string;
  surname: string;
  birth_date: string; // ISO date string from your API
  entry_date: string; // ISO date string from your API
};

const API_BASE: string =
  process.env.NEXT_PUBLIC_API_BASE ?? "http://api.emp.athenalabo.com";

export default function Home() {
  const [name, setName] = useState<string>("");
  const [result, setResult] = useState<Employee[] | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const search = async (): Promise<void> => {
    setError(null);
    setResult(null);

    const q = name.trim();
    if (!q) {
      setError("Please enter a name.");
      return;
    }

    try {
      setLoading(true);
      const url = new URL("/api/employee", API_BASE);
      url.searchParams.set("name", q);

      const res = await fetch(url.toString(), { method: "GET" });
      if (!res.ok) {
        const msg = await res.text().catch(() => "");
        throw new Error(`API ${res.status}${msg ? `: ${msg}` : ""}`);
      }

      const data: Employee[] = await res.json();
      setResult(data);
    } catch (e) {
      const message =
        e instanceof Error ? e.message : "Unknown error while fetching";
      setError(message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <main style={{ padding: 20, maxWidth: 720 }}>
      <h1>Employee Lookup</h1>
      <div style={{ display: "flex", gap: 8, margin: "12px 0" }}>
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Enter name (e.g. Bob)"
          style={{ flex: 1, padding: 8 }}
        />
        <button onClick={search} disabled={loading || !name.trim()}>
          {loading ? "Searching..." : "Search"}
        </button>
      </div>

      {error && <p style={{ color: "crimson" }}>{error}</p>}

      {result && (
        <pre>{JSON.stringify(result, null, 2)}</pre>
      )}
    </main>
  );
}
