# Rancangan tabel datar — NB Treaty In **dan** EDM Treaty In

> ⛔ **`T_POLIS_BREAKDOWN_SPREAD` DIBATALKAN 23-09-2026.** `[keputusan work owner]` Tabel itu **tidak ada**. Keempat medannya turunan: dua dari master `POOLDATA.PROPORTIONALARRG`, dua dihitung dari `TotalPremium` dan `TotalClaim` yang sudah tersimpan. Rujukan di bawah dicoret, bunyinya tidak dihapus.

---


*Disusun 22 September 2026. Satu rancangan untuk dua modul, sebab keduanya menulis ke penyimpanan
yang sama lewat `POOLDATA.PEGA_JSON_POLIS_TREATYIN`.*

⛔ **Bukan DDL.** Rancangan, bukan naskah. Presisi fisik dicocokkan DBA di dalam tiket.
⛔ **Nol nama orang, nol nomor polis harfiah, nol contoh JSON tersimpan.**

---

## 0 · Apa yang dimaksud "flat table" di sini

Dua bacaan mungkin. Yang saya pakai: **bacaan (a)**.

| | Bacaan | Akibat |
| :---: | --- | --- |
| **(a)** ⭐ | berhenti menyimpan dokumen; simpan sebagai **kolom relasional** | dipilih — sejalan dengan meninggalkan `DATA_JSON` dan meninggalkan pemanggilan procedure |
| (b) | **satu tabel lebar tunggal** | ditolak — daftar berulang *(angsuran, layer XOL, spreading)* akan menggandakan baris induk dan merusak penjumlahan uang |

⭐ Tetapi bacaan (a) **tidak berarti banyak tabel**. Sapuan korpus menemukan bahwa sarangnya
**serupa-diri**, sehingga jumlah tabel dapat ditekan sampai **enam**. Lihat Bab 3.

---

## 1 · Dua temuan sapuan yang mengubah rancangan

Disapu **9.430 berkas** seluruh korpus, ditambah 322 jalur properti unik dari `NB Treaty In` +
`EDM Treaty In`. Perintahnya di Bab 7.

### 1.1 ⭐⭐ Pohon `LocationList` **bukan milik Treaty** — ia milik Fac In

Rancangan lama *(prompt yang ditarik 22 September)* memperkirakan pohon **lima tingkat**:
`LocationList → OccupationList → AnekaList → CoverageList → ClauseList/DeductibleList`.
**Perkiraan itu salah untuk Treaty In.** `[terverifikasi]`

| Simpul | NB Treaty In | EDM Treaty In | Pemilik sebenarnya |
| --- | ---: | ---: | --- |
| `LocationList` | ⭐ **0** | ⭐ **0** | NB FacIn 1.527 · RNW Fac In 1.520 · Endorsment Fac In 1.389 |
| `OccupationList` | **0** | **0** | Endorsment Fac In 693 · NB FacIn 631 |
| `AnekaList` | **0** | **0** | Endorsment Fac In 897 · NB FacIn 811 |
| `CoverageList` | **0** | **0** | Endorsment Fac In 2.864 · NB FacIn 2.343 |
| `ClauseList` | **0** | **0** | Endorsment Fac In 343 · NB FacIn 250 |
| `DeductibleList` | 8 | 0 | Endorsment Fac In 202 · NB FacIn 182 |

⇒ ⭐ **Pohon lima tingkat hilang dari lingkup ini.** Rancangan Treaty In tidak perlu menampungnya.

⛔ `[terbuka]` **Tetapi satu hal harus dipastikan work owner.** Contoh JSON pertama yang Anda
berikan memuat pohon `LocationList`. Karena `POOLDATA.json_polis` **dipakai bersama** oleh NB Treaty
In, NB FacIn, dan EDM Treaty In, contoh itu **mungkin baris milik Fac In, bukan Treaty In**.
**Perlu dipastikan sebelum migrasi data lama dijalankan** — bila salah tebak, pemecah dokumen akan
diuji terhadap bentuk yang bukan bentuknya.

### 1.2 ⭐⭐ Sistem lama menyimpan **bentuk yang sama tiga kali**

| Salinan | Tempatnya | Jalur | Modul |
| --- | --- | ---: | --- |
| **sekarang** | `PolicyTreatyIn.*` | 94 di tingkat atas | NB + EDM |
| **lama, beku** | `PolicyTreatyIn.OldData.*` | **93** | EDM *(ditulis juga saat realisasi NB)* |
| **selisih** | `TreatyDifference.*` + `TreatyXOLDifferenceList` | 27 + 24 | ⭐ EDM 199 dari 201 |

Ketiganya memakai **nama medan yang sama persis** — `PremiOgp` `ResultOgp1` `Claim` `NetPremium`
`Deduction1` `TotalClaim` `TotalSharePercentagePremium` dan seterusnya. `[terverifikasi]`

⭐⭐ **Inilah penyederhanaan terbesar rancangan ini.** Dengan kunci **generasi**
`(NOPOLIS, PRODKE)`, ketiganya runtuh menjadi **satu** keluarga tabel:

| Yang lama | Yang baru |
| --- | --- |
| `OldData` — salinan beku | ⇒ **baris pada `PRODKE − 1`**; beku karena baris generasi lampau **tidak boleh disunting** |
| `TreatyDifference` — selisih | ⇒ **dihitung**, bukan disimpan |
| `TreatyXOLDifferenceList` | ⇒ **dihitung** |

⇒ ⭐ **Sekitar 144 kolom cermin lenyap dari rancangan**, dan butir `[terbuka]` *"`OldData` di dalam
`OldData`"* **hilang dengan sendirinya** — tidak ada lagi salinan untuk disarangkan.

`[terverifikasi]` Sapuan menegaskan susunan bersarang itu memang **tidak pernah dibaca**: jalur
berbentuk `OldData…OldData` = **0** dari 322.

⛔ `[terbuka]` **Selisih: dihitung atau disimpan?** Saran saya **dihitung**, sebab selisih yang
disimpan membawa galat presisi milik kedua operannya dan menumpuk pada endorsemen berlapis *(P57)*.
⛔ **Pemutusnya work owner**, dan terikat pilihan (a)/(b)/(c) pada **P29**.

---

## 2 · Yang **tidak** dirancang ulang — sudah relasional

| Tabel | Isi | Perlakuan |
| --- | --- | --- |
| `POOLDATA.TREATY_IN` | `ID` + **20 kolom datar** kontrak | ⭐ **dipertahankan apa adanya** |
| `POOLDATA.ACHIEVEMENT` | 18 kolom pencapaian | dipertahankan; ⛔ penjaga duplikatnya diperbaiki *(Bab 5)* |
| `POOLDATA.JSON_POLIS_MONITORING` | arsip pesan yang dikirim ke produksi | ⭐ **dokumennya dipertahankan** — ini arsip pesan keluar, sah berbentuk dokumen |
| `POOLDATA.M_TREATY_IN` | `ID`, `JSONDATA` | ⛔ `[terbuka]` — 20 kolom `TREATY_IN` sudah memuat kontraknya; perlu dipastikan apakah `JSONDATA` membawa sesuatu yang tidak ada di sana |

Kedua puluh kolom `TREATY_IN`: `PROPORTIONTYPE` `TREATYCONTRACTNAME` `TERITORIALSCOPE`
`COMMENCEMENT` `TERMINATION` `CLASSOFBUSINESS` `LEADINGREINSSOURCE` `LEADINGREINSSOURCEID` `CEDING`
`CEDINGID` `LEADINGREINSID` `NUSARESHAREPCT` `BROKERAGEPCT` `INFORMATION` `POSITIONUSERNAME`
`POSITION` `STATUSAKSEPTASI` `CHOOSESTATUSAKSEPTASI` `TREATYYEAR`.

⭐ **Sasaran penggantian hanya satu kolom:** `POOLDATA.json_polis.DATA_JSON`.
Tujuh kolom lainnya — `IDPEGA` `TGL_INPUT` `NOPOLIS` `NOENDORS` `PRODKE` `TGL_PROD` `USERNAME` —
**sudah datar dan pindah apa adanya** menjadi kolom tabel inti.

---

## 3 · Enam tabel

