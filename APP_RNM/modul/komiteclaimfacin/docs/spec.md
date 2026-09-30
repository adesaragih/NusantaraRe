# Spec — Komite Claim Fac In
## Persetujuan berjenjang atas penyesuaian klaim fakultatif masuk

> **SENSUS BERKAS INI**
>
> ⛔ **Jendelanya disebut lebih dulu: seluruh berkas ini DIKURANGI blok sensus ini sendiri**
> — **909 baris** menurut cacah baris-baru, **910** menurut pemisahan teks. ⚠️ **Selisih satu,
> sebabnya diketahui:** berkas berakhir dengan baris baru, sehingga pemisahan teks menghasilkan satu
> baris kosong penutup yang bukan baris sungguhan. ⭐ **Angka yang dipakai: 909.**
> ⭐ Blok sensus dikecualikan dengan sengaja — bila ia menghitung dirinya sendiri, angkanya
> berkejaran tanpa henti.
>
> ⭐ **Acceptance Criteria: 100** · rentang **1–100** · ⭐ **nomor ganda NOL** · **nomor lompat
> NOL** · ⭐ **AC tanpa penanda NOL** · **user story 38**
>
> `[terverifikasi]` **69** · `[keputusan work owner]` **63** · `[penyimpangan sadar]` **14** ·
> `[terbuka]` **12** · `[data DBA]` **5** · `[work owner]` **6** · `[DBA]` **3** ·
> `[pengembang Pega lama]` **1** · `[tim operasi]` **1**
> ⚠️ **52** · ⭐ **126** · ⛔ **117** · ✅ **9** · RALAT **3** · tabel **15**
>
> ⭐ **Dihitung DUA CARA:** *(a)* cacah penanda atas badan berkas sebagai satu teks utuh;
> *(b)* cacah ulang **baris demi baris**. ✅ **Keduanya sepakat pada ketujuh belas angka**, termasuk
> cacah AC.
>
> ⚠️ **Catatan jujur:** penanda `[terbuka]` terhitung **12**, sedangkan register mencatat
> **14 butir terbuka** — ⭐ jendelanya berbeda: yang satu **penanda tertulis**, yang lain **butir
> dalam register**.

---

## Cara membaca berkas ini

Berkas ini **memerikan perilaku**, bukan skema dan bukan kode. Ia disusun dari **empat ronde
grilling** atas 114 berkas ekspor Pega modul `Komite Claim FacIn`, dibandingkan bila perlu dengan
`Claim Fac In` (482), `Komite Claim Prop` (80), dan `Claim Prop` (329) — **1.005 berkas** seluruhnya.

**Penanda yang dipakai:**

| Penanda | Artinya |
| --- | --- |
| `[terverifikasi]` | Terbaca langsung dari korpus ekspor, dengan jendela sensus yang disebut |
| `[keputusan work owner]` | Diputuskan work owner, bertanggal, **bukan** simpulan asisten |
| `[penyimpangan sadar]` | Sistem baru **sengaja berbeda** dari sistem lama — alasannya selalu ditulis |
| `[terbuka]` | Belum diketahui. ⛔ **Tidak ditebak.** |
| `[data DBA]` | Jawabannya ada di basis data, bukan di korpus |
| ⭐ | Yang menentukan bentuk rancangan |
| ⚠️ | Yang mudah salah dibaca atau salah dibangun |
| ⛔ | Larangan, atau kesimpulan yang membatalkan kesimpulan lain |

⭐ **Kosakata domain** dipakai konsisten: **jenjang** *(satu tingkat pada tangga persetujuan)*,
**pemegang jenjang** *(akun yang memegang jabatan pada jenjang itu)*, **giliran** *(jenjang yang
sedang menunggu keputusan)*, **pengaju** *(orang yang mengajukan penyesuaian ke komite)*,
**ringkasan akseptasi** *(rekapitulasi nilai yang disetujui)*.

⛔ **Nama rule Pega TIDAK menjadi nama komponen sistem baru.** Ia muncul hanya **di dalam AC sebagai
rujukan bukti**.

⚠️ **Tiga RALAT lintas-ronde yang mengubah perilaku, dan wajib dibaca sebelum yang lain:**

| # | Yang dibatalkan | Yang benar |
| --- | --- | --- |
| ⛔ **1** | ronde 3 menyimpulkan *"penolakan tersimpan permanen; **penerimaan tidak**"* | ⭐ **Keduanya tersimpan.** Langkah penyimpan **diangkat ke rule pemanggil**, bergerbang *penyetuju terakhir + keputusan terima* |
| ⛔ **2** | ronde 3 menuduh Komite Claim FacIn membawa **tabel jabatan ditanam keras** | ⭐ **Itu sisa editor Pega**, bukan rule berjalan. Jalur hidupnya memakai **roster** |
| ⛔ **3** | ronde 3 menyebut label jabatan tercetak pada **dokumen yang keluar perusahaan** | ⭐ **Sasarannya catatan kronologi kasus** — jejak internal. ⚠️ Akibatnya tetap serius: **jejak audit mencatat jabatan yang keliru** |

---

## Problem Statement

Penyesuaian nilai klaim fakultatif masuk di atas kewenangan seorang penilai harus **disetujui
berjenjang** oleh komite klaim. Di sistem lama, proses itu berjalan — **tetapi tiga hal membuatnya
tidak dapat dipercaya, tidak dapat diaudit, dan tidak dapat diubah tanpa rilis.**

**Pertama — tidak ada yang benar-benar menjaga siapa boleh memutuskan.** `[terverifikasi]` Sisiran
menyeluruh atas 114 berkas menemukan **tepat satu** pemeriksaan pemilik giliran, dan ia menempel
pada **tombol Submit di layar**. ⛔ Siapa pun yang dapat memanggil lapisan layanan dapat menyimpan
keputusan untuk jenjang mana pun. ⚠️ Di modul saudaranya, `Komite Claim Prop`, **tidak ada
pemeriksaan itu sama sekali** — tidak di aktivitas, tidak di layar.

**Kedua — identitas orang dicampur dengan wewenangnya.** `[terverifikasi]` Daftar komite menyimpan
**nama orang**, bukan jabatan. Akibatnya sebuah rule menukar dua akun menjadi nama akun ketiga
**sebelum** pemeriksaan pemilik giliran dijalankan — sehingga **jejak audit mencatat orang yang
salah**. Pada modul saudaranya, jabatan yang tercetak di catatan kronologi disusun dengan
**mencocokkan ID pengguna perorangan**, dengan **nilai bawaan jabatan direksi** bagi siapa pun yang
tak dikenal.

**Ketiga — perilaku penting bersandar pada langkah yang sudah mati.** `[terverifikasi]` Pemutus
putaran persetujuan bukan pencacah jenjang, melainkan status keputusan terakhir — dan **kedua
langkah yang mengosongkan status itu ber-remark**, begitu pula penutup paksa kasus pada jalur utama.
⚠️ Modul bekerja **karena efek samping**, bukan karena aturan yang tertulis.

⚠️ **Dan satu cacat yang tak pernah ketahuan:** dua lini usaha — MBU dan Travel — diuji lewat medan
data yang **tidak pernah disalin** ke objek kerja komite. ⛔ Cabang keduanya **tidak pernah terbit**,
diam-diam, tanpa galat.

---

## Solution

Membangun ulang persetujuan berjenjang komite klaim fakultatif masuk dengan **tiga hal dipisahkan
yang di sistem lama menyatu**:

| Yang dipisahkan | Isinya | Kapan berubah |
| --- | --- | --- |
| ⭐ **Jenjang komite** | daftar **jabatan** berurutan | saat kebijakan berubah |
| ⭐ **Pemegang jabatan** | jabatan → akun | saat orang berganti |
| ⭐ **Giliran pada kasus** | jabatan **dan** akun pemegangnya, **dibekukan saat kasus dibuat** | tidak pernah |

⭐ **Dari pemisahan itu, tiga masalah di atas selesai sekaligus:** wewenang dapat ditegakkan di
lapisan layanan karena giliran punya pemilik yang pasti; jejak audit mencatat **akun sebenarnya**
karena tak ada lagi penukaran nama; dan tangga berhenti karena **keadaan tiap jenjang**, bukan
karena efek samping langkah yang kebetulan hidup.

