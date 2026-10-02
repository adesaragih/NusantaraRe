# Telusur Flow — Endorsment Fac In

STEP D2, Tahap 4 konteks #15. Ditelusur 2026-09-13. Konvensi: `_METHOD.md`.
**2.061 file** (`find "Endorsment Fac In" -name '*.xml' | wc -l`).

Artefak ini **mengutamakan yang khas endorsement/addendum**. Untuk perilaku yang sama dengan siklus
new business, rujuk **`NB FacIn.md`** — tidak diulang di sini.

**Dua rule `Flow`**, keduanya class `ASM-FW-GISFW-WORK`:

| # | File | `pyStartActivity` | Hash | Status |
| ---: | --- | --- | --- | --- |
| 1 | `Flow/InputAddendumFacultativeIn.xml` | **`Start2`** | `f0b2f4ea86` | khas modul ini — §1 |
| 2 | `Flow/OfferFacRetro.xml` | `Start1` | **`b370146c63`** | **VARIAN BERBEDA** — §2 |

`[terverifikasi]` Modul ini **tidak punya** `InputInwardFacultativeRISlip`, `OfferFacOut`,
`InputQuotation`, `InputInwardFacultativeOffer`, maupun `InputRealizationTreatyIn`. Jumlah rule
`Flow` paling sedikit dari ketiga modul facultative (2 vs 6 vs 4).

---

## 1. Diagram alur — `InputAddendumFacultativeIn`

**15 Assignment, 37 Decision, 9 Utility, 3 SubProcess, 133 connector** `[terverifikasi]`.

Titik masuk **`Start2`** (sama dengan RNW Fac In, berbeda dari NB FacIn yang `Start1`):

```
Start2 ─(Always)─> Assignment7 "MARKETING" [WorkBasket, Custom]
                      │ FlowAction: Endorsement_FlowAct
                      v
                   Decision16 "Accept?" (IsUWAccepted)
```

Tangga persetujuan dan routing `LetterNo` **sama bentuknya** dengan NB FacIn §1.2:

```
Decision11 "Limit Akseptasi"
   ├─ When ToKadivTeknik ──────> Assignment15 "KADIV TEKNIK"
   ├─ When ToManagerTeknik ────> Assignment2  "MANAGER TEKNIK"
   ├─ When ToKadivFin ─────────> Assignment3  "KADIV FINANCIAL"
   ├─ When ToDepHeadUW ────────> Assignment8  "DEP HEAD UW FACULTATIVE"
   ├─ When ToDirMarketing ─────> Assignment10 "DIREKTUR MARKETING"
   ├─ When ToKadivFacultative ─> Assignment12 "KADIV FACULTATIVE"
   ├─ When ToDirTeknik ────────> Assignment14 "DIREKTUR TEKNIK"
   └─ Else ────────────────────> Decision39 "Err Konversi?"

Decision18 / Decision38 "Which team?" / "Limit Akseptasi"
   ├─ When ToUW ───────> Assignment9  "UNDERWRITING"
   ├─ When ToSeniorUW ─> Assignment13 "SENIOR UNDERWRITING"
   ├─ When ToJUW_A ────> Assignment16 "JUNIOR UNDERWITER A"
   ├─ When ToDepHeadUW > Assignment8  "DEP HEAD UW FACULTATIVE"
   └─ When IsTBonding ─> Assignment1  "UNDERWRITING FINANCIAL"
```

### 1.1 Percabangan yang **hanya ada di siklus endorsement** `[terverifikasi]`

Enam shape Decision yang tidak punya padanan di NB FacIn maupun RNW Fac In:

```
Decision23 "Is it EDM RISLIP?"
   ├─ When IsEDMRiSlip ─> Assignment4  "TEAMLEADER"
   └─ Else ────────────> Decision30 "Is EDM PPN PPH?"

Decision29 "Is it EDM RISLIP?"
   ├─ When IsEDMRiSlip ─> Decision39 "Err Konversi?"
   └─ Else ────────────> Assignment14 "DIREKTUR TEKNIK"

Decision30 "Is EDM PPN PPH?"
   ├─ When IsEdmPPNPPH ─> Assignment4  "TEAMLEADER"
   └─ Else ────────────> Assignment11 "JUNIOR UNDERWITER B"

Decision31 "Is EDM PPN PPH?"
   ├─ When IsEdmPPNPPH ─> Assignment12 "KADIV FACULTATIVE"
   └─ Else ────────────> Decision39 "Err Konversi?"

Decision28 "Is EDM Internal Retro"
   ├─ When IsEdmInternalRetro ─> Utility7 "SET LIMIT_AKSEPTASI" (GetLimitAkseptasi_JUW_UW)
   └─ Else ────────────────────> Decision18 "Which team?"

Decision19 "IS LIFE?"
   ├─ When IsLife ─> Decision39 "Err Konversi?"
   └─ Else ───────> Decision5  "Is it group?"
```

Jalur konversi dan penyimpanan:

```
Decision39 "Err Konversi?" ├─ When IsFacRetro ─> Decision13 "PRINT ?"
                           │                        ├─ When IsNotPrintRISlip ─> SubProcess3 -> Flow OfferFacRetro
                           │                        └─ Else ──────────────────> End2
                           └─ Else ───────────> Utility3 "SAVE EDM JSON_POLICY" (SaveEDMToJsonPolicy_Act)

Decision32 "Err Konversi?" ├─ When IsSuccessHitService ─> End2
                           └─ Else ───────────────────> Utility6 "SetToInbox_ACT"

Decision3 "FAC OUT?"  ├─ When IsFacRetro ─> Decision4 "INPUT FAC OUT?"
                      │                        ├─ When IsInputFacRetro ─> Utility2 "SET LIMIT_AKSEPTASI"
                      │                        └─ Else ─────────────────> SubProcess2 -> Flow OfferFacRetro
                      └─ Else ───────────> Utility2 "SET LIMIT_AKSEPTASI" (GetLimitAkseptasi_Act)

Decision7 "FAC OUT?"  ├─ When IsFacRetro ─> SubProcess1 -> Flow OfferFacRetro
                      └─ Else ───────────> Decision29 "Is it EDM RISLIP?"
```

