import './Hero.css'

export default function Hero() {
  return (
    <section className="hero" id="home">
      <div className="hero-bg" />

      <div className="hero-wrapper">
        <div className="spike-container">

          {/* Outer spike ring — rotates clockwise */}
          <svg className="spike-ring" viewBox="0 0 420 420" xmlns="http://www.w3.org/2000/svg">
            <g transform="translate(210,210)">
              {[...Array(16)].map((_, i) => (
                <polygon
                  key={i}
                  points="0,-205 5,-135 -5,-135"
                  fill={i % 2 === 0 ? 'rgba(200,0,26,0.7)' : 'rgba(200,0,26,0.5)'}
                  transform={`rotate(${i * 22.5})`}
                />
              ))}
              {[...Array(8)].map((_, i) => (
                <polygon
                  key={'thin-' + i}
                  points="0,-195 2,-150 -2,-150"
                  fill="rgba(255,0,34,0.25)"
                  transform={`rotate(${i * 45 + 11.25})`}
                />
              ))}
              <circle cx="0" cy="0" r="132" fill="none" stroke="rgba(200,0,26,0.15)" strokeWidth="1" />
              <circle cx="0" cy="0" r="138" fill="none" stroke="rgba(200,0,26,0.08)" strokeWidth="1" />
            </g>
          </svg>

          {/* Inner spike ring — rotates counter-clockwise */}
          <svg className="spike-ring spike-ring--inner" viewBox="0 0 420 420" xmlns="http://www.w3.org/2000/svg">
            <g transform="translate(210,210)">
              {[...Array(12)].map((_, i) => (
                <polygon
                  key={i}
                  points="0,-170 3,-115 -3,-115"
                  fill={i % 2 === 0 ? 'rgba(255,0,34,0.3)' : 'rgba(255,0,34,0.2)'}
                  transform={`rotate(${i * 30})`}
                />
              ))}
              <circle cx="0" cy="0" r="110" fill="none" stroke="rgba(200,0,26,0.1)" strokeWidth="1" />
            </g>
          </svg>

          {/* Photographer photo */}
          <div className="photo-circle">
            <img
              src="https://images.unsplash.com/photo-1531746020798-e6953c6e8e04?w=400&q=80"
              alt="Photographer"
            />
          </div>

          {/* Name & tagline */}
          <div className="hero-text">
            <div className="hero-name">ACEE</div>
            <div className="hero-desc">Художественная фотография · Сюрреализм · Свет</div>
          </div>
        </div>
      </div>
    </section>
  )
}
