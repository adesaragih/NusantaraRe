# Daftar serah terima — Claim Life, PremiumList Life, Komite Claim Life

*Disusun 29 September 2026, GILIRAN-17 paket 6; diperbarui 30 September 2026, GILIRAN-18 (PL-09 dan N13 keluar, lihat §7).
Sumbernya lembar keputusan work owner ("rekomendasi"). Butir di bawah **tidak dapat diputuskan dari kode maupun XML**;
masing-masing menunggu penerimanya.*

Tiap butir berisi satu pertanyaan, buktinya (baris XML/berkas), dan akibatnya bila tidak dijawab. Rinciannya tetap di
tiket dan di `.scratch/claim-life/OQ-untuk-tim.md`.

## 1. DBA

| Butir | Pertanyaan | Bukti | Akibat bila tidak dijawab |
| --- | --- | --- | --- |
| **M2** | Apakah hilir membaca `LAPSE_DATE`/`CLAIM_RECEIVED_DATE`/`COMPLETE_DATE`/`CONFIRMATION_DATE` dari `OS_AKSEPTASI_KLAIM_LIFE`? Bila ya, dengan kunci apa baris cermin dicari tanpa nama tertanggung, dan apakah `LAPSE_DATE = DOL` disengaja? | `UpdateDateClaimLife_SQL` b85–90 (`WHERE CASEID AND NAME_OF_INSURED AND CERTIFICATE_NO`, `LAPSE_DATE = <DOL>`) | Cermin tanggal dialog Edit Date tidak ditiru. Hilir yang membaca tanggal itu dari tabel warisan melihat nilai lama |
| **M7** | Dari mana `OUTWARDRATEID` produk dibaca **tanpa** mengurai `M_PRODUCT_LIFE.JSONDATA` (AC 38)? | `SpreadingClaimLife_Act` 5.1 b1502/b1527–1528, `GetProductLife` b84–86, Java b1016. Pembaca `RATE_LIFE` sudah ada (`repository/ratelife.go`, GILIRAN-17) | Spreading tetap tidak dipanggil. Panel `RetroDetailClaimLife` dan `T_CLAIMLF_ADJUSTMENT_SPREADING` klaim baru tetap kosong |
| **K-04a** | Rentang nomor akseptasi Komite (`GENERATE_SEQUENCE_NUMBER`, per tahun) dan Claim Life (`ACCEPTATIONNOLIFE_SEQ`, global) menghasilkan bentuk yang sama. Apakah rentangnya dijamin tidak bertabrakan? | tiket 04a Komite, bab OQ-K-04a. Uji keunikan: nomor diperiksa di tabel datar dan `T_CLAIMLF_ADJUSTMENT.ACCEPTED_NO` | Tabrakan menahan keputusan akhir (409, transaksi batal) sampai ada kebijakan |
| **PL-17** | Berapa `pyLastReservedID` awalan `NBLF-` di `PC_DATA_UNIQUEID`? | migrasi `058` (`SEQ_WORK_POLIS START WITH 22374`); audit `JSON_POLIS` → 33 baris, maksimum 22373 (tiket 00 PremiumList) | Bila penghitung Pega lebih tinggi dari 22373, nomor `NBLF-` baru dapat bertabrakan dengan nomor lama |
| **OQ-001** | DDL produksi (kolom, tipe, presisi, nullability, PK/FK, index) `OS_AKSEPTASI_KLAIM_LIFE`, `JSON_KLAIM`, `M_LIFE_PREMIUM_DETAIL`, dan tabel klaim Life terkait | tidak ada DDL di korpus; skema uji memakai tipe longgar | Migrasi data produksi (tiket 13 Claim Life) tidak dapat dibuktikan |
| **OQ-002** | Kontrak `PROC_GENERATE_SEQUENCE_NUMBER` (format retro/non-retro, reset per tahun/lini/global) dan `GET_TOKEN_STORAGE` | badan procedure tidak ada di korpus (tiket 02, 12 Claim Life) | Penomoran ditiru dari bacaan SQL; format sebenarnya tidak terverifikasi |
| **OQ-013** | Apakah semua procedure `POOLDATA.*` commit sendiri, dan untuk apa identitas pengguna di dalamnya? | `{OperatorID.pyUserIdentifier}` dikirim sebagai parameter; `COMMIT` di blok PL/SQL | Batas transaksi migrasi (tiket 13) tidak dapat dipastikan |
| **OQ-018** | Konfirmasi lingkungan `pxHostId` `pega-nusre` dan satu id-hash | jboss1073 = produksi, jboss117 = dev (sudah dijawab) | Pemetaan lingkungan saat cutover belum lengkap |
| **OQ-047** | Daftar endpoint aktual di `M_LINK_SERVICE` beserta strukturnya | alamat integrasi dibaca dari tabel itu (ADR-0013) | Resolusi alamat efek keluar tidak dapat diuji terhadap isi nyata |
| **G1** | Mohon **skema uji Oracle kosong** (`ORACLE_DSN`/`ORACLE_SCHEMA`) untuk uji tag `db` | 59 uji `db` DILEWATI di setiap angka GILIRAN-17 (M1/M5/M6/N2/K-05), **61** sesudah GILIRAN-18 (tambah PL-09 dan N13) | Jalur yang hanya terbukti terhadap Oracle (migrasi 021/022, `INSERT … SELECT`/`UPDATE … SELECT` cermin, 5.1 Komite) tetap tanpa bukti Oracle |

## 2. Pemilik ekspor Pega

