# STATUS

Diperbarui: **2026-09-18** — sesi grilling Claim Non Prop selesai; ringkasan serah-terima di `RINGKASAN.md`

> **Folder dirapikan 20 September 2026** menjadi empat folder hasil: `1-grilling/` (22 berkas) · `2-to-spec/` (39) · `3-to-tickets/` (2) · `4-erd-dan-tabel-datar/` (9). Nama berkas **tidak ada yang berubah**, hanya tempatnya. Nama → tempat: [`PETA-FOLDER.md`](PETA-FOLDER.md); isi tiap folder: `ISI-FOLDER.md` di dalamnya. Yang tetap di akar: register ini, `KEADAAN-AKHIR.md`, `alat/`, `pengetahuan/`, `docs/adr/`, `_selesai/`, `sampah/`.

Golongan: **PERMANEN** tidak pernah dihapus · **BERUMUR** hidup selama pertanyaannya terbuka · **SEMENTARA** dihapus begitu isinya terserap · **TURUNAN** dibangkitkan dari berkas lain dan tidak pernah disunting — bila ia berbeda dari sumbernya, alatnya yang salah.

Tata letak: **akar** berisi kesimpulan dan register yang hidup · **`pengetahuan/`** berisi bahan yang masuk — sumber, data mentah, alat · **`docs/adr/`** keputusan · **`_selesai/`** arsip.

Pemisahan `pengetahuan/` dikerjakan 18 September 2026. Aturannya: **kalau isinya tidak saya karang, ia di sana.** Tafsir harus bisa dibantah oleh bahannya, jadi keduanya tidak dicampur. Rujukan di seluruh folder induk sudah disetel ke jalur baru (16 jalur unik, semuanya terverifikasi). **Berkas di `_selesai/` sengaja tidak disetel** — arsip mencatat keadaan pada saat ia diarsipkan.

