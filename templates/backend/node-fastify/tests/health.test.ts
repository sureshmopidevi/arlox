import { afterEach, describe, expect, it } from "vitest";

import { buildApp } from "../src/app.js";

const app = buildApp();

afterEach(async () => app.close());

describe("GET /health", () => {
  it("reports service health", async () => {
    const response = await app.inject({ method: "GET", url: "/health" });

    expect(response.statusCode).toBe(200);
    expect(response.json()).toEqual({ status: "ok" });
  });
});
