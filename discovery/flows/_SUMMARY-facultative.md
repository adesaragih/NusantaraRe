# Sintesis Lintas-Modul — Domain Facultative Inward

STEP D2 Tahap 4. Disusun 2026-09-13 setelah menelusur **NB FacIn**, **RNW Fac In**, dan
**Endorsment Fac In**. Konvensi: `_METHOD.md`.

Dokumen ini **menyajikan bukti dan membandingkan**. Ia **tidak** menetapkan bounded context —
itu STEP D4, setelah gate manusia.

**Cakupan:** 6.071 file (2.083 + 1.927 + 2.061), **domain terbesar korpus** (64,8% dari 9.369 file).
Perintah audit: `for m in "NB FacIn" "RNW Fac In" "Endorsment Fac In"; do find "$m" -name '*.xml' | wc -l; done`

---

## 1. Dua belas file `Flow`, **delapan isi berbeda**

`[terverifikasi]` Hash ternormalisasi 18 tag atas seluruh file `Flow` ketiga modul:

| Hash | File | Ada di | Telusur |
| --- | --- | --- | --- |
| `d30995fece` | `InputQuotation` | NB FacIn | `NB FacIn.md` §1.1 |
| `39cfbbec4f` | `InputInwardFacultativeOffer` | NB FacIn | `NB FacIn.md` §1.2 |
| `c099bf4ebb` | `InputInwardFacultativeRISlip` | NB FacIn **=** RNW Fac In | `NB FacIn.md` §1.3 — **sekali** |
| `832fb12b9d` | `OfferFacOut` | NB FacIn **=** RNW Fac In | `NB FacIn.md` §1.4 — **sekali** |
| `bfd6252070` | `OfferFacRetro` | NB FacIn **=** RNW Fac In | `NB FacIn.md` §1.5 — **varian A** |
| **`b370146c63`** | `OfferFacRetro` | **Endorsment Fac In** | `Endorsment Fac In.md` §2 — **varian B** |
| `1306f68d56` | `InputRealizationTreatyIn` | NB FacIn **=** NB Treaty In | `NB Treaty In.md` — **dipakai ulang** |
| `f3dff3940e` | `InputRenewalFacultativeIn` | RNW Fac In | `RNW Fac In.md` §1 |
| `f0b2f4ea86` | `InputAddendumFacultativeIn` | Endorsment Fac In | `Endorsment Fac In.md` §1 |

Perintah audit:
```
sh nhash.sh "NB FacIn/Flow/"*.xml "RNW Fac In/Flow/"*.xml "Endorsment Fac In/Flow/"*.xml \
            "NB Treaty In/Flow/InputRealizationTreatyIn.xml" | sort
```

`[terverifikasi]` **Empat file Flow ditelusur sekali saja** karena hash-nya membuktikan isinya
identik; `OfferFacRetro` ditelusur **dua kali** karena hash-nya membuktikan isinya berbeda.

---

## 2. Tiga siklus dibandingkan

### 2.1 Tahapan

