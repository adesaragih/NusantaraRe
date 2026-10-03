# Grilling Ronde 2 — Komite Claim Prop

**Tanggal:** 2026-09-18 · **Korpus:** `D:\XML\RNM_BRD\Komite Claim Prop\` (READ-ONLY)
**Lingkup:** Komite Claim Prop saja `[keputusan work owner]` 2026-09-18.
**Metode:** seluruh angka dihitung ulang sendiri dengan parser XML sungguhan. Bukti = path berkas +
nama rule + nomor step Pega. **Tidak ada nomor baris XML di berkas ini.**

> **Urutan kerja ditaati.** §A1 dan §A2 ditulis **sebelum** satu pun kolom disimpulkan. Sensus §B
> dijalankan **sesudah** aturannya berdiri.

---

## 0. Ringkasan

| # | Hasil | Status |
| --- | --- | --- |
| **§A** | Aturan baca Section / FlowAction | ✅ **BERDIRI** — 5 aturan, §A2 |
| **§B1** | Kolom layar `ShowTransfer` | ✅ **85 sel medan · 55 properti berbeda** |
| **§B2** | `ViewTransferDtl` | ⭐ **NOL kolom sendiri** — ia pembungkus, bukan layar |
| **§B3** | Properti layar tanpa penulis di modul ini | ⚠️ **15 dari 55** |
| **§B4** | Silang delapan properti ronde 1 | ✅ 5 tampil · **3 tidak** |
| **§C · F3** | Rule Connect-SQL step 16.5 / 16.7 | ✅ **TERPECAHKAN** — kelimanya bernama |
| **§D · F6** | Pewarisan kelas kerja | ⛔ **tidak ada di jendela ini** |
| **§E · K10** | Kode arah 1 / 4 / 5 | 🔴 **tetap `[terbuka]`** — angka ronde 1 terkonfirmasi |
| **RALAT** | Nama rule insert riwayat | ⚠️ ronde 1 kurang akhiran `_SQL`, §G |

---

## §A — Aturan baca Section dan FlowAction

### A1 — Struktur XML apa adanya `[terverifikasi]`

Kedua berkas berakar `<pagedata>` dan terurai penuh oleh parser XML.

```
<pySections>
  <rowdata>                        <- SATU LAYOUT
    pyIsVisibilityOption           <- ALWAYS | CONDITION | ExpressionCondition
    pyContainerVisibleWhen         <- teks syaratnya
    <pyHeaderTable>                <- judul layout
    <pySectionBody>
      <rowdata>
        pyDisplayWhen
        <pyTable>
          <pyRows>
            <rowdata>              <- SATU BARIS
              <pyCells>
                <rowdata>          <- SATU SEL  (pxObjClass = Embed-Display-Table-Cell)
                    pyType         <- FIELD | LABEL | LAYOUT | SUB_SECTION
                    pyValue        <- NAMA PROPERTI, diawali titik
                    pyReadOnly     <- true | false
                    pyEditOptions  <- Read-only | Auto | Editable
                    pyRequired     <- true | false
                    pyFormat       <- pxNumber | pxDropdown | pxTextInput | ...
                    pyLabelFieldValue / pyLabelFor   <- teks label
                    <pySections>   <- ⚠️ LAYOUT BERSARANG DI DALAM SEL
