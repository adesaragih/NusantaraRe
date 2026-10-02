# Context Map — Nusantara Re (FASE A, STEP D4)

Disusun 2026-09-14 dari **20 inventaris modul (D1)** + **20 telusur konteks / 428 rule (D2)** +
**20 catatan modul (D3)**. Korpus: 9.369 file, READ-ONLY.

**Batas dokumen ini.** Ia **memetakan bounded context yang terbaca dari korpus dan menandai
kualitas buktinya**. Ia **tidak** mengusulkan paket Go, tidak merancang skema, tidak menilai
benar/salah — itu FASE B / ADR, setelah GATE manusia.

Konteks yang **tidak dapat diturunkan dari korpus** ditandai **absent** dengan pertanyaan
terbukanya, bukan diisi tebakan.

---

## 1. Ringkas: sembilan konteks + satu absent

| # | Bounded context | Modul | File | Cakupan bukti |
| ---: | --- | --- | ---: | --- |
| 1 | **Facultative Inward** | NB FacIn, RNW Fac In, Endorsment Fac In | 6.071 | **partial** |
| 2 | **Treaty Inward — Master & Akseptasi** | Treaty In, Treaty In Adjustment | 708 | **partial** |
| 3 | **Treaty Inward — Realisasi & Endorsement** | NB Treaty In, EDM Treaty In | 441 | **full** |
| 4 | **Treaty Arrangement (master term)** | Treaty Contract Out | 303 | **partial** |
| 5 | **Claim — Non-Life** | Claim Fac In, Claim Prop, Claim Non Prop | 1.031 | **partial** |
| 6 | **Claim — Life** | Claim Life | 136 | **full** |
| 7 | **Komite (tangga persetujuan klaim)** | Komite Claim FacIn, Prop, Non Prop, Life | 300 | **full** |
| 8 | **Life — Penawaran & Premium List** | PremiumList Life, Endorsement Life | 199 | **full** |
| 9 | **Life — Master** | Master Product Name Life, Master Contract Retro Life | 180 | **partial** |
| — | **Identity & Access (RBAC)** | — | **0** | **ABSENT** |

Total 9.369 file. Audit:
```
for m in <20 modul>; do find "$m" -name '*.xml' | wc -l; done
```

**Arti label cakupan** (ditetapkan di sini agar dapat diaudit):

| Label | Kriteria |
| --- | --- |
| **full** | titik masuk terbaca, alur end-to-end terekam, dan tidak ada activity/rumus besar yang menentukan perilaku inti masih tertutup |
| **partial** | alur terekam, tetapi **perilaku inti tertentu berada di dalam artefak yang belum/ tidak dapat dibaca** (rumus di activity besar, stored procedure, tabel keputusan tak terekspor) |
| **absent** | tidak ada rule di korpus yang mewujudkan konteks itu |

---

## 2. Konteks satu per satu

### 2.1 Facultative Inward — **partial**

**Modul:** `NB FacIn` (2.083), `RNW Fac In` (1.927), `Endorsment Fac In` (2.061) — **6.071 file,
64,8 % korpus**.

`[terverifikasi]` **Ini SATU ruleset dengan diskriminator siklus, bukan tiga basis kode.**
Empat bukti independen:

| Bukti | Ukuran |
| --- | --- |
| Rule pembeda siklus ada di **ketiga** modul | `When/IsNB.xml` (`Quotation.StatusBusiness = 1`), `When/IsRenewal.xml` (`= 2`), `When/IsEDM.xml` (`= 3`), `When/IsNotEDM.xml` (`!= 3`) — **12 file** |
| Tumpang tindih identitas | **1.906 dari 1.926 identitas `RNW Fac In` (99,0 %)** juga ada di `NB FacIn` |
| Mesin bersama ber-hash identik | `Activity/SetOldData.xml` (`df291257aa`), `Activity/CountPremiNusareRetro_Act.xml` (`8f1268a660`), `When/IsGroup.xml` (`2cff7b9a67`) |
| Kode siklus lain ikut terbawa | modul "new business" memuat 35 activity ber-nama EDM dan 15 rule `IsEdm*` |

