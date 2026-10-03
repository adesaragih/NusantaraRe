# Spec — Claim Fac In (migrasi Pega → Go + React + Oracle)

**Tanggal:** 2026-09-19 · **Lingkup:** Claim Fac In saja
**Sumber:** `grilling-ronde-1.md` · `grilling-ronde-1-ulang-docs.md` · `grilling-ronde-2.md` ·
`grilling-ronde-3.md` · `grilling-ronde-4.md` · `grilling-ronde-5.md` — di folder yang sama.
**Korpus:** `D:\XML\RNM_BRD\Claim Fac In\` — **482 berkas rule**, READ-ONLY.

---

## Cara membaca berkas ini

> ### Ringkasan — 2026-09-19
>
> ```
> user story            46
> acceptance criteria  114
> ── sebaran penanda atas 114 AC ──────────────────────────
> [penyimpangan sadar]  14    ⚠️ total            46
> jebakan (⚠️ tanpa penggolongan)                 32
> AC ber-[terbuka] aktif 20
> titik penyimpangan sadar (Bab 17)               12
> AC tanpa penanda        0
> butir register terdaftar 25 · tertutup 5 · memblokir 5
> ```
>
> **Dihitung dua cara, dan keduanya sepakat** *(aturan `CLAUDE.md` §4a)*:
> *(a)* pengurai blok per nomor yang **membuang baris kutipan `>` satu per satu**;
> *(b)* mesin-keadaan baris yang **memotong blok di baris kutipan atau sub-judul pertama** —
> aturan kutipan yang **berbeda**, bukan skrip yang sama dijalankan dua kali.
> **Jendela:** bab *Acceptance Criteria* berkas ini saja, dari judulnya sampai judul `##` berikutnya.
>
> ⛔ **RALAT 2026-09-19, dalam blok yang sama.** Angka yang pertama saya tulis adalah **perkiraan
> sebelum dihitung**, dan dikutip di sini tidak dihapus: *"[penyimpangan sadar] **25** · ⚠️ total
> **48** · jebakan **23** · AC ber-[terbuka] aktif **14**"*. ⭐ Sesudah kedua cara dijalankan:
> **14 · 44 · 30 · 23**. ⚠️ Dicatat karena inilah persis bentuk kesalahan yang melahirkan aturan
> §4a — **angka yang ditulis sebelum diukur**.
>
> ⛔ **RALAT 2026-09-19 — kedua.** Sesudah **butir 22 dijawab** dan **lima keputusan work owner
> dicatat**, angkanya bergeser lagi. Kalimat lamanya dikutip: *"acceptance criteria **112** ·
> [penyimpangan sadar] **14** · ⚠️ total **44** · jebakan **30** · AC ber-[terbuka] aktif **23**"*.
> ⭐ Yang berlaku: **114 · 14 · 46 · 32 · 20**.
>
> **Dihitung dua cara, dan keduanya sepakat.** **Delta yang menjelaskan selisihnya:** AC bertambah
> **2** *(113 ketergantungan lintas modul · 114 data lama yang kini menjadi galat)*; AC
> ber-`[terbuka]` `23 − 4 + 1 = 20` — **empat ditutup** *(AC 4 · 27 · 56 · 109)* dan **satu
> ditambah** *(AC 113)*; AC 114 sudah terhitung di antara keduanya.
>
> ⚠️ **Uji instrumen** *(butir yang jawabannya sudah diketahui sebelum dihitung)*: AC **1** tanpa ⚠️ ·
> AC **30** ber-`[penyimpangan sadar]` · AC **114** butir terakhir · AC **58** ber-`[terbuka]`.
> **Keempatnya terbaca benar oleh kedua cara.**
>
> ⭐ **Instrumennya terbukti bekerja — dan bukan karena lulus, melainkan karena ia GAGAL saat
> seharusnya gagal.** Pada jalankan pertama sesudah suntingan 2026-09-19, butir uji masih
> berbunyi *"AC **112** butir terakhir"* dan hasilnya **⛔ SALAH** — sebab dua AC baru sudah
> ditambahkan. ⚠️ **Uji yang tidak pernah gagal tidak membuktikan apa pun.**

**Tanda golongan.** Setiap butir dibuka tanda golongan dan — di bab *Acceptance Criteria* — ditutup
Bab asalnya.

| Tanda | Artinya |
| --- | --- |
| `[terverifikasi]` | perilaku **ditiru** dari Pega; buktinya ada di ronde yang disebut |
| `[keputusan work owner]` | perilaku **diputuskan**, bukan ditiru |
| `[keputusan work owner — atas rekomendasi asisten]` | pilihannya datang dari rekomendasi asisten; **pencabutannya murah** |
| `[data DBA]` | bersandar pada DDL atau isi basis data produksi |
| `[terbuka]` | **belum terjawab** — menyebut nomor butirnya di Bab 19 |
| `[penyimpangan sadar]` | **sengaja berbeda** dari Pega — daftar lengkapnya di **Bab 17** |
| ⚠️ | jebakan, kehati-hatian, atau penyimpangan |
| ⭐ | temuan yang mengubah pembacaan |

⛔ **Butir ber-`[terbuka]` bukan sasaran uji.** Ia menandai tempat yang **belum punya** sasaran uji,
supaya ketiadaannya terlihat dan tidak ditambal dengan tebakan.

**Aturan RALAT.** ⛔ **Tidak ada kalimat lama yang dihapus.** Setiap koreksi mengutip kalimat lamanya
utuh di sebelahnya. Enam berkas sumber penuh blok `> ⛔ RALAT`; **kalimat yang sudah diralat di sana
tidak dipakai sebagai bahan spec ini** — lihat Bab 20.

⚠️ **Bab *Acceptance Criteria* tidak memutuskan apa pun.** Ia menyatakan ulang, dalam bentuk yang
dapat diuji dari luar, apa yang sudah ada di Bab 1–18. **Bila sebuah butir terasa seperti keputusan
baru, ia salah tulis: laporkan, jangan laksanakan.**

---

## Problem Statement

Klaim **fakultatif masuk** *(facultative inward)* Nusantara Re hari ini berjalan di atas Pega —
**482 rule**, **179 activity**, **2.509 langkah**, dan **satu berkas alur** yang memegang seluruh
daur hidup kasus. Yang dirasakan orang yang bekerja dengannya:

1. **Aturan bisnisnya tidak terbaca dari layar.** Yang menentukan sebuah klaim boleh dikirim ke
   komite atau tidak adalah sebuah penanda bernama `IsError` yang bernilai **2**, dan arti angka
   itu **tidak tertulis sebagai aturan di mana pun** — ia hanya ada di **catatan bebas seorang
   pengembang** pada rule lain.
2. **Angka uang tidak dijaga tipenya.** Total penyesuaian dideklarasikan sebagai bilangan pecahan
   biner, tiga nilai uang lain lewat sebagai teks, dan satu di antaranya harus **ditambal
   koma-ke-titik** di dalam SQL sebelum dapat dihitung.
3. **Nama tidak dapat dipercaya.** Kolom hasil SQL dinamai hal yang bukan isinya — tanggal kejadian
   dinamai *kode cabang*, nomor klaim dinamai *kode bisnis* — dan **dua rule yang mirip berbohong
   dengan pemetaan yang berbeda**. Rule bernama `RDBList`, yang menurut namanya membaca, ternyata
   **menyimpan dan menutup transaksi**.
4. **Kegagalan tidak terlihat.** Dari **52 langkah** simpan/baca objek, hanya **dua** punya jalur
   kegagalan. Delapan rule menutup transaksi basis data dari dalam dirinya sendiri, tiga di
   antaranya **di tengah pekerjaan pemanggilnya**.
5. **Percabangan produk disembunyikan di dalam lompatan.** Modul ini bercabang per lini produk
   dengan **lompatan bernama label** — `MBU` · `PA` · `TRAVEL` · `MBD` — bukan dengan gerbang yang
   dapat dibaca berdampingan.
6. **Daur hidup kasus tinggal di ruleset bernama akun perorangan.** Satu-satunya berkas alur modul
   ini ada di `ADESAMUEL@`, terpisah dari tempat 467 rule lain tinggal.

---

## Solution

Membangun ulang Claim Fac In sebagai **Go + React + Oracle**, dengan perilaku Pega **ditiru apa
adanya** kecuali pada **dua belas titik yang sengaja diubah** *(Bab 17)*.

Yang berubah bagi orang yang memakainya:

- **Uang menjadi desimal di seluruh jalur** — tidak ada lagi bilangan pecahan biner dan tidak ada
  lagi nilai uang yang disimpan sebagai teks.
- **Aturan validasi menjadi galat yang terbaca**, bukan penanda bernilai angka yang artinya hanya
  diketahui dari catatan pengembang.
- **Kolom hasil dinamai sesuai isinya**, disertai tabel pemetaan nama-lama → nama-benar sehingga
  perbandingan keluaran lama dan baru tidak disalahbaca.
- **Satu aksi pengguna menjadi satu transaksi** — tidak ada lagi `COMMIT` yang menutup pekerjaan
  orang lain di tengah jalan.
- **Penyerahan ke komite gagal terang-terangan** bila daftar penyetujunya kosong, alih-alih
  melahirkan tangga nol tingkat yang senyap.
- **Klasifikasi lini dibaca sekali** saat kasus dimuat, bukan lewat 38 rule terpisah.

Yang **tidak** berubah: bentuk data tiga tingkat *(objek → item objek → penyesuaian)*, roster
penyetuju, penyimpanan berkas di Google Storage, dan penomoran lewat stored procedure.

---

## User Stories

### Registrasi dan polis

1. Sebagai **petugas klaim**, saya ingin mendaftarkan klaim fakultatif masuk terhadap sebuah polis,
   sehingga klaim tercatat pada polis yang benar.
2. Sebagai **petugas klaim**, saya ingin mencari polis dan melihat rinciannya sebelum mendaftarkan
   klaim, sehingga saya tidak salah polis.
3. Sebagai **petugas klaim**, saya ingin melihat riwayat klaim pada polis yang sama, sehingga saya
   tahu apakah klaim ini pernah diajukan.
4. Sebagai **petugas klaim**, saya ingin diperingatkan bila klaim dengan tanggal kejadian yang sama
   sudah ada pada polis itu, sehingga saya tidak membuat klaim ganda.
5. Sebagai **organisasi**, saya ingin nomor klaim dibuat oleh satu sumber tunggal, sehingga tidak
   ada dua klaim bernomor sama.
6. Sebagai **petugas klaim**, saya ingin mencatat penyebab kerugian dari daftar yang sudah baku,
   sehingga pelaporan dapat dikelompokkan.

### Lini produk

7. Sebagai **petugas klaim**, saya ingin sistem mengenali lini produk klaim ini — kendaraan
   bermotor, kecelakaan diri, perjalanan, atau lainnya — sehingga layar yang muncul sesuai isinya.
8. Sebagai **petugas klaim**, saya ingin rincian objek pertanggungan yang tampil sesuai lini
   produknya, sehingga saya tidak mengisi medan yang tidak relevan.
9. Sebagai **pengembang**, saya ingin klasifikasi lini ditentukan **satu kali di satu tempat**,
   sehingga tidak ada dua bagian sistem yang menggolongkan klaim yang sama secara berbeda.

### Objek pertanggungan dan itemnya

10. Sebagai **petugas klaim**, saya ingin mencatat beberapa objek pertanggungan dalam satu klaim,
    sehingga klaim atas polis berobjek banyak tetap utuh.
11. Sebagai **petugas klaim**, saya ingin tiap objek punya daftar itemnya sendiri, sehingga rincian
    per unit terjaga.
12. Sebagai **petugas klaim**, saya ingin nilai pertanggungan tiap objek tercatat lengkap dengan
    mata uang dan kursnya, sehingga klaim multi-mata-uang tidak tercampur.

### Estimasi

13. Sebagai **petugas klaim**, saya ingin memasukkan nilai estimasi kerugian per objek, sehingga
    cadangan dapat dibentuk.
14. Sebagai **petugas klaim**, saya ingin ditolak saat memasukkan nilai estimasi negatif, sehingga
    kesalahan ketik tidak masuk ke angka cadangan.
15. Sebagai **petugas klaim**, saya ingin ditolak saat tanggal estimasi tidak masuk akal, sehingga
    urutan waktu klaim tetap benar.
16. Sebagai **petugas klaim**, saya ingin diperingatkan bila estimasi melebihi nilai pertanggungan,
    sehingga saya memeriksanya sebelum melanjutkan.
17. Sebagai **petugas klaim**, saya ingin diperingatkan bila estimasi melampaui batas tanggung
    jawab, sehingga batas polis terjaga.
18. Sebagai **petugas klaim**, saya ingin tahu **dengan jelas** mana peringatan yang menghentikan
    saya dan mana yang hanya memberi tahu, sehingga saya tidak menebak.
19. Sebagai **petugas klaim**, saya ingin total estimasi terhitung ulang otomatis setiap kali satu
    baris berubah, sehingga angkanya selalu konsisten.

### Spreading

20. Sebagai **petugas klaim**, saya ingin kerugian dibagi ke para penanggung sesuai porsinya,
    sehingga beban tiap pihak benar.
21. Sebagai **petugas klaim**, saya ingin pembagian itu dipisahkan per mata uang, sehingga nilai
    tidak tercampur antar-mata-uang.
22. Sebagai **organisasi**, saya ingin jumlah porsi pada satu mata uang tepat seratus persen,
    sehingga tidak ada nilai yang hilang atau terhitung dua kali.

### Penyesuaian

23. Sebagai **petugas klaim**, saya ingin mencatat baris penyesuaian terhadap estimasi, sehingga
    perubahan nilai klaim terekam.
24. Sebagai **petugas klaim**, saya ingin mencatat nilai sisa barang selamat, sehingga klaim bersih
    terhitung benar.
25. Sebagai **petugas klaim**, saya ingin mencatat biaya adjuster, sehingga biaya penanganan ikut
    diperhitungkan.
26. Sebagai **petugas klaim**, saya ingin menunjuk adjuster dari daftar yang sudah ada, sehingga
    tidak ada nama yang diketik bebas.
27. Sebagai **petugas klaim**, saya ingin total penyesuaian terhitung dari baris-barisnya, sehingga
    tidak ada angka yang dimasukkan dua kali.

### Komite

28. Sebagai **petugas klaim**, saya ingin menyerahkan klaim ke komite ketika nilainya melewati
    kewenangan saya, sehingga keputusan diambil pihak yang berwenang.
29. Sebagai **petugas klaim**, saya ingin penyerahan **ditolak dengan pesan yang terlihat** bila
    daftar penyetujunya kosong, sehingga tidak ada kasus komite tanpa penyetuju.
30. Sebagai **penyetuju**, saya ingin menerima kasus komite lengkap dengan penunjuk ke baris
    penyesuaian yang dinilai, sehingga saya tahu persis apa yang saya putuskan.
31. Sebagai **petugas klaim**, saya ingin tahu apakah sebuah klaim sedang berada di komite, sehingga
    saya tidak mengubahnya saat sedang dinilai.

### Dokumen

32. Sebagai **petugas klaim**, saya ingin mencetak lembar muka klaim, sehingga berkas fisik lengkap.
33. Sebagai **petugas klaim**, saya ingin menerbitkan nota kerugian sementara dan tetap, sehingga
    pihak luar menerima dokumen resmi.
