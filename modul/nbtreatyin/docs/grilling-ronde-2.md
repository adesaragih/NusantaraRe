# NB Treaty In — Grilling Ronde 2
## Yang tidak terekspor: 942 langkah penetapan nilai tanpa satu pun nilainya

Konteks: **Treaty Inward — Realisasi & Endorsement**. Modul: **hanya `NB Treaty In`**.
Korpus READ-ONLY `D:\XML\RNM_BRD\NB Treaty In\`. Aturan: `CLAUDE.md` §4 dan §4a.
Dasar: `.scratch\nb-treaty-in\grilling-ronde-1.md`.

⛔ **Ronde ini MENEMUKAN dan BERTANYA.** Nol keputusan rancangan · nol kode · nol DDL ·
nol daftar kolom usulan · nol nomor baris XML · ⛔ **nol butir terbuka yang saya tutup sendiri** ·
⛔ **nol nilai nama orang disalin**.

⛔ `grilling-ronde-1.md` **TIDAK disunting.** Tiga ralat atasnya ditulis **di sini**, dengan
kalimat lamanya dikutip utuh.

> **SENSUS BERKAS INI**
>
> ⭐ **Jendelanya: berkas ini MINUS blok sensus ini sendiri.** ⛔ Disebut terang-terangan supaya
> blok ini tidak mengubah angka yang ia klaim.
>
> **1023 baris di jendela** · **12 bab** `## ` · **33 sub-bab** `### ` ·
> **11 sub-sub** `#### ` · **31 tabel** · **22 pernyataan berpenanda
> terverifikasi** · **6 berpenanda dugaan** · ⚠️ **0 berpenanda data DBA** — ⭐ nol
> fakta basis data BARU dipakai di ronde ini, dan itu benar: ronde ini membaca korpus, bukan DBA.
>
> ⭐ **Register Bab C.2: 13 butir baru — DUA CARA.** *(a)* baris tabel ⇒ **13**; *(b)* nomor di
> kolom pertama ⇒ **1–13 tanpa nomor hilang**. ✅ **Sepakat.**
> ⭐ **Pemilik — DUA CARA** *(baris yang memuat pemilik, dan kemunculan literalnya)*:
> `[pengembang Pega lama]` **8** · `[Finance]` **3** · `[Product+Underwriting]` **3** ·
> `[pemilik export Pega]` **2** · `[IAM]` **1** · `[DBA]` **0** · `[work owner]` **0**.
> ✅ **Sepakat.** ⚠️ **Sebutan 17 lawan butir 13** — ⭐ sebab **empat butir berpemilik ganda**.
> ⛔ **Nol butir ditutup.**
>
> ⭐ **Bab D — DUA CARA:** *(a)* pola `**Pn ` ⇒ **11**; *(b)* judul `#### ` ⇒ **11**, bernomor
> **P18–P28 tanpa nomor hilang**. ✅ **Sepakat.**
>
> ⚠️ **Ungkapan terlarang muncul hanya sebagai SEBUTAN**, di kalimat yang justru menyatakannya
> nol, dan di perintah audit yang membuktikan nol. ⛔ Dinyatakan di sini apa adanya.
>
> ⛔ **Nol keputusan rancangan · nol kode · nol DDL · nol daftar kolom usulan · nol nomor baris
> XML · nol butir ditutup sendiri · nol nilai nama orang · nol nomor OQ dikarang sendiri.**

---

## §0 — Pencocokan dan ujian instrumen

### 0.1 Delapan angka ronde 1 dicocokkan — ⭐ tujuh COCOK, satu perlu keterangan

| Yang dicocokkan | Harap | Terukur | |
| --- | --- | --- | :---: |
| berkas `.xml` | 278 | **278** | ✅ |
| md5 korpus *(isi disambung, urut path)* | `62a3735e…` | **`62a3735ebb2eabb788c8cd8abdda78ca`** | ✅ |
| tipe rule, dua cara | 10 tipe | **sepakat baris per baris** | ✅ |
| `RDBList` ber-`pyBrowseSQL` | 41/41 | **41/41** | ✅ |
| blok PL/SQL `BEGIN…END` | 5 | **5** | ✅ |
| stored procedure dipanggil | 3 | **3** | ✅ |
| `Section`+`Harness` byte | 16.635.834 · 46,6 % | **16.635.834 · 46,6 %** | ✅ |
| berkas pola hitung uang | 30 | **30** *(29 `Activity` + 1 `When`)* | ✅ |

Perintah audit byte: `total 35.678.284 B`, `Section+Harness 16.635.834 B`, `= 46,6 %`.
Perintah audit hitung uang — persis seperti brief:
`find Activity When -name '*.xml' | grep -icE 'count|sum|total|calc|premi|rate|limit|share|comm|brokerage|tax'`

⚠️ **Satu daftar di brief bukan daftar terbesar, dan itu perlu disebut.** Brief menyebut
`CheckSpreadingProtectAnekaGolf_ACT.xml` **398.150 B** sebagai anggota "yang terbesar". Terukur,
lima `Activity` terbesar sebenarnya: `SumTSIPremiSpreadedRNM_FIRE_Act` **996.052** ·
`ProtectFIREMBUPA_Act` **875.119** · `SumTSIPremiSpreadedRNM_ANEKA_Act` **748.216** ·
`SumTSIPremiSpreadedRNM_Act` **516.395** · ⭐ `InputPolicyTreatyInDetail_preACT` **488.903**.
⭐ Jadi daftar brief adalah **pilihan tangan atas berkas hitung/proteksi besar**, bukan lima
terbesar. ⛔ **Bukan alasan berhenti** — keempat berkas pertamanya benar.

### 0.2 ⛔⛔ UJIAN INSTRUMEN — instrumen saya gagal EMPAT KALI di ronde ini

⭐ `CLAUDE.md` §4a menuntut tiap alat diuji atas butir berjawaban-diketahui. ⛔ **Empat kali alat
saya salah, dan keempatnya tertangkap sebelum masuk berkas ini.** Ditulis apa adanya, bukan
diperbaiki diam-diam.

| # | Yang alat saya katakan | Yang benar | Sebab galatnya |
| --- | --- | --- | --- |
| ⛔ **1** | `pyRequired` ada **2.190** ⇒ 2.190 medan wajib | ⛔ **160 berisi**, dan **1 nilai unik** | ⛔ Pega menulis **tiap tag untuk tiap sel**, mayoritas kosong. Kehadiran tag **bukan** isi |
| ⛔⛔ **2** | ambang uang ter-hardcode: `pxCreateDateTime > 20211227` **1.216×**, `REPEATINGINDEX = "1"` **3.568×** | ⛔ **NOL**. Keduanya **markup ekspor XML** — `pxCreateDateTime` metadata aturan, `REPEATINGINDEX` atribut `rowdata` | ⛔ saya menyisir **teks mentah berkas**, bukan medan pembawa ekspresi |
| ⛔ **3** | *(dugaan yang hampir saya tulis)* `pyConditionString = "[Double click to add condition]"` ⇒ **11 penggolong KOSONG** | ⛔ **NOL kosong.** **75 dari 75** punya syarat terbaca di bentuk `[properti][operator][nilai]` | ⛔ saya membaca **satu** medan dan menyimpulkan tentang aturannya |
| ⚠️ **4** | 19 objek Oracle baru di luar `RDBList` | ⛔ **NOL.** Kandidatnya `THE` **74×**, `FULLY` **68×**, `PROPERTY`, `MATH` — ⭐ **kata prosa Inggris** dari teks dokumentasi | ⚠️ pola `FROM <nama>` memungut prosa |

⭐ **Dan dua ujian yang alat saya LULUSI**, dilaporkan bersama yang gagal:

| Butir uji | Jawaban yang sudah diketahui | Hasil alat | |
| --- | --- | --- | :---: |
| `pyAssociatedPrivileges` **295 kemunculan** di 29 berkas layar, sementara **OQ-007** menyatakan korpus tanpa aturan wewenang | harusnya **nol berisi** | ⭐ **0 berisi** | ✅ |
| Tujuh `When` dari jalur alur — apakah sampai ke layar? | `isApproved` dipakai gerbang pertama | ⭐ `isApproved` di **5 berkas layar**; enam lainnya **0** | ✅ |

⛔⛔ **Aturan baca yang lahir dari galat #1 dan #2, dan ia mengikat ronde berikutnya:**
⭐ **Di ekspor Pega, cacah apa pun WAJIB atas MEDAN yang membawa isi, dan WAJIB membuang yang
kosong.** ⛔ Menyisir teks berkas memungut metadata aturan dan atribut ekspor, dan **inflasinya
mencapai 3.568 lawan nol** — bukan meleset sedikit, melainkan **seluruhnya palsu**.

### 0.3 ⚠️ Sisa editor — dan di modul ini perannya berbeda dari dugaan

`[terverifikasi]` `pyExpressionGadget` memuat `<pxObjClass>PegaGadget-ExpressionBuilder</pxObjClass>`
— ⭐ **wadah sisa editor**, **1.568 kemunculan** di 30 berkas hitung.

⭐ **Diuji dua arah, dan hasilnya mengejutkan ke arah sebaliknya:**

| Medan | Dengan sisa editor | Tanpa | Inflasi |
| --- | ---: | ---: | :---: |
| `pyValue` | 323 | **323** | **1,00×** |
| `pyStepsPreCondition` | 342 | **342** | **1,00×** |
| ⛔⛔ `PropertiesValue` *(30 berkas hitung)* | **29** | ⛔ **0** | ⛔ **tak terhingga** |
| ⛔⛔ `PropertiesValue` *(92 `Activity`)* | **74** | ⛔ **0** | ⛔ **tak terhingga** |

⭐ **Bacaannya:** gerbang langkah **terekspor normal**; ⛔ **nilai yang ditetapkan langkah TIDAK** —
ia hanya bertahan di **salinan editor**. Lihat **§1**.

---

## §1 — ⛔⛔ Temuan terberat: 942 langkah menetapkan nilai, dan NOL nilainya terekspor

⭐ **Ini bukan "saya belum membacanya". Ini "nilainya tidak ada di ekspor".**

`[terverifikasi]` Di **92 `Activity`**:

| | Jumlah |
| --- | ---: |
| langkah ber-metode **`Property-Set`** | ⭐ **942** |
| `PropertiesValue` yang ada **termasuk salinan editor** | **74** |
| ⛔ `PropertiesValue` pada **aturan hidup** *(salinan editor dibuang)* | ⛔ **0** |