Audit:
```
for w in IsNB IsRenewal IsEDM IsNotEDM; do
  for m in "NB FacIn" "RNW Fac In" "Endorsment Fac In"; do ls "$m/When/$w.xml"; done
done            # -> 12 file, semuanya ada
```

`[terverifikasi]` Penanda kedua: **prefix ID case** — `"EDM-"` 262×, `"NB-"` 127×, `"RNW-"` 43×;
diperiksa antara lain oleh `RNW Fac In/Activity/ProtectRenewal_Act.xml`
(`@contains(pyWorkPage.pxInsName,"NB-")`).

`[terverifikasi]` **Dua belas file `Flow` hanya memuat delapan isi berbeda** — empat dipakai ulang
lewat bukti hash, satu (`OfferFacRetro`) bercabang menjadi dua varian dengan **perbedaan perilaku
nyata**: tangga fac out 2 anak tangga (NB/RNW) vs **4** (Endorsment).

**Mengapa partial:** rumus **spreading / capacity / scoring** — 34–38 activity per modul, terbesar
`SumTSIPremiSpreadedRNM_FIRE_Act.xml` **996.052 byte** — **belum dibaca**; **baris tabel keputusan
`IsUWAccepted` tidak ikut terekspor** (OQ-043) sehingga pemetaan `ProposalAcceptStatus` → hasil
tidak dapat dinyatakan; input rating berada di **8 objek db-link di luar sistem** (OQ-017).

Rujukan: `flows/_SUMMARY-facultative.md`, `modules/NB FacIn.md`, `modules/RNW Fac In.md`,
`modules/Endorsment Fac In.md`.

### 2.2 Treaty Inward — Master & Akseptasi — **partial**

**Modul:** `Treaty In` (329), `Treaty In Adjustment` (379) — **708 file**. Keduanya **tanpa rule
`Flow`** (OQ-005).

`[terverifikasi]` **Berbagi basis kode secara masif**: dari 323 file bernama sama, **277 identik**
(85,8 %) setelah normalisasi 18 tag; 43 tetap berbeda setelah normalisasi 21 tag; 3 konflik semu.
D1 mengukurnya sebagai **320 identitas rule dibagi** — angka D2 (277 file identik dari 323) adalah
pengukuran ulang yang lebih ketat pada tingkat file.

`[terverifikasi]` **Mesin status adalah satu rule yang sama**: `DataTransform/Akseptasi_DT.xml`
(`DATA-PORTAL!AKSEPTASI_DT`, hash `58b8e650`) — 10 `WHEN`, 18 `OTHERWISE_WHEN`, 64 `SET`.
Tangga: `Admin` → `SecHead` → `DeptHead` → `Director` → `Resolve Complete`.

`[terverifikasi]` **Pembeda dua proses terbaca**: `Param.type` = `"revision"` / `"adjustment"` →
`TreatyIn.EDMState` = `"1"` / `"3"`; `RevisionState = "1"` **memendekkan tangga** (SecHead langsung
`Resolve Complete`).

`[terverifikasi]` **Pola sejajar dengan §2.1** — satu basis rule, dibedakan saat runtime. Di
facultative pembedanya `Quotation.StatusBusiness`; di sini `Param.type` + `EDMState`.

**Mengapa partial:** 31 activity limit/layer/spreading/ROL **belum dibaca** (terbesar
`TreatyInNPSetTotal.xml` 825.279 B); **4 dari 5 stored procedure adalah penulis master treaty**
(`PEGA_M_TREATY_IN*`) sehingga aturan simpan ada di database; `When/IsTreatyUser.xml`
(6 kondisi) **tidak terbaca**; selisih isi 43 rule **tidak dikarakterisasi**.

Rujukan: `flows/Treaty In.md`, `flows/Treaty In Adjustment.md`, `flows/_SUMMARY-treaty-master.md`.

### 2.3 Treaty Inward — Realisasi & Endorsement — **full**

**Modul:** `NB Treaty In` (278), `EDM Treaty In` (163) — **441 file**.

