# Fac Out / Fac Retro — Temuan Korpus & Rancangan Menu Terpisah

Dokumen ini merekam **mekanisme Fac Out (retrocession)** yang ditemukan di korpus dan mengusulkan
rancangan sebagai **menu terpisah dari NB/RNW/EDM namun tetap terhubung**.

> Kepatuhan CLAUDE.md §3: tiap klaim menyebut path + isi tag + label. **Nama rule bukan bukti** —
> yang dikutip di bawah adalah isi langkah/kondisi, bukan sekadar nama berkas.

---

## 0. Ringkasan verifikasi klaim work owner

Klaim work owner: *"Saat UW melakukan spreading, jika ada spreading Fac Out, saat UW Accept
LocationList yang punya spreading Fac Out disalin ke FacRetroList (via SetValidateDateUW_PostAct →
SetDataFacOut_Act) dan masuk menu Fac Out dulu. Di flow InputInwardFacultativeOffer ada yang memanggil
flow OfferFacRetro."*

| Bagian klaim | Status | Bukti |
| --- | --- | --- |
| `SetValidateDateUW_PostAct` memanggil `SetDataFacOut_Act` | **[terverifikasi]** | NB `Activity\SetValidateDateUW_PostAct.xml` step 3 `<pyStepsActivityName>Call SetDataFacOut_Act</pyStepsActivityName>` |
| Penyalinan risiko ke halaman Fac Retro terjadi | **[terverifikasi]** | NB `Activity\SetDataFacOut_Act.xml` step 1 Property-Remove atas `pyWorkPage.OfferFacIn.FacRetro.{LocationList,PersonList,CargoList,VehicleList}`; step 3-6 memanggil `SetDataFacOut{Fire,AnekaGolf,CargoMBU,PATravel}_Act` |
| Pemicu berbasis "spreading Fac Out" | **[terverifikasi, dengan koreksi lokasi]** | penanda ada di sub-aktivitas per-COB, bukan di `SetDataFacOut_Act` induk. Lihat §2 |
| Target adalah `FacRetroList` | **[terverifikasi — klaim work owner BENAR]** | ADA DUA struktur: staging `OfferFacIn.FacRetro` (tunggal, diisi `SetDataFacOutFire_Act` step 7 `FacRetro.LocationList(<APPEND>)`) DAN `OfferFacIn.FacRetroList` berindeks (diisi `CopyFacRetroFire_ACT` step 1, `FacRetroList(counterFacRetro).LocationList(...)`). Koreksi saya sebelumnya (yang menyingkirkan `FacRetroList`) KELIRU |
| `InputInwardFacultativeOffer` memanggil `OfferFacRetro` | **[terverifikasi, via shape]** | NB `Flow\InputInwardFacultativeOffer.xml` menyebut `OfferFacRetro` 5×; `pyMOName` memuat shape `[Fac Out]`, `FAC OUT?`, `INPUT FAC OUT?`. **Bukan** lewat `<pySubFlowName>` (itulah kenapa pencarian tag polos gagal — sesuai peringatan "UI/shape mengalahkan XML") |
| Masuk menu Fac Out "terlebih dahulu" | ~~**[terverifikasi sebagian]** When `IsFacout` = `pyWorkPage.ProposalAcceptStatus = 4`; flow punya shape `INPUT FAC OUT?` sebelum produksi.~~ ⛔ **DIBATALKAN** — lihat §0.1. Penggantinya: gerbang `Transition92` ber-`pyTaskWhen` = **`IsFacRetro`**. `[terverifikasi]` |

**Kesimpulan:** substansi klaim benar seluruhnya. Satu koreksi tersisa: logika pemicu ada di
sub-aktivitas per-COB (`SetDataFacOut{Fire,…}_Act`, `CopyFacRetro{Fire,…}_ACT`), bukan di
`SetDataFacOut_Act` induk. Klaim `FacRetroList` **benar** — ada dua struktur (`FacRetro` staging +
`FacRetroList` berindeks), lihat §2.1.

---

## 0.1 ⛔ KOREKSI PEMICU — 21 September 2026

> **Rumusan lama di §0, §3 dan §4.2 DIBATALKAN.** Bunyinya:
>
> ~~"UW Accept (`ProposalAcceptStatus = 4` → When `IsFacout`)"~~
>
> **Keliru pada DUA lapis sekaligus**, dan keduanya berdiri sendiri:
>
> 1. ⛔ **`IsFacout` tidak pernah terpanggil di jalur Fac Out.** `[terverifikasi]`
>    `NB FacIn\Flow\InputInwardFacultativeOffer.xml` memuat `IsFacout` **nol kali**, peka huruf
>    maupun tidak. Di `SetValidateDateUW_PostAct.xml` satu-satunya kemunculan `isFacOut` adalah
>    **deklarasi variabel lokal** — tag `<pyLocalParameters>` → `<pyParametersParamName>` =
>    `isFacOut`, bertipe `String` — bukan pemanggilan rule. Kesamaan nama variabel lokal dengan
>    nama rule inilah yang menyesatkan pembacaan lama (`PANDUAN-KERJA` §3: nama bukan bukti).
> 2. ⛔ **`ProposalAcceptStatus = 4` bukan Accept, melainkan BANDING.** `[terverifikasi]` **K-029**
>    dan `steering\GLOSARIUM.md` §"Hasil keputusan": `1` Accept · `2` Reject · `3` Ask ·
>    **`4` Banding** · `7` Decline · `9` Revise. GLOSARIUM §"Kosakata yang dihindari" bahkan sudah
>    mencantumkan `IsFacout` sebagai nama yang **tidak menggambarkan isinya** — "isinya menguji
>    banding". Jadi rumusan lama memakai rule yang tidak terpanggil, untuk menguji keadaan yang
>    bukan Accept.
>
> ⚠️ **Keputusan K-019 tidak berubah** oleh koreksi ini. K-019 menetapkan fitur fac out usang secara
> bisnis tetapi **kodenya tetap diport apa adanya**; itu tetap berlaku. Yang dikoreksi di sini adalah
> **dasar pengukuran jumlah perujuknya** — lihat `10-audit\04-perujuk-isfacout-peka-huruf.md`.

### Pemicu Fac Out yang BENAR — rantai lima langkah

`[terverifikasi]` Saat kasus berada di posisi Underwriter, UW melakukan spreading. Bila spreading yang
dipilih memuat spreading Fac Out, maka **saat UW Accept** lokasi yang punya spreading Fac Out disalin
ke halaman retrosesi, dan kasus masuk ke menu Fac Out lebih dulu.

⚠️ **Rantai langkah 1–3 diukur pada folder NB.** Langkah 4–5 berlaku di ketiga siklus, tetapi **lewat
flow masuk masing-masing** — lihat catatan di bawah tabel.

