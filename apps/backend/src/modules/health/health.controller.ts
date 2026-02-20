import type { Request, Response } from "express";
import type { ApiHealth } from "@acee/shared";
import { pool } from "../../db/pool";

export async function healthHttp(_req: Request, res: Response) {
  const payload: ApiHealth = { ok: true };
  res.json(payload);
}

export async function healthDb(_req: Request, res: Response) {
  try {
    const r = await pool.query("SELECT 1 AS ok");
    res.json({ ok: true, db: r.rows[0] });
  } catch (e: any) {
    res.status(500).json({ ok: false, error: e?.message ?? String(e) });
  }
}
