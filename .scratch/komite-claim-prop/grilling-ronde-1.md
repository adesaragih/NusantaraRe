# Grilling Ronde 1 — Komite Claim Prop

**Tanggal:** 2026-09-18 · **Korpus:** `D:\XML\RNM_BRD\Komite Claim Prop\` (READ-ONLY)
**Metode:** seluruh angka dihitung ulang sendiri dari korpus. Nomor step = **REPEATINGINDEX Pega**
(step bersarang ditulis `16.5`, `26.1`). Tidak ada nomor baris XML di berkas ini.

> ## ⚠️ RALAT 2026-09-18 — aturan identitas rule
>
> Draf pertama berkas ini memakai **Aturan 4** dari brief ronde 1 (*"identitas rule =
> `pzOriginalInstanceKey`"*). **Aturan itu DICABUT.** Aturan proyek yang lama dan benar tidak
> pernah perlu diganti. Seluruh temuan turunannya sudah diralat di tempatnya.
>
> `[terverifikasi]` **Identitas rule = class/nama dari `pxInsName`.** Nama berkas ekspor selalu
> sama dengan nama rule — **1023 dari 1024** berkas di enam modul, dan satu-satunya pengecualian
> adalah `OutstandingClaim(1).xml`, artefak ekspor ganda Windows yang `pxInsName`-nya tetap
> `OUTSTANDINGCLAIM`. Praktisnya **nol pengecualian**. `Call` juga menyebut nama itu, jadi
> penelusuran pemanggil dan penelusuran identitas memakai field yang sama.
>
> `[terverifikasi]` **`pzOriginalInstanceKey` = jejak Save-As** — rule yang di-Save As untuk
> membuat rule ini. Dipakai **hanya** untuk membaca asal-usul; bukan identitas, bukan versi, bukan
> isi. **Field pertamanya tetap dipakai** untuk mengetahui jenis rule (`RULE-OBJ-ACTIVITY`,
> `RULE-CONNECT-SQL`, dan seterusnya).
>
> Contoh yang paling jelas:
>
> ```
> GetSequenceNumber_SQL.xml
>   pxInsName             ASM-FW-GISFW-INT-POLICYJSON!RNM!GETSEQUENCENUMBER_SQL   <- identitas, cocok
>   pzOriginalInstanceKey ... RNM!GETTANGGALCLOSING_SQL #20231208T040022.143 GMT  <- leluhur Save-As
> ```

> **Alat baca.** Step diurai dengan **parser XML sungguhan**, bukan urutan baris. Anak `<pySteps>`
> berada **di dalam** elemen `rowdata` induknya, sehingga `pyStepsPreCondParams` milik induk tidak
> pernah tertukar dengan milik anak — walau di berkas ia terserialisasi **sesudah** anaknya. Struktur
> nyatanya: `<pySteps>` → `<rowdata REPEATINGINDEX=n>` (satu step) → `<pyStepsPreCondParams>` →
> `<rowdata REPEATINGINDEX=k>` (satu baris syarat).

---

## 0. Ringkasan — apa yang ronde ini hasilkan

| # | Hasil | Status |
| --- | --- | --- |
| **K1** | Arah "terbalik" pada step simpan-akseptasi / cetak-PDF / hit-kasir | ✅ **[tertutup]** — bukan cacat produksi. Hipotesis (c) benar |
| **K10** | Kode arah yang belum punya arti | 🔴 **terbuka** — 3 baris hidup di modul ini (§9.5); sensus korpus dipindah ke §11 |
| **§C2** | Penghubung dua baris `IsKomiteLoop` | ✅ **DITUTUP — `AND`** |
| **§C5** | Asal `ParamSeq.HASIL3` sebelum generate sequence | ✅ **DITUTUP** |
| **§C4** | Penangan `TransferType==1` | 🟡 **sebagian** — ditemukan tempatnya. ⚠️ "pemanggil luar terdaftar" **DICABUT** (§3.5) |
| **§A4** | Aturan identitas `pzOriginalInstanceKey` | ⛔ **DICABUT** — identitas = `pxInsName`. Temuan "20 ketidakcocokan" **nol tersisa** |
| **Silsilah** | Jejak Save-As sebagai temuan sah | ✅ **301 berkas** di 6 modul, **20** di modul ini (§2a) |
| **§D1 · §D2** | Dua ralat ke ronde sebelumnya | ✅ **terbukti**, ditulis di §6 |
| **§D3** | Ralat `IsCLMNP` | ⛔ **DIBATALKAN** — ketiganya tiga rule berbeda; kesimpulan lama benar |
| **§E5** | `ShowTransfer` · `ViewTransferDtl` | ⛔ **BELUM** — masuk frontier ronde 2 |
| **§E9** | 55 rule `When` baru vs P14 | ⛔ **BELUM** — masuk frontier ronde 2 |

---

## 1. §B — inventaris, dihitung ulang

### B1 — isi folder `[terverifikasi]`

| Subfolder | Jumlah |
| --- | --- |
| Activity | 35 |
| RDBList | 23 |
| When | 6 |
| ReportDefinition | 4 |
| ConnectREST | 3 |
| DataTransform | 2 |
| FlowAction | 2 |
| Section | 2 |
| Flow | 1 |
| DecisionTable | 1 |
| SystemSettings | 1 |
| **TOTAL XML** | **80** |

**COCOK** dengan §B1. Berkas non-XML: **1** — `Struktur_KomiteTreaty_Flow.xlsx`. Folder `Harness`
**ada dan kosong** — cocok. Nama berkas berbeda: **78**; dua nama muncul dua kali —
`InsertLogServiceClaim.xml` (Activity + RDBList) dan `ViewDetailInterest.xml` (FlowAction + Section).
**COCOK.**

⚠️ **Koreksi kecil ke §B1.** Seluruh **80** berkas punya `pzOriginalInstanceKey`. Pembacaan awal saya
sempat melaporkan 2 berkas "tanpa key"; itu **artefak jendela baca saya sendiri** (4 KB pertama),
bukan temuan korpus. Dibatalkan.

---

## 2. Nama berkas — ⚠️ TEMUAN LAMA DICABUT SELURUHNYA

> ⚠️ **RALAT.** Draf pertama bab ini berjudul *"nama berkas ≠ nama rule: 20, bukan 6"* dan
> menyimpulkan antara lain bahwa *"dua berkas sequence saling tertukar isinya"* dan bahwa
> *"`GetSequenceNumber_SQL.xml` membaca rule `GETTANGGALCLOSING_SQL`"*. **Seluruhnya DICABUT.**
> Itu **artefak dari aturan identitas yang salah**, bukan cacat korpus.

`[terverifikasi]` **Nol ketidakcocokan nama berkas.** 1023 dari 1024 berkas ber-`pxInsName` sama
dengan nama berkasnya; satu-satunya sisa adalah `OutstandingClaim(1).xml`, artefak ekspor ganda
Windows yang `pxInsName`-nya tetap `OUTSTANDINGCLAIM`. Yang tampak seperti ketidakcocokan adalah
**jejak Save-As** di `pzOriginalInstanceKey`. Lihat §2a.

---

## 2a. Silsilah Save-As — asal-usul rule

`[terverifikasi]` Sensus enam modul (`Claim Prop` · `Komite Claim Prop` · `Claim Life` ·
`Claim Non Prop` · `PremiumList Life` · `Endorsement Life`), 1024 berkas:

| Ukuran | Jumlah |
| --- | --- |
| Berkas ber-`pxInsName` | **1024** |
| `pxInsName` == nama berkas | **1023** (sisa 1 = artefak `(1)`) |
| Berkas dengan **jejak Save-As** | **301** |
| — yang **leluhurnya beda kelas** | **41** |

Sebaran per modul: `Claim Non Prop` 78 · `Claim Prop` 76 · `Claim Life` 57 · `PremiumList Life` 42 ·
`Endorsement Life` 28 · **`Komite Claim Prop` 20**.

> Angka work owner **141** dihitung atas **518** berkas; saya menyisir **1024** berkas dan mendapat
> **301**. Cakupannya berbeda, bukan hasilnya. **Tidak dikejar agar cocok.**

**Kedua puluh jejak Save-As di modul ini** — inilah, seluruhnya, yang draf pertama salah baca
sebagai "ketidakcocokan nama berkas":

| Rule (`pxInsName`) | ← Leluhur Save-As | Kelas leluhur bila beda |
| --- | --- | --- |
| `ASM!CURRENCYSTANDARD` | `ASM!UPDATEMASTERCURRENCYSTANDARD` | |
| `GCNM!SAVEDATATOOSAKSEPTASINP` | `GCNM!SAVEDATATOOSAKSEPTASI` | |
| `GETBASE64ATTACHMENT` | `LOADATTACHMENTDATA` | **`WORK-`** |
| `GETEMAILUSER_RD` | `OPERATORSBYSKILL` | |
| `GETURLGOOGLESTORAGE_ACT` | `INSERTGOOGLESTORAGE_ACT` | |
| `HITSERVICETOKASIRKMT_ACT` | `HITSERVICETOKASIR_ACT` | |
| `ISCLMNP` | `ISCLMP` | |
| `ISCLMP` | `ISCLM` | |
| `KONVERSIKLAIMNONLIFE` | `INSERTREJECTCLAIM` | **`ASM-FW-GCNMFW-WORK-PNC`** |
| `PRINTFILEACCEPTANCE_TKMT` | `PRINTFILEACCEPTANCE` | |
| `RNM!GETKODEPRODNONLIFE_SQL` | `RNM!GETSEQUENCENUMBER_SQL` | |
| `RNM!GETLINKSTORAGE_SQL` | `RNM!GETTOKENSTORAGE_SQL` | |
| `RNM!GETSEQUENCENUMBER_SQL` | `RNM!GETTANGGALCLOSING_SQL` | |
| `RNM!GETTOKENSTORAGE_SQL` | `RNM!INSERT_T_STORAGE_SQL` | |
| `RNM!INSERTLOGDIRECTKASIR_SQL` | `RNM!INSERTLOGMOP_SQL` | **`ASM-FW-GISFW-WORK`** |
| `RNM!INSERT_T_STORAGE_SQL` | `RNM!GENERATEIMAGEID_SQL` | |
| `RNM!UPDATE_T_STORAGE_SQL` | `RNM!INSERT_T_STORAGE_SQL` | |
| `SAVEACCEPTATIONTREATY_TKMT` | `SAVEACCEPTATIONTREATY_ACT` | |
| `SENDEMAILKLAIM_KMT` | `SENDEMAILKLAIM` | |
| `SERVICEGOOGLE` | `GOOGLESTORAGE_UPLOAD` | |

Beberapa membentuk **rantai tiga tingkat** `[terverifikasi]`:

```
ISCLMNP                 <-- ISCLMP                <-- ISCLM
UPDATE_T_STORAGE_SQL    <-- INSERT_T_STORAGE_SQL  <-- GENERATEIMAGEID_SQL
GETKODEPRODNONLIFE_SQL  <-- GETSEQUENCENUMBER_SQL <-- GETTANGGALCLOSING_SQL
GETLINKSTORAGE_SQL      <-- GETTOKENSTORAGE_SQL   <-- INSERT_T_STORAGE_SQL
```

**Gunanya untuk migrasi.** Rule hasil Save-As **mewarisi kolom, parameter `CARI*`, dan gerbang**
dari leluhurnya. Kolom yang tampak nyasar di rule anak sering hanya **sisa leluhurnya**. Sebut
silsilah ini setiap kali menemukan kolom tanpa penulis.

⛔ **Silsilah BUKAN bukti dua rule berperilaku sama.** Jangan menyalin kesimpulan dari leluhur ke
anak tanpa membaca anaknya.

---

## 3. §C — sensus ulang, dengan penomoran saya sendiri

### 3.1 `KomitePostAdjustment` — angka saya berbeda dari §C5

| Ukuran | §C5 | Hitungan saya | Catatan |
| --- | --- | --- | --- |
| Step | 49 | **51** | 41 tingkat atas + 10 bersarang (`16.1`–`16.9`, `26.1`) |
| Di-remark `//` | 5 (15 16 17 18 47) | **5** | `16.1 16.2 16.3 16.4 39` |
| Gerbang **tidak berlaku** | 5 | **5** | **COCOK** — `3 · 4 · 16.5 · 36 · 40` |
| Baris syarat | 51 | **57** | |
| Arah TERBALIK | 4 | **5 step** | `6 · 16 · 17 · 21 · 34` |
| Label lompatan | EXT, ENDKASIR | **3** | `EXT` · `ENDKASIR` · **`ENDSERVICE`** (baru) |