```
T_WORK_POLIS                        akar, LINTAS-LINI          (nama ditetapkan work owner)
  +- T_GENERAL_POLIS                inti, satu baris per generasi
       |- T_POLIS_QUOTATION         QuotationData              1:1
       |    +- T_POLIS_CEDING       QuotationData.CedingCoList 1:N
       |- T_POLIS_INSTALMENT        ListInstallment            1:N   induk
       |    +- T_POLIS_INSTALMENT_DETAIL   .InstallmentList    1:N   anak
       |- T_POLIS_SPREADING         SpreadingRiskList          1:N
       |- T_POLIS_XOL               TreatyXOLList              1:N   NONPROP saja
       |    +- T_POLIS_XOL_LAYER    .ValueList                 1:N   NONPROP saja
       +- POOLDATA.HISTORYAKSEPTASIPRODUCTION   SuggestList    1:N   SUDAH datar
```

⚠️ **Nama tabel adalah usulan**, tunduk pada standar penamaan DBA.

### 3.1 ⭐⭐ Sarang serupa-diri — kenapa dua tingkat cukup satu tabel

Dua daftar di dalam dokumen bersarang ke dalam dirinya sendiri, **dengan medan yang sama persis**:

| Induk | Anak | Medan sama? |
| --- | --- | :---: |
| `ListInstallment[]` | `ListInstallment().InstallmentList[]` | ⭐ **ya** — `InstallmentNo` `InstallmentPercentage` `Premium` `Currency` `IDCurrency` `DueDate` `PaymentTotal` |
| `TreatyXOLList[]` | `TreatyXOLList().ValueList[]` | ⭐ **ya** — `Layer` `LayerType` `LayerPart` `LayerPartType` `Currency` `IDCurrency` `GrossPremi` `NetPremi` `Deduction` `DueTo` `DueToValue` |

⇒ ⭐ **Satu tabel dengan rujukan-diri `INDUK_ID`** menampung **kedua bentuk** tanpa penanda apa pun.
Dokumen berbentuk datar menghasilkan baris ber-`INDUK_ID` kosong; dokumen bersarang menghasilkan
baris anak yang menunjuk induknya.

⭐⭐ **Ini jawaban langsung atas bahaya yang dicatat berkas keadaan** — *"pengurai yang menganggap
bentuknya seragam akan patah"*. Dengan rancangan ini pengurai **tidak perlu tahu bentuk mana yang
datang**. `[terverifikasi]` 58 rujukan bentuk bersarang di 6 berkas, 101 rujukan bentuk datar di 11.

---

### 3.2 ⛔⛔ RALAT 23 September 2026 — "serupa-diri" **DITARIK**, keduanya dipecah jadi dua tabel

> Bunyi Bab 3.1: ⛔ *"Satu tabel dengan rujukan-diri `INDUK_ID` menampung **kedua bentuk** tanpa
> penanda apa pun."*

⛔ **Premisnya keliru.** Work owner mempertanyakannya, dan uji ulang membuktikan keduanya
**bukan dua bentuk alternatif dari hal yang sama**, melainkan **dua tingkat 1:N sungguhan**.

#### Uji: berapa berkas memakai kedua tingkat sekaligus?

| | `ListInstallment` / `InstallmentList` | `TreatyXOLList` / `ValueList` |
| --- | ---: | ---: |
| berkas memakai **kedua** tingkat | **7** | **9** |
| berkas memakai **induk saja** | 7 | ⭐ **0** |
| berkas memakai **anak saja** | ⭐ **0** | 1 |

`[terverifikasi]` Kalau keduanya hanya dua bentuk alternatif, **tidak akan ada** berkas yang
memakai keduanya sekaligus. Ada tujuh dan sembilan.

#### Dan medannya tidak bermakna sama

**Angsuran** — nama kolom sama, bobotnya berlawanan:

| Medan | induk | anak |
| --- | ---: | ---: |
| `PaymentTotal` | ⭐ **19** | 1 |
| `InstallmentPercentage` | 27 | 13 |
| `Premium` | 17 | 12 |
| `PaymentDate` | **0** | 1 |

⇒ induk memegang **total**, anak memegang **baris angsurannya**.

**XOL** — lebih tegas lagi, ada medan yang **sama sekali tidak ada** di induk:

| Medan | induk | anak |
| --- | ---: | ---: |
| `Layer` · `LayerType` · `LayerPart` · `LayerPartType` | ⛔ **0** | **12** masing-masing |

⇒ ⭐ **Penanda layer hanya hidup di `ValueList`.** `TreatyXOLList` bukan daftar layer; ia satu
tingkat di atasnya.

#### ⛔ Kenapa satu tabel berbahaya, bukan sekadar kurang rapi

Baris induk dan baris anak memakai nama kolom yang sama dengan makna berbeda. Di satu tabel,
`SUM(PAYMENT_TOTAL)` akan **menjumlahkan total bersama rinciannya** — uang terhitung dua kali.
Satu query lupa menyaring `PARENT_ID IS NULL`, dan angkanya salah **tanpa pesan galat apa pun**.

⚠️ Pertanyaan *"bisa ditarik ulang jadi view?"* **bukan pembedanya** — kedua rancangan bisa.
Pembedanya: dengan dua tabel, salah jumlah **tidak mungkin terjadi karena lupa menyaring**.

#### Rancangan yang berlaku

| Tabel | Dari | Peran |
| --- | --- | --- |
| `T_POLIS_INSTALMENT` | `ListInstallment()` | induk — memegang total |
| `T_POLIS_INSTALMENT_DETAIL` | `ListInstallment().InstallmentList()` | anak — baris angsuran, punya `PAYMENT_DATE` |
| `T_POLIS_XOL` | `TreatyXOLList()` | induk — uang, tanpa penanda layer |
| `T_POLIS_XOL_LAYER` | `TreatyXOLList().ValueList()` | anak — `LAYER*` + uang |

⭐⭐ **TERTUTUP 23 September 2026 — TIDAK.** Work owner memastikan dari **data produksi** bahwa
dokumen proporsional berhenti di `.ListInstallment`, tanpa sarang `.InstallmentList`.

> Bunyi lama, dikutip: ⛔ *"`[terbuka]` Apakah tingkat anak angsuran hidup juga pada dokumen
> proporsional. Korelasinya kuat dengan NonProp … tetapi `DataTransform\SetInstallmentValue`
> memakai kedua tingkat tanpa penanda bentuk sama sekali."*

`[terverifikasi]` Croscheck korpus mendukungnya. Tujuh berkas NB/EDM Treaty In menyentuh
`.ListInstallment().InstallmentList()`, dan **ketujuhnya kini terjelaskan**:

| Berapa | Berkas | Kenapa bukan proporsional |
| ---: | --- | --- |
| **5** | `InputPolicyTreatyInDetail_NonProp` · `InputPolicyTreatyOutDetail_NonProp` · `InputPolicyTreatyEDMDetail_NP` *(NB dan EDM)* · `InputPolicyTreatyEDMDetail_NP_AdjPremi` | bernama `_NonProp` / `_NP` |
| **1** | `FillPaymentInstallmentEDMT` | keterangannya sendiri berbunyi *"Generate Installment based on **XOLDifferenceList**"* dan *"Loop per currency List in XOLDifferenceList"* ⇒ jalur XOL |
| **1** | `SetInstallmentValue` | dipanggil `FillMasterInstallment`, dipanggil `EDMChooseBusiness_Act`. Rantai ini **menyalin dari master kontrak** `pyWorkPage.TreatyIn.Installment(n).InstallmentList`, bukan membuat sarang baru. Kelasnya `ASM-FW-GISFW-Int-treaty_in_edm` — lapisan integrasi, jalur EDM |

⇒ ⭐ **Nol bukti korpus bahwa dokumen proporsional membawa `.InstallmentList`.**

⭐ Dan ini menerangkan **asal-usul** sarang itu: jadwal angsuran **master kontrak** memang
bertingkat dua. `T_POLIS_INSTALMENT_DETAIL` karena itu **hanya hidup pada bentuk non-proporsional**,
dan pada bentuk proporsional bernilai nol baris — bukan kolom kosong, melainkan tidak ada barisnya.

⛔ `[terbuka]` Nama `T_POLIS_XOL_LAYER` adalah usulan saya, bukan penamaan work owner.

---

## 4 · Kolom

Penanda asal: `[dari DATA_JSON]` · `[dari json_polis]` · `[dari TREATY_IN]` · `[baru]`

### 4.1 `TREATY_IN_POLIS` — inti

**Kunci dan generasi**

