import { useEffect, useState, type FormEvent } from "react";
import { api } from "../api/client";

const settingKeys = [
  { key: "site_title", label: "Site Title" },
  { key: "about_text", label: "About Text", multiline: true },
  { key: "phone", label: "Phone" },
  { key: "email", label: "Email" },
  { key: "telegram", label: "Telegram" },
  { key: "instagram", label: "Instagram" },
  { key: "vk", label: "VK" },
  { key: "location", label: "Location" },
];

export default function Settings() {
  const [values, setValues] = useState<Record<string, string>>({});
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState("");

  useEffect(() => {
    api<{ settings: Record<string, string> }>("/settings").then(
      (d) => {
        setValues(d.settings || {});
      }
    );
  }, []);

  async function handleSave(e: FormEvent) {
    e.preventDefault();
    setSaving(true);
    setMessage("");
    try {
      await api("/settings", {
        method: "PUT",
        body: JSON.stringify({ settings: values }),
      });
      setMessage("Saved!");
      setTimeout(() => setMessage(""), 2000);
    } catch (err: unknown) {
      setMessage(err instanceof Error ? err.message : "Error");
    } finally {
      setSaving(false);
    }
  }

  return (
    <div>
      <h1 style={{ margin: "0 0 24px", fontSize: 24 }}>Settings</h1>

      <form
        onSubmit={handleSave}
        style={{
          background: "#1a1a1a",
          borderRadius: 8,
          padding: 24,
          display: "flex",
          flexDirection: "column",
          gap: 14,
          maxWidth: 500,
        }}
      >
        {settingKeys.map((s) => (
          <label key={s.key} style={{ display: "flex", flexDirection: "column", gap: 4 }}>
            <span style={{ fontSize: 13, color: "#999" }}>{s.label}</span>
            {s.multiline ? (
              <textarea
                value={values[s.key] || ""}
                onChange={(e) =>
                  setValues((v) => ({ ...v, [s.key]: e.target.value }))
                }
                style={{ ...inputStyle, minHeight: 80, resize: "vertical" }}
              />
            ) : (
              <input
                value={values[s.key] || ""}
                onChange={(e) =>
                  setValues((v) => ({ ...v, [s.key]: e.target.value }))
                }
                style={inputStyle}
              />
            )}
          </label>
        ))}

        <div style={{ display: "flex", alignItems: "center", gap: 12 }}>
          <button
            type="submit"
            disabled={saving}
            style={{
              padding: "8px 20px",
              background: "#fff",
              color: "#000",
              border: "none",
              borderRadius: 6,
              cursor: saving ? "wait" : "pointer",
              fontSize: 13,
              fontWeight: 600,
            }}
          >
            {saving ? "Saving..." : "Save"}
          </button>
          {message && (
            <span style={{ fontSize: 13, color: message === "Saved!" ? "#4f8" : "#f66" }}>
              {message}
            </span>
          )}
        </div>
      </form>
    </div>
  );
}

const inputStyle: React.CSSProperties = {
  padding: "10px 12px",
  background: "#111",
  border: "1px solid #333",
  borderRadius: 6,
  color: "#eee",
  fontSize: 14,
  outline: "none",
};
