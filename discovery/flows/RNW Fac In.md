# Telusur Flow — RNW Fac In

STEP D2, Tahap 4 konteks #14. Ditelusur 2026-09-13. Konvensi: `_METHOD.md`.
**1.927 file** (`find "RNW Fac In" -name '*.xml' | wc -l`).

Artefak ini **mengutamakan yang BEDA dari new business**. Untuk perilaku yang sama, rujuk
**`NB FacIn.md`** — tidak diulang di sini.

**Empat rule `Flow`**, seluruhnya class `ASM-FW-GISFW-WORK`:

| # | File | `pyStartActivity` | Hash | Status telusur |
| ---: | --- | --- | --- | --- |
| 1 | `Flow/InputRenewalFacultativeIn.xml` | **`Start2`** | `f3dff3940e` | **khas modul ini — ditelusur di §1** |
| 2 | `Flow/InputInwardFacultativeRISlip.xml` | `Start1` | `c099bf4ebb` | **identik** dgn NB FacIn → `NB FacIn.md` §1.3 |
| 3 | `Flow/OfferFacOut.xml` | `Start62` | `832fb12b9d` | **identik** dgn NB FacIn → `NB FacIn.md` §1.4 |
| 4 | `Flow/OfferFacRetro.xml` | `Start1` | `bfd6252070` | **identik** dgn NB FacIn → `NB FacIn.md` §1.5 |

`[terverifikasi]` Ketiga flow #2–#4 punya hash ternormalisasi **sama persis** dengan salinannya di
`NB FacIn`. **Tidak ditelusur ulang.** Perintah audit:

```
sh nhash.sh "RNW Fac In/Flow/InputInwardFacultativeRISlip.xml" "NB FacIn/Flow/InputInwardFacultativeRISlip.xml"
sh nhash.sh "RNW Fac In/Flow/OfferFacOut.xml"                  "NB FacIn/Flow/OfferFacOut.xml"
sh nhash.sh "RNW Fac In/Flow/OfferFacRetro.xml"                "NB FacIn/Flow/OfferFacRetro.xml"
```

`[terverifikasi]` Modul ini **tidak punya** `InputInwardFacultativeOffer`, `InputQuotation`, maupun
`InputRealizationTreatyIn` — tiga flow yang ada di `NB FacIn`.

---

## 1. Diagram alur — `InputRenewalFacultativeIn`

**17 Assignment, 36 Decision, 11 Utility, 5 SubProcess, 140 connector** `[terverifikasi]`
(`grep -o "<pyShapeType>[^<]*" | sort | uniq -c`).

Titik masuk **`Start2`**, bukan `Start1` seperti flow offer di NB FacIn `[terverifikasi]`.

```
Start2 ─(Always)─> Assignment10 "MARKETING" [WorkBasket, Custom]
```

Tulang punggungnya **sama bentuknya** dengan `InputInwardFacultativeOffer` (NB FacIn §1.2):
tangga persetujuan 17 workbasket, gerbang `IsUWAccepted` dengan enam hasil, routing lewat token
`LetterNo`, jalur banding, dan percabangan fac out/retro. Yang dicatat di bawah hanya **selisihnya**.

### 1.1 Yang BEDA dari siklus new business

| Aspek | NB FacIn (`InputInwardFacultativeOffer`) | RNW Fac In (`InputRenewalFacultativeIn`) |
| --- | --- | --- |
| Titik masuk | `Start1` → `Assignment12` MARKETING | **`Start2`** → `Assignment10` MARKETING |
| Assignment | 21 | **17** |
| Decision | 42 | **36** |
| Utility | 12 | 11 |
| SubProcess | 3 (semuanya ke `OfferFacRetro`/`RISlip`) | **5** |
| Cabang Life | **ada** — `Decision36`/`Decision37` "IS LIFE?", `Utility5` "SET LIMIT AKSEPTASI LIFE", Assignment `MEDICAL LIFE`, `UNDERWRITING LIFE`, `DEPTHEAD UNDERWRITING LIFE` | **tidak ada** — tak satu pun shape "IS LIFE?" atau workbasket life |
| Activity limit | `GetLimitAkseptasi_ActFlow` (539.687 B) + `GetLimitAkseptasiLife_Act` | **`GetLimitAkseptasi_Act`** + `GetLimitAkseptasi_JUW_UW`; **tanpa varian Life** |
| Shape Arasapas | `Utility2` → `serviceInsertArasapas_act` | `Utility1` → **`serviceInsertArasapasRNW_act`** (§4.1) |
| Simpan JSON | `Utility9`/`Utility16` "SAVE JSON_OFFER" → `SaveJsonOfferFacIn_Act` | `Utility3` "SAVE JSON_POLICY" → `SaveJsonPolicyFacIn_Act`; **tidak ada shape SAVE JSON_OFFER** |
| SubProcess ke RI Slip | `SubProcess1` "POLICY" | `SubProcess1` "POLICY" (sama) |
| `ISPKSASM` | 2 Decision | 2 Decision (`Decision24`, `Decision33`) |

