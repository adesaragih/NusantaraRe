# KOREKSI-2-DIJALANKAN — P1 terjawab, penahan proyek habis

> **Untuk:** pemilik pekerjaan migrasi Treaty Inward.
>
> **Sebab:** naskah **empat** stored procedure dan **tiga contoh** `DATA_JSON` diterima work owner
> pada 22 September 2026.
>
> ⚠️ **Brief ronde ini direvisi di tengah jalan** — dari dua contoh menjadi tiga, ditambah satu
> temuan yang menyentuh keputusan uang. Laporan ini adalah keadaan **sesudah** revisi itu
> dikerjakan seluruhnya.
>
> **Lingkup tulis:** hanya **dokumen hidup** di dalam `OUTPUT_HASIL_RNM\`. ⛔ Korpus tidak disentuh.
> ⛔ `grilling-ronde-*.md`, `VERIFIKASI-*.md`, `KOREKSI-P18-DIJALANKAN.md`, dan `PROMPT-*.md`
> **tidak disunting** — pernyataan basi di sana sengaja dibiarkan, sesuai Bab 0 brief.
>
> ⚠️ **Dua hal yang perlu keputusan Anda ada di Bab 3.3 dan Bab 4 uji 4.**

---

## Bab 1 — Daftar suntingan

**Dua belas berkas disunting.** Tidak satu pun di luar kelompok "dokumen hidup".

### 1.1 Pencabutan "P1 menahan"

⚠️ **Cacah ulang saya berbeda dari perkiraan brief, dan berkasnya sebagian lain.** Brief menduga
20 pernyataan di 12 berkas, dengan tiket **12** dan **13** di antaranya. Terukur **dua cara**:

| Cara | Hasil |
| --- | ---: |
| **A** — baris ber-`P1` **dan** berkata *menahan / tertahan / menunggu / blocked* | **28** baris di 5 kelompok |
| **B** — seluruh baris ber-`P1`, lalu disaring manual | **110** baris |
| ⭐ **di dokumen HIDUP saja** *(yang boleh disunting)* | ⭐ **29** baris di **11** berkas |

⛔ **Tiket 12 dan 13 tidak memuat satu pun rujukan P1** — brief memperkirakan 2 dan 3. Yang justru
memuatnya dan tidak disebut brief: `KEADAAN-NB-TREATY-IN.md`, `PERTANYAAN-SIAP-KIRIM`,
`PERTANYAAN-untuk-pemilik-export-Pega.md`, dan `spec.md` §9.1.

| Berkas | Yang lama | Yang baru |
| --- | --- | --- |
| `PERTANYAAN-untuk-DBA.md` | *"**P1** dan **P29** menahan seluruh penulisan spesifikasi penyimpanan data."* — di kepala **dan** ekor | ⭐⭐ *"P1 dan P29 SUDAH TERJAWAB, dan keduanya tidak lagi menahan apa pun. **Penahan proyek kini nol.**"* + RALAT dengan bunyi lama dikutip |
| `spec.md` §9.1 | judul *"**Satu** yang MENAHAN"*, satu baris `P1` | ⭐⭐ *"Yang MENAHAN — **NIHIL**"*; baris P1 dikutip utuh di dalam RALAT |
| `spec.md` §1.3 | *"Butir `[terbuka]` **20** terbagi: ⛔ **1 MENAHAN** (P1) · 19 tidak"* | *"Butir `[terbuka]` **23** terbagi: ⭐ **0 MENAHAN** · 23 tidak menahan"* |
| `spec.md` §1.3 sensus | `butir [terbuka]` **20 / 20** | **23 / 23**, dengan **seluruh** angka lama dikutip *(21 → 20 → 22 → 23)* |
| `spec.md` Cara B | *"9.1 + 9.2 = **1 + 19**"* | *"**0 + 23**"* |
| `spec.md` §9.2 | 19 butir | ditambah **4**: **P61** · **P62** · **sensus properti** · ⛔⛔ **galat angka uang tersimpan** |
| `issues/15-migrasi-paritas…` | baris `P29 \| ⛔ MENAHAN migrasi` | ✅ **P29 TERJAWAB**; ⚠️ ditambah bab baru **"Ambang paritas uang"** — lihat Bab 3.4 |
| `KEADAAN-NB-TREATY-IN.md` | *"⭐ **Tinggal satu yang menahan.** Suratnya sudah siap: `SURAT-P1-KE-DBA.md`."* | ⭐⭐ *"**Penahan proyek NIHIL.**"* + RALAT; ditambah catatan bahwa yang tersisa **bukan penahan**, melainkan sensus properti — pekerjaan tim migrasi sendiri |
| `KEADAAN` tabel §6 | baris `P1 \| naskah tiga stored procedure \| DBA` | dicoret, ✅ **TERJAWAB — empat naskah** |
| `PERTANYAAN-YANG-MASIH-KOSONG.md` | judul *"**Dua** yang MENAHAN"*; *"Untuk DBA · **6** kosong"*; judul P1 dan P29 bertanda ⛔⛔ **MENAHAN**; sebaran DBA **6** | ⭐ *"Yang MENAHAN — **NIHIL**"*; **5 kosong**; P1 dan P29 bertanda ✅ **TERJAWAB**; sebaran DBA **4** |
| `PERTANYAAN-untuk-pengembang-Pega-lama.md` | *"`[terbuka]` **P1 menjadi lebih penting** … permintaan P1 berdiri utuh untuk **ketiga** procedure."* | ✅ **TERJAWAB — dan bukan tiga, melainkan EMPAT**; bunyi lama dikutip utuh |
| `PERTANYAAN-YANG-MASIH-KOSONG.md` | blok yang sama *(salinan)* | disamakan, kata demi kata |
| `PERTANYAAN-RONDE-4.md` | *"satu hitungan … ditumpangkan ke surat **P1**"* | ⚠️ catatan: **P1 terjawab, tetapi hitungan itu TIDAK ikut terjawab** — tetap `[terbuka]` sebagai **P49** |
| `PERTANYAAN-SIAP-KIRIM-NB-TREATY-IN.md` | judul *"**P1** — Kami butuh isi tiga program…"* | dicoret, ✅ **TERJAWAB 2026-09-22** |

### 1.2 Tiket `00-skema-penyimpanan-dan-migrasi`

| | Lama | Baru |
| --- | --- | --- |
| judul | *"tiket gantung, menunggu bahan"* | *"menunggu **SENSUS PROPERTI**, bukan lagi bahan dari DBA"* |
| `Blocked by` | `P29` *(contoh isi JSON polis)* · `P1` *(naskah tiga stored procedure)* | ⭐ **sensus properti `PolicyTreatyIn`** — ronde tersendiri, pekerjaan tim migrasi |
| Blocker | dua baris `[DBA]` | keduanya ✅ **DITERIMA**, diganti satu baris sensus + penjelasan **kenapa contoh tidak cukup** |
| badan | *"tidak dapat dirancang ulang **tanpa melihat contoh isinya**"* | ⚠️ RALAT — ⛔ **contoh saja ternyata tidak cukup, dan itu terbukti justru sesudah contohnya datang** |
| "Yang sudah diketahui" | satu paragraf | ⭐ tabel empat baris: data kontrak **sudah relasional** *(20 kolom)* · **tujuh dari delapan** kolom `json_polis` sudah datar · **hanya `DATA_JSON`** yang perlu dipecah · sisi baca lewat view 39 kolom |

⭐ **Statusnya tetap `needs-info`** — tetapi yang ditunggu berubah dari **bahan pihak luar** menjadi
**pekerjaan tim migrasi sendiri**. Itu perubahan sifat, bukan sekadar perubahan daftar.

---

## Bab 2 — P51 dan P46 ditajamkan

### 2.1 P51 — penghapusan ternyata **bersyarat**

⚠️ **Penajaman, bukan pembatalan.** Keputusan work owner *"(a) Benar — data produksi HARUS dihapus
bila konversi gagal"* **tetap berlaku utuh**. Yang bertambah adalah penjaga yang **tidak terlihat
dari sisi Pega**, sebab ia ada di dalam naskah procedure.

```
hanya menghapus bila   STS_KONVERSI IS NULL   atau   STS_KONVERSI <> '1'
lingkup selalu         WHERE IDPEGA = <satu kasus>
```

| `Bisnis` | Tabel yang dihapus |
| --- | --- |
| `'T'` — treaty | `JSON_POLIS` · `TREATYINPRODUCTION` · `TREATYINPRODUCTION_BACKUP` |
| `'F'` — fakultatif | `JSON_POLIS` · `FACINPRODUCTION` · `FACINPRODUCTION_BACKUP` · `FACOUTPRODUCTION` |

⭐⭐ **Yang baru diketahui: data yang sudah berhasil dikonversi tidak dapat tersentuh.**
⭐ Dan pertanyaan lama *"berapa baris yang terhapus"* kini terjawab: baris milik **satu kasus**,
pada **tiga tabel** untuk treaty dan **empat** untuk fakultatif — **nol baris** bila kasus itu sudah
bertanda konversi berhasil.

**Yang mengikat penulisan ulang di Go:** penjaga `STS_KONVERSI` **wajib ikut ditulis**.
Menghilangkannya membuat kegagalan berulang dapat menghapus konversi yang sudah sah.
⚠️ `[penyimpangan sadar]` yang sudah tercatat **tetap berlaku** — logikanya ditulis ulang, bukan
memanggil procedure.

⭐ **Dua butir `[terbuka]` di dalam P51 ikut tertutup**, dan keduanya tertutup oleh jawaban P1 —
bukan oleh saya:

| Butir | Keadaan |
| --- | --- |
| *"`PEGA_DELETE_ERROR_KONVERSI` adalah procedure KEEMPAT yang belum diminta"* | ✅ **naskahnya diterima** |
| *"Blok pemanggilnya tidak memuat `COMMIT` — perlu dipastikan apakah procedure meng-commit di dalam"* | ✅ **meng-commit DI DALAM**, sesudah seluruh penghapusan berhasil; utuh atau tidak sama sekali |

### 2.2 P46 — potongan 2,2 % kini terverifikasi dari data produksi

Rumusnya semula hanya terbaca dari XML. Ia kini terbukti dari satu kasus produksi:

```
TypeTax                = "Inclusive"
Deduction1             = 32.516
BrokerageFeeSebenarnya = 31.8160        <- nilai yang benar-benar tersimpan

