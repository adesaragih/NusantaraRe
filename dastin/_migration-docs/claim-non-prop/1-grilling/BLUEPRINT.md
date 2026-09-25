# BLUEPRINT — Claim Non Prop (fakta teknis)

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Scope: `D:\XML_NURE\Claim Non Prop` saja. Folder `Komite Claim Non Prop` tidak dibuka — lihat §8.
Label: **EVIDENCED** = terbukti dari XML · **DECIDED** = keputusan kita (lihat ADR) · **EXTERNAL** = butuh data luar (lihat _selesai/OPEN-QUESTIONS.md).

---

## 1. Inventaris Rule — EVIDENCED

279 berkas XML + 2 non-XML (`Struktur_Flow_TreatyIn.xlsx`, `~$Struktur_Flow_TreatyIn.xlsx`).

| Jml | Tipe Rule | Folder |
|---:|---|---|
| 114 | `RULE-OBJ-ACTIVITY` | `Activity\` |
| 50 | `RULE-CONNECT-SQL` | `RDBList\` |
| 36 | `RULE-HTML-SECTION` | `Section\` |
| 22 | `RULE-OBJ-REPORT-DEFINITION` | `ReportDefinition\` |
| 16 | `RULE-OBJ-FLOWACTION` | `FlowAction\` |
| 14 | `RULE-HTML-HARNESS` | `Harness\` |
| 10 | `RULE-OBJ-MODEL` (Data Transform) | `DataTransform\` |
| 7 | `RULE-CONNECT-REST` | `ConnectREST\` |
| 7 | `RULE-OBJ-WHEN` | `When\` |
| 1 | `RULE-DECLARE-DECISIONTABLE` | `DecisionTable\` |
| 1 | `RULE-OBJ-FLOW` | `Flow\` |
| 1 | `RULE-ADMIN-SYSTEM-SETTINGS` | `SystemSettings\` |

**Tidak ada di export**: `Rule-Obj-Class`, `Rule-Obj-Property`, `Rule-Declare-Expressions`, `Rule-Declare-Pages`, `Rule-Declare-OnChange`, `Rule-Access-When`, `Rule-Obj-FieldValue`, `Rule-Obj-Corr`, `Rule-Obj-HTML`.

**Applies-to class** (48 class): `…Work-ClaimTreatyNonProp` 84 · `…Data-Adjustment` 27 · `ASM-FW-GISFW-Data-PolicyTreatyIn` 19 · `@baseclass` 18 · `…GCNMFW-Work` 16 · `…Int-V_POLIS` 14 · `…Int-T_STORAGE_IMAGE` 10 · `…Work-ClaimTreaty` 7 · sisanya ≤6.

**Graf pemanggilan** (sumber: `Struktur_Flow_TreatyIn.xlsx`): 261 simpul, 387 edge, akar tunggal `Flow_TreatyIn`. **Orphan = 0** — setiap rule terpakai.

> **Diperiksa ulang 18 September 2026.** Angka 261 itu **benar dan bukan kekurangan**: ekspor memuat **279 jalur berkas** tetapi hanya **261 nama berkas berbeda** — 17 nama muncul di lebih dari satu folder, umumnya pasangan Harness dan Section dengan nama sama (`EditXOLAlokasi`, `Hitung_Test`, `OutstandingClaim`, `InputAcceptation`, `AdjustmentDetailNP`, dan 12 lainnya). Silang-periksa nama berkas pohon terhadap isi ekspor: **nol di satu sisi, nol di sisi lain.** Orphan = 0 bertahan.
>
> Kedudukan berkas ini tetap seperti ditetapkan `FINDING-003`: **turunan, bukan sumber.** Ia boleh dipakai untuk menyempitkan pertanyaan dan untuk mencacah struktur; klaim tentang **perilaku** tetap harus berdiri di atas XML. Butir E3 dan E4 sudah dibasiskan ulang ke XML atas aturan itu.

---

## 2. Struktur Klaim — EVIDENCED

### 2.1 PageList langsung di bawah `.ClaimData` (calon tabel anak)

| PageList | Ref | Isi |
|---|---:|---|
| `SpreadingRisk` | 164 | Hasil alokasi per Layer (termasuk baris Retensi Cedant bertanda `TreatyName="UR"`) |
| `ReceiverClaim` | 101 | Rekening penerima pembayaran |
| `AdjustmentList` | 48 | Transaksi Adjustment |
| `ListClaimAmount` | 37 | Nilai kerugian per mata uang |
| `CNPSpreadLoss` | 23 | Pembagian kerugian per porsi |
| `InterestList` | 19 | Objek pertanggungan |
| `SpreadingClaim` | 19 | Penyebaran ke penerima berikutnya |
| `ListTotalEstimation` | 17 | Rekap total per mata uang |
| `EstimationList` | 14 | Estimasi awal |
| `ReinstatementList` | 13 | Perhitungan Premi Pemulihan |
| `SpreadingAdjustment` | 13 | Agregasi penyebaran (XOL) |
| `SpreadingAdjustmentQS` | 12 | Agregasi penyebaran (quota share) |
| `SpreadingBreakQS` | 11 | Rincian per tipe reasuransi |
| `TotalInterestInsured` | 7 | Rekap nilai pertanggungan |
| `Attachment` | 7 | Lampiran |
| `ObjectList` | 7 | Objek |
| `ListClaimAcceptation` | 4 | Rekap akseptasi |
| `ClaimComitee` | 2 | Komite level klaim (jalur tutup/tolak) |

**18 PageList.** Selebihnya (± 100 properti) skalar — daftar lengkap ada di hasil sapuan.

### 2.2 Tanggal Kejadian bersifat tetap

Sapuan seluruh 114 activity: **tidak ada satu pun `Property-Set` yang menulis `.DateOfLoss`.** Nilai hanya masuk lewat input layar (`Section\OutstandingClaim.xml`, kontrol `pxDateTime`, `REQ`).
→ Mendukung penetapan satu klaim = satu Tanggal Kejadian.

### 2.3 Natural key — EXTERNAL

**TIDAK DITEMUKAN DI XML**: tidak ada rule yang mencegah duplikat kombinasi `PolicyNo` + `DateOfLoss` + `IDMaster`.
`CheckDateDOL_Act` memang menjalankan `CekHistoryClaimNonProp_SQL` dan `CariHistoryClaim_SQL` untuk **menampilkan riwayat klaim atas polis yang sama**, dan `RejectedClaim_RD` untuk mengecek klaim yang pernah ditolak — tetapi keduanya bersifat informasi, bukan penolakan. *Disimpulkan dari ketiadaan rule penolakan, bukan dari adanya rule.*
→ Apakah duplikat dicegah secara prosedural oleh manusia: **EXTERNAL**.

---

### 2.4 Retensi Cedant adalah **baris di dalam `SpreadingRisk`**, bukan entitas terpisah — EVIDENCED

`CountLossAllocation_act` menambahkan Retensi Cedant ke PageList yang sama dengan Layer:

```
pyWorkPage.ClaimData.SpreadingRisk(<APPEND>).TreatyType      = "UR"
pyWorkPage.ClaimData.SpreadingRisk(<LAST>).TreatyName        = "UR"
pyWorkPage.ClaimData.SpreadingRisk(<LAST>).ClaimEstimation   = Local.UR
pyWorkPage.ClaimData.SpreadingRisk(<LAST>).ClaimAmountAdjust = Local.UR
```
`Activity\CountLossAllocation_act.xml` baris 7015–7120.

Konsekuensinya: **setiap loop atas `SpreadingRisk` ikut melihat baris Retensi Cedant**, kecuali loop itu menyaringnya sendiri. Dua ejaan penanda dipakai bergantian — `.TreatyName=="UR"` dan `.TreatyType=="UR"` — dan keduanya memang ditulis pada baris yang sama, jadi bukan kontradiksi.

Ada ejaan ketiga: `"Previously Calculated UR"` (`GenerateCACNP_Act` baris 7223, 7425). Nilai ini **tidak pernah ditulis** di folder ini, hanya dibaca — asal-usulnya **TIDAK DITEMUKAN DI XML**.

**Rule yang menyentuh `SpreadingRisk`, dihitung ulang 18 September 2026 — KOREKSI**

Tabel versi pertama mencantumkan "16–66 rujukan `SpreadingRisk`". **Angka itu salah.** Ia menghitung seluruh kemunculan teks `SpreadingRisk`, yang sebagian besar adalah **nama class di metadata langkah** (`pyStepsClassName`, `pxRuleClassName`, `pyStepsParentClass`) — bukan rujukan ke PageList.

| Rule | Kemunculan teks | Rujukan PageList **sebenarnya** | Bentuknya |
|---|---|---|---|
| `GenerateCFS_act` | 66 | **5** | menyalin ke `TempDataOutStanding.ClaimData.SpreadingRisk` |
| `CountTotalInsterest_Act` | 32 | **2** | `@SizeOfPropertyList(...)` sebagai penjaga `>0` |
| `CopyOldataCurr_act` | 23 | **2** | `@LengthOfPageList(...)>0` + satu `Property-Set` `repeat=EMBEDDED` |
| `CountSpreadingXOL` | 31 | **0** | — |
| `CountSpreading_Act` | 30 | **0** | — |
| `DeleteAkseptasi_Act` | 28 | **0** | — |
| `CountClaimTNP_Act` | 22 | **0** | — |
| `SetAccoutNo_Act` | 16 | **0** | — |
| `SendEmailKlaim` | 16 | **0** | — |

**Enam dari sembilan rule tidak menyentuh PageList itu sama sekali.** Kehadiran baris `"UR"` tidak dapat memengaruhi mereka, karena mereka tidak membacanya.

Tinggal tiga rule yang benar-benar membacanya, dan ketiganya sudah tuntas:

| Rule | Akibat baris `"UR"` |
|---|---|
| `CountTotalInsterest_Act` | **tidak ada** — cacahan hanya dipakai sebagai penjaga `>0` |
| `CopyOldataCurr_act` | **tidak ada** — `>0` sebagai penjaga; `Property-Set` mengubah `.Currency` pada setiap baris termasuk baris `"UR"`, dan baris itu memang membawa `Currency` sendiri, jadi konsisten |
| `GenerateCFS_act` | **ada** — baris `"UR"` tersalin ke struktur retrosesi bila ia hadir (lihat `FINDING-003` bagian 4.2) |

**Batas klaim yang berlaku sekarang**: satu rule terbukti terdampak bila baris `"UR"` hadir, dua terbukti tidak terdampak, enam tidak membacanya sama sekali. Apakah baris itu hadir saat `GenerateCFS_act` berjalan bergantung pada urutan tindakan pengguna — lihat `FINDING-003` bagian 3.3.


### 2.5 `.IsEditClaim` — ditulis, tidak pernah dibaca — EVIDENCED

`IsEditClaim` adalah `Rule-Obj-Property` milik class `ASM-FW-GISFW-Data-SpreadingRisk` (`pxRuleObjClass=Rule-Obj-Property`), **bukan** `Rule-Obj-When`.

Seluruh kemunculannya di folder ini — lima, tidak lebih:

| Rule | Langkah | Aksi |
|---|---|---|
| `EditXOLAlokasi` | `pyStepsClassName=ASM-FW-GISFW-Data-SpreadingRisk` | `.IsEditClaim = 1` |
| `CountLossAllocation_act` | dua langkah | `SpreadingRisk(<LAST>).IsEditClaim = 0` |

**Tidak ada satu pun pembacaan**: tidak dipakai di `pyStepsPreCondParamsWhen`, tidak di `pyVisible` section mana pun, tidak di Report Definition, tidak di SQL. Bendera ini **hanya ditulis**.

Dua kemungkinan, dan XML folder ini tidak bisa memilih di antaranya: (a) pembacanya ada di modul Komite → `DEFERRED-TO-KOMITE-SESSION`; (b) bendera ini mati. *Disimpulkan dari ketiadaan pembacaan, bukan dari adanya rule.*

---

## 3. Status Klaim — EVIDENCED

### 3.1 `CNPStatusCase` — enum lengkap **di dalam folder ini**

| Nilai | Ditulis oleh | Tahap |
|---|---|---|
| `"INPUT ACCEPTATION CLAIM"` | `InputOutStandingCTNP_PostAct` step 2 | selesai Registrasi → masuk Akseptasi |
| `"COMITEE ACCEPTANCE (DEPT. HEAD)"` | `CreateChildKomiteCNP_Act` step 31 | Adjustment dikirim ke Komite |
| `"COMITEE ACCEPTANCE (DEPT. HEAD)"` | `CreateChildKomiteCloseNP_Act` step 12 | Penutupan/Penolakan dikirim ke Komite |

Hanya **3 penulisan, 2 nilai unik** di folder ini. Nilai `"CLAIM ACCEPTED"` dan `"CLAIM REJECTED"` **tidak ditulis di folder ini** — ditulis modul Komite (DEFERRED-TO-KOMITE-SESSION).

**Tidak ada satu pun rule di folder ini yang membaca `CNPStatusCase`** — ketiga rule di atas hanya menulis. *Disimpulkan dari ketiadaan pembacaan.* Konsumennya kemungkinan UI atau modul Komite.

### 3.2 Status case sebenarnya

| Fakta | Bukti |
|---|---|
| Satu-satunya status akhir | `Flow\Flow_TreatyIn.xml` shape `End1` → `pyWorkStatus = Resolved-Completed` |
| Satu-satunya penetapan status dari activity | `CloseClaimTNonProp` step 9 → `Call ASMForceCaseClose` dengan `WorkStatus = Resolved-Completed` |
| `Resolved-Rejected` / `Resolved-Withdrawn` | **TIDAK DITEMUKAN DI XML** pada folder ini |
| Jalur reopen | **TIDAK DITEMUKAN DI XML**. *Disimpulkan dari ketiadaan rule, bukan dari adanya rule.* |

`Open-By` yang muncul di `ASMForceCaseClose`, `AttachCAPDF_Act`, `AttachRISlipToWork`, `GetBase64Attachment`, `GetDetailPolis_act` adalah bagian nama method `Obj-Open-By-Handle`, **bukan** status.

---

## 4. Dua Jalur Pengakhiran Klaim — EVIDENCED

| | `CloseClaimNP` | `CloseClaimMD` |
|---|---|---|
| Pre-activity | `CloseClaimNP_preAct` | — |
| Section | `Section\CloseClaimNP.xml` | `Section\CloseClaimMD.xml` |
| Activity inti | `CreateChildKomiteCloseNP_Act` | `CloseClaimTNonProp` |
| Lewat Komite? | **Ya** — `pxAddChildWork`, `ChildClass=…Work-KomiteTreatyNonProp`, `FlowName=KomiteTreaty_Flow` | **Tidak** |
| Penjaga | `@hasMessages(myStepPage)` sebelum membuat child | **Ada**: langkah 1–3 memindai `AdjustmentList`; bila ada `.AcceptanceStatus=="0"` (masih menunggu Komite) → pesan `ErrPendAcc` → `EXIT-ACT-FAIL` |
| Simpan outstanding | (di modul Komite) | `SaveDataToOsAkseptasiNP` (step 5) |
| Kronologi | (di modul Komite) | `InsertChronology_DT` (step 6) |
| Update `json_klaim` | (di modul Komite) | `InsertJsonClaimTreaty_act` (step 7) |
| Kirim ke Arasapas | (di modul Komite) | `Connect-REST insertClaimFinalOrClosed_NP` (step 8) |
| Tutup case | (di modul Komite) | `ASMForceCaseClose` → `Resolved-Completed`, `CloseAllSubCases=true` (step 9) |
| Notifikasi | `SendEmailKlaimRejectClose` bila `IsPEGAPROD` | — |

### 4.0 Status: **BYPASS BERSYARAT** — `CANDIDATE-NOT-MIGRATED`

Dugaan awal "bypass tanpa penjaga" saya cabut, lalu koreksinya sendiri saya turunkan lagi karena kelewatan. Yang terbukti dari XML:

**(a) Jalur MD menutup klaim tanpa pernah membuat case Komite.** Pemindaian `Activity\CloseClaimTNonProp.xml`, `Section\CloseClaimMD.xml`, dan `FlowAction\CloseClaimMD.xml` atas kata kunci `pxAddChildWork`, `CreateChildKomite*`, `KomiteTreaty_Flow`, `KomiteTreatyNonProp`, `ComiteeClaim`, `KomiteList`: **nihil di ketiganya**. Jalur ini langsung `ASMForceCaseClose` → `Resolved-Completed`.

**(b) Penjaganya hanya menangkap Adjustment yang sedang di Komite, bukan yang belum pernah ke Komite.** `.AcceptanceStatus` ditulis **satu kali saja di seluruh 279 berkas**: `CreateChildKomiteCNP_Act` langkah 31 → `= 0`. Perpindahan `0` → `1`/`2` tidak ada di folder ini (dilakukan modul Komite, DEFERRED).
Akibatnya nilai `.AcceptanceStatus` yang mungkin ada: **kosong** (belum pernah dikirim), `0` (sedang di Komite), `1`/`2` (sudah diputus).
Penjaga di `CloseClaimTNonProp` langkah 1 hanya menguji `.AcceptanceStatus=="0"`. Adjustment yang **belum pernah dikirim ke Komite** bernilai kosong, bukan `"0"`, sehingga **tidak tertangkap** dan penutupan tetap berjalan.

**(c) Niat penjaganya memang bukan otorisasi.** Pesan galatnya sendiri berbunyi `"Can not close claim, there is adjustment in comitee!"` — ia dirancang mencegah penutupan **saat Komite sedang bersidang**, bukan mensyaratkan persetujuan Komite.

**Kesimpulan yang terbukti**: `CloseClaimMD` adalah **bypass bersyarat** — melewati Komite sepenuhnya, tetapi menolak berjalan bila ada Adjustment yang sedang menggantung di Komite. Keputusan dimigrasi atau tidak menunggu REQ-006 (statistik pemakaian).

### 4.1 Otorisasi pada jalur penutupan

Lihat §5.0 — ketiadaan otorisasi bukan khusus jalur ini, melainkan berlaku di seluruh modul.

---

## 5. Otorisasi — EVIDENCED

### 5.0 TEMUAN UTAMA: tidak ada model otorisasi yang bisa diwarisi

> **Modul ini tidak memiliki otorisasi di tingkat rule sama sekali.**
>
> - `pyPrivilegeName` → **kosong di seluruh 279 berkas**. Tag-nya ada, nilainya tidak pernah diisi. Artinya `pyPrivilegeClass` yang terisi (36× `…Work-ClaimTreatyNonProp`, dst.) menunjuk ke privilege **tanpa nama** — tidak mengikat apa pun.
> - `Rule-Access-When` → **tidak ada satu pun di export**.
> - `pyVisible` → `ALWAYS` pada seluruh elemen UI yang diperiksa, termasuk kedua jalur penutupan klaim.
>
> Seluruh kontrol akses yang benar-benar berjalan bertumpu pada **access group dan routing flow**, dan keduanya **tidak ikut ter-export**.
>
> **Konsekuensi terhadap cakupan pekerjaan**: tidak ada model otorisasi yang dapat diwarisi dari modul ini. Seluruh RBAC sistem baru berstatus **DECIDED** — dirancang dari nol bersama pemilik proses — **bukan EVIDENCED**. Ini menambah pekerjaan yang belum ada di rencana mana pun. Lihat [ADR-0006](./docs/adr/0006-rbac-dirancang-dari-nol.md).

### 5.1 Sisa sinyal yang ada di XML — EVIDENCED + EXTERNAL

| Sinyal | Isi |
|---|---|
| `pyPrivilegeClass` | `…Work-ClaimTreatyNonProp` 36, `…Data-Adjustment` 15, `ASM-FW-GISFW-Data-PolicyTreatyIn` 14, `…Work-ClaimTreaty` 13, **`ASM-FW-GCNMFW-Work-PNC` 9**, `…GCNMFW-Work` 9, `…Data-ObjectItem` 6, `@baseclass` 5 |
| `pyPrivilegeName` | **kosong seluruhnya** |
| `pyRuleSet` | `GCNMFW` 214, `GISFW` 72, `Pega-RULES` 10, `GCNMFWInt` 6, **`worts@` 4**, `Pega-IntegrationArchitect` 1 |
| `pyApplicationName` | `GCNMFW` 6, **`RaniRApp` 1** |
| `pyWorkBasket` / `pyWorkbasketName` | **TIDAK DITEMUKAN DI XML** |

**Label peran yang benar-benar ada di XML** — `DataTransform\InsertChronology_DT.xml`, properti `.IsCedingConfirm`:
`"Claim Admin"` (default) · `"Claim Dept. Head"` (bila pelaku `CHRISTINEANGELINA`) · `"Operational Director"` (bila `Himawan`) · `"Technical Director"` (bila `NANDINA`).
→ Empat peran ini **EVIDENCED**, tetapi pemetaannya ke orang bersifat hardcode (§7).

**Routing assignment** (`Flow\Flow_TreatyIn.xml`):
- `Assignment2` Registrasi → `pyImplementation = ToCurrentOperator`
- `Assignment1` Akseptasi → `pyImplementation = ToWorkbasket`, **nama workbasket tidak ada di XML** → EXTERNAL

Anomali yang perlu dijelaskan: `ASM-FW-GCNMFW-Work-PNC`, ruleset `worts@`, aplikasi `RaniRApp` — muncul di folder ini tetapi tidak punya rule sendiri. → OPEN-QUESTIONS.

---

## 6. Inventaris Aritmetika — EVIDENCED

Sapuan seluruh 114 activity, seluruh `PropertiesValue` yang mengandung operasi aritmetika dan menyentuh properti bernuansa nilai: **673 ekspresi**. Daftar lengkap: [`pengetahuan/arithmetic-inventory.tsv`](./pengetahuan/arithmetic-inventory.tsv) (kolom: rule, step, properti tujuan, ekspresi, skala, berkas).

### 6.1 Temuan utama

**518 dari 673 ekspresi (77%) tidak memakai `@divide` sama sekali** — memakai operator `*`, `/`, `+`, `-` biasa, sehingga **tidak punya skala eksplisit** dan bergantung pada perilaku desimal bawaan Pega. Hanya 155 ekspresi (23%) yang menyatakan skala.

### 6.2 Skala `@divide` yang dipakai

| Skala | Jml | Dipakai di |
|---:|---:|---|
| 20 | 126 | Seluruh rule inti klaim non-proporsional |
| 10 | 42 | `CountClaimTNP_Act`, `CountLossAllocation_act` |
| 4 | 10 | Rule sisi treaty/premi (`ASM-FW-GISFW-Data-PolicyTreatyIn`) |
| 8 | 4 | `SetPPNPPH` — konstanta pajak |
| 5 | 2 | `CountClaimTNP_Act` → `.PctProrateClaim` |
| 2 | 2 | `SetTPLNote_Act` — pembulatan tampilan |
| 0 | 2 | `SendEmailKlaim` — pembulatan tampilan |

### 6.3 Polanya: skala mengikuti **modul**, bukan jenis nilai

| Kelompok | Skala | Rule |
|---|---:|---|
| Inti klaim non-prop | **20** | `AdjClaimCNP_Act` (28), `AdjClaimAmount_Act` (14), `CreateChildKomiteCNP_Act` (13), `SaveCNPLayerList_Act` (13), `SetActualPremium_ACT` (10), `GenerateCFS_act` (9), `SaveAdjustmentToOSAksep_Act_Tes` (5), `GeneratePlaCNP2_Act` (4), `SetCurrency_Act` (2), `GenerateCACNP_Act` (1) |
| Treaty / premi | **4** | `DetailCalculation` (3), `CountNetPremi_act` (2), `CountResult1_Act` (2), `SetValidateInstallment_Act` (2), `CountTotalInsterest_Act` (1) |
| Pajak & fee | **8** | `SetPPNPPH`: `@divide(2.5,100,8)` brokerage, `@divide(102.2,100,8)`, `@divide(2,100,8)` PPh, `@divide(2.2,100,8)` PPN |
| Tampilan | **0 / 2** | `SendEmailKlaim`: `@divide(.TotalSharePersen,1,0)`, `@divide(.SharePercentage,1,0)` · `SetTPLNote_Act`: `@divide(.TPLPct,1,2)`, `@divide(.TPLAmount,1,2)` — semuanya **bagi dengan 1**, murni pemformatan |

**Hanya 2 rule yang mencampur skala**: `CountClaimTNP_Act` (10×28, 20×13, 5×2) dan `CountLossAllocation_act` (10×14, 20×14 — tepat berimbang). Keduanya adalah rule alokasi inti. → OPEN-QUESTIONS.

### 6.4 Tarif hardcoded di `SetPPNPPH` — EVIDENCED

`.BrokerageFee = @divide(2.5,100,8)` · `.BrokerageFeeSebenarnya = @divide(102.2,100,8)` · `.PPHValue = @divide(2,100,8)` · `.PPNValue = @divide(2.2,100,8)`
→ Tarif brokerage, PPh, dan PPN tertanam di kode. Bila tarif pajak berubah, kode harus diubah.

---

## 7. Tambalan Per-Case di Dalam Kode — EVIDENCED

Sapuan seluruh folder atas literal `CLMNP-…`, `pzInsKey` literal, `IDMaster` literal, dan angka mati.
**Hasil: 9 identitas berbeda, tersebar di 4 rule, 14 langkah.** (Dugaan awal hanya 3 — meleset.)

| # | Identitas | Rule | Step | Yang ditimpa / di-bypass |
|---|---|---|---|---|
| 1 | `CLMNP-975` | `AdjClaimCNP_Act` | 11.2 | `.CNPReinstatement = 881928.966808370` bila `.Currency=="IDR"` |
| 2 | `CLMNP-975` | `AdjClaimCNP_Act` | 11.3 | `.CNPReinstatement = 806851161.1895` bila `.Currency=="USD"` |
| 3 | `CLMNP-232` | `AdjClaimCNP_Act` | 4 | `.SpreadingAdjustment(3)` disalin dari `AdjustmentList(6).SpreadingAdjustment(3)` |
| 4 | `CLMNP-232` | `AdjClaimCNP_Act` | 5 | `.SpreadingQuotaShare(5)` disalin dari `AdjustmentList(7).SpreadingQuotaShare(5)` |
| 5 | `CLMNP-232` | `AdjClaimCNP_Act` | 15.5.1–15.5.4 | `Local.AdjusterFee` dipetakan manual per `IdxQS` (1/2→SA(1), 3/4→SA(2), 5→SA(3), 6→SA(4)) |
| 6 | `CLMNP-861` | `AdjClaimCNP_Act` | 10 | `.SpreadingRisk(IdxLastLayer).AdjusterFee = TotalClaim × ClaimPercentage/100` (rumus khusus) |
| 7 | `IDMaster 1000393` & `1000393/R01` | `CountLossAllocation_act` | 17.2.1 | `Local.LayerLimit` memakai urutan pengali `Kurs` yang berbeda |
| 8 | `CLMNP-367` & `CLMNP-382` (via `pzInsKey`) | `CountLossAllocation_act` | 17.2.1 | `Local.URLimit` memakai `Deductible2/Kurs`, bukan `Deductible` |
| 9 | `IDMaster 1001130` & `1001130/R01` | `GetHistoryMasterID_NP` | 4, 5 | `InputSpreading.CARI22` diisi literal master-id |
| 10 | `CLMNP-50` & `CLMNP-232` | `GetHistoryMasterID_NP` 7, `InputOutStandingClmTNP_PreAct` 5 | | menambahkan baris `"FACOUT"` ke `ListLossAllocation` |

Catatan: pada `InputOutStandingClmTNP_PreAct` step 5, `pyStepsDescription` berbunyi `pyWorkPage.pyID=="CLMNP-50"` sementara kondisi eksekusi sebenarnya `pyWorkPage.pyID=="CLMNP-232"` — **deskripsi dan kondisi tidak sinkron**, indikasi tambalan disalin-tempel.

**Angka mati lain yang bukan tambalan per-case** (ambang kewenangan Komite, `CreateChildKomiteCNP_Act` step 10 & 26.5): `Local.LimitMax = 30000000.00`, `Local.LimitMaxDivHead = 50000000.00`, `Param.LIMIT_BOTTOM = 25000001`.

---


### 7.1 Kelas tambalan kedua: **identitas ORANG**, bukan hanya nomor case — EVIDENCED

Sapuan ini sebelumnya hanya mencari `CLMNP-…` dan `IDMaster`. Sapuan ulang atas seluruh `pyStepsPreCondParamsWhen` menemukan kelas yang terlewat: kondisi yang menyebut **nama operator dan nama orang**.

| Identitas | Rule | Langkah | Bentuk kondisi |
|---|---|---|---|
| `VINCENTVERNANDO_1` | `CreateChildKomiteCNP_Act` | 4 langkah | `pyWorkPage.pxCreateOperator=="VINCENTVERNANDO_1"` |
| `VINCENTVERNANDO_1` | `SaveDataToOSAksep_Act` | 1 | idem |
| `VINCENTVERNANDO_1` | `SaveToOS` | 1 | idem |
| `VINCENTVERNANDO_1` | `SendEmailKlaimRejectClose` | 2 | idem |
| `Himawan`, `NANDINA`, `CHRISTINEANGELINA`, `CHRISTOPMARHASAK` | `CreateChildKomiteCNP_Act` | 4 langkah | `.KomiteID=="<nama>"` |
| `Himawan`, `Nandina C`, `Christine Angelina Hutagalung` | `SethistoryKlaimTreaty` | 3 langkah | `.PICSuggest=="<nama>"` |
| `Himawan` | `SethistoryKlaimTreaty` | 4 langkah | `@contains(.CommentSuggest,"Accepted by Himawan")` / `"Rejected by Himawan"` |
| `Himawan` | `DataTransform\InsertChronology_DT` | 1 | `.PICSuggest=="Himawan"` |

Yang terakhir paling rapuh: keputusan alur ditentukan dengan **mencari potongan teks nama orang di dalam kolom komentar bebas**. Satu salah ketik dari pengguna mengubah jalur.

**Total gabungan per-case + per-orang: 29 langkah, 18 ekspresi unik, 8 rule.**

### 7.2 Umur tambalan — EVIDENCED

Diambil dari `pxCreateDateTime`/`pxCreateOperator` pada blok langkah yang bersangkutan.

| Tanggal | Penulis | Tambalan |
|---|---|---|
| 2018-01-16 | `YOSUAAMBIKA` | `VINCENTVERNANDO_1` di `CreateChildKomiteCNP_Act` (4 langkah) |
| 2020-02-18 | `MESDISILITONGA` | 4 nama `KomiteID` |
| 2020-02-27 | `REVIRUNDUPADANG` | `CLMNP-50`, dan `CLMNP-232` di `InputOutStandingClmTNP_PreAct` |
| 2022-01-03 | `GABRIELAMILITIA` | `VINCENTVERNANDO_1` disalin ke `SaveDataToOSAksep_Act` dan `SaveToOS` |
| 2022-12-01 | `AnanSosmita` | 3 nama `PICSuggest` |
| 2023-08-30 | `AnanSosmita` | `CLMNP-232` di `AdjClaimCNP_Act` (6 langkah), `CLMNP-861` |
| 2023-11-20 | `ArlexyVarian` | `VINCENTVERNANDO_1` disalin ke `SendEmailKlaimRejectClose` |
| 2025-07-04 | `JEFRIHARI` | `IDMaster 1001130` dan `1001130/R01` |
| 2026-07-16 | `JEFRIHARI` | `CLMNP-975` (2 langkah) |

Tiga hal yang dibuktikan deret ini:

1. **`CLMNP-232` ditambal dua kali, berjarak 3,5 tahun, oleh dua orang berbeda, di dua rule berbeda** (2020-02-27 `REVIRUNDUPADANG`; 2023-08-30 `AnanSosmita`). Tambalan kedua tidak menggantikan yang pertama — keduanya masih aktif.
2. **Tambalan tertua berumur lebih dari 8 tahun** dan masih hidup, serta **disalin ke rule lain dua kali sesudahnya** (2022, 2023). Polanya menyebar, bukan mengendap.
3. Yang terbaru berjarak **dua bulan dari hari ini**. Praktik ini **belum berhenti**.

### 7.3 Angka & alamat mati lain — EVIDENCED

| Nilai | Lokasi | Keterangan |
|---|---|---|
| `http://192.168.105.116:80/prweb` | `SaveToOS`, `SendEmailKlaim`, `SendEmailKlaimRejectClose` | alamat IP server tertanam di kode |
| `http://pega.nusantarare.com:80/prweb` | rule yang sama | alamat kedua, berdampingan dengan yang pertama |
| `"10026"` | `SetCurrency_Act` baris 1798, 2498 | `Param.IDCurr=="10026"` — id mata uang tertanam |
| `"EDITCLAIMXOL"`, `"CountXOL"`, `"AccEdit"`, `"ChildCase"` | tersebar | penanda alur berupa string bebas |