**Penomorannya berbeda, temuannya sebagian besar sama.** Penanda yang mengunci pemetaannya:
COMMIT = step **41** saya (§C5 menyebut 49) · insert history akseptasi = step **32–33** saya (§C5: 40)
· blok nomor akseptasi = **16.5–16.9** saya (§C5: 19–22). **Saya tidak mengejar angka §C agar cocok.**

✅ **Gerbang "tidak berlaku" — §C6 BENAR, lima, dan berpadanan satu-satu.** Aturan 1 membedakan
tiga keadaan, dan hanya `pyStepsPreCondition = **false**` yang berarti "gerbang tersimpan tetapi
**dimatikan**". Hitungan tepat: **false 5 · kosong-tapi-punya-syarat 16 · true 30**.

| Saya | §C6 | Akibat sebenarnya |
| --- | --- | --- |
| **3** | 3 | `Local.Operator:=1` dan `Local.AdjusmentID:=.IndexAdjustment` **selalu** di-set |
| **4** | 4 | `Obj-Open-By-Handle` kasus klaim **selalu** dibuka |
| **16.5** | 19 | RDB **"AMBIL KODE PRODUKSI"** **selalu** jalan — inilah pengisi `HASIL3` (§3.4) |
| **36** | 44 | `Obj-Save` "save pnc" **selalu** jalan |
| **40** | 48 | `KomiteCount` **SELALU** naik, termasuk sesudah penyetuju terakhir |

> ⚠️ **Ralat ke draf pertama berkas ini.** Saya sempat menulis "10, bukan 5" — itu **keliru**,
> lahir dari definisi longgar (`precond ≠ true`) yang ikut menghitung 16 step ber-`precond` **kosong**.
> Step ber-precond kosong **tidak pernah punya gerbang**; syarat tersimpannya tidak berlaku dan itu
> **bukan** temuan. §C6 benar sejak awal.

### 3.2 ⭐ K1 — **TERPECAHKAN. Tidak ada cacat produksi.**

`[terverifikasi]` `Komite Claim Prop/Activity/KomitePostAdjustment.xml`.
Sebabnya persis hipotesis **(c)** di §C7: **ketiga step itu punya DUA baris syarat, bukan satu.**

| Step saya | (≈ §C7) | Baris 1 | Baris 2 |
| --- | --- | --- | --- |
| **17** `Call SaveAcceptation_Act` — insert ke OS akseptasi | 24 | **NORMAL** `KomiteCount==Local.TotalKomite && AcceptStatus=="1"` | **TERBALIK** `pyWorkPage.IsSubjectivity==true` |
| **21** `Call SaveAcceptationTreaty_TKMT` — generate pdf | 28 | **NORMAL** `KomiteCount==Local.TotalKomite && AcceptStatus=="1"` | **TERBALIK** `pyWorkPage.IsSubjectivity==true` |
| **34** `Call HitServiceToKasirKMT_Act` — hit kasir | 42 | **TERBALIK** `pyWorkPage.IsSubjectivity==true` | **NORMAL** `KomiteCount==Local.TotalKomite && AcceptStatus=="1"` |

Kedua baris harus lolos. Arti sebenarnya ketiga step:

> **jalan hanya bila penyetuju TERAKHIR + AcceptStatus="1" + BUKAN subjectivity.**

**Tidak ada pertentangan dengan blok penomoran.** Blok itu (step **16** beserta anak `16.5`–`16.9`)
bergerbang:

```
step 16   baris 1  NORMAL    ...AdjustmentList(Local.AdjusmentID).AcceptedNo == ""
          baris 2  TERBALIK  pyWorkPage.IsSubjectivity == true
step 16.6 · 16.7 · 16.8 · 16.9   NORMAL  KomiteCount==Local.TotalKomite && AcceptStatus=="1"
```

Jadi urutannya benar: **nomor dibuat dulu (16.x), baru disimpan (17), baru dicetak (21), baru hit
kasir (34)** — seluruhnya pada penyetuju terakhir dan hanya bila bukan subjectivity.

✅ **§C7 DIBANTAH.** Kesimpulan lama lahir dari membaca **baris syarat pertama saja**. Pasangan
kontrol step **6** (TERBALIK, satu baris) vs step **7** (NORMAL, satu baris) tetap sah dan tetap
menjadi bukti aturan 3 yang paling bersih.

### 3.3 §C2 — **DITUTUP: `AND`**

`[terverifikasi]` `When/IsKomiteLoop.xml`, kelas `ASM-FW-GCNMFW-WORK-KOMITETREATY`:

```
pyLogic            = "A AND B"
pyLogicalOperator  = AND
```

Jadi: `.AcceptStatus = "1"` **AND** `.KomiteCount <= .KomiteLoop`.

**Akibat perilaku:** komite berputar **hanya selama disetujui**. Begitu seorang penyetuju menolak
(`AcceptStatus="2"`), `IsKomiteLoop` langsung salah dan tangga **berhenti saat itu juga** — tidak
menunggu sisa penyetuju. Ini konsisten dengan step **25** (`AcceptStatus=="2"` →
`KomiteCount := KomiteLoop`) dan step **26.1** (sisa penyetuju di-auto-reject).

### 3.4 §C5 — **DITUTUP: dari mana `ParamSeq.HASIL3` terisi**

`[terverifikasi]` Urutan sebenarnya di dalam blok step 16:

| Step | Gerbang | Isi |
| --- | --- | --- |
| **16.5** | `pyStepsPreCondition=false` → **gerbang MATI, selalu jalan** | `RDB-List` **"AMBIL KODE PRODUKSI"** |
| **16.6** | NORMAL, penyetuju terakhir | `ParamSeq.CARI1 := pyWorkCover.pxObjClass` · `ParamSeq.CARI2 := ParamSeq.HASIL3 + "A"` |
| **16.7** | NORMAL, penyetuju terakhir | `RDB-List` **"generate MM.YYYY DAN SEQUENCE"** |
| **16.8** | NORMAL, penyetuju terakhir | `OutputData.START_DATE := CARI2 + TempOpenPage.OfferFacIn.QuotationData.BusinessOldId + "." + HASIL1 + ".TP" + HASIL2` |
| **16.9** | NORMAL, penyetuju terakhir | `AdjustmentList(n).AcceptedNo := START_DATE` · `AcceptanceStatus := 1` · `AcceptedDate := @CurrentDateTime()` · `InputParam.CARI41 := START_DATE` |

**Jawabannya: `HASIL3` diisi oleh step 16.5 "AMBIL KODE PRODUKSI", yang gerbangnya mati sehingga
selalu jalan — dan ia berjalan SEBELUM 16.6.** Kandidat yang §C5 duga **benar**.

⚠️ `[terbuka]` Rule Connect-SQL mana persisnya yang dipanggil 16.5 dan 16.7 **belum terpastikan** —
step `RDB-List` tidak menyebut nama rule di field yang saya baca. Kandidatnya ada di 23 RDBList
modul ini. *(Draf pertama menyebut alasan lain — "dua berkas sequence saling tertukar namanya" —
**itu DICABUT**; nama berkas selalu sama dengan nama rule.)*

