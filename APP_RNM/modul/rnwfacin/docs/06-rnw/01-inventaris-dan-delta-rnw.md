# Discovery Renewal — Putaran 1: inventaris lengkap dan peta delta NB↔RNW

> **Sumber tunggal:** `D:\migrasi\RNM\RNW Fac In\` — **1.927 berkas `.xml`**, READ-ONLY, dibandingkan
> terhadap `D:\migrasi\RNM\NB FacIn\` (2.083 berkas). `Endorsment Fac In\` **tidak dibaca** (di luar
> fokus, K-030). Korpus Treaty tersendiri **tidak dibaca** (K-005).
>
> Label mengikuti `CLAUDE.md` §3. Setiap angka disertai perintah audit yang menghasilkannya.

---

## Ringkasan temuan yang mengubah rencana

### ⛔ T-1 — RNW bukan "mirip NB". Ia **adalah** NB, minus 176 berkas, plus 20 berkas

`[terverifikasi]` Dari **1.907** berkas yang bernama sama di kedua folder, **1.907 identik
byte-per-byte**. Nol berbeda. Bukan "mirip", bukan "setara secara semantik" — **byte yang sama
persis**.

| Ukuran | Nilai |
| --- | ---: |
| Berkas bersama NB↔RNW | **1.907** |
| — identik **byte-per-byte** | **1.907** (100 %) |
| — berbeda | **0** |
| Hanya ada di RNW | **20** |
| Hanya ada di NB | **176** |

Penjumlahannya menutup sempurna: `1.907 + 20 = 1.927` (RNW) dan `1.907 + 176 = 2.083` (NB).

**Konsekuensi yang mengubah rencana:** pertanyaan "apakah RNW memakai ulang modul NB atau perlu
implementasi tersendiri" **terjawab untuk 1.907 berkas sekaligus — pakai ulang**. Seluruh delta
renewal muat dalam **20 berkas**. Discovery RNW karena itu tidak perlu membaca ulang 1.907 berkas;
mereka sudah didiscovery sebagai NB.

### ⛔ T-2 — Jebakan metodologis #2 **tidak berlaku** untuk pasangan NB↔RNW

`steering\PANDUAN-KERJA.md` §2 mencatat: *"`Get-FileHash` mentah untuk membandingkan antar folder →
buang baris `px*`/`pz*`, metadata ekspor berbeda di setiap berkas."*

`[terverifikasi]` **Itu benar untuk perbandingan yang melibatkan EDM, dan tidak benar untuk NB↔RNW.**
Hash mentah kedua folder cocok pada **1.907 dari 1.907** berkas — metadata ekspornya pun identik.

Ini bukan sekadar kemudahan teknis. Metadata ekspor yang identik berarti kedua folder **diekspor dari
keadaan ruleset yang sama**, menguatkan T-1 dari arah yang berbeda.

⚠️ Jebakan itu **tetap berlaku** untuk NB↔EDM dan RNW↔EDM. Yang dikoreksi hanya cakupannya.

### ⚠️ T-3 — `isApproved` hilang dari RNW tetapi **masih dirujuk dua activity**

`[terverifikasi]` DecisionTable `isApproved` **tidak ada** di `RNW Fac In\DecisionTable\`, tetapi
dirujuk oleh dua berkas yang **ada** di RNW: `Activity\InsertHistoryAkseptasiPega.xml` dan
`Activity\SaveViewSuggest.xml`.

Ini **pola keempat kalinya** di proyek ini — sama bentuknya dengan K-006 (11 activity), K-019
(`IsFacout`), dan K-024 (`JUW_A`): rujukan yang menggantung karena berkas tujuannya tidak ada di
ekspor.

`[pertanyaan terbuka]` Dua kemungkinan, dan **tidak boleh ditebak**:

1. Renewal memang tidak pernah mencapai cabang itu — ketiadaannya disengaja.
2. Ekspor RNW **tidak lengkap** untuk rule itu.

⚠️ Kemungkinan (2) **lebih berbobot di sini daripada di kasus sebelumnya**, justru karena T-1: kedua
folder terbukti diekspor dari keadaan yang sama, sehingga satu berkas yang hadir di satu sisi saja
lebih mencurigakan. **"Hilang dari ekspor" tetap ≠ "usang"** (`CLAUDE.md` §4.5).

### ✅ T-4 — Kelas kerja khusus renewal memang ada

`[terverifikasi]` `ReportDefinition\RenewalList_RD` berkelas **`ASM-FW-GISFW-Work-Renewal`** — kelas
kerja tersendiri untuk renewal, satu-satunya berkas di korpus RNW yang memakainya.

`[dugaan]` Ini menyiratkan renewal punya **antrean kasus tersendiri**, bukan sekadar penanda pada
kasus Fac In biasa. Belum dibuktikan; bahan putaran berikutnya.

### ⚠️ T-5 — Gerbang masuk NB **tidak dipakai renewal sama sekali**

`[terverifikasi]` `IsOfferFacIn` — gerbang masuk siklus NB yang dikunci **K-002** — termasuk 21 rule
`When` yang tidak ada di RNW, dan **dirujuk nol kali** di seluruh korpus RNW.

`[pertanyaan terbuka]` Bagaimana renewal menggerbangi masuknya? Kandidat jawaban ada di
`Flow\InputRenewalFacultativeIn` (1.086 KB, berkas alur terbesar di RNW). Belum dibaca; bahan putaran
berikutnya.

---

## 1. Inventaris per tipe rule

```powershell
$nb="D:\migrasi\RNM\NB FacIn"; $rw="D:\migrasi\RNM\RNW Fac In"
foreach ($t in (Get-ChildItem $nb -Directory).Name) {
  "{0,-24} {1,6} {2,6}" -f $t,
    @(Get-ChildItem "$nb\$t" -Filter *.xml -File -EA SilentlyContinue).Count,
    @(Get-ChildItem "$rw\$t" -Filter *.xml -File -EA SilentlyContinue).Count }
