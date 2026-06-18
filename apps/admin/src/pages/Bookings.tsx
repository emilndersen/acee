import { useEffect, useState } from "react";
import { api } from "../api/client";

interface Booking {
  id: string;
  name: string;
  contact: string;
  shoot_type: string;
  date: string;
  idea: string;
  status: string;
  created_at: string;
}

const statuses = ["new", "confirmed", "completed", "cancelled"] as const;

const statusColors: Record<string, string> = {
  new: "#4af",
  confirmed: "#4f8",
  completed: "#999",
  cancelled: "#f66",
};

export default function Bookings() {
  const [bookings, setBookings] = useState<Booking[]>([]);

  function load() {
    api<{ bookings: Booking[] }>("/bookings").then((d) =>
      setBookings(d.bookings)
    );
  }

  useEffect(load, []);

  async function changeStatus(id: string, status: string) {
    await api(`/bookings/${id}/status`, {
      method: "PATCH",
      body: JSON.stringify({ status }),
    });
    load();
  }

  async function handleDelete(id: string) {
    if (!confirm("Delete this booking?")) return;
    await api(`/bookings/${id}`, { method: "DELETE" });
    load();
  }

  return (
    <div>
      <h1 style={{ margin: "0 0 24px", fontSize: 24 }}>Bookings</h1>

      {bookings.length === 0 && (
        <p style={{ color: "#666" }}>No bookings yet.</p>
      )}

      <div style={{ display: "flex", flexDirection: "column", gap: 12 }}>
        {bookings.map((b) => (
          <div
            key={b.id}
            style={{
              background: "#1a1a1a",
              borderRadius: 8,
              padding: 16,
              display: "flex",
              flexDirection: "column",
              gap: 8,
            }}
          >
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
              <span style={{ fontSize: 16, fontWeight: 600 }}>{b.name}</span>
              <span
                style={{
                  fontSize: 12,
                  padding: "3px 8px",
                  borderRadius: 4,
                  background: `${statusColors[b.status] || "#999"}22`,
                  color: statusColors[b.status] || "#999",
                }}
              >
                {b.status}
              </span>
            </div>

            <div style={{ fontSize: 13, color: "#999", display: "flex", gap: 16, flexWrap: "wrap" }}>
              <span>Contact: {b.contact}</span>
              <span>Type: {b.shoot_type}</span>
              {b.date && <span>Date: {b.date}</span>}
            </div>

            {b.idea && (
              <p style={{ color: "#bbb", margin: 0, fontSize: 13 }}>{b.idea}</p>
            )}

            <div style={{ display: "flex", gap: 6, flexWrap: "wrap", marginTop: 4 }}>
              {statuses
                .filter((s) => s !== b.status)
                .map((s) => (
                  <button
                    key={s}
                    onClick={() => changeStatus(b.id, s)}
                    style={{
                      padding: "4px 10px",
                      background: "transparent",
                      color: statusColors[s],
                      border: `1px solid ${statusColors[s]}44`,
                      borderRadius: 4,
                      cursor: "pointer",
                      fontSize: 11,
                      textTransform: "capitalize",
                    }}
                  >
                    {s}
                  </button>
                ))}
              <button
                onClick={() => handleDelete(b.id)}
                style={{
                  marginLeft: "auto",
                  padding: "4px 10px",
                  background: "transparent",
                  color: "#f66",
                  border: "1px solid #f663",
                  borderRadius: 4,
                  cursor: "pointer",
                  fontSize: 11,
                }}
              >
                Delete
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
