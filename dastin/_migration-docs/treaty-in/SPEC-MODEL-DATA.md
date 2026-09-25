# SPEC-MODEL-DATA — Treaty In

**Tanggal:** 23 September 2026
**Keadaan berkas:** **langkah 1–10 sudah dijalankan.** Bentuk awalnya menuangkan langkah 1–7 ke
disk; langkah 8 (invarian, kini **68**), 9 (peta telusur per jalur — **selesai, semestanya kurang**,
lihat `4-erd-dan-tabel-datar/AUDIT-PENYEBUT-POHON.md`), dan 10 (ERD) **sudah selesai** dan §11.1
mencatat keadaannya masing-masing.

> **Koreksi 24 September 2026 (langkah 4 sesi to-spec).** Kepala berkas ini berbunyi *"langkah 8,
> 9, 10 belum dijalankan"* sampai hari ini, padahal §11.1 di dalam berkas yang sama sudah
> menyatakan ketiganya selesai sejak 23 September. **Kepala berkas dibaca orang yang tidak membaca
> §11**, dan selama dua hari ia memberi tahu mereka bahwa sepertiga pekerjaan belum ada. Sisa
> pekerjaan yang sebenarnya ada di §11.2, bukan di sini.
**Skema sasaran:** `TREATY_MASUK` · akun aplikasi `TREATY_MASUK_APP` (CONTEXT.md §3.5, ADR-0028)

---

## 0. Kadar kepercayaan berkas ini, dibaca sebelum isinya

Berkas ini lahir dari sebuah kegagalan proses: sesi to-spec menghasilkan tujuh langkah kerja tanpa
menuliskan keluarannya, lalu konteks kerjanya hilang. Aturan kerja baru menutup sebabnya
(CONTEXT.md §2.6), dan berkas ini memulihkan akibatnya.

Pemulihan itu **tidak dilakukan dari ingatan**. Perkakas yang menghasilkan adjudikasi langkah 3
— `inv2.py` (penyisir ekspor) dan `klas.py` (pengklasifikasi) — selamat di scratchpad, dan seluruh
angka di bawah **direproduksi ulang dengan menjalankannya kembali atas ekspor XML**. Karena itu:

| Bagian | Dasarnya | Kadar |
|---|---|---|
| §2 daftar entitas dan pembagian KONTRAK/VERSI | keputusan sesi to-spec, diratifikasi pemilik proses | **diputuskan** |
| §3 daftar 187 jalur masukan | dijalankan ulang dari ekspor, 23 Sep 2026 | **terverifikasi ulang** |
| §4 turunan, §5 dibuang, §6 ditunda | dijalankan ulang dari ekspor | **terverifikasi ulang** |
| §7 penggolongan tingkat paket uang | keputusan sesi to-spec | **diputuskan**, 5 di antaranya asumsi |
| §8 bentuk pewarisan cabang | keputusan sesi to-spec, diuji terhadap ekspor | **diputuskan** |
| §9 keadaan siklus hidup | ADR-0055 | **diputuskan** |
| **§10 penamaan dan tipe per atribut** | **belum ada** | **PEKERJAAN BARU, belum pernah ditinjau siapa pun** |

### Tiga hal yang TIDAK dapat dipulihkan, dan tidak saya rekonstruksi

1. **Angka 163.** Sesi to-spec melaporkan "163 atribut model dari 667 jalur (24%)". Yang dapat
   direproduksi hari ini adalah **187 jalur masukan berupa daun**. Selisih 24 adalah langkah
   **peleburan** — beberapa jalur yang bersama-sama membentuk satu paket uang (nilai + mata uang +
   kurs) menjadi **satu** atribut. Keputusan peleburan itu tidak pernah ditulis, dan saya tidak
   merekonstruksinya. **Baseline resmi berkas ini adalah 187**, dan peleburan ke bentuk akhir
   dikerjakan ulang di §10 sebagai pekerjaan yang ditinjau.

2. **Alasan per jalur untuk sebagian yang dibuang.** Sebelas jalur dibuang dapat direproduksi
   beserta polanya; alasan satu-kalimat per jalur sebagian saya turunkan ulang di §5 dan
   **ditandai** mana yang diturunkan ulang.

3. **Nama Indonesia per atribut.** Tidak pernah ada. Tidak ada satu pun dari 187 yang pernah
   dinamai, ditipekan, atau diperiksa terhadap §16. Itu seluruh isi §10.

### Satu kesalahan adjudikasi yang ditemukan saat menjalankan ulang

`ProportionType` terklasifikasi sebagai **turunan** oleh pengklasifikasi, karena pola pencocokannya
mengandung kata `Proportion`. Itu **salah**: `ProportionType` adalah **masukan**, dan ia bagian
lapisan beku `KONTRAK` (§2). Hal yang sama terjadi pada `TotalEgnpiProportion`, yang memang
turunan tetapi karena alasan lain.

Ini persis jenis temuan yang dicari: sebuah jalur yang tidak pernah benar-benar diadili, hanya
tercocokkan. Langkah 9 harus menyisir seluruh 667 jalur dengan pemeriksaan semacam ini, bukan
menyalin keluaran skrip.

---

## 1. Angka pokok, direproduksi 23 September 2026

```
JALUR BERSIH YANG DIADILI: 667

  CERMIN              272    40.8%      ValueDifference 162 · ActualValue 108
                                        OLDDATA 1 · ValueBeforeProrate 1
  MASUKAN             213    31.9%      187 di antaranya DAUN
  TURUNAN-agregat     113    16.9%       65 di antaranya DAUN
  TURUNAN-hitungan     54     8.1%       54 di antaranya DAUN
  DIBUANG              11     1.6%
  DITUNDA               4     0.6%
```

**Cermin bukan kategori final.** 272 jalur (40,8%) adalah bayangan pohon utama —
`ValueDifference`, `ActualValue`, `OLDDATA`, `ValueBeforeProrate` — dan **tidak menyumbang satu pun
atribut baru**. 104 bentuknya identik dengan pohon utama; 73 sisanya ternyata **lubang pada
penyisiran pohon utama saya**, dan itu dicatat sebagai temuan, bukan sebagai catatan metode.

Sisanya 395 jalur non-cermin. Angka yang dilaporkan sesi to-spec (DIPETAKAN 207 · DITURUNKAN 172 ·
DIBUANG 12 · DITUNDA 4) adalah angka **setelah koreksi manual** pemilik proses; keluaran skrip di
atas adalah baseline **sebelum** koreksi. Keduanya berjumlah 395. Selisihnya: `BrokeragePercentP`
dipulihkan sebagai atribut, `BrokeragePct` dibuang, `ContractRefNo` diberi alasan, `TreatyYear`
dipindah ke turunan, dan kategori kelima yang sempat saya buat dilebur kembali ke DIPETAKAN.

---

## 1a. Angka sesudah §10.23c — apa yang bergerak dan apa yang TIDAK

**Ditulis 24 September 2026, langkah 4 sesi to-spec.** §11.2 menuntut *"memperbarui angka §1, §3,
§4, §7, karena keputusan §10.23c mengubah penyebutnya"*. Diperiksa satu per satu; **hanya satu yang
bergerak**, dan menyatakan yang lain **tidak** bergerak sama pentingnya.

| Angka | Sebelum | Sesudah | Sebab |
|---|---:|---:|---|
| **jumlah entitas penyerahan pertama** | 27 | **28** | `PEMULIHAN_LIMIT` diterima (§10.23c butir 1). Tiga usul lain gugur |
| §1 — 667 jalur bersih diadili | 667 | **667** | tidak bergerak — lihat kotak di bawah |
| §3 — 187 jalur masukan daun | 187 | **187** | tidak bergerak — sebab yang sama |
| §4 — 119 jalur turunan daun | 119 | **119** | tidak bergerak |
| §7 — 43 paket uang disimpan | 43 | **43** | `PEMULIHAN_LIMIT` membawa dua **persentase**, bukan paket uang |

> ### Kenapa ketiga angka jalur TIDAK bergerak, dan kenapa itu kabar buruk
>
> `PEMULIHAN_LIMIT` bersumber pada `SetReinstatementPct`, yang menulis `ReinstatementAmount1`,
> `ReinstatementAmount2`, dan `ReinstatementNote` pada `Data-TreatyInLimits`. **Ketiganya ada di
> daftar titik buta** `4-erd-dan-tabel-datar/datar-titik-buta-pohon.csv`, bertanda `TIDAK` — tidak
> pernah ada di pohon 985 simpul (`L-8`).
>
> Maka ketiganya **tidak pernah termasuk dalam 667**, dan menambahkan entitasnya **tidak mengubah
> penyebutnya**. Angka 667, 187, dan 119 bukan salah — **ia dihitung atas semesta yang kurang**,
> dan itu persis yang `AUDIT-PENYEBUT-POHON.md` catat.
>
> **Maka angka-angka itu dibiarkan apa adanya dan diberi tanda**, bukan ditambal dengan taksiran.
> Menaikkan 667 menjadi angka karangan akan membuat semesta yang bolong **terlihat utuh**, dan itu
> lebih berbahaya daripada bolong yang bertanda.

### Satu lubang yang ditemukan saat memeriksa angka ini, dan TIDAK ditambal

`5-tiket/KEPUTUSAN-PEMBAGIAN-TIKET.md` §0.2 butir 4 menetapkan: **soal daftar entitas,
`4-erd-dan-tabel-datar/STRUKTUR-DATA.md` yang mengikat.** Berkas itu **belum memuat
`PEMULIHAN_LIMIT`.**

Artinya wewenang tertinggi atas daftar entitas sekarang **menyebut 27**, sementara §10.23c sudah
memutuskan **28** — dan urutan wewenang mengatakan yang mengikat adalah yang kurang.

| | |
|---|---|
| **Siapa menutup** | sesi yang memperbarui `4-erd-dan-tabel-datar/` — alur B `PROMPT-TO-TICKET-DAN-STRUKTUR-DATA.md` |
| **Yang menagih** | butir 4 urutan wewenang itu sendiri. Siapa pun yang menghitung ruang lingkup **wajib** membuka `STRUKTUR-DATA.md`, dan akan menemukan 27 |
| **Kenapa tidak ditambal di sini** | berkas induk tidak disunting tanpa persetujuan, dan menyuntingnya diam-diam berarti mengubah **sumber yang mengikat** lewat sesi yang bukan pemiliknya |

---

## 2. Daftar entitas, dan pembagian `KONTRAK` / `VERSI_KONTRAK`

### 2.1 Pembagiannya, dan kenapa

**`KONTRAK` hanya memuat lapisan identitas yang dibekukan grilling.** Semua yang lain — seluruh 60
atribut kepala dan **seluruh entitas anak** — tinggal di `VERSI_KONTRAK`.

Alasannya tercatat sebagai koreksi pemilik proses: membekukan sesuatu yang dapat bertentangan
dengan tetangganya berarti membekukan sebuah kontradiksi. Lapisan beku hanya boleh memuat hal yang
**mendefinisikan kontrak itu sebagai kontrak**; sisanya adalah hal yang dapat berubah lewat
addendum, dan karena itu milik versi.

**Satu konsekuensi jatuh gratis dari bentuk ini**, dan ia memenuhi arahan pemilik proses tanpa
mekanisme apa pun: karena seluruh nilai tinggal di versi, **`OLDDATA` menjadi SELECT atas versi
sebelumnya**. Tidak ada tabel bayangan, tidak ada penyalinan. Lihat `SEAM-ADJUSTMENT.md`.

### 2.2 `KONTRAK` — lapisan beku

| Atribut (sistem lama) | Sifat |
|---|---|
| *(pengenal baru)* | pengenal buatan sistem, tanpa makna — ADR-0040 |
| `ID` | identitas warisan, disimpan berdampingan, bukan pengganti — ADR-0042 |
| *(relasi baru)* `DISALIN_DARI` | pecahan pertama dari `OLDID` — ADR-0040 §1 |
| `CedingID` + `Ceding` | cedant |
| `LeadingReinsSourceID` + `LeadingReinsSource` | source of business |
| `ProportionType` | sifat proporsi |
| `Commencement` | tanggal mulai |
| `Termination` | tanggal berakhir |

Empat terakhir ditambah cedant membentuk **kunci alami** kontrak: cedant + SoB + periode +
`ProportionType`. Ia **memperingatkan, tidak melarang** — ADR-0040 §2.

`TreatyYear` **tidak ada di sini.** Ia diputuskan **turunan** dari `Commencement` dan tidak
disimpan. Itu asumsi, dan Uji W mengukurnya: bila ada kontrak yang `TREATYYEAR`-nya berbeda dari
tahun `COMMENCEMENT`, `TreatyYear` naik menjadi masukan dan aturan mana yang dipercaya saat
keduanya berselisih harus ditetapkan bisnis.

`OLDID` **dibuang sebagai kolom** karena ia menampung dua arti — "salinan dari" dan "addendum atas"
— dan diganti dua relasi bernama. Pecahan kedua adalah relasi versi itu sendiri.

### 2.3 Daftar entitas lengkap

**Milik bersama kedua cabang, menggantung pada `VERSI_KONTRAK`:**

| Entitas | Sumber di JSON |
|---|---|
| `MATA_UANG_KONTRAK` | `CurrencyList[]` |
| `RETENSI_CEDANT` | `Retention[]` |
| `EGNPI` | `EGNPI[]` |
| `PORTOFOLIO` | `Portfolio[]` |
| `PERIODE_PELAPORAN` | `ReportingPeriodList[]` |
| `PERIODE_AKUMULASI` | `AccumulationList[]` |
| `TERMIN` | `Installment[]` |
| `SKALA_KOASURANSI` | `CoInScale` — **daftar, bukan skalar** |
| `CATATAN_PERSETUJUAN` | `CommentList[]` — menggantung pada **versi** |
| `DOKUMEN_KONTRAK` | lampiran (`M_ATTACHMENTTREATY_2`) — dirujuk, tidak dimiliki (ADR-0027) |
| `JEJAK_PERUBAHAN` | tidak ada di sistem lama — ADR-0045 |

**Bagian, potongan, penyebaran — satu bentuk untuk kedua cabang:**

| Entitas | Sumber di JSON |
|---|---|
| `BAGIAN` | `Share[]` — **hanya cabang non-proporsional**, lihat §12.3 |
| `POTONGAN` | daftar potongan pada bagian |
| `PENYEBARAN` | `SpreadingList` / `SpreadingListXOL` |
| `RINCIAN_PENYEBARAN` | rincian per pihak |
| `NILAI_PENYEBARAN` | nilai per pihak per mata uang |

**Cabang — satu-satunya pasangan yang benar-benar terpisah:**

| Entitas | Sumber di JSON | Cabang |
|---|---|---|
| `LAYER` | `Limits[]` | non-proporsional |
| `DETAIL_PROPORSIONAL` | `Limits[].Detail[]` | proporsional |

**Ditambahkan 24 September 2026 — tujuh entitas yang hilang dari §2.3 (K-5, menutup `L-6`):**

> `5-tiket/LUBANG-SPESIFIKASI.md` **L-6** mencatat bahwa siapa pun yang menghitung ruang lingkup
> dari §2.3 **kekurangan tujuh entitas**, dan menambahkan butir 4 urutan wewenang sebagai penjaga:
> soal daftar entitas, `4-erd-dan-tabel-datar/STRUKTUR-DATA.md` yang mengikat.
>
> **Butir 4 tetap berlaku dan tidak dicabut.** Yang berubah: §2.3 tidak lagi salah. Penjaga
> prosedural dipertahankan **di samping** perbaikannya, bukan sebagai penggantinya — karena §2.3
> dapat tertinggal lagi, dan butir 4 tidak.

| Entitas | Sumber di JSON | Kenapa hilang semula |
|---|---|---|
| `BATAS_PER_BAHAYA` | `Earthquake` · `FloodJab` · `FloodNation` · `RSMDLimit` + mata uangnya | diperkenalkan §10.0b, ditegakkan INV-14 — **sesudah** §2.3 ditulis |
| `MATA_UANG` | tabel acuan (ADR-0038) | keenam tabel acuan tidak pernah masuk §2.3 sama sekali |
| `JENIS_POTONGAN` | tabel acuan | idem |
| `JENIS_REASURANSI` | tabel acuan | idem |
| `BAHAYA` | tabel acuan | idem |
| `KELOMPOK_TREATY` | tabel acuan | idem |
| `KELAS_BISNIS` | tabel acuan | idem |

**Dan satu entitas kedelapan, lahir sesudah L-6 ditulis:**

| Entitas | Sumber | Kedudukan |
|---|---|---|
| `PEMULIHAN_LIMIT` | `SetReinstatementPct` — pemulihan limit berulang | **§10.23c butir 1, DITERIMA.** Ia yang mengubah 27 menjadi **28** |

> **Maka jumlah entitas penyerahan pertama adalah 28**, dan §2.3 kini memuat seluruhnya. Kemampuan
> *"syarat berbeda tiap pemulihan"* bergolongan **BARU** (P-59) — sistem lama menulis kedua
> persentasenya tetap `"100"` di dalam kode.

**Ditambahkan 24 September 2026 — entitas kesembilan, dari modul Adjustment (diff `D-2`):**

| Entitas | Sumber di JSON | Kedudukan |
|---|---|---|
| **`DOKUMEN_ADDENDUM`** | **tidak ada** — sistem lama tidak pernah merekamnya | **BARU**, `GRL-19`; berdiri sendiri, memayungi banyak versi lintas kontrak |

> **`GRL-01` dikoreksi, bukan dibongkar.** Addendum **tetap** `VERSI_KONTRAK`; dokumen adalah
> **pembungkus** beberapa addendum, satu tingkat di atasnya. Persetujuan **tetap per versi per
> kontrak**.

**Di luar batas, dicatat agar tidak terlupa:**

> **Koreksi 24 September 2026 — `G2` dan `G3` di tabel bawah ini berarti GELOMBANG, bukan gerbang.**
> Huruf `G` sudah dipakai untuk tiga hal berbeda di proyek ini: gelombang penyerahan (di sini),
> gerbang sambungan Adjustment (`4-erd-dan-tabel-datar/KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md`), dan
> gerbang pembagian tiket (`5-tiket/KEPUTUSAN-PEMBAGIAN-TIKET.md`). Sejak folder `5-tiket/`
> gelombang ditulis **`GEL-2`** dan **`GEL-3`**. Baca dua baris di bawah sebagai **GEL-2** dan
> **GEL-3**.

| Entitas | Kedudukan |
|---|---|
| `RETRO_KELUAR` | **G2** ⟦= **GEL-2**⟧ — `ShareFacultativeReinsurers[]`, arahnya keluar; ADR-0020 **tidak berlaku** karena pembedanya **arah** (fac IN vs fac OUT), bukan pokoknya |
| `PENCAPAIAN` | **G3** ⟦= **GEL-3**⟧ — `ACHIEVEMENT`; dimodelkan sekarang, tidak dibangun sekarang |

---

## 3. Daftar 187 jalur masukan, per entitas

Kolom "cabang" adalah hasil penyisiran penulis aturan: `PROP`, `NONPROP`, `KEDUA`, `UMUM`, atau
`TAK-TERTULIS` (tidak pernah ditulis oleh aturan mana pun — dibaca saja, biasanya diisi dari layar).

> **`TAK-TERTULIS` bukan alasan membuang.** Ia berarti tidak ada *aturan Pega* yang menulisnya,
> bukan tidak ada yang mengisinya. Field yang diisi langsung dari layar tampak persis begitu.
> Satu-satunya jalur yang dibuang karena pola ini adalah `ContractRefNo`, dan itu karena ia
> `TAK-TERTULIS` **sekaligus** tidak punya kolom di DDL mana pun **sekaligus** `pyReadOnly=true`
> di seluruh layar — tiga hal, bukan satu.

### 3.1 `KONTRAK` — 8 jalur

```
Ceding                  UMUM          CedingID                UMUM
Commencement            UMUM          Termination             UMUM
LeadingReinsSource      UMUM          LeadingReinsSourceID    UMUM
ID                      UMUM          ProportionType          (salah terklasifikasi — lihat §0)
```

### 3.2 `VERSI_KONTRAK` — 60 jalur kepala

```
AccountingMode           TAK-TERTULIS   AccountingModeNonProp    TAK-TERTULIS
AccumulationPeriod       TAK-TERTULIS   Bordeaux                 UMUM
BordereauxNote           UMUM           BrokeragePercent         TAK-TERTULIS
BrokeragePercentP        TAK-TERTULIS   CedingStatusActive       TAK-TERTULIS
ClassofBusiness          TAK-TERTULIS
Comment                  UMUM           Currency                 UMUM
CurrencyEarthquake       TAK-TERTULIS   CurrencyEgnpiAmount      UMUM
CurrencyFloodJab         TAK-TERTULIS   CurrencyFloodNat         TAK-TERTULIS
CurrencyInstallmentAmount TAK-TERTULIS  CurrencyRSMD             TAK-TERTULIS
EDMEffective             TAK-TERTULIS   EDMState                 TAK-TERTULIS
Earthquake               TAK-TERTULIS   Exclusions               UMUM
ExclusionsP              UMUM           FacShare                 NONPROP
FacShareBrokerage        NONPROP        FacultativeShare         UMUM
FacultativeShareBrokerage UMUM          FloodJab                 TAK-TERTULIS
FloodNation              TAK-TERTULIS   Information              UMUM
InstallmentNo            UMUM           IsMultipleRetro          TAK-TERTULIS
IsProRate                TAK-TERTULIS   LeadingReinsID           TAK-TERTULIS
LeadingReinsName         TAK-TERTULIS   MaxCoGroup               TAK-TERTULIS
MaxCoNonGroup            TAK-TERTULIS   NusareSharePct           TAK-TERTULIS
OptionLimit              TAK-TERTULIS
RNMShare                 UMUM           RNMShareAcrossTheBoard   TAK-TERTULIS
RNMShareP                UMUM           RSMDLimit                TAK-TERTULIS
ReminderDays             TAK-TERTULIS   ReportingConfirmation    TAK-TERTULIS
ReportingEnd             TAK-TERTULIS   ReportingInterval        TAK-TERTULIS
ReportingPeriod          UMUM           ReportingSettlement      TAK-TERTULIS
ReportingStart           TAK-TERTULIS   ReportingSubmission      TAK-TERTULIS
RetroList                TAK-TERTULIS   RnmShareDeducted         UMUM
SourceStatusActive       TAK-TERTULIS   SpecialConditions        UMUM
SpecialConditionsP       UMUM           TeritorialScope          UMUM
TreatyContractName       UMUM           TreatyLeader             TAK-TERTULIS
```

**Empat pasangan bersufiks `P` adalah pasangan cabang, bukan salah eja.** `BrokeragePercent` /
`BrokeragePercentP`, `Exclusions` / `ExclusionsP`, `SpecialConditions` / `SpecialConditionsP`, dan
`FacShare` / `FacultativeShare`. Sufiks `P` = **Proportional**: `BrokeragePercentP` tidak pernah
ditulis satu kali pun, tetapi dirujuk di delapan berkas dan **seluruhnya di cabang proporsional**
(`…OfferProportional_Act`, `TreatyInPropshare`), sementara `BrokeragePercent` muncul di cabang
non-proporsional. Ini ditemukan setelah saya keliru membuangnya sebagai duplikat ejaan.

`BrokeragePct` **dibuang** — lihat §5.

#### 3.2a Dua properti yang menulis ke tabel datar dan TIDAK PERNAH DIADILI — 24 September 2026

Disapu saat mencari empat atribut yang hilang dari §10.2 (lihat blok selisih di §10.2). Keduanya
**tidak menjelaskan selisih itu** — mereka persoalan tersendiri, dan lebih berat.

| Properti | Muncul di ekspor | Penulisnya | Diadili di §10 / §4 / §5 / §6? |
|---|---|---|---|
| **`TreatyIn.NusareSharePct`** | **dua berkas saja** — `RDBList/SaveTreatyIn.xml` dan `SaveTreatyInEDM.xml` | **nol**, atas kelima bentuk | **tidak, di mana pun.** Satu-satunya sebutannya di berkas ini adalah baris daftar §3.2 |
| **`TreatyIn.BrokeragePct`** | dua berkas yang sama | **nol** | **tidak** — dan ia **tidak ada di daftar §3.2 sama sekali**, sehingga berada **di luar semesta 667 jalur yang pernah diadili** |

