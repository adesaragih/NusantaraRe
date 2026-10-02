# Arsitektur Target — Go + React + Oracle

> Sumber: sintesis dokumen discovery di `D:\migrasi\RNM\OUTPUT\`, yang seluruhnya dibangun dari
> korpus `D:\migrasi\RNM\{NB FacIn, RNW Fac In, Endorsment Fac In}\` (6.071 berkas `.xml`).
> Label mengikuti `CLAUDE.md` §3. Dokumen ini memuat **usulan rancangan**; setiap usulan menyebut
> temuan korpus yang memaksanya. Usulan yang tidak dipaksa korpus ditandai **[usulan]**.

---

## 0. Angka dasar korpus

```powershell
# jumlah berkas per folder korpus
foreach ($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")) {
  "{0,-20} {1}" -f $f, (Get-ChildItem "D:\migrasi\RNM\$f" -Recurse -Filter *.xml -File | Measure-Object).Count }
# -> NB FacIn 2083 / RNW Fac In 1927 / Endorsment Fac In 2061   (total 6071)
```

```powershell
# jumlah berkas per tipe rule
$rows = foreach ($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")) {
  Get-ChildItem "D:\migrasi\RNM\$f" -Directory | ForEach-Object {
    [PSCustomObject]@{ Tipe=$_.Name; Jml=(Get-ChildItem $_.FullName -Filter *.xml -File -Recurse | Measure-Object).Count } } }
$rows | Group-Object Tipe | ForEach-Object {
  [PSCustomObject]@{ Tipe=$_.Name; Total=($_.Group | Measure-Object Jml -Sum).Sum } } | Sort-Object Total -Descending
