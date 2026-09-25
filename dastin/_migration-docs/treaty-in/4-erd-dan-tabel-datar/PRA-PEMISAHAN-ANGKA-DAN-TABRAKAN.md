# Dua angka sebelum menggambar, dan lima tabrakan yang tidak saya selesaikan sendiri

**Tanggal:** 24 September 2026
**Untuk:** arahan pemisahan tabel proporsional / non-proporsional

> **Dilaporkan sebelum ERD digambar**, sesuai urutan yang ditetapkan. Angka dilaporkan, **tidak
> dinilai**. Tabrakan dilaporkan, **tidak diselesaikan diam-diam**.

---

## 1. `M_TREATY_IN` dua kolom — dari DDL asli, bukan dari pembacaan

```sql
CREATE TABLE "POOLDATA"."M_TREATY_IN"
   ( "ID" VARCHAR2(10),
     "JSONDATA" CLOB,
     CONSTRAINT "ENSURE_M_TREATY_IN_JSON" CHECK (JSONDATA IS JSON) … )
```

**Dua kolom. Tidak ada kolom waktu.** Tiga akibat, dan ketiganya menutup sesuatu:

| | |
|---|---|
| **a. Eskalasi "nilai yang pernah disetujui tidak dapat direproduksi"** | selama ini bersandar pada **pembacaan** bahwa tidak ada kolom waktu. Sekarang bersandar pada **DDL**. Dasarnya berubah **DECIDED → EVIDENCED**, sumber: `Table/PEGA_M_TREATY_IN_DETAIL.txt`. Ini penerapan pertama aturan `BENTUK-TIKET.md` §7.2, dan ia **berlaku mundur** |
| **b. Butir 5.2 arahan terjawab** | *"ada satu tabel yang berasal dari `M_TREATY_IN`, yang sebelumnya memakai JSON"* — **`M_TREATY_IN` ADALAH tabel JSON itu.** Di model baru ia **arsip tanpa jalur baca** (ADR-0034, INV-61). Tidak ditebak |
| **c. Penulis `TREATY_IN` — diperiksa, dan INV-60 AMAN** | **nol** `INSERT`/`UPDATE` langsung ke `TREATY_IN` dari aturan mana pun. Satu-satunya penulisnya adalah prosedur `PEGA_TREATY_IN`, yang menulis `M_TREATY_IN` **dan** `TREATY_IN` **dalam panggilan yang sama**. Jadi faktanya memang ada di dua tempat, tetapi **penulisnya satu** |

> Akibat (c) pada arahan 5.1: `TREATY_IN` adalah **proyeksi yang dipelihara penulis yang sama**,
> bukan sumber mandiri. Di model baru ia jatuh ke aturan tabel datar — **tidak pernah ditulis
> aplikasi, ia diturunkan** (`TABEL-DATAR.md` §0).

---

## 2. Inventaris ekspor — empat jenis berkas yang tidak pernah dibuka

Kejadian ketiga dengan bentuk yang sama, jadi kali ini disapu lebih dulu.

**Pertama, temuan yang mengubah cara membaca folder ini:** `Table/` (7 DDL) dan `PROCEDURE/`
(5 badan prosedur) **BUKAN bagian ekspor Pega**. Keduanya ada di `_migration-docs/treaty-in/` —
disediakan terpisah, dan karena itu tidak pernah ikut tersapu perkakas mana pun yang menyisir
`D:\XML_NURE\Treaty In`.

| Folder ekspor | Berkas | Pernah dibaca? |
|---|---:|---|
| `Activity` | 149 | ya, berkali-kali |
| `Section` | 50 | ya |
| `RDBList` | 40 | **sebagian** — `pyBrowseSQL`-nya belum pernah disisir menyeluruh |
| `DataTransform` | 35 | sebagian |
| `FlowAction` | 32 | sebagian |
| `ReportDefinition` | 16 | **baru hari ini**, dan langsung menjawab susunan retro |
| `Harness` | 3 | **belum** |
| `SystemSettings` | 1 | **belum — dibuka hari ini, lihat di bawah** |
| `ConnectREST` | 1 | **belum** |
| `DecisionTable` | 1 | **belum** |
| `When` | 1 | **belum — dibuka hari ini** |
| `Struktur_InputTreatyInOffer.xlsx` | 1 | **belum** |