```

⚠️ **Jebakan bersarang yang sama seperti Activity, dan lebih dalam.** Sebuah sel dapat memuat
`<pySections>` lagi. Di `ShowTransfer` sarangnya mencapai **lima tingkat**
(`pySections → pyCells → pySections → pyCells → pySections → pyCells → …`). **Pembacaan urutan
baris pasti salah atribusi.** Hanya parser XML yang aman — persis pelajaran ronde 1.

⭐ **Elemen yang menyebut satu kolom layar adalah `<pyValue>`**, bukan `pyName` atau `pyReference`.
`[terverifikasi]` Sensus tag yang isinya berupa nama properti Pega di `ShowTransfer`: `pyValue`
**109**, `pyLabelFor` **41**, `pyPrompt` **7**, `pyDisabledWhen` **4**, sisanya satuan. Hanya
`pyValue` yang duduk di dalam sel `Embed-Display-Table-Cell`.

### A2 — Aturan baca `visible-when` `[terverifikasi]`

**Aturan 1 — gerbang tampil hidup di tingkat LAYOUT, bukan di tingkat sel.**
`[terverifikasi]` Dari **85** sel medan `ShowTransfer`, yang punya `pyVisible` = **NOL**,
`pyCondition` = **NOL**, `pyVisibleWhen` = **NOL**, `pyReadOnlyCondition` = **NOL**. Seluruh
gerbang tampil ada di `<rowdata>` layout di atasnya. **Sel tidak pernah bergerbang sendiri.**

**Aturan 2 — `pyIsVisibilityOption` menentukan ada-tidaknya gerbang; `pyContainerVisibleWhen` berisi
syaratnya.** Sebarannya di `ShowTransfer`:

| `pyIsVisibilityOption` | Jumlah | Arti |
| --- | --- | --- |
| **`ALWAYS`** | **32** | selalu tampil — `pyContainerVisibleWhen` diabaikan |
| **`CONDITION`** | **7** | bergerbang, syarat di `pyContainerVisibleWhen` |
| **`ExpressionCondition`** | **2** | bergerbang, syaratnya ungkapan |
| *(tidak ada elemennya)* | 495 | bukan layout — elemen lain |

⛔ **Inilah bedanya dari Activity.** Di Activity, "ada gerbang atau tidak" dibaca dari
`pyStepsPreCondition` (`true` / `false` / kosong). Di Section, dibaca dari
**`pyIsVisibilityOption`**. Nilai `ALWAYS` **bukan** berarti gerbangnya mati — ia berarti
**tidak pernah ada gerbang**.

**Aturan 3 — gerbang mati dikenali dari isi syaratnya, bukan dari saklar terpisah.**
`[terverifikasi]` Seluruh nilai `pyContainerVisibleWhen` yang terisi di `ShowTransfer`:

| Syarat | Jumlah | Pembacaan |
| --- | --- | --- |
| `.TransferType =2` | **7** | hidup — hanya jalur **adjustment** |
| `pyWorkCover.ClaimData.StsKatastrofe == 'Catastrophe' \|\| .ClaimData…` | 1 | hidup |
| `pyWorkPage.IsEditEstimation=='true'` | 1 | hidup |
| **`1=2`** | **1** | ⛔ **MATI** — tidak pernah benar |
| **`NEVER`** | **1** | ⛔ **MATI** — kata kunci Pega |
| *(kosong)* | 30 | tidak bergerbang |

⭐ **`1=2` dan `NEVER` adalah dua bentuk gerbang mati di Section.** Keduanya sepadan dengan
`pyStepsBlockName = //` di Activity: isinya masih ada di ekspor, tetapi **tidak pernah tampil**.

**Aturan 4 — sunting-atau-baca dibaca di tingkat SEL, dari dua elemen yang harus sepakat.**
`[terverifikasi]` `pyReadOnly` dan `pyEditOptions` sejalan sempurna pada 85 sel:
`pyReadOnly=true` ⇔ `pyEditOptions=Read-only` (**70**); `pyReadOnly=false` ⇔ `Auto` (**14**) atau
`Editable` (**1**). **Tidak ada satu pun yang berselisih**, jadi salah satu saja cukup dibaca.

**Aturan 5 — satu properti boleh muncul berkali-kali, dengan gerbang berbeda.**
`[terverifikasi]` **85 sel medan** memuat hanya **55 properti berbeda**. Contoh paling jelas:
`.Currency` muncul **7×**, `.Adjustment.Currency` **6×**, `.ClaimSpreaded` dan `.SharePercentage`
masing-masing **5×**, `.TreatyName` dan `.CurrencyID` masing-masing **4×**. Sebagian salinan berada
di balik gerbang `1=2` yang mati.
⛔ **Jangan menghitung kolom layar dari jumlah sel.** Yang dipakai rancangan adalah **properti
berbeda**, bukan kemunculan.

