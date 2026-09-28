# Grilling — PremiumList Life + Endorsement Life — Ronde 1

Status: answered (work owner, 2026-09-15)
Konteks: `premiumlist-life` — "Life — Penawaran & Premium List" (context-map §2.8, cakupan `full`)
Modul: **PremiumList Life** (124 berkas) + **Endorsement Life** (75 berkas)
Skill: `/mattpocock-skills:grill-with-docs` (grilling + domain-modeling)

> **Konvensi penandaan.** `[terverifikasi]` = terbukti korpus dengan **class + nama + path**;
> `[keputusan work owner]`; `[data DBA]`; `[terbuka]` = OQ. Identitas rule **wajib menyertakan
> class** — nama sama di class berbeda = rule berbeda.

---

## Bagian 0 — KOREKSI atas fondasi di brief

Dua fakta di brief tidak cocok dengan korpus versi sekarang. Saya baca ulang; berikut yang benar.

`[terverifikasi]` `PremiumList Life/Activity/InsertJsonPolisLife_Act.xml`
(`ASM-FW-GISFW-WORK-LIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY`, **327.332 byte**)
tersimpan **28 Juli 2026** (`20260728T024422`) — lebih baru dari catatan D2 yang dipakai brief.

### K1. **17 langkah, bukan 21**

Nomor langkah dari `<pyStepPageReference>RH_1.pySteps(n)`: **1…17**. Tidak ada step 18–21.

### K2. `Commit` **dan** `Connect-REST` **ter-REMARK — keduanya mati**

`[terverifikasi]` Hanya **dua** `<pyStepsBlockName>` berisi `//` di seluruh berkas:

| Baris `//` | Step | Langkah |
| ---: | ---: | --- |
| 5294 | **16** | `Commit` (ref baris 5283) |
| 5383 | **17** | `Connect-REST` (ref baris 5372) |

Brief menyebut "Commit eksplisit (step 20)" dan "efek keluar … 21 Connect-REST". **Keduanya tidak
berlaku lagi.** Konsekuensinya besar dan dibahas di Q5 dan Q6.

---

## Bagian A — Peta 17 langkah `[terverifikasi]`

| Step | Langkah | Deskripsi | Precondition |
| ---: | --- | --- | --- |
| 1 | `Obj-Refresh-And-Lock` | | — |
| 2 | `Property-Set` | | — |
| 3 | `Property-Set` | set prodatetime | — |
| **4** | `Property-Set` | **set prodatetime > 25** | `@toDecimal(Local.currentdate)>25` |
| 5 | `Property-Set` | Pega to jsondata | — |
| 6 | `RDB-List` | Find list outward from product name | — |
| **7** | `RDB-List` + 4 sub-step | **percabangan treaty** | lihat §B |
| 8 | `Property-Set` + 2 sub-step | Insert to table summary | — |
| 9 | `Property-Set` | Pega to json_offer & lifeinproduction | — |
| 10 | `RDB-List` | Insert to table json_polis | — |
| 11 | `RDB-List` | Insert to table Lifeinproduction | — |
| 12 | `RDB-List` | **Cek sudah masuk atau blm datanya** | — |
| 13 | `Property-Set` | | `OutDataLife.pxResults(1).PL_NUMBER==""` |
| **14** | `call @baseclass.SendEmailNotification` | **Kalau blm, email errornya** | `IsPEGAPROD` |
| **15** | `call serviceInsertArasapasLife_act` | | `IsPEGAPROD` |
| ~~16~~ | ~~`Commit`~~ | | **REMARK** |
| ~~17~~ | ~~`Connect-REST`~~ | | **REMARK** (`IsPEGAPROD`) |

## Bagian B — Empat ID literal **punya label di korpus** (OQ-031)

`[terverifikasi]` Keempat precondition di sub-step step 7, beserta **deskripsi langkahnya**:

| Sub-step | Precondition | `pyStepsDescription` |
| --- | --- | --- |
| 7.1 | `@contains(.ID,"1000032")` (baris 1714) | **QS** |
| 7.2 | `@contains(.ID,"1000033")` (baris 1856) | **2nd QS** |
| 7.3 | `@contains(.ID,"1000034")` (baris 1998) | **SURPLUS** |
| 7.4 | `@contains(.ID,"1000035")` (baris 2140) | **2nd SURPLUS** |

Step 7 didahului step 6 **"Find list outward from product name"** — jadi ia mencabang atas daftar
**outward** hasil pencarian per nama produk.