**Yang dibuka hari ini:**

| Berkas | Isinya |
|---|---|
| `SystemSettings/LinkService.xml` | `pyValue = 3`, dan keterangannya sendiri memetakan **3 = Quality assurance** (1 Sandbox · 2 Development · 4 Staging · 5 Production). Dibuat 4 September 2025 |
| `When/TreatyMasterInEDM.xml` | membandingkan `TreatyIn.EDMState` dengan `"1"`, `"2"`, `"3"` — menguatkan ADR-0049 (jenis addendum satu sumbu). Kelasnya `Data-Portal`, jadi ia **syarat menu**, bukan aturan bisnis |
| `DecisionTable/GetMimeType.xml` · `ConnectREST/ServiceGoogle.xml` | lampiran dan penyimpanan berkas — **di luar model kontrak** |

> ### `LinkService = 3` adalah temuan yang harus dibaca dua arah
>
> **Bacaan pertama:** ekspor ini diambil dari lingkungan **Quality assurance**, sehingga aturan yang
> kita baca sebelas sesi adalah **ruleset QA**, dan kesamaannya dengan produksi **tidak pernah
> diperiksa**.
>
> **Bacaan kedua:** ini ruleset produksi, dan *link service*-nya memang menunjuk QA — yang berarti
> **produksi memanggil lingkungan QA**, dan itu cacat tersendiri.
>
> **Keduanya berakibat, dan tidak satu pun dapat dipilih dari ekspor.** Ia **pertanyaan ke kantor**,
> satu kalimat: dari lingkungan mana ekspor ini diambil. **Penagihnya:** setiap kalimat berbentuk
> *"sistem lama melakukan X"* di seluruh dokumen ini bergantung padanya.

---

## 3. 98 rujukan kolom → **nol yang patut dicurigai**, tanpa DBA

Kebenaran dasar kedua yang sudah ada dan belum dipakai: **penulisnya**. Badan prosedur memuat
`INSERT`/`UPDATE` yang menyebutkan kolomnya satu per satu. Perkakas
`alat/panen-kolom-dari-penulis.py`.

| Tabel | Kolom dipanen | Sumber |
|---|---:|---|
| `TREATY_IN` | **20** | `PEGA_TREATY_IN.txt` |
| `TREATY_IN_EDM` | **24** | `PEGA_M_TREATY_IN_EDM.txt` |
| `TREATYINDETAIL` | 62 | `PEGA_M_TREATY_IN_DETAIL.txt` |
| `TREATYINDETAILEDM` | 55 | `PEGA_M_TREATY_IN_DETAIL_EDM.txt` |
| `M_TREATY_IN`, `M_TREATY_IN_EDM` | 2 dan 4 | idem |

**Angka 20 pada `TREATY_IN` bukan kebetulan** — ia persis klaim yang dipegang sejak awal
(*"dua puluh kolom bisnis"*), dan kini **terbukti dari sumber kedua yang mandiri**.

Hasil pemeriksaan berkas permintaan DBA:

| | |
|---|---:|
| rujukan kolom **terverifikasi tanpa DBA** | **42** |
| **patut dicurigai** | **0** |
| belum terperiksa — tak ada sumber apa pun | **5** (`M_TREATY_YEAR` 3, `TREATYEXCHANGEYEARLY` 1, `TREATYINOFFER` 1) |

> **Batas kemampuan, dinyatakan:** perkakas ini dapat menyatakan sebuah kolom **ADA** — karena ada
> yang menulisinya. Ia **tidak dapat** menyatakan sebuah kolom tidak ada. **Uji A-0 tetap dikirim
> dan tetap pertama**; yang berubah, bila ia kembali dengan galat, jumlahnya kecil dan sudah diduga.

