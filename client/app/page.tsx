import Image from "next/image";
import styles from "./page.module.css";

"use client";
const base = process.env.NEXT_PUBLIC_API_BASE!; // will be inlined at build-time

import { useState } from "react";

export default function Home() {
  const [name, setName] = useState("");
  const [result, setResult] = useState(null);

  const search = async () => {
    const res = await fetch(`${base}/api/employee?name=${encodeURIComponent(name)}`);

    // const res = await fetch(`/api/employee?name=${name}`);
    const data = await res.json();
    setResult(data);
  };

  return (
    <main style={{ padding: 20 }}>
      <h1>Employee Lookup</h1>
      <input value={name} onChange={(e) => setName(e.target.value)} placeholder="Enter name" />
      <button onClick={search}>Search</button>
      <pre>{JSON.stringify(result, null, 2)}</pre>
    </main>
  );
}



