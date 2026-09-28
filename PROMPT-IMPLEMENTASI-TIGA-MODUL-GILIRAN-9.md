# PROMPT — GILIRAN 9 *(sesi tunggal di `OUTPUT_HASIL_RNM`, cabang `main` @ `0e469cd`)*: **PL 05a → 05b → 06 + 08 + 09 → Komite 01 → 09 → Claim Life §3.1** — laporan pertama paling cepat sesudah **05b**

> GILIRAN-3 *(A/B/C)*, 4–8 **tetap berlaku**, termasuk aturan berhenti GILIRAN-6. Mekanisme giliran = lanjutan 8 §1.
> Paket 0 GILIRAN-8 *(bh, bi)* dan OQ-PL-08 **sudah tuntas**; brief ini menetapkan keadaan awal dan urutan sisa.

## 0. KEADAAN AWAL — DIVERIFIKASI ASISTEN 28-09-2026

| Klaim laporan | Diperiksa ulang | Hasil |
| --- | --- | --- |
| 3 commit `8b6d001`, `e8487b0`, `0e469cd`; berhenti "konteks menipis 8%" | 3 commit; pohon bersih; alasan sah | ✅ |
| Go 504 · 0 · 38; JS 339 | dijalankan ulang asisten *(lihat pesan verifikasi)* | ✅ |
| **bh** — pembaca view ditulis utuh, penjaga dipersempit lima batas | `repository/ambangproduk.go` *(`NamaViewProdukLife = "PRODUCTINWARD_LIFE"`)*; nama tidak dirakit dari potongan | ✅ — pembatalan tipu daya penjaga dicatat dan diterima |
| **bi** — Claim Life ke folder modulnya | `pages/claimlife/` *(4 halaman)*, `components/claimlife/` *(6 komponen)*, `assets/labels.claimlife.ts`; perilaku tidak berubah | ✅ |
| ⚠️ **ralat atas GILIRAN-8 §1**: `.PREMIUM` bukan Σ empat kolom, melainkan **empat `WHEN` saling meniadakan** atas `pyWorkPage.Type` | `AppendCurrencySummary_DT.xml` b2268 `Type=="QR"` → `+ .GROSS_PREMIUM` b2298; b2328 `"QP"` → `+ .GROSS_PREMIUM_REFUND` b2359; b2388 `"TP"` → `+ .GROSS_PREMIUM_RETRO` b2419; b2448 `"TR"` → `+ .GROSS_PREMIUM_REFUND_RETRO` b2478 — **satu polis, satu cabang** | ✅ **ralat diterima**; bacaan asisten yang meratakan `pyPropertiesName` keliru; `.BALANCE` dan `.COMMISSION` mengikuti pola yang sama *(cabang per `Type`, syarat `Param.Currency == .CURRENCY`)* |
| tiga keanehan warisan `.BALANCE` dipertahankan | `+BROKERAGE_FEE_RETRO` pada TP/TR lawan `−BROKERAGE_FEE` pada QR; cabang QP memakai `TAX`/`PROF_COMM`/`CLAIM` bukan `*_REFUND`; `@divide(…,1,4)` sebagai pembulatan | ✅ VERBATIM, dicatat di tiket 05a |

## 1. URUTAN — semua di `main`, satu commit per tiket

| # | Paket | Isi | Rujukan |
| ---: | --- | --- | --- |
| 1 | **PL 05a** summary uang | 32 kolom uang per `CURRENCY`; `.PREMIUM`/`.COMMISSION`/`.BALANCE` per cabang `Type` **seperti DT, satu cabang per polis**; **pl2** `M_LIFE_PREMIUM_SUMMARY`/`M_LIFE_PREMIUM_DETAIL` satu transaksi *(decision o: prosedur tidak dipanggil)*; `SubmitPremiumList_Act` langkah 15 diputuskan dengan bukti; contoh terhitung dari literal untuk keempat `Type`; uji menolak polis ber-`Type` di luar QR/QP/TP/TR | tiket 05a; brief 3-PREMIUMLIST §2 |
| 2 | **PL 05b** JSON polis | `InsertJsonPolisLife_Act` + `finishAssignment` b26442; keputusan tiket/spec diikuti, ralat berbukti bila bertentangan | tiket 05b |
| — | **laporan pertama boleh di sini** | | |
| 3 | **PL 06 + 08 + 09** | outbox Arasapas/email + stub; kontrak hilir dua sisi; verifikasi 050–056 lawan STRUKTUR | brief 3-PREMIUMLIST §2 |
| 4 | **Komite 01 → 09** | Inbox Komite pola grid; `ShowTransfer` layar keputusan *(GILIRAN-3-KOMITE §1)*; `KomitePostAdjustment` utuh; km1–km5 | brief 3-KOMITE |
| 5 | **Claim Life §3.1** | 41 baris sensus diputuskan; celah dibangun | brief 3-CLAIM-LIFE §3.1 |

## 2. ATURAN · LAPORAN · TELEMETRI

Sama dengan GILIRAN-4 §2–§3 dan GILIRAN-6 aturan berhenti. Pelajaran giliran 8 dijadikan aturan: **struktur `WHEN` di
DataTransform dibaca dari urutan aksi**, bukan dari pasangan `pyPropertiesName`/`pyPropertiesValue` yang diratakan; rumus uang
yang dibaca brief **diverifikasi ulang** executor sebelum dibangun. Laporan: baris pertama alasan berhenti; tabel **tiket → commit →
rule XML → rute/komponen**; ralat tiket; OQ; angka uji tiap commit; bab **TELEMETRI EKSEKUSI**.

---

*Disusun 28 September 2026 sesudah verifikasi `0e469cd` (uji dijalankan ulang; `AppendCurrencySummary_DT.xml` b2268–b2478
dibaca ulang sebagai urutan aksi; folder Claim Life dan `ambangproduk.go` dicek).*
