import express from "express";

export const app = express();

app.use(express.json());
app.get("/health", (_request, response) => {
  response.json({ status: "ok" });
});
