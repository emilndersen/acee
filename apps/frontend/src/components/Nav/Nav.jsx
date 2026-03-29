import './Nav.css'

const links = [
  { label: 'Главная', href: '#home' },
  { label: 'Портфолио', href: '#portfolio' },
  { label: 'Инфо', href: '#info' },
  { label: 'Контакты', href: '#contacts' },
]

export default function Nav() {
  return (
    <nav className="nav">
      {links.map((link) => (
        <a key={link.href} className="nav-link" href={link.href}>
          {link.label}
        </a>
      ))}
    </nav>
  )
}
