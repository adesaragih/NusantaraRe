# EDM Treaty In — Grilling Ronde 1
## Tiga pohon data sejajar, dua panggilan ke sistem luar, dan satu langkah yang menghapus produksi

Modul utama: **`D:\XML\RNM_BRD\EDM Treaty In`**. `NB Treaty In` dibaca **hanya sebagai pembanding** —
⛔ temuannya **tidak masuk sensus modul ini**. Korpus **READ-ONLY**; menulis hanya ke
`OUTPUT_HASIL_RNM\`. ⛔ `D:\XML\nusantara-re\` **tidak disentuh**.

⛔ **Nol kode · nol DDL · nol usulan daftar kolom · nol butir terbuka ditutup · nol nama orang
disalin.** ⭐ Pertanyaan baru ada di **`PERTANYAAN-RONDE-1.md`**, bernomor mulai **P50**.

---

## §0 — Apa yang SAMA dengan NB Treaty In

⭐ **Dibaca lebih dulu supaya pembaca tahu apa yang TIDAK perlu dibaca ulang.**

### 0.1 Sensus — enam angka brief dicocokkan

| Yang dicocokkan | Brief | Terukur | |
| --- | ---: | ---: | :---: |
| berkas `.xml` | 163 | **163** | ✅ |
| byte | 17.639.246 | **17.639.246** | ✅ |
| bernama sama dengan NB | 73 | **73** | ✅ |
| kelas berbeda | 3 | **3** | ✅ |
| jumlah langkah berbeda | 1 | **1** | ✅ |
| hanya ada di EDM | 90 | **90** | ✅ |
| `RDBList` | 36 | **36**, ⭐ **36 ber-SQL berisi** | ✅ |

**Sebaran tipe rule** `[terverifikasi]`: `Activity` **66** · `RDBList` **36** · `Section` **18** ·
`When` **11** · `DataTransform` **10** · `ReportDefinition` **9** · `FlowAction` **6** ·
`Harness` **3** · `ConnectREST` **1** · `DecisionTable` **1** · `Flow` **1** · `SystemSettings` **1**.

### 0.2 ⚠️ "Isinya sama — beda hanya metadata ekspor" — hampir benar, dan selisihnya perlu disebut

⭐ **Dihitung dua cara, dan keduanya tidak sepakat:**

| Cara | Hasil |
| --- | ---: |
| **A** — hash isi sesudah **24 tag volatil** dibuang | ⛔ **0 identik** |
| ⭐ **B** — tanda tangan perilaku *(kelas + jumlah langkah + jumlah SQL)* | ⭐ **69 identik** |

⭐ **Yang dipercaya: cara B — 69**, dan ⭐ keempat yang berbeda menurut cara B **persis** yang brief
sebut: `CountSpreading_Act` · `SetCategoryAttach` · `serviceInsertArasapas_act` · `IsUW`.

⚠️ **Kenapa cara A gagal, dan ini instrumen saya:** metadata ekspor yang berbeda ternyata
**26 tag**, bukan 24. Satu yang saya lewatkan — `pyShowJavaWindowName` — ditemukan dengan
mendiff satu berkas apa adanya. ⛔ Bahkan sesudah ditambahkan, cara A tetap 0, sebab sisanya
menyebar tipis di berkas berbeda.

⛔⛔ **Dan satu koreksi atas bunyi brief, yang perlu dicatat:** *"beda hanya metadata ekspor"*
⚠️ **tidak seluruhnya benar.** `[terverifikasi]` **Enam berkas ber-`pyMemo` berbeda** — catatan
pengembang, bukan metadata:

| Berkas | Catatan NB | Catatan EDM |
| --- | --- | --- |
| `Activity\CountOGPONP_Act` | *"buka when step 8"* | *"BUKAN WHEN STEP 8"* |
| `Activity\CountResult1_Act` | *"perbaiki step 4 utk nonprop"* | *"hitung GrossPremium utk nonprop"* |
| `Activity\CountSpreading_Act` | *"ubah jangan pake param"* | *"."* |
| `Activity\SetCategoryAttach` | *"save life"* | *"tambah tarikan rd"* |
| ⭐ `Activity\serviceInsertArasapas_act` | *(kosong)* | ⭐ *"ADD Hit service Arasapas 2 - FACOUT."* |
| ⭐ `When\IsUW` | *"hps ReasFacInMarketing"* | ⭐ *"add conditon not show if Underwriting upper"* |

⭐ **Empat di antaranya ada di dalam kelompok 69 yang "identik"** — ⛔ jadi *"beda hanya metadata"*
melewatkan bahwa **catatan pengembangnya berbeda**, dan dua catatan itu **menjelaskan perilaku**.

### 0.3 ⭐ Dua belas ketetapan NB — berlaku, dengan satu yang sudah TERBUKTI berbeda

⭐ Berkas keadaan NB memuat 12 ketetapan work owner. ⭐ **Sebelas berlaku tanpa bantahan** di EDM.

⛔⛔ **Satu berbeda, dan itu temuan:**

> ⭐ **Ketetapan NB nomor 6** — *"setiap query menulis skema `POOLDATA.` eksplisit"*.
> ⭐ Di EDM, `RDBList\UpdateErrorNoteJsonPolisMonitoring.xml` **sudah menulisnya eksplisit**
> `[terverifikasi]` — ⭐ jadi ketetapan itu **tidak dilanggar**, melainkan **sudah dipenuhi
> sebagian di sistem lama**. ⚠️ Bukan pertentangan; dicatat karena menunjukkan praktiknya **tidak
> seragam** di seluruh korpus.

⚠️ **Dan satu ketetapan yang EDM justru sudah lebih dekat:** ⭐ **ketetapan nomor tangga tiga
jenjang**. `[terverifikasi]` Alur EDM menyebut **hanya tiga antrean** — `ReasTreatyInAdmin`,
`ReasTreatyInSecHead`, `ReasTreatyInDeptHead`. ⛔ **Nol** penyebutan `GroupLeader` atau `Director`.
⭐ Jadi keputusan work owner untuk NB **sudah menjadi kenyataan di EDM**.

---

## §A — Sasaran 1: data lama, data baru, dan selisihnya

### A.1 ⭐⭐ Polanya bukan perbandingan — melainkan TIGA POHON SEJAJAR

⛔ **Dugaan yang wajar dan ternyata salah:** bahwa ada satu kumpulan medan yang dibandingkan
lama-lawan-baru. ⭐ **Yang sebenarnya:** nama medan **sama persis**, dan yang membedakannya adalah
**awalan halaman**.

| Pohon | Bentuk properti | Contoh |
| --- | --- | --- |
| ⭐ **baru** | **telanjang** | `.GrossPremium` · `.NetPremium` · `.Claim` · `.Deduction1` |
| ⭐ **lama** | ⭐ **`.OldData.`** | `.OldData.GrossPremium` · `.OldData.NetPremium` |
| ⭐ **selisih** | ⭐ **`.TreatyDifference.`** | `.TreatyDifference.GrossPremium` · `.TreatyDifference.NetPremium` |

`[terverifikasi]` Perintah audit: urai sel ber-`pyType` = `FIELD`, ambil `pyValue` berawalan titik,
lalu bandingkan himpunannya antar berkas.

### A.2 Medan apa yang dibandingkan — dan apa yang TIDAK

| Berkas | byte | sel `FIELD` | properti unik | `pxCurrency` | `pxNumber` |
| --- | ---: | ---: | ---: | ---: | ---: |
| `…PropOldData` | 745.191 | 40 | **38** | 14 | 16 |
| `…PropNewData` | 674.337 | 34 | **33** | 12 | 16 |
| ⭐ `…PropValueDifference` | 745.389 | 40 | **38** | 14 | 16 |
| `…PropOldData2` | 746.314 | 44 | **42** | 14 | 20 |
| ⭐ `…PropNewData2` | 947.721 | 43 | **39** | 15 | 16 |

⭐⭐ **Jawaban atas "medan apa yang tidak dibandingkan": NOL.**
`[terverifikasi]` ⭐ **Himpunan `Old ∩ New` yang tidak muncul di `Difference` = 0.**
⛔ Artinya **setiap medan yang punya versi lama dan versi baru juga punya versi selisih** —
tidak ada yang dilewatkan.

⚠️ **Dan selisih dihitung untuk medan uang DAN bukan-uang.** `[terverifikasi]` Berkas selisih
memuat **14 medan ber-kendali mata uang** dan **16 ber-kendali angka** — ⭐ jadi **bukan hanya
uang**; persentase dan cacah ikut diselisihkan.

⭐ **Dua puluh delapan medan ber-awalan `.TreatyDifference.`** terbaca, di antaranya:
`BalanceBeforePPH` · `BalanceBeforeTax` · `BalanceDueTo` · `Claim` · `Deduction1` · `Deduction2` ·
`ExcessLoss` · `GrossPremium` · `Installment` · `NetPremium` · `OveriddingCommOgp` ·
`OveriddingCommOnp` · `PPHValue` · `PPNValue`.

### A.3 ⭐⭐ Arti akhiran `2` — bukan dua tampilan, melainkan SARANG

⭐ **Terjawab dari korpus.** `[terverifikasi]` `…PropOldData2` memuat properti berawalan
⭐ **`.OldData.TreatyDifference.`** — ⛔ **awalan bersarang dua tingkat.**

| Berkas | Awalan yang dipakai |
| --- | --- |
| `…PropOldData` | `.OldData.` |
| ⭐ `…PropOldData2` | ⭐ **`.OldData.TreatyDifference.`** |

⭐ **Bacaannya:** berkas ber-akhiran `2` menampilkan **selisih yang tersimpan DI DALAM data lama** —
`[dugaan]` yaitu **selisih dari endorsemen sebelumnya**. ⛔ **Belum terverifikasi** apakah itu berarti
endorsemen berlapis, atau sekadar penyimpanan selisih terdahulu. → **P57**.

`[terverifikasi]` Irisan `OldData` dan `OldData2` = **15**; hanya di `OldData2` = **27**; hanya di
`OldData` = **23**.

### A.4 Bagaimana selisih dihitung

`[terverifikasi]` `Activity\CalculateDifferenceEDM_act.xml` — **216.239 B**, kelas
**`ASM-FW-GISFW-Int-treaty_in_edm`**, ⭐ **9 langkah, seluruhnya `Property-Set`**.

⭐ **Catatan langkahnya terbaca, dan ia menerangkan bentuknya:**

| Catatan langkah | Yang ia beri tahu |
| --- | --- |
| *"Set the nested value"* · *"Set parent value & subscript"* | ⭐ struktur **bersarang berindeks** |
| *"Valuelist on treatydifference"* | ⭐ selisih disimpan sebagai **daftar nilai** |
| ⚠️ *"Jika master ada prorate"* | ⭐ ada cabang **prorata** dari polis induk |
| ⭐⭐ *"EDM BATAL"* | ⛔ ada keadaan **endorsemen dibatalkan** — → **P56** |
| *"Calculate total value"* · *"Set Total to corresponding currency"* | ⭐ total **per mata uang** |

⛔⛔ **Tetapi rumusnya TIDAK terbaca** — kesembilan langkah adalah `Property-Set`, dan
`[terverifikasi]` **nilainya tidak terekspor**, persis seperti di NB Treaty In. ⭐ Jadi **bentuk
selisih diketahui; cara menghitungnya tidak.**

⚠️ **Struktur bersarang yang terbaca dari catatan:**
`pyWorkPage.PolicyTreatyIn.TreatyXOLDifferenceList(<indeks>).ValueList(<indeks>)` — ⭐ **daftar di
dalam daftar**.

⚠️ **Apakah data lama DISALIN atau DIBACA dari polis induk belum terjawab** — ⭐ adanya awalan
`.OldData.` sebagai **halaman tersendiri** `[dugaan]` menunjukkan **salinan**, ⛔ tetapi belum
terverifikasi. → **P58**.

---

## §B — Sasaran 2: adendum dan premi tambahan

### B.1 ⭐ Keduanya HAL BERBEDA — dibuktikan dari kelas dan isi medannya

| | Adendum | Premi tambahan |
| --- | --- | --- |
| Layar | `Section\DetailPolicyTreatyInAddendum` **913.122 B** | `Section\DetailPolicyTreatyInAddPremi` **688.183 B** |
| Kelas | `ASM-FW-GISFW-Data-PolicyTreatyIn` | `ASM-FW-GISFW-Data-PolicyTreatyIn` |
| Properti unik | ⭐ **34** | **14** |
| Isinya | ⭐ **medan kontrak lengkap** — nama bisnis, ceding, jenis klaim, mata uang, tanggal | ⭐ **hanya medan uang** — premi bruto, neto, potongan, pajak |
| Layar rinci | — | ⭐ `Section\DetailPolicyAddPremiDetail` **249.693 B**, kelas ⭐ **`ASM-FW-GISFW-Data-XOLRealisasiData`**, memuat **lapisan** |

⭐⭐ **Bacaannya:** **adendum** mengubah **isi kontrak**; **premi tambahan** menambah **uang** pada
kontrak yang sudah ada, dan rinciannya **per lapisan XOL**. ⛔ **Bukan dua nama untuk satu hal.**

⚠️ **Kapan masing-masing terbit belum terjawab** — → **P55**.

### B.2 ⭐⭐ Adendum punya PENOMORAN SENDIRI

`[terverifikasi]` `RDBList\GenerateNoEDMTreaty.xml`:

> ⭐ Pola: awalan **`RNM-E`** · kode bisnis lama · sebuah slot parameter · **nomor urut 5 digit**
> dari deret **`POOLDATA.ENDORSEMENT_SEQ`**

⭐ **Deret terpisah, awalan terpisah.** ⛔ Berbeda dari penomoran realisasi NB, yang memakai deret
lain dan penanda `.T`.

### B.3 ⭐⭐ Aturan batas tanggal 25 — tidak ada padanannya di NB

`[terverifikasi]` `Activity\GeneratePolicyNoTreatyAddendum_Act.xml` — **124.368 B, 13 langkah**.
⭐ Dua catatan langkahnya menyatakan aturan bisnis yang **belum pernah terlihat**:

| Catatan langkah |
| --- |
| ⭐ *"when not above 25 in that month"* |
| ⭐⭐ *"Determine if it's past 25 on current month, if true move to next month"* |

⛔ **Artinya: endorsemen yang dibuat sesudah tanggal 25 diberi nomor pada bulan BERIKUTNYA.**
⚠️ Itu aturan **periode produksi**, dan ⛔ **alasannya belum diketahui** — → **P54**.

⚠️ Ditambah catatan *"Set error blm generate No EDM"* — ⭐ ada keadaan nomor belum terbentuk.

### B.4 ⭐ Penyisipan ke produksi bersifat idempoten

`[terverifikasi]` `Activity\InsetTreatyInProdAddendum_Act.xml` — **307.657 B**, kelas
`ASM-FW-GISFW-Work`, **22 langkah** *(16 `Property-Set` + 6 `RDB-List`)*.

⭐ **Catatan langkah yang menentukan:** *"when there is data already, exit activity"* —
⭐ **penjaga ganda-sisip yang eksplisit**, dan ⭐ itu praktik yang lebih baik daripada yang
ditemukan di modul lain.

Urutannya terbaca: salin kunci kasus → ambil nomor polis dari kunci → ambil tanggal input →
tetapkan nomor polis dan tanggal → salin tanggal mulai dan akhir → ambil data produksi →
⭐ **keluar bila sudah ada** → cabang *"Others than Non Prop"* → tetapkan nilai bersarang.

---

## §C — Sasaran 3: panggilan ke sistem luar

⛔⛔ **NB Treaty In tidak punya ini sama sekali.** ⭐ Ini bagian paling berisi di modul EDM, dan
⭐⭐ **ia menutup lubang yang menahan NB** — butir **OQ-025** dan **P8**.

### C.1 ⭐⭐ Rantai lengkap — 20 langkah, berurutan

`[terverifikasi]` `Activity\serviceInsertArasapas_act.xml` — kelas **`ASM-FW-GISFW-Work`**,
**232.000 B**, ⭐ **20 langkah**. ⛔ Bandingkan NB: **18.013 B, satu langkah** — sekadar pembungkus.

| Lgk | Metode | Catatan langkah |
| ---: | --- | --- |
| 1 | `Obj-Refresh-And-Lock` | ⭐ **kasus dikunci lebih dulu** |
| 2 | `Property-Set` | *"Get value CaseId"* |
| 3 | `RDB-List` | *"Get PolicyNo by CaseId"* |
| 4 | `Property-Set` | *"Set PolicyNo and Set Sysdate"* |
| ⭐ 5 | `Call …M_LINK_SERVICE.GetLinkService` | ⭐ *"GET LINK SERVICE"* |
| ⭐⭐ **6** | **`Connect-REST`** | ⭐⭐ *"Hit service Arasapas 1 — **FACIN**"* — ⭐ **tanpa gerbang** |
| ⭐⭐ **7** | **`Connect-REST`** | ⭐⭐ *"Hit service Arasapas 2 — **FACOUT**"* — ⚠️ **bergerbang** |
| 8 | `Property-Set` | *"SET ERROR GAGAL KONVERSI"* — bergerbang |
| 9 | `Obj-Save` | tanpa gerbang |
| ⛔ **10** | `Page-Set-Messages` | ⛔⛔ **MATI** — *"SET ERROR GAGAL KONVERSI"* |
| 11 | `RDB-List` | halaman `stsKonversi` |
| ⛔⛔ **12** | `RDB-List` | ⛔⛔ *"**DELETE PRODUKSI JIKA ERROR (PROCEDURE)**"* — bergerbang |
| 13 | `Property-Set` | contoh nilai kunci kasus |
| 14 | `Property-Set` | *"Get Data Pega To JSON_POLIS_MONITORING"* |
| 15 | `RDB-List` | *"Insert to JSON_POLIS_MONITORING"* — bergerbang |
| 16 | `RDB-List` | *"Set err_note sts_konversi 9"* — bergerbang |
| 17 | `Page-Remove` | pembersihan halaman |
| ⛔ **18** | `Call ASMForceCaseClose` | ⛔⛔ **MATI** |
| 19 | `Property-Set` | *"SET PARAM EMAIL"* |
| ⭐ 20 | `Call SendEmailWithAttachments` | ⭐ *"SEND EMAIL JIKA ERROR KONVERSI"* |

⭐ **Perintah audit:** urai `rowdata` secara **rekursif** dengan `ElementTree`, baca
`pyStepsActivityName` · `pyStepsDescription` · `pyStepsBlockName` · `pyStepsPreCondition`.

### C.2 ⭐ Tiga temuan yang menentukan

**1. ⭐⭐ DUA panggilan, bukan satu — dan yang kedua bergerbang.**
`[terverifikasi]` Langkah **6** *(FACIN)* ber-`pyStepsPreCondition` = `false` ⇒ ⭐ **tanpa gerbang,
selalu jalan**. Langkah **7** *(FACOUT)* ber-`pyStepsPreCondition` = `true` ⇒ ⭐ **bergerbang**.
⛔ **Syarat gerbangnya belum saya baca** — → **P50**.

**2. ⛔⛔ Sebuah langkah MENGHAPUS DATA PRODUKSI bila terjadi galat.**
`[terverifikasi]` Langkah **12**: *"DELETE PRODUKSI JIKA ERROR (PROCEDURE)"*, dijalankan lewat
`RDB-List`, **bergerbang**. ⛔ **Apa yang dihapus, seluas apa, dan apakah dapat dikembalikan
belum diketahui** — → **P51**.

**3. ⛔⛔ Dua langkah MATI, dan keduanya mengubah arti jalur galat.**
`[terverifikasi]` `pyStepsBlockName` berawalan `//`:

| Lgk | Yang mati | Akibatnya |
| ---: | --- | --- |
| ⛔ **10** | pemasang **pesan galat** | ⭐ **pengguna TIDAK PERNAH melihat galat konversi** di layar |
| ⛔ **18** | **penutup paksa kasus** | ⭐ kasus **tidak ditutup** sesudah galat |

⭐ ⇒ **Jalur galat yang hidup hanyalah:** tulis ke tabel pemantauan, perbarui catatan galat, dan
**kirim surel berlampiran**. ⛔ **Tidak ada pemberitahuan di layar.** → **P52**.

### C.3 Alamat layanan, autentikasi, dan penilaian berhasil

| Hal | Terbaca |
| --- | --- |
| Layanan | `ConnectREST\convertJsonNusareToProduction` — **14.283 B**, kelas `ASM-FW-GISFW-Work` |
| Alamat | `pyBaseURLSelectionType` = **`SETTING`**, `pyBaseURLSetting` = `LinkService!LinkService` |
| Setelan | `SystemSettings\LinkService` — `pySetting` = **`=ResponLink.URL`** ⭐ **ungkapan, bukan nilai tetap** |
| ⛔ **Autentikasi** | ⛔⛔ `pyUseAuthentication` = **`false`** |
| Batas waktu | **30.000 ms** |
| Dikirim | `.OfferFacIn.PolicyData.PolicyNo` · `.pzInsKey` · `.pxCreateDateTime` — dari **Clipboard ke JSON** |
| Diterima | dipetakan ke **`.StatusService`** |

