# Understanding Report — Bagaimana Aplikasi Nusantara Re Bekerja

FASE A Discovery, **STEP D4**. Disusun 2026-09-14 dari D1 (20 inventaris), D2 (20 telusur konteks /
428 rule), D3 (20 catatan modul + glossary 169 entri). Korpus 9.369 file, READ-ONLY.

**Ini narasi berbukti, bukan spesifikasi.** Setiap pernyataan menyebut `path + rule`; yang tidak
pasti ditandai dan didaftarkan ke `open-questions.md`. Tidak ada ADR/BRD/spec/tiket di sini, tidak
ada rancangan Go/React/skema.

---

## 1. Gambaran satu halaman

`[terverifikasi]` Aplikasi ini menjalankan **bisnis reasuransi** di tiga lini —
**facultative inward**, **treaty inward**, dan **life** — ditambah **siklus klaim** dan
**tangga persetujuan komite**. Sembilan bounded context terbaca dari korpus; satu konteks yang
seharusnya ada (**Identity & Access**) **tidak ada sama sekali**.

`[terverifikasi]` Empat sifat struktural mewarnai seluruh sistem:

1. **Logika bisnis inti banyak yang berada di database, bukan di rule** — 67 stored procedure
   kustom tanpa body di korpus; pada dua modul **seluruh jalur tulis master melewati procedure**.
2. **Otorisasi tidak dimodelkan** — tidak ada satu pun rule peran; otorisasi tersebar sebagai nilai
   data dalam empat pola berbeda, termasuk **5 identitas orang ter-hardcode sebagai tujuan rute**.
3. **Dua pola pengorganisasian yang berlawanan hidup berdampingan** — facultative & treaty master
   adalah satu ruleset bercabang saat runtime (99,0 % / 85,8 % tumpang tindih); trio Claim adalah
   tiga basis kode terpisah (12–40 %).
4. **Uang dan ambang ter-hardcode tanpa mata uang** — pangsa `0.45`/`0.05`, ambang
   `"3000000000"` (dibandingkan **sebagai string**), limit `30000000.00` / `50000000.00` /
   `57750000.00`.

---

## 2. Narasi per konteks

### 2.1 Facultative Inward — offer → akseptasi → R/I slip → fac out/retro → produksi

`[terverifikasi]` Titik masuk `NB FacIn/Flow/InputQuotation.xml`
(`ASM-FW-GISFW-WORK-NB!INPUTQUOTATION`) memutuskan apakah sebuah quotation berjalan sebagai
**facultative** (`SubProcess9` → `InputInwardFacultativeOffer`) atau **treaty inward**
(`SubProcess10` → `InputRealizationTreatyIn`). Inilah simpul yang menyatukan dua lini.

**Alurnya**, dari `NB FacIn/Flow/InputInwardFacultativeOffer.xml`
(`ASM-FW-GISFW-WORK!INPUTINWARDFACULTATIVEOFFER`, **graf terbesar korpus**: 21 Assignment,
42 Decision, 12 Utility, 170 connector):

```
MARKETING [WorkBasket]
   → "Accept?" (DecisionTable IsUWAccepted)
        ├─ confirm → "IS LIFE?" → MEDICAL LIFE  (hanya siklus new business)
        │             atau → SAVE JSON_OFFER → Declaration/SpecialCase → "Is it group?"
        ├─ banding → SetBanding_ACT
        └─ decline → End
   → Group? / ISPKSASM? / UW FINANCIAL? 
   → "SET LIMIT_AKSEPTASI" (GetLimitAkseptasi_ActFlow)
   → "Limit Akseptasi" merutekan berdasarkan token LetterNo:
        KADIVTEKNIK → KADIV TEKNIK      DIREKTURTEKNIK → DIREKTUR TEKNIK
        KADIVFINANCIAL → KADIV FINANCIAL  DEPHEADUNDERWRITER → DEP HEAD UW FACULTATIVE  …
   → MARKETING (BINDING) → "Accept?" 
        ├─ confirm → SubProcess POLICY → InputInwardFacultativeRISlip
        └─ banding → "Letter No Null" → SetBanding_ACT
   → R/I Slip: Check Spreading (CekLimitSpreading_Act) → tangga akseptasi
             → SAVE JSON_POLICY → HIT SERVICE ARASAPAS → cek status konversi
   → Fac out/retro (OfferFacRetro): ReasFacOutAdmin → ReasFacOutHead → [END]
```

