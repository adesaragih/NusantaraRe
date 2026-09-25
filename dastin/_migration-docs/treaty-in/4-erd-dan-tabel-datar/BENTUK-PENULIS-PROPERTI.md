# Bentuk penulis properti — daftar yang disusun, bukan yang terpikir

**Tanggal:** 24 September 2026
**Sebab:** `pyPropertyTarget` ditemukan **secara kebetulan** setelah sebelas sesi — dan itu berarti
daftar bentuk penulis kami **tidak pernah disusun**, hanya dipakai.

> **`CONTEXT.md` §2.0-i menyuruh mendaftar bentuk penulisnya. Berkas ini adalah daftarnya.**
> Aturan yang menyuruh mendaftar tidak berguna sampai daftarnya benar-benar ada.

---

## 1. Cara menyusunnya — mekanis, tidak bergantung pada apa yang terpikir

Seluruh nama elemen XML di korpus dikeluarkan beserta jumlahnya, lalu disaring pada nama yang
berbentuk penugasan (`Target`, `SetValue`, `Default`, `Assign`, `Post`, `Source`, `Property`,
`Value`).

| | |
|---|---:|
| nama elemen berbeda di korpus | **1.823** |
| kandidat berbentuk penugasan | **174** |
| **bentuk yang benar-benar menulis properti kontrak** | **4** |

---

## 2. ~~Empat~~ **LIMA** bentuk, dan hanya satu yang pernah kami sapu

> ### Bentuk kelima ditemukan 24 September 2026 — dan penyaringan §1 TIDAK dapat menemukannya
>
> ```xml
> <pyStepsCallParams>
>   <CopyFrom>TreatyIntemp</CopyFrom>
>   <CopyInto>TreatyIn.OLDDATA</CopyInto>
> </pyStepsCallParams>
> ```
>
> **Salin halaman.** Sasarannya duduk di **parameter langkah**, bukan di pasangan nama/nilai — jadi
> ia tidak tersaring oleh kata `Target`, `Property`, `Value`, atau `Assign`. **64 langkah** di
> korpus, **4** menyasar halaman potret kontrak.
>
> **Yang mengungkapnya bukan penyaringan melainkan NOL YANG MUSTAHIL:** `OLDDATA` ditampilkan di
> empat berkas layar `…OldData.xml` dan sapuan empat-bentuk mengembalikan **nol penulis**. Sesuatu
> yang ditampilkan di empat layar tidak mungkin tidak punya penulis.
>
> **Pelajaran yang berdiri sendiri:** daftar nama elemen menjawab *"nama mana lagi yang dipakai
> untuk gagasan ini"*; ia **tidak** menjawab *"gagasan ini dapat berbentuk apa lagi"*. Yang menjawab
> yang kedua adalah **hasil yang mustahil** — nol yang seharusnya tidak nol. Uraian lengkapnya
> `7-2-ACTUALVALUE.md` §4.



| Bentuk | Di mana | Nilai | Sasaran berbeda | Pernah disapu? |
|---|---|---:|---:|---|
| `<PropertiesName>` | **Activity** — `Property-Set` | 7.609 | 1.579 | **ya**, sejak awal |
| **`<pyPropertiesName>`** | **Data Transform** | 430 | 114 | **TIDAK** |
| **`<pyPropertyTarget>`** | **Harness / Section** — pemetaan hasil pencarian | 450 | 15 | **TIDAK** |
| `<pyTargetProperty>` | **Report Definition** | 351 | 118 | tidak perlu — lihat §3 |

### 2.1 `pyPropertiesName` — **73 berkas Data Transform tidak pernah ikut tersapu**

Ini lubang terbesar dari keempatnya, dan ia sederhana: penyisir kami mencari `<PropertiesName>`,
sementara *data transform* memakai `<pyPropertiesName>` — **beda awalan `py`**.

Satu jebakan yang harus disaring: **elemen yang sama memuat penugasan DAN kondisi.** Dari 114
sasaran, **33 memuat operator** (`==`, `||`) dan merupakan ekspresi *when*, bukan penugasan.

