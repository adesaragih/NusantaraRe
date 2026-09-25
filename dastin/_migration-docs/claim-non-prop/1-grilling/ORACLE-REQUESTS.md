# ORACLE REQUESTS — Claim Non Prop

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

**REQ: aktif 35 · selesai 1 · mati 1**

> Selisih 19 September 2026: **REQ-004 dicabut** menjadi `DEAD` — satu-satunya yang melayaninya (B6) ditutup sebagai *di luar kepemilikan*. Tujuh REQ lain turun menjadi **verifikasi**; keduanya tidak mengubah cacah aktif kecuali lewat REQ-004.

> Pencacah ini menghitung **REQ** — permintaan data ke DBA. Pertanyaan dan tiket punya pencacahnya sendiri dengan format yang sama; bendanya diberi label supaya tidak tertukar.

Register permintaan objek Oracle. Status: **OPEN** → **SENT** → **ANSWERED** → **DEAD** (objek tidak ada).
Nomor REQ tidak pernah dipakai ulang. Selama ada REQ berstatus OPEN/SENT, frontier **tidak boleh dinyatakan kosong**.

| REQ | Objek | Jenis | Confidence | Prioritas | Status | Tanggal |
|---|---|---|---|---|---|---|
| REQ-001 | Isi definisi `Rule-Obj-Property` di schema PegaRULES (`PR4_*`) | TABLE | **AKSES BELUM ADA** — struktur standar bawaan produk, bukan misteri | **VERIFIKASI** — turun dari BLOCKER 19 Sep 2026 (ADR-0003 bagian penutupan B1/B2) | OPEN | 2026-09-17 |
| REQ-002 | `ALL_TAB_COLUMNS` untuk 44 tabel bisnis | TABLE | CONFIRMED + DERIVED | **BLOCKER** | OPEN | 2026-09-17 |
| REQ-004 | `EMAILKOMITE` — struktur + isi konfigurasi jenjang | TABLE + DATA-PROFILE | DERIVED | **DICABUT** — lihat catatan di bawah tabel | DEAD | 2026-09-17 |
| REQ-005 | Profil presisi kolom uang & kurs — **dipersempit: hanya tabel bisnis `POOLDATA`**. Untuk `.ClaimData.*` jalur ini **tertutup permanen** (IN-BLOB terbukti) | DATA-PROFILE | CONFIRMED | PENTING | OPEN | 2026-09-18 |
| REQ-006 | Statistik pemakaian jalur tutup-langsung (mengambil alih A5) | DATA-PROFILE | CONFIRMED | PENTING | OPEN | 2026-09-17 |
| REQ-007 | Verifikasi keberadaan 8 objek DERIVED | TABLE/VIEW/SYNONYM | DERIVED | PENTING | OPEN | 2026-09-17 |
| REQ-008 | Tabel work Pega untuk `ASM-FW-GCNMFW-Work` | TABLE | GUESS | PENTING | OPEN | 2026-09-17 |
| REQ-009 | Struktur JSON di `json_klaim.DATA_JSON` & `OS_AKSEPTASI_KLAIM.DATA_JSON` | DATA-PROFILE | CONFIRMED | **VERIFIKASI** — turun dari PENTING 19 Sep 2026 (SPEC §21.8) | OPEN | 2026-09-17 |
| REQ-010 | Constraint, index, dan volume tabel inti | INDEX/CONSTRAINT | CONFIRMED | PENTING — turun dari penahan index 19 Sep 2026 (SPEC §21.7) | OPEN | 2026-09-17 |
| REQ-011 | Kuantifikasi Adjustment non-IDR vs ambang kewenangan Komite | DATA-PROFILE | CONFIRMED | **VERIFIKASI** — turun dari BLOCKER 19 Sep 2026 (BLUEPRINT §8.7 syarat 3) | OPEN | 2026-09-18 |
| REQ-012 | Isi `Data-Admin-DB-Table` untuk class `ASM-FW-%` | TABLE | **AKSES BELUM ADA** — struktur standar bawaan produk | **BLOCKER** | OPEN | 2026-09-18 |
| REQ-013 | Cacah pola teks pada `.CommentSuggest` | DATA-PROFILE | DERIVED | PENTING | OPEN | 2026-09-18 |
| REQ-014 | Cacah klaim dengan lebih dari satu baris `TreatyType='UR'` | DATA-PROFILE | DERIVED | PENTING | OPEN | 2026-09-18 |
| REQ-015 | Status buka/tutup 8 klaim bertambalan | DATA-PROFILE | DERIVED | PENTING | OPEN | 2026-09-18 |
| REQ-016 | Profil case buatan `VINCENTVERNANDO_1` (mengambil alih A10a) — cacah, rentang tanggal, nilai, jumlah persetujuan, berapa yang tidak sampai ke Arasapas, **dan apakah 59 kolom skalar `OS_AKSEPTASI_KLAIM` kosong untuk case-case itu** (menguji ADR-0023) | DATA-PROFILE | DERIVED | **BLOCKER** | OPEN | 2026-09-18 |
| REQ-017 | DDL **14 objek** yang dirujuk dari procedure/view tetapi tidak disertakan. Daftar lengkapnya kini terbaca (`SAPUAN-S10-S16.md` §2): `M_CURRENCYSTANDARD`, `M_CURRENCY`, `M_NATION`, `M_PROVINCE`, `M_CAUSE_OF_LOSS`, `D_CAUSE_OF_LOSS`, `M_SITE_DATABASE`, `MST_USER_TEKNIS`, `OS_AKSEPTASI_SUBJECTIVITY`, `TANGGAL_CLOSING`, `GCP_IMAGE`, `BRANCH`, `CITYINPUT`, `DISTRICTINPUT`, `RWINPUT`. **Workbook kedua menutup NOL di antaranya** — `CURRENCYSTANDARD` adalah view yang membaca `m_CurrencyStandard`, objek berbeda (D19) | TABLE/VIEW | CONFIRMED | **BLOCKER** (untuk `M_CURRENCYSTANDARD`) | OPEN | 2026-09-18 |
| REQ-018 | Cacah `PYID` yang muncul lebih dari sekali di tabel work, **dan** cacah baris ganda `(CASEID, TypeLoss, Currency)` di `OS_AKSEPTASI_KLAIM` — satu REQ, dua tabel, pola yang sama | DATA-PROFILE | DERIVED | **BLOCKER** | OPEN | 2026-09-18 |
| REQ-019 | Cacah `FLAGONGOINGCOMMITTE` menyala saat klaim ditutup, sebaran `KOMITECOUNT`, **dan klaim tertutup dipecah menurut ada-tidaknya Adjustment** — yang terakhir menguji hipotesis I2 | DATA-PROFILE | CONFIRMED | **VERIFIKASI** — turun 19 Sep 2026; ramalan I2 memverifikasi, tidak menahan | OPEN | 2026-09-18 |
| REQ-020 | Berapa klaim melewati penjaga `CheckDateDOL_Act` dengan `DateOfLoss` di luar masa treaty (FINDING-005). **Keterjangkauan sudah diperiksa: dapat dijalankan dengan SQL biasa** — `TO_DATE(DATEOFLOSS,'YYYYMMDD')` vs `TREATYINDETAIL.TERMINATION` (`DATE`). Tidak perlu bongkar blob | DATA-PROFILE | CONFIRMED | **PRASYARAT CUTOVER** | OPEN | 2026-09-18 |
| REQ-021 | Grant dari `DATAPEGA` ke `POOLDATA`, dan adakah source `POOLDATA` yang menulis ke `DATAPEGA.PC_*` | METADATA | DERIVED | **BLOCKER TEKNIS + BLOCKER PERKIRAAN BIAYA** | OPEN | 2026-09-18 |
| REQ-022 | Mata uang **bukan IDR** yang kursnya tersimpan `1` — mengukur jalur `NO_DATA_FOUND` pada `GETCURRENCYSTANDARD` | DATA-PROFILE | CONFIRMED | PENTING | OPEN | 2026-09-18 |
| REQ-023 | Siapa menulis dan siapa membaca ~60 kolom skalar `OS_AKSEPTASI_KLAIM` yang tidak diisi `PEGA_JSON_OS_AKSEP_KLAIMTNP` | METADATA | DERIVED | **BLOCKER** | OPEN | 2026-09-18 |
| REQ-024 | `PYREOPENCOUNT > 0` dan `PYREOPENTIMESTAMP` — apakah klaim pernah di-reopen (menjawab A4, mengonfirmasi atau membatalkan ADR-0002) | DATA-PROFILE | CONFIRMED | **BLOCKER** | OPEN | 2026-09-18 |
| REQ-025 | `PXAPPLICATION`, `PXAPPLICATIONVERSION` dikelompokkan dengan `PXOBJCLASS` — aplikasi/ruleset lain yang menyentuh class ini (menjawab A9) | DATA-PROFILE | CONFIRMED | PENTING | OPEN | 2026-09-18 |
| REQ-026 | `PXCREATEOPERATOR` / `PXUPDATEOPERATOR` + `PXUPDATEDATETIME` untuk kelima nama di FINDING-002 — aktivitas terakhir tiap orang (menjawab A11a) | DATA-PROFILE | CONFIRMED | PENTING | OPEN | 2026-09-18 |
| REQ-027 | Sebaran nilai `OS_AKSEPTASI_KLAIM.PAYMENTTYPE` — nilai mana yang benar-benar dipakai (menjawab A6a) | DATA-PROFILE | CONFIRMED | PENTING | OPEN | 2026-09-18 |
| REQ-028 | Cacah baris berkode setoran `"100115"` di `DIRECTTOKASIR_LOG` atau tabel kasir — apakah klaim non-prop syariah pernah ada (menutup G1) | DATA-PROFILE | CONFIRMED | **VERIFIKASI** — turun dari PENTING 19 Sep 2026 (SPEC §21.3) | OPEN | 2026-09-18 |
| REQ-029 | Akseptasi ber-`STS_SUBJECTIVITY='1'` yang **sudah punya catatan pembayaran** — apakah uang keluar atas syarat yang belum terpenuhi | DATA-PROFILE | CONFIRMED | **BLOCKER** | OPEN | 2026-09-18 |
| REQ-030 | Adakah baris `json_polis` ber-`BusinessFac='F'` yang dirujuk klaim non-prop (menguji ADR-0020) | DATA-PROFILE | CONFIRMED | PENTING | OPEN | 2026-09-18 |
| REQ-031 | **Cacah baris per entitas yang dimiliki** (ADR-0026): klaim, akseptasi, alokasi, adjustment — beserta sebarannya per tahun dan per status buka/tutup. Menentukan volume migrasi ADR-0018 dan ukuran tabel korelasi | DATA-PROFILE | CONFIRMED | PENTING | OPEN | 2026-09-18 |
| REQ-032 | **Versi instance tujuan** — dua nilai saja: `SELECT version_full FROM v$instance` (atau `v$version`), dan `SHOW PARAMETER COMPATIBLE`. ~~Menentukan batas panjang pengenal~~ — **tidak lagi**: batas 30 byte ditetapkan sebagai aturan tetap 19 Sep 2026, dan tabel singkatan tertutup **dibatalkan**. Yang tersisa: memastikan tidak ada kejutan saat DDL dijalankan, dan mengetahui fitur mana yang tersedia saat pemasangan | VERSI | CONFIRMED | **VERIFIKASI** — turun dari BLOCKER 19 Sep 2026: batas 30 byte dipilih sebagai aturan tetap, sah di setiap versi | OPEN | 2026-09-18 |
| REQ-033 | **Sebaran `LAYERPART` pada struktur treaty** — `SELECT ... GROUP BY ID, LAYER, LAYERTYPE HAVING COUNT(DISTINCT LAYERPART) > 1` atas `POOLDATA.PROPORTIONALARRG` **dan** `POOLDATA.TREATYINDETAIL`. Menjawab apakah kunci halus (klaim, layer, mata uang) pernah melahirkan dua baris di tempat sistem lama melihat satu. Read-only, tanpa data nasabah, **tanpa membongkar blob** — keempat kolomnya nyata. Tidak dititipkan ke REQ-018: `LAYERPART` tidak pernah sampai ke `OS_AKSEPTASI_KLAIM`. **Menyempit 18 September 2026**: struktur kedua tabel sudah terbaca (`TABLE_PROPORTIONALARRG.sql`, `TABLE_TREATYINDETAIL.sql`), jadi yang tersisa murni **sebaran**, bukan struktur | DATA-PROFILE | CONFIRMED | **BLOCKER** (untuk pemetaan view kompatibilitas) | OPEN | 2026-09-18 |
| REQ-034 | **Menguji ramalan FINDING-008 (usulan).** Cacah baris `SpreadingRisk` yang `Currency`-nya **bukan** `IDR` **dan** `ClaimAmountIDR`-nya **sama persis** dengan nilai mata uang aslinya — yaitu baris yang kolom rupiahnya tidak pernah dikonversi. Sertakan sebaran `KursIDR` pada baris-baris itu. Nilainya ada di dalam `PZPVSTREAM`, jadi ia **menumpang REQ-001**; bila REQ-001 gagal, cadangannya membongkar blob. Read-only, tanpa data nasabah, maksimum 20 baris contoh dengan kolom teks bebas ditutup | DATA-PROFILE | DERIVED | **VERIFIKASI** — turun dari PENTING 19 Sep 2026 (FINDING-008 bagian 5) | OPEN | 2026-09-18 |
| REQ-035 | **Isi 42 baris `Rule-Declare-DecisionTable GetMimeType`** pada class `ASM-FW-GISFW-Int-T_STORAGE_IMAGE`. Ekspor memuat 42 slot baris yang **nilainya menutup sendiri** — kerangkanya ada, pemetaannya tidak. Tabel ini memetakan ekstensi berkas ke tipe MIME. **Celah bahan, bukan celah metode**: tidak ada pola sapuan yang dapat memulihkannya (D45) | TABLE | CONFIRMED | PENTING (menahan tiket `20`, dokumen klaim) | OPEN | 2026-09-19 |
| REQ-036 | **Perilaku fungsi bawaan Pega `@month()`** — berbasis nol atau berbasis satu. `HitServiceToKasir_Act` menulis `CARI20 = @toDecimal(@month(CARIDATETIME)) + 1`. Bila berbasis nol penambahan itu benar; bila berbasis satu ia menghasilkan bulan ke-13 pada `TglBolehBayar` yang dikirim ke kasir. Perilakunya **tidak ada di ekspor**. Cukup satu contoh keluaran, atau rujukan dokumentasi versi Pega yang dipakai | VERSI | DERIVED | PENTING (tanggal boleh bayar di muatan kasir) | OPEN | 2026-09-19 |
| REQ-037 | **Seluruh jalur tulis yang melewati hak objek — enam jalur, bukan satu.** (1) Hak sistem berakhiran `ANY`: `SELECT/INSERT/UPDATE/DELETE/ALTER ANY TABLE`, `ANY INDEX`. (2) **Prosedur definer's rights** milik akun berhak tulis — siapa pun yang memegang `EXECUTE` menulis sebagai pemiliknya; ini **pola yang sudah dipakai di instance ini**, bukan kemungkinan teoretis: `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP`. (3) **Hibah ke `PUBLIC`** — tidak muncul saat memeriksa akun satu per satu. (4) **`CREATE ANY TRIGGER`** — menaruh penulis di dalam tabel kita sendiri. (5) **Peran yang memuat hak `ANY`**: `DBA`, `IMP_FULL_DATABASE`, `DATAPUMP_IMP_FULL_DATABASE`, dan peran buatan lokal — periksa **isi peran** lewat `role_sys_privs`, bukan hanya hibah langsung di `dba_sys_privs`. (6) **`GRANT ANY OBJECT PRIVILEGE` dan `GRANT ANY PRIVILEGE`.** Ditambah dua nilai untuk `Z01_PENGAWASAN_TULIS.sql`: apakah **Unified Auditing aktif** (`v$option`), dan siapa pemegang `AUDIT_ADMIN`. Read-only, metadata saja, tanpa data nasabah | METADATA | DERIVED | **BLOCKER** (menentukan apakah ADR-0017 berlaku sungguhan) | OPEN | 2026-09-19 |
| | > **REQ ini mengukur satu titik waktu; klaim ADR-0017 tentang keadaan yang bertahan.** Jalur (6) yang membuat pembedaan itu perlu ditulis: pemegang `GRANT ANY PRIVILEGE` atau `GRANT ANY OBJECT PRIVILEGE` **dapat memulihkan jalur tulis besok**, sesudah REQ ini dijawab bersih hari ini. Karena itu jawabannya tidak pernah menutup ADR-0017 sendirian — ia harus berpasangan dengan pengawasan yang berjalan terus (`Z01_PENGAWASAN_TULIS.sql`), dan itulah sebabnya rumusan klaim ADR-0017 diubah dari *mencegah* menjadi *tidak meninggalkan jejak*. | | | | |

