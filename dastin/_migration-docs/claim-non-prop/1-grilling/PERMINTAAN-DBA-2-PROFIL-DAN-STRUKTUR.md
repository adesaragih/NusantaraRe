# Permintaan ke DBA — 34 permintaan data, terurut menurut apa yang ditahannya

**Untuk**: DBA Oracle, instance tempat `POOLDATA` berjalan
**Dari**: tim migrasi Claim Non Prop
**Tanggal**: 19 September 2026

> **Tidak ada satu pun permintaan di kedua berkas ini yang menahan keputusan rancangan kami.** Empat di antaranya — **REQ-018**, **REQ-033**, **REQ-021**, **REQ-037** — menahan **pelaksanaan** tiket tertentu, dan itu sebabnya mereka di urutan atas. Sisanya memverifikasi.

> Pembedaan itu penting bagi yang menjawab: ia memberi tahu mana yang menunda orang bekerja, dan mana yang menajamkan apa yang sudah diputuskan.

> **Berkas ini bukan yang pertama.** `PERMINTAAN-DBA-1-VERSI-INSTANCE.md` memuat satu permintaan berisi dua nilai, dipisahkan karena ia paling murah dijawab dan paling banyak menahan penamaan. Kerjakan yang itu lebih dulu bila waktunya terbatas.

## Aturan yang berlaku untuk seluruh isi berkas ini

- **READ-ONLY.** Tidak ada `INSERT`, `UPDATE`, `DELETE`, `ALTER`, `DROP`, maupun `CREATE` di mana pun. Bila sebuah permintaan tampak menuntut perubahan, itu salah tulis kami — tolak dan beri tahu.
- **Tanpa data nasabah.** Yang diminta cacah, sebaran, struktur, dan metadata. Bukan isi baris.
- **Teks bebas dimasking.** Bila sebuah kolom teks perlu ikut, kirim panjangnya atau polanya, bukan isinya.
- **Maksimum 20 baris contoh** bila sebuah permintaan memang butuh contoh. Tidak pernah lebih.
- **Bentuk jawaban bebas.** CSV, tangkapan layar, atau tempelan teks sama-sama diterima. Yang kami butuhkan angkanya, bukan formatnya.
- Bila sebuah permintaan **tidak dapat dijalankan** — objeknya tidak ada, aksesnya tidak ada, biayanya terlalu besar — **katakan begitu**. Jawaban "tidak dapat" adalah jawaban, dan kami mencatatnya apa adanya. Yang tidak dapat kami pakai adalah diam.

## Urutan

Diurutkan menurut **berapa banyak pekerjaan yang tertahan**, bukan menurut nomor. Nomor `REQ-` adalah pengenal internal kami dan sengaja tidak berurutan di sini — tolong sebut nomornya saat menjawab supaya jawabannya tidak tertukar.

