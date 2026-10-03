# Grilling Claim Fac In — **Ronde 5**

> **Ronde MEMBACA.** ⛔ Nol keputusan work owner dibuat. ⛔ Nol spec, nol tiket. Tiap butir ditutup
> dengan **pertanyaan + rekomendasi**, bukan dengan jawaban.

**Sensus ronde ini:** butir dikerjakan **5** · ⛔ **RALAT terhadap ronde 4: LIMA** · butir
`[terbuka]` ditutup **3** · `[terbuka]` baru **2** · keputusan work owner **2** *(§K)*

> ⛔ **RALAT 2026-09-19** — baris sensus lamanya **dikutip, tidak dihapus**: *"butir `[terbuka]`
> ditutup **3** · `[terbuka]` baru **1** · **pertanyaan untuk work owner 2**"*. Kedua pertanyaan
> §H **sudah dijawab** pada hari yang sama dan menjadi **§K**; jawaban **K2** melahirkan **satu
> butir `[terbuka]` baru**, sehingga `[terbuka]` baru menjadi **dua**.
>
> **Dihitung dua cara, dan keduanya sepakat:** *(a)* mencacah tanda `[terbuka]` pada §G2 dan §K2
> — **2**; *(b)* delta terhadap baris lama — `1 + 1 = 2`.

**Jendela:** seluruh **482** berkas `.xml` `Claim Fac In`; pembanding **329** `Claim Prop` dan
**80** `Komite Claim Prop`. ⛔ Nol berkas korpus dibuka untuk ditulis. ⛔ Nol kredensial disalin.

---

## §A — ⭐⭐ Butir 1: sumber tipe parameter — **TIDAK ADA PERSELISIHAN**

### A1 — Ronde 4 salah, dan sebabnya alat, bukan korpus

Ronde 4 melaporkan **8 parameter bertipe berselisih** antara dua sumber, dan menjadikannya alasan
menggeser dasar vonis **ADR-0003**. ⛔ **Perselisihan itu tidak ada.**

`[terverifikasi]` **`pyParametersParamName` tinggal di ENAM wadah yang berbeda**, dan hanya **satu**
di antaranya memuat parameter milik activity itu sendiri:

| Wadah induk | Kemunculan | Isinya apa |
| --- | ---: | --- |
| `pyStepsParamUI` *(di bawah `rowdata`)* | 2.817 | nama **slot parameter metode** — `PropertiesName`, `PropertiesValue`, `Property` |
| `pyStepsParamUI` *(di bawah `pyStepsCallParams`)* | 2.247 | parameter **panggilan** — `pyReportName`, `pyReportClass` |
| `pySystemParameters` | 1.068 | bawaan Pega — `flowName`, `flowType`, `ReferencePageName` |
| `pyLocalParameters` | 1.059 | **nyaris seluruhnya bernama kosong** |
| `pzRuleParameters` | 903 | parameter rule **layar** *(Harness · Section)* |
| ⭐ **`pyParameters` langsung di bawah `pagedata`** | **522** | ⭐ **parameter milik rule itu sendiri** |
| `pyParameters` di bawah `pyFunctionData` | 294 | parameter **fungsi** di rule `When` — `lValue`, `comparator`, `rValue` |

⛔ **Ronde 4 memasangkan nama ke tipe dengan `zip` atas seluruh berkas** — nama ke-*n* dipasangkan
dengan tipe ke-*n*, padahal nama datang dari tujuh wadah dan tipe tidak. **Pasangannya bergeser**,
dan pergeseran itulah yang terbaca sebagai "perselisihan".

### A2 — Dibaca dari wadah yang benar: **kedelapan SEPAKAT**

`[terverifikasi]`

| Activity | Parameter | sumber A *(tanda tangan)* | sumber B *(wadah benar)* | cara ronde 4 |
| --- | --- | --- | --- | --- |
| ⭐⭐ `SetKomiteList_ACT` | **`TotalAdjustment`** | `Double` | ⭐ **`Double`** | ~~`Decimal`~~ |
| ⭐ `GetCoverageAneka_Act` | **`TSI`** | `Decimal` | ⭐ **`Decimal`** | ~~`String`~~ |
| ⭐ `SetCedant_act` | **`Share`** | `STRING` | ⭐ **`STRING`** | ~~`Decimal`~~ |
| `CheckDate_Act` | `DOL` | `Date` | `Date` | ~~`String`~~ |
| `GeneratePLA` | `idxobj` | `INTEGER` | `INTEGER` | ~~`int`~~ |
| `SetIndex_Act` | `Name` | `STRING` | `STRING` | ~~`String`~~ |
| `SetTreatyname` | `IndexObject` | `STRING` | `STRING` | ~~`int`~~ |
| `SetTreatyname` | `IndexObjectItem` | `STRING` | `STRING` | ~~`String`~~ |

