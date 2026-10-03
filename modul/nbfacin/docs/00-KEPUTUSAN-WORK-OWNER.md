# Keputusan Work Owner

> Catatan keputusan yang **menutup** pertanyaan terbuka. Setiap keputusan menyebut tanggal,
> pertanyaan asalnya, dan konsekuensinya terhadap dokumen lain.
> Keputusan di sini **mengalahkan** label `[pertanyaan terbuka]` di dokumen mana pun.

---

## K-001 · ~~Treaty Inward tidak termasuk lingkup proyek~~ — **DIBATALKAN oleh K-004**

> ⛔ **Keputusan ini dibalik pada hari yang sama.** Lihat **K-004** di bawah: Treaty Inward
> **termasuk** lingkup. Entri ini dipertahankan agar jejak keputusannya terbaca, **bukan** untuk
> dijalankan. Batas operasional di dalamnya sudah tidak berlaku.

**Tanggal:** 15 September 2026 · **Status: dibatalkan** · **Semula menutup:** blocker #1 di
`00-RINGKASAN-NB.md` §6

Pertanyaan asal: 80 dari 139 rule eksklusif NB adalah Treaty Inward, bukan Facultative Inward.

**Keputusan: tidak termasuk.**

### ⚠️ Keputusan ini **tidak dapat diterapkan secara mekanis**

Ini bukan keberatan atas keputusannya — keputusannya jelas dan saya jalankan. Ini peringatan tentang
**cara** menerapkannya, karena korpus tidak menyediakan batas yang bersih.

`[terverifikasi]` Dua kriteria yang tampaknya jelas ternyata sangat tidak sepakat:

```powershell
# untuk tiap berkas NB: cocokkan '(?i)treaty' pada nama berkas dan pada <pyClassName>
```

| Kriteria | Jumlah berkas NB |
| --- | ---: |
| **Nama berkas** mengandung `treaty` | 69 |
| **`pyClassName`** mengandung `treaty` | 58 |
| Keduanya | **21** |
| Gabungan | 106 |

Artinya **48 berkas bernama Treaty tetapi berkelas Fac In**, dan **37 berkas berkelas Treaty tetapi
bernama generik**.

Contoh yang membuktikan keduanya tidak dapat dipercaya sendirian `[terverifikasi]`:

| Berkas | Nama | Kelas | Kenyataan |
| --- | --- | --- | --- |
| `Section\FormulaTreatyCapacityDesc.xml` | "Treaty" | `Data-OfferFacIn` | dipanggil dari layar Fac In |
| `RDBList\InsertTreatyProduction_Sql.xml` | "Treaty" | `Int-policyjson` | bukan kelas Treaty |
| `Activity\GetKapasitasTreaty.xml` | "Treaty" | bukan Treaty | mesin spreading terbesar NB (1,3 MB) |
| `Activity\CalculatePremi_Act.xml` | generik | `Data-PolicyTreatyIn` | kelas Treaty, nama tidak |
| `Activity\CountSpreading_Act.xml` | generik | `Data-PolicyTreatyIn` | kelas Treaty, nama tidak |

Ini pengulangan `CLAUDE.md` §3 butir 3 pada tingkat lingkup proyek: **nama bukan bukti, dan di sini
kelas pun tidak cukup.**

### ⛔ 42 rule berkelas Treaty **dipanggil dari alur Fac In**

`[terverifikasi]` Dari 55 rule berkelas Treaty bernama ≥ 8 huruf, **42 dirujuk dari berkas berkelas
`Data-OfferFacIn` atau `ASM-FW-GISFW-Work`**.

Yang paling material:

> `Activity\serviceInsertArasapas_act.xml` berkelas **`ASM-FW-GISFW-Data-PolicyTreatyIn`** —
> kelas Treaty — tetapi ia adalah **panggilan servis konversi ke produksi** yang dirujuk 5 berkas
> Fac In. Discovery sebelumnya sudah mencatat bahwa versi kelas `ASM-FW-GISFW-Work` dari activity ini
> **tidak ada di korpus**; yang ada hanya versi kelas Treaty ini, yang mendelegasikan ke sana.

**Menghapus seluruh yang berkelas Treaty akan menghapus satu-satunya salinan terbaca dari jalur
konversi produksi Fac In.**

Rule berkelas Treaty lain yang dirujuk dari Fac In termasuk `FillPaymentInstallment` (6 rujukan),
`SpreadingRiskList` (4), `SetValidateInstallment_Act` (4), `GetCurrencyIDByName` (3), dan
`GenerateNoPolicy` (2).

### Batas operasional yang dipakai

`[usulan]` — sampai work owner menetapkan lain:

1. **Yang dikeluarkan dari lingkup:** rule berkelas Treaty yang **tidak** dirujuk dari kelas
   `Data-OfferFacIn` maupun `ASM-FW-GISFW-Work`. Dari 55 yang diperiksa, itu **13 rule**.
2. **Yang tetap di dalam lingkup meski berkelas Treaty:** 42 rule yang dirujuk dari alur Fac In —
   termasuk `serviceInsertArasapas_act`. Diperlakukan sebagai **dependensi Fac In**, bukan sebagai
   fitur Treaty Inward.
3. **Yang tetap di dalam lingkup meski bernama Treaty:** 48 berkas berkelas Fac In. Nama tidak
   mengeluarkan apa pun dari lingkup.
4. Setiap rule yang dikeluarkan **dicatat**, tidak dihapus diam-diam, supaya keputusan ini dapat
   ditinjau ulang.

⚠️ Angka 13 / 42 berasal dari pencocokan nama token. Sebelum dipakai memotong pekerjaan, daftar
13 rule itu **wajib ditelaah satu per satu** — ini batas ketelitian metode, bukan hasil final.

### Konsekuensi terhadap dokumen lain

| Dokumen | Perubahan |
| --- | --- |
| `00-RINGKASAN-NB.md` §6 | blocker #1 ditutup |
| `03-celah\03-rule-eksklusif-nb.md` | 80 rule Treaty Inward = di luar lingkup; 22 activity premi Treaty tidak diport |
| `04-spec\01-modul-go.md` | paket `production` kehilangan sebagian jalur sampai `serviceInsertArasapas_act` diputuskan statusnya |

⛔ **Tabel di atas TIDAK PERNAH DITERAPKAN.** K-001 dibatalkan sebelum konsekuensinya dijalankan.
Konsekuensi yang benar-benar berlaku ada di tabel K-004 dan K-005.

---

## K-002 · `IsOfferFacIn` memakai **`.Quotation.BusinessFac = "F"`**

**Tanggal:** 15 September 2026 · **Menutup:** blocker #7 di `00-RINGKASAN-NB.md` §6

Pertanyaan asal `[terverifikasi]`: rule `When/IsOfferFacIn` memuat **dua kondisi berbeda** — teks
tampilan `<pyConditionString>` berbunyi `Kode Bisnis = "02"/"58"/"SB"/"SG"`, sedangkan ekspresi
tersimpan `<pyConditionValue1String>` berbunyi `.Quotation.BusinessFac = "F"`. Rule ini menggerbangi
flow masuk NB.

**Keputusan: yang berlaku adalah ekspresi tersimpan, `.Quotation.BusinessFac = "F"`.**

```go
// Asal: NB FacIn/When/IsOfferFacIn.xml (pyConditionValue1String)
// Keputusan work owner K-002, 2026-09-15: ekspresi tersimpan yang berlaku.
// Teks tampilan rule menyebut kondisi lain (Kode Bisnis "02"/"58"/"SB"/"SG") —
// dicatat sebagai kandidat perbaikan, TIDAK diimplementasikan.
func IsOfferFacIn(c *Case) bool { return c.Quotation.BusinessFac == "F" }
```

⚠️ Teks tampilan yang berbeda **tetap dicatat** sebagai kandidat perbaikan. Ia tidak dihapus dari
dokumentasi, karena bila rekonsiliasi paralel run nanti memperlihatkan selisih pada gerbang masuk NB,
inilah tersangka pertamanya.

---

## K-003 · `IsSpreadingDepan` — tidak perlu `panic()`

**Tanggal:** 15 September 2026 · **Menutup:** usulan memasukkannya ke daftar gagal-keras

Pertanyaan asal: rule dirujuk `NB FacIn\Activity\SpreadingAdditionalProtection.xml` sebagai
pra-kondisi langkah, tetapi berkasnya tidak ada di `NB FacIn\When\`. Sempat diusulkan masuk daftar
`panic()` `CLAUDE.md` §4.5.

**Keputusan: tidak perlu.**

Dasarnya `[terverifikasi]`: rule itu **ada di folder Endorsment** dan kondisinya terbaca penuh —

```
pyLogic = A OR B OR C OR D
  A: Rule IsPA     evaluates to true
  B: Rule IsTravel evaluates to true
  C: Rule IsMBU    evaluates to true
  D: Rule IsFire   evaluates to true
```

Ini rule yang hilang dari **satu folder**, bukan rule yang kondisinya tidak diketahui. Jumlah rule
`When` yang kondisinya benar-benar tidak terbaca tetap **0**.

⚠️ **Satu hal yang perlu ditegaskan work owner.** Frasa "tidak perlu" dibaca di sini sebagai
**"tidak perlu `panic()`"** — predikatnya **tetap diimplementasikan** sesuai kondisi di atas, karena
menghapusnya akan mengubah perilaku pra-kondisi langkah di
`SpreadingAdditionalProtection.xml`. Bila yang dimaksud adalah "predikatnya tidak perlu ada sama
sekali", itu keputusan yang berbeda dan mengubah perilaku — mohon dikoreksi.

⚠️ Caveat yang dibawa ke implementasi: salinan yang terbaca berversi `01-01-52` (commit 2025-05-13)
dari folder EDM. Karena salinan NB tidak ada, **tidak dapat dipastikan** NB memakai versi yang sama —
95 rule di korpus ini terbukti berbeda versi antar folder.

---

## K-004 · Treaty Inward **termasuk lingkup proyek** — membatalkan K-001

**Tanggal:** 15 September 2026 · **Membatalkan:** K-001

**Keputusan: Treaty Inward masuk lingkup.**

### Yang langsung ikut selesai

Kekhawatiran terbesar pada K-001 **hilang dengan sendirinya**: tidak perlu lagi memisahkan 42 rule
berkelas Treaty yang dirujuk dari alur Fac In, dan `serviceInsertArasapas_act` — panggilan servis
konversi ke produksi — aman berada di dalam lingkup. Batas yang tidak bersih itu tidak perlu
ditarik sama sekali.

Yang ikut masuk lingkup:

| Cakupan | Jumlah |
| --- | ---: |
| Rule eksklusif NB bertema Treaty Inward | 80 dari 139 |
| — di antaranya activity perhitungan premi Treaty | 22 |
| Berkas NB bernuansa Treaty (nama atau kelas) | 106 |

### ⛔ Tetapi korpus di dalam batas kerja **tidak memuat Treaty Inward secara utuh**

Ini konsekuensi baru yang diciptakan keputusan ini, dan ia **memblokir**.

`[terverifikasi]` Treaty Inward punya **folder korpus tersendiri**, dan **tidak satu pun ada di dalam
`D:\migrasi\RNM\`** — seluruhnya hanya ada di `D:\migrasi\RNM_BRD\`, di luar batas kerja yang
ditetapkan work owner:

```powershell
# HANYA menghitung berkas; isi tidak dibaca (di luar batas kerja).
Get-ChildItem "D:\migrasi\RNM_BRD" -Directory | Where-Object { $_.Name -match '(?i)treaty' } |
  ForEach-Object { "{0,-28} {1}" -f $_.Name, (Get-ChildItem $_.FullName -Recurse -Filter *.xml -File | Measure-Object).Count }
```

| Folder | Berkas `.xml` | Di dalam batas kerja? |
| --- | ---: | :-: |
| `Treaty In Adjustment` | 379 | ❌ |
| `Treaty In` | 329 | ❌ |
| `Treaty Contract Out` | 303 | ❌ (Treaty **Out**, arah berlawanan — juga tidak dibaca per K-005) |
| `NB Treaty In` | 278 | ❌ |
| `EDM Treaty In` | 163 | ❌ |
| **Treaty Inward saja (tanpa Contract Out)** | **1.149** | ❌ |

Bandingkan: berkas bernuansa Treaty **di dalam** `NB FacIn\` hanya **106**.

**Artinya 106 berkas yang sudah didiscovery adalah pecahan, bukan gambaran Treaty Inward.** Rule
Treaty yang ada di folder Fac In hanyalah yang kebetulan ikut terekspor bersamanya.

### Pertanyaan yang timbul — **sudah dijawab K-005**

Ketiga pertanyaan tentang perluasan batas kerja ditutup oleh **K-005**: batas kerja **tidak**
diperluas.

### Konsekuensi terhadap dokumen lain

| Dokumen | Perubahan |
| --- | --- |
| Dokumen | Perubahan | Status |
| --- | --- | :-: |
| `00-RINGKASAN-NB.md` §2.6 | K-001 dicabut; cakupan Treaty terbatas dinyatakan | ✅ diterapkan |
| `00-RINGKASAN-NB.md` §6 | blocker "korpus Treaty di luar batas kerja" **tidak jadi ditambahkan** — K-005 menutupnya sebelum sempat berlaku | ✅ diterapkan |
| `03-celah\03-rule-eksklusif-nb.md` | 80 rule Treaty Inward **masuk** lingkup; §6.2 dan §7 butir 1–2 ditandai tercabut/terjawab | ✅ diterapkan |
| `04-spec\01-modul-go.md` §0 | peringatan K-001 dicabut; `production` tidak lagi terpotong | ✅ diterapkan |
| `03-celah\01-spreading-capacity-scoring-nb.md` | pertanyaan lingkup Treaty realization terjawab | ✅ diterapkan |

**Discovery NB yang sudah selesai tidak perlu diulang.** Dan karena K-005 menetapkan batas kerja
tidak diperluas, **tidak ada pass Treaty Inward tersendiri** dalam upaya ini.

---

## K-005 · Batas kerja **tetap** `D:\migrasi\RNM\NB FacIn\` — tidak diperluas

**Tanggal:** 15 September 2026 · **Menutup:** tiga pertanyaan perluasan batas kerja di K-004

**Keputusan: tidak perlu memperluas batas kerja. Fokus hanya `NB FacIn\`.**

Empat folder korpus Treaty Inward di `RNM_BRD\` (1.149 berkas) **tidak dibaca**, demikian pula
`Treaty Contract Out\` (303 berkas).

### Konsekuensi yang diterima

Treaty Inward **tetap di dalam lingkup proyek** (K-004), tetapi **cakupan discovery-nya terbatas**
pada apa yang muncul di dalam `NB FacIn\`:

| | Berkas |
| --- | ---: |
| Bernuansa Treaty **di dalam** `NB FacIn\` — didiscovery | **106** |
| Korpus Treaty Inward tersendiri — **tidak dibaca** | 1.149 |

`[terverifikasi]` Karena itu apa pun yang dikatakan dokumen discovery ini tentang Treaty Inward
berlaku **hanya untuk jejak Treaty di dalam folder Fac In**, bukan untuk modul Treaty Inward secara
utuh. Ini **batas cakupan yang disadari dan diterima**, bukan pertanyaan terbuka dan bukan blocker.

Konsekuensi praktisnya:

1. Rule Treaty yang dirujuk dari alur Fac In — termasuk `serviceInsertArasapas_act`,
   `FillPaymentInstallment`, `SpreadingRiskList` — **diport sebagai dependensi Fac In**, memakai apa
   yang terbaca di `NB FacIn\`.
2. Rule Treaty yang **hanya** dirujuk dari sesama rule Treaty diport apa adanya bila terbaca, dan
   ditandai *cakupan terbatas* bila rujukannya menunjuk ke luar folder.
3. Bila implementasi Treaty Inward yang utuh diperlukan nanti, itu **upaya terpisah** dengan batas
   kerja yang diperluas — bukan bagian dari pekerjaan ini.

⚠️ Satu hal yang perlu diingat saat rekonsiliasi: bila angka Treaty tidak cocok di paralel run,
**cakupan terbatas ini adalah tersangka pertama** — bukan cacat implementasi.

## K-006 · 11 activity yang hilang dari ekspor **sudah tidak digunakan** — tidak perlu diminta

**Tanggal:** 15 September 2026 · **Menutup:** butir 1 `_EKSTRAKSI-PEGA-SELAGI-HIDUP.md`, blocker
"badan 19 activity yang hilang" di `00-RINGKASAN-NB.md` §6

**Keputusan: kesebelas activity di bawah sudah tidak digunakan. Tidak diminta diekspor ulang secara
khusus.**

⚠️ **Dampaknya terbatas pada permintaan ekspor.** Keputusan ini **tidak** dengan sendirinya
mengeluarkan cabang pemanggilnya dari lingkup — lihat "Konsekuensi terhadap enam cabang pemanggil"
di bawah, yang **direvisi** setelah versi pertama catatan ini keliru menyimpulkan kode mati.

| Nama teknis | Dirujuk oleh (di korpus) |
| --- | --- |
| `GetLimitAkseptasi_Act2` | `InputOfferFacInEngineerUW_preACT`, `SetValidateDate_PostAct` |
| `GetLimitAkseptasi1SA_Act` | `GetLimitAkseptasi_Act`, `GetLimitAkseptasi_ActFlow` |
| `GetLimitAkseptasi2SA_Act` | `GetLimitAkseptasi_Act`, `GetLimitAkseptasi_ActFlow` |
| `CountTreatyCapacity_Act` | `SumTreatyCapacity_Act` |
| `CheckLimitAdditionalTreatyType_Act` | `SpreadingAdditionalProtection` |
| `SaveFacinProdEDMBonding_Act` | `SaveFacinProdAllEDM_Act` |
| `UpdateErrorNoteJsonPolis` | `serviceInsertArasapasEDM_act`, `UpdateErrorNoteJsonPolisMonitoring` |
| `CountASMGrossPremi_ACT` | rujukan internal |
| `CountASMNetPremi_ACT` | rujukan internal |
| `AgentSourceBizTreatyIn_Act` | rujukan internal |
| `BlankValuePropertyItem_Act` | rujukan internal |

### Konsekuensi terhadap enam cabang pemanggil — **DITANGGUHKAN, bukan dibuang**

> **Direvisi 15 September 2026 oleh work owner.** Versi pertama catatan ini menyimpulkan bahwa
> keenam cabang pemanggil adalah **kode mati** dan tidak perlu diport. **Kesimpulan itu dicabut.**
>
> Alasannya mengikat: **"hilang dari ekspor" ≠ "usang".** Ekspor ini diambil dari **titik waktu yang
> berbeda antar folder** (lihat R0 di `_ARSIP-lintas-siklus\_BACA-INI.md` — 95 rule berbeda versi,
> arah tidak konsisten). Ketiadaan sebuah rule di ekspor karena itu **bukan bukti** rule tersebut
> tidak dipakai. Dan Special Acceptance serta tangga akseptasi adalah **inti aplikasi**
> (`CLAUDE.md` §1) — membuangnya atas dasar inferensi adalah persis jenis "perbaikan diam-diam" yang
> dilarang aturan proyek.

`[terverifikasi]` Kesebelas activity **masih dirujuk** oleh activity yang ada di ekspor. Enam cabang
pemanggil berikut karena itu berstatus **DITANGGUHKAN**:

| Cabang pemanggil | Status |
| --- | --- |
| **Jalur Special Acceptance** di `GetLimitAkseptasi_Act` / `_ActFlow` | ⏸ DITANGGUHKAN |
| **Tangga akseptasi putaran kedua** di `SetValidateDate_PostAct` | ⏸ DITANGGUHKAN |
| **Limit tambahan treaty type** di `SpreadingAdditionalProtection` | ⏸ DITANGGUHKAN |
| **Kapasitas treaty** di `SumTreatyCapacity_Act` | ⏸ DITANGGUHKAN |
| **Simpan produksi endorsement bonding** di `SaveFacinProdAllEDM_Act` | ⏸ DITANGGUHKAN |
| **Penanganan galat konversi produksi** di `serviceInsertArasapasEDM_act` | ⏸ DITANGGUHKAN |

**Arti "ditangguhkan":** tidak diport **sekarang**, tidak dibuang **sama sekali**. Keenamnya
menunggu verifikasi terhadap **ekspor produksi tunggal** (butir 1 dan 2 permintaan di
`_EKSTRAKSI-PEGA-SELAGI-HIDUP.md`).

| Hasil verifikasi | Tindakan |
| --- | --- |
| Activity **ADA** di ekspor produksi | **Wajib diport.** Cabang pemanggil dipulihkan penuh. |
| Activity **benar-benar tidak ada** di ekspor produksi | Baru dikeluarkan dari lingkup — **dengan bukti**, dicatat sebagai amandemen K-006. |

⚠️ Sampai verifikasi itu selesai, **tidak ada cabang yang boleh dihapus dari rancangan**, dan tidak
ada klaim "kode mati" yang boleh dibuat atasnya.

---

## 🔄 Amandemen K-006 — 17 September 2026: **9 rule tiba, 6 cabang PULIH, 2 dikeluarkan**

Work owner menambahkan berkas rule ke `D:\migrasi\RNM\DDL\`. **Aturan K-006 dijalankan apa adanya:**
rule yang **ADA** → cabang pemanggilnya **dipulihkan**; yang **dinyatakan tidak dipakai** → dikeluarkan
dengan bukti.

`[terverifikasi]` Kesembilan berkas hadir, nama dan kelasnya cocok:

| Rule | Kelas |
| --- | --- |
| `GetLimitAkseptasi_Act2` | `ASM-FW-GISFW-Work` |
| `GetLimitAkseptasi1SA_Act` · `GetLimitAkseptasi2SA_Act` · `UpdateErrorNoteJsonPolis` | `ASM-FW-GISFW-Int-policyjson` |
| `CountTreatyCapacity_Act` | `ASM-FW-GISFW-Data-OfferFacIn` |
| `CheckLimitAdditionalTreatyType_Act` | `ASM-FW-GISFW-Data-Coverage` |
| `SaveFacinProdEDMBonding_Act` | `ASM-FW-GISFW-Work` |
| `AgentSourceBizTreatyIn_Act` | `ASM-FW-GISFW-Data-Agent` |
| `BlankValuePropertyItem_Act` | `ASM-FW-GISFW-Data-PropertyItem` |

`[terverifikasi]` Ketujuh activity pemanggil **ada** di `NB FacIn\Activity\`.

### Peta 6 cabang → status baru, **dengan bukti rujukan**

Setiap baris diverifikasi langsung ke korpus — **tidak ada yang dinyatakan pulih tanpa bukti**.

| # | Cabang | Pemanggil | Rule yang dituju | Bukti rujukan | Status |
| ---: | --- | --- | --- | --- | :-: |
| 1 | **Special Acceptance** | `GetLimitAkseptasi_Act` | `GetLimitAkseptasi1SA_Act` | L4008 `<RequestType>` | ✅ **PULIH** |
| 1 | | `GetLimitAkseptasi_Act` | `GetLimitAkseptasi2SA_Act` | L4265 `<RequestType>` | ✅ **PULIH** |
| 1 | | `GetLimitAkseptasi_ActFlow` | `GetLimitAkseptasi1SA_Act` | L4003 `<RequestType>` | ✅ **PULIH** |
| 1 | | `GetLimitAkseptasi_ActFlow` | `GetLimitAkseptasi2SA_Act` | L4260 `<RequestType>` | ✅ **PULIH** |
| 2 | **Tangga akseptasi putaran kedua** | `SetValidateDate_PostAct` | `GetLimitAkseptasi_Act2` | L2031 `<pyStepsActivityName>Call GetLimitAkseptasi_Act2` | ✅ **PULIH** |
| 3 | **Limit tambahan treaty type** | `SpreadingAdditionalProtection` | `CheckLimitAdditionalTreatyType_Act` | L1963 `<pyStepsActivityName>Call …` | ✅ **PULIH** |
| 4 | **Kapasitas treaty** | `SumTreatyCapacity_Act` | `CountTreatyCapacity_Act` | L2562 `<pyStepsActivityName>call …` | ✅ **PULIH** |
| 5 | **Simpan produksi endorsement bonding** | `SaveFacinProdAllEDM_Act` | `SaveFacinProdEDMBonding_Act` | L884 `<pyStepsActivityName>Call …` | ✅ **PULIH** |
| 6 | **Penanganan galat konversi produksi** | `serviceInsertArasapasEDM_act` | `UpdateErrorNoteJsonPolis` **lewat rantai** | lihat koreksi di bawah | ✅ **PULIH** |

**Keenam cabang pulih. Tidak ada yang masih tertahan.**

### ⚠️ Dua koreksi terhadap catatan K-006 lama

**1. Mekanisme pemanggilan tidak seragam — tiga di antaranya BUKAN `Call` activity.**

`GetLimitAkseptasi1SA_Act` dan `GetLimitAkseptasi2SA_Act` **tidak** dipanggil lewat
`<pyStepsActivityName>`. `[terverifikasi]` keduanya muncul sebagai **`<RequestType>` di dalam
`<pyStepsCallParams>` sebuah langkah `RDB-List`**, berdampingan dengan
`<ClassName>ASM-FW-GISFW-Int-policyjson</ClassName>`, `<Access>ASM</Access>`, dan
`<BrowsePage>LimitAkseptasi</BrowsePage>`.

Artinya keduanya **rule integrasi/SQL yang dieksekusi sebagai `RequestType`**, bukan Activity yang
di-`Call`. Menyebutnya "activity" — seperti catatan K-006 lama — **menyesatkan porting**: yang perlu
ditulis adalah query dan pemanggilan RDB-nya, bukan sebuah fungsi layanan.

**2. Cabang 6 adalah rantai DUA langkah, bukan rujukan langsung.**

Catatan lama menulis `serviceInsertArasapasEDM_act` sebagai perujuk langsung
`UpdateErrorNoteJsonPolis`. `[terverifikasi]` **tidak demikian** — yang dirujuknya adalah
**`UpdateErrorNoteJsonPolisMonitoring`** (L3273 `<RequestType>`), nama yang **berawalan sama**.
Rantai sebenarnya:

```
serviceInsertArasapasEDM_act
    └─> UpdateErrorNoteJsonPolisMonitoring      ← SUDAH ADA di NB FacIn\, tidak pernah hilang
            └─> UpdateErrorNoteJsonPolis        ← mata rantai yang hilang, kini tiba
                (UpdateErrorNoteJsonPolisMonitoring.xml L57:
                 RULE-CONNECT-SQL ASM-FW-GISFW-INT-POLICYJSON ASM!UPDATEERRORNOTEJSONPOLIS)
```

```powershell
# rujukan PERSIS (bukan diikuti "Monitoring") di seluruh korpus NB -> 1, dan itu di dalam
# UpdateErrorNoteJsonPolisMonitoring.xml sendiri
Select-String -Path "D:\migrasi\RNM\NB FacIn\*\*.xml" -Pattern 'UpdateErrorNoteJsonPolis(?!Monitoring)' -List
```

⚠️ **Pelajaran metodologis yang berulang:** pencocokan **awalan** menghasilkan kesimpulan yang salah.
Ini kekeliruan sejenis dengan hitungan `IsFacout` (11 vs 10, yang ke-11 adalah definisinya sendiri)
dan `JUW_A` (terlewat karena hanya satu kolom diperiksa). **Cocokkan batas kata, bukan awalan.**

### Dua rule DIKELUARKAN — dan rujukan menggantung yang ditinggalkannya

**Keputusan work owner: `CountASMGrossPremi_ACT` dan `CountASMNetPremi_ACT` tidak dipakai lagi.**
Keduanya **tidak** termasuk tujuh rule yang melayani enam cabang di atas.

⛔ `[terverifikasi]` **Tetapi keduanya masih dipanggil** oleh `NB FacIn\Activity\SetErrorMessage_Act.xml`
sebagai **langkah activity sungguhan**:

```
L518  <pyStepsActivityName>call CountASMGrossPremi_ACT</pyStepsActivityName>
L618  <pyStepsActivityName>call CountASMNetPremi_ACT</pyStepsActivityName>
L849  <pyRuleName>CountASMNetPremi_ACT</pyRuleName>        ← blok Embed-Reference-Rule
L885  <pyRuleName>CountASMGrossPremi_ACT</pyRuleName>
```

**Mengeluarkan keduanya meninggalkan dua rujukan menggantung di `SetErrorMessage_Act`** — pola yang
sama persis dengan `JUW_A` (K-024) dan `IsFacout` (K-019).

`[pertanyaan terbuka]` **Dua kemungkinan, dan memilihnya bukan hak migrasi:**

| Kemungkinan | Konsekuensi |
| --- | --- |
| Langkah pemanggil di `SetErrorMessage_Act` **ikut usang** — sisa kelewat saat penghapusan | Diport apa adanya (`CLAUDE.md` §1); jalurnya mati sendiri karena rule tujuannya tidak ada |
| `SetErrorMessage_Act` **sendiri perlu ditinjau** | Perlu keputusan tersendiri sebelum diport |

⛔ **Yang dilarang sampai ini diputus:** menghapus langkah pemanggil diam-diam, dan menyimpulkan
`SetErrorMessage_Act` sebagai kode mati tanpa bukti. **Tidak ada bukti** bahwa activity itu mati —
yang terbukti hanya bahwa dua langkah di dalamnya menunjuk rule yang dinyatakan usang.

### Sikap porting untuk rule bernuansa Treaty — **K-005 tetap berlaku**

`CountTreatyCapacity_Act`, `CheckLimitAdditionalTreatyType_Act`, dan `AgentSourceBizTreatyIn_Act`
tunduk pada **cakupan Treaty terbatas** (K-005): diport sebagai **dependensi Fac In** memakai apa yang
terbaca, dan **ditandai cakupan terbatas** bila rujukannya menunjuk ke luar folder `NB FacIn\`.
Korpus Treaty Inward tersendiri (1.149 berkas) tetap **tidak dibaca**.

### Yang tersisa dari kesebelas

Dua rule lain yang tiba — `AgentSourceBizTreatyIn_Act` dan `BlankValuePropertyItem_Act` — tercatat
di K-006 lama sebagai **"rujukan internal"**, bukan penggerbang salah satu dari enam cabang.
Kehadirannya melengkapi ekspor; tidak ada cabang yang bergantung padanya untuk pulih.

---

# Keputusan dari sesi grilling — 16 September 2026

Tujuh keputusan berikut lahir dari sesi grilling tiga ronde. Ringkasan sesi:
`03-keputusan/RINGKASAN-GRILLING.md`.

---

## K-007 · Rekonsiliasi paralel run **eksak**, dicapai bertahap

**Asal:** Ronde 1 Q1 + Ronde 2 Q3 + Ronde 2 Q4 · **ADR:** `adr/0001-rekonsiliasi-eksak-bertahap.md`

Paralel run tetap menjadi **ukuran keberhasilan** migrasi, dengan **nol selisih sampai digit
terakhir** — bukan toleransi. Paralel run **penuh** diakui baru mungkin setelah ekspor produksi
tunggal tiba, sehingga dijalankan bertahap:

1. **sekarang** — perhitungan murni atas masukan terekam;
2. **setelah mesin akseptasi ditulis** — per-modul, termasuk tangga dengan fixture;
3. **setelah `ALL_SOURCE` 35 prosedur tiba** — end-to-end.

**Konsekuensi:** toleransi ditolak karena hanya menyembunyikan salah-port; penyeragaman pembulatan
adalah perbaikan terpisah, bukan migrasi.

---

## K-008 · Satu flag fase: `panic` saat paralel run, `decline`+log saat produksi

**Asal:** Ronde 1 Q3 + Ronde 2 Q1 · **ADR:** `adr/0002-flag-fase-panic-vs-decline.md`

Mesin keadaan akseptasi dibangun **dari sisi penulis status** sekarang, diverifikasi terhadap ekspor
produksi nanti. Untuk jalur yang belum terverifikasi, perilaku bergantung fase — dikendalikan **satu
flag eksplisit, bukan dua basis kode**:

| Situasi | Paralel run | Produksi |
| --- | --- | --- |
| (a) status dikenali, baris keputusan belum terverifikasi | `panic` | `decline` + log |
| (b) kondisi tidak dikenali sama sekali | `panic` | `panic` |

**Syarat work owner:** transisi jalur (a) dari `panic` ke `decline`+log **hanya** setelah jalur itu
diverifikasi terhadap ekspor produksi.

⚠️ **Konflik yang diselesaikan:** usulan awal mengganti default `decline` dengan `panic` di semua
fase. Itu **ditarik** setelah terbukti `decline` adalah `<pyDefaultResult>` rule `IsUWAccepted` —
perilaku terekam, bukan celah pengetahuan. Menggantinya melanggar `CLAUDE.md` §1.

---

## K-009 · Mesin akseptasi ditulis sekarang; fixture menjadi kontrak data ke DBA

**Asal:** Ronde 2 Q2 · **ADR:** `adr/0003-mesin-akseptasi-fixture-sebagai-kontrak.md`
**Mencabut:** sikap sebelumnya yang menunda `services/acceptance` sampai `M_LIMIT_*` tiba

Yang hilang dari tangga akseptasi adalah **datanya**, bukan logikanya. `services/acceptance` ditulis
sekarang dengan fixture tabel limit, dan fixture itu **dikirim ke DBA sebagai kontrak bentuk data**
(butir 4 `_EKSTRAKSI-PEGA-SELAGI-HIDUP.md`).

**Konsekuensi:** ketidakcocokan bentuk tabel ditemukan saat permintaan dikirim, bukan di akhir.
~~Ejaan nilai kolom `JABATAN` tetap asumsi yang ditandai, bukan fakta.~~

### 🔄 Amandemen 17 September 2026 — data sudah tiba, asumsi bentuknya dikoreksi

**Keputusan intinya tetap berlaku dan terbukti benar.** Menulis fixture lebih dulu memang memunculkan
ketidakcocokan bentuk saat data datang — tepat seperti yang dirancang.

Yang dikoreksi: K-009 semula mengandaikan **satu bentuk tabel untuk tujuh lini**. `[terverifikasi]`
Kenyataannya **dua bentuk, enam tabel**:

| Bentuk | Tabel | Kolom | Ejaan `JABATAN` | Punya `WORKBASKET` / `JABATAN_ATASAN`? |
| --- | ---: | ---: | --- | :-: |
| **A — standar** | 5 | 14 | **tanpa spasi** (`DIREKTURTEKNIK`) | ya |
| **B — financial** | 1 | 7 | **pakai spasi** (`DIREKTUR TEKNIK`) | **tidak** |

Rincian, bukti kolom, dan konsekuensinya: **K-023** di bawah dan amandemen
`adr/0003-mesin-akseptasi-fixture-sebagai-kontrak.md`.

**Ejaan `JABATAN` tidak lagi asumsi** — untuk bentuk A ia kini `[terverifikasi]` dan **cocok** dengan
token routing di rule. Untuk bentuk B ia `[terverifikasi]` **tidak cocok**, dan itu temuan yang
mengikat implementasi.

`M_LIMIT_LIFE` dikonfirmasi **tidak ada** di basis data — bukan kekurangan. Target efektif enam tabel.

---

## K-010 · `Money` dan `Ratio` tipe terpisah; skala melekat pada nilai

**Asal:** Ronde 3 Q1 + Q2 · **ADR:** `adr/0004-money-dan-ratio-tipe-terpisah.md`

`Money{Amount, Currency}` dan `Ratio{Value, Scale}` **tidak dapat dijumlahkan**; satu-satunya
jembatan adalah `Money × Ratio → Money` yang eksplisit. Skala diisi **resolver COB terpusat** yang
mengikuti **rumus**, bukan label layar. COB tak dikenal → `panic`.

**Konsekuensi:** cacat satuan `.Rate` 10× berubah dari kesalahan runtime senyap menjadi kegagalan
kompilasi. Penyimpanan skala di tabel konfigurasi **ditolak** — ia fakta struktural korpus, dan tabel
yang dapat diubah tanpa jejak mengulang kerentanan `M_PROMPT_AI`.

---

## K-011 · Presisi pembulatan literal per-langkah, bukan registry

**Asal:** Ronde 3 Q3 · **ADR:** `adr/0005-presisi-pembulatan-literal-per-langkah.md`

Setiap pembulatan menuliskan presisinya di tempatnya, dengan komentar
`// Asal: <rule>, langkah <n>, presisi <p>`. Registry ditolak karena satu rule tidak punya satu
presisi — satu pembagi saja muncul dengan 9 presisi berbeda.

**Syarat work owner:** uji rekonsiliasi **wajib** memuat kasus pembulatan **di dalam loop
akumulasi**, bukan hanya nilai tunggal.

---

## K-012 · Mata uang `Unknown` sebagai keadaan eksplisit

**Asal:** Ronde 3 Q4 · **ADR:** `adr/0006-mata-uang-unknown-eksplisit.md`

Baca dan tampilkan boleh `Unknown`; **aritmetika lintas mata uang dan tulis ke Oracle `panic`**.
Default ke IDR ditolak (menebak); `panic` tanpa kecuali ditolak (menghentikan 112 layar yang di
sistem lama berjalan mulus).

**Konsekuensi:** implementasi menghasilkan **hitungan nilai produksi yang tiba tanpa mata uang** →
masuk kuesioner sebagai ukuran, bukan dugaan.

### ✅ Diperkuat pengukuran D1 — 17 September 2026

**Keputusan tidak berubah.** Ini konfirmasi silang, bukan amandemen.

`[terverifikasi]` `D:\migrasi\RNM\DDL\D1.xml`: `<TOTAL>691925</TOTAL>` ·
`<TANPA_MATA_UANG>4</TANPA_MATA_UANG>` → **4 dari 691.925 baris (0,0006 %)**.

Angka ini menjawab persis apa yang K-012 dirancang untuk mengukur, dan **menguatkan kedua sisi
keputusannya sekaligus**:

- **`Unknown` nyata, bukan teoretis** — ia benar-benar ada di data produksi. Default diam-diam ke IDR
  akan salah pada keempat baris itu, tanpa jejak.
- **`Unknown` sangat langka** — 0,0006 % berarti `panic` tanpa kecuali akan menghentikan sistem demi
  kasus yang praktis tidak pernah muncul, sementara memperlakukannya sebagai keadaan eksplisit
  praktis tidak berbiaya.

⚠️ Kelangkaannya justru menaikkan risikonya: cacat yang muncul 4 kali dari 691.925 **tidak akan
tertangkap pengujian sampel**. Ia hanya tertangkap oleh tipe yang memaksanya terlihat — yang memang
alasan keberadaan K-012.

---

## K-013 · Urutan kerja, dan penundaan keputusan tim

**Asal:** Ronde 1 Q2 + Q4

Urutan mengikuti §4 `_ARSIP-lintas-siklus/05-migrasi/03-risiko-dan-pertanyaan.md`: ekstraksi Pega →
`pkg/money` + `internal/rules` → mesin akseptasi (kini sejalan K-009) → jalur produksi terakhir.
Urutan final menyesuaikan prioritas bisnis.

**Siapa menulis kode Go/React: ditunda** — keputusan tim. Sementara itu pekerjaan dimulai dari modul
yang tidak menunggu jawaban bisnis.

---

# Jawaban kuesioner — 16 September 2026

Tiga kuesioner dikembalikan terisi. **Dua jawaban dicatat sebagai keputusan; tiga lainnya
bertabrakan dengan bukti korpus dan ditahan** sebagai konflik terbuka di bagian akhir.

---

## K-014 · `Decline` = penolakan final tanpa hak banding; `Reject` = dapat dibanding

**Asal:** kuesioner Underwriting butir 1 · **Pemilik jawaban:** Underwriting

Dugaan kami dikonfirmasi: **Decline** adalah penolakan final tanpa hak banding; **Reject** adalah
penolakan yang masih dapat dibanding.

**Konsekuensi:** keduanya **tidak** disatukan di sistem baru. Efek samping `Reject` yang terverifikasi
(mematikan konfirmasi binding dan penerimaan R/I slip, serta menjadi syarat munculnya jalur banding
lewat hitungan riwayat berstatus reject) adalah **perilaku yang dikehendaki**, bukan cacat — jadi
direproduksi apa adanya.

---

## K-015 · Renewal dinilai atas nilai pertanggungan **penuh** — disengaja

**Asal:** kuesioner Underwriting butir 2 · **Pemilik jawaban:** Underwriting

Perlakuan yang sudah terverifikasi dari korpus dikonfirmasi sebagai **kebijakan yang disengaja**:
renewal dibandingkan dengan tabel limit atas nilai pertanggungan penuh polis baru, bukan atas
kenaikannya terhadap polis yang berakhir.

**Konsekuensi yang menjangkau fase berikutnya:** asimetri antara renewal (nilai penuh) dan
endorsement (selisih) **juga disengaja**. Pertanyaan D-3 di
`04-kuesioner/_DITUNDA-fase-endorsement.md` — yang menanyakan asimetri ini dari sisi endorsement —
karena itu **sudah terjawab sebelum sempat dikirim**. Yang tersisa di fase endorsement hanyalah
pertanyaan nilai mutlak (penurunan diperlakukan seberat kenaikan), yang berdiri sendiri.

---

## K-016 · Satuan rate **MBU = persen (%)**

**Asal:** kuesioner Product/Aktuaria butir 1 · **Pemilik jawaban:** Product/Aktuaria

`[terverifikasi]` Jawaban ini **sejalan dengan rumus korpus**: seluruh rumus MBU
(`CountProrateExtension_Act` langkah 7, cabang `IsMBU` di `CountGPWMarinePAMbu_Act`,
`FillPremiMBU_FacIn`) berskala persen.

**Konsekuensi:** label `Standard Rate (‰)` pada layar input spreading MBU **salah**, dan diperbaiki
menjadi `(%)` di sistem baru. Ini perbaikan **label saja** — tidak ada angka yang berubah, sehingga
tidak mengganggu rekonsiliasi.

⚠️ Jawaban yang sama menyebut **PA dan Layering juga persen**. Itu **bertentangan dengan korpus** dan
**tidak** dicatat sebagai keputusan — lihat Konflik K-1 di bawah.

---

# ⚠️ Konflik terbuka — jawaban yang bertabrakan dengan bukti korpus

Ketiganya **tidak dicatat sebagai keputusan** sampai diselesaikan. `CLAUDE.md` §3 melarang menimpa
bukti terverifikasi dengan jawaban yang bertentangan tanpa menandai konfliknya.

---

## Konflik K-1 · ~~Satuan rate PA dan Layering~~ — **DISELESAIKAN oleh K-018**

> ✅ **Ditutup 16 September 2026.** Work owner mengunci satuan rate berdasarkan bukti korpus —
> lihat **K-018** di bawah. Jawaban kuesioner "PA = %" dinyatakan **keliru terhadap korpus**.
> Entri ini dipertahankan agar jejak konfliknya terbaca.

### Rumusan asli konflik (arsip)

**Asal:** kuesioner Product/Aktuaria butir 1 dan 2

Jawaban menyatakan **MBU, PA, dan Layering semuanya persen**, dan bahwa "rumusnya %, labelnya yang
salah". Untuk MBU itu benar (K-016). Untuk **PA dan Layering, korpus menunjukkan sebaliknya**:

| Lini | Ekspresi terverifikasi | Pembagi efektif untuk `.Rate` | Skala |
| --- | --- | ---: | :-: |
| **PA** | `@Math.divide((.TSI*.Rate),1000,4) * @Math.divide(.PctShortPeriod,100,4)` | 1.000 | **‰** |
| **PA** (varian) | `@Math.divide((.TSI*.Rate),100000,4) * ProRatePercent` | 100.000 ÷ 100 = 1.000 | **‰** |
| **PA** (varian) | `(@Math.divide((.TSI*.Rate),1000,4)*1) - .Discount` | 1.000 | **‰** |
| **Layering** | `@Math.divide((.TSI*.Rate*ProRatePercent),100000,4)` | 100.000 ÷ 100 = 1.000 | **‰** |
| MBU (pembanding) | `@Math.divide(.TSI*.Rate*.ProRatePercent,10000,4)` | 10.000 ÷ 100 = 100 | **%** |

Diperkuat lapisan layar: `Section\InputSpreadingPA_FacIn.xml` **mencetak tanda per mille U+2030
secara harfiah** — `‰) Policy Rate`.

**Mengapa ini memblokir:** bila sistem baru memakai % untuk PA dan Layering sementara sistem lama
memakai ‰, **premi kedua lini itu meleset 10×**. Ini bukan selisih pembulatan — ini satu ordo
besaran.

**Tiga kemungkinan, dan kami tidak dapat memilih sendiri:**

1. Jawaban terfokus pada MBU saja, dan PA/Layering tidak ikut ditinjau.
2. PA/Layering **seharusnya** persen — berarti sistem lama selama ini salah hitung, dan ini temuan
   jauh lebih besar dari migrasi.
3. Ada perbedaan antara satuan **di slip** (yang ditanyakan) dan satuan **yang disimpan** (yang
   dihitung rumus) — keduanya bisa berbeda dan tetap konsisten.

**Penyelesai yang sudah berjalan:** Pengukuran 2 kuesioner DBA meminta rentang nilai `RATE` per lini
bisnis. **Besaran angkanya membuktikan satuannya tanpa perlu pendapat** — rate PA tipikal ~2,5
berarti ‰; ~0,25 berarti %. Tahan keputusan sampai angka itu masuk.

---

## Konflik K-2 · ~~`IsFacout` "tidak dipakai lagi" vs 10 rule yang merujuknya~~ — **DISELESAIKAN oleh K-019**

> ✅ **Ditutup 16 September 2026.** Work owner memilih **kemungkinan ke-3**: fitur fac out **usang
> secara bisnis**, tetapi kodenya **masih aktif dan tetap diport apa adanya** — lihat **K-019** di
> bawah. Entri ini dipertahankan agar jejak konfliknya terbaca.

### Rumusan asli konflik (arsip)

**Asal:** kuesioner Underwriting butir 3

Jawaban: *"aturan `IsFacout` itu tidak dipakai lagi."*

`[terverifikasi]` Korpus menunjukkan rule ini **dirujuk 10 berkas di setiap folder**, termasuk lima
activity penyetel data fac out per lini bisnis (`SetDataFacOutFire_Act`, `SetDataFacOutAnekaGolf_Act`,
`SetDataFacOutCargoMBU_Act`, `SetDataFacOutPATravel_Act`, `SetDataFacOut_Act`), dua post-activity
validasi tanggal, satu activity penjumlah fac out PA, dan dua Section input coverage.

```powershell
Select-String -Path "D:\migrasi\RNM\NB FacIn\*\*.xml" -Pattern '\bIsFacout\b' -List
```

**Ini pola yang sama dengan K-006**, dan pelajarannya berlaku: menghapus predikat yang masih dirujuk
akan mematikan gerbang lima jalur penyetelan data fac out sekaligus.

**Yang perlu ditegaskan — tiga kemungkinan berbeda:**

1. **Predikatnya** tidak dipakai lagi → maka kesepuluh pemanggilnya adalah kode mati. Perlu
   dikonfirmasi eksplisit, karena konsekuensinya besar.
2. **Namanya** yang usang — predikatnya tetap berjalan, hanya namanya menyesatkan. Maka kami
   pertahankan perilakunya dan ganti namanya.
3. **Fitur fac out**-nya yang tidak dipakai lagi secara bisnis, tetapi kodenya masih aktif berjalan.
   Maka ia tetap diport (perilaku terekam), dan penghapusannya adalah perubahan terpisah.

~~Sampai ditegaskan, `IsFacout` dan kesepuluh pemanggilnya berstatus **⏸ ditangguhkan** — diport,
tidak dihapus.~~ → **status ditangguhkan dicabut oleh K-019: porting normal.**

---

## Konflik K-3 · ~~Alamat email "tidak usah dipindahkan" vs `CLAUDE.md` §4.4~~ — **DISELESAIKAN oleh K-020**

> ✅ **Ditutup 16 September 2026.** Work owner menerima jawaban Keamanan/IT dan mencatatnya sebagai
> **pengecualian eksplisit** terhadap §4.4 — persis bentuk yang diminta di akhir entri ini. Lihat
> **K-020** di bawah. Entri ini dipertahankan agar jejak konfliknya terbaca.

### Rumusan asli konflik (arsip)

**Asal:** kuesioner Keamanan/IT butir 1

Jawaban: *"tidak usah dipindahkan."*

`CLAUDE.md` §4.4 menyatakan: **endpoint dan host adalah konfigurasi/env var, tidak pernah literal di
kode.** Mempertahankan alamat email sebagai literal melanggar aturan yang mengikat proyek ini.

**Dua konsekuensi yang perlu disadari sebelum ini ditetapkan:**

1. Sedikitnya satu penerima adalah **alamat perorangan**. Notifikasi akan berhenti atau salah alamat
   ketika orang itu berpindah peran, dan tidak ada yang tahu sampai ada yang mengeluh.
2. Perubahan daftar penerima akan **menuntut deployment** dan tidak meninggalkan jejak — sama seperti
   kerentanan `M_PROMPT_AI` yang sudah kita catat sebagai isu keamanan.

**Yang kami minta:** bila keputusannya tetap "tidak dipindahkan", catat sebagai **pengecualian
eksplisit terhadap §4.4** dengan alasannya — bukan sebagai kelalaian. Pengecualian yang tercatat
dapat ditinjau ulang; pelanggaran yang tidak tercatat tidak.

---

---

## K-018 · Satuan rate per lini bisnis **DIKUNCI**

**Tanggal:** 16 September 2026 · **Menutup:** Konflik K-1 · **ADR:** `adr/0004-money-dan-ratio-tipe-terpisah.md`

| Lini bisnis | Satuan `.Rate` | Pembagi yang menempel pada `.Rate` |
| --- | :-: | ---: |
| **PA** | **‰** | 1.000 |
| **Layering** | **‰** | 1.000 |
| **FIRE** | **‰** | 1.000 |
| **MBU** | **%** | 100 |
| **ANEKA** · **BONDING** · **GOLF** · **MARINE CARGO** | **%** | 100 |

> ⚠️ **AMANDEMEN 1 Oktober 2026** `[keputusan work owner]` (`KEPUTUSAN-30-09-2026.md` butir 30):
> **Layering tidak dipakai lagi di sistem baru.** Baris Layering di atas tetap dikutip sebagai fakta
> korpus (‰, `GenerateLayerList_ACT.xml` L909), tetapi **tidak diport**; peta yang berlaku di kode menjadi
> **7 lini**.

### Aturan penguraian pembagi komposit — mengikat

Pembagi gabungan pada satu `@Math.divide` adalah **hasil kali faktor-faktor satuan**, bukan satu
satuan tunggal. Satuan disimpulkan dari **faktor yang menempel pada `.Rate` saja**, bukan dari
pembagi total:

| Faktor | Sumbangan ke pembagi |
| --- | ---: |
| `.Rate` ber-satuan ‰ | 1.000 |
| `.Rate` ber-satuan % | 100 |
| `ProRatePercent` (persen) | 100 |

Karena itu pembagi `100000` pada PA **menegaskan** `.Rate` = ‰ (1.000 × 100), bukan membantahnya.

### Bukti — dikutip persis, diverifikasi langsung

**PA = ‰** — `NB FacIn\Activity\CalculatePremiPA_FacIn.xml`:

```
L713:  (@Math.divide((.TSI * .Rate),100000,4)*pyWorkPage.OfferFacIn.ProRatePercent) - .Discount
L858:  (@Math.divide((.TSI*.Rate),1000,4)*@Math.divide((.PctShortPeriod),100,4))  - .Discount
L1003: (@Math.divide((.TSI*.Rate),1000,4)*1)-.Discount
L1146: @Math.divide((.TSI*pyWorkPage.OfferFacIn.ProRatePercent*.Rate),100000,20)- .Discount
```

Dua baris memberi pembagi **1.000 langsung**; dua lainnya memberi `100000` yang terurai menjadi
1.000 × 100. Keempatnya sepakat.

**MBU = %** — `NB FacIn\Activity\FillPremiMBU_FacIn.xml`:

```
L1144: @Math.divide((@Math.divide((.TSI*(.Rate+.Loading)),100,4)*@if(.ProRatePercent=="",100,.ProRatePercent)),100,4)
L1626: @Math.divide((@Math.divide((.TSI * .Rate),100,4)*@if(.ProRatePercent=="",100,.ProRatePercent)),100,4) + …
```

> ⚠️ **AMANDEMEN 1 Oktober 2026** `[keputusan work owner]` (`KEPUTUSAN-30-09-2026.md` butir 34, A4):
> **L1626 bukan rumus `.Premium`** — `[terverifikasi]` ia ekspresi di dalam syarat `When` langkah
> `.MinPremium`, bukan nilai yang ditulis ke `.Premium`. Rumus premi MBU yang diport hanya **L1144**.
> Kesimpulan satuan MBU = % tidak berubah: L1144 sendiri sudah membuktikannya.

⚠️ **MBU membuktikan aturan penguraian itu sendiri.** Bentuknya **bersarang**: pembagi **dalam**
(`100`) menempel pada `.Rate`, pembagi **luar** (`100`) menempel pada `ProRatePercent`. Jadi
keterikatan faktor→pembagi bukan tafsir kami, melainkan tertulis eksplisit dalam struktur ekspresinya.

**Layering = ‰** — `NB FacIn\Activity\GenerateLayerList_ACT.xml`:

```
L909: @Math.divide((.TSI * .Rate * pyWorkPage.OfferFacIn.ProRatePercent),100000,4)
```

`[terverifikasi]` Ini **satu-satunya** ekspresi ber-`.Rate` di seluruh 15 berkas bernuansa Layer di
korpus NB — tidak ada varian yang bertentangan. Bentuknya identik dengan PA L1146, dan terurai sama:
1.000 × 100.

```powershell
Get-ChildItem "D:\migrasi\RNM\NB FacIn" -Recurse -Filter *.xml -File |
  Where-Object { $_.Name -match '(?i)layer' } | ForEach-Object {
    $l=[IO.File]::ReadAllLines($_.FullName)
    for ($i=0;$i -lt $l.Count;$i++) { if ($l[$i] -match '\.Rate' -and $l[$i] -match 'divide') { "$($_.Name) L$($i+1)" } } }
# -> hanya GenerateLayerList_ACT.xml L909
```

**Penguat lapisan layar** `[terverifikasi]`: `NB FacIn\Section\InputSpreadingPA_FacIn.xml` mencetak
tanda per mille **U+2030 secara harfiah** pada label rate PA — `‰) Policy Rate`.

### ⚠️ Jawaban kuesioner "PA = %" dinyatakan keliru

Kuesioner Product/Aktuaria menjawab bahwa MBU, PA, dan Layering **ketiganya persen**. Untuk MBU itu
benar (K-016). Untuk **PA dan Layering jawaban itu keliru terhadap korpus** — dan bila diterapkan,
**premi kedua lini meleset 10×**. Bukan selisih pembulatan: satu ordo besaran.

Keputusan ini **tidak** menyalahkan penjawab. Kemungkinan besar pertanyaannya menanyakan satuan
**di slip**, sementara yang dikunci di sini adalah satuan **yang disimpan dan dihitung** — keduanya
bisa berbeda dan tetap konsisten bila ada konversi di antaranya.

### Status pengukuran D2 — ✅ **MASUK 17 September 2026, tidak bertentangan**

~~`[pertanyaan terbuka — menunggu D2]`~~ **Peta skala di atas TIDAK BERUBAH.** D2 berperan sebagai
konfirmasi silang, persis seperti yang ditetapkan — bukan sebagai pemetaan baru.

`[terverifikasi]` `D:\migrasi\RNM\DDL\D2.xml`: 30 nilai `RATE` terbanyak beserta jumlah kemunculannya.
Sebarannya **campuran dalam satu kolom** — banyak nilai jauh di bawah 0,1 (mis. `0,0244` dengan 72.022
kemunculan, `0,00000476581`) berdampingan dengan nilai di kisaran `0,45`–`1,43`.

**Kesimpulan yang boleh ditarik, dan yang tidak:**

| ✅ Boleh | ⛔ Tidak boleh |
| --- | --- |
| Satu kolom `RATE` memuat **besaran dari dua ordo yang berbeda jauh** — konsisten dengan dua skala yang hidup berdampingan, seperti yang K-018 tetapkan | Menyimpulkan **lini bisnis mana** memakai skala mana dari data ini |

⚠️ **D2 tidak dipecah per lini bisnis.** Karena itu ia **tidak dapat** mengonfirmasi maupun
membantah pemetaan per-COB — ia hanya memperlihatkan bahwa dua ordo besaran memang berdampingan.
Peta skala tetap bersandar pada **rumus di korpus**, yang sudah terverifikasi per baris di atas.

📌 D2 juga memunculkan temuan terpisah: nilainya **berkoma desimal** (`0,0244`), dan satu nilai
mempertahankan nol di belakang (`0,054900`). Ditangani **K-027** sebagai kontrak parsing.

### Konsekuensi

1. Resolver COB → skala di `internal/rules` memakai peta terkunci di atas.
2. **COB yang belum ada di peta → `panic`**, bukan skala default (ADR-0004).
3. Label `Standard Rate (‰)` pada layar MBU tetap diperbaiki menjadi `(%)` (K-016) — label saja.
4. Label `Rate (%)` pada tampilan PA dan `LayerListDtl` **salah** dan diperbaiki menjadi `(‰)`.
   Label saja; **nol angka berubah**, sehingga rekonsiliasi tidak terganggu.

---

## K-019 · `IsFacout` — fitur usang secara bisnis, **kode tetap diport apa adanya**

**Tanggal:** 16 September 2026 · **Menutup:** Konflik K-2 · **Mencabut:** status ⏸ ditangguhkan pada
`IsFacout` dan kesepuluh perujuknya

Dari tiga kemungkinan yang diajukan, work owner memilih yang **ketiga**:

> **Fitur fac out sudah tidak dipakai lagi sebagai kebijakan bisnis, tetapi kodenya masih aktif
> berjalan.** Karena itu ia **tetap diport** — perilaku terekam direproduksi apa adanya — dan
> penghapusannya adalah **perubahan terpisah**, bukan bagian migrasi (`CLAUDE.md` §1).

### Yang diport — 10 perujuk, dikonfirmasi ulang

`[terverifikasi]` Diukur ulang 16 September 2026:

```powershell
Select-String -Path "D:\migrasi\RNM\NB FacIn\*\*.xml" -Pattern '\bIsFacout\b' -List
# -> 11 baris
```

⚠️ **Perintah ini menghasilkan 11, bukan 10** — karena ikut menghitung berkas **definisinya sendiri**
(`When\IsFacout.xml`). Perujuknya tetap **10**. Catatan ini ada supaya audit ulang tidak menyimpulkan
angkanya keliru.

| Kelompok | Berkas | Jumlah |
| --- | --- | ---: |
| Activity penyetel data fac out per lini bisnis | `SetDataFacOutFire_Act` · `SetDataFacOutAnekaGolf_Act` · `SetDataFacOutCargoMBU_Act` · `SetDataFacOutPATravel_Act` · `SetDataFacOut_Act` | 5 |
| Post-activity validasi tanggal | `SetValidateDateUW_PostAct` · `SetValidateDate_PostAct` | 2 |
| Activity penjumlah fac out PA | `SumFacOutPA_Act` | 1 |
| Section input coverage | `InputCoverageFire` · `InputCoverageFire_IsUW` | 2 |
| **Perujuk seluruhnya** | | **10** |
| Berkas definisi rule | `When\IsFacout` | 1 |

### Konsekuensi

1. **Status ⏸ ditangguhkan dicabut.** Kesepuluh perujuk dan predikatnya berpindah ke **porting
   normal** — diperlakukan sama seperti rule aktif lainnya, tanpa penanda khusus.
2. **Tidak ada yang dihapus, dan tidak ada klaim "kode mati" yang boleh dibuat** atas kesepuluhnya.
   Ini pola yang sama dengan K-006, dan alasannya sama: menghapus predikat yang masih dirujuk akan
   mematikan gerbang lima jalur penyetelan data fac out sekaligus.
3. **Nama `IsFacout` dipertahankan** di kode port sebagai identitas rule Pega asal, demi ketertelusuran
   (`CLAUDE.md` §4.6) — **tetapi namanya bukan dokumentasi artinya.** `[terverifikasi]` isinya menguji
   `ProposalAcceptStatus = 4` (**Banding**), bukan fac out. Sudah tercatat di
   `steering/GLOSARIUM.md` §"Kosakata yang dihindari"; komentar di kode wajib menyebutkannya.
4. **Penghapusan fitur fac out dari sistem baru memerlukan keputusan tersendiri.** Keputusan ini
   **tidak** memberi mandat untuk menghapusnya nanti — ia hanya menetapkan bahwa migrasi tidak
   menghapusnya sekarang.

⚠️ Catatan untuk rekonsiliasi: karena fiturnya dinyatakan usang secara bisnis, **volume kasus yang
melewati jalur ini mungkin nol di data mutakhir**. Bila paralel run tidak pernah menyentuh jalur fac
out, itu **bukan bukti port-nya benar** — hanya bukti jalurnya tidak terpakai. Uji jalur ini dengan
fixture, bukan dengan mengandalkan data produksi.

---

## K-020 · Alamat email tetap literal — **pengecualian eksplisit terhadap `CLAUDE.md` §4.4**

**Tanggal:** 16 September 2026 · **Menutup:** Konflik K-3 · **Pemilik jawaban asal:** Keamanan/IT

**Keputusan: jawaban Keamanan/IT diterima.** Alamat email tujuan dan BCC pada jalur notifikasi polis
**tidak dipindahkan** ke tabel/env var pada tahap migrasi ini, dan tetap seperti keadaan sistem lama.

### Ini pengecualian **tercatat**, bukan kelalaian

`CLAUDE.md` §4.4 menyatakan endpoint dan host adalah konfigurasi/env var dan tidak pernah literal di
kode. Keputusan ini **menyimpang dari aturan itu secara sadar**, dan dicatat di sini persis supaya
penyimpangannya terlihat dan dapat ditinjau ulang. **Pengecualian yang tercatat dapat ditinjau ulang;
pelanggaran yang tidak tercatat tidak.**

⚠️ `CLAUDE.md` **tidak diubah** — §4.4 tetap berbunyi seperti semula. Batas kerja hanya pernah
mengizinkan penulisan pada §4.5. Pengecualian ini hidup **di register keputusan**, bukan di piagam
proyek.

### Dua risiko yang tetap berlaku dan diterima

1. `[terverifikasi]` Sedikitnya satu penerima adalah **alamat perorangan**. Notifikasi akan berhenti
   atau salah alamat ketika orang itu berpindah peran — dan tidak ada yang tahu sampai ada yang
   mengeluh.
2. Perubahan daftar penerima **menuntut deployment** dan tidak meninggalkan jejak audit — bentuk
   kerentanan yang sama dengan `M_PROMPT_AI` yang sudah dicatat sebagai isu keamanan.

Keduanya **bukan keberatan atas keputusannya** — keputusannya jelas dan dijalankan. Keduanya dicatat
supaya bila gejalanya muncul nanti, penyebabnya langsung dikenali alih-alih didiagnosis dari nol.

### Konsekuensi

1. Jalur notifikasi diport dengan daftar penerima sebagaimana terbaca di korpus.
2. **Nilai alamat emailnya tetap TIDAK disalin** ke dokumen, tiket, spec, maupun test mana pun —
   termasuk ke keputusan ini. Yang dicatat hanya **mekanisme dan jumlahnya**, sejalan dengan praktik
   yang berlaku sejak Prompt 1.
3. Peninjauan ulang dibuka kembali bila: daftar penerima berubah, alamat perorangan diganti alamat
   fungsional, atau audit keamanan memintanya. Cukup tambahkan `K-0nn` baru yang merujuk K-020.

---

## K-021 · Domain `ProposalAcceptStatus` dikunci ke **{1, 2, 3, 4, 7, 9}**

**Tanggal:** 16 September 2026 · **Asal:** kuesioner Underwriting butir 4 ·
**Pemilik jawaban:** Underwriting

Jawaban: *"Reject Ceding dan Ask Ceding sudah tidak dipakai."*

**Keputusan: kolom status keputusan underwriting hanya berisi enam nilai.**

| Nilai | Arti |
| ---: | --- |
| `1` | Accept |
| `2` | Reject |
| `3` | Ask |
| `4` | Banding |
| `7` | Decline |
| `9` | Revise |

Nilai **`5` (Reject Ceding)** dan **`6` (Ask Ceding)** `[terverifikasi]` hidup di **properti ceding
yang terpisah** dan tidak pernah menulis ke `ProposalAcceptStatus` — dikonfirmasi Underwriting bahwa
keduanya juga sudah tidak dipakai. Nilai **`8`** `[terverifikasi]` **nol jejak** di seluruh 6.071
berkas korpus.

### ⚠️ Validasi domain ini adalah **perilaku baru** — perlakukan sebagaimana mestinya

Ini bagian terpenting dari keputusan ini, dan mengabaikannya mengubahnya menjadi pelanggaran §1.

`[terverifikasi]` Sistem lama **tidak punya validasi domain pada kolom ini** — tidak punya validasi
apa pun selain wajib-isi. Sumbernya `04-spec\02-model-data.md` §1.2: nol aturan panjang, format,
rentang, atau pola di seluruh lapisan Section NB; tipe rule `Edit Validate`/`Validate` tidak ada di
ekspor. Menambahkan penolakan nilai karena itu **menambah perilaku yang sebelumnya tidak ada**.

Karena itu validasi ini **tidak boleh** berbentuk penolakan diam-diam:

| Fase | Nilai di luar {1,2,3,4,7,9} tiba | Alasan |
| --- | --- | --- |
| **Paralel run** | **`panic` + log** | Nilai tak terduga berarti dugaan kita salah, dan paralel run memang ada untuk menemukannya. Menolak diam-diam menyembunyikannya persis di fase yang seharusnya menangkapnya |
| **Produksi** | **gagal keras di jalur tulis + log** | Menulis nilai yang tidak dikenali ke kolom yang dipakai mengambil keputusan alur lebih berbahaya daripada berhenti |

**Dicatat sebagai kandidat perbaikan, bukan sebagai bagian migrasi.** Bila kelak paralel run
memperlihatkan nilai lain benar-benar ada di data produksi, itu **membatalkan premis keputusan ini** —
bukan membuktikan datanya rusak. Dalam hal itu, buka keputusan baru yang merujuk K-021, jangan
melonggarkan validasinya diam-diam.

### Konsekuensi

1. Tipe enum di `internal/models` memakai keenam nilai ini, dengan komentar menyebut rule Pega
   asalnya (`SaveViewSuggest.xml`, dekoder literal) sesuai `CLAUDE.md` §4.6.
2. Properti ceding yang menampung `5`/`6` tetap **dimodelkan terpisah** — tidak digabung ke dalam
   enum ini, karena korpus memang memisahkannya.
3. Keputusan ini **tidak** menutup pertanyaan baris keputusan `IsUWAccepted` (K-008): yang dikunci
   di sini adalah **kosakata kolomnya**, bukan pemetaan status → konektor flow.

---

> **Catatan penomoran.** **K-022 sudah dipesan** untuk keputusan riwayat akseptasi
> (`HISTORYAKSEPTASIPEGA` tetap dipakai · `DATAPEGA.pc_History_*` tidak dipakai lagi), sesuai
> `_CHECKLIST-KESIAPAN.md` §A.4 — **tetapi entrinya belum ditulis ke register ini.** Nomor itu
> **dicadangkan, bukan dilewati**; jangan dipakai untuk keputusan lain. Lihat juga K-017 yang memang
> kosong permanen.

---

## K-023 · Tabel limit ada **DUA BENTUK**, bukan satu — dan ejaannya berbeda

**Tanggal:** 17 September 2026 · **Mengamandemen:** K-009 dan `adr/0003-…` ·
**Sumber:** `D:\migrasi\RNM\DDL\M_LIMIT_*.txt` (DDL) + `M_LIMIT_*.xls` (isi)

**Keputusan: kedua bentuk diperlakukan sebagai dua bentuk yang berbeda. Tidak diseragamkan.**

### Bentuk A — standar, 14 kolom, 5 tabel `[terverifikasi]`

`M_LIMIT_PROPERTYY` · `M_LIMIT_ENGINEERINGG` · `M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL` ·
`M_LIMIT_PROPERTY_NON_PREFERREDD` · `M_LIMIT_NONPROPANDENGG`

```
ID · JABATAN · MAX_LIMIT_IDR · LIMIT_BOTTOM · MAX_LIMIT_USD · BATAS_WAKTU · NAMA
TGL_UPDATE · EFFECTIVE_DATE · TEAM_GROUP · LOGIN · WORKBASKET · JABATAN_ATASAN · LIMIT_BOTTOM2
```

Ejaan `JABATAN` **tanpa spasi**: `DIREKTURTEKNIK` · `DIREKTURMARKETING` · `KADIVTEKNIK` ·
`KADIVFACULTATIVE` · `DEPHEADUNDERWRITER` · `MANAGERTEKNIK` · `SENIORUW` · `UNDERWRITER` · `LEADER` ·
`JUW_B`. **Cocok dengan token routing di rule.**

### Bentuk B — financial, 7 kolom, 1 tabel `[terverifikasi]`

`M_LIMIT_FINANCIALINS`

```
ID · JABATAN · LIMITBOND_BOTTOM · LIMITCREDITCL_BOTTOM · LIMITCREDITNCL_BOTTOM · LIMITTRADE_BOTTOM · NAMA
```

Ejaan `JABATAN` **pakai spasi**: `DIREKTUR TEKNIK` · `DIREKTUR MARKETING` · `KADIV KEUANGAN` ·
`SENIOR UNDERWRITER`. **Tidak cocok dengan token routing.**

Bentuk ini **tidak punya** `WORKBASKET`, `JABATAN_ATASAN`, `TEAM_GROUP`, `EFFECTIVE_DATE`, `LOGIN`.
Kolom limitnya empat ambang per jenis pertanggungan, bukan `MAX_LIMIT_IDR`/`LIMIT_BOTTOM`.

```powershell
# menghasilkan 14 / 14 / 14 / 14 / 14 / 7
foreach ($n in @('M_LIMIT_PROPERTYY','M_LIMIT_ENGINEERINGG','M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL',
                 'M_LIMIT_PROPERTY_NON_PREFERREDD','M_LIMIT_NONPROPANDENGG','M_LIMIT_FINANCIALINS')) {
  $c=[IO.File]::ReadAllText("D:\migrasi\RNM\DDL\$n.txt")
  "{0,-42} {1,2}" -f $n, ([regex]::Matches($c,'"([A-Z0-9_]+)"\s+(?:VARCHAR2|NUMBER|DATE|CLOB|CHAR)')).Count }
```

### Konsekuensi

1. ⛔ **Mesin akseptasi financial mencocokkan jabatan berspasi.** Bila kedua bentuk disamakan, **tidak
   ada approver financial yang pernah ditemukan** dan tangga lini itu macet total. Ini persis
   kegagalan yang ADR-0003 peringatkan — sekarang terbukti nyata.
2. **Normalisasi ejaan DILARANG.** Menghapus spasi agar seragam adalah perubahan perilaku, bukan
   migrasi (`CLAUDE.md` §1).
3. ~~`[pertanyaan terbuka]` **Bagaimana eskalasi bentuk B bekerja?**~~ → ✅ **TERTUTUP oleh K-026**
   (17 September 2026): **tidak ada eskalasi berjenjang.** Penyaringnya **filter ambang limit** per
   jenis pertanggungan — `WHERE <kolom_limit> <= nilai` — dan `WHERE` itu **tetap diport**.
4. ~~`[pertanyaan terbuka]` **Antrean untuk lini financial?**~~ → ✅ **TERTUTUP oleh K-026**: memang
   **tidak ada antrean per jabatan**, karena kolomnya tidak ada dan tidak satu pun query membacanya.
5. ~~Bentuk A boleh berpindah dari fixture ke data nyata; bentuk B **tetap fixture** sampai butir 3
   dan 4 terjawab.~~ → **Butir 3 dan 4 sudah terjawab (K-026).** Kedua bentuk kini punya mekanisme
   yang terverifikasi penuh; keduanya boleh memakai data nyata, masing-masing **dengan mekanismenya
   sendiri** — tangga untuk bentuk A, filter ambang untuk bentuk B. Tidak disatukan.

### `M_LIMIT_LIFE` — bukan kekurangan

`[terverifikasi]` Dikonfirmasi **tidak ada** di basis data (keputusan work owner 16 September 2026).
Target efektif **enam** tabel limit, bukan tujuh. Jangan dicatat lagi sebagai item tertunggak.

---

## K-024 · ~~`JUW_A` — fitur usang, tetapi **barisnya masih ada di data**~~ — **DITUTUP**

> ## ✅ Ditutup 17 September 2026 — data sudah diperbaiki work owner
>
> **`JUW_A` dihapus dari `M_LIMIT_PROPERTYY` oleh work owner.** Fitur yang sudah usang secara bisnis
> kini **bersih juga dari data**.
>
> `[terverifikasi]` Diukur ulang dari `D:\migrasi\RNM\DDL\M_LIMIT_PROPERTYY.xls`
> (LastWrite 17 September 2026 13:36), 37 baris data:
>
> | | Hasil |
> | --- | --- |
> | `JUW_A` di seluruh isi tabel | **0 kemunculan** |
> | `JUW_B` | masih ada — 4 kemunculan di kolom `JABATAN` |
> | `JABATAN_ATASAN` | 33 dari 37 terisi; **`JUW_A` tidak lagi termasuk nilainya** |
> | Token jabatan lain | utuh, data tidak rusak |
>
> **Pertanyaan terbuka "disaring vs dibiarkan" GUGUR.** Tidak ada lagi baris yang menunjuk naik ke
> `JUW_A`, sehingga tidak ada yang perlu disaring maupun dibiarkan. Fixture dibangun dari data yang
> sudah bersih.
>
> ⚠️ **Angka rinci per-kolom di badan entri ini (di bawah) menggambarkan keadaan berkas SEBELUM
> perbaikan.** Berkas itu sudah ditimpa, sehingga angka tersebut **tidak dapat diukur ulang** dan
> **tidak boleh dikutip sebagai fakta yang berlaku**. Yang berlaku: *data lama memuat `JUW_A`; data
> yang diperbaiki tidak lagi memuatnya.* Badan entri dipertahankan sebagai jejak riwayat.

### ⚠️ Yang TIDAK ikut berubah — rule `ToJUW_A` tetap diport

**Ini konsekuensi DATA, bukan penonaktifan rule.** `ToJUW_A` **tidak** disebut rule mati, **tidak**
diubah, dan **tetap diimplementasikan sesuai kondisi terbacanya** (`CLAUDE.md` §4.5 dan §1).

`[terverifikasi]` Alur **masih menulis** nilai `JUW_A` ke `LetterNo`, diverifikasi langsung ke korpus
17 September 2026:

| Sumber | Isi |
| --- | --- |
| `NB FacIn\Activity\GetLimitAkseptasi_JUW_UW.xml` **L3105** | `pyWorkPage.LetterNo` ← `"JUW_A"` — penetapan **literal**, bukan hasil pembacaan tabel limit |
| `NB FacIn\Activity\GetLimitAkseptasi_JUW_UW.xml` L971 | `DataSearch.CARI5` ← `"JUW_A"` |
| `NB FacIn\DataTransform\SetLetterNo.xml` L194 | menulis nilai `"JUW_A"` |
| `NB FacIn\When\ToJUW_A.xml` | kondisi `pyWorkPage.LetterNo = "JUW_A"` — tidak diubah |

