# Bahan spesifikasi **pemuatan** ke tabel flat

> **24 September 2026 · direvisi 25 September 2026.** Bahan, **bukan** spec. `/to-spec`
> ber-`disable-model-invocation`; berkas ini menyiapkan bahannya lalu berhenti.
>
> **Revisi 25 September** memasukkan **K-070** (migrasi dua fase), **K-071** (J-1…J-10),
> **K-072** (J-11…J-14) dan **K-073** (Jalan B disetujui). Yang berubah paling besar: **§5** —
> rekonstruksi `ROW_UID` **tidak lagi diperlukan**; dan **§6** — kelima butir yang dulu terbuka di
> sana **sudah tertutup**, tersisa dua.
>
> **Lingkup:** memindahkan isi satu penawaran Fac In ke **78 tabel flat**
> (`DDL-tabel-flat-draf.sql`). **Di luar lingkup:** jalur tulis produksi ke `FACINPRODUCTION` /
> `FACOUTPRODUCTION` — keduanya sudah flat dan sudah ada (lembar `Tabel Sudah Ada`).

---

## 0. ⛔ Temuan yang mengubah bentuk loader — masukan adalah **JSON**, bukan XML

`[terverifikasi]` `DDL\JSON_POLIS.txt`:

```
"DATA_JSON" CLOB
CONSTRAINT "VALID_JSON" CHECK (DATA_JSON IS JSON(STRICT)) DISABLE
CONSTRAINT "PK_JSON_POLIS" PRIMARY KEY ("IDPEGA")
```