`[terverifikasi]` Keduanya ber-`Flow` dan ditelusur end-to-end:
`NB Treaty In/Flow/InputRealizationTreatyIn.xml` (`ASM-FW-GISFW-WORK!INPUTREALIZATIONTREATYIN`,
`Start1`) dan `EDM Treaty In/Flow/InputAddendumTreatyIn.xml`
(`ASM-FW-GISFW-WORK-ENDORSEMENTTREATY!INPUTADDENDUMTREATYIN`, `Start1`).

`[terverifikasi]` Realisasi dan endorsement bertulang sama; pembedanya FlowAction berakhiran
`Addendum` dan activity berakhiran `EDM`, serta `param.isFOR` = `"EDM"` / `"POLICY"` sebagai dua
mode simpan.

`[terverifikasi]` **Format nomor polis terbaca penuh** (`RDBList/GenerateNoPolicy.xml`).

**Mengapa full:** alur, guard, routing, objek Oracle, dan format nomor polis seluruhnya terekam.
Yang tersisa adalah **satu blocker yang sudah dipetakan** (OQ-025) dan batas korpus umum (SP, JSON)
— bukan perilaku inti yang belum dibaca.

### 2.4 Treaty Arrangement (master term) — **partial**

**Modul:** `Treaty Contract Out` (303). **Tanpa rule `Flow`** (OQ-005).

`[terverifikasi]` **Nama modul menyesatkan — ini BUKAN treaty outward.** Lima uji berbukti
(`flows/Treaty Contract Out.md` §2):

| Uji | Hasil |
| --- | --- |
| Objek database | **0** objek `*_OUT*`; yang ada `M_PROPORTIONALARRG` (10), `MTREATYSECURITY` (5), `TREATYREINSURER` (3), `TREATYBUSINESS` (3), `TREATYCONTRACT`, `M_TREATYYEAR`, `TREATYEXCHANGE` |
| Class rule | **0** class `*OUT*`; justru **5 rule berclass `ASM-FW-GISFW-INT-TREATY_IN`**, 11 file menyebutnya |
| Penamaan rule | 81 activity `*TreatyArr*` vs **4** `*Out*` — keempatnya soal **lampiran** |
| Arah tulis | **9 Connect-SQL MENULIS master** lewat `PEGA_TREATYCONTRACT`, `PEGA_TREATYYEAR`, `PEGA_TREATYREINSURER`, `PEGA_TREATYBUSINESS`, `PEGA_PROPORTIONALARRG`, `PEGA_M_PROPORTIONALARRG_CHILD`, `PROSESCOPY` |
| Pembanding korpus | penyentuh objek outward: `NB Treaty In` **8** file, `Claim Non Prop` **4**, `EDM Treaty In` **4**, `Treaty In Adjustment` 2, `Claim Prop` 1, **modul ini 0** |

**→ OQ-022 terjawab sebagian: BUKAN outward.** Ia **editor master *term / arrangement* kontrak
treaty** dengan **16 jenis klausul** (`CancelActivity*`: BordereAux, CashLossLimit,
ClaimCoorperation, EPI, ExGratia, FacIn, PLA, ProfitCommision, Ricomm, TreatyLimit, Portfolio,
TerrLimit, …).

**Mengapa partial:** **seluruh jalur tulis melewati stored procedure** (11 SP) sehingga aturan
penyimpanan master ada di database; **202 dari 303 rule (66,7 %) berada di `@BASECLASS`** sehingga
pemetaan rule → entitas tidak dapat diturunkan dari class (OQ-009); 168 Activity belum ditelusur.

**Yang masih terbuka:** *mengapa* dinamai "Out", dan ke konteks mana 303 rule ini ditempatkan.
**Itu keputusan manusia di GATE**, bukan di sini.

### 2.5 Claim — Non-Life — **partial**

**Modul:** `Claim Fac In` (482), `Claim Prop` (270), `Claim Non Prop` (279) — **1.031 file**.

`[terverifikasi]` **KEBALIKAN dari facultative: tiga basis kode terpisah, bukan satu ruleset.**

| Pasangan | Identitas dibagi | ~% |
| --- | ---: | ---: |
| Claim Fac In ↔ Claim Prop | 108 / 270 | 40 % |
| Claim Non Prop ↔ Claim Prop | 102 / 270 | 38 % |
| Claim Fac In ↔ Claim Non Prop | 83 / 279 | 30 % |