Keduanya diteruskan sebagai parameter ke `POOLDATA.PEGA_TREATY_IN` dan menulis ke
**`TREATY_IN.NUSARESHAREPCT`** dan **`TREATY_IN.BROKERAGEPCT`** pada **setiap penyimpanan**.

> **Dengan nol penulis, keduanya selalu kosong.** Maka **dua dari dua puluh kolom bisnis
> `TREATY_IN`** — dan dua dari dua puluh empat di `TREATY_IN_EDM` — **tidak pernah berisi apa pun**,
> sementara nilai yang sebenarnya hidup di `RNMShare` dan `BrokeragePercent`. Siapa pun yang membaca
> persentase bagian NuRe atau brokerage dari tabel datar memperoleh kosong.
>
> **Batas klaim "nol penulis":** `L-10` — `Declare Expression` dan `Declare Trigger` tidak terekspor
> sama sekali. Yang terbukti: **tidak ada penulis berbentuk yang dapat diperiksa**. Ujinya data, dan
> ia sudah masuk lampiran pengukuran `DAFTAR-ESKALASI-MANAJEMEN.md`.

**Akibat pada model: nol.** Keduanya tidak menjadi atribut — `PERSEN_BAGIAN_NURE` dan
`PERSEN_BROKERAGE` sudah bersumber pada `RNMShare` dan `BrokeragePercent`. Yang terdampak adalah
**pembaca tabel datar di luar sistem ini** (ADR-0051).

### 3.3 `LAYER` (non-proporsional) — 14 jalur

```
Limits[].ID                 Limits[].Layer              Limits[].LayerPart
Limits[].LayerType          Limits[].LayerPartType      Limits[].Cover
Limits[].Currency           Limits[].Limit              Limits[].Deductible
Limits[].AdjRate            Limits[].MDPPct             Limits[].ReinstatementPct
Limits[].ReinstatementValue Limits[].TreatyType
```

Ditambah anak `Limits[].TreatyGroupList[]` → kelompok treaty dan kelas bisnis per layer:
`TreatyGroup`, `TreatyGroupID`, `ClassOfBusinessList[].ClassOfBusiness`,
`ClassOfBusinessList[].ClassOfBusinessID`.

**`Limits[].ReinstatementPct` dan `Limits[].ReinstatementValue` adalah dua besaran terpisah**, dan
hubungannya **perkalian**, bukan persamaan: porsi limit yang dipulihkan, dan tarif premi yang
ditagih untuk pemulihan itu. "Pulihkan penuh, gratis" adalah ketentuan yang sah. Keduanya
**masukan**, dan di sistem lama rasio limit menimpa tarif premi — cacat yang tersembunyi karena
keduanya disemai 100.

### 3.4 `DETAIL_PROPORSIONAL` — 7 jalur langsung + 7 daftar anak

```
Limits[].Detail[].QSPct            Limits[].Detail[].Surplus
Limits[].Detail[].RIOGR            Limits[].Detail[].RIONR
Limits[].Detail[].TreatyGroup      Limits[].Detail[].TreatyGroupID
Limits[].Detail[].TreatyType
```

Daftar anaknya, masing-masing berpasangan `Currency` + `Value`:

```
Limits[].Detail[].RetentionList[]    retensi cedant, sebagai UANG (Surplus)
Limits[].Detail[].IOOLimitList[]     -> KAPASITAS_SURPLUS
Limits[].Detail[].EPIList[]          estimated premium income
Limits[].Detail[].ClaimCoopList[]    claim cooperation
Limits[].Detail[].CashLossList[]     cash loss
Limits[].Detail[].PLAList[]          PLA = PRELIMINARY LOSS ADVICE (terpecahkan 24 Sep 2026)
Limits[].Detail[].COBList[]          kelas bisnis (ClassOfBusiness + ID)
```

**`TreatyType` hidup di `Detail[]`, bukan di `Limits[]`** — dan pencari baris Quota Share adalah
**gelung dua tingkat** (`TreatyIn.Limits` di luar, `.Detail` di dalam) yang **menyeberangi** baris
`Limits[]`. Karena itu ketergantungan Quota Share ↔ Surplus berlingkup **versi kontrak**, bukan
berlingkup baris `Limits[]`. Ini invarian, dan ia masuk langkah 8:

> Kapasitas Surplus dihitung dari retensi bagian Quota Share pada **versi kontrak yang sama**.
> Ketiadaannya adalah kegagalan yang dilaporkan, bukan nol (ADR-0035).

`KAPASITAS_SURPLUS` = retensi × jumlah lines. Ia **kelipatan**, bukan partisi — dan karena itu
**dikecualikan dengan nama** dari invarian rekonsiliasi (§ langkah 8).

### 3.5 Entitas anak bersama

```
MATA_UANG_KONTRAK   CurrencyList[].CurrencyID  .Currency  .Conversion  .PeriodStart  .PeriodEnd
RETENSI_CEDANT      Retention[].ID  .Amount  .Currency  .CurrencyID  .TreatyGroup  .TreatyGroupID
EGNPI               EGNPI[].ID  .Amount  .Currency  .CurrencyID  .AsDate  .TreatyGroup  .TreatyGroupID
PORTOFOLIO          Portfolio[].Description  .Type  .TypePortfolio
PERIODE_PELAPORAN   ReportingPeriodList[].Period  .InitialDate  .SubmissionDue
                    .ConfirmationDue  .SettlementDue
PERIODE_AKUMULASI   AccumulationList[].Period  .ReportDate  .SubDays  .SubDueDate
TERMIN              Installment[].ID  .Currency  .CurrencyID  .Installment  .InstallmentPct
                    .Amount  .DueDate  .PaymentDate  .WPC
SKALA_KOASURANSI    CoInScale[].PctLimit  .CoInShare
POTONGAN            DeductionList[].Deduction  .DeductionPct  .Currency  .CurrencyID  .Comment
PENYEBARAN          SpreadingList[].ReinsTypeID  .ReinsTypeName  .ParentReinsTypeID  .Pct
RINCIAN_PENYEBARAN  SpreadingList[].BreakDownSprdList[]  (dan varian XOL)
NILAI_PENYEBARAN    .Currency  .Value
CATATAN_PERSETUJUAN CommentList[].Date  .OperatorName  .IsApproved  .Suggest
```

**`PERIODE_PELAPORAN` menyimpan tanggal jatuh tempo, TIDAK menyimpan bendera terlambat.**
Keterlambatan adalah turunan yang dihitung saat dibaca — keputusan pemilik proses, OPSI 2.

### 3.6 `BAGIAN` dan penyebaran

```
BAGIAN            -> relasi ke LAYER                             (bukan atribut, lihat §12.2)
                  Share[].Cover  .SharePct  .ClassofBusinessList  .TreatyGroupList
                  Share[].GrossPremiumList[].Currency  .Value        (NONPROP, masukan)
PENYEBARAN XOL    Share[].SpreadingListXOL[].GrossPremiumList[].Currency  .Value
```

`Share[].Layer`, `.LayerPart`, `.LayerType`, `.LayerPartType` **bukan atribut `BAGIAN`** — lihat
§12.2.

### 3.7 `RETRO_KELUAR` — G2, di luar batas, dicatat saja

```
ShareFacultativeReinsurers[].ID  .ReinsID  .ReinsName  .BrokerName  .Layer  .SharePct
ShareFacultativeReinsurers[].FacultativeLimits[].TreatyType  .TreatyTypeID
ShareFacultativeReinsurers[].FacultativeLimits[].Detail[].CessionList  .IOOLimitList
   .QSPct  .RNMShare  .ShareNote  .TreatyGroup  .TreatyGroupID  .TreatyType
FacultativeShareList[]…  (cermin sisi non-proporsional, 16 jalur)
ShareReins[].ID
```

`CessionList` / `CessionPct` → `NILAI_PENYERAHAN` / `PERSEN_PENYERAHAN`. **Kata "sesi" dilarang**
di nama mana pun — CONTEXT.md §3.3b.

---

## 4. Turunan — 119 jalur daun, dari apa masing-masing diturunkan

### 4.1 Turunan hitungan — 54 daun

| Turunan | Diturunkan dari |
|---|---|
| `Limits[].ROLPct` | premi layer ÷ limit layer |
| `Limits[].MDPList[]` | `Limits[].MDPPct` × EGNPI — **deposit** |
| `Limits[].PremiumEarnedList[]` | premi × porsi periode berjalan |
| `Limits[].EgnpiTotalList[]` | jumlah `EGNPI[]` per mata uang |
| `Limits[].Reinstatement_List[].AdditionalAmount1/2` | limit × `ReinstatementPct` × tarif |
| `ProRateDays`, `ProRateTotalDays`, `ProRatePercent` | **rumusnya terbaca** (`DataTransform/TreatyCalculateProratePct.xml`, 24 Sep 2026) — `ProRateDays = @DateTimeDifference(EDMEffective, Termination,'D')`; `ProRateTotalDays = @DateTimeDifference(Commencement, Termination,'D')`; `ProRatePercent = @divide(ProRateDays, ProRateTotalDays) * 100`. **`EDMEffective` disemai `= Commencement`**, sehingga hasilnya 100 selama tanggalnya tidak digeser — **Uji AJ** |
| `Share[].GrossPremiumMinList[]` | premi bruto minimum |
| `Share[].NetPremiumList` | bruto − potongan |
| `Share[].RnmLimitList[]` | limit × bagian NuRe |
| `Share[].RNMSpreadedList{Net,Deduct}{,RI}XOL[]` | penyebaran bagian NuRe, 4 varian |
| `Share[].SpreadingListXOL[].RNMSpreadedList…[]` | 11 varian penyebaran manual XOL |
| `EGNPI[].AmountIDR` | `Amount` × kurs |
| `FacultativeShareList[].NetPremiumList`, `.RnmLimitList[]` | sisi fakultatif, bentuk sama |
| `…Detail[].RetentionPct`, `.CessionPct`, `.RNMShareList` | partisi retensi/penyerahan |

### 4.2 Turunan agregat — 65 daun, tidak satu pun disimpan

Seluruhnya berawalan `Total…` atau berakhiran `SummaryList`. Contoh keluarga:
`TotalLimitsDeductible`, `TotalLimitsIOOLimit`, `TotalLimitsMdp`, `TotalLimitsPremiumEarned`,
`TotalLimitsROL`, `TotalLimitsAdjPct`, `TotalShareGross`, `TotalShareNet`, `TotalShareRnmLimit`,
`TotalRetentionAmount`, `TotalEgnpiAmount`, `TotalInstallmentAmount`, `TotalInstallmentPct`,
ditambah varian ber-`NP` yang berulang per mata uang, dan empat `…SummaryList[]`.

**Tidak ada agregat yang disimpan.** Ia dihitung saat dibaca — ADR-0037.

### 4.3 Empat rumus partisi yang dipastikan dari sumbernya

```
Quota Share      : @divide(QSPCT,100,9)*.Value      dan  @divide((100-QSPCT),100,9)*.Value
Retensi + sesi   : .Value*@divide(RetentionPct,100,4) dan .Value*@divide(CessionPct,100,4)
Sisa bruto       : Value = Share(idx).GrossPremiumList(<CURRENT>).Value - .Value
Penyebaran       : SpreadingList.Value = (.Pct/100) * primary.RNMShareList(1).Value
```

Bukan partisi, dan **dikecualikan dengan nama**: `KAPASITAS_SURPLUS = .Value * Primary.Surplus`
(kelipatan), dan reinstatement (boleh berulang).

---

## 5. Dibuang — 12 jalur, masing-masing dengan alasannya

| Jalur | Alasan | Asal alasan |
|---|---|---|
| `ViewState` | turunan yang disimpan dan disetel oleh klik; dilebur ke keadaan siklus hidup | ADR-0046 |
| `IsEditData` | idem | ADR-0046 |
| `RevisionState` | idem; jalur pintas revisi dihapus | ADR-0046, 0052 |
| `Position` | idem; nilai kosongnya menandai dua keadaan berlawanan | ADR-0055 |
| `StatusAkseptasi` | idem | ADR-0046 |
| `ChooseStatusAkseptasi` | bendera pemilihan di layar, bukan fakta kontrak | ADR-0046 |
| `ContractRefNo` | **tidak punya jalur pengisian sama sekali**: tak pernah ditulis, tidak ada di DDL/prosedur mana pun, `pyReadOnly=true` di seluruh layar | diputuskan sesi to-spec |
| `OLDID` | menampung dua arti; diganti dua relasi bernama | ADR-0040 §1 |
| `FacultativeShareList[].Limit2` | pasangan kolom kembar untuk mata uang kedua | pola terlarang |
| `ShareFacultativeReinsurers[].pxCreateDateTime` | fakta mesin Pega; digantikan jejak perubahan sendiri | ADR-0045 — *diturunkan ulang* |
| `ShareFacultativeReinsurers[].pxCreateOpName` | idem | ADR-0045 — *diturunkan ulang* |
| `BrokeragePct` | bukan properti kontrak — ia **nama parameter panggilan Oracle** di dua RDB List | diputuskan sesi to-spec |

---

## 6. Ditunda — 4 jalur, apa yang ditunggu dan siapa pemiliknya

| Jalur | Menunggu | Pemilik |
|---|---|---|
| ~~`EDMMaterialType`~~ | **TIDAK LAGI DITUNDA — 24 September 2026.** Yang ditunggu adalah *"konfirmasi bahwa tidak ada pemakaian lain"*; jawabannya **ada pemakaian lain, dan itu justru yang memutuskan golongannya.** Ia **MASUKAN**, dan ia kini atribut §10.2 — lihat **§10.2a**. ADR-0049 **tidak berubah** | — |
| `Limits[].Detail[].ProfitCommision` | *profit commission* tidak dapat dihitung sebelum periode selesai; bentuk penyimpanannya menunggu keputusan apakah ia terminologi potongan atau peristiwa tersendiri | bisnis (bagian teknik) |
| `Limits[].Detail[].ProfitME` | idem, bersama di atas | bisnis |
| `Limits[].Detail[].ProfitYDCF` | idem, bersama di atas | bisnis |

`CurrencyRelation` **tidak lagi ditunda**: ADR-0053 menetapkannya **pasif** — daftar mata uang yang
memutuskan, bukan bendera relasi.

---

## 7. Paket uang — 43 yang disimpan, dan tingkat pencatatannya

Dari 65 keluarga uang di pohon utama, 22 adalah agregat `Total*` yang tidak disimpan. **43 menjadi
paket uang tersimpan.**

Setiap paket membawa: **nilai, mata uang, kurs, tanggal/sumber kurs, nilai IDR, TINGKAT
PENCATATAN**, dan — bila ia besaran tingkat bagian — **bagian yang dipakai**. ADR-0007 + 0029 +
0039.

| Tingkat | Jumlah |
|---|---|
| 100% treaty | 16 |
| bagian NuRe | 22 |
| **belum dapat ditentukan** | **5** |

### Lima yang belum dapat ditentukan, dipisah menurut akibat kesalahannya

| Besaran | Akibat bila tingkatnya salah | Kedudukan |
|---|---|---|
| `EPIList` (+ `EGNPI`) | **penyajian ulang historis** — rumus pencapaian mengandaikan 100% | **dipisah sendiri → eskalasi**, tidak dibundel |
| `ClaimCoopList` | perbaikan model | dijawab satu slip kontrak |
| `PLAList` | perbaikan model | dijawab satu slip kontrak |
| `CashLossList` | **salah memicu cash call** — kepercayaan diri paling rendah dari kelimanya; kata-kata pasar biasanya menempatkannya di tingkat bagian reasuradur | dijawab satu slip kontrak |

Ketiga yang terakhir adalah **pertanyaan klausul**, bukan pertanyaan data: dijawab dengan
**meminta satu slip treaty non-proporsional yang memuat klausul claim cooperation, cash loss, dan
PLA** — bukan dengan kueri.

### `PLA` TERPECAHKAN — 24 September 2026

Nama sementara yang sengaja tidak dapat dikirim ke produksi, **`TBD_PLA_ARTI_BELUM_DIKETAHUI`**,
**tidak lagi diperlukan.** Kepanjangannya terbaca dari ekspor, dan ia tidak pernah tersembunyi —
ia hanya ada di kelas yang penyisiran pohon tidak pernah lihat (L-8).

`Activity/TreatyInMappingDataconvertProp.xml` mengisi `PLAList` dengan halaman langkah
**`.PreliminaryLossAdvice`**, sebuah wadah pada kelas `ASM-FW-GISFW-Data-TreatyInBusinessLimit`:

```
langkah page   : .PreliminaryLossAdvice        (kelas Data-TreatyInBusinessLimit)
TreatyIn.Limits(<LAST>).Detail(<LAST>).PLAList(<APPEND>).Value    = .Amount
TreatyIn.Limits(<LAST>).Detail(<LAST>).PLAList(<LAST>).Currency   = .Currency.Currency
```

> **`PLA` = *Preliminary Loss Advice*.** Satu besaran, satu mata uang — ambang yang mewajibkan
> cedant memberi pemberitahuan awal kerugian.

Kelas yang sama juga memberi kepanjangan tetangganya, dan keempatnya cocok dengan daftar `Detail[]`:
`ClaimCooperation` → `ClaimCoopList`, `CashLossLimit` → `CashLossList`, `AmountEPI` → `EPIList`,
`AmountRetention` → `RetentionList`.

**Yang BELUM terpecahkan tetap satu: tingkat pencatatannya.** Kepanjangan menjawab *apa*, bukan
*pada tingkat mana angkanya dicatat*. `PLA` tetap di daftar lima §7 dan tetap dijawab dengan satu
slip kontrak. Nama tetapnya ditentukan saat §10.4 ditulis ulang.

> **Dan itu memenuhi CONTEXT.md §3.3b**, yang mewajibkan setiap istilah pasar yang dipertahankan
> punya entri glosarium. Sebelumnya `PLA` adalah singkatan tanpa kepanjangan di mana pun — persis
> yang dilarang.

---

## 8. Bentuk pewarisan cabang — satu entitas dengan pembeda, kecuali satu pasang

**Aturan uji yang dipakai:** bedakan "**berbeda arti**" dari "**arti sama, hanya tidak dipakai di
salah satu cabang**". Hanya yang pertama membenarkan dua entitas.

| Kandidat | Hasil uji | Putusan |
|---|---|---|
| `BAGIAN` | kolom "berbeda arti" **kosong** | **satu entitas** |
| `POTONGAN` | persoalannya larut bersama BAGIAN | **satu entitas**, tanpa tabel penghubung |
| `PENYEBARAN` | kedalaman fan-out milik **mekanisme**, bukan cabang | **satu entitas** |
| `LAYER` vs `DETAIL_PROPORSIONAL` | **nol tumpang tindih atribut** | **dua entitas terpisah** |

**Bukti untuk `PENYEBARAN`:** `SetSpreadName` (proporsional manual) ber-`BreakDownSprd = 8`;
`SetSpreadingXOL` (XOL manual) ber-`BreakDownSprd = 9`; **kedua jalur otomatis 0**. Kedalaman
mengikuti mekanisme otomatis-versus-manual, bukan cabang proporsional-versus-XOL. Pemisahan runtuh.

**Koreksi yang saya cabut sendiri:** saya sempat membandingkan `Share[]` terhadap
`Limits[].Detail[]` dan menemukan nol tumpang tindih — tetapi itu membandingkan **bagian**
non-proporsional terhadap **ketentuan** proporsional. Perbandingan bagian-lawan-bagian yang benar
menghasilkan kolom "berbeda arti" yang kosong.

---

## 9. Keadaan siklus hidup — enam, ditambah satu keadaan warisan

Isi lengkapnya di **ADR-0055**. Ringkasnya:

```
DRAFT · MENUNGGU_SEC_HEAD · MENUNGGU_DEPT_HEAD · MENUNGGU_DIREKTUR
DISETUJUI (terminal) · DITOLAK (terminal) · WARISAN_TAK_TERPETAKAN (ADR-0054)
```

Dua belas perpindahan. Tidak ada keadaan yatim.

**Dua hal yang dibuktikan, bukan diduga:**

- `ReasTreatyInGroupLeader` punya **tiga perpindahan keluar dan nol masuk** — tidak ada satu pun
  aturan di seluruh ekspor yang menyetelnya. Keadaan yang tidak pernah terjadi.
- `DIKEMBALIKAN` **bukan keadaan**. Kata `Reject` muncul hanya di **dua berkas** pada seluruh
  ekspor, dan pada kondisi tombol ia diperlakukan **identik** dengan `Accept`. Tidak ada aturan
  kelengkapan, wewenang, maupun keterbukaan field yang pernah membacanya. Ia `DRAFT` yang punya
  riwayat penolakan, dan riwayat itu sudah ada di `CATATAN_PERSETUJUAN`.

Keadaan disimpan pada **versi kontrak**, bukan pada kontrak.

---

## 10. Penamaan, tipe konseptual, dan keterisian — **27 entitas penyerahan pertama**

> **Pemutakhiran 24 September 2026 — dua entitas menunggu di luar angka 27**, dan keduanya tidak
> diselipkan ke dalamnya diam-diam: **`PEMULIHAN_LIMIT`** (§10.3b, diterima, P-59 BARU) dan
> **`PERISTIWA_KONTRAK`** (§10.21a, baru, dari §10.20a). Bila keduanya berdiri, penyerahan pertama
> memuat **29** entitas. Angka 27 di seluruh berkas ini **tidak diubah** sampai keduanya disahkan,
> supaya tidak ada dua angka yang berlaku sekaligus.

**§10.1–§10.4 dikerjakan 23 September 2026** — empat entitas penentu.
**§10.5–§10.22 dikerjakan 24 September 2026** — **23 entitas selebihnya**, daftarnya dari
`5-tiket/DAFTAR-PEKERJAAN.md` §5.1 dan §5.2, **dan hanya itu**. `RETRO_KELUAR`, `PENCAPAIAN`,
`NILAI_SELISIH`, dan `BESARAN_DAPAT_DISESUAIKAN` **sengaja tidak dikerjakan** — alasannya §5.4 di
berkas yang sama.

### KEWAJIBAN yang berlaku pada SETIAP subbagian §10

Ditetapkan pemilik proses 24 September 2026, dan berlaku surut atas §10.1–§10.4:

> **Setiap entitas menyatakan KUNCI ALAMINYA dan LINGKUPNYA secara terang** — **Bentuk A** (unik di
> dalam versi) atau **Bentuk B** (unik di dalam induk langsung), memakai penggolongan
> `SPEC-INVARIAN.md` §4.1.
>
> Entitas yang **tidak punya** kunci alami **menyatakannya juga, beserta alasannya** — karena
> entitas tanpa kunci alami **tidak dapat dipadankan antar versi sama sekali**, dan itu
> **keputusan, bukan kekosongan**.

> **Dan setiap entitas yang punya BARIS TURUNAN menyatakan `SUMBU_REKONSILIASI`-nya** — **nama
> kolom** yang dipakai mengelompokkan saat baris turunannya dijumlahkan. Entitas yang tidak punya
> baris turunan menyatakan itu juga.

Sebabnya: entitas tanpa kunci alami yang tidak menyatakan dirinya begitu **terlihat persis sama**
dengan entitas yang lupa diberi kunci, dan bedanya baru ketahuan saat seam Adjustment tidak bekerja.

**Kewajiban kedua punya pembaca yang berbeda dari yang pertama.** Kunci alami dibaca seam
Adjustment; `SUMBU_REKONSILIASI` dibaca **sesi DDL saat menulis `GROUP BY`** pada *materialized
view* INV-47. Sejak `SPEC-INVARIAN.md` §4.4c, **INV-47 menunjuk ke §10 ini** untuk nama kolomnya —
rumusan sebelumnya tidak dapat gagal dan tidak dapat dikompilasi.
Lima entitas yang punya kunci alami tanpa nomor invarian ditemukan dengan cara ini — `SPEC-INVARIAN.md`
§4.1.