| Tahap | NB (new business) | RNW (renewal) | EDM (endorsement) |
| --- | :---: | :---: | :---: |
| Dispatcher FacIn/TreatyIn (`InputQuotation`) | **ya** | tidak | tidak |
| Marketing → gerbang akseptasi | ya | ya | ya |
| Cabang **Life** lengkap (Medical, UW Life, DeptHead UW Life) | **ya** | **tidak** | sebagian — hanya `Decision19 "IS LIFE?"` |
| Gerbang `ISPKSASM` | ya (2×) | ya (2×) | ya (2×) |
| Gerbang `UW FINANCIAL?` / `IsTBonding` | ya | ya | ya |
| Tangga akseptasi berbasis `LetterNo` | ya | ya | ya |
| Tangga `Which team?` (UW / Senior UW / JUW A / DepHead UW) | ya | ya | ya |
| Banding (`SetBanding_ACT` + `LetterNoNull`) | ya | ya | ya |
| **RI Slip** (`InputInwardFacultativeRISlip`) | ya | ya | **tidak** (tak ada flow-nya) |
| **Check Spreading** sebagai shape | ya (`CekLimitSpreading_Act`) | ya | **tidak** (`CekLimitSpreading_Act` tidak ada di modul) |
| Fac out / retro (`OfferFacRetro`) | ya — tangga **2** | ya — tangga **2** | ya — tangga **4** |
| Screen flow `OfferFacOut` (WorkList) | ya | ya | **tidak** |
| Gerbang khas EDM (`IsEDMRiSlip`, `IsEdmPPNPPH`, `IsEdmInternalRetro`) | tidak | tidak | **ya** |
| Simpan | `SaveJsonOfferFacIn_Act` + `SaveJsonPolicyFacIn_Act` | `SaveJsonPolicyFacIn_Act` | **`SaveEDMToJsonPolicy_Act`** |
| Arasapas | `serviceInsertArasapas_act` | `serviceInsertArasapasRNW_act` | `serviceInsertArasapasEDM_act` |
| Cek status konversi (`IsSuccessHitService`) | ya | ya | ya |

### 2.2 Ukuran graf

`[terverifikasi]` (`grep -o "<pyShapeType>[^<]*" | sort | uniq -c`)

| Flow utama | Assignment | Decision | Utility | SubProcess | Connector |
| --- | ---: | ---: | ---: | ---: | ---: |
| `InputInwardFacultativeOffer` (NB) | **21** | **42** | **12** | 3 | **170** |
| `InputInwardFacultativeRISlip` (NB/RNW) | 16 | 32 | 8 | 4 | 119 |
| `InputRenewalFacultativeIn` (RNW) | 17 | 36 | 11 | **5** | 140 |
| `InputAddendumFacultativeIn` (EDM) | 15 | 37 | 9 | 3 | 133 |
| `OfferFacRetro` varian A (NB/RNW) | 3 | 4 | 2 | 0 | 13 |
| `OfferFacRetro` varian B (EDM) | **5** | **6** | 2 | 0 | **19** |
| `OfferFacOut` (NB/RNW) | 2 | 1 | 0 | 0 | 2 |
| `InputQuotation` (NB) | 0 | 4 | 2 | 2 | 6 |

`InputInwardFacultativeOffer` adalah **graf terbesar korpus** — sebagai pembanding, flow terbesar
domain Claim (`Claim Fac In/Register_Flow`) punya 7 shape.

---

## 3. Titik `ProposalAcceptStatus` — bukti dan batasnya

### 3.1 Yang terverifikasi

`[terverifikasi]` Enam hasil yang dideklarasikan `DecisionTable/IsUWAccepted.xml` lewat
`<pyTaskStatusXml>`, **sama di ketiga modul**:

```
confirm   reject   ask   banding   revise   decline
```

`[terverifikasi]` Nilai literal yang **ditulis** ke `.ProposalAcceptStatus`:

| Data Transform | Nilai | NB FacIn | RNW Fac In | Endorsment Fac In |
| --- | ---: | :---: | :---: | :---: |
| `SetAkseptasiProposal` | `1` | ada | ada | ada |
| `SetRejectProposal` | `2` | ada | ada | ada |
| `SetAskProposal_DT` | `3` | ada | ada | ada |
| `SetAkseptasiCeding_DT` | `3` | ada | ada | **tidak ada** |
| `SetBandingProposal_DT` | `4` | ada | ada | ada |
| `SetDeclineProposal_DT` | `7` | ada | ada | ada |
| `SetReviseProposal` | `9` | ada | ada | **tidak ada** |
| `SetAkseptasiProposal_DT` | `1,2,3,9,4` berkondisi | ada | ada | ada |

`[terverifikasi]` Nilai yang **diuji**: `When/IsUWAccepted.xml` → `= "1"`;
`When/IsFacout.xml` → `= 4`. Di flow, connector membawa `<pyConditionType>Status` dengan keenam
nilai hasil di atas.

### 3.2 Yang **tidak** dapat dinyatakan