| Kolom | Tipe | Asal | Catatan |
| --- | --- | --- | --- |
| `ID` | angka, surrogate | `[baru]` | kunci utama |
| `NOPOLIS` | teks | `[dari json_polis]` | ⭐ bersama `PRODKE` membentuk kunci alami |
| `PRODKE` | angka | `[dari json_polis]` | ⭐ **0** = polis baru · **1..n** = endorsemen ke-n |
| `NOENDORS` | teks | `[dari json_polis]` | |
| `EDM_NO` | teks, **terhitung** | `[dari DATA_JSON]` | `NOPOLIS + "/E" + lpad(PRODKE,2,'0')` — **P55** |
| `IDPEGA` | teks | `[dari json_polis]` | ⛔ `[terbuka]` jembatan migrasi; perlu diputuskan apakah bertahan sesudah Pega mati |
| `TGL_INPUT` `TGL_PROD` | tanggal | `[dari json_polis]` | |
| `USERNAME` | teks | `[dari json_polis]` | ⭐ identitas akses login — **P4** |

⭐ **Baris generasi lampau tidak boleh disunting.** Itulah yang menggantikan pembekuan `OldData`
*(P58)*, dan itulah yang membuat selisih *(P57)* dapat dihitung ulang kapan saja.

**Penentu bentuk** — wajib ada, sebab bentuk dokumen berbeda menurutnya *(Bab 4.4 spec)*

| Kolom | Nilai yang diketahui | Rujukan korpus |
| --- | --- | ---: |
| `PROPORTIONAL_TYPE` | `Proportional` · `NonProportional` | **170** |
| `IS_NEW_POLICY_NON_PROP` | `0` · `1` | **87** |
| `EDM_TYPE` | jenis endorsemen, termasuk pembatalan — **P56** | 29 |

**Penanda alur** `[dari DATA_JSON]`

`IS_APPROVED` *(⭐ teks, bukan angka — P24/P6: `0` = ditolak, selain itu = disetujui)* ·
`IS_APPROVED_TO_DEPT_HEAD` · `FLAG_PPH` · `FLAG_RETRO_TREATY` · `HAS_FAC_OUT` · `IS_EDM_INPUT_ON_NB`

**Tanggal** `[dari DATA_JSON]`

`START_DATE` · `END_DATE` · `PRODUCTION_DATE` · `STATEMENT_DATE` · `SUGGEST_DATE`

⚠️ **Dokumen lama menyimpan dua format dalam satu berkas** — `YYYYMMDD` dan cap waktu Pega
bersufiks ` GMT`. Pemecah wajib menerima keduanya. `[terverifikasi]`

**Uang** — seluruhnya angka berpresisi tetap, ⛔ **tidak pernah `float`** *(ADR-0003)*

`GROSS_PREMIUM` · `NET_PREMIUM` · `TOTAL_PREMIUM` · `TOTAL_CLAIM` · `CLAIM` · `OUTSTANDING_CLAIM` ·
`PREMI_OGP` · `PREMI_ONP` · `RI_COMM_OGP` · `RI_COMM_ONP` · `OVERIDDING_COMM_OGP` *(⚠️ ejaan
sistem lama dipertahankan)* · `OVERIDDING_COMM_ONP` · `RESULT_OGP1` · `RESULT_OGP2` · `RESULT_ONP1` ·
`RESULT_ONP2` · `SALVAGE_VALUE` · `EXCESS_LOSS` · `BALANCE_DUE_TO` · `BALANCE_BEFORE_TAX` ·
`BALANCE_BEFORE_PPH` · `PPN_VALUE` · `PPH_VALUE` · `SHARE_VALUE` · `DUE_TO`

**Persen** — ⛔ **bukan uang**; `12.5` berarti 12,5 persen *(ketetapan 2, P29)*

`DEDUCTION1` · `DEDUCTION2` · `TOTAL_SHARE_PERCENTAGE_PREMIUM` · `TOTAL_SHARE_PERCENTAGE_CLAIM`

**Kode dan teks** `[dari DATA_JSON]`

`POLICY_NO` · `NO_OFFER` · `BIZ_CODE` · `BIZ_NAME` · `OJK_BUSINESS_ID` · `ID_NEW_BISNIS` ·
`SOB` · `SOB_NAME` · `TREATY_TYPE` · `TREATY_YEAR` · `TREATY_GROUP_ID` · `TREATY_GROUP_NAME` ·
`TREATY_GROUP_OLD_ID` · `MASTER_ID` · `CURRENCY` · `ID_CURRENCY` · `TYPE_TAX` · `QUARTAL` ·
`YEAR_OF_QUARTAL` · `STATEMENT_TYPE` · `CLAIM_TYPE` · `CLAIM_PAYMENT_TYPE` · `INSTALLMENT` ·
`LAYER` · `LAYER_TYPE` · `LAYER_PART` · `LAYER_PART_TYPE` · `SUGGEST`

**Pihak** — ⛔ **nilai berupa nama orang tidak disalin ke berkas mana pun**; kolomnya tetap ada

`CEDING_CO` · `CEDING_CO_NAME` · `INSURED_ID` · `INSURED_NAME` · `MARKETING_OFFICER` ·
`OPERATOR_NAME`

⚠️ `CEDING_CO` di dokumen **berakhiran `"; "`** — daftar yang digabung, bukan nilai tunggal.
⛔ `[terbuka]` dipecah jadi baris tersendiri atau dipertahankan sebagai teks? Lihat Bab 6.

**Dari `QuotationData{}`** — ⭐ **10 medan** *(RALAT: semula tertulis 4)*, dilipat ke tabel inti,
tidak perlu tabel sendiri

`PROPORTIONAL_TYPE` *(di atas)* · `MO_ID` · `BUSINESS_CODE` · `BUSINESS_OLD_ID` · `GROUP_PANEL` ·
`SOURCE_OF_BUSINESS` · `TYPE` · `EDM_TYPE` · `OLD_POLICY_NO` · `MARKETING_NAME`

⭐ ⚠️ **RALAT.** Masukan penggolong `BusinessType_DeT` adalah **`GROUP_PANEL` + `BUSINESS_OLD_ID`**,
bukan `BUSINESS_CODE` + `BUSINESS_OLD_ID` seperti tertulis semula. `[terverifikasi]` — sejalan dengan
uji terhadap data produksi: `GroupPanel=006` + `BusinessOldId=01` → baris 31 → `"FireStyle2"`.
36 baris, 128 kode, bawaan `"UNKNOWN"`, berhenti di baris pertama yang cocok.

⚠️ `OLD_POLICY_NO` dan `MARKETING_NAME` baru terlihat pada sapuan ini — keduanya **tidak** masuk
hitungan 79 medan tingkat atas, sebab letaknya di dalam `QuotationData`.

⛔ **Tidak dimigrasi:** `pxObjClass` *(internal Pega, ada di setiap simpul)* · `Show` · `ViewState` ·
`pxResults` · `FillPaymentInstallmentEDMT` — keadaan layar, bukan data.

### 4.2 `TREATY_IN_POLIS_XOL` — layer non-proporsional

| Kolom | Catatan |
| --- | --- |
| `ID` · `POLIS_ID` | induk ke tabel inti |
| ⭐ `INDUK_ID` | **rujukan-diri**; kosong = baris tingkat atas |
| `LAYER` `LAYER_TYPE` `LAYER_PART` `LAYER_PART_TYPE` | penanda layer |
| `CURRENCY` `ID_CURRENCY` | mata uang |
| `GROSS_PREMI` `NET_PREMI` `DUE_TO_VALUE` `DUE_TO` | uang |
| `DEDUCTION` | ⚠️ `[dugaan]` — perlu dipastikan persen atau uang; namanya sama dengan `DEDUCTION1/2` yang **persen** |
| `BROKERAGE_FEE_SEBENARNYA` `PPH_VALUE` `PPN_VALUE` | uang · ⭐ **P46**: `nilai / 1,022`, dibulatkan 4 desimal |
| `NET_PREMI_AFTER_PPH` `NET_PREMI_AFTER_PPN` `NET_PREMI_AFTER_TAX` | uang |

Dipakai bila `PROPORTIONAL_TYPE = 'NonProportional'`. 169 rujukan `TreatyXOLList`, 109 rujukan
bentuk bersarang `ValueList`.

### 4.3 `TREATY_IN_POLIS_ANGSURAN` — jadwal angsuran