---

### 7.4 `CLMNP-232` ditambal dua kali — dua gejala berbeda, bukan satu yang bertabrakan — EVIDENCED

Pertanyaan yang belum terjawab ronde lalu: apakah kedua tambalan menyentuh properti yang sama, sehingga yang menang bergantung urutan pemanggilan? **Tidak.** Keduanya menyentuh halaman yang berlainan.

| | Tambalan 2020-02-27 (`REVIRUNDUPADANG`) | Tambalan 2023-08-30 (`AnanSosmita`) |
|---|---|---|
| Rule | `InputOutStandingClmTNP_PreAct` | `AdjClaimCNP_Act` |
| Halaman yang disentuh | `ListLossAllocation.pxResults`, `ListTreaty.pxResults` | `Primary.SpreadingAdjustment`, `.SpreadingQuotaShare`, `Local.*` |
| Yang dilakukan | menyisipkan baris `"FACOUT"` dan `"OR"` ke daftar **masukan** | memetakan ulang nilai pada struktur **keluaran** |
| Jenis | menambah baris | menimpa nilai |

Contoh isi tambalan 2023:
```
.SpreadingAdjustment(3) = pyWorkPage.ClaimData.AdjustmentList(6).SpreadingAdjustment(3)
.SpreadingQuotaShare(5) = pyWorkPage.ClaimData.AdjustmentList(7).SpreadingQuotaShare(5)
Local.AdjusterFee       = Primary.SpreadingAdjustment(1..4).AdjusterFee   (dipetakan manual per indeks)
```

