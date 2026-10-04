/// <reference types="vite/client" />

// Memberi tahu TypeScript env var apa saja yang dibaca frontend, supaya
// `import.meta.env.VITE_API_BASE_URL` dikenal sebagai teks (atau kosong).
// Hanya env var berawalan VITE_ yang diteruskan Vite ke browser.
interface ImportMetaEnv {
  readonly VITE_API_BASE_URL?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
