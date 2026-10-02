# Spec — Master Product Name Life (migrasi Pega → Go + React + Oracle)

Status: **ready-for-agent** — ✅ **OQ-001, OQ-002, OQ-JSON-PRODUCT, OQ-056, OQ-057, OQ-058 DITUTUP**
Konteks: `master-product-name-life` — **konteks/menu sendiri** `[keputusan work owner]`
Modul: **Master Product Name Life** (114 berkas)
Tanggal: 2026-09-15
Sumber: `.scratch/master-product-name-life/grilling-ronde-1-jawaban.md` (**10 verdict final**),
`dba-procedures-and-ddl.md` `[data DBA]`, `grilling-ronde-1.md`, `CONTEXT.md`,
`docs/adr/ADR-0001`–`ADR-0015`, `discovery/modules/Master Product Name Life.md`,
`discovery/context-map.md` §2.9
Skill: `/mattpocock-skills:to-spec`

> **Konvensi penandaan.** `[terverifikasi]` = terbukti korpus dengan **class + nama + path**;
> `[keputusan work owner]`; `[fakta bisnis — work owner]`; `[data DBA]`; `[terbuka]` = OQ.
> **Identitas rule wajib menyertakan class** — nama sama di class berbeda = rule berbeda.

> **Sumber tunggal.** Korpus Pega `D:\XML\RNM_BRD\` (READ-ONLY) dan artefak di
> `OUTPUT_HASIL_RNM\`, sekaligus repo target tunggal. ⛔ `D:\XML\nusantara-re\` di-blacklist.

> `[keputusan work owner]` **Master Contract Retro Life adalah konteks terpisah** dan sudah selesai
> (spec + 12 tiket). Modul ini punya spec dan tiketnya sendiri.

---

## Problem Statement

Setiap produk asuransi jiwa yang ditanggung ulang punya **definisi master**: siapa cedingnya, siapa
pemegang polisnya, mata uangnya, rentang usia yang diterima, batas uang pertanggungan, batas retensi
ceding, share Nusantara Re, brokerage, dan syarat-syarat penerimaan lain. Hari ini definisi itu
dikelola di satu layar Pega dan disimpan sebagai **satu dokumen JSON** ke dua tabel Oracle.

Lima hal membuat modul ini berisiko dipindahkan secara salah:

1. **Satu produk, dua tabel, dan tidak atomik.** `[data DBA]` Kedua procedure penulis **commit
   sendiri**. Bila yang kedua gagal setelah yang pertama berhasil, produk tersimpan **separuh** —
   dan tidak ada yang memberitahu siapa pun.
2. **Datanya JSON, sehingga bentuknya tidak dapat dipercaya.** `[data DBA]` `M_PRODUCTINWARD_LIFE`
   hanya punya `ID` + `JSONDATA`. Tidak ada kolom, tidak ada tipe, tidak ada `NOT NULL`. Apa pun
   dapat masuk, dan **field bernilai kosong tidak muncul sama sekali** di dokumennya.
3. **Empat jalur besar di dalam modul ternyata mati.** Jalur treaty inward, jalur klaim, jalur `OR`,
   dan satu jalur simpan pintas — seluruhnya **tidak dipakai**, tetapi ikut terekspor dan tampak
   nyata. Memigrasikannya berarti membangun ulang sesuatu yang tidak ada gunanya.
4. **Dua salah ketik terbawa sampai ke data.** `POLICYHODER` (kurang satu `L`) adalah **nama field
   di dalam JSON produksi**, dan `PoductName` (kurang satu `r`) menjaga justru langkah yang memanggil
   penyimpanan.
5. **Lampiran bercampur dengan metadata.** Unggahan berkas ke Google Storage berada di jalur yang
   sama dengan penyimpanan produk, sehingga gangguan pada penyimpanan berkas dapat menghalangi
   pekerjaan yang sebenarnya tidak memerlukannya.

## Solution

Membangun ulang Master Product Name Life sebagai **editor master produk life murni** di atas
Go + React + Oracle — bukan proses berjenjang: tidak ada persetujuan, tidak ada Submit/Decline,
tidak ada status akseptasi.

`[fakta bisnis — work owner]` **Satu produk**, disimpan ke **dua tabel** yang terhubung lewat
**`ID` yang sama**. Bukan dua produk berbeda.

### Lima penyimpangan sadar

| # | Penyimpangan | Alasan | Sumber |
| --- | --- | --- | --- |
| 1 | **Satu transaksi atomik** (satu baris `product_life` + tabel anak) — procedure JSON lama **tidak dipakai** | `[data DBA]` keduanya `COMMIT` sendiri, sehingga atomik mustahil selama ia dipakai | `[keputusan work owner]` |
| 1b | **Dua tabel induk existing digabung jadi satu `product_life`** (relasi 1:1 by `ID`) | JSON dibuang + atomik → pemisahan tak beralasan; gabung menghapus kelas bug "produk timpang" | `[keputusan work owner]` |
| 2 | **Skema relasional penuh**, bukan JSON — atribut menjadi kolom bernama, list bersarang menjadi tabel anak; data JSON dimigrasikan | bentuk JSON tidak dapat divalidasi, tidak dapat di-index, dan field kosong tidak terlihat | `[keputusan work owner]` |
| 3 | **Gagal terang-terangan**; identitas `'1' + lpad(sequence, 5, '0')` **ditiru** | keluaran procedure lama tidak dipakai lagi; format identitas tetap dipertahankan | `[keputusan work owner]` |
| 4 | **Lampiran opsional dan terpisah** dari metadata produk; unggah = efek keluar | lampiran pelengkap, bukan syarat sah produk | `[keputusan work owner]` |
| 5 | **Salah ketik dibuang** (`PoductName` → `ProductName`, `POLICYHODER` → `POLICYHOLDER`) dan **cabang mati `IsORS` dibuang** | nama yang salah merambat ke skema baru bila ditiru | `[keputusan work owner]` |

### Kontrak batas

| Arah | Isi |
| --- | --- |
| **Keluar** | Modul ini **menghasilkan** master produk life, dikonsumsi **Claim Life**, **PremiumList Life**, dan **Endorsement Life** (membaca `M_PRODUCT_LIFE` / `PRODUCT_LIFE`) |
| **Masuk** | Tujuh master pendukung (lihat §4); berkas lampiran ke **Google Storage** |
| **Tidak ada** | ⚠️ `[keputusan work owner]` Modul ini **TIDAK** menulis master treaty inward. Jalur `M_TREATY_IN` adalah **salah ekspor** — lihat §Out of Scope |

---

## User Stories

### Mengelola produk

1. Sebagai **admin master produk life**, saya ingin membuat definisi **produk** baru beserta nama,
   ceding, pemegang polis, dan sumber bisnisnya, supaya produk dapat dipakai lini life.
2. Sebagai **admin master**, saya ingin **tidak perlu menentukan nomor identitas** produk, supaya
   penomoran tidak pernah bentrok antar petugas.
3. Sebagai **admin master**, saya ingin identitas produk tetap berbentuk seperti yang sudah dikenal
   perusahaan, supaya rujukan lama tetap terbaca.
4. Sebagai **admin master**, saya ingin mengubah produk yang sudah ada tanpa membuat duplikat.
5. Sebagai **underwriter**, saya ingin menetapkan **syarat penerimaan** produk — rentang usia, batas
   uang pertanggungan, batas retensi ceding, share Nusantara Re, brokerage — supaya batas risiko
   jelas sejak awal.
6. Sebagai **underwriter**, saya ingin menetapkan **masa berlaku** produk (mulai, jatuh tempo, STNC),
   supaya produk tidak dipakai di luar periodenya.
7. Sebagai **Finance**, saya ingin **tidak ada nilai uang yang berubah** saat menyeberang batas
   penyimpanan, supaya batas pertanggungan persis seperti yang ditetapkan.

### Satu produk, dua tabel

8. Sebagai **organisasi**, saya ingin kedua sisi produk tersimpan **bersama atau tidak sama sekali**,
   supaya tidak pernah ada produk yang tersimpan separuh. `[keputusan work owner]`
9. Sebagai **organisasi**, saya ingin kedua sisi berbagi **identitas yang sama**, supaya keduanya
   selalu dapat dipasangkan. `[fakta bisnis — work owner]`
10. Sebagai **admin master**, saya ingin **diberi tahu bila penyimpanan gagal**, supaya saya tidak
    mengira produk tersimpan padahal tidak.

### Tujuh pemilih master

11. Sebagai **admin master**, saya ingin memilih **ceding** dari daftar master, supaya tidak ada
    nama yang diketik bebas.
12. Sebagai **admin master**, saya ingin memilih **pemegang polis** dari daftar master.
13. Sebagai **admin master**, saya ingin memilih **mata uang** dari daftar master.
14. Sebagai **admin master**, saya ingin memilih **sumber bisnis** (*Source of Business*) dari daftar
    master. `[fakta bisnis — work owner]`
15. Sebagai **admin master**, saya ingin memilih **penyebab klaim** dari daftar master.
16. Sebagai **underwriter**, saya ingin memilih **RI Risk** dari masternya sendiri.
17. Sebagai **underwriter**, saya ingin memilih **RI Rate** dari masternya sendiri — **master yang
    berbeda** dari RI Risk. `[fakta bisnis — work owner]`

### Rincian berulang

18. Sebagai **underwriter**, saya ingin mencatat **beberapa plan** dalam satu produk, masing-masing
    dengan benefit dan RI Rate-nya, supaya satu produk dapat menawarkan beberapa paket.
19. Sebagai **underwriter**, saya ingin mencatat **beberapa batas underwriting** — deskripsi, status
    medis, rentang usia, rentang uang pertanggungan — supaya syarat bertingkat dapat dinyatakan.
20. Sebagai **admin master**, saya ingin mencatat **dokumen klaim** yang disyaratkan produk.
21. Sebagai **admin master**, saya ingin meninggalkan **komentar** pada produk beserta tanggal dan
    nama penulisnya, supaya riwayat pertimbangan tidak hilang.
22. Sebagai **underwriter**, saya ingin mencatat **klausul lien**, **underwriting finansial**, dan
    **daftar outward** bila produk memerlukannya.
23. Sebagai **admin master**, saya ingin menambah dan menghapus baris pada daftar-daftar itu tanpa
    menyentuh produk induknya.

### Lampiran

24. Sebagai **admin master**, saya ingin **melampirkan berkas** pada produk, supaya nota dan
    dokumen pendukung tersimpan bersamanya.
25. Sebagai **admin master**, saya ingin produk **tetap tersimpan meski lampiran gagal diunggah**,
    supaya gangguan penyimpanan berkas tidak menghentikan pekerjaan saya. `[keputusan work owner]`
26. Sebagai **admin master**, saya ingin **melihat status lampiran** — terunggah, gagal, atau belum
    — supaya saya tahu apa yang perlu diulang.
27. Sebagai **admin master**, saya ingin **mengunduh** dan **menghapus** lampiran.
28. Sebagai **tim operasi**, saya ingin alamat penyimpanan berkas **dibaca dari konfigurasi**, supaya
    perpindahan lingkungan tidak menuntut penempelan ulang.

### Skema dan migrasi

29. Sebagai **tim migrasi**, saya ingin setiap atribut produk menjadi **kolom bernama**, supaya
    bentuknya dapat diperiksa dan dicari. `[keputusan work owner]`
30. Sebagai **tim migrasi**, saya ingin daftar bersarang menjadi **tabel anak**, supaya barisnya
    dapat dihitung dan dirujuk.
31. Sebagai **tim migrasi**, saya ingin data JSON lama **dipindahkan ke kolom** tanpa kehilangan
    satu nilai pun.
32. Sebagai **tim migrasi**, saya ingin field yang **kosong di JSON** menjadi **kolom kosong**, bukan
    kegagalan migrasi. `[keputusan work owner]`

### Yang sengaja tidak dibawa

33. Sebagai **tim migrasi**, saya ingin **jalur treaty inward tidak ikut pindah** — ia salah ekspor.
34. Sebagai **tim migrasi**, saya ingin **jalur `OR` tidak ikut pindah** — ia mati di balik gerbang
    yang selalu salah.
35. Sebagai **tim migrasi**, saya ingin **salah ketik tidak ikut pindah** ke nama kolom baru.
36. Sebagai **tim migrasi**, saya ingin **jalur simpan pintas tidak ikut pindah**.

---

## Implementation Decisions

### 1. Batas konteks

`[keputusan work owner]` **Konteks ini = Master Product Name Life saja.** Master Contract Retro Life
punya spec dan tiketnya sendiri.

`[terverifikasi]` Titik masuk `Master Product Name Life/Harness/InwardProductName.xml`
(`ASM-FW-GISFW-INT-PRODUCT_LIFE` / `INWARDPRODUCTNAME` / `RULE-HTML-HARNESS`, 540.636 byte) —
**satu-satunya Harness**, dan satu-satunya titik masuk Tahap 5 yang menempel langsung pada class
entitas data, bukan `DATA-PORTAL`.

`[terverifikasi]` **Bukan proses berjenjang**: nol rujukan `StatusAkseptasi`, tanpa `Akseptasi_DT`,
tanpa tombol Submit/Decline. **Nol rule `Flow`.**

Konteks hilir: **Claim Life**, **PremiumList Life**, **Endorsement Life**.

### 2. Penempatan modul

`[terverifikasi]` Struktur mengikat `CLAUDE.md` §5: kode di-scaffold **di dalam
`OUTPUT_HASIL_RNM\`**, arah dependency `handlers` → `services` → `repository`.

| Lapisan | Tanggung jawab |
| --- | --- |
| `models` | Produk (satu entitas: sisi umum + sisi inward); tujuh master pendukung; lima jenis baris anak |
| `repository` | **Satu transaksi** menulis baris `product_life` + tabel anak; klien penyimpanan berkas terpisah |
| `services` | Wajib-isi; orkestrasi simpan atomik; pemisahan metadata dan lampiran |
| `handlers` | Endpoint CRUD produk; endpoint tiap daftar anak; endpoint lampiran |
| `frontend` | Layar produk; tujuh dialog pemilih master; grid tiap daftar anak; panel lampiran |

### 3. Entitas — satu produk, SATU tabel induk, lima tabel anak

> ✅ `[keputusan work owner]` **Dua tabel induk existing DIGABUNG jadi satu `product_life`.**
> Di Pega dipisah (`M_PRODUCT_LIFE` + `M_PRODUCTINWARD_LIFE`, dua JSON blob) — relasi **1:1** by `ID`.
> Karena JSON dibuang (D2) & simpan atomik (D1), pemisahan tak lagi punya alasan; menggabung menghapus
> seluruh kelas bug "produk timpang antar dua induk". Field yang muncul di kedua sisi (`ID`, `CEDING`,
> `POLICYHOLDER`, `POLICYHODERNAME`, `BIRTHDAY`, `TREATYNUMBER`) → **satu kolom**. **Penyimpangan sadar.**
> `[terverifikasi]` **Lima tabel anak.** `LienClause` = field skalar (kolom di `product_life`, bukan tabel);
> `OutwardList` = dead code (jalur OR mati `1==2`). Anak nyata: comment, document_claim, plan, uw_limit,
> fin_uw (`MinInsured`/`MaxInsured`/`Employee`/`Non_Employee`).

`[fakta bisnis — work owner]` **Satu produk** (existing: dua tabel 1:1 by `ID`; skema baru: **satu
tabel `product_life`** hasil gabung):

```
product_life   (ID + ~20 field umum + 40 field inward, field duplikat → satu kolom;
                TERMASUK LIENCLAUSE skalar)
  ├─ product_life_comment          ← CommentList
  ├─ product_life_document_claim   ← DocumentClaim
  ├─ product_life_plan             ← PlanList
  ├─ product_life_uw_limit         ← UnderwritingLimitList
  └─ product_life_fin_uw           ← FinancialUnderwritingList