| # | Langkah | Tag pembawa | Isi |
| ---: | --- | --- | --- |
| 1 | `Activity\SetValidateDateUW_PostAct` | `<pyStepPageReference>` = `RH_1.pySteps(3)` | `Call SetDataFacOut_Act`, **TANPA precondition** — `<pyStepsPreCondParamsWhen>` kosong |
| 2 | `Activity\SetDataFacOut_Act` | `<Property>` ×4 pada `RH_1.pySteps(2)`; `<PropertiesName>` pada `RH_1.pySteps(3)` | `RH_1.pySteps(2)` `Property-Remove` atas `OfferFacIn.FacRetro.LocationList` · `.PersonList` · `.CargoList` · `.VehicleList`; `RH_1.pySteps(3)` `Property-Set` `.OfferFacIn.IsFacRetro = 0` (**reset**); `RH_1.pySteps(4..7)` memanggil `SetDataFacOut{Fire,AnekaGolf,CargoMBU,PATravel}_Act`. ⚠️ Struktur `FacRetroList` dibersihkan **terpisah**, di dalam iterasi `RH_1.pySteps(8).pySteps(1)` — lihat §3 |
| 3 | `Activity\SetDataFacOut{Fire,AnekaGolf,CargoMBU,PATravel}_Act` | `<pyStepsPreCondParamsWhen>` + `<PropertiesName>` | **keempatnya** menyetel `pyWorkPage.OfferFacIn.IsFacRetro = 1`, bergerbang `.TreatyType=="10015"` / `Local.isFacOut==1` / `Local.isFacOutLoc==1` |
| 4 | flow masuk siklus ybs. | `<pyTaskStatusOrWhen>`=`WHEN`, `<pyTaskWhen>`=`IsFacRetro`, `<pyTaskStatus>`=`IsFacRetro`, `<pyLikelihood>`=100 | transisi menuju shape `<pyMOName>` = `[Fac Out]`. Di NB `InputInwardFacultativeOffer` transisi itu ber-`<pyID>`/`<pyMOId>` = **`Transition92`** (⚠️ `<pyTaskName>` **kosong** — nama dibawa `pyID`, bukan `pyTaskName`) |
| 5 | shape `[Fac Out]` | **`<pyImplementation>`** = `OfferFacRetro` | ⚠️ dirujuk lewat tag **`pyImplementation`**, **bukan** `pySubFlowName` — `<pySubFlowName>` memuat `OfferFacRetro` **0 kali**. Itulah sebab pencarian tag polos dahulu gagal |