`[terverifikasi]` **Tidak ada rule pembeda siklus** seperti `IsNB`/`IsRenewal`/`IsEDM` di
facultative. Nama file `Flow` yang sama (`Flow_TreatyIn.xml`, `Register_Flow.xml`) berarti
**class berbeda → rule berbeda**, bukan satu rule bercabang.

`[terverifikasi]` Tahapan siklus **berbeda-beda per modul**: Fac In = (B2B?) → Register → Estimasi
→ Choose Surveyor; Prop & Non Prop = Outstanding → Input Acceptation. Routing **campuran**
(`WorkList` untuk input, `WorkBasket` untuk penilaian) — berbeda dari facultative yang seluruhnya
`WorkBasket`.

`[terverifikasi]` **Mekanisme limit tidak konsisten di dalam konteks yang sama**:
`Claim Non Prop` meng-hardcode (`LimitMax = 30000000.00`, `LimitMaxDivHead = 50000000.00`,
`LimitPersenMax = 30.00`) sementara `Claim Prop` mengambilnya dari database → OQ-040.

**Mengapa partial:** 20+ activity > 380 KB **belum dibaca** (terbesar `CLaimFaceSheet_Act.xml`
804.401 B); perhitungan alokasi kerugian XOL ada di dalamnya; **12 varian `ISCLM*` berkonflik
DAN kondisinya tidak terbaca** (OQ-041); 8 SP di Claim Prop saja.

### 2.6 Claim — Life — **full**

**Modul:** `Claim Life` (136).

`[terverifikasi]` **Dipisahkan dari §2.5 karena bukti memaksanya**: overlap identitas hanya
**12–14 %** terhadap ketiga modul Claim non-life; **tidak memakai `PaymentType` sama sekali**;
punya tahap **Medical Check** yang tidak ada di modul lain; seluruh objek Oracle-nya bersufiks
`_LIFE`.

`[terverifikasi]` Alur terekam penuh: Register → Outstanding → Medical Check → Claim Analis.

**Mengapa full:** alur, kode status (`STS_REJECT` 0/1/2), objek Oracle, dan titik temu dengan
Komite seluruhnya terekam; hanya **3 stored procedure** (paling sedikit di domain Claim) dan satu
activity besar (`SaveOutStandingLife_Act` ± 643 KB) yang tertutup — tidak menentukan alur inti.

### 2.7 Komite (tangga persetujuan klaim) — **full**

**Modul:** `Komite Claim FacIn` (114), `Komite Claim Prop` (80), `Komite Claim Non Prop` (59),
`Komite Claim Life` (47) — **300 file**.

`[terverifikasi]` **Tangga approval berbasis data, bukan berbasis struktur proses.** Keempat modul
punya graf yang **persis sama bentuknya** — 1 Assignment "KomiteRouter" + 1 Decision "KomiteLoop",
4 connector:

```
Start ─(Always)─> Assignment "KomiteRouter" [WorkList, Custom]
                      v
                  Decision "KomiteLoop"
                      ├─ When IsKomiteLoop ─> Assignment   ← LOOP
                      └─ Else ─────────────> End
```

dengan kondisi **sama teksnya di keempat modul**: `.AcceptStatus = "1"` **DAN**
`.KomiteCount <= .KomiteLoop`. Berapa tingkat dan siapa penyetujunya ditentukan **data**
(`KomiteList(KomiteCount)`), bukan shape.

`[terverifikasi]` **Bukan wrapper tipis di atas Claim**: 59 dari 133 identitas modul Komite
(**44,4 %**) tidak muncul di modul Claim mana pun; punya 4 class kerja sendiri, model data sendiri
(`KomiteList`), roster sendiri (`ASM-FW-GCNMFW-Int-EMAILKOMITE`), dan 4 rule penomoran akseptasi
berbeda.

`[terverifikasi]` **Empat implementasi berbeda untuk satu konsep** — properti batas tangga,
sasaran routing, guard identitas, dan rule penomoran semuanya berlainan antar modul.

