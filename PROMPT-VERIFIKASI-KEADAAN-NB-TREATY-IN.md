# PROMPT VERIFIKASI — Keadaan NB Treaty In

> Salin seluruh isi berkas ini sebagai prompt ke sesi eksekutor.
> Ronde ini **hanya memeriksa**. Tidak menulis spec, tidak menemukan hal baru.
> Jalankan ini **sebelum** `PROMPT-TO-SPEC-NB-TREATY-IN.md`.

---

## 0. LINGKUP — DIKUNCI

Hanya modul **`D:\XML\RNM_BRD\NB Treaty In`**.

`D:\XML\RNM_BRD\` adalah korpus **READ-ONLY**. Menulis hanya ke `OUTPUT_HASIL_RNM\`.
`D:\XML\nusantara-re\` **terlarang** — jangan dibaca, dikutip, dibandingkan, atau dijadikan sasaran.

Keluaran satu berkas: **`OUTPUT_HASIL_RNM\.scratch\nb-treaty-in\VERIFIKASI-KEADAAN.md`**.

⛔ Jangan menyunting `KEADAAN-NB-TREATY-IN.md`. Laporkan selisihnya; perbaikannya dilakukan
sesudah laporan dibaca manusia.

⛔ Jangan menyunting `grilling-ronde-1..4.md`. Berkas itu tersegel.

---

## 1. YANG DIPERIKSA DAN KENAPA

Berkas **`.scratch\nb-treaty-in\KEADAAN-NB-TREATY-IN.md`** menyatakan keadaan modul yang berlaku
sekarang. Ia akan menjadi dasar penulisan spec.

Berkas itu disusun oleh sesi asisten dan **belum diperiksa siapa pun**. Penyusunnya melakukan
**dua kesalahan yang diketahui** pada hari yang sama:

1. menyatakan naskah SQL tidak ada di korpus, padahal ada 1.151 pernyataan di tag `<pyBrowseSQL>`;
2. mengulang pernyataan "lapisan layar nol dibuka" tanpa memeriksa silang ke ronde 2, yang sudah
   menyisirnya.

Jangan percaya berkas itu begitu saja. Bab 7 di dalamnya memuat kode Python siap jalan.

---

## 2. SEPULUH PERNYATAAN YANG WAJIB DIUJI

| # | Pernyataan | Angka yang ditulis |
| ---: | --- | ---: |
| 1 | berkas dalam modul | 278 |
| 2 | ukuran modul | 35.492.317 B |
| 3 | langkah `Property-Set` di 92 `Activity` | 938 |
| 4 | langkah `Page-Clear-Messages` | 16 |
| 5 | baris `DecisionTable\BusinessType_DeT` | 36 |
| 6 | `Activity` terjangkau / yatim | 64 / 28 |
| 7 | langkah `Property-Set` terjangkau / yatim | 377 / 561 |
| 8 | medan wajib (`pyRequired` = true) | 27 medan di 6 layar |
| 9 | medan terkunci permanen | 38 — 36 di layar Dept Head, 2 di layar biasa |
| 10 | kolom view `POOLDATA.TREATYINDETAILJOINEDM` | 39; 33 dari 33 medan laporan tersedia |

**Nomor 3, 6, 7, dan 8 dihitung dua cara** yang benar-benar berbeda — bukan skrip yang sama
dijalankan dua kali. Bila hasilnya berbeda, tulis keduanya dan katakan mana yang dipercaya.

### Yang paling perlu diragukan — nomor 6 dan 7

Penjangkauan disusun dengan cara **mencari nama aturan sebagai teks** di dalam berkas lain, bukan
mengurai langkah `Call`. Itu **batas atas**: daftar "terjangkau" mungkin memuat aturan yang
sebenarnya tidak pernah dipanggil.

Bila Anda punya cara yang lebih tepat — mengurai `pyStepsActivityName` yang berawalan `Call`, lalu
memetakan nama sasarannya — **pakai itu**, dan laporkan selisihnya terhadap 64/28 dan 377/561.

Titik masuk yang dipakai penyusun: aturan yang dirujuk dari `Flow`, `Section`, `Harness`,
`FlowAction`, atau `DataTransform`, lalu diikuti sampai tidak ada aturan baru yang tercapai.

---

## 3. TIGA PEMERIKSAAN TAMBAHAN

### 3.1 Pertentangan antar sumber

Urutan kewenangan:

| # | Sumber | Sifat |
| ---: | --- | --- |
| 1 | `PERTANYAAN-untuk-*.md` — enam lembar jawaban | keputusan work owner, mengikat |
| 2 | `KEADAAN-NB-TREATY-IN.md` | keadaan terukur |
| 3 | `grilling-ronde-1..4.md` | latar, tersegel |

**Daftarkan setiap pertentangan yang Anda temukan** antara sumber 1 dan 2. Pertentangan antara 3
dan yang lain sudah diketahui dan tercatat di Bab 5 berkas keadaan — periksa apakah daftar itu
lengkap, dan tambahkan yang terlewat.

### 3.2 Kelengkapan lembar jawaban

Cacah: berapa pertanyaan seluruhnya, berapa terjawab, berapa kosong. **Dua cara** — lewat judul
`##` dan lewat slot `> *(tulis di sini)*`.