## Perubahan kedudukan — 19 September 2026, sesudah penutupan register

**Tujuh REQ turun menjadi verifikasi.** Yang berubah bukan isinya melainkan **apa yang bergantung padanya**: butir register yang dulu menunggu jawabannya sudah ditutup oleh keputusan rancangan, sehingga jawaban REQ kini **menguatkan atau menggugurkan**, tidak lagi **menahan**.

| REQ | Dulu menahan | Yang menutupnya | Yang tersisa untuk REQ |
|---|---|---|---|
| REQ-032 | penamaan seluruh skema; tiket `01`, `02`, `03`, `04` | keputusan 19 Sep 2026: **batas 30 byte berlaku selamanya**, bukan sebagai pengamanan sementara — 30 byte sah di setiap versi Oracle, 128 hanya di sebagian | memastikan tidak ada kejutan saat DDL dijalankan, dan mengetahui fitur mana yang tersedia saat pemasangan |
| REQ-001 | B1, B2, F (uji `IsEditClaim`) | ADR-0003 — tidak ada presisi yang dapat dipulihkan, jadi presisi adalah keputusan maju | memverifikasi isi blob; skala 20 menampung apa pun yang ditemukan |
| REQ-009 | B4 | SPEC §21.8 — muatan mendarat utuh, yang dikenali naik, sisanya tercatat | memperbanyak daftar medan yang dikenali |
| REQ-010 | B7 | SPEC §21.7 — index ditentukan dari pola akses nyata pasca-migrasi | ukuran tabel untuk perkiraan kapasitas |
| REQ-011 | B8, FINDING-001 | BLUEPRINT §8.7 syarat 3 — nilai tanpa mata uang ditolak di batas | menutup FINDING-001 sebagai catatan atas sistem lama |
| REQ-019 | I1/I2 | F/H/I ditutup sebagai hipotesis permanen; rancangan memenuhi kedua cabang | menguji ramalan I2 |
| REQ-028 | G1, A15 | SPEC §21.3 — penanda lini usaha jadi atribut data apa pun jawabannya | apakah klaim non-prop syariah pernah ada |
| REQ-034 | E14 / FINDING-008 | ADR-0029 — nilai rupiah tidak dapat lahir tanpa kurs | berapa baris lama yang terdampak (AK-2a, AK-2b) |

