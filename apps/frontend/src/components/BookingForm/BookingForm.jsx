import { useEffect, useState } from 'react'
import { useLang } from '../../i18n/LangContext'
import './BookingForm.css'

export default function BookingForm() {
  const { t } = useLang()
  const [open, setOpen] = useState(false)
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

  useEffect(() => {
    document.body.style.overflow = open ? 'hidden' : ''
    return () => { document.body.style.overflow = '' }
  }, [open])

  useEffect(() => {
    function onKey(e) { if (e.key === 'Escape') setOpen(false) }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
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
          idea: form.idea,
        }),
      })
      if (!res.ok) throw new Error()
      setSent(true)
      setForm({ name: '', contact: '', telegram: '', type: '', shoot_date: '', idea: '' })
    } catch {
      alert(t.booking.error)
    }
  }

  const isBusy = busyDates.includes(form.shoot_date)

  return (
    <>
      <button className="booking-trigger" onClick={() => setOpen(true)}>
        <span>{t.booking.trigger}</span>
      </button>

      {open && (
        <div className="booking-overlay" onClick={(e) => { if (e.target === e.currentTarget) setOpen(false) }}>
          <div className="booking-modal">
            <button className="booking-close" onClick={() => setOpen(false)}>✕</button>
            <h2 className="booking-modal__title">{t.booking.title}</h2>
            <p className="booking-modal__subtitle">{t.booking.subtitle}</p>

            {sent ? (
              <div className="booking-success"><span>{t.booking.success}</span></div>
            ) : (
              <form className="booking-form" onSubmit={handleSubmit}>
                <div className="form-group">
                  <label className="form-label">{t.booking.name}</label>
                  <input className="form-input" name="name" value={form.name} onChange={handleChange} placeholder={t.booking.namePlaceholder} required />
                </div>
                <div className="form-group">
                  <label className="form-label">{t.booking.phone}</label>
                  <input className="form-input" name="contact" value={form.contact} onChange={handleChange} placeholder={t.booking.phonePlaceholder} required />
                </div>
                <div className="form-group">
                  <label className="form-label">{t.booking.telegram}</label>
                  <input className="form-input" name="telegram" value={form.telegram} onChange={handleChange} placeholder={t.booking.telegramPlaceholder} />
                </div>
                <div className="form-group">
                  <label className="form-label">{t.booking.shootType}</label>
                  <select className="form-select" name="type" value={form.type} onChange={handleChange} required>
                    <option value="">{t.booking.shootTypePlaceholder}</option>
                    {t.booking.types.map((tp) => (
                      <option key={tp} value={tp}>{tp}</option>
                    ))}
                  </select>
                </div>
                <div className="form-group">
                  <label className="form-label">{t.booking.date}</label>
                  <input className="form-input" type="date" name="shoot_date" value={form.shoot_date} onChange={handleChange} />
                  {isBusy && <span className="form-busy-warning">{t.booking.busy}</span>}
                </div>
                <div className="form-group form-group--full">
                  <label className="form-label">{t.booking.idea}</label>
                  <textarea className="form-textarea" name="idea" value={form.idea} onChange={handleChange} placeholder={t.booking.ideaPlaceholder} />
                </div>
                <div className="form-submit">
                  <button className="btn-book" type="submit"><span>{t.booking.submit}</span></button>
                </div>
              </form>
            )}
          </div>
        </div>
      )}
    </>
  )
}