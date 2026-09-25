# SPEC-INVARIAN — Treaty In

**Tanggal:** 23 September 2026 · **direvisi 24 September 2026**
**Revisi 24 Sep 2026:** §0 aturan **kalimat invarian harus benar sendiri** · kalimat INV-20, INV-22,
INV-23, INV-24 diperbaiki (bukan hanya penandanya) · **§4.3a trigger anak `BEFORE INSERT OR UPDATE
OR DELETE`** — INSERT sebelumnya tidak dijaga · §4.1 dipisah **Bentuk A (versi) / Bentuk B (induk)** ·
§1 empat yang belum terbukti dipisah · §4.3c `WARISAN_TAK_TERPETAKAN` memakai jatah versi tak-terminal
**Revisi keempat 24 Sep 2026:** **INV-12 kunci alaminya kurang satu ruas** — termin disusun per
mata uang (§4.1) · invarian ketiga yang lingkupnya keliru
**Revisi ketiga 24 Sep 2026:** **INV-47 dan INV-50 menyebut sumbu yang tidak ada** — "pihak"
diperbaiki menjadi jenis reasuransi (§4.4a, calon penyangkal `ReisuredParticipant` dibuka dan gugur) ·
§4.4b kenapa INV-50 tidak pernah gagal di sistem lama · **§4.4c INV-47 menunjuk `SUMBU_REKONSILIASI`
di §10**, karena rumusan sebelumnya tidak dapat gagal dan tidak dapat dikompilasi
**Revisi kedua 24 Sep 2026:** §4.1 kunci alami tanpa nomor invarian **tiga → LIMA** (`PORTOFOLIO` dan
`DOKUMEN_KONTRAK` ikut, keduanya Bentuk A) · sapuan arah sebaliknya dijalankan dan **bersih**
**Keadaan:** hasil langkah 8. Ditulis sebagai **data**, bukan uraian.
**Sumber:** `SPEC-MODEL-DATA.md`, ADR-0034…0055, ADR warisan 0003–0033, `SEAM-ADJUSTMENT.md`

---

## 0. Cara membacanya

Setiap invarian punya nomor tetap `INV-nn`, satu kalimat isinya, dan **tempat penegakannya**. Empat
tempat, dan urutannya bukan selera — ia urutan yang dicoba:

| Tempat | Kapan dipakai |
|---|---|
| **CONSTRAINT** | dapat dinyatakan deklaratif: `CHECK`, `UNIQUE`, `FOREIGN KEY`, `NOT NULL` |
| **INDEKS UNIK** | keunikan **bersyarat** yang tidak dapat dinyatakan `UNIQUE` biasa |
| **TRIGGER** | tidak ada cara deklaratif **dan** akibat pelanggarannya menyentuh uang |
| **APLIKASI** | selebihnya — dan **wajib** menyebut alasan kenapa tidak di basis data |

**Tidak ada invarian yang boleh hilang dari daftar ini.** Bila sebuah invarian tidak ditegakkan di
mana pun, itu keputusan, dan keputusan ditulis.

Kolom **§** merujuk sumbernya di `SPEC-MODEL-DATA.md` atau nomor ADR.

### KALIMAT INVARIAN HARUS BENAR SENDIRI

> **Penanda ⚠ adalah catatan tentang RIWAYAT perubahan, bukan KOREKSI atas kalimat yang salah.**
> Bila kalimatnya keliru, yang diperbaiki **kalimatnya**; penandanya menerangkan **sejak kapan**.

Alasannya praktis, bukan kerapian: sesi DDL menulis `CHECK` dari kalimat INV-20 dan trigger dari
kalimat INV-23. **Yang dibaca orang saat menulis kode adalah kalimatnya.** Penanda ⚠ dibaca orang
yang sedang menelusuri sejarah, dan itu bukan orang yang sama.

**Uji yang dijalankan pada setiap baris ber-⚠:**
*bila penandanya dihapus, apakah kalimatnya masih benar?* Yang menjawab tidak, **diperbaiki**.

> **Ini pernah gagal, di berkas ini, pada 24 September 2026.** INV-20 masih berbunyi "tujuh"
> dan INV-23 masih menyebut dua keadaan, sementara penandanya bilang delapan dan tiga. Lebih buruk:
> §4.3 mengutip INV-23 dengan isi **tiga keadaan** — mengutip invarian dengan isi yang berbeda dari
> kalimatnya, di berkas yang sama, sembilan puluh baris kemudian.
>
> Itu persis keluarga kegagalan yang diperingatkan §4.3 sendiri: **dua tempat, satu benar satu
> basi, dan yang basi terbaca sebagai berlaku.**

---

## 1. Hitungan

| Tempat penegakan | Jumlah | Bagian |
|---|---:|---:|
| **CONSTRAINT — terbukti** | 38 | 60,3 % |
| **CONSTRAINT — BELUM TERBUKTI** | **4** | **6,3 %** |
| **INDEKS UNIK** | 4 | 6,3 % |
| **TRIGGER** | 6 | 9,5 % |
| **APLIKASI** | 11 | 17,5 % |
| **Total** | **63** | |

> **Keempat yang belum terbukti diberi barisnya sendiri, tidak dilebur ke dalam CONSTRAINT.**
> Alasannya sama persis dengan alasan yang dipakai di §5.2 ke arah sebaliknya: di sana enam invarian
> yang sebenarnya tentang **rancangan** tetap dihitung sebagai APLIKASI **supaya angkanya tidak
> diperkecil secara semu**. Menghitung empat klaim di dalam CONSTRAINT akan **membesarkannya** dengan
> cara yang sama semunya.
>
> Sampai uji negatifnya lulus, keempatnya **klaim, bukan fakta**.

**Bila keempatnya gagal**, keempat invarian itu **turun ke APLIKASI**, dan angkanya menjadi:

| | |
|---|---:|
| CONSTRAINT | 38 — **60,3 %** |
| APLIKASI | **15 dari 63** — **23,8 %** |

Angka 23,8 % itu bukan angka baru: ia persis bagian APLIKASI **sebelum** dua keluarga invarian
ditarik kembali ke basis data (§5.1). **Yang dipertaruhkan pada uji yang belum dijalankan adalah
seluruh hasil penarikan itu.** Kedua angka ditulis berdampingan supaya pembaca tahu besarnya tanpa
harus menghitung sendiri.

Ditambah **2 invarian yang TIDAK DAPAT DILANGGAR** — §2.8.

**Empat constraint berstatus BELUM DIBUKTIKAN**, dan itu bukan formalitas. INV-47, INV-50, INV-51
dan bentuk lintas baris INV-31 ditegakkan lewat *materialized view* ber-`REFRESH ON COMMIT` dengan
`CHECK` padanya. Teknik itu **belum dijalankan di Oracle**, dan sampai uji negatifnya lulus ia
**klaim, bukan fakta**. Rancangan ujinya di `UJI-NEGATIF-INVARIAN.md`. Bila salah satu gagal,
invarian itu **turun ke aplikasi dan tabel ini berubah sebelum dipakai sesi DDL**.

Bagian **APLIKASI 17,5 %** berada di bawah ambang sepertiga. Dua keluarga yang biasanya jatuh ke
aplikasi ditarik kembali ke basis data lebih dulu (§5.1); tanpa penarikan itu bagiannya 23,8 %.
**Angka itu hanya berarti bila kesebelasnya memang yang tidak mungkin** — §5.2 menyebut kesebelasnya
satu per satu dengan alasannya, karena daftar dapat diperiksa orang lain dan persentase tidak.

---

## 2. Daftar invarian

### 2.1 Identitas dan kunci — 19

| # | Invarian | Tempat | § |
|---|---|---|---|
| INV-01 | Setiap tabel punya kunci utama. Tanpa kecuali. | CONSTRAINT | larangan |
| INV-02 | Pengenal dibangkitkan `SEQUENCE`, tidak pernah dari cap waktu maupun teks yang dapat disunting. | CONSTRAINT | 0040, seam §1 |
| INV-03 | Pengenal tidak pernah dipakai ulang; sequence tanpa `CYCLE`. | CONSTRAINT | seam §1 |
| INV-04 | `NOMOR_URUT_VERSI` unik di dalam satu `KONTRAK`. | CONSTRAINT | 10.2 |
| INV-05 | `NOMOR_LAYER` + `BAGIAN_LAYER` unik di dalam satu versi. | CONSTRAINT | 10.3 |
| INV-06 | `ID_KELOMPOK_TREATY` unik di dalam satu `LAYER`. | CONSTRAINT | 10.4 |
| INV-07 | Kode mata uang unik di dalam satu versi (`MATA_UANG_KONTRAK`). | CONSTRAINT | 3.5 |
| INV-08 | Kelompok treaty + mata uang unik di dalam satu versi (`RETENSI_CEDANT`). | CONSTRAINT | 3.5 |
| INV-09 | Kelompok treaty + mata uang unik di dalam satu versi (`EGNPI`). | CONSTRAINT | 3.5 |
| INV-10 | Periode unik di dalam satu versi (`PERIODE_PELAPORAN`). | CONSTRAINT | 3.5 |
| INV-11 | Periode unik di dalam satu versi (`PERIODE_AKUMULASI`). | CONSTRAINT | 3.5 |
| INV-12 | **Kode mata uang + nomor termin** unik di dalam satu versi (`TERMIN`). ⚠ *semula hanya "nomor termin"; diperbaiki 24 Sep 2026 — §4.1* | CONSTRAINT | 12.4, 10.16a |
| INV-13 | Persen limit unik di dalam satu versi (`SKALA_KOASURANSI`). | CONSTRAINT | 12.4 |
| INV-14 | Bahaya unik di dalam satu versi (`BATAS_PER_BAHAYA`). | CONSTRAINT | 10.0b |
| INV-15 | Jenis potongan unik di dalam satu induk (`POTONGAN`). | CONSTRAINT | 14.2 |
| INV-16 | Jenis reasuransi unik di dalam satu induk (`PENYEBARAN`). | CONSTRAINT | 12.4 |
| INV-17 | Setiap kunci asing menunjuk **satu** sasaran tetap; tidak ada rujukan yang sasarannya bergantung nilai kolom lain. | CONSTRAINT | larangan |
| INV-18 | Perilaku hapus setiap kunci asing ditetapkan sadar, tidak dibiarkan bawaan. | CONSTRAINT | urutan kerja 5 |
| INV-19 | **Kunci alami kontrak tidak berubah sepanjang hidup kontrak.** | **TRIGGER** | permintaan (b) |

> **INV-05 sampai INV-16 adalah yang membuat seam Adjustment berfungsi.** Tanpa kunci alami yang
> unik **di dalam lingkupnya**, pemadanan baris antar versi tidak mungkin dan selisih hanya dapat
> dihitung atas total.
>
> **Lingkupnya tidak seragam** — sembilan berlingkup **versi**, tiga berlingkup **induk langsung**
> (INV-06, INV-15, INV-16). Lihat §4.1. Akibatnya pada pemadanan: untuk yang berlingkup induk,
> `KUNCI_PADANAN` harus **majemuk** — kunci padanan induknya lebih dulu, baru kuncinya sendiri.
> Lihat `SEAM-ADJUSTMENT.md` §3, `KUNCI_PADANAN`.

