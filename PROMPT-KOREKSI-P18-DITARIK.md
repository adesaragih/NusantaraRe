# PROMPT KOREKSI — P18 ditarik, rantai uang dibuka

> Salin seluruh isi berkas ini sebagai prompt ke sesi eksekutor.
>
> ⚠️ **Ronde ini MENYUNTING berkas.** Berbeda dari seluruh ronde sebelumnya, yang hanya membaca
> dan menulis berkas baru. Kerjakan dengan hati-hati, dan laporkan setiap suntingan.

---

## 0. LINGKUP — DIKUNCI

`D:\XML\RNM_BRD\` adalah korpus **READ-ONLY**. Menyunting hanya di dalam `OUTPUT_HASIL_RNM\`.
`D:\XML\nusantara-re\` **terlarang**.

⛔ **`grilling-ronde-1..4.md` TERSEGEL.** Jangan disunting, walau angkanya kini terbukti keliru.
Koreksinya hidup di berkas keadaan, bukan di berkas ronde.

---

## 1. SEBABNYA

`VERIFIKASI-P18.md` memutuskan: **P18 tidak diperlukan.** Isi langkah penetapan nilai **ada** di
dalam ekspor, di tag `PropertiesName`/`PropertiesValue` — tanpa awalan `py`.

| | NB Treaty In | EDM Treaty In |
| --- | ---: | ---: |
| langkah `Property-Set` | 939 | 380 |
| punya isinya sendiri | **938** (99,9 %) | **379** (99,7 %) |
| pasangan nama=nilai terisi | 2.481 | 1.309 |
| tanda nilai terpotong, 7 uji | **0** | **0** |

Akibat kedua yang didalilkan P18 juga gugur: dari 34 medan wajib-sekaligus-terkunci, **28 punya
langkah pengisi**; tiga sisanya dapat diketik di layar lain. Layar jenjang ketiga **dapat
disimpan**.

**Penahan turun dari dua menjadi satu: tinggal P1.**

---

## 2. DAFTAR KERJA

**Sumbernya `VERIFIKASI-P18.md` Bab 6.1.** Itu daftar yang mengikat — bukan daftar di bawah ini,
yang hanya ringkasan untuk membantu Anda mengenali bentuk pekerjaannya.

| Berkas | Yang berubah |
| --- | --- |
| `SURAT-P18-KE-PEMILIK-EXPORT-PEGA.md` | dibatalkan |
| `LAMPIRAN-P18-ATURAN-DIMINTA.md` | dibatalkan |
| `spec.md` §8.2 dan AC 79 | karantina dibuka, AC dicabut |
| `spec.md` baris 51, 113, 805-807, 881, 949 | penahan 2 → 1 |
| `issues/12-layar-jenjang-ketiga.md` | `blocked` → siap dikerjakan |
| `issues/13-rantai-perhitungan-uang-dikarantina.md` | `blocked` → siap atas sisi P18 |
| `issues/00-PETA-AC.md` baris 83 | "tertahan P18" → nihil |
| `issues/07`, `issues/15` | rujukan P18 disesuaikan |
| `PERTANYAAN-untuk-pemilik-export-Pega.md` | P18 **ditarik oleh tim migrasi**, bukan dijawab |
| `PERTANYAAN-YANG-MASIH-KOSONG.md` | disamakan |
| `KEADAAN-NB-TREATY-IN.md` | pernyataan isi langkah hilang diralat |

### Cara membatalkan surat dan lampiran

**Jangan dihapus.** Ganti isinya dengan satu blok RALAT yang memuat: tanggal, sebab pembatalan,
**kutipan utuh bunyi permintaan yang ditarik**, dan kalimat tegas bahwa tidak ada yang perlu
dikirim. Bentuknya seperti ralat P41 dan P48 yang sudah ada.

### Cara menarik P18 di lembar pemilik export

P18 **ditarik oleh tim migrasi, bukan dijawab pemilik export**. Bedanya penting dan harus terbaca:
tidak seorang pun di luar tim yang pernah menjawabnya. Sebutkan bahwa kekeliruannya milik tim
migrasi, dan sebutkan apa penyebabnya — dua huruf `py`.

---

## 3. DUA PEKERJAAN TAMBAHAN

### 3.1 Tutup P60 — EDM Treaty In

`P60` menanyakan apa yang dikerjakan satu langkah tambahan pada `CountSpreading_Act` versi EDM.
Karena isinya kini terbaca, **jawablah dari korpus**, jangan ditanyakan kepada siapa pun.

Yang sudah terbaca — **verifikasi ulang sendiri, jangan disalin percaya**:

```
langkah tambahan EDM :  .SpreadingRiskList(1).SplitRNMSharePct = 100