**Ini jawaban parsial OQ-031 yang belum ada di register.** Register menanyakan "produk? mata uang?
jenis premi?" — labelnya menunjuk **jenis treaty proporsional**: Quota Share dan Surplus, masing-
masing dua lapis. `[dugaan]` kepanjangan QS = Quota Share; glossary korpus **belum**
memverifikasinya (`discovery/glossary.md` baris 216: `QS` — "kepanjangan belum terverifikasi").

## Bagian C — Aturan tanggal 25 terbaca penuh (OQ-030)

`[terverifikasi]` Mekanismenya jelas:

| Step | Nilai `ProdDateTime` |
| ---: | --- |
| 3 | `@CurrentDateTime()` — normal |
| **4** | `@CurrentDate("yyyy","Asia/Jakarta") + Local.NextMonth + "01T050000.000 GMT"` |

dengan `Local.currentdate = @CurrentDate("dd","Asia/Jakarta")` dan gerbang
`@toDecimal(Local.currentdate) > 25`.

**Artinya: transaksi yang diproses setelah tanggal 25 dibukukan ke periode bulan berikutnya** —
tanggal 1, pukul `05:00 GMT` (= 12:00 WIB).

Yang **belum** diketahui: aturan bisnisnya (tutup buku?) dan apakah ambang 25 tetap berlaku
pasca-migrasi. → Q2.

## Bagian D — DecisionTable: baris tidak terekspor, **hasilnya terbaca** (OQ-043)

`[terverifikasi]` Keduanya **nol baris keputusan** terekspor, tetapi nilai keluarannya terbaca:

| Rule | Identitas | Nilai keluaran |
| --- | --- | --- |
| `PremiumList Life/DecisionTable/IsLifeAccepted.xml` | `ASM-FW-GISFW-WORK-LIFE` / `ISLIFEACCEPTED`, 15.213 byte | `Confirm`, `Decline`, `Reject` |
| `PremiumList Life/DecisionTable/IsFlagOnGoingPolicy.xml` | `ASM-FW-GISFW-WORK-LIFE` / `ISFLAGONGOINGPOLICY`, 15.464 byte | `Decline`, `Offer`, `Premium` |

**OQ-043 menyempit:** yang tidak diketahui bukan *hasil apa saja* (kini terbaca), melainkan
**masukan apa yang menghasilkan masing-masing**. → Q1, pemblokir utama.

## Bagian E — Validasi upload CSV menyentuh **kolom uang**

`[terverifikasi]`

| Rule | Identitas |
| --- | --- |
| `PremiumList Life/Activity/UploadCSVLifePremium_Act.xml` | `@BASECLASS` / `UPLOADCSVLIFEPREMIUM_ACT` — **rule base-class, dipakai bersama** |
| `PremiumList Life/Activity/ValidasiUploadPL_act.xml` | `ASM-FW-GISFW-WORK-LIFE` / `VALIDASIUPLOADPL_ACT` |

Deskripsi langkah `ValidasiUploadPL_act` menyebut kolom yang divalidasi: `SET ERROR MESSAGE`,
**`NET_PREMIUM`**, **`GROSS_PREMIUM`**, **`SHARE_NUSANTARA_RE`**, **`SUM_INSURED`**,
**`CEDING_RETENTION`**, **`SUM_REASURED`**, **`CLAIM`**.

Seluruhnya **kolom uang** → menyentuh langsung **ADR-0003** (uang non-float, desimal presisi
arbitrer).

## Bagian F — Tautan Endorsement ↔ PremiumList

`[terverifikasi]` Kelas **berbeda** → rule berbeda, bukan konflik:

| Rule | PremiumList Life | Endorsement Life |
| --- | --- | --- |
| `InsertJsonPolisLife_Act` | `ASM-FW-GISFW-WORK-LIFE!…` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE!…` |
| `IsLifeAccepted` | `ASM-FW-GISFW-WORK-LIFE!…` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE!…` |
| titik masuk flow | `InputPolicyHolder` (`…WORK-LIFE`) | `InputEDMLife` (`…WORK-ENDORSEMENTLIFE`) |

`[terverifikasi]` **`PL_NUMBER_EDM` adalah tautannya** — 20 kemunculan di Endorsement Life, 2 di
PremiumList Life:

| Berkas | Modul |
| --- | --- |
| `Activity/GenerateNoEDM_Life.xml` | Endorsement — penomoran EDM |
| `RDBList/InsertJsonPolisEDM.xml`, `RDBList/SaveMasterLPDet.xml`, `Section/ConfirmSubmitEDM.xml` | Endorsement |
| `RDBList/InsertPLSummary.xml` | **keduanya** |