34. Sebagai **petugas klaim**, saya ingin dokumen hanya dapat diterbitkan bila datanya lolos
    validasi, sehingga tidak ada dokumen dengan angka yang salah.

### Berkas lampiran

35. Sebagai **petugas klaim**, saya ingin melampirkan berkas pada klaim, sehingga bukti tersimpan.
36. Sebagai **petugas klaim**, saya ingin tiap berkas punya pengenal yang **pasti unik**, sehingga
    dua unggahan tidak saling menimpa.
37. Sebagai **petugas klaim**, saya ingin membuka kembali berkas yang sudah diunggah, sehingga
    pemeriksaan ulang mungkin dilakukan.

### Pembayaran dan efek keluar

38. Sebagai **organisasi**, saya ingin data akseptasi klaim terkirim ke sistem Kasir, sehingga
    pembayaran dapat diproses.
39. Sebagai **organisasi**, saya ingin kegagalan pengiriman **tercatat dan terlihat**, sehingga
    tidak ada kiriman yang hilang diam-diam.
40. Sebagai **petugas klaim**, saya ingin menerima pemberitahuan lewat surel pada tahap yang
    ditentukan, sehingga saya tidak perlu memantau layar terus-menerus.

### Keutuhan data dan jejak

41. Sebagai **organisasi**, saya ingin satu aksi pengguna menghasilkan **satu transaksi**, sehingga
    tidak ada keadaan separuh tersimpan.
42. Sebagai **auditor**, saya ingin setiap perubahan penting terekam beserta pelaku dan waktunya,
    sehingga jejaknya dapat ditelusuri.
43. Sebagai **auditor**, saya ingin tahu angka mana yang disimpan dan angka mana yang dihitung
    ulang, sehingga selisih dapat diterangkan.

### Wewenang

44. Sebagai **organisasi**, saya ingin wewenang ditentukan **aturan peran yang tertulis**, bukan
    keadaan data yang kebetulan, sehingga hak akses dapat diaudit.
45. Sebagai **pengguna**, saya ingin tidak melihat tombol yang memang bukan wewenang saya, sehingga
    saya tidak mencoba hal yang akan ditolak.

### Migrasi

46. Sebagai **organisasi**, saya ingin data klaim lama terbawa beserta statusnya, sehingga tidak ada
    klaim berjalan yang tertinggal di sistem lama.

---

## Implementation Decisions

### Bab 1 — Lingkup, kelas kerja, dan identitas rule

`[terverifikasi]` *(ronde 1 §B2)* Objek kerja modul ini berkelas **`ASM-FW-GCNMFW-WORK-PNC`**
*(89 rule)*. ⛔ **Bukan** `WORK-CLAIMTREATY` seperti Claim Prop.

⭐ **Dua kelas yang tidak punya padanan di Claim Prop:** **`DATA-OBJECT`** *(50 rule)* dan
**`DATA-OBJECTITEM`** *(56 rule)*. Di Claim Prop konsep itu hidup sebagai halaman di dalam kasus; di
sini ia **kelas rule tersendiri dengan 106 rule**.

⚠️ **Akibatnya tegas, dan ia membatalkan satu pinjaman:** penggolongan Claim Prop yang menyatakan
`ObjectList`/`ObjectItemList` **bukan tabel** ⛔ **tidak berlaku di sini**. `[terverifikasi]`
*(ronde ulang §F4)* keduanya **ditulis oleh puluhan rule** — `.ObjectItemList` **19 berkas penulis /
54 pembaca**, `ClaimData.ObjectList` **26 penulis / 57 pembaca** — jadi keduanya **memegang data**.

#### Identitas rule EMPAT bagian

⭐ `[terverifikasi]` *(ronde ulang §A1)* **Identitas rule = `pxInsName` + kelas + `pyRuleSet` +
`pyRuleSetVersion`.** `pxInsName` saja **tidak cukup**: dua rule bernama sama terbukti hidup di
ruleset berbeda — `BROWSECURRENCY_RD` dan `CURRENCYSTANDARD`, satu di `GISFW`, satu di `GCNMFW`.

Dibaca dengan identitas empat bagian: rule bersama dengan Claim Prop **161** *(bukan 162)*, waktu
perubahan identik **159** *(bukan 158)*, beda waktu **2**, dan **4 pasangan** bernama sama tetapi
beda ruleset/versi.

#### Ruleset `ADESAMUEL@`

⭐ `[terverifikasi]` *(ronde 2 §G1)* **Lima rule tinggal di ruleset bernama akun perorangan** —
dan salah satunya **`Register_Flow`, satu-satunya berkas alur modul ini**.

✅ `[keputusan work owner]` **2026-09-19** — *"ikuti aja begitu, itu yang akan dipake."*
⭐ **`Register_Flow` di `ADESAMUEL@` versi `01-01-01` ADALAH daur hidup yang berlaku**, dan itulah
yang dipindahkan. ⛔ Tidak dicari penggantinya di ruleset resmi, dan ⛔ tidak diperlakukan sebagai
ekspor yang kurang lengkap.

⚠️ **Fakta yang tetap dicatat:** siapa pun yang kelak mencari `Register_Flow` di ruleset resmi
**tidak akan menemukannya** — dan itu **bukan tanda ekspor rusak**.

#### Tiga sistem sumber

⚠️ `[terverifikasi]` *(ronde ulang §A2)* Korpus mencampur ekspor dari **tiga** sistem:
`pega` **320** · `pegadevnusare2` **146** · `pegaprdnusare` **16**. Pola yang **sama** ada di Claim
Prop. ⭐ Dua salinan rule beridentitas empat bagian **sama** terbukti **berbeda isinya**
*(`GENERATEIMAGEID_SQL`)*.

✅ `[keputusan work owner — atas rekomendasi asisten]` **2026-09-19** — ⭐ **rule yang ditiru
adalah salinan dari sistem PRODUKSI.**

⭐ **Satu akibat langsung yang sudah terbaca:** pembuat pengenal berkas dan penentu jenis berkas di
modul ini **berasal dari salinan produksi** — maka **versi berketelitian nanodetik disertai nilai
unik sejagat** dan **19 jenis berkas** itulah yang berlaku. ✅ Menguatkan AC 54, dan ⭐ **menutup
butir 9**: kedua jenis berkas video **ikut**, karena keduanya ada di salinan produksi.

⚠️ **Akibat yang menyentuh modul lain, dan wajib disebut:** pola tiga sistem itu **ada juga di
Claim Prop**, dan seluruh berkas keluaran yang sudah jadi — spec, struktur tabel, relasi tabel —
**disusun tanpa pernah menyatakan dari sistem mana rule-nya berasal**. ⛔ **Modul lain tidak
disentuh dari sini**; dicatat sebagai pekerjaan yang lahir dari jawaban ini.

### Bab 2 — Daur hidup kasus

`[terverifikasi]` *(ronde 5 §E)* Berkas alur dibaca dari wadah **`pyShapes`**, dan isinya:

| | Jumlah |
| --- | ---: |
| **bentuk** | **8** — `Start1` · `Assignment1` · `Assignment3` · `Assignment7` · `Decision4` · `Decision5` · `Decision8` · `END52` |
| **penghubung** | **10** |
| **pengubah** *(tiket tingkat alur)* | **2** |
| **tiket** *(menggantung pada bentuk)* | **7** |

⭐ **Bentuk selesai ADA** — `END52`, berstatus akhir **`Resolved-Completed`**.

⚠️ `[terverifikasi]` **Pengenal `END52` sama dengan Komite Claim Prop, tetapi kelasnya berbeda**
*(`Work-PNC` lawan `Work-KomiteTreaty`)*. ⭐ Itu **warisan salinan** — alur disalin lalu diubah
kelasnya — **bukan rule bersama dan bukan kebetulan**.

#### Tiket: dua bernama, lima cangkang kosong

`[terverifikasi]` *(ronde 5 §E1)*

| Menggantung di | Nama |
| --- | --- |
| `Assignment7` | **`SendToEstAdmin`** |
| `Assignment3` | **`AcceptanceClaim`** |
| `Assignment1` · `Decision4` · `Decision5` · `Decision8` · `END52` | ⛔ **tanpa nama — cangkang kosong** |

⭐⭐ **Ketujuhnya LABEL MATI.** `[terverifikasi]` Sapuan **seluruh 20 korpus**: metode
`Obj-Set-Tickets` **NOL**, activity `SetTicket` **NOL** di modul ini. Tidak ada yang melemparnya.

⚠️ **Batas kejujuran:** mesin pembangkitnya **ada di aplikasi yang sama** — `SetTicket` hidup di
`NB FacIn`, `RNW Fac In`, `Endorsment Fac In`, dan ia rule `@baseclass` bawaan Pega.
⛔ Tidak ada buktinya menyentuh modul ini. Bab 19 butir 11.

### Bab 3 — Percabangan per lini produk

⭐ `[terverifikasi]` *(ronde 1 §C1 · ronde ulang §F2)* **Modul ini bercabang per lini produk dengan
LOMPATAN, bukan dengan gerbang.** **23 lompatan** di seluruh modul, dan **sembilan** di antaranya di
satu rule — `GetAllData_Act`, dengan label `MBU` · `MBU2` · `MBU3` · `PA` · `PA2` · `PA3` ·
`TRAVEL` · `TRAVEL2` · `TRAVEL3`. Lini keempat, **`MBD`**, muncul di `GettsiAneka_Act`.

⭐ **Claim Prop punya NOL lompatan.** Bentuk ini **tidak ada padanannya** di sana.

⚠️ `[terverifikasi]` **Pasangan gerbang→tanda tertukar** di dua blok bersarang `GetAllData_Act`:
`IsMBU` melompat ke tanda Travel dan `IsTravel` ke tanda MBU.

> ⛔ **DICABUT** `[keputusan work owner]` 2026-09-19 — *"ikuti apa adanya, jangan jadikan
> permasalahan."* Kalimat keputusan lamanya dikutip, tidak dihapus: *"sepertinya itu kesalahan
> developer sebelumnya, perbaiki logic nya di go."*
>
> ⭐ **Sebabnya bersandar pada bukti:** `[terverifikasi]` *(ronde 2 §B4)* `GetAllData_Act` punya
> **61 penugasan properti** dan ⭐ **NOL menuju `pyWorkPage`** — seluruhnya ke `Local`,
> `OutputData.pxResults`, `MataUang`, `TempOutstanding`, `InputData`, `DataView`. Cacat itu
> **hanya menyentuh tampilan satu popup**, bukan data tersimpan.

`[terverifikasi]` **`TempOutstanding` tidak disimpan rule lain mana pun** *(ronde 4 §E2)*.

**Satu lompatan menggantung** `[terverifikasi]`: `SetProtectionEstimation` langkah 26 menunjuk label
yang tidak ada, **dan gerbangnya mati**. **Satu lompatan kode 4** *(keluar-iterasi)*:
`ProtectionObjectItem_Act` langkah 3, **gerbangnya hidup** — bila syaratnya tidak terpenuhi,
**perulangan diputus**, bukan langkahnya dilewati.

### Bab 4 — Klasifikasi lini dan `Quotation`

✅ `[keputusan work owner]` **2026-09-19** — *"ada, tapi saat migrasi quotation langsung select dari
table polis tidak di pyWorkPage.Quotation"*.

**Dua hal sekaligus:**

1. ✅ **`pyWorkPage.Quotation` ADA dan propertinya BENAR-BENAR DIISI** — ⭐ **kebalikan Claim
   Prop**. ⭐⭐ **Buktinya BUKAN yang semula ditulis di sini** — lihat RALAT di bawah. Rantai
   pengisiannya, dibaca ujung ke ujung:

   | # | Mata rantai | `[terverifikasi]` |
   | --- | --- | --- |
   | 1 | ⛔ **Di dalam 482 berkas Claim Fac In, `BusinessType` · `BusinessCode` · `BusinessName` · `StatusBusiness` punya NOL penulis** | disisir **dua cara**: *(a)* `PropertiesName` di seluruh 482 berkas — **0**; *(b)* **setiap tag** yang memuat keempat kata itu — **675 kemunculan, nol di antaranya sebagai sasaran penugasan** |
   | 2 | ⭐ **`pyWorkPage.Quotation` diisi dengan PENYALINAN HALAMAN UTUH** dari `pyWorkPage.OfferFacIn.QuotationData` | dua activity melakukannya, keduanya tanpa gerbang |
   | 3 | ⛔ **`OfferFacIn.QuotationData` sendiri tidak pernah ditulis sebagai halaman di dalam 482 berkas** — ia hanya **disalin DARI** | jadi ia **datang dari luar ekspor ini** |
   | 4 | ⭐⭐ **Modul saudara yang mengisinya** — penawaran fakultatif dibuat, diperbarui, dan diperbarui-ulang oleh modul lain, dan di sanalah keempat properti itu **ditulis** | ⚠️ jendela di luar 482: **NB FacIn 2.083 berkas · RNW Fac In 1.927 · Endorsment Fac In 2.061**; modul terakhir menulis `OfferFacIn.QuotationData.BusinessType`, `…BusinessName`, `…BusinessCode` — **halaman yang persis disalin modul ini** |

   ⭐ **Karena itu ke-60 rule `When` DAPAT bernilai benar**, dan **nol** di antaranya boleh disebut
   mati. ✅ **Vonis butir 22: BERDIRI.**

   ⚠️ **Tetapi ia melahirkan ketergantungan yang belum pernah tercatat:** ketepatan klasifikasi
   lini modul ini **bergantung pada modul lain yang mengisi halaman penawaran**. ⛔ Bila penawaran
   dibuat tanpa keempat properti itu, rule-nya **diam-diam bernilai salah** — dan **tidak ada
   satu pun pemeriksaan** di modul ini yang akan mengatakannya.

   > ⛔ **RALAT 2026-09-19.** Kalimat lamanya **dikutip utuh, tidak dihapus**:
   >
   > > *"✅ **`pyWorkPage.Quotation` ADA** pada objek kerja Claim Fac In — ⭐ **kebalikan Claim
   > > Prop**. `[terverifikasi]` *(ronde ulang §A3)* **kedua jalur halaman punya penulis**: jalur-1
   > > `pyWorkPage.Quotation.*` ditulis **2 berkas / 10 penugasan**; jalur-2
   > > `pyWorkPage.OfferFacIn.QuotationData.*` ditulis **4 berkas / 27 penugasan**."*
   >
   > ⚠️ **Angkanya benar; yang salah adalah apa yang ia buktikan.** Ke-10 dan ke-27 penugasan itu
   > menulis **properti LAIN** — `CedingCo`, `SobName`, `SourceOfBusiness`, `QQName`,
   > `TotalShareCeding`, dan kerabatnya — ⛔ **tidak satu pun menulis keempat properti yang diuji
   > rule klasifikasi**. Kalimat lama membuktikan **halamannya** ditulis, lalu **diam-diam
   > dipakai** untuk menyimpulkan **propertinya** terisi. ⭐ Ronde ulang §H2 sendiri sudah menunjuk
   > lompatan itu sebagai **kesimpulannya yang paling rawan** — dan ia memang keliru.
   >
   > ✅ **Vonisnya tidak berubah**, tetapi sekarang berdiri di atas rantai yang benar.