⭐ **Pengisi alamatnya terbaca** `[terverifikasi]`: `Activity\GetLinkService.xml` — **57.104 B**,
kelas ⭐ **`ASM-FW-GISFW-Int-M_LINK_SERVICE`**, **4 langkah**: `Page-New` → **`Obj-Browse`** →
`Property-Set` → `Page-Remove`.
⭐ ⇒ **alamat dibaca dari sebuah tabel layanan saat berjalan** — ⭐ **sejalan dengan ADR-0013**
*(resolusi endpoint via `M_LINK_SERVICE`)*, ⛔ **bukan** ADR-0004 *(env var, superseded)*.
⭐ **Tidak ada pertentangan dengan ADR mana pun; nol ADR baru dibuat.**

⭐ **Apa yang dianggap berhasil** `[terverifikasi]` — `When\IsSuccessHitService.xml`:

> ⭐ **`pyWorkPage.StatusService.StsKonversiFacIn = 1`**

⚠️ **Namanya menyebut `FacIn`**, ⛔ padahal ada dua panggilan — FACIN **dan** FACOUT.
⭐ **Nol rule yang menguji status FACOUT terbaca.** → **P50**.

⭐ **Penanganan galat, dan naskahnya terbaca penuh** `[terverifikasi]`
`RDBList\UpdateErrorNoteJsonPolisMonitoring.xml`:

> `UPDATE POOLDATA.JSON_POLIS_MONITORING SET ERR_NOTE = {…ResponseMessage} WHERE NOPOLIS = {…PolicyNo}`

⭐ Ditambah `RDBList\INSERTJSON_JSONPOLISMONITORING_FACIN.xml` — ⚠️ **namanya pun menyebut FACIN
saja**.

---

## §D — Sasaran 4: alur

`[terverifikasi]` **Satu `Flow`: `Flow\InputAddendumTreatyIn.xml`, 202.269 B.**

| | NB Treaty In | ⭐ EDM Treaty In |
| --- | --- | --- |
| Nama | `InputRealizationTreatyIn` | ⭐ **`InputAddendumTreatyIn`** |
| Kelas kerja | `ASM-FW-GISFW-WORK` | ⭐⭐ **`ASM-FW-GISFW-WORK-ENDORSEMENTTREATY`** |
| Aktivitas mulai | `Start1` | `Start1` |
| Kotak | 31 | ⭐ **36** — Decision **8** · Assignment **4** · Utility **2** · Connector **22** |
| Sambungan | 32 | **24** |
| ⭐ **Antrean** | **lima** disebut | ⭐⭐ **TIGA** — `Admin` · `SecHead` · `DeptHead` |

⭐⭐ **Dua temuan dari perbandingan ini:**

