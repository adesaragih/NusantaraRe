# Grilling ronde 3 — Claim Fac In

**Tanggal:** 2026-09-19 · **Korpus:** `D:\XML\RNM_BRD\Claim Fac In\` — **482 berkas `.xml`**, READ-ONLY.
**Berkas ini MENAMBAH.** ⛔ `grilling-ronde-1.md`, `grilling-ronde-1-ulang-docs.md`, dan
`grilling-ronde-2.md` tidak disentuh satu byte pun.

> **Cara membaca angka di berkas ini.** Setiap sensus menyebut **jendelanya**. Angka yang berbeda
> dari ronde sebelumnya **dilaporkan berdampingan**, tidak ditimpa. Bukti selalu **path berkas +
> nama rule + nomor langkah Pega**; ⛔ **nol nomor baris XML**.

> ### ⚠️ Aturan yang dipatuhi ketat di ronde ini
>
> `pyStepsPreCondition = "false"` mematikan **GERBANGNYA**, bukan **LANGKAHNYA**. Langkah dimatikan
> **hanya** oleh `pyStepsBlockName = "//"`. ⛔ Gerbang mati berarti langkah berjalan **lebih
> sering**. Ronde 2 sempat menyimpulkan sebaliknya dan diralat work owner.

---

## LANGKAH 0

```
0a  isi .scratch\claim-facin\        3 berkas   ✅ grilling-ronde-3.md belum ada
0b  panjang ketiganya                712 · 735 · 1145   ✅ cocok
0c  berkas .xml korpus               482        ✅
0d  §H1 butir yang masih berdiri     10         ✅
0e  docs\adr\*.md                    15         ✅
```

---

## §A — Dasar, dan butir terbuka yang dibawa

### A1 — Aturan baca yang dipakai

Sama seperti ronde 2 §A1, dan ⛔ **berkas kanonnya tetap tidak ada**:

| Dasar | Sumber | Status |
| --- | --- | --- |
| Sembilan aturan baca ekspor Pega | lampiran di `komite-claim-prop\spec.md` | dipakai |
| Aturan **DataPage** | ronde 1 §A2 | **USULAN** |
| Empat butir tambahan | ronde ulang §C3·C4 | **USULAN** |

⛔ **Usulan tetap usulan.** Yang berhak mengesahkan adalah **work owner**.

### A2 — Butir `[terbuka]` yang dibawa — **sembilan, nol ditutup**

| # | Butir | Sesudah ronde 3 |
| --- | --- | --- |
| 1 | Kenapa `CreateKMTNo_Act` langkah **13.1** ditimpa **13.1.1** | **masih terbuka** |
| 2 | Apakah `FlagOnGoingCommitte` Fac In kolom yang **sama** dengan Claim Prop | **masih terbuka** |
| 3 | Apakah `TempOutstanding` disimpan oleh rule **lain** | **masih terbuka** |
| 4 | Arti `REQUIRED = -1` versus `0` | **masih terbuka** — ⚠️ dipakai di §C, **ditandai** |
| 5 | Isi tiap penampung `CARI<n>` | **masih terbuka** — ⭐ bertambah bukti, lihat §D |
| 6 | `IndexObject` versus `.ObjectIndex` | **masih terbuka** |
| 7 | Apakah panggilan **arasapas tanpa saringan** memang dikehendaki | **masih terbuka** |
| 8 | Pemetaan *"aturan A–J"* terhadap sembilan aturan bernomor | **masih terbuka** |
| 9 | Di mana catatan aturan baca disimpan | **masih terbuka** |

⛔ **Nol di antaranya ditutup di ronde ini.**

---

## §B — 79 kandidat `pyMemo` dibuka satu per satu

### B1 — Angka dasar: **ketiganya COCOK** dengan ronde 1 dan ronde 2

`[terverifikasi]` **459 catatan** di **449 berkas**, **457 pasangan berbeda**. Golongan menurut
jendela ronde 2 *(kata kuncinya dikutip penuh di ronde 2 §D1)*: **(d) 32 · (c) 91 · sisanya 336**.
Dari 91, **12 sudah diuji ronde 2**, jadi **sisa 79** — dan **ketujuh puluh sembilannya dibuka di
ronde ini**.

### B2 — ⭐ Yang paling mengubah gambaran: **hanya 12 dari 79 adalah Activity**

`[terverifikasi]` Sebaran jenis rule ke-79:

```
Section 25 · Activity 12 · RDBList 11 · Harness 8 · ReportDefinition 6
DataTransform 5 · FlowAction 5 · When 3 · DataPage 2 · ConnectREST 1 · Flow 1
```

⭐ **67 dari 79 bukan Activity** — jadi *"membuka rule-nya"* untuk mayoritas berarti membaca **isi
rule**-nya, bukan langkahnya. Ronde 1 menggolongkan seluruh 459 catatan dengan cara yang sama untuk
semua jenis; itu sebabnya golongan (c) menggelembung.

### B3 — Pemeriksaan keberadaan: **67 diuji, 67 ADA, NOL tidak ketemu**

**Jendela, disebut penuh:** untuk tiap catatan, ambil setiap kata beralfabet sepanjang ≥ 4 huruf
yang **bukan** kata umum *(daftar henti: tambah · ubah · fix · test · coba · buat · make · new · to ·
the · di · ke · untuk · dan · atau · yg · yang · agar · unutk · dipake · dipakai · jadi · biar ·
sblm · dulu · lansung · langsung · kolom · pilih)*, lalu cari kata itu di **isi berkas rule-nya**,
**tidak peka huruf besar-kecil**.

```
diuji keberadaan              67 dari 79
yang dicarinya ADA            67          ⭐ seluruhnya
TIDAK KETEMU                   0
tak punya kata uji            12          catatannya hanya "fix" · "test" · "ubah"
```

⭐ **Nol catatan yang menjanjikan sesuatu dan tidak ada jejaknya.** ⛔ Ini **bukan** bukti bahwa
perubahannya sudah benar — hanya bahwa **yang disebut catatan itu ada di rule-nya**.

### B4 — Kedua belas yang tak punya kata uji, dibuka satu per satu

| Rule | Catatan | Hasil |
| --- | --- | --- |
| ⭐ `RDBList/InsertClaimPNC.xml` | *"test"* | ⭐ **BUKAN uji coba** — ia **pemanggil procedure sungguhan**: `BEGIN POOLDATA.PEGA_JSON_KLAIM_PNC(…); COMMIT; END;` ⚠️ **ber-`COMMIT` sendiri** |
| ⭐ `RDBList/GetLinkStorage_SQL.xml` | *"TEST"* | ⭐ **BUKAN uji coba** — kueri nyata ke `T_STORAGE_IMAGE`, memberi alamat dan nama berkas |
| `RDBList/CekLunasPremi_Sql.xml` | *"fix"* | rule nyata; **apa yang diperbaiki tidak tertulis** — `[terbuka]` |
| `RDBList/GetProgressClaim_SQL.xml` | *"tambah asc"* | ✅ `asc` **ADA 2×** di kuerinya |
| `When/whenRiskType.xml` · `When/IsBonding.xml` | *"test"* · *"Tambah untuk PNC"* | `PNC` **ADA 8×**; *"test"* tidak dapat diuji |
| `Activity/DeleteValueEstimation.xml` | *"fix"* | 34 langkah · **2 gerbang mati** · 4 titik pesan |
| `Activity/SetAdjsuter_act.xml` · `SetConsultant_Act.xml` | *"fix"* | **4 langkah** masing-masing, 1 `Call` |
| `Activity/GetPayAttachment_Act.xml` | *"ubah"* | **1 langkah** saja |
| `DataTransform/SetEstimation_DT.xml` · `Section/InputEstimasi.xml` | *"fix"* | tidak dapat diputuskan dari isinya |

⭐ **Dua nama berbohong ke arah yang berbahaya:** rule bernama/berkatatan **"test"** ternyata
**menulis ke basis data** dan **menjalankan `COMMIT` sendiri**. ⛔ Siapa pun yang membuang rule
ber-catatan *"test"* akan membuang penulis data.

### B5 — ⭐ Alias keenam yang ronde 2 lewatkan — dan ia **UANG**

`[terverifikasi]` `RDBList/CariHistoryClaim_SQL.xml`, catatan *"tambah replace di TSI"* — **SUDAH
DIPATUHI**, dan isinya:

```
REPLACE(a.DATA_JSON.ClaimEstimate , ',' , '.') AS "TSI"
```

⭐ **Tiga hal sekaligus:**

1. **Estimasi klaim dinamai `TSI`** — alias **keenam** yang berbohong di rule ini.
   ⚠️ **Ronde 2 §F2 hanya menemukan lima**, karena polanya menuntut spasi sebelum `AS`, sedangkan di
   sini `)AS` menempel. **Ralat di §H3.**
2. ⭐ **Koma diganti titik pada nilai UANG** — inilah tambalan yang ronde ulang §B1 sebut sebagai
   dasar pertentangan **ADR-0003**, kini terbaca kalimatnya.
3. **Perbaikan itu dikerjakan di dalam SQL**, bukan di aplikasi.

---

## §C — 147 parameter `STRING`: mana yang UANG

**Jendela:** ke-**169** parameter dari **178** tanda tangan activity yang terurai *(ronde 2 §E2)*,
disaring dengan pola nama uang/angka **tidak peka huruf besar-kecil**.

### C1 — Hasil saringan: **15 dari 147**, tetapi hanya **tiga** benar-benar uang

| Activity | Parameter | Tipe | Arah | Apa sebenarnya |
| --- | --- | --- | --- | --- |
| ⭐ `SetInitial_ACT` | **`AdjustmentValue`** | **STRING** | IN | ⭐ **nilai penyesuaian — UANG** |
| ⭐ `SetCedant_act` | **`Share`** | **STRING** | IN | ⭐ **persentase share** |
| ⭐ `CountSpreadingClaim_ACT` | **`Estimation`** | **STRING** | IN | ⭐ **nilai estimasi — UANG** |
| `CountSpreadingClaim_ACT` · `ProtectionDate_Act` · `ValidateInputEstimate_act` | `estimationdate` | STRING | IN | tanggal, bukan uang |
| `GetPaymentList_Act` · `SetProtectionEstimation` · `SetTreatNameAdjustment` · `SetKomiteList_ACT` · `SaveCFS_ACT` | `IndexAdjustment` · `idxadjust` · `IdxEstimation` | STRING | IN | **penunjuk posisi**, bukan uang |
| `GetUrlGoogleStorage_Act` · `InsertGoogleStorage_Act` | `Durasi` | STRING | IN | lama berkas, bukan uang |
| `PrintPDFAccep_MultiAksep` | `NoAkseptasi` | STRING | IN | nomor, bukan uang |
| `SetTreatyname` | `Adjustment` | STRING | IN | nama treaty, bukan uang |

### C2 — ⭐⭐ Dan yang bertipe angka justru memuat **satu `Double`**

`[terverifikasi]` Dari 169 parameter, **22 bukan `STRING`**:

```
INTEGER 13 · Date 4 · DateTime 2 · TrueFalse 1 · Decimal 1 · ⭐ Double 1
```

| Activity | Parameter | Tipe | Apa sebenarnya |
| --- | --- | --- | --- |
| `GetCoverageAneka_Act` | `TSI` | **Decimal** | ✅ uang, dan tipenya **tepat** |
| ⭐⭐ **`SetKomiteList_ACT`** | ⭐⭐ **`TotalAdjustment`** | ⭐⭐ **`Double`** | ⭐⭐ **TOTAL PENYESUAIAN — UANG, bertipe bilangan pecahan biner** |

### C3 — Vonis **ADR-0003** *(uang non-float)*: ⚠️ **MENENTANG**

**Dasarnya dua, keduanya terbaca langsung:**

1. ⭐ **`SetKomiteList_ACT.TotalAdjustment` bertipe `Double`** — **bilangan pecahan biner untuk total
   uang**, persis yang ADR-0003 larang.
2. ⭐ **Tiga nilai uang lewat sebagai `STRING`** — `AdjustmentValue` · `Share` · `Estimation` —
   sehingga ketelitiannya **tidak dijaga tipe apa pun**, dan ⭐ **satu di antaranya sudah terbukti
   perlu ditambal koma-ke-titik di SQL** *(§B5)*.

⛔ **Dilarang mengusulkan revisi ADR.** ⛔ Dan ⛔ **saya tidak menyimpulkan apa pun untuk aplikasi
Go** — vonis ini menerangkan **korpus**, bukan menetapkan rancangan.

⚠️ **Arti `REQUIRED = -1` versus `0` tetap `[terbuka]`.** Ia **tidak dipakai** sebagai dasar vonis di
atas; disebut hanya sebagai kolom.

---

## §D — 52 titik `pyDisableSubmit = true` dan aksi lokal → ADR-0012

### D1 — Ke-52 titik: terpusat di **sepuluh** berkas

`[terverifikasi]` **Jendela:** 106 berkas layar — Section 57 · Harness 16 · FlowAction 33.

| Berkas | Titik |
| --- | --- |
| `Section/InputAdjustment.xml` | **11** |
| `Section/Estimasi.xml` | **10** |
| `Section/EstimasiPA.xml` | 6 |
| `Section/PropertyItemListGridEstimation.xml` | 5 |
| `Section/InputRegisterDetail.xml` · `PropertyItemListGridEstimationMBU` · `…Travel` · `ShowItemPA` | 4 masing-masing |
| `Section/MstAdjusterConsultant.xml` · `Harness/MstAdjusterConsultant.xml` | 2 masing-masing |

⭐ **Seluruhnya di layar estimasi dan penyesuaian** — nol di layar komite, nol di layar akseptasi.

### D2 — ⭐ Aksi lokal bernama: **SEMBILAN, bukan empat**

`[terverifikasi]` Ronde 2 §F1 menyebut **empat**; sensus penuh memberi **sembilan**:

| Aksi | Kali | Di berkas |
| --- | --- | --- |
| `SetPassWordSP` | 6 | `Estimasi` · `EstimasiMarine` · `EstimasiPA` |
| ⭐ `PreventRejectClaim` | 4 | `ClaimSurvey` |
| `ShowRetro` | 4 | `Estimasi` · `InputAdjustment` |
| `MessageBeforeDeleteTreatyGroup` | 4 | `MstAdjusterConsultant` *(Harness + Section)* |
| `CatastrofeList` · `SureRejectClaim` · `RejectSurveyClaim` · `ProtectDOL` · `InputSubProgressClaim` | 2 masing-masing | — |

**Ralat angka 4 → 9 tercatat di §H3.**

### D3 — ⭐⭐ Medan wewenang: **ADA banyak, TERISI hampir nol**

`[terverifikasi]` Sensus **seluruh** medan berbau `privileg|authoriz|workgroup|accessgroup|role` di
106 berkas layar — **15 nama medan berbeda**:

| Medan | Hadir | Terisi | Isinya |
| --- | --- | --- | --- |
| `pyPrivilegeView` | 699 | ⭐ **0** | — |
| `pyPrivilegeUpdate` | 699 | ⭐ **0** | — |
| `pyPrivilege` | 419 | ⭐ **0** | — |
| `pyWhenNotPrivilege` | 419 | ⭐ **0** | — |
| `pyActionPrivilegeList` · `pyPrivilegeName` · `pySectionReferencePrivilege*` | 33 | ⭐ **0** | — |
| `WorkGroup` | 16 | ⭐ **0** | — |
| `pyAssociatedPrivileges` | 520 | 431 | ⛔ **seluruhnya `false`** — baku |
| ⚠️ `pyPrivilegeClass` | 33 | **33** | ⚠️ **hanya nama KELAS**: `…Data-Object` · `…Work-PNC` · `…Data-ObjectItem` |

⭐⭐ **Pola yang sama persis dengan Komite Claim Prop:** daftar hak akses **menyebut kelas tanpa
menyebut hak**. Nol nama hak terisi di seluruh modul.

### D4 — `SetDisable_ACT` dan `DisableSendComite_Act`, dibaca

`[terverifikasi]`

| Rule | Yang terbaca |
| --- | --- |
| `SetDisable_ACT` *(GCNMFW 01-01-07, 12 langkah)* | gerbang **satu-satunya** di **7.1**: `.IsFacretro==1` · langkah **7.2** menyetel `.IsPrintAccept := 1` dan `.IsKomite := 1` · ⭐ **langkah 8 · 9 · 10 BER-REMARK**, termasuk `pyWorkPage.AktifButton := 0` |
| `DisableSendComite_Act` *(GCNMFW 01-01-23, 6 langkah)* | **nol gerbang** · langkah 4 menyetel `pyWorkPage.AktifButton := 0` dan `.IsKomite := ""` |

⛔ **Tidak satu pun menonaktifkan berdasarkan WEWENANG atau PERAN.** Yang menggerbanginya adalah
**keadaan data** — `.IsFacretro`, nilai penyesuaian — bukan siapa penggunanya.

### D5 — ⭐ Vonis **ADR-0012**: **TIDAK MENYENTUH** — dan kini **dapat diputuskan**

ADR-0012 menetapkan **wewenang kirim ke Komite bergantung `Type` klaim** *(`TP`/`TR` siapa pun;
`QP`/`QR` hanya SPV)*, dibawa apa adanya sebagai paritas.

`[terverifikasi]` **Penyisiran seluruh 482 berkas `.xml`** untuk pola `.Type == "TP|TR|QP|QR"`,
tidak peka huruf besar-kecil: ⭐ **NOL berkas**.

⭐ **Vonis: TIDAK MENYENTUH.** Mekanisme yang diatur ADR-0012 **tidak punya padanan apa pun di Claim
Fac In** — bukan ditentang, melainkan **tidak ada**.

⚠️ **Tetapi satu fakta harus ikut terbawa ke spec, dan ia lebih besar dari ADR-0012 sendiri:**
`[terverifikasi]` **Claim Fac In tidak punya satu pun gerbang layar berdasarkan wewenang** — ke-15
medan hak akses yang menamai hak **kosong seluruhnya**, dan kedua rule penonaktif tombol
menggerbangi **keadaan data**, bukan peran. ⛔ **Saya tidak menyatakan itu cacat**, dan ⛔ **tidak
menyimpulkan apa pun untuk aplikasi Go.**

⛔ **Dua ronde sebelumnya menulis "tidak dapat diputuskan".** Ronde ini **dapat memutuskannya**
karena jendelanya diperluas dari 11 medan menjadi **121 + 15**, ditambah sensus `.Type` se-modul.
**Ralat tercatat di §H3.**

---

## §E — Alias SQL, dan siapa pemanggilnya

### E1 — ⭐ Tiga jendela, tiga angka — dan **tidak satu pun sama dengan ronde ulang**

**Jendela disebut penuh**, seluruhnya atas **63 berkas `RDBList`**, medan `pyBrowseSQL` + `pySQL`,
pembanding memakai potongan terakhir sesudah titik, tidak peka huruf besar-kecil:

| Jendela | Pasangan | TIDAK COCOK | Cocok |
| --- | --- | --- | --- |
| **A** — alias berkutip, sumber longgar *(boleh memuat fungsi dan koma)* | **48** | **48** | 0 |
| **B** — alias berkutip, sumber ketat *(satu kata)* — **jendela ronde 2** | **42** | **39** | 3 |
| **C** — termasuk alias **tanpa kutip** | **371** | **368** | 3 |
| *(ronde ulang §E2, jendelanya tidak tertulis)* | *50* | *38* | *12* |

⛔ **Keempat angka berdiri berdampingan.** Beda empat cara = **belum punya data** *(aturan I3/P4)*.
⚠️ **Jendela C memperlihatkan sebabnya:** menghitung alias tanpa kutip melipatgandakan pasangan
menjadi **371**, karena Oracle memperlakukan hampir setiap `X Y` sebagai alias. **Angka 50 milik
ronde ulang ada di antara A dan C**, jadi jendelanya kemungkinan besar setengah longgar.

⛔ **Ketiga positif palsu yang ronde ulang §H5 butir 3 sebut TIDAK KETEMU** — berkas itu **tidak
menyebutkan yang mana**, dan ketiga jendela saya memberi himpunan berbeda. ⛔ **Tidak saya karang.**
`[terbuka]`

### E2 — ⭐ Pemanggilnya **satu-satu**, dan **tidak ada yang memakai keduanya**

`[terverifikasi]` Penyisiran seluruh 482 berkas `.xml`:

| Rule | Pemanggil | Jumlah |
| --- | --- | --- |
| `CariHistoryClaim_SQL` *(6 alias, 5 berbohong)* | `Activity/CheckDoubleClaim_Act.xml` | **1** |
| `BrowseHistoryClaim` *(3 alias, 2 berbohong)* | `Activity/CallActivityInputRegister.xml` | **1** |

⭐ **Nol pemanggil memakai keduanya** — jadi kedua pemetaan bohong yang berbeda itu **tidak pernah
bertemu dalam satu jalur**. ⚠️ Tetapi keduanya membaca **tabel yang sama**, `json_klaim`, dengan
**nama kolom hasil yang berbeda-beda** — siapa pun yang membandingkan keluaran keduanya akan
mengira datanya berbeda.

---

## §F — Lima rule `ADESAMUEL@`, isinya

### F1 — ⭐ `Register_Flow` — **16 bentuk, dan NOL bentuk SELESAI**

`[terverifikasi]` `Flow/Register_Flow.xml` — 174.848 byte, kelas `ASM-FW-GCNMFW-Work-PNC`,
ruleset **`ADESAMUEL@` 01-01-01**:

| Jenis bentuk | Jumlah |
| --- | --- |
| `Data-MO-Connector-Transition` *(penghubung)* | **9** |
| `Data-MO-Gateway-Decision` *(gerbang keputusan)* | **3** |
| `Data-MO-Activity-Assignment` *(tahap menunggu manusia)* | **3** |
| `Data-MO-Event-Start` *(mulai)* | **1** |
| ⭐ **`Data-MO-Event-End` (selesai)** | ⭐ **NOL** |

**Nama bentuk yang terbaca:** `Input Register` · `InputRegister` · `Input Estimasi` ·
`InputEstimasi` · `Choose Surveyor` · `InputSurveyor` · `SendToEstAdmin` · `AcceptanceClaim` ·
`B2B` · `IsBack` · `IsBackStage` · `IsSPK` · `[Always]` · `[Else]` · `[Result]`.

⭐⭐ **Nol bentuk "selesai" di seluruh berkas alur.** ⚠️ Bandingkan Komite Claim Prop, yang punya
`END52` berstatus akhir `Resolved-Completed`. ⛔ **Saya tidak menyimpulkan kasus Fac In tidak pernah
selesai** — hanya bahwa **berkas alur ini tidak memuat bentuk selesainya**. `[terbuka]` apakah ia
diselesaikan lewat jalan lain *(misalnya activity)* atau bentuknya memang tidak terekspor.

### F2 — Keempat Section

| Section | Byte | `pyUserData` | `pyVisible` terisi | `pyCondition` terisi | Kelas |
| --- | --- | --- | --- | --- | --- |
| `ViewPolis` | 380.257 | 48 | 55 | **5** | `…Work-PNC` |
| `ClaimSurvey` | 346.988 | 48 | 55 | 3 | `…Work-PNC` |
| `ViewHistoryClaim` | 293.448 | 44 | 56 | 3 | `…Work-PNC` |
| `InputInwardFacultativeDtl` | 189.602 | 19 | 20 | 3 | ⚠️ `…Work` *(bukan `Work-PNC`)* |

⚠️ **`InputInwardFacultativeDtl` berkelas berbeda** dari keempat lainnya — `ASM-FW-GCNMFW-Work`,
bukan `…Work-PNC`. ⭐ Dan `ClaimSurvey` inilah yang memuat aksi lokal **`PreventRejectClaim`**
*(§D2)*.

⚠️ **Temuan tentang BERKAS, bukan tentang orang.** ⛔ Nol penilaian terhadap siapa pun.

---

## §G — Arah gerbang kode 4 · 5 · 6, disensus se-modul

**Jendela:** seluruh **179 berkas `Activity`**, kedua keluarga gerbang *(prakondisi dan transisi)*,
kedua sisi *(benar dan salah)*, dibaca dengan pengurai XML rekursif.

```
kode 4  keluar-iterasi     1
kode 5  lewati-when        7      ⚠️ mengubah rantai DAN menjadi ATAU
kode 6  keluar-activity   40
```

### G1 — ⭐ Satu lompatan kode 4, **TERISOLASI**

`[terverifikasi]` **`Activity/ProtectionObjectItem_Act.xml` langkah 3**, keluarga **prakondisi**,
sisi **SALAH**, **flag `true` → gerbangnya HIDUP**, blok tidak ber-remark.

⭐ **Gagal sejak ronde ulang §F3, ketemu di sini.** Artinya: bila syaratnya **tidak** terpenuhi,
**perulangan diputus** — bukan langkahnya yang dilewati.

### G2 — ⚠️ Tujuh kode 5, dan **enam di antaranya HIDUP**

| Rule | Langkah | Keluarga | Flag |
| --- | --- | --- | --- |
| `AgentSourceBiz_Act` | 2 | prakondisi | **true** |
| `CheckEstimateValue` | **24** dan **25** | prakondisi | **true** |
| `SetMOClaim_Act` | 2 | prakondisi | **true** |
| `SpreadingCheckSP` | **3.2** *(dua baris)* | prakondisi | **true** |
| `SendEmailDLA_ACT` | 3 | transisi | (kosong) |

⚠️ **Di keenam titik hidup itu, rantai gerbangnya berperilaku ATAU, bukan DAN.** ⛔ Siapa pun yang
membacanya sebagai DAN akan salah — termasuk ⭐ **`SpreadingCheckSP`, yang menyentuh spreading**, dan
⭐ **`CheckEstimateValue`, rule validasi 60 langkah**.

### G3 — Kode 6, dan ⭐ nilai `0` muncul lagi

`[terverifikasi]` **40 titik keluar-activity**. ⭐ Dua di antaranya ber-flag **`0`** — nilai tak
berdokumen yang sudah menjadi **Pertanyaan 2** ronde 1: `CLaimFaceSheet_Act` langkah **39** dan
`DLAFacintoTreaty_Act` langkah **13.4**.

⚠️ Dan `CheckPeriodePolicy` langkah **1** punya **dua** baris berarah kode 6 dengan flag **`false`**
— gerbangnya tersimpan tetapi **tidak ditegakkan**, sehingga langkahnya berjalan **tanpa saringan**.

---

## §H — Penutup

### H1 — Yang paling menahan untuk ronde 4

| # | Yang dikerjakan | Kenapa menahan |
| --- | --- | --- |
| **1** | ⭐⭐ **Adu Claim Prop dan Komite Claim Prop dengan 15 ADR** | **utang lintas modul**; jawaban work owner menjadikannya *"harus"*, dan kedua modul **tuntas tanpa pernah diadu dengan dokumen** |
| **2** | ⭐ **Baca ulang kesimpulan kurs standar dan penyimpanan berkas di Claim Prop** | rule yang dipakai di sana **4 tahun 3 bulan lebih tua** |
| **3** | ⭐ **Sesuaikan kalimat Claim Prop *"52 dari 61 rule `When` tidak dimigrasikan"*** | keputusan *"rule dipisahkan dari kehidupannya"* menuntutnya |
| **4** | ⭐ **Baca `SetKomiteList_ACT`** — pemilik parameter `TotalAdjustment` bertipe **`Double`** | satu-satunya bukti langsung pertentangan **ADR-0003** |
| **5** | **Baca `SpreadingCheckSP` dan `CheckEstimateValue`** | keduanya memuat kode 5 hidup → rantainya **ATAU**, dan keduanya menyentuh uang |
| **6** | **Telusuri bentuk "selesai" alur Fac In** | `Register_Flow` **nol bentuk selesai** |
| **7** | **Sisir 336 catatan golongan (a)+(b)** | ronde 1–3 hanya menguji golongan (c) dan (d) |
| **8** | **Tetapkan satu jendela alias yang disepakati** | empat angka berdiri: 50 · 48 · 42 · 371 |

### H1b — Daftar penyimpangan sadar modul ini — **tetap 5 terurai + 6 dari ADR**

⛔ **Ronde ini tidak menambah dan tidak mengurangi.** Ia **membaca**, tidak memutuskan.
⚠️ Tetapi **§C3 memperkuat** satu di antaranya: pertentangan **ADR-0003** kini punya bukti langsung
berupa parameter `Double` dan tiga nilai uang bertipe `STRING`.

### H2 — Kesimpulan ronde 3 yang **PALING RAWAN salah**

⚠️ **Yang paling rawan adalah §C1 — memilah "mana yang uang" dari NAMA parameter.**

**Kenapa rawan:** saya menilai dari **nama**, bukan dari **pemakaiannya**. `Share` bisa saja bukan
persentase uang; `Estimation` bisa saja pengenal, bukan nilai. ⛔ **Satu-satunya cara memastikannya
adalah membaca pemanggil tiap parameter**, dan itu **tidak dikerjakan ronde ini**.

Yang **tidak** rawan dan berdiri kokoh: ⭐ **`TotalAdjustment` bertipe `Double`** — tipenya terbaca
langsung dari tanda tangan, bukan ditebak dari nama. **Vonis ADR-0003 bertumpu pada butir ini, bukan
pada saringan nama.**

**Kedua paling rawan:** §B3. *"67 dari 67 ADA"* hanya membuktikan **kata yang disebut catatan ada di
rule-nya** — ⛔ **bukan** bahwa perubahannya sudah benar atau lengkap.

### H3 — RALAT terhadap ronde 1, ronde ulang, dan ronde 2

| # | Yang diralat | Angka/kalimat LAMA, dikutip | Yang benar |
| --- | --- | --- | --- |
| **1** | Alias `CariHistoryClaim_SQL` | ronde 2 §F2: *"di `CariHistoryClaim_SQL.xml` modul ini, **5 alias dan 5-lima-nya berbohong**"* | ⭐ **ENAM alias, LIMA berbohong** — yang keenam `REPLACE(…ClaimEstimate…) AS "TSI"`, terlewat karena `)AS` menempel tanpa spasi |
| **2** | Jumlah aksi lokal | ronde 2 §F1: *"**4 aksi lokal bernama**"* | ⭐ **SEMBILAN** |
| **3** | Vonis ADR-0012 | ronde ulang §B10 dan ronde 2 §G3: *"tidak dapat diputuskan"* | ⭐ **TIDAK MENYENTUH** — nol `.Type==TP/TR/QP/QR` di 482 berkas |
| **4** | Pasangan alias SQL | ronde ulang §E2 *"38 dari 50"* · ronde 2 §F2 *"39 dari 42"* | ⛔ **belum punya data** — empat jendela memberi 50 · 48 · 42 · 371 |
| **5** | Lompatan kode 4 | ronde ulang §F3: *"ronde ini gagal mengisolasinya"* | ⭐ **TERISOLASI** — `ProtectionObjectItem_Act` langkah 3, prakondisi sisi salah, gerbang **hidup** |

### H4 — Yang seharusnya dikerjakan tetapi **TIDAK diperintahkan** blok ini

| # | Butir |
| --- | --- |
| **1** | **Membaca pemanggil tiap parameter uang** — satu-satunya cara menutup kerawanan §C1 |
| **2** | **Menyisir 336 catatan golongan (a)+(b)** — belum pernah disentuh sama sekali |
| **3** | **Memeriksa apakah `TempOutstanding` disimpan rule lain** *(butir terbuka 3)* |
| **4** | **Membaca isi `InsertClaimPNC` sampai ke procedure-nya** — ia ber-`COMMIT` sendiri |
| **5** | **Menyisir 40 titik kode 6** satu per satu — ronde ini hanya menghitungnya |
| **6** | **Memeriksa `pyPrivilegeClass` yang terisi 33** — kelasnya disebut, haknya tidak |

---

## Lampiran — bukti berkas lain tidak disentuh

Sidik jari MD5 diambil **sebelum** ronde ini dan dibandingkan **sesudahnya**:

```
korpus Claim Fac In        482 berkas .xml
berkas lama claim-facin      3 berkas  (712 + 735 + 1.145 baris)
docs\adr\                   15 berkas
modul lain di .scratch\    175 berkas  (claim-prop dan komite-claim-prop TERMASUK)
```

⛔ **Nol berkas korpus dibuka untuk ditulis** — seluruh pembacaan memakai pengurai XML hanya-baca.
⛔ **Nol nilai rahasia disalin.** ⛔ **Nol butir `[terbuka]` ditutup.** ⛔ **Nol keputusan §A5 diubah.**