---

## 4. Delapan dari 67 kalimat "satu baris" memang salah turun

Dugaan Anda tepat, dan angkanya **8 dari 67**. Ketujuh puluh dua sisanya benar.

Pemeriksaannya struktural, bukan pembacaan: komentar dibuang, kedalaman kurung dihitung, `SELECT`
terluar diambil sebagai `SELECT` terakhir di kedalaman nol, lalu diperiksa apakah ada `GROUP BY` di
kedalaman nol **sesudahnya**. Kedelapan yang gagal diganti dengan kalimat yang diturunkan dari
`SELECT` terluar, dan membawa peringatannya sendiri:

```
-- SATU BARIS = satu angka ringkasan atas seluruh himpunan; kueri TERLUAR
--              tidak mengelompokkan, jadi hasilnya tepat satu baris.
--              (dirumuskan dari SELECT TERLUAR; GROUP BY di dalam kurung
--               TIDAK menentukan bentuk hasil akhir)
```

Pemeriksaan ulang: **nol yang masih salah turun.** Dan seperti Anda duga, kedelapannya berada di
kelompok kueri yang sama — yang bertingkat — yaitu kelompok yang sama yang penuh risiko pada setiap
hal yang dikerjakan mekanis.

---

## 5a. ENTITAS YANG DIPAKAI KEDUA CABANG — dibuktikan dari ekspor

Metodenya: seluruh sasaran `Property-Set` disapu, dikelompokkan menurut wadah akarnya, lalu
**penulisnya didaftar**. Bukan dari contoh JSON.

### `Limits[]` — **DIPAKAI KEDUA CABANG. Terbukti.**

Tujuh aktivitas menulisnya, dan keduanya hadir:

| Cabang | Penulis |
|---|---|
| **proporsional** | `TreatyInMappingDataconvertProp.xml`, `TreatyInPropAdd.xml` |
| **non-proporsional** | `TreatyInMappingDataconvert.xml`, `CopyLastLimitNP.xml`, `TreatyInNonAddItem.xml` |
| lain | `TotalEgnpi.xml`, `TreatyInFixSpl.xml` |

> **`LAYER` tetap SATU tabel.** Hanya anaknya yang terpisah. Dugaan pemilik proses benar, dan
> sekarang ia terbukti dari ekspor.

### `Share[]` — **non-proporsional saja. Terbukti.**

Kelima penulisnya `SetSpreadingXOL`, `TreatyInXOLAddSpreadingDetail`,
`TreatyInXOLAddSpreadingDetailActual`, `TreatyInNonAddItem`, `TreatyInSetBrokerage` — tidak satu pun
proporsional. Menguatkan §12.3, yang sampai sekarang bersandar pada bentuk kelas.

### Daftar lengkap, dan apa yang terjadi padanya

| Golongan | Entitas | Perlakuan pada pemisahan |
|---|---|---|
| **kepala** | `KONTRAK`, `VERSI_KONTRAK` | lihat 5b — **tabrakan T-1** |
| **dipakai kedua cabang** | `LAYER` | **tidak dipisah** |
| **anak versi, dipakai kedua cabang** | `MATA_UANG_KONTRAK`, `RETENSI_CEDANT`, `EGNPI`, `PORTOFOLIO`, `PERIODE_PELAPORAN`, `PERIODE_AKUMULASI`, `TERMIN`, `SKALA_KOASURANSI`, `BATAS_PER_BAHAYA`, `DOKUMEN_KONTRAK`, `CATATAN_PERSETUJUAN`, `JEJAK_PERUBAHAN` | **tidak dipisah** — 12 entitas |
| **sudah khas cabang** | `DETAIL_PROPORSIONAL` (prop) · `BAGIAN` (non-prop) · `PEMULIHAN_LIMIT` (non-prop) | sudah terpisah **sejak model disusun**; pemisahan tabel tidak menambah apa pun |
| **berinduk dua** | `POTONGAN`, `PENYEBARAN`, `RINCIAN_PENYEBARAN`, `NILAI_PENYEBARAN` | **tabrakan T-3** |