**Yang dijanjikan kepada pengguna:**

- Penilai klaim mengajukan penyesuaian, dan tahu **siapa yang sedang memegang giliran**.
- Anggota komite melihat **hanya kasus yang menjadi gilirannya**, dan keputusannya **tidak dapat
  diambil alih orang lain**.
- Setiap keputusan meninggalkan jejak: **akun siapa, jabatan apa saat itu, kapan, dan apa
  keputusannya** — dan jejak itu **tidak berubah** ketika orangnya kelak naik jabatan.
- Satu penolakan **menutup kasus seketika**; persetujuan seluruh jenjang **menyelesaikannya**.
- Pembayaran ke kasir **tidak akan terkirim dua kali**.

---

## User Stories

**Pengajuan dan pembentukan komite**

1. Sebagai **penilai klaim**, saya ingin mengajukan penyesuaian nilai klaim ke komite, agar
   penyesuaian di atas kewenangan saya mendapat persetujuan yang sah.
2. Sebagai **penilai klaim**, saya ingin sistem menentukan sendiri **berapa jenjang** yang harus
   menyetujui, agar saya tidak perlu menebak siapa yang berwenang.
3. Sebagai **penilai klaim**, saya ingin melihat **status pengajuan saya** — jenjang keberapa yang
   sedang berjalan — agar saya dapat menjawab pertanyaan pemegang polis tanpa bertanya ke komite.
4. Sebagai **penilai klaim**, saya ingin pengajuan *tutup klaim* dan *tolak klaim* melewati
   **satu jenjang saja**, agar keputusan administratif tidak tertahan oleh tangga penuh.
5. Sebagai **penilai klaim**, saya ingin data penyesuaian dan data kutipan **tersalin lengkap** ke
   berkas komite, agar komite menilai angka yang sama dengan yang saya lihat.

**Menilai dan memutuskan**

6. Sebagai **anggota komite**, saya ingin melihat **hanya kasus yang menjadi giliran saya**, agar
   antrean saya tidak tercampur dengan jenjang lain.
7. Sebagai **anggota komite**, saya ingin melihat rekapitulasi penyesuaian **beserta totalnya dalam
   mata uang asli dan dalam rupiah**, agar saya dapat menilai besarannya.
8. Sebagai **anggota komite**, saya ingin menyetujui atau menolak dengan satu tindakan yang jelas,
   agar tidak ada keraguan apa yang saya putuskan.
9. Sebagai **anggota komite**, saya ingin menuliskan **catatan** pada keputusan saya, agar alasannya
   terbaca oleh jenjang berikutnya.
10. Sebagai **anggota komite jenjang pertama**, saya ingin menandai **usul menutup klaim** dan
    **usul mencadangkan**, agar jenjang di atas saya menilai usulan yang sudah terbentuk.
11. Sebagai **anggota komite jenjang kedua ke atas**, saya ingin **menilai usulan itu, bukan
    menggantinya**, agar yang dinilai setiap jenjang adalah hal yang sama.
12. Sebagai **anggota komite**, saya ingin keputusan saya **tidak dapat diubah** setelah tersimpan,
    agar jejaknya dapat dipertanggungjawabkan.
13. Sebagai **anggota komite**, saya ingin **tidak menerima giliran atas klaim yang saya ajukan
    sendiri**, agar tidak ada benturan kepentingan pada jenjang pertama.

**Berjalannya tangga**

14. Sebagai **pemilik proses klaim**, saya ingin kasus komite **selesai ketika seluruh jenjang
    menyetujui**, agar tidak ada persetujuan yang terlewat.
15. Sebagai **pemilik proses klaim**, saya ingin **satu penolakan menutup kasus seketika**, agar
    jenjang berikutnya tidak membuang waktu menilai yang sudah ditolak.
16. Sebagai **pemilik proses klaim**, saya ingin tangga berhenti karena **keadaan tiap jenjang**,
    bukan karena sebuah pencacah, agar daftar jenjang yang berubah di awal tidak membuatnya salah
    hitung.
17. Sebagai **administrator**, saya ingin dapat **menaikkan kasus satu jenjang** ketika pemegang
    giliran berhalangan, agar kasus tidak macet.
18. Sebagai **auditor**, saya ingin eskalasi itu **terekam** — siapa memindahkan, kapan, dari jenjang
    mana ke mana — agar perpindahan wewenang tidak tak terlihat.

**Wewenang**

19. Sebagai **pemilik proses klaim**, saya ingin **hanya pemegang giliran** yang dapat menyimpan
    keputusan, dan penolakannya terjadi **di lapisan layanan**, agar layar bukan satu-satunya
    penjaga.
20. Sebagai **auditor**, saya ingin percobaan menyimpan keputusan oleh bukan pemegang giliran
    **ditolak dan terekam**, bukan diabaikan diam-diam.
21. Sebagai **administrator kepegawaian**, saya ingin mengubah pemegang jabatan **tanpa menyentuh
    kode**, agar pergantian orang tidak menuntut rilis.
22. Sebagai **pemilik kebijakan**, saya ingin mengubah **susunan jenjang komite** dengan mengubah
    baris data, agar perubahan kebijakan tidak menuntut rilis.
23. Sebagai **auditor**, saya ingin perubahan jabatan seseorang **masuk jejak audit**, sebab jabatan
    menentukan siapa boleh menyetujui.

**Jejak dan pelaporan**

24. Sebagai **auditor**, saya ingin setiap keputusan mencatat **akun** pemutusnya, bukan alias atau
    nama orang lain.
25. Sebagai **auditor**, saya ingin jabatan yang tercatat adalah **jabatan saat keputusan diambil**,
    dan **tidak ikut berubah** ketika orangnya naik jabatan.
26. Sebagai **pembaca laporan**, saya ingin tulisan jabatan **seragam** di seluruh modul, agar orang
    yang sama tidak tampak sebagai dua orang.
27. Sebagai **pembaca laporan**, saya ingin dapat **menyaring** catatan akseptasi yang dibangun
    dengan aturan lama, agar angkanya tidak tercampur dengan yang dibangun aturan baru.

**Efek akhir**

28. Sebagai **penilai klaim**, saya ingin **nomor akseptasi terbit** begitu jenjang terakhir
    menyetujui, agar dokumen dapat dicetak.
29. Sebagai **penilai klaim**, saya ingin **ringkasan akseptasi tersimpan** sekali, agar tidak ada
    catatan ganda.
30. Sebagai **bagian keuangan**, saya ingin instruksi pembayaran terkirim ke kasir **tepat satu
    kali**, agar tidak terjadi pembayaran ganda.
31. Sebagai **bagian keuangan**, saya ingin **diberi tahu ketika pengiriman ke kasir gagal** — dan
    **hanya** ketika gagal — agar pemberitahuan itu berarti sesuatu.
32. Sebagai **penilai klaim**, saya ingin dokumen akseptasi tercetak dan surel terkirim **hanya
    setelah** keputusan terakhir, agar tidak ada dokumen setengah jadi beredar.
33. Sebagai **penerima surel**, saya ingin surel hanya terkirim **di lingkungan produksi**, agar
    pengujian tidak mengirimi mitra usaha.

**Lini usaha**

34. Sebagai **penilai klaim lini MBU**, saya ingin ringkasan akseptasi saya **terbentuk**, sama
    seperti lini lain.
35. Sebagai **penilai klaim lini Travel**, saya ingin hal yang sama.
36. Sebagai **pemilik proses klaim**, saya ingin penggolongan lini usaha membaca **data kutipan yang
    lengkap**, agar tidak ada cabang yang diam-diam tak pernah terbit.

**Jalur khusus**

37. Sebagai **penilai klaim retro fakultatif**, saya ingin jalurnya **dikenali secara eksplisit**,
    agar tidak tercampur dengan jalur biasa.
38. Sebagai **auditor**, saya ingin mengetahui bahwa jalur retro **melewati penyiapan wewenang**,
    agar hal itu dapat dinilai, bukan tersembunyi.

