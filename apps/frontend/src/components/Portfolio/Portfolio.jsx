import { useState, useEffect, useRef } from 'react'
import Lightbox from '../Lightbox/Lightbox'
import { PHOTOS } from '../../data/photos'
import './Portfolio.css'

export default function Portfolio() {
  const [activePhoto, setActivePhoto] = useState(null)
  const itemRefs = useRef([])

  // Staggered scroll reveal for each collage item
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
  }, [])

  return (
    <section className="portfolio" id="portfolio">
      <div className="section-header reveal-once">
        <h2 className="section-title">ПОРТФОЛИО</h2>
        <span className="section-sub">Нажми на фото → перейти к альбому</span>
      </div>

      <div className="collage-grid">
        {PHOTOS.map((photo, i) => (
          <div
            key={photo.id}
            className="collage-item"
            ref={(el) => (itemRefs.current[i] = el)}
            onClick={() => setActivePhoto(photo)}
          >
            <img src={photo.thumb} alt={photo.title} loading="lazy" />
            <div className="collage-overlay">
              <div className="collage-info">
                <div className="collage-info-title">{photo.album}</div>
                <span className="collage-info-link">Перейти к альбому →</span>
              </div>
            </div>
          </div>
        ))}
      </div>

      {activePhoto && (
        <Lightbox photo={activePhoto} onClose={() => setActivePhoto(null)} />
      )}
    </section>
  )
}
