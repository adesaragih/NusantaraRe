# PUTUSAN-01 — Adjudikasi atas `GRILL-01.md`

> Tahap 1.5 · masukan: `PENGETAHUAN.md` (19 Sep 2026) + `GRILL-01.md` (20 Sep 2026) · diputus 20 September 2026.
> Bukti tambahan yang dipakai penilai dan **tidak** dipakai interogator: `_migration-docs/claim-non-prop/pengetahuan/ddl/` — 49 berkas source Oracle yang sudah ada di repo sejak 18 September.

---

## 1. Putusan pokok

1. Dari 20 temuan: **6 berdiri apa adanya**, **13 diubah** (putusan atau tingkat), **1 ditutup** sebagai tidak layak diselidiki. Tingkat setelah adjudikasi: **P0 = 3** (sebelumnya 6) · **P1 = 6** · **P2 = 8** · **P3 = 3**.
2. **Pemblokir sesungguhnya: dua.** (a) satu paket **empat kueri baca-saja** ke basis data berjalan; (b) satu **ekspor ulang terbatas** tiga rule Pega. Bukan satu, karena dua paket itu dipegang pihak berbeda — DBA dan admin Pega — tidak dapat saling menggantikan, dan dapat berjalan **sejajar**; menyatukannya jadi satu "pemblokir" akan menyembunyikan fakta bahwa pertanyaan uang tidak menunggu pihak Pega sama sekali.
3. Bukti termurah yang interogator lewatkan: **source prosedur dan DDL tabel lama sudah ada di repo.** `PROCEDURE_PROC_GENERATE_SEQUENCE_NUMBER.sql` menutup Q-11 seluruhnya, `PROCEDURE_PEGA_JSON_OS_AKSEP_KLAIMTNP.sql` mengubah dasar G-09/G-11, dan `TABLE_DIRECTTOKASIR_LOG.sql` membuktikan muatan Kasir yang terkirim **tersimpan permanen** — sehingga usulan Tracer untuk Q-10 tidak perlu ada.
4. Dua `P0` interogator runtuh di tempat: **G-05 tertutup** (prosedur menangani tanggal NULL secara sengaja) dan **G-04 melingkar** (memakai ketiadaan bukti sebagai bukti ketiadaan, pada jenis data yang temuannya sendiri nyatakan tidak terekspor).
5. Satu `P0` justru **meluas**: G-09. Prosedur penulis OS mengisi `STS_REJECT` dari argumen bernilai `1`/`2`/`4`, sementara kueri pembalik menyaring `STS_REJECT = 0` — jadi persoalannya bukan hanya cakupan klaim vs adjustment, melainkan kemungkinan kueri itu **tidak mengembalikan apa pun**.
6. Gerbang dibelah ulang per aliran kerja: **tiga aliran boleh mulai besok**, dua beku, satu beku ringan. Penundaan seluruh modul tidak beralasan — dua dari tiga aliran yang boleh mulai tidak tersentuh satu pun temuan yang digugat.

---

## 2. Koreksi metodologis atas interogasi

**M-01 · Lingkaran pada G-04 — melanggar U-1.**
G-01 menetapkan bahwa parameter metode `Page-Copy` dan `Page-New` **tidak ada di ekspor**. G-04 lalu menyimpulkan, dari `grep` atas teks 338 berkas, bahwa "tidak ada yang mengisi `TempSpreadingRisk`". Sasaran salin sebuah `Page-Copy` hidup persis di parameter yang G-01 nyatakan hilang — jadi nama halaman itu **tidak akan pernah muncul** di grep sekalipun pengisinya ada. G-04 memakai ketiadaan bukti sebagai bukti ketiadaan, pada jenis data yang ia sendiri nyatakan tak terekspor. **Akibat: G-04 tidak boleh dipakai menopang G-03**, dan bacaan "muatan mungkin tidak pernah dibangun" kehilangan penopangnya.