### 2.2 Siklus hidup — 9

| # | Invarian | Tempat | § |
|---|---|---|---|
| INV-20 | `KEADAAN_SIKLUS_HIDUP` hanya bernilai salah satu dari **delapan** — tujuh keadaan sah ditambah `WARISAN_TAK_TERPETAKAN`. ⚠ *`DIBATALKAN` ditambahkan ADR-0055 perubahan 24 Sep 2026; sebelumnya tujuh* | CONSTRAINT | 0055 |
| INV-21 | `KEADAAN_WARISAN_ASLI` terisi **jika dan hanya jika** keadaannya `WARISAN_TAK_TERPETAKAN`. | CONSTRAINT | 0054 |
| INV-22 | Perpindahan keadaan hanya yang ada di **daftar tiga belas** (§3). ⚠ *semula dua belas; `DRAFT` → `DIBATALKAN` ditambahkan ADR-0055 perubahan 24 Sep 2026* | **TRIGGER** | 0055 |
| INV-23 | Tidak ada perpindahan keluar dari `DISETUJUI`, `DITOLAK`, maupun `DIBATALKAN`. ⚠ *`DIBATALKAN` ditambahkan ADR-0055 perubahan 24 Sep 2026* | **TRIGGER** | 0055 |
| INV-24 | **Versi berkeadaan TERMINAL — `DISETUJUI`, `DITOLAK`, `DIBATALKAN` — tidak pernah berubah nilainya, dan tidak pernah bertambah baris anaknya.** ⚠ *semula hanya `DISETUJUI`; diperluas 24 Sep 2026* | **TRIGGER** | permintaan (c), 0036, §4.3 |
| INV-25 | Paling banyak **satu** versi per kontrak berada di keadaan tak-terminal. ⚠ *`DIBATALKAN` terhitung terminal — ADR-0055 perubahan 24 Sep 2026; lihat §4.3c untuk `WARISAN_TAK_TERPETAKAN`* | **INDEKS UNIK** | 0040 |
| INV-26 | `WARISAN_TAK_TERPETAKAN` hanya dapat dimasuki lewat migrasi, tidak pernah oleh sistem berjalan. | APLIKASI | 0054 |
| INV-27 | Kelengkapan per perpindahan terpenuhi sebelum perpindahan terjadi (lihat §3). | APLIKASI | langkah 8 |
| INV-28 | Setiap perpindahan meninggalkan satu baris `CATATAN_PERSETUJUAN` dengan pelaku dan waktunya. | **TRIGGER** | 0045 |

> **INV-28 dijaga bentuknya, 24 Sep 2026.** Sistem lama menyimpan **peristiwa** dan **keputusan** di satu daftar (`CommentList`). Bila keduanya dibawa ke satu tabel, INV-28 menjadi **bersyarat** — triggernya harus tahu baris mana yang dihitung — dan setiap hitungan di atasnya tercemar. Baris peristiwa karena itu **keluar** ke `PERISTIWA_KONTRAK` (`SPEC-MODEL-DATA.md` §10.21a), dan INV-28 tetap sesederhana bunyinya.
>
> **Satu tabrakan BELUM tertutup dan dilaporkan apa adanya:** baris **pengajuan** di sistem lama juga membawa `IsApproved = "Accept"`, sehingga migrasi yang menuruti aturan mekanis akan menaruhnya di `CATATAN_PERSETUJUAN` dan membuat INV-28 **terpenuhi kelebihan satu**. Bentuk barisnya identik dengan baris persetujuan. Diukur oleh **Uji AK**; uraiannya `4-erd-dan-tabel-datar/PERIKSA-BUTIR-CONTOH.md` §7.10-Z.5.

### 2.3 Cabang dan pembeda — 7

| # | Invarian | Tempat | § |
|---|---|---|---|
| INV-29 | `SIFAT_PROPORSI` hanya bernilai `PROPORSIONAL` atau `NON_PROPORSIONAL`. | CONSTRAINT | 10.1 |
| INV-30 | `JENIS_TREATY` hanya bernilai `QUOTA_SHARE` atau `SURPLUS`. | CONSTRAINT | 10.4 |
| INV-31 | **Tepat satu** dari `PERSEN_QUOTA_SHARE` dan `JUMLAH_LINES_SURPLUS` terisi, ditentukan `JENIS_TREATY`. | CONSTRAINT · bentuk lintas baris **belum dibuktikan** | 10.4 |
| INV-32 | `DETAIL_PROPORSIONAL` hanya ada di bawah kontrak berSIFAT_PROPORSI `PROPORSIONAL`. | **INDEKS UNIK** | 12.3 |
| INV-33 | `BAGIAN` hanya ada di bawah kontrak berSIFAT_PROPORSI `NON_PROPORSIONAL`. | **INDEKS UNIK** | 12.3 |
| INV-34 | `CARA_PEMBUKUAN_XOL` hanya terisi pada kontrak non-proporsional. | **INDEKS UNIK** | 14.4 |
| INV-35 | Baris `SURPLUS` menuntut adanya baris `QUOTA_SHARE` pada **versi yang sama**; ketiadaannya kegagalan yang dilaporkan, bukan nol. | APLIKASI | 3.4, 0035 |

### 2.4 Paket uang — 11

| # | Invarian | Tempat | § |
|---|---|---|---|
| INV-36 | Nilai uang terisi ⟹ mata uangnya terisi. | CONSTRAINT | 0007 |
| INV-37 | `NILAI_IDR` terisi ⟹ kurs, tanggal kurs, dan sumber kurs terisi. | CONSTRAINT | 0007, 0029 |
| INV-38 | Kurs selalu lebih besar dari nol; tidak pernah nol, tidak pernah dipaksa satu. | CONSTRAINT | 0035 |
| INV-39 | `TINGKAT_PENCATATAN` hanya bernilai `TREATY_100_PERSEN` atau `BAGIAN_NURE`. | CONSTRAINT | 0039 |
| INV-40 | Besaran bertingkat `BAGIAN_NURE` ⟹ `PERSEN_BAGIAN_DIPAKAI` terisi. | CONSTRAINT | 0039 |
| INV-41 | Persentase berada di rentang 0–100. | CONSTRAINT | — |
| INV-42 | Nilai uang tidak pernah memuat angka penanda kegagalan; ketidakmampuan menghitung menghasilkan **kosong**, bukan angka. | CONSTRAINT | 0035 |
| INV-43 | Kurs yang dipakai dicatat **sebagai nilai pada versi**, tidak dibaca ulang saat dibutuhkan. | CONSTRAINT | 0036, seam §2 |
| INV-44 | Kode mata uang merujuk tabel acuan mata uang. | CONSTRAINT | 0053 |
| INV-45 | Presisi seragam untuk besaran sejenis di seluruh skema. | CONSTRAINT | 0003 |
| INV-46 | Tidak ada pasangan kolom kembar untuk dua mata uang. | CONSTRAINT | larangan |

### 2.5 Rekonsiliasi dan agregat — 6

| # | Invarian | Tempat | § |
|---|---|---|---|
| INV-47 | **Jumlah seluruh baris turunan atas suatu besaran, dijumlahkan menurut `SUMBU_REKONSILIASI` yang dinyatakan §10 entitas itu di `SPEC-MODEL-DATA.md`, sama dengan besaran induknya.** Baris turunan tidak pernah dijumlahkan bersama induknya. ⚠ *semula "sumbu pihak", lalu sempat "sumbu yang membedakan baris turunannya" — keduanya diperbaiki 24 Sep 2026, §4.4c* | **CONSTRAINT — BELUM DIBUKTIKAN** (§5.1) | langkah 4 |
| INV-48 | Pengecualian bernama INV-47 (1): **kelipatan** — `KAPASITAS_SURPLUS` = retensi × jumlah lines. Bukan partisi. | CONSTRAINT | 3.4 |
| INV-49 | Pengecualian bernama INV-47 (2): **besaran berulang** — reinstatement boleh terjadi lebih dari sekali. Bukan partisi. | CONSTRAINT | 3.3 |
| INV-50 | Jumlah persen penyebaran **per induk penyebaran, dikelompokkan menurut jenis reasuransi**, sama dengan 100. ⚠ *semula tertulis "per sumbu pihak"; diperbaiki 24 Sep 2026 — §4.4a* | **CONSTRAINT — BELUM DIBUKTIKAN** (§5.1) | 0050 |
| INV-51 | Retensi + penyerahan sama dengan 100 persen. | **CONSTRAINT — BELUM DIBUKTIKAN** (§5.1) | 4.3 |
| INV-52 | Premi bruto dikurangi seluruh potongan sama dengan premi bersih. | APLIKASI | 4.3, 14.1 |

### 2.6 Tanggal dan periode — 5

| # | Invarian | Tempat | § |
|---|---|---|---|
| INV-53 | `TANGGAL_MULAI` ≤ `TANGGAL_BERAKHIR`; keduanya batas **inklusif**. | CONSTRAINT | 0022 |
| INV-54 | `TANGGAL_BERLAKU_ADDENDUM` berada di dalam periode kontraknya. | **TRIGGER** | 0040 |
| INV-55 | Periode pelaporan berada di dalam periode kontraknya. | APLIKASI | 3.5 |
| INV-56 | Periode akumulasi berada di dalam periode kontraknya. | APLIKASI | 3.5 |
| INV-57 | Tanggal kurs tidak lebih akhir dari tanggal transaksi yang memakainya. | CONSTRAINT | 0029 |

### 2.7 Turunan, salinan, dan bentuk — 5

| # | Invarian | Tempat | § |
|---|---|---|---|
| INV-58 | Tidak ada kolom turunan yang disimpan, kecuali ia fakta terbukukan yang membawa penunjuk asalnya. | APLIKASI | 0037 |
| INV-59 | Tidak ada nilai yang disalin dari entitas lain; hilir diberi rujukan. | APLIKASI | 0023, 0041 |
| INV-60 | Satu fakta punya tepat satu penulis. | APLIKASI | 0041 |
| INV-61 | Dokumen JSON arsip tidak punya jalur baca; ia bukan sumber kanonik. | APLIKASI | 0034 |
| INV-62 | Himpunan yang dapat bertambah tanpa mengubah arti disimpan sebagai tabel acuan, bukan `CHECK`. | CONSTRAINT | 0038 |

### 2.8 TIDAK DAPAT DILANGGAR — 2

Golongan ini **bukan penegakan yang lebih lemah; ia lebih kuat daripada constraint**, karena
pelanggarannya tidak dapat dinyatakan sama sekali. Tidak ada yang perlu memeriksa, karena tidak ada
keadaan yang perlu ditolak.