### URUTAN: tiga didahulukan

`BAGIAN`, `RINCIAN_PENYEBARAN`, dan `NILAI_PENYEBARAN` dikerjakan **lebih dulu**, dan karena itu
bernomor lebih dulu di sini. Alasannya bukan kepentingan melainkan ketergantungan: keduanya menahan
dua hal yang sudah dinyatakan selesai — seam Adjustment dan pesan kegagalan INV-47/INV-50/INV-51.
Uraiannya `5-tiket/DAFTAR-PEKERJAAN.md` §5.0a.

`PENYEBARAN` dan `POTONGAN` **ikut di blok yang sama meski bukan bagian dari ketiga itu**: rantai
`BAGIAN → PENYEBARAN → RINCIAN_PENYEBARAN → NILAI_PENYEBARAN` hanya dapat dinyatakan utuh, karena
kunci alami Bentuk B sebuah anak menyebut induknya.

### TINGKAT PENCATATAN paket uang — §7 memberi HITUNGAN, tidak pernah DAFTAR

> **Lubang yang baru terlihat saat §10 ditulis, dan dilaporkan di tempatnya muncul.**
>
> §7 menyatakan 43 paket uang tersimpan: **16 tingkat `TREATY_100_PERSEN`, 22 tingkat `BAGIAN_NURE`,
> 5 belum dapat ditentukan.** Ketiga angka itu **tidak pernah dituangkan sebagai daftar** — tidak di
> berkas ini, tidak di `PETA-TELUSUR-JSON.md`, tidak di mana pun. Disapu 24 September 2026.
>
> Akibatnya nyata, bukan kerapian: **INV-39** menuntut setiap paket bernilai `TREATY_100_PERSEN`
> atau `BAGIAN_NURE`, dan **INV-40** menuntut `PERSEN_BAGIAN_DIPAKAI` terisi pada yang bertingkat
> bagian. Keduanya *constraint*. Tanpa daftarnya, **migrasi tidak tahu nilai apa yang harus ditulis
> ke kolom itu**, dan sesi DDL akan menulis constraint yang tidak ada satu pun barisnya dapat
> memenuhinya dengan benar.
>
> **Di §10.5–§10.22 tingkatnya ditulis hanya bila ada sumber tertulis yang menyatakannya.** Yang
> tidak ada ditulis **`BELUM DITENTUKAN`** — bukan ditebak. Rekapnya §10.23.

Tipe konseptual memakai kelompok ADR-0003 yang diwarisi dari modul Claim Non Prop:

| Kelompok | Isi | Bentuk konseptual |
|---|---|---|
| **U1** | nilai uang | bilangan eksak presisi tinggi |
| **P1** | persentase dan porsi | bilangan eksak presisi tinggi |
| **P2** | tarif brokerage dan pajak | bilangan eksak presisi menengah |
| **K1** | kurs | bilangan eksak presisi tinggi |
| **C1** | cacah dan nomor urut | bilangan bulat |
| **T** | teks | teks |
| **D** | tanggal | tanggal |
| **E** | himpunan tertutup | enumerasi |
| **R** | rujukan ke entitas lain | kunci asing |

Angka presisinya **tidak ditetapkan di sini** — ia keputusan gerbang sesi DDL.

---

### 10.0 Empat keputusan tingkat atribut yang baru muncul saat menamai

Keempatnya terlihat hanya ketika daftar ditulis sebagai daftar, bukan sebagai hitungan.

#### a. Keempat pasangan bersufiks `P` MELEBUR jadi satu atribut

`BrokeragePercent`/`BrokeragePercentP`, `Exclusions`/`ExclusionsP`,
`SpecialConditions`/`SpecialConditionsP`, dan `RNMShare`/`RNMShareP` adalah **satu fakta yang
disimpan di dua kolom, dibedakan cabang**. Itu pola terlarang: kolom yang artinya bergantung nilai
kolom lain di baris yang sama.

Ia tidak perlu ada. `SIFAT_PROPORSI` di `KONTRAK` sudah menyatakan cabangnya, dan sebuah kontrak
hanya satu cabang. Sistem lama butuh dua kolom karena satu halaman clipboard melayani dua set tab;
model baru tidak.

**Empat pasangan menjadi empat atribut.** Delapan kolom menjadi empat.

> **Dikoreksi §14.3 dan §14.4.** `FacShare`/`FacultativeShare` ternyata **bukan pasangan**
> melainkan salinan, dan `AccountingMode`/`…NonProp` **belum diputuskan**. Yang melebur pasti
> ada **empat**.

> **Ini keputusan saya.** Yang membatalkannya: satu baris data yang mengisi **kedua** varian
> sekaligus dengan nilai berbeda. Ujinya murah, dan ditambahkan sebagai **Uji X**.

Pasangan `AccountingMode`/`AccountingModeNonProp` berbentuk sama, hanya berakhiran lain. Ia ikut
melebur, dan Uji X mencakupnya.

#### b. Empat kolom bernama bahaya adalah SATU DAFTAR BATAS PER BAHAYA

`Earthquake`, `FloodJab`, `FloodNation`, `RSMDLimit` — masing-masing berpasangan dengan
`CurrencyEarthquake`, `CurrencyFloodJab`, `CurrencyFloodNat`, `CurrencyRSMD`.

Delapan kolom untuk **empat bahaya yang namanya menjadi nama kolom**. Bentuk itu akan menuntut
kolom kesembilan begitu ada bahaya baru — dan bahaya baru adalah hal yang pasti terjadi di
reasuransi.

**Bahayanya adalah nilai, bukan nama kolom.** Delapan kolom menjadi satu entitas
`BATAS_PER_BAHAYA`: rujukan ke versi, rujukan ke **bahaya** (tabel acuan, karena ia bertambah tanpa
mengubah arti apa pun — ADR-0038), dan satu paket uang.

`RSMD` istilah pasar yang dipertahankan dan **wajib punya entri glosarium** (CONTEXT.md §3.3b).
Kepanjangannya belum tercatat di mana pun; ia masuk daftar wawancara.

#### c. Dua atribut mata uang milik agregat yang tidak disimpan

`CurrencyEgnpiAmount` berpasangan dengan `TotalEgnpiAmount`, dan `CurrencyInstallmentAmount` dengan
`TotalInstallmentAmount`. Kedua nilainya **agregat turunan yang tidak disimpan** (§4.2), jadi mata
uangnya juga tidak. **Dibuang, dua atribut.**

#### d. Di sinilah selisih 187 ke 163 sebagian besar berada

| Sebab | Atribut hilang dari kepala |
|---|---|
| peleburan pasangan bersufiks `P` dan `NonProp` | 5 |
| batas per bahaya menjadi entitas tersendiri | 7 |
| mata uang milik agregat turunan | 2 |
| `PositionUsername` dibuang (§12.5) | 1 |
| pasangan nilai + mata uang menjadi satu paket uang di entitas anak | sisanya |

Peleburan paket uang **tidak menghilangkan informasi**: mata uang tetap ada, ia menjadi bagian
paket, bukan atribut berdiri sendiri.

---

### 10.1 `KONTRAK` — lapisan beku, 8 atribut

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_KONTRAK` | *baru* | C1 | tidak | pengenal buatan sistem; bukan dari teks, bukan dari cap waktu |
| `NOMOR_KONTRAK_WARISAN` | `ID` | T | **ya** | hanya terisi pada baris hasil migrasi (ADR-0042) |
| `ID_KONTRAK_DISALIN_DARI` | `OLDID` pecahan 1 | R | **ya** | rujukan ke `KONTRAK` lain |
| `ID_CEDANT` | `CedingID` | R | tidak | nama cedant **tidak disalin** |
| `ID_ASAL_BISNIS` | `LeadingReinsSourceID` | R | tidak | *source of business* |
| `SIFAT_PROPORSI` | `ProportionType` | E | tidak | `PROPORSIONAL` / `NON_PROPORSIONAL` |
| `TANGGAL_MULAI` | `Commencement` | D | tidak | batas **inklusif** (ADR-0022) |
| `TANGGAL_BERAKHIR` | `Termination` | D | tidak | batas **inklusif** |

**`Ceding` dan `LeadingReinsSource` (nama) tidak disimpan.** Keduanya salinan dari master, dan
menyimpannya berarti dua penulis untuk satu fakta (ADR-0041). Apakah rujukannya kunci asing lintas
skema ke `POOLDATA` atau nilai yang divalidasi aplikasi adalah **keputusan sesi DDL**, bukan
keputusan model.

**Kunci alami dan lingkupnya** — `ID_CEDANT` + `ID_ASAL_BISNIS` + `TANGGAL_MULAI` +
`TANGGAL_BERAKHIR` + `SIFAT_PROPORSI`, berlingkup **seluruh skema**. **Bukan Bentuk A maupun B**:
keduanya menggambarkan keunikan *di dalam* sesuatu, sedangkan `KONTRAK` tidak punya induk.

Dan ia **memperingatkan, tidak melarang** (ADR-0040 §2) — karena itu ia **bukan invarian**, dan
ketiadaan nomornya di INV-05…INV-16 disengaja. Lihat `SPEC-INVARIAN.md` §4.1.

---

### 10.2 `VERSI_KONTRAK` — **45 atribut** (43 di sini + 1 di §10.2a + 1 di §10.19a)

> ### `F-2` DITUTUP 24 September 2026 — angka judulnya yang salah, bukan tabelnya
>
> **Judul bagian ini pernah berbunyi *"47 atribut sesudah peleburan, 48 sejak §10.2a"*.** Angka itu
> **tidak pernah dapat diturunkan dari apa pun** dan kini diganti dengan yang dihitung: **45**.
>
> | Dari mana | Berapa |
> |---|---:|
> | tabel-tabel di §10.2 sendiri | 43 |
> | §10.2a `SIFAT_MATERIAL_ADDENDUM` | 1 |
> | §10.19a — `ID_DOKUMEN_ADDENDUM`, ditambahkan diff `D-2` | **1** |
> | **jumlah** | **45** |
>
> **Baris ketiga itu baru diketahui saat putaran ini menghitung, dan ia lubang tersendiri —
> `F-17`.** Lihat kotak di bawah.
>
> **Selisihnya SATU, bukan dua.** "Kurang 3" dan "kurang 4" adalah jarak yang sama diukur terhadap
> dua angka di dalam satu judul: 47 tanpa §10.2a, 48 dengan §10.2a. Pengurai membaca angka pertama,
> pembaca manusia membaca yang kedua, dan keduanya melaporkan lubang yang berbeda besarnya untuk
> lubang yang sama.
>
> #### `F-17` — kolom yang mendarat di tabel yang salah, dan hilang dari tabel yang benar
>
> §10.19a memuat **dua** tabel: atribut `DOKUMEN_ADDENDUM`, lalu — sesudah kalimat *"Dan pada
> **`VERSI_KONTRAK`** ditambahkan satu atribut"* — sebuah tabel berisi satu baris,
> `ID_DOKUMEN_ADDENDUM`, yang **milik `VERSI_KONTRAK`**.
>
> Pengurai `alat/buat-kamus-dan-ddl.py` memberi seluruh tabel di dalam sebuah subbagian kepada
> entitas judulnya. Akibatnya **dua** sekaligus, dan keduanya diam:
>
> | Akibat | Yang terlihat sebelum putaran ini |
> |---|---|
> | `DOKUMEN_ADDENDUM` memperoleh kolom **kedua** bernama `ID_DOKUMEN_ADDENDUM` | selisih **+1** di `KAMUS-KOLOM.md` §0 — terbaca sebagai *"judulnya menulis 3, tabelnya memuat 4"*, seolah cacah judulnya yang salah |
> | `VERSI_KONTRAK` **kehilangan** kunci asingnya ke `DOKUMEN_ADDENDUM` | **tidak terlihat sama sekali.** `ddl-usulan/11_VERSI_KONTRAK.sql` berdiri tanpa relasi itu, dan tidak ada yang memeriksanya |
>
> **Yang kedua jauh lebih berbahaya daripada yang pertama**, dan itu polanya: selisih cacah
> **berbunyi**, relasi yang hilang **diam**. Diff `D-2` diterapkan dengan benar ke §10 — yang gagal
> adalah pembacaannya.
>
> **Ditutup di perkakas, bukan dengan memindahkan tabelnya.** §10.19a adalah tempat yang benar untuk
> baris itu: ia bagian dari keputusan `GRL-19` dan harus terbaca bersamanya. Yang diperbaiki
> pengurainya — ia kini mengenali kalimat *"Dan pada **`X`** ditambahkan"* sebagai perpindahan
> entitas, dan **mencetak setiap perpindahan yang dipakainya**, sehingga aturan ini tidak dapat
> bekerja diam-diam.
>
> | | |
> |---|---|
> | **Label** | **KOREKSI PERKAKAS** — spesifikasinya tidak berubah satu huruf pun |
> | **Arah dampak bila salah** | bila kalimat pemindah itu dipakai di tempat lain untuk maksud yang berbeda, kolom berikutnya akan berpindah ke entitas yang salah. Karena itu perkakas **mencetak daftar perpindahannya**, dan daftar itu diperiksa setiap kali ia dijalankan — sekarang isinya **satu** |
> | **Syarat pembalikan** | §10 berhenti memakai bentuk kalimat itu, atau memakainya untuk sesuatu yang bukan pemindahan entitas |
>
> #### Apa yang dikerjakan sebelum angkanya diturunkan
>
> Menurunkan angka judul adalah tindakan yang **menghapus jejak**, dan itu sebab putaran sebelumnya
> menolak melakukannya: *"begitu ia berbunyi 44, tidak akan ada seorang pun yang mencari empat yang
> hilang."* Maka pencariannya **diselesaikan lebih dulu**, bukan dilewati.
>
> | Langkah | Hasil |
> |---|---|
> | sumber **mandiri** dari §3.2 dan dari §10 — `4-erd-dan-tabel-datar/datar-treatyin-lama.csv`, 985 simpul dari sapuan 329 berkas | 99 skalar kedalaman-1 `TreatyIn` dibandingkan terhadap 131 nama *Asal* di seluruh tabel §10 |
> | yang tidak tercocokkan | **42** |
> | sudah punya sebab tertulis sebelum putaran ini | 37 |
> | **diadili putaran ini** | **5** — `TreatyYear`, `RevisionDate`, `CoInScale`, `Ceding`, `LeadingReinsSource` |
> | **yang menjadi atribut `VERSI_KONTRAK`** | **NOL.** Rinciannya `2-to-spec/COCOK-SILANG-CACAH-ATRIBUT.md` §5 |
>
> **Tidak satu pun dari keempatnya ditemukan, dan tidak satu pun dikarang.** Yang ditemukan adalah
> bahwa **tidak ada calon yang tersisa** di dalam semesta yang terjangkau.
>
> #### Yang TIDAK ikut tertutup, dan ia harus tetap terbaca
>
> Semesta pencarian itu **bolong**, dan bolongnya bernama: **340 properti titik buta belum diperiksa
> siapa pun** (`L-8`, ditagih `M-4`). Maka kalimat yang berlaku sesudah putaran ini **bukan**
> *"`VERSI_KONTRAK` lengkap"*, melainkan:
>
> > **Tidak ada atribut yang hilang di antara yang pernah dilihat. Yang belum pernah dilihat ada
> > 340, dan ia milik gelombang 2.**
>
> | | |
> |---|---|
> | **Label** | **KOREKSI ANGKA** — bukan pelestarian, bukan perubahan perilaku |
> | **Arah dampak bila salah** | bila keempatnya ternyata memang ada dan terbaca dari titik buta, mereka **ditambahkan sebagai atribut baru bernomor temuan**, bukan dipulihkan diam-diam. Tidak ada data yang hilang karena keputusan ini — belum ada data |
> | **Syarat pembalikan** | penutupan `L-8` menemukan skalar kedalaman-1 `TreatyIn` yang tidak punya rumah. Angka judul lalu naik **beserta barisnya**, tidak pernah tanpa barisnya |
>
> #### Angka 47 dikutip di tempat lain — disebutkan, tidak disunting dari sini
>
> | Berkas | Bunyinya |
> |---|---|
> | `4-erd-dan-tabel-datar/AUDIT-PENYEBUT-POHON.md` | *"(`VERSI_KONTRAK`, 47 atribut) tidak tersentuh sama sekali"* dan satu baris tabel |
> | `4-erd-dan-tabel-datar/PRA-PEMISAHAN-ANGKA-DAN-TABRAKAN.md` | *"46 dari 47 atribut identik"*, dan `T-1` yang menghitung *"membuang 35 atribut"* (20 versus 8+47) |
> | `PENGETAHUAN-MIGRASI-TREATY-IN.md` | dua kali — *"`VERSI_KONTRAK` 1/47"* dan *"`VERSI_KONTRAK` (47)"* |
>
> **Keempatnya tidak disunting di putaran ini**, dan itu keputusan: ketiganya berkas **rekaman**
> putaran lain, dan mengubah angka di dalam rekaman membuat rekaman itu berhenti merekam apa yang
> waktu itu diketahui. Yang berubah hanya **sumbernya**. Siapa pun yang mencocokkan angka **wajib
> mulai dari sini**, sesuai urutan wewenang.
>
> Satu di antaranya berubah **artinya**, dan itu perlu disebut: `T-1` menghitung *"membuang 35
> atribut"*; dengan 44 angkanya menjadi **32**. Arah temuannya tidak berubah.

**Identitas dan keadaan**

| Nama | Asal | Tipe | Boleh kosong |
|---|---|---|---|
| `ID_VERSI_KONTRAK` | *baru* | C1 | tidak |
| `ID_KONTRAK` | *baru* | R | tidak |
| `NOMOR_URUT_VERSI` | *baru* | C1 | tidak |
| `KEADAAN_SIKLUS_HIDUP` | `Position` + `StatusAkseptasi` + 3 bendera | E | tidak | ⚠ **DELAPAN nilai** — `DIBATALKAN` ditambahkan ADR-0055 perubahan 24 Sep 2026 |
| `KEADAAN_WARISAN_ASLI` | *baru* | T | **ya** — hanya bila keadaannya `WARISAN_TAK_TERPETAKAN` |
| `JENIS_ADDENDUM` | `EDMState` | E | **ya** — kosong pada versi pertama |
| `TANGGAL_BERLAKU_ADDENDUM` | `EDMEffective` | D | **ya** — idem |

**Penamaan dan lingkup**

| Nama | Asal | Tipe | Boleh kosong |
|---|---|---|---|
| `NAMA_KONTRAK` | `TreatyContractName` | T | tidak |
| `LINGKUP_WILAYAH` | `TeritorialScope` | T | ya |
| `KELAS_BISNIS_KONTRAK` | `ClassofBusiness` | R | ya |
| `KODE_MATA_UANG_KONTRAK` | `Currency` | R | tidak |
| `ID_REASURADUR_PEMIMPIN` | `LeadingReinsID` | R | ya |
| `ID_KETUA_TREATY` | `TreatyLeader` | R | ya |

`LeadingReinsName` **tidak disimpan** — nama perusahaan dibaca dari master.

**Bagian NuRe dan biaya** *(sesudah peleburan §10.0a)*

| Nama | Asal | Tipe | Boleh kosong |
|---|---|---|---|
| `PERSEN_BAGIAN_NURE` | `RNMShare` + `RNMShareP` | P1 | tidak |
| `BAGIAN_NURE_SERAGAM` | `RNMShareAcrossTheBoard` | E | tidak |
| `PERSEN_BAGIAN_NURE_DIPOTONG` | `RnmShareDeducted` | P1 | ya |
| `PERSEN_BROKERAGE` | `BrokeragePercent` + `…P` | P2 | ya |
| `PERSEN_BAGIAN_FAKULTATIF` | `FacultativeShare` (`FacShare` dibuang, salinan — §14.3) | P1 | ya |
| `PERSEN_BROKERAGE_FAKULTATIF` | `FacultativeShareBrokerage` (`FacShareBrokerage` dibuang, salinan) | P2 | ya |

**Ketentuan tertulis**

| Nama | Asal | Tipe | Boleh kosong |
|---|---|---|---|
| `PENGECUALIAN` | `Exclusions` + `ExclusionsP` | T | ya |
| `KETENTUAN_KHUSUS` | `SpecialConditions` + `SpecialConditionsP` | T | ya |
| `KETERANGAN` | `Information` | T | ya |
| `CATATAN` | `Comment` | T | ya |
| `CATATAN_BORDEREAUX` | `BordereauxNote` | T | ya |
| `MEMAKAI_BORDEREAUX` | `Bordeaux` | E | tidak |

> `Bordeaux` hampir pasti **salah eja** `Bordereaux`, dan keduanya berdampingan di kepala yang sama.
> Nama barunya memakai ejaan yang benar; ejaan lamanya tercatat di peta telusur.

**Cara pembukuan dan pelaporan**

| Nama | Asal | Tipe | Boleh kosong |
|---|---|---|---|
| `CARA_PEMBUKUAN` | `AccountingMode` | E | tidak |
| `CARA_PEMBUKUAN_XOL` | `AccountingModeNonProp` | E | ya — **menunggu Uji X-2** (§14.4) |
| `PERIODE_PELAPORAN_KONTRAK` | `ReportingPeriod` | E | ya |
| `SELANG_PELAPORAN` | `ReportingInterval` | C1 | ya |
| `TANGGAL_MULAI_PELAPORAN` | `ReportingStart` | D | ya |
| `TANGGAL_AKHIR_PELAPORAN` | `ReportingEnd` | D | ya |
| `HARI_BATAS_PENYERAHAN` | `ReportingSubmission` | C1 | ya |
| `HARI_BATAS_KONFIRMASI` | `ReportingConfirmation` | C1 | ya |
| `HARI_BATAS_PELUNASAN` | `ReportingSettlement` | C1 | ya |
| `HARI_PENGINGAT` | `ReminderDays` | C1 | ya |
| `PERIODE_AKUMULASI_KONTRAK` | `AccumulationPeriod` | E | ya |

**Termin dan prorata**

| Nama | Asal | Tipe | Boleh kosong |
|---|---|---|---|
| `JUMLAH_TERMIN` | `InstallmentNo` | C1 | ya |
| `MEMAKAI_PRORATA` | `IsProRate` | E | tidak |

**Kapasitas dan batas**

| Nama | Asal | Tipe | Boleh kosong |
|---|---|---|---|
| `BATAS_MAKSIMUM_KELOMPOK` | `MaxCoGroup` | U1 paket uang | ya |
| `BATAS_MAKSIMUM_NON_KELOMPOK` | `MaxCoNonGroup` | U1 paket uang | ya |
| `BATAS_PILIHAN` | `OptionLimit` | U1 paket uang | ya |

*(`Earthquake`, `FloodJab`, `FloodNation`, `RSMDLimit` beserta keempat mata uangnya **pindah** ke
entitas `BATAS_PER_BAHAYA` — §10.0b.)*

**Retro dan status master**

| Nama | Asal | Tipe | Boleh kosong |
|---|---|---|---|
| `RETRO_BERGANDA` | `IsMultipleRetro` | E | tidak |
| `ID_DAFTAR_RETRO` | `RetroList` | R | ya |

> `CedingStatusActive` dan `SourceStatusActive` **dibuang** — salinan status data master, dan
> keduanya tidak pernah ditulis satu aturan pun sehingga tidak mungkin fakta beku. Lihat §14.5.

**Kunci alami dan lingkupnya** — `NOMOR_URUT_VERSI`, unik di dalam `ID_KONTRAK` (INV-04). Pada
penggolongan §4.1 ia **Bentuk B**: lingkupnya **induk langsung**, dan induk langsung
`VERSI_KONTRAK` adalah `KONTRAK`. Ia tidak dapat berbentuk A, karena A berarti "unik di dalam
versi" dan entitas ini **adalah** versinya.

**Tiga paket uang di entitas ini belum punya tingkat pencatatan tertulis** —
`BATAS_MAKSIMUM_KELOMPOK`, `BATAS_MAKSIMUM_NON_KELOMPOK`, `BATAS_PILIHAN`. Lihat §10.23.

---

#### 10.2a `SIFAT_MATERIAL_ADDENDUM` — atribut **ke-44**, ditambahkan 24 September 2026

> ### DIFF `D-1` **DITERAPKAN** 24 September 2026 — parkirnya dicabut oleh penutupan `F-2`
>
> | | |
> |---|---|
> | **Apa yang berubah** | **catatan** baris di bawah, bukan barisnya. Dasarnya berpindah dari `GRL-12` — yang **BATAL** — ke **`GRL-20`**, yang mengembalikan materialitas sebagai **atribut masukan** |
> | **Kenapa ia sempat diparkir** | `D-1` menyebut baris ini *"atribut ke-48"*. Selama `F-2` terbuka, menerapkannya **membekukan angka yang salah ke dalam catatan sebuah baris** — dan angka di dalam catatan adalah angka yang dikutip |
> | **Apa yang mencabut parkirnya** | **`F-2` ditutup** di blok kepala §10.2: pencarian keempat atribut diselesaikan atas sumber mandiri, hasilnya **nol calon tersisa**, dan angka judul diturunkan ke **44** yang dihitung. Nomor urut baris ini karenanya **ke-44**, dan itu angka yang dapat dipertanggungjawabkan |
> | **Isi diffnya** | **tidak diubah.** Yang ditambahkan hanya nomor urut yang benar — syarat yang `D-1` sendiri pasang |
>
> **Akibatnya pada modul Adjustment:** butir `D-1` di
> `treaty-in-adjustment/2-to-spec/USULAN-DIFF-KE-INDUK.md` **berhenti diparkir**. Tujuh diff
> Adjustment kini **enam diterapkan, satu diparkir** (`D-7`, menunggu satu kalimat izin).

`EDMMaterialType` **keluar dari daftar ditunda (§6)** dan menjadi atribut masukan di entitas ini.

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `SIFAT_MATERIAL_ADDENDUM` | `EDMMaterialType` | E | **ya** — kosong pada versi pertama | dua nilai **`MATERIAL`** / **`TIDAK_MATERIAL`**. **Masukan, bukan turunan** (`GRL-20`; `GRL-12` BATAL). Beku sejak `AJUKAN`, berjejak selama `DRAFT` (`KTV-1` Adjustment) |

**Kenapa ia MASUKAN dan bukan turunan, dan kenapa itu tidak menyentuh ADR-0049:**

> Di sistem lama `EDMMaterialType` **dipilih saat addendum dibuat** — ia parameter
> `Param.MaterialType` pada `TreatyInEDMSetValue` — dan ia **menjaga ruas mana yang boleh disunting
> di layar** (brokerase, bagian seragam, prorata, penambahan baris non-proporsional dimatikan saat
> nilainya `2`).
>
> **Sesuatu yang menentukan apa yang boleh diubah tidak mungkin dihitung dari apa yang berubah.**
> Maka ia bukan materialitas yang dimaksud ADR-0049, melainkan **pilihan pembuatnya di muka**.

Keduanya berdampingan tanpa bertabrakan, dan bedanya ditulis di sini supaya tidak tertukar lagi:

| | Apa | Kapan ada | Disimpan |
|---|---|---|---|
| **`SIFAT_MATERIAL_ADDENDUM`** *(atribut ini)* | **pilihan pembuat** — menentukan ruas mana yang terbuka | **sebelum** perubahan | **ya** |
| **materialitas** *(ADR-0049)* | **turunan** — ada-tidaknya baris selisih | **sesudah** perubahan | **tidak** |

> ### ⚠ TABRAKAN NAMA — dilaporkan, tidak diselesaikan diam-diam
>
> Nama yang diusulkan untuk atribut ini adalah **`JENIS_ADDENDUM`**. **Nama itu sudah dipakai di
> entitas yang sama**, bersumber pada properti yang **berbeda**:
>
> | Atribut §10.2 | Asal | Nilainya di sistem lama |
> |---|---|---|
> | `JENIS_ADDENDUM` *(sudah ada)* | **`EDMState`** | `1` Internal Edit · `2` External Addendum · `3` Addendum Premium |
> | `SIFAT_MATERIAL_ADDENDUM` *(ini)* | **`EDMMaterialType`** | `1` · `2` |
>
> **Keduanya properti yang berbeda dan hidup berdampingan di halaman yang sama.** Bukti pemetaan
> `EDMState`: `DataTransform/TreatyInSetEditPre.xml` mencabangkan `EDMState` 1/2/3 dan menuliskan
> ketiga kalimat *"Had Created …"* yang berbeda. Golongan **`EVIDENCED`**.
>
> **Maka nama `JENIS_ADDENDUM` tidak dipakai ulang**, dan atribut ini diberi nama yang menyebut
> dirinya sendiri. Bila pemilik proses menghendaki nama lain, yang diganti **nama atribut ini**,
> bukan pemetaan `EDMState` — karena pemetaan itu terbukti dari ekspor.

**Satu hal yang BELUM terbukti dan tidak ditebak: mana dari nilai `1` dan `2` yang berarti
material.** Yang terbaca: `EDMMaterialType = 2` **mematikan** ruas-ruas angka uang, dan
`EDMMaterialType == '1'` **menampilkan** sesuatu pada addendum bukan-premi. Pembacaan yang wajar
adalah `1` = `MATERIAL` dan `2` = `TIDAK_MATERIAL`, **tetapi itu pembacaan, bukan pemetaan yang
tertulis di mana pun.** Penagihnya: **saat tabel acuan enumerasi ditulis di sesi DDL**, dan ia
ditutup oleh satu pertanyaan ke bagian teknik atau oleh **Uji AN** *(sebaran `EDMMaterialType`
disandingkan dengan ada-tidaknya baris `ValueDifference`)*.

§16: `SIFAT_MATERIAL_ADDENDUM` = 23 bita. Muat.

---

### 10.3 `LAYER` — **21 atribut** (13 + 6 yang tidak pernah terlihat + 2 pulih dari `L-8`, `TDA-17`)

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_LAYER` | `Limits[].ID` | C1 | tidak | |
| `ID_VERSI_KONTRAK` | *baru* | R | tidak | |
| `NOMOR_LAYER` | `Layer` | C1 | tidak | **kunci alami** di dalam versi |
| `BAGIAN_LAYER` | `LayerPart` | C1 | ya | bagian dari kunci alami |
| `JENIS_LAYER` | `LayerType` | E | tidak | |
| `JENIS_BAGIAN_LAYER` | `LayerPartType` | E | ya | |
| `CAKUPAN` | `Cover` | T | ya | |
| `LIMIT` | `Limit` + `Currency` | U1 paket uang | tidak | tingkat **100% treaty** |
| `DEDUCTIBLE` | `Deductible` | U1 paket uang | tidak | tingkat **100% treaty** |
| `PERSEN_PENYESUAIAN` | `AdjRate` | P1 | ya | |
| `PERSEN_MINIMUM_DEPOSIT` | `MDPPct` | P1 | ya | |
| `PORSI_PEMULIHAN_LIMIT` | `ReinstatementPct` | P1 | ya | porsi limit yang dipulihkan |
| `TARIF_PREMI_PEMULIHAN` | `ReinstatementValue` | P1 | ya | tarif premi untuk pemulihan itu — **pembacaannya punya saingan**, §10.3b |
| `LIMIT_AGREGAT` | `AgregateLimit` *(salah eja ada di sumbernya)* | U1 paket uang | ya | **BARU** — batas total sepanjang periode, terpisah dari limit per kejadian. Tingkat **BELUM DITENTUKAN** |
| `MDP` | `MDPList[]` | U1 paket uang | ya | deposit premium; tingkat **BELUM DITENTUKAN** |
| `MDP_MINIMUM` | `MDPMinList[]` | U1 paket uang | ya | **BARU** — tanpa ia, *"deposit premium sama dengan minimum premium"* tidak dapat dinyatakan |
| `PERSEN_MDP_MINIMUM` | `MDPMinPct` | P1 | ya | **BARU** |
| `MDP_DIGABUNG` | `IsCombineMDP` | E | tidak | **BARU** — MDP dihitung per layer atau gabungan |
| `TANPA_HITUNG_PREMI_PEMULIHAN` | `NoRIPCalculation` | E | tidak | **BARU** — mematikan perhitungan premi pemulihan pada layer ini |
| `DEDUCTIBLE_KEDUA` | `Limits[].Deductible2` | U1 | ya | **PULIH dari `L-8`** — `TDA-17`. Tersimpan sebagai **teks** di sistem lama (`@toDecimal` di kedua sisi rumus selisih); migrasi wajib mengubah tipe, dan baris yang gagal konversi **tidak disamarkan jadi nol** (ADR-0035) |
| `PERSEN_ROL` | `Limits[].ROLPct` | P2 | ya | **PULIH dari `L-8`** — `TDA-17`. Rate on line. ⚠ **DIBANTAH `F-15`** — pertanyaan *"selalu dihitung atau pernah disepakati"* **terjawab dari sumber**: `DetailCalculationROL` menghitungnya `premi ÷ limit × 100`. Kolomnya **tidak dicabut diam-diam**; keputusannya pemilik proses, `2-to-spec/F-15-PREMI-DIPEROLEH-DAN-ROL-TURUNAN.md` |