**Kesimpulan: tidak ada yang "menang".** Keduanya berjalan, pada tahap yang berbeda, atas halaman yang berbeda. Artinya pada satu klaim yang sama terdapat **dua gejala terpisah** — satu pada daftar masukan alokasi, satu pada hasil pemetaan Adjustment.

**Hipotesis yang belum diuji** (ditandai supaya tidak terbaca sebagai fakta): tambalan 2023 memperbaiki akibat dari tambalan 2020 — daftar masukan yang disisipi baris tambahan menghasilkan indeks yang bergeser, lalu indeks itu dipetakan ulang secara manual tiga tahun kemudian. Menguji ini memerlukan penelusuran pengaruh `ListLossAllocation` terhadap penomoran `SpreadingAdjustment`, dan **itu belum dikerjakan**.

Catatan terpisah: pada `InputOutStandingClmTNP_PreAct` langkah 5, `pyStepsDescription` berbunyi `pyWorkPage.pyID=="CLMNP-50"` sementara kondisi eksekusinya `pyWorkPage.pyID=="CLMNP-232"` — deskripsi dan kondisi tidak sinkron.

### 7.5 Laju tambalan — EVIDENCED

Sumber: `pxCreateDateTime` / `pxCreateOperator` pada blok langkah yang bersangkutan.