Karena itu kata "ditegakkan" salah arah dan tidak dipakai: ia mengundang pertanyaan "oleh apa", dan
pertanyaan itu tidak punya jawaban. Pertanyaan yang benar adalah **keadaan seperti apa yang harus
ada agar pelanggarannya dapat dinyatakan** — dan kenapa keadaan itu tidak ada di model ini.

| # | Invarian | Keadaan yang diperlukan untuk melanggarnya, dan kenapa ia tidak ada |
|---|---|---|
| INV-B1 | Lapisan beku addendum — cedant, asal bisnis, sifat proporsi, tanggal mulai, dan tanggal berakhir tidak berubah antar versi. | Melanggarnya menuntut **dua versi dari satu kontrak yang memuat nilai berbeda** untuk salah satu dari kelimanya. Kelimanya hanya ada di `KONTRAK`, satu baris per kontrak; tidak ada tempat kedua yang dapat memuat nilai yang berbeda. *(Mengubah nilainya di `KONTRAK` bukan pelanggaran ini melainkan INV-19, dan itu ditegakkan trigger.)* |
| INV-B2 | `OLDDATA` tidak pernah disalin. | Melanggarnya menuntut **sebuah kolom yang memuat nilai milik versi lain**. Seluruh nilai ada di versinya sendiri, dan "nilai lama" adalah versi sebelumnya — sebuah SELECT. Tidak ada kolom salinan untuk diisi. |

### 2.8b Satu yang TURUN dari golongan ini setelah diuji

INV-B3 semula ditulis di sini: *"aturan potongan tidak tertulis dua kali"*. Ia **tidak lulus** uji
di atas.

Keadaan yang melanggarnya **dapat disusun**: sesi DDL membuat dua tabel potongan, masing-masing
dengan `CHECK`-nya sendiri, dan aturannya tertulis dua kali. Keadaan itu tidak mustahil — ia bahkan
bentuk yang paling mudah dipilih orang yang belum membaca §14.1 `SPEC-MODEL-DATA.md`.

Maka ia bukan dijamin bentuk. Ia turun menjadi **INV-63, APLIKASI** — diperiksa saat tinjauan skema,
bukan saat sebuah baris ditulis:

| # | Invarian | Tempat | § |
|---|---|---|---|
| INV-63 | Aturan potongan tertulis **sekali**, berlaku pada kedua pelekatannya. | APLIKASI (tinjauan skema) | 14.1 |
| INV-64 | Satu `BAGIAN` per `LAYER` — `ID_LAYER` unik di `BAGIAN`. | CONSTRAINT | 10.5 · langkah 3 to-spec |
| INV-65 | Jenis reasuransi unik di dalam satu `PENYEBARAN` (`RINCIAN_PENYEBARAN`). | CONSTRAINT | 10.7 · langkah 3 to-spec |
| INV-66 | Arah + jenis portofolio unik di dalam satu versi (`PORTOFOLIO`). | CONSTRAINT | 10.13 · langkah 3 to-spec |
| INV-67 | Dokumen unik di dalam satu versi (`DOKUMEN_KONTRAK`). | CONSTRAINT | 10.19 · langkah 3 to-spec |
| INV-68 | `KODE` unik di seluruh tabel pada keenam tabel acuan. | CONSTRAINT | 10.22 · langkah 3 to-spec |
| INV-69 | Versi ber-`SIFAT_MATERIAL_ADDENDUM` = `TIDAK_MATERIAL` **tidak punya** baris `NILAI_SELISIH` bertipe uang atau porsi. | **CONSTRAINT lewat MV** — **KLAIM** | GRL-20 · diff `D-4` |
| INV-70 | Versi ber-`SIFAT_MATERIAL_ADDENDUM` = `MATERIAL` **tidak mengubah** teks kesepakatan dan identitas kontrak terhadap versi dasarnya. | **CONSTRAINT lewat MV** — **KLAIM** | GRL-20 · diff `D-4` |
| INV-71 | `NOMOR_DOKUMEN` unik di seluruh `DOKUMEN_ADDENDUM`. | CONSTRAINT | GRL-19 · diff `D-4` |

---

## 3. Kelengkapan per perpindahan — sebagai data

Kolom berisi apa yang **harus terpenuhi** sebelum perpindahan itu boleh terjadi. Kosong berarti
tidak ada syarat tambahan di luar invarian yang selalu berlaku.

| Perpindahan | Syarat kelengkapan |
|---|---|
| `LAHIR` → `DRAFT` | kunci alami kontrak lengkap; `SIFAT_PROPORSI` terisi |
| `DRAFT` → `MENUNGGU_SEC_HEAD` | **daftar K1** di bawah terpenuhi seluruhnya |
| `MENUNGGU_SEC_HEAD` → `MENUNGGU_DEPT_HEAD` | pemberi keputusan bukan pengisi kontraknya |
| `MENUNGGU_SEC_HEAD` → `DRAFT` | alasan pengembalian terisi |
| `MENUNGGU_SEC_HEAD` → `DITOLAK` | alasan penolakan terisi |
| `MENUNGGU_DEPT_HEAD` → `MENUNGGU_DIREKTUR` | pemberi keputusan berbeda dari tingkat sebelumnya |
| `MENUNGGU_DEPT_HEAD` → `DRAFT` | alasan pengembalian terisi |
| `MENUNGGU_DEPT_HEAD` → `DITOLAK` | alasan penolakan terisi |
| `MENUNGGU_DIREKTUR` → `DISETUJUI` | **daftar K2** terpenuhi; pemberi keputusan berbeda dari kedua tingkat sebelumnya |
| `MENUNGGU_DIREKTUR` → `DRAFT` | alasan pengembalian terisi |
| `MENUNGGU_DIREKTUR` → `DITOLAK` | alasan penolakan terisi |
| `DRAFT` → `DIBATALKAN` | pelakunya **pembuat versi itu sendiri**; barisnya tidak dihapus — ADR-0055 perubahan 24 Sep 2026 |
| `WARISAN_TAK_TERPETAKAN` → *(keadaan sah)* | keadaan sasaran dipilih eksplisit; nilai asli tetap tersimpan |

### Daftar K1 — syarat pengajuan

Disalin dari kondisi asli `TreatyInCheckID` di sistem lama, **bukan dari parafrase**. Aktivitas itu
ber-blok `//` alias **mati**; daftarnya dipulihkan, bukan disusun ulang.

| # | Syarat | Keadaan — **jalur kontrak** | Keadaan — **jalur addendum** |
|---|---|---|---|
| K1-1 | `ID_CEDANT` terisi | **berjalan** | **MATI** |
| K1-2 | `ID_ASAL_BISNIS` terisi | **berjalan** | **MATI** |
| K1-3 | `SIFAT_PROPORSI` terisi | mati | **MATI** |
| K1-4 | sedikitnya satu `MATA_UANG_KONTRAK` | mati | **MATI** |
| K1-5 | kelas bisnis terisi — **kedua cabang** | mati | **MATI** |
| K1-6 | cabang proporsional: `JENIS_TREATY` terisi pada setiap `DETAIL_PROPORSIONAL` | mati | **MATI** |
| K1-7 | cabang non-proporsional: `NOMOR_LAYER`, `BAGIAN_LAYER`, `JENIS_LAYER`, `JENIS_BAGIAN_LAYER` terisi | mati | **MATI** |
| K1-8 | cabang non-proporsional: `LIMIT` terisi dan bukan nol | mati | **MATI** |

> ### Kolom ketiga ditambahkan 24 September 2026 — **jalur addendum berbeda dari jalur kontrak**
>
> Sampai hari ini tabel ini punya **satu** kolom keadaan, dan kalimat *"enam dari delapan syarat K1
> mati"* dibaca sebagai berlaku untuk seluruh pengajuan. **Ia hanya berlaku untuk jalur kontrak.**
>
> Kedua jalur dibaca langkah demi langkah:
>
> | | `TreatyInSubmit` — **kontrak** | `TreatyInSubmitEDM` — **addendum** |
> |---|---|---|
> | `Call TreatyInCheckID` | **mati** | **mati** |
> | `call TreatyInCheckError` | **hidup** | **mati** |
> | langkah keluar bergalat | **hidup** | **mati** |
>
> **Untuk jalur addendum, KEDELAPANNYA mati — termasuk K1-1 dan K1-2, yang untuk kontrak masih
> berjalan.** Seluruh pemeriksaan sebelum kirim dimatikan di sisi addendum, dan tidak ada satu pun
> penggantinya di langkah yang hidup.
>
> Bukti: `4-erd-dan-tabel-datar/CABANG-MATI-DI-JALUR-PERSETUJUAN.md` §2.3. Golongan
> **`EVIDENCED`**, dengan nomor langkah.

**Dua berjalan, enam tidak.** Keenam yang mati adalah **kemampuan baru** di sistem hasil
migrasi, bukan pelestarian — jangan dihitung gratis.

### Konsekuensi yang belum ditarik: enam syarat baru dinyalakan di hari cut-over

Bila enam syarat itu tidak pernah berjalan, maka **pengguna tidak pernah mengalaminya**. Kontrak
yang selama ini lolos pengajuan adalah kontrak yang sebagian di antaranya **tidak memenuhi** keenam
syarat itu.

Menyalakan keenamnya serentak di hari cut-over berarti orang yang bekerja seperti biasa tiba-tiba
ditolak sistem — pada hari yang sudah paling rapuh.

**Itu bukan alasan tidak menyalakannya. Itu alasan merencanakan cara menyalakannya:**

| # | Langkah |
|---|---|
| CO-6 | jalankan K1-3…K1-8 terhadap data produksi **sebelum** cut-over, dan ketahui berapa persen kontrak yang ada akan gagal. Uji S sudah mengukur sebagiannya. **DIPECAH 24 Sep 2026 — lihat CO-6a dan CO-6b** |
| **CO-6a** | jalankan **K1-3…K1-8** atas baris **kontrak**, terpisah |
| **CO-6b** | jalankan **K1-1…K1-8** atas baris **addendum**, terpisah — **kedelapannya**, karena di jalur itu K1-1 dan K1-2 pun tidak pernah dijaga |

> **Kenapa dipecah, dan kenapa satu angka gabungan akan menyesatkan:** persentase yang akan gagal
> **berbeda jauh** di antara kedua jalur. Jalur kontrak dijaga dua syarat selama ini, jalur addendum
> tidak dijaga sama sekali — maka baris addendum berpeluang jauh lebih besar memuat `ID_CEDANT` atau
> `ID_ASAL_BISNIS` kosong. Satu angka rata-rata akan **menyembunyikan** yang buruk di balik yang
> baik, dan perencanaan cut-over dibuat di atasnya.
| CO-7 | untuk syarat yang tingkat kegagalannya tinggi, pastikan itu **cacat data** dan bukan tanda syaratnya sendiri terlalu ketat — syarat yang gagal pada mayoritas kontrak biasanya salah, bukan datanya |
| CO-8 | nyalakan **berperingatan lebih dulu**, memblokir setelah orang terbiasa; tetapkan berapa lama masa peringatan itu sebelum cut-over, bukan sesudahnya |

