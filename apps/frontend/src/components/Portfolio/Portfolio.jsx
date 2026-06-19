import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import useScrollReveal from "../../hooks/useScrollReveal";
import "./Portfolio.css";

export default function Portfolio() {
  const ref = useScrollReveal();
  const [albums, setAlbums] = useState([]);
  const [categories, setCategories] = useState([]);
  const [activeCategory, setActiveCategory] = useState("");

  useEffect(() => {
    fetch("/api/categories")
      .then((r) => r.json())
      .then((d) => setCategories(d.categories || []))
      .catch(() => {});
  }, []);

  useEffect(() => {
    const url = activeCategory
      ? `/api/albums?category=${activeCategory}`
      : "/api/albums";
    fetch(url)
      .then((r) => r.json())
      .then((d) => setAlbums(d.albums || []))
      .catch(() => {});
  }, [activeCategory]);

  return (
    <section id="portfolio" className="portfolio" ref={ref}>
      <h2 className="portfolio__heading">Portfolio</h2>

      {categories.length > 0 && (
        <div className="portfolio__filters">
          <button
            className={`portfolio__filter ${activeCategory === "" ? "portfolio__filter--active" : ""}`}
            onClick={() => setActiveCategory("")}
          >
            All
          </button>
          {categories.map((c) => (
            <button
              key={c.id}
              className={`portfolio__filter ${activeCategory === c.slug ? "portfolio__filter--active" : ""}`}
              onClick={() => setActiveCategory(c.slug)}
            >
              {c.name}
            </button>
          ))}
        </div>
      )}

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
              {a.category_name && (
                <span className="portfolio__category">{a.category_name}</span>
              )}
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