| Tanggal | Penulis | Rule | Jenis | Langkah |
|---|---|---|---|---|
| 2018-01-16 | `YOSUAAMBIKA` | `CreateChildKomiteCNP_Act` | identitas akun | 4 |
| 2020-02-18 | `MESDISILITONGA` | `CreateChildKomiteCNP_Act` | identitas orang | 4 |
| 2020-02-27 | `REVIRUNDUPADANG` | `GetHistoryMasterID_NP` | nomor case | 1 |
| 2020-02-27 | `REVIRUNDUPADANG` | `InputOutStandingClmTNP_PreAct` | nomor case | 2 |
| 2022-01-03 | `GABRIELAMILITIA` | `SaveDataToOSAksep_Act` | identitas akun *(salinan)* | 1 |
| 2022-01-03 | `GABRIELAMILITIA` | `SaveToOS` | identitas akun *(salinan)* | 1 |
| 2022-12-01 | `AnanSosmita` | `SethistoryKlaimTreaty` | identitas orang | 3 |
| 2023-08-30 | `AnanSosmita` | `AdjClaimCNP_Act` | nomor case | 7 |
| 2023-11-20 | `ArlexyVarian` | `SendEmailKlaimRejectClose` | identitas akun *(salinan)* | 2 |
| 2025-07-04 | `JEFRIHARI` | `GetHistoryMasterID_NP` | ID master | 2 |
| 2026-07-16 | `JEFRIHARI` | `AdjClaimCNP_Act` | nomor case | 2 |

**Langkah baru per tahun:**

```
2018  ####                4
2019                      0
2020  #######             7
2021                      0
2022  ####                4
2023  #########           9
2024                      0
2025  ##                  2
2026  ##                  2   (per 16 Juli; ekspor diambil September)
```

Yang dapat dibaca dari deret ini:

1. **Tidak melambat sampai 2023, lalu turun.** Puncaknya 2023 (9 langkah, 2 penulis). Sesudahnya 2 langkah per tahun. Laju turun, tetapi **tidak berhenti**.
2. **Enam penulis berbeda dalam delapan tahun.** Tidak ada satu orang pun yang memegang pola ini; ia diwariskan.
3. **Menyebar, bukan mengendap.** Kondisi `pxCreateOperator=="VINCENTVERNANDO_1"` ditulis sekali pada 2018 lalu **disalin ke tiga rule lain** pada 2022 dan 2023 oleh dua orang yang berbeda dari penulis aslinya.
4. **Yang terbaru berjarak dua bulan dari hari ini.** Selama migrasi berjalan, tambalan baru masih mungkin bertambah ke sistem yang sedang dipetakan.

**Batas klaim:** deret ini dibaca dari satu ekspor bertanggal September 2026. Tambalan yang **dihapus** sebelum tanggal itu tidak meninggalkan jejak dan tidak terhitung di sini. Jadi angka ini adalah **batas bawah**, bukan hitungan lengkap. *Disimpulkan dari isi ekspor, bukan dari riwayat perubahan rule.*

---

## 8. Boundary Contract Draft — Claim → Komite (EVIDENCED)

Titik kopling yang memaksa kedua modul menjadi satu unit. Detail sisi Komite: **DEFERRED-TO-KOMITE-SESSION**.

### 8.1 Pembuatan child case

| Pemicu | Rule | Parameter |
|---|---|---|
| Kirim Adjustment ke Komite | `CreateChildKomiteCNP_Act` step 29 | `pxAddChildWork`: `ChildClass=ASM-FW-GCNMFW-Work-KomiteTreatyNonProp`, `FlowName=KomiteTreaty_Flow`, `CopyPageData=true`, `Commit=true`, `UpdateHistory=true` |
| Kirim Penutupan/Penolakan ke Komite | `CreateChildKomiteCloseNP_Act` step 10 | idem |

### 8.2 Data yang dikirim Claim → Komite

Disalin ke `ChildWorkPage` sebelum `pxAddChildWork`: `KomiteCount=1`, `KomiteLoop`, `CLMNO`, `TransferType`, seluruh `Adjustment` (`Page-Copy` dari `AdjustmentList(idx)`), `Adjustment.CNPLayerList[].CNPCurrencyList[]`, `KomiteList[]`, `Komite.CircumtansesCouseOfLoss`, `Komite.Remarks`, `Komite.Occupation`, `QuotationData.*`, `AttachmentStream` (Base64).

### 8.3 Data yang ditunggu balik Komite → Claim

Ditulis balik oleh modul Komite ke `AdjustmentList(IndexObject)`: `AcceptanceStatus`, `AcceptedNo`, `AcceptedDate`, `ComiteeClaim[].KomiteAproval`, `.KomiteComment`, `.DateApproval`; ditambah `CNPStatusCase` dan `IsSubjectivity`/`SubjectivityNote` pada level klaim.

### 8.4 Kunci relasi

| Arah | Kunci |
|---|---|
| Claim → Komite | `pyWorkPage.pxCoveredInsKeys(<LAST>)`; `.KomiteNo = @substring(pxCoveredInsKeys(<LAST>),19,30)` |
| Komite → Claim | `pxCoverInsKey` |
| Komite → Adjustment tertentu | `Adjustment.IndexObject` = **posisi numerik** di `AdjustmentList` — bukan surrogate key |

#### 8.4a Ketiganya pengenal sistem lama — tidak satu pun ikut pindah sebagai kunci

*Ditambahkan 18 September 2026. Ini satu masalah, bukan tiga.*

`IndexObject` sudah ditandai di atas: **posisi numerik**, bukan surrogate key. Posisi berubah bila daftar disisipi atau diurut ulang, jadi ia bukan pengenal sama sekali — ia alamat.

`pzInsKey` dan `pxCoverInsKey` punya cacat yang **berbeda tetapi sekelas**. Keduanya **pegangan instance bawaan Pega**, bukan pengenal bisnis. `pzInsKey` berbentuk `<class tabel> <pyID>`; `pxCoverInsKey` menyimpan `pzInsKey` induk pada baris anak. Keduanya menamai *di mana barang itu tersimpan di Pega*, bukan *barang apa itu*. ADR-0021 menetapkan pengenal dinamai menurut isinya, jadi ketiganya sama-sama gugur sebagai kunci sistem baru — nasib yang sama dengan `CASEID`.

**Tetapi `pzInsKey` berguna justru karena ia pegangan instance.** Untuk `MIGRASI_KORELASI` ia **pengenal sisi lama terbaik yang kita punya**, dan alasannya berbukti: `PYID` tidak punya unique constraint di tabel work, dan `POOLDATA.JSON_KLAIM` mewajibkan `MNK_NO_KLAIM` ada tetapi berkunci primer `IDPEGA` — nomor klaim lama **tidak dipaksa unik oleh struktur mana pun**. Pegangan instance lebih dapat diandalkan daripada keduanya.

| Pengenal lama | Nasibnya | Alasan |
|---|---|---|
| `pzInsKey` | **tabel korelasi saja**, berumur terbatas | pegangan instance; paling andal di sisi lama, tidak bermakna di sisi baru |
| `pxCoverInsKey` | **tabel korelasi saja**, berumur terbatas | menyimpan `pzInsKey` induk; relasi induk-anak baru memakai kunci sendiri |
| `IndexObject` | **tidak dipakai sama sekali** | posisi, bukan pengenal — tidak stabil bahkan di sisi lama |
| `CASEID` | **tabel korelasi saja** | salinan `pzInsKey` di sisi basis data |

Garis yang berlaku untuk keempatnya: **tidak pernah masuk skema kanonik.** Relasi induk-anak di sistem baru memakai kuncinya sendiri, dan jembatan migrasi dibongkar setelah rekonsiliasi selesai.


### 8.4b Penghubung di tingkat data — EVIDENCED

Sumber: `pengetahuan/ddl/TABLE_PC_ASM_FW_GCNMFW_WORK.sql`. Tabel ini **tidak dimigrasi** — ia buatan Pega, dan sistem baru punya model sendiri. Yang diambil hanya apa yang tampak.

#### Batas antar modul adalah nilai satu kolom, bukan batas tabel

`pzInsKey` di Pega berbentuk `<class tabel> <pyID>`. Literal di XML memperlihatkan **tiga jenis case memakai nama class yang sama**:

```
ASM-FW-GCNMFW-WORK CLMNP-...     Claim Non Prop
ASM-FW-GCNMFW-WORK KMTNP-...     Komite Non Prop
ASM-FW-GCNMFW-WORK PNC-...       PNC
```

Ketiganya memakai class induk `ASM-FW-GCNMFW-Work` — yaitu tabel `PC_ASM_FW_GCNMFW_WORK`. Subclass `ClaimTreatyNonProp` (3.691 rujukan), `ClaimTreaty` (1.391), `PNC` (322), dan `KomiteTreatyNonProp` (56) **tidak punya tabel sendiri**.

Dikuatkan indeksnya: `PC_WORK_IDX1_2AED8 (PXCOVERINSKEY, PXOBJCLASS, PYSTATUSWORK)` menempatkan `PXOBJCLASS` tepat di tengah pola akses induk-anak.

> **Claim dan Komite duduk di baris-baris tabel yang sama, dibedakan `PXOBJCLASS`.** Batas antara dua modul bukan batas tabel, melainkan nilai satu kolom.

#### Penghubung induk-anak