**Mengapa full:** mekanisme tangga, routing, roster, otorisasi, dan penulisan tabel akseptasi
seluruhnya terekam, termasuk **55 nama kolom `OS_AKSEPTASI_KLAIM_LIFE`** yang terbaca langsung.

### 2.8 Life — Penawaran & Premium List — **full**

**Modul:** `PremiumList Life` (124), `Endorsement Life` (75) — **199 file**.

`[terverifikasi]` Keduanya ber-`Flow`, alur terekam penuh, dan **berbagi mesin yang sama**:
`Activity/InsertJsonPolisLife_Act.xml` dan `DecisionTable/IsLifeAccepted.xml` dipakai kedua modul;
`PremiumListSummary.PL_NUMBER_EDM` merujuk silang.

`[terverifikasi]` Keduanya memakai **`WorkList`** (bukan `WorkBasket`) → OQ-028.

**Mengapa full:** alur, gerbang, kode status connector (`Confirm`/`Decline`/`Reject`/`Offer`/
`Premium`), dan objek Oracle terekam. Yang tertutup adalah **isi tabel keputusan** (OQ-043) —
batas korpus yang berlaku universal, bukan celah telusur konteks ini.

### 2.9 Life — Master — **partial**

**Modul:** `Master Product Name Life` (114), `Master Contract Retro Life` (66) — **180 file**.
Keduanya **tanpa rule `Flow`** (OQ-005).

`[terverifikasi]` Keduanya editor master murni: nol rujukan `StatusAkseptasi`, tanpa `Akseptasi_DT`,
tanpa tombol Submit/Decline.

`[terverifikasi]` `Master Contract Retro Life` adalah **modul paling terisolasi di korpus**: tanpa
`ConnectREST`, tanpa `When`, tanpa `SystemSettings`, tanpa db-link, tanpa kolom JSON, tanpa jejak
guard identitas; hanya 4 objek Oracle, seluruhnya `POOLDATA.*_LIFE`.

**Mengapa partial:** **seluruh jalur tulis kedua modul melewati stored procedure** (5 + 5 SP,
sepuluhnya tanpa body di korpus); rumus `CountingPercentShare_Act` / `TreatyLimit_TypeProtect`
**belum dibaca**; **`REINSTYPEID` dan kode `OR` tidak pernah muncul sebagai nilai literal**
(OQ-057) sehingga enumerasi jenis reasuransi life tidak dapat dinyatakan.

### 2.10 Identity & Access (RBAC) — **ABSENT**

`[terverifikasi]` **Korpus tidak memuat satu pun rule identitas, peran, atau otorisasi** (OQ-007).
Tidak ada `Rule-Access-Role-Obj`, `Rule-Access-Deny`, `Data-Admin-Operator-*`, atau sejenisnya di
17 tipe rule yang diinventarisasi.

`[terverifikasi]` Sebagai gantinya, otorisasi tersebar sebagai **nilai data dan ekspresi
visibilitas** dalam **empat pola berbeda**:

| Pola | Mekanisme | Bukti | OQ |
| --- | --- | --- | --- |
| 1. Identitas orang sebagai **guard** | `OperatorID.pyUserIdentifier`, `pyPosition` diuji terhadap konstanta | `Komite Claim FacIn` 4 file, `Komite Claim Prop` 3 file | OQ-021 |
| 2. Field bisnis menyimpan **kode peran** | `OperatorID.pyTelephone` = `TREATY1`/`TREATY2`/`SPVTREATY1`/`SPVTREATY2`; `pyWorkPage.LetterNo` = 10 token (`KADIVTEKNIK`, `DIREKTURTEKNIK`, …) | `Treaty In/DataTransform/Akseptasi_DT.xml`; `NB FacIn/When/ToKadivTeknik.xml` dkk | OQ-027, **OQ-045** |
| 3. Otorisasi lewat **indeks tetap** | `OperatorID.pyWorkBasketList(2).pyWorkBasketName` — 36 kemunculan | `Treaty In/Section/TreatyInActionButtons.xml` | **OQ-051** |
| 4. Identitas orang **ditetapkan sebagai tujuan rute** | `TreatyIn.PositionUsername := "<orang>"` — **5 identitas ter-hardcode** | `Treaty In/DataTransform/Akseptasi_DT.xml` | **OQ-053** |

