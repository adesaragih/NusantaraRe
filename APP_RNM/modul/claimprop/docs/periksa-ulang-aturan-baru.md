# Periksa ulang dengan tujuh aturan baca baru — Claim Prop

**Tanggal:** 2026-09-19 · **Korpus:** `D:\XML\RNM_BRD\Claim Prop\` (READ-ONLY)
**Sumber aturan:** tujuh aturan di §A brief, seluruhnya `[terverifikasi]` atau `[data work owner]`
dari modul lain — **dipakai apa adanya, tidak diuji ulang**.

> ⛔ **Berkas ini menambal, bukan membongkar.** Ronde 1 dan ronde 2 Claim Prop dikerjakan dengan
> baik; ketujuh aturan ini semuanya ditemukan **sesudah** ronde 2 ditutup.
>
> ⛔ **Tidak satu pun berkas lama disunting** — `grilling-ronde-1.md`, `grilling-ronde-2.md`,
> `spec.md`, dan keenam belas tiket **utuh**.
>
> ⛔ **Berkas ini tidak menyatakan butir tertutup maupun terbuka** — itu keputusan work owner.

---

## 0. Ringkasan

| # | Hasil | Status |
| --- | --- | --- |
| **§B1** | Keluarga gerbang **pertama** | **1428 baris** · dua cara hitung **selisih NOL** |
| **§B2** | Keluarga gerbang **kedua** | ⭐ **1365 baris** — **belum pernah disensus di modul ini** |
| **§B3** | Lompatan **menggantung** | **1**, dan ia **bergerbang mati** |
| **§B4** | Baris kode **5** *(Skip Whens = OR)* | ⚠️ **7 baris**, di **2 activity** |
| **§B5** | Langkah berulang | **224** — `EMBEDDED` 192 · `REPEAT` 32 |
| **§B6** | Langkah `Java` | **29** · nol sentuh basis data · **9 menangani gagal** |
| **§B7** | `Obj-*` tanpa saringan | ⭐ **NOL cacat** |
| **§B8** | Parameter di `pyXMLSignature` | ⚠️ **35 belum pernah tercatat** |
| **§C** | Keluarga gerbang Section | ronde 2 memakai **sumbu berbeda** — keduanya benar |
| **§D** | Kesimpulan **TERBALIK** | ⭐ **NOL** |
| **§D** | Kesimpulan **GOYAH** | ⚠️ **3** |

⭐ **Kalimat penutup lebih dulu: tidak ada satu pun kesimpulan ronde 1 atau ronde 2 yang terbalik.**
Yang ada tiga yang **goyah** — dua di antaranya karena ronde 2 menggabung dua hal yang kini
terbukti berbeda, dan satu karena ada wadah yang belum pernah dibaca.

---

## §A — Ketujuh aturan, dipakai apa adanya

Ditulis ulang di sini supaya berkas ini berdiri sendiri. ⛔ **Tidak diuji ulang.**

**A1 — Peta kode arah lengkap** `[data work owner, dari layar Pega]`

| Kode | Nama di Pega | Arti |
| --- | --- | --- |
| **1** | Jump to Later Step | lompat ke langkah bertanda; sasarannya di kolom parameter |
| **2** | Continue Whens | periksa baris berikutnya; habis baris → **langkah jalan** |
| **3** | Skip Step | langkah **tidak** jalan |
| **4** | Exit Iteration | putus perulangan |
| **5** | Skip Whens | **berhenti memeriksa**, langkah **langsung jalan** |
| **6** | Exit Activity | hentikan activity |
| *(kosong)* | — | tidak ada tindakan dipilih |

⭐ **Ronde 1 mencatat arti enum ini sebagai butir terbuka P13** — *"Arti enum
`pyStepsPreCondParamsWhenTrue/False`, pemilik: pemilik export Pega"*. **P13 kini terjawab.**

**A2 — Baris syarat adalah RANTAI**, dibaca dari atas. *Continue Whens* + *Skip Step* menyusun
**AND**; satu *Skip Whens* mengubahnya jadi **OR**.

**A3 — Keluarga gerbang kedua**: `pyStepsTransParams`, diuji **sesudah** langkah jalan, terhadap
hasilnya. Terpisah dari keluarga pertama.

**A4 — Keluarga ketiga**: `pyStepsRepeatDef`, penanda langkah berulang, bernilai `EMBEDDED` atau
`REPEAT`.

**A5 — Saringan `Obj-Browse`/`Obj-List`** ada di `pyParamArray/rowdata` sebagai baris
`Select`/`Field`/`Condition`/`Value`, **bukan tag tunggal**. ⚠️ `pyParamArray` **bermuka dua**.

**A6 — Parameter activity** dideklarasikan di `pyXMLSignature` milik rule-nya.

**A7 — Langkah `Java`** dikenali dari `pyStepsJavaSource` **terisi**, bukan dari nama metodenya.
⚠️ Ekspor Pega tidak konsisten kapitalnya — **setiap penyaringan harus tidak peka huruf**.

---

## §B — Penerapan ke seluruh 112 Activity `Claim Prop`

### B1 — Keluarga pertama `pyStepsPreCondParams` `[terverifikasi]`

Dihitung **dua cara**: menurun lewat pohon langkah, dan lewat jangkar kelas
`Embed-ActivityPreConditions` (tidak peka huruf).

| Cara | Hasil |
| --- | --- |
| Turun pohon langkah | **1428** |
| Jangkar kelas | **1428** |
| **Selisih** | ⭐ **NOL** |

**Sebaran lengkap:**

| Kode | Jumlah | | Kode | Jumlah |
| --- | --- | --- | --- | --- |
| `2/2` | **815** | | `6/2` | 9 |
| `2/3` | **450** | | **`5/2`** | **7** |
| `3/2` | **56** | | `_/_` | 4 |
| `_/3` | 53 | | `6/_` | 2 |
| `3/_` | 21 | | `_/2` | 1 |
| `2/6` | 10 | | | |

⭐ **Nol baris berkode 1. Nol baris berkode 4.** Di keluarga pertama modul ini, satu-satunya kode
yang dulu tidak diketahui artinya adalah **5**, dan jumlahnya **7**.

#### Perbandingan dengan angka ronde 2

Ronde 2 menulis: *"Sensus baris syarat pada gerbang AKTIF: **380 normal · 44 TERBALIK · 89 memakai
kode 5/6**."*

Dihitung ulang pada dasar yang sama — saklar gerbang bernilai `true`, langkah tidak di-remark,
syarat terisi:

| Ukuran | Ronde 2 | Hitungan saya | |
| --- | --- | --- | --- |
| Baris pada gerbang aktif | — | **482** | |
| **TERBALIK `3/2`** | **44** | **44** | ✅ **cocok persis** |
| NORMAL `2/3` | 380 | **363** | selisih 17 |
| memakai kode **5 atau 6** | **89** | **23** | ⚠️ selisih 66 |

⛔ **Selisih 17 dan 66 tidak saya kejar agar cocok.** Definisi *"gerbang aktif"* milik ronde 2 tidak
tertulis persis, jadi dasarnya bisa berbeda. Yang penting bukan totalnya melainkan **pemisahannya**
— lihat §D GOYAH 1.

### B2 — Keluarga kedua `pyStepsTransParams` ⭐ **belum pernah disensus di modul ini**

`[terverifikasi]`

| Ukuran | Jumlah |
| --- | --- |
| **TOTAL BARIS** | **1365** |
| Langkah seluruhnya | **1366** |
| Langkah punya **kedua** keluarga | **1365** |

**Sebaran:**

| Kode | Jumlah | Arti |
| --- | --- | --- |
| `2/2` | **1350** | polos — tidak ada tindakan khusus |
| `_/_` | 8 | kedua dropdown kosong |
| `6/2` · `6/_` | 2 · 2 | **Exit Activity** |
| **`1/_`** | **1** | ⚠️ **Jump to Later Step** |
| **`1/2`** | **1** | ⚠️ **Jump to Later Step** |
| `2/6` | 1 | Exit Activity bila salah |

⭐ **Hanya 15 dari 1365 baris membawa tindakan khusus.** Sisanya polos.

### B3 — Baris berkode **1** dan sasarannya `[terverifikasi]`

| Activity · step | Keluarga | Sasaran | Saklar | Tandanya ada? |
| --- | --- | --- | --- | --- |
| `SaveAdjusterConsultant_Act` **4** | kedua | **`FAIL`** | **`true`** — **HIDUP** | ✅ **langkah 8** |
| `CountValueADJTreaty_Act` **17** | kedua | **`Count`** | **`false`** — **mati** | ⛔ **MENGGANTUNG** |

⭐ **Satu lompatan menggantung, dan ia bergerbang mati.** Pola yang sama dengan modul lain:
**lompatan yang rusak adalah lompatan yang dimatikan**. Satu-satunya lompatan **hidup** menemukan
tandanya.

### B4 — Baris berkode **4** dan **5** `[terverifikasi]`

**Kode 4 (*Exit Iteration*): NOL** di kedua keluarga.

**Kode 5 (*Skip Whens* = OR): 7 baris, di DUA activity** — dan inilah temuan terpenting §B.

#### `CopyOldataCurr_act` — enam baris, empat langkah

⭐ Keempatnya membentuk **rantai OR** yang **mustahil kalau dibaca AND**:

```
step 7    baris 1  2/3  Local.CurrencyIDOld != ""
          baris 2  2/3  panjang daftar SpreadingRisk > 0
          baris 3  5/2  penanda ganti-mata-uang == "Interest"      <- Skip Whens
          baris 4  2/3  penanda ganti-mata-uang == "ClaimAmount"
