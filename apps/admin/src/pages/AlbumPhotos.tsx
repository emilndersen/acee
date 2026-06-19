import { useEffect, useState, useRef, type DragEvent } from "react";
import { useParams } from "react-router-dom";
import { api, uploadFile } from "../api/client";

interface Photo {
  id: string;
  album_id: string;
  title: string;
  image_url: string;
  thumb_url: string;
  media_type: string;
  sort_order: number;
}

interface Album {
  id: string;
  slug: string;
  title: string;
  cover_photo_id: string | null;
}

export default function AlbumPhotos() {
  const { slug } = useParams<{ slug: string }>();
  const [album, setAlbum] = useState<Album | null>(null);
  const [photos, setPhotos] = useState<Photo[]>([]);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState("");
  const [dragging, setDragging] = useState(false);
  const [dragOverId, setDragOverId] = useState<string | null>(null);
  const dragItem = useRef<string | null>(null);
  const fileRef = useRef<HTMLInputElement>(null);

  function load() {
    if (!slug) return;
    api<{ album: Album }>(`/albums/${slug}`).then((d) => setAlbum(d.album));
    api<{ photos: Photo[] }>(`/albums/${slug}/photos`).then((d) => setPhotos(d.photos || []));
  }

  useEffect(load, [slug]);

  async function doUpload(files: File[]) {
    if (!files.length || !slug) return;
    setUploading(true);
    try {
      for (let i = 0; i < files.length; i++) {
        const file = files[i];
        setUploadProgress(`Uploading ${i + 1}/${files.length}: ${file.name}`);
        const { url, thumb_url, media_type } = await uploadFile(file);
        await api(`/albums/${slug}/photos`, {
          method: "POST",
          body: JSON.stringify({
            title: file.name.replace(/\.[^.]+$/, ""),
            image_url: url,
            thumb_url: thumb_url,
            media_type: media_type,
          }),
        });
      }
      load();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Upload failed");
    } finally {
      setUploading(false);
      setUploadProgress("");
      if (fileRef.current) fileRef.current.value = "";
    }
  }

  function handleDrop(e: DragEvent) {
    e.preventDefault();
    setDragging(false);
    const files = Array.from(e.dataTransfer.files);
    if (files.length) doUpload(files);
  }

  function handleDragOver(e: DragEvent) {
    e.preventDefault();
    setDragging(true);
  }

  function handleDragLeave(e: DragEvent) {
    if (e.currentTarget.contains(e.relatedTarget as Node)) return;
    setDragging(false);
  }

  async function handleReorder(fromId: string, toId: string) {
    if (fromId === toId) return;
    const newPhotos = [...photos];
    const fromIdx = newPhotos.findIndex((p) => p.id === fromId);
    const toIdx = newPhotos.findIndex((p) => p.id === toId);
    const [moved] = newPhotos.splice(fromIdx, 1);
    newPhotos.splice(toIdx, 0, moved);
    setPhotos(newPhotos);

    for (let i = 0; i < newPhotos.length; i++) {
      if (newPhotos[i].sort_order !== i) {
        await api(`/photos/${newPhotos[i].id}/order`, {
          method: "PATCH",
          body: JSON.stringify({ sort_order: i }),
        });
      }
    }
  }

  async function handleSetCover(photoId: string) {
    if (!slug) return;
    await api(`/albums/${slug}`, {
      method: "PUT",
      body: JSON.stringify({ cover_photo_id: photoId }),
    });
    load();
  }

  async function handleDelete(id: string) {
    if (!confirm("Delete this media?")) return;
    await api(`/photos/${id}`, { method: "DELETE" });
    load();
  }

  if (!album) return <p style={{ color: "#999" }}>Loading...</p>;

  return (
    <div
      onDrop={handleDrop}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      style={{ minHeight: "80vh", position: "relative" }}
    >
      {dragging && (
        <div style={{
          position: "fixed", inset: 0, background: "rgba(200,0,26,0.15)",
          border: "3px dashed #c8001a", zIndex: 100,
          display: "flex", alignItems: "center", justifyContent: "center",
          fontSize: 24, color: "#fff", fontWeight: 700, pointerEvents: "none",
        }}>
          Drop files to upload
        </div>
      )}

      <h1 style={{ margin: "0 0 8px", fontSize: 24 }}>{album.title}</h1>
      <p style={{ color: "#666", margin: "0 0 24px", fontSize: 13 }}>/{album.slug}</p>

      <div style={{ marginBottom: 24, display: "flex", alignItems: "center", gap: 12 }}>
        <input
          ref={fileRef}
          type="file"
          multiple
          accept="image/*,video/*,.mp4,.mov,.webm,.avi,.mkv"
          onChange={(e) => doUpload(Array.from(e.target.files || []))}
          style={{ display: "none" }}
        />
        <button onClick={() => fileRef.current?.click()} disabled={uploading} style={btnStyle}>
          {uploading ? "Uploading..." : "+ Upload Media"}
        </button>
        {uploadProgress && <span style={{ fontSize: 12, color: "#999" }}>{uploadProgress}</span>}
        <span style={{ fontSize: 11, color: "#555", marginLeft: "auto" }}>
          Drag & drop files anywhere on this page
        </span>
      </div>

      <div style={{
        display: "grid",
        gridTemplateColumns: "repeat(auto-fill, minmax(180px, 1fr))",
        gap: 12,
      }}>
        {photos.map((p) => {
          const isCover = album.cover_photo_id === p.id;
          return (
            <div
              key={p.id}
              draggable
              onDragStart={() => { dragItem.current = p.id; }}
              onDragOver={(e) => { e.preventDefault(); setDragOverId(p.id); }}
              onDrop={(e) => {
                e.preventDefault();
                e.stopPropagation();
                setDragOverId(null);
                if (dragItem.current) handleReorder(dragItem.current, p.id);
                dragItem.current = null;
              }}
              onDragEnd={() => { dragItem.current = null; setDragOverId(null); }}
              style={{
                background: "#1a1a1a",
                borderRadius: 8,
                overflow: "hidden",
                cursor: "grab",
                border: isCover ? "2px solid #4f8" : dragOverId === p.id ? "2px solid #c8001a" : "2px solid transparent",
                transition: "border-color 0.2s",
              }}
            >
              {p.media_type === "video" ? (
                <div style={{ position: "relative" }}>
                  <video src={p.image_url} style={{ width: "100%", height: 160, objectFit: "cover" }} muted preload="metadata" />
                  <div style={{
                    position: "absolute", top: 8, right: 8, background: "rgba(200,0,26,0.85)",
                    color: "#fff", padding: "2px 8px", borderRadius: 4, fontSize: 10, fontWeight: 700,
                  }}>VIDEO</div>
                </div>
              ) : (
                <img src={p.thumb_url || p.image_url} alt={p.title} style={{ width: "100%", height: 160, objectFit: "cover" }} />
              )}
              <div style={{ padding: "8px 10px" }}>
                <div style={{ fontSize: 13, color: "#eee", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>
                  {p.title || "Untitled"}
                </div>
                <div style={{ display: "flex", gap: 6, marginTop: 6, flexWrap: "wrap" }}>
                  {!isCover && (
                    <button onClick={() => handleSetCover(p.id)} style={smallBtnStyle}>
                      Set Cover
                    </button>
                  )}
                  {isCover && (
                    <span style={{ fontSize: 10, color: "#4f8", fontWeight: 600 }}>★ Cover</span>
                  )}
                  <button onClick={() => handleDelete(p.id)} style={{ ...smallBtnStyle, color: "#f66", borderColor: "#f663" }}>
                    Delete
                  </button>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {photos.length === 0 && (
        <p style={{ color: "#666" }}>No media yet. Upload or drag & drop files!</p>
      )}
    </div>
  );
}

const btnStyle: React.CSSProperties = {
  padding: "8px 16px", background: "#fff", color: "#000",
  border: "none", borderRadius: 6, cursor: "pointer", fontSize: 13, fontWeight: 600,
};

const smallBtnStyle: React.CSSProperties = {
  padding: "3px 8px", background: "transparent", color: "#999",
  border: "1px solid #333", borderRadius: 4, cursor: "pointer", fontSize: 10,
};
