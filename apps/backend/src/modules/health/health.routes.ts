import { Router } from "express";
import { healthDb, healthHttp } from "./health.controller";

export const healthRouter = Router();

healthRouter.get("/", healthHttp);
healthRouter.get("/db", healthDb);