```

| Tipe rule | NB | RNW | Selisih |
| --- | ---: | ---: | ---: |
| Activity | 609 | **556** | −53 |
| Section | 432 | **397** | −35 |
| FlowAction | 250 | **239** | −11 |
| RDBList | 217 | **200** | −17 |
| When | 210 | **189** | −21 |
| DataTransform | 148 | **138** | −10 |
| ReportDefinition | 123 | **118** | −5 |
| Harness | 41 | **41** | 0 |
| DataPage | 30 | **30** | 0 |
| DecisionTable | 12 | **10** | −2 |
| Flow | 6 | **4** | −2 |
| ConnectREST | 3 | **3** | 0 |
| DecisionTree | 1 | **1** | 0 |
| SystemSettings | 1 | **1** | 0 |
| **TOTAL** | **2.083** | **1.927** | **−156** |

⚠️ Selisih −156 pada tabel ini **bukan** ukuran delta yang sebenarnya: ia hasil bersih dari 176 berkas
yang hilang dikurangi 20 yang ditambah.

Berkas non-XML di RNW: **`Struktur_InputRenewalFacultativeIn.xlsx`** (7,2 MB) — layar input renewal,
belum dibaca (butir cakupan 7, putaran berikutnya).

---

## 2. Metode perbandingan — dan kontrol yang membuktikannya sahih

Klaim "nol perbedaan dari 1.907 berkas" terlalu kuat untuk diterima tanpa menguji instrumennya lebih
dulu. Empat kontrol dijalankan:

| Kontrol | Harapan | Hasil |
| --- | --- | --- |
| Dua rule yang memang berbeda (`ToUW` vs `ToJUW_A`) | terdeteksi **berbeda** | ✅ berbeda — instrumen peka |
| Berkas bersama NB↔RNW (`ToUW`) | sama | ✅ sama |
| Porsi baris yang dipertahankan normalisasi | tidak menggunting isi | ✅ **70–89 %** dipertahankan |
| Hash **mentah** NB↔RNW | diharapkan berbeda (metadata) | ⚠️ **identik** → melahirkan T-2 |

Normalisasi yang dipakai: buang baris bertag metadata `px*`/`pz*`, pangkas spasi, buang baris kosong,
urutkan dengan **`[StringComparer]::Ordinal`** (jebakan #3 `PANDUAN-KERJA` §2), lalu SHA-256.

```powershell
# hasil: 1907 / 1907 identik byte-per-byte
$byteSame=0;$byteDiff=0
foreach($t in @('Activity','ConnectREST','DataPage','DataTransform','DecisionTable','DecisionTree',
                'Flow','FlowAction','Harness','RDBList','ReportDefinition','Section','SystemSettings','When')){
  $A=@{}; Get-ChildItem "$nb\$t" -Filter *.xml -File -EA SilentlyContinue | ForEach-Object { $A[$_.BaseName]=$_.FullName }
  Get-ChildItem "$rw\$t" -Filter *.xml -File -EA SilentlyContinue | ForEach-Object {
    if($A.ContainsKey($_.BaseName)){
      if((Get-FileHash $A[$_.BaseName]).Hash -eq (Get-FileHash $_.FullName).Hash){$byteSame++} else {$byteDiff++} } } }
