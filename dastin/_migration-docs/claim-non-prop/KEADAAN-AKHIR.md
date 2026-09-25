# Keadaan akhir — lapisan data Claim Non Prop

**Tanggal**: 19 September 2026
**Untuk**: siapa pun yang melanjutkan pekerjaan ini tanpa pernah ikut di dalamnya.

Berkas ini menjawab empat pertanyaan yang akan Anda tanyakan lebih dulu: **sudah sampai mana**, **apa yang masih belum diketahui**, **bagaimana cara memasangnya**, dan **apa yang harus terjadi supaya pekerjaan ini berlanjut**.

Ia tidak menjelaskan sistem lama dan tidak mengulang keputusan. Untuk itu ada `BLUEPRINT.md` (apa yang terbaca dari sistem lama), `SPEC-MODEL-DATA.md` (bentuk sistem baru), `docs/adr/` (kenapa bentuknya begitu), dan `RIWAYAT-SESI-2026-09-19.md` (urutan keputusan dan sebab premis yang gugur).

---

## 1. Apa yang ada di sini

Modul **Claim Non Prop** di Nusantara Re berjalan di atas Pega, dan akan dipindahkan ke React + Golang + Oracle. Pekerjaan yang sudah selesai adalah **lapisan datanya saja** — skema Oracle, tabel, view, hak akses. Tidak ada satu baris Golang, React, endpoint, layar, service, repository, maupun ORM; itu disengaja dan dinyatakan di setiap tiket.

Seluruh perilaku yang dirujuk di sini dibaca dari **ekspor XML 279 berkas** (`D:\XML_NURE\Claim Non Prop`, bertanggal 2026-09-08/09) dan dari **49 berkas DDL** sistem lama (`ls pengetahuan/ddl/*.sql | wc -l`). Bila sebuah dokumen berbeda dari XML, **XML yang menang untuk perilaku**.

> **Satu aturan yang tidak boleh Anda lewatkan**: setiap pernyataan di seluruh folder ini berlabel — EVIDENCED, EVIDENCED-NIHIL, DECIDED, DECIDED-TEKNIS, DERIVED, EXTERNAL, atau **TIDAK DITEMUKAN**. Yang tidak ada di sumber ditulis **TIDAK DITEMUKAN**, bukan ditebak. Bila Anda menambahkan sesuatu, ia juga berlabel.

---

## 2. Empat pencacah, beserta perintahnya

Angka tanpa perintah tidak dihitung sebagai jawaban. Keempatnya dapat Anda jalankan sekarang.

### 2.1 Pertanyaan — **aktif 0 · selesai 57 · dihapus 2**

```
cd _migration-docs/claim-non-prop
( sed -n '/^# SELESAI$/,/^# Lampiran/p' _selesai/OPEN-QUESTIONS.md \
    | grep -oE '^\| ~*[ABCEGK][0-9]+[ab]?~*' | tr -d '|~ ' \
    | grep -vE '^(A14|A10b|A11b)$' ; printf 'F\nH\nI\n' ) | sort -u | wc -l
```

Menghasilkan **57**, dan penyesuaiannya ada **di dalam perintah**, bukan di prosa sesudahnya:

| Langkah | Kenapa |
|---|---|
| kumpulkan kode berhuruf di bagian SELESAI | 57 baris mentah |
| **buang `A14`** | catatan perpindahan nomor, bukan pertanyaan |
| **buang `A10b`, `A11b`** | **dihapus** dari register, bukan ditutup — jawabannya tidak mengubah satu kolom pun |
| **tambahkan `F`, `H`, `I`** | tiga hipotesis sejajar; berkode huruf saja, tanpa angka, jadi pola di atas tidak menangkapnya |

> **Peringatan yang perlu dibaca**: keluaran mentah sebelum penyesuaian **juga 57**, dan kesamaan itu **kebetulan** — komposisinya berbeda (tiga dibuang, tiga ditambah). Jangan memakai bentuk perintah yang lebih pendek lalu mengira ia membuktikan angkanya.

### 2.2 REQ — **aktif 35 · selesai 1 · mati 1**

```
grep -cE '^\| REQ-0[0-9]+ \|.*\| OPEN \|' ORACLE-REQUESTS.md
```

