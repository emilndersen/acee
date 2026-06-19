import { useEffect } from 'react'
import './Lightbox.css'

export default function Lightbox({ media, src, onClose }) {
  useEffect(() => {
    const handler = (e) => { if (e.key === 'Escape') onClose() }
    document.addEventListener('keydown', handler)
    return () => document.removeEventListener('keydown', handler)
  }, [onClose])

  const isVideo = media?.media_type === 'video'
  const mediaSrc = media?.image_url || src

  if (!mediaSrc) return null

  return (
    <div className="lightbox" onClick={(e) => e.target === e.currentTarget && onClose()}>
      <div className="lightbox-inner">
        <button className="lightbox-close" onClick={onClose}>✕</button>
        {isVideo ? (
          <video
            className="lightbox-video"
            src={mediaSrc}
            controls
            autoPlay
          />
        ) : (
          <img className="lightbox-img" src={mediaSrc} alt="" />
        )}
      </div>
    </div>
  )
}