⚠️ `[terbuka]` Apa yang membedakan `CONDITION` dari `ExpressionCondition` **tidak terbaca dari
struktur** — keduanya sama-sama membawa teks syarat, dan di modul ini keduanya berisi
`.TransferType =2`. Tidak ditebak.

---

## §B — Sensus layar

### B1 — `Section/ShowTransfer.xml` `[terverifikasi]`

`pxInsName` = **`ASM-FW-GCNMFW-WORK-KOMITETREATY!SHOWTRANSFER`** · `pxObjClass` =
`Rule-HTML-Section` · kelas **`ASM-FW-GCNMFW-Work-KomiteTreaty`** — kelas kerja modul ini.

| Ukuran | Jumlah |
| --- | --- |
| Sel `Embed-Display-Table-Cell` | **348** |
| — `pyType = FIELD` | 129 |
| — `LABEL` | 149 |
| — `LAYOUT` | 35 |
| — `SUB_SECTION` | 12 |
| — kosong | 23 |
| **Sel yang menunjuk properti nyata** | **85** |
| **Properti berbeda** | **55** |
| — hanya dibaca (`pyReadOnly=true`) | **70 sel** |
| — dapat disunting | **15 sel** |
| — wajib isi (`pyRequired=true`) | **6 sel** |
| — bergerbang | **51 sel** |
| — **di balik gerbang MATI** | **8 sel** |

**Daftar 55 properti** — `n` = berapa kali muncul, `ro` = berapa di antaranya hanya-baca:

| Properti | n | ro | Gerbang |
| --- | --- | --- | --- |
| `.AcceptStatus` | 1 | 0 | — · **wajib isi** |
| `.AcceptanceStatus` | 1 | 1 | `.TransferType =2` |
| `.AcceptedDate` | 1 | 1 | `.TransferType =2` |
| `.AcceptedNo` | 1 | 1 | `.TransferType =2` |
| `.Adjustment.AdjustmentValue` | 1 | 1 | `.TransferType =2` |
| `.Adjustment.BranchOfBank` | 1 | 1 | `.TransferType =2` |
| `.Adjustment.Currency` | 6 | 6 | `.TransferType =2` |
| `.Adjustment.GrossAdjustment` | 1 | 1 | `.TransferType =2` |
| `.Adjustment.GrossValue` | 1 | 1 | `.TransferType =2` |
| `.Adjustment.IndividualRiskPercentage` | 1 | 1 | `.TransferType =2` |
| `.Adjustment.IndividualRiskRNM` | 1 | 1 | `.TransferType =2` |
| `.Adjustment.IndividualRiskType` | 1 | 1 | `.TransferType =2` |
| `.Adjustment.IndividualRiskValue` | 1 | 1 | `.TransferType =2` |
| `.Adjustment.IsPropReserved` | 1 | 0 | — |
| `.Adjustment.IsProposeClose` | 1 | 0 | — |
| `.Adjustment.NameOfBank` | 1 | 1 | `.TransferType =2` |
| `.Adjustment.NoAccount` | 1 | 1 | `.TransferType =2` |
| `.Adjustment.Payable` | 1 | 1 | `.TransferType =2` |
| `.Adjustment.PayableTo` | 1 | 1 | `.TransferType =2` |
| `.Adjustment.PersenRNM` | 1 | 1 | — |
| `.Adjustment.ProposeAdjustmentValue` | 1 | 1 | `.TransferType =2` |
| `.Adjustment.SwiftCode` | 1 | 1 | `.TransferType =2` |
| `.Adjustment.Type` | 1 | 1 | — |
| `.AdjustmentGross` | 1 | 1 | `.TransferType =2` |
| `.AdjustmentValue` | 2 | 2 | `.TransferType =2` |
| `.ClaimEstimation` | 1 | 1 | — |
| `.ClaimSpreaded` | 5 | 4 | `.TransferType =2` · ⛔ **`1=2`** |
| `.Comment` | 1 | 0 | — · **wajib isi** |
| `.ConvertValue` | 1 | 1 | — |
| `.Currency` | 7 | 6 | `.TransferType =2` · ⛔ **`1=2`** |
| `.CurrencyID` | 4 | 4 | — · **wajib isi** |
| `.DateApprove` | 1 | 1 | — |
| `.EstimationDate` | 1 | 1 | — |
| `.EstimationValue` | 1 | 1 | — |
| `.GrossAdjustment` | 1 | 1 | `.TransferType =2` |
| `.GrossEstimationPct` | 1 | 1 | — |
| `.IDKomite` | 1 | 1 | — |
| `.IDR` | 1 | 0 | `.TransferType =2` |
| `.IsSubjectivity` | 1 | 0 | — |
| `.KomiteAproval` | 1 | 1 | — |
| `.KomiteComment` | 1 | 1 | — |
| `.KomiteNo` | 1 | 0 | `.TransferType =2` |
| `.KursObjectItem` | 1 | 1 | — |
| `.KursValue` | 1 | 1 | — |
| `.ObjectName` | 1 | 1 | — · **wajib isi** |
| `.SharePercentage` | 5 | 4 | `.TransferType =2` · ⛔ **`1=2`** |
| `.SubjectivityNote` | 1 | 0 | — · **wajib isi** |
| `.TSIPerObject` | 1 | 1 | — · **wajib isi** |
| `.TreatyName` | 4 | 4 | `.TransferType =2` · ⛔ **`1=2`** |
| `.TreatyType` | 2 | 2 | ⛔ **`1=2`** |
| `.Type` | 2 | 1 | `.TransferType =2` |
| `.TypeLoss` | 1 | 1 | — |
| `.USD` | 1 | 1 | — |
| `.Value` | 2 | 1 | `.TransferType =2` |
| `.pyTemplateButton` | 2 | 0 | — *(tombol, bukan data)* |

