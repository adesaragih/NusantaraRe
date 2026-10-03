# Spec — NB Treaty In
## Realisasi penutupan treaty inward: dari penawaran yang disetujui menjadi polis treaty

Modul: **`NB Treaty In`**. Korpus `D:\XML\RNM_BRD\` **READ-ONLY**.
Disusun **22 September 2026**, sesudah empat ronde grilling, satu ronde verifikasi, dan
**47 dari 49 pertanyaan terjawab**.

---

## 1 · Cara membaca berkas ini

### 1.1 Urutan kewenangan

| # | Sumber | Sifat |
| ---: | --- | --- |
| 1 | `PERTANYAAN-untuk-*.md` — enam lembar jawaban | ⭐ **keputusan work owner, mengikat** |
| 2 | `VERIFIKASI-KEADAAN.md` | hasil pemeriksaan; menang atas 3 |
| 3 | `KEADAAN-NB-TREATY-IN.md` | keadaan terukur |
| 4 | `grilling-ronde-1..4.md` | ⛔ **latar; angkanya sebagian basi** |
| 5 | korpus XML | selalu boleh dipakai membuktikan ulang |

### 1.2 Penanda

| Penanda | Artinya |
| --- | --- |
| `[terverifikasi]` | ada **jalur berkas + tipe rule + nama rule**, dan perintah auditnya |
| `[keputusan work owner]` | diputuskan orang yang berwenang — ⛔ **tidak boleh diubah spec** |
| `[penyimpangan sadar]` | ⚠️ **sengaja berbeda dari Pega**, alasannya tertulis |
| `[terbuka]` | belum terjawab — ⛔ **tidak ditutup oleh spec ini** |
| `[data DBA]` | hanya dapat dijawab dari basis data |
| `[dugaan]` | ⛔ **tidak dipakai sebagai dasar keputusan mana pun** |

### 1.3 ⭐ Sensus berkas ini — dihitung DUA CARA

> ⭐ **Jendela: seluruh berkas ini MINUS blok sensus 1.3 sendiri.** — **1.032 baris**.
>
> | Yang dicacah | Cara A | Cara B | |
> | --- | ---: | ---: | :---: |
> | User story | **48** | **48** | ✅ |
> | Acceptance Criteria | **96** | **96** | ✅ |
> | ⚠️ di antaranya **dicabut lalu DIGANTI** *(79, 80)* | **2** | **2** | ✅ | nomor **tetap dipakai** |
> | ⛔ AC **tanpa penanda** | **0** | **0** | ✅ |
> | ⛔ AC **tanpa rujukan bab** | **0** | **0** | ✅ |
> | butir `[terbuka]` | **23** | **23** | ✅ | ⚠️ **RALAT** — semula **21**, lalu **20** *(P18 ditarik)*, lalu **22**; kini **23** — P1 tertutup, **empat** butir baru masuk §9.2 |
> | butir `[penyimpangan sadar]` | **6** | **6** | ✅ |
>
> ⭐ **Cara A** — cacah pola penomoran di dalam babnya *(`^N. **Sebagai**` bab 4; `` ^N. `[ `` bab 7;
> baris tabel bab 9 dan 10.1)*.
> ⭐ **Cara B** — periksa **keberurutan nomornya**: user story **1–48 tanpa lompat**, AC
> **1–96 tanpa lompat** *(79 dan 80 dicabut lalu diganti di nomor yang sama)*; butir terbuka
> dipecah **9.1 + 9.2 = 0 + 23**;
> penyimpangan bernomor **1–6**.
>
> ⭐ **Butir `[terbuka]` 23 itu terbagi:** ⭐⭐ **0 MENAHAN** · **23 tidak menahan**.
>
> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — kalimat ini semula berbunyi:
> > *"⭐ **Butir `[terbuka]` 20 itu terbagi:** ⛔ **1 MENAHAN** *(P1)* · **19 tidak menahan**."*
>
> **P1 terjawab.** Tiga butir baru ditambahkan di §9.2 — **P61**, **P62**, dan **sensus properti**.
>
> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — kalimat ini semula berbunyi:
> > *"⭐ **Butir `[terbuka]` 21 itu terbagi:** ⛔ **2 MENAHAN** *(P1, P18)* · **19 tidak menahan**."*
>
> **P18 ditarik oleh tim migrasi** — bukan dijawab pemilik export. Isi langkah ada di dalam
> ekspor, di tag `PropertiesName`/`PropertiesValue` **tanpa awalan `py`**. Lihat `VERIFIKASI-P18.md`.
>
> **Sebaran penanda, kemunculan di jendela:** `[keputusan work owner]` **87** ·
> `[terverifikasi]` **52** · `[terbuka]` **12** · `[penyimpangan sadar]` **11** ·
> `[data DBA]` **5** · `[dugaan]` **2**.
> ⚠️ **Kemunculan penanda ≠ jumlah butir** — mis. `[penyimpangan sadar]` disebut **11 kali** untuk
> **6 penyimpangan**, sebab ia diulang di AC yang menguji penyimpangan itu.
>
> ✅ **Nol nama orang · nol nomor polis apa adanya · nol `CREATE TABLE`** — disisir, ketiganya **0**.
>
> > ⛔⛔ **RALAT — ditangkap sensus ini sendiri, sebelum berkas dilaporkan.** Rancangan pertama
> > blok ini menulis angka **sebelum** mengukurnya. Kalimatnya dikutip utuh, tidak dihapus:
> >
> > *"User story **48** · Acceptance Criteria **96** · butir `[terbuka]` aktif **19** · butir
> > `[penyimpangan sadar]` **6** … Sebaran penanda: `[terverifikasi]` **64** ·
> > `[keputusan work owner]` **41** · `[terbuka]` **19** · `[penyimpangan sadar]` **6** ·
> > `[data DBA]` **5** · `[dugaan]` **3**."*
> >
> > ⭐ **Yang benar:** user story **48** ✅ dan AC **96** ✅ sudah tepat; ⛔ tetapi butir `[terbuka]`
> > **21**, bukan 19 — saya menghitung bab 9.2 saja dan melewatkan **dua penahan** di 9.1.
> > ⛔ Dan **empat dari enam angka sebaran penanda meleset**: `[terverifikasi]` **52** bukan 64,
> > `[keputusan work owner]` **87** bukan 41, `[terbuka]` **12** bukan 19, `[penyimpangan sadar]`
> > **11** bukan 6.
> >
> > ⚠️ **Sebab kekeliruannya, dan ia berulang:** saya mencampur **jumlah butir** dengan
> > **kemunculan penanda**. ⭐ Keduanya kini dipisahkan tegas di atas.
>
> ⚠️ **Rancangan pertama juga menulis** *"User story **45** · Acceptance Criteria **90**"* pada
> tahap lebih awal lagi. ⭐ Terukur **48** dan **96**.

### 1.4 Bacalah bab 8 dan 9 sebelum membangun

⛔ **Lingkup modul ini menyusut 60 %** sesudah satu berkas korpus yang salah diperbaiki. Bab 8
menyebut apa yang **tidak** dibangun; bab 9 menyebut apa yang **belum** dapat dibangun.
⛔ Membangun dari bab 4–7 saja akan membangun hal yang sudah dinyatakan di luar lingkup.

---

## 2 · Problem Statement

Nusantara Re menerima penawaran treaty inward dari ceding company. Ketika sebuah penawaran
disetujui, ia harus **direalisasikan** menjadi polis treaty: datanya dilengkapi, diperiksa
berjenjang, diberi nomor polis, dan disimpan sebagai kontrak yang mengikat.

Hari ini pekerjaan itu berjalan di atas Pega, dan **empat hal membuatnya mahal dan rapuh**:

1. ⛔ **Data kontrak disimpan sebagai satu dokumen teks.** Seluruh halaman kerja diubah menjadi
   JSON apa adanya lalu ditulis ke satu kolom. `[terverifikasi]` Fungsi pembentuknya tidak memilih
   apa pun — ia memotret halaman pada saat penyimpanan. ⛔ Akibatnya **bentuk data tidak dapat
   dinyatakan**, dan tidak ada laporan yang dapat dibuat di atasnya tanpa membongkar teks.
2. ⛔ **Peran pengguna disimpan di kolom nomor telepon.** `[terverifikasi]` Kode seperti `TREATY1`
   dibandingkan terhadap `OperatorID.pyTelephone`. ⛔ Tidak ada model peran; nol aturan otorisasi
   di seluruh 9.430 berkas korpus.
3. ⛔ **Identitas orang ditanam langsung di dalam aturan.** `[terverifikasi]` **12 berkas** modul ini
   memeriksa "apakah pengguna ini orang tertentu". ⛔ Aturan seperti itu berhenti bekerja saat
   orangnya pindah, tanpa pemberitahuan.
4. ⛔ **Penyimpanan tidak punya jaminan keutuhan.** `[data DBA]` Dua dari tiga stored procedure
   melakukan `COMMIT` sendiri, dan blok pemanggilnya `COMMIT` lagi. ⛔ Kegagalan pada langkah
   berikutnya meninggalkan data yang sudah permanen.

⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — paragraf ini semula berbunyi:
> *"⚠️ **Dan satu hal yang bukan masalah Pega, melainkan masalah bahan:** isi **268 langkah
> penetapan nilai** — tempat seluruh aritmetika uang tinggal — ⛔ **tidak ikut dalam ekspor** yang
> diterima tim migrasi. Itu **P18**, dan ia menahan sebagian modul ini."*

⭐⭐ **Kalimat itu keliru, dan kekeliruannya milik tim migrasi.** Isi 268 langkah penetapan nilai
**ada di dalam ekspor**, di tag `PropertiesName` dan `PropertiesValue` — **tanpa awalan `py`** — di
dalam `pyParamArray` milik tiap langkah. Terukur: **269 langkah** pada ke-51 aturan terjangkau,
**707 pasangan nama=nilai terisi**, **nol** tanda terpotong. **P18 ditarik.** Rinciannya di
`VERIFIKASI-P18.md`.

---

## 3 · Solution

Realisasi treaty inward dibangun ulang di Go + React + Oracle, dengan **lima perubahan pokok**
terhadap sistem lama — ⭐ seluruhnya `[keputusan work owner]`:

| # | Perubahan | Dari | Menjadi |
| ---: | --- | --- | --- |
| ⭐ **1** | **Sumber data** | kolom dokumen `JSONDATA` | ⭐ **view relasional `POOLDATA.TREATYINDETAILJOINEDM`** — 39 kolom |
| ⭐ **2** | **Tangga persetujuan** | lima posisi | ⭐ **tiga** — Admin → Sec Head → Dept Head |
| ⭐ **3** | **Peran** | kode di kolom nomor telepon | ⭐ **medan peran tersendiri** |
| ⭐ **4** | **Keutuhan penyimpanan** | commit ganda, tak dapat dibatalkan | ⭐ **satu transaksi**, gagal ⇒ batal seluruhnya |
| ⭐ **5** | **Tanggal** | teks, dua format tercampur | ⭐ **tipe tanggal Oracle**, satu format |