**Pemetaan `1/2/3/4/7/9` → `confirm/reject/ask/banding/revise/decline` tidak ada di korpus.**
Baris tabel keputusan tidak ikut terekspor di **seluruh 49 file `DecisionTable`** korpus (OQ-043).
Korelasi nama rule dengan nama hasil bersifat sugestif; `_METHOD.md` §2.2 melarang memakai nama
sebagai bukti perilaku. **Tidak ditebak.** → OQ-020, OQ-043.

**Nilai `4` ambigu `[terverifikasi]`:** ditulis `SetBandingProposal_DT` (banding) dan diuji
`When/IsFacout.xml` (fac out). Satu nilai, dua nama pemakai → **OQ-044**.

`[terverifikasi]` `DecisionTable/IsUWAccepted.xml` **berbeda isinya** di Endorsment Fac In
(`9226a7db77`) dibanding NB FacIn/RNW Fac In (`f99bc43c45`). Jadi aturan persetujuan **memang
bercabang antar siklus** — **apa** yang bercabang tidak dapat dibaca.

### 3.3 Rantai `LetterNo` yang **utuh terbaca** — pengecualian yang berguna

`[terverifikasi]` Berbeda dari §3.2, satu rantai keputusan dapat dibaca ujung ke ujung:

```
Decision "Accept?" -> hasil banding
  -> SetBandingProposal_DT:  .ProposalAcceptStatus := 4
                             .OfferFacIn.IsBanding := "true"
                             .LetterNo             := .BandingTo
                             FlagBanding.CARI11    := .BandingTo
  -> Decision "Limit Akseptasi" menguji When To*:
       LetterNo = "KADIVTEKNIK"        -> Assignment KADIV TEKNIK
       LetterNo = "MANAGERTEKNIK"      -> Assignment MANAGER TEKNIK
       LetterNo = "KADIVFINANCIAL"     -> Assignment KADIV FINANCIAL
       LetterNo = "KADIVFACULTATIVE"   -> Assignment KADIV FACULTATIVE
       LetterNo = "DEPHEADUNDERWRITER" -> Assignment DEP HEAD UW FACULTATIVE
       LetterNo = "DIREKTURTEKNIK"     -> Assignment DIREKTUR TEKNIK
       LetterNo = "DIREKTURMARKETING"  -> Assignment DIREKTUR MARKETING
       LetterNo = "SENIORUW"           -> Assignment SENIOR UNDERWRITING
       LetterNo = "DEPTHEADUWLIFE"     -> Assignment DEPTHEAD UNDERWRITING LIFE
```

**Field nomor surat adalah mekanisme routing persetujuan.** Pola sejajar dengan
`OperatorID.pyTelephone` di domain treaty (OQ-027). → **OQ-045**.

---

## 4. Fac out / retro

`[terverifikasi]` Tiga gerbang berulang di ketiga siklus:

| `When` | Kondisi literal | Guna |
| --- | --- | --- |
| `IsFacRetro` | `.OfferFacIn.IsFacRetro = 1` | masuk jalur fac out/retro |
| `IsInputFacRetro` | `.OfferFacIn.IsInputFacRetro = 1` | "sudah pernah input fac retro" |
| `IsNotPrintRISlip` | `.OfferFacIn.IsRISlip != 1` | gerbang cetak slip |

`[terverifikasi]` Empat workbasket fac out: `ReasFacOutAdmin` (47), `ReasFacOutHead` (42),
`ReasFacOutGroupLeader` (13), `ReasFacOutTechnicalDirector` (13), ditambah `ReasFacOut` (132).

**Selisih perilaku terbesar antar siklus ada di sini:**

| | NB / RNW (`bfd6252070`) | EDM (`b370146c63`) |
| --- | --- | --- |
| Anak tangga | **2** (`Admin` → `Head`) | **4** (`Admin` → `Head` → `GroupLeader` → `TechnicalDirector`) |
| Decision `IsUWAccepted` | 2 | 4 |
| Tujuan `reject` di tangga ≥2 | `END52` / `Assignment1` | **seluruhnya kembali ke `Assignment1`** |

