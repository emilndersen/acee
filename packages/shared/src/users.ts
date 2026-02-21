export type User = {
  id: number;
  name: string;
  email: string;
  created_at: string; // timestamptz as ISO string
};

export type CreateUserInput = {
  name: string;
  email: string;
};