### A3 — Sensus penuh: **nol perselisihan di seluruh korpus**

`[terverifikasi]` **Jendela: seluruh 482 berkas `.xml`.** Diuji **dua arah**:

```
sumber A  pyXMLSignature                   169 parameter di  72 berkas
sumber B  pyParameters < pagedata          265 parameter di 116 berkas
beririsan (nama sama)                      169
⭐ TIPE BERSELISIH                            0
hanya di A                                   0      <-- B superset murni
hanya di B                                  96      <-- seluruhnya rule DataPage
```

⭐ **Angka 169 sumber A cocok persis dengan angka ronde 4** — jadi yang keliru **bukan sumber A**,
melainkan **cara sumber B dibaca**. Dan **hanya-di-A = NOL** membuktikan sumber B **memuat seluruh
parameter sumber A**: keduanya bukan dua pendapat, melainkan **satu daftar yang sama dilihat dari
dua sisi**.

⚠️ Ke-96 yang hanya ada di B adalah rule **`DataPage`**, yang memang **tidak punya
`pyXMLSignature`**. Itu **bukan selisih pendapat**, melainkan **beda cakupan**.

### A4 — ⭐⭐ Akibatnya: dasar vonis **ADR-0003 PULIH ke bacaan ronde 3**

Ronde 4 menggeser dasar vonis ADR-0003 dari *"`TotalAdjustment` bertipe `Double`"* ke
*"uang lewat sebagai teks"*, **karena tipe `Double` dianggap tidak pasti**. ⛔ **Ketidakpastian itu
tidak pernah ada.**

`[terverifikasi]` **`TotalAdjustment` memang `Double`, dan kedua sumber menyetujuinya.**
`Share` memang **`STRING`**. `TSI` memang **`Decimal`**.

⭐ **Vonis ADR-0003 tetap MENENTANG, dan kini berdiri di atas DUA kaki, bukan satu:**
**(a)** uang dideklarasikan sebagai **`Double`** — bilangan pecahan biner, persis yang dilarang;
**(b)** uang lewat sebagai **teks** dan ditambal **koma-ke-titik di SQL**.
⛔ Kaki (a) tidak lagi perlu dicabut.

⚠️ **Yang TIDAK saya simpulkan:** bahwa presisi benar-benar hilang saat berjalan. Yang dibaca
**deklarasi tipe**, bukan jalannya nilai.

### A5 — ⛔ Satu angka ronde 4 yang **TIDAK DAPAT DIREKONSILIASI**

Ronde 4 menulis *"`pyParametersParamName` muncul **4.071** kali di 179 berkas `Activity`"*.
`[terverifikasi]` Dihitung ulang, **tidak satu pun cara menghasilkan 4.071**:

```
seluruh wadah, seluruh berkas                 9.566  di 390 berkas
yang teksnya BERISI                           8.320
di berkas Activity saja                       7.481  di 178 berkas
berisi DAN di berkas Activity                 7.315  di 178 berkas
```

⭐ **Jumlah berkas cocok (178 lawan 179), jumlah kemunculan tidak.** ⛔ **Belum punya data** tentang
bagaimana 4.071 diperoleh. ⚠️ Angka itu **tidak dipakai sebagai dasar apa pun di ronde ini**.

---

## §B — Butir 2: muara `Local.TotalAdjustment`

### B1 — ⚠️ Dua nama yang sama, dan keduanya nyata

⛔ **`Local.TotalAdjustment`** *(properti pada halaman `Local`)* **bukan** parameter
**`TotalAdjustment`**. Ronde 4 sudah memperingatkannya; ronde ini **memisahkan keduanya di setiap
kalimat**.

### B2 — Kedelapan berkas, dibaca satu per satu

`[terverifikasi]` **Jendela: seluruh 482 berkas, medan `PropertiesName` · `PropertiesValue` ·
`pyStepsPage` · `pyStepsActivityName` · `pyStepsPreCondParams` · `pyStepsTransParams`.**

