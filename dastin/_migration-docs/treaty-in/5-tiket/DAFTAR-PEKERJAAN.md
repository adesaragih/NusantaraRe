# Daftar pekerjaan — Treaty In, penyerahan pertama

**Tanggal:** 24 September 2026 · **Langkah:** 2 dari urutan kerja.

> ### INI BUKAN TIKET, DAN BUKAN TAKSIRAN.
>
> Berkas ini mendaftar **hal yang harus ada di sistem baru**, dengan asal-usulnya. Ia **belum**
> dibentuk menjadi tiket dan **belum** ditaksir.
>
> ~~sengaja, karena pembentukan tiket menuntut `SPEC-MODEL-DATA.md` §10 selesai untuk entitas yang
> terpakai, dan §10 belum selesai.~~
>
> **DIPERBARUI 24 September 2026.** **§10 SELESAI** — 28 entitas, dan §11.1 menyatakannya. Penahan
> itu **jatuh**. Penahan yang **benar-benar tersisa** ada di §1.1.
>
> **Justru daftar inilah yang menentukan entitas mana yang perlu §10.** Lihat §5.

**Dasar:** `KEPUTUSAN-PEMBAGIAN-TIKET.md` (gerbang G1, G2, G3), `SPEC-MODEL-DATA.md`,
`SPEC-INVARIAN.md`, ADR-0034…0055, `4-erd-dan-tabel-datar/STRUKTUR-DATA.md`.

---

## 1. Cara membaca

Satu baris = **satu kemampuan**, menurut G1: satu hal yang dapat dilakukan **pelaku bernama**, dari
ujung ke ujung. Belum tentu satu tiket — sebagian akan dipecah menurut G3 saat dibentuk nanti.

| Kolom | Isinya |
|---|---|
| **Kode** | `P-nn`, tetap. Dipakai `PETA-LIPUTAN.md` untuk menunjuk |
| **Kemampuan** | berbunyi *"pelaku X dapat Y"*. Kalimat ini **tidak memerlukan satu pun nama atribut** — itu sebabnya daftar ini dapat ditulis sebelum §10 selesai |
| **Gol** | **L** = PELESTARIAN · **U** = PERUBAHAN · **B** = BARU |
| **Asal** | bagian spesifikasi, ADR, atau invarian yang mewajibkannya |
| **Entitas** | entitas yang **tersentuh** — kolom inilah yang dijumlahkan di §5 |
| **Keadaan** | **tidak terhalang luar** atau **TERTAHAN** (dengan penghalangnya). Lihat peringatan di §1.1 — kata ini **tidak** berarti "boleh dikerjakan" |

Pelaku yang dipakai, seluruhnya di luar tiket (aturan U-1):
**PK** pengisi kontrak · **SH** section head · **DH** dept head · **DR** direktur ·
**PM** pelaksana migrasi · **PJ** pemeriksa jejak · **KH** konsumen hilir.

### 1.1 "Tidak terhalang luar" BUKAN berarti "boleh dikerjakan"

> **TIDAK SATU PUN baris di berkas ini siap dibentuk menjadi tiket hari ini.**
>
> ~~Seluruhnya masih menunggu **§10 untuk 23 entitas** lalu **sesi DDL**~~ — **penahan ini jatuh
> 24 September 2026.** Kolom **Keadaan** hanya menyatakan ada-tidaknya penghalang **DI LUAR**
> urutan itu.

#### Penahan yang BENAR-BENAR tersisa — diperbarui 24 September 2026

| # | Penahan | Siapa mencabutnya | Kenapa ia menahan pembentukan tiket |
|---|---|---|---|
| ~~1~~ | ~~butir 3 urutan wewenang masih DICORET~~ | — | **DICABUT sebagai penahan 24 September 2026.** Butir 3 **tetap dicoret**, tetapi kebutuhannya dipenuhi **butir 5** urutan wewenang: soal nama kolom dan tipe, `2-to-spec/KAMUS-KOLOM.md` mengikat — **240 kolom bernama**, dibangkitkan dari §10 sehingga tidak dapat menyimpang darinya |
| 2 | **angka presisi per kelompok tipe** belum diputuskan | gerbang sesi DDL | tipe kolom di `ddl-usulan/` masih **mewarisi** modul tetangga sebagai usulan |
| 3 | **bentuk fisik dua induk polimorfik** — `POTONGAN`, `PENYEBARAN` | gerbang sesi DDL | `INV-17` melarang rujukan yang sasarannya bergantung nilai kolom lain; bentuknya menentukan jumlah kolom |
| 4 | **enam paket uang golongan C** berdiri tanpa denominasi | **teknik treaty**, `T-6` | constraint `INV-39`/`INV-40` tidak dipasang untuk keenamnya |
| 5 | paket **`REV-1` … `REV-6`** belum ditanggapi | pemilik ADR induk | lima ADR yang tiket akan bersandar padanya |

> **Yang TIDAK lagi menahan: §10, dan sejak 24 September juga butir 3 urutan wewenang.** §10 Ia selesai untuk **28 entitas**, dan `2-to-spec/KAMUS-KOLOM.md`
> kini memberi **240 kolom bernama** — sehingga kriteria selesai yang menyebut nama field **sudah
> dapat ditulis**. Yang tersisa bukan ketiadaan nama, melainkan **ketiadaan pemutus** bila nama itu
> bertengkar dengan bentuk tabel.

Kata "siap" dipakai di bentuk awal berkas ini dan **ditarik** pada 24 September 2026, karena ia
mengerjakan dua pekerjaan sekaligus: yang dimaksud "tidak terhalang jawaban dari luar", yang
terbaca "boleh mulai".

Contoh yang menunjukkan bedanya bukan soal kata-kata: **P-53** — pelaksana migrasi dapat mematikan
penegakan trigger lalu menyalakannya kembali. Kriteria selesainya harus menyebut **akun mana** yang
boleh mematikan dan akun mana yang tidak. Itu rancangan hak akses, dan ia milik **sesi DDL yang
belum berjalan**. P-53 tidak terhalang dari luar, dan tetap tidak dapat ditulis sekarang.

---

## 2. Kemampuan — **66 didaftar, 63 aktif** *(dihitung ulang 24 September 2026, gerbang pembuka to-ticket)*

> ### CACAH DIPERBAIKI, DAN ANGKA LAMANYA DISEBUT
>
> Judul ini pernah berbunyi *"60 didaftar, 57 aktif"* dan §4 berbunyi *"58 didaftar, 55 aktif"* —
> **dua angka berbeda untuk hal yang sama, di satu berkas**. Keduanya dihitung ulang mekanis:
>
> | | |
> |---|---:|
> | kode yang **punya baris tabel** | **65** — `P-01`…`P-58` dan `P-60`…`P-66` |
> | **`P-59` — tidak punya baris tabel sama sekali**, diputuskan di §2.2b | 1 |
> | **kemampuan seluruhnya** | **66** |
> | dicoret sebagai SIFAT (U-6) | 3 — `P-24`, `P-25`, `P-26` |
> | **aktif** | **63** |
>
> **Baris untuk `P-59` tidak dikarang.** Mengisi kolom Gol, Asal, dan Entitas-nya berarti memutuskan
> hal yang tidak pernah ditetapkan siapa pun; isinya ada di §2.2b dan itu yang berlaku.

### A. Kontrak dan versinya ada

| Kode | Kemampuan | Gol | Asal | Entitas | Keadaan |
|---|---|---|---|---|---|
| P-01 | **PK** dapat membuat kontrak baru dan melihatnya kembali dengan pengenalnya; versi pertamanya lahir langsung dalam keadaan `DRAFT` | U | ADR-0040, ADR-0055 (`LAHIR`), INV-01…04, INV-20 | `KONTRAK`, `VERSI_KONTRAK` | tidak terhalang luar — **DIBELAH ANTAR BATCH 24 Sep 2026**, lihat §2.4 |
| P-02 | **PK** diperingatkan bahwa kontrak yang diisinya berkunci alami sama dengan kontrak lain, **tanpa dihalangi menyimpan** | B | ADR-0040 §2 | `KONTRAK` | tidak terhalang luar |
| P-03 | **PK** dapat menemukan kembali kontrak memakai **nomor warisan** sistem lama, berdampingan dengan pengenal baru | U | ADR-0042 | `KONTRAK` | tidak terhalang luar |
| P-04 | **PK** dapat mengisi dan mengubah seluruh kepala versi selama versinya `DRAFT`; perubahan atas **jenis** maupun **materialitas** meninggalkan baris jejak, dan **sesudah diajukan keduanya beku** — mengubahnya menuntut versinya dikembalikan ke `DRAFT` lebih dulu | U | §10.2, GRL-18, KTV-1, ADR-0045 | `VERSI_KONTRAK`, `JEJAK_PERUBAHAN` | tidak terhalang luar — **DIUBAH 24 Sep 2026**, dari *"dapat diubah selama `DRAFT`"* yang diam tentang apa yang terjadi sesudahnya |
| P-05 | **PK** tidak dapat mengubah kunci alami kontrak sesudah kontraknya lahir, dan penolakannya menyebut apa yang ia coba ubah | U | INV-19 (trigger), ADR-0040 | `KONTRAK` | tidak terhalang luar |

