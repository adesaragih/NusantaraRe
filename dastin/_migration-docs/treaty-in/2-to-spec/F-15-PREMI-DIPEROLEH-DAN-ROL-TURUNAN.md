# `F-15` — `PremiumEarnedList` dan `ROLPct` **DIHITUNG aturan yang hidup**

**Tanggal:** 24 September 2026 · **Putaran penutup to-spec, butir 1**
**Nomor:** `F-15`. Sempat ditulis `F-12` — nomor itu **sudah terpakai** sejak putaran
sebelumnya (`SISA-DAN-SERAH-TERIMA.md` §7.3), dan satu nomor untuk dua temuan adalah
cacat yang sama bentuknya dengan satu nama untuk dua arti.
**Sifat:** temuan pembatal. **Entitas `NILAI_PREMI_DIPEROLEH` TIDAK dibuat.**
**Aturan yang dipakai:** *"Bila saat menerapkan Anda menemukan isinya keliru, **berhenti dan
laporkan** — jangan betulkan diam-diam."*

> ## RINGKAS
>
> Perintah putaran ini berbunyi: *"Bentuknya sudah menentukan rumahnya — `PremiumEarnedList` adalah
> `Page List` berisi `Currency` dan `Value`, golongan A menurut `P-8`."*
>
> **Bentuk menentukan BAGAIMANA sesuatu disimpan. Ia tidak menentukan APAKAH sesuatu disimpan.**
> Nasib `PremiumEarnedList` sudah diadili **DITURUNKAN** di dua artefak induk, dan putaran ini
> menemukan **rumus yang menghitungnya, hidup, di dalam ekspor**. Menjadikannya tabel anak berarti
> menyimpan besaran turunan sebagai fakta — persis yang ADR-0037 larang dan yang ADR-0041 sebut
> "dua penulis untuk satu fakta".
>
> **Dan temuan yang sama mengenai kolom yang SUDAH mendarat:** `PERSEN_ROL` juga dihitung aturan
> yang hidup. Ia sudah berdiri di §10.3 atas putusan `TDA-17`. **Tidak saya cabut sendiri** — ia
> ditandai dan ditagihkan di sini.

---

## 1. Dua artefak induk sudah mengadilinya, dan keduanya berbunyi DITURUNKAN

| Artefak | Barisnya, dikutip apa adanya |
|---|---|
| `2-to-spec/PENELUSURAN-JSON-KE-KOLOM.md` baris 357–358 | `Limits[].PremiumEarnedList[].Currency` → **DITURUNKAN** · *"turunan — tidak disimpan (ADR-0037, §4)"* · dasar **`EVIDENCED(PETA-TELUSUR-JSON@ekspor-2026-09)`** |
| `SPEC-MODEL-DATA.md` §10.3a, baris putusan `PremiumEarned` | **"turunan"** · *"kembaran skalar dari `PremiumEarnedList`, **yang sudah DITURUNKAN di peta telusur**"* |

`TDA-17-PUTUSAN.md` §3.3 mengusulkan **DISIMPAN** tanpa menyebut satu pun dari keduanya. Itu bukan
pertikaian yang boleh diselesaikan dengan memilih yang lebih baru: **yang satu bersandar pada bentuk,
yang lain pada nasib yang sudah diadili.**

---

## 2. Bukti baru — rumusnya dibaca, bukan disimpulkan

Disapu atas `alat/datar-penulis-properti.csv` (kelima bentuk penulis, kedua ekspor, 708 berkas).

### 2.1 `PremiumEarnedList` — EGNPI dikali tarif penyesuaian

`Activity/DetailCalculation.xml`, **langkah 25, HIDUP** (diperiksa `alat/langkah-hidup.py`:
33 langkah, 32 hidup, 1 mati — langkah 25 **bukan** yang mati):

```
pyStepsObjectName : .EgnpiTotalList                       <- gelungnya per mata uang EGNPI
Primary.PremiumEarnedList(<APPEND>).Currency = .Currency
Primary.PremiumEarnedList(<LAST>).Value      = @if(Local.TNOP=="tnop",
                                                   .Value * local.adjrate * Local.Prorate,
                                                   .Value * local.adjrate)
```