**Tahapan siklus endorsement `[terverifikasi]`:** Marketing → tangga akseptasi berbasis limit, dengan
**tiga gerbang tambahan khas endorsement** (EDM RISLIP, EDM PPN PPH, EDM Internal Retro) →
**SAVE EDM JSON_POLICY** (bukan SAVE JSON_POLICY biasa) → fac retro → cek status konversi.

### 1.2 Selisih terukur dari kedua siklus lain

| Aspek | NB FacIn | RNW Fac In | **Endorsment Fac In** |
| --- | --- | --- | --- |
| Rule `Flow` | 6 | 4 | **2** |
| Assignment di flow utama | 21 | 17 | **15** |
| Decision | 42 | 36 | **37** |
| Utility | 12 | 11 | **9** |
| Titik masuk | `Start1` | `Start2` | `Start2` |
| Gerbang khas | — | — | **`IsEDMRiSlip` ×2, `IsEdmPPNPPH` ×2, `IsEdmInternalRetro`** |
| Simpan | `SaveJsonOfferFacIn_Act` + `SaveJsonPolicyFacIn_Act` | `SaveJsonPolicyFacIn_Act` | **`SaveEDMToJsonPolicy_Act`** |
| Shape Arasapas | `serviceInsertArasapas_act` | `serviceInsertArasapasRNW_act` | **`serviceInsertArasapasEDM_act`** |
| Cabang Life | ada (lengkap: Medical, UW Life, DeptHead UW Life) | **tidak ada** | **ada, tetapi hanya `Decision19` "IS LIFE?"** — tanpa workbasket life |
| `GeminiAIGoogle_Act` | ada | ada | **tidak ada** |

Perintah audit: `grep -o "<pyShapeType>[^<]*" "<flow>.xml" | sort | uniq -c` untuk ketiga flow utama.

### 1.3 Shape yatim `[terverifikasi]`

**Tanpa connector masuk:** `Decision11`, `Decision32`, `Decision38`, `Utility1`.
**Tanpa connector keluar:** `Utility1`, `Utility2`, `Utility3`, `Utility5`, `Utility6`, `Utility7`,
`Utility9`, `Utility10`.

`Utility1` "HIT SERVICE ARASAPAS" (→ `serviceInsertArasapasEDM_act`) **tidak punya connector masuk
maupun keluar** — ia terputus sepenuhnya dari graf. → OQ-023.

---

## 2. `OfferFacRetro` — **varian berbeda**, ditelusur terpisah (OQ-011 #338)

`[terverifikasi]` `Endorsment Fac In/Flow/OfferFacRetro.xml` beridentitas sama
(`ASM-FW-GISFW-WORK!OFFERFACRETRO`) dengan salinan NB FacIn/RNW Fac In, tetapi **isinya berbeda**:
hash `b370146c63` vs `bfd6252070`. Varian NB/RNW **tidak dipinjam**; yang ditelusur di sini adalah
file milik modul ini.

**5 Assignment, 6 Decision, 2 Utility** (varian NB/RNW: 3 Assignment, 4 Decision, 2 Utility).

```
Start1 ─(Always)─> Decision1 "IsRISlip"
                      ├─ When IsRISlip ─> Utility2 "Send Email Print RI Slip" (SendEmailPolicy)
                      └─ Else ─────────> Assignment1 "ReasFacOutAdmin"
                                            │ FlowAction: OfferFacOut
                                            v
                                         Decision2 "Accept?" (IsUWAccepted)
                                            ├─ confirm ─> Decision3 "Is it group?"
                                            │                ├─ When IsGroup ─> END52
                                            │                └─ Else ────────> Assignment3 "ReasFacOutHead"
                                            │                                     v  Decision4 "Accept?"
                                            │                                        ├─ confirm ─> Assignment4 "ReasFacOutGroupLeader"
                                            │                                        │                v Decision5 "Accept?"
                                            │                                        │                   ├─ confirm ─> Assignment5 "ReasFacOutTechnicalDirector"
                                            │                                        │                   │                v Decision6 "Accept?"
                                            │                                        │                   │                   ├─ confirm ─> END52
                                            │                                        │                   │                   └─ reject ──> Assignment1
                                            │                                        │                   └─ reject ──> Assignment1
                                            │                                        └─ reject ──> Assignment1
                                            └─ reject ──> END52

Assignment2 "PRINT R/I SLIP" ─Action PrintRISlip_FlowAction─> Utility1 "InsertFacoutProd"
```

### 2.1 Perbedaan perilaku yang terukur `[terverifikasi]`

| | NB FacIn / RNW Fac In (`bfd6252070`) | **Endorsment Fac In** (`b370146c63`) |
| --- | --- | --- |
| Anak tangga persetujuan fac out | **2** | **4** |
| Workbasket | `ReasFacOutAdmin`, `ReasFacOutHead` | `ReasFacOutAdmin`, `ReasFacOutHead`, **`ReasFacOutGroupLeader`**, **`ReasFacOutTechnicalDirector`** |
| Decision `IsUWAccepted` | 2 | **4** |
| Tujuan `reject` | `END52` (dari `Decision2`), `Assignment1` (dari `Decision4`) | **seluruh `reject` di tangga 2–4 kembali ke `Assignment1`** (Admin) |