⭐ **Temuan utama layar: `.TransferType = 2` adalah saklar besar layar ini.** Dari 55 properti,
**24** hanya tampil bila `TransferType = 2`, yaitu **jalur adjustment**. Pada jalur *reject* (3) dan
*close* (4), sebagian besar layar **tidak tampil sama sekali**. Ini cocok dengan §3.5 ronde 1:
`KomitePost` memilah tepat tiga arah 2 / 3 / 4.

⛔ **`.TreatyType` hanya hidup di balik gerbang `1=2`.** Kedua kemunculannya bergerbang mati, jadi
**properti itu tidak pernah tampil di layar mana pun**. Dicatat sebagai temuan.

**Sembilan properti yang dapat disunting penyetuju** (di luar dua tombol):
`.AcceptStatus` · `.Adjustment.IsPropReserved` · `.Adjustment.IsProposeClose` · `.ClaimSpreaded` ·
`.Comment` · `.Currency` · `.IDR` · `.IsSubjectivity` · `.KomiteNo` · `.SharePercentage` ·
`.SubjectivityNote` · `.Type` · `.Value` — **13 properti, 15 sel**.
Yang **tidak bergerbang** (selalu dapat disunting): `.AcceptStatus` · `.Adjustment.IsPropReserved` ·
`.Adjustment.IsProposeClose` · `.Comment` · `.IsSubjectivity` · `.SubjectivityNote` — **enam**.

⭐ Keenam itulah **masukan nyata penyetuju komite**: keputusan (`AcceptStatus`), catatan
(`Comment`), penanda subjectivity berikut catatannya, dan dua penanda usul (`IsPropReserved`,
`IsProposeClose`). Sisanya hanya dapat disunting di jalur adjustment.

### B2 — `FlowAction/ViewTransferDtl.xml` `[terverifikasi]`

`pxInsName` = **`ASM-FW-GCNMFW-WORK-KOMITETREATY!VIEWTRANSFERDTL`** · `pxObjClass` =
`Rule-Obj-FlowAction` · kelas **`ASM-FW-GCNMFW-Work-KomiteTreaty`**.

⭐ **Ia NOL kolom sendiri.** Hanya **2** sel `Embed-Display-Table-Cell`, keduanya
`pyType` kosong dan **tanpa `pyValue`** — wadah tata letak, bukan medan. Seluruh isi layar datang
dari Section yang dirujuknya.

`[terverifikasi]` `pxRuleReferences` menyebut **tepat lima** rule, seluruhnya berkelas
`ASM-FW-GCNMFW-Work-KomiteTreaty`:

| Rule | Jenis | Peran |
| --- | --- | --- |
| `SETKOMITELIST_ACT` | `Rule-Obj-Activity` | **pra-proses** — `pyPreProcessingActivity = SetKomiteList_Act` |
| `SHOWTRANSFER` | `Rule-HTML-Section` | **badan layar** |
| `KOMITEPOST` | `Rule-Obj-Activity` | **aksi simpan** |
| `PYCAPTION!VIEWTRANSFERDTL` | `Rule-Obj-FieldValue` | judul layar |
| `ASM-FW-GCNMFW-WORK-KOMITETREATY` | `Rule-Obj-Class` | kelas terapan |

⛔ **Kontrak layar komite, utuh dalam satu baris:**
`SetKomiteList_Act` **memuat** → `ShowTransfer` **menampilkan** → `KomitePost` **menyimpan**.
Tombol: `pySubmitLabel = Submit`, `pyCancelLabel = Cancel`, `pyShowFAButtons = true`,
`pySaveButtonVisibility = Always`.

⚠️ `[terbuka]` **BARU** — `SetKomiteList_Act` belum pernah dibaca di ronde mana pun. Ia yang
menyiapkan isi layar sebelum tampil. Perlu disisir sebelum spec ditulis.

### B3 — Properti layar tanpa penulis di Activity modul ini `[terverifikasi]`

Disisir atas **35 Activity** modul ini. **15 dari 55** properti layar tidak punya satu pun
`PropertiesName` yang menulisnya di modul ini:

| Properti | Catatan |
| --- | --- |
| `.Adjustment.GrossAdjustment` · `.Adjustment.IndividualRiskPercentage` · `.Adjustment.IndividualRiskRNM` · `.Adjustment.IndividualRiskType` · `.Adjustment.IndividualRiskValue` · `.Adjustment.IsPropReserved` · `.Adjustment.IsProposeClose` · `.Adjustment.ProposeAdjustmentValue` · `.Adjustment.SwiftCode` | **9** — seluruhnya di halaman `Adjustment`, kelas `ASM-FW-GCNMFW-Data-Adjustment` |
| `.Comment` · `.GrossAdjustment` · `.KursObjectItem` · `.ObjectName` · `.TSIPerObject` | **5** — properti tingkat kasus |
| `.pyTemplateButton` | 1 — tombol bawaan Pega, bukan data |

⚠️ **Dua di antaranya dapat disunting penyetuju dan tetap tanpa penulis di modul ini:**
`.Adjustment.IsPropReserved` dan `.Adjustment.IsProposeClose`. Keduanya **ditulis langsung dari
layar**, tanpa Activity perantara.

⚠️ `[terbuka]` Ke-14 properti data (di luar `.pyTemplateButton`) **kemungkinan diisi modul lain**.
⛔ **Tidak ditelusuri** — penyisirannya keluar dari lingkup modul ini. Perbandingannya dengan
preseden K8 (data Claim Prop yang mengalir ke komite) **tidak diputuskan di sini**.

### B4 — Silang dengan delapan properti ronde 1 `[terverifikasi]`

| Properti | Di layar? | Penulis Activity modul ini |
| --- | --- | --- |
| `KomiteComment` | ✅ `.KomiteComment`, **hanya baca** | 3 |
| `KomiteAproval` | ✅ `.KomiteAproval`, **hanya baca** | 3 |
| `DateApprove` | ✅ `.DateApprove`, **hanya baca** | 3 |
| `AcceptStatus` | ✅ `.AcceptStatus`, **dapat disunting, wajib isi** | 1 |
| `IsSubjectivity` | ✅ `.IsSubjectivity`, **dapat disunting** | 1 |
| **`KomiteCount`** | ⛔ **TIDAK** | 4 |
| **`KomiteLoop`** | ⛔ **TIDAK** | 0 |
| **`TransferType`** | ⛔ **TIDAK sebagai kolom** | 0 |

