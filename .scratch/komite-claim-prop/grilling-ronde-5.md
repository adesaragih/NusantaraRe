# Grilling Ronde 5 — Komite Claim Prop

**Tanggal:** 2026-09-18 · **Korpus:** `D:\XML\RNM_BRD\Komite Claim Prop\` (READ-ONLY)
**Lingkup:** Komite Claim Prop saja `[keputusan work owner]` 2026-09-18.
**Metode:** parser XML bersarang. Bukti = path berkas + nama rule + nomor step Pega.
**Tidak ada nomor baris XML di berkas ini.**

> **Ralat dicatat DI SINI**, berkas ronde 1–4 tidak disunting.

---

## 0. Ringkasan

| # | Hasil | Status |
| --- | --- | --- |
| **§A** | Tujuh keputusan work owner | ✅ tertulis · **7 butir terbuka ditutup** |
| **§A4** | ⛔ **RALAT ronde 2 §B3** | `.Adjustment.*` **milik kasus komite**, bukan Claim Prop |
| **§B** | Efek keluar | ⭐ **8 hal keluar — SELURUHNYA sebelum `COMMIT`** |
| **§B4** | Pola kegagalan | ⚠️ **TIDAK seragam** — 2 punya penanganan, 4 tidak |
| **§C** | Pilah 117 | **73 komite** · **44 klaim** — ✅ berjumlah 117 |
| **§C4** | Properti komite yang sudah punya tempat | **4 dari 73** |
| **§D** | Baris syarat keluarga pertama | ⛔ **570**, bukan 566 — ralat ronde 2 |
| **§E2** | Keluarga elemen **KEEMPAT** | ⚠️ **ADA** — `pyActivityPrivilegeList`, terisi **34 dari 35** |

⭐ **Temuan terbesar ronde ini:** `COMMIT` berada di **langkah 41**, dan **kedelapan efek keluar
berjalan sebelum itu**. Bila penyimpanan gagal, kiriman ke Kasir, arasapas, PDF, dan email
**sudah terlanjur keluar**.

---

## §A — Tujuh keputusan work owner `[keputusan work owner]` 2026-09-18

> Ditulis apa adanya. **Nol penyisiran** di bagian ini.

### A1 — `.TreatyType` **DIBUANG**

Ia hanya ada di balik gerbang mati `1=2`, **tidak pernah tampil**. **Tidak dibuat di sistem baru.**
✅ Menutup butir terbuka *"`.TreatyType` di balik gerbang mati"*.

### A2 — Nomor klaim dan dua kotak total estimasi di balik `NEVER` **DIBUANG**

**Tidak dibuat di layar baru.** ✅ Menutup butir terbuka
*"`pyWorkCover.ClaimData.NoClaim` di balik `NEVER`"*.

### A3 — Kelima perintah lompat menggantung **DIBUANG**

**Tidak dipindahkan sama sekali.** Daftarnya, seluruhnya sudah diverifikasi ronde 3–4:

| Activity | Jalur | Sasaran | Keadaan |
| --- | --- | --- | --- |
| `KomitePostAdjustment` | transisi benar | `CHK` | tanda tak ada |
| `KomitePostAdjustment` | gerbang salah | `JMP` | tanda tak ada |
| `KomitePost_Reject` | transisi benar | `SKP` | activity **nol tanda** |
| `KomitePost_Reject` | gerbang salah | `EXT` | activity **nol tanda** |
| `KomitePost_Close` | gerbang salah | `EXT` | activity **nol tanda** |

⛔ **Tidak ditebak ke mana seharusnya melompat. Tidak ditandai cacat. Dibuang.**
✅ Menutup butir terbuka *"lima sasaran lompatan menggantung"*.

### A4 — **ATURAN AWALAN**: pemilah data layar

Berlaku di **seluruh berkas mulai sekarang**:

| Awalan | Pemilik |
| --- | --- |
| **`.Xxx`** | kasus **KOMITE** |
| **`pyWorkPage.Xxx`** | kasus **KOMITE** *(ditulis lengkap)* |
| **`pyWorkCover.Xxx`** | kasus **KLAIM induk** (`CLMP-`) |

Komite **MEMBACA** pengajuan baris `AdjustmentList` yang dikirim kasus klaim induknya.
**Data klaim di layar komite DIBACA dari Claim Prop, tidak digarap modul ini.**

> ### ⛔ RALAT ronde 2 §B3
>
> Ronde 2 menulis bahwa sembilan properti `.Adjustment.*` tanpa penulis *"kemungkinan diisi modul
> lain"* dan mempertanyakan lingkupnya. **Itu keliru.**
>
> **Awalannya kosong — hanya titik — jadi ia milik kasus KOMITE**, bukan Claim Prop.
> **Tidak ada tulis-menulis lintas modul.** Isinya sampai ke halaman komite lewat **penyalinan saat
> kasus dibuat**, bukan lewat penulisan jarak jauh.
>
> ✅ Butir terbuka *"14 properti tanpa penulis"* **TUTUP** dengan alasan ini.
> ✅ Butir terbuka *"definisi kolom layar 55 atau 117"* **TUTUP** — dipilah dengan aturan awalan,
> lihat §C.

### A5 — Jalur **REJECT** dan **CLOSE**: penyetuju memang mengisi **dua** hal

`.AcceptStatus` dan `.Comment`. Empat isian lain hanya di jalur **TRANSFER ADJUSTMENT**.
**Ditiru apa adanya.**

**Nama resmi ketiga jalur** — dipakai mulai sekarang:

| `TransferType` | Nama |
| --- | --- |
| **2** | **TRANSFER ADJUSTMENT** |
| **3** | **REJECT** |
| **4** | **CLOSE** |

✅ Menutup butir terbuka *"isian jalur tolak/tutup hanya dua"*.

### A6 — Kiriman ke Kasir gagal → email **TETAP** dikirim

Kasus tetap jalan, **nol percobaan ulang**. **DITIRU APA ADANYA.**

⚠️ **Risiko yang diterima sadar:** kegagalan **tidak terlihat** sampai ada orang yang memeriksa
belakangan. Tidak ada pemberitahuan ke siapa pun bahwa transfer gagal; yang sampai ke penerima
justru **email keberhasilan akseptasi**. Selisih antara "email terkirim" dan "uang terkirim"
**hanya ketahuan dari pemeriksaan manual**.
✅ Menutup butir terbuka *"kegagalan kirim Kasir tetap mengirim email"*.

### A7 — `StatusKonversi`: syarat resmi yang **dipertahankan**

Kiriman ke Kasir **hanya boleh jalan bila NOMOR AKSEPTASI SUDAH TERCATAT** di OS akseptasi klaim.

`[data work owner]` Buktinya — `getStatusKonversi_Act`
(`ASM-FW-GCNMFW-DATA-ADJUSTMENT!GETSTATUSKONVERSI_ACT`):

```
ParamCekKonversi.CARI1 := @replaceAll(.AcceptedNo, ".", "")
RDB-List getStatusKonversi_SQL   class ASM-FW-GCNMFW-Int-OS_AKSEPTASI_KLAIM
Param.STS_KONVERSI     := @if(HASIL1 > 0, 1, HASIL1)
```

⛔ **SQL-nya di luar modul ini — tidak ditelusuri.**
✅ Menutup butir terbuka *"arti `.StatusKonversi=="1"`"*.

### Butir yang ditutup §A — **tujuh**

`.TreatyType` · nomor klaim `NEVER` · lima lompatan menggantung · 14 properti tanpa penulis ·
definisi kolom layar 55/117 · isian jalur REJECT/CLOSE · kegagalan Kasir · arti `StatusKonversi`.

---

## §B — Efek keluar: enam berkas yang belum pernah dibuka

### B1 — Identitas dan peran `[terverifikasi]`

| Berkas | `pxInsName` | Apa yang dikerjakan | Dipanggil dari |
| --- | --- | --- | --- |
| `HitServiceToKasirKMT_Act` | `ASM-FW-GCNMFW-DATA-ADJUSTMENT!HITSERVICETOKASIRKMT_ACT` | rakit payload pembayaran, panggil REST **`SendAcceptationToKasir`**, catat log | `KomitePostAdjustment` **step 34** |
| `HTMLToPDF` | ⭐ **`@BASECLASS!HTMLTOPDF`** | ubah markup HTML jadi berkas PDF | rantai cetak akseptasi |
| `InsertDocument_Act` | `ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM!INSERTDOCUMENT_ACT` | simpan baris dokumen (`Obj-Save`), lalu panggil `InsertGoogleStorage_Act` | rantai cetak |
| `SendEmailWithAttachments` | ⭐ **`@BASECLASS!SENDEMAILWITHATTACHMENTS`** | satu langkah: `Call SendEmailNotification` | rantai email |
| `GetLinkService` | `ASM-FW-GISFW-INT-M_LINK_SERVICE!GETLINKSERVICE` | **ambil alamat layanan** dari tabel | `HitServiceToKasirKMT_Act` **step 8** · `InsertGoogleStorage_Act` step 10 |
| `InsertGoogleStorage_Act` | `ASM-FW-GISFW-INT-T_STORAGE_IMAGE!INSERTGOOGLESTORAGE_ACT` | unggah berkas ke **Google Storage** lewat REST `ServiceGoogle`, lalu catat ke `T_STORAGE_IMAGE` | `InsertDocument_Act` step 4 |

⚠️ **Dua di antaranya rule bawaan Pega, bukan buatan Nusantara Re** — `HTMLToPDF` dan
`SendEmailWithAttachments`, keduanya berkelas **`@baseclass`**. ⛔ Keduanya **tidak perlu
dipindahkan**; yang perlu dipindahkan adalah **pemakaiannya**.

### B2 · B5 — ⭐ Delapan hal yang keluar, **seluruhnya sebelum `COMMIT`** `[terverifikasi]`

Urutan nyata di `KomitePostAdjustment`:

| Step | Rule yang dipanggil | Yang keluar | Tujuan | Gerbang |
| --- | --- | --- | --- | --- |
| **17** | `SaveAcceptation_Act` | baris akseptasi | **OS akseptasi klaim** *(basis data)* | penyetuju terakhir + disetujui + bukan subjectivity |
| **21** | `SaveAcceptationTreaty_TKMT` | **PDF akseptasi** | berkas → penyimpanan dokumen | *(idem)* |
| **28** | `InsertJsonClaimTreaty_act` | JSON klaim | tabel JSON klaim | `true` |
| **29** | `KonversiKlaim_Act` | panggilan layanan | ⭐ **arasapas** *(REST)* | `true` · ⚠️ transisi gagal **HIDUP** |
| **31** | `InsertLogServiceClaim` | baris log | `POOLDATA.MONITORING_KLAIM_LOG` | `true` |
| **33** | *(RDB-List)* `InsertHistoryAkseptasiPega_Sql` | 6 kolom riwayat | `POOLDATA.HISTORYAKSEPTASIPEGA` | — |
| **34** | `HitServiceToKasirKMT_Act` | payload pembayaran | ⭐ **Kasir** *(REST `SendAcceptationToKasir`)* | `true` · ⚠️ transisi gagal **HIDUP** |
| **35** | `SendEmailKlaim_KMT` | ⭐ **email** | penerima akseptasi | `true` · label **`ENDKASIR`** |
| | | | | |
| **41** | **`Commit`** | — | — | — |

⛔ **Kedelapan efek keluar berjalan SEBELUM `COMMIT` di langkah 41. NOL yang sesudah.**

⚠️ **Akibat yang harus dicatat:** bila `Commit` di langkah 41 gagal, **uang sudah dikirim ke Kasir,
arasapas sudah dipanggil, PDF sudah dibuat, email sudah terkirim, dan tiga tabel log sudah terisi**
— sementara kasusnya sendiri **tidak tersimpan**. Ini sejalan dengan K7 (§10.1 ronde 1): sembilan
rule ber-`COMMIT` sendiri juga sudah menulis sebelum titik itu.

⭐ Step **39** (`ASMForceCaseClose`) **di-remark `//`** — penutupan kasus paksa **tidak berjalan**.

