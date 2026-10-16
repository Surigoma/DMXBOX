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
      routeFileIgnorePattern: "\\.test\\.[jt]sx?$",
    }),
    react(),
  ],
  test: {
    testTimeout: 3000,
    fileParallelism: true,
    sequence: {
      concurrent: false,
    },
    isolate: true,
    browser: {
      enabled: true,
      provider: playwright(),
      headless: false,
      instances: [{ browser: "chromium" }],
      // Vitest 5 locators are exact/case-sensitive by default; keep the v4 behavior.
      locators: { exact: false },
    },
    coverage: {
      reportsDirectory: "./test/coverage/",
      provider: "v8",
      include: ["src/**/*.{ts,tsx}"],
      exclude: ["src/**/*.test.{ts,tsx}"],
    },
    reporters: ["default", "html"],
    outputFile: {
      html: "./test/unit/index.html",
    },
  },
  build: {
    rolldownOptions: {
      output: {
        minify: true,
      },
    },
  },
});