---

## 5b. BERAPA BANYAK KOLOM YANG SAMA — dua angka, dengan dasarnya masing-masing

### `LAYER` — **13 dari 15 ruas hanya ditulis cabang non-proporsional**

Dari penyapuan penulis atas ruas langsung di bawah `Limits`:

| Golongan | Jumlah | Ruas |
|---|---:|---|
| **hanya NON-PROP** | **13** | `Limit`, `Deductible`, `Cover`, `AdjRate`, `MDPPct`, `ROLPct`, `ReinstatementPct`, `Layer`, `LayerPart`, `LayerType`, `LayerPartType`, `Currency`, … |
| hanya PROP | 1 | `TreatyType` |
| kedua cabang | 1 | `ID` |

> **Pembacaan angkanya — dan ini laporan, bukan penilaian:** sebuah "layer" proporsional hampir
> tidak punya atribut sendiri. Ia wadah bagi `Detail[]`. Hampir seluruh atribut `LAYER` di §10.3
> adalah besaran non-proporsional.

### `VERSI_KONTRAK` — **46 dari 47 atribut identik**

Dasarnya `SPEC-MODEL-DATA.md` §10.2 sendiri, bukan sapuan penulis:

| | |
|---|---:|
| atribut `VERSI_KONTRAK` | 47 |
| khas satu cabang | **1** — `CARA_PEMBUKUAN_XOL`, dan ia **masih menunggu Uji X-2** |
| identik bila dipisah | **46** |

Keempat pasangan berakhiran `P` **sudah dilebur** §10.0a menjadi empat atribut tunggal, jadi
keduanya tidak lagi menyumbang perbedaan.

> **Angkanya dilaporkan, penilaiannya tidak.** Yang perlu pemilik proses ketahui: memisahkan kepala
> menduplikasi **46 dari 47** kolom, sementara memisahkan anak-cabang **tidak menduplikasi apa pun**
> karena keduanya memang sudah berbeda entitas.

**Batas angka pertama, dinyatakan:** penggolongan penulis bersandar pada daftar aktivitas
proporsional dan non-proporsional yang saya susun dari nama dan perilakunya. Untuk ruas **kepala**
metode itu meninggalkan **31 ruas tak terklasifikasi** — karena penulisnya aktivitas simpan/muat
yang melayani keduanya — sehingga **angka kepala tidak saya ambil dari sapuan**, melainkan dari
§10.2. Dua angka, dua dasar, dan keduanya disebut.

---

## 6. LIMA TABRAKAN yang tidak saya selesaikan sendiri

| # | Tabrakan | Dengan apa |
|---|---|---|
| **T-1** | **Arahan 5.1 menyatakan kepala mengikuti `TREATY_IN`.** `TREATY_IN` punya **20 kolom**; `VERSI_KONTRAK` §10.2 punya **47**, dan `KONTRAK` §10.1 punya 8. Mengikuti `TREATY_IN` secara harfiah **membuang 35 atribut** yang §10 baru saja tetapkan | `SPEC-MODEL-DATA.md` §10.1, §10.2 |
| **T-2** | **Arahan 5.4 memisahkan tabel prop/non-prop.** §8 memutuskan sebaliknya: *"satu entitas dengan pembeda, kecuali satu pasang"* | `SPEC-MODEL-DATA.md` §8 |
| **T-3** | **`POTONGAN` dan `PENYEBARAN` berinduk dua** (§14.1) dengan **syarat mengikat: aturannya tidak boleh tertulis dua kali** (INV-63). Bila tabel prop dan non-prop terpisah, keduanya menjadi dua tabel per cabang — dan aturannya tertulis dua kali, persis yang dilarang | §14.1, INV-63 |
| **T-4** | **INV-32, INV-33, INV-34 berubah golongan** bila tabelnya terpisah — dari INDEKS UNIK menjadi TIDAK DAPAT DILANGGAR. Hitungan §1 `SPEC-INVARIAN.md` berubah: INDEKS UNIK **4 → 1** | `SPEC-INVARIAN.md` §1, §2.8 |
| **T-5** | **Arahan 5.3 "tidak ada JSON lagi"** berdampingan dengan **ADR-0034**, yang justru mempertahankan dokumen JSON **sebagai arsip tanpa jalur baca**. Keduanya dapat berdiri bersama — arsip bukan bentuk simpan — tetapi itu **pembacaan saya**, dan ia perlu disahkan | ADR-0034, INV-61 |

