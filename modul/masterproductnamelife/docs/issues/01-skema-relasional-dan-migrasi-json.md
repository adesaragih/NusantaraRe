# 01: Skema relasional penuh + migrasi JSON → kolom — **PREFACTOR**

**Status:** ⭐ **aktif (02-10-2026)** — OQ-MPNL-01 dibuka ulang dan ditutup **flat** (K5 keputusan work owner 02-10-2026); rinciannya di bab *"Keputusan bertanggal 02-10-2026"* di ekor tiket ini. *(Status lama dikutip: "ditangguhkan (01-10-2026) — RALAT P1: produk tetap JSON di dua tabel lama seperti Pega, nol DDL; menunggu OQ-MPNL-01 (paket 0 `71c35b1`)")*

**Blocked by:** **CL-01** (kerangka aplikasi + seam API — scaffolding lintas konteks, tidak dibuat
di sini)

⚠️ **Ini PREFACTOR, dan ia tiket PERTAMA — bukan tiket migrasi yang biasanya terakhir.** Bentuk
skema berubah total (JSON → relasional), sehingga **tidak ada slice lain yang dapat berdiri**
sebelum bentuk barunya ada. *"Make the change easy, then make the easy change."*

## Hasil & nilai pengguna

Sebagai **tim migrasi**, saya ingin setiap atribut produk menjadi **kolom bernama** dan setiap daftar
bersarang menjadi **tabel anak**, dengan data lama pindah tanpa kehilangan satu nilai pun — sehingga
bentuk produk dapat diperiksa, dicari, dan divalidasi, bukan tersembunyi di dalam satu dokumen.
*(User story 29–32 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `migrations/` | DDL **satu tabel induk** `product_life` + **lima** tabel anak + sequence; skrip migrasi JSON → relasional (gabung dua tabel existing jadi satu) |
| `internal/models` | Bentuk produk (sisi umum + sisi inward) dan **lima** jenis baris anak |
| — | Skrip rekonsiliasi |

## Keadaan lama `[data DBA]`

| Tabel | Bentuk sekarang |
| --- | --- |
| `POOLDATA.M_PRODUCTINWARD_LIFE` | **hanya** `ID` + `JSONDATA` |
| `POOLDATA.M_PRODUCT_LIFE` | `JSONDATA` (constraint `IS JSON`) + empat kolom hasil flatten: `RIRISKID`, `RIRISK`, `PRODUCTNAME`, `BEGIN_DATE` |
| Kunci | **tanpa PK** — hanya index |

⚠️ **Penyimpangan sadar 2 — skema relasional penuh.** `[keputusan work owner]` Seluruh atribut
menjadi kolom bernama; seluruh daftar bersarang menjadi tabel anak; data JSON dimigrasikan.

## Bentuk baru

```
product_life   (ID + ~20 field umum + 40 field inward, field duplikat → satu kolom;
                TERMASUK LIENCLAUSE skalar)
  ├─ product_life_comment          ← CommentList
  ├─ product_life_document_claim   ← DocumentClaim
  ├─ product_life_plan             ← PlanList
  ├─ product_life_uw_limit         ← UnderwritingLimitList
  └─ product_life_fin_uw           ← FinancialUnderwritingList
```

⚠️ **Penyimpangan sadar (baru) — dua tabel induk existing DIGABUNG jadi satu `product_life`.**
`[keputusan work owner]` Pega memisah `M_PRODUCT_LIFE` + `M_PRODUCTINWARD_LIFE` (dua JSON blob, 1:1
by `ID`). Karena JSON dibuang & simpan atomik, pemisahan tak beralasan lagi; gabung menghapus kelas
bug "produk timpang". Field yang muncul di kedua sisi (`ID`, `CEDING`, `POLICYHOLDER`, `POLICYHODERNAME`,
`BIRTHDAY`, `TREATYNUMBER`) → **satu kolom**.

`[terverifikasi]` **Lima tabel anak.** Dua yang semula disangka tabel dicoret:
- **`LienClause` → BUKAN tabel** — `LIENCLAUSE` adalah **field skalar** (kolom di `product_life`).
- **`OutwardList` → BUKAN tabel** — hanya diisi `GetReinsTypeOR_Life` (jalur OR **mati** `1==2`).

### Kolom sisi inward — **40 field** `[terverifikasi]`

Sumber kebenaran: field yang **di-SET** di `SetProductNameInward`
(`ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SETPRODUCTNAMEINWARD` / `RULE-OBJ-ACTIVITY`):

```
ID, PRODUCTID, BEGIN, MATURE, STNC, BIRTHDAY, MONTHS, CURRENCY, PAYMENT, MAXCONTRACT,
MINAGE, MAXAGE, MINSUMINSURED, MAXSUMINSURED, MAXSUMREASURED,
CEDING, CEDINGLIMIT, CEDINGLIMITXPN, CEDINGRETENTIONNUM, CEDINGRETENTIONPCT,
RNMLIMITNUM, RNMLIMITPCT, RNMSHARE, BROKERAGE, EXTRAPREMI, EXTRAMORTALITY,
MAXDATARECEIVE, MAXEXPIREDCLAIM, INSURED, SUBJECTTO,
ADDENDUMNO, ADDENDUMWORD, AMANDEMENTNO, AMANDEMENTSCHD,
INWARDTREATYNM, PROPORTIONALTABLE, LIENCLAUSE, TREATYNUMBER,
POLICYHODER (sic), POLICYHODERNAME
```

⚠️ **Keempat puluh dipakai** `[keputusan work owner]` — termasuk **`ADDENDUMNO`**,
**`AMANDEMENTNO`**, dan **`CEDING`** yang hilang dari rekap catatan sumber.

### Kolom sisi umum — **20 field** `[terverifikasi]`

Dari `SetProductName` (`ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SETPRODUCTNAME` / `RULE-OBJ-ACTIVITY`):

```
ID, PRODUCTNAME, PRODUCTCODE, PRODUCTTYPE, TYPE, GRUP, INWARDNAME,
CEDING, CAUSE, RIRISK, RIRATE, RICOMM, BENEFIT, BIRTHDAY, PAYMENTTYPE,
ISFACULTATIVE, IsView, OUTWARDNAME, OUTWARDRATE, OUTWARDCOMM
```

Ditambah lima field yang di-set `SaveProductName_Act` (`ASM-FW-GISFW-INT-PRODUCT_LIFE` /
`SAVEPRODUCTNAME_ACT`): `ID`, `CREATEOP`, `UPDATEOP`, `POLICYHODER` (sic), `POLICYHODERNAME`.

### Kolom tabel anak `[data DBA]`

| Tabel | Kolom baris |
| --- | --- |
| `product_life_comment` | tanggal, nama operator, isi saran |
| `product_life_document_claim` | nama dokumen (+ audit) |
| `product_life_plan` | plan, id plan, nama, benefit, **RI Rate + id-nya** |
| `product_life_uw_limit` | deskripsi, status medis (`FCL`/`NM`/`MEDIS`), usia min/maks, uang pertanggungan min/maks |
| `product_life_fin_uw` | `[terverifikasi]` dari `CopyFinancialWriting`: `MinInsured`, `MaxInsured`, `Employee`, `Non_Employee` |

### Tipe kolom `[data DBA]`

| Kelompok | Tipe |
| --- | --- |
| Uang: `*LIMIT*`, `*SUMINSURED*`, `*SUMREASURED*`, batas uang pada baris underwriting | **desimal presisi arbitrer** (**ADR-0003**) |
| Persen: `*PCT`, `RNMSHARE`, `BROKERAGE`, `RICOMM` | desimal |
| Tanggal: `BEGIN`, `MATURE`, `STNC`, `BIRTHDAY` | **`DATE`** |
| Usia, jumlah hari, jumlah kontrak | bilangan bulat |
| Teks panjang: `INSURED`, `SUBJECTTO`, `INWARDNAME` | teks besar |

### Aturan pemilihan kolom `[keputusan work owner]`

1. **Kolom = field yang di-SET di Activity simpan.** Itu data yang benar-benar diisi.
2. Field yang **hanya** muncul di Section dengan visibilitas mati (`1=2` / `1==2` / `never`) **dan
   tidak di-set Activity** — **jangan diambil**. Tampilan mati bukan data.

## ADR terkait

**ADR-0003** (uang non-float), **ADR-0006** (identitas lewat sequence basis data),
**ADR-0009** (migrasi penuh; koeksistensi ditolak).

## Acceptance criteria

- [ ] ⚠️ Skema target **relasional penuh**: setiap atribut menjadi **kolom bernama**, setiap daftar
      bersarang menjadi **tabel anak**. Test yang menemukan kolom JSON sebagai penyimpan atribut
      produk **gagal**. *(AC 40 spec; `[keputusan work owner]` — penyimpangan sadar 2)*
- [ ] Daftar kolom sisi inward memuat **keempat puluh field** yang di-set Activity — termasuk
      `ADDENDUMNO`, `AMANDEMENTNO`, dan `CEDING`. *(AC 41 spec)*
- [ ] Field yang hanya tampil dengan **visibilitas mati** dan tidak di-set Activity **tidak** menjadi
      kolom. *(AC 42 spec)*
- [ ] ⚠️ **Seluruh kolom hasil migrasi NULLABLE.** Test yang menemukan `NOT NULL` pada kolom hasil
      migrasi **gagal**. Wajib-isi ditegakkan **di Go**. *(AC 43, 17 spec; `[keputusan work owner]`)*
- [ ] ⚠️ Field yang **absen** di dokumen JSON lama menjadi **kolom kosong**, bukan kegagalan migrasi.
      *(AC 44 spec; `[keputusan work owner]` — field bernilai kosong tidak muncul di JSON)*
- [ ] Migrasi mengurai **seluruh daftar bersarang** menjadi baris tabel anak; **tidak ada** daftar
      yang tertinggal di dalam dokumen. *(AC 45 spec)*
- [ ] Daftar yang **kosong** menghasilkan **nol baris anak**, bukan kegagalan. *(AC 27 spec)*
- [ ] Nilai uang pindah **tanpa berubah satu digit pun**; rekonsiliasi membandingkan **secara tepat**,
      bukan dengan toleransi. *(AC 46 spec; **ADR-0003**)*
- [ ] Nilai tanggal pindah **tanpa pergeseran zona waktu**.
- [ ] Sequence `M_PRODUCT_LIFE_SEQ` dan `M_PRODUCT_INWARD_LIFE_SEQ` pindah dengan **nilai berjalan
      yang benar**, sehingga identitas baru **tidak bertabrakan** dengan yang lama. *(AC 47 spec)*
- [ ] ⚠️ Kolom baru bernama **`POLICYHOLDER`**, bukan `POLICYHODER`; migrasi **memetakan** ejaan lama
      ke ejaan benar. *(AC 51 spec; `[keputusan work owner]` — penyimpangan sadar 5)*
- [ ] ⚠️ **Tidak ada kolom `IsORS`** di skema baru. *(AC 52 spec)*
- [ ] Migrasi dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**.
      *(AC 48 spec)*
- [ ] Skema uji yang dipakai seluruh tiket lain dibangun **dari DDL yang sama** dengan produksi —
      bukan dari tiruan yang ditulis terpisah. *(AC 49 spec)*

## Blocker

**Tidak ada.** ✅ **OQ-001, OQ-002, dan OQ-JSON-PRODUCT ditutup** `[data DBA]` — DDL, body procedure,
dan isi `JSONDATA` seluruhnya diterima.

## Catatan

⚠️ **Procedure lama tidak ikut pindah.** `[data DBA]` `POOLDATA.PEGA_M_PRODUCT_LIFE` dan
`POOLDATA.PEGA_M_PRODUCT_INWARD_LIFE` adalah **upsert JSON yang `COMMIT` sendiri**. Keduanya
**tidak** dipanggil dan **tidak** dimigrasikan — alasannya di tiket **03**.

⚠️ **Rekap catatan sumber keliru.** `dba-procedures-and-ddl.md` menyebut "28 field" lalu mendaftar
37; sensus korpus atas `SetProductNameInward` menghasilkan **40**. Aturan pemilihan kolom menetapkan
**Activity sebagai sumber kebenaran**, jadi yang berlaku **40**. `[keputusan work owner]`

✅ **Tabel anak terselesaikan dari korpus (bukan ditunda).** `FinancialUnderwritingList` = tabel
nyata, kolom `[terverifikasi]` dari `CopyFinancialWriting` (`MinInsured`, `MaxInsured`, `Employee`,
`Non_Employee`). `LienClause` = field skalar (bukan tabel). `OutwardList` = dead code jalur OR.
Semua bentuk pasti — tidak ada yang ditebak dari nama.

## Seam & perintah verifikasi

**Seam: API HTTP** (dipakai ulang dari CL-01) terhadap **skema uji Oracle nyata**. Migrasi diuji
terhadap basis data sungguhan; procedure **tidak** di-mock karena memang tidak dipakai.

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Ralat bertanggal 01-10-2026 — sesi implementasi (paket 0)

> Sumber: `../RALAT-DEV-30-09-2026.md` (P1–P6 katalog DEV, R7–R18 pembacaan ulang XML) dan `../PARITAS-LAYAR-DAN-AKSI.md`. Kalimat
> di atas **tidak dihapus**; yang berlaku adalah ralat ini.

**Status (01-10-2026): ditangguhkan → OQ-MPNL-01.** Kalimat lama *"**Status:** ready-for-agent"* tidak berlaku.

| Kalimat lama | Ralat |
| --- | --- |
| *"Sebagai **tim migrasi**, saya ingin setiap atribut produk menjadi **kolom bernama** dan setiap daftar bersarang menjadi **tabel anak**"*; *"⚠️ **Penyimpangan sadar 2 — skema relasional penuh.**"* | **P1**: tiga view (`PRODUCT_LIFE`, `PRODUCTINWARD_LIFE`, `DOCUMENTCLAIM_LIFE`), dua prosedur, dan Claim Life membaca `JSONDATA` kedua tabel lama — produk di tabel baru tidak terlihat oleh mereka. Bawaan: **ikut Pega** — `M_PRODUCT_LIFE` + `M_PRODUCTINWARD_LIFE` dengan `JSONDATA` berkunci Pega, kolom datar `RIRISKID`/`RIRISK`/`PRODUCTNAME`/`BEGIN_DATE`. **Nol tabel baru, nol DDL**, rentang 140–179 kosong |
| *"⚠️ **Penyimpangan sadar (baru) — dua tabel induk existing DIGABUNG jadi satu `product_life`.**"* | tidak berlaku (P1): dua tabel lama tetap dua, ditulis dalam **satu transaksi** (P4) |
| *"**`LienClause` → BUKAN tabel** — `LIENCLAUSE` adalah **field skalar**"* | **R10**: ada daftar `ProductName.LienClause` (`Usia`, `Manfaat`, grid b12201) **dan** skalar `ProductNameInward.LIENCLAUSE` |
| *"**`OutwardList` → BUKAN tabel** — hanya diisi `GetReinsTypeOR_Life` (jalur OR **mati** `1==2`)"* | **R9 / P3**: pengisinya **hidup** — checkbox `On Retention` b47312 → `GetReinsTypeOR_Life` (langkah ber-`PRE=false`); kunci `OutwardList` ditulis bentuk Pega |
| *"⚠️ Kolom baru bernama **`POLICYHOLDER`**, bukan `POLICYHODER`; migrasi **memetakan** ejaan lama"*; *"⚠️ **Tidak ada kolom `IsORS`**"* | P1: nol kolom baru. Kunci JSON tetap ejaan Pega `POLICYHODER` (dibaca view); `IsORS` tetap ditulis (R9). Nama medan Go boleh `PolicyHolder` — ejaan Pega hanya di repository |
| *"Sequence `M_PRODUCT_LIFE_SEQ` dan `M_PRODUCT_INWARD_LIFE_SEQ` pindah dengan **nilai berjalan yang benar**"* | tidak ada pemindahan (P1). `ID` baru dari `M_PRODUCT_LIFE_SEQ`; `ID` inward = `ID` produk (P6, R14, OQ-MPNL-02) |

---

## Ralat bertanggal 01-10-2026 — lanjutan 1 (katalog DEV)

> Sumber: brief `PROMPT-LANJUTAN-MASTER-PRODUCT-NAME-LIFE-1.md` §1–§2 (katalog DEV `ALL_TAB_COLUMNS`/`ALL_OBJECTS` dan agregat `JSONDATA`, dibaca asisten, baca-saja). Kalimat di atas tidak dihapus.

| Kalimat lama | Ralat |
| --- | --- |
| *"`JSONDATA` (constraint `IS JSON`) + empat kolom hasil flatten: `RIRISKID`, `RIRISK`, `PRODUCTNAME`, `BEGIN_DATE`"* | DEV `M_PRODUCT_LIFE` hanya `ID`, `JSONDATA`, `RIRISKID`, `RIRISK` — **dua** kolom datar; `PRODUCTNAME`/`BEGIN_DATE` tidak ada dan tidak ditulis (OQ-MPNL-08 ditutup, `6fd539c`) |

## Status 01-10-2026 — OQ-MPNL-01 ditutup (`PROMPT-LANJUTAN-TIGA-MODUL-LIFE-KEPUTUSAN-OQ.md` §2)

**ditutup 01-10-2026 — keputusan work owner: ikut rekomendasi asisten (bawaan dipertahankan)**: produk tetap **JSON seperti Pega** di `M_PRODUCT_LIFE` / `M_PRODUCTINWARD_LIFE` (P1); tiket ini **tetap
ditangguhkan** — nol tabel relasional baru, nol migrasi JSON.

---

## Keputusan bertanggal 02-10-2026 — pindah ke tabel FLAT `[keputusan work owner 02-10-2026]`

> Sumber: brief `PROMPT-PINDAH-FLAT-MASTER-PRODUCT-NAME-LIFE.md` (bab 0–2: katalog dan agregat DEV 02-10-2026, SELECT saja), jawaban
> work owner 02-10-2026 atas dua penolakan penjaga inti (di bawah), dan pembacaan ulang XML 02-10-2026. Keputusan grilling **D2**
> (relasional penuh) dan **Q1b** (dua tabel induk digabung, `grilling-ronde-1-jawaban.md`) **berlaku lagi**; RALAT P1 dicabut untuk
> penyimpanan produk. Kalimat di atas tidak dihapus — yang berlaku bab ini.

### Keputusan

| # | Keputusan |
| ---: | --- |
| K1 | Hanya modul ini. Treaty Contract Retro Life tidak — kelima tabelnya sudah flat |
| K2 | Tabel induk **`M_PRODUCTNAME_LIFE`** (sisi umum + sisi inward digabung, Q1b); tabel anak berawalan `M_PRODUCTNAME_LIFE_`. Nama `PRODUCT_LIFE` sudah dipakai view DEV |
| K3 | Nilai lama tidak sah — 1 `MATURE` bukan tanggal, 1 `UnderwritingLimitList[*].MaxInsured` bukan angka (DEV): **kolom NULL + dicatat di laporan pindah** (ID produk + nama kunci, tanpa nilainya); tabel JSON lama tetap ada sebagai cadangan |
| K4 | `OutwardList`: objek **kosong** tidak dipindahkan (DEV: 2 baris berisi dari 191); cacah yang dibuang dicatat |
| K5 | OQ-MPNL-01 dibuka ulang dan **ditutup flat** (menggantikan penutupan 01-10-2026); tiket ini tidak lagi ditangguhkan |
| K6 | ⛔ Penjaga inti `TestNolNumberTanpaPresisi` menolak `NUMBER` tanpa presisi dan `NUMBER(1)`/`(3)`/`(4)` rancangan brief §2 (dibuktikan di salinan `git archive`). Jawaban work owner: **patuhi penjaga** — uang, persen, rate, faktor → `NUMBER(38,8)` (bentuk uang sah, keputusan work owner c 26-09-2026); bilangan kecil (usia, nomor addendum/amandemen, hari, kontrak, tahun, `URUT`, `IS_ORS`) → `NUMBER(5)` |
| K7 | ⛔ Penjaga inti `TestSeluruhCreateDapatDibacaNamanya` menolak `CREATE OR REPLACE VIEW` di berkas migrasi (pola `migrasi.NamaObjekDibuat` hanya `TABLE`/`INDEX`/`SEQUENCE`). Jawaban work owner: *"tidak ada table view yang dipake, semua simpan dan baca dari table flat"* — **ketiga view tidak dibangun ulang** (T5 brief dibatalkan); modul ini menulis dan membaca tabel flat saja |

⚠️ **Akibat K7 di hilir** (dicatat, tidak diubah — di luar folder modul ini): Claim Life `repository/ambangproduk.go` membaca view
`PRODUCTINWARD_LIFE` (`NamaViewProdukLife`) — view itu tetap membaca `M_PRODUCTINWARD_LIFE.JSONDATA`, yang **berhenti diperbarui**
sesudah peralihan. Produk baru dan ubahan sesudah peralihan tidak terlihat di sana → **OQ-FLAT-04** (register OQ).

### Bentuk baru — satu induk, TUJUH anak

```
M_PRODUCTNAME_LIFE                  (PK ID; sisi umum + sisi inward)
  ├─ M_PRODUCTNAME_LIFE_LIEN        ← LienClause[*]                (Usia, Manfaat — teks)
  ├─ M_PRODUCTNAME_LIFE_DOCCLAIM    ← DocumentClaim[*]             (Document)
  ├─ M_PRODUCTNAME_LIFE_PLAN        ← PlanList[*]                  (PlanID, Plan, Name, Benefit, RIRATEID, RIRATE)
  ├─ M_PRODUCTNAME_LIFE_FINUW       ← FinancialUnderwritingList[*] (MinInsured, MaxInsured, Employee, Non_Employee)
  ├─ M_PRODUCTNAME_LIFE_UWLIMIT     ← UnderwritingLimitList[*]     (MinInsured, MaxInsured, MinAge, MaxAge, Medical, Description)
  ├─ M_PRODUCTNAME_LIFE_OUTWARD     ← OutwardList[*] berisi (K4)   (REINSTYPEID, REINSTYPENAME, TRANSACTIONYEAR, TREATYCONTRACTID, UNDERWRITINGYEAR, OVR_COMM)
  └─ M_PRODUCTNAME_LIFE_COMMENT     ← CommentList[*]               (Date → TANGGAL, OperatorName, Suggest)
```

Setiap anak: `PRODUCTID` VARCHAR2(6) NOT NULL FK → induk `ON DELETE CASCADE` · `URUT` NUMBER(5) NOT NULL (urutan baris grid,
mulai 1) · PK (`PRODUCTID`, `URUT`) — index PK berawalan `PRODUCTID` sekaligus melayani FK. Simpan = baris anak satu produk ditulis
ulang di **transaksi yang sama** dengan induk. Kolom dan tipe persisnya: `STRUKTUR-TABEL-MASTER-PRODUCT-NAME-LIFE.md`.

| Kalimat lama (dikutip) | Ralat 02-10-2026 |
| --- | --- |
| *"DDL **satu tabel induk** `product_life` + **lima** tabel anak + sequence"*; *"**Lima tabel anak.**"* | **TUJUH** anak: `LienClause` daftar hidup (R10, grid b12201) dan `OutwardList` daftar hidup (R9, `On Retention` b47312 → `GetReinsTypeOR_Life`). Nama induk `M_PRODUCTNAME_LIFE` (K2). Sequence **tidak** dibuat — `M_PRODUCT_LIFE_SEQ` warisan dipakai terus |
| *"**`LienClause` → BUKAN tabel**"*, *"**`OutwardList` → BUKAN tabel**"* | dicabut (R9, R10 — `RALAT-DEV-30-09-2026.md`) |
| *"Kolom sisi umum — **20 field** `[terverifikasi]` Dari `SetProductName`"* | ⛔ **Keliru: kedua puluh medan itu dari langkah ter-remark.** Baca ulang 02-10-2026 (`sed -e 's/></>\n</g' Activity/SetProductName.xml`, `pyStepsBlockName` dicetak): langkah 5 `Page-Copy` b1067 dan langkah 6 `Property-Set` b1206 ber-**`//`** — tidak pernah jalan; daftar `BENEFIT … TYPE` (b1279–b1636) ada di langkah 6. Jalur hidup: langkah 3 `RDB-List` b782 (`BrowseUnderwritingList`) + langkah 4 `Java` b959 mengadopsi **seluruh** JSON halaman `ProductName`; langkah 8 b2147 `IsView ← true`. Kolom sisi umum karena itu = kunci yang **ditulis** `SaveProductName_Act` (langkah 8 b1651 `@GetPageJSONString()`, halaman utuh) **dan terisi di DEV** (brief bab 0), bukan daftar langkah 6 |
| *"Kolom sisi inward — **40 field** `[terverifikasi]` … yang **di-SET** di `SetProductNameInward`"* | Benar untuk **pembaca**: langkah 3.1 b1048 menyalin 39 kunci (b1074–b1860, `pyStepsBlockName` kosong, hidup) dari `BrowseProductInward`, ditambah `PRODUCTID` langkah 1 b746 = 40. Tetapi ia hanya kunci yang ada di view `PRODUCTINWARD_LIFE`: lima medan layar hidup **tidak** di sana — `EXPIRYAGE`, `PREMIUMFACTOR`, `AnnuityInterest`, `PremiumRefundFactor`, `CURRENCYID` — dan tetap menjadi kolom. Sebaliknya delapan kunci warisan (`CEDING`, `TREATYNUMBER`, `INWARDTREATYNM`, `CEDINGRETENTIONPCT`, `CEDINGLIMITXPN`, `RNMLIMITPCT`, `LIENCLAUSE`, `MONTHS`) **0 terisi di DEV** dan tanpa medan form → **bukan** kolom |
| *"Field yang muncul di kedua sisi (`ID`, `CEDING`, `POLICYHOLDER`, `POLICYHODERNAME`, `BIRTHDAY`, `TREATYNUMBER`) → **satu kolom**"* | `POLICYHODER`/`POLICYHODERNAME` sama persis di 196 produk DEV → satu pasang kolom `POLICYHOLDER`/`POLICYHOLDERNAME` (sumbernya sisi inward, `SaveProductName_Act` 1 b535/b556). `CEDING`/`TREATYNUMBER`: kolom sisi umum (inward 0 terisi). `BIRTHDAY` hanya sisi inward |
| Tipe: *"Uang … **desimal presisi arbitrer**"*, *"Persen … desimal"*, *"Usia, jumlah hari, jumlah kontrak — bilangan bulat"*, *"Tanggal: `BEGIN`, `MATURE`, `STNC`, `BIRTHDAY` — `DATE`"* | uang/persen/rate/faktor `NUMBER(38,8)`, bilangan kecil `NUMBER(5)` (K6; desimal tetap tidak pernah float — ADR-0003); tanggal `DATE` untuk `BEGIN` (kolom `BEGIN_DATE`), `STNC`, `MATURE`. ⛔ **`BIRTHDAY` bukan tanggal**: DEV 40 terisi, panjang 1, semuanya angka — kolom `VARCHAR2(1)` (pilihan `associated`, OQ-MPNL-05) |
| AC *"⚠️ **Tidak ada kolom `IsORS`** di skema baru."* (AC 52) | dicabut: `IsORS` medan HIDUP (R9, checkbox `On Retention` b47312) → kolom `IS_ORS` NUMBER(5) `CHECK (IS_ORS IN (0, 1))` |
| AC *"Sequence `M_PRODUCT_LIFE_SEQ` dan `M_PRODUCT_INWARD_LIFE_SEQ` pindah dengan **nilai berjalan yang benar**"* (AC 47) | tidak ada sequence baru: `M_PRODUCT_LIFE_SEQ` warisan dipakai terus (`PilihIdentitasBebas` melewati ID terpakai di induk baru **dan** kedua tabel JSON); `M_PRODUCT_INWARD_LIFE_SEQ` tidak dipakai (R14, OQ-MPNL-02). ID inward = ID produk: pasangan bersilang DEV 100079 ↔ 100081 (inward dicari lewat `PRODUCTID`) pindah ke baris produknya |
| AC *"Migrasi dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**."* (AC 48) | DDL: berkas `_down` membuang ketujuh anak lalu induk. Data: **bukan berkas migrasi** — alat Go `backend/alat/pindahflat` (mode `-uji` baca-saja, `-jalankan` satu transaksi: hapus isi tabel flat lalu isi ulang dari JSON — aman diulang). Tabel JSON tidak pernah disentuh |
| AC *"Nilai uang pindah **tanpa berubah satu digit pun**; rekonsiliasi membandingkan **secara tepat**"* (AC 46) | tetap, diperketat: **setiap** medan setiap produk dibandingkan **teks demi teks** (JSON → model → baris flat → model); beda yang bukan K3/K4/pasangan bersilang = gagal. Normalisasi teks angka (mis. nol ekor) dicatat per kunci dan menghentikan `-jalankan` untuk keputusan work owner |
| *"Kunci asing tiap baris (`Asli`) dibawa bolak-balik"* (model, R-baca ulang) | dibuang (D2: JSON dibuang): kunci internal Pega baris (`pxObjClass`, jejak `pxCreate*`, …) dan kunci halaman yang tidak menjadi kolom tidak pindah; cacahnya dicatat laporan pindah |

### Tidak menjadi kolom *(0 terisi di DEV, keadaan layar, atau masukan)*

`TYPE` `TYPE_CEDING` `GRUP` `PRODUCTCODE` `PRODUCTTYPE` `PRODUCTTYPEID` `RICOMMID` `OUTWARDNAME(ID)` `OUTWARDRATE(ID)` `OUTWARDCOMM(ID)`
`BENEFIT(ID)` (medan layar mati, 0 terisi) · `IsView` (keadaan layar, `SetViewEdit`) · `Comment` (masukan popup; isinya menjadi baris
`_COMMENT` lewat `AddCommentList_Act`, langkah 7 b1515) · kunci inward warisan (delapan di atas) · `PRODUCTID` dan `ID` inward (satu baris
= satu produk). ⛔ Alat pindah **gagal** bila salah satunya ternyata terisi pada produk mana pun — kehilangan nilai tidak pernah diam-diam.

### Acceptance criteria tambahan 02-10-2026

- [ ] Rentang migrasi 140–179 hanya membuat/membuang objek `M_PRODUCTNAME_LIFE*`; **nol** `DROP`/`ALTER`/DML atas `M_PRODUCT_LIFE` /
      `M_PRODUCTINWARD_LIFE` (penjaga modul pengganti `TestMPNLNolMigrasiDiRentang`).
- [ ] Penulis dan pembaca modul hanya menyentuh tabel flat; `PilihIdentitasBebas` memeriksa induk baru dan kedua tabel JSON.
- [ ] Alat pindah: `-uji` nol tulisan; `-jalankan` menolak `IS_PEGA_PROD=true`, satu transaksi, aman diulang; laporan agregat saja.
- [ ] Nilai yang tidak muat kolomnya (teks melebihi lebar, angka > 8 desimal atau > 30 digit bulat, bilangan `NUMBER(5)` tidak bulat)
      ditolak services **berkalimat** sebelum SQL tulis — tidak pernah dipotong atau dibulatkan Oracle diam-diam.

### Pelaksanaan 02-10-2026

| Paket | Commit | Isi |
| ---: | --- | --- |
| 1 | `95ccb0c` | baca ulang XML + bab ini, spec, OQ, RALAT, dba-view |
| 2 | `873ebee` | DDL 140–147 + penjaga rentang + STRUKTUR + kaskade `MODUL.md` |
| 3 | `bf9c4e0` | pemetaan flat, alat `pindahflat` (`-uji`/`-jalankan`), rekonsiliasi |
| 4 | `353f9e80` *(sebagian terbawa commit kerja bersama `1e22ecd8`, `8914ad7e`)* | gudang, services, tiruan, uji ke tabel flat |
| 5 | — | **dibatalkan** (K7: view tidak dibangun ulang) |
| 6 | commit dokumen ini | `PANDUAN-PINDAH-FLAT.md`, `LAPORAN-MIGRASI-FLAT.md` (uji kering DEV: 0 gagal, menunggu OQ-FLAT-07) |

### Ralat 02-10-2026 sesudah `/code-review` — premis K3 dan K4 keliru (uji kering DEV 18:50)

Kalimat brief bab 0 dikutip: *"`OutwardList` 191 / 153 / maks 4 — **hanya 2 baris berisi**, 189 objek kosong"* dan K3
*"1 `MATURE` bukan tanggal"*. Uji kering dengan aturan yang diperketat `/code-review` (`LAPORAN-MIGRASI-FLAT.md`): 189 objek
itu kosong HANYA pada keenam kunci OR — keempat kunci `OUTWARDNAME`/`OUTWARDNAMEID`/`OUTWARDRATE`/`OUTWARDRATEID` berisi
(OQ-FLAT-08); `MATURE` produk 100175 tanggal dalam bentuk lain (OQ-FLAT-09). Alat menahan `-jalankan` sampai keduanya
diputuskan; K3 kini 1 nilai.

### Keputusan work owner 02-10-2026 malam — *"ikuti rekomendasi"* (OQ-FLAT-01/02/07/08/09)

K4 diralat: yang tidak dipindah hanya objek `OutwardList` yang **seluruh** kuncinya kosong; objek reasuradur outward
(`OUTWARDNAMEID`, `OUTWARDNAME`, `OUTWARDRATEID`, `OUTWARDRATE` — empat kolom baru `M_PRODUCTNAME_LIFE_OUTWARD`, migrasi 146 *(Ralat 02-10-2026 malam: keempat kolom itu ditambah migrasi BARU `148_m_productname_life_outward_kolom` (`ALTER TABLE ... ADD`), bukan di 146 - 146 ternyata sudah dijalankan di DEV pukul 15:39 (`T_MIGRASI`), sehingga isinya dikembalikan ke bentuk yang dijalankan.)*)
dipindah. K3 diralat: tanggal inward berbentuk lain **dikonversi** (OQ-FLAT-09), hanya nilai yang tidak terbaca yang
di-NULL-kan. Uji kering DEV 20:46: gagal 0, OUTWARD 191 baris (`LAPORAN-MIGRASI-FLAT.md`).

### Keputusan work owner 02-10-2026 malam — aplikasi hanya tabel flat, tanpa view

Kalimat work owner dikutip: *"HANYA MODUL PRODUCTNAME LIFE!! UBAH SEMUA JANGAN ADA YANG SIMPAN KE TABLE JSON SIMPAN KE TABLE FLAT SEMUA. DAN JANGAN GUNAKAN TABLE VIEW NYA"*.

| Jalur | Sebelum | Sesudah |
| --- | --- | --- |
| Simpan produk (Add / Edit / Copy → Save), komentar, baris outward | tabel flat | tabel flat (tidak berubah) |
| Baca produk (grid, View) | tabel flat | tabel flat (tidak berubah) |
| Penerbitan ID produk baru | induk flat **+ kedua tabel JSON** (ID terpakai dilewati) | **induk flat saja** |
| Ketiga view produk (`PRODUCT_LIFE`, `PRODUCTINWARD_LIFE`, `DOCUMENTCLAIM_LIFE`) | tidak dipakai (K7) | tidak dipakai — kini dijaga uji |
| Lampiran `M_ATTACHMENTPRODUCTNAME.DATA_JSON` | `NULL` | `NULL` (tidak berubah; isi lampiran bukan JSON) |
| Alat pindah `backend/alat/pindahflat` | membaca tabel JSON, menulis tabel flat | sama — satu-satunya jalan 196 produk lama ke tabel flat |

Penjaga `TestMPNLAplikasiHanyaTabelFlat`: kode produksi modul tidak menyebut kedua tabel JSON (kecuali definisi nama,
alat pindah, dan kodek JSON-nya) dan tidak menyebut ketiga view produk. Akibatnya alat pindah wajib dijalankan sebelum
aplikasi dipakai: ID produk lama baru terlihat oleh penerbitan ID setelah berada di induk flat. View master milik master
lain (`CURRENCY`, `CAUSEOFLOSS_LIFE`, `PRODUCT_TYPE_LIFE`, `RIRISK_LIFE_SUMMARY`, `RATE_LIFE_SUMMARY`, `RATE_LIFE`) tetap *(RALAT 07-10-2026: ringkasan rate kini tabel `M_RATE_LIFE_SUMMARY` berkolom ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID - keputusan work owner 07-10-2026, `modul/riratelife/MODUL.md` RALAT R6; modul ini membacanya `SELECT ID, USEDBY`, tetap baca-saja)*
dibaca untuk pilihan dropdown — bukan view produk modul ini.

### Permintaan work owner 03-10-2026 — tombol Copy Old

Kalimat work owner dikutip: *"TOLONG BUATKAN DI SAMPING TOMBOL ADD, TOMBOL "COPY OLD", FUNGSINYA UNTUK COPY DATA DARI TABEL LAMA YANG DARI JSON ... SAAT DI BUKA, MUNCUL POPUP, MUNCUL SEMUA LIST DARI TABEL LAMA YANG BELUM DI MIGRASI, DISETIAP LIST BISA DI CENTANG ... TOMBOL PROCESS COPY UNTUK MENGCOPY YANG DI CENTANG LALU MASUK KE TABLE BARU"*.

| Hal | Dibangun |
| --- | --- |
| Tombol | `Copy Old` tepat di samping `Add` (halaman daftar) |
| Popup | `GET /produk-lama`: SEMUA produk tabel JSON lama yang ID-nya belum ada di `M_PRODUCTNAME_LIFE`, urut ID, kolom grid + `Product Name` + `Notes`; `Search`, centang per baris, centang semua; produk yang ditolak rekonsiliasi tampil dengan alasannya (kolom, tanpa nilai) dan tidak dapat dicentang |
| `Process Copy` | `POST /produk-lama/salin`: per produk SATU transaksi - tulis induk + tujuh anak, baca ulang, bandingkan; gagal satu tidak membatalkan yang lain; hasil per produk (`Copied` / `Already in the new tables` / `Cannot be copied` / `Failed`). Daftar dibaca ulang: yang tersalin hilang dari popup dan tampil di grid |
| Aturan salin | SAMA dengan alat pindah: rekonsiliasi seluruh sumber, K3 (nilai tak sah dikosongkan, dicatat), OQ-FLAT-07 (`koma desimal`, `nol depan` diterima; jenis lain menolak), OQ-FLAT-09 (tanggal dikonversi, dicatat) |
| Tabel JSON | hanya DIBACA (jalur pindah `mpnl_pindah.go`); tidak ada tulisan ke JSON maupun view |

*Ralat atas bab sebelumnya (02-10-2026 malam), baris "Penerbitan ID produk baru → induk flat saja":* sejak Copy Old, ID produk
baru kembali **melewati ID produk lama** (`idLamaTerpakai`, baca saja). Tanpa itu produk baru dapat merebut nomor produk lama
yang belum disalin, dan produk lama itu tidak dapat lagi disalin dengan ID-nya.