### B. Daftar yang menggantung pada versi

| Kode | Kemampuan | Gol | Asal | Entitas | Keadaan |
|---|---|---|---|---|---|
| P-06 | **PK** dapat mencatat **lebih dari satu** mata uang berlaku pada satu kontrak, masing-masing dengan kurs dan periodenya | U | ADR-0053, INV-07, INV-44 | `MATA_UANG_KONTRAK`, `MATA_UANG` | tidak terhalang luar |
| P-07 | **PK** dapat mencatat retensi cedant per kelompok treaty per mata uang | L | INV-08 | `RETENSI_CEDANT`, `KELOMPOK_TREATY` | tidak terhalang luar |
| P-08 | **PK** dapat mencatat EGNPI per kelompok treaty per mata uang | L | INV-09 | `EGNPI`, `KELOMPOK_TREATY` | tidak terhalang luar |
| P-09 | **PK** dapat mencatat portofolio masuk dan keluar yang menyertai kontrak | L | §2.3 | `PORTOFOLIO` | tidak terhalang luar |
| P-10 | **PK** dapat mencatat periode pelaporan beserta tanggal jatuh temponya, dan ditolak bila periodenya di luar periode kontrak | L | INV-10, INV-55 | `PERIODE_PELAPORAN` | tidak terhalang luar |
| P-11 | **PK** dapat mencatat periode akumulasi, dan ditolak bila di luar periode kontrak | L | INV-11, INV-56 | `PERIODE_AKUMULASI` | tidak terhalang luar |
| P-12 | **PK** dapat mencatat termin pembayaran premi bernomor urut | L | INV-12 | `TERMIN` | tidak terhalang luar |
| P-13 | **PK** dapat mencatat skala ko-asuransi sebagai **beberapa baris**, bukan satu nilai | U | INV-13 | `SKALA_KOASURANSI` | tidak terhalang luar |
| P-14 | **PK** dapat mencatat batas tanggungan untuk **bahaya apa pun yang ada di daftar bahaya**, termasuk bahaya yang ditambahkan kemudian | U | INV-14, ADR-0038 | `BATAS_PER_BAHAYA`, `BAHAYA` | tidak terhalang luar |
| P-15 | **PK** dapat melampirkan dokumen pada kontrak, dan dokumennya **dirujuk, bukan disalin** | L | ADR-0027, INV-59 | `DOKUMEN_KONTRAK` | tidak terhalang luar |

### C. Cabang non-proporsional

| Kode | Kemampuan | Gol | Asal | Entitas | Keadaan |
|---|---|---|---|---|---|
| P-16 | **PK** dapat mencatat layer beserta limit, deductible, MDP, dan ketentuan pemulihan limitnya | L | §10.3, INV-05 | `LAYER` | tidak terhalang luar |
| P-17 | **PK** dapat mencatat bagian NuRe atas sebuah layer, dan bagian hanya dapat dibuat pada kontrak non-proporsional | L | INV-33 | `BAGIAN` | tidak terhalang luar |
| P-18 | **PK** dapat menyebarkan bagian NuRe ke susunan retro internal **per jenis reasuransi** per mata uang, dan penyebaran yang tidak berjumlah seratus persen **ditolak** | **L?** — lihat §2.3 | INV-16, INV-47, INV-50, INV-51 | `PENYEBARAN`, `RINCIAN_PENYEBARAN`, `NILAI_PENYEBARAN`, `JENIS_REASURANSI` | **TERTAHAN — DUA penghalang**: (1) INV-50/51 belum terbukti (L-3); (2) **golongannya belum pasti** — uji yang sama dengan `RINCIAN_PENYEBARAN`, §2.3 |
| P-19 | **PK** dapat mencatat potongan atas premi, dengan aturan potongan yang **sama** di kedua tempat ia melekat | L | INV-15, INV-63, §14.1 | `POTONGAN`, `JENIS_POTONGAN` | tidak terhalang luar |
| P-20 | **PK** dapat mencatat cara pembukuan XOL, dan hanya pada kontrak non-proporsional | L | INV-34, §14.4 | `VERSI_KONTRAK` | **TERTAHAN** — menunggu Uji X-2 |

### D. Cabang proporsional

| Kode | Kemampuan | Gol | Asal | Entitas | Keadaan |
|---|---|---|---|---|---|
| P-21 | **PK** dapat mencatat ketentuan proporsional per kelompok treaty di dalam sebuah layer | L | §10.4, INV-06, INV-32 | `DETAIL_PROPORSIONAL`, `KELOMPOK_TREATY` | tidak terhalang luar |
| P-22 | **PK** mengisi **tepat satu** dari persen quota share atau jumlah lines surplus, ditentukan jenis treaty-nya; mengisi keduanya atau tidak satu pun **ditolak** | L | INV-30, INV-31 | `DETAIL_PROPORSIONAL` | tidak terhalang luar |
| P-23 | **PK** yang mencatat baris surplus tanpa baris quota share pada versi yang sama **menerima kegagalan yang menyebutkan apa yang kurang** — bukan nilai nol | U | INV-35, ADR-0035 | `DETAIL_PROPORSIONAL` | tidak terhalang luar |

### E. Paket uang

> **Tiga baris di sini DICORET sebagai kemampuan menurut U-6.** Lihat §2.1 di bawah.

| Kode | Kemampuan | Gol | Asal | Entitas | Keadaan |
|---|---|---|---|---|---|
| ~~P-24~~ | ~~Setiap nilai uang membawa mata uangnya~~ — **SIFAT, bukan kemampuan** | — | INV-36 | — | **DICORET → §2.1** |
| ~~P-25~~ | ~~Setiap nilai IDR membawa kurs, tanggal kurs, dan sumber kurs~~ — **SIFAT** | — | INV-37, INV-38 | — | **DICORET → §2.1** |
| ~~P-26~~ | ~~Setiap besaran uang menyatakan tingkat pencatatannya~~ — **SIFAT** | — | INV-39, INV-40 | — | **DICORET → §2.1** |
| P-27 | **PK** melihat angka rupiah yang **sama** pada kontrak yang sudah disetujui, hari ini maupun tahun depan — kursnya dibekukan pada versinya, tidak dibaca ulang saat ditampilkan | U | ADR-0036, INV-43 | `VERSI_KONTRAK`, `MATA_UANG_KONTRAK` | tidak terhalang luar |
| P-28 | **PK** yang perhitungannya gagal melihat **keterangan apa yang gagal**, bukan angka nol dan bukan kurs satu | U | ADR-0035, INV-38, INV-42 | `VERSI_KONTRAK` *(tempat keterangannya)* | tidak terhalang luar |

#### 2.1 Sifat paket uang — dibawa setiap tiket, bukan dikerjakan sendiri

P-24, P-25, dan P-26 gugur uji U-6: **tidak ada satu pelaku yang "dapat melakukan" P-24**, dan
tidak ada titik ia dapat dinyatakan selesai. Kolom entitasnya berbunyi *"seluruh entitas berpaket
uang"* — itu tandanya.

Perlakuannya, menurut U-6:

| | |
|---|---|
| **INV-36, INV-37, INV-38, INV-39, INV-40, INV-41, INV-42, INV-43** | masuk medan **INVARIAN YANG HARUS DIPENUHI** pada **setiap** tiket yang menyimpan uang — tanpa kecuali |
| **bentuk fisik paket uang** | dibangun tiket ber-**PEMBUAT PERTAMA** untuk entitas berpaket uang pertama, bersama tabel dan constraint-nya |

P-27 dan P-28 **dipertahankan** karena keduanya lolos uji U-6: masing-masing punya pelaku nyata
(**PK**) dan akibat yang ia lihat sendiri — angka yang tidak bergeser, dan keterangan kegagalan
alih-alih nol. Keduanya ditulis ulang dari bentuk "sistem menjamin X" menjadi "PK melihat X",
karena bentuk pertama tidak dapat dinyatakan selesai dan bentuk kedua dapat.

### F. Siklus hidup dan persetujuan