bila SharePercentage kosong
   NB  :  100 / jumlah baris spreading            presisi 10
   EDM :  SplitRNMSharePct / RNMShare             presisi 20
```

Jawab juga bagian kedua pertanyaannya: **kenapa endorsemen memerlukannya sementara polis baru
tidak.** Bila tidak terbaca dari korpus, katakan `[terbuka]` — jangan dikarang.

Catat di `.scratch\edm-treaty-in\PERTANYAAN-RONDE-1.md`.

### 3.2 Sensus `DataTransform` setara

`VERIFIKASI-P18.md` mencatat bahwa `DataTransform` memakai **nama tag kebalikannya** dan belum
disensus setara. Kerjakan sekarang, jangan dibiarkan.

Pertanyaannya sama: **apakah ada lagi yang kita kira hilang padahal ada?**

Sensus untuk `DataTransform` dan `FlowAction`, kedua modul, dua cara. Bila ditemukan pola serupa,
sebutkan berkas mana dan berapa banyak — dan **apakah ada butir `[terbuka]` lain yang ikut gugur
karenanya**.

---

## 4. YANG TIDAK BOLEH DIKERJAKAN

- ⛔ Jangan menyunting `grilling-ronde-*.md`. Tersegel.
- ⛔ Jangan menyentuh **P1**. Ia tetap menahan, dan suratnya tetap dikirim.
- ⛔ Jangan menyentuh **P8** bila masih menahan sebagian tiket 13 — periksa, jangan asal dibuka.
- ⛔ Jangan membuat keputusan baru, jangan membuat ADR baru, jangan menutup butir `[terbuka]` lain
  yang tidak disebut di sini.
- ⛔ Jangan menulis kode, DDL, atau usulan kolom.

---

## 5. SESUDAH MENYUNTING — WAJIB DIPERIKSA ULANG

| Uji | Yang diharapkan |
| --- | --- |
| cacah pertanyaan NB, dua cara | 48 dari 49 terjawab; tinggal **P1** |
| cacah pertanyaan EDM | **11 dari 11** terjawab |
| sinkronisasi tiga salinan lembar pertanyaan | nol nomor berbeda isi |
| tiket berstatus `blocked` | **nihil** karena P18; sebutkan bila ada yang masih blocked karena sebab lain |
| `00-PETA-AC.md` | ke-96 AC tetap tertutup sesudah AC 79 dicabut — **cacah ulang dua cara** |
| kata "P18" di seluruh `OUTPUT_HASIL_RNM\` | hanya tersisa di berkas tersegel dan di blok RALAT |

⚠️ **AC 79 dicabut berarti jumlah AC berubah.** Pastikan blok sensus di `spec.md` §1.3 ikut
diperbarui, dan **kutip angka lamanya** — jangan dihapus. Itu aturan rumah.

---

## 6. KELUARAN

`.scratch\nb-treaty-in\KOREKSI-P18-DIJALANKAN.md`, lima bab.

| Bab | Isi |
| --- | --- |
| 1 | daftar suntingan — berkas, baris, bunyi lama, bunyi baru |
| 2 | P60 dijawab, dengan verifikasi ulangnya |
| 3 | sensus `DataTransform` dan `FlowAction` |
| 4 | hasil enam pemeriksaan Bab 5 |
| 5 | **TELEMETRI EKSEKUSI** |

Telemetri diukur dari luar bila bisa:

```
claude --print --output-format json "<prompt>" > hasil.json
```

Bila tidak, katakan begitu. Bila hasil pengurangan baseline, katakan itu juga.

---

## 7. YANG MENANDAKAN RONDE INI BERHASIL

- Seluruh berkas pada Bab 6.1 `VERIFIKASI-P18.md` disunting, dan **tidak ada yang lain**.
- Pembatalan surat berbentuk **RALAT dengan kutipan utuh**, bukan penghapusan.
- P18 tercatat **ditarik oleh tim migrasi**, dengan sebabnya tertulis.
- P60 terjawab dari korpus, diverifikasi ulang, bukan disalin dari prompt ini.
- `DataTransform` disensus, dan pertanyaan *"adakah lagi yang kita kira hilang"* terjawab.
- Keenam pemeriksaan Bab 5 dijalankan, hasilnya ditulis apa adanya.
- Nol keputusan baru, nol ADR baru, nol berkas tersegel tersentuh.

---

*Disusun 22 September 2026, sesudah `VERIFIKASI-P18.md` memutuskan bahwa penahan terberat proyek
ini tidak pernah ada — dan bahwa sebabnya dua huruf.*
