# TICKETS — lapisan data modul Claim Non Prop

<!-- STEMPEL ASAL -->
> **Dasar bukti**: `SPEC-MODEL-DATA.md` (804 baris), `KAMUS-KOLOM.md` (22 tabel · 391 kolom · 9 view), `ddl-usulan/` (32 berkas), `REGISTER-RATIFIKASI.md` (8 butir AK), `_selesai/OPEN-QUESTIONS.md`, `ORACLE-REQUESTS.md`.
> Disusun 18 September 2026. **Hanya lapisan data.** Tidak ada tiket Golang, React, endpoint, layar, service, repository, ORM, pipeline, atau infrastruktur — `SPEC-MODEL-DATA.md` menyatakan seluruhnya di luar lingkup.
> Tidak ada estimasi waktu di dokumen ini. Yang ditulis urutan dan ketergantungan.

## Gerbang — lolos 18 September 2026

| # | Prasyarat | Hasil pemeriksaan |
|---|---|---|
| G1 | View DDL ada beserta entri kamusnya | **9 objek `CREATE OR REPLACE VIEW`** di `ddl-usulan/V01…V09`; **9 entri** di `KAMUS-KOLOM.md` |
| G2 | Sisi kedua GRANT ada | **18 `GRANT SELECT`** atas view untuk `POOLDATA` dan `KLAIMNP_HILIR`; **42 `REVOKE ALL`** atas tabel kanonik di `V00_HAK_AKSES.sql`; hak hilir atas tabel kanonik: **nol** |
| G3 | `MIGRASI_KORELASI` berbentuk komposit | **lima kolom terpisah**: `PZINSKEY_LAMA`, `PYID_LAMA`, `CASEID_LAMA`, `INDEX_OBJECT_LAMA`, `URUTAN_BARIS_LAMA` |

---

## 1. Ringkasan

**39 tiket**, delapan jalur. **Tujuh sudah selesai atau selesai sebagian** sebelum dokumen ini ditulis — ditandai di tempatnya, bukan dihapus, supaya nomor yang dirujuk tetap menunjuk.

| Jalur | Tiket | Isi |
|---|---|---|
| A — fondasi | T-01, T-02 | skema, akun, tabel singkatan |
| B — tabel kanonik | T-03 … T-11 | sembilan rumpun, berurut menurut FK |
| C — aturan lintas tabel | T-12 … T-15 | domain keadaan, pasangan IDR–kurs, penamaan, index penopang |
| D — lapisan hilir | T-16 … T-20 | view dan penutupan hak akses |
| E — migrasi | T-21 … T-23, T-39 | pendaratan, nilai ditolak, korelasi, arsip muatan keluar |
| F — uji | T-24 … T-31 | satu tiket per baris Testing Decisions |
| G — penyelidikan | T-32 … T-37 | sapuan yang hasilnya mengubah kolom |
| H — keputusan | T-38 | tipe desimal di sisi Golang |

**Dapat dimulai hari pertama, tanpa prasyarat**: T-01 · T-34 · T-35 · T-36 · T-37 · T-38.

Lima dari enam itu tiket penyelidikan dan keputusan — dan itu bukan kebetulan. **T-35 dan T-36 menahan T-05**, jadi menjalankannya lebih dulu mencegah `ALOKASI_LAYER` dibongkar sesudah terpasang.

### Yang berubah dari pemotongan yang diusulkan, beserta alasannya

| Tiket | Perubahan | Alasan |
|---|---|---|
| T-02 | **SELESAI** | Tabel singkatan tertutup sudah ditulis di `SPEC-MODEL-DATA.md` bagian 16, dan **ditegakkan mesin**: pembangkit menghentikan keluaran bila ada pengenal > 30 byte |
| T-12 … T-15 | **RANCANGAN SELESAI**, pemasangan ikut Jalur B | Keempatnya sudah ada di DDL. Kekhawatiran "dikerjakan sembilan kali dengan sembilan tafsir" **tidak berlaku**: constraint-nya dipasang oleh satu pembangkit atas seluruh tabel sekaligus, bukan ditulis tangan per tabel. Angka "nol CHECK" dan "1 dari 10" di pemotongan usulan berasal dari keadaan sebelum 18 Sep |
| T-32, T-33 | **SELESAI** | Sapuan S1 dan S2 sudah dijalankan atas empat lapisan; hasilnya di `SPEC-MODEL-DATA.md` bagian 18 |
| T-34 | **tidak lagi menahan** | §4.1 menetapkan tempatnya dibuat dan nilainya menunggu — pola yang sama dengan `TARIF_BERLAKU`. T-34 berubah dari penahan menjadi pengisi |
| T-39 | **baru** | Arsip muatan keluar, lahir dari penutupan temuan S1 |

---

## 2. Urutan pengerjaan

```
T-01 ── T-03 ─┬─ T-04   T-05   T-06   T-07   T-08   T-09   T-10   T-11   T-39
              │  (seluruhnya sejajar — FK-nya hanya ke KLAIM)
              └─ T-21   T-22   T-23

T-35, T-36 ──> menahan T-05 saja
T-37, T-34 ──> tidak menahan apa pun
T-16 ← T-06, T-23, T-39      T-19 ← T-23
T-20 ← seluruh tabel dan view terpasang
T-38 ──────> menahan seluruh tiket aplikasi di spec berikutnya
```

**Koreksi 18 September, jawaban atas §7.1.** Rantai `T-04 → T-05 → T-06` pada versi pertama **terserialisasi tanpa sebab**. Diperiksa ulang terhadap DDL: `ALOKASI_LAYER`, `RETENSI_CEDANT`, `AKSEPTASI`, dan `ADJUSTMENT` seluruhnya ber-FK **hanya ke `KLAIM`** — tidak satu pun menunjuk `NILAI_KLAIM_MATA_UANG`. Urutan itu urutan penalaran saya, bukan urutan yang ditegakkan basis data.

Akibatnya **sepuluh tiket Jalur B dan E dapat berjalan sejajar** begitu T-03 selesai, bukan berbaris. Satu-satunya penahan nyata di dalam Jalur B adalah T-35 dan T-36 terhadap T-05, dan itu penahan **bentuk kolom**, bukan penahan FK.

Berurut, dengan prasyarat terlihat:

| # | Tiket | Prasyarat |
|---|---|---|
| 1 | T-01 skema dan akun | — |
| 2 | T-35, T-36 sapuan penahan bentuk kolom | — |
| 3 | T-34 sapuan domain jenis pembayaran | — |
| 4 | T-03 Klaim | T-01 |
| 5 | T-04, T-06 … T-11, T-21 … T-23, T-39 — **sejajar** | T-03 |
| 6 | T-05 alokasi dan Retensi Cedant | T-03, T-35, T-36 |
| 7 | T-07 premi pemulihan | T-05 |
| 8 | T-16 … T-19 view | T-06, T-23, T-39 |
| 9 | T-20 penutupan hak akses | seluruh tabel dan view terpasang |
| 10 | T-24 … T-31 uji | tiket yang diujinya |

---

## 3. Tiket

### T-01 — Skema `KLAIMNP` dan akun aplikasi berdiri, dan tidak ada akun lain yang dapat menulis

**Hasil** Skema, akun pemilik, akun aplikasi `KLAIMNP_APP`, dan akun hilir `KLAIMNP_HILIR` ada; akun pemilik hanya dipakai saat pemasangan.
**Dasar** DECIDED(ADR-0028) bagian *Satu pintu tulis*; DECIDED(ADR-0017). Premisnya EVIDENCED: `pengetahuan/ddl/TABLE_PC_ASM_FW_GCNMFW_WORK.sql` memberi `POOLDATA` hak `ALTER, DELETE, INSERT, UPDATE` atas tabel work Pega.
**Prasyarat** —
**Lingkup** Pembuatan skema dan ketiga akun; tablespace; peran dasar.
**Tidak termasuk** `GRANT` dan `REVOKE` atas objek — belum ada objeknya. → T-20.
**Kriteria terima**
- Akun aplikasi dapat membuat sesi dan melihat skema.
- Akun hilir dapat membuat sesi dan **tidak** melihat satu objek pun (belum ada).
- Akun pemilik terpisah dari akun aplikasi — bukan akun yang sama dengan nama berbeda.

