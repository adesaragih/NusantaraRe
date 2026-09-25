# PROMPT RANCANGAN PENYIMPANAN — NB Treaty In

> Salin seluruh isi berkas ini sebagai prompt ke sesi eksekutor.
> Ronde ini **merancang bentuk penyimpanan**. Bahan untuk itu kini lengkap.

---

## 0. LINGKUP — DIKUNCI

Modul: **`NB Treaty In`**. `EDM Treaty In` dipakai **hanya** untuk memastikan rancangan menampungnya
juga — sebab keduanya berbagi satu penyimpanan.

`D:\XML\RNM_BRD\` **READ-ONLY**. Menulis hanya ke `OUTPUT_HASIL_RNM\`.
`D:\XML\nusantara-re\` **terlarang**.

⛔ `grilling-ronde-*.md` tersegel. Jangan disunting.

---

## 1. KENAPA RONDE INI BISA JALAN SEKARANG

Dua bahan yang menahannya sudah diterima 22 September 2026, keduanya dari work owner:

| Bahan | Isinya |
| --- | --- |
| Naskah **empat** stored procedure | tabel dan kolom mana yang ditulis sistem lama |
| Satu contoh `DATA_JSON` sungguhan | bentuk dokumen yang harus digantikan |

Rinciannya di jawaban **P1** dan **P29** pada `.scratch\nb-treaty-in\PERTANYAAN-untuk-DBA.md`.

### Yang sudah pasti dari naskah procedure

⭐ **Data kontrak treaty SUDAH relasional.** `POOLDATA.PEGA_TREATY_IN` menulis ke **dua** tabel:
`M_TREATY_IN` (`ID`, `JSONDATA`) **dan** `POOLDATA.TREATY_IN` dengan **20 kolom datar**.

⇒ Yang perlu dipecah **hanya data polis**, bukan data kontrak. Ini memperkecil ronde ini jauh
dari perkiraan semula.

`POOLDATA.PEGA_JSON_POLIS_TREATYIN` menulis satu tabel — `POOLDATA.json_polis` dengan delapan
kolom: `IDPEGA`, `DATA_JSON`, `TGL_INPUT`, `NOPOLIS`, `NOENDORS`, `PRODKE`, `TGL_PROD`, `USERNAME`.

⇒ **Tujuh kolom sudah datar dan dipertahankan apa adanya. Hanya `DATA_JSON` yang dipecah.**

---

## 2. BENTUK DOKUMEN YANG DIPECAH

```
PolicyTreatyIn
  +- LocationList[]
  |    +- OccupationList[]
  |         +- AnekaList[]
  |              +- CoverageList[]
  |                   +- ClauseList[]
  |                   +- DeductibleList[]
  +- SuggestList[]
  +- QuotationData{}        1:1, bukan daftar