| Kolom | Catatan |
| --- | --- |
| `ID` · `POLIS_ID` | |
| ⭐ `INDUK_ID` | **rujukan-diri** — menampung bentuk datar **dan** bersarang |
| `INSTALLMENT_NO` | urutan |
| `DUE_DATE` `PAYMENT_DATE` | tanggal |
| `INSTALLMENT_PERCENTAGE` | ⭐ **persen** |
| `PREMIUM` `PAYMENT_TOTAL` | uang |
| `PREMIUM_AFTER_PPH` `PREMIUM_AFTER_PPN` `PREMIUM_AFTER_TAX` | uang |
| `CURRENCY` `ID_CURRENCY` | |

### 4.4 `TREATY_IN_POLIS_SPREADING` — sebaran risiko

`ID` · `POLIS_ID` · `TREATY_NAME` · `TREATY_TYPE` · `CURRENCY` · `CURRENCY_ID` ·
`SHARE_PERCENTAGE` *(persen)* · `SPLIT_RNM_SHARE_PCT` *(persen)* · `PREMIUM_SPREADED` *(uang)* ·
`CLAIM_SPREADED` *(uang)* · `CLAIM_PERCENTAGE` *(persen)*

⭐ **P60 — beda NB dan EDM ada di sini.** EDM menetapkan `SPLIT_RNM_SHARE_PCT = 100` pada baris
pertama lalu membagi dengan presisi **20**; NB membagi `100 / jumlah baris` dengan presisi **10**.
⛔ **Ditiru apa adanya**, dan perbedaannya wajib punya test sendiri.

⚠️ Sapuan menemukan `SpreadingRiskList` **hanya di NB** *(65 rujukan, EDM 0 di tingkat itu)*;
EDM mencapainya lewat `TreatyDifference.SpreadingRiskList` *(12)*. Konsisten dengan EDM yang
bekerja pada selisih.

### 4.5 `TREATY_IN_POLIS_USULAN` — usulan dan persetujuan per baris

`ID` · `POLIS_ID` · `SUGGEST` · `SUGGEST_DATE` · `OPERATOR_NAME` · `IS_APPROVED`

⭐ ⚠️ **`IS_APPROVED` di sini berbeda dari `IS_APPROVED` tingkat atas.** Yang satu keputusan per
usulan, yang satu keputusan atas polis. Menggabungkan keduanya adalah cacat.
`[terverifikasi]` `SuggestList` muncul 1.849 kali korpus-wide, dominan di keluarga Fac In — di NB
Treaty In hanya 5 rujukan, ⛔ sehingga **cakupan medannya di Treaty In belum lengkap terukur**.

### 4.6 `TREATY_IN_POLIS_NILAI_LAIN` — penampung sementara

⛔ **Ini bukan bagian rancangan akhir.** Ia ada supaya migrasi tidak kehilangan medan yang belum
tergolong, dan **wajib kosong sebelum rancangan dinyatakan selesai**.

`POLIS_ID` · `JALUR` · `NILAI`

Kenapa perlu: sensus medan yang dipakai rancangan ini berasal dari **rujukan di dalam aturan**, dan
aturan hanya menyentuh medan yang dihitung atau ditampilkan. **Medan yang hanya lewat — masuk
dokumen lalu keluar tanpa disentuh — tidak terlihat oleh sapuan ini.** Tiga contoh nyata memberi
37, 64, dan 47 medan skalar tingkat atas dengan gabungan **74**, sementara sapuan korpus memberi
**86 skalar daun** — ⛔ **kedua angka tidak dapat dipertemukan tanpa `JSON_DATAGUIDE`.**

---

## 4bis · ⛔⛔ RALAT 22 September 2026 — dua tabel usulan ternyata **sudah ada**

Sapuan naskah SQL pada `RDBList` *(36 berkas EDM, 41 NB, seluruhnya membawa `pyBrowseSQL`)* menemukan
bahwa sebagian pekerjaan yang saya rancang **sudah dikerjakan sistem lama secara relasional**.
⛔ **Bunyi lama dikutip, tidak dihapus.**

### 4bis.1 `TREATY_IN_POLIS_USULAN` — ditarik

> Bunyi Bab 4.5: *"`TREATY_IN_POLIS_USULAN` — usulan dan persetujuan per baris ·
> `SUGGEST` `SUGGEST_DATE` `OPERATOR_NAME` `IS_APPROVED`"*

⛔ **Tabel ini tidak perlu dibuat.** Riwayat persetujuan sudah tersimpan relasional di
`POOLDATA.historyakseptasiproduction` — **15 kolom**, ditulis `RDBList\InsertViewSuggest_SQL`.

`[terverifikasi]` Sapuan anggota `SuggestList` di NB + EDM Treaty In menghasilkan **nol** medan
skalar. ⇒ `SuggestList` di dalam dokumen adalah **daftar tampilan**, bukan penyimpanan.

Pemetaan kolomnya *(lengkap, dari naskah SQL)*:

| Kolom | Asal |
| --- | --- |
| `IDPEGA` | `pyWorkPage.pzInsKey` |
| `PIC` | `InputData.CARI4` — ⭐ nama tampilan, sejalan **P33** |
| `APPROVAL` | `InputData.CARI9` — keputusan aksep/tolak |
| `KETERANGAN` | `substr(InputData.CARI10, 0, 3990)` — ⚠ **dipotong di 3990** |
| `AKSES_LOGIN` | `OperatorID.pyUserIdentifier` — ⭐ identitas login, sejalan **P4** |
| `B2B` | `pyWorkPage.OfferFacIn.IsB2B` |
| `PERCENT_RNM` | `pyWorkPage.OfferFacIn.PercentShare` |

⚠️ `B2B` dan `PERCENT_RNM` datang dari halaman **`OfferFacIn`** — bukti aturan ini **dipakai bersama
modul Fac In**, bukan milik Treaty In sendiri.

### 4bis.2 ⭐⭐ Tabel datar yang diminta **sudah ada**: `TREATYINPRODUCTION`

`RDBList\InsertTreatyInProdEDMT_SQL` menyisip ke `TREATYINPRODUCTION` dengan **58 kolom datar**,
dan **seluruh 58** sudah tertelusur ke propertynya — 28 langsung, 30 lewat parameter `InputTreaty.CARIn`
yang saya telusuri ke pemanggilnya `Activity\InsetTreatyInProdAddendum_Act`. `[terverifikasi]`

Kolomnya mencakup justru yang paling sulit: `PREMI_OGP` `RI_COMM_OGP` `OVERRIDING_COMM_OGP`
`RESULT*` `PREMI_ONP` `DEDUCTION1` `DEDUCTION2` `NET_PREMIUM` `BALANCE_DUE_TO` `PPHVALUE` `PPNVALUE`
`BALANCE_BEFORE_PPH` `BALANCE_BEFORE_TAX` `SHARE_VALUE` `OUTSTANDING_CLAIM` `LAYER*` `QUARTER`
`GUARANTEE_FUND`.

⇒ ⭐ **Rancangan tabel inti sebaiknya bertolak dari 58 kolom ini**, bukan dari daftar property
mentah — sebab inilah bentuk yang sudah diterima sistem produksi hilir.

Pemetaan penuh ada di sheet **`Treaty In Prop`** dan **`Treaty In NonProp`** pada
`Diagram-Skema-Tabel-NusantaraRe.xlsx`.

### 4bis.3 Empat temuan dari naskah SQL itu

1. ⛔⛔ **Uang menyeberang sebagai teks lalu diubah di dalam SQL:**
   `TO_NUMBER(REPLACE(InputTreaty.CARI32, ',', '.'))` — koma desimal diganti titik **di lapisan basis
   data**, pada **22 kolom uang**. Bila aplikasi baru mengirim format lain, konversinya diam-diam
   meleset. ⇒ Sistem baru **mengirim bilangan, bukan teks**; konversi lokal tidak boleh ada di SQL.
2. ⚠️ `[dugaan]` **`PCT_SHARE_PREMI` ← `.ClaimPercentage`** sedangkan **`PCT_SHARE_CLAIM` ←
   `.SharePercentage`** — ⛔ **tampak tertukar**. Perlu dipastikan work owner; bila memang keliru,
   data lama membawa kekeliruan itu.
