# Telusur Flow — NB FacIn

STEP D2, Tahap 4 konteks #13. Ditelusur 2026-09-13 dari korpus READ-ONLY `D:\XML\RNM_BRD\`.
Konvensi: `_METHOD.md`. Modul terbesar korpus: **2.083 file**
(`find "NB FacIn" -name '*.xml' | wc -l`).

**Enam rule `Flow`**, class `ASM-FW-GISFW-WORK` kecuali satu:

| # | File | `pxInsName` | `pyStartActivity` | Hash ternormalisasi | Ditelusur di sini? |
| ---: | --- | --- | --- | --- | --- |
| 1 | `Flow/InputQuotation.xml` | `ASM-FW-GISFW-WORK-NB!INPUTQUOTATION` | `Start1` | `d30995fece` | ya — §1.1 |
| 2 | `Flow/InputInwardFacultativeOffer.xml` | `ASM-FW-GISFW-WORK!INPUTINWARDFACULTATIVEOFFER` | `Start1` | `39cfbbec4f` | ya — §1.2 |
| 3 | `Flow/InputInwardFacultativeRISlip.xml` | `ASM-FW-GISFW-WORK!INPUTINWARDFACULTATIVERISLIP` | `Start1` | `c099bf4ebb` | ya — §1.3 |
| 4 | `Flow/OfferFacOut.xml` | `ASM-FW-GISFW-WORK!OFFERFACOUT` | `Start62` | `832fb12b9d` | ya — §1.4 |
| 5 | `Flow/OfferFacRetro.xml` | `ASM-FW-GISFW-WORK!OFFERFACRETRO` | `Start1` | `bfd6252070` | ya — §1.5 (**varian NB/RNW**) |
| 6 | `Flow/InputRealizationTreatyIn.xml` | `ASM-FW-GISFW-WORK!INPUTREALIZATIONTREATYIN` | `Start1` | `1306f68d56` | **tidak — identik dengan `NB Treaty In.md`** |

`[terverifikasi]` **Flow #6 tidak ditelusur ulang.** Hash ternormalisasi `1306f68d56` **sama persis**
dengan `NB Treaty In/Flow/InputRealizationTreatyIn.xml`. Perilakunya sudah direkam di
`NB Treaty In.md` dan berlaku apa adanya di sini. Perintah audit:

```
sh nhash.sh "NB FacIn/Flow/InputRealizationTreatyIn.xml" "NB Treaty In/Flow/InputRealizationTreatyIn.xml"
```

`[terverifikasi]` Flow #3 (`c099bf4ebb`) dan #4 (`832fb12b9d`) juga **identik** dengan salinannya di
`RNW Fac In` — ditelusur **sekali di sini**, dirujuk dari `RNW Fac In.md`.

`[terverifikasi]` Flow #5 **berbeda isi** dari salinan `Endorsment Fac In` (`b370146c63`) → OQ-011
entri #338; varian Endorsment ditelusur terpisah di `Endorsment Fac In.md` §2.

**Catatan OQ-018:** setiap pernyataan menyebut file yang dibaca. Korpus memuat hostname DEV dan
dirakit dari >1 server Pega — **belum tentu cerminan production**.

---

## 1. Diagram alur

### 1.1 `InputQuotation` — titik percabangan facultative vs treaty

Class **`ASM-FW-GISFW-WORK-NB`** — satu-satunya rule `Flow` facultative yang **tidak** berclass
`ASM-FW-GISFW-WORK` `[terverifikasi]`.

```
Start1 ─(Always)─> Decision16 "Transfer Marketing ?"
                       └─ Else ─> Utility2 "SetBusinessType" (Activity SetBusinessType_Act)

Decision27 "Is FacIn or TreatyIn ?"
   ├─ When IsOfferFacIn ─> SubProcess9  "Offer FacultativeIn"  -> Flow InputInwardFacultativeOffer
   └─ When IsTreatyIn ──> Decision28 "Is Treaty In Offer ?"
                              └─ Else ─> SubProcess10 "TreatyIn Realization" -> Flow InputRealizationTreatyIn
```

`[terverifikasi]` **Inilah simpul yang menyatukan dua lini**: satu flow memutuskan apakah sebuah
quotation berjalan sebagai facultative inward atau treaty inward. Bukti: `<pyImplementation>` pada
`SubProcess9` = `InputInwardFacultativeOffer`, pada `SubProcess10` = `InputRealizationTreatyIn`.

`[pertanyaan terbuka]` `Decision27`, `Decision5` ("Is Accepted ?") dan `Utility5` ("Go To Publish" →
`pzChangeStageWrapper`) **tidak dituju connector mana pun**. Jalur dari `Start1` berakhir di
`Utility2`; **percabangan FacIn/TreatyIn tidak terjangkau dari titik masuk yang terbaca**. Pola sama
dengan OQ-023 (NB Treaty In §1.1).

Perintah audit: dari pasangan `<pyFrom>`/`<pyTo>`, tidak ada baris ber-`pyTo = Decision27`.

### 1.2 `InputInwardFacultativeOffer` — siklus penawaran (flow terbesar korpus)

**21 Assignment, 42 Decision, 12 Utility, 3 SubProcess, 170 connector** `[terverifikasi]`
(`grep -o "<pyShapeType>[^<]*" | sort | uniq -c`). Ini **graf terbesar yang ditemui sepanjang D2** —
Claim Fac In sebagai pembanding punya 7 shape.

Jalur utama dari titik masuk:

```
Start1 ─(Always)─> Assignment12 "MARKETING" [WorkBasket, Custom]
                      │ FlowAction: InwardFacultative
                      v
                   Decision3 "Accept?"  (DecisionTable IsUWAccepted)
                      ├─ confirm ─> Decision36 "IS LIFE?"
                      │                ├─ When IsLife ─> Assignment18 "MEDICAL LIFE"
                      │                └─ Else ───────> Utility9 "SAVE JSON_OFFER"
                      │                                    │ Activity SaveJsonOfferFacIn_Act
                      │                                    v
                      │                                 Decision35 "Declaration Policy IsSpecialCase"
                      │                                    ├─ When IsSpecialCase ─> Utility4 SetToJsonOffer_ACT
                      │                                    └─ Else ──────────────> Decision1 "Is it group?"
                      ├─ banding ─> Utility3 "SetBanding_ACT"
                      └─ decline ─> End5
```

Percabangan kelompok dan jalur fac out:

```
Decision1 "Is it group?"
   ├─ When IsGroup ─> Decision30 "pxCreateOperator <orang>?"      <- guard identitas, §7.1
   │                     ├─ When IsGroupCreate ─> Assignment15 "MARKETING"
   │                     └─ Else ──────────────> Decision20 "FAC OUT?"
   └─ Else ─────────> Decision27 "ISPKSASM"
                         ├─ When IsPKSASM ─> Assignment3 "UNDERWRITING"
                         └─ Else ─────────> Decision25 "UW FINANCIAL?"
                                               ├─ When IsTBonding ─> Assignment10 "UNDERWRITING FINANCIAL"
                                               └─ Else ───────────> Decision29 "Declaration Policy"

Decision20 "FAC OUT?"
   ├─ When IsFacRetro ─> SubProcess2 "[Fac Out]" -> Flow OfferFacRetro
   └─ Else ───────────> Decision37 "IS LIFE?"
                           ├─ When IsLife ─> Utility5 "SET LIMIT AKSEPTASI LIFE" (GetLimitAkseptasiLife_Act)
                           └─ Else ───────> Assignment4 "DIREKTUR TEKNIK"

Decision24 "FAC OUT?"
   ├─ When IsFacRetro ─> Decision22 "INPUT FAC OUT?"
   │                        ├─ When IsInputFacRetro ─> Utility2 "SET LIMIT_AKSEPTASI"
   │                        └─ Else ─────────────────> SubProcess3 "[Fac Out]" -> Flow OfferFacRetro
   └─ Else ───────────> Utility2 "SET LIMIT_AKSEPTASI" (GetLimitAkseptasi_ActFlow)
```