### B3 — Alamat tujuan datang dari mana `[terverifikasi]`

⛔ **Bukan tertulis keras, bukan SystemSettings** — melainkan **tabel basis data**.

`GetLinkService` melakukan `Obj-Browse` atas kelas **`ASM-FW-GISFW-Int-M_LINK_SERVICE`** ke halaman
`linkService`, lalu:

```
ResponLink.URL := linkService.pxResults(1).URL
```

⭐ **Alamat diambil dari BARIS PERTAMA hasil browse**, tanpa saringan apa pun yang terbaca.

`[terverifikasi]` Rule `SystemSettings/LinkService.xml` (`LINKSERVICE!LINKSERVICE`) **ada**, tetapi
di ekspor ini **hanya memuat `pyPurpose = LinkService`** — **nilainya tidak ikut terekspor**.
⛔ **Nilai rahasia tidak disalin ke berkas ini**; yang dicatat hanya **nama penyimpannya**:
tabel `M_LINK_SERVICE`, kolom `URL`.

⚠️ `[terbuka]` **BARU** — apakah `M_LINK_SERVICE` menyimpan **satu baris untuk semua layanan** atau
satu baris per layanan. Pengambilan `pxResults(1)` tanpa saringan menunjukkan yang pertama, tetapi
**isi tabelnya tidak ada di korpus**.