Perintah audit: cacah `pyStepsActivityName == "Property-Set"` sesudah
`<pyExpressionGadget>…</pyExpressionGadget>` dibuang; lalu cacah
`REPEATINGINDEX="PropertiesValue"` dengan dan tanpa blok itu.

⭐ **Diperiksa apakah nilainya bersembunyi di tempat lain — TIDAK.** `[terverifikasi]`
`pyStepsCallParams` berisi **536**, tetapi isinya **hanya placeholder** —
`<pyTempPlaceHolder>TempPlaceHolder</pyTempPlaceHolder>` dan `<pyStepsParamUI/>` kosong.
⛔ **Nol nilai bisnis di situ.**

⚠️ **Satu-satunya tempat rumus uang masih terbaca adalah salinan editor**, dan ⛔ **menurut
aturan baca proyek ini salinan editor BUKAN aturan yang dieksekusi.** Bentuk yang terlihat di
sana, dikutip satu kali sebagai contoh bentuk **bukan sebagai fakta perilaku**:
`divide(.GrossPremium * pyWorkPage.TreatyIn.RNMShareP,100.…` — ⛔ **terpotong di tengah**, dan
itu memperjelas kenapa ia tidak dapat dipakai.

⛔⛔ **Akibatnya, dan ia sekelas dengan blocker OQ-025:** **rantai hitung uang modul ini TIDAK
DAPAT DIREKONSTRUKSI dari ekspor ini.** ⭐ Bukan sulit — **tidak mungkin**, dengan bahan yang ada.

⚠️ **Yang MASIH terbaca tentang langkah-langkah itu**, dan hanya ini:

| Yang terbaca | Jumlah | Nilainya |
| --- | ---: | --- |
| metode yang dipanggil | **388** `Property-Set` di 30 berkas hitung | ⭐ tahu **bahwa** ia menetapkan, ⛔ tidak tahu **apa** |
| `RDB-List` | 15 | jembatan ke SQL |
| `Page-Set-Messages` · `Property-Set-Messages` | 12 · 10 | ⭐ **22 tempat memasang pesan galat** |
| aktivitas yang dipanggil | ⭐ **6 varian `SumTSIPremiSpreadedRNM_*`** | `_FIRE` · `_ANEKA` · `_MBU` · `_MARINECARGO` · `_GOLF` · `_PA` |
| `Obj-Save` | 1 | satu titik simpan |

### 1.1 ⭐ Catatan pengembang menyelamatkan sebagian arti — dan ia menjawab P14 separuh

`[terverifikasi]` `pyStepsDescription` **316 berisi**, dan sebagiannya **menuliskan arti slot
generik secara eksplisit**:

| Catatan langkah, dikutip apa adanya | Yang ia beri tahu |
| --- | --- |
| ⭐ `Deduction (CARI44[deduction])` | slot **44** = pengurangan |
| ⭐ `Net Premi (CARI31[CURRENCY], CARI39, CARI47[BALANCEDUETO], CARI13[NET PREMIUM])` | slot **31** = mata uang · **47** = saldo terutang · **13** = premi neto |
| ⭐ `Gross Premi (CARI32, CARI31[CURRENCY])` | slot **32** = premi bruto |
| `Gross Premi Retro (CARI32, CARI31[CURRENCY])` | ⭐ slot **sama** dipakai jalur retro |

⭐⭐ **Ini menjawab P14 sebagian, dari korpus, tanpa menunggu siapa pun** — ⭐ dan **memperkuat**
butir baru #1 ronde 1: arti slot **ditentukan per rule**, dan di sini pengembangnya **menuliskan
pemetaannya di dalam catatan langkah**. ⛔ **Tetapi hanya untuk sebagian slot**, dan ⛔ **catatan
bukan aturan** — ia dapat basi tanpa ada yang tahu.

⚠️ **Dan dua catatan langkah membawa aturan bisnis yang tidak ada di mana pun selain di situ:**

| Catatan | Apa yang ia nyatakan |
| --- | --- |
| ⚠️ `set protect kalau master polis premi harus 0` | ⛔ sebuah **kewajiban nilai nol** |
| ⚠️⚠️ `TOTAL PREMI SHARE CEDANT HARUS = PREMI RNM` | ⛔ sebuah **invarian keseimbangan uang** |

⛔ **Keduanya hanya hidup sebagai catatan.** ⚠️ Apakah kode benar-benar menegakkannya **tidak
terbaca**, sebab langkah penegaknya adalah `Property-Set` yang nilainya tidak terekspor.

---

## §2 — Sasaran 1: 31 berkas layar, kini dibuka

⛔ **Tidak satu pun dibaca utuh** — aturan §4 butir 4 dipatuhi: pola dulu, rentang seperlunya.
⭐ Yang saya ekstrak: **medan wajib · syarat tampil · syarat hanya-baca · syarat mati · tombol ·
rujukan penggolong**. ⛔ **Nol tata letak, nol urutan kolom, nol gaya.**

### 2.1 ⭐ Vonis per-berkas — dan ia MENJAWAB P17 dari korpus sendiri

`[terverifikasi]` **733 syarat berisi**, dan **229 bermakna** sesudah tautologi
*(`always`, `true`, `1=1`, `1==1`, `1 = 1`)*, kontradiksi *(`1=2`, `1==2`, `NEVER`)*, dan
`-1` / `Other Property` dibuang.

| Berkas layar | byte | berisi | ⭐ **BERMAKNA** |
| --- | ---: | ---: | ---: |
| ⭐ `GeneralPolicyTreatyIn` | 1.926.378 | 155 | ⭐ **59** |
| ⭐ `DetailPolicyTreatyIn` | 1.916.240 | 155 | ⭐ **59** |
| ⭐ `GeneralDeptHeadTreatyIn_UW` | 1.584.498 | 132 | **27** |
| ⭐ `DetailDeptHeadTreatyIn_UW` | 1.567.128 | 132 | **27** |
| `SFAPortalOpportunities` | 285.893 | 13 | 8 |
| `SFAPortalOpportunitiesHeader` | 133.582 | 10 | 6 |
| `SpreadingRiskList` · `SFAPortal_OpportunitiesList_Header` | 287.046 · 95.626 | 7 · 5 | 5 · 5 |
| `SOB` · `DetailPoliciesNonProportional` | 185.724 · 102.444 | 5 · 4 | 4 · 4 |
| ⚠️ **`DetailPolicyTreatyInNonProportional`** | ⚠️ **1.662.941** | 25 | ⚠️ **5** |
| ⚠️ `DetailPolicyTreatyOutNonProportional` | 978.997 | 14 | ⚠️ **3** |
| ⚠️ `DetailPolicyTreatyInNonProportionalEDM` | 915.926 | 14 | ⚠️ **1** |
| `ListSuggest` · `SFAPortal_OpportunitiesList` | 218.982 · 287.715 | 13 · 5 | 3 · 3 |
| `BusinessAndSOBList`/`Retro` *(Section)* · `HistoricalSurveyReport` · `HistoricalSurveyReportUW` · `SourceHierarki` | — | 6 · 6 · 4 · 4 · 3 | 2 tiap satu |
| ⛔ **11 berkas: MURNI TATA LETAK** | — | 0–4 | ⛔ **0** |

⛔ **Sebelas berkas murni tata letak**, dinyatakan satu baris seperti brief minta:
`BusinessAndSOBList` *(Harness)* · `BusinessAndSOBListRetro` *(Harness)* ·
`HistoricalSurveyReportDtl` · `HistoricalSurveyReportDtlUW` · `Installments_ReadOnly` ·
`InstallmentList` · `ShowPolicyNoTreaty_SC` · `PolicyTreatyInDeclineConfirm` ·
`InputHistoricalSurveyReportDtl` · `InputHistoricalSurveyReportDtlUW` ·
`SFAPortal_Opportunities`.

⭐⭐ **Jawaban P17 terbaca dari angka ini, bukan dari orang:** aturan bisnis layar **memusat di
`GeneralPolicyTreatyIn` dan `DetailPolicyTreatyIn`** — **59 syarat bermakna masing-masing, 118
dari 229 atau 52 %** hanya di dua berkas. ⚠️ Dan kejutannya: **`DetailPolicyTreatyInNonProportional`
berukuran 1,66 MB tetapi hanya 5 syarat bermakna** — ⭐ **berkas ketiga-terbesar modul hampir
seluruhnya tata letak.**

### 2.2 ⛔⛔ Sembilan puluh satu elemen layar MATI SEJAK DIBANGUN

`[terverifikasi]` Syarat yang **tidak mungkin benar**, dan jumlahnya:

| Bentuk | Kemunculan | Di mana |
| --- | ---: | --- |
| ⛔ `1=2` | **62** *(`pyCondition`)* + **10** *(`pyContainerVisibleWhen`)* | 11+8 berkas layar |
| ⛔ `NEVER` | **14** + **1** | `Detail`/`General` `PolicyTreatyIn`, `…NonProportionalEDM` |
| ⛔ `1==2` | **2** | `DetailDeptHeadTreatyIn_UW`, `GeneralDeptHeadTreatyIn_UW` |
| ⛔ `.IsNewPolicyNonProp != 1 && NEVER` | **2** | dua layar Kepala Departemen |
| ⭐ **TOTAL** | ⭐ **91** | — |

⭐ **Dan pasangannya, yang selalu benar:**

| Bentuk | Kemunculan | Medan | Akibatnya |
| --- | ---: | --- | --- |
| ⛔ `1==1` | **32** | `pyReadOnlyCondition` | ⛔ **hanya-baca permanen** |
| ⛔ `1=1` · `1 = 1` · `ALWAYS` | 2 · 2 · 2 | `pyReadOnlyCondition` | sama |
| ⭐ **TOTAL** | ⭐ **38** | — | ⛔ **38 elemen tak dapat disunting, selamanya** |

⛔⛔ **Dan seluruh 38 itu ada di DUA berkas saja** — `DetailDeptHeadTreatyIn_UW` dan
`GeneralDeptHeadTreatyIn_UW`, ⭐ **layar Kepala Departemen**. ⚠️ Artinya `[dugaan]` **layar
Kepala Departemen dibuat hanya-baca dengan cara dipaku**, bukan dengan syarat peran.

⚠️⚠️ **Tiga ejaan untuk tautologi yang sama** — `1==1`, `1=1`, `1 = 1` — dan **dua untuk
kontradiksi** — `1=2`, `1==2`. ⭐ Ini **jebakan sensus nomor 4 `CLAUDE.md` §4a di dalam korpus
sendiri**: mencari satu ejaan saja akan melewatkan sisanya.

