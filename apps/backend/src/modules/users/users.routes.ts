import { Router } from "express";
import { usersCreate, usersList } from "./users.controller";

export const usersRouter = Router();

usersRouter.get("/", usersList);
usersRouter.post("/", usersCreate);
