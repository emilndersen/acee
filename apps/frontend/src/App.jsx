import './styles/globals.css'

import Cursor      from './components/Cursor/Cursor'
import Nav         from './components/Nav/Nav'
import Hero        from './components/Hero/Hero'
import Portfolio   from './components/Portfolio/Portfolio'
import BookingForm from './components/BookingForm/BookingForm'
import Contacts    from './components/Contacts/Contacts'

export default function App() {
  return (
    <>
      {/* Utilities */}
      <Cursor />
      <Nav />

      {/* Sections */}
      <Hero />

      <div className="divider" />

      <Portfolio />

      <div className="divider" />

      <BookingForm />

      <div className="divider" />

      <Contacts />
    </>
  )
}