`[terverifikasi]` `FACOUTPRODUCTION` dirujuk 4 rule di Endorsment Fac In dan **tidak** masuk daftar
teratas di dua modul lain — sejalan dengan tangga yang lebih panjang.

`[terverifikasi]` Penomoran fac retro lewat procedure `POOLDATA.GENERATE_FACRETRO_NO`, dipanggil
dari **ketiga** modul. Isinya tidak ada di korpus (OQ-002).

`[terverifikasi]` Satu rumus retro **terbaca penuh dan identik di ketiga modul** (`8f1268a660`):
`CountPremiNusareRetro_Act` — pangsa **0,45** / **0,05**, ambang **3.000.000.000** (dibandingkan
sebagai string), dan cabang `ProRateType != 3`. **Mata uang tidak disebut; arti belum
terverifikasi.** → OQ-046. Rincian di `NB FacIn.md` §5.5.

---

## 5. Spreading, capacity, scoring — batas cakupan telusur

`[terverifikasi]` Jumlah activity bernama spreading/capacity/scoring:

| Modul | Jumlah | Terbesar | Ukuran |
| --- | ---: | --- | ---: |
| NB FacIn | **38** | `SumTSIPremiSpreadedRNM_FIRE_Act.xml` | **996.052** |
| RNW Fac In | **34** | idem | 996.052 |
| Endorsment Fac In | **38** | idem | 995.970 |

Perintah audit: `ls -l "<modul>/Activity/" | grep -icE "spread|capacit|scoring"`

`[terverifikasi]` Varian per lini bisnis terbaca dari penamaan: `_FIRE`, `_ANEKA`, `_GOLF`,
`_MARINECARGO`, `_MBU`, `_PA`, `_TRAVEL`.

Selisih antar modul `[terverifikasi]`:
- **NB FacIn saja:** `BreakDownSpreading_Act`, `CountSpreading_Act`, `SpreadingAdditionalProtection`
- **NB FacIn + Endorsment:** `SumTreatyCapacity_Act`
- **NB FacIn + RNW:** `CekLimitSpreading_Act` (dipanggil shape "Check Spreading")
- **Endorsment saja:** `CopySpreading_Act` (627.836 B), `SetSpreadingCoverage_ACT`,
  `SetSumTSISpreading_Act`, `SpreadingProtection`

**Tak satu pun habis dibaca.** Rumus spreading, capacity, dan scoring **tidak dinyatakan** di
artefak D2 mana pun. Ini **batas cakupan telusur** (volume), bukan batas pengetahuan korpus —
berbeda dari SP dan db-link.

`[terverifikasi]` Limit akseptasi **diambil dari database**: `GetLimitAkseptasi_ActFlow` (NB FacIn)
memanggil **12 RDB-List** berbeda di class `ASM-FW-GISFW-Int-policyjson`, terpisah per lini bisnis
(Prefered Comm, Non Prefer, Non Fire, Kredit NCL, Kredit CL, Bond, Engineering) dengan **varian
"Banding" tersendiri**. Tabel limit: `POOLDATA.M_LIMIT_PROPERTYY`,
`POOLDATA.M_LIMIT_NONPROPANDENGG`. → OQ-040.

`[terverifikasi]` Input rating berada **di luar sistem**: 8 objek lewat `@ASMD.SINARMAS.CO.ID`
(NB FacIn 8, RNW Fac In **7** — tanpa `LST_KURS_STANDARD`, Endorsment 8). → OQ-017.

`[terverifikasi]` Perhitungan prorata tanggal juga di luar korpus:
**`FIRE.CEK_PRORATA_TANGGAL`** dipanggil **hanya dari Endorsment Fac In**
(`RDBList/SearchSQLRateKPR.xml`, `SearchSQLRateNonKPR.xml`). → OQ-002.

---

## 6. Integrasi

