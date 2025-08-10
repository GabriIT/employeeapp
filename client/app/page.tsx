"use client";

import { useState } from "react";

const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE ?? "http://api.emp.athenalabo.com";

export default function Home() {
  const [name, setName] = useState("");
  const [result, setResult] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const search = async () => {
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

      const res = await fetch(url.toString(), {
        method: "GET",
        // cache: "no-store", // optional
      });

      if (!res.ok) {
        const msg = await res.text().catch(() => "");
        throw new Error(`API ${res.status}${msg ? `: ${msg}` : ""}`);
      }

      const data = await res.json();
      setResult(data);
    } catch (e: String | any) {
      setError(e?.message || "Request failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <main style={{ padding: 20, maxWidth: 600 }}>
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

      {error && (
        <p style={{ color: "crimson", whiteSpace: "pre-wrap" }}>{error}</p>
      )}
      {result && <pre>{JSON.stringify(result, null, 2)}</pre>}
    </main>
  );
}