| Kolom | Tipe | Peran |
|---|---|---|
| `PXCOVERINSKEY` | `VARCHAR2(255 CHAR)` | menyimpan `pzInsKey` case induk pada baris case anak |
| `PXCOVEREDCOUNT` | `NUMBER(18)` | cacah anak |
| `PXCOVEREDCOUNTOPEN` | `NUMBER(18)` | cacah anak yang masih terbuka |
| `PXCOVEREDCOUNTUNSATISFIED` | `NUMBER(18)` | cacah anak yang belum tuntas |

Di sisi XML, `pxAddChildWork` dipakai 6 kali dan `pxCoveredInsKeys` 8 kali. Jadi relasi Claim ke Komite yang selama ini hanya terbaca sebagai pemanggilan rule **punya wujud kolomnya**: `PXCOVERINSKEY` pada baris Komite berisi `pzInsKey` klaim induknya.

#### `PYID` tidak dijamin unik

| | |
|---|---|
| Primary key | `PZINSKEY` (`VARCHAR2(255 CHAR) NOT NULL`) |
| `PYID` | `VARCHAR2(32 CHAR)` — **tanpa unique constraint, tanpa index tersendiri** |

Nomor klaim `CLMNP-...` tidak dijamin unik oleh basis data. Pola yang sama dengan `OS_AKSEPTASI_KLAIM`. Diukur REQ-018.

#### `BYTE` versus `CHAR` — campur, bahkan di dalam 19 kolom itu sendiri

Kolom bawaan Pega seluruhnya `CHAR`. Sembilan belas kolom yang di-expose **tidak seragam**:

| Semantik | Kolom |
|---|---|
| `BYTE` (14) | `SOBNAME`, `BUSINESSNAME`, `CEDINGCONAME`, `INSUREDNAME`, `NOPOLIS`, `MASTERID`, `DATEOFLOSS`, `ISSUBJECTIVITY`, `PLA_CEDING`, `PLA_SOB`, `DLA_CEDING`, `DLA_SOB`, `POLICY_CEDING`, `CAUSEOFLOSS` |
| `CHAR` (3) | `NOCLAIM`, `FLAGONGOINGCOMMITTE`, `CLMNO` |
| bukan teks (2) | `KOMITECOUNT` `NUMBER(18)`, `DATERECEIVED` `DATE` |

Untuk data non-ASCII, `VARCHAR2(n BYTE)` menampung lebih sedikit karakter daripada yang terlihat. Paling terpapar: `INSUREDNAME` (255 BYTE), `CEDINGCONAME` (128 BYTE), dan empat kolom `*_CEDING`/`*_SOB` (500 BYTE). Masuk risiko migrasi data; sistem baru memakai satu semantik saja.

#### `PZPVSTREAM` dan hak akses

`PZPVSTREAM BLOB`, `SECUREFILE ... ENABLE STORAGE IN ROW`. Di luar 19 kolom itu dan kolom bawaan Pega, **tidak ada satu pun kolom nilai uang** — menegaskan bagian 13.5.

```
GRANT ALTER, DELETE, INDEX, INSERT, REFERENCES, SELECT, UPDATE, ... TO POOLDATA
```

`POOLDATA` memiliki `INSERT`, `UPDATE`, `DELETE`, bahkan `ALTER` atas tabel work Pega. Premis REQ-021 terbukti dari DDL; yang belum diketahui hanya apakah hak itu dipakai.

### 8.5 Properti bersama

**Tambahan — kopling yang belum tercatat: `.IsEditClaim`.**
`EditXOLAlokasi` menulis `1`, `CountLossAllocation_act` menulis `0`, dan **tidak ada pembaca di sisi Claim** — bertahan setelah jaring diperlebar 30 → 45 tag dan setelah 480 baris Java dibaca sebagai kode (D38). Bila pembacanya ada di modul Komite, ini kopling lintas modul berarah Claim → Komite.

**Ditutup 19 September 2026 sebagai hipotesis permanen.** Folder Komite tertutup untuk pekerjaan ini, jadi cabangnya tidak akan dipilih. Entrinya masuk kontrak batas **bertanda belum terselesaikan, dengan arah dan pemilik dikosongkan** (§8.7), dan `.IsEditClaim` sendiri masuk golongan **penanda tak berpemilik**: dimigrasi apa adanya, ditandai, tidak diberi perilaku.


`KomiteTreatyNonProp` (4 rule) · `pxAddChildWork` (2) · `KomiteList`/`KomiteCount`/`KomiteLoop` (2) · `KomiteNo` (5) · `ComiteeClaim` (4) · `AcceptanceStatus` (**11 rule** — properti paling terkopel) · `IsKomite` (2) · `KomiteAproval` (4) · `KomiteComment` (3) · `KomitePost` (4) · `TotalKomite` (2) · `DataCommitteeTreaty` (5) · `EMAILKOMITE` (3).

### 8.6 Penentu jenjang Komite

`ReportDefinition\FilterEmailKomiteWithLimit.xml` atas class `ASM-FW-GCNMFW-Int-EMAILKOMITE`.
Kolom: `ID, NAME, EMAIL, LIMIT_BOTTOM, LIMIT_TOP, DEGREE, STS_AKTIF, OPERATOR_ID, STS_REJECT, STS_ADJ, TYPE_BUSINESS, TYPE_KOMITE, STS_REG, STS_SURVEY, STS_SALVAGE, STS_ADJUSTER, STS_KLAIM, JABATAN`.
Filter `A AND C AND B`: `LIMIT_BOTTOM = Param.LIMIT_BOTTOM`, `STS_KLAIM = Param.STS_KLAIM`, `STS_AKTIF = "1"`.
Jumlah baris hasil → `ChildWorkPage.KomiteLoop` = jumlah jenjang persetujuan.

### 8.7 Syarat penerimaan di batas — ditetapkan 19 September 2026

*Tiga syarat, lahir dari penutupan C4, C5, dan C6. Bagian sebelumnya memerikan **apa yang menyeberang hari ini**; bagian ini menetapkan **dengan syarat apa sistem baru menerimanya**.*

Folder `Komite Claim Non Prop` **tertutup untuk pekerjaan ini**. Ketiga butir karena itu tidak ditutup dengan jawaban tentang isi modul seberang — ia ditutup dengan **batas**: apa yang kita terima, dan apa yang kita tolak.

#### Syarat 1 — pihak luar merujuk Adjustment dengan **kunci stabil**, bukan posisi (menutup C4)

`IndexObject` adalah **posisi numerik di dalam `AdjustmentList`**, bukan pengenal (§8.4a). Posisi berubah bila daftar disisipi atau diurut ulang.

**`IndexObject` tidak ada di sistem baru.** Penggantinya **kunci surrogate yang stabil pada baris Adjustment**, dan kunci itulah yang menyeberang. Pihak luar yang mengembalikan `AcceptanceStatus`, `AcceptedNo`, atau `AcceptedDate` menyebut Adjustment dengan kunci itu; **rujukan berdasarkan posisi ditolak di batas.**

Pertanyaan lama — *bagaimana `AcceptanceStatus` dikembalikan ke `AdjustmentList(IndexObject)`* — karena itu tidak dijawab melainkan **diganti**: ia bertanya bagaimana meniru sesuatu yang tidak dimigrasi.

#### Syarat 2 — `AlokasiXOLPaid` adalah **medan masuk**, tidak pernah diturunkan sendiri (menutup C5)

Dicari menyeluruh pada empat lapisan, pada enam tag yang sempat terlewat saat jaring diperlebar 30 → 45, dan di 480 baris Java tertanam: **nol penulis** (D11, D38). Yang ada hanya sasaran perulangan, pembacaan selisih, penjaga panjang, dan PageList layar.

Itu cukup untuk merancang, dan inilah bentuknya: `AlokasiXOLPaid` **nilai yang datang dari luar**, boleh kosong, dan **mesin kita tidak pernah menghitungnya sendiri**. Bila kelak ternyata ada penulis di modul seberang, bentuk ini sudah benar; bila ternyata tidak ada penulis sama sekali, kolomnya kosong dan tidak ada yang rugi.

#### Syarat 3 — `.ValueAdjustment` **wajib datang bersama mata uang dan kursnya** (menutup C6 dan H1/H2)

`.ValueAdjustment` muncul di **tepat satu berkas** (`CreateChildKomiteCNP_Act`), hanya dibaca. Asal nilainya **TIDAK DITEMUKAN DI XML**, dan satuannya karena itu tidak diketahui — itu seluruh isi FINDING-001.

**ADR-0029 menutupnya sebagai syarat penerimaan, bukan sebagai tebakan.** Nilai yang menyeberang membawa nilai asli, kode mata uang, dan kurs yang dipakai. **Nilai tanpa mata uang ditolak di batas** — tidak diterima lalu dianggap rupiah, dan tidak diterima lalu dianggap mata uang asli.

Dengan itu kedua cabang H terpenuhi sekaligus: bila hulu memang sudah mengonversi (H1), ia menyatakannya; bila tidak (H2), ia juga menyatakannya. **Ambang kewenangan dibandingkan terhadap nilai IDR** (ADR-0007), dan perbandingan itu sah dalam kedua keadaan.

#### Entri yang tetap bertanda belum terselesaikan

Satu entri **sengaja ditinggalkan tidak lengkap**, dan itu bukan kelalaian:

| Entri | Keadaan | Kenapa |
|---|---|---|
| `.IsEditClaim` | **arah dan pemilik dikosongkan** | pembacanya ada di modul seberang atau tidak ada sama sekali, dan folder itu tertutup. Arah yang ditebak lebih buruk daripada arah yang kosong |

Perlindungan hasil suntingan **tidak menunggu entri ini**: ADR-0008 menuntutnya terlepas dari apakah sistem lama pernah punya jalurnya.

---

## 9. Java Step — EVIDENCED

17 activity, 26 blok. Dua pola dominan:

| Pola | Rule |
|---|---|
| PDF byte → Base64 via `tools.getParameterPage()` | `HTMLToPDF` (5), `HTMLRISlipToPDF` (3), `AttachRISlipToWork` (2), `GenerateCFS_act` (1), `GetBase64Attachment` (1) |
| Dedup `HashSet<String>` atas PageList mata uang/estimasi | `CountTotalInsterest_Act` (2), `GetHistoryMasterID_NP` (1), `CountSpreading_act` (1), `AdjClaimCNP_Act` (1), `AdjClaimAmount_Act` (1), `AddAkseptasiCNP_Act` (1) |
| Serialisasi page → JSON | `SetValueClaimTNP_Act`, `GetDetailPolis_act`, `GeneratePlaCNP_Act`, `HitServiceToKasir_Act` (2) |
| Deteksi ekstensi berkas → MIME | `InsertGoogleStorage_Act`, `InsertDocument_Act` |
| Penutupan case paksa | `ASMForceCaseClose` (2) — menyalin page via `tools.findPageByHandle` |

---

## 10. Penyimpanan Properti: EXPOSED vs IN-BLOB — EVIDENCED (analisis)

Pega menyimpan properti work object di kolom BLOB `pzPVStream` kecuali properti itu di-*expose* sebagai kolom nyata. Bukti terkuat EXPOSED adalah pemakaian properti sebagai kolom/filter/sort di Report Definition.

**Pemeriksaan 22 Report Definition di folder ini: tidak ada satu pun yang berjalan di atas class work `ASM-FW-GCNMFW-Work-ClaimTreatyNonProp` atau `ASM-FW-GCNMFW-Work`.** Seluruhnya berjalan di atas class `-Int-` (tabel/view Oracle eksternal), `Link-Attachment`, atau `Data-Admin-Operator-ID`.