> ⛔ **KOREKSI 21 September 2026 atas kalimat pembuka.** Sesi sebelumnya menulis
> ~~"Berlaku **sama di NB, RNW dan EDM**"~~ sambil menyebut `Flow\InputInwardFacultativeOffer`.
> **Berkas itu HANYA ADA DI NB** `[terverifikasi]` — `Flow\` RNW berisi 4 berkas dan Endorsment 2,
> tidak satu pun bernama demikian.
>
> **Substansinya tetap berlaku lintas siklus, lewat flow masuk masing-masing** `[terverifikasi]`:
>
> | Flow | `pyTaskWhen=IsFacRetro` | `pyImplementation=OfferFacRetro` |
> | --- | ---: | ---: |
> | NB `InputInwardFacultativeOffer` | 2 | 2 |
> | NB `InputInwardFacultativeRISlip` | 3 | 3 |
> | RNW `InputInwardFacultativeRISlip` | 3 | 3 |
> | RNW `InputRenewalFacultativeIn` | 3 | 3 |
> | Endorsment `InputAddendumFacultativeIn` | 3 | 3 |
>
> 📌 **NB punya dua flow** yang berbelok ke Fac Out, bukan satu. Rinciannya di
> `10-audit\06-tinjau-ulang-suntingan-facout.md` §2.8.

**Bentuk pemicunya, dinyatakan tegas:** `[terverifikasi]` Fac Out dipicu oleh **flag `.IsFacRetro`**
yang di-*reset* ke `0` di langkah 2 lalu **dinyalakan ke `1` per COB** di langkah 3 ketika penanda
spreading Fac Out ditemukan. Penyaringan terjadi **di precondition per COB**, bukan di titik masuk.
Flag itu kemudian menjadi gerbang `Transition92` di flow.

`[terverifikasi]` Jumlah kemunculan di **`NB FacIn\`**`Flow\InputInwardFacultativeOffer.xml`:
`IsFacout` **0** (peka huruf maupun tidak) · `IsFacRetro` **17** · `IsInputFacRetro` **9** ·
`IsFacRetroOffer` **0**.

> `[pertanyaan terbuka]` Indeks rujukan flow menuliskan resolusi `IsFacRetro` ke app **SFAGIS** /
> workType **ClaimLife**, dan `OfferFacRetro` ke app **GISFW** / workType **LIFE**, sementara
> `Embed-Reference-Rule` menyebut kelas `ASM-FW-GISFW-Work` dan definisi ekspornya berkelas
> `ASM-FW-GISFW-Data-OfferFacIn`. **Butuh UI Pega work owner. Jangan disimpulkan.**

---

## 1. Berkas Fac Out/Retro yang [terverifikasi ADA] di korpus

Semua tiga siklus (NB/RNW/EDM) punya set berkas ini (jumlah bervariasi per folder):
- **Flow:** `OfferFacRetro.xml` (NB/RNW/EDM), `OfferFacOut.xml` (Flow & FlowAction).
- **Activity salin:** `SetDataFacOut_Act` (induk) + `SetDataFacOut{Fire,AnekaGolf,CargoMBU,PATravel}_Act`
  (per COB) + `CopyFacRetro{Fire,AnekaGolf,CargoMBU,PATravel,Life,EDMLoc}_ACT` + `CopyFacRetro_ACT`
  + `CopyToAllSpreading_ACT` / `CopyToAllLocSpreading_ACT`.
- **Activity produksi Fac Out:** `InsertFacoutProd`, `InsertFacoutProduction`,
  `InsertFacoutProductionEDM`, `InsertFacoutProductionEDMCurr`.
- **When:** `IsFacout` (`ProposalAcceptStatus=4`), `IsFacRetro` (`.IsFacRetro=1`),
  `IsInputFacRetro` (`.IsInputFacRetro=1`).
- **DataTransform:** `SetIndexFacRetroList_DT`, `SetInFacRetroStatus_PreDT`.
- **Harness/Section/RDBList:** `ViewFacretro`, banyak `*FacOut` (Section/FlowAction), `GetOPFacOut_Sql`,
  `CekFacoutProd_Sql`, `GenerateOurRefFacOut_SQL`, `GetFacoutList_SQL` (Endorsment).

Perintah audit: `Get-ChildItem -Recurse "D:\migrasi\RNM\NB FacIn" -Filter *FacRetro* , *FacOut*`.

---

## 2. Logika pemicu penyalinan — [terverifikasi] dari kondisi PreCond

Dari NB `Activity\SetDataFacOutFire_Act.xml` (step-precondition `<pyStepsPreCondParamsWhen>`, isi
langsung dikutip):

```
.TreatyType=="10015"                         (penanda spreading Fac Out)
.CoverageList(1).CoverageBasis!=5  / ==5     (percabangan basis coverage)
Local.isFacOutLoc==1                         (penanda lokasi punya Fac Out)
Local.isFacOut==1                            (penanda ada Fac Out → lakukan salin)
Param.Status=="UW" && pyWorkPage.OfferFacIn.IsInputFacRetro==1   (gerbang re-input saat UW)
```

Dan dari NB `Activity\CopyToAllSpreading_ACT.xml` + `CopyToAllLocSpreading_ACT.xml`
(`<pyStepsPreCondParamsWhen>`):

```
@contains(.TreatyType,"10015") || @contains(.TreatyName,"SPL") || …
```

**Sumber & target (dari `<pyStepsObjectName>` iterasi):**
- baca: `pyWorkPage.OfferFacIn.LocationList` → `.Property.PropertyItemList` →
  `.CoverageList(1).SpreadingList` (+ `.LayerList`).
- tulis: `.OfferFacIn.FacRetro.LocationList`.

**[terverifikasi] Penanda spreading Fac Out — daftar lengkap dari korpus (dikutip apa adanya):**
- `SetDataFacOutFire_Act` & `CopyFacRetroFire_ACT`: **`.TreatyType=="10015"`** (kesetaraan ketat).
- `CopyToAllSpreading_ACT` & `CopyToAllLocSpreading_ACT`:
  **`@contains(.TreatyType,"10015") || @contains(.TreatyName,"SPL") || @contains(.TreatyType,"10007")`**.

> Dua konteks berbeda, **jangan digabung**: penyalinan khusus ke FacRetro memakai `10015` ketat;
> penyalinan spreading umum (`CopyToAll*`) juga mencakup `10007` dan TreatyName `"SPL"`. Port
> masing-masing apa adanya.

**[terverifikasi] Struktur di dalam FacRetro:** dari `CopyFacRetroFire_ACT` —
`FacRetroList(i).LocationList(j).Property.PropertyItemList(k).CoverageList(1).FacOutObjectList(1).Rate`
dan `.FacOutTSI`. Jadi tiap coverage yang diretro punya **`FacOutObjectList`** berisi `Rate` + `FacOutTSI`.

**Arti kode (dari keterangan work owner, dicatat K-054):** `10015` = **FACOUT**, `10007` = **ORS**
(Own Retention Share). Belum ditemukan file korpus yang memetakan angka→label secara eksplisit
(`M_TREATY_IN` = `(ID, JSONDATA CLOB)`, labelnya di dalam JSON). Status: **[keputusan work owner]**,
bukan [terverifikasi korpus]. TreatyName `"SPL"` = **[pertanyaan terbuka]**.

**[terverifikasi] `CoverageBasis` — daftar nilai lengkap** dari `DDL\CoverageBasis.xml`
(`pyPromptTableList`): `1`=Sum Insured Basis · `2`=First Loss Basis · `3`=EML/PML Basis ·
`4`=Sub Limit Basis · **`5`=Layering Basis**. Jadi `CoverageBasis==5` di gerbang FacOut = basis
**Layering** (berkaitan `LayerList`).

### 2.1 Dua struktur berbeda — [terverifikasi]

| Struktur | Diisi oleh | Bentuk | Peran (dugaan dari pola) |
| --- | --- | --- | --- |
| `OfferFacIn.FacRetro` | `SetDataFacOutFire_Act` step 7 (`FacRetro.LocationList(<APPEND>)`) | tunggal, punya `.LocationList` | **staging** saat UW Accept |
| `OfferFacIn.FacRetroList` | `CopyFacRetroFire_ACT` step 1 (`FacRetroList(counterFacRetro)…`) | **list berindeks**, tiap entri punya `.LocationList` | **daftar final** retrosesi (bisa >1 reasuradur) |

Gerbang `CopyFacRetroFire_ACT`: `Param.status=="CopyAll"`, `.TreatyType=="10015"`, `local.Facout=="0"`,
`IsFire`. Flag disetel: `IsFacRetro`, `IsInputFacRetro`. **[dugaan]** peran staging→list belum
dipastikan urutannya (UI work owner mengalahkan XML); yang [terverifikasi] hanya keberadaan keduanya.

**[UI work owner mengalahkan XML]** Urutan langkah pasti (kapan `isFacOut` di-set 1, iterasi ganda
Location×PropertyItem×Coverage×Spreading) tidak sepenuhnya terserialisasi jelas; yang dikutip di atas
adalah kondisi yang terbaca, bukan rekonstruksi alur penuh.

---

## 3. Rantai yang terbukti (peta)

> ⛔ **Baris pertama peta lama DIBATALKAN** — ~~`UW Accept (ProposalAcceptStatus = 4 → When IsFacout)`~~.
> Lihat §0.1. Peta di bawah ini sudah memakai pemicu yang benar.

```
UW Accept
   │   (TANPA gerbang When di titik masuk — penyaringan ada di precondition per COB)
   │
   ├─ Activity\SetValidateDateUW_PostAct
   │     pyStepPageReference = RH_1.pySteps(3) : Call SetDataFacOut_Act
   │        RH_1.pySteps(2) : Property-Remove  <Property> x4 =
   │                            pyWorkPage.OfferFacIn.FacRetro.LocationList
   │                            pyWorkPage.OfferFacIn.FacRetro.PersonList
   │                            pyWorkPage.OfferFacIn.FacRetro.CargoList
   │                            pyWorkPage.OfferFacIn.FacRetro.VehicleList
   │        RH_1.pySteps(3) : Property-Set   .OfferFacIn.IsFacRetro = 0        (RESET)
   │        RH_1.pySteps(4..7) : Call SetDataFacOut{Fire,AnekaGolf,CargoMBU,PATravel}_Act  (per COB)
   │                     └─ bergerbang pyStepsPreCondParamsWhen:
   │                          .TreatyType=="10015" / Local.isFacOut==1 / Local.isFacOutLoc==1
   │                        → Property-Set .OfferFacIn.IsFacRetro = 1   (NYALA)
   │                        → salin Location → OfferFacIn.FacRetro.LocationList
   │
   └─ Flow\InputInwardFacultativeOffer
         Transition92 : pyTaskStatusOrWhen=WHEN, pyTaskWhen=IsFacRetro, pyLikelihood=100
            → shape pyMOName = [Fac Out]
                 pyImplementation = OfferFacRetro      ← BUKAN pySubFlowName
                    → menu/harness ViewFacretro
                    → InsertFacoutProd → InsertFacoutProduction* (produksi Fac Out)