## Bagian G — Frontier Ronde 1

Enam pertanyaan, diajukan di percakapan. Tiga **fakta bisnis** (OQ tim — tidak saya tebak) dan tiga
**keputusan desain** (saya beri rekomendasi).

*(Jawaban ditempel ke berkas ini pada ronde berikutnya.)*

---

# Jawaban work owner — enam verdict (2026-09-15)

## V1. Keputusan alur — **manual, bukan formula** (OQ-043 TUTUP)

`[keputusan work owner]` Ketiga keluaran `IsLifeAccepted` punya akibat berbeda:

| Keluaran | Akibat |
| --- | --- |
| **Confirm** (accept) | naik ke tahap berikutnya |
| **Reject** | **selalu balik ke input** — inputor/admin mengisi ulang |
| **Decline** | case **ditutup/dibuang** — berhenti |

`IsFlagOnGoingPolicy`: **`1` = offer** → berhenti di tahap penawaran; **`2` = premium** → lanjut ke
Input Premium List Detail.

**Keputusan accept/reject/decline dibuat MANUAL oleh inputor/admin** — bukan formula otomatis.
Karena itu **kriteria masukannya bukan aturan yang perlu direplikasi**; yang perlu direplikasi adalah
**tiga akibat** di atas.

Konsisten dengan graf: `Decision1` dan `Decision3` memakai `IsLifeAccepted`.
`[terverifikasi]` Nilai keluaran terbaca di kedua DecisionTable meski barisnya tidak terekspor.

## V2. Tutup buku dari **tabel**, bukan hardcode (OQ-030 TUTUP)

`[terverifikasi]` Ambangnya **bukan** `25`. `PremiumList Life/Activity/SubmitPremiumList_Act.xml`
(`ASM-FW-GISFW-WORK-LIFE` / `SUBMITPREMIUMLIST_ACT` / `RULE-OBJ-ACTIVITY`, 214.154 byte, tersimpan
**`20260122T072213`**):

| Baris | Isi |
| ---: | --- |
| 773 | memanggil `GETTanggalClosing_SQL` — `SELECT * FROM POOLDATA.TANGGAL_CLOSING` |
| 918-919 | `Local.TglProd = TglProd.pxResults(1).TANGGAL` |
| **966** | `@if(Local.TglProd=="", 25, Local.TglProd)` — **`25` hanyalah FALLBACK bila tabel kosong** |
| 1211 | `@if(@toDecimal(Local.currentdate)>Local.TglProd, @toDecimal(Local.CurrentMonth)+1, ...)` |
| 3491 | `@toDecimal(Local.currentdate) > Local.TglProd` |

**Aturan:** transaksi setelah tanggal closing dibukukan ke periode **bulan berikutnya** (tanggal 1,
`05:00 GMT` = 12:00 WIB).

⚠️ **Dua versi berbeda hidup berdampingan** `[terverifikasi]`:

| Rule | Ambang | Tersimpan |
| --- | --- | --- |
| `SubmitPremiumList_Act` | dari `TANGGAL_CLOSING` | `20260122T072213` |
| `InsertJsonPolisLife_Act` | **hardcode `>25`** (baris ~1170) | `20260728T024422` |

⛔ *Ralat 28-09-2026 (sensus remark GILIRAN-12):* `SubmitPremiumList_Act` langkah 2–4 membaca
`TANGGAL_CLOSING`, tetapi hasilnya hanya mengalir ke `TempGenerate.CARI1/2` (tak dibaca rule mana
pun) dan ke langkah 15 yang **ter-remark** (`//` b3398; b3491 miliknya). Dua aturan yang **hidup**:
nomor PL digulir `PROC_GENERATE_SEQUENCE_NUMBER` (membaca tabel itu sendiri), dan `ProdDateTime`
digeser `InsertJsonPolisLife_Act` langkah 4 (`>25` tertanam, b1170). Keputusan "ikuti yang dari DB"
tetap; penerapannya pada `ProdDateTime` (yang di Pega memakai 25) dikonfirmasi ulang — **OQ-PL-13**.

`[keputusan work owner]` **Ikuti yang dari DB.** Sistem baru membaca `TANGGAL_CLOSING`/konfigurasi;
**jangan tanam `25`**.

## V3. Empat treaty ID = **logika polis lama** (OQ-031 TUTUP)