3. ⚠️ `[dugaan]` **`STATEMENT_DATE` dan `BEGINDATE` dua-duanya dari `CARI11`** = `.StartDate`.
4. ⛔ `[terbuka]` **Kardinalitasnya belum pasti.** `.PremiumSpreaded` `.ClaimSpreaded`
   `.SharePercentage` adalah anggota `SpreadingRiskList`, sehingga `TREATYINPRODUCTION` tampaknya
   **satu baris per baris spreading**, bukan per polis. ⛔ Tidak saya putuskan.

### 4bis.4 Empat sasaran tulis langsung — seluruh SQL disapu

| Perintah | Tabel | Modul |
| --- | --- | --- |
| `INSERT` | `HISTORYAKSEPTASIPEGA` *(6 kolom)* | NB **dan** EDM, naskah identik |
| `INSERT` | `POOLDATA.HISTORYAKSEPTASIPRODUCTION` *(15)* | NB |
| `INSERT` | `TREATYINPRODUCTION` *(58)* | EDM |
| `UPDATE` | `POOLDATA.JSON_POLIS_MONITORING` | EDM — ⭐ inilah langkah 16 yang menyetel `STS_KONVERSI`, sejalan **P51** dan **P52** |

⚠️ **Dua dari empat ditulis tanpa awalan skema** — menambah bukti pelanggaran ketetapan 6.

---

## 4ter · ⭐ KEPUTUSAN 23 September 2026 — `OldData` dan `TreatyDifference`

Empat keputusan work owner, diambil berurutan dalam satu sesi. **Mengikat.**

### 4ter.1 `OldData` → penunjuk, bukan tabel

⭐ **`T_GENERAL_POLIS.OLD_POLIS_ID` → `T_WORK_POLIS.ID`** · nullable · kosong pada polis baru
*(PRODKE 0)*, terisi pada tiap endorsemen. `[keputusan work owner]`

Sistem lama **tidak punya penunjuk ini** — ia mencari tiap kali:

```sql
-- EDM\RDBList\FetchPolisJsonPolis
Select DATA_JSON from POOLDATA.JSON_POLIS
 where NOPOLIS = ?  order by PRODKE desc  fetch first 1 row only
```

⚠️⚠️ Dan pencarian keduanya membawa jebakan:

```sql
-- EDM\RDBList\SelectProdKe
select PRODKE from json_polis where substr(nopolis, 1, 24) = ? order by prodke desc
```

`substr(nopolis, 1, 24)` — nomor polis **dianggap selalu 24 karakter**. Penunjuk eksplisit
menghapus jebakan itu; FK tidak peduli panjang teks.

**Endorsemen berlapis** *(P57)* menjadi senarai berantai, bukan sarang:

```
ID=1  PRODKE=0  OLD_POLIS_ID=NULL      polis baru
ID=2  PRODKE=1  OLD_POLIS_ID=1         .../E01
ID=3  PRODKE=2  OLD_POLIS_ID=2         .../E02
```

⭐ Cacat `OldData` di dalam `OldData` **hilang sendiri** — lapisan itu tidak pernah lahir.

**Tiga aturan wajib:**

| Aturan | Cara |
| --- | --- |
| satu generasi hanya boleh punya satu penerus | ⭐ `UNIQUE (OLD_POLIS_ID)` — melarang percabangan |
| nomor generasi tidak boleh kembar | `UNIQUE (NOPOLIS, PRODKE)` |
| baris generasi lampau tidak boleh disunting | aturan di `services` |

⚠️ Aturan kedua sekaligus memperbaiki bahaya lama: `ProdKe` baru diambil dari `MAX + 1`, sehingga
dua endorsemen serentak pada polis yang sama akan membaca angka yang sama. Di sistem baru yang kedua
**gagal**, bukan bentrok diam-diam.

### 4ter.2 `TreatyDifference` → dihitung **di Go**, bukan di Oracle

`[keputusan work owner]` ⛔ **Bukan tabel sumber. Bukan view.**

> ⛔ **RALAT.** Saya sempat menyarankan `V_POLIS_DIFFERENCE`. **Ditarik.** Work owner menunjukkan
> bahwa rumusnya ada di Pega — lapisan aplikasi — bukan di Oracle. View berarti rumus ditulis
> **dua kali**, di Go dan di SQL, dan pasti bercabang.

Repository hanya mengambil **dua baris**: baris ini dan baris `OLD_POLIS_ID`. `services` yang
menurunkan selisihnya.

**Empat aturan turunannya** `[terverifikasi]` dari `EDMTCalculateTreatyDifference` dan
`CalculateDifferenceEDM_act`:

| Jenis | Rumus |
| --- | --- |
| uang | `baru.X − lama.X` |
| ⭐ persen | `baru.X` — **TIDAK dikurangi** |
| kunci *(`Layer*` `Currency` `TreatyName` `DueDate` `InstallmentNo`)* | `baru.X` |
| ⚠️ `DUE_TO` | `@If(selisih DueToValue > 0, "DUE TO US", "DUE TO YOU")` — dari **tanda**, bukan nilai |

⇒ Seluruhnya **fungsi murni** dari dua baris. Tidak ada satu medan pun yang butuh keadaan tersimpan.

### 4ter.3 `T_POLIS_DIFFERENCE` → tabel **proyeksi**, dibangun

`[keputusan work owner]` Pertanyaan *"siapa yang membaca angka selisih, lewat apa"* dijawab:
⭐ **layar aplikasi DAN SQL sendiri.**

⇒ Karena ada pembaca Oracle langsung, proyeksinya **wajib ada**. Kalau tidak, laporan mereka mati
pada hari Pega dimatikan.

⛔ **Bukan tabel sumber — tabel proyeksi.** Tiga aturan yang membedakannya:

1. **Hanya ditulis Go**, di dalam transaksi yang sama dengan generasinya *(P2)*. Tidak pernah diisi
   manual, tidak pernah lewat procedure.
2. **Boleh dihapus total dan dibangun ulang** kapan saja dari `T_GENERAL_POLIS`. Bila tidak bisa,
   ada yang salah.
3. **Bila isinya berbeda dari hasil hitung ulang, tabelnya yang salah** — bukan operannya. Sumber
   kebenaran tetap kedua baris generasi.

⇒ Galat uang tidak dibekukan: ia diperbaiki dengan membangun ulang proyeksinya, tanpa menyentuh
data polis.

⛔ `[terbuka]` Anak-anaknya — `_INSTALMENT`, `_SPREADING`, dan untuk nonprop `T_POLIS_XOL_DIFFERENCE`
— **belum dibangun**. Keduanya juga dapat dibangun ulang, jadi menundanya nyaris tanpa biaya.
Dibangun bila pembaca SQL memintanya.

⚠️ Karena pembacanya menulis query sendiri, **penamaan kolomnya mengikat**. Ikuti gaya
`TREATYINPRODUCTION` supaya terasa akrab, dan sertakan kunci yang mereka saring: `NOPOLIS` `PRODKE`
`EDM_NO` `IDPEGA` — supaya mereka tidak perlu join balik.

### 4ter.4 ⭐⭐ `T_POLIS_DIFFERENCE_ARSIP` → **tangkap saat migrasi, atau hilang selamanya**

`[keputusan work owner]` Inilah satu-satunya yang **tidak dapat dibangun ulang**.

| | Bisa dibangun ulang nanti? |
| --- | :---: |
| `T_POLIS_DIFFERENCE` *(proyeksi)* | ✔ ya |
| `T_POLIS_DIFFERENCE_ARSIP` *(selisih lama)* | ⛔ **TIDAK** |

Selisih yang dihitung sistem lama hanya ada di `DATA_JSON`, lengkap dengan galat presisinya —
dan justru galat itulah yang diperlukan supaya laporan lama dapat direkonsiliasi. Begitu migrasi
jalan dan dokumen lama tidak lagi jadi sumber, angka itu **hilang permanen**.

⇒ Baca-saja · tidak pernah diperbarui · tidak pernah dibangun ulang · hanya untuk rekonsiliasi.

⭐ **Berlaku apa pun pilihan presisi uang pada P29** — justru karena isinya angka lama apa adanya.

### 4ter.5 Pemeriksaan wajib saat migrasi

Untuk tiap endorsemen, bandingkan `OldData` yang tersimpan terhadap baris yang ditunjuk
`OLD_POLIS_ID` hasil penurunan.

⛔ Bila tidak cocok, itu **dilaporkan, bukan diam-diam dirapikan** — berarti pembekuan di sistem
lama terjadi pada saat yang berbeda, atau ada generasi yang hilang.