### 2.3 ⛔ Nama orang di dalam syarat tampil layar — dicacah, TIDAK disalin

⛔ **Aturan 10 dipatuhi: nilai nama TIDAK ditulis.**

`[terverifikasi]` **12 kemunculan** di **4 berkas layar**, memuat **4 nama berbeda**:

| Berkas layar | Kemunculan | Medan yang dicocokkan |
| --- | ---: | --- |
| `DetailDeptHeadTreatyIn_UW` | 3 | `OperatorID.pyUserIdentifier` |
| `GeneralDeptHeadTreatyIn_UW` | 3 | `OperatorID.pyUserIdentifier` |
| ⭐ `ListSuggest` | **4** | `OperatorID.pyUserIdentifier` |
| ⚠️ `DetailPoliciesNonProportional` | 2 | ⚠️ **`OperatorID.pxInsName`** |

Perintah audit: `(pyUserIdentifier|pxInsName|pxCreateOperator)\s*[!=]=?\s*'[A-Z_][A-Z_0-9]{3,}'`
atas 31 berkas layar, teks ter-unescape dua kali.

⚠️ **Dua hal yang memperburuknya dibanding OQ-021:**
1. ⛔ Nama dipakai **dua arah** — ada syarat `== nama` **dan** `!= nama` di berkas yang sama,
   ⭐ sehingga **menghapus orangnya mengubah perilaku pada DUA cabang**.
2. ⚠️ Satu berkas mencocokkan **`pxInsName`**, medan yang berbeda dari yang disapu OQ-021
   *(`pyUserIdentifier`/`pyUserName`/`pyPosition`)* ⛔ **dan berbeda dari OQ-027** *(`pyTelephone`)*
   — ⭐ **pola guard KEEMPAT**.

### 2.4 Validasi yang hanya hidup di layar

`[terverifikasi]` **Medan wajib isi:**

| Medan | Berisi | Nilai |
| --- | ---: | --- |
| `pyRequired` | **160** | `true` — **1 nilai unik** |
| `pyRequiredNew` | **84** | ⚠️ `always` **47** + `true` **37** — ⭐ **dua ejaan** |
| ⭐ `pyRequiredWhen` | **37** | ⭐ **5 syarat sungguhan** |
| `pyParameterRequired` | 48 | `-1` |

⭐ **Kelima syarat wajib-bersyarat, dikutip:**

| Syarat | Kali | Layar |
| --- | ---: | --- |
| ⭐ `pyWorkPage.Quotation.ProportionalType != 'NonProportional'` | **26** | `Detail`+`GeneralPolicyTreatyIn` |
| `.QuotationData.ProportionalType = 'Proportional'` | 4 | sama |
| ⚠️ `.Claim != '' && .Claim != 0` | 4 | sama |
| `.FlagPPH = true` | 2 | sama |
| ⛔ `.IsApproved == 1 && (pyUserIdentifier == …nama…)` | 1 | `ListSuggest` |

⛔⛔ **`.Claim != '' && .Claim != 0` adalah pola yang brief peringatkan, dan ia ADA di sini:**
⭐ satu medan bernuansa uang dibandingkan **terhadap teks kosong DAN terhadap angka nol** dalam
satu syarat. ⚠️ Artinya `[dugaan]` medan itu **kadang teks, kadang angka** — dan itu menentukan
apakah "belum diisi" dapat dibedakan dari "diisi nol".

⛔⛔ **Dan hal yang sama menimpa gerbang PERTAMA seluruh alur:**

| Ejaan di layar | Kali | Bentuknya |
| --- | ---: | --- |
| ⛔ `.IsApproved = '0'` | 2 | ⛔ **TEKS berkutip** |
| ⛔ `.IsApproved==0` | 2 | ⛔ **ANGKA telanjang** |
| `.IsApproved==1` · `.IsApproved == 1` | 2 · 5 | angka |

⭐ Dan **varian `When`-nya** menguji `pyWorkPage.PolicyTreatyIn.IsApproved = 1`.
⛔ **Jadi satu bendera yang sama dibandingkan sebagai teks di satu tempat dan angka di tempat
lain** — tepat peringatan brief, di medan yang paling menentukan.

### 2.5 Penyaring tampilan: penggolong, jabatan, dan satu yang rapuh

`[terverifikasi]` Penggolong yang dipakai **sebagai penyaring layar**: `IsUW` · `IsClaim` ·
`IsOperatorLife` · `TreatyMasterInEDM` · `IsNotAdmin` · `isSellingModeB2B` / `B2C` / `B2BB2C` ·
`crmCreateOpportunity` · `crmIsReview` · `pyIsIpadOrDesktop`.

⭐ **Dua tempat yang menyebut jabatan sebagai TEKS, dan ini menyentuh OQ-024:**

| Syarat | Layar | Kenapa penting |
| --- | --- | --- |
| ⭐ `pyWorkPage.PositionNote != 'ReasTreatyInSecHead'` | `ListSuggest` | ⭐⭐ **nama antrean terbaca sebagai teks di sebuah syarat** — OQ-024 menyebutnya *"tidak terbaca"*; ⭐ **di sini SATU pemetaan terbaca** |
| ⚠️ `TempEmail.CARI28 != pyWorkPage.PositionNote` | `Detail`+`GeneralPolicyTreatyIn` *(4×)* | ⛔ **slot generik menggerbangi LAYAR**, bukan hanya SQL — memperluas butir baru #1 ronde 1 |

⛔⛔ **Satu syarat yang rapuh secara struktural:**
`OperatorID.pyWorkGroup!='ReasLife' && OperatorID.pyWorkBasketList(2).pyWorkBasketName=='…'`
⭐ dan di penggolong: `pxRequestor.OperatorID.pyWorkBasketList(1)` **3×**.
⛔ **Keduanya menunjuk antrean menurut POSISI di dalam daftar** — ⚠️ posisi 1 dan posisi 2.
⛔ Menambah atau mengurutkan ulang antrean seorang pengguna **mengubah perilaku** tanpa ada yang
menyentuh aturan.

### 2.6 Tombol dan skrip peramban

`[terverifikasi]` **938 aksi berisi, 18 jenis.** Yang membawa perilaku:
`runActivity` **43** · `runScript` **44** · `finishAssignment` **15** · `createWork` **12** ·
`localAction` **12** · `showHarness` **16** · `save` **9** · `openWorkByHandle` **2** ·
`deleteRow` **10** · `addRow` **22**.

⛔ **`pyFunctionName` — 44 berisi, dan keduanya JavaScript peramban mentah:**
`script:window.close` **24×** · `script:opener.location.reload` **20×**.
⭐ `[dugaan]` arsitekturnya **jendela sembul**: satu layar membuka jendela lain, lalu
**memaksa jendela induk memuat ulang**. ⛔ Apa yang gagal bila muat-ulang itu tidak terjadi
**tidak terbaca**.

⛔ **`pyDisableSubmit` 397 kemunculan ⇒ 0 berisi** — ⭐ tidak satu pun pengiriman diblokir dari layar.

⚠️ **Satu nilai yang tampak salah tulis:** `pyDisabledNew` bernilai ⚠️ **`truewhn`** **4×**
*(di samping `true` 17× dan `always` 12×)*. ⭐ `[dugaan]` **salah ketik dari `truewhen`**.
⛔ Apakah Pega mengabaikannya — sehingga penonaktifan **diam-diam tidak berlaku** — **tidak dapat
saya nyatakan** dari korpus.

### 2.7 Pembeda proporsional di modul ini

`[terverifikasi]` Medan pembedanya **`ProportionalType`**, bernilai `'Proportional'` /
`'NonProportional'`. ⚠️ **Dibawa DUA halaman berbeda**:
`pyWorkPage.Quotation.ProportionalType` *(14+26 kemunculan)* dan
`.QuotationData.ProportionalType` *(4+4+2)*. ⛔ **Apakah keduanya satu nilai yang disalin, atau
dua nilai yang dapat berbeda, tidak terbaca.**

⭐ Pembeda lain yang terbaca dari layar: `.IsNewPolicyNonProp` · `.IsNewPolicyListFormat` ·
`.ClaimType = 'XOL Retro'` · `.TreatyType='XOL'` · `.FlagPPH` · `.BalanceDueTo >= 0` / `< 0` ·
`pyWorkPage.TreatyIn.FacultativeShare = 0` / `!= 0` / `> 0` · `TreatyIn.ViewState = 1` ·
`pyWorkPage.FlagViewPolicy = 1` · `pyWorkPage.IsOldData = 1` ·
`pyWorkPage.PolicyTreatyIn.IsEDMInputOnNB` · `InputParam.CARI12 != 'treaty' && != 'TreatyPolicy'`.

⚠️ `InputParam.CARI12` — ⭐ **slot generik lagi**, kali ini dibandingkan dengan **dua teks harfiah
yang berbeda hanya pada huruf besar-kecil dan sisipan**: `'treaty'` dan `'TreatyPolicy'`.

---

## §3 — Sasaran 3: dua DecisionTable, dan RALAT atas ronde 1

### 3.1 ⛔⛔ Baris keputusan KEDUA DecisionTable tidak terekspor

`[terverifikasi]` Dibaca utuh — keduanya kecil:

| Berkas | byte | Identitas 4 bagian | Baris keputusan |
| --- | ---: | --- | --- |
| `DecisionTable\isApproved.xml` | 16.902 | `ASM-FW-GISFW-WORK!ISAPPROVED` · `GISFW` · **01-01-52** · `20170608T025010.808` | ⛔ **NOL** |
| `DecisionTable\BusinessType_DeT.xml` | 36.776 | `ASM-FW-GISFW-WORK!BUSINESSTYPE_DET` · `GISFW` · **01-01-87** · `20250922T082411.240` | ⛔ **NOL** |

⛔ **Nol `pyCriteriaValue`, nol `pyResult`, nol `pyReturnValue`, nol `pyPropertyName`, nol
`pyOtherwiseResult`.** ⭐ Yang bertahan hanya `pyLabel` berisi nama aturan dan daftar versi.
⛔ **Persis pola yang brief peringatkan** *(sudah terjadi pada 49 DecisionTable lain di korpus)*.
⛔ **Saya TIDAK menebak isinya.**