| Kode | Kemampuan | Gol | Asal | Entitas | Keadaan |
|---|---|---|---|---|---|
| P-29 | **PK** dapat mengajukan versi `DRAFT` untuk persetujuan, dan pengajuan yang tidak memenuhi **K1-1 dan K1-2** ditolak dengan menyebut apa yang kurang | L | ADR-0055 (`AJUKAN`), INV-27, K1-1…2 | `VERSI_KONTRAK` | tidak terhalang luar |
| P-30 | Pengajuan yang tidak memenuhi **K1-3 sampai K1-8** menghasilkan **peringatan yang tercatat**, belum penolakan | **B** | INV-27, K1-3…8, CO-6…CO-8 | `VERSI_KONTRAK` | tidak terhalang luar — wajib mengisi **CARA MENYALAKANNYA** (4 butir) |
| P-31 | **SH** dapat menyetujui versi yang menunggunya, dan versinya berpindah ke antrian **DH** | L | ADR-0055 (`SETUJUI`) | `VERSI_KONTRAK`, `CATATAN_PERSETUJUAN` | tidak terhalang luar |
| P-32 | **DH** dapat menyetujui versi yang menunggunya, dan versinya berpindah ke antrian **DR** | L | ADR-0055 | `VERSI_KONTRAK`, `CATATAN_PERSETUJUAN` | tidak terhalang luar |
| P-33 | **DR** dapat menyetujui versi yang menunggunya, dan versinya menjadi `DISETUJUI` — hanya bila **seluruh K2** terpenuhi | L | ADR-0055, K2-1…6 | `VERSI_KONTRAK`, `CATATAN_PERSETUJUAN` | tidak terhalang luar |
| P-34 | **SH / DH / DR** dapat mengembalikan versi ke `DRAFT`, dan pengembalian tanpa alasan ditolak | L | ADR-0055 (`KEMBALIKAN`), §3 | `VERSI_KONTRAK`, `CATATAN_PERSETUJUAN` | tidak terhalang luar |
| P-35 | **SH / DH / DR** dapat menolak versi **kontrak**, dan penolakan tanpa alasan ditolak | L | ADR-0055 (`TOLAK`) | `VERSI_KONTRAK`, `CATATAN_PERSETUJUAN` | tidak terhalang luar |
| P-58 | **SH / DH / DR** dapat menolak versi **addendum**, dan penolakannya **meninggalkan catatan** | **B** | ADR-0045, ADR-0055 (`TOLAK`) | `VERSI_KONTRAK`, `CATATAN_PERSETUJUAN` | tidak terhalang luar — wajib **CARA MENYALAKANNYA** |

> ### Dasar P-58 DIPERBAIKI 24 September 2026 — dan ke arah yang lebih berat
>
> Sampai hari ini dasarnya berbunyi *"sistem lama tidak pernah mencatat penolakan addendum"*. Itu
> **kurang tepat**. `Activity/TreatyInDeclineConfirmation_postactEDM.xml` dibaca langkah demi
> langkah:
>
> ```
> 1  Property-Set   Param.Info / Param.Comment                        MATI
> 2  Call AddCommentList_Act                                          MATI
> 3  Property-Set   CommentList(<LAST>).IsApproved = "Decline"        MATI
>                   StatusAkseptasi               = "Decline"         MATI
> 4  Call SaveTreatyIn_Act                                            MATI
> ------------------------------------------------------------- hidup:
> 5  Property-Set   InputData.CARI1 = TreatyIn.ID
> 6  RDB-List       RDB remove EDM
> 7  RDB-List       RDB remove EDM
> 8  call TreatyInInputVis
> ```
>
> **KODE PENCATATANNYA ADA, DITULIS ORANG, LALU DIMATIKAN — sementara kode penghapusannya
> dibiarkan hidup.** Bukan kelalaian, melainkan **penggantian**.
>
> **Golongan P-58 TETAP BARU** — kemampuan yang mati bukan pelestarian, dan tidak ada satu baris
> data pun yang membuktikan ia pernah berjalan. Yang berubah **dasarnya**, dan dasar yang tepat
> lebih kuat di hadapan orang yang berwenang: yang diminta bukan *"buatkan yang belum pernah ada"*
> melainkan *"hidupkan kembali yang pernah ditulis dan dimatikan"*.
>
> Bukti: `4-erd-dan-tabel-datar/CABANG-MATI-DI-JALUR-PERSETUJUAN.md` §2.2.
> Golongan dasarnya **`EVIDENCED`**, dengan nomor langkah.
> **P-58 dipisah dari P-35 pada 24 September 2026, dan golongannya BARU — bukan PELESTARIAN.**
> Ini mudah terlewat: karena penolakan **kontrak** memang sudah meninggalkan jejak di sistem lama,
> penolakan addendum yang berjejak tampak seperti pelestarian. Ia bukan.
>
> Dua jalur sistem lama diperlakukan **berbeda**, tanpa alasan bisnis yang masuk akal:
>
> | Jalur | Yang terjadi saat ditolak |
> |---|---|
> | `TreatyInDeclineConfirmation_postact` — kontrak | 4 langkah hidup, **hanya menyetel keadaan**. Barisnya bertahan |
> | `TreatyInDeclineConfirmation_postactEDM` — addendum | langkah 6 dan 7 **hidup**, **menghapus baris** di `M_TREATY_IN_EDM` dan `TREATY_IN_EDM` |
>
> **Untuk addendum, mencatat penolakan adalah kemampuan BARU**, karena sistem lama tidak pernah
> melakukannya. Menghitungnya sebagai pelestarian akan membuat taksirannya terlalu rendah dan
> pengujiannya salah sasaran — pelestarian diuji terhadap data lama, dan **data lamanya tidak ada**.
>
> **Akibat yang ditarik satu langkah lebih jauh:** karena penolakan addendum tidak pernah berjejak,
> **tidak ada seorang pun yang pernah dapat menjawab berapa kali sebuah kontrak gagal diubah.**
> Itu pertanyaan yang biasanya ditanyakan auditor. Dicatat di `DAFTAR-ESKALASI-MANAJEMEN.md` butir 4.

| P-36 | Pemberi keputusan yang **sama dengan pengisi kontraknya**, atau sama dengan pemberi keputusan tingkat sebelumnya, **ditolak** | U | INV-27, ADR-0052, §3 | `VERSI_KONTRAK`, `CATATAN_PERSETUJUAN` | tidak terhalang luar |
| P-37 | Perpindahan di luar **tiga belas** yang sah **ditolak**, termasuk setiap perpindahan keluar dari `DISETUJUI` dan `DITOLAK` | U | INV-22, INV-23 (trigger), ADR-0055 | `VERSI_KONTRAK` | tidak terhalang luar |
| P-38 | Versi berkeadaan `DISETUJUI` **tidak dapat diubah nilainya oleh siapa pun**, dan percobaannya ditolak | U | INV-24 (trigger), ADR-0036 | `VERSI_KONTRAK` dan seluruh anaknya | tidak terhalang luar |
| P-39 | Satu kontrak hanya dapat punya **satu** versi yang belum selesai; membuat versi kedua yang belum selesai ditolak | U | INV-25 (indeks unik) | `VERSI_KONTRAK` | tidak terhalang luar — **baca bersama §2.2: tanpa cara membuang draf, baris ini mengunci kontraknya** |
| P-40 | **PK** tidak dapat menyetujui kontrak tanpa melewati ketiga tingkat — tidak ada jalan pintas | U | ADR-0052, eskalasi butir 1 | `VERSI_KONTRAK` | tidak terhalang luar |
| P-41 | Wewenang persetujuan dibatasi **nilai kontrak** menurut batas yang berlaku | U | eskalasi butir 5 | `VERSI_KONTRAK` | **TERTAHAN** — dokumen batas wewenang belum ada; **tidak dapat dipecah** |

### G. Versi baru dan addendum

| Kode | Kemampuan | Gol | Asal | Entitas | Keadaan |
|---|---|---|---|---|---|
| P-42 | **PK** dapat membuat versi baru dari **versi berlaku terakhir**, dan versi baru itu mulai dari `DRAFT` serta melewati keempat tingkat; rujukan dasarnya **disimpan eksplisit** — wajib terisi pada versi penyesuaian, wajib kosong pada versi pertama | U | ADR-0055, ADR-0052, GRL-10, B-3, B-4 | `VERSI_KONTRAK` | tidak terhalang luar — **DIUBAH 24 Sep 2026**, dari *"dari versi yang sudah `DISETUJUI`"*, yang mengizinkan dasar yang bukan versi terakhir |
| P-43 | Lapisan beku — cedant, asal bisnis, sifat proporsi, tanggal mulai, tanggal berakhir — **tidak berubah antar versi** | U | INV-B1, ADR-0040 | `KONTRAK`, `VERSI_KONTRAK` | tidak terhalang luar |
| P-44 | **PK** melihat nilai versi sebelumnya berdampingan dengan versi yang sedang disusun, dan nilai lama itu **di-SELECT, tidak disalin** | U | INV-B2, `SEAM-ADJUSTMENT.md` §1 | `VERSI_KONTRAK` | tidak terhalang luar |
| P-45 | Tanggal berlaku **pada versi** yang jatuh di luar periode kontraknya **ditolak** — dan tanggal berlaku **pada dokumen addendum adalah ruas yang berbeda** (`P-62`) | L | INV-54 (trigger), GRL-15, KTV-2 | `VERSI_KONTRAK` | tidak terhalang luar — **DIUBAH 24 Sep 2026**: pembawanya disebut, sebab "tanggal berlaku" kini punya **dua** pembawa |
| P-57 | **PK** dapat membatalkan versi `DRAFT` yang ia buat sendiri, dan kontraknya langsung terbuka untuk versi berikutnya — **barisnya tetap tersimpan, nomornya tidak dipakai ulang** | **B** | ADR-0055 perubahan 24 Sep 2026, INV-04 | `VERSI_KONTRAK` | tidak terhalang luar — wajib **CARA MENYALAKANNYA** |