---

## 4quater · ⭐ KEPUTUSAN 23 September 2026 — ronde penutupan butir

Dua belas butir ditutup berurutan dalam satu sesi tanya-jawab. **Mengikat.**

### 4q.1 `[keputusan work owner]` P29 — presisi uang: **ikuti apa adanya**

Nilai lama dipindahkan **utuh termasuk ekor galatnya**. Tidak dibulatkan saat migrasi.

| Akibat yang mengikat | |
| --- | --- |
| skala kolom uang | **minimal 9 desimal** — `NUMBER(p,9)` atau lebih, presisi pasti dari DBA |
| penjaga duplikat `ACHIEVEMENT` | ⛔ **tidak boleh** berbasis nilai uang — kuncinya medan pengenal |
| pembandingan uang | ⛔ **tidak boleh sama-persis** — bertoleransi, atau dibandingkan terbulatkan |
| ADR-0003 | tetap: bukan `float`, melainkan desimal presisi tetap berskala besar |

⚠️ Galat itu lahir di rantai perhitungan **Pega**. Go memakai desimal, jadi perhitungan baru **tidak
menghasilkan galat yang sama**. Satu kolom akan memuat dua macam nilai — baris lama bergalat, baris
baru bersih. Kolom `SUMBER` membedakannya.

### 4q.2 ⛔⛔ RALAT — `PCT_SHARE_PREMI` **tidak** tertukar

> Bunyi laporan saya: ⛔ *"`PCT_SHARE_PREMI ← .ClaimPercentage` dan `PCT_SHARE_CLAIM ←
> .SharePercentage` tampak TERTUKAR."*

**Keliru, dan kekeliruannya milik metode saya.** Skrip penelusuran `CARI` mengumpulkan nilai dari
**seluruh berkas** lalu mengambil yang paling sering — bukan yang ada di aktivitas pengisinya.

`[terverifikasi]` Pemetaan sebenarnya, dari `InsetTreatyInProdAddendum_Act`
*(`pxUpdateDateTime` 2026-09-23 07:26)*:

```
PCT_SHARE_PREMI  <-  CARI15  <-  .SharePercentage     benar
PCT_SHARE_CLAIM  <-  CARI16  <-  .ClaimPercentage     benar
```

⭐ **Dan sapuan ulang itu memunculkan yang lebih penting: hampir setiap `CARI` punya dua varian** —
satu jalur proporsional, satu jalur XOL:

| Kolom | jalur prop | jalur XOL |
| --- | --- | --- |
| `PREMI_OGP` | `.PremiOgp` | `.GrossPremi` |
| `NET_PREMIUM` | `.PremiumSpreaded` | `.NetPremi` |
| `BALANCE_DUE_TO` | `.BalanceDueTo` | `.DueToValue` |
| `DEDUCTION1` | `.Deduction1` *(persen)* | `.Deduction` *(uang)* |
| `LAYERTYPE` `LAYER` `LAYERPARTTYPE` `LAYERPART` | ⭐ **`"0"`** | `.LayerType` dst. |

⛔ **Satu kolom `DEDUCTION1` menampung dua satuan.** Siapa pun yang membacanya lewat SQL **wajib
menyaring `PROPORTIONALTYPE` lebih dulu**, kalau tidak ia menjumlahkan persen dengan rupiah.
Rancangan kita tidak menularkannya — `T_GENERAL_POLIS.DEDUCTION1` dan `T_POLIS_XOL.DEDUCTION`
sudah dua kolom terpisah.

⭐ Dan `LAYER*` pada polis proporsional **bernilai `"0"`, bukan kosong**.

### 4q.3 `DEDUCTION` di XOL adalah **uang** `[terverifikasi]`

Buktinya ia **dijumlahkan** antar layer — persentase tidak dijumlahkan:

```
CARI44          = .Deduction + @toDecimal(InputXOL.CARI44)
local.Deduction = local.Deduction + .Deduction
```

Pola identik dengan `GrossPremi` `NetPremi` `DueToValue`, dan ia dikurangkan seperti uang.

### 4q.4 `[keputusan work owner]` Contoh JSON pertama milik **Fac In**

⇒ Pohon lima tingkat *(`LocationList` → `CoverageList` → …)* **tetap di luar lingkup Treaty In**.
Rancangan tidak perlu lima tabel tambahan. Butir terbuka Bab 1.1 **tertutup**.

### 4q.5 `[keputusan work owner]` Tiga penentu bentuk — **biarkan ketiganya**

`QuotationData.ProportionalType` *(170 rujukan)* · `IsNewPolicyNonProp` *(87)* ·
`TreatyType == "XOL"` *(11, dipakai `InsetTreatyInProdAddendum_Act`)*.

Ketiganya satu polis yang sama. Diikuti apa adanya, tidak digabung.
⚠️ Catatan pelaksanaan: Go menulis ketiganya dari **satu titik keputusan**, supaya tidak menyimpang.

### 4q.6 `[keputusan work owner]` `STATEMENT_DATE` = tanggal mulai — **disengaja**

`STATEMENT_DATE` dan `BEGINDATE` sama-sama dari `CARI11` = `.StartDate`, padahal property
`PolicyTreatyIn.StatementDate` ada dan dirujuk 19 kali. Itu memang maksudnya.

### 4q.7 `[keputusan work owner]` P61 — `_BACKUP` ikut terhapus: **disengaja**

`_BACKUP` bukan jaring pengaman. Ditiru apa adanya.
⚠️ Dicatat di spec bahwa **namanya menyesatkan**, supaya pembaca berikutnya tidak mengira ada
jalur pemulihan.

### 4q.8 `[keputusan work owner]` P62 — tanggal mati: **biarkan**

`IF TRUNC(v_now) <= TO_DATE('02/01/2026') THEN v_mm_yyyy := '12.2025'` adalah **modifikasi akhir
tahun** yang dipakai bila tanggal closing bergeser. Perilakunya dipertahankan.

⚠️ **Usulan, menunggu persetujuan:** di sistem baru tidak ada procedure untuk disunting. Kalau
penyesuaian ini ditulis di Go, tiap akhir tahun butuh rilis kode. Sebaiknya nilainya jadi **data**
— satu baris setelan, bukan baris program.

### 4q.9 `[keputusan work owner]` `M_TREATY_IN.JSONDATA` — **di luar lingkup sekarang**

`[terverifikasi]` Dibaca **11 SQL di 8 modul**: NB Treaty In · EDM Treaty In · Treaty In ·
Treaty In Adjustment · NB FacIn · Claim Prop · Claim Non Prop · Master Product Name Life.
`Claim Prop` dan `Claim Non Prop` memakainya untuk **memeriksa batas limit treaty** lewat
`GetLimitsTreatyIn_SQL`.

⭐ Oracle bahkan menanyai isinya langsung: `a.jsondata.Note`, `a.jsondata.ID`.

⛔ **RALAT dua klaim saya:**

> *"Data kontrak treaty SUDAH relasional — yang perlu dipecah hanya data polis."*

Terlalu jauh. Yang relasional hanya **20 kolom ringkasan** `TREATY_IN`. Rincian kontrak masih
dokumen di `M_TREATY_IN.JSONDATA`, dan ada keluarga tabel serupa: `M_TREATY_IN_EDM` ·
`M_TREATY_IN_DETAIL_EDM` · `M_TREATY_OUT` · `M_REINSURANCE…`.

> *"Tidak ada lagi dokumen JSON."*

Tidak benar seluruhnya. Untuk sisi **kontrak**, dokumen itu tetap ada. ⇒ **Go harus tetap mampu
MEMBACA JSON**, walau tidak menulisnya.

### 4q.10 `LAYER*` tingkat polis — **pantulan, bukan data** `[terverifikasi]`

```
InputPolicyTreatyInDetail_preACT:
  POLIS.LayerType = pyReportContentPage.pxResults(1).LAYERTYPE
  POLIS.Layer     = pyReportContentPage.pxResults(1).LAYER