---

## Implementation Decisions

### ID-1 — Tiga benda dipisahkan

⭐ `[keputusan work owner]` **K12 · 2026-09-20** — Jenjang komite disusun dari **jabatan**, bukan
nama orang.

**Yang dibangun:**

| Benda | Peran |
| --- | --- |
| **Daftar induk jabatan** | kode jabatan → tulisan tampilnya |
| **Susunan jenjang komite** | jabatan apa saja, pada urutan berapa |
| **Pemegang jabatan** | kolom `jabatan` di **tabel login**, berisi **kode** |
| **Giliran pada kasus** | jenjang, jabatan, dan **akun pemegangnya**, dibekukan saat kasus dibuat |

⭐ `[keputusan work owner]` **K13 · 2026-09-20** — kolom `jabatan` di tabel login, dengan **empat
syarat**: *(1)* isinya **kode**, bukan teks bebas; *(2)* **urutan jenjang** ada di daftar terpisah;
*(3)* catatan keputusan **menyimpan salinan kode jabatan** saat itu; *(4)* perubahan kolom `jabatan`
**masuk jejak audit**.

### ID-2 — Aturan rujuk-atau-salin

⭐ `[keputusan work owner]` **K13** — aturan pembeda yang mengikat seluruh pencatatan:

> ⭐ **Identitas orang DIRUJUK · jabatan DISALIN · tulisan jabatan DIRUJUK.**

| Yang dicatat | Caranya | Alasannya |
| --- | --- | --- |
| **Akun** | rujukan | orangnya tetap orang yang sama; bila namanya berubah, catatan **memang seharusnya** ikut nama baru |
| **Jabatan** | **salinan kode** | jabatan saat itu adalah **fakta sejarah** — ⛔ tidak boleh ikut naik ketika orangnya naik jabatan |
| **Tulisan jabatan** | rujukan ke daftar induk | bila label jabatan yang sama ditulis ulang, catatan lama boleh ikut; ⛔ yang dilarang adalah catatan lama **berpindah ke jabatan berbeda** |

### ID-3 — Pemutus tangga adalah keadaan jenjang

⭐ `[keputusan work owner]` **K1 · 2026-09-20** — komite berakhir bila **seluruh jenjang
menyetujui**; satu menolak ⇒ **kasus tutup seketika**.

`[terverifikasi]` Di sistem lama, penugasan berhenti karena router **tidak menghasilkan penerima**
ketika setiap jenjang sudah memutuskan — sebuah **perilaku diam** yang tidak dapat diuji.
⚠️ `[penyimpangan sadar]` Di sistem baru, syarat berakhir **dinyatakan eksplisit** dan dapat diuji.

⛔ **Pencacah jenjang boleh ada sebagai tampilan, tetapi bukan kebenaran.** `[terverifikasi]` Daftar
jenjang dapat **berubah di awal** — pengaju dikeluarkan — dan pencacah tidak ikut tahu.

### ID-4 — Penegakan wewenang di lapisan layanan

⭐ `[keputusan work owner]` **K10 · 2026-09-20** — ADR-0014 **berdiri**. Penegakan pemilik giliran
**tetap di lapisan layanan**, bacaan **sempit**.

Lapisan layanan memeriksa **identitas akun pemanggil** terhadap **akun beku pada jenjang berjalan**
sebelum menyimpan keputusan. ⛔ Gagal ⇒ **tolak**, bukan abaikan diam-diam. ⛔ Layar boleh
menyembunyikan tombol; itu kenyamanan, **bukan** penegakan.

⭐ `[keputusan work owner]` **K3** — pembandingan memakai **identitas akun**, ⛔ bukan alias.

### ID-5 — Larangan menyetujui klaim sendiri, ditiru apa adanya

⚠️ `[keputusan work owner]` **K9 · 2026-09-20** — ⛔ **BERBEDA dari rekomendasi asisten**, dan itu
disengaja.

Saat kasus komite dibentuk, bila **pemegang jenjang pertama adalah pengaju**, ia **dikeluarkan dari
daftar**, dan **jumlah jenjang dihitung ulang**.

⛔⛔ **Akibat yang diterima dengan sadar, ditulis terang-terangan:**

> ⛔ **Pemegang jenjang KEDUA ke bawah tetap dapat menyetujui klaim yang ia ajukan sendiri.**
> Pada modul saudaranya `Komite Claim Prop`, **tidak ada pemeriksaan sama sekali** — sehingga lubang
> itu berlaku **pada jenjang mana pun**. ⭐ **Ini lubang yang dibawa masuk dengan sadar**, bukan
> kelalaian.

### ID-6 — Data kutipan disalin utuh

⭐ `[keputusan work owner]` **K2 · 2026-09-20** — ⚠️ `[penyimpangan sadar]`, dicatat sebagai
**cacat yang diperbaiki**.

`[terverifikasi]` Di sistem lama, objek kerja komite hanya menerima **dua medan** data kutipan pada
jalur penerimaan, sedangkan **sembilan** penggolong lini usaha menguji medan lain — **dua di
antaranya dipakai hidup lima kali**. ⛔ Cabang MBU dan Travel karenanya **tidak pernah terbit**.

⭐ **Di sistem baru, seluruh data kutipan yang dibutuhkan komite disalin dari kasus induk** — ⛔
bukan daftar medan bernama satu per satu, sebab **daftar bernama itulah yang melahirkan cacat ini**.

### ID-7 — Penomoran akseptasi

`[terverifikasi]` Sistem lama membangkitkan nomor lewat **satu rangkaian** yang dikunci pada **kelas
kasus induk**, bukan per lini usaha. Tujuh pembangkit per lini yang ada di korpus **ber-remark
seluruhnya**, dan sisiran menyeluruh **tidak menemukan penggantinya**.

⭐ Sistem baru **mempertahankan satu rangkaian**, dikunci pada **jenis kasus induk**.

### ID-8 — Penjaga ganda-bayar dibuat eksplisit

⭐ `[keputusan work owner]` **K5 · 2026-09-20** — penanda pembayaran **tersimpan dalam transaksi
yang sama** dengan pengiriman; panggilan kedua **ditolak**.

⚠️ `[penyimpangan sadar]` `[terverifikasi]` Sistem lama menjaganya dengan sebuah gerbang bertanda
sama dengan **tunggal** — ⛔ belum dapat dipastikan apakah itu pembandingan atau penugasan.
⭐ Keputusan ini **mengurung** butir `[terbuka]` **26**: rancangan tidak lagi bergantung pada
jawabannya.

### ID-9 — Pemberitahuan galat hanya pada kegagalan

⭐ `[keputusan work owner]` **K6 · 2026-09-20** — ⚠️ `[penyimpangan sadar]`.

`[terverifikasi]` Sistem lama mengirim pemberitahuan galat **setiap kali**, sebab gerbangnya
bertanda bendera mati; medan yang diujinya pun **salah eja**. ⭐ Sistem baru memberitahu **hanya
pada kegagalan**, dan **menamai medannya dengan benar**.

### ID-10 — Jalur Fac Retro

⚠️ `[keputusan work owner]` **K7 · 2026-09-20** — **ditiru apa adanya**. ⛔ **BERBEDA dari
rekomendasi asisten**, dan itu disengaja.

Jenis treaty tertentu menandai penyesuaian sebagai **Fac Retro**, dan penyesuaian itu **melompati
penyiapan wewenang penyetujuan**.

⛔⛔ **Dua akibat yang wajib dibaca:**

1. ⚠️ Nilai penanda jenis treaty tetap **tetapan di dalam kode** ⇒ **setiap perubahan daftar treaty
   menuntut rilis perangkat lunak**.
2. ⛔⛔ **Ada jalur yang melompati pemeriksaan wewenang.** ⭐ Karena itu spec ini **menyebutnya
   sebagai jalur tersendiri** — bukan menguburnya di dalam detail — supaya siapa pun yang membaca
   aturan **melihat jalur itu ada**.

⚠️ `[terverifikasi]` Penanda Fac Retro **dipicu setidaknya tiga jalur**, bukan hanya jenis treaty —
salah satunya **tanpa gerbang sama sekali**. `[terbuka]` **butir 28**: arti nilai penandanya, dan
apakah ketiga jalur itu memang dikehendaki.

