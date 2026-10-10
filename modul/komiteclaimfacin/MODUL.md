# Modul `komiteclaimfacin` — Komite Claim FacIn

Komite klaim fakultatif inward: kasus `ASM-FW-GCNMFW-Work-Komite` (`Flow/Komite_Flow.xml` korpus `Komite Claim FacIn`:
assignment "KomiteRouter" → `ViewTransferDtl` (`SetValueKomite` / `ShowTransfer` / `KomitePostAct`) → decision
IsKomiteLoop). Kasus `KMT-` **dilahirkan Claim Fac In**: TT2 `CreateKMTNo_Act` (penyerahan adjustment), TT3
`SendRejectClaimToKomite2` (Reject Claim), TT4 `SendCloseClaimToKomite` (Close Without Payment). Modul ini menampilkan
layar `ShowTransfer` dan menjalankan `KomitePost_Adjustment` / `KomitePost_Reject` / `KomitePost_CloseClaim` dalam satu
transaksi Submit. Dibangun 10-10-2026 (prompt work owner "IMPLEMENTASI KOMITE CLAIM FAC IN, TAHAP 2 DARI 2") — pola
Komite Claim Prop / Non Prop: **tanpa menu** (modul `layar.ts`, baris menunya dibuang migrasi inti 949), tabel kasusnya
tampil di bawah inbox Claim Fac In dan dibuka di tempat; jenjang komite memakai workbasket (roster FACIN, migrasi 640).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `komiteclaimfacin` |
| Folder korpus | `Komite Claim FacIn` |
| GROUPMENU | `KLAIM` |
| Pemilik | `@PEMILIK-KOMITECLAIMFACIN` |
| Status | dimigrasi |
| Rentang migrasi | `640-679` |
| Slot menu | `—` |
| Prefix rute API | `/api/komite-claim-fac-in` (`/api/komite` milik Komite Claim Life, `/api/komite-claim-prop` Komite Claim Prop, `/api/komite-claim-non-prop` Komite Claim Non Prop) |
| Kontrak disediakan | — |
| Kontrak dipakai | `kontrak.KlaimFacInKomite` (`inti/backend/kontrak/klaimfacin.go`, disediakan `claimfacin`) — prompt work owner tahap 2 §2 butir 1 |

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `PARITAS.md` (setiap isian / tombol / langkah ↔ status), `OQ.md`, `STRUKTUR-TABEL-KOMITE-CLAIM-FACIN.md`, `spec.md` + `issues/` + `keputusan-sebelum-to-spec.md` (September, ber-RALAT 10-10-2026), grilling (tersegel) |
| `backend/` | `modul.go` (`Pendaftaran`), `models/` (murni: kasus, pra-proses `SetValueKomite`, perluasan tangga `ApprovalKomite_Act`, rencana Submit, layar, OS / kasir / surel), `repository/` (SQL; `tangga.go` satu-satunya penyebut tabel tangga), `services/` (aturan + satu transaksi Submit + pelaksana outbox), `handlers/`, `tiruan/` (uji, kontrak palsu Claim Fac In), `migrations/`, `konfigurasi/kasir.json` / `email.json` |
| `frontend/` | `layar.ts` (modul TANPA menu), `rute.tsx`, layar `ShowTransfer` (`pages/KasusKomite.tsx`, dibuka dari tabel komite inbox Claim Fac In) |

## Rute

Ketiganya juga **rute pinjaman** pemegang menu `claimfacin` (`cmd/api/rakit.go` `ruteDipinjam`): tabel komite di bawah
inbox Claim Fac In dan layar keputusan yang dibuka di tempat.

| Metode | Jalur | Isi |
| --- | --- | --- |
| GET | `/api/komite-claim-fac-in/kasus` | daftar kerja pelaku (baris tangga tingkat berjalan `KOMITE_URUT = KOMITE_COUNT` yang menunggu, akun atau workbasket aktif; anggota `ReasClaimSPVB` juga melihat tingkat `ReasClaimSPVA`) |
| GET | `/api/komite-claim-fac-in/kasus/{id}` | layar `ShowTransfer` (pra-proses `SetValueKomite` + perluasan tangga ditampilkan, nol tulisan) |
| POST | `/api/komite-claim-fac-in/kasus/{id}/putuskan` | Submit → `KomitePostAct`: 403 bukan pemegang, 409 kasus / klaim tertutup atau berubah, 422 isian |

## Keputusan work owner 10-10-2026

- **KCF-01** roster `EMAILKOMITE` STS_KLAIM FACIN → workbasket (UPDATE di tempat menurut DEGREE + JABATAN, migrasi 640):
  1 `ReasClaimSPVA`, 2 `ReasClaimDeptHead`, 3 `ReasClaimTechDivHead`, 4 `ReasClaimOpsDir`, 5 `ReasClaimTechDir`. Anggota
  `ReasClaimSPVB` juga memutus tingkat 1. Pita SPV B: pemutus tingkat 1 anggota `ReasClaimSPVB` dan 30.000.000 < total
  ≤ 57.750.000 → satu jenjang atas.
- **KCF-02** perluasan tangga ikut XML (`ApprovalKomite_Act` L4–L8): dihitung saat tingkat 1 membuka (tanpa tulisan
  saat GET), disimpan saat tingkat 1 Submit; Fac Retro melewatinya.
