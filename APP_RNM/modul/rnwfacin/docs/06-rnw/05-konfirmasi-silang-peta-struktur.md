# Discovery Renewal — Konfirmasi silang dari peta struktur rule (RNW-6)

> **Sumber:** `D:\migrasi\RNM\RNW Fac In\Struktur_InputRenewalFacultativeIn.xlsx`
> (OOXML sah, 7,02 MB, sheet tunggal `Struktur`, **81.379 baris × 28 kolom**).
>
> ⛔ **Kolom `XML` (kolom 26) TIDAK DIBACA sama sekali** — ia memuat cuplikan XML rule Pega yang
> terbukti mengandung nama operator. Yang dibaca hanya: `Nama Rule` · `Jenis Rule` · `Dipanggil dari` ·
> `Status XML`. **Tidak ada nama orang yang disalin ke dokumen ini.**

---

## Ringkasan

✅ **Konfirmasi silang mendukung keputusan yang sudah diambil**, dari sumber yang independen dari
korpus XML. Dan ✅ **tidak ditemukan rule delta yang terlewat** oleh discovery 20-berkas.

---

## 1. Bentuk berkasnya

`[terverifikasi]` **81.378 baris data** (baris 1 = header), **1.825 nama rule unik** setelah
normalisasi.

Karena jumlah baris jauh melebihi jumlah rule, ini **bukan** satu baris per rule — melainkan
**daftar rujukan (edge)**: satu baris per pemakaian rule, sejalan dengan adanya kolom
`Dipanggil dari`.

| Jenis Rule | Baris | | Jenis Rule | Baris |
| --- | ---: | --- | --- | ---: |
| When | 37.776 | | Flow Action | 4.325 |
| Activity | 11.895 | | Data Page | 1.031 |
| Data Transform | 7.594 | | Harness | 525 |
| Report Definition | 6.885 | | Decision Table | 140 |
| Section | 6.138 | | Connect REST | 46 |
| RDB List | 4.971 | | System Settings · Flow · Decision Tree | 40 · 8 · 4 |

`[terverifikasi]` Kolom `Status XML` bernilai **`Sudah diupload` untuk seluruh 81.378 baris** — satu
nilai saja, tidak ada variasi. Kolom ini karena itu **tidak membedakan apa pun** dan tidak dapat
dipakai sebagai penanda rule yang hilang.

---

## 2. ⚠️ Satu koreksi pengukuran — dicatat agar tidak terulang

Perbandingan pertama saya menghasilkan **221 nama "tidak ada sebagai berkas"**, yang sempat tampak
seperti temuan besar. **Itu keliru.**

`[terverifikasi]` Kedua ratus dua puluh satu nama itu adalah **kunci rule berformat lengkap** —
`kelas + access group + nama` — bukan rule yang hilang:

```
ASM-FW-GISFW-Int-BUSINESS ASM CariBusinessGID
ASM-FW-GISFW-Int-ACCUMULATION ASM GetAccumulation_SQL
```

Baris pertama itu **adalah** `CariBusinessGID`, salah satu dari 20 berkas delta. Peta mencampur dua
format penamaan: nama pendek dan kunci lengkap, terutama pada **RDB List** (203 dari 221).

Setelah dinormalisasi (ambil token terakhir), angkanya runtuh dari 221 menjadi **17**.

📌 **Pola yang sama dengan pelajaran K-032:** perbandingan berbasis nama menipu, dan kali ini yang
tertipu adalah pengukuran saya sendiri. Normalisasi format wajib dilakukan **sebelum** membandingkan.

---

## 3. Konfirmasi (a) — ketiga rujukan yang sudah diputuskan

Uji: apakah `GenerateNoPolicy`, `SumTreatyCapacity_Act`, `SearchJobID` muncul di peta?

`[terverifikasi]` setelah normalisasi:

| Rule | Sebagai `Nama Rule` | Sebagai `Dipanggil dari` |
| --- | :-: | :-: |
| `GenerateNoPolicy` | **tidak** | **tidak** |
| `SumTreatyCapacity_Act` | **tidak** | **tidak** |
| `SearchJobID` | **tidak** | **tidak** |

