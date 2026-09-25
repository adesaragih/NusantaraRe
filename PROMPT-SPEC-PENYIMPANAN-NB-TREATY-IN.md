# PROMPT — spec penyimpanan relasional, NB Treaty In

> Salin seluruh isi berkas ini sebagai prompt ke sesi eksekutor.
>
> ✅ **Prasyarat terpenuhi.** Rancangan tabelnya sudah selesai dan seluruh butir penahannya
> ditutup work owner pada 22–23 September 2026.
>
> ⚠️ **Skill `to-spec` tidak dapat dipanggil sendiri oleh agen** (CLAUDE.md §8). Ketikkan
> `/to-spec` sebagai manusia, lalu berkas ini menjadi briefnya. Bila skill tidak dipakai, tulis
> spec langsung dalam bentuk rumah yang dijelaskan di Bab 3.

---

## 0. LINGKUP — DIKUNCI

Hanya **`NB Treaty In`**. EDM Treaty In **tidak** digarap di ronde ini.

`D:\XML\RNM_BRD\` adalah korpus **READ-ONLY**. Menulis hanya ke `OUTPUT_HASIL_RNM\`.
`D:\XML\nusantara-re\` **terlarang**.

Keluaran satu berkas:
**`OUTPUT_HASIL_RNM\.scratch\nb-treaty-in\spec-penyimpanan-relasional.md`**

⛔ `grilling-ronde-1..4.md` tersegel. ⛔ `spec.md` yang sudah ada **jangan disunting** — ralatnya
ronde tersendiri.

⭐ **Satu skema untuk NB dan EDM.** Tabelnya sama; yang membedakan hanya baris mana yang terisi.
Spec ini menulis skemanya dari sudut NB, dan **wajib menyatakan** bahwa EDM memakai tabel yang sama.

---

## 1. SUMBER DAN URUTAN KEWENANGAN

Bila dua sumber bertentangan, yang di atas menang. **Tuliskan pertentangan yang Anda temukan.**

| # | Sumber | Sifat |
| ---: | --- | --- |
| 1 | `.scratch\nb-treaty-in\rancangan-tabel-datar-treaty-in.md` **Bab 4bis, 4ter, 4quater** | keputusan work owner 22–23 Sep 2026, **mengikat** |
| 2 | `.scratch\nb-treaty-in\PERTANYAAN-untuk-*.md` — enam lembar, P1–P49 | keputusan work owner |
| 3 | `Diagram-Skema-Tabel-NusantaraRe.xlsx` sheet **`NB Treaty In Prop`** dan **`NB Treaty In NonProp`** | bentuk tabel dan relasinya |
| 4 | `.scratch\nb-treaty-in\KEADAAN-NB-TREATY-IN.md` | keadaan terukur |
| 5 | korpus XML | selalu boleh dipakai membuktikan ulang |

⚠️ Berkas rancangan memuat **tiga bab ralat** yang mengutip bunyi lamanya. Yang berlaku selalu
yang **terbaru**; kutipan lama ada supaya jejaknya tidak hilang, bukan untuk dipakai.

Bentuk rumah: `.scratch\claim-life\spec-penyimpanan-relasional.md` *(33 KB, sembilan bab)*.

---

## 2. TABEL YANG DISPEC

### Proporsional — 7 tabel

```
T_WORK_POLIS                          akar · LINTAS-LINI · sudah dipakai PremiumList
 └ T_GENERAL_POLIS                    1:1 SHARED PK · satu baris per GENERASI
    ├ T_POLIS_QUOTATION               1:1
    │   └ T_POLIS_CEDING              1:N
    ├ T_POLIS_INSTALMENT              1:N · SATU tingkat pada proporsional
    ├ T_POLIS_SPREADING               1:N
    └ POOLDATA.HISTORYAKSEPTASIPRODUCTION   1:N · SUDAH datar, tidak dibuat ulang
```

### Non-proporsional — tambah 3

```
    ├ T_POLIS_INSTALMENT_DETAIL       1:N   ← ListInstallment().InstallmentList()
    ├ T_POLIS_XOL                     1:N   ← TreatyXOLList()
    │   └ T_POLIS_XOL_LAYER           1:N   ← .ValueList()
