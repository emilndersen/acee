import { Pool } from "pg";
import { env } from "../config/env";

export const pool = new Pool({
  host: env.DB_HOST=localhost,
  port: env.DB_PORT=5432,
  database: env.DB_NAME=acee_db,
  user: env.DB_USER=ace_user,
  password: env.DB_PASSWORD=123456789,
});