| Berkas | Perannya terhadap `Local.TotalAdjustment` |
| --- | --- |
| `CreateKMTNo_Act` | **menulis** *(nol di awal, lalu akumulasi `+ .ValueAdjustment`)* · **membaca sekali** ke `Local.IsADj` |
| `SendPICProtect_Act` | **menulis** *(nol, lalu akumulasi)* · ⭐ **dipakai sebagai GERBANG** terhadap batas bawah |
| `SetKomiteList_ACT` | **menulis** dari `.ValueAdjustment` |
| `SetSalvageValue` | **menulis** dari `.ValueAdjustment` |
| `SetValueAdjusterFee` | **menulis** dari `.ValueAdjustment` |
| `SetNilaiResikoSendiri` | ⛔ **nol di dalam langkah** — hanya **mendeklarasikan parameter** bernama sama |
| `Section/InputAdjustment` | ⛔ **slot parameter KOSONG** di `pyDeferLoadRetrievalActivityParams` |
| `Section/ItemListEstimation` | ⛔ **slot parameter KOSONG** di wadah yang sama |

**Ringkas:** ⭐ **7 penulisan, 1 pembacaan, 1 gerbang** — dan **seluruh 7 penulisan menuju halaman
`Local`**.

### B3 — ⭐⭐ Vonis: **tidak pernah sampai ke kolom tersimpan**

`[terverifikasi]` `CreateKMTNo_Act` punya **71 penugasan properti**, **8** menuju `pyWorkPage` —
dan ⭐ **tidak satu pun berasal dari `Local.TotalAdjustment` maupun `Local.IsADj`**. Kedelapan itu
mengisi daftar objek, tanggal, mata uang, dan penanda komite.

⭐ **`Local.IsADj` ditulis dua kali dan dibaca NOL kali** — di seluruh 482 berkas, ia **tidak pernah
muncul sebagai nilai**.

⚠️ **Jebakan nama yang hampir menggigit:** pencarian `IsADj` yang tidak peka huruf ikut menangkap
**`IsAdjustment`**, **`IsAdjVal`**, dan **`IsAdjValue`** di enam berkas lain — **properti yang
berbeda sama sekali**. Ketiganya **dibuang dari hitungan** sesudah dibaca satu per satu.

⭐ **Kesimpulan:** `Local.TotalAdjustment` adalah **akumulator kerja**. Ia menentukan **apakah
sesuatu dikirim ke PIC** *(lewat gerbang di `SendPICProtect_Act`)*, tetapi **angkanya sendiri tidak
pernah disimpan**. ✅ Butir `[terbuka]` ronde 4 **ditutup dengan bukti.**

### B4 — ⛔ RALAT: pemanggil `SetKomiteList_ACT` **SATU, bukan dua**

Ronde 4 menulis *"Pemanggil `SetKomiteList_ACT`: **dua** — `SetListKomite_act` dan
`SetNilaiResikoSendiri`"*.

`[terverifikasi]` ⛔ **`SetListKomite_act` tidak memanggilnya.** Satu-satunya kemunculan nama itu di
dalam berkasnya adalah **namanya sendiri** *(`pyRuleName`, `pyLabel`, `pyActivityName`,
`pyJavaClassName`)*. ⭐ **Ia rule yang BERBEDA dengan nama yang nyaris sama.**

`[terverifikasi]` Keduanya sama-sama nyata dan sama-sama dipanggil — **oleh pemanggil yang berbeda**:

| Rule | Pemanggil |
| --- | --- |
| **`SetKomiteList_ACT`** | ⭐ **satu** — `SetNilaiResikoSendiri` |
| **`SetListKomite_act`** | ⭐ **tiga** — `CreateKMTNo_Act` · `SetSalvageValue` · `SetValueAdjusterFee` |