```

Sumbernya **kolom tabel**, dibaca balik saat layar rincian dibuka.

⇒ ⭐ `T_GENERAL_POLIS.LAYER` dkk. **dicoret** — turunan dari `T_POLIS_XOL_LAYER`.

⚠️ `pxResults(1)` hanya mengambil **baris pertama**. Nilai di tingkat polis itu **layer pertama
saja**, bukan ringkasan seluruh layer. Laporan yang memakainya sebagai "layer polis ini"
mengabaikan layer kedua dan seterusnya.

### 4q.11 Induk `TreatyXOLDifferenceList` — **hanya salinan kunci**, tabelnya dibuang

`[terverifikasi]` `CalculateDifferenceEDM_act` menulis **tujuh medan kunci yang sama persis** ke
kedua tingkat — `Layer` `LayerType` `LayerPart` `LayerPartType` `Currency` `IDCurrency` `DueTo`.
Induk tidak membawa informasi yang tidak ada di anak; hanya anak yang menerima 10 kolom uang.

Perbedaan bentuknya dari `TreatyXOLList` semata karena **dibangun aturan berbeda**, bukan karena
maknanya berbeda.

⇒ ⭐ **`T_POLIS_XOL_DIFFERENCE` dibuang.** Cukup `T_POLIS_XOL_LAYER_DIFFERENCE`, menempel langsung
ke `T_POLIS_DIFFERENCE`. Pengelompokan per mata uang tetap bisa lewat `ID_CURRENCY` di tiap baris.

| | Semula | Sekarang |
| --- | ---: | ---: |
| tabel proyeksi, NonProp EDM | 5 | **4** |
| total tabel, EDM NonProp | 15 | **14** |

### 4q.12 `[keputusan work owner]` Tipe kolom — **konversi sebagian**

| Golongan | Tipe | Alasan |
| --- | --- | --- |
| uang | angka presisi tetap, skala ≥ 9 | ADR-0003 · 4q.1 |
| persen | angka presisi tetap | ketetapan 2 P29 |
| tanggal | `DATE` | dua format masuk: `YYYYMMDD` dan cap waktu Pega bersufiks ` GMT` |
| cacah | bilangan bulat | `NOURUT` `PRODKE` `INSTALLMENT_NO` |
| ⭐ **kode** | **TETAP teks** | `GroupPanel = "006"` · `BusinessOldId = "01"` — nol di depan **wajib utuh**, kalau tidak penggolong `BusinessType_DeT` 36 baris gagal |
| ⭐ **penanda** | **TETAP teks** | `IsApproved` dibandingkan sebagai **teks** *(ketetapan 1, P24/P6)*; `""` keadaan sah yang **berbeda** dari `"0"` |

Konversi terjadi **sekali saat masuk**, bukan tiap kali dibaca.

⛔ `[terbuka]` **Satu butir kecil tersisa:** nilai teks kosong `""` yang masuk kolom bertipe angka
atau tanggal menjadi `NULL` atau `0`? Saya memakai **`NULL`** sebagai bawaan — sebab `0` adalah
nilai uang yang bermakna, dan tanggal tidak punya nol. **Perlu dikonfirmasi.**

---

## 4quinque · ⭐ KOREKSI 23 September 2026 — sesudah data guide dan satu dokumen nyata

Work owner menyerahkan **data guide** `POOLDATA.JSON_POLIS.DATA_JSON` dan **satu dokumen EDM
sungguhan**. Keduanya mengubah beberapa hal, dan **membatalkan tiga kesimpulan saya.**

⛔ **Nilai dokumen itu tidak disimpan di berkas mana pun** — ia memuat nama orang pada medan
pemasar, operator dan tertanggung. Yang dicatat hanya **bentuk**, **nama medan**, dan **jumlah**.

### 4q5.1 ⛔⛔ Data guide **TIDAK LENGKAP** — dan saya sempat memperlakukannya sebagai kebenaran

`[terverifikasi]` **Satu dokumen tunggal memuat 95 jalur yang tidak ada di data guide.**

| | Jalur |
| --- | ---: |
| data guide | 378 |
| satu dokumen EDM | 319 |
| ⛔ ada di dokumen, **tidak** di data guide | **95** |

Yang terlewat antara lain: `QuotationData.CedingCoList` dan ketiga anggotanya · seluruh
`BreakDownSpreadList` · `Claim` · `ExcessLoss` · `GrossPremium` · `BalanceBeforePPH` ·
`BalanceBeforeTax` · `ListInstallment.DueDate` · `ListInstallment.InstallmentNo` ·
`ResultOgp1/2` · `ResultOnp1/2` · `RiComm*` · `OveriddingComm*` · `SalvageValue` · `MasterID` ·
`ClaimType` · `ClaimPaymentType` · `IsSOAUpload` · `QuotationList` · `NoOfferSlip` · `EdmDate` ·
`RNWDate` · `btnQuotation`.

⇒ Indeksnya tampaknya **basi**. `[keputusan work owner]` DBA diminta **menyegarkannya**.

⚠️ **Sampai disegarkan, data guide berkedudukan PELENGKAP, bukan sumber kebenaran.** Ia sah untuk
`o:length`; ⛔ **tidak sah** untuk membuktikan ketiadaan sebuah medan.

### 4q5.2 ⛔ Tiga kesimpulan saya yang ditarik

> ⛔ Bunyi lama: *"`T_POLIS_CEDING` gugur — `CedingCoList` tidak ada di dokumen polis Treaty."*

**Keliru.** `[terverifikasi]` `QuotationData.CedingCoList` **ada**, dengan `CedingCo` dan
`CedingCoName` per baris — di tingkat polis **dan** di dalam `OldData`.
⇒ ⭐ **`T_POLIS_CEDING` TETAP DIBUAT.** `[keputusan work owner]`

> ⛔ Bunyi lama: *"Tabel inti 48 kolom, bukan 79."*

**Tidak terbukti.** Angka 48 berasal dari data guide yang tidak lengkap. Sensus korpus **79** justru
lebih dekat, dan tetap **batas bawah**.

> ⛔ Bunyi lama: *"`GroupPanel` `BusinessOldId` `BusinessCode` tidak ada di dokumen."*

**Tidak terbukti.** Ketiadaan di data guide **bukan bukti ketiadaan**. ⛔ Tetap `[terbuka]`.

### 4q5.3 ⭐ Tabel baru — ~~`T_POLIS_BREAKDOWN_SPREAD`~~ ⛔

`[terverifikasi]` `BreakDownSpreadList` nyata, dua baris pada dokumen contoh, **empat medan**:

```
TreatyType  ·  SharePercentage  ·  PremiumSpreaded  ·  ClaimSpreaded
```

Ada di tingkat polis **dan** di dalam `OldData`. Sensus korpus melewatkannya karena tidak satu pun
aturan merujuk anggotanya — persis jebakan "medan yang hanya lewat".

⚠️ Bentuknya mirip `T_POLIS_SPREADING` tetapi **tidak sama**: spreading punya `Currency` dan
`CurrencyID`, breakdown tidak.

⛔ `[terbuka]` **Apa bedanya secara dagang** — pemiliknya `[work owner]`.

### 4q5.4 ⭐⭐ Rumus selisih — naik dari `[dugaan]` ke `[terverifikasi]`

Dokumen contoh membuktikannya dari **data**, bukan dari membaca aturan:

```
OldData.NetPremium              130.463.146,76
NetPremium (generasi baru)                0,00
TreatyDifference.NetPremium    −130.463.146,76     = baru − lama     ✔
```

Idem `PremiOgp` −200.712.533,47 · `ResultOgp1` −70.249.386,71 · `Claim` −10.502.027,00 ·
`ListInstallment.Premium` −119.961.119,76 · `SpreadingRiskList.PremiumSpreaded` −130.463.146,76.

⇒ ⭐ **`baru − lama.nilai`**, bukan `baru − lama.selisih`. Keputusan work owner 23-09-2026
**terbukti benar**, dan sekarang ada buktinya dari data produksi.

⭐ Dan penjelasan varian ganjil itu ikut ketemu: `OldData.TreatyDifference` isinya
**kosong** — `{pxObjClass, ListInstallment:[], SpreadingRiskList:[]}`. Jadi varian kedua memang
mengurangi terhadap **nol**, bukan terhadap angka yang bermakna.

⭐ Persentase **disalin, tidak dikurangi** — juga terbukti:
`TotalSharePercentagePremium` bernilai `"100"` di ketiga tempat *(baru, lama, selisih)*.

### 4q5.5 `[keputusan work owner]` Presisi uang — **`NUMBER(20,8)`**

⚠️ Ini **menyempurnakan P29**, dan sebagian membalikkannya. Bunyi lama dikutip:

> ⛔ *"P29 — ikuti apa adanya, jangan bulatkan. Skala kolom uang minimal 9 desimal."*

Yang berlaku sekarang: **`NUMBER(20,8)`** — 12 digit di depan koma, 8 di belakang.

```
tersimpan sekarang   592.629.512,880000276      9 desimal
NUMBER(20,8)         592.629.512,88000028       dibulatkan ke 8
ShareValue           …,000000000000000000000000  24 desimal  →  8
PremiumSpreaded      …,00000000000000000000      20 desimal  →  8
```

⚠️ ⭐ **Bukan pembatalan penuh.** Delapan desimal jauh lebih halus daripada dua, dan ekor galat
`2,76 × 10⁻⁷` **tetap terlihat** — ia jatuh di desimal ketujuh. Yang berubah: nilai berdesimal
lebih dari delapan **dibulatkan saat dimuat**, jadi tidak lagi "apa adanya" secara harfiah.

⇒ Yang ikut berubah:

| | Semula | Sekarang |
| --- | --- | --- |
| skala kolom uang | minimal **9** | **8** |
| AC "tersimpan tanpa kehilangan satu digit pun" | `592629512.880000276` utuh | utuh **sampai 8 desimal** |
| migrasi | tidak mengubah satu pun nilai | **membulatkan** yang berdesimal > 8 |

⛔ `[terbuka]` **12 digit di depan koma cukup atau tidak** belum diuji terhadap nilai terbesar
yang pernah tersimpan. Nilai terbesar yang terlihat 9 digit; premi rupiah bisa lebih panjang.

### 4q5.6 Panjang kolom — dari `o:length`

`[terverifikasi]` Data guide sahih untuk ini, dan menghapus kebutuhan menebak:

| Medan | Panjang | | Medan | Panjang |
| --- | ---: | --- | --- | ---: |
| `PolicyNo` · `EDMNo` | 32 | | `Suggest` | **256** |
| `CedingCoName` · `InsuredName` | 32 | | `Remark` | **128** |
| `SOBName` · `SobName` | 64 | | `TypeTax` · `NoOffer` | 16 |
| `IsApproved` · `ProdKe` · `Layer` | **1** | | `Currency` · `TreatyType` | 4 |

⚠️ Tetap **batas bawah** — data guide tidak melihat seluruh baris.

### 4q5.7 Tiga medan yang sensus korpus lewatkan

`Remark` *(128)* · `Show` · `ViewState`.

`Show` dan `ViewState` **tidak dimigrasi** — keadaan layar. `Remark` **dimigrasi**, dan ini medan
baru yang belum pernah masuk rancangan.

---

## 5 · Tiga hal yang **tidak** ditiru

| Yang tidak ditiru | Sebab |
| --- | --- |
| ⛔ penjaga duplikat `ACHIEVEMENT` berbasis **17 kolom termasuk uang** | dipatahkan galat `2,76 x 10^-7`; kunci diganti **medan pengenal** |
| ⛔ `OldData` sebagai salinan tersimpan | digantikan generasi `PRODKE − 1` |
| ⛔ `TreatyDifference` sebagai nilai tersimpan | dihitung; menyimpannya menumpuk galat pada endorsemen berlapis |

---

## 6 · Enam butir yang **harus diputuskan sebelum DDL ditulis**

⛔ **Tidak satu pun saya tutup.** Urut menurut yang paling menahan.

| # | Butir | Saran saya | Pemilik |
| ---: | --- | --- | --- |
| 1 | **Presisi uang** — ikuti apa adanya · bulatkan saat migrasi · bulatkan saat dihitung | **bulatkan saat dihitung**, simpan apa adanya; data lama tidak diubah dan galat tidak menyebar | `[work owner]` + `[Finance]` |
| 2 | **`JSON_DATAGUIDE(DATA_JSON)` dari DBA** | ⭐ minta — hanya ini yang dapat menutup selisih 74 lawan 86 medan | `[data DBA]` |
| 3 | **Contoh JSON pertama milik Treaty atau Fac?** | pastikan sebelum pemecah diuji | `[work owner]` |
| 4 | **Selisih disimpan atau dihitung?** | **dihitung** | `[work owner]` |
| 5 | **`CEDING_CO` berakhiran `"; "`** — dipecah jadi baris atau tetap teks? | **dipecah**, sebab satu polis dapat punya lebih dari satu ceding | `[work owner]` |
| 6 | **Migrasi dokumen lama** — dipindahkan atau dibaca lewat jalur lama? | **dipindahkan**, dengan `DATA_JSON` dipertahankan sebagai arsip baca-saja sampai rekonsiliasi lulus | `[work owner]` |

⛔ Butir 1 **menahan** butir 4, dan butir 2 **menahan** kelengkapan Bab 4.

---

## 7 · Cara menguji ulang rancangan ini

⚠️ Pakai Python — jalur korpus memuat spasi. ⚠️ Buang `pyExpressionGadget` lebih dulu.
⛔ **Jangan `html.unescape` sebelum mencocokkan pola struktur** — entitas `&lt;` memecah polanya.

```python
import os, re, html, collections
ROOTS = {"NB": r"D:\XML\RNM_BRD\NB Treaty In", "EDM": r"D:\XML\RNM_BRD\EDM Treaty In"}
GADGET = re.compile(r"<pyExpressionGadget>.*?</pyExpressionGadget>", re.S)
REF = re.compile(r"(?:pyWorkPage\.|MergePage\.|TempPage\.|\.)?"
                 r"PolicyTreatyIn((?:\.[A-Za-z0-9_]+(?:\([^)]*\))?)+)")