**Sesudah disaring: 39 sasaran `TreatyIn.*` yang ditulis data transform dan tidak pernah terlihat.**

| Sasaran | Berkas |
|---|---|
| `Position`, `PositionUsername`, `StatusAkseptasi`, `ViewState`, `RevisionState` | `Akseptasi_DT` |
| `ProRatePercent`, `ProRateDays`, `ProRateTotalDays` | `TreatyCalculateProratePct` |
| `Limits(…).Reinstatement_List(…)` | `CalculateReinstatement`, `…Pct`, `ReCalculateReinstatement` |
| **`CommentList(<APPEND>).IsApproved`, `.OperatorName`, `.Date`, `.Suggest`** | **`TreatyInSetEditPre`** — §4 |
| **`ActualValue.EGNPI`**, **`AddendumPremi`**, `EDMEffective`, `EDMState` | `TreatyInSetEditPre`, `TreatyCreateEDM` |
| `Ceding`, `CedingID`, `LeadingReinsSource`, `LeadingReinsSourceID` | `TreatyInSetReinsured` |
| `Exclusions`, `SpecialConditions` | `TreatyInCopyConditions` |
| `TreatyYear`, `Termination` | `TreatyInSetTreatyYear` |
| `FacultativeShare`, `FacultativeShareBrokerage`, `ShareFacultativeReinsurers` | `ShareRetroFacultative_DT` |
| `AccumulationList`, `ReportingPeriod`, `IsEditData`, `Comment`, `Limits` | lima berkas lain |

---

## 3. `pyTargetProperty` — **KOREKSI: bukan penulis halaman kontrak**

Hitungan pertama saya berbunyi *"118 sasaran, 116 berawalan `TreatyIn.` atau `.`"*, dan itu
**menyesatkan**. Sasarannya memang berawalan titik — tetapi titik itu relatif terhadap **kelas
laporannya sendiri**, bukan terhadap halaman kontrak.

Seluruh 18 berkasnya adalah `Browse*_RD.xml` — **Report Definition**, dan `pyTargetProperty` di sana
memetakan **kolom keluaran laporan**: `.ID`, `.TreatyYear`, `.ReinsTypeID`, `.Pct`.

> **Jebakan yang layak dicatat: nama sasaran berawalan titik belum tentu menyasar halaman yang
> sama.** Periksa **jenis aturan** tempat bentuk itu muncul sebelum menghitungnya sebagai penulis.

---

## 4. `pyPropertyTarget` tanpa saringan — dan apa artinya bagi pasangan NAMA + PENGENAL

Diulang **tanpa** saringan `TreatyIn.*`: **15 sasaran berbeda di seluruh korpus**, dan hanya
**empat** yang menyasar halaman kontrak — `CedingID`, `CedingStatusActive`,
`LeadingReinsSourceID`, `SourceStatusActive`. Sisanya menyasar halaman pencarian sementara.

**Jadi angka empat bertahan untuk bentuk ini.** Tetapi kecurigaan tentang **pasangan nama +
pengenal** tidak gugur — ia hanya pindah bentuk: `TreatyInSetReinsured.xml`, sebuah **data
transform**, menulis `Ceding` **dan** `CedingID`, `LeadingReinsSource` **dan**
`LeadingReinsSourceID` berpasangan.

**Pola tanda tangannya nyata**, dan contoh JSON memperlihatkannya berulang:

```
SpreadingType "2023 QS 150M TRT"  ↔  SpreadingTypeID "10236"
TreatyGroup   "MARINE HULL"       ↔  TreatyGroupID   "10010"
ReinsTypeName "QS (R/I)"          ↔  ReinsTypeID     "10004"
Ceding        "ASURANSI …"        ↔  CedingID        "G0000024"
```

### Uji yang mengadili tiap pasangan — satu pertanyaan