### B4 — Apa yang terjadi bila **GAGAL**: ⚠️ **polanya TIDAK seragam**

| Efek keluar | Penanganan gagal | Pola |
| --- | --- | --- |
| **Kasir** — step 34 | ⚠️ transisi **HIDUP**: `StepStatusFail` → lompat **`ENDKASIR`** (step 35 = **kirim email**) | **sesuai A6** |
| **arasapas** — step 29 | ⚠️ transisi **HIDUP**: `StepStatusFail` → lompat **`ENDSERVICE`** (step 30 = *Set Param Insert Log*) | **sesuai A6** — alur lanjut |
| **REST di dalam `HitServiceToKasirKMT_Act`** — step 10.7 dan 14.4 | transisi **HIDUP**: `StepStatusFail` → lompat **`END`** (step 16 = *set error handling*), lalu step 17 `Call SendErrorDirectKasir` **bergerbang `false` → tetap jalan** | ⭐ **BEDA** — ada **email error** tersendiri |
| **PDF** `HTMLToPDF` | ⛔ **nol** — satu transisi bergerbang `false` | **BEDA — tidak ditangani** |
| **Google Storage** `InsertGoogleStorage_Act` step 11 | ⛔ **nol kode khusus** | **BEDA — tidak ditangani** |
| **Email** `SendEmailWithAttachments` | ⛔ **nol** — satu langkah, saklar transisi `0` | **BEDA — tidak ditangani** |

⭐ **Jawaban B4: pola A6 hanya berlaku pada DUA efek keluar** — Kasir dan arasapas. **Empat sisanya
tidak punya penanganan gagal sama sekali**; bila gagal, langkah berikutnya tetap jalan **tanpa ada
yang mencatat bahwa ia gagal**.