### Kontrol — membuktikan instrumennya memang mendeteksi

Ketiadaan hanya bermakna bila yang **seharusnya ada** memang terdeteksi. Empat nama kontrol dari 20
berkas delta:

| Kontrol | `Nama Rule` | `Dipanggil dari` |
| --- | :-: | :-: |
| `CariBusinessGID` | **ya** | tidak |
| `GetBusinessGroup_Act` | **ya** | **ya** |
| `InputRenewal` | **ya** | **ya** |
| `RenewalList_RD` | **ya** | tidak |

✅ **Instrumen peka.** Ketiadaan ketiga rujukan karena itu **nyata**, bukan artefak pengukuran.

### Apa yang boleh dan tidak boleh disimpulkan

| ✅ Boleh | ⛔ Tidak boleh |
| --- | --- |
| `[terverifikasi]` Ketiganya **tidak tercatat** di peta struktur rule renewal | Menyatakan ini **membuktikan** ketiganya usang |

`[dugaan]` Ketiadaannya **konsisten dengan** keputusan work owner (K-032): dua usang, satu tidak ada
di ruleset renewal. Tetapi peta ini **buatan tim**, bukan ekstraksi mesin dari ruleset — bila sebuah
rule sudah dihapus, ia memang tidak akan terdaftar. **Konsisten-dengan, bukan bukti.**

**K-032 tidak dinaikkan tingkat buktinya atas dasar berkas ini.**

---

## 4. Konfirmasi (b) — adakah rule delta yang terlewat?

`[terverifikasi]` Setelah normalisasi:

| Ukuran | Nilai |
| --- | ---: |
| Nama rule unik di peta | **1.825** |
| — ada sebagai berkas RNW | **1.808** |
| — tidak ada sebagai berkas RNW | **17** |
| — di antaranya ada di NB | **0** |
| Berkas RNW yang tidak muncul di peta | **4** |

✅ **Tidak ada rule delta yang terlewat discovery.** Tidak satu pun nama di peta merupakan berkas NB
yang absen dari RNW — angka itu **nol**, dan itu sekaligus konfirmasi independen bahwa daftar 20
berkas delta lengkap.

### Selisih kecil yang tersisa

Perbandingan peka-huruf memberi **17** dan **4**; perbandingan tidak-peka-huruf memberi **2** dan
**3**. Selisihnya adalah **varian kapitalisasi** — mis. `isCAR` vs `IsCAR`.

⚠️ Ini menyentuh `[pertanyaan terbuka]` lama yang belum terjawab: **apakah resolusi rule Pega di
instalasi ini peka huruf?** Bila ya, varian kapitalisasi adalah rule yang berbeda. Peta ini
**tidak menjawabnya**, hanya memperlihatkan gejalanya lagi.

Yang tersisa setelah penyamaan huruf:

| Di peta, bukan berkas | Jenis |
| --- | --- |
| `isCAR` · `IsGroupUWFac` | When |

| Berkas RNW, tidak di peta | Jenis |
| --- | --- |
| `AddCoverageProRate` | DataTransform |
| `FlagOldData` | When |
| `ViewDtlAdditionalDetail` | FlowAction |

`[dugaan]` Ketiga berkas yang tidak terdaftar kemungkinan tidak tersentuh alur renewal yang dipetakan
tim — bukan berarti tidak ada. Tidak dikejar lebih jauh; dampaknya kecil dan tidak mengubah kesimpulan.

---

## 5. Kesimpulan

| Tujuan | Hasil |
| --- | --- |
| (a) Konfirmasi ketiga rujukan | ✅ **Konsisten** dengan K-032 — ketiganya tidak tercatat, dan instrumen terbukti peka. `[dugaan]`, bukan bukti |
| (b) Rule delta yang terlewat | ✅ **Tidak ada.** Nol nama peta yang merupakan berkas NB absen dari RNW |

**Nilai tambahan yang tidak diduga:** peta ini memberi **graf pemanggilan renewal** (874 pemanggil
unik) dari sumber independen — berguna bila kelak diperlukan penelusuran jalur, meski tidak dipakai
pada putaran ini.

---

*Kolom `XML` tidak dibaca. Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
