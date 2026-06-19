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
  media_type: string;
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
  const [uploadProgress, setUploadProgress] = useState("");
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
    const total = files.length;
    try {
      for (let i = 0; i < total; i++) {
        const file = files[i];
        setUploadProgress(`Uploading ${i + 1}/${total}: ${file.name}`);
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

  async function handleDelete(id: string) {
    if (!confirm("Delete this media?")) return;
    await api(`/photos/${id}`, { method: "DELETE" });
    load();
  }

  function isVideo(p: Photo) {
    return p.media_type === "video";
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
          accept="image/*,video/*,.mp4,.mov,.webm,.avi,.mkv"
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
          {uploading ? "Uploading..." : "+ Upload Media"}
        </button>
        {uploadProgress && (
          <span style={{ marginLeft: 12, fontSize: 12, color: "#999" }}>
            {uploadProgress}
          </span>
        )}
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
              position: "relative",
            }}
          >
            {isVideo(p) ? (
              <div style={{ position: "relative" }}>
                <video
                  src={p.image_url}
                  style={{ width: "100%", height: 160, objectFit: "cover" }}
                  muted
                  preload="metadata"
                />
                <div style={{
                  position: "absolute",
                  top: 8,
                  right: 8,
                  background: "rgba(200,0,26,0.85)",
                  color: "#fff",
                  padding: "2px 8px",
                  borderRadius: 4,
                  fontSize: 10,
                  fontWeight: 700,
                }}>
                  VIDEO
                </div>
              </div>
            ) : (
              <img
                src={p.thumb_url || p.image_url}
                alt={p.title}
                style={{ width: "100%", height: 160, objectFit: "cover" }}
              />
            )}
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
        <p style={{ color: "#666" }}>No media yet. Upload some!</p>
      )}
    </div>
  );
}