⚠️ Satu pengecualian yang menarik: di **dalam** `HitServiceToKasirKMT_Act` ada jalur error sendiri
(`SendErrorDirectKasir`) yang **mengirim email khusus kegagalan**. Jadi kegagalan REST ke Kasir
**dilaporkan**, sementara kegagalan PDF, penyimpanan berkas, dan email **tidak**.

---

## §C — Sensus properti milik komite, dengan aturan awalan §A4

> ⛔ **Ini pencatatan, bukan rancangan.** Tidak ada usulan kolom, tidak ada penamaan tabel.

### C1 — Pilah 117 `[terverifikasi]`

| Awalan | Properti berbeda | Pemilik |
| --- | --- | --- |
| `.Xxx` | **55** | KOMITE |
| `pyWorkPage.Xxx` | **18** | KOMITE |
| `pyWorkCover.Xxx` | **44** | KLAIM induk |
| `TempOpenPage.Xxx` | 0 | — |
| | **TOTAL 117** | |

⭐ **KOMITE 73 · KLAIM 44.** Berjumlah **117** — cocok dengan angka ronde 3.

### C2 — Daftar lengkap properti milik **KOMITE**

**Bagian 1 — awalan titik, 55 properti.** *(n = kemunculan · ro = hanya-baca · A/B = famili gerbang)*

| Properti | n | ro | Penulis di modul ini | Gerbang |
| --- | --- | --- | --- | --- |
| `.AcceptStatus` | 1 | 0 | ADA | — · **wajib isi** |
| `.AcceptanceStatus` | 1 | 1 | ADA | B `.TransferType =2` |
| `.AcceptedDate` | 1 | 1 | ADA | B `.TransferType =2` |
| `.AcceptedNo` | 1 | 1 | ADA | B `.TransferType =2` |
| `.Adjustment.AdjustmentValue` | 1 | 1 | ADA | B `.TransferType =2` |
| `.Adjustment.BranchOfBank` | 1 | 1 | ADA | B `.TransferType =2` |
| `.Adjustment.Currency` | 6 | 6 | ADA | B `.TransferType =2` |
| `.Adjustment.GrossAdjustment` | 1 | 1 | **NOL** | B `.TransferType =2` |
| `.Adjustment.GrossValue` | 1 | 1 | ADA | B `.TransferType =2` |
| `.Adjustment.IndividualRiskPercentage` | 1 | 1 | **NOL** | B `.TransferType =2` |
| `.Adjustment.IndividualRiskRNM` | 1 | 1 | **NOL** | B `.TransferType =2` |
| `.Adjustment.IndividualRiskType` | 1 | 1 | **NOL** | B `.TransferType =2` |
| `.Adjustment.IndividualRiskValue` | 1 | 1 | **NOL** | B `.TransferType =2` |
| `.Adjustment.IsPropReserved` | 1 | 0 | **NOL** | A `.TransferType =2` |
| `.Adjustment.IsProposeClose` | 1 | 0 | **NOL** | A `.TransferType =2` |
| `.Adjustment.NameOfBank` | 1 | 1 | ADA | B `.TransferType =2` |
| `.Adjustment.NoAccount` | 1 | 1 | ADA | B `.TransferType =2` |
| `.Adjustment.Payable` | 1 | 1 | ADA | B `.TransferType =2` |
| `.Adjustment.PayableTo` | 1 | 1 | ADA | B `.TransferType =2` |
| `.Adjustment.PersenRNM` | 1 | 1 | ADA | — |
| `.Adjustment.ProposeAdjustmentValue` | 1 | 1 | **NOL** | B `.TransferType =2` |
| `.Adjustment.SwiftCode` | 1 | 1 | **NOL** | A `SwiftCode != ''` · B `.TransferType =2` |
| `.Adjustment.Type` | 1 | 1 | ADA | — |
| `.AdjustmentGross` | 1 | 1 | ADA | B `.TransferType =2` |
| `.AdjustmentValue` | 2 | 2 | ADA | B `.TransferType =2` |
| `.ClaimEstimation` | 1 | 1 | ADA | — |
| `.ClaimSpreaded` | 5 | 4 | ADA | B `.TransferType =2` · ⛔ B `1=2` |
| `.Comment` | 1 | 0 | **NOL** | — · **wajib isi** |
| `.ConvertValue` | 1 | 1 | ADA | — |
| `.Currency` | 7 | 6 | ADA | B `.TransferType =2` · ⛔ B `1=2` |
| `.CurrencyID` | 4 | 4 | ADA | — · **wajib isi** |
| `.DateApprove` | 1 | 1 | ADA | — |
| `.EstimationDate` | 1 | 1 | ADA | — |
| `.EstimationValue` | 1 | 1 | ADA | — |
| `.GrossAdjustment` | 1 | 1 | **NOL** | B `.TransferType =2` |
| `.GrossEstimationPct` | 1 | 1 | ADA | — |
| `.IDKomite` | 1 | 1 | ADA | — |
| `.IDR` | 1 | 0 | ADA | B `.TransferType =2` |
| `.IsSubjectivity` | 1 | 0 | ADA | A `.AcceptStatus = 1 && .TransferType =2` |
| `.KomiteAproval` | 1 | 1 | ADA | — |
| `.KomiteComment` | 1 | 1 | ADA | — |
| `.KomiteNo` | 1 | 0 | ADA | B `.TransferType =2` |
| `.KursObjectItem` | 1 | 1 | **NOL** | — |
| `.KursValue` | 1 | 1 | ADA | — |
| `.ObjectName` | 1 | 1 | **NOL** | — · **wajib isi** |
| `.SharePercentage` | 5 | 4 | ADA | B `.TransferType =2` · ⛔ B `1=2` |
| `.SubjectivityNote` | 1 | 0 | ADA | A `.IsSubjectivity = true` · **wajib isi** |
| `.TSIPerObject` | 1 | 1 | **NOL** | — · **wajib isi** |
| `.TreatyName` | 4 | 4 | ADA | B `.TransferType =2` · ⛔ B `1=2` |
| `.TreatyType` | 2 | 2 | ADA | ⛔ B `1=2` — **DIBUANG §A1** |
| `.Type` | 2 | 1 | ADA | B `.TransferType =2` |
| `.TypeLoss` | 1 | 1 | ADA | — |
| `.USD` | 1 | 1 | ADA | — |
| `.Value` | 2 | 1 | ADA | B `.TransferType =2` |
| `.pyTemplateButton` | 2 | 0 | **NOL** | — *(tombol bawaan)* |

