# F05: Empat struktur penampung Fac Retro — dipisahkan, bukan disatukan

**What to build:** Data retrosesi hidup di **empat** struktur berbeda sepanjang satu kasus. Tiket ini
memodelkan keempatnya beserta perannya, sehingga tidak ada yang tertukar.

⛔ **EMPAT, bukan dua.** Catatan lama yang menyebut "dua struktur" (`09-facout\01` §2.1) **tidak
lengkap**:

| # | Struktur | Bentuk | Diisi oleh |
| ---: | --- | --- | --- |
| 1 | `OfferFacIn.FacRetro` | **tunggal** (staging) | `SetDataFacOutFire_Act` — `FacRetro.LocationList(<APPEND>)` |
| 2 | `OfferFacIn.FacRetroList` | **list berindeks** | `CopyFacRetroFire_ACT` — `FacRetroList(counterFacRetro)…` |
| 3 | `OfferFacIn.FacRetroDetails` | tunggal, 4 properti | `OfferFacOut_PreAct` · `GetRISlipDataFromDB_Act` · `SetPropertyToOfferFacIn_Act` |
| 4 | `pyWorkPage.Policy.FacOfferList(n)` | list berindeks | sumber salinan bagi struktur 3 |

`[terverifikasi]` **`FacRetroDetails` punya tepat empat properti**, seluruhnya dari tag
`<PropertiesName>` / `<PropertiesValue>` / `<pyValue>` / `<pyStepsPreCondParamsWhen>`:

| Properti | Ditulis | Dibaca |
| --- | --- | --- |
| `.BackUpStatus` | `GetRISlipDataFromDB_Act` | `CheckProtectFacout_Act` (`==3` **telanjang**) · `OfferFacOut_PostAct` (`=="3"` **berkutip**) |
| `.FacNo` | `OfferFacOut_PreAct` | — |
| `.AdditionalInfo` | `OfferFacOut_PreAct` | — |
| `.DocumentPosition` | `OfferFacOut_PreAct` · `SetPropertyToOfferFacIn_Act` | `Harness\ViewLetter` · `Section\FacOutPrintRISlipSectionInside` · `Section\FacultativeLetter` |

⚠️ **Nilai yang sama dibandingkan dua gaya** — `==3` telanjang dan `=="3"` berkutip. Angka
diperlakukan sebagai string di satu tempat dan sebagai angka di tempat lain. **Kandidat K-046, diport
apa adanya.**

### Anak langsung `FacRetroList` — **19**, bukan enam

> ⛔ **KOREKSI 21 September 2026.** Tiket ini semula menulis ~~"`FacRetroList` punya **enam** koleksi
> anak: `LocationList` · `CargoList` · `PersonList` · `VehicleList` ·
> `Property.RiskLocation.AnekaList` · `Property.RiskLocation.OccupationList(n).AnekaList`"~~.
> **Salah pada dua hal:** dua yang terakhir **bukan anak langsung** (letaknya bersarang di bawah
> `LocationList(n).Property.RiskLocation`), dan daftarnya **melewatkan `CurrencyList`**.

`[terverifikasi]` `FacRetroList(n)` punya **19 anak langsung** — **5 koleksi** dan **14 non-koleksi**:

| Koleksi (`*List`) — 5 | Non-koleksi — 14 |
| --- | --- |
| `CargoList` · **`CurrencyList`** · `LocationList` · `PersonList` · `VehicleList` | `EndPeriod` · `OurRef` · `PctPremiAllObj` · `PctPremiAllObjUSD` · `PCTPremiIDR` · `PctShareAllObj` · `PrintRISlip` · `ReinsurerID` · `ReinsurerName` · `RiCommAllObj` · `StartPeriod` · `TFAllObj` · `UjrahAllObj` · `pyExpanded` |

⚠️ `pyExpanded` adalah properti tampilan Pega, bukan data bisnis — dicatat agar sensusnya dapat
direproduksi.

📌 Lini Aneka memang dijangkau lewat `LocationList(n).Property.RiskLocation.AnekaList` dan
`…OccupationList(o).AnekaList`, tetapi itu **jalur bersarang**, bukan anak `FacRetroList`.

**Sumber angka:** `..\..\10-audit\03-struktur-facretro-dan-populasi.md` §2.1.

`[terverifikasi]` **`FacOutTSI` melekat di coverage, bukan di `FacOutObjectList`.** Seluruh jalur
berbentuk `.CoverageList(1).FacOutTSI` atau saudara PA-nya `.ASMCoverage(1).FacOutTSI`; **nol** jalur
berbentuk `FacOutObjectList(…).FacOutTSI`. Anggota `FacOutObjectList(1)` yang sebenarnya: `Rate`,
`ShareOffered`, `PercentOffered`, `ObjectPremi`, `RiComm`, `commision`, `TF`, `Ujrah`.