### 3.5 §C4 — `TransferType==1`: ditemukan tempatnya

`[terverifikasi]` `KomitePost` memilah **tiga** arah saja (`TransferType` 2 → Adjustment, 3 → Reject,
4 → Close). Nilai **1 tidak pernah dipilah di sini**. Ia muncul **di dalam** `KomitePost_Close` dan
`KomitePost_Reject`, dan **selalu berpasangan**:

```
pyWorkPage.TransferType=="1" && pyWorkPage.Komite.TypeComentAnalysis=="5"
```

Jadi `TransferType==1` bukan "arah keempat yang masuk lewat pintu lain" — ia **sub-keadaan di dalam
dua jalur yang sudah ada**, dibedakan oleh `TypeComentAnalysis`.

> ## ⚠️ RALAT 2026-09-18 — "pemanggil luar" DICABUT
>
> `[keputusan work owner]` 2026-09-18 — **"itu hanya nama yang sama. Ada nama yang sama tapi
> tetap menjalankan sesuai modulnya masing-masing."**
>
> **Tabel di bawah ini DICABUT sebagai bukti pemanggilan lintas modul.** Ia lahir dari
> **mencocokkan nama berkas**, bukan dari identitas rule. Nama yang sama di modul berbeda adalah
> **rule berbeda**, masing-masing berjalan di modulnya sendiri.
>
> `[terverifikasi]` Bukti dari `pxInsName` — identitas rule yang sah menurut aturan proyek:
>
> ```
> Komite Claim Prop/Activity/KomitePost_Reject.xml     ASM-FW-GCNMFW-WORK-KOMITETREATY ! KOMITEPOST_REJECT
> Komite Claim FacIn/Activity/ApprovalKomite_Act.xml   ASM-FW-GCNMFW-WORK-KOMITE       ! APPROVALKOMITE_ACT
> Komite Claim Life/Activity/KomitePostAdjustment.xml  berkas terpisah, modul terpisah
> ```
>
> **Nama sama, KELAS BEDA.** Tiap modul menjalankan rule-nya sendiri.
>
> ### ⭐ Pelajaran, satu baris
>
> **Kecocokan nama berkas bukan bukti rule yang sama. Identitas rule = class/nama dari `pxInsName`.**
>
> Tabel lama dibiarkan berdiri di bawah **hanya sebagai catatan nama yang bertabrakan** — bukan
> sebagai peta pemanggil. Jangan dipakai menyimpulkan apa pun tentang lintas modul.

~~`[terverifikasi]` **Pemanggil `KomitePost*` di luar modul ini:**~~ — **DICABUT**, lihat ralat di
atas. Yang berikut hanyalah **daftar nama berkas yang kebetulan sama**, per modul:

| Nama yang bertabrakan | Berkas bernama sama di modul lain (**rule berbeda**) |
| --- | --- |
| `KomitePost` | `Claim Non Prop/Activity/CreateChildKomiteCloseNP_Act` · `CreateChildKomiteCNP_Act` · `Claim Non Prop/Section/AdjustmentDetailNP` · `InputAcceptation` · seluruh keluarga `Komite Claim FacIn` |
| `KomitePostAdjustment` | `Komite Claim Life/Activity/KomitePostAdjustment` · `Komite Claim Life/FlowAction/ViewTransferDtl` · `Komite Claim Non Prop/Activity/KomitePostAdjustment` · `KomitePostAdjustmentCWP` · `Komite Claim Non Prop/FlowAction/ViewTransferDtl` |
| `KomitePost_Close` / `_Reject` | `Komite Claim FacIn/Activity/KomitePostAct` dan saudaranya |

⛔ **Kalimat lama di sini berbunyi "jalur komite dipakai bersama lintas lini klaim, bukan rule
khas Prop". DICABUT** — ia kesimpulan dari kecocokan nama, bukan dari `pxInsName`. `KomitePost*`
di modul ini berkelas `ASM-FW-GCNMFW-Work-KomiteTreaty` dan **hanya dijalankan oleh modul ini**.

### 3.6 Peta kelas kerja `[terverifikasi]`

Kelas yang benar-benar muncul di 80 berkas — **lima kelas kerja** seperti §E7 duga, ditambah kelas
integrasi/data:

| Kelas | Berkas | Peran |
| --- | --- | --- |
| `ASM-FW-GCNMFW-Work-ClaimTreaty` | 11 | kasus **klaim** (yang dibuka komite lewat `TempOpenPage`) |
| `ASM-FW-GCNMFW-Work-KomiteTreaty` | 7 Activity + 1 When + 1 Flow + 1 FlowAction + 1 Section | kasus **komite** — `pyWorkPage` |
| `ASM-FW-GCNMFW-Data-Adjustment` | 5 Activity + 1 RDB + 1 REST | **baris adjustment** |
| `ASM-FW-GCNMFW-Work` | 4 Activity + 1 RDB | induk lintas-lini |
| `ASM-FW-GCNMFW-Work-PNC` | 1 Model + 1 REST | |
| `ASM-FW-GISFW-Int-*` | 16 | integrasi (policyjson, T_STORAGE_IMAGE, currency, dll) |
| `@baseclass` / `Work-` | 5 | When umum |

> ⚠️ **RALAT.** Draf pertama mencantumkan `ASM-FW-GCNMFW-Work-Komite` sebagai **kelas kerja
> kelima** atas dasar `KomitePost_Reject`. **DICABUT.** `pxInsName`-nya
> `ASM-FW-GCNMFW-WORK-KOMITETREATY!KOMITEPOST_REJECT` — kelasnya **KomiteTreaty**, dan manifes
> `.xlsx` **BENAR**. Yang berkelas `Work-Komite` adalah **leluhur Save-As**-nya (§2a).
> Kelas kerja di modul ini ada **empat**, bukan lima.
>
> ⛔ **RALAT KEDUA 2026-09-18.** Draf ini sempat menulis bahwa *"Kesimpulan Q2 (komite digarap
> lintas-lini) TIDAK gugur"* karena **"pemanggil dari Claim Life / Claim Non Prop / Komite Claim
> FacIn"**. **Penopang itu DICABUT** — ia kecocokan nama berkas, bukan identitas rule (§3.5).
> Yang tersisa hanyalah `ClaimData.ObjectList` milik Non Prop yang **dirujuk dari dalam**
> `KomitePost_Close`/`_Reject` **milik modul ini sendiri** — itu soal **nama properti** yang
> dirujuk, bukan bukti rule dijalankan modul lain. **Lingkup garapan tetap satu modul** (§10.2).

⚠️ `[terbuka]` **Hubungan pewarisan keempat kelas kerja belum terbukti dari korpus ini** — tidak ada
rule `Rule-Obj-Class` di ekspor. Properti komite (`KomiteCount`, `KomiteLoop`, `KomiteID`,
`KomiteAproval`, `AcceptStatus`, `TransferType`, `Comment`, `IsSubjectivity`, `SubjectivityNote`,
`IndexAdjustment`) **dipakai** di `Work-KomiteTreaty` lewat `pyWorkPage`, tetapi **tempat
deklarasinya tidak ada di ekspor ini**. Frontier ronde 2.

### 3.7 Tabel yang disentuh + siapa COMMIT `[terverifikasi]`

23 RDBList disisir. **9 memuat `COMMIT` sendiri:**

| Rule (nama sebenarnya) | Tabel | COMMIT sendiri |
| --- | --- | --- |
| `GCNM!INSERTCLAIMREJECTED_SQL` | `POOLDATA.CLAIMREJECTED` | **YA** |
| `ASM!INSERTHISTORYAKSEPTASIPEGA` | `HISTORYAKSEPTASIPEGA` ⚠️ **tanpa awalan skema** | **YA** |
| `RNM!INSERTLOGMOP_SQL` | `POOLDATA.DIRECTTOKASIR_LOG` | **YA** |
| `RNM!INSERTLOGSERVICECLAIM` | `POOLDATA.MONITORING_KLAIM_LOG` | **YA** |
| `RNM!GENERATEIMAGEID_SQL` | `T_STORAGE_IMAGE` ⚠️ tanpa skema | **YA** |
| `RNM!GETTANGGALCLOSING_SQL` | (blok PL/SQL) | **YA** |
| `RNM!INSERT_T_STORAGE_SQL` | `T_STORAGE_IMAGE` | **YA** |
| `GCNM!INSERTCLAIMPNC` | (blok PL/SQL) | **YA** |
| `GCNM!SAVEDATATOOSAKSEPTASI` · `ASM!SAVEOSCLAIM_SQL` | OS akseptasi | **YA** |

Baca saja (tanpa COMMIT): `KODE_PRODUKSI` · `PROPORTIONALARRG` · `TREATYREINSURER` · `BUSINESS` ·
`TREATYBUSINESS` · `TREATYYEAR` · `T_FOLDER_IMAGE` · `REINSURANCE.TRLOSS_DETAIL_T` ·
`CURRENCYSTANDARD`.

⚠️ Ini **memperluas K7**: bukan dua tabel yang COMMIT sendiri di luar COMMIT step 41, melainkan
**sembilan rule**. Bila langkah sesudahnya gagal, seluruhnya tetap terisi.

---

## 4. §C8 — dua tabel, dikonfirmasi

`[terverifikasi]` `RDBList/InsertClaimRejected_Sql.xml` (`ASM-FW-GCNMFW-Int-CLAIMREJECTED`) →
`INSERT INTO POOLDATA.CLAIMREJECTED` **12 kolom**, `COMMIT` sendiri. **COCOK §C8.**

`[terverifikasi]` `RDBList/InsertHistoryAkseptasiPega_Sql.xml` (`ASM-FW-GISFW-Int-policyjson`) →
`INSERT INTO HISTORYAKSEPTASIPEGA` **6 kolom**, **tanpa awalan skema**, `COMMIT` sendiri.
**COCOK §C8.** Diisi dari step **32** (param) → **33** (RDB) dalam penomoran saya.

