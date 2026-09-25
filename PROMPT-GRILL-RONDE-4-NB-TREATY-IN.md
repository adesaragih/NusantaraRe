# PROMPT GRILLING — RONDE 4 · NB Treaty In

> Salin seluruh isi berkas ini sebagai prompt ke sesi eksekutor.
> Ronde 1–3 sudah selesai. Ronde ini menutup bagian modul yang **belum pernah dibuka sama sekali**.

---

## 0. LINGKUP — DIKUNCI

Hanya modul **`D:\XML\RNM_BRD\NB Treaty In`**.

Modul lain boleh dibaca **hanya** untuk membuktikan bahwa sesuatu ada atau tidak ada di luar
NB Treaty In. Temuan dari modul lain tidak masuk sensus.

`D:\XML\RNM_BRD\` adalah korpus **READ-ONLY**. Menulis hanya ke `OUTPUT_HASIL_RNM\`.
`D:\XML\nusantara-re\` **terlarang** — jangan dibaca, dikutip, dibandingkan, atau dijadikan sasaran.

---

## 1. YANG SUDAH DIKETAHUI — JANGAN DIULANG

Ronde 1–3 sudah menghasilkan 43 pertanyaan, **29 di antaranya sudah terjawab**. Jangan menemukan
ulang hal-hal berikut:

- 25.986 langkah `Property-Set` mengekspor nol referensi properti. Sudah diminta ekspor ulang (P18).
- `IsApproved`: `0` = ditolak, selain itu = disetujui. Aturan hidup ada di `DecisionTable`,
  bukan `When`. Sudah diputuskan.
- `InputData.CARIn` dan 33 halaman slot lain: arti slot **per-aturan**, tidak ada makna tetap.
- 75 rule `When`: bila keterangan berbeda dari syarat yang dijalankan, **yang dijalankan benar**.
- 17 langkah `Page-Clear-Messages`: 12 pembersihan awal, 5 ditahan menunggu pemastian.

### KOREKSI PENTING — naskah SQL ADA di dalam korpus

Ronde sebelumnya pernah menyimpulkan naskah SQL tidak ikut terekspor. **Itu keliru.**

Naskah SQL tersimpan di tag **`<pyBrowseSQL>`** pada rule tipe **`RDBList`**:
**41 pernyataan di NB Treaty In, 1.151 di seluruh korpus.**

Begitu pula baris tabel keputusan — tersimpan di `pyCondition`, `pyOrConditions`, dan `pyResults`,
**bukan** di `pyCriteriaValue` atau `pyResult`.

Sebelum menyatakan sesuatu "tidak ada di korpus", **cari dulu tag yang benar**, dan sebutkan tag
apa saja yang sudah dicari.

---

## 2. SASARAN RONDE INI

### Sasaran 1 — LAPISAN LAYAR, BELUM DIBUKA SAMA SEKALI (utama)

| | Berkas | Ukuran |
| --- | ---: | ---: |
| `Section` | 25 | 14.229.383 B |
| `Harness` | 6 | 2.406.451 B |
| **jumlah** | **31** | **16.635.834 B** |

Itu **46,6 %** dari 35.678.284 B seluruh modul (278 berkas). **Nol berkas dibuka** di ronde 1–3.

Lima terbesar, semuanya `Section`:

| Ukuran | Nama |
| ---: | --- |
| 1.926.378 | `GeneralPolicyTreatyIn` |
| 1.916.240 | `DetailPolicyTreatyIn` |
| 1.662.941 | `DetailPolicyTreatyInNonProportional` |
| 1.584.498 | `GeneralDeptHeadTreatyIn_UW` |
| 1.567.128 | `DetailDeptHeadTreatyIn_UW` |

**Yang dicari: aturan yang hanya hidup di layar dan tidak ada di tempat lain.**

- syarat tampil-sembunyi, dan apa yang mengendalikannya
- medan wajib, dan kapan wajibnya berlaku
- medan hanya-baca, dan setelan mana yang menguncinya
- validasi sisi layar: batas nilai, format, ketergantungan antar medan
- urutan pengisian yang dipaksakan
- nilai awal yang ditetapkan oleh layar, bukan oleh aturan
- medan yang ada di layar tetapi tidak ada di aktivitas mana pun, dan sebaliknya

Berkas-berkas ini besar. Bila satu berkas tidak muat dibaca sekaligus, baca bertahap dan
**catat berapa byte yang benar-benar dibaca** — jangan menyimpulkan dari bagian yang belum dibuka.

### Sasaran 2 — P22, SENSUS ULANG ELEMEN MATI DAN TERKUNCI

Ronde 2 mencatat **91 bagian layar dibuat agar tidak pernah muncul** (`1=2` 72×, `NEVER` 15×,
`1==2` 2×) dan **38 dikunci permanen** (`pyReadOnlyCondition` `1==1` 32× ditambah 6 varian).

**Hitung ulang dua cara.** Angka ini akan dikirim ke Product & Underwriting sebagai dasar keputusan
"masih perlu dibangun atau tidak", jadi ia harus benar. Katakan bila berbeda dari 91 dan 38.

Untuk tiap elemen mati, catat **apa yang disembunyikannya** — satu medan, satu bagian, atau satu tab
penuh. Bedanya besar bagi beban kerja pembangunan.

### Sasaran 3 — P3, EJAAN NAMA TABEL

Empat nama tabel muncul dengan dan tanpa awalan skema di dalam 1.151 naskah SQL:

| Nama | Dengan `POOLDATA.` | Tanpa awalan |
| --- | ---: | ---: |
| `REINSURANCETYPE` | 18 | 11 |
| `PROPORTIONALARRG` | 2 | 25 |
| `HISTORYAKSEPTASIPEGA` | 0 | 15 |
| `CURRENCY` | 4 | **perlu dihitung ulang** |

**Peringatan: angka `CURRENCY` di atas tercemar.** Penelusuran yang menghasilkannya ikut menangkap
rujukan properti seperti `DATA_JSON.CURRENCY` dan `osAkseptasi.CURRENCY`, yang **bukan** nama tabel.
Hitung ulang dengan memisahkan rujukan tabel — yang muncul sesudah `FROM`, `JOIN`, `INTO`, `UPDATE` —
dari rujukan properti.

**Yang dihasilkan:** untuk tiap nama, daftar berkas mana memakai ejaan mana.
**Jangan menyimpulkan apakah keduanya tabel yang sama.** Itu keputusan DBA.

### Sasaran 4 — P4, KOLOM `OPERATORID`

Kolom `OPERATORID` pada `POOLDATA.HISTORYAKSEPTASIPEGA` disebut tidak pernah diisi dari modul ini.

Buktikan atau bantah dari 1.151 naskah SQL: adakah pernyataan mana pun, di modul mana pun, yang
menulis ke kolom itu.

**Peringatan:** jangan tertukar dengan halaman Pega `OperatorID` — misalnya
`OperatorID.pyUserIdentifier` — yang sering muncul sebagai **parameter**, bukan sebagai nama kolom.
Keduanya dieja sama. Penelusuran naif akan menghasilkan puluhan hasil palsu.

---

## 3. DISIPLIN — WAJIB

### Anti-halusinasi

Setiap pernyataan membawa bukti: **jalur berkas + tipe rule + nama rule**. Tanpa itu, jangan ditulis.

Beri label setiap temuan:
`[terverifikasi]` · `[dugaan]` · `[terbuka]` · `[data DBA]` · `[keputusan work owner]` ·
`[penyimpangan sadar]`

Setiap angka membawa **perintah yang menghasilkannya**, supaya orang lain dapat mengujinya ulang.

### Sensus dua cara

Setiap angka dihitung **dua cara yang benar-benar berbeda**. Bila keduanya tidak sama, tulis keduanya,
lalu katakan mana yang dipercaya dan mengapa.

### Tujuh jebakan yang sudah pernah menjerat ronde sebelumnya

1. **Spasi di dalam jalur.** `NB Treaty In/...` terpecah oleh shell. Pakai Python, bukan loop shell.
2. **`<rowdata REPEATINGINDEX="n"/>` yang menutup sendiri** tidak tertangkap pola
   `<rowdata...>(.*?)</rowdata>`. Sel kosong hilang, daftar menjadi pendek, pasangan bergeser.
3. **`pyRowNum` pada `pyOrConditions` berbasis NOL**, sedangkan baris tabel berbasis satu.
4. **Memasangkan dua daftar menurut URUTAN**, bukan menurut kunci. Hanya sah bila keduanya berada di
   dalam `rowdata` yang sama.
5. **Pencocokan substring.** `grep -i 'ITL'` juga menangkap `pyTitleEnabled`. Pakai pola tag utuh.
6. **Menebak nama tag.** `pyShapeName`, `pyTaskLabel`, `pyStepsPage`, `pyCriteriaValue` — keempatnya
   **tidak ada** di ekspor ini. Periksa daftar tag lebih dulu.
7. **Menyimpulkan "tidak ada di korpus" dari penelusuran yang sempit lingkupnya.** Sebutkan lingkup
   penelusuran setiap kali menyatakan sesuatu nihil.

### Batas

- **Nol kode, nol DDL, nol `CREATE TABLE`, nol usulan daftar kolom.** Ronde ini membaca, bukan merancang.
- **Jangan menutup pertanyaan terbuka sendiri.** Penutupan milik work owner, DBA, Product &
  Underwriting, Aktuaria, Finance, atau IAM.
- **Jangan menyalin nilai yang berupa nama orang.** Catat nama rule, nama medan, dan jumlahnya saja.
- Nol rahasia, nol token, nol data nasabah, nol cuplikan data produksi.

---

## 4. KELUARAN — DUA BERKAS, BUKAN SATU

Ronde 1–3 menaruh pertanyaan di dalam berkas grilling, dan pertanyaannya terkubur di antara ribuan
baris. Jangan diulang.

### Berkas 1 — `OUTPUT_HASIL_RNM\.scratch\nb-treaty-in\grilling-ronde-4.md`

Temuan lengkap, dengan bukti dan perintah audit. Bab A sampai D mengikuti keempat sasaran di atas.

### Berkas 2 — `OUTPUT_HASIL_RNM\.scratch\nb-treaty-in\PERTANYAAN-RONDE-4.md`

**Hanya pertanyaan.** Tidak ada temuan, tidak ada analisis, tidak ada ringkasan.

Bentuk tiap pertanyaan, sama seperti ronde sebelumnya:

```
## Pn — <judul dalam bahasa bisnis, bukan bahasa Pega>