⛔ **`RiComm` di `FacOutObjectList` BUKAN `.RIComm` di coverage.** Beda kapitalisasi, beda properti,
beda pemilik. Begitu pula `commision` (ejaan pendek) terhadap `COMMISION` di tabel produksi.

⚠️ **Gaya nilai flag berbeda antar properti.** `[terverifikasi]` `.IsFacRetroOffer` disetel **`"0"`
berkutip** di `OfferFacOut_PreAct`, sedangkan `.IsFacRetro` memakai `0`/`1` **telanjang**. Properti
berbeda, gaya berbeda — **diport apa adanya**.

⚠️ `[terverifikasi]` Gerbang `CopyFacRetroFire_ACT` menguji **`local.Facout=="0"`** — angka
dibandingkan **sebagai string**. Kandidat K-046, diport apa adanya.

📌 `[terverifikasi]` **Nol kemunculan `FacRetroDetails` di `DDL\`** → ia struktur clipboard/JSON,
**bukan** tabel Oracle. Tidak ada DDL yang perlu dicari untuknya.

### ⛔ Empat properti bernama mirip — **TIDAK BOLEH disamakan**

`[terverifikasi]` Empat properti berbeda dengan nama yang nyaris sama hidup berdampingan. Menyamakan
dua di antaranya akan menggeser gerbang alur:

| Properti | Nilai & gaya | Ditulis / dibaca di | Perannya |
| --- | --- | --- | --- |
| `OfferFacIn.IsFacRetro` | `0` / `1` **telanjang** | disetel `0` di `SetDataFacOut_Act` (`RH_1.pySteps(3)`, `<PropertiesName>`); disetel `1` di **keempat** `SetDataFacOut{Fire,AnekaGolf,CargoMBU,PATravel}_Act` | **pemicu Fac Out** — menjadi `<pyTaskWhen>` transisi di flow masuk siklus ybs. |
| `OfferFacIn.IsInputFacRetro` | `0` / `1` | `<pyStepsPreCondParamsWhen>` `Param.Status=="UW" && …IsInputFacRetro==1` | gerbang **re-input** saat UW; **9** kemunculan di `NB FacIn\Flow\InputInwardFacultativeOffer` |
| `pyWorkPage.IsInFacRetro` | dibanding `!= 1` | `<pyContainerVisibleWhen>` di `Section\InputCoverageFire` | **kondisi tampilan** layar, bukan alur |
| `.IsFacRetroOffer` | **`"0"` berkutip** | disetel di `Activity\OfferFacOut_PreAct` (`<PropertiesName>`) | penanda penawaran retro |

⚠️ **Gayanya berbeda dan itu disengaja dipertahankan:** `.IsFacRetro` memakai angka telanjang,
`.IsFacRetroOffer` memakai string berkutip. **Diport apa adanya** (K-046).

⚠️ `[terverifikasi]` Reset memakai jalur **relatif** `.OfferFacIn.IsFacRetro`; penyalaan memakai jalur
**absolut** `pyWorkPage.OfferFacIn.IsFacRetro`. Dua bentuk penulisan untuk properti yang sama.

**Sumber angka:** `..\..\10-audit\03-struktur-facretro-dan-populasi.md` ·
`..\..\10-audit\06-tinjau-ulang-suntingan-facout.md` §2.5.

**Asal (Pega).** `Activity\SetDataFacOutFire_Act` · `CopyFacRetroFire_ACT` · `CopyFacRetro_ACT` ·
keluarga `CopyAllObjFacOut{FireAneka,GolfCargo,PAMBU}_ACT` · `OfferFacOut_PreAct` ·
`GetRISlipDataFromDB_Act` · `SetPropertyToOfferFacIn_Act`

**Keputusan.** K-053 · **K-046** · K-010/K-012 · `CLAUDE.md` §1, §4.6

**Blocked by:** F02

**Status:** blocked

- [ ] **Empat struktur** dimodelkan terpisah; tidak ada yang dilebur
- [ ] `FacRetroDetails` punya **tepat empat properti**, tidak lebih
- [ ] **Enam koleksi anak** `FacRetroList` terwakili
- [ ] `FacOutTSI` melekat di **coverage**, bukan di dalam `FacOutObjectList`
- [ ] `RiComm` (obyek fac out) dan `.RIComm` (coverage) adalah **dua properti berbeda** yang tidak dapat tertukar
- [ ] **K-046** `K046_BackUpStatus_DuaGayaPembanding` — `==3` dan `=="3"` berdampingan
- [ ] **K-046** `K046_LocalFacout_AngkaSebagaiString` — `local.Facout=="0"`
- [ ] **K-046** `K046_IsFacRetroOffer_NolBerkutip` — `"0"` berkutip vs `0` telanjang
- [ ] Nilai uang bertipe `Money`; rasio bertipe `Ratio`
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