**Kunci alami dan lingkupnya** — `NOMOR_LAYER` + `BAGIAN_LAYER`, **Bentuk A**: unik di dalam
`ID_VERSI_KONTRAK` (INV-05).

**`SUMBU_REKONSILIASI`** — `ID_LAYER` pada entitas anaknya. `LAYER` sendiri **punya** baris turunan
(`BAGIAN`, dan lewat `BAGIAN` seluruh rantai penyebaran), sehingga jumlah yang direkonsiliasi
INV-47 dikelompokkan menurut `ID_LAYER`.

#### 10.3a Enam belas properti yang tidak pernah terlihat — diadili

**`LAYER` diselesaikan atas sumber yang kurang**, bukan dibuka kembali: daftar 13 atribut di atas
disusun dari pohon, dan pohon tidak pernah melihat 16 properti kelas `Data-TreatyInLimits` (L-8).
Tiga di antaranya agregat `Sum…` dan langsung jatuh ke turunan. Tiga belas diadili di sini.

| Properti | Putusan | Dasar |
|---|---|---|
| `AgregateLimit` *(salah eja ada di sumbernya)* | **ATRIBUT BARU — `LIMIT_AGREGAT`, paket uang** | lihat kotak di bawah |
| `AgregateLimit2` | **dibuang** | kembaran mata uang kedua, INV-46 |
| `CurrencyRelation` | **tidak disimpan** | sudah diputuskan ADR-0053 — ia pasif |
| `PremiumEarned` | **turunan** | kembaran skalar dari `PremiumEarnedList`, yang sudah DITURUNKAN di peta telusur |
| `MDPMinList` · `MDPMinPct` | **ATRIBUT BARU — minimum deposit premium**: satu paket uang + satu persentase | berpasangan dengan `MDPList`/`MDPPct` yang sudah ada. Tanpa keduanya, "deposit premium sama dengan minimum premium" **tidak dapat dinyatakan**, padahal ia andaian yang menopang lantai premi (`PENGETAHUAN.md`) |
| `IsCombineMDP` | **ATRIBUT BARU — `MDP_DIGABUNG`, enumerasi** | menentukan apakah MDP dihitung per layer atau gabungan; tanpa ia, migrasi tidak dapat memilih rumus |
| `NoRIPCalculation` | **ATRIBUT BARU — `TANPA_HITUNG_PREMI_PEMULIHAN`, enumerasi** | mematikan perhitungan premi pemulihan pada layer ini |
| `LayerList` · `AdditionalPct` · `ReinstatementAmount1` · `ReinstatementAmount2` · `ReinstatementNote` | **DITAHAN — keputusan pemulihan limit, kotak berikutnya** | kelimanya hidup di dalam `Reinstatement_List` |

> **`AgregateLimit` hilang tanpa disengaja, dan begini caranya.** Daftar di bawah membuang
> `AggregateLimit2` sebagai kembaran mata uang kedua — **benar** — tetapi **kembaran pertamanya
> tidak pernah ikut dibawa masuk**, karena pohon tidak melihatnya. Sebuah batas agregat adalah
> ketentuan pokok kontrak non-proporsional: ia membatasi total tanggungan sepanjang periode,
> terpisah dari limit per kejadian. **Menghilangkannya mengubah arti kontrak.**
>
> Ini bentuk kegagalan yang pantas dicatat: **membuang kembaran lebih mudah daripada memastikan
> yang aslinya terbawa.** Setiap baris "dibuang — kembaran mata uang kedua" di berkas ini diperiksa
> ulang karena temuan ini, dan keempat lainnya (`Limit2`, `Deductible2`, `Currency2`, `MDP2`)
> **aslinya memang ada** di daftar atribut.

#### 10.3b PEMULIHAN LIMIT ADALAH DAFTAR, BUKAN DUA SKALAR — diangkat, tidak diputuskan sendiri

`Limits.Reinstatement_List` adalah **Page List** berkelas `Data-TreatyInLimits` sendiri — kelas yang
bersarang pada dirinya — dan `Activity/SetReinstatementPct.xml` **menambah baris** ke dalamnya:

```
.Reinstatement_List(<APPEND>).ReinstatementPct   = "100"
.Reinstatement_List(<LAST>).ReinstatementValue   = Local.runningidx + 1
.Reinstatement_List(<LAST>).AdditionalPct        = "100"
.Reinstatement_List(<LAST>).ReinstatementNote    = .ReinstatementNote
.Reinstatement_List(<LAST>).ReinstatementAmount1 = .Limit
```

Tiga hal terbaca sekaligus, dan ketiganya baru:

1. **Pemulihan limit berulang, dan tiap kejadian punya barisnya sendiri.** INV-49 sudah menyatakan
   *"reinstatement boleh terjadi lebih dari sekali"*; yang belum pernah terlihat adalah bahwa
   sistem lama **memang menyimpannya sebagai baris**.
2. **Seluruh ketentuannya disemai `100`.** `ReinstatementPct` dan `AdditionalPct` **ditulis tetap
   `"100"`**, bukan diambil dari masukan. Maka sistem lama **tidak dapat menyatakan** ketentuan
   pasar yang biasa — *"pemulihan pertama 100%, kedua 50%"*. Ini instans baru dari pola **cacat
   bersembunyi di balik nilai bawaan**: selama semuanya 100, satu baris dan sepuluh baris
   menghasilkan angka yang sama.
3. **`ReinstatementValue` berarti DUA hal di kelas yang sama.** Di dalam daftar ia
   `runningidx + 1` — **nomor urut pemulihan**. Di baris `Limits` induknya, satu-satunya
   penulisnya adalah `TreatyInMappingDataconvert.xml` dengan `= .Reinstatement`. Nama yang sama,
   dua arti, dibedakan **kedalaman** — persis pola yang dilarang.

> **Akibatnya pada §10.3 di atas: `TARIF_PREMI_PEMULIHAN` (dari `Limits[].ReinstatementValue`)
> berdiri di atas pembacaan yang kini punya saingan.** Saya **tidak** mengubahnya sendiri —
> pembacaan "tarif premi" berasal dari temuan grilling tentang rasio limit yang menimpa tarif, dan
> menggantinya dari satu penulis konversi akan menukar satu bukti dengan bukti lain tanpa
> menimbangnya.

**Yang diminta keputusan**, karena ia menambah entitas dan daftar entitas mengikat:

| Pilihan | Ongkos bila salah |
|---|---|
| **entitas anak `PEMULIHAN_LIMIT`** di bawah `LAYER`, satu baris per kejadian | satu tabel yang barisnya selalu seragam — **murah**, dan ia merosot dengan sendirinya menjadi bentuk sekarang |
| tetap dua skalar pada `LAYER` | bila bisnis memakai ketentuan bertingkat, ia **perubahan skema + migrasi ulang** — **mahal** |

Ongkosnya **tidak setangkup**, dan aturan yang sudah dipakai berkali-kali di proyek ini memilih sisi
murah tanpa menunggu bukti.

> ### DIPUTUSKAN 24 September 2026 — entitas `PEMULIHAN_LIMIT` BERDIRI
>
> Dasarnya bukan bahwa datanya berbentuk daftar, melainkan pertanyaan pokok: **apakah "setiap
> pemulihan bersyarat sama" itu KEHENDAK BISNIS atau AKIBAT cara sistem lama dibangun?**
>
> `"100"` dan `"100"` adalah **satu baris kode**, bukan kebijakan yang dinyatakan di mana pun.
> Sementara di reasuransi non-proporsional, pemulihan dengan syarat berbeda — yang pertama
> cuma-cuma, berikutnya prorata — adalah hal **biasa**, bukan aneh.
>
> Satu baris per pemulihan, masing-masing membawa syaratnya sendiri.

**Dua hal yang menyertainya, dan keduanya mengikat.**

**(a) Kemampuan "syarat berbeda tiap pemulihan" bergolongan BARU.** Sistem lama **tidak dapat**
menyatakannya. Yang **PELESTARIAN** hanyalah menyimpan pemulihan sebagai daftar; yang **BARU**
adalah nilainya boleh berbeda-beda. Ia menuntut barisnya sendiri di `5-tiket/DAFTAR-PEKERJAAN.md`
dan medan **`CARA MENYALAKANNYA`** — kemampuan **P-59** *(bernomor P-58 sampai 24 Sep 2026; lihat `5-tiket/DAFTAR-PEKERJAAN.md` §2.2b)*.

**(b) Nama `ReinstatementValue` TIDAK DIBAWA.** Ia berarti dua hal di kelas yang sama: nilai di satu
tempat, `runningidx + 1` di tempat lain. Di model baru keduanya bernama menurut dirinya sendiri —
**`NOMOR_URUT_PEMULIHAN`** untuk nomor urut, dan nama yang menyebut nilai untuk nilainya.
**Membawa satu nama untuk dua arti adalah menanam cacat yang sama dengan tangan sendiri.**

Atribut `PEMULIHAN_LIMIT`:

| Nama | Asal | Tipe | Boleh kosong |
|---|---|---|---|
| `ID_PEMULIHAN_LIMIT` | *baru* | C1 | tidak |
| `ID_LAYER` | *baru* | R | tidak |
| `NOMOR_URUT_PEMULIHAN` | `Reinstatement_List[].ReinstatementValue` | C1 | tidak |
| `PERSEN_PEMULIHAN` | `Reinstatement_List[].ReinstatementPct` | P1 | tidak |
| `PERSEN_TAMBAHAN` | `Reinstatement_List[].AdditionalPct` | P1 | ya |
| `CATATAN` | `Reinstatement_List[].ReinstatementNote` | T | ya |

**Kunci alami** — `NOMOR_URUT_PEMULIHAN`, **Bentuk B** di dalam `ID_LAYER`. **Belum bernomor
invarian** — bergabung dengan kelima yang lain. **`SUMBU_REKONSILIASI`** — tidak punya baris
turunan; **INV-49 mengecualikannya dari partisi** karena pemulihan adalah besaran berulang.

**Tidak dibawa:** `ReinstatementAmount1`/`2` dan `AdditionalAmount1`/`2` — kembaran mata uang yang
peta telusur §2.3 sudah adili sebagai turunan, dan yang mata uang ketiganya menjadi nol (ADR-0035).
`LayerList` — sarang kelas pada dirinya sendiri, perancah clipboard.

**`Limit2`, `Deductible2`, `Currency2`, `AggregateLimit2`, `MDP2` dibuang** — pasangan kolom kembar
untuk mata uang kedua. Paket uang menggantikannya.

**`Limits[].TreatyType` dibuang**: fakta itu hidup di `Detail[]`, dan keberadaannya di dua tempat
adalah dua penulis untuk satu fakta (ADR-0041). Ia hanya ditulis oleh
`TreatyInMappingDataconvertProp` — aktivitas **proporsional** — sehingga pada layer
non-proporsional ia tidak pernah terisi.

Dua nama terakhir sengaja **tidak** memakai kata `reinstatement` telanjang untuk keduanya. Keduanya
besaran berbeda yang di sistem lama saling menimpa; istilah `reinstatement` tetap ada di glosarium,
yang dilarang adalah dua kolom yang namanya tidak membedakan isinya.

---