✅ Ini sekaligus **menutup butir §F4-3 ronde 4** *("baca `SetListKomite_act`, rule bernama nyaris
sama")*. ⚠️ **Dan ia contoh keempat aturan identitas rule di modul ini: nama mirip bukan bukti.**

---

## §C — Butir 3: tujuh pesan `CheckEstimateValue`

### C1 — ⭐ Mekanisme pemblokirnya ADA, dan pengembangnya menuliskannya sendiri

`[terverifikasi]` Satu-satunya langkah yang menyetel **`pyWorkPage.IsError`** di rule ini adalah
**langkah 12.7**, bernilai **`2`**, bergerbang `Local.CompareDate=="1" || Local.CompareTGL=="false"`
— yaitu **kedua pesan TANGGAL**.

⭐⭐ **Arti nilai `2` terbaca dari catatan langkah pengembang di rule lain** — `SetSalvageValue`
mencatat: *"set `pyWorkPage.IsError` ==2 **to disable send komite**"*, dan `SetValueAdjusterFee`
mencatat *"enable button send …"*.

`[terverifikasi]` **`pyWorkPage.IsError` DIBACA** — di **tiga gerbang**: `ProtectionObjectItem_Act`
dan `SetInitial_ACT` *(keduanya `IsError>1`)*, serta `SetValueAdjusterFee` *(`IsError<2`)*.
⛔ **Ini kebalikan preseden Claim Prop**, tempat `IsError` **tidak pernah diset**.

### C2 — Tabel tujuh pesan

`[terverifikasi]`

| Pesan | Syarat munculnya | Penanda yang menyusul | Akibat yang terbaca |
| --- | --- | --- | --- |
| lgk **7** | persentase estimasi **< 0** | ⛔ tidak ada | **peringatan** |
| lgk **8** | nilai estimasi **< 0** | *(syaratnya ikut di gerbang 23)* | **peringatan** — kecuali bila syarat 23 lengkap terpenuhi |
| ⭐ lgk **12.2** | tanggal estimasi ditolak *(`CompareDate`)* | ⭐ **12.7 → `IsError = 2`** | ⭐⭐ **MEMBLOKIR** — tombol kirim komite mati |
| ⭐ lgk **12.5** | tanggal estimasi ditolak *(`CompareTGL`)* | ⭐ **12.7 → `IsError = 2`** | ⭐⭐ **MEMBLOKIR** — tombol kirim komite mati |
| lgk **12.6** | persentase **dan** nilai keduanya **0** | **12.8 → `Local.ErrSts = 1`** | **tombol CFS saja** |
| lgk **21** | batas tanggung jawab kosong/nol, dan hasil melampauinya *(`IsFire`)* | ⛔ tidak ada | **peringatan** |
| lgk **22** | estimasi melebihi TSI | **23 → `Local.ErrSts = 1`** | **tombol CFS saja** |

⭐ **`Local.ErrSts` hanya dibaca di langkah 24**, yang menyetel **`CFS`** dan **`PrintFaceClaim`** —
✅ **menguatkan ronde 4 §C2**: kedua titik kode-5 itu memang **soal tombol unduh, bukan uang**.

### C3 — ⛔ Yang tetap **`[terbuka]`**, dan ia penting

⛔ **Apakah pesan Pega itu sendiri menghalangi penyerahan** *(perilaku bawaan: pesan yang menempel
pada properti membuat flow action gagal disimpan)* **tidak terbaca dari ekspor.**

`[terverifikasi]` Rule ini **tidak dipanggil sebagai rule validasi penyerahan**. Ketiga layar
pemanggil — `Section/Estimasi` *(8 titik)*, `EstimasiMarine` *(6)*, `EstimasiPA` *(4)* — memanggilnya
lewat **`pyActionSets` dan `pyBehaviors`**, yaitu **aksi sisi-klien**. Peristiwa yang terdaftar di
berkas itu: `change` **81** · `click` **36** · `keyboard` **14** · `doubleclick` **4**; tindakannya
`refresh` **76** · `postValue` **28** · `runActivity` **6** · `save` **6**.

⭐ **Artinya ia validasi yang berjalan SAAT PENGGUNA MENGETIK**, bukan saat menekan simpan. Apakah
pesannya bertahan sampai penyerahan **bergantung pada perilaku runtime Pega**, dan ⛔ **itu tidak
ada di korpus.** `[terbuka]`

---

## §D — Butir 4: tiga belas rule ber-`COMMIT`

### D1 — ⭐ Ternyata **DUA jenis yang sangat berbeda**, bukan satu kelompok

`[terverifikasi]`

| Jenis | Berapa | Bentuknya | Artinya |
| --- | ---: | --- | --- |
| **Metode langkah Pega `Commit`** | **5** *(Activity)* | langkah bermetode `Commit`, berparameter *"release the lock when we commit the changes"* | **commit milik Pega sendiri** — menutup pekerjaannya dan melepas kunci |
| ⭐⭐ **`COMMIT;` di dalam SQL** | **8** *(RDBList)* | perintah `COMMIT;` tertulis di dalam badan SQL | ⭐ **commit tingkat basis data** — menutup **transaksi sesi**, termasuk pekerjaan pemanggilnya yang belum selesai |

⛔ **Yang berbahaya adalah kelompok kedua**, dan ronde 4 menyatukan keduanya sebagai satu angka.

### D2 — Pemanggil tiap satu

`[terverifikasi]`

| Rule | Pemanggil | Catatan |
| --- | ---: | --- |
| ⭐ `SaveOSClaim_SQL` | **3** — `CloseClaim` · `SaveAcceptation` · `SaveCFS_ACT` | ⭐ lihat D3 |
| ⭐ `GetSequenceNumber_SQL` | **4** — `CLaimFaceSheet_Act` · `DLAFacintoTreaty_Act` · `GenerateDLAFacin_Act` · `GeneratePLATreaty_Act` | **pembangkit nomor urut** |
| `GetTokenStorage_SQL` | 2 — `GetUrlGoogleStorage_Act` · `InsertGoogleStorage_Act` | penerbitan token berkas |
| `InsertClaimPNC` | 2 — `InsertJsonClaimNonMBU_act` · `SaveAdjustmenttoDB_ACT` | |
| `InsertProgressClaim_SQL` · `InsertSUBProgressClaim_SQL` | 1 — keduanya `InsertProgressClaim` | |
| `InsertLOGDirectKasir_SQL` | 1 — `HitServiceToKasir_Act` | |
| `UpdateDCauseOfLoss` | 1 — `CNMInsertDetailCauseOfLoss` | |
| `InsertLogServiceClaim` | 2 — `CLaimFaceSheet_Act` · `ChooseDla_Act` | ⭐ **`COMMIT;` di SQL, tetapi rule-nya Activity** |
| `DLAFacintoTreaty_Act` | 1 — `ChooseDla_Act` | metode `Commit` |
| `SetProtectionEstimation` | 1 — `SetInitial_ACT` | metode `Commit` |
| `SetSalvageValue` | 1 — `SetGrossAdjustment_act` | metode `Commit` |
| `GetPayAttachmentAdj_Act` | 1 — `GetPayAttachment_Act` | metode `Commit` |

### D3 — ⭐ Di mana ia jatuh relatif terhadap penyimpanan

`[terverifikasi]`

| Pemanggil | Letak panggilan | Keadaan |
| --- | --- | --- |
| `SaveAcceptation` | `SaveOSClaim_SQL` di langkah **10**, `Obj-Save` di langkah **14** | ⭐ **panggilan itu BER-REMARK `//`** — mati. `Obj-Save` hidup |
| ⭐ `CloseClaim` | `SaveOSClaim_SQL` di langkah **8** dan **8.3** | **hidup**, dan ⛔ **nol `Obj-Save` di activity ini** |
| ⭐ `SaveCFS_ACT` | `SaveOSClaim_SQL` di langkah **19** | **hidup** |
| `SaveAdjustmenttoDB_ACT` | `InsertClaimPNC` di langkah **6** | **hidup** *(activity 6 langkah)* |

⚠️ **Bentuk yang perlu diwaspadai terbaca jelas:** sebuah rule **bernama "List"** *(`RDBList`)* —
yang secara nama **membaca** — ternyata **menyimpan dan menutup transaksi**. ⛔ Nama rule di korpus
ini sudah berkali-kali terbukti menyesatkan.

### D4 — Bandingkan Claim Prop

Claim Prop mencatat *"15 dari 58 rule `COMMIT` sendiri"* dan **tidak mereplikasinya**.
⛔ **Bentuknya tidak sama persis:** angka Claim Prop **tidak memisahkan** metode `Commit` Pega dari
`COMMIT;` dalam SQL, sedangkan ronde ini **memisahkannya**. ⛔ **Keputusan Claim Prop tidak boleh
disalin begitu saja**, dan ronde ini **tidak menyalinnya**.

---

## §E — Butir 5: tiket `Register_Flow`

### E1 — ⭐ Tujuh tiket itu sebenarnya **DUA**

`[terverifikasi]`

| Menggantung di | Nama tiket |
| --- | --- |
| `Assignment7` | ⭐ **`SendToEstAdmin`** |
| `Assignment3` | ⭐ **`AcceptanceClaim`** |
| `Assignment1` · `Decision4` · `Decision5` · `Decision8` · `END52` | ⛔ **cangkang kosong** — tanpa nama, tanpa pengenal |

**Pengubah:** `Ticket5` = **`SendToEstAdmin`** · `Ticket1` = **`AcceptanceClaim`** — ⭐ **definisi
dari dua tiket yang sama**.

⭐ **Jadi: 2 tiket bernama · 5 cangkang kosong · 2 definisi.** ⛔ Bukan sembilan jalur perkecualian.

### E2 — ⭐⭐ `AcceptanceClaim` adalah tiket yang SAMA dengan Claim Prop

`[terverifikasi]` `Register_Flow` berkelas **`ASM-FW-GCNMFW-Work-PNC`**, dan `Flow_TreatyIn` Claim
Prop memikul `pyBaseClass` **`ASM-FW-GCNMFW-Work-PNC`** pada langkah yang menggantungi tiket bernama
sama. ⭐ **Ketiga modul klaim — Prop, Non Prop, Fac In — mewarisi satu model alur yang sama.**

✅ **Dan ia tetap LABEL MATI:** `Obj-Set-Tickets` **NOL** di seluruh **20 korpus**; `SetTicket` hanya
hidup di `NB FacIn`, `RNW Fac In`, `Endorsment Fac In`. ⛔ Tidak ada yang melemparnya.

### E3 — ⭐ `END52`: **warisan salinan, bukan kebetulan**

`[terverifikasi]`

| Modul | Kelas bentuk | `pyWorkStatus` |
| --- | --- | --- |
| Claim Fac In | `ASM-FW-GCNMFW-Work-PNC` | `Resolved-Completed` |
| Komite Claim Prop | `ASM-FW-GCNMFW-Work-KomiteTreaty` | `Resolved-Completed` |

⭐ **Pengenal bentuknya sama, kelasnya berbeda.** Itu tanda **alur disalin lalu diubah kelasnya** —
bukan rule bersama, dan bukan kebetulan. ✅ Menutup butir §F4-5 ronde 4.

---

## §F — RALAT terhadap ronde sebelumnya — **LIMA**

| # | Yang diralat | Kalimat/angka LAMA, dikutip | Yang benar |
| --- | --- | --- | --- |
| ⭐⭐ **1** | tipe parameter | ronde 4 §B4: *"⭐ **TIPENYA BERSELISIH** 8"* dan *"Delapan parameter **belum punya tipe yang pasti**"* | ⭐ **NOL berselisih.** Sumber B dibaca dari wadah yang salah; dibaca dari `pyParameters` di bawah `pagedata`, **169 dari 169 sepakat** |
| ⭐⭐ **2** | dasar ADR-0003 | ronde 4 §F5: *"dasarnya **bergeser** dari tipe `Double` ke uang-sebagai-teks"* | ⭐ **tidak perlu bergeser.** `TotalAdjustment` memang `Double` menurut **kedua** sumber. Vonis **MENENTANG** kini berkaki **dua** |
| **3** | pemanggil | ronde 4 §B2: *"Pemanggil `SetKomiteList_ACT`: **dua** — `SetListKomite_act` dan `SetNilaiResikoSendiri`"* | ⭐ **satu.** `SetListKomite_act` **rule berbeda**, dan ia sendiri dipanggil **tiga** activity lain |
| **4** | cacah nama parameter | ronde 4 §B4: *"`pyParametersParamName` muncul **4.071** kali di 179 berkas `Activity`"* | ⛔ **belum punya data** — tidak satu pun cara menghasilkan 4.071 *(9.566 · 8.320 · 7.481 · 7.315)*. Jumlah **berkas** cocok |
| **5** | jalur perkecualian | ronde 4 §D4: *"⭐ **9 `Data-MO-Event-Exception`**"* | ⭐ angkanya benar, **tetapi isinya bukan sembilan jalur**: **2 tiket bernama + 5 cangkang kosong + 2 definisi** |

---

## §G — Butir `[terbuka]`

### G1 — ✅ DITUTUP ronde ini — **tiga**

| Butir | Ditutup oleh |
| --- | --- |
| **muara `Local.TotalAdjustment`** | §B3 — nol penulisan ke halaman tersimpan, dari 71 penugasan |
| **`SetListKomite_act` yang belum dibaca** | §B4 — rule berbeda, bukan pemanggil |
| **kecocokan `END52` dua modul** | §E3 — warisan salinan, kelas berbeda |

### G2 — ⚠️ BARU — **satu**

| Butir | Pemilik |
| --- | --- |
| ⚠️ **Apakah pesan Pega menghalangi penyerahan** — rule ini dipanggil sebagai **aksi sisi-klien saat mengetik**, bukan validasi penyerahan. Perilakunya **tidak ada di korpus** | work owner + tim Pega |
| ⭐ **BARU dari keputusan K2** — **pesan mana yang menggagalkan simpan dan mana yang hanya memperingatkan**; ditambah **berapa banyak kasus lama** yang memuat keadaan yang akan menjadi galat | work owner |

### G3 — ⛔ Masih terbuka dari ronde sebelumnya

**Arti `REQUIRED = -1` lawan `0`** pada tanda tangan parameter — ⛔ tidak dipakai sebagai dasar apa
pun sampai sekarang.

---

## §H — Dua pertanyaan untuk work owner

### H1 — Perlakuan delapan `COMMIT;` di dalam SQL

**Apa yang ditanyakan.** Delapan `RDBList` menutup transaksi basis data dari dalam dirinya sendiri.
Tiga di antaranya — `SaveOSClaim_SQL` di `CloseClaim` dan `SaveCFS_ACT`, serta `InsertClaimPNC` di
`SaveAdjustmenttoDB_ACT` — dipanggil **di tengah pekerjaan yang lebih besar**.

➡️ **Rekomendasi: jangan replikasi `COMMIT;`-nya; satu aksi pengguna = satu transaksi.**
⭐ Alasannya bukan kerapian: `COMMIT;` di tengah membuat **kegagalan separuh jalan meninggalkan data
separuh tersimpan**, dan itu **tidak dapat dipulihkan** oleh pemanggilnya.
⚠️ **Kecualikan `GetSequenceNumber_SQL`** — ia pembangkit nomor urut, dan **ADR-0006** sudah
memutuskan penomoran tetap lewat stored procedure. Nomor urut **memang harus** bertahan walau
transaksi induknya batal, kalau tidak nomornya terpakai ulang.

### H2 — Kedua pesan tanggal yang memblokir

**Apa yang ditanyakan.** Dari tujuh pesan, **hanya dua** yang benar-benar menghalangi — dan caranya
**mematikan tombol kirim komite** lewat `IsError = 2`, bukan menggagalkan simpan.

➡️ **Rekomendasi: di sistem baru, jadikan ketujuhnya GALAT VALIDASI yang tegas, dengan pembedaan
yang terbaca** — mana yang **menggagalkan penyimpanan** dan mana yang **hanya memperingatkan**.
⭐ Alasannya: bentuk sekarang **menyembunyikan aturan bisnis di dalam nilai sebuah penanda** yang
artinya hanya diketahui dari catatan pengembang. ⚠️ **Ini penyimpangan sadar** bila diterima, dan
⛔ **daftar mana-memblokir-mana-tidak harus datang dari work owner** — korpus hanya memberi
**keadaan sekarang**, bukan **niatnya**.

---

## §I — Kesimpulan ronde ini yang **PALING RAWAN SALAH**

⚠️⚠️ **§C2 — tabel tujuh pesan.**

**Kenapa rawan:** ia memetakan pesan ke penanda **lewat kesamaan syarat gerbang**, bukan lewat
tautan yang tertulis. Langkah 12.7 tidak menyebut *"pesan 12.2 dan 12.5"* — ia hanya kebetulan
bergerbang pada **dua variabel yang sama**. ⛔ Kalau sebuah gerbang punya sumber lain yang belum
saya baca, pemetaan itu bergeser. ⚠️ Pesan **8** paling rawan: syaratnya muncul di gerbang langkah
23, tetapi **tergabung dalam rantai DAN** dengan syarat lain.

**Yang kokoh dan tidak rawan:** ⭐ **§A** — ia berdiri di atas **uji dua arah yang tutup sempurna**
*(hanya-di-A = NOL, irisan = 169 = angka ronde 4 sendiri)*, dan ⭐ **§B4** — dibaca langsung dari
medan, bukan disimpulkan.

**Rawan kedua:** §D3 — *"nol `Obj-Save` di `CloseClaim`"*. Saya menyisir **nama metode langkah**;
penyimpanan bisa terjadi lewat rule yang dipanggilnya.

---

## §J — ⭐ Kesiapan untuk spec

⭐⭐ **VONIS: SIAP dispesifikasikan — dengan dua butir yang harus dijawab lebih dulu.**

**Yang sudah cukup untuk menulis spec:** daur hidup kasus *(8 bentuk, 10 penghubung, bentuk selesai
`Resolved-Completed`)* · tipe seluruh 265 parameter **tanpa perselisihan** · muara
`Local.TotalAdjustment` · perilaku ketujuh pesan · peta 13 rule ber-`COMMIT` · enam keputusan
rancangan ronde 2 §A5 · vonis ADR-0003 dan ADR-0012.

**Yang menahan, dan keduanya pertanyaan di §H:** perlakuan **`COMMIT;`** dan **daftar pesan yang
memblokir**. ⛔ Keduanya menentukan **AC**, bukan sekadar uraian — menebaknya berarti menulis ulang
spec belakangan.

⚠️ **Dan satu hal yang bukan pertanyaan melainkan kekurangan bahan:** modul ini **belum punya
struktur tabel maupun relasi tabel**. Keduanya **tidak menahan spec**, tetapi **menahan tiket** —
dan di Claim Prop keduanya dikerjakan **sesudah** spec.

---

## §K — ✅ Dua keputusan work owner — 2026-09-19

`[keputusan work owner — atas rekomendasi asisten]` **2026-09-19**. ⛔ Tanda *"atas rekomendasi
asisten"* dipakai **sengaja**: pilihannya datang dari rekomendasi di §H, dan **pencabutannya harus
murah**. ⛔ **Nol di antaranya mengubah korpus** — keduanya tentang sistem baru.

### K1 — `COMMIT;` di dalam SQL **tidak direplikasi**

> **Satu aksi pengguna, satu transaksi.** ⛔ Kedelapan `COMMIT;` yang tertulis di dalam badan SQL
> **tidak dibawa** ke sistem baru.

⚠️ **Satu pengecualian: `GetSequenceNumber_SQL`.** Ia **pembangkit nomor urut**, dan nomornya
**harus bertahan walau transaksi induknya batal** — kalau ikut dibatalkan, nomor yang sama
**terpakai ulang** oleh kasus berikutnya. ✅ Sejalan **ADR-0006**, yang sudah menetapkan penomoran
tetap lewat stored procedure.

**Kenapa diputuskan begitu — alasannya, bukan sekadar isinya.** `COMMIT;` di tengah pekerjaan orang
lain membuat sebuah kegagalan separuh jalan meninggalkan **data separuh tersimpan**, dan
⛔ **pemanggilnya tidak dapat memulihkannya** — transaksinya sudah ditutup oleh rule yang ia
panggil, bukan oleh dirinya.

⚠️ `[terverifikasi]` **Tiga panggilan memang jatuh di tengah pekerjaan yang lebih besar:**
`SaveOSClaim_SQL` di **`CloseClaim`** dan **`SaveCFS_ACT`**, serta `InsertClaimPNC` di
**`SaveAdjustmenttoDB_ACT`**. ⭐ **Dan satu yang justru menenangkan:** di `SaveAcceptation`
panggilan itu **ber-remark** — mati — sedangkan `Obj-Save` di langkah berikutnya **hidup**. Jadi
jalur akseptasi **sudah** berperilaku seperti yang diputuskan di sini.

⚠️ **Perlu diingat saat membangun:** rule-rule ini bernama **`RDBList`** — yang menurut namanya
**membaca** — padahal **menyimpan dan menutup transaksi**. ⛔ Nama rule di korpus ini sudah
berkali-kali terbukti menyesatkan.

⭐ **Lima rule ber-metode `Commit` Pega TIDAK termasuk keputusan ini** — itu commit milik Pega
sendiri yang menutup pekerjaannya dan melepas kunci, bukan `COMMIT;` yang disisipkan ke dalam SQL.

### K2 — Ketujuh pesan `CheckEstimateValue` menjadi **galat validasi yang tegas**

> **Ketujuhnya menjadi galat validasi**, dengan pembedaan **yang terbaca** antara yang
> **menggagalkan penyimpanan** dan yang **hanya memperingatkan**.

⚠️ **Penyimpangan sadar.**

**Kenapa diputuskan begitu.** Bentuk sekarang **menyembunyikan aturan bisnis di dalam nilai sebuah
penanda**: `pyWorkPage.IsError = 2`. ⛔ Arti angka **2** itu **tidak tertulis di mana pun sebagai
aturan** — ia hanya dapat diketahui dari **catatan langkah pengembang** pada rule lain, yang
berbunyi *"to disable send komite"*. ⭐ Aturan yang hanya hidup di dalam sebuah catatan bebas
**tidak dapat diuji**, dan **hilang begitu catatannya hilang**.

⛔ **Yang BELUM diputuskan, dan ia `[terbuka]`:** **pesan mana masuk golongan mana.** Korpus hanya
memberi **keadaan sekarang** — 2 menahan *(kedua pesan tanggal)* · 2 hanya mematikan tombol CFS ·
2 peringatan · 1 belum punya data — dan **bukan niatnya**. ⚠️ Keadaan sekarang **tidak boleh
dipakai sebagai jawaban**: justru ketimpangannya yang melahirkan keputusan ini. Pemilik:
**work owner**.

⚠️ **Akibat yang harus ikut dijawab:** bila sebuah pesan yang hari ini **hanya peringatan**
dinaikkan menjadi **penggagal simpan**, maka data yang **hari ini bisa disimpan** menjadi
**tidak bisa** — dan **kasus lama hasil migrasi mungkin sudah memuat keadaan itu**. ⛔ Belum
diperiksa berapa banyak. `[terbuka]`

---

## Lampiran — bukti berkas lain tidak disentuh

```
korpus Claim Fac In        482 berkas .xml   — nol dibuka untuk ditulis
korpus Claim Prop          329 berkas .xml   — nol dibuka untuk ditulis
korpus Komite Claim Prop    80 berkas .xml   — nol dibuka untuk ditulis
docs\adr\                   15 berkas        — nol disunting
berkas lama claim-facin      5 berkas        — nol disunting
spec.md kedua modul Prop                     — nol disunting
```

⛔ **Nol kode Go/React · nol `CREATE TABLE` · nol DDL · nol nomor baris XML dikutip · nol nilai
rahasia disalin · nol keputusan work owner dibuat · nol ADR direvisi.**