**Tidak satu pun saya selesaikan.** T-1 yang paling berakibat: ia menentukan apakah kepala punya
20 kolom atau 47, dan seluruh ERD bergantung padanya.

---

## 7. Uji §2.8 atas INV-32, INV-33, INV-34 — dijalankan

Pertanyaannya: **keadaan seperti apa yang harus ada agar pelanggarannya dapat dinyatakan.**

| Invarian | Keadaan yang diperlukan untuk melanggarnya, bila tabelnya terpisah |
|---|---|
| **INV-32** — `DETAIL_PROPORSIONAL` hanya di bawah kontrak proporsional | Melanggarnya menuntut sebuah baris `DETAIL_PROPORSIONAL` yang induknya kontrak non-proporsional. Bila tabel detail proporsional **hanya** berkunci asing ke tabel layer proporsional, **tidak ada kolom yang dapat memuat pengenal kontrak non-proporsional**. → **TIDAK DAPAT DILANGGAR** |
| **INV-33** — `BAGIAN` hanya di bawah kontrak non-proporsional | cermin INV-32. → **TIDAK DAPAT DILANGGAR** |
| **INV-34** — `CARA_PEMBUKUAN_XOL` hanya terisi pada non-proporsional | **TIDAK.** Kolom itu ada di **kepala**, dan kepala **tidak dipisah** (5b: 46 dari 47 identik). Selama kepala satu tabel, kolom itu dapat terisi pada baris berkontrak proporsional. → **tetap INDEKS UNIK** |

> **INV-34 bertahan justru karena kepala tidak dipisah** — dan ia menjadi alasan ketiga kenapa
> keputusan T-1 harus diambil lebih dulu. **Dan butir 7.9 contoh JSON menunjukkan INV-34 mungkin
> salah sejak awal:** `AccountingMode` dan `AccountingModeNonProp` **keduanya terisi** pada kontrak
> **proporsional**. Bila pola itu berulang, INV-34 sebagaimana tertulis akan **menolak kontrak yang
> sah** — dan itu diperiksa di langkah berikutnya, bukan di sini.

**Hitungan §1 yang berubah bila T-2 disahkan:** INDEKS UNIK **4 → 2**, TIDAK DAPAT DILANGGAR
**2 → 4**. Jumlah invarian tetap 63.

---

## 1b. Sapuan asal sistem — **MENYEMPITKAN, TIDAK MEMUTUSKAN**

Uji yang diusulkan dijalankan atas 708 berkas, dengan penanda **milik aturannya sendiri** — kemunculan
`pxCreateSystemID` **pertama**, yang duduk di blok `Rule-Obj-*` di kepala berkas. Disiplin yang sama
dengan mengambil `pyClassName` pertama, dan diperiksa langsung pada satu berkas sebelum dipakai.

| Sistem penulis | Berkas |
|---|---:|
| `pega` | 465 |
| `pegadevnusare2` | 206 |
| **`pegaprdnusare`** | **32** |
| `SLSHYD`, `SLSV72HYD`, `sls-envhyd83`, `sde` | 5 |

### Kenapa ia tidak memutuskan, dan saya menyatakannya alih-alih membungkusnya

> **`pxCreateSystemID` mencatat tempat aturan DIBUAT, bukan tempat ekspor DIAMBIL.** Aturan yang
> ditulis di lingkungan pengembangan lalu dipromosikan ke produksi **tetap menyebut lingkungan
> pengembangan** selamanya.

