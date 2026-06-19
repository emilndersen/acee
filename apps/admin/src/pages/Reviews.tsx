import { useEffect, useState } from "react";
import { api } from "../api/client";

interface Review {
  id: string;
  author_name: string;
  text: string;
  rating: number;
  is_visible: boolean;
  created_at: string;
}

export default function Reviews() {
  const [reviews, setReviews] = useState<Review[]>([]);

  function load() {
    api<{ reviews: Review[] }>("/reviews/all").then((d) => setReviews(d.reviews || []));
  }
  useEffect(load, []);

  async function toggleVisible(id: string, visible: boolean) {
    await api(`/reviews/${id}/visible`, {
      method: "PATCH",
      body: JSON.stringify({ is_visible: visible }),
    });
    load();
  }

  async function handleDelete(id: string) {
    if (!confirm("Delete review?")) return;
    await api(`/reviews/${id}`, { method: "DELETE" });
    load();
  }

  return (
    <div>
      <h1 style={{ margin: "0 0 24px", fontSize: 24 }}>Reviews</h1>
      {reviews.length === 0 && <p style={{ color: "#666" }}>No reviews yet.</p>}

      <div style={{ display: "flex", flexDirection: "column", gap: 12 }}>
        {reviews.map((r) => (
          <div key={r.id} style={{
            background: "#1a1a1a", borderRadius: 8, padding: 16,
            border: r.is_visible ? "1px solid #4f833" : "1px solid #333",
          }}>
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 8 }}>
              <span style={{ fontWeight: 600, fontSize: 15 }}>{r.author_name}</span>
              <span style={{ fontSize: 14 }}>{"⭐".repeat(r.rating)}</span>
            </div>
            <p style={{ color: "#bbb", margin: "0 0 12px", fontSize: 14 }}>{r.text}</p>
            <div style={{ display: "flex", gap: 8 }}>
              <button
                onClick={() => toggleVisible(r.id, !r.is_visible)}
                style={{
                  padding: "4px 12px", background: "transparent", borderRadius: 4,
                  border: r.is_visible ? "1px solid #4f8" : "1px solid #666",
                  color: r.is_visible ? "#4f8" : "#666", cursor: "pointer", fontSize: 11,
                }}
              >
                {r.is_visible ? "✓ Visible" : "Hidden"}
              </button>
              <button onClick={() => handleDelete(r.id)} style={{
                padding: "4px 10px", background: "transparent", color: "#f66",
                border: "1px solid #f663", borderRadius: 4, cursor: "pointer", fontSize: 11,
              }}>
                Delete
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
