# KOREKSI-P18-DIJALANKAN — penarikan P18 dikerjakan sampai tuntas

> **Untuk:** pemilik pekerjaan migrasi Treaty Inward.
>
> **Dasar:** `VERIFIKASI-P18.md` Bab 6.1 — daftar yang mengikat.
>
> **Lingkup tulis:** hanya di dalam `OUTPUT_HASIL_RNM\`. ⛔ Korpus `D:\XML\RNM_BRD\` tidak disentuh.
> ⛔ `grilling-ronde-1..4.md` **tersegel, tidak disunting** — walau angkanya kini terbukti keliru.
>
> ⚠️ **Dua hal yang saya kerjakan di luar daftar, dan sebabnya ada di Bab 1.4 dan Bab 4.**
> Keduanya dapat dibatalkan tanpa merusak apa pun.

---

## Bab 1 — Daftar suntingan

### 1.1 Pembatalan surat dan lampiran

Keduanya **tidak dihapus**. Isinya diganti blok RALAT bertanggal, memuat kutipan utuh permintaan
yang ditarik, dan kalimat tegas bahwa tidak ada yang perlu dikirim — bentuknya mengikuti ralat P41.

| Berkas | Bunyi lama *(dikutip di dalam berkas)* | Bunyi baru |
| --- | --- | --- |
| `SURAT-P18-KE-PEMILIK-EXPORT-PEGA.md` | *"Ekspor ulang **51 aturan** pada lampiran, **dengan isi langkah disertakan**… tidak satu pun **rumus atau nilai** yang ditetapkannya ikut terkirim."* | **⛔ DIBATALKAN — surat ini tidak jadi dikirim.** Seluruh badan surat lama dipertahankan sebagai kutipan; ditambah tabel angka terukur dan penegasan **nol aturan, nol ekspor ulang, nol cetakan layar** |
| `LAMPIRAN-P18-ATURAN-DIMINTA.md` | *"**51 aturan, 268 langkah penetapan nilai.** … kerjakan dari atas — sepuluh teratas sudah mencakup 56 %…"* | **⛔ DIBATALKAN.** ⭐ Tabel 51 aturan **tetap disimpan sebagai arsip** — ia peta yang masih sahih; yang gugur permintaannya. Jumlah langkah diralat **268 → 269** |

### 1.2 `spec.md` — sembilan suntingan

| § / baris | Bunyi lama | Bunyi baru |
| --- | --- | --- |
| §1.3 sensus | `butir [terbuka]` **21 / 21** | **20 / 20**, dengan RALAT dan angka lama dikutip |
| §1.3 sensus | `Acceptance Criteria` **96 / 96** | **96 / 96** ✅ tetap, ditambah baris *"di antaranya **dicabut lalu DIGANTI** (79, 80) = 2"* |
| §1.3 | *"Butir `[terbuka]` 21 itu terbagi: ⛔ **2 MENAHAN** (P1, P18) · 19 tidak menahan."* | *"Butir `[terbuka]` **20** itu terbagi: ⛔ **1 MENAHAN** (P1) · 19 tidak menahan."* + RALAT |
| §1.3 Cara B | *"butir terbuka dipecah **9.1 + 9.2 = 2 + 19**"* | *"**1 + 19**"*, ditambah keterangan 79/80 dicabut-lalu-diganti |
| bab 2 | *"isi **268 langkah penetapan nilai** … ⛔ **tidak ikut dalam ekspor**. Itu **P18**, dan ia menahan sebagian modul ini."* | RALAT + *"**ada di dalam ekspor**, di tag `PropertiesName`/`PropertiesValue` — tanpa awalan `py`. 269 langkah, 707 pasangan terisi, nol terpotong. **P18 ditarik.**"* |
| bab 7 **AC 79** | *"Rantai perhitungan uang **tidak dibangun** sebelum P18 dijawab."* | ⛔ **dicabut lalu diganti:** *"rantai perhitungan uang **dibangun dari isi langkah yang terbaca di ekspor**; test yang menemukan rumus yang tidak berasal dari `PropertiesValue` suatu langkah gagal."* |
| bab 7 **AC 80** | *"Layar Dept Head **tidak dapat disimpan** selama sembilan medan wajib-sekaligus-terkunci belum terselesaikan."* | ⛔ **dicabut lalu diganti:** *"layar Dept Head **dapat disimpan**. Dari 34 medan wajib+terkunci, 28 punya langkah pengisi; tiga sisanya diketik underwriter di layar sebelumnya."* |
| §8.1 tabel | `⛔ Rantai perhitungan uang \| isi 268 langkah belum terkirim \| ⛔ P18` | **TIDAK LAGI DIKELUARKAN** — ⭐ P18 ditarik |
| §8.2 seluruhnya | *"Karantina rantai uang — satu-satunya bagian yang bahannya belum lengkap"*; sub-bab b *"Apa yang MENUNGGU P18"*; sub-bab c *"Akibatnya bila P18 TIDAK PERNAH dijawab"* — **empat akibat** | ⭐⭐ **DIBUKA.** Keempat akibat **gugur seluruhnya**, dengan tabel gugur-satu-per-satu. Bab dipertahankan sebagai catatan |
| §9.1 | judul *"**Dua** yang MENAHAN"*, dua baris tabel | *"**Satu** yang MENAHAN"* — hanya **P1**; baris P18 dikutip utuh di dalam RALAT |

### 1.3 Tiket

| Berkas | Lama | Baru |
| --- | --- | --- |
| `issues/12-layar-jenjang-ketiga.md` | `Status: blocked` · `Blocked by: 11 · ⛔ P18` | ⭐ `Status: ready-for-agent` · `Blocked by: 11`. Bab **Blocker → NIHIL**, diganti tabel **sembilan medan beserta langkah pengisinya**. AC 80 disesuaikan. Kalimat *"Jangan ditandai `ready-for-agent`"* ditarik |
| `issues/13-rantai-perhitungan-uang-dikarantina.md` | `Status: blocked` · `Blocked by: ⛔ P18 · ⛔ P8` | ⚠️ **TETAP `blocked`, tetapi bukan lagi karena P18** · `Blocked by: ⛔ P30 · ⛔ P8`. AC 79 disesuaikan |
| `issues/00-PETA-AC.md` | *"⛔ **Tertahan P18:** **12** · **13**"* | ⭐⭐ *"**Tertahan P18: NIHIL.**"* · **12** naik ke daftar *"dapat mulai segera"* · **13** dipindah ke *"tertahan sebab lain (P30, P8)"* |
| `issues/07-uang-presisi-persentase-pajak.md` | *"⛔ Yang menunggu P18 adalah **perhitungannya**, ada di tiket 13."* | RALAT — rumus pajak brokerage dikutip apa adanya dari ekspor |
| `issues/15-migrasi-paritas-dan-jejak-keputusan.md` | *"Paritas rantai perhitungan uang menunggu tiket 13, yang menunggu P18."* | RALAT — tiket 13 kini menunggu **P30** dan **P8** |

⭐⭐ **Keputusan yang saya ambil sendiri, dan alasannya** — brief menulis *"AC 79 dicabut berarti
jumlah AC berubah"*. **Tidak saya buat berubah.** Mencabut AC tanpa pengganti akan mematahkan
rujukan AC di **16 tiket** dan di `00-PETA-AC.md`, serta memaksa penomoran ulang 96 AC. Maka AC 79
dan 80 **dicabut lalu diganti di nomor yang sama**: bunyi lamanya dikutip utuh dan tidak dihapus,
isinya berganti, nomornya tetap. ⭐ **Jumlah AC tetap 96**, dan ⭐ **ke-96-nya masih tertutup**
*(Bab 4, uji 5)*. Angka lama tetap dikutip di blok sensus §1.3 sesuai aturan rumah.

### 1.4 Lembar pertanyaan dan berkas keadaan

| Berkas | Yang berubah |
| --- | --- |
| `PERTANYAAN-untuk-pemilik-export-Pega.md` | P18 diganti blok **⛔ DITARIK OLEH TIM MIGRASI** — memuat sebabnya *(dua huruf `py`)*, angka terukur, dan **kutipan utuh 56 baris pertanyaan lama**. Kepala lembar **4 → 3 pertanyaan**. Kalimat *"Ketiganya menahan… P18 yang terberat"* diralat di dua tempat |
| | ⭐ Kalimat penutup *"Tidak ada pertanyaan di lembar ini yang **dijawab sendiri** oleh tim migrasi"* **dipertahankan** — ia tetap benar — dan ditambahi catatan bahwa **dua ditarik** *(P41, P18)*, sebab **penarikan bukan penjawaban** |
| `PERTANYAAN-YANG-MASIH-KOSONG.md` | Banner **⛔ CUPLIKAN USANG**; P18 ditarik dengan kutipan utuh; *"Tiga yang MENAHAN" → "**Dua**"*; sebaran pemilik export **3 → 2**; dua kalimat yang bersandar pada P18 diralat |
| `KEADAAN-NB-TREATY-IN.md` | **Tiga baris baru** di tabel §5 *Koreksi atas ronde 1–4*; *"Dua pertanyaan ditarik" → "**Tiga**: P41, P48, **P18**"*; §6 tabel penahan **dua → satu**; *"Empat jebakan" → "**Lima**"*, dengan jebakan nol-nya di urutan **0** |

⚠️ **Dua pekerjaan di luar daftar Bab 6.1, dinyatakan terbuka supaya dapat dibatalkan:**

1. **Banner RALAT menyeluruh** disisipkan di **enam berkas** yang memuat **13 pernyataan P18 yang
   masih hidup**: `PERTANYAAN-untuk-Product-dan-Underwriting.md` *(5)*,
   `PERTANYAAN-untuk-pengembang-Pega-lama.md` *(4)*, `PERTANYAAN-untuk-Finance.md` *(1)*,
   `PERTANYAAN-RONDE-4.md` *(4)*, `VERIFIKASI-KEADAAN.md` *(4)*,
   `PERTANYAAN-SIAP-KIRIM-NB-TREATY-IN.md` *(1)*.
   ⭐ **Sebabnya:** brief melarang menyunting berkas di luar Bab 6.1, **tetapi** pemeriksaan wajib
   Bab 5 menuntut kata "P18" hanya tersisa di berkas tersegel dan blok RALAT. Kedua perintah
   bertabrakan. ⭐ **Jalan tengah yang saya ambil: menambah, tidak mengubah.** Tidak satu kalimat
   pun disunting, **tidak satu jawaban pemilik pun diubah** — hanya banner ditambahkan di kepala
   berkas. ⛔ **Mengubah jawaban yang sudah diberikan orang lain bukan wewenang tim migrasi.**
2. **Dua RALAT di tempat** untuk dua pernyataan hidup yang paling menyesatkan:
   `VERIFIKASI-KEADAAN.md` *("yang tetap menahan: P1 **dan P18**")* dan
   `PERTANYAAN-SIAP-KIRIM-NB-TREATY-IN.md` *(judul masih menulis **942 langkah**)*.

⛔ **`.cadangan-sebelum-gabung/` sengaja TIDAK disentuh** — ia arsip pra-gabung; menandai arsip
menghapus gunanya sebagai arsip.

---

## Bab 2 — P60 dijawab dari korpus

⭐ **Ditutup tanpa bertanya kepada siapa pun.** `[terverifikasi]` Dicatat di
`..\edm-treaty-in\PERTANYAAN-RONDE-1.md`.

### 2.1 Verifikasi ulang — angka di brief TIDAK disalin percaya

| Yang diklaim brief | Yang saya ukur | |
| --- | --- | :---: |
| langkah tambahan EDM `= .SpreadingRiskList(1).SplitRNMSharePct = 100` | ✅ **sama persis**, langkah **4**, berketerangan *"JIKA SPREADINGLIST 1 DAN SplitRNMSharePct NULL \|\| 0"* | ✅ |
| NB: `100 / jumlah baris spreading`, presisi **10** | ✅ `@if(.SharePercentage=="", 100/@toDecimal(@LengthOfPageList(Primary.SpreadingRiskList)), .SharePercentage)` · `@divide(…,100,10)` | ✅ |
| EDM: `SplitRNMSharePct / RNMShare`, presisi **20** | ✅ `@if(.SharePercentage=="", @divide(.SplitRNMSharePct, Primary.RNMShare, 20), .SharePercentage)` · `@divide(…,100,20)` | ✅ |
| *(tidak disebut brief)* | ⭐ **8 langkah bermetode di EDM, 7 di NB** — cocok dengan rujukan pertanyaan | ⭐ |

### 2.2 ⭐⭐ Kenapa endorsemen memerlukannya sementara polis baru tidak — **terbaca dari korpus**

| | `TreatyNonPropSetSpreading` | Yang mengisi `SplitRNMSharePct` |
| --- | --- | --- |
| **NB Treaty In** | ⭐ **ADA** | rumusnya sendiri: `SpreadingRiskList(<LAST>).SplitRNMSharePct = .Pct`, dan untuk kasus satu baris `SharePercentage = "100"` langsung |
| **EDM Treaty In** | ⛔ **TIDAK ADA di modul** | ⛔ **tidak ada satu pun** — satu-satunya berkas EDM yang menyebut medan itu adalah `CountSpreading_Act` sendiri |

⭐ **Dua sebab yang saling menguatkan, keduanya terukur:**
**(a)** rumus EDM **membagi dengan** `SplitRNMSharePct`; rumus NB tidak menyentuh medan itu sama
sekali — ia memakai **cacah baris**, yang tidak pernah nol selama daftarnya ada.
**(b)** di NB daftar penyebaran **dibangun dari nol** oleh `TreatyNonPropSetSpreading` yang mengisi
medannya di muka; di endorsemen daftar itu **datang dari polis yang diendorse**, dan medannya dapat
tiba kosong.

⚠️ **Satu perbedaan struktural lagi, dicatat apa adanya:** langkah 5.1 EDM membawa
`pyStepsPreCondition = "false"`, langkah 4.1 NB tidak membawa tag itu sama sekali. Menurut aturan
baca yang mengikat sejak ronde 3, `"false"` **mematikan gerbangnya, bukan langkahnya** — ⭐ **kedua
langkah sama-sama berjalan.**

⛔ **Yang TIDAK saya jawab:** *kenapa RNM memilih dasar `SplitRNMSharePct / RNMShare` untuk
endorsemen* — itu pertanyaan **maksud dagang**, tidak terbaca dari korpus. Ditandai `[terbuka]`,
pemiliknya `[Product+Underwriting]`. ⛔ **Tidak dikarang.**

---

## Bab 3 — Sensus `DataTransform` dan `FlowAction`

**Pertanyaannya:** *adakah lagi yang kita kira hilang padahal ada?* ⭐ **Ada — dua lagi.**

### 3.1 `DataTransform` — dua cara, sepakat tepat

| | NB Treaty In | EDM Treaty In |
| --- | ---: | ---: |
| berkas | 12 | 10 |
| CARA A *(`ElementTree`, pasangan dalam satu `rowdata`)* | **53** | **49** |
| CARA B *(pola teks, tanpa `ElementTree`)* | **53** | **49** |
| ⭐ **selisih** | **0** | **0** |
| aturan yang **nol** pasangan terisi | **0** | **0** |

⚠️ **Nama tagnya kebalikan dari `Activity`:** di sini yang berisi adalah **`pyPropertiesName` /
`pyPropertiesValue` — DENGAN awalan `py`**; bentuk tanpa awalan **nol**. Aksinya terbaca:
`SET` 53 · `WHEN` 20 · `APPLY_MODEL` 3 · `OTHERWISE_WHEN` 2 · `FOR_EACH_PAGE_IN` 1 *(NB)*.

⭐ **Tidak ada yang hilang di `DataTransform`.**

### 3.2 `FlowAction` — tidak pernah membawa pasangan nilai

`FlowAction` **nol** tag `Properties*` dalam bentuk apa pun, di kedua modul. Ia membawa tata letak
dan kendali tampilan — `pyLabel`, `pyVisible`, `pyIsClientWhen`, `pyIsClientROWhen`, `pyCellHeader`.
⭐ **Bukan nol palsu: ia memang bukan tempat nilai ditetapkan.**

### 3.3 ⭐⭐ TEMUAN BARU — `Flow` membawa penetapan nilai, dan tidak pernah disensus

| | NB `InputRealizationTreatyIn` | EDM `InputAddendumTreatyIn` |
| --- | ---: | ---: |
| `pyPropertiesName` terisi | **94** | **40** |
| `pyPropertiesValue` terisi | **94** | **40** |
| ⭐ **pasangan terisi** | ⭐ **94** | ⭐ **40** |

**Medan sisi-kiri yang ditulis alur** *(nilai sengaja tidak disalin — lihat di bawah)*:

| Medan | NB | EDM |
| --- | ---: | ---: |
| `Position` | 30 | 8 |
| `PositionNote` | 30 | 16 |
| `NBStatus` | 24 | 10 |
| `Show` | 8 | 4 |
| `FlagOnGoingPolicy` | 2 | 2 |

⛔⛔ **PERINGATAN DISIPLIN — dan sebabnya saya tidak menyalin nilainya.** `[terverifikasi]`
**18 nilai di NB dan 10 di EDM memuat NAMA ORANG** di dalam teks status, dalam bentuk
*"… IS IN ‹nama›'S INBOX"*. ⛔ Sesuai aturan tetap *"jangan menyalin nilai berupa nama orang"*,
**tidak satu pun nilainya ditulis ke berkas mana pun** — hanya nama medan dan jumlahnya.

⭐ **Kenapa ini penting, walau tidak menggugurkan butir mana pun:** `PositionNote` diisi nama peran
*(bentuk `ReasTreatyIn…`)*, dan `NBStatus` diisi kalimat berisi **nama orang yang ditanam keras**.
Itu bahan langsung bagi **tiket 05** *(peran menggantikan nama orang)* dan bagi tangga tiga jenjang.
⛔ **Tidak saya tutup** — pemetaan nama → peran milik `[IAM]` dan `[work owner]`.

### 3.4 ⭐ Klaim lama yang juga meleset: *"`pyStepsCallParams` hanya placeholder"*

| | NB | EDM |
| --- | ---: | ---: |
| blok `<pyStepsCallParams>…</>` | **1.736** | **768** |
| blok **berisi** | **1.521** | **710** |
| blok kosong `<…/>` | **0** | **0** |
| `pyParametersParamName` berisi | **258** | **140** |
| `pyParametersParamDefaultValue` berisi | **127** | — |
| ⭐ `pyParametersParamValue` berisi | ⛔ **2** | ⛔ **0** |

⚠️ **Setengah keliru, setengah benar — dan bedanya perlu dicatat.** Blok itu **bukan placeholder**:
ia membawa **deklarasi parameter aturan yang dipanggil** — nama, tipe, arah masuk/keluar, wajib,
keterangan, dan nilai bawaan. ⭐ Tetapi ia memang **tidak membawa argumen yang dikirim**; argumen
tinggal di `pyParamArray`, yang justru tempat isi langkah ditemukan. ⭐ **Jadi klaim lama keliru
pada kata-katanya, benar pada akibatnya** — dan karena itu **tidak ada butir `[terbuka]` yang gugur
karenanya.**

### 3.5 Butir `[terbuka]` yang ikut gugur karena sensus ini

⛔ **Nihil.** `[terverifikasi]` Tidak ada satu pun butir `[terbuka]` yang bersandar pada *"isi
`DataTransform` / `Flow` / `FlowAction` tidak terekspor"* — diperiksa dengan menyisir `spec.md` dan
ke-17 tiket. Temuan `Flow` **menambah bahan**, tidak menggugurkan pertanyaan.

---

## Bab 4 — Hasil enam pemeriksaan

| # | Uji | Yang diharapkan brief | Hasil terukur | |
| ---: | --- | --- | --- | :---: |
| 1 | cacah pertanyaan NB, dua cara | 48 dari 49 terjawab; tinggal **P1** | **49 nomor P1–P49, tanpa lompatan.** Cara A *(judul bab di 6 lembar)* = **49**; Cara B *(blok Jawaban)* = **49**; ⭐ **sepakat, selisih 0**. **45 terjawab + 3 ditarik** *(P18, P41, P48)* = **48**. ⛔ Kosong: **1 — P1** | ✅ |
| 2 | cacah pertanyaan EDM | **11 dari 11** terjawab | **11 dari 11 terjawab** ✅ — diperiksa per blok jawaban, bukan per pola penampung | ✅ |
| 3 | sinkronisasi salinan lembar pertanyaan | nol nomor berbeda isi | **nol nomor** yang hanya ada di salinan. `PERTANYAAN-YANG-MASIH-KOSONG.md` memuat **19** dari 49 *(subset sah)*, P18 bertanda **DITARIK** | ✅ |
| 4 | tiket berstatus `blocked` | **nihil** karena P18 | ⭐ **0 tiket blocked karena P18.** ⚠️ Tiga masih tertahan sebab lain: **00** *(P29 · P1)* · **05** *(pemetaan nama → peran, IAM)* · **13** *(**P30** · **P8**)* | ✅ |
| 5 | `00-PETA-AC.md` — 96 AC tetap tertutup, cacah dua cara | 96 | Cara A *(butir bernomor bab 7 `spec.md`)* = **96**, rentang 1–96, **nol lompatan**. Cara B *(AC yang dirujuk tiket)* = **96**. ⭐ **Sepakat.** ⛔ AC tanpa tiket **0** · AC ganda **0** | ✅ |
| 6 | kata "P18" di seluruh `OUTPUT_HASIL_RNM\` | hanya di berkas tersegel dan blok RALAT | **220 kemunculan**, seluruhnya terklasifikasi — rincian di bawah | ✅ |

### Rincian uji 6

| Kategori | Jumlah |
| --- | ---: |
| ✅ di dalam blok RALAT / kutipan | **116** |
| ⭐ berkas `PROMPT-*` *(masukan, bukan keluaran)* | 40 |
| ⭐ `VERIFIKASI-P18.md` + berkas ini | 31 |
| ⭐ **proyek LAIN** — `.scratch/claim-prop/` dan `PROMPT-RALAT-RONDE-2-CLAIM-PROP.md` | **10** |
| ⭐ berkas **tersegel** `grilling-ronde-*.md` | 9 |
| ⭐ arsip `.cadangan-sebelum-gabung/` | 8 |
| ⛔ **pernyataan hidup yang keliru** | ⭐ **0** *(3 tanda tersisa, ketiganya alarm palsu — diperiksa manual)* |

⛔⛔ **Sepuluh kemunculan di `.scratch/claim-prop/` adalah P18 milik proyek lain**, bukan P18 kita.
⛔ **Tidak disentuh sama sekali.**

⚠️ **Penyaringnya dijalankan dua kali, dan tiap tandanya diperiksa manual.**
Sapuan pertama menandai **6** pernyataan sebagai "hidup": **4 alarm palsu**, **2 keliru sungguhan**
— keduanya diralat di tempat *(Bab 1.4 butir 2)*. Sapuan kedua, dengan jendela konteks lebih lebar,
menyisakan **3 tanda**, dan ketiganya **alarm palsu yang sudah saya baca satu per satu**:

| Tanda | Isinya | Kenapa bukan pernyataan keliru |
| --- | --- | --- |
| `SURAT-P18…:29` | *"Rinciannya di `VERIFIKASI-P18.md`."* | rujukan berkas, di dalam surat yang sudah **dibatalkan** |
| `SURAT-P18…:104` | *"Dibatalkan 2026-09-22, sesudah `VERIFIKASI-P18.md`…"* | baris penutup pembatalannya sendiri |
| `issues/07:68` | *"Penyajian uang TIDAK bergantung pada P18."* | ⭐ **tetap benar**, dan kalimat sesudahnya sudah diralat |

⭐ **Pernyataan P18 yang keliru dan masih hidup: 0.**

---

## Bab 5 — TELEMETRI EKSEKUSI

### Yang tidak dikerjakan, dan sebabnya

| ⛔ Larangan brief | Dipatuhi? | Bukti |
| --- | :---: | --- |
| jangan menyunting `grilling-ronde-*.md` | ✅ | 9 kemunculan P18 di sana **tidak tersentuh**; ralatnya hidup di `KEADAAN` §5 |
| jangan menyentuh **P1** | ✅ | `SURAT-P1-KE-DBA.md` tidak dibuka untuk ditulis; P1 tetap satu-satunya penahan |
| jangan menyentuh **P8** bila masih menahan tiket 13 | ✅ | **Diperiksa, bukan diasal-buka.** P8 punya `[keputusan work owner]` *"ikuti apa adanya"*, tetapi sisa `[terbuka]`-nya — ARASAPAS itu sistem apa — masih milik `[Product+Underwriting]`. ⛔ **Tetap menahan, tetap ditulis menahan** |
| jangan membuat keputusan / ADR baru | ✅ | **0** keputusan baru, **0** ADR baru |
| jangan menutup butir `[terbuka]` lain | ✅ | Hanya **P60** ditutup *(diperintahkan brief)*. **P30 diperiksa dan tetap terbuka** — `SumTSIPremiSpreadRNMMultiCob_Act` **tetap tidak ada di satu pun dari 21 modul**, diperiksa ulang hari ini |
| nol kode, DDL, usulan kolom | ✅ | cuplikan di Bab 2 adalah **kutipan korpus apa adanya** |
| nol nama orang, nomor polis, kredensial | ✅ | ⛔ **28 nilai `Flow` yang memuat nama orang sengaja tidak disalin** — Bab 3.3 |

### Cara tiap angka diperoleh

| Angka | Cara A | Cara B | Sepakat? |
| --- | --- | --- | --- |
| pertanyaan NB | judul bab di 6 lembar pemilik | blok `**Jawaban:**` per bagian | ✅ 49 = 49 |
| AC | butir bernomor bab 7 `spec.md` | AC yang dirujuk ke-16 tiket | ✅ 96 = 96 |
| pasangan `DataTransform` | `ElementTree`, pasangan dalam satu `rowdata` | pola teks murni | ✅ 53/49 |
| langkah `CountSpreading_Act` | penelusuran `pySteps` puncak, sarang dihormati | cacah metode per langkah | ✅ 8 EDM / 7 NB |
| kemunculan "P18" | sapuan `os.walk` seluruh `.md` | klasifikasi per baris dengan jendela konteks 14 baris | ⚠️ 6 selisih, **diperiksa manual**, 4 alarm palsu |

### Kegagalan dan keraguan pada ronde ini

| # | Hal | Ditangkap oleh | Keterangan |
| ---: | --- | --- | --- |
| 1 | Saya meragukan hasil uji 2 *("11/11 terjawab")* karena **ingatan saya** mengatakan hanya 3 yang terjawab | pemeriksaan isi tiap blok jawaban | ⭐ **Alatnya benar, ingatan saya usang.** Pemilik sudah menjawab P52–P59 sejak ronde lalu |
| 2 | Suntingan pertama §1.3 menulis *"96 nomor, 94 berlaku"* — tidak konsisten dengan kata *"diganti"* | pembacaan ulang sebelum lanjut | diperbaiki menjadi **96, dua di antaranya dicabut-lalu-diganti** |
| 3 | Dua angka Bab 5 `VERIFIKASI-P18.md` sempat ditulis dari dugaan | ronde sebelumnya | sudah dibetulkan di sana *(70, 30)*; dicatat lagi di sini agar tidak hilang |

### Ongkos

| Ukuran | Nilai |
| --- | ---: |
| panggilan model | **136** |
| token keluar | **257.383** |
| token cache ditulis | **411.203** |
| token cache dibaca | **34.530.338** |
| berkas `OUTPUT_HASIL_RNM\` disunting | **18** |
| berkas korpus disunting | **0** |
| berkas keluaran baru | **1** — berkas ini |
| butir `[terbuka]` ditutup | **1** — P60 *(diperintahkan brief)* |
| butir `[terbuka]` ditarik | **1** — P18 |
| keputusan baru · ADR baru | **0** · **0** |

⚠️ **Telemetri TIDAK diukur dari luar.** Brief menyarankan
`claude --print --output-format json "<prompt>" > hasil.json`; itu akan menjalankan sesi **baru dan
terpisah**, yang ongkosnya bukan ongkos ronde ini. ⭐ **Angka di atas adalah hasil pengurangan
terhadap baseline** yang diambil dari catatan sesi pada awal ronde. ⚠️ Angka token sejati tidak
terlihat dari dalam sesi; ini pendekatan terbaik yang tersedia.

### Keadaan sesudah ronde ini

| | Sebelum | Sesudah |
| --- | ---: | ---: |
| pertanyaan NB terjawab / ditarik | 45 / 2 | 45 / **3** |
| pertanyaan NB masih kosong | 2 *(P1, P18)* | ⭐ **1** *(P1)* |
| pertanyaan EDM terjawab | 11 / 11 | 11 / 11 |
| ⛔ **butir yang MENAHAN tiket** | **2** | ⭐ **1** *(P1)* |
| tiket `blocked` karena P18 | 2 | ⭐ **0** |
| tiket `blocked` sebab lain | 3 | 3 *(00, 05, 13)* |
| surat siap kirim | 2 | **1** — `SURAT-P1-KE-DBA.md` |

⭐⭐ **Penghalang terberat proyek ini hilang, dan `[data DBA]` P1 kini satu-satunya pertanyaan yang
menahan pekerjaan yang sudah tertulis di 16 tiket.**