`[terverifikasi]` **Rantai routing terbaca ujung-ke-ujung** — ini salah satu temuan paling berguna:
`DataTransform/SetBandingProposal_DT.xml` men-`SET` `.ProposalAcceptStatus := 4`,
`.OfferFacIn.IsBanding := "true"`, dan **`.LetterNo := .BandingTo`**; lalu shape
"Limit Akseptasi" memilih Assignment berdasarkan nilai `LetterNo`. Jadi **field nomor surat adalah
mekanisme routing persetujuan** → OQ-045.

`[terverifikasi]` **Renewal dan endorsement memakai ruleset yang sama**, dibedakan
`Quotation.StatusBusiness` (1 = `IsNB`, 2 = `IsRenewal`, 3 = `IsEDM`) — keempat rule `When` itu ada
di ketiga modul. Selisih perilakunya terukur: renewal **tanpa cabang Life sama sekali**;
endorsement punya tiga gerbang khas (`IsEDMRiSlip`, `IsEdmPPNPPH`, `IsEdmInternalRetro`) dan
tangga fac out **4 anak tangga** (vs 2).

`[terverifikasi]` **Taksonomi endorsement terbaca penuh**: 14 rule `When` menguji
`StatusBusiness = 3` + `EdmType = 4` + `QuotationData.Type` menghasilkan **12 nilai** —
`0` RefNo/RISlip, `1` ExtendPeriod, `2` AdjTSI, `3` AdjRate, `4` AdjSpreading, `5` AddObject,
`6` AdjPeriod, `7` AdjCurrency, `8` AdjRIC, `9` AdjInsured, `11` AdjShareCedant,
`12` AdjCeding/AdjRIC/PPNPPH. Nilai `10` tidak dipakai. **Arti tiap nilai belum terverifikasi.**

**Yang tidak dapat dinyatakan:** pemetaan `ProposalAcceptStatus` (`1/2/3/4/7/9`) ke enam hasil
(`confirm`/`reject`/`ask`/`banding`/`revise`/`decline`) — baris `DecisionTable/IsUWAccepted.xml`
**tidak ikut terekspor** (OQ-043); nilai `4` bahkan **ambigu** (ditulis rule "banding", diuji rule
"facout" — OQ-044); dan **rumus spreading/capacity/scoring** ada di 34–38 activity per modul yang
belum dibaca (terbesar 996.052 byte).

### 2.2 Treaty Inward — master, akseptasi, realisasi, endorsement

`[terverifikasi]` **Dua konteks yang bertemu di satu master.**

**Master & akseptasi** (`Treaty In`, `Treaty In Adjustment` — tanpa rule `Flow`). Prosesnya nyata
tetapi ditulis sebagai **satu Data Transform bersarang**, bukan graf:
`DataTransform/Akseptasi_DT.xml` (`DATA-PORTAL!AKSEPTASI_DT`) — 10 `WHEN`, 18 `OTHERWISE_WHEN`,
64 `SET`.

```
tombol Submit  →  TreatyInSubmit (7 langkah)
                     TreatyInCheckID → TreatyInCheckError
                     precondition: OutputParam.ERRMSG == ""
                     Apply-DataTransform Akseptasi_DT   ← mesin status
                     AddCommentList_Act → SaveTreatyIn_Act → TreatyInInputVis

Akseptasi_DT:
  Position ""/Admin    + pyTelephone TREATY1|TREATY2|SPVTREATY1|SPVTREATY2 → SecHead, Accept
  Position SecHead     + Accept  → DeptHead, Accept
  Position DeptHead    + Accept  → Director, Accept
  Position GroupLeader + Accept  → Director, Accept
  Position Director    + Accept  → "",       Resolve Complete
  mana pun             + Reject  → Admin,    Reject   (PositionUsername := CommentList(1).OperatorName)
  mana pun             + Decline → "",       Decline
  RevisionState == 1   + Accept di SecHead → "", Resolve Complete   ← memotong 2 tingkat
```

`[terverifikasi]` **Pengguna memilih 3 nilai** (`Accept`/`Reject`/`Decline`) tetapi **status punya
4** — `Resolve Complete` hanya dihasilkan mesin. **Beda `Reject` vs `Decline` tidak dijelaskan
korpus** → OQ-052.

`[terverifikasi]` **Revisi dikodekan di dalam ID**: `TreatyInRevisi_post` menambahkan sufiks
`/R01`, `/R02`, … pada posisi karakter 10–12, dan menyalin nilai lama ke `TreatyIn.OLDID`; sisi
database mengambil nomor berikutnya dengan `TO_NUMBER(SUBSTR(ID,10,2))+1`. **Ter-hardcode di dua
tempat** → OQ-055.

