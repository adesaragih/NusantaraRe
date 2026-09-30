# Modul `komiteclaimfacin` — Komite Claim FacIn

⚠️ **Kerangka — belum dimigrasi.** Folder ini dibuat struktur tim satu folder per modul (keputusan work
owner 30-09-2026) supaya pemilik, rentang migrasi, dan slot menu modul ini TETAP sejak awal — satu
modul, satu folder, satu pemilik. Belum ada kode: tanpa `backend/modul.go` modul ini tidak terdaftar
(daftar Go bangkitan `inti/backend/daftar`, `import.meta.glob` frontend), dan kelompoknya di sidebar
tetap "belum dimigrasi" (`M_NAV_MENU.DIMIGRASI = '0'`). Cara memulainya:
`docs/bersama/PANDUAN-TIM-PER-MODUL.md` (akar repo) bab 4.

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `komiteclaimfacin` |
| Folder korpus | `Komite Claim FacIn` |
| GROUPMENU | `KLAIM` |
| Pemilik | `@PEMILIK-KOMITECLAIMFACIN` |
| Status | belum dimigrasi |
| Rentang migrasi | `640-679` |
| Slot menu | `984-985` |
| Prefix rute API | — (ditetapkan spec modul ini) |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

`Pemilik` adalah penanda; akun sebenarnya diisi work owner di `.github/CODEOWNERS`.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | spec, tiket (`issues/`), grilling, catatan — dulu `.scratch/komite-claim-facin/` (dipindah dengan `git mv`, isi byte-identik) |
| `backend/` | belum ada — lahir bersama `backend/modul.go` (`Pendaftaran()`) saat modul dimulai |
| `frontend/` | belum ada — lahir bersama `frontend/menu.ts` dan `rute.tsx` saat modul dimulai |

## Migrasi

Rentang `640-679` (tabel R2, urut hulu ke hilir: migrasi modul hilir yang merujuk tabel modul hulu
selalu berjalan sesudahnya). Butir menu modul ini ditulis di slot `984-985`, di folder
`backend/migrations/` modul ini sendiri — bentuk SQL-nya di `APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`
bab 6. Nomor selalu tiga digit.