**REQ-004 dicabut (status `DEAD`).** Ia hanya melayani **B6** — arti `EMAILKOMITE.DEGREE` dan `LIMIT_TOP`. `EMAILKOMITE` adalah tabel milik **modul Komite** (ADR-0026, batas kepemilikan mengikuti nama class): tidak dimigrasi, tidak dirancang, tidak ditanyakan. Dicabut karena tidak ada lagi yang menunggunya, bukan karena terjawab. Nomornya tidak dipakai ulang.

> `BLUEPRINT.md` §8.6 tetap memuat kolom `EMAILKOMITE` beserta filternya sebagai **gambaran sistem lama**. Mengetahui bentuknya tidak sama dengan memilikinya.

**REQ-032 turun jadi VERIFIKASI — 19 September 2026.** Ia tidak lagi menahan apa pun.

Sebabnya: penamaan diputuskan pada **batas 30 byte sebagai aturan tetap**, bukan sebagai pengamanan sementara sampai `COMPATIBLE` diketahui. Tiga puluh byte sah di **setiap** versi Oracle; seratus dua puluh delapan hanya di sebagian. Memilih 30 membuat skema ini benar tanpa bergantung pada nilai yang tidak kita pegang — dan **pertanyaan yang jawabannya tidak mengubah apa pun bukan pertanyaan yang perlu ditunggu**.