**Ketidakpastian** Tidak ada. **REQ-032 turun jadi verifikasi 19 September 2026**: batas 30 byte ditetapkan sebagai **aturan tetap**, sah di setiap versi Oracle. Bila versinya 12.2+, yang ada hanyalah kelonggaran yang sengaja tidak dipakai. Tabel singkatan tertutup **dibatalkan** — tiket `02` mati.

---

### T-02 — Tabel singkatan tertutup ditulis dan dibekukan — **SELESAI**

**Hasil** Tidak ada singkatan ad-hoc sesudah ini.
**Dasar** `SPEC-MODEL-DATA.md` bagian 16. DECIDED-TEKNIS.
**Prasyarat** —
**Lingkup** Sudah ada: 23 istilah, awalan `PK_` `UQ_` `FK_` `CK_` `IX_`, aturan *nama tabel memakai istilah selengkapnya, nama constraint selalu memakai singkatan*.
**Kriteria terima**
- Pengenal terpanjang terhitung: tabel `KLAIM_PENJAGA_TANGGAL` **28 byte**, constraint `UQ_TARIF_BERLAKU_1` **20 byte**.
- Pengenal di atas 30 byte: **nol**, dan pembangkit **menghentikan keluaran** bila ada — bukan melewatkannya.

**Ketidakpastian** Tidak ada. `UQ_AKSEPTASI_KLAIM_LAYER_MATA_UANG` (34 byte) yang menjadi alasan daftar ini menjadi `UQ_AKSEPTASI_1` (18 byte).

---

### T-03 — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap

**Hasil** Klaim dapat disisipkan; nomor klaim ganda ditolak; masa berlaku treaty terbalik ditolak.
**Dasar** DECIDED(ADR-0001, ADR-0021, ADR-0022). EVIDENCED: `BLUEPRINT.md` §2.2 — sapuan 114 activity, **tidak satu pun `Property-Set` menulis `.DateOfLoss`**.
**Prasyarat** T-01
**Lingkup** `KLAIM` beserta `UQ_KLAIM_1`, `CK_KLAIM_1`, `CK_KLAIM_2`, `CK_KLAIM_3`, `IX_KLAIM_1`.
**Tidak termasuk** `CHECK` domain untuk `STATUS_KLAIM` — sengaja tidak ada, lihat T-12. Kolom jejak pelaku dipasang di sini tetapi tanpa foreign key; tabel pengguna tidak dirancang (ADR-0006).
**Kriteria terima**
- Klaim dengan nomor yang sudah ada **ditolak**.
- Klaim ber-`TREATY_MULAI` sesudah `TREATY_AKHIR` **ditolak**.
- Klaim dengan `KEADAAN_BARIS` di luar tiga nilai domain **ditolak**.
- Klaim tanpa `DIBUAT_OLEH` **ditolak**; `DIBUAT_ATAS_NAMA` kosong **diterima** — kosong berarti tidak ada perwakilan, bukan tidak diketahui.

**Ketidakpastian** **REQ-018** (duplikat `PYID`) belum terjawab. `UQ_KLAIM_1` dipasang sekarang. Bila data lama melanggarnya, yang berubah adalah **penanganan migrasi**, bukan kuncinya — ADR-0024 sudah menetapkan pemilahan tiga kelompoknya.

---

### T-04 — Nilai kerugian per mata uang, dengan nilai IDR yang tidak dapat lahir tanpa kurs

**Hasil** Nilai per mata uang tersimpan; nilai IDR tanpa kurs ditolak basis data.
**Dasar** DECIDED(ADR-0007, ADR-0014, ADR-0019). EVIDENCED: FINDING-006 — `GETCURRENCYSTANDARD` mengembalikan `1` bila kurs tidak ditemukan, dan `1` tidak dapat dibedakan dari kurs yang sah.
**Prasyarat** T-03
**Lingkup** `NILAI_KLAIM_MATA_UANG`, empat besaran uang berpasangan IDR, `UQ_NILAI_KLAIM_MATA_UANG_1`, empat `CK` pasangan IDR–kurs, `CK` domain keadaan.
**Tidak termasuk** Penularan keadaan ke baris turunan — diuji di T-28 setelah tabel turunannya ada.
**Kriteria terima**
- Baris dengan `NILAI_KERUGIAN_IDR` terisi sementara `KURS` kosong **ditolak** — berlaku untuk keempat besaran.
- Baris dengan `KURS` terisi dan `*_IDR` kosong **diterima** — belum dikonversi, bukan salah.
- Baris kedua dengan (klaim, mata uang) sama **ditolak**.
- Nilai nol pada `*_IDR` **diterima** dan dapat dibedakan dari `NULL` lewat kueri.

**Ketidakpastian** Tidak ada REQ yang menyentuhnya.

---

### T-05 — Alokasi per Layer terpisah dari Retensi Cedant, dan suntingan tidak dapat tertimpa hitung ulang

**Hasil** Hitung ulang menimpa `_HITUNG` dan tidak menyentuh `_SUNTING`; `_DIPAKAI` mengikuti tanpa ditulis siapa pun.
**Dasar** DECIDED(ADR-0010, ADR-0008, ADR-0011). EVIDENCED: `BLUEPRINT.md` §2.5 — `.IsEditClaim` **ditulis lima kali, tidak pernah dibaca**; FINDING-004.
**Prasyarat** T-03, T-35, T-36 — *bukan* T-04; FK-nya hanya ke `KLAIM`
**Lingkup** `ALOKASI_LAYER` dan `RETENSI_CEDANT`; kolom berpasangan `_HITUNG`/`_SUNTING`/`_DIPAKAI`; parameter yang di-snapshot; `UQ_ALOKASI_LAYER_1`, `UQ_RETENSI_CEDANT_1`.
**Tidak termasuk** Kolom jenis reasuransi pada `RETENSI_CEDANT` — **sengaja tidak ada**. EVIDENCED: baris Retensi Cedant menulis `TreatyName = "UR"` secara literal melewati lookup (FINDING-003 bagian 2); `"UR"` artefak penyatuan yang dibatalkan ADR-0010, bukan nilai domain.
**Kriteria terima**
- Menjalankan ulang pengisian `_HITUNG` atas baris yang `_SUNTING`-nya terisi **tidak mengubah** nilai yang dibaca lewat `_DIPAKAI`.
- Menulis langsung ke `_DIPAKAI` **ditolak** — ia kolom turunan.
- Baris kedua dengan (klaim, keempat field layer, mata uang) sama **ditolak**.
- Baris Retensi Cedant kedua untuk (klaim, mata uang) sama **ditolak**.
- `DISUNTING_OLEH` kosong sementara `_SUNTING` terisi **ditolak** — suntingan selalu punya pelaku.

**Ketidakpastian** **T-36** menentukan besaran mana yang benar-benar berpasangan `_HITUNG`/`_SUNTING`; memasangnya di kolom yang tidak pernah disunting adalah beban mati. **T-35** menentukan apakah Biaya Penilaian, Salvage, dan Biaya Lain perlu pasangan IDR di tabel ini. **REQ-033** (sebaran `LayerPart`) menentukan apakah kunci halus pernah melahirkan dua baris di tempat sistem lama melihat satu — itu mengubah **view** di T-16, bukan tabel ini.

---

### T-06 — Akseptasi dengan kunci alami yang ditegakkan, dan Adjustment sebagai unit pembayaran