> **Bila master berubah besok, nilai di kontrak ini harus ikut berubah atau tidak?**
>
> **Ikut berubah** → ia **salinan**; dibuang, hilir diberi rujukan (ADR-0023, ADR-0041).
> **Tidak** → ia **fakta terbukukan**, dan ia menuntut **penunjuk asal** seperti `ID_SUSUNAN_RETRO`
> (ADR-0036, INV-58).

**Jangan berlebihan ke arah mana pun.** Nama cedant pada saat kontrak disetujui memang harus tetap
meskipun master berubah — itu bukan salinan yang dibuang.

**Keadaan model sekarang, dan ia sudah benar untuk empat pasangan:** §10.1 dan §10.2 **sudah**
membuang `Ceding`, `LeadingReinsSource`, dan `LeadingReinsName`, menyimpan hanya pengenalnya. Yang
**belum** diadili dengan pertanyaan di atas: `TreatyGroup`, `ReinsTypeName`, `SpreadingType`,
`TreatyDescName`, `TreatyGroupName` — seluruhnya di entitas anak.

---

## 5. Yang harus disapu ULANG dengan daftar lengkap ini — **selesai, 24 Sep 2026**

**Kelima** bentuk digabung menjadi satu perkakas, `alat/sapu-penulis-properti.py`, yang menyebut
nama tag yang dipakainya di kepalanya sendiri. Hasilnya: **8.462 baris penugasan, 1.667 sasaran
berbeda, 693 berawalan `TreatyIn.`** — dan 91 sasaran ditolak karena memuat operator (ia ekspresi
*when*).

*(Angka semula 8.398 / 1.638 / 679 adalah hasil empat bentuk; selisihnya tepat 64 langkah salin
halaman, 14 di antaranya menyasar `TreatyIn.*`.)*

| Klaim | Keadaan |
|---|---|
| `CedingStatusActive` / `SourceStatusActive` | **selesai** — terbantah, dan §14.5 justru menguat |
| `ReisuredParticipant` "nol aturan menulisnya" | **selesai** — tetap **nol** atas **kelima** bentuk. Putusan membuang (`SPEC-MODEL-DATA.md` §10.4a) bertahan |
| `ReasTreatyInGroupLeader` "disetel nol kali" | **selesai** — tetap **nol** atas **kelima** bentuk; ia hanya pernah muncul sebagai **kondisi**, dan `PERIKSA-BUTIR-CONTOH.md` §7.10-Z memperlihatkan kenapa |
| masukan versus turunan — 187 / 119 | **selesai** — **tidak berubah**, §5.2 |
| `TREATY_IN` "nol `INSERT`/`UPDATE` langsung" | **tidak terdampak** — bentuknya SQL, bukan penugasan properti |

### 5.1 Berapa banyak yang benar-benar tersembunyi: **18 sasaran, 12 akar properti**

*(Dihitung ulang sesudah bentuk kelima ditambahkan: **tetap 18 dan 12.** Salin halaman menyasar
`TreatyIn.*` di 14 tempat, dan tidak satu pun berimpit dengan ke-18 itu.)*

Dari 679 sasaran `TreatyIn.*`, yang **hanya** punya penulis *data transform* — yaitu yang benar-benar
tak terlihat sapuan lama — ada **18**, atas **12 akar properti**:

| Akar properti | Berkas yang menulisnya |
|---|---|
| `AccumulationList` | `TreatyInAddAccumulation`, `TreatyInDeleteAccumulationLists` |
| `ActualValue.EGNPI` | `TreatyInSetEditPre` |
| `AddendumPremi` | `TreatyInSetEditPre` |
| `CommentList` | `TreatyInSetEditPre`, `TreatyInRemoveLastComment` |
| `EDMEffective` | `TreatyInSetEditPre` |
| `Limits[]`, `Limits[].Reinstatement_List[]` | `SetIndexLayer_DT`, `CalculateReinstatement`, `…Pct`, `ReCalculateReinstatement` |
| `ProRateDays`, `ProRateTotalDays`, `ProRatePercent` | `TreatyCalculateProratePct` |
| `ReportingPeriod` | `TreatyInSetPeriod` |
| `ShareFacultativeReinsurers` | `ShareRetroFacultative_DT` |

