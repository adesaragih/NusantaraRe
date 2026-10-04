# Spec — Claim Prop (Klaim Treaty Inward Proporsional)

**Status:** ready-for-agent · **Tanggal:** 2026-09-18 · **Modul korpus:** `Claim Prop` (329 berkas)

> ### Ringkasan — 2026-09-19
>
> ```
> user story          66
> acceptance criteria  135
> ── sebaran penanda atas 135 AC ──────────────────────────
> [penyimpangan sadar]  48    ⚠️ total            72
> jebakan (⚠️ tanpa penggolongan)                 24
> AC ber-[terbuka] aktif  3    (AC 104 · 115 · 120)
> butir [terbuka] ringan  BELUM PUNYA DATA — lihat RALAT
> AC tanpa penanda        0
> ```
>
> ⛔ **RALAT 2026-09-19.** Angka lamanya **dikutip utuh, tidak dihapus**: *"acceptance criteria
> **129** · [penyimpangan sadar] **46** · ⚠️ total **67** · jebakan **21** · butir [terbuka] ringan
> **5**"*. **Enam AC ditambahkan** — **130** s/d **135**, dari lima keputusan work owner
> 2026-09-19 sesudah modul ini diadu dengan 15 ADR untuk pertama kalinya.
>
> **Dihitung dua cara, dan keduanya sepakat:** *(a)* pengurai bab AC yang memotong blok per nomor
> dan membuang baris kutipan; *(b)* rekonsiliasi **delta terhadap sensus lama yang tertulis di
> atas** — AC `129→135 (+6)` · penyimpangan `46→48 (+2)` · ⚠️ `67→72 (+5)` ·
> jebakan `21→24 (+3)`. Penambahnya: **133** dan **134** penyimpangan sadar; **130 · 131 · 132 ·
> 135** jebakan; dan **AC 72 kehilangan ⚠️**-nya karena penggolongannya sudah diputuskan.
>
> ⚠️⚠️ **`butir [terbuka] ringan` ditulis BELUM PUNYA DATA dengan sengaja — dan sebabnya BUKAN
> blok ini.** Ketiga sumber di dalam berkas ini **sudah berselisih sebelum hari ini**: kepala
> sensus menulis **5**, catatan di atas tabelnya menulis *"jumlah yang berlaku sejak 2026-09-19
> adalah **enam**"*, sedangkan **tabelnya sendiri memuat sembilan baris**. ⛔ Selisih itu **tidak
> diperbaiki di sini** — memilih salah satu tanpa dasar hanya menyembunyikannya. Yang pasti dan
> terbaca: **butir 5 tabel itu — penggolongan AC 72 — DITUTUP hari ini**, sehingga berapa pun
> angka yang benar, ia **berkurang satu**. Pemilik: **work owner**.
>
> ✅ **Kewajiban memakai `sensus.py` SUDAH DICABUT** — `[keputusan work owner]` **2026-09-19**.
> `[terverifikasi]` Berkasnya **pernah ada** dan kini **hilang tanpa salinan**; proyek ini bukan
> repositori git, jadi ia **tidak dapat dipulihkan**. Penggantinya ada di `CLAUDE.md` §4a
> *(Aturan sensus)*: ⭐ **tiap sensus dihitung DUA CARA yang berbeda, jendelanya disebut, dan bila
> keduanya berselisih ditulis "belum punya data" beserta kedua angkanya.**
>
> ⛔ **Kedua kalimat lama dikutip utuh, tidak dihapus:**
>
> > *"⛔ **Dan satu hal yang wajib disebut:** CLAUDE.md §4a mewajibkan setiap sensus dijalankan
> > lewat **`OUTPUT_HASIL_RNM\sensus.py`** … Aturannya **tidak dapat dipenuhi** … ⚠️ Selama
> > `sensus.py` belum ada, **§4a tidak dapat ditegakkan oleh siapa pun.**"*
> >
> > *"⚠️ **Sensus wajib lewat `OUTPUT_HASIL_RNM\sensus.py`**, bukan grep yang ditulis ulang tiap
> > kali — jendelanya tertanam di dalam kode, dicetak di setiap hasil, dan skripnya **menolak
> > jalan tanpa uji instrumen**. Lihat `CLAUDE.md`, bagian *Sensus wajib lewat sensus.py*."*
>
> ⚠️ **Yang TIDAK ikut gugur: jendelanya.** Jendela baku di bawah **tetap berlaku** dan sama persis
> dengan yang dulu dicetak skrip itu — yang gugur **alatnya**, bukan **ukurannya**.
>
> ⚠️ **RALAT 2026-09-19** — angka berubah karena **lima AC ditambahkan sekaligus** (baris
> penyesuaian yang beku · hapus klaim · kolom `FLAG_ON_GOING_COMMITTEE` yang dibuang):
> AC 124 → **129**, tanda `[penyimpangan sadar]` 41 → **46**, ⚠️ total 61 → **67**.
> ⚠️ **Jebakan 20 → 21** — bukan karena AC baru: kelima AC baru bertanda
> `[penyimpangan sadar]`. Yang bertambah adalah **AC 5**, yang ralatnya membawa ⚠️ tanpa
> tanda penyimpangan, sehingga menurut definisi bab ini ia terhitung **jebakan**. Angka ini
> **dihitung ulang, tidak dikira-kira**. AC ber-`[terbuka]` tetap **3**, butir `[terbuka]`
> ringan tetap **5**,
> AC tanpa penanda tetap **0**. ⛔ **Nol butir `[terbuka]` ditutup.** Bukan penyimpangan sensus;
> penyebabnya tercatat di **14c** dan di kelima AC itu sendiri.
>
> ⚠️ **RALAT 2026-09-18** — angka berubah karena **AC 124** ditambahkan (urutan pemanggilan
> procedure penulis): AC 123 → **124**, tanda `[penyimpangan sadar]` 40 → **41**, ⚠️ total 60 →
> **61**. Jebakan tetap **20**, butir `[terbuka]` ringan tetap **5**, AC ber-`[terbuka]` tetap
> **3**. Bukan penyimpangan sensus; penyebabnya tercatat di AC 124 sendiri.
>
> ⚠️ **RALAT 2026-09-18** — angka berubah karena **AC 123** ditambahkan (migrasi memilih baris yang
> berlaku): AC 122 → **123**, tanda `[penyimpangan sadar]` 39 → **40**, ⚠️ total 59 → **60**.
> Jebakan tetap 20. Bukan penyimpangan sensus; penyebabnya tercatat di AC 123 sendiri.

> ### Status `[terbuka]` — 2026-09-17
>
> ```
> [terbuka] yang memblokir : NIHIL
> [terbuka] ringan         : satu — enam kolom STS_* di roster komite
> ```
>
> Ketujuh `[terbuka]` pada terbitan pertama sudah dijawab; rinciannya di bab
> **Pertanyaan terbuka di dalam spec**. Yang ringan tidak memblokir: perilaku yang ditiru adalah
> perilaku yang berjalan.
>
> #### ⚠️ RALAT 2026-09-17 — angka korpus berubah
>
> Korpus naik dari **328** menjadi **329** berkas; `When/` dari **60** menjadi **61**, karena
> `When/IsCustomBonds.xml` ditambahkan. Seluruh sensus yang menyebut 328 atau 60 sudah dijalankan
> ulang, dengan angka lama dan baru berdampingan.
>
> `[terverifikasi]` Akibat langsungnya: nama rule `When` yang **dirujuk tetapi berkasnya tidak ada**
> turun dari **1** menjadi **nol**. `IsCustomBonds` bukan salah ketik — ia rule hilang keenam, dan
> kini ada.
>
> ⚠️ Sensus **rule `When` yatim** (§25d ronde 2: 14 dari 60) **tidak** ikut diperbarui di sini.
> Pemeriksaan cepat atas 61 berkas memberi angka berbeda dengan metode yang berbeda pula, sehingga
> menimpanya akan menukar satu angka yang tidak terverifikasi dengan angka lain yang juga belum.
> **Dibiarkan sebagaimana tercatat; bukan pemblokir tiket.**

> Korpus `D:\XML\RNM_BRD\` **READ-ONLY**. Semua tulisan hanya ke `OUTPUT_HASIL_RNM\`.
> ⛔ `D:\XML\nusantara-re\` di-blacklist. Target: **Pega → Go + React + Oracle** (DB tetap).
> Bukti = **path berkas dari akar modul + nama rule + nomor step Pega**. Nomor baris XML tidak pernah
> dikutip.
> Tanda: `[terverifikasi]` (terbukti korpus) · `[keputusan work owner]` · `[data DBA]` · `[terbuka]` ·
> ⚠️ **penyimpangan sadar**.
>
> #### ⚠️ RALAT 2026-09-18 — penggolongan pindah dari prosa ke tanda berkurung
>
> ⚠️ dipakai untuk **dua hal** yang keduanya sah: **penyimpangan sadar** (Go sengaja tidak meniru
> Pega) dan **jebakan** (perilaku ditiru apa adanya, tetapi mudah salah dibaca — nama yang
> berbohong, arah gerbang terbalik, dua enum bernilai sama). Selama tiga ronde penggolongannya
> bergantung pada frasa di dalam prosa, yang **bisa terlipat di pergantian baris** dan karenanya
> lolos dari sensus. Mulai sekarang setiap AC penyimpangan sadar membawa tanda
> **`[penyimpangan sadar]`** tepat sesudah nomornya, sebelum ⚠️ — posisi yang tidak pernah
> terlipat. Sensusnya menjadi satu perintah dan eksak.
>
> **Jendela sensus baku.** Satu unit AC = baris `<n>. ` ditambah seluruh baris lanjutannya sampai
> baris bernomor berikutnya atau sebuah judul. **Kecualikan** setiap baris yang diawali `>` — itu
> blok kutipan RALAT berisi teks lama yang sudah dibatalkan. **Normalkan** sebelum mencari: gabung
> jadi satu string, rapatkan spasi ganda, buang penanda markdown `*` `` ` `` `_`.
> ⚠️ **Tidak ada PII di berkas ini.** Korpus memuat nama orang, alamat email, operator ID, dan IP
> internal; seluruhnya dirujuk secara deskriptif.

**Sumber:** `.scratch/claim-prop/grilling-ronde-1.md` dan `grilling-ronde-2.md`. **Ronde 2 menang**
bila bertentangan — ia memuat ralat atas ronde 1.

**ADR yang mengikat:** ADR-0003 (uang non-float) · ADR-0006 (penomoran lewat stored procedure) ·
ADR-0011 (unit keputusan = baris `AdjustmentList`) · ADR-0013 (endpoint lewat `M_LINK_SERVICE`) ·
ADR-0014 (keputusan komite ditegakkan di **lapisan layanan** — ⚠️ cakupan asli ADR-nya Komite Claim
Life; yang dipakai di sini **prinsipnya**) · ADR-0015 (efek keluar wajib berhasil).

---

## Problem Statement

Penanganan klaim treaty inward proporsional hari ini berjalan di atas Pega, dan **data klaim intinya
tidak pernah ada sebagai baris Oracle**. `[terverifikasi]` Peserta, baris adjustment, estimasi,
spreading, dan interest hidup sebagai *page list* di dalam work object Pega dan dipersist lewat
`Obj-Save` ke blob; yang menyentuh Oracle hanyalah **proyeksi** — `OS_AKSEPTASI_KLAIM` dan
`JSON_KLAIM`, keduanya hibrida satu kolom dokumen `DATA_JSON` plus beberapa kolom kunci.

Akibatnya bagi pengguna dan bagi perusahaan:

- **Angka klaim tidak dapat ditanyakan.** Tidak ada tabel yang bisa di-`SELECT` untuk menjawab
  "berapa outstanding per treaty" tanpa membongkar JSON.
- **Nomor terbakar.** `[terverifikasi]` Lima generator nomor `COMMIT` sendiri di dalam procedure,
  sebelum pemanggilnya selesai — bila langkah sesudahnya gagal, nomor sudah terpakai dan tidak dapat
  ditarik. `Activity/TryMakePLA_Act.xml` bahkan **tidak punya `Obj-Save` sama sekali** setelah
  sequence diambil.
- **Presisi rusak sebelum disimpan.** `[data DBA]` Nilai turunan dihitung dari angka yang **sudah
  dipotong 4 desimal** di hulu, lalu dirambatkan.
- **Tidak ada satu pun gerbang wewenang di sisi server.** `[terverifikasi]` Nol `pyValidateActivity`
  di 12 dari 12 FlowAction; nol pemeriksaan peran di Activity, 418 ekspresi Section, 19
  ReportDefinition, dan 11 Harness. Gerbang penyerahan ke Komite sepenuhnya *client-side*.
- **Jejak audit bocor dan salah label.** Aksi satu peran tidak terekam sama sekali; memo penutupan
  klaim tidak pernah sampai ke jejak.

## Solution

Membangun ulang siklus klaim treaty inward proporsional di Go + React di atas Oracle yang sama,
dengan **skema relasional yang dirancang baru** (belum pernah ada), **satu jalur perhitungan uang**,
**batas transaksi yang eksplisit untuk penomoran**, dan **penegakan wewenang di lapisan layanan**.

Perilaku bisnis **ditiru apa adanya** — termasuk cacat yang sudah diterima sadar — kecuali pada
titik-titik yang ditandai ⚠️ **penyimpangan sadar**.

### Aturan induk — tiru PERILAKU, bukan NIAT

`[keputusan work owner]` Berlaku untuk **seluruh** spec ini:

| Hal di Pega | Bentuk di Go | Cakupan `[terverifikasi]` |
| --- | --- | --- |
| Gerbang tertulis tetapi flag `pyStepsPreCondition` tidak `true` | **WHEN tidak ditulis; STEP tetap ditulis** dan berjalan tanpa saringan | **84** di Activity · **57** gerbang UI di Section |
| Langkah ber-`pyStepsBlockName` = `//` | **STEP tidak ditulis sama sekali** | **68** step |
| Elemen UI bersyarat selalu-salah (`NEVER`, `1=2`, `1=3`) | **Elemen tidak dibuat** | **75** elemen |

⚠️ **Bedanya:** syarat yang tertinggal karena Pega tidak menghapusnya bersih **tidak membawa niat** →
diabaikan. Sedangkan `NEVER` / `1=2` yang **diketik developer** adalah **niat menyembunyikan** →
elemennya tidak dibuat. Keduanya berujung tidak-ditulis, alasannya berbeda.

---

## User Stories

### Registrasi dan validasi masuk

1. Sebagai **Claim Admin**, saya ingin mendaftarkan klaim baru dengan nomor polis treaty, sehingga
   klaim tercatat dan masuk antrean kerja.
2. Sebagai **Claim Admin**, saya ingin sistem menolak nomor polis yang tidak ada di master treaty,
   sehingga klaim tidak menggantung tanpa induk.
3. Sebagai **Claim Admin**, saya ingin diperingatkan bila ada klaim lain dengan Date of Loss yang
   sama pada polis yang sama, sehingga saya tidak membuat klaim ganda.
4. Sebagai **Claim Admin**, saya ingin peringatan duplikat itu **tidak muncul** bila klaim lama sudah
   ditolak, sehingga klaim pengganti tetap bisa dibuat.
5. Sebagai **Claim Admin**, saya ingin sistem memeriksa apakah premi polis sudah lunas sebelum
   akseptasi, sehingga klaim tidak diproses di atas polis yang belum dibayar.
6. Sebagai **Claim Admin**, saya ingin pemeriksaan premi itu dapat **dibebaskan** lewat data proteksi,
   sehingga kasus yang sudah disetujui manajemen tetap bisa jalan.
7. Sebagai **Claim Admin**, saya ingin mengisi Cause of Loss dari master dua tingkat, sehingga
   penyebab kerugian tercatat seragam.
8. Sebagai **Claim Admin**, saya ingin sistem menolak simpan bila Cause of Loss belum diisi, sehingga
   tidak ada klaim tanpa penyebab.
9. Sebagai **Claim Admin**, saya ingin menandai klaim sebagai katastrofa atau bukan, dan
   mengaitkannya ke satu **event katastrofa**, sehingga beberapa klaim dari satu peristiwa dapat
   dikelompokkan.

### Insured Interest (objek pertanggungan / TSI)

10. Sebagai **Claim Admin**, saya ingin mendaftar objek pertanggungan beserta nilai TSI-nya, sehingga
    dasar perhitungan klaim jelas.
11. Sebagai **Claim Admin**, saya ingin tiap objek punya **mata uang dan kursnya sendiri**, sehingga
    klaim multi-mata-uang dapat ditangani.
12. Sebagai **Claim Admin**, saya ingin melihat total TSI per mata uang dan totalnya dalam IDR,
    sehingga saya tahu plafon klaim.
13. Sebagai **Claim Admin**, saya ingin daftar objek wajib terisi sebelum simpan, sehingga klaim
    tidak berjalan tanpa dasar nilai.

### Estimasi

14. Sebagai **Claim Admin**, saya ingin sistem membentuk baris estimasi **per treaty** dari hasil
    loss allocation, sehingga estimasi terbagi sesuai penempatan risiko.
15. Sebagai **Claim Admin**, saya ingin tiap baris estimasi punya mata uang dan kursnya sendiri,
    sehingga estimasi mengikuti mata uang objeknya.
16. Sebagai **Claim Admin**, saya ingin tanggal estimasi dibatasi antara Date of Loss dan hari ini,
    sehingga tidak ada estimasi bertanggal mustahil.
17. Sebagai **Claim Admin**, saya ingin diperingatkan bila total estimasi melebihi TSI, sehingga
    estimasi tidak melampaui pertanggungan.
18. Sebagai **Claim Admin**, saya ingin diperingatkan bila total estimasi melebihi plafon cash call,
    sehingga batas pembayaran dini terjaga.
19. Sebagai **Claim Admin**, saya ingin menghapus baris estimasi dan melihat seluruh total ikut
    menyesuaikan, sehingga angka tetap konsisten.

### Loss Allocation

20. Sebagai **Claim Admin**, saya ingin membagi nilai kerugian ke beberapa treaty menurut
    persentase share, sehingga tiap treaty menanggung porsinya.
21. Sebagai **Claim Admin**, saya ingin pembagian itu **disegmentasi per mata uang**, sehingga
    klaim multi-mata-uang tidak tercampur.
22. Sebagai **organisasi**, saya ingin total share pada satu mata uang **tepat 100 %**, sehingga
    tidak ada nilai kerugian yang hilang atau terhitung dua kali.
23. Sebagai **Claim Admin**, saya ingin ditolak saat total share **kurang dari** 100 % maupun **lebih
    dari** 100 %, sehingga kesalahan tertangkap di kedua arah.

### Baris adjustment

24. Sebagai **Claim Admin**, saya ingin menambahkan baris adjustment pada klaim, sehingga unit
    keputusan klaim terbentuk (**ADR-0011**).
25. Sebagai **Claim Admin**, saya ingin tiap baris adjustment punya statusnya sendiri, sehingga satu
    klaim dapat memuat beberapa keputusan berbeda.
26. Sebagai **Claim Admin**, saya ingin adjuster dan consultant wajib terisi sebelum baris adjustment
    **pertama** dibuat, sehingga penanganan klaim punya penanggung jawab.
27. Sebagai **Claim Admin**, saya ingin mengisi data bank penerima pada baris adjustment, sehingga
    pembayaran dapat diproses.
28. Sebagai **Claim Admin**, saya ingin nilai adjustment dibandingkan terhadap estimasi, sehingga
    saya tahu bila melampaui.
29. Sebagai **Claim Admin**, saya ingin menghapus baris adjustment yang salah, sehingga klaim dapat
    diperbaiki.

### Spreading

30. Sebagai **organisasi**, saya ingin nilai klaim tersebar ke treaty dan ke breakdown quota share,
    sehingga porsi tiap pihak terhitung.
31. Sebagai **organisasi**, saya ingin spreading pada baris adjustment terpisah dari spreading di
    tingkat klaim, sehingga penyelesaian dan estimasi tidak tercampur.
32. Sebagai **Finance**, saya ingin jumlah seluruh baris spreading **sama persis** dengan nilai
    induknya, sehingga tidak ada selisih yang menggantung.

### Deductible

33. Sebagai **Claim Admin**, saya ingin deductible dihitung sebagai **yang terbesar** antara
    persentase dikali basis dan nilai flat, sehingga ketentuan polis terpenuhi.
34. Sebagai **Claim Admin**, saya ingin memilih basis deductible — dari TSI atau dari nilai klaim —
    sehingga sesuai jenis kontraknya.
35. Sebagai **Claim Admin**, saya ingin deductible dikurangkan **sesudah** share ceding diterapkan,
    sehingga angka akhir konsisten dengan praktik yang berlaku.

### Penyerahan ke Komite

36. Sebagai **Claim Admin**, saya ingin menyerahkan baris adjustment ke Komite untuk diputuskan,
    sehingga keputusan di atas wewenang saya diteruskan.
37. Sebagai **organisasi**, saya ingin jumlah tingkat komite dihitung dari roster aktif yang pita
    limitnya mencakup nilai klaim, sehingga tangga persetujuan mengikuti data, bukan konstanta.
38. Sebagai **Claim Admin**, saya ingin ditolak menyerahkan bila data bank belum lengkap, sehingga
    keputusan tidak diambil di atas data yang belum siap.
39. Sebagai **Claim Admin**, saya ingin ditolak menyerahkan bila lampiran wajib belum lengkap,
    sehingga Komite punya bahan yang cukup.
40. Sebagai **organisasi**, saya ingin penyerahan ke Komite **ditolak di lapisan layanan**, bukan
    hanya disembunyikan di layar, sehingga tidak dapat dilewati lewat API.
41. Sebagai **Claim Admin**, saya ingin menutup klaim **tanpa pembayaran** lewat persetujuan pihak
    berwenang, sehingga klaim yang tidak dibayar tetap punya jejak keputusan.

### Penutupan klaim

42. Sebagai **Claim Admin**, saya ingin ditolak menutup klaim bila masih ada baris adjustment yang
    belum diputus, sehingga klaim tidak tertutup setengah jalan.
43. Sebagai **Claim Admin**, saya ingin ditolak menutup klaim bila pengiriman ke Kasir belum
    berhasil, sehingga pembayaran tidak tertinggal.