**M-02 · Tracer diusulkan untuk perkara yang sudah punya saksi tertulis — melanggar U-7.**
Lima dari delapan pertanyaan baru (Q-9…Q-16) mengusulkan Tracer atau perekaman sesi. Padahal `PENGETAHUAN.md` §9.2 — yang interogator baca dan kutip — mencatat tiga tabel saksi: `DIRECTTOKASIR_LOG` (menyimpan **muatan JSON yang benar-benar terkirim**), `HISTORYAKSEPTASIPEGA` (menyimpan tiap keputusan), `OS_AKSEPTASI_KLAIM` (menyimpan tiap baris akseptasi dan pembaliknya). Ketiganya menjawab Q-10, Q-12, Q-14 dengan satu `SELECT`. Lebih jauh: DDL ketiga tabel itu **sudah ada di repo**. Tracer turun dari lima usulan menjadi **nol**.

**M-03 · Tingkat akibat melebihi tingkat sebabnya — melanggar U-5.**
G-01 disebut "cacat metodologis terbesar" tetapi diberi **P1**, sementara G-02 — yang tidak lain adalah satu kejadian dari G-01 — diberi **P0**. Sebab dan akibat tertukar peringkat. Selain itu G-14 dan G-20 adalah **hasil hitungan** (296 kendali, 137 baca-saja, 15 wajib) namun diberi putusan `RAGU`; hitungan tidak ragu, ia benar atau salah.

**M-04 · Beban pembuktian tidak dibalik pada bacaan yang menyiratkan fungsi inti tak pernah jalan — melanggar U-2 dan U-3.**
Gabungan G-03+G-04 menghasilkan bacaan bahwa muatan pembayaran mungkin tidak pernah tersusun, sampai JSON-nya cacat. Bila itu benar, tidak ada satu pun akseptasi non-prop yang pernah sampai ke Kasir selama bertahun-tahun — sementara sistem lama memelihara tabel log khusus untuk kiriman itu, lengkap dengan kolom nomor akseptasi. Bacaan semacam itu wajib turun ke `RAGU` dengan beban pembuktian di pihak penuduh, bukan naik ke `P0`.

**M-05 · Gerbang per seksi mencampur dua jenis bukti — melanggar U-4.**
§9 dinyatakan "BELUM SIAP" seluruhnya karena pengikatan langkah→rule tidak terverifikasi. Namun **isi** rule SQL dan REST tidak digugat sama sekali — interogator sendiri memakainya sebagai bukti sah di G-09, G-10, dan G-13. Satu seksi memuat bukti yang runtuh dan bukti yang berdiri; membekukan seluruhnya menghentikan pekerjaan yang sudah siap. Karena itu gerbang dibelah ulang di bagian 6.

---

## 3. Adjudikasi per temuan — hanya yang berubah