```powershell
Select-String -Path "D:\migrasi\RNM\NB FacIn\*\*.xml" -Pattern '"JUW_A"' -List
# -> GetLimitAkseptasi_JUW_UW.xml, SetBanding_ACT.xml, SetToInbox_ACT.xml,
#    DataTransform\SetLetterNo.xml, When\ToJUW_A.xml
```

**Artinya pembedaan berikut wajib dipahami sebelum fixture dibangun:**

| Mekanisme | Terpengaruh penghapusan baris? |
| --- | --- |
| Naik tangga **lewat tabel limit** (`JABATAN_ATASAN` → jabatan berikutnya) | **Ya** — `JUW_A` tidak lagi dapat menjadi tujuan naik |
| **Penetapan literal** `LetterNo = "JUW_A"` di `GetLimitAkseptasi_JUW_UW` L3105 | **Tidak** — jalur ini tidak membaca tabel limit untuk nilai itu |

Jadi `ToJUW_A` **tetap dapat tercapai**, dan menghapus predikatnya akan mengubah perilaku. Port
mereproduksi keduanya apa adanya.

`[pertanyaan terbuka]` Bila sebuah kasus tiba di `JUW_A` lewat penetapan literal itu, lalu tangga
mencari baris `JABATAN = 'JUW_A'` di tabel limit dan **tidak menemukannya** — apakah hasilnya
**tangga selesai** (perilaku "tidak ada yang cocok" per ADR-0003) memang yang dikehendaki? Belum
diverifikasi terhadap alur sesungguhnya. **Jangan ditebak**; angkat bila fixture menyentuh jalur ini.

### Rumusan asli (arsip — keadaan data SEBELUM perbaikan)

**Tanggal:** 17 September 2026 · **Pola serupa:** K-019 (`IsFacout`)

**Keputusan work owner: fitur `JUW_A` sudah TIDAK DIPAKAI LAGI** (usang secara bisnis).

### ⚠️ Tetapi jejaknya masih ada — dan bentuknya lebih spesifik dari dugaan awal

`[terverifikasi]` Diukur dari `D:\migrasi\RNM\DDL\M_LIMIT_PROPERTYY.xls`, 17 September 2026:

| Kolom | `JUW_A` | `JUW_B` |
| --- | ---: | ---: |
| `JABATAN` (baris yang mendefinisikan limit) | **0×** | 4× |
| `JABATAN_ATASAN` (tujuan naik berikutnya) | **2×** | 0× |

⛔ **Artinya `JUW_A` dirujuk sebagai tujuan, tetapi tidak punya baris yang mendefinisikan limitnya.**
Dua baris menunjuk naik ke `JUW_A`, dan tidak ada baris `JABATAN = 'JUW_A'` untuk dituju.

⚠️ **Koreksi atas laporan sebelumnya.** Laporan audit DDL saya menyatakan *"`JUW_A` tidak ada di satu
pun tabel limit"* — **itu keliru**. `JUW_A` **ada**, di kolom `JABATAN_ATASAN`. Pemindaian saya hanya
memeriksa kolom `JABATAN` dan `WORKBASKET`. Kekeliruan yang sama bentuknya dengan jebakan
"nama bukan bukti": saya menyimpulkan ketiadaan dari pencarian yang tidak menyeluruh.

### `[pertanyaan terbuka]` — WAJIB diputus SEBELUM fixture dibangun

> **Apakah baris ber-`JABATAN_ATASAN = 'JUW_A'` DISARING saat porting, atau dibiarkan apa adanya?**

Keduanya punya dasar, dan **pilihannya bukan milik migrasi**:

| Pilihan | Dasarnya | Akibatnya |
| --- | --- | --- |
| **Disaring** | Fitur sudah usang menurut work owner; sisa data adalah kelalaian penghapusan | Dua baris itu kehilangan tujuan naik — perlu ditetapkan tujuan penggantinya, atau menjadi ujung tangga |
| **Dibiarkan** | Reproduksi perilaku terekam (`CLAUDE.md` §1) | Kasus yang naik ke `JUW_A` tidak menemukan baris tujuan. `[terverifikasi]` perilaku "tidak ada yang cocok" di sistem lama = **tangga selesai**, jadi kasus itu berhenti di situ |

**Tidak ditebak di sini.** Ini pola yang sama dengan **K-019**: fitur dinyatakan usang secara bisnis
sementara jejaknya masih hidup di sistem. Pada K-019 keputusannya "diport apa adanya"; apakah keputusan
yang sama berlaku untuk data — bukan kode — adalah pertanyaan tersendiri.

⛔ **Fixture akseptasi tidak boleh dibangun sebelum butir ini diputus**, karena kedua pilihan
menghasilkan tangga yang berbeda untuk dua baris nyata.

---

## K-025 · Kolom `NAMA` dan `LOGIN` **wajib dibuang** dari fixture apa pun

**Tanggal:** 17 September 2026 · **Menegakkan:** `CLAUDE.md` §3 butir 5 dan batas kerja Prompt 1

`[terverifikasi]` Berkas isi tabel limit (`D:\migrasi\RNM\DDL\M_LIMIT_*.xls`) memuat **nama orang**
pada kolom **`NAMA`** dan **`LOGIN`** — keduanya ada di bentuk A maupun bentuk B.

**Keputusan: kedua kolom dibuang saat fixture, test, spec, atau artefak apa pun dibangun dari berkas
itu.** Bukan disamarkan, bukan diganti inisial — **dibuang**.

### Mengapa ini perlu keputusan tersendiri

Aturan "nama orang tidak pernah disalin" sudah berlaku sejak Prompt 1, tetapi selama ini berlaku pada
**dokumen**. Ini pertama kalinya nama orang masuk lewat **data yang akan menjadi kode** — dan fixture
yang di-commit ke repositori adalah tempat paling mudah untuk melanggarnya tanpa sadar.

### Konsekuensi

1. Fixture memuat kolom yang **dibaca query** saja. Untuk bentuk A: `JABATAN`, `JABATAN_ATASAN`,
   `LIMIT_BOTTOM`, `LIMIT_BOTTOM2`, `MAX_LIMIT_IDR`, `MAX_LIMIT_USD`, `BATAS_WAKTU`, `WORKBASKET`,
   `TEAM_GROUP`. **`NAMA` dan `LOGIN` tidak termasuk.**
2. ⚠️ Bila ternyata ada rule yang **membaca** `LOGIN` untuk mengambil keputusan, itu **guard berbasis
   identitas** — dicatat sebagai **jumlah dan mekanismenya saja** (`CLAUDE.md` §3 butir 5), tidak
   pernah nilainya. Belum diperiksa; ditandai `[pertanyaan terbuka]`.
3. Pemeriksaan kebocoran nama dijalankan atas fixture sebelum di-commit, memakai pencocokan
   **batas kata** — bukan substring, yang terbukti menghasilkan positif palsu.

---

## K-026 · Bentuk B (financial) — **filter ambang limit saja, tanpa tangga berjenjang**

**Tanggal:** 17 September 2026 · **Menutup:** dua `[pertanyaan terbuka]` di K-023 (eskalasi & antrean
bentuk B) · **Sumber:** `D:\migrasi\RNM\NB FacIn\RDBList\`

**Keputusan work owner: lini financial TIDAK bereskalasi berjenjang.** Satu-satunya penyaring adalah
**ambang limit per jenis pertanggungan**. Tidak ada rantai atasan, tidak ada antrean per jabatan —
karena kolomnya memang tidak ada di tabelnya (K-023 bentuk B).

> ⚠️ **AMANDEMEN 1 Oktober 2026** `[keputusan work owner]` (`KEPUTUSAN-30-09-2026.md` butir 46): **lini financial
> BERESKALASI** — tiga tingkat menurut antrean saat ini (UW Financial → Kadiv Keuangan → Direktur Marketing →
> Direktur Teknik), dirutekan Decision23 ke antrean `ReasFacInFinDivHead` / `ReasFacInMarketingDirector` /
> `ReasFacInTechnicalDirector`. Tangganya ditulis di ACTIVITY (`GetLimitAkseptasi_ActFlow` langkah 22,
> `GetLimitAkseptasi_Act` langkah 23), bukan di SQL — K-026 hanya membaca SQL. Yang tetap benar: ketiga SQL dan
> `WHERE`-nya, ejaan berspasi, dan kedua bentuk tidak disatukan. Juga keliru untuk bentuk A: tangga bentuk A
> **tidak** memakai `JABATAN_ATASAN`/`WORKBASKET` (butir 42, ralat C.1) — ia `ORDER BY LIMIT_BOTTOM`.

### Bukti — tiga rule SQL, dikutip dari tag `<pyBrowseSQL>`

`[terverifikasi]` Diverifikasi langsung 17 September 2026:

| Rule | Jenis pertanggungan | Isi `WHERE` |
| --- | --- | --- |
| `GetLimitAkseptasiBond_SQL.xml` | bond | `LIMITBOND_BOTTOM <= {DataSearch.CARID2}` |
| `GetLimitAkseptasiKreditCL_SQL.xml` | credit CL | `LIMITCREDITCL_BOTTOM <= {DataSearch.CARID2}` |
| `GetLimitAkseptasiKreditNCL_SQL.xml` | credit NCL | `LIMITCREDITNCL_BOTTOM <= {DataSearch.CARID2}` |

Bentuk lengkapnya seragam pada ketiganya:

```sql
SELECT JABATAN AS CARI1, <kolom_limit> CARI2, NAMA AS CARI5
FROM M_LIMIT_FINANCIALINS
WHERE <kolom_limit> <= {DataSearch.CARID2}
```

`[terverifikasi]` Ketiganya **tanpa `JOIN`**, **tanpa `JABATAN_ATASAN`**, **tanpa `WORKBASKET`**, dan
**tanpa `ORDER BY`**. Ketiadaan itu bukan kelalaian pembacaan — ia yang membuktikan tidak ada tangga.

```powershell
Select-String -Path "D:\migrasi\RNM\NB FacIn\RDBList\GetLimitAkseptasi*_SQL.xml" -Pattern '<pyBrowseSQL>' -Context 0,4
```

`[terverifikasi]` Pemanggilnya **sama dengan jalur bentuk A**: `Activity\GetLimitAkseptasi_Act.xml`
dan `Activity\GetLimitAkseptasi_ActFlow.xml`. Jadi satu activity menangani kedua bentuk, bercabang
menurut lini bisnis — bukan dua jalur terpisah.

### ⚠️ `WHERE` DIPERTAHANKAN — yang tidak ada adalah eskalasinya, bukan filternya

Ini pembedaan yang paling mudah salah baca, dan salah membacanya menghapus satu-satunya kontrol
wewenang yang dimiliki lini financial.

| Yang **TIDAK ADA** di bentuk B | Yang **TETAP ADA** dan wajib diport |
| --- | --- |
| Eskalasi berjenjang (`JABATAN` → `JABATAN_ATASAN`) | **Filter ambang `WHERE <kolom_limit> <= nilai`** |
| Antrean per jabatan (`WORKBASKET`) | Pemilihan kolom limit menurut jenis pertanggungan |
| Urutan eskalasi (`ORDER BY`) | Ketiga rule SQL apa adanya |

**Membuang `WHERE` dilarang** — ia perilaku terekam (`CLAUDE.md` §1). Yang dinyatakan tidak ada hanya
mekanisme naik-berjenjangnya.

### Konsekuensi terhadap `services/acceptance`

| | Bentuk A — standar (5 tabel) | Bentuk B — financial (1 tabel) |
| --- | --- | --- |
| Mekanisme | **Tangga**: `JABATAN` → `JABATAN_ATASAN`, satu keputusan = satu transisi | **Filter ambang**: satu query per jenis pertanggungan |
| Keluaran | jabatan tujuan **berikutnya** | **daftar jabatan** yang limitnya menampung nilai |
| Langkah naik | ada | **tidak ada** |
| Antrean | dari `WORKBASKET` | **tidak tersedia** |
| Ejaan `JABATAN` | tanpa spasi | **pakai spasi** (K-023) |

Keduanya **tidak disatukan di balik satu abstraksi**. Memaksa bentuk B ke dalam bentuk tangga berarti
mengarang langkah naik yang tidak ada di sistem lama.

### Catatan: kolom `NAMA` ikut di-`SELECT`

`[terverifikasi]` Ketiga query memilih `NAMA AS CARI5`. `[pertanyaan terbuka]` Apakah nilai itu
dipakai mengambil keputusan belum terbukti — `CARI5` adalah slot `DataSearch` serbaguna yang disebut
**63 Activity** di NB, jadi kehadirannya di sini **bukan bukti** ia menggerakkan alur.

⛔ Terlepas dari jawabannya, **K-025 tetap berlaku penuh**: nilai `NAMA` **tidak pernah** disalin ke
fixture, test, spec, atau artefak mana pun. Bila kelak terbukti ia menggerakkan keputusan, itu
**guard berbasis identitas** dan dicatat sebagai **jumlah serta mekanismenya saja**
(`CLAUDE.md` §3 butir 5).

---

## K-027 · `RATE` dan `PCTLIMIT` disimpan **pakai koma desimal** — kontrak parsing

**Tanggal:** 17 September 2026 · **Sumber:** `D:\migrasi\RNM\DDL\D2.xml`, `DDL\TABLEOFLIMIT.xml`

`[terverifikasi]` Nilai rasio di basis data produksi disimpan sebagai **teks dengan koma sebagai
pemisah desimal**, bukan titik:

| Berkas | Tag | Contoh nilai terbaca |
| --- | --- | --- |
| `DDL\D2.xml` | `<RATE>` | `0,0244` · `0,025` · `1,43` · `0,00000476581` · `0,054900` |
| `DDL\TABLEOFLIMIT.xml` | `<PCTLIMIT>` | `70,000 ` (dengan **spasi di belakang**) |

Diperkuat tipe kolomnya: `[terverifikasi]` `FACINOFFER.RATE` bertipe **`VARCHAR2(100)`** — teks, bukan
angka (lihat `_CHECKLIST-KESIAPAN.md` §A.3).

**Keputusan: parser Go WAJIB menangani koma sebagai pemisah desimal** di batas input, sebelum nilai
masuk ke `Ratio`. Konversi terjadi **di batas**, bukan tersebar di dalam perhitungan.

### Tiga jebakan yang wajib ditangani, bukan diasumsikan

1. ⛔ **`70,000` bersifat ambigu.** Dibaca sebagai koma-desimal ia bernilai **70**; dibaca sebagai
   pemisah ribuan ia bernilai **70.000** — selisih seribu kali. Parser yang memakai konvensi lokal
   otomatis akan menebak, dan tebakan itu tidak selalu sama. Konvensi **ditetapkan eksplisit**, tidak
   diserahkan ke locale.
2. **Spasi di belakang nilai** (`70,000 `) — dipangkas di batas input.
3. **Nol di belakang dipertahankan** (`0,054900`). Ini bukti tambahan bahwa nilainya disimpan sebagai
   teks: angka sejati tidak membawa nol tak berarti. Presisi tampilannya **tidak boleh** dipakai
   menyimpulkan presisi perhitungan.

### Ini kandidat perbaikan, bukan perbaikan

`CLAUDE.md` §1. Menyimpan rasio sebagai teks berkoma **tampak keliru** dan membuat perbandingan angka
menjadi perbandingan string — konsisten dengan temuan lama bahwa satu ambang dibandingkan sebagai
string (§4.1). **Dicatat, direproduksi apa adanya, tidak diperbaiki diam-diam.** Mengubah penyimpanan
menjadi numerik adalah perubahan terpisah milik bisnis.

⚠️ Catatan rekonsiliasi: bila paralel run memperlihatkan selisih pada nilai rasio, **parsing koma
adalah tersangka pertama** — bukan rumus preminya.

---

## K-028 · `HISTORYAKSEPTASIPEGA` tidak diarsip; kolom `ID_KOMITE` diabaikan ke depan

**Tanggal:** 17 September 2026 · **Keputusan work owner** · **Melengkapi:** `_CHECKLIST-KESIAPAN.md` §A.4

**Dua penegasan:**

1. **Tabel ini tidak ada di arsip** — tidak pernah diarsip maupun dipangkas. Pertanyaan lama
   *"apakah tabel ini pernah diarsip atau dipangkas?"* (butir 5 `_EKSTRAKSI-PEGA-SELAGI-HIDUP.md`)
   karena itu **tertutup**: riwayat yang ada adalah riwayat penuh, dan tidak ada periode yang hilang
   diam-diam.
2. **Kolom `ID_KOMITE` diabaikan ke depan.** Dahulu dipakai untuk **klaim**; pada sistem baru
   **tidak dipakai lagi**. Pertanyaan lama *"siapa yang mengisi `ID_KOMITE`?"* **gugur**.

`[terverifikasi]` Kolom itu memang ada di strukturnya — `DDL\HISTORYAKSEPTASIPEGA.txt` mendeklarasikan
`ID_KOMITE VARCHAR2(50)` di antara tujuh kolomnya.

### Konsekuensi

1. **Tabelnya tetap dipakai** — tidak berubah. `[terverifikasi]` ia dibaca untuk **mengambil
   keputusan alur**, bukan sekadar audit: jalur banding dan flag penolakan keduanya membacanya
   (`CLAUDE.md` §4.3).
2. `ID_KOMITE` **tetap ada di skema** (skema Oracle tidak berubah, §4.3) tetapi **tidak dibaca dan
   tidak ditulis** oleh sistem baru. Ia bukan kolom yang dihapus — ia kolom yang diabaikan.
3. Karena riwayatnya utuh, **hitungan riwayat berstatus reject** yang menentukan ketersediaan jalur
   banding dapat dipercaya apa adanya — tidak perlu penyesuaian untuk periode yang terpangkas.

---

## K-029 · Arti enumerasi **terverifikasi** — seluruhnya tertutup, tidak ada lagi yang menunggu Product

**Tanggal:** 17 September 2026 · **Menutup:** butir 6 Out of Scope spec · `_CHECKLIST-KESIAPAN.md` §C.1 ·
seluruh baris ☐ enumerasi di `steering/GLOSARIUM.md`

`[terverifikasi]` Keenam berkas property Pega di `D:\migrasi\RNM\DDL\` memuat pasangan
`<pyStandardValue>` (kode) dan `<pyLocalizedValue>` (arti) di dalam `<pyPromptTableList>`.
**Nilai dan artinya terbaca langsung dari ekspor — tidak satu pun ditebak.**

### `Type` — jenis penyesuaian · `DDL\Type.xml` (10 kode)

| Kode | Arti | Kode | Arti |
| ---: | --- | ---: | --- |
| `0` | Adjustment Reff. Number | `8` | Adjustment Deduction |
| `1` | Extend Period | `9` | Adjustment Insured Name |
| `2` | Adjustment TSI / Add Object / Rate / Premium | `11` | Adjustment Share Cedant |
| `4` | Adjustment Spreading | `12` | Adjustment PPN/PPH |
| `6` | Adjustment Period | | |
| `7` | Adjustment Currency | | |

⚠️ **Tanpa `3` dan `5`** — lihat bagian kode usang.

### `EdmType` — jenis endorsement · `DDL\EdmType.xml` (3 kode)

| Kode | Arti |
| ---: | --- |
| `1` | Batal Sejak Semula |
| `2` | Batal Prorata |
| `4` | Penambahan / Pengurangan / Perubahan |

⚠️ **Tanpa `3`** — lihat bagian kode usang.

### `TypeDeductible` · `DDL\TypeDeductible.xml` (8 kode)

`0` NIL · `1` % Of Claim · `2` % Of Loss · `3` % Of TSI · `4` **% Of Approved Loss Value** ·
`5` % Of Recoverable Claim Amount · `6` % Of Recoverable Amount · `7` In Amount

### `TypeDeductible2` · `DDL\TypeDeductible2.xml` (5 kode)

`0` NIL · `1` % Of Claim · `2` % Of Loss · `3` % Of TSI · `4` **% Of TSI Whichever Is Higher**

### ⛔ `TypeDeductible` **≠** `TypeDeductible2` — dua property berbeda

Keduanya berbagi kode `0`–`3` dengan arti yang sama, lalu **berpisah pada kode `4`**:

| Kode `4` di | Artinya |
| --- | --- |
| `TypeDeductible` | % Of **Approved Loss Value** |
| `TypeDeductible2` | % Of **TSI Whichever Is Higher** |

**Jangan disatukan menjadi satu enum.** Satu nilai yang sama berarti dua hal berlainan, dan
menggabungkannya menghasilkan salah hitung deductible tanpa gejala. Ini juga menjelaskan catatan lama
tentang "dua ruang nilai" pada deductible.

### `TeamGroup` · `DDL\TeamGroup.xml` (5 kode)

`1` Group 1 · `2` Group 2 · `3` Group 3 · `4` Group 4 · `5` **Bonding**

⚠️ Kode `5` **memutus pola** "Group N". Artinya tidak boleh diturunkan dari nomornya.

### `ProRateType` · `DDL\ProRateType.xml` (4 kode)

`1` Based on age and period (At Once) · `2` Based on period (At Once) ·
`3` Based on age (Yearly) · `4` Monthly

### `IsB2B` — **flag penanda**, ditutup · `DDL\B2B.xml`

**Keputusan work owner:** `IsB2B` hanya **flag penanda**. Nilai sahnya **`ASM` · `KBRU` · `BDX` ·
`SRB`** `[terverifikasi]`, dan **tidak ada arti bisnis yang lebih dalam** — nilainya sendiri adalah
isinya. Diperlakukan **apa adanya sebagai kode string**, bukan dipetakan ke makna lain.

⚠️ Sekali lagi: `IsB2B` **bukan boolean**, meski namanya berawalan `Is`.

**Tidak ada lagi enumerasi yang menunggu Product.**

---

### Kode usang yang sudah dihapus — `Type = 3`, `Type = 5`, `EdmType = 3`

**Keputusan work owner:** ketiganya **dahulu ada, kini sudah dihapus** karena usang. Jejak yang masih
tertulis di rule adalah **sisa yang kelewat saat penghapusan**.

Ini **pola ketiga kalinya** di proyek ini — sama dengan **K-019** (`IsFacout`) dan **K-024** (`JUW_A`):
fitur usang secara bisnis, jejaknya masih hidup di sistem.

#### ⚠️ Sikap portingnya BERBEDA untuk keduanya — jangan disamakan

| | `Type = 3` dan `Type = 5` | `EdmType = 3` |
| --- | --- | --- |
| Di mana jejaknya | Kondisi rule `When` | **Logika perhitungan premi yang aktif** |
| Dampak bila dihapus | **Nihil** | **Mengubah perilaku** |
| Sikap | Diport apa adanya | **WAJIB diport apa adanya** |

**`Type = 3` / `Type = 5` — aman.** `[terverifikasi]` (`_ARSIP-lintas-siklus\04-aturan\02-formula-dan-status.md`
L821–825): `IsEdmAdjRate` berkondisi `Type = 3 | Rule IsEdmAdjTSI evaluates to true` dan
`IsEdmAddObject` berkondisi `Type = 5 | Rule IsEdmAdjTSI evaluates to true`. Karena `Type = 3` dan
`Type = 5` **tidak pernah ditawarkan UI**, kedua rule itu **efektif setara `IsEdmAdjTSI`** — sejalan
dengan label `Type = 2` yang memang menggabungkan "TSI / Add Object / Rate / Premium".

⛔ **`EdmType = 3` — BERBEDA, dan ini yang paling mudah salah.** Cabangnya **masih tertulis di logika
perhitungan premi yang aktif**, diverifikasi **langsung ke korpus** 17 September 2026:

```
NB FacIn\Activity\SaveFacinProdEDMFire_Act.xml
  L10261  pyStepsPreCondParamsWhen:
          …QuotationData.EdmType==1 || …EdmType==2 || …EdmType==3     ← ketiga selisih persentase × −1
  L10810  …QuotationData.EdmType==3 && Datain1.CARI9==0
  L17883  …QuotationData.EdmType==3 && Datain1.CARI9==0
```

`[terverifikasi]` Pola yang sama muncul di **7 berkas folder NB** — bukan hanya folder endorsement:
`SaveFacinProdEDMFire_Act` · `SaveFacinProdEDMGolf_Act` · `SaveFacinProdEDMMarineCargo_Act` ·
`SaveFacinProdEDMMBU_Act` · `SaveFacinProdEDMPA_Act` · `SaveTreatyProductionEDMAneka_Act` ·
`InputPolicyTreatyInPre_Act`.

```powershell
Select-String -Path "D:\migrasi\RNM\NB FacIn\*\*.xml" -Pattern 'EdmType\s*==?\s*"?3"?' -List
# -> 7 berkas
```

**Karena itu:**

1. ⛔ **Cabang `||EdmType==3` DIPORT APA ADANYA.** **Jangan hapus `==3` dari kondisinya.**
   Menghapusnya mengubah perilaku terhadap **record lama** yang masih bernilai `3`, dan itu akan
   muncul sebagai selisih tak terjelaskan saat paralel run — pelanggaran `CLAUDE.md` §1.
2. **Status usangnya dicatat, bukan dieksekusi.** Karena tidak akan ada data `EdmType = 3` baru,
   cabang itu **menjadi jalur mati dengan sendirinya** seiring waktu — **tanpa satu baris kode pun
   diubah**. Persis mekanisme yang dipakai pada K-019.
3. ⚠️ Untuk rekonsiliasi: bila record lama ber-`EdmType = 3` ikut diuji, cabang `× −1` **harus**
   menghasilkan angka yang sama dengan sistem lama. Ini bukan jalur yang boleh dilewati dalam
   pengujian hanya karena statusnya usang.

---

### `IsPKSASM` **TETAP `panic`** — tidak disentuh

Daftar nilai `IsB2B` yang kini terverifikasi **tidak mencabut** `panic` pada `IsPKSASM`.

`[terverifikasi]` Penyebab `panic`-nya (`CLAUDE.md` §4.5) adalah **kondisinya belum ter-resolve** —
`<pyConditionString>` masih berisi placeholder `[Double click to add condition]` dan `<pyTempText>`
bernilai `true`, di ketiga folder. **Bukan** soal nilai `ASM` tidak dikenal.

Mengetahui nilai apa saja yang sah dan membuktikan kondisi itu dieksekusi adalah **dua pertanyaan
berbeda**; hanya yang kedua yang menyebabkan `panic`, dan yang kedua belum terjawab. Dicabut hanya
lewat ekspor produksi yang labelnya sudah ter-resolve.

---

## K-030 · Fokus aktif bergeser ke **`RNW Fac In`** — NB selesai sampai tiket

**Tanggal:** 17 September 2026 · **Mengamandemen:** K-005 · **Tidak membatalkan** pekerjaan NB

**Keputusan: siklus New Business dinyatakan SELESAI sampai tahap tiket.** Fokus aktif berpindah ke
siklus **Renewal**, dengan target yang sama: discovery → grilling → ADR/spec → tiket.

### Yang berubah dan yang tidak

| | Keadaan |
| --- | --- |
| Pekerjaan NB (discovery, spec, 16 tiket di `05-tickets\`) | **Final.** Tidak diulang, tidak dibuka kembali tanpa alasan baru |
| Batas **baca** | Kini **termasuk** `D:\migrasi\RNM\RNW Fac In\` |
| `Endorsment Fac In\` | **Tetap di luar fokus** — fase endorsement masih ditunda |
| Batas **tulis** | Tidak berubah: hanya `D:\migrasi\RNM\OUTPUT\` |
| Korpus Treaty Inward tersendiri di luar `RNM\` | Tetap **tidak dibaca** (K-005 bagian ini tetap berlaku) |

**K-005 tidak dibatalkan.** Ia mengunci fokus ke `NB FacIn\` saat pekerjaan NB berjalan, dan
pekerjaan itu selesai. Yang diamandemen hanya **fokus aktifnya**.

### Apa yang SUDAH diketahui tentang RNW — supaya discovery tidak mulai dari nol

Sumber: `_ARSIP-lintas-siklus\04-aturan\01-katalog-when.md`, diverifikasi ulang 17 September 2026.

**1. `[terverifikasi]` RNW tidak punya satu pun rule `When` eksklusif.**

| Folder | Rule `When` eksklusif |
| --- | ---: |
| `NB FacIn` | 18 |
| `Endorsment Fac In` | 16 |
| **`RNW Fac In`** | **0** |

Himpunan rule RNW adalah **himpunan bagian murni** dari gabungan NB+EDM. Siklus renewal **tidak
menambah predikat baru apa pun** — ia memakai ulang predikat NB.

**2. `[terverifikasi]` NB dan RNW selalu sepakat.** Dari **192** nama rule yang hadir di ≥2 folder,
**39 benar-benar bercabang** isinya antar siklus — dan **seluruh 39 melibatkan EDM** sebagai pihak
yang berbeda. NB dan RNW **tidak pernah** berbeda satu sama lain.

### ⚠️ Batas pengetahuan — apa yang TIDAK dibuktikan oleh kedua fakta itu

Ini bagian yang paling mudah disalahartikan, dan salah membacanya akan membuat discovery RNW
melewatkan hal yang penting.

**Katalog membuktikan rule-nya sama. Ia tidak menutup semua nuansa.**

- "`When` sudah terpetakan" berarti **daftar dan kondisi** rule `When` RNW sudah ada di katalog
  lintas-siklus, dan RNW memakai rule yang sama dengan NB. Itu **tidak otomatis** berarti setiap
  **perilaku runtime** RNW identik dengan NB.
- Bila discovery menemukan **cabang siklus yang membaca rule `When` secara berbeda** — misalnya
  activity atau flow yang menggerbangi predikat yang sama dengan cara lain — itu **temuan baru** yang
  wajib dicatat, bukan anomali yang boleh diabaikan karena "katalog bilang sama".
- `[terverifikasi]` **RNW memuat 189 berkas `When`, NB memuat 210** — selisih **21 rule yang ada di NB
  tetapi tidak ada di RNW**. "Himpunan bagian murni" karena itu **bukan** berarti setara: ada rule NB
  yang memang tidak dipakai renewal. Komposisi kedua puluh satu rule itu **belum ditelusuri** dan
  menjadi bahan discovery, bukan asumsi.

### Arah discovery RNW — **hipotesis kerja, bukan fakta**

`[dugaan]` Karena lapisan predikat sudah terpetakan, nilai terbesar discovery RNW kemungkinan ada di
**delta proses**: Activity, Flow, DataTransform, dan Section yang **khas renewal** — misalnya
pengambilan polis lama dan perhitungan ulang.

⚠️ Ini **arah pencarian, bukan kesimpulan**. Delta itu tetap harus **ditemukan dari korpus**, bukan
diasumsikan ada atau diasumsikan hanya itu. Bila ternyata deltanya lebih luas atau lebih sempit,
korpus yang menentukan.

### Yang sudah diputuskan dan tidak ditanyakan ulang

**K-015** — renewal dinilai atas **nilai pertanggungan penuh**, dan itu **disengaja**. Ini acuan
discovery RNW, bukan pertanyaan terbuka. Asimetri terhadap endorsement (yang memakai selisih) juga
sudah dinyatakan disengaja.

---

## K-031 · Temuan discovery Renewal — **RNW memakai ulang seluruh modul NB**

**Tanggal:** 17 September 2026 · **Sumber:** `06-rnw\01`…`04` · **Menguatkan:** K-015

### `[terverifikasi]` RNW adalah NB, minus 176 berkas, plus 20

| Ukuran | Nilai |
| --- | ---: |
| Berkas bersama NB↔RNW | **1.907** |
| — identik **byte-per-byte** | **1.907** (100 %) |
| — berbeda | **0** |
| Hanya di RNW | **20** |
| Hanya di NB | **176** |

`1.907 + 20 = 1.927` (RNW) · `1.907 + 176 = 2.083` (NB) — penjumlahannya menutup sempurna.

**Keputusan rancangan yang mengikat: siklus Renewal MEMAKAI ULANG modul NB, tidak diimplementasikan
tersendiri.** Pertanyaan "pakai ulang atau implementasi terpisah" terjawab untuk 1.907 berkas
sekaligus. Yang perlu dibangun khusus renewal hanyalah **20 berkas** — seluruhnya lapisan **masuk dan
tampilan**, bukan perhitungan.

```powershell
# menghasilkan 1907 / 1907 identik byte-per-byte
foreach($t in @('Activity','ConnectREST','DataPage','DataTransform','DecisionTable','DecisionTree',
                'Flow','FlowAction','Harness','RDBList','ReportDefinition','Section','SystemSettings','When')){
  $A=@{}; Get-ChildItem "D:\migrasi\RNM\NB FacIn\$t" -Filter *.xml -File -EA SilentlyContinue |
    ForEach-Object { $A[$_.BaseName]=$_.FullName }
  Get-ChildItem "D:\migrasi\RNM\RNW Fac In\$t" -Filter *.xml -File -EA SilentlyContinue | ForEach-Object {
    if($A.ContainsKey($_.BaseName)){
      (Get-FileHash $A[$_.BaseName]).Hash -eq (Get-FileHash $_.FullName).Hash } } }
