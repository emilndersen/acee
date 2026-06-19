import { useEffect, useState } from 'react'
import { useScrollReveal } from '../../hooks/useScrollReveal'
import './Reviews.css'

export default function Reviews() {
  const ref = useScrollReveal()
  const [reviews, setReviews] = useState([])
  const [form, setForm] = useState({ author_name: '', text: '', rating: 5 })
  const [sent, setSent] = useState(false)

  useEffect(() => {
    fetch('/api/reviews')
      .then((r) => r.json())
      .then((d) => setReviews(d.reviews || []))
      .catch(() => {})
  }, [])

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      const res = await fetch('/api/reviews', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(form),
      })
      if (!res.ok) throw new Error()
      setSent(true)
      setForm({ author_name: '', text: '', rating: 5 })
      const d = await fetch('/api/reviews').then((r) => r.json())
      setReviews(d.reviews || [])
    } catch {
      alert('Не удалось отправить отзыв')
    }
  }

  return (
    <section className="reviews" ref={ref}>
      <h2 className="reviews__heading">Отзывы</h2>

      <div className="reviews__grid">
        {reviews.map((r) => (
          <div key={r.id} className="reviews__card">
            <div className="reviews__stars">{'⭐'.repeat(r.rating)}</div>
            <p className="reviews__text">{r.text}</p>
            <span className="reviews__author">— {r.author_name}</span>
          </div>
        ))}
      </div>

      {!sent ? (
        <form className="reviews__form" onSubmit={handleSubmit}>
          <h3 className="reviews__form-title">Оставить отзыв</h3>
          <input
            className="reviews__input"
            placeholder="Ваше имя"
            value={form.author_name}
            onChange={(e) => setForm((f) => ({ ...f, author_name: e.target.value }))}
            required
          />
          <div className="reviews__rating-select">
            {[1, 2, 3, 4, 5].map((n) => (
              <button
                key={n}
                type="button"
                className={`reviews__star-btn ${form.rating >= n ? 'active' : ''}`}
                onClick={() => setForm((f) => ({ ...f, rating: n }))}
              >
                ★
              </button>
            ))}
          </div>
          <textarea
            className="reviews__textarea"
            placeholder="Ваш отзыв..."
            value={form.text}
            onChange={(e) => setForm((f) => ({ ...f, text: e.target.value }))}
            required
          />
          <button type="submit" className="reviews__submit">Отправить</button>
        </form>
      ) : (
        <p className="reviews__thanks">Спасибо за отзыв! ✦</p>
      )}
    </section>
  )
}