`[terverifikasi]` **Realisasi & endorsement** (`NB Treaty In`, `EDM Treaty In`) berbasis `Flow`.
Alur realisasi: Input Realitation [WorkBasket] → "Is Correct?" → percabangan SPV/Treaty →
Acceptance by Head. Treaty → (opsional) Dept. Head → cek nomor polis kosong → SAVE JSON policy →
**HIT SERVICE ARASAPAS**. Format nomor polis terbaca penuh:
`'RNM-' || CARI20 || '.T' || BusinessOldId || '.' || MM.yyyy || '.' || LPAD(seq,5,'0')`.

`[terverifikasi]` **Langkah terakhir terhenti** — `SERVICEINSERTARASAPAS_ACT` yang dipanggil hanya
ada di dua modul lain dan berkonflik → **OQ-025**.

### 2.3 Treaty Arrangement — master term yang namanya menyesatkan

`[terverifikasi]` `Treaty Contract Out` **bukan modul treaty outward**. Lima uji independen
(§2.4 `context-map.md`) menunjukkan: nol objek `*_OUT*`, nol class `*OUT*`, 81 activity
`*TreatyArr*` vs 4 `*Out*` yang semuanya soal lampiran, 9 Connect-SQL **menulis** master treaty,
dan penyentuh objek outward sesungguhnya adalah `NB Treaty In` (8 file), `Claim Non Prop` (4),
`EDM Treaty In` (4) — **modul ini 0**.

`[terverifikasi]` Yang dilakukannya: mengelola **16 jenis klausul** kontrak treaty lewat pola
`New…` / `Set…` / `Save…` / `CancelActivity…`, dengan seluruh penyimpanan lewat stored procedure.

**→ OQ-022 terjawab sebagian.** Yang tersisa: mengapa dinamai "Out", dan ke konteks mana ia
ditempatkan.

### 2.4 Klaim — register → akseptasi → komite

`[terverifikasi]` **Empat siklus kerugian yang berbeda-beda**, dan tidak ada dua modul dengan
rangkaian tahap sama:

| Modul | Tahapan |
| --- | --- |
| `Claim Life` | Register → Outstanding → **Medical Check** → Claim Analis |
| `Claim Fac In` | (percabangan **B2B** via `IsSPK`) → Register → Estimasi → Choose Surveyor |
| `Claim Prop` | Outstanding → Input Acceptation |
| `Claim Non Prop` | Outstanding → Input Acceptation |

`[terverifikasi]` **Adjustment, Close, dan Reject tidak muncul sebagai shape di modul mana pun** —
seluruhnya Activity di luar graf → OQ-039.

`[terverifikasi]` **Jembatan ke Komite adalah Activity, bukan shape**:
`Claim Non Prop/Activity/CreateChildKomiteCNP_Act.xml` (756.836 B) dan
`AddKomiteTreatyChild_ACT` (420.247 B), keduanya di class `ASM-FW-GCNMFW-DATA-ADJUSTMENT`,
**digerbangi perbandingan limit**.

`[terverifikasi]` **Limit wewenang dikelola dua cara berbeda dalam satu konteks**: `Claim Non Prop`
meng-hardcode (`LimitMax = 30000000.00`, `LimitMaxDivHead = 50000000.00`,
`LimitPersenMax = 30.00`) sementara `Claim Prop` mengambilnya dari database
(`GetLimitDirekturUtama_SQL`, `GetLimitPLATreatyin`, `GetLimitsTreatyIn_SQL`).
Batas bawah **30.000.000,00 sama persis** dengan ambang roster komite di `Komite Claim FacIn`,
tetapi batas atasnya berbeda (50 jt vs 57,75 jt). **Mata uang tidak disebut** → OQ-040, OQ-037.

### 2.5 Komite — tangga persetujuan berbasis data

`[terverifikasi]` Keempat modul Komite punya graf yang **persis sama bentuknya** — 1 Assignment
"KomiteRouter" [`WorkList`] + 1 Decision "KomiteLoop", 4 connector — dengan kondisi loop **sama
teksnya**: `.AcceptStatus = "1"` **DAN** `.KomiteCount <= .KomiteLoop`.

**Tangga persetujuan tidak dimodelkan sebagai shape.** Berapa tingkat dan siapa penyetujunya
ditentukan **data**: tiap tingkat menulis entri `KomiteList(KomiteCount)` berisi `KomiteAproval`,
`KomiteComment`, `DateApprove`.

