# Modul `komiteclaimnonprop` — Komite Claim Non Prop

Komite klaim treaty inward non proporsional (XoL): kasus `ASM-FW-GCNMFW-Work-KomiteTreatyNonProp`
(`Flow/KomiteTreaty_Flow.xml` korpus `Komite Claim Non Prop`: assignment "KomiteRouter" → `KomitePostAdjustment` →
Resolved-Completed). Kasus `KMTNP-` **dilahirkan Claim Non Prop** (`CreateChildKomiteCNP_Act`); modul ini menampilkan
layar `ShowTransfer` dan menjalankan `KomitePostAdjustment` dalam satu transaksi. Dibangun 09-10-2026 atas perintah
work owner — pola Komite Claim Prop: **tanpa menu sendiri** (menu disembunyikan lewat data), daftar kasusnya tampil di
bawah inbox Claim Non Prop dan dibuka di tempat; jenjang komite memakai workbasket (roster NONPROP, migrasi
claimnonprop 611).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `komiteclaimnonprop` |
| Folder korpus | `Komite Claim Non Prop` |
| GROUPMENU | `KLAIM` |
| Pemilik | `@PEMILIK-KOMITECLAIMNONPROP` |
| Status | dimigrasi |
| Rentang migrasi | `720-759` |
| Slot menu | `988-989` |
| Prefix rute API | `/api/komite-claim-non-prop` (`/api/komite` milik Komite Claim Life, `/api/komite-claim-prop` milik Komite Claim Prop) |
| Kontrak disediakan | — |
| Kontrak dipakai | `kontrak.KlaimTreatyNonPropKomite` (`inti/backend/kontrak/klaimtreatynonprop.go`, disediakan `claimnonprop`) — perintah work owner 09-10-2026 |

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `PARITAS.md` (setiap isian / tombol / langkah `KomitePostAdjustment` ↔ status) |
| `backend/` | `modul.go` (`Pendaftaran`), `models/` (murni; `templat/` = stream HTML `EmailKlaim_HTML_KMT`), `repository/` (SQL), `services/` (aturan + satu transaksi Submit + perakit isi efek + pelaksana outbox), `handlers/`, `tiruan/` (uji, kontrak palsu Claim Non Prop), `konfigurasi/kasir.json` / `email.json`, `migrations/` |
| `frontend/` | `menu.ts`, `rute.tsx`, daftar kerja penyetuju, layar `ShowTransfer` |

## Rute

Ketiganya juga **rute pinjaman** pemegang menu `claimnonprop` (`cmd/api/rakit.go` `ruteDipinjam`): tabel komite di
bawah inbox Claim Non Prop dan layar keputusan yang dibuka di tempat.

| Metode | Jalur | Isi |
| --- | --- | --- |
| GET | `/api/komite-claim-non-prop/kasus` | daftar kerja pelaku (KomiteRouter S6.1: baris tangga pertama yang menunggu, akun atau workbasket aktif) |
| GET | `/api/komite-claim-non-prop/kasus/{id}` | layar `ShowTransfer` (pra-proses `SetKomiteList_Act`) |
| POST | `/api/komite-claim-non-prop/kasus/{id}/putuskan` | Submit → `KomitePostAdjustment`: 403 bukan pemegang, 409 kasus / klaim tertutup atau berubah, 422 isian |

## Batas modul

- Klaim induk dibaca dan ditulis **hanya** lewat `kontrak.KlaimTreatyNonPropKomite`; nol impor `modul/claimnonprop`.
- Satu-satunya berkas yang menyebut tabel tangga: `backend/repository/tangga.go` (entri penjaga Claim Life
  `komite_statik_test.go`, izin work owner 09-10-2026).
- Semua kueri menyaring `w.LINI = 'NONPROP'` (ketat) dan awalan `KMTNP-`.
- Kolom kepala kasus `KOMITE_USUL_TUTUP` / `KOMITE_USUL_CADANG` / `KOMITE_SUBJECTIVITY(_NOTE)` milik migrasi
  komiteclaimprop 680 / 682 (tabel bersama `T_GENERAL_KOMITE`).

## Konfigurasi berdokumen

`backend/konfigurasi/kasir.json` — kode tetap muatan `SendAcceptationToKasir` (`HitServiceToKasirKMT_Act` S10.3 cabang
IsCLMNP: CompanyName / LjtdId / LdcId; S10.4 IsPEGASyariah LdcId). Konfigurasi, bukan literal kode; env tidak dipakai
(env hanya dibaca `inti/backend/config`, ADR-U-0013).

`backend/konfigurasi/email.json` — akun notifikasi `SendEmailKlaim_KMT` (`NUSARE` / `NUSARESYARIAH` bila alamat tujuan
memuat "syariah") dan CC kotak surat klaim (hanya IsPEGAPROD). BCC pribadi yang di-hardcode XML **tidak disalin**.

## Efek keluar (outbox, hanya produksi)

| Jenis | Asal | MUATAN | Pelaksana |
| --- | --- | --- | --- |
| `konversi-klaim` | `KonversiKlaim_Act` (S14.22) | CASEID, NOPOLIS, STS_REJECT "1" | `ErrArasapasBelumDisetujui` |
| `kasir` | `HitServiceToKasirKMT_Act` (S19.3) | muatan `SendAcceptationToKasir` per mata uang Spreading In | `ErrKasirBelumDisetujui` |
| `email-komite` | `SendEmailKlaim_KMT` (S19.2) | jenis, ID akun penerima, ID baris tangga | rakit isi (`SusunEmailKomite`) → `ErrEmailBelumDisetujui` |

MUATAN email hanya pengenal (claimlife/015: tanpa nama / alamat); isi dirakit saat dikirim.

## Migrasi

Rentang `720-759`. `-migrate` dijalankan **work owner**.

| Nomor | Isi |
| --- | --- |
| 988 | slot menu: `UPDATE M_NAV_MENU SET DIMIGRASI = '1' WHERE KODE = 'komiteclaimnonprop'` |

Menu `komiteclaimnonprop` **disembunyikan lewat data** (langkah work owner, bukan migrasi — pola Komite Claim Prop):
`M_NAV_MENU.STATUS_AKTIF = '0'` untuk kode itu dan kode menu dicabut dari akun. Pemegang menu `claimnonprop` membuka
kasus komite lewat rute pinjaman.

## Data lama

Kasus komite warisan Pega **tidak dimigrasi** (pola Komite Claim Prop). Nol alat pemuat di modul ini.
