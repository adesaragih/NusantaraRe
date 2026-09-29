# 05a: Simpan premium summary — konversi uang di satu batas, dan urutan procedure yang mengikat

**Status:** sebagian — satu transaksi polis penuh + injeksi kegagalan, `M_LIFE_PREMIUM_SUMMARY` (OQ-PL-09 — **terhalang DBA**, daftar serah terima), dan uji pulang-pergi Oracle belum ada; kosong = 0 di warisan (**OQ-PL-10 ditutup** GILIRAN-17)

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
      transaksi Go**, dan transaksi itu **commit sebelum** procedure lain dipanggil. *(AC 20 spec)* — belum: kedua procedure tidak dipanggil (penomoran lewat tabel penghitung; `M_LIFE_PREMIUM_SUMMARY` menunggu OQ-PL-09)
- [ ] ~~`INSERTJSONPOLISLIFE` dan `INSERTJSONOFFERLIFE` **tidak** dipanggil dari dalam transaksi
      itu.~~ ⚠️ **TIDAK BERLAKU 2026-09-16** — keduanya **dibuang**; lihat AC pengganti di bawah. — belum: dicoret, tidak berlaku (2026-09-16)

### Penyimpanan relasional ⚠️ BARU 2026-09-16 — spec §6, §12

- [ ] ⚠️ **Satu polis ditulis dalam SATU transaksi**: header, seluruh baris rekap mata uang, seluruh
      peserta, seluruh spreading, seluruh spreading retro, dan seluruh riwayat penawaran — lalu
      **commit sekali**. Kegagalan di tingkat mana pun **membatalkan seluruhnya**. Dibuktikan dengan
      menyuntikkan kegagalan pada baris peserta ke-N dan memastikan **tidak ada** polis tersimpan.
      *(AC 35 spec; penyimpangan sadar 1)* — belum: `simpanDalam` memuat nomor + rekap + salinan warisan saja; header, peserta, spreading, dan riwayat penawaran tidak ikut, dan uji injeksi kegagalan belum ada
- [ ] Penomoran (`PROC_GENERATE_SEQUENCE_NUMBER`) berada **di dalam** transaksi itu: tidak pernah ada
      nomor tanpa polis, tidak pernah ada polis tanpa nomor. *(AC 36 spec)* — belum: di dalam transaksi `simpanDalam` (uji `TestSimpanDalamUrutanTerkunci`), tetapi `POST …/nomor` (`Terbitkan`) tetap menerbitkan nomor tanpa rekap
- [ ] ⚠️ **Tidak ada procedure JSON yang dipanggil.** Test yang menemukan pemanggilan
      `INSERTJSONPOLISLIFE`, `INSERTJSONOFFERLIFE`, atau padanan `@ASM.GetPageJSONString()`
      **gagal**. *(AC 32, 34 spec; penyimpangan sadar 1)* — belum: nol pemanggilan procedure JSON di kode, tetapi uji penjaga yang gagal bila pemanggilan muncul belum ada
- [x] ⚠️ **Rekap uang per mata uang tersimpan** di tabel rekap — bukan dihitung lalu dibuang. Polis
      bermata uang ganda menghasilkan **satu baris rekap per mata uang**, masing-masing dengan nilai
      uangnya. *(AC 37 spec; penyimpangan sadar 2)* — bukti: `repository/polis_summary.go:GantiRekap`; uji `TestSatuBarisRekapPerMataUang`, `TestNilaiSisipRekapSejajarDanTepat`
- [ ] ⚠️ `M_LIFE_PREMIUM_SUMMARY`, `M_LIFE_PREMIUM_DETAIL`, dan `LIFEINPRODUCTION` **tidak ditulis**.
      Test yang menemukan tulisan ke ketiganya **gagal**. *(AC 33 spec)* — belum: pl2 membalik — `M_LIFE_PREMIUM_DETAIL` ditulis (`repository/polis_warisan.go:Ganti`); `M_LIFE_PREMIUM_SUMMARY`/`LIFEINPRODUCTION` tidak ditulis, tetapi uji penjaganya belum ada
- [ ] ⚠️ Bila jalur warisan `SaveMasterLPDet` masih dipakai selama transisi, ia berada **di luar**
      transaksi polis dan **dapat diulang** — ia satu-satunya titik potong yang tersisa.
      *(spec §6; `[data DBA]` `COMMIT` di dalam procedure)*
      *(AC 21 spec)* — belum: pl2 menaruh salinan warisan DI DALAM transaksi `simpanDalam`; dapat diulang (hapus lalu sisip, uji `TestHapusWarisanDikurungNomorDanWork`)