**Bagian 2 — awalan `pyWorkPage.`, 18 properti** *(seluruhnya hanya dibaca di layar)*:

`pyWorkPage.Edit` **8×** · `pyWorkPage.Type` **3×** · `pyWorkPage.Adjustment.AcceptedNo` ·
`pyWorkPage.CLMNO` · `pyWorkPage.ClaimData.InsuredRelationship` · `pyWorkPage.ClaimData.Payable` ·
`pyWorkPage.IsEditEstimation` · `pyWorkPage.Komite` · `pyWorkPage.Komite.AdjusterFee` ·
`pyWorkPage.Komite.CircumtansesCouseOfLoss` · `pyWorkPage.Komite.Occupation` ·
`pyWorkPage.Komite.Remarks` · `pyWorkPage.Komite.Salvage` ·
`pyWorkPage.Komite.SpreadingAdjustment` · `pyWorkPage.Komite.SpreadingQuotaShare` ·
`pyWorkPage.KomiteList` · `pyWorkPage.pxCreateDateTime` · `pyWorkPage.pxCreateOpName`.

⭐ **Halaman `pyWorkPage.Komite` membawa tujuh properti yang belum pernah tercatat di ronde mana
pun** — `AdjusterFee`, `CircumtansesCouseOfLoss`, `Occupation`, `Remarks`, `Salvage`,
`SpreadingAdjustment`, `SpreadingQuotaShare`. `[terbuka]` **BARU**.

**Tanpa penulis di modul ini: 15 dari 55** — sembilan `.Adjustment.*` ditambah `.Comment`,
`.GrossAdjustment`, `.KursObjectItem`, `.ObjectName`, `.TSIPerObject`, `.pyTemplateButton`.

### C3 — Siapa mengisi `.Adjustment.*` `[terverifikasi]`

**19 properti `.Adjustment.*` di layar; 9 tanpa penulis di modul ini.** Pengisinya **dipastikan**:

```
Claim Prop\Activity\AddKomiteTreatyChild_ACT.xml
pxInsName : ASM-FW-GCNMFW-DATA-ADJUSTMENT!ADDKOMITETREATYCHILD_ACT
```

Ia memakai **`Page-Copy`** dan menulis **36 penugasan** yang menyebut `Adjustment`, di antaranya:

```
childPageKomite.Adjustment.CoverageSubscript := .IndexCoverage
childPageKomite.Adjustment.EstimationValue   := .EstimationValue
childPageKomite.Adjustment.Currency          := .Currency
childPageKomite.Adjustment.CurrencyID        := .CurrencyID
childPageKomite.IndexAdjustment              := .pxListSubscript
```

⭐ **Cocok dengan §A4:** isinya **disalin ke halaman komite saat kasus dibuat**, bukan ditulis
jarak jauh. Sesudah penyalinan, nilainya **milik kasus komite**.
⛔ **Isi modul Claim Prop tidak disisir** — hanya nama rule dan fakta penyalinan yang dipastikan.

### C4 — Silang dengan sembilan kolom yang sudah dikunci

| Kolom terkunci | Properti komite yang menampungnya |
| --- | --- |
| `ID` | *(struktural)* |
| `DATA_KOMITE_ID` | *(struktural)* |
| `KOMITE_URUT` | *(struktural)* |
| `KOMITE_OPERATORID` | `.KomiteID` — ⚠️ **tidak muncul di layar** |
| `KOMITE_JABATAN` | ✅ `.IDKomite` |
| `KOMITE_EMAIL` | `.KomiteEmail` — ⚠️ **tidak muncul di layar** |
| `KOMITE_APPROVAL` | ✅ `.KomiteAproval` |
| `KOMITE_COMMENT` | ✅ `.KomiteComment` |
| `DATE_APPROVE` | ✅ `.DateApprove` |