| ID | Putusan interogator | **Putusan saya** | Alasan (menyebut uji) |
|---|---|---|---|
| G-01 | TAK TERVERIFIKASI · P1 | **SAHIH SEBAGIAN · P2** | U-7: pengikatan dapat dipulihkan gratis lewat pencocokan nama halaman bind — tiap rule SQL memakai halaman yang khas (`InsertHistory.*`, `ParamSeq.*`, `ParamKasir.*`, `InputParamOs.*`, `UploadDoc.*`), dan langkah `Property-Set` sebelum tiap `RDB-List` menyebut halaman itu. Yang benar-benar tersisa hanya 3 langkah. |
| G-02 | TAK TERVERIFIKASI · P0 | **TAK TERVERIFIKASI · P1** | U-5: akibat tidak boleh berperingkat lebih tinggi dari sebabnya (G-01); dan U-2: penulisan balik ke induk terbukti berhasil karena nomor akseptasi memang mendarat di data. |
| G-04 | RAGU · P0 | **TAK TERVERIFIKASI · P1** | U-1: grep tidak dapat melihat sasaran `Page-Copy` (lihat M-01); U-2: bacaan kuatnya menyiratkan pembayaran tak pernah terkirim. |
| G-05 | SALAH · P0 | **SALAH (penyajian) · P2 · TERTUTUP** | U-7: `PROCEDURE_PROC_GENERATE_SEQUENCE_NUMBER.sql` di repo menunjukkan `IF p_proddate IS NULL THEN v_now := SYSTIMESTAMP AT TIME ZONE 'Asia/Jakarta'` — tanggal kosong **ditangani secara sengaja**. Cacatnya tinggal cacat dokumentasi. |
| G-06 | SAHIH · P1 | **SAHIH · P2 (sebagian tertutup)** | U-7: `TABLE_OS_AKSEPTASI_KLAIM.sql` memuat kolom `CLAIMOLD`, sejalan memo rule "set CLAIMOLD" — arti `CARI20` praktis tertutup; hanya `CARI17` yang tersisa. |
| G-07 | RAGU · P1 | **RAGU · P3 · HENTIKAN PENYELIDIKAN** | U-4 + T-01: `stsReject` tidak punya kolom pendaratan (G-12), jadi kosong atau tidak kosong sama-sama tidak berakibat; nilai yang penting (`CARI10`) sudah sampai ke prosedur di langkah 3→6, sebelum penghapusan di langkah 7. |
| G-09 | SAHIH · P0 | **SAHIH · P0 · DIPERLUAS** | U-4: `PEGA_JSON_OS_AKSEP_KLAIMTNP` mengisi `STS_REJECT` dari argumen yang di jalur Pega bernilai `1`, `2`, atau `4`, sementara `GetDataCNPOS` menyaring `STS_REJECT = 0`. Persoalannya bertambah satu tingkat: bukan hanya cakupan, tetapi apakah kueri pembalik mengembalikan baris sama sekali. |
| G-10 | SAHIH · P1 | **SAHIH · P3 · TERTUTUP** | U-2 + U-3: idiom yang sama (`TO_DATE(TO_CHAR(sysdate,…),…)`) muncul di dalam `PEGA_JSON_OS_AKSEP_KLAIMTNP` yang terbukti jalan bertahun-tahun; ini gaya rumah untuk memangkas jam, bukan pergeseran tanggal diam-diam. |
| G-11 | SALAH · P0 | **SALAH · P0 · BUKAN PEMBLOKIR** | U-7: tidak menunggu bukti apa pun — DDL usulan dan sumber prosedur keduanya sudah di tangan. Ini pekerjaan `to-spec`, bukan pertanyaan penyelidikan. |
| G-13 | RAGU · P2 | **dibelah: `PostEmailKomiteCNP` P1 · enam lainnya P2** | U-5: satu rule menentukan penerima surat jenjang berikutnya (perilaku berbeda); enam sisanya utilitas kerangka (ketidakjelasan). |
| G-14 | RAGU · P2 | **SAHIH · P2** | U-5: hitungan bukan dugaan. |
| G-15 | SAHIH · P1 | **SAHIH · P2** | U-5 + T-05: keempat penanda hanya bermuara ke `InputParam.CARI51–54` (tampilan) dan parameter cetak — tidak ada nilai uang yang bergantung padanya. |
| G-17 | RAGU · P1 | **RAGU · P2** | U-7: satu `SELECT` atas properti case menutupnya; dan penolakan sirkulasi nol jenjang sudah boleh dibekukan apa pun hasilnya (bagian 9 butir 8). |
| G-19 | RAGU · P1 | **KEPUTUSAN BISNIS · P1** | U-4: tidak ada bukti yang ditunggu; ini pilihan model yang harus diputuskan, bukan keraguan atas fakta. |

**Diterima tanpa perubahan:** G-03, G-08, G-12, G-16, G-18, G-20.

---

## 4. Temuan turunan

**T-01 · G-07 dinetralkan oleh G-12.**
Mengikuti: G-07 + G-12. Bila `stsReject` memang tidak punya kolom pendaratan di model sasaran, maka apakah `InputData.CARI10` terhapus sebelum dibaca **tidak berakibat apa pun** pada sistem baru; yang berakibat hanyalah bahwa nilai itu sudah masuk ke prosedur OS sebelum penghapusan. Yang berubah di `PENGETAHUAN.md`: baris `stsReject` di §7.2 berhenti menjadi "kontrak tulis" dan menjadi catatan sejarah. Ditutup bersama: tidak ada — justru tidak perlu ditutup.

**T-02 · G-01 tidak menyentuh dasar bukti sebagian besar dokumen.**
Mengikuti: G-01. Seluruh rumus di §8.1, §8.3, §8.7, seluruh mesin keadaan §3, routing §4, dan inventaris layar §10 bersandar pada `Property-Set`, pra-syarat, transisi, rule `When`, dan XML Section — kelimanya **terekspor utuh** (166 dari 166 `Property-Set` membawa parameter). Yang berubah: gerbang "§8 BELUM SIAP" terlalu luas; yang runtuh hanya pengikatan langkah→rule di §9.2 dan penamaan rule per langkah di §5.2. Ditutup bersama: Q-9, yang setelah pencocokan nama bind menyusut jadi 3 langkah.