> **P-57 ditambahkan 24 September 2026, dan ia menutup L-7.** ADR-0055 memutuskan keadaan kedelapan
> `DIBATALKAN`, terminal, satu perpindahan masuk `DRAFT → BATALKAN → DIBATALKAN`.
>
> **Golongannya BARU, bukan PELESTARIAN**, dan itu bukan kehalusan penggolongan: sistem lama tidak
> punya cara membuang draf yang tidak merusak apa pun. Yang ia punya hanyalah menolak — dan untuk
> addendum, penolakan **menghapus barisnya**.
>
> **Penarikan sesudah diajukan TIDAK termasuk** dan belum diputuskan. Ia kemampuan tersendiri dengan
> pelaku dan akibat yang berbeda.

### H. Jejak

| Kode | Kemampuan | Gol | Asal | Entitas | Keadaan |
|---|---|---|---|---|---|
| P-46 | Setiap perpindahan keadaan meninggalkan satu catatan persetujuan bersama pelaku dan waktunya, **tanpa dapat dilewati** | U | INV-28 (trigger), ADR-0045 | `CATATAN_PERSETUJUAN` | tidak terhalang luar |
| P-47 | **PJ** dapat melihat siapa mengubah fakta apa dan kapan, untuk kontrak mana pun | **B** | ADR-0045, eskalasi butir 4 | `JEJAK_PERUBAHAN` | tidak terhalang luar |

### I. Himpunan acuan

| Kode | Kemampuan | Gol | Asal | Entitas | Keadaan |
|---|---|---|---|---|---|
| P-48 | Himpunan yang dapat bertambah tanpa mengubah arti — mata uang, jenis potongan, jenis reasuransi, bahaya, kelompok treaty, kelas bisnis — dapat **ditambah tanpa mengubah skema** | U | ADR-0038, INV-62 | `MATA_UANG`, `JENIS_POTONGAN`, `JENIS_REASURANSI`, `BAHAYA`, `KELOMPOK_TREATY`, `KELAS_BISNIS` | tidak terhalang luar |
| P-49 | Tabel acuan pembagian kapasitas punya **pemilik yang bertanggung jawab atas isinya** | U | ADR-0050, eskalasi butir 7 | *(belum ditetapkan)* | **TERTAHAN** — pemiliknya belum ditetapkan |

### J. Migrasi

| Kode | Kemampuan | Gol | Asal | Entitas | Keadaan |
|---|---|---|---|---|---|
| P-50 | **PM** dapat memindahkan kontrak lama apa adanya, **tanpa menghitung ulang apa pun** | U | ADR-0042, ADR-0043 | seluruh entitas | tidak terhalang luar — **payung, DIPECAH per kelompok entitas SEBELUM penaksiran**; **DIBELAH ANTAR BATCH 24 Sep 2026**, lihat §2.4 |
| P-51 | Kontrak lama yang keadaannya tidak punya padanan mendarat di `WARISAN_TAK_TERPETAKAN` dengan **nilai aslinya tetap tersimpan** | U | ADR-0054, INV-21, INV-26 | `VERSI_KONTRAK` | tidak terhalang luar |
| P-52 | **PK** dapat memindahkan kontrak ber-`WARISAN_TAK_TERPETAKAN` ke keadaan sah **yang dipilihnya secara eksplisit** | **B** | ADR-0054 | `VERSI_KONTRAK` | tidak terhalang luar |
| P-53 | **PM** dapat mematikan penegakan trigger selama pemindahan dan menyalakannya kembali, dan **keadaan sakelar itu terlihat** | **B** | §4.3b, CO-1…CO-8 | — *(mekanisme, bukan entitas)* | tidak terhalang luar — wajib **CARA MENYALAKANNYA** |

> **P-50 dipecah SEBELUM penaksiran, bukan sambil menaksir.** Kemampuan payung yang ditaksir
> sebagai satu akan **selalu** ditaksir terlalu rendah, karena yang ditaksir adalah **kalimatnya**,
> bukan isinya. Pemecahannya mengikuti kelompok entitas di §5.1.

### K. Hilir dan arsip

| Kode | Kemampuan | Gol | Asal | Entitas | Keadaan |
|---|---|---|---|---|---|
| P-54 | **KH** yang tidak dikenal tetap dapat membaca identitas kontrak **dalam bentuk lama** | U | ADR-0051 | bentuk baca atas `KONTRAK`, `VERSI_KONTRAK` | tidak terhalang luar |
| P-55 | **PM** dapat menyimpan dokumen JSON sistem lama sebagai arsip pada saat pemindahan, dan dapat menunjukkan arsip itu ada untuk kontrak mana pun | U | ADR-0034 | arsip | tidak terhalang luar |

> **P-55 ditulis ulang.** Bentuk lamanya — *"arsip tidak punya jalur baca"* — gugur uji U-6: ia
> **larangan**, bukan kemampuan, dan tidak ada pelaku yang "melakukannya". INV-61 karena itu menjadi
> **invarian yang dibawa** tiket P-55 dan setiap tiket yang menyentuh arsip, sementara kemampuannya
> yang nyata — **menyimpan** arsipnya — dipegang **PM** dan dapat dinyatakan selesai.

### L. Menemukan kembali

| Kode | Kemampuan | Gol | Asal | Entitas | Keadaan |
|---|---|---|---|---|---|
| P-56 | **PK** dapat mencari kontrak yang ada — menurut cedant, asal bisnis, periode, sifat proporsi, atau keadaan siklus hidupnya — dan melihat hasilnya sebagai daftar | L | §10.1, §10.2, ADR-0055 | `KONTRAK`, `VERSI_KONTRAK` | tidak terhalang luar |

### J. Kemampuan Adjustment — `P-60` … `P-66`

> **Disambung, bukan didaftar terpisah** (`GRL-01`: satu model, satu papan). Ketujuhnya lolos
> adjudikasi di [`ADJUDIKASI-KEMAMPUAN-ADJUSTMENT.md`](ADJUDIKASI-KEMAMPUAN-ADJUSTMENT.md) — dua
> calon lain bergolongan **SUDAH ADA** dan tiga **PERLU DIUBAH** (`P-04`, `P-42`, `P-45`, ditandai
> di tempatnya).

| Kode | Kemampuan | Gol | Asal | Entitas | Keadaan |
|---|---|---|---|---|---|
| P-60 | **PK** dapat menyatakan **materialitas** sebuah versi, dan materialitas itu **mengunci ruas mana yang boleh disunting** — dua arah — dengan penolakannya ditegakkan **di sisi simpan**, bukan hanya di layar | **B** | GRL-20, INV-69, INV-70 | `VERSI_KONTRAK` | **TERTAHAN** — `DIASUMSIKAN-CLEAR(DB-20)`; wajib **CARA MENYALAKANNYA** dan **pemantau kebasian** (`L-3`) |
| P-61 | **PK** dapat mencatat **dokumen addendum** bernomor sendiri, dan **satu dokumen dapat menyentuh beberapa kontrak** — persetujuannya tetap per kontrak | **B** | GRL-19, INV-71 | `DOKUMEN_ADDENDUM`, `VERSI_KONTRAK` | **TERTAHAN** — `DIASUMSIKAN-CLEAR(DB-16a)`; wajib **CARA MENYALAKANNYA** |
| P-62 | Dokumen addendum dapat membawa **tanggal berlakunya sendiri**, boleh kosong | **B** | KTV-2 | `DOKUMEN_ADDENDUM` | **TERTAHAN** — `DIASUMSIKAN-CLEAR(DB-16b)`; **bertenggat**: kolomnya dicabut **sebelum data masuk** bila dibantah |
| P-63 | **PK** melihat **baris selisih per besaran yang berubah**, dipadankan lewat **kunci bisnis termasuk mata uang** — bukan menurut posisi baris | U | ADR-0048 butir 3, E-1a, GRL-14, GRL-16 | `NILAI_SELISIH`, `VERSI_KONTRAK` | tidak terhalang luar — **PEMBUAT PERTAMA** tabel selisih |
| P-64 | Perubahan **persentase reinstatement menghasilkan** baris selisih | U | TDA-18 | `NILAI_SELISIH`, `PEMULIHAN_LIMIT` | tidak terhalang luar — **PERUBAHAN perilaku**, menuntut pemberitahuan pra-peralihan |
| P-65 | **PM** dapat memindahkan addendum warisan dengan **nomor urut diberikan ulang menurut kronologi**, sementara pengenal `/Rnn` **dilestarikan apa adanya** | U | GRL-17, ADR-0042, ADR-0043 | `VERSI_KONTRAK` | tidak terhalang luar — **pecahan payung `P-50`**; perubahan lebar, perluas–pindahkan–kerutkan |
| P-66 | **PK** dapat mengisi **nomor dokumen** baris warisan **dari arsip kertas** | **B** | GRL-19 §3.1 | `DOKUMEN_ADDENDUM` | tidak terhalang luar — **pekerjaan orang**, bukan sistem; migrasi tidak punya sumber |


