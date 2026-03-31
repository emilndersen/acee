import { useEffect, useState } from 'react'
import './Hero.css'

export default function Hero() {
  const [progress, setProgress] = useState(0)

  useEffect(() => {
    const onScroll = () => {
      const max = window.innerHeight * 0.75
      const value = Math.min(window.scrollY / max, 1)
      setProgress(value)
    }

    onScroll()
    window.addEventListener('scroll', onScroll, { passive: true })
    return () => window.removeEventListener('scroll', onScroll)
  }, [])

  return (
    <section className="hero" id="home">
      <div className="hero-bg" />

      <div className={`hero-brand-fixed ${progress > 0.08 ? 'visible' : ''}`}>
        Ace Nikelsky
      </div>

      <div className="hero-wrapper">
        <div className="spike-container">
          <svg className="spike-ring" viewBox="0 0 420 420" xmlns="http://www.w3.org/2000/svg">
            <g transform="translate(210,210)">
              {[...Array(16)].map((_, i) => (
                <polygon
                  key={i}
                  points="0,-205 5,-135 -5,-135"
                  fill={i % 2 === 0 ? 'rgba(200,0,26,0.68)' : 'rgba(200,0,26,0.42)'}
                  transform={`rotate(${i * 22.5})`}
                />
              ))}
              <circle cx="0" cy="0" r="132" fill="none" stroke="rgba(200,0,26,0.12)" strokeWidth="1" />
            </g>
          </svg>

          <svg className="spike-ring spike-ring--inner" viewBox="0 0 420 420" xmlns="http://www.w3.org/2000/svg">
            <g transform="translate(210,210)">
              {[...Array(12)].map((_, i) => (
                <polygon
                  key={i}
                  points="0,-170 3,-115 -3,-115"
                  fill={i % 2 === 0 ? 'rgba(255,0,34,0.22)' : 'rgba(255,0,34,0.12)'}
                  transform={`rotate(${i * 30})`}
                />
              ))}
              <circle cx="0" cy="0" r="110" fill="none" stroke="rgba(200,0,26,0.08)" strokeWidth="1" />
            </g>
          </svg>

          <div className="photo-circle">
            <img
              src="https://images.unsplash.com/photo-1531746020798-e6953c6e8e04?w=400&q=80"
              alt="Ace Nikelsky"
            />
          </div>

          <div
            className="hero-text"
            style={{
              transform: `translate(-50%, ${-progress * 180}px)`,
              opacity: 1 - progress * 0.95,
            }}
          >
            <h1 className="hero-name">ACE</h1>
            <p className="hero-desc">Художественная фотография · Сюрреализм · Свет</p>
          </div>
        </div>
      </div>

      <div className="hero-scroll-indicator">scroll</div>
    </section>
  )
}