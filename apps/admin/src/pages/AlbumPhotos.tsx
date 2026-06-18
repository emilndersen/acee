import { useEffect, useState, useRef } from "react";
import { useParams } from "react-router-dom";
import { api, uploadFile } from "../api/client";

interface Photo {
  id: string;
  album_id: string;
  title: string;
  description: string;
  image_url: string;
  thumb_url: string;
  sort_order: number;
  created_at: string;
}

interface Album {
  id: string;
  slug: string;
  title: string;
}

export default function AlbumPhotos() {
  const { slug } = useParams<{ slug: string }>();
  const [album, setAlbum] = useState<Album | null>(null);
  const [photos, setPhotos] = useState<Photo[]>([]);
  const [uploading, setUploading] = useState(false);
  const fileRef = useRef<HTMLInputElement>(null);

  function load() {
    if (!slug) return;
    api<{ album: Album }>(`/albums/${slug}`).then((d) => setAlbum(d.album));
    api<{ photos: Photo[] }>(`/albums/${slug}/photos`).then((d) =>
      setPhotos(d.photos)
    );
  }

  useEffect(load, [slug]);

  async function handleUpload(files: FileList | null) {
    if (!files || !slug) return;
    setUploading(true);
    try {
      for (const file of Array.from(files)) {
        const { url, thumb_url } = await uploadFile(file);
        await api(`/albums/${slug}/photos`, {
          method: "POST",
          body: JSON.stringify({
            title: file.name.replace(/\.[^.]+$/, ""),
            image_url: url,
            thumb_url: thumb_url,
          }),
        });
      }
      load();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Upload failed");
    } finally {
      setUploading(false);
      if (fileRef.current) fileRef.current.value = "";
    }
  }

  async function handleDelete(id: string) {
    if (!confirm("Delete this photo?")) return;
    await api(`/photos/${id}`, { method: "DELETE" });
    load();
  }

  if (!album) return <p style={{ color: "#999" }}>Loading...</p>;

  return (
    <div>
      <h1 style={{ margin: "0 0 8px", fontSize: 24 }}>{album.title}</h1>
      <p style={{ color: "#666", margin: "0 0 24px", fontSize: 13 }}>/{album.slug}</p>

      <div style={{ marginBottom: 24 }}>
        <input
          ref={fileRef}
          type="file"
          multiple
          accept="image/*"
          onChange={(e) => handleUpload(e.target.files)}
          style={{ display: "none" }}
        />
        <button
          onClick={() => fileRef.current?.click()}
          disabled={uploading}
          style={{
            padding: "8px 16px",
            background: "#fff",
            color: "#000",
            border: "none",
            borderRadius: 6,
            cursor: uploading ? "wait" : "pointer",
            fontSize: 13,
            fontWeight: 600,
          }}
        >
          {uploading ? "Uploading..." : "+ Upload Photos"}
        </button>
      </div>

      <div
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(auto-fill, minmax(180px, 1fr))",
          gap: 12,
        }}
      >
        {photos.map((p) => (
          <div
            key={p.id}
            style={{
              background: "#1a1a1a",
              borderRadius: 8,
              overflow: "hidden",
            }}
          >
            <img
              src={p.thumb_url || p.image_url}
              alt={p.title}
              style={{ width: "100%", height: 160, objectFit: "cover" }}
            />
            <div style={{ padding: "8px 10px" }}>
              <div
                style={{
                  fontSize: 13,
                  color: "#eee",
                  overflow: "hidden",
                  textOverflow: "ellipsis",
                  whiteSpace: "nowrap",
                }}
              >
                {p.title || "Untitled"}
              </div>
              <button
                onClick={() => handleDelete(p.id)}
                style={{
                  marginTop: 6,
                  padding: "4px 8px",
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

      {photos.length === 0 && (
        <p style={{ color: "#666" }}>No photos yet. Upload some!</p>
      )}
    </div>
  );
}