paths = collections.Counter()
for tag, root in ROOTS.items():
    for dp, dn, fn in os.walk(root):
        for f in fn:
            if not f.lower().endswith(".xml"): continue
            raw = GADGET.sub("", open(os.path.join(dp, f), encoding="utf-8",
                                      errors="replace").read())
            for m in REF.finditer(html.unescape(raw)):
                paths[re.sub(r"\([^)]*\)", "()", m.group(1))] += 1
print("jalur unik:", len(paths))                      # 322
top = [p[1:] for p in paths if p.count(".") == 1 and "(" not in p]
cont = {n for n in top if any(p.startswith("." + n + ".") or p.startswith("." + n + "(")
                             for p in paths)}
print("tingkat atas:", len(top), "wadah:", len(cont), "skalar:", len(top) - len(cont))
#                                          94              8              86
print("sarang OldData dalam OldData:",
      sum(1 for p in paths if p.startswith(".OldData.") and "OldData" in p[9:]))   # 0
```

| Pernyataan | Angka | Cara A | Cara B |
| --- | ---: | --- | --- |
| jalur properti unik | **322** | pola `REF` di atas | — |
| medan tingkat atas | **94** | ruas tunggal tanpa `()` | 8 wadah + 86 daun |
| `OldData` bersarang di `OldData` | ⭐ **0** | sapuan jalur | sejalan keadaan EDM Bab 9 |
| `LocationList` di NB/EDM Treaty In | ⭐ **0** | sapuan 9.430 berkas korpus | terhitung per modul |
| `ListInstallment` bersarang / datar | **58 / 101** | 6 berkas / 11 berkas | — |
| `TreatyDifference` milik EDM | **199 dari 201** | cacah per modul | — |

---

## 8 · TELEMETRI EKSEKUSI

| | |
| --- | ---: |
| Berkas korpus disapu | **9.430** `.xml` |
| Berkas proyek dibaca | 4 |
| Skrip sapuan dijalankan | 6 |
| Panggilan alat | 11 |
| ⚠️ Taksiran token masuk | ± 95–110 rb |
| ⚠️ Taksiran token keluar | ± 16–19 rb |

⚠️ **Taksiran, bukan ukuran.** Angka token sejati tidak terlihat dari dalam sesi. Yang terukur hanya
bila dijalankan lewat `claude --print --output-format json "<prompt>" > hasil.json`. ⛔ Pengukuran
dari luar **tidak dilakukan** untuk ronde ini.

---

*Berlaku untuk `NB Treaty In` dan `EDM Treaty In`. Rancangan ini menggantikan perkiraan pohon lima
tingkat pada `PROMPT-RANCANGAN-TABEL-NB-TREATY-IN.md` yang ditarik 22 September 2026 — bunyi lamanya
dikutip di Bab 1.1, tidak dihapus.*