1. ⭐ **Kelas kerjanya BERBEDA** — `…WORK-ENDORSEMENTTREATY`, sebuah **subkelas tersendiri**.
   ⭐ ⇒ endorsemen adalah **jenis kasus sendiri**, ⛔ bukan status di dalam kasus realisasi.
2. ⭐⭐ **Alur EDM hanya menyebut tiga antrean** — ⛔ **nol** `GroupLeader`, **nol** `Director`.
   ⭐ **Keputusan work owner untuk NB *(tangga tiga jenjang)* sudah menjadi kenyataan di EDM.**

⚠️ **Empat Assignment lawan enam di NB**, dan **delapan Decision lawan dua belas**.
⛔ **Tahap apa yang ada di EDM tetapi tidak di NB belum saya petakan kotak demi kotak** — ⭐ nama
kotaknya tidak terbaca dari medan yang saya sisir; `[terbuka]`.

⭐ **Di mana panggilan ke sistem luar duduk:** ⛔ **tidak terbaca dari `Flow`**. `[terverifikasi]`
Satu-satunya pemanggil `convertJsonNusareToProduction` adalah
`Activity\serviceInsertArasapas_act.xml`, ⭐ dan alur memuat **2 Utility** — `[dugaan]` salah
satunya memanggilnya, ⛔ **belum terverifikasi**.