> **P-56 ditambahkan 24 September 2026.** Ia sebelumnya tidak ada — P-03 hanya menemukan kembali
> lewat **nomor warisan**, yang mengandaikan penggunanya sudah tahu nomornya. Kemampuan mencari
> sendiri terlewat karena terasa terlalu jelas untuk didaftar.
>
> Di bawah irisan tegak, **yang terasa terlalu jelas untuk didaftar akan dibangun lima kali atau
> nol kali** — lima kali bila tiap tiket membuat pencariannya sendiri, nol kali bila tiap tiket
> mengira tetangganya yang mengerjakannya.

---

## 2.2 LUBANG DI ADR-0055 — tidak ada cara membuang draf  — **DITUTUP 24 Sep 2026**

> **Lubang ini ditemukan di sini, dilaporkan sebagai lubang, dan diputuskan pemilik proses pada
> hari yang sama.** ADR-0055 kini punya keadaan kedelapan `DIBATALKAN`. Kemampuannya **P-57**.
> Uraian di bawah dipertahankan apa adanya karena ia alasan keputusannya.

### Yang diperiksa

Dua belas perpindahan ADR-0055 disisir. **`DRAFT` punya tepat SATU perpindahan keluar: `AJUKAN`.**
Tidak ada `BATALKAN`, `BUANG`, `TARIK`, maupun `HAPUS` di seluruh dua belas.
*(Angka "dua belas" di sini adalah keadaan SEBELUM perbaikan; sesudahnya tiga belas.)*

### Akibatnya, dan ia nyata

Gabungkan dengan dua baris yang sudah ada di daftar ini:

| | |
|---|---|
| **P-39** | satu kontrak hanya boleh punya **satu** versi yang belum selesai (INV-25) |
| **P-37** | perpindahan di luar daftar yang sah **ditolak** (INV-22) |

Maka: **pengisi kontrak yang membuat versi karena salah pencet, atau memulai addendum yang
ternyata tidak jadi, tidak punya jalan keluar.** Kontraknya terkunci — versi yang belum selesai
tidak dapat dibuang, dan versi kedua tidak boleh dibuat.

Satu-satunya jalan yang tersisa: **mengajukannya ke persetujuan supaya ada yang menolaknya** —
memakai jalur persetujuan untuk membuang sampah, dan mengotori riwayat persetujuan dengan
pengajuan yang tidak pernah dimaksudkan serius.

### Sistem lama menanganinya bagaimana

Dijawab dari `PENGETAHUAN.md`, dan **keberadaan langkahnya diperiksa hidup-matinya** karena baris
itu tidak bertanda:

| Berkas | Keadaan |
|---|---|
| `TreatyInDeclineConfirmation_postact` | **4 langkah, seluruhnya hidup** — hanya menyetel keadaan; **tidak menghapus apa pun** |
| `TreatyInDeclineConfirmation_postactEDM` | **8 langkah, 4 hidup.** Langkah 6 dan 7 — *"RDB remove EDM"* — **HIDUP**: menghapus baris di `M_TREATY_IN_EDM` dan `TREATY_IN_EDM` |

**Sistem lama juga tidak punya "buang draf".** Yang ia punya:

- **kontrak baru** — decline saja; barisnya tetap ada, keadaannya menjadi ditolak;
- **addendum** — decline, dan penanganan decline-nya **menghapus barisnya**.

Artinya dugaan di atas bukan ramalan: **orang selama ini memang memakai jalur penolakan sebagai
tempat sampah**, dan untuk addendum sistem lama bahkan membuat buktinya lenyap.

### Kenapa lubang ini harus ditutup sebelum sesi to-spec melanjutkan

Di sistem lama, jalan keluarnya ada meski buruk. Di rancangan baru **jalan itu hilang**, dan tidak
ada penggantinya:

| | Sistem lama | Rancangan baru |
|---|---|---|
| buang addendum draf | decline → **baris dihapus** | tidak ada — `DITOLAK` terminal, dan ADR-0045 melarang jejak hilang |
| buang kontrak draf | decline → tercatat ditolak | hanya lewat ajukan-lalu-tolak |
| buat versi pengganti | tidak dibatasi | **dilarang** selama versi lama belum terminal (INV-25) |

Kita sudah melihat pola ini sepanjang tujuh sesi: **orang akan mencari jalan pintas, dan jalan
pintas itu akan menjadi tombol yang dimatikan `NEVER` bertahun-tahun kemudian.**

### Yang diminta

Keputusan pada **ADR-0055**, bukan di sini. Bentuk yang mungkin, seluruhnya milik pemilik proses:

| Pilihan | Yang perlu ditimbang |
|---|---|
| **✓ DIPILIH** — perpindahan `DRAFT → DIBATALKAN` (terminal baru) | menambah keadaan kedelapan; jejaknya utuh, ADR-0045 aman |
| tambah perpindahan `DRAFT → DITOLAK` langsung oleh pengisinya | tidak menambah keadaan, tetapi membuat `DITOLAK` berarti dua hal berbeda |
| longgarkan INV-25 | versi belum selesai boleh lebih dari satu — mengubah invarian, bukan menambah perpindahan |

**Sudah dipilih** — yang pertama, dengan tiga syarat mengikat. Alasan lengkapnya di ADR-0055
bagian "Perubahan 24 September 2026". Satu syaratnya, **nomor revisi tidak dipakai ulang**, ternyata
**tidak menuntut mekanisme baru**: INV-04 sudah membuat `NOMOR_URUT_VERSI` unik per kontrak, jadi
selama barisnya tidak dihapus, constraint yang sudah ada menolak pemakaian ulangnya sendiri.

---

## 2.2b P-59 — SYARAT BERBEDA TIAP PEMULIHAN LIMIT, bergolongan **BARU**

> **TABRAKAN NOMOR, dikoreksi 24 September 2026.** Kemampuan ini dan *"penolakan addendum
> meninggalkan catatan"* **sama-sama diberi nomor P-58**, keduanya pada hari yang sama, oleh saya.
> Dua kemampuan dengan satu nomor akan membuat setiap rujukan ke P-58 berarti dua hal — persis
> kegagalan yang kita larang pada nama kolom. **Yang di tabel utama tetap P-58; yang ini menjadi
> P-59**, dan seluruh rujukannya disunting.

Lahir 24 September 2026 dari sapuan tetapan di kode (`TETAPAN-DI-KODE.md` §3).

| Medan | Isi |
|---|---|
| **Kemampuan** | **PK** dapat menetapkan **syarat yang berbeda untuk tiap pemulihan limit** pada sebuah layer — pemulihan pertama, kedua, dan seterusnya boleh berpersentase berbeda |
| **Golongan** | **BARU** |
| **Kenapa BARU, bukan PELESTARIAN** | `SetReinstatementPct.xml` menulis `ReinstatementPct = "100"` dan `AdditionalPct = "100"` sebagai **tetapan**, bukan sebagai masukan. Sistem lama **tidak dapat menyatakan** ketentuan yang berbeda antar pemulihan. Yang PELESTARIAN hanyalah **menyimpan pemulihan sebagai daftar**; yang BARU adalah **nilainya boleh berbeda** |
| **Entitas** | `PEMULIHAN_LIMIT` — `SPEC-MODEL-DATA.md` §10.3b |
| **Invarian** | INV-49 (pemulihan adalah besaran berulang, dikecualikan dari partisi) |
| **Tidak terhalang luar** | ya |

**`CARA MENYALAKANNYA` — empat butir, wajib pada golongan BARU:**