---

## 5. §C3 · §C9 · §C10 — dikonfirmasi

`[terverifikasi]` **`KomiteRouter`**: 8 step, **6 di-remark** (`1 2 3 4 5 7`). Yang hidup: step **6**
(`.KomiteAproval==0` → `param.AssignTo := .KomiteID`) dan step 8. **COCOK §C3** — tangga workbasket
`komitepnc1..4` memang mati.

`[terverifikasi]` **`SetKomiteList_Act`**: 12 step (§C9 menyebut 9 — beda karena step bersarang).
Nol remark. **COCOK** pada intinya: activity ini menghitung total adjustment per mata uang, **bukan**
menyusun daftar anggota komite.

`[terverifikasi]` **`KomitePost_Close`** 38 step · **`KomitePost_Reject`** 36 step · **nol remark** di
keduanya · **nol arah terbalik** di keduanya. **COCOK §C10** (angka step beda karena bersarang).
⚠️ **RALAT** — draf pertama menulis *"`KomitePost_Reject` berkelas `ASM-FW-GCNMFW-WORK-KOMITE`,
kelas kerja kelima"*. **DICABUT.** `pxInsName`-nya `ASM-FW-GCNMFW-WORK-KOMITETREATY!KOMITEPOST_REJECT`;
`WORK-KOMITE` adalah **leluhur Save-As**. Manifes `.xlsx` benar.

⚠️ **Juga dicabut** (temuan brief §B4, tidak pernah saya tulis sebagai kesimpulan sendiri):
*"dua versi `HitServiceToKasir_Act`, 2023 lawan 2022"*. `[terverifikasi]`
`HITSERVICETOKASIRKMT_ACT` dan `HITSERVICETOKASIR_ACT` adalah **DUA RULE BERBEDA**, bukan dua
versi — yang kedua adalah **leluhur Save-As** yang pertama. Keduanya **wajib dibaca terpisah**;
temuan lama tentang `HITSERVICETOKASIR_ACT` tetap berlaku untuk rule itu saja.

---

## 6. §D — tiga ralat, ketiganya TERBUKTI

### D1 — `.AcceptanceStatus` **punya** penulis ⚠️ RALAT

`[terverifikasi]` `Komite Claim Prop/Activity/KomitePostAdjustment.xml`:

- step **16.9** → `TempOpenPage.ClaimData.AdjustmentList(Local.AdjusmentID).AcceptanceStatus := 1`
- step **26** → gerbangnya menguji `AdjustmentList(...).AcceptanceStatus=="2"`, dan jalur tolak
  menuliskannya

Pernyataan ronde 2 Claim Prop **P3 · "nol penulisan di seluruh Activity keempat modul"** adalah
**KELIRU**. Penulisnya ada **di dalam korpus**, di modul yang belum disensus saat itu.

### D2 — `ComiteeClaim` adalah page-**LIST** ⚠️ RALAT

`[terverifikasi]` step **6** menulis `…AdjustmentList(Local.AdjusmentID).ComiteeClaim(local.count)`
dan step **7** menulis `…ComiteeClaim(<LAST>)`, dengan `local.count` berasal dari
`pyWorkPage.KomiteCount`. **Satu baris `ComiteeClaim` = satu penyetuju**, bukan satu kasus komite.

➡️ Yang berbanding lurus dengan `ComiteeClaim` adalah **`T_KOMITE_KOMITELIST`** (satu baris per
penandatangan), **bukan** `T_GENERAL_KOMITE`. `T_GENERAL_KOMITE` berbanding lurus dengan **kasus
komite** (`pyWorkPage`, kelas `Work-KomiteTreaty`).

⛔ **Tidak satu pun berkas rancangan tabel disentuh di ronde ini.** Keputusan diserahkan ke work owner.

### D3 — ⛔ DIBATALKAN. Tidak ada yang perlu diralat soal `IsCLMNP`.

> ⚠️ **Draf pertama berkas ini menyatakan "rule `ISCLMNP` tidak ada di korpus" dan menyuruh
> meralat setiap kalimat lama yang berbunyi "digerbangi IsCLMNP". ITU SALAH — DICABUT.**
> Kekeliruan itu lahir dari Aturan 4 yang salah (lihat ralat di kepala berkas). Perintah §D3 di
> brief ronde 1 juga ditarik oleh work owner.
>
> **Tidak ada berkas ronde lama yang disunting** atas dasar D3 — ralat itu hanya sempat tertulis
> di berkas ini dan kini sudah dikembalikan. Tidak ada yang perlu dipulihkan di tempat lain.

`[terverifikasi]` **Ketiganya tiga rule BERBEDA**, masing-masing menguji prefiks namanya sendiri,
dengan logika **`A OR B`**:

| Berkas | `pxInsName` (identitas) | Syarat |
| --- | --- | --- |
| `When/IsCLM.xml` | `@BASECLASS!ISCLM` | `pyWorkCover.pyWorkIDPrefix = "CLM-"` **OR** `pyWorkPage.pyWorkIDPrefix = "CLM-"` |
| `When/IsCLMP.xml` | `@BASECLASS!ISCLMP` | `pyWorkCover.pyWorkIDPrefix = "CLMP-"` **OR** `pyWorkPage.pyWorkIDPrefix = "CLMP-"` |
| `When/IsCLMNP.xml` | `@BASECLASS!ISCLMNP` | `pyWorkCover.pyWorkIDPrefix = "CLMNP-"` **OR** `pyWorkPage.pyWorkIDPrefix = "CLMNP-"` |

`[terverifikasi]` `pyLogic = "A OR B"` pada ketiganya. **Setiap kesimpulan lama "digerbangi
`IsCLMNP`" sudah benar sejak awal.** Yang berjejak Save-As hanya asal-usulnya
(`ISCLMNP` ← `ISCLMP` ← `ISCLM`, §2a) — dan silsilah **bukan** identitas.

### D3a — temuan baru: gerbang prefiks memeriksa INDUK atau DIRI

⚠️ `[terverifikasi]` Ketiga rule memakai **`pyWorkCover` ATAU `pyWorkPage`**. Jadi gerbang itu
bernilai **benar juga ketika yang berprefiks adalah kasus INDUKNYA**, bukan kasus yang sedang
diproses.

Ini **langsung relevan untuk komite**: kasus komite (`Work-KomiteTreaty`) induknya adalah kasus
**klaim**. Sebuah `IsCLMP` yang dievaluasi di dalam kasus komite akan bernilai **benar** karena
`pyWorkCover` — si kasus klaim — berprefiks `CLMP-`, walaupun kasus komite itu sendiri tidak.
Jangan membaca gerbang prefiks sebagai "lini kasus ini"; ia berarti **"lini kasus ini ATAU
induknya"**.

---

## 7. Frontier ronde 2 — 4 cabang (lingkup Komite Claim Prop)

| # | Cabang | Prasyarat | Pemilik |
| --- | --- | --- | --- |
| F1 | `Section/ShowTransfer.xml` (2,3 MB) + `FlowAction/ViewTransferDtl.xml` — **belum dibaca sama sekali**. Tetapkan dulu aturan baca `visible-when` (beda mekanisme dari activity), tulis aturannya, baru simpulkan. Dari sini daftar kolom layar komite didapat | — | korpus |
| ~~F2~~ | ~~Nama rule sebenarnya untuk `Call ..._TKMT` / `..._KMT`~~ | — | ✅ **GUGUR** — `Call` menyebut nama rule, dan nama rule = nama berkas (`pxInsName`) |
| F3 | Rule Connect-SQL mana yang dipanggil step **16.5** dan **16.7** — step `RDB-List` tidak menyebutnya di field yang terbaca | — | korpus |
| ~~F4~~ | ~~K2 — akibat `KomiteCount` yang selalu naik~~ | — | ✅ **GUGUR** — urutan step diturunkan dengan parser bersarang, §9.1 |
| ~~F5~~ | ~~K4 — `.KomiteLoop` diisi dari mana~~ | — | ✅ **GUGUR** — §9.2 |
| F6 | Pewarisan lima kelas kerja + tempat deklarasi properti komite — tidak ada `Rule-Obj-Class` di ekspor | — | korpus, lalu work owner |
| ~~F7~~ | ~~§E9 — 55 rule `When` baru vs P14~~ | — | ↳ **KELUAR** — pokoknya di modul **Claim Prop**, §11 butir 6 |

**Frontier Komite Claim Prop → ronde 2: 4 cabang terbuka** — **F1** · **F3** · **F6** · **arti kode
1 / 4 / 5 pada ketiga baris hidup di modul ini**. **11 tertutup atau gugur**: K1 · K2 · K3 · K4 ·
K5 · K7 · K11 · §C2 · §C5 · §C4 · F2. **1 keluar ke modul lain**: F7 (§11).

⛔ Ketiga cabang yang tersisa seluruhnya berkas **`Komite Claim Prop`** — `Section/ShowTransfer` dan
`FlowAction/ViewTransferDtl` ada di folder modul ini (§1 B1: Section 2, FlowAction 2), 23 RDBList
kandidat F3 juga, dan kelima kelas kerja F6 memuat kelas komite Prop.

---

## 8. Pemblokir — keadaan sesudah ronde 1