Enam lembar pemilik adalah sumbernya. `PERTANYAAN-RONDE-4.md` dan `PERTANYAAN-YANG-MASIH-KOSONG.md`
adalah **salinan** — periksa apakah ketiganya sinkron, dan laporkan nomor yang berbeda isi.

### 3.3 Berkas korpus yang mencurigakan

Pada 22 September 2026 ditemukan satu berkas salah di dalam folder modul: `Protection_Act` versi
Fac In (kelas `ASM-FW-GISFW-Work`, 343.271 B) menggantikan versi Treaty yang benar (kelas
`ASM-FW-GISFW-Data-PolicyTreatyIn`, 157.304 B). Berkas itu sudah diperbaiki.

Penelusuran pola yang sama menemukan nol kasus lain, dengan cara: aturan di folder Treaty yang
memakai kelas berbagi, padahal modul Treaty lain punya varian berkelas Treaty.

**Ujilah dengan cara lain.** Misalnya: bandingkan md5 setiap berkas NB Treaty In dengan berkas
bernama sama di modul Fac, lalu periksa mana yang seharusnya punya varian Treaty tetapi tidak
punya. Laporkan kandidat yang Anda temukan, walau ragu.

---

## 4. KELUARAN

Satu berkas: `.scratch\nb-treaty-in\VERIFIKASI-KEADAAN.md`.

### Bab 1 — hasil sepuluh pernyataan

| # | Pernyataan | Ditulis | Terukur | Cocok? | Cara |
| ---: | --- | ---: | ---: | --- | --- |

Kolom **Cara** memuat perintah yang Anda jalankan, supaya orang lain dapat mengulanginya.

### Bab 2 — pertentangan antar sumber

Setiap pertentangan: sumber mana, bunyi masing-masing, dan mana yang menang menurut urutan
kewenangan.

### Bab 3 — kelengkapan dan sinkronisasi lembar jawaban

### Bab 4 — kandidat berkas korpus yang mencurigakan

Boleh kosong. Bila kosong, sebutkan cara apa saja yang sudah dijalankan untuk mencarinya.

### Bab 5 — putusan

Satu kalimat: **berkas keadaan layak dipakai sebagai dasar spec, atau tidak**, dan bila tidak,
apa yang harus diperbaiki lebih dulu.

### Bab 6 — TELEMETRI EKSEKUSI

Wajib. Angka token sejati tidak terlihat dari dalam sesi; ukur dari luar:

```
claude --print --output-format json "<prompt>" > hasil.json
```

Catat token keluaran · cache-read · jumlah panggilan alat · durasi · biaya · byte dibaca.
Bila pengukuran dari luar tidak dilakukan, **katakan begitu** — jangan menaksir lalu menyajikannya
sebagai angka terukur.

---

## 5. DISIPLIN

Setiap angka membawa **perintah yang menghasilkannya**.

Penanda: `[terverifikasi]` · `[dugaan]` · `[terbuka]`.

**Jangan menutup pertanyaan terbuka.** **Jangan menyalin nilai berupa nama orang.** Nol rahasia,
nol data nasabah, nol nomor polis apa adanya.

Jalur korpus memuat spasi — pakai Python, bukan loop shell.

### Empat jebakan yang sudah pernah menjerat

1. `<rowdata REPEATINGINDEX="n"/>` yang menutup sendiri tidak tertangkap pola
   `<rowdata...>(.*?)</rowdata>` — sel kosong hilang, pasangan bergeser.
2. `pyRowNum` pada `pyOrConditions` berbasis **nol**; baris tabel berbasis satu.
3. `pyStepsPage`, `pyCriteriaValue`, `pyResult`, `pyShapeName`, `pyTaskLabel` **tidak ada** di
   ekspor ini. Periksa daftar tag lebih dulu.
4. Menyatakan sesuatu nihil tanpa menyebut lingkup penelusuran.

---

## 6. YANG MENANDAKAN RONDE INI BERHASIL

- Sepuluh pernyataan diuji, hasilnya ditulis apa adanya — **termasuk yang tidak cocok**.
- Nomor 3, 6, 7, 8 dihitung dua cara, dan kedua caranya dijelaskan.
- Penjangkauan diuji dengan cara yang berbeda dari penyusun, atau dikatakan terang-terangan bila
  caranya sama.
- Bab 5 memuat putusan yang tegas, bukan ringkasan.
- Bab telemetri terisi angka terukur, bukan taksiran.

Ronde ini **berhasil juga bila menemukan banyak kesalahan**. Menemukan nol kesalahan bukan
ukuran keberhasilan.

---

*Disusun 22 September 2026. Ronde berikutnya: `PROMPT-TO-SPEC-NB-TREATY-IN.md`.*