Kepala berkasnya memuat angka berjalannya. **Tidak satu pun menahan keputusan rancangan**; empat menahan **pelaksanaan** tiket tertentu — REQ-018, REQ-033, REQ-021, REQ-037.

### 2.3 Tiket — **aktif 0 · menunggu instance 8 · tertahan 1 · selesai 29 · mati 1 = 39**

```
cd .scratch/claim-non-prop-lapisan-data/issues
grep -h '^status:' *.md _selesai/*.md _tertahan/*.md _mati/*.md | sort | uniq -c
ls *.md _selesai/*.md _tertahan/*.md _mati/*.md | grep -v README | wc -l
```

Kedua perintah harus sepakat di **39**. Status adalah **field di dalam berkas**, bukan lokasi folder — foldernya hanya kerapian, dan pembangkit papan (`alat/papan.py`) membaca field-nya.

`alat/papan.py` **menulis satu berkas dan hanya satu**: `issues/README.md`. Berkas tiketnya sendiri tidak pernah disentuh.

### 2.4 Pengenal DDL — **355 unik · nol di atas 30 byte · nol singkatan**

```
python alat/periksa-penamaan.py
```

Alat itu **hanya membaca** — ia tidak pernah mengubah apa pun, dan itu disengaja: pemeriksa yang ikut mengubah berkas saat dijalankan adalah jebakan persis bagi orang yang cuma ingin memeriksa. Versi pertamanya begitu, dan itu diperbaiki.

Ia memeriksa **tiga** hal dan **keluar dengan kode 1** bila salah satunya dilanggar: nama di atas 30 byte, sisa singkatan dari daftar yang dibatalkan, dan nama constraint yang dipakai dua kali. Itu aturan 4 penamaan sebagai kode, bukan sebagai niat.

Keluarannya mencacah **delapan jenis** pengenal — tabel, view, sequence, constraint, index, akun/peran/tablespace, kebijakan pengawasan, dan kolom. Versi pertamanya hanya dua, dan itu sebabnya angka `316` dan `333` yang pernah ditulis di `RIWAYAT` **dicabut**.

---

## 3. Tiga batas yang berdiri — dan kenapa ketiganya bukan lubang yang terlupa

**Register pertanyaan yang kosong bukan pengetahuan yang lengkap.** Yang tertutup adalah pertanyaan yang **menahan pekerjaan**, bukan pertanyaan yang ada. Tiga hal tetap tidak diketahui, dan ketiganya **tercatat sebagai batas**:

### 3.1 Empat hipotesis permanen, delapan cabang

**Satuannya hipotesis, bukan cabang, dan angkanya empat.** Versi pertama berkas ini menulis "tujuh", lalu "empat", lalu "ketiganya" untuk hal yang sama — aturan 1 dilanggar di berkas yang menetapkannya. Dicabut; yang berlaku angka di bawah.

| Hipotesis | Tentang | Cabangnya |
|---|---|---|
| **F** | `.IsEditClaim` | F1 pembaca ada di modul Komite · F2 tidak ada pembaca sama sekali |
| **H** | `.ValueAdjustment` | H1 sudah IDR di hulu · H2 mata uang aslinya |
| **I** | transisi `.AcceptanceStatus` | I1 ditulis modul Komite · I2 tidak pernah ditulis |
| **G1** | syariah | G1a lini usaha nyata · G1b hanya pemisahan teknis |

```
grep -cE '^\| \*\*[FHI]\*\* |^\| ~~G1~~ ' _selesai/OPEN-QUESTIONS.md
```

Menghasilkan **4** — tiga baris di tabel penutupan F/H/I, satu baris G1 di tabel penutupan bagian G.

**Empat hipotesis, delapan cabang, nol dipilih.** Cabangnya tertulis utuh di lampiran `_selesai/OPEN-QUESTIONS.md`.

**Keempatnya dapat ditutup sebagai hipotesis justru karena rancangannya benar di kedua cabang.** Bila Anda menemukan diri sedang merancang sesuatu yang **hanya benar bila salah satu cabang benar**, berhenti — itu tanda Anda keluar dari batas ini.

### 3.2 Tiga puluh lima permintaan belum dijawab

Tiga puluh empat di `PERMINTAAN-DBA-2-PROFIL-DAN-STRUKTUR.md`, satu (REQ-032) di `PERMINTAAN-DBA-1-VERSI-INSTANCE.md`. Keduanya siap dikirim; tidak ada yang menahan rancangan.

