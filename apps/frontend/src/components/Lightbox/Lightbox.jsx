import { useEffect } from 'react'
import './Lightbox.css'

export default function Lightbox({ photo, onClose }) {
  // Close on Escape
  useEffect(() => {
    const handler = (e) => { if (e.key === 'Escape') onClose() }
    document.addEventListener('keydown', handler)
    return () => document.removeEventListener('keydown', handler)
  }, [onClose])

  if (!photo) return null

  return (
    <div className="lightbox" onClick={(e) => e.target === e.currentTarget && onClose()}>
      <div className="lightbox-inner">
        <button className="lightbox-close" onClick={onClose}>[ ЗАКРЫТЬ ]</button>
        <img className="lightbox-img" src={photo.src} alt={photo.title} />
        <div className="lightbox-meta">
          <div className="lightbox-album">// {photo.album}</div>
          <div className="lightbox-title">{photo.title.toUpperCase()}</div>
          <div className="lightbox-desc">{photo.desc}</div>
          <a href="#portfolio" className="lightbox-go">→ Перейти к альбому</a>
        </div>
      </div>
    </div>
  )
}