```

| Tipe rule | NB | RNW | EDM | Total |
| --- | ---: | ---: | ---: | ---: |
| Activity | 609 | 556 | 582 | 1.747 |
| Section | 432 | 397 | 434 | 1.263 |
| FlowAction | 250 | 239 | 265 | 754 |
| RDBList | 217 | 200 | 188 | 605 |
| When | 210 | 189 | 202 | 601 |
| DataTransform | 148 | 138 | 157 | 443 |
| ReportDefinition | 123 | 118 | 138 | 379 |
| Harness | 41 | 41 | 37 | 119 |
| DataPage | 30 | 30 | 41 | 101 |
| DecisionTable | 12 | 10 | 10 | 32 |
| Flow | 6 | 4 | 2 | 12 |
| ConnectREST | 3 | 3 | 3 | 9 |
| DecisionTree | 1 | 1 | 1 | 3 |
| SystemSettings | 1 | 1 | 1 | 3 |
| **Total** | **2.083** | **1.927** | **2.061** | **6.071** |

---

## 1. Keputusan arsitektur yang **dipaksa korpus**

Enam keputusan di bawah bukan preferensi. Masing-masing memiliki temuan yang membuat alternatifnya
salah.

### 1.1 Satu layanan, bukan tiga — `StatusBusiness` sebagai atribut case

`[terverifikasi]` Dari 2.439 identitas rule unik (Tipe + Nama) di seluruh korpus, **1.680 hadir di
ketiga folder**; dari 1.680 itu **1.220 logikanya identik** setelah metadata ekspor dibuang.

```powershell
# skrip lengkap: lihat §4 dokumen ini (normalisasi wajib: buang px*/pz*, urutkan ordinal)
# -> identitas di 3 folder: 1680 ; logika identik: 1220 ; bercabang: 460
```

`[terverifikasi]` **NB dan RNW tidak pernah berbeda satu sama lain** — 1.680 dari 1.680 identitas
bersama berlogika sama antara `NB FacIn\` dan `RNW Fac In\`. Seluruh 460 percabangan adalah
**EDM vs NB+RNW**. Di lapisan UI angkanya bahkan mutlak: 575 dari 575 identitas UI bersama identik
antara NB dan RNW (`03-ui/01-inventaris-layar.md` §5).

**Konsekuensi:** membangun tiga layanan menggandakan mayoritas kode tanpa dasar. Yang benar adalah
**satu layanan** dengan `StatusBusiness` sebagai atribut case, dan percabangan EDM sebagai cabang
eksplisit di dalamnya.

⚠️ `[pertanyaan terbuka]` Mengapa NB dan RNW identik sempurna belum dijelaskan korpus — apakah
Renewal memang memakai rule NB apa adanya, atau ekspor RNW adalah salinan ruleset NB. Jawabannya
menentukan apakah diskriminator siklus perlu dibawa ke lapisan UI sama sekali.

⛔ **Angka 460 adalah batas atas dan tidak boleh dipakai sebagai ukuran volume kerja.**
`[terverifikasi]` **86 rule bernama sama punya `pyRuleSetVersion` berbeda** antara NB dan EDM, dengan
arah **tidak konsisten** — sebagian salinan EDM lebih lama dan kehilangan cabang, sebagian justru
lebih baru dan menambah klausa. Ketiga folder diekspor dari titik waktu yang berbeda, sehingga
sebagian dari 460 itu adalah **selisih versi ekspor, bukan percabangan bisnis**. Korpus ini tidak
dapat memisahkan keduanya.

Keputusan §1.1 (satu layanan) **tidak terpengaruh** — ia ditopang 1.220 rule identik dan NB=RNW
1.680 dari 1.680. Yang terpengaruh adalah estimasi kerja. Lihat
`05-migrasi/03-risiko-dan-pertanyaan.md` R0.

### 1.2 Mesin tangga akseptasi = **mesin keadaan satu langkah**, bukan loop

`[terverifikasi]` (`01-flow/04-mesin-akseptasi.md` §0) Tangga persetujuan **bukan** loop yang
berjalan sampai selesai dalam satu proses. Satu putaran = satu keputusan manusia:

1. Petugas menekan tombol → DataTransform menulis `.ProposalAcceptStatus`.
2. Post-activity memanggil `GetLimitAkseptasi_Act` → SQL ke tabel limit → menulis **satu** properti:
   `pyWorkPage.LetterNo` (kode jabatan tujuan berikutnya).
3. DecisionTable `IsUWAccepted` memetakan `.ProposalAcceptStatus` → status konektor flow.
4. Gerbang "Limit Akseptasi" membaca `LetterNo` lewat rule `When` bernama `To*` → mengarahkan case
   ke antrean yang sesuai.
5. Bila `LetterNo` tidak cocok satu pun `To*`, cabang `Else` diambil — **itulah akhir tangga, bukan
   galat**.

**Konsekuensi:** di Go ini adalah *state machine* yang dipanggil per transisi, dengan keadaan
dipersistensi di antara panggilan. **Bukan** fungsi rekursif yang menghitung seluruh rantai
approver sekaligus. Merancangnya sebagai loop akan menghasilkan perilaku berbeda saat rantai
terputus di tengah.

### 1.3 Tiga field state, bukan satu "status"

`[terverifikasi]` (`01-flow/04-mesin-akseptasi.md` §6) Alur digerakkan tiga properti terpisah:

| Properti Pega | Isi | Nama kanonik usulan |
| --- | --- | --- |
| `PositionNote` | antrean yang **sedang** memegang case | `current_queue` |
| `LetterNo` | kode jabatan **tujuan berikutnya** | `next_approver_position` |
| `ProposalAcceptStatus` | hasil **keputusan terakhir** | `last_decision` |

⚠️ `LetterNo` **tidak menyimpan nomor surat.** Menyalin namanya mewariskan kebingungan yang sama;
`CONTEXT.md` §4 sudah menetapkan `next_approver_position` sebagai nama kanonik.

⚠️ `[terverifikasi]` **Dua ruang nama yang tidak boleh disatukan**: token antrean (`PositionNote`,
mis. `ReasFacInUnderwriting`) dan kode jabatan (`LetterNo`/`JABATAN`, mis. `SENIORUW`). Menyatukannya
merusak routing.

### 1.4 Uang tidak pernah `float`

`[terverifikasi]` (`04-aturan/02-formula-dan-status.md` §7) Bukti yang memaksa:

- **1.057 kemunculan di 91 berkas** konversi koma→titik (`@replaceAll(x,",",".")`). Komentar step
  pengembang lama menyatakannya eksplisit: *"karena property bukan decimal"*.
- **8 bentuk perbandingan ambang uang sebagai string berkutip**, mis. `.TSILiability > "3000000000"`.
  Ambang **yang sama** (30 miliar) dibandingkan numerik di rule lain — dua perilaku berbeda hidup
  berdampingan di korpus yang sama.
- Produk antara rumus premi mencapai orde 10²³ sebelum dibagi; `float64` punya 15–17 digit
  signifikan → **pasti** kehilangan presisi.
- Mata uang adalah dimensi eksplisit (`CurrencyList`, loop atas `CurrencyMaster.pxResults`).

**Konsekuensi mengikat** (rinci di `04-aturan/02-formula-dan-status.md` §7.2):

```go
// pkg/money
type Money struct {
    Amount   decimal.Decimal
    Currency string // WAJIB — korpus mengelola mata uang sebagai dimensi terpisah
}
```

`float32`/`float64` dilarang untuk uang, TSI, rate, persen, dan limit. Oracle: `NUMBER`. JSON API:
**string desimal**, bukan `number` JSON (yang IEEE-754 double di sebagian besar parser).

⚠️ **Presisi pembulatan adalah parameter per-rule, bukan konstanta global.** Aritmetika yang sama
(`TSI*Rate*Prorate*…/1e9`) dibulatkan **4 desimal** di `HitungPremi_FacInDT` dan **20 desimal** di
`CountPremi_ACT`. Menyeragamkannya mengubah angka.

### 1.5 Nama field JSON tidak boleh diubah

`[terverifikasi]` (`02-data/01-model-domain.md` §7) Agregat case diserialkan ke kolom
`JSON_POLIS.DATA_JSON` dengan nama properti identik. Mengubah nama field Go/JSON memutus pembacaan
data historis yang sudah tersimpan.

### 1.6 Versioning polis bersifat *append-only*

`[terverifikasi]` `PRODKE = COUNT(NOPOLIS) − 1`. Riwayat endorsement adalah **baris bertambah**,
bukan update in-place. Repository tidak boleh menyediakan operasi update terhadap baris polis lama.

⚠️ `[terverifikasi]` **Tiga mekanisme "nilai lama" yang berbeda** hidup berdampingan dan tidak boleh
disatukan menjadi satu konsep (`02-data/01-model-domain.md` §5).

---

## 2. Struktur target

Mengikuti `CLAUDE.md` §5. Arah dependency `handlers → services → repository` searah, tanpa memotong
lapisan (`CLAUDE.md` §4.2).

```
cmd/api/main.go
internal/
├── config/           endpoint & secret dari env — TIDAK PERNAH literal (CLAUDE.md §4.4)
├── handlers/         HTTP; tidak memuat aturan bisnis
├── models/           agregat case + Money; nama field JSON dikunci ke DATA_JSON
├── repository/
│   ├── oracle/       satu berkas per tabel; setiap query menyebut rule RDBList asalnya
│   └── lookup/       M_LINK_SERVICE dan master lain
├── rules/            registry predikat When — 226 nama unik
└── services/
    ├── faccase/      siklus hidup case + StatusBusiness
    ├── acceptance/   mesin tangga (§1.2) — inti aplikasi
    ├── premium/      rumus premi; presisi per-rule
    ├── spreading/    BELUM DIDISCOVERY — lihat §5
    ├── production/   konversi ke produksi, pasangan nilai-sesudah/delta
    └── integration/  panggilan servis luar
