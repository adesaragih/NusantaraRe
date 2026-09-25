# Telusur Flow — Komite Claim Life

STEP D2, Tahap 2 konteks #1. Ditelusur 2026-09-13. Konvensi: `_METHOD.md`.

**Titik masuk:** `Komite Claim Life/Flow/KomiteLife_Flow.xml`
→ `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITELIFE_FLOW`, tipe **`Rule-Obj-Flow` dari `<pxObjClass>`**
(file ini tidak punya `<pzOriginalInstanceKey>` — lihat `../inventory/_summary.md` §19.4),
`<pyStartActivity>Start2`.

Modul terkecil ber-`Flow` di domain Komite: 47 file.

---

## 1. Diagram alur

**1 Assignment, 1 Decision, 4 connector** `[terverifikasi]` — flow paling kecil yang ditelusur
sejauh ini di seluruh D2.

```
Start2 ─(Always)─> Assignment1 "KomiteRouter"  [WorkList, pyRouteTo=Custom]
                        │ FlowAction: ViewTransferDtl
                        v
                    Decision1 "KomiteLoop"
                        ├─ When IsKomiteLoop ─> Assignment1   ← LOOP kembali ke assignment yang sama
                        └─ Else ─────────────> End1
```

### 1.1 Tangga persetujuan **tidak** dimodelkan sebagai shape — ia data

`[terverifikasi]` Ini temuan struktural utama konteks ini. Graf hanya punya **satu** Assignment yang
di-loop, bukan rangkaian shape per tingkat persetujuan.