> **Pembaca berikutnya akan membaca K1 sebagai "aturan yang sudah ada dipindahkan".** Itu keliru
> untuk enam dari delapan, dan kekeliruan itu akan muncul sebagai perkiraan pekerjaan yang terlalu
> rendah dan sebagai kejutan di hari cut-over.

### Daftar K2 — syarat persetujuan akhir

| # | Syarat |
|---|---|
| K2-1 | seluruh K1 tetap terpenuhi |
| K2-2 | INV-47 terpenuhi untuk setiap besaran yang punya baris turunan |
| K2-3 | INV-50 terpenuhi: penyebaran berjumlah 100 per sumbu pihak |
| K2-4 | setiap paket uang punya kurs, tanggal kurs, dan sumber kurs (INV-37) |
| K2-5 | tidak ada besaran yang gagal dihitung dan tersimpan kosong tanpa keterangan |
| K2-6 | tiga tingkat persetujuan sebelumnya masing-masing meninggalkan catatan |

---

## 4. Empat yang diminta perhatian khusus

### 4.1 Kunci alami per entitas — **CONSTRAINT**, INV-05…INV-16

Deklaratif seluruhnya. Tidak ada alasan menundanya ke aplikasi, dan **tanpa ini seam Adjustment
tidak berfungsi** — selisih hanya dapat dihitung atas total, dan baris yang hilang atau muncul antar
versi tidak terlihat sama sekali.

> ### Lingkupnya TIDAK SERAGAM. Ada dua bentuk, bukan satu.
>
> **Diperbaiki 24 September 2026.** Rumusan semula — satu baris `UNIQUE (ID_VERSI_KONTRAK, ‹kunci
> alami›)` untuk kedua belasnya — **salah untuk tiga**, dan rumus itulah yang akan disalin sesi DDL
> karena ia tampak seperti petunjuk pelaksanaan.

#### Bentuk A — berlingkup VERSI, sembilan entitas

`UNIQUE (ID_VERSI_KONTRAK, ‹kunci alami›)`

| INV | Entitas | Kunci alami |
|---|---|---|
| INV-05 | `LAYER` | `NOMOR_LAYER` + `BAGIAN_LAYER` |
| INV-07 | `MATA_UANG_KONTRAK` | kode mata uang |
| INV-08 | `RETENSI_CEDANT` | kelompok treaty + mata uang |
| INV-09 | `EGNPI` | kelompok treaty + mata uang |
| INV-10 | `PERIODE_PELAPORAN` | periode |
| INV-11 | `PERIODE_AKUMULASI` | periode |
| INV-12 | `TERMIN` | **kode mata uang + nomor termin** |
| INV-13 | `SKALA_KOASURANSI` | persen limit |
| INV-14 | `BATAS_PER_BAHAYA` | bahaya |

Kesembilannya menggantung **langsung** pada `VERSI_KONTRAK`, sehingga lingkup versi memang lingkup
induknya. Diperiksa satu per satu, bukan disimpulkan dari bentuknya.

> **INV-12 nyaris salah, dan diperbaiki 24 September 2026.** Ia tetap Bentuk A — lingkupnya memang
> versi — tetapi **kunci alaminya kurang satu ruas**. Di sistem lama termin disusun **per mata
> uang**: `TreatyIn.Installment[]` adalah baris per mata uang, dan termin bernomor 1..N hidup di
> dalam `InstallmentList[]` di bawahnya. Kontrak bermata uang dua karena itu punya **dua** rangkaian
> bernomor 1..N pada versi yang sama, dan `UNIQUE (ID_VERSI_KONTRAK, NOMOR_TERMIN)` akan menolak
> kontrak yang sah. Lihat `SPEC-MODEL-DATA.md` §10.16a.
>
> **Ini invarian ketiga yang lingkupnya keliru**, sesudah INV-47 dan INV-50, dan ketiganya keliru
> dengan sebab yang sama: **lingkup ditulis dari bentuk yang terlihat, bukan dari penulisnya.**

#### Bentuk B — berlingkup INDUK LANGSUNG, tiga entitas

`UNIQUE (‹kunci asing ke induk›, ‹kunci alami›)` — **bukan** `ID_VERSI_KONTRAK`.

| INV | Entitas | Induknya | Kenapa versi SALAH |
|---|---|---|---|
| INV-06 | `DETAIL_PROPORSIONAL` | `LAYER` | menaikkannya ke versi **melarang kelompok treaty yang sama muncul di dua layer berbeda pada versi yang sama** — dan itu justru hal yang biasa pada kontrak berlayer banyak |
| INV-15 | `POTONGAN` | **dua induk** (§14.1) | menaikkannya ke versi **melarang jenis potongan yang sama melekat pada kedua induknya** — padahal satu konsep di dua tempat justru yang §14.1 tetapkan |
| INV-16 | `PENYEBARAN` | induk penyebarannya | sama bentuknya dengan INV-15 |

**Daftar invariannya sendiri sudah benar** — ketiganya memang berbunyi "di dalam satu `LAYER`" dan
"di dalam satu induk". Yang salah hanya rumus ringkas di bagian ini.

#### Kunci alami yang BELUM punya nomor invarian — **LIMA, bukan tiga**

> **Dikoreksi 24 September 2026, beberapa jam setelah ditulis.** Daftar ini semula berisi **tiga**.
> Ia disusun dengan cara yang salah: saya berangkat dari entitas yang **disebut** INV-05…INV-16 lalu
> melihat tetangganya, bukan dari kolom *kunci alami* `STRUKTUR-DATA.md` disapu seluruhnya. Cara
> pertama hanya dapat menemukan yang dekat dengan yang sudah ada. **Sapuan mekanis dua arah
> menemukan dua lagi**, dan keduanya justru berbentuk **A** — bukan B seperti ketiga yang pertama.

Kelimanya punya kunci alami di `4-erd-dan-tabel-datar/STRUKTUR-DATA.md` dan **tidak punya invarian
bernomor** di §2.1. Dilaporkan sebagai lubang, **tidak ditambal di sini**.

| Entitas | Kunci alaminya di `STRUKTUR-DATA.md` | Induk | Bentuk |
|---|---|---|---|
| `PORTOFOLIO` | jenis + jenis portofolio | `VERSI_KONTRAK` | **A** |
| `DOKUMEN_KONTRAK` | pengenal dokumen | `VERSI_KONTRAK` | **A** |
| `BAGIAN` | satu per layer | `LAYER` | **B** |
| `RINCIAN_PENYEBARAN` | pihak | `PENYEBARAN` | **B** |
| `NILAI_PENYEBARAN` | mata uang | `RINCIAN_PENYEBARAN` | **B** |

Menomori kelimanya adalah pekerjaan **sesi to-spec**, dan ketiga yang berbentuk B **didahulukan** —
alasannya di `5-tiket/DAFTAR-PEKERJAAN.md` §5.0.

**Dua hal yang BUKAN lubang, dicatat supaya tidak dibaca sebagai lubang:**

| | |
|---|---|
| `KONTRAK` | kunci alaminya ada, dan ketiadaan nomornya **disengaja**: ia **memperingatkan, tidak melarang** (ADR-0040 §2). Sebuah peringatan bukan invarian |
| `CATATAN_PERSETUJUAN`, `JEJAK_PERUBAHAN` | `STRUKTUR-DATA.md` menyatakan keduanya **tidak punya** kunci alami — urutan waktu yang membedakan barisnya. Ketiadaan yang **dinyatakan** adalah keputusan, dan keputusan tidak perlu ditambal |

**Enam tabel acuan tidak punya kunci alami tertulis sama sekali** di `STRUKTUR-DATA.md` §3 — bukan
"tidak punya", melainkan **tidak dinyatakan**. Kodenya hampir pasti kunci alaminya, tetapi
"hampir pasti" bukan dasar menulis `UNIQUE`. Itu bagian dari pekerjaan §10 untuk keenamnya.

#### Sapuan arah sebaliknya — dijalankan 24 September 2026, **bersih**

Arah yang belum pernah diperiksa: **adakah invarian yang menyebut entitas yang tidak ada di
`STRUKTUR-DATA.md`?**

Ke-78 token berhuruf kapital di berkas ini disapu terhadap daftar entitas. Yang tidak ada di
`STRUKTUR-DATA.md` berjumlah 52, dan **tidak satu pun entitas**: seluruhnya nama atribut
(`NOMOR_LAYER`, `SIFAT_PROPORSI`, …), nilai enumerasi (`DRAFT`, `DISETUJUI`, `DIBATALKAN`, …),
kata kunci Oracle (`CHECK`, `UNIQUE`, `SEQUENCE`, `NOVALIDATE`, …), atau kamus data
(`USER_CONSTRAINTS`, `USER_TRIGGERS`).

**Tidak ada invarian yang menggantung pada entitas hantu.** Sapuannya sekali jalan, dan hasilnya
dicatat di sini supaya tidak dijalankan ulang tanpa sebab.

### 4.2 Kunci alami tidak berubah sepanjang hidup kontrak — **TRIGGER**, INV-19

Tidak dapat dinyatakan deklaratif: Oracle tidak punya kolom "hanya boleh ditulis sekali".

**Dipilih trigger, bukan aplikasi**, dengan alasan yang memenuhi syarat §0: kunci alami menentukan
**kontrak mana** sebuah angka melekat. Mengubahnya memindahkan seluruh riwayat uang sebuah kontrak
ke kontrak lain tanpa satu pun angka berubah — kerusakan yang tidak terlihat dari mana pun.

Bentuknya: trigger `BEFORE UPDATE` yang menolak perubahan pada kelima kolom lapisan beku.

### 4.3 Versi berkeadaan TERMINAL tidak pernah berubah nilainya — **TRIGGER**, INV-24