Perintah audit (selisih graf, bukan selisih teks):
```
awk -f flow.awk "NB FacIn/Flow/OfferFacRetro.xml"          > a.tsv
awk -f flow.awk "Endorsment Fac In/Flow/OfferFacRetro.xml" > b.tsv
diff <(sort -u a.tsv) <(sort -u b.tsv)
```
→ 2 Assignment baru, 2 Decision baru, 6 connector baru, 1 connector berubah tujuan.

`[terverifikasi]` **Ini perbedaan perilaku, bukan perbedaan kosmetik**: persetujuan fac out pada
siklus endorsement melewati **dua tingkat tambahan**. Alasannya **tidak ada di korpus**
(`_METHOD.md` §0 — D2 tidak menyimpulkan *mengapa*). → OQ-011 tetap terbuka untuk entri #338.

---

## 3. Status / state dan kode

### 3.1 Taksonomi tipe endorsement — **terbaca penuh**

`[terverifikasi]` **Empat belas** rule `When` menguji kombinasi
`QuotationData.StatusBusiness = 3` **dan** `QuotationData.EdmType = 4` **dan** nilai
`QuotationData.Type` — menghasilkan **12 nilai `Type` berbeda**:

| `QuotationData.Type` | Rule `When` yang mengujinya |
| ---: | --- |
| `0` | `IsEdmAdjRefNo`, `IsEDMRiSlip` (via `EdmTypeNew = 4`) |
| `1` | `IsEdmExtendPeriod` |
| `2` | `IsEdmAdjTSI` |
| `3` | `IsEdmAdjRate` |
| `4` | `IsEdmAdjSpreading` |
| `5` | `IsEdmAddObject` |
| `6` | `IsEdmAdjPeriod` |
| `7` | `IsEdmAdjCurrency` |
| `8` | `IsEdmAdjRIC` |
| `9` | `IsEdmAdjInsured` |
| `11` | `IsEdmAdjShareCedant` |
| `12` | `IsEdmAdjCeding`, `IsEdmAdjRIC`, `IsEdmPPNPPH` |

Perintah audit:
```
cd "Endorsment Fac In/When"
for f in IsEdm*.xml IsEDM*.xml; do
  t=$(grep -oE '<pyLabel>[^<]{1,160}' "$f" | sed 's/<pyLabel>//;s/&amp;#61;/=/g;s/\]\[/ /g;s/[][]//g' \
      | grep -oE "QuotationData\.Type = [0-9]+" | sed 's/.*= //' | sort -u | tr '\n' ',')
  printf "%-22s Type=%s\n" "${f%.xml}" "$t"
done
```

`[terverifikasi]` Nilai **`10` tidak dipakai**; nilai **`12` dipakai tiga rule berbeda**.
**Arti setiap nilai belum terverifikasi** (OQ-020) — nama rule bukan bukti (`_METHOD.md` §2.2).

`[terverifikasi]` `When/IsEdmPerubahan.xml` adalah **komposit**: kondisinya berupa rujukan
"Rule[X] evaluates to true" untuk `IsEdmAddObject`, `IsEdmAdjCurrency`, `IsEdmAdjPeriod`,
`IsEdmAdjRate`, `IsEdmAdjRIC`, `IsEdmAdjSpreading`, `IsEdmAdjTSI` — tujuh dari dua belas tipe.

`[terverifikasi]` `When/IsEdmInternalRetro.xml` — `<pyLabel>` hanya template kosong
`[first value][relation][second value]`. Kondisinya **tidak terbaca**, dan ia menggerbangi
`Decision28` (§1.1) → OQ-029.

`[terverifikasi]` Nilai `EdmType` yang muncul di modul ini: **1, 2, 3, 4**.
```
grep -rhoE "EdmType(New)?[^<]{0,18}" "Endorsment Fac In" --include="*.xml" \
| sed 's/&amp;#61;/=/g;s/&amp;quot;/"/g;s/\]\[/ /g;s/[][]//g;s/ //g' | grep -E "=" | sort | uniq -c | sort -rn
```
→ `EdmType=4` 44×, `EdmType==2` 33×, `EdmType==4` 19×, `EdmType==1` 12×, `EdmType=="2"` 8×,
`EdmType==3` 5×.

Ini melanjutkan bukti `EdmType` yang ditemukan di `Endorsement Life.md` (Tahap 1) → OQ-014.

### 3.2 `StatusBusiness = 3` — siklus ini

`[terverifikasi]` `When/IsEDM.xml`: `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3`;
`When/IsNotEDM.xml`: `!= 3`. Rule `When/IsNB.xml` (`= 1`) dan `When/IsRenewal.xml` (`= 2`) **juga ada
di modul ini** — lihat `RNW Fac In.md` §2.1 untuk bukti lengkap diskriminator siklus (OQ-015).

### 3.3 `ProposalAcceptStatus`

`[terverifikasi]` Sama dengan NB FacIn (`NB FacIn.md` §2.1), **dengan dua selisih terukur**:

1. Modul ini punya **6** Data Transform ber-`ProposalAcceptStatus`, bukan 8 —
   **`SetAkseptasiCeding_DT` dan `SetReviseProposal` tidak ada**.
2. Hash `DecisionTable/IsUWAccepted.xml` = **`9226a7db77`**, berbeda dari NB FacIn/RNW Fac In
   (`f99bc43c45`). `When/IsUWAccepted.xml` justru **identik** (`fa24354fb9` di NB FacIn dan di sini).