⭐ **Empat dari 73 properti komite sudah punya tempat.** Sisanya **69 belum** — dan itu wajar:
kesembilan kolom itu adalah **satu baris penyetuju**, sedangkan 69 sisanya adalah **isi kasus
komite dan baris adjustment**, yang tempatnya ada di tabel lain.

⛔ **Didaftarkan saja. Tidak ada usulan kolom baru** — itu keputusan work owner.

**Yang BELUM punya tempat, dikelompokkan apa adanya:**

| Kelompok | Jumlah | Contoh |
| --- | --- | --- |
| `.Adjustment.*` — baris adjustment | **19** | `AdjustmentValue` · `PayableTo` · `SwiftCode` · `IndividualRisk*` |
| Angka uang tingkat kasus | 14 | `.AdjustmentGross` · `.Value` · `.IDR` · `.USD` · `.ConvertValue` · `.KursValue` |
| Keputusan & akseptasi | 6 | `.AcceptStatus` · `.AcceptedNo` · `.AcceptedDate` · `.AcceptanceStatus` · `.KomiteNo` |
| Subjectivity | 2 | `.IsSubjectivity` · `.SubjectivityNote` |
| Estimasi | 5 | `.ClaimEstimation` · `.EstimationValue` · `.EstimationDate` · `.GrossEstimationPct` |
| Mata uang & treaty | 6 | `.Currency` · `.CurrencyID` · `.TreatyName` · `.Type` · `.TypeLoss` |
| Objek & spreading | 5 | `.ObjectName` · `.TSIPerObject` · `.ClaimSpreaded` · `.SharePercentage` · `.KursObjectItem` |
| `pyWorkPage.Komite.*` — ⭐ **baru** | 7 | `AdjusterFee` · `Salvage` · `Occupation` · `Remarks` |
| Penanda layar & bawaan Pega | 5 | `pyWorkPage.Edit` · `.pyTemplateButton` · `pxCreateOpName` |

---

## §D — Ralat angka: **570**, bukan 566

`[terverifikasi]` Dihitung **dua cara**, atas 35 Activity modul ini:

| Cara | Hasil |
| --- | --- |
| Turun pohon `<pySteps>` → `<pyStepsPreCondParams>` → `<rowdata>` | **570** |
| Jangkar `<pxObjClass>Embed-ActivityPreConditions</pxObjClass>` di seluruh dokumen | **570** |
| **Selisih antar cara** | **0** |

⛔ **Angka ronde 2 §E "566" SALAH. Yang benar 570.** Selisihnya **4**, dan sebabnya **bukan
metode** melainkan **penyaringan saya sendiri**: skrip ronde 2 membuang baris yang **kedua kodenya
kosong**. Keempat baris itu:

| Berkas | Step |
| --- | --- |
| `HTMLToPDF.xml` | 9 |
| `KomitePost_Close.xml` | 12.5 |
| `KomitePost_Reject.xml` | 12.5 |
| `PrintFileAcceptance_TKMT.xml` | 11 |

**Sebaran lengkap yang benar:**

| Kode | Jumlah | | Kode | Jumlah |
| --- | --- | --- | --- | --- |
| `2/2` | **298** | | `6/2` | 7 |
| `2/3` | **205** | | `_/_` | **4** ⭐ |
| `3/2` | 26 | | `_/1` | 3 |
| `_/3` | 14 | | `1/3` · `_/4` · `5/2` | 1 masing-masing |
| `2/6` | 10 | | | |

⭐ **Keempat baris tambahan tidak membawa kode apa pun**, jadi **tidak satu pun kesimpulan K10
berubah** — enam baris berkode 1/4/5 tetap enam, tiga tetap hidup. Yang berubah hanya **penyebut**.

---

## §E — Apa lagi yang seharusnya dikerjakan

### E1 — Urutan sesudah ronde 5

| # | Yang dikerjakan | Kenapa perlu | Besar | Menunggu |
| --- | --- | --- | --- | --- |
| **1** | **Berkas `Flow`** — daur hidup kasus komite | satu-satunya yang menentukan **urutan tahap**; belum pernah dibuka. Murah dan menutup lubang besar | sekali sisir | korpus |
| **2** | **Keluarga ketiga `pyStepsRepeatDef`** — 85 langkah berulang | perulangan menentukan **berapa kali** langkah jalan; `Exit Iteration` sudah terbukti nyata | sedang | korpus |
| **3** | **Keluarga keempat `pyActivityPrivilegeList`** — §E2 | **wewenang**: siapa boleh menjalankan activity apa. Di modul lain ini selalu satu tiket | sedang | korpus, lalu work owner |
| **4** | **`ViewDetailInterest`** — Section 115 KB + FlowAction | **layar kedua** modul ini, belum pernah dibuka | sedang | korpus |
| **5** | **`pyActionSets` tombol + `pyWorkPage.Edit`** | `.pyTemplateButton` 2× tak tertelusur; `pyWorkPage.Edit` muncul 8× sebagai penanda | sekali sisir | korpus |
| **6** | **Sisa 26 berkas** — 12 activity penghitung · 8 RDBList · 3 ReportDefinition · 2 DataTransform · 1 DecisionTable · 1 When | melengkapi peta; angka uang di layar berasal dari sini | sedang | korpus |
| **7** | **`pyWorkPage.Komite.*`** — tujuh properti baru §C2 | belum pernah tercatat; muncul di layar | sekali sisir | korpus |