44. Sebagai **Claim Admin**, saya ingin memo alasan penutupan **benar-benar tercatat** di riwayat
    klaim, sehingga keputusan dapat ditelusuri.

### Dokumen

45. Sebagai **Claim Admin**, saya ingin mencetak **PLA** pada tahap outstanding, sehingga pemberitahuan
    awal kerugian terkirim.
46. Sebagai **Claim Admin**, saya ingin mencetak **DLA** (*Definite Loss Advise*) per reinsurer per
    baris adjustment, sehingga tiap reinsurer menerima rincian porsinya.
47. Sebagai **Claim Admin**, saya ingin Acceptance Note tercetak otomatis saat akseptasi disimpan,
    sehingga bukti persetujuan selalu ada.
48. Sebagai **Claim Admin**, saya ingin dokumen tersimpan dan dapat dibuka kembali lewat tautan,
    sehingga tidak perlu mencetak ulang.
49. Sebagai **organisasi**, saya ingin nomor dokumen **tidak terbakar** saat pembuatan dokumen gagal,
    sehingga urutan nomor tetap rapat.

### Efek keluar

50. Sebagai **organisasi**, saya ingin data akseptasi terkirim ke Kasir, sehingga pembayaran dapat
    diproses.
51. Sebagai **Claim Admin**, saya ingin melihat status pengiriman ke Kasir, sehingga saya tahu apakah
    perlu tindakan.
52. Sebagai **organisasi**, saya ingin efek keluar yang gagal **diantre dan diulang**, sehingga tidak
    hilang diam-diam (**ADR-0015**).
53. Sebagai **organisasi**, saya ingin klaim tersinkron ke sistem inti reinsurance non-life, sehingga
    pembukuan sejalan.
54. Sebagai **Claim Admin**, saya ingin email pemberitahuan terkirim ke anggota komite yang berwenang,
    sehingga keputusan tidak tertunda.

### Wewenang

55. Sebagai **organisasi**, saya ingin tingkat wewenang setiap pelaku dibaca dari **satu sumber
    data**, sehingga tidak ada dua daftar yang berbeda.
56. Sebagai **organisasi**, saya ingin pelaku yang tidak terdaftar di roster komite otomatis
    diperlakukan sebagai **Claim Admin**, sehingga tidak ada tingkat yang menggantung.
57. Sebagai **organisasi**, saya ingin batas nilai Direktur Utama diperlakukan sebagai wewenang
    **terpisah**, bukan tingkat kelima tangga komite, sehingga dua hal berbeda tidak tercampur.
58. Sebagai **administrator**, saya ingin mengganti pemegang jabatan cukup dengan mengubah baris
    tabel, sehingga tidak perlu menyentuh kode.

### Jejak audit

59. Sebagai **auditor**, saya ingin setiap aksi pada klaim meninggalkan entri riwayat berisi aksi,
    pelaku, waktu, dan tingkat wewenang, sehingga klaim dapat direkonstruksi.
60. Sebagai **auditor**, saya ingin riwayat tampil **urut waktu**, sehingga urutan kejadian terbaca
    benar.
61. Sebagai **auditor**, saya ingin **semua** peran terekam tanpa kecuali, sehingga tidak ada aksi
    yang tak berjejak.
62. Sebagai **auditor**, saya ingin tingkat wewenang tersimpan sebagai nilai terbatas dan nama pelaku
    di kolom terpisah, sehingga riwayat dapat disaring dan dikelompokkan.

### Penyimpanan dan pelaporan

63. Sebagai **Finance**, saya ingin menanyakan outstanding klaim lewat query biasa, sehingga laporan
    tidak perlu membongkar JSON.
64. Sebagai **organisasi**, saya ingin data klaim lama ikut termigrasi, sehingga riwayat tidak putus.
65. Sebagai **organisasi**, saya ingin satu aksi pengguna menghasilkan **satu transaksi**, sehingga
    tidak ada keadaan setengah tersimpan.

### Lintas-lini

66. Sebagai **organisasi**, saya ingin modul ini tetap melayani tiga prefix klaim beserta varian
    Syariah-nya, sehingga tidak ada lini yang tertinggal.

---

## Implementation Decisions

### 1. ⚠️ PREFACTOR — skema relasional dirancang baru

`[terverifikasi]` **Skema relasional Claim Prop tidak dapat direkayasa-balik dari korpus — ia belum
pernah ada.** Bukti: `Activity/DeleteAjsutment_Act.xml` nol `RDB-List`; `Activity/AddAdjustment_Act.xml`
dan `Activity/AddEstimation_Act.xml` hanya memanipulasi page dan melakukan lookup kurs; persistensi
terjadi lewat `Obj-Save` atas `pyWorkPage` di `Activity/SaveOutstanding_Act.xml`.

**Bahan perancangan yang tersedia:**

| Sumber | Isi |
| --- | --- |
| `[data DBA]` bentuk JSON nyata — satu baris contoh | nama dan tipe properti, ronde 2 §23 |
| Hierarki **45 pasangan page-list** | `pyPageListProperty` + `pyPageListPropertyClass` di 36 dari 36 Section, ronde 2 §14b |
| Sensus kolom proyeksi | atribut `DATA_JSON` yang terbaca, ronde 2 §8h |

**Hierarki yang terbukti definisional** `[terverifikasi]` — dua akar:

```
ClaimData ─┬─ EstimationList         (estimasi, per treaty)
           ├─ InterestList           (objek pertanggungan / TSI)
           ├─ InterestListDtl        (salinan tampilan)
           ├─ AdjustmentList         (unit keputusan — ADR-0011)
           │     ├─ SpreadingAdjustment
           │     ├─ SpreadingQuotaShare
           │     └─ LossAllocation   (snapshot dari SpreadingRisk)
           ├─ ListTotalEstimation · ListClaimAmount · TotalInterestInsured  (subtotal per mata uang)
           ├─ SuggestList            (jejak audit)
           ├─ ComiteeClaim           (roster komite pada kasus)
           └─ SpreadingRisk · SpreadingClaim · SpreadingBreakQS            (loss allocation & spreading)

PaymentData ─┬─ AcceptationList
             └─ CurrencyList ─ DetailPayment
```

⚠️ `[terverifikasi]` **Enam properti spreading berbagi satu class** (`Data-SpreadingRisk`). Di Oracle
mereka **tabel berbeda** dengan peran berbeda — jangan disatukan hanya karena class-nya sama.

⚠️ `[terverifikasi]` `SuggestList` menumpang class milik modul lain (`Data-OfferFacIn-SuggestList`).
Di Go ia **tabel jejak audit milik Claim Prop sendiri**.

#### 1a. Tanda tangan procedure penulis, dan arti ketujuh kolom penandanya

> ⚠️ **RALAT 2026-09-18 — butir `[terbuka]` ke-6 DITUTUP.** Arti ketujuh kolom penanda pada
> `OS_AKSEPTASI_KLAIM` dan `JSON_KLAIM` sebelumnya tercatat sebagai terbuka. Terjawab dari korpus
> ditambah tanda tangan procedure `[data DBA]`; jawabannya di bawah.

`[terverifikasi]` **Jendela sensus:** **329 berkas**, seluruh teks, dan penelusuran **diikat ke
halaman pemanggilnya** — `InputData.CARI<n>`, bukan `CARI<n>` telanjang. ⚠️ `TempKasir.CARI16` dan
`TempCFS.CARI20` milik halaman lain dan memberi **jawaban palsu**; itu jebakan lingkup yang sama
seperti `DataChronology.CARI1` pada jejak audit. *(Brief menyebut 328; sensus dijalankan atas 329,
angka yang tercatat di kepala spec ini.)*

`[data DBA]` Tanda tangan `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM`, dipetakan **posisional** ke
pemanggilnya `RDBList/SaveOSClaim_SQL.xml`:

```
(NoClaim, PolisNO, PegaID, DataPega, STATUS_RJT, KONVERSI, STS_PLA, IDMASTER, OldClaim,
 ErrMsg OUT, StsSimpan OUT)
   ↑CARI29  ↑CARI3  ↑CARI2  ↑CARI1   ↑CARI10    ↑CARI16   ↑CARI17  ↑CARI18  ↑CARI20
```

| Kolom | Parameter | Diisi dari | Nilai yang benar-benar masuk | Tanda |
| --- | --- | --- | --- | --- |
| `STS_REJECT` | `STATUS_RJT` | `InputData.CARI10` ← `TempOSAkseptasi.Type` | **0** Save Outstanding · **2** Acceptation · **4** Close Claim | `[terverifikasi]` |
| `MASTERID` | `IDMASTER` | `InputData.CARI18` | ⚠️ **dua sumber** — `ClaimData.IDMaster` dari `SaveOutstanding_Act`, `TreatyInMaster.ID` dari `CloseClaimProp` | `[terverifikasi]` |
| `STS_KONVERSI` | `KONVERSI` | `InputData.CARI16` | ⚠️ **NOL PENULIS** di 329 berkas — masuk **kosong** | `[terverifikasi]` |
| `STS_DLA` | `STS_PLA` | `InputData.CARI17` | ⚠️ **NOL PENULIS** — masuk **kosong** | `[terverifikasi]` |
| `CLAIMOLD` | `OldClaim` | `InputData.CARI20` | ⚠️ **NOL PENULIS** — masuk **kosong** | `[terverifikasi]` |
| `IDPROD` | — | hardcode di `PEGA_JSON_KLAIM_PNC` | **1** | `[data DBA]` |
| `TGL_KONVERSI` | — | hardcode | **NULL** | `[data DBA]` |

**Empat temuan yang menempel:**

1. ⚠️ `[terverifikasi]` **Status konversi tidak dibaca dari kolomnya.**
   `RDBList/getStatusKonversi_SQL.xml` berbunyi
   `select COUNT(1) … from reinsurance.trloss_detail_t WHERE NO_AKSEP = …` — konversi diperiksa
   dengan **menghitung baris di sistem tujuan**, bukan membaca `STS_KONVERSI`. Kolomnya bukan hanya
   kosong, ia juga **tidak dipakai**.
2. ⚠️ `[data DBA]` **Parameter bernama `STS_PLA` mengisi kolom `STS_DLA`.** PLA dan DLA **dua
   dokumen berbeda**. Masuk daftar nama yang berbohong — lihat **§16**.
3. ⚠️ `[terverifikasi]` **`NOCLAIM` diisi dua kali di activity yang sama.**
   `Activity/SaveOutstanding_Act.xml` menulis `InputData.CARI29` dari `ClaimData.NoClaim` **lalu**
   dari `ClaimData.ClaimNo`. **Yang terakhir menang.**
4. `[terverifikasi]` **Dua jalur, satu tabel.** `SaveOSClaim_SQL` hanya dipanggil dari
   `Activity/SaveOutstanding_Act.xml`, dan `SaveDataToOsAkseptasiNP` hanya dari
   `Activity/CloseClaimProp.xml`.

#### 1b. ⚠️ Tabel proyeksi diikuti APA ADANYA — bukan penyimpangan

`[keputusan work owner]` **2026-09-18.** Ketujuh penanda pada `OS_AKSEPTASI_KLAIM` dan `JSON_KLAIM`
**ditiru apa adanya. Tidak ada yang dianggap cacat:**

- **`STS_KONVERSI` · `STS_DLA` · `CLAIMOLD`** — kolomnya **tetap ada** dan **tetap ditulis kosong**.
  Nol penulis di Pega, nol penulis di Go. **Bukan cacat; jangan diisi, jangan dibuang.**
- **`IDPROD = 1`** dan **`TGL_KONVERSI = NULL`** — hardcode **ditiru**.
- Parameter **`STS_PLA` yang mengisi kolom `STS_DLA`** — nama **dibiarkan**.
- **`NOCLAIM` yang ditulis dua kali** dengan yang terakhir menang — **ditiru**.
- **`MASTERID` yang bersumber beda menurut jalurnya** — **ditiru**.

⚠️ **Ini BUKAN penyimpangan sadar.** Tidak ada tanda penyimpangan pada butir-butir ini, dan tidak ada
AC baru yang ditambahkan untuknya.

### 2. ⚠️ Empat schema Oracle, dan mayoritas tabel tanpa prefix

`[terverifikasi]` **EMPAT** schema, bukan tiga: `pooldata` (mayoritas) · `arasapas`
(`CekLunasPremi_Sql`) · `reinsurance` (`getStatusKonversi_SQL`) · **`gl`** (`gl.f_get_email` di
`RDBList/GetEmailCeding_SQL.xml`).

⚠️ `[terverifikasi]` **Mayoritas tabel tidak diprefiks schema sama sekali** — resolusinya bergantung
default schema koneksi. Kasus terburuk: `RDBList/CekLunasPremi_Sql.xml` memakai `arasapas.invoice`
berdampingan dengan `detail_invoice` **tanpa prefix**. **Di Go setiap objek diprefiks eksplisit.**
⚠️ **penyimpangan sadar** — memperbaiki resolusi nama, bukan meniru kerapuhannya.

### 3. ⚠️ Batas transaksi penomoran — keputusan eksplisit

`[terverifikasi]` **15 dari 58** rule Connect-SQL `COMMIT` sendiri di dalam `<pyBrowseSQL>`, dan
**lima di antaranya generator nomor**: `RDBList/GetSequenceNumber_SQL.xml` ·
`GenerateNoCLMTreatyIn.xml` · `GenerateNoPLATreatyIn.xml` · `GenerateNoDlaTreatyIn.xml` ·
`GenerateNoTRTInTemp.xml`.

**Akibatnya nomor terbakar begitu di-generate.** Gejalanya terlihat di
`Activity/TryMakePLA_Act.xml`: sequence diambil di step **11**, hasilnya hanya ditaruh di clipboard
di step **14**, step **26** dapat melempar exception, dan **rule itu tidak punya `Obj-Save` sama
sekali**.

⚠️ **penyimpangan sadar — keputusan:** di Go, **pengambilan nomor dan penyimpanan hasilnya berada
dalam satu transaksi**. Bila penyimpanan gagal, transaksi *rollback* dan nomor tidak terpakai.
Logika pembentukan nomor **tetap di stored procedure** dan **tidak direplikasi** (**ADR-0006**);
yang berubah hanya **di mana `COMMIT` terjadi** — dipindahkan ke lapisan aplikasi.

`[terverifikasi]` Tanda tangan keempat procedure penomoran:

| Procedure | Masuk | Keluar |
| --- | --- | --- |
| `GENERATE_NOCLMTREATYIN` | `KODE, KODE_BIS, BULAN, TAHUN, TIPE` | `START_DATE, TSI` |
| `GENERATE_NOPLATREATYIN` | sama | sama |
| `GENERATE_NODLATREATYIN` | sama | sama |
| `GENERATE_NOCLMTRTYINTEMP` | `KODE, TAHUN, TIPE` | sama |

⚠️ **Nomor temp tidak dapat dipromosikan menjadi nomor final** — ia kekurangan `KODE_BIS` dan
`BULAN`.

### 4. ⚠️ Uang — satu jalur perhitungan

`[keputusan work owner]` **Simpan penuh tanpa pembulatan. Pembulatan hanya untuk tampilan.**
Semua nilai uang bertipe **desimal presisi arbitrer**, tidak pernah melewati *binary floating point*
(**ADR-0003**).

`[terverifikasi]` Hari ini `.ClaimSpreaded` ditulis di **20 titik, 7 activity, dengan 4 perlakuan
pembulatan berbeda** (`/100` polos · `@divide(…,100,4)` · `,100,10` · `,100,20`), dan nilainya
**ditimpa berulang** oleh rule berbeda — nilai akhir bergantung urutan pemanggilan, bukan kebijakan.

⚠️ **penyimpangan sadar — keputusan:**

1. **Satu fungsi perhitungan** menggantikan 20 titik itu.
2. **Urutan operasi dibalik: kali dulu, bagi terakhir** — presisi antara menjadi tidak relevan.
3. **Rantai "hasil dipakai sebagai masukan" diputus.** `[data DBA]` Satu baris contoh membuktikan
   `Value` dihitung dari `GrossValue` yang **sudah dipotong 4 desimal**. Perhitungan berjalan
   **sekali jalan dari basis asli**, tidak membaca balik hasil antara.
4. `.ClaimSpreaded` adalah **detail/turunan** — boleh dihitung ulang dari sumber.
5. ⚠️ **Angka yang sudah terbit dibekukan** — nomor akseptasi sudah terbit, DLA sudah dicetak, atau
   data sudah dikirim ke Kasir. Itu **fakta**, tidak ikut dihitung ulang.

> #### ⭐ TAMBAHAN 2026-09-19 — share ceding, empat keputusan work owner
>
> `[data work owner 2026-09-19]` **Label `//` memang mematikan langkah.** Dipastikan dari layar
> Pega: work owner membuka `KomitePostAdjustment` dan memeriksa step **16** anak **1–4** —
> keempatnya **tampak dicoret**. Aturan baca ronde 1–3 karena itu **berdiri**.
>
> `[keputusan work owner 2026-09-19]` **Perkalian ganda share ceding tidak pernah berjalan.**
> #### TAMBAHAN 2026-09-19 (ronde 6) — rantai syarat `CopyOldataCurr_act` terbaca
>
> `[terverifikasi]` Gerbang langkah induk **7 · 8 · 9 · 10** `Activity/CopyOldataCurr_act.xml`
> dibaca apa adanya; keempatnya berflag prakondisi **`true`**, jadi gerbangnya berlaku. Enam baris
> berkode **5** berada di keempat langkah itu; baris ketujuh ada di
> `Activity/SetMOClaimTreaty_Act.xml` langkah 1.
>
> **Rantai yang berlaku**, diturunkan dengan kode 5 = *lewati sisa syarat, jalankan langkah*:
>
> ```
> langkah 7 (SpreadingRisk)      CurrencyIDOld != ""  AND  daftar tidak kosong
>                                AND ( Interest OR ClaimAmount )
> langkah 8 (EstimationList)     sama bentuknya dengan langkah 7
> langkah 9 (SpreadingClaim)     CurrencyIDOld != ""  AND  daftar tidak kosong
>                                AND ( Interest OR ClaimAmount OR Estimation )
> langkah 10 (SpreadingBreakQS)  sama bentuknya dengan langkah 9
> ```
>
> **Keempat langkah HIDUP.** Rantainya **dapat terpenuhi**. Bacaan lama yang menyatakannya
> mustahil berasal dari membaca baris berkode 5 sebagai **AND** — dengan aturan yang benar ia
> **OR**. Butir goyah ini **tidak dinyatakan tertutup**; keputusannya milik work owner.
>
> Bandingkan langkah **6** (`ListClaimAmount`) yang **tidak** punya baris berkode 5: rantainya
> murni AND — `CurrencyIDOld != "" AND daftar tidak kosong AND Interest`.

> `Activity/CountPersen_act.xml` blok **5** ber-remark, jadi perkalian kedua atas share ceding
> **tidak pernah dieksekusi di produksi**. ⭐ **Data lama AMAN** — tidak ada nilai tersimpan yang
> perlu diperbaiki. Blok itu **tidak dimigrasikan**.
>
> ⚠️ **Selisih yang dilaporkan, bukan dikejar supaya cocok.** Catatan lisan work owner menyebut
> yang ber-remark adalah langkah **5.1** dan bahwa langkah 5 berlabel kosong. Ekspor korpus
> menunjukkan **kebalikannya**: `pyStepsBlockName` berisi `//` pada langkah **5** (induk), dan
> **kosong** pada 5.1. ⛔ **Akibat keduanya sama** — blok 5 mati sepenuhnya kalau induknya mati,
> dan mati sebagian kalau hanya 5.1 yang mati; dalam kedua bacaan nilai yang dipakai hilir tetap
> berasal dari **blok 6**. `[terbuka]` mana yang benar; lihat `grilling-ronde-5.md` §B.
>
> `[keputusan work owner 2026-09-19]` **Lima titik perkalian share ceding di kelas Adjustment
> DIPERTAHANKAN** — `CountGrossAdjTreaty_Act` langkah **3 · 4 · 5.2** dan `CountValueADJTreaty_Act`
> langkah **12** (dua titik). ⛔ Perkaliannya **tidak dibuang** dan rumusnya **tidak diubah**. Yang
> diseragamkan hanyalah **ketelitiannya**, mengikuti butir 1–3 di atas: **hitung dengan ketelitian
> penuh, nol pemotongan di tengah jalan**. ⭐ Akibatnya angka berubah **hanya di belakang koma**,
> bukan sebesar satu faktor.

### 5. Mata uang per baris — invariant lama tidak berlaku

⚠️ `[terverifikasi]` Invariant "satu klaim satu mata uang" yang berlaku di Claim Life **TIDAK
berlaku** di Claim Prop. Tiga bukti bebas: `Activity/CurencyEstimation_Act.xml` menetapkan mata uang
**per baris**; baris penyelaras kurs hanya berlaku bila `CurrencyID` sama — **secara eksplisit
mengizinkan** baris bermata uang berbeda; dan `Activity/CountEstimation_Act.xml` **bercabang
eksplisit** antara jalur satu mata uang dan multi mata uang.

**Bentuk uang di Go: `(nilai, mata uang, kurs)` per baris**, dengan subtotal per mata uang sebagai
entitas tersendiri. Hal sama berlaku untuk Insured Interest — tiap objek punya kursnya sendiri.

### 6. ⚠️ Total share wajib tepat 100 % — penjaga dua sisi

`[keputusan work owner]` **Total share reinsurer wajib tepat 100 %.**

`[terverifikasi]` Pega hanya menjaga **sebelah**: `Activity/CountSpreading_Act.xml` step **6.2.3**
memeriksa `@greaterThan(Local.TotalPersen, 100)` — **satu-satunya** pemeriksaan total share di 329
berkas, dan ia **tidak pernah menangkap kurang dari 100**.

⚠️ **penyimpangan sadar:** di Go dipasang **penjaga dua sisi** saat input.

`[terverifikasi]` Karena total selalu 100 % **dan** presisi disimpan penuh, **alokasi sisa pembulatan
tidak diperlukan** — `[data DBA]` satu baris contoh membuktikan penjumlahan baris spreading **tepat
sampai digit terakhir**.

### 7. Wewenang — satu sumber, ditegakkan di lapisan layanan

`[keputusan work owner]`

