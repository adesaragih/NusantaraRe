# Struktur Tabel — Endorsement Life

Acuan bentuk tabel untuk aplikasi Go. Dibuat 2026-09-18 atas perintah work owner.
Presisi fisik adalah `[data DBA]` dan TIDAK ditetapkan di sini.
Berkas ini menggambarkan BENTUK, bukan alasan — alasannya ada di `spec.md` dan `issues/`.

⚠️ **Berkas ini adalah acuan TUNGGAL nama kolom. Seluruh nama kolom snake_case**
`[keputusan work owner 2026-09-18]`. Ejaan yang muncul di `spec.md`, `issues/`, dan `revisi-*.md`
adalah **ejaan korpus** — itu **bukti asal kolom**, bukan nama kolom.

---

## Endorsement tidak punya tabel sendiri

`[terverifikasi]` **New Business dan Endorsement memakai TABEL YANG SAMA.** Endorsement **tidak
menambah satu tabel pun** — yang membedakan adalah **BARIS**, bukan tabel. Sebuah endorsement adalah
**versi baru** dari polis yang sama: satu baris `T_PREMIUM_LIST` tambahan ber-`PROD_KE` lebih besar,
beserta anak-anaknya. Seluruh versi hidup berdampingan.

Buktinya: **sembilan rule SQL** di bawah isinya **sama persis, byte demi byte**, di
`PremiumList Life/RDBList/` dan `Endorsement Life/RDBList/`:

| # | Rule | Yang disentuhnya |
| --- | --- | --- |
| 1 | `SaveMasterLPDet` | `POOLDATA.M_LIFE_PREMIUM_DETAIL` — **80 kolom**, identik di kedua modul |
| 2 | `InsertPLSummary` | `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY` |
| 3 | `GetProductDtlPL` | baca produk |
| 4 | `GetRateLifePM` | baca rate |
| 5 | `InsertDataUploadLife` | `POOLDATA.M_TEMPUPLOADLIFE` |
| 6 | `DeleteTempUploadDataLife` | `POOLDATA.M_TEMPUPLOADLIFE` |
| 7 | `SelectComm_SQL` | baca komisi |
| 8 | `CekDoubleInsured` | pemeriksaan peserta ganda |
| 9 | `GetPolicyNoByCaseId` | baca nomor polis |

⛔ **Bentuk tabelnya TIDAK diulang di sini.** Sumber tunggalnya adalah
**`.scratch/premiumlist-life/STRUKTUR-TABEL-PREMIUMLIST-LIFE.md`** — seluruh tujuh tabel, kolom,
index, pohon relasi, dan tabel relasi ada di sana. Berkas ini hanya menambahkan apa yang **khas
EDM**.

Satu-satunya rule yang **berbeda** antara kedua modul adalah `SaveLifeinProduction_SQL`: jalur NB
menulis **36 kolom** ke `POOLDATA.LIFEINPRODUCTION`, jalur EDM menulis **27 kolom** dan **menambah
`NO_ENDORS`**. Gabungannya **37 kolom** — semuanya sudah tercatat di berkas PremiumList.

---

## Kolom EDM pada `T_PREMIUM_LIST`

Kolom di bawah **seluruhnya nullable** dan **kosong pada baris new business**. Daftar ini adalah
irisan ber-`Dipakai = EDM saja` dari tabel `T_PREMIUM_LIST` di berkas PremiumList — bukan tabel baru.

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `NO_ENDORS` | teks | ya | | **EDM saja** | korpus `NOENDORS` — `Endorsement Life/RDBList/SaveLifeinProduction_SQL.xml` |
| `EDM_TYPE` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement — `1`=Perubahan Data, `3`=Batal |
| `OLD_POLICY_NO` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `EDM_DATE` | DATE | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `EDM_NOTE` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `PREMI_PROPOSED` | angka desimal | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `UANG_PERTANGGUNGAN` | angka desimal | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `SUM_INSURED` | angka desimal | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `JENIS_PRODUK` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `SISTEM_REASURANSI` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `STATUSS` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `STATUS_UPDATE` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `STATUS_SERVICE` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `START_DATE` | DATE | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `END_DATE` | DATE | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `PL_NUMBER_EDM` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement |
| `PROD_KE` | bilangan bulat | ya | | **EDM saja** | keputusan tiket 00 Endorsement — versi berjalan = `PRODKE` terbesar |
| `EDM_STATUS` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement — `Old`/`New`/`Delete`/`Batal` |
| `STATUS_OLD` | teks | ya | | **EDM saja** | keputusan tiket 00 Endorsement |

⚠️ **`TYPE_CEDING` tidak ada di daftar ini** — ia **kolom bersama (NB + EDM)**, bersumber korpus.
Tiket 00 Endorsement hanya menetapkan **domain nilainya**. Lihat
`STRUKTUR-TABEL-PREMIUMLIST-LIFE.md`.

⚠️ `EDM_TYPE` hanya menerima **`1`** (Perubahan Data) dan **`3`** (Batal). Nilai `2` **memang tidak
ada** — bukan kolom yang terlewat.

⚠️ **Bentuk `NO_ENDORS`** `[terverifikasi]` — keterangan, **bukan** kolom baru:

```
NO_ENDORS = <NO_POLIS> / <PROD_KE dua digit>

rantainya:  GetProdKeOldData_SQL      -> Local.Prodke
            GenerateNoEDM_Life        -> CARI4 = Prodke + 1
                                         CARI14 = CARI4 dipad nol jadi 2 digit
            Generate_NoEndorsmentLife -> NO_POLIS || '/' || CARI14
```

---

## `PARENT_ID` pada `T_PREMIUM_LIST_DETAIL`

`PARENT_ID` adalah **FK self-reference ke peserta versi sebelumnya** — ia menjawab "baris peserta
ini menggantikan baris yang mana pada versi lama".

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `PARENT_ID` | teks | ya | FK | **EDM saja** | keputusan tiket 00 Endorsement — self-reference → `T_PREMIUM_LIST_DETAIL.ID`, ber-index |

`NULL` untuk **new business** dan untuk **peserta baru** di endorsement.

⚠️ **Hanya tabel peserta yang punya `PARENT_ID`.** `T_PREMIUM_LIST_SPREADING`,
`T_PREMIUM_LIST_SPREADING_RETRO`, dan `T_PREMIUM_LIST_SUMMARY` **TIDAK**.

---

## Bagaimana baris EDM dibedakan dari baris NB

| Tabel | Baris NB | Baris EDM |
| --- | --- | --- |
| `T_PREMIUM_LIST` | `PROD_KE` = 1, `NO_ENDORS` kosong | `PROD_KE` = 2, 3, … · `NO_ENDORS` terisi |
| `T_PREMIUM_LIST_DETAIL` | `PARENT_ID` NULL | `PARENT_ID` menunjuk versi lama |
| `_SPREADING` / `_RETRO` | — | ikut tersalin di bawah peserta versi baru, **tanpa** `PARENT_ID` |
| `T_PREMIUM_LIST_SUMMARY` | — | **tidak disalin** antar versi |
| `T_WORK_POLIS` | — | baris NB dan baris EDM **sejajar**, tidak saling menunjuk |

⛔ **TIDAK ADA tabel baru untuk endorsement. Nol.**

---

## Catatan — belum ditetapkan, TIDAK menghambat berkas ini

- `[terverifikasi]` `COMMISION` (satu S) dan `COMMISSION` (dua S) **SAMA-SAMA ejaan korpus** — 369
  berkas lawan 83 berkas. Bukan salah ketik artefak. `T_PREMIUM_LIST_SPREADING_RETRO` memakai
  `COMMISION`, `T_PREMIUM_LIST_SUMMARY` memakai `COMMISSION`, keduanya mengikuti sumbernya
  masing-masing. **JANGAN diseragamkan tanpa keputusan work owner.**
- `[terbuka]` `POLIVYHOLDER` `[terverifikasi]` ejaan korpus — `InsertDataUploadLife.xml` di **kedua** modul. Ditiru atau diperbaiki belum diputuskan
- `[terbuka]` `STATUSS` dua huruf S di akhir, tetapi `[terverifikasi]` **itu ejaan korpus** — 73 berkas XML. Bukan salah ketik artefak. Perbaikan belum diputuskan
- `[data DBA]` daftar kolom `M_LIFE_PREMIUM_SUMMARY` + isi procedure `PEGA_M_LIFE_PREMIUM_SUMMARY`
- `[data DBA]` daftar kolom `DOCUMENT_POLIS`
- `[data DBA]` presisi fisik seluruh kolom
- `[terbuka]` delapan kolom produk/layer (`PRODUCT_NAME_ID`, `PRODUCT_NAME`, `LAYER_1..4`, `ANNUITY_INTEREST`, `PREMIUM_REFUND_FACTOR`) tidak ditulis jalur EDM — kosong di baris versi EDM
- `[terbuka]` rantai versi dicari lewat (`NO_POLIS`, `PROD_KE`), belum ada kolom penunjuk versi sebelumnya
- `[terbuka]` polis lama hasil migrasi belum punya baris `T_WORK_POLIS`, padahal PK-nya diambil dari sana
- `[terbuka]` format identitas kerja `EDMLF-<n>`
- `[terbuka]` kolom `T_PREMIUM_LIST_SPREADING`, `_SPREADING_RETRO`, dan `T_VIEW_SUGGEST` **sudah ditetapkan** di `spec.md` §12 PremiumList; yang belum ditetapkan hanya presisi dan nullability per kolom
- `[terbuka]` `DESCRIPTION` dan `WPC` juga tidak ditulis jalur EDM, di luar delapan kolom produk/layer
- `[terbuka]` `SUM_INSURED`, `PL_NUMBER_EDM`, `EDM_STATUS`, `STATUS_OLD` ada di header **dan** di peserta
- `[terbuka]` `PL_NUMBER_EDM`, `EDM_STATUS`, `STATUS_OLD`, `STATUS` pada peserta ditulis rule yang identik di kedua modul, tetapi hanya terisi pada jalur EDM
- `[terbuka]` `ON DELETE` untuk `PARENT_ID` belum ditetapkan
- `[terbuka]` `EDM_TYPE` menerima `1` dan `3`; siapa yang menolak nilai lain — basis data atau lapisan layanan — belum ditetapkan