```

⛔ **Tiga tabel proyeksi selisih** *(`T_POLIS_DIFFERENCE` dan anaknya)* **di luar ronde ini** —
di NB nilainya nol baris, sebab `OLD_POLIS_ID` selalu kosong. Sebutkan keberadaannya di
Further Notes, jangan dispec.

---

## 3. BENTUK SPEC

Bab wajib, berurutan:

```
Cara membaca berkas ini    Problem Statement    Solution
User Stories               Implementation Decisions
Acceptance Criteria        Testing Decisions
Out of Scope               Butir [terbuka] — daftar penuh
Further Notes
```

### Blok ringkasan di "Cara membaca berkas ini"

Cacah: user story · acceptance criteria · sebaran penanda · butir `[terbuka]` aktif ·
butir `[penyimpangan sadar]`. **Dihitung dua cara berbeda**, sebutkan jendela hitungnya.
Bila angka diperbaiki, **kutip yang lama** (CLAUDE.md §4a).

### Bentuk Acceptance Criteria

```
N. `[terverifikasi]` <pernyataan>. Test yang menemukan <keadaan sebaliknya> **gagal**. *(Bab X)*
```

Setiap butir membawa **penanda** dan **rujukan bab**. Butir tanpa penanda adalah cacat.
Bab AC **tidak memutuskan apa pun** — ia menyatakan ulang keputusan yang sudah ada.

---

## 4. KEPUTUSAN MENGIKAT — masing-masing wajib jadi satu AC atau lebih

### 4.1 Kunci dan generasi

| # | Ketetapan |
| ---: | --- |
| 1 | `T_GENERAL_POLIS` **satu baris per generasi**; kunci alami `(NOPOLIS, PRODKE)` · NB = `PRODKE 0` |
| 2 | `OLD_POLIS_ID → T_WORK_POLIS.ID` · nullable · **di NB selalu kosong** · `UNIQUE (OLD_POLIS_ID)` melarang percabangan |
| 3 | `UNIQUE (NOPOLIS, PRODKE)` — mencegah dua endorsemen serentak mendapat nomor sama |
| 4 | Baris generasi lampau **tidak boleh disunting** — pembekuan `OldData` *(P58)* |
| 5 | `T_WORK_POLIS` **LINTAS-LINI**, dipakai bersama PremiumList — baris NB dan EDM sejajar, tidak saling menunjuk |

### 4.2 `NOURUT` — wajib di setiap tabel anak

| # | Ketetapan |
| ---: | --- |
| 6 | Setiap tabel anak punya `NOURUT`, dan itulah kunci pasangan antar generasi |
| 7 | ⭐ **Di NB baris boleh dihapus, dan `NOURUT` DINOMORI ULANG rapat `1..n`** — aman karena belum ada generasi untuk dipasangkan |
| 8 | Begitu `PRODKE 0` ditutup, `NOURUT` **beku**. Di EDM tidak ada penghapusan dan nomornya terbawa apa adanya |

### 4.3 Tipe kolom

| # | Ketetapan |
| ---: | --- |
| 9 | Uang dan persen → angka presisi tetap, **skala minimal 9 desimal**. P29: galat lama **diikuti apa adanya**, tidak dibulatkan |
| 10 | Tanggal → `DATE`. **Dua format masuk**: `YYYYMMDD` dan cap waktu Pega bersufiks ` GMT` |
| 11 | Cacah → bilangan bulat: `NOURUT` `PRODKE` `INSTALLMENT_NO` |
| 12 | ⭐ **Kode TETAP teks** — `GroupPanel "006"` dan `BusinessOldId "01"`; nol depannya wajib utuh, kalau tidak penggolong `BusinessType_DeT` 36 baris gagal |
| 13 | ⭐ **Penanda TETAP teks** — `IsApproved` dibandingkan sebagai teks *(P24/P6)*, dan `""` keadaan sah yang **berbeda** dari `"0"` |
| 14 | Teks kosong `""` masuk kolom angka atau tanggal → **`NULL`**, bukan `0` |
| 15 | ⛔ Uang tidak pernah `float` (ADR-0003) · pembandingan uang **tidak boleh sama-persis** |

### 4.4 Isi tabel

| # | Ketetapan |
| ---: | --- |
| 16 | `T_GENERAL_POLIS`: 79 medan skalar `PolicyTreatyIn` + 7 kolom datar dari `POOLDATA.json_polis` |
| 17 | ⛔ `LAYER` `LAYER_TYPE` `LAYER_PART` `LAYER_PART_TYPE` **TIDAK disimpan** di `T_GENERAL_POLIS` — di sistem lama ia pantulan `pxResults(1)`, dan itu **layer pertama saja** |
| 18 | `T_POLIS_QUOTATION` 1:1, 10 medan · memuat `ProportionalType` `GroupPanel` `BusinessOldId` `OldPolicyNo` |
| 19 | `T_POLIS_CEDING` ← `QuotationData.CedingCoList()`, 2 medan: `.CedingCo` dan `.CedingCoName` |
| 20 | ⭐ `CEDING_CO_NAME` dan `CEDING_CO` di `T_GENERAL_POLIS` **DISALIN APA ADANYA**, tidak pernah dirangkai ulang. Keduanya daftar yang digabung berakhiran `"; "` |
| 21 | ⚠️ Konsekuensi diterima: bentuk gabungan dan `T_POLIS_CEDING` **bisa tidak sinkron**, tanpa penjaga |
| 22 | ⭐ **Pada proporsional `T_POLIS_INSTALMENT` SATU tingkat** — dokumen berhenti di `ListInstallment`, tanpa sarang `InstallmentList` |
| 23 | `T_POLIS_SPREADING` — **P60**: NB membagi `100 / jumlah baris` dengan **presisi 10** *(EDM berbeda: presisi 20)* |
| 24 | `T_POLIS_XOL` induk **tanpa penanda layer**; `T_POLIS_XOL_LAYER` yang memegang `LAYER*` — 12 rujukan, di induk **nol** |
| 25 | `DEDUCTION` di XOL adalah **UANG**, bukan persen — terbukti karena **dijumlahkan** antar layer |
| 26 | `HISTORYAKSEPTASIPRODUCTION` ← `PolicyTreatyIn.SuggestList` · `.Suggest → KETERANGAN` *(substr 0,3990)* · `.IsApproved → APPROVAL` *(`"1"`=Accept, `"0"`=Reject)* — **per baris usulan, beda dari tingkat polis** |

### 4.5 Dari ketetapan lama yang tetap berlaku

| # | Ketetapan | Butir |
| ---: | --- | --- |
| 27 | Setiap query menulis skema `POOLDATA.` eksplisit | P3 |
| 28 | Seluruh urutan penyimpanan **satu transaksi** | P2 |
| 29 | `OPERATORID` dari identitas akses login; `PIC` dari nama tampilan | P4, P33 |
| 30 | `DEDUCTION1` `DEDUCTION2` `TOTAL_SHARE_PERCENTAGE_*` adalah **persen** | P29 ket.2 |
| 31 | Arah ketergantungan `handlers → services → repository` | CLAUDE.md §5 |

---

## 5. YANG WAJIB MASUK OUT OF SCOPE

| Yang dikeluarkan | Sebab |
| --- | --- |
| Pohon `LocationList → OccupationList → AnekaList → CoverageList → ClauseList` | ⭐ milik **Fac In** · nol rujukan di NB Treaty In · dipastikan work owner dari contoh JSON |
| `M_TREATY_IN.JSONDATA` dan keluarganya | kontrak · dibaca 11 SQL di 8 modul · ⭐ **Go tetap harus MAMPU MEMBACA JSON**, walau tidak menulisnya |
| Tiga tabel proyeksi selisih | milik EDM · di NB nol baris |
| `TREATYINPRODUCTION` · `HISTORYAKSEPTASIPEGA` · `JSON_POLIS_MONITORING` · `TREATY_IN` · `ACHIEVEMENT` | **dipertahankan apa adanya**, tidak dirancang ulang |
| Penjaga duplikat `ACHIEVEMENT` berbasis nilai uang | ⛔ dipatahkan galat presisi — kuncinya medan pengenal |
| Penghapusan ceding lewat `@replaceAll` pada teks gabungan | tidak ditiru — cukup hapus barisnya |

⚠️ **Wajib ditulis sebagai peringatan, bukan dihilangkan:** `TREATYINPRODUCTION.DEDUCTION1`
menampung **dua satuan** — persen dari jalur prop *(`.Deduction1`)*, uang dari jalur XOL
*(`.Deduction`)*. Pembaca SQL **wajib menyaring `PROPORTIONALTYPE` lebih dulu**.
Dan `LAYER*` pada polis proporsional bernilai `"0"`, bukan kosong.

---

## 6. BUTIR `[terbuka]` YANG WAJIB DIBAWA, JANGAN DITUTUP

| Butir | Pemilik |
| --- | --- |
| ⭐ `JSON_DATAGUIDE(DATA_JSON)` dari DBA — satu-satunya yang dapat menutup selisih **74 lawan 86** medan | `[data DBA]` |
| Presisi fisik kolom uang — skala pasti | `[data DBA]` |
| ⚠️ **Treaty Out masuk lingkup atau tidak.** Folder `NB Treaty In` memuat `InsertToTreatyOutXOLList` · `InputPolicyTreatyOutDetail_preACT` · `DetailPolicyTreatyOutNonProportional` · `BrowseTreatyOut` → `M_treaty_out`. Sapuan berjangkar `PolicyTreatyIn` **tidak menampungnya** | `[work owner]` |
| Migrasi dokumen lama — dipindahkan atau dibaca lewat jalur lama | `[work owner]` |

⛔ **Jangan menutup satu pun.** Penutupan milik work owner, DBA, Product & Underwriting,
Aktuaria, Finance, atau IAM.

---

## 7. DISIPLIN

Setiap pernyataan membawa bukti: **jalur berkas + tipe rule + nama rule**.

Penanda: `[terverifikasi]` · `[dugaan]` · `[terbuka]` · `[data DBA]` · `[keputusan work owner]` ·
`[penyimpangan sadar]`.

⛔ **Jangan menulis `CREATE TABLE` atau DDL apa pun.** Spec, bukan naskah — presisi fisik
dicocokkan DBA di dalam tiket.
⛔ Jangan membuat ADR baru. Bila spec menuntutnya, **katakan di Further Notes**.
⛔ Nol nama orang · nomor polis **tidak** ditulis apa adanya · nol contoh JSON tersimpan ·
nol rahasia, nol data nasabah.

### ⚠️ Empat jebakan yang sudah menjerat proyek ini

1. **Menebak nama tag.** `Activity` memakai `PropertiesName`/`PropertiesValue` **tanpa** awalan
   `py`; `DataTransform` dan `Flow` memakai **dengan** awalan. Menyapu satu saja kehilangan yang lain.
2. **Sapuan berjangkar tidak melihat rujukan relatif.** Di dalam kalang, anggota daftar dirujuk
   `.Suggest` atau `.CedingCoName`, bukan `PolicyTreatyIn.…`. Tiga kali terbukti: `SuggestList`
   terbaca nol medan · `CedingCoList` tidak terbaca sama sekali · lima property EDM terlewat.
3. **Mengambil nilai mayoritas korpus, bukan nilai di aturan yang bersangkutan.** Sekali ini
   menghasilkan laporan palsu bahwa dua kolom tertukar.
4. **`pyStepsPreCondParams` adalah percabangan lompat**, bukan gerbang hidup/mati.

⇒ **Angka medan apa pun adalah BATAS BAWAH, bukan total.** Nyatakan itu di spec.

Jalur korpus memuat spasi — pakai Python, bukan loop shell.
⚠️ Buang blok `pyExpressionGadget` lebih dulu. ⚠️ Jangan `html.unescape` sebelum mencocokkan
pola struktur.

---

## 8. BAB WAJIB — TELEMETRI EKSEKUSI

Spec ditutup dengan bab `## TELEMETRI EKSEKUSI`. Angka token sejati tidak terlihat dari dalam
sesi — ukur dari luar:

```
claude --print --output-format json "<prompt>" > hasil.json
```

Catat token keluaran · cache-read · panggilan alat · durasi · biaya · byte dibaca. Bila pengukuran
dari luar tidak dilakukan, **katakan begitu** — jangan menaksir lalu menyajikannya sebagai angka
terukur.

---

## 9. YANG MENANDAKAN RONDE INI BERHASIL

- Spec lengkap sembilan bab, seluruh AC berpenanda dan berujuk bab.
- **Ketiga puluh satu ketetapan Bab 4** seluruhnya terwakili di Acceptance Criteria.
- Perbedaan NB dan EDM dinyatakan terang-terangan — terutama **`NOURUT` dinomori ulang di NB,
  tidak di EDM**, dan **satu tingkat angsuran pada proporsional**.
- Peringatan `DEDUCTION1` dua satuan dan `LAYER* = "0"` tertulis, bukan dihilangkan.
- Keempat butir `[terbuka]` dibawa utuh, nol yang ditutup sendiri.
- Blok ringkasan dihitung dua cara, jendela hitung disebutkan.
- Nol `CREATE TABLE`, nol ADR baru, nol nama orang.
- Bab telemetri jujur tentang cara pengukurannya.

---

*Disusun 23 September 2026, sesudah rancangan tabel selesai dan dua belas butir penahan ditutup
work owner dalam satu ronde tanya-jawab.*