`[terverifikasi]` Step 6-7 `InsertJsonPolisLife_Act` menggerbangi empat `Property-Set` dengan
`@contains(.ID,"1000032")` … `"1000035"`, atas class `ASM-FW-GISFW-Int-TREATYYEAR_LIFE`, berlabel
**QS / 2nd QS / SURPLUS / 2nd SURPLUS**.

`[keputusan work owner]` **Ini logika POLIS LAMA, sudah tidak dipakai** — kemungkinan kelupaan
di-remark. ⚠️ Di korpus `blockname` masih **kosong (aktif)**, tetapi work owner menyatakan tidak
dipakai. **Sistem baru: JANGAN replikasi keempat gerbang ID ini.**

`[terverifikasi]` **Arti kode jenis treaty terbaca penuh** di baris ~3787 berkas yang sama:

```
@if(pyWorkPage.TypeCeding="1","QS",
 @if(pyWorkPage.TypeCeding="2","SURPLUS",
  @if(pyWorkPage.TypeCeding="3","QS + SURPLUS",
   @if(pyWorkPage.TypeCeding="4","XOL",""))))
```

**`TypeCeding`: 1 = QS (Quota Share), 2 = SURPLUS, 3 = QS + SURPLUS, 4 = XOL.**

Bila sistem baru butuh jenis treaty, ambil dari **data `TypeCeding`**, bukan konstanta ID.

## V4. Efek keluar — dua aktif, satu di antaranya **alarm**

`[terverifikasi + keputusan work owner]`

| Step | Efek | Sifat |
| ---: | --- | --- |
| 15 | `serviceInsertArasapasLife_act` | **efek bisnis**; endpoint via `M_LINK_SERVICE` (**ADR-0013**) |
| 14 | `SendEmailNotification` | **JALUR ALARM** — bergerbang `OutDataLife.pxResults(1).PL_NUMBER==""`, deskripsi "Kalau blm, email errornya". **Bukan notifikasi bisnis** |
| ~~17~~ | ~~`ConvertJsonNusareToProduction` (Connect-REST)~~ | **SUDAH TIDAK DIPAKAI** — dibuang, jangan dihidupkan |

Keduanya bergerbang `IsPEGAPROD` — flag lingkungan (env var di sistem baru, **ADR-0005**).

## V5. Batas transaksi — Go yang memegang

`[terverifikasi]` Step 1 `Obj-Refresh-And-Lock` mengunci; step 16 `Commit` **ter-remark**.

`[keputusan desain]` Sistem baru: **Go memegang batas transaksi eksplisit** — buka, simpan polis +
generate nomor, commit/rollback. **Commit segera setelah nomor** agar lock `SELECT ... FOR UPDATE`
lekas lepas. Jangan andalkan commit implisit Pega. Konsisten **OQ-013** (Claim Life).

## V6. Endorsement & `PL_NUMBER_EDM` — dua jalur penomoran sejajar

`[keputusan work owner]`

- **Field yang boleh diubah endorsement ditentukan kondisi `When` di TIAP FORM** — **tidak ada
  daftar global**.
- **`PL_NUMBER_EDM` = nomor PL khusus EDM**, dibuat **terpisah** agar tidak mengganggu `PL_NUMBER`
  dari NB (new business). **Dua jalur penomoran yang tidak saling timpa** — **BUKAN pointer ke PL
  asal.**
- Endorsement diperlakukan sebagai **entri baru dengan penomoran EDM sendiri**, berdampingan dengan
  `PL_NUMBER` NB.