```

### `[dilaporkan Claude]` Renewal **tidak punya mesin hitung ulang** — dan itu menjelaskan K-015

Hanya **empat** Activity RNW menggerbangi `StatusBusiness = 2`: masa berlaku tanggal, validasi
tanggal, proteksi spreading, input pembayaran. **Tidak satu pun perhitungan premi**, dan keempatnya
berkas bersama yang identik dengan NB.

Renewal menjalankan **perhitungan NB yang sama persis**. Kaitan ke polis lama adalah **rujukan**
(`OldPolicyNo`), **bukan pembawaan data**.

📌 **Ini menjelaskan K-015 dari sisi mekanisme.** Renewal dinilai atas nilai pertanggungan **penuh**
bukan karena ada kebijakan khusus yang memilih demikian, melainkan karena **tidak ada perhitungan
selisih sama sekali** — yang ada hanya perhitungan baru di atas kasus baru. Keputusan bisnis (K-015)
dan mekanismenya sejalan, dan K-015 **tidak berubah**.

### Konsekuensi

1. Tiket dan spec lima modul NB (`05-tickets\`, `04-spec\03`) **berlaku juga untuk renewal** tanpa
   perubahan pada lapisan perhitungan, registry predikat, maupun tangga akseptasi.
2. Yang perlu tiket tersendiri kelak hanyalah **20 berkas delta** — lapisan masuk dan tampilan.
3. ⚠️ `[dilaporkan Claude]` Gerbang masuk renewal **tidak ada di korpus** — lihat K-033 butir 2.

---

## K-032 · ~~⛔ **TIGA rujukan menggantung di RNW** — kandidat ekspor tidak lengkap~~ — **DITUTUP**

> ## ✅ Ditutup 18 September 2026 — ketiganya terjawab, **bukan ekspor tidak lengkap**
>
> Dugaan "ekspor RNW tidak lengkap" **tidak terbukti**. Ketiga rujukan punya sebab masing-masing, dan
> tidak satu pun memerlukan ekspor ulang.
>
> | # | Rule | Jawaban |
> | ---: | --- | --- |
> | 1 | `GenerateNoPolicy` | **Usang** — tidak dipakai lagi. Panggilan yang tersisa = **sisa kelewat** |
> | 2 | `SumTreatyCapacity_Act` | ✅ **Sudah ditambahkan ke `DDL\`** — tidak menggantung lagi |
> | 3 | `SearchJobID` | **Usang** — sudah lama dihapus. Panggilan tersisa = **sisa kelewat** |
>
> ⛔ **Konsekuensi: permintaan P-9 ke IT DIBATALKAN** — lihat `_PAKET-PERMINTAAN-DBA-IT-PRODUCT.md`.
>
> Badan entri di bawah dipertahankan sebagai jejak riwayat.
>
> ---
>
> ### 1 · `GenerateNoPolicy` — usang; jalur simpan renewal **sama dengan NB**
>
> **Keputusan work owner: sudah tidak dipakai lagi.** Simpan produksi NB **dan** Renewal memakai
> activity **yang sama**: `SaveJsonPolicyFacIn_Act`.
>
> `[terverifikasi]` 18 September 2026 — activity itu **identik byte-per-byte** antara NB dan RNW:
>
> ```powershell
> (Get-FileHash "D:\migrasi\RNM\NB FacIn\Activity\SaveJsonPolicyFacIn_Act.xml").Hash -eq
> (Get-FileHash "D:\migrasi\RNM\RNW Fac In\Activity\SaveJsonPolicyFacIn_Act.xml").Hash
> # -> True   (371.819 byte di kedua folder)
> ```
>
> Ia salah satu dari **1.907 berkas identik** (K-031) — memang activity simpan yang sama, bukan dua
> jalur berbeda.
>
> ⚠️ **Panggilan yang tersisa adalah sisa kelewat saat penghapusan**, `[terverifikasi]`:
>
> ```
> RNW\Activity\SaveJsonPolicyFacIn_Act.xml  L2204  <RequestType>GenerateNoPolicy</RequestType>
> RNW\Activity\GenerateNopolis_Act.xml      L3179  <RequestType>GenerateNoPolicy</RequestType>
> ```
>
> **Pola yang sama dengan K-019 (`IsFacout`) dan K-024 (`JUW_A`)**: usang secara bisnis, jejak
> panggilannya masih hidup di kode.
>
> ⛔ **Diport apa adanya** (`CLAUDE.md` §1). **Jangan hapus panggilannya diam-diam.** Status usangnya
> **dicatat, tidak dieksekusi** — tanpa rule tujuan, jalur itu mati dengan sendirinya tanpa satu baris
> pun diubah.
>
> ### 2 · `SumTreatyCapacity_Act` — ✅ sudah masuk, tidak menggantung lagi
>
> `[terverifikasi]` 18 September 2026: berkasnya **ada** di `D:\migrasi\RNM\DDL\SumTreatyCapacity_Act.xml`
>
> | | |
> | --- | --- |
> | Ukuran | 139,9 KB |
> | Kelas | `ASM-FW-GISFW-Data-OfferFacIn` |
> | Langkah (kedalaman penuh) | **17** |
>
> **Cabang kapasitas treaty (cabang 4 K-006) kini punya sumber terbaca** — bukan hanya di NB, tetapi
> sebagai berkas rule tersendiri. Rantainya lengkap:
> `CountASMShareTotal_ACT` / `SetValidateDate_Act` → `SumTreatyCapacity_Act` (DDL\) →
> `CountTreatyCapacity_Act` (DDL\, amandemen K-006).
>
> ### 3 · `SearchJobID` — usang, sudah lama dihapus
>
> **Keputusan work owner: sudah lama dihapus.** Panggilan yang tersisa `[terverifikasi]`:
>
> ```
> RNW\Activity\UploadCSVPerson_PostAct.xml  L2404  <RequestType>SearchJobID</RequestType>
> ```
>
> Sama seperti butir 1: **sisa kelewat, diport apa adanya, status usang dicatat.**
>
> ---
>
> ### 📌 Catatan untuk tahap tiket RNW — satu blocker yang ternyata tidak ada
>
> `[terverifikasi]` **Jalur simpan produksi renewal = jalur simpan produksi NB.** Keduanya memakai
> `SaveJsonPolicyFacIn_Act` yang identik byte-per-byte.
>
> **Renewal tidak memerlukan penghasil nomor polis tersendiri.** Blocker yang sempat dikira ada —
> dan yang melahirkan permintaan P-9 — **tidak pernah ada**. Ini sejalan dengan K-031: renewal
> memakai ulang modul NB, termasuk jalur simpannya.
>
> ### ⚠️ Yang tetap berlaku dari entri ini
>
> Pelajaran metodologisnya **tidak ikut gugur**: dari 40 nama yang diperiksa, muncul **empat mekanisme
> rujukan berbeda**, dan dua nama identik berperilaku berbeda menurut tipe rule-nya. **Pemeriksaan per
> nama lewat tag pembawa tetap tidak tergantikan** — lihat bagian "Empat yang TERNYATA BUKAN
> menggantung" di badan entri.

### Rumusan asli (arsip — sebelum dijawab work owner)

**Tanggal:** 17 September 2026 · ~~Status: `[pertanyaan terbuka]`~~ · **Menyentuh:** K-006 (butir 1) ·
jalur simpan produksi renewal (butir 2)

Tiga rule **dirujuk dari berkas yang ada di RNW**, tetapi **berkasnya tidak ada di RNW** — ketiganya
**ada di NB**.

| # | Rule | Tipe | Mekanisme rujukan | Perujuk di RNW | Bukti |
| ---: | --- | --- | --- | --- | --- |
| 1 | **`SumTreatyCapacity_Act`** | Activity | `<pyStepsActivityName>Call …` | `SetValidateDate_Act` L1129 · `CountASMShareTotal_ACT` L444 | `[terverifikasi]` |
| 2 | **`GenerateNoPolicy`** | RDBList | `<RequestType>` | `SaveJsonPolicyFacIn_Act` L2204 · `GenerateNopolis_Act` L3179 | `[terverifikasi]` |
| 3 | `SearchJobID` | RDBList | `<RequestType>` | `UploadCSVPerson_PostAct` L2404 | `[dilaporkan Claude]` |

### ⛔ Butir 2 paling berdampak

`GenerateNoPolicy` adalah **penghasil nomor polis**, dan dirujuk dari **jalur simpan polis**
(`SaveJsonPolicyFacIn_Act`). Bila benar hilang di renewal, **jalur simpan produksi renewal kehilangan
penghasil nomor polisnya**.

### ⚠️ Kandidatnya adalah "ekspor RNW tidak lengkap" — BUKAN "usang"

`CLAUDE.md` §4.5 berlaku penuh: **"hilang dari ekspor" ≠ "usang"**. Dan di sini bobotnya **lebih
besar** daripada kasus-kasus sebelumnya (K-006, K-019, K-024), karena satu fakta tambahan:

`[terverifikasi]` **1.907 dari 1.907 berkas bersama identik byte-per-byte, metadata ekspor termasuk.**
Kedua folder terbukti diekspor dari **keadaan ruleset yang sama**. Dalam keadaan itu, sebuah berkas
yang hadir di satu sisi saja menjadi **lebih** mencurigakan sebagai kelalaian ekspor, bukan kurang.

⛔ **Yang dilarang sampai ini diputus:** memvonis ketiganya usang, menghapus cabang pemanggilnya, atau
menyimpulkan renewal memang tidak memakai jalur itu. **Tidak ada bukti** untuk satu pun kesimpulan itu.

### Yang diminta

📌 **Permintaan ke IT: ekspor rule RNW yang lengkap** — analog permintaan B.1 untuk ekspor produksi
tunggal. Cukup untuk memastikan ketiga rule di atas memang tidak ada di ruleset renewal, atau sekadar
tidak ikut terekspor.

### Empat yang **TERNYATA BUKAN** menggantung — jangan diulang pemeriksaannya

`[terverifikasi]` Keempatnya tampak menggantung dari namanya, dan tidak satu pun benar:

| Nama | Kenyataan |
| --- | --- |
| `isApproved` | **properti** (`Rule-Obj-Property`) pada dua kelas, bukan DecisionTable |
| `IsUWAccepted` | **DecisionTable** bernama sama **ada** di RNW; hanya rule `When`-nya NB-only |
| `GetLimitAkseptasi_ActFlow` | hanya di `<pzOriginalInstanceKey>` — **metadata *save-as***, bukan panggilan |
| `GetInsuredID` | **`Activity\GetInsuredID` ADA di RNW**; yang NB-only adalah `RDBList\GetInsuredID` |

📌 **Pelajaran metodologis, dan ini yang paling berharga dari pemeriksaan ini.** Dari 40 nama yang
diperiksa satu per satu, muncul **empat mekanisme rujukan berbeda** — `pyStepsActivityName`,
`RequestType`, `Rule-Obj-Property`, `pzOriginalInstanceKey` — dan **dua nama yang identik
(`GetInsuredID`) berperilaku berbeda menurut tipe rule-nya**. **Tidak ada satu aturan umum** yang
memisahkan rujukan sungguhan dari tabrakan nama. Pemeriksaan **per nama, lewat tag pembawa**, tidak
dapat digantikan.

### Celah K-004 **tetap terbuka**

`[terverifikasi]` Discovery RNW **tidak** mengisinya:

| Activity | NB | RNW | Kelas | Langkah |
| --- | :-: | :-: | --- | ---: |
| `serviceInsertArasapas_act` | ada | ada | `Data-PolicyTreatyIn` | 1 (stub) |
| `serviceInsertArasapasEDM_act` | ada | ada | `ASM-FW-GISFW-Work` | 20 |
| `serviceInsertArasapasRNW_act` | — | ada | `ASM-FW-GISFW-Work` | 10 |
| **`serviceInsertArasapas_act` kelas `Work`** | **tidak ada** | **tidak ada** | sasaran delegasi stub | — |

Yang bertambah hanyalah **dua saudara terbaca** yang memperlihatkan polanya:
`GetLinkService` (dari `M_LINK_SERVICE`) → `Connect-REST` → catat log layanan. `[terverifikasi]` Pola
itu menguatkan `CLAUDE.md` §4.4 langsung dari korpus — tetapi **bukan rule yang dicari**.

---

## K-033 · Pertanyaan terbuka Renewal yang menunggu keputusan work owner

**Tanggal:** 17 September 2026 · **Status: seluruhnya `[pertanyaan terbuka]` — tidak satu pun
diputuskan sendiri**

| # | Pertanyaan | Menyentuh |
| ---: | --- | --- |
| ~~**RNW-1**~~ | ~~Ketiga rujukan menggantung (K-032)~~ → ✅ **SELESAI 18 September 2026.** `GenerateNoPolicy` **usang** · `SumTreatyCapacity_Act` **sudah masuk `DDL\`** · `SearchJobID` **usang**. **Bukan ekspor tidak lengkap**; permintaan ke IT dibatalkan. Dua yang usang **diport apa adanya** | ~~K-006~~ ✅ |
| ~~**RNW-2**~~ | ~~Gerbang masuk renewal tidak terekam~~ → ✅ **SELESAI 18 September 2026 — K-034.** Alur masuk **ditetapkan**: No. Polis + Renewal Date + Note → OK → activity pembuat renewal (di dalamnya salin data polis) | ~~rancangan~~ ✅ |
| ~~**RNW-3**~~ | ~~Bolehkah pola konversi dipakai sebagai acuan merancang?~~ → ✅ **SELESAI — K-035.** Pola konversi renewal **sama seperti NB**. Celah K-004 tetap terbuka | ~~K-004~~ ✅ |
| ~~**RNW-4**~~ | ~~Tinjauan keamanan `Section\EditMarketing`~~ → ✅ **SELESAI — K-036.** Section **dikeluarkan dari lingkup port**; pemilihan marketing mengikuti NB. **Temuan keamanan `.Password` GUGUR** — komponennya tidak diport | ~~§6~~ ✅ |
| ~~**RNW-5**~~ | ~~`PROSESCOPY` diport apa adanya atau ditinjau?~~ → ✅ **SELESAI — K-037.** Blok `pySaveSQL`-nya **DIBUANG**. ⚠️ Dicatat sebagai **perubahan perilaku yang disengaja**, bukan porting | ✅ |
| ~~**RNW-6**~~ | ~~Perlukah isi `Struktur_…xlsx` dibaca sebagai konfirmasi silang?~~ → ✅ **SELESAI 18 September 2026.** Sudah dibaca (kolom `XML` **tidak disentuh**). Hasil: `06-rnw\05-konfirmasi-silang-peta-struktur.md` | ✅ |

✅ **Seluruh enam butir K-033 tertutup.** Ringkasan hasil RNW-6:

- `[terverifikasi]` Ketiga rujukan (`GenerateNoPolicy`, `SumTreatyCapacity_Act`, `SearchJobID`)
  **tidak tercatat** di peta struktur — dan **instrumen terbukti peka** lewat empat nama kontrol dari
  20 berkas delta yang semuanya terdeteksi. `[dugaan]` **konsisten dengan** K-032, **bukan bukti**;
  K-032 **tidak dinaikkan tingkat buktinya**.
- `[terverifikasi]` **Tidak ada rule delta yang terlewat** discovery: **nol** nama di peta yang
  merupakan berkas NB absen dari RNW. Daftar 20 berkas delta **lengkap**.
- ⚠️ Koreksi pengukuran yang saya catat apa adanya: perbandingan pertama memberi "221 nama hilang",
  ternyata **kunci rule berformat lengkap** (`kelas + access group + nama`), bukan rule hilang.
  Setelah normalisasi angkanya runtuh menjadi 17 — lalu menjadi 2 setelah penyamaan huruf.

---

## K-034 · Gerbang masuk Renewal **DITETAPKAN** (menjawab RNW-2)

**Tanggal:** 18 September 2026 · **Menutup:** RNW-2 (K-033) · **Sifat: keputusan bisnis** yang
mengisi celah yang korpus memang tidak merekamnya

### Alur masuk yang ditetapkan

> User memasukkan **No. Polis + Renewal Date + Note** → klik **OK** → activity pembuat renewal
> dijalankan; di dalamnya ada proses **salin data** dari polis yang diinput.

### Yang mendukungnya dari korpus

`[terverifikasi]` `RNW Fac In\FlowAction\Renewal_FlowAct.xml`:

| Tag | Nilai |
| --- | --- |
| `pySectionReference` | `InputRenewal` |
| `pyPreProcessingActivity` | `InputOfferFacInEngineer_preACT` |
| `pyLocalActionActivity` | `SetValidateDate_PostAct` |
| `pyActionTransformRule` | `AddToListSuggestOfferFacIn_DT` |

Kedua activity terkonfirmasi sebagai `Rule-Obj-Activity` pada blok `Embed-Reference-Rule`.

Referensi layar: berkas gambar `DDL\halaman depan renewal.JPG` **ada** (30 KB).
⛔ **Isinya tidak dibaca sebagai fakta** — hanya keberadaannya yang dirujuk.

### Mengapa ini keputusan, bukan porting

`[terverifikasi]` `IsOfferFacIn` — gerbang masuk NB — **dirujuk nol kali** di RNW, dan
`Flow\InputRenewalFacultativeIn` **tidak dirujuk berkas mana pun**. Penentu "kapan sebuah kasus masuk
alur renewal" ada di **konfigurasi work type/portal yang tidak ikut terekspor**.

Karena tidak ada perilaku terekam yang dapat direproduksi, sistem baru **menetapkan** titik masuknya
secara eksplisit sesuai alur di atas. Ini **bukan** pelanggaran `CLAUDE.md` §1 — §1 melarang
mereka-ulang perilaku **yang terekam**; di sini tidak ada yang terekam untuk direka-ulang.

---

## K-035 · Pola konversi produksi renewal **sama seperti NB** (menjawab RNW-3)

**Tanggal:** 18 September 2026 · **Menutup:** RNW-3 (K-033) · **Sejalan:** K-031

**Keputusan: rancangan konversi produksi renewal mengikuti NB.** Tidak ada jalur konversi tersendiri
untuk renewal.

Sejalan dengan K-031 (renewal memakai ulang modul NB) dan dengan penutupan K-032, yang
`[terverifikasi]` membuktikan `SaveJsonPolicyFacIn_Act` **identik byte-per-byte** antara NB dan RNW.

⚠️ **Celah K-004 tetap terbuka.** Rule asli kelas `Work` (`serviceInsertArasapas_act` tanpa sufiks)
tetap tidak ada di korpus mana pun. Yang diputuskan di sini adalah **arah rancangan**, bukan
ditemukannya rule yang hilang.

---

## K-036 · `Section\EditMarketing` **DIKELUARKAN dari lingkup port** (menjawab RNW-4)

**Tanggal:** 18 September 2026 · **Menutup:** RNW-4 (K-033)

**Keputusan: `Section\EditMarketing` dan `FlowAction\EditMarketing` tidak diport.** Pemilihan
marketing pada renewal **mengikuti NB**.

### Konsekuensi: temuan keamanan **gugur**

`[terverifikasi]` Section itu mengikat `.ID` · **`.Password`** · `.QuotationData.MOID`, dan sempat
diangkat sebagai temuan keamanan yang memerlukan tinjauan `CLAUDE.md` §6.

**Karena komponennya tidak diport, tinjauan §6 atas `EditMarketing` tidak lagi relevan.** Temuannya
**dicatat dan ditutup** — bukan dilupakan.

⚠️ Ini **mengurangi permukaan** berkas delta renewal dari 20 menjadi **18 yang perlu dirancang**.
Kedua berkas `EditMarketing` tetap tercatat di dokumen discovery sebagai bagian korpus, dengan
penanda dikeluarkan dari lingkup.

---

## K-037 · ⚠️ `PROSESCOPY` **DIBUANG** — perubahan perilaku yang disengaja (menjawab RNW-5)

**Tanggal:** 18 September 2026 · **Menutup:** RNW-5 (K-033)

> ⛔ **Ini PERUBAHAN PERILAKU, bukan porting apa adanya.** Dicatat sebagai keputusan sadar terpisah
> justru karena `CLAUDE.md` §1 melarang "perbaikan diam-diam" — dan yang membedakan keduanya adalah
> **tercatat atau tidak**.

**Keputusan: blok `pySaveSQL` pada `RDBList\CariBusinessGID` yang memanggil `POOLDATA.PROSESCOPY`
tidak diport.**

### Alasan — empat sifat terbaca, seluruhnya `[terverifikasi]`

```sql
DECLARE
  vTHN_TREATY VARCHAR2(10);  vTOP_ID VARCHAR2(10);
  vIDTreatyYear VARCHAR2(10); errmsg VARCHAR2(4000);
BEGIN
  vTHN_TREATY   := '2007';
  vTOP_ID       := '10018';
  vIDTreatyYear := '1000134';
  POOLDATA.PROSESCOPY(vTHN_TREATY, vTOP_ID, vIDTreatyYear, errmsg);
  dbms_output.put_line(errmsg);
END;
```

1. **`POOLDATA.PROSESCOPY` tidak ada di basis data** — dikonfirmasi DBA (checklist §A.1).
2. Seluruh argumennya **literal keras** — tahun `'2007'`, dua ID tetap. Tidak ada parameter masuk.
3. Keluarannya ke **`dbms_output`** — kanal yang tidak dibaca aplikasi.
4. **Tidak berhubungan** dengan `pyBrowseSQL` di rule yang sama, yang mencari kelompok bisnis.

### Yang TIDAK ikut dibuang

⚠️ **Hanya blok `pySaveSQL` itu** yang dibuang. `pyBrowseSQL` pada rule yang sama — pencarian
`BusinessGroupID` dari tabel `business` (K-031, R-4) — **tetap diport apa adanya**.

### Catatan untuk rekonsiliasi

Karena prosedur yang dipanggil tidak ada, blok ini **tidak dapat pernah berhasil dieksekusi**.
`[dugaan]` membuangnya karena itu diperkirakan **tidak mengubah satu angka pun** pada paralel run.
Bila ternyata ada selisih yang menunjuk ke jalur ini, **keputusan ini adalah tersangka pertama** —
dan itulah gunanya ia dicatat bernomor.

---

## K-038 · Fitur **pilih-tertanggung DIBUANG** dari lingkup port (4 berkas)

**Tanggal:** 18 September 2026 · **Pola sama:** K-036 (`EditMarketing`)

**Keputusan: keempat berkas berikut tidak diport.**

| Berkas | Kelas `[terverifikasi]` |
| --- | --- |
| `Harness\ChooseInsured` | `ASM-FW-GISFW-Data-Quotation` |
| `Section\ChooseInsuredDtl` | `ASM-FW-GISFW-Data-Quotation` |
| `ReportDefinition\BrowseAccountInsuredEDM` | `ASM-FW-SFAGISFW-Work-Account` (SFA) |
| `Activity\SetDataInsuredEDM_Act` | `ASM-FW-SFAGISFW-Work-Account` (SFA) |

**Alasan:** renewal mengambil tertanggung dari **polis lama** — disalin saat pembuatan kasus — bukan
dari pemilihan CRM/SFA. Keempatnya **satu fitur**, dibuang bersama.

`[terverifikasi]` Dua di antaranya berkelas **SFA**, bukan Fac In, dan `BrowseAccountInsuredEDM`
memuat field CRM (`Industry`, `OppCountRep.Count`, `Org.Name`, `crmImageFileName`) — menguatkan bahwa
ia pemilih akun SFA, bukan komponen Fac In.

### Konsekuensi

**Delta renewal turun dari 20 menjadi 14 berkas** — dikurangi 2 (`EditMarketing`, K-036) dan 4 (fitur
ini). Keenamnya tetap tercatat di dokumen discovery sebagai bagian korpus, dengan penanda dikeluarkan
dari lingkup.

---

## K-039 · `.OldTSI` dan penyalinan data polis lama saat pembuatan kasus

**Tanggal:** 18 September 2026 · **Menutup:** pertanyaan `.OldTSI` · **Meluruskan:** narasi discovery

### Keputusan work owner

Saat **OK** ditekan, renewal membuat kasus baru dengan **menyalin seluruh data polis lama sebagai
nilai awal**. Perhitungan kemudian berjalan **seperti NB** di atas hasil salinan itu.

### Narasi discovery diluruskan — keduanya benar, pada titik yang berbeda

| Pernyataan | Berlaku pada |
| --- | --- |
| "Kaitan polis lama = **rujukan** `OldPolicyNo`, bukan pembawaan data" | **Perhitungan** — tetap benar (K-031) |
| "Data polis lama **disalin** jadi nilai awal" | **Titik pembuatan kasus** (K-034) |

**Tidak bertentangan:** salin di awal, lalu hitung seperti NB atas hasil salinan. K-031 tetap berdiri.

### ⚠️ Koreksi terhadap atribusi `.OldTSI` — dua mekanisme, nama yang menyesatkan

`[terverifikasi]` Pemeriksaan korpus 18 September 2026 menemukan bahwa `.OldTSI` **bukan** diisi oleh
proses salin saat pembuatan. Penulisnya adalah activity **perhitungan**, dari nilai terhitung:

```
Activity\CountPremiumNet       L2792 · L5505 · L7325   .OldTSI = Local.TSISpreadTotal
Activity\CountPremiumNetElse   L1781 · L3916 · L5751   .OldTSI = Local.TSISpreadTotal
Activity\SumTSIPremiSpreadedRNM_Act L8391              .OldTSI = Local.TsiLoLSpread + .OldTSI
Activity\GetEDMData            L942 · L1565            Local.OldTSI = .TSI   (untuk cetak R/I slip)
```

Seluruh penulisnya **berkas bersama** yang identik NB↔RNW.

**Mekanisme salin yang memang ada, tetapi memakai nama lain:** `[terverifikasi]`
`Activity\SetOldData` — **63 langkah**, kelas `ASM-FW-GISFW-Work`, **identik byte-per-byte NB↔RNW**,
dipanggil `Activity\GetPolicyData_ACT` — melakukan **52 penetapan properti ber-sufiks `Old`**:

```
…PropertyItemList(idx).TSIObjectItemOld  = @If(.TSIObjectItem!="", .TSIObjectItem, 0)
…TotalTSIPremiGrossList(idx).TSIOld      = @If(.TSI!="", .TSI, 0)
…CurrencyList(idx).PremiumOld            = @If(.Premium!="", .Premium, 0)
```

⛔ **`TSIOld` (sufiks) dan `OldTSI` (prefiks) adalah DUA properti berbeda.** Di `SetOldData`, string
`OldTSI` hanya muncul sebagai **nama halaman clipboard** (`<pxPageName>OldTSI</pxPageName>`), bukan
properti yang ditulis.

**Ringkasnya:**

| Mekanisme | Properti | Kapan | Sumber nilai |
| --- | --- | --- | --- |
| Snapshot nilai awal | `*Old` (**sufiks**) — 52 penetapan | saat data polis dimuat | nilai berjalan (`@If(.Xxx!="",.Xxx,0)`) |
| Nilai terhitung | `.OldTSI` (**prefiks**) | saat perhitungan premi/spreading | total spreading terhitung |

`[pertanyaan terbuka]` Apakah `SetOldData` adalah proses salin yang dimaksud keputusan ini, atau
mekanisme terpisah, **belum dipastikan** — `GetPolicyData_ACT` yang memanggilnya belum ditelusuri.
**Keputusan bisnisnya tidak berubah**; yang dikoreksi hanya atribusi teknis `.OldTSI`.

📌 **Ini jebakan "nama bukan bukti" yang ketiga kalinya** di proyek ini, setelah `isApproved`
(properti vs DecisionTable) dan `GetInsuredID` (Activity vs RDBList). Kali ini pembedanya **posisi
kata**: prefiks vs sufiks.

---

## K-040 · Empat section yang disisipkan layar renewal **DIWARISI dari NB** — bukan menggantung

**Tanggal:** 18 September 2026 · **Menutup:** `[pertanyaan terbuka]` "5 section" pada spec RNW ·
**Konsekuensi langsung dari:** K-031

`[terverifikasi]` Keempatnya **ada di `NB FacIn\Section\`** dan **tidak ada di `RNW Fac In\Section\`**:

| Section | NB | RNW |
| --- | :-: | :-: |
| `InputInwardFacultativeSuggest` | **ada** | tidak |
| `InwardFacIn` | **ada** | tidak |
| `OfferFacIn_NusaRe` | **ada** | tidak |
| `OfferFacIn_NusaRe_IsUW` | **ada** | tidak |

**Keputusan: keempatnya BUKAN rujukan menggantung.** Ini **konsekuensi wajar K-031** — renewal berbagi
basis kode dengan NB, dan layar renewal **menyisipkan section NB apa adanya**, persis seperti 1.907
berkas identik lainnya.

⚠️ **Pembedaan yang penting.** Sebuah section yang hadir di NB dan disisipkan dari RNW **bukan hal
yang sama** dengan rule yang dipanggil tetapi tidak ada di folder mana pun. Yang pertama adalah
**pemakaian bersama**; yang kedua **kekurangan**. Ketiga rujukan pada K-032 termasuk kategori kedua;
keempat section ini kategori pertama.

### Konsekuensi

1. Keempatnya **tidak menambah delta renewal** — tetap **14 berkas**.
2. Di sistem baru, komponen ini dibangun **sekali** (sebagai bagian NB) dan **dipakai kedua siklus**.
3. Tidak ada yang perlu diminta ke IT.

---

## K-041 · `InputDtlObject` **usang**, digantikan `InputDtlObject_FacIn`

**Tanggal:** 18 September 2026 · **Pola sama:** K-024 (`JUW_A`) · K-037 (`PROSESCOPY`) · K-032

**Keputusan work owner: `InputDtlObject` sudah tidak dipakai lagi, digantikan
`InputDtlObject_FacIn`.** Keempat sisipan yang tersisa adalah **sisa kelewat** saat penggantian.

### Bukti

`[terverifikasi]` Keberadaan berkas:

| Section | NB | RNW |
| --- | :-: | :-: |
| `InputDtlObject` (tanpa sufiks) | **tidak ada** | **tidak ada** |
| `InputDtlObject_FacIn` (dengan sufiks) | **ada** | **ada** |

`[terverifikasi]` Sisipan yang tersisa — `RNW Fac In\Section\InputRenewalDtl.xml`:

```
L2764  <pySection>InputDtlObject</pySection>
L2944  <pySection>InputDtlObject</pySection>
L3207  <pySection>InputDtlObject</pySection>
L3390  <pySection>InputDtlObject</pySection>
```

`[terverifikasi]` Di **berkas yang sama**, `InputDtlObject_FacIn` muncul **26 kali** — penggantinya
memang sudah dipakai luas, dan keempat sisipan lama tertinggal di antaranya.

### ⚠️ `InputDtlObject` ≠ `InputDtlObject_FacIn`

**Dua section berbeda, dibedakan hanya oleh sufiks.** Jangan disatukan, jangan dianggap salah ketik
dari yang lain. Ini **pelajaran "nama bukan bukti" yang keempat kalinya** di proyek ini — setelah
`isApproved` (beda tipe rule), `GetInsuredID` (nama identik, tipe berbeda), dan `TSIOld` vs `.OldTSI`
(beda posisi kata). Kali ini pembedanya **ada-tidaknya sufiks**.

### Sikap porting

⛔ **Diport apa adanya** (`CLAUDE.md` §1). **Jangan hapus keempat sisipan itu diam-diam.**

Status usangnya **dicatat, tidak dieksekusi**: karena section tujuannya tidak ada, sisipan itu
**tidak menghasilkan apa pun** dan menjadi tidak berdampak dengan sendirinya — tanpa satu baris pun
diubah. Mekanisme yang sama dengan K-019, K-024, dan K-032.

---

## K-042 · Fokus aktif bergeser ke **`Endorsment Fac In`** — NB & RNW selesai sampai tiket

**Tanggal:** 18 September 2026 · **Mengamandemen fokus aktif:** K-005 (NB) dan K-030 (RNW) ·
**Tidak membatalkan** pekerjaan keduanya

**Keputusan: siklus New Business dan Renewal dinyatakan SELESAI sampai tahap tiket.** Fokus aktif
berpindah ke siklus **Endorsement**, dengan target yang sama: discovery → keputusan → spec → tiket.

| | Keadaan |
| --- | --- |
| **NB** — discovery · spec · **16 tiket** (`05-tickets\01`…`16`) | ✅ **final**, tidak diulang |
| **RNW** — discovery · spec · **8 tiket** (`05-tickets\rnw\R01`…`R08`) | ✅ **final**, tidak diulang |
| **EDM** — `Endorsment Fac In\` | 🔵 **fokus aktif**, mulai dari discovery |
| Korpus Treaty tersendiri (`RNM_BRD\`) | tetap **tidak dibaca** (K-005) |

### ⛔ EDM **BUKAN** "NB dengan pintu masuk berbeda" — strategi RNW tidak berlaku

Ini pembeda terpenting keputusan ini, dan mengabaikannya akan membuat discovery EDM salah arah sejak
awal.

`[terverifikasi]` Perbandingan langsung terhadap NB:

| Ukuran | **RNW** | **EDM** |
| --- | ---: | ---: |
| Total berkas `.xml` | 1.927 | **2.061** |
| Bernama sama dengan NB | 1.907 | **1.707** |
| — **identik byte-per-byte** | **1.907** (100 %) | **0** (0 %) |
| — berbeda isinya | **0** | **1.707** |
| Eksklusif (tidak ada di NB) | 20 | **354** |

**RNW memakai ulang seluruh modul NB** karena berkasnya memang berkas yang sama (K-031).
**EDM tidak.** Nol berkas identik berarti **setiap berkas bernama sama tetap berbeda isinya** —
sehingga:

1. ⛔ **Tidak ada satu pun berkas EDM yang boleh diasumsikan sama dengan NB** hanya karena namanya
   sama. Angka nol itu **bukti** bahwa asumsi semacam itu akan salah.
2. **Strategi "pakai ulang modul NB" tidak berlaku.** Discovery EDM lebih dekat ke discovery NB dari
   awal daripada ke discovery RNW.
3. Delta EDM adalah **1.707 berkas berbeda + 354 eksklusif**, bukan 354 saja.

```powershell
# menghasilkan 1707 bernama-sama / 0 identik / 354 EDM-only
$nb="D:\migrasi\RNM\NB FacIn"; $ed="D:\migrasi\RNM\Endorsment Fac In"
$same=0;$diff=0;$onlyE=0
foreach($t in (Get-ChildItem $ed -Directory).Name){
  $A=New-Object 'System.Collections.Generic.HashSet[string]'
  Get-ChildItem "$nb\$t" -Filter *.xml -File -EA SilentlyContinue | ForEach-Object { [void]$A.Add($_.BaseName) }
  Get-ChildItem "$ed\$t" -Filter *.xml -File -EA SilentlyContinue | ForEach-Object {
    if($A.Contains($_.BaseName)){
      if((Get-FileHash "$nb\$t\$($_.BaseName).xml").Hash -eq (Get-FileHash $_.FullName).Hash){$same++} else {$diff++} }
    else { $onlyE++ } } }