2. ⭐ **Di sistem baru `Quotation` diambil LANGSUNG dari tabel polis**, bukan dari halaman.
   ⚠️ **Penyimpangan sadar** — Bab 17 titik 1.

⭐ `[keputusan work owner — atas rekomendasi asisten]` **`Quotation` dibaca SEKALI saat kasus
dimuat** sebagai satu nilai turunan, dan **ke-38 rule `When` diganti SATU fungsi klasifikasi lini**.

**Alasannya:** ke-38 rule itu menguji **satu halaman yang sama** untuk memilah lini; memindahkannya
satu per satu berarti membawa **38 potong logika yang masing-masing berisi satu keputusan**.

✅ `[keputusan work owner]` **2026-09-19** — *"MASIH PROSES DEVELOP, (T_QUOTATIONDATA)"*.

⭐ **Nama tabelnya `T_QUOTATIONDATA`.** ⚠️ **Tetapi ia MASIH DIBANGUN**, dan `[terverifikasi]`
**tidak ada di korpus mana pun** — **0 berkas** dari 482 menyebutnya. ⭐ Ia **tabel baru**, bukan
tabel Pega yang sudah berjalan.

⛔ **Butir 1 MENYEMPIT, tidak tertutup.** Yang tertutup: **nama tabelnya**. Yang **masih terbuka**:
**nama kolomnya**, dan **kapan tabel itu siap**. ⚠️ **AC 14 tetap tanpa sasaran uji yang konkret**
sampai kolomnya ada. Bab 19 butir 1.

#### Rule dipisahkan dari kehidupannya

⭐ `[keputusan work owner — atas rekomendasi asisten]` Ke-60 rule `When` **identik pada keempat
bagian** dengan Claim Prop *(ronde ulang §D3 — beda: 0)*. **Rule-nya satu, dipindahkan sekali,
disimpan satu salinan.** Yang berbeda per modul **bukan rulenya**, melainkan **apakah halaman yang
diujinya ada pada objek kerja modul itu**.

✅ Dengan begitu keputusan Claim Prop *(di sana halamannya tidak ada, ke-52 tidak pernah bernilai
benar)* dan keadaan Fac In *(di sini halamannya ada, rule hidup)* **berdiri bersama tanpa saling
membatalkan**.

⚠️ **Bila fungsi klasifikasi tunggal di atas dibangun, butir ini selesai sendiri** — satu fungsi
meniadakan pertanyaan *"rule mana hidup di modul mana"*.

### Bab 5 — Estimasi dan validasinya

`[terverifikasi]` *(ronde 5 §C)* Rule validasi estimasi punya **60 langkah** dan **tujuh pesan**,
ditambah satu pembersih pesan di langkah 1. **Yang menentukan penghalangnya adalah langkah 12.7**,
yang menyetel **`pyWorkPage.IsError = 2`** — dan ⭐ **arti angka 2 hanya terbaca dari catatan langkah
pengembang pada rule lain**: *"set pyWorkPage.IsError ==2 **to disable send komite**"*.

`[terverifikasi]` **`IsError` DIBACA** di **tiga gerbang** — dua menguji `IsError>1`, satu menguji
`IsError<2`. ⛔ **Kebalikan preseden Claim Prop**, tempat `IsError` tidak pernah diset.

| Pesan | Syarat | Penanda yang menyusul | Akibat yang terbaca |
| --- | --- | --- | --- |
| lgk 7 | persentase estimasi **< 0** | ⛔ tidak ada | peringatan |
| lgk 8 | nilai estimasi **< 0** | syaratnya ikut di gerbang 23 | ⚠️ **belum punya data** |
| ⭐ lgk 12.2 | tanggal estimasi ditolak | ⭐ **12.7 → `IsError = 2`** | ⭐ **MENGHALANGI** |
| ⭐ lgk 12.5 | tanggal estimasi ditolak | ⭐ **12.7 → `IsError = 2`** | ⭐ **MENGHALANGI** |
| lgk 12.6 | persentase **dan** nilai keduanya **0** | 12.8 → `Local.ErrSts = 1` | **tombol CFS saja** |
| lgk 21 | batas tanggung jawab kosong/nol dan terlampaui | ⛔ tidak ada | peringatan |
| lgk 22 | estimasi melebihi nilai pertanggungan | 23 → `Local.ErrSts = 1` | **tombol CFS saja** |

⭐ **`Local.ErrSts` hanya dibaca di langkah 24**, yang menyetel penanda **CFS** dan
**PrintFaceClaim** — ✅ menguatkan bahwa kedua titik kode-5 di rule ini **soal tombol unduh, bukan
uang**.

✅ `[keputusan work owner — atas rekomendasi asisten]` **2026-09-19** — **ketujuhnya menjadi galat
validasi yang tegas**, dengan pembedaan **yang terbaca** antara yang menggagalkan penyimpanan dan
yang hanya memperingatkan.

**Alasannya:** bentuk sekarang **menyembunyikan aturan bisnis di dalam nilai sebuah penanda**, dan
aturan yang hanya hidup di dalam catatan bebas **tidak dapat diuji** dan **hilang begitu catatannya
hilang**. ⚠️ **Penyimpangan sadar** — Bab 17 titik 12.

✅ `[keputusan work owner — atas rekomendasi asisten]` **2026-09-19** — ⭐ **KETUJUHNYA
MENGGAGALKAN PENYIMPANAN.** ⛔ Tidak ada yang berstatus *"peringatan saja"*.

⚠️ **Alasannya:** ketujuhnya menguji **angka uang atau tanggal** — estimasi negatif, persentase
negatif, tanggal tidak masuk akal, nilai dan persentase keduanya nol, melebihi nilai pertanggungan,
melampaui batas tanggung jawab. ⛔ **Tak satu pun masuk akal dibiarkan tersimpan.**

✅ **Butir 12 DITUTUP.** ⚠️ **Tetapi ia melahirkan satu butir BARU, dan butir itu menyentuh data
lama:** di sistem lama hanya **dua** yang menghalangi. ⭐ **Karena itu data klaim lama MUNGKIN sudah
memuat keadaan yang kini menjadi galat** — estimasi bernilai-dan-berpersentase nol, atau estimasi
melebihi nilai pertanggungan. ⛔ **Berapa banyak belum diperiksa.** ⚠️ Itu **memblokir migrasi**,
bukan pembangunan. Bab 19 butir **25**.

⛔ **Yang tetap terbuka:** **apakah pesan Pega sendiri menghalangi penyerahan** tidak terbaca dari
ekspor — rule ini dipanggil sebagai **aksi sisi-klien saat pengguna mengetik**
*(`pyActionSets`/`pyBehaviors`, peristiwa `change` 81 · `click` 36)*, bukan sebagai rule validasi
penyerahan. Bab 19 butir 13. ⚠️ **Ia tidak lagi memblokir** — keputusan di atas menetapkan
perilaku sistem baru tanpa menunggu jawabannya.

### Bab 6 — Uang dan ketelitiannya

⛔ `[terverifikasi]` **Uang di modul ini tidak dijaga tipenya**, dan buktinya **berkaki dua**
*(ronde 5 §A4)*:

1. ⭐ **`SetKomiteList_ACT.TotalAdjustment` bertipe `Double`** — bilangan pecahan biner untuk total
   uang. ⭐ **Kedua sumber tipe sepakat** *(tanda tangan dan medan parameter)*.
2. ⭐ **Tiga nilai uang lewat sebagai teks** — `AdjustmentValue` · `Share` · `Estimation` — dan satu
   di antaranya **terbukti perlu ditambal koma-ke-titik** di dalam SQL sebelum dapat dihitung.

`[terverifikasi]` Sensus penuh tipe parameter, **jendela seluruh 482 berkas**:

```
sumber A  pyXMLSignature               169 parameter di  72 berkas
sumber B  pyParameters < pagedata      265 parameter di 116 berkas
beririsan                              169
⭐ TIPE BERSELISIH                        0
hanya di A                               0      <-- B superset murni
hanya di B                              96      <-- seluruhnya rule DataPage
```

⭐ **Tidak ada ketidakpastian tipe.** `TotalAdjustment` memang `Double`, `Share` memang teks, `TSI`
memang desimal.

✅ **Di sistem baru, uang bertipe desimal di seluruh jalur** — ⚠️ **penyimpangan sadar**, Bab 17
titik 7 *(pertentangan ADR-0003)*.

#### ✅ Ketelitian yang dipakai

`[keputusan work owner — atas rekomendasi asisten]` **2026-09-19** — ⭐ **uang dihitung sampai
20 angka dengan 8 di belakang koma, tanpa pembulatan di tengah jalan; satu-satunya titik pembulatan
ada di batas penyimpanan; tampilan layar 4 angka di belakang koma.**

⚠️ **Alasannya:** modul ini dan Claim Prop **berbagi tabel akseptasi yang sama**, dan
`[terverifikasi]` tabel itu **disentuh 27 kali** di modul ini — jauh mendominasi seluruh objek
Oracle yang disentuh. ⛔ **Ketelitian berbeda pada satu tabel bersama akan melahirkan selisih yang
tidak dapat diterangkan**, dan selisih itu akan tampak seperti cacat migrasi padahal bukan.

✅ Sejalan **ADR-0003** dan aturan ketelitian yang sudah dipakai `STRUKTUR-TABEL-CLAIM-PROP`.
✅ **Butir 14 DITUTUP.**

#### Muara `Local.TotalAdjustment`

⚠️ **Jebakan nama:** `Local.TotalAdjustment` *(properti pada halaman `Local`)* **bukan** parameter
`TotalAdjustment`. Keduanya bernama sama.

`[terverifikasi]` *(ronde 5 §B)* Nama itu muncul di **delapan berkas**: **7 penulisan, seluruhnya ke
halaman `Local`** · **1 pembacaan** ke `Local.IsADj` · **1 gerbang** terhadap batas bawah di
`SendPICProtect_Act` · dan **dua Section** yang hanya memuat **slot parameter kosong**.

⭐ **Tidak pernah sampai ke kolom tersimpan.** `CreateKMTNo_Act` punya **71 penugasan properti**,
**8 menuju `pyWorkPage`**, dan **nol** berasal darinya. `Local.IsADj` **ditulis dua kali, dibaca nol
kali** di seluruh 482 berkas.

### Bab 7 — Spreading

`[terverifikasi]` `.SpreadingList` adalah halaman **kedua paling sering diulang** *(31 kali)*.
Pembagian dipisahkan **per mata uang**, dan jumlah porsi per mata uang diperiksa terhadap seratus
persen.

⚠️ `[terverifikasi]` *(ronde 3 §G2)* **`SpreadingCheckSP` langkah 3.2 memuat kode arah 5**, yang
**mengubah rantai gerbang dari DAN menjadi ATAU** — dan ⭐ ia **menyentuh spreading**. Siapa pun
yang membacanya sebagai DAN akan salah.

### Bab 8 — Penyesuaian, objek, dan item objek

`[terverifikasi]` Unit data modul ini **tiga tingkat**: objek pertanggungan → item objek → baris
penyesuaian. ⭐ **Dua tingkat lebih dalam daripada Claim Life**, yang unit statusnya hanya baris
penyesuaian.

`[terverifikasi]` *(ronde 2 §C2)* `CreateKMTNo_Act` langkah **6.9** membawa ⭐ **LIMA penunjuk
posisional** ke kasus komite — `Adjustment.IDObject` · `Adjustment.CoverageSubscript` ·
`IndexObjectItem` · `IndexObject` · **`IndexAdjustment`** — dan **tidak satu pun bergerbang**.

⚠️ **Dua kejanggalan:** `IndexObject` diisi dari `.ObjectIndex` sedangkan `Adjustment.IDObject`
diisi dari `.IndexObject` — dua nama sumber berbeda untuk dua penunjuk bernama nyaris sama
*(Bab 19 butir 4)*; dan `local.IndexAdjust` ditulis huruf kecil sedangkan empat lainnya tidak —
petunjuk bahwa baris itu **ditulis belakangan**.

⛔ **Nol pemeriksaan** bila yang ditunjuk hilang atau bergeser — sama seperti Claim Prop, tetapi
**dengan tiga tingkat**, sehingga peluang bergesernya lebih besar.

### Bab 9 — Penyerahan ke komite

`[terverifikasi]` *(ronde 2 §C)* `CreateKMTNo_Act` punya **42 langkah**, **nol ber-remark**, **nol
lompatan**. ⭐ **Satu-satunya penjaga di seluruh rule** adalah langkah **10**: `Call pxAddChildWork`
bergerbang *"ada pesan kesalahan pada kasus"* → keluar-activity.

`[terverifikasi]` Roster penyetuju dibangun `SetListKomite_act`: `KomiteID` dari akun operator,
`IDKomite` dari jabatan, `KomiteAproval` mulai `0`. ⭐ **Bentuknya identik dengan Life dan Prop.**

⚠️ **Langkah 9 tanpa gerbang:** bila daftar penyetuju kosong, jumlah tangga menjadi **0** sementara
pencacah tetap **1** — **tangga nol tingkat, senyap**.

✅ `[keputusan work owner — atas rekomendasi asisten]` **Penyerahan itu DITOLAK dengan galat yang
terlihat.** ⚠️ **Penyimpangan sadar** — Bab 17 titik 2. ✅ Sejalan preseden **Komite Claim Life AC 4**.

#### `.KomiteNo`

`[terverifikasi]` Pega menulisnya **dua kali**: langkah 13.1 dari pengenal kasus komite, lalu
13.1.1 dari **potongan teks mulai huruf ke-19** atas kunci internal Pega.

✅ `[keputusan work owner — atas rekomendasi asisten]` **Dipakai yang pertama** — pengenal kasus
komitenya. **Alasannya:** angka 19 hanyalah panjang awalan kelas; begitu nama kelas berubah satu
huruf, potongan itu **diam-diam salah**. ⚠️ **Penyimpangan sadar** — Bab 17 titik 3.

⛔ **Yang belum ditemukan:** kenapa 13.1 ditimpa 13.1.1. Bab 19 butir 2.

#### `FlagOnGoingCommitte`

`[terverifikasi]` Pega mengisinya dengan **teks `"Send Commite"`** di langkah 3, tanpa gerbang.

✅ `[keputusan work owner — atas rekomendasi asisten]` **Kolom itu tidak dibuat** — keadaan *"komite
sedang berjalan"* **diturunkan dari kasus komitenya**. ⚠️ **Penyimpangan sadar** — Bab 17 titik 4.
✅ Sejalan **keputusan 27 Komite Claim Prop**.

⛔ **Belum dibuktikan** kolom Fac In dan kolom Claim Prop adalah kolom yang sama — **nama sama bukan
bukti**. Bab 19 butir 3.

### Bab 10 — Dokumen

`[terverifikasi]` Modul ini menerbitkan **lembar muka klaim**, **nota kerugian sementara**, dan
**nota kerugian tetap**. ⭐ Enam rule dokumen, dan penomorannya memakai **stored procedure pembangkit
nomor urut** yang sama dengan modul lain.

⚠️ `[terverifikasi]` Penerbitan lembar muka digerbangi penanda **CFS** dan **PrintFaceClaim**, yang
disetel rule validasi estimasi — lihat Bab 5.

### Bab 11 — Penyimpanan berkas

✅ `[terverifikasi]` Berkas **tetap disimpan di Google Storage**; token diterbitkan **Oracle**, bukan
aplikasi.