32.516 / 1.022 = 31.81604696673189...
                 dibulatkan 4 desimal = 31.8160   ✓ COCOK
```

⭐ **Hitungannya saya ulang sendiri dengan `Decimal`**, bukan disalin dari brief: pembagi
`102.2/100`, pembulatan `ROUND_HALF_UP`. Cocok tepat.

⭐⭐ **Bonus yang tidak disebut brief:** nilai tersimpan berhenti di **4 desimal**, bukan 8.
Artinya kasus ini melewati **cabang empat desimal**
*(`InputPolicyTreatyInDetail_NonProp` / `InputPolicyTreatyOutDetail_NonProp`)* — pada cabang delapan
desimal nilainya akan `31.81604697`.

⛔ **Butir `[terbuka]` presisi TETAP TERBUKA.** Satu kasus menunjukkan **cabang mana yang
dilewatinya**, bukan **cabang mana yang seharusnya**. Pertanyaan asalnya masih milik `[Finance]`
dan `[Product+Underwriting]`.

---

## Bab 3 — Butir `[terbuka]` baru

### 3.1 Dua yang diperintahkan brief — dicatat, **tidak ditutup**

Keduanya di `PERTANYAAN-untuk-DBA.md`, lengkap dengan konteks, bentuk jawaban yang diharapkan,
dampak bila salah, dan baris `rujukan:` — sama seperti pertanyaan lain.

| # | Butir | Pemutus |
| --- | --- | --- |
| **P61** | Tabel `_BACKUP` **ikut terhapus** bersama data utamanya. Namanya menyiratkan jaring pengaman, tetapi sesudah penghapusan **tidak ada apa pun yang tersisa untuk dipulihkan** | `[work owner]` |
| **P62** | `IF TRUNC(v_now) <= TO_DATE('02/01/2026') THEN v_mm_yyyy := '12.2025'` — seluruh penomoran sebelum 2 Januari 2026 dipaksa ke periode **Desember 2025**. Masih dikehendaki? Perlu lagi saat pindah ke 2027? | `[work owner]` · `[Finance]` |

⭐ **Keduanya ditandai TIDAK MENAHAN** — perlu dipastikan sebelum **jalur terkait dibangun ulang**,
bukan sebelum pekerjaan dimulai.

### 3.2 ⭐⭐ Temuan yang mengubah rancangan — bentuk JSON **mengikuti jenis treaty**

| | Contoh 1 | Contoh 2 | Contoh 3 |
| --- | ---: | ---: | ---: |
| jenis | proporsional | proporsional *(SOA)* | ⛔ **non-proporsional XOL** |
| medan skalar tingkat atas | **37** | **64** | **47** |
| ⭐ ada di **ketiganya** | | **23** | |
| ⭐ gabungan ketiganya | | **74** | |

⛔ **Hanya 23 dari 74 medan muncul di ketiga contoh** — lebih dari dua pertiga gabungan tidak ada
di sebagian kasus.

⚠️ **Selisih kecil yang saya periksa:** catatan lama untuk contoh 1 menulis **38** medan; selisihnya
`pxObjClass`, properti internal Pega. Tanpa itu **37**, dan angka brief cocok.

### ⛔⛔ Yang paling berbahaya: satu nama daftar, DUA kedalaman

**`ListInstallment` bersarang di contoh 3 tetapi datar di contoh 2.**

⭐⭐ **Dan ini dapat saya periksa dari korpus, bukan hanya dari contoh** — disapu atas NB + EDM,
blok penyunting dibuang:

| Bentuk yang dirujuk aturan | Rujukan | Berkas |
| --- | ---: | ---: |
| ⛔ `ListInstallment(n).InstallmentList` — **bersarang** | **58** | 6 |
| ⛔ `ListInstallment(n).<medan skalar>` — **datar** | **101** | 11 |
| `TreatyXOLList(n).ValueList` | 207 | 8 |
| `OldData.TreatyXOLList` | 130 | 9 |
| `BreakDownSpreadList` | 11 | 5 |

⇒ ⭐ **Kedua bentuk sama-sama nyata di dalam aturan sistem lama.** Contoh 2 dan contoh 3 bukan
kebetulan — keduanya sah, dan **pengurai yang menganggapnya seragam akan patah.**

⭐ **Penentunya terlihat:** `QuotationData.ProportionalType` *(korpus: **170** kemunculan)* dan
`IsNewPolicyNonProp` *(**87**)*. ⭐ Contoh 3 juga memunculkan **`OldData` di dalam dokumen polis
baru**, berisi `TreatyXOLList` kosong — didukung korpus, 130 rujukan di 9 berkas.

### Sensus awal properti — dua cara

| Cara | Hasil |
| --- | ---: |
| **A** — sapuan rujukan teks `PolicyTreatyIn.<nama>` | **94** nama unik, **1.869** rujukan |
| **B** — hanya sisi-kiri penetapan nilai | **86** nama unik, **867** penulisan |
| irisan · hanya di A · hanya di B | 86 · 8 · **0** |

⭐ **9 simpul bersarang** — cocok dengan angka brief. Empat di antaranya — `OldData`,
`TreatyXOLList`, `TreatyXOLDifferenceList`, `TreatyDifference` — **tidak muncul di ketiga contoh**
padahal dirujuk ratusan kali.

⚠️ **Angka rujukan saya berbeda dari brief, dan saya tidak memaksakannya cocok:**

| Simpul | brief | saya: `PolicyTreatyIn.<n>` | saya: `<n>` di mana pun |
| --- | ---: | ---: | ---: |
| `OldData` | 258 | **224** | 344 |
| `TreatyXOLList` | 228 | **186** | 347 |
| `TreatyXOLDifferenceList` | 123 | **99** | 108 |
| `TreatyDifference` | 97 | **96** | 192 |

⇒ **Tiga pola, tiga hasil.** ⛔ **Ketiganya batas bawah, bukan sensus.** Properti yang hanya ada
di dalam dokumen JSON dan **tidak pernah dirujuk satu aturan pun** tidak tertangkap cara mana pun.

⛔⛔ **Dicatat sebagai `[terbuka]` pada P29 dan sebagai butir 22 di `spec.md` §9.2:** rancangan tabel
**tidak dapat** dibangun dari contoh.

### 3.3 ⛔⛔ Galat angka uang yang sudah tersimpan permanen — dan satu koreksi atas sebabnya

```
premium angsuran   148157378.220000069     x 4  =  592629512.880000276
NetPremium         592629512.880000276          <- sama persis
nilai bulat 2 desimal                           =  592629512.88
selisih                                         =  2,76 x 10^-7
```

⭐ **Hitungannya saya ulang dengan `Decimal` presisi 40**, bukan disalin: perkaliannya **cocok
persis**, selisihnya **tepat 2,76 x 10^-7**.

⚠️ **Tetapi dugaan sebabnya saya koreksi — dan koreksinya MEMPERBURUK, bukan meringankan.**
Galat ini **bukan** sekadar galat `float64`:

| | Nilai |
| --- | --- |
| `float64` terdekat dari `592629512.88` | `592629512.87999999523...` — meleset **ke bawah**, ~**4,8 x 10^-9** |
| galat yang benar-benar tersimpan | meleset **ke ATAS**, **2,76 x 10^-7** |
| ⭐ nisbahnya | ~**57 kali lebih besar**, **arah berlawanan** |
| `592629512.88 / 4` dalam desimal | `148157378.22` — **tepat, tanpa sisa** |

⇒ ⛔ **Pembagian empat yang bersih tidak akan menghasilkan ekor ini sama sekali.** Galatnya lahir
di **rantai perhitungan**, bukan di penyimpanan. ⛔ Sebab pastinya **tidak dapat dinyatakan dari
satu kasus, dan tidak dikarang.**

⭐⭐ **Kenapa koreksi ini penting:** bila galatnya ada di rantai hitung, sistem baru yang memakai
aritmetika desimal yang benar akan menghasilkan angka **yang berbeda dari sistem lama** — bukan
hanya pada data lama, tetapi pada **perhitungan baru**. ⛔ **Uji paritas yang menuntut kesamaan
persis akan GAGAL justru karena sistem barunya benar.** Itu dicatat di tiket **15**, bab baru *"Ambang paritas uang"*.

**Tiga pilihan tercatat di P29 — ⛔ tidak dipilih:**

| | Pilihan | Akibat | Ambang paritas |
| --- | --- | --- | --- |
| **a** | ikuti apa adanya | cacat diwariskan, tiap penjumlahan menambahnya | paritas **persis** |
| **b** | bulatkan saat migrasi | **angka historis berubah** | paritas **tidak berlaku** untuk data lama |
| **c** | simpan apa adanya, bulatkan saat dihitung | sejarah utuh, galat **berhenti berlipat** | paritas **dengan toleransi** |

⛔⛔ **Pemutusnya `[work owner]`; presisi wajar per mata uang milik `[Finance]`.**
⭐ **ADR-0003 tetap ditegakkan apa pun pilihannya** — sistem baru tidak memakai `float`, sehingga
**tidak menambah galat baru**. ⛔ **Tidak ada ADR baru yang diperlukan.**

### 3.4 ⚠️ Keputusan yang saya ambil sendiri

**P61 dan P62 saya beri nomor.** Brief hanya menyuruh *"catat pada lembar pemilik yang tepat"*,
tanpa menyebut penomoran. ⭐ Saya menomorinya karena aturan proyek sendiri berbunyi *"satu nomor
berarti satu pertanyaan di seluruh proyek"* dan *"nomor dipakai untuk melacak jawaban ke temuan
aslinya"* — butir tanpa nomor tidak akan tertangkap cacah pertanyaan mana pun.
⚠️ **Bila Anda lebih suka keduanya tanpa nomor, keduanya dapat diturunkan tanpa merusak apa pun** —
belum ada berkas lain yang merujuknya.

**Dan keduanya saya letakkan di lembar `[DBA]`**, walau pemutusnya `[work owner]`, sebab tidak ada
lembar work owner dan faktanya berasal dari basis data. Tercatat jelas di tiap butir siapa
pemutusnya.

**Dan satu berkas saya sunting di luar daftar brief: `issues/15-migrasi-paritas-dan-jejak-keputusan.md`.**
⭐ Sebabnya: galat angka uang **menentukan ambang uji paritas**, dan tiket 15-lah yang memiliki uji
itu. Tanpa catatan di sana, uji paritas dapat ditulis dengan ambang yang dikarang, lalu selisihnya
diam-diam dianggap wajar. ⛔ **Yang saya tambahkan hanya catatan — nol butir ditutup, nol ambang
ditetapkan**, dan baris `P29` di tabelnya ditandai terjawab. Dapat dibatalkan tanpa merusak apa pun.

---

## Bab 4 — Hasil enam pemeriksaan

| # | Uji | Diharapkan brief | Terukur | |
| ---: | --- | --- | --- | :---: |
| 1 | cacah pertanyaan NB, dua cara | **49 dari 49**, nol terbuka | ⭐ **P1–P49: 49 nomor, nol lompatan, nol kosong** — 46 terjawab + 3 ditarik *(P18, P41, P48)*. Cara A *(judul bab)* = 51, Cara B *(blok Jawaban)* = 51, ⭐ **sepakat** — 51 karena **P61 dan P62 baru** | ✅ |
| 2 | cacah pertanyaan EDM | 11 dari 11 | **11 dari 11 terjawab**, nol kosong | ✅ |
| 3 | sinkronisasi salinan | nol nomor berbeda isi | **nol** nomor yang hanya ada di salinan. `KOSONG` 19 nomor, `SIAP-KIRIM` 38 — keduanya subset sah dari 51 | ✅ |
| 4 | tiket `blocked` | hanya **13**, karena **P30** | ⭐ **nol tiket blocked karena P1.** Tiket **13** `blocked` karena **P30** dan **P8** — sesuai. ⚠️ Tetapi **dua lagi `needs-info`**: **00** *(sensus properti — pekerjaan kita sendiri)* dan **05** *(pemetaan nama → peran, `[IAM]`)* | ⚠️ |
| 5 | "P1 menahan" di dokumen hidup | **nihil** | ⭐ **0** — disapu ulang sesudah menyunting. ⚠️ Penyaring menandai **4** baris, keempatnya di **laporan ini sendiri** *(judul bab "Pencabutan 'P1 menahan'" dan barisnya)* — teks laporan, bukan pernyataan hidup | ✅ |
| 6 | butir `[terbuka]` aktif di `spec.md`, dua cara | dicacah ulang; angka lama **dikutip** | Cara A *(baris bernomor §9.2)* = **23**, Cara B *(blok sensus §1.3)* = **23**, ⭐ **sepakat**, nol lompatan. §9.1 = **NIHIL**. Seluruh angka lama dikutip: **21 → 20 → 22 → 23** | ✅ |

⚠️ **Uji 4 tidak sepenuhnya seperti yang diharapkan brief, dan itu bukan kesalahan suntingan.**
Brief menulis *"hanya tiket 13"*. Terukur **tiga** tiket masih berstatus tertahan — tetapi ⭐ **nol
di antaranya karena P1**, dan sifat ketiganya berbeda:

| Tiket | Status | Menunggu | Sifat |
| --- | --- | --- | --- |
| **13** | `blocked` | **P30** · **P8** | ⛔ pihak luar — sesuai harapan brief |
| **00** | `needs-info` | **sensus properti** | ⭐ **pekerjaan tim migrasi sendiri**, bukan pihak luar |
| **05** | `needs-info` | pemetaan nama → peran | ⛔ `[IAM]` · `[work owner]` — tidak disebut brief, tidak saya sentuh |

⭐ **Butir 22 di `spec.md` §9.2 sengaja saya tulis sebagai butir `[terbuka]`** supaya sensus properti
itu tidak hilang dari neraca hanya karena ia pekerjaan kita sendiri.

---

## Bab 5 — TELEMETRI EKSEKUSI

### Yang tidak dikerjakan, dan sebabnya

| ⛔ Larangan brief | Dipatuhi? | Bukti |
| --- | :---: | --- |
| jangan menyunting berkas pada tabel Bab 0 | ✅ | `grilling-ronde-*`, `VERIFIKASI-KEADAAN`, `VERIFIKASI-P18`, `KOREKSI-P18-DIJALANKAN`, `PROMPT-*` — **nol** tersentuh; pernyataan basi di sana **dibiarkan**, sesuai perintah |
| jangan menutup butir `[terbuka]` yang tidak disebut | ✅ | ditutup: **P1**, **P29** *(keduanya disebut brief)* dan **dua butir di dalam P51** *(ditutup oleh jawaban P1, bukan oleh saya)*. ⛔ Butir presisi P46 **tetap terbuka** walau data produksi menajamkannya. ⛔⛔ Butir **galat angka uang** dicatat dengan **tiga pilihan dan akibatnya**, **tanpa dipilih** |
| jangan merancang tabel | ✅ | **nol** rancangan. Yang ditulis justru **kenapa rancangan belum dapat dibuat** |
| jangan membuat ADR baru | ✅ | **0**. ⚠️ Dan tidak ada yang perlu — penjaga `STS_KONVERSI` adalah rincian logika, bukan keputusan arsitektur |
| nol kode, DDL, nama orang | ✅ | tiga baris pseudo-syarat dan satu baris `IF` dikutip sebagai **bunyi naskah yang dipertanyakan**, bukan sebagai kode yang akan dibangun |
| nol contoh JSON tersimpan | ✅ | ⛔ **Nol.** Dari **tiga** contoh hanya **lima angka medan uang** yang dikutip *(32.516 · 31.8160 · TypeTax · 148157378.220000069 · 592629512.880000276)*; nol nama orang, nol nomor polis, nol nama tertanggung |

### Cara tiap angka diperoleh

| Angka | Cara A | Cara B | Sepakat? |
| --- | --- | --- | --- |
| pernyataan "P1 menahan" | baris ber-`P1` **dan** kata menahan | seluruh baris ber-`P1`, disaring manual | ⚠️ 28 lawan 110 — **disaring manual**, 29 di dokumen hidup |
| pertanyaan NB | judul bab di 6 lembar | blok `**Jawaban:**` per bagian | ✅ 51 = 51 |
| properti `PolicyTreatyIn` | sapuan rujukan teks | sisi-kiri penetapan nilai | ⚠️ 94 lawan 86, irisan 86, **hanya di B = 0** |
| butir `[terbuka]` spec | baris bernomor §9.2 | angka di blok sensus §1.3 | ✅ 22 = 22 |
| uji silang P46 | `Decimal` presisi penuh | pembulatan `ROUND_HALF_UP` 4 desimal | ✅ cocok tepat |

### Keraguan dan selisih — dicatat, bukan dirapikan

| # | Hal | Keterangan |
| ---: | --- | --- |
| 1 | Brief memperkirakan tiket **12** dan **13** memuat rujukan P1 | ⛔ **Keduanya tidak memuat satu pun.** Yang memuat justru empat berkas yang tidak disebut brief. Brief sendiri menyuruh mencacah ulang |
| 2 | Brief menulis **97** properti; saya ukur **94** | Selisih bentuk pola. ⭐ **Keduanya batas bawah**, dan keduanya bukan sensus |
| 3 | Brief menulis contoh 1 punya **37** medan; catatan lama menulis **38** | ⭐ Selisihnya `pxObjClass`. Brief benar untuk medan bisnis |
| 4 | Brief mengharapkan **hanya tiket 13** yang tertahan | Terukur **tiga**; nol karena P1. Dilaporkan apa adanya, tidak dipaksakan cocok |
| 5 | Uji 1 kini membaca **51** nomor, bukan 49 | ⭐ Karena **P61 dan P62 yang saya buat sendiri**. P1–P49 tetap **49/49 nol kosong** |

### Ongkos

| Ukuran | Nilai |
| --- | ---: |
| panggilan model | **86** |
| token keluar | **195.444** |
| token cache ditulis | **1.959.121** |
| token cache dibaca | **33.342.231** |
| berkas disunting | **12** |
| berkas keluaran baru | **1** — berkas ini |
| berkas korpus disunting | **0** |
| butir `[terbuka]` ditutup | **4** — P1 · P29 · dua di dalam P51 |
| butir `[terbuka]` baru | **4** — P61 · P62 · sensus properti · ⛔⛔ **galat angka uang** |
| ADR baru · rancangan tabel | **0** · **0** |
| keputusan uang yang saya ambil sendiri | ⛔ **0** — tiga pilihan dicatat, **tidak dipilih** |

⚠️ **Telemetri TIDAK diukur dari luar.** `claude --print --output-format json` akan menjalankan
sesi **baru dan terpisah**, yang ongkosnya bukan ongkos ronde ini. ⭐ Angka di atas **hasil
pengurangan terhadap baseline** yang diambil dari catatan sesi pada awal ronde. ⚠️ Angka token
sejati tidak terlihat dari dalam sesi.

### Keadaan sesudah ronde ini

| | Sebelum | Sesudah |
| --- | ---: | ---: |
| pertanyaan NB P1–P49 terjawab | 48 / 49 | ⭐ **49 / 49** |
| pertanyaan NB masih kosong | **1** *(P1)* | ⭐ **0** |
| pertanyaan EDM | 11 / 11 | 11 / 11 |
| ⭐ pertanyaan baru | — | **2** *(P61, P62)* — tidak menahan |
| ⛔ **butir yang MENAHAN pekerjaan** | **1** *(P1)* | ⭐⭐ **0** |
| butir `[terbuka]` di `spec.md` | 20 | **23** |
| tiket tertahan pihak luar | 3 | **2** *(13, 05)* |
| tiket menunggu pekerjaan sendiri | 0 | **1** *(00 — sensus properti)* |
| surat siap kirim | 1 | ⭐ **0** |

⭐⭐ **Penahan proyek habis. Yang tersisa adalah pekerjaan kita sendiri: sensus properti
`PolicyTreatyIn`, yang menentukan bentuk tabel penyimpanan — dan itu ronde berikutnya.**
