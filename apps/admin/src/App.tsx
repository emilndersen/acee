import { Routes, Route, Navigate } from "react-router-dom";
import { isLoggedIn } from "./api/client";
import Login from "./pages/Login";
import Dashboard from "./pages/Dashboard";
import Albums from "./pages/Albums";
import AlbumPhotos from "./pages/AlbumPhotos";
import Bookings from "./pages/Bookings";
import Settings from "./pages/Settings";
import Layout from "./components/Layout";

function PrivateRoute({ children }: { children: React.ReactNode }) {
  return isLoggedIn() ? <>{children}</> : <Navigate to="/login" />;
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route
        path="/"
        element={
          <PrivateRoute>
            <Layout />
          </PrivateRoute>
        }
      >
        <Route index element={<Dashboard />} />
        <Route path="albums" element={<Albums />} />
        <Route path="albums/:slug" element={<AlbumPhotos />} />
        <Route path="bookings" element={<Bookings />} />
        <Route path="settings" element={<Settings />} />
      </Route>
    </Routes>
  );
}