Perintah audit selisih shape:
```
awk -F'\t' '$2!~/Connector/{print $3}' t4-rnw.tsv | sort -u
awk -F'\t' '$2!~/Connector/{print $3}' t4-nbfacin-offer.tsv | sort -u
```

`[terverifikasi]` **Tidak ada jalur Life di siklus renewal.** Ini selisih perilaku paling besar:
new business punya cabang medis/underwriting life lengkap; renewal tidak punya sama sekali.

### 1.2 Jalur fac out / retro — identik bentuknya

```
Decision1  "FAC OUT?"       ├─ When IsFacRetro ─> Decision19 "INPUT FAC OUT?"
                            │                        ├─ When IsInputFacRetro ─> Utility2 "SET LIMIT_AKSEPTASI"
                            │                        └─ Else ─────────────────> SubProcess8 "[Fac Out]" -> OfferFacRetro
                            └─ Else ───────────> Utility2 "SET LIMIT_AKSEPTASI" (GetLimitAkseptasi_Act)

Decision5  "FAC OUT?"       ├─ When IsFacRetro ─> SubProcess4 "[Fac Out]" -> OfferFacRetro
                            └─ Else ───────────> Assignment6 "DIREKTUR TEKNIK"

Decision13 "FAC OUT?"       ├─ When IsFacRetro ─> Decision2 "PRINT ?"
                            │                        ├─ When IsNotPrintRISlip ─> Utility4 "UPDATE STS KONVERSI"
                            │                        └─ Else ──────────────────> End2
                            └─ Else ───────────> End2

Decision25 "It is Group?"   ├─ When IsGroup ─> Utility3 "SAVE JSON_POLICY"
                            └─ Else ────────> Utility9 "SendEmailBind" (SendEmailPolicy)

Decision29 "pxCreateOperator <orang>?"  <- guard identitas
                            ├─ When IsGroupCreate ─> Assignment15 "ADMIN"
                            └─ Else ──────────────> Decision5 "FAC OUT?"

Decision36 "Letter No Null" ├─ When LetterNoNull ─> Assignment14 "MARKETING (BINDING)"
                            └─ Else ─────────────> Utility10 "SetBanding_ACT"

Decision24 "ISPKSASM"       ├─ When IsPKSASM ─> Assignment8 "UNDERWRITING"
                            └─ Else ─────────> Decision28 "Declaration Policy"
Decision33 "ISPKSASM"       ├─ When IsPKSASM ─> Utility11 "INSERT TO PRODUCTION" (SaveToProduction_ACT)
                            └─ Else ─────────> Utility6 "SET LIMIT_AKSEPTASI" (GetLimitAkseptasi_JUW_UW)
```

**Tahapan siklus renewal `[terverifikasi]`:** Marketing → Group?/PKS?/UW Financial? → tangga
akseptasi berbasis limit → binding → RI Slip → fac out/retro → produksi.
**Tanpa tahap Life, tanpa tahap pembuatan offer baru.**

### 1.3 Workbasket

`[terverifikasi]` 17 Assignment, **seluruhnya `WorkBasket` + `Custom`** — sama seperti NB FacIn.
Label yang muncul sama persis dengan NB FacIn **dikurangi** tiga label Life (`UNDERWRITING LIFE`,
`MEDICAL LIFE`, `DEPTHEAD UNDERWRITING LIFE`) dan `JUNIOR UNDERWITER B` tetap ada (`Assignment4`).
Pemetaan shape → nama workbasket tetap tidak terbaca (OQ-024).

---

## 2. Status / state — pembeda siklus yang terbaca penuh

### 2.1 `Quotation.StatusBusiness` — inilah diskriminator siklus `[terverifikasi]`