> ### Diperluas 24 September 2026 — dari satu keadaan menjadi tiga
>
> Judul dan seluruh uraian bagian ini semula ditulis untuk **`DISETUJUI` saja**. Itu lubang, dan ia
> baru terlihat ketika keadaan kedelapan dipasang.
>
> **Bandingkan dua invarian yang bersebelahan:**
>
> | | Menyebut |
> |---|---|
> | INV-23 | `DISETUJUI`, `DITOLAK`, `DIBATALKAN` — **tiga** |
> | INV-24 *(sebelum perluasan)* | `DISETUJUI` — **satu** |
>
> Artinya versi ber-`DITOLAK` **dapat diubah nilainya oleh siapa pun, kapan saja, selamanya.**
> Keadaannya tidak dapat berpindah, tetapi isinya terbuka.
>
> **Akibatnya tidak sepele:**
>
> - Versi yang ditolak adalah **bukti atas keputusan penolakan**. Bila isinya berubah sesudah
>   ditolak, `CATATAN_PERSETUJUAN` menunjuk pada sesuatu yang bukan lagi yang dinilai — pemberi
>   keputusan menolak angka A, yang tersimpan sekarang angka B.
> - Versi yang dibatalkan adalah **satu-satunya jejak** apa yang sedang disusun ketika ia dibuang.
>   Bila ia dapat disunting, jejak itu tidak berarti.
> - `JEJAK_PERUBAHAN` akan memuat perubahan atas versi yang **sudah mati**, tanpa satu pun
>   perpindahan yang menjelaskannya — membingungkan siapa pun yang membacanya.
>
> **Yang benar: TERMINAL BERARTI BEKU, BUKAN HANYA BERHENTI BERPINDAH.**
>
> **Bentuk kegagalannya adalah keluarga yang sudah dikenal di proyek ini:** dua invarian yang
> bersandingan, hampir sama bunyinya, satu menyebut tiga keadaan dan satu menyebut satu — dan tidak
> ada yang menyadarinya selama daftar keadaan terminalnya masih dua. **Keadaan kedelapan tidak
> menciptakan lubang ini; ia hanya membuatnya cukup besar untuk terlihat.**
>
> **Yang berubah pada rancangan trigger: tidak ada.** Trigger `BEFORE UPDATE OR DELETE` di bawah
> menjangkaunya tanpa perubahan bentuk; yang berubah hanya **syaratnya**, dari satu keadaan menjadi
> tiga.

**Ini prasyarat seluruh seam Adjustment** — dan sejak perluasan di atas, **jangkauannya lebih luas
daripada alasan yang melahirkannya**: seam Adjustment hanya menuntut `DISETUJUI` beku, tetapi
`DITOLAK` dan `DIBATALKAN` ikut dibekukan karena keduanya bukti, bukan karena seam memerlukannya.

Pertanyaannya bukan "bisa tidak" melainkan "seberapa jauh".

| Cara | Menegakkan? | Putusan |
|---|---|---|
| `CHECK` | tidak — `CHECK` tidak dapat membaca nilai lama | — |
| Hak akses per objek | sebagian — mencegah penulis yang salah, bukan penulis yang benar salah menulis | dipakai sebagai lapis kedua |
| **Trigger — bentuknya BERBEDA antara induk dan anak** | **ya**, lihat §4.3a | **dipilih** |

#### 4.3a Bentuk triggernya berbeda antara induk dan anak, dan bedanya bukan gaya

> **Ditulis 24 September 2026.** Bentuk semula — `BEFORE UPDATE OR DELETE` untuk induk **dan** anak
> — **lengkap untuk induk dan BOLONG untuk anak.**

| Tabel | Bentuk trigger | Kenapa |
|---|---|---|
| **`VERSI_KONTRAK`** (induk) | `BEFORE UPDATE OR DELETE` | keadaan terminal ada **pada baris itu sendiri**. Baris yang sudah ada tidak dapat diubah maupun dihapus, dan itu sudah seluruh permukaannya |
| **seluruh entitas anak** | **`BEFORE INSERT OR UPDATE OR DELETE`** | `UPDATE OR DELETE` hanya menjaga baris yang **sudah ada**. Ia **tidak menolak INSERT** |

**Kenapa INSERT pada anak membatalkan maksud INV-24 seluruhnya.** Tanpa penjagaan INSERT, terhadap
versi yang sudah `DISETUJUI` seseorang masih dapat menambahkan `LAYER` baru, `POTONGAN` baru pada
layer yang sudah ada, baris `PENYEBARAN` atau `RINCIAN_PENYEBARAN` baru, `MATA_UANG_KONTRAK`,
`RETENSI_CEDANT`, `EGNPI`, dan seterusnya.

**Tidak satu pun angka yang sudah ada berubah. Yang berubah adalah jumlahnya** — dan di model ini
**jumlah baris anak justru yang menentukan angkanya**, karena INV-47, INV-50, dan INV-51 seluruhnya
tentang penjumlahan baris turunan.

> Jadi versi yang disetujui tidak dapat **diubah** nilainya, tetapi dapat **bertambah** nilainya.
> Itu pintu yang tidak dijaga, dan ia membatalkan INV-24 tanpa melanggar satu kata pun dari
> rumusan lamanya.

**Dua hal yang harus dinyatakan di muka, karena sesi DDL akan menanyakannya:**

| | |
|---|---|
| **a. Trigger anak membaca keadaan INDUKNYA, bukan keadaan dirinya** | baris anak **tidak punya keadaan**. Triggernya membaca `VERSI_KONTRAK`, dan **ongkosnya satu pembacaan per baris yang disisipkan atau diubah**. Disebutkan di muka, seperti ongkos `REFRESH ON COMMIT` disebutkan di muka (§5.1) |
| **b. Pengecualian migrasi tetap lewat §4.3b** | trigger dimatikan selama jendela migrasi, dan **CO-3** membuktikannya menyala kembali. **Tidak ada jalan keluar kedua**, dan tidak perlu ada |

**Trigger INV-19 diperiksa dengan pertanyaan yang sama, dan ia benar apa adanya.**
`BEFORE UPDATE` pada `KONTRAK` sudah cukup: kunci alami adalah **kolom pada baris yang sudah ada**,
dan tidak ada INSERT yang dapat mengubah kunci alami baris yang sudah lahir. Penjagaan INSERT di
sana akan menolak pembuatan kontrak itu sendiri.

**Dan ini harus dinyatakan terang:** trigger menegakkannya untuk setiap jalur yang melewati basis
data, **termasuk** perbaikan data manual dan prosedur — yang justru jalur paling berbahaya dan yang
tidak dapat dijangkau aplikasi.

Yang **tidak** dijangkau trigger: `ALTER TABLE ... DISABLE TRIGGER`, dan pemuatan langsung yang
melewati trigger. Keduanya menuntut hak istimewa yang **tidak diberikan** kepada akun aplikasi
(§8 sesi DDL). Jadi jaminannya bersandar pada **basis data**, bukan pada aplikasi — dan itu
menurunkan bobot risikonya, bukan menaikkannya.

### 4.3b Trigger punya sakelar mati, dan migrasi AKAN memerlukannya

Trigger dapat dimatikan satu perintah. Dan migrasi **akan** memerlukannya: mengimpor versi warisan
yang sudah berkeadaan `DISETUJUI` berarti menulis baris yang INV-24 justru ada untuk menolak.

Maka prosedur di sekelilingnya **bagian dari penegakan, bukan pelengkapnya**, dan ditetapkan di
sini — bukan nanti.

| Hal | Ketetapan |
|---|---|
| Siapa boleh mematikan | **hanya pemilik objek** skema `TREATY_MASUK`. Akun aplikasi `TREATY_MASUK_APP` **tidak pernah** diberi `ALTER` pada objek mana pun, sehingga ia tidak dapat mematikan trigger walaupun kodenya mencoba. |
| Kapan boleh dimatikan | **hanya selama jendela migrasi**, dan hanya oleh pekerjaan migrasi. Tidak ada keadaan operasi normal yang membenarkannya. |
| Apa yang tercatat | pematian dan penyalaan kembali masing-masing menulis satu baris ke jejak operasi: siapa, kapan, trigger mana, dan alasannya. Jejak ini **di luar** `JEJAK_PERUBAHAN`, karena ia fakta tentang skema, bukan tentang kontrak. |
| Pemeriksaan wajib sesudahnya | kueri `USER_TRIGGERS` di `UJI-NEGATIF-INVARIAN.md` §3 dijalankan **sebagai langkah cut-over yang harus lulus**, bukan sebagai pemeriksaan sukarela. Hasil bukan nol menghentikan cut-over. |

> **Trigger yang lupa dinyalakan adalah bentuk yang sama dengan materialized view yang basi**, dan
> keduanya bentuk yang sama dengan aturan mati di sistem lama: penjaga yang terbaca dan tidak
> menjaga. Ketiganya sekarang punya pemantaunya.

**Masuk daftar cut-over, dan harus ada di sana sebelum sesi DDL menulis trigger-nya:**

| # | Langkah cut-over |
|---|---|
| CO-1 | sebelum migrasi: catat daftar trigger yang akan dimatikan |
| CO-2 | sesudah migrasi: nyalakan kembali seluruhnya |
| CO-3 | **buktikan** menyala — `USER_TRIGGERS` tidak mengembalikan satu baris pun berstatus bukan `ENABLED` |
| CO-4 | jalankan N-5a…N-5f satu kali terhadap data yang sudah dimigrasi; kegagalan menghentikan cut-over |
| CO-5 | `USER_CONSTRAINTS` tidak mengembalikan satu baris pun `NOVALIDATE` atau `DISABLED` |

CO-4 ada karena trigger yang menyala belum tentu trigger yang bekerja: menyalakannya kembali setelah
tabel dasarnya berubah dapat menghasilkan trigger `ENABLED` tetapi tidak valid.

### 4.3c `WARISAN_TAK_TERPETAKAN` memakai satu-satunya jatah versi tak-terminal — INV-25

**Ditulis 24 September 2026.** Bukan cacat rancangan; akibat yang harus diketahui **sebelum**
cut-over, bukan ditemukan di hari pertama.

`WARISAN_TAK_TERPETAKAN` punya perpindahan keluar (`PERBAIKAN_WARISAN`), jadi ia **bukan terminal**.
INV-25 menetapkan paling banyak **satu** versi per kontrak berada di keadaan tak-terminal.

Gabungkan keduanya:

> **Setiap kontrak warisan yang mendarat di `WARISAN_TAK_TERPETAKAN` tidak dapat dibuatkan versi
> baru** sampai seseorang memilihkan keadaan sah untuknya.

Itu **mungkin memang yang dikehendaki** — kontrak yang keadaannya tidak diketahui sebaiknya tidak
diubah dulu. Tetapi tiga hal harus diketahui lebih dulu:

| # | Pertanyaan | Jenisnya | Di mana dijawab |
|---|---|---|---|
| **a** | **Berapa banyak** kontrak yang akan mendarat di sana? | **data** | **Uji S-7** di `PERMINTAAN-DBA-1-UJI-A-SAMPAI-H.sql` — angkanya tinggal dibaca |
| **b** | **Siapa** yang berwenang memilihkan keadaan sahnya, dan atas dasar apa? | **bisnis** | belum ditanyakan |
| **c** | Adakah kontrak yang perlu di-addendum **segera** setelah cut-over dan kebetulan mendarat di sana? | **cut-over** | daftar cut-over |

**Kenapa (a) mendesak:** bila jumlahnya ribuan, maka di hari pertama **ribuan kontrak tidak dapat
di-addendum**, dan pekerjaan menyelesaikannya satu per satu menjadi pekerjaan nyata yang **belum ada
di rencana siapa pun**.

