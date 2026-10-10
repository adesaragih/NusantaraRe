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
>
> ⛔ **RALAT 10-10-2026, dalam blok yang sama.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⛔ **Jendelanya disebut
> lebih dulu: seluruh berkas ini DIKURANGI blok sensus ini sendiri** — **909 baris** menurut cacah baris-baru, **910**
> menurut pemisahan teks."* → semua
> angka di blok ini adalah cacah **sebagaimana tertulis 20-09-2026**. Blok RALAT 10-10-2026 (prompt work owner tahap 2
> §7 butir 2) yang ditambahkan sesudahnya **tidak ikut dicacah** dan **tidak mencabut satu tanda pun**; kalimat lama
> tetap utuh di tempatnya. Yang mengikat untuk pembangunan: XML Komite Claim FacIn, keputusan KCF-01–04 (10-10-2026),
> dan pola Komite Claim Prop / Non Prop, sebagaimana dirinci di setiap blok RALAT.

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

> ⛔ **RALAT 10-10-2026.** Istilah lamanya dikutip utuh, tidak dihapus: *"**pemegang jenjang** *(akun yang memegang
> jabatan pada jenjang itu)*"* → sejak **KCF-01** jenjang dipegang **workbasket** (`KOMITE_OPERATORID` = workbasket roster
> FACIN), dan "pemegang jenjang" = **setiap anggota** workbasket itu; tingkat 1 = anggota `ReasClaimSPVA` atau
> `ReasClaimSPVB`. Istilah **pengaju** tetap dipakai, tetapi tidak ada aturan pengeluaran pengaju (K9 gugur, RALAT ID-5).

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

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"`[terverifikasi]` Daftar komite menyimpan
> **nama orang**, bukan jabatan. Akibatnya sebuah rule menukar dua akun menjadi nama akun ketiga **sebelum** pemeriksaan
> pemilik giliran dijalankan — sehingga **jejak audit mencatat orang yang salah**."* →
>
> - Roster `EMAILKOMITE` STS_KLAIM FACIN menyimpan **JABATAN** dan **akun orang** di `OPERATOR_ID` (DEV 10-10-2026: ID
>   1 / 2 / 10 / 3 / 4 = DEGREE 1–5, `OPERATOR_ID` masih akun orang), dan `KomiteID` tangga = `OPERATOR_ID` itu
>   (`ApprovalKomite_Act` L7.1). **KCF-01** mengganti `OPERATOR_ID` dengan **workbasket** (UPDATE baris yang ada menurut
>   DEGREE, JABATAN dan LIMIT tetap).
> - Penukaran akun (`SetProteksiSubmiteKomite` L1, pemetaan akun orang tertulis mati) hanya dipakai gerbang tombol Submit
>   (L2.1). Langkah pencatat memakai akun pelaku sebenarnya: `KomitePost_Adjustment` S20 `InsertHistory.CARI4 :=
>   OperatorID.pyUserName` dan `ChronologyInsertion_DT` `.ASMUser := OperatorID.pyUserName`. Yang tercatat keliru hanya
>   baris tangga, yang tetap menyebut akun roster. Penukaran itu **dibuang** (prompt §5 #1), dan `KOMITE_OPERATORID`
>   ditimpa akun pemutus (pola Komite Prop).
> - Kalimat tentang modul saudara tidak dinilai ulang di sini (di luar lingkup).

**Ketiga — perilaku penting bersandar pada langkah yang sudah mati.** `[terverifikasi]` Pemutus
putaran persetujuan bukan pencacah jenjang, melainkan status keputusan terakhir — dan **kedua
langkah yang mengosongkan status itu ber-remark**, begitu pula penutup paksa kasus pada jalur utama.
⚠️ Modul bekerja **karena efek samping**, bukan karena aturan yang tertulis.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"`[terverifikasi]` Pemutus putaran persetujuan
> bukan pencacah jenjang, melainkan status keputusan terakhir — dan **kedua langkah yang mengosongkan status itu
> ber-remark**, begitu pula penutup paksa kasus pada jalur utama."* → pemutusnya **keduanya**. Decision `IsKomiteLoop` di
> `Komite_Flow` = `AcceptStatus == "1"` **AND** `KomiteCount <= KomiteLoop`. `KomitePost_Adjustment` S24 (pre=false)
> **selalu** menaikkan `KomiteCount`, dan S14 menyetel `KomiteCount := KomiteLoop` saat ditolak, sehingga alur berakhir
> karena pencacah melewati `KomiteLoop` atau karena status `"2"`. Yang ter-remark memang benar: pengosong status
> `KomitePostAct` S6 dan `KomiteRouter` S5, serta `ASMForceCaseClose` `KomitePost_Adjustment` S28. Status `"1"` yang
> tidak dikosongkan tidak mengubah hasil, sebab pencacahlah yang mengakhiri. Sistem baru menulis syarat itu eksplisit
> (pola Komite Prop `MasihBerjalan`).

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

> ⛔ **RALAT 10-10-2026.** Baris dan janji lamanya dikutip utuh, tidak dihapus: *"| ⭐ **Giliran pada kasus** | jabatan
> **dan** akun pemegangnya, **dibekukan saat kasus dibuat** | tidak pernah |"*, *"Anggota komite melihat **hanya kasus
> yang menjadi gilirannya**, dan keputusannya **tidak dapat diambil alih orang lain**."* dan *"Pembayaran ke kasir
> **tidak akan terkirim dua kali**."* →
>
> - **Giliran pada kasus** = baris tangga `T_KOMITE_KOMITELIST` (`KOMITE_OPERATORID` = workbasket, `KOMITE_JABATAN` =
>   salinan JABATAN roster). Ia **tidak dibekukan pada akun**, dan tangga awal dapat **diperluas** di tingkat 1
>   (KCF-02). **Pemegang jabatan** = keanggotaan workbasket (KCF-01).
> - Anggota komite melihat kasus `KMT-` yang menunggu workbasket / akunnya di **tabel komite di bawah inbox Claim Fac In**
>   (modul tanpa menu, pola `dcd1522e`). **Setiap anggota** workbasket tingkat berjalan boleh memutus, dan keputusan
>   pertama yang tersimpan mengunci tingkat itu; tanpa larangan rangkap (pola Komite Prop 09-10-2026).
> - Kasir: gerbang `StatusKasir` kosong tetap; dedupe muatan yang masih antre = **OQ-CFI-26** (RALAT AC 61–63).

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

> ⛔ **RALAT 10-10-2026.** User story lamanya dikutip utuh, tidak dihapus:
>
> - US 13 *"Sebagai **anggota komite**, saya ingin **tidak menerima giliran atas klaim yang saya ajukan sendiri**, agar
>   tidak ada benturan kepentingan pada jenjang pertama."* → **gugur** (K9 salah baca, RALAT ID-5).
> - US 17 *"Sebagai **administrator**, saya ingin dapat **menaikkan kasus satu jenjang** ketika pemegang giliran
>   berhalangan, agar kasus tidak macet."* dan US 18 *"Sebagai **auditor**, saya ingin eskalasi itu **terekam** — siapa
>   memindahkan, kapan, dari jenjang mana ke mana — agar perpindahan wewenang tidak tak terlihat."* → **gugur**
>   (eskalasi: nol bukti di korpus, prompt tahap 2 §3).
> - US 20 *"Sebagai **auditor**, saya ingin percobaan menyimpan keputusan oleh bukan pemegang giliran **ditolak dan
>   terekam**, bukan diabaikan diam-diam."* → "ditolak" tetap (403); "terekam" **gugur** (di luar XML).
> - US 23 *"Sebagai **auditor**, saya ingin perubahan jabatan seseorang **masuk jejak audit**, sebab jabatan menentukan
>   siapa boleh menyetujui."* → **gugur** (K13 digantikan workbasket, KCF-01).
> - US 27 *"Sebagai **pembaca laporan**, saya ingin dapat **menyaring** catatan akseptasi yang dibangun dengan aturan
>   lama, agar angkanya tidak tercampur dengan yang dibangun aturan baru."* → **gugur** (KCF-04).
> - US 32 *"Sebagai **penilai klaim**, saya ingin dokumen akseptasi tercetak dan surel terkirim **hanya setelah**
>   keputusan terakhir, agar tidak ada dokumen setengah jadi beredar."* → dokumen akseptasi tetap hanya di tingkat akhir,
>   tetapi **surel terkirim di setiap keputusan** (`KomitePost_Adjustment` S7.2.1.15 `SendEmailKlaim_KMT`: ke tingkat
>   berikut, atau ke pembuat dengan "(Approval)" / "(Reject)").
> - US 38 *"Sebagai **auditor**, saya ingin mengetahui bahwa jalur retro **melewati penyiapan wewenang**, agar hal itu
>   dapat dinilai, bukan tersembunyi."* → jalur retro hanya melewati **perluasan tangga** (`ApprovalKomite_Act` L3),
>   bukan wewenang (RALAT ID-10).

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

> ⛔ **RALAT 10-10-2026.** Baris dan kalimat lamanya dikutip utuh, tidak dihapus: *"| **Pemegang jabatan** | kolom
> `jabatan` di **tabel login**, berisi **kode** |"*, *"| **Giliran pada kasus** | jenjang, jabatan, dan **akun
> pemegangnya**, dibekukan saat kasus dibuat |"* dan *"⭐ `[keputusan work owner]` **K13 · 2026-09-20** — kolom `jabatan` di
> tabel login, dengan **empat syarat**: *(1)* isinya **kode**, bukan teks bebas; *(2)* **urutan jenjang** ada di daftar
> terpisah; *(3)* catatan keputusan **menyimpan salinan kode jabatan** saat itu; *(4)* perubahan kolom `jabatan` **masuk
> jejak audit**."* → **K13 digantikan workbasket** (prompt tahap 2 §3 "Gugur", **KCF-01**):
>
> | Benda | Wujud 10-10-2026 |
> | --- | --- |
> | Daftar induk jabatan | **tidak ada**; label = teks `EMAILKOMITE.JABATAN` |
> | Susunan jenjang komite | roster `EMAILKOMITE` STS_KLAIM FACIN: DEGREE, JABATAN, LIMIT_BOTTOM / LIMIT_TOP (`ApprovalKomite_Act` L4) |
> | Pemegang jabatan | **workbasket** di `OPERATOR_ID` roster (DEGREE 1 `ReasClaimSPVA` + cadangan `ReasClaimSPVB`, 2 `ReasClaimDeptHead`, 3 `ReasClaimTechDivHead`, 4 `ReasClaimOpsDir`, 5 `ReasClaimTechDir`) dan anggotanya |
> | Giliran pada kasus | baris `T_KOMITE_KOMITELIST` tingkat `KOMITE_COUNT`; lahir di claimfacin, diperluas di tingkat 1 (KCF-02); **tidak dibekukan** pada akun |
>
> K12 (jenjang dari jabatan, bukan nama orang) tetap terpenuhi dalam bentuk workbasket.

### ID-2 — Aturan rujuk-atau-salin

⭐ `[keputusan work owner]` **K13** — aturan pembeda yang mengikat seluruh pencatatan:

> ⭐ **Identitas orang DIRUJUK · jabatan DISALIN · tulisan jabatan DIRUJUK.**

| Yang dicatat | Caranya | Alasannya |
| --- | --- | --- |
| **Akun** | rujukan | orangnya tetap orang yang sama; bila namanya berubah, catatan **memang seharusnya** ikut nama baru |
| **Jabatan** | **salinan kode** | jabatan saat itu adalah **fakta sejarah** — ⛔ tidak boleh ikut naik ketika orangnya naik jabatan |
| **Tulisan jabatan** | rujukan ke daftar induk | bila label jabatan yang sama ditulis ulang, catatan lama boleh ikut; ⛔ yang dilarang adalah catatan lama **berpindah ke jabatan berbeda** |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"| **Jabatan** | **salinan kode** | jabatan saat
> itu adalah **fakta sejarah** — ⛔ tidak boleh ikut naik ketika orangnya naik jabatan |"* dan *"| **Tulisan jabatan** |
> rujukan ke daftar induk | bila label jabatan yang sama ditulis ulang, catatan lama boleh ikut; ⛔ yang dilarang adalah
> catatan lama **berpindah ke jabatan berbeda** |"* → dengan K13 gugur, **jabatan disalin sebagai teks** JABATAN roster
> ke `KOMITE_JABATAN` saat baris tangga dibentuk (`ApprovalKomite_Act` L7.1 `IDKomite := .JABATAN`; `SetListKomite_act`
> sisi klaim), lalu ikut tertulis di teks kronologi. Tidak ada daftar induk, jadi baris "Tulisan jabatan" gugur. Akun
> tetap rujukan: `KOMITE_OPERATORID` ditimpa akun pemutus (pola Komite Prop). Aturan "identitas dirujuk · jabatan
> disalin" (ADR-0039 di `OUTPUT_HASIL_RNM/docs/bersama/adr/`) tetap terpenuhi untuk dua baris pertama.

### ID-3 — Pemutus tangga adalah keadaan jenjang

⭐ `[keputusan work owner]` **K1 · 2026-09-20** — komite berakhir bila **seluruh jenjang
menyetujui**; satu menolak ⇒ **kasus tutup seketika**.

`[terverifikasi]` Di sistem lama, penugasan berhenti karena router **tidak menghasilkan penerima**
ketika setiap jenjang sudah memutuskan — sebuah **perilaku diam** yang tidak dapat diuji.
⚠️ `[penyimpangan sadar]` Di sistem baru, syarat berakhir **dinyatakan eksplisit** dan dapat diuji.

⛔ **Pencacah jenjang boleh ada sebagai tampilan, tetapi bukan kebenaran.** `[terverifikasi]` Daftar
jenjang dapat **berubah di awal** — pengaju dikeluarkan — dan pencacah tidak ikut tahu.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"`[terverifikasi]` Di sistem lama, penugasan
> berhenti karena router **tidak menghasilkan penerima** ketika setiap jenjang sudah memutuskan — sebuah **perilaku diam**
> yang tidak dapat diuji."* dan *"`[terverifikasi]` Daftar jenjang dapat **berubah di awal** — pengaju dikeluarkan — dan
> pencacah tidak ikut tahu."* →
>
> - K1 tetap (semua setuju ⇒ selesai; satu tolak ⇒ tutup seketika). **Mekanismenya** di XML: Decision `IsKomiteLoop`
>   (`AcceptStatus == "1"` AND `KomiteCount <= KomiteLoop`) sesudah `KomitePostAct`; `KomitePost_Adjustment` S24
>   (pre=false) selalu menaikkan `KomiteCount`, S14 menyetel `KomiteCount := KomiteLoop` saat ditolak, dan S7.2.1.9.1.1
>   menandai sisa baris menunggu `"2"`. Router tidak berperan sebagai pemutus. `KomiteRouter` S6.1 tanpa transisi keluar
>   memberi `AssignTo` = baris menunggu **terakhir**; sistem baru memakai baris tingkat berjalan (`KomiteCount`), selaras
>   prompt §5 #1.
>   ⛔ **RALAT 10-10-2026 (kedua)** atas kalimat di atas, dikutip utuh: *"`KomiteRouter` S6.1 tanpa transisi keluar memberi `AssignTo` = baris menunggu **terakhir**"* → `KomiteRouter` S6.1 **punya**
>   transisi pasca-langkah (`pyStepsTransParams` `true` → 6 keluar; alat dump sebelumnya hanya mencetak prasyarat):
>   Pega menugaskan baris menunggu **pertama**. Baris tingkat berjalan (`KomiteCount`) sistem baru setara pada alur
>   berurutan; prompt §5 #1 tetap berlaku untuk `SetProteksiSubmiteKomite` L2.1 (lolos untuk baris menunggu mana saja).
> - **Pencacah ikut tahu.** Daftar jenjang memang berubah di awal, tetapi karena **perluasan** (KCF-02,
>   `ApprovalKomite_Act` L4–L7), bukan karena pengaju dikeluarkan (K9 gugur). Sesudahnya L8 menghitung ulang
>   `KomiteLoop`. Di sistem baru `KOMITE_COUNT` dan `KOMITE_LOOP` adalah kebenaran yang tersimpan dan maju per tingkat
>   (pola Komite Prop `MasihBerjalan`), bersama `POSITION`. ADR-0038
>   (`OUTPUT_HASIL_RNM/docs/bersama/adr/0038-pemutus-tangga-persetujuan-adalah-keadaan-jenjang.md`) bersumber ID-3 ini
>   dan tidak disunting di sini.

### ID-4 — Penegakan wewenang di lapisan layanan

⭐ `[keputusan work owner]` **K10 · 2026-09-20** — ADR-0014 **berdiri**. Penegakan pemilik giliran
**tetap di lapisan layanan**, bacaan **sempit**.

Lapisan layanan memeriksa **identitas akun pemanggil** terhadap **akun beku pada jenjang berjalan**
sebelum menyimpan keputusan. ⛔ Gagal ⇒ **tolak**, bukan abaikan diam-diam. ⛔ Layar boleh
menyembunyikan tombol; itu kenyamanan, **bukan** penegakan.

⭐ `[keputusan work owner]` **K3** — pembandingan memakai **identitas akun**, ⛔ bukan alias.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Lapisan layanan memeriksa **identitas akun
> pemanggil** terhadap **akun beku pada jenjang berjalan** sebelum menyimpan keputusan."* → akun pemanggil diperiksa
> terhadap **keanggotaan workbasket** `KomiteID` tingkat berjalan: `Pemegang(akun, peran)`, pola Komite Prop 09-10-2026.
> Tingkat 1 = anggota `ReasClaimSPVA` atau `ReasClaimSPVB` (KCF-01); TT3 / TT4 = anggota `ReasClaimDeptHead` (KCF-03);
> selainnya **403**. ADR-0014 kini di
> `OUTPUT_HASIL_RNM/docs/bersama/adr/0014-pemutus-komite-ditegakkan-per-komiteid-tingkat-berjalan.md`, dan aturan perannya
> di ADR-0030 (`OUTPUT_HASIL_RNM/docs/bersama/adr/0030-aturan-peran-ditetapkan-sekali-lintas-modul.md`). K10 dan K3
> tetap.

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

> ⛔ **RALAT 10-10-2026 — ID-5 GUGUR.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Saat kasus komite dibentuk, bila
> **pemegang jenjang pertama adalah pengaju**, ia **dikeluarkan dari daftar**, dan **jumlah jenjang dihitung ulang**."* →
> **K9 salah baca** (prompt tahap 2 §3 "Gugur"). `ApprovalKomite_Act` L6.2 (`Property-Remove` when
> `pyWorkPage.KomiteList(1).KomiteID == .OPERATOR_ID`) hanya membuang **calon perluasan** (roster DEGREE > 1 hasil L4 /
> L5) yang sama dengan **anggota tingkat 1**. Ia tidak membandingkan pengaju dan tidak mengeluarkan siapa pun dari
> tangga. Hitung ulang `KomiteLoop` (L8) terjadi sesudah calon ditambahkan (L7). Korpus Komite Claim FacIn **tidak
> memuat larangan menyetujui klaim sendiri** di jenjang mana pun, jadi tidak ada yang dibangun, dan tanpa larangan
> rangkap (pola Komite Prop 09-10-2026). Kalimat "lubang yang dibawa masuk dengan sadar" kehilangan dasarnya. ADR-0040
> (`OUTPUT_HASIL_RNM/docs/bersama/adr/0040-larangan-menyetujui-pekerjaan-sendiri.md`) bersumber ID-5 ini dan **tidak
> disunting di sini** — pemutakhirannya
> menjadi pertanyaan untuk work owner.

### ID-6 — Data kutipan disalin utuh

⭐ `[keputusan work owner]` **K2 · 2026-09-20** — ⚠️ `[penyimpangan sadar]`, dicatat sebagai
**cacat yang diperbaiki**.

`[terverifikasi]` Di sistem lama, objek kerja komite hanya menerima **dua medan** data kutipan pada
jalur penerimaan, sedangkan **sembilan** penggolong lini usaha menguji medan lain — **dua di
antaranya dipakai hidup lima kali**. ⛔ Cabang MBU dan Travel karenanya **tidak pernah terbit**.

⭐ **Di sistem baru, seluruh data kutipan yang dibutuhkan komite disalin dari kasus induk** — ⛔
bukan daftar medan bernama satu per satu, sebab **daftar bernama itulah yang melahirkan cacat ini**.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⭐ **Di sistem baru, seluruh data kutipan yang
> dibutuhkan komite disalin dari kasus induk** — ⛔ bukan daftar medan bernama satu per satu, sebab **daftar bernama
> itulah yang melahirkan cacat ini**."* → komite **tidak menyalin** apa pun. Halaman kasus klaim induk, termasuk
> `OfferFacIn.QuotationData` / `PolicyData`, **dibaca** lewat kontrak `kontrak.KlaimFacInKomite` (`BacaKlaimFacIn`,
> `inti/backend/kontrak/klaimfacin.go`), disediakan claimfacin (`services.KlaimUntukKomite`), tanpa impor modul
> claimfacin. Di XML pun komite membaca halaman induk langsung: `pyWorkCover.OfferFacIn.QuotationData` (`SetValueKomite`
> S12) dan `TempOpenPage.OfferFacIn.QuotationData.BusinessOldId` (`KomitePost_Adjustment` S7.2.1.2.7). Tujuan K2 (semua
> medan tersedia bagi penggolong lini) tercapai tanpa salinan. Fakta "dua medan disalin" (`SetValueKomite` S12:
> `BusinessType`, `GroupPanel`) tidak dinilai ulang di sini.

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

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⭐ `[keputusan work owner]` **K5 · 2026-09-20**
> — penanda pembayaran **tersimpan dalam transaksi yang sama** dengan pengiriman; panggilan kedua **ditolak**."* → lihat
> RALAT AC 61–63: kasir lewat **outbox** (hanya produksi), sehingga penanda `StatusKasir` baru terisi sesudah pelaksana
> outbox menjalankan muatan; dedupe saat muatan masih antre = **OQ-CFI-26**. Tanda `=` tunggal dibaca **pembandingan**
> (prompt §5 #2).

### ID-9 — Pemberitahuan galat hanya pada kegagalan

⭐ `[keputusan work owner]` **K6 · 2026-09-20** — ⚠️ `[penyimpangan sadar]`.

`[terverifikasi]` Sistem lama mengirim pemberitahuan galat **setiap kali**, sebab gerbangnya
bertanda bendera mati; medan yang diujinya pun **salah eja**. ⭐ Sistem baru memberitahu **hanya
pada kegagalan**, dan **menamai medannya dengan benar**.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"`[terverifikasi]` Sistem lama mengirim
> pemberitahuan galat **setiap kali**, sebab gerbangnya bertanda bendera mati; medan yang diujinya pun **salah eja**."* →
> **bukan setiap kali.** Di `HitServiceToKasirKMT_Act`, L15 `Exit-Activity` ("exit disini jika tidak ada error") menutup
> jalur biasa. L16 (label `END`) dan L17 `SendErrorDirectKasir` (yang ber-pre=false) hanya dicapai lewat lompatan
> kegagalan (`pyOnException` = `END` pada kedua panggilan luar). `SendErrorDirectKasir` L2 (pre=true) lalu memeriksa
> ulang `.StatusServiceKasir.ReponseCode != "1"`. Jadi K6 **sama dengan XML**, bukan penyimpangan sadar. Surelnya lewat
> outbox hanya produksi, BCC orang dibuang (prompt §6 butir 8). Ejaan `ReponseCode` adalah nama properti XML.

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

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Jenis treaty tertentu menandai penyesuaian
> sebagai **Fac Retro**, dan penyesuaian itu **melompati penyiapan wewenang penyetujuan**."* dan *"2. ⛔⛔ **Ada jalur
> yang melompati pemeriksaan wewenang.** ⭐ Karena itu spec ini **menyebutnya sebagai jalur tersendiri** — bukan
> menguburnya di dalam detail — supaya siapa pun yang membaca aturan **melihat jalur itu ada**."* → **tidak ada jalur
> yang melompati pemeriksaan wewenang.** Yang dilompati Fac Retro hanya:
>
> - **perluasan tangga**: `ApprovalKomite_Act` L2.2 (`IsFacRetro == 1` pada adjustment KMT) → L3 `Exit-Activity`
>   melewati L4–L8. Ini **KCF-02** ("Kasus Fac Retro melewati perluasan (L3); tangganya sudah dibentuk sisi klaim").
>   Tangga retro lahir di claimfacin (`RosterKomiteCalon`: batas = `ValueAdjustment`, tanpa saringan bila melampaui
>   `LIMIT_TOP` Technic Div. Head);
> - **konversi** dan **log "AKSEPATSI"** di tingkat akhir (`KomitePost_Adjustment` S17 / S19, `Local.FacRetro == 1
>   [T=3]`).
>
> Gerbang Submit (`SetValueKomite` S15 `SetProteksiSubmiteKomite`) tetap berjalan untuk retro, dan di sistem baru
> `Pemegang` berlaku untuk semua kasus. Akibat 1 (konstanta `"10015"` di kode, `SaveAcceptation_KMT` L2.1) tetap. Asal
> K7 lihat RALAT Q3 di `keputusan-sebelum-to-spec.md`. OQ jalur kasir retro (S25 pre=false) di
> `issues/10-jalur-fac-retro.md`.

### ID-11 — Jabatan dibaca dari data pengguna

⭐ `[keputusan work owner]` **K11 · 2026-09-20** — ⚠️ `[penyimpangan sadar]`.

⛔ Tabel jabatan yang ditanam di dalam ekspresi **tidak dibawa**, **berikut nilai bawaannya**.
⭐ Bila jabatan seseorang **tidak diketahui**, sistem **menolak** — ⛔ tidak memberi jabatan bawaan.

`[terverifikasi]` Di sistem lama, nilai bawaannya adalah sebuah **jabatan direksi**, diberikan
kepada **siapa pun yang tak dikenal**, dan jejak audit mencatatnya sebagai fakta.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⭐ Bila jabatan seseorang **tidak diketahui**,
> sistem **menolak** — ⛔ tidak memberi jabatan bawaan."* → jabatan dibaca dari **baris tangga** (`KOMITE_JABATAN` =
> JABATAN roster; `KomitePost_Adjustment` S3 / S4 `KomiteList(KomiteCount).IDKomite`), bukan dari data pengguna, sehingga
> tidak pernah "tidak diketahui". Tanpa jabatan bawaan tetap berlaku. Rantai jabatan + bawaan direksi di modul ini adalah
> sisa editor (RALAT lintas-ronde #2), bukan rule berjalan.

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

> ⛔ **RALAT 10-10-2026 — ID-12 GUGUR.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⭐ `[keputusan work owner]` **K8 ·
> 2026-09-20** — catatan akseptasi historis MBU dan Travel **tidak dibangun ulang**, tetapi **ditandai**."* → **KCF-04**
> (10-10-2026): kasus komite lama Pega (`KMT-`) tidak dimigrasi; K8 `DIBANGUN_ATURAN_LAMA` **gugur**. Tanpa kolom penanda.
> Lihat RALAT AC 91–95.

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

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"| ⚠️⚠️ **U-5** | **Pengeluaran pengaju — DUA
> ARAH** *(ID-5)*: pengaju di jenjang **pertama** dikeluarkan dan jumlah jenjang berkurang; pengaju di jenjang **kedua** ⛔
> **TETAP DI DAFTAR** | ⭐ **Arah kedua mengunci penyimpangan sadar K9** supaya tidak "terperbaiki" diam-diam oleh
> pengembang berikutnya |"* dan *"| ⭐⭐ **U-7** | **Salinan jabatan** *(ID-2)* — catat sebuah keputusan, **naikkan jabatan
> orangnya**, lalu periksa **catatan lama TIDAK berubah** | ⛔ inilah yang membedakan **salinan** dari **rujukan** |"* →
>
> - **U-5 gugur** (K9 salah baca). Penggantinya **U-5′ Perluasan tangga** (KCF-02): tangga satu baris + total adjustment ⇒
>   baris DEGREE > 1 ber-`LIMIT_BOTTOM` < total, urut DEGREE; calon yang sama dengan anggota tingkat 1 tidak masuk dua
>   kali (L6.2); pita SPV B (pemutus anggota `ReasClaimSPVB`, 30.000.000 < total ≤ 57.750.000) ⇒ satu jenjang atas;
>   Fac Retro ⇒ tanpa perluasan; GET tidak menulis.
> - **U-7** dibaca ulang: catat keputusan, **ubah JABATAN baris roster** ⇒ `KOMITE_JABATAN` baris tangga lama dan teks
>   kronologinya tidak berubah. Jabatan melekat pada workbasket, bukan pada akun (KCF-01).
> - **U-6** memakai halaman klaim lewat kontrak (RALAT ID-6). **U-3** tunduk pada OQ-CFI-26 (RALAT AC 61–63).

---

## Acceptance Criteria

### A · Pembentukan kasus komite

**AC 1** `[terverifikasi]` Ketika penilai mengajukan penyesuaian ke komite, sistem membentuk **satu
kasus komite** yang tertaut ke klaim induk.

**AC 2** `[terverifikasi]` Saat pembentukan, sistem menetapkan **jumlah jenjang** dari **susunan
jenjang komite** yang berlaku.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Saat pembentukan, sistem menetapkan **jumlah
> jenjang** dari **susunan jenjang komite** yang berlaku."* → jumlah jenjang ditetapkan **dua kali**:
>
> - **saat lahir** oleh claimfacin tahap 1 (`CreateKMTNo_Act` 8 → `SetListKomite_act`, commit `8c3b2e71`): roster FACIN
>   aktif ber-`LIMIT_BOTTOM ≤ 1` untuk non-retro (praktis hanya DEGREE 1), atau batas `ValueAdjustment` untuk retro;
>   `KOMITE_LOOP` = cacah baris;
> - **saat tingkat 1** oleh modul ini (**KCF-02**, `ApprovalKomite_Act` L4–L8): bila `KomiteCount = 1` dan tangga satu
>   baris, roster DEGREE > 1 ber-`LIMIT_BOTTOM` < total adjustment ditambahkan (atau satu jenjang saja pada pita SPV B,
>   L5), lalu `KOMITE_LOOP` dihitung ulang. Ditampilkan saat dibuka tanpa menulis, disimpan saat tingkat 1 Submit. Fac
>   Retro dilewati (L3).

**AC 3** `[terverifikasi]` Giliran dimulai pada **jenjang pertama**.

**AC 4** `[terverifikasi]` Jalur **tutup klaim** dan jalur **tolak klaim** selalu membentuk komite
**berjenjang satu**.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Jalur **tutup klaim** dan jalur **tolak klaim**
> selalu membentuk komite **berjenjang satu**."* → tetap satu jenjang, dengan rincian **KCF-03**: anggota tunggal
> `SendRejectClaimToKomite2` (TT3) / `SendCloseClaimToKomite` (TT4) di XML adalah akun orang tertulis mati, sedangkan di
> sistem baru satu tingkat = **`ReasClaimDeptHead`** (prompt §5 #7). Kasusnya tanpa adjustment, jadi migrasi
> komiteclaimfacin (640–679) `MODIFY T_GENERAL_KOMITE.ADJUSTMENT_ID` menjadi boleh kosong dan menambah kolom jenis
> penyerahan `TRANSFER_TYPE` (2 / 3 / 4). Tombol Yes di claimfacin (`SureRejectClaim` / `PreventRejectClaim`) dinyalakan.
> Hanya Fac In. Yang dimaksud "tutup klaim" di sini adalah **Close Without Payment** (TT4).

**AC 5** `[keputusan work owner]` **K13** Saat pembentukan, tiap jenjang **dibekukan** bersama
**jabatan** dan **akun pemegangnya**. Perubahan pemegang jabatan sesudah itu **tidak memindahkan
kasus**.

> ⛔ **RALAT 10-10-2026 — AC 5 GUGUR.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Saat pembentukan, tiap jenjang
> **dibekukan** bersama **jabatan** dan **akun pemegangnya**. Perubahan pemegang jabatan sesudah itu **tidak memindahkan
> kasus**."* → K13 digantikan workbasket (**KCF-01**). Baris tangga menyimpan `KOMITE_JABATAN` (salinan teks) dan
> `KOMITE_OPERATORID` = **workbasket**, bukan akun. Pemutusnya = siapa pun anggota workbasket tingkat berjalan saat
> Submit (pola Komite Prop 09-10-2026), jadi perubahan anggota workbasket **memang** memindahkan siapa yang dapat memutus.
> Tangga juga tidak beku: ia diperluas di tingkat 1 (KCF-02).

**AC 6** `[keputusan work owner]` **K13** Bila sebuah jabatan pada susunan jenjang **tidak memiliki
pemegang aktif**, pembentukan kasus **ditolak dengan galat yang menyebut jabatannya**.

> ⛔ **RALAT 10-10-2026 — AC 6 GUGUR.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Bila sebuah jabatan pada susunan
> jenjang **tidak memiliki pemegang aktif**, pembentukan kasus **ditolak dengan galat yang menyebut jabatannya**."* → nol
> bukti XML (prompt §1: "Nol … aksi tanpa bukti XML"). `ApprovalKomite_Act` L4 / L5.2 hanya menyaring roster `STS_AKTIF
> = "1"` dan `STS_KLAIM = "FACIN"`; tidak ada pemeriksaan pemegang. Kelahiran kasus milik claimfacin. Workbasket tanpa
> anggota aktif berarti kasus menunggu di tingkat itu sampai anggotanya diatur (Kelola User).

**AC 7** `[terbuka]` **butir 39** Bila sebuah jabatan memiliki **lebih dari satu** pemegang aktif,
perilakunya **belum diputuskan**. ⭐ **Usul asisten:** tolak dengan galat yang jelas — ⛔ jangan
memilih diam-diam. ⛔ **Belum menjadi keputusan.**

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Bila sebuah jabatan memiliki **lebih dari satu**
> pemegang aktif, perilakunya **belum diputuskan**."* → **tertutup** pola Komite Prop 09-10-2026 (prompt tahap 2 §3):
> penyetuju = anggota workbasket tingkat berjalan (`Pemegang`), **siapa pun** anggotanya boleh memutus, keputusan pertama
> yang tersimpan mengunci tingkat itu, surel ke **semua** anggota workbasket. Tingkat 1 punya dua workbasket (`ReasClaimSPVA`
> dan cadangan `ReasClaimSPVB`, KCF-01). Butir 39 tidak menahan.

**AC 8** `[keputusan work owner]` **K2** Seluruh data kutipan yang dibutuhkan penggolongan lini
usaha **tersalin dari kasus induk** ke kasus komite pada saat pembentukan.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Seluruh data kutipan yang dibutuhkan
> penggolongan lini usaha **tersalin dari kasus induk** ke kasus komite pada saat pembentukan."* → **tidak disalin**:
> dibaca dari kasus induk lewat kontrak `kontrak.KlaimFacInKomite` (`BacaKlaimFacIn`) setiap kali dibutuhkan (layar,
> pra-proses, Submit). Kasus `KMT-` tahap 1 hanya menyimpan tautan (`COVER_KEY`, `T_GENERAL_KOMITE.ADJUSTMENT_ID`). Lihat
> RALAT ID-6.

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

> ⛔ **RALAT 10-10-2026 — BAB B (AC 9–13) GUGUR.** AC lamanya dikutip utuh, tidak dihapus:
>
> - AC 9 *"Bila pemegang **jenjang pertama** adalah **pengaju** penyesuaian, ia **dikeluarkan** dari susunan jenjang
>   kasus itu."*
> - AC 10 *"Sesudah pengeluaran, **jumlah jenjang berkurang satu**, dan jenjang berikutnya **naik menjadi jenjang
>   pertama**."*
> - AC 11 *"Bila pengaju memegang jenjang **kedua atau lebih rendah**, ia **TETAP berada di susunan jenjang** dan **tetap
>   dapat menyetujui klaim yang ia ajukan sendiri**. ⛔ **Ini lubang yang dibawa masuk dengan sadar, bukan kelalaian.**"*
> - AC 12 *"Perilaku AC 11 **wajib memiliki uji yang membuktikannya** *(U-5 arah kedua)*, agar tidak diubah tanpa
>   keputusan work owner."*
> - AC 13 *"Bila pengeluaran pengaju membuat susunan jenjang menjadi **kosong**, kasus komite **tidak dibentuk** dan
>   pengajuan **dikembalikan** kepada penilai dengan alasannya."*
>
> → **K9 salah baca** (prompt tahap 2 §3 "Gugur"; RALAT ID-5). `ApprovalKomite_Act` L6.2 membuang calon perluasan yang
> sama dengan anggota tingkat 1, bukan pengaju. Tidak ada pengeluaran, tidak ada pengurangan jenjang, tidak ada kasus
> yang dikembalikan, dan uji U-5 diganti U-5′ (perluasan). Larangan menyetujui klaim sendiri tidak ada di XML pada
> jenjang mana pun, sehingga tidak dibangun.

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

> ⛔ **RALAT 10-10-2026 — AC 20–22 GUGUR.** AC lamanya dikutip utuh, tidak dihapus: AC 20 *"Administrator dapat
> **menaikkan kasus satu jenjang** ketika pemegang giliran berhalangan. ⛔ Menurunkan jenjang **dilarang**."*, AC 21
> *"Eskalasi **terekam**: siapa memindahkan, kapan, dari jenjang mana ke jenjang mana."*, dan AC 22 *"Jenjang yang
> **dilewati** oleh eskalasi tercatat **tanpa keputusan** — ⛔ bukan tercatat sebagai menyetujui."* → **eskalasi: nol
> bukti di korpus** (prompt tahap 2 §3 "Gugur"), walau ketiganya bertanda `[terverifikasi]`. Sisiran 10-10-2026 atas 114
> berkas XML Komite Claim FacIn: nol kemunculan `escalat` / `eskalasi` / `pyServiceLevel`; `KomiteRouter` S1–S5
> ter-remark dan S6 hanya menetapkan `AssignTo`.
> Tidak dibangun. Pemegang yang berhalangan digantikan lewat keanggotaan workbasket (KCF-01).

### D · Wewenang

**AC 23** ⭐ `[keputusan work owner]` **K10** Hanya **akun beku pada jenjang berjalan** yang dapat
menyimpan keputusan. Penolakannya terjadi **di lapisan layanan**.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Hanya **akun beku pada jenjang berjalan** yang
> dapat menyimpan keputusan. Penolakannya terjadi **di lapisan layanan**."* → "akun beku" diganti **anggota workbasket
> tingkat berjalan** (`Pemegang(akun, peran)`, pola Komite Prop 09-10-2026, prompt §6 butir 5): tingkat 1 = anggota
> `ReasClaimSPVA` atau `ReasClaimSPVB` (KCF-01); tingkat lain = anggota workbasket `KomiteID` baris `KomiteCount`;
> TT3 / TT4 = anggota `ReasClaimDeptHead` (KCF-03). Selainnya **403**. Submit juga ditolak bila kasus tertutup atau
> tingkatnya sudah memutus. Penegakan di lapisan layanan tetap (K10). XML `SetProteksiSubmiteKomite` L2.1 membuka Submit
> bagi baris menunggu **mana saja** dan bagi "IT Developer"; keduanya diperbaiki (prompt §5 #1).

**AC 24** `[keputusan work owner]` **K10** Percobaan menyimpan oleh akun lain **ditolak dengan
galat**, ⛔ bukan diabaikan diam-diam.

**AC 25** `[keputusan work owner]` **K10** Percobaan yang ditolak **terekam** di jejak audit.

> ⛔ **RALAT 10-10-2026 — AC 25 GUGUR.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Percobaan yang ditolak **terekam**
> di jejak audit."* → **di luar XML** (prompt tahap 2 §3 "Gugur"). Tidak ada rule pencatat percobaan ditolak; jawaban 403
> cukup (AC 24).

**AC 26** `[keputusan work owner]` **K3 · K12** Pembandingan memakai **identitas akun**. ⛔ Tidak ada
penukaran nama sebelum pembandingan.

**AC 27** ⚠️ `[penyimpangan sadar]` **K12** `[terverifikasi]` Sistem lama menukar **dua akun** menjadi
akun ketiga **sebelum** pemeriksaan pemilik giliran. ⛔ **Penukaran itu tidak dibawa.**

**AC 28** `[keputusan work owner]` **K10** Layar boleh menyembunyikan tindakan yang tak berwenang,
⛔ **tetapi itu bukan penegakan** — penegakan tetap di lapisan layanan.

**AC 29** `[terbuka]` **butir 6** Daftar jabatan dan jabatan mana yang membentuk susunan jenjang
**belum diketahui**. ⭐ **Dikurung sebagai titik sambung bernama:** aturan wewenang di atas dapat
ditulis dan diuji tanpa isinya; ⛔ **pembangunannya menunggu**.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Daftar jabatan dan jabatan mana yang membentuk
> susunan jenjang **belum diketahui**."* dan *"⛔ **pembangunannya menunggu**."* → **butir 6 tertutup.** Penutupnya ADR-0030
> *"Aturan peran ditetapkan sekali dan berlaku lintas modul"*
> (`OUTPUT_HASIL_RNM/docs/bersama/adr/0030-aturan-peran-ditetapkan-sekali-lintas-modul.md`) + keanggotaan workbasket +
> **KCF-01**. Isi roster FACIN (DEV 10-10-2026): DEGREE 1 Claim Supervisor (`LIMIT_TOP` 57.750.000) → `ReasClaimSPVA`
> (cadangan `ReasClaimSPVB`), 2 Claim Dept. Head (`LIMIT_BOTTOM` 57.750.001) → `ReasClaimDeptHead`, 3 Technic Div. Head
> (189.750.001) → `ReasClaimTechDivHead`, 4 Operational Director (495.000.001) → `ReasClaimOpsDir`, 5 Technical
> Director (660.000.001) → `ReasClaimTechDir`. Keenam workbasket ada di `M_WORKBASKET` DEV. Pembangunan tidak menunggu.
> AC 97 merujuk AC 29 dan ikut tertutup.

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

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: AC 37 *"Jabatan tersimpan sebagai **salinan
> kode**. ⛔ Ketika orangnya kelak naik jabatan, **catatan lama tidak berubah**."*, AC 39 *"**Tulisan jabatan** dirujuk
> dari **daftar induk**. ⛔ Catatan lama **tidak boleh berpindah ke jabatan yang berbeda**."* dan AC 40 *"Bila jabatan
> pemutus **tidak diketahui**, penyimpanan keputusan **ditolak**. ⛔ **Tidak ada jabatan bawaan.**"* → dengan K13
> digantikan workbasket (KCF-01): jabatan tersimpan sebagai **salinan teks** JABATAN roster di `KOMITE_JABATAN`, bukan
> kode. Catatan lama tidak berubah ketika JABATAN roster diubah, dan "naik jabatan" orang tidak bermakna karena jabatan
> melekat pada workbasket. AC 39 **gugur** (tidak ada daftar induk). AC 40: jabatan dibaca dari baris tangga, jadi tidak
> pernah "tidak diketahui"; tanpa bawaan tetap berlaku. AC 36 / 38 / 42 tetap (`KOMITE_OPERATORID` ditimpa akun
> pemutus, pola Komite Prop; Submit ditolak bila tingkatnya sudah memutus). AC 41: rantai bawaan direksi adalah sisa
> editor (RALAT lintas-ronde #2), tetap tidak dibawa. Rinciannya di tiket 05.

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

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: AC 46 *"**Nomor akseptasi dibangkitkan** setelah
> ringkasan tersimpan."* dan AC 48 *"**Surel pemberitahuan terkirim** setelah dokumen siap."* → urutan XML
> `KomitePost_Adjustment` berbeda: nomor (S7.2.1.2) → tulis adjustment (S7.2.1.5) → surel (S7.2.1.15) →
> `SaveAcceptation_KMT` (S7.2.1.16) → `SaveAccept_ACT` (S7.2.1.17) → `OS_AKSEPTASI_KLAIM` (S8) → PDF (S9). Nomor terbit
> **sebelum** ringkasan disimpan. Surel terkirim di **setiap** keputusan (`SendEmailKlaim_KMT`: ke tingkat berikut, atau
> ke pembuat "(Approval)" / "(Reject)"), **sebelum** PDF, dan tidak menunggu dokumen. Di sistem baru keduanya efek outbox
> sesudah commit transaksi Submit; PDF sesudah commit (isi berkas OQ, stream tidak diekspor).

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

> ⛔ **RALAT 10-10-2026.** AC lamanya dikutip utuh, tidak dihapus: AC 61 *"Instruksi pembayaran terkirim ke kasir **tepat
> satu kali** per akseptasi."*, AC 62 *"Penanda "sudah terkirim" tersimpan **dalam transaksi yang sama** dengan
> pengirimannya."*, dan AC 63 *"Panggilan kedua atas akseptasi yang sama **ditolak**."* →
>
> - **Kasir lewat outbox, hanya produksi** (prompt §2 butir 6 dan §6 butir 8). Gerbang XML dipertahankan:
>   `KomitePost_Adjustment` S25.2.1.1 `.DirectToKasir == "true" && .StatusKasir = ""` (tanda `=` tunggal dibaca
>   **pembandingan**, prompt §5 #2) + `HitServiceToKasirKMT_Act` L3 `StatusKonversi = 1`; satu panggilan per KMT
>   (S25.2.1 "HANYA LOOPING 1 KALI"); jalur `CLM` L14.2, jumlah per `AcceptedNo`.
> - **AC 62 tidak dapat dipenuhi apa adanya**: muatan diantre dalam transaksi Submit, sedangkan `StatusKasir` baru terisi
>   dari tanggapan kasir (`HitServiceToKasirKMT_Act` L14.5, `DIRECTTOKASIR_LOG`) sesudah pelaksana outbox menjalankan
>   muatan.
> - **AC 61 / 63**: panggilan kedua ditolak bila `StatusKasir` sudah terisi. Dedupe ketika muatan pertama **masih antre** =
>   **OQ-CFI-26** (`modul/claimfacin/docs/OQ.md`), belum diputuskan. K5 berlaku sebatas itu.

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

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Sistem lama mengirimnya **setiap kali**, sebab
> gerbangnya bertanda bendera mati."* → **bukan setiap kali**; lihat RALAT ID-9 (`HitServiceToKasirKMT_Act` L15
> `Exit-Activity` sebelum label `END`, `SendErrorDirectKasir` L2 memeriksa ulang). AC 66 sama dengan XML; AC 67 bukan
> penyimpangan.

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

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Penyesuaian bertanda **Fac Retro** mengikuti
> **jalur tersendiri** yang **melompati penyiapan wewenang penyetujuan**."* → yang dilompati hanya **perluasan tangga**
> (`ApprovalKomite_Act` L2.2 / L3, **KCF-02**) serta **konversi** dan **log "AKSEPATSI"** di tingkat akhir
> (`KomitePost_Adjustment` S17 / S19). Tangga retro dibentuk sisi klaim saat kelahiran. Wewenang (`Pemegang`) tetap
> ditegakkan. Lihat RALAT ID-10.

**AC 72** ⚠️ `[penyimpangan sadar]` **K7** ⛔ **BERBEDA dari rekomendasi asisten**, yang mengusulkan
menjadikan penanda jenis treaty sebagai data. ⛔ **Ditiru apa adanya, dengan sadar.**

**AC 73** ⚠️ `[keputusan work owner]` **K7** Nilai penanda jenis treaty tetap **tetapan di dalam
kode** ⇒ ⛔ **setiap perubahan daftar treaty menuntut rilis perangkat lunak.**

**AC 74** ⛔⛔ `[keputusan work owner]` **K7** Spec ini **menyebut jalur Fac Retro secara terbuka
sebagai jalur yang melompati pemeriksaan wewenang**, ⭐ supaya siapa pun yang membaca aturan
**melihat jalur itu ada** — ⛔ bukan menemukannya setelah sesuatu terjadi.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Spec ini **menyebut jalur Fac Retro secara
> terbuka sebagai jalur yang melompati pemeriksaan wewenang**, ⭐ supaya siapa pun yang membaca aturan **melihat jalur
> itu ada** — ⛔ bukan menemukannya setelah sesuatu terjadi."* → **keliru menurut XML.** `SetValueKomite` S14
> (`ApprovalKomite_Act`, yang keluar seketika untuk retro) diikuti S15 `SetProteksiSubmiteKomite` yang tetap berjalan,
> dan di sistem baru `Pemegang` berlaku untuk setiap kasus. Yang disebut terbuka kini: **Fac Retro tidak mendapat
> perluasan tangga, tidak memicu konversi, dan tidak menulis log "AKSEPATSI"**. Jalur kasir untuk retro tetap berjalan
> (S25 ber-pre=false); maksudnya OQ, lihat `issues/10-jalur-fac-retro.md`.

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

> ⛔ **RALAT 10-10-2026 — BAB L DIBACA ULANG (K13 digantikan workbasket, KCF-01).** AC lamanya dikutip utuh, tidak
> dihapus:
>
> - AC 77 *"Jabatan seseorang disimpan sebagai **kolom di tabel login**, berisi **kode**."* → **gugur**. Tidak ada kolom
>   jabatan di tabel login; jabatan melekat pada baris roster / workbasket.
> - AC 78 *"Kode itu menunjuk **daftar induk jabatan**, yang memuat tulisan tampilnya. ⛔ **Bukan teks bebas.**"* →
>   **gugur**. Label = teks `EMAILKOMITE.JABATAN`, disalin ke `KOMITE_JABATAN`; tidak ada daftar induk.
> - AC 79 *"Di sistem lama, orang yang sama tertulis dengan **jabatan berbahasa berbeda** antar modul, dan **dibaca
>   lewat medan yang berbeda**. ⭐ Kode tunggal menutup keduanya."* → faktanya tidak dinilai ulang (menyangkut modul
>   lain). Di modul ini jabatan hanya dibaca dari satu sumber: JABATAN roster FACIN lewat baris tangga.
> - AC 80 *"**Susunan jenjang komite** — jabatan apa saja, pada urutan berapa — disimpan **terpisah** dari tabel login."*
>   dan AC 81 *"Mengubah susunan jenjang adalah **mengubah baris data**, ⛔ **bukan rilis**."* → **isinya tetap benar**,
>   tetapi sumbernya kini XML + KCF-01, bukan K13: susunan = roster `EMAILKOMITE` STS_KLAIM FACIN (DEGREE / JABATAN /
>   LIMIT, `ApprovalKomite_Act` L4), dan perubahannya = UPDATE baris roster (pola migrasi claimprop 537 / claimnonprop
>   611).
> - AC 82 *"Perubahan kolom `jabatan` **masuk jejak audit**: siapa mengubah, kapan, dari kode apa ke kode apa."* dan
>   AC 83 *"⭐ **Alasan AC 82 dinyatakan:** begitu jabatan menentukan siapa boleh menyetujui, ⛔ **mengubah satu baris
>   berarti mengubah wewenang.** Sejalan **ADR-0007** dan **ADR-0014**."* → **gugur** bersama kolom jabatan. Wewenang
>   kini berubah lewat keanggotaan workbasket (Kelola User, di luar modul ini). Rujukan ADR-0007 dan ADR-0014 kini di
>   `OUTPUT_HASIL_RNM/docs/bersama/adr/`.
>
> Tiket 11 yang menutup bab ini **gugur**; AC 80 / 81 / 84 diliput migrasi roster (KCF-01) dan tiket 03.

**AC 84** `[keputusan work owner]` **K12** Roster komite berbasis **jabatan**. ⛔ Nama orang **tidak
menjadi kunci** di mana pun.

**AC 85** ⭐ `[keputusan work owner]` **K12** ⭐ **Alasan dinyatakan work owner:** pemetaan akun
menjadi nama orang lain ada **karena persetujuan komite dikirim memakai nama**. ⭐ Memakai jabatan
**menghapus akarnya**, bukan menambal gejalanya.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⭐ Memakai jabatan **menghapus akarnya**, bukan
> menambal gejalanya."* → akarnya kini dihapus dengan **workbasket** (KCF-01). `KomiteID` tangga = workbasket, dan
> pemutus = anggotanya, sehingga tidak ada lagi akun yang perlu dipetakan ke akun lain. Pemetaan akun di
> `SetProteksiSubmiteKomite` L1 (akun orang tertulis mati) **dibuang** (prompt §5 #1 dan §3 "Gugur").

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

> ⛔ **RALAT 10-10-2026.** AC lamanya dikutip utuh, tidak dihapus: AC 89 *"Catatan kronologi disusun dari **medan
> tersimpan** — akun, kode jabatan, keputusan, waktu — ⛔ **bukan disimpan sebagai kalimat jadi**."* dan AC 90 *"⭐
> **Alasan AC 89 dinyatakan:** justru karena sistem lama menyimpan **kalimat jadi**, orang yang sama tercetak berbeda
> antar modul, dan nilai bawaan diberikan kepada yang tak dikenal. ⛔ Dengan medan terpisah, keduanya **tidak mungkin
> terjadi**."* →
>
> - **Kronologi = teks di `T_VIEW_SUGGEST`**, ikut XML dan pola claimfacin (OQ-CFI-19). `KomitePost_Adjustment` S3 / S4
>   (juga `KomitePost_Reject` S3 / S4, `KomitePost_CloseClaim` S2 / S3) menyusun `Data.CARI12` = `"Accepted by " +
>   jabatan + " - " + nomor KMT` atau `"Rejected by " + …`, dengan jabatan = `KomiteList(KomiteCount).IDKomite` (JABATAN
>   roster). `ChronologyInsertion_DT` lalu menambahkan satu baris ke kronologi **kasus klaim induk**. Ditulis lewat
>   kontrak di transaksi Submit.
> - Alasan AC 90 gugur untuk modul ini: jabatan diambil dari baris tangga, sehingga jabatan berbeda atau bawaan tidak
>   muncul. Rantai `@if` + bawaan direksi adalah sisa editor (RALAT lintas-ronde #2). Pengecualian `pyPosition != "IT
>   Developer"` di `ChronologyInsertion_DT` dibuang: semua peran terekam (pola claimfacin).

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

> ⛔ **RALAT 10-10-2026 — BAB N DIBACA ULANG (KCF-04).** AC lamanya dikutip utuh, tidak dihapus: AC 91 *"Catatan
> akseptasi historis MBU dan Travel **tidak dibangun ulang**."*, AC 92 *"Catatan itu **ditandai** dengan kolom
> **`DIBANGUN_ATURAN_LAMA`**, `CHAR(1)`, bernilai `'1'` / `'0'`."*, AC 93 *"⭐ **Alasan dinyatakan:** penandanya **kolom,
> bukan catatan prosa**, supaya **laporan dapat menyaringnya** — ⛔ bukan hanya pembaca manusia yang kebetulan membaca
> catatan kaki."*, AC 94 *"Cacah baris terdampak **belum diketahui**."* dan AC 95 *"Catatan baru **tidak** menerima
> penanda itu."* → **KCF-04** (10-10-2026): kasus komite lama Pega (`KMT-`) tidak dimigrasi; **K8
> `DIBANGUN_ATURAN_LAMA` gugur**. AC 91 tetap dalam arti "tidak ada yang dibangun ulang". AC 92–95 **gugur**: tanpa
> kolom penanda, tanpa DDL, dan butir 29 (cacah baris) tidak dibutuhkan.

### O · Batas dan titik sambung

**AC 96** `[terbuka]` **butir 1** `[data DBA]` Kolom tabel data kutipan **belum diketahui**.
⭐ **Dikurung:** ia menyentuh **bentuk tabel**, ⛔ bukan perilaku — dan tahap struktur tabel
**belum dimulai** sesuai `[keputusan work owner]` 2026-09-19.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Kolom tabel data kutipan **belum diketahui**."*
> dan *"⭐ **Dikurung:** ia menyentuh **bentuk tabel**, ⛔ bukan perilaku — dan tahap struktur tabel **belum dimulai**
> sesuai `[keputusan work owner]` 2026-09-19."* → **tertutup untuk modul ini.** Data polis dan kutipan dibaca dari
> `JSON_POLIS.DATA_JSON` (RALAT claimfacin AC 15) lewat kontrak `kontrak.KlaimFacInKomite`, jadi komite tidak butuh
> tabel kutipan sendiri. Tahap struktur sudah berjalan: `modul/claimfacin/docs/STRUKTUR-TABEL-CLAIM-FACIN.md` (§5 sisi
> komite) dan perubahan `T_GENERAL_KOMITE` KCF-03 di migrasi komiteclaimfacin 640–679.

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

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"| ⚠️ **9** | ⚠️ **Larangan menyetujui klaim
> sendiri untuk jenjang kedua ke atas** | ⛔ **Sengaja TIDAK dibangun** — `[keputusan work owner]` **K9**. ⭐ Ditulis di
> sini **agar tidak terbaca sebagai kelalaian** |"* dan *"| **10** | ⭐ **Struktur tabel dan relasinya** | `[keputusan work
> owner]` 2026-09-19 — dikerjakan **bersamaan dengan Claim Fac In**, sesudah spec ini |"* →
>
> - **Baris 9**: larangan menyetujui klaim sendiri tidak dibangun di **jenjang mana pun**, termasuk jenjang pertama.
>   XML tidak memuatnya (K9 salah baca, RALAT ID-5). Tanpa larangan rangkap (pola Komite Prop 09-10-2026).
> - **Baris 10**: sudah dikerjakan. `modul/claimfacin/docs/STRUKTUR-TABEL-CLAIM-FACIN.md` §5 (sisi komite, nol tabel
>   baru). Satu-satunya perubahan tabel komite adalah KCF-03 (`MODIFY T_GENERAL_KOMITE.ADJUSTMENT_ID` boleh kosong +
>   kolom `TRANSFER_TYPE`) dan roster FACIN → workbasket (KCF-01, UPDATE baris), keduanya di migrasi komiteclaimfacin
>   640–679.
> - **Baris 1** ("pintu belakang") dan **baris 7** ("pemetaan akun") tetap tidak dibangun, kini juga berdasar prompt
>   tahap 2 §3 "Gugur" dan §5 #1.
> - **Baris 8**: tetap tidak dibangun ulang, tetapi kini juga **tidak ditandai** (KCF-04, K8 gugur).

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

> ⛔ **RALAT 10-10-2026.** Kalimat dan baris lamanya dikutip utuh, tidak dihapus: *"⭐ **Satu pemblokir tersisa: butir
> 6**, dan ia menahan **PEMBANGUNAN**, ⛔ bukan penulisan."* dan *"| ⛔ **6** | **Isi daftar jabatan dan susunan
> jenjang** | `[work owner]` | ⛔ **YA — PEMBANGUNAN** |"* → keadaan register 10-10-2026, nomor lama tetap:
>
> - **6** tertutup: ADR-0030 + keanggotaan workbasket + KCF-01 (RALAT AC 29). **Nol pemblokir.**
> - **39** tertutup: pola Komite Prop 09-10-2026, setiap anggota workbasket tingkat berjalan boleh memutus (RALAT AC 7).
> - **1** tertutup untuk modul ini: data dibaca lewat kontrak dari `JSON_POLIS.DATA_JSON` (RALAT AC 96).
> - **25** (jalur pengganti pintu belakang) gugur bersama pintu belakang "IT Developer" (prompt §3 "Gugur").
> - **26**: untuk rancangan, `=` tunggal dibaca pembandingan (prompt §5 #2). Yang terbuka = dedupe kasir OQ-CFI-26.
> - **29** gugur bersama K8 (KCF-04).
> - **31** tidak dibutuhkan: penomoran di aplikasi (`inti/backend/penomor`, ADR-0043); format nomor terbaca dari
>   `KomitePost_Adjustment` S7.2.1.2.4–7 (RALAT tiket 07).
> - **11 · 12 · 13 · 15 · 27 · 28 · 30** tidak dinilai ulang di sini. Angka pita SPV B dan batas roster kini konstanta
>   bernama (KCF-01), tetapi asal-usulnya (butir 15) tetap terbuka.

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

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"⛔ **Pembangunan menunggu butir 6.**"* → butir 6
> tertutup (RALAT AC 29). Pembangunan dimulai 10-10-2026 atas perintah work owner (prompt "IMPLEMENTASI KOMITE CLAIM FAC
> IN, TAHAP 2 DARI 2"), dengan XML sebagai patokan dan keputusan KCF-01–04.

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

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"| `docs\adr\` — 15 berkas, termasuk ADR-0014 | ✅
> nol disunting |"* → ADR kini tinggal di **`OUTPUT_HASIL_RNM/docs/bersama/adr/`**: 43 catatan (0001–0043) + `00-INDEKS.md`.
> Yang dirujuk spec ini: ADR-0007 (jejak audit), ADR-0014 (pemutus komite per `KomiteID` tingkat berjalan), ADR-0015
> (efek keluar komite lewat outbox), ADR-0029 (batas transaksi), ADR-0030 (aturan peran), ADR-0043 (penomoran di
> aplikasi, menggantikan ADR-0006). Tiga ADR bersumber spec ini: ADR-0038 (ID-3), ADR-0039 (ID-2 / ID-6), dan ADR-0040
> (ID-5). Ketiganya **tidak disunting** oleh RALAT ini; ADR-0040 dan sebagian ADR-0038 / ADR-0039 bertentangan dengan
> RALAT ID-3 / ID-5 / ID-6 (pertanyaan untuk work owner). Kolom "Keadaan" tabel ini tetap fakta 20-09-2026.
| `OUTPUT_HASIL_RNM\CLAUDE.md` | ✅ nol disunting |
| korpus 114 · 482 · 80 · 329 | ✅ nol disunting |

⛔ Kode Go/React **NOL** · `CREATE TABLE` **NOL** · DDL **NOL** · nomor baris XML **NOL** ·
nilai kredensial **NOL** · tiket **NOL** · revisi ADR **NOL** · butir `[terbuka]` ditutup **NOL** ·
keputusan work owner baru **NOL**.