| # | Pertanyaan | Keadaan |
| --- | --- | --- |
| **K1** | Arah step simpan/cetak/kasir | ✅ **[tertutup]** — lihat kalimat penutup di bawah tabel |
| **K2** | `KomiteCount` selalu naik | ✅ **[tertutup]** — gerbang mati **tidak merusak**; justru ia yang menghentikan loop. Urutan step juga terbukti aman. §9.1 |
| **K3** | `DateApproval` vs `DateApprove` | ✅ **[tertutup]** — satu kolom saja, `DATE_APPROVE`. §9.3 |
| **K4** | `.KomiteLoop` diisi dari mana | ✅ **[tertutup]** — jumlah baris `KomiteList`, ditetapkan saat kasus dibuat. §9.2. Pengecualian lini Fac In ↳ §11 butir 2 |
| **K5** | `.KomiteID` — operator, jabatan, atau workbasket | ✅ **[tertutup]** — **operator**. `.IDKomite` = **jabatan**. §9.4 |
| **K6** | `HISTORYAKSEPTASIPEGA` tanpa awalan skema | ✅ **[tertutup]** — `[data DBA]` 2026-09-18: **dua tabel**, `POOLDATA.HISTORYAKSEPTASIPEGA` (tablespace `POOLMASTER`) dan `POOLDATA.T_STORAGE_IMAGE` (tablespace `TBS_POOLDATA`). Keduanya terselesaikan lewat **skema bawaan koneksi, `POOLDATA`**. Tidak ada kejanggalan. §12 |
| ⚠️ **BARU** | `HISTORYAKSEPTASIPEGA.OPERATORID` tidak pernah diisi jalur komite | 🔴 **[terbuka]** — adakah penulis lain di luar modul ini, atau kolom itu memang selalu kosong. §12.1 |
| **K7** | COMMIT sendiri di luar COMMIT akhir | ✅ **[tertutup]** — `[keputusan work owner]` 2026-09-18: **ditiru apa adanya**, kesembilan rule tetap COMMIT sendiri. §10.1 |
| **K8** | `.KursIDR` | ✅ **[tertutup]** — kurs **dibekukan**. Pengisinya rule Connect-SQL `ASM!CURRENCYSTANDARD`; masukannya **mata uang + `SYSDATE`**, jadi kurs terkunci = **kurs tanggal lookup dijalankan**. `[keputusan work owner]` 2026-09-18: **"ikuti aja query di xml CurrencyStandard."** §9.6 |
| K9 | `Workbasket` selalu literal `"KLAIM"` | ↳ **KELUAR dari daftar modul ini** — sisi **TULIS** selesai dan benar (§9.7); pembacanya `GETFLAGREJECT_SQL` milik **Fac In / Treaty In**, §11 butir 1 |
| **K10** | Kode arah yang belum punya arti | 🔴 **terbuka — lingkup modul ini**: 2 baris hidup berkode 1/4 di `KomitePostAdjustment` + 1 baris kode 5 di `SendErrorDirectKasir`. §9.5. Sensus korpus ↳ §11 butir 3 |
| **K11** | Lingkup garapan komite | ✅ **[tertutup]** — `[keputusan work owner]` 2026-09-18: **hanya modul ini**. §10.2. ⚠️ Salah satu alasan pendukungnya ("pemanggil lintas lini") **DICABUT** — kecocokan nama berkas, bukan `pxInsName` (§3.5). **Keputusannya tidak berubah.** |

### K1 — kalimat penutup `[tertutup]`

`[terverifikasi]` `KomitePostAdjustment` step **24 · 28 · 42** (penomoran brief; **17 · 21 · 34**
dalam penomoran `REPEATINGINDEX` saya) masing-masing punya **dua baris syarat**:
`KomiteCount==TotalKomite && AcceptStatus=="1"` (**NORMAL**) dan `IsSubjectivity==true`
(**TERBALIK**). **Keduanya harus lolos.** Artinya: **penyetuju terakhir, disetujui, bukan
subjectivity.** Tidak ada pertentangan dengan penomoran akseptasi. Dugaan "cacat produksi" di
brief lahir dari membaca **baris pertama saja**. **[tertutup]**

### K10 — kode arah yang belum punya arti (MEMBESAR)

`[terverifikasi]` Sebaran kode arah atas **57 baris syarat** `KomitePostAdjustment`:

| `WhenTrue`/`WhenFalse` | Jumlah | Arti |
| --- | --- | --- |
| `2`/`3` | **27** | NORMAL — jalan bila syarat terpenuhi |
| `2`/`2` | **16** | seluruhnya di step **tanpa gerbang** (`precond` kosong, `when` kosong) → "tidak ada gerbang" |
| `3`/`2` | **5** | TERBALIK |
| `(kosong)`/`3` | **4** | ⚠️ **K10 — belum punya arti** |
| `6`/`2` | 1 | ⚠️ K10 |
| `2`/`6` | 1 | keluar activity bila syarat tidak terpenuhi |
| `(kosong)`/`1` | 1 | ⚠️ **K10 — kode `1` belum punya arti** |
| `1`/`3` | 1 | ⚠️ K10 |
| `(kosong)`/`4` | 1 | ⚠️ **K10 — kode `4` belum punya arti** |

⚠️ **Pola `WhenTrue` kosong bukan keanehan langka** — 4× di activity ini saja, ditambah
`KomiteRouter`. **Salah satunya adalah step yang menulis `AcceptedNo` dan `AcceptanceStatus := 1`**
(step `16.9` saya). Kode `1`, `4`, dan `6` juga belum punya arti — ketiganya masuk K10.

> Selisih hitungan dengan work owner: sebaran work owner menyebut **54** baris, saya **57**.
> Selisih 3 belum saya telusuri — tidak dikejar agar cocok.

### Rekonsiliasi — 4 lawan 5 baris TERBALIK

**Definisi saya sudah ketat:** `WhenTrue == "3"` **dan** `WhenFalse == "2"`. Dengan definisi itu
saya dapat **5 baris di 5 step**, seluruhnya ber-`precond=true` (jadi gerbangnya benar-benar
berlaku):

| Step (saya) | ≈ brief | Syarat |
| --- | --- | --- |
| **6** | 6 | `...AdjustmentList(Local.AdjusmentID).IsSubjectivity == true` |
| **17** | 26 | `pyWorkPage.IsSubjectivity==true` |
| **21** | 30 | `pyWorkPage.IsSubjectivity==true` |
| **34** | 44 | `pyWorkPage.IsSubjectivity==true` |
| **16** | — | `pyWorkPage.IsSubjectivity==true` ← **INILAH YANG KELIMA** |

⭐ **Yang kelima adalah step `16` — INDUK blok penomoran akseptasi** (`16.1`–`16.9`).
Ia terlewat pada pembacaan urutan-baris karena **`pyStepsPreCondParams` milik induk terserialisasi
SESUDAH `</pySteps>` anak-anaknya** — jebakan nesting yang sama yang pernah membakar Claim Prop.
Parser XML menangkapnya karena anak berada **di dalam** elemen induk.

**Akibatnya penting:** seluruh blok penomoran akseptasi bergerbang **bukan-subjectivity** di
tingkat induk, lalu **penyetuju-terakhir-dan-disetujui** di tingkat anak. Justru inilah yang
membuat cerita K1 utuh.

---

## 9. Sensus korpus-PENUH — tujuh pemblokir

> **Jendela saya berbeda dari jendela work owner.** Work owner memakai **555 berkas ter-stage**;
> saya menyisir **seluruh korpus, 9.429 berkas XML di 20 modul**. Setiap angka di bawah dihitung
> ulang sendiri. Di mana berbeda, **saya tulis apa adanya dan tidak saya kejar agar cocok** —
> selisihnya hampir selalu cakupan, tetapi **dua di antaranya membalik kesimpulan**.

### 9.1 K2 — `KomiteCount` yang selalu naik: **[tertutup]**

`[terverifikasi]` Penulis `.KomiteCount` korpus-penuh: **24 penulisan di 19 berkas**
(jendela 555: 10/7). Polanya seragam di **keempat lini**: `:= 1` saat kasus komite **dibuat**
(8 berkas `Claim*`), lalu `:= .KomiteCount + 1` dan `:= .KomiteLoop` saat **putaran**.

✅ **`KomitePostAdjustment` punya TEPAT DUA penulisan**, bukan empat — koreksi work owner
**terbukti**. Angka empat memang artefak pemisahan datar.

**Jawaban K2:** gerbang mati pada step penambah **tidak merusak apa pun — justru ia yang
menghentikan loop.** `IsKomiteLoop = .AcceptStatus = "1" AND .KomiteCount <= .KomiteLoop`:

```
penyetuju terakhir setuju -> KomiteCount naik jadi KomiteLoop+1 -> syarat kedua gagal -> keluar
ditolak                   -> KomiteCount dipaksa = KomiteLoop, lalu naik -> keluar juga
                             (syarat pertama sudah gagal duluan karena AcceptStatus != "1")
```

Operator `<=` memastikan penyetuju terakhir tetap kebagian giliran. Gerbang itu memang **mubazir
sejak awal**.

#### ⭐ Urutan step — diturunkan sendiri dengan parser bersarang, **aman**

`[terverifikasi]` Urutan dokumen di `KomitePostAdjustment` (penomoran `REPEATINGINDEX`):

| Step | Gerbang | Isi |
| --- | --- | --- |
| **5** | kosong | `local.count := pyWorkPage.KomiteCount` — **potret diambil di awal** |
| **6** | true | tulis `KomiteList(local.count).KomiteComment` / `.KomiteAproval` / `.DateApprove` |
| **8** | true | **baca** `"Accepted by " + KomiteList(pyWorkPage.KomiteCount).IDKomite` |
| **9** | true | **baca** `"Rejected by " + KomiteList(pyWorkPage.KomiteCount).IDKomite` |
| **25** | true | `KomiteCount := KomiteLoop` (jalur tolak) |
| **26.1** | true | auto-reject sisa penyetuju |
| **40** | **false** | `KomiteCount := KomiteCount + 1` ← **PALING AKHIR** |

✅ **Pembacaan riwayat (step 8/9) terjadi JAUH SEBELUM penambahan (step 40).** Tidak ada
pembacaan indeks di luar daftar. Butir yang work owner sisakan `[terbuka]` **kini tertutup**.
Bonus: step 25 juga berada **sesudah** 8/9, sehingga teks riwayat pada jalur tolak masih memakai
nilai sebelum dipaksa — perilaku yang benar.