**Hasil** Akseptasi kedua untuk klaim, layer, dan mata uang yang sama ditolak. Satu klaim dapat memiliki banyak Adjustment bernomor urut tetap.
**Dasar** DECIDED(ADR-0024, ADR-0015, ADR-0023). EVIDENCED: blok yang **dikomentari** di `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP` memakai persis kombinasi itu; `OS_AKSEPTASI_KLAIM` **tanpa primary key maupun unique**, dan prosedurnya **selalu `INSERT`**.
**Prasyarat** T-03 — *bukan* T-05; T-34 **tidak lagi menahan** (§4.1)
**Lingkup** `AKSEPTASI`, `ADJUSTMENT`, `REKENING_PENERIMA`; `UQ_AKSEPTASI_1`; `UQ_ADJUSTMENT_1`; `FK_ADJUSTMENT_2` beserta `IX_ADJUSTMENT_1`; `CK_AKSEPTASI_1`.
**Tidak termasuk** `CHECK` domain untuk `KEPUTUSAN_KOMITE` — **sengaja tidak ada**, domainnya EXTERNAL. Tabel Komite dan keanggotaannya tidak dirancang.
**Kriteria terima**
- `INSERT` kedua dengan (klaim, layer, mata uang) sama **ditolak**.
- Adjustment dengan nomor urut yang sudah dipakai pada klaim yang sama **ditolak**.
- Menghapus rekening yang masih dirujuk Adjustment **ditolak**.
- `KEADAAN_AKSEPTASI` di luar tiga nilai domain **ditolak**.
- `KEPUTUSAN_KOMITE` bernilai apa pun **diterima** — dan itu keputusan, bukan kelalaian.

**Ketidakpastian** **I1/I2** — transisi status akseptasi `0 → 1/2` **TIDAK DITEMUKAN** di 279 berkas; ia ditulis satu kali saja sebagai `0`. **H1/H2** — apakah `.ValueAdjustment` sudah IDR di hulu belum diketahui, sehingga kolom nilai Adjustment dipasang berpasangan mata uang dan **tidak mengandaikan** salah satunya. **REQ-018** mengukur baris ganda pada kunci lama, yang lebih longgar dari kunci ini.

---

### T-07 — Premi pemulihan dapat dihitung ulang dari barisnya sendiri

**Hasil** Setiap masukan rumus tersimpan bersama hasilnya; limit layer nol ditolak.
**Dasar** EVIDENCED: `MEMORI_PEMAHAMAN.MD` §6.3 — `CNPLimit` adalah **penyebut**. DECIDED-TEKNIS(AK-2b): acuannya `CountReinstatement_Act`; `AdjClaimCNP_Act` baris 3109 **tidak diwarisi** (FINDING-007).
**Prasyarat** T-05
**Lingkup** `PREMI_PEMULIHAN` beserta seluruh masukan rumus; `CK_PREMI_PEMULIHAN_1`; `UQ_PREMI_PEMULIHAN_1`.
**Tidak termasuk** Perhitungannya sendiri — ini lapisan data.
**Kriteria terima**
- Baris ber-`LIMIT_LAYER` nol **ditolak**.
- Baris kedua untuk (klaim, layer, mata uang) sama **ditolak**.
- Seluruh masukan rumus dapat dibaca dari satu baris, tanpa menyentuh tabel master treaty.

**Ketidakpastian** AK-2b menempatkan selisih terhadap sistem lama sebagai **pengecualian bernama** pada shadow-run, bukan kegagalan cutover. Besaran selisihnya belum terukur.

---

### T-08 — Masukan mesin alokasi dan penyebaran tersimpan, bukan hanya keluarannya

**Hasil** Pembagian kerugian, penyebaran, dan estimasi awal tersimpan sebagai data.
**Dasar** EVIDENCED: `MEMORI_PEMAHAMAN.MD` §6.2 Langkah 2 — `CNPSpreadLoss` adalah **masukan** mesin alokasi. DECIDED(ADR-0011, ADR-0013): hitung ulang idempoten mensyaratkan seluruh masukan tersimpan. `CONTEXT.md` **Penyebaran**.
**Prasyarat** T-03
**Lingkup** `PEMBAGIAN_KERUGIAN`, `PENYEBARAN`, `ESTIMASI_AWAL`.
**Tidak termasuk** `SpreadingAdjustment` dan `SpreadingAdjustmentQS` — keduanya **view**, bukan tabel. → T-18.
**Kriteria terima**
- Nomor urut ganda dalam satu klaim **ditolak**, di ketiga tabel.
- Baris penyebaran dapat dijumlahkan per jenis reasuransi per mata uang tanpa membaca tabel lain.

**Ketidakpastian** Tidak ada REQ yang menyentuhnya.

---

### T-09 — Objek, kronologi, dan dokumen klaim

**Hasil** Objek pertanggungan, kronologi, dan rujukan dokumen tersimpan; URL tanpa masa berlaku ditolak.
**Dasar** DECIDED(ADR-0027): berkas di Google Cloud Storage, yang tersimpan hanya metadata dan URL beserta `EXPDATE`. DECIDED(ADR-0009): kronologi terstruktur, keterangan tidak pernah dibaca mesin.
**Prasyarat** T-03
**Lingkup** `OBJEK_PERTANGGUNGAN`, `KRONOLOGI_KLAIM`, `DOKUMEN_KLAIM` beserta `CK_DOKUMEN_KLAIM_1`.
**Tidak termasuk** `ObjectList` — **lubang tercatat**, bukan tabel. Apa yang membedakannya dari `InterestList` **TIDAK DITEMUKAN DI XML**.
**Kriteria terima**
- Dokumen dengan URL terisi tetapi tanpa kedaluwarsa **ditolak**, dan sebaliknya.
- Dokumen tanpa URL sama sekali **diterima** — URL belum pernah diterbitkan.
- Kronologi tanpa jenis tindakan **ditolak**; tanpa keterangan **diterima**.

**Ketidakpastian** Tidak ada REQ. Lubang `ObjectList` dicatat di `SPEC-MODEL-DATA.md` bagian 10.

---

### T-10 — Tanggal tutup buku dan tarif menjadi data bertanggal berlaku

**Hasil** Satu sumber kebenaran untuk tanggal tutup buku dan tarif; keduanya terisi baris awalnya.
**Dasar** DECIDED(ADR-0025). DECIDED-TEKNIS(AK-5, AK-6.3). EVIDENCED: angka `25` tertanam di `HitServiceToKasir_Act`; `PROC_GENERATE_SEQUENCE_NUMBER` membacanya dari `POOLDATA.TANGGAL_CLOSING` — **dua sumber kebenaran untuk satu aturan**.
**Prasyarat** T-03
**Lingkup** `TUTUP_BUKU`, `TARIF_BERLAKU`; baris awal: tutup buku `25` berlingkup global; tarif brokerage 2,5%, PPh 2%, PPN 2,2%, seluruhnya berlingkup global dan berlaku sejak sebelum data tertua.
**Tidak termasuk** Riwayat tarif — **tidak ada bukti tarif pernah berubah** (AK-5), jadi tidak ada riwayat yang dibuat-buat. Faktor `102,2` **tidak disimpan**; ia diturunkan sebagai `(100 + tarif PPN)`.
**Kriteria terima**
- Dua baris dengan (lingkup, tanggal berlaku) sama **ditolak**.
- Tanggal tutup buku di luar 1–31 **ditolak**.
- Tarif terbaca lewat satu kueri tanpa membaca kode mana pun.

**Ketidakpastian** **A8** dan **ASK-AKUNTANSI no. 5** — apakah ketiga tarif masih berlaku **belum dikonfirmasi**. Yang dibuat tempatnya; nilainya berlabel DECIDED-TEKNIS menunggu ratifikasi, dan **ratifikasi tidak menahan tiket ini**.

---

### T-11 — Koreksi bernilai tercatat menggantikan tambalan di dalam kode

**Hasil** Setiap penyimpangan nilai punya nilai sebelum, nilai sesudah, alasan, pelaku, dan waktu.
**Dasar** DECIDED(ADR-0004, ADR-0018). EVIDENCED: `BLUEPRINT.md` §7 — **29 langkah, 18 ekspresi unik, 8 rule** menambal per-case; tambalan terbaru `CLMNP-975` bertanggal **2026-07-16**.
**Prasyarat** T-03
**Lingkup** `KOREKSI_NILAI`, `KLAIM_PENJAGA_TANGGAL`, `IX_KOREKSI_NILAI_1`.
**Tidak termasuk** Pengisian daftar klaim terdampak — menunggu migrasi.
**Kriteria terima**
- Koreksi tanpa alasan **ditolak**; tanpa pelaku **ditolak**.
- Koreksi dengan `NILAI_SESUDAH` kosong **diterima** — koreksi yang membatalkan nilai.
- Satu klaim hanya dapat muncul sekali di daftar klaim penjaga tanggal.
- Riwayat koreksi atas satu kolom tertentu terbaca lewat satu kueri.