⭐ **Yang TIDAK berubah, dan itu disengaja:** aturan bisnisnya. Nilai penanda persetujuan,
penggolongan jenis usaha, rumus pajak brokerage, pengecualian mata uang, dan antrean bersama
**ditiru apa adanya** — termasuk ketidakseragamannya.

---

## 4 · User Stories

### Realisasi dan daur hidup berkas

1. **Sebagai** admin treaty, **saya ingin** membuka penawaran yang sudah disetujui dan melengkapi
   datanya, **supaya** ia dapat direalisasikan menjadi polis treaty.
2. **Sebagai** admin treaty, **saya ingin** sistem mengisi tanggal mulai, tanggal akhir, dan
   tanggal laporan secara otomatis bila saya biarkan kosong, **supaya** saya tidak perlu mengetik
   tanggal hari ini berulang kali.
3. **Sebagai** admin treaty, **saya ingin** berkas yang saya kerjakan menunggu di antrean bersama,
   **supaya** rekan sejawat dapat melanjutkannya ketika saya tidak ada.
4. **Sebagai** admin treaty, **saya ingin** diberi peringatan ketika saya hendak membuat penawaran
   yang sudah pernah ada, **supaya** tidak terjadi penawaran ganda.
5. **Sebagai** admin treaty, **saya ingin** nomor polis terbentuk otomatis dengan pola yang tetap,
   **supaya** penomorannya konsisten dan tidak bentrok.
6. **Sebagai** admin treaty, **saya ingin** menyimpan berkas hanya bila seluruh medan wajib terisi,
   **supaya** kontrak tidak tersimpan separuh jadi.

### Tangga persetujuan

7. **Sebagai** admin treaty, **saya ingin** mengajukan berkas ke Kepala Seksi, **supaya** ia
   diperiksa sebelum mengikat.
8. **Sebagai** Kepala Seksi, **saya ingin** melihat berkas yang menunggu keputusan saya di antrean
   bersama, **supaya** siapa pun yang memegang posisi itu dapat menanganinya.
9. **Sebagai** Kepala Seksi, **saya ingin** menyetujui berkas dan meneruskannya ke Kepala
   Departemen, **supaya** jenjang berikutnya dapat memutuskan.
10. **Sebagai** Kepala Seksi, **saya ingin** menolak berkas dan mengembalikannya ke admin,
    **supaya** kekeliruannya dapat diperbaiki dan diajukan ulang.
11. **Sebagai** Kepala Departemen, **saya ingin** mengisi tujuh medan keputusan saya sendiri,
    **supaya** keputusan saya terekam beserta angkanya.
12. **Sebagai** Kepala Departemen, **saya ingin** menyetujui berkas sebagai keputusan terakhir,
    **supaya** realisasi selesai dan polis terbentuk.
13. **Sebagai** Kepala Departemen, **saya ingin** menolak berkas dan mengembalikannya ke admin,
    **supaya** ia diperbaiki, bukan dibuang.
14. **Sebagai** admin treaty, **saya ingin** berkas yang **saya sendiri** tolak langsung
    diselesaikan sebagai ditolak, **supaya** tidak beredar tanpa tujuan.
15. **Sebagai** siapa pun di tangga, **saya ingin** melihat riwayat siapa memutuskan apa dan kapan,
    **supaya** keputusan dapat ditelusuri.
16. **Sebagai** siapa pun di tangga, **saya ingin** menambahkan catatan pada berkas, **supaya**
    alasan keputusan tersimpan bersama keputusannya.

### Wewenang

17. **Sebagai** pemilik sistem, **saya ingin** wewenang ditentukan oleh **peran**, bukan oleh nama
    orang, **supaya** kepergian seseorang tidak mematahkan alur.
18. **Sebagai** pemilik sistem, **saya ingin** peran disimpan di medan peran, bukan di kolom nomor
    telepon, **supaya** data pengguna berarti apa yang namanya katakan.
19. **Sebagai** pemilik sistem, **saya ingin** pemeriksaan wewenang bertanya *"apakah pengguna ini
    anggota antrean X"*, bukan *"antrean nomor dua Anda namanya apa"*, **supaya** perubahan urutan
    daftar tidak mengubah wewenang.
20. **Sebagai** auditor, **saya ingin** setiap tindakan mencatat **identitas akses login**-nya,
    **supaya** jejaknya tidak bergantung pada nama tampilan yang dapat berubah.

### Data kontrak

21. **Sebagai** pengguna, **saya ingin** melihat batas, retensi, premi bersih, dan pendapatan premi
    diperkirakan beserta mata uangnya, **supaya** saya tahu besaran kontraknya.
22. **Sebagai** pengguna, **saya ingin** angka uang ditampilkan apa adanya seperti tersimpan,
    **supaya** layar dan laporan tidak pernah berselisih.
23. **Sebagai** pengguna, **saya ingin** melihat penempatan keluar (retro) dari layar treaty masuk,
    **supaya** saya melihat gambaran utuh penempatan risikonya.
24. **Sebagai** pemilik sistem, **saya ingin** data treaty keluar **hanya dibaca** dari konteks ini,
    **supaya** batas tanggung jawab antar modul tetap jelas.
25. **Sebagai** pengguna, **saya ingin** jenis usaha digolongkan otomatis dari kode bisnisnya,
    **supaya** saya tidak perlu memilihnya manual.
26. **Sebagai** pengguna, **saya ingin** penggolongan yang tidak dikenali diberi nilai bawaan yang
    jelas, **supaya** berkas tetap dapat diproses.
27. **Sebagai** pengguna, **saya ingin** daftar pilihan mata uang tidak memuat mata uang mati,
    **supaya** saya tidak salah pilih.

### Layar

28. **Sebagai** pengguna, **saya ingin** medan yang tidak boleh saya ubah tampil terkunci,
    **supaya** saya tahu batas kewenangan saya tanpa mencoba.
29. **Sebagai** pengguna, **saya ingin** medan wajib ditandai sesuai tingkat saya, **supaya** saya
    tidak diminta mengisi hal yang bukan urusan saya.
30. **Sebagai** Kepala Departemen, **saya ingin** layar saya menuntut medan yang berbeda dari layar
    admin, **supaya** keputusan saya terekam lengkap tanpa mengulang pekerjaan admin.
31. **Sebagai** pengguna, **saya ingin** bagian layar yang tidak berlaku bagi jenis kontrak saya
    disembunyikan, **supaya** layarnya tidak membingungkan.
32. **Sebagai** pengguna, **saya ingin** melihat kode jenis kontrak apa adanya bila keterangannya
    belum tersedia, **supaya** saya tetap dapat bekerja.

### Penyimpanan dan keutuhan

33. **Sebagai** pemilik sistem, **saya ingin** seluruh urutan penyimpanan dibungkus **satu
    transaksi**, **supaya** kegagalan di langkah mana pun tidak meninggalkan data separuh.
34. **Sebagai** pemilik sistem, **saya ingin** setiap query menyebut skema secara eksplisit,
    **supaya** perilakunya sama di lingkungan uji dan produksi.
35. **Sebagai** pemilik sistem, **saya ingin** tanggal disimpan sebagai tipe tanggal, bukan teks,
    **supaya** urutan dan perbandingannya benar.
36. **Sebagai** pemilik sistem, **saya ingin** nilai uang disimpan berpresisi penuh, **supaya**
    pembulatan tidak merambat ke laporan.
37. **Sebagai** pengguna, **saya ingin** diberi tahu ketika data kontrak gagal dibaca, **supaya**
    saya tidak bekerja di atas data separuh tanpa menyadarinya.

### Pajak dan potongan

38. **Sebagai** petugas keuangan, **saya ingin** potongan brokerage dihitung dengan memperhitungkan
    pajak bila jenis pajaknya *inclusive*, **supaya** angkanya sama seperti sistem lama.
39. **Sebagai** petugas keuangan, **saya ingin** persentase dibaca sebagai persentase, bukan uang,
    **supaya** `12.5` berarti 12,5 persen.

### Efek keluar dan riwayat

40. **Sebagai** pemilik sistem, **saya ingin** setiap perpindahan tahap tercatat di tabel riwayat,
    **supaya** alurnya dapat diaudit.
41. **Sebagai** pemilik sistem, **saya ingin** catatan riwayat memisahkan **identitas akses** dari
    **nama tampilan**, **supaya** keduanya tidak tertukar.
42. **Sebagai** pengguna, **saya ingin** pemberitahuan menyebut nama orang yang **diambil dari
    data**, **supaya** pesannya tetap benar ketika orangnya berganti.

### Migrasi

43. **Sebagai** pemilik sistem, **saya ingin** data lama terbaca di sistem baru, **supaya** riwayat
    kontrak tidak hilang.
44. **Sebagai** pemilik sistem, **saya ingin** tahu berkas mana yang tanggalnya ambigu sebelum
    migrasi, **supaya** saya dapat memutuskan penanganannya.
45. **Sebagai** pemilik sistem, **saya ingin** aturan yang tidak pernah dipanggil **tidak**
    dimigrasi, **supaya** sistem baru tidak mewarisi beban mati.
46. **Sebagai** pemilik sistem, **saya ingin** rantai perhitungan uang yang bahannya belum lengkap
    **dikarantina**, bukan ditebak, **supaya** tidak ada angka yang dibangun dari tebakan.
47. **Sebagai** pemilik sistem, **saya ingin** setiap penyimpangan dari perilaku lama tercatat
    beserta alasannya, **supaya** keputusannya dapat ditinjau ulang.
48. **Sebagai** pemilik sistem, **saya ingin** setiap butir yang masih terbuka tertulis di satu
    tempat, **supaya** tidak ada yang terlupakan saat go-live.

---

## 5 · Implementation Decisions

### 5.1 ⭐ Sumber data pindah ke view relasional

`[keputusan work owner]` `[data DBA]` **P29, P15, P42**

⛔ **JSON tidak dipakai di sistem baru** — tidak untuk membaca, tidak untuk menulis. Seluruh
mekanisme `adoptJSONObject` + pembacaan kolom dokumen **tidak dimigrasi**.

⭐ Sebagai gantinya, data ditarik dari view relasional yang sudah ada:

```
POOLDATA.TREATYINDETAILJOINEDM
   = pooldata.treatyindetail  UNION ALL  pooldata.treatyindetailedm
```

**39 kolom**, berkelompok: pengenal · kontrak · penggolongan · **uang berpasangan mata uang** ·
potongan · lapisan.