### ID-11 — Jabatan dibaca dari data pengguna

⭐ `[keputusan work owner]` **K11 · 2026-09-20** — ⚠️ `[penyimpangan sadar]`.

⛔ Tabel jabatan yang ditanam di dalam ekspresi **tidak dibawa**, **berikut nilai bawaannya**.
⭐ Bila jabatan seseorang **tidak diketahui**, sistem **menolak** — ⛔ tidak memberi jabatan bawaan.

`[terverifikasi]` Di sistem lama, nilai bawaannya adalah sebuah **jabatan direksi**, diberikan
kepada **siapa pun yang tak dikenal**, dan jejak audit mencatatnya sebagai fakta.

### ID-12 — Penandaan data lama

⭐ `[keputusan work owner]` **K8 · 2026-09-20** — catatan akseptasi historis MBU dan Travel
**tidak dibangun ulang**, tetapi **ditandai**.

| | |
| --- | --- |
| Nama kolom penanda | **`DIBANGUN_ATURAN_LAMA`** |
| Jenis | `CHAR(1)` |
| Nilai | `'1'` / `'0'` |

⭐ **Alasannya dinyatakan:** supaya **laporan dapat menyaringnya**, bukan hanya pembaca manusia yang
kebetulan membaca catatan kaki. `[data DBA]` cacah baris terdampak.

### ID-13 — Dua kotak centang usulan

⭐ `[keputusan work owner]` **K4 · 2026-09-20** — **ditiru apa adanya**.

`[terverifikasi]` Dari **92 medan** pada layar komite sistem lama, **tepat dua** yang terkunci bagi
jenjang selain yang pertama: **usul menutup klaim** dan **usul mencadangkan**.

⭐ Aturannya sempit dan jernih: **usulan tindak lanjut dibuat pemegang jenjang pertama; jenjang
berikutnya menyetujui atau menolak usulan itu, tidak menggantinya.**

⚠️ **Jebakan alih:** pembanding di sistem lama adalah **teks**, padahal pencacahnya **bilangan**.
Pega memaafkan; ⛔ sistem baru tidak akan.

### ID-14 — Penggolong lini usaha

`[terverifikasi]` Dari **49** penggolong lini usaha di korpus modul ini, **13 dipakai hidup** dan
**36 tidak dipakai sama sekali** — terbawa ekspor karena kelasnya dibagi.

⭐ **Hanya yang 13 dialihkan.** ⛔ Yang 36 **tidak dibangun**.

---

## Testing Decisions

⭐ **Yang membuat uji ini baik:** ia menguji **perilaku yang dijanjikan kepada pengguna**, ⛔ bukan
bentuk dalamnya. Sebuah uji yang memeriksa nilai pencacah, nama rule, atau urutan langkah **akan
menolak perbaikan yang benar** dan **meloloskan kesalahan yang nyata** — ⚠️ persis yang terjadi di
sistem lama, tempat perilaku bersandar pada langkah yang kebetulan hidup.

⭐ **Jahitan uji: satu — lapisan layanan komite.** Seluruh uji di bawah masuk lewat jahitan itu,
bukan lewat layar dan bukan lewat basis data. ⭐ **Alasannya:** wewenang ditegakkan di sana *(ID-4)*,
jadi menguji dari atasnya berarti menguji penjaga yang sesungguhnya.

**Prior art:** jahitan yang sama dipakai spec `Claim Fac In` dan `Komite Claim Prop` — bab
*Testing Decisions* keduanya.

### Jahitan uji — tujuh

| # | Uji | Yang dibuktikan |
| --- | --- | --- |
| **U-1** | **Tangga penyetuju** — masukkan susunan jenjang, jalankan urutan keputusan, periksa kapan berhenti | ⛔ **Uji perilaku, bukan pencacah** *(ID-3)* |
| **U-2** | **Efek akhir terbit sekali** — penyimpanan, penomoran, pemanggilan kasir **tepat satu kali** walau tangga berjenjang banyak | efek akhir hanya oleh jenjang terakhir |
| **U-3** | **Penjaga ganda-bayar** — panggilan kedua **ditolak** *(ID-8)* | ⭐ diuji sebagai **perilaku**, bukan sebagai isi kolom |
| **U-4** | **Wewenang** — bukan pemegang giliran **tidak dapat** menyimpan keputusan, **lewat lapisan layanan** *(ID-4)* | ⛔ bukan lewat layar |
| ⚠️⚠️ **U-5** | **Pengeluaran pengaju — DUA ARAH** *(ID-5)*: pengaju di jenjang **pertama** dikeluarkan dan jumlah jenjang berkurang; pengaju di jenjang **kedua** ⛔ **TETAP DI DAFTAR** | ⭐ **Arah kedua mengunci penyimpangan sadar K9** supaya tidak "terperbaiki" diam-diam oleh pengembang berikutnya |
| **U-6** | **Penggolong lini usaha** — ketiga belas yang hidup diuji dengan data kutipan **lengkap** *(ID-6)*, termasuk **MBU dan Travel** yang dulu tak pernah terbit | membuktikan cacat lama tertutup |
| ⭐⭐ **U-7** | **Salinan jabatan** *(ID-2)* — catat sebuah keputusan, **naikkan jabatan orangnya**, lalu periksa **catatan lama TIDAK berubah** | ⛔ inilah yang membedakan **salinan** dari **rujukan** |

⚠️ **U-5 dan U-7 bukan uji fungsi.** ⭐ Keduanya **mengunci keputusan work owner** — U-5 mengunci
lubang yang sengaja dibawa, U-7 mengunci pilihan salinan yang tampak "kurang rapi" dibanding
rujukan. ⛔ **Tanpa keduanya, pengembang berikutnya akan mengubahnya karena terlihat seperti bug.**

---

## Acceptance Criteria

### A · Pembentukan kasus komite

**AC 1** `[terverifikasi]` Ketika penilai mengajukan penyesuaian ke komite, sistem membentuk **satu
kasus komite** yang tertaut ke klaim induk.

**AC 2** `[terverifikasi]` Saat pembentukan, sistem menetapkan **jumlah jenjang** dari **susunan
jenjang komite** yang berlaku.

**AC 3** `[terverifikasi]` Giliran dimulai pada **jenjang pertama**.

**AC 4** `[terverifikasi]` Jalur **tutup klaim** dan jalur **tolak klaim** selalu membentuk komite
**berjenjang satu**.

**AC 5** `[keputusan work owner]` **K13** Saat pembentukan, tiap jenjang **dibekukan** bersama
**jabatan** dan **akun pemegangnya**. Perubahan pemegang jabatan sesudah itu **tidak memindahkan
kasus**.

**AC 6** `[keputusan work owner]` **K13** Bila sebuah jabatan pada susunan jenjang **tidak memiliki
pemegang aktif**, pembentukan kasus **ditolak dengan galat yang menyebut jabatannya**.

**AC 7** `[terbuka]` **butir 39** Bila sebuah jabatan memiliki **lebih dari satu** pemegang aktif,
perilakunya **belum diputuskan**. ⭐ **Usul asisten:** tolak dengan galat yang jelas — ⛔ jangan
memilih diam-diam. ⛔ **Belum menjadi keputusan.**

**AC 8** `[keputusan work owner]` **K2** Seluruh data kutipan yang dibutuhkan penggolongan lini
usaha **tersalin dari kasus induk** ke kasus komite pada saat pembentukan.

### B · Pengeluaran pengaju

**AC 9** `[keputusan work owner]` **K9** Bila pemegang **jenjang pertama** adalah **pengaju**
penyesuaian, ia **dikeluarkan** dari susunan jenjang kasus itu.

**AC 10** `[keputusan work owner]` **K9** Sesudah pengeluaran, **jumlah jenjang berkurang satu**,
dan jenjang berikutnya **naik menjadi jenjang pertama**.

**AC 11** ⚠️ `[penyimpangan sadar]` **K9** Bila pengaju memegang jenjang **kedua atau lebih rendah**,
ia **TETAP berada di susunan jenjang** dan **tetap dapat menyetujui klaim yang ia ajukan sendiri**.
⛔ **Ini lubang yang dibawa masuk dengan sadar, bukan kelalaian.**