**`Primary` di sini adalah `LAYER`, dan itu diperiksa, bukan ditebak.** Kelas terapan aturannya —
dibaca `alat/halaman-aktivitas.py`, perkakas yang ditulis putaran ini justru karena nama halaman
relatif tidak terbaca sapuan mana pun:

```
pyClassName = ASM-FW-GISFW-Data-TreatyInLimits
```

Itu kelas yang sama dengan `TreatyIn.Limits[]` pada `4-erd-dan-tabel-datar/datar-treatyin-lama.csv`
baris 340. **Jadi `Primary.PremiumEarnedList` adalah `Limits[].PremiumEarnedList`.**

> **Bacaannya:** premi diperoleh = **EGNPI per mata uang × tarif penyesuaian**, dikali prorata bila
> jenisnya `tnop`. Keduanya sudah menjadi atribut tersimpan — `EGNPI.NILAI_EGNPI` dan
> `LAYER.PERSEN_PENYESUAIAN`. Menyimpan hasilnya berarti menyimpan besaran yang **dapat dihitung
> ulang dari dua kolom yang sudah ada**.

Penulis keduanya, dan ia tidak membatalkan yang di atas:

```
TreatyInMappingDataconvert.xml:
  TreatyIn.Limits(<LAST>).PremiumEarnedList(<APPEND>).Value = .EarnedPremium
```

Itu **penyalinan dari borang penawaran** saat penawaran dikonversi menjadi treaty — bukan masukan
baru. Nilai yang disalin itu kemudian **ditimpa** `DetailCalculation` pada jalur `adj`.

### 2.2 `ROLPct` — premi dibagi limit

`Activity/DetailCalculationROL.xml`:

```
.ROLPct = @toDecimal(@if(Local.TotalLimit==0, 9989998,
                         @divide(local.TotalPremi, Local.TotalLimit, 8))) * 100
```

`Activity/TreatyInActualUpdateValueLimits.xml`:

```
.ROLPct = @divide(Local.Value, .Limit, 4) * 100
```

**Rate on line = premi ÷ limit × 100.** `TDA-17-PUTUSAN.md` §3.2 sudah menduga ini dan menjadikannya
pertanyaan bisnis: *"apakah rate on line pernah disepakati sebagai angka, atau selalu dihitung"*.
**Pertanyaan itu terjawab dari sumber, dan jawabannya "selalu dihitung"** — sehingga mengirimkannya
ke teknik treaty menjadi uji yang jawabannya sudah terbaca, yaitu **biaya**, bukan ketelitian.

### 2.3 Kalibrasi — supaya "ada rumusnya" tidak terbaca sebagai tuduhan kosong

Sapuan yang sama dijalankan atas `Deductible2`, dan hasilnya **berbeda**:

| Ruas | Penulis yang **menghitungnya dari atribut lain** | Bacaan |
|---|---|---|
| `PremiumEarnedList` | **ada** — `DetailCalculation` langkah 25 | **turunan** |
| `ROLPct` | **ada** — `DetailCalculationROL`, `TreatyInActualUpdateValueLimits` | **turunan** |
| `Deductible2` | **tidak ada.** Yang ditemukan hanya transformasi pada halaman ringkasan (`CalculateRetroSumary`, `TreatyInDifferenceLimitsSumary`) dan penskalaan prorata (`TreatyEDMProRateCalculation`) — tidak satu pun **menurunkan nilainya dari atribut lain** | **masukan — `DEDUCTIBLE_KEDUA` BERTAHAN** |

> Kalibrasi ini yang membuat temuan ini dapat dipercaya. Kalau ketiganya kembali "ada rumusnya",
> sapuannya yang salah. Satu dari tiga kembali bersih, dan yang bersih itu **justru yang sudah
> mendarat dengan alasan paling kuat**.

---

## 3. Satu tetapan di dalam kode, dan ia menentukan bentuk kolom