**T-03 · Periode buku ditentukan tiga kali, dan tambalan Pega hanya menyala ketika ketiganya berselisih.**
Mengikuti: G-05 + K-06 + sumber prosedur. Urutannya: (a) `PROC_GENERATE_SEQUENCE_NUMBER` sudah menggeser periode sendiri, memakai hari tutup yang dibaca dari tabel `TANGGAL_CLOSING`; (b) `KomitePostAdjustment` 14.6 menggeser lagi memakai angka `25` yang ditulis di rule; (c) 14.14 menambal hasilnya dengan `replaceAll` memakai ambang ketiga, `TglProd`. Karena (a) sudah benar, tambalan (c) hanya menemukan polanya — dan karena itu hanya mengubah nomor — **persis ketika ambang Pega dan ambang Oracle tidak sama**. K-06 menyebut "dua rumus"; sebenarnya tiga, dan yang ketiga menambal keluaran yang sudah benar. Yang berubah: §8.2 harus ditulis ulang dari sisi Oracle lebih dulu. Ditutup bersama: satu kueri isi `TANGGAL_CLOSING`.

**T-04 · OS akseptasi lama adalah buku besar tambah-saja tanpa kunci — jadi G-11 bukan soal kolom.**
Mengikuti: G-09 + G-11 + sumber prosedur. `PEGA_JSON_OS_AKSEP_KLAIMTNP` **selalu** `INSERT`; cabang `UPDATE`-nya dikomentari; variabel `id_count` dihitung lalu tidak pernah dipakai; dan seluruh nilai uang hanya hidup di dalam `DATA_JSON` (kolom `CLOB` ber-`CHECK … IS JSON`), sementara puluhan kolom skalar tabel itu tidak diisi prosedur ini. Artinya: akseptasi dan pembaliknya adalah **dua baris tambahan**, bukan dua keadaan satu baris. Yang berubah: §12.1 memperlakukan perbedaan ini sebagai kekurangan kolom; ia sesungguhnya perbedaan **jenis penyimpanan**. Ditutup bersama: dua kueri cacah di Lapis 2.

**T-05 · Empat penanda §8.3 bukan aturan uang.**
Mengikuti: G-15. Penanda `CNPAccNo*` hanya dibaca kembali di langkah 30 (`InputParam.CARI51–54`) dan sebagai parameter cetak yang — terbukti dari G-15 — dikirim kosong. Tidak satu pun nilai uang bergantung padanya. Yang berubah: §8.3 turun dari "aturan yang wajib pindah" menjadi aturan tampilan. Ditutup bersama: membuka satu PDF Acceptance Note (Lapis 4), dan itu pun opsional.

**T-06 · Pemetaan nama→jabatan adalah sumber kebenaran kedua, bukan sekadar hardcode.**
Mengikuti: K-03 + `TABLE_EMAILKOMITE.sql`. Tabel roster komite sudah memuat `NAME`, `EMAIL`, `OPERATOR_ID`, `DEGREE`, `LIMIT_BOTTOM`, `LIMIT_TOP`, `STS_AKTIF`, `TYPE_KOMITE` — yaitu persis jabatan dan ambang yang `SetKomiteList_Act` tulis ulang sebagai lima cabang nama orang. Yang berubah: K-03 berhenti menjadi keluhan kebersihan kode dan menjadi temuan **duplikasi sumber kebenaran**, yang menjelaskan mengapa jabatan di layar dan di kronologi bisa berbeda. Ditutup bersama: satu kueri isi `EMAILKOMITE`, yang juga menutup Q-1.

---

## 5. Dua puluh temuan, berapa akar

