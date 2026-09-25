# Penelusuran jalur `JSONDATA` ke kolom skema baru — Treaty In

**Tanggal:** 24 September 2026 · **Langkah 1 sesi to-spec**
**Dibangkitkan** `alat/telusur-jalur-ke-kolom.py` + `alat/tulis-penelusuran.py`.
Jangan disunting dengan tangan.

> ## SEMESTA BERKAS INI — dibaca sebelum satu baris pun dipercaya
>
> **Klaim *"seluruh jalur JSON punya rumah"* DILARANG di berkas ini.** Yang dinyatakan:
> *seluruh jalur **di dalam semesta yang diperiksa** punya nasib, dan semestanya kurang
> sebanyak angka di bawah.*
>
> | Fakta | Angka |
> |---|---:|
> | jalur beradjudikasi di `PETA-TELUSUR-JSON.md` §6 | **667** |
> | **dikecualikan** — pohon cermin, milik modul Adjustment | **272** |
> | **lingkup berkas ini** | **395** |
> | titik buta penjaga `primary_ok` (`L-8`) | **414** properti |
> | sudah diperiksa untuk 27 entitas gelombang 1 | 74, yang **41** diadili |
> | **belum diperiksa siapa pun** | **340** |
>
> Pohon cermin yang dikecualikan: `ActualValue`, `ValueDifference`, `OLDDATA`,
> `ValueBeforeProrate`. `struktur-treatyin-lama.md` §5 menghitungnya **526 dari 985 simpul**
> (217 + 152 + 142 + 15). **Penyebut 395 di atas dihitung ulang dari daftar jalur**, bukan
> diwarisi dari hitungan simpul — keduanya cara menghitung yang berbeda dan **tidak harus sama**.
>
> **`PETA-TELUSUR-JSON.md` sendiri berlabel *"semestanya kurang"***; lihat
> `4-erd-dan-tabel-datar/AUDIT-PENYEBUT-POHON.md`.

> ### BATAS PERKAKAS INI, dinyatakan di dalam keluarannya sendiri
>
> Ia dapat menunjukkan sebuah jalur **punya nasib tertulis**, dan untuk jalur DIPETAKAN ia dapat
> menunjuk kolom yang kolom *Asal*-nya menyebut ruas terakhir jalur itu.
>
> Ia **tidak** dapat membuktikan pemetaan itu **benar**. Pencocokannya **leksikal**: dua ruas
> bernama sama di kelas berbeda akan tercocokkan ke kolom yang sama. Setiap baris bertanda
> **COCOK-GANDA** adalah **calon**, bukan putusan — yang ditampilkan hanya kandidat pertama.

---

## 1. Hitungan penutup

| Nasib | Jumlah |
|---|---:|
| DIPETAKAN·cocok tunggal | 99 |
| DIPETAKAN·cocok ganda | 61 |
| DIPETAKAN·tidak tercocokkan | 53 |
| DITURUNKAN | 166 |
| DIBUANG | 12 |
| DITUNDA | 4 |
| **JUMLAH lingkup** | **395** |
| dikecualikan — pohon cermin | 272 |
| **JUMLAH semesta §6** | **667** |

**Menutup:** 395 + 272 = 667, sama dengan penyebut §6.

**Tidak ada jalur tanpa nasib.**

**Penolakan perkakas:** tidak ada.

---

## 2. Jalur DIPETAKAN yang tidak tercocokkan — 53, digolongkan bersebab

**Tidak tercocokkan bukan berarti tidak punya rumah.** Enam dari delapan golongan punya alasan
tertulis di berkas induk; dua terakhir yang menuntut pekerjaan.

| Golongan | Jumlah | Artinya |
|---|---:|---|
| **GEL-2** | 17 | cabang fakultatif keluar — `RETRO_KELUAR`, **di luar gelombang 1** |
| **TABEL** | 13 | jalur menamai **daftarnya**, bukan sebuah ruas — ia menjadi **tabel** |
| **NAMA-MASTER** | 7 | nama dibaca dari master; hanya pengenalnya disimpan (**ADR-0041**) |
| **DIBUANG-ULANG** | 4 | sudah diputuskan dibuang atau melebur di §12.5, §13, §14.5 |
| **MELEBUR-B** | 4 |  |
| **BERSUSUN** | 4 | daftar kelas bisnis / kelompok treaty bersusun — **belum punya rumah bernama** |
| **AGREGAT** | 2 | berpasangan dengan besaran **agregat** yang tidak disimpan (§4.2) |
| **YATIM** | 1 | **tidak ada rumah, dan tidak ada alasan tertulis** |
| **TURUNAN-F2** | 1 |  |

### 2.1 Dua jalur YATIM — dan keduanya menulis ke tabel datar setiap simpan

* **`NusareSharePct`**

**`NusareSharePct`** sudah tercatat di `SPEC-MODEL-DATA.md` §3.2a: **nol penulis** atas kelima
bentuk, dan ia menulis ke `TREATY_IN.NUSARESHAREPCT` pada setiap penyimpanan — kolomnya
**selalu kosong**.

**`TreatyYear` adalah temuan baru berkas ini.** Ia bertanda **DIPETAKAN** di
`PETA-TELUSUR-JSON.md`, **tidak punya satu pun kolom** di skema baru — tidak ada `TAHUN_TREATY`
maupun padanan lain — dan ia **salah satu dari 20 kolom bisnis `TREATY_IN`**, ditulis setiap
penyimpanan lewat `POOLDATA.PEGA_TREATY_IN`.

> **Tidak ditambal di sini.** Tahun treaty **mungkin** turunan dari tanggal mulai kontrak, dan
> **mungkin** tahun *underwriting* yang disepakati terpisah. Keduanya lazim di pasar, dan
> keduanya menghasilkan angka yang **sama pada kebanyakan kontrak** — persis bentuk *"cacat
> bersembunyi di balik nilai bawaan"*. Menebaknya sebagai turunan akan benar untuk hampir semua
> baris, dan **salah tanpa terlihat** untuk sisanya.
>
> | | |
> |---|---|
> | **Siapa menutup** | **teknik treaty** — satu pertanyaan, diusulkan sebagai `T-7` |
> | **Yang menagih** | ketiadaan `TAHUN_TREATY` di `KAMUS-KOLOM.md`, dan berkas ini |