### 10.4 `DETAIL_PROPORSIONAL` — **12 atribut** (9 + 3 yang tidak pernah terlihat)

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_DETAIL_PROPORSIONAL` | *baru* | C1 | tidak | |
| `ID_LAYER` | *baru* | R | tidak | |
| `ID_KELOMPOK_TREATY` | `TreatyGroupID` | R | tidak | **kunci alami** di dalam **`LAYER`** — INV-06 |
| `JENIS_TREATY` | `TreatyType` | E | tidak | `QUOTA_SHARE` / `SURPLUS` |
| `PERSEN_QUOTA_SHARE` | `QSPct` | P1 | **ya** | terisi hanya bila `JENIS_TREATY = QUOTA_SHARE` |
| `JUMLAH_LINES_SURPLUS` | `Surplus` | C1 | **ya** | terisi hanya bila `JENIS_TREATY = SURPLUS` |
| `PERSEN_KOMISI_KOTOR` | `RIOGR` | P1 | ya | kepanjangan **belum diketahui** |
| `PERSEN_KOMISI_BERSIH` | `RIONR` | P1 | ya | kepanjangan **belum diketahui** |
| `PERSEN_CADANGAN_PREMI` | `PremiumReservePct` | P1 | ya | dari inventaris kelas Pega |
| `CADANGAN_PREMI` | `ReserveList[]` | U1 paket uang | ya | **BARU** — pasangan nilai dari persentase di atas; persentase tanpa nilainya adalah setengah fakta. Tingkat **BELUM DITENTUKAN** |
| `PERSEN_KAPASITAS_SURPLUS` | `IOOPct` | P1 | ya | **BARU** |
| `ID_SUSUNAN_RETRO` | `SpreadingTypeID` | R | ya | **BARU — penunjuk susunan baku yang menyemai penyebaran.** Syarat berdirinya `RINCIAN_PENYEBARAN` sebagai fakta terbukukan (INV-58); lihat §10.4a |

> **`RIOGR` dan `RIONR` adalah singkatan yang kepanjangannya tidak ada di mana pun** — persis yang
> CONTEXT.md §3.3b larang. Nama yang saya usulkan adalah **tebakan berdasarkan posisi** (`GR` ≈
> gross, `NR` ≈ net), **bukan pengetahuan**. Keduanya **tidak boleh masuk DDL** sebelum
> kepanjangannya dipastikan. Masuk daftar wawancara.

`PERSEN_QUOTA_SHARE` dan `JUMLAH_LINES_SURPLUS` **tidak** melanggar larangan "kolom yang artinya
bergantung kolom lain": arti masing-masing **tidak berubah**, ia hanya tidak berlaku pada cabang
yang lain. Yang dilarang adalah satu kolom yang **berganti arti**; ini dua kolom yang masing-masing
berarti satu hal, salah satunya kosong. Aturannya ditegakkan sebagai invarian di langkah 8: tepat
satu dari keduanya terisi, ditentukan `JENIS_TREATY`.

**Kunci alami dan lingkupnya** — `ID_KELOMPOK_TREATY`, **Bentuk B**: unik di dalam `ID_LAYER`,
**bukan** di dalam versi (INV-06). Menaikkannya ke versi melarang kelompok treaty yang sama muncul
di dua layer berbeda pada versi yang sama, dan itu hal yang biasa pada kontrak berlayer banyak.

**`SUMBU_REKONSILIASI`** — `ID_DETAIL_PROPORSIONAL` pada entitas anaknya (`POTONGAN`, `PENYEBARAN`
dan rantainya, serta keenam daftar uang di bawah).

#### 10.4a Tiga puluh delapan properti yang tidak pernah terlihat — diadili

**`DETAIL_PROPORSIONAL` diselesaikan atas sumber yang kurang.** Kelas
`Data-TreatyInLimitsDetail` punya 76 properti; **38 tidak pernah ada di pohon** (L-8) — jumlah
terbesar di seluruh modul. Lima jatuh ke GEL-3 dan enam belas ke agregat `Total…`/`Sum…`. **Tujuh
belas diadili**, dan satu di antaranya mengubah sesuatu yang sudah dinyatakan mengikat.

| Properti | Putusan | Dasar |
|---|---|---|
| **`SpreadingType`** | **ATRIBUT BARU — penunjuk susunan retro** | kotak di bawah |
| `CessionListSurplus` · `IOOLimitListSurplus` · `RetentionListSurplus` | **MELEBUR** ke `CessionList`, `IOOLimitList`, `RetentionList` | bentuknya persis pola berakhiran `P` (§10.0a): satu fakta di dua daftar, dibedakan **nilai kolom lain** (`JENIS_TREATY`). Yang membatalkannya: satu baris yang mengisi **kedua** varian dengan nilai berbeda — masuk **Uji X**, yang memang sudah dirancang untuk bentuk ini |
| `IOOPct` · `CurrencyIOOLimit` | **ATRIBUT BARU** — persentase dan mata uang `KAPASITAS_SURPLUS` | `KAPASITAS_SURPLUS` = retensi × jumlah *lines* (§3.4). Mata uangnya melebur ke paket uang; `IOOPct` berdiri sendiri |
| `ReserveList` | **ATRIBUT BARU — paket uang cadangan premi** | berpasangan dengan `PremiumReservePct` yang **sudah** ada di §10.4. Persentase tanpa nilainya adalah setengah fakta |
| `PremiumReservePct` | **sudah ada** | dibawa masuk §12.4 lewat inventaris kelas — bukti bahwa penambalan itu bekerja, untuk lima entitas yang kebetulan ditanyakan |
| `LowerBand` · `UpperBand` | **DITUNDA — ikut *profit commission*** | keduanya berdampingan dengan `LossRatio`, `ProfitCommision`, `ProfitME`, `ProfitYDCF` di kolom layar yang sama, dan kelas `Data-TreatyInProfitCommision` memuat `ME`, `ProfitComm`, `YDCF`. Itu **skala luncur komisi laba**, dan *profit commission* memang sudah DITUNDA (§6) |
| `LossRatio` · `Periode` | **GEL-3** | keluarga pencapaian |
| `RNMSpreadedList` · `RNMSpreadedListRI` | **turunan, tidak disimpan** | kembaran proporsional dari 10 daftar `…XOL` yang §12.3 sudah nyatakan turunan (ADR-0037) |
| `SpreadingTotalPct` · `Brokerage` | **turunan** | jumlah persen penyebaran, dan brokerage = persen × premi |
| `ParentID` | **dibuang** | salinan pengenal induk; di model baru ia **relasi**, sekeluarga dengan §12.2 |
| **`ReisuredParticipant`** *(salah eja ada di sumbernya)* | **dibuang** | **kolom layar mati**: nol aturan menulisnya, nol membacanya, dan kendalinya sendiri ber-`pyDisabled = true` sementara kolom tetangganya `QSPct` tidak. Sekalipun hidup, nama pihak sebagai teks bebas pada kontrak melewati ADR-0044 — sekeluarga dengan `PositionUsername` (§12.5). Pemeriksaan lengkapnya `SPEC-INVARIAN.md` §4.4a |

> ### `SpreadingType` / `SpreadingTypeID` ADALAH PENUNJUK YANG ADR-0036 TUNTUT
>
> ADR-0036 mewajibkan hasil yang dibekukan **membawa penunjuk ke masukan yang dipakai** —
> *"susunan mana, periode mana, kurs berapa dari sumber apa, bagian berapa, dan kapan dihitung"*.
> Kata pertamanya, **susunan mana**, tidak pernah punya kolom di model ini.
>
> Ia ada di sistem lama: `Limits.Detail.SpreadingTypeID` (8 rujukan) menunjuk **susunan retro baku
> mana** yang menyemai penyebaran, dan pasangannya `SpreadingType` adalah namanya. Di cabang
> non-proporsional pasangannya `Share.SpreadingTypeXOL` / `SpreadingTypeIDXOL`.
>
> **Tanpa kolom ini, `RINCIAN_PENYEBARAN` tidak dapat menjadi "fakta terbukukan yang membawa
> penunjuk asalnya"** — dan itu justru kedudukan yang ADR-0036 dan INV-58 berikan kepadanya
> (`TITIK-BUTA-POHON.md` §6.4). Pengecualian INV-58 **menuntut** penunjuk itu ada; selama ia tidak
> ada, `RINCIAN_PENYEBARAN` adalah turunan tersimpan tanpa alasan.
>
> Maka: **`ID_SUSUNAN_RETRO` menjadi atribut `DETAIL_PROPORSIONAL` dan `BAGIAN`**. Ia bukan
> kemewahan telusur — ia syarat berdirinya sebuah entitas.
>
> **Dan ia merujuk MASTER DI LUAR MODUL, bukan tabel acuan ketujuh.** Diperiksa 24 September 2026
> atas permintaan pemilik proses, yang menahan penggolongannya karena keenam tabel acuan yang ada
> berbentuk seragam dan kecil sementara susunan retro jelas tidak.
>
> `Activity/SetSpreadName.xml` mengambilnya dari `ReportDefinition/BrowseTreatyArrangement_Limit_RD`,
> dan laporan itu berkelas **`ASM-FW-GISFW-Int-PROPORTIONALARRG`** — kelas **`Int-`**, yaitu tabel
> di luar skema, sekeluarga dengan `Int-CURRENCY` dan `Int-TREATYGROUP`. Kolomnya:
>
> ```
> ID · Line · TreatyYear · TreatyYearID · TreatyGroupID · TreatyDescID
> ReinsTypeID · ReinsTypeName · ParentReinsTypeID · Pct · Rp · Usd
> ```
>
> | Yang terbaca | Akibatnya |
> |---|---|
> | ia **di luar** modul | kolomnya **rujukan ke master luar**, seperti cedant dan asal bisnis (ADR-0023, ADR-0041). **Tidak ada entitas yang ditambahkan** |
> | ia berkunci **tahun treaty** + kelompok treaty + deskripsi | susunannya **berubah antar tahun** — itu justru sebab ADR-0036 menuntut penunjuknya dibekukan, dan sebab *"perubahan pada susunan baku tidak menjalar"* |
> | `ParentReinsTypeID` adalah kolomnya | **susunan bertingkat jenis reasuransi TERBUKTI dari ekspor.** Dugaan di §10.6 tidak lagi dugaan, dan **Uji AE dicabut** — ia tidak perlu lagi |
> | ia **bukan** `Data-TreatyInShareReins` | yang terakhir itu retro **keluar** per **pihak**, GEL-2. Keduanya bersaudara tetapi bukan benda yang sama, dan batas itu dinyatakan di sini alih-alih ditemukan saat DDL |
>
> **Yang ditambahkan bukan entitas melainkan KEWAJIBAN**, dan bentuknya sama dengan P-49 yang
> tertahan karena pemilik tabel acuan kapasitas belum ditetapkan: **master susunan retro ada di
> mana, dan siapa pemiliknya.** Selama itu belum dijawab, `ID_SUSUNAN_RETRO` adalah kolom yang
> rujukannya tidak dapat ditegakkan — dan itu **penghalang**, bukan catatan.
>
> `PETA-TELUSUR-JSON.md` §5 **sudah menamai `SpreadingTypeID` sebagai salah satu dari dua properti
> yang hilang dari sapuannya**, dan menyebutnya dua kali. Yang tidak pernah terjadi: ia tidak pernah
> berjalan sampai menjadi atribut. **Lubang yang dicatat tetapi tidak ditindaklanjuti berperilaku
> persis seperti lubang yang tidak diketahui.**

`RetentionPct` dan `CessionPct` **tidak** ada di sini — keduanya turunan (§4.1). Masukannya adalah
`RetentionList[]`, yang memuat retensi **sebagai uang**.

---

### 10.5 `BAGIAN` — 8 atribut · **DIDAHULUKAN**

Cabang **non-proporsional saja** (§12.3). Sumber: pohon `TreatyIn.Share[]` + kelas
`Data-TreatyInShare` (47 properti, **6 tidak pernah terlihat**).

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_BAGIAN` | *baru* | C1 | tidak | |
| `ID_LAYER` | *baru* | R | tidak | **relasi**, bukan salinan — §12.2 |
| `CAKUPAN` | `Cover` | T | ya | tetap atribut bagian (§12.2) |
| `PERSEN_BAGIAN_NURE` | `Share[].RNMShare` | P1 | ya | terisi bila `BAGIAN_NURE_SERAGAM` mati; berbeda dari atribut senama pada versi (§10.2) yang berlaku seragam |
| `PREMI_BRUTO` | `GrossPremiumList[]` | U1 paket uang | tidak | tingkat **BELUM DITENTUKAN** — §10.23 |
| `PREMI_BRUTO_MINIMUM` | `GrossPremiumMinList[]` | U1 paket uang | ya | **BARU** — tidak pernah terlihat pohon |
| `ID_SUSUNAN_RETRO` | `SpreadingTypeIDXOL` | R | ya | penunjuk susunan baku yang menyemai penyebaran — **syarat INV-58**, §10.4a |
| `PERSEN_BAGIAN_DIPAKAI` | *baru* | P1 | ya | wajib bila paket uangnya bertingkat `BAGIAN_NURE` (INV-40) |

**Kunci alami dan lingkupnya** — **`ID_LAYER` sendiri**: *satu `BAGIAN` per `LAYER`*. **Bentuk B**,
dan bentuk yang tidak lazim: kunci alaminya **habis oleh kunci asing induknya**, sehingga
constraintnya `UNIQUE (ID_LAYER)` tanpa ruas tambahan. **Ia belum punya nomor invarian** — salah
satu dari lima, `SPEC-INVARIAN.md` §4.1.

**`SUMBU_REKONSILIASI`** — `ID_BAGIAN` pada `POTONGAN` dan `PENYEBARAN` di bawahnya.

**Yang tidak dibawa, dan sebabnya satu per satu:**

| Properti | Sebab |
|---|---|
| `Layer` · `LayerPart` · `LayerType` · `LayerPartType` | identitas layer, bukan atribut bagian — §12.2 |
| `ClassofBusinessList` · `TreatyGroupList` | **salinan** daftar milik `LAYER` (`Limits.TreatyGroupList`). Satu fakta satu penulis (ADR-0041). *Yang membatalkannya:* satu baris yang isinya berbeda dari daftar layernya — **Uji AF**, baru |
| `GrossPremium` · `NetPremium` *(tidak pernah terlihat)* | kembaran **skalar** dari daftar yang sudah ada; `NetPremium` juga turunan (INV-52) |
| `NusareLimit` *(tidak pernah terlihat)* | turunan — limit dikali bagian NuRe; kembaran skalar `RnmLimitList` |
| `IsLayerActive` *(tidak pernah terlihat)* | salinan keadaan layer, sekeluarga dengan baris pertama |
| `SpreadingTypeXOLRetro` · `…RetroID` *(tidak pernah terlihat)* | susunan retro **dari retro** — arah keluar, **GEL-2** |
| 10 daftar `RNMSpreadedList…XOL`, `BrokerageList`, `NetPremiumList`, `RnmLimitList`, `RnmGrossPremiDisplay`, `RnmLimitListDisplay`, `SpreadingTotalPctXOL` | turunan, tidak disimpan — §12.3, ADR-0037 |

> **`SharePct` BUKAN atribut `BAGIAN`.** §3.6 mendaftarnya sebagai `Share[].SharePct`; pohon tidak
> memuatnya di bawah `Share`, dan satu-satunya `SharePct` yang hidup ada di
> `ShareFacultativeReinsurers[]` (§3.7, GEL-2) dan di `BreakDownSprdList[]`. **§3.6 dikoreksi.**

---

### 10.6 `PENYEBARAN` — 5 atribut

Satu bentuk untuk kedua cabang; induknya **`BAGIAN` atau `DETAIL_PROPORSIONAL`**.

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_PENYEBARAN` | *baru* | C1 | tidak | |
| `ID_INDUK_PENYEBARAN` | *baru* | R | tidak | dua induk; bentuk fisiknya sesi DDL, syarat mengikat INV-63 |
| `ID_JENIS_REASURANSI` | `ReinsTypeID` | R | tidak | **kunci alami** |
| `ID_JENIS_REASURANSI_INDUK` | `ParentReinsTypeID` | R | ya | jenis reasuransi **bersusun**; rujukan ke baris lain di tabel acuan yang sama |
| `PERSEN_PENYEBARAN` | `Pct` | P1 | tidak | INV-50 menjumlahkannya |

**Kunci alami dan lingkupnya** — `ID_JENIS_REASURANSI`, **Bentuk B**: unik di dalam
`ID_INDUK_PENYEBARAN` (INV-16). **Kunci padanan seam-nya majemuk** dan harus menyebut **induk mana**
dari kedua pelekatan — `SEAM-ADJUSTMENT.md` §3.

**`SUMBU_REKONSILIASI`** — `ID_PENYEBARAN` pada `RINCIAN_PENYEBARAN`.

**Tidak dibawa:** `ReinsTypeName` (salinan tabel acuan, INV-59) · `TreatyGroupID` (salinan induk) ·
`Rp`, `Usd` (kembaran mata uang, §12.4) · `GrossPremiumList`, `GrossPremiumMinList`,
`NetPremiumList`, `DeductionTotalList`, `RnmLimitList`, dan 10 `RNMSpreadedList…` (turunan).

> **`ParentReinsTypeID` dibaca sebagai susunan bertingkat jenis reasuransi**, bukan penunjuk ke
> baris penyebaran lain — ia skalar yang duduk berdampingan dengan `ReinsTypeID` di baris yang sama.
> **Tidak lagi dugaan, dan Uji AE dicabut 24 September 2026.** `ParentReinsTypeID` adalah **kolom
> tabel master** `Int-PROPORTIONALARRG`, berdampingan dengan `ReinsTypeID` dan `Pct` — terbaca dari
> `ReportDefinition/BrowseTreatyArrangement_Limit_RD`. Ia susunan bertingkat **jenis reasuransi**,
> dan itu terbukti dari ekspor, bukan dari uji data. Lihat §10.4a.

---

### 10.7 `RINCIAN_PENYEBARAN` — 4 atribut · **DIDAHULUKAN**

**Sumber: `BreakDownSprdList[]` (proporsional) dan `BreakDownSprdListXOL[]` (non-proporsional)** —
keduanya **tidak pernah ada di pohon** (L-8). Ini entitas yang membongkar titik butanya.

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_RINCIAN_PENYEBARAN` | *baru* | C1 | tidak | |
| `ID_PENYEBARAN` | *baru* | R | tidak | |
| `ID_JENIS_REASURANSI` | `ReinsID` (prop) / `ReinsTypeID` (XOL) | R | tidak | **kunci alami**. Ruas proporsional **bernama pihak tetapi berisi jenis** — §4.4a |
| `PERSEN_RINCIAN` | `SharePct` | P1 | tidak | INV-50 menjumlahkannya per induk |

**Kunci alami dan lingkupnya** — `ID_JENIS_REASURANSI`, **Bentuk B**: unik di dalam
`ID_PENYEBARAN`. **Belum punya nomor invarian** — salah satu dari lima.

**`SUMBU_REKONSILIASI`** — **`ID_JENIS_REASURANSI`**, dikelompokkan dalam `ID_PENYEBARAN`.
**Bukan pihak.** Inilah kolom yang INV-47 dan INV-50 tunjuk, dan yang sesi DDL pakai di `GROUP BY`.

**Kedudukannya:** **fakta terbukukan yang membawa penunjuk asalnya** — ADR-0036 menyebut penyebaran
dengan nama, INV-58 menyediakan pengecualiannya, dan penunjuk asalnya adalah `ID_SUSUNAN_RETRO`
pada induknya (§10.4a, §10.5).

> ### PENUNJUKNYA BELUM TENTU MENJELASKAN SELURUH ISINYA — jangan dibaca sebagai jaminan
>
> `ID_SUSUNAN_RETRO` bermakna **hanya bila seluruh isi hasilnya datang dari yang ditunjuk**. Satu
> baris diketahui **tidak** datang dari sana: `FetchQSfromMasterXOL.xml` menulis sendiri
> `ReinsTypeID = "10007"`, `ReinsTypeName = "ORS"`, `Pct = "15.00"` — bukan mengambilnya dari
> `Int-PROPORTIONALARRG` (`TETAPAN-DI-KODE.md` §1).
>
> Maka sebuah baris penyebaran beku dapat menunjuk ke susunan retro yang **tidak menjelaskan seluruh
> isinya**. Orang yang lima tahun lagi menelusuri *"dari mana angka ini"* akan sampai ke master itu,
> **tidak menemukan ORS 15%, dan menyimpulkan datanya rusak.**
>
> | Yang dijawab **Uji AH** | Akibatnya |
> |---|---|
> | **ada** baris ORS di master | kodenya **menduplikasi** master; penunjuknya menjelaskan seluruhnya. Yang diperlukan hanya memastikan barisnya ada sebelum migrasi — **langkah cut-over**, bukan perubahan model |
> | **tidak ada** | *"ORS selalu 15% pada XOL"* adalah aturan bisnis yang **tidak ada di mana pun selain berkas itu**, dan ia **hilang dalam migrasi** kecuali ditangkap sekarang. Model bertambah **satu kolom** yang menyatakan *"baris ini tidak berasal dari susunan"* |
>
> **Sampai Uji AH terjawab, kalimat "membawa penunjuk asalnya" di atas adalah niat, bukan jaminan.** **Golongan kemampuannya TERTAHAN**, bukan bentuknya —
`5-tiket/DAFTAR-PEKERJAAN.md` §2.3.

**Tidak dibawa:** `ReinsName` / `ReinsTypeName` (salinan tabel acuan) · `Amount` — **turunan**,
dihitung *nilai × bagian NuRe × persen master* oleh `SetSpreadName.xml`, dan ADR-0037 melarang
menyimpannya. Yang dibekukan adalah **persennya**, bukan hasil kalinya.

---

### 10.8 `NILAI_PENYEBARAN` — kedudukannya **DIANGKAT**, bukan diputuskan sendiri · DIDAHULUKAN

`STRUKTUR-DATA.md` §1.4 menyatakannya entitas: *"nilai sebuah rincian penyebaran pada satu mata
uang"*, kunci alami **mata uang**, induk `RINCIAN_PENYEBARAN`. Sesudah `BreakDownSprdList` terbaca,
**bukti untuk bentuk itu tidak ada**:

| Yang diperiksa | Hasil |
|---|---|
| bentuk `Currency` dan `Amount` di `BreakDownSprdList` | **skalar**, bukan daftar — **satu** mata uang per baris rincian |
| dari mana `Currency` diisi | `.Currency` milik `Detail[]` induknya — satu nilai, bukan gelung |
| adakah daftar nilai per mata uang di bawah rincian | **tidak ada**, di kelas maupun di ekspor |

> Maka yang terbaca: **`NILAI_PENYEBARAN` adalah paket uang milik `RINCIAN_PENYEBARAN`, bukan
> entitas anak.** Satu baris per rincian, bukan satu baris per mata uang.

**Saya tidak membuangnya sendiri.** Daftar entitas mengikat (wewenang butir 4), dan ini
**mengurangi** hitungan 27 menjadi 26 — perubahan yang berlawanan arah dengan seluruh temuan hari
ini, dan justru karena itu layak diperiksa orang lain.

| Pilihan | Ongkos bila salah |
|---|---|
| lebur menjadi paket uang pada `RINCIAN_PENYEBARAN` | bila ternyata sebuah rincian dapat bernilai dalam beberapa mata uang, ia **perubahan skema** — mahal |
| **pertahankan sebagai entitas** | satu tabel yang selalu berisi satu baris per induk — murah, dan ia merosot sendiri |

Di sini ongkosnya **berlawanan arah** dengan kasus pemulihan limit (§10.3b): yang murah adalah
**mempertahankan**. **Usul saya: pertahankan**, dengan catatan jujur bahwa buktinya satu baris per
induk, dan kunci alaminya **mata uang** tetap benar meski himpunannya kini beranggota satu.

Bila dipertahankan, atributnya:

| Nama | Asal | Tipe | Boleh kosong |
|---|---|---|---|
| `ID_NILAI_PENYEBARAN` | *baru* | C1 | tidak |
| `ID_RINCIAN_PENYEBARAN` | *baru* | R | tidak |
| `NILAI` | `Amount` + `Currency` | U1 paket uang | tidak |

**Kunci alami** — kode mata uang, **Bentuk B** di dalam `ID_RINCIAN_PENYEBARAN`; **belum punya
nomor invarian**. **`SUMBU_REKONSILIASI`** — tidak punya baris turunan.

---

### 10.9 `POTONGAN` — 5 atribut

**Atributnya sudah ditetapkan §14.2 dan tidak diulang di sini** — mengulangnya membuat dua tempat,
satu akan basi (`SPEC-INVARIAN.md` §4.3). Yang ditambahkan hanya kedua kewajiban §5.0b:

**Kunci alami dan lingkupnya** — `ID_JENIS_POTONGAN`, **Bentuk B**: unik di dalam
`ID_INDUK_POTONGAN` (INV-15). Karena induknya **dua** (§14.1), kunci padanan seam-nya harus
menyebut **induk mana** — `SEAM-ADJUSTMENT.md` §3.

**`SUMBU_REKONSILIASI`** — tidak punya baris turunan. Yang menjumlahkan potongan adalah **INV-52**
(premi bruto dikurangi seluruh potongan sama dengan premi bersih), dan ia menjumlahkan **ke arah
induk**, bukan ke arah anak.

**Satu properti tidak pernah terlihat:** `DeductionPctCalculate` — **tetap dibuang**, bendera mode
hitung, bukan fakta potongan (§12.4). Titik buta tidak mengubah putusan yang sudah punya alasan.

---

### 10.10 `MATA_UANG_KONTRAK` — 6 atribut

Kelas `Data-TreatyInCurrencyList` (8 properti, **3 tidak pernah terlihat**).

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_MATA_UANG_KONTRAK` | *baru* | C1 | tidak | |
| `ID_VERSI_KONTRAK` | *baru* | R | tidak | |
| `KODE_MATA_UANG` | `CurrencyID` | R | tidak | **kunci alami**; `Currency` (nama) tidak disimpan — INV-59 |
| `KURS` | `Conversion` | K1 | tidak | INV-38: selalu lebih besar dari nol |
| `TANGGAL_MULAI_BERLAKU` | `PeriodStart` | D | tidak | batas inklusif, ADR-0022 |
| `TANGGAL_AKHIR_BERLAKU` | `PeriodEnd` | D | tidak | batas inklusif |

**Kunci alami dan lingkupnya** — `KODE_MATA_UANG`, **Bentuk A** di dalam `ID_VERSI_KONTRAK`
(INV-07). **`SUMBU_REKONSILIASI`** — tidak punya baris turunan.

**Diadili:** `AchievementPctGross`, `LossRatioGross` → **GEL-3**. `Parameter` → **dibuang**, nama
yang tidak menyatakan apa pun dan tidak ditulis satu aturan pun; membawanya berarti membawa kolom
yang artinya harus ditanyakan setiap kali dibaca.

> **Kurs dicatat sebagai nilai pada versi** (INV-43), bukan dibaca ulang saat dibutuhkan. Itu yang
> membuat entitas ini ada: tanpa ia, konversi memakai kurs hari ini atas peristiwa tahun lalu —
> cacat ketiga yang melahirkan ADR-0036.

---

### 10.11 `RETENSI_CEDANT` — 6 atribut

Kelas `Data-TreatyInRetention` (7 properti, **0 tidak terlihat** — entitas ini utuh di pohon).

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_RETENSI_CEDANT` | *baru* | C1 | tidak | |
| `ID_VERSI_KONTRAK` | *baru* | R | tidak | |
| `ID_KELOMPOK_TREATY` | `TreatyGroupID` | R | tidak | bagian **kunci alami**; `TreatyGroup` (nama) tidak disimpan |
| `NILAI_RETENSI` | `Amount` + `Currency`/`CurrencyID` | U1 paket uang | tidak | bagian **kunci alami** lewat kode mata uangnya. Tingkat **BELUM DITENTUKAN** — §10.23 |
| `CATATAN` | `Note` | T | ya | **BARU** — ada di kelas, tidak ada di §3.5 |
| `PERSEN_BAGIAN_DIPAKAI` | *baru* | P1 | ya | bila paket uangnya bertingkat `BAGIAN_NURE` (INV-40) |

