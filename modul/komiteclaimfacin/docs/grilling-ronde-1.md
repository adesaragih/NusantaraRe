# Grilling Ronde 1 — Komite Claim Fac In · SENSUS STRUKTURAL

**Tanggal:** 2026-09-19 · **Korpus:** `D:\XML\RNM_BRD\Komite Claim FacIn\` — **114 berkas**, READ-ONLY
**Pembanding:** `Claim Fac In` *(482)* · `Komite Claim Prop` *(80)*

> ⛔ **Ronde ini MEMBACA, bukan merancang.** Nol kesimpulan untuk aplikasi Go. Nol spec, nol tiket.
> ⛔ Nol butir `[terbuka]` dinyatakan tertutup · nol keputusan work owner dibuat · nol ADR direvisi.
> ⛔ Nol berkas korpus dibuka untuk ditulis · nol nilai rahasia disalin.

---

## §1 — Inventaris dan identitas

### T1a — Dihitung dua cara, dan keduanya sepakat

| | |
| --- | ---: |
| **cara A** — glob rekursif `**/*.xml` | **114** |
| **cara B** — `os.walk` | **114** |

✅ **SEPAKAT.** `[terverifikasi]`

| | Jumlah |
| --- | ---: |
| identitas **empat bagian** unik | **112** |
| identitas muncul di **lebih dari satu** berkas | **2** |
| ⭐ **uji tutup** — `112 unik + 2 lebihan = 114 berkas` | ✅ **COCOK** |
| ⭐ bersama **Claim Fac In** | ⭐ **81** |
| bersama **Komite Claim Prop** | **40** |
| ⭐ **khas modul ini** | ⭐ **23** |

**Kedua identitas rangkap** *(masing-masing di dua berkas)*: `ShowSecurityReinsurer` dan
`SpreadingDetail` — keduanya **`FlowAction` + `Section`**. ⭐ Pola yang sama dengan Claim Fac In,
tempat satu rule layar sering diekspor sebagai dua jenis sekaligus.

⭐⭐ **Modul ini sebagian besar PINJAMAN: hanya 23 dari 112 identitas yang tidak ada di dua modul
saudaranya.**

**Sebaran jenis rule** `[terverifikasi]`:

| Jenis | Jumlah |
| --- | ---: |
| ⭐ **`When`** | ⭐ **49** |
| `Activity` | 28 |
| `RDBList` | 17 |
| `FlowAction` · `Section` | 5 masing-masing |
| `ConnectREST` · `ReportDefinition` | 3 masing-masing |
| `DataTransform` · `DecisionTable` · `Flow` · `SystemSettings` | 1 masing-masing |

⭐ **43% modul ini adalah rule `When`.** Bandingkan Claim Fac In: **60 dari 482** *(12%)*.

### T1b — Kelas: **23 berbeda**

| Kelas | Jumlah | Mewakili |
| --- | ---: | --- |
| `@baseclass` | **26** | rule lintas-aplikasi |
| ⭐ **`ASM-FW-GCNMFW-Work-Komite`** | ⭐ **15** | ⭐ **objek kerja kasus komite modul ini** |
| `ASM-FW-GISFW-Data` | 13 | data bersama |
| `ASM-FW-GISFW-Int-T_STORAGE_IMAGE` | 10 | penyimpanan berkas |
| `ASM-FW-GCNMFW-Data-Adjustment` | 9 | baris penyesuaian |
| `ASM-FW-GCNMFW-Work` | 8 | objek kerja induk |
| `ASM-FW-GISFW-Work` | 6 | objek kerja GISFW |
| `ASM-FW-GISFW-Data-SpreadingRisk` | 4 | pembagian risiko |
| `ASM-FW-GCNMFW-Data-ObjectItem` | 4 | item objek pertanggungan |
| `ASM-FW-GISFW-Data-FacOffer` | 3 | penawaran fakultatif |

> ### ⭐⭐ TIGA MODUL KOMITE, TIGA KELAS KERJA BERBEDA
>
> | Modul | Kelas kerja |
> | --- | --- |
> | Claim Fac In *(klaimnya)* | `ASM-FW-GCNMFW-Work-PNC` |
> | Komite Claim Prop | `ASM-FW-GCNMFW-Work-KomiteTreaty` |
> | ⭐ **Komite Claim Fac In** | ⭐ **`ASM-FW-GCNMFW-Work-Komite`** |
>
> ⚠️ **Namanya paling pendek dan paling umum dari ketiganya** — `Work-Komite`, tanpa penanda lini.
> ⛔ **Tidak saya simpulkan** apakah itu berarti ia kelas komite generik yang dipakai lini lain;
> korpus ini **tidak memuat modul lini lain**. `[terbuka]` **butir 1**

### T1c — Ke-23 rule khas, satu per satu

⭐ **Di sinilah isi modul ini yang sebenarnya.**

| Berkas | Kelas |
| --- | --- |
| `Activity\ApprovalKomite_Act` | `Work-Komite` |
| `Activity\InsertJsonClaimNonMBU_act` | `Work-Komite` |
| ⭐ `Activity\KomitePostAct` | `Work-Komite` |
| ⭐ `Activity\KomitePost_Adjustment` | `Work-Komite` |
| ⭐ `Activity\KomitePost_CloseClaim` | `Work-Komite` |
| ⭐ `Activity\KomitePost_Reject` | `Work-Komite` |
| ⭐ `Activity\KomiteRouter` | `Work-Komite` |
| `Activity\PrintPDFAccep_MultiAksep_KMT` | `Work-PNC` |
| `Activity\SaveAccept_ACT` | `Data-ObjectItem` |
| `Activity\SaveAcceptation_KMT` | `Data-Adjustment` |
| `Activity\SaveReject_ACT_KMT` | `Data-ObjectItem` |
| `Activity\SetDataForInformation_Act` | `Work-Komite` |
| `Activity\SetProteksiSubmiteKomite` | `Work-Komite` |
| `Activity\SetValueKomite` | `Work-Komite` |
| `FlowAction\DetailAdjustmentFac` · `Section\DetailAdjustmentFac` | `Data-Adjustment` |
| `FlowAction\ViewTransferDtl` | `Work-Komite` |
| ⭐ `Flow\Komite_Flow` | `Work-Komite` |
| `RDBList\UpdateSubProgresKlaim` | `Work-Komite` |
| `Section\ShowTransfer` | `Work-Komite` |
| `Section\SpreadingDetail` | `Data-ObjectItem` |
| `When\IsCustomBonds` | `Work` |
| ⭐ `When\IsKomiteLoop` | `Work-Komite` |

⭐ **Empat rule `KomitePost_*` adalah jantung modul ini** — `Adjustment` · `CloseClaim` · `Reject`,
ditambah `KomitePostAct`. ⚠️ Bandingkan Komite Claim Prop, yang punya `KomitePostAdjustment` ·
`KomitePost_Close` · `KomitePost_Reject` — **nama nyaris sama, identitas empat bagian BERBEDA**.
⛔ **Nama mirip bukan bukti.**

### T1d — Ruleset dan sistem sumber

| Ruleset | Berkas |
| --- | ---: |
| `GISFW` | **56** |
| `GCNMFW` | **52** |
| `SFAGIS` | 5 |
| ⚠️ **`ADESAMUEL@`** | ⚠️ **1** |

**Versi ruleset berbeda: 32.**

⚠️ **Pola ruleset bernama akun perorangan BERULANG** — di Claim Fac In ada **5** rule di
`ADESAMUEL@` *(termasuk berkas alurnya)*; di sini **satu**. ⛔ **Yang mana belum saya isolasi** —
kekurangan ronde ini. `[terbuka]` **butir 2**

| `pxUpdateSystemID` | Berkas | % |
| --- | ---: | ---: |
| `pegadevnusare2` | **57** | 50,0% |
| `pega` | **53** | 46,5% |
| ⭐ `pegaprdnusare` | **4** | 3,5% |

⭐⭐ **Pola TIGA sistem sumber berulang untuk modul ketiga.** ⚠️ Di Claim Fac In sebarannya
`pega` 320 · `pegadevnusare2` 146 · `pegaprdnusare` 16; di sini **`pegadevnusare2` justru
terbanyak**. ⛔ Artinya tidak disimpulkan.

⭐ `[keputusan work owner]` 2026-09-19 sudah menetapkan **salinan produksi yang ditiru** — jadi
⚠️ **keempat berkas `pegaprdnusare` itu perlu diisolasi**, sebab merekalah yang berlaku bila isinya
berbeda. ⛔ Belum dikerjakan ronde ini. `[terbuka]` **butir 3**

---

## §2 — Sensus langkah activity

**Jendela:** **28 berkas `Activity`**, dibaca dengan **pengurai XML rekursif**, langkah bersarang
terbaca dengan nomor bertitik.

### T2a — Bentuk dasar

| | Jumlah |
| --- | ---: |
| langkah | **541** |
| ber-remark `//` | **23** |
| bermetode `Java` | **11** |

### T2b — Keluarga gerbang dan arah

| Keluarga | Baris |
| --- | ---: |
| **1** `pyStepsPreCondParams` | **591** — ⚠️ lebih banyak dari jumlah langkah |
| **2** `pyStepsTransParams` | **541** — ⭐ **tepat satu per langkah** |

⭐ **Aturan "satu baris transisi per langkah" BERLAKU PENUH di sini** — 541 = 541, nol pengecualian.
⚠️ Di Claim Fac In ditemukan **satu** pengecualian.

**Arah keluarga-1** *(prakondisi)*:

```
2/2  266    2/3  242    _/3   25    3/2   24    2/6    9    6/2    8
3/_    5    _/1    3    _/_    3    5/2    3    _/6    1    5/_    1
```

**Arah keluarga-2** *(transisi)*: `2/2` **525** · `_/_` 6 · `1/2` 4 · `6/2` 1 · `2/6` 1 · `1/1` 1 ·
`1/_` 1 · `6/6` 1.

### T2c — ⭐ Kode khusus, tiap kemunculan diisolasi

| Kode | Arti | Jumlah |
| --- | --- | ---: |
| **1** Jump to Later Step | lompatan | **12** |
| **4** Exit Iteration | putus perulangan | ⭐ **0** |
| ⚠️ **5** Skip Whens | ⚠️ **mengubah rantai DAN menjadi ATAU** | **4** |
| **6** Exit Activity | keluar activity | **22** |

**Keempat titik kode 5, seluruhnya berflag `true` — jadi HIDUP:**

| Rule | Langkah | Sisi | Syarat |
| --- | --- | --- | --- |
| `SaveAccept_ACT` | 3 | benar | `IsFire` |
| `SaveAccept_ACT` | 3 | benar | `IsAneka` |
| `SendErrorDirectKasir` | 2.3 | benar | `IsCLMP` |
| ⚠️⚠️ `SetProteksiSubmiteKomite` | 2.1 | benar | ⭐ **`OperatorID.pyPosition=="IT Developer"`** |

> ### ⚠️⚠️ TEMUAN YANG PALING MENONJOL DI RONDE INI
>
> ⭐ **Pintu belakang pengembang muncul lagi — dan kali ini di dalam rantai ATAU.**
>
> `SetProteksiSubmiteKomite` langkah **2.1** menguji **`OperatorID.pyPosition == "IT Developer"`**
> dengan **kode 5 pada sisi benar**, artinya rantai gerbangnya **berperilaku ATAU**: begitu syarat
> itu terpenuhi, **sisa syaratnya tidak diuji lagi**.
>
> ⚠️ Rule bernama **"set proteksi submit komite"** — ⛔ **apa yang dilindunginya belum saya baca**,
> dan **saya tidak menyimpulkan apa pun**. ⭐ Tetapi bentuknya **persis** yang sudah tercatat di
> Claim Prop dan Komite Claim Prop: **wewenang digantungkan pada TEKS JABATAN**.
>
> ⭐ **Ini kemunculan KETIGA pola itu, dan yang pertama di dalam rantai ATAU.** `[terbuka]`
> **butir 4**

### T2d — Flag prakondisi

| Nilai | Langkah |
| --- | ---: |
| *(kosong)* | **268** |
| `true` | **233** |
| `false` | **37** |
| ⚠️ **`0`** | ⚠️ **3** |

⭐ **Nilai `0` yang tak berdokumen muncul lagi.** ✅ `[keputusan work owner]` 2026-09-19 sudah
menetapkannya **setara kosong** — jadi ketiganya dibaca sebagai **tanpa gerbang**.

⭐⭐ **76 baris syarat berada pada langkah berflag MATI** — syaratnya **tertulis tetapi tidak
berlaku**. ⚠️ Siapa pun yang membaca syarat itu tanpa memeriksa flagnya akan **salah membaca 76
aturan**.

### T2e — `Obj-*` dan `Commit`

| Metode | Jumlah |
| --- | ---: |
| `Obj-Save` | **8** |
| `Obj-Browse` | 7 |
| `Obj-Open-By-Handle` | 5 |
| `Obj-Refresh-And-Lock` | 5 |
| ⭐ **`Commit`** | ⭐ **NOL** |

**Total 25 langkah · ber-remark NOL.**

⭐⭐ **Modul ini tidak punya satu pun metode langkah `Commit` Pega.** ⚠️ Bandingkan Claim Fac In yang
punya 5. ⛔ **Artinya penutupan transaksinya bergantung sepenuhnya pada `COMMIT;` di dalam SQL** —
lihat §5.

⛔ **Berapa dari 25 itu punya jalur kegagalan belum saya hitung** — kekurangan ronde ini.
`[terbuka]` **butir 5**

---

## §3 — Berkas alur

`Flow\Komite_Flow.xml` · kelas **`ASM-FW-GCNMFW-Work-Komite`** · ruleset **`GCNMFW 01-01-21`**.

### T3a — Dibaca wadah demi wadah, dengan uji penjumlahan

| Wadah | Isi | Jumlah |
| --- | --- | ---: |
| ⭐ **`pyShapes`** | ⭐ **bentuk** | ⭐ **4** |
| `pyConnectors` | penghubung | **4** |
| `pyModifiers` | pengubah *(tiket tingkat alur)* | **1** |
| `pyTicketShapes` | tiket menggantung di bentuk | **3** |
| `pyRouterProp` | perute | 4 |
| `pyNotifyProp` | pemberitahu | 1 |

⭐ **UJI PENJUMLAHAN:** `4 + 4 + 1 + 3 + 4 + 1 = 17`, dan sapuan mentah seluruh `Data-MO-*` juga
**17**. ✅ **COCOK** — nol entri hilang, nol terhitung dua kali.

### T3b — Keempat bentuk

| Pengenal | Jenis | Nama | Status akhir |
| --- | --- | --- | --- |
| ⭐ `END52` | **Event-End** | — | ⭐ **`Resolved-Completed`** |
| `ASSIGNMENT63` | Activity-Assignment | **`KomiteRouter`** | — |
| `Decision1` | Gateway-Decision | **`KomiteLoop`** | — |
| `Start1` | Event-Start | — | — |

⭐ **Bentuk selesai ADA**, dan statusnya **selesai-tuntas**.

### T3c — Tiket: **satu bernama, dua cangkang kosong**

| Menggantung di | Nama |
| --- | --- |
| `END52` | ⛔ *(tanpa nama — cangkang kosong)* |
| `ASSIGNMENT63` | ⛔ *(tanpa nama — cangkang kosong)* |
| ⭐ `Decision1` | ⭐ **`komiteAccept_ticket`** |

⭐⭐ **Namanya SAMA PERSIS dengan Komite Claim Prop.**

`[terverifikasi]` **Jendela: seluruh 114 berkas.** Metode **`Obj-Set-Tickets`** → **NOL berkas**;
activity **`SetTicket`** → **NOL berkas**.

⭐ **Tiketnya LABEL MATI** — tidak ada yang membangkitkannya. ✅ **Sama dengan dua modul lain.**

⚠️ **Batas kejujurannya sama:** mesin pembangkit bawaan Pega **ada di aplikasi yang sama**
*(tiga modul lain punya `SetTicket`)*, jadi rule **di luar ekspor ini** secara teknis dapat
melemparnya. ⛔ Tidak ada buktinya. `[terbuka]` **butir 6**

### T3d — ⭐⭐ Pengenal bentuk: **warisan salinan, terbaca ketiga kalinya**

| Modul | Kelas | Pengenal bentuk |
| --- | --- | --- |
| **Komite Claim Prop** | `…Work-KomiteTreaty` | `END52` · `ASSIGNMENT63` · `Decision1` · `Start1` |
| ⭐ **Komite Claim Fac In** | `…Work-Komite` | ⭐ **`END52` · `ASSIGNMENT63` · `Decision1` · `Start1`** |
| Claim Fac In | `…Work-PNC` | `Decision8` · `Assignment7` · `Decision5` · `Decision4` · `Assignment3` · `END52` · `Assignment1` · `Start1` |

⭐⭐ **Keempat pengenal bentuk kedua modul komite IDENTIK, dan jumlahnya sama — tetapi kelasnya
berbeda.** ⚠️ Itu **bukan kebetulan dan bukan rule bersama**: identitas empat bagiannya berbeda
*(`Komite_Flow` termasuk 23 rule khas)*. ⭐ **Alur disalin lalu diubah kelasnya.**

⚠️ **Akibat praktisnya:** daur hidup kasus komite **kemungkinan besar sama persis** di kedua modul —
tetapi ⛔ **itu dugaan dari bentuk, bukan pembacaan isi**. `[terbuka]` **butir 7**

---

## §4 — Ke-49 rule `When`

**Jendela:** 49 berkas `When`, **tujuh medan** disisir — `pyConditionString` ·
`pyConditionValue1String` · `pyConditionValue1` · `pySimpleCondition` · `pyWhenExpression` ·
`pyExpression` · `pyLabel`.

| | |
| --- | ---: |
| kondisi **terbaca** | **48 dari 49** |
| ⛔ **tidak terbaca di tujuh medan** | **1** — ⭐ **`IsKomiteLoop`** |

⚠️ **`IsKomiteLoop` adalah salah satu dari 23 rule KHAS modul ini**, dan justru ia yang kondisinya
**tidak ketemu di tujuh medan yang disisir**. ⛔ **Bukan berarti kosong** — ia menuntut medan
kedelapan yang belum saya kenal. `[terbuka]` **butir 8**

### T4b · T4c — Properti yang diuji, dan penulisnya

| Properti yang diuji | Rule | Penulis di modul ini |
| --- | ---: | ---: |
| ⭐ `pyWorkPage.Quotation.BusinessType` | **40** | ⭐ **1** |
| `pyWorkPage.Quotation.BusinessCode` | 10 | ⛔ **0** |
| `pyWorkCover.pyWorkIDPrefix` | 3 | ⛔ 0 |
| `pyWorkPage.pyWorkIDPrefix` | 3 | ⛔ 0 |
| `pyWorkPage.OfferFacIn.QuotationData.BusinessType` | 2 | ⛔ 0 |
| `pyWorkPage.Quotation.BusinessName` | 1 | ⛔ 0 |
| `pxProcess.pzProductionLevel` | 1 | ⛔ 0 |
| `pxProcess.pxSystemNodeID` | 1 | ⛔ 0 |

**Properti berbeda 8 · punya penulis 1 · nol penulis 7.**

> ### ⭐⭐ SATU PENULIS ITU MENUTUP RANTAI YANG DI CLAIM FAC IN HARUS DITELUSURI EMPAT MATA
>
> `[terverifikasi]` **`Activity\SetValueKomite`** — salah satu dari **23 rule khas** —
> menulis:
>
> ```
> pyWorkPage.Quotation.BusinessType  =  pyWorkCover.OfferFacIn.QuotationData.BusinessType
> ```
>
> ⭐ **`pyWorkCover` adalah kasus klaim INDUK.** Jadi nilai klasifikasi lini **disalin dari kasus
> klaim ke kasus komite**, satu properti, satu penugasan, terbaca langsung.
>
> ✅ **Ini menguatkan temuan Claim Fac In** — di sana rantainya empat mata dan berakhir di modul
> saudara; di sini **mata rantai terakhirnya terbaca di dalam modul ini sendiri**.

⚠️ **Ketujuh properti tanpa penulis TIDAK berarti mati** — pola yang sama di Claim Fac In terbukti
diisi lewat **penyalinan halaman dari luar**. ⛔ **Belum diperiksa di sini.** `[terbuka]` **butir 9**

### T4d — Berapa yang dipakai bersama

⭐ **47 dari 49** rule `When` **identik empat bagian** dengan **Claim Fac In**.
Hanya **5** identik dengan Komite Claim Prop.

⭐⭐ **Himpunan `When` modul ini hampir seluruhnya milik Claim Fac In.** ⚠️ Keputusan
*"rule dipisahkan dari kehidupannya"* **berlaku langsung di sini** — satu salinan, dan hidup-matinya
bergantung pada halaman yang diujinya ada atau tidak pada objek kerja modul ini.

---

## §5 — Layar · RDB · ConnectREST

### T5a — Gerbang layar

**Jendela:** 5 `Section` + 5 `FlowAction`.

| Medan | Muncul | Tidak bernilai mati |
| --- | ---: | ---: |
| `pyVisible` | **519** | 519 |
| `pyIsVisibilityOption` | 52 | 52 |
| `pyDisabledNew` | **90** | ⭐ **11** |
| `pyRequiredNew` | **59** | ⭐ **2** |
| `pyReadOnlyCondition` | 3 | 3 |

⚠️ **`pyRequiredNew` muncul 59 kali tetapi hanya 2 yang bersyarat hidup** — pola yang sama dengan
Claim Fac In. ⛔ **"Tidak bernilai mati" bukan sama dengan "hidup"** — saya hanya menyaring nilai
`NEVER` / `1=2` / `false`; ⛔ **penyaring itu belum diuji instrumennya**. `[terbuka]` **butir 10**

### T5b — ⭐⭐ Medan wewenang: **MODUL KETIGA DENGAN POLA YANG SAMA**

**Jendela:** seluruh **114** berkas, **setiap tag** yang namanya mengandung
`privileg|authoriz|workgroup|accessgroup|role`.

| Medan | Hadir | Terisi | Isinya |
| --- | ---: | ---: | --- |
| `pyPrivilegeView` · `pyPrivilegeUpdate` | 82 masing-masing | ⭐ **0** | — |
| `pyPrivilege` · `pyWhenNotPrivilege` | 38 masing-masing | ⭐ **0** | — |
| `pyPrivilegeName` | 36 | ⭐ **0** | — |
| `pyActivityPrivilegeList` | 28 | ⭐ **0** | — |
| `pyPrivilegeList` | 8 | ⭐ **0** | — |
| `pySectionReferencePrivilege*` · `pyActionPrivilegeList` | 5 masing-masing | ⭐ **0** | — |
| `pyAllowedToEditValuesPrivileges` | 1 | ⭐ **0** | — |
| `pyAssociatedPrivileges` | 66 | 58 | ⛔ seluruhnya **`false`** — nilai baku |
| ⚠️ `pyPrivilegeClass` | 35 | **35** | ⚠️ **hanya nama KELAS** — `@baseclass`, `…Data-Adjustment` |

⭐⭐ **Tiga belas medan hak akses, dan nama haknya KOSONG SELURUHNYA — untuk modul ketiga
berturut-turut.** ✅ Claim Fac In: pola sama. ✅ Komite Claim Prop: pola sama.

⚠️ ⭐ **Itu memperkuat butir register 6 Claim Fac In** *(aturan peran belum ditulis)* menjadi
temuan **lintas tiga modul**, bukan ciri satu modul. ⛔ **Nol kesimpulan untuk Go ditarik dari sini.**

### T5c — 17 `RDBList`

| Aksi SQL | Jumlah |
| --- | ---: |
| `BEGIN … END` *(blok anonim)* | **9** |
| `SELECT` | 6 |
| `UPDATE` | 2 |

⭐⭐ **Sembilan dari 17 membawa `COMMIT` tertanam di dalam SQL:**

```
GetSequenceNumber_SQL · GetTokenStorage_SQL · InsertClaimPNC · InsertClaimRejected_Sql
InsertHistoryAkseptasiPega_Sql · InsertLOGDirectKasir_SQL · InsertLogServiceClaim
Insert_T_Storage_SQL · SaveOSClaim_SQL
```

⚠️ **Digabung dengan §2e — nol metode `Commit` Pega — artinya SELURUH penutupan transaksi modul ini
terjadi di dalam SQL.** ⛔ Apakah itu dikehendaki **tidak terbaca dari korpus**.

**Objek Oracle** *(10 teratas)*: `t_storage_image` 3 · `dual` 2 · `pooldata.t_folder_image` ·
`pooldata.kode_produksi` · `pooldata.claimrejected` · `historyakseptasipega` ·
`pooldata.directtokasir_log` · `pooldata.monitoring_klaim_log` · `pooldata.subprogressclaim` ·
`reinsurance.trloss_detail_t` — 1 masing-masing.

⚠️ **Prefiks schema tidak konsisten** — `pooldata.` ditulis pada sebagian, `historyakseptasipega`
dan `t_storage_image` **telanjang**. ⛔ Sama dengan cacat yang sudah tercatat di modul lain.

### T5d — Alias SQL

⛔ **Jendela yang dipakai memberi NOL pasangan.**

**Jendela disebut penuh:** alias **berkutip ganda**, sumber **satu kata**, pembanding memakai
**potongan terakhir sesudah titik**, **tidak peka** huruf besar-kecil, atas medan `pyBrowseSQL`
ke-17 berkas `RDBList`.

⚠️ **Sebabnya terbaca:** **9 dari 17** SQL modul ini adalah **blok `BEGIN … END`**, yang **tidak
memakai alias kolom sama sekali**. ⛔ **Bukan berarti nol alias berbohong** — berarti **jendela ini
tidak menjangkaunya**. `[terbuka]` **butir 11**

### T5e — 3 `ConnectREST`

| Rule | Autentikasi | Sumber alamat |
| --- | --- | --- |
| `KonversiKlaimNonLife` | `false` | **`SETTING`** → tabel tautan layanan |
| ⚠️ **`SendAcceptationToKasir`** | ⚠️ **`true`** | **`SETTING`** → tabel tautan layanan |
| `ServiceGoogle` | `false` | **`SETTING`** → tabel tautan layanan |

⭐⭐ **Ketiganya mengambil alamat dari tabel tautan layanan — NOL URL langsung.** ⚠️ Bandingkan
Claim Fac In, yang punya **satu dari delapan** memakai URL langsung. ✅ **Modul ini lebih bersih.**

⚠️ **Satu berautentikasi** — rule Kasir, memakai profil autentikasi bernama. ⛔ **Nilai
kredensialnya TIDAK dibaca, TIDAK dicetak, TIDAK disalin.** ⚠️ Karena kredensial itu beredar di
dalam ekspor, ia **sebaiknya diganti sesudah migrasi** — urusan tim pemilik layanan.

---

## §6 — Kontrak dengan Claim Fac In

### T6a — Kasus komite **tidak dilahirkan dari sini**

`[terverifikasi]` **Jendela: seluruh 114 berkas.**

| Yang dicari | Hasil |
| --- | ---: |
| `pxAddChildWork` | ⛔ **NOL** |
| `pxCoveredInsKeys` | ⛔ **NOL** |

⭐ **Benar dan sesuai bentuk:** modul ini **kasus ANAK**. Pembuatnya ada di Claim Fac In, yang
`[terverifikasi]` membawa **lima penunjuk posisional** dan **nol pemeriksaan** bila yang ditunjuk
bergeser.

### T6b — Yang ditulis balik ke kasus induk: **18 penugasan**

`[terverifikasi]` Penugasan bersasaran halaman induk *(`Primary.*`)*:

| Sasaran | Kali |
| --- | ---: |
| `Primary.ClaimData.Attachment(…)` — `AttachStream` · `IMAGEID` · `URLPUBLIC` · `pyCategory` · `PNOTE` · `pyMemo` | **7** |
| `Primary.PrintRISlip.FacOfferList(1).SecurityReinsurer…` | 4 |
| `Primary.ShowRetroClaim(…)` — `ReinsurerName` · `ReinsurerID` · `PctShareAllObj` | 6 |
| ⭐ `Primary.StatusKasir` | **1** |

⭐ **Modul komite menulis ke kasus klaim induk** — lampiran berkas, daftar retro, dan **status
Kasir**. ⚠️ Bandingkan Komite Claim Prop, yang mencatat **delapan penulisan lintas modul tingkat
klaim**. ⛔ **Belum saya adu satu per satu.** `[terbuka]` **butir 12**

### T6c — Tabel yang disentuh kedua modul

| Penanda | Kemunculan | Berkas |
| --- | ---: | ---: |
| ⭐ `os_akseptasi_klaim` | **24** | **4** |
| `AcceptStatus` | **90** | **10** |
| `KomiteList` | **65** | **7** |
| `IsPEGAPROD` | **46** | **10** |
| `ComiteeClaim` | 18 | 2 |
| `KomiteID` | 11 | 3 |

⭐⭐ **`os_akseptasi_klaim` disentuh kedua modul** — **24 kali** di sini, **27 kali** di Claim Fac In.
⚠️ **Itu bukti terkuat sejauh ini bahwa struktur tabel keduanya HARUS dibuat bersamaan**, dan
✅ membenarkan keputusan work owner 2026-09-19.

### T6d — ⭐ Uang: **modul ini BERSIH**

⚠️ **Sapuan kasar memberi alarm palsu.** Pencarian kata `Double|Float` **tidak peka huruf** atas
seluruh teks memberi **312 kemunculan di 35 berkas** — angka yang mengkhawatirkan.

⭐ **Dibaca dari wadah yang benar** *(`pyParameters` langsung di bawah `pagedata`)*:
⭐⭐ **NOL parameter bertipe `Double` atau `Float`.**

`[terverifikasi]` Kata itu sebenarnya hanya muncul di **7 tempat**: `pyConditionString` **4** ·
`PropertiesValue` **2** · `pyRuleName` **1**.

⭐ **Sama dengan Komite Claim Prop**, yang juga nol. ⚠️ Tetapi uang **masuk dari modul klaim**, dan
di Claim Fac In pintu masuknya **bertipe `Double`**. ⛔ **Belum diperiksa apakah jalur masuk ke sini
juga begitu.** `[terbuka]` **butir 13**

---

## §7 — Penutup

### 7.1 — ⚠️ Kesimpulan ronde ini yang **PALING RAWAN SALAH**

⚠️⚠️ **§5a — tabel gerbang layar.**

**Kenapa rawan:** kolom *"tidak bernilai mati"* disusun dari penyaring **tiga nilai**
*(`NEVER` · `1=2` · `false`)*, dan ⛔ **penyaring itu belum diuji instrumennya**. Di Claim Fac In
jendela yang lebih luas menemukan bentuk mati lain. ⛔ **Angka 519 · 52 · 11 · 2 · 3 karena itu
BUKAN "hidup"** — ia hanya *"tidak termasuk tiga bentuk mati yang saya kenal"*.

**Yang kokoh dan tidak rawan:** ⭐ **§3 sensus alur** — ia **menutup secara aritmetika** *(17 = 17)*
dan tiap entri terbaca **wadah demi wadah**; ⭐ **§4 penulis `BusinessType`** — terbaca langsung
sebagai satu penugasan, bukan disimpulkan.

**Rawan kedua:** §1d *"pola ADESAMUEL@ berulang"*. Saya menghitung **satu** berkas di ruleset itu
tetapi ⛔ **tidak mengisolasi yang mana**. Bila ternyata ia berkas sepele, kalimatnya kehilangan
bobot.

### 7.2 — ⭐ Pertanyaan untuk work owner — **tiga**

❓ **Q1** — **Kelas kerja `Work-Komite` yang tanpa penanda lini**

Ketiga modul komite memakai **tiga kelas kerja berbeda**: `Work-KomiteLife`, `Work-KomiteTreaty`,
dan — di sini — ⭐ **`Work-Komite`** tanpa penanda lini. ⛔ Korpus ini tidak memuat modul lini lain,
jadi apakah `Work-Komite` adalah **kelas komite generik yang dipakai beberapa lini** **tidak
terbaca**.

**Bedanya kalau A atau B.** **A — ia khusus Fac In**, hanya namanya yang pendek: satu tabel kasus
komite per lini, seperti dua modul lain. **B — ia generik**: satu tabel kasus komite **dipakai
beberapa lini**, dan struktur tabelnya harus memuat penanda lini.

➡️ **Rekomendasi: A — perlakukan sebagai khusus Fac In sampai terbukti sebaliknya.** ⭐ Alasannya:
ekspor ini **hanya memuat Fac In**, dan menganggapnya generik berarti merancang kolom untuk lini
yang **belum pernah kita lihat**. ⚠️ Tetapi ⛔ **ini pertanyaan untuk DBA juga** — satu pandangan ke
tabel produksi menutupnya dalam semenit.

---

❓ **Q2** — **Penutupan transaksi seluruhnya di dalam SQL**

`[terverifikasi]` Modul ini punya ⭐ **NOL metode langkah `Commit` Pega**, tetapi ⭐ **9 dari 17
`RDBList` membawa `COMMIT` tertanam di dalam SQL-nya**. ⛔ Artinya seluruh penutupan transaksi
terjadi **di dalam rule yang menurut namanya membaca**.

**Bedanya kalau A atau B.** **A — ikuti keputusan Claim Fac In:** satu aksi pengguna, satu
transaksi; kesembilan `COMMIT;` **tidak direplikasi**, kecuali pembangkit nomor urut.
**B — modul komite diperlakukan berbeda**, karena keputusan komite memang titik akhir.

➡️ **Rekomendasi: A, dengan pengecualian yang sama.** ⭐ Alasannya: keputusan K1 Claim Fac In
2026-09-19 sudah menetapkannya, dan ⚠️ **kedua modul berbagi rule** — `GetSequenceNumber_SQL`,
`SaveOSClaim_SQL`, `InsertClaimPNC`, `GetTokenStorage_SQL` **ada di kedua daftar**. ⛔ Memutuskan
berbeda berarti **satu rule dua nasib** pada rule yang identitas empat bagiannya sama.

---

❓ **Q3** — ⚠️⚠️ **Pintu belakang `"IT Developer"` di dalam rantai ATAU**

`[terverifikasi]` `SetProteksiSubmiteKomite` langkah **2.1** menguji
**`OperatorID.pyPosition == "IT Developer"`** dengan **kode arah 5**, sehingga rantai gerbangnya
**berperilaku ATAU** — syarat itu terpenuhi, **sisanya tidak diuji**.

⭐ Ini kemunculan **ketiga** pola *"wewenang digantungkan pada teks jabatan"*, dan **yang pertama
di dalam rantai ATAU**. ⚠️ Rule-nya bernama *"set proteksi submit komite"* — ⛔ **apa yang
dilindunginya belum dibaca**.

**Bedanya kalau A atau B.** **A — ia pintu belakang yang sama** dengan dua modul lain: perlakuannya
mengikuti keputusan yang sudah ada — **izin eksplisit lewat peran, bukan teks jabatan**.
**B — ia sesuatu yang lain**, dan menyamakannya akan salah.

➡️ **Rekomendasi: baca isi rule-nya lebih dulu, jangan putuskan sekarang.** ⭐ Alasannya: dua
kemunculan sebelumnya mengatur **penghapusan catatan kronologi**; yang ini menyentuh **penyerahan
ke komite** — ⛔ **wilayah yang jauh lebih berat**, dan menyamakan keduanya tanpa membaca adalah
persis kesalahan yang sudah berulang di proyek ini.

### 7.3 — Butir `[terbuka]` yang lahir ronde ini — **tiga belas**

| # | Butir | Pemilik |
| --- | --- | --- |
| **1** | Apakah `Work-Komite` kelas komite **generik** atau khusus Fac In | work owner + DBA |
| **2** | Berkas mana yang ada di ruleset `ADESAMUEL@` | asisten |
| **3** | Isolasi **4 berkas `pegaprdnusare`** — merekalah yang berlaku | asisten |
| ⚠️ **4** | Isi `SetProteksiSubmiteKomite` dan pintu belakang `"IT Developer"` | asisten → work owner |
| **5** | Berapa dari 25 langkah `Obj-*` punya jalur kegagalan | asisten |
| **6** | Apakah rule di luar ekspor dapat melempar `komiteAccept_ticket` | work owner + tim Pega |
| **7** | Apakah daur hidup kedua modul komite benar-benar sama isinya | asisten |
| **8** | Kondisi `IsKomiteLoop` — tidak ketemu di tujuh medan | asisten |
| **9** | Apakah ketujuh properti tanpa penulis diisi dari luar | asisten |
| **10** | Penyaring gerbang layar mati belum diuji instrumennya | asisten |
| **11** | Alias SQL — jendela ini tidak menjangkau blok `BEGIN…END` | asisten |
| **12** | Adu 18 penulisan ke induk dengan 8 penulisan lintas modul Komite Claim Prop | asisten |
| **13** | Apakah uang masuk ke modul ini lewat gerbang bertipe pecahan biner | asisten |

⛔ **Nol di antaranya ditutup ronde ini.**

### 7.4 — Yang paling menahan untuk ronde 2

| # | Yang dikerjakan | Kenapa menahan |
| --- | --- | --- |
| ⭐ **1** | **Baca keempat rule `KomitePost_*`** | ⭐ **jantung modul** — efek keluar keputusan komite |
| ⭐ **2** | **Baca `SetProteksiSubmiteKomite`** *(butir 4)* | menyentuh **wewenang penyerahan** |
| **3** | **Baca `KomiteRouter` dan `ApprovalKomite_Act`** | penempatan tugas dan penyimpanan keputusan |
| **4** | **Adu daur hidup dengan Komite Claim Prop** *(butir 7)* | bentuknya identik; isinya belum diuji |
| **5** | **Isolasi 4 berkas produksi** *(butir 3)* | keputusan work owner menjadikannya **yang berlaku** |
| **6** | **Kondisi `IsKomiteLoop`** *(butir 8)* | satu-satunya `When` yang kondisinya belum terbaca |

### 7.5 — ⭐ Beda dengan Komite Claim Prop — tabel berdampingan

| | Komite Claim Prop | ⭐ Komite Claim Fac In |
| --- | --- | --- |
| berkas | **80** | **114** |
| kelas kerja | `Work-KomiteTreaty` | ⭐ **`Work-Komite`** |
| rule `When` | — | ⭐ **49** *(43% modul)* |
| bentuk alur | **4** · penghubung 4 | ⭐ **4** · penghubung 4 — **identik** |
| pengenal bentuk | `END52` · `ASSIGNMENT63` · `Decision1` · `Start1` | ⭐ **sama persis** |
| tiket bernama | `komiteAccept_ticket` | ⭐ **sama persis** |
| `Obj-Set-Tickets` | **NOL** | ⭐ **NOL** |
| metode `Commit` Pega | ada | ⭐ **NOL** |
| `ConnectREST` URL langsung | — | ⭐ **NOL** — semua lewat tabel tautan |
| parameter `Double`/`Float` | **NOL** | ⭐ **NOL** |
| medan hak akses terisi | **NOL** | ⭐ **NOL** |
| `AcceptStatus` | 87 di 8 berkas | **90** di **10** berkas |

⭐⭐ **Kesimpulan bentuk: kedua modul komite adalah SATU RANCANGAN yang disalin.** ⛔ **Isinya belum
diadu** — itu ronde 2.

---

## Lampiran — bukti berkas lain tidak disentuh

```
korpus Komite Claim FacIn   114 berkas .xml   — nol dibuka untuk ditulis
korpus Claim Fac In         482 berkas .xml   — nol dibuka untuk ditulis
korpus Komite Claim Prop     80 berkas .xml   — nol dibuka untuk ditulis
docs\adr\                    15 berkas        — nol disunting
.scratch\claim-facin\         7 berkas        — SATU sisipan pada spec.md (T8)
.scratch\claim-prop\         27 berkas        — nol disunting
.scratch\komite-claim-prop\  27 berkas        — nol disunting
CLAUDE.md                                     — nol disunting
```

⛔ **Nol kode Go/React · nol `CREATE TABLE` · nol DDL · nol nomor baris XML dikutip · nol nilai
rahasia disalin · nol keputusan work owner dibuat · nol butir `[terbuka]` ditutup · nol ADR
direvisi · nol kesimpulan untuk aplikasi Go.**