Tangga persetujuan berbasis token `LetterNo` (§2.2):

```
Utility2 "SET LIMIT_AKSEPTASI" ─(Always)─> Decision23 "Limit Akseptasi"
   ├─ When ToKadivTeknik ──────> Assignment1  "KADIV TEKNIK"
   ├─ When ToManagerTeknik ────> Assignment5  "MANAGER TEKNIK"
   ├─ When ToDirTeknik ────────> Assignment4  "DIREKTUR TEKNIK"
   ├─ When ToDirMarketing ─────> Assignment6  "DIREKTUR MARKETING"
   ├─ When ToKadivFacultative ─> Assignment8  "KADIV FACULTATIVE"
   ├─ When ToKadivFin ─────────> Assignment11 "KADIV FINANCIAL"
   ├─ When ToDepHeadUW ────────> Assignment17 "DEP HEAD UW FACULTATIVE"
   └─ Else ────────────────────> Utility16 "SAVE JSON_OFFER"

Decision41 "Which team?" / Decision42 "Limit Akseptasi"
   ├─ When ToUW ────────> Assignment3  "UNDERWRITING"
   ├─ When ToSeniorUW ──> Assignment7  "SENIOR UNDERWRITING"
   ├─ When ToJUW_A ─────> Assignment21 "JUNIOR UNDERWITER A"
   ├─ When ToDepHeadUW ─> Assignment17 "DEP HEAD UW FACULTATIVE"
   └─ When IsTBonding ──> Assignment10 "UNDERWRITING FINANCIAL"

Decision38 "LIMIT AKSEPTASI LIFE"
   ├─ When ToDeptHeadUWLife ─> Assignment19 "DEPTHEAD UNDERWRITING LIFE"
   ├─ When ToDirTeknik ──────> Assignment20 "DIREKTUR TEKNIK"
   └─ Else ──────────────────> Utility16 "SAVE JSON_OFFER"
```

Jalur banding dan penutupan:

```
Decision13 "Accept?" (dari Assignment9 "MARKETING (BINDING)")
   ├─ confirm ─> SubProcess1 "POLICY" -> Flow InputInwardFacultativeRISlip
   ├─ banding ─> Decision33 "Letter No Null"
   │                ├─ When LetterNoNull ─> Assignment9 "MARKETING (BINDING)"
   │                └─ Else ─────────────> Utility10 "Banding" (SetBanding_ACT)
   ├─ revise ──> Assignment12 "MARKETING"
   ├─ reject ──> Utility8 "SendEmailReject/Ask" ─(Always)─> Assignment12
   └─ decline ─> End3

Decision28 "ISPKSASM"
   ├─ When IsPKSASM ─> Utility11 "INSERT TO PRODUCTION" (SaveToProduction_ACT)
   └─ Else ─────────> Utility7 "SET LIMIT_AKSEPTASI" (GetLimitAkseptasi_JUW_UW)

Decision16 "It is Group?"
   ├─ When IsGroup ─> SubProcess1 "POLICY" -> Flow InputInwardFacultativeRISlip
   └─ Else ────────> Utility6 "SendEmailBind" (SendEmailPolicy)
```

**Tahapan siklus new business `[terverifikasi]`:**
Marketing → (Life? → Medical Life) → simpan JSON offer → Group?/PKS?/UW Financial? →
tangga akseptasi berbasis limit → binding → **RI Slip** (`InputInwardFacultativeRISlip`) →
fac out/retro (`OfferFacRetro`) → produksi.

#### 1.2.1 Dua puluh satu workbasket, satu pola routing `[terverifikasi]`

Seluruh 21 Assignment memakai `<pyImplementation>WorkBasket` dan `<pyRouteTo>Custom` — **tidak ada
satu pun `WorkList`**. Ini kebalikan dari domain Claim (Tahap 3) yang memakai routing campuran.
→ memperkuat OQ-028.

Label shape yang muncul (label, **bukan** perilaku — `_METHOD.md` §2.2): `ADMIN`, `MARKETING`,
`MARKETING (BINDING)`, `TEAM LEADER`, `UNDERWRITING`, `SENIOR UNDERWRITING`,
`JUNIOR UNDERWITER A`, `JUNIOR UNDERWITER B`, `UNDERWRITING FINANCIAL`, `UNDERWRITING LIFE`,
`MEDICAL LIFE`, `MANAGER TEKNIK`, `KADIV TEKNIK`, `KADIV FINANCIAL`, `KADIV FACULTATIVE`,
`DEP HEAD UW FACULTATIVE`, `DEPTHEAD UNDERWRITING LIFE`, `DIREKTUR TEKNIK`, `DIREKTUR MARKETING`.

`<pyWorkBasket>` kosong dan routing `Custom` → pemetaan shape → nama workbasket **tidak terbaca**
(OQ-024). Nama workbasket yang ada di korpus ketiga modul facultative `[terverifikasi]`:

```
grep -rhoE "ReasFacIn[A-Za-z]*|ReasFacOut[A-Za-z]*" "NB FacIn" "RNW Fac In" "Endorsment Fac In" \
  --include="*.xml" | sort | uniq -c | sort -rn
```
→ `ReasFacInMarketing` (377), `ReasFacInAdmin` (143), `ReasFacOut` (132),
`ReasFacInUnderwriting` (114), `ReasFacInTeamLeader` (109), `ReasFacInSeniorUnderwriting` (105),
`ReasFacInUnderwritingFinancial` (91), `ReasFacInDepHeadUnderwriting` (91),
`ReasFacInManagerTeknik` (80), `ReasFacInGroupLeader` (80), `ReasFacInJuniorUnderwritingA` (78),
`ReasFacInTechnicalDirector` (68), `ReasFacOutAdmin` (47), `ReasFacOutHead` (42),
`ReasFacOutGroupLeader` (13), `ReasFacOutTechnicalDirector` (13).

#### 1.2.2 Enam belas shape yatim `[terverifikasi]`

`Decision14`, `Decision16`, `Decision21`, `Decision38`, `Decision42` **tidak dituju connector mana
pun**. Sebelas shape **tidak punya connector keluar**: `SubProcess1`, `SubProcess2`, `Utility1`,
`Utility3`, `Utility4`, `Utility5`, `Utility6`, `Utility7`, `Utility10`, `Utility11`, `Utility16`.

Perintah audit: dari daftar pasangan `<pyFrom>`/`<pyTo>`, hitung ID shape yang tidak pernah muncul
sebagai `pyTo` (masuk) atau `pyFrom` (keluar).

`[pertanyaan terbuka]` Pola sama dengan OQ-023 tetapi **skalanya jauh lebih besar** (16 shape vs 5).
Tidak dapat dipastikan apakah ini jalur mati, sub-flow yang dipanggil dari luar graf, atau ekspor
tidak lengkap. **Tidak ditebak.** → OQ-023 diperluas.