| Akar | Temuan yang bersumber darinya | Bentuk penutupan |
|---|---|---|
| **R-1 · Pengikatan langkah→rule tidak ikut terekspor** | G-01, G-02, G-13, Q-4, K-11, sebagian K-01 | Pencocokan nama halaman bind (gratis, menutup mayoritas) + ekspor ulang tiga rule |
| **R-2 · Gerbang dan muatan pembayaran Kasir belum dipahami** | G-03, G-04, K-10, K-12 | Satu kueri atas tabel log kiriman |
| **R-3 · Buku besar OS akseptasi: cakupan, penyaring, dan bentuk simpan** | G-09, G-11, G-08, G-06, T-04 | Dua kueri cacah + dua source prosedur yang belum dipegang |
| **R-4 · Periode buku ditetapkan berulang kali dengan ambang berbeda** | G-05, G-10, K-05, K-06, T-03 | Sudah tertutup dari berkas di tangan; sisa satu kueri isi tabel hari tutup |
| **R-5 · Kontrak tulis ke induk dan cakupan layar belum dipetakan tuntas** | G-07, G-12, G-14, G-15, G-17, G-18, G-19, G-20, T-01, T-05, K-03, K-04, K-08 | Tidak menunggu bukti; diselesaikan sebagai pekerjaan pemetaan di `to-spec` |

`to-tickets` berangkat dari tabel ini, bukan dari daftar dua puluh temuan — lima akar menghasilkan lima berkas kerja, bukan dua puluh tiket yang saling mengulang.

---

## 6. Putusan gerbang — pembelahan aliran

| Aliran | Isi | Temuan penopang | Putusan | Syarat pembukaan |
|---|---|---|---|---|
| **A-1 · Sirkulasi & jenjang persetujuan** | mesin keadaan, urutan jenjang, wewenang, penolakan, penutupan case | Flow, `When`, pra-syarat, `Property-Set` — **tak satu pun digugat** (T-02); K-01, K-03, K-04, K-08, G-17, G-18 | **BOLEH MULAI** | — |
| **A-2 · Layar persetujuan & aturan medan** | tujuh daftar, medan wajib, keadaan baca-saja, aturan tampil | Section XML utuh; G-14, G-20 | **BOLEH MULAI** | lampirkan daftar penuh 296 kendali sebelum tiket UI dibuat |
| **A-3 · Penomoran akseptasi & periode buku** | bentuk nomor, ambang hari tutup, larangan sunting ulang | G-05 (tertutup), K-05, K-06, T-03 | **BOLEH MULAI TERBATAS** | satu kueri isi tabel hari tutup sebelum ambangnya ditulis sebagai aturan |
| **A-4 · Pembayaran ke Kasir** | gerbang kirim, isi muatan, penjaga kiriman ganda, penanganan gagal | G-03, G-04, K-10, K-12 | **BEKU** | satu kueri atas log kiriman: apakah muatan terisi, baris mana yang terkirim, berapa kiriman per nomor akseptasi |
| **A-5 · OS akseptasi & pembalikan** | akseptasi, akseptasi bersyarat, rincian layer, pembalik | G-06, G-08, G-09, G-11, T-04 | **BEKU** | dua kueri cacah (`STS_REJECT`, `TYPE`) + source dua prosedur akseptasi yang belum dipegang |
| **A-6 · Dokumen cetak & surat** | PDF akseptasi/tutup/tolak, penerima surat | G-13 (`PostEmailKomiteCNP`), G-15, T-05 | **BEKU RINGAN** | daur hidup dokumen boleh dirancang sekarang; **isi** cetak dan daftar penerima menunggu ekspor satu rule |

**Pencabutan ruang lingkup — terbatas.** Keputusan 17 September 2026 (gap ditangguhkan sampai aplikasi jadi) **dicabut khusus untuk A-4 dan A-5**, dan hanya sebatas enam kueri baca-saja serta dua source prosedur. Alasannya satu kalimat: kedua aliran itu menentukan angka uang yang akan ditulis sistem baru, dan biaya penutupannya adalah satu sesi kueri — menangguhkannya berarti membayar dengan risiko untuk menghemat satu jam.

---

## 7. Rencana penutupan berjenjang

### Lapis 1 — berkas yang sudah dipegang, atau ekspor ulang