Perintah audit:
```
for m in "NB FacIn" "RNW Fac In" "Endorsment Fac In"; do
  grep -rl "ProposalAcceptStatus" "$m/DataTransform" --include="*.xml" | xargs -n1 basename | sort
done
for m in "NB FacIn" "RNW Fac In" "Endorsment Fac In"; do sh nhash.sh "$m/DecisionTable/IsUWAccepted.xml"; done
```

`[terverifikasi]` Tabel keputusan yang menentukan hasil persetujuan **berbeda isinya di siklus
endorsement** — tetapi karena baris tabelnya tidak ada di ekspor (OQ-043), **apa yang berbeda tidak
dapat dinyatakan**. Ini gabungan terburuk OQ-011 × OQ-043.

### 3.4 Rule berkonflik lain, dan **pengukuran ulang register OQ-011**

`[terverifikasi]` **438 dari 533 identitas berkonflik korpus menyentuh ketiga modul facultative
sekaligus** (82,2%). Perintah audit:
```
awk -F'|' '/^\| [0-9]+ \|/ { if ($0 ~ /NB FacIn\// && $0 ~ /RNW Fac In\// && $0 ~ /Endorsment Fac In\//) n++ }
 END{print n}' inventory/_oq011-konflik-isi.md
```

Sebarannya per tipe rule: `RULE-OBJ-ACTIVITY` 171, `RULE-HTML-SECTION` 78, `RULE-OBJ-FLOWACTION` 42,
`RULE-CONNECT-SQL` 29, `RULE-OBJ-MODEL` 27, `RULE-DECLARE-PAGES` 27, `RULE-OBJ-WHEN` 24,
`RULE-OBJ-REPORT-DEFINITION` 21, `RULE-HTML-HARNESS` 8, `RULE-DECLARE-DECISIONTABLE` 6,
`RULE-CONNECT-REST` 3, `RULE-OBJ-FLOW` 1, `RULE-ADMIN-SYSTEM-SETTINGS` 1.

**Pengukuran baru `[terverifikasi]`.** §5.1 menemukan bahwa varian
`ASM-FW-GISFW-WORK!SERVICEINSERTARASAPAS_ACT` hanya berbeda pada **tiga tag metadata Pega**
(`<pyDelete>`, `<pyVersionSecure>`, `<pzIsPrivateCheckOut>`) yang **tidak** ada di daftar 18 tag
normalisasi. Seluruh 438 identitas di atas dihitung ulang dengan normalisasi **21 tag**:

| Ukuran | Hasil |
| --- | ---: |
| Identitas berkonflik menyentuh ketiga modul | **438** |
| Masih >1 isi dengan normalisasi 18 tag | **438** |
| Masih >1 isi dengan normalisasi **21 tag** | **429** |
| **Runtuh menjadi identik** oleh 3 tag baru | **9** (2,1%) |

Kesembilan konflik semu itu: **#12** `ASM!GETDATAKLAIM_SQL`, **#13** `ASM!GETEDMOLDIDPEGA`,
**#324** `LINKSERVICE!LINKSERVICE`, **#367** `RNM!DELETESTORAGE_SQL`, **#369** `RNM!GETAPPNAME_SQL`,
**#371** `RNM!GETLINKSTORAGE_SQL`, **#372** `RNM!GETMKTANDLEADER_SQL`,
**#376** `RNM!INSERT_T_STORAGE_SQL`, **#377** `RNM!UPDATE_T_STORAGE_SQL`.

Perintah audit (dapat dijalankan ulang penuh):
```
VOL21='<(…18 tag…|pyDelete|pyVersionSecure|pzIsPrivateCheckOut)[/>]'
# untuk tiap path di baris konflik: grep -vE "$VOL21" "$p" | sort | md5sum
# lalu hitung berapa identitas yang hash-nya tinggal satu
```

`[terverifikasi]` **Register OQ-011 pada dasarnya tetap berlaku**: 429 dari 438 (97,9%) tetap
berbeda isi. Yang berubah adalah **#324 `LINKSERVICE`** — rule SystemSettings yang memuat base URL
ternyata **identik** di ketiga modul, sehingga konfigurasi endpoint tidak bercabang antar siklus.
Ini relevan langsung ke OQ-018/OQ-047.

---

## 4. Objek Oracle yang disentuh

```
awk -F'\t' '$1 ~ "^Endorsment Fac In/" && $2=="RULE-CONNECT-SQL"{print $8}' all-rules.tsv \
  | grep -oE "tables=[^;]*" | sed 's/tables=//' | tr ',' '\n' | sort | uniq -c | sort -rn
```

| Objek | Rule perujuk |
| --- | ---: |
| **`JSON_POLIS`** | 17 |
| `FACINPRODUCTION` | 6 |
| `RW`, `POOLDATA.REINSURANCETYPE`, `OCCUPATION`, `CURRENCY` | 5 masing-masing |
| `TREATYBUSINESS`, `PROPORTIONALARRG`, `POOLDATA.M_LIMIT_PROPERTYY`, `POOLDATA.JSON_POLIS`, `MARKETINGOFFICER`, **`FACOUTPRODUCTION`**, `ACCUMULATION` | 4 masing-masing |
| `T_STORAGE_IMAGE`, `POOLDATA.M_LIMIT_NONPROPANDENGG` | 3 masing-masing |

`[terverifikasi]` **`FACOUTPRODUCTION` (4 rule) hanya menonjol di modul ini** — di NB FacIn dan
RNW Fac In ia tidak masuk daftar teratas. Sejalan dengan §2: siklus endorsement punya tangga fac out
yang lebih panjang.

---

