# STEP D2 — LAPORAN PENUTUP

FASE A Discovery, migrasi Pega → Go + React (Nusantara Re).
Ditutup **2026-09-13**. Korpus READ-ONLY `D:\XML\RNM_BRD\` — **9.369 file**, 20 modul.

D2 **merekam perilaku existing apa adanya**. Dokumen ini **tidak** menetapkan bounded context,
tidak menilai, dan tidak merancang Go/React/skema — itu D4 dan FASE B, setelah gate manusia.

---

## 1. Cakupan yang selesai

`[terverifikasi]` **20 konteks ditelusur, 428 rule**, di 5 tahap.

| Tahap | Domain | Konteks | Sintesis |
| ---: | --- | ---: | --- |
| 1 | treaty/life offer | 4 | — |
| 2 | Komite / tangga persetujuan | 4 | `flows/_SUMMARY-komite.md` |
| 3 | Claim / siklus kerugian | 4 | `flows/_SUMMARY-claim.md` |
| 4 | Facultative Inward | 3 | `flows/_SUMMARY-facultative.md` |
| 5 | Treaty & Master (tanpa `Flow`) | 5 | `flows/_SUMMARY-treaty-master.md` |

Konvensi telusur: **`flows/_METHOD.md`** (semua tahap) + **`flows/_METHOD-noflow.md`** (Tahap 5).

Perintah audit:
```
cd discovery/flows
for f in *.md; do case "$f" in _*) continue;; esac
  grep -ohE '\*\*[0-9]+ rule ditelusur' "$f" | head -1 | grep -oE '[0-9]+'