`[terverifikasi]` **Roster berasal dari database** (class `ASM-FW-GCNMFW-Int-EMAILKOMITE`),
tetapi **sasaran routing justru ter-hardcode** (`komitepnc`…`komitepnc4` di tiga modul) → OQ-036.
**Komposisi roster bergantung besaran nilai**: `Local.TotalAdj > 30000000.00 && <= 57750000.00`
menentukan `Obj-Browse` mana yang dijalankan → OQ-037.

`[terverifikasi]` **Empat implementasi berbeda untuk satu konsep**: properti batas tangga
(`.KomiteLoop` vs `Local.TotalKomite`), sasaran routing, guard identitas (FacIn 4 file, Prop 3,
Life 0, Non Prop 0), dan rule penomoran akseptasi (2 per kode `Type` / 2 per `IsFire` / 1 / 1).

### 2.6 Life — penawaran, premium list, endorsement, master

`[terverifikasi]` `PremiumList Life` menjalankan penawaran hingga daftar premi:

```
Input Offer [WorkList] → "Accept" (IsLifeAccepted)
   ├─ Decline → End
   └─ Confirm → "FlagOnGoingPolicy" (IsFlagOnGoingPolicy)
        ├─ Offer   → END52   (selesai di tahap penawaran)
        └─ Premium → Input Premium List Detail [WorkList] → "Accept"
                        ├─ Reject  → balik ke Input Offer
                        └─ Confirm → InsertJsonPolisLife_Act → END52
```

`[terverifikasi]` `Endorsement Life` memakai **mesin yang sama** (`InsertJsonPolisLife_Act`,
`IsLifeAccepted`) dengan graf paling sederhana di korpus (2 Assignment, 1 Decision, 1 Utility).

`[terverifikasi]` Kedua modul life memakai **`WorkList`**, berbeda dari treaty inward yang
seluruhnya `WorkBasket` → OQ-028.

`[terverifikasi]` **Master life** (`Master Product Name Life`, `Master Contract Retro Life`) adalah
editor master murni tanpa status proses. `Master Contract Retro Life` adalah modul paling
terisolasi di korpus (tanpa `ConnectREST`/`When`/JSON/db-link/guard identitas).

### 2.7 Antar-konteks: bagaimana data mengalir

`[terverifikasi]`

```
  InputQuotation ──┬─> Facultative Inward (offer → akseptasi → RI slip → fac out) ──> FACINPRODUCTION
                   └─> Treaty Inward Realisasi ──> JSON policy ──> [ARASAPAS: terblokir]

  Treaty In / Adjustment ──(Akseptasi_DT)──> M_TREATY_IN  <──(ASM!SAVETREATYIN, rule yang SAMA)──┐
                                                                                                 │
  Master Product Name Life ────────────────────────────────────────────────────────────────────┘
                                            (penulis kedua — OQ-056)

  Treaty Contract Out ──(9 Connect-SQL)──> TREATYCONTRACT, M_TREATYYEAR, TREATYREINSURER,
                                           TREATYBUSINESS, M_PROPORTIONALARRG
                                                   │ dibaca oleh
                                                   v
  Claim Fac In / Claim Prop / Komite Claim Prop

  Claim (4 modul) ──(CreateChildKomiteCNP_Act / AddKomiteTreatyChild_ACT, digerbangi limit)──>
  Komite (4 modul) ──> OS_AKSEPTASI_KLAIM(_LIFE)  <── juga ditulis/dibaca sisi Claim
                                                     (UpdateOsAkseptasiClaimLife_sql = rule yang SAMA)

  Claim Non Prop ──> M_TREATY_OUT / M_TREATY_OUT_DETAIL / TREATY_OUT   (OQ-042)
```

---

## 3. Peta integrasi eksternal total

`[terverifikasi]` **51 file `RULE-CONNECT-REST`, 15 nama layanan berbeda.**
Audit: `find . -path ./OUTPUT_HASIL_RNM -prune -o -type f -path "*/ConnectREST/*.xml" -print | wc -l`