**Ketidakpastian** Tidak ada REQ. Bila `KOREKSI_NILAI` dan kolom suntingan di tingkat baris berbeda, **`KOREKSI_NILAI` yang berlaku** — `SPEC-MODEL-DATA.md` bagian 20.

---

### T-12 — Domain kolom keadaan — **RANCANGAN SELESAI**, pemasangan ikut Jalur B

**Hasil** Kolom keadaan yang domainnya kita tetapkan sendiri terikat `CHECK`; yang domainnya warisan dan belum lengkap **sengaja tidak terikat**.
**Dasar** DECIDED(ADR-0019 lapis 2). EVIDENCED: sapuan S2 empat lapisan — `CNPStatusCase` hanya **2 nilai ditulis**, **nol rule menguji**, dan `"CLAIM ACCEPTED"`/`"CLAIM REJECTED"` **nihil di 279 berkas**.
**Prasyarat** T-03 … T-11 (pemasangannya ikut di sana)
**Lingkup** `CK` `KEADAAN_BARIS` di **11 dari 11** tabel yang punya kolom itu; `CK` `KEADAAN_AKSEPTASI`.
**Tidak termasuk** `STATUS_KLAIM` dan `KEPUTUSAN_KOMITE` — **tanpa `CHECK`, dan itu keputusan**. Domain tertutup yang isinya belum lengkap akan menolak nilai yang sah pada hari pertama modul Komite tersambung.
**Kriteria terima**
- Nilai `KEADAAN_BARIS` di luar `LENGKAP` · `MENUNGGU_KURS` · `GAGAL_URAI` **ditolak**, di kesebelas tabel.
- `STATUS_KLAIM` bernilai apa pun **diterima**, dan alasannya tercatat.

**Ketidakpastian** Kelengkapan enum `CNPStatusCase` hanya dapat dibuktikan dengan membuka folder `Komite Claim Non Prop` — **belum diizinkan**.

---

### T-13 — Pasangan IDR–kurs — **RANCANGAN SELESAI**, pemasangan ikut Jalur B

**Hasil** Tidak ada nilai IDR yang dapat lahir tanpa kurs yang menghasilkannya, di tabel mana pun.
**Dasar** DECIDED(ADR-0007, ADR-0014). EVIDENCED: FINDING-006 — `RETURN 1` pada kurs yang tidak ditemukan.
**Prasyarat** T-03 … T-11
**Lingkup** **23 `CHECK`** pasangan, dipasang **oleh pembangkit** atas setiap kolom `*_IDR` di tabel yang punya `KURS`.
**Tidak termasuk** Penularan keadaan ke baris turunan — itu aturan prosedur, diuji di T-28.
**Kriteria terima**
- Di **setiap** tabel bernilai uang, `*_IDR` terisi tanpa `KURS` **ditolak**.
- Menambahkan tabel bernilai uang baru tanpa `CHECK`-nya **tidak mungkin** — pembangkit memasangnya sendiri.

**Ketidakpastian** **T-35** dapat menambah kolom `*_IDR` di `ALOKASI_LAYER`; bila itu terjadi, `CHECK`-nya ikut terpasang tanpa perubahan tangan.

---

### T-14 — Penamaan yang tidak nyaris kembar — **SELESAI**

**Hasil** Tidak ada pasangan nama berjarak satu kata di skema yang belum berjalan.
**Dasar** EVIDENCED: `MEMORI_PEMAHAMAN.MD` §10.4 butir 4 mendaftar nama nyaris kembar sebagai technical debt sistem lama — `CountTotalInterest_Act` vs `CountTotalInsterest_Act`, `ProtectNilaiClaim` vs `ProteksiNilaiClaim`.
**Prasyarat** —
**Lingkup** `ADJUSTMENT.STATUS_AKSEPTASI` → **`KEPUTUSAN_KOMITE`**, dinamai menurut isinya: ia hasil keputusan Komite atas Adjustment, bukan keadaan entitas akseptasi. `FK_DPT_KLM_FK` → `FK_KLAIM_PENJAGA_TANGGAL_1`; sufiks `_FK` tersisa: **nol**.
**Kriteria terima**
- Tidak ada dua pengenal di skema yang berbeda hanya pada satu kata.
- Tidak ada constraint bersufiks `_FK`.

**Ketidakpastian** Tidak ada.

---

### T-15 — Index penopang setiap foreign key — **SELESAI**

**Hasil** Tidak ada FK yang kolomnya bukan awalan sebuah index atau unique.
**Dasar** DERIVED: di Oracle, FK tanpa index membuat `DELETE` dan `UPDATE` pada tabel induk mengambil kunci di tingkat tabel.
**Prasyarat** —
**Lingkup** `IX_ADJUSTMENT_1` untuk `FK_ADJUSTMENT_2`; sapuan ulang seluruh FK oleh pembangkit.
**Kriteria terima**
- Setiap FK punya index atau unique yang kolom pertamanya sama dengan kolom pertama FK itu.
- Empat belas FK ke `KLAIM` tetap tertopang `UNIQUE` berawalan `ID_KLAIM` — **tidak boleh hilang** saat kunci disentuh lagi.

**Ketidakpastian** Tidak ada.

---

### T-16 — Arasapas dan kasir menerima bentuk lama, dan yang tidak dapat diringkas tidak hilang diam-diam

**Hasil** View kompatibilitas menghasilkan bentuk yang dikenal hilir; kelompok yang besaran per-layernya tidak sepakat **ditolak** dan muncul di view penolakan.
**Dasar** DECIDED(ADR-0023, ADR-0021 pengecualian `CASEID`). EVIDENCED: sapuan S1 — 17 parameter `InputParamOs.*` di `Activity\SaveDataToOSAksep_Act.xml`; `TypeLoss = .TreatyName` adalah **satuan kasar**.
**Prasyarat** T-06, T-23, **T-39**
**Lingkup** `V_AKSEPTASI_KOMPATIBEL`, `V_AKSEPTASI_DITOLAK`. Keduanya membaca tabel kanonik **dan arsip muatan keluar**, lalu menghasilkan **selisih**.
**Tidak termasuk** Menegakkan kunci kasar sebagai `UNIQUE` kedua — **ditolak dengan alasan**: itu menjadikan keterbatasan bentuk lama sebagai aturan bisnis sistem baru.
**Kriteria terima**
- Kolom dan tipenya sama dengan yang diterima hilir hari ini, termasuk `CASEID`.
- Besaran aditif **dijumlahkan**; Limit Layer, Premi Deposit, persen premi pemulihan, kurs, dan porsi **tidak**.
- Bila besaran per-layer berbeda antar baris yang dilebur, view **tidak menghasilkan baris** untuk kelompok itu, dan kelompoknya muncul di view penolakan beserta sebabnya.
- `LayerPart` dan `LayerPartType` tidak muncul di keluaran.

**Ketidakpastian** **REQ-033** menentukan apakah agregasi halus→kasar pernah menemui baris yang tidak sepakat.

**Temuan S1 sudah TERTUTUP, dan tidak lewat REQ-018.** Sumber `OutOSAcc` ditemukan: `RDBList\GetDataOS.xml` dan `RDBList\GetDataCNPOS.xml` — EVIDENCED:

```sql
select sum(nvl(a.data_json.Value,0)) as "Value", ...
  from OS_AKSEPTASI_KLAIM a
 WHERE CASEID = {pyWorkPage.pzInsKey}
   AND a.data_json.TypeLoss = {InputParamOs.TypeLoss}
   AND a.data_json.Currency = {InputParamOs.Currency}
   AND STS_REJECT = 0
```

Penyaringnya **persis kunci alami ADR-0024**, dan agregatnya `SUM` atas seluruh baris berkunci sama. Jadi baris memang **tambahan**, bukan keadaan — dan itu pasangan wajib dari prosedur yang **selalu `INSERT`**. Naik dari DERIVED ke **EVIDENCED**.

**Satu syarat yang tidak terduga dan ikut mengikat**: penyaringnya memuat **`STS_REJECT = 0`**. Muatan yang ditolak **tidak** ikut dijumlahkan. Karena itu `ARSIP_MUATAN_KELUAR` berkolom `DITOLAK`, dan view menjumlahkan hanya yang tidak ditolak.