## 5. Integrasi eksternal

`[terverifikasi]` Tiga rule `RULE-CONNECT-REST` yang sama dengan dua modul lain
(`ServiceGoogle`, `convertJsonNusareToProduction`, `getPremiumPaidOn`), seluruhnya
`pyBaseURLSelectionType = SETTING`, `pyBaseURLSetting = LinkService!LinkService`, **tanpa URL literal
endpoint**.

`[terverifikasi]` **Tidak ada `GeminiAIGoogle_Act` di modul ini** — berbeda dari NB FacIn dan
RNW Fac In (OQ-048). Tiga activity Google Storage lainnya ada.

### 5.1 `serviceInsertArasapas_act` — di sini varian miliknya **ADA**, jadi telusur TIDAK terhenti

`[terverifikasi]` Berbeda dari NB Treaty In (OQ-025), NB FacIn, dan RNW Fac In, modul ini
**memiliki** varian implementasinya sendiri:

`Endorsment Fac In/Activity/serviceInsertArasapas_act.xml`
→ `ASM-FW-GISFW-WORK!SERVICEINSERTARASAPAS_ACT`, **231.975 byte**, hash `7314b6c49c`.

Sesuai `_METHOD.md` §1.1: varian **milik modul konteks** ditelusur; varian `EDM Treaty In`
(`ab3ae3c9ba`) **tidak dipinjam** untuk menjelaskan modul lain.

**Dua puluh langkahnya `[terverifikasi]`:**

| # | Method | Catatan |
| ---: | --- | --- |
| 1 | `Obj-Refresh-And-Lock` | mengunci case |
| 2 | `Property-Set` | |
| 3 | `RDB-List` | |
| 4 | `Property-Set` | |
| 5 | **`Call ASM-FW-GISFW-Int-M_LINK_SERVICE.GetLinkService`** | resolusi endpoint dari tabel DB (OQ-047) |
| 6 | **`Connect-REST`** | panggilan keluar #1 |
| 7 | **`Connect-REST`** | panggilan keluar #2 |
| 8 | `Property-Set` | |
| 9 | `Obj-Save` | |
| 10 | `Page-Set-Messages` | jalur galat |
| 11–12 | `RDB-List` ×2 | |
| 13–14 | `Property-Set` ×2 | |
| 15–16 | `RDB-List` ×2 | |
| 17 | `Page-Remove` | |
| 18 | `Call ASMForceCaseClose` | menutup case |
| 19 | `Property-Set` | |
| 20 | `Call SendEmailWithAttachments` | |

Rule yang dirujuk `[terverifikasi]`:
`ASM-FW-GISFW-Int-policyjson ASM GetPolicyNoByCaseId`,
`ASM-FW-GISFW-Int-policyjson ASM INSERTJSON_JSONPOLISMONITORING_FACIN`,
`ASM-FW-GISFW-Int-policyjson ASM UpdateErrorNoteJsonPolisMonitoring`,
**`ASM-FW-GISFW-Int-policyjson RNM CekSTSKonversiJson`**,
**`ASM-FW-GISFW-Int-policyjson RNM DeleteDataProduction`**,
serta fungsi kustom **`@ASM.GetPageJSONString()`**.

`[terverifikasi]` Dua rule memakai prefix **`RNM`**, sisanya `ASM`, di class yang sama — bukti
konteks baru untuk OQ-008 (arti prefix `ASM` vs `RNM` pada kunci RDBList).

`[terverifikasi]` `@ASM.GetPageJSONString()` — fungsi Pega kustom yang **source-nya tidak ada di
korpus** — dipakai lagi di sini untuk membentuk JSON yang dikirim keluar. Sama seperti temuan
`NB Treaty In.md` §5.3 → memperkuat OQ-012.

### 5.2 Kedua varian ternyata berbeda **hanya pada metadata** — temuan baru

`[terverifikasi]` Perbandingan baris demi baris setelah normalisasi 18 tag antara
`EDM Treaty In/Activity/serviceInsertArasapas_act.xml` (`ab3ae3c9ba`) dan
`Endorsment Fac In/Activity/serviceInsertArasapas_act.xml` (`7314b6c49c`) menghasilkan
**tepat 3 baris berbeda**:

```
<pyDelete>            false   vs  true
<pyVersionSecure>     true    vs  false
<pzIsPrivateCheckOut> true    vs  (kosong)
```

Perintah audit:
```
VOL='<(…18 tag…)>'
diff <(grep -vE "$VOL" "EDM Treaty In/Activity/serviceInsertArasapas_act.xml" | sort) \
     <(grep -vE "$VOL" "Endorsment Fac In/Activity/serviceInsertArasapas_act.xml" | sort)
diff <(grep -oE "<pyStepsActivityName>[^<]*" "EDM Treaty In/Activity/serviceInsertArasapas_act.xml") \
     <(grep -oE "<pyStepsActivityName>[^<]*" "Endorsment Fac In/Activity/serviceInsertArasapas_act.xml")
```
→ 3 baris berbeda; **urutan 20 langkah identik**; **daftar rule yang dirujuk identik**.

`[terverifikasi]` Ketiganya adalah **metadata pengelolaan rule Pega** (status hapus, keamanan versi,
status check-out privat) — **bukan perilaku**. Artinya konflik OQ-011 entri **#415 adalah artefak
check-out/versi, bukan percabangan perilaku**.

