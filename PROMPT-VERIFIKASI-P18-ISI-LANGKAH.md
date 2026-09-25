# PROMPT VERIFIKASI — apakah isi langkah penetapan nilai benar-benar tidak terekspor?

> Salin seluruh isi berkas ini sebagai prompt ke sesi eksekutor.
> Ronde ini **hanya memeriksa satu klaim**. Tidak menulis spec, tidak menulis tiket, tidak
> menemukan hal baru.

---

## 0. LINGKUP — DIKUNCI

Modul: **`D:\XML\RNM_BRD\NB Treaty In`** dan **`D:\XML\RNM_BRD\EDM Treaty In`**.

`D:\XML\RNM_BRD\` adalah korpus **READ-ONLY**. Menulis hanya ke `OUTPUT_HASIL_RNM\`.
`D:\XML\nusantara-re\` **terlarang**.

Keluaran satu berkas: **`OUTPUT_HASIL_RNM\.scratch\nb-treaty-in\VERIFIKASI-P18.md`**.

⛔ Jangan menyunting berkas grilling, lembar jawaban, spec, atau berkas keadaan. Ronde ini
**melaporkan**, tidak memperbaiki.

---

## 1. KLAIM YANG DIPERIKSA

Selama empat ronde, proyek ini bekerja di atas satu anggapan:

> **P18** — *"Isi 268 langkah penetapan nilai tidak ada di dalam berkas yang kami terima. Kami
> dapat melihat bahwa langkah itu ada dan urutannya, tetapi tidak satu pun rumus atau nilai yang
> ditetapkannya ikut terkirim."*

Anggapan itu menjadikan P18 **penahan terberat proyek** — rantai perhitungan uang dikarantina di
spec, sejumlah tiket bertanda `blocked`, dan satu surat permintaan ekspor ulang sudah disiapkan.

**Sesi asisten kini menduga anggapan itu KELIRU**, dan bahwa isinya ada di dalam ekspor.

### Dugaan yang harus diuji

Isi langkah tersimpan di tag **`PropertiesName`** dan **`PropertiesValue`** — **tanpa awalan
`py`** — di dalam blok `rowdata` milik langkah.

Ronde-ronde sebelumnya mencari `pyPropRef`, `pyPropertiesName`, dan `pyPropertiesValue` — **dengan**
awalan — menemukan nol, lalu menyimpulkan isinya tidak terekspor.

### Angka yang diklaim sesi asisten

| Pernyataan | NB Treaty In | EDM Treaty In |
| --- | ---: | ---: |
| `Activity` | 92 | 66 |
| langkah `Property-Set` | 938 | 378 |
| **pasangan `PropertiesName`=`PropertiesValue` terisi** | **2.080** | **1.062** |
| `Activity` punya langkah tetapi **nol pasangan terisi** | **0** | — |
| nilai yang tampak terpotong | **0** | — |

⚠️ **Jangan percaya angka ini.** Ukur sendiri.

---

## 2. YANG WAJIB DIKERJAKAN

### 2.1 Cacah, dua cara

Hitung pasangan `PropertiesName`=`PropertiesValue` yang **keduanya tidak kosong**, dengan **dua
cara yang benar-benar berbeda** — misalnya pengurai `ElementTree` dan pencocokan pola teks.

Laporkan keduanya. Bila berbeda, katakan mana yang dipercaya dan mengapa.

### 2.2 Uji keutuhan — ini yang paling menentukan

Klaim lama menyebut isi yang ada hanyalah *"salinan sementara dari layar penyunting"* yang
**terpotong di tengah**. Uji itu:

- berapa nilai yang **kurungnya tidak seimbang**;
- berapa yang **berakhir dengan operator** (`+ - * / ,`);
- berapa yang berakhir dengan tanda kutip tak tertutup;
- sebaran panjang nilai — terpendek, median, terpanjang;
- adakah nilai yang berakhir tepat pada batas bulat yang mencurigakan (255, 256, 1024) — tanda
  pemotongan oleh batas kolom.

**Bila nol terpotong, katakan nol.** Bila ada, sebutkan berapa dan di aturan mana.

### 2.3 Uji kecukupan per langkah

Yang dicacah di 2.1 adalah **pasangan**, bukan langkah. Satu langkah `Property-Set` dapat menetapkan
beberapa properti sekaligus.

Yang harus dibuktikan: **apakah setiap langkah `Property-Set` punya isinya.**

Periksa **di dalam blok langkah itu sendiri**, bukan di tingkat berkas. Laporkan:

- berapa langkah `Property-Set` yang blok-nya memuat sekurangnya satu pasangan terisi;
- berapa yang **tidak**, dan di aturan mana.

⚠️ Ini pemeriksaan yang paling mudah salah. Blok `rowdata` bersarang — pastikan Anda memasangkan
di dalam blok yang **sama**, bukan menurut urutan kemunculan di berkas.

### 2.4 Uji kebermaknaan — ambil sampel

Ambil **lima aturan** dari daftar teratas `LAMPIRAN-P18-ATURAN-DIMINTA.md`, dan untuk masing-masing
tampilkan **seluruh** pasangan nama=nilainya.

Pertanyaannya sederhana, dan jawabannya harus penilaian manusia, bukan angka:

> **Cukupkah ini untuk menulis ulang perhitungannya di sistem baru?**

Bila ada bagian yang tetap tidak terbaca — misalnya rumus yang memanggil fungsi yang naskahnya
tidak ada — sebutkan apa.

### 2.5 Periksa rule tipe lain

Langkah `Property-Set` ada juga di luar `Activity` — periksa `DataTransform` dan `FlowAction`.
Berlaku pola yang sama atau tidak?

---

## 3. YANG TIDAK BOLEH DIKERJAKAN

- ⛔ Jangan menyunting `PERTANYAAN-untuk-pemilik-export-Pega.md`, `spec.md`, `KEADAAN-…`, atau
  berkas grilling.
- ⛔ Jangan menutup P18. Penutupannya keputusan work owner sesudah laporan ini dibaca.
- ⛔ Jangan menulis kode, DDL, atau usulan kolom.
- ⛔ Jangan menyalin nilai berupa nama orang, dan jangan menulis nomor polis apa adanya.

---

## 4. KELUARAN

`.scratch\nb-treaty-in\VERIFIKASI-P18.md`, enam bab.

| Bab | Isi |
| --- | --- |
| 1 | cacah dua cara, kedua modul |
| 2 | uji keutuhan — terpotong atau tidak, dengan angkanya |
| 3 | uji kecukupan per langkah — berapa langkah tanpa isi, di aturan mana |
| 4 | lima sampel aturan, seluruh pasangannya, beserta penilaian cukup atau tidak |
| 5 | `DataTransform` dan `FlowAction` — pola sama atau berbeda |
| 6 | **putusan** |

### Bab 6 menjawab satu pertanyaan, dengan tegas

> **P18 masih diperlukan, atau tidak?**

Tiga kemungkinan, pilih satu dan sebutkan dasarnya:

- **tidak diperlukan** — isinya lengkap, permintaan ekspor ulang dibatalkan;
- **diperlukan sebagian** — sebutkan aturan mana saja yang isinya benar-benar kurang;
- **tetap diperlukan** — sebutkan di mana dugaan sesi asisten keliru.

Bila putusannya "tidak diperlukan", sebutkan juga **apa lagi yang ikut berubah**: karantina rantai
uang di `spec.md`, tiket bertanda `blocked`, dan surat `SURAT-P18-KE-PEMILIK-EXPORT-PEGA.md`.

### Bab telemetri

Wajib, di akhir. Ukur dari luar bila bisa:

```
claude --print --output-format json "<prompt>" > hasil.json
```

Bila tidak, katakan begitu. Bila angkanya hasil pengurangan baseline, katakan itu juga.

---

## 5. DISIPLIN

Setiap angka membawa **perintah yang menghasilkannya**.
Penanda: `[terverifikasi]` · `[dugaan]` · `[terbuka]`.

Jalur korpus memuat spasi — pakai Python, bukan loop shell.

### Lima jebakan yang sudah pernah menjerat proyek ini

1. **Menebak nama tag.** Inilah jebakan yang melahirkan ronde ini. `pyPropRef`,
   `pyPropertiesValue`, `pyStepsPage`, `pyCriteriaValue`, `pyResult`, `pyShapeName` — semuanya
   pernah dicari dan tidak ada. **Periksa daftar tag yang benar-benar ada lebih dulu.**
2. `<rowdata REPEATINGINDEX="n"/>` yang menutup sendiri tidak tertangkap pola
   `<rowdata...>(.*?)</rowdata>`.
3. `pyRowNum` pada `pyOrConditions` berbasis nol; baris tabel berbasis satu.
4. Memasangkan dua daftar menurut **urutan**, bukan menurut blok yang sama.
5. Menyatakan sesuatu nihil tanpa menyebut lingkup penelusuran.

---

## 6. YANG MENANDAKAN RONDE INI BERHASIL

- Keempat angka yang diklaim diuji ulang, hasilnya ditulis apa adanya — **termasuk bila ternyata
  sesi asisten salah**.
- Uji keutuhan dijalankan dengan sekurangnya tiga cara berbeda.
- Lima sampel ditampilkan **utuh**, bukan diringkas.
- Bab 6 memuat putusan tegas, bukan ringkasan.
- Bila putusannya membatalkan P18, akibatnya ke spec, tiket, dan surat ikut disebut.

⭐ Ronde ini **berhasil juga bila membuktikan sesi asisten keliru.** Menemukan bahwa P18 tetap
diperlukan adalah hasil yang sama berharganya — dan jauh lebih murah diketahui sekarang daripada
sesudah surat terlanjur dibatalkan.

---

*Disusun 22 September 2026, sesudah satu berkas — `EDM Treaty In\Activity\CountSpreading_Act.xml` —
diekspor ulang dan perbandingannya menunjukkan isi langkah ternyata terbaca.*