Tiga rule `When` beridentitas jelas, **ada di ketiga modul facultative**:

| Rule `When` | Kondisi literal |
| --- | --- |
| `When/IsNB.xml` | `pyWorkPage.Quotation.StatusBusiness = 1` |
| `When/IsRenewal.xml` | `pyWorkPage.Quotation.StatusBusiness = 2` |
| `When/IsEDM.xml` | `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3` |
| `When/IsNotEDM.xml` | `pyWorkPage.Quotation.StatusBusiness != 3` |

Perintah audit:
```
for w in IsNB IsRenewal IsEDM IsNotEDM; do
  for m in "NB FacIn" "RNW Fac In" "Endorsment Fac In"; do
    [ -f "$m/When/$w.xml" ] && grep -oE "<pyLabel>[^<]{1,120}" "$m/When/$w.xml" | grep -v "Available)"
  done
done
```
→ keempat rule **ada di ketiga modul**.

`[terverifikasi]` Nilai literal `StatusBusiness` yang muncul di korpus facultative: **1, 2, 3**.

```
grep -rhoE "StatusBusiness[^<]{0,14}" "NB FacIn" "RNW Fac In" "Endorsment Fac In" --include="*.xml" \
| sed 's/&amp;#61;/=/g;s/&amp;quot;/"/g;s/\]\[/ /g;s/[][]//g;s/ //g' \
| grep -oE 'StatusBusiness[=!<>]+"?[0-9]' | sort | uniq -c | sort -rn
```
→ `==3` 241×, `=3` 153×, `=="3` 133×, `!=3` 125×, `=1` 17×, `=2` 12×, `==2` 9×, `=="1` 9×, `!=2` 9×,
`!=1` 6×, `==1` 4×.

**Ini bukti langsung untuk OQ-015**: ketiga modul membawa **rule pembeda siklus yang sama**, jadi
ketiganya menjalankan satu basis rule yang membedakan tahap siklus **saat runtime**, bukan tiga
basis kode terpisah. `[dugaan]` ekspor per modul adalah pemotongan administratif, bukan pemisahan
aplikasi. **Arti nilai 1/2/3 tetap belum terverifikasi** (OQ-020) — korelasi nama rule bukan bukti.

### 2.2 Prefix ID case sebagai penanda siklus `[terverifikasi]`

```
grep -rhoE '"(NB|RNW|EDM)-"' "NB FacIn" "RNW Fac In" "Endorsment Fac In" --include="*.xml" \
| sort | uniq -c | sort -rn
```
→ `"EDM-"` 262×, `"NB-"` 127×, `"RNW-"` 43×.

`[terverifikasi]` `RNW Fac In/Activity/ProtectRenewal_Act.xml`
(`ASM-FW-GISFW-WORK!PROTECTRENEWAL_ACT`, 46.258 byte, 3 langkah) memuat:

```
precondition langkah: When IsRenewal
ekspresi:             @contains(pyWorkPage.pxInsName,"NB-")
                      Local.DateDif >= 0
langkah 3:            Page-Set-Messages  (Local.ErrMsg)
```

Jadi **ID case membawa penanda siklus** dan diperiksa untuk menolak/menandai kondisi tertentu pada
renewal. Rincian pesan dan arti `Local.DateDif` **belum terverifikasi** — isi `Local.ErrMsg` tidak
terbaca dari tag.

### 2.3 `ProposalAcceptStatus`

`[terverifikasi]` **Sama persis** dengan NB FacIn: delapan Data Transform dengan nilai
`1 / 2 / 3 / 3 / 4 / 7 / 9` dan varian berkondisi. Lihat `NB FacIn.md` §2.1 — tidak diulang.

**Satu selisih terukur:** `RNW Fac In` **tidak punya** `When/IsUWAccepted.xml`; hanya
`DecisionTable/IsUWAccepted.xml`. Perintah audit:

```
for m in "NB FacIn" "RNW Fac In" "Endorsment Fac In"; do ls "$m/When/IsUWAccepted.xml" 2>&1; done
```

`[terverifikasi]` Ini **memperkuat kesimpulan §2.1 `NB FacIn.md`**: di modul yang hanya punya satu
tipe, tipe itu adalah `DecisionTable`. Tetapi isi baris tabelnya tetap tidak ada di ekspor (OQ-043),
sehingga pemetaan nilai → hasil tetap **tidak dapat dinyatakan**.

