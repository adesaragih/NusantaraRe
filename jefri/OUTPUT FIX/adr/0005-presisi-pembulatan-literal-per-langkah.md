---
status: accepted
---

# Presisi pembulatan ditulis literal di tiap langkah, bukan disentralkan

Setiap pembulatan hasil port menuliskan presisinya **sebagai literal di tempatnya**, disertai
komentar `// Asal: <rule>, langkah <n>, presisi <p>` (`CLAUDE.md` §4.6). Presisi **tidak**
disentralkan ke registry maupun dilekatkan pada nilai. Ini satu-satunya tempat dalam rancangan ini
yang memilih duplikasi ketimbang sentralisasi, dan alasannya spesifik.

## Mengapa

`[terverifikasi]` Di sistem lama, presisi **bukan properti nilai maupun properti rule** — ia properti
**satu langkah tertentu di dalam** sebuah rule:

- Aritmetika yang sama dibulatkan **4 desimal** di satu rule dan **20 desimal** di rule lain.
- Satu pembagi muncul dengan **9 presisi berbeda**.
- Sebuah pembulatan berada **di dalam** loop akumulasi, sehingga galatnya menumpuk per iterasi.

## Considered Options

- **Registry terpusat `nama rule → presisi`** — ditolak. Ia mengandaikan satu rule punya satu
  presisi, dan itu tidak benar.
- **Presisi melekat pada nilai dan mengalir lewat operasi** — ditolak. Presisi tidak mengalir; ia
  diterapkan di satu titik lalu selesai.

## Consequences

Uji rekonsiliasi **wajib** memuat kasus khusus untuk pembulatan **di dalam loop akumulasi**, bukan
hanya nilai tunggal. Pembulatan per-iterasi menumpuk galat, sehingga port yang benar untuk satu nilai
masih bisa meleset untuk daftar panjang — dan itu justru bentuk kesalahan yang paling sulit terlihat.

Menyeragamkan presisi adalah **perbaikan terpisah**, bukan bagian migrasi.