**Konsekuensi untuk OQ-025:** blocker di `NB Treaty In`, `NB FacIn`, dan `RNW Fac In` tetap berdiri
secara prosedural — ketiga modul itu **tidak memuat** implementasinya, jadi apa yang dieksekusi di
lingkungan mereka tetap tidak dapat dipastikan dari korpus. Tetapi **kedua varian yang ada di korpus
berperilaku sama**, sehingga risiko salah tafsirnya jauh lebih kecil dari yang diperkirakan D1.
→ OQ-025 diperbarui, **belum ditutup** (keputusan versi mana yang aktif tetap milik pemilik export).

---

## 6. Batas pengetahuan

### 6.1 Mesin data lama (*before-image*) — class yang belum pernah muncul di D2

`[terverifikasi]` Delapan activity penyalin nilai lama, di **dua class endorsement yang berbeda**:

| Activity | Ukuran | Langkah | Class |
| --- | ---: | ---: | --- |
| `SetValueToEDMWork.xml` | 445.703 | 41 | **`ASM-SFAGIS-WORK-ENDORSEMENT`** |
| `SetErrorBatalEndorsement_Act.xml` | 425.872 | 39 | `ASM-SFAGIS-WORK-ENDORSEMENT` |
| `CountEndorsementData.xml` | 357.791 | 25 | `ASM-SFAGIS-WORK-ENDORSEMENT` |
| `SetOLDValueToEDMWork_FIRE.xml` | 170.097 | 13 | `ASM-SFAGIS-WORK-ENDORSEMENT` |
| `SetOLDValueToEDMWork_LIFE.xml` | 74.841 | 4 | `ASM-SFAGIS-WORK-ENDORSEMENT` |
| `SetEdmType.xml` | 73.041 | 4 | `ASM-SFAGIS-WORK-ENDORSEMENT` |
| `SetOLDValueToEDMWork_{Aneka,GOLF,MBU,MC,PA}.xml` | — | — | `ASM-SFAGIS-WORK-ENDORSEMENT` |
| `EDMRetro_Act.xml` | 99.051 | 12 | **`ASM-FW-GISFW-WORK-ENDORSEMENT`** |

**Isinya belum dibaca.** Yang dicatat: ada **tujuh varian per lini bisnis**
(`_Aneka`, `_FIRE`, `_GOLF`, `_LIFE`, `_MBU`, `_MC`, `_PA`) dari satu pola penyalinan nilai lama —
seluruhnya **eksklusif modul ini** (tidak ada di NB FacIn maupun RNW Fac In).

`[terverifikasi]` **Prefix class `ASM-SFAGIS-` baru bagi D2.** Sebarannya di korpus:
```
awk -F'\t' 'toupper($3) ~ /SFAGIS/ {print $3}' all-rules.tsv | sort | uniq -c | sort -rn
```
→ `ASM-SFAGIS-WORK-ENDORSEMENT` 31, `ASM-FW-SFAGISFW-WORK-OPPORTUNITY` 5,
`ASM-FW-SFAGISFW-WORK-ACCOUNT` 5.

Empat class bersufiks `ENDORSEMENT` hidup berdampingan di korpus:
`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` (38), `ASM-SFAGIS-WORK-ENDORSEMENT` (31),
`ASM-FW-GISFW-WORK-ENDORSEMENTTREATY` (10), `ASM-FW-GISFW-WORK-ENDORSEMENT` (5).
→ **OQ-049**.

`[terverifikasi]` Rule yang menyebut `OldData` di modul ini: **133** — terbanyak dari ketiga modul
(NB FacIn 129, RNW Fac In 114). `Activity/SetOldData.xml` sendiri **identik** di ketiganya
(`df291257aa`, 450.829 byte, 59 langkah) — lihat `RNW Fac In.md` §5.5.

### 6.2 Prorate dan selisih

`[terverifikasi]` Pemakaian `ProRateType` di modul ini **paling padat** dari ketiganya:
`ProRateType==3||…` 8×, `=="3"||…` 8×, `==4)` 6×, `==3,.Premium-.PremiumRetro` 6×,
`=="4")&&primary.` 6×, `==4` 5×, `!=4&&local.count` 5× — **23 kemunculan berkondisi**
(NB FacIn 12, RNW Fac In 9).

Activity prorate yang **eksklusif modul ini**: `SetLocalNProrate` (76.187 B, 5 langkah,
class `ASM-FW-GISFW-DATA-COVERAGE`), `SetLocalNonMbuProrate` (297.476 B, 26 langkah, class sama),
`ChangeProRatePreActPA_FacIn`. **Isinya belum dibaca.**

`[terverifikasi]` Procedure Oracle **`FIRE.CEK_PRORATA_TANGGAL`** dipanggil dari
`Endorsment Fac In/RDBList/SearchSQLRateKPR.xml` dan `SearchSQLRateNonKPR.xml`. Ia
**tidak dipanggil dari NB FacIn maupun RNW Fac In** `[terverifikasi]`. Isinya tidak ada di korpus
(OQ-002) → **perhitungan prorata tanggal berada di database**.

Dua procedure lain yang juga **eksklusif modul ini**: `FIRE.PEGA_FIRE_SET_RATE`
(`SearchTemplateMainCoverageSQL.xml`) dan `NEW_GENERAL.CEK_PLAT_NO`. Skema `NEW_GENERAL` **baru**
→ OQ-016.

Perintah audit:
```
awk -F'\t' '$1 ~ "^Endorsment Fac In/" && $8 ~ /CEK_PRORATA_TANGGAL|PEGA_FIRE_SET_RATE|CEK_PLAT_NO/ {print $5" <- "$1}' all-rules.tsv
```