| Bukti yang diambil | Menutup |
|---|---|
| `pengetahuan/ddl/PROCEDURE_PROC_GENERATE_SEQUENCE_NUMBER.sql` — **sudah di repo** | Q-11 seluruhnya · G-05 · dasar K-05, K-06 · T-03 · membuka A-3 |
| `pengetahuan/ddl/PROCEDURE_PEGA_JSON_OS_AKSEP_KLAIMTNP.sql` — **sudah di repo** | bentuk simpan OS · T-04 · memperluas G-09 · memperkuat G-11 |
| `pengetahuan/ddl/TABLE_OS_AKSEPTASI_KLAIM.sql` — **sudah di repo** | arti `CARI20` (kolom `CLAIMOLD`) · domain `STS_REJECT` · keberadaan kolom `TYPE` · sebagian G-06 |
| `pengetahuan/ddl/TABLE_DIRECTTOKASIR_LOG.sql` + `TABLE_EMAILKOMITE.sql` — **sudah di repo** | bentuk log kiriman (menyiapkan kueri Lapis 2) · roster & ambang komite · T-06 · dasar Q-1 |
| Pencocokan nama halaman bind terhadap 16 rule SQL — **tanpa biaya, tanpa permintaan** | pengikatan 14 dari 22 langkah `RDB-List` dan seluruh 6 langkah `Connect-REST` · sebagian besar G-01 · T-02 |
| Ekspor ulang **tiga** rule: `PostEmailKomiteCNP`, `ASMForceCaseClose`, dan rule RDB yang mengisi halaman `TglProd` | G-13 (bagian P1) · Q-4 · sisa G-01 · membuka A-6 |
| Ekspor ulang dua activity inti **dengan parameter metode** bila opsi itu tersedia | G-02 · Q-9 |

### Lapis 2 — kueri baca-saja ke basis data berjalan (paket pemblokir, empat kueri)

| Kueri | Menutup |
|---|---|
| Satu baris `DATA_JSON` pada log kiriman kasir untuk akseptasi non-prop, **plus** cacah baris per nomor akseptasi | G-03 (arah penyaring terbukti dari hasilnya) · G-04 · K-10 (kiriman ganda) · K-12 · uji paritas P-08 · membuka A-4 |
| Cacah baris OS akseptasi dikelompokkan menurut `STS_REJECT` dan `TYPE` | G-09 (apakah penyaring `= 0` pernah kena) · sebagian G-06 · membuka A-5 |
| Baris OS untuk satu klaim ber-dua-adjustment yang salah satunya ditutup | G-09 (cakupan klaim vs adjustment) · uji paritas P-07 |
| Isi tabel hari tutup, **plus** sampel nomor akseptasi beserta tanggalnya | T-03 · ambang A-3 · uji paritas P-09/P-10/P-11 |

Satu kueri tambahan yang tidak memblokir: baris rincian layer untuk akseptasi bersyarat → menutup **G-08** dan bagian bersyarat **K-07** sekaligus.

### Lapis 3 — Tracer

**Tidak diperlukan.** Kelima usulan Tracer di `GRILL-01.md` diganti: Q-10 oleh kueri log kasir, Q-11 oleh berkas prosedur yang sudah dipegang, Q-12 dan Q-14 oleh kueri OS, Q-4 oleh ekspor ulang satu rule, Q-8 oleh penutupan K-11 (bagian 8). Satu-satunya keadaan yang mengembalikan Tracer: **bila kueri log kasir tidak menemukan satu pun baris untuk akseptasi non-prop.** Itu pun bukan alasan menjalankan Tracer lebih dulu — ketiadaan baris adalah jawaban tersendiri, dan jawaban itu mengubah G-03/G-04 dari dugaan menjadi kepastian.

---

## 8. Koreksi atas sidang K

**K-11 · DIPERLEMAH → DITUTUP, bukan temuan.**
Perilaku `Obj-Refresh-And-Lock` tidak perlu dipastikan: akibatnya sudah terlihat di data. Bila refresh membuang suntingan, nomor akseptasi dan keputusan komite yang ditulis pada langkah 6–20 tidak akan pernah mendarat di induk — padahal keduanya ada. Menyelidikinya lebih jauh berarti membayar untuk mengetahui apa yang sudah dibuktikan hasilnya.
· **Hentikan penyelidikan K-11, dan tetapkan batas transaksi sistem baru dari invarian bisnis, bukan dari perilaku Pega.**

**K-12 · SALAH KAPRAH → DIKUATKAN, dasarnya diganti.**
Interogator menggugurkan dasar K-12 karena §8.4 terbalik. Itu salah sasaran: risiko perangkaian JSON berdiri di atas **source Java yang terekspor utuh**, bukan di atas penyaring `UR`. Arah penyaring menentukan *baris mana* yang disusun, bukan *bagaimana* ia disusun, dan `catch(Exception)` yang melanjutkan tanpa melempar ulang terbaca langsung.
· **Susun muatan dengan pembentuk JSON yang menolak nilai kosong pada medan angka, dan jadikan kegagalan penyusunan sebagai kegagalan transaksi.**