Ke-21 sasaran sisanya dari 39 **juga ditulis aktivitas**, jadi penulisnya sudah terlihat sejak awal.
Angka yang berarti bukan 39 melainkan **18**.

### 5.2 Cocokkan ke §10.1 dan §10.2 — **nol berubah, dan berikut hitungannya**

| Pertanyaan | Jawaban |
|---|---:|
| dari 12 akar, berapa yang tercatat **TURUNAN** | **3** — `ProRateDays`, `ProRateTotalDays`, `ProRatePercent` (§4.1) |
| berapa yang tercatat **DIBUANG** (§5) | **0** |
| berapa yang **sudah** tercatat sebagai atribut | **2** di §10.2 — `EDMEffective`, `ReportingPeriod` |
| berapa yang milik **entitas anak**, bukan §10.1/§10.2 | **5** — `AccumulationList`, `CommentList`, `Limits[]`, `ShareFacultativeReinsurers`, `Reinstatement_List[]` |
| berapa yang milik **modul Adjustment** | **2** — `ActualValue`, `AddendumPremi` |

**Ketiga yang tercatat TURUNAN diadili ulang, dan ketiganya TETAP TURUNAN** — tetapi dasarnya
berubah dari terkaan menjadi rumus yang terbaca:

```
ProRateDays      = @DateTimeDifference(EDMEffective,  Termination, 'D')
ProRateTotalDays = @DateTimeDifference(Commencement,  Termination, 'D')
ProRatePercent   = @divide(ProRateDays, ProRateTotalDays) * 100
```

§4.1 berbunyi *"selisih tanggal terhadap periode"* — tidak salah, tetapi tidak dapat diperiksa.
Sekarang ia **`EVIDENCED`**, dan `SPEC-MODEL-DATA.md` §4.1 diperbarui dengan rumusnya.

> **Ketiganya digolongkan turunan karena rumusnya dikenali, BUKAN karena penulisnya tidak
> ditemukan.** Itu sebabnya sapuan yang buta tidak merusaknya. Penggolongan yang berdiri di atas
> *"tidak ditemukan penulis"* tidak ada di antara ke-12 — dan itulah kesimpulan §5.2.

> **Maka 187 / 119 TIDAK BERUBAH.** Bukan karena diperiksa sepintas, melainkan karena ke-12 akar
> ditelusuri satu per satu dan tidak satu pun berpindah golongan.

### 5.3 Satu temuan lain yang lahir dari ke-12 itu — dan ia bukan soal penggolongan

`TreatyInSetEditPre.xml` menyemai **`EDMEffective = TreatyIn.Commencement`** saat addendum dibuat.

Masukkan ke rumus di atas: bila `EDMEffective` **tidak diubah orang**, maka
`ProRateDays == ProRateTotalDays`, dan **`ProRatePercent` selalu 100**. Prorata yang sesungguhnya
hanya muncul bila seseorang **menggeser tanggal berlakunya**.

> Ini instans kelima dari *"cacat bersembunyi di balik nilai bawaan"*, dan bentuknya persis: satu
> rumus yang benar, satu nilai bawaan yang membuat hasilnya tak terbedakan dari tidak memakai rumus
> itu sama sekali.
>
> **Ujinya data, dan ia dapat dijalankan:** berapa baris yang `EDMEffective != Commencement`. Nol
> berarti prorata **belum pernah dipakai** di produksi walaupun mesinnya ada; bukan nol memberi
> angkanya. Ditambahkan sebagai **Uji AJ**.

**Dan `Rule-Declare-Expressions` / `Rule-Declare-Trigger` tetap nol di ekspor (L-10)** — jadi
walaupun keempat bentuk di atas disapu seluruhnya, lubang itu **belum tertutup**. Daftar ini
menutup lubang yang **dapat** ditutup dari ekspor yang ada; L-10 menutup sisanya.