**Istilah "menjadi/selisih"** (bentuk lazim penyajian endorsement) `[terverifikasi]`: tidak ditemukan
sebagai nama properti, nama rule, maupun label di korpus modul ini. Yang ada adalah pasangan
**nilai lama** (`SetOLDValueToEDMWork_*`, `OldData`, `getOldDeduct`, `ViewOldDataEDM_Act`) dan
**nilai baru** (`SetValueToEDMWork`, `EdmTypeNew`), plus activity penghitung selisih bernama
`CountEdmAdjTSI_Act` (169.196 B, 19 langkah) dan `CountEdmExtPeriode_Act` (228.536 B, 24 langkah).
**Isinya belum dibaca**, sehingga **bentuk penyajian selisih tidak dinyatakan**.

### 6.3 Stored procedure — 28 dipanggil

`[terverifikasi]` Sama dengan `NB FacIn.md` §5.2 **dikurangi** `POOLDATA.PEGA_JSON_POLIS_TREATYIN`,
`POOLDATA.PEGA_TREATY_IN`, `POOLDATA.PEGA_M_JSON_OFFER`, `POOLDATA.RDBMASTERSHIP`,
`UTL_MATCH.EDIT_DISTANCE_SIMILARITY` — **ditambah tiga yang eksklusif**:
`FIRE.CEK_PRORATA_TANGGAL`, `FIRE.PEGA_FIRE_SET_RATE`, `NEW_GENERAL.CEK_PLAT_NO`.

Isinya tidak ada di korpus (OQ-002).

### 6.4 Database link — delapan objek

`[terverifikasi]` Delapan objek lewat `@ASMD.SINARMAS.CO.ID`, **termasuk `LST_KURS_STANDARD`**
(yang tidak dipakai RNW Fac In). Sama dengan NB FacIn → OQ-017.

### 6.5 Spreading — 38 activity, tak satu pun habis dibaca

`[terverifikasi]` `ls -l "Endorsment Fac In/Activity/" | grep -icE "spread|capacit|scoring"` → **38**.

Selisih dari NB FacIn: **`CekLimitSpreading_Act` tidak ada** di modul ini (padahal ia yang dipanggil
`Utility5` "Check Spreading" di `InputInwardFacultativeRISlip`); sebagai gantinya ada
**`CopySpreading_Act` (627.836 B)**, `SetSpreadingCoverage_ACT`, `SetSumTSISpreading_Act`,
`SpreadingProtection` — **eksklusif modul ini**.

Terbesar tetap `SumTSIPremiSpreadedRNM_FIRE_Act.xml` **995.970 byte** — **belum dibaca**
(ukurannya berbeda 82 byte dari salinan NB FacIn; keduanya terdaftar berkonflik di OQ-011).

Rumus spreading/capacity/scoring **tidak dinyatakan** — batas cakupan telusur.

### 6.6 `CountPremiNusareRetro_Act`

`[terverifikasi]` **Identik** dengan NB FacIn (`8f1268a660`). Rumus pangsa 0,45/0,05 dan ambang
3 miliar berlaku sama di siklus endorsement — lihat `NB FacIn.md` §5.5 → OQ-046.

---

## 7. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `Endorsment Fac In/`) | OQ-011? |
| --- | --- | --- |
| `ASM-FW-GISFW-WORK` / `INPUTADDENDUMFACULTATIVEIN` / `RULE-OBJ-FLOW` | `Flow/InputAddendumFacultativeIn.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `OFFERFACRETRO` / `RULE-OBJ-FLOW` | `Flow/OfferFacRetro.xml` | **YA — #338**, **varian modul ini** (§2) |
| `ASM-FW-GISFW-WORK` / `SERVICEINSERTARASAPAS_ACT` / `RULE-OBJ-ACTIVITY` (231.975 B) | `Activity/serviceInsertArasapas_act.xml` | **YA — #415**, varian modul ini (§5.1–5.2) |
| `ASM-FW-GISFW-WORK` / `SERVICEINSERTARASAPASEDM_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/serviceInsertArasapasEDM_act.xml` | **YA — #414** |
| `ASM-FW-GISFW-WORK` / `ISEDM`, `ISNOTEDM`, `ISNB`, `ISRENEWAL` / `RULE-OBJ-WHEN` | `When/` | tidak — diskriminator siklus |
| `ASM-FW-GISFW-WORK` / `ISEDMRISLIP`, `ISEDMPPNPPH`, `ISEDMINTERNALRETRO`, `ISEDMPERUBAHAN` / `RULE-OBJ-WHEN` | `When/` | sebagian |
| `ASM-FW-GISFW-WORK` / 10 `ISEDMADJ*` + `ISEDMADDOBJECT`, `ISEDMEXTENDPERIOD` / `RULE-OBJ-WHEN` | `When/` | sebagian — **taksonomi §3.1** |
| `ASM-FW-GISFW-WORK` / `ISLIFE`, `ISFACRETRO`, `ISINPUTFACRETRO`, `ISNOTPRINTRISLIP`, `ISSUCCESSHITSERVICE`, `ISGROUP`, `ISGROUPCREATE`, `ISTBONDING`, `ISPKSASM` / `RULE-OBJ-WHEN` | `When/` | sebagian |
| `ASM-FW-GISFW-WORK` / 9 `TO*` (routing `LetterNo`) / `RULE-OBJ-WHEN` | `When/` | sebagian |
| `ASM-SFAGIS-WORK-ENDORSEMENT` / `SETVALUETOEDMWORK` (445.703 B), `SETERRORBATALENDORSEMENT_ACT` (425.872 B), `COUNTENDORSEMENTDATA` (357.791 B), `SETEDMTYPE`, 7× `SETOLDVALUETOEDMWORK_*` / `RULE-OBJ-ACTIVITY` | `Activity/` | sebagian |
| `ASM-FW-GISFW-WORK-ENDORSEMENT` / `EDMRETRO_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/EDMRetro_Act.xml` | tidak |
| `ASM-FW-GISFW-DATA-OFFERFACIN` / `COUNTEDMEXTPERIODE_ACT`, `SETADJNOREFFEDM_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/` | sebagian |
| `ASM-FW-GISFW-DATA-ANEKA` / `COUNTEDMADJTSI_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/CountEdmAdjTSI_Act.xml` | sebagian |
| `ASM-FW-GISFW-DATA-COVERAGE` / `SETLOCALNPRORATE`, `SETLOCALNONMBUPRORATE` / `RULE-OBJ-ACTIVITY` | `Activity/` | sebagian |
| `ASM-FW-GISFW-WORK` / `SAVEEDMTOJSONPOLICY_ACT`, `GETLIMITAKSEPTASI_ACT`, `GETLIMITAKSEPTASI_JUW_UW`, `SETBANDING_ACT`, `SAVETOPRODUCTION_ACT`, `SETOLDDATA`, `COUNTPREMINUSARERETRO_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/` | sebagian |
| `ASM-FW-GISFW-INT-M_LINK_SERVICE` / `GETLINKSERVICE` / `RULE-OBJ-ACTIVITY` | `Activity/GetLinkService.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `ISUWACCEPTED` / `RULE-DECLARE-DECISIONTABLE` (hash `9226a7db77`) | `DecisionTable/IsUWAccepted.xml` | tidak (isi baris tak ada — OQ-043) |
| `ASM-FW-GISFW-WORK` / `ISUWACCEPTED` / `RULE-OBJ-WHEN` (hash `fa24354fb9`) | `When/IsUWAccepted.xml` | tidak |
| `ASM-FW-GISFW-WORK` / 6 Data Transform `SET*PROPOSAL*` / `RULE-OBJ-MODEL` | `DataTransform/` | sebagian |
| `ASM-FW-GISFW-INT-POLICYJSON` / `SEARCHSQLRATEKPR`, `SEARCHSQLRATENONKPR`, `SEARCHTEMPLATEMAINCOVERAGESQL` / `RULE-CONNECT-SQL` | `RDBList/` | sebagian |
| (FlowAction) `ENDORSEMENT_FLOWACT`, `ENDORSEMENT_FLOWACT_ISUW`, `OFFERFACOUT`, `PRINTRISLIP_FLOWACTION` | `FlowAction/` | sebagian |
| 3 rule `RULE-CONNECT-REST` (§5) | `ConnectREST/` | sebagian |