### 2.2 Empat jalur BERSUSUN — daftar yang belum punya rumah bernama

* `Limits[].Detail[].COBList[].ClassOfBusinessID`
* `Limits[].TreatyGroupList[].ClassOfBusinessList[].ClassOfBusinessID`
* `Share[].ClassofBusinessList`
* `Share[].TreatyGroupList`

Keempatnya daftar **kelas bisnis** atau **kelompok treaty** yang bersusun di bawah induknya.
`EGNPI` punya `ID_KELAS_BISNIS` dan `VERSI_KONTRAK` punya `KELAS_BISNIS_KONTRAK` — tetapi
**tidak ada tempat bagi daftar kelas bisnis per kelompok treaty per layer**.

> **Ini kandidat kekeliruan lingkup yang sudah tiga kali terjadi** (`CONTEXT.md` §2.2c): daftar
> yang tampak menggantung langsung pada induk besarnya, padahal ada satu tingkat pengelompokan
> di antaranya. **Lingkupnya harus dibaca dari aktivitas yang menambah barisnya**, bukan dari
> bentuk pohon. Dilaporkan; pemiliknya sesi to-spec berikutnya.

---

## 3. Tabel penelusuran — satu baris per jalur, 395 baris

`BUKTI` berbentuk `EVIDENCED(sumber@ekspor-2026-09)`. Ekspor ini **mungkin** dari lingkungan QA (`L-9`);
penanda itu yang memungkinkan **satu sapuan** mengeluarkan daftar terdampak bila jawabannya datang.