### 1.3 `InputInwardFacultativeRISlip` — siklus RI slip & konversi

**16 Assignment, 32 Decision, 8 Utility, 4 SubProcess, 119 connector** `[terverifikasi]`.

```
Start1 ─(Always)─> Decision8 "Is Life?"
                      ├─ When IsLife ─> Utility1 "SAVE JSON_POLICY" (SaveJsonPolicyFacIn_Act)
                      └─ Else ───────> Decision9 "Is it group?"
                                          ├─ When IsGroup ─> Assignment3 "MARKETING"
                                          └─ Else ────────> Decision27 "UW FINANCIAL?"
                                                               ├─ When IsTBonding ─> Assignment1 "UNDERWRITING FINANCIAL"
                                                               └─ Else ───────────> Utility5 "Check Spreading"
                                                                                       │ Activity CekLimitSpreading_Act
                                                                                       v
                                                                                    Assignment6 "TEAM LEADER"
```

Tangga akseptasi dan jalur fac out:

```
Decision24 "Accept?" ├─ confirm ─> Decision21 "Is it UW?"
                     │                ├─ When IsFlagUW ─> Utility8 "SET LIMIT_AKSEPTASI" (GetLimitAkseptasi_JUW_UW)
                     │                └─ Else ─────────> Utility1 "SAVE JSON_POLICY"
                     └─ reject ──> End2

Decision18 "Limit Akseptasi" / Decision33 "Limit Akseptasi" / Decision14 "Which team?"
   -> Assignment per peran (KADIV TEKNIK, DIREKTUR TEKNIK, KADIV FACULTATIVE, MANAGER TEKNIK,
      KADIV FINANCIAL, DIREKTUR MARKETING, DEP HEAD UW FACULTATIVE, SENIOR UNDERWRITING,
      JUNIOR UNDERWITER A, UNDERWRITING, UNDERWRITING FINANCIAL)

Decision4 "FAC OUT?" ├─ When IsFacRetro ─> Decision3 "INPUT FAC OUT?"
                     │                        ├─ When IsInputFacRetro ─> Utility1 SAVE JSON_POLICY
                     │                        └─ Else ─────────────────> SubProcess3 -> Flow OfferFacRetro
                     └─ Else ───────────> Utility1 "SAVE JSON_POLICY"

Decision12 "Is it group?" ├─ When IsGroup ─> Assignment2 "ADMIN (FAC OUT)"
                          │                     └─Action StartScreenFlowAuto─> SubProcess4 -> Flow OfferFacOut
                          └─ Else ────────> SubProcess1 "[Print R/I Slip]" -> Flow OfferFacRetro
                                               └─(Always)─> Utility2 "HIT SERVICE ARASAPAS"   <- BLOCKER §5.3

Decision30 "Err Konversi?" ├─ When IsSuccessHitService ─> Decision5 "FAC OUT?"
                           └─ Else ───────────────────> Utility6 "SetToInbox_ACT"

Decision5 ├─ When IsFacRetro ─> Decision6 "PRINT ?"
          │                        ├─ When IsNotPrintRISlip ─> Utility4 "UPDATE STS KONVERSI"
          │                        │                              │ Activity UpdateStsKonversiFacOut_Act
          │                        │                              └─(Always)─> SubProcess1 -> OfferFacRetro
          │                        └─ Else ──────────────────> END52
          └─ Else ───────────> END52
```

**Tahapan siklus RI slip `[terverifikasi]`:** Life?/Group?/UW Financial? → **Check Spreading** →
tangga akseptasi ("Limit Akseptasi" / "Which team?") → simpan JSON policy → cetak R/I slip / fac out
→ **HIT SERVICE ARASAPAS** → cek status konversi → update status konversi.

### 1.4 `OfferFacOut` — screen flow, satu-satunya `WorkList` di modul ini

`<pyStartActivity>Start62`, `<pyShapeType>` titik masuk = **`Data-MO-Event-Start-StartScreenFlow`**
`[terverifikasi]`. Hanya 4 shape:

```
Start62 ─> Decision1 "IsRISlip"
              ├─ When IsRISlip ─> AssignmentSF2 "Print R/I Slip"  [WorkList, route=Current operator]
              └─ Else ─────────> Assignment1   "Offer Fac Retro"  [WorkList]
```

`[terverifikasi]` Kedua Assignment berjenis `Data-MO-Activity-Assignment-WorkAction` dan memakai
`WorkList`. **Satu-satunya rule `Flow` di tiga modul facultative yang memakai `WorkList`, bukan
`WorkBasket`** — sisanya seluruhnya `WorkBasket`/`Custom`. → OQ-028.

Tidak ada connector keluar dari kedua Assignment `[terverifikasi]` — screen flow berakhir pada
layar, bukan pada shape End.

### 1.5 `OfferFacRetro` — varian NB FacIn / RNW Fac In (`bfd6252070`)

**3 Assignment, 4 Decision, 2 Utility, 1 Start** `[terverifikasi]`.

```
Start1 ─(Always)─> Decision1 "IsRISlip"
                      ├─ When IsRISlip ─> Utility2 "Send Email Print RI Slip" (SendEmailPolicy)
                      └─ Else ─────────> Assignment1 "ReasFacOutAdmin"  [WorkBasket, Custom]
                                            │ FlowAction: OfferFacOut
                                            v
                                         Decision2 "Accept?" (IsUWAccepted)
                                            ├─ confirm ─> Decision3 "Is it group?"
                                            │                ├─ When IsGroup ─> END52
                                            │                └─ Else ────────> Assignment3 "ReasFacOutHead"
                                            │                                     │ FlowAction: OfferFacOut
                                            │                                     v
                                            │                                  Decision4 "Accept?"
                                            │                                     ├─ confirm ─> END52
                                            │                                     └─ reject ──> Assignment1
                                            └─ reject ──> END52

Assignment2 "PRINT R/I SLIP" ─Action PrintRISlip_FlowAction─> Utility1 "InsertFacoutProd"
```

`[terverifikasi]` **Tangga persetujuan fac out di sini hanya DUA anak tangga**: `ReasFacOutAdmin` →
`ReasFacOutHead`. Varian `Endorsment Fac In` punya **empat** (lihat `Endorsment Fac In.md` §2).
Ini perbedaan perilaku yang terukur, bukan sekadar perbedaan teks.

`[terverifikasi]` Di flow ini `IsUWAccepted` hanya memakai **dua** hasil (`confirm`, `reject`) dari
enam yang tersedia (§2.1).

---

## 2. Status / state dan kode

### 2.1 `ProposalAcceptStatus` — enam hasil, nilai literal terbaca, **pemetaannya tidak**

`[terverifikasi]` Rule `ASM-FW-GISFW-WORK!ISUWACCEPTED` menggerbangi **22 shape Decision** di
`InputInwardFacultativeOffer` saja. Ia ada sebagai **dua tipe rule**:

| Path | Tipe | Isi yang terbaca |
| --- | --- | --- |
| `NB FacIn/When/IsUWAccepted.xml` | `RULE-OBJ-WHEN` | `[.ProposalAcceptStatus][=]["1"]` — **boolean** |
| `NB FacIn/DecisionTable/IsUWAccepted.xml` | `RULE-DECLARE-DECISIONTABLE` | properti masukan `ProposalAcceptStatus`; 5 kolom |

`[terverifikasi]` Enam hasil yang dideklarasikan tabel keputusan, dari `<pyTaskStatusXml>` di
`NB FacIn/DecisionTable/IsUWAccepted.xml`:

```
confirm   reject   ask   banding   revise   decline
```

Perintah audit:
```
grep -o "<pyTaskStatusXml>[^<]*" "NB FacIn/DecisionTable/IsUWAccepted.xml"
```

Keenamnya **persis** nilai `<pyExpression>` pada connector ber-`<pyConditionType>Status` di flow.
Karena `RULE-OBJ-WHEN` hanya dapat mengembalikan benar/salah, `[dugaan]` yang dipakai flow adalah
varian `DecisionTable`. Ini **mempersempit OQ-026**, tidak menutupnya — flow tidak menyebut tipe.

**Nilai literal yang ditulis ke `ProposalAcceptStatus` `[terverifikasi]`:**

| Rule `RULE-OBJ-MODEL` (Data Transform) | `SET` | Efek samping yang terbaca |
| --- | ---: | --- |
| `DataTransform/SetAkseptasiProposal.xml` | **1** | — |
| `DataTransform/SetRejectProposal.xml` | **2** | `.OfferFacIn.ConfirmBinding=0`, `.OfferFacIn.ReceivedRiSlip=false` |
| `DataTransform/SetAskProposal_DT.xml` | **3** | — |
| `DataTransform/SetAkseptasiCeding_DT.xml` | **3** | — |
| `DataTransform/SetBandingProposal_DT.xml` | **4** | `.OfferFacIn.IsBanding="true"`, **`.LetterNo = .BandingTo`**, `FlagBanding.CARI11 = .BandingTo` |
| `DataTransform/SetDeclineProposal_DT.xml` | **7** | — |
| `DataTransform/SetReviseProposal.xml` | **9** | `.OfferFacIn.ConfirmBinding=0`, `.OfferFacIn.ReceivedRiSlip=false` |
| `DataTransform/SetAkseptasiProposal_DT.xml` | **1, 2, 3, 9, 4** (berkondisi `WHEN`) | juga men-set `QuotationList(1).Email` |

Perintah audit (memasangkan `<pyPropertiesName>` dengan `<pyPropertiesValue>`):
```
awk '/<pyPropertiesName>/{n=$0;gsub(/<[^>]*>/,"",n)} /<pyPropertiesValue>/{v=$0;gsub(/<[^>]*>/,"",v);
     if(n!=""){print n" := "v; n=""}}' "NB FacIn/DataTransform/SetBandingProposal_DT.xml"
```

**Yang TIDAK dapat dinyatakan:** pemetaan nilai `1/2/3/4/7/9` → `confirm/reject/ask/banding/
revise/decline`. Baris tabel keputusan (kondisi → hasil) **tidak ada di ekspor** — §5.1. Korelasi
nama rule dengan nama hasil bersifat sugestif, dan `_METHOD.md` §2.2 melarang memakai nama sebagai
bukti perilaku. Arti bisnisnya **belum terverifikasi** (OQ-020). → **OQ-043**, **OQ-044**.

`[terverifikasi]` **Satu nilai, dua pemakai berbeda nama.** `NB FacIn/When/IsFacout.xml` berisi
`[pyWorkPage.ProposalAcceptStatus][=][4]` — nilai yang sama yang ditulis `SetBandingProposal_DT`.
Jadi `4` menggerbangi jalur bernama "facout" **dan** ditulis oleh rule bernama "banding". Mana yang
dimaksud **tidak dapat diputuskan dari korpus**. → **OQ-044**.

### 2.2 `LetterNo` — field nomor surat dipakai sebagai token routing

`[terverifikasi]` Sepuluh token literal diuji terhadap `pyWorkPage.LetterNo` di korpus facultative:

```
grep -rhoE "LetterNo\]\[&amp;#61;\]\[&amp;quot;[A-Z0-9]+" "NB FacIn" "RNW Fac In" "Endorsment Fac In" \
  --include="*.xml" | sed 's/.*quot;//' | sort | uniq -c | sort -rn
```

| Token | Muncul | Rule `When` |
| --- | ---: | --- |
| `KADIVTEKNIK` | 6 | `When/ToKadivTeknik.xml` |
| `KADIVFACULTATIVE` | 6 | `When/ToKadivFacultative.xml` |
| `DIREKTURTEKNIK` | 6 | `When/ToDirTeknik.xml` |
| `DIREKTURMARKETING` | 6 | `When/ToDirMarketing.xml` |
| `DEPHEADUNDERWRITER` | 6 | `When/ToDepHeadUW.xml` |
| `SENIORUW` | 3 | `When/ToSeniorUW.xml` |
| `MANAGERTEKNIK` | 3 | `When/ToManagerTeknik.xml` |
| `KADIVFINANCIAL` | 3 | `When/ToKadivFin.xml` |
| `DEPTHEADUWLIFE` | 2 | `When/ToDeptHeadUWLife.xml` |
| `TREATYINDEPTHEAD` | 2 | `When/ToTREATYDEPTHEAD.xml` (dicatat di `NB Treaty In.md` §2) |

`[terverifikasi]` **Rantainya lengkap dan terbaca**: `SetBandingProposal_DT` men-set
`.LetterNo = .BandingTo`, lalu `Decision23`/`Decision38`/`Decision42` memilih Assignment berdasarkan
nilai `LetterNo`. Jadi **field nomor surat adalah mekanisme routing persetujuan**, bukan sekadar
penanda dokumen. Pola sama dengan `OperatorID.pyTelephone` (OQ-027) di domain treaty: field bisnis
dipakai menyimpan kode peran. → **OQ-045**.

`When/ToDepHeadUW.xml` juga menguji `pyWorkPage.OfferFacIn.IsB2B != "ASM"` `[terverifikasi]`;
arti `IsB2B` dan nilai `"ASM"` **belum terverifikasi**.

### 2.3 Kode dan flag lain

| Properti | Nilai literal | Diuji/diubah oleh | Arti |
| --- | --- | --- | --- |
| `.OfferFacIn.QuotationData.IsGroup` | `"Group"` | `When/IsGroup.xml` | **belum terverifikasi** |
| `.OfferFacIn.IsFacRetro` | `1` | `When/IsFacRetro.xml` | gerbang jalur fac out/retro |
| `.OfferFacIn.IsInputFacRetro` | `1` | `When/IsInputFacRetro.xml` | "Is Already Input Fac Retro" |
| `.OfferFacIn.IsRISlip` | `!= 1` | `When/IsNotPrintRISlip.xml` | gerbang cetak slip |
| `pyWorkPage.Quotation.PolicyType` | `2` | `When/IsDeclarationPolicy.xml` | **belum terverifikasi** |
| `pyWorkPage.IsSpecialCase` | `true` | `When/IsSpecialCase.xml` | **belum terverifikasi** |
| `.OfferFacIn.QuotationData.TeamGroup` | `5` | `When/IsTBonding.xml` | **belum terverifikasi** |
| `pyWorkPage.Quotation.BusinessOldId` | `"L1"`, `"L2"`, … | `When/IsLife.xml` | kode lini bisnis — **sama dengan OQ-038** |
| `.StatusService.StsKonversiFacIn` | `1` | `When/IsSuccessHitService.xml` | status konversi |
| `.StatusService.StsKonversiFacOut` | `""`, `1` | `When/IsSuccessHitService.xml` | status konversi |
| `.OfferFacIn.ProRateType` | `3`, `4` | precondition `CountPremiNusareRetro_Act` + ~20 rule | **belum terverifikasi** |

