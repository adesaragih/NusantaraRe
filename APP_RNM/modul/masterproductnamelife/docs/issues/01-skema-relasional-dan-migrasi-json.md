# 01: Skema relasional penuh + migrasi JSON → kolom — **PREFACTOR**

**Status:** ditangguhkan (01-10-2026) — RALAT P1: produk tetap JSON di dua tabel lama seperti Pega, nol DDL; menunggu OQ-MPNL-01 (paket 0 `71c35b1`)

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
