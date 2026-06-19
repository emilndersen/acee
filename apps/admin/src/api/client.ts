const API_BASE = "/api";

let accessToken = localStorage.getItem("access_token") || "";
let refreshToken = localStorage.getItem("refresh_token") || "";

export function setTokens(access: string, refresh: string) {
  accessToken = access;
  refreshToken = refresh;
  localStorage.setItem("access_token", access);
  localStorage.setItem("refresh_token", refresh);
}

export function clearTokens() {
  accessToken = "";
  refreshToken = "";
  localStorage.removeItem("access_token");
  localStorage.removeItem("refresh_token");
}

export function isLoggedIn() {
  return !!accessToken;
}

async function tryRefresh(): Promise<boolean> {
  if (!refreshToken) return false;
  try {
    const res = await fetch(`${API_BASE}/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });
    if (!res.ok) return false;
    const data = await res.json();
    setTokens(data.tokens.access_token, data.tokens.refresh_token);
    return true;
  } catch {
    return false;
  }
}

export async function api<T = unknown>(
  path: string,
  opts: RequestInit = {}
): Promise<T> {
  const headers: Record<string, string> = {
    ...(opts.headers as Record<string, string>),
  };

  if (accessToken) {
    headers["Authorization"] = `Bearer ${accessToken}`;
  }

  if (!(opts.body instanceof FormData)) {
    headers["Content-Type"] = "application/json";
  }

  let res = await fetch(`${API_BASE}${path}`, { ...opts, headers });

  if (res.status === 401 && (await tryRefresh())) {
    headers["Authorization"] = `Bearer ${accessToken}`;
    res = await fetch(`${API_BASE}${path}`, { ...opts, headers });
  }

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `HTTP ${res.status}`);
  }

  return res.json();
}

export async function login(email: string, password: string) {
  const data = await api<{
    tokens: { access_token: string; refresh_token: string };
  }>("/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
  setTokens(data.tokens.access_token, data.tokens.refresh_token);
  return data;
}

export async function uploadFile(file: File): Promise<{ url: string; thumb_url: string; media_type: string }> {
  const form = new FormData();
  form.append("file", file);
  const data = await api<{ file: { image_url: string; thumb_url: string; media_type: string } }>("/upload", { method: "POST", body: form });
  return { url: data.file.image_url, thumb_url: data.file.thumb_url, media_type: data.file.media_type };
}