⚠️⚠️ **Dan `BusinessType_DeT` lebih buruk dari sekadar tak terbaca:** ⭐ `pxUpdateDateTime`-nya
**22 September 2025** — ⭐ **aturan paling baru diperbarui yang saya lihat di modul ini** — dan
`pyLabel`-nya mendaftar **delapan versi sebelumnya** *(01-01-87, 77, 74, 73, 55, 54, 53, …)*.
⛔ **Jadi ia penggolong yang HIDUP dan SERING DIUBAH, yang isinya tidak terlihat sama sekali.**

### 3.2 ⭐ Varian `When`-nya justru TERBACA PENUH

`[terverifikasi]` `When\isApproved.xml`, 17.806 byte:
`pyConditionString` = **`pyWorkPage.PolicyTreatyIn.IsApproved = 1`**, dan bentuk terstrukturnya
`[pyWorkPage.PolicyTreatyIn.IsApproved][=][1]`.

⭐⭐ **Inilah pernyataan OQ-026 yang jauh lebih tajam dari ronde 1:**
⛔ **Bila alur memakai varian `When`, logikanya diketahui. Bila alur memakai varian
`DecisionTable`, logikanya TIDAK ADA — bukan sulit dicari, melainkan tidak terekspor.**
⭐ Jadi OQ-026 berhenti menjadi "mana yang dipakai" dan menjadi **"apakah gerbang pertama seluruh
alur dapat diketahui sama sekali"**.

### 3.3 ⛔ RALAT ronde 1 — "±26 menit" dihitung dari medan yang SALAH

> ⛔⛔ **RALAT.** Kalimat ronde 1 berikut **SALAH**, dan dikutip utuh di sini alih-alih dihapus:
>
> *"`[terverifikasi]` keduanya terbaca di sensus cara B, dan **`pxUpdateDateTime`-nya BERBEDA** —
> `#20170608T022121.683` *(When)* lawan `#20170608T024733.462` *(DecisionTable)*, ⚠️ **selisih
> ±26 menit di hari yang sama**."*
>
> ⛔ **Cara mana yang keliru, dan kenapa:** kedua cap waktu yang saya kutip itu diambil dari
> **`pzOriginalInstanceKey`**, ⛔ **bukan dari `pxUpdateDateTime`** — walau kalimat saya menyebut
> `pxUpdateDateTime`. ⭐ Dan `CLAUDE.md` §4 butir 3 menyatakan tegas:
> **`pzOriginalInstanceKey` adalah asal salinan Save-As, BUKAN identitas.**
> ⭐ Jadi saya memberi label medan yang benar pada angka dari medan yang salah.
>
> ⭐ **Yang benar, dari `pxUpdateDateTime` sungguhan:**
> `When` **`20170608T023134.700`** lawan `DecisionTable` **`20170608T025010.808`**
> ⇒ ⭐ **selisih 18 menit 36 detik**, bukan ±26 menit.
>
> ⭐⭐ **Dan yang lebih penting dari angkanya:** keduanya berbagi **`pyRuleSet` `GISFW` DAN
> `pyRuleSetVersion` `01-01-52` yang identik**. ⛔ Jadi dari empat bagian identitas, ketiganya
> **sama persis** — keduanya berbeda **hanya** pada cap waktu dan pada **tipe aturan**.
> ⚠️ Itu membuat tabrakan namanya **lebih ketat**, bukan lebih longgar.

---

## §4 — Sasaran 4: 75 penggolong, terbaca seluruhnya

### 4.1 ⭐ 75 dari 75 punya syarat yang dapat dibaca — nol kosong

`[terverifikasi]` Perintah audit: sesudah sisa editor dibuang, cacah `pyConditionString`,
lalu bentuk `[properti][operator][nilai]` di `pyLabel`, lalu
`pyCriteriaValue`/`pyPropertyName`/`pyExpression`.

| | Jumlah |
| --- | ---: |
| ⭐ `When` dengan syarat terbaca | ⭐ **75 dari 75** |
| ⛔ `When` benar-benar kosong | ⛔ **0** |

