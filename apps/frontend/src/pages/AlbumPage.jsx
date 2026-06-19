import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import Lightbox from "../components/Lightbox/Lightbox";
import "../styles/album-page.css";

export default function AlbumPage() {
  const { slug } = useParams();
  const [album, setAlbum] = useState(null);
  const [photos, setPhotos] = useState([]);
  const [lightboxIndex, setLightboxIndex] = useState(-1);

  useEffect(() => {
    fetch(`/api/albums/${slug}`)
      .then((r) => r.json())
      .then((d) => setAlbum(d.album))
      .catch(() => {});

    fetch(`/api/albums/${slug}/photos`)
      .then((r) => r.json())
      .then((d) => setPhotos(d.photos || []))
      .catch(() => {});
  }, [slug]);

  useEffect(() => {
    function onKey(e) {
      if (lightboxIndex < 0) return;
      if (e.key === "ArrowRight") setLightboxIndex((i) => Math.min(i + 1, photos.length - 1));
      if (e.key === "ArrowLeft") setLightboxIndex((i) => Math.max(i - 1, 0));
    }
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [lightboxIndex, photos.length]);

  if (!album) return <div className="album-loading">Loading...</div>;

  const current = lightboxIndex >= 0 ? photos[lightboxIndex] : null;

  return (
    <section className="album-page">
      <Link to="/" className="album-back">&larr; Back</Link>
      <h1 className="album-title">{album.title}</h1>
      {album.description && <p className="album-desc">{album.description}</p>}

      <div className="album-grid">
        {photos.map((p, i) => (
          <div key={p.id} className="album-grid__item" onClick={() => setLightboxIndex(i)}>
            {p.media_type === "video" ? (
              <div className="album-grid__video-wrap">
                <video src={p.image_url} muted preload="metadata" />
                <div className="album-grid__play">&#9654;</div>
              </div>
            ) : (
              <img src={p.thumb_url || p.image_url} alt="" loading="lazy" />
            )}
          </div>
        ))}
      </div>

      {photos.length === 0 && <p className="album-empty">No media yet.</p>}

      {current && (
        <Lightbox
          media={current}
          onClose={() => setLightboxIndex(-1)}
        />
      )}
    </section>
  );
}