⭐ **Uji kecukupan sudah dilakukan** `[terverifikasi]`: `ReportDefinition\BrowseTreatyInDetail.xml`
merujuk **33 medan**, dan **ke-33-nya ada di view**. ⛔ Nol medan yang dipakai tetapi tidak
tersedia.

⚠️ **Kenapa ini menyelesaikan masalah pokok:** `[terverifikasi]` naskah
`@ASM.GetPageJSONString()` — kini diterima — ternyata **tidak memilih apa pun**; ia memotret
seluruh halaman langkah. ⛔ Jadi isi kolom dokumen **bukan skema**, melainkan potret. ⭐ View
relasional memberi bentuk yang dapat dinyatakan.

### 5.2 ⭐ Tangga persetujuan menyusut menjadi tiga

`[keputusan work owner]` **P13**

```
ReasTreatyInAdmin  →  ReasTreatyInSecHead  →  ReasTreatyInDeptHead
```

⛔ **Dua posisi DIBUANG:** `ReasTreatyInGroupLeader` dan `ReasTreatyInDirector`.
⛔ Jalur klaim → Manajer Klaim → Direktur **tidak dibangun** *(P5)*.

⚠️ **Angka yang menyertai keputusan, supaya dapat ditinjau ulang** `[terverifikasi]`: kedua posisi
yang dibuang punya jejak **sama banyak** dengan Dept Head yang dipertahankan — dua berkas.
⛔ Korpus **tidak mendukung dan tidak membantah** pembuangan itu; ia sepenuhnya keputusan bisnis.

⭐ **Antrean bersama dipertahankan** `[terverifikasi]` **P10** — keenam Assignment memakai
`ToWorkBasket`, **nol** `ToWorklist`. ⛔ Tidak ada kotak masuk pribadi.

### 5.3 ⭐ Penanda persetujuan dan empat cabang putusan

`[keputusan work owner]` **P24, P6**

| Nilai `IsApproved` | Artinya |
| --- | --- |
| ⭐ **`0`** | **ditolak** |
| ⭐ **selain itu**, termasuk kosong | **disetujui** |

⛔ **Dibandingkan sebagai TEKS** `[terverifikasi]` — kolom tabel keputusan ber-`pyColumnDataType`
= `text`.

⭐ **Aturan hidup ada di `DecisionTable\isApproved.xml`**, ⛔ **bukan** `When\isApproved.xml`.
`[terverifikasi]` Keenam kotak Decision di `Flow\InputRealizationTreatyIn.xml` menyambung ke
`Rule-Declare-DecisionTable`; **nol** menyambung ke rule `When`.
⚠️ **Perbedaannya nyata:** rule `When` memperlakukan nilai kosong sebagai **tidak disetujui**;
tabel keputusan memperlakukannya sebagai **disetujui**. ⛔ Membaca yang salah membalik perilaku
pada setiap berkas yang penandanya belum pernah diisi.

⭐ **Empat cabang putusan:**

| Siapa | Putusan | Akibat |
| --- | --- | --- |
| Admin | terima | naik ke Sec Head |
| ⭐ **Admin** | **tolak** | ⭐ **berkas langsung diselesaikan sebagai ditolak** |
| Atasan *(Sec Head, Dept Head)* | terima | naik |
| ⭐ **Atasan** | **tolak** | ⭐ **kembali ke admin** |

### 5.4 ⭐ Peran menggantikan nama orang dan nomor telepon

`[keputusan work owner]` **P11, P12, P28, P25**

| Yang diganti | Menjadi |
| --- | --- |
| ⛔ kode peran di `OperatorID.pyTelephone` | ⭐ **medan/tabel peran tersendiri**; kolom telepon kembali berisi nomor telepon |
| ⛔ pemeriksaan "apakah pengguna ini orang tertentu" *(12 berkas)* | ⭐ **pemeriksaan peran** |
| ⛔ pencarian antrean menurut **nomor urut** | ⭐ **pemeriksaan keanggotaan antrean** |

⭐ **Dua peran terbaca** `[terverifikasi]`: `TREATY1` *(staf)* dan `SPVTREATY1` *(supervisor)*.

⛔⛔ **Dan di sinilah keputusan ini belum dapat dijalankan sepenuhnya** `[terbuka]`:
⭐ tangga punya **tiga** posisi, sedangkan peran yang terbaca baru **dua**. ⛔ Pemetaan
**nama → peran** untuk ke-12 tempat **tidak ada di korpus**, dan pada **dua layar** nama yang sama
dipakai **dua arah** — satu bagian muncul hanya untuk orang itu, bagian lain untuk semua kecuali
orang itu. ⛔ **Menebak arahnya berarti menampilkan bagian layar kepada orang yang salah.**

⭐ **Pencarian antrean menurut nomor urut dapat langsung dikerjakan** — tidak menunggu siapa pun.
`[terverifikasi]` Tiga tempat di modul ini; **41 berkas pada 9 modul** di seluruh korpus.

### 5.5 ⭐ Penggolongan jenis usaha

`[keputusan work owner]` `[terverifikasi]` **P19, P23**

`DecisionTable\BusinessType_DeT.xml` — **36 baris**, dua kolom penguji
*(`.Quotation.GroupPanel` dan `.Quotation.BusinessOldId`)*, menghasilkan `BusinessType`.

| Sifat | Nilai |
| --- | --- |
| baris | **36** |
| nilai `BusinessOldId` | ⭐ **128**, seluruhnya **unik**, nol muncul dua kali |
| ⭐ bawaan | ⭐ **`"UNKNOWN"`** |
| ⭐ cara evaluasi | ⭐ **berhenti di baris pertama yang cocok** *(`pyEvaluateAllRows` = `no`)* |
| baris ber-daftar-OR | 8, memuat 106 nilai |
| baris menguji panel saja | 6 — penampung per panel |

⛔ **Keterangan rule `When` diabaikan seluruhnya** `[keputusan work owner]` **P23** — spec ditulis
dari **syarat yang dijalankan**. `[terverifikasi]` Dari 75 rule `When`: 41 keterangannya cocok,
**9 bertentangan**, 21 tidak pernah diisi, 4 syaratnya tidak terekspor.

### 5.6 ⭐ Uang, persentase, dan pajak

`[keputusan work owner]` **P29, P46**

| Ketetapan | Isinya |
| --- | --- |
| ⭐ presisi | **penuh**; ⛔ pembulatan **hanya di titik penyajian**, tidak pernah di repository |
| ⛔ tipe | ⛔ **tidak pernah `float`** *(CLAUDE.md §7)* |
| ⭐ persentase | `DEDUCTION1` `DEDUCTION2` `BROKERAGE` `RNM_SHARE` adalah **persentase** — `12.5` = **12,5 persen**, bukan uang |
| ⭐ pajak brokerage | ditiru apa adanya, termasuk ketidakseragaman presisi |

⭐ **Rumus pajak brokerage** `[terverifikasi]`, ada di **12 tempat pada 7 Activity**:

```
BrokerageFeeSebenarnya = @if(TypeTax == "Inclusive", Deduction / (102.2/100), Deduction)
```

⚠️ **Tiga sifat yang ikut disalin, seluruhnya disadari** `[terverifikasi]`:
⛔ nilai `"Exclusive"` **tidak pernah tertulis di mana pun** — cabang kedua hanya tercapai bila
medan berisi sesuatu yang lain, **termasuk kosong**; ⛔ pembandingnya **teks persis**, sehingga
huruf kecil atau spasi di belakang jatuh ke cabang kedua dan brokerage **tidak** dibagi 1,022 —
selisih **2,2 %** tanpa pesan galat; ⚠️ presisi tidak seragam untuk rumus yang sama.

⭐ **Angka uang ditampilkan apa adanya** `[keputusan work owner]` **P43** — ⛔ **tidak dihitung
ulang** dari data dasar. `[terverifikasi]` Kedelapan medan uang adalah **kolom tersimpan** pada
view; nol `Activity` atau `DataTransform` menghasilkannya.
⭐ **Akibatnya menguntungkan:** penyajiannya **tidak bergantung pada P18**.

### 5.7 ⭐ Keutuhan penyimpanan — satu transaksi

`[keputusan work owner]` `[penyimpangan sadar]` **P2**

⭐ **Seluruh urutan penyimpanan dibungkus satu transaksi.** Kegagalan di langkah mana pun
membatalkan seluruhnya.

⚠️ **Ini LEBIH KETAT daripada sistem lama, dan disengaja.** `[data DBA]`

