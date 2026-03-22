/// <reference types="vitest/config" />
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import { playwright } from "@vitest/browser-playwright";

// https://vite.dev/config/
export default defineConfig({
  base: "/gui",
  plugins: [
    tanstackRouter({
      target: "react",
      autoCodeSplitting: true,
    }),
    react(),
  ],
  test: {
    browser: {
      enabled: true,
      provider: playwright(),
      headless: true,
      instances: [{ browser: "chromium" }],
    },
    coverage: {
      reportsDirectory: "./test/coverage/",
      provider: "v8",
      include: ["./src/**/*.ts", "./src/**/*.tsx"],
      exclude: ["./src/**/*.test.ts", "./src/**/*.test.tsx"],
    },
    reporters: ["default", "html"],
    outputFile: {
      html: "./test/unit/index.html",
    },
  },
});