Jadi ketiga cabang yang diharapkan — "seluruhnya produksi", "ada yang lain", "tidak ada penanda" —
tidak berlaku: penandanya ada di **seluruh** berkas, dan yang dijawabnya pertanyaan yang berbeda.

### Yang tetap diperoleh, dan tiga di antaranya berguna

| | |
|---|---|
| **32 aturan dibuat LANGSUNG di `pegaprdnusare`** | itu temuan tersendiri — **penulisan aturan di lingkungan produksi**. Ia tentang cara kerja tim, bukan tentang asal ekspor, dan ia layak masuk daftar eskalasi |
| **Tidak satu pun temuan kita bersandar pada aturan yang dibuat di produksi** | `FetchQSfromMasterXOL`, `SetSpreadName`, `TreatyInSetValueInstallment`, `SetSpreadingXOL`, `SaveTreatyInDetail_Act` seluruhnya `pegadevnusare2`; `SetReinstatementPct` dan kedua `TreatyInMappingDataconvert` dibuat di `pega` pada 2020 |
| **`LinkService` sendiri dibuat di `pegadevnusare2`**, 4 September 2025 | nilai `3` karena itu **konsisten dengan setelan lingkungan bukan-produksi** — tetapi konsistensi bukan pembuktian |

**Kesimpulan: asal lingkungan sebuah ekspor adalah pertanyaan ke ORANG, bukan ke berkas.**
Dicatat sebagai **L-9**, dengan penagihnya.

---

## 8. USUL: pemisahan mengikuti DATA, dengan ambang berangka

Kedua angka 5b menunjuk arah berlawanan, dan justru itu yang membuat usul ini dapat dinilai.

| Entitas | Angka | Arah |
|---|---|---|
| `LAYER` | **13 dari 15** ruas khas satu cabang | tabel gabungan punya **13 kolom yang selalu kosong** pada setiap baris proporsional → **dipisah** |
| `VERSI_KONTRAK` | **46 dari 47** atribut identik | dua tabel **98% sama**; setiap perubahan dilakukan dua kali, **satu akan basi** → **tidak dipisah** |

> **Ambang yang diusulkan, berangka dan diterapkan ke SETIAP entitas — bukan dinilai satu per satu:**
> sebuah entitas dipisah menurut cabang bila **lebih dari separuh ruasnya hanya dipakai satu
> cabang**.

**Kasus ujinya `POTONGAN`** — dan ini penerapan T-3 sebagai alat, bukan sebagai tabrakan.
`POTONGAN` berinduk dua dan **INV-63 melarang aturannya tertulis dua kali**. Ruasnya identik di
kedua pelekatan (§14.2), jadi ambang di atas **tidak** menyuruh memisahkannya. **Bila suatu ambang
menyuruh memisahkan `POTONGAN`, ambang itu salah** — itu penguji ambangnya.

### Hasil penerapan ambang pada 28 entitas

| Dipisah | Tidak dipisah |
|---|---|
| `LAYER` **dan hanya bila pemilik proses menyetujui ambangnya** | kepala `KONTRAK` + `VERSI_KONTRAK` · 12 anak versi · `POTONGAN`, `PENYEBARAN`, `RINCIAN_PENYEBARAN`, `NILAI_PENYEBARAN` · 6 tabel acuan |
| *(`DETAIL_PROPORSIONAL`, `BAGIAN`, `PEMULIHAN_LIMIT` **sudah** terpisah sejak model disusun — pemisahan tabel tidak menambah apa pun)* | |

---

## 9. KOREKSI ANGKA §1 — milik saya yang benar, dan bedanya satu

Pemilik proses menulis **INDEKS UNIK 4 → 3**; saya menulis **4 → 2**. Keduanya tidak dapat sama-sama
benar, dan daftarnya memutuskan.