`[terverifikasi]` Satu pola **sudah berbasis data** dan dapat menjadi titik awal:
`Komite Claim Non Prop/Activity/KomitePostAdjustment.xml` mencocokkan pengguna dengan **roster**
(`@contains(…KomiteList(…).KomiteID, …OperatorID.pyUserIdentifier)`), dan roster itu berasal dari
tabel database (class `ASM-FW-GCNMFW-Int-EMAILKOMITE`).

`[pertanyaan terbuka]` Model peran yang sesungguhnya — siapa boleh melakukan apa — **harus
dirancang, bukan dimigrasikan**. Ini keputusan manusia (IAM + Product+UW), bukan temuan korpus.
→ OQ-007, OQ-021, OQ-024, OQ-027, OQ-028, OQ-036, OQ-045, OQ-051, OQ-053.

**Nilai nama orang tidak disalin** ke artefak D1–D4 mana pun.

---

## 3. Kebocoran batas antar-konteks — berbukti

`[terverifikasi]` Tujuh kebocoran terukur. "Kebocoran" di sini berarti **satu konteks membaca atau
menulis data milik konteks lain**, bukan penilaian.

| # | Dari | Ke | Bukti | OQ |
| ---: | --- | --- | --- | --- |
| 1 | **Claim Non Prop** | Treaty **outward** | `RDBList/GetDataMasterTOutNP.xml` → `M_TREATY_OUT`; `BrowseDtlTreatyOutNP.xml` → `M_TREATY_OUT_DETAIL`; `GetLimitTONPPLA.xml` → `TREATY_OUT` (4 file) | **OQ-042** |
| 2 | **Claim Fac In** | Treaty inward (master) | membaca `TREATYCONTRACT`, `TREATYBUSINESS`, `PROPORTIONALARRG`; `Activity/DLAFacintoTreaty_Act.xml` (644.505 B) | OQ-042 |
| 3 | **Master Product Name Life** | Treaty inward (master) | membawa `RDBList/SaveTreatyIn.xml` (`ASM!SAVETREATYIN`) → procedure `POOLDATA.PEGA_TREATY_IN`; **penulis `M_TREATY_IN` kedua** | **OQ-056** |
| 4 | **NB FacIn** | Treaty inward (realisasi) | `Flow/InputQuotation.xml` bercabang ke `InputRealizationTreatyIn`; `Flow` itu **identik** dgn NB Treaty In (`1306f68d56`) | — |
| 5 | **Treaty In / Treaty In Adjustment** | Claim & Komite | `POOLDATA.OS_AKSEPTASI_KLAIM` (2 rule masing-masing) — tabel akseptasi klaim | — |
| 6 | **Komite Claim Prop / Non Prop** | Claim | membawa rule berclass `…WORK-CLAIMTREATY` (12) / `…WORK-CLAIMTREATYNONPROP` (5) | — |
| 7 | **Treaty Contract Out** | Treaty inward | 5 rule berclass `ASM-FW-GISFW-INT-TREATY_IN`; 12 rule bernama FacIn/TreatyIn | **OQ-022** |

### 3.1 Titik temu data yang dipakai bersama — bukan kebocoran

`[terverifikasi]` Tiga rule terbukti **satu identitas, satu isi, dipanggil dari dua sisi** —
ini berbagi rule, bukan duplikasi:

| Rule | Hash | Dipakai oleh |
| --- | --- | --- |
| `ASM-FW-GISFW-INT-TREATY_IN!ASM!SAVETREATYIN` | `e371c194` | **6 modul**: Treaty In, Treaty In Adjustment, NB Treaty In, EDM Treaty In, NB FacIn, Master Product Name Life |
| `ASM-FW-GISFW-INT-TREATY_IN!ASM!BROWSETREATYIN` | `196d49b5` | 6 modul yang sama |
| `UpdateOsAkseptasiClaimLife_sql` | `c50bfd9a12` | Claim Life **dan** Komite Claim Life |

Ketiganya **tidak terdaftar di register OQ-011** — jadi bukan varian yang menyimpang.

### 3.2 Satu ketergantungan yang terhenti — blocker