```

> ⛔ **KOREKSI SESI SEBELUMNYA DIBATALKAN — 21 September 2026.**
>
> Sesi sebelumnya menulis di sini:
>
> ~~"Peta lama juga menulis target salin sebagai `OfferFacIn.FacRetro.{Location,Person,Cargo,Vehicle}List`
> untuk langkah 1. Yang di-`Property-Remove` adalah `pyWorkPage.OfferFacIn.FacRetroList`. Kedua
> struktur itu berbeda."~~
>
> **Itu KELIRU, dan rumusan lama DIPULIHKAN.** Sebabnya: koreksi itu menemukan bentuk lain yang
> **juga ada** di korpus, lalu memperlakukannya sebagai bukti bahwa yang lama salah. Menemukan
> alternatif **bukan** bukti pembatalan. Keduanya ada, di **langkah yang berbeda**:
>
> | Fakta | Tag pembawa | Alamat `pyStepPageReference` | Isi |
> | --- | --- | --- | --- |
> | **A** — yang dimaksud rumusan lama | `<Property>` ×4 (parameter metode `Property-Remove`) | **`RH_1.pySteps(2)`** | `pyWorkPage.OfferFacIn.FacRetro.LocationList` · `.PersonList` · `.CargoList` · `.VehicleList` |
> | **B** — yang ditemukan sesi lalu | `<pyStepsObjectName>` | **`RH_1.pySteps(8).pySteps(1)`** | `pyWorkPage.OfferFacIn.FacRetroList` |
>
> ⛔ **B bukan pengganti A.** Keduanya berbeda **langkah**, berbeda **tag**, dan berbeda **peran**:
> `<Property>` adalah **parameter metode `Property-Remove`** (apa yang dihapus), sedangkan
> `<pyStepsObjectName>` **`[dugaan]` berperan sebagai sasaran iterasi langkah** (apa yang dilewati satu
> per satu) — ⚠️ label "sasaran iterasi" **diturunkan dari `[terverifikasi]` menjadi `[dugaan]`**;
> yang `[terverifikasi]` hanyalah: **metodenya kosong**, **objeknya `OfferFacIn.FacRetroList`**, dan
> ia punya **8 langkah anak**. Arti tag itu sendiri **`[di luar korpus]`** —
> `RH_1.pySteps(8).pySteps(1)` ber-`pyStepsActivityName` **kosong**, jadi ia tidak menghapus apa pun.
>
> 📌 `[terverifikasi]` Di dalam iterasi itu memang ada `Property-Remove` kedua, di
> `RH_1.pySteps(8).pySteps(1).pySteps(1)`, dengan `<Property>` ×4 **berjalur relatif**: `.CargoList`,
> `.LocationList`, `.VehicleList`, `.PersonList`. Jadi **kedua struktur sama-sama dibersihkan** —
> `FacRetro` di langkah 2, dan tiap entri `FacRetroList` di dalam iterasi langkah 8.

`[terverifikasi]` **Penomoran alamat dapat ditetapkan.** Arah tag sudah diuji empiris pada tiga
activity (32/32 `rowdata`): `<pyStepPageReference>` selalu muncul **sesudah** `<pyStepsActivityName>`
di dalam `rowdata` yang sama, dan keduanya **anak langsung `rowdata` itu** — sehingga keanggotaan
bersifat struktural, bukan tekstual. Perinciannya di `10-audit\06-tinjau-ulang-suntingan-facout.md` §3.

⚠️ **Penomoran "langkah N" pada rumusan lama ≠ `pySteps(N)`.** `[terverifikasi]` Rumusan lama memakai
**nomor entri `pyStepsActivityName` yang tidak kosong** (14 entri), sedangkan alamat memakai indeks
`pySteps` (17 simpul, **tiga** di antaranya ber-metode kosong: `pySteps(1)`, `pySteps(8)`, dan
`pySteps(8).pySteps(1)` — **17 − 3 = 14** entri bernama). Karena itu "step 1 Property-Remove" pada
rumusan lama = entri ke-1 = alamat **`RH_1.pySteps(2)`**, dan "step 3-6" = entri ke-3…6 = alamat
**`RH_1.pySteps(4..7)`**. **Isi klaimnya benar; hanya konvensi penomorannya berbeda.**

`[terverifikasi]` Strukturnya **empat**, bukan dua: `OfferFacIn.FacRetro` (staging tunggal) ·
`OfferFacIn.FacRetroList` (berindeks) · `OfferFacIn.FacRetroDetails` ·
`pyWorkPage.Policy.FacOfferList(n)`. Rinciannya di `10-audit\03-struktur-facretro-dan-populasi.md`
dan tiket `05-tickets\facout\F05-empat-struktur-staging.md`.

---

## 4. Rancangan target: menu Fac Out TERPISAH tapi TERHUBUNG

Prinsip: Fac Out adalah **retrocession keluar** — porsi risiko yang RNM teruskan lagi ke reasuradur
lain. Data awalnya **diturunkan** dari OfferFacIn (NB/RNW/EDM), maka menu terpisah tetapi berbagi akar.

### 4.1 Pemisahan modul (arsitektur Go, sejalan CLAUDE.md §4.2 & §5)
```
internal/services/
├── faccase/        (NB/RNW/EDM — sudah dirancang)
├── acceptance/     (tangga akseptasi — Seam 2)
├── premium/        (Seam 3)
├── production/     (InsertJson/flat inward)
└── facout/         ← BARU: retrocession keluar
      ├── derive.go     (SetDataFacOut: turunkan FacRetro dari OfferFacIn saat Accept)
      ├── offer.go      (flow OfferFacRetro: tinjauan/penawaran Fac Out)
      └── production.go (InsertFacoutProduction*)
```
Ketergantungan tetap searah: `facout` membaca hasil `faccase`/`acceptance`, tidak sebaliknya.

### 4.2 Titik sambung (seam) yang [terverifikasi] perlu
- **Sambungan 1 — derivasi.** ~~Dipicu saat akseptasi mencapai kondisi Fac Out
  (`ProposalAcceptStatus=4` + ada spreading `TreatyType~"10015"`).~~ ⛔ **Rumusan pemicu itu
  DIBATALKAN — lihat §0.1.** Penggantinya: `[terverifikasi]` dipanggil **tanpa syarat** dari
  `SetValidateDateUW_PostAct` `RH_1.pySteps(3)` saat UW Accept; penyaringan Fac Out terjadi di
  **precondition per COB** di dalam `SetDataFacOut{Fire,…}_Act`, dan hasilnya adalah flag
  `.IsFacRetro` (0 → 1). `ProposalAcceptStatus` **tidak ikut** menggerbangi jalur ini.
  Fungsi murni: `DeriveFacRetro(offer OfferFacIn) FacRetro` — menyalin
  Location/PropertyItem/Coverage/Spreading yang ber-treaty Fac Out ke agregat `FacRetro`. Ini port
  dari `SetDataFacOut*_Act`.
- **Sambungan 2 — penanda spreading Fac Out.** Predikat `IsFacOutSpreading(s Spreading) bool` =
  `contains(s.TreatyType,"10015") || contains(s.TreatyName,"SPL") || …`. **Port apa adanya**;
  daftar lengkap alternatif `||` harus dibaca utuh dari `CopyToAllSpreading_ACT` sebelum implementasi
  (di dokumen ini baru sebagian terkutip). Masuk registry predikat (K-050).
- **Sambungan 3 — status Fac Out.** ~~`IsFacout` (`ProposalAcceptStatus=4`),~~ `IsFacRetro`
  (`.IsFacRetro=1`), `IsInputFacRetro` (`.IsInputFacRetro=1`) → **dua** predikat di registry.
  ⛔ **`IsFacout` dikeluarkan dari daftar ini** (§0.1): ia tidak terpanggil di jalur Fac Out, dan
  isinya menguji **banding**. Porting-nya sudah ditangani tiket New Business
  `05-tickets\09-predikat-sikap-khusus.md` di bawah K-019 — **jangan diulang** di lingkup Fac Out.

  ⚠️ **`IsFacRetro` — ekspresi tersimpan yang berlaku, bukan label grid.** `[terverifikasi]`
  Yang dieksekusi: `<pyLogic>` = `A` → `<pyConditionValue1>` = `compareTwoValues(.IsFacRetro,"=",1)`.
  Di `<pyNestedConditions>` `rowdata(1)` masih ada `<pyConditionString>` berbunyi
  `pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName = ReasFacInGroupLeader` — itu **sisa
  dari rule asalnya**: `<pzOriginalInstanceKey>` = `RULE-OBJ-WHEN ASM-FW-GISFW-DATA ISUW #…` dan
  `<pyJavaClassName>` = `Rule_Obj_When_…_IsFacRetro_When_…`. Rule ini **salinan `IsUW` yang grid
  kondisinya tidak dibersihkan**. Sikap portingnya **sama dengan `IsOfferFacIn`** (K-002, tiket NB
  `05-tickets\09-predikat-sikap-khusus.md`): **ekspresi tersimpan yang berlaku**, label grid
  diabaikan.