| Butir | Pertanyaan | Bukti | Akibat bila tidak dijawab |
| --- | --- | --- | --- |
| **N11** | Rule mana (Declare Expression, activity kelas lain) yang mengisi `.CLAIM_GROSS` baris adjustment? Mohon ekspor rule itu pada kelas `Int-LIFE_PREMIUM_DETAIL`/`AdjustmentList` | `AdjustmentDetail_Section` b2961/b2970 (Read-only), b2977/b3026 (wajib); pembaca 4, penulis 0 di korpus Claim Life | `CLAIM_GROSS` = `CLAIM_AMOUNT` **sementara** (OQ-N12 (a)), dan gerbang Save to RNM 11.17.1 praktis tidak pernah menolak |
| **N4** | Pemberitahuan tiga residu yang tidak ditiru: (a) `.Protect` nol penulis, (b) `ADJUSTMENT_DATE`/`PrintFaceClaim` tanpa kolom, (c) `IsAccept` tanpa kolom (dipakai `IS_CHECK`). Apakah ada rule di luar ekspor yang mengisinya? | `SaveOutStandingLife_Act` langkah 5, 22.1.3.2, b1181 | Ketiganya tetap tidak ditiru; pesan `.Protect` tidak pernah muncul |

## 3. Pemilik Arasapas

| Butir | Pertanyaan | Bukti | Akibat bila tidak dijawab |
| --- | --- | --- | --- |
| **PL-11** | `convertJsonNusareToProduction` hanya mengirim `noPolis`/`caseId`/`tglInput`. Dari mana layanan membaca polis sesudah `JSON_POLIS` tidak lagi ditulis (pl1)? Dan dari mana `tglInput` (`.pxCreateDateTime`), yang tidak berkolom? | `ConnectREST/ConvertJsonNusareToProduction.xml` b200/b209/b221; `serviceInsertArasapasLife_act` langkah 5 b963 | Pelaksana outbox tetap **stub** (`PelaksanaPremiumList`, GILIRAN-17). Tak satu polis pun terkirim ke produksi |

## 4. Product + Underwriting / Finance

| Butir | Pertanyaan | Bukti | Akibat bila tidak dijawab |
| --- | --- | --- | --- |
| **OQ-032** | Apa yang menentukan jumlah tingkat komite (`KomiteLoop`) sebuah klaim Life? | sumbernya tidak ada di korpus (tiket 10 Claim Life) | Tangga komite dibentuk dari roster tiruan/konfigurasi, bukan aturan bisnis |
| **OQ-037** | Aturan ambang nominal roster komite Life: mata uang IDR? dapat dikonfigurasi? | ambang ter-hardcode di modul lain (mis. 30 jt/50 jt) | Komposisi roster tidak dapat diverifikasi |
| **OQ-060** | Apakah satu klaim Life selalu satu mata uang, atau boleh campur antarbaris? | `CURRENCY` di tingkat baris, disalin dari baris 1 | Bentuk tipe uang rekam akseptasi tetap per baris |

## 5. Work owner — baru dari GILIRAN-17

*Kosong sejak GILIRAN-18.* N13 diputuskan dari bukti XML b176 (lihat §7). Work owner tetap dapat **memveto**; bila diveto,
yang dibalik satu pemanggilan di `simpanrnm.go` langkah 22.

## 6. Menunggu jawaban lain

| Butir | Pertanyaan | Menunggu | Akibat bila tidak dijawab |
| --- | --- | --- | --- |
| **M3** | Migrasi dua kolom `PCT_CLAIM`/`CLAIM_PAID` (tipe, presisi) dan rute sunting adjustment untuk `CountClaimAmountLife_Act` | **N11** (`CLAIM_GROSS`): rumusnya `CLAIM_PAID = CLAIM_GROSS × PCTClaim/100` (b702/b723) | Claim Paid / `Percent Claim (%)` tidak tampil; aturan murni `models/klaimbayar.go` tetap tanpa sambungan |

---

*Migrasi `021_kolom_komentar_jejak.sql` dan `022_kolom_sts_hapus_peserta.sql` (+ `_down`) dijalankan **work owner** sesudah
giliran ini; executor tidak menjalankan `-migrate` dan tidak menyentuh Oracle.*

⛔ **Urutan deploy (temuan /code-review):** binari dari commit GILIRAN-17 **menuntut** kedua kolom itu. Sebelum `021` terpasang,
setiap penulisan jejak (8 kolom) gagal ORA-00904; sebelum `022`, setiap pembacaan peserta Claim Life gagal ORA-00904. Migrasi
tidak berjalan saat binari mulai, jadi jalankan `-migrate` (`make migrate`) **lebih dulu**, baru jalankan binarinya.

## 7. Keluar dari daftar ini — GILIRAN-18 (30 September 2026)

| Butir | Dari | Ditutup oleh | Commit |
| --- | --- | --- | --- |
| **PL-09** | §1 DBA | Asisten membaca badan `PEGA_M_LIFE_PREMIUM_SUMMARY` dari `ALL_SOURCE` DEV (`.scratch/premiumlist-life/dba-procedure-PEGA_M_LIFE_PREMIUM_SUMMARY.md`, `526fc93`). Isinya ditiru di Go: 37 kolom + `ID` dari sequence, dalam transaksi simpan summary yang sama. Tipe kolom dan index tabel itu tetap `[belum terverifikasi]`; keduanya dicatat di tiket 05a, bukan butir DBA baru | `0021df1` |
| **N13** | §5 work owner | Ikut XML b176 `[keputusan asisten dari bukti; veto work owner]`: Save to RNM menyetel `STS_REJECT = '0'` cermin. Larangan baca-saja dicabut untuk kolom itu **saja**; `ACCEPTATION_DATE` (b175) tidak ikut (tiket 03) | `1b51183` |