⭐ **Tiga yang tidak tampil adalah tiga penggerak mesin.** `KomiteCount` dan `KomiteLoop`
menjalankan tangga penyetuju **tanpa pernah terlihat pengguna**; `TransferType` **tidak pernah
menjadi kolom** tetapi **menggerbangi 24 dari 55 properti** (§B1). Ketiganya **keadaan dalam**, bukan
medan layar.

✅ **Menguatkan §9.3 ronde 1:** `DateApprove` memang dipakai tampilan, sedangkan `DateApproval`
**nol** di kedua berkas layar modul ini. Keputusan satu kolom `DATE_APPROVE` berdiri.

---

## §C — F3: rule Connect-SQL step 16.5 dan 16.7 ✅ **TERPECAHKAN**

⛔ **Sebab ronde 1 gagal:** saya membaca `pyStepsActivityName`, yang untuk step ini hanya berisi
metode `RDB-List` — **bukan nama rule**. Nama rulenya ada di **`<pyStepsCallParams>`**, dalam tiga
elemen yang harus dibaca bersama: **`ClassName` + `Access` + `RequestType`**.

`[terverifikasi]` Kelima step `RDB-List` di `KomitePostAdjustment`:

| Step | Gerbang | `ClassName` | `Access` | `RequestType` | Keterangan Pega |
| --- | --- | --- | --- | --- | --- |
| **16.1** | ⛔ remark `//` | `ASM-FW-GISFW-Int-policyjson` | `RNM` | `GETTanggalClosing_SQL` | *get tanggal produksi* |
| **16.4** | ⛔ remark `//` | `ASM-FW-GCNMFW-Int-V_POLIS` | `GCNM` | `GenerateNoAcceptTreaty` | *generate no akseptasi* |
| **16.5** | `false` — **gerbang dimatikan, step tetap jalan** | `ASM-FW-GISFW-Int-policyjson` | `RNM` | `GetKodeProdNonLife_SQL` | *AMBIL KODE PRODUKSI* |
| **16.7** | `true` | `ASM-FW-GISFW-Int-policyjson` | `RNM` | `GetSequenceNumber_SQL` | *generate MM.YYYY DAN SEQUENCE* |
| **33** | kosong | `ASM-FW-GISFW-Int-policyjson` | `ASM` | `InsertHistoryAkseptasiPega_Sql` | *RDB Insert ke table* |

**Identitas ketiga rule yang ada di modul ini, diverifikasi dari `pxInsName`:**

```
RDBList/GetKodeProdNonLife_SQL.xml          ASM-FW-GISFW-INT-POLICYJSON!RNM!GETKODEPRODNONLIFE_SQL
RDBList/GetSequenceNumber_SQL.xml           ASM-FW-GISFW-INT-POLICYJSON!RNM!GETSEQUENCENUMBER_SQL
RDBList/InsertHistoryAkseptasiPega_Sql.xml  ASM-FW-GISFW-INT-POLICYJSON!ASM!INSERTHISTORYAKSEPTASIPEGA_SQL
```

⭐ **Pola identitas terbaca:** `ClassName` + `!` + `Access` + `!` + `RequestType`, seluruhnya
dilipat huruf besar — **persis `pxInsName`**. Jadi step `RDB-List` **memang menyebut rulenya**, hanya
terpecah tiga elemen.

⭐ **Ketiganya adalah rantai Save-As §2a ronde 1**, dan **ketiganya dipakai di activity yang sama**:
`GETKODEPRODNONLIFE_SQL ← GETSEQUENCENUMBER_SQL ← GETTANGGALCLOSING_SQL`. Jejak Save-As bukan
kebetulan penamaan — ia **keluarga rule sequence yang saling menurunkan**.

⚠️ **Dua rule TIDAK ADA di modul ini** `[terverifikasi]`: `GETTanggalClosing_SQL` (step 16.1) dan
`GenerateNoAcceptTreaty` (step 16.4). **Keduanya persis step yang di-remark `//`.** Konsisten:
langkah yang dimatikan, rulenya pun tidak ikut diekspor. ⛔ Di mana keduanya berada **tidak
ditelusuri** — di luar modul ini.

⭐ **Step 33 mengukuhkan §4 ronde 1** — `InsertHistoryAkseptasiPega_Sql` memang yang menulis
`HISTORYAKSEPTASIPEGA`, dan jalurnya `Access = ASM`, cocok dengan `ASM!` pada `pxInsName`.