"$byteSame / $($byteSame+$byteDiff)"
```

---

## 3. Dua puluh berkas yang **hanya ada di RNW** — permukaan renewal yang sesungguhnya

Diurai dengan `[System.Xml.XmlDocument]`; langkah Activity dihitung sampai **kedalaman penuh**.

| Berkas | KB | Kelas | Langkah |
| --- | ---: | --- | ---: |
| `Flow\InputRenewalFacultativeIn` | 1.086,2 | `ASM-FW-GISFW-Work` | — |
| `Section\InputRenewalDtl_IsUW` | 2.401,0 | `ASM-FW-GISFW-Work` | — |
| `Section\InputRenewalDtl` | 2.250,7 | `ASM-FW-GISFW-Work` | — |
| `Section\PeriodeRenewal` | 428,9 | `ASM-FW-GISFW-Data-OfferFacIn` | — |
| `Section\SFAPortal_Renewal` | 298,6 | `Data-Portal` | — |
| `Section\InputRenewal` | 292,1 | `ASM-FW-GISFW-Work` | — |
| `Section\InputRenewal_IsUW` | 270,6 | `ASM-FW-GISFW-Work` | — |
| `Section\PeriodeRenewal_IsUW` | 243,2 | `ASM-FW-GISFW-Data-OfferFacIn` | — |
| `Harness\ChooseInsured` | 183,1 | `ASM-FW-GISFW-Data-Quotation` | — |
| `Section\ChooseInsuredDtl` | 132,7 | `ASM-FW-GISFW-Data-Quotation` | — |
| `Section\EditMarketing` | 101,9 | `ASM-FW-GISFW-Data-OfferFacIn` | — |
| `Activity\serviceInsertArasapasRNW_act` | 87,8 | `ASM-FW-GISFW-Work` | **10** |
| `ReportDefinition\BrowseAccountInsuredEDM` | 81,8 | `ASM-FW-SFAGISFW-Work-Account` | — |
| `ReportDefinition\RenewalList_RD` | 68,6 | **`ASM-FW-GISFW-Work-Renewal`** | — |
| `Activity\SetDataInsuredEDM_Act` | 68,1 | `ASM-FW-SFAGISFW-Work-Account` | **5** |
| `Activity\GetBusinessGroup_Act` | 45,7 | `ASM-FW-GISFW-Data-Quotation` | **4** |
| `FlowAction\Renewal_FlowAct_IsUW` | 33,4 | `ASM-FW-GISFW-Work` | — |
| `FlowAction\Renewal_FlowAct` | 33,2 | `ASM-FW-GISFW-Work` | — |
| `FlowAction\EditMarketing` | 29,8 | `ASM-FW-GISFW-Data-OfferFacIn` | — |
| `RDBList\CariBusinessGID` | 6,9 | `ASM-FW-GISFW-Int-BUSINESS` | — |

### Pengelompokan menurut fungsi `[dugaan]`

| Kelompok | Berkas |
| --- | --- |
| **Layar & alur masuk renewal** | flow masuk · 4 Section input renewal · 2 Section periode · 2 FlowAction renewal · portal |
| **Pemilih tertanggung** | `Harness\ChooseInsured` · `Section\ChooseInsuredDtl` · `ReportDefinition\BrowseAccountInsuredEDM` · `Activity\SetDataInsuredEDM_Act` |
| **Kelompok bisnis** | `Activity\GetBusinessGroup_Act` · `RDBList\CariBusinessGID` |
| **Sunting marketing** | `FlowAction\EditMarketing` · `Section\EditMarketing` |
| **Daftar renewal** | `ReportDefinition\RenewalList_RD` |
| **Konversi ke produksi** | `Activity\serviceInsertArasapasRNW_act` |

⚠️ **Nama bukan bukti, contoh baru.** `SetDataInsuredEDM_Act` dan `BrowseAccountInsuredEDM` bernama
**EDM** tetapi hanya ada di **RNW** — dan keduanya berkelas **`ASM-FW-SFAGISFW-Work-Account`**, kelas
SFA, bukan kelas Fac In maupun endorsement. Namanya menyesatkan di dua arah sekaligus.

---

## 4. Seratus tujuh puluh enam berkas yang **tidak dipakai renewal**

| Tipe | Jumlah | Tema nama yang menonjol |
| --- | ---: | --- |
| Activity | 56 | Treaty 15 · Life 2 · lain 39 |
| Section | 44 | Treaty 12 · perangkat/portal 4 · lain 28 |
| When | **21** | B2B/SFA/CRM 5 · Spv/approval 5 · Treaty 4 · Life 2 · perangkat 2 · lain 3 |
| RDBList | 18 | Treaty 7 · Life 1 · lain 10 |
| FlowAction | 14 | Treaty 5 · lain 9 |
| DataTransform | 10 | Treaty 7 · lain 3 |
| ReportDefinition | 7 | B2B/SFA/CRM 2 · Treaty 1 · Life 1 · lain 3 |
| Flow | 3 | Treaty 1 · lain 2 |
| DecisionTable | 2 | `isApproved` · `BusinessType_DeT` |
| Harness | 1 | perangkat |
| **TOTAL** | **176** | |

`[dugaan]` Tema yang berulang — **Treaty (≈53 berkas), Life, B2B/SFA/CRM, dan perangkat/portal** —
menyiratkan renewal tidak menangani jalur-jalur itu. Belum dibuktikan per berkas.

### Ke-21 rule `When` yang tidak ada di RNW — terkonfirmasi

```
crmCreateOpportunity · isApproved · isClaimTreaty · IsNotAdmin · IsOfferFacIn · IsOperatorLife
IsProposalTransfer · isSellingModeB2B · isSellingModeB2BB2C · isSellingModeB2C · IsSPVCreate
IsSPVTreaty1 · IsTreaty1 · IsUWAccepted · NopolisEmpty · pyIsIpadOrDesktop · pyIsMobile
StepStatusFail · ToDeptHeadUWLife · ToTREATYDEPTHEAD · TreatyMasterInEDM
```

⚠️ **Ejaan berbeda dari daftar rujukan** pada lima nama: `isClaimTreaty` · `NopolisEmpty` ·
`IsSPVCreate` · `ToTREATYDEPTHEAD` · `TreatyMasterInEDM`. Ini bukan kesalahan ketik — katalog
lintas-siklus sudah menandai **kepekaan huruf sebagai `[pertanyaan terbuka]`** yang belum terjawab.
Ejaan di atas diambil apa adanya dari nama berkas.

### Penelusuran lima rule yang paling berdampak

```powershell
foreach($n in @('IsUWAccepted','isApproved','IsOfferFacIn','NopolisEmpty','IsProposalTransfer')){
  "{0,-20} {1}" -f $n, @(Select-String -Path "$rw\*\*.xml" -Pattern "\b$n\b" -List).Count }