`[terverifikasi]` `When/IsLife.xml` memakai kode `L1`… yang sama dengan yang ditemukan di
`Komite Claim FacIn` (OQ-038) — bukti bahwa daftar lini bisnis `L1`…`L11` dipakai lintas domain,
bukan khas komite.

### 2.4 Empat `When` yang kondisinya tidak terbaca

`[terverifikasi]` `<pyLabel>` hanya berisi template kosong `[first value][relation][second value]`:

| Rule | Menggerbangi |
| --- | --- |
| `NB FacIn/When/IsPKSASM.xml` | `Decision27`, `Decision28` — jalur `UNDERWRITING` vs `UW FINANCIAL`, dan `Utility11 "INSERT TO PRODUCTION"` |
| `NB FacIn/When/LetterNoNull.xml` | `Decision33` — jalur banding kembali ke `MARKETING (BINDING)` |
| `NB FacIn/When/ToJUW_A.xml` | `Decision41`, `Decision42` — pemilihan tim underwriter |
| `NB FacIn/When/ToUW.xml` | `Decision41`, `Decision42` — pemilihan tim underwriter |

`PKS` dan `ASM` **kepanjangan belum terverifikasi**. Ini keluarga keempat pola When tak terbaca di
D2 (setelah `IsPEGAPROD`, `IsSendtoMedical`/`IsSPK`, `ISCLM*`) → OQ-029.

---

## 3. Objek Oracle yang disentuh

Perintah audit:
```
awk -F'\t' '$1 ~ "^NB FacIn/" && $2=="RULE-CONNECT-SQL"{print $8}' all-rules.tsv \
  | grep -oE "tables=[^;]*" | sed 's/tables=//' | tr ',' '\n' | sort | uniq -c | sort -rn
```

| Objek | Rule perujuk |
| --- | ---: |
| **`JSON_POLIS`** | 18 |
| `FACINPRODUCTION` | 7 |
| `RW`, `ACCUMULATION` | 6 masing-masing |
| `TREATYBUSINESS`, `POOLDATA.REINSURANCETYPE`, `POOLDATA.JSON_POLIS`, `OCCUPATION`, `CURRENCY` | 5 masing-masing |
| `PROPORTIONALARRG`, `POOLDATA.M_LIMIT_PROPERTYY`, `POOLDATA.CLIENT`, `MARKETINGOFFICER` | 4 masing-masing |
| `T_STORAGE_IMAGE`, `POOLDATA.M_LIMIT_NONPROPANDENGG` | 3 masing-masing |

`[terverifikasi]` `JSON_POLIS` adalah objek paling dirujuk — **18 rule**, lebih banyak dari modul
mana pun yang ditelusur sejauh ini. Struktur isinya tidak ada di korpus → OQ-012.

`[terverifikasi]` Dua tabel limit: `POOLDATA.M_LIMIT_PROPERTYY` (4) dan
`POOLDATA.M_LIMIT_NONPROPANDENGG` (3). Ini **sumber limit dari database** — pasangan dari pola limit
ter-hardcode yang ditemukan di `Claim Non Prop` (OQ-040).

### 3.1 Sentuhan ke domain treaty `[terverifikasi]`

Modul facultative ini merujuk `TREATYBUSINESS` (5 rule) dan `PROPORTIONALARRG` (4 rule), dan
memanggil procedure **`POOLDATA.PEGA_TREATY_IN`** serta **`POOLDATA.PEGA_JSON_POLIS_TREATYIN`** —
procedure yang sama yang dipanggil `NB Treaty In` (`NB Treaty In.md` §5.1). Ditambah rule `Flow`
`InputRealizationTreatyIn` yang identik.

Pola sejajar dengan OQ-042 (Claim → Treaty), tetapi di sini **ada penjelasan struktural**:
`InputQuotation` (§1.1) memang bercabang ke kedua lini. **Dicatat untuk D4; penetapan konteks bukan
di sini.**

---

## 4. Integrasi eksternal

`[terverifikasi]` Tiga rule `RULE-CONNECT-REST`, **seluruhnya berbasis konfigurasi**:

| File | `pyServiceName` | `pyBaseURLSelectionType` | `pyBaseURLSetting` |
| --- | --- | --- | --- |
| `ConnectREST/ServiceGoogle.xml` | `ServiceGoogle` | `SETTING` | `LinkService!LinkService` |
| `ConnectREST/convertJsonNusareToProduction.xml` | `convertJsonNusareToProduction` | `SETTING` | `LinkService!LinkService` |
| `ConnectREST/getPremiumPaidOn.xml` | `getPremiumPaidOn` | `SETTING` | `LinkService!LinkService` |

`[terverifikasi]` **Tidak ada URL literal endpoint.** Satu-satunya URL dalam ketiga file adalah
`<pyHelpURI>` — tautan dokumentasi Pega, bukan endpoint. Perintah audit:

```
grep -ohE "<[a-zA-Z]+>https?://[^<]*" "NB FacIn/ConnectREST/"*.xml | sort -u
```

Ini **berbeda dari `Claim Fac In`** yang punya 8 `ConnectREST` dengan satu URL literal (OQ-018).
Base URL sesungguhnya ada di `RULE-ADMIN-SYSTEM-SETTINGS` `LinkService` — nilainya **tidak disalin**
ke sini (OQ-018).

Integrasi lain yang tersentuh flow, sebagai **shape**, bukan `ConnectREST`:

| Shape | Activity yang dirujuk | Catatan |
| --- | --- | --- |
| `Utility2` (RI Slip) "HIT SERVICE ARASAPAS" | `serviceInsertArasapas_act` | **terblokir** — §5.3 |
| `Utility4` (RI Slip) "UPDATE STS KONVERSI" | `UpdateStsKonversiFacOut_Act` | status konversi; isi belum dibaca |
| `Utility6`/`Utility8` (Offer) "SendEmail…" | `SendEmailPolicy` | memanggil `SendEmailWithAttachments` |
| `Utility11` (Offer) "INSERT TO PRODUCTION" | `SaveToProduction_ACT` | memanggil `serviceInsertArasapas_act` + `serviceInsertArasapasEDM_act` |

`[terverifikasi]` `Activity/SaveToProduction_ACT.xml` (94.952 byte) punya 10 langkah:
`Call SaveJsonPolicyFacIn_Act`, `Call serviceInsertArasapas_act`, `RDB-List`, `Call SetToInbox_ACT`,
`Call ASMForceCaseClose`, `Call SaveEDMToJsonPolicy_Act`, `Call serviceInsertArasapasEDM_act`,
`RDB-List`, `Call SetToInbox_ACT`, `Call ASMForceCaseClose`.

`[terverifikasi]` `Activity/SetToJsonOffer_ACT.xml` punya 4 langkah: `Call SaveJsonOfferFacIn_Act`,
`Call SaveJsonPolicyFacIn_Act`, `Call serviceInsertArasapas_act`, `Call SetToInbox_ACT`.

**Kasir dan Google Storage:** `[terverifikasi]` procedure `POOLDATA.GET_TOKEN_STORAGE` dipanggil
(§5.2) dan `ServiceGoogle` ada sebagai `ConnectREST`; kaitan keduanya **belum terverifikasi**. Tidak
ada rule bernama Kasir di modul ini (`ls "NB FacIn" | grep -i kasir` → kosong) — berbeda dari
`Claim Fac In` yang punya `SendAcceptationToKasir`.

---

## 5. Batas pengetahuan