⭐ **Ini membatalkan dugaan yang hampir saya tulis** *(galat instrumen #3, §0.2)*.

### 4.2 ⛔⛔ `pyConditionString` adalah LABEL MANUSIA, dan 35 % di antaranya TIDAK COCOK

⭐⭐ **Temuan yang mengubah cara membaca seluruh 75 penggolong.**

`[terverifikasi]` Contoh yang membuktikannya — `When\IsBondingAndCustomBonds.xml`:

| Medan | Isinya |
| --- | --- |
| `pyConditionString` | ⚠️ **`Kode Bisnis = "02"`** dan **`Kode Bisnis = "58"`** — ⭐ prosa Indonesia, **dua** nilai |
| bentuk terstruktur di `pyLabel` | ⛔ **`[pyWorkPage.Quotation.BusinessType][=]["Bonding"]`** |

⛔ **Keduanya tidak berbicara tentang hal yang sama**: label menyebut **KODE bisnis 02 dan 58**,
aturan yang dieksekusi menguji **TIPE bisnis = "Bonding"**.

`[terverifikasi]` **Diukur di seluruh 75:**

| | Jumlah |
| --- | ---: |
| `When` yang punya **kedua** bentuk | **52** |
| ⛔ yang label-nya **tidak cocok** dengan bentuk terstrukturnya | ⛔ **18** |
| ⭐ **persentase** | ⛔ **35 %** |

⛔ **Akibatnya mengikat:** ⭐ **membaca `pyConditionString` saja akan salah pada lebih dari
sepertiga penggolong.** ⚠️ Dan prosanya **berbahasa Indonesia dan berkata "Kode"** di tempat
aturannya menguji "Tipe" — ⛔ **jenis kekeliruan yang tidak akan tertangkap oleh pembacaan
sekilas.**

⚠️ Bentuk label lain yang terbaca tanpa awalan halaman: `BusinessType = "Car"`,
`BusinessType = "BondingKBG"` — ⛔ **melanggar aturan awalan halaman di dalam korpusnya sendiri**.

### 4.3 Apa yang sebenarnya diuji 75 penggolong itu

`[terverifikasi]` Dari bentuk terstruktur, **halaman pembawa**:

| Halaman | Kemunculan |
| --- | ---: |
| ⭐ `pyWorkPage.Quotation` | ⭐ **67** |
| ⭐ `pyWorkPage.OfferFacIn.QuotationData` | **20** |
| ⚠️ `.OfferFacIn.QuotationData` *(tanpa `pyWorkPage`)* | **9** |
| ⛔ **tanpa halaman sama sekali** | ⛔ **9** |
| `.PolicyTreatyIn` · `.Quotation` · `pyWorkPage` · `OperatorID` | 6 · 4 · 4 · 4 |
| ⛔ `pxRequestor.OperatorID.pyWorkBasketList(1)` | **3** |
| `pyWorkPage.TreatyIn` | 3 |

**Operator:** `=` **120** · `relation` **7** · `!=` **5** · ⚠️ **`EQUALS` 1**
— ⭐ **dua ejaan untuk operator yang sama**.

⚠️⚠️ **Dua halaman sumber untuk penggolongan yang sama** — `Quotation` **67×** lawan
`OfferFacIn.QuotationData` **29×**. ⛔ **Apakah keduanya selalu sama isinya tidak terbaca**,
dan ⛔ **9 penggolong tidak menyebut halaman sama sekali**, sehingga halamannya ditentukan
**konteks pemanggil**.

### 4.4 ⛔⛔ Dua puluh nomor polis PRODUKSI dipaku di dalam aturan hidup

⛔⛔ **Temuan paling mengkhawatirkan untuk `[Finance]`**, dan ia analog langsung dengan pola
kasus-uji yang sudah dikenal.

`[terverifikasi]` **Dihitung dua cara:** *(a)* kemunculan literal ⇒ **47**; *(b)* nomor polis
**unik** ⇒ **20**. Jendela: seluruh 278 berkas, sisa editor dibuang, teks ter-unescape dua kali.
Perintah: pola `"RNM-[A-Z0-9.\-]{8,}"`.

| Berkas | Kemunculan | Perannya |
| --- | ---: | --- |
| ⛔⛔ `Activity\SumTSIPremiSpreadedRNM_Act.xml` | ⛔ **24** | ⭐ **aturan penjumlah TSI dan premi** — inti rantai uang |
| ⛔⛔ `Activity\ProtectFIREMBUPA_Act.xml` | ⛔ **12** | aturan proteksi lini kebakaran |
| ⛔ `When\IsErrorSpreading.xml` | **11** | ⭐ **penggolong bernama "ada galat penyebaran"** |

⛔ **Nomor polisnya bertanggal 2023, 2024, dan 2025** — ⭐ **ditambahkan bertahap selama
bertahun-tahun**, bukan satu kali.
⭐ `[dugaan]` bentuknya **daftar pengecualian yang tumbuh**: setiap kali sebuah polis berperilaku
salah, nomornya **ditempelkan ke dalam aturan** alih-alih sebabnya diperbaiki. ⛔ **Belum
terverifikasi** — hanya pemiliknya dapat memastikan.

⚠️ **Kenapa ini menyentuh uang, bukan kerapian:** ⛔ dua dari tiga berkas itu **menghitung TSI,
premi, dan proteksi**. ⭐ Bila sistem baru dibangun tanpa daftar ini, **dua puluh polis akan
dihitung dengan cara berbeda dari sekarang** — dan ⛔ **tidak ada yang akan melihatnya**, sebab
tidak ada pesan galat.

### 4.5 Penggolong yang menguji jabatan dan lingkungan

`[terverifikasi]` `OperatorID.pyPosition` · `OperatorID.pyAccessGroup` ·
`pxRequestor.OperatorID.pyWorkBasketList(1)` · `pyWorkPage.pyWorkIDPrefix = "CLM-"`
*(`IsClaim`)* · `pyWorkPage.LetterNo` · `pyWorkPage.PolicyTreatyIn.PolicyNo` ·
`QuotationData.StatusBusiness = 3` *(`IsEDM`)* · `QuotationData.EdmTypeNew = 4`
*(`IsEDMRiSlip`)* · `BusinessCode = "10085"` *(`IsBuilderRisk`)*.

⚠️ **Kode angka tanpa arti terverifikasi:** `StatusBusiness = 3` · `EdmTypeNew = 4` ·
`BusinessCode = "10085"` · `Kode Bisnis "02"` / `"58"`. ⛔ **Tidak saya tebak.**

---

## Bab A — Tambahan stored procedure dan SQL mentah

⭐ **NIHIL tambahan, dan dibuktikan.**

`[terverifikasi]` Jendela: **237 berkas non-`RDBList`** *(278 − 41)*, sisa editor dibuang,
teks ter-unescape dua kali. Pola: `(POOLDATA|DATAPEGA|ARASAPAS|ASM)\.[A-Z0-9_]+\s*\(`.

| Hasil | Jumlah |
| --- | ---: |
| ⭐ pemanggilan procedure Oracle **baru** | ⭐ **0** |
| yang ditemukan | **1** — `@ASM.GetPageJSONString` **2 kemunculan** |

⭐ `@ASM.GetPageJSONString` **sudah tercatat** — OQ-012, dan ronde 1 §A.3. ⛔ **Bukan temuan baru.**

⚠️ **14 berkas non-`RDBList` memuat kata kunci SQL**, tetapi diperiksa: seluruhnya **teks
dokumentasi dan tanda tangan metode**, bukan SQL yang dijalankan — 8 `Activity` dan 4 `Harness`,
**masing-masing 1 kemunculan saja**. ⭐ Satu kemunculan per berkas adalah tanda prosa, bukan
pernyataan.

⭐ **Jadi ketiga procedure ronde 1 tetap ketiganya**, dan ⛔ **nol badannya masih ada.**

## Bab B — Tambahan objek Oracle

⭐ **NIHIL tambahan, dan cara pembuktiannya sekaligus galat instrumen keempat.**

⛔ Sisiran pertama melaporkan **19 objek baru**. ⚠️ Diperiksa: teratasnya **`THE` 74×**,
**`FULLY` 68×**, `PROPERTY` 11×, `MATH` 6×, `DATA`, `ENCRYPTED`, `ACTION`, `ANOTHER` —
⛔ **seluruhnya kata prosa Inggris** dari teks dokumentasi yang kebetulan mengikuti kata `FROM`
atau `UPDATE` dalam kalimat biasa.

⭐ **Sesudah disaring: 0 objek Oracle baru.** ⛔ Ke-31 objek ronde 1 tetap 31.

---

## Bab C — Register `[terbuka]` berpemilik

⭐ Nomor OQ dari `discovery/open-questions.md`. ⭐ **72 OQ terdaftar; nomor bebas berikutnya
OQ-073.** ⭐ Butir baru ronde 1 dirujuk sebagai **"ronde 1 baru #1"** sampai **#7**.

⛔⛔ **NOL butir saya tutup.** Penutupan hanya oleh pemilik peran.

### C.1 Butir yang BERUBAH STATUSNYA karena ronde 2

⛔ **Berubah artinya diperdalam atau dipersempit — ⛔ BUKAN ditutup.**

| Butir | Perubahannya | Pemilik |
| --- | --- | --- |
| ⛔⛔ **OQ-026** | ⭐ **Dipertajam menentukan.** Bukan lagi "mana dari dua yang dipakai": varian `When` **terbaca penuh**, varian `DecisionTable` ⛔ **nol baris keputusannya terekspor**. ⭐ Ditambah: ketiga bagian identitas lain **identik** *(ruleset `GISFW`, versi `01-01-52`)*; selisih waktu sebenarnya **18 mnt 36 dtk** — lihat RALAT §3.3 | `[Product+Underwriting]` |
| ⭐ **OQ-024** | ⭐ **Dipersempit.** OQ-024 berbunyi pemetaan antrean *"tidak terbaca"*. ⭐ `[terverifikasi]` **Satu pemetaan TERBACA** — `ListSuggest` menguji `PositionNote != 'ReasTreatyInSecHead'` sebagai teks. ⛔ Empat antrean lain tetap tak terbaca | `[IAM]` · `[Product+Underwriting]` |
| ⭐ **OQ-007** *(dirujuk OQ-024 ronde 1)* | ⭐ **Diperkuat dari sudut baru.** `pyAssociatedPrivileges` **295 kemunculan di 29 berkas layar ⇒ 0 berisi** — ⛔ tidak satu pun elemen layar membawa hak-akses Pega | `[IAM]` |
| ⛔ **OQ-021** | ⭐ **Diperluas.** Pola guard **KEEMPAT** ditemukan: `OperatorID.pxInsName` dicocokkan ke nama orang di sebuah layar — berbeda dari medan yang disapu OQ-021 dan dari `pyTelephone` OQ-027. **12 kemunculan, 4 berkas, 4 nama** | `[IAM]` |
| ⭐ **ronde 1 baru #1** *(slot generik)* | ⭐ **Dijawab SEPARUH dari korpus.** Catatan langkah menuliskan arti sebagian slot — 44=deduction, 31=currency, 47=balance due to, 13=net premium, 32=gross premi. ⛔ Sebagian saja, dan ⛔ **catatan bukan aturan**. ⭐ **Diperluas juga**: slot generik kini terbukti menggerbangi **LAYAR**, bukan hanya SQL | `[pengembang Pega lama]` |
| ⚠️ **ronde 1 baru #6** *(31 berkas nol dibuka)* | ⭐ **Dikerjakan.** 31 dari 31 kini disisir; **11 murni tata letak**; aturan bisnis memusat di 2 berkas. ⛔ **Tidak satu pun dibaca utuh** — tetap terbuka sebagai kedalaman | `[pengembang Pega lama]` |
| ⛔ **OQ-002 · OQ-013** | ⭐ **Tidak berubah**, ⛔ dan §1 membuatnya **lebih menekan**: bahkan bila badan procedure diserahkan, nilai yang **masuk** ke procedure ditetapkan langkah `Property-Set` yang **tidak terekspor** | `[DBA]` |
| ⛔ **OQ-025** | ⭐ **Tidak berubah, tetap BLOCKER.** ⛔ `EDM Treaty In` tidak dibuka; varian tidak dipinjam | `[Product+Underwriting]` · `[pemilik export Pega]` |

⚠️ **Sepuluh butir ronde 1 lainnya** — OQ-012, 059, 023, 028, 020, 022, 011, 018, 009, 054, dan
baru #2 · #3 · #4 · #5 · #7 — ⭐ **tidak berubah statusnya di ronde ini.** ⛔ Tidak diulang di sini;
rujuk Bab C ronde 1.

### C.2 ⭐ Butir yang LAHIR di ronde 2 — belum bernomor register

⛔ **Saya tidak memberi nomor OQ sendiri.** Dirujuk **"ronde 2 baru #1"** dan seterusnya.

| # | Butir | Pemilik | Kenapa memblokir | Akibat bila dijawab salah |
| --- | --- | --- | --- | --- |
| ⛔⛔ **1** | **942 langkah `Property-Set` di 92 Activity, NOL nilainya terekspor** — hanya 74 salinan editor, dan salinan editor bukan aturan hidup | `[pemilik export Pega]` · `[pengembang Pega lama]` | ⛔⛔ **rantai hitung uang tidak dapat direkonstruksi** — setara blocker OQ-025 | ⛔⛔ seluruh perhitungan uang disalin dari tebakan, dan selisihnya baru terlihat di laporan keuangan |
| ⛔⛔ **2** | **20 nomor polis produksi dipaku di 3 aturan hidup**, 2 di antaranya menghitung TSI/premi/proteksi; bertanggal 2023–2025 | ⭐ `[Finance]` · `[Product+Underwriting]` | ⛔ dua puluh polis **dihitung berbeda** dari aturan umum | ⛔ 20 polis berpindah ke perhitungan lain **tanpa pesan galat** |
| ⛔⛔ **3** | **Baris keputusan kedua `DecisionTable` tidak terekspor**, dan salah satunya **diperbarui 22 Sep 2025 dengan 8 versi sebelumnya** | `[pemilik export Pega]` | ⛔ penggolong **hidup dan sering diubah** yang isinya tak terlihat | ⛔ penggolongan bisnis salah, memengaruhi treaty mana yang dipakai |
| ⛔⛔ **4** | **35 % penggolong (18 dari 52) label manusianya TIDAK COCOK dengan aturan yang dieksekusi** — contoh: label "Kode Bisnis 02/58", aturan menguji "BusinessType = Bonding" | `[pengembang Pega lama]` · `[Product+Underwriting]` | ⛔ **membaca label saja salah pada sepertiga penggolong** | ⛔ spesifikasi ditulis dari label, sistem baru menggolongkan berbeda dari Pega |
| ⛔ **5** | **Gerbang pertama alur dibandingkan sebagai TEKS dan sebagai ANGKA** — `.IsApproved = '0'` berkutip 2× lawan `.IsApproved==0` telanjang 2× | `[pengembang Pega lama]` | ⛔ "kosong" tak dapat dibedakan dari "nol" pada bendera paling menentukan | ⛔ pengajuan lolos atau tertahan salah |
| ⛔ **6** | **91 elemen layar mati sejak dibangun** *(`1=2` 72×, `NEVER` 15×, `1==2` 2×)*, dan **38 elemen hanya-baca permanen** — ⭐ **seluruh 38 di dua layar Kepala Departemen** | `[Product+Underwriting]` | ⚠️ apakah medan itu memang tidak boleh dipakai, atau dimatikan sementara lalu terlupakan | ⛔ medan mati ikut dibangun, atau medan yang masih perlu ikut dibuang |
| ⛔ **7** | **Antrean ditunjuk menurut POSISI dalam daftar** — `pyWorkBasketList(1)` 3× di penggolong, `(2)` 1× di layar | `[IAM]` | ⛔ mengurutkan ulang antrean pengguna **mengubah perilaku** tanpa menyentuh aturan | ⛔ wewenang berpindah diam-diam saat keanggotaan antrean berubah |
| ⚠️ **8** | **`pyDisabledNew` bernilai `truewhn` 4×** — ⭐ `[dugaan]` salah ketik `truewhen` | `[pengembang Pega lama]` | ⚠️ bila Pega mengabaikannya, penonaktifan **diam-diam tidak berlaku** | ⚠️ medan yang seharusnya terkunci dapat disunting |
| ⚠️ **9** | **Dua halaman sumber untuk penggolongan yang sama** — `Quotation` 67× lawan `OfferFacIn.QuotationData` 29×, plus **9 penggolong tanpa halaman** | `[pengembang Pega lama]` | ⛔ halaman menentukan nilai; tanpa halaman, konteks pemanggil yang menentukan | ⛔ penggolong membaca nilai dari tempat yang salah |
| ⚠️ **10** | **Dua invarian uang hanya hidup sebagai CATATAN langkah** — kewajiban premi nol, dan "total premi share cedant harus = premi RNM" | ⭐ `[Finance]` | ⛔ aturan uang yang **tidak dapat dipastikan ditegakkan**, sebab langkah penegaknya tak terekspor | ⛔ ketidakseimbangan lolos tanpa terdeteksi |
| ⚠️ **11** | **`.Claim` dibandingkan terhadap teks kosong DAN angka nol** dalam satu syarat wajib-isi | ⚠️ `[Finance]` · `[pengembang Pega lama]` | ⚠️ medan uang bertipe ganda | ⚠️ klaim bernilai nol diperlakukan sebagai belum diisi, atau sebaliknya |
| ⚠️ **12** | **JavaScript peramban mentah di 44 tombol** — `window.close`, `opener.location.reload`; arsitektur jendela sembul | `[pengembang Pega lama]` | ⚠️ perilaku bergantung jendela induk yang mungkin tidak ada | ⚠️ data tampak tidak tersimpan padahal tersimpan, atau sebaliknya |
| ⚠️ **13** | **Tautologi dan kontradiksi dieja 3 dan 2 cara** di korpus *(`1==1`/`1=1`/`1 = 1`; `1=2`/`1==2`)*; `pyRequiredNew` dieja `always` dan `true`; operator dieja `=` dan `EQUALS` | `[pengembang Pega lama]` | ⚠️ **setiap sensus atas korpus ini akan meleset** bila satu ejaan saja dicari | ⚠️ angka migrasi salah, dan salahnya tidak terlihat |

### C.3 Rekap per pemilik

⭐ **13 butir lahir di ronde 2.** ⛔ **0 ditutup.** ✅ **Dihitung dua cara** *(baris tabel C.2
yang memuat pemilik, dan kemunculan literal nama pemilik)*; sesudah ralat di bawah, **sepakat**.

> ⛔⛔ **RALAT — dua galat, ditangkap sensus berkas ini sebelum dilaporkan. Dan yang kedua
> adalah PENGULANGAN kesalahan saya sendiri di ronde 1.**
>
> **Galat 1.** Tabel di bawah semula berbunyi *"`[pengembang Pega lama]` **7**"*.
> ⭐ **Yang benar: 8** — butir #1, #4, #5, #8, #9, #11, #12, #13. ⚠️ Saya mendaftar delapan
> nomor lalu menuliskan angka tujuh.
>
> **Galat 2 — dan ini yang penting.** Cacahan pertama melaporkan
> `[Product+Underwriting]` = **0**, padahal sebenarnya **3**. ⛔ **Sebabnya saya mengeja pemilik
> yang sama dua cara di dalam berkas ini** — `[Product+UW]` di **5 tempat** dan
> `[Product+Underwriting]` di **2 tempat**.
>
> ⛔⛔ **Ronde 1 sudah mencatat galat yang PERSIS sama, dengan pemilik yang PERSIS sama, dan
> saya mengulanginya.** ⭐ Itu memperkuat butir **ronde 2 baru #13**: ejaan ganda bukan hanya
> penyakit korpus — **ia penyakit berkas laporan ini juga**, dan satu-satunya yang menangkapnya
> adalah sensus yang dijalankan dengan perintah, bukan dengan mata.
>
> ⭐ **Ejaannya kini diseragamkan menjadi `[Product+Underwriting]` di seluruh berkas.**

| Pemilik | Butir ronde 2 | Nomornya |
| --- | ---: | --- |
| ⛔ `[pengembang Pega lama]` | ⭐ **8** | #1 *(bersama)* · #4 *(bersama)* · #5 · #8 · #9 · #11 *(bersama)* · #12 · #13 |
| ⭐ `[Finance]` | ⭐ **3** | #2 *(bersama)* · #10 · #11 *(bersama)* |
| ⛔ `[Product+Underwriting]` | **3** | #2 *(bersama)* · #4 *(bersama)* · #6 |
| ⛔ `[pemilik export Pega]` | **2** | #1 *(bersama)* · #3 |
| ⚠️ `[IAM]` | **1** | #7 |
| `[DBA]` | **0** | ⭐ tidak ada butir basis data **baru** — ⚠️ tetapi OQ-002 dan OQ-013 **menjadi lebih menekan**, §1 |
| `[work owner]` | **0** | seluruh butir jatuh ke pemilik lebih spesifik |

### C.4 ⭐⭐ `[Finance]` — lubang ronde 1 kini terisi, dan isinya buruk

⛔ **Brief menuntut ini disebut eksplisit, jadi disebut eksplisit.**

⭐ **Ronde 1: `[Finance]` = 0 butir**, dan ronde 1 menyebut sendiri bahwa nol itu **lubang, bukan
kebersihan**.
⭐ **Ronde 2: `[Finance]` = 3 butir.** ⛔ **Lubangnya terisi, dan yang keluar dari dalamnya
lebih buruk dari yang saya kira:**

1. ⛔⛔ **20 nomor polis produksi dipaku di dalam aturan penghitung uang** — bukan di satu tempat,
   melainkan **24 kemunculan di aturan penjumlah TSI dan premi**, ditambah 12 di aturan proteksi.
2. ⛔ **Dua invarian uang hanya berupa catatan** — termasuk sebuah **kewajiban keseimbangan**.
3. ⚠️ **Sebuah medan uang bertipe ganda** — dibandingkan sebagai teks dan angka.

⛔⛔ **Dan yang paling menentukan bagi `[Finance]`, walau pemiliknya bukan Finance:** butir #1 —
⛔ **rumus uangnya tidak ada di ekspor.** ⭐ Jadi `[Finance]` **belum dapat memeriksa apa pun
tentang perhitungan**, sebab perhitungannya belum dapat ditunjukkan kepada mereka.

---

## Bab D — PERTANYAAN SIAP KIRIM

⛔ **Hanya yang BARU.** Penomoran lanjut dari **P18**. Penjawabnya **tidak membaca XML**.
⛔ **Nol saya jawab sendiri.** Urut dari yang paling memblokir.

---

### Untuk pemilik export Pega

#### **P18 — Isi 942 langkah penghitungan tidak ada di dalam berkas yang kami terima. Bisa dikirim ulang?**

Sistem lama punya **942 langkah** yang menetapkan nilai — hasil perhitungan premi, bagian,
pengurangan, dan pajak. Kami dapat melihat **bahwa** langkah itu ada, dan **urutannya**,
⛔ **tetapi tidak satu pun rumus atau nilai yang ditetapkannya ikut terkirim.** Yang ada hanya
salinan sementara yang tertinggal dari layar penyunting, dan salinan itu **terpotong di tengah**
sehingga tidak dapat dipakai.

**Konteks:** tanpa isi langkah-langkah ini, perhitungan uang sistem lama tidak dapat ditiru —
hanya ditebak.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh berkas** — ekspor ulang dengan **isi langkah
disertakan**. Bila ekspor jenis ini memang tidak pernah menyertakannya, katakan begitu, dan
sebutkan **cara lain yang tersedia** *(cetakan layar aturan, dokumentasi, atau akses baca ke
sistem lama)*.

**Dampak bila salah:** ⛔⛔ seluruh perhitungan uang dibangun dari tebakan, dan **selisihnya baru
muncul di laporan keuangan berbulan-bulan kemudian**.

rujukan: 942 langkah `Property-Set` di 92 `Activity`; `PropertiesValue` = 74 dengan salinan editor,
**0** tanpa; `pyStepsCallParams` hanya placeholder — ronde 2 baru #1

---

#### **P19 — Isi dua tabel keputusan juga tidak terkirim, dan satu di antaranya masih sering diubah**

Ada dua **tabel keputusan** — daftar "kalau begini maka begitu" — dan ⛔ **baris-barisnya tidak
ikut terkirim.** Yang tersisa hanya namanya. Salah satunya **terakhir diubah 22 September 2025**
dan punya **delapan versi sebelumnya**, jadi ia jelas masih dipakai dan masih disesuaikan.

**Konteks:** salah satu dari keduanya adalah **pemeriksaan pertama** yang dilewati setiap
pengajuan.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh berkas** — kedua tabel keputusan itu beserta
barisnya. Bila tidak bisa, **cetakan layarnya** sudah cukup.

**Dampak bila salah:** ⛔ penggolongan bisnis menjadi berbeda, dan itu menentukan **treaty mana
yang dipakai** untuk sebuah penutupan.

rujukan: `DecisionTable\isApproved.xml` · `DecisionTable\BusinessType_DeT.xml`, nol
`pyCriteriaValue`/`pyResult` — OQ-026, ronde 2 baru #3

---

### Untuk Finance dan Product + Underwriting

#### **P20 — Dua puluh nomor polis tertulis langsung di dalam aturan penghitung uang. Apakah itu disengaja?**

⛔ Di dalam tiga aturan sistem lama terdapat **dua puluh nomor polis tertentu yang ditulis
langsung**, bukan dibaca dari data. ⚠️ **Dua dari tiga aturan itu menghitung nilai pertanggungan
dan premi.** Nomor-nomornya bertanggal **2023, 2024, dan 2025** — jadi ditambahkan bertahap
selama tiga tahun.

**Konteks:** artinya dua puluh penutupan itu **dihitung dengan cara berbeda** dari semua penutupan
lain, dan perbedaannya tidak tercatat di mana pun selain di dalam aturan itu.

**Bentuk jawaban yang diharapkan:** pilihan ganda + keterangan — **(a)** memang perlakuan khusus
yang disengaja, jelaskan alasannya · **(b)** tambalan sementara untuk memperbaiki data yang salah,
sudah tidak perlu · **(c)** tidak tahu, perlu diperiksa satu per satu.
⭐ Bila (a) atau (c): ⭐ **kami perlu daftar itu dibaca bersama**, sebab sistem baru harus tahu
apakah dua puluh polis ini tetap istimewa.

**Dampak bila salah:** ⛔⛔ dua puluh penutupan berpindah ke perhitungan yang berbeda dari
sekarang, **tanpa satu pun pesan galat** yang memberi tahu.

rujukan: `Activity\SumTSIPremiSpreadedRNM_Act.xml` 24× · `Activity\ProtectFIREMBUPA_Act.xml` 12× ·
`When\IsErrorSpreading.xml` 11×; 47 kemunculan, 20 nomor unik — ronde 2 baru #2

---

#### **P21 — Dua aturan uang hanya tertulis sebagai catatan. Apakah keduanya benar-benar berlaku?**

Di dalam catatan pengembang pada langkah perhitungan, ada **dua aturan uang** yang dinyatakan
sebagai kewajiban: bahwa **premi polis induk harus nol** dalam keadaan tertentu, dan bahwa
**total bagian premi yang diberikan kepada pemberi bisnis harus sama dengan premi perusahaan**.
⛔ Keduanya **hanya ada sebagai catatan** — kami tidak dapat memastikan apakah sistem benar-benar
menolak ketika keduanya dilanggar.

**Konteks:** kalau ini benar-benar kewajiban, sistem baru harus menegakkannya; kalau bukan,
membangunnya akan menolak data yang sah.

**Bentuk jawaban yang diharapkan:** untuk masing-masing, pilihan ganda — **(a)** kewajiban keras,
harus ditolak bila dilanggar · **(b)** peringatan saja, boleh dilanjutkan · **(c)** bukan aturan,
hanya catatan lama. ⭐ Bila (a): **siapa yang boleh menyimpang**, kalau ada.

**Dampak bila salah:** ⛔ ketidakseimbangan premi **lolos tanpa terdeteksi**, atau sebaliknya
penutupan yang sah **ditolak**.

rujukan: `pyStepsDescription` — `set protect kalau master polis premi harus 0` ·
`TOTAL PREMI SHARE CEDANT HARUS = PREMI RNM` — ronde 2 baru #10

---

#### **P22 — Sembilan puluh satu bagian layar dibuat agar tidak pernah muncul, dan tiga puluh delapan dikunci permanen. Masih perlu?**

⛔ **Sembilan puluh satu** bagian layar diberi syarat tampil yang **tidak mungkin pernah benar** —
setara dengan menuliskan "tampilkan jika 1 sama dengan 2". Dan ⛔ **tiga puluh delapan** medan
diberi syarat **hanya-baca yang selalu benar**, sehingga tidak pernah dapat disunting.
⭐ **Seluruh tiga puluh delapan itu ada di layar Kepala Departemen.**

**Konteks:** kami perlu tahu mana yang memang tidak boleh dipakai lagi, dan mana yang dimatikan
sementara lalu terlupakan — sebab yang pertama tidak perlu dibangun, yang kedua perlu.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** memang sudah tidak dipakai, boleh
hilang · **(b)** dimatikan sementara, seharusnya kembali · **(c)** campuran, perlu ditinjau
bersama. ⭐ Untuk layar Kepala Departemen: **apakah jabatan itu memang hanya boleh melihat, tidak
mengubah?**

**Dampak bila salah:** ⛔ medan mati ikut dibangun dan dirawat selamanya, atau medan yang masih
dibutuhkan **hilang** dari sistem baru.

rujukan: `1=2` 72× · `NEVER` 15× · `1==2` 2×; `pyReadOnlyCondition` `1==1` 32× + varian 6×,
seluruhnya di `DetailDeptHeadTreatyIn_UW` dan `GeneralDeptHeadTreatyIn_UW` — ronde 2 baru #6

---

### Untuk pengembang Pega lama

#### **P23 — Pada sepertiga penggolong, keterangan yang tertulis berbeda dari yang dikerjakan aturan. Mana yang benar?**

Ada 75 aturan kecil yang menggolongkan jenis bisnis. Masing-masing punya **dua hal**: keterangan
yang ditulis manusia, dan syarat yang benar-benar dijalankan. ⛔ **Pada 18 dari 52 yang punya
keduanya — sepertiga — keduanya tidak berbicara tentang hal yang sama.** Satu contoh: keterangannya
berbunyi *"Kode Bisnis 02 dan 58"*, sedangkan yang dijalankan memeriksa **jenis** bisnis bernama
*"Bonding"*.

**Konteks:** kami menulis spesifikasi dari aturan ini. Bila kami membaca keterangannya, kami salah
pada sepertiga kasus; bila kami membaca yang dijalankan, mungkin keterangannya yang mencerminkan
maksud sebenarnya.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** yang dijalankan selalu benar,
keterangannya basi · **(b)** keterangannya menyatakan maksud, yang dijalankan mungkin bug ·
**(c)** kasus per kasus, ⭐ **perlu ditinjau bersama 18 aturan itu**.

**Dampak bila salah:** ⛔ sistem baru menggolongkan bisnis **berbeda dari sistem lama**, dan
penggolongan itu menentukan treaty, premi, dan laporan.

rujukan: 52 `When` punya `pyConditionString` + bentuk bertanda kurung; **18 tidak cocok (35 %)`;
contoh `When\IsBondingAndCustomBonds.xml` — ronde 2 baru #4

---

#### **P24 — Bendera "sudah disetujui" dibandingkan sebagai teks di satu tempat dan sebagai angka di tempat lain. Mana yang dimaksud?**

Pemeriksaan **pertama** yang dilewati setiap pengajuan membaca sebuah penanda "sudah disetujui".
⛔ Di satu tempat penanda itu dibandingkan **sebagai teks** — nol di dalam tanda kutip — dan di
tempat lain **sebagai angka** — nol tanpa kutip.

**Konteks:** bila penanda itu belum pernah diisi, "kosong" dan "nol" adalah dua hal berbeda dalam
satu perbandingan dan satu hal yang sama dalam perbandingan lainnya.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** penanda ini angka, pembandingan teks
adalah kekeliruan · **(b)** penanda ini teks · **(c)** boleh keduanya, tidak pernah jadi masalah.
⭐ Dan: **apa arti penanda ini saat belum pernah diisi** — belum diperiksa, atau ditolak?

**Dampak bila salah:** ⛔ pengajuan yang belum diperiksa **dianggap ditolak**, atau sebaliknya
**diteruskan** seolah sudah disetujui.

rujukan: `.IsApproved = '0'` berkutip 2× vs `.IsApproved==0` telanjang 2× di layar;
`When\isApproved.xml` menguji `= 1` — ronde 2 baru #5

---

#### **P25 — Wewenang di beberapa tempat ditentukan oleh POSISI seseorang dalam sebuah daftar. Apakah urutan itu dijamin?**

Beberapa aturan memutuskan apa yang boleh dilihat seseorang dengan melihat **antrean kerja
nomor satu** atau **nomor dua** dalam daftar antrean pengguna itu — ⛔ **bukan dengan menyebut
nama antreannya.**

**Konteks:** kalau seseorang ditambahkan ke antrean baru, atau daftarnya diurutkan ulang,
⛔ **perilaku aturan berubah tanpa ada yang menyentuh aturan itu.**

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** urutannya memang dijamin, jelaskan
bagaimana · **(b)** tidak dijamin, ini memang rapuh · **(c)** tidak tahu. ⭐ Bila (a) atau (c):
⭐ **antrean mana yang seharusnya dimaksud** pada posisi satu dan posisi dua.

**Dampak bila salah:** ⛔ wewenang seseorang **berpindah diam-diam** ketika keanggotaan antreannya
berubah karena alasan yang sama sekali lain.

rujukan: `pxRequestor.OperatorID.pyWorkBasketList(1)` 3× di `When`;
`OperatorID.pyWorkBasketList(2).pyWorkBasketName` 1× di layar — ronde 2 baru #7

---

#### **P26 — Satu nilai pengaturan tampak salah ketik. Apakah ia bekerja?**

Di empat tempat, sebuah pengaturan yang mengunci medan diberi nilai yang **tampaknya salah
ketik** — satu huruf hilang dari kata yang seharusnya. Di tempat lain nilai yang sama ditulis
dengan benar.

**Konteks:** kalau sistem lama **mengabaikan** nilai yang salah ketik itu, maka keempat medan itu
sebenarnya **tidak pernah terkunci** — padahal seseorang berniat menguncinya.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** salah ketik, penguncian tidak bekerja ·
**(b)** nilai itu sah, bekerja normal · **(c)** tidak tahu, perlu dicoba di sistem lama.

**Dampak bila salah:** ⚠️ medan yang seharusnya terkunci **dapat disunting**, atau sebaliknya kami
mengunci medan yang selama ini bebas.

rujukan: `pyDisabledNew` = `truewhn` 4×, di samping `true` 17× dan `always` 12× — ronde 2 baru #8

---

#### **P27 — Wadah slot bernomor: sebagian artinya tertulis di catatan. Apakah catatan itu masih benar, dan di mana sisanya?**

Melanjutkan pertanyaan sebelumnya tentang wadah keterangan bernomor: ⭐ **kami menemukan sebagian
artinya tertulis di dalam catatan pengembang** — slot 44 untuk pengurangan, 31 untuk mata uang,
47 untuk saldo terutang, 13 untuk premi neto, 32 untuk premi bruto.

**Konteks:** catatan bisa basi tanpa ada yang tahu, dan kami hanya menemukan arti **sebagian**
slot — bukan semuanya.

**Bentuk jawaban yang diharapkan:** ⭐ **konfirmasi + kelengkapan** — apakah kelima arti di atas
masih benar, dan **di mana arti slot lainnya dapat dibaca**. Bila tidak ada tempatnya,
katakan begitu.

**Dampak bila salah:** ⛔ nilai uang masuk ke kolom yang salah, dan **kesalahannya tidak
menimbulkan pesan galat** — hanya angka yang salah tempat.

rujukan: `pyStepsDescription` — `Deduction (CARI44[deduction])` ·
`Net Premi (CARI31[CURRENCY], CARI39, CARI47[BALANCEDUETO], CARI13[NET PREMIUM])` ·
`Gross Premi (CARI32, CARI31[CURRENCY])` — OQ-059, ronde 1 baru #1, ronde 2 pendalaman

---

### Untuk IAM

#### **P28 — Nama orang juga dipakai di layar, lewat medan keempat yang belum pernah kami laporkan**

Selain tempat-tempat yang sudah kami tanyakan, ⛔ **nama orang tertentu juga menentukan apa yang
muncul di layar** — di **empat layar**, **dua belas tempat**, memuat **empat nama berbeda**.
⛔ **Nilainya tidak kami salin ke berkas mana pun.** Dan salah satu layar memakai **medan identitas
yang berbeda lagi** dari yang pernah kami laporkan.

**Konteks:** pada dua layar, nama yang sama dipakai **dua arah** — ada bagian yang muncul hanya
untuk orang itu, dan bagian lain yang muncul untuk **semua kecuali** orang itu. ⭐ Mengeluarkan
orang itu mengubah **dua** perilaku sekaligus.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** ganti dengan peran, sebutkan perannya ·
**(b)** memang harus orang tertentu, jelaskan alasannya · **(c)** sudah tidak relevan.
⭐ Dan: **apa yang seharusnya terjadi** pada bagian layar itu untuk orang lain.

**Dampak bila salah:** ⛔ bagian layar **hilang** bagi orang yang berhak, atau **terlihat** oleh
yang tidak berhak.

rujukan: `pyUserIdentifier` di `DetailDeptHeadTreatyIn_UW` 3× · `GeneralDeptHeadTreatyIn_UW` 3× ·
`ListSuggest` 4×; ⚠️ `pxInsName` di `DetailPoliciesNonProportional` 2× — OQ-021 diperluas,
ronde 2 baru *(pola guard keempat)*

---

## Bab E — Lampiran bukti isolasi

| Yang dibuktikan | Keadaan | Perintah audit |
| --- | --- | --- |
| korpus `NB Treaty In` tidak berubah | ✅ **`62a3735ebb2eabb788c8cd8abdda78ca`** — cocok ronde 1 | `find . -name '*.xml' \| sort \| xargs cat \| md5sum` |
| ⛔ **`EDM Treaty In` NOL dibuka** | ✅ **nol** — ⭐ OQ-025 tetap dicatat sebagai pertanyaan, ⛔ varian **tidak dipinjam** | folder itu **tidak pernah** menjadi argumen perintah apa pun; seluruh skrip ber-`ROOT = D:\XML\RNM_BRD\NB Treaty In` |
| ⛔ folder korpus modul lain | ✅ **NOL dibuka** | sama — satu `ROOT` di setiap skrip |
| ⛔ `D:\XML\nusantara-re\` | ✅ **NOL disentuh** | tidak pernah muncul sebagai argumen |
| ⛔ `grilling-ronde-1.md` **tidak disunting** | ✅ **nol suntingan** — ⭐ tiga ralat ditulis **di berkas ini** | `ls -la .scratch\nb-treaty-in\` — mtime ronde 1 tetap **2026-09-22 09:25** |
| berkas dibuat | ✅ **tepat SATU** — `grilling-ronde-2.md` | `ls .scratch\nb-treaty-in\` |
| ⛔ spec / tiket | ✅ **NOL** | tidak ada berkas ber-nama `spec`/`issue`/`tiket` |
| kode · DDL · `CREATE TABLE` · daftar kolom usulan | ✅ **NOL** | ⚠️ ungkapan terlarangnya **disebut** di kalimat yang menyatakannya nol |
| ⛔ nilai nama orang | ✅ **NOL disalin** | §2.3 dan P28 menyebut **medan, berkas, dan jumlah**; ⛔ **nol nama** |
| butir terbuka yang saya tutup | ✅ ⛔ **NOL** | Bab C — setiap butir tetap berpemilik |
| nomor OQ dikarang sendiri | ✅ **NOL** | butir baru dirujuk **"ronde 2 baru #n"**; ⭐ nomor bebas OQ-073 disebut, **tidak dipakai** |

### ⚠️ Baris yang TIDAK bersih — dilaporkan sesuai perintah brief

⛔ Brief memerintahkan: *bila sambungan sesi menyuntik `.scratch` konteks lain tanpa diminta,
catat kejadiannya dan nyatakan apa yang tidak merembes.* ⭐ **Itu terjadi, dan ini catatannya.**

⚠️ Sesi ini **panjang dan pernah disambung**. Pada penyambungan, harness **menyuntikkan ulang dua
bacaan lama** dari `.scratch` konteks lain — dua berkas `urutan-tiket.md` milik konteks klaim.
⛔ **Bukan saya yang memanggilnya**, dan ⛔ **tidak saya panggil sekali pun di ronde 2.**

⭐ **Apa yang tidak merembes, dan buktinya:**
⛔ **Nol nama tabel, nol nama kolom, nol pola rancangan** dari konteks lain dipakai di berkas ini.
⭐ Seluruh angka di sini berasal dari **sisiran korpus `NB Treaty In` sendiri** dan dari **berkas
yang brief perintahkan dibaca** — `grilling-ronde-1.md`, `discovery/`, `CONTEXT.md`, `docs/adr/`,
`CLAUDE.md`.
⚠️ **Satu penyebutan lintas-konteks yang memang ada**, dan ia dari register bersama bukan dari
korpus: Bab A ronde 1 menyebut tiga procedure jalur Life yang badannya sudah diserahkan.
⭐ Di ronde 2 itu **tidak diulang**, dan penyebutan aslinya justru untuk membuktikan bahwa
ketiganya **BUKAN** procedure modul ini.

---

## Penutup

### Cakupan sesudah ronde 2

| Golongan | Ronde 1 | ⭐ Ronde 2 | Keadaan sekarang |
| --- | --- | --- | --- |
| `RDBList` **41** | ✅ pengurai penuh | — | ✅ **selesai** |
| blok PL/SQL **5** | ✅ dibaca utuh, SQL dikutip | — | ✅ **selesai** |
| `DecisionTable` **2** | ⛔ nol dibuka | ⭐ **dibaca UTUH** | ⛔ **isinya tidak terekspor** |
| `When` **75** | ⚠️ sisiran pola | ⭐ **75/75 syaratnya dibaca** | ✅ syarat terbaca; ⛔ arti kode belum |
| `Section` **25** + `Harness` **6** | ⛔ **nol dibuka** | ⭐ **31/31 disisir**, 11 murni tata letak | ⛔ **nol dibaca utuh** |
| `Activity` **92** | ⚠️ sisiran pola | ⭐ **30 hitung uang disisir mendalam** | ⛔ **62 belum disisir mendalam**; ⛔ **nilai langkah tidak terekspor untuk seluruh 92** |
| `DataTransform` **12** · `FlowAction` **9** · `ReportDefinition` **15** · `Flow` **1** | ⛔ nol dibuka | ⛔ **nol dibuka** | ⛔ **37 berkas MASIH nol dibuka** |
| `.xlsx` **2** | ⛔ tidak dibuka, sengaja | ⛔ sama | ⭐ artefak turunan |

⛔⛔ **Yang MASIH nol dibuka: 37 berkas** — seluruh `DataTransform` **12**, seluruh `FlowAction`
**9**, seluruh `ReportDefinition` **15**, dan `Flow` **1** *(⭐ yang terakhir sudah ditelusur di
D2, jadi bukan lubang)*.
⚠️ **Dan `DataTransform` 12 adalah lubang yang nyata** — ⭐ di modul lain jenis aturan inilah yang
memindahkan tahap kasus.

⭐ **Dibaca UTUH sepanjang dua ronde: 8 berkas** *(5 blok PL/SQL + 2 `DecisionTable` + 1 `When`
`isApproved`)*. ⛔ **Nol dari 31 berkas layar dibaca utuh** — aturan §4 butir 4 dipatuhi.

### Rekap butir

| | Jumlah |
| --- | ---: |
| butir **lahir** di ronde 2 | ⭐ **13** |
| butir ronde 1 yang **berubah status** *(diperdalam/dipersempit)* | **8** |
| ⛔ butir **ditutup** | ⛔ **0** |
| ⭐ `[Finance]`: ronde 1 **0** → ronde 2 | ⭐ **3** |

### Apakah frontier masih terbuka?

⛔⛔ **MASIH TERBUKA — dan ronde 2 mengubah SIFATNYA, bukan hanya mengurangi ukurannya.**

⭐ Ronde 1 berhenti karena **saya belum membaca cukup**. ⛔ **Ronde 2 berhenti karena bahannya
tidak ada.** Dua batas itu berbeda jenis, dan yang kedua **tidak dapat diatasi dengan ronde
tambahan**:

1. ⛔⛔ **942 langkah penghitungan tanpa nilai** — bukan belum dibaca, **tidak terekspor**.
2. ⛔⛔ **Dua tabel keputusan tanpa baris** — sama.
3. ⛔ **35 % penggolong label-nya berselisih dengan aturannya** — terbaca, tetapi **mana yang
   benar bukan pertanyaan korpus.**

⭐ **Yang MASIH dapat dimajukan tanpa menunggu siapa pun:** 62 `Activity` non-uang, 12
`DataTransform`, 9 `FlowAction`, 15 `ReportDefinition` — ⭐ **98 berkas**.

### Syarat ronde 3 produktif — dan kapan modul ini siap `to-spec`

| # | Syarat | Sifat |
| --- | --- | --- |
| ⭐ **1** | **12 `DataTransform` + 9 `FlowAction` dibaca** | ⭐ **tidak menunggu siapa pun** — dan `DataTransform` adalah tempat tahap kasus biasanya berpindah |
| **2** | **62 `Activity` non-uang disisir** | tidak menunggu siapa pun; ⚠️ nilainya terbatas selama butir #1 berdiri |
| ⛔ **3** | **P18 dijawab** — isi 942 langkah | ⛔⛔ **MENAHAN**, dan ⭐ **ronde 3 tidak dapat menggantinya dengan pembacaan** |
| ⛔ **4** | **P19 dijawab** — baris dua tabel keputusan | ⛔ **MENAHAN** gerbang pertama alur |
| ⛔ **5** | **P20 dijawab** — 20 nomor polis dipaku | ⛔ **MENAHAN** paritas perhitungan |

### ⭐ Kesiapan `to-spec` — dinyatakan terang-terangan

⛔⛔ **Modul ini BELUM siap masuk `to-spec`, dan syaratnya LEBIH BANYAK dari yang brief
perkirakan.**

⭐ Brief menawarkan rumusan: *"siap `to-spec` begitu OQ-025 dan naskah tiga procedure tersedia."*
⛔ **Berdasarkan ronde 2, rumusan itu belum cukup.** ⭐ **Tiga syarat, bukan dua:**

| Syarat | Keadaan |
| --- | --- |
| ⛔ **OQ-025** — implementasi efek keluar terakhir | ⭐ tidak berubah, tetap blocker |
| ⛔ **Naskah tiga stored procedure** | ⭐ tidak berubah, nol badan ada |
| ⛔⛔ **BARU — isi 942 langkah penghitungan** *(P18)* | ⛔ **ronde 2 menemukan ini, dan ronde 1 tidak dapat mengetahuinya** |

⚠️ **Sebab syarat ketiga wajib ada:** ⭐ naskah procedure memberi tahu apa yang terjadi
**di dalam basis data**. ⛔ **Tetapi nilai yang DIKIRIM ke procedure ditetapkan oleh langkah
`Property-Set` yang tidak terekspor.** ⭐ Jadi walau ketiga naskah procedure diserahkan besok,
⛔ **apa yang masuk ke dalamnya tetap tidak diketahui** — dan sebuah spesifikasi perhitungan uang
**tidak dapat ditulis dari setengah rantai.**

⛔ **Saya berhenti di sini.** Ketiga syarat itu **keputusan pemilik peran**, bukan saya.
