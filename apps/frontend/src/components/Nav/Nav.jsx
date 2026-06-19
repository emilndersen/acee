import { Link, useLocation } from "react-router-dom";
import { useLang } from "../../i18n/LangContext";
import "./Nav.css";

export default function Nav() {
  const { pathname } = useLocation();
  const { lang, toggleLang, t } = useLang();
  const isHome = pathname === "/";

  function scrollTo(id) {
    const el = document.getElementById(id);
    if (el) el.scrollIntoView({ behavior: "smooth" });
  }

  return (
    <nav className="nav">
      {isHome ? (
        <>
          <a className="nav-link" href="#portfolio" onClick={(e) => { e.preventDefault(); scrollTo("portfolio"); }}>
            {t.nav.portfolio}
          </a>
        </>
      ) : (
        <Link className="nav-link" to="/">{t.nav.home}</Link>
      )}
      <button className="nav-lang" onClick={toggleLang}>
        {lang === 'ru' ? 'EN' : 'RU'}
      </button>
    </nav>
  );
}