`[terverifikasi]` Hash `DecisionTable/IsUWAccepted.xml`: NB FacIn = RNW Fac In = `f99bc43c45`;
Endorsment Fac In = `9226a7db77` (**berbeda**).

### 2.4 Guard identitas — identik di ketiga modul

`[terverifikasi]` Rule ber-guard identitas orang di modul ini, dengan **jumlah** identitas
(nilai nama **tidak disalin**, `_METHOD.md` §1.4):

| Rule | Tag | Identitas berbeda | Hash ternormalisasi (3 modul) |
| --- | --- | ---: | --- |
| `When/IsGroup.xml` | `OperatorID.pyUserIdentifier` | **3** | `2cff7b9a67` — **identik** |
| `When/IsGroupCreate.xml` | `pyWorkPage.pxCreateOperator` | **2** | identik |
| `When/IsTBonding.xml` | `.OfferFacIn.QuotationData.MarketingName` | 1 | identik |

Perintah audit (menghitung tanpa menampilkan nilai):
```
for m in "NB FacIn" "RNW Fac In" "Endorsment Fac In"; do
  grep -oE "pyUserIdentifier\]\[&amp;#61;\]\[&amp;quot;[A-Z]+" "$m/When/IsGroup.xml" | sort -u | wc -l
done
```

`[terverifikasi]` Ketiga rule **tidak berkonflik** antar modul — hash ternormalisasinya sama.
Guard identitas facultative adalah **satu himpunan yang sama**, bukan variasi per modul.
→ OQ-021 / OQ-027.

---

## 3. Objek Oracle yang disentuh

```
awk -F'\t' '$1 ~ "^RNW Fac In/" && $2=="RULE-CONNECT-SQL"{print $8}' all-rules.tsv \
  | grep -oE "tables=[^;]*" | sed 's/tables=//' | tr ',' '\n' | sort | uniq -c | sort -rn
```

| Objek | Rule perujuk |
| --- | ---: |
| **`JSON_POLIS`** | 18 |
| `FACINPRODUCTION` | 7 |
| `RW`, `ACCUMULATION` | 6 masing-masing |
| `POOLDATA.REINSURANCETYPE`, `POOLDATA.JSON_POLIS`, `OCCUPATION`, `CURRENCY` | 5 masing-masing |
| `TREATYBUSINESS`, `POOLDATA.M_LIMIT_PROPERTYY`, `MARKETINGOFFICER` | 4 masing-masing |
| `T_STORAGE_IMAGE`, `PROPORTIONALARRG`, `POOLDATA.M_LIMIT_NONPROPANDENGG`, `POOLDATA.FACINPRODUCTION` | 3 masing-masing |

`[terverifikasi]` **Selisih dari NB FacIn:** `POOLDATA.CLIENT` (4 rule di NB FacIn) tidak muncul di
daftar teratas modul ini, dan `TREATYBUSINESS` turun dari 5 ke 4. Sisanya identik. Sebaran objek
secara keseluruhan **hampir sama persis** — konsisten dengan §2.1.

---

## 4. Integrasi eksternal

`[terverifikasi]` Tiga rule `RULE-CONNECT-REST`, sama persis dengan NB FacIn: `ServiceGoogle`,
`convertJsonNusareToProduction`, `getPremiumPaidOn` — seluruhnya `pyBaseURLSelectionType = SETTING`
dengan `pyBaseURLSetting = LinkService!LinkService`, **tanpa URL literal endpoint**.

### 4.1 `serviceInsertArasapasRNW_act` — eksklusif modul ini, dan **rantai integrasinya terbaca**

`[terverifikasi]` `RNW Fac In/Activity/serviceInsertArasapasRNW_act.xml`
(`ASM-FW-GISFW-WORK!SERVICEINSERTARASAPASRNW_ACT`, 89.947 byte, hash `cb70714c1a`) **hanya ada di
modul ini**:

```
find . -path ./OUTPUT_HASIL_RNM -prune -o -type f -iname "serviceInsertArasapasRNW_act.xml" -print
# -> hanya ./RNW Fac In/Activity/serviceInsertArasapasRNW_act.xml
```

Sepuluh langkahnya `[terverifikasi]`:

| # | Method | Catatan |
| ---: | --- | --- |
| 1 | `Call serviceInsertArasapas_act` | class `ASM-FW-GISFW-Work` → **terblokir, §5.1** |
| 2 | `Property-Set` | |
| 3 | `RDB-List` | ke `ASM-FW-GISFW-Data-Search` / `Code-Pega-List` |
| 4 | `Property-Set` | |
| 5 | **`Call ASM-FW-GISFW-Int-M_LINK_SERVICE.GetLinkService`** | **resolusi endpoint — §4.2** |
| 6 | **`Connect-REST`** | panggilan keluar |
| 7 | `Property-Set` | |
| 8 | `Call InsertLogServiceProd` | pencatatan log |
| 9 | `Page-Remove` | |
| 10 | `RDB-List` | |

Class yang dirujuk termasuk **`ASM-FW-GISFW-Work-Renewal`** `[terverifikasi]` — satu-satunya
kemunculan class bersufiks `-Renewal` yang menyentuh alur; di dataset korpus class itu hanya memiliki
**satu** rule lain, `RNW Fac In/ReportDefinition/RenewalList_RD.xml`. Perintah audit:

```
awk -F'\t' 'toupper($3) ~ /RENEWAL/ {print $2"\t"$3"\t"$5"\t"$1}' all-rules.tsv | sort -u
```

### 4.2 Endpoint diambil dari **tabel database**, bukan hanya SystemSettings — temuan baru

`[terverifikasi]` `RNW Fac In/Activity/GetLinkService.xml` beridentitas
**`ASM-FW-GISFW-INT-M_LINK_SERVICE!GETLINKSERVICE`**, empat langkah:
`Page-New`, **`Obj-Browse`**, `Property-Set`, `Page-Remove`. Ia mem-browse class
`ASM-FW-GISFW-Int-M_LINK_SERVICE` dengan kunci **`KATEGORI_1`** dan **`KATEGORI_2`**.

Perintah audit:
```
grep -oE "<pyStepsActivityName>[^<]*|<pyRuleName>[^<]*" "RNW Fac In/Activity/GetLinkService.xml" \
| sed 's/<[^>]*>//' | sort -u
```

`[terverifikasi]` Jadi alamat layanan keluar **berada di tabel Oracle `M_LINK_SERVICE`**, dicari
dengan dua kolom kategori — **bukan** semata di `RULE-ADMIN-SYSTEM-SETTINGS`. Isi tabel itu
**tidak ada di korpus**, sehingga daftar endpoint sesungguhnya **tidak diketahui**.

Ini memperluas OQ-018 secara material: sapuan D1 hanya mencari URL literal di dalam file rule;
**sumber konfigurasi yang sebenarnya ada di database**. → **OQ-047**.

Sembilan rule modul ini memakai `M_LINK_SERVICE` `[terverifikasi]`
(`grep -rl "M_LINK_SERVICE" "RNW Fac In" --include="*.xml" | wc -l`), termasuk
`DeleteGoogleStorage_Act`, `GetUrlGoogleStorage_Act`, `InsertGoogleStorage_Act`,
`GeminiAIGoogle_Act`, `serviceInsertArasapasEDM_act`.

### 4.3 Google Storage dan `GeminiAIGoogle_Act`

`[terverifikasi]` Empat activity integrasi Google di modul ini: `InsertGoogleStorage_Act`,
`GetUrlGoogleStorage_Act`, `DeleteGoogleStorage_Act`, **`GeminiAIGoogle_Act`**. Ketiganya yang
pertama juga ada di `NB FacIn` dan `Endorsment Fac In`; **`GeminiAIGoogle_Act` hanya ada di
`NB FacIn` dan `RNW Fac In`, tidak di `Endorsment Fac In`** `[terverifikasi]`:

```
for m in "NB FacIn" "RNW Fac In" "Endorsment Fac In"; do ls "$m/Activity/" | grep -iE "google|gemini"; done
```

Isi keduanya **belum dibaca**. Bahwa ada pemanggilan model AI pihak ketiga di dalam alur
underwriting adalah fakta yang dicatat; apa yang dikirim dan apa yang dilakukan hasilnya **belum
terverifikasi** → **OQ-048**.

---

## 5. Batas pengetahuan

### 5.1 Blocker `serviceInsertArasapas_act` berlaku juga di sini

`[terverifikasi]` `RNW Fac In/Activity/serviceInsertArasapas_act.xml` beridentitas
`ASM-FW-GISFW-DATA-POLICYTREATYIN!SERVICEINSERTARASAPAS_ACT`, hash `b013027e0c` — **sama persis**
dengan salinan di `NB Treaty In` dan `NB FacIn`. Ia memanggil
`ASM-FW-GISFW-WORK!SERVICEINSERTARASAPAS_ACT` yang **tidak ada di modul ini**.