Yang tersisa di **REQ-018** hanya rekonsiliasi: apakah jumlah seluruh baris lama benar-benar sama dengan nilai sekarang di produksi. Bila tidak, itu temuan tentang **data lama**, bukan alasan mengubah rancangan.

---

### T-17 — Rekap per klaim per mata uang, dengan Retensi Cedant yang tidak pernah tercampur

**Hasil** Penjumlahan alokasi tidak pernah diam-diam memuat Retensi Cedant.
**Dasar** DECIDED(ADR-0010). EVIDENCED: `BLUEPRINT.md` §2.4 — di sistem lama setiap loop atas `SpreadingRisk` ikut melihat baris Retensi Cedant kecuali menyaringnya sendiri.
**Prasyarat** T-05
**Lingkup** `V_REKAP_KLAIM_MATA_UANG`.
**Kriteria terima**
- Retensi Cedant muncul sebagai **kolom tersendiri**, bukan bagian dari jumlah alokasi layer.
- Menjumlahkan kolom alokasi saja menghasilkan angka tanpa Retensi Cedant — tanpa penyaring apa pun di sisi pemanggil.

**Ketidakpastian** Baris Retensi Cedant di sistem lama membawa `TotalClaim` yang **selalu nol** pada cabang mata uang sama (ADR-0010 tambahan). Sistem baru menghitungnya seperti Retensi Cedant biasa (AK-2b), jadi selisih terhadap sistem lama **diharapkan** dan masuk pengecualian bernama shadow-run.

---

### T-18 — Lima PageList lama tersedia sebagai view, tanpa menyimpan hasil yang dapat menyimpang

**Hasil** `ListTotalEstimation`, `SpreadingAdjustment`, `SpreadingAdjustmentQS`, `TotalInterestInsured`, `ListClaimAcceptation` terbaca tanpa tabel penyimpan.
**Dasar** DECIDED(ADR-0013). EVIDENCED: `MEMORI_PEMAHAMAN.MD` §6.4 — kunci agregasi `SpreadingAdjustment` adalah `(Currency, TreatyName)`, yaitu **satuan kasar**; dan `CountSpreadingXOL` **menghapus lalu membangun ulang** isinya setiap kali dijalankan.
**Prasyarat** T-05, T-06, T-08, T-09
**Lingkup** `V_TOTAL_ESTIMASI`, `V_PENYEBARAN_AGREGAT`, `V_PENYEBARAN_AGREGAT_QS`, `V_TOTAL_NILAI_PERTANGGUNGAN`, `V_REKAP_AKSEPTASI`.
**Kriteria terima**
- Setiap view menghasilkan angka yang sama dengan penjumlahan langsung atas tabel sumbernya.
- Tidak ada tabel yang menyimpan angka-angka ini.

**Ketidakpastian** Penyaring quota share memakai `JENIS_REASURANSI_ID IN ('10028','10004')` — EVIDENCED dari `CountLossAllocation_act` Langkah 10. Apakah hanya dua nilai itu yang berarti quota share **TIDAK DITEMUKAN**; daftar penuhnya ada di `POOLDATA.REINSURANCETYPE`, yang **tidak punya satu pun constraint**.

---

### T-19 — Setiap baris baru dapat ditunjuk balik ke barisnya di sistem lama

**Hasil** Perbandingan baris per baris mungkin dilakukan.
**Dasar** DECIDED(ADR-0005). EVIDENCED: `BLUEPRINT.md` §8.4 — `IndexObject` adalah **posisi numerik, bukan surrogate key**, sehingga jembatannya hilang begitu baris berpindah ke kunci sendiri.
**Prasyarat** T-23
**Lingkup** `V_PARITAS_SHADOW`.
**Tidak termasuk** Pembandingnya sendiri dan toleransinya — AK-3.
**Kriteria terima**
- Setiap baris kanonik yang berasal dari migrasi muncul tepat sekali.
- Baris ber-`KEADAAN_BARIS` bukan `LENGKAP` dapat **dikeluarkan** dari perbandingan lewat satu penyaring, bukan dihitung sebagai selisih.

**Ketidakpastian** Baseline pembandingnya sendiri belum tentu benar: `TotalUR` selalu nol dan dua rumus premi pemulihan bekerja atas nilai berbeda. AK-2b menempatkan ketiganya sebagai pengecualian bernama.

---

### T-20 — Satu pintu tulis ditegakkan hak akses, bukan kesepakatan

**Hasil** Tidak ada akun selain akun aplikasi yang dapat menulis ke tabel kanonik; hilir membaca hanya lewat view.
**Dasar** DECIDED(ADR-0017, ADR-0028). EVIDENCED: `GRANT ALTER, DELETE, INDEX, INSERT, REFERENCES, SELECT, UPDATE ... TO POOLDATA` pada tabel work Pega — hak yang ada akan dipakai cepat atau lambat.
**Prasyarat** seluruh tabel (T-03 … T-11, T-21 … T-23) dan seluruh view (T-16 … T-19) terpasang
**Lingkup** `V00_HAK_AKSES.sql`: **44 `REVOKE ALL`** atas tabel kanonik untuk `POOLDATA` dan `KLAIMNP_HILIR`; `GRANT SELECT` atas view. **Termasuk `ARSIP_MUATAN_KELUAR`**: akun aplikasi memegang `INSERT` dan `SELECT` saja — **tanpa `UPDATE` dan tanpa `DELETE`**, karena arsip itu tulis-sekali.
**Tidak termasuk** Perintah yang dijalankan tangan — `GRANT` dan `REVOKE` adalah **objek DDL**. Hak yang tidak tertulis di DDL tidak dapat diaudit.
**Kriteria terima**
- `POOLDATA` gagal `INSERT`, `UPDATE`, dan `DELETE` ke setiap tabel kanonik.
- `POOLDATA` berhasil `SELECT` lewat view kompatibilitas.
- Akun hilir gagal `SELECT` langsung dari tabel kanonik mana pun.

**Ketidakpastian** **REQ-021** — apakah `POOLDATA` benar-benar menulis ke tabel Pega hari ini belum diketahui. Bila ya, migrasinya memerlukan pokok tersendiri untuk manajemen; itu **tidak menahan tiket ini**, karena skema baru tidak mengulangi hak itu apa pun jawabannya.

---

### T-21 — Bentuk lama boleh masuk utuh, di satu tempat saja

**Hasil** Satu-satunya tempat JSON hidup di skema ini.
**Dasar** DECIDED(ADR-0028). EVIDENCED: `BLUEPRINT.md` §13.2, §13.5 — nilai uang lama tersimpan sebagai **teks di dalam JSON**, dan `CLAIMXOL` mengeluarkannya sebagai `varchar2`.
**Prasyarat** T-01
**Lingkup** `MIGRASI_PENDARATAN` beserta `CK_MIGRASI_PENDARATAN_1`.
**Kriteria terima**
- Muatan yang bukan JSON sah **ditolak**.
- Tidak ada kolom JSON di tabel kanonik mana pun — diperiksa dengan sapuan atas seluruh skema.

**Ketidakpastian** Tidak ada REQ.

---

### T-22 — Nilai yang tidak dapat diurai tercatat, tidak dibulatkan, tidak dibuang

**Hasil** Tidak ada angka yang hilang tanpa jejak.
**Dasar** DECIDED(ADR-0014). DECIDED-TEKNIS(AK-4): migrasi berhenti bila ada **satu** baris yang tidak dapat diurai **dan** tidak dapat diselesaikan. Tidak ada baris yang dibuang karena "cuma sedikit".
**Prasyarat** T-01
**Lingkup** `MIGRASI_NILAI_DITOLAK`.
**Kriteria terima**
- Nilai mentah tersimpan **apa adanya**, tanpa pembulatan dan tanpa konversi.
- Setiap baris membawa sebab penolakan dan pengenal barisnya di sistem lama.

**Ketidakpastian** Tidak ada REQ. Volume belum terukur — **REQ-031**.

---

### T-23 — Jembatan ke sistem lama berdiri sebagai tabel terpisah