⭐ `[terverifikasi]` *(ronde ulang §B7)* Salinan pembuat pengenal berkas di modul ini **lebih baru**
daripada salinan Claim Prop, dan bedanya **bermakna**: ketelitian cap waktu naik dari **milidetik ke
nanodetik**, ditambah **nilai unik sejagat** dari basis data. ⭐ **Itu bentuk sebuah tambalan
tabrakan** — versi lama dapat menghasilkan **pengenal kembar** bila dua unggahan jatuh dalam
milidetik yang sama.

✅ `[keputusan work owner]` **2026-09-19** *(A4-1)* — **versi Claim Fac In yang berlaku** untuk
keempat rule bersama yang berbeda.

⭐ `[terverifikasi]` Salinan penentu jenis berkas modul ini mengenali **19 jenis**, dua lebih banyak
daripada salinan Claim Prop — tambahannya **`avi`** dan **`mp4`**.
⛔ **Apakah keduanya ikut ditiru belum diputuskan.** Bab 19 butir 9.

### Bab 12 — Efek keluar

`[terverifikasi]` *(ronde ulang §B6)* Sensus efek keluar, **jendela 179 Activity + 63 RDBList +
8 ConnectREST**:

| Jalur | Langkah | Berkas |
| --- | ---: | ---: |
| basis data — `RDB-List` | **118** | 69 |
| basis data — `Obj-Save` / `Commit` | **31** | — |
| layanan luar — `Connect-REST` | **11** | 9 |
| surel | **7** | 7 |
| Google Storage | — | 5 |
| Kasir | — | 9 |
| dokumen | — | 6 |

> ## ⚠️ ANGKA YANG PALING KERAS
>
> ⭐ **Dari 52 langkah simpan/baca objek, hanya DUA punya jalur kegagalan.** Enam lainnya ber-remark,
> jadi mati.

⛔ `[terverifikasi]` **Nol jejak outbox, nol antre-ulang, dan nol penanganan kegagalan** pada hampir
seluruh efek keluar. ⚠️ **Penyimpangan sadar** — Bab 17 titik 10 dan 11 *(pertentangan ADR-0008 dan
ADR-0015)*.

#### Alamat layanan keluar

`[terverifikasi]` **Delapan `ConnectREST`**, **tujuh** mengambil alamat dari tabel tautan layanan,
⚠️ **satu memakai URL langsung** — `getPremiumPaidOnMarine`. ⚠️ **Penyimpangan sadar** — Bab 17
titik 9 *(pertentangan ADR-0013)*.

⚠️ **Satu berautentikasi** — rule Kasir, dan ia **rule bersama** dengan Claim Prop.
⛔ **Nilai kredensialnya tidak disalin ke berkas mana pun.** ⚠️ Karena kredensial itu beredar di
dalam ekspor, ia **sebaiknya diganti sesudah migrasi**. Itu urusan tim pemilik layanan.

#### `IsPEGAPROD`

`[terverifikasi]` Rule itu **hidup** dan menggerbangi **10 perujuk**, sebagian besar efek keluar.
Kondisinya menguji **tingkat produksi proses berjalan**.

⚠️ **Penyimpangan sadar** — Bab 17 titik 8 *(pertentangan ADR-0005)*: di sistem baru ia menjadi
**flag lingkungan eksplisit**.

⭐ `[terverifikasi]` **`IsPEGAPROD` tidak menjelaskan campuran tiga sistem di Bab 1** — yang satu
mengukur perilaku saat berjalan, yang satu jejak penyuntingan. Keduanya **bukan pengukur yang sama**.

#### Arasapas

⚠️ `[keputusan work owner]` **2026-09-19** — *"SAYA CEK DI FACIN TIDAK MATI"*.
⭐ `[terverifikasi]` **Panggilan konversi TIDAK MATI.** Gerbang berflag `false` mematikan
**GERBANGNYA**, bukan **LANGKAHNYA** — sehingga panggilan itu justru berjalan **tanpa saringan**,
dari **empat** pemanggil, **nol bergerbang hidup**.
⛔ Apakah itu memang dikehendaki **belum diputuskan**. Bab 19 butir 7.

### Bab 13 — Transaksi dan `COMMIT`

⭐ `[terverifikasi]` *(ronde 5 §D)* **Tiga belas rule menutup transaksi sendiri, dan ia DUA jenis
yang sangat berbeda:**

| Jenis | Berapa | Artinya |
| --- | ---: | --- |
| metode langkah **`Commit`** Pega | **5** *(Activity)* | commit milik Pega sendiri — menutup pekerjaannya dan melepas kunci |
| ⭐⭐ **`COMMIT;` di dalam SQL** | **8** *(RDBList)* | ⭐ **commit tingkat basis data** — menutup transaksi sesi, **termasuk pekerjaan pemanggilnya yang belum selesai** |

⭐ **Tiga panggilan jatuh di tengah pekerjaan yang lebih besar:** rule penyimpan akseptasi dipanggil
`CloseClaim` dan `SaveCFS_ACT`; rule penyisip klaim dipanggil `SaveAdjustmenttoDB_ACT`.
⭐ **Dan satu yang menenangkan:** di `SaveAcceptation` panggilan itu **ber-remark** — mati —
sedangkan `Obj-Save` di langkah berikutnya **hidup**.

✅ `[keputusan work owner — atas rekomendasi asisten]` **2026-09-19** — **Satu aksi pengguna, satu
transaksi.** ⛔ Kedelapan `COMMIT;` **tidak direplikasi**.

⚠️ **Satu pengecualian: pembangkit nomor urut.** Nomornya **harus bertahan walau transaksi induknya
batal** — kalau ikut dibatalkan, nomor yang sama **terpakai ulang**. ✅ Sejalan **ADR-0006**.

**Alasannya:** `COMMIT;` di tengah pekerjaan orang lain membuat kegagalan separuh jalan meninggalkan
**data separuh tersimpan**, dan ⛔ **pemanggilnya tidak dapat memulihkannya** — transaksinya sudah
ditutup oleh rule yang ia panggil.

⚠️ **Perlu diingat saat membangun:** rule-rule ini bernama **`RDBList`** — yang menurut namanya
membaca — padahal **menyimpan dan menutup transaksi**.

### Bab 14 — Jejak audit

✅ `[terverifikasi]` Jejak kronologi **ditulis** oleh rule khusus. ⛔ **Kelengkapannya belum diuji.**

⚠️ `[terverifikasi]` *(ronde ulang §D4)* Dua catatan pengembang berbunyi *"hapus yg set
chronologi"*; salah satunya **masih ada jejaknya, hidup**. ⛔ Apakah penghapusan itu dikehendaki
**tidak terbaca**.

### Bab 15 — Wewenang

⭐⭐ `[terverifikasi]` *(ronde 3 §D3)* **Claim Fac In tidak punya satu pun gerbang layar berdasarkan
wewenang.** Sensus **seluruh medan berbau hak akses** di 106 berkas layar — **15 nama medan
berbeda**:

| Medan | Hadir | Terisi |
| --- | ---: | ---: |
| `pyPrivilegeView` · `pyPrivilegeUpdate` | 699 tiap | ⭐ **0** |
| `pyPrivilege` · `pyWhenNotPrivilege` | 419 tiap | ⭐ **0** |
| `pyActionPrivilegeList` · `pyPrivilegeName` · `pySectionReferencePrivilege*` | 33 | ⭐ **0** |
| `WorkGroup` | 16 | ⭐ **0** |
| `pyAssociatedPrivileges` | 520 | 431 — ⛔ seluruhnya nilai baku |
| ⚠️ `pyPrivilegeClass` | 33 | **33** — ⚠️ **hanya nama KELAS, bukan nama hak** |

⭐ **Pola yang sama persis dengan Komite Claim Prop:** daftar hak akses **menyebut kelas tanpa
menyebut hak**.

`[terverifikasi]` Kedua rule penonaktif tombol menggerbangi **keadaan data** — bukan siapa
penggunanya. **52 titik penonaktif kiriman** terpusat di **sepuluh** berkas, ⭐ **seluruhnya di layar
estimasi dan penyesuaian** — nol di layar komite, nol di layar akseptasi. **Sembilan aksi lokal
bernama** tercatat.

#### ✅ Aturan peran ditetapkan SEKALI, lintas modul

`[keputusan work owner — atas rekomendasi asisten]` **2026-09-19** — **aturan peran sistem baru
ditetapkan SEKALI dan berlaku lintas modul**, bukan dirancang ulang per modul.

