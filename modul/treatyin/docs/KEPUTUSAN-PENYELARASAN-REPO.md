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

---

## 17 · Tab Retro **TIDAK DIBANGUN — karena JARANG**

**Keputusan pemilik proses, 4 Oktober 2026.** Menggantikan [§14](#14), yang tetap berdiri sebagai
riwayat.

| | |
| --- | --- |
| **Apa** | Retro tidak dibangun. Tabnya tetap menyatakan belum ada kode, dengan petunjuk yang menyebut **jarang**. Nol tabel, nol pemuat, nol layar. |
| **Kenapa** | `TreatyIn.IsMultipleRetro` bernilai `"true"` pada **5 dari 1.854** dokumen (`false` 1.531, kunci absen 318). `RetroList` berisi pada **2** kontrak, 4 elemen. Delapan dari tujuh belas medannya turunan (`Total*`, INV-58), jadi isi nyatanya sekitar sembilan. |
| **⛔ Kenapa BUKAN karena kode mati** | Ronde sebelumnya hendak mencabutnya atas dasar penjaga `1=2`, dan **verifikasi membatalkan dasar itu**. `Section/ShareRetro.xml` memang memuat tiga penjaga `1=2` — tetapi ketiganya membungkus **sebuah tombol** (@340.829), **satu blok `BAR`** (@1.442.638), dan **satu tombol kepala** (@1.507.551). **Nol yang membungkus tabnya.** Tab non-prop berdiri **tanpa syarat tampil sama sekali**; tab prop bersyarat `TreatyIn.IsMultipleRetro`, sebuah syarat DATA. `FlowAction/ShareRetro.xml` berbunyi `pyRuleAvailable = Yes`. **Retro HIDUP di Pega** — ia jarang, bukan mati. |
| **Pembalikan** | kontrak **keenam** ber-`IsMultipleRetro = true`, **atau** kontrak **ketiga** ber-`RetroList` berisi. Saat salah satunya terjadi, dasar keputusan ini berubah — **tinjau ulang keputusannya, jangan sekadar perbarui angkanya.** |

### 17.1 Kedua ambangnya DIJAGA, bukan dihafal

§17 berdiri di atas dua angka, jadi dua penjaga mengukurnya ulang dari Oracle:

| Penjaga | Mengukur | Ambang |
| --- | --- | ---: |
| `TestAmbangRetroJarangMasihDuaKontrak` | kontrak ber-`RetroList` berisi | **2** |
| `TestAmbangRetroMultipleMasihLima` | dokumen ber-`IsMultipleRetro = "true"` | **5** |

Keduanya **merah** pada hari ambangnya terlampaui, dan pesannya menyebut §17 beserta syarat
pembalikannya — bukan sekadar "angka berubah".

⚠️ `repository.KontrakRetroTertunda` **berganti nama** menjadi `KontrakRetroJarang` dan berganti
arti: dari *"penanda pekerjaan tertunda"* menjadi **ambang pembalikan §17**. Tidak ada pekerjaan
yang menunggu; yang ada ambang yang diawasi. `cetakPenandaRetro()` di pemuat idem.

---

## 18 · Aturan nama berkas berlaku **di pintu masuk**, tidak surut

**Keputusan pemilik proses, 4 Oktober 2026.**

| | |
| --- | --- |
| **Apa** | `services.NamaBerkasAman` menegakkan spanduk biru *"Recommended safe substitute should be . or _"* pada unggahan **baru**. Ke-43 baris yang sudah tersimpan **dibiarkan apa adanya** — tidak diganti namanya, tidak ditolak, **dan tidak ditandai di layar**. |
| **Kenapa dibiarkan** | **37 dari 43** nama melanggar aturannya (spasi, tanda kurung) — misalnya `Signed Schedule Marine Hull Quota Share 2025 (3).pdf`. Mengganti nama berarti **menulis ke tabel warisan** `M_ATTACHMENTTREATY_2`, dan `FILENAME` adalah penghubung ke `T_STORAGE_IMAGE`: menggantinya dapat **memutus tautan ke berkasnya**. |
| **Kenapa tidak ditandai** | penanda yang **tidak dapat ditindaklanjuti** hanya melatih orang mengabaikan penanda. Keputusannya membiarkan; menandai 37 baris yang memang tidak akan diapa-apakan membuat setiap penanda di layar ini kehilangan bobotnya. |
| **Pembalikan** | **ketika `FILENAME` berhenti menjadi penghubung ke penyimpanan** — yaitu ketika `T_STORAGE_ID` menjadi satu-satunya tautan — pembersihan nama menjadi mungkin tanpa risiko memutus berkas. Saat itu ke-37 nama dapat dirapikan dalam satu migrasi data. |

⚠️ `NamaBerkasAman` **tidak diubah** oleh keputusan ini: ia sudah benar, dan ia diuji dua arah —
yang ditolak dan yang diterima. Yang ditetapkan di sini **lingkup berlakunya**, bukan isinya.

---

## 19 · `NILAI_SELISIH` + `NILAI_SEBELUM_PRO_RATE` **PINDAH** ke modul Adjustment

**Keputusan pemilik proses, 4 Oktober 2026** — kepemilikan **dijernihkan**, bukan dikecualikan.

| | |
| --- | --- |
| **Apa** | `treatyin/.../435_cabut_nilai_selisih.sql` mencabut keduanya; `treatyinadjustment/.../442_nilai_selisih.sql` membangunnya kembali. **Bentuknya disalin apa adanya** — `NUMBER(38,8)`, `SEQ_TRIA_*`, `FK_NILAI_SELISIH_1`, `FK_NILAI_SEBELUM_PRO_RATE_1`, `IX_*`. Diverifikasi identik: `diff` atas kedua berkas tanpa komentar mengembalikan **nol selisih**. |
| **Kenapa** | tiket **76** dan **77** ada di papan Adjustment, dan tabelnya kini ada di sana juga. Selama terpisah, tiket `06`, `11`, `13` tertahan di belakangnya. |
| **Kenapa AMAN** | keduanya **nol baris**, diukur ulang langsung ke POOLDATA **sebelum** pencabutan ditulis dan **sesudah** pembangunan kembali. Nol baris berarti pemindahan tidak kehilangan apa pun. |
| **`426` TIDAK dihapus** | pola migrasi `434` dipakai: yang dicabut **bendanya**, bukan catatannya. Kepala `426` memuat empat alasan rancangan yang masih mengikat — kardinalitas 1:N lawan 1:1 ERD, penghalang `BESARAN_DAPAT_DISESUAIKAN`, induk `VERSI_KONTRAK` lawan `KONTRAK`, dan INV-58 — dan `442` **menunjuk balik** ke sana alih-alih menyalinnya dan membiarkan salinannya membeku. |
| **Pembalikan** | `435_..._down.sql` dan `442_..._down.sql`, dan keduanya **WAJIB berjalan bersama**. Dua modul yang sama-sama membuat tabel yang sama bertabrakan dengan ORA-00955, dan pemasangan berhenti di tengah. |

### 19.1 Ketergantungan urutan, dan apa yang menjaganya

Kedua tabel berkunci asing keluar ke `VERSI_KONTRAK` — **milik modul `treatyin`**, migrasi `401`.
Yang menjaga urutannya hanyalah **pilihan rentang nomor**: `treatyin` memakai `400`–`439`, modul
Adjustment `440`+, dan pelari migrasi berjalan menurut nomor.

⛔ **Tidak ada penjaga yang menegakkan ini.** Karena itu ia dinyatakan di kepala `442` dan di
sini: siapa pun yang mengusulkan menomori ulang rentang modul Adjustment memutus kedua kunci
asing itu.

### 19.2 Dua penjaga modul Adjustment menjadi daftar-IZIN, bukan dicabut

| Penjaga | Sebelum | Sesudah |
| --- | --- | --- |
| `TestNolTabelBaru` | nol `CREATE TABLE` di modul ini | **hanya** `NILAI_SELISIH` dan `NILAI_SEBELUM_PRO_RATE`; tabel ketiga tetap ditolak, dan daftar-izin yang tidak terpakai juga ditolak |
| `TestNolKaskadeDiModulIni` → `TestKaskadeHanyaPadaBerkasTerdaftar` | nol `ON DELETE` di modul ini | **hanya** berkas `442_`; berkas lain tetap ditolak |

⛔ **Melonggarkannya menjadi "modul ini boleh membuat tabel" akan membuang pagar yang masih
diperlukan** — model datanya memang sebagian besar masih satu dengan `treatyin`, dan percabangan
diam-diam adalah persis yang pagar itu cegah.

⚠️ Kedua baris kunci asing **IKUT PINDAH** ke `TestPerilakuHapusSesuaiERD` modul Adjustment,
**tidak digandakan**. Dua tempat yang menyatakan perilaku hapus yang sama akan berselisih, dan
yang salah satunya basi tidak akan berbunyi.

⚠️ Di sisi `treatyin`, kedua kunci asing berhenti dinilai dengan sendirinya: penyaring
`dicabut()` — yang lahir bersama migrasi `434` — mengeluarkan kunci asing yang tabelnya dibuang
migrasi berikutnya. Teks `CREATE`-nya tetap ada di `426`, sebab migrasi tidak pernah disunting
mundur.

### 19.3 ⛔ Yang DITAGIH, dan bukan milik ronde ini

`treaty-in/SPEC-MODEL-DATA.md` masih menyatakan model datanya **SATU**, dan pesan
`TestNolTabelBaru` yang lama mengutipnya apa adanya. Pernyataan itu **kini tidak lagi benar
seluruhnya** — dua tabel berdiri di luar `treatyin`.

Yang dapat memperbaikinya **pemilik spec**, bukan berkas migrasi dan bukan penjaga. Sampai
diperbaiki, §19 ini adalah satu-satunya tempat yang menyatakan selisihnya.

---

## 20 · Syarat tampil tab — keadaan KELIMA, dan ia BERCABANG

Ronde 5 Oktober 2026 membangun **syarat tampil**, bukan pemangkasan daftar. Tab yang syaratnya
tidak terpenuhi **tidak dirender sama sekali** — keadaan yang berbeda dari `Kosong` (tabnya ada,
isinya nol), dari `.trin__belum` (tabnya ada, kodenya belum ditulis), dan dari **jarang** (§17).

**Terukur**, dengan `<pyIncludedRuleXML>` dibuang lebih dulu:

| Berkas | wadah tab | bersyarat |
| --- | ---: | --- |
| `Section/TreatyInTabsProportional.xml` | 11 `BAR` (nol `TABBED`) | `Retro` @1.031.030 — `TreatyIn.IsMultipleRetro` |
| `Section/TreatyInTabsNonProportional.xml` | 11 `TABBED` | `Value Difference` @3.848.982 — `TreatyIn.EDMState != 3 && TreatyIn.EDMMaterialType == 1` |

Syarat `Value Difference` yang sama muncul di **tiga** ekspor — `TreatyInTabsNonProportional.xml`
@3.848.982, `InputTreatyInOffer.xml` @5.479.082, `TreatyInNONProportional.xml` @5.840.210.

### 20.1 ⛔ Syaratnya milik CABANG, bukan milik nama tab

`Retro` ada di **kedua** daftar, dan hanya yang proporsional bersyarat — wadah `Retro` non-prop
(TABBED @3.291.822) **nol** `pyContainerVisibleWhen`.

Bentuk pertama ronde ini memakai satu peta berkunci nama tab, dan ia menyembunyikan `Retro`
non-proporsional pada **770 dari 775** kontrak. Yang menangkapnya ujinya
(`syarat-tab.test.ts`), bukan pembacaan ulang. Petanya kini dua:
`SYARAT_TAB_PROPORSIONAL` dan `SYARAT_TAB_NON_PROPORSIONAL`.

### 20.2 ⚠️ Kedua syarat itu tidak pernah terpenuhi oleh data hari ini

| Syarat | Terpenuhi di POOLDATA |
| --- | ---: |
| `IsMultipleRetro = "true"` **dan** kontraknya proporsional | **0** dari 1.079 |
| `EDMMaterialType = 1` | **0** dari 1.854 |

Kelima kontrak `IsMultipleRetro = "true"` (`1000493`, `1000755`, `1001043`, `1001405`,
`1001853`) **seluruhnya non-proporsional** — dan di cabang itu tabnya tidak bersyarat.

⛔ **Itu BUKAN kesimpulan bahwa kedua tab mati.** Pega menilai syaratnya atas *clipboard*, yang
dapat diisi Activity saat jalan tanpa pernah tersimpan di `JSONDATA`. Yang terukur hanya bahwa
nilainya tidak **tersimpan**. Pertanyaannya di `PERTANYAAN-TERBUKA-LAYAR-PEGA.md` §4.

### 20.3 Kontrak baru memperlihatkan SELURUH tab

`tabUntuk(jenis)` tanpa argumen kedua mengembalikan daftar penuh. Kontrak yang belum punya
dokumen tidak punya jawaban syaratnya, dan menyembunyikan tab karena pertanyaannya belum dapat
dijawab menyembunyikannya **selamanya**: tab tersembunyi tidak pernah diisi, dan yang tidak
pernah diisi tidak pernah memenuhi syaratnya.

---

## 21 · `Accounting Mode` adalah DUA properti, dan `Bordereaux` satu cabang

**Yang memutuskan bukan jaraknya di berkas melainkan `pyCondition` sel masing-masing**, di
`Section/TreatyInNONProportional.xml` — satu-satunya Section yang memuat kepala kontrak, untuk
kedua cabang:

| Properti | `pyValue` | `pyCondition` sel |
| --- | ---: | --- |
| `TreatyIn.Bordeaux` | @56.974 | @60.649 `TreatyIn.ProportionType='Proportional'` |
| `TreatyIn.TreatyYear` | @102.100 | @105.380 `ALWAYS` |
| `TreatyIn.AccountingMode` | @108.019 | @111.777 `TreatyIn.ProportionType='Proportional'` |
| `TreatyIn.AccountingModeNonProp` | @114.500 | @118.211 `TreatyIn.ProportionType='NonProportional'` |

Jaraknya ajeg — tiap `pyCondition` ~3.700 bita sesudah `pyValue` selnya — dan keempatnya
mengikuti pola yang sama, jadi pemasangannya bukan kebetulan.

⚠️ **Dan itu aturan TAMPIL, bukan aturan data.** Sapuan 1.854 dokumen:

| Kunci | Proporsional (1.079) | Non-proporsional (775) |
| --- | --- | --- |
| `Bordeaux` | `nonreporting` 626 · `reporting` 453 | `reporting` 711 · `nonreporting` 61 · tidak ada 3 |
| `AccountingMode` | `underwriting` 1.046 · `accounting` 33 | `accounting` 708 · `underwriting` 64 · tidak ada 3 |
| `AccountingModeNonProp` | **`loss` 1.079** | `loss` 751 · `risk` 21 · tidak ada 3 |

⛔ Ketiganya ada di **kedua** cabang. Jadi "medan ini milik cabang itu" tidak dapat disimpulkan
dari isi dokumen; ia harus dibaca dari `pyCondition`, dan itulah yang dilakukan.

⚠️ **Angka briefing berbeda, dan yang dipakai yang terukur.** Briefing menyebut `AccountingMode`
`underwriting` 758 · `accounting` 442 dan `Bordeaux` `reporting` 634 · `nonreporting` 566;
totalnya terukur 1.110/741 dan 1.164/687 — angka yang sudah tercatat di `models` sejak 3 Oktober.

### 21.1 Keempat labelnya TIDAK dikarang

Kedua dropdown `pyListSource = associated` — daftar pilihannya milik `Rule-Obj-Property`, dan
korpus **tidak punya folder `Property`** (1.048 berkas XML disapu). Yang ditampilkan nilai
tersimpan apa adanya. Permintaan ekspornya di `PERTANYAAN-TERBUKA-LAYAR-PEGA.md` §1, lengkap
dengan contoh bentuk yang sudah ada: `EDMState.xml` memperlihatkan `pyPromptTableList` dengan
`1 = Internal` dan `2 = External`.

---

## 22 · `Update Total` ada di LIMA tab, dan rumusnya satu Activity

Briefing menyebut `Update Total` sebagai tombol tunggal di bawah Maximum Retention. Terukur di
`Section/TreatyInTabsNonProportional.xml`, ia ada di **lima** tab:

| Tab | `pyLabel` |
| --- | ---: |
| Maximum Retention | @202.657 |
| EGNPI | @843.732 |
| Limits | @1.641.962 |
| Share | @3.188.309 |
| Installment | @3.828.165 |

Yang dibangun ronde ini **satu** — yang briefing tunjuk. Keempat sisanya tercatat di
`TOTAL_RETENSI.tabLain` dan dijaga ujinya.

⚠️ Ada tombol **kedua** tepat di atas yang hidup, @194.468, dan ia **mati**: penjaganya
`pyCondition 1=2` @198.452, tanpa label (`.pyTemplateButton`). Yang dibangun yang hidup — dan
yang mati dicatat di sini, bukan dilewatkan diam-diam.

Tombol yang hidup: `pyAction refresh` @203.023 (param `retention` @203.951),
`pyDisabledWhen TreatyIn.EDMMaterialType = 2` @202.558,
`pyVisible` → `TreatyIn.IsEditData !='1'` @207.693.

⛔ Di layar kita ia **dinonaktifkan, bukan dihilangkan**: layar ini baca-saja, dan totalnya sudah
dihitung setiap kontrak dibaca. Tombol hidup yang tidak mengubah apa pun berbohong; tombol hilang
menyembunyikan bahwa layar lama punya langkah ini.

### 22.1 `Total Retention Amount` — dua kolom, dan yang pertama mata uang

`pyPageListProperty TreatyIn.TotalRetentionAmountNP` @147.731; judul kolom `Total Retention
Amount` @152.591 (isinya `.Currency` @161.105) dan `Value` @156.684 (isinya `.Value` @166.072).

⛔ Jadi nama panelnya **adalah judul kolom pertama**, dan kolom itu berisi mata uang. Briefing
menyebutnya berkolom `Value` saja.

**Rumusnya dari `Activity/TreatyInNPSetTotal.xml` cabang `param.type=retention`**
(`pyRuleAvailable = Yes`, nol penjaga mati): jumlah `.Amount` tab Maximum Retention,
dikelompokkan per `.Currency`, berurut **kemunculan pertama** (`<APPEND>`), dengan baris
bermata-uang sama dijumlahkan ke baris yang sudah ada (`Appendflag`).

⚠️ Dihitung di **services**, bukan di layar. Rumus yang hidup di dua tempat punya satu yang basi.

⚠️ Penjumlahannya `math/big`, bukan `float64`: `0,1 + 0,2` biner memberi `0,30000000000000004`,
dan angka itu lalu tampil sebagai nilai yang tidak pernah diketik siapa pun.

---

## 23 · Ketujuh rumus §6 — yang dijalankan satu, dan sebabnya enam sisanya tidak

⛔ **Enam panel sisanya TIDAK dibangun atas tebakan.** Activity-nya jelas; yang tidak jelas
**masukannya** di jalur baca kita.

| Panel | Activity + langkah | Masukan yang dituntut | Ada pada kita? |
| --- | --- | --- | --- |
| `Total Retention Amount` | `TreatyInNPSetTotal` `param.type=retention` | `.Currency` · `.Amount` | ⭐ **ya** — `Retensi.MataUang`/`.Jumlah`. **Dibangun.** |
| `Total Share` (prop) | `TreatyInNonSetTotal` langkah `share` | `.NusareLimit` · `.GrossPremium` · `.NetPremium` | ⛔ nol padanan pasti di `M_TREATY_IN2` |
| `Summary of Limit` | `TreatyInNPSetTotal` `param.type=limits` | `.Limit` · **`.Deductible`** | ⚠️ `.Deductible` dipetakan ke `CEDANT_RETENTION` **51%** — bukan kepastian |
| `Summary of MDP` | idem, `TotalLimitPremiEarnNP` / `TotalLimitMDPNP` | `.PremiumEarnedList(1).Value` · `.MDPList(1).Value` | ⚠️ keduanya **larik bersarang** per layer; tabel datar kita satu nilai per baris |
| `Total All Layers` | idem, jumlah lintas layer per mata uang | sama dengan dua baris di atas | ⚠️ idem |
| `Total All Layers RNM Share` | idem, `TotalShareRnmNP` | `.Value` baris RNM share | ⚠️ `RNMSHARE` ada, tetapi ia **persen**, bukan nilai |
| ROL | `TreatyInNonSetTotal`: `(TotalLimitsPremiumEarned / TotalLimitsIOOLimit) * 100` | dua total di atas | ⚠️ menunggu keduanya |

⚠️ **Perbedaan bentuknya yang mengikat, bukan kemalasan:** Activity-nya menjumlahkan larik
bersarang (`.PremiumEarnedList(1).Value`) yang hidup di clipboard Pega; `M_TREATY_IN2` memuat satu
baris **datar** per layer. Menjumlahkan kolom datar **mungkin** memberi angka yang sama dan
mungkin tidak, dan "mungkin" bukan dasar bagi angka uang di layar.

Yang diperlukan agar keenamnya dapat dibangun: padanan pasti untuk `.Deductible`, dan jawaban
apakah `PremiumEarnedList` / `MDPList` pernah berisi lebih dari satu elemen per layer.

⚠️ Kedua Activity-nya **hidup dan bersih** — `TreatyInNPSetTotal` (825.279 bita bersih) dan
`TreatyInNonSetTotal` (99.026 bita bersih), keduanya `pyRuleAvailable = Yes` dengan **nol**
penjaga `1=2`/`Never`. Jadi yang menahan keenam panel itu bukan kode mati.

---

## 24 · `Treaty Year` — rumusnya KETEMU, dan ia tetap tidak dihitung ulang

`DataTransform/TreatyInSetTreatyYear.xml` (`pyRuleAvailable = Yes`, memo *"termination
automatically add 1 year from start date"*):

```
TreatyIn.TreatyYear  := @substring(TreatyIn.Commencement,0,4)
TreatyIn.Termination := @addCalendar(TreatyIn.Commencement,"1","0","0","0","0","0","0")
```

**Diadu dengan POOLDATA, 1.854 baris:**

| Yang diuji | Cocok |
| --- | ---: |
| `TREATYYEAR = SUBSTR(COMMENCEMENT,1,4)` | **1.849** |
| `TERMINATION = COMMENCEMENT + 12 bulan` | **9** |

⭐ Baris pertama mengesahkan rumusnya. ⛔ Baris kedua memutuskan **cara memakainya**: kalau
`Termination` sungguh terikat pada `Commencement`, ia akan cocok pada ribuan baris, bukan
sembilan. Keduanya **nilai awal saat masuk**, bukan aturan yang terus berlaku.

Karena itu layar tetap menampilkan **kolomnya apa adanya**. Menghitung ulang akan menimpa lima
kontrak yang tahunnya sengaja berbeda dari tahun mulainya (`1000858`, `1000894`, `1001149`,
`1001761`, `1001772`).

⚠️ `05/1` di tangkapan layar **tetap tidak terjelaskan** — pada kontrak yang sama (`1001856`)
`TREATYYEAR` bernilai `2026`. Pertanyaannya di `PERTANYAAN-TERBUKA-LAYAR-PEGA.md` §2.

---

## 25 · Medan tanggal kosong — CSS saja tidak cukup, dan sebabnya diukur

`<input type="date">` yang kosong menuliskan `dd/mm/yyyy` sendiri, dan di layar penuh medan
terisi teks itu terbaca seperti nilai.

⛔ **Jalur CSS murni diperiksa dan ia TIDAK dapat bekerja.** Tidak ada pemilih yang membedakan
medan tanggal kosong dari yang terisi:

| Yang dicoba | Mengapa gagal |
| --- | --- |
| `:placeholder-shown` | tidak cocok untuk `type="date"` |
| `:invalid` | medan tanggal kosong tanpa `required` itu **`:valid`** |
| `[value='']` | React menyetel `value` sebagai **properti**, bukan atribut |

Kaitnya karena itu harus datang dari markup — dan markup yang boleh disunting **bukan**
`inti/frontend/components/ui/dasar.tsx`. Jalan keluarnya pembungkus di berkas modul:
`TanggalRedup` di `FormKontrakTreatyIn.tsx` menambahkan `.trin__tanggal--kosong` selama nilainya
kosong, dan `treatyin.css` mewarnai `::-webkit-datetime-edit` transparan selama kosong **dan**
tidak difokus.

⭐ Tetap dapat diisi dan tetap terbaca: `:focus` membatalkan aturannya, kelasnya hilang begitu
nilainya terisi, dan ikon pemilih tanggal **tidak** disembunyikan.

⚠️ **Batasnya dinyatakan:** `::-webkit-datetime-edit` hanya ada di Chromium dan WebKit. Di
Firefox `dd/mm/yyyy` tetap tampak, dan tidak ada padanan standarnya. Menggantinya dengan
`type="text"` akan membuang pemilih tanggal bawaan — harga yang lebih mahal daripada masalahnya.

---

## 26 · ⛔ `1001856` adalah kontrak DRAFT

Briefing meminta bukti §0 dijalankan atasnya. Dijalankan, dan sifatnya perlu dinyatakan:

| | `1001856` | kontrak nyata |
| --- | ---: | ---: |
| nama kontrak | `dasdas` | *nama kontrak sebenarnya* |
| panjang `JSONDATA` | **1.764** bita | 33.000 – 56.000 bita |
| kunci kepala yang ada | **0** dari 8 | 8 dari 8 |

⛔ **Medan yang kosong di tangkapan layar itu kosong karena KONTRAKNYA kosong**, bukan karena
layar membaca salah. Uji `TestKontrakBuktiAdalahDraftKosong` menjaga sifat itu dan akan merah
pada hari kontrak itu terisi — hari itu bukti §0 perlu kontrak lain.

---

## 27 · Dokumen desain — 43 tangkapan layar yang MENGALAHKAN tiga dugaan kami

`D:\NUSANTARA RE APP\design treaty in.docx`, diekstrak ke
`D:\NUSANTARA RE APP\design-treaty-in-gambar\` (43 `.png`), diserahkan 5 Oktober 2026. Ia
memperlihatkan sistem **berjalan** dengan data nyata — kontrak `1001846` (prop) dan `1001841`
(non-prop) — dan itu kelas bukti yang berbeda dari ekspor XML: ekspor mengatakan apa yang
dirender, gambar mengatakan apa yang **terlihat**.

Sapuan lengkapnya di [`SAPUAN-43-GAMBAR-DESAIN.md`](SAPUAN-43-GAMBAR-DESAIN.md).

### 27.1 ⭐ Tiga dugaan yang DITUTUP

| | Dugaan lama | Yang gambar perlihatkan |
| --- | --- | --- |
| `Treaty Year` | `05/1`, mungkin turunan (§24) | **`2025` — tahun polos** (gambar 01). `05/1` artefak form belum tersimpan |
| `Accounting Mode` | label tidak diketahui (§21.1) | `underwriting` -> **"Underwriting Year"** (01); `loss` -> **"Loss Occuring"** (26) |
| `RNM Share` | panel atau tab? (pertanyaan §3) | **SUB-TAB di dalam `Share`**, kedua cabang (16/17/34) |

⭐ Ketiganya **sejalan dengan pengukuran ekspor ronde sebelumnya** — nol yang dibantah, tiga
yang dilengkapi. `RNM Share` khususnya: §20 mengukur ia bukan satu dari 11 wadah `TABBED` dan
jatuh di dalam wilayah tab `Share`; gambar memperlihatkan persis itu.

### 27.2 ⭐ Dua keputusan ronde sebelumnya DIBENARKAN oleh gambar

1. **`Bordereaux` proporsional-saja.** Gambar 26 memperlihatkan layar non-prop dengan
   `Bordereaux Note` **tanpa** `Bordereaux` di atasnya — padanan independen bagi `pyCondition`
   @60.649.
2. **`RNM as Treaty Leader` TIDAK disembunyikan.** §20 menolak menyembunyikannya walau rule
   `TreatyMasterInEDM` salah pada setiap kontrak. Gambar 01 dan 26 memperlihatkannya **tampil**
   di keduanya. Penolakan itu benar.

⚠️ Keduanya layak dicatat bukan karena menyenangkan, melainkan karena keduanya adalah keputusan
yang diambil **tanpa** gambar — dan keduanya dapat saja salah.

---

## 28 · Panel `Existing Policy for Master ID` — rantai sumber tiga langkah

Panel yang nol jejaknya di kode kita, dan sumbernya ditelusuri sampai SQL-nya:

| Langkah | Rule | Isi |
| --- | --- | --- |
| panel | `Section/InputTreatyInOffer.xml` @89.949 | `pyPageListProperty = PolisList.pxResults` @105.310 |
| pengisi | `Activity/FetchTreatyExistingProduction.xml` | kunci cari `@if(TreatyIn.EDMState="", TreatyIn.ID, TreatyIn.OLDID)`, lalu `RDB-List`, lalu `@substring(.CARI2,18)` |
| SQL | `RDBList/FetchTreatyInProductionUsingNooffer.xml` | `SELECT DISTINCT NOPOLIS, IDPEGA, QUARTER, QUARTER_YEAR FROM POOLDATA.TREATYINPRODUCTION WHERE SUBSTR(NOOFFER,1,7) = SUBSTR({InputData.CARI1},1,7)` |

⭐ **Dibuktikan dua arah.** Kontrak `1001841`: gambar memperlihatkan satu baris
`RNM-QR.T02.05.2025.11987` / `NB-147044`; Oracle memberi `NOPOLIS` yang sama dan `IDPEGA =
ASM-FW-GISFW-WORK NB-147044` — prefiksnya tepat **18 aksara**, persis `@substring(.CARI2,18)`.
Kontrak `1001846`: gambar `No items`, Oracle nol baris.

⚠️ `TREATYINPRODUCTION` **41.936 baris**, tabel warisan yang BELUM pernah dibaca modul ini.
Dibaca saja; `TestTreatyInProductionTidakBerubah` menjaga cacahnya.

---

## 29 · Dua bentuk tanggal, dan aturannya bukan "kepala lawan grid"

| Tempat | Bentuk | Gambar |
| --- | --- | --- |
| `Commencement`, `Termination` | `01/01/2025` | 01 |
| grid Rate of Exchange, Reporting Period, EGNPI | `01/01/25` | 01, 29 |
| **rincian** EGNPI `As At` | `18/01/2025` | 29 |
| grid Installment `Due Date` (isian) | `18/01/2024` | 38 |

⛔ **Installment membantah "grid memakai dua digit".** Yang terbaca dari kelima contoh:
**baris grid yang TERLIPAT** memakai dua digit; **medan** — kepala, rincian terbuka, dan kolom
yang isinya medan isian — memakai empat.

Ronde ini memasang `TanggalTampilPanjang` pada `Commencement` dan `Termination`; grid tetap
`TanggalTampil`. Installment dan rincian EGNPI **belum** — keduanya memerlukan pembedaan baris
terlipat dari baris terbuka, dan grid kita belum punya baris yang dapat dibuka.

---

## 30 · ⚠️ Aturan angka TUNGGAL terbantah — dan tidak diubah

§7 briefing menetapkan satu jumlah desimal per jenis. Gambar memperlihatkan jumlah desimal
**per medan**, dan dua medan bernilai sama di baris yang sama pun berbeda:

| Gambar | Medan | Tampil | Desimal |
| --- | --- | --- | ---: |
| 29 | EGNPI `Amount` | `137.849.315.068,00` | 2 |
| 29 | EGNPI `Amount in IDR` | `137.849.315.068` | **0** |
| 29 | `Proportion %` grid lawan rincian | `100,00` lawan `100,0000000000` | 2 lawan **10** |
| 30 | `100% Limits ( IDR )` lawan `Summary of Limit` | `750.000.000,00` lawan `750.000.000` | 2 lawan **0** |
| 32 | `Adjustment Rate %` · `ROL %` | `0,367` · `67,424` | **3** |
| 38 | `% Installment` lawan `% Total` | `25,00` lawan `100,0000` | 2 lawan **4** |

Dan **nol di ekor TIDAK dibuang**: `1,00` tetap `1,00` (gambar 01).

⛔ **Nol baris aturan angka diubah.** Yang menentukan di Pega adalah `pyDecimalPlaces` tiap
kontrol, dan nilai itu **tidak ikut diekspor**. Yang dapat disimpulkan dari gambar adalah bahwa
aturan tunggal itu salah — **bukan** aturan benarnya apa, dan menebaknya akan mengubah setiap
angka di setiap grid modul ini sekaligus. Pertanyaannya di
[`PERTANYAAN-TERBUKA-LAYAR-PEGA.md`](PERTANYAAN-TERBUKA-LAYAR-PEGA.md) §9.

---

## 31 · ⚠️ Gambar 36/37 TIDAK mengubah dasar §17

§6 briefing meminta ini dijawab sebelum Retro disentuh.

| Gambar | Isi |
| --- | --- |
| 36 | centang `Has Share To Retro:` **tidak tercentang** → nol isi dirender |
| 37 | **tercentang** → `Other Treaty Retro`, `Summarry of Other Treaty Retro`, `Total All Layers Other Treaty Retro` (4 panel) |

⛔ **Dan di gambar 37 SETIAP grid dan SETIAP panel berbunyi `No items`.** Kontrak `1001841` —
kontrak contoh non-proporsional milik pemilik proses sendiri — punya nol baris Retro.

§17 menyatakan Retro tidak dibangun **karena jarang**, bukan karena tidak ada spesifikasinya.
Gambar memberi **bentuknya**, bukan bukti pemakaiannya, dan isi kosongnya justru memperkuat §17.

⭐ Tiga hal yang gambar tambahkan dan layak dicatat untuk hari Retro dibangun:
1. Isinya digerbangi satu centang `Has Share To Retro`, bukan oleh syarat tab.
2. Strukturnya sejajar dengan `Share`: grid layer → `Summarry of …` → `Total All Layers …`.
3. Ejaan ekspornya `Summarry` dengan **dua r**, sama dengan `Summarry of RNM Share`.

**Nol baris Retro ditulis ronde ini.**

---

## 32 · Tab `Limits` proporsional bersarang TIGA tingkat — dan datanya belum punya rumah

Empat belas gambar (`03`–`15`) untuk satu tab, dan sebabnya struktural:

```
Limits → Kind of Treaty [Add] → SPECIAL SURPLUS [Delete]
       → Treaty Type (dropdown)
       → Treaty Group [Add] → ENGINEERING [Delete]
          → Class of Business (grid)
          → 100% Limit / Retention / Cession to R/I  (berdaftar mata uang, ☑ Auto calculate)
          → 11 sub-tab: Event Limits · Deduction In A · Deduction · Reserve ·
            Experience Premium Refund · PLA · Cash Loss Limit · Claim Cooperation ·
            LPC · EPI · Achievement
```

⭐ Kesebelas sub-tabnya **persis** kesebelas wadah `TABBED` di `Section/DetailLimits.xml` yang
diukur ronde sebelumnya — pengukuran itu kini punya gambarnya.

⛔ **Belum dibangun, dan bukan karena terlewat.** Tiga entitas yang hilang dari jalur baca kita:
`Kind of Treaty`, `Treaty Type`, `Class of Business` per treaty group. `M_TREATY_IN2` memuat
**satu baris datar per layer**; ia tidak memuat pohon ini. Membangun bentuknya di atas data
datar memberi pohon bercabang satu di mana Pega punya banyak — layar yang terlihat benar dan
berbohong tentang strukturnya.

---

## 33 · Permintaan perubahan `View File` — satu-satunya penyimpangan yang DIIZINKAN

Pemilik proses, pada keterangan gambar `25`: *"dan view upload saran dibuatkan pop up dan ini
diubah menjadi design nya bagus"*.

Di Pega, `View File` membuka **jendela Chrome terpisah** berjudul `ShowAttachmentTreaty`,
lengkap dengan bilah alamat `appdev.nusantarare.com/prweb/...`. Di sini ia menjadi modal di
dalam halaman.

⭐ `Modal` **dipanggil** dari `inti/frontend/components/ui/dasar.tsx`, tidak ditulis ulang —
modal kedua berarti dua perilaku tutup, dua perangkap fokus, dan dua tempat untuk salah.

⚠️ Isinya **dua kolom** (`File Name` · `Type`), bukan empat seperti panel induknya. Gambar 25
memperlihatkan tepat dua, dan menyamakan keduanya membuat modal ini bukan modal yang gambar
itu perlihatkan.

⛔ **Selebihnya nol penyimpangan.** Lebar kolom, urutan kolom, nama tombol, dan kalimat galat
tetap salinan.

---

## 34 · Judul kolom DIRALAT — empat tab memakai nama kunci dokumen, bukan judul layar

| Tab | Sebelumnya | Gambar | Sekarang |
| --- | --- | --- | --- |
| Portfolio | `Type` · `Type Portfolio` · `Description` | 02 | `Portfolio Type` · `Premium / Loss Type` · `Description` |
| Accumulation | `Period` · `Report Date` · `Sub Days` · `Sub Due Date` | 19 | `Period` · `Reporting Date` · `Submission Days` · `Submission Due` |
| EGNPI | 8 kunci **berurut abjad** | 29 | 6 kolom **urut layar**, lalu 2 yang tidak terlihat |
| Maximum Retention | 5 kunci berurut abjad | 26, 27 | `Treaty Group` · `Currency` · `Amount`, lalu 2 yang tidak terlihat |

⛔ **Nol kolom dihapus.** `Class of Business` dan `Note` tidak terlihat di gambar mode-baca, dan
§0 briefing memperingatkan persis itu: *"jangan menghapus kolom yang hanya tidak tampak di mode
ini"*. Keduanya pindah ke ujung. `Note` memang terlihat di **rincian** baris (gambar 27, 29),
jadi ia ada — hanya bukan di gridnya.

⚠️ `JENIS_*` ikut diurutkan ulang bersama `KOLOM_*`-nya. Larik golongan yang tidak ikut pindah
akan memformat kolom teks sebagai uang tanpa satu uji pun berbunyi — dan itulah sebabnya
keduanya dijaga sepanjang di `desain-pega.test.ts`.

---

## 35 · Kategori lampiran BERCABANG — temuan yang tidak ada di ekspor

Gambar `24` (prop) dan `42` (non-prop) memberi kesebelas nama kategori, dan menemukan hal baru:
**daftarnya berbeda menurut cabang.** Sepuluh nama identik; yang kesepuluh berbeda —
`Pega Proportional Calculation /Perhitungan Pega Proportional` lawan
`Pega Non Proportional Calculation /Perhitungan Pega Non Proportional`.

⛔ **Ini TIDAK menutup `PERTANYAAN-TERBUKA-KODE-KATEGORI-LAMPIRAN.md`.** Yang terbuka pasangan
**kode↔nama**; gambar hanya memberi namanya. Urutannya di layar alfabetis, dan ronde 4 Oktober
2026 melarang keras menyimpulkan kode dari urutan abjad.

---

## 36 · ⛔⛔ `M_TREATY_IN2` DICABUT sebagai sumber — satu-satunya sumber `M_TREATY_IN`

Keputusan pemilik proses, 5 Oktober 2026, dinyatakan tiga kali. Bukan pemindahan sebagian, bukan
soal empat tab: **nol jalur baca** di seluruh modul.

⚠️ Tabelnya **tetap ada** di Oracle (7.281 baris) dan **tetap terlarang disentuh** — nol `DROP`,
nol `INSERT`, nol `UPDATE`, nol migrasi. Yang dicabut jalur bacanya di aplikasi.

### 36.1 Berkas yang berubah

| Berkas | Yang terjadi |
| --- | --- |
| `repository/warisan_in2.go` | **DICABUT** — pembacanya |
| `repository/warisan_in2_test.go` | **DICABUT** — 3 uji, sebabnya di §36.4 |
| `repository/warisan_in2_db_test.go` | **DICABUT** — 4 uji, sebabnya di §36.4 |
| `repository/warisan_layer_dokumen.go` | **BARU** — penggantinya |
| `repository/warisan_layer_dokumen_test.go` | **BARU** — 9 uji |
| `repository/warisan_layer_dokumen_db_test.go` | **BARU** — 3 uji Oracle |
| `backend/nol_m_treaty_in2_test.go` | **BARU** — penjaga, 2 uji |
| `repository/warisan_kontrak.go` | `jsonWarisan` memperoleh `Limits []jsonLimit`; `k.Layer` terisi di sini |
| `services/services.go` | seam `BacaLayerWarisan` **dicabut** |
| `services/warisan_kontrak.go` | pemanggilannya dicabut |
| `services/services_test.go` · `handlers/rute_treaty_in_test.go` | gudang tiruannya dicabut |
| `models/warisan_layer.go` | +7 medan Event Limits (§37); komentar sumbernya diperbarui |
| `cmd/ceklimit_main.go` | **DICABUT** — alat ukur sekali pakai, satu-satunya kueri nyata yang tersisa |
| `frontend/labels.ts` | petunjuk kosong dicabut (§38); `KOLOM_EVENT_LIMITS` 3 → 9 kolom |
| `docs/PEMETAAN-M-TREATY-IN2.md` | ditandai **catatan sejarah**, tidak dihapus |

⚠️ `repository/pendaratan_muat.go` dan `migrations/430_tabel_tab_treatyin.sql` hanya menyebutnya
di **komentar** sebagai yang TIDAK disentuh. Komentar itu tetap benar dan tinggal.

### 36.2 ⭐ Penjaganya, dan bukti ia merah

`TestNolKueriMTreatyIn2` menyapu **seluruh berkas Go repositori** dan merah bila ada yang
menyebut `M_TREATY_IN2` di luar daftar-izin. Daftar-izinnya pendek dan tiap entrinya bersebab;
`TestDaftarIzinMTreatyIn2TidakBasi` menolak entri yang tidak terpakai lagi.

**Dibuktikan dua arah.** Satu berkas palsu dipasang:

```go
const kueriPalsu = `SELECT * FROM POOLDATA.M_TREATY_IN2`
```

→ `FAIL: TestNolKueriMTreatyIn2 … bukti_penjaga_sementara.go menyebut M_TREATY_IN2`.
Dicabut → hijau kembali.

⭐ Dan penjaga kedua menangkap kesalahan saya sendiri saat ronde ini ditulis: saya memasang
`services/services.go` ke daftar-izin, lalu membersihkan komentarnya sehingga entri itu basi —
`TestDaftarIzinMTreatyIn2TidakBasi` merah, dan entrinya dicabut.

### 36.3 ⭐ Bentuk `Limits[]` — DIUKUR sebelum dibangun

Pohon tiga tingkatnya **ADA**, dan ronde sebelumnya mencarinya di tempat yang salah. Ia bukan di
kunci puncak; ia **di dalam** elemen `Limits[]`:

```
Limits[]                  satu elemen per LAYER              4.210 elemen
  .Layer .LayerType .Cover .Currency .Currency2 .Limit .Limit2
  .Deductible .Deductible2 .AgregateLimit .AdjRate .ROLPct .MDPPct
  .CurrencyRelation .TreatyType .MDPList[] .PremiumEarnedList[]
  .Reinstatement_List[] .TreatyGroupList[] .LayerList[]
  └ .Detail[]             satu elemen per TREATY GROUP       5.716 silang
      .TreatyGroup .TreatyType .CessionPct .QSPct .Brokerage .RNMShare
      .RIOGR .Earthquake .FloodJab .FloodNation .RSMDLimit
      .SpreadingType .SpreadingList[] .IOOLimitList[] .RetentionList[]
      .CessionList[] .EPIList[] .PLAList[] .CashLossList[]
      .ClaimCoopList[] .DeductionList[] .ReserveList[]
      .RNMShareList[] .RNMSpreadedList[] .RNMSpreadedListRI[]
      .AchievementLists[]
      └ .COBList[]        satu elemen per CLASS OF BUSINESS
          .ClassOfBusiness .ClassOfBusinessID .TreatyGroup .TreatyGroupID
```

⭐ `Kind of Treaty` → `Detail[].TreatyType`; `Treaty Group` → `Detail[].TreatyGroup`;
`Class of Business` → `Detail[].COBList[].ClassOfBusiness`. Ketiganya ada.

**Pemetaan kolom → jalur DIBUKTIKAN dengan pengaduan nilai**, 1.340 kontrak yang ada di kedua
sumber, disejajarkan menurut **nilai `Layer`** (bukan indeks baris — penyejajaran indeks memberi
13–48%, penyejajaran nilai memberi angka di bawah):

| Kolom | Jalur | Cocok |
| --- | --- | ---: |
| `LAYER` | `.Layer` | **100,0%** |
| `MDP_RATIO` | `.MDPPct` | **100,0%** |
| `CURRENCYRELATION` | `.CurrencyRelation` | **99,8%** |
| `BASIS_COVER` | `.Cover` | **98,6%** |
| `LAYERTYPE` | `.LayerType` | **91,6%** |
| `ADJ_RATE` | `.AdjRate` | **90,3%** |
| `ROL` | `.ROLPct` | **87,7%** |
| `TREATYTYPE` | `.TreatyType` | **81,2%** |
| `EARN_PREMIUM` | `.PremiumEarnedList[].Value` | **80,0%** |
| `CEDANT_RETENTION` | **`.Deductible`** | **70,6%** |

⭐ Baris terakhir menutup dugaan `PEMETAAN-M-TREATY-IN2.md` §5 yang berkeyakinan **51%** — dan
gambar `30` dokumen desain memberi kolom `Deductible ( IDR )` dengan nama itu di layar.

⚠️ **Kolom ber-angka rendah (12–45%) SELURUHNYA di tingkat `Detail[]`**, dan itu menjelaskan
dirinya sendiri: tabel memipihkan LAYER × TREATY GROUP, sementara pengadu hanya membaca
`Detail[0]`. Cocok ketika barisnya kebetulan treaty group pertama. Angka rendah itu **bukan**
tanda padanan yang salah — ia tanda tabelnya pipih dan dokumennya bersarang.

### 36.4 Ketujuh uji lama — dicabut beserta sebabnya, atau diubah artinya

⛔ Nol yang hilang diam-diam.

| Uji lama | Nasib |
| --- | --- |
| `TestPindaiLayerSejajarDenganDaftarKolom` | **DICABUT** — nol daftar kolom dan nol pemindai tersisa; JSON bernama tidak punya urutan untuk melenceng |
| `TestKolomLayerSamaDenganKatalogOracle` | **DICABUT** — digantikan `TestNolKueriMTreatyIn2`, yang menjaga hal lebih keras: bukan "namanya cocok" melainkan "tabelnya tidak disentuh" |
| `TestKueriLayerBerurutAngka` | **DIUBAH** → `TestUrutanDokumenDipertahankan`. Urutan dahulu milik `ORDER BY`; kini milik dokumen, dan yang dijaga **berbalik** menjadi "jangan diurutkan ulang" |
| `TestBacaLayerMengisiMedanDariKolomYangBenar` | **DIUBAH** → `TestLayerDokumenMengisiMedanDariJalurYangBenar` |
| `TestLayerBerurutMenurutAngka` | **DIUBAH** → sama dengan baris ketiga, dan tidak lagi memerlukan Oracle |
| `TestKontrakTanpaBarisLayerMengembalikanKosong` | **DIUBAH** → `TestNolLimitsMemberiIrisanKosong`, non-db: nol `Limits[]` sifat dokumen |
| `TestJangkauanTabelLayerTerukur` | **DIUBAH** → `TestLubang510Tertutup`. Dahulu MENGUNCI lubangnya; kini membuktikannya **tertutup**, sambil tetap mengunci cacah 7.281 |

**Cacah uji db `treatyin/backend/repository`: 93 → 98.** Akuntansinya persis:

```
  93  sebelum
−  3  warisan_in2_test.go      (paket `repository`, ikut terhitung di bawah -tags db)
−  4  warisan_in2_db_test.go
+  9  warisan_layer_dokumen_test.go
+  3  warisan_layer_dokumen_db_test.go
= 98
```

⚠️ Briefing memperkirakan angkanya **turun**. Ia naik, sebab penggantinya lebih banyak daripada
yang dicabut — tujuh uji atas satu kueri menjadi dua belas uji atas satu pohon.

Dua uji penjaga (`nol_m_treaty_in2_test.go`) tinggal di paket `treatyin/backend`, bukan
`repository`, jadi di luar hitungan itu.

### 36.5 ⭐ Lubang 510 kontrak TERTUTUP

| | |
| --- | ---: |
| kontrak dengan `Limits[]` berisi | **1.850** |
| di antaranya ada di `M_TREATY_IN2` | 1.340 |
| **lubang** | **510** kontrak · 1.210 elemen |

Sesudah pencabutan, jangkauannya **1.850**. `TestLubang510Tertutup` membuktikannya pada kontrak
nyata: `1001588` ada di dalam lubang dan kini memberi **12 baris layer** di mana dahulu nol.

⚠️ Penyaring teks `'"Limits":['` memberi **1.851**; selisih satu itu dokumen ber-`"Limits":[]` —
larik kosong, bukan kunci yang hilang. Yang membedakan pengurai, bukan substring.

---

## 37 · ⭐ Pencabutan itu MENGEMBALIKAN tiga kolom Event Limits

`M_TREATY_IN2` hanya berkolom `EARTHQUAKE`. Layar lama menampilkan **empat** batas — gambar `28`
dokumen desain: `RSMD Limit` · `Earthquake Limit` · `Flood Limit (Jabodetabek)` ·
`Flood Limit (Nationwide)`, masing-masing dengan mata uangnya sendiri.

Dokumen punya keempatnya di `Limits[].Detail[]`:

| Jalur | Kemunculan |
| --- | ---: |
| `.RSMDLimit` + `.CurrencyRSMD` | 493 |
| `.Earthquake` + `.CurrencyEarthquake` | 569 |
| `.FloodJab` + `.CurrencyFloodJab` | 218 |
| `.FloodNation` + `.CurrencyFloodNat` | 550 |

⛔ Jadi pencabutan bukan sekadar memindahkan sumber: ia **mengembalikan tiga kolom yang tabel itu
tidak pernah bisa berikan**. `KOLOM_EVENT_LIMITS` 3 → 9.

---

## 38 · Dua petunjuk kosong DICABUT — keduanya menjelaskan sebab yang sudah tidak ada

| Petunjuk | Bunyi lama | Mengapa dicabut |
| --- | --- | --- |
| `petunjukLayer` | *"…M_TREATY_IN2, yang mencakup 1.340 dari 1.854 kontrak … 510 kontrak punya limit di dokumennya tanpa satu baris pun di sana."* | Lubangnya tertutup (§36.5). Kosong kembali punya **satu** arti. |
| `petunjukEventLimits` | *"…tiga batas lain … tidak ada di tabel itu dan belum punya sumber tabel."* | Ketiganya kini kolom di layar (§37). |

⛔ Membiarkan keduanya berarti layar menjelaskan sebab yang sudah tidak ada — dan petunjuk yang
keliru lebih buruk daripada tidak ada petunjuk, sebab ia menghentikan orang mencari.

Keempat uji layar yang menuntut bunyi lama **dibalik, bukan dihapus**, masing-masing dengan
sebabnya tertulis di tempatnya.

---

## 24 · Jangkauan ronde 69 DIPERSEMPIT — desimal per kolom di form Treaty In

**Keputusan pemilik proses 5 Oktober 2026**, atas bacaan yang ronde sebelumnya tawarkan ketika
ia **menolak** membangun §23.

| | |
| --- | --- |
| **Apa** | Ronde 69 berlaku pada **ketiga layar yang ia sebut sendiri saja**. Form Treaty In mengikuti gambar desain: desimal **per kolom**, nol di ekor **DIPERTAHANKAN** sampai presisi kolomnya. |
| **Kenapa** | Gambar desain memperlihatkan `1,00` dan `16.000,00`. Ronde 69 memperlihatkan `2484250` tanpa ekor. **Keduanya benar — pada layar yang berbeda.** |
| **Yang TIDAK berubah** | ⛔ `format.ts` tidak disentuh. `formatNumber` tetap membuang nol di ekor; yang memadankan kembali `padankanDesimal` di lapis modul, dan hanya untuk layar ini. |
| **Pembalikan** | Ekspor `pyDecimalPlaces` tiba dan memperlihatkan presisi kolom yang berbeda dari gambar. Saat itu angkanya diambil dari ekspor, bukan dari gambar. |

### 24.1 ⛔ Ini MEMPERSEMPIT ronde 69, bukan membatalkannya

Ronde 69 **tetap berlaku penuh** pada ketiga layar yang ia sebut sendiri, dan ketiganya disebut
di sini supaya pembaca berikutnya tidak membacanya sebagai dicabut seluruhnya:

1. **Laporan Realisasi** (view existing)
2. **Limit Treaty In**
3. **grid XOL**

Nol di antaranya form Treaty In. Kalimat ronde 69 — *"`desimal` adalah BATAS ATAS, bukan panjang
tetap… ditiru per jenis: 6 uang realisasi · 3 share/komisi · 2 limit/TSI · 1 Pct arrangement"* —
dan kutipan pemilik proyeknya — *"DI PEGA 2484250 MAKA DI SISTEM BARU JANGAN
2484250.000000000"* — tidak berubah satu huruf pun untuk ketiga layar itu.

⚠️ **Dan penghentian ronde sebelumnya benar.** Premis briefing saat itu — bahwa ronde 69 menolak
pemadanan ke batas *global* — keliru: angka `6` dalam *"pada batas 6"* **adalah** presisi kolom
uang realisasi. Menurutinya berarti membangun kembali persis perilaku yang pemilik proyek cabut.
Yang menyelamatkannya membaca keputusan lama sebelum menimpanya.

### 24.2 Di mana pemadanannya hidup

`padankanDesimal` di `modul/treatyin/frontend/pages/FormKontrakTreatyIn.tsx`, dipanggil
`selAngka` **sesudah** `formatNumber`.

⛔ Ia **bukan pemformat kedua**, dan bedanya penting: ia tidak menguraikan angka, tidak
membulatkan, tidak mengelompokkan ribuan, dan tidak menyentuh apa pun di depan koma. Ia hanya
**menambahkan** nol di ekor sampai panjang yang kolomnya minta.

Tiga sifatnya dijaga uji:

- **Tidak pernah memotong.** Yang sudah lebih panjang lewat apa adanya — memotong berarti
  membulatkan, dan membulatkan menghilangkan digit berarti.
- **Teks bukan-angka lewat apa adanya.** `>=30% up to < 50%` tidak punya ekor untuk dipadankan.
- **Tanpa presisi kolom, aturan ronde 69 tetap berlaku.** `selAngka('uang', '2484250')` tetap
  memberi `2.484.250`.

### 24.3 Tabel kolom → desimal, dengan gambar yang membuktikannya

| Kolom | Desimal | Gambar | Nilai di layar |
| --- | ---: | ---: | --- |
| `Value to IDR` (Rate of Exchange) | **2** | 01, 26 | `1,00` · `16.000,00` · `15.500,00` |
| EGNPI `Amount` | **2** | 29 | `137.849.315.068,00` |
| EGNPI `Amount in IDR` | **0** | 29 | `137.849.315.068` — **baris yang sama** |
| EGNPI `Proportion %` (grid) | **2** | 29 | `100,00` |
| Maximum Retention `Amount` (grid) | **0** | 26 | `3.500.000.000` |
| Installment `Amount` | **2** | 38 | `161.168.704,90` |
| Installment `Pct` | **2** | 38 | `25,00` |

⭐ Baris kedua dan ketiga adalah buktinya sendiri: **dua kolom di baris yang sama, presisi
berbeda.** Satu aturan per golongan tidak dapat menghasilkannya.

### 24.4 ⛔ Yang TIDAK terbaca — didaftarkan, bukan ditebak

Keempat ini tampil di gambar tetapi **bukan sebagai kolom grid**; layar kita belum punya
tempatnya, jadi angkanya dicatat tanpa dipasang:

| Medan | Desimal | Gambar | Mengapa belum dipasang |
| --- | ---: | ---: | --- |
| `Proportion %` di RINCIAN baris | 10 | 29 | grid kita belum punya baris yang dapat dibuka |
| `Adjustment Rate %` | 3 | 32 | medan rincian layer, belum dirender |
| `ROL %` (rincian layer) | 3 | 32 | idem |
| `% Total` (Installment) | 4 | 38 | medan panel, bukan kolom grid |

Dan **seluruh kolom tab Limits · Share · RNM Share** tidak terbaca di gambar mana pun — ketiganya
tetap memakai aturan lama, dan `desimalPadan` mengembalikan `null` untuk semuanya. Dijaga uji
`⚠️ kolom yang TIDAK terbaca di gambar tetap null — bukan ditebak`.

⚠️ `null` **bukan nol**. Nol berarti "padankan sampai nol desimal" (`Amount in IDR`); `null`
berarti "tidak diketahui, pakai aturan lama". Membedakan keduanya adalah seluruh pokok §24.

---

## 25 · Tab Limits — pohon tiga tingkat, dan kedua cabang berbeda di PUNCAK

Delapan belas gambar dibaca (`02`–`15` prop, `30`–`33` non-prop) sebelum satu baris ditulis.

### 25.1 Susunannya

```
PROP (gambar 02–15)                NON-PROP (gambar 30–33)
Kind of Treaty          [Add]      Layers                  [add Layer]
 └ Treaty Type (dropdown)           └ Layer N / Part of Layer N
 └ Treaty Group         [Add]       └ Treaty Group  [Add Treaty Group] + ROL Profile
    └ Class of Business                └ Class of Business
    └ 100% Limit / Retention /         └ Egnpi this layer
      Cession to R/I                   └ Cover · Currency Relation · 100% Limit ×2
    └ 11 sub-tab                          Agregate Year Limit ×2 · Deductible ×2
                                       └ Reinstatement · MDP · ROL %
                                     ⤷ Summary of Limit · Total All Layers · Total ROL
```

### 25.2 ⭐ Mengapa puncaknya berbeda — DIUKUR, bukan dikira

Sapuan seluruh dokumen, 5 Oktober 2026:

| Medan | Proporsional | Non-proporsional |
| --- | --- | --- |
| `Limits[].TreatyType` | **terisi**, 1.360 elemen, 9 nilai (QUOTA SHARE 724 · SURPLUS 590 · SPECIAL SURPLUS 12 …) | **nil** pada 2.847 dari 2.850 |
| `Limits[].LayerType` | **nil** pada SELURUH 1.360 | **terisi**, 2.850 elemen (`layer` 2.680 · `sublayer` 170) |

⭐ **Setiap cabang memakai medan pengelompokan yang cabang lain kosongkan.** Satu bentuk layar
untuk keduanya akan memberi satu kelompok bernama kosong di salah satunya.

Dijaga `TestMedanPuncakBercabang`, yang mengukurnya lewat jalur baca sungguhan: prop
`jenisTreaty` 105 / `jenisLayer` 0; non-prop `jenisTreaty` 0 / `jenisLayer` 83.

### 25.3 Pemetaan ketiga tingkat

| Tingkat | Jalur dokumen |
| --- | --- |
| `Kind of Treaty` (prop) | `Limits[].TreatyType` |
| `Layers` (non-prop) | `Limits[].Layer` + `.LayerType` |
| `Treaty Group` | `Limits[].Detail[].TreatyGroup` |
| `Class of Business` | `Limits[].Detail[].COBList[].ClassOfBusiness` |

⚠️ **`Detail[]` jarang**: dari 4.210 elemen `Limits[]`, hanya **1.365** punya `Detail[]` berisi.
Itu menjebak uji pertama ronde ini — kontraknya dipilih menurut UKURAN dokumen, dan kedelapan
dokumen terbesar ternyata punya `Detail: []` kosong (besarnya datang dari `RevisionHistory`).
Ujinya memilih menurut **adanya `Detail[]`**, dan sebabnya tertulis di tempatnya.

### 25.4 ⚠️ `Deductible` ditandai, tidak ditampilkan seolah pasti

Padanan `CEDANT_RETENTION` ↔ `Limits[].Deductible` cocok **70,6%** atas 1.340 kontrak — jauh di
atas dugaan 51% dan cukup untuk memetakan, **tidak** cukup untuk diam. Layar menandainya
`(padanan 70,6%)` dengan keterangan lengkap pada `title`, dan uji menjaga bahwa **hanya** medan
itu yang ditandai.

### 25.5 Yang BELUM dibangun di tab ini

| | Sebabnya |
| --- | --- |
| `Treaty Type` dropdown (gambar 03) | calonnya `Detail[].SpreadingType` — 10 nilai yang bentuknya persis sama (`2023 QS 150M TRT`) — tetapi ia tingkat `Detail[]` sementara layar menaruhnya di tingkat `Kind of Treaty`. Memasangnya berarti memilih `Detail[0]` dan menyebutnya milik seluruh kelompok |
| 11 sub-tab tingkat terdalam (gambar 05–15) | isinya ada di dokumen (`PLAList`, `CashLossList`, `ClaimCoopList`, `DeductionList`, `ReserveList`, `EPIList`, `AchievementLists`) — yang belum ada bentuk layarnya |
| `Summary of Limit` · `Total All Layers` · `Total ROL` (gambar 30) | rumusnya menunggu jawaban apakah `PremiumEarnedList`/`MDPList` pernah berisi >1 elemen |

⛔ Nol di antaranya dibangun atas tebakan.

---

## 26 · ⚠️ Tiga penjaga lintas-aplikasi MERAH akibat migrasi `436` — dua di antaranya CACAT PENJAGANYA

Migrasi `436` (pergantian nama tabel tab, dijalankan pemilik proses) membuat tiga uji di
`inti/backend/penjaga` merah. ⛔ **Nol di antaranya disentuh** — ketiganya milik pemilik penjaga
lintas-aplikasi. Yang dikerjakan ronde ini hanya **membuktikan sebabnya**.

### 26.1 ⭐ Oracle TIDAK PUNYA bentuk ber-skema untuk mengganti nama sequence

Dibuktikan **terhadap instans yang dipakai** — Oracle 12.2.0.1.0 — dengan nama yang dijamin
tidak ada, supaya setiap pernyataan pasti gagal dan nol objek nyata tersentuh. Yang dibedakan:
galat **sintaks** lawan galat **objek tidak ada**; yang kedua berarti sintaksnya sah.

| Pernyataan | Jawaban Oracle | Artinya |
| --- | --- | --- |
| `ALTER SEQUENCE POOLDATA.x RENAME TO y` | **ORA-02286** *no options specified for ALTER SEQUENCE* | ⛔ `ALTER SEQUENCE` **nol klausa `RENAME`** |
| `ALTER SEQUENCE x RENAME TO y` | **ORA-02286** | idem, tanpa skema pun |
| `RENAME POOLDATA.x TO y` | **ORA-01765** *specifying owner's name of the table is not allowed* | ⛔ awalan skema **dilarang keras** |
| `RENAME x TO POOLDATA.y` | **ORA-01765** | idem di sisi sasaran |
| `RENAME x TO y` | **ORA-04043** *object does not exist* | ⭐ **sintaksnya SAH** — hanya objeknya tiada |
| `ALTER TABLE POOLDATA.x RENAME TO y` *(pembanding)* | **ORA-00942** *table or view does not exist* | ⭐ TABEL memang punya bentuk ber-skema |

⭐ **Kesimpulan: tidak ada cara menulis pergantian nama sequence yang menyebut skema.** Baris
pembanding terakhir yang mengunci bacaan itu — tabel punya bentuknya, sequence tidak.

⛔ **Jadi `TestSetiapPernyataanSahDanBerskema` dan `TestPernyataanMulaiDenganPerintah` menuntut
sesuatu yang tidak dapat ditulis. Itu CACAT PENJAGANYA, bukan cacat migrasi `436`.**

⚠️ `RENAME` pun **perintah SQL yang sah** di Oracle — ia pernyataan tersendiri, bukan klausa.
Penjaga kedua menolaknya hanya karena ia tidak ada di daftar perintahnya.

Perbaikan keduanya milik pemilik penjaga lintas-aplikasi: menerima `RENAME` sebagai perintah,
dan mengecualikan `RENAME` dari syarat ber-skema. ⛔ **Tidak dikerjakan di sini.**

### 26.2 Yang ketiga menunggu keputusan, bukan perbaikan

`TestTCONolNamaTabelBaruDiKode` merah karena menganggap awalan `T_TREATY_*` milik tabel warisan
modul lain. Apakah nama itu sah bergantung pada bentrok §4.3
[`PEMETAAN-SKEMA-NUSANTARARE.md`](PEMETAAN-SKEMA-NUSANTARARE.md) — dan keputusan itu milik
pemilik proses. ⛔ Penjaganya tidak disentuh.

---

## 27 · Layar dipecah meniru `masterproductnamelife` — dan §1.3 diputuskan per butir

### 27.1 Pemecahan — MURNI, dan buktinya angka uji tidak bergerak

| | Sebelum | Sesudah |
| --- | ---: | ---: |
| `pages/FormKontrakTreatyIn.tsx` | **1.826** baris | **752** |
| `frontend/components/` | tidak ada | **12 berkas** |
| Komponen terbesar | — | **184** (`PanelLampiran`) |
| Pembanding: komponen terbesar modul contoh | 372 | — |

⭐ **Uji yang berubah: NOL.** Sebelum 172, sesudah 172 (yang 178→187 kemudian datang dari uji
BARU ronde ini dan dari `saring.test.ts` milik agen sebelah).

⚠️ Caranya: `FORM` di ketiga berkas uji kini membaca **halaman + seluruh `components/`**, jadi
tiap pernyataan tetap menanyakan hal yang sama — *"apakah kode layar modul ini memuat ini"* —
tanpa satu pun disunting. Dua pernyataan yang menanyakan **LETAK di dalam halaman** (irisan,
`lastIndexOf`) memakai konstanta `HALAMAN` yang hanya halaman; menggabungkan komponen akan
menggeser letaknya dan membuat pernyataannya menanyakan hal yang berbeda.

### 27.2 ⚠️ Cara yang GAGAL lebih dulu, dicatat supaya tidak diulang

Percobaan pertama memakai **pencocokan kurung** dari `s.index('{', awalFungsi)`. Ia tersandung
parameter berdestrukturisasi: pada `function Panel({ baris }: {...}) {`, kurung pertama adalah
milik `{ baris }`, sehingga blok "berakhir" di situ dan fungsi terbelah di tengah — lima komponen
sekaligus, dan halamannya rusak.

⭐ Yang dipakai akhirnya **berbasis baris**: tiap deklarasi tingkat atas di berkas ini berakhir
pada baris yang TEPAT `}`. Sederhana, dan tidak punya cara gagal diam-diam.

⚠️ Pemulihannya mungkin karena badan yatimnya masih utuh di halaman. ⛔ Pelajarannya: **salin
berkas lebih dulu** sebelum pembedahan berskrip — ronde ini hampir kehilangan `selAngka` dan
lima saudaranya sebab penulisan `angka.ts` gagal (direktorinya belum ada) **sesudah** halaman
ditulis ulang tanpa mereka.

### 27.3 §1.3 — keempat komponen modul contoh, diputuskan per butir

| Komponen contoh | Ditiru? | Sebabnya |
| --- | --- | --- |
| **`DropdownCari.tsx`** (262) | ⭐ **YA** — polanya, pada `DropdownWarisan` | Pemilih kita nol papan tik; 131 cedant tidak dapat digulir dengan nyaman |
| **`Medan.tsx`** (151) | ⛔ **TIDAK** | `inti/.../dasar.tsx` sudah memberi `Field`, `Pilih`, `Area`, `FieldTanggal`, dan layar ini memakainya. Menyalin lapis ketiga berarti tiga tempat mengatur satu kotak isian |
| **`Dialog.tsx`** (57) | ⛔ **TIDAK** | `Modal` dari `inti` sudah dipakai `View File` (§33). Dialog kedua berarti dua perilaku tutup dan dua perangkap fokus |
| **`PanelLampiran.tsx`** (372) | ⛔ **TIDAK disalin** | Kita **sudah punya** `PanelLampiran` sendiri (184 baris), dari ekspor Pega modul ini — `WorkAttachments.xml` dan `ShowAttachmentTreaty.xml`. Bentuknya berbeda sebab layarnya berbeda |
| **`Saran.tsx`** (112) | ⛔ **TIDAK** | Ia panel saran modul contoh; layar ini punya panel `History` sendiri dari `M_TREATYIN_COMMENT`. Nol kebutuhan yang belum terlayani |

⛔ **Komponen yang ditiru tanpa kebutuhan adalah kode mati yang terlihat rapi.**

### 27.4 Papan tik pemilih — yang dipasang

Enter/panah-bawah **membuka** · panah atas-bawah dan PageUp/PageDown **menggeser** · Enter
**memilih** · Escape **menutup tanpa memilih** · Tab menutup. Baris aktif digulirkan ke dalam
pandangan, disebut lewat `aria-activedescendant`, dan diberi kelas sendiri — ⚠️ **bukan
`:hover`**, sebab papan tik tidak menghasilkan hover dan tanpa kelas itu yang memakai papan tik
tidak melihat di mana ia berada.

⛔ **Geserannya DIJEPIT, tidak melingkar** — sama dengan `geserAktif` modul contoh. Daftar yang
melingkar membuat yang menekan terus tidak pernah tahu ia sudah di ujung.

⚠️ **Dua beda yang disengaja dari modul contoh**, dicatat supaya yang membandingkan nanti tahu
keduanya pilihan:

1. `DropdownCari` **menolak** ketikan bebas; pemilih ini **harus** dapat diketik — permintaan
   pemilik proses.
2. **Nama kembar diberi pengenalnya** (`ASURANSI ADIRA DINAMIKA — G0000006`), hanya yang kembar.
   Modul contoh tidak menghadapi soal itu; kita menghadapinya pada **36 nama cedant**.

---

## 28 · ⭐ `PremiumEarnedList`/`MDPList` — terjawab dengan MENGUKUR

Pertanyaan yang sempat masuk daftar "menunggu orang" ternyata **menunggu pengukuran**. Sapuan
seluruh 1.854 dokumen, 4.210 elemen `Limits[]`:

| Larik | Panjang 0 | 1 | 2 | >1 | di antaranya BEDA mata uang |
| --- | ---: | ---: | ---: | ---: | ---: |
| `PremiumEarnedList` | 2 | 2.748 | **100** | 100 | **100** |
| `MDPList` | 3 | 2.741 | **106** | 106 | **106** |
| `EgnpiTotalList` | 2 | 2.748 | 100 | 100 | 100 |
| `MDPMinList` | 25 | 274 | 11 | 11 | 11 |
| `IOOLimitList` | 3 | 1.354 | 5 | 5 | 5 |
| `RetentionList` | 2 | 1.356 | 4 | 4 | 4 |

⭐ **Jawabannya: ya, dan maksimum DUA** — dan pada **setiap** larik berpanjang 2, kedua elemennya
**bermata uang berbeda**. Jadi lariknya *satu elemen per mata uang*, tidak pernah dua elemen
semata uang.

⇒ Panel `Summary of Limit` dan `Total All Layers` **terbuka**: rumusnya jumlah per mata uang,
pola yang sama dengan `Total Retention Amount` yang sudah berdiri.

⛔ **Dan ia memperlihatkan cacat pada yang sudah dibangun:** `nilaiPertama()` di
`warisan_layer_dokumen.go` mengambil elemen **[0]** saja, sehingga pada **100–106 elemen layer**
mata uang kedua tidak terlihat di layar. Itu **bukan tebakan melainkan ukuran**, dan
perbaikannya menuntut bentuk baris layer berubah dari satu nilai menjadi per-mata-uang —
pekerjaan ronde berikutnya, dicatat di sini supaya tidak terlupa.

---

## 29 · `nilaiPertama` — cacat ditutup, dan ia ternyata menyentuh EMPAT medan

### 29.1 Yang hilang, terukur

`nilaiPertama()` membaca elemen `[0]` saja. Terukur atas 1.854 dokumen:

| Larik | Berpanjang 2 | Kedua elemennya bermata uang BERBEDA |
| --- | ---: | ---: |
| `MDPList` | 106 layer | **106** |
| `PremiumEarnedList` | 100 layer | **100** |

⛔ Jadi yang hilang **bukan nilai kembar** melainkan nilai yang berdiri sendiri — 106 angka MDP
dan 100 angka premi yang tidak pernah sampai ke layar.

⭐ **Dan cacatnya lebih luas daripada yang dilaporkan:** `Limit2` (2.441 kemunculan) dan
`Deductible2` (2.436) **tidak pernah dibaca sama sekali** — bukan dipotong, melainkan tidak
dipetakan. Jadi empat medan, bukan dua.

### 29.2 Bentuk penggantinya — DIBACA DARI GAMBAR

Gambar `30`: grid Layers berkolom `100% Limits ( IDR )` · `Deductible ( IDR )` ·
`100% Limits ( USD )` · `Deductible ( USD )`; panel `Summary of Limit` berkolom `MDP (IDR)` dan
`MDP (USD)`. **Pega menampilkan keduanya berdampingan, bukan memilih salah satu.**

⇒ `nilaiKe(l, i)` menggantikan `nilaiPertama(l)`, dan model memperoleh empat medan pasangan:
`Limit100Kedua` · `RetensiCedantKedua` · `MDPKedua` · `PremiEarnedKedua`.

⚠️ **Pasangannya sejajar, dan itu terbukti**: pada kontrak `1000003` tiap layer ber-`Currency =
IDR` dan `Currency2 = USD`, sementara `MDPList[0]` bermata uang IDR dan `MDPList[1]` USD.

⛔ Lariknya **tidak dijumlahkan** — dua mata uang berbeda yang dijumlahkan memberi bilangan yang
bukan uang apa pun.

### 29.3 Ujinya — positif memakai layer bermata uang DUA

⚠️ Layer bermata uang tunggal **lulus bahkan dengan cacatnya**; uji positif yang memakainya
tidak membuktikan apa pun. Karena itu:

| Uji | Isinya |
| --- | --- |
| `TestMataUangKeduaTerbaca` | bentuk disalin dari kontrak **`1000003` Limits[0]**, keempat pasangnya diadu |
| `TestMataUangKeduaKosongBilaTunggal` | ⛔ medan kedua **kosong**, bukan salinan yang pertama |
| `TestLarikNilaiTidakPernahLebihDariDua` | merah pada hari elemen ketiga muncul |
| `TestMataUangKeduaPadaKontrakNyata` | Oracle: **7 dari 7** baris layer `1000003` membawa mata uang kedua |
| `TestCacahLayerMataUangDuaTerukur` | Oracle: **MDP 106 · PremiEarned 100**, persis angka dokumen |

---

## 30 · Panel total — Activity DIBACA, dan ia MENGHITUNG

`Activity/TreatyInNPSetTotal.xml`, cabang `param.type=limits`, langkah `Sum Total Limit` dan
`Set Total Limit`:

```
local.Value  := local.Value  + .Limit      ⇒ TotalLimitIOONP(<APPEND>).Currency := Limits(1).Currency
local.Value2 := local.Value2 + .Limit2     ⇒ TotalLimitIOONP(<APPEND>).Currency := Limits(1).Currency2
```

dan serupa untuk `TotalLimitDeductblNP` dengan `.Deductible` / `.Deductible2`.

⭐ **Ia MENGHITUNG, bukan menyalin** — berbeda dari preseden
`TreatyEDMProRateCalculation` yang ternyata satu `Copy value`. Rumusnya **jumlah per SLOT mata
uang**: slot 1 = `Currency`, slot 2 = `Currency2`.

⚠️ Itu persis bentuk yang §29 buka. Panel `Summary of Limit` dan `Total All Layers` karena itu
**siap dibangun**, dan ⛔ **belum dibangun ronde ini** — dinyatakan, bukan disamarkan.

---

## 31 · Tab Event Limits — empat baris berlabel, bukan grid

Gambar `28` menyusunnya sebagai empat baris berlabel, masing-masing dengan mata uangnya sendiri:
`RSMD Limit` · `Earthquake Limit` · `Flood Limit (Jabodetabek)` · `Flood Limit (Nationwide)`.
Bentuk sebelumnya grid sembilan kolom — **benar isinya, salah bentuknya**.

⚠️ **Satu beda dari gambar, disengaja:** gambar memperlihatkan SATU himpunan, sementara dokumen
menyimpannya per LAYER × TREATY GROUP. Menampilkan satu himpunan berarti memilih satu baris dan
menyembunyikan sisanya — persis cacat §29 yang baru ditutup. Jadi tiap baris layer memperoleh
himpunannya sendiri, berkepala layer dan treaty group-nya.

### 31.1 ⚠️ RALAT: tab yang belum punya layar ada TIGA, bukan delapan

Briefing menyebut delapan. Diukur terhadap cabang `tabTampil` di layar:

| Tab | Keadaan |
| --- | --- |
| `Retro` | ⛔ sengaja tidak dibangun — §17, dan gambar 37 memperkuatnya |
| `Achievement In IDR` | belum punya sumber |
| `Value Difference` | syaratnya tidak pernah terpenuhi; isinya milik modul Adjustment |

Kesebelas tab lain **sudah merender**. Yang kurang pada sebagiannya bukan layar melainkan
**bentuk** — grid datar di tempat gambar memperlihatkan bentuk lain, seperti Event Limits
sebelum ronde ini.

---

## 32 · Tab Portfolio — kolomnya tertukar, dan kendalinya salah bentuk

Permintaan pemilik proses 6 Oktober 2026: *"`Portfolio Type` dan `Premium / Loss Type` itu
dropdown … pastikan juga nilainya didapat dari mana, lihat dari XML-nya Treaty In, jangan ada
ngarang."*

Yang ditemukan sewaktu membaca ekspornya **dua** hal, dan yang kedua lebih berat daripada yang
diminta.

### 32.1 ⛔ RALAT: kolom 1 dan kolom 2 TERTUKAR sejak 5 Oktober 2026

Grid `TreatyIn.Portfolio` di `Section/TreatyInTabsProportional.xml` punya dua baris sel yang
berpasangan lurus — judul di baris 1, properti di baris 2, empat kolom masing-masing:

| kolom | judul | sel | properti | sel |
| --- | --- | --- | --- | --- |
| 1 | `Portfolio Type` | 107 | `.TypePortfolio` | 112 |
| 2 | `Premium / Loss Type` | 108 | `.Type` | 113 |
| 3 | `Description` | 109 | `.Description` (`pxTextArea`) | 114 |
| 4 | `Add` | 110 | `Delete` | 115 |

Catatan di `labels.ts` berbunyi **persis sebaliknya** — *"kunci `Type` menjadi kolom `Portfolio
Type`"* — dan layar mengikutinya. Akibatnya kolom berjudul `Portfolio Type` menampilkan
`Premium`/`Loss`.

⭐ **Tiga saksi yang berdiri sendiri-sendiri:**

1. Pasangan sel di atas, yang tidak dapat ditafsirkan dua cara.
2. **Judul kolom 2 menyebut nilainya sendiri.** `Premium / Loss Type`, dan yang bernilai
   `Premium`/`Loss` adalah kunci `Type` — terukur, §32.2.
3. Ekspor modul **Treaty In Adjustment**, dua seksi terpisah (`TreatyInTabsProportional` dan
   `TreatyInTabsProportionalOldData`), berpasangan persis sama.

⚠️ Ejaannya `TypePortfolio`, **tanpa spasi**. Catatan lama menulis `Type Portfolio`; nol dari
1.925 elemen memakai ejaan itu.

### 32.2 Daftar pilihannya TIDAK ADA di ekspor — jadi ia diukur

Kedua sel ber-`pyControlDisplayTitle` = *"Control inherited from property"* dengan
`pyListDataSource` **kosong**: daftarnya hidup di `Rule-Obj-Property`, dan ⛔ **aturan properti
tidak ikut diekspor** (korpus `D:\XML_NURE\Treaty In` hanya Activity · ConnectREST ·
DataTransform · DecisionTable · FlowAction · Harness · RDBList · ReportDefinition · Section ·
SystemSettings · When).

Jadi daftarnya **diukur**, bukan dikarang dan bukan pula diambil dari tempat yang tidak
memuatnya. Sapuan atas **seluruh 1.855** dokumen `M_TREATY_IN` (nol `ROWNUM`, nol gagal urai) —
1.925 elemen di 845 dokumen:

| Kunci | Nilai | Cacah |
| --- | --- | ---: |
| `TypePortfolio` | `Withdrawal` · `Assumption` | 1.578 · 347 |
| `Type` | `Premium` · `Loss` | 1.032 · 893 |

Nol kosong pada keduanya. Keempat kombinasi terpakai — sejalan dengan
`UQ_PORTOFOLIO (ID_VERSI_KONTRAK, ARAH_PORTOFOLIO, JENIS_PORTOFOLIO)` migrasi `406`.

⚠️ Yang pengukuran tidak dapat katakan: nilai **sah tetapi belum pernah dipakai**. Pertanyaannya
di `PERTANYAAN-TERBUKA-PILIHAN-PORTFOLIO.md`, dan harganya murah — nilai asing **tidak
dijatuhkan**, ia ditawarkan di bawah daftar.

### 32.3 Satu sakelar mengatur seluruh tab: `TreatyIn.IsEditData`

```
sel 112·113·114   pyReadOnlyCondition = TreatyIn.IsEditData= 1
sel 110·115       pyCondition         = TreatyIn.IsEditData!='1'
```

Satu tanda, dua arah — persis `mode` di layar ini.

⭐ **Karena itu `GRID_BERTOMBOL_TAMBAH` diralat.** `Portfolio` dulu terdaftar sebagai grid
**tanpa** `Add`, atas dasar gambar `02`. Gambar itu tangkapan layar **mode-BACA**; tombolnya ada,
yang menyembunyikannya keadaan layarnya. Uji yang menegakkan ketiadaannya diganti berikut
sebabnya tertulis.

⚠️ Tab ini memakai `IsEditData` sementara grid lain di seksi yang sama memakai
`TreatyIn.ViewState` (terukur: 9 sel `ViewState`, 3 sel `IsEditData`). Keduanya dipetakan ke
`mode` yang sama, dan bedanya dicatat supaya tidak terbaca sebagai kelalaian.

⛔ **`pyDisabledWhen` = `TreatyIn.EDMMaterialType = 2` TIDAK dibawa.** Ia ada di keempat sel,
tetapi medannya belum terbaca layar ini — menyalakannya berarti menebak nilainya.

### 32.4 ⚠️ Temuan menyamping, di luar repo dan TIDAK diubah

`2-to-spec/PENELUSURAN-JSON-KE-KOLOM.md` 267–268 memetakan `Type` → `ARAH_PORTOFOLIO` dan
`TypePortfolio` → `JENIS_PORTOFOLIO`. Pengukuran mengatakan sebaliknya, dan
`TestTiket24PortofolioEmpatKombinasiDiterima` (`MASUK`/`KELUAR` × `PREMI`/`KLAIM`) berpihak pada
pengukuran. Berkasnya di luar repo, tabel `PORTOFOLIO` belum dimuat satu baris pun, dan
menukarnya keputusan pemilik proses.
