import { Link, useLocation } from "react-router-dom";
import "./Nav.css";

export default function Nav() {
  const { pathname } = useLocation();
  const isHome = pathname === "/";

  function scrollTo(id) {
    const el = document.getElementById(id);
    if (el) el.scrollIntoView({ behavior: "smooth" });
  }

  return (
    <nav className="nav">
      {isHome ? (
        <>
          <a className="nav-link" href="#portfolio" onClick={(e) => { e.preventDefault(); scrollTo("portfolio"); }}>Portfolio</a>
          <a className="nav-link" href="#booking" onClick={(e) => { e.preventDefault(); scrollTo("booking"); }}>Booking</a>
        </>
      ) : (
        <Link className="nav-link" to="/">Home</Link>
      )}
    </nav>
  );
}
