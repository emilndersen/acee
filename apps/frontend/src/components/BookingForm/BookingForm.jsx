import { useState } from 'react'
import { useScrollReveal } from '../../hooks/useScrollReveal'
import './BookingForm.css'

const SHOOT_TYPES = ['Портрет', 'Арт-съёмка', 'Фэшн', 'Репортаж', 'Другое']

export default function BookingForm() {
  const titleRef = useScrollReveal(0.1)
  const formRef = useScrollReveal(0.1, 100)

  const [form, setForm] = useState({
    name: '',
    contact: '',
    type: '',
    date: '',
    idea: '',
  })
  const [sent, setSent] = useState(false)

  const handleChange = (e) =>
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }))

  const handleSubmit = (e) => {
    e.preventDefault()
    // TODO: wire up to Go backend POST /api/bookings
    console.log('Booking:', form)
    setSent(true)
  }

  return (
    <section className="booking">
      <div className="booking-title reveal" ref={titleRef}>
        ЗАПИСАТЬСЯ<br />НА СЪЁМКУ
      </div>
      <p className="booking-subtitle reveal" ref={useScrollReveal(0.1, 80)}>
        Расскажи мне о своей идее — воплотим вместе
      </p>

      {sent ? (
        <div className="booking-success">
          <span>Заявка отправлена — скоро свяжусь ✦</span>
        </div>
      ) : (
        <form className="booking-form reveal" ref={formRef} onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="form-label">Имя</label>
            <input
              className="form-input"
              name="name"
              value={form.name}
              onChange={handleChange}
              placeholder="Твоё имя"
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label">Телефон / Telegram</label>
            <input
              className="form-input"
              name="contact"
              value={form.contact}
              onChange={handleChange}
              placeholder="+7 или @username"
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label">Тип съёмки</label>
            <select
              className="form-select"
              name="type"
              value={form.type}
              onChange={handleChange}
              required
            >
              <option value="">Выбери формат</option>
              {SHOOT_TYPES.map((t) => (
                <option key={t} value={t}>{t}</option>
              ))}
            </select>
          </div>

          <div className="form-group">
            <label className="form-label">Дата (примерно)</label>
            <input
              className="form-input"
              name="date"
              value={form.date}
              onChange={handleChange}
              placeholder="Месяц / число"
            />
          </div>

          <div className="form-group form-group--full">
            <label className="form-label">Твоя идея</label>
            <textarea
              className="form-textarea"
              name="idea"
              value={form.idea}
              onChange={handleChange}
              placeholder="Расскажи о концепции, настроении, референсах..."
            />
          </div>

          <div className="form-submit">
            <button className="btn-book" type="submit">
              <span>ОТПРАВИТЬ ЗАЯВКУ</span>
            </button>
          </div>
        </form>
      )}
    </section>
  )
}