---

## §E — Sasaran 5: sembilan puluh berkas EDM-saja

### E.1 Sebaran menurut tipe

`[terverifikasi]` `Activity` **28** · `RDBList` **22** · `Section` **15** · `DataTransform` **8** ·
`When` **5** · `FlowAction` **3** · `Harness` **3** · `ReportDefinition` **3** · `ConnectREST` **1** ·
`Flow` **1** · `SystemSettings` **1** = **90**.

### E.2 ⭐ Sebaran menurut PERAN

| Kelompok | Berkas | Isinya |
| --- | ---: | --- |
| ⭐ **EDM inti** | **31** | pembuatan, pembatalan, penomoran, penyimpanan endorsemen |
| ⭐ **adendum & premi tambahan** | **13** | layar, alur, transformasi, penomoran adendum |
| ⭐ **data lama / baru / selisih** | **9** | lima layar + empat aktivitas penghitung |
| ⭐ **layanan luar & pemantauan** | **6** | REST, setelan alamat, tabel pemantauan, uji berhasil |
| spreading & XOL | **1** | `Activity\FillSpreading` |
| lain | **30** | termasuk `SendEmailWithAttachments`, `ProtectionNonProp_Act`, `CheckNopolisAvailability` |

### E.3 ⭐⭐ Jangkauan — dan hasilnya berlawanan dengan NB

⭐ **Dihitung dua cara**, sama seperti NB:

| Cara | Terjangkau | Yatim |
| --- | ---: | ---: |
| ⭐ **A — urai langkah `Call`** dari 34 titik masuk | ⭐ **66** | ⭐ **0** |
| **B — nama aturan sebagai teks** | **66** | **0** |
| ✅ **sepakat** | | |

⭐⭐ **Nol aturan yatim, dari 66 `Activity`.** ⛔ **Berlawanan tajam dengan NB Treaty In**, yang
punya **28 yatim dari 92**.

⭐ **Bacaannya, dan ia penting untuk perencanaan:** ⭐ **folder EDM Treaty In adalah closure yang
RAPAT** — hampir seluruh isinya benar-benar dipakai. ⛔ Berbeda dari folder NB, yang memuat
keluarga aturan milik modul Fac. ⇒ ⭐ **lingkup EDM tidak akan menyusut seperti NB menyusut 60 %.**

⚠️ **Lingkup penelusuran dinyatakan:** titik masuk = **34** `Activity` yang namanya muncul di
`Flow`, `Section`, `Harness`, `FlowAction`, atau `DataTransform` **milik EDM**; ⛔ tidak disisir
di luar modul.

---

## §F — Empat aturan yang berbeda dari NB, diperiksa satu per satu

### F.1 ⭐ `CountSpreading_Act` — satu langkah lebih banyak, dan gerbangnya terbaca

| | NB | EDM |
| --- | ---: | ---: |
| byte | 86.087 | **91.344** |
| kelas | `Data-PolicyTreatyIn` | **sama** |
| langkah | 7 | ⭐ **8** |
| selisih | — | ⭐ **satu `Property-Set` lebih banyak** |

⭐⭐ **Catatan langkah yang HANYA ada di EDM, dikutip apa adanya:**

> ⭐ *"JIKA SPREADINGLIST 1 DAN SplitRNMSharePct NULL || 0"*

⭐ **Bacaannya:** langkah tambahan itu menangani keadaan **daftar penyebaran berisi satu baris**
dan **persentase pembagian kosong atau nol**. ⛔ **Apa yang ditetapkannya tidak terbaca** — nilai
`Property-Set` tidak terekspor. → **P60**.

### F.2 ⭐⭐ `SetCategoryAttach` — EDM adalah PEMBUNGKUS yang memanggil versi NB

| | NB | EDM |
| --- | ---: | ---: |
| byte | 87.436 | ⭐ **53.522** *(lebih kecil)* |
| kelas | `Data-OfferFacIn` | ⭐ **`Work`** |
| langkah | 7 | ⭐ **4** |

⭐⭐ `[terverifikasi]` **Langkah EDM memuat `Call ASM-FW-GISFW-Data-OfferFacIn.SetCategoryAttach`**
— ⭐ **ia memanggil versi berkelas NB.** ⭐ Dua catatan langkahnya: *"Get list for spreading"* dan
*"Set Spreading into Temp Page"*.

⭐ ⇒ **EDM menambahkan penanganan penyebaran, lalu mendelegasikan sisanya.** ⛔ **Bukan dua
implementasi yang bersaing** — melainkan **lapisan di atas lapisan**.

### F.3 ⭐⭐ `serviceInsertArasapas_act` — lihat §C

⭐ **Inilah implementasi yang menahan NB.** 1 langkah → **20 langkah**; kelas
`Data-PolicyTreatyIn` → **`Work`**. ⭐ Seluruh rinciannya di **§C**.

### F.4 ⚠️ `IsUW` — kelas berbeda, syarat tampak SAMA, tetapi jumlah syaratnya tidak

| | NB | EDM |
| --- | ---: | ---: |
| kelas | `Work` | ⭐ **`Data-OfferFacIn-LocationReinsurance`** |
| `pyConditionString` | `…pyWorkBasketList(1).pyWorkBasketName = ReasFacInGroupLeader` | ⭐ **identik** |
| ⭐ wadah syarat | **4** | ⭐ **5** |

⚠️⚠️ **Keterangannya sama, tetapi EDM punya SATU WADAH SYARAT LEBIH BANYAK.** ⭐ Dan catatan
pengembang EDM berbunyi *"add conditon not show if Underwriting upper"* — ⭐ **menegaskan ada
syarat tambahan**.

⛔ **Syarat tambahan itu tidak muncul di `pyConditionString`** — ⭐ ini **persis pola yang sudah
diputuskan work owner di P23**: *"bila keterangan berbeda dari yang dijalankan, yang dijalankan
benar"*. ⭐ **Tidak perlu ditanyakan ulang** — ⛔ tetapi **wajib dibaca dari bentuk terstrukturnya**
saat spec ditulis.