**Kunci alami dan lingkupnya** — `ID_KELOMPOK_TREATY` + kode mata uang, **Bentuk A** di dalam
`ID_VERSI_KONTRAK` (INV-08). **`SUMBU_REKONSILIASI`** — tidak punya baris turunan.

> `Note` menunjukkan sesuatu yang berlaku lebih luas: **penyisiran jalur menangkap yang ditulis
> aturan, dan ruas yang hanya diketik orang di layar dapat lolos.** Itu sebab kedua sebuah properti
> hilang, berbeda dari L-8, dan ia sudah tercatat di §3 sebagai `TAK-TERTULIS`.

---

### 10.12 `EGNPI` — 8 atribut

Kelas `Data-TreatyInEGNPI` (12 properti, **0 tidak terlihat**). **Lima ada di kelas tetapi tidak di
§3.5** — daftar §3.5 menyebut tujuh.

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_EGNPI` | *baru* | C1 | tidak | |
| `ID_VERSI_KONTRAK` | *baru* | R | tidak | |
| `ID_KELOMPOK_TREATY` | `TreatyGroupID` | R | tidak | bagian **kunci alami** |
| `ID_KELAS_BISNIS` | `ClassOfBusiness` | R | ya | **BARU** — tidak ada di §3.5 |
| `NILAI_EGNPI` | `Amount` + `Currency`/`CurrencyID` | U1 paket uang | tidak | tingkat **BELUM DITENTUKAN**, dan ini **satu dari lima §7 yang naik ke eskalasi** |
| `TANGGAL_BERLAKU` | `AsDate` | D | ya | |
| `PROPORSI` | `Proportion` | P1 | ya | **BARU** — belum terpakai di rumus mana pun; dibawa karena membuangnya menghilangkan fakta yang tidak dapat dipulihkan |
| `CATATAN` | `Note` | T | ya | **BARU** |

**Kunci alami dan lingkupnya** — `ID_KELOMPOK_TREATY` + kode mata uang, **Bentuk A** di dalam
`ID_VERSI_KONTRAK` (INV-09). **`SUMBU_REKONSILIASI`** — tidak punya baris turunan.

**Tidak dibawa:** `AmountIDR` — **turunan** dari nilai dan kurs, dan paket uang sudah membawa
`NILAI_IDR` sebagai bagian dirinya. `AddendumStatus` — bendera warisan alur addendum, digantikan
`JENIS_ADDENDUM` pada versi (§10.2).

> **Tingkat pencatatan `EGNPI` adalah butir eskalasi, bukan pertanyaan model** (§7): salah tingkat
> di sini berarti **penyajian ulang historis**, karena rumus pencapaian mengandaikan 100%.

---

### 10.13 `PORTOFOLIO` — 5 atribut

Kelas `Data-TreatyInPortfolio` (3 properti, **0 tidak terlihat**).

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_PORTOFOLIO` | *baru* | C1 | tidak | |
| `ID_VERSI_KONTRAK` | *baru* | R | tidak | |
| `ARAH_PORTOFOLIO` | `Type` | E | tidak | masuk atau keluar; bagian **kunci alami** |
| `JENIS_PORTOFOLIO` | `TypePortfolio` | E | tidak | bagian **kunci alami** |
| `KETERANGAN` | `Description` | T | ya | |

**Kunci alami dan lingkupnya** — `ARAH_PORTOFOLIO` + `JENIS_PORTOFOLIO`, **Bentuk A** di dalam
`ID_VERSI_KONTRAK`. **Ia belum punya nomor invarian** — salah satu dari lima, dan satu dari dua
yang berbentuk A. **`SUMBU_REKONSILIASI`** — tidak punya baris turunan.

---

### 10.14 `PERIODE_PELAPORAN` — 7 atribut

Kelas `Data-TreatyInAccountReport` (6 properti, **1 tidak terlihat**).

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_PERIODE_PELAPORAN` | *baru* | C1 | tidak | |
| `ID_VERSI_KONTRAK` | *baru* | R | tidak | |
| `PERIODE` | `Period` | T | tidak | **kunci alami** |
| `TANGGAL_AWAL` | `InitialDate` | D | tidak | |
| `BATAS_PENYERAHAN` | `SubmissionDue` | D | tidak | |
| `BATAS_KONFIRMASI` | `ConfirmationDue` | D | tidak | |
| `BATAS_PELUNASAN` | `SettlementDue` | D | tidak | |

**Kunci alami dan lingkupnya** — `PERIODE`, **Bentuk A** di dalam `ID_VERSI_KONTRAK` (INV-10).
**`SUMBU_REKONSILIASI`** — tidak punya baris turunan.

**Diadili:** `AutoCalculate` → **dibuang**, bendera mode hitung yang menyatakan **bagaimana** ketiga
tanggal itu diperoleh, bukan **apa** tanggalnya. Sekeluarga dengan `DeductionPctCalculate` (§12.4).
*Yang membatalkannya:* bila bisnis membedakan periode yang tanggalnya dihitung dari yang diketik
tangan sebagai fakta yang berbeda — **pertanyaan wawancara**, satu kalimat.

> **Tidak ada bendera terlambat**, dan itu keputusan pemilik proses (OPSI 2): keterlambatan adalah
> turunan yang dihitung saat dibaca — ADR-0036.

---

### 10.15 `PERIODE_AKUMULASI` — 6 atribut

Kelas `Data-TreatyInAccumulation` (4 properti; `SubDays` dan `SubDueDate` **sudah** ditambal §12.4).

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_PERIODE_AKUMULASI` | *baru* | C1 | tidak | |
| `ID_VERSI_KONTRAK` | *baru* | R | tidak | |
| `PERIODE` | `Period` | T | tidak | **kunci alami** |
| `TANGGAL_LAPOR` | `ReportDate` | D | tidak | |
| `HARI_BATAS_PENYERAHAN` | `SubDays` | C1 | ya | dari inventaris kelas, §12.4 |
| `BATAS_PENYERAHAN` | `SubDueDate` | D | ya | dari inventaris kelas, §12.4 |

**Kunci alami dan lingkupnya** — `PERIODE`, **Bentuk A** di dalam `ID_VERSI_KONTRAK` (INV-11).
**`SUMBU_REKONSILIASI`** — tidak punya baris turunan.

> `HARI_BATAS_PENYERAHAN` dan `BATAS_PENYERAHAN` adalah **hari** dan **tanggal** — satu turunan dari
> yang lain. Keduanya dibawa karena yang tersimpan di sistem lama keduanya, dan mana yang masukan
> **belum terbaca**. Sesi DDL tidak boleh menyimpan keduanya tanpa menetapkan penulisnya
> (ADR-0041) — dicatat sebagai keputusan yang menunggu, bukan sebagai dua kolom biasa.

---

### 10.16 `TERMIN` — 8 atribut, **dan tingkat yang terbalik**

Kelas `Data-TreatyInInstallment` (12 properti).

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_TERMIN` | *baru* | C1 | tidak | |
| `ID_VERSI_KONTRAK` | *baru* | R | tidak | |
| `NOMOR_TERMIN` | `InstallmentList[].Installment` | C1 | tidak | bagian **kunci alami**, bersama kode mata uang — §10.16a |
| `PERSEN_TERMIN` | `InstallmentPct` | P1 | tidak | |
| `NILAI_TERMIN` | `Amount` + `Currency`/`CurrencyID` | U1 paket uang | ya | tingkat **BELUM DITENTUKAN** — §10.23 |
| `TANGGAL_JATUH_TEMPO` | `DueDate` | D | tidak | |
| `TANGGAL_BAYAR` | `PaymentDate` | D | ya | kosong selama belum dibayar |
| `WPC` | `WPC` | T | ya | **singkatan yang kepanjangannya tidak ada di korpus** — lihat di bawah |

**Kunci alami dan lingkupnya** — `NOMOR_TERMIN`, **Bentuk A** di dalam `ID_VERSI_KONTRAK` (INV-12).
**`SUMBU_REKONSILIASI`** — **bergantung keputusan di bawah**: bila `InstallmentList` menjadi entitas,
sumbunya `ID_TERMIN`; bila tidak, `TERMIN` tidak punya baris turunan.

**Tidak dibawa:** `PctTotal`, `AmountTotal` — turunan (§12.4).

> **`WPC` masuk daftar wawancara bersama `RSMD`.** Teknik §2.10 `CONTEXT.md` dijalankan: ia tidak
> punya saudara bernama panjang di kelas mana pun, dan tidak muncul di modul lain. Ia **tidak boleh
> masuk DDL dengan nama itu** sebelum kepanjangannya dipastikan — aturan yang sama yang menahan
> `RIOGR` dan `RIONR` (§10.4).

#### 10.16a `Installment[].InstallmentList[]` — bukan anak, melainkan TERMIN itu sendiri

`peta-nama-tabel-treatyin.tsv` sudah memuat tabel datar `T_TREATY_INSTALLMENT_ITEM`
(`RINCIAN_ANGSURAN`) dengan induk `T_TREATY_INSTALLMENT` dan ruas `DueDate`, `InstallmentPct`,
`PaymentDate`, `WPC` — **dan buktinya tercatat sebagai "bentuknya hanya terbaca lewat salinan
`OLDDATA`/`ValueDifference`"**. Kelas `Data-TreatyInInstallment` memang memuat `InstallmentList`.

**`STRUKTUR-DATA.md` tidak memuat entitas untuk itu.** Peta tabel datar punya, daftar entitas tidak
— dan keduanya tidak pernah dicocokkan ke dua arah, persis bentuk kegagalan yang melahirkan
kewajiban §5.0b.

> ### DIPERIKSA 24 September 2026 — BUKAN entitas baru. **TINGKATNYA YANG TERBALIK.**
>
> Pemilik proses menahan usul entitas baru dan meminta dibuktikan dulu bahwa ia bukan `TERMIN`
> dengan nama lain — uji yang sama yang membongkar `FacShare` sebagai salinan. Ujinya dijalankan,
> dan jawabannya bukan salah satu dari dua kemungkinan yang saya ajukan.

`Activity/TreatyInSetValueInstallment.xml` — **27 langkah, seluruhnya hidup** — menunjukkan:

```
langkah  6  "Set currency list (also total's currency)"
            TreatyIn.Installment(<APPEND>).Currency = .Currency
langkah  8  "100 / TreatyIn.InstallmentNo"
langkah  9  "Set installment list"
langkah 14  "set installment no & due date"
            .InstallmentList(<APPEND>).Installment  = Local.InstallmentNo
            .InstallmentList(<LAST>).DueDate        = @CurrentDate(…)
            .InstallmentList(<LAST>).InstallmentPct = Local.Percentage
            .InstallmentList(<LAST>).Amount         = pct/100 × value
```

| Tingkat | Apa sebenarnya |
|---|---|
| `TreatyIn.Installment[]` — yang saya sebut `TERMIN` | **satu baris per MATA UANG.** Ruas miliknya sendiri hanya `Currency`, `CurrencyID`, `ID`, dan dua total turunan |
| `Installment[].InstallmentList[]` | **termin yang sebenarnya** — nomor 1..N, persentase, jatuh tempo, nilai, tanggal bayar, `WPC` |

**Maka §10.16 di atas mencampur dua tingkat**, dan koreksinya **mengurangi** entitas alih-alih
menambah: baris pengelompokan per mata uang **larut**, karena di model ini setiap besaran uang
membawa mata uangnya **di dalam paket uangnya sendiri**. Baris luar itu ada di Pega hanya karena
daftar clipboard tidak dapat memikul kunci majemuk.

> **Tidak ada entitas baru. `TERMIN` tetap satu, tetapi isinya diambil dari `InstallmentList`.**

**Dan satu invarian ikut salah lingkup.** INV-12 berbunyi *"nomor termin unik di dalam satu versi"*.
Bila sebuah kontrak bermata uang dua, ada **dua** rangkaian termin bernomor 1..N pada versi yang
sama — jadi nomor termin **berulang** di dalam versi. **INV-12 berbentuk B, bukan A**, dan kunci
alaminya **kode mata uang + nomor termin**. Diperbaiki di `SPEC-INVARIAN.md` §4.1.

**Ini invarian ketiga yang lingkupnya keliru** — sesudah INV-47 dan INV-50 — dan ketiganya keliru
dengan sebab yang sama: **lingkup ditulis dari bentuk yang terlihat, bukan dari penulisnya.**

---

### 10.17 `SKALA_KOASURANSI` — 4 atribut

Kelas `Data-TreatyInCoInScale` (2 properti, keduanya **sudah** ditambal §12.4; ia memang **daftar**,
bukan skalar).

| Nama | Asal | Tipe | Boleh kosong |
|---|---|---|---|
| `ID_SKALA_KOASURANSI` | *baru* | C1 | tidak |
| `ID_VERSI_KONTRAK` | *baru* | R | tidak |
| `PERSEN_LIMIT` | `PctLimit` | P1 | tidak |
| `PERSEN_BAGIAN` | `CoInShare` | P1 | tidak |

**Kunci alami dan lingkupnya** — `PERSEN_LIMIT`, **Bentuk A** di dalam `ID_VERSI_KONTRAK` (INV-13).
**`SUMBU_REKONSILIASI`** — tidak punya baris turunan.

---

### 10.18 `BATAS_PER_BAHAYA` — 5 atribut

Lahir §10.0b dari delapan kolom bernama bahaya di **kepala**, bukan dari kelas tersendiri.

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_BATAS_PER_BAHAYA` | *baru* | C1 | tidak | |
| `ID_VERSI_KONTRAK` | *baru* | R | tidak | |
| `ID_BAHAYA` | nama kolom lama | R | tidak | **kunci alami**; tabel acuan, ADR-0038 |
| `NILAI_BATAS` | `Earthquake` / `FloodJab` / `FloodNation` / `RSMDLimit` + mata uangnya | U1 paket uang | tidak | tingkat **BELUM DITENTUKAN** — §10.23 |
| `PERSEN_BAGIAN_DIPAKAI` | *baru* | P1 | ya | bila bertingkat `BAGIAN_NURE` (INV-40) |

**Kunci alami dan lingkupnya** — `ID_BAHAYA`, **Bentuk A** di dalam `ID_VERSI_KONTRAK` (INV-14).
**`SUMBU_REKONSILIASI`** — tidak punya baris turunan.

**Isi awal tabel acuan `BAHAYA`: empat** — gempa, banjir Jabodetabek, banjir nasional, dan `RSMD`.
**Empat adalah awal, bukan batas**, dan justru itu sebabnya bentuk ini dipilih: kolom kesembilan
akan diminta begitu ada bahaya baru, dan bahaya baru pasti terjadi.

> **`RSMD` tetap tebakan.** Korpus sudah disapu habis — 44 kemunculan, nol kepanjangan
> (`CONTEXT.md` §2.10). Yang pasti hanya golongannya: ia **nama bahaya**, sekeluarga dengan gempa
> dan banjir, dan itulah yang mengukuhkan tabel acuan ini.

---

### 10.19 `DOKUMEN_KONTRAK` — 5 atribut

**Dirujuk, tidak dimiliki** (ADR-0027). Sumbernya lampiran `M_ATTACHMENTTREATY_2`, bukan kelas
clipboard.

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_DOKUMEN_KONTRAK` | *baru* | C1 | tidak | |
| `ID_VERSI_KONTRAK` | *baru* | R | tidak | |
| `ID_DOKUMEN` | lampiran | R | tidak | **kunci alami**; rujukan ke luar skema |
| `JENIS_DOKUMEN` | lampiran | E | ya | |
| `TANGGAL_LAMPIR` | lampiran | D | tidak | |

**Kunci alami dan lingkupnya** — `ID_DOKUMEN`, **Bentuk A** di dalam `ID_VERSI_KONTRAK`. **Belum
punya nomor invarian** — salah satu dari lima, dan yang kedua berbentuk A.
**`SUMBU_REKONSILIASI`** — tidak punya baris turunan.

> **Isi dokumennya tidak disalin**, hanya pengenalnya. Tabel lampiran adalah milik modul lain, dan
> menyalin namanya berarti dua penulis untuk satu fakta (ADR-0041).

---

### 10.19a `DOKUMEN_ADDENDUM` — 3 atribut · **ENTITAS BARU, ditambahkan 24 Sep 2026 (diff `D-2`)**

Dari **`GRL-19`**, ronde D grilling Adjustment. **Tidak punya padanan di sistem lama.**

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_DOKUMEN_ADDENDUM` | *baru* | C1 | tidak | |
| `NOMOR_DOKUMEN` | **tidak ada di sistem lama** | T | tidak | **kunci alami**, unik **global** — `INV-71`. Nomor yang beredar di luar sistem: ditulis di kertas dan disebut orang |
| `TANGGAL_BERLAKU` | *baru* | D | **ya** | `KTV-2` — bila kosong, berlaku mengikuti versi yang dipayunginya. **Dicabut bila `DB-16b` dibantah** |

**Kunci alami dan lingkupnya** — `NOMOR_DOKUMEN`, **bukan Bentuk A maupun B**: ia tidak punya induk,
sehingga lingkupnya **seluruh tabel** (`INV-71`).

**`SUMBU_REKONSILIASI`** — tidak punya baris turunan.

Dan pada **`VERSI_KONTRAK`** ditambahkan satu atribut:

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_DOKUMEN_ADDENDUM` | *baru* | R | **ya** | satu dokumen memayungi banyak versi, **lintas kontrak** (`DB-3`, `DB-4` dibantah). **Persetujuan tetap per versi** |

#### 10.19b Arti "kosong untuk seluruh baris warisan" — DIBACA BEGINI, bukan begitu

> **Kalimat ini punya dua bacaan yang menghasilkan DDL berbeda, dan hanya satu yang benar.**

| Bacaan | Akibat pada `NOMOR_DOKUMEN NOT NULL` |
|---|---|
| **✅ BENAR** — tabel `DOKUMEN_ADDENDUM` **tidak punya satu pun baris warisan**; yang kosong adalah **`VERSI_KONTRAK.ID_DOKUMEN_ADDENDUM`**, bernilai `NULL` | **aman** — tidak ada baris dokumen, jadi tidak ada `NOMOR_DOKUMEN` kosong |
| ❌ salah — barisnya **ada** tetapi nomornya kosong | **pecah saat migrasi**, dan penambalnya akan mencabut `NOT NULL` |

> **Kenapa ini ditulis eksplisit.** Mencabut `NOT NULL` dari `NOMOR_DOKUMEN` **menghapus kunci
> alaminya** — `INV-71` tidak lagi dapat ditegakkan, dan dokumen kehilangan satu-satunya pengenal
> yang dipakai manusia. Orang berikutnya yang membaca *"kosong untuk seluruh baris warisan"* tanpa
> §10.19b ini akan melakukannya, dan kerusakannya tidak akan terlihat sampai ada dua dokumen
> bernomor sama.

**Migrasi tidak menyisipkan satu pun baris `DOKUMEN_ADDENDUM`.** Pengisiannya **pekerjaan orang dari
arsip kertas**, dan itu masuk **rencana peralihan**, bukan spesifikasi.

### 10.20 `CATATAN_PERSETUJUAN` — 6 atribut

Kelas `Data-SuggestList` (5 properti, **1 tidak terlihat**).

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_CATATAN_PERSETUJUAN` | *baru* | C1 | tidak | |
| `ID_VERSI_KONTRAK` | *baru* | R | tidak | menggantung pada **versi**, bukan kontrak |
| `WAKTU_KEPUTUSAN` | `Date` | D | tidak | |
| `NAMA_PEMUTUS` | `OperatorName` | T | tidak | **dibekukan sebagai teks dengan sengaja** — ia fakta historis siapa memutuskan apa dan kapan (§12.5, ADR-0045) |
| `DISETUJUI` | `IsApproved` | E | tidak | |
| `ALASAN` | `Suggest` | T | ya | |

**Kunci alami — TIDAK ADA, dan itu keputusan.** Dua keputusan pada versi yang sama oleh orang yang
sama pada hari yang sama adalah keadaan yang sah; yang membedakan barisnya **urutan waktu**, bukan
sebuah nilai. Akibatnya: **baris entitas ini tidak dapat dipadankan antar versi**, dan itu benar —
ia catatan peristiwa, bukan ketentuan kontrak. **`SUMBU_REKONSILIASI`** — tidak punya baris turunan.

**Diadili:** `ConvertDate` → **dibuang**, cap waktu proses konversi, bukan fakta keputusan.

> ### 10.20a Bentuknya diuji 24 Sep 2026 dan **BERTAHAN** — `DISETUJUI` tetap tidak boleh kosong
>
> Sapuan penulis yang lengkap menemukan bahwa `CommentList` di sistem lama memuat **dua jenis
> baris**: keputusan (ber-`IsApproved` terisi) dan **peristiwa** (ber-`IsApproved` kosong —
> *"Had Created Internal Edit"*, *"Copied from ID …"*, *"Create Revision"*).
>
> Usul sementara saya adalah **nilai enumerasi ketiga `PERISTIWA`**. Usul itu **dicabut**: ia akan
> membuat satu tabel yang separuh kolomnya kosong pada separuh barisnya, dan mencemari **INV-28**,
> yang menghitung satu baris per perpindahan. Baris peristiwa **keluar** ke entitas tersendiri —
> §10.24. **Entitas ini tidak berubah sama sekali.**
>
> Pemeriksaan lengkapnya, beserta satu **tabrakan yang dilaporkan dan tidak diputuskan** (baris
> pengajuan juga membawa `IsApproved = "Accept"`): `4-erd-dan-tabel-datar/PERIKSA-BUTIR-CONTOH.md`
> §7.10-Z.

---

### 10.21 `JEJAK_PERUBAHAN` — **8 atribut** (6 sebelumnya + 1 dipecah + 1 ditambahkan)

> ### `F-3` DITUTUP 24 September 2026 — yang hilang **PERAN**, bukan nilainya
>
> Blok sebelumnya menulis: *"yang hilang **tampaknya** nilai yang berubah itu sendiri — dan
> 'tampaknya' bukan dasar."* **Dua hal ternyata salah di dalamnya**, dan keduanya disebut supaya
> pembaca berikutnya tahu ke mana harus melihat.
>
> **Pertama, nilainya SUDAH ada.** Baris `NILAI_SEBELUM` / `NILAI_SESUDAH` ditambahkan pada
> putaran yang sama; blok selisihnya yang tertinggal. Yang membuatnya tidak terlihat: **satu baris
> memuat dua nama**, sehingga pengurai `KAMUS-KOLOM.md` menolaknya — nama sel pertamanya bukan
> pengenal tunggal. Barisnya kini **dipecah dua**. Cacah baris dan cacah kolom menjadi satu angka — **7** — dan
> tidak ada lagi yang perlu memilih cara menghitung. Dengan `PERAN_PELAKU` di bawah, judulnya
> berbunyi **8**.
>
> **Kedua, yang benar-benar hilang terbaca dari ADR-0045**, dan tidak perlu ditebak. Isi minimalnya
> disebut di sana dalam satu kalimat:
>
> > *"Isi minimal: siapa, kapan, apa yang berubah dari nilai apa ke nilai apa, dan **di bawah peran
> > apa** — peran yang berlaku saat itu, bukan peran orangnya hari ini."*
>
> Empat yang pertama ada. **Butir kelima tidak ada sama sekali**, dan ADR-0045 menyebutnya bukan
> pelengkap melainkan bagian isi minimal.
>
> | | |
> |---|---|
> | **Yang ditambahkan** | `PERAN_PELAKU`, bertipe `T`, **tidak boleh kosong** |
> | **Kenapa `T` dan bukan rujukan** | alasannya **sama persis** dengan `PELAKU`, yang sudah bertipe `T` dan **dibekukan**: jejak harus menyatakan peran yang berlaku **saat itu**. Sebuah kunci asing ke tabel peran akan ikut berubah ketika penugasannya berubah, dan jejaknya berhenti membuktikan apa pun |
> | **Label** | **BARU** — seluruh entitas ini baru; sistem lama tidak punya jejak perubahan sama sekali |
> | **Arah dampak bila salah** | bila peran ternyata tidak perlu dibekukan, satu kolom teks berisi nilai yang dapat dihitung ulang. Salah ke arah sebaliknya — tidak dipasang, lalu dibutuhkan — **tidak dapat dipulihkan**: peran yang berlaku di masa lalu tidak dapat diadakan kemudian |
> | **Syarat pembalikan** | ADR-0044 dibatalkan, atau penugasan peran ditetapkan **tidak** bertanggal. Keduanya perubahan ADR, bukan perubahan di sini |
>
> #### Lubang yang ditemukan saat menutupnya, dan TIDAK ditambal — `F-16`
>
> **ADR-0044 menempatkan entitas peran dan penugasan bertanggal di gelombang 1**, dan §10 tidak
> memuat satu pun dari keduanya — sapuan kata `PERAN` atas seluruh `SPEC-MODEL-DATA.md` kembali
> **nol** sebelum putaran ini.
>
> | | |
> |---|---|
> | **Akibatnya** | `PERAN_PELAKU` di bawah berdiri sebagai **teks beku tanpa daftar nilai yang sah**. Ia dapat diisi apa saja, dan tidak ada yang dapat memeriksanya |
> | **Yang SUDAH diputuskan, 24 Sep 2026 sore** | **bentuknya** — `KTV-D` menetapkan ia **POTRET**, bukan rujukan, dan bentuk itu **tidak berubah** apa pun jadinya entitas peran nanti. Karena itu `F-16` **tidak menahan** irisan tiket `39`; yang masih ditagihnya **dari mana nilainya diambil** |
> | **Kenapa tidak ditambal di sini** | dua entitas baru adalah keputusan pemilik proses, bukan akibat samping penutupan `F-3`. §10.23c sudah menjadi presedennya: tiga dari empat usul entitas saya **gugur** setelah diperiksa |
> | **Siapa menutup** | pemilik proses, bersama sesi yang mengerjakan ADR-0044 |
> | **Yang menagih** | baris ini, dan `PERAN_PELAKU` yang berdiri tanpa acuan |