| Layanan | File | Domain |
| --- | ---: | --- |
| `ServiceGoogle` | **15** | seluruh domain |
| `convertJsonNusareToProduction` | 6 | facultative, treaty, life |
| `SendAcceptationToKasir` | 6 | claim, komite |
| `KonversiKlaimNonLife` | 6 | claim, komite |
| `getPremiumPaidOn` | 4 | facultative |
| `getPayAttachment` | 3 | claim |
| `insertClaimFinalOrClosed_NP`, `getPremiumPaidOnTreatyIn` | 2 masing-masing | claim |
| `insertClaimReject_NP`, `getPremiumPaidOnMarine`, `getPaymentClaim`, `convertJsonNusareToProductionClaimLife`, `InsertClaimOutstanding_NP`, `HitDLAClaimFacin`, `GetDtlPaymentClaim` | 1 masing-masing | claim, facultative |

### 3.1 Alamat ditentukan dalam dua lapis — keduanya di luar rule

`[terverifikasi]`

1. **`RULE-ADMIN-SYSTEM-SETTINGS` `LinkService!LinkService`** — `pyBaseURLSelectionType = SETTING`
   pada hampir seluruh `ConnectREST`. Rule ini **identik** di modul-modul yang membawanya setelah
   normalisasi 21 tag (konflik semu OQ-011 #324) → konfigurasi base URL **tidak bercabang**.
2. **Tabel Oracle `M_LINK_SERVICE`** — `Activity/GetLinkService.xml`
   (`ASM-FW-GISFW-INT-M_LINK_SERVICE!GETLINKSERVICE`) melakukan `Obj-Browse` dengan kunci
   `KATEGORI_1` / `KATEGORI_2`, hasilnya dipakai langkah `Connect-REST` berikutnya.
   **Isi tabel tidak ada di korpus** → **OQ-047**.

`[terverifikasi]` **Satu-satunya URL literal endpoint** di seluruh korpus ada di
`Claim Fac In/ConnectREST/getPremiumPaidOnMarine.xml` (`pyBaseURLSelectionType = URL`); nilainya
**tidak disalin** ke artefak mana pun karena memuat contoh data. Korpus juga memuat hostname DEV
`appdev.nusantarare.com` (7×) di 81 file / 15 modul → OQ-018.

### 3.2 Integrasi per nama

| Integrasi | Cara tersentuh | Status telusur |
| --- | --- | --- |
| **Arasapas** | shape "HIT SERVICE ARASAPAS" → `serviceInsertArasapas_act` / `…RNW_act` / `…EDM_act`; skema Oracle `ARASAPAS` | **20 langkah terbaca** (varian `Endorsment Fac In`); **terblokir** di NB Treaty In / NB FacIn / RNW Fac In — OQ-025 |
| **Kasir** | `SendAcceptationToKasir` (6 file), `POOLDATA.DIRECTTOKASIR_LOG`, `HitServiceToKasir_Act.xml` (482 KB), `HitServiceToKasirKMT_Act` | isi activity **belum dibaca** |
| **Konversi** | `KonversiKlaimNonLife` (6), `convertJsonNusareToProduction` (6); `.StatusService.StsKonversiFacIn/FacOut`; `UpdateStsKonversiFacOut_Act` | gerbang `IsSuccessHitService` terbaca; isi activity **belum dibaca** |
| **Google Storage** | `ServiceGoogle` (15), `InsertGoogleStorage_Act`, `GetUrlGoogleStorage_Act`, `DeleteGoogleStorage_Act`, `T_STORAGE_IMAGE`, `POOLDATA.GET_TOKEN_STORAGE` | rantai terbaca; isi procedure **tidak ada di korpus** |
| **Gemini AI** | `Activity/GeminiAIGoogle_Act.xml` — **hanya di `NB FacIn` dan `RNW Fac In`** | isi **belum dibaca** → **OQ-048** |
| **`M_LINK_SERVICE`** | `Activity/GetLinkService.xml`, dipanggil dari Arasapas & Google Storage | isi tabel **tidak ada di korpus** → **OQ-047** |
| **officeapps** | **tidak ditemukan** sebagai `ConnectREST` maupun nama rule | **tidak ada bukti** — tidak dicatat sebagai integrasi |

---

## 4. Batas pengetahuan total

Dibedakan tegas: **batas korpus** (informasinya memang tidak ada — perlu jawaban manusia atau akses
database) vs **batas cakupan telusur** (ada di korpus, belum dibaca karena volume).

### 4.1 Batas korpus

| # | Batas | Ukuran terverifikasi | OQ |
| ---: | --- | --- | --- |
| 1 | **Stored procedure tanpa body** | **70** dipanggil (**67 kustom** + 3 bawaan Oracle `DBMS_LOB`, `DBMS_OUTPUT`, `UTL_MATCH`), tersebar di **10 skema aplikasi**: `POOLDATA` (61), `FIRE` (2), `GENERAL`, `GL`, `MBU`, `NEW_GENERAL`, `NEW_UNDERWRITING`, `ARASAPAS`, `DATAPEGA`, `REINSURANCE` | **OQ-002**, OQ-013, OQ-016 |
| 2 | **Objek lewat database link** | **8** via `@ASMD.SINARMAS.CO.ID`: `M_EQS_RATE`, `M_TERORISME_RATE`, `M_RSMD_RATE`, `M_FLEXAS_RATE`, `M_FLOOD_AREA`, `FIRE.M_FLOOD_RATE`, `FIRE.M_BI_INDEMNITY`, `LST_KURS_STANDARD` — **hanya di 3 modul facultative** | **OQ-017** |
| 3 | **Objek berlabel JSON** | **10**: `JSON_POLIS`, `JSON_OFFER`, `JSON_KLAIM`, `JSON_FOLLOWING`, `JSON_POLIS_ERROR`, `POOLDATA.JSON_POLIS`, `POOLDATA.JSON_OFFER`, `POOLDATA.JSON_OFFER_LIFE`, `POOLDATA.JSON_FOLLOWING`, `POOLDATA.JSON_POLIS_MONITORING`; + kolom `JSONDATA`/`DATA_JSON`. Pembentuknya **`@ASM.GetPageJSONString()`** — fungsi kustom **tanpa source di korpus** | **OQ-012** |
| 4 | **Tabel keputusan tanpa baris** | **49 file `DecisionTable`**, **tak satu pun** memuat baris kondisi→hasilnya. Termasuk `IsUWAccepted` yang menggerbangi **22 shape Decision** di flow terbesar korpus | **OQ-043** |
| 5 | **Keluarga `When` tak terbaca** | **5**: `IsPEGAPROD`; `IsSendtoMedical`/`IsSPK`; `ISCLM*` (12 varian); `IsPKSASM`/`LetterNoNull`/`ToJUW_A`/`ToUW`/`IsEdmInternalRetro`; `IsTreatyUser` (6 kondisi)/`recordEvent` | **OQ-029**, OQ-041 |
| 6 | **Identitas rule berkonflik** | **533**; 438 (82,2 %) menyentuh ketiga modul facultative, **429 tetap berbeda** setelah normalisasi 21 tag | **OQ-011** |
| 7 | **Tidak ada DDL / definisi properti** | nol file; tipe kolom tidak dapat disimpulkan dari SQL | **OQ-001** |
| 8 | **Arti kode & enumerasi** | `ProposalAcceptStatus`, `StatusAkseptasi`, `PaymentType`, `EdmType`, `QuotationData.Type`, `StatusBusiness`, `ProRateType`, `AcceptStatus`, `TransferType`, `L1`–`L11`; **`REINSTYPEID` dan kode `OR` tanpa satu pun nilai literal** | **OQ-020**, OQ-038, **OQ-057** |
| 9 | **Slot parameter generik** | `CARI1`…`CARI30` dipakai sebagai slot menuju SQL di ≥5 modul | **OQ-059** |
| 10 | **Kelengkapan ekspor** | folder `excludeXML`; folder **`Claude outputs`** (`.xlsx`, `.diff`) di 4 modul; hostname DEV; 2 `pxHostId` berbeda | OQ-003, **OQ-054**, OQ-018 |
| 11 | **Tidak ada rule identitas/otorisasi** | nol rule peran di 17 tipe rule | **OQ-007** |

Audit untuk baris 1:
```
awk -F'\t' '$2=="RULE-CONNECT-SQL"{print $8}' all-rules.tsv | grep -oE "procs=[^;]*" \
 | sed 's/procs=//' | tr ',' '\n' | sed 's/^ *//;s/ *$//' | grep -v "^$" | sort -u | wc -l   # 70
```

### 4.2 Batas cakupan telusur

| # | Yang belum dibaca | Ukuran |
| ---: | --- | --- |
| 1 | Rumus **spreading / capacity / scoring** facultative | 34–38 activity per modul; terbesar `SumTSIPremiSpreadedRNM_FIRE_Act.xml` **996.052 B** |
| 2 | Rumus **limit / layer / ROL** treaty inward | 31 activity; terbesar `TreatyInNPSetTotal.xml` **825.279 B** |
| 3 | Activity besar domain **Claim** | `CLaimFaceSheet_Act.xml` 804 KB, `CreateChildKomiteCNP_Act.xml` 757 KB, dan 20+ lainnya |
| 4 | Mesin **data lama** endorsement | `SetValueToEDMWork` 446 KB + 7 `SetOLDValueToEDMWork_*` |
| 5 | Rumus **share retro life** | `CountingPercentShare_Act`, `TreatyLimit_TypeProtect`, `SetValueRetroLimit_TreatyYearLife` |
| 6 | Selisih isi rule yang tetap berbeda | 43 (`Treaty In` ↔ `Adjustment`), 12 (`Master Product Name Life` ↔ `Treaty In`) |
| 7 | **Harness terbesar korpus** | `InputTreatyInOffer.xml` **8.225.095 B** — hanya di-grep |
| 8 | `GeminiAIGoogle_Act`, `HitServiceToKasir_Act` | isi belum dibaca |

`[terverifikasi]` **§4.1 memerlukan jawaban manusia atau akses database. §4.2 hanya memerlukan
waktu telusur tambahan** dan tidak memblokir siapa pun.

---

## 5. Risiko migrasi teridentifikasi

Seluruhnya berbukti; **ini identifikasi risiko, bukan rekomendasi mitigasi** (itu FASE B).

### 5.1 Korpus belum tentu mencerminkan production

`[terverifikasi]` Korpus memuat hostname DEV **`appdev.nusantarare.com` (7×)** di **81 file /
15 modul**, dan dirakit dari **lebih dari satu server Pega** (dua `<pxHostId>` berbeda).
Ditambah **10 tombol ber-label `(dev)`** di panel aksi treaty inward — termasuk
`Force Resolve Complete(dev)` yang **melompati seluruh tangga persetujuan** dan `Force Edit (dev)`.
Ditambah **`DBMS_OUTPUT.PUT_LINE`** di SQL, dan folder **`Claude outputs`** (`.xlsx`, `.diff`)
di dalam ekspor 4 modul.
→ OQ-018, **OQ-050**, OQ-054.

**Risikonya:** memigrasikan perilaku yang tidak berlaku di production, atau melewatkan perilaku
production yang tidak ada di korpus. **533 identitas rule berkonflik** memperbesar ini — versi mana
yang berjalan **tidak diketahui** (OQ-011).

### 5.2 Rumus bisnis berada di luar korpus

`[terverifikasi]` **67 stored procedure kustom tanpa body.** Pada `Treaty Contract Out` dan
`Master Contract Retro Life`, **seluruh jalur tulis master melewati procedure** — aturan validasi,
versi, kunci, dan kaskade **tidak ada di korpus sama sekali**. Ditambah **8 tabel rate & kurs di
database lain** lewat db-link (OQ-017), **perhitungan prorata tanggal** di
`FIRE.CEK_PRORATA_TANGGAL`, dan **struktur JSON** yang dibentuk fungsi kustom tanpa source
(OQ-012).

**Risikonya:** perhitungan premi, spreading, dan limit **tidak dapat direplikasi** tanpa akses ke
database dan jawaban DBA. Ini pemblokir FASE B yang paling luas.

### 5.3 RBAC harus dirancang, bukan dimigrasikan

`[terverifikasi]` **Tidak ada satu pun rule identitas/otorisasi di korpus** (OQ-007). Yang ada
adalah empat pola berbeda (`context-map.md` §2.10), termasuk:

- **5 identitas orang ter-hardcode sebagai pemilik tugas berikutnya** di `Akseptasi_DT`
  (OQ-053) — bila orang itu berganti jabatan, rule harus diubah dan di-deploy ulang;
- **otorisasi lewat indeks tetap** `OperatorID.pyWorkBasketList(2)` — bergantung urutan daftar
  workbasket di profil operator (OQ-051);
- **field bisnis menyimpan kode peran**: `pyTelephone` (4 nilai), `LetterNo` (10 token)
  (OQ-027, OQ-045);
- **4 sasaran routing komite ter-hardcode** `komitepnc`…`komitepnc4` (OQ-036).

**Satu pola sudah berbasis data** dan dapat menjadi titik awal: pencocokan roster di
`Komite Claim Non Prop`.

**Risikonya:** tidak ada yang bisa disalin. Model peran harus dibuat dari nol, dan **setiap
keputusan siapa-boleh-apa memerlukan jawaban bisnis**.

### 5.4 Uang material tidak boleh `float`

`[terverifikasi]` Nilai uang dan ambang yang ditemukan, **seluruhnya tanpa penyebutan mata uang**:

| Nilai | Lokasi |
| --- | --- |
| pangsa `0.45` / `0.05` | `NB FacIn/Activity/CountPremiNusareRetro_Act.xml` (identik di 3 modul facultative) |
| ambang **`"3000000000"`** — dibandingkan **sebagai string berkutip** | idem (`.TSILiability <= "3000000000"`) |
| `LimitMax = 30000000.00`, `LimitMaxDivHead = 50000000.00`, `LimitPersenMax = 30.00` | `Claim Non Prop/Activity/CreateChildKomiteCNP_Act.xml` |
| pita `> 30000000.00 && <= 57750000.00` | `Komite Claim FacIn/Activity/ApprovalKomite_Act.xml` |

`[terverifikasi]` Perhitungan retro memakai `@toDecimal(...)`, tetapi **ambangnya dibandingkan
sebagai string** — perilakunya pada nilai di luar rentang perlu dikonfirmasi.
→ OQ-037, OQ-040, **OQ-046**.

**Risikonya:** representasi uang harus ditetapkan eksplisit (aturan proyek: **jangan `float`**),
dan mata uang harus dijawab sebelum angka mana pun dipindahkan.

### 5.5 Identitas dan host ter-hardcode

`[terverifikasi]` Guard identitas orang tersebar: `NB FacIn` **43 file**, `RNW Fac In` 36,
`Endorsment Fac In` 33, `NB Treaty In` 12, domain Claim 2–7 per modul.
Hostname/URL literal: **81 file / 15 modul**. Endpoint sesungguhnya di **tabel database**
(OQ-047), bukan di rule.

**Risikonya:** endpoint & host harus menjadi konfigurasi/env var (aturan proyek), dan **inventaris
endpoint tidak lengkap** sampai isi `M_LINK_SERVICE` diketahui.

### 5.6 Model AI pihak ketiga di dalam alur underwriting

`[terverifikasi]` `Activity/GeminiAIGoogle_Act.xml` ada di **`NB FacIn`** dan **`RNW Fac In`**
(tidak di `Endorsment Fac In`), dan termasuk 9 rule yang memakai `M_LINK_SERVICE`.
**Isinya belum dibaca.** → **OQ-048**.

**Risikonya:** apa yang dikirim (apakah termasuk data tertanggung / data pribadi) dan apakah
hasilnya memengaruhi keputusan underwriting **tidak diketahui**. Perlu penilaian kepatuhan dan
privasi sebelum migrasi — pemilik: Product+UW + Security.

### 5.7 Ketergantungan yang terputus di korpus

`[terverifikasi]` `ASM-FW-GISFW-WORK!SERVICEINSERTARASAPAS_ACT` **dipanggil** `NB Treaty In`,
`NB FacIn`, `RNW Fac In` tetapi **tidak ada** di ketiganya (OQ-025). Ditambah **OQ-035**
(activity dipanggil tetapi salinannya hanya di modul lain).

**Risikonya:** ada jalur eksekusi yang **tidak dapat dinyatakan** dari korpus. Risiko salah
tafsirnya lebih kecil dari perkiraan D1 (kedua varian berbeda hanya 3 tag metadata), tetapi
**keputusan versi mana yang aktif tetap milik pemilik export**.

---

## 6. Apa yang sudah pasti, dan apa yang belum

`[terverifikasi]` **Sudah pasti dan dapat dipakai sebagai dasar FASE B:**

- struktur proses 20 konteks, dengan titik masuk dan guard yang terbaca;
- pemetaan modul → bounded context dengan kualitas bukti per konteks (`context-map.md`);
- inventaris integrasi eksternal (51 ConnectREST, 15 layanan) sebagai **konfigurasi**;
- daftar objek Oracle per modul dan arah baca/tulis lintas konteks;
- nilai literal seluruh kode/status (**apa adanya**, tanpa tafsir);
- 169 entri glossary dengan kolom Bukti (`glossary.md`).

`[pertanyaan terbuka]` **Belum pasti dan memblokir FASE B:**

- arti setiap kode dan status (OQ-020 dan keluarganya);
- seluruh rumus yang berada di stored procedure dan db-link (OQ-002, OQ-017);
- struktur JSON (OQ-012) dan tipe kolom (OQ-001);
- aturan persetujuan sesungguhnya (OQ-043 — tabel keputusan tak terekspor);
- versi rule mana yang berlaku di production (OQ-011);
- model peran (OQ-007 dan keluarganya);
- penetapan bounded context final (OQ-010, OQ-015, OQ-019, OQ-022).

Rincian per pemilik peran dan penanda pemblokir: **`D3-D4-CLOSING-REPORT.md`**.
