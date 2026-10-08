# Modul `komiteclaimprop` — Komite Claim Prop

Komite klaim treaty inward proporsional: kasus `ASM-FW-GCNMFW-Work-KomiteTreaty` (`Flow/KomiteTreaty_Flow.xml`:
assignment "KomiteRouter" → Decision `KomiteLoop` → Resolved-Completed). Kasus `TKMT-` **dilahirkan Claim Prop**
(`AddKomiteTreatyChild_ACT`, opsi B 07-10-2026); modul ini menampilkan daftar kerja penyetuju, layar `ShowTransfer`
wajah TT 2 (ADJUSTMENT), dan menjalankan `KomitePostAdjustment` dalam satu transaksi. Dimulai 08-10-2026
(`OUTPUT_HASIL_RNM/_brief/PROMPT-IMPLEMENTASI-MODUL-KOMITE-CLAIM-PROP.md`).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `komiteclaimprop` |
| Folder korpus | `Komite Claim Prop` |
| GROUPMENU | `KLAIM` |
| Pemilik | `@PEMILIK-KOMITECLAIMPROP` |
| Status | dimigrasi |
| Rentang migrasi | `680-719` |
| Slot menu | `986-987` |
| Prefix rute API | `/api/komite-claim-prop` (`/api/komite` milik Komite Claim Life) |
| Kontrak disediakan | — |
| Kontrak dipakai | `kontrak.KlaimTreatyKomite` (`inti/backend/kontrak/klaimtreaty.go`, disediakan `claimprop`) — keputusan work owner 08-10-2026 |

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | spec, tiket (`issues/`), grilling (tersegel), STRUKTUR / RELASI, `PARITAS.md` (setiap isian / tombol / langkah `KomitePost*` ↔ tiket ↔ status) |
| `backend/` | `modul.go` (`Pendaftaran`), `models/` (murni), `repository/` (SQL), `services/` (aturan + satu transaksi Submit + pelaksana outbox), `handlers/`, `tiruan/` (uji, kontrak palsu Claim Prop), `alat/pemuatlama/` (data lama, uji-kering), `konfigurasi/kasir.json`, `migrations/` |
| `frontend/` | `menu.ts`, `rute.tsx`, daftar kerja penyetuju, layar `ShowTransfer` |

## Rute

| Metode | Jalur | Isi |
| --- | --- | --- |
| GET | `/api/komite-claim-prop/kasus` | daftar kerja pelaku (KomiteRouter S6.1: baris tangga pertama yang menunggu) |
| GET | `/api/komite-claim-prop/kasus/{id}` | layar `ShowTransfer` (pra-proses `SetKomiteList_Act`) |
| POST | `/api/komite-claim-prop/kasus/{id}/putuskan` | Submit → `KomitePostAdjustment`: 403 bukan pemegang, 409 kasus / klaim tertutup atau berubah, 422 isian |

## Batas modul

- Klaim induk dibaca dan ditulis **hanya** lewat `kontrak.KlaimTreatyKomite`; nol impor `modul/claimprop`.
- Satu-satunya berkas yang menyebut tabel tangga: `backend/repository/tangga.go` (entri penjaga Claim Life
  `komite_statik_test.go`, izin work owner 08-10-2026).
- Semua kueri menyaring `w.LINI = 'PROP'` (ketat) dan awalan `TKMT-`.

## Konfigurasi berdokumen

`backend/konfigurasi/kasir.json` — kode tetap muatan `SendAcceptationToKasir` (`HitServiceToKasirKMT_Act` S14.1.3:
CompanyName / LjtdId / LdcId; S14.1.4 IsPEGASyariah LdcId). `[penyimpangan sadar]` CLAUDE.md §10: konfigurasi, bukan
literal kode; env tidak dipakai (env hanya dibaca `inti/backend/config`, ADR-U-0013). Alamat CC/BCC pribadi yang
di-hardcode XML (`SendEmailKlaim_KMT` S3-S4) **tidak disalin**.

## Migrasi

Rentang `680-719`. `-migrate` dijalankan **work owner**.

| Nomor | Isi |
| --- | --- |
| 680 | `T_GENERAL_KOMITE` ADD `KOMITE_USUL_TUTUP` / `KOMITE_USUL_CADANG` `CHAR(1) DEFAULT '0' NOT NULL` + CHECK `'0'`/`'1'` |
| 681 | `UX_GENERAL_KOMITE_ADJ` (unik, `claimlife/013`) → `IX_GENERAL_KOMITE_ADJ` biasa (keputusan work owner 08-10-2026) |
| 986 | slot menu: `UPDATE M_NAV_MENU SET DIMIGRASI = '1' WHERE KODE = 'komiteclaimprop'` |