Telusur **dihentikan** di langkah 1 `serviceInsertArasapasRNW_act` (§4.1). Varian dari modul lain
**tidak dipinjam** → OQ-025. Rincian lengkap di `NB FacIn.md` §5.3.

### 5.2 Stored procedure — 26 dipanggil

`[terverifikasi]` Sama dengan daftar `NB FacIn.md` §5.2 **dikurangi tiga**:
`POOLDATA.PEGA_JSON_POLIS_TREATYIN`, `POOLDATA.PEGA_TREATY_IN`, `POOLDATA.PEGA_M_JSON_OFFER`
**tidak dipanggil** dari modul ini.

Perintah audit:
```
diff <(awk -F'\t' '$1~"^NB FacIn/" && $2=="RULE-CONNECT-SQL"{print $8}' all-rules.tsv \
       | grep -oE "procs=[^;]*" | sed 's/procs=//' | tr ',' '\n' | sort -u) \
     <(awk -F'\t' '$1~"^RNW Fac In/" && $2=="RULE-CONNECT-SQL"{print $8}' all-rules.tsv \
       | grep -oE "procs=[^;]*" | sed 's/procs=//' | tr ',' '\n' | sort -u)
```

`[terverifikasi]` **Konsisten dengan §1.1**: modul renewal tidak menyentuh jalur treaty-in maupun
pembuatan JSON offer baru. Isi seluruh procedure tetap tidak ada di korpus (OQ-002).

### 5.3 Database link — **tujuh objek, bukan delapan**

`[terverifikasi]` Modul ini mengakses **7** objek lewat `@ASMD.SINARMAS.CO.ID`:
`M_EQS_RATE`, `M_TERORISME_RATE`, `M_RSMD_RATE`, `M_FLEXAS_RATE`, `M_FLOOD_AREA`,
`FIRE.M_FLOOD_RATE`, `FIRE.M_BI_INDEMNITY`.

**`LST_KURS_STANDARD` tidak dirujuk dari modul ini** — ia dirujuk `NB FacIn` dan
`Endorsment Fac In`. Ini mencocokkan angka OQ-017 (NB FacIn 8, RNW Fac In 7, Endorsment 8).

`[dugaan]` renewal memakai kurs yang sudah tersimpan pada data lama alih-alih mengambil kurs standar
baru. **Belum terverifikasi** — tidak ada rule yang membuktikannya.

### 5.4 Spreading — 34 activity, tak satu pun habis dibaca

`[terverifikasi]` `ls -l "RNW Fac In/Activity/" | grep -icE "spread|capacit|scoring"` → **34**
(NB FacIn 38). Yang **tidak ada** di sini tetapi ada di NB FacIn: `BreakDownSpreading_Act`,
`CountSpreading_Act`, `SpreadingAdditionalProtection`, `SumTreatyCapacity_Act`.

Ukuran terbesar sama: `SumTSIPremiSpreadedRNM_FIRE_Act.xml` **996.052 byte** — **belum dibaca**.
Rumus spreading/capacity/scoring **tidak dinyatakan** (batas cakupan telusur, `NB FacIn.md` §5.4).

### 5.5 `SetOldData` — mesin data lama, identik di ketiga modul

`[terverifikasi]` `RNW Fac In/Activity/SetOldData.xml` (`ASM-FW-GISFW-WORK!SETOLDDATA`,
**450.829 byte**, **59 langkah `Property-Set`**). Hash `df291257aa` — **identik di NB FacIn,
RNW Fac In, dan Endorsment Fac In**.

```
for m in "NB FacIn" "RNW Fac In" "Endorsment Fac In"; do sh nhash.sh "$m/Activity/SetOldData.xml"; done
```

**Isinya belum habis dibaca** — 59 penugasan properti. Yang dicatat: rule ini adalah mekanisme
salinan data lama (*before-image*) dan **dipakai bersama oleh ketiga siklus**, bukan khas renewal
maupun endorsement.

Activity data-lama lain yang ada di modul ini dan **belum dibaca**:
`GetOldDataRetro_ACT` (55.459 B, 4 langkah: `RDB-List`, `Property-Set`, `Obj-Open-By-Handle`,
`Property-Set`), `ViewOldDataEDM_Act` (213.120 B, 26 langkah), `getOldDeduct`.