**Tabel singkatan tertutup dibatalkan, bukan ditunda.** Tekanan panjang selama ini datang dari satu kebiasaan, bukan dari bahasa: mencantumkan daftar kolom di nama constraint. Nama tabel dan kolom sendiri tidak pernah menabrak batas. Aturan penamaan final ada di `SPEC-MODEL-DATA.md` bagian 16.

**Tidak ada lagi REQ yang menahan keputusan rancangan.** Empat REQ menahan **pelaksanaan** tiket tertentu — dan itu jenis penahan yang berbeda:

| REQ | Menahan pelaksanaan | Tiket |
|---|---|---|
| **REQ-018** | cacah `PYID` ganda dan baris ganda akseptasi | `15`, dan penanganan migrasi di `12`, `18`, `29`, `34` |
| **REQ-033** | sebaran `LAYERPART` | pemetaan view kompatibilitas — `29`, `32` |
| **REQ-021** | grant `DATAPEGA` → `POOLDATA` | `35`, `39` |
| **REQ-037** | enam jalur tulis yang melewati hak objek | klaim ADR-0017 di `35` |

Sisanya memverifikasi.

**REQ-037 tetap BLOCKER untuk klaim ADR-0017**, bukan untuk pekerjaan model data. Ia menentukan apakah "satu pintu tulis" berlaku sungguhan di instance tujuan; skema tetap dapat dirancang dan dibaca tanpanya.

