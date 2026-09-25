# 05a: Simpan premium summary — konversi uang di satu batas, dan urutan procedure yang mengikat

**Status:** ready-for-agent

**Blocked by:** **00 (skema tujuh tabel — PREFACTOR)**, 03 (penomoran — `PL_NUMBER` adalah masukan
rekam polis)

> ⚠️ **Diselaraskan 2026-09-16 — revisi penyimpanan.** Judul asli menyebut *"urutan procedure yang
> mengikat"*; **urutan itu lenyap**. Kedua procedure yang dulu memaksa titik potong —
> `INSERTJSONPOLISLIFE` dan `INSERTJSONOFFERLIFE` — adalah **procedure JSON**, dan keduanya
> **dibuang** (spec §6, §12). Sekarang: **satu polis = satu transaksi**. Lihat blok AC "Penyimpanan
> relasional" di bawah; AC lama tentang urutan **tidak berlaku**.

## Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin nomor premium list dan rekam summary-nya lahir **bersama atau
tidak sama sekali**, supaya tidak pernah ada nomor yang terbit tanpa rekam, maupun rekam tanpa nomor —
dan saya ingin setiap rupiah yang dikirim ke basis data kembali **persis sama** ketika dibaca.
*(User story 25–32 di spec)*

## Area codebase

`internal/repository` (satu transaksi meliputi penomoran + summary; konversi teks ↔ desimal **hanya di
sini**), `internal/services` (perakitan 37 kolom summary), `internal/handlers` (respons memuat
`PL_NUMBER` yang terbit), `frontend/` (tampilan ringkasan premium list).

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path |
| --- | --- | --- |
| `InsertPLSummary` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/InsertPLSummary.xml` **dan** `Endorsement Life/RDBList/InsertPLSummary.xml` |
| `GetSequenceNumber_SQL` | `ASM-FW-GISFW-INT-POLICYJSON` / `RNM!GETSEQUENCENUMBER_SQL` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/GetSequenceNumber_SQL.xml` |
| `InsertJsonPolisLife_Act` | `ASM-FW-GISFW-WORK-LIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/InsertJsonPolisLife_Act.xml` — pemanggil, step **8** |

`[terverifikasi]` **`InsertPLSummary` adalah satu rule yang sama** di kedua modul, bukan dua rule
serupa: identitas identik, `<pxUpdateDateTime>` identik (`20241101T073100.421 GMT`),
`<pyRuleSetVersion>` identik (`01-01-81`). Diff atas dua ekspor yang dinormalisasi menyisakan
**8 baris**, seluruhnya cap waktu ekspor (`pyRuleFormStatusTime`, `pyShowJavaWindowName`).
**PremiumList Life dan Endorsement Life menulis summary lewat jalur yang sama.**

`[terverifikasi]` Bentuk panggilannya:

```
BEGIN
POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY(
  {pyWorkPage.BusinessName},
  {pyWorkPage.PremiumListSummary.PL_NUMBER},
  {pyWorkPage.PremiumListSummary.PL_NUMBER_EDM},
  {TempInputData.CARI2} … {TempInputData.CARI34},
  {pyWorkPage.pzInsKey},
  {InputParam.ERRMSG out}, {InputParam.STSSAVE out});
COMMIT;
END;
```

`[data DBA]` `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY` → `INSERT` ke `M_LIFE_PREMIUM_SUMMARY`, PK dari
`M_LIFE_PREMIUM_SUMMARY_SEQ`, **37 kolom**: identitas (`ID`, `COB`, `PL_NUMBER`, `PL_NUMBER_EDM`,
`CURRENCY`, `IDPEGA`), uang gross (`PREMIUM`, `COMMISSION`, `BROKERAGE_FEE`, `OVR_COMM`, `TAX`,
`PROF_COMM`, `CLAIM`, `CLAIM_AMOUNT`, `BALANCE`, `RI_ADMIN_FEE`, `DEDUCTION`), dan turunan
`*_REFUND` / `*_RETRO` / `*_REFUND_RETRO`. **Procedure tidak commit sendiri**; `INSERT` + rollback
saat error.

⚠️ `[data DBA]` **Seluruh parameter procedure bertipe `VARCHAR2` — termasuk kolom uang.** Uang
menyeberang batas sebagai **teks**.

## ⚠️ Urutan procedure relatif terhadap commit `[keputusan desain]`

Batas transaksi jalur Life **campuran**:

| Procedure | Commit di dalam body? |
| --- | --- |
| `PROC_GENERATE_SEQUENCE_NUMBER` | **tidak** |
| `PEGA_M_LIFE_PREMIUM_SUMMARY` | **tidak** |
| `INSERTJSONPOLISLIFE` | **ya** |
| `INSERTJSONOFFERLIFE` | **ya** |

**Urutan yang ditetapkan:**

1. **Satu transaksi Go** memuat `PROC_GENERATE_SEQUENCE_NUMBER` **dan**
   `PEGA_M_LIFE_PREMIUM_SUMMARY`, lalu **commit**.
2. Barulah `INSERTJSONPOLISLIFE` dipanggil — di luar transaksi itu (tiket **05b**).

⚠️ `[terverifikasi]` **Ini perbaikan yang disengaja, bukan tiruan.** Di Pega, keempat pembungkus
`RULE-CONNECT-SQL` menerbitkan `COMMIT;` **tepat sesudah** panggilan procedure, di dalam blok
`<pyBrowseSQL>` yang sama:

| Pembungkus | Baris `COMMIT;` |
| --- | --- |
| `GetSequenceNumber_SQL` | 88 |
| `InsertPLSummary` | 125 |
| `InsertJsonPolis` | 102 |
| `SaveOfferJsonLife_SQL` | 141 |