`[terverifikasi]` Jumlah rule yang menyebut `OldData`: NB FacIn **129**, RNW Fac In **114**,
Endorsment Fac In **133**.
```
for m in "NB FacIn" "RNW Fac In" "Endorsment Fac In"; do grep -rl "OldData" "$m" --include="*.xml" | wc -l; done
```

### 5.6 Prorate

`[terverifikasi]` `ProRateType` diuji dengan nilai literal **3** dan **4** di modul ini
(`ProRateType==3` 4×, `!=4` 3×, `==4` 2×, `!=3` 2×, `=="3"` 2×, `=="4"` 2×). Frekuensinya **lebih
rendah** dari `Endorsment Fac In` (23 kemunculan berkondisi vs 9 di sini).

Perintah audit:
```
grep -rhoE "ProRateType[^<]{0,24}" "RNW Fac In" --include="*.xml" \
| sed 's/&amp;#61;/=/g;s/&amp;quot;/"/g;s/\]\[/ /g;s/[][]//g' | sort | uniq -c | sort -rn
```

**Prorate renewal tidak muncul sebagai mekanisme tersendiri di modul ini.** Procedure
`FIRE.CEK_PRORATA_TANGGAL` **tidak dipanggil** dari `RNW Fac In` — ia hanya dipanggil dari
`Endorsment Fac In` (§5.2 `Endorsment Fac In.md`). Arti `ProRateType` **belum terverifikasi**.

---

## 6. Rule yang terlibat