Seluruh query siap jalan ada di [`pengetahuan/RECON.sql`](./pengetahuan/RECON.sql), dipetakan per bagian:

| REQ | Bagian pengetahuan/RECON.sql |
|---|---|
| REQ-001 | BAGIAN 1 (1A–1D, tiga varian) |
| REQ-002 | BAGIAN 3 |
| REQ-003 | BAGIAN 5 (5A, 5B, 5C) |
| REQ-004 | BAGIAN 3 + BAGIAN 7B |
| REQ-005 | BAGIAN 6 (6A–6D) |
| REQ-006 | BAGIAN 8 |
| REQ-007 | BAGIAN 2 (2A, 2B) |
| REQ-008 | BAGIAN 2C |
| REQ-009 | BAGIAN 9 |
| REQ-010 | BAGIAN 4 + BAGIAN 7 |
| REQ-011 | BAGIAN 10 |
| REQ-012 | BAGIAN 11 — *belum ditulis, menunggu T5a* |
| REQ-013 | BAGIAN 12 — *belum ditulis, menunggu T2* |
| REQ-014 | BAGIAN 13 — *belum ditulis, menunggu T2* |
| REQ-015 | BAGIAN 14 — *belum ditulis* |
| REQ-016 | BAGIAN 15 — *belum ditulis* |
| REQ-017 | BAGIAN 3 (perluasan daftar objek) |
| REQ-018 | BAGIAN 16 — *belum ditulis* |
| REQ-019 | BAGIAN 17 — *belum ditulis* |
| REQ-020 | BAGIAN 18 — *belum ditulis* |
| REQ-021 | BAGIAN 19 — *belum ditulis* |
| REQ-022 | BAGIAN 20 — *belum ditulis* |
| REQ-023 | BAGIAN 21 — *belum ditulis* |
| REQ-024 | BAGIAN 22 — *belum ditulis, menunggu DDL tabel work* |
| REQ-025 | BAGIAN 22 — idem |
| REQ-026 | BAGIAN 22 — idem |
| REQ-027 | BAGIAN 23 — *belum ditulis* |
| REQ-028 | BAGIAN 24 — *belum ditulis* |
| REQ-029 | BAGIAN 25 — *belum ditulis* |
| REQ-030 | BAGIAN 26 — *belum ditulis* |
| REQ-031 | BAGIAN 27 — *belum ditulis* |
| REQ-032 | BAGIAN 28 — *belum ditulis* |
| REQ-033 | BAGIAN 29 — *belum ditulis* |

