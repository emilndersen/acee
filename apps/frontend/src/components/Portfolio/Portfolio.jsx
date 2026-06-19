import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import useScrollReveal from "../../hooks/useScrollReveal";
import "./Portfolio.css";

export default function Portfolio() {
  const ref = useScrollReveal();
  const [albums, setAlbums] = useState([]);

  useEffect(() => {
    fetch("/api/albums")
      .then((r) => r.json())
      .then((d) => setAlbums(d.albums || []))
      .catch(() => {});
  }, []);

  return (
    <section id="portfolio" className="portfolio" ref={ref}>
      <h2 className="portfolio__heading">Portfolio</h2>
      <div className="portfolio__grid">
        {albums.map((a) => (
          <Link to={`/albums/${a.slug}`} key={a.id} className="portfolio__card">
            {a.cover_url ? (
              <img src={a.cover_url} alt={a.title} className="portfolio__img" loading="lazy" />
            ) : (
              <div className="portfolio__img portfolio__img--empty" />
            )}
            <div className="portfolio__overlay">
              <span className="portfolio__label">{a.title}</span>
            </div>
          </Link>
        ))}
      </div>
      {albums.length === 0 && (
        <p style={{ color: "#666", textAlign: "center", padding: "40px 0" }}>
          No albums yet
        </p>
      )}
    </section>
  );
}