| BERKAS | GOLONGAN | STATUS | DIPERBARUI |
|---|---|---|---|
| `RINGKASAN.md` | PERMANEN | berkas serah-terima — dibaca pertama sesudah `CONTEXT.md` | 2026-09-18 |
| `PENGETAHUAN.md` | PERMANEN | **gabungan 40 artefak** jadi satu bacaan, 243 KB. Turunan — bila berbeda, berkas aslinya yang berlaku | 2026-09-18 |
| `KEADAAN-AKHIR.md` | **PERMANEN** | berkas penutup untuk orang yang datang tanpa sesi ini: empat pencacah beserta perintahnya, tiga batas yang berdiri, urutan pemasangan 36 berkas, dan apa yang harus terjadi supaya pekerjaan berlanjut. **Setiap angka di dalamnya berperintah** — diperiksa dengan sapuan penutup 19-09 | 2026-09-19 |
| `Diagram-Skema-Tabel-ClaimNonProp.xlsx` + `ERD-CLAIM-NON-PROP.html` | **TURUNAN** | **rancangan tabel Claim Non Prop, penamaan `T_CLAIM_`** — awalan `T_CLAIMNP_` DIHAPUS 20-09-2026 mengikuti [keputusan work owner] di berkas rujukan. 33 tabel · 30 relasi · 10 halaman yang bukan tabel. 11 tabel kini **dipakai bersama** Claim Prop / Claim Fac In, 16 khas Non Prop. Tiga tabel dipecah agar cocok dengan struktur lama sekaligus berkas rujukan. Direkonsiliasi: 22/22 `ddl-usulan/` terpadankan, induk BERBEDA dari rujukan **nihil**, 3 tabel baru tanpa padanan (`T_CLAIM_ADJ_LAYER`, `_ADJ_LAYER_CURRENCY`, `T_CLAIM_RETRO`) | 2026-09-20 |
| `alat/buat-skema-claimnp.py` | PERMANEN | pembangkitnya. **Hanya membaca** `struktur-claimdata-lama.md`, `MEMORI_PEMAHAMAN.MD`, `pengetahuan/`, `ddl-usulan/`. Memeriksa tiap halaman Pega yang dikutip memang ada di pohonnya | 2026-09-20 |
| `ERD-ORACLE.xlsx` + `ERD-ORACLE.md` | **TURUNAN** | **ERD kotak-entitas dalam Excel, TABEL Oracle saja** — sepuluh sheet. Sistem lama: 32 tabel, 627 kolom, **nol foreign key**, 15 relasi **tersirat** yang dibaca dari SQL ekspor XML. Skema baru: 22 tabel, 396 kolom, 17 foreign key tertulis. VIEW, PROCEDURE, FUNCTION, dan halaman clipboard Pega **tidak digambar**. Dibangkitkan `alat/buat-erd-excel.py` dari `pengetahuan/SCHEMA-ACTUAL.csv`, `ddl-usulan/`, dan ekspor XML **langsung** | 2026-09-19 |
| `Diagram-Skema-Tabel-Gabungan.xlsx` + `ERD-GABUNGAN.html` | **TURUNAN** | **satu peta untuk KEDUA berkas rancangan, sesudah penyatuan nama 20-09-2026.** 45 tabel · 38 relasi pohon · **10 tali penghubung** · lima lajur: INTI · **T_CLAIM_ dipakai bersama** (Prop · Fac In · Non Prop) · khas Non Prop · Claim Life · Polis. Garis penghubung **digambar** — border sel di Excel, `<line>` di HTML; **nol** titik garis jatuh di dalam kotak. Direkonsiliasi: **51/51** relasi sheet *Daftar Relasi* tertutup | 2026-09-20 |
| `alat/buat-skema-gabungan.py` | PERMANEN | pembangkitnya. **Hanya membaca** kedua workbook rancangan. Memeriksa ulang seluruh titik garis terhadap seluruh sel yang ditempati kotak — bila ada yang bertumpuk, alat berhenti, bukan menggambar | 2026-09-20 |
| `alat/buat-erd-excel.py` | PERMANEN | pembangkit workbook ERD. **Hanya membaca** sumber; mencetak nilai yang seharusnya dikeluarkan tiap rumus supaya dapat diadu dengan isi sel | 2026-09-19 |
| `4-erd-dan-tabel-datar/PETA-NAMA-TABEL.md` + `peta-nama-tabel.tsv` | **TURUNAN** | peta nama lama → nama berlaku, 27 tabel `T_CLAIM_`, beserta sebab tiap pemecahan dan daftar lini pemakainya. TSV-nya dibaca `buat-erd-excel.py` untuk sheet **PETA-NAMA-T** — satu sumber, dua keluaran | 2026-09-20 |
| `PETA-FOLDER.md` + `*/ISI-FOLDER.md` | **TURUNAN** | peta nama-berkas → folder, 121 berkas terdaftar, dibangkitkan `os.walk`. Dokumen lama menyebut berkas lain **hanya dengan namanya** — berkas inilah yang membuat sebutan itu tetap dapat dilacak sesudah pemindahan | 2026-09-20 |
| `sampah/` | **MENUNGGU DIHAPUS** | bentuk yang digantikan: `ERD.html`, `ERD-KLAIMNP.html`, `keluaran/datar-*.tsv`, `datar-claimdata-*.csv`, beserta alatnya. Dipindahkan 19-09 atas permintaan pemilik proyek — yang diminta Excel, bukan HTML, dan isinya TABEL Oracle, bukan halaman Pega. Aturannya tetap: *saya yang memindahkan, Anda yang menghapus* | 2026-09-19 |
| `alat/papan.py` | PERMANEN | **disalin dari scratchpad 19-09** — sebelumnya hanya hidup di sesi, sementara `KEADAAN-AKHIR.md` menunjuknya. Ia **menulis satu berkas dan hanya satu**: `issues/README.md` | 2026-09-19 |
| `alat/periksa-penamaan.py` | PERMANEN | **hanya membaca** — ditulis ulang 19-09; versi sebelumnya ikut mengubah berkas DDL, dan itu berbahaya diserahkan ke orang yang belum tahu | 2026-09-19 |
| `RIWAYAT-SESI-2026-09-19.md` | PERMANEN | **baru** — riwayat sesi: urutan keputusan, premis yang gugur beserta sebabnya, dan lima angka yang dicabut. Transkrip mentah 20 MB tidak disalin; letaknya disebut di kepala berkas | 2026-09-19 |
| `KAMUS-KOLOM.md` | SEGAR | 22 tabel **396 kolom**, 9 view **56 kolom**. *Angka 452 dicabut 19-09 — ia jumlah keduanya, bukan kolom tabel; `sampah/keluaran/RINGKASAN-TABEL-DATAR.md` bagian 2.1.* Lima kolom baru 19-09: `LINI_USAHA` dan `PENUTUPAN_LAMA` di `KLAIM`; `BENTUK_TERURAI` dan `SEBAB_GAGAL_URAI` di `MIGRASI_PENDARATAN`; `ID_PENDARATAN` di `MIGRASI_NILAI_DITOLAK`. *Angka 391/368 pada versi lama dicabut — tidak dapat direproduksi dari berkasnya* | 2026-09-19 |
| `ddl-usulan/` (36 berkas) | SEGAR | 22 tabel + 9 view + skema/akun + hak akses + pengawasan + kunci pemilik. **Penamaan final 19-09**: 105 objek bernama ditulis ulang ke bentuk `<peran>_<tabel>[_n]`; satu tabel diganti nama (`DAFTAR_KLAIM_PENJAGA_TANGGAL` → `KLAIM_PENJAGA_TANGGAL`); pengenal di atas 30 byte **nol** | 2026-09-19 |
| `TICKETS.md` | PERMANEN | 39 tiket lapisan data, 8 jalur; 7 selesai | 2026-09-18 |
| `.scratch/.../issues/` | SEGAR | papan terbit — **aktif 0 · menunggu instance 8 · tertahan 1 · selesai 29 · mati 1**. **Seluruh tiket lapisan data ditutup 19-09.** Yang tersisa: 8 uji yang menunggu instance nyata, dan `11` yang menunggu batch aplikasi | 2026-09-19 |
| `.scratch/claim-non-prop-lapisan-data/ISI-FOLDER.md` | **PERMANEN** | **baru** — indeks isi folder papan: 39 tiket dikelompokkan menurut status, empat aturan papan beserta sebabnya, dan **sebab tiap tiket belum selesai** dibedakan (menunggu mesin, menunggu lingkup, premis gugur). Ditujukan kepada orang yang membuka folder itu tanpa konteks | 2026-09-19 |
| `REGISTER-RATIFIKASI.md` | PERMANEN | 8 butir AK, DECIDED-TEKNIS, menunggu ratifikasi akuntansi | 2026-09-18 |
| `CONTEXT.md` | PERMANEN | hidup | 2026-09-18 |
| `SPEC-MODEL-DATA.md` | SEGAR | bagian **16 ditulis ulang** (penamaan final, tabel singkatan dibatalkan) dan bagian **21 baru** (delapan aturan dari penutupan register). Satu keluaran tersisa menunggu **REQ-012**: pemetaan migrasi kolom demi kolom | 2026-09-19 |
| `BLUEPRINT.md` | PERMANEN | hidup | 2026-09-18 |
| `docs/adr/` (28 berkas, ADR-0001…0028) | PERMANEN | hidup | 2026-09-18 |
| `FINDING-001-threshold-currency.md` | PERMANEN | menunggu REQ-011 | 2026-09-18 |
| `FINDING-002-percabangan-identitas.md` | PERMANEN | menunggu REQ-016 | 2026-09-18 |
| `FINDING-003-baris-ur-tanpa-penyaring.md` | PERMANEN | tuntas dari XML | 2026-09-18 |
| `FINDING-004-suntingan-manual-tertimpa.md` | PERMANEN | bersyarat, menunggu sesi Komite | 2026-09-18 |
| `FINDING-005-perbandingan-tanggal-beda-format.md` | PERMANEN | menunggu REQ-020 | 2026-09-18 |
| `FINDING-006-kurs-tidak-ditemukan-bernilai-satu.md` | PERMANEN | menunggu REQ-022 | 2026-09-18 |
| `FINDING-007-dua-rumus-reinstatement.md` | PERMANEN | tuntas dari XML | 2026-09-18 |
| `pengetahuan/arithmetic-inventory.tsv` | PERMANEN | data mentah, 673 ekspresi | 2026-09-17 |
| `pengetahuan/rekonsiliasi-kolom-vs-properti.tsv` | PERMANEN | data mentah, 408 kolom × 659 properti | 2026-09-18 |
| `pengetahuan/README.md` | PERMANEN | penjelas isi `pengetahuan/` dan batas keamanannya | 2026-09-18 |
| `pengetahuan/ddl/` (49 berkas) | PERMANEN | untuk dibaca, bukan dijalankan | 2026-09-18 |
| `pengetahuan/DDL_Script_ClaimNonProp.xls` | PERMANEN | sumber 48 objek `POOLDATA` | 2026-09-18 |
| `pengetahuan/DDL_Script_ClaimNonProp2.xls` | PERMANEN | **belum dibaca** — muncul 2026-09-18 10:57, belum diperiksa isinya | 2026-09-18 |
| `SAPUAN-S3-S9.md` | PERMANEN | hasil sapuan S3, S4, S5–S9; melahirkan D9–D15 dan A16. **Lima koreksi diterapkan 18-09**: 273→79 berkas, tujuh baris/sembilan nama, satuan byte, `pzInsKey`, dan pencabutan klaim “sumber baru” | 2026-09-18 |
| `SAPUAN-CELAH-DAN-JAVA.md` | PERMANEN | penutupan celah tag (30→45), S11/S12/S16 dijalankan ulang beserta selisihnya, 480 baris Java dibaca sebagai kode, tiket `05` dan `06`. Melahirkan **D37–D46** | 2026-09-19 |
| `PERMINTAAN-DBA-1-VERSI-INSTANCE.md` | BERUMUR | **siap kirim** — REQ-032 sendirian, dua nilai, satu halaman | 2026-09-19 |
| `PERMINTAAN-DBA-2-PROFIL-DAN-STRUKTUR.md` | BERUMUR | 32 REQ, diurutkan menurut apa yang ditahannya; REQ-018, REQ-033, REQ-021 di atas | 2026-09-19 |
| `FINDING-008-kolom-idr-tanpa-kurs.md` | PERMANEN | **usulan** — kolom IDR terisi tanpa kurs; naik jadi temuan tetap hanya bila ramalannya terukur benar lewat REQ-034 | 2026-09-18 |
| `SAPUAN-S10-S16.md` | PERMANEN | sapu ulang S1/S2 dengan pola tiga-bentuk; rekonsiliasi DDL; rantai `.GrossValue` dan `.AdjusterFeeValue`; penanda yatim; tipe ganda; `CountLossAllocation_act` utuh; S14–S16. Melahirkan **D16–D36** | 2026-09-18 |
| `ddl-usulan/00_SKEMA_DAN_AKUN.sql` | PERMANEN | tiket `01`: dua tablespace, tiga akun, tiga peran, hak sistem, Bagian 5b (lubang hak `ANY`). **Sembilan `REVOKE` diperbaiki** — versi pertama menggugurkan pemasangan di instance bersih | 2026-09-19 |
| `ddl-usulan/Z01_PENGAWASAN_TULIS.sql` | PERMANEN | **baru** — pengawasan tulis dari luar `KLAIMNP_APP`, bersyarat `SESSION_USER`. Dijalankan DBA | 2026-09-19 |
| `ddl-usulan/Z02_KUNCI_PEMILIK.sql` | PERMANEN | **baru** — penguncian akun pemilik, langkah terakhir pemasangan. Dikeluarkan dari komentar | 2026-09-19 |
| `ddl-usulan/01_KLAIM.sql` | PERMANEN | tiket `12` **ditutup**: aggregate root ditulis ulang lengkap. Dua kolom baru — `LINI_USAHA` (SPEC §21.3) dan `PENUTUPAN_LAMA` (SPEC §21.6); empat aturan dinyatakan jatuh ke aplikasi | 2026-09-19 |
| `ddl-usulan/21_MIGRASI_PENDARATAN.sql` | PERMANEN | tiket `13` **ditutup**: gerbang `IS JSON` **dicabut** — muatan mendarat utuh, hasil penguraian dicatat sebagai fakta (D41, D43). `CK_..._2` menolak gagal-tanpa-sebab dan berhasil-tapi-bersebab | 2026-09-19 |
| `ddl-usulan/20_MIGRASI_NILAI_DITOLAK.sql` | PERMANEN | tiket `14` **ditutup**: `ID_PENDARATAN` menjadi seam ke tiket `13` — nilai tak terurai menunjuk **baris muatannya**, bukan hanya nama sumbernya | 2026-09-19 |
| `ddl-usulan/19_MIGRASI_KORELASI.sql` | PERMANEN | tiket `15`: `CK` sekurangnya satu pengenal lama; **dua index arah lama → baru** — arah yang justru dijanjikan judul tiketnya | 2026-09-19 |
| `ddl-usulan/V00_HAK_AKSES.sql` | PERMANEN | tiket `35`: **44 `REVOKE` diperbaiki** — versi lama menggugurkan pemasangan di instance bersih (`ORA-01927`), kelas cacat yang sama dengan `00_SKEMA_DAN_AKUN.sql` | 2026-09-19 |
| `ddl-usulan/V01`, `V02` | PERMANEN | tiket `29`: **view kompatibilitas kini menghasilkan SELISIH** terhadap `ARSIP_MUATAN_KELUAR`, bukan nilai mutlak. Pengelompokan kedua view disamakan — sebelumnya berbeda, sehingga ada kelompok yang jatuh dari keduanya | 2026-09-19 |
| `ddl-usulan/V03`, `V05`–`V09` | PERMANEN | tiket `30`, `31`: `CACAH_BELUM_LENGKAP` di keenam view agregat — ADR-0019 di tingkat agregat, bukan hanya di tingkat baris | 2026-09-19 |
| `ddl-usulan/V04_V_PARITAS_SHADOW.sql` | PERMANEN | tiket `23`: ditulis ulang dari **1 tabel menjadi 19**, ditambah cabang `TABEL_TIDAK_DIKENALI`. Versi lama membuang 18 tabel lewat penyaring keadaan, tanpa pesan | 2026-09-19 |
| `ddl-usulan/Z00_ISIAN_AWAL.sql` | PERMANEN | **baru** — tiket `21`: baris awal tutup buku dan tiga tarif. ⚠ **Satu-satunya berkas DML di seluruh `ddl-usulan/`**, dipisahkan supaya kekecualiannya terlihat. ID-nya lewat `NEXTVAL`, bukan angka tertulis tangan | 2026-09-19 |
| **22 sequence** `SQ_<tabel>` | PERMANEN | **baru** — sampai 19-09 tidak ada pembangkit pengenal sama sekali: nol `IDENTITY`, nol `SEQUENCE`, nol `DEFAULT` di 36 berkas, padahal hak `CREATE SEQUENCE` sudah diberikan. Sequence dipilih, bukan `IDENTITY`, dengan alasan yang sama yang menetapkan batas 30 byte | 2026-09-19 |
| `ddl-usulan/15_TUTUP_BUKU.sql`, `16_TARIF_BERLAKU.sql` | PERMANEN | dua kolom diganti nama menurut isinya (ADR-0021): `TANGGAL_TUTUP` → **`HARI_TUTUP_BUKU`**, `NILAI_TARIF` → **`NILAI_TARIF_PERSEN`** | 2026-09-19 |
| **23 constraint pasangan IDR–kurs** | PERMANEN | tiket `16`, `17`: arah terbalik **dan** saling mengunci — mengisi satu nilai IDR memaksa tiga lainnya terisi. Diperbaiki di 10 berkas; ditambah `KURS ⇒ KURS_SUMBER` (ADR-0029 butir kelima) di 10 tabel | 2026-09-19 |
| `docs/adr/0028-*.md` | PERMANEN | **`accepted` 19 Sep 2026** lewat konfirmasi manusia, dua alasan dicatat. Gerbang DDL terangkat 19-09; **gerbang penamaan ikut dicabut 19-09** — batas 30 byte jadi aturan tetap, bukan menunggu REQ-032 | 2026-09-19 |
| `docs/adr/0029-*.md` | PERMANEN | **baru** — kurs ikut tersimpan; `proposed`, dua cabang terbuka menunggu `ASK-AKUNTANSI` butir 7 | 2026-09-18 |
| `_selesai/OPEN-QUESTIONS.md` | **PERMANEN** | **Register TERTUTUP dan berkasnya DIPINDAHKAN** 19-09 — aktif **0** · selesai **57** · dihapus **2**. Tidak akan diperbarui lagi **kecuali sebuah cabang gugur**. Lampirannya **bahan hidup**: empat hipotesis, delapan cabang, nol dipilih. Penunjuk ditinggalkan di jalur lama | 2026-09-19 |
| `OPEN-QUESTIONS.md` (jalur lama) | PERMANEN | **penunjuk**, bukan salinan — satu paragraf yang menyatakan register ditutup, jalur barunya, dan bahwa lampirannya masih bahan hidup. Berdiri karena jalur lama masih dipegang dokumen di luar folder ini | 2026-09-19 |
| `ORACLE-REQUESTS.md` | SEGAR | REQ: **aktif 35 · selesai 1 · mati 1**. **REQ-004 dicabut**; tujuh REQ turun jadi **verifikasi** (001, 009, 010, 011, 019, 028, 034); **REQ-032 ikut turun jadi verifikasi**; tidak ada lagi REQ yang menahan keputusan rancangan. Empat menahan **pelaksanaan** tiket: REQ-018, REQ-033, REQ-021, REQ-037. **REQ-037 BLOCKER diperluas jadi enam jalur** — hak `ANY`, prosedur definer's rights, hibah `PUBLIC`, `CREATE ANY TRIGGER`, peran pemuat `ANY`, dan `GRANT ANY PRIVILEGE` yang membuatnya ukuran satu titik waktu | 2026-09-19 |
| `pengetahuan/PULL-LIST.csv` | BERUMUR | **kanonik**, 74 objek — isi v2, versi lama di `_selesai/` | 2026-09-18 |
| `pengetahuan/SCHEMA-ACTUAL.csv` | BERUMUR | 698 baris (598 + 100 tabel work), hidup sampai seluruh REQ terjawab | 2026-09-18 |
| `pengetahuan/PREFLIGHT.sql` | BERUMUR | T1–T3 terjawab dari DDL; T4, T5, T7 masih menunggu | 2026-09-18 |
| `pengetahuan/RECON.sql` | BERUMUR | perlu ditulis ulang untuk TOAD | 2026-09-18 |
| `_selesai/ASK-AKUNTANSI.md` | SEGAR | butir 1–6 tertutup lewat delapan butir AK; **butir 7 turun jadi verifikasi 19-09** — cabang A/B ADR-0029 dihapuskan, ADR naik `accepted` dan berlaku seragam. Berkas kembali arsip; **tidak ada butir yang menahan** | 2026-09-19 |
| `_selesai/TABLE-EXTRACTION-REQUEST.md` | BERUMUR → selesai | dipindah 2026-09-18 | 2026-09-18 |
| `_selesai/PULL-LIST-lama-2026-09-18.csv` | BERUMUR → selesai | 56 baris, digantikan v2 | 2026-09-18 |
| `_selesai/SCHEMA-ACTUAL-WORK-terserap-2026-09-18.csv` | SEMENTARA → selesai | isinya sudah masuk `pengetahuan/SCHEMA-ACTUAL.csv` | 2026-09-18 |

## Catatan

**Dua rotasi yang sempat tertunda sudah selesai 2026-09-18.** `pengetahuan/PULL-LIST.csv` kini berisi versi kanonik 74 objek; `pengetahuan/SCHEMA-ACTUAL.csv` kini 698 baris setelah 100 kolom `PC_ASM_FW_GCNMFW_WORK` disalin masuk, dipisahkan dari baris lain lewat `sumber=DDL_SQL` dan `req_id=DDL-WORK-2026-09-18`.

**`pengetahuan/DDL_Script_ClaimNonProp2.xls` belum pernah dibaca.** Berkas itu muncul di folder pada 2026-09-18 10:57 dan tidak disebut dalam instruksi mana pun. Dicatat, tidak diperiksa.