**F3 tertutup.**

---

## §D — F6: pewarisan kelas kerja ⛔ **tidak ada di jendela ini**

`[terverifikasi]` **Tidak ada satu pun berkas `Rule-Obj-Class` di ekspor modul ini.** Yang ada
hanyalah **rujukan** ke kelas dari dalam berkas layar — rujukan menyebut nama kelas, **bukan
induknya**.

Kelas yang dirujuk `ShowTransfer` (298 rujukan):

| Kelas | Rujukan |
| --- | --- |
| `ASM-FW-GCNMFW-Work-KomiteTreaty` | **126** |
| `ASM-FW-GCNMFW-Data-Adjustment` | 38 |
| `ASM-FW-GCNMFW-Data-ClaimData` | 34 |
| `ASM-FW-GCNMFW-Data-Comitee` | 16 |
| *(kosong)* | 21 |

Rujukan ber-`pxRuleObjClass = Rule-Obj-Class`: **enam** — `@BASECLASS`,
`ASM-FW-GCNMFW-WORK-CLAIMTREATY`, `ASM-FW-GCNMFW-DATA-CLAIMDATA`,
`ASM-FW-GCNMFW-WORK-KOMITETREATY`, `ASM-FW-GCNMFW-DATA-ADJUSTMENT`, `CODE-PEGA-LIST`.

⛔ **Itu daftar kelas yang DIPAKAI, bukan pohon pewarisan.** Hubungan induk–anak **tidak terbaca**,
dan **tidak saya simpulkan dari pemakaian** — persis yang diperintahkan. **F6 tetap `[terbuka]`,
jawabannya tidak ada di jendela ini.**

⚠️ Satu fakta yang bertambah: kelas **`ASM-FW-GCNMFW-Data-Comitee`** muncul 16× di layar tetapi
**belum pernah tercatat** di peta kelas §3.6 ronde 1. Dicatat sebagai temuan, **tidak ditafsirkan**.

---

## §E — K10: kode arah 1 / 4 / 5 🔴 **tetap `[terbuka]`**

`[terverifikasi]` Sensus ulang **35 Activity modul ini, 566 baris syarat** — dengan parser
bersarang:

| Kode `WhenTrue`/`WhenFalse` | Jumlah | Arti yang sudah berdiri |
| --- | --- | --- |
| `2`/`2` | **298** | tanpa gerbang |
| `2`/`3` | **205** | NORMAL |
| `3`/`2` | **26** | TERBALIK |
| `_`/`3` | 14 | — |
| `2`/`6` | 10 | keluar activity bila syarat gagal |
| `6`/`2` | 7 | — |
| **`_`/`1`** | **3** | ⚠️ **K10** |
| **`1`/`3`** | **1** | ⚠️ **K10** |
| **`_`/`4`** | **1** | ⚠️ **K10** |
| **`5`/`2`** | **1** | ⚠️ **K10** |

**Enam baris berkode 1 / 4 / 5, tiga di antaranya HIDUP** `[terverifikasi]`:

| Berkas | Step | Kode | `pyStepsPreCondition` | Hidup? | Syarat |
| --- | --- | --- | --- | --- | --- |
| `KomitePostAdjustment` | 4 | `_/1` | `false` | ⛔ tidak | `pyWorkPage.TransferType=="2"` |
| `KomitePost_Close` | 1 | `_/1` | `false` | ⛔ tidak | `TransferType=="1" && …TypeComentAnalysis` |
| `KomitePost_Reject` | 1 | `_/1` | `false` | ⛔ tidak | `TransferType=="1" && …TypeComentAnalysis` |
| **`KomitePostAdjustment`** | **12** | `1/3` | `true` | ✅ **HIDUP** | `pyWorkPage.AcceptStatus=="2"` |
| **`KomitePostAdjustment`** | **26** | `_/4` | `true` | ✅ **HIDUP** | `…AdjustmentList(…).AcceptanceStatus` |
| **`SendErrorDirectKasir`** | **2.3** | `5/2` | `true` | ✅ **HIDUP** | `IsCLMP` |