**Konsekuensi**: tidak ada bukti XML bahwa properti klaim mana pun di-expose sebagai kolom. Dugaan default untuk seluruh `.ClaimData.*` dan `.TreatyInMaster.*` adalah **IN-BLOB**.
→ `DESC` tabel work Pega **tidak akan** menampilkan properti-properti ini. Sumber tipe yang sahih adalah definisi property di schema PegaRULES (REQ-001), bukan `ALL_TAB_COLUMNS`.
→ Nilai klaim yang terlihat di tabel Oracle bisnis (`OS_AKSEPTASI_KLAIM`, `json_klaim`, `claimxol2`) adalah **salinan hasil**, bukan sumber. Presisi riil yang diterima akuntansi justru ada di sana — itu yang harus diprofilkan.

---

## 11. Rekonsiliasi terhadap DDL sistem lama — EVIDENCED

**Sumber**: `pengetahuan/DDL_Script_ClaimNonProp.xls`. **Versi 2026-09-18 10:36 — 48 objek**, seluruhnya owner `POOLDATA`. (Versi sebelumnya 46 objek; dua VIEW ditambahkan, lihat bagian 13.)
**Stempel asal berkas**: dibuat 2026-09-18 03:06:58 oleh `Raynold`, WPS Spreadsheets, format BIFF (OLE Compound Document); diperbarui 2026-09-18 10:36.
**Otoritas**: tertinggi untuk struktur data. Bentuk skripnya `CREATE OR REPLACE EDITIONABLE ...` beserta klausa `SEGMENT CREATION`, `STORAGE(INITIAL ...)`, dan `TABLESPACE` — itu keluaran `DBMS_METADATA.GET_DDL` dari basis data yang berjalan, bukan rancangan yang ditulis orang.

### 11.1 Cakupan

| Jenis | Jumlah |
|---|---|
| TABLE | 31 |
| VIEW | 10 |
| PROCEDURE | 6 |
| FUNCTION | 1 |

**Tidak memuat satu pun objek Pega** (`pc_*`, `pr4_*`, `pr_data_admin`). Karena itu REQ-001 dan REQ-012 **tetap BLOCKER**, dan dugaan `.ClaimData.*` tersimpan IN-BLOB di `pzPVStream` **tetap terbuka**.

### 11.2 Hasil rekonsiliasi

| | Jumlah |
|---|---|
| Naik dari GUESS ke CONFIRMED | **10** |
| Sudah CONFIRMED, kini ber-owner pasti | 36 |
| Ada di daftar saya, **tidak ada di DDL** | 9 |
| Objek baru yang belum pernah saya daftarkan | **11** |
| Koreksi jenis objek (saya salah) | **5** |

**Naik ke CONFIRMED**: `EMAILKOMITE`, `V_POLIS`, `V_M_CAUSE_OF_LOSS`, `V_D_CAUSE_OF_LOSS`, `ADJUSTERCONSULTANT`, `LST_BANK_GROUP`, `CATASTROPHE`, `PROVINCE`, `TREATYGROUP`, `CURRENCYSTANDARD`.

**Koreksi jenis — saya salah lima kali**:

| Objek | Saya daftarkan | Sebenarnya |
|---|---|---|
| `CURRENCY` | TABLE | **VIEW** |
| `CITY` | TABLE | **VIEW** |
| `PROVINCE` | TABLE | **VIEW** |
| `CURRENCYSTANDARD` | TABLE | **VIEW** |
| `GET_TOKEN_STORAGE` | FUNCTION | **PROCEDURE** |

**Prediksi saya yang meleset paling jauh**: saya menulis *"Dugaan saya: objek `V_POLIS` ini TIDAK ADA"*, dengan alasan nama class Pega hanya wadah. **Salah.** `V_POLIS` ada, berupa VIEW dengan definisi sepanjang 4.472 karakter — view terbesar di seluruh DDL. Pelajarannya: ketiadaan sebuah nama di dalam `FROM` tidak berarti objeknya tidak ada; Pega dapat mengakses tabel lewat pemetaan class, tanpa nama objek pernah muncul di SQL mana pun.

**Tidak ada di DDL** (tetap di antrean Oracle, **bukan** disimpulkan mati — DDL ini bisa saja parsial): `TREATY_OUT`, `CLAIMXOL`, `TRLOSS_DETAIL_T` (owner `REINSURANCE`), `V_MST_USER_TEKNIS`, `F_GET_EMAIL` (owner `GL`), `PLATNP_SEQ`, dan seluruh objek Pega.

`TREATY_OUT` yang absen patut diperhatikan: ia BLOCKER, terbaca literal di SQL produksi, tetapi DDL-nya tidak ikut. `CLAIMXOL` absen padahal `CLAIMXOL2` ada.

### 11.3 Objek baru — dirujuk di dalam DDL, tidak pernah terbaca di XML

Sebelas objek yang tidak meninggalkan jejak apa pun di 279 berkas XML, karena hanya disentuh dari dalam stored procedure:

| Objek | Ditemukan di | Kenapa penting |
|---|---|---|
| **`M_CURRENCYSTANDARD`** | `GETCURRENCYSTANDARD`, view `CURRENCYSTANDARD` | **tabel kurs yang sebenarnya**. Seluruh konversi mata uang bertumpu padanya |
| **`TANGGAL_CLOSING`** | `PROC_GENERATE_SEQUENCE_NUMBER` | **tanggal tutup buku**. Menentukan periode produksi sebuah nomor — berkaitan dengan aturan cut-off tanggal 25 di `MEMORI_PEMAHAMAN.MD` |
| `CITYINPUT`, `DISTRICTINPUT`, `RWINPUT` | `GET_TOKEN_STORAGE` | perlu diverifikasi — bisa jadi alias di dalam view, bukan objek tersendiri |
| `M_CAUSE_OF_LOSS`, `D_CAUSE_OF_LOSS` | view `V_*_CAUSE_OF_LOSS` | tabel dasar penyebab kerugian |
| `M_CURRENCY`, `M_PROVINCE` | view `CURRENCY`, `PROVINCE` | tabel dasar |
| `M_NATION` | `GET_TOKEN_STORAGE` | — |
| `BRANCH` | dalam DDL | peran belum terbaca |

Masuk `pengetahuan/PULL-LIST.csv` sebagai **REQ-017**.

**Koreksi atas sapuan pertama saya.** Daftar pertama memuat `V_DAY_CLOSING` dan `V_PERIODE_DATE` sebagai objek baru. **Keduanya bukan objek** — mereka variabel lokal PL/SQL (`v_day_closing NUMBER`, `v_periode_date DATE`) yang terbaca sebagai nama objek karena sapuan saya menormalkan huruf besar. Penggantinya adalah objek yang sebenarnya dibaca di tempat itu: **`POOLDATA.TANGGAL_CLOSING`**.

**Nilai mati di dalam basis data.** `PROC_GENERATE_SEQUENCE_NUMBER` memuat:
```sql
IF TRUNC(v_now) <= TO_DATE('02/01/2026','DD/MM/YYYY') THEN
    v_mm_yyyy := '12.2025';  v_tahun := 2025;
```
Tanggal dan periode tertanam langsung di prosedur. Ini kelas tambalan yang sama dengan yang terdaftar di bagian 7, **tetapi berada di basis data, bukan di Pega** — jadi sapuan XML tidak akan pernah menemukannya. Sapuan berkala menurut ADR-0012 perlu mencakup DDL, bukan hanya XML.

### 11.4 Yang TIDAK tertutup oleh DDL

1. **Presisi riil tetap terbuka.** Dari 120 kolom `NUMBER`, **88 tidak menyatakan presisi maupun scale sama sekali** — `NUMBER` polos, yang di Oracle berarti sampai 38 digit signifikan dengan scale apa pun yang datang. Basis data **tidak memaksakan pembulatan apa pun**. Apa pun hasil `@divide(...,20)` dari Pega tersimpan apa adanya. **REQ-005 tetap hidup; ADR-0003 tetap draft.**

   Pengecualiannya satu-satunya: **`TREATYINPRODUCTION` memakai `NUMBER(20,4)` pada 23 kolom nilai** (`PREMI_OGP`, `CLAIM`, `NET_PREMIUM`, `PPHVALUE`, `PPNVALUE`, `BALANCE_DUE_TO`, dst.). Jadi ada **satu** tabel yang menetapkan 4 desimal, dan seluruh sisanya bebas. Itu menguatkan temuan §6.3 bahwa skala mengikuti modul, bukan jenis nilai.

2. **Properti IN-BLOB justru terbukti dengan tidak munculnya.** DDL ini tidak memuat tabel work Pega sama sekali, jadi ketiadaan `.ClaimData.*` di sini **bukan** bukti propertinya tidak ada — ia konsisten dengan dugaan IN-BLOB, bukan membantahnya.

3. **Pemetaan properti Pega ke kolom tetap belum terjawab.** DDL memberi sisi kolom; XML memberi sisi properti; jembatannya tetap `Data-Admin-DB-Table` — **REQ-012, masih BLOCKER**.

### 11.5 T1, T2, T3 — terjawab sebagian tanpa menjalankan PREFLIGHT

| | Jawaban | Dari mana |
|---|---|---|
| **T2** tipe `DATA_JSON` | **`CLOB` dengan `CHECK (DATA_JSON IS JSON)`**, disimpan `SECUREFILE ... ENABLE STORAGE IN ROW`. Berlaku untuk `JSON_KLAIM` dan `OS_AKSEPTASI_KLAIM` | DDL langsung |
| **T1** versi Oracle | **12c ke atas** — `IS JSON` dan notasi titik `a.data_json.TypeLoss` keduanya butuh 12c+; `EDITIONABLE` butuh 12.1+ | disimpulkan dari sintaks yang dipakai |
| **T3** owner | **`POOLDATA`** untuk ke-46 objek ini | DDL langsung |

`pengetahuan/PREFLIGHT.sql` tetap perlu dijalankan untuk **T4** (hak akses), **T5** (letak PegaRULES), **T7** (ukuran tabel), dan untuk memastikan 9 objek yang tidak ada di DDL.

---

## 12. Kenapa sistem ini butuh tambalan per-case — satu rantai sebab akibat

Tiga temuan yang selama ini terpisah ternyata satu cerita. Ditulis di sini supaya tidak dibaca sebagai tiga keluhan yang berdiri sendiri.

**Sebab.** Tidak ada satu pun activity yang memanggil `CountClaimTNP_Act`; ia dipicu dari Section lewat `<pyActivity>`, 42 rujukan di tiga berkas. Urutan jalannya `CountClaimTNP_Act`, `CountLossAllocation_act`, `GenerateCFS_act`, dan `SaveToOS` **tidak ditetapkan di kode mana pun** — ia ditentukan oleh kontrol mana yang ditekan pengguna (§bagian 3 `FINDING-003`).

**Akibat pertama — hasil tidak dapat direproduksi.** Dua pengguna dengan masukan sama tetapi urutan tindakan berbeda dapat menghasilkan keadaan berbeda. `CountLossAllocation_act` meng-`<APPEND>` baris `"UR"` dan tidak ada penghapusan daftar sebelum itu; `PEGA_JSON_OS_AKSEP_KLAIMTNP` selalu `INSERT` dan tidak pernah `UPDATE` (logika upsert-nya dikomentari). Keduanya menumpuk, bukan mengulang.