```

`[terverifikasi]` Isi tiap daftar anak:

| Tabel anak | Kolom baris |
| --- | --- |
| `product_life_comment` | `Date`, `OperatorName`, `Suggest` (isi saran) |
| `product_life_document_claim` | `Document` (+ audit) |
| `product_life_plan` | plan, id plan, nama, benefit, **RI Rate + id-nya** |
| `product_life_uw_limit` | deskripsi, status medis (`FCL`/`NM`/`MEDIS`), usia min/maks, uang pertanggungan min/maks |
| `product_life_fin_uw` | `MinInsured`, `MaxInsured`, `Employee`, `Non_Employee` (`[terverifikasi]` `CopyFinancialWriting`) |

**Bukan tabel:** `LienClause` (field skalar di `productinward_life`), `OutwardList` (dead code jalur OR).

### 4. Tujuh pemilih master

`[fakta bisnis — work owner]`

| Pemilih | Arti | FlowAction `[terverifikasi]` |
| --- | --- | --- |
| Ceding | perusahaan ceding | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSECEDING` |
| Pemegang polis | pemegang polis | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSEPOLICYHOLDER` |
| Mata uang | mata uang | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSECURRENCY` |
| **SOB** | **Source of Business** | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSESOB` |
| Penyebab klaim | penyebab klaim/rugi | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSECAUSEOFLOSS` |
| **RI Risk** | master RI Risk — mengisi `RIRISKID`/`RIRISK` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `CHOOSERIRISK` |
| **RI Rate** | master RI Rate — **berbeda** dari RI Risk | ⚠️ `ASM-FW-GISFW-DATA-PLAN` / `CHOOSERIRATE` — **class berbeda** |