```

Lima tingkat pada pohon objek pertanggungan. Perkiraan **delapan tabel** — buktikan atau koreksi.

### Lima sifat yang mengikat, seluruhnya `[terverifikasi]`

1. **Seluruh nilai bertipe teks** di dalam dokumen — termasuk uang, tanggal, dan penanda.
2. **Dua format tanggal** dalam satu dokumen: `YYYYMMDD` dan cap waktu Pega bersufiks ` GMT`.
3. `pxObjClass` ada di **setiap simpul** — internal Pega, **tidak dimigrasi**.
4. `CedingCo` di dalam `QuotationData` **berakhiran `"; "`** — daftar bergabung, bukan nilai tunggal.
5. `SuggestList` membawa `IsApproved` **per baris**, berbeda dari `IsApproved` tingkat atas.

⚠️ **Contoh JSON memuat nama orang** pada medan pemasar, operator, dan tertanggung.
⛔ **Jangan menyalin nilainya** ke berkas mana pun. Nama medan dan jumlahnya saja.

---

## 3. YANG HARUS DIHASILKAN

Dua berkas, mengikuti jalan yang sudah ditempuh `claim-life` dan `premiumlist-life`.

### 3.1 `.scratch\nb-treaty-in\revisi-penyimpanan-nb-treaty-in.md`

Rancangan tabelnya. Untuk tiap tabel: nama, induknya, kardinalitas, kunci, dan **daftar kolom
lengkap** dengan tipenya.

Wajib menyatakan, per kolom, **dari mana asalnya**:

| Asal | Artinya |
| --- | --- |
| `[dari DATA_JSON]` | dipecah dari dokumen |
| `[dari json_polis]` | tujuh kolom datar yang sudah ada |
| `[dari TREATY_IN]` | 20 kolom kontrak yang sudah relasional |
| `[baru]` | tidak ada padanannya di sistem lama — **sebutkan alasannya** |

### 3.2 `.scratch\nb-treaty-in\spec-penyimpanan-relasional.md`

Spec tersendiri untuk penyimpanan, bentuk rumah seperti `claim-life/spec-penyimpanan-relasional.md`:
Problem Statement · Solution · User Stories · Implementation Decisions · Acceptance Criteria ·
Testing Decisions · Out of Scope · Butir `[terbuka]` · Further Notes.

---

## 4. KEPUTUSAN YANG SUDAH MENGIKAT — JANGAN DIPUTUSKAN ULANG

Dari `KEADAAN-NB-TREATY-IN.md` Bab 4 dan jawaban yang sudah ada:

| # | Ketetapan | Butir |
| ---: | --- | --- |
| 1 | Uang berpresisi penuh; pembulatan hanya di titik penyajian | P29 |
| 2 | `DEDUCTION1` `DEDUCTION2` `BROKERAGE` `RNM_SHARE` **persentase**, bukan uang | P29 |
| 3 | `COMMENCEMENT` `TERMINATION` dibaca apa adanya; **P32** berlaku untuk penulisan baru | P29 |
| 4 | Skema `POOLDATA.` ditulis eksplisit di setiap query | P3 |
| 5 | Seluruh urutan penyimpanan **satu transaksi** | P2 |
| 6 | `OPERATORID` dari identitas akses login; `PIC` dari nama tampilan | P4, P33 |
| 7 | Uang tidak pernah `float` | ADR-0003 |
| 8 | Arah ketergantungan `handlers → services → repository` | CLAUDE.md §5 |

### Dan satu yang khusus untuk penyimpanan

⭐ **EDM Treaty In memakai penyimpanan yang sama.** `PEGA_JSON_POLIS_TREATYIN` dipanggil
NB Treaty In, NB FacIn, **dan** EDM Treaty In dengan delapan parameter yang sama — EDM mengisi
`NOENDORS` dan `PRODKE` yang NB kirim sebagai `NULL` dan `'0'`.

⇒ **Rancangan wajib menampung endorsemen sejak awal.** Uji rancangan Anda terhadap bentuk
endorsemen: nomor `…/E01`, data lama beku, selisih per medan, dan endorsemen berlapis (**P57**).

---

## 5. YANG TIDAK BOLEH DIKERJAKAN

- ⛔ Jangan menulis `CREATE TABLE` atau DDL apa pun. **Rancangan, bukan naskah.** Presisi fisik
  dicocokkan DBA belakangan, di dalam tiket — bukan prasyarat.
- ⛔ Jangan memutuskan ulang butir pada Bab 4.
- ⛔ Jangan menutup butir `[terbuka]` milik pihak lain.
- ⛔ Jangan membuat ADR baru. Bila rancangan menuntutnya, **katakan di Further Notes**.
- ⛔ Jangan menyalin nilai berupa nama orang, dan jangan menyimpan contoh JSON apa adanya.

---

## 6. EMPAT PERTANYAAN YANG HARUS DIJAWAB RANCANGAN

1. **Kolom mana yang dipertahankan sebagai teks**, dan mana yang diubah jadi bilangan atau tanggal?
   Ingat: di dokumen lama **semuanya teks**. Setiap perubahan tipe adalah keputusan, dan perlu
   alasannya.
2. **`CedingCo` bergabung dengan `"; "`** — dipecah jadi baris tersendiri, atau dipertahankan
   sebagai teks? Sebutkan akibat masing-masing.
3. **Pohon lima tingkat** — delapan tabel, atau sebagian dilipat? Melipat mengurangi jumlah tabel
   tetapi menghilangkan kemampuan menanyai tingkat yang dilipat.
4. **Migrasi data lama** — dokumen yang sudah ada di `DATA_JSON` dipindahkan, atau dibiarkan dan
   dibaca lewat jalur lama? Sebutkan mana yang Anda sarankan dan mengapa.

---

## 7. BAB WAJIB — TELEMETRI EKSEKUSI

Di akhir `revisi-penyimpanan-nb-treaty-in.md`. Ukur dari luar bila bisa:

```
claude --print --output-format json "<prompt>" > hasil.json
```

Bila tidak, katakan begitu. Bila hasil pengurangan baseline, katakan itu juga.

---

## 8. YANG MENANDAKAN RONDE INI BERHASIL

- Setiap kolom pada rancangan menyebut **asalnya**, dan kolom `[baru]` menyebut alasannya.
- Keempat pertanyaan Bab 6 dijawab dengan pilihan **dan** akibatnya, bukan hanya pilihannya.
- Rancangan diuji terhadap bentuk **endorsemen**, bukan hanya polis baru.
- Nol `CREATE TABLE`, nol nama orang, nol contoh JSON tersimpan.
- Nol keputusan Bab 4 yang diputuskan ulang, nol ADR baru.
- Bab telemetri jujur tentang cara pengukurannya.

---

*Disusun 22 September 2026, sesudah naskah empat stored procedure dan satu contoh `DATA_JSON`
diterima — dua bahan yang menahan ronde ini sejak awal proyek.*
