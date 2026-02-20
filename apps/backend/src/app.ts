import express from "express";
import cors from "cors";
import { healthRouter } from "./modules/health/health.routes";

export const app = express();

app.use(cors());
app.use(express.json());

app.use("/health", healthRouter);