⚠️ `[fakta bisnis — work owner]` **RI Rate ≠ RI Risk.** Dua master berbeda; **jangan digabung**.
RI Rate dipakai pada **baris plan** (tabel anak), RI Risk pada **produk induk**.

`[terbuka]` Kepanjangan "RI Rate" dan "RI Risk" — **OQ kecil**, tidak memblokir.

### 5. Simpan — satu transaksi atomik, tanpa procedure lama

⚠️ **Penyimpangan sadar 1.** `[keputusan work owner]`

`[data DBA]` Procedure lama `POOLDATA.PEGA_M_PRODUCT_LIFE` dan
`POOLDATA.PEGA_M_PRODUCT_INWARD_LIFE` adalah **upsert JSON** yang **`COMMIT` sendiri**. Selama
keduanya dipakai, penyimpanan dua tabel **tidak mungkin atomik**.

**Keputusan:** sistem baru **tidak memanggil dan tidak memigrasikan** kedua procedure itu. Sisi umum
+ sisi inward digabung jadi **satu baris `product_life`**; Go menulis baris itu **beserta seluruh
tabel anak** di dalam **satu transaksi**; kegagalan di mana pun **membatalkan seluruhnya**.

⚠️ Ini **berbeda** dari pola "panggil procedure apa adanya" yang berlaku di empat konteks Life
sebelumnya (**ADR-0006**). Alasannya tunggal dan tegas: **procedure itu menghalangi atomik yang
diminta**.

`[terverifikasi]` Rule Pega yang menjadi rujukan — **tidak dimigrasikan**, hanya sebagai jejak:

| Rule | Class / Nama / Tipe |
| --- | --- |
| `SaveProductNameLIfe` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!SAVEPRODUCTNAMELIFE` / `RULE-CONNECT-SQL` |
| `SaveProductNameInwardLIfe` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!SAVEPRODUCTNAMEINWARDLIFE` / `RULE-CONNECT-SQL` |
| `SaveProductName_Act` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SAVEPRODUCTNAME_ACT` / `RULE-OBJ-ACTIVITY` (165.682 byte, 16 langkah) |
| `SaveInwardProductName_Act` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SAVEINWARDPRODUCTNAME_ACT` / `RULE-OBJ-ACTIVITY` (64.796 byte, 6 langkah) |

### 6. Identitas dan hasil simpan

⚠️ **Penyimpangan sadar 3.** `[keputusan work owner]`

`[data DBA]` Format identitas lama **ditiru**: **`'1' + lpad(sequence, 5, '0')`** — lima digit,
misalnya `100001`. Sumbernya sequence `M_PRODUCT_LIFE_SEQ` dan `M_PRODUCT_INWARD_LIFE_SEQ`.
Konsisten **ADR-0006**: aplikasi **tidak menyusun** identitas sendiri.

`[data DBA]` Kontrak keluaran procedure lama — `StsSave` **100 = sukses / 99 = gagal**, `ErrMsg`
teks (terisi **juga saat sukses**), `IDPegaOut` identitas final — **tidak dipakai**, karena
procedure-nya dibuang. Ia dicatat sebagai **rujukan**, bukan perilaku.