### 4.3 Data — [terverifikasi] `FACOUTPRODUCTION` sudah flat & sudah ada

**Temuan besar: produksi Fac Out TIDAK butuh tabel flat baru.** DDL `FACOUTPRODUCTION.txt` (65 kolom,
[terverifikasi]) sudah flat, multi-COB, dan mempertahankan pasangan `_MENJADI`/`_SELISIH`:

- Kunci & penaut: `IDPEGA VARCHAR2(150)`, `POLICYNO`, `NOENDORS`, **`OBJECTNO_FACIN`** (tautan balik ke
  objek Fac In induk — inilah "keterhubungan" yang dijaga).
- Reasuradur tujuan: `REINSURER_ID`, `REINSURER_NAME`.
- Porsi retro (semua `NUMBER`): `TSIRNM`, `TSISPREADED`, `PCTOFFERED`, `SHAREOFFERED`, `OBJECTPREMI`,
  `RICOMM`, `COMMISION`, `RATE`, `PRORATE`.
- Pasangan `_MENJADI`/`_SELISIH` (EDM): `SHAREOFFERED_SELISIH`, `OBJECTPREMI_SELISIH`,
  `COMMISION_SELISIH`, `PREMI_COVERAGE_MENJADI/_SELISIH`, `COMMISION_COVERAGE_MENJADI/_SELISIH`,
  `RATE_COVERAGE`.
- Kolom per-COB: Fire (`OCCUPATION`,`OBJECTNAME`,`ZIPCODE`,`ROADNAME`), Marine
  (`SHIPID`,`FROMRUTE`,`TORUTE`,`SAILDATE`,`PACKINGNOTE`), Vehicle
  (`LICENSEPLATE`,`ENGINENUMBER`,`BRANDNAME`,`MODELNAME`), PA (`CLASSPA`,`DOB` — PII).

> [terverifikasi] Uang di `FACOUTPRODUCTION` = `NUMBER`, konsisten CLAUDE.md §4.1. `CURRENCY`/
> `CURRENCYID` disediakan → tetap pasangkan Money{Amount,Currency}.

**Yang perlu dimodelkan (belum ada di DDL):** hanya **staging JSON** `FacRetro`/`FacRetroList` +
`FacOutObjectList{Rate,FacOutTSI}` (bagian dari agregat OfferFacIn saat proses UW). Ini masuk rancangan
flat §08 sebagai tabel `FLAT_FACRETRO_*` paralel bila JSON_POLIS diganti flat; ber-FK ke `IDPEGA`+
`PRODKE` induk.

**[terverifikasi] Query daftar Fac Out** — `Endorsment\RDBList\GetFacoutList_SQL.xml` tag
`pyBrowseSQL`: `select IDPEGA as HASIL1, policyno HASIL2 from facoutproduction where policyno =
{InputData.CARI17}`. Jadi daftar Fac Out dibaca dari tabel **`facoutproduction`** berdasar `policyno`.

**[terverifikasi] Insert produksi = SQL INSERT langsung, BUKAN stored procedure.**
`DDL\InsertTreatyProd_Sql.xml` (`pxObjClass=Rule-Connect-SQL`, tag `<pyBrowseSQL>`):
```sql
BEGIN
  INSERT INTO FACOUTPRODUCTION (IDPEGA, POLICYNO, ... 65 kolom ..., OBJECTNO_FACIN)
  VALUES ({DataIN.CARI1}, {DataIN.CARI2}, ..., {DataIN1.CARI14});
  COMMIT;
END;
```
Jadi Pega menjalankan INSERT langsung ke `FACOUTPRODUCTION` dengan parameter `{DataIN.CARIn}` /
`{DataINCargo.*}` / `{DataINMBU.*}` / `{DataINPA.*}` (per COB). **Tidak butuh ALL_SOURCE** — SQL lengkap
di korpus.

> ⛔ **KOREKSI 21 September 2026.** Dua contoh pemetaan yang pernah tertulis di paragraf ini —
> `TSIRNM`←`{DataIN.CARI23}` dan `SHAREOFFERED_SELISIH`←`{DataIN.CARI36}` — **KELIRU**, dan
> **dibatalkan**. Keduanya berasal dari pembacaan yang memecah daftar `VALUES` pada setiap koma,
> sehingga koma di dalam `To_date(x, 'format')` ikut terbaca sebagai pemisah dan seluruh pemasangan
> bergeser. Yang benar: `{DataIN.CARI23}`→`CURRENCY` (kolom 42) dan `{DataIN.CARI36}`→`RATE_COVERAGE`
> (kolom 57); `TSIRNM` diisi `{DataIN.CARI30}` dan `SHAREOFFERED_SELISIH` diisi `{DataIN.CARI33}`.
> **Penggantinya adalah tabel lengkap di bawah ini** — satu geseran kolom = korupsi data produksi
> yang senyap, jadi tabel ini yang dipakai, bukan contoh.

#### 4.3.1 Pemetaan kolom→parameter — **65 baris lengkap** `[terverifikasi]`

Sumber: `DDL\InsertTreatyProd_Sql.xml`, tag `<pyBrowseSQL>`. Tipe kolom dari `DDL\FACOUTPRODUCTION.txt`.
`[terverifikasi]` **65 kolom seimbang dengan 65 nilai**, diukur dengan pemecah sadar-kurung.

Perintah audit:

```powershell
$c   = [IO.File]::ReadAllText('D:\migrasi\RNM\DDL\InsertTreatyProd_Sql.xml')
$sql = [System.Net.WebUtility]::HtmlDecode([regex]::Match($c,'(?s)<pyBrowseSQL>(.*?)</pyBrowseSQL>').Groups[1].Value)
$colBlock = [regex]::Match($sql,'(?s)FACOUTPRODUCTION\s*\((.*?)\)\s*VALUES').Groups[1].Value
$valBlock = [regex]::Match($sql,'(?s)VALUES\s*\((.*)\)\s*;').Groups[1].Value
function Split-TopLevel([string]$s){                      # WAJIB sadar-kurung
  $o=New-Object Collections.Generic.List[string]; $d=0; $b=''
  foreach($ch in $s.ToCharArray()){
    if($ch -eq '('){ $d++; $b+=$ch } elseif($ch -eq ')'){ $d--; $b+=$ch }
    elseif($ch -eq ',' -and $d -eq 0){ $o.Add($b); $b='' } else { $b+=$ch } }
  if($b.Trim() -ne ''){ $o.Add($b) }
  $o | ForEach-Object { ($_ -replace '\s+',' ').Trim() } | Where-Object { $_ -ne '' } }
$cols = Split-TopLevel $colBlock; $vals = Split-TopLevel $valBlock
"KOLOM=$($cols.Count) NILAI=$($vals.Count)"               # -> KOLOM=65 NILAI=65
for($i=0;$i -lt $cols.Count;$i++){ "{0,3}  {1,-28} <- {2}" -f ($i+1),$cols[$i],$vals[$i] }
```

| # | Kolom `FACOUTPRODUCTION` | Tipe Oracle | Parameter Pega |
| ---: | --- | --- | --- |
| 1 | `IDPEGA` | `VARCHAR2(150)` | `{DataIN.CARI1}` |
| 2 | `POLICYNO` | `VARCHAR2(50)` | `{DataIN.CARI2}` |
| 3 | `GROUPPANEL` | `VARCHAR2(10)` | `{DataIN.CARI3}` |
| 4 | `REINSURER_ID` | `VARCHAR2(15)` | `{DataIN.CARI4}` |
| 5 | `REINSURER_NAME` | `VARCHAR2(150)` | `{DataIN.CARI5}` |
| 6 | `TGL_PRINT` | `DATE` | `To_date({DataIN.CARI6}, 'DD/MM/YYYY HH24:MI:SS')` |
| 7 | `RISTARTPERIOD` | `DATE` | `To_date({DataIN.CARI8}, 'DD/MM/YYYY HH24:MI:SS')` |
| 8 | `RIENDPERIOD` | `DATE` | `To_date({DataIN.CARI9}, 'DD/MM/YYYY HH24:MI:SS')` |
| 9 | `START_DATE` | `DATE` | `To_date({DataIN.CARI11}, 'DD/MM/YYYY HH24:MI:SS')` |
| 10 | `END_DATE` | `DATE` | `To_date({DataIN.CARI12}, 'DD/MM/YYYY HH24:MI:SS')` |
| 11 | `RISLIPNO` | `VARCHAR2(600)` | `{DataIN.CARI7}` |
| 12 | `NOENDORS` | `VARCHAR2(100)` | `{DataIN.CARI13}` |
| 13 | `PACKINGID` | `VARCHAR2(200)` | `{DataINCargo.CARI1}` |
| 14 | `PACKINGNOTE` | `VARCHAR2(900)` | `{DataINCargo.CARI2}` |
| 15 | `GOODNOTE` | `VARCHAR2(500)` | `{DataINCargo.CARI3}` |
| 16 | `TRADINGNOTE` | `VARCHAR2(500)` | `{DataINCargo.CARI4}` |
| 17 | `SHIPID` | `VARCHAR2(100)` | `{DataINCargo.CARI6}` |
| 18 | `FROMRUTE` | `VARCHAR2(500)` | `{DataINCargo.CARI7}` |
| 19 | `TORUTE` | `VARCHAR2(500)` | `{DataINCargo.CARI8}` |
| 20 | `SAILDATE` | `DATE` | `To_date({DataINCargo.CARI9}, 'DD/MM/YYYY HH24:MI:SS')` |
| 21 | `CONVEYANCENOTE` | `VARCHAR2(500)` | `{DataINCargo.CARI5}` |
| 22 | `OBJECTNO` | `VARCHAR2(50)` | `{DataIN1.CARI1}` |
| 23 | `ZIPCODE` | `VARCHAR2(50)` | `{DataIN1.CARI2}` |
| 24 | `PROVINCE` | `VARCHAR2(100)` | `{DataIN1.CARI3}` |
| 25 | `CITY` | `VARCHAR2(100)` | `{DataIN1.CARI4}` |
| 26 | `DISTRICT` | `VARCHAR2(100)` | `{DataIN1.CARI5}` |
| 27 | `RW` | `VARCHAR2(1000)` | `{DataIN1.CARI6}` |
| 28 | `ADDRESS` | `VARCHAR2(4000)` | `{DataIN1.CARI7}` |
| 29 | `BUILDINGNO` | `VARCHAR2(30)` | `{DataIN1.CARI8}` |
| 30 | `ROADNAME` | `VARCHAR2(4000)` | `{DataIN1.CARI9}` |
| 31 | `OBJECTNAME` | `VARCHAR2(4000 CHAR)` | `{DataIN1.CARI10}` |
| 32 | `OCCUPATION` | `VARCHAR2(4000)` | `{DataIN1.CARI11}` |
| 33 | `OBJECTITEM` | `VARCHAR2(4000)` | `{DataIN1.CARI12}` |
| 34 | `OBJECTITEMID` | `VARCHAR2(100)` | `{DataIN1.CARI13}` |
| 35 | `BRANDNAME` | `VARCHAR2(500 CHAR)` | `{DataINMBU.CARI2}` |
| 36 | `LICENSEPLATE` | `VARCHAR2(50)` | `{DataINMBU.CARI3}` |
| 37 | `TYPENAME` | `VARCHAR2(50)` | `{DataINMBU.CARI4}` |
| 38 | `MODELNAME` | `VARCHAR2(500)` | `{DataINMBU.CARI5}` |
| 39 | `ENGINENUMBER` | `VARCHAR2(100)` | `{DataINMBU.CARI6}` |
| 40 | `CLASSPA` | `VARCHAR2(100)` | `{DataINPA.CARI1}` |
| 41 | `DOB` | `VARCHAR2(100)` | `{DataINPA.CARI2}` |
| 42 | `CURRENCY` | `VARCHAR2(10)` | `{DataIN.CARI23}` |
| 43 | `CURRENCYID` | `VARCHAR2(10)` | `{DataIN.CARI29}` |
| 44 | `TSIRNM` | `NUMBER(20,4)` | `{DataIN.CARI30}` |
| 45 | `TSISPREADED` | `NUMBER(20,4)` | `{DataIN.CARI31}` |
| 46 | `PCTOFFERED` | `NUMBER(20,4)` | `{DataIN.CARI32}` |
| 47 | `SHAREOFFERED` | `NUMBER(20,4)` | `{DataIN.CARI24}` |
| 48 | `OBJECTPREMI` | `NUMBER(20,4)` | `{DataIN.CARI25}` |
| 49 | `RICOMM` | `NUMBER(20,4)` | `{DataIN.CARI26}` |
| 50 | `COMMISION` | `NUMBER(20,4)` | `{DataIN.CARI27}` |
| 51 | `RATE` | `NUMBER(25,20)` | `{DataIN.CARI28}` |
| 52 | `SHAREOFFERED_SELISIH` | `NUMBER(20,4)` | `{DataIN.CARI33}` |
| 53 | `OBJECTPREMI_SELISIH` | `NUMBER(20,4)` | `{DataIN.CARI34}` |
| 54 | `COMMISION_SELISIH` | `NUMBER(20,4)` | `{DataIN.CARI35}` |
| 55 | `TYPEFACULTATIVE` | `VARCHAR2(50)` | `{DataIN.CARI39}` |
| 56 | `PRORATE` | `NUMBER(25,20)` | `{DataIN.CARI40}` |
| 57 | `RATE_COVERAGE` | `NUMBER(20,4)` | `{DataIN.CARI36}` |
| 58 | `PREMI_COVERAGE_MENJADI` | `NUMBER(20,8)` | `{DataIN.CARI37}` |
| 59 | `PREMI_COVERAGE_SELISIH` | `NUMBER(20,8)` | `{DataIN.CARI38}` |
| 60 | `COVERAGE_NAME` | `VARCHAR2(1000)` | `{DataIN.CARI42}` |
| 61 | `COVERAGE_ID` | `VARCHAR2(100)` | `{DataIN.CARI41}` |
| 62 | `COMMISION_COVERAGE_PCT` | `NUMBER(20,8)` | `{DataIN.CARI43}` |
| 63 | `COMMISION_COVERAGE_MENJADI` | `NUMBER(20,8)` | `{DataIN.CARI44}` |
| 64 | `COMMISION_COVERAGE_SELISIH` | `NUMBER(20,8)` | `{DataIN.CARI45}` |
| 65 | `OBJECTNO_FACIN` | `VARCHAR2(10)` | `{DataIN1.CARI14}` |