**Tidak ada di sistem lama.** Ia lahir dari ADR-0045: jejak perubahan sebagai fakta mesin.

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_JEJAK_PERUBAHAN` | *baru* | C1 | tidak | |
| `ID_VERSI_KONTRAK` | *baru* | R | tidak | |
| `WAKTU_PERUBAHAN` | *baru* | D | tidak | **`DATE` beresolusi detik.** Dua perubahan dalam detik yang sama tidak dapat dipisahkan olehnya — dan itu tidak merusak apa pun di sini, sebab entitas ini **sengaja tanpa kunci alami** |
| `PELAKU` | *baru* | T | tidak | dibekukan, sekeluarga dengan `NAMA_PEMUTUS` |
| `PERAN_PELAKU` | *baru* | T | tidak | **peran yang berlaku SAAT ITU** — ADR-0045 isi minimal butir kelima. **POTRET, bukan rujukan** (`KTV-D`): teks, final, **tidak dinormalisasi ulang** ketika entitas peran kelak lahir. Dari mana nilainya diambil masih `F-16` |
| `RUAS_YANG_BERUBAH` | *baru* | T | tidak | |
| `NILAI_SEBELUM` | *baru* | T | ya | teks, karena ruasnya beragam tipe |
| `NILAI_SESUDAH` | *baru* | T | ya | idem |

**Kunci alami — TIDAK ADA, dan itu keputusan**, dengan alasan yang sama seperti
`CATATAN_PERSETUJUAN`. **`SUMBU_REKONSILIASI`** — tidak punya baris turunan.

> **Bentuk `NILAI_SEBELUM`/`NILAI_SESUDAH` sebagai teks adalah pengecualian sadar** terhadap
> larangan satu kolom bermakna banyak: ia **bukan** fakta kontrak melainkan **catatan tentang
> fakta**, dan tidak ada rumus yang membacanya. Sesi DDL tidak boleh menegakkan tipe di atasnya.

---

### 10.21a `PERISTIWA_KONTRAK` — 5 atribut · **ENTITAS BARU, ditambahkan 24 Sep 2026**

**Dilaporkan sebagai penambahan, bukan diselipkan.** Ia lahir dari §10.20a: sistem lama menyimpan
**peristiwa** dan **keputusan** di satu daftar `CommentList`; sistem baru memisahkannya. **Itu
PERUBAHAN, bukan pelestarian** — dan tiket yang membawanya harus bergolongan demikian.

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_PERISTIWA_KONTRAK` | *baru* | C1 | tidak | |
| `ID_VERSI_KONTRAK` | *baru* | R | tidak | menggantung pada **versi** |
| `WAKTU` | `CommentList[].Date` | D | tidak | |
| `NAMA_PELAKU` | `CommentList[].OperatorName` | T | tidak | dibekukan sebagai teks, sekeluarga `NAMA_PEMUTUS` (§12.5) |
| `JENIS_PERISTIWA` | `CommentList[].Suggest` **diurai** | E | tidak | `SUNTINGAN_INTERNAL` · `ADDENDUM_EKSTERNAL` · `ADDENDUM_PREMI` · `DISALIN` · `REVISI_DIBUAT` |

> **`JENIS_PERISTIWA` adalah enumerasi, bukan teks bebas**, dan itu disengaja. Di sistem lama
> kelimanya adalah **kalimat yang dirangkai di dalam kode** — `OperatorID.pxInsName + "  " + "Had
> Created Internal Edit"`. Membawa kalimatnya berarti membawa nama pelakunya **dua kali**, sekali di
> kolomnya sendiri dan sekali di dalam teks. Peristiwa *"Copied from ID 1001378"* **tidak** membawa
> nomor kontraknya ke sini: ia sudah ada sebagai `ID_KONTRAK_DISALIN_DARI` di §10.1.

**Kunci alami — TIDAK ADA, dan itu keputusan**, dengan alasan yang sama seperti §10.20:
peristiwa yang sama dapat terjadi dua kali pada versi yang sama. **`SUMBU_REKONSILIASI`** — tidak
punya baris turunan.

**Batas yang disebut supaya tidak dibaca lebih jauh:** kelima nilai enumerasi adalah **kelima yang
ada di ekspor**. Bila baris `CommentList` ber-`IsApproved` kosong di produksi memuat `Suggest` di
luar kelimanya, ia **tidak punya rumah** — dan itu dihitung oleh **Uji AK**, bukan ditebak sekarang.

---

### 10.22 Enam tabel acuan — bentuk seragam

Seluruhnya dipakai P-48 dan dirujuk kunci asing dari entitas di atas. **Wajib tabel, bukan `CHECK`**
— himpunannya bertambah tanpa mengubah arti apa pun (ADR-0038, INV-62).

| Entitas | Isi awal dari | Kunci alami |
|---|---|---|
| `MATA_UANG` | `Int-CURRENCY` | `KODE` |
| `JENIS_POTONGAN` | empat nama di §14.2 | `KODE` |
| `JENIS_REASURANSI` | `Int-REINSURANCETYPE` | `KODE` |
| `BAHAYA` | empat kolom bernama bahaya, §10.18 | `KODE` |
| `KELOMPOK_TREATY` | `Int-TREATYGROUP` | `KODE` |
| `KELAS_BISNIS` | `Int-BUSINESS` / `Int-BUSINESSGROUP` | `KODE` |

Bentuk seragam keenamnya:

| Nama | Tipe | Boleh kosong | Catatan |
|---|---|---|---|
| `ID_‹entitas›` | C1 | tidak | |
| `KODE` | T | tidak | **kunci alami**, **Bentuk — tak berinduk**: unik di seluruh tabel |
| `NAMA` | T | tidak | |
| `AKTIF` | E | tidak | nilai lama **tidak pernah dihapus**; ia dimatikan |
| `ID_INDUK` | R | ya | **hanya `JENIS_REASURANSI`** — ia bersusun, §10.6 |

**Kunci alami dan lingkupnya** — `KODE`, **bukan Bentuk A maupun B**: keenamnya tidak punya induk,
sehingga lingkupnya seluruh tabel. **`SUMBU_REKONSILIASI`** — tidak punya baris turunan.

> **`STRUKTUR-DATA.md` §3 tidak menyatakan kunci alami keenamnya sama sekali** — bukan "tidak
> punya", melainkan **tidak dinyatakan**, dan itu persis kekosongan yang §5.0b larang. Dinyatakan di
> sini: **`KODE`**, dan alasan ia bukan `NAMA` adalah nama berubah ejaannya sementara kode tidak
> (`Overiding Commision` diperbaiki di tabel acuan, §14.2).
>
> **`AKTIF` bukan penghapusan.** Kontrak lama merujuk kode yang mungkin tidak dipakai lagi; menghapus
> barisnya memutus rujukan kontrak yang sudah disetujui.

---

### 10.23 Rekap — tingkat pencatatan, sumbu rekonsiliasi, dan yang masih terbuka

#### a. Paket uang yang tingkatnya BELUM DITENTUKAN

> ### DIPERBARUI 24 September 2026 — delapan dari empat belas terjawab dari ekspor
>
> Sapuan penulis dijalankan atas jalur sumber tiap paket (kolom *Asal* §10), dan **basis gelungnya
> dibaca**, bukan hanya ekspresinya. Rinciannya `2-to-spec/TINGKAT-PENCATATAN-14-PAKET-UANG.md`.
>
> **Kalibrasi lulus** — `NILAI_PENYEBARAN`, yang jawabannya sudah diketahui `BAGIAN_NURE`, memang
> terdeteksi: `Primary.SpreadingList(<LAST>).Value = (.Pct/100) * primary.RNMShareList(1).Value`.
>
> | Paket uang | Entitas | **Tingkat** | Dasarnya |
> |---|---|---|---|
> | `PREMI_BRUTO`, `PREMI_BRUTO_MINIMUM` | `BAGIAN` | **`BAGIAN_NURE`** | `.Value * @divide(Local.ShareCalc,100,4)`, dan `Local.ShareCalc = TreatyIn.RNMShare` |
> | `NILAI_TERMIN` | `TERMIN` | **`BAGIAN_NURE`** | basis gelungnya `TreatyIn.TotalShareNetNP` — total share net |
> | `NILAI_RETENSI` | `RETENSI_CEDANT` | **`TREATY_100_PERSEN`** | dua cabang: QS `× (100 - QSPct)`, Surplus masukan manual — keduanya membagi risiko treaty |
> | `KAPASITAS_SURPLUS` | `DETAIL_PROPORSIONAL` | **`TREATY_100_PERSEN`** | retensi × jumlah *lines* |
> | `MDP`, `MDP_MINIMUM` | `LAYER` | **`TREATY_100_PERSEN`** | basis gelungnya `TreatyIn.Limits` — limit layer, yang §10.3 sudah tetapkan 100% |
> | `CADANGAN_PREMI` | `DETAIL_PROPORSIONAL` | **`TREATY_100_PERSEN`** | basis gelungnya `.CessionList` — jumlah yang diserahkan cedant |
> | **`BATAS_MAKSIMUM_KELOMPOK`, `BATAS_MAKSIMUM_NON_KELOMPOK`, `BATAS_PILIHAN`** | `VERSI_KONTRAK` | **BELUM — `T-1`** | **nol penulis**; diketik di layar (18 kemunculan `Section` masing-masing) |
> | **`NILAI_BATAS`** | `BATAS_PER_BAHAYA` | **BELUM — `T-2`** | **nol penulis**; diketik di layar |
> | **`LIMIT_AGREGAT`** | `LAYER` | **BELUM — `T-3`** | tidak ada penulis yang menghitungnya; hanya penjumlah ringkasan lintas layer |
> | **`NILAI_EGNPI`** | `EGNPI` | **BELUM — `T-4`** | `= .AmountEPI`, **dibawa apa adanya dari penawaran**; tingkatnya diwarisi dari `EPI` |
>
> **Dan aturan pemilahannya sendiri dikoreksi.** *"Kali `Share`/`Pct` maka `BAGIAN_NURE"* menjawab
> **salah pada tiga dari empat belas**: `RetentionPct = 100 - QSPct` adalah porsi **cedant**,
> `Surplus` adalah **cacah lines**, dan `MDPPct`/`PremiumReservePct` adalah **tarif syarat kontrak**.
> Yang menentukan **apa yang dikalikan dan apa basis gelungnya**, bukan adanya perkalian.
>
> **Angka §7 semula — 16 / 22 / 5 — tidak dihapus.** Ia jejak, bukan salah ketik. Sesudah langkah
> ini hitungannya **21 `TREATY_100_PERSEN` · 25 `BAGIAN_NURE` · 6 belum**, dan keenam yang belum
> **masing-masing menyebut nomor pertanyaannya** supaya tidak ada lubang tanpa penagih (§2.9b).



§7 menyatakan 16 `TREATY_100_PERSEN`, 22 `BAGIAN_NURE`, 5 belum ditentukan — **sebagai hitungan,
tidak pernah sebagai daftar**. Yang ditemui saat §10 ditulis:

| Paket uang | Entitas | Tingkat |
|---|---|---|
| `LIMIT`, `DEDUCTIBLE` | `LAYER` | **`TREATY_100_PERSEN`** — tertulis §10.3 |
| `NILAI_PENYEBARAN` / `NILAI` | rantai penyebaran | **`BAGIAN_NURE`** — tertulis `STRUKTUR-DATA.md` §1.4: penyebaran adalah **bagian NuRe** ke susunan retro |
| `NILAI_EGNPI` | `EGNPI` | **BELUM — dan ia butir ESKALASI**, §7 |
| `BATAS_MAKSIMUM_KELOMPOK`, `BATAS_MAKSIMUM_NON_KELOMPOK`, `BATAS_PILIHAN` | `VERSI_KONTRAK` | **BELUM** |
| `LIMIT_AGREGAT`, `MDP`, `MDP_MINIMUM` | `LAYER` | **BELUM** |
| `NILAI_RETENSI` | `RETENSI_CEDANT` | **BELUM** |
| `NILAI_TERMIN` | `TERMIN` | **BELUM** |
| `NILAI_BATAS` | `BATAS_PER_BAHAYA` | **BELUM** |
| `PREMI_BRUTO`, `PREMI_BRUTO_MINIMUM` | `BAGIAN` | **BELUM** |
| `CADANGAN_PREMI`, `KAPASITAS_SURPLUS`, dan keempat daftar uang `Detail[]` | `DETAIL_PROPORSIONAL` | **BELUM** |

**Dua ditetapkan dari sumber tertulis, selebihnya BELUM.** Tidak satu pun ditebak. Ini **lubang yang
menahan sesi DDL**, karena INV-39 dan INV-40 adalah constraint dan migrasi tidak tahu nilai apa yang
ditulis ke kolomnya.

#### b. Sumbu rekonsiliasi — hanya empat entitas punya baris turunan

| Entitas | `SUMBU_REKONSILIASI` |
|---|---|
| `VERSI_KONTRAK` | `ID_VERSI_KONTRAK` |
| `LAYER` | `ID_LAYER` |
| `DETAIL_PROPORSIONAL` | `ID_DETAIL_PROPORSIONAL` |
| `BAGIAN` | `ID_BAGIAN` |
| `PENYEBARAN` | `ID_PENYEBARAN` |
| **`RINCIAN_PENYEBARAN`** | **`ID_JENIS_REASURANSI`, dikelompokkan dalam `ID_PENYEBARAN`** |
| **seluruh entitas lain** | **tidak punya baris turunan — dinyatakan, bukan dikosongkan** |

Hanya baris terakhir yang menjadi `GROUP BY` sebuah *materialized view* INV-50.

#### c. Empat keputusan — **seluruhnya dijawab 24 September 2026**

Keempatnya diangkat karena mengubah jumlah entitas. **Tiga dari empat berakhir berbeda dari usul
saya**, dan ketiganya berubah karena pemilik proses menuntut pemeriksaan sebelum jawaban.

| # | Soal | Usul semula | **Hasil** | Arah |
|---|---|---|---|---|
| 1 | pemulihan limit berulang (§10.3b) | entitas `PEMULIHAN_LIMIT` | **DITERIMA** — dan kemampuan *"syarat berbeda tiap pemulihan"* bergolongan **BARU**, P-59 | **+1** |
| 2 | rincian pembayaran termin (§10.16a) | entitas `RINCIAN_TERMIN` | **DITOLAK** — bukan entitas baru; **tingkatnya yang terbalik**, dan baris pengelompokan per mata uang justru **larut** | **0** |
| 3 | `NILAI_PENYEBARAN` (§10.8) | pertahankan | **DITANGGUHKAN** — bukti saya tidak memisahkan apa pun; **Uji AD** ditulis | **0** |
| 4 | `ID_SUSUNAN_RETRO` (§10.4a) | tabel acuan ketujuh | **DITOLAK** — ia **master di luar modul**, `Int-PROPORTIONALARRG`. Yang ditambahkan **kewajiban**, bukan entitas | **0** |

> **27 menjadi 28, bukan 30.** Satu entitas bertambah, dan tiga usul saya gugur setelah diperiksa.

**Pelajaran yang ditulis di sini karena ia akan terulang:** ketiga usul yang gugur punya bentuk yang
sama — **daftar di sistem lama disimpulkan sebagai entitas di sistem baru**. Dua kali kesimpulan itu
keliru karena daftarnya adalah **perancah clipboard** (pengelompokan per mata uang; master di luar),
dan sekali karena **buktinya tidak memisahkan apa pun**. Daftar bukan entitas sampai terbukti punya
identitas sendiri.

---

## 11. PEKERJAAN YANG BELUM DIKERJAKAN — penamaan dan tipe entitas selebihnya
### 11.1 Langkah yang belum dijalankan

| Langkah | Keluaran | Keadaan |
|---|---|---|
| 8 | `SPEC-INVARIAN.md` — invarian, kelengkapan per perpindahan | **selesai** — 63 invarian |
| 9 | `PETA-TELUSUR-JSON.md` — adjudikasi **per jalur**, 667 jalur | **selesai, semestanya kurang** — lihat `4-erd-dan-tabel-datar/AUDIT-PENYEBUT-POHON.md` |
| 10 | ERD — teks mengikat, gambar alat bantu | **selesai**, dijalankan ulang sesudah keputusan §10.23c |
| — | §10 berkas ini: nama, tipe, keterisian per atribut | **SELESAI 24 September 2026 — 27 entitas** |

### 11.2 Yang tersisa sesudah §10 selesai

§10 kini meliputi **seluruh 27 entitas penyerahan pertama**: §10.1–§10.2 tidak tersentuh titik buta,
§10.3–§10.4 **diselesaikan ulang atas sumber yang lengkap**, dan §10.5–§10.22 ditulis 24 September
2026 dengan kedua kewajiban `5-tiket/DAFTAR-PEKERJAAN.md` §5.0b berlaku pada seluruhnya.

| Yang tersisa | Siapa |
|---|---|
| **empat keputusan §10.23c** — ketiganya mengubah jumlah entitas | **pemilik proses** |
| **tingkat pencatatan 14 paket uang** (§10.23a) | sebagian uji data, satu butir eskalasi (`EGNPI`), sisanya slip kontrak |
| menomori **lima kunci alami** yang belum punya invarian | sesi to-spec, sesudah §10.23c |
| kepanjangan `RSMD`, `WPC`, `RIOGR`, `RIONR` | **wawancara** — korpus sudah disapu habis |
| menjalankan ulang peta nama tabel dan ERD struktur | perkakas, sesudah §10.23c |
| memperbarui angka §1, §3, §4, §7 | sesudah §10.23c, karena keputusannya mengubah penyebut |

**Yang TIDAK tersisa:** entitas tanpa §10. Tidak ada lagi.

### 11.3 Entitas yang sengaja TIDAK dikerjakan

| Entitas | Sebab |
|---|---|
| `RETRO_KELUAR` | GEL-2 |
| `PENCAPAIAN` | GEL-3 |
| `NILAI_SELISIH` | modul Adjustment, embargo |
| `BESARAN_DAPAT_DISESUAIKAN` | tabel acuan sambungan; tidak tersentuh satu pun kemampuan penyerahan pertama |

---

## 12. Koreksi 23 September 2026 — hasil penambalan terhadap inventaris kelas Pega

Lima persoalan diangkat pemilik proses setelah membaca bentuk awal berkas ini. Seluruhnya
diselesaikan dengan **sumber yang mandiri dari penyisiran JSON**: indeks rujukan aturan di dalam
ekspor, yang memasangkan `pxRuleClassName` dengan `pyRuleName` dan karena itu memberi daftar
properti **per kelas Pega** tanpa melewati awalan `TreatyIn.`.

Itu menutup lubang metode yang sudah saya sebut sendiri: penyisiran hanya menangkap jalur berawalan
`TreatyIn.`, sehingga daftar yang ditulis lewat *step page* tidak tertangkap.

### 12.1 Satu jebakan metode baru: berkas Section MEMBUNDEL aturan lain

`Section/TreatyInNONProportional.xml` memuat entri indeks milik **empat aturan berbeda**:

| Pemilik indeks | Entri |
|---|---|
| `TREATYINNONPROPORTIONAL` | 105 |
| `TREATYINTABSPROPORTIONAL` | **193** |
| `TREATYINTABSNONPROPORTIONAL` | 337 |
| `TREATYINTABSNONPROPORTIONALADJUSTPREMI` | 73 |

Sebagai perbandingan, `TreatyInTabsProportional.xml` memuat **satu** pemilik saja.

**Akibatnya:** menyimpulkan cabang sebuah properti dari **nama berkasnya** dapat salah total. Saya
sempat mengira `RNMShareP` dipakai di cabang non-proporsional karena ia muncul di berkas bernama
`TreatyInNONProportional.xml`; ia sebenarnya milik seksi proporsional yang **dibundel di dalamnya**.

Aturan yang ditambahkan: **cabang sebuah properti ditentukan oleh `pzIndexOwnerKey`, bukan oleh nama
berkas.**

### 12.2 `BAGIAN` tidak membawa identitas layer sebagai atribut

`Share[].Layer`, `.LayerPart`, `.LayerType`, `.LayerPartType` adalah **identitas layer**, bukan
atribut bagian. Di sistem lama ia disalin ke baris share supaya baris itu dapat dicocokkan kembali
ke layernya — denormalisasi yang lahir karena clipboard Pega tidak punya kunci asing.

Di model baru hubungan itu adalah **relasi**: `BAGIAN` → `LAYER`. Membawanya sebagai empat kolom
salinan berarti mewarisi denormalisasi ke DDL, dan dua salinan identitas layer akan berpisah —
bentuk yang sudah dihapus di tempat lain.

`.Cover` **tetap atribut `BAGIAN`**: ia muncul pada `Data-TreatyInShare` maupun `Data-TreatyInLimits`
dengan arti masing-masing, dan bukan salinan pengenal.

### 12.3 `BAGIAN` BUKAN satu entitas untuk kedua cabang — §8 dikoreksi

Uji "berbeda arti versus tidak dipakai" dijalankan pada sumbu yang salah. Sumbu yang menentukan
adalah **masukan versus turunan**, dan pada sumbu itu kedua sisi tidak sebanding:

| | Cabang non-proporsional | Cabang proporsional |
|---|---|---|
| Wadahnya | kelas `Data-TreatyInShare`, **49 properti** | **tidak ada kelas tersendiri** |
| Letaknya | baris `Share[]` di bawah kontrak | atribut **pada** `Data-TreatyInLimitsDetail` |
| Masukannya | `GrossPremiumList[].Value` — diketik orang | `RNMShare` per baris, bila `RNMShareAcrossTheBoard` mati |
| Sisanya | 10 daftar `RNMSpreadedList…XOL` — turunan | `RNMShareList`, `RNMSpreadedList(RI)` — turunan |

Pemilik proses benar: **bagian proporsional tidak punya baris sendiri.** Ia sekumpulan atribut pada
baris `DETAIL_PROPORSIONAL`, bukan entitas.

**Keputusan yang dikoreksi:**