`[terverifikasi]` **Tiga rule `RULE-CONNECT-REST` yang sama di ketiga modul**: `ServiceGoogle`,
`convertJsonNusareToProduction`, `getPremiumPaidOn` — seluruhnya `pyBaseURLSelectionType = SETTING`
dengan `pyBaseURLSetting = LinkService!LinkService`. **Tidak ada URL literal endpoint** (satu-satunya
URL adalah `<pyHelpURI>` dokumentasi Pega). Berbeda dari `Claim Fac In` yang punya 8 ConnectREST
dengan satu URL literal → OQ-018.

### 6.1 Endpoint sesungguhnya ada di **tabel database** — temuan Tahap 4

`[terverifikasi]` `Activity/GetLinkService.xml` beridentitas
`ASM-FW-GISFW-INT-M_LINK_SERVICE!GETLINKSERVICE`, melakukan `Obj-Browse` atas class
`ASM-FW-GISFW-Int-M_LINK_SERVICE` dengan kunci **`KATEGORI_1`** dan **`KATEGORI_2`**. Ia dipanggil
dari activity Arasapas ketiga siklus dan dari seluruh activity Google Storage.

**Daftar endpoint sesungguhnya berada di tabel Oracle `M_LINK_SERVICE`, yang isinya tidak ada di
korpus.** Sapuan OQ-018 di D1 hanya mencari URL literal di dalam file rule, sehingga sumber
konfigurasi yang sebenarnya luput. → **OQ-047**.

`[terverifikasi]` Rule SystemSettings `LINKSERVICE!LINKSERVICE` (`RULE-ADMIN-SYSTEM-SETTINGS`)
terdaftar berkonflik di OQ-011 entri **#324**, tetapi pengukuran ulang §8 menunjukkan ia
**identik** di ketiga modul setelah tiga tag metadata dinormalisasi. Konfigurasi base URL **tidak**
bercabang antar siklus.

### 6.2 Arasapas — tiga pintu, satu implementasi yang terkunci

`[terverifikasi]`

| Siklus | Activity pembungkus | Implementasi `ASM-FW-GISFW-WORK!SERVICEINSERTARASAPAS_ACT` ada di modul? |
| --- | --- | --- |
| NB | `serviceInsertArasapas_act` (18.013 B, `…DATA-POLICYTREATYIN`) | **tidak** → blocker |
| RNW | `serviceInsertArasapasRNW_act` (89.947 B) → memanggil pembungkus yang sama | **tidak** → blocker |
| EDM | `serviceInsertArasapasEDM_act` **dan** `serviceInsertArasapas_act` (231.975 B) | **ya** → telusur jalan |

`[terverifikasi]` Class sasaran panggilan kini **tertulis eksplisit** di
`<pyStepsClassName>ASM-FW-GISFW-Work` — menaikkan temuan `NB Treaty In.md` §5.2 dari `[dugaan]`
menjadi `[terverifikasi]`.

`[terverifikasi]` **Kedua varian implementasi yang ada di korpus (EDM Treaty In `ab3ae3c9ba` dan
Endorsment Fac In `7314b6c49c`) berbeda hanya pada 3 tag metadata Pega** — `<pyDelete>`,
`<pyVersionSecure>`, `<pzIsPrivateCheckOut>`. Urutan 20 langkah dan daftar rule yang dirujuk
**identik**. Jadi konflik OQ-011 #415 adalah **artefak check-out/versi, bukan percabangan
perilaku**. Blocker OQ-025 tetap berdiri secara prosedural untuk 3 modul yang tidak memuat
implementasinya, tetapi **risiko salah tafsirnya jauh lebih kecil** dari perkiraan D1.

Rincian 20 langkah (varian Endorsment) di `Endorsment Fac In.md` §5.1.

### 6.3 Model AI pihak ketiga di dalam alur

`[terverifikasi]` `GeminiAIGoogle_Act` ada di **NB FacIn** dan **RNW Fac In**, **tidak** di
Endorsment Fac In. Tiga activity Google Storage (`Insert`, `GetUrl`, `Delete`) ada di ketiganya.
Isinya **belum dibaca**. → **OQ-048**.