## Penomoran — tabrakan diselesaikan 18 September 2026

Karena jawaban Round 5 dan artefak Round 5 ditulis berbarengan, dua nomor bertabrakan. Diselesaikan tanpa menomori ulang apa pun yang sudah ditulis:

| Tetap | Dipindahkan |
|---|---|
| REQ-018 `PYID` ganda + duplikat akseptasi | rasio kurs ≈ 1 → **REQ-022** |
| REQ-019 `FLAGONGOINGCOMMITTE` / `KOMITECOUNT` | penulis & pembaca kolom skalar → **REQ-023** |
| REQ-020 penjaga `CheckDateDOL_Act` | — |
| `FINDING-005` bug perbandingan tanggal | `RETURN 1` pada kurs → **`FINDING-006`** |

## PREFLIGHT — didahulukan dari seluruh REQ

[`pengetahuan/PREFLIGHT.sql`](./pengetahuan/PREFLIGHT.sql) menjawab T1–T7 dari kamus data, tanpa menyentuh tabel bisnis.
Hasilnya menentukan bagaimana `pengetahuan/RECON.sql` ditulis ulang:

| Jawaban | Yang berubah |
|---|---|
| T1 versi Oracle | boleh/tidaknya `FETCH FIRST`, `search_condition_vc` |
| T2 tipe kolom JSON | BAGIAN 6 dan BAGIAN 10 dipakai apa adanya atau ditulis ulang dengan `JSON_VALUE` |
| T3 owner 25 nama tanpa prefix | kolom `OWNER` di `pengetahuan/PULL-LIST.csv`; yang AMBIGU tetap `?` |
| T4 hak akses | arti "tidak ditemukan": **TIDAK ADA** vs **TIDAK TERLIHAT** |
| T5 letak PegaRULES | REQ-001 dan REQ-012 diajukan ke DBA yang mana |
| T7 ukuran tabel | BAGIAN 6, 7, 10: sampling atau full scan |
| T8 (dijawab: tidak ada UAT) | seluruh pengetahuan/RECON.sql harus aman di produksi |

