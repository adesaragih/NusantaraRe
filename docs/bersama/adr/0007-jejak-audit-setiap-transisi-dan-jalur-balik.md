---
status: accepted
tanggal: 2026-09-14
sumber: grilling Ronde 2 Q11 (`.scratch/claim-life/grilling-ronde-2.md`), keputusan work owner
---

# Jejak audit merekam siapa + kapan untuk setiap transisi status dan setiap jalur balik

Sistem baru merekam **siapa** dan **kapan** untuk **setiap** transisi status klaim **dan setiap
jalur balik** (`SendtoAdmin`, `SendtoMedical`). Ini **penyimpangan sadar dari sistem lama** —
sebuah perbaikan, bukan paritas.

## Keadaan sekarang, dan apa yang hilang `[terverifikasi]`

Yang direkam Pega pada rekam akseptasi (`POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`, ditulis oleh
`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL`):

| Kolom | Isi |
| --- | --- |
| `CREATEOPNAME` | satu nama operator |
| `ACCEPTATION_DATE`, `CONFIRMATION_DATE`, `CLAIM_RECEIVED_DATE`, `COMPLETE_DATE` | empat tanggal |

Identitas pengguna juga dipakai sebagai **data jejak** di
`Claim Life/Activity/RejectOSClaimLife_Act.xml` dan `Claim Life/Activity/SendEmailKlaimLF.xml`
(`OperatorID.pyUserIdentifier`, `pyUserName`), dan ada pencatatan layanan di
`Claim Life/RDBList/InsertLogServiceClaim.xml` → `INSERT INTO pooldata.monitoring_klaim_log`.

**Yang tidak terekam sekarang:** siapa yang mengembalikan kasus lewat `SendtoAdmin` atau
`SendtoMedical`. Kedua penanda hanya menyimpan **nilai `1`**, tanpa pelaku dan tanpa waktu:

- `Claim Life/When/IsSendtoAdmin.xml` → `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOADMIN` /
  `RULE-OBJ-WHEN` — `pyWorkPage.SendtoAdmin = 1`
- `Claim Life/When/IsSendtoMedical.xml` → `…` / `ISSENDTOMEDICAL` / `RULE-OBJ-WHEN`

Jalur balik itu bukan kejadian langka: `SendtoAdmin` menggerbangi pengembalian dari **tiga titik**
di `Claim Life/Flow/Register_Flow.xml` (`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `REGISTER_FLOW` /
`RULE-OBJ-FLOW`).

## Considered Options

- **Rekam siapa + kapan untuk setiap transisi dan setiap jalur balik** — dipilih
- Paritas dengan Pega (`CREATEOPNAME` + empat tanggal) — ditolak: pengembalian kasus menjadi tidak
  dapat ditelusuri, padahal ia jalur yang sering dipakai

## Consequences

- Ini **penyimpangan sadar** dari non-goal "tidak merapikan alur" (Ronde 1 Q2). Alurnya tetap sama;
  yang bertambah adalah perekamannya.
- Jejak audit menjadi **riwayat transisi**, bukan sekadar beberapa kolom tanggal pada satu rekam.
  Bentuk penyimpanannya adalah keputusan implementasi, bukan ditetapkan di sini.
- Peran pelaku terbatas pada tiga peran di **ADR-0002**; pencatatan pelaku memakai identitas akun,
  bukan nama ter-hardcode.
- Rekam akseptasi lama tetap ditulis sebagaimana adanya — kontrak dengan Komite (**ADR-0001**)
  tidak berubah karena keputusan ini.

## OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-001** | Tidak ada DDL — struktur `monitoring_klaim_log` dan tipe kolom tanggal tidak diketahui |
| **OQ-013** | `COMMIT` di dalam blok PL/SQL — batas transaksi ada di sisi database, sehingga atomisitas "tulis akseptasi + tulis jejak audit" perlu ditetapkan bersama DBA |