⭐ **Efek keluar turun dari daftar — selesai ronde ini.**

### E2 — Keluarga elemen **KEEMPAT**: ⚠️ **ADA**

⛔ **Dicari, tidak ditunggu.** Wadah di **tingkat activity** *(bukan tingkat langkah)*, 35 Activity:

| Wadah | Muncul | **Berisi** | Sudah dibaca? |
| --- | --- | --- | --- |
| **`pyActivityPrivilegeList`** | 35 | **34** | ⚠️ **BELUM** — ⭐ keluarga keempat |
| `pyLocalParameters` | 35 | **35** | ⚠️ **BELUM** |
| `pyParameters` | 35 | **35** | ⚠️ **BELUM** |
| `pySystemParameters` | 35 | 35 | ⚠️ BELUM |
| `pyPagesAndClasses` | 35 | 33 | sebagian |
| `pxNamedPageReferences` | 33 | 33 | ⚠️ BELUM |
| `pySteps` | 35 | 35 | ✅ ronde 1–4 |
| `pxRuleReferences` | 35 | 35 | ✅ ronde 2–5 |
| `pyActivityRolesList` · `pyCustomFields` · `pxMessageSummary` | 35 | **0** | — kosong |
| `pyRuleVersionsList` · `pyRuleCompareList` · `pyToolBarSettings` · `pxWarnings` · `pyRMAction` | 35 | 35 | metadata, bukan perilaku |

⭐ **`pyActivityPrivilegeList` adalah keluarga keempat yang berperilaku** — ia menyimpan
**`pyPrivilegeClass`**, dan sebarannya menunjuk sembilan kelas berbeda:
`Work-ClaimTreaty` **9** · `Data-Adjustment` **5** · **`Work-Komite` 4** · `Work` **3** ·
`Work-KomiteTreaty` **3** · `Int-T_STORAGE_IMAGE` 2 · `@baseclass` 2 · `Data-Estimasi` 1 ·
`Int-M_LINK_SERVICE` 1.

⚠️ Perhatikan **`ASM-FW-GCNMFW-Work-Komite`** muncul **4×** — kelas yang di ronde 1 §3.6 sempat
dicabut dari peta kelas karena dikira keliru. Di sini ia muncul lagi, **sebagai kelas hak akses**.
`[terbuka]` **BARU**. ⛔ Tidak ditafsirkan.

### E3 — Kesimpulan yang **paling rawan salah** — **3**

1. ⛔ **§B4 ronde ini — "empat efek keluar tidak punya penanganan gagal"**. Rawan karena saya
   menyimpulkan ketiadaan dari **dua keluarga** (`PreCond` dan `Trans`); **keluarga ketiga dan
   keempat belum dibaca**. Penanganan gagal bisa saja ada di tempat yang belum saya buka —
   persis kesalahan yang saya buat di ronde 2 dengan `pyUserData`.
2. ⚠️ **§A2 ronde 4 — "keluarga kedua diuji SESUDAH langkah jalan"**. Masih belum teruji silang.
   Seluruh §B4 dan §B5 bersandar padanya.
3. ⚠️ **Ronde 1 §3.7 "9 rule ber-`COMMIT` sendiri"** — pencocokan teks, bukan parser. Sekarang
   **lebih penting** karena §B5 menunjukkan seluruh efek keluar mendahului `COMMIT`.

⛔ **Ditunjuk, tidak diperbaiki.**

### E4 — Pertanyaan **BARU** untuk work owner — **3**

> ⚠️ **Nol pertanyaan lama yang menggantung** — kesembilan sudah dijawab §A. Ketiga ini memang baru.

#### Pertanyaan 10 — Seluruh efek keluar mendahului penyimpanan. Ditiru apa adanya?

**Apa yang ditanyakan.** Kedelapan hal yang keluar dari sistem — kiriman ke Kasir, panggilan
arasapas, PDF akseptasi, email, dan empat penulisan tabel — **semuanya berjalan sebelum kasus
disimpan**. Penyimpanan baru terjadi di langkah terakhir. Apakah urutan itu ditiru apa adanya.

**Kenapa muncul.** Baru terbaca ronde ini. Bila penyimpanan gagal, **uang sudah dikirim, email
sudah sampai, dokumen sudah dibuat** — sementara kasusnya sendiri tidak tersimpan, dan tidak ada
rule pembatal di korpus.

**Bedanya jawaban A atau B.** Bila **ditiru**, perilakunya sama persis dan risikonya diterima
seperti K7 — konsisten dengan preseden. Bila **dibalik** (simpan dulu, baru kirim), perilakunya
**berubah dari Pega**, dan itu keputusan besar: lebih aman, tetapi bukan lagi pemindahan.