Akibatnya **di Pega, nomor dan summary tidak atomik**: kegagalan di antara keduanya meninggalkan
nomor yatim yang sudah ter-commit. Sistem baru menutup celah itu.

## ADR terkait

**ADR-0003** (uang non-float; DDL `NUMBER` tanpa presisi → desimal presisi arbitrer — **diperkuat**
oleh temuan bahwa parameter procedure seluruhnya `VARCHAR2`), **ADR-0006** (penomoran),
**ADR-0015** (Go memegang batas transaksi), **ADR-0011** (bentuk rekam premium yang dikonsumsi hilir).

## Acceptance criteria

- [ ] `PROC_GENERATE_SEQUENCE_NUMBER` dan `PEGA_M_LIFE_PREMIUM_SUMMARY` dipanggil di dalam **satu
      transaksi Go**, dan transaksi itu **commit sebelum** procedure lain dipanggil. *(AC 20 spec)*
- [ ] ~~`INSERTJSONPOLISLIFE` dan `INSERTJSONOFFERLIFE` **tidak** dipanggil dari dalam transaksi
      itu.~~ ⚠️ **TIDAK BERLAKU 2026-09-16** — keduanya **dibuang**; lihat AC pengganti di bawah.

### Penyimpanan relasional ⚠️ BARU 2026-09-16 — spec §6, §12

- [ ] ⚠️ **Satu polis ditulis dalam SATU transaksi**: header, seluruh baris rekap mata uang, seluruh
      peserta, seluruh spreading, seluruh spreading retro, dan seluruh riwayat penawaran — lalu
      **commit sekali**. Kegagalan di tingkat mana pun **membatalkan seluruhnya**. Dibuktikan dengan
      menyuntikkan kegagalan pada baris peserta ke-N dan memastikan **tidak ada** polis tersimpan.
      *(AC 35 spec; penyimpangan sadar 1)*
- [ ] Penomoran (`PROC_GENERATE_SEQUENCE_NUMBER`) berada **di dalam** transaksi itu: tidak pernah ada
      nomor tanpa polis, tidak pernah ada polis tanpa nomor. *(AC 36 spec)*
- [ ] ⚠️ **Tidak ada procedure JSON yang dipanggil.** Test yang menemukan pemanggilan
      `INSERTJSONPOLISLIFE`, `INSERTJSONOFFERLIFE`, atau padanan `@ASM.GetPageJSONString()`
      **gagal**. *(AC 32, 34 spec; penyimpangan sadar 1)*
- [ ] ⚠️ **Rekap uang per mata uang tersimpan** di tabel rekap — bukan dihitung lalu dibuang. Polis
      bermata uang ganda menghasilkan **satu baris rekap per mata uang**, masing-masing dengan nilai
      uangnya. *(AC 37 spec; penyimpangan sadar 2)*
- [ ] ⚠️ `M_LIFE_PREMIUM_SUMMARY`, `M_LIFE_PREMIUM_DETAIL`, dan `LIFEINPRODUCTION` **tidak ditulis**.
      Test yang menemukan tulisan ke ketiganya **gagal**. *(AC 33 spec)*
- [ ] ⚠️ Bila jalur warisan `SaveMasterLPDet` masih dipakai selama transisi, ia berada **di luar**
      transaksi polis dan **dapat diulang** — ia satu-satunya titik potong yang tersisa.
      *(spec §6; `[data DBA]` `COMMIT` di dalam procedure)*
      *(AC 21 spec)*
- [ ] Kegagalan **sebelum** commit summary tidak meninggalkan nomor maupun rekam separuh: test
      menyuntikkan kegagalan di `PEGA_M_LIFE_PREMIUM_SUMMARY` dan memastikan **tidak ada** baris
      `M_LIFE_PREMIUM_SUMMARY` **dan** sequence tidak bergerak. *(AC 22 spec)*
- [ ] Kegagalan **setelah** commit summary meninggalkan nomor + rekam summary **utuh** dan keadaan itu
      **terdeteksi** — bukan senyap. *(AC 23 spec)*
- [ ] **Ada test yang gagal bila urutan pemanggilan diubah** — urutannya bagian dari kebenaran, bukan
      kebetulan. *(AC 24 spec)*
- [ ] Konversi teks ↔ desimal terjadi **hanya di lapisan repository**, di satu tempat; lapisan
      services dan handlers hanya mengenal desimal. *(AC 15 spec; **ADR-0003**)*
- [ ] Ke-37 kolom terisi dari sumber yang benar, dan pemetaannya diuji kolom demi kolom — **bukan**
      lewat posisi `CARI2`…`CARI34` yang tidak bernama.
- [ ] Nilai uang yang dikirim dan dibaca kembali **identik**, termasuk nilai berpecahan panjang dan
      nilai negatif. Tidak ada pembulatan diam. *(AC 16 spec)*
- [ ] Baik `PL_NUMBER` maupun `PL_NUMBER_EDM` tersimpan pada rekam summary yang sama sebagai **dua
      nilai terpisah**; jalur new business mengisi yang pertama, endorsement yang kedua.
      *(AC 13 spec)*
- [ ] Keluaran galat procedure (`ERRMSG`, `STSSAVE`) **diperiksa**; galat yang dilaporkan procedure
      tidak boleh diabaikan sehingga transaksi tampak berhasil.

## Blocker

**Tidak ada pemblokir.** Terkait tetapi tidak memblokir: **OQ-001 sisa** (DDL fisik — memblokir
tiket **09**, bukan tiket ini; kolomnya sudah terbaca dari body procedure).

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```
