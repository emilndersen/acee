import type { Request, Response } from "express";
import { createUser, getUsers } from "./users.service";
import type { CreateUserInput } from "@acee/shared";

export async function usersList(_req: Request, res: Response) {
  const users = await getUsers();
  res.json({ ok: true, users });
}

export async function usersCreate(req: Request, res: Response) {
  const body = req.body as Partial<CreateUserInput>;

  if (!body?.name || !body?.email) {
    return res.status(400).json({ ok: false, error: "name and email are required" });
  }

  const user = await createUser({ name: body.name, email: body.email });
  res.status(201).json({ ok: true, user });
}