```

| Rule `When` NB-only | Masih dirujuk di RNW? | Catatan |
| --- | :-: | --- |
| `IsUWAccepted` | **4 berkas** | ⚠️ **Bukan rujukan menggantung.** DecisionTable `IsUWAccepted` **ADA** di RNW; yang tidak ada hanya rule `When` bernama sama. Dua tipe rule, satu nama — contoh "nama bukan bukti" |
| `isApproved` | **2 berkas** | ⛔ **RUJUKAN MENGGANTUNG** — DecisionTable-nya juga tidak ada di RNW. Lihat T-3 |
| `IsOfferFacIn` | **0** | Bersih. Gerbang masuk NB memang tidak dipakai renewal (T-5) |
| `NopolisEmpty` | **0** | Bersih |
| `IsProposalTransfer` | **0** | Bersih |

---

## 5. Pertanyaan terbuka dari putaran ini

| # | Pertanyaan | Kandidat sumber jawaban |
| ---: | --- | --- |
| R-1 | `isApproved` hilang dari RNW tetapi dirujuk 2 activity — disengaja, atau ekspor tidak lengkap? | Ekspor produksi tunggal; atau konfirmasi IT |
| R-2 | Bagaimana renewal menggerbangi masuknya, bila `IsOfferFacIn` tidak dipakai? | `Flow\InputRenewalFacultativeIn` |
| R-3 | Apakah `Work-Renewal` berarti antrean kasus tersendiri, atau sekadar kelas pelaporan? | `ReportDefinition\RenewalList_RD` + flow masuk |
| R-4 | `BusinessType_DeT` (DecisionTable penentu lini bisnis) tidak ada di RNW — bagaimana renewal menentukan lini bisnisnya? | penelusuran rujukan di korpus RNW |
| R-5 | Benarkah renewal **mengambil polis lama dan menghitung ulang**? Hipotesis `[dugaan]` yang belum diuji | `Section\InputRenewalDtl` · `Harness\ChooseInsured` · `RDBList\CariBusinessGID` |
| R-6 | Apakah 176 berkas NB-only benar-benar tak relevan renewal, atau sebagian menandakan ekspor tidak lengkap (sejalan R-1)? | pemeriksaan rujukan per berkas |

---

## 6. Yang belum dikerjakan — cakupan putaran berikutnya

- **Struktur langkah bersarang penuh** untuk 556 Activity RNW → *sebagian besar sudah tercakup
  discovery NB karena byte-identik; yang tersisa hanya 3 Activity RNW-only (19 langkah total)*
- `Flow\InputRenewalFacultativeIn` — alur masuk renewal, titik masuk dan transisi (R-2, R-3)
- Empat Section input renewal + dua Section periode — bagaimana polis lama diambil (R-5)
- `RDBList\CariBusinessGID` dan `Activity\GetBusinessGroup_Act` — jalur baca (R-4, R-5)
- `Struktur_InputRenewalFacultativeIn.xlsx` (7,2 MB) — layar input renewal
- Penelusuran per berkas atas 176 NB-only (R-6)

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