---

## 7. Konfirmasi temuan D1: **satu ruleset, diskriminator siklus** — KEBALIKAN dari Claim

### 7.1 Bukti D1 (OQ-015)

`[terverifikasi]` 6.071 file ketiga modul hanya memuat **2.552 identitas unik**;
**1.567 identitas (61,4%) ada di ketiganya**; **1.906 dari 1.926 identitas RNW Fac In (99,0%) juga
ada di NB FacIn**.

### 7.2 Bukti baru dari D2 Tahap 4 — **konfirmasi, dengan mekanismenya**

`[terverifikasi]` **Diskriminator siklusnya terbaca sebagai rule, dan ada di ketiga modul:**

| Rule `When` | Kondisi literal | NB FacIn | RNW Fac In | Endorsment Fac In |
| --- | --- | :---: | :---: | :---: |
| `IsNB` | `Quotation.StatusBusiness = 1` | **ada** | **ada** | **ada** |
| `IsRenewal` | `Quotation.StatusBusiness = 2` | **ada** | **ada** | **ada** |
| `IsEDM` | `OfferFacIn.QuotationData.StatusBusiness = 3` | **ada** | **ada** | **ada** |
| `IsNotEDM` | `Quotation.StatusBusiness != 3` | **ada** | **ada** | **ada** |

Perintah audit:
```
for w in IsNB IsRenewal IsEDM IsNotEDM; do
  for m in "NB FacIn" "RNW Fac In" "Endorsment Fac In"; do ls "$m/When/$w.xml"; done
done
```
→ 12 file, seluruhnya ada.

`[terverifikasi]` Penanda kedua: **prefix ID case**. `"EDM-"` 262×, `"NB-"` 127×, `"RNW-"` 43×
di korpus facultative, diperiksa antara lain oleh
`RNW Fac In/Activity/ProtectRenewal_Act.xml` (`@contains(pyWorkPage.pxInsName,"NB-")`).

`[terverifikasi]` Penanda ketiga: **mesin bersama**. `Activity/SetOldData.xml`
(450.829 byte, 59 langkah) **identik di ketiga modul** (`df291257aa`) — mesin salinan data lama
dipakai ketiga siklus, bukan khas renewal maupun endorsement. Hal yang sama untuk
`CountPremiNusareRetro_Act` (`8f1268a660`) dan ketiga rule guard identitas (`IsGroup` `2cff7b9a67`).

`[terverifikasi]` Penanda keempat: **modul "new business" membawa 35 activity ber-nama EDM** dan
`IsEdm*` ×15 — kode siklus lain ikut terbawa di setiap ekspor.

**Kesimpulan berbukti:** ketiga modul menjalankan **satu basis rule** yang membedakan tahap siklus
**saat runtime** lewat `Quotation.StatusBusiness` dan prefix ID case. Pemisahan menjadi tiga folder
ekspor adalah **pemotongan administratif**, bukan pemisahan aplikasi. Yang berbeda per siklus adalah
**subhimpunan rule yang aktif**, bukan basis kodenya.

**Arti nilai 1/2/3 tetap belum terverifikasi** (OQ-020) — korelasi nama rule bukan bukti.

### 7.3 Kontras dengan domain Claim (Tahap 3)

`[terverifikasi]`

| | Facultative (NB / RNW / EDM) | Claim (Life / Prop / Non Prop / Fac In) |
| --- | --- | --- |
| Tumpang tindih identitas | **99,0%** (RNW ⊂ NB) | 12–40% |
| Rule pembeda siklus | **ada, terbaca** (`IsNB`/`IsRenewal`/`IsEDM`) | **tidak ada** |
| Nama Flow sama, class sama | ya → satu rule, atau varian isi | **tidak** — nama sama, **class berbeda** → rule berbeda |
| Bentuk graf | satu pola, panjang berbeda | **empat rangkaian tahap yang berbeda-beda** |
| Routing assignment | **seluruhnya `WorkBasket`/`Custom`** (kecuali `OfferFacOut`) | **campuran** `WorkList` + `WorkBasket` |
| Mekanisme limit | dari database (12 RDB-List per lini bisnis) | **campur**: hardcode (Non Prop) dan database (Prop) |

