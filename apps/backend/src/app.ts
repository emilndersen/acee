import express from "express";
import cors from "cors";
import { healthRouter } from "./modules/health/health.routes";
import { usersRouter } from "./modules/users/users.routes";

export const app = express();

app.use(cors());
app.use(express.json());

app.use("/health", healthRouter);
app.use("/api/users", usersRouter);