| # | REQ | Menahan | Jenis | Prioritas kami |
|---|---|---|---|---|
| 1 | **REQ-018** | 6 tiket | DATA-PROFILE | BLOCKER |
| 2 | **REQ-033** | 4 tiket | DATA-PROFILE | BLOCKER (untuk pemetaan view kompatibili |
| 3 | **REQ-021** | 2 tiket | METADATA | BLOCKER TEKNIS + BLOCKER PERKIRAAN BIAYA |
| 4 | **REQ-037** | 2 tiket | METADATA | BLOCKER (menentukan apakah ADR-0017 berl |
| 5 | **REQ-001** | 1 tiket | TABLE | VERIFIKASI — turun dari BLOCKER 19 Sep 2 |
| 6 | **REQ-020** | 1 tiket · temuan 005 | DATA-PROFILE | PRASYARAT CUTOVER |
| 7 | **REQ-015** | 1 tiket | DATA-PROFILE | PENTING |
| 8 | **REQ-027** | 1 tiket | DATA-PROFILE | PENTING |
| 9 | **REQ-031** | 1 tiket | DATA-PROFILE | PENTING |
| 10 | **REQ-036** | 1 tiket | VERSI | PENTING (tanggal boleh bayar di muatan k |
| 11 | **REQ-011** | temuan 001 | DATA-PROFILE | VERIFIKASI — turun dari BLOCKER 19 Sep 2 |
| 12 | **REQ-016** | temuan 002 | DATA-PROFILE | BLOCKER |
| 13 | **REQ-017** | temuan 001 | TABLE/VIEW | BLOCKER (untuk `M_CURRENCYSTANDARD`) |
| 14 | **REQ-002** | — | TABLE | BLOCKER |
| 15 | **REQ-012** | — | TABLE | BLOCKER |
| 16 | **REQ-023** | — | METADATA | BLOCKER |
| 17 | **REQ-024** | — | DATA-PROFILE | BLOCKER |
| 18 | **REQ-029** | — | DATA-PROFILE | BLOCKER |
| 19 | **REQ-013** | temuan 002 | DATA-PROFILE | PENTING |
| 20 | **REQ-022** | temuan 006 | DATA-PROFILE | PENTING |
| 21 | **REQ-034** | temuan 008 | DATA-PROFILE | VERIFIKASI — turun dari PENTING 19 Sep 2 |
| 22 | **REQ-005** | — | DATA-PROFILE | PENTING |
| 23 | **REQ-006** | — | DATA-PROFILE | PENTING |
| 24 | **REQ-007** | — | TABLE/VIEW/SYNONYM | PENTING |
| 25 | **REQ-008** | — | TABLE | PENTING |
| 26 | **REQ-009** | — | DATA-PROFILE | VERIFIKASI — turun dari PENTING 19 Sep 2 |
| 27 | **REQ-010** | — | INDEX/CONSTRAINT | PENTING — turun dari penahan index 19 Se |
| 28 | **REQ-014** | — | DATA-PROFILE | PENTING |
| 29 | **REQ-019** | — | DATA-PROFILE | VERIFIKASI — turun 19 Sep 2026; ramalan  |
| 30 | **REQ-025** | — | DATA-PROFILE | PENTING |
| 31 | **REQ-026** | — | DATA-PROFILE | PENTING |
| 32 | **REQ-028** | — | DATA-PROFILE | VERIFIKASI — turun dari PENTING 19 Sep 2 |
| 33 | **REQ-030** | — | DATA-PROFILE | PENTING |
| 34 | **REQ-035** | — | TABLE | PENTING (menahan tiket `20`, dokumen kla |

---

## 1. REQ-018

**Menahan**: tiket `10`, `12`, `15`, `18`, `29`, `34`

**Jenis**: DATA-PROFILE · **Prioritas kami**: BLOCKER

**Yang diminta**: Cacah `PYID` yang muncul lebih dari sekali di tabel work, **dan** cacah baris ganda `(CASEID, TypeLoss, Currency)` di `OS_AKSEPTASI_KLAIM` — satu REQ, dua tabel, pola yang sama

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 2. REQ-033

**Menahan**: tiket `10`, `17`, `29`, `32`

**Jenis**: DATA-PROFILE · **Prioritas kami**: BLOCKER (untuk pemetaan view kompatibilitas)

**Yang diminta**: **Sebaran `LAYERPART` pada struktur treaty** — `SELECT ... GROUP BY ID, LAYER, LAYERTYPE HAVING COUNT(DISTINCT LAYERPART) > 1` atas `POOLDATA.PROPORTIONALARRG` **dan** `POOLDATA.TREATYINDETAIL`. Menjawab apakah kunci halus (klaim, layer, mata uang) pernah melahirkan dua baris di tempat sistem lama melihat satu. Read-only, tanpa data nasabah, **tanpa membongkar blob** — keempat kolomnya nyata. Tidak dititipkan ke REQ-018: `LAYERPART` tidak pernah sampai ke `OS_AKSEPTASI_KLAIM`. **Menyempit 18 September 2026**: struktur kedua tabel sudah terbaca (`TABLE_PROPORTIONALARRG.sql`, `TABLE_TREATYINDETAIL.sql`), jadi yang tersisa murni **sebaran**, bukan struktur

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 3. REQ-021

**Menahan**: tiket `35`, `39`

**Jenis**: METADATA · **Prioritas kami**: BLOCKER TEKNIS + BLOCKER PERKIRAAN BIAYA

**Yang diminta**: Grant dari `DATAPEGA` ke `POOLDATA`, dan adakah source `POOLDATA` yang menulis ke `DATAPEGA.PC_*`

**Bentuk jawaban yang diterima**: daftar nama beserta atributnya

---

## 4. REQ-037

**Menahan**: tiket `01`, `35`

**Jenis**: METADATA · **Prioritas kami**: BLOCKER (menentukan apakah ADR-0017 berlaku sungguhan)

**Yang diminta**: **Seluruh jalur tulis yang melewati hak objek — enam jalur, bukan satu.** (1) Hak sistem berakhiran `ANY`: `SELECT/INSERT/UPDATE/DELETE/ALTER ANY TABLE`, `ANY INDEX`. (2) **Prosedur definer's rights** milik akun berhak tulis — siapa pun yang memegang `EXECUTE` menulis sebagai pemiliknya; ini **pola yang sudah dipakai di instance ini**, bukan kemungkinan teoretis: `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP`. (3) **Hibah ke `PUBLIC`** — tidak muncul saat memeriksa akun satu per satu. (4) **`CREATE ANY TRIGGER`** — menaruh penulis di dalam tabel kita sendiri. (5) **Peran yang memuat hak `ANY`**: `DBA`, `IMP_FULL_DATABASE`, `DATAPUMP_IMP_FULL_DATABASE`, dan peran buatan lokal — periksa **isi peran** lewat `role_sys_privs`, bukan hanya hibah langsung di `dba_sys_privs`. (6) **`GRANT ANY OBJECT PRIVILEGE` dan `GRANT ANY PRIVILEGE`.** Ditambah dua nilai untuk `Z01_PENGAWASAN_TULIS.sql`: apakah **Unified Auditing aktif** (`v$option`), dan siapa pemegang `AUDIT_ADMIN`. Read-only, metadata saja, tanpa data nasabah

**Bentuk jawaban yang diterima**: daftar nama beserta atributnya

---

## 5. REQ-001

**Menahan**: tiket `07`

**Jenis**: TABLE · **Prioritas kami**: VERIFIKASI — turun dari BLOCKER 19 Sep 2026 (ADR-0003 bagian penutupan B1/B2)

**Yang diminta**: Isi definisi `Rule-Obj-Property` di schema PegaRULES (`PR4_*`)

**Bentuk jawaban yang diterima**: daftar kolom beserta tipe, panjang, presisi, skala, dan nullable — satu baris per kolom

---

## 6. REQ-020

**Menahan**: tiket `24`; FINDING-005

**Jenis**: DATA-PROFILE · **Prioritas kami**: PRASYARAT CUTOVER

**Yang diminta**: Berapa klaim melewati penjaga `CheckDateDOL_Act` dengan `DateOfLoss` di luar masa treaty (FINDING-005). **Keterjangkauan sudah diperiksa: dapat dijalankan dengan SQL biasa** — `TO_DATE(DATEOFLOSS,'YYYYMMDD')` vs `TREATYINDETAIL.TERMINATION` (`DATE`). Tidak perlu bongkar blob

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 7. REQ-015

**Menahan**: tiket `34`

**Jenis**: DATA-PROFILE · **Prioritas kami**: PENTING

**Yang diminta**: Status buka/tutup 8 klaim bertambalan

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 8. REQ-027

**Menahan**: tiket `07`

**Jenis**: DATA-PROFILE · **Prioritas kami**: PENTING

**Yang diminta**: Sebaran nilai `OS_AKSEPTASI_KLAIM.PAYMENTTYPE` — nilai mana yang benar-benar dipakai (menjawab A6a)

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 9. REQ-031

**Menahan**: tiket `14`

**Jenis**: DATA-PROFILE · **Prioritas kami**: PENTING

**Yang diminta**: **Cacah baris per entitas yang dimiliki** (ADR-0026): klaim, akseptasi, alokasi, adjustment — beserta sebarannya per tahun dan per status buka/tutup. Menentukan volume migrasi ADR-0018 dan ukuran tabel korelasi

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 10. REQ-036

**Menahan**: tiket `05`

**Jenis**: VERSI · **Prioritas kami**: PENTING (tanggal boleh bayar di muatan kasir)

**Yang diminta**: **Perilaku fungsi bawaan Pega `@month()`** — berbasis nol atau berbasis satu. `HitServiceToKasir_Act` menulis `CARI20 = @toDecimal(@month(CARIDATETIME)) + 1`. Bila berbasis nol penambahan itu benar; bila berbasis satu ia menghasilkan bulan ke-13 pada `TglBolehBayar` yang dikirim ke kasir. Perilakunya **tidak ada di ekspor**. Cukup satu contoh keluaran, atau rujukan dokumentasi versi Pega yang dipakai

**Bentuk jawaban yang diterima**: dua baris teks disalin apa adanya

---

## 11. REQ-011

**Menahan**: FINDING-001

**Jenis**: DATA-PROFILE · **Prioritas kami**: VERIFIKASI — turun dari BLOCKER 19 Sep 2026 (BLUEPRINT §8.7 syarat 3)

**Yang diminta**: Kuantifikasi Adjustment non-IDR vs ambang kewenangan Komite

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 12. REQ-016

**Menahan**: FINDING-002

**Jenis**: DATA-PROFILE · **Prioritas kami**: BLOCKER

**Yang diminta**: Profil case buatan `VINCENTVERNANDO_1` (mengambil alih A10a) — cacah, rentang tanggal, nilai, jumlah persetujuan, berapa yang tidak sampai ke Arasapas, **dan apakah 59 kolom skalar `OS_AKSEPTASI_KLAIM` kosong untuk case-case itu** (menguji ADR-0023)

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 13. REQ-017

**Menahan**: FINDING-001

**Jenis**: TABLE/VIEW · **Prioritas kami**: BLOCKER (untuk `M_CURRENCYSTANDARD`)

**Yang diminta**: DDL **14 objek** yang dirujuk dari procedure/view tetapi tidak disertakan. Daftar lengkapnya kini terbaca (`SAPUAN-S10-S16.md` §2): `M_CURRENCYSTANDARD`, `M_CURRENCY`, `M_NATION`, `M_PROVINCE`, `M_CAUSE_OF_LOSS`, `D_CAUSE_OF_LOSS`, `M_SITE_DATABASE`, `MST_USER_TEKNIS`, `OS_AKSEPTASI_SUBJECTIVITY`, `TANGGAL_CLOSING`, `GCP_IMAGE`, `BRANCH`, `CITYINPUT`, `DISTRICTINPUT`, `RWINPUT`. **Workbook kedua menutup NOL di antaranya** — `CURRENCYSTANDARD` adalah view yang membaca `m_CurrencyStandard`, objek berbeda (D19)

**Bentuk jawaban yang diterima**: daftar kolom beserta tipe, panjang, presisi, skala, dan nullable — satu baris per kolom

---

## 14. REQ-002

**Menahan**: — tidak menahan pekerjaan tertentu; dibutuhkan untuk kelengkapan model

**Jenis**: TABLE · **Prioritas kami**: BLOCKER

**Yang diminta**: `ALL_TAB_COLUMNS` untuk 44 tabel bisnis

**Bentuk jawaban yang diterima**: daftar kolom beserta tipe, panjang, presisi, skala, dan nullable — satu baris per kolom

---

## 15. REQ-012

**Menahan**: — tidak menahan pekerjaan tertentu; dibutuhkan untuk kelengkapan model

**Jenis**: TABLE · **Prioritas kami**: BLOCKER

**Yang diminta**: Isi `Data-Admin-DB-Table` untuk class `ASM-FW-%`

**Bentuk jawaban yang diterima**: daftar kolom beserta tipe, panjang, presisi, skala, dan nullable — satu baris per kolom

---

## 16. REQ-023

**Menahan**: — tidak menahan pekerjaan tertentu; dibutuhkan untuk kelengkapan model

**Jenis**: METADATA · **Prioritas kami**: BLOCKER

**Yang diminta**: Siapa menulis dan siapa membaca ~60 kolom skalar `OS_AKSEPTASI_KLAIM` yang tidak diisi `PEGA_JSON_OS_AKSEP_KLAIMTNP`

**Bentuk jawaban yang diterima**: daftar nama beserta atributnya

---

## 17. REQ-024

**Menahan**: — tidak menahan pekerjaan tertentu; dibutuhkan untuk kelengkapan model

**Jenis**: DATA-PROFILE · **Prioritas kami**: BLOCKER

**Yang diminta**: `PYREOPENCOUNT > 0` dan `PYREOPENTIMESTAMP` — apakah klaim pernah di-reopen (menjawab A4, mengonfirmasi atau membatalkan ADR-0002)

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 18. REQ-029

**Menahan**: — tidak menahan pekerjaan tertentu; dibutuhkan untuk kelengkapan model

**Jenis**: DATA-PROFILE · **Prioritas kami**: BLOCKER

**Yang diminta**: Akseptasi ber-`STS_SUBJECTIVITY='1'` yang **sudah punya catatan pembayaran** — apakah uang keluar atas syarat yang belum terpenuhi

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 19. REQ-013

**Menahan**: FINDING-002

**Jenis**: DATA-PROFILE · **Prioritas kami**: PENTING

**Yang diminta**: Cacah pola teks pada `.CommentSuggest`

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 20. REQ-022

**Menahan**: FINDING-006

**Jenis**: DATA-PROFILE · **Prioritas kami**: PENTING

**Yang diminta**: Mata uang **bukan IDR** yang kursnya tersimpan `1` — mengukur jalur `NO_DATA_FOUND` pada `GETCURRENCYSTANDARD`

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 21. REQ-034

**Menahan**: FINDING-008

**Jenis**: DATA-PROFILE · **Prioritas kami**: VERIFIKASI — turun dari PENTING 19 Sep 2026 (FINDING-008 bagian 5)

**Yang diminta**: **Menguji ramalan FINDING-008 (usulan).** Cacah baris `SpreadingRisk` yang `Currency`-nya **bukan** `IDR` **dan** `ClaimAmountIDR`-nya **sama persis** dengan nilai mata uang aslinya — yaitu baris yang kolom rupiahnya tidak pernah dikonversi. Sertakan sebaran `KursIDR` pada baris-baris itu. Nilainya ada di dalam `PZPVSTREAM`, jadi ia **menumpang REQ-001**; bila REQ-001 gagal, cadangannya membongkar blob. Read-only, tanpa data nasabah, maksimum 20 baris contoh dengan kolom teks bebas ditutup

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 22. REQ-005

**Menahan**: — tidak menahan pekerjaan tertentu; dibutuhkan untuk kelengkapan model

**Jenis**: DATA-PROFILE · **Prioritas kami**: PENTING

**Yang diminta**: Profil presisi kolom uang & kurs — **dipersempit: hanya tabel bisnis `POOLDATA`**. Untuk `.ClaimData.*` jalur ini **tertutup permanen** (IN-BLOB terbukti)

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 23. REQ-006

**Menahan**: — tidak menahan pekerjaan tertentu; dibutuhkan untuk kelengkapan model

**Jenis**: DATA-PROFILE · **Prioritas kami**: PENTING

**Yang diminta**: Statistik pemakaian jalur tutup-langsung (mengambil alih A5)

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 24. REQ-007

**Menahan**: — tidak menahan pekerjaan tertentu; dibutuhkan untuk kelengkapan model

**Jenis**: TABLE/VIEW/SYNONYM · **Prioritas kami**: PENTING

**Yang diminta**: Verifikasi keberadaan 8 objek DERIVED

**Bentuk jawaban yang diterima**: daftar kolom beserta tipe, panjang, presisi, skala, dan nullable — satu baris per kolom

---

## 25. REQ-008

**Menahan**: — tidak menahan pekerjaan tertentu; dibutuhkan untuk kelengkapan model

**Jenis**: TABLE · **Prioritas kami**: PENTING

**Yang diminta**: Tabel work Pega untuk `ASM-FW-GCNMFW-Work`

**Bentuk jawaban yang diterima**: daftar kolom beserta tipe, panjang, presisi, skala, dan nullable — satu baris per kolom

---

## 26. REQ-009

**Menahan**: — tidak menahan pekerjaan tertentu; dibutuhkan untuk kelengkapan model

**Jenis**: DATA-PROFILE · **Prioritas kami**: VERIFIKASI — turun dari PENTING 19 Sep 2026 (SPEC §21.8)

**Yang diminta**: Struktur JSON di `json_klaim.DATA_JSON` & `OS_AKSEPTASI_KLAIM.DATA_JSON`

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 27. REQ-010

**Menahan**: — tidak menahan pekerjaan tertentu; dibutuhkan untuk kelengkapan model

**Jenis**: INDEX/CONSTRAINT · **Prioritas kami**: PENTING — turun dari penahan index 19 Sep 2026 (SPEC §21.7)

**Yang diminta**: Constraint, index, dan volume tabel inti

**Bentuk jawaban yang diterima**: daftar nama beserta atributnya

---

## 28. REQ-014

**Menahan**: — tidak menahan pekerjaan tertentu; dibutuhkan untuk kelengkapan model

**Jenis**: DATA-PROFILE · **Prioritas kami**: PENTING

**Yang diminta**: Cacah klaim dengan lebih dari satu baris `TreatyType='UR'`

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 29. REQ-019

**Menahan**: — tidak menahan pekerjaan tertentu; dibutuhkan untuk kelengkapan model

**Jenis**: DATA-PROFILE · **Prioritas kami**: VERIFIKASI — turun 19 Sep 2026; ramalan I2 memverifikasi, tidak menahan

**Yang diminta**: Cacah `FLAGONGOINGCOMMITTE` menyala saat klaim ditutup, sebaran `KOMITECOUNT`, **dan klaim tertutup dipecah menurut ada-tidaknya Adjustment** — yang terakhir menguji hipotesis I2

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 30. REQ-025

**Menahan**: — tidak menahan pekerjaan tertentu; dibutuhkan untuk kelengkapan model

**Jenis**: DATA-PROFILE · **Prioritas kami**: PENTING

**Yang diminta**: `PXAPPLICATION`, `PXAPPLICATIONVERSION` dikelompokkan dengan `PXOBJCLASS` — aplikasi/ruleset lain yang menyentuh class ini (menjawab A9)

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 31. REQ-026

**Menahan**: — tidak menahan pekerjaan tertentu; dibutuhkan untuk kelengkapan model

**Jenis**: DATA-PROFILE · **Prioritas kami**: PENTING

**Yang diminta**: `PXCREATEOPERATOR` / `PXUPDATEOPERATOR` + `PXUPDATEDATETIME` untuk kelima nama di FINDING-002 — aktivitas terakhir tiap orang (menjawab A11a)

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 32. REQ-028

**Menahan**: — tidak menahan pekerjaan tertentu; dibutuhkan untuk kelengkapan model

**Jenis**: DATA-PROFILE · **Prioritas kami**: VERIFIKASI — turun dari PENTING 19 Sep 2026 (SPEC §21.3)

**Yang diminta**: Cacah baris berkode setoran `"100115"` di `DIRECTTOKASIR_LOG` atau tabel kasir — apakah klaim non-prop syariah pernah ada (menutup G1)

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 33. REQ-030

**Menahan**: — tidak menahan pekerjaan tertentu; dibutuhkan untuk kelengkapan model

**Jenis**: DATA-PROFILE · **Prioritas kami**: PENTING

**Yang diminta**: Adakah baris `json_polis` ber-`BusinessFac='F'` yang dirujuk klaim non-prop (menguji ADR-0020)

**Bentuk jawaban yang diterima**: satu angka, atau satu tabel cacah bila diminta dipecah per kelompok

---

## 34. REQ-035

**Menahan**: — tidak menahan pekerjaan tertentu; dibutuhkan untuk kelengkapan model

**Jenis**: TABLE · **Prioritas kami**: PENTING (menahan tiket `20`, dokumen klaim)

**Yang diminta**: **Isi 42 baris `Rule-Declare-DecisionTable GetMimeType`** pada class `ASM-FW-GISFW-Int-T_STORAGE_IMAGE`. Ekspor memuat 42 slot baris yang **nilainya menutup sendiri** — kerangkanya ada, pemetaannya tidak. Tabel ini memetakan ekstensi berkas ke tipe MIME. **Celah bahan, bukan celah metode**: tidak ada pola sapuan yang dapat memulihkannya (D45)

**Bentuk jawaban yang diterima**: daftar kolom beserta tipe, panjang, presisi, skala, dan nullable — satu baris per kolom

---


## Bila ada yang ingin ditanyakan balik

Setiap permintaan di sini lahir dari sesuatu yang **tidak terbaca** di ekspor rule Pega yang kami punya. Bila sebuah permintaan tampak aneh, kemungkinan besar kami salah menebak bentuk datanya — dan koreksi Anda lebih berharga daripada jawaban atas pertanyaan yang salah.