**AC 12** `[keputusan work owner]` **K9** Perilaku AC 11 **wajib memiliki uji yang membuktikannya**
*(U-5 arah kedua)*, agar tidak diubah tanpa keputusan work owner.

**AC 13** `[terverifikasi]` Bila pengeluaran pengaju membuat susunan jenjang menjadi **kosong**,
kasus komite **tidak dibentuk** dan pengajuan **dikembalikan** kepada penilai dengan alasannya.

### C · Berjalannya tangga

**AC 14** `[terverifikasi]` Giliran diberikan kepada **pemegang jenjang terendah yang belum
memutuskan**.

**AC 15** `[keputusan work owner]` **K1** Kasus komite **selesai** ketika **seluruh jenjang telah
menyetujui**.

**AC 16** `[keputusan work owner]` **K1** Kasus komite **tertutup seketika** ketika **satu jenjang
menolak** — ⛔ jenjang berikutnya **tidak** menerima giliran.

**AC 17** ⭐ `[keputusan work owner]` **K1** Syarat berakhir dinilai dari **keadaan tiap jenjang**,
⛔ **bukan** dari pencacah.

**AC 18** ⚠️ `[penyimpangan sadar]` **K1** `[terverifikasi]` Sistem lama berhenti karena penugasan
**tidak menghasilkan penerima** — perilaku diam yang tak dapat diuji. ⭐ Sistem baru menyatakan
syaratnya eksplisit.

**AC 19** `[terverifikasi]` Sebuah jenjang **hanya dapat memutuskan sekali**. Percobaan kedua
**ditolak**.

**AC 20** `[terverifikasi]` Administrator dapat **menaikkan kasus satu jenjang** ketika pemegang
giliran berhalangan. ⛔ Menurunkan jenjang **dilarang**.

**AC 21** `[terverifikasi]` Eskalasi **terekam**: siapa memindahkan, kapan, dari jenjang mana ke
jenjang mana.

**AC 22** `[terverifikasi]` Jenjang yang **dilewati** oleh eskalasi tercatat **tanpa keputusan** —
⛔ bukan tercatat sebagai menyetujui.

### D · Wewenang

**AC 23** ⭐ `[keputusan work owner]` **K10** Hanya **akun beku pada jenjang berjalan** yang dapat
menyimpan keputusan. Penolakannya terjadi **di lapisan layanan**.

**AC 24** `[keputusan work owner]` **K10** Percobaan menyimpan oleh akun lain **ditolak dengan
galat**, ⛔ bukan diabaikan diam-diam.

**AC 25** `[keputusan work owner]` **K10** Percobaan yang ditolak **terekam** di jejak audit.

**AC 26** `[keputusan work owner]` **K3 · K12** Pembandingan memakai **identitas akun**. ⛔ Tidak ada
penukaran nama sebelum pembandingan.

**AC 27** ⚠️ `[penyimpangan sadar]` **K12** `[terverifikasi]` Sistem lama menukar **dua akun** menjadi
akun ketiga **sebelum** pemeriksaan pemilik giliran. ⛔ **Penukaran itu tidak dibawa.**

**AC 28** `[keputusan work owner]` **K10** Layar boleh menyembunyikan tindakan yang tak berwenang,
⛔ **tetapi itu bukan penegakan** — penegakan tetap di lapisan layanan.

**AC 29** `[terbuka]` **butir 6** Daftar jabatan dan jabatan mana yang membentuk susunan jenjang
**belum diketahui**. ⭐ **Dikurung sebagai titik sambung bernama:** aturan wewenang di atas dapat
ditulis dan diuji tanpa isinya; ⛔ **pembangunannya menunggu**.

### E · Layar dan penyuntingan

**AC 30** `[keputusan work owner]` **K4** Pemegang **jenjang pertama** dapat menyunting **usul
menutup klaim** dan **usul mencadangkan**.

**AC 31** `[keputusan work owner]` **K4** Pemegang jenjang **kedua ke atas** melihat kedua usulan itu
**terkunci**.

**AC 32** `[terverifikasi]` Medan lain pada layar komite **tidak terkunci** oleh aturan AC 31.
⚠️ Dari 92 medan, **tepat dua** yang terkunci.

**AC 33** ⚠️ `[terverifikasi]` **Jebakan alih:** pembanding di sistem lama **bertipe teks**, padahal
pencacahnya **bilangan**. ⛔ Sistem baru membandingkan **bilangan dengan bilangan**.

**AC 34** `[terverifikasi]` Anggota komite dapat menuliskan **catatan** pada keputusannya.

**AC 35** `[terverifikasi]` Layar menampilkan rekapitulasi penyesuaian **dalam mata uang asli** dan
**total dalam rupiah**.

### F · Penyimpanan keputusan

**AC 36** `[terverifikasi]` Keputusan tiap jenjang tersimpan bersama: **akun** pemutus, **salinan
kode jabatan** saat itu, **waktu**, dan **keputusannya**.

**AC 37** ⭐ `[keputusan work owner]` **K13** Jabatan tersimpan sebagai **salinan kode**. ⛔ Ketika
orangnya kelak naik jabatan, **catatan lama tidak berubah**.

**AC 38** `[keputusan work owner]` **K13** **Akun** tersimpan sebagai **rujukan**. Bila nama tampil
orangnya berubah, catatan **ikut nama baru** — ⭐ itu memang dikehendaki.

**AC 39** `[keputusan work owner]` **K13** **Tulisan jabatan** dirujuk dari **daftar induk**.
⛔ Catatan lama **tidak boleh berpindah ke jabatan yang berbeda**.

**AC 40** `[keputusan work owner]` **K11** Bila jabatan pemutus **tidak diketahui**, penyimpanan
keputusan **ditolak**. ⛔ **Tidak ada jabatan bawaan.**

**AC 41** ⚠️ `[penyimpangan sadar]` **K11** `[terverifikasi]` Sistem lama memberi **jabatan direksi**
sebagai nilai bawaan kepada siapa pun yang tak dikenal. ⛔ **Tidak dibawa.**

**AC 42** `[terverifikasi]` Keputusan yang tersimpan **tidak dapat diubah**.

### G · Efek akhir

**AC 43** ⭐ `[terverifikasi]` Efek akhir terbit **hanya** ketika **jenjang terakhir menyetujui** —
⛔ tidak pada jenjang mana pun sebelumnya.

**AC 44** `[terverifikasi]` **Ringkasan akseptasi disusun**, lalu **disimpan**. ⭐ Penyusunan dan
penyimpanan adalah dua langkah, dan penyimpanan terjadi **sekali**.

**AC 45** `[terverifikasi]` ⛔ **RALAT terhadap ronde 3.** Kalimat lama dikutip: *"penolakan tersimpan permanen;
**penerimaan tidak**."* ⭐ **Yang benar: keduanya tersimpan** — langkah penyimpan **diangkat ke rule
pemanggil**, bergerbang *keputusan terima* **dan** *jenjang terakhir*.

**AC 46** `[terverifikasi]` **Nomor akseptasi dibangkitkan** setelah ringkasan tersimpan.

**AC 47** `[terverifikasi]` **Dokumen akseptasi dicetak** setelah nomor terbit.

**AC 48** `[terverifikasi]` **Surel pemberitahuan terkirim** setelah dokumen siap.

**AC 49** `[terverifikasi]` Surel terkirim **hanya di lingkungan produksi**.

**AC 50** `[terverifikasi]` **Instruksi pembayaran terkirim ke kasir** setelah efek di atas.

**AC 51** ⚠️ `[terverifikasi]` Kegagalan pada salah satu efek akhir **tidak boleh** membatalkan efek
yang sudah terbit tanpa jejak. ⭐ Sistem lama punya **nol jalur kegagalan pada 25 langkah
penyimpanan**; ⛔ sistem baru **wajib** punya.

### H · Penomoran akseptasi

**AC 52** `[terverifikasi]` Nomor akseptasi berasal dari **satu rangkaian**, dikunci pada **jenis
kasus induk** — ⛔ bukan per lini usaha.

