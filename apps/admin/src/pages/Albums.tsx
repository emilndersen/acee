import { useEffect, useState, type FormEvent } from "react";
import { Link } from "react-router-dom";
import { api } from "../api/client";

interface Album {
  id: string;
  slug: string;
  title: string;
  cover_url: string;
  description: string;
  is_public: boolean;
  sort_order: number;
  created_at: string;
}

export default function Albums() {
  const [albums, setAlbums] = useState<Album[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [slug, setSlug] = useState("");
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [error, setError] = useState("");

  function load() {
    api<{ albums: Album[] }>("/albums").then((d) => setAlbums(d.albums));
  }

  useEffect(load, []);

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    setError("");
    try {
      await api("/albums", {
        method: "POST",
        body: JSON.stringify({ slug, title, description, is_public: true }),
      });
      setSlug("");
      setTitle("");
      setDescription("");
      setShowForm(false);
      load();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Error");
    }
  }

  async function handleDelete(albumSlug: string) {
    if (!confirm("Delete this album?")) return;
    await api(`/albums/${albumSlug}`, { method: "DELETE" });
    load();
  }

  return (
    <div>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 24 }}>
        <h1 style={{ margin: 0, fontSize: 24 }}>Albums</h1>
        <button onClick={() => setShowForm(!showForm)} style={btnStyle}>
          {showForm ? "Cancel" : "+ New Album"}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleCreate} style={{ background: "#1a1a1a", padding: 20, borderRadius: 8, marginBottom: 24, display: "flex", flexDirection: "column", gap: 10 }}>
          {error && <p style={{ color: "#f66", margin: 0 }}>{error}</p>}
          <input placeholder="Slug (url-friendly)" value={slug} onChange={(e) => setSlug(e.target.value)} required style={inputStyle} />
          <input placeholder="Title" value={title} onChange={(e) => setTitle(e.target.value)} required style={inputStyle} />
          <textarea placeholder="Description" value={description} onChange={(e) => setDescription(e.target.value)} style={{ ...inputStyle, minHeight: 60, resize: "vertical" }} />
          <button type="submit" style={{ ...btnStyle, alignSelf: "flex-start" }}>Create</button>
        </form>
      )}

      <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(280px, 1fr))", gap: 16 }}>
        {albums.map((a) => (
          <div key={a.id} style={{ background: "#1a1a1a", borderRadius: 8, padding: 16, display: "flex", flexDirection: "column", gap: 8 }}>
            {a.cover_url && (
              <img src={a.cover_url} alt="" style={{ width: "100%", height: 160, objectFit: "cover", borderRadius: 6 }} />
            )}
            <Link to={`/albums/${a.slug}`} style={{ color: "#fff", textDecoration: "none", fontSize: 18, fontWeight: 600 }}>
              {a.title}
            </Link>
            <span style={{ color: "#666", fontSize: 12 }}>/{a.slug}</span>
            {a.description && <p style={{ color: "#999", margin: 0, fontSize: 13 }}>{a.description}</p>}
            <div style={{ display: "flex", gap: 8, marginTop: "auto" }}>
              <Link to={`/albums/${a.slug}`} style={{ ...btnSmall, textDecoration: "none", textAlign: "center" }}>
                Photos
              </Link>
              <button onClick={() => handleDelete(a.slug)} style={{ ...btnSmall, color: "#f66", borderColor: "#f663" }}>
                Delete
              </button>
            </div>
          </div>
        ))}
      </div>

      {albums.length === 0 && !showForm && (
        <p style={{ color: "#666" }}>No albums yet. Create one to get started.</p>
      )}
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

const btnStyle: React.CSSProperties = {
  padding: "8px 16px",
  background: "#fff",
  color: "#000",
  border: "none",
  borderRadius: 6,
  cursor: "pointer",
  fontSize: 13,
  fontWeight: 600,
};

const btnSmall: React.CSSProperties = {
  padding: "6px 12px",
  background: "transparent",
  color: "#eee",
  border: "1px solid #333",
  borderRadius: 6,
  cursor: "pointer",
  fontSize: 12,
};