```

Dibaca dengan A2: **`CurrencyIDOld ≠ "" DAN daftar > 0 DAN (Interest ATAU ClaimAmount)`**.

⚠️ **Kalau keempat baris dibaca AND**, syaratnya menjadi *"penanda bernilai Interest **dan**
ClaimAmount"* — **tidak mungkin terpenuhi**, dan langkahnya akan disimpulkan **tidak pernah jalan**.

| Step | Baris | Bentuk rantainya |
| --- | --- | --- |
| **7** | 4 | `CurrencyIDOld ≠ ""` **DAN** `SpreadingRisk > 0` **DAN** (`Interest` **ATAU** `ClaimAmount`) |
| **8** | 4 | idem, tetapi daftar yang diperiksa `EstimationList` |
| **9** | 5 | `CurrencyIDOld ≠ ""` **DAN** `SpreadingClaim > 0` **DAN** (`Interest` **ATAU** `ClaimAmount` **ATAU** `Estimation`) |
| **10** | 5 | idem, tetapi daftar yang diperiksa `SpreadingBreakQS` |

#### `SetMOClaimTreaty_Act` — satu baris

```
step 1    baris 1  5/2  nomor polis == ""            <- Skip Whens
          baris 2  3/2  BranchDetailID == ""
```

Dibaca dengan A2: **`nomor polis kosong` ATAU `BranchDetailID TIDAK kosong`**.

⚠️ Ronde 2 mendaftarkan `SetMOClaimTreaty_Act` sebagai rule dengan **1 baris terbalik** — itu
**benar** untuk baris kedua. Tetapi **bentuk rantai penuhnya tidak pernah dinyatakan**.

### B5 — Keluarga ketiga `pyStepsRepeatDef` `[terverifikasi]`

| `HasRepeat` | Jumlah |
| --- | --- |
| **`EMBEDDED`** | **192** |
| **`REPEAT`** | **32** |
| **TOTAL** | **224** langkah berulang, di **63 activity** |

Terbanyak: `CountTotalInsterest_Act` **19** · `CountEstimation_Act` **18** · `CountSpreading_Act`
**10** · `SaveOutstanding_Act` **10** · `CountPersen_act` **8** · `CountSpreadingADJ_Act` **8** ·
`DeleteEstimation_Act` **8** · `SaveAcceptationTreaty_Act` **8** · `SetTreatyNameSpreading_Act` **8**
· `TryMakePLA_Act` **8** · `HitServiceToKasir_Act` **7**.

#### ⚠️ Langkah berulang yang MEMUAT penulisan basis data atau efek keluar

⛔ Langkah berulangnya sendiri tidak menulis — **anaknya** yang menulis. Artinya **berapa kali
penulisan terjadi = berapa kali induknya berputar**.

| Activity · langkah berulang | Yang ada di dalamnya |
| --- | --- |
| `AddEstimation_Act` **5** | `RDB-List` konversi kurs |
| `CNMSetDetailCauseOfLoss_act` **4** | `RDB-List` sebab kerugian |
| `CekPremiLunas_Act` **7** | `RDB-List` cek proteksi |
| `CloseClaimProp` **6** | `RDB-List` simpan ke OS akseptasi |
| `GetDtlPaymentPremi_act` **7** | `RDB-List` mata uang |
| ⭐ **`GetPayAttachmentAdj_Act` 3** | **`Connect-REST`** · **dua `Obj-Save`** · **`Commit`** |
| `GetPayAttachmentAdj_Act` **3.3** | satu `Obj-Save` |

⚠️ **`GetPayAttachmentAdj_Act` langkah 3 memuat panggilan layanan luar, dua penyimpanan, dan satu
`Commit` — seluruhnya di dalam perulangan.** Berapa kali ketiganya terjadi **belum terbaca**.

### B6 — Langkah `Java` `[terverifikasi]`

Dikenali dari `pyStepsJavaSource` terisi, tidak peka huruf.

| Cara menyaring | Hasil |
| --- | --- |
| Metode **== `"Java"`** *(peka huruf)* | 29 |
| Metode **lower == `"java"`** | 29 |
| **`pyStepsJavaSource` terisi** | **29** |

⭐ **Ketiganya sama** — berbeda dari modul lain, di sini **tidak ada** ejaan huruf kecil. Aturan A7
tetap berlaku sebagai kehati-hatian, tetapi **modul ini tidak terpengaruh**.

| Ukuran | Jumlah |
| --- | --- |
| **Menyentuh basis data** | ⭐ **0** |
| **Memanggil layanan luar** | **1** |
| **Menangani kegagalan** | **9** |

**Perilakunya, dikelompokkan** ⛔ *(kode Java tidak disalin)*:

| Kelompok | Berapa | Apa yang dikerjakan |
| --- | --- | --- |
| **Pembersih duplikat** | ~12 | menyisir daftar di halaman kerja, membuang baris yang kunci mata uangnya sudah pernah muncul |
| **Perakit pesan ke Kasir** | 2 | `HitServiceToKasir_Act` **9.6** dan **13.3** — merakit teks pesan dari puluhan potongan, dibungkus `try`/`catch`, mencatat galat ke log |
| **Pengubah berkas jadi teks** | 4 | `CreatClaimAnalysis_Act` **10** · `PrintDLATreatyIn` **15.18** · `PrintFileAcceptance` **14** · `TryMakePLA_Act` **26** — membaca berkas PDF dan mengubahnya jadi teks, **melempar galat** bila gagal |
| **Pengurai jawaban JSON** | 2 | `GetDetailPolis_act` **4** · `SetValueToClaim_Act` **4** — mengurai jawaban layanan menjadi properti kasus, **mencatat galat** bila gagal |
| **Pengambil lampiran** | 1 | `GetBase64Attachment` **5.3** — ⭐ **satu-satunya yang memanggil layanan luar**: membuka aliran dari sebuah URL, membacanya, mengubahnya jadi teks; `try`/`catch` |
| **Penentu jenis berkas** | 2 | `InsertDocument_Act` · `InsertGoogleStorage_Act` — mengambil akhiran nama berkas |
| Sisanya | ~6 | penyalin dan penghitung panjang daftar |

⭐ **Nol perhitungan uang di dalam `Java`** — seluruh rumus uang terbaca dari penugasan properti
biasa. Menguatkan bab rumus di ronde 2.

### B7 — `Obj-Browse` / `Obj-List` / `Obj-Open` `[terverifikasi]`

**Lima langkah** di 112 activity. **Empat punya saringan:**

| Activity · step | Kelas | Saringan |
| --- | --- | --- |
| `GetBase64Attachment` **3** | dokumen klaim | `.IDPEGA` + empat nilai `.KATEGORI_1` |
| `GetInvoiceAttachments` **2** | dokumen klaim | `.IDPEGA` |
| `GetLinkService` **2** | alamat layanan | `.KATEGORI_1` + `.KATEGORI_2` |
| `HitServiceToKasir_Act` **11** | rekening bank | nama bank + cabang + nomor rekening |

**Satu tanpa saringan** — `GetBase64Attachment` **5.6**, tetapi metodenya **membuka satu instance
lewat kunci**, jadi **tidak perlu saringan**; kuncinya memang dikirim.

⭐ **NOL cacat nyata.** ⚠️ `GetLinkService` di modul ini **sama bentuknya** dengan di modul lain —
alamat layanan diambil dengan **dua kunci**, bukan baris pertama tanpa saringan.

### B8 — `pyXMLSignature` `[terverifikasi]`

| Ukuran | Jumlah |
| --- | --- |
| Activity yang memakainya | **48 dari 112** |
| Parameter seluruhnya | **109** |
| ⚠️ **Belum pernah tercatat di ronde 1/2** | **35** |

Yang terlewat, dikelompokkan:

| Kelompok | Contoh |
| --- | --- |
| **Penunjuk indeks** | `IndexObject` · `AdjusIdx` · `IDXADJ` · `idx` · `idxObj` · `Indezx` *(ejaan Pega)* |
| **Tanggal pemeriksaan** | `StartDate` · `EndDate` · `BeginDate` · `ReportDate` · `ReceivedDate` · `EndFor` |
| **Muatan dokumen** | `NAMAFILE` · `NOPREKAS` · `PAYMENTDATE` · `Namafile` |
| **Muatan log layanan** | `ParamInsert` · `JenisService` · `StsMessage` · `ResponMessage` |
| Lain-lain | `dcolid` · `Quarter` |

⚠️ Keempat parameter log layanan menyentuh **apa yang dicatat saat memanggil layanan luar** —
menyentuh bab efek keluar.

---

## §C — Keluarga gerbang di Section: ronde 2 dan modul lain memakai **sumbu berbeda**

⭐ **Keduanya benar, dan tidak bertentangan** — mereka memotong hal yang sama dengan pisau berbeda.

| | Sumbu | Yang dihitung |
| --- | --- | --- |
| **Ronde 2 Claim Prop** | **JENIS gerbang** | tampil · disable · wajib · hanya-baca |
| **Modul lain** | **WADAH** | tingkat sel vs tingkat layout |

### C1 — Ketiga keluarga ronde 2, ditulis ulang dengan tepat `[terverifikasi]`

#### Keluarga 1 — gerbang **tampil**

| | |
| --- | --- |
| **Wadah** | `pyUserData` — **2192 blok**, dua kelas: *HeaderElements* **1329** · *UserData* **863** |
| **Selektor mode** | `pyVisible` |
| **Nilai yang berarti HIDUP** | **`OTHER`** |
| **Nilai lain** | `ALWAYS` **1893** *(selalu tampil, syarat residu)* · kosong **129** *(tampil)* · `NOTBLANK` **9** |
| **Syaratnya di** | `pyCondition` |
| **Tingkat** | **SEL / kontrol** |

⭐ **Wadah kedua: `pyDefaultUserData` — 74 blok, seluruhnya berkelas *UserData*, dan seluruhnya
ber-`pyVisible = ALWAYS`. NOL yang hidup.**

#### Keluarga 2 — gerbang **tampil container**

| | |
| --- | --- |
| **Wadah** | elemen layout — **325 blok** |
| **Selektor mode** | `pyIsVisibilityOption` |
| **Nilai HIDUP** | **`CONDITION`** **47** · **`ExpressionCondition`** **2** |
| **Nilai lain** | `ALWAYS` **276** |
| **Syaratnya di** | `pyContainerVisibleWhen` |
| **Tingkat** | **LAYOUT** |

#### Keluarga 3 — gerbang **disable** dan **wajib**

| | Disable | Wajib |
| --- | --- | --- |
| **Selektor mode** | `pyDisabledNew` | `pyRequiredNew` |
| **Nilai HIDUP** | `true` **46** · `always` **18** | `always` **45** · `true` **10** · `truewhn` **2** |
| **Nilai mati** | `false` **535** | `false` **344** |
| **Syaratnya di** | `pyDisabledWhen` — **46 terisi** | `pyRequiredWhen` — **12 terisi** dari 26 |

### C2 — Keluarga **tanpa selektor**: gerbang **hanya-baca**

| | |
| --- | --- |
| **Wadah** | `pyReadOnlyCondition` — **3097 muncul, 128 terisi** |
| **Selektor mode** | ⭐ **TIDAK ADA** |
| **Cara membacanya** | **dibaca langsung** — bila terisi, ia berlaku |

✅ **Terbaca dari struktur.** Tidak ada tag pengaktif di sekitarnya, dan tidak ada nilai yang
berfungsi sebagai saklar. Syarat terbanyak: `pyWorkPage.IsOutstanding==1` **47** ·
`.IsAnyAcceptation =1` **14** · `.PrintFaceClaim==1` **8** · `.IsKomite=1` **6**.

⚠️ `[terbuka]` **Blok yang bersyarat TANPA selektor apa pun: NOL** `[terverifikasi]` — saya sisir
seluruh 59 berkas tampilan (36 Section + 12 FlowAction + 11 Harness) dan tidak menemukan satu pun
blok `pyUserData`/`pyDefaultUserData` yang bersyarat tetapi selektornya kosong. Jadi **"keluarga
tanpa selektor" yang dimaksud ronde 2 adalah gerbang hanya-baca di atas**, bukan sebuah blok yatim.

### C3 — HIDUP dan MATI per keluarga `[terverifikasi]`

Penanda mati: **`NEVER`** dan **`1=2`**.

| Keluarga | HIDUP | MATI |
| --- | --- | --- |
| **tampil, tingkat SEL** (`pyUserData`) | **108** | **53** |
| **tampil, tingkat SEL** (`pyDefaultUserData`) | ⭐ **0** | 0 |
| **tampil, tingkat LAYOUT** | **36** | **13** |
| **disable** | **64** *(46 bersyarat + 18 selalu)* | — |
| **wajib** | **57** *(12 bersyarat + 45 selalu)* | — |
| **hanya-baca** | **128** | — |

✅ **`pyActionConditions` = 17 baris** di **3 berkas** — `Section/AdjustmentDetail.xml` **10** ·
`Section/OutstandingClaim.xml` **6** · `Section/InputAcceptation_Est.xml` **1**.
⭐ **Angka ronde 2 COCOK PERSIS.** R7-baru berdiri.

### C4 — ⭐ Untuk dibawa ke modul lain

> ⛔ **Berkas modul lain tidak disentuh.** Ini ringkasan untuk dibaca, bukan untuk dipindahkan
> sendiri.

**Gerbang di Section punya LIMA jenis, bukan dua.** Modul yang hanya mengenali dua *(tampil-sel dan
tampil-layout)* **melewatkan tiga**: disable, wajib, dan hanya-baca.

| Jenis | Selektor mode | Nilai HIDUP | Syarat di | Tingkat |
| --- | --- | --- | --- | --- |
| **tampil — sel** | `pyVisible` | `OTHER` | `pyCondition` | sel |
| **tampil — layout** | `pyIsVisibilityOption` | `CONDITION` / `ExpressionCondition` | `pyContainerVisibleWhen` | layout |
| **disable** | `pyDisabledNew` | `true` / `always` | `pyDisabledWhen` | sel |
| **wajib** | `pyRequiredNew` | `true` / `always` / `truewhn` | `pyRequiredWhen` | sel |
| **hanya-baca** | ⭐ **tidak ada** | — | `pyReadOnlyCondition` *(dibaca langsung)* | sel |

**Dua wadah untuk gerbang tampil-sel, bukan satu:** `pyUserData` **dan** `pyDefaultUserData`.
Keduanya berisi kelas *HeaderElements* atau *UserData*. ⚠️ Di Claim Prop `pyDefaultUserData`
**nol yang hidup** — tetapi itu fakta modul ini, **bukan aturan umum**.

**Penanda mati sama di semua jenis:** `NEVER` dan `1=2`.

**Dan satu jenis lagi di luar kelima itu:** `pyActionConditions` — **gerbang pada aksi tombol**,
jarang tetapi nyata. ⚠️ `[terbuka]` cara menggabungkan beberapa barisnya *(AND atau OR)* **tidak
terbaca dari korpus** — jangan disimpulkan.

---

## §D — Kesimpulan mana yang goyah

### D5 — Pernyataan tegas lebih dulu

> ⭐ **TIDAK ADA satu pun kesimpulan ronde 1 atau ronde 2 yang TERBALIK.**
>
> Disensus **100 %**: 1428 baris keluarga pertama · 1365 baris keluarga kedua · 224 langkah
> berulang · 29 langkah `Java` · 5 langkah `Obj-*` · 109 parameter · 59 berkas tampilan.
>
> Yang ada **tiga yang GOYAH**, dan ketiganya **kekurangan keterangan**, bukan kekeliruan.

### D1 — Daftar lengkap

#### ✅ P13 — arti enum kode arah — **TERJAWAB, bukan goyah**

Ronde 1 mendaftarkannya sebagai butir terbuka, pemiliknya *"pemilik export Pega"*. A1 menjawabnya.
Seluruh pembacaan arah di ronde 1 dan 2 dibuat **tanpa** peta ini, memakai tiga aturan yang
diturunkan sendiri — dan **ketiganya ternyata benar**: `2/3` normal, `3/2` terbalik, `6` keluar
activity. **Nol yang perlu diperbaiki.**

#### ⚠️ GOYAH 1 — *"89 memakai kode 5/6"* menggabung dua hal yang berbeda

**Di mana:** ronde 2, sensus baris syarat pada gerbang aktif.

**Mengapa goyah:** kode **5** dan kode **6** digabung menjadi **satu ember** bernama *"belum
diketahui artinya"*. Kini keduanya terbukti **sangat berbeda** — 5 mengubah rantai syarat menjadi
**OR**, 6 **menghentikan activity**. Satu mengubah **makna syarat**, satu mengubah **alur**.

**Buktinya:** pada dasar yang sama, kode **5 = 7 baris**, kode **6 = 16 baris**.

**Statusnya: GOYAH, bukan terbalik.** Tidak ada kesimpulan turunan yang bersandar pada ember
gabungan itu — ia hanya angka sensus. Tetapi **siapa pun yang memakainya untuk memperkirakan
kerumitan akan salah**, karena tujuh di antaranya mengubah logika syarat.

#### ⚠️ GOYAH 2 — bentuk rantai syarat tujuh langkah berkode 5 **tidak pernah dinyatakan**

**Di mana:** ronde 2 mencatat rumus uang di `CopyOldataCurr` **7.1 · 9.1 · 10.1** *(anak)*, dan
mendaftarkan `SetMOClaimTreaty_Act` sebagai punya **1 baris terbalik**.

**Mengapa goyah:** gerbang **langkah induk** — 7, 8, 9, 10 — **tidak pernah dituliskan**. Dengan A2
ia terbaca **OR**; kalau dibaca **AND** ia **mustahil terpenuhi**, dan langkahnya akan disimpulkan
**mati**.

**Akibat praktisnya:** ⭐ **kesimpulan ronde 2 tentang rumus uang MASIH BERDIRI** — dengan OR,
keempat langkah itu memang berjalan. Yang hilang hanyalah **alasannya**. Bahayanya muncul bila
orang lain membaca ulang langkah-langkah itu dan menyimpulkan sebaliknya.

**Statusnya: GOYAH ringan — lubang keterangan, bukan kekeliruan.**

#### ⚠️ GOYAH 3 — keluarga kedua belum pernah dibaca, dan ada satu jalur kegagalan yang hilang

**Di mana:** seluruh ronde 1 dan 2 — keluarga kedua tidak pernah disensus.

**Apa yang ditemukan:** dari 1365 baris, hanya **15** membawa tindakan khusus. Satu di antaranya
**hidup dan bermakna**:

```
SaveAdjusterConsultant_Act  step 4  Obj-Save
    keluarga kedua:  StepStatusFail  ->  lompat ke tanda FAIL, yaitu langkah 8
