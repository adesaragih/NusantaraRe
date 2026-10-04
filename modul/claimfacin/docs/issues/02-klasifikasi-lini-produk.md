# 02: Klasifikasi lini produk — 13 penggolong hidup, 36 tidak dialihkan

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 01 (registrasi)
**Menutup:** AC 13 · 14 · 15 · 16 · 17 · 18 *(6 AC)* — US 34–36

## Hasil & nilai pengguna

Hari ini Penggolongan lini usaha **membaca medan data kutipan yang tidak pernah disalin** ke objek kerja. ⛔ `[terverifikasi]` Akibatnya cabang **MBU** dan **Travel** tidak pernah terbit — diam-diam, tanpa galat.

Sesudah tiket ini, Setiap lini usaha **tergolong benar**, termasuk MBU dan Travel yang dulu tak pernah terbit. ⭐ Penggolongan membaca **data kutipan yang lengkap**.

## Area codebase

- Lapisan layanan: penggolongan lini usaha
- Penyalinan data kutipan dari kasus induk

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Penggolong lini usaha | 49 rule penggolong; ⭐ **13 dipakai hidup**, ⛔ **36 tidak dipakai sama sekali** |
| Sumber data kutipan | salinan data kutipan ke objek kerja |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

## ADR terkait

- **ADR-0003** — dua kaki dasar klasifikasi

## Acceptance criteria

- [ ] **AC 13–18** — klasifikasi lini produk
- [ ] ⭐ **13 penggolong** dialihkan; ⛔ **36 tidak dibangun**
- [ ] ⭐ Cabang **MBU** dan **Travel** **terbit** bila datanya memenuhi — ⚠️ inilah cacat lama yang ditutup

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | Isi rule peran | ⛔ tidak menyentuh tiket ini |

## Perintah verifikasi

1. Jalankan satu klaim tiap lini usaha yang hidup — ⭐ ketiga belasnya tergolong benar.
2. Jalankan klaim lini **MBU** dan **Travel** — ⭐ cabangnya **terbit**; ⛔ di sistem lama tidak.
3. Cari pemanggilan salah satu dari 36 penggolong yang tidak dialihkan — ⛔ **nihil**.