**Gagal terang-terangan:** kegagalan penyimpanan **ditampilkan**; penyimpanan yang gagal **tidak
pernah** tampak berhasil.

### 7. Skema relasional penuh — dan dari mana daftar kolomnya

⚠️ **Penyimpangan sadar 2.** `[keputusan work owner]`

`[data DBA]` **Keadaan lama:**

| Tabel | Bentuk |
| --- | --- |
| `M_PRODUCTINWARD_LIFE` | **hanya** `ID` + `JSONDATA` |
| `M_PRODUCT_LIFE` | `JSONDATA` (dengan constraint `IS JSON`) + empat kolom hasil flatten: `RIRISKID`, `RIRISK`, `PRODUCTNAME`, `BEGIN_DATE` |
| Kunci | **tanpa PK** — hanya index |

**Keputusan:** skema baru **relasional penuh**. Setiap atribut menjadi **kolom bernama**; setiap
daftar bersarang menjadi **tabel anak**; data JSON lama **dimigrasikan ke kolom**.

#### Aturan pemilihan kolom `[keputusan work owner]`

1. **Kolom = field yang di-SET di Activity simpan.** Itu data yang benar-benar diisi.
2. **Field yang hanya muncul di Section dengan visibilitas mati** (`1=2` / `1==2` / `never`) **dan
   tidak di-set Activity** — **jangan diambil**. Tampilan mati bukan data.

`[terverifikasi]` **Sisi inward — 40 field** di-set di `SetProductNameInward`
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

⚠️ **Koreksi terhadap catatan sumber.** `dba-procedures-and-ddl.md` menyebut "**28 field**" lalu
mendaftar **37**. Sensus korpus atas `SetProductNameInward` menghasilkan **40**. Tiga yang hilang
dari daftar itu: **`ADDENDUMNO`**, **`AMANDEMENTNO`**, **`CEDING`**. Karena aturan pemilihan kolom
menetapkan **Activity sebagai sumber kebenaran**, yang berlaku adalah **40**.

`[terverifikasi]` **Sisi umum — 20 field** di-set di `SetProductName`
(`ASM-FW-GISFW-INT-PRODUCT_LIFE` / `SETPRODUCTNAME` / `RULE-OBJ-ACTIVITY`):

```
ID, PRODUCTNAME, PRODUCTCODE, PRODUCTTYPE, TYPE, GRUP, INWARDNAME,
CEDING, CAUSE, RIRISK, RIRATE, RICOMM, BENEFIT, BIRTHDAY, PAYMENTTYPE,
ISFACULTATIVE, IsView, OUTWARDNAME, OUTWARDRATE, OUTWARDCOMM
```

`[terverifikasi]` Ditambah lima field yang di-set `SaveProductName_Act`: `ID`, `CREATEOP`,
`UPDATEOP`, `POLICYHODER` (sic), `POLICYHODERNAME`.

#### Tipe kolom

`[data DBA]`

| Kelompok | Tipe |
| --- | --- |
| Uang: `*LIMIT*`, `*SUMINSURED*`, `*SUMREASURED*`, batas uang pada baris underwriting | **desimal presisi arbitrer** (**ADR-0003**) — **tidak** lewat `float` |
| Persen: `*PCT`, `RNMSHARE`, `BROKERAGE`, `RICOMM` | desimal |
| Tanggal: `BEGIN`, `MATURE`, `STNC`, `BIRTHDAY` | **`DATE`** |
| Usia, jumlah hari, jumlah kontrak | bilangan bulat |
| Teks panjang: `INSURED`, `SUBJECTTO`, `INWARDNAME` | teks besar |

#### Aturan migrasi `[keputusan work owner]`

1. ⚠️ **Field bernilai kosong TIDAK muncul di JSON.** Migrasi membaca tiap field dari dokumen dan
   menghasilkan **kosong bila field absen** — bukan kegagalan.
2. ⚠️ **Seluruh kolom hasil migrasi NULLABLE.** Tidak boleh ada `NOT NULL` yang disandarkan pada
   asumsi "field ini pasti ada". **Wajib-isi ditegakkan di Go**, bukan oleh basis data.
3. **Daftar kolom berasal dari Activity**, bukan dari superset visibilitas Section.

### 8. Validasi

`[terverifikasi]` `SaveProductName_Act` memeriksa lima hal sebelum menyimpan:

| Pemeriksaan | Precondition korpus |
| --- | --- |
| Tipe & grup terisi | `ProductName.TYPE=="" \|\| ProductName.GRUP==""` |
| Nama produk terisi | `ProductName.PRODUCTNAME==""` |
| Ceding terisi | `ProductName.CEDING==""` |
| Pemegang polis terisi | `ProductNameInward.POLICYHODERNAME==""` |
| Sumber bisnis terisi | `ProductName.SOBNAME==""` |

⚠️ `[terverifikasi]` Di `SaveInwardProductName_Act`, dua pemeriksaan padanannya **ter-remark**
(`<pyStepsBlockName>//` baris 414 dan 601). Sistem baru **menegakkan validasi yang sama di kedua
sisi** — dua jalur atas satu produk tidak boleh berbeda ketatnya.

### 9. Salah ketik dan cabang mati yang dibuang

⚠️ **Penyimpangan sadar 5.** `[keputusan work owner]`

| Yang dibuang | Bukti |
| --- | --- |
| **`PoductName`** → `ProductName` | `[terverifikasi]` ejaan tanpa `r` muncul **20 kali di dua berkas saja** — keduanya activity simpan — dan **tidak pernah dideklarasikan sebagai halaman**; `ProductName` muncul **714 kali** |
| **`POLICYHODER`** → `POLICYHOLDER` | `[data DBA]` nama field di dalam JSON produksi, kurang satu `L` |
| **`IsORS`** dan seluruh cabangnya | `[terverifikasi]` dibandingkan **3× `=="true"`** dan **2× `=="'true'"`**, dengan **nol** tempat men-setnya; jalur `OR` yang dilayaninya **mati** |

⚠️ Membuang `PoductName` **mengubah perilaku**: guard pada langkah simpan yang selama ini mungkin
tak pernah menyala akan **mulai menolak** penyimpanan yang dahulu lolos. Itu disadari, bukan
diselundupkan.

### 10. Lampiran dan Google Storage

⚠️ **Penyimpangan sadar 4.** `[keputusan work owner]`

`[terverifikasi]` Rantai berkas ber-class **`ASM-FW-GISFW-INT-T_STORAGE_IMAGE`**:

| Rule | Class / Nama / Tipe |
| --- | --- |
| `ServiceGoogle` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `SERVICEGOOGLE` / `RULE-CONNECT-REST` |
| `GetTokenStorage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RNM!GETTOKENSTORAGE_SQL` / `RULE-CONNECT-SQL` |
| `GetLinkStorage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RNM!GETLINKSTORAGE_SQL` / `RULE-CONNECT-SQL` |
| `Insert_T_Storage_SQL` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `RNM!INSERT_T_STORAGE_SQL` / `RULE-CONNECT-SQL` |
| `GetMimeType` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` / `GETMIMETYPE` / `RULE-OBJ-DECISIONTABLE` |

**Keputusan:**

- **Lampiran OPSIONAL dan TERPISAH dari metadata.** Produk tersimpan lebih dulu; lampiran menyusul
  dengan **status yang jelas**. Kegagalan unggah **tidak** membatalkan produk.
- Unggah adalah **efek keluar** (**ADR-0015**): kegagalannya **tidak boleh diam-diam**, tercatat, dan
  **dapat diulang**.
- Alamat penyimpanan **di-resolve runtime** dari konfigurasi / `M_LINK_SERVICE` (**ADR-0013**) —
  dilarang sebagai literal, konstanta, maupun env var.
- `[data DBA]` Token mengikuti pola `GET_TOKEN_STORAGE`: **MD5**, **di-cache**, **berlaku 1 menit**.

⚠️ Berbeda dari Master Contract Retro Life yang **tanpa integrasi luar sama sekali**, konteks ini
punya efek keluar nyata — **ADR-0013** dan **ADR-0015** berlaku.

---

## Testing Decisions

### Apa yang membuat test baik di sini

Test memeriksa **perilaku yang terlihat dari luar** — apa yang tersimpan, apa yang ditolak, apa yang
tampil — bukan susunan internal.

Tiga hal di konteks ini **wajib** diuji terhadap data nyata: **atomik** (baris produk dan seluruh
tabel anak tersimpan bersama atau tidak sama sekali), **angka** (uang tidak berubah), dan
**migrasi** (JSON → kolom tanpa kehilangan nilai).

### Seam — **memakai ulang seam yang sudah ada**

`[terverifikasi]` Repo target belum di-scaffold. Seam yang ditetapkan spec Claim — Life adalah
**API HTTP**, dan konteks ini memakainya kembali — **tidak menambah seam**:

> **Seam utama: API HTTP Master Product Name Life.** Test menggerakkan CRUD produk lewat endpoint
> REST dan memeriksa hasilnya lewat endpoint REST, dengan `handlers → services → repository`
> terpasang sungguhan, terhadap skema uji Oracle.

**Batas proses difake:**

| Batas | Perlakuan |
| --- | --- |
| Oracle | **skema uji nyata, bukan mock** — atomik lintas satu induk `product_life` + lima tabel anak hanya berperilaku benar pada basis data sungguhan |
| Google Storage | ***fake* di balik interface** — test memeriksa **efeknya**: unggah terjadi/tidak, status lampiran, kegagalan tercatat dan dapat diulang |

**Tidak ada seam kedua.** Konteks ini tidak punya worker asinkron.

### Modul yang diuji

| Yang diuji | Lewat seam |
| --- | --- |
| CRUD produk; identitas dibuat basis data | API HTTP + skema uji |
| **Atomik**: gagal di tabel anak mana pun membatalkan seluruh simpan (baris produk + anak) | API HTTP + skema uji |
| Sisi umum & inward = satu baris `product_life` (field duplikat satu kolom) | API HTTP + skema uji |
| Tujuh pemilih master; RI Rate dan RI Risk **tidak tertukar** | API HTTP + skema uji |
| Lima daftar anak: tambah, ubah, hapus baris | API HTTP + skema uji |
| Validasi lima field wajib, **sama di kedua sisi** | API HTTP |
| Uang tidak berubah menyeberang batas | API HTTP |
| Lampiran opsional: produk tersimpan meski unggah gagal | API HTTP + fake |
| Status lampiran terlihat dan dapat diulang | API HTTP + fake |
| Alamat penyimpanan di-resolve runtime | API HTTP + skema uji |
| **Migrasi**: JSON → kolom, field absen menjadi kosong | skema uji |

### Prior art

`[terverifikasi]` **Tidak ada** — nol kode, nol test di `OUTPUT_HASIL_RNM`. Spec Claim — Life,
Komite Claim Life, PremiumList Life, Endorsement Life, dan Master Contract Retro Life menetapkan
bentuknya; konteks ini mengikuti bentuk yang sama.

Perintah verifikasi wajib ditulis eksplisit di tiap tiket selama `Makefile` belum ada. Target:
`go test ./internal/...` dan `cd frontend && npm test`.

---

## Acceptance Criteria

**Produk dan identitas**

1. Produk dapat dibuat dengan nama, tipe, grup, ceding, pemegang polis, dan sumber bisnis.
2. **Aplikasi tidak menetapkan identitas produk** — identitas dibuat basis data lewat sequence.
   *(**ADR-0006**)*
3. ⚠️ Identitas berbentuk **`'1' + lima digit`** (misalnya `100001`) — format lama **ditiru**.
   `[keputusan work owner]`
4. Menyimpan produk yang sudah ada **memperbarui** baris itu, tidak membuat duplikat.
5. Produk dapat dibaca kembali utuh lewat API, termasuk seluruh daftar anaknya.

**Satu produk, satu tabel induk — atomik**

6. ⚠️ Baris `product_life` (sisi umum + inward) **beserta seluruh baris tabel anak** ditulis dalam
   **satu transaksi**. `[keputusan work owner]`
7. ⚠️ Kegagalan pada **tabel anak mana pun** membatalkan **seluruh** simpan; dibuktikan dengan
   menyuntikkan kegagalan pada satu tabel anak lalu memastikan **tidak ada baris** `product_life`.
   `[keputusan work owner]`
8. ⚠️ Sisi umum & sisi inward = **satu baris** `product_life` (dua tabel existing digabung, 1:1 by
   `ID`); field yang muncul di kedua sisi jadi **satu kolom**. `[keputusan work owner]`
9. ⚠️ **Procedure JSON lama tidak dipanggil.** Test yang menemukan pemanggilan
   `PEGA_M_PRODUCT_LIFE` atau `PEGA_M_PRODUCT_INWARD_LIFE` **gagal**. `[keputusan work owner]`
10. ⚠️ **Gagal terang-terangan**: penyimpanan yang gagal **ditampilkan** dan **tidak pernah** tampak
    berhasil. `[keputusan work owner]`

**Validasi**

11. Produk ditolak bila **tipe** atau **grup** kosong.
12. Produk ditolak bila **nama produk** kosong.
13. Produk ditolak bila **ceding** kosong.
14. Produk ditolak bila **pemegang polis** kosong.
15. Produk ditolak bila **sumber bisnis** kosong.
16. Kelima pemeriksaan berlaku **sama** pada kedua sisi produk — tidak ada sisi yang lebih longgar.
    *(§8; Pega me-remark dua di antaranya pada sisi inward)*
17. Wajib-isi ditegakkan **di Go**, bukan oleh basis data — seluruh kolom hasil migrasi nullable.

**Tujuh pemilih master**

18. Ceding, pemegang polis, mata uang, sumber bisnis, penyebab klaim, RI Risk, dan RI Rate dipilih
    dari masternya masing-masing — **tidak** diketik bebas.
19. ⚠️ **RI Rate dan RI Risk adalah dua master berbeda** dan tidak pernah tertukar. RI Risk mengisi
    produk induk; RI Rate mengisi **baris plan**. `[fakta bisnis — work owner]`
20. Nilai yang dipilih tersimpan **beserta identitas masternya**, bukan hanya namanya.
21. Pilihan yang tidak ada di master **ditolak**.

**Daftar anak**

22. Produk dapat memuat **beberapa plan**, masing-masing dengan benefit dan RI Rate-nya.
23. Produk dapat memuat **beberapa batas underwriting** dengan deskripsi, status medis, rentang usia,
    dan rentang uang pertanggungan.
24. Produk dapat memuat **dokumen klaim**, **komentar**, **klausul lien**, **underwriting finansial**,
    dan **daftar outward**.
25. Komentar tersimpan beserta **tanggal dan nama penulisnya**.
26. Baris dapat ditambah dan dihapus tanpa menyentuh produk induk — kecuali ketika penyimpanan
    dilakukan sebagai satu kesatuan (AC 7).
27. Daftar yang **kosong** tersimpan sebagai daftar kosong, bukan kegagalan.

**Uang dan tanggal**

28. Nilai uang diperlakukan sebagai **desimal presisi arbitrer**; **tidak ada** yang melewati `float`
    di lapisan mana pun maupun di JSON API. *(**ADR-0003**)*
29. Nilai uang yang ditulis dan dibaca kembali **identik** — tidak ada pembulatan diam.
30. Persentase (share, brokerage, komisi, retensi) juga desimal, bukan pembulatan ke bilangan bulat.
31. Tanggal tersimpan sebagai **tanggal**, bukan teks.
32. Batas bawah yang lebih besar dari batas atas — usia maupun uang pertanggungan — **ditolak**.

**Lampiran**

33. ⚠️ **Lampiran OPSIONAL**: produk tersimpan **tanpa** lampiran tanpa keluhan.
    `[keputusan work owner]`
34. ⚠️ **Lampiran TERPISAH dari metadata**: kegagalan unggah **tidak** membatalkan produk yang sudah
    tersimpan. `[keputusan work owner]`
35. ⚠️ Kegagalan unggah **tercatat dan terlihat** — bukan senyap — dan **dapat diulang**.
    *(**ADR-0015**)*
36. Status tiap lampiran terbaca lewat API: terunggah, gagal, atau belum.
37. Alamat penyimpanan di-resolve **runtime** dari konfigurasi; test yang menemukan URL sebagai
    literal, konstanta, atau env var **gagal**. *(**ADR-0013**)*
38. Lampiran dapat **diunduh** dan **dihapus**.
39. Token penyimpanan **di-cache** dan diperbarui sebelum kedaluwarsa; kedaluwarsanya tidak
    menggagalkan permintaan pengguna secara langsung.

**Skema dan migrasi**

40. ⚠️ Skema target **relasional penuh**: setiap atribut menjadi **kolom bernama**, setiap daftar
    bersarang menjadi **tabel anak**. Test yang menemukan kolom JSON sebagai penyimpan atribut
    produk **gagal**. `[keputusan work owner]`
41. Daftar kolom sisi inward memuat **keempat puluh field** yang di-set Activity simpan — termasuk
    `ADDENDUMNO`, `AMANDEMENTNO`, dan `CEDING` yang **hilang** dari catatan sumber.
42. Field yang hanya tampil di layar dengan **visibilitas mati** dan tidak di-set Activity **tidak**
    menjadi kolom.
43. ⚠️ **Seluruh kolom hasil migrasi NULLABLE.** Test yang menemukan `NOT NULL` pada kolom hasil
    migrasi **gagal**. `[keputusan work owner]`
44. ⚠️ Field yang **absen** di dokumen JSON lama menjadi **kolom kosong**, bukan kegagalan migrasi.
    `[keputusan work owner]`
45. Migrasi mengurai **seluruh daftar bersarang** menjadi baris tabel anak; tidak ada daftar yang
    tertinggal di dalam dokumen.
46. Nilai uang pindah **tanpa berubah satu digit pun**; rekonsiliasi membandingkan **secara tepat**,
    bukan dengan toleransi.
47. Sequence pindah dengan **nilai berjalan yang benar**, sehingga identitas baru tidak bertabrakan
    dengan yang lama.
48. Migrasi dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**.
49. Skema uji yang dipakai seluruh test dibangun **dari DDL yang sama** dengan produksi.

**Salah ketik dan cabang mati**

50. ⚠️ **Tidak ada `PoductName`** di kode baru — hanya `ProductName`. `[keputusan work owner]`
51. ⚠️ **Tidak ada `POLICYHODER`** sebagai nama kolom — hanya `POLICYHOLDER`, meski data lama
    memakai ejaan yang salah. Migrasi memetakannya. `[keputusan work owner]`
52. ⚠️ **Tidak ada `IsORS`** maupun cabang yang bergantung padanya. `[keputusan work owner]`
53. Tidak ada padanan jalur simpan pintas yang memperbarui sebagian kolom di luar jalur utama.

---

## Pertanyaan terbuka di dalam spec

**Nol OQ pemblokir.**

| OQ | Status | Pemilik | Catatan |
| --- | --- | --- | --- |
| **OQ-002** | ✅ **DITUTUP** `[data DBA]` | — | Body tiga procedure diterima: kedua penulis produk = **upsert JSON, `COMMIT` sendiri**, `StsSave` 100/99; `GET_TOKEN_STORAGE` = MD5, cache, 1 menit |
| **OQ-001** (modul ini) | ✅ **DITUTUP** `[data DBA]` | — | DDL dua tabel diterima; **tanpa PK**, hanya index |
| **OQ-JSON-PRODUCT** | ✅ **DITUTUP** `[data DBA]` | — | Isi `JSONDATA` terurai dari contoh produk nyata; daftar kolom lengkap diturunkan dari Activity |
| **OQ-056** | ✅ **DITUTUP** `[keputusan work owner]` | — | Jalur treaty inward = **salah ekspor**; modul ini tidak menulis `M_TREATY_IN` |
| **OQ-057** | ✅ **DITUTUP** | — | `OR` = `10200`; pemakainya di modul ini **mati** di balik gerbang `1==2` |
| **OQ-058** | ✅ **DITUTUP** `[keputusan work owner]` | — | `GetCountClaim` / `PC_ASM_FW_GCNMFW_WORK` tidak dipakai |
| **OQ-047** | `[terbuka]` — **tidak memblokir** | — | Alamat fisik Google Storage dari `M_LINK_SERVICE`; **ADR-0013** mengatur cara resolusinya |
| **OQ kecil** | `[terbuka]` — **tidak memblokir** | Product+UW | Kepanjangan "RI Rate" dan "RI Risk" |
| **OQ kecil** | `[terbuka]` — **tidak memblokir** | — | Field baru yang muncul saat migrasi nyata (skema evolutif) |

⚠️ **Selisih daftar kolom yang perlu diketahui, bukan ditebak.** Catatan sumber menyebut "28 field"
lalu mendaftar 37; sensus korpus menghasilkan **40**. Spec ini memakai **40** sesuai aturan pemilihan
kolom (Activity sebagai sumber kebenaran). Bila tim menghendaki tiga field tambahan itu dikecualikan,
itu **keputusan**, bukan temuan.

---

## Out of Scope

### Empat jalur yang **tidak dimigrasikan** — dikonfirmasi mati

| Jalur | Bukti | OQ |
| --- | --- | --- |
| **Seluruh jalur treaty inward** — `SetTreatyIn_Act` (`DATA-PORTAL` / `SETTREATYIN_ACT`), `SaveTreatyIn` (`ASM-FW-GISFW-INT-TREATY_IN` / `ASM!SAVETREATYIN`), `BrowseTreatyIn`, gerbang `revisionstate==1`, `POOLDATA.PEGA_TREATY_IN`, `M_TREATY_IN` / `M_TREATY_IN_EDM` | `[keputusan work owner]` **salah ekspor XML** — bukan fungsi modul ini | OQ-056 |
| **Jalur klaim** — `GetCountClaim` dan akses langsung ke `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` | `[keputusan work owner]` tidak dipakai | OQ-058 |
| **Jalur `OR`** — `GetReinsTypeOR_Life` (`ASM-FW-GISFW-INT-PRODUCT_LIFE` / `GETREINSTYPEOR_LIFE`), `BrowseReinstypeOR_SQL`, bendera `IsORS` | `[terverifikasi]` pemanggilnya dibungkus container bergerbang `pyContainerVisibleWhen` = `1==2` — **selalu salah** | OQ-057 |
| **Jalur simpan pintas** — `SaveProductNameLIfeFlat` (`ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!SAVEPRODUCTNAMELIFEFLAT`), `UPDATE` parsial atas `RIRISKID`/`RIRISK` | `[keputusan work owner]` tidak dipakai; pembaruan RI Risk menjadi bagian pembaruan produk biasa | — |

⚠️ **Pelajaran OQ-066, berlaku berkali-kali di modul ini: berkas ikut terekspor ≠ dipakai.** Empat
sumber kompleksitas terbesar ternyata mati. Tanpa konfirmasi work owner dan tanpa gerbang `1==2`
yang terbaca, keempatnya akan terbangun ulang tanpa guna.

### Lainnya

- **Master Contract Retro Life — seluruhnya.** Konteks/menu terpisah, sudah selesai.
- **Claim Life, PremiumList Life, Endorsement Life.** Konteks hilir yang **membaca** master ini.
- **Master pendukung itu sendiri** — ceding, pemegang polis, mata uang, sumber bisnis, penyebab
  klaim, RI Rate, RI Risk. Modul ini **memilih dari** mereka, tidak mengelolanya.
- **Lini non-Life.**
- **Identity & Access.** `[terverifikasi]` Guard identitas hanya muncul di tiga berkas; model peran
  dirancang terpisah (OQ-007 dan keluarganya).
- **Scaffolding kode.** Pekerjaan terpisah yang mendahului tiket mana pun.

---

## Further Notes

**Ukuran pekerjaan.** `[terverifikasi]` **114 berkas**: 37 Activity, 22 Connect-SQL,
17 ReportDefinition, 12 Section, 11 FlowAction, 10 DataTransform, 1 When, 1 Harness,
1 DecisionTable, 1 ConnectREST, 1 SystemSettings. **Nol rule `Flow`.**

⚠️ **Modul ini jauh lebih ramping dari angkanya.** Setelah empat jalur mati dikeluarkan, yang
tersisa adalah **editor master produk dengan lampiran** — satu layar, **satu tabel induk**
`product_life` (dua tabel existing digabung), lima tabel anak, tujuh pemilih master.

**Mengapa konteks ini menyimpang dari pola empat konteks Life sebelumnya.** Di Claim Life, Komite,
PremiumList Life, Endorsement Life, dan Master Contract Retro Life, keputusannya selalu **"panggil
procedure apa adanya"**. Di sini **tidak** — dan alasannya tunggal:

> `[data DBA]` Kedua procedure penulis **`COMMIT` sendiri**. Selama keduanya dipakai, "satu produk
> tersimpan bersama atau tidak sama sekali" **mustahil**. Work owner memilih atomik; maka procedure
> harus dilepas.

Keputusan itu menarik keputusan kedua: bila Go menulis sendiri, menulis **JSON** tidak lagi masuk
akal — sehingga skema menjadi **relasional penuh**. D1 dan D2 adalah **satu keputusan dengan dua
akibat**, bukan dua keputusan terpisah.

**Nama dan bentuk yang menyesatkan di korpus** — dicatat agar tidak ditiru:

| Yang tertulis | Yang sebenarnya |
| --- | --- |
| `PoductName` (20×, hanya di dua activity simpan) | salah ketik `ProductName`; menjaga langkah simpan |
| `POLICYHODER` | salah ketik `POLICYHOLDER` — **sudah masuk data produksi** |
| `IsORS` dibandingkan dua cara berbeda | bendera yang tak pernah di-set; cabangnya mati |
| `GetReinsTypeOR_Life` (66 KB, menjembatani tiga class) | mati di balik gerbang `1==2` |
| `SetTreatyIn_Act` (159 KB, 15 langkah, jalur tulis nyata) | **salah ekspor** — bukan fungsi modul ini |
| `ErrMsg` terisi **juga saat sukses** | bukan penanda kegagalan; `StsSave` yang menentukan |
| "28 field" pada catatan sumber | **40** menurut sensus Activity |

**Baca kodenya, jangan namanya** — dan di sini juga: **hitung sendiri, jangan percaya rekapnya.**

**Urutan yang saya sarankan untuk `/to-tickets`** — vertical slice:

1. Skema relasional + migrasi JSON → kolom (**prefactor**: seluruh slice lain berdiri di atasnya).
2. Produk: CRUD sisi umum + identitas dari sequence + gagal terang-terangan.
3. Sisi inward (kolom pada baris produk yang sama) + **simpan atomik** (baris produk + tabel anak).
4. Tujuh pemilih master (RI Rate dan RI Risk terpisah tegas).
5. Validasi lima field wajib, seragam di kedua sisi.
6. Daftar plan + batas underwriting.
7. Daftar komentar, dokumen klaim, lien, underwriting finansial, outward.
8. Lampiran: unggah, unduh, hapus, status.
9. Efek keluar Google Storage: resolusi alamat, token, kegagalan yang dapat diulang.

⚠️ **Slice 1 adalah prefactor**, bukan tiket migrasi yang biasanya terakhir. Karena skema berubah
bentuk (JSON → relasional), tidak ada slice lain yang dapat berdiri sebelum bentuk barunya ada.

**Keputusan yang masih terbuka, tidak memblokir:** pembagian paket domain di dalam `internal/` —
sama seperti lima spec sebelumnya.

**Catatan sumber.** Spec ini bersandar pada korpus Pega `D:\XML\RNM_BRD\` (READ-ONLY), artefak di
`OUTPUT_HASIL_RNM\`, dan berkas `[data DBA]` yang diterima 2026-09-15. Sumber ADR tunggal:
`docs/adr/ADR-0001`…`ADR-0015`.

---

## Ralat bertanggal 01-10-2026 — sesi implementasi (paket 0)

> Sumber: `RALAT-DEV-30-09-2026.md` (P1–P6 katalog DEV, R7–R18 pembacaan ulang XML) dan `PARITAS-LAYAR-DAN-AKSI.md`. Kalimat di atas
> **tidak dihapus**; yang berlaku adalah ralat ini. Register OQ: `OQ-MASTER-PRODUCT-NAME-LIFE.md`.

| Bagian spec *(dikutip)* | Ralat |
| --- | --- |
| Penyimpangan sadar 1b, 2; §3 *"satu produk, SATU tabel induk, lima tabel anak"*; §7 *"Skema relasional penuh"*; AC 40–49 | **P1** — ditangguhkan (OQ-MPNL-01): JSON di dua tabel lama seperti Pega, nol tabel baru, nol DDL. AC 40–49 tidak berlaku |
| Penyimpangan sadar 1 *"Satu transaksi atomik"*; §5 | tetap (P4) — atas dua tabel lama: `M_PRODUCT_LIFE` + `M_PRODUCTINWARD_LIFE` dalam satu transaksi; prosedur tidak dipanggil |
| Penyimpangan sadar 5 *"Salah ketik dibuang (`POLICYHODER` → `POLICYHOLDER`) dan cabang mati `IsORS` dibuang"*; AC 51–52 | kunci JSON tetap ejaan Pega (`POLICYHODER`) — dibaca tiga view; `IsORS` **hidup** (R9). `PoductName` tetap tidak dipakai, tetapi bukan karena perilakunya: prakondisi bernama itu **PRE=false**, tak pernah dievaluasi (R7) |
| §8 *"`SaveProductName_Act` memeriksa lima hal … Tipe & grup terisi"*; AC 11; §9 *"guard pada langkah simpan … akan mulai menolak"* | **R7**: langkah 1 adalah `Property-Set` ber-PRE=false, medan `TYPE`/`GRUP` mati (`1=2`). Validasi = empat pesan VERBATIM `Product Name Empty`, `Ceding Empty`, `Policy Holder Empty`, `SOB Empty` (R13). AC 11 tidak berlaku |
| §8 *"Di `SaveInwardProductName_Act`, dua pemeriksaan padanannya ter-remark … menegakkan validasi yang sama di kedua sisi"*; AC 16 | **R11**: jalur inward tak terjangkau (`Inward` b75368 `1=2`) — satu jalur simpan, satu aturan |
| Out of Scope *"Jalur `OR` … mati di balik gerbang yang selalu salah"*; *"Jalur simpan pintas — `SaveProductNameLIfeFlat` … tidak dipakai"* | **R9**: jalur OR hidup lewat checkbox `On Retention`; **R8**: `SaveProductNameLIfeFlat` jalan setiap simpan (PRE=false) |
| §1 *"Titik masuk `Master Product Name Life/Harness/InwardProductName.xml`"* | **R11**: halaman awal = `Section/InboxProductName.xml` |
| §3 *"`LienClause` = field skalar"*, *"`OutwardList` = dead code"* | **R10**, **P3** |
| §6 *"Sumbernya sequence `M_PRODUCT_LIFE_SEQ` dan `M_PRODUCT_INWARD_LIFE_SEQ`"* | **P6 / R14**: `ID` inward = `ID` produk (OQ-MPNL-02) |
| §4 *"RI Rate … master RI Rate"*; AC 19 | **R18**: sumber RI Rate menunggu OQ-MPNL-03 (503 berkalimat) |
| §10 *"Alamat penyimpanan di-resolve runtime"* | **P5**: pengiriman berkas = stub outbox; alamat nyata tidak dipanggil dan tidak ditulis (OQ-MPNL-10, OQ-MPNL-11) |

---

## Ralat bertanggal 01-10-2026 — lanjutan 1 (katalog DEV)

> Sumber: brief `PROMPT-LANJUTAN-MASTER-PRODUCT-NAME-LIFE-1.md` §1–§2 (katalog DEV `ALL_TAB_COLUMNS`/`ALL_OBJECTS` dan agregat `JSONDATA`, dibaca asisten, baca-saja). Kalimat di atas tidak dihapus.

| Kalimat lama | Ralat |
| --- | --- |
| *"`JSONDATA` (dengan constraint `IS JSON`) + empat kolom hasil flatten: `RIRISKID`, `RIRISK`, `PRODUCTNAME`, `BEGIN_DATE`"* (tabel induk `M_PRODUCT_LIFE`) | DEV `M_PRODUCT_LIFE` hanya `ID`, `JSONDATA`, `RIRISKID`, `RIRISK` — **dua** kolom datar; `PRODUCTNAME`/`BEGIN_DATE` tidak ada dan tidak ditulis (OQ-MPNL-08 ditutup, `6fd539c`) |

## Ralat bertanggal 02-10-2026 — penyimpanan FLAT (§3, §7) `[keputusan work owner 02-10-2026]`

> Sumber: brief `PROMPT-PINDAH-FLAT-MASTER-PRODUCT-NAME-LIFE.md` dan tiket 01 bab *"Keputusan bertanggal 02-10-2026"*. Kalimat lama
> dikutip; ralat 01-10-2026 (P1, *"Nol tabel baru, nol DDL"*) **dicabut untuk penyimpanan produk**.

| Kalimat lama | Ralat |
| --- | --- |
| §3 *"Entitas — satu produk, SATU tabel induk, lima tabel anak"* | satu produk = **satu** baris `M_PRODUCTNAME_LIFE` + **tujuh** tabel anak `M_PRODUCTNAME_LIFE_{LIEN,DOCCLAIM,PLAN,FINUW,UWLIMIT,OUTWARD,COMMENT}` (FK `ON DELETE CASCADE`, PK `PRODUCTID`+`URUT`) |
| §7 tipe *"desimal presisi arbitrer"* | `NUMBER(38,8)` untuk uang/persen/rate/faktor, `NUMBER(5)` untuk bilangan kecil — penjaga inti `TestNolNumberTanpaPresisi`, jawaban work owner 02-10-2026 (K6). Tetap nol float |
| §1/§7 *"tiga view … membaca"* | **ketiga view tidak dibangun ulang** (K7, jawaban work owner 02-10-2026): modul ini menulis dan membaca tabel flat saja. Pembaca hilir view `PRODUCTINWARD_LIFE` (Claim Life) → OQ-FLAT-04 |
| §7 *"Aturan migrasi"* | data pindah lewat alat Go `backend/alat/pindahflat` (`-uji` / `-jalankan`), rekonsiliasi teks demi teks; tabel JSON tetap ada sebagai cadangan, tidak pernah disentuh |