**Facultative = satu proses bercabang. Claim = empat proses berbeda.** Keduanya ekstrem yang
berlawanan, dan keduanya bukti untuk D4 — **bukan keputusan D4**.

---

## 8. Register OQ-011 diukur ulang untuk domain ini

`[terverifikasi]` **438 dari 533 identitas berkonflik korpus (82,2%) menyentuh ketiga modul
facultative sekaligus.**

```
awk -F'|' '/^\| [0-9]+ \|/ { if ($0 ~ /NB FacIn\// && $0 ~ /RNW Fac In\// && $0 ~ /Endorsment Fac In\//) n++ } END{print n}' \
  ../inventory/_oq011-konflik-isi.md
```

Sebaran per tipe: `RULE-OBJ-ACTIVITY` 171, `RULE-HTML-SECTION` 78, `RULE-OBJ-FLOWACTION` 42,
`RULE-CONNECT-SQL` 29, `RULE-OBJ-MODEL` 27, `RULE-DECLARE-PAGES` 27, `RULE-OBJ-WHEN` 24,
`RULE-OBJ-REPORT-DEFINITION` 21, `RULE-HTML-HARNESS` 8, `RULE-DECLARE-DECISIONTABLE` 6,
`RULE-CONNECT-REST` 3, `RULE-OBJ-FLOW` 1, `RULE-ADMIN-SYSTEM-SETTINGS` 1.

**Pengukuran baru `[terverifikasi]`** — setelah §6.2 menemukan tiga tag metadata yang luput dari
daftar 18 tag, seluruh 438 identitas dihitung ulang dengan normalisasi **21 tag**
(+ `pyDelete`, `pyVersionSecure`, `pzIsPrivateCheckOut`):

| Ukuran | Hasil |
| --- | ---: |
| Identitas berkonflik menyentuh ketiga modul | **438** |
| Masih >1 isi — normalisasi 18 tag | **438** |
| Masih >1 isi — normalisasi **21 tag** | **429** |
| **Konflik semu** (runtuh jadi identik) | **9** — 2,1% |

Kesembilannya: #12 `ASM!GETDATAKLAIM_SQL`, #13 `ASM!GETEDMOLDIDPEGA`,
**#324 `LINKSERVICE!LINKSERVICE`**, #367 `RNM!DELETESTORAGE_SQL`, #369 `RNM!GETAPPNAME_SQL`,
#371 `RNM!GETLINKSTORAGE_SQL`, #372 `RNM!GETMKTANDLEADER_SQL`, #376 `RNM!INSERT_T_STORAGE_SQL`,
#377 `RNM!UPDATE_T_STORAGE_SQL`.

**Register OQ-011 tetap berlaku: 97,9% konflik facultative adalah perbedaan isi yang nyata.**
Yang berubah secara material hanya **#324 `LINKSERVICE`** (§6.1) dan **#415
`SERVICEINSERTARASAPAS_ACT`** (§6.2, di luar 438 karena tidak menyentuh ketiganya).

---

## 9. Batas pengetahuan domain facultative — ringkas

