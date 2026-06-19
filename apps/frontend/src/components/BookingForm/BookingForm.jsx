import { useEffect, useState } from 'react'
import { useScrollReveal } from '../../hooks/useScrollReveal'
import './BookingForm.css'

const SHOOT_TYPES = ['Портрет', 'Арт-съёмка', 'Фэшн', 'Репортаж', 'Другое']

export default function BookingForm() {
  const titleRef = useScrollReveal(0.1)
  const formRef = useScrollReveal(0.1, 100)
  const [busyDates, setBusyDates] = useState([])

  const [form, setForm] = useState({
    name: '',
    contact: '',
    telegram: '',
    type: '',
    shoot_date: '',
    idea: '',
  })
  const [sent, setSent] = useState(false)

  useEffect(() => {
    fetch('/api/bookings/busy-dates')
      .then((r) => r.json())
      .then((d) => setBusyDates(d.dates || []))
      .catch(() => {})
  }, [])

  const handleChange = (e) =>
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }))

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      const res = await fetch('/api/bookings', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: form.name,
          contact: form.contact,
          telegram: form.telegram,
          shoot_type: form.type,
          date: form.shoot_date,
          shoot_date: form.shoot_date || null,
          idea: form.idea,
        }),
      })
      if (!res.ok) throw new Error('Ошибка отправки')
      setSent(true)
      setForm({ name: '', contact: '', telegram: '', type: '', shoot_date: '', idea: '' })
    } catch (error) {
      console.error(error)
      alert('Не удалось отправить заявку')
    }
  }

  const isBusy = busyDates.includes(form.shoot_date)

  return (
    <section className="booking" id="booking">
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
            <input className="form-input" name="name" value={form.name} onChange={handleChange} placeholder="Твоё имя" required />
          </div>

          <div className="form-group">
            <label className="form-label">Телефон</label>
            <input className="form-input" name="contact" value={form.contact} onChange={handleChange} placeholder="+7 999 123 45 67" required />
          </div>

          <div className="form-group">
            <label className="form-label">Telegram</label>
            <input className="form-input" name="telegram" value={form.telegram} onChange={handleChange} placeholder="@username" />
          </div>

          <div className="form-group">
            <label className="form-label">Тип съёмки</label>
            <select className="form-select" name="type" value={form.type} onChange={handleChange} required>
              <option value="">Выбери формат</option>
              {SHOOT_TYPES.map((t) => (
                <option key={t} value={t}>{t}</option>
              ))}
            </select>
          </div>

          <div className="form-group">
            <label className="form-label">Дата съёмки</label>
            <input className="form-input" type="date" name="shoot_date" value={form.shoot_date} onChange={handleChange} />
            {isBusy && <span className="form-busy-warning">Этот день уже занят</span>}
          </div>

          <div className="form-group form-group--full">
            <label className="form-label">Твоя идея</label>
            <textarea className="form-textarea" name="idea" value={form.idea} onChange={handleChange} placeholder="Расскажи о концепции, настроении, референсах..." />
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