```
tangga    Claim Admin → Claim Dept. Head → Operational Director → Technical Director
sumber    POOLDATA.EMAILKOMITE.DEGREE
dasar     tidak ketemu di roster → Claim Admin   (admin BUKAN anggota komite)
luar      Direktur Utama = wewenang TERPISAH, BUKAN tingkat kelima
```

`[terverifikasi]` `Claim Admin` di-set **tanpa syarat** lebih dulu di
`DataTransform/InsertChronology_DT.xml` langkah **1.1.4**, lalu **ditimpa** oleh tiga cabang bernama
(langkah 1.1.5–1.1.7). Bentuknya **bukan empat tingkat sejajar** melainkan **nilai dasar + tiga
tingkat komite**.

⚠️ **penyimpangan sadar:** **empat nama orang yang di-hardcode dibuang** — tiga di
`DataTransform/InsertChronology_DT.xml` dan satu di `RDBList/GetLimitDirekturUtama_SQL.xml`. Perilaku
tidak berubah (orang yang sama tetap mendapat tingkat yang sama), tetapi pergantian pemegang jabatan
cukup lewat baris tabel.

> ⚠️ **Catatan supersesi:** jawaban **Q5** (§21 ronde 2) berbunyi *"ikuti apa adanya —
> `GetLimitDirekturUtama_SQL` tetap mencari nama literal"*. Jawaban itu **digantikan** oleh §28e,
> yang diambil **setelah** diketahui roster dan jejak audit adalah tangga yang sama. **§28e berlaku.**

⚠️⚠️ `[terverifikasi]` **Pega tidak punya satu pun gerbang blocking di sisi server:** nol
`pyValidateActivity` di 12 dari 12 FlowAction; nol pemeriksaan peran di Activity, **418** ekspresi
gerbang di 36 dari 36 Section, **19** ReportDefinition, dan **11** Harness. Gerbang penyerahan komite
sepenuhnya *client-side* (`pyIsClientWhen=true`), dievaluasi dua kali dari sumber yang sama.

**Di Go dengan REST API ini dibangun dari nol.** ⚠️ Ini **bukan penyimpangan** — melainkan **mengisi
lubang yang memang tidak pernah ditutup**. Sejalan **ADR-0014**.

> #### ⚠️ RALAT 2026-09-17 — kunci pencocokan SUDAH final
>
> Teks lama: *"`[data DBA]` **Kunci pencocokan pelaku ke roster belum final** — kolom `EMAIL` yang
> sudah ada, atau kolom `USER_ID` baru. **Tidak menahan spec**; wajib dijawab sebelum tiket
> ditulis."*

✅ `[data DBA]` **Kunci pencocokan = `POOLDATA.EMAILKOMITE.OPERATOR_ID`**, dikonfirmasi sama dengan
user ID Pega. **Kolom `USER_ID` baru BATAL** — tidak perlu ditambahkan. **Kolom `EMAIL` tidak dipakai
sebagai kunci**; ia tetap hanya alamat pengiriman surat.

`[terverifikasi]` Kolomnya memang sudah diambil hari ini:
`ReportDefinition/FilterEmailKomiteWithLimit.xml` memilih `.OPERATOR_ID` bersama `.ID` · `.NAME` ·
`.EMAIL` · `.JABATAN` · `.DEGREE` · `.TYPE_KOMITE` · `.TYPE_BUSINESS`. **Nol perubahan skema** untuk
bagian ini.

#### ⚠️ Roster punya enam saklar per jenis aksi — nol dipakai

`[terverifikasi]` RD yang sama juga memilih **enam kolom `STS_*`** di luar yang menyaring:
`.STS_ADJ` · `.STS_ADJUSTER` · `.STS_REG` · `.STS_REJECT` · `.STS_SALVAGE` · `.STS_SURVEY`.
Bentuknya menunjukkan roster **dirancang** punya wewenang **per jenis aksi** — registrasi, adjuster,
salvage, survey, reject — bukan hanya per status klaim.

`[terverifikasi]` **Nol dari keenam dipakai sebagai penyaring oleh rule mana pun** (sensus 329
berkas). Penyaring yang sebenarnya hanya tiga, terbaca dari `pyFilterLogic` = `A AND C AND B`:
**`.LIMIT_BOTTOM`** (A) · **`.STS_KLAIM`** (C) · **`.STS_AKTIF`** (B, dibandingkan terhadap `"1"`).

⚠️ **Jebakan nama, bukan pemakaian.** `STS_REJECT` juga muncul di
`Activity/KonversiKlaim_Act.xml` dan `ConnectREST/KonversiKlaimNonLife.xml`, tetapi di sana ia
**`InService.STS_REJECT`** — properti muatan layanan konversi non-life, bukan kolom roster. Sensus
yang membaca nama telanjang akan salah menyimpulkan kolom roster itu terpakai.

`[terbuka]` **ringan** — Apakah wewenang per jenis aksi memang dimaksudkan berlaku, atau keenam kolom
itu warisan yang tidak terpakai? Pemilik: **Finance + work owner**. **Tidak memblokir** — yang ditiru
adalah perilaku yang berjalan, yaitu tiga penyaring.

### 8. Klasifikasi lini bisnis — hanya satu yang nyata

`[terverifikasi]` + `[data DBA]` **`InputData.CARI27` = `TreatyGroupID`**. Ia menentukan template
teks Insured Interest, di `Activity/SetValueToClaim_Act.xml` step **8** (penautan) dan step
**9, 10, 11, 12** (template):

| Kode | Kelas | Template |
| --- | --- | --- |
| `10007` | Property | *"Constrution Class : Occupation: Risk Category: Risk Location: Coverage :"* |
| `10016` | Motor Vehicle | *"Brand: Type/Year: Serial No/Engine No: Model:"* |
| `10009` | Marine Cargo | *"Project Name: Risk Location: Construction Period: Project Type:"* |
| 16 kode lain | Aneka | **template sama dengan Marine Cargo** |

⚠️ **penyimpangan TIDAK dilakukan — ikuti apa adanya, termasuk cacatnya:** Marine Cargo memakai
template **proyek konstruksi** yang sama dengan Aneka, dan typo **`Constrution`** ikut terbawa.
Bila kelak diperbaiki, catat sebagai penyimpangan sadar tersendiri.

⭐ `[keputusan work owner]` **2026-09-19** · **satu salinan, dua nasib**

`[terverifikasi]` Ke-61 rule klasifikasi lini **dipindahkan sekali sebagai satu himpunan**, dan
**hidup-matinya tidak ikut dipindahkan**. Sebuah rule bernilai benar hanya bila halaman yang
diujinya ada pada objek kerja modul yang memanggilnya. **Pada Claim Prop, 52 dari 61 menguji
`BusinessType` pada halaman `OfferFacIn`/`Quotation` yang tidak ada di sini — sehingga di modul
ini ke-52-nya tidak pernah bernilai benar.** Pada Claim Fac In halaman itu **ada**, dan rule yang
sama **hidup**. ⭐ Karena itu ia **bukan** rule yang "tidak dimigrasikan", melainkan rule yang
**mati di satu modul dan hidup di modul lain** — satu salinan, dua nasib.

⚠️ Klasifikasi lini di sistem baru tetap memakai **satu kunci `TreatyGroupID`** *(AC 71)*; ke-61
rule itu **tidak dipakai untuk klasifikasi**, dan keberadaannya semata agar satu salinan logika tetap
ada bagi modul yang memerlukannya.

`[data DBA]` Data klaim nyata memuat `BusinessName` dan `BusinessCode`, bukan `BusinessType`.

> ⛔ **RALAT 2026-09-19** — kalimat lamanya **dikutip utuh, tidak dihapus**:
>
> > *"⚠️ `[terverifikasi]` **52 dari 61 rule `When` adalah sisa impor model Fac In** — menguji
> > `BusinessType` pada halaman `OfferFacIn`/`Quotation` yang **tidak ada** pada objek kerja Claim
> > Prop. `[data DBA]` Data klaim nyata memuat `BusinessName` dan `BusinessCode`, bukan
> > `BusinessType`. **Dihidupkan pun semuanya bernilai salah.**"*
>
> ⚠️ **Angkanya tidak salah — yang salah pembingkaiannya.** *"Sisa impor"* menyarankan rule itu
> **milik modul lain yang tercecer di sini**; keputusan **A5-6 Claim Fac In** *(rule dipisahkan dari
> kehidupannya)* menunjukkan ia **satu rule bersama** yang **hidup di sana dan mati di sini**.

> #### ⚠️ RALAT 2026-09-17 — sensus dijalankan ulang atas 61 berkas
>
> Penyebut berubah dari 60 menjadi 61. Angka lama dan baru berdampingan:
>
> | Sensus di `When/` | Lama (60 berkas) | **Baru (61 berkas)** |
> | --- | ---: | ---: |
> | menguji `OfferFacIn` **atau** `Quotation` | 51 | **52** |
> | `BusinessType` | 44 | **45** |
> | `BusinessCode` | 10 | **10** |
> | `BusinessName` | 1 | **1** |
> | `TreatyGroupID` | 0 | **0** |
>
> `[terverifikasi]` Berkas baru **`When/IsCustomBonds.xml`** justru **menambah** barisan sisa itu,
> bukan mengurangi: kondisinya menguji `pyWorkPage.Quotation.BusinessType` terhadap `"CustomBond"` —
> halaman `Quotation` yang sama, yang tidak ada pada objek kerja Claim Prop. Identitasnya
> `ASM-FW-GCNMFW-WORK!ISCUSTOMBONDS`. **Kesimpulan bab ini tidak berubah, hanya menguat.**

**Spec tidak perlu memilih di antara dua mekanisme — hanya satu yang berlaku.**

⚠️ **`CARI27` `10004`** (keranjang Aneka) dan **`TreatyType` `10004`** (`QS (R/I)`) adalah **dua tabel
kode berbeda**. **Jangan digabung jadi satu enum.**

### 9. Tiga prefix klaim dan varian Syariah

`[terverifikasi]` Modul melayani **tiga prefix**: `CLM-` · `CLMP-` · `CLMNP-`, masing-masing dengan
varian ber-`S` untuk server Syariah (`CLMS-` · `CLMPS-` · `CLMNPS-`) yang **diperlakukan identik
dengan induknya** (`[keputusan work owner]`).

`[terverifikasi]` Konsekuensinya **tiga rule simpan OS akseptasi terpisah**:
`PEGA_JSON_OS_AKSEP_KLAIM` · `..._KLAIMTRT` (Treaty) · `..._KLAIMTNP` (Treaty Non Prop).

### 10. ⚠️ Jejak audit — empat cacat yang diperbaiki

`[terverifikasi]` Jejaknya adalah `ClaimData.SuggestList` dengan empat bidang: `CommentSuggest` ·
`DateSuggest` · `PICSuggest` · `IsCedingConfirm`. Dua penulis:
`DataTransform/InsertChronology_DT.xml` (27 pemanggil) dan `Activity/SethistoryKlaimTreaty.xml`.

| # | Cacat `[terverifikasi]` | Perbaikan |
| --- | --- | --- |
| 1 | `IsCedingConfirm` menyimpan **tingkat wewenang**, bukan boolean meski namanya berawalan `Is` | ⚠️ di Go menjadi **enum tertutup**, dan **nama pelaku dipisah ke kolomnya sendiri** |
| 2 | Tingkat terbawah dieja **dua bentuk** oleh dua penulis — `"Claim Admin"` versus `"Admin Claim <nama>"`. Korpus sendiri mengakalinya dengan `@contains(…,"Admin")` | ⚠️ **migrasi wajib memetakan keduanya** ke tingkat `Claim Admin` + nama ke kolom pelaku |
| 3 | `InsertChronology_DT` langkah **1** digerbangi `OperatorID.pyPosition != "IT Developer"`, dan **seluruh penulisan bersarang di bawahnya** → aksi peran itu **nol entri**. `SethistoryKlaimTreaty` **tidak** punya pengecualian ini | ⚠️ di Go **semua peran terekam tanpa kecuali** |
| 4 | Memo penutupan klaim **tidak pernah sampai** — `Activity/CloseClaimProp.xml` step **5** menulis ke `DataChronology.CARI12` sementara `InsertChronology_DT` langkah **1.1.1** membaca `DataChronology.CARI1`. ⚠️ **RALAT 2026-09-17:** sensus dijalankan ulang atas **329** berkas (dulu 328) dan **angkanya tidak berubah** — `DataChronology.CARI1` **34** kemunculan di 27 berkas, `DataChronology.CARI12` **satu** di **satu** berkas, yaitu langkah yang salah itu. ⚠️ Lingkup halaman `DataChronology.` **wajib** disebut: grep nama telanjang memberi `CARI1` 315 dan `CARI12` 57, karena `CARI12` juga slot parameter umum di `RDBList/` | ⚠️ diperbaiki; memo penutupan tercatat |

⚠️ `[data DBA]` **`SuggestList` tidak terurut waktu** — satu baris contoh memuat entri terakhir
bertanggal lebih lambat dari sebelas lainnya. **Di Go urutkan dari kolom tanggal, bukan urutan baris.**

### 11. Dokumen dan penomoran

`[keputusan work owner]` **DLA = *Definite Loss Advise***. Satu PDF **per reinsurer per baris
adjustment**, kategori lampiran `DLA`, dengan saudara `DLARetro`.

`[terverifikasi]` Ketiga dokumen (PLA, DLA, Acceptance Note) **disimpan**, bukan hanya di-stream:
alurnya token → unggah → simpan metadata di `T_STORAGE_IMAGE`, dengan URL berbatas waktu dan jalur
penyegaran (**ADR-0010**).

`[terverifikasi]` **Pencetakan bukan gerbang** — nol `pyValidateActivity`, dan
`Activity/CloseClaimProp.xml` tidak memeriksa nomor PLA maupun DLA. Yang ada hanya gerbang UI
**anti-cetak-ganda**. **Ditiru apa adanya.**

> #### ⚠️ RALAT 2026-09-17 — step 15.11.1 TETAP dimigrasikan
>
> Teks lama menyatakan **ketiga** step tidak ditulis, termasuk 15.11.1 yang hidup, dengan alasan ia
> satu-satunya cabang tersisa dari percabangan yang sudah tidak utuh. `[keputusan work owner]`
> **membalik bagian itu**: aturan induk berlaku atas yang di-remark saja.

⚠️ `[terverifikasi]` `Activity/PrintDLATreatyIn.xml` step **15.11.2** (salvage) dan **15.11.3**
(adjuster fee) **di-remark** → **tidak dimigrasikan**, tanpa pengecualian.

`[keputusan work owner]` **Step 15.11.1 diikuti apa adanya**, termasuk pola `TempData` dengan indeks
tetap **(1)**. **Konsekuensi yang diterima sadar dan dicatat:** DLA untuk **salvage** dan **adjuster
fee** menampilkan **nilai klaim**, persis seperti produksi hari ini.

### 12. Efek keluar — dari nol jaring pengaman menjadi outbox

`[terverifikasi]` `Activity/HitServiceToKasir_Act.xml` menjalankan 14 langkah **tanpa penanganan
kegagalan**: nol `Exit-Activity`, nol `Page-Set-Messages`, nol `Rollback`, dan **39 dari 39**
`pyOnException` kosong — termasuk pada kedua langkah `Connect-REST`. Tidak ada retry.

`[terverifikasi]` `POOLDATA.DIRECTTOKASIR_LOG` adalah **jejak, bukan antrean**: nol kolom
status/retry, satu INSERT dan satu SELECT di seluruh korpus, nol UPDATE, nol job pemroses ulang.

⚠️ **penyimpangan sadar:** di Go dibangun **outbox transaksional** (**ADR-0015**) — efek keluar
dicatat dalam transaksi yang sama dengan perubahan data, lalu dikirim dan **diulang** bila gagal.

`[terverifikasi]` Endpoint **tidak pernah literal** — keenam Connect-REST memakai satu setting yang
isinya ekspresi, dan URL sebenarnya berasal dari tabel `M_LINK_SERVICE` dikunci sepasang kategori
(**ADR-0013**). ⚠️ Pemisahan lingkungan **sepenuhnya bergantung isi tabel** — kelima *production
level* Pega bernilai identik.


> #### ⚠️ RALAT 2026-09-17 — TIDAK ADA BUG; label cacat dicabut
>
> Teks lama menyebut pembacaan `ReponseCode` sebagai **cacat kosmetik** yang membuat teks status
> gagal muncul, dan menggantungkan `[terbuka]` atas nama properti yang benar. **Keduanya dicabut.**

✅ `[data DBA]` **Respons Kasir memuat kedua ejaan sekaligus, dengan nilai identik:**
`"ReponseCode":"1"` **dan** `"ResponseCode":"1"`; `"ResponseMsg":"Success"` **dan**
`"ResponseMessage":"Success"`. Pega membaca ejaan tanpa `s` — dan karena Kasir mengirim keduanya,
teks **"Akseptasi Sudah Masuk ke Kasir"** memang muncul selama ini.

`[terverifikasi]` Sensus 329 berkas menajamkannya: `ReponseCode` ada di **1 berkas**
(`Activity/HitServiceToKasir_Act.xml`), `ResponseCode` di **nol berkas**. Pega memang tidak pernah
membaca ejaan yang benar — dan itu tidak pernah menjadi masalah, karena sisi luar mengirim dua-duanya.

### 13. Entitas pendukung

`[terverifikasi]` **Adjuster/Consultant: satu master, dua peran.** Master-nya **tidak punya kolom
tipe**; kedua rule pemilihnya klon dengan properti tujuan berbeda. Satu klaim = maksimal satu
adjuster + satu consultant. **Wajib terisi sebelum baris adjustment pertama** —
`Activity/AddAdjustment_Act.xml` step **9**, syarat kelengkapan di-AND dengan panjang daftar = 1,
sehingga **hanya berlaku pada baris pertama**. Ditiru apa adanya.

`[terverifikasi]` **Katastrofa punya entitas event sendiri** dengan nomor berprefiks `CTS-`, dan satu
nomor event **mengelompokkan banyak klaim**. Master disimpan lewat `Obj-Save` di
`Activity/SaveCatasrtope_Act.xml` step **3** — step SQL-nya (step 2) di-remark.

`[terverifikasi]` **Cause of Loss dua tingkat** dihubungkan `M_COL_ID`, ditambah tingkat ketiga
(relasi ke lini bisnis). ⚠️ Kedua procedure pembaruannya menerima **payload JSON** lewat parameter
bernama seolah kolom teks. **Wajib** sebelum simpan.

### 14b. ⭐ Titik yang sengaja diubah — 2026-09-19

Ronde 4 mencatat *"rumus share ceding diperbaiki"* sebagai **perubahan sadar ketiga**. Ronde 5
memeriksanya ulang dan **mencabut penggolongan itu**.

**Keputusannya: BUKAN perubahan sadar.** Alasannya satu baris — perkalian gandanya **tidak pernah
berjalan** (blok ber-remark), jadi tidak ada perilaku berjalan yang disimpangi; dan yang benar-benar
berubah menurut A3 hanyalah **ketelitian**, yang sudah tercakup keputusan uang di bab 4 butir 1–3.

⛔ Jumlah perubahan sadar karena itu **tetap dua**, tidak menjadi tiga.