```
@if(Local.TotalLimit==0, 9989998, …)
```

**`9989998` adalah kegagalan yang menyamar sebagai nilai** — persis yang ADR-0035 larang. Ketika
limitnya nol, sistem lama tidak menolak dan tidak mengosongkan; ia **menuliskan angka sembilan juta
sekian sebagai persentase**.

Dua akibatnya, dan keduanya praktis:

1. Setiap baris warisan yang berlimit nol membawa `ROL_PCT = 9989998`. Migrasi yang membacanya
   sebagai persentase **memindahkan kegagalan sebagai fakta**.
2. **Ia tidak muat di `NUMBER(11,8)`** — tipe yang kelompok `P2` usulkan, yang hanya menampung tiga
   angka di depan koma. Itu salah satu dasar `KTV-A` menyeragamkan presisi angka.

Golongan temuannya: **tetapan yang ditulis langsung di kode adalah calon** — alat cari yang sudah
berdiri di `TETAPAN-DI-KODE.md`. Butir ini **ditambahkan ke sana**.

---

## 4. Yang saya kerjakan, dan yang TIDAK

| | |
|---|---|
| **TIDAK dibuat** | entitas `NILAI_PREMI_DIPEROLEH`, berkas `ddl-usulan/`-nya, dan barisnya di `KAMUS-KOLOM.md` |
| **TIDAK dicabut** | kolom `LAYER.PERSEN_ROL` yang sudah mendarat. Mencabutnya diam-diam sama buruknya dengan memasangnya diam-diam |
| **Dikerjakan** | penandaan `PERSEN_ROL` di §10.3 yang menunjuk berkas ini; butir `9989998` di `TETAPAN-DI-KODE.md`; penagih §11.1 **dipersempit menjadi satu kolom dan diarahkan ulang** |
| **Yang memutuskan** | **pemilik proses.** Dua jalan, keduanya tertulis di §5 |
| **Yang menagih** | berkas ini, dan ketiadaan `PREMIUM_EARNED` di `KAMUS-KOLOM.md` |

---

## 5. Dua jalan, beserta ongkos salahnya — supaya keputusannya satu kalimat

### Jalan A — ikuti nasib yang sudah diadili: **keduanya turunan, tidak disimpan**

| | |
|---|---|
| **Yang terjadi** | `NILAI_PREMI_DIPEROLEH` tidak pernah berdiri; `LAYER.PERSEN_ROL` **dicabut** sebelum data dimuat; keduanya dihitung saat dibaca, dari `NILAI_EGNPI` × `PERSEN_PENYESUAIAN`, dan premi ÷ `LIMIT` |
| **Arah dampak bila salah** | bila ternyata pernah ada yang **menyunting** hasilnya dengan tangan, suntingan itu hilang. **Dugaan itu tidak dapat disahkan data** — suntingan yang ditimpa rumus tidak meninggalkan jejak, jadi "tidak pernah disunting" dan "selalu ditimpa" menghasilkan data yang persis sama |
| **Yang menyempitkannya** | **`Uji AQ`** (baru, §6) — sebaran `ROL_PCT` terhadap premi ÷ limit yang dihitung ulang dari kolom lain. Baris yang **menyimpang** adalah baris yang rumusnya tidak menjelaskan |

### Jalan B — pertahankan putusan `TDA-17`: **keduanya disimpan sebagai potret beku**

| | |
|---|---|
| **Yang terjadi** | `NILAI_PREMI_DIPEROLEH` berdiri sebagai tabel anak golongan A; `PERSEN_ROL` bertahan. Keduanya diberi label **potret**, bukan fakta — sekeluarga dengan ADR-0036 "beku saat disetujui, hitung saat dibaca" |
| **Dasar yang sah untuknya, dan ia nyata** | ADR-0036 membekukan angka yang disetujui. Bila premi diperoleh **ikut dibekukan saat persetujuan**, menyimpannya bukan menyimpan turunan melainkan menyimpan **apa yang disetujui** — dan itu berbeda |
| **Arah dampak bila salah** | dua penulis untuk satu fakta (ADR-0041). Ketika `NILAI_EGNPI` atau `PERSEN_PENYESUAIAN` berubah pada versi berikutnya, kolom potretnya **tidak ikut berubah** — dan tidak ada yang memberi tahu pembacanya mana yang berlaku |
| **Syarat pembalikan** | ia **harus** memakai nama yang menyatakan dirinya potret, bukan `PREMI_DIPEROLEH` polos. Nama yang berarti dua hal adalah cacat yang ditanam dengan tangan sendiri |

