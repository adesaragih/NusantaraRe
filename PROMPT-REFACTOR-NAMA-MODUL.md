# PROMPT — PENYERAGAMAN NAMA FOLDER MODUL *(sesi baru, folder `OUTPUT_HASIL_RNM`, cabang `main` @ `c2fd123` atau lebih baru)*

> Keputusan work owner 30-09-2026: *"oke"* atas usul asisten — **satu tabel nama untuk semua tempat, tanpa singkatan**. Refactor bentuk B
> *(`3ff9301` → `5a281cf`)* memakai nama singkat `treaty`, `premiumlist`, `komite`; itu diseragamkan di sini. **Nol perubahan perilaku**; angka
> uji sama persis sebelum dan sesudah. Pindahkan dengan `git mv`; commit dengan jalur eksplisit. **Hanya sesi ini** yang bekerja di
> `OUTPUT_HASIL_RNM` selama penggantian nama.

## 0. VERIFIKASI BENTUK B *(asisten, 30-09-2026)*

`internal/` kosong; `inti/` *(config, db, galat, jejak, kontrak, layanan, migrasi, outbox, penjaga, penomor, uang, unggah, utils)*;
`modul/{claimlife,komite,premiumlist,treaty}` + `modul/daftar.go`; migrasi per modul: claimlife **22**, komite **1**, premiumlist **9**, treaty **0**;
frontend `src/inti/` + `src/modul/{claimlife,komite,premiumlist,treaty}` + `daftar.ts`. Uji: Go **957 PASS · 0 FAIL**, dengan tag `db`
**59 SKIP**; vitest **678**; tsc, vet bersih; build **98** modul. ✅

## 1. TABEL NAMA — satu-satunya sumber

| Modul korpus | Backend Go `APP_RNM/modul/…` *(tanpa tanda hubung)* | Frontend `frontend/src/modul/…` dan `.scratch/…` | Nilai `MODUL_AKTIF` |
| --- | --- | --- | --- |
| Claim Life | `claimlife` *(tetap)* | `claim-life` | `claimlife` |
| PremiumList Life | `premiumlist` → **`premiumlistlife`** | `premiumlist` → **`premiumlist-life`** | `premiumlistlife` |
| Komite Claim Life | `komite` → **`komiteclaimlife`** | `komite` → **`komite-claim-life`** | `komiteclaimlife` |
| Treaty Contract Out | `treaty` → **`treatycontractout`** | `treaty` → **`treaty-contract-out`** | `treatycontractout` |
| *(modul berikutnya, mis. NB FacIn)* | `nbfacin` | `nb-facin` | `nbfacin` |

Aturan: **nama backend = nama dokumen `.scratch` tanpa tanda hubung**; frontend memakai nama `.scratch` persis. Nama paket Go = nama folder.

## 2. YANG TIDAK BERUBAH — dijaga uji

- **Nama berkas migrasi** *(`T_MIGRASI` mencatat nama, bukan folder)*: `001_…`–`058_…` tetap.
- **URL rute API** *(`/api/klaim-life`, `/api/polis-life`, `/api/komite`, `/api/treaty-contract-out`, …)* dan teks layar.
- Berkas dan awalan di dalam modul *(`tco_*`, `polis_*`, `komite_*`)*.
- `.scratch/<modul>/` *(sudah memakai nama panjang)*.

## 3. URUTAN — satu commit per paket

| # | Paket | Commit |
| ---: | --- | --- |
| 1 | backend: `git mv` tiga folder `modul/*`; ubah `package` dan seluruh jalur impor; `modul/daftar.go`; penjaga impor lintas modul dan penjaga lain yang menyebut nama folder; `MODUL_AKTIF` menerima nama baru **dan** menolak nama lama dengan galat berkata-kata *(bukan diam)* | `refactor: nama folder modul backend sesuai tabel nama` |
| 2 | frontend: `git mv` empat folder `src/modul/*` *(termasuk `claimlife` → `claim-life`)*; jalur impor; `daftar.ts`; uji sinkron sidebar ↔ palet | `refactor: nama folder modul frontend sesuai tabel nama` |
| 3 | dokumen: `.env.example`, `PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`, PANDUAN-MENJALANKAN, README, brief induk — tabel §1 disalin utuh; `.env` work owner **tidak** disunting, tetapi laporan menyebut bila `MODUL_AKTIF` di sana perlu diubah | `docs: tabel nama modul` |

Setiap commit: build, vet *(dengan dan tanpa tag `db`)*, test *(dengan dan tanpa tag `db`)*, gofmt, tsc, vitest, vite build — **angka sama dengan §0**.

## 4. LAPORAN

Satu pesan: tabel **paket → commit → folder dipindah → angka uji**; nilai `MODUL_AKTIF` lama dan baru; bab **TELEMETRI EKSEKUSI**.

---

*Disusun 30 September 2026 sesudah verifikasi refactor bentuk B (`3ff9301..5a281cf`, uji dijalankan ulang di `main` `c2fd123`).*