> ⛔ **RALAT 2026-09-19** — kalimat lamanya **dikutip, tidak dihapus**: *"⛔ Jumlah perubahan
> sadar karena itu **tetap dua**, tidak menjadi tiga."* Ia benar sampai ronde 5. `[keputusan work
> owner]` 2026-09-19 menambahkan **tiga** perubahan sadar sekaligus — baris penyesuaian yang beku,
> larangan hapus klaim permanen, dan satu kolom Pega yang dibuang — sehingga jumlahnya menjadi
> **lima**. Rinciannya di **14c** di bawah.

#### 14c. ⭐ Baris penyesuaian yang membeku, dan kolom yang dibuang — 2026-09-19

`[keputusan work owner]` 2026-09-19 — ⚠️ **Baris penyesuaian membeku saat diserahkan ke komite.**
Sebuah baris dianggap **beku bila ada kasus komite yang menunjuknya** — **berjalan maupun selesai**.
⛔ **Tidak ada kolom penanda: bekunya diturunkan, bukan disimpan.** Baris beku **tidak dapat diubah,
selamanya**, dan **tidak dapat dihapus satu per satu**. **Penolakan komite tidak mencairkannya**;
perbaikan memakai **baris penyesuaian baru**.

⚠️ **Klaim yang pernah punya kasus komite tidak dapat dihapus — selamanya**, walau komitenya
**sudah selesai maupun menolak**, dan **jalur komitenya mana pun** — penyerahan baris penyesuaian
maupun **penutupan tanpa pembayaran** `[keputusan work owner]` 2026-09-19. Klaim yang **belum
pernah** ke komite **tetap dapat dihapus** dan **mengkaskade seperti biasa**.

⚠️ **Klaim lama hasil migrasi tidak dapat dihapus sama sekali**, apa pun riwayat komitenya
`[keputusan work owner — atas rekomendasi asisten]` 2026-09-19. Sebabnya: bukti *"pernah punya
kasus komite"* **diturunkan dari adanya kasus komite**, dan untuk data lama keberadaan kasus itu
bergantung pada apa yang berhasil dimigrasikan — sementara kolom yang dulu membuktikannya sudah
dibuang. ⛔ **Mudah dicabut**: bila yang dimaksud ternyata *semua klaim tanpa kecuali*, pencabutannya
menyentuh **AC 5**, separuh **AC 128**, dan butir uji kaskade di tiket 00. Lihat ralat di **AC 128**.

⚠️ **Kasus komite jalur penutupan TIDAK membekukan baris penyesuaian** — ia tidak menunjuk baris
mana pun. Yang dikuncinya adalah **klaimnya**. **AC 125 karena itu tidak berubah.**

⚠️ **Kolom `FLAG_ON_GOING_COMMITTEE` pada `T_GENERAL_CLAIM` dibuang.** Layar akseptasi
**menghitung sendiri** *"komite sedang berjalan"* dari kasus komitenya. `[terverifikasi]` Kolom itu
**ada di Pega** — ditulis `Activity/AddKomiteTreatyChild_ACT.xml` langkah **3**, juga
`KomitePostAdjustment` langkah **26.1** dan **27**, dan **dibaca layar akseptasi** — tetapi
**sengaja tidak dibawa** karena nilainya **dapat diturunkan**. ⚠️ **Nilai lamanya tidak
dimigrasikan** — lebih tegas lagi `[keputusan work owner]` 2026-09-19: **migrasi tidak membaca
kolom itu sama sekali**, dan nilainya **tidak dipakai sebagai bukti** apa pun.

⚠️ **Perilaku ini BARU.** `[terverifikasi]` Pega tidak memilikinya: **nol pemeriksaan, nol pesan
kesalahan, nol penanganan** untuk baris penyesuaian yang berubah atau hilang sesudah diserahkan.
**Karena itu ia dibangun, bukan dimigrasikan.** **Lingkup:** Claim Prop dan Claim Non Prop;
**Non Prop menyusul**, **Claim — Life tidak termasuk**.

⚠️ **Ketiganya penyimpangan sadar**, dan ketiganya membawa tanda `[penyimpangan sadar]` di
AC-nya: **AC 125 · 126 · 127** *(baris beku)* · **AC 128** *(hapus klaim)* · **AC 129** *(kolom
dibuang)*.

**Pembagian pekerjaannya:** tiket **11** — penyerahan **melahirkan** kasus komite, dan keberadaan
kasus itulah yang membekukan baris induknya; ⛔ **nol penanda dipasang** · tiket **08** —
**menegakkan** bekunya di **lapisan layanan** · tiket **00** — **membuang kolomnya** dan
**menyesuaikan jalur kaskade hapus**.

### 14. Risiko yang diterima sadar — ⚠️ kini DUA, bukan tiga### 14. Risiko yang diterima sadar — ⚠️ kini DUA, bukan tiga

⚠️ `[keputusan work owner]` Konsekuensi langsung aturan induk. **Diterima, bukan kelalaian migrasi.**

1. **Uang — spreading berjalan untuk semua payment type.**
   `[terverifikasi]` `Activity/CountSpreadingADJ_Act.xml` step **6**: saringan payment type
   `1|2|5` **tidak ditulis**, sehingga perhitungan spreading berjalan untuk **semua** payment type.
2. **Data di sistem luar — klaim yang ditolak komite tetap punya data di Arasapas.**
   `[terverifikasi]` `Activity/SaveOutstanding_Act.xml` step **39** dikirim **tanpa menunggu komite**,
   dan **nol rule pembatal di 329 berkas**.

> #### ⚠️ RALAT 2026-09-17 — risiko ketiga DIANGKAT, bukan dipertahankan
>
> Teks lama, butir 3: *"**Uang — satu baris adjustment yang melebihi estimasi tidak ditandai
> error.** `[terverifikasi]` `Activity/CountValueADJTreaty_Act.xml` step **16** punya dua baris
> syarat: baris pertama (`.AdjustmentValue > .TotalEstimasiValue`) **tidak menyaring** — kedua
> sisinya lanjut — sedangkan baris kedua (total) menyaring. **Pemeriksaan per baris ada di korpus
> tetapi tidak berpengaruh.** ⚠️ `[terbuka]` Niat aslinya tidak terbaca."*
>
> Niatnya kini terbaca, dan `[keputusan work owner]` **memperbaiki inkonsistensinya**. Lihat
> butir 14a di bawah dan **AC 49**.

#### 14a. ⚠️ Pesan dan penanda blokir dibuat konsisten

`[terverifikasi]` Yang sebenarnya terjadi di `Activity/CountValueADJTreaty_Act.xml` — dua step
berbeda membandingkan **dua besaran berbeda**, dan hanya satu yang memblokir:

| Step | Perannya | Baris syarat | Kode arah | Hidup? |
| --- | --- | --- | --- | --- |
| **12** | menghitung | `.AdjustmentValue` = propose × persen RNM × loss allocation × share ceding | — | — mata uang **asal** |
| **12** | menghitung | `.ValueAdjustment` = `Local.KursValue` × `.AdjustmentValue` | — | — **IDR** |
| **16** | SET `IsError = 2` — **MEMBLOKIR** | `.AdjustmentValue > .TotalEstimasiValue` | (kosong)/**2** | **MATI** |
| **16** | SET `IsError = 2` — **MEMBLOKIR** | `Local.SumProposeAdjustment > .TotalEstimasiValue` | (kosong)/**3** | hidup |
| **17** | `Property-Set-Messages` — **TAMPIL** | `.ValueAdjustment > .TotalEstimasiValue` | (kosong)/**3** | hidup |
| **17** | `Property-Set-Messages` — **TAMPIL** | `SumProposeAdjustment > … && .Type=="1"` | (kosong)/**3** | hidup |

`[terverifikasi]` `IsError` dibaca sebagai gerbang di **dua** tempat, keduanya `>1`:
`Activity/ProteksiInitialandDate_Act.xml` step **5** (aktif) dan `Section/AdjustmentDetail.xml`
lewat `pyDisabledWhen` (`.IsKomite==1 || pyWorkPage.IsError>1`). Step **19** justru
**membersihkannya** kembali ke `""` saat total masih di bawah estimasi.

⚠️ **Akibat di Pega:** bila **satu baris** melebihi estimasi dalam IDR, **pesan tampil tetapi
`IsError` tidak diset** — pengguna melihat teks error, tetapi tidak ada yang memblokir. Dan kedua
step membandingkan besaran berbeda: step 16 memakai **mata uang asal**, step 17 memakai **IDR**.

`[keputusan work owner]` **Alasan mengapa kelebihan per baris memang tidak boleh memblokir terbaca
di korpus:** baris salvage **wajib bernilai negatif** — `Activity/CountValueADJTreaty_Act.xml`
step **2** menyiapkan pesan *"Adjustment RNM for Salvage should be minus"*, ditegakkan step **18**
atas `.ProposeAdjustmentValue>0 && .Type=="3"`. Karena itu satu baris positif bisa melebihi estimasi
sementara **hasil bersihnya tidak**.

⚠️ **penyimpangan sadar — keputusan:** pesan dan penanda blokir memakai **satu besaran yang sama,
dalam IDR**, dibandingkan terhadap **estimasi terkini** (bukan potret — lihat AC 40). **Kelebihan
total memblokir penyimpanan; kelebihan per baris ditampilkan sebagai peringatan berlabel jelas** —
bukan error yang tampak memblokir padahal tidak. ⚠️ **Akar masalahnya lebih dalam dari ambang `>1`**
— satu properti mengerjakan tiga urusan sekaligus. Ditutup di **§15**.

### 15. ⚠️ `IsError` dipecah menjadi tiga hal terpisah

`[keputusan work owner]` 2026-09-17. **⚠️ penyimpangan sadar.**

`[terverifikasi]` Sensus 329 berkas: properti `IsError` disentuh **5 berkas** — **8 titik tulis di 3
rule** (enam di antaranya memasang nilai; dua sisanya mengembalikan ke `""`), dan **2 pembaca**.

**Penulis**

| Rule | Step | Nilai | Maksudnya |
| --- | --- | --- | --- |
| `Activity/CountValueADJTreaty_Act.xml` | **2** | `""` | reset di awal |
| `Activity/CountValueADJTreaty_Act.xml` | **14** | `2` | propose adjustment kosong |
| `Activity/CountValueADJTreaty_Act.xml` | **16** | `2` | melebihi estimasi (total) |
| `Activity/CountValueADJTreaty_Act.xml` | **19** | `""` | bersihkan — *enable* kirim ke komite |
| `Activity/CountValueADJTreaty_Act.xml` | **20** | `3` | *disable* kirim ke komite |
| `Activity/ProteksiInitialandDate_Act.xml` | **4.1** | `2` | — |
| `Activity/CheckDateDOL_Act.xml` | **25** | `1` | ⚠️ **yatim — tidak pernah lolos `>1`** |
| `Activity/CheckDateDOL_Act.xml` | **26** | `0` | tidak ada duplikat |

**Pembaca — keduanya menguji `> 1`**

| Tempat | Ekspresi |
| --- | --- |
| `Activity/ProteksiInitialandDate_Act.xml` step **5** | `pyWorkPage.IsError>1` |
| `Section/AdjustmentDetail.xml` `pyDisabledWhen` | `.IsKomite==1 \|\| pyWorkPage.IsError>1` |

⚠️ **Nilai `0` · `1` · `2` · `3` · `""` tidak punya arti yang tertulis di mana pun** — bukan di nama
properti, bukan di keterangan step, bukan di rule terpisah. Artinya hanya dapat disimpulkan dari
tempat ia ditulis, dan dua di antaranya bertabrakan: `2` berarti "salah", `3` berarti "belum boleh
kirim ke komite" — dua urusan berbeda pada satu pita angka yang sama.

**Di Go dipecah menjadi tiga hal bernama sendiri:**

| Hal | Bentuk |
| --- | --- |
| kesalahan yang memblokir | ya / tidak |
| peringatan yang ditampilkan | daftar pesan |
| boleh kirim ke komite | ya / tidak |

Duplikat Date of Loss masuk ke **kesalahan yang memblokir dan peringatan sekaligus** (§15a, AC 102).

⚠️ **Menaikkan ambang menjadi `>= 1` BUKAN perbaikan yang benar** — nilai `3` adalah saklar
kirim-ke-komite, dan dengan ambang itu ia akan ikut memblokir hal yang bukan urusannya. Pemisahan
menjadi tiga hal adalah satu-satunya perbaikan yang tidak menukar satu cacat dengan cacat lain.

#### 15a. ⚠️ Duplikat Date of Loss MEMBLOKIR penyimpanan

`[keputusan work owner]` **Perbaikan, bukan paritas.** `[terverifikasi]` Perilaku Pega hari ini di
`Activity/CheckDateDOL_Act.xml`:

| Step | Aktif? | Perbuatannya |
| --- | --- | --- |
| **20.1** | aktif | `Local.DOL==.START_DATE && Local.CLMNo!=.BRANCH_NAME` → `Local.Flagerr = 1`, dan simpan INSKEY klaim lama ke `Data.CARI42` |
| **21–22** | aktif | ambil data klaim lama lewat `ReportDefinition/RejectedClaim_RD` dengan `Param.INSKEY` dari `Data.CARI42` |
| **23.1** | aktif | klaim lama **ada** di tabel reject (`.INSKEY==Data.CARI42`) → `Local.Flagerr = 0` — **pengecualian** |
| **24** | aktif | `Local.Flagerr==1` → `Property-Set-Messages` — **peringatan tampil** |
| **25** | aktif | `Local.Flagerr==1` → `.IsError = 1` |
| **26** | aktif | `Local.Flagerr!=1` → `.IsError = 0` |

⚠️ **Nilai `1` tidak pernah memblokir**, karena kedua pembaca menguji `> 1`. Keterangan step 25
berbunyi *"make disable submit if there is similar date"* — **niatnya memblokir**, tetapi nilai yang
dipasang tidak melewati ambangnya. Jadi **peringatan tampil dan penyimpanan tetap lolos**.

⚠️ **penyimpangan sadar — keputusan:** klaim dengan Date of Loss sama pada polis yang sama
**memblokir penyimpanan** dan **menampilkan peringatan**. **Pengecualian tetap berlaku**: bila klaim
lama sudah ditolak, tidak ada blokir dan tidak ada peringatan — pengecualian itu sudah berjalan di
hulu (step **23.1**), jadi alur klaim pengganti tidak terganggu.

#### 15b. ⚠️ Temuan tambahan — satu nomor polis dikecualikan lewat hardcode

> ### ⚠️ RALAT 2026-09-17 — lokasinya **step 20**, bukan step 20.1
>
> Teks lama menempatkan baris syarat polis pada **step 20.1**. **Itu salah.** Verifikasi ulang atas
> struktur bersarang `Activity/CheckDateDOL_Act.xml` membuktikan ia milik **step 20**.
>
> **Sebabnya jebakan serialisasi:** `pyStepsPreCondParams` milik sebuah step **induk** ditulis
> **sesudah** daftar `pySteps` anaknya ditutup. Dibaca berurutan tanpa memperhatikan pembungkus,
> syarat milik induk tampak menempel pada anak terakhir. Ini **jebakan ketiga** sejenis dalam proyek
> ini, setelah step 16-versus-17 dan pemisahan `CARI*` berlingkup halaman.
>
> ⚠️ **Arahnya TIDAK diralat.** Teks lama sudah berbunyi *"dibebaskan dari pemeriksaan duplikat"*,
> dan itu benar. Yang diralat hanya lokasinya.

`[terverifikasi]` Struktur sebenarnya, dibaca dari pembungkus `<pySteps>`:

| Step | Jenis | Flag `pyStepsPreCondition` | Baris syarat | Kode arah |
| --- | --- | --- | --- | --- |
| **20** | loop atas `TempHistoryClaim.pxResults` | **`false`** | `pyWorkPage.OfferFacIn.PolicyData.PolicyNo == " RNM-F01.01.2018.00040"` | BENAR=**3** · SALAH=(kosong) |
| **20.1** | `Property-Set` — *"set error dan lempar inskey"* | **`true`** | `Local.DOL==.START_DATE && Local.CLMNo!=.BRANCH_NAME` | BENAR=2 · SALAH=3 |

**Step 20.1 hanya punya satu baris syarat, arah normal.** Seluruh carve-out ada di step 20.

⚠️ **Arahnya MEMBEBASKAN, bukan membatasi.** Dengan aturan *sisi kosong = continue when*:
polisnya **sama** dengan nomor itu → **lewati**; **bukan** → **jalan**. Jadi niat aslinya
**mengecualikan satu polis** dari pemeriksaan duplikat.

⚠️ Literalnya diawali **spasi**. Sensus 329 berkas: nomor itu muncul di **satu berkas**, yaitu step
ini sendiri. `[terverifikasi]` `OfferFacIn.PolicyData.PolicyNo` **tidak pernah ditulis** di 329 berkas
— ia hanya dibaca (7 berkas).

**Kesimpulan akhir tidak berubah, tetapi alasannya berubah.** Flag step 20 = **`false`**, sehingga
menurut aturan induk **WHEN tidak ditulis dan loop tetap ditulis** — pengecualian itu **tidak berlaku
hari ini dan tidak dimigrasikan**.

⚠️ `[keputusan work owner]` **Alasan baru yang wajib dicatat:** ini **bekas pengecualian bisnis yang
sengaja dipasang lalu dimatikan**, **bukan** sisa saringan uji coba. Bedanya penting — sisa uji coba
dibuang tanpa bertanya; pengecualian bisnis yang dimatikan **layak ditanyakan kembali**.

`[terbuka]` **ringan tetap berdiri** — apakah pembebasan satu polis ini masih dikehendaki, dan bila
ya, apakah ia menjadi baris data seperti keputusan **AC 55**? Pemilik: **work owner**.
**Tidak memblokir** — default yang ditulis adalah **tidak dimigrasikan**.

### 16. ⚠️ Nama alias SQL yang berbohong — `CariHistoryClaim_SQL`

`[keputusan work owner]` **⚠️ penyimpangan sadar. Di Go, nama dibuat sesuai isinya.**

`[terverifikasi]` `RDBList/CariHistoryClaim_SQL.xml` membaca tabel `json_klaim` dan memberi **enam
alias**, **lima di antaranya tidak ada hubungannya dengan isinya**:

| Sumber | Alias | Isinya sebenarnya |
| --- | --- | --- |
| `a.DATA_JSON.DateOfLoss` | `"START_DATE"` | ⚠️ tanggal kejadian, **bukan** tanggal mulai |
| `NOPOLIS` | `"POLICY_NO"` | ✅ benar |
| `IDPEGA` | `"BRANCH_CODE"` | ⚠️ id Pega, **bukan** kode cabang |
| `a.DATA_JSON.ClaimNo` | `"BRANCH_NAME"` | ⚠️ nomor klaim, **bukan** nama cabang |
| `a.DATA_JSON.CauseOfLoss` | `"BUSINESS_CODE"` | ⚠️ penyebab kerugian, **bukan** kode bisnis |
| `a.DATA_JSON.ClaimEstimate` | `"TSI"` | ⚠️ estimasi klaim, **bukan** TSI |

⚠️ **Akibatnya terbaca langsung di pemanggilnya.** `Activity/CheckDateDOL_Act.xml` step **20.1**
berbunyi `Local.CLMNo != .BRANCH_NAME` — membandingkan nomor klaim dengan kolom bernama *"nama
cabang"* — dan menyimpan `.BRANCH_CODE` sebagai INSKEY klaim lama. **Bila alias itu terbawa, kodenya
tidak terbaca oleh siapa pun yang tidak ikut grilling ini.**

⚠️ `[terverifikasi]` **Temuan yang menempel, relevan bagi uang.** SQL yang sama membungkus estimasi
dengan `REPLACE(a.DATA_JSON.ClaimEstimate, ',', '.')` — angka uang itu **tersimpan sebagai teks
dengan koma sebagai pemisah desimal**, lalu ditambal saat dibaca. Ini menguatkan **§4**: di Go nilai
uang disimpan sebagai **desimal**, bukan teks, sehingga tidak ada tambalan pemisah desimal di jalur
baca mana pun.

⚠️ Penyaringnya hanya `where nopolis = {TempPolis.PolicyNo}`, diurutkan `TGL_INPUT ASC` — **tanpa
saringan tanggal**. Seluruh riwayat klaim satu polis ditarik, lalu disaring di dalam aktivitas.

#### 16a. ⚠️ Dua alias berbohong lagi, dan satu tabel dieja dua cara

`[terverifikasi]` Ditemukan saat menutup celah cakupan registrasi:

| Rule | Yang ditulis | Masalahnya |
| --- | --- | --- |
| `RDBList/SetPolicyTreatyProp.xml` | `TREATYGROUP AS "BusinessName"` | ⚠️ kelompok treaty dinamai **nama bisnis** |
| `RDBList/SetPolicyTreatyProp.xml` | `'0' AS "Prodke"` | ⚠️ literal yang dipalsukan jadi kolom |
| `RDBList/SetPolicyTreatyProp.xml` | `FROM TREATYINPRODUCTION` | ⚠️ **tanpa prefiks schema** |
| `RDBList/CekPolicyNumber_SQL.xml` | `from pooldata.TREATYINPRODUCTION` | ✅ **tabel yang sama, diprefiks** |

⚠️ **Tabel yang sama dieja dua cara di dua rule.** Ini bukti kedua yang berdiri sendiri untuk **§2**:
resolusi nama bergantung schema bawaan koneksi, dan seorang pembaca tidak dapat tahu dari SQL-nya
sendiri ke mana `TREATYINPRODUCTION` menunjuk.

⚠️ `[terverifikasi]` **Tambalan koma-desimal muncul untuk kedua kalinya**, kali ini di jalur premi:
`Activity/CekPremiLunas_Act.xml` step **5** menjalankan `@replaceAll(…,",",".")` atas hasil jumlah
premi, persis seperti `REPLACE(…,',','.')` di `CariHistoryClaim_SQL`. **Dua jalur berbeda, satu
penyakit yang sama** — angka uang tersimpan sebagai teks. Ditutup oleh **AC 107**.

#### 16b. ⚠️ Satu tambahan — parameter yang salah dokumen

`[data DBA]` Ditemukan saat menutup butir `[terbuka]` ke-6:

| Yang ditulis | Yang sebenarnya diisi | Di mana |
| --- | --- | --- |
| parameter `STS_PLA` | kolom **`STS_DLA`** | `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM` |

⚠️ **PLA dan DLA dua dokumen berbeda** — PLA adalah pemberitahuan awal kerugian, DLA adalah
*Definite Loss Advise* per reinsurer (**AC 84**). Nama parameternya menyebut dokumen yang salah.
`[keputusan work owner]` **Namanya dibiarkan apa adanya** (§1b) — dicatat di sini supaya pembaca
berikutnya tidak menyangka ada jalur PLA yang hilang.

### 17. ⚠️ Batas `COMMIT` — procedure yang sama, `COMMIT` dijalankan aplikasi

`[keputusan work owner]` **2026-09-18.**

Go memanggil **stored procedure yang sama** — `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM` ·
`POOLDATA.PEGA_JSON_KLAIM_PNC` · `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP` — dan menjalankan `COMMIT`
**sendiri sesudah panggilan itu**. **Logika penulisan tidak direplikasi di aplikasi.**

`[terverifikasi]` **Ini sejalan dengan AC 7, bukan melawannya.** AC 7 menuntut **nol `COMMIT` di
dalam teks SQL aplikasi**. Di Pega `COMMIT` tertanam di dalam blok anonim rule — `BEGIN … ; COMMIT;
END;` pada `RDBList/InsertClaimPNC.xml`. Di Go `COMMIT` **dikeluarkan oleh aplikasi**, bukan ditanam
di dalam teks SQL. **Objek yang dipanggil tetap sama.**

#### 17a. ⚠️ Akibat yang mengubah cara memanggil

`[data DBA]` Ketiga procedure memasang **`EXCEPTION WHEN OTHERS THEN ROLLBACK` di dalam dirinya
sendiri**. ⚠️ **`ROLLBACK` telanjang di Oracle membatalkan SELURUH transaksi**, bukan hanya pekerjaan
procedure itu — termasuk apa pun yang aplikasi tulis **sebelumnya** dalam transaksi yang sama. Dan
`ROLLBACK TO SAVEPOINT` **tidak menolong**, karena yang dipasang bukan itu.

**Akibatnya, urutan pemanggilan menjadi wajib:**

1. Panggilan procedure ini adalah **langkah terakhir sebelum `COMMIT`**.
2. **Tidak boleh ada pekerjaan penting yang masih menggantung** sebelumnya dalam transaksi yang sama.
3. Aplikasi **memeriksa `StsSimpan`** sesudah panggilan — **1 berhasil, 0 gagal**.
4. Pada **0**, aplikasi **tidak boleh melaporkan berhasil**: pada saat itu transaksinya **sudah
   dibatalkan oleh procedure**, sehingga `COMMIT` sesudahnya **tidak menyimpan apa pun**.

⚠️ **AC 6 belum menampung urutan ini.** Bunyinya hanya *"Satu aksi pengguna menghasilkan satu
transaksi. Test yang menemukan keadaan separuh tersimpan gagal."* — ia menuntut hasil, bukan urutan.
Usulan AC tambahan **tidak dituliskan di sini**; ia dilaporkan ke work owner untuk diputuskan.

---

### §18 — Lima keputusan lintas modul, `[keputusan work owner]` 2026-09-19

*Modul ini dinyatakan tuntas **sebelum** pernah diadu dengan ke-15 ADR. Pengaduan itu dilakukan
2026-09-19 dan menghasilkan lima keputusan. Buktinya di `utang-lintas-modul.md`; yang dicatat di
sini adalah **alasannya**, bukan sekadar isinya.*

**§18a — Uang yang menyeberang tetap desimal.** `[terverifikasi]` **12 parameter berbau uang
bertipe `Double`** di 5 berkas, dan yang terberat ada di **`AddKomiteTreatyChild_ACT`** — rule yang
**melahirkan kasus komite**. ⭐ **Buktinya lebih kuat daripada di Claim Fac In:** di sana kedua
sumber tipe berselisih sehingga vonisnya harus bersandar pada tambalan koma-ke-titik di SQL; **di
sini keduanya sepakat**. ⛔ Yang **tidak** disimpulkan: bahwa presisi benar-benar hilang — yang
dibaca **deklarasi tipe**, bukan jalannya nilai. Yang tegas: **ADR-0003 melarang representasi itu,
dan representasi itu ada di sini.** *(AC 130 · 131)*

**§18b — Migrasi penuh, tanpa koeksistensi.** ADR-0009 **belum pernah disebut** di modul ini, dan
spec **tidak pernah memutuskan** penuh lawan koeksistensi — yang ada hanya *"tidak dimigrasikan"*
tentang **rule**, bukan **data**. ⭐ Alasan penolakan koeksistensi di sini **lebih kuat daripada di
Claim Life**: aturan beku modul ini **permanen** *(AC 126 · 127 · 128)*, dan sistem lama **tidak
mengenalnya**. *(AC 132)*

**§18c — Wewenang lewat peran, bukan teks jabatan.** `[terverifikasi]` Sapuan `pyPosition` memberi
**115 kemunculan / 28 berkas**, dan itu **menipu**: **38** bernilai `Top` dan **18** bernilai
`AFTER` — keduanya **tata letak layar** di ReportDefinition dan Section. ⭐ **Hanya dua** yang
wewenang sejati, dan bunyinya perbandingan terhadap **teks jabatan**. ⛔ **Nama peran gaya Life
— `ReasLifeAdmin` dan kerabatnya — NOL**, sehingga model tiga peran ADR-0002 memang tidak ada di
sini. *(AC 133)*

**§18d — Pengenal berkas: ambil versi yang sudah ditambal.** `[terverifikasi]` Dari **162** rule
beridentitas sama dengan Claim Fac In, **158 identik** dan **4 beda versi**. Rule pembuat pengenal
berkas adalah satu-satunya di antara keempatnya yang **isinya benar-benar berbeda**: versi baru
menaikkan ketelitian waktu **dari milidetik ke nanodetik** dan menambahkan **nilai unik sejagat**.
⭐ **Bentuknya persis sebuah tambalan tabrakan** — dan ia **membenarkan temuan modul ini sendiri**
di ronde 2. ✅ Kesimpulan lama **bukan gugur, melainkan dikuatkan**, dan kini punya jalan keluar
yang **sudah terbukti dipakai**. *(AC 134)*

**§18e — Jenis berkas: paritas.** Rule penentu jenis berkas versi baru memuat **48** baris lawan
**42**; ⛔ **tidak diambil**. Menambah jenis yang diterima **mengubah perilaku**, dan tidak ada yang
memintanya. *(AC 135)*

**§18f — Kurs standar: kecemasan yang GUGUR.** ⭐ Rule kurs standar di modul ini **4 tahun 3 bulan
lebih tua** daripada salinan Claim Fac In, dan itu sempat dicurigai membatalkan kesimpulan kurs.
`[terverifikasi]` **Isinya identik** — **12 medan berbeda dari 74**, dan **nol** di antaranya
logika; perintah SQL-nya **sama kata demi kata** di ketiga modul. ⛔ **Tidak ada AC yang lahir dari
butir ini**, karena tidak ada yang berubah. Dicatat supaya tidak dicurigai ulang.

---

## Testing Decisions

### Seam

**Satu seam: API HTTP.** Ini seam yang sudah dipakai konteks lain di proyek ini dan dipertahankan —
makin sedikit seam makin baik.

- **Oracle tidak pernah di-mock.** Test berjalan terhadap skema sungguhan. Ini penting karena seluruh
  keputusan presisi, transaksi, dan penomoran hanya bermakna terhadap database nyata.
- **Yang di-fake hanya klien luar:** penyimpanan berkas, Kasir, konversi ke sistem inti, dan email.
- Test memeriksa **perilaku dari luar** — apa yang dikembalikan API dan apa yang berubah di database
  — **bukan** detail implementasi.

### Yang wajib diuji

| Area | Yang dibuktikan |
| --- | --- |
| Skema & migrasi | pohon klaim tersimpan utuh; hapus klaim mengkaskade sampai tingkat terdalam |
| Transaksi | satu aksi = satu transaksi; gagal di tengah tidak meninggalkan keadaan separuh |
| **Penomoran** | **nomor tidak terbakar** — gagal sesudah pengambilan nomor membuat nomor tidak terpakai |
| Uang | presisi penuh tersimpan; hasil tidak berubah oleh urutan pemanggilan; nol *floating point* |
| Total share | ditolak pada `< 100 %` **dan** `> 100 %` |
| Multi mata uang | satu klaim dengan dua mata uang tersimpan dan terhitung benar |
| Wewenang | penyerahan komite **ditolak di lapisan layanan** meski kontrol UI ditampilkan |
| Jejak audit | semua peran terekam; memo penutupan tercatat; urutan dari kolom tanggal |
| Efek keluar | kegagalan masuk outbox dan terkirim ulang; tidak hilang diam-diam |
| Risiko sadar | test yang **menuntut** perilaku "diperbaiki" pada butir 14 justru **gagal** |

⚠️ **Prior art:** pola test Claim Life dan Treaty Contract Out — satu seam HTTP, Oracle nyata, klien
luar dipalsukan.

---

## Acceptance Criteria

### Skema dan transaksi

1. ⚠️ `[terverifikasi]` Skema relasional Claim Prop **dirancang baru** — test yang mengandalkan tabel
   klaim warisan **gagal**, karena tabel itu tidak pernah ada.
2. `[terverifikasi]` Pohon klaim tersimpan relasional: klaim → objek pertanggungan · estimasi · baris adjustment ·
   spreading · loss allocation, masing-masing baris dapat di-`SELECT` tanpa membongkar JSON.
3. ⚠️ Enam properti spreading menjadi **tabel berbeda** meski di Pega berbagi satu class.
   *(Implementation Decision 1)*
4. `[terverifikasi]` Jejak audit menjadi **tabel milik Claim Prop sendiri**, bukan menumpang struktur modul lain.
5. `[keputusan work owner]` Menghapus klaim **mengkaskade sampai tingkat terdalam**; test wajib memeriksa cicit.
   ⚠️ **Satu pengecualian:** klaim yang **pernah** punya kasus komite **tidak dapat dihapus sama
   sekali** — permintaannya **ditolak** dan **kaskadenya tidak pernah dijalankan**, walau komitenya
   **sudah selesai maupun menolak**.
   > ⛔ **RALAT 2026-09-19** — kalimat lama **tidak mengenal klaim yang pernah masuk komite**;
   > ia **dikutip utuh, tidak dihapus**: *"5. `[keputusan work owner]` Menghapus klaim
   > **mengkaskade sampai tingkat terdalam**; test wajib memeriksa cicit."* Pengecualiannya
   > ditambahkan `[keputusan work owner]` 2026-09-19 — lihat **AC 128** dan **14c**.
   > ⛔ Nomornya **tetap 5**, dan baris **Menutup:** tiket 00 **tidak berubah karenanya**.
6. `[keputusan work owner]` Satu aksi pengguna menghasilkan **satu transaksi**. Test yang menemukan keadaan separuh tersimpan
   **gagal**.
7. `[penyimpangan sadar]` ⚠️ `[terverifikasi]` Nol `COMMIT` di dalam SQL aplikasi — kelima belas rule yang `COMMIT` sendiri
   **tidak direplikasi**. *(§8b ronde 2)* **⚠️ penyimpangan sadar.**
   > ⚠️ **RALAT 2026-09-18** — frasa penggolongan ditambahkan; sebelumnya AC ini ber-⚠️ tanpa
   > penggolongan padahal ia jelas menyimpang: Pega menutup transaksinya sendiri di dalam SQL,
   > Go memindahkan batas transaksi ke lapisan aplikasi — sejajar dengan **AC 13** yang sudah
   > tergolong penyimpangan sadar atas keputusan yang sama.
8. `[penyimpangan sadar]` ⚠️ Setiap objek Oracle diprefiks **schema eksplisit**. Empat schema dikenal: `pooldata`,
   `arasapas`, `reinsurance`, `gl`. **⚠️ penyimpangan sadar.**
9. `[keputusan work owner]` Migrasi memindahkan data klaim lama dari bentuk JSON ke bentuk relasional.
10. ⚠️ `[terverifikasi]` Parser JSON warisan **membuang seluruh kunci berawalan `px`, `py`, `pz`**,
    dan **tidak mengandalkan kunci tertentu selalu hadir**.
11. `[penyimpangan sadar]` ⚠️ `[terverifikasi]` Status tampilan (`pyExpanded`) **tidak ikut disimpan** — kebocoran lapisan UI
    ke penyimpanan **tidak ditiru**. **⚠️ penyimpangan sadar.**
12. `[data DBA]` Pembaca JSON warisan menerima **dua format tanggal berdampingan** dalam satu dokumen.

### Penomoran

13. `[penyimpangan sadar]` ⚠️ **Pengambilan nomor dan penyimpanan hasilnya berada dalam satu transaksi.** Bila penyimpanan
    gagal, **nomor tidak terpakai**. **⚠️ penyimpangan sadar.**
14. `[keputusan work owner]` Test: gagalkan langkah sesudah pengambilan nomor → nomor berikutnya **berurutan rapat**, tidak
    melompat.
15. `[keputusan work owner]` Logika pembentukan nomor **tetap di stored procedure**, tidak direplikasi di aplikasi
    (**ADR-0006**).
16. ⚠️ `[terverifikasi]` Nomor **temp** tidak dapat dipromosikan menjadi nomor final — ia kekurangan
    `KODE_BIS` dan `BULAN`.
17. `[terverifikasi]` Nomor klaim, PLA, dan DLA memakai **pola tanda tangan yang sama** — lima parameter masuk, dua
    keluar.

### Uang

18. ⚠️ Nilai uang bertipe **desimal presisi arbitrer**; **nol** *binary floating point* di lapisan
    mana pun maupun di kontrak API (**ADR-0003**).
19. `[keputusan work owner]` Nilai **disimpan penuh tanpa pembulatan**; pembulatan **hanya untuk
    tampilan**.
20. `[penyimpangan sadar]` ⚠️ Perhitungan spreading dilakukan **satu fungsi**, menggantikan 20 titik dengan 4 perlakuan
    pembulatan berbeda. **⚠️ penyimpangan sadar.**
21. ⚠️ Urutan operasi **kali dulu, bagi terakhir**. Test: hasil **tidak berubah** oleh urutan
    pemanggilan.
22. `[penyimpangan sadar]` ⚠️ Nilai turunan dihitung **sekali jalan dari basis asli** — test yang menemukan perhitungan
    membaca balik hasil antara **gagal**. **⚠️ penyimpangan sadar.**
23. ⚠️ Angka yang **sudah terbit** — nomor akseptasi terbit, DLA dicetak, data terkirim ke Kasir —
    **dibekukan** dan tidak ikut dihitung ulang.
24. `[terverifikasi]` `.ClaimSpreaded` diperlakukan sebagai **turunan**, boleh dihitung ulang dari
    sumber.

### Mata uang

25. ⚠️ `[terverifikasi]` Bentuk uang adalah **`(nilai, mata uang, kurs)` per baris**. Test yang
    menuntut satu klaim satu mata uang **gagal** — invariant itu **tidak berlaku** di Claim Prop.
26. `[terverifikasi]` Subtotal per mata uang tersimpan sebagai entitas tersendiri.
27. `[terverifikasi]` Tiap objek pertanggungan punya **mata uang dan kursnya sendiri**.
28. `[terverifikasi]` Satu klaim dengan **dua mata uang** tersimpan dan terhitung benar dari ujung ke ujung.

### Loss allocation dan spreading

29. `[terverifikasi]` Loss allocation membagi nilai kerugian ke **treaty**, disegmentasi **per mata
    uang** — tanpa dimensi tahun maupun coverage.
30. `[penyimpangan sadar]` ⚠️ `[keputusan work owner]` **Total share wajib tepat 100 %.** Ditolak pada `< 100 %` **dan**
    `> 100 %`. **⚠️ penyimpangan sadar** — Pega hanya menjaga sebelah.
31. `[terverifikasi]` **Alokasi sisa pembulatan tidak diperlukan** — karena total 100 % dan presisi
    penuh, penjumlahan baris spreading **tepat sampai digit terakhir**.
32. `[terverifikasi]` Spreading di tingkat klaim terpisah dari spreading pada baris adjustment.
33. `[terverifikasi]` Loss allocation disalin sebagai **snapshot** ke baris adjustment saat baris dibuat.

### Estimasi

34. `[terverifikasi]` Baris estimasi dibentuk **per treaty**, bersumber dari hasil loss allocation.
35. `[terverifikasi]` Tanggal estimasi wajib **antara Date of Loss dan hari ini**, inklusif di kedua
    ujung, zona `Asia/Jakarta`.
36. `[terverifikasi]` Total estimasi melebihi TSI → peringatan.
37. `[terverifikasi]` Total estimasi melebihi plafon cash call → peringatan.
38. `[terverifikasi]` Menghapus baris estimasi menyesuaikan seluruh subtotal.
39. `[penyimpangan sadar]` ⚠️ `[terverifikasi]` Dua nama kolom Pega **berbohong** dan **tidak dibawa apa adanya** ke nama
    kolom baru: satu kolom bernama persen berisi **nilai uang**, satu kolom bernama jenis kerugian
    berisi **identitas treaty**. **⚠️ penyimpangan sadar** — penamaan diluruskan.
40. `[penyimpangan sadar]` ⚠️ `[keputusan work owner]` **Pagar nilai mengikat estimasi TERKINI, bukan potret.** Di Pega
    nilainya potret: `Activity/AddAdjustment_Act.xml` step **5** menulis `.TotalEstimasiValue` dari
    `.ClaimData.TotalEstimasiIDR` saat baris dibuat, dan **nol rule lain memperbaruinya** — sensus
    329 berkas menemukan properti itu di **2 berkas**: satu penulis (step 5) dan satu pembaca
    (`Activity/CountValueADJTreaty_Act.xml`, hanya di dalam ekspresi syarat). Potret **tetap
    disimpan** untuk jejak audit, tetapi **tidak** dipakai sebagai pembanding.
    **⚠️ penyimpangan sadar.**
    > ⚠️ **RALAT 2026-09-17** — teks lama: *"`[terverifikasi]` Nilai estimasi pada baris adjustment
    > **tidak** berasal dari daftar estimasi — `[terbuka]` sumber pengisinya tidak ditemukan di
    > korpus."* Sumbernya kini ditemukan dan pertanyaannya dijawab.

### Insured Interest

41. ⚠️ `[terverifikasi]` "Interest" di modul ini adalah **Insured Interest (objek pertanggungan
    dengan TSI)**, **bukan bunga finansial**. Nol rate, nol jumlah hari, nol basis 360/365 di seluruh
    rule perhitungannya.
42. `[terverifikasi]` Daftar objek pertanggungan wajib terisi sebelum simpan.
43. `[terverifikasi]` Total TSI per mata uang dan total dalam IDR tersedia.

### Baris adjustment

44. `[terverifikasi]` Unit keputusan adalah **baris `AdjustmentList`**, masing-masing dengan statusnya
    sendiri (**ADR-0011**).
45. `[terverifikasi]` Adjuster dan consultant wajib terisi **hanya sebelum baris adjustment pertama** —
    baris ke-2 dan seterusnya **tidak** menuntutnya. Ditiru apa adanya.
46. `[penyimpangan sadar]` ⚠️ `[terverifikasi]` Validasi gagal **tidak meninggalkan baris adjustment yatim**. Di Pega,
    `Activity/AddAdjustment_Act.xml` step **11** keluar activity sebelum step 12 melakukan rollback,
    sehingga baris kosong tertinggal. **⚠️ penyimpangan sadar.**
47. `[terverifikasi]` Baris adjustment menyimpan nama bank, id bank, nomor rekening, dan kode SWIFT.
48. `[penyimpangan sadar]` ⚠️ **Risiko diterima sadar:** perhitungan spreading berjalan untuk **semua** payment type —
    saringan `1|2|5` **tidak ditulis**. **⚠️ penyimpangan sadar (dipertahankan).**
49. `[penyimpangan sadar]` ⚠️ `[keputusan work owner]` **Pesan dan penanda blokir memakai satu besaran yang sama, dalam
    IDR, dibandingkan terhadap estimasi terkini** (AC 40). **Kelebihan TOTAL memblokir penyimpanan;
    kelebihan per baris ditampilkan sebagai peringatan berlabel jelas** — bukan error yang tampak
    memblokir padahal tidak. Bukti dan alasannya di Implementation Decision **14a**.
    **⚠️ penyimpangan sadar.**
    > ⚠️ **RALAT 2026-09-17** — teks lama: *"⚠️ **Risiko diterima sadar:** satu baris adjustment
    > yang melebihi estimasi **tidak ditandai error** selama total masih di bawah estimasi.
    > **⚠️ penyimpangan sadar (dipertahankan).** `[terbuka]` niat aslinya."* Niatnya kini terbaca,
    > dan inkonsistensinya **diperbaiki**, bukan dipertahankan.
125. `[penyimpangan sadar]` ⚠️ `[keputusan work owner]` Sebuah baris penyesuaian **beku** bila ada
    kasus komite yang menunjuknya, **berjalan maupun selesai**. ⛔ **Tidak ada kolom penanda beku**;
    keadaan itu **diturunkan**. Test yang mencari kolom penanda **gagal**.
    **⚠️ penyimpangan sadar.**
126. `[penyimpangan sadar]` ⚠️ `[keputusan work owner]` Baris beku **tidak dapat diubah,
    selamanya**; **setiap jalur sunting menolaknya di lapisan layanan**, bukan hanya di layar.
    **Penolakan komite tidak mencairkannya.** **⚠️ penyimpangan sadar.**
127. `[penyimpangan sadar]` ⚠️ `[keputusan work owner]` Baris beku **tidak dapat dihapus satu per
    satu**; upaya menghapusnya **ditolak**. Perbaikan memakai **baris penyesuaian baru**.
    **⚠️ penyimpangan sadar.**

### Deductible

50. `[keputusan work owner]` Deductible = **MAX(persentase × basis, nilai flat)**.
51. `[terverifikasi]` Basis dipilih antara TSI dan nilai klaim.
52. `[keputusan work owner]` Rumus yang berlaku adalah versi **2024** — deductible dikurangkan
    **sesudah** share ceding. Versi 2022 **ditinggalkan**.
53. `[penyimpangan sadar]` ⚠️ `[terverifikasi]` Properti bernama sama berarti berbeda di dua rule lama; di Go **satu makna
    saja**. **⚠️ penyimpangan sadar.**

### Wewenang

54. `[penyimpangan sadar]` ⚠️ `[keputusan work owner]` Tingkat wewenang dibaca dari **satu sumber**: roster komite. Tidak
    ketemu → **`Claim Admin`**. **⚠️ penyimpangan sadar.**
    > ⚠️ **RALAT 2026-09-18** — frasa penggolongan ditambahkan. Di Pega tingkat wewenang tidak
    > berasal dari satu sumber: `DataTransform/InsertChronology_DT.xml` langkah **1.1.4** memasang
    > nilai dasar tanpa syarat, lalu langkah **1.1.5–1.1.7** menimpanya lewat tiga cabang bernama
    > orang. Membacanya dari roster saja **mengubah asal datanya** — sejajar **AC 55**, yang sudah
    > tergolong penyimpangan sadar sebagai konsekuensinya.
55. `[penyimpangan sadar]` ⚠️ `[keputusan work owner]` **Empat nama orang yang di-hardcode dibuang.** Pergantian pemegang
    jabatan cukup lewat baris tabel. **⚠️ penyimpangan sadar.**
56. `[keputusan work owner]` Batas nilai **Direktur Utama** adalah wewenang **terpisah**, **bukan**
    tingkat kelima tangga komite.
57. ⚠️ Seluruh gerbang wewenang **ditegakkan di lapisan layanan** (**ADR-0014**). Test yang menemukan
    gerbang hanya di UI **gagal**.
58. `[keputusan work owner]` Test: panggil endpoint penyerahan komite langsung, tanpa UI → **ditolak** bila syarat tidak
    terpenuhi.
59. `[terverifikasi]` Jumlah tingkat komite = **COUNT baris roster aktif** yang pita limitnya mencakup
    nilai klaim — dihitung saat penyerahan, **tidak** dari konstanta.
60. `[keputusan work owner]` **`LIMIT_TOP` TIDAK dijadikan penyaring.** Penyaring tetap
    **`LIMIT_BOTTOM` + `STS_KLAIM` + `STS_AKTIF`**. Alasannya: jumlah tingkat komite = jumlah baris
    roster yang lolos; bila `LIMIT_TOP` ikut menyaring, hanya **satu** pita yang lolos dan tangga
    berjenjang hilang. `[terverifikasi]` `ReportDefinition/FilterEmailKomiteWithLimit.xml` memang
    **memilih** `.LIMIT_TOP`, tetapi `pyFilterLogic`-nya `A AND C AND B` dan ketiganya bukan
    `LIMIT_TOP` — dipilih untuk ditampilkan, tidak untuk menyaring.
    > ⚠️ **RALAT 2026-09-17** — teks lama: *"⚠️ `[terverifikasi]` Batas **atas** pita roster **tidak
    > pernah dipakai sebagai filter** di Pega. `[terbuka]` apakah ia seharusnya dipakai."*
    > Pertanyaannya dijawab: tidak, dan sekarang ada alasannya.
61. ✅ `[data DBA]` **Kunci pencocokan pelaku ke roster = `EMAILKOMITE.OPERATOR_ID`**, dikonfirmasi
    sama dengan user ID Pega. **Kolom `USER_ID` baru batal; `EMAIL` tidak dipakai sebagai kunci.**
    **Nol perubahan skema** untuk bagian ini — kolomnya sudah diambil RD roster hari ini.
    > ⚠️ **RALAT 2026-09-17** — teks lama: *"`[data DBA]` Kunci pencocokan pelaku ke roster
    > **menunggu konfirmasi** — kolom email yang sudah ada, atau kolom id pengguna baru."*

### Penyerahan ke Komite dan penutupan

62. `[keputusan work owner]` Penyerahan membuat rekam kasus komite yang membawa **penunjuk stabil** ke baris adjustment.
63. `[penyimpangan sadar]` ⚠️ `[terverifikasi]` Rujukan memakai **ID stabil**, bukan indeks posisi. Di Pega penunjuknya
    indeks posisi, sehingga muatan digemukkan 24 field sebagai snapshot. **⚠️ penyimpangan sadar.**
64. `[terverifikasi]` Penyerahan ditolak bila data bank belum lengkap.
65. `[terverifikasi]` Penyerahan ditolak bila lampiran wajib belum lengkap; daftar lampiran wajib **bergantung jenis
    pembayaran**.
66. `[penyimpangan sadar]` ⚠️ `[terverifikasi]` **Dua jalur ke Komite** ada — penyerahan adjustment dan penutupan tanpa
    pembayaran — dan keduanya menghasilkan kasus anak berkelas sama. **Keduanya dapat aktif pada
    klaim yang sama.** ⚠️ Di Go dipasang **saling-kunci**. **⚠️ penyimpangan sadar.**
67. `[penyimpangan sadar]` ⚠️ `[terverifikasi]` Jalur penutupan tanpa pembayaran memakai **roster tertanam satu orang** di
    Pega. Di Go ia diambil dari **roster data**. **⚠️ penyimpangan sadar.**
68. `[terverifikasi]` Penutupan ditolak bila masih ada baris adjustment yang belum diputus.
69. `[terverifikasi]` Penutupan ditolak bila pengiriman ke Kasir belum berhasil.
70. `[penyimpangan sadar]` ⚠️ **Risiko diterima sadar:** klaim yang **ditolak** komite **tetap** memiliki data di sistem
    luar — pengiriman terjadi tanpa menunggu komite, dan **nol rule pembatal**.
    **⚠️ penyimpangan sadar (dipertahankan).**
128. `[penyimpangan sadar]` ⚠️ `[keputusan work owner]` Menghapus klaim yang **pernah** punya kasus
    komite **ditolak — selamanya**, walau komitenya **sudah selesai maupun menolak**, dan **jalur
    komitenya mana pun — penyerahan baris penyesuaian maupun penutupan tanpa pembayaran**. Pesannya
    **menyebut berapa kasus komite yang pernah ada**. Klaim yang **belum pernah** ke komite **tetap
    terhapus** dan **mengkaskade** — **kecuali klaim lama hasil migrasi**, yang **tidak dapat dihapus
    sama sekali**, apa pun riwayat komitenya. **⚠️ penyimpangan sadar.**
    > ⛔ **RALAT 2026-09-19 — dua tambahan.** Kalimat lamanya **dikutip utuh, tidak dihapus**:
    >
    > > *"128. `[penyimpangan sadar]` ⚠️ `[keputusan work owner]` Menghapus klaim yang **pernah**
    > > punya kasus komite **ditolak — selamanya**, walau komitenya **sudah selesai maupun
    > > menolak**. Pesannya **menyebut berapa kasus komite yang pernah ada**. Klaim yang **belum
    > > pernah** ke komite **tetap terhapus** dan **mengkaskade**."*
    >
    > **① Jalur komite mana pun dihitung** — `[keputusan work owner]` 2026-09-19, menjawab
    > pertanyaan A. Kasus komite jalur **penutupan tanpa pembayaran** tidak menunjuk baris
    > penyesuaian mana pun, jadi ia **tidak membekukan baris** *(AC 125 tidak berubah)*, tetapi ia
    > **tetap mengunci klaimnya**.
    >
    > **② Klaim lama hasil migrasi tidak dapat dihapus sama sekali** —
    > `[keputusan work owner — atas rekomendasi asisten]` 2026-09-19. ⚠️ **Tanda ini sengaja
    > dipakai:** work owner menyerahkan pilihannya kepada asisten di antara dua bacaan — *"hanya
    > klaim migrasi"* atau *"semua klaim"* — dan yang dipilih adalah **yang pertama**, karena ia
    > **tidak mencabut apa pun**. ⛔ **Mudah dicabut** bila yang dimaksud ternyata *semua klaim*;
    > pencabutannya akan menyentuh **AC 5**, separuh **AC 128**, dan butir uji kaskade di tiket 00.
    >
    > **Sebabnya ② ada:** bukti *"pernah punya kasus komite"* di sistem baru **diturunkan dari
    > adanya kasus komite**, sedangkan untuk data lama keberadaan kasus itu bergantung pada apa yang
    > berhasil dimigrasikan — dan kolom yang dulu membuktikannya sudah dibuang **(AC 129)**.
129. `[penyimpangan sadar]` ⚠️ `[keputusan work owner]` **Tidak ada kolom `FLAG_ON_GOING_COMMITTEE`**
    di sistem baru; layar akseptasi menampilkan *"komite sedang berjalan"* **dari kasus komite yang
    berjalan**. Test yang **mengandalkan kolom itu gagal**. ⚠️ **Migrasi tidak membaca kolom itu
    sama sekali** — nilainya **tidak dipakai sebagai bukti** apa pun, termasuk sebagai bukti bahwa
    sebuah klaim pernah ke komite. **⚠️ penyimpangan sadar.**
    > ⛔ **RALAT 2026-09-19** `[keputusan work owner]` — *"abaikan saja properti itu saat migrasi."*
    > Kalimat lamanya **dikutip utuh, tidak dihapus**: *"129. … layar akseptasi menampilkan
    > *"komite sedang berjalan"* **dari kasus komite yang berjalan**. Test yang **mengandalkan kolom
    > itu gagal**."* Yang ditambahkan: migrasi **tidak membacanya**, bukan sekadar **tidak
    > memindahkannya**.

### Klasifikasi lini bisnis

71. `[terverifikasi]` + `[data DBA]` Klasifikasi lini bisnis memakai **satu kunci: `TreatyGroupID`**.
72. `[terverifikasi]` + `[keputusan work owner]` **2026-09-19** Ke-61 rule klasifikasi lini
    **dipindahkan sekali sebagai satu himpunan**; **hidup-matinya tidak ikut dipindahkan**. Sebuah
    rule bernilai benar hanya bila halaman yang diujinya **ada** pada objek kerja modul yang
    memanggilnya. **Pada Claim Prop, 52 dari 61 menguji `BusinessType` pada halaman
    `OfferFacIn`/`Quotation` yang tidak ada di sini, sehingga ke-52-nya tidak pernah bernilai
    benar**; pada Claim Fac In halaman itu ada dan rule yang sama **hidup**. ⭐ **Satu salinan, dua
    nasib.** Test yang menemukan **dua salinan terpisah** dari rule yang sama **gagal**; test yang
    menemukan salah satu dari ke-52 itu **bernilai benar pada objek kerja Claim Prop** **gagal**.
    *(Bab Klasifikasi lini bisnis)*

    > ✅ **RALAT 2026-09-18 TERJAWAB** — `[keputusan work owner]` **2026-09-19**. Butirnya
    > **ditutup**, dan penggolongannya adalah **penerapan aturan berdiri**, bukan penyimpangan
    > tersendiri. ⭐ **Sebabnya: pertanyaannya sendiri terbukti salah pertanyaan.** Selama
    > pilihannya dibingkai *"dimigrasikan atau tidak"*, ketiga aturan induk memang tidak
    > menjawabnya. Keputusan **A5-6 Claim Fac In** membingkainya ulang — **rule dipisahkan dari
    > kehidupannya** — dan dengan bingkai itu tidak ada yang "tidak dimigrasikan": rule-nya
    > dipindahkan, yang tidak dipindahkan adalah **anggapan bahwa ia hidup di sini**.
    >
    > ⛔ **Teks lama dikutip utuh, tidak dihapus:**
    >
    > > *"72. ⚠️ `[terverifikasi]` Rule klasifikasi warisan yang menguji model modul lain **tidak
    > > dimigrasikan** — 52 dari 61 rule `When` (dulu **51 dari 60**) menguji properti yang tidak
    > > ada pada objek kerja Claim Prop."*
    > >
    > > *"⚠️ **RALAT 2026-09-18** — `[terbuka]` **penggolongannya belum diputuskan.** Ketiga
    > > aturan induk menjawab gerbang berflag mati, langkah ber-remark, dan **elemen UI**
    > > selalu-salah — tidak satu pun menyebut sebuah **rule utuh** yang selalu bernilai salah.
    > > Apakah ini penerapan aturan berdiri atau penyimpangan tersendiri diserahkan ke **work
    > > owner**; **tidak memblokir**, karena perilaku yang ditiru sama saja dalam kedua bacaan."*
    >
    > ⚠️ **Akibat pada sensus:** AC 72 **tidak lagi ber-`[terbuka]`** dan **tidak lagi
    > ber-⚠️ penyimpangan**. Baris sensus di kepala berkas disesuaikan.
73. ⚠️ **Cacat template ikut dibawa apa adanya:** Marine Cargo memakai template proyek konstruksi yang
    sama dengan Aneka, **termasuk typo**. Bila kelak diperbaiki, perbaikannya dicatat tersendiri
    dan ditandai ⚠️.
    > ⚠️ **RALAT 2026-09-18** — anak kalimat lama berbunyi *"catat sebagai penyimpangan sadar
    > tersendiri"*; frasa itu menggambarkan tindakan di masa depan, bukan menggolongkan AC ini,
    > sehingga setiap sensus memungutnya keliru. **AC 73 adalah paritas, bukan penyimpangan** —
    > ⚠️-nya tetap berdiri karena ia jebakan.
74. ⚠️ Kode `10004` pada tabel kelompok treaty dan `10004` pada tabel jenis treaty adalah **dua enum
    berbeda**. Test yang menyatukannya **gagal**.

### Tiga prefix dan Syariah

75. `[terverifikasi]` Modul melayani **tiga prefix klaim**; varian ber-`S` diperlakukan **identik**
    dengan induknya.
76. `[terverifikasi]` **Tiga jalur simpan OS akseptasi terpisah** dipertahankan sesuai jenis.

### Jejak audit

77. `[terverifikasi]` Setiap aksi meninggalkan entri berisi **aksi, pelaku, waktu, dan tingkat wewenang**.
78. `[penyimpangan sadar]` ⚠️ Tingkat wewenang disimpan sebagai **enum tertutup**; **nama pelaku di kolom terpisah**.
    **⚠️ penyimpangan sadar.**
79. `[penyimpangan sadar]` ⚠️ Migrasi memetakan **kedua ejaan** tingkat terbawah ke satu nilai enum, dan memindahkan nama ke
    kolom pelaku. **⚠️ penyimpangan sadar.**
    > ⚠️ **RALAT 2026-09-18** — frasa penggolongan ditambahkan. Data warisan menyimpan **dua ejaan**
    > tingkat terbawah dari dua penulis berbeda, dan korpus sendiri mengakalinya dengan pencocokan
    > substring. Memetakan keduanya ke satu nilai **mengubah data tersimpan** — sejajar **AC 78**,
    > yang sudah tergolong penyimpangan sadar atas keputusan enum yang sama.
80. `[penyimpangan sadar]` ⚠️ **Semua peran terekam tanpa kecuali.** Test yang menemukan satu peran tidak berjejak **gagal**.
    **⚠️ penyimpangan sadar.**
81. `[penyimpangan sadar]` ⚠️ **Memo penutupan klaim tercatat di riwayat.** **⚠️ penyimpangan sadar.**
82. `[penyimpangan sadar]` ⚠️ `[data DBA]` Riwayat ditampilkan **urut dari kolom tanggal**, bukan urutan baris.
    **⚠️ penyimpangan sadar.**
    > ⚠️ **RALAT 2026-09-17** — teks lama tidak bertanda ⚠️ dan karena itu tidak terhitung sebagai
    > penyimpangan. Pemeriksaan korpus membuktikan sebaliknya: `[terverifikasi]` grid jejak audit
    > **tidak punya konfigurasi pengurutan sama sekali** — nol tag pengurutan di dalam bloknya, dan
    > nol `pySortOrder` berisi pada `Section/InputAcceptation.xml` maupun
    > `Section/OutstandingClaim.xml`. Pega menampilkan **urutan baris tersimpan**, dan `[data DBA]`
    > satu baris contoh membuktikan urutan itu **bukan kronologis**. Mengurutkan menurut tanggal
    > karena itu **mengubah apa yang dilihat pengguna** — penyimpangan sadar yang sesungguhnya,
    > bukan penegasan perilaku yang sudah ada.
83. `[terverifikasi]` Jejak audit mencatat **jenis aksi**, bukan nilai sebelum/sesudah — batas yang
    diwarisi dan **tidak** diperluas dalam spec ini.

### Dokumen

84. `[keputusan work owner]` **DLA = *Definite Loss Advise*** — satu PDF **per reinsurer per baris
    adjustment**, dengan saudara retro.
85. `[terverifikasi]` Ketiga dokumen **disimpan** ke penyimpanan berkas dengan metadata di database (**ADR-0010**).
86. `[terverifikasi]` **Pencetakan bukan gerbang** — penutupan klaim tidak memeriksa nomor dokumen.
    Ditiru apa adanya.
87. `[terverifikasi]` Gerbang **anti-cetak-ganda** dipertahankan.
88. `[penyimpangan sadar]` ⚠️ `[keputusan work owner]` Cabang **salvage** (step **15.11.2**) dan **adjuster fee** (step
    **15.11.3**) pada DLA **di-remark → tidak dimigrasikan**, tanpa pengecualian. **Step 15.11.1
    diikuti apa adanya**, termasuk pola `TempData` berindeks tetap **(1)**. **Konsekuensi yang
    diterima:** DLA untuk salvage dan adjuster fee menampilkan **nilai klaim**, sama seperti
    produksi. **⚠️ penyimpangan sadar.**
    > ⚠️ **RALAT 2026-09-17** — teks lama menyatakan **ketiga** step tidak ditulis, dan
    > menggantungkan *"`[terbuka]` apakah keduanya seharusnya punya cabang sendiri."* Step 15.11.1
    > **tetap** dimigrasikan, dan pertanyaannya dijawab.

### Efek keluar

89. `[penyimpangan sadar]` ⚠️ Efek keluar memakai **outbox transaksional** — dicatat dalam transaksi yang sama, dikirim dan
    **diulang** bila gagal (**ADR-0015**). **⚠️ penyimpangan sadar.**
90. `[keputusan work owner]` Test: gagalkan klien luar → pekerjaan **tetap ada di outbox** dan terkirim pada percobaan
    berikutnya.
91. `[terverifikasi]` Endpoint **di-resolve saat runtime** dari tabel tautan layanan; **nol URL
    literal** di kode (**ADR-0013**).
92. `[penyimpangan sadar]` ⚠️ `[terverifikasi]` Pemisahan lingkungan **tidak** bergantung pada identitas node aplikasi.
    **⚠️ penyimpangan sadar** — di Pega satu cabang digerbangi nama node JBoss.
93. ✅ `[data DBA]` **Tidak ada bug.** Kasir mengirim **kedua ejaan sekaligus dengan nilai identik**
    (`ReponseCode` **dan** `ResponseCode`; `ResponseMsg` **dan** `ResponseMessage`), sehingga teks
    **"Akseptasi Sudah Masuk ke Kasir"** memang muncul selama ini. **Label cacat dicabut.**
    > ⚠️ **RALAT 2026-09-17** — teks lama: *"⚠️ **Cacat kosmetik dipertahankan:** teks status Kasir
    > dibaca dari properti salah eja, sehingga pesan sukses tidak muncul. **Alur tidak
    > terpengaruh.** `[terbuka]` nama properti yang benar."* Kedua pernyataan itu keliru.
94. `[terverifikasi]` Status pengiriman ke Kasir dapat dilihat pengguna.

### Entitas pendukung

95. `[terverifikasi]` **Adjuster dan Consultant satu master, dua peran** — master tanpa kolom tipe.
    Satu klaim maksimal satu adjuster dan satu consultant.
96. `[terverifikasi]` **Katastrofa punya entitas event sendiri** yang **mengelompokkan banyak klaim**.
97. `[terverifikasi]` **Cause of Loss dua tingkat** dihubungkan kunci induk, ditambah relasi ke lini
    bisnis; **wajib** sebelum simpan.
98. `[penyimpangan sadar]` ⚠️ `[terverifikasi]` Kedua procedure pembaruan Cause of Loss menerima **payload JSON** lewat
    parameter bernama seolah kolom teks. Di Go parameternya **dinamai sesuai isinya**.
    **⚠️ penyimpangan sadar.**

### Efek keluar — tambahan 2026-09-17 (dump respons Kasir)

99. `[penyimpangan sadar]` ⚠️ `[data DBA]` **`CaseIDCashier` disimpan bersama jejak pengiriman ke Kasir.** Itu nomor kasus
    di sistem Kasir — satu-satunya pegangan rekonsiliasi antara klaim kita dan pembayaran di sana.
    `[terverifikasi]` Pega **membuangnya**: sensus 329 berkas menemukan `CaseIDCashier` di **nol
    berkas**. Yang dicatat `Activity/HitServiceToKasir_Act.xml` step **9.8** dan **13.5** ke
    `POOLDATA.DIRECTTOKASIR_LOG` hanya empat kolom — `DATA_JSON`, `IDPEGA` (`pyWorkPage.pzInsKey`,
    kunci internal Pega), `NOAKSEPTASI` (`.AcceptedNo`), dan `KET`
    (`.StatusServiceKasir.ResponseMsg`). **⚠️ penyimpangan sadar** — tanpa ini rantai rekonsiliasi
    dengan Kasir putus.
100. `[penyimpangan sadar]` ⚠️ `[data DBA]` **Kode respons Kasir dibandingkan sebagai string, atau dinormalisasi secara
    eksplisit.** Kasir mengirim `"1"` sebagai **string**; Pega membandingkannya sebagai **angka**
    (`@if(.StatusServiceKasir.ReponseCode=1, …)`) dan bergantung pada konversi diam-diam. Go tidak
    melakukan konversi itu. Test: `"01"` dan `" 1"` **tidak boleh** lolos diam-diam sebagai sukses,
    dan **tidak boleh** gagal tanpa suara. **⚠️ penyimpangan sadar.**
101. `[data DBA]` **Parser menerima kedua ejaan** — `ReponseCode`/`ResponseCode` dan
    `ResponseMsg`/`ResponseMessage` — **dan mencatat di log** bila yang datang hanya ejaan yang
    tidak diharapkan, agar kita tidak bergantung diam-diam pada sisi luar yang tidak kita kendalikan.

### Registrasi dan validasi masuk — tambahan 2026-09-17

> ⚠️ **Catatan cakupan — diperbarui 2026-09-17.** Bab Acceptance Criteria terbitan sebelumnya
> **tidak punya sub-bab registrasi sama sekali**: US **1–9** hanya tertutup sebagian lewat AC 96
> (katastrofa) dan AC 97 (Cause of Loss). AC **102–104** menutup US **3** dan US **4**.
>
> > #### ⚠️ RALAT 2026-09-17 — celah sisanya SUDAH DITUTUP
> >
> > Teks lama berbunyi: *"sisanya (US 1, 2, 5, 6, 7, 8) **masih belum punya AC** dan dicatat di
> > tabel Pertanyaan terbuka sebagai celah cakupan, bukan sebagai pertanyaan."* Itu **tidak lagi
> > benar** — AC **108–122** menutup US 1, 2, 5, 6, 7, 8, dan melengkapi US 9. Peta penuhnya ada
> > di bab **Pertanyaan terbuka di dalam spec**, dan celah cakupan kini **NIHIL**.

102. `[penyimpangan sadar]` ⚠️ `[keputusan work owner]` **Klaim dengan Date of Loss sama pada polis yang sama MEMBLOKIR
     penyimpanan** dan **menampilkan peringatan**. Di Pega niat itu tertulis — keterangan
     `Activity/CheckDateDOL_Act.xml` step **25** berbunyi *"make disable submit if there is similar
     date"* — tetapi step itu memasang `.IsError = 1` sedangkan **kedua pembacanya menguji `> 1`**,
     sehingga peringatan tampil dan penyimpanan tetap lolos. Bukti lengkap di **§15a**.
     **⚠️ penyimpangan sadar.**
103. `[terverifikasi]` **Pengecualian klaim yang sudah ditolak tetap berlaku** — bila klaim lama ada
     di tabel reject, **tidak ada blokir dan tidak ada peringatan**. Pengecualian itu sudah berjalan
     di hulu: `Activity/CheckDateDOL_Act.xml` step **23.1** mengembalikan penanda ke nol setelah
     step **21–22** mengambil data lewat `ReportDefinition/RejectedClaim_RD`. **Alur klaim pengganti
     tidak terganggu**, dan AC 102 tidak mengubahnya.
104. `[penyimpangan sadar]` ⚠️ `[terverifikasi]` **Pembebasan satu nomor polis yang di-hardcode TIDAK dimigrasikan.**
     `Activity/CheckDateDOL_Act.xml` step **20** — sebuah loop atas hasil riwayat klaim — membebaskan
     satu nomor polis tertentu dari pemeriksaan duplikat; literalnya bahkan **diawali spasi**, dan
     nomor itu muncul di **satu berkas** dari 329. Flag prakondisi step 20 adalah **`false`**,
     sehingga pengecualian itu **tidak berlaku hari ini**. ⚠️ Ia **bekas pengecualian bisnis yang
     sengaja dipasang lalu dimatikan**, bukan sisa saringan uji coba. Sejenis dengan empat nama orang
     di **AC 55**. **⚠️ penyimpangan sadar.** `[terbuka]` ringan — lihat **§15b**.
     > ⚠️ **RALAT 2026-09-17** — teks lama menempatkannya pada **step 20.1**. Lokasinya **step 20**;
     > arahnya (membebaskan) sudah benar sejak semula.
105. `[penyimpangan sadar]` ⚠️ `[keputusan work owner]` **Penanda kesalahan dipecah menjadi tiga hal bernama sendiri** —
     *kesalahan yang memblokir* (ya/tidak) · *peringatan yang ditampilkan* (daftar pesan) · *boleh
     kirim ke komite* (ya/tidak). Di Pega ketiganya berbagi satu properti angka dengan nilai
     `0`·`1`·`2`·`3`·`""` yang **artinya tidak tertulis di mana pun**. ⚠️ Menaikkan ambang menjadi
     `>= 1` **bukan** perbaikan yang benar: nilai `3` adalah saklar kirim-ke-komite dan akan ikut
     memblokir hal yang bukan urusannya. Bukti di **§15**. **⚠️ penyimpangan sadar.**
106. `[penyimpangan sadar]` ⚠️ `[keputusan work owner]` **Nama kolom dibuat sesuai isinya**, bukan disalin dari alias
     warisan. `[terverifikasi]` `RDBList/CariHistoryClaim_SQL.xml` memberi enam alias yang **lima di
     antaranya berbohong** — tanggal kejadian dinamai *start date*, nomor klaim dinamai *nama
     cabang*, id Pega dinamai *kode cabang*, penyebab kerugian dinamai *kode bisnis*, estimasi klaim
     dinamai *TSI*. Akibatnya `Activity/CheckDateDOL_Act.xml` step **20.1** membandingkan nomor
     klaim dengan kolom bernama *nama cabang*. Rincian di **§16**. **⚠️ penyimpangan sadar.**
107. `[penyimpangan sadar]` ⚠️ `[terverifikasi]` **Nilai uang tidak pernah disimpan sebagai teks.** SQL riwayat klaim
     menambal estimasi dengan `REPLACE(…, ',', '.')` — angka itu tersimpan sebagai teks berkoma dan
     diperbaiki saat dibaca. Di Go ia **desimal** sejak awal, sehingga **nol tambalan pemisah
     desimal** di jalur baca mana pun. Menguatkan **§4** dan **ADR-0003**. **⚠️ penyimpangan sadar.**

#### Pendaftaran klaim dan nomor polis (US 1–2)

108. `[terverifikasi]` **Nomor polis treaty wajib berformat yang dikenali.**
     `Activity/CheckNoPolicy.xml` step **2** menampilkan *"Format Policy No wrong, please fill the
     right one"*. ⚠️ **Arah gerbangnya terbalik** — syaratnya `@contains(PolicyNo,"RNM-Q")` dengan
     BENAR=**3** (lewati), sehingga pesan justru muncul saat nomor polis **tidak** mengandung
     `RNM-Q`. Aturan sebenarnya: **nomor polis harus mengandung `RNM-Q`**.
109. `[terverifikasi]` **Data polis ditarik dari master saat klaim didaftarkan.**
     `Activity/CheckNoPolicy.xml` step **4** memanggil `RDBList/GetDataNopolisTreatyin.xml` atas tabel
     `policyjson`; `RDBList/SetPolicyTreatyProp.xml` (dipanggil `Activity/SetMasterID.xml`) menarik
     nomor polis, nomor offer, kelompok treaty, sumber bisnis, tahun underwriting, dan kuartal dari
     `TREATYINPRODUCTION`.
110. `[penyimpangan sadar]` ⚠️ `[terverifikasi]` **Dua nama kolom lagi yang berbohong, dan satu tabel tanpa prefiks schema.**
     `RDBList/SetPolicyTreatyProp.xml` memberi alias `TREATYGROUP AS "BusinessName"` — kelompok
     treaty dinamai *nama bisnis* — dan menyisipkan literal `'0' AS "Prodke"`. ⚠️ Tabel yang sama
     dieja **`TREATYINPRODUCTION` tanpa prefiks** di sini tetapi **`pooldata.TREATYINPRODUCTION`** di
     `RDBList/CekPolicyNumber_SQL.xml`. Di Go: nama sesuai isinya, schema selalu eksplisit (**§2**,
     **§16**). **⚠️ penyimpangan sadar.**
111. `[terverifikasi]` **Nomor polis yang tidak ada di master treaty ditolak.**
     `RDBList/CekPolicyNumber_SQL.xml` mencari `IDPEGA` di `pooldata.TREATYINPRODUCTION` menurut
     nomor polis; `Activity/CheeckNoRNM_Act.xml` step **5** menampilkan *"policy number not found,
     please re-choose policy number"* ketika hasilnya **nol baris**
     (`@LengthOfPageList(PolicyList.pxResults)=0`). ⚠️ Gerbang pencariannya (step **4**) memakai
     BENAR=**3** — **terbalik** — sehingga pencarian berjalan justru saat nomor polis **tidak**
     kosong. Keduanya flag `true`, jadi **berlaku**.
112. `[terverifikasi]` **Nomor polis kosong ditolak terpisah dari nomor polis tidak dikenal.**
     `Activity/CheckNopolicy_Act.xml` step **3** menampilkan *"Policy No is null "* saat nomor polis
     kosong (flag `true`, arah normal); step **4** dan step **5** — yang menjalankan pembaruan dan
     transformasi data — hanya berjalan saat nomor polis **tidak** kosong, ditulis dalam dua bentuk
     berlawanan (`==""` dengan BENAR=3, dan `!=""` dengan arah normal).

#### Kelunasan premi dan pembebasannya (US 5–6)

113. `[terverifikasi]` **Premi polis diperiksa lunas sebelum akseptasi.**
     `RDBList/CekLunasPremi_Sql.xml` menjumlahkan `sum(ivd_trans_sign * ivd_total)` atas
     `arasapas.invoice` dan `detail_invoice`. `Activity/CekPremiLunas_Act.xml` step **6** (flag
     `true`) menandai **"BELUM LUNAS"** dan menyiapkan pesan *"Akseptasi tidak dapat dilanjutkan
     dikarenakan Premi belum Lunas"* bila hasilnya **`> 0`**; step **8** menampilkannya.
     ⚠️ **"Lunas" berarti jumlah bertanda ≤ 0, bukan = 0.**
114. `[penyimpangan sadar]` ⚠️ `[terverifikasi]` **Tabel `detail_invoice` tidak diprefiks schema** sementara pasangannya
     `arasapas.invoice` diprefiks — dalam satu perintah yang sama. Di Go keduanya diprefiks eksplisit
     (**§2**). **⚠️ penyimpangan sadar.**
115. ⚠️ `[terverifikasi]` **Pemeriksaan premi tidak punya cabang untuk prefix `CLMNP-`.**
     `Activity/CekPremiLunas_Act.xml` step **2** mengisi parameter untuk prefix **`CLMP-`** dan step
     **3** untuk **`CLM-`**; **tidak ada step untuk `CLMNP-`**. Step **4** tetap menjalankan
     pencarian **tanpa gerbang**, sehingga untuk klaim `CLMNP-` pencarian berjalan dengan parameter
     yang tidak pernah diisi. `[terbuka]` **ringan** — apakah `CLMNP-` memang dikecualikan dari
     pemeriksaan premi, atau ini cabang yang tertinggal? Pemilik: **work owner**. **Tidak memblokir.**
116. `[terverifikasi]` **Pemeriksaan premi dapat dibebaskan lewat data proteksi.**
     `Activity/CekPremiLunas_Act.xml` step **7** — yang gerbangnya `ParamData.HASIL1=="BELUM LUNAS"`
     — memanggil `RDBList/CekProteksiKlaim.xml` di step **7.1**, lalu step **7.2** mengembalikan
     status menjadi **"Lunas"** bila proteksi ditemukan. Jadi pembebasan **hanya dievaluasi ketika
     premi belum lunas**.
117. `[penyimpangan sadar]` ⚠️ `[terverifikasi]` **Dua kode proteksi di-hardcode di dalam SQL.**
     `RDBList/CekProteksiKlaim.xml` menyaring `pooldata.openproteksi_edm` dengan `type = '5'` dan
     `STS_AKSEP = '1'`. Di Go keduanya menjadi **nilai bernama**, bukan literal di dalam perintah.
     **⚠️ penyimpangan sadar.**

#### Cause of Loss dan katastrofa (US 7–9)

118. `[terverifikasi]` **Cause of Loss dipilih dari master dua tingkat yang bertaut.** Tingkat induk
     dibaca `ReportDefinition/BrowseVMCauseOfLoss_RD.xml` (`M_COL_ID`, `COL_DESC`, `OLD_M_COL_ID`);
     tingkat anak dibaca `ReportDefinition/BrowseVDCauseOfLoss_RD.xml`, **disaring oleh `M_COL_ID`
     induk yang sedang dipilih**. Perawatannya lewat `Activity/CNMInsertCauseOfLoss_act.xml` dan
     `Activity/CNMInsertDetailCauseOfLoss_act.xml`.
119. `[terverifikasi]` **Ada tingkat ketiga: relasi Cause of Loss ke lini bisnis.**
     `RDBList/GetLBUID_SQL.xml` menggabungkan `V_D_CAUSE_OF_LOSS_BUSINESS` dengan `BUSINESS` menurut
     `D_COL_ID` tingkat anak. ⚠️ **Kedua tabel tanpa prefiks schema** dan digabung dengan koma gaya
     lama — di Go keduanya diprefiks eksplisit dengan `JOIN` yang tertulis (**§2**).
120. ⚠️ `[terverifikasi]` **Satu rujukan kelas salah ketik.**
     `ReportDefinition/BrowseVDCauseOfLoss_RD.xml` memuat `ASM-FW-CNMFW-Int-V_D_CAUSE_OF_LOSS`
     — **kurang huruf `G`** — berdampingan dengan ejaan benar `ASM-FW-GCNMFW-Int-…` di berkas yang
     sama. `[terbuka]` **ringan** — salah ketik lama; perbaikannya sepele dan **tidak memblokir**.
121. `[terverifikasi]` **Simpan ditolak bila Cause of Loss belum diisi.**
     `Activity/ProteksiData_act.xml` step **14** — keterangannya *"make sure cause of loss has been
     choose"*, flag **`true`**, arah normal — menampilkan *"Please fill Cause of loss"* saat
     `pyWorkPage.ClaimData.CauseOfLoss` kosong. ⚠️ Rule yang sama menyiapkan **dua pesan kembar**
     (`ErrMsg8` dan `ErrMsg9`, keduanya berbunyi *"Please Cause of loss"*) dan memakai penanda angka
     kedua `Local.Error` dengan ambang `>1` — pola yang sama dengan yang dipecah di **§15**.
122. `[penyimpangan sadar]` ⚠️ `[terverifikasi]` **Penanda katastrofa ya/tidak dicatat terpisah dari event-nya** —
     melengkapi **AC 96**, yang hanya menutup entitas event dan pengelompokannya.
     `Activity/SetDefNonCatastrope_Act.xml` memakai `ClaimData.StsKatastrofe` bernilai
     `"Catastrophe"` dan `ClaimData.NonKatastrofeType` bernilai `"Claim"`. ⚠️ Nama propertinya
     bahasa Indonesia sedangkan **nilainya** bahasa Inggris, dan modul mengeja konsep ini dalam
     **empat bentuk** — `Catastrope` (82 kemunculan) · `Catastrophe` (16) · `Catasrtope` (8) ·
     `Catastrofe` (3 berkas). Di Go: **satu ejaan, satu enum**. **⚠️ penyimpangan sadar.**

### Migrasi baris yang berlaku — tambahan 2026-09-18

123. `[penyimpangan sadar]` ⚠️ `[data DBA]` **Migrasi memilih satu baris yang berlaku per klaim**
     dari `OS_AKSEPTASI_KLAIM` yang menumpuk, memakai urutan `TANGGAL` **naik** lalu
     `DATA_JSON.AcceptedNo` **turun** — urutan yang sama dipakai
     `RDBList/DataOutstandingTreatyin.xml`. Baris sisanya **disimpan sebagai riwayat**, tidak
     dibuang. **⚠️ penyimpangan sadar.**
     **Alasan menyimpang:** kedua stored procedure penulis menjalankan
     `SELECT count(1) INTO id_count` lalu **tidak pernah memakainya**, dan di
     `PEGA_JSON_OS_AKSEP_KLAIMTNP` cabang `UPDATE`-nya dikomentari habis — sehingga satu `CASEID`
     punya banyak baris **tanpa penanda mana yang berlaku**.
     > ⚠️ **RALAT 2026-09-18** — AC ini **baru**, bukan penggantian. Ia menutup celah yang
     > ditemukan saat menanam fakta Oracle ke tiket 00: **AC 9** hanya menuntut migrasi dari bentuk
     > JSON ke bentuk relasional, dan **tidak** menuntut pemilihan baris ketika satu `CASEID`
     > memiliki banyak baris. Bunyi AC 9 tidak diubah.

### Urutan pemanggilan procedure penulis — tambahan 2026-09-18

124. `[penyimpangan sadar]` ⚠️ `[data DBA]` **Panggilan stored procedure penulis adalah langkah
     TERAKHIR sebelum `COMMIT`** dalam satu transaksi, tanpa pekerjaan penting yang masih
     menggantung sebelumnya. Aplikasi memeriksa `StsSimpan` — **`1` berhasil, `0` gagal** — dan
     pada `0` melaporkan **gagal**. **⚠️ penyimpangan sadar.**
     **Alasan menyimpang:** ketiga procedure memasang `EXCEPTION WHEN OTHERS THEN ROLLBACK` di
     dalam dirinya sendiri; `ROLLBACK` telanjang di Oracle membatalkan **seluruh transaksi**,
     termasuk pekerjaan aplikasi sebelumnya, dan `ROLLBACK TO SAVEPOINT` tidak menolong karena yang
     dipasang bukan itu. Pada `StsSimpan = 0` transaksinya **sudah dibatalkan procedure**, sehingga
     `COMMIT` sesudahnya **tidak menyimpan apa pun**.
     **AC 6** menuntut **hasil**, bukan **urutan**; urutan itulah yang wajib di sini. Bukti dan
     rinciannya di **Implementation Decision §17** dan **§17a** — bunyi keduanya, dan bunyi AC 6,
     **tidak diubah**; AC ini hanya menautkannya.
     > ⚠️ **RALAT 2026-09-18** — AC ini **baru**, bukan penggantian. Ia lahir dari usulan yang
     > dilaporkan bersama §17a dan disetujui `[keputusan work owner]` pada hari yang sama.

### Lintas modul — lima keputusan 2026-09-19

*Keenam AC di bawah lahir dari **satu blok keputusan work owner** tertanggal 2026-09-19, sesudah
modul ini diadu dengan **15 ADR** untuk pertama kalinya. Buktinya di `utang-lintas-modul.md`.*

130. `[keputusan work owner]` **ADR-0003** Uang yang menyeberang ke modul komite **tetap desimal di
     kedua sisi**. Test yang menemukan nilai uang melewati **bilangan pecahan biner** di jalur mana
     pun — termasuk jalur pembuatan kasus komite — **gagal**. ⚠️ `[terverifikasi]` Di Pega,
     `AddKomiteTreatyChild_ACT` mendeklarasikan **`TSISpread`** dan **`ClaimSpread`**, dan
     `SetKomiteTreaty_ACT` mendeklarasikan **`TotalAdjustment`**, bertipe **`Double`** — dan
     **kedua sumber tipe sepakat** *(tanda tangan dan medan parameter)*. ⛔ **Itu cacat deklarasi
     Pega, bukan keputusan penyimpanan**: nilainya tidak pernah *harus* melewati pecahan biner hanya
     karena Pega menamainya begitu. *(Bab Efek keluar · Penyerahan ke Komite)*
131. `[keputusan work owner]` **Tabel pemetaan wajib ada** untuk ketiga parameter itu —
     **`TSISpread`** · **`ClaimSpread`** · **`TotalAdjustment`** — memasangkan **nama dan tipe lama**
     dengan **bentuk barunya**. Test yang menemukan ketiganya berpindah **tanpa tabel pemetaan**
     **gagal**. ⚠️ **Sebabnya sama dengan AC 106:** tanpa tabel itu, orang yang membandingkan
     keluaran lama dan baru akan mengira **datanya berubah**, padahal hanya **bentuknya** yang
     dibetulkan. *(Bab Penyerahan ke Komite)*
132. `[keputusan work owner]` **ADR-0009** **Seluruh data klaim dipindahkan; tidak ada koeksistensi
     dua penulis.** Test yang menemukan klaim **berjalan** diselesaikan di sistem lama sementara
     sistem baru sudah menerima klaim baru **gagal**. ⚠️ **Alasan khas modul ini:** **AC 128**
     menetapkan klaim yang pernah punya kasus komite **tidak dapat dihapus selamanya** dan **AC 126
     · 127** menetapkan baris beku **tidak mencair**. Dua penulis atas klaim yang sama, dengan aturan
     beku yang **hanya dipahami satu sisi**, **akan** melahirkan selisih yang tidak dapat
     diterangkan. *(Bab Migrasi)*
133. `[penyimpangan sadar]` ⚠️ `[keputusan work owner]` **Wewenang menghapus catatan kronologi
     adalah izin eksplisit lewat peran**, **bukan** perbandingan terhadap **teks jabatan** operator.
     Test yang menemukan wewenang bergantung pada **isi medan jabatan** **gagal**.
     **Alasan menyimpang:** `[terverifikasi]` di Pega ia berupa perbandingan terhadap teks jabatan,
     dan itu **satu-satunya pemeriksaan wewenang sejati di modul ini** — rule yang **sama persis**
     ada di Komite Claim Prop. ⛔ Teks jabatan adalah **data kepegawaian**: begitu satu huruf
     diganti, wewenangnya **berubah diam-diam**. *(Bab Wewenang)*
134. `[penyimpangan sadar]` ⚠️ `[keputusan work owner]` **Pengenal berkas dibuat dengan cara yang
     sudah ditambal** — berketelitian **nanodetik** **dan** disertai **nilai unik sejagat** dari
     basis data. Test yang menemukan **dua unggahan berdekatan menghasilkan pengenal yang sama**
     **gagal**. **Alasan menyimpang:** `[terverifikasi]` versi yang dipakai modul ini bersandar pada
     ketelitian **milidetik tanpa nilai unik**, sehingga dua unggahan dalam milidetik yang sama
     menghasilkan **pengenal kembar** — cacat yang **sudah dicatat sendiri** oleh modul ini pada
     ronde 2. ⭐ **Tambalannya bukan rancangan baru:** ia **sudah berjalan** pada salinan rule yang
     sama di modul saudara, yang **4 tahun lebih baru**. *(Bab Dokumen)*
135. `[terverifikasi]` + `[keputusan work owner]` Daftar **jenis berkas yang diterima** mengikuti
     daftar lama **apa adanya — tidak ditambah**. Test yang menemukan jenis berkas di luar daftar
     lama diterima **gagal**. ⚠️ Salinan rule penentu jenis berkas yang lebih baru di modul saudara
     memuat **48** baris lawan **42** di sini *(catatan pengembangnya berbunyi "add avi")*;
     ⛔ **menambah jenis yang diterima adalah PERUBAHAN PERILAKU**, dan tidak ada yang memintanya.
     *(Bab Dokumen)*

---

## Pertanyaan terbuka di dalam spec

### ✅ Ketujuh `[terbuka]` terbitan pertama SUDAH DIJAWAB — 2026-09-17

> Blok ini menyimpan pertanyaan lama apa adanya. Tidak ada di antaranya yang masih berstatus terbuka.

| # | Pertanyaan lama | Jawaban | AC |
| --- | --- | --- | --- |
| 1 | Niat gerbang `Activity/CountValueADJTreaty_Act.xml` step **16** | `[keputusan work owner]` hanya kelebihan **total** yang memblokir; per baris menjadi peringatan — dan inkonsistensi pesan-versus-blokir **diperbaiki** | **49** |
| 2 | Nama properti respons Kasir yang sebenarnya | `[data DBA]` **tidak ada bug** — Kasir mengirim **kedua** ejaan dengan nilai identik | **93** |
| 3 | Cabang salvage & adjuster fee pada DLA | `[keputusan work owner]` step 15.11.2 dan 15.11.3 tidak dimigrasikan; **step 15.11.1 diikuti apa adanya** | **88** |
| 4 | Haruskah `LIMIT_TOP` dipakai menyaring | `[keputusan work owner]` **tidak** — penyaring tetap `LIMIT_BOTTOM` + `STS_KLAIM` + `STS_AKTIF`, agar tangga berjenjang tidak runtuh jadi satu pita | **60** |
| 5 | Sumber pengisi nilai estimasi pada baris adjustment | `[keputusan work owner]` ditemukan; pagar mengikat estimasi **terkini**, potret tetap disimpan untuk audit | **40** |
| 6 | Salah ketik `IsCustomBonds` versus `IsCustomBond` | `[terverifikasi]` **bukan salah ketik** — rule hilang keenam; `When/IsCustomBonds.xml` kini ada, dan nama `When` menggantung turun 1 → **nol** | — |
| 7 | Kunci pencocokan pelaku ke roster | `[data DBA]` **`EMAILKOMITE.OPERATOR_ID`**, sama dengan user ID Pega; kolom `USER_ID` batal | **61** |

### Status sekarang

```
[terbuka] yang memblokir : NIHIL
[terbuka] ringan         : empat
celah cakupan            : NIHIL — 66 dari 66 user story tertutup AC
```

> ⚠️ **RALAT 2026-09-19 — daftar `[terbuka]` disusun ulang sesudah ronde 4 dan 5.**
>
> ```
> [terbuka] yang memblokir : NIHIL
> [terbuka] ringan         : ENAM   (lima lama + satu baru)
> celah cakupan            : NIHIL
> ```
>
> **HILANG — nol.** Tidak ada butir `[terbuka]` ringan yang hilang dari tabel di bawah.
> Dua butir yang sempat digantung ronde 4 **ditutup oleh §A ronde 5 dan tidak pernah masuk tabel
> ini**: *nasib nilai lama share ceding* (ditutup — perkaliannya tidak pernah berjalan) dan
> *arti label `//`* (ditutup — dipastikan dari layar Pega).
>
> **BERTAMBAH — satu, yaitu butir 6 di bawah.**
>
> **RALAT 2026-09-19 (ronde 6)** — sesudah sapu bersih korpus, **tiga butir ringan lagi
> bertambah** (butir **7 · 8 · 9**). Jumlah yang berlaku menjadi **SEMBILAN**.
> **Hilang: nol.** **Yang memblokir: tetap NIHIL.** **Jumlah AC tidak berubah** — temuan ronde 6
> dicatat sebagai keterangan dan butir terbuka, bukan sebagai AC baru.
>
> ⚠️ Angka tabel *"empat"* di blok sebelumnya adalah tulisan lama; jumlah yang berlaku sejak
> 2026-09-19 adalah **enam**. ⛔ Selisih empat-versus-lima pada terbitan sebelumnya **tidak
> diperbaiki di sini** — ia milik RALAT 2026-09-18 (b) dan dibiarkan apa adanya.

