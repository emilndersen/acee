import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { clearTokens } from "../api/client";

const links = [
  { to: "/", label: "Dashboard" },
  { to: "/albums", label: "Albums" },
  { to: "/bookings", label: "Bookings" },
  { to: "/reviews", label: "Reviews" },
  { to: "/settings", label: "Settings" },
];

export default function Layout() {
  const navigate = useNavigate();

  function handleLogout() {
    clearTokens();
    navigate("/login");
  }

  return (
    <div style={{ display: "flex", minHeight: "100vh" }}>
      <aside
        style={{
          width: 220,
          background: "#0d0d0d",
          padding: "24px 16px",
          display: "flex",
          flexDirection: "column",
          gap: 8,
        }}
      >
        <h2 style={{ color: "#fff", margin: "0 0 24px", fontSize: 20 }}>
          Acee Admin
        </h2>
        {links.map((l) => (
          <NavLink
            key={l.to}
            to={l.to}
            end={l.to === "/"}
            style={({ isActive }) => ({
              display: "block",
              padding: "10px 12px",
              borderRadius: 6,
              color: isActive ? "#fff" : "#999",
              background: isActive ? "#1a1a1a" : "transparent",
              textDecoration: "none",
              fontSize: 14,
            })}
          >
            {l.label}
          </NavLink>
        ))}
        <div style={{ flex: 1 }} />
        <button
          onClick={handleLogout}
          style={{
            background: "none",
            border: "1px solid #333",
            color: "#999",
            padding: "8px 12px",
            borderRadius: 6,
            cursor: "pointer",
            fontSize: 13,
          }}
        >
          Logout
        </button>
      </aside>
      <main style={{ flex: 1, background: "#111", padding: 32, color: "#eee" }}>
        <Outlet />
      </main>
    </div>
  );
}