| Butir | Isi |
|---|---|
| rujukan CO | menunggu penomoran CO saat tiketnya dibentuk |
| **siapa membaca, seberapa sering** | **PK yang menyusun layer non-proporsional**, pada setiap kontrak berpemulihan — bukan pemantau berkala |
| **ambang berangka** | nyala penuh ketika **satu kontrak** tercatat dengan dua pemulihan berpersentase berbeda dan lolos persetujuan. Sebelum itu ia tersedia tetapi **bernilai bawaan seragam**, sehingga tidak mengubah apa pun |
| **siapa boleh menyalakan** | pemilik proses, setelah bisnis mengonfirmasi bahwa pemulihan bertingkat memang dipakai NuRe |

> **Satu pertanyaan wawancara menyertainya, dan ia tidak menahan pembangunannya:** apakah
> keseragaman 100/100 di sistem lama **kehendak bisnis** atau **akibat cara sistem dibangun**?
> Datanya **tidak dapat menjawab** — seluruhnya disemai 100, sehingga *"selalu seragam"* dan
> *"tidak pernah bisa berbeda"* menghasilkan data yang persis sama.

---

## 2.4 DUA KEMAMPUAN YANG DIBELAH ANTAR BATCH — 24 September 2026

Ronde to-ticket membagi kemampuan Treaty In menjadi **batch 1** — yang **tidak** menyentuh daftar
keadaan — dan **batch 2**, yang menyentuhnya dan karena itu menunggu tanggapan **`REV-3`** atas
ADR-0055 §4.

Dua kemampuan jatuh **persis di garis itu**, dan `G3` memerintahkan **dipecah tepat di tempat
kepastian berhenti**. **Tanpa catatan ini, pembaca berikutnya membaca keduanya sebagai sudah
ditiketkan seluruhnya.**

### `P-01` — identitas di batch 1, keadaan lahirnya di batch 2

| Bagian | Batch | Tiket | Isinya |
|---|---|---|---|
| **identitas** | **1** | **`14`** | `KONTRAK` dan `VERSI_KONTRAK` berdiri beserta kunci alami, kunci asing, dan pengenalnya. Kolom `KEADAAN_SIKLUS_HIDUP` **ada** — tabelnya tidak dapat berdiri tanpanya |
| **keadaan lahirnya** | **2** | belum ditiketkan | *"versi pertamanya lahir langsung dalam keadaan `DRAFT`"*, beserta `INV-20`, `INV-22`, `INV-23`, `INV-24`, `INV-25` — **daftar nilai sah dan mesin perpindahannya** |

> **Garis pecahnya berprinsip:** yang dibangun batch 1 adalah **tempat**; yang ditunda **daftar nilai
> dan mesin perpindahannya** — dan `REV-3` merevisi persis yang kedua. Bila `REV-3` ditolak
> seluruhnya, tiket `14` **tidak berubah satu kriteria pun**.

### `P-50` — hanya pecahan tabel acuannya di batch 1

| Pecahan | Batch | Tiket | Sebab |
|---|---|---|---|
| **isi enam tabel acuan** | **1** | **`44`** | **nol keadaan ditulis**; dan setiap kunci asing di seluruh skema menunggunya |
| kepala kontrak | 2 | belum ditiketkan | **menulis `KEADAAN_SIKLUS_HIDUP`**, termasuk `WARISAN_TAK_TERPETAKAN` (`P-51`) |
| daftar anak versi · cabang non-prop · cabang prop · potongan dan penyebaran · catatan dan jejak | 2 | belum ditiketkan | seluruhnya menggantung pada pemindahan kepala |
| **addendum warisan** — nomor urut menurut kronologi | — | `05`, `10`, `12` | **sudah ditiketkan** ronde Adjustment sebagai `P-65`, pecahan payung ini |

---

## 2.3 GOLONGAN P-18 BELUM PASTI — akibat yang ditarik 24 September 2026

`SPEC-INVARIAN.md` §4.4b menyatakan: **INV-50 mulai bekerja pada hari sistem baru menyala, bukan
sebelumnya.** Sebabnya persentase rincian penyebaran di sistem lama **disalin dari tabel master**
(`SharePct = pxResults(idx).Pct`), dan master itu berjumlah 100 — sehingga invariannya tidak pernah
dapat dilanggar oleh data yang dihasilkan jalur itu.

Satu langkah lagi ditarik, dan ia mengubah sebuah golongan:

> Bila **ADR-0036 yang mengizinkan** penyebaran disunting, maka **sebelum ADR-0036 ia tidak
> disunting — ia dihitung.** Dan menyunting penyebaran adalah hal yang orang **belum pernah
> lakukan**.

Maka kalimat P-18 — *"PK dapat menyebarkan bagian NuRe…"* — mungkin bukan pemindahan kemampuan yang
sudah ada, melainkan **kemampuan BARU**. Yang memutuskannya **uji yang sama** yang menahan golongan
`RINCIAN_PENYEBARAN`:

| Hasil uji | Golongan P-18 | Akibat lanjutannya |
|---|---|---|
| persentase **pernah** menyimpang dari master | **L — PELESTARIAN** | tidak ada tambahan medan |
| persentase **tidak pernah** menyimpang | **B — BARU** | P-18 **menuntut medan `CARA MENYALAKANNYA`** seperti keenam BARU yang lain |

**Itu bukan tambahan pekerjaan; itu tambahan yang ketahuan sekarang alih-alih saat tiket ditulis.**

Dan kalimat kemampuannya sendiri sudah diperbaiki di satu hal yang tidak menunggu uji apa pun:
**"per pihak" menjadi "per jenis reasuransi"**, kata yang sama yang kini ada di INV-50. Tidak ada
sumbu pihak di keluarga penyebaran — `SPEC-INVARIAN.md` §4.4a.

> **Penghalang yang hanya separuh tertulis akan dibaca sebagai penghalang yang hanya separuh
> berat.** Karena itu kedua penghalang P-18 ditulis sebagai dua baris di §4.2, bukan satu.

---

## 3. Yang TIDAK ada di daftar ini, dan sebabnya

Supaya pembaca berikutnya tidak mengira daftar ini lupa.

| Tidak ada | Sebab |
|---|---|
| `RETRO_KELUAR` | **GEL-2** |
| `PENCAPAIAN` | **GEL-3** |
| `NILAI_SELISIH` dan perilaku penyesuaian | milik modul Adjustment; yang dibangun di sini hanya bentuk yang membuatnya mungkin (P-44) |
| jalur penerbitan ke luar pengganti `TREATYINOFFER` | penggantinya **belum diketahui** — eskalasi butir 3 |
| perbaikan dan hitung ulang data lama | ADR-0043: peristiwa bisnis tersendiri |
| uji data A…AB | bukan pekerjaan pengembang |
| tombol pengembang sistem lama | sudah diputuskan tidak dibawa |

---

## 4. Hitungan

**Dihitung ulang 24 September 2026** di gerbang pembuka to-ticket; angka lamanya disebut di kotak
§2 supaya pemeriksaan berikutnya tidak mencocokkan ke angka yang salah.

| | |
|---|---:|
| Kemampuan **didaftar** | ~~58~~ **66** — 65 berbaris + `P-59` yang tidak berbaris |
| — **DICORET** sebagai kemampuan (SIFAT, U-6) | 3 — P-24, P-25, P-26 |
| Kemampuan **aktif** | ~~55~~ **63** |
| **TEGAS TIDAK MASUK penyerahan pertama** (`G2`) | **3** — P-20, P-41, P-49 |
| Kemampuan **di dalam penyerahan pertama** | **60** |
| **TERTAHAN** | ~~4~~ **4** — P-18, P-60, P-61, P-62 |

> **Golongan PELESTARIAN / PERUBAHAN / BARU tidak dihitung ulang di sini**, dan itu dinyatakan:
> ketiganya menuntut membaca tiap baris, dan sapuan kata **tidak dapat membacanya**. Angka lamanya
> — 23 / 25 / 7 — **dicabut**, bukan diperbarui: angka yang diketahui salah lebih berbahaya
> daripada angka yang tidak ada.

**Empat yang tertahan, beserta penghalangnya** — `P-20`, `P-41`, dan `P-49` **pindah ke tabel
berikutnya**, sebab mereka bukan tertahan melainkan **tidak masuk**:

| Kode | Menunggu | Siapa yang dapat menjawab |
|---|---|---|
| P-18 | **(1)** INV-50 dan INV-51 belum terbukti — uji negatif *materialized view* belum pernah dijalankan | **kantor** — sediakan instans Oracle yang terjangkau (L-3, lubang lingkungan) |
| P-18 | **(2)** golongannya `L` atau `B` belum pasti — bergantung apakah penyebaran pernah disunting orang di sistem lama | **kantor** — izin kueri baca-saja ke produksi; uji yang sama yang menahan golongan `RINCIAN_PENYEBARAN` |
| P-60 | `DB-20` — titik beku materialitas | bisnis; pertanyaannya **belum dikirim**. Tiket `02` |
| P-61 | `DB-16a` — apakah dokumen addendum dikirim ke luar | bisnis; **belum dikirim**. Tiket `04` |
| P-62 | `DB-16b` — apakah dokumen punya tanggal berlaku sendiri | bisnis; **belum dikirim**. Tiket `08`, dan ia **bertenggat** |