Tetapi **seluruh 115 contoh** di `DDL\CONTOH\` berbentuk **XML** (`<pagedata>`), dan **seluruh
rancangan tabel flat diturunkan dari XML itu**. Loader produksi akan membaca **JSON**.

### Apakah keduanya membawa data yang sama? — **ya**

`[terverifikasi]` Dua kasus kebetulan ada dalam **kedua** bentuk. Jalur daun terisi dibandingkan
sebagai himpunan:

| Kasus | JSON | XML | Sama | Hanya JSON | Hanya XML |
| --- | ---: | ---: | ---: | ---: | ---: |
| `NB-181231` *(sesudah normalisasi)* | 621 | 622 | **621** | **0** | **1** |
| `NB-176005` *(sesudah normalisasi)* | 224 | 240 | 214 | 10 | 26 |

⛔ **`NB-181231` cocok 621 dari 621.** Satu-satunya jalur XML-saja adalah `pzStatus` — metadata Pega.

⚠️ `NB-176005` berselisih lebih banyak, tetapi selisihnya **bukan soal bentuk**: yang XML-saja
hampir seluruhnya tag `px*`/`py*` (dibuang kontrak 23 tag), sisanya medan bisnis yang terisi di satu
ekspor dan kosong di ekspor lain. `[dugaan]` **beda waktu ekspor**, bukan beda format. ⚠️ Berkas ini
juga salah satu dari **dua pasang duplikat byte-identik**, jadi bukan saksi terbaik.

### ⛔ Satu perbedaan bentuk yang sistematis dan WAJIB ditangani

| | XML (contoh) | JSON (produksi) |
| --- | --- | --- |
| Bentuk | `CurrencyList/<USD>/Name` | `CurrencyList` adalah **array**; `Name` di dalam unsurnya |
| Kode mata uang | nama elemen **dan** medan `Name` | **hanya** medan `Name` |

`[terverifikasi]` di JSON `NB-181231`, `CurrencyList` bertipe **`Object[]`** dengan **1 unsur**,
ber-medan `Name`.

📌 **Ini menutup perselisihan lama** tentang bentuk `CurrencyList`: keterangan "anaknya `rowdata`,
bukan elemen bernama kode" dan pengukuran "elemen bernama kode" **dua-duanya benar** — untuk
**artefak yang berbeda**. Di ekspor XML ia elemen bernama kode; di JSON produksi ia unsur array.

⛔ **Konsekuensi mengikat:** seluruh jalur di lembar `Jalur Sumber` diturunkan dari **XML**. Loader
JSON **harus membuang ruas kode mata uang** dari jalur itu. Pembawa mata uang tetap **`Name`** di
kedua bentuk, sehingga kolom `CURRENCY_CODE` (K-063 kelompok b) **tetap sah**.

✅ **DITUTUP K-066 (24 September).** Hanya 2 kasus tersedia dalam dua bentuk, keduanya NB — tetapi
work owner menetapkan **sepadan untuk seluruh COB**. Tidak perlu contoh berpasangan tambahan.

⛔ **Yang tetap berlaku:** perbedaan **bentuk** di atas nyata dan wajib ditangani loader. Yang
ditutup adalah pertanyaan **kesepadanan isi**, bukan **kesamaan bentuk**.

---

## 1. Ukuran pekerjaan

`[terverifikasi]` atas 115 contoh:

| Ukuran | Nilai |
| --- | ---: |
| Simpul calon baris flat | **22.622** |
| — bernama `rowdata` | **20.937** |
| — bernama kode mata uang | **1.685** |
| — tumpang tindih | **0** |
| Rata-rata per berkas | **~197** |
| Terbanyak dalam satu berkas | **13.249** |
| **Kedalaman `rowdata` bersarang maksimum** | **8** |

⚠️ **Angka ini membawa ketidakpastian.** Metode kedua (`SelectNodes('//rowdata')`) menghasilkan
**22.139** simpul `rowdata`, bukan 20.937, dan selisih **1.202** itu **tidak dapat saya rekonsiliasi**
— padahal pemindaian nama menunjukkan hanya **satu** ejaan (`rowdata`). Angka di atas memakai
pemindaian nama, yang **jumlahnya konsisten** (20.937 + 1.685 = 22.622 persis). Untuk penentuan
ukuran ini memadai; **jangan dipakai sebagai angka rekonsiliasi**.

Wadah berulang terbesar — inilah yang menentukan biaya loader:

| Wadah | Baris |
| --- | ---: |
| `DeductibleList` | 6.557 |
| `CoverageList` | 4.943 |
| `SpreadingList` | 4.526 |
| `ViewSuggest` | 936 |
| `PropertyItemList` | 934 |

---

## 2. Yang HARUS dikecualikan

`[terverifikasi]` keputusan yang sudah ada, bukan usulan baru:

| Dibuang | Dasar | Catatan |
| --- | --- | --- |
| Seluruh cabang **`OldData`** | **V-19** | `[terverifikasi]` pada satu sapuan 9 medan saja, **326 kemunculan** dibuang |
| `OldOfferedPayment` | **V-44** | muncul di 2 berkas EDM |
| Cabang `FacRetro` tanpa `List` | **V-25** | salinan kerja; nol `rowdata` |
| `Parameters`, `OutGoList` | **V-34** | parameter layar |
| `TotalTSIList`, `TotalTSIPremiGrossList`, `TotalTSIPremiSpreadRNM` | **V-28** | ⚠️ **butir 8 MASIH TERBUKA** — lihat §6 |
| Medan berakhiran `Old` dan berpola `EDMOld*` | **V-28b**, **V-45** | ⚠️ 4 berkas **NB** mengisi `EDMOld*` |
| Tag `px*`/`pz*` kontrak 23 tag | **K-042/K-043** | metadata ekspor |

⛔ **`OldData` bukan versi sebelumnya.** V-43a: ia potret salinan kerja **sesudah** baris baru
dibentuk. Selisih endorsement dihitung terhadap baris versi sebelumnya, bukan terhadap `OldData`.

---

## 3. Kolom sistem — dari mana isinya

| Kolom | Isi | Status |
| --- | --- | --- |
| `ID` | surrogate, dibangkitkan | jelas |
| `IDPEGA` | kunci kerja Pega | jelas |
| `PARENT_ID` | `ID` baris induk | jelas |
| `PARENT_TABLE` | nama tabel induk | ⛔ **wajib** untuk 12 tabel berinduk ganda |
| `SRC_PATH` | jalur sumber | hanya tabel berjalur ganda (**V-23**) |
| `COB_GROUP` | hasil aturan **V-16** | lihat §4 |
| ~~`LINI`~~ | ~~penanda facultative vs treaty~~ | ⛔ **DIHAPUS — K-069.** Penanda lini **tidak ada di proyek ini**; kolomnya sudah dicabut dari rancangan |
| `SEQ_NO` | **posisi baris saat muat** (V-41) | jelas |
| `ROW_UID` | identitas baris | ✅ **dibangkitkan, bukan direkonstruksi** (K-073) — §5 |
| `PROD_KE` | nomor versi polis | ⛔ sumber `VARCHAR2(5)` → sasaran **`NUMBER`** (K-071 J-3). Nilai bukan angka wajib **`panic`** |
| `OLD_POLIS_ID` | penunjuk polis yang di-renew | ⛔ **hanya terisi saat RENEWAL** (K-071 J-5). NB dan EDM `NULL` |

### ⛔ Dua belas tabel yang keutuhannya tidak dijaga basis data

`[terverifikasi]` dari `Daftar Relasi`: 64 pasangan (anak, kunci tamu) berinduk **tunggal** → punya
`FOREIGN KEY`. **12 berinduk ganda** → **tidak bisa**. Yang terparah `T_COVERAGELIST.PARENT_ID`,
menunjuk **5** tabel berbeda.

`[terverifikasi]` **12 dari 12** punya kolom `PARENT_TABLE`. Loader **wajib mengisinya** — tanpa itu
baris anak menjadi yatim yang tidak dapat ditelusuri, dan **tidak ada constraint yang akan
mengeluh**.

---

## 4. `COB_GROUP` — aturan berjenjang V-16

Dibaca **dari atas**, berhenti di yang pertama cocok:

1. ada `LocationList/Property/PropertyItemList` → **FIRE**
2. ada `AnekaList` di bawah `LocationList` bisnis → **Aneka**
3. ada `VehicleList` → **MBUCar**
4. ada `CargoList` → **MarineCargo**
5. sisanya → **Life** atau **PA**

⛔ **Langkah 5 satu-satunya yang bergantung pada NILAI, bukan bentuk** (V-16a): Life dan PA tidak
dapat dibedakan secara struktural; pembedanya `QuotationData.BusinessType`.

`[terverifikasi]` `BusinessType` terbaca **128 kemunculan**, **18 nilai berbeda** — cocok dengan
"18 BusinessType" pada V-16. Nilainya: `Aneka` `AviationHull` `Bonding` `BondingKBG` `CustomBond`
`ElectronicEquipment` `FireStyle1` `FireStyle2` `GolfInsurance` `HE` `LandRig` `Liability` `Life`
`MarineCargo` `MarineHull` `MBD` `MBUCar` `PA`.

⚠️ 128 kemunculan atas 115 berkas — **sebagian berkas memuatnya lebih dari sekali**. Loader harus
menetapkan **mana yang mengikat** bila lebih dari satu.

---

## 5. ✅ `ROW_UID` — **bukan lagi bagian tersulit** (K-070 · K-073)

**V-42** menuntut tiap baris tabel berulang punya identitas yang **tetap sama di seluruh versi
polis**. Alasannya mengikat: `[terverifikasi]` di ekspor Pega **tidak ada satu pun identitas baris
yang stabil** — `pxListSubscript` hanya posisi, `OBJECT_NO` juga sekadar 1..N.

⛔ **Yang berubah: tidak ada lagi UID yang perlu DIREKONSTRUKSI.**

`BAHAN` versi sebelumnya menandai rekonstruksi UID lintas versi sebagai bagian tersulit dari seluruh
sisa pekerjaan. **K-070 + K-073 membubarkannya:**

| Fase | Perlakuan `ROW_UID` |
| :-: | --- |
| **1** — versi terakhir | **dibangkitkan** saat baris lahir · ⛔ **bersifat sementara** |
| **2** — seluruh generasi | ⛔ tabel **dikosongkan**; UID **dibangkitkan sekali lagi**, atas riwayat lengkap, dalam **satu jalan** |
| sesudah fase 2 | **beku** — tidak ada muat ulang lagi |

Karena fase 2 memandang riwayat lengkap sekaligus, ia **tidak perlu menjodohkan** ke UID mana pun
yang sudah ada. Yang tersisa hanya **membangkitkan**, dan itu pekerjaan biasa.

### ⛔ Satu larangan yang menggantikan seluruh kesulitan lama

**`ROW_UID` fase 1 tidak boleh dijadikan sandaran apa pun di luar tabel flat** (K-073 konsekuensi 1):
tidak diekspor, tidak dikirim ke tim lain, tidak ditampilkan sebagai identitas, tidak dipakai sebagai
kunci integrasi. Ia **akan berganti**. Melanggarnya membuat muat ulang terlihat keluar, dan Jalan B
kehilangan sifat dapat-dibatalkannya.

### Yang tetap berlaku dan memudahkan pembangkitan

- **V-43** — baris baru **selalu ditambahkan di belakang**; baris tidak pernah dihapus fisik,
  melainkan **dinolkan**. `[terverifikasi]` di `EDM-13470`: 126 lokasi lama tetap di posisi 1..126,
  55 lokasi baru menempati 127..181, **nol** disisipkan di tengah.
- **V-40** — penjodohan dua tingkat: `OBJECT_NO` di tingkat objek, **posisi baris** di dalam objek.
- **V-40b** — `PROPERTY_ITEM_NO` **gagal** sebagai kunci item (kosong pada kedua item di `EDM-13240`).

⚠️ **Bunyi V-38 dan V-40 masih perlu ditulis ulang** — keduanya menyebut pencarian `PROD_KE`
tertinggi sebagai cara menaut versi, yang digantikan `OLD_POLIS_ID` untuk renewal (K-068, dipersempit
K-071 J-5) dan `POLICY_NO` + `PROD_KE` untuk endorsement (K-072 J-11).

---

## 5a. Aturan memilih **versi terakhir** untuk fase 1 — K-071 · K-072

```
KELOMPOKKAN menurut  POLICY_NO
URUTKAN      menurut  TGL_INPUT DESC        <- penentu
PERIKSA SILANG dengan PRODKE tertinggi      <- wajib sepakat
```

**Mengapa `TGL_INPUT` yang jadi penentu, bukan `PRODKE`:** `[terverifikasi]` `DDL\JSON_POLIS.txt`
baris 7 — `"PRODKE" VARCHAR2(5)`, **teks, bukan angka**. Mengurutkan teks membuat `'9'` terbaca lebih
besar daripada `'10'`. `TGL_INPUT` bertipe `DATE` dan bebas dari jebakan itu. Sistem lama sendiri
`[terverifikasi]` memakai `ORDER BY TGL_INPUT DESC` pada `GetProdKeOldData_SQL`.

⛔ **Work owner menyatakan keduanya sama-sama tertinggi (K-071 J-2).** Karena itu keduanya saling
menguji: bila `TGL_INPUT` terbaru dan `PRODKE` tertinggi **menunjuk baris berbeda**, anggapan itu
patah → **`panic`**, bukan diam-diam memilih salah satu.

### ⚠️ Batas pengetahuan — aturan pengelompokan tidak dapat diuji dari korpus

`[terverifikasi]` 115 berkas contoh memuat **113 `PolicyNo` unik**; hanya **2** nilai berulang,
masing-masing tepat **2 kali** — dan itu **persis** dua pasang duplikat byte-identik (V-45). ⛔ **Nol
pasangan generasi sejati** ada di contoh. Aturan di atas bersandar pada jawaban work owner dan bentuk
skema, **bukan** pengukuran.

⛔ **`POLICY_MASTER_NUMBER` bukan kunci pengelompokan** — `[terverifikasi]` terisi **1 dari 115** dan
pada contoh itu **berbeda** dari `PolicyNo`. Artinya `belum terverifikasi` (**J-17**).

---

## 6. Yang masih terbuka dan menyentuh loader

✅ **Kelima butir yang dulu berdiri di sini SUDAH TERTUTUP.** Dicatat agar tidak dibuka lagi:

| Butir | Ditutup oleh | Dampaknya ke loader sekarang |
| --- | --- | --- |
| ~~**8** tabel `Total*`~~ | **K-067** | nilainya informasi saja, **tidak dihitung ulang** → tidak ada yang dibandingkan |
| ~~**5b** `OLD_POLIS_ID`~~ | **K-068**, dipersempit **K-071 J-5** | ⛔ **hanya diisi saat RENEWAL**; NB dan EDM `NULL` |
| ~~**7b** `IsCedingConfirm`~~ | **K-069** | punya **kolom sendiri**, bukan `POSISI` |
| ~~**1c** nama `LINI`~~ | **K-069** | ⛔ **di luar proyek ini** — kolomnya dicabut |
| ~~**10** mata uang `NOT NULL`~~ | **K-069** | **24 kolom** `DEFAULT 'UNKNOWN' NOT NULL` → loader **tidak boleh** menulis mata uang kosong; yang kosong masuk sebagai `UNKNOWN` dan **bersuara** |
| ~~fan-out mata uang~~ | diukur 24-09 | 13.666 baris, **432** berkode, **seluruhnya tepat satu**, **nol** lebih dari satu → satu kolom `CURRENCY_CODE` memadai; loader **tidak perlu memecah baris** |

### ⬜ Yang benar-benar masih terbuka — **dua**

| # | Butir | Dampak ke loader |
| :-: | --- | --- |
| **J-16** | Penyaring **"endorsement sudah jadi"** untuk fase 1. ⛔ `STS_KONVERSI` **bukan** jawabannya — `[keterangan work owner]` kolom itu menandai **konversi ke tim lain**, bukan selesainya endorsement | ⚠️ tidak menentukan hasil akhir (fase 2 memuat ulang), **tetapi** menguji atas generasi yang salah menghasilkan pengujian menyesatkan |
| **J-17** | Arti `POLICY_MASTER_NUMBER` / `POLICY_MASTER_ID_PEGA` | tidak memblokir — keduanya terisi 1 dari 115 |

---

## 7. Format nilai

`[terverifikasi]` sapuan seluruh nilai daun di 115 contoh XML:

| Pola | Jumlah |
| --- | ---: |
| bulat (`123`) | **185.389** |
| **titik desimal** (`1.23`) | **58.124** |
| koma desimal (`1,23`) | **71** |
| ribuan titik (`1.234,56`) | **0** |
| ribuan koma (`1,234.56`) | **0** |

⛔ **Di korpus ini titik adalah pemisah desimal, dan tidak ada pemisah ribuan sama sekali.** Koma
desimal muncul hanya **71 kali** — tangani sebagai **pengecualian**, jangan sebagai bentuk baku.

⚠️ **Sumber lain berbeda dan jangan dicampur.** `FACINOFFER.RATE` (`VARCHAR2`) dan
`TABLEOFLIMIT.PCTLIMIT` `[terverifikasi]` memakai **koma desimal** (`0,0244`, `70,000`). Itu **kolom
Oracle**, bukan isi `DATA_JSON`. Parser loader dan parser tabel lookup **bukan parser yang sama**.

📌 Presisi dan urutan operasi **tidak diubah** (ADR-0005, ADR-0001): nilai masuk apa adanya; loader
**tidak membulatkan**.

---

## 8. Urutan muat

### 8a. Urutan **fase** — K-070 · K-073

| | Fase 1 | Fase 2 |
| --- | --- | --- |
| Populasi | **versi terakhir** tiap polis (§5a) | ⛔ **seluruh generasi**, satu jalan |
| Tabel flat | dimuat dari kosong | ⛔ **dikosongkan lebih dulu**, lalu dimuat penuh |
| `ROW_UID` | dibangkitkan, **sementara** | dibangkitkan **sekali**, lalu **beku** |
| `OLD_POLIS_ID` | `NULL` seluruhnya *(belum ada generasi sebelumnya untuk ditunjuk)* | terisi **hanya untuk renewal** |
| `PROD_KE` | nilai asli apa adanya — ⛔ **tidak dinomori ulang** | idem |
| Rekonsiliasi ADR-0001 | **tidak** dijalankan sebagai putusan akhir | ⛔ **di sinilah** nol selisih ditegakkan |

⛔ **Syarat yang menjaga bentuk ini** — enam butir di **K-072 J-12**. Yang paling mudah dilanggar
tanpa sadar: **Oracle lama wajib tetap terbaca sampai fase 2 selesai** (syarat 2, dan itu syarat
K-070 secara umum, bukan khusus Jalan B).

### 8b. Urutan **tabel** di dalam satu jalan

Diturunkan dari `Daftar Relasi`, bukan dari selera:

1. `T_WORK_POLIS` → `T_GENERAL_POLIS` (**SHARED PK** lewat `IDPEGA`, 1:1)
2. anak langsung `T_GENERAL_POLIS`
3. turun mengikuti pohon, **maksimum 8 tingkat**
4. 12 tabel berinduk ganda: isi `PARENT_ID` **dan** `PARENT_TABLE` dalam transaksi yang sama

⚠️ `T_FR_FACOFFERLIST`, `T_FR_OBJECT`, `T_FR_POLICY` adalah **wadah murni** — `[terverifikasi]` nol
medan skalar terisi. Barisnya tetap dibuat (V-27), berisi kolom sistem saja, karena anak-anaknya
menggantung padanya.

---

## 9. Yang TIDAK dapat saya verifikasi

| # | Klaim | Sebab |
| :-: | --- | --- |
| ~~1~~ | ✅ **DITUTUP K-066** — work owner: JSON dan XML **sepadan untuk seluruh COB** | — |
| 2 | Selisih 1.202 antara dua cara menghitung `rowdata` | tidak terekonsiliasi; angka §1 memakai pemindaian nama yang konsisten |
| 3 | Apakah `NB-176005` berselisih karena **beda waktu ekspor** | `[dugaan]`; berkas itu juga salah satu pasangan duplikat |
| ~~4~~ | ✅ **GUGUR — K-065.** Diagram skema rumah tidak dipakai; risiko irisan kolom diterima sadar | — |
| ~~5~~ | ✅ **DITUTUP K-070 + K-073.** Migrasi **dua fase**: fase 1 versi terakhir, fase 2 seluruh generasi dengan **muat ulang dari nol** | — |
| ~~6~~ | ✅ **DITUTUP K-066** — work owner: **struktur RNW = NB**. Sel RNW tetap berbunyi "warisan NB": diwarisi, bukan diukur | — |
| 7 | Apakah loader membaca `DATA_JSON` langsung atau lewat ekspor | belum ditetapkan |

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan. Seluruh skrip yang menyapu `DDL\CONTOH\`
hanya mengeluarkan hitungan, nama tag, nama kelas Pega, dan kode mata uang.*