### 5.1 Baris tabel keputusan **tidak ada di ekspor** — temuan baru

`[terverifikasi]` **Tidak satu pun dari 49 file `DecisionTable` di seluruh korpus memuat baris
keputusannya** (kondisi → hasil). Yang ada hanya: nama properti masukan (`<pyRuleName>` di
`pxRuleReferences`), daftar hasil yang diizinkan (`<pyTaskStatusXml>`), lebar kolom
(`<pyColumnPreferences>`), dan metadata versi.

Perintah audit:
```
find . -path ./OUTPUT_HASIL_RNM -prune -o -type f -path "*/DecisionTable/*.xml" -print | wc -l
# -> 49
find . -path ./OUTPUT_HASIL_RNM -prune -o -type f -path "*/DecisionTable/*.xml" -print \
| while IFS= read -r f; do grep -oE "pyAllowedResults|pyReturnValue|pyDecisionTableData|pyRowData|pyConditions" "$f"; done \
| sort | uniq -c
# -> hanya "pyAllowedResults", dan itu bagian dari <pzIsAllowedResultsPropSetExist>
```

**Dampaknya besar dan bukan sekadar detail ekspor**: `IsUWAccepted` menentukan **setiap** hasil
persetujuan di flow terbesar korpus, dan isi tabelnya **tidak diketahui**. Hal yang sama berlaku
untuk `isApproved` (OQ-026), `IsLifeAccepted`, `MappingCoverage`, `BusinessType_DeT`,
`LicensePlatRegion_DeT`, `MappingOutgoIndex`, `MappingAdditionalCoverageIndex`, `SetUploadHubAW1-3`,
dan 40 tabel lain. → **OQ-043**.

### 5.2 Stored procedure — 29 dipanggil, isinya tidak diketahui

`[terverifikasi]`:

```
POOLDATA.PEGA_JSON_POLIS_TREATYIN   POOLDATA.PEGA_TREATY_IN          POOLDATA.PEGA_M_JSON_OFFER
POOLDATA.INSERTJSONPOLIS            POOLDATA.INSERTJSONPOLISMONITORING
POOLDATA.INSERTUPDATECEDINGPRODUCTION POOLDATA.INSERTUPDATERISKADDRESS
POOLDATA.GENERATE_FACRETRO_NO       POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER
POOLDATA.GETCURRENCYSTANDARD        POOLDATA.GET_TOKEN_STORAGE
POOLDATA.PEGA_DELETE_ERROR_KONVERSI POOLDATA.PEGA_MARKETINGOFFICER
POOLDATA.PEGA_M_ACCUMULATION_LIFE   POOLDATA.FACINFORBACKUP          POOLDATA.RDBINSERTCLIENT
POOLDATA.RDBMASTERACCUMULATEDTYPE   POOLDATA.RDBMASTERACCUMULATION   POOLDATA.RDBMASTERBRANCH
POOLDATA.RDBMASTERCITY              POOLDATA.RDBMASTERDISTRICT       POOLDATA.RDBMASTERNATION
POOLDATA.RDBMASTERPROVINCE          POOLDATA.RDBMASTERRW             POOLDATA.RDBMASTERSHIP
GENERAL.F_GET_NM_ASURADUR           MBU.F_CEK_HURUF
DBMS_LOB.CREATETEMPORARY            UTL_MATCH.EDIT_DISTANCE_SIMILARITY
```

**Isinya tidak ada di korpus** (OQ-002). `POOLDATA.GENERATE_FACRETRO_NO` — penomoran fac retro — dan
`POOLDATA.PEGA_M_JSON_OFFER` hanya muncul di modul facultative `[terverifikasi]`.

Dua skema yang belum tercatat sebelumnya muncul di sini: `GENERAL`, `MBU` → OQ-016.

### 5.3 `serviceInsertArasapas_act` — blocker OQ-025 **meluas**, dan identitas sasarannya kini terverifikasi

`[terverifikasi]` Ada **dua identitas rule berbeda** dengan nama file yang sama:

| `pxInsName` | Ukuran | Hash | Ada di |
| --- | ---: | --- | --- |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN!SERVICEINSERTARASAPAS_ACT` | 18.013 | `b013027e0c` | NB Treaty In, **NB FacIn**, **RNW Fac In** |
| `ASM-FW-GISFW-WORK!SERVICEINSERTARASAPAS_ACT` | 232.000 | `ab3ae3c9ba` | EDM Treaty In |
| `ASM-FW-GISFW-WORK!SERVICEINSERTARASAPAS_ACT` | 231.975 | `7314b6c49c` | **Endorsment Fac In** |

Perintah audit:
```
find . -path ./OUTPUT_HASIL_RNM -prune -o -type f -iname "serviceInsertArasapas_act.xml" -print \
| while IFS= read -r f; do grep -o "<pxInsName>[^<]*" "$f" | head -1; sh nhash.sh "$f"; done
```

`[terverifikasi]` `NB FacIn/Activity/serviceInsertArasapas_act.xml` berisi **satu langkah**:
`Call serviceInsertArasapas_act` dengan `<pyStepsObjectName>pyWorkPage` dan — ini yang baru —
**`<pyStepsClassName>ASM-FW-GISFW-Work`**.

**Ini menaikkan status temuan OQ-025 dari `[dugaan]` menjadi `[terverifikasi]`.** Di `NB Treaty In.md`
§5.2 class sasaran panggilan hanya dapat **diduga** dari class `pyWorkPage`; di sini class sasaran
**tertulis eksplisit** di tag langkah. Sasarannya `ASM-FW-GISFW-WORK!SERVICEINSERTARASAPAS_ACT` —
identitas berkonflik (OQ-011 entri #415, 2 varian, 2 isi berbeda) yang **tidak ada di NB FacIn**.

Sesuai `_METHOD.md` §1.1 aturan 4: telusur **dihentikan** di shape `Utility2` "HIT SERVICE ARASAPAS"
(§1.3) dan di dalam `SetToJsonOffer_ACT` / `SaveToProduction_ACT`. Varian dari modul lain **tidak
dipinjam**. → OQ-025 diperluas ke NB FacIn dan RNW Fac In.

### 5.4 Spreading, capacity, dan scoring — 38 activity, tak satu pun habis dibaca

`[terverifikasi]` Modul ini memuat **38 activity** bernama spreading/capacity/scoring. Perintah
audit: `ls -l "NB FacIn/Activity/" | grep -icE "spread|capacit|scoring"` → 38.

| Activity | Ukuran | Terbaca |
| --- | ---: | --- |
| `SumTSIPremiSpreadedRNM_FIRE_Act.xml` | **996.052** — terbesar di seluruh D2 | **belum** |
| `SumTSIPremiSpreadedRNM_ANEKA_Act.xml` | 748.342 | **belum** |
| `SumTSIPremiSpreadedRNM_Act.xml` | 516.395 | **belum** |
| `CheckSpreadingProtectAnekaGolf_ACT.xml` | 398.150 | **belum** |
| `CopyToAllLocSpreading_ACT.xml` | 394.959 | **belum** |
| `CheckSpreadingProtectFire_ACT.xml` | 366.123 | **belum** |
| `CopyToAllSpreadingAnekaGolf_ACT.xml` | 341.207 | **belum** |
| `CopyToAllSpreading_ACT.xml` | 321.542 | **belum** |
| `SetDataScoringRisk_act.xml` | 304.035 | **belum** |
| `CekLimitSpreading_Act.xml` | 296.620 | **belum** — dipanggil `Utility5` §1.3 |
| `ValueScoringRisk_act.xml` | 267.195 | **belum** |
| `ScoringResult.xml` | 143.652 | **belum** |
| `SumTreatyCapacity_Act.xml` | 143.327 | **belum** |

**Rumus spreading, capacity, dan scoring tidak dinyatakan di artefak ini.** Ini **batas cakupan
telusur** (volume), berbeda dari batas pengetahuan korpus (SP/db-link). Pola per lini bisnis
terbaca dari penamaan — FIRE, ANEKA, GOLF, MARINECARGO, MBU, PA, TRAVEL — tetapi `_METHOD.md` §2.2
melarang memakai nama sebagai bukti perilaku.

### 5.5 Satu rumus yang **terbaca penuh** — `CountPremiNusareRetro_Act`

`[terverifikasi]` Berbeda dari §5.4, `NB FacIn/Activity/CountPremiNusareRetro_Act.xml`
(`ASM-FW-GISFW-WORK!COUNTPREMINUSARERETRO_ACT`, 64.707 byte, 5 langkah `Property-Set`) **dapat
dibaca seluruhnya**. Hash `8f1268a660` — **identik di ketiga modul facultative**.

```
local.sharernm1 := 0.45
local.sharernm2 := 0.05
local.tampungpreminet1 := 0