"$($same+$diff) / $same / $onlyE"
```

### ⚠️ Angka bergantung pada kepekaan huruf — keduanya dicatat

`[terverifikasi]` Perbandingan **peka huruf** memberi **1.707 / 354**; **tidak peka huruf** memberi
**1.719 / 342**. Selisihnya **12 berkas** yang namanya **hanya berbeda kapitalisasi**:

| Tipe | Pasangan |
| --- | --- |
| `When` (11) | `IsAllRisk`↔`isAllRisk` · `IsAviationHull`↔`isAviationHull` · `IsBillboardNeonSyariah`↔`isBillboardNeonSyariah` · `IsBurglary`↔`isBurglary` · `IsCar`↔`IsCAR` · `IsElectronicEquipment`↔`isElectronicEquipment` · `isExclusion`↔`IsExclusion` · `IsFidelity`↔`isFidelity` · `IsGroupUWFac`↔`IsGroupUwFac` · `isMaintenance`↔`IsMaintenance` · `IsTravelTime`↔`isTravelTime` |
| `Section` (1) | `periode`↔`Periode` |

Angka **1.707 / 354 (peka huruf)** dipakai sebagai angka resmi, karena bila resolusi rule Pega memang
peka huruf maka kedua belas itu **rule yang berbeda**.

⚠️ `[pertanyaan terbuka]` **Apakah resolusi rule Pega di instalasi ini peka huruf?** Pertanyaan lama
yang belum terjawab, dan di EDM ia **berdampak langsung pada 12 berkas**. Perlu dijawab IT.

⚠️ **Catatan metodologis:** `Test-Path` di Windows **tidak peka huruf**, sehingga tidak dapat dipakai
membedakan varian kapitalisasi. Perbandingan wajib memakai `HashSet[string]` peka huruf.

### Yang diwarisi sebagai titik awal — **wajib dikonfirmasi ulang, jangan diterima buta**

| Sumber arsip | Yang dicatat |
| --- | --- |
| `_ARSIP-lintas-siklus\04-aturan\01-katalog-when.md` | Dari **39** rule `When` bercabang antar siklus, **seluruhnya melibatkan EDM** sebagai yang berbeda; NB↔RNW selalu sepakat. **Pola P1**: EDM membaca agregat tersimpan `.OfferFacIn.QuotationData.*`, sedangkan NB/RNW membaca halaman aktif `pyWorkPage.Quotation.*` |
| `_ARSIP-lintas-siklus\01-flow\03-alur-endorsement.md` | Alur endorsement + before/after image sebagian terpetakan — **titik awal**, bukan hasil final |
| **K-029** | `EdmType`: `1` Batal Sejak Semula · `2` Batal Prorata · `4` Penambahan/Pengurangan/Perubahan. Kode `3` **usang**, tetapi cabang `EdmType==3` **diport apa adanya** |
| **K-039** | `Activity\SetOldData` dan properti ber-sufiks `TSIOld` — **ditunda dari RNW ke EDM**; ditelusuri di fase ini |

⛔ Seluruh butir di atas berasal dari discovery lintas-siklus Prompt 2 dan **wajib dikonfirmasi ulang
ke korpus EDM**. Nol berkas identik membuat pewarisan temuan menjadi berisiko.

---

## K-043 · Lingkup kerja EDM = **866 berkas**, dengan 1.195 identik **dicek ulang perilakunya** (pilihan B)

**Tanggal:** 18 September 2026 · **Menutup:** E-Q4 (lingkup kerja EDM) ·
**Mengoreksi angka "0 identik" di K-042** · **Melanjutkan:** K-042

**Keputusan work owner: pilihan (B).** Lingkup kerja discovery & pembangunan EDM adalah **866 berkas**
= **512 berbeda** + **354 EDM-only**. Yang **1.195 identik** dipakai ulang dari modul NB **tetapi tidak
diterima buta** — perilakunya **dicek ulang saat integrasi**, tidak langsung dipercaya hanya karena teks
berkasnya identik.

### Koreksi angka terhadap K-042: "0 identik" → **1.195 identik / 512 berbeda**

K-042 mencatat **0 identik byte-per-byte** dari 1.707 bernama-sama. `[terverifikasi]` Angka itu benar
**untuk metode hash mentah**, tetapi hash mentah membandingkan **stempel waktu ekspor** yang selalu
berbeda antar folder — jadi ia mengukur kapan berkas diekspor, bukan apakah logikanya sama.

Angka yang berlaku sekarang **1.195 identik (70,0 %) / 512 berbeda (30,0 %)**, diperoleh dengan
membuang 23 tag volatile (stempel waktu, host, checksum, indeks) **sambil mempertahankan tag pembawa
kelas** (`pxRuleClassName`, `pxRuleFamilyName`). Metode & skrip: `07-edm\02-e1-pola-perbedaan-when.md`
§1.1 (daftar 23 tag + alasan) dan §1.7 (skrip reproducible).

⚠️ `[terverifikasi Kiro — uji-silang independen]` Kiro tidak memakai metode 23-tag Claude; Kiro
menguji titik-sampel dengan metode berbeda (banding isi `pyConditionValue1String` + hitung pola).
Hasil sampel **konsisten dengan arah** 1.195/512, tetapi **angka penuh 512 belum direproduksi Kiro di
korpus penuh**. Karena itu, sebagai **syarat pilihan (B)**:

> Sebelum 512 dipakai memotong pekerjaan, skrip §1.7 **wajib dijalankan ulang** dan hasilnya
> dicocokkan. Ini pengaman yang sama seperti caveat "daftar 13 rule wajib ditelaah satu per satu" di
> K-001.

### Arti operasional pilihan (B)

| Kelompok | Jumlah | Perlakuan |
| --- | ---: | --- |
| Identik (metode 23-tag) | **1.195** | Pakai ulang modul NB, **TETAPI** ditandai untuk **cek ulang perilaku saat integrasi** — bukan diterima buta |
| Berbeda isinya | **512** | Masuk lingkup discovery/pembangunan EDM |
| EDM-only | **354** | Masuk lingkup discovery/pembangunan EDM |
| **Total lingkup kerja EDM** | **866** | Basis discovery E-2…E-6 |

### Kenapa 1.195 identik **tidak** diterima buta (inti pilihan B)

`[terverifikasi Kiro]` **Teks identik ≠ perilaku identik**, dan EDM sudah membuktikannya:

- Enam predikat COB (`IsFire`, `IsPA`, `IsAneka`, `IsMarineCargo`, `IsGolfInsurance`, `IsMBU`) di EDM
  membaca **hasil query** (`OutData.pxResults`, terbaca di EDM: IsFire 28×, IsAneka 128×, IsMBU 23×;
  **0× di NB**), sedangkan NB membaca properti/agregat atau mendelegasi ke sub-rule. Nama sama, tipe
  sama, tetapi **sumber data yang dibaca berbeda**. Lihat pola P2 di
  `07-edm\02-e1-pola-perbedaan-when.md`.
- Karena itu sebuah berkas yang **teksnya identik** pun bisa berperilaku beda bila **data yang masuk**
  ke berkas itu berasal dari jalur berbeda. Cek ulang perilaku saat integrasi adalah pengaman terhadap
  kelas kesalahan ini.

⚠️ **Yang dilarang oleh pilihan (B):** menyatakan sebuah berkas EDM "beres" hanya karena hash/teksnya
cocok dengan NB, tanpa memverifikasi jalur data yang memberinya masukan. Kalau paralel run EDM nanti
selisih pada modul yang "identik", **inilah tersangka pertamanya** — bukan cacat implementasi.

### Dua koreksi Kiro terhadap laporan discovery E-1 (agar tidak diport keliru)

`[terverifikasi Kiro]` Dua frasa di laporan E-1 perlu dikoreksi sebelum masuk spec:

1. **"Beda kelas deklarasi"** untuk 6 predikat COB — **tidak akurat**. `pxObjClass` keenam rule tetap
   `Rule-Obj-When` di NB **dan** EDM (sama). Yang berbeda adalah **kelas sumber data yang dibaca
   kondisi** (`OutData.pxResults`), bukan kelas tempat rule dideklarasikan. Ditulis benar: *"beda kelas
   sumber data yang dibaca", bukan "beda kelas deklarasi".*
2. **"NB mendelegasikan ke sub-rule"** — hanya **sebagian**. Terbukti untuk `IsFire` (20×) dan
   `IsAneka` (125×); `IsPA`/`IsMBU`/`IsMarineCargo`/`IsGolfInsurance` di NB **0× delegasi** (baca
   properti langsung). Kesimpulan "registry predikat NB tak dapat dipakai ulang apa adanya untuk COB
   EDM" **tetap berlaku**, hanya alasannya bukan "delegasi vs literal" yang seragam.

### Konsekuensi terhadap dokumen lain

| Dokumen | Perubahan |
| --- | --- |
| `07-edm\02-e1-pola-perbedaan-when.md` | angka 1.195/512 dikonfirmasi sebagai basis lingkup; dua koreksi Kiro di atas ditambahkan sebagai catatan |
| Discovery E-2…E-6 | fokus pada 866 berkas (512 + 354); 1.195 identik = pakai-ulang-dengan-cek-ulang |
| Spec EDM (nanti) | resolver jalur properti & registry predikat COB **per-rule**, tidak diwarisi buta dari NB |

---

## 🔄 Amandemen K-043 — 19 September 2026: kontrak volatile + **isi blok `pzIndexes` diabaikan** → lingkup **866 → 708**

**Menutup:** E-Q23 · **Mengamandemen:** kontrak metode K-042/K-043 (definisi "berkas identik") ·
**Tidak membatalkan** pilihan (B): 1.195 (kini 1.353) identik tetap dicek ulang perilakunya saat integrasi.

**Keputusan work owner: kontrak volatile diperluas menjadi "23 tag volatile + isi setiap blok
`pzIndexes` diabaikan".** Lingkup kerja EDM turun dari **866** menjadi **708 berkas**
(**354 berbeda** + 354 EDM-only).

### Angka — tereproduksi INDEPENDEN oleh work owner (bukan diteruskan dari laporan)

`[terverifikasi work owner]` Skrip 23-tag disalin persis dari `07-edm\02-e1-pola-perbedaan-when.md`
§1.1/§1.7 (daftar 23 tag volatile + 6 tag berkunci di-`<TS>`-kan + urut Ordinal sebelum hash SHA-256),
dijalankan ulang atas **1.707** pasang berkas bernama-sama. Hasil **cocok bit-per-bit** dengan laporan
Claude:

| Kontrak | Identik | Berbeda |
| --- | ---: | ---: |
| 23-tag (K-042/K-043 lama) | 1.195 | **512** |
| 23-tag **+ buang isi `pzIndexes`** | 1.353 | **354** |
| **Berbeda hanya karena `pzIndexes`** | | **158** |

⚠️ Verifikasi work owner sebelumnya (metode volatile ad-hoc) sempat **gagal mereproduksi** angka ini
(dapat 0) — sebabnya kontrak ad-hoc membuang terlalu sedikit tag, tidak menormalkan `<TS>` pada tag
berkunci, dan tidak mengurut Ordinal sebelum hash. Dengan **daftar 23 tag yang benar**, 158
tereproduksi persis. Ini menegaskan prinsip §0 dokumen E-1: angka hanya dapat dipercaya bila daftar
tag-nya eksplisit dan tertutup.

### Dasar struktural — mengapa `pzIndexes` aman diabaikan

`[terverifikasi work owner]` Blok `pzIndexes` di korpus:

- Tag sebenarnya **`<pzIndexes REPEATINGTYPE="PropertyGroup">`** (beratribut) — `<pzIndexes>` polos = 0.
  (Pencarian tag-polos memberi 0 dan sempat menyesatkan; ini jebakan yang sama dengan
  `<pxRuleReferences REPEATINGTYPE=...>`.)
- Di `Activity\GetLimitAkseptasi_Act.xml`: **82 blok** di NB **dan** 82 di EDM; isinya berselisih hanya
  **offset penomoran** (NB `RuleReference=86…`, EDM `=6…`), bukan isi. Cocok dengan selisih ukuran
  **78 byte dari 544.882** (0,014 %).
- Struktur seluruh blok pzIndexes EDM: **elemen non-`rowdata` = 0, nilai non-bilangan-bulat = 0,
  blok yang memuat tag pembawa kelas (`pxRuleClassName`/`pxRuleFamilyName`/`pyDesignatedClass`) = 0.**
  Jadi mengabaikan isinya **tidak menyembunyikan perbedaan kelas** — kesalahan yang menggugurkan
  angka #2 di §0 tidak terulang.
- 5 uji silang tetap lolos: `FlagOldData` tetap **BEDA** (kelas terjaga); 3 `When` + `D_AnekaList`
  tetap identik.

### Arti operasional

| Kelompok | Jumlah | Perlakuan |
| --- | ---: | --- |
| Identik (kontrak baru: 23-tag + pzIndexes diabaikan) | **1.353** | Pakai ulang modul NB, **tetap dicek ulang perilaku saat integrasi** (pilihan B tak berubah) |
| Berbeda isinya | **354** | Masuk lingkup discovery/pembangunan EDM |
| EDM-only | **354** | Masuk lingkup discovery/pembangunan EDM |
| **Total lingkup kerja EDM** | **708** | Basis discovery E-2…E-6 |

⚠️ **Kebetulan yang perlu ditegaskan:** angka "354 berbeda" hasil amandemen ini **kebetulan sama**
dengan cacah "354 EDM-only" — **keduanya tidak berhubungan**. Total lingkup = 354 + 354 = 708.

### Batas keputusan — apa yang terbukti vs apa yang diputuskan

- **Terbukti (dapat dihitung):** angka 158, struktur pzIndexes (0 tag kelas), reproduksibilitas.
- **Diputuskan (kontrak, bukan fakta):** *mengabaikan* isi pzIndexes sebagai kebijakan pembandingan.
  Dasarnya kuat (bukti struktural di atas), tetapi ia tetap keputusan work owner. Jaring pengaman
  pilihan (B) — cek ulang perilaku saat integrasi — **tetap berlaku** untuk 1.353 identik.

### Yang berdiri SENDIRI, tidak bergantung amandemen ini

`[terverifikasi]` **4 dari 6** berkas EDM+RNW-bukan-NB (`GetBusinessGroup_Act`, `CariBusinessGID`,
`BrowseAccountInsuredEDM`, `SetDataInsuredEDM_Act`) sudah **identik di bawah kontrak 23-tag lama** —
tanpa perlu amandemen pzIndexes. Ini **menjatuhkan premis E-Q3/E-Q25** (dugaan pilih-tertanggung
"berbeda" di EDM berasal dari hash mentah). Hanya "6 dari 6 identik" yang bergantung amandemen ini
(2 sisanya, `ChooseInsured` + `ChooseInsuredDtl`, beda hanya di blok indeks).

### Konsekuensi terhadap dokumen lain

| Dokumen | Perubahan |
| --- | --- |
| `07-edm\02-e1-pola-perbedaan-when.md` | kontrak metode diperluas: 23 tag + isi `pzIndexes` diabaikan; angka lingkup 512→354 |
| `07-edm\05-e4-edm-only.md` | E-Q23 ditutup dengan angka tereproduksi; lingkup 866→708 |
| Discovery E-2…E-6 | fokus pada **708** berkas (354 berbeda + 354 EDM-only) |
| E-Q3 / E-Q25 (pilih-tertanggung) | premis "berbeda" gugur untuk 4 berkas; nasib fitur pilih-tertanggung EDM **tetap** pertanyaan work owner terpisah |

---

## K-044 · Lingkup EDM: Life **ikut**, pilih-tertanggung **dibuang**

**Tanggal:** 19 September 2026 · **Menutup:** W-1 (E-Q26/E-Q31), W-2 (E-Q25/E-Q3)

- **W-1 — EDM Life IKUT dikerjakan.** Endorsement produk Life masuk lingkup, meski jalurnya terpisah
  (melompati tangga akseptasi, hanya Life jalan untuk `EdmType==1`, `SetOldData` keluar bila `IsLife`,
  jalur produksi Life lewat `SaveFacinLive_Act`/`SaveFacinSpreadLife_Sql`). Life **tidak** dikeluarkan.
- **W-2 — Fitur pilih-tertanggung EDM DIBUANG.** Tertanggung diambil dari **polis yang di-endors**,
  tidak dipilih ulang. **K-038 (dibuang untuk RNW) berlaku juga untuk EDM.** Premis E-Q3/E-Q25 ("EDM
  berbeda dari RNW") gugur — `[terverifikasi]` 4 dari 6 berkas identik di bawah kontrak 23-tag; 2 sisanya
  (`ChooseInsured`, `ChooseInsuredDtl`) beda hanya blok `pzIndexes`. Keempat berkas pilih-tertanggung
  (`ChooseInsured`, `ChooseInsuredDtl`, `BrowseAccountInsuredEDM`, `SetDataInsuredEDM_Act`) **tidak
  diport** — dicatat, bukan dihapus diam-diam.

---

## K-045 · Tangga akseptasi Putaran 2 (`GetLimitAkseptasi_Act2`) **di-remark** — Putaran 1 tetap hidup

**Tanggal:** 19 September 2026 · **Menutup:** W-3 (E-Q28), sebagian mengamandemen dugaan E-Q28

**Sumber bukti: UI Pega work owner** (status remark `//` tidak terserialisasi ke ekspor XML — sama
seperti `pyStepsPreCondition`, UI Pega mengalahkan pembacaan XML).

- `GetLimitAkseptasi_Act2` ("Putaran 2") **di-remark di KEDUA pemanggilnya**:
  `InputOfferFacInEngineerUW_preACT` (step 6, label `//`) **dan** `SetValidateDate_PostAct`.
  → **Putaran 2 tidak dipakai lagi di EDM.**
- ⚠️ **`GetLimitAkseptasi_Act` ("Putaran 1") TETAP HIDUP.** `[terverifikasi Kiro]` masih dipanggil
  `InputAddendumFacIn_PreAct`, `InputAddendumFacultativeIn` (flow utama EDM), dan `SetValidateDate_PostAct`
  (di `InputOfferFacInEngineerUW_preACT` step 5 memang di-remark, tapi pemanggil lain tidak). **Jangan
  divonis mati.**