Guard loop `Komite Claim Life/When/IsKomiteLoop.xml`
(`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `ISKOMITELOOP` / `RULE-OBJ-WHEN`):

```
.AcceptStatus = "1"   DAN   .KomiteCount <= .KomiteLoop
```

- `.KomiteLoop` = **jumlah tingkat** tangga (batas atas)
- `.KomiteCount` = **tingkat yang sedang berjalan** (pencacah)
- `.AcceptStatus` = hasil keputusan tingkat itu

Artinya: **berapa tingkat tangga dan siapa yang menyetujui ditentukan data, bukan struktur flow.**
Nilai `.KomiteLoop` tidak di-set di modul ini → **OQ-032**.

### 1.2 Routing

`[terverifikasi]` `<pyImplementation>WorkList`, `<pyRouteTo>Custom`, router
`<pyImplementation>KomiteRouter`.

`Komite Claim Life/Activity/KomiteRouter.xml`
(`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEROUTER` / `RULE-OBJ-ACTIVITY`) berisi **satu langkah**:

| Step | Method | Precondition | Aksi |
| ---: | --- | --- | --- |
| 1 | `Property-Set` | `Primary.TransferType == '2'` | `param.AssignTo = .KomiteID` |

**Sasaran penugasan adalah nilai data `.KomiteID`**, dan hanya di-set bila `TransferType == '2'`.

`[pertanyaan terbuka]` Apa yang terjadi bila `TransferType != '2'` — `param.AssignTo` tidak di-set
sama sekali. Arti `TransferType = '2'` **belum terverifikasi** (OQ-020). → **OQ-033**.

Modul ini memakai `WorkList` (antrean per-pengguna), sama dengan Endorsement Life dan PremiumList
Life, berbeda dari treaty inward yang memakai `WorkBasket` → OQ-028.

---

## 2. Status / state yang berubah

Seluruhnya dari `Komite Claim Life/Activity/KomitePostAdjustment.xml`
(`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`, **36 step**)
kecuali disebut lain. Rule ini **tidak** terdaftar di register OQ-011.

### 2.1 Jejak persetujuan ditulis per tingkat `[terverifikasi]`

Tiga `Property-Set` menulis ke **elemen list** yang diindeks pencacah tangga:

| Properti | Nilai |
| --- | --- |
| `...AdjustmentList(IndexAdjustment).KomiteList(Local.Komite).KomiteAproval` | `pyWorkPage.AcceptStatus` |
| `...KomiteList(Local.Komite).KomiteComment` | `pyWorkPage.Comment` |
| `...KomiteList(Local.Komite).DateApprove` | `@CurrentDateTime()` |

dengan `Local.Komite = pyWorkPage.KomiteCount`. Hal yang sama juga ditulis ke
`pyWorkPage.KomiteList(Local.Komite)`.

**Jadi setiap tingkat tangga menghasilkan satu entri berisi keputusan, komentar, dan waktu.**

### 2.2 Kode dan nilai literal

| Properti | Nilai literal | Peran | Arti |
| --- | --- | --- | --- |
| `AcceptStatus` | `1`, `2` | menentukan lanjut/berhenti tangga; ditulis ke `KomiteAproval` | **belum terverifikasi** |
| `TransferType` | `'2'` | gerbang routing (§1.2) | **belum terverifikasi** (OQ-020) |
| `PolicyDataLife.Type` | `"QP"`/`"QR"` vs `"TP"`/`"TR"` | memilih **dua RDB-List berbeda** untuk penomoran akseptasi | **belum terverifikasi** (OQ-020) |
| `STS_REJECT` | `"0"` | bagian kondisi `.ACCEPTEDNO=="" && .IsCheck = true && .STS_REJECT == "0"` | **belum terverifikasi** |
| `ACCEPTEDNO` | `""` | penanda belum bernomor akseptasi | — |
| `RetroID` | `"1000013"` | gerbang `Call InsertJsonClaimLife_Act` | **belum terverifikasi** — pola sama dengan OQ-031 |
| `ProdDateTime` | `< "20250207T000000.000 GMT"` | gerbang 2 `Property-Set` | **tanggal cutover ter-hardcode** — §2.3 |
| `Local.currentdate` | `@CurrentDate("dd","Asia/Jakarta")` | dipakai bersama `GETTanggalClosing_SQL` | pola ambang tanggal (OQ-030) |

### 2.3 Tanggal cutover ter-hardcode `[terverifikasi]`

Dua `Property-Set` digerbangi:

```
TempOpenPage.PolicyDataLife.OfferFacIn.PolicyData.ProdDateTime < "20250207T000000.000 GMT"
```

Sebuah **tanggal absolut (7 Februari 2025)** tertanam di dalam rule sebagai pembeda perilaku lama
dan baru. `[dugaan]` ini penanda perubahan aturan bisnis pada tanggal itu; **belum terverifikasi**
→ **OQ-034**.

### 2.4 Titik akhir tangga `[terverifikasi]`

Dua langkah hanya berjalan **di tingkat terakhir** (`KomiteCount == KomiteLoop`):

| Precondition | Langkah |
| --- | --- |
| `AcceptStatus = 1 && KomiteCount == KomiteLoop` | `Call LoadDocumentLife_ACT` |
| `AcceptStatus == 2 && KomiteCount == KomiteLoop` | `RDB-List` |

Jadi `AcceptStatus` **1** dan **2** memicu jalur akhir yang berbeda. Mana yang "setuju" dan mana
"tolak" **belum terverifikasi** — tidak disimpulkan dari nama.

---

## 3. Objek Oracle yang disentuh

Tujuh pemanggilan `RDB-List` dari `KomitePostAdjustment`:

| Class | RequestType | Precondition |
| --- | --- | --- |
| `ASM-FW-GISFW-Int-policyjson` | `GETTanggalClosing_SQL` | — |
| `ASM-FW-GISFW-Int-LIFE_PREMIUM_DETAIL` | `Generate_NoAccept_KMT_Life` | `Type=="QP"\|\|"QR"` |
| `ASM-FW-GISFW-Int-LIFE_PREMIUM_DETAIL` | `Generate_NoAccept_KMT_LifeRetro` | `Type=="TP"\|\|"TR"` |
| `ASM-FW-GISFW-Int-policyjson` | `GetKodeProdLife_SQL` | `ACCEPTEDNO == ""` |
| `ASM-FW-GISFW-Int-policyjson` | `GetSequenceNumber_SQL` | `ACCEPTEDNO == ""` |
| `ASM-FW-GISFW-Int-LIFE_PREMIUM_DETAIL` | `UpdateOsAkseptasiClaimLife_sql` | `.ACCEPTEDNO=="" && .IsCheck = true && .STS_REJECT == "0"` |
| `ASM-FW-GISFW-Int-LIFE_PREMIUM_DETAIL` | `UpdateOsAkseptasiClaimLife_sql` | `AcceptStatus==2 && KomiteCount == KomiteLoop` |

`[terverifikasi]` Kode `QP/QR` dan `TP/TR` memilih **dua rule penomoran akseptasi yang berbeda** —
`Generate_NoAccept_KMT_Life` versus `..._LifeRetro`. Ini pemakaian kode OQ-020 yang paling
konsekuensial sejauh ini: ia menentukan **skema penomoran** mana yang dipakai.

### 3.1 Tabel akseptasi klaim — **dapat dibaca langsung** `[terverifikasi]`

`Komite Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` adalah blok PL/SQL yang melakukan
`INSERT INTO POOLDATA.OS_AKSEPTASI_KLAIM_LIFE (...)` diikuti `COMMIT` — **bukan** pemanggilan
stored procedure. Isinya terbaca, sehingga **bukan** batas pengetahuan.

**55 kolom** yang ditulis:

```
CASEID NO_CLAIM POLICY_NO POLICY_HOLDER CERTIFICATE_NO NAME_OF_INSURED SEX DOB AGE PLAN
BEGIN_DATE LAPSE_DATE EXPIRED_DATE STATUS EM_PERCENT CURRENCY SUM_INSURED CEDING_RETENTION
SUM_REASURED SHARE_NUSANTARA_RE CLAIM_AMOUNT WPC PL_NUMBER DISEASE ICD_CODE NOTES CEDINGCO
CEDINGCONAME SOB SOBNAME BUSINESSID BUSINESSNAME SHARE_RETRO KETERANGAN ACCEPTATION_DATE
NO_ACCEPTATION STS_REJECT CLAIM_RETRO RETROID RETRONAME SECURITYREINSURERID SECURITYREINSURER
TYPECEDING TYPE CONFIRMATION_DATE CLAIM_RECEIVED_DATE COMPLETE_DATE ID NAME_OF_BANK IDBANK
ACCOUNTNO CREATEOPNAME PRODUCTNAMEID PRODUCTNAME RETROCEDED_SHARE
```

Perintah audit:
`awk '/<pyBrowseSQL>/,/<\/pyBrowseSQL>/' "Komite Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml"`

**Ini fragmen skema nyata** — meringankan sebagian OQ-001 untuk tabel ini. **Tipe kolom tetap tidak
diketahui.** Kolom `ICD_CODE` `[dugaan]` kode diagnosis medis; `WPC` **kepanjangan belum
terverifikasi**. `COMMIT` berada di dalam blok → OQ-013.

---

## 4. Integrasi eksternal

`[terverifikasi]` Tiga pemanggilan berturut-turut, **seluruhnya digerbangi `IsPEGAPROD`**:

| Step | Call | Precondition |
| --- | --- | --- |
| `Call serviceInsertArasapasClaimLife_act` | service eksternal | `IsPEGAPROD` |
| `Call SendEmailKlaimLife` | notifikasi email | `IsPEGAPROD` |
| `Call HitServiceToKasirKMTLife_Act` | service "kasir" | `IsPEGAPROD` |

Pola identik dengan PremiumList Life (D2 Tahap 1): **seluruh jalur keluar berada di belakang satu
gerbang yang kondisinya tidak terbaca**.

**`IsPEGAPROD` tidak ditelusur isinya** — `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN` terdaftar di
register **OQ-011 entri #305 (14 modul, 2 isi berbeda)**, dan kondisinya **tidak terbaca dari tag**
(`<pyLabel>` hanya template kosong). **Cabang mana yang aktif tidak ditebak** → OQ-029.

---

## 5. Batas pengetahuan

### 5.1 Dua target panggilan tidak ada di modul ini

`[terverifikasi]` `KomitePostAdjustment` memanggil dua activity yang **tidak ada** di
`Komite Claim Life`:

| Call | Satu-satunya salinan di korpus | Status register OQ-011 |
| --- | --- | --- |
| `UpdateWorkObject` | `Endorsment Fac In/Activity/UpdateWorkObject.xml` | tidak terdaftar |
| `serviceInsertArasapasClaimLife_act` | `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml` | tidak terdaftar |

Keduanya **tidak berkonflik**, jadi ini bukan blocker seperti OQ-025 — tetapi resolusinya
bergantung pada pewarisan class Pega yang **tidak dapat dipastikan dari ekspor**. Perilakunya
**tidak dinyatakan** di sini. → **OQ-035**.

### 5.2 Nilai `.KomiteLoop` tidak di-set di modul ini

`.KomiteLoop` menentukan **berapa tingkat** tangga persetujuan, tetapi tidak ada rule di
`Komite Claim Life` yang mengisinya (`grep -rl "KomiteLoop" "Komite Claim Life"` → hanya
`KomitePostAdjustment`, `SendEmailKlaimLife`, `Flow`, `When/IsKomiteLoop`, seluruhnya **membaca**).
Sumber nilainya **tidak diketahui** → **OQ-032**.

### 5.3 Rule yang belum ditelusur

`PrintAkseptasiPDF`, `LoadDocumentLife_ACT`, `InsertJsonClaimLife_Act`, `SendEmailKlaimLife`,
`HitServiceToKasirKMTLife_Act`, `SetInformationData`, `ViewTransferDtl` (FlowAction) — di luar
jalur inti.

---

## 6. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `Komite Claim Life/`) | OQ-011? |
| --- | --- | --- |
| `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITELIFE_FLOW` / `Rule-Obj-Flow` (dari `<pxObjClass>`) | `Flow/KomiteLife_Flow.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEROUTER` / `RULE-OBJ-ACTIVITY` | `Activity/KomiteRouter.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY` | `Activity/KomitePostAdjustment.xml` | tidak |
| `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `ISKOMITELOOP` / `RULE-OBJ-WHEN` | `When/IsKomiteLoop.xml` | tidak |
| (FlowAction) `VIEWTRANSFERDTL` | `FlowAction/ViewTransferDtl.xml` | tidak |
| `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN` | `When/IsPEGAPROD.xml` | **YA — entri #305** |
| 6 rule `RULE-CONNECT-SQL` unik (§3) | `RDBList/` | sebagian (lihat `_SUMMARY-komite.md`) |

**12 rule ditelusur.**

---

## 7. Pertanyaan terbuka baru

| OQ | Ringkas | Pemilik |
| --- | --- | --- |
| **OQ-032** | `.KomiteLoop` menentukan jumlah tingkat tangga, tetapi tidak ada rule di modul ini yang mengisinya | Product+UW |
| **OQ-033** | `KomiteRouter` hanya men-set `param.AssignTo` bila `TransferType == '2'`; apa yang terjadi bila tidak | Product+UW |
| **OQ-034** | Tanggal cutover ter-hardcode `20250207T000000.000 GMT` menggerbangi 2 langkah | Product+UW |
| **OQ-035** | Dua activity dipanggil tetapi salinannya hanya ada di modul lain (`UpdateWorkObject`, `serviceInsertArasapasClaimLife_act`) | pemilik export Pega |

**Guard identitas orang:** `[terverifikasi]` **nol** di modul ini — tidak ada
`pyUserIdentifier`/`pyTelephone`/`pxCreateOperator` sebagai literal. Hanya **1 file** memakai
`pyPosition` (`grep -rlE "pyPosition[ ]*[=!]+[ ]*[\"']" "Komite Claim Life"`). Berbeda tajam dari
Komite Claim FacIn dan Komite Claim Prop → lihat `_SUMMARY-komite.md`.

OQ yang dikuatkan: OQ-011 & OQ-029 (§4), OQ-013 (§3.1), OQ-020 (§2.2, §3), OQ-028 (§1.2),
OQ-030 (§2.2), OQ-031 (§2.2 — `RetroID` `"1000013"`).