**Kenapa (b) belum terjawab:** `5-tiket/DAFTAR-PEKERJAAN.md` P-52 menyebut **pengisi kontrak** sebagai
pelakunya. Apakah itu cukup untuk memutuskan bahwa sebuah kontrak warisan sebenarnya **sudah
disetujui**? Itu keputusan yang berakibat uang, dan pelakunya belum tentu tepat. Pertanyaan bisnis,
bukan pertanyaan bentuk.

**Kenapa (c) tidak dapat ditunda:** bila ada, penyelesaiannya harus **didahulukan**, dan urutannya
perlu diketahui sebelum hari pertama — bukan disusun saat orang sudah menunggu.

### 4.4 Rekonsiliasi penyebaran — **CONSTRAINT**, bukan aplikasi. INV-47, INV-50, INV-51

Pemilik proses memperkirakan ini jatuh ke aplikasi. **Tidak perlu**, dan §5.1 menjelaskan caranya.

Pengecualian bernama tetap berlaku dan **bagian dari invariannya, bukan pengecualian terhadapnya**:
kelipatan (INV-48) dan besaran berulang (INV-49) tidak pernah masuk perhitungan partisi.

#### 4.4a SUMBUNYA JENIS REASURANSI, BUKAN PIHAK — diperbaiki 24 September 2026

INV-47 dan INV-50 semula berbunyi **"per sumbu pihak"**. Sumbu itu **tidak ada** di keluarga
penyebaran, dan salahnya bukan salah ketik: sesi DDL menulis `GROUP BY` dari kalimat ini.
**Constraint yang mengelompokkan menurut kolom yang keliru akan lulus pada data yang seharusnya
gagal.**

**Yang dijalankan.** Seluruh 866 pasangan kelas–properti disapu untuk ruas beridentitas pihak —
`ReinsID`, `ReinsName`, `BrokerName`, `ReinsurerID`, `MemberID`, `ParticipantID`, dan sejenisnya.
Hasilnya empat kelas, dan **tidak satu pun ada di keluarga penyebaran**:

| Kelas | Ruas | Kedudukan |
|---|---|---|
| `Data-TreatyInShareReins` | `ReinsID`, `ReinsName`, `BrokerName` | `RETRO_KELUAR` — **GEL-2** |
| `Int-AGENT`, `Int-TREATYREINSURER` | `ClientID`, `ClientName`, `ReinsurerID` | master di luar skema |
| `Data-TreatyInLimitsDetail` | `ReisuredParticipant` *(salah eja ada di sumbernya)* | satu ruas, **belum pernah terlihat** — diadili di §10.4 |

Kelas `Data-TreatyInLimitsSpreading` — penyandang `PENYEBARAN`, `RINCIAN_PENYEBARAN`, dan
`NILAI_PENYEBARAN` — punya **26 properti dan nol di antaranya identitas pihak**. Sumbunya
`ReinsTypeID`, dan `ParentReinsTypeID` di sebelahnya menunjukkan jenis reasuransi itu **bersusun**.

Bahkan ruas yang **bernama** pihak ternyata berisi jenis. `Activity/SetSpreadName.xml` mengisi
`BreakDownSprdList(…).ReinsID` dari `pxResults(…).ReinsTypeID` dan `.ReinsName` dari
`.ReinsTypeName`. **Namanya menyebut pihak; kondisinya menulis jenis, dan kondisinya yang menang.**

> **Maka pengelompokannya: `GROUP BY ‹kunci asing ke induk penyebaran›, ID_JENIS_REASURANSI`** —
> dan tidak pernah menurut pihak, karena pihak baru muncul di GEL-2.

**Satu-satunya calon penyangkal dibuka, dan ia gugur.** Sapuan di atas menyisakan satu ruas yang
namanya menyebut pihak di dalam keluarga limit: **`ReisuredParticipant`** pada
`Data-TreatyInLimitsDetail` — kelas dengan **146 rujukan jalur yang ditolak penjaga**, ketiga
terbesar, dan karena itu salah satu tempat yang paling lama tidak terlihat. Menyimpulkan "tidak ada
sumbu pihak" tanpa membukanya akan menjadi kesimpulan yang calon penyangkalnya belum diperiksa.

Yang ditemukan setelah dibuka:

| Yang diperiksa | Hasil |
|---|---|
| **bentuknya** | **skalar**, bukan daftar. Satu nilai per baris `Detail[]` — ia tidak dapat memuat baris per pihak, dan sumbu menuntut baris |
| **siapa menulisnya** | **tidak ada satu aturan pun**. Ia muncul **hanya** di `Section/DetailLimits.xml`, sebagai kolom bertajuk *"Reisured Participant"* |
| **siapa membacanya** | **tidak ada**. Nol aktivitas, nol *data transform*, nol *when* |
| **dapatkah diisi** | **tidak.** Blok kendalinya sendiri ber-`pyDisabled = true`; kolom tetangganya yang memang disunting, `QSPct`, tidak punya penanda itu |

> Jadi ia **kolom layar mati**: tidak diisi, tidak dibaca, tidak dapat disentuh. Ia **label tanpa
> angka di sebelahnya** — dan satu ruas pihak sendirian adalah label, bukan sumbu.

**Kesimpulannya berdiri, dan sekarang berdiri lebih kuat**, karena ia menyebutkan bahwa calonnya
sudah dibuka. Yang tersisa: `ReisuredParticipant` **tetap diadili di §10.4** sebagai salah satu dari
17 properti `LimitsDetail` yang menuntut keputusan — dugaan awalnya **dibuang**, sekeluarga dengan
`PositionUsername` (§12.5), karena nama orang atau pihak sebagai teks bebas pada kontrak melewati
ADR-0044.

#### 4.4c INV-47 HARUS MENYEBUT NAMA KOLOM, BUKAN GAGASAN

Perbaikan pertama INV-47 — *"menurut sumbu yang membedakan baris turunannya"* — **ditolak pemilik
proses pada hari yang sama**, dengan dua keberatan yang keduanya praktis:

| Keberatan | Isinya |
|---|---|
| **ia tidak dapat gagal** | untuk sekumpulan baris mana pun **selalu** ada sumbu yang membedakannya — kalau tidak ada, barisnya bukan baris yang berbeda. Yang benar untuk semua keadaan **tidak menolak apa pun**, dan kriteria yang tidak bisa gagal bukan kriteria |
| **ia tidak dapat dikompilasi** | sesi DDL menulis *materialized view* ber-`GROUP BY` yang nyata. "Sumbu yang membedakan" **tidak menghasilkan satu pun nama kolom**, dan yang menulisnya akan memilih sendiri — diam-diam, tanpa ditandai |

INV-50 tidak punya masalah ini karena ia menyebut kolomnya. INV-47 tidak bisa: ia berlaku atas
banyak entitas sekaligus, dan sumbunya berbeda-beda.

**Penyelesaiannya memberi INV-47 alamat, bukan gagasan.** `5-tiket/DAFTAR-PEKERJAAN.md` §5.0b
mewajibkan setiap entitas menyatakan **`SUMBU_REKONSILIASI`-nya** — nama kolom yang dipakai
mengelompokkan — berdampingan dengan kunci alami dan lingkupnya, dan entitas yang **tidak punya**
baris turunan menyatakan itu juga.

INV-47 lalu menunjuk ke tempat yang berisi **nama kolom**. Satu invarian umum tetap satu invarian,
dan penegakannya punya alamat.

#### 4.4b INV-50 TIDAK PERNAH GAGAL DI SISTEM LAMA, DAN ITU BUKAN BUKTI

Persentase rincian penyebaran di sistem lama **disalin dari tabel master jenis reasuransi** —
`SharePct = pyReportContentPage.pxResults(idx).Pct`, dan `Amount` dihitung darinya. Master itu
berjumlah 100.

Maka INV-50 **tidak pernah dapat dilanggar oleh data yang dihasilkan jalur itu**. Invariannya benar,
penegakannya benar, dan ia **tidak pernah menguji apa pun**.

> **Ini bentuk paling bersih dari pola "cacat bersembunyi di balik nilai bawaan"** yang sudah
> tercatat di `PENGETAHUAN.md`: dua hal yang berbeda menghasilkan angka yang sama selama nilainya
> bawaan, sehingga perbedaannya tidak pernah terlihat.

**Ditulis di sini supaya orang yang melihat INV-50 "selalu lulus" tidak menyimpulkan ia tidak
perlu.** Yang benar sebaliknya: begitu sistem baru mengizinkan penyebaran **disunting** — dan
ADR-0036 mengizinkannya, dengan asal-usul *"disemai lalu disunting"* — INV-50 menjadi satu-satunya
hal yang mencegah kapasitas terbagi lebih atau kurang dari seratus persen. **Ia mulai bekerja pada
hari sistem baru menyala, bukan sebelumnya.**

---

## 5. Yang ditarik kembali ke basis data, dan yang benar-benar tinggal di aplikasi

### 5.1 Agregat lintas baris **dapat** deklaratif di Oracle

Invarian "jumlah baris anak sama dengan nilai induk" biasanya dianggap mustahil dinyatakan
deklaratif. Di Oracle ia mungkin, dan caranya baku:

> **Materialized view ber-`REFRESH COMPLETE ON COMMIT`** yang memuat selisih antara jumlah anak dan
> nilai induk, **ditambah `CHECK` constraint pada view itu** yang menuntut selisihnya nol.
> Pelanggaran membuat `COMMIT` gagal.

Empat invarian ditarik kembali dengan cara ini: INV-47, INV-50, INV-51, dan INV-31 pada bentuk
lintas barisnya.

**Ongkosnya nyata dan disebut di muka:** `ON COMMIT` menyerialkan commit yang menyentuh kelompok
baris yang sama. Bila profil beban kelak menunjukkan itu tidak tertahankan, keempatnya turun ke
aplikasi — dan penurunan itu **keputusan sadar yang ditulis**, bukan bawaan yang tidak pernah
dicoba. Sesi DDL menuliskannya di `CATATAN-DDL.md`.

### 5.2 Kesebelas yang menuntut aplikasi, satu per satu

Angka 17,5 % hanya kabar baik bila kesebelasnya memang **yang tidak mungkin**, bukan **yang belum
dicoba**. Daftar dengan alasan per baris dapat diperiksa orang lain; persentase tidak.