**Pemetaannya TIDAK berurutan.** Perhatikan kolom 6–11 (`CARI6, 8, 9, 11, 12`, lalu `CARI7` menyusul
di `RISLIPNO`), kolom 17–21 (`CARI6, 7, 8, 9`, lalu `CARI5` di `CONVEYANCENOTE`), dan kolom 42–57
(`CARI23, 29, 30, 31, 32`, lalu `CARI24…28`, lalu `CARI33…35`, lalu `CARI39, 40`, lalu `CARI36`).
Urutan ini **diport apa adanya**; jangan "dirapikan".

**Parameter yang TIDAK terpakai** `[terverifikasi]` — sebelas seluruhnya:

| Halaman | Dipakai | Tidak terpakai |
| --- | --- | --- |
| `DataIN` | 35 dari rentang 1..45 | `CARI10`, `CARI14`…`CARI22` (**sepuluh**) |
| `DataIN1` | 14 dari 1..14 | — |
| `DataINCargo` | 9 dari 1..9 | — |
| `DataINMBU` | 5 dari 1..6 | `CARI1` |
| `DataINPA` | 2 dari 1..2 | — |

⚠️ `{Halaman.CARIn}` adalah **substitusi teks** ke dalam SQL, bukan bind variable. Di sistem baru ini
**wajib** menjadi parameter terikat. Itu perubahan mekanisme penyampaian nilai, **bukan** perubahan
perilaku — nilai yang tertulis ke kolom tetap sama, sehingga tidak melanggar `CLAUDE.md` §1.

> **Koreksi catatan sebelumnya:** `InsertFacoutProduction.xml` bukan pemanggil SP DB; insert nyata =
> Connect-SQL `InsertTreatyProd_Sql` di atas. Item "isi SP" DICABUT dari paket DBA.

### 4.4 Alur kerja Fac Out — [terverifikasi] tangga persetujuan retro SENDIRI (4 tingkat di EDM)

Flow `OfferFacRetro` punya tangga peran retrosesi **terpisah** dari tangga akseptasi Fac In.