| # | `[terbuka]` ringan | Pemilik | Memblokir? |
| --- | --- | --- | --- |
| 1 | `[terverifikasi]` Enam kolom `STS_*` ada di roster komite — `.STS_ADJ` · `.STS_ADJUSTER` · `.STS_REG` · `.STS_REJECT` · `.STS_SALVAGE` · `.STS_SURVEY` — tetapi **nol** dipakai menyaring oleh rule mana pun (sensus 329 berkas). Apakah wewenang **per jenis aksi** memang dimaksudkan berlaku, atau keenam kolom itu warisan yang tidak terpakai? | Finance + work owner | **tidak** — yang ditiru adalah perilaku yang berjalan, yaitu tiga penyaring |
| 2 | `[terverifikasi]` `Activity/CheckDateDOL_Act.xml` **step 20** membebaskan **satu nomor polis yang di-hardcode** dari pemeriksaan duplikat Date of Loss (flag prakondisinya `false`, jadi mati; literalnya diawali spasi; muncul di 1 dari 329 berkas). ⚠️ Bekas **pengecualian bisnis yang sengaja dimatikan**, bukan sisa uji coba — karena itu layak ditanyakan. Masih dikehendaki? Bila ya, jadikan baris data seperti **AC 55**? Lihat **§15b**. | work owner | **tidak** — default yang ditulis adalah tidak dimigrasikan (**AC 104**) |
| 3 | `[terverifikasi]` `Activity/CekPremiLunas_Act.xml` mengisi parameter pemeriksaan premi untuk prefix `CLMP-` (step 2) dan `CLM-` (step 3) tetapi **tidak untuk `CLMNP-`**, sedangkan pencariannya (step 4) berjalan tanpa gerbang. Apakah `CLMNP-` memang dikecualikan, atau cabangnya tertinggal? Lihat **AC 115**. | work owner | **tidak** |
| 4 | `[terverifikasi]` `ReportDefinition/BrowseVDCauseOfLoss_RD.xml` memuat rujukan kelas salah ketik `ASM-FW-CNMFW-Int-…` (kurang huruf `G`) berdampingan dengan ejaan benar di berkas yang sama. Lihat **AC 120**. | pemilik export Pega | **tidak** — perbaikannya sepele |
| ~~5~~ | ✅ **DITUTUP `[keputusan work owner]` 2026-09-19 — penerapan aturan berdiri.** ⭐ Sebabnya: pertanyaannya sendiri terbukti **salah pertanyaan** — selama dibingkai *"dimigrasikan atau tidak"*, ketiga aturan induk memang tidak menjawabnya; keputusan **A5-6 Claim Fac In** *(rule dipisahkan dari kehidupannya)* membingkainya ulang menjadi **satu salinan, dua nasib**. ⛔ Teks lamanya dikutip utuh, tidak dihapus: ~~`[terverifikasi]` **Penggolongan AC 72 belum dapat diputuskan dari aturan berdiri.** Ketiga aturan induk menyebut gerbang berflag mati, langkah ber-remark, dan **elemen UI** selalu-salah — tidak satu pun menyebut sebuah **rule utuh** yang selalu bernilai salah. Penerapan aturan berdiri, atau penyimpangan tersendiri? Lihat **AC 72** dan tiket 04.~~ | work owner — **sudah menjawab** | **tidak** |
| 6 | ⭐ **BARU 2026-09-19.** `[terverifikasi]` `Activity/CountPersen_act.xml` — label `//` berada pada langkah **5** (induk) menurut ekspor, sedangkan catatan lisan work owner menyebut **5.1**. Kalau yang mati hanya 5.1, maka langkah **5.2.1** tetap berjalan dengan `Local.Value` yang **tidak pernah diisi** di blok itu, sehingga menulis `.ClaimSpreaded` dari nilai kosong. Mana yang benar? | work owner | **tidak** — dalam kedua bacaan nilai yang dipakai hilir berasal dari **blok 6**, yang berjalan belakangan dan menimpa |
| 7 | **BARU 2026-09-19 (ronde 6).** `[terverifikasi]` `Activity/CheckPeriodPolicy_Act.xml` langkah **1** ber-remark, padahal ia **satu-satunya** pengisi `Local.CompareDate` dan `Local.EqualDate`; langkah 2 dan 4 membacanya. Bentuk sama persis dengan butir 6. Apakah pemeriksaan periode polis ini memang dimaksudkan mati seluruhnya? | work owner | **tidak** — pemeriksaannya mati dalam kedua bacaan |
| 8 | **BARU 2026-09-19 (ronde 6).** `[terverifikasi]` `.ClaimData.PeriodPolicyTBA` dan pesan *"Policy period is TBA"* ada di 3 berkas, tetapi **konsep periode polis "belum pasti" belum pernah tercatat**. Apakah klaim berperiode TBA memang boleh berjalan sampai akhir? | work owner | **tidak** — perilaku yang ditiru adalah yang berjalan |
| 9 | **BARU 2026-09-19 (ronde 6).** ⚠️ `Activity/CountValueADJTreaty_Act.xml` bercatatan *"samain dengan prod"*, menyiratkan rule ini pernah berbeda antar-lingkungan. Masih berbeda? Yang mana yang benar? | work owner | **tidak** — korpus hanya memuat satu salinan |
> ⚠️ **RALAT 2026-09-18 (a)** — butir **5** dan **6** sebelumnya hanya tercatat di dalam blok RALAT
> (AC 72) dan di tiket 00 (ketujuh kolom), sehingga **tidak terlihat** oleh sensus berjendela baku
> yang mengecualikan baris berawalan `>`. Keduanya kemudian berdiri di tabel ini.
>
> ⚠️ **RALAT 2026-09-18 (b)** — **butir 6 kini DITUTUP** dan pindah ke tabel *sudah dijawab* di
> bawah. Teks lamanya: *"`[data DBA]` **Arti tujuh kolom penanda di kedua tabel proyeksi belum
> dijelaskan DBA** dan tidak boleh disimpulkan dari namanya: `STS_KONVERSI` (di kedua tabel) ·
> `STS_REJECT` · `STS_DLA` · `IDPROD` · `MASTERID` · `CLAIMOLD` · `TGL_KONVERSI`."* Jumlah butir
> ringan **enam → lima**; tidak ada yang memblokir.