| Procedure | `COMMIT` di dalam |
| --- | --- |
| `POOLDATA.PEGA_TREATY_IN` | **ada** |
| `POOLDATA.PEGA_JSON_POLIS_TREATYIN` | **ada** |
| `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | **tidak ada** |

⛔ Untuk dua yang pertama terjadi **commit ganda** — procedure meng-commit, lalu blok pemanggil
meng-commit lagi. ⛔ Sistem lama **tidak punya jaminan keutuhan** untuk urutan yang memanggil
keduanya.

⭐ **Setiap query menulis skema `POOLDATA.` eksplisit** `[keputusan work owner]` **P3** —
⛔ tidak mengandalkan skema bawaan pengguna. `[terverifikasi]` Empat nama tabel muncul dua ejaan
di 1.151 naskah SQL; keempatnya **satu tabel**, milik `POOLDATA`.

### 5.8 ⭐ Tanggal

`[keputusan work owner]` `[penyimpangan sadar]` **P32, P35**

⭐ **Tanggal disimpan sebagai tipe tanggal Oracle**, ⛔ bukan teks, dengan **satu format** di
seluruh sistem. Format tampilan diurus di lapisan layar.

⚠️ **Perilaku Pega sekarang** `[terverifikasi]`: dua format di berkas yang sama —
`dd/MM/yyyy` untuk tanggal mulai dan `MM/dd/yyyy hh:mm a` untuk tanggal akhir, keduanya di
`DataTransform\InputPolicyTreatyIn_preDT.xml`. `[data DBA]` Penyimpanannya bertipe teks:
`POOLDATA.TANGGAL_CLOSING.TANGGAL` adalah `VARCHAR2(10)`.

⭐ **Tanggal akhir kosong tetap diisi HARI INI** `[keputusan work owner]` **P35** — perilaku lama
dipertahankan, walau hasilnya kontrak bermasa berlaku nol hari.
⚠️ **Konsekuensi dicatat, bukan dibantah:** migrasi data lama menghasilkan kontrak
ber-`EndDate` = `StartDate`, dan laporan "masih berlaku" mengeluarkannya sejak hari pertama.
⛔ **Disengaja, bukan cacat migrasi.**

### 5.9 ⭐ Kegagalan pembacaan menghentikan proses

`[keputusan work owner]` `[penyimpangan sadar]` **P31**

⭐ Bila pembacaan data kontrak gagal, sistem baru **berhenti** dan **menampilkan galat**.
⛔ Proses tidak diteruskan dengan data separuh terisi.

⚠️ **Perilaku Pega sekarang** `[terverifikasi]`: **enam langkah Java identik** membaca data lalu
`adoptJSONObject`; bila gagal, galatnya **hanya ditulis ke log** dan aktivitas **LANJUT**.
⛔ Tidak ada penghentian, tidak ada pemberitahuan.
⭐ **Alasan penyimpangan:** di sistem yang mencatat komitmen keuangan, **kegagalan yang terlihat
lebih murah daripada kegagalan yang tersembunyi**.

### 5.10 ⭐ Jejak audit

`[keputusan work owner]` **P4, P33, P40**

| Kolom | Diisi dari |
| --- | --- |
| ⭐ `OPERATORID` | ⭐ **identitas akses login** |
| ⭐ `PIC` | ⭐ **nama tampilan** |

⚠️ **Perubahan perilaku yang disadari** `[penyimpangan sadar]` `[terverifikasi]`: sistem lama
**tidak pernah** menulis ke `OPERATORID` dari aplikasi — **nol dari 1.151 naskah SQL** di 21 modul
menyebutnya sebagai kolom. ⭐ Sistem baru **mulai mengisi kolom yang selama ini kosong** — itu
**perbaikan jejak audit, bukan peniruan**.

⭐ **`OperatorName` berarti satu hal saja: nama tampilan** `[keputusan work owner]` **P33**.
⛔ Pengisian dari pengenal akun di `DataTransform\DeptHeadTreatyInUW_preDT.xml` adalah **bug**,
⛔ bukan perbedaan maksud antar tahap, dan **tidak ditiru**.

⭐ **Nama orang tidak ditulis di dalam teks pemberitahuan** `[keputusan work owner]` **P40** —
diambil dari data, seperti pesan lain di tempat yang sama yang sudah benar.

### 5.11 Layar

`[keputusan work owner]` **P45, P47, P44, P39, P37**

| Ketetapan | Angka |
| --- | ---: |
| ⭐ medan **wajib** | **27** medan berbeda di **6** layar |
| ⭐ medan **terkunci permanen** | **38** — ⭐ **36** di layar Dept Head, **2** di layar biasa |
| ⭐ bagian layar **dimatikan** | **91** — ⛔ **80 dibuang tanpa ditanyakan**, 11 menyembunyikan 128 medan |

⭐ **Medan wajib BERBEDA menurut tingkat** `[terverifikasi]`. Layar Dept Head **tidak** mewajibkan
enam medan yang wajib di layar admin — `ClaimPaymentType` `ClaimType` `IDCurrency` `Quartal`
`TypeTax` `YearOfQuartal` — ⭐ tetapi **mewajibkan `ResultOnp1`** yang tidak wajib di sana.

⭐ **Setiap kunci mengenai tepat SATU medan**, ⛔ bukan bagian atau tab. `[terverifikasi]`
`pyReadOnlyCondition` selalu-benar: `1==1` 32× · `1=1` 2× · `1 = 1` 2× · `ALWAYS` 2×.

⭐ **Kepala Departemen BUKAN hanya melihat** `[terverifikasi]` — ia mengisi **tujuh medan** sendiri,
lalu memutuskan.

⭐ **Delapan puluh elemen mati dibuang tanpa ditanyakan** `[keputusan work owner]` **P44** —
`[terverifikasi]` sebabnya struktural: `pyCondition` menempel pada **sel tunggal** dan
**tidak menyembunyikan satu medan pun**, sedangkan `pyContainerVisibleWhen` menempel pada **wadah**
dan menyembunyikan **128 medan**. ⭐ Membangun atau membuang yang 80 sama-sama berbiaya nol.

⭐ **Pengecualian mata uang dipertahankan** `[keputusan work owner]` **P39** — `[terverifikasi]`
kode `ITL` disaring di dua `ReportDefinition`. ⚠️ `[dugaan]` itu **pembersihan mata uang mati**
*(Lira Italia, digantikan Euro 2002)*, bukan kebijakan treaty.

⭐ **Kode jenis kontrak disalin apa adanya** `[keputusan work owner]` **P37** — ⛔ tidak
diterjemahkan, tidak dinormalkan, tidak diberi enumerasi bernama. ⛔ Tidak ada validasi daftar
nilai; nilai di luar yang dikenal tetap diterima.

### 5.12 Treaty keluar dibaca, tidak ditulis

`[keputusan work owner]` **P9**

⭐ **Ini jalur retrosesi, bukan pelanggaran batas modul.** Saat Nusantara Re menerima treaty masuk,
sebagian ditempatkan kembali keluar, dan layar inward perlu menampilkannya.

⭐ Data treaty outward **tetap dibaca** dari konteks ini, ⛔ dan **tetap tidak boleh ditulis**.
`[terverifikasi]` 16 berkas menyebut TreatyOut; ⛔ **nol** operasi tulis.

### 5.13 Arsitektur

⭐ Arah ketergantungan **`handlers → services → repository`** *(CLAUDE.md §5)*.
⛔ Uang tidak pernah `float` *(§7)*.

---

## 6 · Testing Decisions

### 6.1 Apa yang membuat sebuah test baik di sini

⭐ **Hanya perilaku luar yang diuji**, ⛔ bukan rincian dalam. Sebuah test yang pecah ketika nama
fungsi berubah — padahal perilakunya tetap — adalah test yang buruk.

⭐ **Setiap AC di bab 7 ditulis agar dapat gagal.** ⛔ Butir yang tidak mungkin gagal bukan
kriteria, melainkan pernyataan.

### 6.2 Jahitan uji — sesedikit mungkin

⭐ **Tiga jahitan, dan ketiganya sudah ada**, ⛔ tidak ada jahitan baru yang diusulkan:

| # | Jahitan | Yang diuji di situ |
| ---: | --- | --- |
| ⭐ **1** | **HTTP handler** | daur hidup berkas, tangga persetujuan, wewenang, validasi medan wajib |
| ⭐ **2** | **repository lawan Oracle sungguhan** | pembacaan view, penulisan riwayat, **keutuhan satu transaksi**, skema eksplisit |
| ⭐ **3** | **fungsi murni penggolong dan penghitung** | 36 baris penggolongan, rumus pajak brokerage, pembacaan persentase |

⭐ **Jahitan 1 dipakai sebanyak mungkin.** Jahitan 3 hanya untuk yang benar-benar fungsi murni.

### 6.3 Yang wajib diuji lawan basis data sungguhan

⛔ **Keutuhan transaksi tidak dapat diuji dengan tiruan.** Uji yang menyuntikkan kegagalan di
tengah urutan penyimpanan, lalu memastikan **nol baris tersisa**, wajib berjalan lawan Oracle.

⛔ Begitu pula **skema eksplisit**: uji yang menyambung sebagai pengguna **selain** `POOLDATA` dan
memastikan query tetap menemukan tabelnya.

### 6.4 Prior art

⭐ Sebelas `spec.md` modul lain di `.scratch\*\` memakai bentuk AC yang sama — pernyataan +
"test yang menemukan keadaan sebaliknya **gagal**" + rujukan bab.

---

## 7 · Acceptance Criteria

⛔ **Bab ini tidak memutuskan apa pun.** Ia menyatakan ulang bab 5 dalam bentuk yang dapat diuji
dari luar. ⛔ Butir yang terasa seperti keputusan baru adalah salah tulis.

### Penanda persetujuan dan putusan alur

1. `[keputusan work owner]` `IsApproved` bernilai `0` diperlakukan **ditolak**. Test yang menemukan
   `0` diperlakukan disetujui **gagal**. *(Bab 5.3)*
2. `[keputusan work owner]` `IsApproved` bernilai apa pun selain `0`, **termasuk kosong**,
   diperlakukan **disetujui**. Test yang menemukan nilai kosong diperlakukan ditolak **gagal**.
   *(Bab 5.3)*
3. `[terverifikasi]` Perbandingan `IsApproved` dilakukan sebagai **teks**. Test yang menemukan
   perbandingan numerik **gagal**. *(Bab 5.3)*
4. `[terverifikasi]` Aturan yang dipakai berasal dari tabel keputusan, bukan dari rule `When`
   bernama sama. Test yang menemukan perilaku `IsApproved = 1` **gagal**. *(Bab 5.3)*
5. `[keputusan work owner]` Admin menolak ⇒ berkas **diselesaikan sebagai ditolak**. Test yang
   menemukan berkas masih beredar **gagal**. *(Bab 5.3)*
6. `[keputusan work owner]` Atasan menolak ⇒ berkas **kembali ke admin**. Test yang menemukan
   berkas diselesaikan **gagal**. *(Bab 5.3)*
7. `[keputusan work owner]` Admin menyetujui ⇒ berkas naik ke Sec Head. Test yang menemukan berkas
   melompati Sec Head **gagal**. *(Bab 5.2, 5.3)*
8. `[keputusan work owner]` Sec Head menyetujui ⇒ berkas naik ke Dept Head. Test yang menemukan
   tujuan lain **gagal**. *(Bab 5.2, 5.3)*
9. `[keputusan work owner]` Dept Head menyetujui ⇒ realisasi **selesai**. Test yang menemukan
   jenjang keempat **gagal**. *(Bab 5.2)*
10. `[keputusan work owner]` Posisi `GroupLeader` dan `Director` **tidak ada** di sistem baru. Test
    yang menemukan berkas dirutekan ke salah satunya **gagal**. *(Bab 5.2)*

### Antrean dan wewenang

11. `[keputusan work owner]` Setiap penugasan masuk ke **antrean bersama**. Test yang menemukan
    kotak masuk pribadi **gagal**. *(Bab 5.2)*
12. `[keputusan work owner]` Wewenang ditentukan **peran**, bukan nama orang. Test yang menemukan
    pemeriksaan identitas orang **gagal**. *(Bab 5.4)*
13. `[keputusan work owner]` Peran disimpan di medan peran. Test yang menemukan kode peran di
    kolom nomor telepon **gagal**. *(Bab 5.4)*
14. `[keputusan work owner]` Keanggotaan antrean diperiksa **menurut nama antrean**. Test yang
    menemukan pemeriksaan menurut **nomor urut** daftar **gagal**. *(Bab 5.4)*

### Sumber data

15. `[keputusan work owner]` Data kontrak dibaca dari view relasional. Test yang menemukan
    pembacaan kolom dokumen JSON **gagal**. *(Bab 5.1)*
16. `[keputusan work owner]` Sistem baru **tidak menulis** JSON. Test yang menemukan penulisan
    kolom dokumen **gagal**. *(Bab 5.1)*
17. `[terverifikasi]` Ke-33 medan yang dipakai laporan tersedia dari view. Test yang menemukan
    medan laporan tidak terisi **gagal**. *(Bab 5.1)*
18. `[keputusan work owner]` Kedelapan medan uang **ditampilkan apa adanya** dari penyimpanan.
    Test yang menemukan perhitungan ulang **gagal**. *(Bab 5.6)*

### Penggolongan

19. `[terverifikasi]` Penggolongan berhenti di **baris pertama yang cocok**. Test yang menemukan
    seluruh baris dievaluasi **gagal**. *(Bab 5.5)*
20. `[terverifikasi]` Kode bisnis yang tidak cocok baris mana pun menghasilkan **`"UNKNOWN"`**.
    Test yang menemukan nilai kosong atau galat **gagal**. *(Bab 5.5)*
21. `[terverifikasi]` Ke-128 nilai kode bisnis menghasilkan penggolongan yang sama seperti sistem
    lama. Test yang menemukan satu saja berbeda **gagal**. *(Bab 5.5)*
22. `[keputusan work owner]` Syarat rule `When` diambil dari **yang dijalankan**, bukan dari
    keterangannya. Test yang menemukan perilaku mengikuti keterangan pada salah satu dari **9**
    rule yang bertentangan **gagal**. *(Bab 5.5)*

### Uang, persentase, pajak

23. `[keputusan work owner]` Nilai uang disimpan **berpresisi penuh**. Test yang menemukan
    pembulatan di lapisan repository **gagal**. *(Bab 5.6)*
24. `[keputusan work owner]` Pembulatan hanya terjadi di titik penyajian. Test yang menemukan nilai
    tersimpan sudah dibulatkan **gagal**. *(Bab 5.6)*
25. `[keputusan work owner]` Uang **tidak pernah** diwakili `float`. Test yang menemukan tipe
    pecahan biner **gagal**. *(Bab 5.6, 5.13)*
26. `[keputusan work owner]` `DEDUCTION1` `DEDUCTION2` `BROKERAGE` `RNM_SHARE` dibaca sebagai
    **persentase**. Test yang menemukan `12.5` diperlakukan sebagai jumlah uang **gagal**.
    *(Bab 5.6)*
27. `[terverifikasi]` `TypeTax` bernilai `"Inclusive"` ⇒ potongan dibagi **1,022**. Test yang
    menemukan potongan dipakai apa adanya **gagal**. *(Bab 5.6)*
28. `[terverifikasi]` `TypeTax` bernilai apa pun selain `"Inclusive"` persis — **termasuk kosong,
    huruf kecil, atau berspasi** — ⇒ potongan dipakai apa adanya. Test yang menemukan pembagian
    1,022 **gagal**. *(Bab 5.6)*

### Keutuhan penyimpanan

29. `[keputusan work owner]` `[penyimpangan sadar]` Seluruh urutan penyimpanan berada dalam **satu
    transaksi**. Test yang menyuntikkan kegagalan di tengah dan menemukan **ada baris tersisa**
    **gagal**. *(Bab 5.7)*
30. `[keputusan work owner]` Setiap query menyebut skema **`POOLDATA.` eksplisit**. Test yang
    menyambung sebagai pengguna lain dan menemukan query gagal menemukan tabel **gagal**.
    *(Bab 5.7)*
31. `[keputusan work owner]` Nomor polis terbentuk dari deret, tanpa bentrok. Test yang menemukan
    dua berkas bernomor sama **gagal**. *(Bab 5.7)*

### Tanggal

32. `[keputusan work owner]` `[penyimpangan sadar]` Tanggal disimpan sebagai **tipe tanggal**,
    bukan teks. Test yang menemukan kolom teks **gagal**. *(Bab 5.8)*
33. `[keputusan work owner]` `[penyimpangan sadar]` Satu format tanggal dipakai di seluruh sistem.
    Test yang menemukan dua format berdampingan **gagal**. *(Bab 5.8)*
34. `[keputusan work owner]` Tanggal akhir kosong diisi **tanggal hari ini**. Test yang menemukan
    penambahan satu tahun **gagal**. *(Bab 5.8)*
35. `[keputusan work owner]` Tanggal mulai dan tanggal laporan kosong diisi tanggal hari ini. Test
    yang menemukan medan tetap kosong **gagal**. *(Bab 5.8)*

### Kegagalan pembacaan

36. `[keputusan work owner]` `[penyimpangan sadar]` Kegagalan pembacaan data kontrak
    **menghentikan** proses. Test yang menemukan proses berlanjut **gagal**. *(Bab 5.9)*
37. `[keputusan work owner]` `[penyimpangan sadar]` Kegagalan pembacaan **menampilkan galat kepada
    pengguna**. Test yang menemukan galat hanya tercatat di log **gagal**. *(Bab 5.9)*
38. `[keputusan work owner]` Tidak ada kasus yang tersimpan dari pembacaan yang gagal. Test yang
    menemukan berkas tersimpan separuh **gagal**. *(Bab 5.9)*

### Jejak audit

39. `[keputusan work owner]` `OPERATORID` diisi dari **identitas akses login**. Test yang menemukan
    nama tampilan **gagal**. *(Bab 5.10)*
40. `[keputusan work owner]` `PIC` diisi dari **nama tampilan**. Test yang menemukan pengenal akun
    **gagal**. *(Bab 5.10)*
41. `[keputusan work owner]` `[penyimpangan sadar]` `OPERATORID` **terisi** pada setiap penulisan
    riwayat. Test yang menemukan kolom kosong **gagal**. *(Bab 5.10)*
42. `[keputusan work owner]` `OperatorName` diisi dari **nama tampilan di setiap tahap**, termasuk
    tahap Dept Head. Test yang menemukan pengenal akun di tahap mana pun **gagal**. *(Bab 5.10)*
43. `[keputusan work owner]` Setiap perpindahan tahap menulis satu baris riwayat. Test yang
    menemukan perpindahan tanpa jejak **gagal**. *(Bab 5.10)*
44. `[keputusan work owner]` Pemberitahuan menyebut nama orang **dari data**. Test yang menemukan
    nama tertanam di dalam teks **gagal**. *(Bab 5.10)*

### Layar — medan wajib

45. `[terverifikasi]` Layar admin mewajibkan medan wajibnya; **27** medan berbeda tersebar di
    **6** layar. Test yang menemukan medan wajib dapat dilewati **gagal**. *(Bab 5.11)*
46. `[terverifikasi]` Layar Dept Head **tidak** mewajibkan `ClaimPaymentType` `ClaimType`
    `IDCurrency` `Quartal` `TypeTax` `YearOfQuartal`. Test yang menemukan salah satunya diwajibkan
    di sana **gagal**. *(Bab 5.11)*
47. `[terverifikasi]` Layar Dept Head **mewajibkan** `ResultOnp1`. Test yang menemukan medan itu
    opsional di sana **gagal**. *(Bab 5.11)*
48. `[keputusan work owner]` Berkas tidak dapat disimpan bila medan wajib pada tingkat itu kosong.
    Test yang menemukan penyimpanan berhasil **gagal**. *(Bab 5.11)*

### Layar — medan terkunci

49. `[keputusan work owner]` **38** medan terkunci permanen. Test yang menemukan salah satunya
    dapat disunting **gagal**. *(Bab 5.11)*
50. `[terverifikasi]` **36** dari 38 berada di layar Dept Head; **2** di layar biasa. Test yang
    menemukan seluruh 38 di layar Dept Head **gagal**. *(Bab 5.11)*
51. `[terverifikasi]` Setiap kunci mengenai **tepat satu medan**. Test yang menemukan satu kunci
    mengunci sebuah bagian atau tab **gagal**. *(Bab 5.11)*
52. `[terverifikasi]` Kepala Departemen **dapat mengisi tujuh medan**. Test yang menemukan layarnya
    sepenuhnya hanya-baca **gagal**. *(Bab 5.11)*

### Layar — tampilan

53. `[keputusan work owner]` **80** elemen mati **tidak dibangun**. Test yang menemukan salah
    satunya muncul **gagal**. *(Bab 5.11)*
54. `[keputusan work owner]` Daftar pilihan mata uang **tidak memuat `ITL`**. Test yang menemukan
    kode itu di daftar **gagal**. *(Bab 5.11)*
55. `[keputusan work owner]` Kode jenis kontrak ditampilkan **apa adanya** bila keterangannya belum
    tersedia. Test yang menemukan kode diterjemahkan sendiri **gagal**. *(Bab 5.11)*
56. `[keputusan work owner]` Nilai kode di luar daftar yang dikenal **tetap diterima dan disimpan**.
    Test yang menemukan penolakan **gagal**. *(Bab 5.11)*

### Treaty keluar

57. `[keputusan work owner]` Data treaty keluar **dapat dibaca** dari konteks realisasi treaty
    masuk. Test yang menemukan data itu tidak tampil **gagal**. *(Bab 5.12)*
58. `[keputusan work owner]` Data treaty keluar **tidak pernah ditulis** dari konteks ini. Test
    yang menemukan operasi tulis **gagal**. *(Bab 5.12)*

### Penawaran ganda

59. `[terverifikasi]` Peringatan dipasang ketika jumlah berkas klaim terhubung **lebih dari nol**.
    Test yang menemukan peringatan tidak muncul **gagal**. *(Bab 5.11)*

### Arsitektur

60. `[keputusan work owner]` Arah ketergantungan **`handlers → services → repository`**. Test
    arsitektur yang menemukan repository memanggil service **gagal**. *(Bab 5.13)*

### Lingkup — yang TIDAK dibangun

61. `[keputusan work owner]` **28** aturan yatim **tidak dimigrasi**. Test yang menemukan salah
    satunya terpanggil **gagal**. *(Bab 8)*
62. `[keputusan work owner]` Enam aturan pembongkar JSON **tidak dimigrasi**. Test yang menemukan
    salah satunya terpanggil **gagal**. *(Bab 8)*
63. `[keputusan work owner]` Rule `When\isApproved` **tidak dimigrasi**. Test yang menemukan
    perilakunya **gagal**. *(Bab 8)*
64. `[keputusan work owner]` Properti `isApprovedtoDeptHead` **tidak dibangun**, dan kolomnya tidak
    dibuat. Test yang menemukan kolom itu **gagal**. *(Bab 8)*
65. `[keputusan work owner]` Medan "nomor surat" **tidak** dipakai menyimpan penanda arah. Test yang
    menemukan nilai routing di sana **gagal**. *(Bab 8)*

### Penggolongan lanjutan

66. `[terverifikasi]` `isFOR` bernilai `EDM` ⇒ jalur **endorsemen**; `POLICY` ⇒ **polis baru**. Test
    yang menemukan keduanya menempuh cara penyimpanan yang sama **gagal**. *(Bab 5.5)*
67. `[terverifikasi]` `IsLife` dimigrasi apa adanya, beserta ke-16 nilai kode lini jiwa. Test yang
    menemukan satu nilai berbeda **gagal**. *(Bab 5.5)*

### Migrasi

68. `[keputusan work owner]` Data lama terbaca di sistem baru. Test yang menemukan berkas lama
    tidak dapat dibuka **gagal**. *(Bab 5.1)*
69. `[keputusan work owner]` Migrasi menghasilkan kontrak ber-`EndDate` = `StartDate` untuk berkas
    yang tanggal akhirnya kosong. Test yang menemukan tanggal lain **gagal**. *(Bab 5.8)*
70. `[keputusan work owner]` Setiap penyimpangan dari perilaku lama tercatat di bab 5 beserta
    alasannya. Test yang menemukan penyimpangan tanpa catatan **gagal**. *(Bab 5, 10)*

### Riwayat dan catatan

71. `[keputusan work owner]` Catatan pengguna tersimpan bersama tanggal dan operatornya. Test yang
    menemukan catatan tanpa salah satunya **gagal**. *(Bab 5.10)*
72. `[terverifikasi]` Riwayat dapat dibaca berurutan waktu. Test yang menemukan urutan acak
    **gagal**. *(Bab 5.10)*

### Nomor polis

73. `[terverifikasi]` Nomor polis memuat awalan tetap, penanda treaty, bulan-tahun, dan nomor urut
    berdigit tetap. Test yang menemukan pola lain **gagal**. *(Bab 5.7)*
74. `[keputusan work owner]` Nomor polis dibentuk **sekali** per berkas. Test yang menemukan nomor
    berubah pada penyimpanan berikutnya **gagal**. *(Bab 5.7)*

### Validasi lintas medan

75. `[terverifikasi]` Berkas non-proporsional ditandai sesuai jenis proporsinya. Test yang menemukan
    penandaan terbalik **gagal**. *(Bab 5.5)*
76. `[terverifikasi]` Penanda "punya penempatan keluar" dipasang bila kode jenis kontrak cocok salah
    satu dari dua kode yang dikenal. Test yang menemukan penandaan pada kode lain **gagal**.
    *(Bab 5.11)*

### Pesan galat

77. `[keputusan work owner]` Pembersihan pesan galat di awal aktivitas diterima apa adanya —
    **12** tempat. Test yang menemukan pesan lama bertahan **gagal**. *(Bab 9)*
78. `[terbuka]` **5** tempat yang menghapus pesan **sesudah** validasi memasangnya **ditahan** dan
    tidak dibangun. Test yang menemukan salah satunya sudah dibangun **gagal**. *(Bab 9)*

### ⚠️ DICABUT LALU DIGANTI — semula "Yang menunggu P18"

> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22
> **Kedua AC di bawah DICABUT, lalu DIGANTI di nomor yang sama.**
>
> ⭐ **Kenapa diganti, bukan dihapus:** menghapusnya akan mematahkan rujukan AC di **16 tiket** dan
> di `00-PETA-AC.md`, serta mengubah jumlah AC dari **96**. Bunyi lamanya dikutip utuh di bawah dan
> **tidak dihapus**; yang berubah adalah isinya, bukan nomornya. ⚠️ **Jumlah AC tetap 96.**
>
> Keduanya lahir dari premis bahwa isi 268 langkah tidak terekspor. **Premis itu keliru** — isinya
> ada, di tag `PropertiesName`/`PropertiesValue` **tanpa awalan `py`**. Lihat `VERIFIKASI-P18.md`.

79. ⛔ **DICABUT.** Bunyi semula, dikutip utuh:
    > *"`[terbuka]` Rantai perhitungan uang **tidak dibangun** sebelum P18 dijawab. Test yang
    > menemukan rumus yang ditebak **gagal**. *(Bab 8.2)*"*

    ⭐ **Ganti:** rantai perhitungan uang **dibangun dari isi langkah yang terbaca di ekspor**.
    Test yang menemukan rumus **yang tidak berasal dari `PropertiesValue` suatu langkah** gagal.
    *(Bab 8.2)*
80. ⛔ **DICABUT.** Bunyi semula, dikutip utuh:
    > *"`[terbuka]` Layar Dept Head **tidak dapat disimpan** selama sembilan medan
    > wajib-sekaligus-terkunci belum terselesaikan. Test yang menemukan penyimpanan berhasil
    > **gagal**. *(Bab 8.2)*"*

    ⭐ **Ganti:** layar Dept Head **dapat disimpan**. `[terverifikasi]` Dari **34** medan
    wajib-sekaligus-terkunci, **28** punya langkah pengisi; tiga sisanya — `.ExcessLoss`,
    `.OutstandingClaim`, `.SalvageValue` — **diketik underwriter** di `DetailPolicyTreatyIn` /
    `GeneralPolicyTreatyIn` *(`bacasaja=false`)* dan hanya **ditampilkan terkunci** di layar Dept
    Head. Test yang menemukan layar Dept Head **tidak dapat disimpan** gagal. *(Bab 8.2)*

### Yang menunggu pemetaan peran

81. `[terbuka]` Ke-12 tempat yang kini memeriksa identitas orang **ditandai tertunda**, tidak
    dibangun dengan peran yang ditebak. Test yang menemukan pemetaan tanpa dasar **gagal**.
    *(Bab 5.4, 9)*
82. `[terbuka]` Pada dua layar, arah pemeriksaan — "muncul untuk peran ini" lawan "muncul untuk
    semua kecuali peran ini" — **tidak ditebak**. Test yang menemukan arah yang dibalik **gagal**.
    *(Bab 5.4, 9)*

### Ketahanan

83. `[keputusan work owner]` Kegagalan menyimpan riwayat membatalkan seluruh transaksi. Test yang
    menemukan berkas tersimpan tanpa riwayat **gagal**. *(Bab 5.7, 5.10)*
84. `[terverifikasi]` Nilai kosong pada penanda persetujuan **tidak** menghentikan alur. Test yang
    menemukan galat **gagal**. *(Bab 5.3)*

### Penyajian

85. `[keputusan work owner]` Angka uang ditampilkan beserta **kode mata uang pasangannya**. Test
    yang menemukan angka tanpa mata uang **gagal**. *(Bab 5.6)*
86. `[terbuka]` Format penyajian uang — jumlah desimal dan pemisah ribuan — **belum ditetapkan**;
    penyajian tidak dibangun sebelum ditetapkan. Test yang menemukan format yang ditebak **gagal**.
    *(Bab 9)*

### Antarmuka luar

87. `[terbuka]` Pengiriman ke layanan luar di akhir alur **tidak dibangun** sebelum muatannya
    diketahui. Test yang menemukan muatan yang ditebak **gagal**. *(Bab 8.2, 9)*

### Lingkup data

88. `[keputusan work owner]` Aturan yang tidak terjangkau dari titik masuk **tidak dibangun**. Test
    yang menemukan salah satu dari **561** langkah yatim terpanggil **gagal**. *(Bab 8)*
89. `[terverifikasi]` Modul ini membaca **39** kolom view; nol medan yang dipakai tidak tersedia.
    Test yang menemukan medan hilang **gagal**. *(Bab 5.1)*

### Konsistensi skema

90. `[keputusan work owner]` Keempat nama tabel berejaan ganda diperlakukan sebagai **satu tabel**
    milik `POOLDATA`. Test yang menemukan dua tabel berbeda **gagal**. *(Bab 5.7)*

### Peran dan tangga

91. `[terbuka]` Peran yang tersedia **belum cukup** membedakan tiga posisi tangga; pemetaan tidak
    diselesaikan dengan tebakan. Test yang menemukan peran karangan **gagal**. *(Bab 5.4, 9)*
92. `[keputusan work owner]` Berkas menunggu **posisi**, bukan orang. Test yang menemukan berkas
    terikat pada satu pengguna **gagal**. *(Bab 5.2)*

### Jejak keputusan

93. `[keputusan work owner]` Setiap AC merujuk bab asalnya. Test dokumen yang menemukan butir tanpa
    rujukan **gagal**. *(Bab 1.2, 7)*
94. `[keputusan work owner]` Setiap AC membawa penanda. Test dokumen yang menemukan butir tanpa
    penanda **gagal**. *(Bab 1.2, 7)*
95. `[terbuka]` Butir terbuka **tidak ditutup** oleh spec ini. Test dokumen yang menemukan butir
    terbuka dinyatakan selesai **gagal**. *(Bab 9)*
96. `[keputusan work owner]` Nilai berupa nama orang **tidak** muncul di artefak mana pun. Test yang
    menemukan nama orang **gagal**. *(Bab 1.2, 10)*

---

## 8 · Out of Scope

### 8.1 Yang dikeluarkan, dan sebabnya

| Yang dikeluarkan | Sebab | Butir asal |
| --- | --- | --- |
| ⛔ **28 aturan yatim, 561 langkah `Property-Set`** | ⛔ **tidak terjangkau** dari titik masuk; keluarga spreading/premi milik **Fac In** | keadaan Bab 2 |
| ⛔ **6 aturan pembongkar JSON** | ⭐ penyimpanan pindah ke **view relasional** | **P29, P15** |
| ⚠️ ~~**Rantai perhitungan uang**~~ — **TIDAK LAGI DIKELUARKAN** | ⛔ ~~isi 268 langkah belum terkirim~~ — **RALAT 2026-09-22: isinya ada di ekspor** | ⭐ **P18 ditarik** |
| ⛔ **20 nomor polis di dalam aturan uang** | ⭐ rantainya milik **Fac In**, tampak di sini karena layar ditanam bersama | **P20** |
| ⛔ **Dua aturan uang yang hanya berupa catatan** | ⭐ rantai yang sama, **tidak terjangkau** dari sini | **P21** |
| ⛔ **Pencarian antrean menurut nomor urut** | ⭐ diganti **pemeriksaan keanggotaan** | **P25** |
| ⛔ **Rule `When\isApproved`** | ⛔ **bukan aturan yang hidup** — nol kotak Decision menyambung ke sana | **P6** |
| ⛔ **Properti `isApprovedtoDeptHead`** | ⛔ **sudah tidak dipakai**; berkas rule-nya nol di seluruh korpus | **P36** |
| ⛔ **Posisi `GroupLeader` dan `Director`** | ⭐ tangga menyusut menjadi **tiga** | **P13** |
| ⛔ **Jalur klaim → Manajer Klaim → Direktur** | ⭐ sudah lama tidak dipakai | **P5** |
| ⛔ **Nama orang sebagai penjaga wewenang** | ⭐ diganti **peran** | **P12, P28** |
| ⛔ **`LetterNo` sebagai penanda arah** | ⭐ routing memakai medan tersendiri | **P7** |
| ⛔ **Pemanggilan `CopyToPolicy`** | ⛔ **tidak digunakan**; tidak perlu diminta | **P30** |
| ⛔ **80 elemen layar mati** | ⭐ **tidak menyembunyikan satu medan pun**; berbiaya nol di kedua arah | **P44** |

⭐ **Folder modul ini adalah *dependency closure*** — ia memuat aturan milik modul lain karena
pewarisan kelas Pega. ⛔ **Bukan daftar pekerjaan Treaty.**

### 8.2 ⭐⭐ Karantina rantai uang — **DIBUKA 2026-09-22**

> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — judul dan kalimat pembuka bab ini semula
> berbunyi:
> > *"### 8.2 ⛔⛔ Karantina rantai uang — ⭐ **Ditulis tersendiri karena ia satu-satunya bagian
> > yang bahannya belum lengkap.**"*
>
> ⭐⭐ **Bahannya lengkap, dan selalu lengkap.** Isi 268 langkah ada di ekspor, di tag
> `PropertiesName`/`PropertiesValue` **tanpa awalan `py`**. **Karantina ini dibuka.**
> Bab dipertahankan sebagai catatan, **tidak dihapus**.

#### a. Apa yang SUDAH diketahui bentuknya

| Hal | Keadaan |
| --- | --- |
| ⭐ **medan** | delapan medan uang berpasangan mata uang — batas, retensi, premi bersih, pendapatan premi diperkirakan |
| ⭐ **sumber data** | ⭐ **kolom tersimpan** pada `POOLDATA.TREATYINDETAILJOINEDM` |
| ⭐ **arah aliran** | ⭐ **dibaca, bukan dihasilkan** — nol `Activity` atau `DataTransform` menghitungnya |
| ⭐ **penyajian** | ⭐ **apa adanya**, tidak dihitung ulang *(P43)* |
| ⭐ **pajak brokerage** | ⭐ rumusnya **terbaca penuh** — `@if(TypeTax=="Inclusive", Deduction/1.022, Deduction)` |
| ⭐ **persentase** | `DEDUCTION1` `DEDUCTION2` `BROKERAGE` `RNM_SHARE` |

⭐⭐ **Karena penyajian menyalin apa adanya, bagian TERBESAR uang tidak lagi menunggu P18.**

#### b. ⛔ ~~Apa yang MENUNGGU P18~~ — **RALAT: tidak ada lagi yang menunggu P18**

> ⛔ Bunyi semula, dikutip utuh:
> > *"⛔ **Isi 268 langkah `Property-Set` yang terjangkau** — di situlah seluruh aritmetika
> > tinggal: penjumlahan, pembagian bagian, pemulihan, angsuran, dan perhitungan antara."*

⭐⭐ **Terbaca seluruhnya.** `[terverifikasi]` **269 langkah** pada ke-51 aturan terjangkau,
**707 pasangan nama=nilai terisi**, **nol** tanda terpotong pada tujuh uji keutuhan. Termasuk yang
selama ini disebut belum diketahui: pembagi **102,2** untuk `TypeTax="Inclusive"`, **PPH 2 %**,
**PPN 2,2 %**, dan **presisi 8 angka** pada tiap `@divide`.

⚠️ **Dua hal berikut TETAP menunggu, dan keduanya bukan P18:**

⛔ Ditambah **`SumTSIPremiSpreadRNMMultiCob_Act`** — dipanggil, ⛔ **tidak ada di seluruh korpus
21 modul**, dan letaknya **di rantai uang** *(P30)*.

⛔ Ditambah **muatan yang dikirim ke layanan luar** di akhir alur — implementasinya ada di modul
lain dan **tidak boleh dipinjam** *(P8)*.

#### c. ⛔ ~~Akibatnya bila P18 TIDAK PERNAH dijawab~~ — **RALAT: keempatnya gugur**

| # | Akibat |
| ---: | --- |
| ⛔ **1** | ⛔ **Layar Dept Head tidak dapat disimpan.** `[terverifikasi]` **Sembilan medan wajib sekaligus terkunci** di sana — `Deduction2` · `ExcessLoss` · `OutstandingClaim` · `OveriddingCommOgp` · `OveriddingCommOnp` · `ResultOnp1` · `RiCommOgp` · `RiCommOnp` · `SalvageValue`. ⭐ Medan itu wajib diisi, ⛔ tidak dapat diisi pengguna, dan ⛔ **yang seharusnya mengisinya adalah perhitungan yang belum diketahui**. |
| ⛔ **2** | ⛔ Tangga persetujuan **berhenti di jenjang kedua** — berkas dapat naik ke Dept Head tetapi tidak dapat keluar dari sana. |
| ⛔ **3** | ⛔ Nilai kedelapan medan uang **hanya dapat ditampilkan, tidak dapat dihitung ulang** — cukup untuk membaca data lama, ⛔ **tidak cukup untuk membuat kontrak baru**. |
| ⚠️ **4** | ⚠️ Cabang `IsLife` di dalam kedua penjaga **tidak dapat dinyatakan** — isinya ada di dalam langkah yang tidak terekspor *(P38)*. |

> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — **keempat akibat di atas GUGUR SELURUHNYA.**
> Tabelnya dipertahankan sebagai catatan, tidak dihapus.
>
> | # | Akibat semula | Keadaan sebenarnya |
> | ---: | --- | --- |
> | 1 | layar Dept Head tidak dapat disimpan | ⭐ **dapat disimpan** — dari 34 medan wajib+terkunci, **28** punya langkah pengisi; 3 sisanya diketik underwriter di layar lain |
> | 2 | tangga berhenti di jenjang kedua | ⭐ **tidak berhenti** — akibat turunan dari nomor 1 |
> | 3 | delapan medan uang hanya dapat ditampilkan | ⭐ **dapat dihitung ulang** — rumusnya terbaca |
> | 4 | cabang `IsLife` tidak dapat dinyatakan | ⚠️ **terbaca**, tetapi **P38 tetap berdiri** sebagai pertanyaan arti, bukan pertanyaan bahan |
>
> ⛔ Kalimat penutup bab ini semula berbunyi:
> > *"⭐ **Jalan keluarnya hanya dua, dan keduanya bukan keputusan teknis:** ekspor ulang **dengan
> > isi langkah disertakan**, atau **pencabutan kewajiban** atas kesembilan medan itu oleh Product
> > & Underwriting."*
>
> ⭐⭐ **Tidak perlu keduanya.** Tidak ada ekspor ulang yang diminta, dan tidak ada kewajiban yang
> perlu dicabut.

---

## 9 · Butir `[terbuka]` — daftar penuh

⛔⛔ **Spec ini TIDAK menutup satu pun butir di bawah.** Penutupan milik pemilik perannya.

### 9.1 ⭐⭐ Yang MENAHAN — **NIHIL**

⭐⭐ **Tidak ada satu pun butir `[terbuka]` yang menahan pekerjaan.** `[terverifikasi]` 2026-09-22

> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — bab ini semula berjudul **"Satu yang MENAHAN"**
> dan memuat baris berikut, dikutip utuh:
> > *"| ⛔⛔ **P1** | naskah **tiga stored procedure** | `[DBA]` | penulisan lapisan penyimpanan |"*
>
> ⭐ **P1 terjawab work owner 2026-09-22** — naskah **empat** stored procedure diterima *(bukan
> tiga: `PEGA_DELETE_ERROR_KONVERSI` menyusul sebagai yang keempat)*, ditambah **dua contoh**
> `DATA_JSON`. Lihat `PERTANYAAN-untuk-DBA.md` §P1.

> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — bab ini semula berjudul **"Dua yang MENAHAN"**
> dan memuat satu baris lagi, dikutip utuh:
> > *"| ⛔⛔ **P18** | isi **268 langkah** penghitungan | `[pemilik export Pega]` | ⛔ rantai uang
> > **dan** kemampuan menyimpan layar Dept Head |"*
>
> ⭐ **P18 ditarik oleh tim migrasi, bukan dijawab pemilik export.** Tidak seorang pun di luar tim
> pernah menjawabnya; kekeliruannya milik tim migrasi. Lihat `VERIFIKASI-P18.md`.

### 9.2 Yang TIDAK menahan

| # | Butir | Pemilik |
| --- | --- | --- |
| 1 | pemetaan **nama → peran** untuk 12 tempat, ⛔ **berikut arahnya** pada dua layar | `[IAM]` · `[work owner]` |
| 2 | peran yang tersedia **belum cukup** untuk tiga posisi tangga | `[IAM]` |
| 3 | penetapan, pembuatan, dan pencabutan peran **belum dibahas sama sekali** | `[IAM]` |
| 4 | `SumTSIPremiSpreadRNMMultiCob_Act` — di bagian mana ia tinggal | `[pemilik export Pega]` |
| 5 | muatan yang dikirim ke layanan luar di akhir alur | `[Product+Underwriting]` |
| 6 | **tipe Oracle** kolom uang — sisi Pega menyimpannya sebagai **teks** | `[DBA]` |
| 7 | **format penyajian** uang — desimal dan pemisah ribuan | `[work owner]` |
| 8 | `RNM_SHARE` persentase **dari apa** | `[Product+Underwriting]` |
| 9 | ketidakseragaman presisi **4 lawan 8 desimal** untuk rumus yang sama | `[Finance]` |
| 10 | arti kode lini bisnis, termasuk `"40"` yang berskala penomoran berbeda | `[DBA]` |
| 11 | arti dua kode penanda "punya penempatan keluar" | `[Product+Underwriting]` |
| 12 | peran `SystemSetOneYear_DT` — kapan aturan satu-tahun berjalan | `[Product+Underwriting]` |
| 13 | **5 tempat** yang menghapus pesan galat **sesudah** validasi memasangnya | `[pengembang Pega lama]` |
| 14 | apakah ada **mata uang mati lain** yang tidak disembunyikan | `[Product+Underwriting]` |
| 15 | data lama bertanggal **ambigu secara mutlak** — `05/06` dapat berarti dua hal | `[work owner]` · `[DBA]` |
| 16 | berapa baris yang `OPERATORID`-nya **sudah terisi** sekarang | `[DBA]` |
| 17 | empat wadah layar berisi **104 medan** — masih dipakai atau ditinggalkan | `[Product+Underwriting]` |
| 18 | apakah ada kasus mencapai **Acceptance by Dir.** dalam 12 bulan terakhir | `[Product+Underwriting]` |
| 19 | `InsertToTreatyOutXOLList` **bernama** menulis tetapi tidak ada perintah tulis | `[pengembang Pega lama]` |
| ⭐ **20** | **P61** — tabel `_BACKUP` ikut terhapus bersama data utamanya; tidak ada pemulihan | `[work owner]` · `[DBA]` |
| ⭐ **21** | **P62** — tanggal ditulis mati di program penomoran memaksa periode Desember 2025 | `[work owner]` · `[Finance]` |
| ⭐ **22** | **sensus properti** `PolicyTreatyIn` belum dikerjakan — 94 nama unik dan 9 simpul bersarang baru **batas bawah**. ⛔ Diperberat: **bentuk dokumen mengikuti jenis treaty** — tiga contoh hanya berbagi **23 dari 74** medan, dan `ListInstallment` **bersarang pada satu kasus, datar pada kasus lain** | tim migrasi, ronde tersendiri |
| ⛔⛔ **23** | **galat angka uang tersimpan permanen di produksi** — ekor `2,76 x 10^-7` pada `NetPremium`, berlipat dari empat angsuran. Tiga pilihan *(ikuti apa adanya · bulatkan saat migrasi · bulatkan saat dihitung)* tercatat di **P29**; ⛔ **tim migrasi tidak memilih** | `[work owner]` · `[Finance]` |

> ⭐ **Empat baris terakhir BARU**, ditambahkan 2026-09-22 dari naskah procedure dan **tiga** contoh
> `DATA_JSON`. ⛔ **Keempatnya tidak menahan penulisan spec.** ⚠️ Butir **23** menahan **keputusan
> migrasi data lama** dan **ambang uji paritas** di tiket **15** — bukan penulisan lapisan
> penyimpanan.

---

## 10 · Further Notes

### 10.1 ⚠️ Enam penyimpangan sadar — dan kenapa masing-masing

⛔ **Setiap penyimpangan dari perilaku Pega dicatat beserta alasannya, sesuai AC 70.**

| # | Penyimpangan | Alasan |
| ---: | --- | --- |
| ⚠️ **1** | **Satu transaksi** menggantikan commit ganda | ⭐ sistem lama **tidak punya jaminan keutuhan**; ini lebih ketat, disengaja |
| ⚠️ **2** | **Berhenti** saat pembacaan gagal, bukan lanjut | ⭐ **kegagalan yang terlihat lebih murah** daripada yang tersembunyi |
| ⚠️ **3** | **Tanggal sebagai tipe tanggal**, satu format | ⛔ dua format tercampur dan disimpan sebagai teks membuat sebagian nilai **ambigu mutlak** |
| ⚠️ **4** | **`OPERATORID` mulai diisi** | ⭐ perbaikan jejak audit; kolomnya selama ini kosong |
| ⚠️ **5** | **`OperatorName` selalu nama tampilan** | ⭐ pengisian dari pengenal akun adalah **bug**, bukan maksud |
| ⚠️ **6** | **Peran menggantikan nama orang dan nomor urut antrean** | ⛔ keduanya **rapuh**: berhenti bekerja saat orang pindah atau daftar diurutkan ulang |

### 10.2 ⭐ Empat hal yang ditiru apa adanya, walau tampak salah

⛔ **Disebut supaya tidak "diperbaiki" diam-diam oleh pembangunnya.**

1. ⭐ **Tanggal akhir kosong ⇒ hari ini**, menghasilkan kontrak bermasa berlaku nol hari *(P35)*.
2. ⭐ **`TypeTax` dibandingkan teks persis** — huruf kecil atau spasi di belakang menggeser
   brokerage **2,2 %** tanpa pesan galat *(P46)*.
3. ⭐ **Nilai kosong pada penanda persetujuan berarti DISETUJUI** *(P24)*.
4. ⭐ **Kode jenis kontrak tanpa validasi daftar nilai** — nilai asing tetap diterima *(P37)*.

### 10.3 ⚠️ Pertentangan antar sumber yang ditemukan saat menulis spec ini

| Pertentangan | Yang menang | Sebab |
| --- | --- | --- |
| grilling ronde 2–4 menulis `Property-Set` **942**; keadaan menulis **938** | ⭐ **938** | ⭐ korpus **berubah** — satu berkas salah diperbaiki, selisihnya **185.967 B** tepat |
| grilling ronde 3 menulis `Page-Clear-Messages` **17**; keadaan **16** | ⭐ **16** | sebab yang sama |
| grilling ronde 4 menulis baris tabel keputusan **tidak terekspor** | ⛔ **keadaan** | ⭐ tag-nya **wadah**, bukan tag berteks — **36 baris ada** |
| grilling ronde 4 menulis **0** medan wajib pada sel medan | ⛔ **keadaan** | ⭐ `pyRequired` menempel lebih tinggi — **27 medan** |
| ronde 2 menulis seluruh **38** medan terkunci di layar Dept Head | ⛔ **keadaan** | ⭐ **36 di sana, 2 di layar biasa** |
| P22 *(ronde 2)* lawan P44–P45 *(ronde 4)* | ⭐ **P44–P45** | ⭐ P22 **ditutup** sebagai induk yang sudah terwakili |
| P4 lawan P49 | ⭐ **P4** | ⭐ P49 **rangkap**, ditulis tanpa mengetahui P4 sudah ada |

⭐ **Nol pertentangan** antara lembar jawaban *(sumber 1)* dan berkas keadaan *(sumber 3)*.
`[terverifikasi]` Ke-13 nomor pertanyaan yang menopang 12 ketetapan **seluruhnya sudah terjawab**.

### 10.4 ⚠️ Yang perlu diketahui sebelum modul berikutnya digarap

⭐ **Tiga temuan modul ini berlaku di modul lain**, dan sebaiknya tidak ditemukan ulang:

1. ⭐ **Pencarian antrean menurut nomor urut ada di 41 berkas pada 9 modul** — perubahan yang sama
   berlaku di sana.
2. ⭐ **`isApprovedtoDeptHead` dirujuk juga oleh `NB FacIn`** — tiga berkasnya merujuk rule yang
   sama, yang berkasnya tidak ada.
3. ⭐ **Rantai spreading/premi milik Fac In tampak di folder Treaty** karena layarnya ditanam
   bersama. ⛔ Folder modul **bukan** daftar pekerjaan.

---

## 11 · Lampiran

### 11.1 Angka yang berlaku

| Hal | Angka |
| --- | ---: |
| berkas dalam modul | **278** |
| ukuran modul | **35.492.317 B** |
| md5 korpus | `96271809212b68bde10ae8612f414baf` |
| `Activity` terjangkau / yatim | **64 / 28** |
| langkah `Property-Set` terjangkau / yatim | **377 / 561** |
| ⭐ langkah yang **diminta** sesudah penyaringan | ⭐ **268** |
| medan wajib | **27** di 6 layar |
| medan terkunci permanen | **38** — 36 + 2 |
| baris penggolong `BusinessType_DeT` | **36**, **128** kode |
| kolom view | **39**, **33/33** medan laporan tersedia |
| pertanyaan | **47 terjawab dari 49** |

### 11.2 Lingkup penelusuran — disebut sesuai disiplin

⭐ Seluruh angka di atas berasal dari **278 berkas `NB Treaty In`**, kecuali yang menyebut korpus:
**1.151 naskah SQL di 21 modul**. ⛔ `D:\XML\nusantara-re\` **tidak disentuh**.
⛔ Stored procedure, pemicu basis data, penjadwal, dan sistem di luar Pega **tidak tersisir** — ⭐
pernyataan nihil mana pun di berkas ini **hanya berlaku untuk lingkup itu**.

### 11.3 Berkas sumber

`PERTANYAAN-untuk-{DBA,Finance,IAM,Product-dan-Underwriting,pemilik-export-Pega,pengembang-Pega-lama}.md` ·
`VERIFIKASI-KEADAAN.md` · `KEADAAN-NB-TREATY-IN.md` · `grilling-ronde-1..4.md` *(tersegel)*

---

## TELEMETRI EKSEKUSI

⛔⛔ **Pengukuran dari luar TIDAK dilakukan.** ⭐ Spec ini ditulis **di dalam sesi interaktif**,
sehingga `total_cost_usd` dan `duration_ms` **tidak tersedia**. ⛔ **Tidak ditaksir.**

| Yang dicatat | Nilai | Sumber |
| --- | ---: | --- |
| ⭐ **token keluaran** | ⭐ **115.495** | ⚠️ **selisih dikurangkan**, lihat catatan di bawah |
| ⭐ **token cache-read** | ⭐ **15.580.830** | sama |
| token cache-write | **261.193** | sama |
| token masuk | **42** | sama |
| ⭐ **jumlah panggilan alat** | ⭐ **21** | sama |
| ⛔ durasi | ⛔ **TIDAK DAPAT DIUKUR** | butuh `--print` dari luar sesi |
| ⛔ biaya | ⛔ **TIDAK DAPAT DIUKUR** | butuh `--print` dari luar sesi |
| ⭐ **byte dibaca** | ⭐ **~150 KB** enam lembar jawaban + 3 berkas keadaan/verifikasi; korpus **tidak dibaca ulang** | spec ditulis dari sumber 1–3, bukan dari korpus |

⚠️ **Tiga keterbatasan pengukuran dari dalam:** *(a)* pesan terakhir belum tertulis ke transkrip
saat pengukuran kedua — angkanya **kurang satu pesan**; *(b)* jalur transkrip **diperiksa** dan
cocok dengan sesi ini; *(c)* yang dilaporkan **SELISIH**, bukan total sesi.

⚠️⚠️ **Cara angka di atas diperoleh, dan ia BUKAN pengukuran langsung.** ⛔ Baseline diambil
**sebelum ronde verifikasi**, bukan sebelum spec ini. ⭐ Yang terukur langsung adalah **gabungan
ronde verifikasi + penulisan spec**: keluaran **208.491** · cache-read **36.381.748** · panggilan
**56**. ⭐ Angka spec di tabel didapat dengan **mengurangkan** angka ronde verifikasi yang sudah
tercatat *(keluaran 92.996 · cache-read 20.800.918 · panggilan 35)*.
⚠️ **Pengurangan itu sah hanya bila tidak ada pekerjaan lain di antaranya** — dan memang tidak ada,
tetapi ⛔ **itu tidak diukur, melainkan diingat.** Dinyatakan apa adanya.

⭐ **Bandingkan dengan enam ronde sebelumnya** — seluruhnya selisih, alat yang sama:

| | R1 | R2 | R3 | R4 | Verifikasi | ⭐ **Spec** |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| keluaran | 192.021 | 225.898 | 178.857 | 154.518 | 92.996 | ⭐ **115.495** |
| cache-read | 13,5 jt | 22,7 jt | 16,7 jt | 26,2 jt | 20,8 jt | **15,6 jt** |
| panggilan | 69 | 72 | 39 | 47 | 35 | ⭐ **21** |

⭐ **Panggilan alat terendah dari seluruh ronde** — ⭐ sebab spec ini ditulis dari **jawaban yang
sudah ada**, bukan dari korpus. ⛔ **Korpus tidak dibaca ulang sama sekali** saat menulis bab 2–11.