⚠️ **Alasannya, dan ia bukan soal kerapian:** lubang yang sama muncul **dua modul beruntun** —
**Claim Fac In** *(nol gerbang wewenang; 15 medan hak akses kosong seluruhnya)* dan **Komite Claim
Prop** *(ADR-0014 tidak ditegakkan; AC 67 modul itu menyerahkan wewenang ke "aturan peran sistem
baru" yang belum ditulis)*. ⛔ Merancangnya per modul akan **mengulang lubang yang sama untuk
ketiga kalinya**.

⛔ **Butir 6 MENYEMPIT, tidak tertutup.** Yang diputuskan **cara menetapkannya**; **isinya sendiri
belum ditulis di mana pun**. Bab 19 butir 6.

### Bab 16 — Vonis 15 ADR

⚠️ **Dasar bab ini:** `[keputusan work owner]` **2026-09-19** *(A4-4)* — *"ADR berlaku
lintas-modul"*. Kelima belas ADR **mengikat seluruh modul**, bukan hanya Claim Life. Karena itu
keenam pertentangan di bawah **menjadi penyimpangan sadar**, dan masuk Bab 17.

| ADR | Pokok | Vonis | Bukti satu baris |
| --- | --- | --- | --- |
| **0001** | batas konteks Claim/Komite Life | **TIDAK MENYENTUH** | lingkupnya Life; bentuknya berulang di sini lewat kasus anak komite |
| **0002** | RBAC tiga peran | **TIDAK MENYENTUH** | korpus **nol** memuat model peran — lihat Bab 15 |
| ⭐ **0003** | uang bukan `float` | ⛔ **MENENTANG** | ⭐ **berkaki dua** — `TotalAdjustment` bertipe `Double` **dan** uang lewat sebagai teks yang ditambal koma-ke-titik |
| **0004** | alamat jadi env var | *(tidak dinilai — `superseded`, digantikan 0013)* | — |
| **0005** | `IsPEGAPROD` jadi flag lingkungan | ⛔ **MENENTANG** | rule **hidup**, 10 perujuk, menggerbangi efek keluar |
| **0006** | penomoran lewat stored procedure | ✅ **MENDUKUNG** | procedure yang **sama** dipakai |
| **0007** | jejak audit tiap transisi | ✅ **MENDUKUNG** | kronologi ditulis; kelengkapannya belum diuji |
| **0008** | efek keluar asinkron, tidak memblokir | ⛔ **MENENTANG sebagian** | kegagalan **tidak dicatat** — 2 dari 52 punya jalur gagal |
| **0009** | migrasi penuh data Life | **TIDAK MENYENTUH** | lingkupnya Life |
| **0010** | berkas tetap di Google Storage | ✅ **MENDUKUNG** | berkas memang di Google Storage |
| **0011** | unit status = baris penyesuaian | ⛔ **MENENTANG sebagian** | ada **dua tingkat lagi** di bawahnya — objek dan item objek |
| ⭐ **0012** | wewenang kirim komite bergantung `Type` | **TIDAK MENYENTUH** | ⭐ **nol** pola `.Type` = `TP`/`TR`/`QP`/`QR` di **482 berkas** |
| **0013** | alamat di-*lookup* dari tabel tautan layanan | ⛔ **MENENTANG sebagian** | satu dari delapan memakai **URL langsung** |
| **0014** | keputusan komite per `KomiteID` | ✅ **MENDUKUNG** | roster per `KomiteID` memang bentuknya |
| **0015** | efek keluar komite wajib berhasil | ⛔ **MENENTANG** | **nol** outbox, **nol** antre-ulang |

**Rekap:** ✅ MENDUKUNG **4** · ⛔ MENENTANG **6** · TIDAK MENYENTUH **4** · tidak dinilai **1**
*(superseded)* — jumlah **15**.

> ⛔ **RALAT terhadap ronde ulang §B10.** Rekap lamanya dikutip utuh, tidak dihapus:
> *"⭐ **Menguatkan 5 · MENENTANG 6 · tidak menyentuh 3 · tidak dapat diputuskan 1.**"*
>
> ⚠️ **Rekap itu tidak menjumlah menjadi 15** *(5+6+3+1 = 15, tetapi tabelnya sendiri memuat
> 4 menguatkan dan 4 tidak menyentuh)*. **Dihitung dua cara atas tabel yang sama:** *(a)* mencacah
> baris per vonis; *(b)* memeriksa jumlah totalnya harus 15. ⭐ **Keduanya memberi 4 · 6 · 4 · 1.**
>
> ⭐ Ditambah satu vonis yang **sudah bergeser** sesudahnya: **ADR-0012** dari *"tidak dapat
> diputuskan"* menjadi **TIDAK MENYENTUH** *(ronde 3 §D5)*.

⛔ **Nol revisi ADR diusulkan.** ADR hanya dicabut work owner.

### Bab 17 — Titik yang **SENGAJA DIUBAH** dari Pega

⭐ **Ada dua belas.** **Dihitung dua cara, dan keduanya sepakat:** *(a)* 5 titik terurai + 6
pertentangan ADR + 1 dari keputusan ronde 5 = **12**; *(b)* delta terhadap ronde 4 §F1b
*("5 terurai + 6 dari ADR" = 11)* ditambah keputusan ronde 5 = **12**.

> ⛔ **RALAT.** Satu kalimat sumber **tidak dipakai** dan dikutip di sini apa adanya:
> ronde 2 §H1b menulis *"**2 yang terurai + 6 yang tercatat di ronde ulang = 8**"*. ⚠️ Kalimat itu
> **sudah usang saat ditulis** — tabel di atasnya pada berkas yang sama sudah menyebut *"Jumlah
> titik yang terurai kini **LIMA**"*. **Yang dipakai: lima.**

| # | Titik | Di Pega | Di sistem baru | Bab |
| --- | --- | --- | --- | --- |
| **1** | **Sumber `Quotation`** | dibaca dari halaman `pyWorkPage.Quotation`, diuji **38 rule `When`** | ⭐ diambil **langsung dari tabel polis**, dibaca **sekali** saat kasus dimuat; ke-38 rule diganti **satu fungsi klasifikasi** | 4 |
| **2** | **Daftar penyetuju kosong** | kasus komite lahir **tanpa penyetuju, tanpa pesan** | penyerahan **DITOLAK** dengan galat yang terlihat | 9 |
| **3** | **Sumber `.KomiteNo`** | potongan teks mulai **huruf ke-19** dari kunci internal | **pengenal kasus komitenya** | 9 |
| **4** | **`FlagOnGoingCommitte`** | disimpan sebagai teks `"Send Commite"` | **tidak disimpan** — diturunkan dari kasus komitenya | 9 |
| **5** | **Nama kolom hasil SQL** | **5 dari 6** dan **2 dari 3** alias berbohong, dengan pemetaan **berbeda** | **dinamai sesuai isinya** di batas pembacaan + **tabel pemetaan** nama-lama → nama-benar | 18 |
| **6** | **Satu rule, satu nasib** | rule yang sama dinilai berbeda per modul | **rule dipisahkan dari kehidupannya** — satu salinan, hidup-matinya bergantung halaman yang diujinya | 4 |
| ⭐ **7** | **Representasi uang** *(ADR-0003)* | `Double` untuk total penyesuaian; tiga nilai uang sebagai **teks**; satu ditambal koma-ke-titik | **desimal di seluruh jalur** | 6 |
| **8** | **`IsPEGAPROD`** *(ADR-0005)* | rule `When` hidup, **10 perujuk** | **flag lingkungan eksplisit** | 12 |
| **9** | **Alamat layanan keluar** *(ADR-0013)* | **1 dari 8** memakai URL langsung | seluruhnya lewat **tabel tautan layanan** | 12 |
| **10** | **Kegagalan efek keluar dicatat** *(ADR-0008)* | **tidak dicatat** — 2 dari 52 punya jalur gagal | **dicatat dan terlihat** | 12 |
| **11** | **Efek keluar wajib berhasil** *(ADR-0015)* | **nol** outbox, **nol** antre-ulang | **outbox transaksional**, at-least-once | 12 |
| ⭐ **12** | **Ketujuh pesan validasi estimasi** | hanya **2** menghalangi, lewat penanda `IsError = 2` yang artinya **hanya ada di catatan pengembang** | **galat validasi yang tegas**, dengan pembedaan **yang terbaca** | 5 |

⚠️ **Titik 11 diturunkan dari ADR-0011** *(unit status)* **tidak dimasukkan sebagai titik
tersendiri** — ia **pertentangan struktur**, bukan perilaku yang diubah: sistem baru **meniru** tiga
tingkatnya apa adanya. ⛔ Dicatat di Bab 16 sebagai vonis, bukan di sini sebagai penyimpangan.

### Bab 18 — Nama yang menyesatkan, dan aturan menamai ulang

⭐ `[terverifikasi]` **Dua rule pembaca riwayat klaim berbohong pada nama kolom hasilnya, dan
berbohong dengan pemetaan yang BERBEDA:**

| Rule | Kolom asli | Alias | Kenyataannya |
| --- | --- | --- | --- |
| pembaca riwayat A | tanggal kejadian | `START_DATE` | tanggal kejadian dinamai *tanggal mulai* |
| pembaca riwayat A | pengenal Pega | `BRANCH_CODE` | dinamai *kode cabang* |
| pembaca riwayat A | nomor klaim | `BRANCH_NAME` | dinamai *nama cabang* |
| pembaca riwayat A | penyebab kerugian | `BUSINESS_CODE` | dinamai *kode bisnis* |
| ⭐ pembaca riwayat B | tanggal kejadian | `BRANCH_CODE` | ⭐ **pemetaan BERBEDA dari A** |
| ⭐ pembaca riwayat B | nomor klaim | `BUSINESS_CODE` | ⭐ **pemetaan BERBEDA dari A** |

⭐ `[terverifikasi]` *(ronde 3 §E2)* **Pemanggilnya satu-satu, dan nol pemanggil memakai keduanya** —
jadi kedua pemetaan bohong itu **tidak pernah bertemu dalam satu jalur**. ⚠️ Tetapi keduanya membaca
**tabel yang sama**.

✅ `[keputusan work owner — atas rekomendasi asisten]` **Kolom diberi nama sesuai isinya di batas
pembacaan**, disertai **tabel pemetaan nama-lama → nama-benar**.

⚠️ **Tabel pemetaan itu WAJIB ada.** Tanpanya, orang yang membandingkan keluaran lama dan baru akan
mengira **datanya berubah**, padahal hanya **namanya** yang dibetulkan. ✅ Sejalan **AC 106 Claim
Prop**. ⚠️ **Penyimpangan sadar** — Bab 17 titik 5.

⛔ **Berapa tepatnya alias yang berbohong: belum punya data.** Empat jendela memberi **50 · 48 · 42 ·
371** pasangan. Bab 19 butir 8.

---

## Testing Decisions

**Apa yang membuat sebuah test baik di sini.** Test menguji **perilaku yang terlihat dari luar** —
apa yang tersimpan, apa yang ditolak, apa yang terkirim — ⛔ **bukan** langkah internal, nama rule,
maupun urutan pemanggilan. Setiap AC di bawah ditulis dalam bentuk *"Test yang menemukan … gagal"*
supaya sasarannya tidak dapat ditafsir dua cara.

**Sambungan uji yang dipakai** *(dari yang paling tinggi)*:

1. ⭐ **Batas layanan** — satu aksi pengguna masuk, keadaan tersimpan keluar. **Sambungan utama**,
   dan sedapat mungkin **satu-satunya**. Di sinilah AC transaksi *(41)*, wewenang *(44)*, dan
   validasi *(14–18)* diuji.
2. **Batas basis data** — untuk AC yang menuntut bentuk tersimpan: uang desimal, pengenal berkas
   unik, jejak audit.
3. **Batas layanan luar** — ganda uji untuk Kasir, Google Storage, surel, dan konversi; dipakai
   hanya untuk AC efek keluar dan antre-ulang.

⚠️ **Yang sengaja TIDAK dijadikan sambungan:** lompatan per lini produk *(Bab 3)*. Ia bentuk
internal Pega; yang diuji adalah **hasilnya** — layar dan rincian yang muncul per lini — bukan
mekanisme lompatannya.

**Preseden.** Bentuk AC dan sambungan uji mengikuti `claim-prop/spec.md` dan
`komite-claim-prop/spec.md`, yang keduanya sudah dipakai menurunkan tiket.

⚠️ **Yang belum dapat diuji sama sekali:** ke-14 AC ber-`[terbuka]`. Ia **menandai tempat yang belum
punya sasaran uji**, dan ⛔ **tidak boleh ditambal dengan tebakan**.

---

## Acceptance Criteria

*Bab ini **tidak memutuskan apa pun.** Ia menyatakan ulang — dalam bentuk yang bisa diuji dari luar —
keputusan yang **sudah** ada di Bab 1–18. Bila sebuah butir terasa seperti keputusan baru, ia salah
tulis: laporkan, jangan dilaksanakan.*

### Lingkup dan identitas

1. `[terverifikasi]` Objek kerja modul ini adalah **klaim fakultatif masuk**, terpisah dari objek
   kerja Claim Prop. Test yang menemukan keduanya berbagi satu objek kerja **gagal**. *(Bab 1)*
2. `[terverifikasi]` Data klaim tersimpan **tiga tingkat**: objek pertanggungan → item objek → baris
   penyesuaian. Test yang menemukan item objek tanpa objek induk, atau baris penyesuaian tanpa item
   objek, **gagal**. *(Bab 1 · 8)*
3. `[terverifikasi]` Objek pertanggungan dan item objek adalah **data tersimpan**, bukan halaman
   kerja sementara. Test yang menemukan keduanya hilang sesudah kasus ditutup dan dibuka kembali
   **gagal**. *(Bab 1)*
4. `[keputusan work owner]` ⭐ **Perilaku yang ditiru adalah perilaku salinan dari sistem
   PRODUKSI.** Test yang membandingkan hasil terhadap salinan pengembangan **gagal**. ⚠️
   `[terverifikasi]` Korpus memuat **tiga** sistem sumber, dan dua salinan beridentitas empat bagian
   **sama** terbukti **berbeda isinya** — jadi *"sama dengan Pega"* **tidak bermakna** tanpa
   menyebut salinan mana. *(Bab 1)*

   > ✅ **butir 10 DITUTUP 2026-09-19.** Teks lamanya **dikutip utuh, tidak dihapus**:
   > *"`[terbuka]` **butir 10** — **rule dari lingkungan mana yang ditiru** belum diputuskan …
   > ⛔ Selama itu belum dijawab, **kesetaraan perilaku terhadap Pega tidak punya sasaran uji yang
   > pasti**."*

### Daur hidup kasus

5. `[terverifikasi]` Sebuah klaim berjalan melalui tahap-tahap yang **ditetapkan satu berkas alur**,
   dan tahap yang tidak ada di sana **tidak dibuat**. Test yang menemukan tahap tambahan **gagal**.
   *(Bab 2)*
6. `[terverifikasi]` Kasus yang selesai berstatus akhir **selesai-tuntas** dan **tidak muncul lagi**
   di kotak kerja siapa pun. Test yang menemukan kasus selesai masih menggantung **gagal**. *(Bab 2)*
7. `[terverifikasi]` **Titik lompat darurat tidak dibuat.** Test yang menemukan jalan masuk ke alur
   selain pendaftaran klaim **gagal**. ⭐ Tidak membangunnya adalah **paritas, bukan penyimpangan**:
   sapuan seluruh korpus menemukan **nol** pembangkit tiket, sehingga di Pega pun ketujuh tiket itu
   **tidak pernah dilempar**. *(Bab 2)*

### Registrasi dan polis

8. `[terverifikasi]` Klaim **wajib** menunjuk satu polis. Test yang berhasil menyimpan klaim tanpa
   polis **gagal**. *(Bab 1)*
9. `[terverifikasi]` Nomor klaim dibuat oleh **stored procedure pembangkit nomor urut**, bukan
   dihitung aplikasi. Test yang menemukan dua klaim bernomor sama **gagal**. *(Bab 10 · ADR-0006)*
10. `[terverifikasi]` Pengguna **diperingatkan** bila klaim dengan tanggal kejadian sama sudah ada
    pada polis yang sama. Test yang tidak memunculkan peringatan itu **gagal**. *(Bab 5)*
11. `[terverifikasi]` Penyebab kerugian diambil dari **daftar baku**, bukan diketik bebas. Test yang
    berhasil menyimpan penyebab di luar daftar **gagal**. *(Bab 1)*
12. `[terverifikasi]` Riwayat klaim pada polis yang sama **dapat dilihat** sebelum klaim baru
    disimpan. Test yang tidak menemukannya **gagal**. *(Bab 18)*

### Klasifikasi lini produk

13. `[keputusan work owner]` Klasifikasi lini ditentukan **satu kali di satu tempat**. Test yang
    menemukan **dua salinan** logika klasifikasi **gagal**. *(Bab 4)*
14. `[keputusan work owner]` ⚠️ `[penyimpangan sadar]` Nilai penawaran diambil **langsung dari tabel
    polis** dan **dibaca sekali** saat kasus dimuat, bukan dari halaman kerja. Test yang menemukannya
    dibaca ulang berkali-kali dalam satu kasus **gagal**. *(Bab 4 · Bab 17 titik 1)*
15. `[terbuka]` **butir 1** — ✅ **nama tabelnya SUDAH ditetapkan: `T_QUOTATIONDATA`**
    `[keputusan work owner]` 2026-09-19. ⛔ **Nama kolomnya belum**, dan `[terverifikasi]` **tabelnya
    belum ada** — **0 dari 482 berkas** menyebutnya; ia **sedang dibangun**. Karena itu **AC 14
    tetap tanpa sasaran uji yang konkret**. ⛔ Jangan ditambal dengan tebakan. *(Bab 4)*
16. `[keputusan work owner]` ⚠️ `[penyimpangan sadar]` Rule klasifikasi warisan **dipindahkan sekali
    sebagai satu himpunan**; hidup-matinya **tidak ikut dipindahkan**. Test yang menemukan rule yang
    sama disalin dua kali untuk dua modul **gagal**. ⭐ **Dasarnya kini terbukti, bukan diasumsikan:**
    di Claim Fac In propertinya **diisi** — lewat penyalinan halaman dari penawaran — sedangkan di
    Claim Prop halaman itu tidak ada sama sekali. **Satu salinan, dua nasib.** *(Bab 4 · Bab 17
    titik 6)*
17. `[terverifikasi]` Lini produk yang dikenali modul ini mencakup **kendaraan bermotor**,
    **kecelakaan diri**, **perjalanan**, dan **lainnya**. Test yang menemukan lini di luar daftar
    yang tercatat **gagal**. *(Bab 3)*
18. `[terverifikasi]` ⚠️ Cacat pasangan lompatan lini pada rule pengumpul data **ditiru apa adanya**.
    Test yang menemukan perilakunya **diperbaiki** **gagal**. ⭐ Dasarnya: rule itu **nol penugasan ke
    kasus** dari 61 penugasan, sehingga cacatnya **hanya menyentuh tampilan satu popup**.
    *(Bab 3)*

### Estimasi

19. `[terverifikasi]` Nilai estimasi dicatat **per objek pertanggungan**. Test yang menemukan
    estimasi tanpa objek **gagal**. *(Bab 5)*
20. `[terverifikasi]` ⚠️ Nilai estimasi **negatif ditolak**. Test yang berhasil menyimpan estimasi
    negatif **gagal**. *(Bab 5)*
21. `[terverifikasi]` ⚠️ Persentase estimasi **negatif ditolak**. Test yang berhasil menyimpannya
    **gagal**. *(Bab 5)*
22. `[terverifikasi]` ⚠️ Tanggal estimasi yang tidak masuk akal **ditolak**. Test yang berhasil
    menyimpannya **gagal**. *(Bab 5)*
23. `[terverifikasi]` ⚠️ Estimasi yang **persentase dan nilainya keduanya nol** ditandai. Test yang
    melewatkannya **gagal**. *(Bab 5)*
24. `[terverifikasi]` ⚠️ Estimasi yang **melebihi nilai pertanggungan** ditandai. Test yang
    melewatkannya **gagal**. *(Bab 5)*
25. `[terverifikasi]` ⚠️ Estimasi yang **melampaui batas tanggung jawab** ditandai. Test yang
    melewatkannya **gagal**. *(Bab 5)*
26. `[keputusan work owner]` ⚠️ `[penyimpangan sadar]` ⭐ **KETUJUH pemeriksaan di atas
    menggagalkan penyimpanan.** Test yang berhasil menyimpan estimasi yang melanggar salah satu dari
    ketujuhnya **gagal**; test yang menemukan salah satunya hanya **memperingatkan** juga **gagal**.
    ⛔ Dan test yang menemukan sebuah aturan hanya dapat diketahui dari **nilai sebuah penanda**
    **gagal**. *(Bab 5 · Bab 17 titik 12)*
27. `[keputusan work owner]` ⭐ **Golongannya tegas: tidak ada pemeriksaan estimasi yang berstatus
    peringatan.** Test yang menggolongkan salah satu dari ketujuhnya sebagai peringatan **gagal**.
    *(Bab 5)*

    > ✅ **butir 12 DITUTUP 2026-09-19.** Teks lamanya **dikutip utuh, tidak dihapus**:
    > *"`[terbuka]` **butir 12** — **pesan mana masuk golongan mana belum diputuskan**. Korpus
    > memberi keadaan sekarang *(2 menghalangi · 2 hanya menonaktifkan tombol unduh · 2 peringatan ·
    > 1 belum punya data)*, ⛔ **bukan niatnya**."*
    >
    > ⚠️ **Dan ia melahirkan butir BARU 25** — lihat AC 114.
28. `[terbuka]` **butir 13** — ⛔ apakah pesan pada sistem lama **menghalangi penyerahan** tidak
    terbaca dari ekspor; rule validasinya dipanggil sebagai **aksi saat pengguna mengetik**, bukan
    validasi penyerahan. *(Bab 5)*
29. `[terverifikasi]` Total estimasi **terhitung ulang** setiap kali satu baris berubah. Test yang
    menemukan total tertinggal dari barisnya **gagal**. *(Bab 5)*
30. `[terverifikasi]` ⚠️ `[penyimpangan sadar]` Penerbitan dokumen lembar muka **terhalang** bila
    validasi estimasi belum lolos. Test yang berhasil menerbitkannya dengan estimasi bertanda galat
    **gagal**. *(Bab 5 · Bab 10)*

### Spreading

31. `[terverifikasi]` Kerugian dibagi ke para penanggung sesuai porsinya, dan pembagian itu
    **dipisahkan per mata uang**. Test yang mencampur dua mata uang dalam satu pembagian **gagal**.
    *(Bab 7)*
32. `[terverifikasi]` Jumlah porsi pada satu mata uang **tepat seratus persen**. Test yang berhasil
    menyimpan jumlah di luar itu **gagal**. *(Bab 7)*
33. `[terverifikasi]` ⚠️ Rantai pemeriksaan spreading berperilaku **ATAU**, bukan DAN — satu syarat
    terpenuhi sudah cukup. Test yang menuntut **seluruh** syarat terpenuhi **gagal**. ⭐ Dasarnya
    kode arah **5** pada rule pemeriksa spreading. *(Bab 7)*

### Penyesuaian

34. `[terverifikasi]` Baris penyesuaian dicatat terhadap **item objek**, bukan langsung terhadap
    klaim. Test yang menemukan penyesuaian tanpa item objek **gagal**. *(Bab 8)*
35. `[terverifikasi]` Nilai sisa barang selamat dan biaya adjuster **dicatat terpisah** dari nilai
    penyesuaian. Test yang menjumlahkan ketiganya menjadi satu medan **gagal**. *(Bab 8)*
36. `[terverifikasi]` Adjuster ditunjuk dari **daftar yang sudah ada**. Test yang berhasil menyimpan
    nama adjuster yang diketik bebas **gagal**. *(Bab 8)*
37. `[terverifikasi]` Total penyesuaian **dihitung dari baris-barisnya**, bukan disimpan sebagai
    angka yang dimasukkan sendiri. Test yang menemukan total berbeda dari jumlah barisnya **gagal**.
    *(Bab 6 · 8)*
38. `[terverifikasi]` Total penyesuaian **tidak disimpan sebagai kolom tersendiri** pada klaim. ⭐
    Dasarnya: di Pega ia **tidak pernah sampai ke kolom tersimpan** — 71 penugasan, 8 ke kasus, nol
    berasal darinya. Test yang menemukannya tersimpan dan **berbeda** dari jumlah barisnya **gagal**.
    *(Bab 6)*

### Uang

39. `[keputusan work owner]` ⚠️ `[penyimpangan sadar]` Seluruh nilai uang **bertipe desimal** di
    setiap lapisan — penyimpanan, perhitungan, dan kontrak layanan, dengan ⭐ **20 angka, 8 di
    belakang koma**. Test yang menemukan nilai uang melewati **bilangan pecahan biner** **gagal**;
    test yang menemukan **pembulatan di tengah jalan** **gagal**; test yang menemukan ketelitian
    **berbeda dari Claim Prop** pada tabel akseptasi bersama **gagal**. ⭐ Tampilan layar **4 angka
    di belakang koma**. *(Bab 6 · Bab 17 titik 7 · ADR-0003)*
40. `[keputusan work owner]` ⚠️ `[penyimpangan sadar]` **Tidak ada nilai uang yang disimpan sebagai
    teks.** Test yang menemukan sebuah nilai uang perlu **ditambal pemisah desimalnya** sebelum dapat
    dihitung **gagal**. *(Bab 6 · Bab 17 titik 7)*
41. `[terverifikasi]` Setiap nilai uang tersimpan lengkap dengan **mata uang** dan **kursnya**. Test
    yang menemukan nilai uang tanpa keduanya **gagal**. *(Bab 1 · 7)*
42. `[keputusan work owner]` ⭐ **Ketelitiannya SAMA dengan Claim Prop** — 20 angka, 8 di belakang
    koma, tampilan 4. Test yang menemukan kedua modul memakai ketelitian berbeda pada tabel
    akseptasi bersama **gagal**. *(Bab 6)*

    > ✅ **butir 14 DITUTUP 2026-09-19.** Teks lamanya **dikutip, tidak dihapus**: *"`[terbuka]`
    > **butir 14** — ⛔ **ketelitian desimal yang dipakai belum ditetapkan** untuk modul ini."*

114. `[terbuka]` **butir 25** — ⭐ **BARU 2026-09-19.** ⛔ **Berapa banyak klaim lama sudah memuat
     keadaan yang kini menjadi galat** — estimasi bernilai-dan-berpersentase nol, atau estimasi
     melebihi nilai pertanggungan — **belum diperiksa**. ⚠️ Di sistem lama hanya **dua** dari tujuh
     pemeriksaan yang menghalangi; kini **ketujuhnya** menggagalkan penyimpanan. ⚠️ **Memblokir
     migrasi, bukan pembangunan.** *(Bab 5)*

### Penyerahan ke komite

43. `[terverifikasi]` Klaim diserahkan ke komite dengan **membuat kasus komite anak**, membawa
    **lima penunjuk posisional** ke baris yang dinilai. Test yang menemukan kasus komite tanpa
    penunjuk lengkap **gagal**. *(Bab 8 · 9)*
44. `[keputusan work owner]` ⚠️ `[penyimpangan sadar]` Penyerahan **DITOLAK dengan galat yang
    terlihat** bila daftar penyetujunya kosong. Test yang berhasil melahirkan kasus komite tanpa satu
    pun penyetuju **gagal**. *(Bab 9 · Bab 17 titik 2)*
45. `[keputusan work owner]` ⚠️ `[penyimpangan sadar]` Nomor komite diambil dari **pengenal kasus
    komitenya**. Test yang menemukannya diambil dari **potongan posisi karakter** atas kunci internal
    **gagal**. *(Bab 9 · Bab 17 titik 3)*
46. `[keputusan work owner]` ⚠️ `[penyimpangan sadar]` Keadaan *"komite sedang berjalan"*
    **diturunkan dari ada-tidaknya kasus komite**, bukan disimpan sebagai kolom penanda. Test yang
    menemukan kolom penanda itu dibuat **gagal**. *(Bab 9 · Bab 17 titik 4)*
47. `[terverifikasi]` Pembuatan kasus komite **dibatalkan** bila kasus induknya masih memuat pesan
    kesalahan. Test yang berhasil membuatnya **gagal**. *(Bab 9)*
48. `[terverifikasi]` Roster penyetuju menyimpan **akun operator** sebagai pengenal penyetuju dan
    **jabatan** sebagai keterangan, dengan keputusan mulai **kosong**. Test yang menukar keduanya
    **gagal**. *(Bab 9 · ADR-0014)*
49. `[terbuka]` **butir 2** — ⛔ **kenapa nomor komite ditulis dua kali berturut-turut** tidak ada
    di korpus, sehingga **tidak diketahui apakah penulisan pertama memang dianggap kurang**.
    *(Bab 9)*
50. `[terbuka]` **butir 3** — ⛔ **belum dibuktikan** kolom penanda komite di modul ini adalah kolom
    yang **sama** dengan yang dibuang di Komite Claim Prop. **Nama sama bukan bukti.** *(Bab 9)*

### Dokumen

51. `[terverifikasi]` Modul menerbitkan **lembar muka klaim**, **nota kerugian sementara**, dan
    **nota kerugian tetap**. Test yang menemukan jenis dokumen di luar ketiganya **gagal**. *(Bab 10)*
52. `[terverifikasi]` Penomoran dokumen memakai **pembangkit nomor urut yang sama** dengan penomoran
    klaim. Test yang menemukan dua dokumen bernomor sama **gagal**. *(Bab 10 · 13)*

### Penyimpanan berkas

53. `[terverifikasi]` Berkas lampiran **tetap disimpan di Google Storage**; token diterbitkan
    **basis data**, bukan aplikasi. Test yang menemukan aplikasi menerbitkan token sendiri **gagal**.
    *(Bab 11 · ADR-0010)*
54. `[keputusan work owner]` Pengenal berkas dibuat dengan cara **yang sudah ditambal** — ketelitian
    **nanodetik** **dan** nilai unik sejagat. Test yang menemukan **dua unggahan berdekatan
    menghasilkan pengenal sama** **gagal**. *(Bab 11)*
55. `[terverifikasi]` Berkas yang sudah diunggah **dapat dibuka kembali**. Test yang tidak menemukan
    jalur itu **gagal**. *(Bab 11)*
56. `[keputusan work owner]` ⭐ **Kedua jenis berkas video ikut diterima** — keduanya ada di
    salinan **produksi**, dan salinan produksilah yang ditiru *(AC 4)*. Test yang menolak keduanya
    **gagal**. *(Bab 11)*

    > ✅ **butir 9 DITUTUP 2026-09-19**, sebagai akibat langsung keputusan lingkungan. Teks lamanya
    > **dikutip, tidak dihapus**: *"`[terbuka]` **butir 9** — ⛔ apakah **dua jenis berkas video**
    > yang hanya dikenali salinan produksi ikut diterima **belum diputuskan**."*

### Efek keluar

57. `[terverifikasi]` Data akseptasi klaim **terkirim ke sistem Kasir**. Test yang tidak menemukan
    kiriman itu **gagal**. *(Bab 12)*
58. `[terbuka]` **butir 7** — ⛔ apakah panggilan konversi yang berjalan **tanpa saringan** memang
    dikehendaki **belum diputuskan**. ⭐ `[terverifikasi]` Gerbang berflag mati mematikan
    **gerbangnya**, bukan **langkahnya** — sehingga panggilan itu berjalan **lebih sering**, bukan
    tidak pernah, dari **empat** pemanggil yang **nol** bergerbang hidup. *(Bab 12)*
59. `[keputusan work owner]` ⚠️ `[penyimpangan sadar]` Kegagalan efek keluar **dicatat dan terlihat**.
    Test yang menemukan kegagalan **hilang diam-diam** **gagal**. *(Bab 12 · Bab 17 titik 10 ·
    ADR-0008)*
60. `[keputusan work owner]` ⚠️ `[penyimpangan sadar]` Efek keluar dijalankan lewat **outbox
    transaksional** dengan jaminan **at-least-once**; keputusan tidak dianggap tuntas sampai seluruh
    efek terkirim atau ditandai **perlu intervensi**. Test yang menemukan efek keluar berjalan
    **tanpa jejak outbox** **gagal**. *(Bab 12 · Bab 17 titik 11 · ADR-0015)*
61. `[keputusan work owner]` ⚠️ `[penyimpangan sadar]` **Seluruh** alamat layanan keluar di-*lookup*
    dari **tabel tautan layanan**. Test yang menemukan satu pun alamat ditanam sebagai URL langsung,
    konstanta, atau env var **gagal**. *(Bab 12 · Bab 17 titik 9 · ADR-0013)*
62. `[keputusan work owner]` ⚠️ `[penyimpangan sadar]` Gerbang lingkungan menjadi **flag lingkungan
    eksplisit**. Test yang menemukan gerbang itu berupa rule yang kondisinya tidak terbaca **gagal**.
    *(Bab 12 · Bab 17 titik 8 · ADR-0005)*
63. `[terverifikasi]` Surel dikirim pada tahap yang ditentukan. Test yang menemukan surel terkirim di
    luar tahap itu **gagal**. *(Bab 12)*
64. ⚠️ `[terverifikasi]` Kredensial layanan keluar **tidak pernah tertulis di dalam rule maupun kode**
    yang beredar. Test yang menemukannya tertanam **gagal**. ⚠️ Kredensial yang kini beredar di dalam
    ekspor **sebaiknya diganti sesudah migrasi** — urusan tim pemilik layanan. *(Bab 12)*

### Transaksi

65. `[keputusan work owner]` **Satu aksi pengguna menghasilkan satu transaksi.** Test yang menemukan
    keadaan **separuh tersimpan** sesudah kegagalan di tengah **gagal**. *(Bab 13)*
66. `[keputusan work owner]` ⛔ **Tidak ada perintah penutup transaksi di dalam rule pembaca.** Test
    yang menemukan sebuah pembacaan menutup transaksi pemanggilnya **gagal**. *(Bab 13)*
67. `[keputusan work owner]` ⚠️ **Pengecualian: pembangkit nomor urut.** Nomor yang sudah terbit
    **tidak ikut dibatalkan** saat transaksi induknya batal. Test yang menemukan nomor yang sama
    **terpakai ulang** sesudah pembatalan **gagal**. *(Bab 13 · ADR-0006)*

### Jejak audit

68. `[keputusan work owner]` Setiap transisi status dan setiap jalur balik merekam **siapa** dan
    **kapan**. Test yang menemukan transisi tanpa keduanya **gagal**. *(Bab 14 · ADR-0007)*
69. `[terverifikasi]` Jejak kronologi **ditulis** pada titik-titik yang sudah ada di sistem lama.
    Test yang tidak menemukannya **gagal**. *(Bab 14)*
70. `[terbuka]` **butir 5** — ⛔ **kelengkapan jejak kronologi belum diuji**, dan dua catatan
    pengembang meminta sebagian penulisnya dihapus — salah satunya **masih hidup**. ⛔ Apakah
    penghapusan itu dikehendaki **tidak terbaca**. *(Bab 14)*
71. `[terverifikasi]` Angka yang **disimpan** dan angka yang **dihitung ulang** dapat dibedakan.
    Test yang tidak dapat membedakannya **gagal**. *(Bab 6 · 8)*

### Wewenang

72. `[terverifikasi]` ⭐ **Sistem lama tidak punya satu pun gerbang layar berdasarkan wewenang** —
    ke-15 medan hak akses yang menamai hak **kosong seluruhnya**, dan penonaktifan tombol
    menggerbangi **keadaan data**, bukan peran. Test yang menemukan gerbang wewenang **ditiru** dari
    sistem lama **gagal** — tidak ada yang dapat ditiru. *(Bab 15)*
73. `[terbuka]` **butir 6** — ✅ **CARA menetapkannya sudah diputuskan** `[keputusan work owner]`
    2026-09-19: **sekali, lintas modul**, bukan per modul. ⛔ **Isinya masih belum ditulis di mana
    pun**, sehingga AC wewenang **tetap tanpa sasaran uji**. ⚠️ Ini **bukan** temuan kecil: tanpa
    aturan peran, modul dapat lulus seluruh AC lain dan tetap membiarkan siapa pun melakukan apa
    pun. *(Bab 15)*
74. `[terverifikasi]` Penonaktifan tombol pada layar estimasi dan penyesuaian **ditiru apa adanya**,
    dan ia bergantung **keadaan data**. Test yang menemukannya bergantung peran **gagal** sampai
    butir 6 dijawab. *(Bab 15)*

### Migrasi

75. `[terbuka]` **butir 15** — ⛔ **migrasi penuh lawan koeksistensi belum diputuskan** untuk modul
    ini. *(Bab 1)*
76. `[keputusan work owner]` Baris data lama **dipindahkan seperti aslinya** — ⛔ tidak dihitung
    ulang, tidak ditambal, tidak ditolak. Test yang menemukan nilai lama **berubah** saat migrasi
    **gagal**. *(Bab 3)*
77. ⚠️ `[keputusan work owner]` **Risiko yang diterima sadar:** baris yang dibuat sesudah migrasi
    dapat **berbeda perilakunya** dari baris lama, dan **tidak ada penanda yang membedakannya**.
    Laporan yang menjumlahkan keduanya **mencampur dua perilaku**. Test yang menuntut keduanya
    identik **gagal**. *(Bab 3)*
78. `[terbuka]` **butir 16** — ⛔ apakah perlu **penanda** yang membedakan baris sebelum dan sesudah
    perbaikan **belum diputuskan**. *(Bab 3)*

### Nama kolom dan pembacaan

79. `[keputusan work owner]` ⚠️ `[penyimpangan sadar]` Kolom hasil pembacaan **dinamai sesuai isinya**
    di batas pembacaan. Test yang menemukan kolom dinamai hal yang bukan isinya **gagal**.
    *(Bab 18 · Bab 17 titik 5)*
80. `[keputusan work owner]` **Tabel pemetaan nama-lama → nama-benar WAJIB ada.** Test yang menemukan
    penamaan ulang dilakukan **tanpa** tabel pemetaan **gagal**. ⚠️ Tanpanya, pembanding keluaran lama
    dan baru akan mengira **datanya** berubah. *(Bab 18)*
81. `[terverifikasi]` ⚠️ **Dua rule pembaca riwayat berbohong dengan pemetaan yang BERBEDA**, dan
    **nol pemanggil memakai keduanya**. Test yang menyamakan pemetaan keduanya **gagal**. *(Bab 18)*
82. `[terbuka]` **butir 8** — ⛔ **berapa tepatnya alias yang berbohong belum punya data**: empat
    jendela memberi **50 · 48 · 42 · 371**. *(Bab 18)*

### Aturan baca yang mengikat pembangunan

83. `[terverifikasi]` ⚠️ Gerbang yang benderanya **mati** berarti **gerbangnya** tidak berlaku —
    **langkahnya tetap berjalan**, dan karena itu berjalan **lebih sering**, bukan tidak pernah.
    Test yang memperlakukan langkah bergerbang mati sebagai **tidak pernah berjalan** **gagal**.
    *(Bab 12 · Lampiran)*
84. `[terverifikasi]` ⚠️ **121 baris syarat** tersimpan pada langkah bergerbang mati — syaratnya
    tertulis **tetapi tidak berlaku**. Test yang menegakkan salah satunya **gagal**. *(Lampiran)*
85. `[terverifikasi]` ⚠️ Rantai gerbang yang memuat **kode arah 5** berperilaku **ATAU**, bukan DAN.
    **Enam titik hidup** tercatat. Test yang membacanya sebagai DAN **gagal**. *(Bab 7 · Lampiran)*
86. `[terverifikasi]` ⚠️ Satu titik **keluar-iterasi** hidup: bila syaratnya tidak terpenuhi,
    **perulangan diputus** — bukan langkahnya dilewati. Test yang melewatkan langkahnya saja
    **gagal**. *(Bab 3)*

### Butir yang belum punya sasaran uji

87. `[terbuka]` **butir 4** — ⛔ **`IndexObject` lawan `ObjectIndex`**: dua penunjuk bernama nyaris
    sama diisi dari dua sumber berbeda, dan mana yang mana **tidak terbaca**. *(Bab 8)*
88. `[terbuka]` **butir 11** — ⛔ mesin pembangkit tiket **ada di aplikasi yang sama**; apakah sebuah
    rule di luar ekspor ini dapat melemparnya **tidak terbukti**. *(Bab 2)*
89. `[terbuka]` **butir 17** — ⛔ arti `REQUIRED = -1` lawan `0` pada tanda tangan parameter **belum
    diketahui**, dan ⛔ **tidak dipakai sebagai dasar apa pun**. *(Bab 6)*
90. `[terbuka]` **butir 18** — ⛔ **isi tiap penampung generik `CARI<n>`** tidak terbaca dari rule
    yang memakainya. *(Bab 18)*
91. `[terbuka]` **butir 19** — ⛔ **327 catatan pengembang belum diuji** dengan membuka rule-nya; dari
    459 catatan, yang benar-benar diuji **12**. ⚠️ Di Claim Prop cara golong-dari-kata-kunci terbukti
    **meleset lima kali lipat**. *(Bab 20)*

### Paritas yang ditegaskan

92. `[terverifikasi]` Bentuk data **tiga tingkat** ditiru apa adanya. Test yang meratakannya menjadi
    dua tingkat **gagal**. *(Bab 1 · 8 · ADR-0011)*
93. `[terverifikasi]` Roster penyetuju ditiru apa adanya. Test yang mengubah bentuknya **gagal**.
    *(Bab 9 · ADR-0014)*
94. `[terverifikasi]` Penyimpanan berkas ditiru apa adanya. Test yang memindahkannya keluar dari
    Google Storage **gagal**. *(Bab 11 · ADR-0010)*
95. `[terverifikasi]` Penomoran lewat stored procedure ditiru apa adanya. Test yang menghitung nomor
    di aplikasi **gagal**. *(Bab 10 · ADR-0006)*
96. `[terverifikasi]` ⭐ Mekanisme wewenang kirim komite bergantung jenis klaim **TIDAK ADA** di modul
    ini — **nol** di seluruh 482 berkas. Test yang membangunnya **gagal**. *(Bab 16 · ADR-0012)*
97. `[terverifikasi]` Daftar jenis berkas yang diterima mengikuti daftar yang ada. Test yang menerima
    jenis di luar daftar **gagal**. *(Bab 11)*
98. `[terverifikasi]` ⚠️ Daur hidup kasus diambil dari berkas alur yang ada di **ruleset bernama akun
    perorangan**, dan itulah yang berlaku. Test yang mencarinya di ruleset resmi dan menyatakan
    ekspor rusak **gagal**. *(Bab 1)*

### Layar

99. `[terverifikasi]` Rincian objek pertanggungan yang tampil **mengikuti lini produk** klaim. Test
    yang menampilkan medan lini lain **gagal**. *(Bab 3)*
100. `[terverifikasi]` Layar estimasi, penyesuaian, dan akseptasi **terpisah**. Test yang
     menggabungkannya menjadi satu layar **gagal**. *(Bab 15)*
101. `[terverifikasi]` ⚠️ **Sembilan aksi lokal bernama** tercatat pada layar modul ini. Test yang
     menemukan aksi lokal di luar daftar itu **gagal**. *(Bab 15)*
102. `[terbuka]` **butir 20** — ⛔ **gerbang aksi tombol tidak ketemu di 11 medan yang disisir**; ⛔
     bukan berarti nol. Claim Prop punya 17 baris di medan yang di sini kosong. *(Bab 15)*

### Keutuhan penunjuk

103. `[terverifikasi]` ⚠️ Sistem lama **nol pemeriksaan** bila penunjuk posisional ke baris yang
     dinilai **hilang atau bergeser**. Test yang menuntut pemeriksaan itu ditiru **gagal** — tidak ada
     yang dapat ditiru. *(Bab 8)*
104. `[terbuka]` **butir 21** — ⛔ **apa yang dilakukan sistem baru bila penunjuk itu salah alamat**
     belum diputuskan. ⚠️ Dengan **tiga tingkat**, peluang bergesernya lebih besar daripada Claim
     Prop. *(Bab 8)*

### Sisa yang ditegaskan

105. `[terverifikasi]` ⚠️ **Satu lompatan menggantung** menunjuk label yang tidak ada, **dan
     gerbangnya mati**. Test yang membangun lompatan ke label yang tidak ada **gagal**. *(Bab 3)*
106. `[terverifikasi]` ⚠️ **112 langkah ber-remark** di sistem lama, **22 di antaranya bermetode
     berat** — jalur simpan atau efek keluar yang **MATI**. Test yang menghidupkannya **gagal**.
     *(Bab 12)*
107. `[terverifikasi]` ⚠️ Halaman kerja tampilan **tidak disimpan**. Test yang menemukan isinya
     tersimpan **gagal**. *(Bab 3)*
108. `[terverifikasi]` Daftar berbasis halaman data **dimuat sekali per utas kerja** dan **tidak
     disegarkan** di dalamnya. Test yang menuntut penyegaran otomatis **gagal**. *(Bab 1)*
109. `[terverifikasi]` ⭐ **Properti yang diuji rule klasifikasi BENAR-BENAR DIISI** — bukan oleh
     modul ini, melainkan oleh **modul saudara yang membuat penawaran fakultatif**, lalu sampai ke
     sini lewat **penyalinan halaman utuh**. Test yang menemukan rule klasifikasi **selalu bernilai
     salah** karena propertinya kosong **gagal**. *(Bab 4)*

     > ✅ **butir 22 DITUTUP 2026-09-19.** Teks lamanya **dikutip utuh, tidak dihapus**:
     > *"`[terbuka]` **butir 22** — ⛔ **apakah properti yang diuji tiap rule klasifikasi
     > benar-benar diisi** belum diperiksa; yang terbukti baru **halamannya** ditulis."*
     > Rantai pengisiannya ada di **Bab 4**, dibaca ujung ke ujung dan dihitung dua cara.

113. `[terbuka]` **butir 24** — ⭐ **BARU 2026-09-19.** ⛔ Ketepatan klasifikasi lini modul ini
     **bergantung pada modul lain** yang mengisi halaman penawaran, dan ⛔ **tidak ada satu pun
     pemeriksaan di modul ini** yang akan mengatakannya bila keempat properti itu kosong. Apa yang
     dilakukan aplikasi dalam keadaan itu **belum diputuskan**. ⚠️ **Nomornya 113, bukan 110** —
     nomor menentukan identitas, sub-judul menentukan tempat. *(Bab 4)*
110. `[terverifikasi]` ⚠️ Kelas rule dan nama rule **tidak cukup** sebagai identitas — ruleset dan
     versinya ikut menentukan. Test yang menyamakan dua rule hanya karena namanya sama **gagal**.
     *(Bab 1)*
111. `[terverifikasi]` ⚠️ Nama rule **menyesatkan**: rule bernama pembaca daftar ternyata
     **menyimpan dan menutup transaksi**. Test yang menggolongkan rule dari namanya **gagal**.
     *(Bab 13 · 18)*
112. `[terbuka]` **butir 23** — ⛔ **bentuk tabel dan relasinya belum ada** untuk modul ini, sehingga
     seluruh AC yang menyebut *"tersimpan"* **belum punya sasaran uji di tingkat kolom**. ⚠️ Ini
     **tidak menahan spec**, tetapi **menahan tiket**. *(Out of Scope)*

---

**Jumlah AC: 114.** **AC tanpa tanda golongan: 0.** **AC ber-`[terbuka]`: 20** — butir
**15 · 28 · 49 · 50 · 58 · 70 · 73 · 75 · 78 · 82 · 87 · 88 · 89 · 90 · 91 · 102 · 104 · 112 ·
113 · 114**.

> ✅ **Dihitung dua cara, keduanya SEPAKAT pada 20**, dengan aturan kutipan yang berbeda — lihat
> kepala berkas. **AC tanpa tanda golongan: 0** pada kedua cara.
>
> ⛔ **RALAT 2026-09-19.** Kalimat lamanya **dikutip utuh, tidak dihapus**: *"**Jumlah AC: 112.**
> … **AC ber-`[terbuka]`: 23** — butir **4 · 15 · 27 · 28 · 42 · 49 · 50 · 56 · 58 · 70 · 73 ·
> 75 · 78 · 82 · 87 · 88 · 89 · 90 · 91 · 102 · 104 · 109 · 112**."* ⭐ **AC 4 · 27 · 56 · 109
> tidak lagi ber-`[terbuka]`** — keempatnya dijawab; **AC 113 · 114** baru.
>
> ⚠️ **Ke-20 AC itu bukan 20 butir register yang berbeda.** Beberapa AC merujuk butir register
> yang sama dari sudut berbeda. **Jangan menyimpulkan pemetaan satu-lawan-satu.**

---

## Out of Scope

1. **Bentuk tabel dan relasi antar-tabel.** ⛔ Modul ini **belum punya** `STRUKTUR-TABEL` maupun
   `RELASI-TABEL`. ⚠️ Keduanya **tidak menahan spec ini**, tetapi **menahan tiket** — di Claim Prop
   keduanya dikerjakan **sesudah** spec dan **sebelum** tiket.
2. **Kode Go dan React.** Spec ini menyatakan perilaku, bukan rancangan modul perangkat lunak.
3. **Modul Komite Claim Fac In.** Ia konteks tersendiri. Yang ada di sini hanyalah **kontrak
   penyerahan** — pembuatan kasus komite anak berikut lima penunjuknya.
4. **Modul Claim Prop dan Komite Claim Prop.** ⛔ Tidak disentuh, walau enam berkasnya dirujuk
   sebagai pembanding.
5. **Revisi ADR.** ⛔ Bab 16 memberi vonis; ia **tidak mengusulkan perubahan** pada satu pun ADR.
6. **Aturan peran sistem baru.** ⛔ Belum ditulis di mana pun — lihat butir register 6.
7. **Nama tabel dan kolom polis** sebagai sumber nilai penawaran — butir register 1.
8. **Penggolongan ketujuh pesan validasi** — butir register 12.
9. **Migrasi penuh lawan koeksistensi** — butir register 15.

---

## Butir `[terbuka]` — daftar penuh

✅ **Lima butir DITUTUP 2026-09-19** — **9 · 10 · 12 · 14 · 22** — dan **dua MENYEMPIT** —
**1 · 6**. ⭐ **Dua butir BARU lahir** — **24 · 25**. ⛔ Baris yang ditutup **tidak dihapus**;
ia dicoret dan tetap terbaca.

| # | Butir | Pemilik | Memblokir? |
| --- | --- | --- | --- |
| **1** | ⚠️ **MENYEMPIT** — nama tabelnya **`T_QUOTATIONDATA`** `[keputusan work owner]` 2026-09-19; ⛔ **kolomnya belum**, dan `[terverifikasi]` **tabelnya belum ada** *(0 dari 482 berkas)* — **sedang dibangun**. Teks lama: ~~Nama tabel polis dan kolomnya sebagai sumber nilai penawaran~~ | work owner + DBA | ⚠️ **ya** — AC 14 tanpa sasaran uji |
| **2** | Kenapa nomor komite ditulis **dua kali berturut-turut** | work owner | tidak |
| **3** | Apakah kolom penanda komite Fac In **sama** dengan yang dibuang di Komite Claim Prop | work owner + DBA | tidak |
| **4** | `IndexObject` lawan `ObjectIndex` — mana yang mana | work owner | tidak |
| **5** | Kelengkapan jejak kronologi, dan apakah dua catatan *"hapus yg set chronologi"* dikehendaki | work owner | tidak |
| ⭐ **6** | ⚠️ **MENYEMPIT** — **CARA menetapkannya diputuskan**: *sekali, lintas modul*, bukan per modul `[keputusan work owner]` 2026-09-19. ⛔ **Isinya masih belum ditulis di mana pun.** Teks lama: ~~Aturan peran sistem baru belum ditulis~~ | work owner + IAM | ⚠️⚠️ **ya** — AC 72 · 73 · 74 |
| **7** | Apakah panggilan konversi **tanpa saringan** dikehendaki | work owner | tidak |
| **8** | Jendela alias SQL yang disepakati — empat angka berdiri: 50 · 48 · 42 · 371 | asisten | tidak |
| ~~**9**~~ | ✅ **DITUTUP 2026-09-19** — **ikut diterima**, sebagai akibat langsung keputusan lingkungan *(butir 10)*. Teks lama: ~~Apakah dua jenis berkas video ikut diterima~~ | work owner — **sudah menjawab** | — |
| ~~⭐ **10**~~ | ✅ **DITUTUP 2026-09-19** — ⭐ **salinan PRODUKSI yang ditiru** `[keputusan work owner]`. ⚠️ Akibatnya menyentuh **seluruh modul yang sudah selesai**, yang disusun tanpa menyebut sistem sumbernya. Teks lama: ~~Rule dari lingkungan mana yang ditiru — korpus memuat tiga sistem~~ | work owner — **sudah menjawab** | — |
| **11** | Apakah rule di luar ekspor dapat melempar tiket | work owner + tim Pega | tidak |
| ~~⭐ **12**~~ | ✅ **DITUTUP 2026-09-19** — ⭐ **ketujuhnya menggagalkan penyimpanan**; nol berstatus peringatan. ⚠️ **Melahirkan butir 25.** Teks lama: ~~Pesan mana menggagalkan simpan, mana yang memperingatkan~~ | work owner — **sudah menjawab** | — |
| **13** | Apakah pesan sistem lama **menghalangi penyerahan** — ⚠️ **tidak lagi memblokir** sejak butir 12 dijawab; perilaku sistem baru ditetapkan tanpa menunggu jawabannya | work owner + tim Pega | tidak |
| ~~**14**~~ | ✅ **DITUTUP 2026-09-19** — **20 angka, 8 di belakang koma, tampilan 4** — **sama dengan Claim Prop**, sebab keduanya berbagi tabel akseptasi. Teks lama: ~~Ketelitian desimal yang dipakai modul ini~~ | work owner — **sudah menjawab** | — |
| **15** | **Migrasi penuh lawan koeksistensi** | work owner | tidak |
| **16** | Apakah perlu **penanda** pembeda baris sebelum/sesudah perbaikan | work owner | tidak |
| **17** | Arti `REQUIRED = -1` lawan `0` | work owner + tim Pega | tidak |
| **18** | Isi tiap penampung generik `CARI<n>` | work owner + DBA | tidak |
| **19** | **327 catatan pengembang belum diuji** dengan membuka rule-nya | asisten | tidak |
| **20** | Gerbang aksi tombol — medan kedua belas belum disisir | asisten | tidak |
| **21** | Perilaku bila penunjuk posisional **salah alamat** | work owner | tidak |
| ~~**22**~~ | ✅ **DITUTUP 2026-09-19 — vonis BERDIRI.** ⭐ Propertinya **benar-benar diisi**, tetapi **oleh modul saudara** yang membuat penawaran; modul ini menyalin halamannya utuh. Rantai empat mata terbaca di Bab 4. ⚠️ **Melahirkan butir 24.** Teks lama: ~~Apakah properti yang diuji tiap rule klasifikasi benar-benar diisi~~ | asisten — **sudah dikerjakan** | — |
| **23** | **Bentuk tabel dan relasinya** — ⭐ `[keputusan work owner]` **2026-09-19: DITUNDA sampai Komite Claim Fac In selesai, lalu keduanya dibuat BERSAMAAN.** ⚠️ Alasannya: `[terverifikasi]` **81 dari 112** identitas rule Komite Claim Fac In juga ada di modul ini, dan ronde 1 modul itu menemukan tabel akseptasi klaim **disentuh 24 kali di sana lawan 27 kali di sini** — keduanya **berbagi tabel**. ⛔ Membangun strukturnya terpisah mengundang **dua bentuk kolom untuk satu tabel yang sama**. ⚠️ **Yang berubah KAPAN ia dikerjakan, bukan statusnya** — ia **tetap memblokir tiket**. | asisten + DBA | ⚠️ **ya untuk tiket**, tidak untuk spec |
| ⭐ **24** | **BARU** — ketepatan klasifikasi lini **bergantung modul lain** yang mengisi halaman penawaran, dan ⛔ **nol pemeriksaan** di modul ini akan mengatakannya bila propertinya kosong. Apa yang dilakukan aplikasi dalam keadaan itu belum diputuskan | work owner | ⚠️ **ya** — AC 113 |
| ⭐ **25** | **BARU** — **berapa banyak klaim lama sudah memuat keadaan yang kini menjadi galat** *(estimasi nol-nol, estimasi melebihi nilai pertanggungan)*. Di sistem lama hanya **2** dari 7 pemeriksaan yang menghalangi; kini **ketujuhnya** | asisten + DBA | ⚠️ **ya untuk MIGRASI**, tidak untuk pembangunan |

**Ringkas:** **25 butir terdaftar** · ✅ **tertutup 5** *(9 · 10 · 12 · 14 · 22)* ·
⭐ **masih terbuka 20** · ⚠️ **memblokir 5** — butir **1 · 6 · 23 · 24 · 25**.

⚠️ **Dihitung dua cara, dan keduanya sepakat:** *(a)* mencacah baris tabel ⇒ **25** terdaftar,
**5** bertanda `DITUTUP`; *(b)* delta terhadap ringkas lama — `23 + 2 baru = 25`, dan
`memblokir 7 − 4 tertutup-atau-tak-lagi-memblokir + 2 baru = 5`.

> ⛔ **RALAT 2026-09-19.** Kalimat lamanya **dikutip utuh, tidak dihapus**: *"**Ringkas:**
> **23 butir** · ⚠️ **memblokir 7** — butir **1 · 6 · 10 · 12 · 14 · 22 · 23**."* Dan kalimat
> pembuka bab ini, yang juga sudah tidak berlaku: *"⛔ **Nol butir ditutup di berkas ini.** Spec
> hanya mendaftar."*
>
> ⚠️ **Butir 1 dan 6 TETAP memblokir walau menyempit** — yang diketahui baru **nama tabelnya**
> dan **cara menetapkan aturan peran**, bukan isinya. ⛔ Menutupnya sekarang akan menyembunyikan
> dua lubang.

---

## Further Notes

### Kesimpulan spec ini yang **PALING RAWAN SALAH**

⚠️⚠️ **Bab 5 — tabel tujuh pesan validasi.**

**Kenapa rawan:** ia memetakan pesan ke penanda **lewat kesamaan syarat gerbang**, bukan lewat
tautan yang tertulis. Langkah yang menyetel penanda tidak menyebut pesan mana pun — ia hanya
kebetulan bergerbang pada **dua variabel yang sama**. ⛔ Kalau sebuah gerbang punya sumber lain yang
belum terbaca, pemetaan itu bergeser. ⚠️ Pesan pada langkah 8 paling rawan: syaratnya muncul di
gerbang lain **tergabung dalam rantai DAN** dengan syarat tambahan.

**Yang kokoh dan tidak rawan:** ⭐ **Bab 6** — sensus tipe parameter menutup **dua arah**
*(hanya-di-A = 0, irisan = 169 = angka yang sumbernya sendiri hasilkan)*, dan ⭐ **Bab 2** — sensus
bentuk alur menutup secara aritmetika, tiap entri terbaca **wadah demi wadah**.

**Rawan kedua:** Bab 15 vonis *"nol gerbang wewenang"*. Ia bersandar pada sapuan **15 nama medan**;
⚠️ pola *"medan keenam belas"* sudah **empat kali** menggigit proyek ini.

### Dua kegagalan alat yang melahirkan aturan sensus sekarang

⭐ Dua ronde beruntun kehilangan kesimpulan karena sensus dihitung dengan **jendela yang salah**:

| Ronde | Yang runtuh | Sebabnya |
| --- | --- | --- |
| **4** | *"`Register_Flow` punya 16 bentuk, dan NOL bentuk selesai"* → lalu *"27 bentuk"* | **empat jenis entri dijumlahkan** sebagai satu; bentuk tinggal di wadah `pyShapes` saja |
| **5** | *"8 parameter tanpa tipe pasti"* — dan ia sempat **menggeser dasar vonis ADR-0003** | nama dipasangkan ke tipe **menyilang tujuh wadah**, bukan lewat kunci |

⚠️ **Keduanya lolos karena tidak ada alat yang jendelanya tetap.** Aturan penggantinya kini ada di
`CLAUDE.md` §4a: **tiap sensus dihitung dua cara, jendelanya disebut, dan bila berselisih ditulis
"belum punya data" beserta kedua angkanya** — aturan yang dipakai di seluruh berkas ini.

### Urutan pengerjaan yang disarankan

1. **Jawab tujuh butir pemblokir** — 1 · 6 · 10 · 12 · 14 · 22 · 23.
2. **`STRUKTUR-TABEL-CLAIM-FACIN.md`** dan **`RELASI-TABEL-CLAIM-FACIN.md`**.
3. **Tiket**, dengan prefactor skema uang lebih dulu — pola yang sama dengan Claim Prop.

### Yang membuat modul ini berbeda dari Claim Prop

| | Claim Prop | Claim Fac In |
| --- | --- | --- |
| kelas kerja | `WORK-CLAIMTREATY` | **`WORK-PNC`** |
| tingkat data | dua | ⭐ **tiga** — objek → item objek → penyesuaian |
| lompatan | ⛔ **nol** | ⭐ **23**, sembilan bernama lini produk |
| halaman penawaran | tidak ada pada objek kerjanya | ⭐ **ada, dan ditulis dua jalur** |
| `IsError` | **tidak pernah diset** | ⭐ **diset, dan dibaca tiga gerbang** |
| berkas alur | ruleset resmi | ⚠️ **ruleset bernama akun perorangan** |
| penunjuk ke kasus komite | tiga | ⭐ **lima** |

---

## Lampiran — Aturan baca ekspor Pega yang mengikat spec ini

⭐ Aturan di bawah **bukan hiasan**. Tiap satu lahir dari kesalahan nyata di proyek ini, dan
mengabaikannya menghasilkan angka palsu.

1. ⭐⭐ **Bendera gerbang bernilai mati mematikan GERBANGNYA, bukan LANGKAHNYA.** Langkahnya tetap
   berjalan — dan karena saringannya hilang, ia berjalan **lebih sering**, bukan tidak pernah.
   Sebuah langkah hanya mati bila **ber-remark**.
2. ⭐ **Bentuk alur dibaca dari wadah `pyShapes` saja.** Penghubung, pengubah, tiket, perute, dan
   pemberitahu tinggal di wadah lain atau **di dalam** sebuah bentuk. Menjumlahkan semuanya memberi
   angka yang **bukan jumlah bentuk**.
3. ⭐ **Identitas rule EMPAT bagian** — nama + kelas + ruleset + versi. Nama berkas **bukan bukti**,
   dan nama rule yang sama **bukan bukti**.
4. ⭐ **Kode arah 5 mengubah rantai gerbang dari DAN menjadi ATAU.** Enam titik hidup di modul ini.
5. ⭐ **Dua daftar hanya boleh dipasangkan lewat KUNCI, tidak pernah lewat urutan** — kecuali
   keduanya terbukti datang dari wadah yang sama.
6. ⭐ **Nama rule menyesatkan.** Rule bernama pembaca daftar terbukti **menyimpan dan menutup
   transaksi**.
7. **Ketiadaan tidak boleh disimpulkan sebelum seluruh medan disisir** — kalimat yang sah berbunyi
   *"tidak ketemu di medan X, Y, Z"*, bukan *"nol"*.
8. **Tiap sensus menyebut jendelanya**, dan **dihitung dua cara**; bila berselisih, tulis **"belum
   punya data"** dan cantumkan keduanya.
9. **Penyaring tidak peka huruf besar-kecil** — ejaan metode di modul ini terbukti tidak konsisten.
10. **Properti korpus dicari dengan awalan halaman**, bukan telanjang.
11. ⚠️ **Dua nilai tak berdokumen** sudah ditetapkan setara dengan yang dikenal
    `[keputusan work owner]` 2026-09-19: bendera gerbang bernilai `0` **setara kosong**, dan penanda
    perulangan `PROPERTYLIST` **setara `EMBEDDED`**.
12. ⚠️ **Aturan "satu baris transisi per langkah" punya pengecualian** — satu langkah di modul ini
    bernilai nol. Pengecualian **wajib dilaporkan**, bukan diam-diam dibulatkan.

---

## Lampiran 2 — Teks sumber yang **TIDAK dipakai**

⛔ Enam berkas sumber memuat blok RALAT yang mengutip kalimat lama. Kalimat berikut **mati** dan
**tidak menjadi bahan spec ini** — dicatat agar tidak masuk kembali lewat pembacaan ulang:

| Kalimat lama | Yang berlaku |
| --- | --- |
| *"8 parameter tanpa tipe pasti"* | ⭐ **nol berselisih** — 169 dari 169 sepakat |
| *"`Register_Flow` punya 16 bentuk"* · *"27 bentuk"* | ⭐ **8 bentuk · 10 penghubung · 2 pengubah · 7 tiket** |
| *"9 `Data-MO-Event-Exception`"* sebagai sembilan jalur | ⭐ **2 tiket bernama + 5 cangkang kosong + 2 definisi** |
| *"NOL bentuk selesai"* | ⭐ **bentuk selesai ADA** — `END52`, `Resolved-Completed` |
| *"pemanggil `SetKomiteList_ACT` dua"* | ⭐ **satu**; rule bernama nyaris sama adalah **rule lain** dengan **tiga** pemanggil |
| *"perbaiki logic `GetAllData_Act` di Go"* | ⛔ **DICABUT** — ditiru apa adanya |
| *"vonis ADR-0012 tidak dapat diputuskan"* | ⭐ **TIDAK MENYENTUH** |
| *"4 aksi lokal bernama"* | ⭐ **sembilan** |
| *"`CariHistoryClaim_SQL` 5 alias, 5 berbohong"* | ⭐ **enam alias, lima berbohong** |
| *"`CreateKMTNo_Act` empat penunjuk posisional"* | ⭐ **lima** |
| *"38 dari 60 rule `When` menguji `Quotation`"* | ⭐ **34 jalur-1 + 4 jalur-2** |
| *"rule bersama 162 · identik 158 · beda versi 4"* | ⭐ **161 · 159 · 2 beda waktu + 4 beda ruleset** |
| *"2 yang terurai + 6 = 8"* penyimpangan | ⭐ **5 terurai + 6 ADR + 1 = 12** |
| *"`pyParametersParamName` 4.071 kali"* | ⛔ **belum punya data** — 9.566 · 8.320 · 7.481 · 7.315 |

---

## Lampiran 3 — bukti berkas lain tidak disentuh

```
korpus Claim Fac In        482 berkas .xml   — nol dibuka untuk ditulis
berkas sumber claim-facin    6 berkas        — nol disunting
docs\adr\                   15 berkas        — nol disunting
.scratch\claim-prop\        27 berkas        — nol disunting
.scratch\komite-claim-prop\ 27 berkas        — nol disunting
CLAUDE.md                                    — nol disunting
```

⛔ **Nol kode Go/React · nol `CREATE TABLE` · nol DDL · nol nomor baris XML dikutip · nol nilai
rahasia disalin · nol keputusan work owner BARU dibuat · nol butir `[terbuka]` ditutup · nol ADR
direvisi.**