- **KCF-03** TT3 / TT4 untuk Fac In; `T_GENERAL_KOMITE.ADJUSTMENT_ID` boleh kosong (641) + `TRANSFER_TYPE` (642); satu
  tingkat `ReasClaimDeptHead`.
- **KCF-04** kasus komite lama Pega tidak dimigrasi.
- Kelainan XML diperbaiki (prompt §5, `docs/PARITAS.md` §7); literal data ("AKSEPATSI", "Febuari", "July") tetap.

## Batas modul

- Klaim induk dibaca dan ditulis **hanya** lewat `kontrak.KlaimFacInKomite`; nol impor `modul/claimfacin`. Daftar
  putih jalur tulis balik: `docs/PARITAS.md` §9.
- Satu-satunya berkas yang menyebut tabel tangga: `backend/repository/tangga.go` (entri penjaga Claim Life
  `komite_statik_test.go`, izin work owner 10-10-2026).
- Semua kueri menyaring `w.LINI = 'FACIN'` (ketat) dan awalan `KMT-` (fixture Claim Life memakai `KMT-` ber-LINI LIFE).
- Kolom kepala kasus `KOMITE_USUL_TUTUP` / `KOMITE_USUL_CADANG` milik migrasi komiteclaimprop 680; `TRANSFER_TYPE` milik
  migrasi 642 modul ini (tabel bersama `T_GENERAL_KOMITE`; baris kolom di STRUKTUR Claim Life / Komite Claim Life /
  Komite Claim Prop).

## Konfigurasi berdokumen

`backend/konfigurasi/kasir.json` — kode tetap muatan `SendAcceptationToKasir` (`HitServiceToKasirKMT_Act`
S14.2.2.1.1.3 jalur IsCLM: CompanyName / LjtdId / LdcId; S14.2.2.1.1.4 IsPEGASyariah LdcId). `[penyimpangan sadar]`
CLAUDE.md §10: konfigurasi, bukan literal kode; env tidak dipakai (env hanya dibaca `inti/backend/config`, ADR-U-0013).

`backend/konfigurasi/email.json` — akun notifikasi `SendEmailKlaim_KMT` / KomitePost_Reject / _CloseClaim S13.10 /
S11.10 (`NUSARE` / `NUSARESYARIAH`) dan CC kotak surat klaim (hanya IsPEGAPROD). BCC pribadi yang di-hardcode XML
**tidak disalin**.

## Efek keluar (outbox, hanya produksi)

| Jenis | Asal | MUATAN | Pelaksana |
| --- | --- | --- | --- |
| `konversi-klaim` | `KonversiKlaim_Act` (TT2 S17 STSREJECT "1", bukan Fac Retro; TT3 S15 "2"; TT4 S12.4 "4") | CASEID, NOPOLIS, STS_REJECT | `ErrArasapasBelumDisetujui` |
| `kasir` | `HitServiceToKasirKMT_Act` (TT2 S25) | muatan `SendAcceptationToKasir` per `AcceptedNo` | `ErrKasirBelumDisetujui` |
| `email-komite` | `SendEmailKlaim_KMT` (S7.2.1.15), TT3 S13.10, TT4 S11.10 | jenis, KomiteID / akun penerima | rakit isi (`SusunEmailKomite`) → `ErrEmailBelumDisetujui` |

MUATAN email hanya pengenal (claimlife/015: tanpa nama / alamat); isi dirakit saat dikirim. Badan email dan PDF menunggu
stream (OQ-KCFI-01). Pekerja outbox tidak dijalankan modul ini (`TanpaPekerja`, pola Komite Claim Prop).

## Menu

**Tidak ada** (perintah work owner 09-10-2026: "kode menu nya di hapus dari repo, anggap menu itu tidak pernah ada,
karena digabung di menu klaim nya masing-masing"). Baris `M_NAV_MENU` dan hak akunnya dibuang migrasi inti 949; slot
menu 984–985 dihapus. Modul ini dipasang bagi pemegang menu `claimfacin` (`MODUL_DIPINJAM` frontend/App.tsx, frontend
`layar.ts`) dan rutenya dipinjam (`ruteDipinjam` cmd/api/rakit.go). Tabel komite inbox tampil bila
`GET /api/claim-fac-in/hak` → `komite`.

## Migrasi

Rentang `640-679` (tabel R2, urut hulu ke hilir). Nol COMMIT / PL-SQL; `-migrate` dijalankan work owner. Berjalan sesudah
claimprop 537 (empat workbasket jenjang 2–5) dan claimfacin 560–567.

| Nomor | Isi | Mundur |
| --- | --- | --- |
| `640` | `M_WORKBASKET` `ReasClaimSPVA` / `ReasClaimSPVB` bila belum ada; `EMAILKOMITE` FACIN `OPERATOR_ID` / `NAME` = workbasket per DEGREE + JABATAN, `EMAIL` kosong (KCF-01) | baris FACIN ber-workbasket dikosongkan; workbasket tidak dibuang |
| `641` | `T_GENERAL_KOMITE.ADJUSTMENT_ID` boleh kosong (KCF-03) | `NOT NULL NOVALIDATE` |
| `642` | `T_GENERAL_KOMITE.TRANSFER_TYPE CHAR(1) DEFAULT '2' NOT NULL` + `CK_GENERAL_KOMITE_TRANSFER` (KCF-03) | kolom + CHECK dibuang |

## Data lama

Kasus komite warisan Pega **tidak dimigrasi** (KCF-04). Nol alat pemuat di modul ini.
