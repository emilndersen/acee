import { useScrollReveal } from '../../hooks/useScrollReveal'
import './Contacts.css'

const SOCIALS = [
  { label: 'ВКонтакте', href: '#' },
  { label: 'Telegram', href: '#' },
  { label: 'Instagram', href: '#' },
]

export default function Contacts() {
  const ref = useScrollReveal(0.1)

  return (
    <section className="contacts" id="contacts">
      <div className="contacts-left reveal" ref={ref}>
        <h2 className="contacts-title">КОНТАКТЫ</h2>

        <div className="contact-item">
          <div className="contact-icon">✆</div>
          <div className="contact-value">+7 (900) 000-00-00</div>
        </div>

        <div className="contact-item">
          <div className="contact-icon">@</div>
          <div className="contact-value">@acee_photo</div>
        </div>

        <div className="contact-social">
          {SOCIALS.map((s) => (
            <a key={s.label} className="social-link" href={s.href}>
              {s.label}
            </a>
          ))}
        </div>
      </div>

      <div className="footer-logo">ACEE</div>
    </section>
  )
}
