import { useState, useEffect, useRef } from 'react'
import { Link } from 'react-router-dom'
import { useLang } from '../../i18n/LangContext'
import './Portfolio.css'

export default function Portfolio() {
  const { t } = useLang()
  const [albums, setAlbums] = useState([])
  const [categories, setCategories] = useState([])
  const [activeCategory, setActiveCategory] = useState(null)
  const gridRef = useRef(null)

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
    if (!gridRef.current) return
    const items = gridRef.current.querySelectorAll('.portfolio__card')
    const observers = []
    items.forEach((el, i) => {
      el.classList.remove('visible')
      const obs = new IntersectionObserver(
        ([entry]) => {
          if (entry.isIntersecting) {
            setTimeout(() => el.classList.add('visible'), i * 100)
            obs.unobserve(el)
          }
        },
        { threshold: 0.05 }
      )
      obs.observe(el)
      observers.push(obs)
    })
    return () => observers.forEach((o) => o.disconnect())
  }, [albums, activeCategory])

  const filtered = activeCategory
    ? albums.filter((a) => a.category_id === activeCategory)
    : albums

  return (
    <section className="portfolio" id="portfolio">
      <h2 className="portfolio__heading">{t.portfolio.heading}</h2>

      {categories.length > 0 && (
        <div className="portfolio__filters">
          <button
            className={`portfolio__filter ${!activeCategory ? 'portfolio__filter--active' : ''}`}
            onClick={() => setActiveCategory(null)}
          >
            {t.portfolio.all}
          </button>
          {categories.map((cat) => (
            <button
              key={cat.id}
              className={`portfolio__filter ${activeCategory === cat.id ? 'portfolio__filter--active' : ''}`}
              onClick={() => setActiveCategory(cat.id)}
            >
              {cat.name}
            </button>
          ))}
        </div>
      )}

      <div className="portfolio__grid" ref={gridRef}>
        {filtered.map((album) => (
          <Link
            key={album.id}
            to={`/albums/${album.slug}`}
            className="portfolio__card"
          >
            {album.cover_url ? (
              <img src={album.cover_url} alt={album.title} className="portfolio__img" loading="lazy" />
            ) : (
              <div className="portfolio__img portfolio__img--empty" />
            )}
            <div className="portfolio__overlay">
              <span className="portfolio__label">{album.title}</span>
            </div>
          </Link>
        ))}
      </div>

      {filtered.length === 0 && (
        <p className="portfolio__empty">{t.portfolio.empty}</p>
      )}
    </section>
  )
}