**52 rule ditelusur** (2 Flow; 33 When; 11 Activity; 1 DecisionTable; 3 RDBList; 3 ConnectREST —
dengan tumpang tindih antar baris tabel).

Rule **dirujuk tapi belum ditelusur**: 38 activity spreading/scoring (§6.5), 8 activity data-lama
(§6.1), 3 activity prorate (§6.2), `SaveEDMToJsonPolicy_Act`, `SetToInbox_ACT`, `ASMForceCaseClose`,
`InsertFacoutProd`, `SendEmailWithAttachments`, dan **94** activity yang eksklusif modul ini terhadap
NB FacIn (`comm -23 <(ls "Endorsment Fac In/Activity"|sort) <(ls "NB FacIn/Activity"|sort) | wc -l`).

---

## 8. Pertanyaan terbuka

| OQ | Ringkas | Pemilik |
| --- | --- | --- |
| **OQ-049** | Empat class bersufiks `ENDORSEMENT` di korpus, salah satunya berprefix **`ASM-SFAGIS-`** (bukan `ASM-FW-GISFW-`) — apakah aplikasi/ruleset berbeda? | Arsitektur Pega + pemilik export |

OQ diperbarui secara material:
- **OQ-025** — §5.2: kedua varian implementasi berbeda **hanya pada 3 tag metadata Pega**; konflik
  #415 adalah artefak check-out, bukan percabangan perilaku. Blocker prosedural tetap berdiri untuk
  3 modul yang tidak memuat implementasinya.
- **OQ-011** — §3.4: pengukuran ulang 438 konflik facultative dengan normalisasi 21 tag →
  **429 tetap berbeda (97,9%)**, 9 konflik semu, termasuk **`LINKSERVICE`**.
- **OQ-018 / OQ-047** — §3.4: `LINKSERVICE!LINKSERVICE` identik di ketiga modul.

OQ dikuatkan: OQ-002 (§6.3), OQ-008 (§5.1 — prefix `RNM` vs `ASM` dalam satu class),
OQ-012 (§4, §5.1 — `@ASM.GetPageJSONString()`), OQ-014 (§3.1 — `EdmType` 1–4),
OQ-015 (§3.2), OQ-016 (§6.2 — skema `NEW_GENERAL`), OQ-017 (§6.4), OQ-020 (§3.1),
OQ-021/OQ-027 (guard identitas identik dgn dua modul lain), OQ-023 (§1.3), OQ-024,
OQ-028 (§1 — seluruhnya WorkBasket), OQ-029 (§3.1 — `IsEdmInternalRetro`),
OQ-043 (§3.3 — tabel keputusan berbeda **dan** isinya tak ada), OQ-045, OQ-046 (§6.6),
OQ-047 (§5.1), OQ-048 (§5 — modul ini **tidak** punya `GeminiAIGoogle_Act`).

**Guard identitas** `[terverifikasi]`: `When/IsGroup.xml` (3 identitas), `When/IsGroupCreate.xml`
(2 identitas), `When/IsTBonding.xml` (1) — **hash identik dengan NB FacIn dan RNW Fac In**.
Nilai nama orang **tidak disalin**.
