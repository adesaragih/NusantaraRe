# Penyelarasan spec Treaty In dengan aplikasi — empat keputusan

Spec Treaty In disusun di `D:\XML_NURE\_migration-docs\treaty-in\` sebelum modul ini punya tempat di
`APP_RNM`. Di empat titik, bentuk yang spec usulkan **tidak dapat berdiri apa adanya** di aplikasi ini.
Berkas ini mencatat keempatnya beserta alasan dan syarat pembalikannya — **bukan menyamarkannya**.

Dibuat 1 Oktober 2026, bersamaan migrasi `400`–`402` (tiket `14` dan `15`).

---

## 1 · Skema `{skema}`, bukan `TREATY_MASUK` — ADR-0028 DIBALIK

| | |
| --- | --- |
| **Spec** | `2-to-spec/ddl-usulan/00_SKEMA_DAN_AKUN.sql` membuat skema `TREATY_MASUK`, akun `TREATY_MASUK_APP`, dan `GRANT` per tabel. ADR-0028 memilihnya sebab `POOLDATA.TREATY_IN` dan `POOLDATA.TREATY_IN_EDM` warisan **sudah ada** di instance yang sama. |
| **Aplikasi** | Seluruh modul bermigrasi ke **satu** skema, lewat penanda `{skema}` yang diganti nama dari `ORACLE_SCHEMA` saat dijalankan (ADR-U-0033, `inti/backend/migrasi/migrasi.go`). Nol `CREATE USER` dan nol `GRANT` di seluruh repo. |
| **Diputuskan** | Ikut aplikasi. Migrasi memakai `{skema}`; `CREATE USER` dan `GRANT` **tidak dibawa**. |
| **Kenapa** | Sebab ADR-0028 hilang dasarnya di sini: skema sasaran bukan lagi `POOLDATA`, melainkan skema yang DBA tunjuk lewat `ORACLE_SCHEMA`, sehingga tabrakan nama yang ADR-0028 hindari tidak terjadi. Mempertahankan `TREATY_MASUK` menuntut dukungan multi-skema di `inti/backend/migrasi` — berkas milik tim inti, dan perubahan yang **seluruh** modul tanggung ongkosnya. |
| **Akibat** | Tabel Treaty In berdiri berdampingan dengan tabel modul lain di satu skema. Nama tabelnya (butir 3) karena itu menanggung beban pembeda sendirian. |
| **Pembalikan** | Bila tim inti kelak mendukung multi-skema, `{skema}` diganti dan berkas `00` dihidupkan kembali. Tidak ada data yang perlu dipindahkan selama pembalikan terjadi **sebelum** tiket `44` memuat data pertama. |

## 2 · `NUMBER(38,20)` → `NUMBER(38,8)` dan `NUMBER(9)` → `NUMBER(10)`

| | |
| --- | --- |
| **Spec** | Butir `KTV-A` (24 September 2026) menetapkan `NUMBER(38,20)` untuk uang, persen, tarif, kurs, dan limit; `NUMBER(9)` untuk cacah dan nomor urut. Ditulis **TANPA VERIFIKASI**, dengan kaidah *"terlalu lebar di Oracle murah, terlalu sempit memotong data"*. |
| **Aplikasi** | `TestNolNumberTanpaPresisi` (`inti/backend/penjaga/strukturkolom_test.go`) menolak **setiap** bentuk `NUMBER` di luar empat ini: `NUMBER(38,8)` — *"uang, share, persen, dan rate, keputusan work owner c, 26 September 2026"* · `NUMBER(5)` · `NUMBER(10)` · `NUMBER(19)`. Penjaga berjalan atas SQL **seluruh** modul, setiap pull request. |
| **Diputuskan** | Uang, persen, tarif → `NUMBER(38,8)`. Cacah dan nomor urut → `NUMBER(10)`. Kunci asing dan pengenal → `NUMBER(19)`, tidak berubah. |
| **Kenapa** | **Dua perubahan berbeda, dan wewenangnya tidak sama — dicatat terpisah supaya tidak terbaca seolah satu klausul membenarkan keduanya.** `NUMBER(38,20)` → `NUMBER(38,8)` adalah **penyempitan**, dan `KTV-A` menulis syarat pembalikannya sendiri: *"sesi DDL boleh MEMPERSEMPIT, dan HANYA SEBELUM DATA DIMUAT."* Kita tepat di titik itu — tiket `44`, yang memuat data pertama, belum berjalan. `NUMBER(9)` → `NUMBER(10)` adalah **pelebaran**, yang klausul itu **tidak** bicarakan sama sekali; satu-satunya dasarnya adalah `presisiSah` penjaga, yang tidak mengenal `NUMBER(9)`. Keduanya bertentangan dengan tuntutan tiket `14` bahwa tipe **mengikat** pada `KAMUS-KOLOM.md`, dan keputusan work owner 26 September 2026 — **lebih baru** daripada `KTV-A`, dan berlaku lintas aplikasi — yang mengalahkannya. |
| **Akibat** | **Angka uang kehilangan 12 angka di belakang koma** (20 → 8). Delapan angka desimal adalah presisi yang modul Claim Life dan PremiumList Life pakai untuk uang, share, persen, dan rate yang sama. `NUMBER(9)` → `NUMBER(10)` adalah **pelebaran**, bukan penyempitan — tidak ada risiko memotong. |
| **Pembalikan** | Menambah `NUMBER(38,20)` ke `presisiSah` adalah suntingan berkas **tim inti**, dan ia melonggarkan penjaga untuk **seluruh** modul. Bila teknik treaty membuktikan 8 desimal memotong angka sungguhan, itu jalurnya — lewat pull request tim inti, dan **sebelum tiket `44`**. |

## 3 · Nama tabel memakai nama spec apa adanya

| | |
| --- | --- |
| **Spec** | `KONTRAK`, `VERSI_KONTRAK`, `MATA_UANG`, `BAHAYA`, `KELAS_BISNIS`, `KELOMPOK_TREATY`, `JENIS_POTONGAN`, `JENIS_REASURANSI` — nama telanjang, sebab spec mengandaikan skema sendiri (butir 1). |
| **Aplikasi** | Kedua puluh tabel yang sudah ada berawalan `T_` (`T_WORK_CLAIM`, `T_PREMIUM_LIST`, `T_CLAIMLF_*`). Tidak ada penjaga yang mewajibkannya. |
| **Diputuskan** | Nama spec dipakai apa adanya `[keputusan work owner 01-10-2026]`. |
| **Akibat** | Nol pemetaan antara `KAMUS-KOLOM.md` dan DDL — nama di spec = nama di basis data, dan setiap rujukan silang di 51 tiket tetap sah tanpa diterjemahkan. Ongkosnya: nama sejenerik `MATA_UANG` dan `BAHAYA` kini **milik Treaty In** di skema bersama, dan modul lain yang membutuhkan tabel bernama sama harus memakai yang ini atau memilih nama lain. Sequence **tidak** ikut aturan ini: ia berawalan `SEQ_TRIN_` sebab nama seperti `SEQ_KONTRAK` tidak menyebutkan pemiliknya sama sekali. |

## 4 · `ddl-usulan/` tidak memuat satu pun `CREATE SEQUENCE`

| | |
| --- | --- |
| **Temuan** | Sapuan atas seluruh 36 berkas `2-to-spec/ddl-usulan/` menemukan **nol** `CREATE SEQUENCE`. Padahal `INV-02` menuntut pengenal datang dari sequence, dan daftar periksa tiket `14` menuntutnya dengan kalimat sendiri: *"pengenal keduanya datang dari `SEQUENCE` tanpa `CYCLE`; tidak ada jalur lain yang dapat memberi pengenal."* |
| **Diputuskan** | Sequence ditulis di berkas **tersendiri**, `402_sequences.sql` — bukan ditambalkan ke dalam berkas tabel. |
| **Kenapa** | Supaya ketiadaannya di `ddl-usulan/` tetap terbaca sebagai **temuan atas spec**, bukan hilang ke dalam berkas yang orang kira salinan spec. |
| **Ditagih** | Pemilik `alat/buat-kamus-dan-ddl.py` — pembangkit DDL-nya tidak pernah membangkitkan sequence, sehingga **ke-35 entitas** lain akan lahir dengan lubang yang sama. |

## 5 · Penyimpangan lain dari `ddl-usulan/`, didaftar supaya tidak ada yang tersembunyi

Keempat butir di atas adalah keputusan BESAR. Yang di bawah lebih kecil, tetapi tetap penyimpangan —
dan daftar yang mengaku "empat penyelarasan" tanpa menyebutnya adalah daftar yang berbohong.

| Yang ditambahkan | Tidak ada di `ddl-usulan/` | Kenapa |
| --- | --- | --- |
| `IX_VERSI_KONTRAK_KONTRAK` (`VERSI_KONTRAK.ID_KONTRAK`) | benar | setiap pembacaan rantai versi menyaring menurut kolom ini; Oracle **tidak** membuat index untuk kunci asing dengan sendirinya |
| `IX_VERSI_KONTRAK_DASAR` | benar | idem, dipasang migrasi `440` modul `treatyinadjustment` |
| `IX_NILAI_MDP_LAYER` · `IX_NILAI_MDP_MIN_LAYER` · `IX_PEMULIHAN_LIMIT_LAYER` · `IX_JEJAK_PERUBAHAN_VERSI` | benar | keempat tabel ini tidak punya kunci alami, jadi tidak ada index UNIQUE yang melayani kunci asingnya |
| `IX_MATA_UANG_KONTRAK_MU` · `IX_RETENSI_CEDANT_KLP` · `IX_EGNPI_KELOMPOK` · `IX_EGNPI_KELAS_BISNIS` · `IX_BATAS_PER_BAHAYA_BHY` | benar | kunci asing ke **tabel acuan**; tidak pernah memimpin kunci alami. Menghapus satu baris `BAHAYA` atau `MATA_UANG` tanpa index ini memindai seluruh tabel anak |
| `FK_VERSI_KONTRAK_MATA_UANG` (`KODE_MATA_UANG_KONTRAK` → `MATA_UANG`) | benar | **INV-44** menuntutnya, dan tiket `15` menulis jalur gagalnya: *"Kode mata uang pada kontrak yang tidak ada di tabel acuan -> ditolak `INV-44`."* `ddl-usulan/09` menghilangkannya — lubang yang sama bentuknya dengan sequence yang hilang (butir 4) |
| `ID_VERSI_KONTRAK_DASAR` | ada di `ddl-usulan`? **tidak**, dan tidak ada di `KAMUS-KOLOM.md` §10.2 | dituntut tiket `01` papan Adjustment. **Utang hulu**: §10.2 perlu diperbarui, dan itu milik pemilik `SPEC-MODEL-DATA.md` |

⛔ **Satu kunci asing yang SENGAJA belum dipasang:** `KELAS_BISNIS_KONTRAK` → `KELAS_BISNIS`. Bentuknya
sama dengan mata uang (NUMBER(19) yang menunjuk tabel acuan), tetapi **tidak ada invarian yang
menamainya** — INV-44 menyebut mata uang saja. Memasang kunci asing atas tebakan berarti menolak data
yang mungkin sah. **Ditagih:** pemilik `SPEC-INVARIAN.md`, satu kalimat.

---

## Yang TIDAK diubah, supaya tidak ada yang mengira semuanya dinegosiasikan

- **Nama dan tipe kolom** mengikat pada `2-to-spec/KAMUS-KOLOM.md` (urutan wewenang butir 5), kecuali
  lebar `NUMBER` di butir 2. `VARCHAR2(1000 CHAR)`, `VARCHAR2(4000 CHAR)`, `VARCHAR2(40 CHAR)`, dan
  `DATE` dibawa persis.
- **`NOT NULL` spec dipertahankan.** Claim Life memakai kaidah *"seluruh kolom nullable kecuali kunci
  utama, wajib-isi di services"* (ADR-U-0027); kaidah itu **milik modul itu**, bukan penjaga lintas
  aplikasi, dan `NOT NULL` di sini adalah invarian yang spec putuskan sadar.
- **Nol `CHECK` daftar nilai** untuk keenam himpunan acuan (INV-62) dan untuk `KEADAAN_SIKLUS_HIDUP`
  (tiket `14` bab "Tidak termasuk"). Keduanya pernyataan keputusan yang tertulis di dalam DDL-nya.
- **Nol trigger, nol procedure, nol `COMMIT`** di teks SQL (ADR-0056 K-4 dan ADR-U-0029 sepakat).

---

## 6 · `L-3` dicabut, dan `ddl-usulan/` kini punya penilai

Keempat penyimpangan di atas ditulis saat `L-3` berlaku — *"tidak ada instans Oracle yang
terjangkau"* — sehingga seluruhnya dinilai dari BENTUK teksnya saja.

**Sejak 2 Oktober 2026 itu tidak lagi benar.** Ke-22 langkah migrasi kedua modul sudah dijalankan
Oracle dan diterima tanpa satu `ORA-` pun, dan sembilan constraint sudah dibuktikan sungguh menolak.
Rinciannya di [`VERIFIKASI-ORACLE-2026-10-02.md`](VERIFIKASI-ORACLE-2026-10-02.md).

**Akibatnya pada keempat penyimpangan:**

| Butir | Keadaan sesudah diuji |
| --- | --- |
| 1 · skema `{skema}` | terbukti jalan — tetapi ia jalan di **POOLDATA**, dan itu masalah tersendiri (§1 berkas verifikasi) |
| 2 · `NUMBER(38,8)` dan `NUMBER(10)` | **diterima Oracle**; penyempitannya tidak menolak satu pun nilai uji |
| 3 · nama tabel telanjang | **tabrakan nyata**: ke-29 nama kini berdiri di POOLDATA bersama 816 tabel warisan. Belum bertabrakan dengan nama warisan, tetapi ruangnya kini dibagi sungguhan — bukan lagi hipotetis |
| 4 · sequence ditulis sendiri | ke-29 sequence berdiri; **belum satu pun terpakai**, sebab belum ada jalur tulis |

⚠️ **Butir 3 naik dari risiko menjadi keadaan.** Keputusan memakai nama telanjang diambil saat
skema sasaran dianggap milik aplikasi ini sendiri. Ia ternyata POOLDATA. Bila DBA kelak menyediakan
skema tersendiri, keputusan itu kembali aman; bila tidak, ia perlu ditinjau ulang oleh pemilik
proses — dan peninjauannya harus terjadi **sebelum tiket `44` memuat data pertama**.

## 7 · Penyimpangan kelima: uji `-tags=db` dan target `test-db`

| | |
| --- | --- |
| **Apa** | `modul/treatyin/backend/repository/invarian_db_test.go` ditambahkan, dan `Makefile` target `test-db` kini menjalankan kedua modul ini di samping Claim Life. |
| **Kenapa** | Papan menuntut *"uji negatif DAN positif"*, dan `migrasi_invarian_test.go` secara terbuka menyatakan ia hanya membaca bentuk teks. Tanpa berkas ini, bukti perilaku yang dikumpulkan 2 Oktober hanya hidup di dalam sebuah dokumen — dan dokumen tidak pernah merah. |
| **Menyentuh berkas bersama** | ya, `Makefile` — tiga baris, milik tim inti. Perubahannya hanya **menambah** dua baris `go test`; nol target yang ada diubah. |
| **Pembalikan** | hapus kedua baris itu. Uji di dalam modul tetap melewati sendiri bila `ORACLE_DSN` kosong. |

---

## 8 · Pencacah pemanggil `skemauji.Buka()` di modul Claim Life

> ⚠️ **UTANG RONDE LALU.** Ronde 2 menyunting berkas ini dan **mengaku mencatatnya**. Tidak
> tercatat di mana pun. Ini pencatatannya, terlambat satu ronde.

| | |
| --- | --- |
| **Apa** | `modul/claimlife/backend/repository/batasanpemakaian_test.go`, konstanta `mau` dinaikkan `15` → `16`, beserta satu paragraf komentar yang menyebut pemanggil barunya. |
| **Kenapa** | `TestSetiapPemanggilBukaMemeriksaBolehDilewati` mencacah pemanggil `skemauji.Buka()` di seluruh repo dan **menolak** bila cacahnya berubah tanpa disertai alasan. `invarian_db_test.go` Treaty In menambah satu pemanggil. Ujinya sendiri yang meminta: *"bila memang bertambah, perbarui angkanya di sini."* |
| **Menyentuh berkas bersama** | **Ya** — berkas itu ada di folder `claimlife`, tetapi berfungsi sebagai sensus lintas modul. Tiga modul lain sudah menambah paragrafnya di sana lebih dulu (`treatycontractout`, `premiumlistlife`, `uji/lintasmodul`), jadi bentuknya mengikuti preseden, bukan menciptakan jalan baru. |
| **Pembalikan** | hapus paragraf dan turunkan `mau` kembali ke `15`, bersamaan dengan membuang `invarian_db_test.go`. |

## 9 · `ERD.md` §2 dipakai sebagai sumber perilaku hapus — dan 21 relasi diperbaiki

| | |
| --- | --- |
| **Apa** | Dua puluh satu kunci asing memperoleh `ON DELETE CASCADE`, satu memperoleh `ON DELETE SET NULL` beserta kunci asing yang sebelumnya tidak ada sama sekali (`FK_KONTRAK_DISALIN_DARI`), dan sebelas dinyatakan `tolak` secara eksplisit. Bab **"Kaskade ON DELETE CASCADE"** ditulis di `MODUL.md` kedua modul. |
| **Kenapa** | Tujuh belas berkas migrasi mengutip `INV-18` **terbalik**: *"TANPA ON DELETE: bawaan Oracle MENOLAK"*. `SPEC-INVARIAN.md` baris 129 menuntut perilaku hapus **"ditetapkan sadar, tidak dibiarkan bawaan"** — membiarkan bawaan lalu menyebut nama invarian yang melarangnya bukan pemenuhan, bahkan ketika bawaannya kebetulan cocok. Keputusan sadarnya sudah ada dan mengikat di `4-erd-dan-tabel-datar/ERD.md` §2, dan **dokumen itu tidak pernah dibuka** sampai 2 Oktober 2026. |
| **Bentuk perbaikan: SUNTING DI TEMPAT, bukan migrasi korektif `420`+** | Tiga alasan. **Pertama**, Oracle tidak dapat meng-`ALTER` aturan hapus sebuah kunci asing — ia harus di-`DROP` lalu dibuat ulang, sehingga migrasi korektif berarti 22 pasang `DROP`/`ADD` yang mengabadikan versi pertama yang salah di setiap pemasangan baru selamanya. **Kedua**, ke-29 tabel di POOLDATA **kosong** — nol data hilang. **Ketiga**, tabel itu memang salah tempat dan sudah dijadwalkan dibuang (`alat/buang-tabel-treatyin-dari-pooldata.sql`), sehingga sejarah `T_MIGRASI` di sana ikut terbuang. |
| **Akibat** | Berkas migrasi kini **menyimpang dari yang tercatat di `T_MIGRASI` POOLDATA**. Itu disengaja dan tidak berbahaya selama POOLDATA dibersihkan sebelum ada pemasangan lain yang membacanya. **Bila POOLDATA tidak jadi dibersihkan, keputusan ini harus dibalik** menjadi migrasi korektif `420`+. |
| **Pembalikan** | kembalikan klausa `ON DELETE`, dan pindahkan perubahannya ke migrasi `420`+ berisi `DROP CONSTRAINT` + `ADD CONSTRAINT`. |

## 10 · Pagar pada jalur `-migrate` maju

| | |
| --- | --- |
| **Apa** | `cmd/api/main.go` — `jalankanMigrasi` kini menerima `config.Config` dan menolak dua hal: `IS_PEGA_PROD=true`, dan `ORACLE_SCHEMA` yang memuat nama skema warisan. `inti/backend/config/config.go` memperoleh satu fungsi baru, `SkemaWarisan(string) bool`. |
| **Kenapa** | Jalur maju **tanpa pagar sama sekali**, sementara `-migrate-down` menolak produksi DAN menolak skema bukan-uji. Asimetri itulah yang menaruh 29 tabel di POOLDATA pada 2 Oktober 2026 — dan `-migrate-down` kemudian menolak membongkarnya. |
| **Sengaja LEBIH LONGGAR daripada `-migrate-down`** | Pagar maju **tidak** menuntut `ORACLE_SKEMA_UJI=true`. Migrasi maju memang dijalankan ke skema sungguhan saat pemasangan; menuntut pengakuan *"skema ini boleh dihapus isinya"* di sana akan salah. Yang ditolak hanya produksi dan skema warisan. |
| **Menyentuh berkas bersama** | **Ya** — `cmd/api/main.go` dan `inti/backend/config/config.go`, keduanya milik tim inti. `SkemaWarisan` **mengekstrak** pemeriksaan yang sudah ada di dalam `pagarSkemaUji` (MEMUAT, bukan sama persis) tanpa mengubah perilakunya; nol pagar yang ada dilonggarkan. |
| **Pembalikan** | kembalikan tanda tangan `jalankanMigrasi(svc *inti.Dasar)` dan hapus `SkemaWarisan`. Pembalikan itu mengembalikan lubang yang menaruh 29 tabel di POOLDATA. |

---

## 11 · RALAT atas §9 dan §10 — POOLDATA adalah skema sasaran

**Keputusan pemilik proses, 2 Oktober 2026:** POOLDATA **tidak** dibersihkan dan **tidak** diganti.
Ke-29 tabel Treaty In dibuat di sana dengan sengaja, dan di sanalah ia tinggal.

Dua keputusan sebelumnya bersandar pada andaian sebaliknya, dan keduanya dibalik di sini.

### 11.1 §9 dibalik — suntingan di tempat menjadi migrasi korektif `420`

§9 menulis syarat pembalikannya sendiri: *"Bila POOLDATA tidak jadi dibersihkan, keputusan ini
harus dibalik menjadi migrasi korektif `420`+."* Syarat itu **kepicu**.

| | |
| --- | --- |
| **Apa** | Migrasi `400`–`419` dan `440` dikembalikan ke bentuk yang tercatat di `T_MIGRASI` — nol klausa `ON DELETE`. Ke-21 kaskade dan satu `SET NULL` pindah ke `420_perilaku_hapus_erd.sql` sebagai `DROP CONSTRAINT` + `ADD CONSTRAINT`. |
| **Kenapa sepasang** | Oracle tidak dapat meng-`ALTER` aturan hapus sebuah kunci asing yang sudah berdiri. |
| **Satu yang bukan pasangan** | `FK_KONTRAK_DISALIN_DARI` **belum pernah ada** — kolomnya berdiri sejak `401` tanpa kunci asing. Ia **dibuat** di `420`, nol `DROP` mendahuluinya. |
| **Yang TIDAK pindah** | Ralat `INV-18` di 18 berkas **tetap di tempatnya**. Yang pindah klausa DDL-nya, bukan penjelasannya. |
| **Akibat pada uji** | `TestPerilakuHapusSesuaiERD` tetap mengadu ke `ERD.md` §2 dengan tabel yang sama; **pembacanya** yang disesuaikan — ia kini mengambil definisi **terakhir** di urutan migrasi, sebab `420` menimpa `403`–`418`. `gabungan()` juga diurutkan: peta Go tidak berurut, dan urutan menentukan siapa menimpa siapa. |
| **Pembalikan** | `420_perilaku_hapus_erd_down.sql` mengembalikan seluruhnya ke bentuk pra-`420`. Ia **mengembalikan pelanggaran INV-18**, dan kepalanya mengatakan itu. |

### 11.2 §10 dibalik SEBAGIAN — pagar POOLDATA dicabut dari jalur maju

§10 mencatat pemasangan pagar yang menolak skema warisan di jalur `-migrate` **maju**. Dengan
POOLDATA sebagai skema sasaran, pagar itu **mematikan `make migrate` sama sekali**.

| | |
| --- | --- |
| **Apa dicabut** | panggilan penolakan skema warisan di `cmd/api/main.go`, dan fungsi `config.SkemaWarisan()` yang §10 tambahkan — ia kehilangan seluruh pemanggilnya, dan fungsi terekspor tanpa pemanggil adalah kode mati. |
| **Apa TETAP** | penolakan `IS_PEGA_PROD=true` di jalur maju. `PagarSkemaUji`/`PastikanSkemaUji` **utuh** — keduanya menjaga `-migrate-down` dan `uji/skemauji`, yang **MENGHAPUS** tabel. `NamaSkemaWarisan` utuh — ia dipakai `uji/skemauji`. |
| **Asimetri yang disengaja** | migrasi **maju** ke POOLDATA dikehendaki; `DROP`, teardown, dan `-migrate-down` dari sana **tidak**. POOLDATA memuat 816 tabel warisan sungguhan. **Jangan ratakan asimetri ini** — ia bukan ketidakkonsistenan, ia keputusan. |
| **Yang hilang bersamanya** | lubang yang menaruh 29 tabel di POOLDATA tanpa peringatan tetap terbuka — tetapi kini itu **perilaku yang dikehendaki**, bukan kecelakaan. |
| **Pembalikan** | pasang kembali pemeriksaan skema warisan di jalur maju, bersamaan dengan keputusan memindahkan tabel Treaty In keluar dari POOLDATA. |

### 11.3 Skrip pembersihan POOLDATA — DIBUANG

`modul/treatyin/alat/buang-tabel-treatyin-dari-pooldata.sql` **dihapus**. Ia tidak akan pernah
dijalankan, dan membiarkannya menunggu di folder `alat/` membuatnya terbaca sebagai rencana yang
masih berlaku. Isinya tercatat di sini: 29 `DROP TABLE`, 29 `DROP SEQUENCE`, dan satu `DELETE` atas
`T_MIGRASI`. Bila kelak dibutuhkan lagi, ia dibangkitkan ulang dari daftar tabel di
`STRUKTUR-TABEL-TREATY-IN.md`.

---

## 12 · `MASTERID` BUKAN kunci asing — tabel pendaratan tab Treaty In

**Keputusan, 3 Oktober 2026.** Migrasi `430` membuat delapan tabel `M_TREATYIN_*` yang
mendaratkan larik di dalam `POOLDATA.M_TREATY_IN.JSONDATA`. Rancangannya menyebut `MASTERID`
sebagai **kunci asing ke `TREATY_IN.ID`**. Kunci asing itu **tidak dibuat**, dan sebabnya bukan
selera.

### 12.1 Kenapa — diukur, bukan diperkirakan

Tiga kueri katalog di POOLDATA, 3 Oktober 2026:

| Kueri | Hasil |
| --- | --- |
| `all_constraints` untuk `TREATY_IN`, `constraint_type IN ('P','U')` | **nol baris** |
| `all_indexes` untuk `TREATY_IN` | satu, `INDEX_ID (ID)`, **NONUNIQUE** |
| `all_tab_columns` untuk `TREATY_IN.ID` | `VARCHAR2(100)`, **`NULLABLE = Y`** |

Oracle menolak `REFERENCES {skema}.TREATY_IN (ID)` dengan **ORA-02270** selama kolom rujukannya
tidak memimpin kunci utama maupun `UNIQUE`. Ini bukan pilihan rancangan yang dapat ditimbang —
pernyataannya **tidak dapat dijalankan**.

⚠️ Datanya sendiri bersih: ke-1.854 `ID` unik dan terisi, panjang maksimum 7, dan nol baris
`M_TREATY_IN` tanpa pasangan di `TREATY_IN`. **Yang hilang constraintnya, bukan integritasnya** —
dan karena itu memasang kunci asing akan berhasil andai `TREATY_IN` boleh disentuh.

### 12.2 Yang dipilih sebagai gantinya

| | |
| --- | --- |
| **Apa** | `MASTERID VARCHAR2(100 CHAR) NOT NULL` pada kedelapan tabel, memimpin `UNIQUE (MASTERID, URUTAN)` — yang sekaligus melayani indexnya. Nol klausa `REFERENCES`. |
| **Siapa yang menjaga keterhubungan** | pemuat (`services.MuatSatuKontrak`), yang hanya pernah menulis `MASTERID` yang baru saja ia baca dari `M_TREATY_IN.ID`, dan rekonsiliasinya, yang mengadu cacah baris dengan cacah elemen dokumen. |
| **Perilaku hapus dari kontrak** | `ikut hapus`, dijalankan `repository.KosongkanKontrak` — bukan oleh basis data. |
| **Satu kunci asing yang ADA** | `FK_MTI_INSTALLMENTITEM_1`, `M_TREATYIN_INSTALLMENTITEM.IDINDUK` → `M_TREATYIN_INSTALLMENT.ID`, `ON DELETE CASCADE`. Keduanya tabel baru, jadi di sana tidak ada penghalang. |
| **Preseden di modul ini** | `MIGRASI_KORELASI.ID_KONTRAK_BARU` dan `MIGRASI_PENDARATAN.KUNCI_WARISAN` — keduanya **nilai, bukan kunci asing**, dengan sebab yang berbeda (jejak asal-usul harus bertahan melewati penghapusan barisnya). Bentuknya sama; alasannya tidak, dan itu disebut supaya tidak terbaca sebagai satu aturan. |
| **Ongkos yang dibayar** | basis data **tidak** menolak `MASTERID` yang menunjuk kontrak yang tidak ada, dan **tidak** membersihkan baris pendaratan ketika barisan `TREATY_IN` dibuang oleh jalur lain. Keduanya nyata. |

### 12.3 Syarat pembalikan

Keputusan ini **dibalik** begitu salah satu dari dua hal terjadi:

1. **`POOLDATA.TREATY_IN` memperoleh kunci utama atau `UNIQUE` pada `ID`** — oleh siapa pun, atas
   sebab apa pun. Sejak saat itu `ORA-02270` berhenti berlaku, dan kedelapan `MASTERID` menjadi
   kunci asing lewat migrasi `ALTER TABLE ... ADD CONSTRAINT FK_MTI_<TAB>_MST`, perilaku hapus
   `ON DELETE CASCADE`, didaftarkan di `TestPerilakuHapusSesuaiERD`.
2. **Sumber pendaratan berpindah dari `TREATY_IN` ke `KONTRAK`/`VERSI_KONTRAK` model baru** —
   yang punya kunci utama sejak migrasi `401`. Pada titik itu `MASTERID` berganti menjadi
   `ID_VERSI_KONTRAK NUMBER(19)` dengan kunci asing biasa, dan kedelapan tabel berhenti menjadi
   tabel pendaratan.

⛔ Sampai salah satunya terjadi, **jangan "memperbaiki" ketiadaan kunci asing ini dengan menyentuh
`TREATY_IN`.** DDL terhadap tabel warisan dilarang, dan `TestPendaratanTidakMerujukTabelWarisanDenganKunciAsing`
menolak jalan pintasnya.

---

## 13 · Persen SHARE tampil **8 desimal** — bukan 2, bukan tanpa batas

**Keputusan pemilik proses, 4 Oktober 2026**, menjawab
[`PERTANYAAN-TERBUKA-PERSEN-SHARE.md`](PERTANYAAN-TERBUKA-PERSEN-SHARE.md).

| | |
| --- | --- |
| **Apa** | `labels.ts` memperoleh `DESIMAL_PERSEN_SHARE = 8`, dipakai cabang `case 'persenShare'` di `selAngka` menggantikan `DESIMAL_TAK_DIBATASI`. `DESIMAL_UANG = 4` dan `DESIMAL_PERSEN = 2` **tidak berubah**. |
| **Kena pada** | `CESSIONPCT` · `RNMSHARE` · `QSOR` · `QSRI` · `BROKERAGEPERCENTP` (dari `M_TREATY_IN2`) dan `PctLimit` (tab Co-Ins Scale) — persen yang beberapa barisnya **dijumlahkan dan harus menghasilkan 100**. |
| **TIDAK kena pada** | `MDP_RATIO` · `ADJ_RATE` · `ROL`, yang tetap 2 desimal: ketiganya tidak dijumlahkan menjadi 100, dan dua di antaranya memang melampaui 100 (109,6 dan 199,4 terukur). |
| **Kenapa bukan 2** | `SD-05`/`BR-01` benar — tiga share `33,333` berjumlah tepat 100, tiga share `33,33` tidak. |
| **Kenapa bukan tanpa batas** | `2,825601535925207120348922139444%` (30 desimal, nyata di satu kontrak) tidak terbaca di dalam sel grid. Dan penyimpanannya **hanya 8 desimal** — `NUMBER(38,8)`, dijaga `TestNolNumberTanpaPresisi`: menampilkan 30 berarti mengaku lebih teliti daripada yang sistem simpan. |
| **Kenapa 8** | sama persis dengan batas penyimpanan, dan untuk share seberapa pun realistis sifat jumlah-tepat-100 tetap terjaga. |
| **Pembalikan** | ganti satu angka di `labels.ts`. Nilai tersimpan tidak pernah diformat, jadi perubahan tampilan tidak menyentuh satu baris data pun — itulah sebab §13 murah dibalik dan tidak perlu syarat tambahan. |

### 13.1 ⚠️ Satu dari dua alasan penolakan 2 desimal TIDAK terpenuhi oleh 8 — diukur

Pertanyaannya menolak 2 desimal dengan **dua** alasan. Yang kedua (`SD-05`/`BR-01`) terpenuhi oleh 8.
Yang **pertama tidak**, dan itu diukur sesudah keputusan diterapkan:

| Batas desimal | `formatPersen('99.999999999999900', n)` |
| ---: | --- |
| 2 · 6 · **8** · 10 · 12 | `100%` |
| 13 ke atas, dan tanpa batas | `99,9999999999999%` |

Pembulatannya setengah-ke-atas dan limpahannya merambat, jadi `PctTotal` `99.999999999999900`
**tetap tampil `100%` pada 8 desimal** — persis hal yang alasan pertama ingin cegah. Selisihnya
baru terlihat pada **13 desimal**, yang melampaui presisi penyimpanan.

⛔ **Keputusan 8 tetap dijalankan apa adanya**, sebab alasan keduanya berdiri sendiri dan kuat.
Yang dicatat di sini: **jangan mengira `PctTotal` yang tampil `100%` sudah pasti tepat 100.**
Nilai yang tersimpan tetap utuh di tabel pendaratan (teks, nol tafsir), jadi pemeriksaan presisi
penuh dilakukan di sana, bukan di layar.

⚠️ **Ditagih balik ke pemilik proses:** bila selisih `PctTotal` memang harus terlihat di layar,
ia menuntut perlakuan tersendiri — bukan perubahan `DESIMAL_PERSEN_SHARE`, sebab 13 desimal akan
melanggar batas penyimpanan yang jadi alasan kedua. `TestPersenShareDipotongPadaDelapan`
mengunci perilaku yang berlaku hari ini, lengkap dengan kasus `100%`-nya.

### 13.2 Pita tidak diberi tanda `%` kedua

Diperbaiki bersama §13, sebab ia muncul saat mengujinya: `formatPersen` mengembalikan teks
bukan-angka apa adanya **lalu menempelkan `%`**, sehingga pita `>=30% up to < 50%` menjadi
`>=30% up to < 50%%`. `selAngka` kini memeriksa lebih dulu apakah nilainya angka murni.

⛔ Ini **bukan pemformat kedua** — `format.ts` tidak disentuh; yang berubah hanya **apakah** ia
dipanggil. Hari ini hanya `CoInShare` berbentuk pita dan ia digolongkan `teks`, jadi jalur ini
tidak pernah terpicu; penjaganya ada untuk salah-golong berikutnya, yang jaraknya satu huruf.

---

## 14 · Tab Retro **DITUNDA**, dan kedua kontraknya **DITANDAI**

**Keputusan pemilik proses, 4 Oktober 2026**, menjawab
[`PERTANYAAN-TERBUKA-RETRO.md`](PERTANYAAN-TERBUKA-RETRO.md).

| | |
| --- | --- |
| **Apa** | nol tabel, nol pemuat, nol layar untuk Retro. Tabnya **tetap** `.trin__belum` — yang kurang kodenya, bukan datanya. |
| **Kenapa** | **2 kontrak dari 1.854** tidak cukup untuk merancang tiga tabel bersarang. Model yang diturunkan dari dua contoh akan salah di tempat yang tidak ada contoh ketiga untuk membantahnya, dan ongkos salahnya mahal: tiga tabel, satu pemuat, satu rekonsiliasi, dan presisi yang tidak dapat dipersempit lagi sesudah data masuk. Delapan dari tujuh belas medannya turunan (`Total*`, `INV-58`), jadi isi nyatanya sekitar sembilan. |
| **Penandaannya** | `repository/warisan_retro.go` — `KontrakRetroTertunda = ["1000493", "1000755"]`, `LarikRetroTertunda`, dan `AlasanRetroTertunda`. Keduanya `NonProportional`, 2 elemen masing-masing. |
| **Pembalikan** | begitu kontrak **ketiga** memperoleh `RetroList`, dasar penundaan berubah dan keputusan ini ditinjau ulang — bukan daftarnya yang diperbarui. |

### 14.1 Penandaannya disapu MESIN, dalam tiga lapis

Catatan yang hanya hidup di dokumen dibaca oleh yang mencarinya. Yang diperlukan di sini adalah
sesuatu yang lewat di depan yang **tidak** mencarinya — orang yang menyimpulkan *"27.238 cocok,
selesai"*. Tiga lapis, dan ketiganya diperlukan:

| Lapis | Wujud | Kapan ia bicara |
| --- | --- | --- |
| **1 · daftar** | `repository.KontrakRetroTertunda` | saat seseorang `grep` pengenalnya |
| **2 · uji** | `TestKontrakRetroTertundaMasihDuaItu` (`-tags db`) — **mengukur ulang dari Oracle**, menyaring calon dengan `DBMS_LOB.INSTR` lalu **mengurai JSON** atas kunci puncak | tiap `make test-db-treatyin` |
| **3 · pemuat** | `cetakPenandaRetro()` di `backend/pemuat/jalankan.go`, tercetak pada **tiap `-cocokkan`** | saat rekonsiliasi dijalankan |

⛔ Lapis 2 sengaja **tidak** mempercayai penyaringan teks saja: kunci dapat muncul bersarang, dan
sapuan substring atas dokumen bersarang sudah pernah menipu ronde ini sekali — `EGNPI` terbaca
155 padahal 846. Yang memutuskan pengurai JSON, atas kunci puncak.

Lapis 2 juga menagih hal kedua: **nol tabel pendaratan boleh memuat `RetroList`.** Bila suatu
hari ada, penanda ini dan tabelnya saling membantah, dan ujinya merah.

### 14.2 ⛔ Kewajiban bila Retro kelak dibangun

**`ShareSumary` WAJIB diadili lebih dulu: turunan dari `Share`, atau bukan.** Bila turunan,
`INV-58` melarangnya dan yang tersisa **dua** tabel, bukan tiga. Kewajiban ini ditulis di sini
dan bukan ditinggalkan sebagai hal yang diingat orang.

*(Ejaan `ShareSumary` adalah ejaan ekspor apa adanya, bukan salah ketik berkas ini.)*

### 14.3 Kapan penandanya dihapus

**Bersama tabelnya, dalam ronde yang sama** — berkas `warisan_retro.go`, ujinya, dan panggilan
`cetakPenandaRetro()` sekaligus. Jangan lebih awal: tenggang waktu saat penandanya sudah hilang
tetapi datanya belum pindah adalah persis lubang yang penanda ini ada untuk mencegahnya.

---

## 15 · Tab teks lewat **Jalan B** — isinya tetap di `JSONDATA`

**Keputusan pemilik proses, 4 Oktober 2026**, menjawab
[`PERTANYAAN-TERBUKA-TAB-TEKS.md`](PERTANYAAN-TERBUKA-TAB-TEKS.md).

| | |
| --- | --- |
| **Apa** | **nol tabel pendaratan** untuk Exclusions dan Special Conditions. Kelima kuncinya dibaca dari `M_TREATY_IN.JSONDATA` lewat jalur warisan yang sudah ada (`jsonWarisan`), dipilih menurut cabang di `services/tab_teks.go`. |
| **Kenapa** | layarnya hari ini **baca-saja**, dan nol tiket meminta penyuntingan. Dokumennya sudah ditarik untuk medan lain, jadi satu medan `CLOB` lagi dari dokumen yang sama **nol tambahan perjalanan** ke basis data. Jalan A tetap menuntut `CLOB` — isinya mencapai 23.453 aksara, jauh di atas batas `VARCHAR2` 4.000 — sehingga keunggulan utamanya hilang sementara ongkosnya (satu migrasi, satu pemuat, satu rekonsiliasi) dibayar penuh. |
| **`ValueDifference` DIKELUARKAN** | ia **objek** berisi `EGNPI`/`Limits`/`Share` dan sembilan medan `Total*` — potret nilai sebelum perubahan, bukan teks yang pemakai ketik. Menaruhnya di tab teks salah dua kali: bentuknya bukan teks, artinya bukan medan layar. Ia milik pertanyaan lain — apakah riwayat nilai ikut dipindahkan sama sekali — wilayah tiket `06`/`11`/`13` (`NILAI_SELISIH`). Tab Value Difference **tetap** `.trin__belum`. |
| **Pembalikan** | **begitu ada tiket yang membuat tab itu dapat DISUNTING sebelum tiket `44` selesai, Jalan A terbayar** — menulis satu baris jauh lebih aman daripada menulis ulang dokumen 158 KB yang 174 medan lain menumpanginya. Datanya tetap di dokumen, jadi perpindahan ke A **tidak kehilangan apa pun**. |

### 15.1 ⛔ Ejaan dipilih menurut CABANG, dan tidak pernah jatuh ke ejaan lain

Terukur atas 1.854 dokumen: dari **303** dokumen yang punya lebih dari satu ejaan
`SpecialConditions*`, **nol** yang isinya identik. Ejaan lain **bukan salinan yang basi** —
ia teks yang berbeda.

Aturannya:

1. ejaan yang **sesuai cabang** dipakai — `…P` proporsional, telanjang non-proporsional;
2. bila ia **tidak ada**, hasilnya **kosong** — bukan ejaan lain;
3. ejaan lain yang berisi tetap **disebut di layar**, sebab pembacanya berhak tahu ada teks yang
   tidak ia lihat.

Sebarannya, terukur:

| | Exclusions | Special Conditions |
| --- | ---: | ---: |
| ejaan cabang ada — proporsional | 1.027 / 1.079 | 896 / 1.079 |
| ejaan cabang ada — non-proporsional | 719 / 772 | 509 / 772 |
| **ejaan LAIN juga berisi** | 175 + 126 = **301** | 157 + 93 = **250** |
| **ejaan cabang kosong, ejaan lain ada** | 2 + 2 = **4** | 20 + 33 = **53** |
| nol ejaan sama sekali | 50 + 51 = 101 | 163 + 230 = 393 |

⚠️ **Ejaan ketiga `SpecialConditionsp`** — huruf kecil di akhir, **bukan salah ketik**: 292
dokumen memakainya, **170 proporsional dan 122 non-proporsional**. Karena ia dipakai kedua
cabang, ia tidak terikat cabang mana pun dan **tidak pernah terpilih**; bila berisi, ia selalu
muncul sebagai "ejaan lain".

### 15.2 Cabang dibaca dari KOLOM, bukan dari dokumen

`SifatProporsional` membaca `TREATY_IN.PROPORTIONTYPE`. Terukur: kolomnya terisi pada seluruh
1.854 baris, sedangkan kunci `ProportionType` **di dalam dokumen TIDAK ADA pada tiga** —
`1001854`, `1001855`, `1001856`. Memilih ejaan dengan kunci yang kadang hilang berarti ketiga
kontrak itu diam-diam diperlakukan sebagai non-proporsional.

*(Ketiganya tidak punya satu pun dari kelima kunci teks, jadi tidak ada teks yang terkena hari
ini — tetapi andaiannya tetap salah, dan andaian yang salah tanpa akibat adalah andaian yang
akan berakibat nanti.)*

### 15.3 Tampilan

Teks **dapat digulir** (`.trin__teks`, `max-height: 60vh`, `white-space: pre-wrap`), bukan satu
baris yang terpotong: isinya mencapai 23.453 aksara dan memuat baris baru sungguhan. Kosong
memakai `Kosong` (*"belum ada DATA"*), **bukan** `.trin__belum` — datanya memang tidak ada,
kodenya ada.

## 16 · `MATA_UANG` dan `MATA_UANG_KONTRAK` **DICABUT** — kurs dari `TREATYEXCHANGEYEARLY`

| | |
| --- | --- |
| **Apa** | Migrasi `434` membuang tabel `MATA_UANG`, `MATA_UANG_KONTRAK`, kunci asing `FK_VERSI_KONTRAK_MATA_UANG`, dan kedua sequence-nya. Himpunan acuan turun dari **enam menjadi lima**: `jenis-potongan`, `kelas-bisnis`, `kelompok-treaty`, `bahaya`, `jenis-reasuransi`. |
| **Kenapa** | Keputusan pemilik proses 4 Oktober 2026: daftar mata uang dan kursnya diambil dari **`TREATYEXCHANGEYEARLY`** — tabel warisan yang sudah hidup, 140 baris, 25 mata uang, terbagi per `TREATYYEAR` — seperti yang layar lama lakukan. Dua tabel model baru itu tidak dipakai lagi. |
| **Aman, dan itu terukur** | Keduanya **nol baris** pada hari pencabutan, dan kunci asing yang masuk hanya dua — satu di antaranya ikut terbawa tabelnya sendiri. Nol data hilang. |
| **Akibat 1 — tiket 20 kehilangan tabelnya** | `MATA_UANG_KONTRAK` adalah tabel tiket `20`. Dua uji Oracle yang membuktikannya (`TestTiket20KursMataUangGandaDitolak` dan pasangan positifnya) dicabut bersamanya. Lapisan skema tiket `20` karena itu **mundur**, dan itu dinyatakan di sini alih-alih ditemukan orang lain di ledger. |
| **Akibat 2 — tiket 57 TIDAK PUNYA RUMAH** | Tiket `57` menuntut kurs **dibekukan** pada versi yang disetujui: *"ubah baris kurs tahunan lalu buka versi disetujui → angkanya tidak berubah"*. `TREATYEXCHANGEYEARLY` adalah justru baris kurs tahunan itu — master bersama yang boleh berubah — dan `MATA_UANG_KONTRAK.KURS` adalah tempat pembekuannya. Sesudah pencabutan ini **tidak ada tempat beku**. Siapa pun yang mengerjakan `57` harus memutuskan di mana lebih dulu. |
| **Akibat 3 — `INV-44` tidak lagi ditegakkan basis data** | Kolom `VERSI_KONTRAK.KODE_MATA_UANG_KONTRAK` **tetap ada** dan kini tanpa penjaga: ia dulu menunjuk `MATA_UANG.ID_MATA_UANG`. Mengubah artinya menjadi kode `TREATYEXCHANGEYEARLY` adalah pekerjaan tiket `20`, bukan migrasi ini. |
| **Satu penjaga yang ikut diperbaiki** | `TestPerilakuHapusSesuaiERD` menilai kunci asing dari teks `CREATE` di migrasi; `DROP` tidak terlihat olehnya, sehingga ia **lulus sambil menagih tiga kunci asing yang sudah tidak ada**. Pembacanya kini mengenal pencabutan, dan **posisi yang menentukan** — migrasi `420` membongkar lalu memasang ulang 21 kunci asing, dan membaca "ada DROP" saja akan membuang kedua puluh satunya. |
| **Pembalikan** | `434_cabut_mata_uang_down.sql` membangun keduanya kembali utuh beserta kedua kunci asingnya, bentuk disalin apa adanya dari `400`/`402`/`403`. Selama keduanya nol baris, pembalikan tidak kehilangan apa pun. Yang perlu ikut dibalik: lima himpunan acuan kembali enam (`models/acuan.go`, `repository/acuan.go`, `services.go`, `frontend/api.ts`, `labels.ts`), dan kedua uji tiket `20`. |