**Hasil** Setiap baris kanonik dapat ditunjuk balik lewat pengenal lamanya, tanpa mengurai teks.
**Dasar** DECIDED(ADR-0005, ADR-0023). EVIDENCED: kelima pengenal lama — `PZINSKEY` dan `PYID` dari DDL tabel work, `CASEID` dari `OS_AKSEPTASI_KLAIM` dan §13.2, `IndexObject` dari §8.4.
**Prasyarat** T-01
**Lingkup** `MIGRASI_KORELASI` dengan **lima kolom terpisah**, `IX_MIGRASI_KORELASI_1`.
**Tidak termasuk** `UNIQUE` atas kombinasi pengenal lama — **ditahan, dengan alasan**: kolom mana yang ada sudah EVIDENCED; **kombinasi mana yang unik** menunggu REQ-018.
**Kriteria terima**
- Kelima pengenal terbaca sebagai kolom tersendiri dan dapat di-`JOIN` tanpa penguraian teks.
- Baris yang bukan berasal dari work object Pega **diterima** dengan `PZINSKEY_LAMA` kosong — kosong berarti bukan dari sana, bukan tidak diketahui.

**Ketidakpastian** **REQ-018**, BLOCKER, OPEN. Yang berubah bila jawabannya lain: **constraint**, bukan kolom.

---

### T-39 — Arsip muatan keluar: apa yang benar-benar dikirim, tersimpan sebagai catatan

**Hasil** Setiap muatan yang dikirim ke Arasapas dan kasir tercatat apa adanya, dan muatan berikutnya dapat dihitung sebagai selisih terhadapnya.

**Dasar** DECIDED-TEKNIS(AK-1) — *"bila muatan keluar perlu diarsipkan, arsipnya tabel tersendiri, bukan kolom di tabel bisnis."* Ini pemakaian pertamanya. EVIDENCED: `RDBList\GetDataOS.xml` menjumlahkan seluruh baris berkunci `(CASEID, TypeLoss, Currency)` dengan `STS_REJECT = 0`, dan `SaveDataToOSAksep_Act` mengurangkan hasilnya — **baris adalah tambahan, bukan keadaan**.

**Prasyarat** T-03

**Lingkup** `ARSIP_MUATAN_KELUAR`: tujuan, kunci (klaim, jenis reasuransi, mata uang), lima besaran uang **berskala 2**, kurs, penanda ditolak, pelaku, waktu kirim. `IX_ARSIP_MUATAN_KELUAR_1`, `CK_ARSIP_MUATAN_KELUAR_1`.

**Tidak termasuk** Proses pengirimannya — lapisan aplikasi. Muatan dalam bentuk JSON — **tidak ada kolom JSON di sini**; kolomnya bertipe, karena daftar kolomnya sudah terbaca dari S1 (ADR-0028).

**Kriteria terima**
- Baris arsip tidak dapat di-`UPDATE` maupun di-`DELETE` oleh akun aplikasi — tulis-sekali ditegakkan hak akses, bukan kesepakatan.
- Dua muatan atas kunci yang sama pada waktu berbeda **keduanya tersimpan**; tidak ada yang menimpa.
- Muatan bertanda ditolak **tidak** ikut terjumlah saat selisih berikutnya dihitung.
- Nilai tersimpan **berskala 2** — ini satu-satunya tempat pembulatan tepi menjadi baris tersimpan.

**Ketidakpastian** Tiga hal yang ditulis terang supaya tidak dibaca sebagai pelanggaran ADR: view kompatibilitas **membaca** arsip sehingga ia tetap view, bukan salinan (ADR-0023 utuh); yang menulis arsip adalah proses pengiriman lewat akun aplikasi (ADR-0017 utuh); dan arsip berisi **apa yang benar-benar dikirim**, bukan apa yang seharusnya — ia catatan, bukan turunan yang dapat dihitung ulang.

---

### T-24 — Uji: kunci alami akseptasi

**Hasil** Terbukti bahwa baris ganda tidak dapat lahir.
**Dasar** DECIDED(ADR-0024). Bentuk uji mengikuti `pengetahuan/PREFLIGHT.sql`: SQL berblok, tiap blok berdiri sendiri, dapat dijalankan dari TOAD. Itu satu-satunya prior art — `MEMORI_PEMAHAMAN.MD` §10.8 mencatat ketiadaan kerangka uji sebagai "aspek yang tidak ada".
**Prasyarat** T-06
**Kriteria terima** `INSERT` kedua dengan (klaim, layer, mata uang) sama **ditolak**; `INSERT` dengan `LayerPart` berbeda **diterima**, dan itu perilaku yang dimaksud — bukan celah.
**Ketidakpastian** **REQ-033**.

---

### T-25 — Uji: pasangan uang dan mata uang

**Hasil** Terbukti bahwa tidak ada nilai uang yang dapat hidup tanpa satuannya, di tabel mana pun.

**Prasyarat** T-04, T-13
**Kriteria terima** Nilai uang tanpa mata uang **ditolak**; nilai IDR tanpa kurs **ditolak**; kurs tanpa IDR **diterima**. Diuji di **setiap** tabel bernilai uang, bukan satu.
**Dasar** DECIDED(ADR-0007). **Ketidakpastian** T-35 dapat menambah kolom yang ikut diuji.

---

### T-26 — Uji: batas tanggal inklusif

**Hasil** Terbukti bahwa klaim yang jatuh tepat di hari terakhir masa berlaku treaty tidak tertolak.

**Prasyarat** T-03
**Kriteria terima** Klaim ber-Tanggal Kejadian **tepat di hari terakhir** masa berlaku treaty **diterima**; sehari sesudahnya **ditolak**.
**Dasar** DECIDED(ADR-0022). EVIDENCED: FINDING-005 — di sistem lama perbandingannya dilakukan antara dua format teks yang berbeda.
**Ketidakpastian** **REQ-020** mengukur berapa klaim lama yang terdampak; tidak menahan uji ini.

---

### T-27 — Uji: NULL bukan nol

**Hasil** Terbukti bahwa nilai yang belum dihitung tidak pernah terbaca sebagai nol.

**Prasyarat** T-04
**Kriteria terima** Baris menunggu kurs ber-IDR `NULL`, `KEADAAN_BARIS` bukan `LENGKAP`, dan agregasi atasnya **tidak** memperlakukannya sebagai nol. Baris bernilai nol yang sudah dihitung dapat dibedakan darinya lewat satu kueri.
**Dasar** DECIDED(ADR-0019). **Ketidakpastian** Tidak ada.

---

### T-28 — Uji: penularan keadaan menunggu kurs

**Hasil** Terbukti bahwa ketiadaan kurs merambat ke seluruh turunannya, bukan berhenti di baris asalnya.

**Prasyarat** T-05, T-06, T-07
**Kriteria terima** Baris turunan dari nilai yang menunggu kurs **ikut** bertanda menunggu kurs, di seluruh tabel turunannya — alokasi, Retensi Cedant, akseptasi, Adjustment, premi pemulihan.
**Dasar** DECIDED(ADR-0014). **Ketidakpastian** Penularannya aturan prosedur; basis data menyimpan keadaannya, tidak menegakkan penularannya. Uji ini karena itu menguji **hasil**, bukan mekanismenya.

---

### T-29 — Uji: satu pintu tulis

**Hasil** Terbukti bahwa satu pintu tulis berlaku sebagai hak akses, bukan sebagai kesepakatan.

**Prasyarat** T-20
**Kriteria terima** Akun selain akun aplikasi **gagal** `INSERT` ke tabel kanonik dan **berhasil** `SELECT` lewat view.
**Dasar** DECIDED(ADR-0017). **Ketidakpastian** **REQ-021**.

---

### T-30 — Uji: bentuk view kompatibilitas

**Hasil** Terbukti bahwa hilir menerima bentuk yang sama dengan hari ini, dan menerima selisih pada muatan kedua.

**Prasyarat** T-16, T-32
**Kriteria terima**
- Kolom dan tipe yang dihasilkan sama dengan yang diterima hilir hari ini, termasuk `CASEID`. Dibandingkan terhadap 17 parameter `InputParamOs.*` hasil sapuan S1 — **dibaca, bukan dirancang**.
- **Muatan kedua atas klaim, jenis reasuransi, dan mata uang yang sama menghasilkan selisih, bukan nilai penuh.**
- Muatan yang ditandai ditolak **tidak** ikut mengurangi muatan berikutnya.
**Dasar** EVIDENCED: `SPEC-MODEL-DATA.md` bagian 18. **Ketidakpastian** Pola selisih `InputParamOs.Value - OutOSAcc...` — lihat T-16.

