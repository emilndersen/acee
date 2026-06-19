import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import Lightbox from "../components/Lightbox/Lightbox";
import "../styles/album-page.css";

export default function AlbumPage() {
  const { slug } = useParams();
  const [album, setAlbum] = useState(null);
  const [photos, setPhotos] = useState([]);
  const [lightboxIndex, setLightboxIndex] = useState(-1);
  const [needsPassword, setNeedsPassword] = useState(false);
  const [password, setPassword] = useState("");
  const [pwError, setPwError] = useState("");
  const [unlocked, setUnlocked] = useState(false);
  const [reviews, setReviews] = useState([]);
  const [reviewForm, setReviewForm] = useState({ author_name: "", text: "", rating: 5 });
  const [reviewSent, setReviewSent] = useState(false);

  useEffect(() => {
    fetch(`/api/albums/${slug}`)
      .then((r) => r.json())
      .then((d) => {
        if (d.requires_password) {
          setNeedsPassword(true);
          setAlbum(d.album);
        } else {
          setAlbum(d.album);
          loadPhotos();
        }
      })
      .catch(() => {});
  }, [slug]);

  useEffect(() => {
    if (album?.id) {
      fetch(`/api/reviews/album/${album.id}`)
        .then((r) => r.json())
        .then((d) => setReviews(d.reviews || []))
        .catch(() => {});
    }
  }, [album?.id]);

  function loadPhotos() {
    fetch(`/api/albums/${slug}/photos`)
      .then((r) => r.json())
      .then((d) => setPhotos(d.photos || []))
      .catch(() => {});
  }

  async function handleUnlock(e) {
    e.preventDefault();
    setPwError("");
    const res = await fetch(`/api/albums/${slug}?password=${encodeURIComponent(password)}`);
    const d = await res.json();
    if (d.ok && !d.requires_password) {
      setAlbum(d.album);
      setNeedsPassword(false);
      setUnlocked(true);
      loadPhotos();
    } else {
      setPwError("Wrong password");
    }
  }

  async function handleReviewSubmit(e) {
    e.preventDefault();
    try {
      const res = await fetch("/api/reviews", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...reviewForm, album_id: album.id }),
      });
      if (!res.ok) throw new Error();
      setReviewSent(true);
      setReviewForm({ author_name: "", text: "", rating: 5 });
      const d = await fetch(`/api/reviews/album/${album.id}`).then((r) => r.json());
      setReviews(d.reviews || []);
    } catch {
      alert("Не удалось отправить отзыв");
    }
  }

  useEffect(() => {
    function onKey(e) {
      if (lightboxIndex < 0) return;
      if (e.key === "ArrowRight") setLightboxIndex((i) => Math.min(i + 1, photos.length - 1));
      if (e.key === "ArrowLeft") setLightboxIndex((i) => Math.max(i - 1, 0));
    }
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [lightboxIndex, photos.length]);

  if (needsPassword) {
    return (
      <section className="album-page">
        <Link to="/" className="album-back">&larr; Back</Link>
        <h1 className="album-title">{album?.title || "Private Album"}</h1>
        <div className="album-password-form">
          <p className="album-password-text">This album is password protected.</p>
          <form onSubmit={handleUnlock} className="album-password-fields">
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Enter password"
              className="album-password-input"
            />
            <button type="submit" className="album-password-btn">Unlock</button>
          </form>
          {pwError && <p className="album-password-error">{pwError}</p>}
        </div>
      </section>
    );
  }

  if (!album) return <div className="album-loading">Loading...</div>;

  const current = lightboxIndex >= 0 ? photos[lightboxIndex] : null;
  const downloadUrl = unlocked
    ? `/api/albums/${slug}/download?password=${encodeURIComponent(password)}`
    : `/api/albums/${slug}/download`;

  return (
    <section className="album-page">
      <div className="album-header">
        <Link to="/" className="album-back">&larr; Back</Link>
        {photos.length > 0 && (
          <a href={downloadUrl} className="album-download">Download All</a>
        )}
      </div>
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
        <Lightbox media={current} onClose={() => setLightboxIndex(-1)} />
      )}

      <div className="album-reviews">
        <h2 className="album-reviews__heading">Отзывы</h2>

        {reviews.length > 0 && (
          <div className="album-reviews__list">
            {reviews.map((r) => (
              <div key={r.id} className="album-reviews__card">
                <div className="album-reviews__stars">{"⭐".repeat(r.rating)}</div>
                <p className="album-reviews__text">{r.text}</p>
                <span className="album-reviews__author">— {r.author_name}</span>
              </div>
            ))}
          </div>
        )}

        {!reviewSent ? (
          <form className="album-reviews__form" onSubmit={handleReviewSubmit}>
            <h3 className="album-reviews__form-title">Оставить отзыв</h3>
            <input
              className="album-reviews__input"
              placeholder="Ваше имя"
              value={reviewForm.author_name}
              onChange={(e) => setReviewForm((f) => ({ ...f, author_name: e.target.value }))}
              required
            />
            <div className="album-reviews__rating-select">
              {[1, 2, 3, 4, 5].map((n) => (
                <button
                  key={n}
                  type="button"
                  className={`album-reviews__star-btn ${reviewForm.rating >= n ? "active" : ""}`}
                  onClick={() => setReviewForm((f) => ({ ...f, rating: n }))}
                >
                  ★
                </button>
              ))}
            </div>
            <textarea
              className="album-reviews__textarea"
              placeholder="Ваш отзыв..."
              value={reviewForm.text}
              onChange={(e) => setReviewForm((f) => ({ ...f, text: e.target.value }))}
              required
            />
            <button type="submit" className="album-reviews__submit">Отправить</button>
          </form>
        ) : (
          <p className="album-reviews__thanks">Спасибо за отзыв! ✦</p>
        )}
      </div>
    </section>
  );
}