**Akibat kedua — selalu ada kasus yang keluar jalur.** Karena keadaan akhir bergantung jalan yang ditempuh, akan selalu ada klaim yang berakhir di keadaan yang tidak diantisipasi rule mana pun.

**Akibat ketiga — satu-satunya perbaikan yang murah adalah menambal klaim itu satu per satu.** Memperbaiki penyebabnya berarti mengubah arsitektur pemicuan; menambal satu klaim berarti menambah satu `pyStepsPreCondParamsWhen`. Yang kedua selesai dalam satu jam.

**Bukti bahwa rantai ini nyata, bukan teori**: 29 langkah tambalan, 18 ekspresi, 8 rule, 6 penulis, sepanjang **delapan tahun** (2018-01-16 s/d 2026-07-16), dengan laju yang naik sampai 2023 lalu turun tetapi **tidak berhenti** (§7.5).

**Karena itu ketiga keputusan berikut harus diambil bersama, bukan sendiri-sendiri**: ADR-0011 (hitung ulang idempoten), ADR-0012 (deteksi tambalan lewat sapuan), dan keputusan Q22 (perhitungan sebagai turunan data, bukan akibat tombol). Mengambil satu tanpa dua lainnya akan mengembalikan pola yang sama dalam bentuk baru.

---

## 13. Pembaruan DDL 2026-09-18 10:36 — dan satu berkas yang belum sampai

### 13.1 Yang bertambah: dua VIEW

| Objek | Menutup |
|---|---|
| `POOLDATA.CLAIMXOL` | sebelumnya terdaftar "tidak ada di DDL" — **ternyata VIEW, bukan TABLE**. Koreksi keenam atas jenis objek. |
| `POOLDATA.V_MST_USER_TEKNIS` | sebelumnya GUESS, sekarang **CONFIRMED** dan ber-owner pasti |

### 13.2 `CLAIMXOL` membuka struktur JSON — REQ-009 terjawab sebagian

```sql
FROM pooldata.os_akseptasi_klaim OS, json_table (
     data_json, '$.CNPLayerList[*]'
     columns( XOL varchar2 path '$.XOL',
       nested path '$.CNPCurrencyList[*]' columns(
         TotalXOLGross    varchar2 path '$.TotalXOLGross',
         CNPReinstatement varchar2 path '$.CNPReinstatement',
         Currency         varchar2 path '$.Currency',
         KursIDR          varchar2 path '$.KursIDR')))
WHERE SUBSTR(CASEID,1,25) = 'ASM-FW-GCNMFW-WORK CLMNP-'
```

Tiga hal terbaca langsung:

1. **Bentuk JSON**: `$.CNPLayerList[*]` bersarang `$.CNPCurrencyList[*]`. Jadi layer adalah tingkat luar dan mata uang tingkat dalam — bukan sebaliknya. Itu memperkuat ADR-0024 (satu akseptasi per klaim per layer per mata uang) dan menetapkan arah sarangnya — buktinya dipindahkan ke badan ADR itu.
2. **Seluruh nilai uang di dalam JSON diekstrak sebagai `varchar2`** — `TotalXOLGross`, `CNPReinstatement`, `KursIDR`. Angka uang di dalam blob **berupa teks**. Presisinya adalah apa pun yang tertulis di teks itu, tanpa batas dan tanpa jaminan.
3. **Bentuk `CASEID`**: `'<NAMA CLASS> <pyID>'`, yaitu `ASM-FW-GCNMFW-WORK CLMNP-…`. Class yang dipakai adalah **class induk** `ASM-FW-GCNMFW-Work`, bukan `ClaimTreatyNonProp`.

### 13.3 Objek baru lagi, dari dalam dua view ini

| Objek | Ditemukan di | Catatan |
|---|---|---|
| `POOLDATA.OS_AKSEPTASI_SUBJECTIVITY` | `CLAIMXOL`, cabang `UNION ALL` dengan `STS_SUBJECTIVITY = '1'` | tabel akseptasi kedua, khusus Subjectivity. Belum pernah terlihat di XML |
| `POOLDATA.MST_USER_TEKNIS` | `V_MST_USER_TEKNIS` | tabel dasar; kolomnya sebagian di dalam `JSON_DATA` |
| **`hrdasm.v_hrd_mst@asmd.sinarmas.co.id`** | `V_MST_USER_TEKNIS` | **database link ke sistem HRD**. Menjawab sebagian T5b: db link memang dipakai di lingkungan ini |

### 13.4 Berkas yang belum sampai: DDL tabel work Pega

`DATAPEGA.PC_ASM_FW_GCNMFW_WORK` **tidak ada di dalam berkas `.xls` versi 2026-09-18 10:36.** Saya baca ulang seluruh 48 baris × 4 kolom dan menyapu seluruh isi sel untuk `DATAPEGA`, `PC_ASM`, dan `PZPVSTREAM`; satu-satunya kecocokan adalah nama parameter `Datapega IN CLOB` pada dua prosedur. Tidak ada berkas baru lain di folder artefak.

Sampai berkas itu sampai, yang berikut **belum dapat dikerjakan** dan tidak boleh dianggap selesai: pemeriksaan apakah Claim dan Komite berbagi satu tabel; `PXCOVERINSKEY` sebagai penghubung induk-anak; ketiadaan unique constraint pada `PYID`; pola `BYTE` versus `CHAR`; `ISSUBJECTIVITY VARCHAR2(5)`; dan pengisian baris tabel work ke `pengetahuan/SCHEMA-ACTUAL.csv`.

Yang **sudah** dapat dikerjakan dari keterangan yang disampaikan langsung, dan sudah dikerjakan: konsekuensi IN-BLOB (bagian 13.5), format `DATEOFLOSS` (`FINDING-005`), dan `FLAGONGOINGCOMMITTE`/`KOMITECOUNT` (bagian 13.6).

### 13.5 `.ClaimData.*` IN-BLOB — status naik dari dugaan menjadi **terbukti**

Dari sembilan belas kolom yang di-expose di luar kolom bawaan Pega, **tidak satu pun kolom nilai uang**. Tidak ada `ClaimAmount`, `Currency`, `Kurs`, `ValueAdjustment`, dan tidak ada satu pun anggota `SpreadingRisk` maupun `AdjustmentList`.

Tiga konsekuensi yang harus dinyatakan terang-terangan:

1. **REQ-005 tidak dapat dijalankan atas properti klaim.** Profil presisi lewat SQL mensyaratkan kolom; nilai-nilai itu tidak punya kolom. `ALL_TAB_COLUMNS` tidak akan pernah memberi tipe atau presisi `ClaimAmount`, berapa kali pun dijalankan. **REQ-005 dipersempit menjadi hanya tabel bisnis `POOLDATA`.**
2. **REQ-001 menjadi satu-satunya jalur** untuk tipe properti klaim. Bila `PR4_*` tidak dapat diakses, satu-satunya cadangan adalah membongkar `PZPVSTREAM`. Rencananya disiapkan sekarang, tidak menunggu REQ-001 gagal — lihat `TABLE-EXTRACTION-REQUEST.md`.
3. **ADR-0003 tetap draft, tetapi alasannya berubah.** Bukan karena datanya belum ditarik, melainkan karena **basis data memang tidak menyimpan presisi untuk nilai-nilai ini sama sekali**. Diperkuat oleh 13.2: bahkan ketika nilai itu dikeluarkan dari JSON lewat view, tipenya `varchar2`.

### 13.6 `FLAGONGOINGCOMMITTE` dan `KOMITECOUNT` — sebagian yang saya sebut "tidak dapat diukur" ternyata dapat

Sapuan 279 berkas:

| Kolom | Kemunculan di folder Claim |
|---|---|
| `FLAGONGOINGCOMMITTE` | **nol** |
| `KOMITECOUNT` | 2, keduanya `ChildWorkPage.KomiteCount` — ditulis pada case anak, bukan pada klaim |

Artinya keduanya **ditulis dari modul Komite**, tetapi karena keduanya punya kolom nyata, keduanya **dapat dihitung lewat SQL tanpa membuka modul Komite**. Itu membuka dua hal yang sebelumnya saya tandai tidak terukur:

- **Q3 / `CloseClaimMD`** — `FLAGONGOINGCOMMITTE` adalah penanda yang dicari: berapa klaim ditutup ketika penanda itu masih menyala.
- **Q9 / berapa kali Adjustment dikirim ke Komite** — `KOMITECOUNT` menjawabnya langsung.

Masuk antrean sebagai **REQ-019**. Maknanya tetap `DEFERRED-TO-KOMITE-SESSION`; yang berubah adalah jumlahnya kini terukur.

---

## 14. Catatan metodologis — setiap lapisan sumber punya titik buta sendiri

Ditulis karena polanya sudah berulang tiga kali, bukan karena satu kejadian.

**Analisis XML saja punya titik buta yang sistematis.** Tiga hal berikut menggerakkan sistem ini dan **tidak meninggalkan jejak apa pun di 279 berkas XML**:

| Yang tidak terlihat dari XML | Terlihat setelah lapisan apa masuk |
|---|---|
| Sebelas objek basis data yang hanya disentuh dari dalam stored procedure — termasuk `M_CURRENCYSTANDARD`, tabel kurs yang menopang seluruh konversi | DDL |
| Tanggal dan periode tertanam di `PROC_GENERATE_SEQUENCE_NUMBER` (`TO_DATE('02/01/2026')` → `'12.2025'`) — kelas tambalan yang sama dengan 29 langkah di bagian 7, tetapi di basis data | DDL |
| Ketergantungan lintas basis data: `hrdasm.v_hrd_mst@asmd.sinarmas.co.id` | DDL |

Dan sebaliknya, **DDL punya titik butanya sendiri**: seluruh nilai uang klaim tidak punya kolom, sehingga `ALL_TAB_COLUMNS` tidak akan pernah menyebutkannya — keberadaannya hanya terbaca dari XML.

Kesimpulan yang mengikat cara kerja, bukan sekadar catatan:

> **Tidak ada satu sumber pun yang memperlihatkan sistem ini secara utuh.** Setiap lapisan baru yang masuk membuka jalur yang lapisan sebelumnya tidak dapat melihatnya. Karena itu tidak boleh ada pernyataan berbentuk "tidak ada X di sistem ini" — yang sah hanya "tidak ada X di lapisan yang sudah saya baca".

**Konsekuensi untuk ADR-0012**: sapuan berkala harus mencakup **tiga** hal, bukan satu — ekspor XML, DDL, dan **daftar database link**. Yang ketiga baru ditambahkan setelah `V_MST_USER_TEKNIS` terbaca; sebelum itu tidak ada yang tahu perlu menyapunya.

Lapisan yang **belum** masuk sama sekali dan titik butanya masih penuh: tabel work Pega (`DATAPEGA`), isi `PR4_*`, pemetaan `Data-Admin-DB-Table`, dan isi blob `PZPVSTREAM`.

---

## 15. Konvensi boolean — diperiksa, dan hasilnya bersih

Dugaan yang diuji: adakah properti yang **ditulis** dengan satu konvensi boolean lalu **dibaca** dengan konvensi lain? Itu akan menjadi kasus kedua dari pola `CheckDateDOL_Act` — dua bagian sistem berselisih tentang bentuk data.