- [ ] Kegagalan **sebelum** commit summary tidak meninggalkan nomor maupun rekam separuh: test
      menyuntikkan kegagalan di `PEGA_M_LIFE_PREMIUM_SUMMARY` dan memastikan **tidak ada** baris
      `M_LIFE_PREMIUM_SUMMARY` **dan** sequence tidak bergerak. *(AC 22 spec)* — belum: procedure tidak dipanggil, dan nol uji injeksi kegagalan
- [ ] Kegagalan **setelah** commit summary meninggalkan nomor + rekam summary **utuh** dan keadaan itu
      **terdeteksi** — bukan senyap. *(AC 23 spec)* — belum: efek keluar sesudah commit tercatat di outbox (uji `TestEfekBerjalanSesudahCommitBukanDiDalamnya`), tetapi `Putuskan` (jalur Confirm) membuang ringkasan `EfekKeluar`
- [x] **Ada test yang gagal bila urutan pemanggilan diubah** — urutannya bagian dari kebenaran, bukan
      kebetulan. *(AC 24 spec)* — bukti: uji `TestSimpanDalamUrutanTerkunci`, `TestSimpanSebelumTutupDalamSatuTransaksi`
- [ ] Konversi teks ↔ desimal terjadi **hanya di lapisan repository**, di satu tempat; lapisan
      services dan handlers hanya mengenal desimal. *(AC 15 spec; **ADR-0003**)* — belum: `services/polis_summary.go:keTampil` memformat desimal ke teks di lapisan services
- [ ] Ke-37 kolom terisi dari sumber yang benar, dan pemetaannya diuji kolom demi kolom — **bukan**
      lewat posisi `CARI2`…`CARI34` yang tidak bernama. — belum: `M_LIFE_PREMIUM_SUMMARY` tidak ditulis — OQ-PL-09
- [ ] Nilai uang yang dikirim dan dibaca kembali **identik**, termasuk nilai berpecahan panjang dan
      nilai negatif. Tidak ada pembulatan diam. *(AC 16 spec)* — belum: nol uji pulang-pergi terhadap Oracle; rekap sengaja dibulatkan empat angka (DT)
- [ ] Baik `PL_NUMBER` maupun `PL_NUMBER_EDM` tersimpan pada rekam summary yang sama sebagai **dua
      nilai terpisah**; jalur new business mengisi yang pertama, endorsement yang kedua.
      *(AC 13 spec)* — belum: rekam summary warisan tidak ditulis (OQ-PL-09); `T_PREMIUM_LIST_SUMMARY` tanpa kolom `PL_NUMBER`
- [ ] Keluaran galat procedure (`ERRMSG`, `STSSAVE`) **diperiksa**; galat yang dilaporkan procedure
      tidak boleh diabaikan sehingga transaksi tampak berhasil. — belum: tidak berlaku — procedure tidak dipanggil (pl1)

## Blocker

**Tidak ada pemblokir.** Terkait tetapi tidak memblokir: **OQ-001 sisa** (DDL fisik — memblokir
tiket **09**, bukan tiket ini; kolomnya sudah terbaca dari body procedure).

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```

## Pembacaan XML 28-09-2026 — peta 33 parameter, dan TIGA yang belum tertentukan

`[terverifikasi]` Sumber tiap parameter `PEGA_M_LIFE_PREMIUM_SUMMARY` ditemukan di
**`Activity/InsertJsonPolisLife_Act.xml`** — bukan di `SubmitPremiumList_Act` maupun
`SavePremiumList_Act`, yang keduanya hanya mengisi sebagian. Perintah audit:

```
grep -o "<PropertiesName>TempInputData\.CARI[0-9]*</PropertiesName>\s*<PropertiesValue>[^<]*" \
  Activity/InsertJsonPolisLife_Act.xml
