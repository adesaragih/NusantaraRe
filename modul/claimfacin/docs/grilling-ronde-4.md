# Grilling ronde 4 — Claim Fac In

**Tanggal:** 2026-09-19 · **Korpus:** `D:\XML\RNM_BRD\Claim Fac In\` — **482 berkas `.xml`**, READ-ONLY.
**Ronde pendek — tiga berkas rule dan satu penelusuran alur.** ⛔ Keempat berkas grilling lama tidak
disentuh satu byte pun.

> ⚠️ **Ronde ini MERALAT ronde 3 di dua tempat**, dan keduanya bukan soal kecil: satu tentang
> **tipe angka uang**, satu tentang **bentuk selesai alur**. Keduanya ditemukan karena diminta
> **menghitung dengan cara kedua** *(aturan I3/P4)*. Rinciannya di **§F3**.

---

## LANGKAH 0

```
0a  isi .scratch\claim-facin\   4 berkas   ✅ grilling-ronde-4.md belum ada
0b  panjang keempatnya          712 · 735 · 1.145 · 470   ✅ cocok
0c  berkas .xml korpus          482        ✅
0d  docs\adr\*.md               15         ✅
```

---

## §A — Sembilan butir `[terbuka]` yang dibawa

⛔ **Nol ditutup.** Tiga tersentuh ronde ini, enam tidak.

| # | Butir | Ronde ini |
| --- | --- | --- |
| 1 | Kenapa `CreateKMTNo_Act` langkah 13.1 ditimpa 13.1.1 | tidak disentuh |
| 2 | `FlagOnGoingCommitte` Fac In kolom yang sama dengan Claim Prop? | tidak disentuh |
| ⭐ 3 | `TempOutstanding` disimpan rule **lain**? | ⭐ **disentuh — §E2** |
| ⭐ 4 | Arti `REQUIRED = -1` versus `0` | ⭐ **disentuh — §B4**, tetap terbuka |
| ⭐ 5 | Isi tiap penampung `CARI<n>` | ⭐ **disentuh — §B2**, tetap terbuka |
| 6 | `IndexObject` versus `.ObjectIndex` | tidak disentuh |
| 7 | Arasapas tanpa saringan — dikehendaki? | tidak disentuh |
| 8 | Pemetaan *"aturan A–J"* ke sembilan aturan bernomor | tidak disentuh |
| 9 | Di mana catatan aturan baca disimpan | tidak disentuh |

---

## §B — `SetKomiteList_ACT` dan `TotalAdjustment`

### B1 — Bentuk rule-nya

`[terverifikasi]` **Identitas empat bagian:**

```
pxInsName  ASM-FW-GCNMFW-DATA-ADJUSTMENT!SETKOMITELIST_ACT
kelas      ASM-FW-GCNMFW-Data-Adjustment
ruleset    GCNMFW        versi 01-01-24
```

**8 langkah · 2 ber-remark · gerbang hidup 2 · gerbang mati 0.**

| Langkah | Keadaan | Isinya |
| --- | --- | --- |
| 1 | hidup | `Property-Remove` |
| 2 | **hidup**, bergerbang `pyWorkPage.ClaimData.PaymentType…` | induk blok |
| 2.1 | hidup | `Local.Currency := .CurrencyDol` |
| ⛔ **2.2** | **REMARK** | `.ValueAdjustment := Local.Currency * .AdjustmentValue` |
| **2.3** | **hidup** | `.ValueAdjustment := .AdjustmentValue * Local.Currency` |
| **2.4** | **hidup**, bergerbang `.IsKomite==""` | **`Local.TotalAdjustment := .ValueAdjustment`** |
| 3 | hidup | `Call SetListKomite_act` |
| ⛔ **4** | **REMARK** | **`Obj-Save`** |

⭐ **Dua hal yang menonjol:**

1. **Langkah 2.2 dan 2.3 adalah perkalian yang SAMA dengan operan tertukar** — `kurs × nilai`
   versus `nilai × kurs`. Yang pertama **di-remark**, yang kedua hidup. ⛔ Alasannya tidak tertulis.
   ⚠️ Hasilnya sama secara aritmetika; **yang berbeda hanyalah urutan penulisannya**. `[terbuka]`
2. ⭐ **`Obj-Save` di langkah 4 DI-REMARK** — jadi rule ini **tidak menyimpan apa pun**.

⚠️ **Jebakan nama:** modul ini punya **`SetKomiteList_ACT`** *(rule ini)* **dan**
**`SetListKomite_act`** *(yang dipanggilnya di langkah 3)* — **dua rule berbeda dengan kata yang
tertukar**. ⛔ Siapa pun yang menyamakannya akan salah.

### B2 — Ke mana `TotalAdjustment` pergi, dan dari mana datang

`[terverifikasi]` **Pemanggil `SetKomiteList_ACT`: dua** — `Activity/SetListKomite_act.xml` dan
`Activity/SetNilaiResikoSendiri.xml`.

⚠️ **Yang ditugaskan di langkah 2.4 adalah `Local.TotalAdjustment`** — properti pada halaman
**`Local`**, ⛔ **bukan parameter `TotalAdjustment`**. Keduanya bernama sama dan **mudah tertukar**.

`[terverifikasi]` Nama `TotalAdjustment` muncul di **delapan berkas**: `CreateKMTNo_Act` ·
`SendPICProtect_Act` · `SetKomiteList_ACT` · `SetNilaiResikoSendiri` · `SetSalvageValue` ·
`SetValueAdjusterFee` · `Section/InputAdjustment` · `Section/ItemListEstimation`.

⛔ **Apakah nilai itu akhirnya DISIMPAN tidak terbaca dari rule ini** — `Obj-Save`-nya di-remark, dan
penelusuran ke delapan berkas di atas **tidak dikerjakan ronde ini**. `[terbuka]`

### B3 — Pembulatan di jalur itu

`[terverifikasi]` **NOL.** Di ke-8 langkah `SetKomiteList_ACT` tidak ada pembulatan, pemotongan,
pemformatan angka, maupun penetapan skala. Yang ada hanya **satu perkalian** *(2.3)* dan **satu
penyalinan** *(2.4)*.

### B4 — ⭐⭐ Dan inilah yang meralat ronde 3: **tipe parameter punya DUA sumber yang BERSELISIH**

`[terverifikasi]` Ronde 3 §C membaca tipe dari **`pyXMLSignature`** saja. Ronde ini membacanya
**dua cara** *(aturan I3/P4)*:

```
sumber A  pyXMLSignature            — atribut TYPE pada tiap parameter        169 parameter
sumber B  pyParametersParamType     — medan tipe di samping pyParametersParamName
muncul di KEDUANYA                                                            169
⭐ TIPENYA BERSELISIH                                                            8
```

| Activity | Parameter | Sumber A *(tanda tangan)* | Sumber B *(medan)* |
| --- | --- | --- | --- |
| ⭐⭐ **`SetKomiteList_ACT`** | **`TotalAdjustment`** | **`Double`** | ⭐ **`Decimal`** |
| ⭐ `GetCoverageAneka_Act` | `TSI` | `Decimal` | ⭐ **`String`** |
| ⭐ `SetCedant_act` | `Share` | `STRING` | ⭐ **`Decimal`** |
| `CheckDate_Act` | `DOL` | `Date` | `String` |
| `GeneratePLA` | `idxobj` | `INTEGER` | `int` |
| `SetIndex_Act` | `Name` | `STRING` | `String` |
| `SetTreatyname` | `IndexObject` | `STRING` | `int` |
| `SetTreatyname` | `IndexObjectItem` | `STRING` | `String` |

⛔ **Tidak ada dasar untuk memilih salah satunya.** Delapan parameter **belum punya tipe yang
pasti** *(aturan I3/P4)*. ⚠️ **Tiga di antaranya menyentuh uang** — `TotalAdjustment` · `TSI` ·
`Share` — dan ketiganya persis yang dipakai ronde 3 untuk memvonis ADR-0003.

⚠️ **Sumber B jauh lebih besar dan TIDAK bersih:** `pyParametersParamName` muncul **4.071** kali di
179 berkas `Activity`, sedangkan parameter activity hanya **169**. ⛔ Jadi sumber B **dipakai juga
untuk hal lain**, dan hanya ke-169 yang beririsan dapat dibandingkan. **Sisanya tidak diperiksa.**

⚠️ **Arti `REQUIRED = -1` versus `0` tetap `[terbuka]`** — ⛔ tidak dipakai sebagai dasar apa pun.

### B5 — ⭐ Vonis **ADR-0003**: **BERTAHAN — tetapi dasarnya BERGESER**

**Yang GUGUR:** ⛔ *"`TotalAdjustment` bertipe `Double`"* **tidak lagi dapat dipakai sendirian**,
karena medan parameter rule yang sama menyebutnya **`Decimal`**.

**Yang BERTAHAN, dan kini menjadi dasar utamanya** `[terverifikasi]`:

1. ⭐ **Nilai uang lewat sebagai teks.** Menurut tanda tangan, `AdjustmentValue` · `Share` ·
   `Estimation` bertipe **`STRING`** — ketelitiannya tidak dijaga tipe apa pun.
2. ⭐⭐ **Dan itu terbukti menimbulkan tambalan di SQL**: `CariHistoryClaim_SQL` menjalankan
   `REPLACE(a.DATA_JSON.ClaimEstimate , ',' , '.') AS "TSI"` — **koma diganti titik pada nilai
   uang**, persis gejala uang-sebagai-teks *(ronde 3 §B5)*.
3. ⚠️ **Tipe uang tidak punya satu sumber kebenaran** — dua deklarasi di rule yang sama berselisih
   pada **tiga** parameter uang.

⛔ **Dilarang mengusulkan revisi ADR.** ⛔ **Nol kesimpulan untuk aplikasi Go.**

---

## §C — Dua rantai yang berperilaku **ATAU**

### C1 — `SpreadingCheckSP` langkah 3.2 — ⭐ ATAU-nya memang disengaja

`[terverifikasi]` `ASM-FW-GCNMFW-DATA-OBJECT!SPREADINGCHECKSP`, ruleset `GCNMFW 01-01-01`,
**12 langkah · nol remark**. Langkah **3.2**, flag **`true`**, keterangan pengembang:
***"Detecting empty fields"***.

| Baris | Syarat | Benar | Salah |
| --- | --- | --- | --- |
| 1 | `.TreatyType == ""` | ⭐ **LEWATI-WHEN** *(kode 5)* | lanjut |
| 2 | `.SharePercentage == ""` | ⭐ **LEWATI-WHEN** *(kode 5)* | lewati-langkah |

**Dibaca sebagai ATAU** — *"jalankan bila **salah satu** medan kosong"* — dan itu **cocok dengan
keterangan pengembangnya**.
**Dibaca sebagai DAN** — *"jalankan hanya bila **keduanya** kosong"* — maka klaim yang hanya kosong
di **satu** medan **lolos tanpa pesan**.

⭐ **Bacaan ATAU-lah yang benar**, dan ia **menyentuh uang**: `.SharePercentage` adalah persen share.
Tiga pesan galat rule ini — langkah **4** *(% > 100)*, **6** *(% < 100)*, **8** *(medan kosong)* —
seluruhnya `Property-Set-Messages` ber-flag `true`.

### C2 — `CheckEstimateValue` langkah 24 dan 25 — ⚠️ ternyata **soal TOMBOL, bukan uang**

`[terverifikasi]` `ASM-FW-GCNMFW-DATA-OBJECTITEM!CHECKESTIMATEVALUE`, ruleset `GCNMFW 01-01-24`,
**60 langkah · 1 ber-remark**.

| Langkah | Keterangan pengembang | Baris kode 5 |
| --- | --- | --- |
| **24** | *"Disable download cfs button when Local.StsErr set"* | baris 1 `@hasMessages(pyWorkPage)` → **LEWATI-WHEN** |
| **25** | *"Set button download CFS enable"* | baris 2 `Local.estimasimotethanTSI` → **LEWATI-WHEN** |

⭐ **Keduanya menggerbangi TOMBOL UNDUH CFS, bukan perhitungan uang.**

⚠️ **RALAT terhadap ronde 3 §G2**, yang menulis *"keduanya menyentuh uang"*. Yang benar:
**`SpreadingCheckSP` menyentuh uang; `CheckEstimateValue` pada kedua titik itu menyentuh tombol.**
⛔ Rule-nya **secara keseluruhan** memang rule validasi estimasi — tetapi **bukan di titik kode 5
itu**.

### C3 — Yang memblokir versus yang memperingatkan

`[terverifikasi]` `CheckEstimateValue` punya **7 `Property-Set-Messages`**, seluruhnya ber-flag
`true`, ditambah satu `Page-Clear-Messages` di langkah 1:

```
lgk 7 · 8      "Nilai Estimasi ke-1 tidak …"
lgk 12.2 · 12.5  "Estimasi Date tidak boleh …"
lgk 12.6       "Nilai Estimasi tidak bole…"
lgk 21 · 22    "Nilai Total Estimasi tidak…"
```

⛔ **Mana yang MEMBLOKIR dan mana yang hanya PERINGATAN tidak terbaca dari struktur langkah.**
Ketujuhnya memakai metode yang sama; pembedanya ada pada **bagaimana pesan itu dibaca layar**, dan
itu **tidak diperiksa ronde ini**. `[terbuka]`

---

## §D — ⛔ Bentuk "selesai" alur — **RONDE 3 SALAH, dan ini ralatnya**

### D1 — Dua cara hitung, dan angkanya **BEDA JAUH**

`[terverifikasi]` `Flow/Register_Flow.xml` — dan modul ini **memang hanya punya SATU berkas Flow**
*(dihitung ulang: 1)*.

| Cara | Yang dihitung | Hasil |
| --- | --- | --- |
| **1** — `pyShapeType` *(dipakai ronde 3)* | medan `pyShapeType` yang terisi | **16** |
| ⭐ **2** — `pxObjClass` berawalan `Data-MO-` | jenis sebenarnya tiap baris | ⭐ **27** |

⭐⭐ **Selisih 11.** Sebabnya: **sebagian bentuk tidak punya medan `pyShapeType` sama sekali** —
jenisnya hanya tertulis di `pxObjClass`. ⛔ Cara 1 **buta terhadapnya**.

### D2 — ⭐ Sensus yang benar

| Jenis | Jumlah |
| --- | --- |
| `Data-MO-Connector-Transition` | **10** |
| ⭐ **`Data-MO-Event-Exception`** | ⭐ **9** |
| `Data-MO-Gateway-Decision` | **3** |
| `Data-MO-Activity-Assignment` | **3** |
| ⭐ **`Data-MO-Event-End`** | ⭐ **1** |
| `Data-MO-Event-Start` | **1** |

### D3 — ⭐⭐ Bentuk selesai **ADA**, dan statusnya sama dengan Komite Claim Prop

`[terverifikasi]` Barisnya ber-`REPEATINGINDEX="END52"`:

```
pxObjClass     Data-MO-Event-End
pyMOId         END52
pyWorkStatus   Resolved-Completed
pyMOName       (kosong)
pyCategory     FlowStandard
```

⭐ **Status akhirnya `Resolved-Completed`** — **sama persis** dengan `END52` di modul Komite Claim
Prop. ⚠️ **Bahkan nomor bentuknya sama: `END52`.** ⛔ Itu **tidak saya jadikan kesimpulan** bahwa
kedua alur berasal dari satu salinan — hanya dicatat sebagai kecocokan yang mencolok. `[terbuka]`

⚠️ **`pyMOName` kosong** — itu sebabnya bentuk ini juga tidak muncul di daftar nama bentuk ronde 3.

### D4 — ⭐ Sembilan bentuk **Event-Exception** yang ronde 3 tidak lihat sama sekali

`[terverifikasi]` **Sembilan** `Data-MO-Event-Exception` di alur yang sama. ⛔ **Apa yang
dipicunya dan ke mana ia mengarah tidak dibaca ronde ini** — ia muncul justru karena sensus cara
kedua. `[terbuka]`

⚠️ **Ini mengubah gambaran daur hidup kasus Fac In:** bukan alur lurus 16 bentuk tanpa akhir,
melainkan **27 bentuk, berakhir `Resolved-Completed`, dengan sembilan jalur perkecualian**.

---

## §E — Tiga tambahan murah

### E1 — `pyPrivilegeClass` terisi 33, dan **nama haknya tetap nol**

`[terverifikasi]` **13 kelas berbeda** di 33 titik:

| Kelas | Kali |
| --- | --- |
| `ASM-FW-GCNMFW-Data-Object` | 10 |
| `ASM-FW-GCNMFW-Work-PNC` | 7 |
| `ASM-FW-GCNMFW-Data-ObjectItem` | 5 |
| `ASM-FW-GCNMFW-Data-ClaimData` | 2 |
| sembilan kelas lain, masing-masing 1 — termasuk `…Data-Quotation` · `…Data-SpreadingRisk` · `…Int-TREATYGROUP` | 9 |

⭐ **Kenapa hanya kelasnya yang terisi?** Karena `pyPrivilegeClass` adalah **kelas tempat hak itu
DICARI**, sedangkan **nama haknya** ada di `pyPrivilegeName` — yang `[terverifikasi]` **terisi NOL
dari 33**. ⚠️ Jadi setiap titik itu berbunyi *"cari hak di kelas X"* **tanpa menyebut hak apa**.

⛔ **Pola yang sama persis dengan Komite Claim Prop**, dan ronde 3 §D3 sudah menandainya. Ronde ini
menambahkan **daftar kelasnya**, ⛔ bukan jawabannya.

### E2 — ⭐ `TempOutstanding`: **NOL rule lain yang menyimpannya**

`[terverifikasi]` Penyisiran **seluruh 482 berkas `.xml`**, tidak peka huruf besar-kecil — nama itu
muncul di **tiga berkas saja**:

| Berkas | Kali |
| --- | --- |
| `Harness/Outstanding.xml` | 8 |
| `Section/Outstanding_SC.xml` | 8 |
| `Activity/GetAllData_Act.xml` | 4 |

⭐ **Ketiganya jalur tampilan** — satu harness, satu section, dan activity yang mengisinya. **Nol
`Obj-Save`, nol rule penyimpan.**

⛔ **Butir terbuka 3 TIDAK saya tutup** — ini bukti dari korpus, dan yang menutup butir adalah work
owner. ⚠️ Tetapi ia **menguatkan** keputusan *"`GetAllData_Act` ikuti apa adanya"* ronde 2: cacat
lompatan itu memang **tidak pernah menyentuh data tersimpan**.

### E3 — ⭐ Rule ber-`COMMIT` sendiri: **TIGA BELAS**

`[terverifikasi]` **5 `Activity` + 8 `RDBList`**:

```
Activity   DLAFacintoTreaty_Act (2) · SetProtectionEstimation (2) · SetSalvageValue (2)
           GetPayAttachmentAdj_Act (1) · InsertLogServiceClaim (1)
