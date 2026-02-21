import type { CreateUserInput, User } from "@acee/shared";
import { insertUser, listUsers } from "./users.repository";

function normalizeEmail(email: string) {
  return email.trim().toLowerCase();
}

export async function createUser(input: CreateUserInput): Promise<User> {
  return insertUser({ ...input, email: normalizeEmail(input.email) });
}

export async function getUsers(): Promise<User[]> {
  return listUsers();
}