**AC 53** `[terverifikasi]` Tujuh pembangkit nomor **per lini usaha** yang ada di korpus
**ber-remark seluruhnya**, dan sisiran menyeluruh **tidak menemukan penggantinya**.
⛔ **Tidak dialihkan.**

**AC 54** `[data DBA]` **butir 31** Perilaku prosedur tersimpan yang membangkitkan nomor **belum
terbaca dari korpus**. ⭐ Sistem baru **wajib memerikan sendiri** aturan pembangkitannya, tidak
mewarisi yang tak terbaca.

**AC 55** `[terverifikasi]` Nomor akseptasi yang sudah terbit **tidak dibangkitkan ulang**.

### I · Penggolongan lini usaha

**AC 56** `[terverifikasi]` **Tiga belas** penggolong lini usaha dipakai hidup di modul ini dan
**dialihkan**.

**AC 57** `[terverifikasi]` **Tiga puluh enam** penggolong lainnya **tidak dipakai sama sekali** —
⛔ **tidak dialihkan**.

**AC 58** ⭐ `[keputusan work owner]` **K2** Penggolongan membaca **data kutipan yang lengkap**,
sehingga **seluruh** cabang lini usaha dapat terbit.

**AC 59** ⚠️ `[penyimpangan sadar]` **K2** `[terverifikasi]` Di sistem lama, cabang **MBU** dan
**Travel** menguji medan yang **tidak pernah disalin** ke objek kerja komite pada jalur penerimaan —
⛔ **keduanya tidak pernah terbit**, diam-diam, tanpa galat. ⭐ **Cacat ini diperbaiki.**

**AC 60** `[terverifikasi]` Penggolongan yang memakai hubungan **ATAU** tetap berperilaku **ATAU**,
⚠️ bukan **DAN**.

### J · Jalur kasir

**AC 61** ⭐ `[keputusan work owner]` **K5** Instruksi pembayaran terkirim ke kasir **tepat satu
kali** per akseptasi.

**AC 62** `[keputusan work owner]` **K5** Penanda "sudah terkirim" tersimpan **dalam transaksi yang
sama** dengan pengirimannya.

**AC 63** `[keputusan work owner]` **K5** Panggilan kedua atas akseptasi yang sama **ditolak**.

**AC 64** ⚠️ `[penyimpangan sadar]` **K5** `[terverifikasi]` Penjaga ganda-bayar sistem lama
bersandar pada gerbang bertanda sama dengan **tunggal**, ⛔ yang **belum dapat dipastikan**
pembandingan atau penugasan. ⭐ Sistem baru **tidak bergantung pada tafsir itu**.

**AC 65** `[terbuka]` **butir 26** ⭐ **Dikurung oleh AC 64.** Jawaban atas tanda sama dengan tunggal
**tetap dibutuhkan** — ⛔ bukan untuk merancang, melainkan untuk **memeriksa apakah data lama pernah
terbayar dua kali**.

**AC 66** ⭐ `[keputusan work owner]` **K6** Pemberitahuan galat kasir terkirim **hanya ketika
pengiriman gagal**.

**AC 67** ⚠️ `[penyimpangan sadar]` **K6** `[terverifikasi]` Sistem lama mengirimnya **setiap kali**,
sebab gerbangnya bertanda bendera mati. ⛔ **Tidak ditiru.**

**AC 68** ⚠️ `[terverifikasi]` Medan yang diuji gerbang itu **salah eja** di sistem lama. ⭐ Sistem
baru **menamainya dengan benar**.

**AC 69** `[terbuka]` **butir 27** Apakah pemberitahuan galat sistem lama **benar-benar sampai** ke
seseorang **belum diketahui**. ⚠️ Bila tidak, maka penanganan galat kasir sistem lama **sebenarnya
tidak ada**.

**AC 70** `[terverifikasi]` Kegagalan pengiriman ke kasir **tidak membatalkan** akseptasi yang sudah
tersimpan; ⭐ ia **dicatat** dan **dapat diulang**.

### K · Jalur Fac Retro

**AC 71** ⚠️ `[keputusan work owner]` **K7** Penyesuaian bertanda **Fac Retro** mengikuti **jalur
tersendiri** yang **melompati penyiapan wewenang penyetujuan**.

**AC 72** ⚠️ `[penyimpangan sadar]` **K7** ⛔ **BERBEDA dari rekomendasi asisten**, yang mengusulkan
menjadikan penanda jenis treaty sebagai data. ⛔ **Ditiru apa adanya, dengan sadar.**

**AC 73** ⚠️ `[keputusan work owner]` **K7** Nilai penanda jenis treaty tetap **tetapan di dalam
kode** ⇒ ⛔ **setiap perubahan daftar treaty menuntut rilis perangkat lunak.**

**AC 74** ⛔⛔ `[keputusan work owner]` **K7** Spec ini **menyebut jalur Fac Retro secara terbuka
sebagai jalur yang melompati pemeriksaan wewenang**, ⭐ supaya siapa pun yang membaca aturan
**melihat jalur itu ada** — ⛔ bukan menemukannya setelah sesuatu terjadi.

**AC 75** ⚠️ `[terverifikasi]` Penanda Fac Retro **dipicu setidaknya tiga jalur**, salah satunya
**tanpa gerbang sama sekali**. `[terbuka]` **butir 28** — arti nilai penandanya dan apakah ketiga
jalur itu dikehendaki.

**AC 76** ⚠️ `[terverifikasi]` Medan jenis treaty di sistem lama **mencampur kode master dengan teks
harfiah**. ⛔ Hal itu **wajib beres** sebelum medan itu boleh menjadi kolom bernilai kode.

### L · Jabatan, tabel login, jejak audit

**AC 77** ⭐ `[keputusan work owner]` **K13** Jabatan seseorang disimpan sebagai **kolom di tabel
login**, berisi **kode**.

**AC 78** `[keputusan work owner]` **K13** Kode itu menunjuk **daftar induk jabatan**, yang memuat
tulisan tampilnya. ⛔ **Bukan teks bebas.**

**AC 79** ⚠️ `[terverifikasi]` Di sistem lama, orang yang sama tertulis dengan **jabatan berbahasa
berbeda** antar modul, dan **dibaca lewat medan yang berbeda**. ⭐ Kode tunggal menutup keduanya.

**AC 80** ⭐ `[keputusan work owner]` **K13** **Susunan jenjang komite** — jabatan apa saja, pada
urutan berapa — disimpan **terpisah** dari tabel login.

**AC 81** `[keputusan work owner]` **K13** Mengubah susunan jenjang adalah **mengubah baris data**,
⛔ **bukan rilis**.

**AC 82** ⭐ `[keputusan work owner]` **K13** Perubahan kolom `jabatan` **masuk jejak audit**: siapa
mengubah, kapan, dari kode apa ke kode apa.

**AC 83** ⚠️ `[keputusan work owner]` **K13** ⭐ **Alasan AC 82 dinyatakan:** begitu jabatan
menentukan siapa boleh menyetujui, ⛔ **mengubah satu baris berarti mengubah wewenang.** Sejalan
**ADR-0007** dan **ADR-0014**.

**AC 84** `[keputusan work owner]` **K12** Roster komite berbasis **jabatan**. ⛔ Nama orang **tidak
menjadi kunci** di mana pun.

**AC 85** ⭐ `[keputusan work owner]` **K12** ⭐ **Alasan dinyatakan work owner:** pemetaan akun
menjadi nama orang lain ada **karena persetujuan komite dikirim memakai nama**. ⭐ Memakai jabatan
**menghapus akarnya**, bukan menambal gejalanya.

### M · Catatan kronologi

**AC 86** `[terverifikasi]` Setiap keputusan komite menambah **catatan kronologi** pada kasus.

**AC 87** `[terverifikasi]` ⛔ **RALAT terhadap ronde 3.** Kalimat lama dikutip: *"menandatangani **dokumen akseptasi**
… **Dokumen itu keluar perusahaan**."* ⭐ **Yang benar: sasarannya catatan kronologi kasus** —
jejak internal, ⛔ bukan dokumen akseptasi.