- **`SaveFacinProdEDMBonding_Act` (di `DDL\`, E-Q28):** status jalur produksi Bonding **tetap terbuka** —
  belum dikonfirmasi apakah cabang `IsBonding` aktif. Tidak dihapus (K-006).

> ⛔ Ini **berbeda** dari amandemen K-006 (yang memulihkan cabang karena rule tiba di DDL). Di sini
> Putaran 2 **di-remark** (dinonaktifkan pemanggilnya), bukan dipulihkan. Vonis usang menyeluruh atas
> rule lain butuh pengecekan tiap pemanggil di UI Pega.

---

## K-046 · Kode usang EDM — diport apa adanya, kecuali **satu** yang diperbaiki

**Tanggal:** 19 September 2026 · **Menutup:** W-4 (12 butir dari `07-edm\06-e5` §8 + E-3)

Prinsip CLAUDE.md §1: keanehan sistem lama ditiru apa adanya (kandidat perbaikan milik bisnis).
Setelah ditinjau work owner, dari 12 butir: **9 diport apa adanya, 1 diperbaiki, 2 bukan keanehan.**

| # | Butir | Bukti | Keputusan |
| --- | --- | --- | --- |
| A.1 | `RateOld` self-referential `@If(.RateOld!="",.Rate,0)` (6× di `SetOldData.xml`) | `[terverifikasi]` | **BENAR / diport apa adanya** (bukan bug menurut work owner) |
| A.2 | `EDMPremiMenjadi` banyak bentuk rumus | `[terverifikasi]` | diport apa adanya (memang berbeda) |
| A.3 | rumus juga di `CountPremiEDM_DT` (DataTransform) | work owner | diport apa adanya |
| A.4 | label langkah "batal" tak cocok kondisi | — | **abaikan label** (kondisi yang berlaku, §3) |
| **A.5** | `PremiNusantaraReOld` diisi string `"0"` (50 field lain numerik) di `SetOldData.xml` L6294/L7527 | `[terverifikasi]` | **⚠️ DIPERBAIKI → angka `0` di sistem baru.** Satu-satunya butir yang diperbaiki, bukan ditiru. |
| B.6 | `PCT_BROKERGARE_FEE` | work owner | **BUKAN keanehan** — kolom penyimpan nilai dalam persen, ejaan tetap apa adanya |
| B.7 | kolom "sesudah" tanpa sufiks | work owner | diport apa adanya (kemungkinan lupa label) |
| B.8 | dua semantik `PRODKE` | `[terverifikasi]` §4.2 | **tertutup** — `GetProdKeOldData_SQL` (order by TGL_INPUT desc) yang sesuai (P-10) |
| C.9 | label "MBU" vs deskripsi | `[terverifikasi]` cuma label | **salah label saja** (klaim "gerbang IsMarineCargo" dicabut — tak terbukti) |
| C.10 | guard fac out hanya `Type=="7"` (`InsertFacoutProductionEDM.xml` L535/L1945/L1900) | `[terverifikasi]` | diport apa adanya |
| C.11 | `JSON_POLIS_MONITORING` hanya prefiks `EDMT-` (`serviceInsertArasapasEDM_act.xml` L3309/L3543) | `[terverifikasi]` | **BENAR** — memang sengaja dilewati untuk EDMT- |
| C.12 | "cabang IsMBU ganda" | — | **DICABUT** — work owner sudah cek, tidak ada yang salah |

⚠️ **A.5 — konsekuensi rekonsiliasi.** Karena diperbaiki (string `"0"` → angka `0`), saat paralel run
nilai bisa tampak "beda tipe" dari sistem lama walau angkanya sama. Ini **perbaikan sadar**, bukan
salah-port — dicatat agar selisih itu tidak dikira bug (K-007 rekonsiliasi eksak).

---

## K-047 · **Seam 4 disetujui** — `services/endorsement.PrepareBeforeImage`

**Tanggal:** 19 September 2026 · **Menutup:** usulan seam di `04-spec\05-spec-edm-before-image.md`
(§Testing Decisions) · **Menambah:** 1 seam ke 3 seam NB (16 Sept 2026)

**Keputusan work owner: Seam 4 DISETUJUI.**

Modul baru `services/endorsement` mendapat **satu** seam pengujian:
`PrepareBeforeImage(polisLama, tanggalEndorsement)` — satu pintu masuk yang mencakup **ketiga lapis
before-image sekaligus** (A `OfferFacIn.OldData`, B properti `*Old`, C penanda `.IsOldData="old"`).

### Batas keputusan (dijaga)

- **Lapis A/B/C BUKAN seam terpisah** — memberi seam per lapis akan mengunci pembagian internal yang
  belum tentu bertahan. Konsisten prinsip NB (226 predikat = 1 seam, `03-spec` L312).
- **Selisih TIDAK dapat seam baru** — diuji lewat **Seam 3 (`services/premium.Calculate`)** yang sudah
  ada, diperluas menerima masukan endorsement (nilai lama + rasio prorata). **Total tetap 4 seam,
  bukan 6.**
- **3 seam NB tidak disentuh** (`premium.Calculate` · `acceptance.Next` · `rules.Eval`,
  disetujui 16 Sept 2026). `[terverifikasi Kiro]` ketiganya ada di `03-spec-modul-terverifikasi.md`.

### Alasan (vs alternatif menolak)

Bila ditolak, lapis B & C diuji **tidak langsung** lewat Seam 3 — bug pengisian before-image baru
ketahuan saat hitung premi (telat, sulit dilacak). Seam 4 membuat before-image dapat diuji langsung,
tanpa memecah jadi 3 seam.

### Konsekuensi terhadap dokumen lain

| Dokumen | Perubahan |
| --- | --- |
| `04-spec\05-spec-edm-before-image.md` | Seam 4 berubah dari "usulan menunggu persetujuan" → **disetujui** |
| `04-spec\03-spec-modul-terverifikasi.md` | catatan "3 seam, tak boleh nambah tanpa keputusan baru" → kini 4 seam (K-047 = keputusan barunya) |

Langkah berikutnya: `/to-tickets` (dijalankan work owner manual) untuk memecah spec before-image jadi
tiket. Tidak dijalankan otomatis.

---

## K-048 · Modul selisih EDM: `ProrateEDMEnd` port apa adanya · `pyDisabled=true` = dilewati

**Tanggal:** 19 September 2026 · **Menutup:** §8.1 & §8.2 di `04-spec\06-spec-edm-selisih.md` (bahan)

### §8.1 — `ProrateEDMEnd` **PORT APA ADANYA (opsi a)**

**Keputusan work owner: satu properti `ProrateEDMEnd`, ditimpa berurutan** — bukan dipecah jadi
beberapa rasio bernama berbeda (opsi b ditolak).

`[terverifikasi]` `ProrateEDMEnd` ditulis di beberapa tahap yang saling menimpa satu properti yang sama:
modul before-image (modul 1) menyetelnya dari selisih DateTime; `CountPremiEDM_DT` langkah 9
menimpanya dengan selisih hari; langkah 11.3 memaksanya `1` bila `EdmType==1`; langkah 12.1 menulis
rasio awal ke properti akhir bila `EdmType==2`. `[terverifikasi Kiro]` `ProrateEDMEnd` muncul di
**27 berkas** EDM (tulis+baca).

- Alasan pilih (a): K-007 (rekonsiliasi eksak) mewajibkan tiru perilaku lama; sistem lama memang
  menimpa berurutan. Memecah jadi rasio bernama beda = "perbaikan diam-diam" (§1) yang berisiko
  selisih tak terjelaskan.
- ⚠️ **Konsekuensi:** nilai `ProrateEDMEnd` **rapuh terhadap urutan eksekusi** antar modul. **Spec
  modul 1 (before-image) WAJIB diberi catatan:** nilai `ProrateEDMEnd` yang disetelnya bersifat
  **sementara** — dapat ditimpa modul selisih. Urutan eksekusi modul 1 → modul selisih harus dijaga.

### §8.2 — `<pyDisabled>true</pyDisabled>` pada DataTransform = **DILEWATI (tidak dieksekusi)**

**Keputusan/konfirmasi work owner (UI Pega): aksi ber-`pyDisabled=true` DILEWATI.**

⚠️ **Beda dari `pyStepsPreCondition=false`** (K-045/P-11 area: itu **tetap jalan**). Dua tanda mirip,
arti berlawanan — jangan tertukar:
| Tanda | Di mana | Arti |
| --- | --- | --- |
| `pyStepsPreCondition=false` | langkah Activity | **tetap jalan** (langkah tak bersyarat) |
| `pyDisabled=true` | aksi DataTransform | **dilewati** (dinonaktifkan) |

`[terverifikasi Kiro]` `CountPremiEDM_DT.xml`: **23** aksi `pyDisabled=true` (`false`=0) — 22 berurutan
L478–L1075 (blok "langkah 1") + 1 di L3167.

- **Konsekuensi:** 22 aksi Set di blok awal `CountPremiEDM_DT` **tidak dieksekusi** — yang berlaku
  langkah penggantinya. Spec modul selisih **mengabaikan blok disabled**, mengikuti langkah aktif.
  Ini **bukan** kode usang yang diport (beda dari K-046) — memang dinonaktifkan by design di Pega.

### Konsekuensi terhadap dokumen lain

| Dokumen | Perubahan |
| --- | --- |
| `04-spec\05-spec-edm-before-image.md` | tambah catatan: `ProrateEDMEnd` nilai sementara, dapat ditimpa modul selisih (§8.1) |
| `04-spec\06-spec-edm-selisih.md` (bahan) | §8.1 = port apa adanya; §8.2 = abaikan 22 aksi disabled di CountPremiEDM_DT |

---

## K-049 · **Seam 5 disetujui** — `services/endorsement.OpenCase` · koreksi lingkup `BrowseOpenProteksiEdm_RD`

**Tanggal:** 19 September 2026 · **Menutup:** §8 & §9 di `04-spec\07-spec-edm-alur-masuk.md` (bahan)

### §8 — Seam 5 DISETUJUI: `services/endorsement.OpenCase`

Modul alur masuk EDM mendapat **satu** seam pengujian: `OpenCase` — mencakup pembuatan case
(dua kelas, prefiks EDM-, tautan tiga arah EndorsementID/EDMHandle/EDMHandle2), **6 gerbang penolakan**,
**4 klep pembatal**, bypass "EDM RI Slip", dan penempatan fase Policy — dalam satu pintu masuk.

**Dasar (terverifikasi Kiro):** 6 gerbang penolakan adalah **kondisi inline `@LengthOfPageList(...)>0`**
di langkah Activity (`SetErrorBatalEndorsement_Act`), **bukan rule When** — 52× `@LengthOfPageList` di
berkas itu. Karena itu **Seam 1 (`rules.Eval`) tidak dapat mengujinya**, dan 3 seam lain
(premium/acceptance/before-image) semua berjalan **setelah** case lahir. Tanpa Seam 5, pembuatan &
penolakan case tidak teruji langsung.

**Alasan:** penolakan endorsement adalah **keluaran bisnis yang sah** (bukan error teknis) — bila tak
teruji, salah-implementasi gerbang muncul sebagai selisih case saat paralel run (K-007). Total seam
kini **5** (3 NB + Seam 4 before-image K-047 + Seam 5 open-case). Modul selisih & modul lain **tidak**
dapat seam baru — disiplin dijaga (seam ditambah hanya saat ada celah uji nyata di modul baru).

### §9 — koreksi lingkup `BrowseOpenProteksiEdm_RD` (mengoreksi K-046/P-10)

`[terverifikasi work owner]` `BrowseOpenProteksiEdm_RD` (ReportDefinition, berkasnya nihil di korpus,
dirujuk 4× di `SetErrorBatalEndorsement_Act`) **DIPAKAI untuk pengecekan polis saat membuat EDM** —
membaca tabel `OPENPROTEKSI_EDM`.

**Koreksi:** sebelumnya (K-046/P-10) `OPENPROTEKSI_EDM` dicatat "di luar lingkup EDM". Klarifikasi
work owner: **tabelnya** memang dipakai bersama menu lain, **tetapi rule pengecekannya adalah bagian
dari alur EDM** (klep pembatal 4 gerbang: Param.Type 1=Batal, 2=Klaim, 3=Pembayaran, 4=Renewal).
Jadi:
- Rule pengecekan (`BrowseOpenProteksiEdm_RD`) → **masuk lingkup EDM**, diport.
- Isi tabel `OPENPROTEKSI_EDM` → **di luar korpus**, dikelola menu lain/DBA. `[pertanyaan terbuka]`
  bentuk & sumbernya — tidak ditebak.
- Konsekuensi (terverifikasi): klep hanya menutup **4 dari 6** gerbang; gerbang 2 (GetListEdm, pakai
  pxRetrieveReportData sebagai daftar) & gerbang 6 (fac out) **tanpa klep** — tidak dapat dibatalkan.
  Diport apa adanya.

### Konsekuensi terhadap dokumen lain

| Dokumen | Perubahan |
| --- | --- |
| `04-spec\07-spec-edm-alur-masuk.md` | Seam 5 dari "usulan" → disetujui; §9 lingkup dikoreksi |
| `04-spec\03-spec-modul-terverifikasi.md` | catatan jumlah seam: 3 → 5 (K-047 + K-049) |
| `_PAKET-…-DBA-IT-PRODUCT.md` P-10 | baris OPENPROTEKSI_EDM "di luar lingkup" dikoreksi: rule pengecekan masuk EDM, isi tabel dari DBA/menu lain |

Langkah berikutnya: `/to-tickets` (work owner manual) setelah modul-modul EDM cukup dispec.

---

## K-050 · Registry predikat COB = **satu registry, sumber data per-rule** (E-Q12) · **PRODKE aman** (E-Q30 tuntas)

**Tanggal:** 19 September 2026 · **Menutup:** E-Q12 (bentuk registry predikat), E-Q30 (kerapatan PRODKE)

### E-Q12 — Pilihan A: satu registry `internal/rules`, sumber data atribut per-rule (Seam 1 diperluas)

**Keputusan work owner: satu registry untuk ketiga siklus; tidak ada registry kedua; tidak ada seam
baru (total tetap 5).**

`rules.Eval` (Seam 1) diperluas: tiap predikat mendeklarasikan **sumber datanya** (properti kasus
atau kolom hasil query); resolver menerima konteks berisi keduanya.

**Fakta bisnis yang memperkuat (work owner):** **COB dipilih manual SEKALI saat membuat NB** (di
halaman depan NB); **RNW dan EDM menyalin data dari polis NB**, tidak memilih ulang. Jadi COB adalah
**satu nilai yang diwariskan**, bukan klasifikasi ulang per siklus.

`[terverifikasi work owner via UI]` Halaman depan pembuatan NB memuat field **`Class Of Business`**
(dropdown, wajib) — inilah titik pilih COB. Field lain di layar: `Type Of Inward` (dropdown wajib),
`Business Status` = "New Business", `Phase` = "Proposal", `Stage` = "Opportunity" (fase awal, bukan
Policy — konsisten dengan NB mulai di fase penawaran sedangkan EDM langsung Policy), `Business Prospect
Name`, `Group Business`, `Opportunity Source`, `Estimated Closing Date`, `Description`. **Nilai
field identitas (Owner dsb.) tidak disalin** (PII, §3.5). Konsekuensi:

- Registry predikat COB **secara logika sama di NB/RNW/EDM** — yang berbeda hanya **jalur baca runtime**:
  - NB baca properti kasus (nilai yang baru dipilih);
  - EDM baca hasil query `GetBusinessType_Sql` → `OutData.pxResults(1).CARI2` = BusinessType polis lama
    (yakni nilai yang dulu dipilih di NB lalu disalin). `[terverifikasi]` `GetBusinessType_Sql` EDM-only,
    dirujuk `SetErrorBatalEndorsement_Act` langkah 7.
- Karena hasilnya **nilai COB yang sama**, memisahkan registry (Pilihan B) hanya menduplikasi 196 dari
  202 predikat yang identik — ditolak.
- `[terverifikasi]` cabang **yang dieksekusi** (`pyConditionValue1String` memuat `OutData.pxResults`):
  IsFire 5 · IsAneka 25 · IsPA 1 · IsMBU 4 · IsMarineCargo 1 · IsGolfInsurance 1 = **37 cabang**
  (BUKAN 203 token mentah). Implementasi mengikuti 37.

### E-Q30 — **PRODKE aman, `COUNT-1` boleh dipakai** (TUTUP)

**Keputusan/konfirmasi work owner: baris polis yang sudah masuk TIDAK PERNAH dihapus.** Maka `PRODKE`
selalu rapat dari 0 tanpa lubang → rumus `PRODKE = COUNT(NOPOLIS)-1` **valid**.

Ini menutup E-Q30 **sepenuhnya**, untuk **ketiga** pemakaian yang bergantung padanya:
- klasifikasi COB (`GetBusinessType_Sql`, temuan modul predikat),
- nilai lama before-image (`GetEDMOldData_SQL`),
- penomoran versi (`GetProdKeEDM_SQL`).

⚠️ Catatan porting: kedua jalur (`COUNT-1` dan `ORDER BY TGL_INPUT DESC`) sama-sama valid karena tak
ada lubang — diport apa adanya, tidak "diseragamkan" (K-007). Bila paralel run kelak menemukan selisih
versi, asumsi "tak ada penghapusan" adalah tersangka pertama.

### Risiko yang ditutup (sebelumnya diangkat sebagai baru)

Temuan modul predikat: `GetBusinessType_Sql` memakai `COUNT-1`, sehingga kerapatan PRODKE memengaruhi
**bukan hanya nilai lama, tetapi juga klasifikasi COB** (→ cabang premi, skala K-018, jalur produksi).
**Ditutup oleh jaminan "tak ada penghapusan" di atas** — tidak jadi blocker DBA.

### Catatan porting JSON Oracle

`GetBusinessType_Sql` membaca `a.Data_json.QuotationData.BusinessType` (dot-notation JSON Oracle) →
di sistem baru diterjemahkan ke `JSON_VALUE(...)` atau kolom tersendiri (sama seperti `GetEDMOldData_SQL`,
E-3). `[terverifikasi]` dot-notation ada di query.

### Konsekuensi terhadap dokumen lain

| Dokumen | Perubahan |
| --- | --- |
| `04-spec\08-spec-edm-predikat-when.md` | registry = Pilihan A; 37 cabang dieksekusi; E-Q5 tuntas (GetBusinessType_Sql) |
| `03-spec-modul-terverifikasi.md` | Seam 1 `rules.Eval` diperluas: konteks berisi properti kasus + kolom query |
| E-Q30 (di E-5 §4.2) | ditutup — PRODKE aman, tak ada penghapusan |
| Paket DBA | risiko PRODKE→COB **tidak jadi** ditambahkan (tertutup) |

---

## K-051 · Modul EDM-only (350 diport) + perhitungan premi: **rumus dasar dipakai ulang, Seam 3 = 4 bentuk**

**Tanggal:** 19 September 2026 · **Menutup:** modul 5 (354 EDM-only), banding premi EDM↔NB, satuan `ProRatePercent`

### Lingkup modul 5 — 350 berkas diport

`[terverifikasi]` 354 EDM-only − 4 pilih-tertanggung (K-044 dibuang) = **350 diport**. Neraca: 4 dibuang,
27 When (modul 4), 28 sudah dibahas modul 1/2/3/E-5, **295 sisa modul ini**. 5 berkas tak-dirujuk (E-Q24)
tidak divonis usang (K-006). Tidak ada kelompok jadi modul domain baru; yang menyentuh premi → services/premium
& spreading (Seam 3), pendukung/tampilan → handler+React, query → repository. **Tidak ada seam baru.**

### Perhitungan premi — rumus dasar SAMA lintas siklus, EDM memakai ulang mesin NB

`[terverifikasi]` **`CalculatePremi_Act` (NB) BUKAN baseline** — kelasnya `ASM-FW-GISFW-Data-PolicyTreatyIn`
(Treaty), `.Rate/.TSI/.Premium = 0/0/0`. Jebakan penamaan ke-8. Baseline sebenarnya = **mesin premi
bersama** (40 activity bertema premi ada di NB **dan** EDM). Hipotesis "EDM punya mesin premi sendiri"
**gugur**.

Rumus dasar (`CalculatePremiFire` L2003, `HitungPremiOnChange` L645, identik):
`Premium = @Math.divide(.TSI × .Rate × .ProRatePercent × FirstLossScale × IndemnityPct, 1e9, 4)`.
Pembagi `1e9` = pembagi komposit K-018 (1.000 rate‰ × 100 × 100 × 100).

- **PA** `[terverifikasi]`: rumus sama, pembagi `100.000` (= 1.000 × 100; PA tak punya FirstLossScale &
  IndemnityPct) → **penguat K-018, bukan angka hafalan**. EDM menambah jalur **pita tarif** (PCT_RATE_BAWAH ×
  MIN_TSI per pita, basis akumulasi premi). ⚠️ `CalculatePremiPA_FacIn` NB↔EDM diklaim identik fungsional —
  **belum diverifikasi Kiro dengan 23-tag** (metode ringkas beri "beda"); ditandai `[dugaan]`.
- **Travel** `[terverifikasi]`: BENTUK BEDA — `.Rate=0`, tak pakai `Premium=Rate×TSI÷1e9`. Premi dibaca
  langsung: `.PREMIUM` + `.PREMIUMMOREWEEK × minggu_lebih`, sumber `fillActPremiTravel → ViewPremiTravel_SQL`
  yang **NIHIL di korpus** → bentuk tabel tarif `[pertanyaan terbuka]`, masukan dari repository.

### Satuan `ProRatePercent` — **PECAHAN, dipakai langsung tanpa ÷100** (TUTUP)

Pertanyaan Claude: di FIRE ProRatePercent di dalam pembagi komposit; di `SetLocalNonMbuProrate`
pengali langsung — risiko 100×.

**Jawaban terverifikasi (Kiro cek `SetLocalNonMbuProrate` L3525/L4713/L4952):** `.ProRatePercent`
dipakai sebagai **pengali langsung TANPA ÷100**; yang dibagi 1.000 adalah `.Rate` (satuan ‰), bukan
ProRatePercent. **Tidak ada bug 100×** — ProRatePercent memang **pecahan** (mis. 0,9899 untuk 98,99%),
konsisten dengan work owner (tipe Decimal). `PctShortPeriod` justru DIBAGI 100 (L4821 `@Math.divide(.PctShortPeriod,100,4)`)
— dua rasio beda perlakuan, jangan disamakan.

### Dua keluarga prorata — beda satuan & penyebut (jangan disatukan)

`[terverifikasi]`
| Properti | Penyebut | Satuan | Untuk |
| --- | --- | --- | --- |
| `ProratePercentEDMEnd`/`ProratePercentStartEDM` (di-SET di `SetLocalNonMbuProrate`) | 365 tetap | persen | premi dasar |
| `ProrateEDMEnd`/`ProrateStartEDM` (E-5) | periode polis aktual | pecahan | delta spreading |

`[terverifikasi Kiro]` `SetLocalNonMbuProrate` L6209: `.ProRatePercent = .ProratePercentStartEDM +
.ProratePercentEDMEnd` — **varian dua-bagian terkonfirmasi** (Kiro sempat gagal temukan karena regex
salah pola; dikoreksi: ADA di L6209).

### Kontrak Seam 3 — **4 bentuk, satu pintu, tanpa seam baru**

`[terverifikasi]` `SetLocalNonMbuProrate` bukan sekadar pengisi prorata — ia **kalkulator delta premi
tingkat coverage** (L3525 dst.: `(ΔRate×TSI_lama)+(ΔTSI×Rate)×ProRatePercent`), plus cabang
`CalculateMethod=="2"` yang ganti prorata dengan `PctShortPeriod` (23 kemunculan CalculateMethod, 20
PctShortPeriod).

**`premium.Calculate` (Seam 3) memilih di antara 4 bentuk:** (1) rumus dasar · (2) dua-bagian
(Start+End) · (3) tabel tarif Travel · (4) dekomposisi delta coverage. **Tetap 5 seam** — bukan "satu
rumus + varian" melainkan "satu pintu memilih 4 bentuk".

`[pertanyaan terbuka]` **4 bentuk = satu fungsi bercabang atau 4 implementasi di balik Seam 3?** →
keputusan desain, milik work owner. Belum diputuskan.

### Kejanggalan K-046 (diport apa adanya)

`[terverifikasi]` `.PremiRp` pakai rasio prorata tertukar terhadap `.Premium` (konsisten 2 berkas —
mungkin disengaja, korpus tak jelaskan); `FirstLossScale` tanpa guard `@if(…="",100,…)` di `.PremiRp`
padahal ada di `.Premium` (kosong→0, bukan 100). Diport apa adanya, test bernama K-046.

### Belum tuntas (sebelum implementasi, bukan sebelum spec)

- Rumus PA & Travel **penuh** (Travel butuh bentuk tabel tarif dari repository).
- 5 dari 6 salinan EDM mesin premi bersama beda isi (23-tag) — rumus dasar sama, beda di tempat lain.
- `CalculatePremiPA_FacIn` identik fungsional NB↔EDM — verifikasi 23-tag.

### Konsekuensi terhadap dokumen lain

| Dokumen | Perubahan |
| --- | --- |
| `04-spec\09-spec-edm-only.md` | 350 diport; kelompok fungsi; 5 E-Q24 terbuka |
| `07-edm\10-banding-premi-edm-nb.md` + `11-banding-premi-pa-travel.md` | rumus dasar sama; PA/Travel; Seam 3 = 4 bentuk |
| `03-spec-modul-terverifikasi.md` (Seam 3) | `premium.Calculate` = satu pintu, 4 bentuk; konteks endorsement (2 pasang rasio + OldCoverage) |

---

## K-052 · Jalur Life EDM (modul 6): flow `InputEDMLife` sendiri · tanpa tangga akseptasi · **Seam 6 disetujui**

**Tanggal:** 19 September 2026 · **Menutup:** modul 6 (jalur Life), tangga akseptasi Life, usulan Seam 6.
**Modul spec EDM lengkap (6/6).**

### Lingkup Life — berlapis (tiga angka, tiga dasar)

`[terverifikasi]` Selisih angka terjelaskan — beda kriteria:
| Kriteria | Jumlah |
| --- | ---: |
| Nama mengandung "Life" | **68** (16 EDM-only + 52 bersama NB) |
| Nama luas / isi bersinyal | 156 |
| Menyebut predikat `IsLife` | 60 |

Lingkup modul 6 = **L1 inti (16) + L2 subsistem medis (±45)**. L3 (60 titik cabang IsLife) milik modul
1-5 yang punya cabang Life; L4 (skoring bersama) milik modul lain. IsLife = `pyWorkPage.Quotation.BusinessOldId = "L1"…"L16"`.

### Flow `InputEDMLife` — DITAMBAHKAN work owner (bukan dari korpus)

`[terverifikasi work owner]` Life EDM punya **flow tersendiri bernama `InputEDMLife`** yang **baru
ditambahkan work owner** — `[terverifikasi Kiro]` NIHIL di korpus (ekspor lama hanya 2 Flow:
`InputAddendumFacultativeIn` + `OfferFacRetro`). Jadi Life dimodelkan sebagai **alur sendiri**, bukan
cabang di flow endorsement umum.

### Tangga akseptasi Life — TIDAK ADA

**Keputusan work owner: EDM Life TIDAK punya tangga akseptasi.** `[terverifikasi Kiro]`
`GetLimitAkseptasiLife_Act` NIHIL di korpus. Life melewati tangga (arsip #6). Yang menggantikan peran
"boleh/tidak" = **skoring underwriting medis** (§ berikut) yang menghasilkan Standard/Postpone/Decline.

### Skoring: DUA sistem berbeda — jangan disatukan

`[terverifikasi Kiro]`
| Sistem | Berkas | Status | Peran |
| --- | --- | --- | --- |
| Skoring risiko | `ScoringRisk` (8 MB), `SetScore_Act`, `ValueScoringRisk_act` | **BERSAMA NB** (ada di NB & EDM) | skoring risiko umum, bukan baru |
| Pemeriksaan medis Life | `Medical_Harnes` (4,1 MB), `Medical_Sec` (4,1 MB), `SetParamLab_Act`, `CalculateScorLife_Act` | **EDM-only** | underwriting medis Life |

⛔ **`CalculateScorLife_Act` BUKAN perhitungan premi.** `[terverifikasi Kiro]` `.Rate/.TSI/.Premium/@Math.divide
= 0/0/0/0`; `Medical` 6.975×, `Age` 3.145×, keluaran **Standard (423×) / Postpone (132×) / Decline (275×)**
dari ambang klinis (AFP/CEA/PSA/HbA1c/eGFR). Ini **underwriting medis** — keluaran KEPUTUSAN, bukan uang.
Ambang klinis = aturan medis, **diport apa adanya, tidak ditebak** (domain sensitif — nilai ambang &
data pasien TIDAK disalin ke artefak, §6/§3.5).

### Rumus premi Life — keluarga sama, bentuk beda (Seam 3 tetap 4 bentuk)

`[terverifikasi Kiro]` `ReCountPremiLifeEDM`: `TSILiability` (15×), `RateLifeAverage`, pembagi `1000` —
`Premium = (TSILiability × RateLifeAverage × (1 + Rate/100)) ÷ 1.000`. Basis `TSILiability` (bukan TSI),
`.Rate` = faktor pembebanan (bukan rate premi). Pembagi 1.000 konsisten K-018.
✅ **Life = parameterisasi bentuk 1 Seam 3, BUKAN bentuk kelima.** Seam 3 tetap 4 bentuk (K-051).

### Empat titik cabang Life (terverifikasi)

`SetOldData` L1 `IF IsLife → T=6` keluar (lapis B dilewati); `SaveEDMToJsonPolicy_Act` L12/L13 pasangan
komplementer nopolis Life; L28-29 `SaveFacinLive_Act`/`SaveFacinSpreadLife_Sql` **sejajar** (bukan di
dalam `SaveFacinProdAllEDM_Act`); lapis C `SetOLDValueToEDMWork_LIFE` setel 2 penanda (FlagOldData +
IsOldData; cacah 4 karena model data Life satu tingkat list — bukan kelalaian).
`GenerateEDMNoLife` NIHIL (generate tertanam di SaveEDMToJsonPolicy_Act, bukan panic).

### Seam 6 DISETUJUI — `services/underwriting.ScoreMedical`

**Keputusan work owner: Seam 6 disetujui.** `CalculateScorLife_Act` tidak masuk seam mana pun (bukan
predikat/uang/before-image/open-case/akseptasi). Ambang klinis → keluaran Decline/Postpone/Standard
adalah **keputusan underwriting berkonsekuensi nyata** (menolak/menunda pertanggungan); salah ambang
baru ketahuan di produksi. **Total seam kini 6** (3 NB + Seam 4 before-image + Seam 5 open-case +
Seam 6 score-medical). Life pakai seam existing untuk premi (Seam 3)/before-image (Seam 4); Seam 6
khusus skoring medis.

### Konsekuensi terhadap dokumen lain

| Dokumen | Perubahan |
| --- | --- |
| `04-spec\10-spec-edm-life.md` | lingkup berlapis; flow InputEDMLife; tanpa tangga akseptasi; 2 sistem skoring; Seam 6 |
| `03-spec-modul-terverifikasi.md` | seam total 3 → 6 (K-047 Seam4, K-049 Seam5, K-052 Seam6) |
| Diagram/temuan "EDM 2 Flow" | ditambah catatan: work owner menambah flow ke-3 `InputEDMLife` (di luar ekspor lama) |

---

## K-053 — Fac Out (retrosesi keluar) menjadi menu terpisah namun terhubung

**Pertanyaan asal (work owner):** saat UW Accept, lokasi dengan spreading Fac Out disalin ke Fac Retro
lalu masuk menu Fac Out; bisakah Fac Out dibuat menu terpisah dari NB/RNW/EDM tetapi tetap terhubung.

**Keputusan:** ya. Modul `internal/services/facout/` terpisah dari `faccase`, dependensi searah
(`facout` membaca hasil akseptasi Fac In, tidak sebaliknya). Keterhubungan dijaga lewat penaut data
`OBJECTNO_FACIN` + `IDPEGA`/`POLICYNO`/`NOENDORS` (kolom `FACOUTPRODUCTION`, terverifikasi).

**Dasar [terverifikasi]:** `SetValidateDateUW_PostAct` step 3 `Call SetDataFacOut_Act`; penyalinan ke
`OfferFacIn.FacRetro`/`FacRetroList`; flow `InputInwardFacultativeOffer` bercabang ke `OfferFacRetro`
(shape `FAC OUT?`/`INPUT FAC OUT?`). Detail di `09-facout\01-temuan-dan-rancangan-facout.md`.

## K-054 — Penanda spreading Fac Out diport apa adanya (dua konteks berbeda)

**Keputusan:** predikat penanda Fac Out di-port persis dari korpus, TIDAK disederhanakan:
- Salin khusus ke FacRetro: `.TreatyType == "10015"` (kesetaraan ketat) — di `SetDataFacOutFire_Act`,
  `CopyFacRetroFire_ACT`.
- Salin spreading umum: `@contains(.TreatyType,"10015") || @contains(.TreatyName,"SPL") ||
  @contains(.TreatyType,"10007")` — di `CopyToAllSpreading_ACT`, `CopyToAllLocSpreading_ACT`.

Dua konteks **tidak digabung**. Predikat masuk registry (K-050).

**Arti kode (keterangan work owner):** `10015` = **FACOUT**, `10007` = **ORS** (Own Retention Share).
Ini keputusan work owner — belum ada file korpus yang memetakan angka→label eksplisit (`M_TREATY_IN` =
`ID`+`JSONDATA CLOB`; label ada di dalam JSON). TreatyName `"SPL"` = [pertanyaan terbuka].
**[terverifikasi korpus]** `CoverageBasis==5` = **Layering Basis** (`DDL\CoverageBasis.xml`
`pyPromptTableList`: 1=Sum Insured, 2=First Loss, 3=EML/PML, 4=Sub Limit, 5=Layering).

## K-055 — Tangga persetujuan Fac Out = tangga retro tersendiri (4 tingkat), bukan Seam 2

**Keputusan:** Fac Out memakai tangga peran retrosesi sendiri, bukan mesin akseptasi limit-berjenjang
Fac In (Seam 2). Urutan target (work owner): **ReasFacOutAdmin → ReasFacOutHead → ReasFacOutGroupLeader
→ ReasFacOutTechnicalDirector**, dengan percabangan `IsGroup`, diakhiri cetak RI Slip + insert produksi.

**Dasar [terverifikasi]:** ekspor ulang `DDL\OfferFacRetro.xml` (dari work owner, `ASM-FW-GISFW-Work`)
memuat keempat peran `ReasFacOutAdmin`/`Head`/`GroupLeader`/`TechnicalDirector` + `Is it group?` +
`Accept?` + `PRINT R/I SLIP` + `InsertFacoutProd` + routing `WorkBasket`. Ekspor lama NB/RNW 2 tingkat
= versi lama, digantikan. Nomor slip retro: `DDL\GENERATE_FACRETRO_NO.txt` → `RNM-Y...` (Fac Code 'Y'
utk 'FAKULTATIF RETROSESI'), sequence `FACRETRO_SEQ`. Urutan transisi persis = UI work owner mengalahkan
XML. Seam Fac Out terpisah dari 6 seam Fac In.

## K-056 — Produksi Fac Out ke tabel `FACOUTPRODUCTION` yang sudah ada (flat)

**Keputusan:** produksi Fac Out menulis ke tabel **`FACOUTPRODUCTION`** existing (65 kolom, flat,
multi-COB, dengan pasangan `_MENJADI`/`_SELISIH`) — TIDAK membuat tabel flat baru untuk produksi Fac Out.
Daftar Fac Out dibaca via `GetFacoutList_SQL`: `select IDPEGA, policyno from facoutproduction where
policyno = {InputData.CARI17}` (terverifikasi, tag `pyBrowseSQL`).

**[terverifikasi] Insert = Connect-SQL langsung, BUKAN stored procedure.** `DDL\InsertTreatyProd_Sql.xml`
(`Rule-Connect-SQL`, `pyBrowseSQL`) = `BEGIN INSERT INTO FACOUTPRODUCTION (65 kolom) VALUES
({DataIN.CARIn}, {DataINCargo.*}, {DataINMBU.*}, {DataINPA.*}); COMMIT; END;`. SQL & pemetaan
kolom→parameter lengkap di korpus — **tidak butuh ALL_SOURCE**. (Koreksi: catatan "SP
INSERTFACOUTPRODUCTION" sebelumnya keliru; itu nama request/connector, bukan SP DB.)

**Yang dimodelkan baru:** hanya staging JSON `FacRetro`/`FacRetroList` + `FacOutObjectList{Rate,
FacOutTSI}` bila `JSON_POLIS` diganti flat (tabel `FLAT_FACRETRO_*`, §08-flat).

---

## K-057 · Pemicu Fac Out ditetapkan = flag `.IsFacRetro`

**Tanggal:** 21 September 2026 · **Mencabut:** rumusan pemicu lama di
`09-facout\01-temuan-dan-rancangan-facout.md` §0, §3, §4.2 · **Sumber bahan:**
`10-audit\05-draf-keputusan-facout.md`

**Pertanyaan asal:** apa yang sesungguhnya memicu kasus masuk ke menu Fac Out.

**Keputusan work owner:** pemicu Fac Out adalah **flag `.IsFacRetro`**, yang di-*reset* ke `0` lalu
dinyalakan ke `1` **per lini bisnis** oleh penanda spreading Fac Out saat UW Accept.
⛔ Rumusan lama **"UW Accept (`ProposalAcceptStatus = 4` → When `IsFacout`)" DICABUT.**

### Dasar `[terverifikasi]` — rantai lima langkah

| # | Berkas | Tag pembawa | Isi |
| ---: | --- | --- | --- |
| 1 | `Activity\SetValidateDateUW_PostAct` | `<pyStepPageReference>` = `RH_1.pySteps(3)` | `Call SetDataFacOut_Act`, **tanpa precondition** |
| 2 | `Activity\SetDataFacOut_Act` | `<Property>` ×4 pada `RH_1.pySteps(2)` · `<PropertiesName>` pada `RH_1.pySteps(3)` | `RH_1.pySteps(2)` `Property-Remove` atas `OfferFacIn.FacRetro.`{`LocationList`,`PersonList`,`CargoList`,`VehicleList`}; `RH_1.pySteps(3)` `Property-Set` `.OfferFacIn.IsFacRetro = 0`; `RH_1.pySteps(4..7)` panggil per COB |
| 3 | `Activity\SetDataFacOut{Fire,AnekaGolf,CargoMBU,PATravel}_Act` | `<pyStepsPreCondParamsWhen>` · `<PropertiesName>` | **keempatnya** menyetel `pyWorkPage.OfferFacIn.IsFacRetro = 1`, bergerbang `.TreatyType=="10015"` / `Local.isFacOut==1` / `Local.isFacOutLoc==1` |
| 4 | flow masuk siklus ybs. | `<pyTaskStatusOrWhen>`=`WHEN` · `<pyTaskWhen>`=`IsFacRetro` · `<pyLikelihood>`=100 | transisi menuju shape `<pyMOName>` = `[Fac Out]`; di NB `InputInwardFacultativeOffer` transisi itu ber-`<pyID>` = `Transition92` |
| 5 | shape `[Fac Out]` | **`<pyImplementation>`** = `OfferFacRetro` | **bukan** `<pySubFlowName>` (0 kemunculan) |

⚠️ **`OfferFacIn.FacRetroList` BUKAN sasaran `Property-Remove` langkah 2.** Ia muncul di
`RH_1.pySteps(8).pySteps(1)` sebagai `<pyStepsObjectName>` pada langkah **bermetode kosong**.
Kekeliruan itu pernah tercatat dan sudah dicabut — lihat `10-audit\06-tinjau-ulang-suntingan-facout.md` §2.4.

### Dua lapis kekeliruan rumusan lama

1. `[terverifikasi]` `IsFacout` muncul **0 kali** di `NB FacIn\Flow\InputInwardFacultativeOffer.xml`,
   peka huruf maupun tidak (`IsFacRetro` **17**, `IsInputFacRetro` **9**). Sensus peka huruf seluruh
   korpus: hanya **3 berkas per folder**, seluruhnya kondisi tampilan Section + berkas definisinya.
2. `[terverifikasi]` `ProposalAcceptStatus = 4` berarti **Banding**, bukan Accept (**K-029**,
   `GLOSARIUM.md`).

### Konsekuensi

- `IsFacout` **keluar** dari registry predikat Fac Out; porting-nya tetap di tiket New Business
  `05-tickets\09-predikat-sikap-khusus.md` di bawah **K-019**, yang **tidak berubah**.
- Tiket **F01** dan **F02** memakai pemicu ini.

**`[pertanyaan terbuka]` yang TIDAK ditutup keputusan ini:** indeks rujukan flow menuliskan resolusi
`IsFacRetro` ke app **SFAGIS** / workType **ClaimLife** dan `OfferFacRetro` ke app **GISFW** /
workType **LIFE**, sementara `Embed-Reference-Rule` menyebut kelas `ASM-FW-GISFW-Work` dan definisi
ekspornya berkelas `ASM-FW-GISFW-Data-OfferFacIn`. **Butuh UI Pega.**

---

## K-058 · Sumber kebenaran = salinan per siklus, **kecuali** rule yang sudah diekspor ulang ke `DDL\`

**Tanggal:** 21 September 2026

**Pertanyaan asal:** bila satu identitas rule hadir di lebih dari satu folder dengan nomor versi sama
tetapi isi berbeda, salinan mana yang menjadi sumber kebenaran spec.

**Keputusan work owner:** **opsi A — salinan masing-masing siklus**, **DENGAN PENGECUALIAN**: rule
yang sudah **diekspor ulang ke `DDL\`** — salinan `DDL\` yang **menang**.

**Tentang `Flow\OfferFacRetro`:** perbedaan antar salinan **hanya pada alurnya** (2 tingkat vs 4
tingkat). **Tampilan menu dan isinya SAMA untuk keempat peran** — `ReasFacOutGroupLeader` dan
`ReasFacOutTechnicalDirector` memakai **layar yang sama** dengan `ReasFacOutAdmin` dan
`ReasFacOutHead`.

### Dasar `[terverifikasi]` — sensus drift 21 September 2026 (kontrak 23 tag)

| Pasangan | Dibandingkan | Versi sama + isi sama | **Versi sama + isi BEDA** | Versi beda |
| --- | ---: | ---: | ---: | ---: |
| NB ↔ RNW | 1.907 | 1.907 | **0** | 0 |
| NB ↔ EDM | 1.707 | 1.353 | **267** | 87 |
| RNW ↔ EDM | 1.674 | 1.330 | **259** | 85 |
| **Total** | 5.288 | — | **526** | **172** |

Pemeriksaan silang `1.353 + 267 + 87 = 1.707` ✅ dan `267 + 87 = 354` ✅ — cocok K-043.

⛔ **Nomor versi tidak dapat dipakai menilai kebaruan.** `[terverifikasi]`
`DDL\CountRateRetroCov.xml` ber-`<pyRuleSetVersion>` **01-01-01** dengan `<pxCommitDateTime>`
**20260921**, sedangkan salinan korpus ber-versi **01-01-95** dengan commit **20260728** — **nomor
lebih rendah pada commit lebih baru**.

### Konsekuensi

- Rule yang ada di `DDL\` dipakai sebagai sumber kebenaran, mengalahkan salinan korpus.
- Daftar prioritas ekspor ulang berikutnya: `10-audit\08-daftar-ekspor-ulang-prioritas.md`.

### ⚠️ `[pertanyaan terbuka]` — TIDAK diputuskan di sini

**Apakah metode deteksi drift diganti** dari banding `pyRuleSetVersion` menjadi banding **hash 23
tag**. `[terverifikasi]` metode lama melewatkan **267 dari 354** perbedaan NB↔EDM (75,4 %), dan angka
"95 rule berbeda versi" di `PANDUAN-KERJA` §5 / `CLAUDE.md` §4.5 **tidak tereproduksi** (hasil ukur:
87 / 85 / 0). **Tetap terbuka — menunggu keputusan work owner.**

---

## K-059 · Populasi kerja Fac Out = nama Fac Out **DAN** pembawa kode `10015`

**Tanggal:** 21 September 2026

**Pertanyaan asal:** berapa luas permukaan Fac Out yang harus diport.

**Keputusan work owner:** populasi kerja = berkas **bernama `FacOut`/`FacRetro`** **DAN** berkas
**pembawa kode `10015`**, **diperlakukan sama** — karena **`10015` adalah kode Fac Out**. Sekitar
**55 activity per folder**.

### Dasar `[terverifikasi]`

| Ukuran | NB | RNW | EDM |
| --- | ---: | ---: | ---: |
| Activity bernama `FacOut`/`FacRetro` (Ordinal, peka huruf) | 41 | 40 | 39 |
| Hadir di ketiga folder | — | **39** | — |
| Activity pembawa `10015` **tanpa** nama Fac Out | 16 | 15 | 16 |
| **Permukaan fungsional** | ≈57 | **≈55** | **≈55** |

Banding 23 tag atas ke-39: **NB vs RNW 39/39 identik**; **NB vs EDM 30/39 identik, 9 berbeda**.

⛔ Angka lama **18/18 dan 10/18 DIBATALKAN** — populasinya tidak pernah terdefinisi.

### Konsekuensi

Tiket **F11** memakai populasi ≈55. Berkas seperti `CopyToAllSpreading_ACT`,
`CopyToAllLocSpreading_ACT`, `cekSpreadingFactIn`, `CountPremiumNet`, `SaveJsonPolicyFacIn_Act`,
`SetReinsurerEndorsement_Act` **wajib ikut ditinjau** meski namanya tidak ber-Fac Out.

---

## K-060 · Premi retro memakai `CountRateRetroCov` — **DUA rule, dipilih menurut COB**

**Tanggal:** 21 September 2026 · **Diamandemen:** 21 September 2026 (bukan nomor K baru)

**Pertanyaan asal:** premi retro sebagai bentuk Seam 3 atau seam tersendiri, dan salinan mana yang
berlaku.

> ## ⛔ AMANDEMEN — rumusan semula TIDAK LENGKAP
>
> Rumusan semula:
>
> ~~"Perhitungan premi retro tetap memakai `CountRateRetroCov`, dan salinan yang berlaku adalah
> ekspor ulang **`DDL\CountRateRetroCov.xml`**."~~
>
> **Tidak salah, tetapi hanya menyebut SATU dari DUA rule.** `CountRateRetroCov` adalah **dua rule
> berbeda di kelas berbeda**, bukan satu rule dengan dua salinan.
>
> **Rumusan yang berlaku:** perhitungan premi retro memakai **`CountRateRetroCov`**, yang berwujud
> **dua rule terpisah yang dipilih menurut lini bisnis (COB)**. **Keduanya berlaku**, keduanya sudah
> diekspor ulang ke `DDL\`, dan **keduanya harus diport** — bukan dipilih salah satu.

### Dua rule, dibedakan `pzInsKey` dan `pyClassName` — bukan nama berkas

`[terverifikasi]`

| | Varian **FIRE** | Varian **ANEKA** |
| --- | --- | --- |
| `pyClassName` | `ASM-FW-GISFW-Data-PropertyItem` | `ASM-FW-GISFW-Data-Aneka` |
| Basis `pzInsKey` | `RULE-OBJ-ACTIVITY ASM-FW-GISFW-DATA-PROPERTYITEM COUNTRATERETROCOV` | `RULE-OBJ-ACTIVITY ASM-FW-GISFW-DATA-ANEKA COUNTRATERETROCOV` |
| **Berkas berlaku** | **`DDL\CountRateRetroCov.xml`** | **`DDL\CountRateRetroCov(ANEKA).xml`** |
| ver / commit berkas berlaku | `01-01-01` / `20260921T095256` | `01-01-95` / `20260831T101051` |
| Ukuran | 131.407 B | 115.823 B |
| Salinan korpus | **NB** dan **RNW** (`01-01-95`, commit `20260728T042513`) | **Endorsment** (commit `20260728T065753`) |
| Kemutakhiran salinan korpus | ⛔ **BASI** — dua perbedaan, lihat di bawah | ✅ **MUTAKHIR** — `[terverifikasi]` **nol beda** terhadap ekspor ulang (metode, gerbang, T/F, penugasan) |

`[terverifikasi]` `pzOriginalInstanceKey` **keduanya sama**
(`…DATA-PROPERTYITEM… #20240627T074538.631 GMT`): varian Aneka **lahir dari menyalin** varian
PropertyItem, lalu menyimpang.

⚠️ **`pzInsKey` memuat stempel waktu per-penyimpanan**, sehingga identitas rule dibaca dari
**basisnya** (`RULE-OBJ-ACTIVITY <KELAS> <NAMA>`), bukan dari kunci utuh.

### Perbandingan kedua varian — `[terverifikasi]` dari ekspor ulang `DDL\`

| Aspek | **PropertyItem (FIRE)** | **Aneka (ANEKA)** |
| --- | --- | --- |
| Langkah | **14** (10 puncak + 4 bersarang) | **12** (9 puncak + 3 bersarang) |
| **Pembagi premi** | **100.000** | ⛔ **10.000** |
| **Uji mata uang** | `.Currency=="IDR"` · `.Currency=="USD"` | ⛔ `.Currency.Name=="IDR"` · `.Currency.Name=="USD"` |
| `Local.RIComIN` **diakumulasi** (langkah 8) | ✅ **ADA** | ⛔ **TIDAK ADA** |
| `Local.RIComIN` **dibaca** (`RICommPercentage`) | ada | **ada** |
| `Local.TotalPremiCov` | ada, diakumulasi (2 penugasan) | ⛔ **tidak ada** (0) |
| Rumus koreksi langkah 10 / 10.1 | ✅ ada | ⛔ **tidak ada** |
| Penugasan `.PremiumRetro` | 3 | 2 |
| Gerbang `pySteps(2)` / `(3)` / `(4)` | `IsEDM` / `IsEdmExtendPeriod` / `IsEDM` | **identik** |

### ⛔ Konsekuensi yang mengubah temuan sebelumnya

**Keanehan `Local.RIComIN` MASIH ADA — tetapi hanya di varian ANEKA.** `[terverifikasi]` Varian Aneka
**membaca** `Local.RIComIN` di `.RICommPercentage = @Math.divide(Local.RIComIN, Local.LengCov, 20)`
tetapi **tidak pernah mengakumulasinya**. Jadi:

- Varian **FIRE**: keanehan **tidak ada** (diakumulasi di langkah 8) — pencabutan sebelumnya **benar
  untuk varian ini**.
- Varian **ANEKA**: keanehan **ada**, **diport apa adanya**, **kandidat perbaikan** (K-046).

⚠️ Pencabutan menyeluruh yang tercatat sebelumnya **terlalu jauh** dan **dibatasi ulang** ke varian
FIRE saja.

### Dasar `[terverifikasi]` — varian FIRE, salinan berlaku vs salinan korpus

| Salinan | Ukuran | `pyRuleSetVersion` | `pxCommitDateTime` |
| --- | ---: | --- | --- |
| `DDL\CountRateRetroCov.xml` | 131.407 B | `01-01-01` | `20260921T095256 GMT` |
| `NB FacIn\Activity\CountRateRetroCov.xml` | 131.835 B | `01-01-95` | `20260728T042513 GMT` |

Keduanya **14 langkah** (10 puncak + 4 bersarang, **2 tingkat**), alamat identik kecuali awalan
halaman harness (`RH_1` korpus vs `RH_2` DDL — artefak ekspor). **Metode identik di seluruh 14
langkah** (0 perbedaan).

### ⛔ DUA perbedaan isi, bukan satu

`[terverifikasi]` diukur dengan parser yang memisahkan **anak langsung** dari keturunan:

| Alamat | Kelas perbedaan | Korpus | DDL |
| --- | --- | --- | --- |
| `pySteps(8)` | **penugasan** | `Local.RateIN += .Rate` · `Local.LengCov := .pxListSubscript` | **ditambah** `Local.RIComIN := Local.RIComIN + .RIComm` |
| `pySteps(3)` | **gerbang + T/F** | **dua** baris precondition: `IsEDM` (T=2/F=3) **dan** `IsEdmExtendPeriod` (T=2/F=3) | **satu** baris: `IsEdmExtendPeriod` (T=2/F=3) — **gerbang `IsEDM` HILANG** |

Jumlah baris precondition seluruh pohon: korpus **15**, DDL **14**; gerbang `IsEDM` persis: korpus
**3**, DDL **2**.

⚠️ `<pyStepsDescription>` langkah 3 **tetap berbunyi** `IsEDM IsEdmExtendPeriod` di **kedua** salinan
— label itu **basi** di salinan DDL. Contoh lagi bahwa **label bukan bukti** (`PANDUAN-KERJA` §3).

### Konsekuensi

1. **Keanehan `Local.RIComIN` pada varian FIRE tidak ada** — akibat salinan korpus yang basi. ⚠️ Pada
   **varian ANEKA keanehan itu ADA** dan tetap kandidat perbaikan (lihat amandemen di atas).
2. Langkah 3 memaksa `Local.prorate := 100` pada **kondisi yang lebih luas** di salinan DDL
   (tanpa syarat `IsEDM`). Ini **perbedaan perilaku pada jalur uang** dan diport mengikuti DDL.
3. `10-audit\01-pohon-langkah-countrateretrocov.md` memuat pohon **kedua** varian.
4. ⛔ **Satu implementasi tidak cukup.** Jalur properti mata uang berbeda (`.Currency` vs
   `.Currency.Name`); implementasi yang membaca mata uang lewat satu jalur akan **salah pada salah
   satu varian**.
5. **Pemilihnya COB** — keputusan work owner. Korpus **tidak** memuat gerbang yang memilih di antara
   kedua rule; pemilihan terjadi lewat **kelas halaman** tempat rule dipanggil.

### `[pertanyaan terbuka]` yang TIDAK ditutup keputusan ini

1. **Kemungkinan varian ketiga berkelas `ASM-FW-GISFW-Data-Cargo`.** `[terverifikasi]` Indeks
   `Embed-Reference-Rule` pada **7 berkas** mencatat `CountRateRetroCov` pada **tiga** kelas —
   `Data-PropertyItem`, `Data-Aneka`, **`Data-Cargo`** — dan tiap pemanggil produksi memuat **tiga**
   langkah `Call CountRateRetroCov`. ⛔ Varian berkelas `Data-Cargo` **tidak ada** di ketiga folder
   korpus maupun di `DDL\`. Work owner menyatakan hanya ada **dua** rule di Pega. Indeks rujukan
   merekam resolusi **saat berkas terakhir disimpan**, jadi ia **bukan bukti** rule itu masih ada
   sekarang. **Jangan disimpulkan salah satu.** Rincian: `10-audit\09-pemanggil-countrateretrocov.md`.
2. Tiga penugasan `Local.prorate` berurutan di langkah 1, yang terakhir `:= 1`, **tidak berubah** di
   ekspor baru. Apakah dua rumus prorata sebelumnya mati bergantung pada apakah urutan
   `REPEATINGINDEX` = urutan eksekusi — **`[di luar korpus]`, butuh UI Pega**.
3. Mengapa gerbang `IsEDM` hilang di `pySteps(3)` varian FIRE — disengaja atau kekeliruan ekspor.

📌 **Temuan "rule kembar" ini digeneralisasi ke seluruh korpus** di
`10-audit\10-rule-kembar-pzinskey.md`: **215 pasangan** ber-kunci/kelas berbeda atas **109 nama
unik**. Sebagian angka drift yang tercatat sebelumnya ternyata rule kembar, bukan drift.

### 📎 Catatan pelengkap — 21 September 2026 · `Local.prorate`

> **Bunyi keputusan K-060 di atas TIDAK berubah.** Ini catatan pelengkap, bukan amandemen dan bukan
> nomor K baru.

Work owner menetapkan bahwa pada `CountRateRetroCov`, **`Local.prorate` bernilai `1` keluar dari
langkah 1** — dua penugasan sebelumnya di langkah yang sama **tidak berlaku**. `[terverifikasi]`
keempat langkah prorata **identik di kedua varian** (`pySteps(1)` tanpa gerbang T=2/F=2; `pySteps(2)`
`IsEDM` T=2/F=3; `pySteps(3)` `IsEdmExtendPeriod` T=2/F=3; `pySteps(4)` `IsEDM` **T=3/F=2**, terbalik),
sehingga **nilai akhir `Local.prorate` per jalur** adalah: **NB/RNW (non-EDM) → `100`**; **EDM biasa →
`(ProrateEDMEnd + ProrateStartEDM) × 100`**; **EDM extend period → `100`**. Konsekuensinya, dua rumus
pertama di langkah 1 — selisih periode dan pembagian terhadap `EDMDay`/365 — adalah **KODE MATI di
semua jalur**; keduanya **tetap diport apa adanya** (`CLAUDE.md` §1) dan didaftarkan sebagai kandidat
perbaikan **K-046** (`K046_Prorate_DuaRumusMati_Langkah1`), **tidak dihapus**. Pemeriksaan silang
**K-018** membenarkan kedua pembagi: pada jalur non-EDM, FIRE menghasilkan `Share × Rate × 100 /
100000` = `Share × Rate / 1000` (**‰**) dan ANEKA `Share × Rate × 100 / 10000` = `Share × Rate / 100`
(**%**) — keduanya persis skala rasio per lini bisnis yang dikunci K-018, sehingga perbedaan pembagi
**bukan kekeliruan** melainkan konsekuensi aturan pembagi komposit. ⚠️ Yang dijawab work owner adalah
**nilai `Local.prorate`**, **bukan** aturan umum bahwa urutan `REPEATINGINDEX` = urutan eksekusi;
butir 2 pada daftar `[pertanyaan terbuka]` di atas **tetap terbuka sebagai aturan umum**. Rincian:
`10-audit\01-pohon-langkah-countrateretrocov.md` §7D.

---

## K-061 · F13 (cetak RI Slip + email) ditulis penuh, ditandai menunggu

**Tanggal:** 21 September 2026

**Pertanyaan asal:** apakah cetak RI Slip dan kirim email masuk lingkup sekarang atau ditunda.

**Keputusan work owner:** **ditulis penuh, ditandai menunggu, mengikuti pola R07**
(`05-tickets\rnw\R07-konversi-produksi-renewal.md`). **Tidak dihapus.**

**Dasar `[terverifikasi]`:** `DDL\OfferFacRetro.xml` `<pyMOName>` memuat `PRINT R/I SLIP`,
`Send Email Print RI Slip`, `IsRISlip`, `PrintRISlip_FlowAction`; empat activity PDF ada di korpus.
⛔ Penghambatnya `CLAUDE.md` §4.4 — daftar endpoint ada di tabel Oracle `M_LINK_SERVICE`, **isinya
tidak ada di korpus**.

**Konsekuensi:** lingkupnya tetap terlihat; hanya eksekusinya menunggu. Menghapusnya akan
menyembunyikan pekerjaan yang pasti ada.

**`[pertanyaan terbuka]`:** arti enumerasi tipe email `1` / `2` / `7` pada jalur retro — tidak
dijelaskan korpus.

---

## K-062 · `GetOPFacOut_Sql` **TIDAK diport dan DIHAPUS** — ⛔ pengecualian pertama terhadap K-006

**Tanggal:** 21 September 2026

**Pertanyaan asal:** `RDBList\GetOPFacOut_Sql` membaca tabel internal Pega yang tidak ikut
bermigrasi. Apa penggantinya.

**Keputusan work owner:** klep itu **TIDAK diport**; **dihapus, tidak dipakai lagi**. **Kolom
pelaksana TETAP ADA** di layar Fac Out, tetapi **ISINYA KOSONG**.

### ⛔ CATAT TEGAS — pengecualian yang DISENGAJA

> **Ini pengecualian PERTAMA yang disengaja terhadap pola "ditangguhkan, bukan dihapus"**
> (**K-006**, `PANDUAN-KERJA` §5, `CLAUDE.md` §4.5).
>
> Pola baku: rule yang dirujuk tetapi tidak dapat dijalankan **ditangguhkan** — cabangnya
> **`panic` bila tercapai**, **tidak dihapus**. **K-062 sengaja menyimpang dari pola itu.**
>
> ⚠️ **Sesi berikutnya JANGAN membaca ini sebagai inkonsistensi lalu membalikkannya.** Ia adalah
> keputusan work owner yang sadar dan tercatat. Membalikkannya memerlukan keputusan baru, bukan
> "perapian".

**Alasan pengecualiannya dapat dipertanggungjawabkan:** yang hilang adalah **tampilan identitas
pelaksana**, bukan angka uang atau gerbang alur. Mengosongkan kolom **tidak** menghasilkan selisih
rekonsiliasi (ADR-0001), sehingga `panic` tidak diperlukan.

### Dasar `[terverifikasi]`

`RDBList\GetOPFacOut_Sql` (`<pxObjClass>` = `Rule-Connect-SQL`, tag `<pyBrowseSQL>`) membaca
`datapega.pc_history_asm_fw_gisfw_work`. `CLAUDE.md` §4.3: `DATAPEGA.PC_*` **hilang bersama Pega —
jangan dimigrasikan apa adanya**.

⚠️ Nilai yang dibacanya adalah **identitas orang** — K-025: mekanismenya saja yang dicatat, nilainya
tidak pernah disalin.

### Konsekuensi

- Sisi pemanggil (`Activity\OfferFacOut_PreAct`) **wajib menghasilkan nilai kosong, bukan galat**.
- Tiket **F12** tidak lagi memuat blocker eksternal untuk klep ini.
- **DRAF K-062 lama** (yang mengusulkan membangun ulang di atas riwayat sendiri) **digantikan** oleh
  keputusan ini.

---

## K-063 · Mata uang pada tabel flat — **empat perlakuan**, dan lima total akar **tanpa mata uang**

**Tanggal:** 24 September 2026 · **Menutup:** butir **9** rekonsiliasi tabel flat ·
**ADR:** `adr/0007-agregat-lintas-mata-uang-bukan-money.md`

**Pertanyaan asal:** apakah tiap kolom uang di rancangan tabel flat diberi kolom mata uang
pendamping, sebagaimana dituntut **K-010**, **K-012**, **ADR-0004** dan **ADR-0006**.

**Keputusan work owner:** ya — **kecuali** lima total skalar di akar penawaran, yang **sengaja tidak
diberi** kolom mata uang. Perlakuannya dibagi **empat kelompok**.

⛔ **Pembawa mata uang ditentukan dari ISI korpus, bukan dari nama kolom.** `[terverifikasi]` di
korpus mata uang dibawa medan bernama **`Currency`**, **`Name`**, dan **`TreatyName`** — penyaringan
berbasis nama melewatkan hampir semuanya. `TreatyName` namanya *treaty*, isinya **kode mata uang**.

### Cakupan sesudah koreksi 24 September

`[terverifikasi]` diukur atas **115** contoh `DDL\CONTOH\*.xml` (populasi efektif **113**):
**33 tabel · 147 kolom uang**, turun dari 35 / 156 setelah **sembilan kolom** terbukti bukan uang.

⛔ **Seluruh hitungan kolom di keputusan ini DIHITUNG TANPA CABANG `OldData`.** Dasarnya **V-19**:
`OldData` **tidak disimpan**, jadi ia tidak boleh masuk inventaris kolom. `[terverifikasi]` **326
kemunculan dibuang** oleh aturan itu pada sapuan kesembilan medan.

📌 Jumlah tabel rancangan **75 → 78** (24 September): tiga tabel wadah yang dipakai sebagai induk di
**27 baris** `Daftar Relasi` tetapi belum pernah didefinisikan — `T_FR_FACOFFERLIST`,
`T_FR_OBJECT`, `T_FR_POLICY`. Ketiganya **wadah murni** (kolom sistem saja), sesuai **V-27**.

| Kelompok | Tabel | Kolom | Tindakan |
| --- | ---: | ---: | --- |
| **(a) AMAN** | **8** | **25** | Sudah punya kolom mata uang yang berlaku. **Biarkan.** |
| **(b) PULIHKAN** | **10** | **61** | Mata uangnya **ADA di korpus** pada jalur tabel itu sendiri, tetapi rancangan tidak menyimpannya. **Kembalikan kolomnya.** |
| **(c) WARISI** | **8** | **33** | Tidak punya sendiri; **leluhurnya punya**. Diisi saat tulis. |
| **(d) TANPA MATA UANG** | — | **5** | Lima total skalar di akar penawaran. **Tidak diberi kolom mata uang.** |

### (a) AMAN — 8 tabel / 25 kolom

`T_CURRENCYLIST` (12) · `T_PROPERTYITEMLIST` (3) · `T_FR_PROPERTYITEMLIST` (3) ·
`T_LOCATIONLIST` (2) · `T_DEDUCTIBLELIST` (2) · `T_FR_DEDUCTIBLELIST` (1) ·
`T_FR_LISTCAUSEOFLOSS` (1) · `T_LISTCAUSEOFLOSS` (1)

⛔ **Dua kolom TIDAK dihitung sebagai pembawa mata uang yang berlaku:**

| Kolom | Isinya | Sebab tidak dihitung |
| --- | --- | --- |
| `FLAG_CURRENCY` | **ANGKA**, 1.709 kemunculan | **penanda**, bukan kode mata uang |
| `CURRENCY_OLD_ID` | **ANGKA**, 844 kemunculan | rujukan mata uang **LAMA** |

### (b) PULIHKAN — 10 tabel / 61 kolom · ✅ **DITERAPKAN 24 September 2026**

`T_FR_CURRENCYLIST` (18) · `T_COVERAGELIST` (17) · `T_FR_COVERAGELIST` (10) ·
`T_FR_NETPERCURRENCY` (4) · `T_FR_OFFEREDPAYMENT` (3) · `T_FR_ANEKALIST` (3) ·
`T_CEDING_CURRENCYLIST` (2) · `T_FR_TOTALTSIPREMIRETRO` (2) · `T_ANEKALIST` (1) · `T_PERSONLIST` (1)

⛔ **Membuangnya berarti membuang informasi yang sudah kita punya.** Mata uangnya terbaca di korpus
pada jalur tabel itu sendiri — rancangan hanya tidak menyimpannya. Mengembalikannya tidak menuntut
data baru dari siapa pun.

#### ✅ Cara penerapannya — **satu `CURRENCY_CODE` per tabel, bukan 61 `*_CCY`**

`[terverifikasi]` Pembawa mata uang **seragam di kesepuluh tabel**: medan **`<Name>`** di dalam
**elemen berkunci kode mata uang**, tepat di simpul tabel itu sendiri. Diukur atas 115 contoh,
cabang `OldData` dikecualikan:

| Tabel | Kemunculan | Berkas | Kode terbaca |
| --- | ---: | ---: | --- |
| `T_COVERAGELIST` | 231 | 90 | EUR · IDR · USD |
| `T_ANEKALIST` | 127 | 81 | EUR · GBP · IDR · USD |
| `T_CEDING_CURRENCYLIST` | 103 | 101 | EUR · IDR · USD |
| `T_FR_CURRENCYLIST` | 34 | 17 | IDR · USD |
| `T_FR_COVERAGELIST` | 20 | 12 | IDR · USD |
| `T_FR_NETPERCURRENCY` | 17 | 17 | IDR · USD |
| `T_FR_ANEKALIST` | 15 | 9 | IDR · USD |
| `T_PERSONLIST` | 5 | 5 | IDR |
| `T_FR_OFFEREDPAYMENT` | 2 | 2 | IDR · USD |
| `T_FR_TOTALTSIPREMIRETRO` | 1 | 1 | USD |

Karena korpus membawa **satu mata uang per baris**, yang ditambahkan adalah **satu kolom
`CURRENCY_CODE VARCHAR2(10)` per tabel** — **10 kolom**, menaungi ke-61 kolom uang itu.

⛔ **Alternatif yang TIDAK dipilih:** 61 kolom `*_CCY`, satu per kolom uang, seperti rancangan lama
`08-flat\01`. Itu **denormalisasi** dan **tidak didukung bentuk korpus**. Bila tetap diinginkan,
perlu keputusan tersendiri.

⚠️ **Jebakan yang hampir menyesatkan:** `<Currency>` muncul **1.936 kali** di bawah `CoverageList`,
tetapi jalurnya `CoverageList/**DeductibleList**/Currency` — itu milik **`T_DEDUCTIBLELIST`** yang
sudah ada di kelompok (a), **bukan** pembawa `CoverageList`.

#### ✅ TERTUTUP 24 September — satu baris **tidak pernah** membawa lebih dari satu mata uang

~~`[pertanyaan terbuka]` Apakah satu baris dapat membawa lebih dari satu mata uang, sehingga perlu
pecah baris, bukan satu kolom.~~

`[terverifikasi]` Sapuan seluruh baris `<rowdata>` di 115 contoh, cabang `OldData` dikecualikan:

| Ukuran | Nilai |
| --- | ---: |
| Baris diperiksa | **13.666** |
| — membawa kode mata uang | **432** |
| — membawa **tepat satu** kode | **432** |
| — membawa **lebih dari satu** | **0** |

Sebarannya: `CoverageList` 263 · `AnekaList` 151 · `ASMCoverage` 13 · `PersonList` 5 — **seluruhnya
satu kode per baris**.

⛔ **Satu kolom `CURRENCY_CODE` per tabel karena itu TERBUKTI memadai.** Tidak ada pemecahan baris
yang diperlukan.

📌 Wadah `CurrencyList`, `TotalTSIList` dan `TotalTSIPremiGrossList` memang dapat memuat **dua** kode
(`[terverifikasi]` 4 simpul, seluruhnya di `NB-173649`) — tetapi itu **list yang barisnya memang
per-mata-uang**, jadi tiap barisnya tetap satu kode. Bukan fan-out.

⚠️ **Catatan bentuk yang sempat menyesatkan saya sendiri:** kode mata uang berada **di dalam baris**,
bukan sebagai anak langsung wadahnya — jalurnya `CoverageList/rowdata/<IDR>/Name`. Pengukuran yang
melihat anak langsung `CoverageList` menemukan **nol** dan menyimpulkan salah.

### (c) WARISI — 8 tabel / 33 kolom

`T_LISTINSTALLMENT` (10) · `T_RETROLIST` (5) · `T_FR_SPREADINGLIST` (5) · `T_SPREADINGLIST` (5) ·
`T_COVERAGEDATALIST` (3) · `T_ADDITIONALCOVERAGE` (2) · `T_FR_FACOUTOBJECTLIST` (2) ·
`T_PROPERTY` (1)

⚠️ **Beberapa tabel punya DUA induk** — jangan menelusuri satu rantai saja.
`T_FR_LISTINSTALLMENT` berinduk `T_FR_CURRENCYLIST` **dan** `T_FR_PAYMENT`; `T_FR_PAYMENT` bersumber
dari `NetPerCurrency` **dan** `OfferedPayment`. **Keduanya ikut beres begitu kelompok (b)
dipulihkan.**

### (d) ⛔ TANPA KOLOM MATA UANG — **inilah keputusan yang sesungguhnya**

Lima total skalar di akar penawaran:

`TOTAL_TSI_NUSA_RE` · `TOTAL_PREMI_NUSA_RE` · `TOTAL_TSI_NUSA_RE_SPREADING` ·
`TOTAL_TSI_TOP_RISK` · `SUM_TOTAL_TSI`

⛔ **Alasannya BUKAN mata uangnya hilang, melainkan TIDAK ADA satu mata uang pun yang benar untuk
diisikan — angkanya campuran.**

`[terverifikasi]` dari `DDL\CONTOH\NB-173649.xml`, penawaran **dua mata uang**:

| Skalar akar | Hubungan yang terukur | Hasil |
| --- | --- | :-: |
| `TotalTSINusaReSpreading` | = jumlah **2 baris** `TotalTSIPremiSpreadRNM.TSISpreaded`, ber-`TreatyName` **USD** dan **IDR** | ✅ **SAMA** |
| `TotalPremiNusaRe` | = jumlah **17 baris** `CoverageList.PremiNusantaraRe` | ✅ **SAMA** |
| `SumTotalTSI` | = satu nilai `LocationList/Property.TotalTSI` | ✅ **SAMA** |
| `TotalTSINusaRe` | vs jumlah `CoverageList.TSINusantaraRe` | ⛔ **BEDA** — ⚠️ **JANGAN TEBAK** |
| `TotalTSITopRisk` | vs `LocationList/Property.TotalTSI` | ⛔ **BEDA** — ⚠️ **JANGAN TEBAK** (bernilai **nol** di berkas ini) |

⛔ **Jadi sistem lama menjumlahkan USD dan IDR tanpa konversi.** Itu **kandidat perbaikan**, tetapi
**diport apa adanya** (`CLAUDE.md` §1) — memperbaikinya diam-diam akan membuat rekonsiliasi paralel
run berbeda dan **menyembunyikan salah-port**.

### Struktur lain sudah benar — yang tidak pecah hanya kelima skalar akar

`[terverifikasi]` di `NB-173649`, struktur lain menangani multi mata uang dengan **memecah baris per
mata uang**:

| Wadah | Bentuk | Isi |
| --- | --- | --- |
| `CurrencyList` | 2 entri | USD · IDR |
| `CedingCedantList/CurrencyList` | 2 entri | USD · IDR |
| `TotalTSIList` | 2 entri | USD · IDR |
| `TotalTSIPremiGrossList` | 2 entri | USD · IDR |
| `TotalTSIPremiSpreadRNM` | 2 `rowdata` | `TreatyName` USD · IDR |
| `PropertyItemList` | 5 baris | `Currency` **3 IDR + 2 USD** |

⚠️ **Catatan bentuk yang mengikat cara mengukur** — *diperjelas 24 September.*

`[terverifikasi]` `CurrencyList`, `TotalTSIList` dan `TotalTSIPremiGrossList` **tidak memakai
`<rowdata>`**. Barisnya adalah **elemen bernama kode mata uang** — `<USD>`, `<IDR>` — dan tiap
elemen itu:

- ber-atribut **`REPEATINGINDEX`** (`1`, `2`, …) seperti baris berulang Pega pada umumnya, dan
- memuat medan **`Name`** yang berisi **kode mata uang itu sendiri**.

Jadi ia **tetap kelompok berulang Pega**; yang tidak lazim hanyalah **nama elemennya** dipakai
sebagai kunci baris, bukan `rowdata`.

Sapuan 115 contoh: `CurrencyList` **115 dari 115** berkas memakai bentuk itu dan **0** memakai
`rowdata`; `TotalTSIPremiGrossList` **103**; `TotalTSIList` **21**; **tidak satu berkas pun
mencampur kedua bentuk**. ⛔ Penghitung baris yang mencari `rowdata` mengembalikan **nol** dan
melewatkan seluruh datanya.

⚠️ **Jebakan sejenis pada XPath.** `//CedingCedantList/CurrencyList` **tidak cocok apa pun**, karena
jalur sesungguhnya `CedingCedantList/rowdata/CurrencyList`. `[terverifikasi]` `CedingCedantList`
memang hanya **1 baris**, tetapi di bawahnya ada **2 entri `CurrencyList`** (USD, IDR) dan
**keduanya ber-`Premium` terisi** — sehingga bukti `TotalPremiNusaRe` berdiri **utuh**.

### Konsekuensi

1. Kelompok **(b)** dikerjakan sebagai **perbaikan rancangan**, tidak menunggu siapa pun.
2. Kelompok **(d)** memerlukan **tipe ketiga** di samping `Money` dan `Ratio` — **ADR-0007**.
3. `[terverifikasi]` **Sembilan kolom keluar** dari daftar kolom uang setelah diuji ke isi korpus
   **tanpa cabang `OldData`**; `T_FR_LOCATIONLIST` dan `T_FR_PRINTRISLIP` **keluar seluruhnya**:

   | Kolom | Nilai terisi | Isi | FacIn · retro |
   | --- | ---: | --- | --- |
   | `T_GENERAL_POLIS.CEDING_RETENTION` | **114** | semuanya 0 | 114 · 0 |
   | `T_GENERAL_POLIS.CEDING_RETENTION_NOMINAL` | **5** | semuanya 0 | 5 · 0 |
   | `T_GENERAL_POLIS.PPN_CHECK` | **74** | boolean `true`/`false` | 74 · 0 |
   | `T_GENERAL_POLIS.PAY_RI_COMMISION` | **0** | tidak pernah terisi | — |
   | `T_FR_LOCATIONLIST.LOSS_RATIO1_YEAR_AMOUNT` | **2** *(196 seluruh korpus)* | semuanya 0 | 194 · 2 |
   | `T_FR_LOCATIONLIST.LOSS_RATIO35_YEAR_AMOUNT` | **2** *(196)* | semuanya 0 | 194 · 2 |
   | `T_FR_PRINTRISLIP.WARR_PAYMENT` | **15** | bukan angka | 0 · 15 |
   | `T_FR_LISTINSTALLMENT.STAMP` | **4** *(164)* | semuanya 0 | 160 · 4 |
   | `T_FR_PERSONLIST.ASMCC_AMOUNT` | **8** | seluruhnya bernilai `1` — **penanda** | 5 · 8 |

   ⚠️ Kolom *"nilai terisi"* menghitung **cabang tabel itu sendiri**; angka dalam kurung adalah
   seluruh korpus termasuk kembaran Fac In (butir terbuka 9b). **Vonis tidak berubah** oleh koreksi
   ini — hanya angkanya.
4. ⚠️ **Butir 8 TIDAK ditutup keputusan ini** — lihat ADR-0007 bagian akibat.
5. **Nullability** *(ditambahkan 24 September)*: lembar `Kolom` kini punya kolom **`NULL`**.
   `NOT NULL` hanya untuk empat kolom sistem — `ID`, `PARENT_ID`, `SEQ_NO`, `ROW_UID` (**252**
   kolom); sisanya **`NULL`** (**1.058**). Dasarnya: di korpus sebagian besar medan kosong di
   sebagian besar berkas, sehingga `NOT NULL` pada kolom data akan **menolak baris lama** saat
   migrasi.
   ⚠️ Sub-aturan *"kunci tamu bernama ikuti catatan Daftar Relasi"* **tidak berefek**:
   `[terverifikasi]` **nol** baris `Daftar Relasi` bertanda `nullable`, dan seluruh kolom `*_ID`
   berjenis **`data`**, bukan sistem.

### ⬜ USULAN yang MENUNGGU persetujuan — belum diterapkan

Kolom mata uang kelompok **(a)**, **(b)** dan **(c)** sebaiknya **`NOT NULL`** dengan **nilai
sentinel `'UNKNOWN'` yang eksplisit** bila tidak diketahui, sesuai **ADR-0006** — supaya ketiadaan
mata uang **bersuara, bukan senyap**. Itu justru maksud keputusan ini.

⛔ **Belum diterapkan.** Lembar `Kolom` saat ini menandai seluruh kolom mata uang `NULL`. Tercatat
sebagai butir **10** di lembar `Butir Terbuka`.

### `[pertanyaan terbuka]` yang TIDAK ditutup keputusan ini

1. **Kembaran di sisi Fac In** yang lolos uji yang sama tetapi **tidak termasuk kesembilan yang
   disetujui** — *angka diukur ulang tanpa `OldData`, 24 September* → **butir terbuka 9b**:

   | Kolom | Nilai terisi | Isi |
   | --- | ---: | --- |
   | `T_LOCATIONLIST.LOSS_RATIO1_YEAR_AMOUNT` | **194** | semuanya **0** |
   | `T_LOCATIONLIST.LOSS_RATIO35_YEAR_AMOUNT` | **194** | semuanya **0** |
   | `T_LISTINSTALLMENT.STAMP` | **160** | semuanya **0** |
   | `T_PERSONLIST.ASMCC_AMOUNT` | **5** | 4 bernilai `1` + 1 bernilai besar |

   ⛔ Bila dikoreksi, `T_LOCATIONLIST` **keluar dari kelompok (a)** karena tidak lagi punya kolom
   uang. **Tidak dikoreksi di sini** — di luar kesembilan yang disetujui work owner.
2. Apakah `TotalTSINusaRe` dan `TotalTSITopRisk` punya rumus lain yang belum terbaca.

---

## K-064 · `T_WORK_POLIS` dan `T_GENERAL_POLIS` adalah **tabel fisik yang SAMA** dengan Treaty In

**Tanggal:** 24 September 2026 · **Menutup:** butir **1b** rekonsiliasi tabel flat

**Pertanyaan asal:** apakah `T_WORK_POLIS` dan `T_GENERAL_POLIS` pada rancangan tabel flat Fac In
adalah tabel yang **sama** dengan yang sudah dipakai Treaty In dan PremiumList di skema rumah
Nusantara Re, atau tabel tersendiri yang kebetulan senama.

**Keputusan work owner:** **tabel yang SAMA.** Fac In **menyambung** ke akar yang sudah ada, tidak
membuat akar sendiri. Pembedanya sebuah **kolom penanda lini** yang membedakan **facultative** dari
**treaty**.

### Mengapa ini konsisten dengan yang sudah terbaca

`[terverifikasi]` Dicatat di `08-flat\00-rekonsiliasi-rancangan-flat.md` §3.1 dari baris relasi **52**
diagram skema rumah:

```
T_WORK_POLIS -> T_GENERAL_POLIS   SHARED PK (ID = ID)   1:1
T_WORK_POLIS dipakai bersama PremiumList
```

dan `[terverifikasi]` skema rumah memakai **43 nama `T_*`, nol `FLAT_*`** — dasar yang sama yang
menutup **butir 1**.

📌 **Preseden pola penanda lini sudah ada di skema rumah**, di sisi klaim: tabel akar lintas-lini
`T_WORK_CLAIM` membawa kolom **`LINI`** (nilai terbaca: `FAC`), dan ID-nya berawalan per lini —
`CLM-`/`KMT-` untuk FAC, `CLMP-`/`TKMT-` untuk PROP, `CLMNP-`/`KMTNP-` untuk NONPROP,
`CLMLF-`/`KMTLF-` untuk LIFE. Diagram itu juga menetapkan: **"kolom khas satu lini WAJIB nullable"**.

### Konsekuensi yang mengikat

1. ⛔ **`T_WORK_POLIS` dan `T_GENERAL_POLIS` bukan lagi milik Fac In untuk didefinisikan sendiri.**
   Kolomnya adalah **gabungan** Fac In + Treaty In + PremiumList.
2. **Kolom khas Fac In WAJIB nullable** — sejalan aturan nullability K-063 konsekuensi 5, yang
   sudah menandai seluruh kolom selain empat kolom sistem sebagai `NULL`.
3. **Kolom penanda lini ditambahkan ke `T_WORK_POLIS`.** `[dugaan]` namanya **`LINI`** mengikuti
   preseden `T_WORK_CLAIM.LINI`; ⚠️ **nama dan daftar nilai sahnya belum dikonfirmasi** dari sisi
   Treaty In.
4. ⚠️ **Butir 5b hampir ikut terjawab.** Treaty In menaut versi dengan **`OLD_POLIS_ID`** (baris
   relasi 53: nullable, UNIK, *"pengganti OldData"*). Satu tabel fisik **tidak dapat memakai dua
   mekanisme versi berbeda**, jadi Fac In kemungkinan besar **mewarisi `OLD_POLIS_ID`** dan bukan
   pencarian `PROD_KE` (V-38). ⛔ **Tidak diputuskan di sini** — butir 5b tetap terbuka, tetapi
   pilihannya kini menyempit.
5. **Ruang ID tidak boleh bertabrakan.** Fac In memakai `IDPEGA` berbentuk
   `ASM-FW-GISFW-WORK NB-nnnnn`; bentuk Treaty In belum terbaca.

### ⛔ Yang MENGHALANGI penerapan keputusan ini

`[terverifikasi]` Berkas **`Diagram-Skema-Tabel-NusantaraRe.xlsx` tidak ada di mana pun di bawah
`D:\migrasi\RNM\`** — dihapus ke Recycle Bin 24 September 2026 pukul 15.31 (`08-flat\00a` §1).

Selama berkas itu tidak kembali, **tiga hal tidak dapat dikerjakan**:

| | Tidak dapat dikerjakan | Sebab |
| :-: | --- | --- |
| a | Rekonsiliasi kolom `T_WORK_POLIS` (**18** kolom Fac In) dan `T_GENERAL_POLIS` (**74**) terhadap sisi Treaty In | daftar kolom Treaty In hanya ada di berkas itu |
| b | Konfirmasi nama dan nilai sah kolom penanda lini | preseden `LINI` hanya terbaca di berkas itu |
| c | Pemeriksaan tabrakan ruang `IDPEGA` antar lini | bentuk ID Treaty In hanya terbaca di berkas itu |

⚠️ **Sebelum 1b dijawab, hilangnya berkas itu hanya catatan ketertelusuran. Sesudah dijawab, ia
menjadi penghalang** — karena tabel bersama menuntut rekonsiliasi kolom, dan rekonsiliasi itu
mustahil tanpa daftar kolom lawannya. **Kembalikan berkas itu ke `DDL\`.**

### `[pertanyaan terbuka]` yang TIDAK ditutup keputusan ini

1. Nama dan daftar nilai sah kolom penanda lini.
2. Berapa dari 74 kolom `T_GENERAL_POLIS` Fac In beririsan dengan Treaty In, dan adakah nama sama
   berarti beda.
3. Apakah `IDPEGA` Fac In dan Treaty In dijamin tidak bertabrakan.
4. Butir **5b** — lihat konsekuensi 4.

---

## K-065 · Diagram skema rumah **TIDAK dipakai** — sumber pembuatan tabel flat adalah `Claude outputs\`

**Tanggal:** 24 September 2026 · **Menutup:** butir **1d** · **Mencabut** penghalang yang dicatat K-064

**Pertanyaan asal:** `Diagram-Skema-Tabel-NusantaraRe.xlsx` hilang dari `DDL\`, dan K-064 mencatatnya
sebagai **penghalang** karena tabel akar bersama menuntut rekonsiliasi kolom terhadap Treaty In.

**Keputusan work owner:** berkas diagram itu **tidak dipakai**. Sumber pembuatan tabel flat adalah
apa yang ada di **`D:\migrasi\RNM\Claude outputs\`**.

### Konsekuensi

1. ⛔ **Butir 1d gugur** — rekonsiliasi kolom terhadap Treaty In **tidak dikerjakan**, dan tidak lagi
   menahan pekerjaan apa pun.
2. **Butir 1c tetap terbuka** tetapi turun bobotnya: nama kolom penanda lini (`LINI`, `[dugaan]`)
   dapat ditetapkan sendiri tanpa menunggu berkas itu.
3. ⛔ **Penghalang yang dicatat K-064 DICABUT.** Bagian *"Yang MENGHALANGI penerapan keputusan ini"*
   pada K-064 **tidak lagi berlaku**; ia dipertahankan sebagai jejak, bukan sebagai tugas.

### ⚠️ Risiko yang dicatat, bukan disembunyikan

**K-064 tetap berlaku** — `T_WORK_POLIS` dan `T_GENERAL_POLIS` adalah tabel yang **sama** dengan
Treaty In. Tanpa rekonsiliasi, kolom kedua tabel itu dibangun **hanya dari sisi Fac In** (18 dan 74
kolom). Tiga hal karena itu **tidak diketahui** dan baru akan muncul saat integrasi:

- berapa kolom Fac In yang beririsan dengan Treaty In, dan adakah **nama sama berarti beda**;
- apakah ruang `IDPEGA` kedua lini **dijamin tidak bertabrakan**;
- daftar nilai sah kolom penanda lini.

⛔ Ini **keputusan work owner yang sadar**, dicatat di sini supaya tidak tampak sebagai kelalaian
(`CLAUDE.md` §1).

---

## K-066 · JSON ≡ XML untuk **seluruh** COB, dan **struktur RNW = NB**

**Tanggal:** 24 September 2026 · **Menutup:** dua butir yang tidak dapat diverifikasi dari korpus

**Pertanyaan asal:** dua celah yang tercatat di `08-flat\BAHAN-SPEC-PEMUATAN.md` §9.

### (a) Kesepadanan JSON dan XML

**Pertanyaan:** loader produksi membaca **JSON** (`JSON_POLIS.DATA_JSON`), tetapi seluruh rancangan
diturunkan dari **XML**. Kesepadanan hanya dapat diuji pada **2** kasus — keduanya NB — karena hanya
dua kasus yang tersedia dalam kedua bentuk.

**Keputusan work owner:** **sepadan untuk seluruh COB.** Tidak perlu contoh berpasangan tambahan.

`[terverifikasi]` Bukti yang ada mendukungnya: pada `NB-181231`, sesudah ruas kode mata uang
dinormalkan, **621 dari 621** jalur daun JSON ada di XML; satu-satunya jalur XML-saja adalah
`pzStatus`, metadata Pega.

⛔ **Yang tetap berlaku dan BUKAN soal kesepadanan:** perbedaan **bentuk** antara keduanya nyata dan
wajib ditangani loader — di JSON `CurrencyList` adalah **array** tanpa kunci kode; di XML ia
**elemen bernama kode**. Pembawa mata uang tetap medan **`Name`** di kedua bentuk.

### (b) Struktur RNW

**Pertanyaan:** `[terverifikasi]` **nol** contoh RNW di antara 115 berkas, sehingga seluruh cakupan
RNW di buku kerja **diwarisi dari NB, tidak diukur**.

**Keputusan work owner:** **struktur RNW = NB.** Warisan itu sah; tidak perlu contoh RNW.

`[terverifikasi]` Konsisten dengan dua temuan lama: `08-flat\01` §0.1 (*"RNW dan EDM memakai struktur
yang SAMA dengan NB"*, dari satu contoh RNW) dan **NB↔RNW nol berbeda** pada perbandingan rule.

### Konsekuensi

1. Butir 1 dan 6 pada daftar *"yang tidak dapat diverifikasi"* di `BAHAN-SPEC-PEMUATAN.md` §9
   **ditutup**.
2. Loader **tetap** wajib menangani perbedaan bentuk array-vs-elemen-bernama-kode.
3. ⚠️ Sel RNW di buku kerja tetap berbunyi **"warisan NB"** — itu tetap **pernyataan yang jujur**:
   diwarisi atas keterangan work owner, bukan diukur dari data.

---

## K-067 · Tabel `Total*` **informasi saja** — aman dihapus, dan **ADR-0001 tidak perlu diubah**

**Tanggal:** 24 September 2026 · **Menutup:** butir **8** rekonsiliasi tabel flat

**Pertanyaan asal:** **V-28** membuang `TotalTSIList`, `TotalTSIPremiGrossList` dan
`TotalTSIPremiSpreadRNM` karena isinya penjumlahan baris yang sudah tersimpan — tetapi catatan V-28
sendiri menyimpulkan rekonsiliasi *"perlu **ambang toleransi**, bukan perbandingan sama persis"*,
sedangkan **ADR-0001** menolak toleransi secara eksplisit. Keduanya tidak dapat berlaku bersama.

**Keputusan work owner:** penjumlahan di ketiga tabel itu **untuk informasi saja** — **tidak dipakai
untuk pembayaran atau hal lain**. **Aman dihapus.**

### ⛔ Pertentangan dengan ADR-0001 BUBAR — dan tidak satu pun jalan A/B/C/D diperlukan

Empat jalan yang sempat disusun (**A** amandemen ADR-0001 · **B** simpan apa adanya · **C** beda
perlakuan data baru/lama · **D** reproduksi urutan pembulatan) **seluruhnya berangkat dari premis
yang sama: bahwa nilai itu akan DIHITUNG ULANG lalu dibandingkan.**

⛔ **Premis itu gugur.** Nilainya **tidak disimpan dan tidak dihitung ulang** — jadi **tidak ada yang
dibandingkan**, dan selisih pembulatan 2,6 per sepuluh juta **tidak pernah terwujud**.

✅ **ADR-0001 tetap utuh, tanpa amandemen.** Ukuran keberhasilan migrasi tidak berubah.

### Keadaan rancangan — sudah sesuai, tidak ada yang perlu disunting

`[terverifikasi]` atas `08-flat\Tabel-Flat-Lintas-Siklus.xlsx` dan `DDL-tabel-flat-draf.sql`:

| | Keadaan |
| --- | --- |
| `T_TOTALTSILIST` · `T_TOTALTSIPREMIGROSSLIST` · `T_TOTALTSIPREMISPREADRNM` | ⛔ **sudah tidak ada** di 78 tabel — V-28 memang sudah diterapkan |
| `T_FR_TOTALTSIPREMIRETRO` | ✅ **tetap ada** — V-28 sendiri mengecualikannya (*"tidak disebut dan Premium-nya tidak cocok rumus mana pun"*) |
| Lima skalar akar di `T_GENERAL_POLIS` | ✅ **tetap ada**, disimpan **apa adanya** (ADR-0007), **tidak** dihitung ulang |

### ⚠️ Dua pagar yang harus ikut dicatat

1. ⛔ **Kelima skalar akar tetap TIDAK dihitung ulang.** Mereka nilai terekam (ADR-0007). Menghitung
   ulangnya akan memunculkan kembali selisih pembulatan yang baru saja dihindari — dan pada
   `TOTAL_TSI_NUSA_RE_SPREADING` juga menuntut penjumlahan **lintas mata uang** (ADR-0007), yang
   selisihnya **bukan pembulatan melainkan beda hasil**.
2. ⚠️ **Angka total yang dihitung ulang untuk tampilan TIDAK BOLEH dipakai dalam rekonsiliasi.**
   Untuk tampilan, selisih 2,6 per sepuluh juta tidak material — keputusan ini menyatakannya begitu.
   Untuk rekonsiliasi, ADR-0001 tetap menuntut nol selisih, dan angka hasil hitung ulang **bukan
   pembanding yang sah**.

### Konsekuensi

1. **Butir 8 tertutup.** Pertanyaan yang sempat disiapkan untuk penyusun V-28 (*"pembulatan
   per-langkah atau sekali di akhir?"*) **tidak lagi perlu diajukan**.
2. **Tidak ada perubahan pada rancangan, DDL, maupun ADR.** Keputusan ini **menegaskan** keadaan yang
   sudah ada, bukan mengubahnya.
3. `[terverifikasi]` Uji aritmetika V-28 (106 contoh) **tidak pernah saya ukur ulang** — dan kini
   **tidak perlu**, karena angkanya tidak dipakai untuk apa pun.

---

## K-068 · Fac In memakai **`OLD_POLIS_ID`** untuk menaut versi — satu kolom

**Tanggal:** 24 September 2026 · **Menutup:** butir **5b** rekonsiliasi tabel flat

**Pertanyaan asal:** bagaimana sistem baru menemukan **versi sebelumnya** sebuah polis — penunjuk
langsung `OLD_POLIS_ID` seperti Treaty In, atau pencarian **`PROD_KE` tertinggi** seperti V-38/V-40.

**Keputusan work owner:** **`OLD_POLIS_ID`** — cukup **satu kolom**.

### Mengapa pilihannya memang tinggal satu

⛔ **K-064** menetapkan `T_GENERAL_POLIS` adalah **tabel fisik yang SAMA** dengan Treaty In. **Satu
tabel tidak dapat memakai dua mekanisme versi.** Treaty In sudah memakai `OLD_POLIS_ID`, jadi Fac In
mewarisinya.

`[terverifikasi]` Bentuknya tercatat di `08-flat\00-rekonsiliasi-rancangan-flat.md` §3.3, dikutip
dari baris relasi **53** diagram skema rumah:

```
T_WORK_POLIS | T_GENERAL_POLIS | OLD_POLIS_ID | 1:1 | di Go | UNIK
             | nullable . penunjuk ke generasi SEBELUMNYA . pengganti OldData . UNIK melarang percabangan
```

📌 Preseden juga sudah ada di basis data yang hidup: `[terverifikasi]` `JSON_POLIS` memuat kolom
**`OLDNOPOLIS VARCHAR2(100)`** — penunjuk ke polis sebelumnya bukan gagasan baru bagi skema ini.

### Yang sudah diterapkan

`[terverifikasi]`

| Berkas | Perubahan |
| --- | --- |
| `Tabel-Flat-Lintas-Siklus.xlsx` lembar `Kolom` | `T_GENERAL_POLIS.OLD_POLIS_ID` `NUMBER` `NULL`, jenis **sistem** — `T_GENERAL_POLIS` kini **75** kolom |
| lembar `Daftar Relasi` | satu baris: `T_GENERAL_POLIS → T_GENERAL_POLIS`, kunci `OLD_POLIS_ID`, **1:1**, **di Go**, **UNIK** |
| `DDL-tabel-flat-draf.sql` | kolom ditambahkan + `CREATE UNIQUE INDEX UX_GENERAL_POLIS_OLDPOLIS` |

⛔ **TANPA `FOREIGN KEY`.** Baris relasi 53 berbunyi **"di Go"** — keutuhannya ditegakkan di kode,
bukan oleh basis data. DDL karena itu memuat **index UNIK saja**, bukan constraint.
`[terverifikasi]` jumlah FK tetap **64**; index naik **76 → 77**.

⚠️ **UNIK itu bukan hiasan:** ia **melarang percabangan** — dua versi tidak boleh mengaku pendahulu
yang sama.

### Konsekuensi

1. **Cara mencari versi sebelumnya berubah.** V-38 dan V-40 menetapkan pencarian lewat `POLICY_NO`
   sama dengan `PROD_KE` tertinggi di bawah baris berjalan. Itu **digantikan** penunjuk langsung.
   ⚠️ **Bunyi V-38 dan V-40 perlu ditulis ulang** supaya tidak terus terbaca bertentangan — sama
   seperti V-40/V-40b yang perlu diperbaiki setelah butir 3.
2. **Butir 6 sebagian besar gugur.** `PROD_KE NUMBER(5)` tetap benar sebagai **nomor versi**, tetapi
   ia **tidak lagi memikul pencarian versi**, sehingga alasan V-38 tentang "pengurutan teks menaruh
   10 sebelum 9" **tidak lagi menentukan**.
3. **Loader terpengaruh.** `BAHAN-SPEC-PEMUATAN.md` §6 mencatat butir 5b menentukan cara loader
   menaut versi — kini tertetapkan.
4. ⚠️ **`ROW_UID` TIDAK ikut terselesaikan.** `OLD_POLIS_ID` menaut **polis** antar versi;
   `ROW_UID` menaut **baris** antar versi. Rekonstruksi `ROW_UID` untuk data lama (V-42) **tetap
   terbuka**, dan tetap bergantung pada apakah migrasi memuat seluruh versi atau versi terakhir saja.

---

## K-069 · Empat butir terakhir rancangan tabel flat — **7b · 9b · 1c · usulan 10**

**Tanggal:** 25 September 2026 · **Menutup:** butir **7b**, **9b**, **1c**, dan **usulan 10**
⛔ **Dengan ini seluruh butir rekonsiliasi tabel flat TERTUTUP.**

---

### (7b) `IsCedingConfirm` mendapat **kolom sendiri**

**Pertanyaan asal:** V-32 mengusulkan memetakannya ke `HISTORYAKSEPTASIPRODUCTION.POSISI`, tetapi
V-48 memakai `POSISI` untuk **posisi tangga akseptasi**.

**Keputusan work owner:** **kolom sendiri.**

⛔ Satu kolom tidak boleh memikul dua arti. Bila `IsCedingConfirm` ditulis ke `POSISI`, maka
`T_WORK_POLIS.POSISI` akan mencerminkan **konfirmasi ceding**, bukan posisi tangga akseptasi —
dan `[terverifikasi]` V-48c mencatat daftar nilai sah `POSISI` **belum pernah terlihat** (hanya DDL,
nol baris data), jadi tabrakan itu tidak akan ketahuan dari data.

---

### (9b) Empat kolom kembaran Fac In **DIKELUARKAN** dari daftar kolom uang

**Keputusan work owner:** **iya, dikeluarkan.**

`[terverifikasi]` diukur atas 115 contoh, cabang `OldData` dikecualikan:

| Kolom | Nilai terisi | Isi |
| --- | ---: | --- |
| `T_LOCATIONLIST.LOSS_RATIO1_YEAR_AMOUNT` | **194** | semuanya **0** |
| `T_LOCATIONLIST.LOSS_RATIO35_YEAR_AMOUNT` | **194** | semuanya **0** |
| `T_LISTINSTALLMENT.STAMP` | **160** | semuanya **0** |
| `T_PERSONLIST.ASMCC_AMOUNT` | **5** | 4 bernilai `1` + 1 besar — **penanda** |

**Akibatnya pada pengelompokan K-063:**

| Kelompok | Sebelum | **Sesudah** |
| --- | --- | --- |
| (a) AMAN | 8 tabel / 25 kolom | **7 / 23** — ⛔ `T_LOCATIONLIST` **keluar** |
| (b) PULIHKAN | 10 / 61 | **9 / 60** — ⛔ `T_PERSONLIST` **keluar** |
| (c) WARISI | 8 / 33 | **8 / 32** |
| (d) TANPA mata uang | 7 / 28 | **7 / 28** |
| **Total** | 33 / 147 | **31 / 143** |

⚠️ Kolom `T_PERSONLIST.CURRENCY_CODE` dan `T_LOCATIONLIST.CURRENCY_ID` **tetap ada** sebagai data
korpus — keduanya terbaca di korpus — tetapi tidak lagi memikul kolom uang, jadi **dibiarkan `NULL`**.

---

### (1c) Kolom penanda lini **TIDAK ADA di proyek ini**

**Keputusan work owner:** *"itu tidak ada di dalam proyek ini"*.

⛔ **Kolom `LINI` yang ditambahkan 24 September DIHAPUS kembali** dari `T_WORK_POLIS`.

📌 Ini **tidak membatalkan K-064**: `T_WORK_POLIS` dan `T_GENERAL_POLIS` tetap tabel yang **sama**
dengan Treaty In. Yang ditetapkan di sini adalah **lingkup**: penanda facultative-vs-treaty bukan
milik proyek ini untuk didefinisikan. Ia akan ada ketika Treaty In diintegrasikan — sejalan dengan
keterangan awal work owner *"nanti ada kolom"* dan dengan **K-065** yang mengeluarkan diagram skema
rumah dari lingkup.

---

### (usulan 10) Mata uang **`NOT NULL` + sentinel `'UNKNOWN'`** — DISETUJUI

**Keputusan work owner:** **setuju.**

`[terverifikasi]` diterapkan pada **24 kolom** pembawa mata uang:

| Kelompok | Kolom dikeraskan |
| --- | ---: |
| (a) | **7** — `CURRENCY` / `CURRENCY_CODE` pada 7 tabel |
| (b) | **9** — `CURRENCY_CODE` yang ditambahkan 24 September |
| (c) | **8** — ⛔ **kolom BARU**, lihat di bawah |

Bentuknya: `VARCHAR2(10) DEFAULT 'UNKNOWN' NOT NULL`.

⛔ **`DEFAULT`-nya yang membuat `NOT NULL` aman.** Tanpa default, `NOT NULL` akan **menolak baris
lama** saat migrasi — persis alasan K-063 membuat seluruh kolom data `NULL`. Dengan default, baris
lama tanpa mata uang masuk sebagai **`UNKNOWN`**, terekam, dan **bersuara**.

#### ⛔ Kelompok (c) menuntut **delapan kolom baru**

K-063 (c) berbunyi *"leluhurnya punya — diisi saat tulis"*. `[terverifikasi]` kedelapan tabelnya
**tidak punya kolom mata uang sama sekali**, jadi tidak ada yang dapat diisi. Kolom `CURRENCY_CODE`
**ditambahkan** ke: `T_LISTINSTALLMENT` · `T_RETROLIST` · `T_FR_SPREADINGLIST` · `T_SPREADINGLIST` ·
`T_COVERAGEDATALIST` · `T_ADDITIONALCOVERAGE` · `T_FR_FACOUTOBJECTLIST` · `T_PROPERTY`.

⚠️ Ini **konsekuensi** usulan 10, bukan keputusan tambahan — tanpa kolomnya, "bersuara" mustahil.

#### Yang TIDAK dikeraskan, dan sebabnya

| Kolom | Sebab |
| --- | --- |
| `CURRENCY_ID` · `CURRENCY_OLD_ID` | `[terverifikasi]` isinya **ANGKA** rujukan master, bukan kode ISO |
| `FLAG_CURRENCY` | penanda, bukan kode mata uang |
| `T_PERSONLIST.CURRENCY_CODE` | keluar dari kelompok (b) lewat butir 9b |

#### ⚠️ Batas terhadap ADR-0006

ADR-0006 menetapkan **`panic` saat menulis `Unknown` ke Oracle**. Sentinel ini **tidak
melanggarnya**: aturan itu berlaku bagi **tulisan BARU aplikasi**, bukan bagi **baris LAMA yang
dimigrasikan**. Sentinel merekam fakta bahwa baris lama memang datang tanpa mata uang —
persis efek samping yang ADR-0006 sebut *"menghasilkan hitungan berapa banyak nilai produksi yang
tiba tanpa mata uang"*.

---

### Keadaan sesudah K-069 — `[terverifikasi]`

| | Nilai |
| --- | ---: |
| Tabel | **78** |
| Kolom | **1.329** |
| Kolom uang | **143** di **31** tabel |
| Kolom mata uang `NOT NULL DEFAULT 'UNKNOWN'` | **24** |
| `PRIMARY KEY` · `FOREIGN KEY` · `INDEX` | **78** · **64** · **77** |

DDL diregenerasi penuh dan lolos delapan pemeriksaan: nol kolom ganda · nol tabel tanpa `ID` ·
**24 `NOT NULL` seluruhnya ber-`DEFAULT`, nol tanpa** · nol FK/index bermasalah · nol kata tercadang ·
nol nama ganda atau melebihi 30 karakter · **78/78 cocok dengan lembar `Kolom`**.

---

## K-070 · Migrasi data lama dipecah **DUA FASE** — fase 1 versi terakhir, fase 2 generasi sebelumnya

**Tanggal:** 25 September 2026 · **Menutup:** butir **A-1** `_DAFTAR-ISSUE-TERBUKA.md` ·
`08-flat\BAHAN-SPEC-PEMUATAN.md` §5 dan §9 butir 5 · lembar `Butir Terbuka` baris 25

**Pertanyaan asal:** apakah migrasi memuat **seluruh versi** polis atau **versi terakhir saja**.

**Keputusan work owner:** ⛔ **seluruh polis TETAP menjadi sasaran akhir** — tetapi dikerjakan
**bertahap**:

| Fase | Isi | Kapan |
| :-: | --- | --- |
| **1** | **Versi terakhir saja** dari tiap polis | sekarang, bersama pembangunan sistem |
| **2** | **Generasi sebelumnya** (riwayat endorsement) | **setelah sistem jadi** |

⛔ **Ini BUKAN "versi terakhir saja" sebagai jawaban akhir.** Pilihan "versi terakhir" pada daftar
A-1 memuat konsekuensi *"nilai sebelum endorsement hilang permanen"* — konsekuensi itu **tidak
berlaku** di sini, karena fase 2 akan memuatnya. Yang berubah hanya **urutan waktu**, bukan cakupan.

### Konsekuensi yang mengikat

1. **`ROW_UID` tetap wajib dirancang sekarang**, bukan ditunda ke fase 2. Alasannya membalik arah
   pekerjaan: di fase 1 UID **dibangkitkan** (baris lahir tanpa riwayat), sehingga di fase 2 generasi
   lama **tidak bebas memilih UID** — ia harus **menjodohkan ke belakang** terhadap UID yang sudah
   terlanjur ada. ⚠️ Ini **lebih sulit** daripada memuat seluruh versi sekaligus, dan itulah harga
   yang dibayar untuk mendapat sistem lebih cepat.
2. **`OLD_POLIS_ID` (K-068) kosong di seluruh baris fase 1** — tidak ada generasi sebelumnya untuk
   ditunjuk. `[terverifikasi]` di Oracle, `UNIQUE` mengizinkan **banyak** `NULL`, jadi kendala unik
   K-068 **tidak dilanggar**. Kolomnya baru bermakna setelah fase 2.
3. **Rekonsiliasi paralel run (ADR-0001) di fase 1 hanya mencakup generasi terakhir.** ADR-0001
   **tidak dilonggarkan** — nol selisih tetap berlaku, hanya populasinya yang lebih sempit.
4. Cabang `OldData` **tetap dibuang** (V-19). ⛔ Ia bukan pengganti generasi sebelumnya: V-43a
   `[terverifikasi]` menyatakan `OldData` adalah potret salinan kerja **sesudah** baris baru
   dibentuk. Membiarkannya masuk untuk "menambal" fase 1 akan menyimpan nilai yang salah.

### `[pertanyaan terbuka]` yang TIDAK ditutup keputusan ini

Sepuluh pertanyaan turunan lahir dari keputusan ini dan **belum dijawab** — seluruhnya dicatat
sebagai golongan **J** di `_DAFTAR-ISSUE-TERBUKA.md`, ditambah satu permintaan DBA baru (**T-10**).
Empat di antaranya **wajib dijawab sebelum loader ditulis**: kunci pengelompokan generasi ·
aturan memilih "terakhir" · perlakuan nilai `PRODKE` · penomoran endorsement baru.

---

## K-071 · Sepuluh pertanyaan turunan K-070 dijawab — **J-1 … J-10**

**Tanggal:** 25 September 2026 · **Menutup:** golongan **J** `_DAFTAR-ISSUE-TERBUKA.md` ·
**Mengamandemen sebagian:** **K-068** (lihat di bawah)

### Jawaban work owner

| # | Pertanyaan | Jawaban |
| :-: | --- | --- |
| **J-1** | Kunci pengelompokan generasi | `JSON_POLIS` **tidak dipakai lagi** ke depan; yang dipakai **tabel flat baru** |
| **J-2** | Penentu "versi terakhir" | **`PRODKE` dan `TGL_INPUT` sama-sama tertinggi** — keduanya menunjuk baris yang sama |
| **J-3** | Tipe `PRODKE` di tabel flat | **`NUMBER`** |
| **J-4** | Penomoran endorsement baru | **Seluruh data ikut dimigrasikan** begitu sistem Go selesai → endorsement berikutnya tidak bermasalah |
| **J-5** | `OLD_POLIS_ID` | ⛔ **Hanya terisi saat membuat RENEWAL.** NB dan EDM **tidak punya** nilai itu |
| **J-6** | Cara menangani penjodohan fase 2 | **Minta rekomendasi** → lihat di bawah |
| **J-7** | Rekonsiliasi fase 1 | Data lama dimasukkan **sambil menguji** sistem baru |
| **J-8** | Riwayat akseptasi | **Hanya beda penamaan kolom** — di polis lama namanya `IDPEGA` |
| **J-9** | `_MENJADI` / `_SELISIH` | **Periksa dulu** → hasil pemeriksaan di bawah |
| **J-10** | Endorsement belum selesai | Yang dimuat **yang sudah jadi**. `STS_KONVERSI` ada di `JSON_POLIS`, tabel yang tidak dipakai lagi |

---

### ⚠️ Amandemen sebagian **K-068** — `OLD_POLIS_ID` adalah tautan **RENEWAL**, bukan tautan generasi

**Bunyi K-068 tidak dibatalkan; cakupannya dipersempit.** K-068 menetapkan `OLD_POLIS_ID` sebagai
penunjuk ke "generasi sebelumnya", menggantikan pencarian `PROD_KE` tertinggi. **J-5 mempersempitnya:**
kolom itu **hanya** terisi saat renewal; NB dan EDM membiarkannya kosong.

`[terverifikasi]` Penyempitan ini **konsisten dengan skema sumber.** `DDL\JSON_POLIS.txt` memuat
`"OLDNOPOLIS" VARCHAR2(100)` — kolom **nomor polis lama**, dan renewal memang satu-satunya peristiwa
yang **mengganti nomor polis**. Endorsement tidak menggantinya.

⛔ **Yang menjadi kosong karenanya:** apa yang menautkan **generasi endorsement** satu sama lain.
`[dugaan]` pembacaan yang paling sejalan dengan skema:

| Peristiwa | Nomor polis | Tautannya |
| --- | --- | --- |
| **Endorsement** | **tetap** | `NOPOLIS` sama + `PROD_KE` naik |
| **Renewal** | **berganti** | **`OLD_POLIS_ID`** (preseden `OLDNOPOLIS`) |

⚠️ **Ini `[dugaan]` saya, bukan pernyataan work owner.** Menunggu konfirmasi sebelum dipakai
menyusun loader.

---

### Rekomendasi **J-6** — ⛔ **jangan menjodohkan ke belakang; muat ulang dari nol**

**Alasan premisnya berubah.** Kebutuhan "menjodohkan ke belakang" (K-070 konsekuensi 1) lahir dari
anggapan bahwa fase 2 berjalan **sesudah** sistem dipakai sungguhan. **J-7 membatalkan anggapan itu:**
data lama dimasukkan **sambil menguji**, dan **J-4** menegaskan seluruh data ikut begitu sistem
selesai. Selama masih menguji, isi tabel **boleh dibuang dan dimuat ulang**.

| | Jalan A — jodohkan ke belakang | **Jalan B — muat ulang dari nol** ⭐ |
| --- | --- | --- |
| Fase 2 | cocokkan baris lama ke `ROW_UID` yang sudah beku | **kosongkan, muat seluruh generasi sekali jalan** |
| Arah kerja | terbalik, sasaran beku | maju, seperti memuat sejak awal |
| `ROW_UID` | dibangkitkan dua kali, harus rukun | **dibangkitkan sekali**, atas riwayat lengkap |
| Risiko | salah jodoh **tidak bersuara** | — |

⛔ **Syarat yang mengikat Jalan B:** fase 2 **wajib selesai sebelum sistem dipakai sungguhan** —
sebelum ada satu pun kasus baru yang dibuat di sistem baru. Sesudah itu tabel tidak boleh dikosongkan
lagi, dan kita **terpaksa** kembali ke Jalan A.

📌 **Konsekuensi yang menyenangkan:** `ROW_UID` fase 1 menjadi **sementara**. Ia tetap wajib ada
(barisnya butuh identitas), tetapi **tidak perlu kekal** — sehingga rekonstruksi UID yang disebut
V-42 sebagai bagian tersulit **turun bobotnya secara drastis**.

**Status: usulan, menunggu persetujuan work owner.**

---

### Hasil pemeriksaan **J-9** — ⛔ `_MENJADI` / `_SELISIH` **tidak dapat** menggantikan nilai sebelum endorsement

`[terverifikasi]` sapuan **80** berkas `DDL\*.txt`:

| Ukuran | Nilai |
| --- | ---: |
| Berkas yang memuat pasangan | **5** |
| Basis unik `_MENJADI` | **9** |
| Basis unik `_SELISIH` | **15** |
| Basis yang **berpasangan lengkap** | **9** |

Kelima berkasnya: `FACINOFFER`, `FACINPRODUCTION`, `FACINPRODUCTION_BACKUP`, `FACOUTPRODUCTION`,
`TREATYPRODUCTION_BACKUP`.

**Tiga sebab mengapa tidak bisa menggantikan:**

1. **Cakupannya terlalu sempit** — 9 ukuran berpasangan (TSI, TSI100, TSITOP, PREMI,
   PREMI_COVERAGE, COMMISION_COVERAGE, BROKERAGE_FEE, DEDUCTION2, LOL) berhadapan dengan **143 kolom
   uang** pada rancangan flat.
2. **Butirannya salah** — kolom-kolom itu ada di **tabel produksi tingkat penawaran**, bukan di
   wadah berulang tingkat baris (coverage, deductible, spreading, location, person, vehicle) yang
   justru menyimpan nilai yang hilang.
3. **Berada di luar lingkup loader** — kelima tabel itu **sudah flat dan sudah ada**
   (lembar `Tabel Sudah Ada`); `BAHAN-SPEC-PEMUATAN.md` §0 menempatkan jalur tulis produksi di luar
   lingkup.

⚠️ **Sekalipun demikian, arahnya benar:** `[dugaan]` nilai lama dapat diturunkan sebagai
`MENJADI − SELISIH` untuk kesembilan basis berpasangan. Korpus **tidak menyatakan** aturan itu; ia
kesimpulan aritmetika dari nama kolom, jadi **jangan dipakai tanpa konfirmasi**.

#### ⛔ Temuan sampingan — enam `_SELISIH` **tanpa** pasangan `_MENJADI`

`[terverifikasi]` `COMMISION` · `OBJECTPREMI` · `PCT_BROKERAGE_FEE` · `PCT_RI_COMM` · `RICOMM` ·
`SHAREOFFERED`. Diperiksa lebih lanjut apakah kolom **polos**-nya ada sebagai titik acuan:

| Basis | Kolom polos | Nilai lama dapat diturunkan? |
| --- | :-: | --- |
| `RICOMM` · `COMMISION` · `OBJECTPREMI` · `SHAREOFFERED` | **ada** | `[dugaan]` ya — `polos − SELISIH` |
| **`PCT_BROKERAGE_FEE`** · **`PCT_RI_COMM`** | ⛔ **tidak ada** | ⛔ **tidak sama sekali** — selisih tanpa titik acuan |

⚠️ **Menyentuh aturan di `CLAUDE.md` §4.3**, yang berbunyi *"Kolom `_MENJADI` / `_SELISIH`
dipertahankan berpasangan."* `[terverifikasi]` **6 dari 15** basis `_SELISIH` tidak punya pasangan
`_MENJADI`, dan **2** di antaranya tidak punya kolom polos juga. Instruksi itu **tidak dapat
dijalankan secara harfiah** atas keenamnya. ⛔ `CLAUDE.md` berada **di luar `OUTPUT\`** — tidak
disunting dari sini; diangkat agar keputusannya milik work owner.

```powershell
# menghasilkan 80 / 5 / 9 / 15 / 9
$nF=0; $bMj=@{}; $bSl=@{}; $nP=0
foreach ($p in [IO.Directory]::EnumerateFiles('D:\migrasi\RNM\DDL','*.txt')) {
  $nF++; $t=[IO.File]::ReadAllText($p); $ada=$false
  foreach ($m in [regex]::Matches($t,'"([A-Z0-9_]*)_MENJADI"')) { $bMj[$m.Groups[1].Value]=1; $ada=$true }
  foreach ($m in [regex]::Matches($t,'"([A-Z0-9_]*)_SELISIH"')) { $bSl[$m.Groups[1].Value]=1; $ada=$true }
  if ($ada) { $nP++ } }
$cocok=0; foreach ($k in $bMj.Keys) { if ($bSl.ContainsKey($k)) { $cocok++ } }
"$nF / $nP / $($bMj.Count) / $($bSl.Count) / $cocok"
```

---

### Konsekuensi terhadap rancangan

1. **J-2 → aturan pemilihan "versi terakhir" ditetapkan.** Karena `PRODKE` dan `TGL_INPUT`
   **sama-sama tertinggi**, keduanya saling menguji. ⛔ Loader **mengurutkan dengan `TGL_INPUT`**
   (bertipe `DATE`, bebas dari jebakan urut-teks) dan memakai `PRODKE` sebagai **pemeriksa silang**.
   Bila keduanya **tidak sepakat**, itu tanda anggapan ini patah → **`panic`**, bukan diam-diam
   memilih salah satu.
2. **J-3 → `PROD_KE NUMBER` menegaskan butir 6.** Sumbernya `VARCHAR2(5)`, sasarannya `NUMBER`,
   sehingga loader **mengubah tipe** saat memuat. Nilai yang **bukan angka** wajib **`panic`**,
   bukan diam-diam jadi `NULL` atau `0`.
3. **J-8 → pertanyaan jalur bandingnya BUBAR.** Kunci riwayat memang sama (`ID_PEGA` = `IDPEGA`),
   dan karena **J-4 + J-7** memastikan seluruh generasi masuk **sebelum** sistem dipakai sungguhan,
   keadaan timpang "riwayat lengkap, polis sebagian" **tidak pernah terwujud di produksi**. Pola
   penutupan yang sama dengan butir 8 (K-067): premisnya yang gugur, bukan jawabannya yang dipilih.
4. **J-10 → penyaring "sudah jadi" tetap perlu untuk fase 1.** `STS_KONVERSI` memang tidak ikut ke
   sistem baru, tetapi loader **membaca** `JSON_POLIS` sebagai sumber, jadi ia tetap tersedia sebagai
   penyaring **saat memuat**. ⚠️ `belum terverifikasi` arti nilainya — ⛔ **jangan ditebak.** Bila
   Jalan B (J-6) disetujui, bobotnya turun: fase 2 memuat ulang seluruhnya.

### `[pertanyaan terbuka]` yang TIDAK ditutup keputusan ini

1. **Tautan antar generasi endorsement** — `[dugaan]` `NOPOLIS` + `PROD_KE`, menunggu konfirmasi.
2. **Persetujuan Jalan B** pada J-6, berikut syaratnya (fase 2 selesai sebelum pemakaian sungguhan).
3. **Arti nilai `STS_KONVERSI`** sebagai penyaring "sudah jadi".
4. **Sikap terhadap enam `_SELISIH` tanpa pasangan**, khususnya dua yang tanpa kolom polos.

---

## K-072 · Empat sisa K-071 dijawab — **J-11 … J-14**

**Tanggal:** 25 September 2026 · **Menutup:** J-11, J-13, J-14 · **Menjawab pertanyaan** J-12
(syarat Jalan B) — ⬜ **persetujuannya sendiri masih ditunggu**

---

### J-11 · Tautan antar generasi endorsement = **`PRODKE`**

**Keputusan work owner:** yang menautkan generasi endorsement adalah **`PRODKE`**.

⛔ **`PRODKE` sendiri tidak cukup — ia butuh sandaran.** `PRODKE` hanyalah nomor urut (1, 2, 3…) yang
**berulang di polis berbeda**; tanpa kunci pengelompokan, "generasi ke-3" tidak menunjuk polis mana
pun. Jadi aturannya: **kelompokkan menurut polis, urutkan menurut `PRODKE`.**

`[terverifikasi]` Kandidat kunci pengelompokan pada rancangan flat, lembar `Kolom` (1.330 baris):

| Tabel · kolom | Terisi pada 115 contoh |
| --- | ---: |
| `T_GENERAL_POLIS.POLICY_NO` | **115 / 115** |
| `T_GENERAL_POLIS.POLICY_MASTER_NUMBER` | **1 / 115** |
| `T_GENERAL_POLIS.POLICY_MASTER_ID_PEGA` | **1 / 115** |

⛔ **`POLICY_MASTER_NUMBER` BUKAN kunci itu.** Namanya menggoda — "master" terdengar seperti polis
induk lintas generasi — tetapi `[terverifikasi]` ia terisi hanya pada **1 dari 115** contoh, dan pada
contoh itu nilainya **berbeda** dari `PolicyNo`. Artinya ia sesuatu yang lain; **`belum
terverifikasi`** apa. Ini jebakan **"nama bukan bukti"** (`CLAUDE.md` §3 butir 3) yang hampir
termakan.

**Aturan yang ditetapkan:** kelompokkan dengan **`POLICY_NO`**, urutkan dengan **`PRODKE`**.

#### ⚠️ Batas pengetahuan yang jujur — aturan ini **tidak dapat diuji dari korpus**

`[terverifikasi]` Dari **115** berkas contoh terdapat **113 `PolicyNo` unik**; hanya **2** nilai
muncul lebih dari sekali, masing-masing **tepat 2 kali**. Angka itu **persis** dua pasang duplikat
byte-identik yang sudah dikenal (`NB-172576`=`NB-176005`, `NB-181622`=`NB-184183`, V-45).

⛔ **Jadi nol pasangan generasi sejati ada di contoh.** Tidak satu pun kasus di mana dua generasi dari
polis yang sama hadir bersamaan. Aturan pengelompokan di atas bersandar pada **jawaban work owner +
bentuk skema**, **bukan** pada pengukuran. ⚠️ Jangan tulis `[terverifikasi]` di atasnya.

📌 Temuan sampingan: angka **113** ini **menguatkan secara mandiri** populasi efektif 113 (V-45),
lewat jalur yang sama sekali berbeda dari perbandingan byte.

```powershell
# menghasilkan 115 / 113 / 2 / 2  - nilai TIDAK dicetak (K-025)
$h=@{}; $n=0
foreach ($p in [IO.Directory]::EnumerateFiles('D:\migrasi\RNM\DDL\CONTOH','*.*')) {
  $n++; $m=[regex]::Match([IO.File]::ReadAllText($p),'(?i)<PolicyNo>([^<]*)</PolicyNo>')
  if ($m.Success) { $v=$m.Groups[1].Value.Trim(); if ($v -ne '') { if(-not $h.ContainsKey($v)){$h[$v]=0}; $h[$v]++ } } }
$ber=0; $maks=0
foreach ($k in $h.Keys) { if ($h[$k] -gt 1) { $ber++ }; if ($h[$k] -gt $maks) { $maks=$h[$k] } }
"$n / $($h.Count) / $ber / $maks"
```

---

### J-12 · **Syarat Jalan B** — enam, seluruhnya dapat diperiksa

Pertanyaan work owner: *"apa saja itu syaratnya"*. Berikut daftarnya. ⬜ **Persetujuan Jalan B
sendiri belum diberikan** — ini penjelasan, bukan keputusan.

| # | Syarat | Bila dilanggar |
| :-: | --- | --- |
| **1** | **Fase 2 selesai sebelum kasus pertama dibuat sungguhan** di sistem baru | syarat pokok — Jalan B batal |
| **2** | ⚠️ **BUKAN syarat Jalan B — syarat K-070 itu sendiri.** Sumber lama masih dapat **dibaca** saat fase 2 berjalan: Oracle lama dan `JSON_POLIS` **belum dimatikan atau dicabut aksesnya** | ⛔ generasi lama **tidak dapat dimuat lewat jalan mana pun**. Jalan A pun tidak menolong |
| **3** | **Data yang lahir saat pengujian dianggap habis pakai** — muat ulang mengosongkan tabel, apa pun yang diketik penguji ikut terhapus | perlu daftar tabel yang dikecualikan, dan pengecualian itu mengembalikan masalah penjodohan |
| **4** | **Fase 2 dimuat dalam SATU jalan** — bukan bertahap per COB atau per lini | `ROW_UID` lahir dari pandangan sebagian → masalah penjodohan kembali |
| **5** | **Satu titik "beku" yang diumumkan** — sesudah fase 2 lolos rekonsiliasi, `ROW_UID` dibekukan dan tidak ada muat-ulang lagi | tanpa titik tegas, "boleh muat ulang" menjalar ke masa produksi |
| **6** | **Rekonsiliasi ADR-0001 dijalankan SESUDAH fase 2**, atas populasi lengkap | rekonsiliasi atas fase 1 menguji populasi yang bukan sasaran akhir |

⛔ **Bila salah satu tidak terpenuhi → kembali ke Jalan A** (menjodohkan ke belakang). Jalan A tetap
bisa dikerjakan; biayanya adalah **salah jodoh yang tidak bersuara** — tidak ada constraint basis data
yang akan mengeluh, persis seperti 12 tabel berinduk ganda.

📌 **Syarat 2 adalah yang paling mudah dilanggar tanpa sadar**, karena mematikan Pega terasa seperti
akhir pekerjaan, padahal fase 2 justru membutuhkannya tetap hidup.

#### ⚠️ Koreksi penempatan — syarat 2 salah letak pada perumusan pertama

Syarat 2 sempat ditulis seolah **syarat Jalan B**. **Itu keliru dan menyesatkan.** Ia syarat
**K-070**: begitu keputusan "fase 2 menyusul" diambil, sumber lama **wajib** tetap terbaca sampai
fase 2 selesai — **apa pun jalan yang dipilih**. Jalan A membutuhkannya persis sama.

**Arti "data hilang" di baris itu, tepatnya:** bukan terhapus, melainkan **tidak lagi dapat
dijangkau**. Tiga hal berbeda yang mudah tertukar:

| Yang dimaksud | Bukan yang dimaksud |
| --- | --- |
| Oracle lama **dimatikan / aksesnya dicabut** sebelum fase 2 → generasi lama ada secara fisik tetapi tak terbaca | ⛔ baris `JSON_POLIS` **dihapus** — tidak ada rencana apa pun yang menghapusnya |
| | ⛔ isi tabel flat terhapus saat muat ulang — itu **syarat 3**, dan yang terhapus hanya **data uji** |

#### Perbandingan risiko — ⛔ Jalan B **lebih rendah**, dan satu-satunya yang **dapat dibatalkan**

| | Jalan A | **Jalan B** |
| --- | --- | --- |
| Risiko khasnya | ⛔ **salah jodoh `ROW_UID` yang tidak bersuara** — permanen, tidak ada constraint yang mengeluh | data **uji** terhapus saat muat ulang (syarat 3) · disiplin titik beku (syarat 5) |
| Bila rencananya meleset | tidak ada jalan mundur — `ROW_UID` sudah beku di atas data produksi | **turun ke Jalan A**, tanpa biaya tambahan |
| Sifat keputusannya | ⛔ **pintu satu arah** | **dapat dibatalkan kapan saja** sampai sistem dipakai sungguhan |

📌 **Memilih B tetap menyimpan A sebagai cadangan; memilih A membuang B.** Itu alasan pokok
rekomendasinya, lebih kuat daripada perbandingan kerumitannya.

#### ⚠️ Satu risiko K-070 yang belum pernah dinamai — berlaku bagi kedua jalan

Selama fase 1, sistem diuji atas data yang **tidak lengkap**. Perilaku yang bergantung pada riwayat
— jalur banding lewat `HISTORYAKSEPTASIPEGA`, penomoran endorsement berikutnya — akan berlaku
berbeda saat pengujian dibanding saat dipakai sungguhan. ⛔ **Pengujian yang lolos di fase 1 karena
itu bukan jaminan.** Ini **bukan** alasan menolak K-070; ia alasan menjalankan **fase 2 sedini
mungkin**, bukan selambat mungkin.

---

### J-13 · `STS_KONVERSI` = status **konversi ke tim lain**, ⛔ **bukan** penanda "sudah jadi"

**Keterangan work owner:** kolom itu dipakai untuk **konversi ke tim lain** — data yang diinput akan
dipakai tim lain untuk keperluan mereka.

⛔ **Konsekuensi yang mengoreksi rencana saya.** Pada K-071 konsekuensi 4 saya mencatat `STS_KONVERSI`
sebagai kandidat penyaring *"endorsement sudah jadi"*. **Itu keliru** — ia menandai apakah data sudah
**diserahkan ke hilir**, bukan apakah endorsementnya selesai. Memakainya sebagai penyaring akan
menyaring hal yang salah.

**Dua akibat:**

1. **Penyaring "sudah jadi" masih belum teridentifikasi.** ⚠️ Bobotnya **turun drastis** bila Jalan B
   disetujui: fase 2 memuat ulang seluruh generasi, sehingga ketepatan pemilihan di fase 1 tidak lagi
   menentukan hasil akhir. ⛔ Tetapi **tidak hilang** — fase 1 tetap dipakai untuk menguji, dan menguji
   dengan generasi yang salah menghasilkan pengujian yang menyesatkan.
2. **Jalur konversi produksi tetap wajib ada di sistem baru.** Karena tim lain memakai keluarannya,
   penulisan ke tabel produksi bukan warisan yang boleh dilepas. Ini **sejalan** dengan lingkup yang
   sudah ada (`services/production`, `_MENJADI`/`_SELISIH` dipertahankan, `CLAUDE.md` §4.3) — dicatat
   di sini karena alasannya baru: **ada konsumen di luar tim ini**.

---

### J-14 · `_SELISIH` tanpa `_MENJADI` menandai **penambahan lokasi**

**Keterangan work owner:** keenam kolom itu berarti **ada penambahan lokasi**.

**Ini menjelaskan bentuknya.** Bila sebuah lokasi **ditambahkan**, tidak ada nilai "sebelum" yang
dapat "menjadi" nilai baru — yang ada hanya **pertambahannya**. Karena itu kolomnya `_SELISIH` saja,
tanpa pasangan `_MENJADI`.

⚠️ `[keterangan work owner]`, bukan `[terverifikasi]` dari korpus. Korpus hanya membuktikan
**bentuknya** (6 dari 15 tanpa pasangan, 2 di antaranya tanpa kolom polos); **artinya** datang dari
work owner.

**Konsekuensi:**

1. **Keenamnya diport apa adanya sebagai kolom tunggal.** Instruksi `CLAUDE.md` §4.3 *"dipertahankan
   berpasangan"* berlaku bagi **9 basis yang memang berpasangan**; keenam sisanya **bukan pasangan
   yang hilang** melainkan **jenis kolom yang berbeda**. ⛔ Jangan "melengkapi" dengan membuat kolom
   `_MENJADI` baru — itu termasuk perbaikan diam-diam (§1).
2. **Menguatkan kesimpulan J-9.** Kolom-kolom itu tidak dapat memberi nilai sebelum endorsement bukan
   karena datanya hilang, melainkan karena **memang tidak pernah ada nilai sebelumnya**.

---

### `[pertanyaan terbuka]` yang TIDAK ditutup keputusan ini

1. **Persetujuan Jalan B** (J-12) — syaratnya kini tertulis; keputusannya belum diambil.
2. **Penyaring "sudah jadi"** untuk fase 1 — `STS_KONVERSI` terbukti bukan jawabannya.
3. **Arti `POLICY_MASTER_NUMBER` / `POLICY_MASTER_ID_PEGA`** — terisi 1 dari 115, berbeda dari
   `PolicyNo`, artinya `belum terverifikasi`.

---

## K-073 · **Jalan B DISETUJUI** — fase 2 mengosongkan dan memuat ulang, tidak menjodohkan ke belakang

**Tanggal:** 25 September 2026 · **Menutup:** **J-15** (usulan J-6 pada K-071, syarat pada K-072)

**Keputusan work owner:** **setuju Jalan B.**

| | |
| --- | --- |
| **Fase 1** | muat **versi terakhir** tiap polis · `ROW_UID` dibangkitkan · ⛔ **bersifat sementara** |
| **Fase 2** | ⛔ **kosongkan tabel flat**, muat **seluruh generasi dalam satu jalan** · `ROW_UID` lahir sekali atas riwayat lengkap |
| **Titik beku** | sesudah fase 2 lolos rekonsiliasi → `ROW_UID` beku, tidak ada muat ulang lagi |

**Enam syaratnya berlaku penuh** sebagaimana tertulis di K-072 J-12, dengan koreksi yang sudah
dicatat di sana: **syarat 2 adalah syarat K-070, bukan syarat Jalan B.**

### Konsekuensi yang mengikat

1. ⛔ **`ROW_UID` fase 1 TIDAK BOLEH dijadikan sandaran apa pun di luar tabel flat.** Tidak
   diekspor, tidak dikirim ke tim lain, tidak dicetak di layar sebagai identitas, tidak dipakai
   sebagai kunci integrasi. Ia akan **berganti** saat fase 2. Melanggar ini membuat muat ulang
   menjadi perubahan yang terlihat keluar — dan Jalan B kehilangan sifat dapat-dibatalkannya.
2. **Bobot V-42 turun drastis.** Rekonstruksi UID lintas versi — yang
   `BAHAN-SPEC-PEMUATAN.md` §5 sebut bagian tersulit — **tidak lagi diperlukan**, karena fase 2
   memuat riwayat lengkap dari nol. Yang tersisa hanya **membangkitkan** UID, bukan
   **merekonstruksi**-nya.
3. **Ketepatan pemilihan "versi terakhir" di fase 1 turun bobotnya, tidak hilang.** Hasil akhirnya
   ditentukan fase 2. Tetapi fase 1 dipakai **menguji**, dan menguji atas generasi yang salah
   menghasilkan pengujian yang menyesatkan → **J-16 tetap terbuka**.
4. **Rekonsiliasi ADR-0001 dijalankan sesudah fase 2**, atas populasi lengkap. Nol selisih sampai
   digit terakhir **tidak dilonggarkan**.
5. **Jalan A tetap tersedia sebagai cadangan** sampai detik sebelum sistem dipakai sungguhan. Bila
   salah satu syarat gugur, turun ke Jalan A tanpa biaya tambahan.

### Risiko yang diterima sadar

⚠️ Selama fase 1 sistem diuji atas data **tidak lengkap**; perilaku yang bergantung riwayat (jalur
banding `HISTORYAKSEPTASIPEGA`, penomoran endorsement) berlaku berbeda dibanding nanti. ⛔ **Lolos uji
di fase 1 bukan jaminan.** Peredamnya satu: **jalankan fase 2 sedini mungkin.**

### `[pertanyaan terbuka]` yang TIDAK ditutup keputusan ini

1. **J-16** — penyaring "endorsement sudah jadi" untuk fase 1.
2. **J-17** — arti `POLICY_MASTER_NUMBER` / `POLICY_MASTER_ID_PEGA` (terisi 1 dari 115).

---

## K-074 · **Seam keempat disetujui** — `loader.Flatten`, murni, tanpa basis data

**Tanggal:** 25 September 2026 · **Mengubah:** batas "tiga seam" pada spec 03 ·
**Melahirkan:** `04-spec\11-spec-pemuatan-data-lama.md`

**Pertanyaan asal:** program pemuatan data lama butuh titik uji. Tidak satu pun dari tiga seam
runtime yang ada cocok — `rules.Eval`, `acceptance.Next`, dan `premium.Calculate` seluruhnya soal
**perhitungan**, sedangkan pemuatan soal **pemindahan bentuk**.

**Keputusan work owner:** **`loader.Flatten` sebagai seam baru** — satu dokumen penawaran masuk,
himpunan baris untuk 78 tabel keluar. ⛔ **Murni**: tidak menyentuh basis data, tidak menyentuh jam,
tidak menyentuh berkas.

Penulisan ke Oracle diuji lewat **seam `repository`**, yang ⛔ **bukan seam baru** — ia **sudah
dicadangkan** `04-spec\03-spec-modul-terverifikasi.md` §5 butir 5, dan di sana tertulis ia *"menunggu
butir 3"*, yaitu bentuk tabel flat. Bentuk itu kini ada (78 tabel, DDL draf lolos delapan
pemeriksaan), sehingga cadangan itu **aktif**, bukan ditambahkan.

### Hitungan seam sesudah keputusan ini

| Seam | Modul | Status |
| --- | --- | --- |
| 1 | `internal/rules.Eval` | lama |
| 2 | `services/acceptance.Next` | lama |
| 3 | `services/premium.Calculate` | lama |
| **4** | **`loader.Flatten`** | ⭐ **baru** |
| — | seam `repository` | **dicadangkan sejak 16 September**, kini aktif |

⛔ **Batas "tidak ada seam lain tanpa keputusan baru" tetap berlaku.** Keputusan ini **adalah**
keputusan baru yang dimaksud kalimat itu; ia **tidak** membuka pintu bagi seam berikutnya.

### Mengapa murni, dan bukan satu pintu yang sekaligus menulis

Dua pilihan lain ditolak, keduanya dengan alasan yang sama:

| Pilihan yang ditolak | Sebab |
| --- | --- |
| Satu seam `loader.Load` yang membaca **dan** menulis | test wajib memakai penulis tiruan, dan kebenaran pemetaan 78 tabel bercampur dengan kebenaran penulisan ke Oracle |
| Tanpa seam baru — uji lewat `repository` saja | kesalahan pemetaan pohon baru terlihat setelah menyentuh basis data |

⛔ **Alasan pokoknya dua belas tabel berinduk ganda.** Keutuhannya **tidak dijaga constraint apa
pun** — satu kolom penunjuk induk menunjuk sampai **5** tabel berbeda. Salah mengisi nama tabel
induk **tidak akan bersuara** di basis data. Memurnikan `Flatten` membuat justru bagian itu dapat
diuji tanpa basis data dan tanpa tiruan.

### Konsekuensi

1. **Test seam 4 berbentuk table-driven**, mengikuti bentuk ketiga seam terdahulu — spec 03 menegaskan
   ketiganya **menjadi** prior art bagi modul berikutnya.
2. **Kontrak `panic` saat menulis mata uang tidak diketahui** (ADR-0006, spec 03 §5 butir 5) akhirnya
   punya tempat: **seam `repository`**. Ia berhenti menjadi "kewajiban yang belum punya tempat".
3. ⛔ **Satu kekurangan cakupan dicatat, tidak ditutupi:** tautan antar generasi endorsement
   `[terverifikasi]` tidak dapat diuji dari korpus — nol pasangan generasi sejati di antara 115
   contoh (113 nomor polis unik; dua yang berulang adalah duplikat byte-identik V-45). Test untuknya
   memakai dokumen yang disusun sendiri dan berstatus `[dugaan]`.

---

*Keputusan dicatat di sini agar asal-usulnya jelas dan tidak tampak sebagai temuan korpus
(`CLAUDE.md` §3). Keputusan yang dibalik **tidak dihapus** — ia ditandai dibatalkan dan menunjuk
penggantinya, supaya alasan perubahan tetap terbaca. Menambah keputusan: tambahkan entri `K-0nn`
dengan tanggal, pertanyaan asal, dan konsekuensinya.*