**Hasil: nol bentrok.** Setiap properti taat asas pada dirinya sendiri.

| Konvensi | Properti |
|---|---|
| Angka (`==1` / `==0`) | `IsOutstanding` (73 pembacaan), `IsCFS`, `IsKomite`, `IsAcceptation`, `IsPLA`, `IsPayment`, `IsSaveToOs`, `Flagkomite`, `FlagErr`, `FlagCurrency`, `FlagError` |
| Teks (`=="true"`) | `IsTPL`, `FlagActualPremium`, `IsSubjectivity` |

Perpecahannya **antar** properti, bukan **di dalam** satu properti. Jadi ia soal keseragaman rancangan, bukan bug.

**Batas pemeriksaan**: pendeteksian penulisan hanya menangkap pasangan `<PropertiesName>` dan `<PropertiesValue>` yang berdampingan di Activity. Penulisan lewat Data Transform (`pyPropertiesName`) dan lewat layar tidak ikut tersapu. Beberapa properti — `IsTreatyIn`, `IsReject`, `IsCloseFile`, `IsAnyAcceptation` — terbaca tetapi penulisnya tidak tertangkap, jadi untuk keempatnya pernyataan "bersih" **belum berlaku**.

Untuk sistem baru: satu konvensi saja, tipe boolean sejati. `ISSUBJECTIVITY VARCHAR2(5)` di basis data lama menjadi contoh yang tidak diwarisi.

---

## 16. `V_POLIS` — bukan gabungan banyak tabel, melainkan proyeksi satu dokumen JSON

Objek terbesar di DDL (4.472 karakter) ternyata sederhana bentuknya: **26 kolom, satu tabel sumber (`json_polis`), 32 pemanggilan `JSON_VALUE`, nol `JOIN`, nol `UNION`.**

Kolomnya: `RN`, `CASEID`, `POLICYNO`, `STARTDATETIME`, `ENDDATETIME`, `SOURCEOFBUSINESS(NAME)`, `CEDINGCO(NAME)`, `BUSINESSCODE`, `BUSINESSNAME`, `CUSTOMERNAME`, `SOURCEOFBUSINESSROOT`, `PRODUCTIONDATE`, `PREMI`, `DISCOUNT`, `TSI`, `TYPE`, `PRODKE`, `BUSINESSTYPE`, `MARKETINGCODE`, `CLIENTIDCON`, `CLIENTIDORG`, `QQ`, `SELECTRENEWAL`, `RNWSTATUS`.

Tiga hal yang mengubah gambaran:

1. **"Polis" di sistem ini tidak punya tabel sendiri.** Ia dokumen JSON di `json_polis`, dan `V_POLIS` hanyalah cara memandangnya sebagai baris dan kolom. `CASEID` pun berasal dari dalam JSON (`a.DATA_JSON.IDNewBisnis`).
2. **Satu tabel menampung dua bentuk dokumen** — Facultative dan Treaty — dibedakan `QuotationData.BusinessFac`, dengan jalur JSON yang berbeda untuk data yang sama. Lihat `FINDING-005` bagian 6.1.
3. **Seluruh kolom tanggalnya adalah teks hasil sambungan `SUBSTR`**, bukan `DATE`. Inilah asal format `dd/MM/yyyy` yang bertabrakan dengan `.DateOfLoss` di `FINDING-005`.

Yang ini menutup: **14 RDB List pada class `ASM-FW-GCNMFW-Int-V_POLIS` membaca proyeksi JSON, bukan tabel relasional.** Dugaan awal saya bahwa `V_POLIS` tidak ada sebagai objek memang salah, tetapi alasan di baliknya — bahwa class itu tidak menunjuk tabel bisnis sungguhan — ternyata setengah benar: ia menunjuk view atas satu dokumen.

---

## 17. Tidak ada pembatas antara yang dicoba dan yang dipakai

Lima hal yang selama ini terdaftar sebagai anomali terpisah (E3, E4, E6, E7, E8, E9) ternyata satu temuan. Ditulis sebagai satu bagian karena memperlakukannya sebagai lima keluhan membuat sebabnya tidak terlihat.

### 17.1 Artefak uji berjalan di alur produksi — dan dapat dicapai pengguna

| Artefak | Ukuran | Jalur masuknya |
|---|---|---|
| `Activity\SaveAdjustmentToOSAksep_Act_Tes` | — | dipanggil `Call` dari `CreateChildKomiteCNP_Act` langkah 25 |
| `Section\Hitung_Test` | 187 KB | **`<pyInclude>Hitung_Test</pyInclude>` di dalam `Section\InputAcceptation.xml`** |
| `Harness\Hitung_Test` | 240 KB | harness tersendiri, `pyHarnessName = Hitung_Test` |

Pertanyaan yang diajukan: terdaftar saja, atau benar-benar dapat dicapai dari layar? **Dapat dicapai.** `Hitung_Test` di-`include` ke dalam Section produksi `InputAcceptation` — bukan sekadar terdaftar di suatu tempat. Dan `SaveAdjustmentToOSAksep_Act_Tes` dipanggil lewat `Call` dari activity produksi, bukan dari harness uji.

Keduanya membawa stempel kompilasi 2024 (`..._Stream_20240229T095743_632_GMT`), jadi bukan sisa lama yang tertinggal — mereka masih dibangun ulang bersama rule lain.

### 17.2 Sambungannya dengan tambalan

| | Artefak uji | Tambalan per-case |
|---|---|---|
| Cara masuk | di-`include` / di-`Call` dari alur produksi | satu `pyStepsPreCondParamsWhen` |
| Umur | stempel 2024, masih terkompilasi | tertua 2018, terbaru 2026-07-16 |
| Yang menghentikan | **tidak ada** | **tidak ada** |

**Satu sebab, dua gejala: tidak ada mekanisme yang memisahkan yang dicoba dari yang dipakai.** Sekali sesuatu masuk ke ruleset produksi, tidak ada apa pun yang mengeluarkannya kembali — baik ia percobaan maupun tambalan. Delapan tahun, enam penulis, dan tidak ada satu pun peristiwa pembersihan yang meninggalkan jejak.

Bukti ketiga dari pola yang sama, kali ini di basis data: `PROC_GENERATE_SEQUENCE_NUMBER` memuat `IF TRUNC(v_now) <= TO_DATE('02/01/2026') THEN v_mm_yyyy := '12.2025'` — nilai mati yang sudah lewat tanggalnya dan tetap ada.

### 17.3 Akibatnya untuk ADR-0012

Pembekuan tambalan yang hanya mengurus `pyStepsPreCondParamsWhen` **akan bocor lewat pintu yang sama**. Artefak uji masuk lewat `pyInclude` dan `Call` — dua jalur yang tidak disentuh sapuan tambalan.

Sapuan berkala perlu menambah satu hal lagi: **artefak yang namanya mengandung penanda percobaan (`Test`, `Tes`, `Tmp`, `Coba`, `Backup`, `Old`) tetapi dirujuk dari alur produksi.** Itu murah diperiksa dan langsung memperlihatkan kebocoran.

---

## 18. Perilaku sistem berbeda menurut node yang mengeksekusi

Dimensi yang tidak pernah masuk analisis selama enam ronde, ditemukan lewat sapuan `pxSystemNodeID`.

| Yang dicari | Hasil |
|---|---|
| Rule yang bercabang pada `pxSystemNodeID` | **satu**: `When\IsPEGASyariah.xml` → `pxProcess.pxSystemNodeID = "jboss1074"` |
| Rule yang memakai When itu | **satu**: `Activity\HitServiceToKasir_Act.xml`, 3 langkah |
| Yang berubah | `TempKasir.CARI15 = "100115"` — kode yang dikirim ke sistem kasir |
| `jboss117` (41x), `jboss122117` (36x) | **bukan percabangan** — seluruhnya `<pxHostId>`, metadata audit |

Cakupannya sempit: satu rule, satu field, jalur setoran ke kasir. Tetapi **kode setoran yang berbeda menunjuk ke entitas pembukuan yang berbeda**, bukan sekadar server yang berbeda — dan itu membalik kesimpulan pertama saya tentang "syariah". Uraiannya di `_selesai/OPEN-QUESTIONS.md` G1, ditulis sebagai dua kemungkinan sejajar dan **tidak** saya putuskan sendiri.

**Batas yang penting untuk cakupan seluruh analisis ini**: karena percabangannya dievaluasi saat berjalan di dalam ruleset yang sama, ekspor 279 berkas berlaku untuk kedua node. Yang tidak dapat dibuktikan dari sini: apakah ada instance Pega lain dengan ruleset yang berbeda sama sekali. Itu `_selesai/OPEN-QUESTIONS.md` A15.

---

## 19. Rekonsiliasi kolom DDL terhadap properti Pega

Dua daftar terpisah, seperti diminta. Hasil lengkapnya di `pengetahuan/rekonsiliasi-kolom-vs-properti.tsv`.

| | Jumlah |
|---|---|
| Kolom unik di seluruh DDL (31 tabel + 8 view) | 408 |
| Properti Pega unik di 279 XML (di luar `px*`/`py*`/`pz*`) | 659 |
| **Berpasangan** (nama sama) | **110** |
| **Kolom tanpa pasangan properti** | **298** |
| **Properti tanpa pasangan kolom** | **549** |

### 19.1 Delapan puluh tiga persen properti tidak punya kolom di mana pun

549 dari 659 properti tidak ditemukan sebagai nama kolom di seluruh DDL. Itu bukan anomali — itu **ukuran kuantitatif dari temuan IN-BLOB** di bagian 13.5, dan mencakup hampir seluruh model klaim: `ADJUSTMENTLIST`, `ACCEPTANCESTATUS`, `ADJCLAIMVALUE`, `ALOKASIIDR`, `AMOUNTIDR`, dan seterusnya.

**Peringatan tafsir**: pencocokan ini berdasarkan **nama**, bukan pemetaan sungguhan. Pemetaan yang sah hanya datang dari `Data-Admin-DB-Table` (REQ-012). Nama yang sama belum tentu kolom yang sama, dan nama yang berbeda belum tentu bukan. Angka di atas adalah **batas atas** jumlah properti IN-BLOB, bukan hitungan pasti.

### 19.2 298 kolom yang tidak disentuh Pega sama sekali

Kolom seperti `BALANCE_BEFORE_PPH`, `BALANCE_DUE_TO`, `BROKERAGE`, `CESSIONVALUE`, `COINS_MIN`/`COINS_MAX`, `CLAIMPAYMENTTYPE` ada di tabel bisnis tetapi tidak pernah muncul sebagai nama properti di folder Claim.

Tiga kemungkinan, dan XML tidak dapat memilih di antaranya:

1. Diisi oleh **modul lain** (Produksi, Akuntansi, Komite).
2. Diisi oleh **proses basis data** — dan itu bertemu dengan REQ-023 (~60 kolom skalar `OS_AKSEPTASI_KLAIM` yang tidak diisi prosedurnya) serta REQ-021 (jalur tulis dari luar aplikasi).
3. Mati.

Yang perlu ditegaskan: **kolom-kolom ini tetap ikut termigrasi bila tabelnya dimigrasi.** Memutuskan mana yang dibuang memerlukan jawaban REQ-021 dan REQ-023 lebih dulu — tidak boleh disimpulkan dari ketiadaannya di folder Claim saja, karena lingkup sesi ini memang satu modul.