Keempat INDEKS UNIK di `SPEC-INVARIAN.md` §2: **INV-25, INV-32, INV-33, INV-34.**

Uji §2.8 menemukan **dua** yang pindah — INV-32 dan INV-33 — sementara **INV-34 bertahan** karena ia
hidup di kepala dan kepala tidak dipisah. Maka:

| | Sebelum | Sesudah |
|---|---:|---:|
| INDEKS UNIK | 4 | **2** — tersisa INV-25 dan INV-34 |
| TIDAK DAPAT DILANGGAR | 2 | **4** |
| **Total** | **63** | **63** |

**4 → 3 hanya benar bila satu invarian yang pindah.** Yang pindah dua, jadi 4 → 2.
Angka TIDAK DAPAT DILANGGAR 2 → 4 sama di kedua hitungan, dan justru ia yang menunjukkan
ketidakseimbangannya: dua masuk, maka dua harus keluar.

---

## 10. T-5 DITUTUP — satu baris, supaya tidak ditemukan lagi sebagai tabrakan

> **Butir 5.3 dan ADR-0034 adalah kalimat yang sama dilihat dari dua sisi.** 5.3 melarang JSON
> sebagai **bentuk simpan yang dibaca** — seluruh isi menjadi kolom sungguhan. ADR-0034
> mempertahankannya sebagai **arsip tanpa jalur baca**, dan INV-61 menegaskan ia **bukan sumber
> kanonik**. Sesuatu yang tidak punya jalur baca bukan bentuk simpan yang dibaca. **Keduanya
> berdiri bersama tanpa perubahan pada salah satunya.**

Baris itu juga ditulis di `SPEC-MODEL-DATA.md` dan di `STRUKTUR-TABEL-KEPALA.md` saat berkas itu
dibuat, karena pembaca yang menemukan tabrakan ini tidak akan membuka berkas angka.

---

## 11. T-1 dikerjakan atas bacaan: §10 BERDIRI

Bacaan yang dipakai, dan dasarnya temuan §1(c) berkas ini sendiri:

> **Butir 5.1 menetapkan BENTUKNYA, bukan daftar kolomnya.** Kepala kontrak berbentuk **relasional
> datar dengan kolom bernama**, seperti `TREATY_IN` — bukan satu kolom dokumen seperti
> `M_TREATY_IN`. Ia **tidak** menetapkan bahwa kolomnya dua puluh.

`TREATY_IN` tidak punya penulis langsung; ia diisi prosedur yang sama yang menulis `M_TREATY_IN`,
dalam panggilan yang sama. **Ia proyeksi yang diturunkan, bukan sumber.** Maka dua puluh kolomnya
bukan inventaris isi kontrak melainkan **ringkasan** — ruas yang dulu dipilih seseorang supaya
kontrak dapat dicari dan dilaporkan tanpa membongkar JSON.

Menurunkan model dari ringkasan berarti membuang 35 atribut yang §10 tetapkan dari sumber yang
lengkap, dan itu **berlawanan dengan maksud butir 5.3** — *"seluruh isi menjadi kolom sungguhan"*.
Seluruh isi, bukan dua puluh.

| | |
|---|---|
| **Kepala** | tetap `KONTRAK` 8 + `VERSI_KONTRAK` 47, sebagaimana §10.1 dan §10.2 |
| **20 kolom `TREATY_IN`** | menjadi **titik mulai `TABEL-DATAR-KEPALA`** — daftar ruas yang orang benar-benar pakai untuk mencari dan melaporkan, teruji bertahun-tahun. Yang ditambahkan di atasnya disebut satu per satu beserta alasannya |
| **`M_TREATY_IN`** | arsip tanpa jalur baca — ADR-0034, INV-61 |

**Bila pemilik proses menegaskan sebaliknya** — bahwa kepala benar-benar menjadi dua puluh kolom —
maka 35 atribut menuntut tempat baru, dan itu **perancangan ulang, bukan penyesuaian**. Dalam hal
itu pekerjaan berhenti dan dilaporkan.