| Class / Nama / Tipe | Path (relatif `RNW Fac In/`) | OQ-011? |
| --- | --- | --- |
| `ASM-FW-GISFW-WORK` / `INPUTRENEWALFACULTATIVEIN` / `RULE-OBJ-FLOW` | `Flow/InputRenewalFacultativeIn.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `INPUTINWARDFACULTATIVERISLIP` / `RULE-OBJ-FLOW` | `Flow/InputInwardFacultativeRISlip.xml` | tidak — identik NB FacIn |
| `ASM-FW-GISFW-WORK` / `OFFERFACOUT` / `RULE-OBJ-FLOW` | `Flow/OfferFacOut.xml` | tidak — identik NB FacIn |
| `ASM-FW-GISFW-WORK` / `OFFERFACRETRO` / `RULE-OBJ-FLOW` | `Flow/OfferFacRetro.xml` | **YA — #338**, varian NB/RNW (identik NB FacIn) |
| `ASM-FW-GISFW-WORK` / `ISNB`, `ISRENEWAL`, `ISEDM`, `ISNOTEDM` / `RULE-OBJ-WHEN` | `When/` | tidak — **diskriminator siklus, §2.1** |
| `ASM-FW-GISFW-WORK` / `PROTECTRENEWAL_ACT` / `RULE-OBJ-ACTIVITY` (46.258 B) | `Activity/ProtectRenewal_Act.xml` | tidak |
| `ASM-FW-GISFW-WORK` / `SERVICEINSERTARASAPASRNW_ACT` / `RULE-OBJ-ACTIVITY` (89.947 B) | `Activity/serviceInsertArasapasRNW_act.xml` | tidak — **eksklusif modul ini** |
| `ASM-FW-GISFW-INT-M_LINK_SERVICE` / `GETLINKSERVICE` / `RULE-OBJ-ACTIVITY` | `Activity/GetLinkService.xml` | tidak — **§4.2** |
| `ASM-FW-GISFW-DATA-POLICYTREATYIN` / `SERVICEINSERTARASAPAS_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/serviceInsertArasapas_act.xml` | tidak (memanggil identitas berkonflik — §5.1) |
| `ASM-FW-GISFW-WORK` / `SETOLDDATA` / `RULE-OBJ-ACTIVITY` (450.829 B) | `Activity/SetOldData.xml` | tidak — identik 3 modul |
| `ASM-FW-GISFW-WORK` / `GETOLDDATARETRO_ACT` / `RULE-OBJ-ACTIVITY` | `Activity/GetOldDataRetro_ACT.xml` | tidak |
| `ASM-FW-GISFW-DATA-QUOTATION` / `GETBUSINESSGROUP_ACT` / `RULE-OBJ-ACTIVITY` (46.765 B) | `Activity/GetBusinessGroup_Act.xml` | tidak — **eksklusif modul ini** |
| `ASM-FW-GISFW-WORK` / `GETLIMITAKSEPTASI_ACT`, `GETLIMITAKSEPTASI_JUW_UW` / `RULE-OBJ-ACTIVITY` | `Activity/` | tidak |
| `ASM-FW-GISFW-WORK` / `ISUWACCEPTED` / `RULE-DECLARE-DECISIONTABLE` | `DecisionTable/IsUWAccepted.xml` | tidak (isi baris tak ada — OQ-043) |
| `ASM-FW-GISFW-WORK-RENEWAL` / `RENEWALLIST_RD` / `RULE-OBJ-REPORT-DEFINITION` | `ReportDefinition/RenewalList_RD.xml` | tidak — satu-satunya rule di class `-RENEWAL` |
| `ASM-FW-GISFW-WORK` / `RENEWAL_FLOWACT`, `RENEWAL_FLOWACT_ISUW` / `RULE-OBJ-FLOWACTION` | `FlowAction/` | sebagian |
| `ASM-FW-GISFW-WORK` / `ISGROUP`, `ISGROUPCREATE`, `ISTBONDING` / `RULE-OBJ-WHEN` | `When/` | tidak — **guard identitas, §2.4** |
| `ASM-FW-GISFW-WORK` / `ISPKSASM`, `LETTERNONULL` / `RULE-OBJ-WHEN` | `When/` | tidak — kondisi tak terbaca |
| 3 rule `RULE-CONNECT-REST` (§4) | `ConnectREST/` | sebagian |

**28 rule ditelusur** (4 Flow — 3 dipakai ulang lewat bukti hash; 8 When; 9 Activity;
1 DecisionTable; 1 ReportDefinition; 2 FlowAction; 3 ConnectREST).

`[terverifikasi]` **Activity yang eksklusif RNW Fac In terhadap NB FacIn hanya TIGA**:
`GetBusinessGroup_Act`, `SetDataInsuredEDM_Act`, `serviceInsertArasapasRNW_act`. Sebaliknya,
**57 activity NB FacIn tidak ada di sini** — sebagian besar terkait treaty inward
(`SaveJsonPolisTreatyIn_Act`, `SetTreatyIn_Act`, `TreatyIn*`, `GeneratePolicyNoTreaty_Act`,
`InsertToTreatyXOLList`, `FetchTreatyGroup*`) dan perhitungan premi/komisi
(`CalculatePremi_Act`, `CountNetPremi_act`, `CountResult1_Act`, `CountResult2Ogp_act`,
`CountRiCommOgp_act`, `CountOverridingCommOgp_Act`, `SetPPNPPH`, `SetReinstatementPct`).

Perintah audit:
```
comm -23 <(ls "RNW Fac In/Activity" | sort) <(ls "NB FacIn/Activity" | sort)   # 3 baris
comm -13 <(ls "RNW Fac In/Activity" | sort) <(ls "NB FacIn/Activity" | sort)   # 57 baris
```

---

## 7. Pertanyaan terbuka

| OQ | Ringkas | Pemilik |
| --- | --- | --- |
| **OQ-047** | Alamat layanan keluar diambil dari **tabel Oracle `M_LINK_SERVICE`** (kunci `KATEGORI_1`/`KATEGORI_2`), bukan hanya SystemSettings — daftar endpoint tidak ada di korpus | DBA + Platform |
| **OQ-048** | `GeminiAIGoogle_Act` memanggil model AI pihak ketiga di dalam alur underwriting (NB FacIn & RNW Fac In, bukan Endorsment) — apa yang dikirim, apa yang dilakukan hasilnya | Product+UW + Security |

OQ dikuatkan: **OQ-015** (§2.1 — **bukti terkuat sejauh ini**), OQ-002 (§5.2 — 26 SP),
OQ-011 (§2.3 — hash DecisionTable berbeda di Endorsment), OQ-012, OQ-017 (§5.3 — 7 objek, bukan 8),
OQ-018 (§4.2 — **meluas ke konfigurasi di database**), OQ-020 (§2.1 — arti 1/2/3 tetap terbuka),
OQ-021/OQ-027 (§2.4), OQ-024 (§1.3), **OQ-025** (§5.1 — meluas ke modul ini),
OQ-026 (§2.3 — dipersempit lagi), OQ-028 (§1.3), OQ-029 (§6), OQ-040, OQ-043 (§2.3).
