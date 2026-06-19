import { useEffect, useState, type FormEvent } from "react";
import { Link } from "react-router-dom";
import { api } from "../api/client";

interface Album {
  id: string;
  title: string;
  slug: string;
  description: string;
  category_id: string | null;
  category_name: string | null;
  cover_url: string | null;
  is_private: boolean;
  view_count: number;
  created_at: string;
}

interface Category {
  id: string;
  name: string;
  slug: string;
}

export default function Albums() {
  const [albums, setAlbums] = useState<Album[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [title, setTitle] = useState("");
  const [slug, setSlug] = useState("");
  const [desc, setDesc] = useState("");
  const [categoryId, setCategoryId] = useState("");
  const [isPrivate, setIsPrivate] = useState(false);
  const [password, setPassword] = useState("");
  const [creating, setCreating] = useState(false);

  function load() {
    api<{ albums: Album[] }>("/albums/admin/list").then((d) => setAlbums(d.albums || []));
    api<{ categories: Category[] }>("/categories").then((d) => setCategories(d.categories || []));
  }

  useEffect(load, []);

  function autoSlug(t: string) {
    return t.toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/(^-|-$)/g, "");
  }

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    setCreating(true);
    try {
      await api("/albums", {
        method: "POST",
        body: JSON.stringify({
          title,
          slug: slug || autoSlug(title),
          description: desc,
          category_id: categoryId || null,
          is_private: isPrivate,
          password: isPrivate && password ? password : null,
        }),
      });
      setTitle("");
      setSlug("");
      setDesc("");
      setCategoryId("");
      setIsPrivate(false);
      setPassword("");
      load();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Error");
    } finally {
      setCreating(false);
    }
  }

  async function handleDelete(s: string) {
    if (!confirm("Delete this album?")) return;
    await api(`/albums/${s}`, { method: "DELETE" });
    load();
  }

  return (
    <div>
      <h1 style={{ margin: "0 0 24px", fontSize: 24 }}>Albums</h1>

      <form onSubmit={handleCreate} style={{
        background: "#1a1a1a", borderRadius: 8, padding: 24,
        display: "flex", flexDirection: "column", gap: 12, maxWidth: 500, marginBottom: 32,
      }}>
        <input placeholder="Title" value={title} onChange={(e) => { setTitle(e.target.value); if (!slug) setSlug(""); }} style={inputStyle} />
        <input placeholder="Slug (auto)" value={slug} onChange={(e) => setSlug(e.target.value)} style={inputStyle} />
        <textarea placeholder="Description" value={desc} onChange={(e) => setDesc(e.target.value)} style={{ ...inputStyle, minHeight: 60, resize: "vertical" }} />

        <select value={categoryId} onChange={(e) => setCategoryId(e.target.value)} style={inputStyle}>
          <option value="">No category</option>
          {categories.map((c) => (
            <option key={c.id} value={c.id}>{c.name}</option>
          ))}
        </select>

        <label style={{ display: "flex", alignItems: "center", gap: 8, fontSize: 13, color: "#ccc" }}>
          <input type="checkbox" checked={isPrivate} onChange={(e) => setIsPrivate(e.target.checked)} />
          Private album (password protected)
        </label>

        {isPrivate && (
          <input placeholder="Password for album" value={password} onChange={(e) => setPassword(e.target.value)} style={inputStyle} />
        )}

        <button type="submit" disabled={creating} style={btnStyle}>
          {creating ? "Creating..." : "+ Create Album"}
        </button>
      </form>

      <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(280px, 1fr))", gap: 16 }}>
        {albums.map((a) => (
          <div key={a.id} style={{
            background: "#1a1a1a", borderRadius: 8, overflow: "hidden",
            border: a.is_private ? "1px solid rgba(255,165,0,0.3)" : "1px solid #333",
          }}>
            {a.cover_url ? (
              <img src={a.cover_url} alt="" style={{ width: "100%", height: 160, objectFit: "cover" }} />
            ) : (
              <div style={{ width: "100%", height: 160, background: "#111", display: "flex", alignItems: "center", justifyContent: "center", color: "#444" }}>
                No cover
              </div>
            )}
            <div style={{ padding: 16 }}>
              <div style={{ display: "flex", justifyContent: "space-between", alignItems: "start", marginBottom: 4 }}>
                <Link to={`/albums/${a.slug}`} style={{ color: "#fff", textDecoration: "none", fontSize: 16, fontWeight: 600 }}>
                  {a.title}
                </Link>
                <span style={{ fontSize: 11, color: "#666", whiteSpace: "nowrap" }}>
                  👁 {a.view_count}
                </span>
              </div>
              <div style={{ fontSize: 12, color: "#666", marginBottom: 8 }}>/{a.slug}</div>
              {a.category_name && (
                <span style={{
                  fontSize: 10, background: "#333", color: "#aaa", padding: "2px 8px",
                  borderRadius: 4, marginBottom: 8, display: "inline-block",
                }}>
                  {a.category_name}
                </span>
              )}
              {a.is_private && (
                <span style={{
                  fontSize: 10, background: "rgba(255,165,0,0.2)", color: "#ffa500", padding: "2px 8px",
                  borderRadius: 4, marginLeft: 6, display: "inline-block",
                }}>
                  🔒 Private
                </span>
              )}
              <div style={{ marginTop: 12 }}>
                <button onClick={() => handleDelete(a.slug)} style={{
                  padding: "4px 10px", background: "transparent", color: "#f66",
                  border: "1px solid #f663", borderRadius: 4, cursor: "pointer", fontSize: 11,
                }}>
                  Delete
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

const inputStyle: React.CSSProperties = {
  padding: "10px 12px", background: "#111", border: "1px solid #333",
  borderRadius: 6, color: "#eee", fontSize: 14, outline: "none",
};

const btnStyle: React.CSSProperties = {
  padding: "8px 20px", background: "#fff", color: "#000",
  border: "none", borderRadius: 6, cursor: "pointer", fontSize: 13, fontWeight: 600,
};
