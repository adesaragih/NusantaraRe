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
| Slot menu | `—` |
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
| `backend/` | `modul.go` (`Pendaftaran`), `models/` (murni; `templat/` = stream HTML `EmailKlaim_HTML_KMT` VERBATIM; `pdf_akseptasi.go` = PDF `FILEAcceptanceNote`), `repository/` (SQL), `services/` (aturan + satu transaksi Submit + perakit isi efek + pelaksana outbox), `handlers/`, `tiruan/` (uji, kontrak palsu Claim Prop), `konfigurasi/kasir.json` / `email.json`, `migrations/` |
| `frontend/` | `layar.ts` (modul TANPA menu), `rute.tsx`, layar `ShowTransfer` (dibuka dari tabel komite inbox Claim Prop) |

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
literal kode; env tidak dipakai (env hanya dibaca `inti/backend/config`, ADR-U-0013).

`backend/konfigurasi/email.json` — akun notifikasi `SendEmailKlaim_KMT` S17 (`NUSARE`) / S18 (`NUSARESYARIAH`, bila
alamat tujuan memuat "syariah") dan CC kotak surat klaim S4 (hanya IsPEGAPROD). BCC pribadi yang di-hardcode XML (S3)
**tidak disalin**.

## Efek keluar (outbox, hanya produksi)

| Jenis | Asal | MUATAN | Pelaksana |
| --- | --- | --- | --- |
| `konversi-klaim` | `KonversiKlaim_Act` (S29) | CASEID, NOPOLIS, STS_REJECT | `ErrArasapasBelumDisetujui` |
| `kasir` | `HitServiceToKasirKMT_Act` (S34) | muatan `SendAcceptationToKasir` | `ErrKasirBelumDisetujui` |
| `email-komite` | `SendEmailKlaim_KMT` (S35) | jenis, ID akun penerima, ID baris tangga | rakit isi (`SusunEmailKomite`) → `ErrEmailBelumDisetujui` |

MUATAN email hanya pengenal (claimlife/015: tanpa nama / alamat); isi dirakit saat dikirim.

## Dokumen akseptasi (bukan outbox)

`PrintFileAcceptance_TKMT` (S21 → S8) di SETIAP lingkungan, sesudah Submit tingkat akhir tersimpan: PDF
(`github.com/go-pdf/fpdf`, `models.PDFAcceptanceNote`) diunggah lewat `inti/backend/penyimpanan` (folder Claim), lalu
T_STORAGE_IMAGE dan baris tabel warisan dokumen klaim dicatat di satu transaksi - pola lampiran Bordereaux. Gagal =
keputusan tetap tersimpan, layar menerima `galatDokumen`.

## Menu

**Tidak ada** (perintah work owner 09-10-2026: "kode menu nya di hapus dari repo, anggap menu itu tidak pernah ada,
karena digabung di menu klaim nya masing-masing"). Baris `M_NAV_MENU` dan hak akunnya dibuang migrasi inti 949; slot
menu 986 dihapus. Modul ini dipasang bagi pemegang menu `claimprop` (`MODUL_DIPINJAM` frontend/App.tsx, frontend
`layar.ts`) dan rutenya dipinjam (`ruteDipinjam` cmd/api/rakit.go); kasus komite dibuka dari tabel komite inbox Claim
Prop.

## Migrasi

Rentang `680-719`. `-migrate` dijalankan **work owner**.

| Nomor | Isi |
| --- | --- |
| 680 | `T_GENERAL_KOMITE` ADD `KOMITE_USUL_TUTUP` / `KOMITE_USUL_CADANG` `CHAR(1) DEFAULT '0' NOT NULL` + CHECK `'0'`/`'1'` |
| 681 | `UX_GENERAL_KOMITE_ADJ` (unik, `claimlife/013`) → `IX_GENERAL_KOMITE_ADJ` biasa (keputusan work owner 08-10-2026) |
| 682 | `T_GENERAL_KOMITE` ADD `KOMITE_SUBJECTIVITY` `CHAR(1) DEFAULT '0' NOT NULL` + CHECK, `KOMITE_SUBJECTIVITY_NOTE` `VARCHAR2(1000)` — isian Subjectivity tingkat 1 antar tingkat (keputusan work owner 08-10-2026, OQ-KCP-01 "a") |

## Data lama

Kasus komite warisan Pega **tidak dimigrasi** (keputusan work owner 08-10-2026, OQ-KCP-02 "b"); klaimnya dimuat pemuat
Claim Prop (baris berlaku `STS_REJECT = 1` → Input Acceptation). Nol alat pemuat di modul ini.