---

### T-31 — Uji: paritas

**Hasil** Terbukti bahwa setiap baris hasil migrasi dapat ditunjuk balik ke barisnya di sistem lama, tepat satu.

**Prasyarat** T-19, T-23
**Kriteria terima** Untuk satu klaim contoh, tiap baris baru punya **tepat satu** pasangan di `MIGRASI_KORELASI`. Kasus uji utama **`CLMNP-975`** — tambalan terbaru, 2026-07-16, EVIDENCED `BLUEPRINT.md` §7.2.
**Dasar** DECIDED(ADR-0005), lewat ADR-0004. **Ketidakpastian** **REQ-018**; **REQ-015** mengukur status buka/tutup 8 klaim bertambalan.

---

### T-32 — Sapuan S1: kolom yang sesungguhnya diterima Arasapas dan kasir — **SELESAI**

**Hasil** 17 parameter `InputParamOs.*` dan 25 parameter `TempKasir.CARI*` terbaca, beserta asal tiap satunya.
**Dasar** EVIDENCED: `Activity\SaveDataToOSAksep_Act.xml`, `Activity\HitServiceToKasir_Act.xml`. Hasil lengkap di `SPEC-MODEL-DATA.md` bagian 18.
**Prasyarat** — *(dijalankan 18 September 2026)*
**Kriteria terima** Terpenuhi. Temuan yang mengubah T-16: keempat nilai uang dikirim sebagai **selisih**, bukan nilai penuh.
**Ketidakpastian** `InsertOSKlaimCNP` **TIDAK DITEMUKAN** sebagai nama berkas; `KonversiKlaim_Act.xml` ada dengan **nol** pasangan properti; `ConnectREST\KonversiKlaimNonLife.xml` belum dibaca isinya. → T-37.

---

### T-33 — Sapuan S2: kelengkapan enum `CNPStatusCase` — **SELESAI**

**Hasil** Empat lapisan disapu; **2 nilai ditulis**, 3 berkas, **nol** rule menguji, dan `"CLAIM ACCEPTED"`/`"CLAIM REJECTED"` **nihil di 279 berkas**.
**Dasar** EVIDENCED. Menguatkan `BLUEPRINT.md` §3.1 — semula berlaku untuk lapisan Activity saja.
**Prasyarat** — *(dijalankan 18 September 2026)*
**Kriteria terima** Terpenuhi. Akibatnya: `STATUS_KLAIM` **tanpa `CHECK`** (T-12).
**Ketidakpastian** Kelengkapan enum hanya dapat dibuktikan dengan membuka folder Komite — **belum diizinkan**.

---

### T-34 — Sapuan S3: perilaku per nilai `PaymentType`

**Hasil** Tabel nilai → perilaku yang dijaganya, disusun dari setiap kondisi dan ekspresi yang mengujinya.
**Dasar** EVIDENCED: `MEMORI_PEMAHAMAN.MD` §7.4 menyebut `1` Final dan `2` Partial/interim, **disimpulkan dari pemakaian**. `Rule-Obj-FieldValue` **tidak ada** di ekspor.
**Prasyarat** —
**Lingkup** Empat lapisan: Activity, Data Transform, Section/Harness/FlowAction, pemetaan RDB/REST.
**Tidak termasuk** **Nama** tiap nilai — itu tetap **A6b**, pertanyaan untuk pemilik proses.
**Kriteria terima** Setiap nilai yang muncul di penjaga terdaftar beserta berkas dan barisnya. Hasil nihil tetap dicatat beserta pola yang dipakai.
**Putusan §4.1 diterima sebagiannya, dan satu bagiannya saya tolak dengan bukti.**

Diterima: **tempatnya dibuat, nilainya menunggu** — tabel referensi bernama yang barisnya kosong sampai arti tiap nilai dikonfirmasi, pola yang sama dengan `TARIF_BERLAKU`. Dan T-34 **berhenti menahan** tiket mana pun.

Ditolak: `CHECK IN ('1'..'7')`. Sapuan empat lapisan atas seluruh 279 berkas menemukan literal berikut, dan hanya ini:

| Nilai | Muncul sebagai | Berkas |
|---|---|---|
| `''` kosong | penjaga | `Activity\ProteksiSendKomiteCNP_Act.xml` |
| `'3'` | penjaga | `Activity\HitServiceToKasir_Act.xml` |
| `'7'` | penjaga | `Activity\ProteksiSendKomiteCNP_Act.xml` |

`1` dan `2` datang dari `MEMORI_PEMAHAMAN.MD` §7.4, yang menyatakannya sendiri **"disimpulkan dari pemakaian"** — bukan dari literal. `4`, `5`, dan `6` **TIDAK DITEMUKAN** di mana pun.

Menulis `CHECK IN ('1'..'7')` karena itu melakukan dua hal yang dilarang sekaligus: **mengarang** keberadaan 4, 5, dan 6; dan **menolak `''`**, yang justru terbukti diuji di kode. `JENIS_PEMBAYARAN` dibiarkan **tanpa `CHECK`**, diperlakukan sama dengan `STATUS_KLAIM` — dan alasannya sama: domain tertutup yang isinya belum lengkap menolak nilai yang sah.

**Ketidakpastian** Sebaran nilai yang sesungguhnya dipakai diukur **REQ-027**. Bila hasilnya menunjukkan himpunan tertutup, `CHECK`-nya dipasang saat itu — satu `ALTER`, bukan pembongkaran.

---

### T-35 — Sapuan S4: adakah padanan IDR untuk biaya penilaian, salvage, dan biaya lain

**Hasil** Jawaban apakah ketiganya memang tanpa pasangan IDR, atau kolomnya yang kurang.
**Dasar** DERIVED dari review; DECIDED(ADR-0007) mewajibkan pasangan untuk setiap nilai uang.
**Prasyarat** —
**Lingkup** Sapuan seluruh properti bersufiks `IDR` pada class `SpreadingRisk`, `Data-Adjustment`, dan `ClaimData`, dibandingkan terhadap daftar properti bernilai uang. Empat lapisan.
**Kriteria terima** Daftar properti uang yang punya padanan IDR dan yang tidak, beserta berkas dan baris.
**Ketidakpastian** **Menahan T-05.** Bila ketiganya perlu pasangan IDR, kolomnya bertambah di `ALOKASI_LAYER` dan `CHECK`-nya ikut terpasang sendiri lewat pembangkit.

---

### T-36 — Sapuan S9: apa persisnya yang ditulis `EditXOLAlokasi`

**Hasil** Daftar properti yang benar-benar dapat disunting petugas.
**Dasar** EVIDENCED: `BLUEPRINT.md` §2.5 — `EditXOLAlokasi` menulis `.IsEditClaim = 1` pada class `ASM-FW-GISFW-Data-SpreadingRisk`, yang di sistem lama memuat baris layer **dan** baris Retensi Cedant.
**Prasyarat** —
**Lingkup** Seluruh `Property-Set` di rule itu, per properti. Empat lapisan.
**Kriteria terima** Daftar properti beserta baris; dan pernyataan tegas properti mana yang **tidak** tersentuh.
**Ketidakpastian** **Menahan T-05.** Memasang `_HITUNG`/`_SUNTING` di kolom yang tidak pernah disunting adalah beban mati; tidak memasangnya di kolom yang disunting berarti membongkar tabel kemudian (ADR-0008).

---

### T-37 — Sapuan S5, S6, S7, S8, dan dua sumber yang belum pernah dibaca

