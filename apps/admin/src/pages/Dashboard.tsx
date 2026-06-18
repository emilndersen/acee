import { useEffect, useState } from "react";
import { api } from "../api/client";

export default function Dashboard() {
  const [albumCount, setAlbumCount] = useState<number | null>(null);
  const [bookingCount, setBookingCount] = useState<number | null>(null);

  useEffect(() => {
    api<{ albums: unknown[] }>("/albums").then((d) =>
      setAlbumCount(d.albums.length)
    );
    api<{ bookings: unknown[] }>("/bookings").then((d) =>
      setBookingCount(d.bookings.length)
    );
  }, []);

  return (
    <div>
      <h1 style={{ margin: "0 0 24px", fontSize: 24 }}>Dashboard</h1>
      <div style={{ display: "flex", gap: 16 }}>
        <Card label="Albums" value={albumCount} />
        <Card label="Bookings" value={bookingCount} />
      </div>
    </div>
  );
}

function Card({ label, value }: { label: string; value: number | null }) {
  return (
    <div
      style={{
        background: "#1a1a1a",
        borderRadius: 8,
        padding: "24px 32px",
        minWidth: 140,
      }}
    >
      <div style={{ fontSize: 32, fontWeight: 700 }}>
        {value ?? "..."}
      </div>
      <div style={{ color: "#999", fontSize: 14, marginTop: 4 }}>{label}</div>
    </div>
  );
}