| # | Batas | Sifat | OQ |
| ---: | --- | --- | --- |
| 1 | Baris tabel keputusan (49 `DecisionTable`, termasuk `IsUWAccepted`) **tidak ada di ekspor** | korpus | **OQ-043** |
| 2 | Rumus spreading/capacity/scoring di 34–38 activity, terbesar 996 KB | **cakupan telusur** | — |
| 3 | Tabel rate & kurs lewat `@ASMD.SINARMAS.CO.ID` (7–8 objek) | korpus | OQ-017 |
| 4 | 26–29 stored procedure, termasuk `GENERATE_FACRETRO_NO`, `CEK_PRORATA_TANGGAL`, `PEGA_FIRE_SET_RATE` | korpus | OQ-002 |
| 5 | Struktur `JSON_POLIS` / `JSON_OFFER` (17–18 rule perujuk); `@ASM.GetPageJSONString()` tanpa source | korpus | OQ-012 |
| 6 | Daftar endpoint di tabel `M_LINK_SERVICE` | korpus | **OQ-047** |
| 7 | Implementasi `SERVICEINSERTARASAPAS_ACT` tidak ada di 3 modul yang memanggilnya | korpus | OQ-025 |
| 8 | Kondisi `IsPKSASM`, `LetterNoNull`, `ToJUW_A`, `ToUW`, `IsEdmInternalRetro` tidak terbaca | korpus | OQ-029 |
| 9 | Pemetaan shape → nama workbasket (`<pyWorkBasket>` kosong, routing `Custom`) | korpus | OQ-024 |
| 10 | 16 shape yatim di `InputInwardFacultativeOffer`, 4 di `InputAddendumFacultativeIn` | korpus | OQ-023 |
| 11 | Isi `GeminiAIGoogle_Act` | **cakupan telusur** | OQ-048 |
| 12 | Isi 8 activity data-lama (`SetValueToEDMWork` 446 KB, `SetOLDValueToEDMWork_*` ×7) | **cakupan telusur** | — |

---

## 10. Pertanyaan terbuka baru dari Tahap 4

| OQ | Ringkas | Pemilik |
| --- | --- | --- |
| **OQ-043** | Baris tabel keputusan tidak ikut terekspor — 49 `DecisionTable`; `IsUWAccepted` menggerbangi 22 Decision di flow terbesar korpus | pemilik export + Product+UW |
| **OQ-044** | `ProposalAcceptStatus = 4` ditulis rule "banding" tetapi diuji rule "facout" | Product+UW |
| **OQ-045** | `pyWorkPage.LetterNo` dipakai sebagai token routing persetujuan (10 kode peran literal) | IAM + Product+UW |
| **OQ-046** | Pangsa retro 0,45/0,05 dan ambang 3 miliar ter-hardcode; mata uang tidak disebut; ambang dibandingkan sebagai string | Product+UW + Actuarial |
| **OQ-047** | Daftar endpoint ada di tabel Oracle `M_LINK_SERVICE` (`KATEGORI_1`/`KATEGORI_2`), bukan hanya SystemSettings | DBA + Platform |
| **OQ-048** | `GeminiAIGoogle_Act` memanggil model AI pihak ketiga di alur underwriting (2 dari 3 modul) | Product+UW + Security |
| **OQ-049** | Empat class bersufiks `ENDORSEMENT`, salah satunya berprefix `ASM-SFAGIS-` | Arsitektur Pega + pemilik export |

**Diperbarui secara material:** OQ-011 (§8), OQ-015 (§7.2 — **bukti mekanisme**), OQ-018 (§6.1),
OQ-025 (§6.2 — naik ke `[terverifikasi]`, risiko turun), OQ-026 (dipersempit ke `DecisionTable`).

**Dikuatkan:** OQ-002, OQ-008, OQ-012, OQ-014, OQ-016, OQ-017, OQ-020, OQ-021, OQ-023, OQ-024,
OQ-027, OQ-028, OQ-029, OQ-038, OQ-040.

---

## 11. Yang **tidak** diputuskan di sini

Sesuai `_METHOD.md` §0, dokumen ini **tidak** menetapkan:

- apakah ketiga modul menjadi satu bounded context atau tiga (→ **D4**);
- apakah `InputQuotation` menjadikan facultative dan treaty inward satu konteks (→ **D4**);
- arti nilai `ProposalAcceptStatus`, `EdmType`, `QuotationData.Type`, `StatusBusiness`,
  `ProRateType` (→ OQ terbuka, pemilik bisnis);
- apakah pola apa pun di sini baik atau buruk (→ **FASE B**, setelah gate).