**K-10 · DIKUATKAN → DIPERLUAS.**
Sistem lama **sudah** punya buku kirim: tabel log kiriman kasir menyimpan muatan, nomor akseptasi, dan keterangan hasil. Jadi gagasan mencatat niat kirim bukan barang baru bagi domain ini — yang kurang hanyalah urutannya (dicatat sesudah panggilan, bukan sebelum). Tabel itu sekaligus alat ukur: berapa kiriman ganda yang benar-benar terjadi selama ini.
· **Catat niat kirim sebelum memanggil, dan ukur dulu kiriman ganda yang sudah tercatat sebelum merancang penjaganya.**

**K-06 · DIKUATKAN → DIPERLUAS (lihat T-03).**
Bukan dua rumus, melainkan tiga ambang di dua lapis teknologi, dengan tambalan di Pega yang hanya menyala ketika lapisnya berselisih.
· **Tetapkan periode buku satu kali, di satu lapis, dan hapus penambalan nomor sesudah terbit.**

**K-03 · DIKUATKAN → DIPERLUAS (lihat T-06).**
Roster komite beserta jabatan dan ambangnya sudah ada sebagai data; pemetaan di rule adalah salinan kedua yang tidak lengkap.
· **Ambil peran dan ambang dari roster, dan salin nilainya ke keputusan pada saat keputusan diambil.**

**Diterima tanpa perubahan:** K-01, K-02, K-04, K-05, K-07, K-08, K-09.

---

## 9. Yang sudah boleh dibekukan

Sembilan keputusan berikut benar di **semua** kemungkinan hasil penyelidikan, dan dapat ditulis sebagai ADR hari ini tanpa menunggu satu kueri pun.

1. Wewenang memutus diperiksa dengan **kesetaraan identitas penuh**, ditegakkan di backend, dan permintaan yang gagal **ditolak**, bukan ditandai. *(K-01)*
2. Keputusan hanya dicatat untuk pihak yang benar-benar bertindak; jenjang yang tidak sempat memutus punya keadaan tersendiri yang bukan "tolak". *(K-04)*
3. Satu waktu keputusan per jenjang — satu fakta, satu kolom. *(K-08)*
4. Peran dan jabatan diambil dari roster sebagai data, lalu **disalin ke keputusan pada saat diambil** sehingga riwayat tidak berubah ketika orangnya pindah jabatan. *(K-03, T-06)*
5. Nilai lingkungan, kode akuntansi, alamat layanan, dan daftar penerima surat adalah **konfigurasi**; tidak ada cabang berdasarkan identitas orang di jalur mana pun. *(K-09)*
6. Niat memanggil layanan luar dicatat **sebelum** panggilan, dengan kunci idempotensi per efek. *(K-10)*
7. Nomor akseptasi diterbitkan **sekali**, dari nilai yang seluruhnya berpenulis, dan **tidak pernah disunting** sesudah terbit. *(K-05, G-05, T-03)*
8. Sirkulasi tanpa jenjang **ditolak saat dibuat**; tidak ada jalur yang menutup sirkulasi tanpa keputusan. *(G-17)*
9. Keadaan sirkulasi dan akibatnya pada klaim disimpan **terpisah**, sehingga "komite setuju" pada usulan tutup/tolak tidak terbaca sebagai "klaim disetujui". *(G-19, §6)*

**Dan yang belum boleh dibekukan:**

- **Unit pembalikan** — klaim atau adjustment: menunggu dua kueri OS (G-09).
- **Baris mana yang dibayar** ke Kasir, dan apakah baris retensi cedant termasuk: menunggu kueri log kiriman (G-03, G-04).
- **Bentuk penyimpanan akseptasi** — baris berkunci unik atau catatan tambah-saja: bergantung pada jawaban di atas (G-11, T-04).
- **Ambang periode buku**: menunggu isi tabel hari tutup (T-03).
- **Apakah rincian layer ditulis untuk akseptasi bersyarat**: menunggu satu kueri (G-08).
- **Penerima surat jenjang berikutnya**: menunggu ekspor satu rule (G-13).
