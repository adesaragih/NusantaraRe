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