✅ **Angka ronde 1 terkonfirmasi persis** — 6 baris, 3 hidup, atribusi sama.

⛔ **Artinya tetap tidak terbukti.** Dengan hanya **6 kemunculan** di 566 baris, tidak ada pola yang
dapat diuji di dalam modul ini: kode `1` muncul pada empat syarat yang tidak sejenis, kode `4` dan
`5` masing-masing sekali. **Tidak ditebak. K10 tetap `[terbuka]`.**

---

## §G — RALAT ke ronde 1

⚠️ **Nama rule insert riwayat kurang akhiran.** Ronde 1 menulis
**`ASM!INSERTHISTORYAKSEPTASIPEGA`** di tiga tempat (§3.7, §9.7, §10.1).
`[terverifikasi]` `pxInsName` sebenarnya adalah
**`ASM-FW-GISFW-INT-POLICYJSON!ASM!INSERTHISTORYAKSEPTASIPEGA_SQL`** — **ada akhiran `_SQL`**, dan
`RequestType` pada step 33 memang `InsertHistoryAkseptasiPega_Sql`.

⛔ Berkas ronde 1 **tidak disunting dari sini** sesuai batasan prompt ini. Ralat dicatat di sini
supaya terbawa ke spec. **Tidak ada kesimpulan ronde 1 yang berubah** — tabel, kolom, dan `COMMIT`
sendiri seluruhnya tetap.

---

## §H — Frontier dan butir terbuka sesudah ronde 2

### Frontier — **1 tersisa** *(dari 4)*

| # | Cabang | Keadaan |
| --- | --- | --- |
| ~~F1~~ | Layar komite | ✅ **GUGUR** — §B |
| ~~F3~~ | Connect-SQL step 16.5 / 16.7 | ✅ **GUGUR** — §C |
| **F6** | Pewarisan kelas kerja | 🔴 **tetap** — tidak ada di jendela ini, §D |
| ~~kode 1/4/5~~ | *(bukan frontier; ia K10)* | — |

### Butir `[terbuka]` modul ini — **7**

| # | Butir | Menunggu |
| --- | --- | --- |
| 1 | **K10** — arti kode arah 1 / 4 / 5 pada tiga baris hidup | korpus |
| 2 | **`OPERATORID`** `HISTORYAKSEPTASIPEGA` tak pernah diisi jalur komite | korpus / DBA |
| 3 | **F6** — pewarisan kelas kerja + tempat deklarasi properti komite | korpus, lalu work owner |
| 4 | ⭐ **BARU** — `SetKomiteList_Act` belum pernah dibaca; ia pra-proses layar | korpus |
| 5 | ⭐ **BARU** — 14 properti layar tanpa penulis di modul ini | work owner (lingkup) |
| 6 | ⭐ **BARU** — beda `CONDITION` vs `ExpressionCondition` tidak terbaca | korpus |
| 7 | ⭐ **BARU** — `.TreatyType` hanya di balik gerbang mati `1=2`: disengaja atau sisa? | work owner |

⚠️ Butir 5 menyentuh lingkup. **Tidak diputuskan di sini** — preseden K8 (data Claim Prop yang
mengalir ke komite boleh ditelusuri) **belum tentu berlaku** untuk `ASM-FW-GCNMFW-Data-Adjustment`.
Keputusannya milik work owner.

### Temuan yang menguatkan ronde 1

- ✅ `.TransferType` memilah tepat tiga arah 2 / 3 / 4 — sekarang **terlihat di layar**: 24 properti
  bergerbang `.TransferType =2` (§B1 ↔ §3.5 ronde 1).
- ✅ `DateApprove` dipakai tampilan, `DateApproval` nol (§B4 ↔ §9.3 ronde 1).
- ✅ Step 33 menulis `HISTORYAKSEPTASIPEGA` (§C ↔ §4 ronde 1).
- ✅ Enam baris kode 1/4/5, tiga hidup (§E ↔ §9.5 ronde 1).
- ✅ Rantai Save-As sequence dipakai nyata di satu activity (§C ↔ §2a ronde 1).