> ⛔ **PEMBATAS LINGKUP (ditambahkan 21 September 2026).** Tangga **empat tingkat bukan keadaan
> seragam lintas siklus.** Ada **empat salinan** `OfferFacRetro`, dan **keempatnya ber-`pyRuleSetVersion`
> IDENTIK `01-01-95`** — sehingga perbedaannya **tidak terlihat** dari nomor versi:
>
> | Salinan | Ukuran | Assignment | Gateway | Workbasket | Commit |
> | --- | ---: | ---: | ---: | --- | --- |
> | `NB FacIn\Flow\OfferFacRetro.xml` | 134.214 B | 3 | 4 | `ReasFacOutAdmin` + `ReasFacOutHead` (**2 tingkat**) | 2026-08-07 |
> | `RNW Fac In\Flow\OfferFacRetro.xml` | 134.214 B | 3 | 4 | idem NB — **identik byte** | 2026-08-07 |
> | `Endorsment Fac In\Flow\OfferFacRetro.xml` | 170.966 B | 5 | 6 | Admin/Head/GroupLeader/TechnicalDirector (**4 tingkat**) | 2026-08-07 |
> | `DDL\OfferFacRetro.xml` (ekspor ulang work owner) | 171.308 B | 5 | 6 | Admin/Head/GroupLeader/TechnicalDirector (**4 tingkat**) | 2026-08-31 |
>
> **Yang [terverifikasi]:** tangga 4 tingkat ada pada salinan **Endorsment** dan pada **ekspor ulang
> `DDL\`**. Salinan **NB dan RNW di korpus hanya memuat 2 tingkat**.
>
> **Yang [pertanyaan terbuka]:** apakah NB/RNW di **produksi** juga sudah 4 tingkat. Ekspor `DDL\`
> lebih baru (2026-08-31) tetapi **tidak menyatakan siklus asalnya**, jadi ia tidak dapat dipakai
> untuk menyimpulkan keadaan NB/RNW. K-055 menyebut ekspor lama NB/RNW "digantikan" — **itu keputusan
> work owner, bukan temuan korpus**, dan berlaku sebagai keputusan (`PANDUAN-KERJA` §1).
>
> ⚠️ Kasus ini sekaligus **bukti** bahwa membandingkan `pyRuleSetVersion` tidak cukup untuk mendeteksi
> drift antar folder — lihat bahan K-057.

**[terverifikasi] EDM `Endorsment Fac In\Flow\OfferFacRetro.xml` — 4 tingkat + gerbang group:**
shape `pyMOName` memuat `ReasFacOutAdmin`, `ReasFacOutHead`, `ReasFacOutGroupLeader`,
`ReasFacOutTechnicalDirector`, plus `Is it group?`/`IsGroup`, beberapa `Accept?`, `IsRISlip`,
`OfferFacOut`, `InsertFacoutProd`, `PRINT R/I SLIP`, `Send Email Print RI Slip`.
Urutan yang dinyatakan work owner: **Admin → Head → GroupLeader → TechnicalDirector**.

**[terverifikasi] Ekspor ulang `DDL\OfferFacRetro.xml` (2026, dari work owner) — 4 tingkat penuh:**
keempat peran `ReasFacOutAdmin`(13) + `ReasFacOutHead`(11) + `ReasFacOutGroupLeader`(9) +
`ReasFacOutTechnicalDirector`(9), plus `Is it group?`/`IsGroup`, `Accept?`, `IsRISlip`,
`PRINT R/I SLIP`, `InsertFacoutProd`, routing `ToWorkbasket`/`WorkBasket`. **Menutup selisih**: tangga
retro 4 tingkat kini terkonfirmasi (bukan hanya EDM). Ekspor lama NB/RNW `Flow\OfferFacRetro.xml`
(2 tingkat) = versi lama, digantikan ekspor ulang ini.

> Urutan transisi persis antar `Accept?` = [UI work owner mengalahkan XML]; yang [terverifikasi]
> keberadaan keempat peran + percabangan group + Accept + routing workbasket.

### 4.4.1 Generator nomor slip retro — [terverifikasi] `DDL\GENERATE_FACRETRO_NO.txt`

Fungsi Oracle `POOLDATA.GENERATE_FACRETRO_NO(FacType, BusinessType, RetroMonth, RetroYear)`:
- Format keluaran: `RNM-{FacCode}{BusinessCode}.{RetroMonth}.{RetroYear}.{LPAD(FACRETRO_SEQ.NEXTVAL,5,'0')}`.
- `FacCode = 'Y'` bila `FacType=='FAKULTATIF RETROSESI'`; `'F'` bila `'FAKULTATIF INWARD'`.
- Sequence Oracle `FACRETRO_SEQ`.

> [terverifikasi] Fac Out (retrosesi) memakai prefix **`RNM-Y...`**, berbeda dari Fac In `RNM-F...`.
> Port fungsi ini apa adanya (termasuk pemetaan `FacType`→`FacCode`).

- **Konsekuensi rancangan:** `services/facout` memakai **tangga peran retro sendiri** (4 langkah),
  bukan memakai ulang mesin `acceptance` limit-berjenjang Fac In (Seam 2). Ini seam Fac Out tersendiri.

### 4.5 Keterhubungan menu (UI React)
- Menu **Fac In** (NB/RNW/EDM) dan menu **Fac Out** terpisah di navigasi.
- Saat Accept memicu derivasi, sistem membuat **item Fac Out** yang tertaut ke objek Fac In induk lewat
  **`OBJECTNO_FACIN`** + `IDPEGA`/`POLICYNO`/`NOENDORS` ([terverifikasi] kolom `FACOUTPRODUCTION`).
- UI Fac Out menampilkan daftar dari `GetFacoutList_SQL` (Endorsment) / `GetOPFacOut_Sql`
  ([terverifikasi] berkas RDBList ada; isi SQL butuh ALL_SOURCE).

---

## 5. Status discovery Fac Out & sisa pertanyaan

**Sudah tuntas (langkah discovery ini):**
- [terverifikasi] Rantai pemicu Accept → SetDataFacOut → salin ke FacRetro/FacRetroList.
- [terverifikasi] Daftar lengkap penanda spreading Fac Out (`10015` ketat; `10015`/`10007`/`"SPL"` umum).
- [terverifikasi] Dua struktur staging: `FacRetro` (tunggal) + `FacRetroList` (berindeks) + `FacOutObjectList{Rate,FacOutTSI}`.
- [terverifikasi] Alur `OfferFacRetro` = Admin→Head + RI Slip + insert produksi (bukan tangga akseptasi berjenjang).
- [terverifikasi] `FACOUTPRODUCTION` = tabel produksi flat 65 kolom, sudah ada, dengan `_SELISIH`/`_MENJADI`.

**Tambahan tuntas:**
- [terverifikasi] `GetFacoutList_SQL` = `select IDPEGA, policyno from facoutproduction where policyno=…`.
- ~~[terverifikasi] Insert via SP `INSERTFACOUTPRODUCTION` (dipanggil `InsertFacoutProduction.xml`).~~
  ⛔ **DIBATALKAN oleh K-056** (dan oleh §4.3 dokumen ini). Label `[terverifikasi]`-nya **keliru**:
  `INSERTFACOUTPRODUCTION` adalah nama request/connector, **bukan** stored procedure basis data.
  **Penggantinya:** insert dijalankan Connect-SQL `DDL\InsertTreatyProd_Sql.xml`
  (`INSERT INTO FACOUTPRODUCTION`, 65 kolom — tabel penuh di §4.3.1).
- [terverifikasi] Tangga retro EDM 4 tingkat (Admin→Head→GroupLeader→TechnicalDirector).
  ⚠️ Berlaku untuk salinan **Endorsment** dan **`DDL\`** saja — lihat pembatas lingkup di §4.4.

**Sisa [pertanyaan terbuka] — tinggal sedikit:**
1. ~~Isi SP insert~~ → **TUNTAS**: insert = Connect-SQL `InsertTreatyProd_Sql` (SQL lengkap di korpus).
2. ~~`CoverageBasis==5`~~ → **TUNTAS**: = Layering Basis (`CoverageBasis.xml`).
3. ~~Tangga NB/RNW 2 tingkat~~ → **TUNTAS**: ekspor ulang `DDL\OfferFacRetro.xml` = 4 tingkat.
4. `10015`=FACOUT, `10007`=ORS → **[keputusan work owner K-054]**; belum ada file korpus pemeta
   angka→label eksplisit. TreatyName `"SPL"` masih [pertanyaan terbuka].

Discovery Fac Out kini praktis lengkap; hanya makna `"SPL"` dan label master treaty yang tersisa.

**Kesiapan tiket:** seam Fac Out kini jelas — derive (K-053), predikat penanda (K-054), tangga retro
4 tingkat (K-055), produksi ke `FACOUTPRODUCTION` (K-056). Cukup untuk membuat tiket F01…
Tabel produksi tidak memblokir (sudah ada).

> ⛔ **DIBATALKAN 21 September 2026 — dua rumusan di paragraf ini:**
>
> 1. ~~"produksi ke `FACOUTPRODUCTION` **via SP**"~~ — **tidak ada SP**. Insert = Connect-SQL
>    `InsertTreatyProd_Sql`. Sudah dikoreksi pada kalimat di atas; dicatat di sini agar perubahannya
>    terbaca (`PANDUAN-KERJA` §7).
> 2. ~~"yang menunda hanya isi tubuh SP insert (`ALL_SOURCE`)"~~ — **tidak ada yang menunda dari sisi
>    SP**, karena SQL-nya lengkap di korpus. Ini **bukan** koreksi baru: §4.3 dan K-056 sudah
>    mencabutnya, tetapi kalimat ini tertinggal masih berlabel tuntas.
>
> **Yang sesungguhnya menunda Fac Out** (per 21 September 2026) bukan `ALL_SOURCE`, melainkan:
> `RDBList\GetOPFacOut_Sql` membaca `datapega.pc_history_asm_fw_gisfw_work`, yakni tabel internal Pega
> yang menurut `CLAUDE.md` §4.3 **hilang bersama Pega dan tidak dimigrasikan apa adanya** — jadi klep
> ini butuh pengganti, dan penggantinya belum diputuskan. `[terverifikasi]` tag `<pyBrowseSQL>`.

**Keputusan tercatat:** K-053 (menu terpisah), K-054 (penanda), K-055 (tangga retro), K-056 (produksi).