**Hasil** Empat pertanyaan terbuka terjawab atau tercatat nihil beserta polanya.
**Dasar** `_selesai/OPEN-QUESTIONS.md` C5 dan C9; `PENGETAHUAN.md` §15.
**Prasyarat** —
**Lingkup**
- **S5** penulis `AlokasiXOLPaid` — termasuk sebagai sasaran `Page-Copy`, `Page-New`, `RDB-List` ke page, dan hasil Report Definition, bukan hanya `Property-Set`.
- **S6** asal `"Previously Calculated UR"` — **jangan cari literal utuh**; cari potongan `"Previously"` dan `"Calculated"` di dalam ekspresi penyambungan, karena nilainya mungkin dirakit.
- **S7** `Property-Remove` di `CountLossAllocation_act` langkah 10 yang parameter sasarannya kosong — baca langkah itu utuh beserta tetangganya.
- **S8** `IsTreatyIn`, `IsReject`, `IsCloseFile`, `IsAnyAcceptation` — keempatnya belum berlaku pernyataan "bersih".
- **Dua sumber**: `pengetahuan/DDL_Script_ClaimNonProp2.xls` (**belum pernah dibaca**) dan `excludeXML/GetBase64Attachment.xml`.

**Tidak termasuk** H1/H2, I1/I2, F1/F2 — berada di modul Komite.
**Kriteria terima** Tiap sasaran menghasilkan jawaban atau pernyataan nihil **beserta lapisan yang disapu dan pola yang dipakai**, sesuai `PENGETAHUAN.md` §14: yang sah hanya *"tidak ada X di lapisan yang sudah saya baca"*. Temuan yang mengubah putusan lama masuk `_selesai/OPEN-QUESTIONS.md` sebagai D9 dan seterusnya. Hasilnya dibandingkan terhadap `BLUEPRINT.md` §7.5 — tambalan per-case baru adalah temuan tersendiri (ADR-0012).
**Ketidakpastian** S6 menentukan apakah `RETENSI_CEDANT` perlu **keadaan ketiga**; bila ya, kunci (klaim, mata uang) **belum cukup** dan T-05 berubah.

---

### T-38 — Tipe desimal di sisi Golang untuk `NUMBER(38,20)`

**Hasil** Satu putusan: pustaka, presisi yang dinyatakan, dan perilaku saat hasil antara melewati 38 digit signifikan.
**Dasar** DECIDED-TEKNIS(AK-1). DERIVED: `NUMBER(38,20)` adalah plafon mutlak Oracle — 38 digit signifikan, sehingga skala 20 menyisakan tepat **18 digit bulat** dan nol kelonggaran. Untuk nilai rupiah memadai: `CNPLimit` terbesar yang tercatat **Rp 12,5 miliar**.
**Prasyarat** —
**Lingkup** Satu pertanyaan dan satu putusan. **Tanpa kode.**
**Tidak termasuk** Implementasi apa pun — memerlukan spec aplikasi.
**Putusan — DECIDED-TEKNIS, 18 September 2026**

| Hal | Putusan |
|---|---|
| Pustaka | **`cockroachdb/apd`**, dengan konteks presisi **38** |
| Yang ditolak | `shopspring/decimal` |
| `float64` | **tidak boleh muncul di jalur nilai uang mana pun** — termasuk JSON keluar dan lapisan React |

Alasannya bukan selera: `apd` memodelkan presisi dan pembulatan sebagai **konteks yang dinyatakan**; `shopspring` memakai skala yang **mengikuti operasi**. AK-1 menuntut skala dinyatakan sepanjang rantai dengan pembulatan hanya di tepi — pustaka yang skalanya mengikuti operasi membuat ADR-0003 ditegakkan lewat **ingatan penulis kode**, persis yang ditinggalkan sistem lama.

**Kriteria terima**
- Konteks presisi 38 dinyatakan di satu tempat, bukan disebar per pemanggilan.
- Perilaku saat hasil antara melewati 38 digit signifikan dinyatakan tegas: galat, pembulatan, atau perluasan — dan yang dipilih tertulis beserta alasannya.
- Sapuan atas jalur nilai uang menemukan **nol** pemakaian `float64`.

**Ketidakpastian** Mengunci setiap tiket aplikasi sesudahnya. Karena itu ia tiket keputusan, bukan tiket implementasi.

---

## 4. Prasyarat luar — bukan tiket

Menjadikannya tiket membuat papan terlihat bergerak padahal yang bergerak adalah antrean orang lain.

| Prasyarat | Menahan | Sebabnya bukan tiket |
|---|---|---|
| **REQ-012** pemetaan `Data-Admin-DB-Table` | pemetaan migrasi kolom demi kolom | menunggu DBA |
| **REQ-018** duplikat `PYID` dan baris ganda akseptasi | `UNIQUE` di T-23; pola selisih di T-16 | menunggu DBA |
| **REQ-021** hak tulis `POOLDATA` ke tabel Pega | perkiraan biaya migrasi | menunggu DBA; bila hasilnya besar, ia pokok tersendiri untuk manajemen |
| **REQ-033** sebaran `LayerPart` | penjaga agregasi di T-16 | menunggu DBA |
| **REQ-031** cacah baris per entitas | volume migrasi | menunggu DBA |
| **Izin membuka folder `Komite Claim Non Prop`** | H1/H2, I1/I2, F1/F2 | **satu permintaan izin, bukan tiga tiket.** `PENGETAHUAN.md` menyebutnya "satu sesi, satu grep" tiga kali untuk tiga hipotesis berbeda — ia menutup lebih banyak butir per satuan usaha daripada seluruh S1–S9 digabung |
| **A8, A10b, A11b, A12, A13** | isi `TARIF_BERLAKU`; model otorisasi | pertanyaan untuk manusia |
| **Ratifikasi delapan butir AK** | — | `REGISTER-RATIFIKASI.md` menyatakan ratifikasi **tidak menahan pekerjaan** |

---

## 5. Menunggu spec berikutnya

Terlihat dari sini, tetapi memerlukan spec aplikasi — **bukan tiket di dokumen ini**.

| Pekerjaan | Kenapa menunggu |
|---|---|
| Mesin alokasi kerugian | Perhitungan sebagai turunan data (ADR-0013) dan idempotensi (ADR-0011) adalah sifat **prosedur**, bukan sifat baris. Basis data hanya dapat mencegah baris ganda |
| Penegakan suntingan manual bertahan | Basis data menyimpan penanda, pelaku, dan waktu; ia **tidak dapat** mencegah proses hitung menimpanya (ADR-0008) |
| Penularan keadaan menunggu kurs | Aturan prosedur (ADR-0014); T-28 menguji hasilnya, bukan mekanismenya |
| Penyegaran URL dokumen | `EXPDATE` berumur terbatas (ADR-0027); mekanisme penyegarannya lapisan aplikasi |
| Model otorisasi dan tabel pengguna | ADR-0006 RBAC dari nol; dan **ADR tentang pemisahan pengguna, jabatan, dan keanggotaan komite tidak pernah ditulis** — jejaknya hanya satu kurung di `docs/adr/0016-…md:12` |
| Seluruh lapisan Golang dan React | `SPEC-MODEL-DATA.md` menyatakan di luar lingkup |

---

## 6. Swa-periksa

| Pertanyaan | Jawaban |
|---|---|
| Setiap tiket punya kriteria terima yang menguji perilaku, bukan bentuk? | Ya — seluruhnya berbentuk *baris apa diterima, baris apa ditolak*. Tidak ada kriteria berbentuk "kolom X bertipe Y" |
| Setiap tiket punya dasar yang dapat ditunjuk? | Ya — ADR, butir AK, atau berkas beserta bagian/barisnya |
| Ada tiket yang memuat nilai yang dikarang? | Tidak. Nilai yang belum diketahui ditulis **TIDAK DITEMUKAN** atau dirujuk ke REQ-nya |
| Ada tiket Golang atau React? | Hanya **T-38**, dan ia tiket **keputusan tanpa kode** — pengecualian yang dinyatakan di bagian 1 |
| Setiap REQ terbuka yang menyentuh sebuah tiket tercatat di Ketidakpastian tiket itu? | Ya — REQ-012, 015, 018, 020, 021, 027, 031, 032, 033 seluruhnya muncul di tiket yang disentuhnya |
| Setiap uji di Testing Decisions punya tiketnya sendiri? | Ya — delapan baris, delapan tiket T-24 … T-31, tidak digabung |
| Ada putusan yang ditolak, dan alasannya berbukti? | Ya — satu: `CHECK IN ('1'..'7')` pada `JENIS_PEMBAYARAN`, ditolak di T-34 dengan sapuan literal empat lapisan |
| Ada estimasi waktu? | Tidak |