| JALUR_JSON | NASIB | TABEL | KOLOM | TIPE | BOLEH_KOSONG | ALASAN | BUKTI |
|---|---|---|---|---|---|---|---|
| `AccountingMode` | DIPETAKAN | VERSI_KONTRAK | CARA_PEMBUKUAN | E | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `AccountingModeNonProp` | DIPETAKAN | VERSI_KONTRAK | CARA_PEMBUKUAN_XOL | E | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `AccumulationList` | DIPETAKAN·TABEL | — | — | — | — | jalur ini menamai **daftarnya**, bukan sebuah ruas — ia menjadi **tabel** `PERIODE_AKUMULASI` | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `AccumulationList[]` | DIPETAKAN·TABEL | — | — | — | — | jalur ini menamai **daftarnya**, bukan sebuah ruas — ia menjadi **tabel** `PERIODE_AKUMULASI` | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `AccumulationList[].Period` | DIPETAKAN | PERIODE_PELAPORAN | PERIODE | T | tidak | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `AccumulationList[].ReportDate` | DIPETAKAN | PERIODE_AKUMULASI | TANGGAL_LAPOR | D | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `AccumulationPeriod` | DIPETAKAN | VERSI_KONTRAK | PERIODE_AKUMULASI_KONTRAK | E | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Bordeaux` | DIPETAKAN | VERSI_KONTRAK | MEMAKAI_BORDEREAUX | E | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `BordereauxNote` | DIPETAKAN | VERSI_KONTRAK | CATATAN_BORDEREAUX | T | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `BrokeragePercent` | DIPETAKAN | VERSI_KONTRAK | PERSEN_BROKERAGE | P2 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `BrokeragePercentP` | DIPETAKAN·DIBUANG-ULANG | — | — | — | — | §13 — akhiran P = proportional; MELEBUR ke PERSEN_BROKERAGE | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Ceding` | DIPETAKAN·NAMA-MASTER | — | — | — | — | nama dibaca dari master; hanya pengenalnya disimpan (ADR-0041) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CedingID` | DIPETAKAN | KONTRAK | ID_CEDANT | R | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CedingStatusActive` | DIPETAKAN·DIBUANG-ULANG | — | — | — | — | §14.5 — bendera status master, bukan fakta kontrak | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ClassofBusiness` | DIPETAKAN | VERSI_KONTRAK | KELAS_BISNIS_KONTRAK | R | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CoInScale` | DIPETAKAN·TABEL | — | — | — | — | jalur ini menamai **daftarnya**, bukan sebuah ruas — ia menjadi **tabel** `SKALA_KOASURANSI` | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Commencement` | DIPETAKAN | KONTRAK | TANGGAL_MULAI | D | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Comment` | DIPETAKAN | VERSI_KONTRAK | CATATAN | T | ya | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CommentList` | DIPETAKAN | PERISTIWA_KONTRAK | WAKTU | D | tidak | tercocokkan **cocok-ganda(3)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CommentList[]` | DIPETAKAN | PERISTIWA_KONTRAK | WAKTU | D | tidak | tercocokkan **cocok-ganda(3)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CommentList[].Date` | DIPETAKAN | CATATAN_PERSETUJUAN | WAKTU_KEPUTUSAN | D | tidak | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CommentList[].IsApproved` | DIPETAKAN | CATATAN_PERSETUJUAN | DISETUJUI | E | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CommentList[].OperatorName` | DIPETAKAN | CATATAN_PERSETUJUAN | NAMA_PEMUTUS | T | tidak | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CommentList[].Suggest` | DIPETAKAN | CATATAN_PERSETUJUAN | ALASAN | T | ya | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Currency` | DIPETAKAN | VERSI_KONTRAK | KODE_MATA_UANG_KONTRAK | R | tidak | tercocokkan **cocok-ganda(17)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CurrencyEarthquake` | DIPETAKAN·MELEBUR-B | — | — | — | — | mata uang batas per bahaya — **melebur** ke `BATAS_PER_BAHAYA.KODE_MATA_UANG` (`KTV-C` koreksi 1) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CurrencyEgnpiAmount` | DIPETAKAN·AGREGAT | — | — | — | — | berpasangan dengan `TotalEgnpiAmount` yang **agregat** dan tidak disimpan (§4.2) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CurrencyFloodJab` | DIPETAKAN·MELEBUR-B | — | — | — | — | mata uang batas per bahaya — **melebur** ke `BATAS_PER_BAHAYA.KODE_MATA_UANG` (`KTV-C` koreksi 1) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CurrencyFloodNat` | DIPETAKAN·MELEBUR-B | — | — | — | — | mata uang batas per bahaya — **melebur** ke `BATAS_PER_BAHAYA.KODE_MATA_UANG` (`KTV-C` koreksi 1) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CurrencyInstallmentAmount` | DIPETAKAN·AGREGAT | — | — | — | — | berpasangan dengan `TotalInstallmentAmount` yang **agregat** dan tidak disimpan (§4.2) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CurrencyList` | DIPETAKAN·TABEL | — | — | — | — | jalur ini menamai **daftarnya**, bukan sebuah ruas — ia menjadi **tabel** `MATA_UANG_KONTRAK` | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CurrencyList[]` | DIPETAKAN·TABEL | — | — | — | — | jalur ini menamai **daftarnya**, bukan sebuah ruas — ia menjadi **tabel** `MATA_UANG_KONTRAK` | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CurrencyList[].Conversion` | DIPETAKAN | MATA_UANG_KONTRAK | KURS | K1 | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CurrencyList[].Currency` | DIPETAKAN | VERSI_KONTRAK | KODE_MATA_UANG_KONTRAK | R | tidak | tercocokkan **cocok-ganda(17)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CurrencyList[].CurrencyID` | DIPETAKAN | LAYER | KODE_MATA_UANG | T | tidak | tercocokkan **cocok-ganda(10)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CurrencyList[].PeriodEnd` | DIPETAKAN | MATA_UANG_KONTRAK | TANGGAL_AKHIR_BERLAKU | D | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CurrencyList[].PeriodStart` | DIPETAKAN | MATA_UANG_KONTRAK | TANGGAL_MULAI_BERLAKU | D | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `CurrencyRSMD` | DIPETAKAN·MELEBUR-B | — | — | — | — | mata uang batas per bahaya — **melebur** ke `BATAS_PER_BAHAYA.KODE_MATA_UANG` (`KTV-C` koreksi 1) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `EDMEffective` | DIPETAKAN | VERSI_KONTRAK | TANGGAL_BERLAKU_ADDENDUM | D | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `EDMState` | DIPETAKAN | VERSI_KONTRAK | JENIS_ADDENDUM | E | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `EGNPI` | DIPETAKAN·TABEL | — | — | — | — | jalur ini menamai **daftarnya**, bukan sebuah ruas — ia menjadi **tabel** `EGNPI` | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `EGNPI[]` | DIPETAKAN·TABEL | — | — | — | — | jalur ini menamai **daftarnya**, bukan sebuah ruas — ia menjadi **tabel** `EGNPI` | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `EGNPI[].Amount` | DIPETAKAN | NILAI_PENYEBARAN | NILAI | U1 paket uang | tidak | tercocokkan **cocok-ganda(4)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `EGNPI[].AsDate` | DIPETAKAN | EGNPI | TANGGAL_BERLAKU | D | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `EGNPI[].Currency` | DIPETAKAN | VERSI_KONTRAK | KODE_MATA_UANG_KONTRAK | R | tidak | tercocokkan **cocok-ganda(17)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `EGNPI[].CurrencyID` | DIPETAKAN | LAYER | KODE_MATA_UANG | T | tidak | tercocokkan **cocok-ganda(10)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `EGNPI[].ID` | DIPETAKAN | KONTRAK | NOMOR_KONTRAK_WARISAN | T | ya | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `EGNPI[].TreatyGroup` | DIPETAKAN·NAMA-MASTER | — | — | — | — | nama dibaca dari master; hanya pengenalnya disimpan (ADR-0041) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `EGNPI[].TreatyGroupID` | DIPETAKAN | DETAIL_PROPORSIONAL | ID_KELOMPOK_TREATY | R | tidak | tercocokkan **cocok-ganda(3)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Earthquake` | DIPETAKAN | BATAS_PER_BAHAYA | NILAI_BATAS | U1 paket uang | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Exclusions` | DIPETAKAN | VERSI_KONTRAK | PENGECUALIAN | T | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ExclusionsP` | DIPETAKAN | VERSI_KONTRAK | PENGECUALIAN | T | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacShare` | DIPETAKAN | VERSI_KONTRAK | PERSEN_BAGIAN_FAKULTATIF | P1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacShareBrokerage` | DIPETAKAN | VERSI_KONTRAK | PERSEN_BROKERAGE_FAKULTATIF | P2 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShare` | DIPETAKAN | VERSI_KONTRAK | PERSEN_BAGIAN_FAKULTATIF | P1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareBrokerage` | DIPETAKAN | VERSI_KONTRAK | PERSEN_BROKERAGE_FAKULTATIF | P2 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList` | DIPETAKAN·GEL-2 | — | — | — | — | cabang fakultatif keluar — `RETRO_KELUAR`, di luar gelombang 1 (§2.3) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[]` | DIPETAKAN·GEL-2 | — | — | — | — | cabang fakultatif keluar — `RETRO_KELUAR`, di luar gelombang 1 (§2.3) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].ClassofBusinessList` | DIPETAKAN·GEL-2 | — | — | — | — | cabang fakultatif keluar — `RETRO_KELUAR`, di luar gelombang 1 (§2.3) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].Cover` | DIPETAKAN | LAYER | CAKUPAN | T | ya | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].GrossPremiumList[].Currency` | DIPETAKAN | VERSI_KONTRAK | KODE_MATA_UANG_KONTRAK | R | tidak | tercocokkan **cocok-ganda(17)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].GrossPremiumList[].Value` | DIPETAKAN | NILAI_PREMI_BRUTO | NILAI | U1 | tidak | tercocokkan **cocok-ganda(5)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].Layer` | DIPETAKAN | LAYER | NOMOR_LAYER | C1 | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].LayerPart` | DIPETAKAN | LAYER | BAGIAN_LAYER | C1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].LayerPartType` | DIPETAKAN | LAYER | JENIS_BAGIAN_LAYER | E | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].LayerType` | DIPETAKAN | LAYER | JENIS_LAYER | E | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].Limit` | DIPETAKAN | LAYER | LIMIT | U1 paket uang | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].ShareFacultativeReinsurers[].Amount` | DIPETAKAN | NILAI_PENYEBARAN | NILAI | U1 paket uang | tidak | tercocokkan **cocok-ganda(4)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].ShareFacultativeReinsurers[].Amount2` | DIPETAKAN·GEL-2 | — | — | — | — | cabang fakultatif keluar — `RETRO_KELUAR`, di luar gelombang 1 (§2.3) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].ShareFacultativeReinsurers[].ReinsID` | DIPETAKAN | RINCIAN_PENYEBARAN | ID_JENIS_REASURANSI | R | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].ShareFacultativeReinsurers[].ReinsName` | DIPETAKAN·GEL-2 | — | — | — | — | cabang fakultatif keluar — `RETRO_KELUAR`, di luar gelombang 1 (§2.3) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].ShareFacultativeReinsurers[].SharePct` | DIPETAKAN | RINCIAN_PENYEBARAN | PERSEN_RINCIAN | P1 | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].TreatyGroupList` | DIPETAKAN·GEL-2 | — | — | — | — | cabang fakultatif keluar — `RETRO_KELUAR`, di luar gelombang 1 (§2.3) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FloodJab` | DIPETAKAN | BATAS_PER_BAHAYA | NILAI_BATAS | U1 paket uang | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FloodNation` | DIPETAKAN | BATAS_PER_BAHAYA | NILAI_BATAS | U1 paket uang | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ID` | DIPETAKAN | KONTRAK | NOMOR_KONTRAK_WARISAN | T | ya | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Information` | DIPETAKAN | VERSI_KONTRAK | KETERANGAN | T | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Installment` | DIPETAKAN | TERMIN | NOMOR_TERMIN | C1 | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `InstallmentNo` | DIPETAKAN | VERSI_KONTRAK | JUMLAH_TERMIN | C1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Installment[]` | DIPETAKAN | TERMIN | NOMOR_TERMIN | C1 | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Installment[].Currency` | DIPETAKAN | VERSI_KONTRAK | KODE_MATA_UANG_KONTRAK | R | tidak | tercocokkan **cocok-ganda(17)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Installment[].ID` | DIPETAKAN | KONTRAK | NOMOR_KONTRAK_WARISAN | T | ya | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `IsMultipleRetro` | DIPETAKAN | VERSI_KONTRAK | RETRO_BERGANDA | E | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `IsProRate` | DIPETAKAN | VERSI_KONTRAK | MEMAKAI_PRORATA | E | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `LeadingReinsID` | DIPETAKAN | VERSI_KONTRAK | ID_REASURADUR_PEMIMPIN | R | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `LeadingReinsName` | DIPETAKAN·NAMA-MASTER | — | — | — | — | nama dibaca dari master; hanya pengenalnya disimpan (ADR-0041) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `LeadingReinsSource` | DIPETAKAN·NAMA-MASTER | — | — | — | — | nama dibaca dari master; hanya pengenalnya disimpan (ADR-0041) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `LeadingReinsSourceID` | DIPETAKAN | KONTRAK | ID_ASAL_BISNIS | R | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits` | DIPETAKAN | LAYER | ID_LAYER | C1 | tidak | tercocokkan **cocok-ganda(3)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[]` | DIPETAKAN | LAYER | ID_LAYER | C1 | tidak | tercocokkan **cocok-ganda(3)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].AdjRate` | DIPETAKAN | LAYER | PERSEN_PENYESUAIAN | P1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Cover` | DIPETAKAN | LAYER | CAKUPAN | T | ya | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Currency` | DIPETAKAN | VERSI_KONTRAK | KODE_MATA_UANG_KONTRAK | R | tidak | tercocokkan **cocok-ganda(17)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Deductible` | DIPETAKAN | LAYER | DEDUCTIBLE | U1 paket uang | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].COBList[].ClassOfBusiness` | DIPETAKAN | EGNPI | ID_KELAS_BISNIS | R | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].COBList[].ClassOfBusinessID` | DIPETAKAN·BERSUSUN | — | — | — | — | daftar kelas bisnis / kelompok treaty bersusun di bawah induknya — **belum punya rumah bernama** | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].CashLossList[].Currency` | DIPETAKAN | VERSI_KONTRAK | KODE_MATA_UANG_KONTRAK | R | tidak | tercocokkan **cocok-ganda(17)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].CashLossList[].Value` | DIPETAKAN | NILAI_PREMI_BRUTO | NILAI | U1 | tidak | tercocokkan **cocok-ganda(5)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].ClaimCoopList[].Currency` | DIPETAKAN | VERSI_KONTRAK | KODE_MATA_UANG_KONTRAK | R | tidak | tercocokkan **cocok-ganda(17)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].ClaimCoopList[].Value` | DIPETAKAN | NILAI_PREMI_BRUTO | NILAI | U1 | tidak | tercocokkan **cocok-ganda(5)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].EPIList[].Currency` | DIPETAKAN | VERSI_KONTRAK | KODE_MATA_UANG_KONTRAK | R | tidak | tercocokkan **cocok-ganda(17)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].EPIList[].Value` | DIPETAKAN | NILAI_PREMI_BRUTO | NILAI | U1 | tidak | tercocokkan **cocok-ganda(5)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].IOOLimitList[].Currency` | DIPETAKAN | VERSI_KONTRAK | KODE_MATA_UANG_KONTRAK | R | tidak | tercocokkan **cocok-ganda(17)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].IOOLimitList[].CurrencyID` | DIPETAKAN | LAYER | KODE_MATA_UANG | T | tidak | tercocokkan **cocok-ganda(10)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].IOOLimitList[].ID` | DIPETAKAN | KONTRAK | NOMOR_KONTRAK_WARISAN | T | ya | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].IOOLimitList[].Value` | DIPETAKAN | NILAI_PREMI_BRUTO | NILAI | U1 | tidak | tercocokkan **cocok-ganda(5)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].PLAList[].Currency` | DIPETAKAN | VERSI_KONTRAK | KODE_MATA_UANG_KONTRAK | R | tidak | tercocokkan **cocok-ganda(17)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].PLAList[].Value` | DIPETAKAN | NILAI_PREMI_BRUTO | NILAI | U1 | tidak | tercocokkan **cocok-ganda(5)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].QSPct` | DIPETAKAN | DETAIL_PROPORSIONAL | PERSEN_QUOTA_SHARE | P1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].RIOGR` | DIPETAKAN | DETAIL_PROPORSIONAL | PERSEN_KOMISI_KOTOR | P1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].RIONR` | DIPETAKAN | DETAIL_PROPORSIONAL | PERSEN_KOMISI_BERSIH | P1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].RetentionList[].Currency` | DIPETAKAN | VERSI_KONTRAK | KODE_MATA_UANG_KONTRAK | R | tidak | tercocokkan **cocok-ganda(17)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].RetentionList[].Value` | DIPETAKAN | NILAI_PREMI_BRUTO | NILAI | U1 | tidak | tercocokkan **cocok-ganda(5)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].Surplus` | DIPETAKAN | DETAIL_PROPORSIONAL | JUMLAH_LINES_SURPLUS | C1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].TreatyGroup` | DIPETAKAN·NAMA-MASTER | — | — | — | — | nama dibaca dari master; hanya pengenalnya disimpan (ADR-0041) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].TreatyGroupID` | DIPETAKAN | DETAIL_PROPORSIONAL | ID_KELOMPOK_TREATY | R | tidak | tercocokkan **cocok-ganda(3)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].TreatyType` | DIPETAKAN | DETAIL_PROPORSIONAL | JENIS_TREATY | E | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].ID` | DIPETAKAN | KONTRAK | NOMOR_KONTRAK_WARISAN | T | ya | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Layer` | DIPETAKAN | LAYER | NOMOR_LAYER | C1 | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].LayerPart` | DIPETAKAN | LAYER | BAGIAN_LAYER | C1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].LayerPartType` | DIPETAKAN | LAYER | JENIS_BAGIAN_LAYER | E | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].LayerType` | DIPETAKAN | LAYER | JENIS_LAYER | E | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Limit` | DIPETAKAN | LAYER | LIMIT | U1 paket uang | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].MDPPct` | DIPETAKAN | LAYER | PERSEN_MINIMUM_DEPOSIT | P1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].ReinstatementPct` | DIPETAKAN | LAYER | PORSI_PEMULIHAN_LIMIT | P1 | ya | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].ReinstatementValue` | DIPETAKAN | LAYER | TARIF_PREMI_PEMULIHAN | P1 | ya | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].TreatyGroupList[].ClassOfBusinessList[].ClassOfBusiness` | DIPETAKAN | EGNPI | ID_KELAS_BISNIS | R | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].TreatyGroupList[].ClassOfBusinessList[].ClassOfBusinessID` | DIPETAKAN·BERSUSUN | — | — | — | — | daftar kelas bisnis / kelompok treaty bersusun di bawah induknya — **belum punya rumah bernama** | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].TreatyGroupList[].TreatyGroup` | DIPETAKAN·NAMA-MASTER | — | — | — | — | nama dibaca dari master; hanya pengenalnya disimpan (ADR-0041) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].TreatyGroupList[].TreatyGroupID` | DIPETAKAN | DETAIL_PROPORSIONAL | ID_KELOMPOK_TREATY | R | tidak | tercocokkan **cocok-ganda(3)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].TreatyType` | DIPETAKAN | DETAIL_PROPORSIONAL | JENIS_TREATY | E | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `MaxCoGroup` | DIPETAKAN | VERSI_KONTRAK | BATAS_MAKSIMUM_KELOMPOK | U1 paket uang | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `MaxCoNonGroup` | DIPETAKAN | VERSI_KONTRAK | BATAS_MAKSIMUM_NON_KELOMPOK | U1 paket uang | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `NusareSharePct` | DIPETAKAN·YATIM | — | — | — | — | **tidak ada rumah, dan tidak ada alasan tertulis** | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `OptionLimit` | DIPETAKAN | VERSI_KONTRAK | BATAS_PILIHAN | U1 paket uang | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Portfolio` | DIPETAKAN·TABEL | — | — | — | — | jalur ini menamai **daftarnya**, bukan sebuah ruas — ia menjadi **tabel** `PORTOFOLIO` | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Portfolio[]` | DIPETAKAN·TABEL | — | — | — | — | jalur ini menamai **daftarnya**, bukan sebuah ruas — ia menjadi **tabel** `PORTOFOLIO` | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Portfolio[].Description` | DIPETAKAN | PORTOFOLIO | KETERANGAN | T | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Portfolio[].Type` | DIPETAKAN | PORTOFOLIO | ARAH_PORTOFOLIO | E | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Portfolio[].TypePortfolio` | DIPETAKAN | PORTOFOLIO | JENIS_PORTOFOLIO | E | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `PositionUsername` | DIPETAKAN·DIBUANG-ULANG | — | — | — | — | §12.5 — pemegang sekarang adalah turunan, bukan fakta | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ProportionType` | DIPETAKAN | KONTRAK | SIFAT_PROPORSI | E | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `RNMShare` | DIPETAKAN | VERSI_KONTRAK | PERSEN_BAGIAN_NURE | P1 | tidak | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `RNMShareAcrossTheBoard` | DIPETAKAN | VERSI_KONTRAK | BAGIAN_NURE_SERAGAM | E | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `RNMShareP` | DIPETAKAN | VERSI_KONTRAK | PERSEN_BAGIAN_NURE | P1 | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `RSMDLimit` | DIPETAKAN | BATAS_PER_BAHAYA | NILAI_BATAS | U1 paket uang | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ReminderDays` | DIPETAKAN | VERSI_KONTRAK | HARI_PENGINGAT | C1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ReportingConfirmation` | DIPETAKAN | VERSI_KONTRAK | HARI_BATAS_KONFIRMASI | C1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ReportingEnd` | DIPETAKAN | VERSI_KONTRAK | TANGGAL_AKHIR_PELAPORAN | D | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ReportingInterval` | DIPETAKAN | VERSI_KONTRAK | SELANG_PELAPORAN | C1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ReportingPeriod` | DIPETAKAN | VERSI_KONTRAK | PERIODE_PELAPORAN_KONTRAK | E | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ReportingPeriodList` | DIPETAKAN·TABEL | — | — | — | — | jalur ini menamai **daftarnya**, bukan sebuah ruas — ia menjadi **tabel** `PERIODE_PELAPORAN` | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ReportingPeriodList[]` | DIPETAKAN·TABEL | — | — | — | — | jalur ini menamai **daftarnya**, bukan sebuah ruas — ia menjadi **tabel** `PERIODE_PELAPORAN` | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ReportingPeriodList[].ConfirmationDue` | DIPETAKAN | PERIODE_PELAPORAN | BATAS_KONFIRMASI | D | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ReportingPeriodList[].InitialDate` | DIPETAKAN | PERIODE_PELAPORAN | TANGGAL_AWAL | D | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ReportingPeriodList[].Period` | DIPETAKAN | PERIODE_PELAPORAN | PERIODE | T | tidak | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ReportingPeriodList[].SettlementDue` | DIPETAKAN | PERIODE_PELAPORAN | BATAS_PELUNASAN | D | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ReportingPeriodList[].SubmissionDue` | DIPETAKAN | PERIODE_PELAPORAN | BATAS_PENYERAHAN | D | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ReportingSettlement` | DIPETAKAN | VERSI_KONTRAK | HARI_BATAS_PELUNASAN | C1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ReportingStart` | DIPETAKAN | VERSI_KONTRAK | TANGGAL_MULAI_PELAPORAN | D | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ReportingSubmission` | DIPETAKAN | VERSI_KONTRAK | HARI_BATAS_PENYERAHAN | C1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Retention` | DIPETAKAN·TABEL | — | — | — | — | jalur ini menamai **daftarnya**, bukan sebuah ruas — ia menjadi **tabel** `RETENSI_CEDANT` | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Retention[]` | DIPETAKAN·TABEL | — | — | — | — | jalur ini menamai **daftarnya**, bukan sebuah ruas — ia menjadi **tabel** `RETENSI_CEDANT` | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Retention[].Amount` | DIPETAKAN | NILAI_PENYEBARAN | NILAI | U1 paket uang | tidak | tercocokkan **cocok-ganda(4)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Retention[].Currency` | DIPETAKAN | VERSI_KONTRAK | KODE_MATA_UANG_KONTRAK | R | tidak | tercocokkan **cocok-ganda(17)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Retention[].CurrencyID` | DIPETAKAN | LAYER | KODE_MATA_UANG | T | tidak | tercocokkan **cocok-ganda(10)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Retention[].ID` | DIPETAKAN | KONTRAK | NOMOR_KONTRAK_WARISAN | T | ya | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Retention[].TreatyGroup` | DIPETAKAN·NAMA-MASTER | — | — | — | — | nama dibaca dari master; hanya pengenalnya disimpan (ADR-0041) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Retention[].TreatyGroupID` | DIPETAKAN | DETAIL_PROPORSIONAL | ID_KELOMPOK_TREATY | R | tidak | tercocokkan **cocok-ganda(3)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `RetroList` | DIPETAKAN | VERSI_KONTRAK | ID_DAFTAR_RETRO | R | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `RnmShareDeducted` | DIPETAKAN | VERSI_KONTRAK | PERSEN_BAGIAN_NURE_DIPOTONG | P1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share` | DIPETAKAN | BAGIAN | PERSEN_BAGIAN_NURE | P1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers` | DIPETAKAN·GEL-2 | — | — | — | — | cabang fakultatif keluar — `RETRO_KELUAR`, di luar gelombang 1 (§2.3) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[]` | DIPETAKAN·GEL-2 | — | — | — | — | cabang fakultatif keluar — `RETRO_KELUAR`, di luar gelombang 1 (§2.3) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].BrokerName` | DIPETAKAN·GEL-2 | — | — | — | — | cabang fakultatif keluar — `RETRO_KELUAR`, di luar gelombang 1 (§2.3) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].CessionList` | DIPETAKAN·GEL-2 | — | — | — | — | cabang fakultatif keluar — `RETRO_KELUAR`, di luar gelombang 1 (§2.3) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].IOOLimitList` | DIPETAKAN·GEL-2 | — | — | — | — | cabang fakultatif keluar — `RETRO_KELUAR`, di luar gelombang 1 (§2.3) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].QSPct` | DIPETAKAN | DETAIL_PROPORSIONAL | PERSEN_QUOTA_SHARE | P1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].RNMShare` | DIPETAKAN | VERSI_KONTRAK | PERSEN_BAGIAN_NURE | P1 | tidak | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].ShareNote` | DIPETAKAN·GEL-2 | — | — | — | — | cabang fakultatif keluar — `RETRO_KELUAR`, di luar gelombang 1 (§2.3) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].TreatyGroup` | DIPETAKAN·GEL-2 | — | — | — | — | cabang fakultatif keluar — `RETRO_KELUAR`, di luar gelombang 1 (§2.3) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].TreatyGroupID` | DIPETAKAN | DETAIL_PROPORSIONAL | ID_KELOMPOK_TREATY | R | tidak | tercocokkan **cocok-ganda(3)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].TreatyType` | DIPETAKAN | DETAIL_PROPORSIONAL | JENIS_TREATY | E | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].FacultativeLimits[].TreatyType` | DIPETAKAN | DETAIL_PROPORSIONAL | JENIS_TREATY | E | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].FacultativeLimits[].TreatyTypeID` | DIPETAKAN·GEL-2 | — | — | — | — | cabang fakultatif keluar — `RETRO_KELUAR`, di luar gelombang 1 (§2.3) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].ID` | DIPETAKAN | KONTRAK | NOMOR_KONTRAK_WARISAN | T | ya | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].Layer` | DIPETAKAN | LAYER | NOMOR_LAYER | C1 | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].ReinsID` | DIPETAKAN | RINCIAN_PENYEBARAN | ID_JENIS_REASURANSI | R | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].ReinsName` | DIPETAKAN·GEL-2 | — | — | — | — | cabang fakultatif keluar — `RETRO_KELUAR`, di luar gelombang 1 (§2.3) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].SharePct` | DIPETAKAN | RINCIAN_PENYEBARAN | PERSEN_RINCIAN | P1 | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareReins` | DIPETAKAN·GEL-2 | — | — | — | — | cabang fakultatif keluar — `RETRO_KELUAR`, di luar gelombang 1 (§2.3) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareReins[]` | DIPETAKAN·GEL-2 | — | — | — | — | cabang fakultatif keluar — `RETRO_KELUAR`, di luar gelombang 1 (§2.3) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareReins[].ID` | DIPETAKAN | KONTRAK | NOMOR_KONTRAK_WARISAN | T | ya | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[]` | DIPETAKAN | BAGIAN | PERSEN_BAGIAN_NURE | P1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].ClassofBusinessList` | DIPETAKAN·BERSUSUN | — | — | — | — | daftar kelas bisnis / kelompok treaty bersusun di bawah induknya — **belum punya rumah bernama** | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].Cover` | DIPETAKAN | LAYER | CAKUPAN | T | ya | tercocokkan **cocok-ganda(2)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].GrossPremiumList[].Currency` | DIPETAKAN | VERSI_KONTRAK | KODE_MATA_UANG_KONTRAK | R | tidak | tercocokkan **cocok-ganda(17)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].GrossPremiumList[].Value` | DIPETAKAN | NILAI_PREMI_BRUTO | NILAI | U1 | tidak | tercocokkan **cocok-ganda(5)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].Layer` | DIPETAKAN | LAYER | NOMOR_LAYER | C1 | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].LayerPart` | DIPETAKAN | LAYER | BAGIAN_LAYER | C1 | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].LayerPartType` | DIPETAKAN | LAYER | JENIS_BAGIAN_LAYER | E | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].LayerType` | DIPETAKAN | LAYER | JENIS_LAYER | E | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].GrossPremiumList[].Currency` | DIPETAKAN | VERSI_KONTRAK | KODE_MATA_UANG_KONTRAK | R | tidak | tercocokkan **cocok-ganda(17)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].GrossPremiumList[].Value` | DIPETAKAN | NILAI_PREMI_BRUTO | NILAI | U1 | tidak | tercocokkan **cocok-ganda(5)** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].TreatyGroupList` | DIPETAKAN·BERSUSUN | — | — | — | — | daftar kelas bisnis / kelompok treaty bersusun di bawah induknya — **belum punya rumah bernama** | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `SourceStatusActive` | DIPETAKAN·DIBUANG-ULANG | — | — | — | — | §14.5 — idem | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `SpecialConditions` | DIPETAKAN | VERSI_KONTRAK | KETENTUAN_KHUSUS | T | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `SpecialConditionsP` | DIPETAKAN | VERSI_KONTRAK | KETENTUAN_KHUSUS | T | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TeritorialScope` | DIPETAKAN | VERSI_KONTRAK | LINGKUP_WILAYAH | T | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Termination` | DIPETAKAN | KONTRAK | TANGGAL_BERAKHIR | D | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TreatyContractName` | DIPETAKAN | VERSI_KONTRAK | NAMA_KONTRAK | T | tidak | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TreatyLeader` | DIPETAKAN | VERSI_KONTRAK | ID_KETUA_TREATY | R | ya | tercocokkan **cocok-tunggal** dari kolom *Asal* §10 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TreatyYear` | DIPETAKAN·TURUNAN-F2 | — | — | — | — | diturunkan `@substring(Commencement,0,4)` (`TreatyInSetTreatyYear`); sisa pertanyaannya dipersempit ke `T-7` — `COCOK-SILANG-CACAH-ATRIBUT.md` §4.3 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `EGNPI[].AmountIDR` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].NetPremiumList` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].RnmLimitList[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].RnmLimitList[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `LimitFacShareSummaryList` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `LimitFacShareSummaryList[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `LimitShareSummaryList` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `LimitShareSummaryList[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `LimitSummaryList` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `LimitSummaryList[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].EgnpiTotalList[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].EgnpiTotalList[].CurrencyID` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].EgnpiTotalList[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].MDPList[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].MDPList[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].PremiumEarnedList[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].PremiumEarnedList[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].ROLPct` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Reinstatement_List[].AdditionalAmount1` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Reinstatement_List[].AdditionalAmount2` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `MDPSummaryList` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `MDPSummaryList[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ProRateDays` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ProRatePercent` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ProRateTotalDays` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].CessionPct` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].RNMShareList` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].RetentionPct` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].GrossPremiumMinList[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].GrossPremiumMinList[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].NetPremiumList` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].RNMSpreadedListDeductRIXOL[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].RNMSpreadedListDeductRIXOL[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].RNMSpreadedListDeductXOL[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].RNMSpreadedListDeductXOL[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].RNMSpreadedListNetRIXOL[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].RNMSpreadedListNetRIXOL[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].RNMSpreadedListNetXOL[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].RNMSpreadedListNetXOL[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].RnmLimitList[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].RnmLimitList[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListDeductRIXOL[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListDeductRIXOL[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListDeductXOL[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListDeductXOL[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListGrossMinXOL[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListGrossMinXOL[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListGrossRIMinXOL[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListGrossRIMinXOL[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListGrossRIXOL[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListGrossRIXOL[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListGrossXOL[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListGrossXOL[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListNetRIXOL[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListNetRIXOL[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListNetXOL[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListNetXOL[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListRIXOL[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListRIXOL[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListXOL[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Share[].SpreadingListXOL[].RNMSpreadedListXOL[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalEgnpiAmount` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalEgnpiAmountNP` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalEgnpiAmountNP[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalEgnpiAmountNP[].AltValue` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalEgnpiAmountNP[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalEgnpiAmountNP[].ID` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalEgnpiAmountNP[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalEgnpiProportion` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalFacShareDeductionNP` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalFacShareDeductionNP[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalFacShareDeductionNP[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalFacShareDeductionNP[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalFacShareGrossNP` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalFacShareGrossNP[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalFacShareGrossNP[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalFacShareGrossNP[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalFacShareNetNP` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalFacShareNetNP[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalFacShareNetNP[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalFacShareNetNP[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalFacShareRnmNP` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalFacShareRnmNP[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalFacShareRnmNP[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalFacShareRnmNP[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalInstallmentAmount` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalInstallmentNP` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalInstallmentNP[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalInstallmentNP[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalInstallmentNP[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalInstallmentPct` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitDeductblNP` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitDeductblNP[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitDeductblNP[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitDeductblNP[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitIOONP` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitIOONP[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitIOONP[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitIOONP[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitMDPMinNP` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitMDPMinNP[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitMDPMinNP[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitMDPMinNP[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitMDPNP` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitMDPNP[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitMDPNP[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitMDPNP[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitPremiEarnNP` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitPremiEarnNP[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitPremiEarnNP[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitPremiEarnNP[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitsAdjPct` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitsDeductible` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitsIOOLimit` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitsMdp` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitsPremiumEarned` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalLimitsROL` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalRetentionAmount` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalRetentionAmountNP` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalRetentionAmountNP[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalRetentionAmountNP[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalRetentionAmountNP[].ID` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalRetentionAmountNP[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareDeductionNP` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareDeductionNP[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareDeductionNP[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareDeductionNP[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareGross` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareGrossMinNP` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareGrossMinNP[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareGrossMinNP[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareGrossMinNP[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareGrossNP` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareGrossNP[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareGrossNP[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareGrossNP[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareNet` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareNetNP` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareNetNP[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareNetNP[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareNetNP[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareRnmLimit` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareRnmNP` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareRnmNP[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareRnmNP[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareRnmNP[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareRnmProp` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareRnmProp[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareRnmProp[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalShareRnmProp[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalSpreadedNetPremi` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalSpreadedNetPremiRI` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalSpreadedNetPremiRI[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalSpreadedNetPremiRI[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalSpreadedNetPremiRI[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalSpreadedNetPremi[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalSpreadedNetPremi[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalSpreadedNetPremi[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalSpreadedRnmProp` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalSpreadedRnmProp[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalSpreadedRnmProp[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalSpreadedRnmProp[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalSpreadedRnmRIProp` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalSpreadedRnmRIProp[]` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalSpreadedRnmRIProp[].Currency` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `TotalSpreadedRnmRIProp[].Value` | DITURUNKAN | — | — | — | — | turunan — tidak disimpan (ADR-0037, §4) | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `BrokeragePct` | DIBUANG | — | — | — | — | dibuang — alasan per jalur di `PETA-TELUSUR-JSON.md` §3 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ChooseStatusAkseptasi` | DIBUANG | — | — | — | — | dibuang — alasan per jalur di `PETA-TELUSUR-JSON.md` §3 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ContractRefNo` | DIBUANG | — | — | — | — | dibuang — alasan per jalur di `PETA-TELUSUR-JSON.md` §3 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `FacultativeShareList[].Limit2` | DIBUANG | — | — | — | — | dibuang — alasan per jalur di `PETA-TELUSUR-JSON.md` §3 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `IsEditData` | DIBUANG | — | — | — | — | dibuang — alasan per jalur di `PETA-TELUSUR-JSON.md` §3 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `OLDID` | DIBUANG | — | — | — | — | dibuang — alasan per jalur di `PETA-TELUSUR-JSON.md` §3 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Position` | DIBUANG | — | — | — | — | dibuang — alasan per jalur di `PETA-TELUSUR-JSON.md` §3 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `RevisionState` | DIBUANG | — | — | — | — | dibuang — alasan per jalur di `PETA-TELUSUR-JSON.md` §3 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].pxCreateDateTime` | DIBUANG | — | — | — | — | dibuang — alasan per jalur di `PETA-TELUSUR-JSON.md` §3 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ShareFacultativeReinsurers[].pxCreateOpName` | DIBUANG | — | — | — | — | dibuang — alasan per jalur di `PETA-TELUSUR-JSON.md` §3 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `StatusAkseptasi` | DIBUANG | — | — | — | — | dibuang — alasan per jalur di `PETA-TELUSUR-JSON.md` §3 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `ViewState` | DIBUANG | — | — | — | — | dibuang — alasan per jalur di `PETA-TELUSUR-JSON.md` §3 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `EDMMaterialType` | DITUNDA | — | — | — | — | ditunda — yang ditunggu dan pemiliknya di `PETA-TELUSUR-JSON.md` §4 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].ProfitCommision` | DITUNDA | — | — | — | — | ditunda — yang ditunggu dan pemiliknya di `PETA-TELUSUR-JSON.md` §4 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].ProfitME` | DITUNDA | — | — | — | — | ditunda — yang ditunggu dan pemiliknya di `PETA-TELUSUR-JSON.md` §4 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |
| `Limits[].Detail[].ProfitYDCF` | DITUNDA | — | — | — | — | ditunda — yang ditunggu dan pemiliknya di `PETA-TELUSUR-JSON.md` §4 | `EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)` |

---

## 4. Jalur yang DIKECUALIKAN — 272, milik modul Adjustment

Pohon `ActualValue`, `ValueDifference`, `OLDDATA`, dan `ValueBeforeProrate`. **Tidak diadili di
sini**, dan itu batas prompt, bukan kelalaian.

> `GRL-14` grilling Adjustment sudah memutuskan nasib `ActualValue`: **dibuang seluruhnya**,
> premi aktual menjadi **nilai versinya sendiri**. Maka 108 jalur `ActualValue.*` di antara yang
> dikecualikan ini **tidak akan melahirkan satu pun kolom**, dan itu keputusan yang sudah
> terkunci — bukan pekerjaan yang menunggu.

| Nasib di §6 | Jumlah |
|---|---:|
| DIPETAKAN | 272 |
| **JUMLAH** | **272** |