pkg/
├── money/            decimal + mata uang wajib
└── utils/
frontend/             React via Vite — src/{assets,components,hooks,pages,services,store}/
```

### 2.1 `internal/rules` — registry predikat

`[terverifikasi]` (`04-aturan/01-katalog-when.md` §1, §5) **226 nama rule `When` unik**
(case-insensitive) dari 601 berkas; **237** bila kapitalisasi dibedakan.

⚠️ **Temuan yang mengoreksi `CLAUDE.md` §4.5:** jumlah rule yang kondisinya **tidak terbaca sama
sekali adalah 0**. Kelima rule yang sebelumnya dinyatakan kosong ternyata terbaca penuh — kondisinya
tersembunyi di `<pyConditionValue1String>` yang tidak dibaca discovery sebelumnya:

| Rule | Kondisi terbaca |
| --- | --- |
| `IsPKSASM` | `OfferFacIn.IsB2B = "ASM"` |
| `ToUW` | `LetterNo = "UNDERWRITER"` |
| `ToJUW_A` | `LetterNo = "JUW_A"` |
| `LetterNoNull` | `LetterNo = ""` |
| `IsEdmInternalRetro` | `QuotationData.EndorsementInternalRetro = 1` |

**Daftar `panic()` di `CLAUDE.md` §4.5 kini kosong.** Mekanisme gagal-kerasnya tetap dipertahankan
untuk hal lain yang benar-benar tidak diketahui — lihat `02-data/03-batas-pengetahuan.md`.

⚠️ `[pertanyaan terbuka]` **11 nama rule dieja dengan kapitalisasi berbeda** antar siklus (semuanya
EDM vs NB+RNW). Bila resolusi rule di Pega case-sensitive, itu 22 predikat berbeda, bukan 11.

### 2.2 `frontend/` — pemetaan layar

`[terverifikasi]` (`03-ui/01-inventaris-layar.md` §1) 45 Harness unik, 520 Section unik, 307
FlowAction unik dari 2.136 berkas (35,2 % korpus).

⚠️ Pola yang **tidak boleh** diterjemahkan lurus ke React: **5 pasangan Section saling merujuk**,
tiga di antaranya membentuk siklus beranggota 3 (`CoverageItem` ↔ `InputCoverageFire` ↔
`PropertyItemListCoverage`) — rekursi tanpa batas bila dipetakan satu-ke-satu menjadi komponen.

---

## 3. Yang **hilang bersama Pega**

| Hal | Status | Konsekuensi |
| --- | --- | --- |
| Tabel internal Pega (`DATAPEGA.PC_*`) | `[terverifikasi]` | `GetHistoryAccPega_SQL` membacanya **untuk mengambil keputusan alur**, bukan sekadar audit. Sumber penggantinya belum ditentukan — **memblokir**. |
| Titik masuk portal (`Rule-Portal`) | `[terverifikasi]` tidak ada di korpus | Layar pertama pengguna harus ditetapkan bisnis, tidak dapat direkam. |
| Kontrak render server-side | `[terverifikasi]` 2.306 `pyPreDataTransform` + 1.810 `pyGridPreActivity` | Memindahkannya ke frontend **mengubah perilaku**; keputusan arsitektur yang harus disadari. |
| Isi `M_LINK_SERVICE` | `[terverifikasi]` tidak ada di korpus | Endpoint = konfigurasi. Nilainya diambil dari Oracle saat runtime, tidak pernah literal di kode. |

---

## 4. Catatan metodologis — cara membandingkan berkas korpus

Tiga jebakan ini menghasilkan angka yang **salah total** bila diabaikan. Semuanya ditemukan dan
dibuktikan selama discovery ini.

**Jebakan 1 — `Get-FileHash` mentah melaporkan 0 berkas identik.** Penyebabnya metadata ekspor,
bukan logika: operator impor, `pxCommitDateTime`, dan `pxHostId`. Nilai `pxHostId` yang berbeda
(`jboss122117` vs `jboss1073`) mengonfirmasi `CLAUDE.md` §3 butir 1 — korpus dirakit dari lebih dari
satu server Pega.

**Jebakan 2 — urutan tag ekspor Pega tidak deterministik.** Dibuktikan pada `When\ToUW.xml`, NB vs
EDM: **142 baris berbeda urutan, 0 berbeda isi**.

```powershell
$volatil = '^\s*<(px|pz)[A-Za-z0-9_]*[ >/]|^\s*</(px|pz)[A-Za-z0-9_]*>|^\s*<pyRuleFormStatusTime>|^\s*<pyShowJavaWindowName>|^\s*<rowdata REPEATINGINDEX="(RuleReference|Warnings)">'
$a = [IO.File]::ReadAllLines("D:\migrasi\RNM\NB FacIn\When\ToUW.xml")          | Where-Object { $_ -notmatch $volatil }
$b = [IO.File]::ReadAllLines("D:\migrasi\RNM\Endorsment Fac In\When\ToUW.xml") | Where-Object { $_ -notmatch $volatil }
(Compare-Object $a $b).Count               # -> 0   (tidak peka urutan)
(Compare-Object $a $b -SyncWindow 0).Count # -> 142 (peka urutan)
```

**Jebakan 3 — `Sort-Object` bawaan case-insensitive dan tidak stabil.** Korpus memuat nama rule yang
hanya berbeda kapitalisasi, jadi pengurutan **wajib** `[StringComparer]::Ordinal`. Dengan
`Sort-Object` bawaan, perbandingan yang sama melaporkan 507 identik; dengan pengurutan ordinal,
**1.220**.

Skrip lengkap yang menghasilkan angka §1.1 tersimpan sebagai lampiran di
`05-migrasi/03-risiko-dan-pertanyaan.md` §5.

---

## 5. Cakupan yang **belum** didiscovery

Jangan diperlakukan sebagai "tidak ada".

| Area | Status |
| --- | --- |
| **Spreading · capacity · scoring** | Sebagian tersentuh lewat rumus (`04-aturan/02` §2.4) dan proteksi spreading, tetapi **belum ada dokumen alurnya sendiri**. |
| Isi 460 rule bercabang EDM | Terhitung, **belum dipetakan field-per-field**. Sebelum dipetakan, tidak boleh diklaim "hanya beda kosmetik". |
| `Flow\` sebagai sumber urutan layar | 12 berkas; dipakai dokumen alur, belum dipakai untuk menurunkan urutan layar React. |
| `ReportDefinition\` (379 berkas) | Belum dibaca sama sekali. |

---

*Setiap query, predikat, dan transisi yang diimplementasikan wajib menyebut rule Pega asalnya dalam
komentar (`CLAUDE.md` §4.6). Itu yang membuat rekonsiliasi paralel run mungkin.*