done | awk '{s+=$1;n++} END{print n" konteks, "s" rule"}'
# -> 20 konteks, 428 rule
```

### 1.1 Daftar lengkap konteks + titik masuk

| # | Konteks | Titik masuk | Rule |
| ---: | --- | --- | ---: |
| 1 | `NB Treaty In` | `Flow/InputRealizationTreatyIn.xml` → `Start1` | 15 |
| 2 | `EDM Treaty In` | `Flow/InputAddendumTreatyIn.xml` → `Start1` | 20 |
| 3 | `Endorsement Life` | `Flow/InputEDMLife.xml` → `Start1` | 14 |
| 4 | `PremiumList Life` | `InputPolicyHolder.xml` (root modul) → `Start1` | 13 |
| 5 | `Komite Claim Life` | `Flow/KomiteLife_Flow.xml` → `Start2` | 12 |
| 6 | `Komite Claim FacIn` | `Flow/Komite_Flow.xml` → `Start1` | 15 |
| 7 | `Komite Claim Prop` | `Flow/KomiteTreaty_Flow.xml` → `Start1` (class `…KOMITETREATY`) | 13 |
| 8 | `Komite Claim Non Prop` | `Flow/KomiteTreaty_Flow.xml` → `Start1` (class `…KOMITETREATYNONPROP`) | 13 |
| 9 | `Claim Life` | `Flow/Register_Flow.xml` → `Start2` (class `…CLAIMLIFE`) | 16 |
| 10 | `Claim Prop` | `Flow/Flow_TreatyIn.xml` → `Start1` (class `…CLAIMTREATY`) | 14 |
| 11 | `Claim Non Prop` | `Flow/Flow_TreatyIn.xml` → `Start1` (class `…CLAIMTREATYNONPROP`) | 15 |
| 12 | `Claim Fac In` | `Flow/Register_Flow.xml` → `Start1` (class `…WORK-PNC`) | 17 |
| 13 | `NB FacIn` | 6 `Flow`; utama `Flow/InputInwardFacultativeOffer.xml` → `Start1` | 45 |
| 14 | `RNW Fac In` | `Flow/InputRenewalFacultativeIn.xml` → `Start2` | 28 |
| 15 | `Endorsment Fac In` | `Flow/InputAddendumFacultativeIn.xml` → `Start2` | 52 |
| 16 | `Treaty In` | **Harness** `InputTreatyInOffer.xml` (`DATA-PORTAL`, 8,2 MB) | 24 |
| 17 | `Treaty In Adjustment` | **Harness** `InputTreatyInAdjustment.xml` (`DATA-PORTAL`) | 31 |
| 18 | `Treaty Contract Out` | **Harness** `InboxTreatyContract.xml` (`DATA-PORTAL`) | 26 |
| 19 | `Master Product Name Life` | **Harness** `InwardProductName.xml` (`…INT-PRODUCT_LIFE`) | 22 |
| 20 | `Master Contract Retro Life` | **Harness** `InboxRetroLifeReinsurersList.xml` (`DATA-PORTAL`) | 23 |

`[terverifikasi]` Konteks 1–15 masuk lewat rule `Flow` → `<pyStartActivity>`;
konteks 16–20 lewat **Harness/FlowAction** karena kelima modul itu tidak punya rule `Flow`
(OQ-005, kini terjawab untuk keperluan D2).

### 1.2 Rule `Flow` yang tidak ditelusur ulang karena terbukti identik

`[terverifikasi]` Empat file `Flow` dipakai ulang berdasarkan bukti hash ternormalisasi,
bukan asumsi:

| Hash | Rule | Dipakai ulang dari |
| --- | --- | --- |
| `1306f68d56` | `InputRealizationTreatyIn` (NB FacIn = NB Treaty In) | `flows/NB Treaty In.md` |
| `c099bf4ebb` | `InputInwardFacultativeRISlip` (NB FacIn = RNW Fac In) | `flows/NB FacIn.md` §1.3 |
| `832fb12b9d` | `OfferFacOut` (NB FacIn = RNW Fac In) | `flows/NB FacIn.md` §1.4 |
| `bfd6252070` | `OfferFacRetro` (NB FacIn = RNW Fac In) | `flows/NB FacIn.md` §1.5 |

Sebaliknya, **enam rule ditelusur terpisah** karena terbukti berbeda: `OfferFacRetro` varian
Endorsment (`b370146c63`), dan dua pasang Flow yang berbagi nama file tetapi berbeda class
(`Register_Flow`, `Flow_TreatyIn`, `KomiteTreaty_Flow`).

---

## 2. Peta integrasi eksternal — seluruhnya konfigurasi

`[terverifikasi]` Korpus memuat **51 file `RULE-CONNECT-REST`** dengan **15 nama layanan berbeda**.

```
find . -path ./OUTPUT_HASIL_RNM -prune -o -type f -path "*/ConnectREST/*.xml" -print | wc -l   # 51
… | while read f; do grep -o "<pyServiceName>[^<]*" "$f" | head -1; done | sort | uniq -c | sort -rn
```

| `pyServiceName` | File | Domain |
| --- | ---: | --- |
| `ServiceGoogle` | **15** | seluruh domain |
| `convertJsonNusareToProduction` | 6 | facultative, treaty |
| `SendAcceptationToKasir` | 6 | claim |
| `KonversiKlaimNonLife` | 6 | claim |
| `getPremiumPaidOn` | 4 | facultative |
| `getPayAttachment` | 3 | claim |
| `insertClaimFinalOrClosed_NP`, `getPremiumPaidOnTreatyIn` | 2 masing-masing | claim, treaty |
| `insertClaimReject_NP`, `getPremiumPaidOnMarine`, `getPaymentClaim`, `convertJsonNusareToProductionClaimLife`, `InsertClaimOutstanding_NP`, `HitDLAClaimFacin`, `GetDtlPaymentClaim` | 1 masing-masing | claim, facultative |

### 2.1 Cara alamat ditentukan — **dua lapis, keduanya di luar rule**

`[terverifikasi]`

1. **`RULE-ADMIN-SYSTEM-SETTINGS` `LinkService!LinkService`** — `pyBaseURLSelectionType = SETTING`
   pada hampir seluruh `ConnectREST`. Rule ini **identik** di modul-modul yang membawanya setelah
   normalisasi 21 tag (konflik semu OQ-011 #324), jadi konfigurasi base URL **tidak bercabang**.
2. **Tabel Oracle `M_LINK_SERVICE`** — `Activity/GetLinkService.xml`
   (`ASM-FW-GISFW-INT-M_LINK_SERVICE!GETLINKSERVICE`) melakukan `Obj-Browse` dengan kunci
   `KATEGORI_1` / `KATEGORI_2`, hasilnya dipakai langkah `Connect-REST` berikutnya. **Isi tabel
   tidak ada di korpus** → **OQ-047**.

`[terverifikasi]` **Satu-satunya URL literal endpoint** di seluruh korpus ada di
`Claim Fac In/ConnectREST/getPremiumPaidOnMarine.xml` (`pyBaseURLSelectionType = URL`) — dan
nilainya **tidak disalin** ke artefak D2 karena memuat contoh data (OQ-018). Korpus juga memuat
hostname DEV `appdev.nusantarare.com` (7×) di 81 file/15 modul.

### 2.2 Integrasi per nama, sebagai konfigurasi

| Integrasi | Cara tersentuh | Status telusur |
| --- | --- | --- |
| **Arasapas** | shape Utility "HIT SERVICE ARASAPAS" → `serviceInsertArasapas_act` / `…RNW_act` / `…EDM_act`; skema Oracle `ARASAPAS` | **20 langkah terbaca** — hanya varian `Endorsment Fac In`; **terblokir** di NB Treaty In / NB FacIn / RNW Fac In (**OQ-025**) |
| **Kasir** | `ConnectREST/SendAcceptationToKasir` (6 file), tabel `POOLDATA.DIRECTTOKASIR_LOG`, `Activity/HitServiceToKasir_Act.xml` (482 KB) | isi activity **belum dibaca** |
| **Konversi** | `ConnectREST/KonversiKlaimNonLife` (6), `convertJsonNusareToProduction` (6); properti `.StatusService.StsKonversiFacIn/FacOut`; `UpdateStsKonversiFacOut_Act` | gerbang `IsSuccessHitService` terbaca; isi activity **belum dibaca** |
| **Google Storage** | `ConnectREST/ServiceGoogle` (15), `InsertGoogleStorage_Act`, `GetUrlGoogleStorage_Act`, `DeleteGoogleStorage_Act`, tabel `T_STORAGE_IMAGE` / `POOLDATA.T_FOLDER_IMAGE`, procedure `POOLDATA.GET_TOKEN_STORAGE` | rantai terbaca; isi procedure **tidak ada di korpus** |
| **Gemini AI** | `Activity/GeminiAIGoogle_Act.xml` — **hanya di `NB FacIn` dan `RNW Fac In`** | isi **belum dibaca** → **OQ-048** |
| **`M_LINK_SERVICE`** | `Activity/GetLinkService.xml`, dipanggil dari activity Arasapas & Google Storage | isi tabel **tidak ada di korpus** → **OQ-047** |
| **officeapps / dokumen** | tidak ditemukan sebagai `ConnectREST` maupun nama rule di korpus | **tidak ada bukti** — tidak dinyatakan |

`[terverifikasi]` Pencarian `officeapps` di seluruh korpus tidak menghasilkan rule apa pun;
istilah itu **tidak dapat dikonfirmasi dari korpus** dan karena itu **tidak dicatat sebagai
integrasi**.

---

## 3. Rekap batas pengetahuan — yang WAJIB dijawab manusia sebelum FASE B

Dibedakan tegas antara **batas korpus** (informasinya memang tidak ada) dan **batas cakupan
telusur** (ada di korpus, belum dibaca karena volume).

### 3.1 Batas korpus — tidak dapat dijawab dengan membaca lebih banyak

| # | Batas | Ukuran terverifikasi | OQ |
| ---: | --- | --- | --- |
| 1 | **Stored procedure tanpa body** | **70** procedure dipanggil (67 kustom + 3 bawaan Oracle `DBMS_LOB`, `DBMS_OUTPUT`, `UTL_MATCH`), tersebar di **10 skema aplikasi**: `POOLDATA` (61), `FIRE` (2), `GENERAL`, `GL`, `MBU`, `NEW_GENERAL`, `NEW_UNDERWRITING`, `ARASAPAS`, `DATAPEGA`, `REINSURANCE` | **OQ-002**, OQ-013, OQ-016 |
| 2 | **Tabel lewat database link** | **8** objek via `@ASMD.SINARMAS.CO.ID`: `M_EQS_RATE`, `M_TERORISME_RATE`, `M_RSMD_RATE`, `M_FLEXAS_RATE`, `M_FLOOD_AREA`, `FIRE.M_FLOOD_RATE`, `FIRE.M_BI_INDEMNITY`, `LST_KURS_STANDARD` — **hanya di 3 modul facultative** | **OQ-017** |
| 3 | **Kolom JSON** | **10** objek berlabel JSON: `JSON_POLIS`, `JSON_OFFER`, `JSON_KLAIM`, `JSON_FOLLOWING`, `JSON_POLIS_ERROR`, `POOLDATA.JSON_OFFER_LIFE`, `POOLDATA.JSON_POLIS_MONITORING`, dll. Ditambah kolom `JSONDATA` / `DATA_JSON`. Pembentuknya fungsi kustom **`@ASM.GetPageJSONString()`** yang **source-nya tidak ada di korpus** | **OQ-012** |
| 4 | **Tabel keputusan tidak terekspor** | **49 file `DecisionTable`**, **tak satu pun memuat baris keputusannya**. Termasuk `IsUWAccepted` yang menggerbangi **22 shape Decision** di flow terbesar korpus | **OQ-043** |
| 5 | **Tidak ada DDL / definisi properti** | nol file DDL; tipe kolom tidak dapat disimpulkan dari SQL | **OQ-001** |
| 6 | **Rule dengan kondisi tak terbaca** | 5 keluarga: `IsPEGAPROD`, `IsSendtoMedical`/`IsSPK`, `ISCLM*` (12 varian), `IsPKSASM`/`LetterNoNull`/`ToJUW_A`/`ToUW`/`IsEdmInternalRetro`, `IsTreatyUser` (6 kondisi) / `recordEvent` | **OQ-029**, OQ-041 |
| 7 | **Identitas rule berkonflik** | **533** identitas dengan >1 isi; 438 (82,2 %) menyentuh ketiga modul facultative; **429 tetap berbeda** setelah normalisasi 21 tag | **OQ-011** |
| 8 | **Arti kode & enumerasi** | `ProposalAcceptStatus` (1/2/3/4/7/9), `StatusAkseptasi` (4 nilai), `PaymentType` (1–7), `EdmType` (1–4), `QuotationData.Type` (12 nilai), `StatusBusiness` (1/2/3), `ProRateType`, `REINSTYPEID` (**tanpa literal**), kode `OR` (**tanpa literal**), `L1`–`L11` | **OQ-020**, OQ-038, OQ-057 |
| 9 | **Slot parameter generik** | `CARI1`…`CARI30` dipakai sebagai slot menuju SQL di beberapa modul | **OQ-059** |
| 10 | **Kelengkapan ekspor** | folder `excludeXML`; folder **`Claude outputs`** (`.xlsx`, `.diff`) di 4 modul; hostname DEV; 2 `pxHostId` berbeda | OQ-003, **OQ-054**, **OQ-018** |

### 3.2 Batas cakupan telusur — ada di korpus, belum dibaca

| # | Yang belum dibaca | Ukuran |
| ---: | --- | --- |
| 1 | **Rumus spreading / capacity / scoring** facultative | 34–38 activity per modul; terbesar `SumTSIPremiSpreadedRNM_FIRE_Act.xml` **996.052 B** |
| 2 | **Rumus limit / layer / ROL** treaty inward | 31 activity; terbesar `TreatyInNPSetTotal.xml` **825.279 B** |
| 3 | **Activity besar domain Claim** | `CLaimFaceSheet_Act.xml` 804 KB, `CreateChildKomiteCNP_Act.xml` 757 KB, dan 20+ lainnya |
| 4 | **Mesin data lama endorsement** | `SetValueToEDMWork` 446 KB + 7 `SetOLDValueToEDMWork_*` |
| 5 | **Selisih isi rule yang tetap berbeda** | 43 (`Treaty In` ↔ `Adjustment`), 12 (`Master Product Name Life` ↔ `Treaty In`) |
| 6 | **Harness besar** | `InputTreatyInOffer.xml` **8.225.095 B** — hanya di-grep |
| 7 | **`GeminiAIGoogle_Act`**, `HitServiceToKasir_Act` | isi belum dibaca |

`[terverifikasi]` Perbedaan ini penting untuk perencanaan: **§3.1 memerlukan jawaban manusia atau
akses database**; **§3.2 hanya memerlukan waktu telusur tambahan** dan dapat dikerjakan tanpa
memblokir siapa pun.

---

## 4. Seluruh OQ terbuka, dikelompokkan per pemilik peran

`[terverifikasi]` **57 terbuka / 2 terjawab** (`open-questions.md`).

```
awk '/^## Terbuka$/{s="T"} /^## Terjawab$/{s="J"} /^### OQ-[0-9]+ /{c[s]++}
     END{print "terbuka: "c["T"]"  terjawab: "c["J"]}' open-questions.md