### 3.3 Folder `Komite Claim Non Prop` tidak pernah dibuka

**Tertutup untuk pekerjaan ini — bukan ditunda.** Apa pun yang hidup di sana tidak pernah masuk ke folder ini, dan itu **diketahui, bukan diabaikan**. Batas kepemilikan ditetapkan ADR-0026: mengikuti nama class.

Akibat yang paling sering terlupa: **26 penanda dimigrasi apa adanya dan ditandai tak berpemilik** — daftarnya di butir **E17**, `grep -c '^| ~~E17~~' _selesai/OPEN-QUESTIONS.md` menunjuk barisnya. Jangan memberi mereka perilaku; perilakunya tidak terbaca dari sini.

---

## 4. Urutan pemasangan

**36 berkas di `ddl-usulan/`. Urutan abjad namanya sudah urutan pemasangannya** — tidak ada langkah yang perlu diingat di luar mengurutkan nama.

```
ls ddl-usulan/*.sql | wc -l
```

| Urutan | Berkas | Isinya | Dijalankan |
|---|---|---|---|
| 1 | `00_SKEMA_DAN_AKUN.sql` | 2 tablespace, 3 akun, 3 peran, hak sistem | DBA |
| 2 | `01_KLAIM.sql` … `22_ARSIP_MUATAN_KELUAR.sql` | **22 tabel**, 115 constraint, 11 index, **22 sequence** — *kelima angka dari keluaran `alat/periksa-penamaan.py`, §2.4* | pemilik skema |
| 3 | `V00_HAK_AKSES.sql` | 44 pencabutan hak objek | pemilik skema |
| 4 | `V01_…` … `V09_…` | **9 view** beserta hibah `SELECT` ke hilir | pemilik skema |
| 5 | `Z00_ISIAN_AWAL.sql` | **baris awal** — 1 tutup buku, 3 tarif | pemilik skema |
| 6 | `Z01_PENGAWASAN_TULIS.sql` | kebijakan pengawasan tulis | **DBA** — pemilik skema sengaja tidak diberi `AUDIT_ADMIN` |
| 7 | `Z02_KUNCI_PEMILIK.sql` | `ALTER USER KLAIMNP ACCOUNT LOCK` | DBA, **paling akhir** |

### Yang harus Anda tahu sebelum menjalankan apa pun

> **Tidak satu pun dari 36 berkas ini pernah dijalankan di mana pun.** Seluruhnya berlabel *USULAN — untuk dibaca, bukan untuk dijalankan*. Mereka ditulis supaya **dapat** dijalankan, dan itu berbeda dari sudah terbukti jalan.

**Tiga hal sengaja dikosongkan** di `00_SKEMA_DAN_AKUN.sql`, ditulis sebagai penanda `&&`: letak dan ukuran datafile, kata sandi ketiga akun, dan profil kata sandi. **Kata sandi tidak pernah ditulis di berkas mana pun**, termasuk berkas itu.

**Satu berkas berisi DML**, dan hanya satu: `Z00_ISIAN_AWAL.sql`. Diperiksa dengan:

```
grep -lE '^(INSERT|UPDATE|DELETE|MERGE)' ddl-usulan/*.sql
```

yang harus mengembalikan **berkas itu saja**. Ke-35 berkas lain menyatakan *"TIDAK ADA DML DI BERKAS INI"*, dan pernyataan itu dapat diperiksa.

**`Z01` dan `Z02` dijalankan DBA, bukan pemilik skema.** Itu bukan kelemahan rancangan; itu bentuknya. Pengawasan yang dapat dimatikan oleh yang diawasi bukan pengawasan.

### Jangan jalankan `Z02` terlalu cepat

`Z02_KUNCI_PEMILIK.sql` mengunci akun pemilik di langkah terakhir, dan **sesudah itu tidak ada lagi yang dapat menjalankan perbaikan DDL tanpa DBA membukanya kembali**. Itu memang bentuk yang diinginkan — bukan sesuatu yang perlu dihindari.

Tetapi ketiga puluh enam berkas ini **belum pernah dijalankan di mana pun**. Pemasangan pertama hampir pasti menemukan sesuatu, dan menemukannya **sesudah** akun terkunci berarti **setiap perbaikan kecil menjadi permintaan ke DBA**.

Urutan yang menjaga Anda:

1. Pasang di **instance uji** lebih dulu — seluruh 36 berkas **kecuali `Z02`**.
2. Jalankan **kedelapan tiket uji** di sana (`24` `32` `33` `34` `36` `37` `38` `39`).
3. Jalankan `Z02` **hanya setelah seluruhnya hijau**.
4. Di instance **produksi**, urutan yang sama berlaku: `Z02` paling akhir, dan **hanya sekali tidak ada lagi yang perlu diperbaiki**.

Ini tidak mengubah rancangan. Ia menjaga orang yang memasang.

---

## 5. Apa yang sudah dijanjikan basis data — dan apa yang tidak

Bagian ini yang paling mudah salah dibaca, jadi ia ditulis dua arah.

### Yang ditegakkan basis data

Keunikan kunci alami tiap entitas · urutan masa berlaku treaty · domain keadaan baris · **nilai rupiah tidak dapat lahir tanpa kurs** · **kurs tidak dapat ada tanpa asal-usulnya** · **suntingan tidak dapat ada tanpa pelakunya** · keharusan pelaku pembuat di setiap baris · arsip muatan keluar **tulis-sekali**, ditegakkan hak akses.

### Yang **tidak** ditegakkan, dan jatuh ke aplikasi

Ketetapan `TANGGAL_KEJADIAN` sesudah baris dibuat · domain `STATUS_KLAIM`, `LINI_USAHA`, `KEPUTUSAN_KOMITE`, `SEBAB_DITOLAK` — **seluruhnya sengaja terbuka**, karena domain tertutup yang belum diketahui akan menolak nilai yang sah · keutuhan rujukan polimorfik `TABEL_TUJUAN`/`ID_TUJUAN` · `UNIQUE` atas kombinasi pengenal lama (menunggu REQ-018).

### Yang tidak dapat dijanjikan sama sekali

> **"Satu pintu tulis" tanpa syarat tidak dapat dipertahankan.** Yang dapat: *tidak ada tulisan dari luar pintu yang tidak meninggalkan jejak.*

`V00_HAK_AKSES.sql` menutup **pintu objek**. Enam jalur lain melewatinya — hak sistem berakhiran `ANY`, prosedur definer's rights, hibah ke `PUBLIC`, `CREATE ANY TRIGGER`, peran yang memuat hak `ANY`, dan `GRANT ANY PRIVILEGE`. Menutupnya di tingkat instance, oleh DBA: **REQ-037**. Yang menutupi sisanya sementara ini: `Z01_PENGAWASAN_TULIS.sql`.

Jalur kedua **bukan kemungkinan teoretis** — ia pola yang sudah dipakai di instance ini: `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP` menulis atas nama pemiliknya.

---

## 6. Apa yang harus terjadi supaya pekerjaan ini berlanjut

Tidak ada pekerjaan yang tertunda. Yang tersisa adalah pekerjaan yang **membutuhkan sesuatu yang belum ada**, dan ketiganya di luar kendali folder ini.

### 6.1 Instance Oracle yang dapat dipasangi → membuka **delapan tiket uji**

`24` `32` `33` `34` `36` `37` `38` `39`, seluruhnya berstatus `menunggu-instance`. **Rancangannya lengkap; yang kurang mesinnya.** Uji ini baru dapat hijau sesudah DDL benar-benar dijalankan.

Ditandai begitu supaya papan **tidak menunjukkan pekerjaan yang tidak dapat dimulai siapa pun** — bukan supaya ia tampak sedikit.

### 6.2 Batch lapisan aplikasi dimulai → membuka **tiket `11`**

Tipe desimal di sisi Golang. Prasyaratnya **tinggal satu**: batch aplikasi dimulai. Ia tertahan **lingkup, bukan pengetahuan** — presisinya sudah tetap, `NUMBER(38,20)`, kelompok AK-1 di ADR-0003.

### 6.3 REQ-012 dijawab → membuka **pemetaan migrasi kolom demi kolom**

Satu-satunya keluaran `SPEC-MODEL-DATA.md` yang belum selesai.

**Ia sengaja tidak ditulis**, dan alasannya perlu Anda pegang: pemetaan yang sah hanya datang dari `Data-Admin-DB-Table`. Menuliskannya sekarang berarti **mengabadikan kecocokan nama sebagai kebenaran** — dan itu persis cara `M_CURRENCYSTANDARD` dan `CURRENCYSTANDARD` hampir tertukar (D19: memiliki view tidak menutup permintaan atas tabelnya).