RDBList    GetSequenceNumber_SQL · GetTokenStorage_SQL · InsertClaimPNC
           InsertLOGDirectKasir_SQL · InsertProgressClaim_SQL · InsertSUBProgressClaim_SQL
           SaveOSClaim_SQL · UpdateDCauseOfLoss
```

⚠️ **Termasuk `GetSequenceNumber_SQL`** — pembangkit nomor urut. ⛔ Bandingkan Claim Prop, yang
mencatat *"15 dari 58 rule `COMMIT` sendiri"* dan **tidak mereplikasinya**. ⛔ **Saya tidak
menyimpulkan apa pun untuk Fac In** — angkanya dicatat, perlakuannya belum diputuskan. `[terbuka]`

---

## §F — Penutup

### F1 — Yang paling menahan untuk ronde berikutnya

| # | Yang dikerjakan | Kenapa menahan |
| --- | --- | --- |
| **1** | ⭐⭐ **Baca sembilan `Data-MO-Event-Exception`** di `Register_Flow` | baru ketahuan ronde ini; **mengubah gambaran daur hidup kasus** |
| **2** | ⭐⭐ **Adu Claim Prop dan Komite Claim Prop dengan 15 ADR** | utang lintas modul, sudah tiga ronde tertunda |
| **3** | ⭐ **Tentukan sumber tipe parameter mana yang sah** — tanda tangan atau medan | **8 parameter tanpa tipe pasti**, tiga di antaranya uang |
| **4** | **Telusuri ke mana `Local.TotalAdjustment` akhirnya disimpan** — delapan berkas | `Obj-Save` rule-nya di-remark; muaranya belum terbaca |
| **5** | **Periksa mana dari 7 pesan `CheckEstimateValue` yang MEMBLOKIR** | tidak terbaca dari struktur langkah |
| **6** | **Baca ulang kesimpulan kurs standar Claim Prop** | rule yang dipakai di sana 4 tahun 3 bulan lebih tua |
| **7** | **Sesuaikan kalimat Claim Prop *"52 dari 61 rule `When`"*** | keputusan *"rule dipisahkan dari kehidupannya"* |
| **8** | **Putuskan perlakuan 13 rule ber-`COMMIT`** | Claim Prop tidak mereplikasinya; Fac In belum diputuskan |

### F1b — Daftar penyimpangan sadar — ⛔ **tidak berubah: 5 terurai + 6 dari ADR**

⛔ Ronde ini **membaca**, tidak memutuskan. ⚠️ Tetapi **§B5 menggeser dasar** salah satunya: vonis
ADR-0003 **bertahan** dengan bukti yang berbeda.

### F2 — Kesimpulan ronde 4 yang **PALING RAWAN salah**

⚠️ **§B5 — vonis ADR-0003 "BERTAHAN".**

**Kenapa rawan:** ia kini bersandar pada *"nilai uang lewat sebagai `STRING`"*, dan tipe `STRING`
itu berasal dari **sumber A** — sumber yang **baru saja terbukti berselisih dengan sumber B pada 8
parameter**. ⛔ **Kalau ternyata sumber B yang sah**, maka `Share` bertipe `Decimal` dan sebagian
dasar vonis ikut bergeser lagi.

Yang **tidak** rawan dan berdiri kokoh: ⭐ **`REPLACE(… , ',' , '.')` pada `ClaimEstimate`** — itu
**kalimat SQL yang terbaca langsung**, bukan deklarasi tipe. Uang yang tidak perlu ditambal
koma-ke-titik **tidak akan ditambal**.

**Kedua paling rawan:** §D4 — sembilan `Event-Exception` **hanya dihitung, tidak dibaca**.

### F3 — ⛔ RALAT terhadap ronde sebelumnya — **empat**

| # | Yang diralat | Kalimat/angka LAMA, dikutip | Yang benar |
| --- | --- | --- | --- |
| **1** | Bentuk alur | ronde 3 §F1: *"`Flow/Register_Flow.xml` punya **16 bentuk** … dan ⭐ **NOL bentuk selesai**"* | ⭐⭐ **27 bentuk**, dan **bentuk selesai ADA** — `END52`, `Resolved-Completed`. Cara 1 memakai `pyShapeType`; bentuk akhir mendeklarasikan jenisnya di `pxObjClass` |
| **2** | Jalur perkecualian | ronde 3 §F1 tidak menyebutnya sama sekali | ⭐ **9 `Data-MO-Event-Exception`** |
| **3** | Tipe `TotalAdjustment` | ronde 3 §C2: *"⭐⭐ **`SetKomiteList_ACT.TotalAdjustment`** bertipe **`Double`**"* | ⚠️ **dua sumber berselisih** — tanda tangan `Double`, medan parameter `Decimal`. **Belum punya tipe pasti** |
| **4** | Kode 5 menyentuh uang | ronde 3 §G2: *"termasuk ⭐ **`SpreadingCheckSP`, yang menyentuh spreading**, dan ⭐ **`CheckEstimateValue`, rule validasi 60 langkah**"* | ⚠️ **`SpreadingCheckSP` benar**; **`CheckEstimateValue` pada kedua titik kode 5 itu menggerbangi TOMBOL UNDUH CFS**, bukan uang |

### F4 — Yang seharusnya dikerjakan tetapi **TIDAK diperintahkan** blok ini

| # | Butir |
| --- | --- |
| **1** | **Membaca isi sembilan `Event-Exception`** — ronde ini hanya menghitungnya |
| **2** | **Menyisir 4.071 `pyParametersParamName`** — hanya 169 yang beririsan yang dibandingkan |
| **3** | **Membaca `SetListKomite_act`** — rule bernama nyaris sama yang dipanggil langkah 3 |
| **4** | **Menelusuri kenapa langkah 2.2 di-remark** dan 2.3 hidup, padahal perkaliannya sama |
| **5** | **Memeriksa kecocokan `END52`** antara Fac In dan Komite Claim Prop |

### F5 — ⭐ **Vonis ADR-0003 sesudah §B: BERTAHAN — MENENTANG, tetapi dasarnya bergeser dari tipe `Double` ke uang-sebagai-teks yang terbukti ditambal koma-ke-titik di SQL.**

---

## Lampiran — bukti berkas lain tidak disentuh

```
korpus Claim Fac In        482 berkas .xml
berkas lama claim-facin      4 berkas  (712 + 735 + 1.145 + 470 baris)
docs\adr\                   15 berkas
modul lain di .scratch\    175 berkas  (claim-prop dan komite-claim-prop TERMASUK)
```

⛔ **Nol berkas korpus dibuka untuk ditulis.** ⛔ **Nol nilai rahasia disalin.**
⛔ **Nol butir `[terbuka]` ditutup.** ⛔ **Nol keputusan §A5 ronde 2 diubah.**