<satu atau dua paragraf, bahasa bisnis>

**Konteks:** <mengapa ini tidak bisa dijawab dari berkas>

**Bentuk jawaban yang diharapkan:** <pilihan ganda / butuh berkas / butuh contoh>

**Dampak bila salah:** <akibat nyata>

rujukan: <bukti teknis, boleh diabaikan pembaca bisnis>

**Jawaban:**

> *(tulis di sini)*
```

Nomor pertanyaan **mulai dari P42.** P1 sampai P41 sudah terpakai; P41 sudah ditarik.

Kelompokkan menurut pemilik: pengembang Pega lama · Product & Underwriting · DBA · Finance · IAM ·
pemilik export Pega.

---

## 5. BAB WAJIB — TELEMETRI EKSEKUSI

Berkas grilling **wajib** ditutup dengan bab `## TELEMETRI EKSEKUSI`.

Angka token yang sejati **tidak terlihat dari dalam sesi**. Ukur dari luar:

```
claude --print --output-format json "<prompt>" > hasil.json
```

Muatan balasannya memuat `total_cost_usd`, `usage` (rincian token per model), `duration_ms`,
`num_turns`, dan `session_id`.

Catat di bab itu:

| Yang dicatat | Nilai |
| --- | --- |
| token keluaran | |
| token cache-read | |
| jumlah panggilan alat | |
| durasi | |
| biaya | |
| berkas dibuka / byte dibaca | |

Bila pengukuran dari luar tidak dilakukan, **katakan begitu**. Jangan menaksir lalu menyajikannya
sebagai angka terukur.

---

## 6. YANG MENANDAKAN RONDE INI BERHASIL

- Ke-31 berkas layar **dibuka**, bukan hanya didaftar. Sebutkan berapa byte benar-benar dibaca.
- Sensus 91 elemen mati dan 38 elemen terkunci **dikonfirmasi atau dikoreksi**, dengan dua cara hitung.
- P3 dan P4 dijawab **dari korpus**, atau dinyatakan tegas tidak termuat di korpus beserta lingkup
  penelusuran yang sudah dijalankan.
- Pertanyaan baru berada di berkasnya sendiri, bernomor mulai P42, siap dikirim tanpa disunting lagi.
- Bab telemetri terisi angka terukur, bukan taksiran.

---

*Disusun dari hasil ronde 1–3 dan 29 jawaban work owner yang sudah tercatat. Seluruh angka di dalam
prompt ini diverifikasi ulang pada 2026-09-22 dan boleh diuji ulang.*