```

⭐ **Itu penanganan kegagalan penyimpanan, dan belum pernah tercatat di mana pun.** Langkah 8
sendiri bergerbang **mati**, sehingga **lompatannya sampai ke langkah yang tidak berbuat apa-apa**,
lalu alur berlanjut ke langkah 9.

**Statusnya: GOYAH.** `issues/08` mencatat langkah **5** activity itu *"di-remark, tidak ditulis
sama sekali"* — **pernyataan itu tetap benar**. Yang kurang adalah **jalur kegagalan langkah 4**.

#### ✅ Yang diperiksa dan **MASIH BERDIRI**

| Kesimpulan | Diperiksa terhadap | Hasil |
| --- | --- | --- |
| **44 baris TERBALIK** | sensus ulang gerbang aktif | ✅ **cocok persis 44** |
| `CheckNoPolicy` step 2 arah terbalik *(AC 108)* | rantai syaratnya | ✅ **satu baris saja** — tidak ada pertanyaan AND/OR |
| `CheeckNoRNM_Act` step 4 arah terbalik *(AC 111)* | idem | ✅ **satu baris saja** |
| `CountValueADJTreaty_Act` step 16/17, gerbang mati dan hidup *(§14a spec)* | keluarga kedua | ✅ **berdiri** — lompatan menggantung di step 17 **bergerbang mati**, tidak mengubah apa pun |
| `pyActionConditions` 17 baris di 3 berkas *(R7-baru)* | hitung ulang anak elemen | ✅ **cocok persis** |
| Rumus uang dan pembulatan bertingkat | langkah `Java` | ✅ **berdiri** — nol perhitungan uang di `Java` |
| Aturan R1–R6 gerbang Section | sensus 59 berkas tampilan | ✅ **keenamnya berdiri** |

⭐ **Tidak ada pernyataan *"nol"* / *"tidak ada"* / *"tidak ketemu"* di ronde 1 atau 2 yang
terbantah** oleh ketujuh aturan baru. Yang paling rawan — *"`pyActionConditions` kosong di
mana-mana"* — **sudah diralat sendiri oleh ronde 2** menjadi R7-baru, dan ralat itu **terbukti
benar**.

### D2 — Bab `spec.md` yang isinya berubah

⛔ **`spec.md` tidak disunting.**

| Bab | Kalimat yang terkena | Perubahannya |
| --- | --- | --- |
| **Implementation Decisions** — bagian rumus uang dan pembulatan | daftar `/100` polos yang menyebut `CopyOldataCurr` **7.1 · 9.1 · 10.1** | **isi tidak berubah**; yang perlu **ditambahkan** adalah bentuk gerbang langkah induknya — `(Interest ATAU ClaimAmount)` |

⭐ **Satu bab, dan sifatnya penambahan keterangan — bukan koreksi.** Bab **Acceptance Criteria**
diperiksa dan **AC 108 · 111 berdiri**; bab **§14a** tentang `CountValueADJTreaty` **berdiri**.

### D3 — Tiket yang isinya terpengaruh

⛔ **Tiket tidak disunting.** Tiga dari 16:

| Tiket | Bagian yang terkena | Sifatnya |
| --- | --- | --- |
| **`00`** prefactor skema relasional transaksi uang | daftar bentuk pembagian `/100` | **penambahan keterangan** — gerbang induk `CopyOldataCurr` |
| **`06`** loss allocation dan spreading | rumus yang menyebut `CopyOldataCurr` | **penambahan keterangan** — idem |
| **`08`** baris adjustment dan adjuster | tabel rule sumber, baris `SaveAdjusterConsultant_Act` | ⭐ **penambahan nyata** — jalur kegagalan langkah 4 belum tercatat |

### D4 — Diurutkan dari yang paling berbahaya

*"Berbahaya"* = paling mungkin sudah dipakai orang untuk mulai membangun.

| # | Temuan | Kenapa berbahaya |
| --- | --- | --- |
| **1** | ⚠️ **Jalur kegagalan `SaveAdjusterConsultant_Act`** *(GOYAH 3)* | Tiket **08** sudah ditulis dan **siap dikerjakan**. Seseorang bisa membangun penyimpanan adjuster **tanpa** jalur kegagalan sama sekali, dan tidak ada yang menandai kekurangannya |
| **2** | ⚠️ **Gerbang induk `CopyOldataCurr` 7/8/9/10** *(GOYAH 2)* | Menyentuh **dua tiket** (00, 06) dan **satu bab spec**, semuanya tentang **angka uang**. Pembaca yang menurunkan gerbangnya sendiri dengan aturan AND akan menyimpulkan langkahnya mati, lalu **membuang rumusnya** |
| **3** | ⚠️ **Perulangan yang membungkus `Commit` dan panggilan layanan** *(`GetPayAttachmentAdj_Act` 3)* | Belum masuk tiket mana pun, tetapi **berapa kali `Commit` terjadi** adalah pertanyaan yang harus dijawab sebelum jalur lampiran dibangun |
| **4** | ⚠️ **35 parameter yang belum tercatat** | Empat di antaranya menyentuh **apa yang dicatat saat memanggil layanan luar** |
| **5** | ⚠️ **"89 kode 5/6" sebagai satu ember** *(GOYAH 1)* | Angka sensus saja; **tidak ada kesimpulan turunan** yang bersandar padanya |

---

## §E — Apa lagi yang seharusnya dikerjakan

### E1 — Urutan sesudah periksa ulang ini

| # | Yang dikerjakan | Kenapa | Besar | Menunggu |
| --- | --- | --- | --- | --- |
| **1** | **Perulangan yang membungkus penulisan dan efek keluar** — mulai dari `GetPayAttachmentAdj_Act` **3** | ia membungkus **`Commit`** dan **panggilan layanan luar**; berapa kali keduanya terjadi belum terbaca | sedang | korpus |
| **2** | **Sisa 224 langkah berulang** — halaman apa yang diulang | menentukan berapa kali setiap rumus uang berjalan; menyentuh seluruh bab angka | berat | korpus |
| **3** | **35 parameter `pyXMLSignature`** — ke mana masing-masing dipakai | melengkapi bab efek keluar dan bab dokumen | sekali sisir | korpus |
| **4** | **61 rule `When`** — disebut 131 kali di ronde 1/2 tetapi **belum pernah disensus sebagai rule** | gerbang bernama dipakai di seluruh modul; isinya belum dibaca | sedang | korpus |
| **5** | **11 `DataTransform` · 19 `ReportDefinition` · 6 `ConnectREST`** | disebut sedikit, isinya belum dibaca | sedang | korpus |

⛔ **Keempat belas tiket dan spec tidak perlu ditulis ulang** — temuan di atas menambah keterangan,
bukan membalik keputusan.

### E2 — Wadah yang belum pernah diparse: ⚠️ **ADA**

⛔ **Dicari, tidak ditunggu.** Di tingkat **langkah**, keempat keluarga sudah habis. Yang tersisa di
tingkat **rule** dan **jenis rule**:

| Wadah / jenis | Sebaran | Keadaan |
| --- | --- | --- |
| **`pyMemo`** | **303 berisi** dari 303 | ⚠️ **BELUM** — catatan perubahan terakhir penulis rule. Di modul lain ini **mengubah pembacaan** |
| **`pyUsage`** | **870 berisi** dari 1205 | ⚠️ **BELUM** |
| `pyActivityPrivilegeList` | 112 berisi | ⚠️ **BELUM** — ⛔ di modul lain `[keputusan work owner]` **diabaikan**; berlaku juga di sini atau tidak, **belum diputuskan** |
| `pyLocalParameters` · `pySystemParameters` | 112 berisi | ⚠️ BELUM |
| `pxNamedPageReferences` | 110 berisi | ⚠️ BELUM |
| `pyPagesAndClasses` | 110 berisi | sebagian |
| **jenis rule `When`** | **61 berkas** | ⚠️ **BELUM sebagai rule** |
| **jenis rule `DataTransform`** | **11 berkas** | ⚠️ BELUM |
| **jenis rule `ReportDefinition`** | **19 berkas** | ⚠️ BELUM |
| **jenis rule `ConnectREST`** | **6 berkas** | ⚠️ BELUM |

⭐ **`pyMemo` yang paling menjanjikan** — di modul lain ia memuat catatan seperti *"add
StepStatusFail"* dan *"FIX ERROR HANDLING"* yang **langsung mengubah pembacaan**. Di sini ada
**303 catatan, seluruhnya terisi**, dan **belum satu pun dibaca**.

### E3 — Kesimpulan yang paling rawan salah

⛔ **Ditunjuk, tidak diperbaiki.**

1. ⛔ **Klaim "nol perhitungan uang di `Java`"** *(§B6 berkas ini)*. **Paling rawan**, dan
   kelemahannya pola yang sama yang sudah membakar modul lain: saya menyimpulkan **ketiadaan**.
   Ia ditarik dari **29 langkah yang saya baca**, tetapi **`pyMemo` dan `pyUsage` belum dibaca** —
   dan keduanya bisa menyebut perubahan yang jejaknya tidak terbaca di tempat lain.
2. ⚠️ **Pembacaan rantai OR pada `CopyOldataCurr`** *(§B4)*. Bersandar pada A2, yang **dipakai apa
   adanya dari modul lain** dan **tidak diuji di modul ini**. Bila A2 keliru, keempat langkah itu
   terbaca terbalik.
3. ⚠️ **Angka "482 baris pada gerbang aktif"** *(§B1)*. Definisi *"aktif"* saya rumuskan sendiri
   karena ronde 2 tidak menuliskannya. Angka terbaliknya cocok persis, tetapi dua angka lain tidak
   — **jadi salah satu dari dua definisi itu tidak sama**, dan saya belum tahu yang mana.

### E4 — Pertanyaan **BARU** untuk work owner — **2**

#### Pertanyaan 1 — Penyimpanan data adjuster punya jalur kegagalan yang berakhir tanpa berbuat apa-apa. Dipertahankan?

**Apa yang ditanyakan.** Ketika penyimpanan data adjuster dan konsultan **gagal**, alur melompat ke
sebuah titik penanganan — dan titik itu **sengaja dimatikan**, sehingga tidak berbuat apa pun.
Alur lalu berlanjut seolah tidak terjadi apa-apa. Apakah perilaku itu dipertahankan.

**Kenapa muncul.** Baru terbaca sekarang; keluarga gerbang yang memuatnya belum pernah dibaca di
modul ini. Tiket **08** sudah ditulis tanpa menyebutnya.

**Bedanya jawaban A atau B.** Bila **dipertahankan**, kegagalan penyimpanan adjuster **diam-diam
diabaikan**, konsisten dengan preseden *"tiru apa adanya"* — dan tiket 08 perlu menyebutnya supaya
tidak dikira kelalaian. Bila **tidak**, perlu ditetapkan apa yang terjadi saat gagal, dan itu
**perubahan perilaku**, bukan pemindahan.

**Yang tertahan.** Tiket 08, dan sikap terhadap kegagalan penyimpanan di seluruh modul.

#### Pertanyaan 2 — Satu perulangan membungkus panggilan layanan luar dan satu penyimpanan-akhir. Berapa kali seharusnya terjadi?

**Apa yang ditanyakan.** Pada jalur lampiran pembayaran, satu langkah berulang membungkus
**panggilan ke layanan luar**, **dua penyimpanan**, dan **satu penyimpanan-akhir**. Artinya
ketiganya terjadi **sekali per putaran**. Berapa kali putaran itu seharusnya terjadi, dan apakah
memang dimaksudkan begitu.

**Kenapa muncul.** Perulangan belum pernah disensus di modul ini. Halaman yang diulang **tidak
terbaca dari struktur**, jadi jumlah putarannya belum diketahui.

**Bedanya jawaban A atau B.** Bila **memang per lampiran**, maka satu klaim dengan sepuluh lampiran
menghasilkan **sepuluh panggilan layanan dan sepuluh penyimpanan-akhir** — dan itu perlu ditiru
apa adanya. Bila **seharusnya sekali**, maka yang ada sekarang adalah **cacat yang belum ketahuan**,
dan memindahkannya akan memindahkan cacatnya.

**Yang tertahan.** Jalur lampiran pembayaran, dan berapa banyak transaksi yang dihasilkan satu
klaim.