### ✅ Butir yang DITUTUP sesudah terbitan pertama — 2026-09-18

| # lama | Pertanyaan | Jawaban | Di mana |
| --- | --- | --- | --- |
| 6 | Arti tujuh kolom penanda pada `OS_AKSEPTASI_KLAIM` dan `JSON_KLAIM` | `[terverifikasi]` `STS_REJECT` = jenis aksi (**0** outstanding · **2** acceptation · **4** close) · `MASTERID` bersumber **dua** tempat menurut jalurnya · `STS_KONVERSI`, `STS_DLA`, `CLAIMOLD` **nol penulis** di 329 berkas sehingga masuk **kosong**. `[data DBA]` `IDPROD` dipaku **1**, `TGL_KONVERSI` dipaku **NULL**. `[keputusan work owner]` **ketujuhnya ditiru apa adanya — bukan cacat, bukan penyimpangan** | **§1a** dan **§1b**; tiket **00** |

⚠️ **Sisa untuk DBA hanya dua, dan keduanya TIDAK memblokir tiket 00:** domain nilai `STS_KONVERSI`
**bila kolomnya dipakai modul lain**, dan arti `IDPROD` di `POOLDATA.KODE_PRODUKSI` — **bukan di
sini**, karena di sini nilainya dipaku **1**.