### 6.4 Empat REQ yang menahan pelaksanaan

| REQ | Menahan | Yang berubah bila jawabannya lain |
|---|---|---|
| **REQ-018** | penanganan migrasi di `12`, `15`, `18`, `29`, `34` | `UNIQUE` atas kombinasi pengenal lama — **constraint, bukan kolom** |
| **REQ-033** | pemetaan view kompatibilitas — `29`, `32` | berapa banyak kelompok yang ditolak, bukan apakah penolakannya ada |
| **REQ-021** | `35`, `39` | perlu-tidaknya pokok tersendiri untuk manajemen |
| **REQ-037** | klaim ADR-0017 di `35` | apakah "satu pintu tulis" berlaku sungguhan di instance tujuan |

---

## 7. Enam aturan yang berlaku untuk pekerjaan berikutnya

Keenamnya lahir dari cacat nyata di pekerjaan ini, bukan dari selera.

1. **Setiap angka membawa perintah yang menghasilkannya.** Angka yang tidak dapat direproduksi dicabut, bukan diwariskan. Empat angka dicabut dengan cara ini pada 19 September: "aktif 41", "selesai 26", "368 kolom" — dan, sesudah berkas ini ditulis, **"452 kolom"** juga, ketika rekonsiliasi tabel datar menemukan bahwa 452 mencacah kolom tabel **dan** kolom view sekaligus.
2. **Setiap kesimpulan lama yang gugur dicabut di tempatnya**, bukan hanya dilaporkan di tempat lain. Yang dicabut lebih mudah hilang daripada yang ditambahkan.
3. **Sekali sebuah kelas cacat ketahuan, ia disapu ke seluruh berkas** — bukan diperbaiki di satu tempat. Aturan ini ada karena dilanggar: `REVOKE` yang gagal di instance bersih diperbaiki di satu berkas pada pagi hari, dan terulang utuh di berkas lain pada sore hari.
4. **Register menyimpan keadaan; laporan menyimpan bukti.** Sebuah sapuan belum selesai sampai setiap butir register yang disentuhnya diperbarui **di tempatnya**.
5. **Setiap pola sapuan mencakup empat bentuk** — berkutip ganda (termasuk `&quot;`), berkutip tunggal, **tanpa kutip**, dan **nilai yang dirakit**. Dua kali sapuan literal melewatkan sesuatu karena hanya mencakup bentuk pertama (D9, C9).
6. **Berkas turunan tidak pernah menjadi sumber.** Bila turunan berbeda dari sumbernya, **alatnya yang salah**. Aturan ini ada karena dua sumber untuk satu kebenaran adalah cacat yang folder ini sudah tutup dua kali — `OS_AKSEPTASI_KLAIM` menyimpan nilai yang sama di `DATA_JSON` dan di 59 kolom skalar, dan angka `25` punya dua sumber di sistem lama. Yang turunan di sini: `PENGETAHUAN.md`, `ERD-ORACLE.xlsx`, dan `issues/README.md`.

Dan satu yang berlaku pada kriteria tiket, bukan pada kode:

> **Kriteria yang ternyata salah diubah, bukan diberi catatan penyimpangan.** Kriteria yang dibiarkan salah akan dibaca orang berikutnya sebagai pekerjaan yang belum selesai — dan ia akan "menyelesaikannya", melemahkan rancangan demi memenuhi kalimat.

---

## 8. Di mana mencari apa