**AC 88** ⚠️ `[terverifikasi]` **Akibat yang benar, dan tetap serius:** jejak audit sistem lama
**mencatat jabatan yang keliru**. ⭐ Menyentuh **ADR-0007**.

**AC 89** `[keputusan work owner]` **K13** Catatan kronologi disusun dari **medan tersimpan** — akun,
kode jabatan, keputusan, waktu — ⛔ **bukan disimpan sebagai kalimat jadi**.

**AC 90** ⭐ `[keputusan work owner]` **K13** ⭐ **Alasan AC 89 dinyatakan:** justru karena sistem lama
menyimpan **kalimat jadi**, orang yang sama tercetak berbeda antar modul, dan nilai bawaan diberikan
kepada yang tak dikenal. ⛔ Dengan medan terpisah, keduanya **tidak mungkin terjadi**.

### N · Data lama

**AC 91** ⭐ `[keputusan work owner]` **K8** Catatan akseptasi historis MBU dan Travel **tidak
dibangun ulang**.

**AC 92** `[keputusan work owner]` **K8** Catatan itu **ditandai** dengan kolom
**`DIBANGUN_ATURAN_LAMA`**, `CHAR(1)`, bernilai `'1'` / `'0'`.

**AC 93** ⭐ `[keputusan work owner]` **K8** ⭐ **Alasan dinyatakan:** penandanya **kolom, bukan
catatan prosa**, supaya **laporan dapat menyaringnya** — ⛔ bukan hanya pembaca manusia yang
kebetulan membaca catatan kaki.

**AC 94** `[data DBA]` **butir 29** Cacah baris terdampak **belum diketahui**.

**AC 95** `[terverifikasi]` Catatan baru **tidak** menerima penanda itu.

### O · Batas dan titik sambung

**AC 96** `[terbuka]` **butir 1** `[data DBA]` Kolom tabel data kutipan **belum diketahui**.
⭐ **Dikurung:** ia menyentuh **bentuk tabel**, ⛔ bukan perilaku — dan tahap struktur tabel
**belum dimulai** sesuai `[keputusan work owner]` 2026-09-19.

**AC 97** `[terbuka]` **butir 6** Isi daftar jabatan **belum diketahui**. ⭐ **Dikurung sebagai titik
sambung bernama** — lihat AC 29.

**AC 98** `[terverifikasi]` Modul ini **mengunci kasus induk** sebelum mengerjakan apa pun.

**AC 99** `[terverifikasi]` Empat jalur masuk dikenali: **penyesuaian**, **tolak**, **tutup klaim**,
dan **survey**. ⛔ Jalur **survey** tidak dialihkan — lihat *Out of Scope*.

**AC 100** `[terverifikasi]` Tidak ada jalur masuk kelima. ⭐ Sisiran menyeluruh membuktikannya.

---

## Out of Scope

| # | Yang tidak dibangun | Alasannya |
| --- | --- | --- |
| **1** | ⭐ **Pintu belakang berbasis teks jabatan pengembang** | `[keputusan work owner]` 2026-09-19 — ⛔ tidak dibawa. `[terverifikasi]` gerbangnya memotong dua syarat lain, termasuk satu-satunya pemeriksaan pemilik giliran |
| **2** | ⭐ **Jalur survey** | ⛔ `[terverifikasi]` rule-nya **tidak ada di korpus mana pun** — hanya pemanggilnya yang ada, dan itu pun ber-remark |
| **3** | ⭐ **36 penggolong lini usaha** | ⛔ tidak dipakai modul ini; terbawa ekspor karena kelasnya dibagi |
| **4** | ⭐ **Tujuh pembangkit nomor per lini usaha** | ⛔ ber-remark seluruhnya; penggantinya **satu rangkaian** yang dikunci jenis kasus induk |
| **5** | ⭐ **Penaik pencacah pada penugasan** | ⛔ ber-remark; ⭐ penggerak tangga adalah **keadaan tiap jenjang** |
| **6** | ⭐ **Tabel jabatan ditanam keras, berikut nilai bawaan jabatan direksi** | `[keputusan work owner]` **K11** — jabatan dibaca dari data pengguna |
| **7** | ⭐ **Pemetaan akun menjadi nama orang lain** | `[keputusan work owner]` **K12** — ⛔ menghapus jejak siapa yang benar-benar memutuskan |
| **8** | ⭐ **Pembangunan ulang catatan akseptasi lama MBU/Travel** | `[keputusan work owner]` **K8** — angkanya mungkin sudah dilaporkan keluar; ⭐ ditandai, bukan dibangun ulang |
| ⚠️ **9** | ⚠️ **Larangan menyetujui klaim sendiri untuk jenjang kedua ke atas** | ⛔ **Sengaja TIDAK dibangun** — `[keputusan work owner]` **K9**. ⭐ Ditulis di sini **agar tidak terbaca sebagai kelalaian** |
| **10** | ⭐ **Struktur tabel dan relasinya** | `[keputusan work owner]` 2026-09-19 — dikerjakan **bersamaan dengan Claim Fac In**, sesudah spec ini |
| **11** | ⭐ **Perbaikan modul saudara** *(`Komite Claim Prop`, `Claim Prop`, `Claim Fac In`)* | ⭐ temuan lintas modul **ditunjuk**, ⛔ tidak dikerjakan di sini |

---

## Butir `[terbuka]` — daftar penuh

⭐ **14 butir.** ⛔ **Nol ditutup oleh spec ini** — spec **mencatat**, tidak memutuskan.
⛔ **Nomor 18 dibatalkan dan tidak dipakai ulang.**

| # | Butir | Pemilik | Memblokir? |
| --- | --- | --- | --- |
| ⛔ **1** | Kolom tabel data kutipan | `[DBA]` | ⭐ **tahap struktur tabel**, bukan spec |
| ⛔ **6** | **Isi daftar jabatan dan susunan jenjang** | `[work owner]` | ⛔ **YA — PEMBANGUNAN** |
| **11** | *(dibawa dari ronde 1)* | asisten | tidak |
| **12** | *(dibawa dari ronde 1)* | asisten | tidak |
| **13** | *(dibawa dari ronde 1)* | asisten | tidak |
| **15** | Asal-usul ambang nilai pada pita wewenang | `[work owner]` | tidak |
| **25** | Bentuk jalur administratif pengganti pintu belakang | `[work owner]` | tidak |
| **26** | Tanda sama dengan tunggal — pembandingan atau penugasan | `[pengembang Pega lama]` | ⭐ **dikurung K5** — untuk memeriksa **data lama** |
| **27** | Apakah pemberitahuan galat kasir benar-benar sampai | `[tim operasi]` | tidak |
| **28** | Arti penanda jenis treaty, dan ketiga jalur pemicu Fac Retro | `[work owner]` | tidak |
| **29** | Cacah baris MBU/Travel terdampak | `[DBA]` | tidak |
| **30** | Gerbang yang membuka kembali penyetuju terakhir | asisten → `[work owner]` | tidak |
| **31** | Isi prosedur tersimpan pembangkit nomor | `[DBA]` | tidak |
| ⚠️ **39** | **Satu jabatan, lebih dari satu pemegang — siapa menerima giliran** | `[work owner]` | tidak |

⭐ **Satu pemblokir tersisa: butir 6**, dan ia menahan **PEMBANGUNAN**, ⛔ bukan penulisan.

**Tiga butir dibawa terbuka ke dalam spec, masing-masing terkurung:**

| # | Kalimat pengurungnya |
| --- | --- |
| **6** | ⭐ Aturan wewenang dapat **ditulis dan diuji lengkap** tanpa mengetahui jabatan apa saja yang ada. Daftarnya menjadi **titik sambung bernama**, ⛔ bukan lubang di tengah kalimat |
| **26** | ⭐ **K5 membuat penjaga ganda-bayar eksplisit**, sehingga rancangan **tidak bergantung** pada jawabannya. Jawabannya dipakai untuk **memeriksa data lama** |
| **1** | ⭐ Spec memerikan **perilaku**, bukan **skema**. ⛔ Tak satu pun AC berubah bunyinya karena daftar kolom. Ia milik **tahap struktur tabel** |