### ✅ Celah cakupan DITUTUP — peta 66 user story terhadap AC

> ⚠️ **RALAT 2026-09-17** — terbitan sebelumnya mencatat *"celah cakupan: satu (US 1, 2, 5, 6, 7, 8
> belum punya AC)"*. Celah itu **ditutup** oleh AC **108–122**. Tabel lama diganti peta penuh di
> bawah.

| Kelompok US | US | AC yang menutupinya |
| --- | --- | --- |
| Registrasi dan validasi masuk | **1–9** | 96 · 97 · 102 · 103 · 104 · 108 · 109 · 110 · 111 · 112 · 113 · 114 · 115 · 116 · 117 · 118 · 119 · 120 · 121 · 122 |
| Insured Interest (TSI) | **10–13** | 25 · 27 · 41 · 42 · 43 |
| Estimasi | **14–19** | 34 · 35 · 36 · 37 · 38 · 39 · 40 |
| Loss Allocation | **20–23** | 29 · 30 · 31 · 33 |
| Baris adjustment | **24–29** | 44 · 45 · 46 · 47 · 48 · 49 · 105 |
| Spreading | **30–32** | 20 · 21 · 22 · 24 · 32 |
| Deductible | **33–35** | 50 · 51 · 52 · 53 |
| Penyerahan ke Komite | **36–41** | 54 · 57 · 58 · 59 · 60 · 61 · 62 · 63 · 64 · 65 · 66 · 67 |
| Penutupan klaim | **42–44** | 68 · 69 · 70 · 81 |
| Dokumen | **45–49** | 13 · 14 · 16 · 84 · 85 · 86 · 87 · 88 |
| Efek keluar | **50–54** | 89 · 90 · 91 · 92 · 93 · 94 · 99 · 100 · 101 |
| Wewenang | **55–58** | 54 · 55 · 56 · 57 · 59 · 61 |
| Jejak audit | **59–62** | 77 · 78 · 79 · 80 · 82 · 83 |
| Penyimpanan dan pelaporan | **63–65** | 1 · 2 · 3 · 4 · 5 · 6 · 7 · 8 · 9 · 10 · 11 · 12 · **123** |
| Lintas-lini | **66** | 75 · 76 |