| Pertanyaan Anda | Berkas |
|---|---|
| Sistem lama bekerja bagaimana? | `BLUEPRINT.md` · `PENGETAHUAN.md` (cuplikan bertanggal 18 Sep — **berkas aslinya yang berlaku**) |
| Bentuk sistem barunya apa? | `SPEC-MODEL-DATA.md` · `KAMUS-KOLOM.md` (**396 kolom tabel** ditambah 56 kolom view, tiap satunya berkolom sumber. *Angka **452 kolom di 22 tabel** dicabut 19 September 2026 oleh rekonsiliasi 1 tabel datar: 452 = **396 kolom tabel + 56 kolom view**, dan perintah yang tertulis di sebelahnya menghasilkan **457**, bukan 452. Ketiga angka beserta perintahnya di `sampah/keluaran/RINGKASAN-TABEL-DATAR.md` bagian 2.1.*) |
| Tabelnya berhubungan bagaimana? | **`ERD-ORACLE.xlsx`** — ERD kotak-entitas dalam Excel, **TABEL Oracle saja**, dua sistem di sheet terpisah. Sheet `ERD-BARU` + `RELASI-BARU`: 22 kotak, 396 kolom, **17 foreign key yang benar-benar tertulis di DDL**. Sheet `ERD-LAMA` + `RELASI-LAMA`: 32 tabel, 627 kolom, **nol foreign key** — 15 relasinya **tersirat**, dibaca dari comma-join dan subquery di SQL ekspor XML. Tidak ada garis yang disimpulkan dari kesamaan nama. Keterangannya di `ERD-ORACLE.md` |
| Kenapa bentuknya begitu? | `docs/adr/` — 29 ADR |
| Apa yang salah di sistem lama? | `FINDING-001` … `FINDING-008` |
| Apa yang belum diketahui? | `_selesai/OPEN-QUESTIONS.md` — register **tertutup**, lampirannya memuat cabang yang tidak dipilih |
| Apa yang diminta ke DBA? | `PERMINTAAN-DBA-1-…` dan `PERMINTAAN-DBA-2-…` |
| Bagaimana sampai ke sini? | `RIWAYAT-SESI-2026-09-19.md` |
| Apa yang dikerjakan berikutnya? | papan `.scratch/claim-non-prop-lapisan-data/issues/README.md` |
| Rancangan tabelnya bagaimana, dalam penamaan `T_*`? | **`ERD-CLAIM-NON-PROP.html`** dan **`Diagram-Skema-Tabel-ClaimNonProp.xlsx`** — 33 tabel, 30 relasi, **hanya PK dan FK**. Sejak **20-09-2026** awalan `T_CLAIMNP_` **tidak dipakai lagi**: `T_CLAIM_`, mengikuti [keputusan work owner] di `Diagram-Skema-Tabel-NusantaraRe.xlsx`. Akibatnya 11 tabel **dipakai bersama** Claim Prop dan Claim Fac In. Peta nama lama → baru: `4-erd-dan-tabel-datar/PETA-NAMA-TABEL.md` |
| Bagaimana Claim Non Prop menyambung ke lini lain? | **`ERD-GABUNGAN.html`** dan **`Diagram-Skema-Tabel-Gabungan.xlsx`** — 51 tabel, 44 relasi pohon, **11 tali penghubung**, lima lajur. **Enam tabel dipakai kedua berkas** dan digambar SEKALI saja di lajur INTI; lajur lain menunjuknya lewat penambat. Sheet *Tali Penghubung* memuat titik temunya satu per satu. Direkonsiliasi: 36/36 relasi `Daftar Relasi` NusantaraRe tertutup |
| Bentuk datar seluruh skema? | `ERD-ORACLE.xlsx`, sheet `KOLOM-LAMA` (627 kolom) · `KOLOM-BARU` (396) · `RELASI-LAMA` (15) · `RELASI-BARU` (17) · `CONSTRAINT-BARU` (76) — **TURUNAN**, dibangkitkan `alat/buat-erd-excel.py`. Tiap sheet sudah berpenyaring, tinggal dibuka. Bentuk TSV sebelumnya ada di `sampah/keluaran/` dan **menunggu dihapus**. Bila workbook berbeda dari sumbernya, **alatnya yang salah** |
| Berkas X ada di mana sekarang? | `PETA-FOLDER.md` — folder dirapikan **20-09-2026** jadi `1-grilling/` · `2-to-spec/` · `3-to-tickets/` · `4-erd-dan-tabel-datar/`. **Tidak satu nama berkas pun berubah**, hanya tempatnya. Keempat alat di `alat/` sudah disambungkan ulang dan dijalankan ulang: rekonsiliasinya sama persis dengan sebelum pindah |
| Apa isi folder papan itu? | `.scratch/claim-non-prop-lapisan-data/ISI-FOLDER.md` — indeks 39 tiket, aturan papan, dan sebab tiap tiket belum selesai |

**Berkas mana yang paling segar** terbaca di `STATUS.md`, yang menandai tiap artefak SEGAR, BERUMUR, PERMANEN, atau **TURUNAN** beserta tanggalnya.