⚠️ Ditambah: pola ini memakai **penunjukan antrean menurut nomor urut** — ⭐ sudah diputuskan di
**P25**, **tidak dimigrasi**.

---

## §G — Bukti isolasi

| Yang dibuktikan | Keadaan |
| --- | --- |
| ⛔ `D:\XML\nusantara-re\` | ✅ **NOL disentuh** |
| korpus | ✅ **READ-ONLY** — hanya dibaca |
| ⭐ `NB Treaty In` dibaca | ⚠️ **ya, sebagai pembanding** — ⭐ **brief mengizinkannya** |
| ⛔ temuan NB masuk sensus EDM | ✅ **NOL** — §0 menyebut jendelanya setiap kali |
| berkas dibuat | ✅ **tepat DUA**, sesuai brief |
| kode · DDL · usulan kolom | ✅ **NOL** |
| ⛔ nilai nama orang | ✅ **NOL disalin** |
| ⛔ nomor polis apa adanya | ✅ **NOL** |
| butir terbuka ditutup | ✅ **NOL** |
| ADR baru | ✅ **NOL** — ⭐ ADR-0013 dirujuk, tidak dibuat ulang |

⚠️ **Baris tidak bersih:** sesi ini panjang dan pernah disambung; harness pernah menyuntik ulang
bacaan `.scratch` konteks lain. ⛔ **Nol nama tabel, kolom, atau pola rancangan konteks lain
merembes** — seluruh angka dari sisiran korpus EDM sendiri dan pembandingnya NB.

---

## TELEMETRI EKSEKUSI

⛔⛔ **Pengukuran dari luar TIDAK dilakukan.** ⭐ Ronde ini berjalan **di dalam sesi interaktif**,
sehingga **biaya dan durasi tidak tersedia**. ⛔ Tidak ditaksir.

| Yang dicatat | Nilai | Cara |
| --- | ---: | --- |
| ⭐ **token keluaran** | ⭐ **145.564** | transkrip `.jsonl`, SELISIH terhadap baseline ronde ini |
| ⭐ **token cache-read** | ⭐ **37.448.906** | sama |
| token cache-write | **188.751** | sama |
| ⭐ **panggilan alat** | ⭐ **42** | sama |
| ⛔ durasi | ⛔ **TIDAK DAPAT DIUKUR** | butuh `--print` dari luar sesi |
| ⛔ biaya | ⛔ **TIDAK DAPAT DIUKUR** | butuh `--print` dari luar sesi |
| ⭐ **byte dibaca** | ⭐ **17.639.246** *(163 berkas EDM, diurai penuh)* ditambah **35.492.317** *(163+278 berkas NB sebagai pembanding)* | perintah audit di tiap bab |

⭐ **Baseline ronde ini diambil SEBELUM ronde dimulai** — ⭐ jadi angkanya **terukur langsung**,
⛔ **bukan hasil pengurangan** seperti dua ronde sebelumnya.

⚠️ **Tiga keterbatasan pengukuran dari dalam:** *(a)* pesan terakhir belum tertulis ke transkrip
saat pengukuran kedua — **kurang satu pesan**; *(b)* jalur transkrip **diperiksa**, cocok dengan
sesi ini; *(c)* yang dilaporkan **SELISIH**, bukan total sesi.

⭐ **Bandingkan dengan seluruh ronde sebelumnya** — selisih, alat yang sama:

| | NB R1 | NB R2 | NB R3 | NB R4 | Verifikasi | Spec | Tiket | ⭐ **EDM R1** |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| keluaran | 192.021 | 225.898 | 178.857 | 154.518 | 92.996 | 115.495 | 111.654 | **145.564** |
| cache-read | 13,5 jt | 22,7 jt | 16,7 jt | 26,2 jt | 20,8 jt | 15,6 jt | 25,5 jt | ⚠️ **37,4 jt** |
| panggilan | 69 | 72 | 39 | 47 | 35 | 21 | 31 | **42** |

⚠️⚠️ **Cache-read TERTINGGI dari seluruh ronde — 37,4 juta.** ⭐ Sebabnya terukur dan wajar:
ronde ini mengurai **dua modul sekaligus** — 163 berkas EDM **dan** 278 berkas NB sebagai
pembanding — ditambah membaca berkas keadaan NB. ⭐ Perbandingan lintas-modul memang mahal, dan
itu **harga yang dibayar sekali** untuk 69 aturan yang kini tidak perlu diurai ulang.