### 9.2 K4 — asal `.KomiteLoop`: **[tertutup]**

`[terverifikasi]` Penulis korpus-penuh: **13 penulisan di 9 berkas** (jendela 555: 3). Pola:
`:= @Utilities.SizeOfPropertyList(childPageKomite.KomiteList)` saat kasus dibuat, atau `:= 1`
untuk jalur *close*/*reject* yang memakai satu penyetuju tetap.

**Jawaban K4:** `.KomiteLoop` = **jumlah baris `KomiteList`**, ditetapkan **sekali saat kasus
komite dibuat**. **Nol penulis di `Komite Claim Prop`** — jadi di jalur modul ini ia memang tidak
bisa berubah di tengah jalan.

↳ Pengecualian korpus-penuh di lini **Fac In** dipindahkan ke **§11 butir 2** — pokoknya di modul
`Komite Claim FacIn`, dan **tidak menahan apa pun di Komite Claim Prop**.

⛔ Butir `[terbuka]` lama **"`KomiteLoop` disimpan atau dihitung" DICABUT.** Ia opini asisten,
bukan keputusan siapa pun. Keputusan yang sah ada di **§10.3**: `KOMITE_LOOP` adalah kolom yang
**disimpan**.

### 9.3 K3 — `DateApprove` dipakai, `DateApproval` tidak: **[tertutup]**

`[terverifikasi]` Sensus korpus-penuh:

| | Jendela 555 | Korpus penuh |
| --- | --- | --- |
| Berkas **tampilan** (Section/Harness) ber-`DateApprove` | 5 | **12** (11 Section + 1 Harness) |
| Berkas **tampilan** ber-`DateApproval` | **NOL** | **NOL** ✅ **dikonfirmasi** |
| Total berkas ber-`DateApproval` | 2 | **4** — + `Komite Claim Non Prop` ×2 |

`DateApproval` hanya hidup di **Activity**, tidak pernah ditampilkan. Di Pega keduanya diisi
`@CurrentDateTime()` dalam langkah yang sama — satu peristiwa, dua properti, nilai identik.

`[keputusan work owner]` 2026-09-18 — **satu kolom saja: `DATE_APPROVE`**, mengikuti berkas
struktur Komite Claim Life. Karena nilainya identik, migrasi data lama **tidak kehilangan apa pun**.
**Pertanyaan ini tertutup di meja work owner** — tidak diusulkan lagi.

↳ Butir ringan `CreatClaimAnalysis_Act` dipindahkan ke **§11 butir 4** — berkasnya milik modul
**Claim Prop**.

### 9.4 K5 — `.KomiteID` operator, `.IDKomite` jabatan: **[tertutup]**

`[terverifikasi]` Dua properti bernama nyaris sama, **isinya berbeda sama sekali**:
`.KomiteID := .OPERATOR_ID` (**akun operator**) dan `.IDKomite := .JABATAN` (**jabatan/posisi**).

**Jawaban K5:** `.KomiteID` = **operator** — bukan jabatan, bukan workbasket. `.IDKomite` =
**jabatan**, dipakai di teks riwayat (*"Accepted by …"*, lihat §9.1) dan di dokumen.

`[keputusan work owner]` 2026-09-18 — kolom `T_KOMITE_KOMITELIST` sudah diganti nama supaya cocok
isinya (`KOMITE_ID` → `KOMITE_OPERATORID`, `OPERATORID_KOMITE` → `KOMITE_JABATAN`).
⛔ **Penggantian itu dikerjakan prompt lain** — berkas struktur **tidak disentuh dari sini**.

#### ⚠️ Nama orang di-hardcode — temuan kedua

`[terverifikasi]` `Claim Prop/Activity/SendCloseClaimToKomite.xml`:

```
pyWorkPage.ClaimData.ClaimComitee(<last>).KomiteID   := "CHRISTINEANGELINA"
childPageKomite.KomiteList(1).KomiteID               := "CHRISTINEANGELINA"
pyWorkPage.ClaimData.ClaimComitee(<APPEND>).IDKomite := "Claim Dept. Head"
childPageKomite.KomiteList(1).IDKomite               := "Claim Dept. Head"
```

Yang **kedua** setelah `GetLimitDirekturUtama_SQL` yang mencari `where name = 'EKA'`. Orangnya
pindah, jalur *close* komite **diam-diam salah alamat**. Dicatat sebagai temuan; **perbaikannya
meja work owner**.

### 9.5 K10 — kode arah **di Komite Claim Prop**: 2 baris hidup

`[terverifikasi]` **Di `Komite Claim Prop`**, dengan parser bersarang, kode 1/4/5 ada **6 baris** —
dan atribusi work owner **terbukti benar**:

| Berkas | Step | T/F | Gerbang | Syarat |
| --- | --- | --- | --- | --- |
| `KomitePostAdjustment` | 4 | `_/1` | **tidak berlaku** (`false`) | `TransferType=="2"` |
| `KomitePost_Close` | 1 | `_/1` | **tidak berlaku** (`false`) | `TransferType=="1" && …` |
| `KomitePost_Reject` | 1 | `_/1` | **tidak berlaku** (`false`) | `TransferType=="1" && …` |
| `KomitePostAdjustment` | **12** | `1/3` | **HIDUP** | `pyWorkPage.AcceptStatus=="2"` |
| `KomitePostAdjustment` | **26** | `_/4` | **HIDUP** | `…AdjustmentList(…).AcceptanceStatus` |
| `SendErrorDirectKasir` | 2.3 | `5/2` | **HIDUP** | `IsCLMP` — kode 5, `[terbuka]` lama |

✅ Tiga baris berkode 1 memang duduk di step bergerbang **mati** — **dikonfirmasi dengan parser
bersarang**, bukan pemisahan datar. **Sisa hidup: 2 baris.**

↳ Sensus korpus-penuh **3.052 Activity / 54.919 baris syarat** dipindahkan ke **§11 butir 3**.

**Yang TETAP milik modul ini:** **2 baris hidup** berkode 1 dan 4 di `KomitePostAdjustment`
(step **12** `1/3`, step **26** `_/4`), ditambah **1 baris berkode 5** di `SendErrorDirectKasir`
step 2.3. Ketiganya **belum punya arti**, dan salah satu di antaranya menggerbangi blok penomoran
akseptasi (§8). **K10 tetap `[terbuka]` untuk modul ini** — bukan karena besarnya di korpus,
melainkan karena ketiga baris itu ada di sini.

### 9.6 K8 — `.KursIDR` dibekukan, pengisinya Connect-SQL: **[tertutup]**

`[terverifikasi]` `Claim Prop/Activity/SetNameCurrency_Act.xml`:
`.KursIDR := @if(.KursIDR == "", SearchNilaiKursOutput.pxResults(1).NILAIKURS, .KursIDR)`.
Pola *ambil-hanya-bila-kosong* = nilainya **dikunci begitu terisi**; lookup jalan **sekali**.
`AddKomiteTreatyChild_ACT` hanya **menyalinnya** ke halaman komite.

**Jawaban sebagian:** kurs **disimpan dan dibekukan**, tidak dihitung ulang.

#### Pengisi `SearchNilaiKursOutput` — **terjawab**

`[data work owner]` 2026-09-18 — `SearchNilaiKursOutput` diisi oleh **rule Connect-SQL**, bukan
activity. **Itu sebabnya penyisiran activity saya kosong.** Berkasnya `CurrencyStandard.xml`:

```sql
SELECT POOLDATA.GETCURRENCYSTANDARD( {SearchNilaiKursInput.Currency}, SYSDATE ) AS nilaiKurs
FROM DUAL
```

**Tiga fakta yang mengikat:**

1. **Dua masukan saja:** kode mata uang + `SYSDATE`. **NOL nomor kasus · NOL lini · NOL tanggal
   adjustment · NOL tanggal kasus.**
2. `SYSDATE` = **tanggal server basis data SAAT QUERY DIJALANKAN**. Jadi kurs yang terkunci
   adalah kurs **TANGGAL LOOKUP DIJALANKAN** — **bukan** tanggal kasus dibuat, dan **bukan**
   tanggal komite menyetujui.
3. `GETCURRENCYSTANDARD` adalah **fungsi tersimpan di Oracle**, bukan logika Pega. **Isinya TIDAK
   ADA di korpus.**

`[keputusan work owner]` 2026-09-18 — **"ikuti aja query di xml CurrencyStandard."** Artinya:
**dipanggil apa adanya, dengan dua masukan yang sama.** Isi fungsinya **tidak diminta ke DBA**.
Konsisten dengan preseden **"tiru apa adanya"** di K7 (§10.1).

**K8 `[tertutup]`.**

⭐ **Identitas rule-nya, diverifikasi — bukan diterka dari nama berkas** (pelajaran §11):
`[terverifikasi]` berkas `Komite Claim Prop/RDBList/CurrencyStandard.xml` ber-`pxInsName`
**`ASM-FW-GISFW-INT-CURRENCYSTANDARD!ASM!CURRENCYSTANDARD`** — kelas
**`ASM-FW-GISFW-Int-CurrencyStandard`**, nama rule **`ASM!CURRENCYSTANDARD`**. Berkasnya ada **di
dalam modul ini**, dan cocok dengan dua catatan lama: baris `ASM!CURRENCYSTANDARD` di silsilah
Save-As (§2a) dan `CURRENCYSTANDARD` di daftar "baca saja, tanpa COMMIT" (§3.7).

⚠️ Nama berkas `CurrencyStandard.xml` **juga muncul di enam modul lain**. Sesuai pelajaran §11,
**itu tidak berarti rule yang sama** — dan **tidak saya sisir**, karena di luar lingkup modul ini.

### 9.7 K9 — sisi TULIS modul ini: **selesai dan benar**

`[terverifikasi]` Yang dikerjakan **Komite Claim Prop** atas `HISTORYAKSEPTASIPEGA` hanyalah
**menulis**: `ASM!INSERTHISTORYAKSEPTASIPEGA`, 6 kolom, `COMMIT` sendiri, diisi step **32**
(param) → **33** (RDB) dalam penomoran saya (§4). Kolom `Workbasket` diisi **literal `"KLAIM"`**,
selalu, tanpa cabang.

**Sisi tulis itu sudah lengkap dan benar; tidak ada yang tersisa untuk modul ini.**

↳ Pembacanya — `ASM!GETFLAGREJECT_SQL` dengan 4 `SELECT` — duduk di **Fac In / Treaty In**, dan
dipindahkan ke **§11 butir 1**. Ia **tidak menahan apa pun** di Komite Claim Prop.

### 9.8 `TransferType == 1` — **sisa jalur lama**

`[terverifikasi]` Penulis `.TransferType` korpus-penuh: **6 penulisan di 6 berkas**, seluruhnya
saat kasus komite **dibuat** — nilainya `2` (adjustment), `3` (reject), `4` (close), dan satu
`:= .Type` di `Claim Non Prop/Activity/CreateChildKomiteCNP_Act.xml`.

⭐ **Tidak ada satu pun penulis yang menetapkan `1`.** Yang ada hanya **pembaca**:
`TransferType=="1"` muncul 7× sebagai teks syarat, dan **ketiga step yang memakainya bergerbang
mati** (§9.5). Ditambah `KomitePost` yang memilah hanya 2/3/4.

**Kesimpulan:** `TransferType == 1` adalah **sisa jalur lama** — bukan arah keempat.

↳ Butir ringan `:= .Type` dipindahkan ke **§11 butir 5** — berkasnya milik modul
**Claim Non Prop**.

---

## 10. Tiga keputusan work owner — 2026-09-18

> Ketiganya sudah diucapkan work owner tetapi belum pernah sampai ke berkas mana pun. Ditulis apa
> adanya. Ketiganya **menutup** butir terbukanya.

### 10.1 K7 — sembilan rule ber-`COMMIT` sendiri: **ditiru apa adanya** `[tertutup]`

`[keputusan work owner]` 2026-09-18 — **"ditiru apa adanya."** Kesembilan rule yang memuat `COMMIT`
sendiri (§3.7) **tetap `COMMIT` sendiri di Go**, termasuk `POOLDATA.CLAIMREJECTED`
(`GCNM!INSERTCLAIMREJECTED_SQL`) dan `HISTORYAKSEPTASIPEGA` (`ASM!INSERTHISTORYAKSEPTASIPEGA`).
Konsisten dengan preseden **"ikuti apa adanya"** yang sudah dipakai di modul-modul sebelumnya.

⚠️ **Risiko yang ikut diterima sadar.** Bila langkah sesudahnya gagal, **kesembilan tabel itu tetap
terisi**, dan **tidak ada rule pembatal di korpus** — tidak ada `DELETE`, tidak ada `ROLLBACK`,
tidak ada penanda batal di satu pun dari 23 RDBList yang disisir. **Data separuh jadi adalah
perilaku yang disengaja ditiru, bukan cacat yang terlewat.** Barangsiapa membaca `CLAIMREJECTED`
atau `HISTORYAKSEPTASIPEGA` kelak harus tahu bahwa sebagian barisnya berasal dari alur yang tidak
pernah selesai.

### 10.2 K11 — lingkup garapan komite: **hanya modul ini** `[tertutup]`

`[keputusan work owner]` 2026-09-18 — **"hanya di modul ini."** Komite digarap sebagai
**Komite Claim Prop saja**. Ia **TIDAK dilebarkan** menjadi modul komite lintas-lini.
**Komite Life / Non Prop / Fac In digarap sendiri-sendiri nanti.**

⛔ **RALAT 2026-09-18 — alasan pendukung dicabut, keputusan TIDAK.** Draf pertama paragraf ini
menopang keputusan di atas dengan temuan *"`KomitePost` dan `KomitePostAdjustment` dipanggil dari
empat lini"*. **Temuan itu TIDAK BENAR dan DICABUT** — ia lahir dari mencocokkan nama berkas,
bukan dari `pxInsName`. `[keputusan work owner]` 2026-09-18: **"itu hanya nama yang sama. Ada nama
yang sama tapi tetap menjalankan sesuai modulnya masing-masing."** Lihat ralat penuh berikut
buktinya di §3.5.

**Keputusan lingkup di atas berdiri utuh tanpa penopang itu** — ia keputusan work owner, bukan
simpulan dari bukti. **Risiko rancangan berselisih antar lini tetap diterima sadar**: bila kelak
Komite Life atau Komite Fac In digarap dan rancangannya berbeda dari yang ditulis di sini, selisih
itu diselesaikan saat itu, bukan sekarang.

### 10.3 `KomiteLoop` — **disimpan**, bukan dihitung `[tertutup]`

`[keputusan work owner]` 2026-09-18 — **"ikuti activity yang berjalan."** `KOMITE_LOOP` adalah
**kolom yang disimpan**, dan **diisi ulang di titik yang sama dengan Pega mengisinya**:

| Titik | Rule | Perlakuan |
| --- | --- | --- |
| saat kasus komite **dibuat** | `AddKomiteTreatyChild_ACT` | **isi ulang** |
| lini **Fac In**, di dalam komite | `ApprovalKomite_Act` | **isi ulang** |
| **selain dua titik itu** | — | **tidak disentuh** |

⛔ **Bukan `COUNT(*)` setiap kali dibaca.**

⛔ Butir ringan `[terbuka]` **"`KomiteLoop` disimpan atau dihitung" DICABUT** — tertutup oleh
keputusan ini. Kalimat lama yang mengusulkan `COUNT(*)` juga sudah **dicabut dari §9.2**: ia opini
asisten, bukan keputusan siapa pun, dan tidak seharusnya pernah masuk ke berkas ini.

---

## 11. Temuan lintas-modul — BUKAN pekerjaan modul ini

> `[keputusan work owner]` 2026-09-18 — **"intinya ini aku mau supaya ini hanya bahas Komite Claim
> Prop."** **Tujuh** butir di bawah ini **terverifikasi**, tetapi **pokok bahasannya ada di modul lain**,
> dan **tidak satu pun menahan apa pun di Komite Claim Prop**. Mereka dipindahkan ke sini utuh —
> tidak dihapus — supaya modul pemiliknya menemukannya ketika gilirannya tiba. **Tidak ada satu pun
> di antaranya yang masuk daftar terbuka atau frontier modul ini.**

### Butir 1 — `HISTORYAKSEPTASIPEGA` **dibaca** oleh `GETFLAGREJECT_SQL`

**Modul pemilik:** `NB FacIn` · `RNW Fac In` · `Endorsment Fac In` · `NB Treaty In` ·
`EDM Treaty In` — **bukan** Komite Claim Prop. *(asal: K9)*

Di jendela 555 work owner, tabel ini hanya disentuh dua tempat dan **nol SELECT**. Korpus-penuh
berkata lain. `[terverifikasi]` Tabel muncul di **56 berkas**, dan **ada 4 `SELECT`** — rule
`ASM!GETFLAGREJECT_SQL` (berkas `GetFlagReject_SQL.xml`, 4 modul):

```sql
SELECT COUNT(*) AS CARI1 FROM HISTORYAKSEPTASIPEGA
 WHERE ID_PEGA = {DataSearch.CARI20}
   AND WORKBASKET != 'ReasFacInMarketing'
   AND WORKBASKET != 'ReasFacInTeamLeader'
   AND STATUS = 'REJECT'
```

⭐ Dua akibat:

1. Tabelnya **BUKAN** murni jejak audit keluar — ia **dibaca dari dalam Pega**, untuk menghitung
   riwayat penolakan sebuah kasus.
2. **Kolom `WORKBASKET` DIPAKAI sebagai saringan** — jadi kolom itu **masih berguna**, di jalur
   Fac In / Treaty In. Literal `"KLAIM"` yang ditulis jalur komite **lolos saringan `!=`**,
   sehingga baris komite **ikut terhitung** sebagai riwayat reject.

⚠️ `[terbuka]` **milik modul Fac In / Treaty In dan work owner**, bukan milik modul ini: apakah
ikut-terhitungnya baris komite itu disengaja. ⛔ **Jangan usulkan membuang tabel atau kolomnya.**

**Yang sudah selesai di sisi Komite Claim Prop:** menulis, 6 kolom, `COMMIT` sendiri, `Workbasket`
literal `"KLAIM"` (§9.7). Tidak ada sisa.

### Butir 2 — `KomiteLoop` diisi ulang di tengah jalan pada lini Fac In

**Modul pemilik:** `Komite Claim FacIn`. *(asal: pengecualian K4)*

`[terverifikasi]` `Komite Claim FacIn/Activity/ApprovalKomite_Act.xml` menulis
`pyWorkPage.KomiteLoop := @LengthOfPageList(pyWorkPage.KomiteList)` — **di dalam modul komite,
bukan saat kasus dibuat**. Untuk lini Fac In, `KomiteLoop` **dapat dihitung ulang di tengah jalan**.
Tidak terlihat di jendela 555 karena Komite Claim FacIn tidak ter-stage.

**Di Komite Claim Prop pengecualian ini tidak ada** — nol penulis `.KomiteLoop` di modul ini
(§9.2), jadi K4 di sini tertutup penuh. Keputusan §10.3 sudah memuat titik `ApprovalKomite_Act`
sebagai titik isi-ulang, sehingga modul Fac In tinggal mengikutinya.

### Butir 3 — sensus kode arah se-korpus

**Modul pemilik:** seluruh 20 modul korpus — **pekerjaan lintas-proyek**, bukan modul ini.
*(asal: K10 sisi luas)*

`[terverifikasi]` Sensus **3.052 Activity / 54.919 baris syarat** korpus-penuh:

| Kode | Jumlah | | Kode | Jumlah |
| --- | --- | --- | --- | --- |
| `2/4` | **202** | | `1/1` | **38** |
| `1/2` | **179** | | `_/4` | **34** |
| `2/1` | **108** | | `1/_` | **17** |
| `_/1` | **64** | | `4/3` | **3** |
| `1/3` | **47** | | `3/1` | **3** |
| `4/2` | **40** | | `1/6` | **2** |

**Lebih dari 700 baris syarat memakai kode 1 atau 4**; kode **5** (`5/2`) dipakai **1.018 baris**.
Ia **lubang sistematis di aturan baca kita**, dan setiap modul akan menabraknya lagi.

**Yang TETAP milik Komite Claim Prop** (§9.5, tetap `[terbuka]` di sini): **2 baris hidup** berkode
1 dan 4 di `KomitePostAdjustment` — step **12** (`1/3`) dan step **26** (`_/4`) — ditambah **1
baris kode 5** di `SendErrorDirectKasir` step 2.3.

### Butir 4 — `CreatClaimAnalysis_Act` menulis `DateApproval`

**Modul pemilik:** `Claim Prop`. *(asal: butir ringan K3)*

`[terverifikasi]` `Claim Prop/Activity/CreatClaimAnalysis_Act.xml` menulis `DateApproval` **2×**.
Ke mana nilai itu mengalir **belum ditelusuri; penelusurannya milik Claim Prop.**

Keputusan `DATE_APPROVE` sendiri **sudah tertutup** dan tidak terpengaruh butir ini.

### Butir 5 — `TransferType := .Type` di Claim Non Prop

**Modul pemilik:** `Claim Non Prop`. *(asal: butir ringan §C4 / §9.8)*

`[terverifikasi]` `Claim Non Prop/Activity/CreateChildKomiteCNP_Act.xml` satu-satunya penulis
`.TransferType` yang **tidak memakai literal** — ia menulis `:= .Type`. Apakah `.Type` di sana
pernah bernilai `1` **belum ditelusuri**.

**Di Komite Claim Prop tidak berpengaruh:** kesimpulan §9.8 berdiri atas penulis modul ini, yang
seluruhnya literal `2` / `3` / `4`, dan atas tiga step pembaca `TransferType=="1"` yang **bergerbang
mati**.

### Butir 6 — 55 rule `When` baru vs P14

**Modul pemilik:** `Claim Prop`. *(asal: F7 / §E9)*

Apakah **55 rule `When` baru di Claim Prop** menutup **P14** (kode jenis usaha) belum dijawab.
Berkasnya seluruhnya di folder `Claim Prop`, **nol di `Komite Claim Prop`**. Keluar dari frontier
modul ini.

### Butir 7 — pemakai `HISTOAKSPEGA_INDEX4` / `INDEX5`

**Modul pemilik:** kelas kerja berawalan `ASM-FW-GISFW-WORK` — **di luar Komite Claim Prop**.
*(asal: `[data DBA]` K6, §12.1)*

`[data DBA]` Dua dari lima index `POOLDATA.HISTORYAKSEPTASIPEGA` memotong awalan kelas
`ASM-FW-GISFW-WORK ` dari `ID_PEGA`:
`HISTOAKSPEGA_INDEX4 (REPLACE("ID_PEGA",'ASM-FW-GISFW-WORK ',''))` dan
`HISTOAKSPEGA_INDEX5 (ID_PEGA, REPLACE("ID_PEGA",'ASM-FW-GISFW-WORK ',''))`.

Keberadaan keduanya membuktikan **ada pihak yang mencari dengan bentuk yang awalannya sudah
dipotong**. Awalan itu menunjuk kelas kerja **di luar modul ini**. **Siapa pihak itu BUKAN
pekerjaan modul ini dan tidak ditelusuri.**

---

### ⛔ Catatan pembatalan — bekas butir "`KomitePost*` dipanggil dari empat lini"

Salah satu butir bab ini pernah berbunyi *"`KomitePost*` dipanggil dari empat lini"*, dengan peta
pemanggil ke Claim Non Prop, Komite Claim Life, Komite Claim Non Prop, dan Komite Claim FacIn.

**Temuan itu TIDAK BENAR. DIBATALKAN seluruhnya** — bukan dipindah ke modul lain, karena tidak
ada modul lain yang memilikinya. `[keputusan work owner]` 2026-09-18: **"itu hanya nama yang sama.
Ada nama yang sama tapi tetap menjalankan sesuai modulnya masing-masing."** Bukti `pxInsName`
selengkapnya di **§3.5**.

⭐ **Pelajaran, satu baris: kecocokan nama berkas bukan bukti rule yang sama. Identitas rule =
class/nama dari `pxInsName`.**

**§11 berisi 7 butir**, nomor 1–7 di atas. Catatan ini **bukan** salah satunya.

---

## 12. `[data DBA]` DDL dua tabel tanpa awalan skema — K6 `[tertutup]`

`[data DBA]` DDL asli diserahkan work owner, **2026-09-18**. Keduanya ada di skema **`POOLDATA`**:

| Tabel | Tablespace |
| --- | --- |
| `POOLDATA.HISTORYAKSEPTASIPEGA` | `POOLMASTER` |
| `POOLDATA.T_STORAGE_IMAGE` | `TBS_POOLDATA` |

**Jawaban K6:** keduanya memang ditulis rule **tanpa awalan skema** (§3.7), dan itu **terselesaikan
lewat skema bawaan koneksi, yaitu `POOLDATA`**. **Tidak ada kejanggalan.** **K6 `[tertutup]`.**

### 12.1 `HISTORYAKSEPTASIPEGA` — **7 kolom**, jalur komite mengisi **6**

`[data DBA]`

| Kolom | Tipe | Diisi jalur komite? |
| --- | --- | --- |
| `ID_PEGA` | teks 150 | ✅ |
| `TGL_TRANSFER` | tanggal | ✅ |
| `STATUS` | teks 150 | ✅ |
| `USERNAME` | teks 150 | ✅ |
| `WORKBASKET` | teks 150 | ✅ |
| `ID_KOMITE` | teks 50 | ✅ |
| `OPERATORID` | teks 50 | ⚠️ **TIDAK** |

`[terverifikasi]` `InsertHistoryAkseptasiPega_Sql` mengisi **enam**: `ID_PEGA` · `Tgl_Transfer` ·
`Status` · `Username` · `Workbasket` · `ID_KOMITE`. **Cocok dengan §4**, yang menghitung 6 kolom.

⚠️ **`OPERATORID` tidak pernah diisi jalur komite.** `[terbuka]` — **adakah penulis lain di luar
modul ini, atau kolom itu memang selalu kosong.** ⛔ **Tidak ditebak**, dan **tidak diusulkan
dibuang**.

**Lima index, semuanya bertumpu pada `ID_PEGA`** `[data DBA]`:

```
HISTOAKSPEGA_INDEX1   (ID_PEGA)
HISTOAKSPEGA_INDEX2   (ID_PEGA, TGL_TRANSFER)
HISTOAKSPEGA_INDEX3   (ID_PEGA, STATUS, USERNAME)
HISTOAKSPEGA_INDEX4   (REPLACE("ID_PEGA",'ASM-FW-GISFW-WORK ',''))
HISTOAKSPEGA_INDEX5   (ID_PEGA, REPLACE("ID_PEGA",'ASM-FW-GISFW-WORK ',''))
```

⭐ **`INDEX4` dan `INDEX5` memotong awalan kelas `ASM-FW-GISFW-WORK ` dari `ID_PEGA`.**
Keberadaannya membuktikan **dua hal**:

1. Isi `ID_PEGA` **bukan nomor kasus polos**, melainkan **kunci instance penuh berikut awalan
   kelas**.
2. **Ada pihak yang mencari dengan bentuk yang awalannya sudah dipotong.**

✅ **Cocok dengan yang sudah tercatat:** jalur komite mengisi `ID_PEGA` dari
`TempOpenPage.pzInsKey` — memang **kunci instance penuh**. Dua sumber yang berbeda, satu
kesimpulan yang sama.

⚠️ **Fakta untuk Go, bukan usulan:** `ID_PEGA` adalah **teks kunci kasus, panjang 150** — **bukan
angka** dan **bukan foreign key basis data**.

⛔ Awalan `ASM-FW-GISFW-WORK` menunjuk **kelas kerja di luar Komite Claim Prop**. **Siapa pemakai
kedua index itu bukan pekerjaan modul ini** — dicatat sebagai **§11 butir 7**, **tidak ditelusuri**.

### 12.2 `T_STORAGE_IMAGE` — **8 kolom**, PK `IMAGEID`

`[data DBA]`

| Kolom | Tipe | Catatan |
| --- | --- | --- |
| `IMAGEID` | teks 200, **NOT NULL** | ← **PRIMARY KEY** (`T_STORAGE_IMAGE_PK`) |
| `URLPUBLIC` | teks 4000 | |
| `APPFOLDER` | teks 4000 | |
| `EXPDATE` | tanggal | |
| `FILENAME` | teks 4000 | |
| `APPNAME` | teks 100 | |
| `STORAGE` | teks 50 | |
| `TANGGAL_UPLOAD` | tanggal | **DEFAULT `SYSDATE`** |

`T_STORAGE_IMAGE_PK` = **unique index atas `IMAGEID`**, tablespace `POOLDATA_IDX`, constraint
**ENABLE VALIDATE**. **Nol foreign key.**

⚠️ **Kolom `APPNAME` dan `APPFOLDER` menunjukkan tabel ini dipakai bersama banyak aplikasi**, bukan
milik satu modul. Ditulis sebagai **fakta**.

### 12.3 Keduanya `[data DBA]`, bukan rancangan proyek ini

⛔ **Bentuk, kolom, tipe, index, dan constraint kedua tabel ini mengikuti DDL yang sudah ada di
produksi, dan TIDAK ditetapkan ulang oleh rancangan mana pun di modul ini.**
