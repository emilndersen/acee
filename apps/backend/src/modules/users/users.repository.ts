import { pool } from "../../db/pool";
import type { CreateUserInput, User } from "@acee/shared";

export async function insertUser(input: CreateUserInput): Promise<User> {
  const q = `
    INSERT INTO users (name, email)
    VALUES ($1, $2)
    RETURNING id, name, email, created_at
  `;
  const r = await pool.query<User>(q, [input.name, input.email]);
  return r.rows[0];
}

export async function listUsers(): Promise<User[]> {
  const r = await pool.query<User>(
    `SELECT id, name, email, created_at FROM users ORDER BY id DESC`
  );
  return r.rows;
}