`[terverifikasi]` `ASM-FW-GISFW-WORK!SERVICEINSERTARASAPAS_ACT` (232 KB) **dipanggil** oleh
`NB Treaty In`, `NB FacIn`, dan `RNW Fac In` — tetapi **hanya ada di** `EDM Treaty In`
(`ab3ae3c9ba`) dan `Endorsment Fac In` (`7314b6c49c`), dan keduanya terdaftar berkonflik.

Telusur **dihentikan** di ketiga modul pemanggil; varian modul lain **tidak dipinjam** → **OQ-025**.

`[terverifikasi]` Risikonya lebih kecil dari perkiraan D1: kedua varian berbeda **hanya pada 3 tag
metadata Pega**; 20 langkah dan rule yang dirujuk identik. Tetapi blocker prosedural tetap berdiri.

---

## 4. Dua arsitektur modul yang berlawanan — temuan struktural utama

`[terverifikasi]` Korpus memuat **dua pola pengorganisasian yang berlawanan**, dan keduanya terukur:

| | **Satu ruleset, bercabang saat runtime** | **Modul benar-benar terpisah** |
| --- | --- | --- |
| Contoh | Facultative NB/RNW/EDM; Treaty In ↔ Adjustment | Trio Claim non-life |
| Tumpang tindih identitas | **99,0 %** (RNW ⊂ NB); **85,8 %** file identik (Treaty In ↔ Adj) | **12–40 %** |
| Rule pembeda siklus | **ada, terbaca** (`IsNB`/`IsRenewal`/`IsEDM`; `Param.type`+`EDMState`) | **tidak ada** |
| Nama `Flow` sama berarti | satu rule, atau varian isi yang terukur | **class berbeda → rule berbeda** |
| Routing assignment | seluruhnya `WorkBasket`/`Custom` (facultative) | **campuran** `WorkList` + `WorkBasket` |

`[terverifikasi]` Domain **Komite** adalah pola ketiga: **empat rule berbeda dengan graf yang
bentuknya identik** — sama bentuk, beda implementasi.

**Konsekuensi untuk penetapan konteks adalah keputusan manusia di GATE**, bukan simpulan dokumen
ini.

---

## 5. Apa yang harus diputuskan manusia sebelum context map ini final

`[terverifikasi]` Lima keputusan menentukan bentuk akhir peta, dan **tidak satu pun dapat
diturunkan dari korpus**:

| # | Keputusan | Bukti yang tersedia | OQ |
| ---: | --- | --- | --- |
| 1 | Facultative NB/RNW/EDM → **satu** konteks atau **tiga**? | satu ruleset, 99,0 % overlap, diskriminator terbaca | **OQ-015** |
| 2 | `Treaty In` + `Adjustment` → **satu** konteks atau **dua**? | 85,8 % file identik, mesin status satu rule | **OQ-010** |
| 3 | Ke mana **`Treaty Contract Out`** (303 rule) ditempatkan? | master arrangement inward, bukan outward | **OQ-022** |
| 4 | Trio Claim non-life → **tiga** konteks atau satu dengan varian? | overlap 12–40 %, tanpa diskriminator runtime | **OQ-019** |
| 5 | Siapa **pemilik** data lintas konteks? | Claim membaca master Treaty; Master Life menulis `M_TREATY_IN` | **OQ-042**, **OQ-056** |

Ditambah dua yang menentukan isi, bukan bentuk:

| # | Keputusan | OQ |
| ---: | --- | --- |
| 6 | Model RBAC — **harus dirancang**, tidak ada di korpus | OQ-007 dan keluarganya |
| 7 | Versi rule mana yang berlaku di production (533 identitas berkonflik) | **OQ-011** |

---

## 6. Yang **tidak** ditetapkan di dokumen ini

Sesuai batas fase:

- **tidak** mengusulkan paket/modul Go, struktur direktori, atau pemisahan service;
- **tidak** merancang skema database atau kontrak API;
- **tidak** memutuskan kelima pertanyaan di §5 — bukti disajikan, keputusan milik manusia;
- **tidak** menilai pola mana pun baik atau buruk.

Seluruhnya adalah wilayah **FASE B / ADR**, setelah GATE.