## Catatan penanganan

- Hasil ditempel ke [`pengetahuan/SCHEMA-ACTUAL.csv`](./pengetahuan/SCHEMA-ACTUAL.csv) mengikuti header yang sudah disediakan.
- Bila sebuah objek dilaporkan **tidak ada**, statusnya menjadi **DEAD** dan temuan itu dicatat di `_selesai/OPEN-QUESTIONS.md` — berarti ada rule yang merujuk objek mati.
- Bila hasil bertentangan dengan dugaan di `TABLE-EXTRACTION-REQUEST.md`, koreksinya ditulis terang-terangan di sini beserta apa yang berubah karenanya.

---

# SELESAI

| REQ | Isi | Jenis | Keyakinan | Prioritas | Status | Tanggal |
|---|---|---|---|---|---|---|
| ~~REQ-003~~ | Source 6 procedure + 1 function | PROCEDURE/FUNCTION | CONFIRMED | PENTING | **ANSWERED** — seluruhnya ada di `pengetahuan/DDL_Script_ClaimNonProp.xls`, tersimpan di `pengetahuan/ddl/` | 2026-09-18 |

`REQ-003` tertutup karena seluruh source procedure dan function diterima lewat `pengetahuan/DDL_Script_ClaimNonProp.xls` dan tersimpan di `pengetahuan/ddl/`. Isinya melahirkan `FINDING-006` (kurs `RETURN 1`), `ADR-0024` (kunci alami akseptasi), dan `ASK-AKUNTANSI.md` bagian 0.
