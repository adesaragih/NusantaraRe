---
status: accepted
tanggal: 2026-09-14
sumber: grilling Ronde 2 Q13 (`.scratch/claim-life/grilling-ronde-2.md`), keputusan work owner
---

# Seluruh data Claim — Life dipindahkan; tidak ada koeksistensi dua penulis

**Semua data dipindah** ke sistem baru. Opsi koeksistensi — sistem baru hanya menerima klaim baru
sementara klaim berjalan diselesaikan di Pega — **ditolak**.

## Considered Options

- **Migrasi penuh seluruh data** — dipilih
- **Koeksistensi**: hanya klaim baru di sistem baru; klaim berjalan (Outstanding, Medical Check,
  Claim Analis) diselesaikan di Pega — **ditolak**

Koeksistensi sempat diusulkan karena secara teknis mungkin: kedua sistem menulis ke tabel yang sama,
`POOLDATA.OS_AKSEPTASI_KLAIM_LIFE`, lewat rule yang identik
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL!RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL`, hash ternormalisasi
`c50bfd9a12`, `[terverifikasi]` sama persis di `Claim Life` dan `Komite Claim Life`).
Penolakannya dicatat di sini supaya tidak diusulkan ulang.

## Consequences

- **Tidak ada periode dua penulis** ke `OS_AKSEPTASI_KLAIM_LIFE` dari sisi Claim Life — cutover
  bersifat memutus, bukan bertahap.
- Klaim yang sedang berada di tengah siklus **harus terbawa beserta statusnya**. Statusnya adalah
  `STS_REJECT` (`0` = Outstanding) ditambah posisi tahap (Register / Outstanding / Medical Check /
  Claim Analis) — lihat `CONTEXT.md` §Status dan kode. Kasus yang sudah diserahkan ke Komite dan
  belum kembali membawa keadaan lintas batas (**ADR-0001**), sehingga pemindahannya perlu
  disepakati dengan pemilik konteks Komite Life.
- Penomoran klaim tetap memanggil `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` (**ADR-0006**).
  **Migrasi tidak boleh membuat sequence melompat atau mengulang** — perilaku procedure terhadap
  data yang dipindahkan perlu dipastikan bersama DBA.
- Jejak audit baru (**ADR-0007**) **tidak dapat direkonstruksi ke belakang**: data lama hanya
  memiliki `CREATEOPNAME` dan empat kolom tanggal. Riwayat transisi lengkap hanya ada untuk kejadian
  setelah cutover.

## OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-001** | Tidak ada DDL — struktur dan tipe kolom sumber tidak diketahui, sehingga rencana pemindahan belum dapat disusun rinci |
| **OQ-002** (dipersempit) | Kontrak `PROC_GENERATE_SEQUENCE_NUMBER` — apakah sequence di-reset per tahun/lini/global menentukan aman-tidaknya migrasi nomor |
| **OQ-018** (terjawab untuk Claim — Life) | `jboss1073` = production, `jboss117` = dev, sistem mirroring. Lingkungan `pega-nusre` (73 berkas di modul ini) **belum dinyatakan** — perlu dipastikan sebelum menentukan sumber data migrasi |
| **OQ-060** | Cakupan `CURRENCY` — memengaruhi bentuk data uang yang dipindahkan |
