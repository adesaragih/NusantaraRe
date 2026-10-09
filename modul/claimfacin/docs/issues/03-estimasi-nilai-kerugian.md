# 03: Estimasi nilai kerugian dan validasinya

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR) · 01 (registrasi)
**Menutup:** AC 19 · 20 · 21 · 22 · 23 · 24 · 25 · 26 · 27 · 28 · 29 · 30 *(12 AC)* — US 7 · 8

## Hasil & nilai pengguna

Hari ini Estimasi kerugian **belum punya tempat menggantung** pada lini FAC, dan aturan validasinya belum ditetapkan.

Sesudah tiket ini, Penilai dapat **mencatat estimasi per item objek**, dengan mata uang, risiko sendiri, dan konversinya — dan sistem **menolak** estimasi yang melanggar aturan.

## Area codebase

- Lapisan layanan klaim: estimasi
- Validasi nilai estimasi terhadap nilai pertanggungan

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Estimasi | daftar estimasi di dalam item objek |
| Validasi | rule validasi masukan estimasi |

⭐ Rincian medan dan asalnya ada di `STRUKTUR-TABEL-CLAIM-FACIN.md` **§2b**.

## ADR terkait

- **ADR-0007** — jejak audit atas perubahan nilai

## Acceptance criteria

- [ ] **AC 19–30** — estimasi dan validasinya
- [ ] ⭐ Estimasi menggantung pada **item objek**, bukan pada klaim

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **12** | ⚠️ **10 dari 13 kolom estimasi tanpa penulis di korpus** — diduga diisi lewat layar, ⛔ belum terbukti | ⚠️ **menahan jalur TULIS**, tidak menahan jalur BACA |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"⚠️ **10 dari 13 kolom estimasi tanpa penulis di
> korpus**"* → ⚠️ **ruang nomor**: butir **12** di sini adalah butir register **`STRUKTUR-TABEL-CLAIM-FACIN.md` §6**,
> bukan butir 12 spec (penggolongan tujuh pesan, sudah ditutup). Jalur TULIS estimasi **sudah dibangun** 10-10-2026:
> kolom = katalog `backend/models/katalog_tabel.go` `TabelEstimasi` (nama DDL Claim Prop `521` + `ADD` `563`, mis.
> `GROSS_ESTIMATION_VALUE`, `ESTIMATION_VALUE_IDR`, `CURRENCY_NAME`), baris FAC mengisi `CLAIM_ID` **dan**
> `OBJECT_ITEM_ID`; uang `NUMBER(38,10)`. Tidak menahan lagi.

## Perintah verifikasi

1. Catat estimasi pada satu item objek — ⭐ tersimpan dan terbaca kembali.
2. Catat estimasi melampaui nilai pertanggungan — ⭐ perilakunya sesuai AC.
3. Hapus item objek — ⭐ estimasinya **ikut terhapus**.