---

## Further Notes

### ⚠️ Peringatan kredensial

`[terverifikasi]` Layanan pengiriman ke kasir **berautentikasi**, dan kredensialnya **beredar di
dalam berkas ekspor Pega**. ⛔ **Nilainya tidak dibaca, tidak dicetak, dan tidak disalin ke berkas
mana pun sepanjang empat ronde grilling maupun spec ini.**

⭐ **Karena ia beredar di dalam ekspor, ia WAJIB diputar ulang sesudah alih** — ⚠️ itu tanggung jawab
**tim pemilik layanan**, bukan tim alih.

### ⭐ Penunjuk lintas modul — bukan pekerjaan spec ini

`[terverifikasi]` Dua temuan menyentuh modul di luar spec ini:

| Temuan | Modul terdampak |
| --- | --- |
| **Tabel jabatan ditanam keras**, berikut nilai bawaan jabatan direksi | `Komite Claim Prop` |
| **Gerbang berakun-ditanam pada pengiriman surel** | `Claim Fac In` · `Claim Prop` |
| ⛔ **Tidak ada pemeriksaan pemilik giliran sama sekali** | `Komite Claim Prop` |

⭐ **Butir yang sama WAJIB masuk `utang-lintas-modul.md` milik modul-modul itu.**
⛔ **Spec ini tidak menyunting berkas modul lain.**

### ⭐ Aturan baca yang mengikat, lahir dari kesalahan asisten

⛔ **Ini bukan catatan kaki — ini aturan.** ⚠️ Ia lahir karena instrumen asisten **dua kali**
menghasilkan kesimpulan yang keliru, dan koreksinya mengubah temuan.

| Aturan | Sebabnya |
| --- | --- |
| ⭐ **Elemen di bawah `pyExpressionGadget` / `PegaGadget-ExpressionBuilder` adalah SISA EDITOR**, bukan rule berjalan — ⛔ wajib dikeluarkan dari setiap sensus gerbang | ⚠️ Tanpa pemisahan itu, cacah gerbang **melebih lebih dari dua kali lipat** |
| ⭐ **Sebelum menyatakan sesuatu TIDAK ADA, sebut tag mana saja yang disisir — dan sisir lebih dari satu** | ⚠️ Penyaring bertag tunggal **dua kali** menghasilkan tuduhan palsu |
| ⭐ **Pernyataan tentang MODUL wajib bersandar pada sisiran tingkat MODUL** | ⛔ Membaca isi dua rule lalu menyimpulkan tentang modul **membatalkan satu vonis utuh** |
| ⭐ **Tiap angka dihitung dua cara, jendelanya disebut**; beda ⇒ **"belum punya data"** + keduanya | ⚠️ Kepala sensus yang ditulis **sebelum** berkasnya jadi pernah salah pada **delapan dari sembilan** angka |

### ⭐ Yang membuat modul ini sulit dibaca, dan pelajarannya

`[terverifikasi]` Modul ini bekerja **karena efek samping**, bukan karena aturan tertulis: pemutus
tangga adalah status yang **tak pernah dikosongkan**; penugasan berhenti karena router **tidak
menghasilkan penerima**; penanda "sudah memutuskan" dipasang **justru karena gerbangnya mati**.

⭐ **Sistem baru menyatakan ketiganya secara eksplisit.** ⚠️ Itu membuat perilakunya **terlihat
berbeda di kode**, padahal **sama di mata pengguna** — dan itulah sebabnya **U-1 menguji perilaku,
bukan pencacah**.

### ⭐ Kesiapan

⭐ **Vonis to-spec: SIAP**, diambil 2026-09-20 sesudah butir 35 tertutup oleh **K10**.
⛔ **Pembangunan menunggu butir 6.**

---

## Lampiran — Aturan baca ekspor Pega yang mengikat spec ini

⭐ Seluruh `[terverifikasi]` di berkas ini bersandar pada aturan berikut. ⛔ Sebuah pembacaan yang
melanggarnya **tidak sah**.

| # | Aturan |
| --- | --- |
| **1** | Bentuk alur hanya ada di penampung bentuk; penampung lain adalah **saudaranya**, bukan isinya. ⭐ Menjumlahkan seluruhnya menghasilkan angka yang **bukan** cacah bentuk |
| **2** | ⭐ Bendera gerbang mati mematikan **GERBANG**, bukan **LANGKAH** — ⚠️ langkahnya justru berjalan **lebih sering** |
| **3** | ⭐ Kode arah **5** mengubah rantai **DAN** menjadi **ATAU** |
| **4** | ⭐ Identitas rule ada **EMPAT bagian**: nama · kelas · ruleset · versi. ⛔ `md5` berkas **bukan alat yang sah** |
| **5** | ⭐ Rujukan penggolong dapat tersimpan di **medan ekspresi** dan **medan rule sebelumnya**, ⛔ bukan hanya di medan nama-penggolong |
| **6** | ⭐ Elemen di bawah penampung *expression builder* adalah **sisa editor**, ⛔ bukan rule berjalan |
| **7** | ⭐ Penampung generik bernomor **tidak berarti apa-apa sendirian** — yang memberinya arti adalah **halaman pembawanya** |
| **8** | ⭐ Penyaring properti **tidak boleh mewajibkan awalan halaman** — banyak rule menguji properti relatif |
| **9** | ⭐ Hanya penampung parameter **langsung di bawah akar** yang memuat parameter rule itu sendiri. ⛔ Memasangkan dua daftar lintas penampung menghasilkan konflik palsu |
| **10** | ⭐ Langkah bersarang dapat mencapai **enam tingkat**; penelusur yang tidak menelusuri anak **kehilangan langkah** |

---

## Lampiran 2 — Teks sumber yang TIDAK dipakai

⛔ **Yang sengaja tidak menjadi bahan spec ini:**

| Yang tidak dipakai | Alasannya |
| --- | --- |
| **Nilai kredensial** pada rule layanan kasir | ⛔ dilarang dibaca, dicetak, dan disalin |
| **Nama akun perorangan** yang ditanam di gerbang | ⭐ terdaftar utuh di `keputusan-sebelum-to-spec.md` §2b atas perintah work owner; ⛔ **tidak diulang di spec**, sebab spec memerikan **aturan**, bukan **daftar orang** |
| **Nomor baris berkas XML** | ⛔ berubah setiap ekspor ulang |
| **Isi langkah Java** | ⭐ bentuknya disebut *(penyusun muatan)*, ⛔ isinya tidak disalin |
| **Ekspresi di bawah penampung *expression builder*** | ⛔ **sisa editor**, bukan rule berjalan — lihat Lampiran aturan **6** |
| **Kesimpulan modul lain** | ⭐ dirujuk boleh, ⛔ diturunkan tidak |

---

## Lampiran 3 — bukti berkas lain tidak disentuh

| Berkas | Keadaan |
| --- | --- |
| `komite-claim-facin\grilling-ronde-1.md` | ✅ md5 tidak berubah |
| `komite-claim-facin\grilling-ronde-2.md` | ✅ md5 tidak berubah |
| `komite-claim-facin\grilling-ronde-3.md` | ✅ md5 tidak berubah |
| `komite-claim-facin\grilling-ronde-4.md` | ✅ md5 tidak berubah |
| `komite-claim-facin\keputusan-sebelum-to-spec.md` | ✅ md5 tidak berubah |
| `claim-facin\spec.md` · `claim-prop\spec.md` · `komite-claim-prop\spec.md` | ✅ utuh |
| `docs\adr\` — 15 berkas, termasuk ADR-0014 | ✅ nol disunting |
| `OUTPUT_HASIL_RNM\CLAUDE.md` | ✅ nol disunting |
| korpus 114 · 482 · 80 · 329 | ✅ nol disunting |

⛔ Kode Go/React **NOL** · `CREATE TABLE` **NOL** · DDL **NOL** · nomor baris XML **NOL** ·
nilai kredensial **NOL** · tiket **NOL** · revisi ADR **NOL** · butir `[terbuka]` ditutup **NOL** ·
keputusan work owner baru **NOL**.