⚠️ **Ini meralat rekomendasi saya di Ronde 1** ("endorsement menunjuk pendahulunya lewat
`PL_NUMBER_EDM`"). Keduanya **sejajar**, bukan berantai.

`[terverifikasi]` Didukung korpus: `PremiumList Life/RDBList/InsertPLSummary.xml`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY` / `RULE-CONNECT-SQL`) memanggil
`POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY(...)` dengan **`PL_NUMBER` dan `PL_NUMBER_EDM` sebagai dua
parameter terpisah** pada pemanggilan yang sama.

---

# Bagian 4 — Telusur Ronde 2

## T1. OQ-023 — Assignment "Input Premium List Summary" **tanpa jalur masuk terbaca**

`[terverifikasi]` `PremiumList Life/InputPolicyHolder.xml`
(`ASM-FW-GISFW-WORK-LIFE` / `INPUTPOLICYHOLDER` / `RULE-OBJ-FLOW`). Shape
**`Assignment1` = "Input Premium List Summary"** (label baris 1234, subscript baris 1248).

Sapuan seluruh kemunculan `Assignment1`:

| Baris | Tag | Artinya |
| ---: | --- | --- |
| 328 | `<pyTaskID>` | id tugasnya sendiri |
| 585 | `<pyFromTaskName>` | ia **sumber** sebuah connector |
| 1215 | `<pyMOId>` | id shape |
| 1248 | `<pxSubscript>` | subscript shape |
| 1900 | `<pyFrom>` | ia **sumber** connector |

**Ia tidak pernah muncul sebagai `<pyTo>` maupun `<pyToTaskName>`** — **nol connector masuk**.

`[terverifikasi]` Dua jalur alternatif juga nihil:

- **Harness:** tidak satu pun dari 9 harness modul merujuk `ShowLifePremiumSummary` maupun
  `Assignment1`.
- **Ticket:** `Ticket1` berjenis `Data-MO-Event-Exception` dan menuju `Decision2`, bukan Assignment1.

Shape itu **punya** FlowAction (`PremiumList Life/FlowAction/ShowLifePremiumSummary.xml`) dan
**punya** connector keluar — jadi ia dapat **ditinggalkan**, tetapi **tidak dapat dicapai** lewat apa
pun yang terbaca di ekspor.

**Verdict:** `[terverifikasi]` **tidak ada jalur masuk di dalam korpus**. Apakah ia dicapai lewat
mekanisme di luar 17 tipe rule yang diekspor (mis. navigasi/menu Pega) **tidak dapat ditentukan dari
ekspor** — OQ-023 **dipersempit**, tetap terbuka.

## T2. Kontrak hulu ke Claim Life

`[terverifikasi]` **Class integrasi `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` adalah tulang punggung
domain Life** — dipakai keempat modul:

| Modul | Jumlah rule yang merujuk |
| --- | ---: |
| Claim Life | **46** |
| Endorsement Life | 25 |
| PremiumList Life | 19 |
| Komite Claim Life | 11 |

`[terverifikasi]` **Penulisan tabel premi lewat stored procedure, bukan `INSERT` langsung:**

| Rule (PremiumList Life/RDBList) | Identitas | Procedure / tabel |
| --- | --- | --- |
| `InsertPLSummary.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY` | `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY(...)` — menerima **`PL_NUMBER` dan `PL_NUMBER_EDM`** |
| `InsertJsonPolis.xml` | `RULE-CONNECT-SQL` | `POOLDATA.INSERTJSONPOLISLIFE` |
| `SaveOfferJsonLife_SQL.xml` | `RULE-CONNECT-SQL` | `POOLDATA.INSERTJSONOFFERLIFE` |
| `GetSequenceNumber_SQL.xml` | `RULE-CONNECT-SQL` | `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` — **sama dengan Claim Life** (**ADR-0006**) |
| `SaveLifeinProduction_SQL.xml` | `RULE-CONNECT-SQL` | `INSERT INTO POOLDATA.LIFEINPRODUCTION` |
| `InsertDataUploadLife.xml` | `RULE-CONNECT-SQL` | `insert into POOLDATA.M_TEMPUPLOADLIFE` — **staging upload CSV** |
| `InsertLogServiceProd.xml` | `RULE-CONNECT-SQL` | `INSERT INTO pooldata.MONITORING_PROD_LOG` |

Pembacaan: `POOLDATA.JSON_OFFER_LIFE` (`GetIdOffer_SQL`, `GetOfferLife_sql`), `POOLDATA.JSON_POLIS`
(`GetPLandNopolis_sql`).

**Kontrak hulu ke hilir:** PremiumList Life dan Endorsement Life **menghasilkan baris premium
summary/detail** (lewat `PEGA_M_LIFE_PREMIUM_SUMMARY`), berkunci `PL_NUMBER` (NB) atau
`PL_NUMBER_EDM` (endorsement). Claim Life kemudian **mengonsumsinya** sebagai
`PremiumListSummary` / `PremiumListDetail` — persis struktur yang dipakai mesin status klaim
(**ADR-0011**: `PremiumListDetail` adalah cerminan baris `AdjustmentList`).

⚠️ **Batas pengetahuan baru:** badan `PEGA_M_LIFE_PREMIUM_SUMMARY`, `INSERTJSONPOLISLIFE`, dan
`INSERTJSONOFFERLIFE` **tidak ada di korpus** — keluarga **OQ-002**. Kolom dan aturan yang mereka
tulis tidak terbaca. Ini **pemblokir kontrak hulu**.
