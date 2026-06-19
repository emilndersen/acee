import { useEffect, useState } from "react";
import { api } from "../api/client";

interface Booking {
  id: string;
  name: string;
  contact: string;
  telegram: string;
  shoot_type: string;
  date: string;
  shoot_date: string | null;
  idea: string;
  status: string;
  created_at: string;
}

const statuses = ["new", "confirmed", "completed", "cancelled"] as const;
const statusColors: Record<string, string> = {
  new: "#4af", confirmed: "#4f8", completed: "#999", cancelled: "#f66",
};
const DAYS = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

export default function Bookings() {
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [currentMonth, setCurrentMonth] = useState(() => {
    const d = new Date();
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}`;
  });

  function load() {
    api<{ bookings: Booking[] }>("/bookings").then((d) => setBookings(d.bookings || []));
  }
  useEffect(load, []);

  async function changeStatus(id: string, status: string) {
    await api(`/bookings/${id}/status`, { method: "PATCH", body: JSON.stringify({ status }) });
    load();
  }

  async function handleDelete(id: string) {
    if (!confirm("Delete?")) return;
    await api(`/bookings/${id}`, { method: "DELETE" });
    load();
  }

  function shiftMonth(dir: number) {
    const [y, m] = currentMonth.split("-").map(Number);
    const d = new Date(y, m - 1 + dir, 1);
    setCurrentMonth(`${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}`);
  }

  const [year, month] = currentMonth.split("-").map(Number);
  const firstDay = new Date(year, month - 1, 1);
  const daysInMonth = new Date(year, month, 0).getDate();
  let startDay = firstDay.getDay() - 1;
  if (startDay < 0) startDay = 6;

  const busyDates = new Set(
    bookings
      .filter((b) => b.shoot_date && b.status !== "cancelled")
      .map((b) => b.shoot_date!)
  );

  const calendarCells: (number | null)[] = [];
  for (let i = 0; i < startDay; i++) calendarCells.push(null);
  for (let d = 1; d <= daysInMonth; d++) calendarCells.push(d);

  function dateStr(day: number) {
    return `${currentMonth}-${String(day).padStart(2, "0")}`;
  }

  const monthName = firstDay.toLocaleString("en", { month: "long", year: "numeric" });

  return (
    <div>
      <h1 style={{ margin: "0 0 24px", fontSize: 24 }}>Bookings</h1>

      {/* Calendar */}
      <div style={{ background: "#1a1a1a", borderRadius: 8, padding: 20, marginBottom: 32, maxWidth: 420 }}>
        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 16 }}>
          <button onClick={() => shiftMonth(-1)} style={navBtn}>&larr;</button>
          <span style={{ fontSize: 15, fontWeight: 600 }}>{monthName}</span>
          <button onClick={() => shiftMonth(1)} style={navBtn}>&rarr;</button>
        </div>
        <div style={{ display: "grid", gridTemplateColumns: "repeat(7, 1fr)", gap: 4, textAlign: "center" }}>
          {DAYS.map((d) => (
            <div key={d} style={{ fontSize: 10, color: "#666", padding: 4 }}>{d}</div>
          ))}
          {calendarCells.map((day, i) => {
            if (day === null) return <div key={`e${i}`} />;
            const ds = dateStr(day);
            const isBusy = busyDates.has(ds);
            const isToday = ds === new Date().toISOString().slice(0, 10);
            return (
              <div key={day} style={{
                padding: 6, borderRadius: 4, fontSize: 13,
                background: isBusy ? "rgba(200,0,26,0.25)" : isToday ? "rgba(255,255,255,0.08)" : "transparent",
                color: isBusy ? "#f66" : isToday ? "#fff" : "#aaa",
                fontWeight: isBusy || isToday ? 700 : 400,
                border: isToday ? "1px solid #555" : "1px solid transparent",
              }}>
                {day}
              </div>
            );
          })}
        </div>
        <div style={{ marginTop: 12, fontSize: 11, color: "#666", display: "flex", gap: 16 }}>
          <span><span style={{ color: "#f66" }}>●</span> Booked</span>
          <span><span style={{ color: "#fff" }}>●</span> Today</span>
        </div>
      </div>

      {/* Bookings list */}
      {bookings.length === 0 && <p style={{ color: "#666" }}>No bookings yet.</p>}

      <div style={{ display: "flex", flexDirection: "column", gap: 12 }}>
        {bookings.map((b) => (
          <div key={b.id} style={{
            background: "#1a1a1a", borderRadius: 8, padding: 16,
            display: "flex", flexDirection: "column", gap: 8,
          }}>
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
              <span style={{ fontSize: 16, fontWeight: 600 }}>{b.name}</span>
              <span style={{
                fontSize: 12, padding: "3px 8px", borderRadius: 4,
                background: `${statusColors[b.status] || "#999"}22`,
                color: statusColors[b.status] || "#999",
              }}>
                {b.status}
              </span>
            </div>
            <div style={{ fontSize: 13, color: "#999", display: "flex", gap: 16, flexWrap: "wrap" }}>
              <span>Contact: {b.contact}</span>
              {b.telegram && <span>TG: {b.telegram}</span>}
              <span>Type: {b.shoot_type}</span>
              {b.shoot_date && <span>📅 {b.shoot_date}</span>}
              {b.date && !b.shoot_date && <span>Date: {b.date}</span>}
            </div>
            {b.idea && <p style={{ color: "#bbb", margin: 0, fontSize: 13 }}>{b.idea}</p>}
            <div style={{ display: "flex", gap: 6, flexWrap: "wrap", marginTop: 4 }}>
              {statuses.filter((s) => s !== b.status).map((s) => (
                <button key={s} onClick={() => changeStatus(b.id, s)} style={{
                  padding: "4px 10px", background: "transparent",
                  color: statusColors[s], border: `1px solid ${statusColors[s]}44`,
                  borderRadius: 4, cursor: "pointer", fontSize: 11, textTransform: "capitalize",
                }}>{s}</button>
              ))}
              <button onClick={() => handleDelete(b.id)} style={{
                marginLeft: "auto", padding: "4px 10px", background: "transparent",
                color: "#f66", border: "1px solid #f663", borderRadius: 4, cursor: "pointer", fontSize: 11,
              }}>Delete</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

const navBtn: React.CSSProperties = {
  background: "transparent", border: "1px solid #333", color: "#aaa",
  padding: "4px 12px", borderRadius: 4, cursor: "pointer", fontSize: 14,
};