### 4.1 TEGAS TIDAK MASUK penyerahan pertama — tiga, dan penandanya **diperbaiki 24 September 2026**

> **`TERTAHAN` dan "tegas tidak masuk" adalah SUMBU YANG BERBEDA**, dan ketiga kemampuan di bawah
> sempat ditandai dengan yang salah:
>
> | Penanda | Artinya |
> |---|---|
> | **`TERTAHAN`** | akan dikerjakan **begitu penghalangnya jatuh** |
> | **tegas tidak masuk** | **tidak dikerjakan** di penyerahan ini, **apa pun jawabannya** |
>
> Siapa pun yang mencabut penghalangnya akan mengira ketiganya jadi dapat dikerjakan — dan
> mengerjakan sesuatu yang `G2` sudah keluarkan dari ruang lingkup.

| Kode | Tepi `G2` | Sebab |
|---|---|---|
| **P-20** | *"invarian atas `AccountingMode` — menunggu Uji X-2"* | separuh keduanya **adalah** invariannya. Memotongnya menjadi *"kolomnya ada, aturannya menyusul"* gagal **U-3** |
| **P-41** | *"apa pun yang mewujudkan butir eskalasi yang belum diputuskan"* — **butir 5** | dokumen batas wewenangnya belum ada; `G3` sudah mencoba memecahnya dan mencatat kegagalannya |
| **P-49** | idem — **butir 7** | kolom Entitas-nya berbunyi *(belum ditetapkan)* |

**Ketiganya tetap berdiri sebagai kemampuan** dan **tidak dihapus**. Yang berubah hanya
kedudukannya: mereka menunggu **ruang lingkup berikutnya**, bukan menunggu jawaban.

---

## 5. KELUARAN LANGKAH 7 — entitas yang tersentuh penyerahan pertama

Ini yang menentukan seberapa besar sisa pekerjaan `SPEC-MODEL-DATA.md` §10.

### 5.0 URUTAN dan KEWAJIBAN — dibaca sebelum §10 dikerjakan

**Ditambahkan 24 September 2026 oleh pemilik proses.** Dua hal: mana yang didahulukan, dan apa yang
wajib dinyatakan setiap entitas.

#### a. Tiga entitas DIDAHULUKAN — bukan karena penting, melainkan karena menahan yang sudah selesai

| Entitas | Perannya | Invarian yang menjumlahkan barisnya |
|---|---|---|
| `BAGIAN` | bagian NuRe atas sebuah layer — **induk dari seluruh penyebaran** | — (ia induknya) |
| `RINCIAN_PENYEBARAN` | rincian penyebaran per pihak | **INV-50** |
| `NILAI_PENYEBARAN` | nilai penyebaran per pihak per mata uang | **INV-47**, **INV-51** |

Ketiganya adalah entitas yang **paling banyak disebut invarian rekonsiliasi**, dan justru ketiganya
yang **tidak punya invarian keunikan sendiri** (`SPEC-INVARIAN.md` §4.1, daftar lima).

Dua akibat, dan yang kedua lebih halus:

**a.1 — Seam.** `SPEC-INVARIAN.md` §2.1 menyatakan bahwa kunci alami yang unik **di dalam
lingkupnya** adalah yang membuat pemadanan baris antar versi mungkin. Untuk ketiga entitas ini
pemadanan itu **belum punya dasar** — dan penyebaran justru salah satu hal yang paling sering
disesuaikan. **Nilai selisih atas penyebaran tidak dapat dihitung per baris** sampai ketiganya punya
kunci yang dijamin unik.

**a.2 — Rekonsiliasinya tetap bekerja, tetapi pesannya menyesatkan.** Tanpa keunikan, satu baris
penyebaran yang tersisip dua kali membuat jumlahnya menjadi dua kali lipat — dan INV-50 **menolaknya
dengan benar**. Yang salah bukan penolakannya melainkan **sebab yang dilaporkan**: orang akan
mengira persentasenya salah, padahal barisnya yang kembar. *Constraint* yang menolak karena alasan
yang tidak terbaca membuat orang mencari di tempat yang salah.

> **Maka ketiganya dikerjakan lebih dulu**, dan sebabnya ditulis di sini, bukan di catatan terpisah:
> tanpa keduanya, **dua hal yang sudah selesai tidak sepenuhnya berlaku**.

**Berdampingan dengan ketiganya, untuk pembaca yang sama:** ruas induk pada `KUNCI_PADANAN`
`POTONGAN`. `POTONGAN` berinduk dua (§14.1 `SPEC-MODEL-DATA.md`), sehingga kunci padanannya **harus
menyebut induk mana** — tanpa ruas itu, potongan pada kedua induk saling tertukar saat dipadankan.
Temuan itu bentuknya sama persis dengan ketiga di atas, dan ia **sudah ditutup** di
`SEAM-ADJUSTMENT.md` §3. Disebut di sini supaya pembaca tahu ia satu keluarga, bukan dua persoalan.

**Dua yang BUKAN didahulukan, meski sama-sama tanpa nomor invarian:** `PORTOFOLIO` dan
`DOKUMEN_KONTRAK` (Bentuk A, ditemukan sapuan mekanis susulan). Keduanya memang belum punya kunci
yang dijamin unik, tetapi **tidak ada invarian yang menjumlahkan barisnya** — sehingga akibat a.2
tidak berlaku pada keduanya, dan hanya a.1 yang berlaku. Mereka dikerjakan pada urutan biasa.

#### b. KEWAJIBAN yang berlaku untuk KESELURUHAN 23 entitas

> **SETIAP ENTITAS MENYATAKAN KUNCI ALAMINYA DAN LINGKUPNYA SECARA TERANG** — **Bentuk A** (unik di
> dalam versi) atau **Bentuk B** (unik di dalam induk langsung).
>
> Entitas yang **TIDAK** punya kunci alami **menyatakannya juga, beserta alasannya** — karena
> entitas tanpa kunci alami **tidak dapat dipadankan antar versi sama sekali**, dan itu
> **keputusan, bukan kekosongan**.


#### b2. SUMBU REKONSILIASI — ditambahkan 24 September 2026

> **Setiap entitas yang punya BARIS TURUNAN menyatakan `SUMBU_REKONSILIASI`-nya** — **nama kolom**
> yang dipakai mengelompokkan saat baris turunannya dijumlahkan — berdampingan dengan kunci alami
> dan lingkupnya.
>
> Entitas yang **tidak punya** baris turunan **menyatakan itu juga**.

Sebabnya sama bentuknya dengan kewajiban di atas, tetapi pembacanya berbeda: yang membaca kunci
alami adalah seam Adjustment, yang membaca sumbu rekonsiliasi adalah **sesi DDL saat menulis
`GROUP BY`**.

**INV-47 bergantung pada kewajiban ini.** Rumusan INV-47 yang tidak menyebut kolom punya dua cacat
yang keduanya praktis: ia **tidak dapat gagal** — untuk sekumpulan baris mana pun selalu ada sumbu
yang membedakannya — dan ia **tidak dapat dikompilasi**, sehingga yang menulis *materialized view*
memilih kolomnya sendiri, diam-diam. Sejak §4.4c, INV-47 berbunyi *"menurut `SUMBU_REKONSILIASI`
yang dinyatakan §10 entitas itu"*, dan **§10-lah yang memuat nama kolomnya**.

Contoh yang sudah pasti: `RINCIAN_PENYEBARAN` bersumbu `ID_JENIS_REASURANSI`, **bukan** pihak —
tidak ada sumbu pihak di keluarga penyebaran (`SPEC-INVARIAN.md` §4.4a).

**Yang dicari bukan kerapian.** Entitas tanpa kunci alami yang tidak menyatakan dirinya tanpa kunci
alami akan **terlihat persis sama** dengan entitas yang lupa diberi kunci — dan bedanya baru
ketahuan saat seam-nya tidak bekerja.

**Kenapa kewajiban ini baru muncul sekarang.** Kelima entitas tanpa nomor invarian ditemukan dengan
mencocokkan **dua berkas**: kolom *kunci alami* di `4-erd-dan-tabel-datar/STRUKTUR-DATA.md` terhadap
daftar invarian bernomor di `SPEC-INVARIAN.md` §2.1 — yang ada di satu berkas dan tidak ada di yang
lain. **Itu lolos tujuh sesi.** Bukan karena ada yang lalai, melainkan karena **tidak ada satu pun
langkah yang pernah mencocokkan kedua berkas itu ke dua arah.** Kewajiban di atas membuat
pencocokan itu tidak lagi bergantung pada ada-tidaknya orang yang terpikir menjalankannya.