**Yang tertahan.** Urutan seluruh alur persetujuan komite, dan sikap terhadap kegagalan
penyimpanan.

#### Pertanyaan 11 — Alamat layanan diambil dari baris pertama tabel, tanpa saringan. Memang satu baris untuk semua?

**Apa yang ditanyakan.** Alamat tujuan untuk Kasir dan untuk penyimpanan berkas sama-sama diambil
dari sebuah tabel alamat layanan, dan **selalu baris pertama**, tanpa penyaring apa pun. Apakah
tabel itu memang hanya berisi satu baris, atau ada beberapa dan yang terpakai kebetulan yang
pertama.

**Kenapa muncul.** Isi tabelnya tidak ada di korpus. Pengambilan tanpa saringan bekerja selama
tabelnya satu baris, dan **diam-diam salah** begitu bertambah.

**Bedanya jawaban A atau B.** Bila **satu baris**, ia sekadar tempat menyimpan satu alamat. Bila
**beberapa baris**, maka pemilihan tanpa saringan itu **cacat yang belum ketahuan**, dan perlu
diketahui saringan apa yang seharusnya dipakai.

**Yang tertahan.** Cara sistem baru menemukan alamat layanan luar.

#### Pertanyaan 12 — Hak akses per activity: ditegakkan, atau diganti aturan peran yang baru?

**Apa yang ditanyakan.** Setiap activity menyimpan daftar hak akses — **34 dari 35** terisi, dan
kelas haknya bermacam-macam. Apakah pengaturan hak itu ikut dipindahkan apa adanya, atau wewenang
di sistem baru ditetapkan ulang dengan aturan peran sendiri.

**Kenapa muncul.** Wadah ini belum pernah dibaca dalam lima ronde. Ia menentukan **siapa boleh
menjalankan apa** — termasuk siapa boleh menyetujui.

**Bedanya jawaban A atau B.** Bila **ditegakkan apa adanya**, isinya perlu disisir dan dipetakan
satu per satu. Bila **ditetapkan ulang**, penyisiran itu **tidak perlu**, dan wewenang komite
dirancang dari aturan bisnis, bukan dari ekspor.

**Yang tertahan.** Apakah butir nomor 3 di daftar pekerjaan perlu dikerjakan sama sekali.

---

## §F — Daftar `[terbuka]` modul ini sesudah ronde 5 — **11**

| # | Butir | Menunggu |
| --- | --- | --- |
| 1 | `OPERATORID` `HISTORYAKSEPTASIPEGA` tak pernah diisi jalur komite | korpus / DBA |
| 2 | **F6** — pewarisan kelas kerja *(permintaan sudah dirumuskan §E2 ronde 3)* | work owner |
| 3 | beda `CONDITION` vs `ExpressionCondition` tidak terbaca | korpus |
| 4 | keluarga ketiga `pyStepsRepeatDef`, 85 langkah berulang | korpus |
| 5 | beda `REPEAT` vs `EMBEDDED`, dan halaman apa yang diulang | korpus |
| 6 | urutan evaluasi bila kedua keluarga terisi | korpus / work owner |
| 7 | ⭐ **BARU** — keluarga keempat `pyActivityPrivilegeList`, terisi 34/35 | work owner *(pertanyaan 12)* |
| 8 | ⭐ **BARU** — `M_LINK_SERVICE` satu baris untuk semua layanan atau tidak | work owner *(pertanyaan 11)* |
| 9 | ⭐ **BARU** — seluruh efek keluar mendahului `COMMIT` | work owner *(pertanyaan 10)* |
| 10 | ⭐ **BARU** — tujuh properti `pyWorkPage.Komite.*` belum pernah tercatat | korpus |
| 11 | ⭐ **BARU** — kelas `ASM-FW-GCNMFW-Work-Komite` muncul lagi sebagai kelas hak akses | korpus |

**Keluar dari daftar ronde 4 — tujuh butir**, ditutup §A.
⛔ **Tidak ada butir lain dinyatakan tertutup.**

---

## §G — Ralat, dikumpulkan

| # | Di mana | Yang salah | Yang benar |
| --- | --- | --- | --- |
| 1 | ronde 2 §B3 | `.Adjustment.*` *"kemungkinan diisi modul lain"*, lingkupnya dipertanyakan | **milik kasus KOMITE** — awalan kosong. Disalin saat kasus dibuat oleh `ADDKOMITETREATYCHILD_ACT` (§A4, §C3) |
| 2 | ronde 2 §E | **566** baris syarat keluarga pertama | ⛔ **570** — 4 baris tanpa kode terbuang oleh saringan skrip saya (§D) |
| 3 | ronde 3 §F1 · ronde 4 §D1 | efek keluar *"berat, belum pernah dibaca"* | **selesai ronde ini** — 6 berkas dibaca, 8 efek keluar terpetakan |
| 4 | ronde 2 §B1 · ronde 3 §B4 | `.TreatyType` dicatat sebagai temuan | **DIBUANG** `[keputusan work owner]` (§A1) |

⛔ **Tidak satu pun berkas ronde 1–4 disunting.**
