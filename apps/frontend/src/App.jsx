import { Routes, Route } from "react-router-dom";
import { LangProvider } from "./i18n/LangContext";
import Cursor from "./components/Cursor/Cursor";
import Nav from "./components/Nav/Nav";
import Home from "./pages/Home";
import AlbumPage from "./pages/AlbumPage";

export default function App() {
  return (
    <LangProvider>
      <Cursor />
      <Nav />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/albums/:slug" element={<AlbumPage />} />
      </Routes>
    </LangProvider>
  );
}