```

Kolom **Blokir**: `D3` = memblokir catatan modul / glossary final; `D4` = memblokir context-map /
understanding-report; `B` = memblokir FASE B (ADR/BRD/spec/tiket).

### 4.1 DBA

| OQ | Ringkas | Blokir |
| --- | --- | --- |
| **OQ-001** | Tidak ada DDL Oracle maupun definisi properti | **D4**, **B** |
| **OQ-002** | Body 67 stored procedure tidak ada di korpus | **D4**, **B** |
| **OQ-008** | Arti prefix `ASM` / `RNM` / `GCNM` pada kunci RDBList | D3 |
| **OQ-012** | Struktur JSON di `JSONDATA` / `DATA_JSON` (bersama Product+UW) | **D4**, **B** |
| **OQ-013** | Batas transaksi & identitas pengguna di dalam procedure | **B** |
| **OQ-016** | Peran 10 skema Oracle; mana yang dalam ruang lingkup (bersama Product+UW) | **D4** |
| **OQ-017** | Database link `ASMD.SINARMAS.CO.ID` — rating di luar sistem (bersama Actuarial, Product+UW) | **D4**, **B** |
| **OQ-047** | Daftar endpoint di tabel `M_LINK_SERVICE` (bersama Platform) | **B** |
| **OQ-055** | Nomor revisi di dalam string ID, ter-hardcode di rule dan SQL (bersama Product+UW) | **B** |
| **OQ-056** | Master produk life menulis master treaty inward (bersama Product+UW) | **D4** |
| **OQ-058** | Tabel internal Pega `DATAPEGA.PC_*` diakses langsung lewat SQL (bersama Arsitektur Pega) | **B** |
| **OQ-059** | Slot parameter generik `CARI<n>` menuju SQL (bersama Product+UW) | **D4**, **B** |

### 4.2 Product + Underwriting

| OQ | Ringkas | Blokir |
| --- | --- | --- |
| **OQ-003** | Folder `excludeXML` — aktif atau tidak | D3 |
| **OQ-009** | Rule di class bawaan Pega (`@BASECLASS`) — 202/303 di `Treaty Contract Out` | **D4** |
| **OQ-010** | `Treaty In` ↔ `Treaty In Adjustment`: dua modul atau satu | **D4** |
| **OQ-011** | 533 identitas berkonflik — versi mana yang berlaku di production | **D3**, **D4**, **B** |
| **OQ-014** | Arti singkatan `EDM` | D3 |
| **OQ-015** | Facultative NB/RNW/EDM — **terjawab sebagian**; arti `StatusBusiness` tersisa | **D4** |
| **OQ-019** | Trio Claim tidak berbagi basis kode — apa pembedanya | **D4** |
| **OQ-020** | Arti seluruh kode/enumerasi (bersama Finance) | **D3**, **B** |
| **OQ-022** | `Treaty Contract Out` — **terjawab: bukan outward**; alasan penamaan & konteks tersisa | **D4** |
| **OQ-023** | Shape/tombol tanpa jalur masuk — 16 di flow terbesar, 4 di endorsement, guard `FALSE &&` | **D4** |
| **OQ-025** | `SERVICEINSERTARASAPAS_ACT` — implementasi tidak ada di 3 modul yang memanggilnya | **B** |
| **OQ-026** | `isApproved` / `IsUWAccepted` ada sebagai dua tipe rule | **B** |
| **OQ-030**–**OQ-034** | Ambang hari-25, ID 1000032–35, `.KomiteLoop`, tanggal cutover `20250207` | D3, **B** |
| **OQ-036** | Empat sasaran routing komite ter-hardcode (bersama IAM) | **B** |
| **OQ-038** | Sebelas kode lini bisnis `L1`…`L11` | **D4** |
| **OQ-039** | Adjustment / Close / Reject tidak muncul di graf flow | **D4** |
| **OQ-041** | `ISCLM*` berkonflik 6 versi **dan** tidak terbaca | **D4**, **B** |
| **OQ-042** | Klaim membaca master treaty outward — recovery/retrosesi? | **D4** |
| **OQ-043** | Baris tabel keputusan tidak terekspor (bersama pemilik export) | **D4**, **B** |
| **OQ-044** | `ProposalAcceptStatus = 4` ambigu (banding vs facout) | **B** |
| **OQ-046** | Pangsa retro 0,45/0,05 & ambang 3 miliar (bersama Actuarial) | **B** |
| **OQ-048** | `GeminiAIGoogle_Act` — model AI di alur underwriting (bersama Security) | **D4**, **B** |
| **OQ-050** | 10 tombol `(dev)`, termasuk `Force Resolve Complete` (bersama pemilik export) | **D4**, **B** |
| **OQ-052** | `Reject` vs `Decline` — bedanya tidak dijelaskan | **B** |
| **OQ-057** | Kode `OR` & `REINSTYPEID` tidak pernah muncul sebagai nilai literal | **D3**, **B** |

### 4.3 IAM

| OQ | Ringkas | Blokir |
| --- | --- | --- |
| **OQ-007** | Tidak ada rule identitas/otorisasi di korpus | **D4**, **B** |
| **OQ-018** | Host/URL ter-hardcode termasuk hostname DEV (bersama DBA) | **B** |
| **OQ-021** | Identitas orang ter-hardcode sebagai guard (bersama Product+UW) | **B** |
| **OQ-024** | Pemetaan Assignment → workbasket tidak terbaca (bersama Product+UW) | **D4**, **B** |
| **OQ-027** | `OperatorID.pyTelephone` menyimpan kode peran — **4 nilai** | **B** |
| **OQ-028** | `WorkList` vs `WorkBasket`: dua model penugasan (bersama Product+UW) | **D4** |
| **OQ-029** | Kondisi `When` tidak terbaca — 5 keluarga (bersama DBA, Product+UW) | **D4**, **B** |
| **OQ-045** | `LetterNo` sebagai token routing persetujuan (bersama Product+UW) | **B** |
| **OQ-051** | Otorisasi lewat **indeks tetap** `pyWorkBasketList(2)` | **B** |
| **OQ-053** | **5 identitas orang ditetapkan sebagai pemilik tugas berikutnya** (bersama Product+UW) | **B** |

### 4.4 Actuarial

| OQ | Ringkas | Blokir |
| --- | --- | --- |
| **OQ-017** | Tabel rate & kurs di luar sistem (bersama DBA, Product+UW) | **D4**, **B** |
| **OQ-046** | Pangsa retro & ambang ter-hardcode (bersama Product+UW) | **B** |

### 4.5 Finance

| OQ | Ringkas | Blokir |
| --- | --- | --- |
| **OQ-020** | Arti kode `PaymentType`, `TransferType`, dan seluruh enumerasi (bersama Product+UW) | **D3**, **B** |
| **OQ-030** | Ambang hari-25 (bersama Product+UW) | **B** |
| **OQ-037** | Ambang nominal ter-hardcode menentukan roster komite (bersama Product+UW) | **B** |
| **OQ-040** | Limit wewenang: ter-hardcode di satu modul, dari database di modul lain (bersama Product+UW) | **B** |

### 4.6 Pemilik export Pega / Arsitektur Pega / Platform / Security

| OQ | Ringkas | Blokir |
| --- | --- | --- |
| **OQ-035** | Activity dipanggil tetapi salinannya hanya ada di modul lain | **B** |
| **OQ-043** | Baris tabel keputusan tidak ikut terekspor | **D4**, **B** |
| **OQ-049** | Empat class bersufiks `ENDORSEMENT`, satu berprefix `ASM-SFAGIS-` | **D4** |
| **OQ-054** | Folder `Claude outputs` berisi berkas non-Pega di dalam ekspor | D3 |
| **OQ-047** | `M_LINK_SERVICE` (Platform + DBA) | **B** |
| **OQ-048** | Gemini AI di alur underwriting (Security + Product+UW) | **D4**, **B** |
| **OQ-058** | Tabel internal Pega diakses langsung (Arsitektur Pega + DBA) | **B** |

### 4.7 Ringkas pemblokir

`[terverifikasi]` Dari 57 OQ terbuka:

| Kategori | Jumlah | Catatan |
| --- | ---: | --- |
| Memblokir **D3** (catatan modul + glossary) | **8** | OQ-003, 008, 011, 014, 020, 030–034 (sebagian), 054, 057 |
| Memblokir **D4** (context-map + understanding-report) | **24** | terutama OQ-001, 002, 009–012, 015–017, 019, 022–024, 028, 029, 038, 039, 041–043, 048–050, 056, 059 |
| Memblokir **FASE B** | **38** | seluruh batas korpus §3.1 + seluruh temuan RBAC + seluruh arti kode |

**Tidak satu pun OQ menghalangi dimulainya D3.** D3 mencatat *apa yang ada* per modul; OQ yang
memblokirnya (8) bersifat penamaan/arti, dan dapat dicatat sebagai "belum terverifikasi" tanpa
menghentikan pekerjaan.

---

## 5. Temuan lintas-domain yang paling menentukan

`[terverifikasi]` Lima pola yang muncul berulang dan mengubah cara memahami sistem:

### 5.1 Dua arsitektur modul yang berlawanan

| | **Satu ruleset, bercabang saat runtime** | **Modul benar-benar terpisah** |
| --- | --- | --- |
| Contoh | Facultative NB/RNW/EDM (**99,0 %** overlap identitas); `Treaty In` ↔ `Adjustment` (**85,8 %** file identik) | Trio Claim (**12–40 %** overlap) |
| Pembeda | rule `When` yang terbaca: `IsNB`/`IsRenewal`/`IsEDM` (`StatusBusiness` 1/2/3); `Param.type` + `EDMState` | **tidak ada** — nama Flow sama, **class berbeda** → rule berbeda |
| Bukti | `flows/_SUMMARY-facultative.md` §7; `flows/_SUMMARY-treaty-master.md` §3 | `flows/_SUMMARY-claim.md` |

### 5.2 Otorisasi hidup di luar model peran — empat pola berbeda

| Pola | Contoh | OQ |
| --- | --- | --- |
| Identitas orang sebagai **guard** | `OperatorID.pyUserIdentifier`, `pyPosition` | OQ-021 |
| Field bisnis menyimpan **kode peran** | `OperatorID.pyTelephone` = `TREATY1`/`TREATY2`/`SPVTREATY1`/`SPVTREATY2`; `pyWorkPage.LetterNo` = 10 token | OQ-027, **OQ-045** |
| Otorisasi lewat **indeks tetap** daftar workbasket | `OperatorID.pyWorkBasketList(2).pyWorkBasketName` (36×) | **OQ-051** |
| Identitas orang **ditetapkan sebagai tujuan rute** | `TreatyIn.PositionUsername := "<orang>"` — 5 identitas | **OQ-053** |

`[terverifikasi]` **Tidak ada satu pun rule identitas/otorisasi di korpus** (OQ-007). Seluruh
otorisasi tersebar sebagai nilai data dan ekspresi visibilitas.

### 5.3 Logika bisnis inti berada di database, bukan di rule

`[terverifikasi]` 67 stored procedure kustom, dan pada dua modul **seluruh jalur tulis master
melewati procedure** (`Treaty Contract Out` 9 jalur, `Master Contract Retro Life` 5 jalur).
Ditambah perhitungan prorata tanggal (`FIRE.CEK_PRORATA_TANGGAL`), tabel rate lewat db-link, dan
struktur JSON yang dibentuk fungsi kustom tanpa source.

### 5.4 Nama menipu, berulang kali

| Klaim dari nama | Kenyataan terverifikasi |
| --- | --- |
| `Treaty Contract Out` = treaty outward | **editor master arrangement**; 0 objek outward (OQ-022) |
| `pzOriginalInstanceKey` = identitas rule | **asal salinan**; identitas ada di `pxInsName` (D1) |
| `AKSEPTASI_ACT` berbeda antar modul | **identik** (`027a42f4`) — contoh OQ-011 yang usang |
| Nama Flow sama = proses sama | campuran: kadang satu rule, kadang class berbeda → rule berbeda (OQ-006) |
| Kode `OR`, `REINSTYPEID` punya arti terbaca | **tidak pernah muncul sebagai nilai literal** (OQ-057) |

### 5.5 Ekspor Pega punya artefak yang menyesatkan pengukuran

`[terverifikasi]` Tiga tag metadata (`pyDelete`, `pyVersionSecure`, `pzIsPrivateCheckOut`) yang
luput dari daftar normalisasi 18 tag menghasilkan **konflik semu**. Diukur: 9 dari 438 di
facultative, 3 dari 46 di `Treaty In` ↔ `Adjustment`, **7 dari 19** di `Master Product Name Life`
↔ `Treaty In`. Polanya konsisten — keluarga rule `T_STORAGE_IMAGE` dan `LINKSERVICE`.
**Register OQ-011 tetap berlaku**: 97,9 % konflik facultative adalah perbedaan isi yang nyata.

---

## 6. Kesiapan D3 dan D4

### 6.1 D3 — catatan pemahaman per modul + glossary final

**Siap dimulai.** Bahan yang tersedia:

| Bahan | Lokasi | Status |
| --- | --- | --- |
| Inventaris 20 modul (7 bagian per modul) | `inventory/<modul>.md` | lengkap |
| Rekap korpus | `inventory/_summary.md` (§22) | lengkap |
| Register konflik isi | `inventory/_oq011-konflik-isi.md` (533 entri) | lengkap |
| Telusur perilaku 20 konteks | `flows/<konteks>.md` | lengkap |
| 4 sintesis domain | `flows/_SUMMARY-*.md` | lengkap |
| Glossary berjalan | `glossary-draft.md` (8 blok tambahan dari D2) | perlu **finalisasi** |

**Yang harus dikerjakan D3:** menggabungkan inventaris (apa yang ada) dengan telusur (apa yang
dilakukan) menjadi satu catatan per modul, **untuk 20 modul** — termasuk 5 modul yang belum punya
artefak telusur tersendiri di luar inventarisnya. Lalu memfinalkan `glossary-draft.md`:
memisahkan istilah `[terverifikasi]` dari yang `[dugaan]`, dan menandai mana yang menunggu
jawaban OQ.

**Peringatan yang mengikat untuk D3:** 8 OQ memblokir sebagian isi (§4.7). Istilah yang artinya
belum terverifikasi **harus tetap ditulis demikian**; jangan diisi tebakan demi kelengkapan.

### 6.2 D4 — context-map + understanding-report

**Belum siap tanpa keputusan manusia.** Bahan bukti sudah lengkap, tetapi **24 OQ memblokir**
(§4.7), dan lima di antaranya menentukan bentuk peta konteks:

| OQ | Keputusan yang ditunggu |
| --- | --- |
| **OQ-015** | Facultative NB/RNW/EDM → satu konteks atau tiga? (bukti: satu ruleset) |
| **OQ-010** | `Treaty In` + `Adjustment` → satu konteks atau dua? (bukti: 85,8 % identik) |
| **OQ-022** | Ke mana `Treaty Contract Out` (303 rule) ditempatkan? (bukti: master arrangement, bukan outward) |
| **OQ-019** | Trio Claim → tiga konteks? (bukti: overlap 12–40 %, tidak ada pembeda runtime) |
| **OQ-042 / OQ-056** | Kepemilikan data lintas konteks: klaim membaca master treaty outward; master produk life menulis master treaty inward |

`[terverifikasi]` D2 **sengaja tidak** memutuskan kelimanya (`flows/_METHOD.md` §0). Bukti sudah
disajikan; keputusan adalah milik manusia di gate.

### 6.3 GATE

**FASE B (ADR/BRD/spec/tiket) tidak boleh dimulai sebelum gate manusia setelah D4.**
Yang di-review di gate: `understanding-report.md`, `context-map.md`, dan `open-questions.md`.

Sebelum FASE B, **38 OQ** perlu jawaban (§4.7) — terutama seluruh batas korpus §3.1, seluruh
temuan RBAC §5.2, dan seluruh arti kode/enumerasi.

---

## 7. Verifikasi penutup D2

`[terverifikasi]`

| Pemeriksaan | Hasil |
| --- | --- |
| Konteks ditelusur | **20** dari 20 modul |
| Rule ditelusur | **428** |
| Berkas `flows/` | **27** (20 konteks + 2 metode + 4 sintesis + README) |
| Format 7 bagian | **16 dari 20** konteks tepat 7 bagian; **4** memakai 8 bagian karena menampung satu pemeriksaan tambahan — `Endorsment Fac In.md` (varian `OfferFacRetro`), `Treaty In.md` (mesin status `Akseptasi_DT`), `Treaty In Adjustment.md` (pengukuran selisih + mesin revisi), `Master Product Name Life.md` (uji salinan subtree). Perintah audit: `grep -c '^## [0-9]' flows/*.md` |
| Rule identik ditelusur sekali | 4 `Flow` + 277 file `Treaty In`↔`Adjustment` + 16 file `Master Product Name Life`↔`Treaty In`, seluruhnya dengan bukti hash |
| Rule berkonflik ditelusur per-varian | `OfferFacRetro` (2 varian), `SERVICEINSERTARASAPAS_ACT` (2 varian), `ISCLM*` (per modul) |
| Korpus tersentuh? | **tidak** — nol file berubah di luar `OUTPUT_HASIL_RNM` |
| `D:\XML\nusantara-re\` tersentuh? | **tidak** |
| Nilai nama orang di artefak? | **tidak ada** — diverifikasi lewat grep |
| Setiap angka punya perintah audit? | ya |

Perintah audit korpus READ-ONLY:
```
find . -path ./OUTPUT_HASIL_RNM -prune -o -type f -newermt "-1 day" -print   # kosong
find "D:/XML/nusantara-re" -type f -newermt "-1 day" -print                  # kosong
```