| # | Isinya | Kenapa tidak dapat dinyatakan di basis data |
|---|---|---|
| INV-26 | `WARISAN_TAK_TERPETAKAN` hanya dimasuki lewat migrasi | Basis data tidak dapat membedakan penulisan migrasi dari penulisan biasa — keduanya `INSERT` dari sesi yang sah. Yang membedakan adalah **asal-usul** perintahnya, dan itu ada di luar data. |
| INV-27 | Kelengkapan per perpindahan terpenuhi sebelum perpindahan | Sebagian syaratnya menyangkut **siapa pelakunya** — "pemberi keputusan bukan pengisi kontraknya", "berbeda dari tingkat sebelumnya" — dan identitas pelaku bukan kolom pada kontrak. |
| INV-35 | Baris `SURPLUS` menuntut baris `QUOTA_SHARE` pada versi yang sama | **Keberadaannya** dapat ditegakkan MV seperti INV-47. Yang **tidak** dapat adalah bagian keduanya: "ketiadaannya adalah kegagalan yang **dilaporkan**" — kepada siapa, dengan pesan apa. Basis data menolak; ia tidak melapor. Bagian penolakannya karena itu **dapat** turun ke constraint bila dikehendaki, dan itu dicatat sebagai pilihan yang terbuka. |
| INV-52 | Premi bruto − potongan = premi bersih | Premi bersih adalah **turunan yang tidak disimpan** (ADR-0037). Tidak ada kolom untuk diperiksa; yang diperiksa adalah kode yang menghitungnya. |
| INV-55 | Periode pelaporan di dalam periode kontrak | Lintas tabel, dan **tidak menyentuh uang** — ambang trigger di §0 tidak terpenuhi. Pelanggarannya menghasilkan jadwal yang salah, bukan angka yang salah. |
| INV-56 | Periode akumulasi di dalam periode kontrak | Sama dengan INV-55. |
| INV-58 | Tidak ada turunan yang disimpan | Invarian **atas rancangan**, bukan atas isi baris: ia dilanggar dengan menambah kolom, bukan dengan menulis nilai. Diperiksa saat tinjauan skema. |
| INV-59 | Tidak ada nilai yang disalin dari entitas lain | Sama — dilanggar dengan menambah kolom salinan. |
| INV-60 | Satu fakta satu penulis | Sama — ia tentang siapa yang boleh menulis apa, sebagian ditegakkan hak akses per objek, sisanya tinjauan kode. |
| INV-61 | Arsip JSON tidak punya jalur baca | Sama — dilanggar dengan menulis kueri yang membacanya, dan basis data tidak dapat melarang `SELECT` atas kolomnya sendiri tanpa mencabut akses ke tabelnya. |
| INV-63 | Aturan potongan tertulis sekali | Sama — dilanggar dengan membuat dua tabel. Diperiksa saat tinjauan skema. |

**Pemeriksaan jujur atas daftar ini.** Enam dari sebelas — INV-58, 59, 60, 61, 63, dan sebagian
INV-35 — **bukan invarian atas data melainkan atas rancangan**: ia dilanggar dengan mengubah skema,
bukan dengan menulis baris. Keenamnya tetap dihitung sebagai aplikasi di §1 supaya angkanya tidak
diperkecil secara semu.

Yang benar-benar menuntut **kode aplikasi yang berjalan** tinggal lima: INV-26, INV-27, INV-52,
INV-55, INV-56 — **7,9 %**.

Dan satu di antara kesebelas **dapat ditarik lebih jauh**: bagian "keberadaan baris `QUOTA_SHARE`"
pada INV-35 dapat menjadi constraint lewat cara §5.1. Yang tinggal di aplikasi hanyalah pelaporannya.
Itu dicatat di sini sebagai **pilihan yang terbuka dan belum diambil**, bukan sebagai batas yang
tidak dapat dilewati — supaya sesi DDL tahu ia boleh mengambilnya.

---

## 6. Yang belum ditulis, dan itu disengaja

| Hal | Sebab |
|---|---|
| Nama constraint | milik sesi DDL; §16 mengatur bentuknya (`CK_x_1`, tidak mengeja kolom) |
| Angka presisi pada INV-45 | keputusan gerbang sesi DDL |
| Bentuk fisik `POTONGAN` berinduk dua | §14.1 — satu konsep sudah diputuskan; bentuknya milik DDL |
| Invarian atas `RETRO_KELUAR` dan `PENCAPAIAN` | keduanya G2/G3, di luar batas gelombang ini |
| Invarian yang menyangkut `AccountingMode` | menunggu Uji X-2 (§14.4) |

---

## 7. Lima kunci alami yang belum bernomor — ditulis 24 September 2026 (langkah 3 sesi to-spec)

`SPEC-MODEL-DATA.md` §11.2 menyisakan *"menomori lima kunci alami yang belum punya invarian"*.
Kelimanya dicari dengan menyisir setiap baris **"Kunci alami dan lingkupnya"** di §10 dan memeriksa
apakah ia menyebut nomor invarian. **Yang ditemukan enam, bukan lima** — dan selisihnya dijelaskan
di §7.3.

> ### Lingkupnya dibaca dari AKTIVITAS YANG MENAMBAH BARIS, bukan dari bentuk pohonnya
>
> `CONTEXT.md` §2.2c: kekeliruan ini sudah terjadi **tiga kali** di modul ini, dan bentuknya selalu
> sama — daftar yang tampak menggantung langsung pada induk besarnya ternyata punya satu tingkat
> pengelompokan di antaranya, sehingga constraint yang ditulis dari bentuk terlihat **menolak data
> yang sah**. Untuk kelima invarian di bawah, penulis barisnya dibuka satu per satu.

### 7.1 Kelimanya

#### INV-64 — satu `BAGIAN` per `LAYER`

| | |
|---|---|
| **Bunyi** | Pada `BAGIAN`, `ID_LAYER` unik. Tidak ada dua baris `BAGIAN` dengan `ID_LAYER` yang sama |
| **Bentuk** | B — lingkupnya `LAYER`, dan kuncinya **hanya** kunci asing itu sendiri |
| **Dapat gagal?** | ya — dua baris `BAGIAN` ber-`ID_LAYER` sama |
| **Dapat dikompilasi jadi nama kolom?** | ya — `UNIQUE (ID_LAYER)` pada `BAGIAN` |
| **Dibaca dari penulisnya** | `TreatyInNonAddItem.xml` menulis `TreatyIn.Share(<APPEND>)` **di dalam gelung ber-`pyStepsObjectName = TreatyIn.Limits`** — tepat satu baris `Share` ditambahkan **per layer**. Bukan dibaca dari bentuk pohon |

#### INV-65 — jenis reasuransi unik di dalam satu `PENYEBARAN`

| | |
|---|---|
| **Bunyi** | Pada `RINCIAN_PENYEBARAN`, pasangan (`ID_PENYEBARAN`, `ID_JENIS_REASURANSI`) unik |
| **Bentuk** | B — lingkupnya `PENYEBARAN`, **bukan** versi |
| **Dapat gagal?** | ya — dua rincian dengan jenis reasuransi sama di bawah satu penyebaran |
| **Dapat dikompilasi jadi nama kolom?** | ya — `UNIQUE (ID_PENYEBARAN, ID_JENIS_REASURANSI)` |
| **Dibaca dari penulisnya** | `SetSpreadName.xml` menulis `Primary.SpreadingList(Local.idxSpreadList).BreakDownSprdList(…)` — **diindeks oleh induknya** `Local.idxSpreadList`, di dalam gelung atas `pyReportContentPage.pxResults` (baris susunan retro baku). Satu rincian per jenis, **di dalam** satu baris penyebaran. `SetSpreadingXOL.xml` menulis bentuk yang sama pada cabang non-proporsional |
| **Catatan lingkup** | ini **persis** bentuk kekeliruan yang §2.2c peringatkan. Ditulis dari bentuk pohon, kuncinya akan naik ke versi — dan itu **melarang jenis reasuransi yang sama muncul di dua penyebaran berbeda**, padahal susunan retro baku memang memberi jenis yang sama kepada beberapa penyebaran |

#### INV-66 — arah + jenis portofolio unik di dalam satu versi

| | |
|---|---|
| **Bunyi** | Pada `PORTOFOLIO`, rangkap tiga (`ID_VERSI_KONTRAK`, `ARAH_PORTOFOLIO`, `JENIS_PORTOFOLIO`) unik |
| **Bentuk** | A — lingkupnya versi |
| **Dapat gagal?** | ya — dua baris "masuk + premium" pada satu versi |
| **Dapat dikompilasi jadi nama kolom?** | ya — `UNIQUE (ID_VERSI_KONTRAK, ARAH_PORTOFOLIO, JENIS_PORTOFOLIO)` |
| **Dibaca dari penulisnya** | `TreatyInMappingDataconvertProp.xml` menulis `TreatyIn.Portfolio(<APPEND>)` di dalam gelung atas `TreatyOffer.OfferTreatyInList(1).PortfolioList` — sumbernya daftar portofolio **penawaran**, yang bertingkat penawaran, bukan bertingkat layer. Lingkup versi terbaca, bukan disimpulkan |

#### INV-67 — dokumen unik di dalam satu versi

| | |
|---|---|
| **Bunyi** | Pada `DOKUMEN_KONTRAK`, pasangan (`ID_VERSI_KONTRAK`, `ID_DOKUMEN`) unik |
| **Bentuk** | A |
| **Dapat gagal?** | ya — satu dokumen dilampirkan dua kali ke versi yang sama |
| **Dapat dikompilasi jadi nama kolom?** | ya — `UNIQUE (ID_VERSI_KONTRAK, ID_DOKUMEN)` |
| **Dibaca dari penulisnya** | **TIDAK DAPAT — dan ini dinyatakan, bukan dilewati.** Sapuan atas kelima bentuk penulis menemukan **nol** penambah baris lampiran. Lampiran ditangani mesin lampiran bawaan Pega, yang **tidak terekspor** (`L-10`). Lingkupnya karena itu **keputusan rancangan**, bukan pembacaan |
| **Syarat pembalikan** | bila ekspor kedua menunjukkan lampiran bertingkat **kontrak**, bukan versi, invarian ini naik satu tingkat |

#### INV-68 — `KODE` unik di keenam tabel acuan

| | |
|---|---|
| **Bunyi** | Pada `MATA_UANG`, `JENIS_POTONGAN`, `JENIS_REASURANSI`, `BAHAYA`, `KELOMPOK_TREATY`, `KELAS_BISNIS`: `KODE` unik **di seluruh tabel** |
| **Bentuk** | tak berinduk — satu-satunya di modul ini |
| **Dapat gagal?** | ya — dua baris berkode sama |
| **Dapat dikompilasi jadi nama kolom?** | ya — `UNIQUE (KODE)` pada masing-masing dari enam tabel |
| **Catatan** | satu invarian bernomor, **enam constraint**. Ia tidak dipecah menjadi enam nomor karena bunyinya identik dan pembacanya sama; `PETA-INVARIAN-KE-DDL.md` yang memecahnya menjadi enam baris |

### 7.2 Dua uji yang harus dilewati, diterapkan pada kelimanya

`aturan-harus-menyebut-kolom`: sebuah rumusan bukan invarian kecuali **(a)** ia dapat gagal dan
**(b)** ia dapat dikompilasi menjadi nama kolom. Kelimanya lulus keduanya — tabel di atas
menyebutkan jawabannya satu per satu, bukan menyatakannya di akhir.

### 7.3 Yang KEENAM, dan kenapa ia tidak dinomori sekarang

`§10.8 NILAI_PENYEBARAN` juga berbunyi *"belum punya nomor invarian"* — kunci alaminya kode mata
uang, Bentuk B di dalam `ID_RINCIAN_PENYEBARAN`.