**User story yang belum tertutup: NOL.**

⚠️ Catatan kejujuran: peta ini memetakan **kelompok** US ke AC, bukan satu-lawan-satu per US. Setiap
kelompok diperiksa punya sekurang-kurangnya satu AC yang menyatakan perilakunya; beberapa AC melayani
lebih dari satu US, dan itu wajar. Yang dijamin tabel ini adalah **tidak ada kelompok yang kosong** —
bukan bahwa tiap kalimat US punya AC kembarannya sendiri.

---

## Out of Scope

- **Isi konteks Komite Claim Prop.** Modul ini membangun **batasnya**, bukan isinya. Tangga
  persetujuan, routing per tingkat, dan penomoran akseptasi milik konteks sebelah.
- **Claim Non Prop dan Claim Fac In.** ⚠️ `[terverifikasi]` Beberapa rule di modul ini mengeksekusi
  cabang untuk prefix modul lain; cabang itu **tidak dimigrasikan ke konteks ini**.
- **Sistem inti reinsurance non-life, Kasir, dan Arasapas.** Diperlakukan sebagai sistem luar di balik
  kontrak.
- **Lini Syariah sebagai lini terpisah.** Varian ber-`S` diperlakukan **identik** dengan induknya;
  aturan perhitungan Syariah tersendiri **tidak** dalam lingkup.
- **Perbaikan cacat pada butir 14 Implementation Decisions.** Ketiganya **dipertahankan** atas
  keputusan work owner.
- **Layar pemeliharaan master** yang hanya terjangkau dari portal admin dan tidak dicapai dari alur
  klaim.
- **Isi fungsi serialisasi JSON Pega.** Tidak ada di korpus; hanya perilakunya yang terbaca dari data.

---

## Further Notes

### Urutan pengerjaan yang disarankan

**Tiket 00 adalah PREFACTOR** — skema relasional + migrasi. Tidak ada tiket lain yang dapat selesai
sebelumnya, karena seluruh perilaku lain menulis ke tabel yang belum ada. Tiga hal menempel padanya:
batas transaksi penomoran, pemutusan rantai presisi, dan prefix schema eksplisit.

### Sapu bersih korpus — 2026-09-19 (ronde 6)

Keempat puluh dua berkas yang belum pernah disebut ronde 1-5 dibuka **satu per satu**.
**Nol berkas yatim** — keempat puluh dua punya perujuk.

#### Yang MENGUBAH pembacaan — sebelas

1. ⚠️ `[terverifikasi]` **Proteksi periode polis DIMATIKAN dengan sengaja.**
   `Activity/CheckReportDate_Act.xml` — catatan pengembangnya *"matikan protek dalam periode
   polis"*, dan langkah **11** serta **14** (pemeriksa Report Date terhadap Start/End Date Policy)
   **ber-remark**. Pasangannya `Activity/CheckDateReceived_Act.xml` langkah **10** dan **13**
   ber-remark dengan pola sama. **Report Date dan Date Received tidak lagi diperiksa terhadap
   periode polis.** Ini **bukan cacat** — dimatikan sengaja; tetapi ia perilaku yang berjalan,
   jadi ditiru apa adanya.

2. ⚠️ `[terverifikasi]` **Periode polis boleh "TBA".** `Activity/AddAdjustment_Act.xml` memuat
   pesan *"Policy period is TBA. Please verify the dates."* dan properti
   `.ClaimData.PeriodPolicyTBA` dipakai `Section/InputAcceptation.xml` serta
   `Section/OutstandingClaim.xml`. **Konsep ini belum pernah tercatat** di spec maupun tiket.

3. ⚠️ `[terverifikasi]` **`Activity/CheckPeriodPolicy_Act.xml` — pengisi menggantung, kejadian
   KEDUA.** Langkah **1** ber-remark, padahal ia satu-satunya pengisi `Local.CompareDate` dan
   `Local.EqualDate`; langkah 2 dan 4 membacanya. Bentuknya sama persis dengan `CountPersen_act`
   blok 5. `[terbuka]` — lihat `grilling-ronde-6.md` §B.

4. ⚠️ `[terverifikasi]` **Nomor komite diiris berdasarkan posisi karakter.**
   `Activity/SetKomiteNo_Act.xml` langkah 1: nomor komite diambil dengan memotong karakter ke-19
   sampai ke-30 dari **kunci internal Pega** kasus anak terakhir. Menguatkan penyimpangan
   *"rujukan memakai ID stabil"* yang sudah berdiri di tiket 11.

5. ⚠️ `[terverifikasi]` **Teks bebas pengguna diturunkan huruf secara merusak.**
   `Activity/MakeLowercase_Act.xml`: `.ClaimData.Location` dan `.ClaimData.ReportDescription`
   ditimpa versi huruf kecilnya — **nilai asli tidak disimpan**.

6. `[terverifikasi]` **Batas waktu panggilan penyimpanan berkas = 300 000 milidetik (5 menit).**
   `ConnectREST/ServiceGoogle.xml`, catatan pengembang *"buat jadi 300rb"*.

7. ⚠️ `[terverifikasi]` **`Activity/SetPayableTreaty_Act.xml` tidak pernah menyimpan.** Langkah
   **4** (`Obj-Save`) dan **5** (`Obj-Refresh-And-Lock`) ber-remark, juga langkah 1.8. Catatan
   pengembangnya *"BUKA PROTEKSI SALVAGE KASIR"*. Termasuk daftar **jalur simpan mati**.

8. ⚠️ `[terverifikasi]` `Activity/CountValueADJTreaty_Act.xml` bercatatan *"samain dengan prod"* —
   menyiratkan rule ini pernah **berbeda antar-lingkungan**. `[terbuka]` apakah masih berbeda.

9. `[terverifikasi]` `Activity/SetDLACedingSOB.xml` menyalin `DLANoCeding` dan `DLANoSOB` ke
   halaman kerja, masing-masing bergerbang *"properti ada isinya"* — dua nomor dokumen yang belum
   tercatat di bab dokumen.

10. `[terbuka]` `RDBList/GetMObyNopol_SQL.xml` bercatatan *"UBAH BIAR GA SINGLE ROW"*; satu
    penanda baris-tunggal masih terbaca di berkasnya. Belum dapat dipastikan mana yang berlaku.

11. `[terverifikasi]` `Activity/RemoveLossAlloction_act.xml` menulis kronologi *"Delete Loss
    Allocation"* — satu jenis peristiwa jejak audit yang belum terdaftar.

#### Enam pasang nama kembar — TIGA sepasang, TIGA kebetulan nama

⚠️ `[terverifikasi]` Dibuktikan lewat `pxInsName`, bukan nama berkas:

| Pasangan | Hasil |
| --- | --- |
| `PrintFileDLA` (Section + FlowAction) · `ViewAttachment` (Section + Harness) · `SummaryOutsClaim`/`SummaryOutSClaim` (Section + Harness) | **rule yang SAMA** — satu `pxInsName`, dua berkas ekspor |
| `ListPaymentClaim_SC` vs `ListPaymentClaim_Harness` · `ListPolicyNoTreaty` vs `ListPolicyNoTreaty_Harness` · `ViewDetailPayment_SC` vs `ViewDetailPayment_FA` | **KEBETULAN NAMA** — `pxInsName` berbeda, rule berbeda |

#### Lima rule `When` terakhir — nol yang hidup

`[terverifikasi]` `IsMBD` · `IsMBU` · `IsMarineHull` · `isBillboardNeonSyariah` · `isMaintenance`
seluruhnya menguji `pyWorkPage.Quotation.*` atau `pyWorkPage.OfferFacIn.*` — **halaman yang tidak
ada** pada objek kerja Claim Prop. Kelimanya termasuk ke-52 rule yang **mati di modul ini dan hidup
di Claim Fac In**; **isi AC 72 tentang kelimanya tidak berubah**.

> ⛔ **RALAT 2026-09-19** — kalimat lamanya **dikutip, tidak dihapus**: *"Kelimanya termasuk 52
> rule sisa impor; **AC 72 tidak berubah**."* Sebutan *"sisa impor"* diganti mengikuti keputusan
> **satu salinan, dua nasib**; ⚠️ **AC 72 sendiri memang berubah** pada hari yang sama — yang
> tidak berubah adalah **kedudukan kelima rule ini di dalamnya**.

> ⚠️ **Pelajaran parser, tercatat karena hampir menggigit keempat kalinya.** `isMaintenance`
> tampak **tanpa kondisi** bila dibaca dari medan *viewer* (`pyConditionValue1StringLabel` berisi
> pola kosong `[first value][relation][second value]`). Kondisi sebenarnya ada di
> **`pyConditionString`**: `pyWorkPage.Quotation.BusinessType = "Maintenance"`. Aturan *"isi viewer
> adalah sisa basi"* berlaku dan terbukti lagi.

### Aturan baca korpus yang berlaku### Aturan baca korpus yang berlaku

Sembilan aturan baca (R1–R9) dan empat aturan kerja tercatat di `grilling-ronde-2.md` §0 dan §13.
Yang paling menentukan saat membaca ulang korpus:

- Gerbang step hanya berlaku bila flag prakondisinya `true`; kode arah dapat **terbalik**; kode arah
  **kosong** berarti *continue when*.
- `pyStepsBlockName` = `//` berarti remark; nilai lain adalah **label blok**.
  ⭐ **Buktinya `[data work owner 2026-09-19]`, bukan korpus.** Ekspor Pega **tidak memuat medan apa
  pun yang mematikan langkah** — `//` hanyalah isi label. Work owner memastikannya dari layar Pega
  (`KomitePostAdjustment`, step 16 anak 1–4 **tampak dicoret**). Di Claim Prop ini menyangkut
  **68 dari 1366** langkah, **32** di antaranya tanpa syarat apa pun — termasuk tiga `Obj-Save`
  dan satu `Commit` yang karenanya **tidak pernah berjalan** (temuan: jalur simpan mati, **bukan**
  cacat).
- Di Section, syarat tampil hanya berlaku bila selektor modenya menunjuk ke syarat itu.
- Di rule `When`, isi *viewer* adalah **sisa basi** — yang berlaku ada di daftar kondisi.
- `pzOriginalInstanceKey` menyimpan **asal save-as**, bukan pemanggil.
- Sensus nama rule **wajib tidak peka huruf besar-kecil**.

> #### RALAT 2026-09-19 (ronde 6) — dua angka `grilling-ronde-2.md` diselaraskan di sini
>
> `grilling-ronde-2.md` **READ-ONLY** dan tidak ditambal; penyelarasannya ditulis di spec.
>
> | Ditulis ronde 2 | Yang berlaku | Dihitung ulang di |
> | --- | --- | --- |
> | *"89 memakai kode 5/6"* | **23** — kode **5** = **7** baris, kode **6** = **16** baris | `periksa-ulang-aturan-baru.md`, dikuatkan ronde 4 |
> | *"380 normal"* | **363** | ronde 4 |
> | *"44 TERBALIK"* | **44** — cocok, tidak berubah | ronde 4 |
>
> Penyebabnya: ronde 2 menggabung kode 5 dan kode 6 menjadi satu ember dan memakai basis hitung
> yang berbeda. Nol kesimpulan turunan bersandar pada angka lama itu.

### Yang membuat modul ini berbeda dari Claim Life

⚠️ Jangan menyalin keputusan Claim Life mentah-mentah. Yang **berbeda dan terbukti**: invariant satu
mata uang **tidak berlaku**; `.Type` bernilai numerik, bukan `QR`/`QP`/`TP`/`TR`; `TransferType`
**hidup** dan menjadi dispatcher lintas-case; bentuk spreading berbeda; dan modul ini punya
Estimation, Loss Allocation, Insured Interest, Deductible, Catastrophe, serta Adjuster/Consultant yang
**tidak ada padanannya** di Life.

Yang **boleh** diwarisi hanya satu: rule pemilihan roster komite, yang isinya **identik** dengan
Claim Life.