- `BAGIAN` adalah entitas **cabang non-proporsional saja**, dengan relasi ke `LAYER`.
- Bagian NuRe pada cabang proporsional adalah **atribut `DETAIL_PROPORSIONAL`** (`RNMShare`), dan
  daftar-daftar turunannya tidak disimpan (ADR-0037).
- `POTONGAN` menggantung pada **keduanya**: `DeductionList` muncul pada `Data-TreatyInShare`
  **dan** pada `Data-TreatyInLimitsDetail`. Karena induknya dua, di sinilah tabel penghubung
  pernah dipertimbangkan — dan tidak diperlukan: `POTONGAN` cukup punya dua kunci asing
  bersyarat yang **saling meniadakan**, atau dua tabel sejenis. Bentuk finalnya ditetapkan sesi
  DDL, dan §12.6 mencatatnya sebagai keputusan yang belum diambil.

### 12.4 Lima entitas dilengkapi dari inventaris kelas

| Entitas | Kelas Pega | Yang hilang di bentuk awal, kini ditambahkan |
|---|---|---|
| `POTONGAN` | `Data-TreatyInDeduction` (7) | **seluruhnya** — `Deduction`, `DeductionPct`, `Currency`, `CurrencyID`, `Comment` |
| `TERMIN` | `Data-TreatyInInstallment` (13) | `Installment`, `InstallmentPct`, `Amount`, `DueDate`, `PaymentDate`, `WPC`, `CurrencyID` |
| `PERIODE_AKUMULASI` | `Data-TreatyInAccumulation` (4) | `SubDays`, `SubDueDate` |
| `SKALA_KOASURANSI` | `Data-TreatyInCoInScale` (2) | **seluruhnya** — `PctLimit`, `CoInShare`; dan ia memang **daftar**, bukan skalar |
| `PENYEBARAN` | `Data-TreatyInLimitsSpreading` (28) | `ReinsTypeID`, `ReinsTypeName`, `ParentReinsTypeID`, `Pct` |

`DeductionPctCalculate` **tetap dibuang** (§5): ia bendera mode hitung, bukan fakta potongan.

`Rp` dan `Usd` pada `Data-TreatyInLimitsSpreading` **dibuang** — pasangan kolom kembar untuk dua
mata uang, pola terlarang. Penggantinya adalah `NILAI_PENYEBARAN` berbaris per mata uang.

`PctTotal` dan `AmountTotal` pada `Data-TreatyInInstallment` **turunan**, tidak disimpan.

### 12.5 `PositionUsername` dibuang, dan pemindaian sejenisnya

`PositionUsername` adalah **nama orang yang menempel di kontrak sebagai teks**. Ia melewati ADR-0044
(peran dan penugasan bertanggal) dan usang begitu orangnya pindah jabatan — penyakit yang sama
dengan TD-01 di tempat lain. **Dibuang.** Penugasan yang berlaku dibaca dari peran, bukan dari teks
pada kontrak.

Pemindaian atribut sejenis di `VERSI_KONTRAK`, dan hasilnya:

| Atribut | Putusan |
|---|---|
| `PositionUsername` | **dibuang** — nama orang sebagai teks |
| `TreatyLeader`, `LeadingReinsName`, `LeadingReinsID` | **dipertahankan** — nama *perusahaan* reasuradur, bukan orang; menjadi rujukan ke master reasuradur |
| `CommentList[].OperatorName` | **dipertahankan** — ia fakta historis siapa memutuskan apa dan kapan; membekukan nama pada saat keputusan justru yang benar (ADR-0045) |

Pembedanya: nama yang menyatakan **siapa memegang kontrak sekarang** adalah turunan yang usang;
nama yang menyatakan **siapa melakukan sesuatu dahulu** adalah fakta yang beku.

### 12.6 Keputusan yang dipindahkan ke sesi DDL

| Hal | Sebab |
|---|---|
| Bentuk `POTONGAN` berinduk dua — dua kunci asing bersyarat atau dua tabel | bentuk fisik, bukan bentuk model |

---

## 13. `RNMShareP` — akhiran `P` berarti PROPORTIONAL, dan satu temuan grilling dikoreksi

### 13.1 Pembacaan yang diuji

Dua pembacaan bersaing, dan keduanya cocok dengan bukti awal:

| | Isi |
|---|---|
| Pembacaan grilling | `RNMShareP` = bagian **global**, dipakai saat `RNMShareAcrossTheBoard` menyala |
| Pembacaan sesi ini | akhiran `P` = **Proportional**, sama seperti tiga pasangan lainnya |

Bukti awal tidak membedakan: `TreatyInPropshare` menyetel `.RNMShare = TreatyIn.RNMShareP` ketika
`RNMShareAcrossTheBoard == true`, tetapi `TreatyInPropshare` **adalah** aktivitas proporsional, jadi
`RNMShareP` dibaca di cabang proporsional pada kedua pembacaan.

### 13.2 Uji yang membedakan: apakah ada `RNMShare` polos di cabang non-proporsional

Jawabannya **ya**, dan itu memutuskan. `TreatyIn.RNMShare` sebagai kata utuh dirujuk **19 kali di
12 berkas**, dan yang menentukan adalah sepasang kembar pemetaan:

| Aktivitas | Menulis |
|---|---|
| `TreatyInMappingDataconvert` (non-proporsional) | `TreatyIn.RNMShare` |
| `TreatyInMappingDataconvertProp` (proporsional) | `TreatyIn.RNMShareP` |

Dua aktivitas kembar yang berbeda **hanya pada akhiran `Prop`**, masing-masing menulis varian yang
berbeda **hanya pada akhiran `P`**. Itu pola yang identik dengan `BrokeragePercent`/`…P`,
`Exclusions`/`…P`, dan `SpecialConditions`/`…P`.

Pembaca sisi non-proporsional lainnya menguatkan: `SaveTreatyInOfferNonProp1_Act`,
`TreatyInNonAddItem` (5 rujukan), dan `Section/TreatyInTabsNonProportional`.

**Pembacaan grilling patah.** Akhiran `P` berarti **Proportional**, bukan global.

### 13.3 Maka apa sebenarnya `RNMShareAcrossTheBoard`

Ia dibaca sebagai **kondisi** hanya di tiga aktivitas, dan ketiganya proporsional:
`TreatyInPropshare`, `TreatyInPropshareDetail`, `CalculateShareList`. Di layar ia kotak centang,
tampil di tab proporsional maupun non-proporsional.

Artinya: ia **bendera cabang proporsional** yang memutuskan apakah bagian NuRe tingkat kontrak
(`RNMShareP`) diberlakukan ke **setiap baris detail**, atau tiap baris membawa bagiannya sendiri.
Ia tidak menyatakan "global versus lokal" pada tingkat kontrak; ia menyatakan **cara sebaran** di
dalam satu cabang.

### 13.4 Temuan grilling yang harus diperbaiki — dicatat sebagai koreksi, bukan penggantian diam-diam

**Temuan lama:** *"`AchievementPct` memakai bagian GLOBAL pada kontrak yang bagiannya per baris."*

**Temuan yang benar:** `GetAchievement` berjalan pada kelas
`ASM-FW-GISFW-Data-TreatyInLimitsDetail` — yaitu **baris detail proporsional** — dan rumusnya

```
.AchievementPct = @if(Local.ValueEPI=="", 0,
     @divide(@divide(.TotalAchPremium, @divide(TreatyIn.RNMShareP,100,20), 20),
             Local.ValueEPI, 20) * 100)
```

membagi dengan **bagian tingkat kontrak** `RNMShareP` **tanpa memeriksa `RNMShareAcrossTheBoard`
sama sekali** — aktivitas itu tidak termasuk ketiga yang membaca bendera tersebut.

Jadi pada kontrak proporsional yang **setiap baris detailnya membawa bagian sendiri**, pencapaian
tiap baris tetap dihitung terhadap bagian tingkat kontrak.

**Apa yang berubah dari temuan lama, dan apa yang tidak:**

| | Lama | Benar |
|---|---|---|
| Cabang terdampak | tersirat lintas cabang | **hanya proporsional** |
| Sebab | memakai "bagian global" | mengabaikan bendera cara sebaran |
| Arah kesalahan | tidak dinyatakan | **bergantung data, tidak searah** — bila bagian baris lebih kecil dari bagian kontrak, pencapaian **terlalu rendah**; bila lebih besar, terlalu tinggi |

Arah yang tidak searah membuatnya **lebih sulit terlihat** daripada kesalahan searah: sebagian baris
naik, sebagian turun, dan totalnya bisa tampak wajar. Ia tidak menggelembungkan dan tidak
menghilangkan secara konsisten — ia **mengaburkan**.

Ini **tidak** mengubah keputusan model mana pun: pencapaian adalah turunan yang dihitung saat
dibaca (ADR-0037), dan bagian yang dipakai adalah bagian yang berlaku pada baris itu. Yang berubah
adalah **uraian temuan** di `RINGKASAN-GRILLING.md` dan kedudukannya di daftar eskalasi.
---

## 14. Koreksi putaran kedua — 23 September 2026

### 14.1 `POTONGAN` adalah SATU KONSEP dengan dua pelekatan — diputuskan sekarang

Bentuk fisiknya milik sesi DDL, tetapi **keputusan logisnya tidak boleh ikut ditunda**, karena ia
menentukan berapa kali invarian potongan ditulis di langkah 8:

| Bila | Akibatnya |
|---|---|
| satu konsep | invarian ditulis **sekali**, berlaku di kedua pelekatan |
| dua konsep | invarian ditulis **dua kali**, dan keduanya akan berpisah dalam dua tahun |

**Bukti yang memutuskan — satu aktivitas melayani kedua pelekatan dengan satu rumus.**
`CalculateDeduction` menghitung

```
.DeductionList(i).Deduction = @divide(.DeductionList(i).DeductionPct, 100, 4)
                              * .GrossPremiumList(1).Value
```

dan pemanggilnya mencakup **kedua sisi**: `Share.xml`, `ShareRetro.xml`,
`TreatyInXOLAddSpreadingDetail` di sisi non-proporsional, dan `DetailLimits.xml` di sisi
proporsional. Rumusnya tidak bercabang.

Nama potongan yang benar-benar dipakai juga identik di kedua sisi: `Brokerage fee`,
`Facultative Brokerage fee`, `Overiding Commision`, `Comm to NuRe`.

**Keputusan: `POTONGAN` satu konsep, dua pelekatan.** Bentuk fisiknya — dua kunci asing yang saling
meniadakan, dua tabel sejenis, atau bentuk lain — ditentukan sesi DDL, **dengan satu syarat yang
mengikat: aturannya tidak boleh tertulis dua kali.**

### 14.2 Atribut `POTONGAN` yang sebenarnya — `Comment` ternyata NAMA potongannya

Pemeriksaan penulisnya membalik pembacaan awal saya:

| Properti lama | Sebenarnya | Perlakuan |
|---|---|---|
| `Comment` | **nama potongan** — diisi teks `"Brokerage fee"`, `"Facultative Brokerage fee"`, `"Overiding Commision"`, `"Comm to NuRe"` | menjadi **rujukan ke tabel acuan jenis potongan** |
| `DeductionPct` | persentase yang disepakati | **masukan** |
| `Deduction` | hasil `pct/100 × premi bruto` | **turunan**, tidak disimpan (ADR-0037) |
| `Currency`, `CurrencyID` | mata uang dari nilai turunan itu | ikut turunan, tidak disimpan |
| `DeductionPctCalculate` | bendera mode hitung | **dibuang** |

Atribut `POTONGAN` karena itu:

| Nama | Asal | Tipe | Boleh kosong |
|---|---|---|---|
| `ID_POTONGAN` | *baru* | C1 | tidak |
| `ID_INDUK_POTONGAN` | *baru* | R | tidak |
| `ID_JENIS_POTONGAN` | `Comment` | R | tidak |
| `DASAR_PERHITUNGAN` | *baru* | E | tidak |
| `PERSEN_POTONGAN` | `DeductionPct` | P2 | tidak |

**Jenis potongan WAJIB tabel acuan, bukan CHECK** — himpunannya bertambah tanpa mengubah arti apa
pun (ADR-0038). Nama seperti `Overiding Commision` yang salah eja di sistem lama diperbaiki di
tabel acuan, dan ejaan lamanya tercatat di peta telusur.

**`DASAR_PERHITUNGAN` adalah atribut baru, dan itu keputusan saya.** Glosarium membedakan potongan
menurut **dasar perhitungannya**: terhadap premi bruto, terhadap laba (*profit commission*), atau
jumlah yang disepakati. Sistem lama hanya pernah menghitung terhadap premi bruto — satu rumus, tanpa
cabang — sehingga dasarnya tersirat dan tidak pernah disimpan. Menyimpannya sekarang membuat
*profit commission* (yang masih DITUNDA, §6) punya tempat tanpa mengubah bentuk nanti.

> **Yang membatalkannya:** bila bisnis menyatakan ketiga dasar itu bukan sifat sebuah potongan
> melainkan jenis potongan yang berbeda sama sekali, `DASAR_PERHITUNGAN` larut ke tabel acuan jenis
> dan kolomnya hilang. Satu kolom, bukan perubahan bentuk.

### 14.3 `FacShare` BUKAN pasangan cabang — ia SALINAN, dan teori pasangan tidak berlaku

Peringatan pemilik proses tepat, dan penandaannya memang menyimpan sesuatu. `FacShare` bertanda
`NONPROP` sementara `FacultativeShare` bertanda `UMUM` — dan itu bukan penandaan yang keliru.

**Bukti yang memutuskan**, di `TreatyInXOLAddSpreading`:

```
TreatyIn.FacShare          <-  TreatyIn.FacultativeShare
TreatyIn.FacShareBrokerage <-  TreatyIn.FacultativeShareBrokerage
```

Penyalinan lurus, tanpa perhitungan, di dalam **satu aktivitas yang membaca yang satu dan menulis
yang lain**. Langkah di sebelahnya bahkan berlabel *"untuk cek ada atau tidak
TreatyIn.FacultativeShare"*.

Jadi keduanya **bukan dua varian cabang dari satu fakta**, melainkan **satu fakta dan salinannya**:
`FacultativeShare` sumbernya, `FacShare` potret yang diambil saat penyebaran XOL disusun.

**Perlakuan berubah, meski hasil akhirnya mirip:**

| | Bila dianggap pasangan | Yang benar |
|---|---|---|
| Tindakan | dilebur jadi satu atribut | `FacShare` dan `FacShareBrokerage` **dibuang sebagai salinan** |
| Dasarnya | `SIFAT_PROPORSI` sudah menyatakan cabang | ADR-0023 dan ADR-0041 — satu bentuk kanonik, satu penulis |
| Yang tertinggal | tidak ada | **temuan**: potret itu bisa basi |

**Temuan yang menyertainya.** Layar membaca `FacShare`, sementara yang disunting orang adalah
`FacultativeShare`. Bila `FacultativeShare` berubah setelah penyebaran XOL disusun, `FacShare`
tetap memegang nilai lama dan **layar menampilkan angka yang sudah tidak berlaku**. Ini bentuk yang
sama dengan seluruh keluarga salinan yang model baru hapus, dan ia hilang dengan sendirinya begitu
salinannya tidak dibawa.

**Maka peleburan pasangan bercabang tinggal TIGA, bukan lima:** `RNMShare`/`…P`,
`BrokeragePercent`/`…P`, `Exclusions`/`…P`, `SpecialConditions`/`…P` — empat — dikurangi
`FacShare`/`FacultativeShare` yang bukan pasangan, ditambah `AccountingMode`/`…NonProp` yang
**belum diputuskan** (§14.4). Jadi **empat pasti, satu menunggu**.

### 14.4 `AccountingMode` / `AccountingModeNonProp` — belum diputuskan, dan tidak ditebak

Keduanya **tidak pernah ditulis oleh satu aturan pun**: keduanya diisi dari layar, dan keduanya
tampil di seksi yang sama. Penulisnya tidak dapat membedakan, jadi bentuknya tidak dapat diputuskan
dari ekspor.

Dua bentuk yang mungkin, dan akibatnya berbeda:

| Bentuk | Artinya | Bila dilebur |
|---|---|---|
| pasangan cabang | tiap cabang mengisi variannya sendiri | benar |
| dasar + pengecualian | satu mode umum, satu penimpaan khusus non-proporsional | **salah** — pembedaan antara "tidak diisi" dan "sengaja sama dengan dasarnya" hilang |

Namanya sendiri memberi petunjuk ke arah kedua: varian pertama **tidak** bernama proporsional.

**Uji X-2** memutuskannya, dan bentuknya persis yang diminta: berapa kontrak mengisi salah satu
saja, **dipilah per sifat proporsi**. Pasangan cabang yang benar menunjukkan pola bersih — silangnya
nyaris kosong; dasar-plus-pengecualian tidak.

Sampai hasilnya kembali, keduanya **tetap dua atribut** di §10.2. Nama varian kedua memakai
`XOL` dan bukan `NON_PROPORSIONAL`, karena yang terakhir 31 bita — melewati batas §16. `XOL`
istilah pasar yang sudah ada di glosarium, jadi ia penggantian **kata**, bukan penyingkatan. Meleburkan lebih dulu adalah
keputusan yang tidak bisa dibatalkan oleh data.

### 14.5 `CedingStatusActive` dan `SourceStatusActive` DIBUANG — tanpa wawancara

> ### ⚠ DASARNYA MELEMAH 24 September 2026 — L-10
>
> Alasan pembuangan di bawah bersandar pada *"keduanya tidak pernah ditulis satu aturan pun"*.
> Sapuan penulis kita mencari `Property-Set` di aktivitas dan *data transform* — dan **ekspor ini
> tidak memuat satu pun `Rule-Declare-Expressions` maupun `Rule-Declare-Trigger`**, dua jenis aturan
> yang **menulis nilai tanpa dipanggil siapa pun**.
>
> **Ini satu-satunya tempat di seluruh modul yang memakai ketiadaan penulis sebagai ALASAN
> MEMBUANG**, dan karena itu ia yang paling patut dicurigai dari kelima klaim terdampak
> (`4-erd-dan-tabel-datar/JENIS-ATURAN-TAK-TEREKSPOR.md` §3.1).
>
> **DIPERIKSA, DAN HASILNYA MEMBALIK ARAH — 24 September 2026.**
>
> Klaim *"tidak pernah ditulis satu aturan pun"* memang **SALAH**, tetapi bukan karena L-10.
> **Karena bentuk penulisan yang tidak pernah kami cari.** `Harness/InputTreatyInOffer.xml`:
>
> ```xml
> <pyDisplayProperty>.StatusActive</pyDisplayProperty>
> <pyPropertyTarget>TreatyIn.CedingStatusActive</pyPropertyTarget>
> <pySetValueOnSelect>true</pySetValueOnSelect>
> ```
>
> Ketika orang memilih satu cedant dari master, kolom `StatusActive` baris itu **disalin** ke
> kontrak. Tidak ada `Property-Set` di mana pun, dan sapuan kami hanya mengenali `<PropertiesName>`.
>
> **Alasan pertama dicabut. Alasan kedua NAIK dari `DECIDED` menjadi `EVIDENCED`** — pemetaan itu
> **memperlihatkan penyalinannya terjadi**, dari baris master, pada saat dipilih, dengan berkas dan
> baris yang dapat ditunjuk.
>
> **Putusan membuang kedua atribut BERTAHAN, di atas dasar yang terbukti**, dan §14.5 **keluar**
> dari daftar yang harus dilihat ulang saat L-10 dijawab.

Pemilik proses benar bahwa ADR-0023 sudah memutuskannya, dan pengecualian yang ia minta diperiksa
**tidak berlaku di sini**.

Pengecualiannya: bila status itu status **pada saat akseptasi** dan tidak pernah diperbarui, ia
fakta beku, bukan salinan — bentuk yang sama dengan `CommentList[].OperatorName`.

**Pemeriksaannya menutup kemungkinan itu:** keduanya **tidak pernah ditulis oleh satu aturan pun**
di seluruh ekspor. Sesuatu yang tidak pernah ditulis tidak mungkin dibekukan pada suatu saat —
tidak ada saat itu. Keduanya nilai yang muncul saat halaman dimuat dan dibaca layar serta
`TreatyInCheckCedingBlacklist`.

**Dibuang, keduanya, sebagai salinan data master** (ADR-0023, ADR-0041). Status aktif cedant dan
asal bisnis dibaca dari masternya saat ditanya.

Ini juga instans §2.9 CONTEXT.md: keduanya menyatakan **keadaan sekarang** sebuah benda lain, bukan
peristiwa yang pernah terjadi pada kontrak ini.

### 14.6 `RSMD` — dugaan yang diperlakukan sebagai dugaan

Dugaan pemilik proses: **Riot, Strike, Malicious Damage** — perluasan polis properti yang baku di
pasar Indonesia, sering disebut bersama *civil commotion* sebagai RSMDCC.

**Diperlakukan sebagai dugaan, bukan fakta.** Yang dapat dipastikan dari ekspor: keempat kolom
bahaya tidak pernah ditulis satu aturan pun — seluruhnya diisi dari layar — dan keempatnya muncul
di `TreatyInTabsNonProportional`, jadi mereka **milik cabang non-proporsional**.

Verifikasinya murah dan tidak menahan apa pun: **Uji Y** menghitung sebaran keempat bahaya menurut
kelas bisnis. Terpusat di kelas properti berarti dugaan bertahan; tersebar rata berarti dugaan
patah dan artinya ditanyakan ke bisnis.

Apa pun hasilnya, **bentuknya tidak berubah**: keempatnya tetap anggota tabel acuan bahaya di dalam
`BATAS_PER_BAHAYA`, bukan empat kolom. Yang bergantung hasil uji hanyalah **entri glosarium**
untuk `RSMD`.

### 14.7 `RIOGR` dan `RIONR` — digabung ke permintaan slip yang sudah ada

Dugaan pemilik proses, dicatat **sebagai dugaan**: keduanya menyangkut **dasar tarif** — apakah
penyerahan dihitung atas *original gross rate* atau *original net rate*. Di treaty proporsional itu
ketentuan yang lazim dan tertulis di slip.

Karena itu keduanya **tidak butuh wawancara tersendiri**; keduanya butuh **dokumen yang sama** yang
sudah diminta.

**Permintaan dokumen itu kini menjawab lima hal sekaligus**, dan itu menjadikannya butir paling
berharga di seluruh daftar tunggu:

| # | Yang dijawab |
|---|---|
| 1 | tingkat pencatatan *claim cooperation* |
| 2 | tingkat pencatatan *cash loss* |
| 3 | arti dan tingkat pencatatan `PLA` |
| 4 | kepanjangan dan arti `RIOGR` / `RIONR` |
| 5 | kemungkinan juga prorata waktu pada reinstatement |

Yang diminta: **satu slip treaty proporsional dan satu slip treaty non-proporsional**, keduanya
lengkap dengan klausulnya. Bukan penjelasan lisan.

### 14.8 Sifat "saling meniadakan" dicatat di daftar eskalasi

Kesalahan pencapaian membuat baris berbagian kecil tercatat terlalu rendah dan yang besar terlalu
tinggi. Keduanya **saling meniadakan di agregat**.

Itu bukan sekadar catatan teknis; ia mengubah cara butir eskalasi dibaca:

> **Kesalahan yang saling meniadakan di agregat lebih sulit ditemukan, dan karena itu lebih lama
> bertahan.** Ia tidak pernah memicu pertanyaan dari siapa pun yang membaca total.

Ia melengkapi arah yang sudah tercatat di butir 6 daftar eskalasi — kesalahan yang menggelembungkan
mengundang pertanyaan, yang menghilangkan tidak. Yang **saling meniadakan** adalah kelas ketiga,
dan yang paling sunyi dari ketiganya.

Rancangan ujinya disesuaikan: **Uji V-5** membandingkan **per baris detail** dan melaporkan sebaran
simpangan ke dua arah, bukan jumlahnya. **Uji V-6** memberi hitungan termurah — berapa kontrak yang
baris detailnya punya bagian berbeda satu sama lain. Nol berarti cacat ini **bom waktu, bukan
kerusakan berjalan**.