```

| Param | Sumber | | Param | Sumber |
| --- | --- | --- | --- | --- |
| `CARI2` | `.CURRENCY` | | `CARI19` | `.BROKERAGE_FEE_RETRO` |
| `CARI3` | `.PREMIUM` ⚠️ | | `CARI20` | `.NET_PREMIUM_RETRO` |
| `CARI4` | `.COMMISSION` ⚠️ | | `CARI21` | `.GROSS_PREMIUM_REFUND_RETRO` |
| `CARI5` | `.NET_PREMIUM_REFUND` | | `CARI22` | `.DISCOUNT_PREMIUM_REFUND_RETRO` |
| `CARI6` | `.GROSS_PREMIUM_REFUND` | | `CARI23` | `.OVR_COMM_REFUND_RETRO` |
| `CARI7` | `.COMM_REFUND` | | `CARI24` | `.BROKERAGE_FEE_REFUND_RETRO` |
| `CARI8` | `.BROKERAGE_FEE_REFUND` | | `CARI25` | `.NET_PREMIUM_REFUND_RETRO` |
| `CARI9` | `.OVR_COMM_REFUND` | | `CARI26` | `.CLAIM` |
| `CARI10` | `.TAX_REFUND` | | `CARI27` | `.BALANCE` ⚠️ |
| `CARI11` | `.BROKERAGE_FEE` | | `CARI28` | `.CLAIM_AMOUNT` |
| `CARI12` | `.OVR_COMM` | | `CARI29` | `.RI_ADMIN_FEE_RETRO` |
| `CARI13` | `.TAX` | | `CARI30` | `.RI_ADMIN_FEE_REFUND_RETRO` |
| `CARI14` | `.PROF_COMM` | | `CARI31` | `.RI_ADMIN_FEE_REFUND` |
| `CARI15` | `.SHARE_RETRO` | | `CARI32` | `.DEDUCTION_REFUND` |
| `CARI16` | `.GROSS_PREMIUM_RETRO` | | `CARI33` | `.RI_ADMIN_FEE` |
| `CARI17` | `.DISCOUNT_PREMIUM_RETRO` | | `CARI34` | `.DEDUCTION` |
| `CARI18` | `.OVR_COMM_RETRO` | | | |

Tiga puluh dua kolom uang plus `CURRENCY` — **33**, sesuai tiket. Dua puluh sembilan di antaranya
punya kolom bernama **sama** di `T_PREMIUM_LIST_DETAIL` (052), jadi rekapnya `SUM(...) GROUP BY
CURRENCY` langsung.

### ⛔ OQ-PL-08 — TIGA parameter belum tertentukan sumbernya

| Param | Medan summary | Kenapa belum tertentukan |
| --- | --- | --- |
| `CARI3` | `.PREMIUM` | detail punya **`GROSS_PREMIUM`** *dan* **`NET_PREMIUM`**; korpus tidak menyatakan yang mana — dan keduanya angka yang sah, sehingga salah pilih tidak akan pernah berbunyi |
| `CARI4` | `.COMMISSION` | detail bernama **`COMM`**; kemiripan nama BUKAN bukti pemetaan, dan `COMM` berdampingan dengan `PROF_COMM`/`OVR_COMM` yang juga komisi |
| `CARI27` | `.BALANCE` | **nol** kolom detail bernama itu. Ia tampaknya angka turunan *(posisi bersih)*, tetapi rumusnya tidak ada di korpus PremiumList Life |

`[terbuka — pemilik kerja]` Ketiganya **kolom uang**, dan uang yang dipetakan dengan tebakan akan
tersimpan, terbaca, dan terlaporkan sebagai angka yang sah. **Executor tidak menebaknya.** Tiket 05a
menunggu ketiga pemetaan ini diputuskan; dua puluh sembilan sisanya sudah siap dibangun.

⚠️ Dicari di: `SubmitPremiumList_Act`, `SavePremiumList_Act`, `InsertJsonPolisLife_Act`,
`ShowLifePremiumSummary.xml`, dan seluruh `Activity/*.xml` lewat pencarian penetapan
`PropertiesName` yang memuat `.PREMIUM`/`.BALANCE`/`.COMMISSION`. Yang ditemukan hanya penetapan
`.CLAIM` pada baris DETAIL, bukan pada summary. Rumusnya kemungkinan di Declare Expression atau
Data Transform yang **tidak ikut** terekspor ke korpus ini.

## OQ-PL-08 DITUTUP 28-09-2026 — dan bacaan brief GILIRAN-8 §1 DIRALAT

`DataTransform/AppendCurrencySummary_DT.xml` dibaca **utuh** (5.992 baris, 121 penetapan).
Perintah audit — urutan aksi apa adanya, bukan pasangan yang diratakan:

```
grep -o "<pyActionName>[^<]*\|<pyPropertiesName>[^<]*\|<pyPropertiesValue>[^<]*" \
  DataTransform/AppendCurrencySummary_DT.xml
```

### ⚠️ RALAT — cabangnya BUKAN "per jenis baris", dan BUKAN dijumlahkan

GILIRAN-8 §1 membaca `.PREMIUM` sebagai *"Σ `.GROSS_PREMIUM` + Σ `.GROSS_PREMIUM_REFUND` +
Σ `.GROSS_PREMIUM_RETRO` + Σ `.GROSS_PREMIUM_REFUND_RETRO` (per cabang jenis baris)"*.

**Struktur sesungguhnya — empat `WHEN` yang SALING MENIADAKAN atas `pyWorkPage.Type`:**

```
WHEN Param.Currency == .CURRENCY          <- pengelompokan per mata uang
  WHEN pyWorkPage.Type == "QR"  ->  Param.Premium += .GROSS_PREMIUM
  WHEN pyWorkPage.Type == "QP"  ->  Param.Premium += .GROSS_PREMIUM_REFUND
  WHEN pyWorkPage.Type == "TP"  ->  Param.Premium += .GROSS_PREMIUM_RETRO
  WHEN pyWorkPage.Type == "TR"  ->  Param.Premium += .GROSS_PREMIUM_REFUND_RETRO
```

Satu polis punya **SATU** `Type`, jadi **tepat satu** cabang menyala. `.PREMIUM` adalah Σ **satu**
kolom yang dipilih `Type` — **bukan** Σ empat kolom.

⛔ **Kenapa selisih ini mahal.** Menjumlahkan keempatnya akan **melipatgandakan** `PREMIUM` untuk
polis yang barisnya memuat nilai di lebih dari satu kolom itu — dan hasilnya tetap angka yang sah,
tersimpan, lalu dilaporkan. Tidak satu pun galat akan menunjukkannya.

⚠️ Bacaan brief tampaknya lahir dari mencocokkan `pyPropertiesName`/`pyPropertiesValue` secara
**berpasangan**, yang meratakan struktur `WHEN`-nya: di ekspor Pega, kondisi sebuah `WHEN` menempati
medan `pyPropertiesName`, sehingga pasangan yang diratakan terbaca seperti empat penetapan berurutan.
Yang benar dibaca dari **urutan aksi** (`pyActionName` `WHEN`/`SET`), bukan dari pasangannya.

### Ketiga parameter, VERBATIM

| Param | Rumus |
| --- | --- |
| `.PREMIUM` | Σ satu kolom menurut `Type`: `QR`→`GROSS_PREMIUM` · `QP`→`GROSS_PREMIUM_REFUND` · `TP`→`GROSS_PREMIUM_RETRO` · `TR`→`GROSS_PREMIUM_REFUND_RETRO` |
| `.COMMISSION` | `Σ .COMM` saja — **tanpa** cabang `Type`. `PROF_COMM` dan `OVR_COMM` parameter tersendiri |
| `.BALANCE` | Σ per `Type`, lalu `@divide(Param.Balance,1,4)` |

**`.BALANCE` per cabang, disalin apa adanya:**

```
QR  Balance + (GROSS_PREMIUM - DEDUCTION
               - (RI_ADMIN_FEE + BROKERAGE_FEE + TAX + PROF_COMM + CLAIM))
QP  Balance + (GROSS_PREMIUM_REFUND + CLAIM_AMOUNT
               - (DEDUCTION_REFUND + BROKERAGE_FEE_REFUND + RI_ADMIN_FEE_REFUND
                  + TAX + PROF_COMM + CLAIM))
TP  Balance + (GROSS_PREMIUM_RETRO - DISCOUNT_PREMIUM_RETRO
               - RI_ADMIN_FEE_RETRO + BROKERAGE_FEE_RETRO)
TR  Balance + (GROSS_PREMIUM_REFUND_RETRO - DISCOUNT_PREMIUM_REFUND_RETRO
               - RI_ADMIN_FEE_REFUND_RETRO + BROKERAGE_FEE_REFUND_RETRO)
```

⚠️ **Keanehan warisan, VERBATIM, tidak "diperbaiki"** — dan brief benar menandainya:

1. `BROKERAGE_FEE_RETRO` **DITAMBAH** pada cabang `TP`/`TR`, sedangkan `BROKERAGE_FEE` **DIKURANGI**
   pada cabang `QR`. Biaya yang menambah saldo di satu cabang dan mengurangi di cabang lain.
2. Cabang `QP` mengurangi `TAX`, `PROF_COMM`, dan `CLAIM` — **bukan** padanan `*_REFUND`-nya,
   padahal ketiga biaya lain di cabang itu memakai `*_REFUND`.
3. `@divide(Param.Balance,1,4)` di akhir: pembagian dengan **1** — yaitu pembulatan ke **4** angka
   desimal, ditulis sebagai pembagian.

Ketiganya disalin apa adanya. Menormalkan tandanya mengubah angka uang yang sudah beredar.

### Bentuk yang akan dibangun

`models` merakit rumusnya **per baris peserta** (murni, dapat diuji tanpa Oracle), lalu dijumlah per
`CURRENCY`. Dua puluh sembilan kolom sisanya `SUM(...) GROUP BY CURRENCY` langsung — pengelompokan
per mata uang itu sendiri VERBATIM dari `WHEN Param.Currency == .CURRENCY`.

⛔ Pembulatan empat angka **hanya di akhir**, sesudah penjumlahan — bukan per baris. Membulatkan per
baris lalu menjumlah menghasilkan angka yang berbeda dari menjumlah lalu membulatkan, dan selisihnya
tumbuh bersama cacah peserta.

## Implementasi bagian 2 — 28-09-2026 (giliran 10): pembaca, penyimpan, rute, layar

### Pohon XML yang dibaca sebelum kode

| Rule | Yang diambil |
| --- | --- |
| `DataTransform/AppendCurrencySummary_DT.xml` | langkah **2.38** perulangan baris; **2.39–2.74** ke-36 keluaran (`pyPropertyStepId` dibaca bersama `pyActionName`) |
| `Activity/SavePremiumList_Act.xml` | langkah **12** `Apply-DataTransform AppendCurrencySummary_DT` — rekap dihitung untuk layar |
| `Activity/SubmitPremiumList_Act.xml` | 17 langkah: 1 `Page-Remove` · 2–3 tanggal closing · 4–5 reset · 6–9 PL Number per `Type` · 10–13 kode produksi + urut · **14** `Set COB & PL Number` · **15** pl7 · **16** `pyMemo=1`, `IsJsonPolis=0` · **17** `Obj-Save` |
| `RDBList/SaveMasterLPDet.xml` | `INSERT INTO POOLDATA.M_LIFE_PREMIUM_DETAIL` — 80 kolom (b87–b167), `COMMIT;` b252 |
| `Activity/InsertLifePremiumDetail_act.xml` | 50 penetapan `TempInputDetail.CARIn` — sumber tiap nilai `SaveMasterLPDet` |
| `Section/ShowLifePremiumSummary.xml` | empat layout `pyContainerVisibleWhen` `.Type = 'QR'/'QP'/'TP'/'TR'`, judul + properti kolom; tombol `Submit` |

### ⚠️ RALAT bagian 1 — SELURUH keluaran dibulatkan, bukan hanya `BALANCE`

Bagian 1 (`021d0f3`) membulatkan `.BALANCE` saja, dan `TestPembulatanHanyaDiAkhir` **menagih**
`PREMIUM` tak dibulatkan (`0.00015`). `[terverifikasi]` DT langkah **2.39**
`.PREMIUM = @divide(Param.Premium,1,4)`, **2.40** `.COMMISSION = @divide(Param.Commission,1,4)`,
**2.41–2.73** tiga puluh tiga kolom, **2.74** `.BALANCE`. Perintah audit:

```
grep -o "<pyPropertyStepId>[^<]*\|<pyActionName>[^<]*\|<pyPropertiesName>[^<]*\|<pyPropertiesValue>[^<]*" \
  DataTransform/AppendCurrencySummary_DT.xml
```

Cacah `<pyPropertiesValue>@divide(…,1,4)` = **36**; cacah target `<pyPropertiesName>.X</pyPropertiesName>`
berhuruf besar = **36** — dua cara, sepakat. Uji lama diralat; `TestSeluruhKeluaranDibulatkanEmpatAngka` lahir.

### ⚠️ RALAT bagian 1 — daftar kolom yang dijumlah BUKAN `KolomUangUnggah`

Bagian 1 menjumlah ke-32 kolom unggahan. DT menjumlah **33 kolom lain** (`models.KolomJumlahSummary`):
tanpa `NET_PREMIUM`, `GROSS_PREMIUM`, `SUM_INSURED`, `EM_PERCENT`, `FLEET_DISCOUNT`, `COMM`,
`SHARE_NUSANTARA_RE`; dengan `DEDUCTION`, `DEDUCTION_REFUND`, `RI_ADMIN_FEE*` (4), `SUM_AT_RISK_GROSS`,
`SHARE_NUSANTARA_RE_GROSS`. `TestKolomJumlahSummaryVERBATIMDariDT` membaca DT langsung dua cara —
**cara A** target keluaran, **cara B** penjumlahan `param.X + .X` — dan menagih kesamaan dua arah.
⚠️ Instrumennya **diuji dulu**: versi pertama cara A memakai jendela `name…value` yang tidak pernah
cocok (tag lain menyela di XML mentah), dan cara B memungut cabang `Param.Premium + .GROSS_PREMIUM`;
keduanya diperbaiki sebelum dipercaya. Himpunan 33 kolom = seluruh kolom uang migrasi 055 di luar
`BALANCE`/`PREMIUM`/`COMMISSION` (`TestKolomRekapSamaDenganMigrasi055`, dua arah).

### Pilihan pembaca — baca lalu hitung dengan rumus murni

Brief memberi dua jalan (SUM/GROUP BY di SQL, atau baca lalu hitung). Dipilih **baca lalu hitung**:
cabang `Type` dan tiga keanehan tanda `BALANCE` hidup satu kali di `models`, dikunci literal
**975 / 978 / 897 / 1789** (`TestBalanceKeempatCabangDariLiteral`, kini `975.0000` dst.). SQL `CASE`
akan menjadi salinan kedua rumus uang yang tidak diuji literal itu. Kesetaraan pembaca ↔ rumus:
`sqlBarisUangPolis` memilih **tepat** `models.KolomBacaSummary()` (33 + `COMM` + kolom `PREMIUM`/`BALANCE`
keempat cabang), seluruhnya ada di 052 (`TestKolomBacaSummaryAdaDiMigrasi052`), seluruhnya lewat
`TO_CHAR` ber-NLS; `TestNilaiSisipRekapSejajarDanTepat` membawa literal 975 dari rumus sampai ke
argumen SQL (`"975.0000"`).

### Satu transaksi — `SummaryPremiumList.Submit`

`POST /api/polis-life/{id}/summary`: gerbang kasus tertutup (`T_WORK_POLIS`) → **satu**
`DalamTransaksi`: `terbitkanDalam` (langkah 10–14; diekstrak dari `Terbitkan` tiket 03 supaya gerbang
lahir-sekali satu fungsi untuk dua jalur) → rekap murni → `GantiRekap` (hapus lalu sisip
`T_PREMIUM_LIST_SUMMARY`) → `SumberWarisan` → `PesertaWarisan.Ganti` (pl2) → commit. Urutan dikunci
`TestSubmitSummaryUrutanTerkunci` (dan tepat **satu** `DalamTransaksi`). Nol `COMMIT` di teks SQL
(`TestNolCommitDiQueryRekapDanWarisan`). `GET` yang sama menghitung **tanpa** menyimpan.

### ⚠️ RALAT AC 33 — `M_LIFE_PREMIUM_DETAIL` DITULIS (pl2), `M_LIFE_PREMIUM_SUMMARY` TIDAK (OQ-PL-09)

AC 33 (2026-09-16) melarang penulisan `M_LIFE_PREMIUM_SUMMARY`/`M_LIFE_PREMIUM_DETAIL`.
**pl2** (brief 3-PREMIUMLIST, 28-09-2026, lebih baru) membaliknya: keduanya ditulis dalam transaksi
yang sama, sebab Claim Life membaca `M_LIFE_PREMIUM_DETAIL` (`GET /api/peserta-life`). Yang dikerjakan:

- **`M_LIFE_PREMIUM_DETAIL` — DITULIS.** Pemetaan 80 kolom `[terverifikasi]` dari `SaveMasterLPDet` +
  50 `CARIn` `InsertLifePremiumDetail_act` (`repository/polis_warisan.go`, `kolomPesertaWarisan`);
  `TestKolomWarisanVERBATIMDariSaveMasterLPDet` membaca korpus dua cara (daftar kolom = 80, butir
  `VALUES` = 80, `To_date` = 11) dan menagih urutannya posisi demi posisi. Idempoten: `DELETE … WHERE
  PL_NUMBER = :1 AND IDPEGA = :2` — dikurung **work**, supaya baris endorsemen yang memuat `PL_NUMBER`
  sama tidak ikut terhapus. `STATUSOLD = '0'` dan `STATUS = CARI48` VERBATIM. `IDPEGA` (`pzInsKey`)
  diisi **pengenal work** (`polisID` = `T_WORK_POLIS.ID`), **bukan** `T_PREMIUM_LIST.ID_PEGA`: kolom
  header itu belum punya satu pun penulis di repo ini, sehingga membacanya membuat setiap submit gagal.
  ⚠️ Berkas penulisnya **tanpa satu pun kueri baca** — penjaga 66,8 juta baris tetap tajam; bahan
  salinan dibaca dari tabel kami di `polis_summary.go` (`TestPenulisWarisanTanpaKueriBaca`).
  `TestPenyaringPesertaHanyaSatuTempat` kini menghitung **3** berkas yang menyebut `EDMSTATUS`
  (penulis ini menyebutnya di daftar kolom `INSERT`, bukan sebagai penyaring).
- **`M_LIFE_PREMIUM_SUMMARY` — TIDAK DITULIS. ⛔ OQ-PL-09 `[terbuka — DBA]`.** Isi procedure
  `PEGA_M_LIFE_PREMIUM_SUMMARY` dan daftar kolom tabelnya **tidak ada di korpus** (`[data DBA]`,
  STRUKTUR-ENDORSEMENT baris 131). `InsertPLSummary` memanggilnya **posisional** — 37 argumen masuk +
  2 keluar (STRUKTUR-PREMIUMLIST §`T_PREMIUM_LIST_SUMMARY`) — jadi nama kolomnya hanya ada di dalam
  procedure. Tiket ini sendiri (baris 56–57) mencatat PK dari `M_LIFE_PREMIUM_SUMMARY_SEQ` **dan**
  "37 kolom" dengan `ID` di antaranya: 1 `ID` + 37 argumen = 38, sehingga satu argumen tidak
  berkolom atau tidak bernama sama — dan yang mana tidak dapat dibaca. Memetakan dengan tebakan nama
  menyimpan uang di kolom yang mungkin salah. ⚠️ Tiket 09 (baris 51) dan tiket EDM 12 (baris 62)
  menyebut body procedure **diserahkan DBA 2026-09-15** — tetapi pemetaan argumen → kolomnya **tidak
  tercatat** di repo mana pun. **Yang dibutuhkan:** salinan body yang diserahkan itu dicatat (atau
  `ALL_SOURCE`/`ALL_TAB_COLUMNS`); sesudahnya penulisnya satu fungsi di `polis_warisan.go`.
  Dua AC ikut **terbuka** karenanya: *"Ke-37 kolom terisi…"* dan *"`PL_NUMBER` maupun `PL_NUMBER_EDM`
  tersimpan pada rekam summary"*.

### AC — keadaan

| AC | Keadaan |
| --- | --- |
| satu polis satu transaksi; kegagalan membatalkan seluruhnya | ✅ struktural (satu `DalamTransaksi`, dikunci statik) · ⚠️ **injeksi kegagalan peserta ke-N terhadap Oracle BELUM dibuktikan** — skema uji belum memasang 050–056 dan mesin ini tanpa Oracle |
| penomoran di dalam transaksi | ✅ |
| nol procedure JSON | ✅ (tidak ada yang dipanggil) |
| rekap per mata uang tersimpan, satu baris per mata uang | ✅ `GantiRekap`; `PengenalRekap` deterministik |
| urutan dikunci uji | ✅ `TestSubmitSummaryUrutanTerkunci` |
| konversi teks ↔ desimal hanya di repository | ✅ |
| uang identik, tanpa pembulatan diam | ✅ teks `FormatDecimal`; pembulatan empat angka hanya yang ditulis DT |
| `ERRMSG`/`STSSAVE` diperiksa | ➖ tidak berlaku — procedure tidak dipanggil (pl1); galat Oracle dikembalikan apa adanya |
| AC 33 | ⚠️ diralat di atas (pl2) |

### Layar

`pages/premiumlist/PremiumListSummary.tsx` — tahap `Input Premium Summary`. Grid per `Type` VERBATIM
(`GRID_REKAP`, `labels.premiumlist.ts`), termasuk keanehan judul **`PREMIUM DEDUCTION` di atas
`.COMMISSION`** pada QP/TP/TR. ⚠️ `.COB` baris rekap **tidak ditetapkan rule mana pun** di korpus
PremiumList Life (grep `COB` di `Activity/`, `DataTransform/` = nol penetapan pada baris mata uang) —
sel ditandai `—` dan layar mengatakannya. Tombol `Submit` memanggil `POST …/summary`;
`InsertJsonPolisLife_Act` + `finishAssignment` = tiket **05b**.

### Langkah 15 — pl7

`[DIPUTUSKAN; veto work owner]` no-op di produksi, **tidak ditiru** (b1142/b3413/b3126). Dicatat
juga di tiket 03. Langkah 16 (`pyMemo`, `IsJsonPolis`) tidak ditiru: properti halaman kerja tanpa
kolom, dan `IsJsonPolis` milik jalur JSON yang dibuang.

### Angka

Go **530 PASS · 0 FAIL** tingkat atas (+16 dari 514; +101 sub-uji), `go vet ./...` dan
`go vet -tags db ./...` bersih, `gofmt` bersih · vitest **346** (+7) · `tsc --noEmit` bersih.
Perintah hitung: `go test -json -count=1 ./...` → cacah `Action=pass` tanpa `/` di nama uji.

### Temuan `/code-review` — diperbaiki sebelum commit, atau dicatat

Dua sumbu (Standards, Spec) dijalankan paralel atas pohon kerja lawan `0cddee0`. Spec memverifikasi
ulang dari korpus: 36 `@divide`, pemetaan 80 kolom baris demi baris, dan keempat grid — ketiganya ✅.

| Sumbu | Temuan | Tindakan |
| --- | --- | --- |
| Spec (c) | ⛔ `IDPEGA` dibaca dari `T_PREMIUM_LIST.ID_PEGA`, yang **nol penulisnya** di repo — setiap submit akan gagal | diganti pengenal work (`polisID`); `ErrSalinanTanpaIDPega` dibuang |
| Spec (a) | argumen "38 vs 37" OQ-PL-09 kurang rujukan; body DBA disebut sudah diserahkan 2026-09-15 | OQ ditulis ulang dengan rujukan baris, dan yang diminta kini salinan body itu |
| Standards 7 | cabang uji mati (`GROSS_PREMIUM` tidak ada di `KolomJumlahSummary`) | dibuang |
| Standards 2 | `Lihat` mengetik ulang pemeriksaan `siapkan`, dan tanpa gerbang tertutup tanpa alasan tertulis | `periksaDasar` bersama; alasan tanpa gerbang ditulis |
| Standards 6 | `AngkaDesimalBalance` berbohong sesudah ralat | → `AngkaDesimalRekap` |
| Standards 10 | `GRID_REKAP` tanpa nomor baris | b4018/b8761/b15196/b20690 dicatat |
| Spec | komentar `terbitkanDalam` "10-13" lawan "10-14" | diselaraskan 10–14 |

**Dicatat, tidak diubah:**

- ⛔ **OQ-PL-10 `[terbuka — pemilik kerja]` kosong = NULL atau 0 di tabel warisan.** Pega mengisi
  `CARIn = @toDecimal(.X)`, yang untuk properti kosong menghasilkan nol, jadi `M_LIFE_PREMIUM_DETAIL`
  warisan memuat `0` di tempat kami menulis `NULL` (ADR-U-0027). Claim Life membaca tabel itu —
  pilihan ini kontrak hilir, dan diputuskan dengan sadar, bukan oleh executor.
- ⚠️ **Urutan mata uang.** DT menyusun `CurrencyList` menurut kemunculan pertama di baris detail;
  pembaca kami `ORDER BY d.ID` (pengenal md5), jadi urutan baris rekap **di layar** dapat berbeda dari
  Pega. Nilai tidak terpengaruh; tabel 052 tidak punya kolom urutan unggah.
- ⚠️ **Layar sebagian.** Blok kepala `ShowLifePremiumSummary` (`.Type`, `.DateReceived`, `.SobName`,
  `.CedingCoName`, `.PolicyHolderName`, `.MarketingName`, `.BusinessName`, `.RetroName`, `.WPC`) belum
  dirender — hanya grid rekap dan `Submit`. `KepalaPolis` belum membawa medan itu.
- ⚠️ Angka ke kolom `NUMBER` warisan dikirim sebagai teks dan bergantung pada NLS sesi — pola yang
  sudah ada (`nilaiSisipPeserta`), tidak diubah di tiket ini. `jenisNilai` menyerupai enum di
  `kolompeserta.go`; disatukan bila keduanya disentuh lagi.

## Keputusan bertanggal — 29 September 2026 (GILIRAN-17 paket 4: OQ-PL-09, OQ-PL-10) `[keputusan work owner 29-09-2026 — lembar keputusan, "rekomendasi"]`

| OQ | Keputusan | Keadaan |
| --- | --- | --- |
| **PL-09** | tulis `M_LIFE_PREMIUM_SUMMARY` sesuai pl2, dalam transaksi simpan summary yang sama, kolom VERBATIM dari prosedur yang ditiru | ⛔ **terhalang — tidak ditulis, tidak ditebak.** `InsertPLSummary` memanggil `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY` secara **posisional** (37 masuk + 2 keluar), jadi nama kolomnya hanya ada di badan prosedur. Badan itu tidak tercatat di repo mana pun, dan hitungannya tidak cocok (1 `ID` + 37 argumen = 38 lawan 37 kolom). Memetakan dengan tebakan menyimpan uang di kolom yang mungkin salah. **Dipindah ke daftar serah terima DBA**: `ALL_SOURCE` `PEGA_M_LIFE_PREMIUM_SUMMARY` dan `ALL_TAB_COLUMNS` `M_LIFE_PREMIUM_SUMMARY`. Sesudah itu penulisnya satu fungsi di `polis_warisan.go`, dipanggil sesudah `GantiRekap` |
| **PL-10** | kolom uang kosong di `M_LIFE_PREMIUM_DETAIL` warisan = **0**, seperti Pega | ✅ `repository.kolomNolBilaKosongWarisan` (45 kolom) di `nilaiSalinWarisan`, satu fungsi di tepi repository warisan. Himpunannya adalah kolom `SaveMasterLPDet` yang memakai `TempInputDetail.CARIn`, dengan `CARIn = @toDecimal(.X)` di `InsertLifePremiumDetail_act` langkah 3.3.3 (b1843, hidup), termasuk `RISK` (CARI50). Kolom teks, tanggal, `PERIOD_YY/MM`, `PASSED_PERIOD`, dan `AGE`/`ENTRY_AGE`/`CURRENT_AGE` (tanpa `@toDecimal`) tetap NULL. ADR-U-0027 tetap berlaku untuk tabel `T_*`. Uji: `TestKosongWarisanJadiNolSepertiToDecimal`, dan `TestKolomNolBilaKosongWarisanDariKorpus`, yang menurunkan himpunannya ulang dari korpus dua arah (mutasi menjadi merah) |