Arah sebaliknya sudah disapu sekali jalan dan **bersih** — tidak ada invarian yang menyebut entitas
yang tidak ada di `STRUKTUR-DATA.md`. Catatannya di `SPEC-INVARIAN.md` §4.1.

---

### 5.1 Entitas milik Treaty In — ~~21~~ **29** *(dihitung ulang 24 September 2026)*

> **Delapan entitas bertambah sesudah daftar ini ditulis**, dan tidak satu pun ditambahkan ke
> tabelnya di bawah — tabel itu dipertahankan **apa adanya** sebagai rekaman keadaan waktu itu.
> Yang **mengikat** soal daftar entitas tetap `4-erd-dan-tabel-datar/STRUKTUR-DATA.md` (urutan
> wewenang butir 4), dan sejak `F-18` ditutup ia **cocok dua arah** dengan `2-to-spec/KAMUS-KOLOM.md`.
>
> | Yang bertambah | Dari |
> |---|---|
> | `PEMULIHAN_LIMIT` | §10.23c butir 1 |
> | `PERISTIWA_KONTRAK` | §10.20a |
> | `DOKUMEN_ADDENDUM` | `GRL-19`, diff `D-2`/`D-3` |
> | `NILAI_MDP`, `NILAI_MDP_MINIMUM`, `NILAI_PREMI_BRUTO`, `NILAI_PREMI_BRUTO_MINIMUM`, `NILAI_CADANGAN_PREMI` | `P-8` golongan A — paket uang yang asalnya **daftar per mata uang** menjadi tabel anak. Kelimanya **tidak pernah terdaftar** di berkas yang mengikat sampai `F-18` menutupnya |


| # | Entitas | Dipakai kemampuan | §10 sudah? |
|---:|---|---|---|
| 1 | `KONTRAK` | P-01, 02, 03, 05, 43, 54, 56 | **sudah** (§10.1) |
| 2 | `VERSI_KONTRAK` | P-01, 04, 20, 27, 28, 29…45, 51, 52, 54, 56 | **sudah** (§10.2) |
| 3 | `LAYER` | P-16 | **sudah** (§10.3) |
| 4 | `DETAIL_PROPORSIONAL` | P-21, 22, 23 | **sudah** (§10.4) |
| 5 | `MATA_UANG_KONTRAK` | P-06, 27 | **sudah** |
| 6 | `RETENSI_CEDANT` | P-07 | **sudah** |
| 7 | `EGNPI` | P-08 | **sudah** |
| 8 | `PORTOFOLIO` | P-09 | **sudah** — *kunci alaminya belum bernomor*, §5.0a |
| 9 | `PERIODE_PELAPORAN` | P-10 | **sudah** |
| 10 | `PERIODE_AKUMULASI` | P-11 | **sudah** |
| 11 | `TERMIN` | P-12 | **sudah** |
| 12 | `SKALA_KOASURANSI` | P-13 | **sudah** |
| 13 | `BATAS_PER_BAHAYA` | P-14 | **sudah** |
| 14 | `DOKUMEN_KONTRAK` | P-15 | **sudah** — *kunci alaminya belum bernomor*, §5.0a |
| 15 | `CATATAN_PERSETUJUAN` | P-31…36, 46 | **sudah** |
| 16 | `JEJAK_PERUBAHAN` | P-47 | **sudah** |
| 17 | `BAGIAN` | P-17 | **sudah** — dikerjakan lebih dulu, §5.0a |
| 18 | `POTONGAN` | P-19 | **sudah** |
| 19 | `PENYEBARAN` | P-18 | **sudah** |
| 20 | `RINCIAN_PENYEBARAN` | P-18 | **sudah** — dikerjakan lebih dulu, §5.0a |
| 21 | `NILAI_PENYEBARAN` | P-18 | **sudah** — dikerjakan lebih dulu, §5.0a |

### 5.2 Tabel acuan — 6

Seluruhnya dipakai P-48 dan dirujuk kunci asing dari entitas di atas. Bentuknya seragam dan kecil
— kode, nama, penanda aktif — sehingga §10 untuknya jauh lebih murah daripada entitas biasa.

`MATA_UANG` · `JENIS_POTONGAN` · `JENIS_REASURANSI` · `BAHAYA` · `KELOMPOK_TREATY` · `KELAS_BISNIS`

### 5.3 Angka yang ditunggu

| | |
|---|---:|
| Entitas tersentuh penyerahan pertama | ~~27~~ **35** |
| — milik Treaty In | ~~21~~ **29** |
| — tabel acuan | 6 |

> **Angka 35 dicocokkan dua arah** terhadap `2-to-spec/KAMUS-KOLOM.md` pada 24 September 2026, dan
> keduanya sepakat. Sebelum `F-18` ditutup, selisihnya **lima** — seluruhnya tabel anak paket uang
> yang berdiri di kamus dan DDL tetapi tidak pernah terdaftar di berkas yang mengikat.
| §10 **sudah** ada, 23 September | 4 |
| §10 **ditulis 24 September** | **23** — 17 entitas + 6 tabel acuan |
| §10 **diselesaikan ulang** sesudah L-8 | 2 — `LAYER`, `DETAIL_PROPORSIONAL` |
| §10 **belum** ada | **0** |

> ### L-2 DITUTUP 24 September 2026
>
> §10 kini meliputi seluruh 27, dan dua di antaranya dikerjakan **dua kali** — bukan karena ada yang
> berubah pikiran, melainkan karena keduanya semula **diselesaikan atas sumber yang kurang** (L-8).
>
> Empat keputusan tersisa dan ketiganya menambah entitas — `SPEC-MODEL-DATA.md` §10.23c. Bila
> diterima, **27 menjadi 30**, dan ketiga entitas barunya lahir dari titik buta yang sama.

> ### Angka ini LEBIH BESAR daripada dugaan, bukan lebih kecil.
>
> Dugaan sebelum daftar ini disusun: entitas yang terpakai akan "jauh di bawah 20", sehingga sisa
> pekerjaan §10 mengecil. **Yang terjadi sebaliknya.**
>
> Sebabnya dua, dan keduanya terbaca dari daftar:
>
> 1. **Hampir seluruh entitas Treaty In terpakai di penyerahan pertama.** Yang dikeluarkan G2 —
>    `RETRO_KELUAR`, `PENCAPAIAN`, `NILAI_SELISIH` — memang hanya **tiga**, dan ketiganya sudah di
>    luar daftar entitas inti sejak awal. Batas G2 memotong **modul dan gelombang**, bukan memotong
>    isi kontraknya.
> 2. **Enam tabel acuan sebelumnya tidak pernah ikut dihitung.** `SPEC-MODEL-DATA.md` §2.3 tidak
>    mendaftarnya, dan §11 hanya menyebut "entitas selebihnya".
>
> **Jalan (iii) tetap benar**, dan justru angka ini yang membuktikannya berguna: tanpa daftar ini,
> sesi to-spec akan mengerjakan §10 untuk **entitas yang salah** — mengerjakan `RETRO_KELUAR` dan
> `PENCAPAIAN` yang tidak dipakai, sambil melewatkan enam tabel acuan yang dipakai. Yang berubah
> bukan besarnya pekerjaan, melainkan **ketepatan sasarannya**.

### 5.4 Yang tetap kosong, dan tercatat kosong

| Entitas | Sebab | Kapan §10-nya dikerjakan |
|---|---|---|
| `RETRO_KELUAR` | GEL-2 | saat GEL-2 dimulai |
| `PENCAPAIAN` | GEL-3 | saat GEL-3 dimulai |
| `NILAI_SELISIH` | milik modul Adjustment | saat embargo Adjustment dicabut |
| `BESARAN_DAPAT_DISESUAIKAN` | tabel acuan sambungan; **tidak tersentuh satu pun kemampuan di atas** | saat sambungan Adjustment dibangun |

`BESARAN_DAPAT_DISESUAIKAN` sengaja **tidak** dimasukkan ke hitungan 27. Ia mudah dibangun dan
menggoda untuk disertakan "sekalian", tetapi tidak ada satu pun kemampuan di penyerahan pertama
yang memakainya — dan memasukkannya berarti mengerjakan §10 untuk sesuatu yang tidak dipakai,
yaitu persis kesalahan jalan (i).

---

## 6. Apa yang terjadi berikutnya

Menurut jalan (iii): **langkah c — §10 dikerjakan untuk 23 entitas di §5.3, dan hanya untuk itu**,
dengan **urutan §5.0a** (tiga didahulukan) dan **kewajiban §5.0b** (kunci alami dan lingkupnya
dinyatakan terang, termasuk ketika tidak ada) berlaku pada keseluruhannya.

Sesudahnya: sesi DDL, lalu pembentukan tiket (langkah 4).