> **Rekomendasi saya: Jalan A**, dan sebabnya satu kalimat — dua artefak induk sudah mengadilinya
> DITURUNKAN dengan dasar `EVIDENCED`, dan putaran ini menambahkan rumusnya. Membalikkan putusan
> ber-`EVIDENCED` menuntut bukti yang lebih kuat daripada bentuk `Page List`.
>
> **Dan bila Jalan B yang dipilih, ia tetap sah** — asalkan dasarnya ditulis **ADR-0036 (potret
> beku)**, bukan *"bentuknya daftar per mata uang"*. Dasar yang salah untuk kesimpulan yang benar
> tetap membuat pembaca berikutnya menarik kesimpulan yang salah dari dasar itu.

---

## 6. `Uji AQ` — uji data baru, dan sifat pembuktiannya disebut di kepalanya

> **Uji ini dapat MEMATAHKAN bacaan "selalu dihitung"; ia TIDAK dapat mengesahkannya.** Bila seluruh
> baris cocok dengan rumusnya, itu konsisten dengan "tidak pernah disunting" **dan** dengan "selalu
> ditimpa rumus" — dua kemungkinan yang menghasilkan data yang sama.

| | |
|---|---|
| **Pertanyaannya** | untuk tiap baris `TREATYINDETAIL` berkolom `ROL_PCT`, berapa selisihnya terhadap premi ÷ limit yang dihitung ulang dari kolom di baris yang sama? |
| **Yang memisahkan** | baris yang **menyimpang lebih dari pembulatan**. Baris semacam itu berarti ada penulis di luar rumus yang kita baca |
| **Dan yang wajib dilaporkan terpisah** | berapa baris bernilai **tepat `9989998`**. Ia bukan penyimpangan, ia kegagalan yang tercatat — dan bila ikut masuk hitungan pokok, ia akan terbaca sebagai sifat aturan umum |
| **Dua kontrak istimewa** | `ID` yang diperlakukan berbeda di enam belas tempat **dikeluarkan dari hitungan pokok dan dilaporkan di baris tersendiri**, sesuai aturan yang berlaku untuk seluruh uji di berkas permintaan DBA |

**Ditambahkan ke:** `PERMINTAAN-DBA-1-UJI-A-SAMPAI-H.sql` pada putaran pengiriman berikutnya —
**bukan sekarang**, sebab berkas itu sudah diserahkan dan menyisipkan uji ke dalam berkas terkirim
membuat dua salinan berbeda beredar.

---

## 7. Apa yang ini TIDAK katakan

- Ia **tidak** membatalkan `TDA-17` seluruhnya. `DEDUCTIBLE_KEDUA` bertahan, dan kalibrasinya justru
  yang membuat temuan ini dapat dipercaya.
- Ia **tidak** menutup `L-8`. Ketiga kolom tetap ditemukan lewat titik buta, dan **340 properti
  masih belum diperiksa siapa pun** (`M-4`). Yang berubah hanya **nasib** dua di antaranya, bukan
  sebab hilangnya.
- Ia **tidak** menyatakan `DetailCalculation` satu-satunya penulis. Batas sapuannya tetap `L-10`:
  **21 jenis aturan nol kemunculan** di ekspor, dan `Declare Expression` — yang justru jenis aturan
  yang paling mungkin menghitung sebuah turunan — adalah salah satunya. Bila ada ekspresi deklaratif
  yang menghitung `PremiumEarned`, sapuan ini **tidak dapat melihatnya**, dan itu menguatkan
  bacaan "turunan", bukan melemahkannya.
