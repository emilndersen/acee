import { useState, useEffect, useRef } from 'react'
import { Link } from 'react-router-dom'
import './Portfolio.css'

export default function Portfolio() {
  const [albums, setAlbums] = useState([])
  const [categories, setCategories] = useState([])
  const [activeCategory, setActiveCategory] = useState(null)
  const itemRefs = useRef([])

  useEffect(() => {
    fetch('/api/albums')
      .then((r) => r.json())
      .then((d) => setAlbums(d.albums || []))
      .catch(() => {})

    fetch('/api/categories')
      .then((r) => r.json())
      .then((d) => setCategories(d.categories || []))
      .catch(() => {})
  }, [])

  useEffect(() => {
    const observers = itemRefs.current.map((el, i) => {
      if (!el) return null
      const obs = new IntersectionObserver(
        ([entry]) => {
          if (entry.isIntersecting) {
            setTimeout(() => el.classList.add('visible'), i * 80)
            obs.unobserve(el)
          }
        },
        { threshold: 0.05 }
      )
      obs.observe(el)
      return obs
    })
    return () => observers.forEach((o) => o?.disconnect())
  }, [albums, activeCategory])

  const filtered = activeCategory
    ? albums.filter((a) => a.category_id === activeCategory)
    : albums

  return (
    <section className="portfolio" id="portfolio">
      <div className="section-header">
        <h2 className="section-title">ПОРТФОЛИО</h2>
        <span className="section-sub">Нажми на альбом → посмотреть фото</span>
      </div>

      {categories.length > 0 && (
        <div className="category-filters">
          <button
            className={`category-btn ${!activeCategory ? 'active' : ''}`}
            onClick={() => setActiveCategory(null)}
          >
            Все
          </button>
          {categories.map((cat) => (
            <button
              key={cat.id}
              className={`category-btn ${activeCategory === cat.id ? 'active' : ''}`}
              onClick={() => setActiveCategory(cat.id)}
            >
              {cat.name}
            </button>
          ))}
        </div>
      )}

      <div className="collage-grid">
        {filtered.map((album, i) => (
          <Link
            key={album.id}
            to={`/albums/${album.slug}`}
            className="collage-item"
            ref={(el) => (itemRefs.current[i] = el)}
          >
            {album.cover_url ? (
              <img src={album.cover_url} alt={album.title} loading="lazy" />
            ) : (
              <div className="collage-placeholder" />
            )}
            <div className="collage-overlay">
              <div className="collage-info">
                <div className="collage-info-title">{album.title}</div>
                <span className="collage-info-link">Смотреть альбом →</span>
              </div>
            </div>
          </Link>
        ))}
      </div>

      {filtered.length === 0 && (
        <p style={{ textAlign: 'center', color: '#666', padding: '40px 0' }}>
          Пока нет альбомов
        </p>
      )}
    </section>
  )
}
