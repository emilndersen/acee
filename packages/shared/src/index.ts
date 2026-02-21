export type ApiHealth = { ok: boolean };

export type ApiError = {
  ok: false;
  error: string;
};

export * from "./users";