**Ia tidak dinomori** karena §10.23c butir 3 **menangguhkan kedudukan entitasnya sendiri**: bukti
yang ada tidak memisahkan apa pun, dan **Uji AD** yang memutuskan apakah ia entitas tersendiri atau
melebur. Menomori invarian atas entitas yang mungkin tidak ada berarti menulis constraint yang
mungkin tidak punya tabel.

| | |
|---|---|
| **Siapa menutup** | sesi to-spec berikutnya, **sesudah** Uji AD kembali |
| **Yang menagih** | §10.8 sendiri tetap berbunyi *"belum punya nomor invarian"*, dan §11.2 tetap memuat barisnya sampai Uji AD menjawab |

### 7.4 Yang TIDAK punya kunci alami, dan itu keputusan — bukan kekosongan

Tiga entitas menyatakan **"Kunci alami — TIDAK ADA, dan itu keputusan"**: `CATATAN_PERSETUJUAN`
(§10.20), `JEJAK_PERUBAHAN` (§10.21), `PERISTIWA_KONTRAK` (§10.21a). Alasannya satu dan sama:
peristiwa yang sama dapat terjadi dua kali pada versi yang sama, dan melarangnya berarti melarang
kenyataan. **Ketiganya diperiksa ulang di langkah ini dan pernyataannya berdiri.**

---

## 8. Bentuk fisik paket uang — dan tiga invarian yang dibebaskan olehnya

**Ditulis 24 September 2026 (P-8).**

`SPEC-MODEL-DATA.md` §0 melebur *"nilai + mata uang + kurs"* menjadi **satu atribut**, dan peleburan
itu **konseptual**. Ketika `2-to-spec/ddl-usulan/` pertama dibangkitkan, akibatnya terlihat: setiap
paket uang menjadi **satu kolom angka tanpa denominasi**, dan **INV-08, INV-09, INV-12 tidak dapat
dikompilasi** — ketiganya menyebut *"kode mata uang"* pada kunci alaminya.

### 8.1 Tiga golongan, dibaca dari kolom *Asal* §10

Keenam belas atribut `U1 paket uang` **tidak satu bentuk**:

| Gol. | Jumlah | Bentuk di sistem lama | Bentuk di skema baru |
|---|---:|---|---|
| **A** | 5 | daftar per mata uang — `MDPList[]`, `MDPMinList[]`, `GrossPremiumList[]`, `GrossPremiumMinList[]`, `ReserveList[]` | **tabel anak**, satu baris per mata uang |
| **B** | 5 | `Amount` + `Currency` pada satu baris, dan **induknya sendiri sudah per mata uang** | kolom **`KODE_MATA_UANG`** pada baris itu |
| **C** | 6 | **skalar polos** — `MaxCoGroup`, `MaxCoNonGroup`, `OptionLimit`, `Deductible`, `AgregateLimit`, `NILAI_BATAS` | **DITAHAN** — menunggu `T-6` |

### 8.2 INV-08, INV-09, INV-12 — DIBEBASKAN

Golongan B menjawabnya sendiri: mata uang memang **atribut baris induknya** di sistem lama, bukan
sesuatu yang hilang. Ketiganya kini tertulis sebagai `UNIQUE` di `Z00_KUNCI_ALAMI.sql`:

| Invarian | Constraint |
|---|---|
| INV-08 | `UNIQUE (ID_VERSI_KONTRAK, ID_KELOMPOK_TREATY, KODE_MATA_UANG)` pada `RETENSI_CEDANT` |
| INV-09 | `UNIQUE (ID_VERSI_KONTRAK, ID_KELOMPOK_TREATY, KODE_MATA_UANG)` pada `EGNPI` |
| INV-12 | `UNIQUE (ID_VERSI_KONTRAK, NOMOR_TERMIN, KODE_MATA_UANG)` pada `TERMIN` |

> **Catatan yang tidak boleh hilang.** Menulis ketiganya **tanpa** mata uang bukan jalan tengah — ia
> mengubah artinya menjadi *"dilarang dua baris bermata uang berbeda"*, dan itu **menolak data yang
> sah**. Itu bentuk kekeliruan lingkup yang sudah tiga kali terjadi di modul ini (§2.2c `CONTEXT.md`).

### 8.3 Lima tabel anak baru, dan kunci alaminya BELUM bernomor

Golongan A melahirkan `NILAI_MDP`, `NILAI_MDP_MINIMUM`, `NILAI_PREMI_BRUTO`,
`NILAI_PREMI_BRUTO_MINIMUM`, `NILAI_CADANGAN_PREMI` — masing-masing berkunci alami **kode mata uang
di dalam induknya**, **belum bernomor**.

| | |
|---|---|
| **Siapa menutup** | sesi to-spec berikutnya — penomorannya sepele, tetapi ia **keputusan**, bukan pembukuan |
| **Yang menagih** | `Z00_KUNCI_ALAMI.sql` memuat 23 constraint yang **seluruhnya menyebut nomor invariannya**; kelima tabel ini tidak muncul di sana, dan ketiadaannya terlihat |

**Bentuknya merosot dengan aman:** bila ternyata satu induk selalu satu mata uang, yang tersisa
hanyalah tabel anak berisi satu baris — penalaran yang sama dengan §10.9 `POTONGAN`.

### 8.4 Golongan C — ditahan, dan kenapa itu bukan kelalaian

Keenamnya **tidak punya mata uang di sumbernya sama sekali**. Menetapkan "mata uang kontrak" untuk
mereka berarti **menuliskan fakta yang sistem lama tidak pernah catat** ke setiap baris migrasi.
Pertanyaannya **`T-6`**, dan sampai ia dijawab keenam kolom itu berdiri tanpa denominasi — dinyatakan
sebagai pernyataan keputusan di dalam berkas DDL-nya masing-masing, bukan sebagai `TODO`.

---

## 9. `INV-69` dan `INV-70` — diterapkan 24 September 2026 (diff `D-4`)

Dari **`GRL-20`**, ronde D grilling Adjustment: materialitas adalah **atribut masukan**, dipilih
**sebelum** perubahan dikerjakan, dan ia **sakelar dua arah**.

### 9.1 Bunyinya, dan kenapa dua arah

`EVIDENCED(Section@ekspor-2026-09)` — hitungannya
`../treaty-in-adjustment/GRILL-D/01-TEMUAN.md` `TD-02`.

> **Angka *"220 kondisi `pyDisabledWhen` di 18 seksi"* DILARANG dipakai** (`MA-13`): ia tidak dapat
> direproduksi. Rujuk hitungan `GRILL-D`.

| Pilihan | Yang **tidak boleh** berubah |
|---|---|
| **`MATERIAL`** | teks kesepakatan — pengecualian, ketentuan khusus, nama kontrak, nomor rujukan kontrak |
| **`TIDAK_MATERIAL`** | besaran uang dan porsi — **428 sel bernama** di sistem lama |

Rumusan lama *"Non Material berarti uang tidak berubah"* hanya menyebut **satu arah**, dan **tidak
dipakai lagi di mana pun**.

### 9.2 Dua uji yang dilewati keduanya

| | `INV-69` | `INV-70` |
|---|---|---|
| **Dapat gagal?** | ya — versi `TIDAK_MATERIAL` dengan satu baris selisih premi | ya — versi `MATERIAL` dengan `PENGECUALIAN` berubah |
| **Dapat dikompilasi jadi nama kolom?** | ya — `NILAI_SELISIH` dikelompokkan `ID_VERSI_KONTRAK`, disaring jenis besarannya, dibandingkan `SIFAT_MATERIAL_ADDENDUM` | ya — kolom teks `VERSI_KONTRAK` dibandingkan terhadap `ID_VERSI_KONTRAK_DASAR` |

### 9.3 KADAR — **KLAIM, bukan penegakan**, dan itu dinyatakan di sini

> **`INV-69` dan `INV-70` berstatus KLAIM sampai dijalankan di instans nyata.**

Keduanya **CONSTRAINT lewat *materialized view***, dan MV punya sifat yang constraint biasa tidak
punya: **ia dapat berhenti menjaga tanpa satu pun kegagalan terlihat.** Bila MV-nya basi, gagal
*refresh*, atau dinonaktifkan, keduanya **tidak menolak apa pun dan tidak menimbulkan galat**.

Dan keduanya **tidak dapat diuji sekarang** — tidak ada instans Oracle yang terjangkau (`L-3`).

> Ini penerapan aturan yang sudah berlaku di berkas ini untuk `INV-47`, `INV-50`, `INV-51`, dan
> bentuk lintas baris `INV-31`: **penegakan yang dapat berhenti tanpa gagal menuntut uji negatif
> yang benar-benar dijalankan dan pemantau kebasian — bukan argumen rancangan.**

**`INV-71` tidak kena syarat ini.** Ia `UNIQUE` biasa, terpasang langsung di DDL, dan **tidak dapat
berhenti diam-diam**.

### 9.4 Pemantau kebasian — bentuknya ditulis, bukan diasumsikan

| | |
|---|---|
| **Apa yang diperiksa** | untuk MV penopang `INV-69` dan `INV-70`: **waktu *refresh* terakhir**, **status *refresh* terakhir** (berhasil / gagal), dan apakah MV masih ber-`REFRESH ON COMMIT` |
| **Seberapa sering** | **setiap hari**, dan **sebelum setiap peralihan gelombang** — dua saat, bukan satu: harian menangkap pembusukan, pra-peralihan menangkap yang dimatikan saat persiapan |
| **Apa yang terjadi bila MV basi** | **bukan peringatan yang dapat dilewati.** Ia dicatat sebagai **kegagalan penegakan bernama**, menyebut invarian mana yang berhenti dijaga dan **sejak kapan** — sebab yang perlu diketahui bukan "ada yang salah", melainkan **berapa lama data masuk tanpa dijaga** |
| **Siapa membacanya** | pemilik proses, bersama daftar invarian berstatus KLAIM |

> **Mematikan MV adalah keputusan yang harus diambil dengan nama** — siapa yang memutuskan, kapan,
> dan invarian mana yang dilepaskan. Bukan tindakan operasional yang diambil diam-diam saat laporan
> performa datang. Bila dilepaskan, keduanya **turun ke aplikasi**, dan
> `PETA-INVARIAN-KE-DDL.md` diperbarui — ia tidak boleh menjadi invarian yang tidak ditegakkan di
> mana pun.

### 9.5 Arah dampak bila salah

**`INV-69`/`INV-70` salah** — materialitas ternyata tidak mengunci apa pun: keduanya dicabut, dan
kolomnya tinggal sebagai keterangan. **Murah.**

**Tidak dipasang lalu dibutuhkan** — data lama dan baru bercampur tanpa pembeda, dan **tidak ada cara
mengetahui mana yang dulu dinyatakan tidak material**. Tidak setangkup.

**Tidak satu pun ditandai terverifikasi.**