bila .TSILiability <= "3000000000":
    .PremiLifeNusantaraRe := .GrossPremiumRetro * local.sharernm1

bila .TSILiability >  "3000000000":
    local.tampungpreminet1 := .GrossPremiumRetro * local.sharernm1
    .PremiLifeNusantaraRe  := (.GrossPremiumRetro * local.sharernm2) + local.tampungpreminet1

.PremiumRetro := @toDecimal(.GrossPremiumRetro) - @toDecimal(.PremiLifeNusantaraRe)

bila pyWorkPage.OfferFacIn.ProRateType != 3:
    .PremiLifeNusantaraRe := .Premium - .PremiumRetro
```

Perintah audit:
```
awk '/<PropertiesName>/{n=$0;gsub(/<[^>]*>/,"",n)} /<PropertiesValue>/{v=$0;gsub(/<[^>]*>/,"",v);
     if(n!=""){print n" := "v;n=""}}' "NB FacIn/Activity/CountPremiNusareRetro_Act.xml"
grep -oE "<pyStepsPreCondParamsWhen>[^<]*" "NB FacIn/Activity/CountPremiNusareRetro_Act.xml"
```

`[terverifikasi]` Fakta yang tercatat apa adanya:
- Pangsa **0,45** dan **0,05** ter-hardcode; **mata uang tidak disebut**; **arti belum terverifikasi**.
- Ambang **3.000.000.000** ter-hardcode dan dibandingkan sebagai **string berkutip**
  (`<="3000000000"`), bukan sebagai angka.
- `ProRateType != 3` mengubah cara `PremiLifeNusantaraRe` dihitung.

→ **OQ-046**.

### 5.6 Tabel rate di luar sistem

`[terverifikasi]` Delapan objek diakses lewat database link `@ASMD.SINARMAS.CO.ID`:
`M_EQS_RATE`, `M_TERORISME_RATE`, `M_RSMD_RATE`, `M_FLEXAS_RATE`, `M_FLOOD_AREA`,
`FIRE.M_FLOOD_RATE`, `FIRE.M_BI_INDEMNITY`, `LST_KURS_STANDARD`.

Ini **cakupan penuh** OQ-017 untuk modul ini (D1 sudah menyapu 100% korpus: hanya 3 modul
facultative yang memakainya). Isi tabel rate dan kurs standar **tidak ada di korpus** — perhitungan
premi facultative tidak dapat direplikasi tanpa akses ke database itu.

---

## 6. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `NB FacIn/`) | OQ-011? |
| --- | --- | --- |
| `ASM-FW-GISFW-WORK-NB` / `INPUTQUOTATION` / `RULE-OBJ-FLOW` | `Flow/InputQuotation.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `INPUTINWARDFACULTATIVEOFFER` / `RULE-OBJ-FLOW` | `Flow/InputInwardFacultativeOffer.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `INPUTINWARDFACULTATIVERISLIP` / `RULE-OBJ-FLOW` | `Flow/InputInwardFacultativeRISlip.xml` | tidak (identik dgn RNW) |
| `ASM-FW-GISFW-WORK` / `OFFERFACOUT` / `RULE-OBJ-FLOW` | `Flow/OfferFacOut.xml` | tidak (identik dgn RNW) |
| `ASM-FW-GISFW-WORK` / `OFFERFACRETRO` / `RULE-OBJ-FLOW` | `Flow/OfferFacRetro.xml` | **YA — #338** (varian NB/RNW) |
| `ASM-FW-GISFW-WORK` / `INPUTREALIZATIONTREATYIN` / `RULE-OBJ-FLOW` | `Flow/InputRealizationTreatyIn.xml` | tidak (identik dgn NB Treaty In — dipakai ulang) |
| `ASM-FW-GISFW-WORK` / `ISUWACCEPTED` / `RULE-DECLARE-DECISIONTABLE` | `DecisionTable/IsUWAccepted.xml` | tidak (isi baris tak ada — §5.1) |
| `ASM-FW-GISFW-WORK` / `ISUWACCEPTED` / `RULE-OBJ-WHEN` | `When/IsUWAccepted.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `ISFACOUT` / `RULE-OBJ-WHEN` | `When/IsFacout.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `ISGROUP`, `ISGROUPCREATE`, `ISTBONDING` / `RULE-OBJ-WHEN` | `When/` | tidak — **guard identitas, §7.1** |
| `ASM-FW-GISFW-WORK` / `ISFACRETRO`, `ISINPUTFACRETRO`, `ISNOTPRINTRISLIP`, `ISSUCCESSHITSERVICE` / `RULE-OBJ-WHEN` | `When/` | tidak |
| `ASM-FW-GISFW-WORK` / `ISLIFE`, `ISSPECIALCASE`, `ISDECLARATIONPOLICY` / `RULE-OBJ-WHEN` | `When/` | tidak |
| `ASM-FW-GISFW-WORK` / `ISPKSASM`, `LETTERNONULL`, `TOJUW_A`, `TOUW` / `RULE-OBJ-WHEN` | `When/` | tidak — **kondisi tak terbaca, §2.4** |
| `ASM-FW-GISFW-WORK` / `TOKADIVTEKNIK`, `TOKADIVFIN`, `TOKADIVFACULTATIVE`, `TODIRTEKNIK`, `TODIRMARKETING`, `TOMANAGERTEKNIK`, `TODEPHEADUW`, `TOSENIORUW`, `TODEPTHEADUWLIFE` / `RULE-OBJ-WHEN` | `When/` | tidak |
| `ASM-FW-GISFW-WORK` / `GETLIMITAKSEPTASI_ACTFLOW` / `RULE-OBJ-ACTIVITY` (539.687 B) | `Activity/GetLimitAkseptasi_ActFlow.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `GETLIMITAKSEPTASI_JUW_UW` / `RULE-OBJ-ACTIVITY` (197.433 B) | `Activity/GetLimitAkseptasi_JUW_UW.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `GETLIMITAKSEPTASILIFE_ACT` / `RULE-OBJ-ACTIVITY` (62.560 B) | `Activity/GetLimitAkseptasiLife_Act.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `SETBANDING_ACT` / `RULE-OBJ-ACTIVITY` (278.875 B) | `Activity/SetBanding_ACT.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `SETTOJSONOFFER_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/SetToJsonOffer_ACT.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `SAVETOPRODUCTION_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/SaveToProduction_ACT.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `SAVEJSONOFFERFACIN_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/SaveJsonOfferFacIn_Act.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `COUNTPREMINUSARERETRO_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/CountPremiNusareRetro_Act.xml` | tidak (identik 3 modul) |
| `ASM-FW-GISFW-WORK` / `SENDEMAILPOLICY` / `RULE-OBJ-ACTIVITY` (191.489 B) | `Activity/SendEmailPolicy.xml` | tidak |
| `@BASECLASS` / `SETTICKET` / `RULE-OBJ-ACTIVITY` | `Activity/SetTicket.xml` | tidak — class bawaan Pega (OQ-009) |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN` / `SERVICEINSERTARASAPAS_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/serviceInsertArasapas_act.xml` | tidak (tapi **memanggil** identitas berkonflik — §5.3) |
| `ASM-FW-GISFW-WORK` / 8 Data Transform `SET*PROPOSAL*` / `RULE-OBJ-MODEL` | `DataTransform/` | sebagian |
| (FlowAction) `INWARDFACULTATIVE`, `INWARDFACULTATIVE_ISUW`, `OFFERFACOUT`, `PRINTRISLIP_FLOWACTION`, `STARTSCREENFLOWAUTO` | `FlowAction/` | sebagian |
| 3 rule `RULE-CONNECT-REST` (§4) | `ConnectREST/` | sebagian (`SERVICEGOOGLE` berkonflik) |

**45 rule ditelusur** (6 Flow — 1 dipakai ulang tanpa telusur baru; 1 DecisionTable; 22 When;
10 Activity; 8 DataTransform; 3 ConnectREST; 5 FlowAction dihitung sebagai kelompok).

`[terverifikasi]` `Activity/GetLimitAkseptasi_ActFlow.xml` memanggil **12 RDB-List berbeda**, semua
di class `ASM-FW-GISFW-Int-policyjson` dengan prefix `ASM`:
`GetLimitAkseptasi_SQL`, `GetLimitAkseptasiBanding_SQL`, `GetLimitAkseptasiPreferedComm_SQL`,
`GetLimitAkseptasiNonPrefer_SQL`, `GetLimitAkseptasiNonPreferBanding_SQL`,
`GetLimitAkseptasiNonFire_SQL`, `GetLimitAkseptasiNonFireBanding_SQL`,
`GetLimitAkseptasiKreditNCL_SQL`, `GetLimitAkseptasiKreditCL_SQL`, `GetLimitAkseptasiBond_SQL`,
`GetLimitAccEngineeringUW_SQL`, `GetLimitAccEngineeringBanding_SQL`.

Perintah audit:
`grep -oE "<pyRuleName>[^<]*" "NB FacIn/Activity/GetLimitAkseptasi_ActFlow.xml" | sort | uniq -c`

`[terverifikasi]` Limit akseptasi facultative **diambil dari database**, berbeda per lini bisnis
(Prefered Comm / Non Prefer / Non Fire / Kredit NCL / Kredit CL / Bond / Engineering) dan punya
**varian "Banding" tersendiri**. Langkah pertama ketiga activity limit adalah
`Call CountTotalTSIPremiNusaRe_Act`. → memperkaya OQ-040.

Rule **dirujuk tapi belum ditelusur**: 38 activity spreading/scoring (§5.4), 12 RDB-List
`GetLimitAkseptasi*`, `SetBusinessType_Act`, `pzChangeStageWrapper`, `UpdateStsKonversiFacOut_Act`,
`SetToInbox_ACT`, `ASMForceCaseClose`, `InsertFacoutProd`, `SaveJsonPolicyFacIn_Act`,
`CountTotalTSIPremiNusaRe_Act`, `SaveOfferProduction_Act`, `SendEmailWithAttachments`.

---

## 7. Pertanyaan terbuka

| OQ | Ringkas | Pemilik |
| --- | --- | --- |
| **OQ-043** | Baris tabel keputusan tidak ada di ekspor — 49 `DecisionTable`, termasuk `IsUWAccepted` yang menggerbangi 22 Decision | pemilik export + Product+UW |
| **OQ-044** | `ProposalAcceptStatus = 4` ditulis `SetBandingProposal_DT` tetapi diuji `When/IsFacout` — satu nilai, dua nama pemakai | Product+UW |
| **OQ-045** | `pyWorkPage.LetterNo` dipakai sebagai token routing persetujuan (10 kode peran literal) | IAM + Product+UW |
| **OQ-046** | Pangsa retro 0,45/0,05 dan ambang 3 miliar ter-hardcode; mata uang tidak disebut; ambang dibandingkan sebagai string | Product+UW + Actuarial |

OQ dikuatkan: **OQ-002** (§5.2 — 29 SP, 2 skema baru), **OQ-011** (§1 tabel, §5.3),
**OQ-012** (§3 — `JSON_POLIS` 18 rule), **OQ-015** (§1), **OQ-016** (`GENERAL`, `MBU`),
**OQ-017** (§5.6), **OQ-018** (§4 — nihil URL literal di sini), **OQ-020** (§2.1),
**OQ-023** (§1.1, §1.2.2 — 16 shape yatim), **OQ-024** (§1.2.1), **OQ-025** (§5.3 — **meluas dan
naik status ke terverifikasi**), **OQ-026** (§2.1 — dipersempit), **OQ-028** (§1.2.1, §1.4),
**OQ-029** (§2.4 — keluarga keempat), **OQ-038** (§2.3), **OQ-040** (§6 — 12 rule limit dari DB).

### 7.1 Guard identitas orang `[terverifikasi]`

Dicatat sesuai `_METHOD.md` §1.4: rule, tag, dan efeknya. **Nilai nama orang tidak disalin.**

| Rule | Tag yang memuat identitas | Jumlah identitas | Efek |
| --- | --- | ---: | --- |
| `When/IsGroup.xml` | `OperatorID.pyUserIdentifier` | **3** | ikut menentukan cabang `Decision1` "Is it group?" |
| `When/IsGroupCreate.xml` | `pyWorkPage.pxCreateOperator` | **2** | menentukan `Decision30` → Assignment `MARKETING` |
| `When/IsTBonding.xml` | `.OfferFacIn.QuotationData.MarketingName` | 1 | menentukan `Decision25`/`Decision41` → `UNDERWRITING FINANCIAL` |

`[terverifikasi]` Ketiga rule **identik isinya** di NB FacIn, RNW Fac In, dan Endorsment Fac In
(`IsGroup` hash `2cff7b9a67`) — satu himpunan guard identitas dipakai ketiga siklus.

`[terverifikasi]` Label `<pyMOName>` pada shape `Decision30` di **ketiga** rule `Flow` utama
facultative memuat nama orang. Dicatat sebagai temuan; nilainya tidak disalin.

**Pola ketiga yang baru** untuk OQ-021: `.MarketingName` — sebuah **field data bisnis**, bukan field
operator — dipakai sebagai guard percabangan. Berbeda dari OQ-021 (`pyUserIdentifier`/`pyPosition`)
dan OQ-027 (`pyTelephone`).

Perintah audit:
```
grep -rlE "pyUserIdentifier|pxCreateOperator|MarketingName" "NB FacIn/When" --include="*.xml"
```
