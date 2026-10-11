# Modul `endorsmentfacin` — Endorsment Fac In

✅ **Terdaftar sejak 11-10-2026** (irisan I0–I3, keputusan work owner W-EDM-01…12 di
`docs/KEPUTUSAN-WORK-OWNER-EDM.md`): `backend/modul.go`, slot menu `966`, portal endorsement + form *Create
Addendum* dengan pemeriksaan gerbang dan tanggal (`POST /api/endorsmentfacin/periksa`, nol tulis) dan pembuatan
kasus `EDM-<5 digit>` (`POST /api/endorsmentfacin/kasus`): MENULIS `T_WORK_POLIS` (tabel milik premiumlistlife,
dipakai bersama) ber-**LINI `FACEDM`** — sengaja bukan `FAC`, supaya kasus EDM tidak tampil di layar NB FacIn
(W-EDM-09) — dan `T_EDM_FACIN_KASUS`. Kerangka folder lahir dari struktur tim satu folder per modul (keputusan work
owner 30-09-2026); cara memulai modul: `APP_RNM/PANDUAN-TIM-PER-MODUL.md` bab 4.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `endorsmentfacin` |
| Folder korpus | `Endorsment Fac In` |
| GROUPMENU | `FACULTATIVE` |
| Pemilik | `@PEMILIK-ENDORSMENTFACIN` |
| Status | dimigrasi |
| Rentang migrasi | `260-299` |
| Slot menu | `966-967` |
| Prefix rute API | `/api/endorsmentfacin` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | spec, tiket (`issues/`), discovery, grilling — dipindah dari `jefri/OUTPUT FIX/` 30-09-2026 (`git mv`, isi tidak diubah). Asal tiap berkas dan tujuan tautan lamanya: `docs/PETA-ASAL.md` |
| `backend/` | `modul.go`; `models/` (agregat penawaran EDM, kabel periksa), `services/` (before-image E06–E11, gerbang E04, `OpenCase` E03, validasi tanggal E05, Life E18, premi/pembayaran E12–E14 sebagian, `medis/` E20, `predikat/` E01/E02, layanan periksa I2), `repository/` (BACA saja: `POOLDATA` `JSON_POLIS`, `FACINPRODUCTION`, `FACOUTPRODUCTION`, `OS_AKSEPTASI_KLAIM`, `ARASAPAS.DETAIL_INVOICE`, `M_EDM_PENGECUALIAN_POLIS`), `handlers/` (`POST /api/endorsmentfacin/periksa`, `POST /api/endorsmentfacin/kasus`), `tiruan/` (gudang memori untuk uji), `migrations/`. MENULIS: `T_WORK_POLIS` (LINI `FACEDM`, POSITION `ReasFacInMarketing`, STATUS_WORK `Pending-Policy`, FLAG_ONGOING_POLICY `1`) + `T_EDM_FACIN_KASUS`, satu transaksi |
| `frontend/` | `menu.ts`, `rute.tsx`, `api.ts`, `labels.ts` (verbatim korpus, diuji), `pages/PortalEndorsement.tsx`, `pages/FormAddendum.tsx`, `endorsmentfacin.css` |

## Migrasi

Terpakai: **260** `SEQ_WORK_POLIS_EDM` (START WITH 1 — nomor pertama `EDM-00001`, W-EDM-08), **261**
`M_EDM_PENGECUALIAN_POLIS` (tabel konfigurasi pengecualian nomor polis, isinya diisi DBA — W-EDM-03; tanpa tabel ini
`/periksa` tetap jalan dengan daftar kosong + catatan), **262** `T_EDM_FACIN_KASUS` (kepala kasus EDM, PK = ID
`T_WORK_POLIS`, tanpa FK), **966** slot menu. Nomor berikut = 263. Urutan jalan: 260 → 261 → 262 → 966 — semuanya
dijalankan DBA, bukan agent. Tanpa 260/262 `POST …/kasus` menjawab 503; tanpa 262 gerbang 2 `belum-diperiksa`.
⛔ **260, 261, 262, dan 966 SUDAH JALAN di DEV** (`T_MIGRASI` 11-10-2026 17:12 / 17:25; dibaca sesi asisten
`nusantarare-0b`) — berkasnya tidak boleh diubah lagi; ralat lewat migrasi baru (263 dst.). Isi
`M_EDM_PENGECUALIAN_POLIS` masih kosong (menunggu DBA).

Rentang `260-299` (tabel R2, urut hulu ke hilir: migrasi modul hilir yang merujuk tabel modul hulu
selalu berjalan sesudahnya). Slot menu `966-967` hanya menyalakan `DIMIGRASI` baris modul ini (satu `UPDATE`,
nol `INSERT` — menu datar 30-09-2026) saat modul mendapat layar pertamanya, di folder
`backend/migrations/` modul ini sendiri — bentuk SQL-nya di `APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`
bab 6. Nomor selalu tiga digit.
