# UJI NEGATIF — membuktikan invarian benar-benar menegakkan

**Tanggal:** 23 September 2026
**Berlaku untuk:** Treaty In
**Sebabnya:** empat invarian rekonsiliasi (INV-47, INV-50, INV-51, INV-31 bentuk lintas baris)
ditegakkan lewat *materialized view* ber-`REFRESH ON COMMIT` dengan `CHECK` padanya. Teknik itu
**belum dibuktikan menegakkan**, dan sampai dibuktikan ia klaim.

---

## 0. Batas berkas ini, dinyatakan di muka

**Uji-uji di bawah BELUM DIJALANKAN.** Tidak ada instance Oracle yang dapat saya capai dari sesi
ini. Yang ada di berkas ini adalah **rancangan ujinya** — apa yang disisipkan, apa yang harus
terjadi, dan bagaimana membaca hasilnya.

Konsekuensinya mengikat dan sudah dicatat di `SPEC-INVARIAN.md` §1:

> Keempat invarian itu berstatus **CONSTRAINT (BELUM DIBUKTIKAN)**. Ia dihitung sebagai constraint
> **hanya setelah** uji negatifnya lulus di Oracle versi yang sama dengan sasaran. Bila salah satu
> gagal, invarian itu **turun ke aplikasi**, dan hitungan langkah 8 berubah **sebelum** dipakai
> sesi DDL.

Objeknya sendiri — tabel, materialized view, constraint — dibuat sesi DDL. Berkas ini tidak membuat
objek apa pun; ia menetapkan apa yang harus dibuktikan tentangnya.

---

## 1. Kenapa uji negatif, dan bukan argumen

Ini instans kelima dari `CONTEXT.md §2.0` — **struktur yang terlihat bukan struktur yang berlaku** —
dan yang pertama **lahir dari rancangan kami sendiri**, bukan diwarisi sistem lama.

Sebuah materialized view yang berhenti me-refresh — karena log-nya kurang kolom, karena kueri
agregatnya melanggar syarat *fast refresh*, karena seseorang mengubah tabel dasarnya, atau karena
seseorang menjalankan `ALTER MATERIALIZED VIEW … REFRESH ON DEMAND` — **berhenti menegakkan apa
pun**. Constraint-nya tetap ada di skema. Ia tetap terbaca sebagai penjaga oleh siapa pun yang
memeriksa katalog, termasuk auditor. Dan ia tidak menjaga.

Empat kegagalan sistem lama yang kami temukan berbentuk persis begitu: aturan yang ada, terbaca,
dan tidak berjalan. Membangun yang kelima sendiri tanpa membuktikannya adalah mengulangi kesalahan
yang seluruh sesi ini dibangun untuk menghindari.

**Ketidakpastian teknis yang harus ikut dijawab uji ini**, dan saya nyatakan sebagai
ketidakpastian, bukan sebagai fakta: apakah bentuk kueri agregat yang dibutuhkan keempat invarian
ini memenuhi syarat `REFRESH FAST ON COMMIT`, atau hanya `REFRESH COMPLETE ON COMMIT`. Keduanya
menegakkan; yang kedua jauh lebih mahal. Uji N-0 menjawabnya sebelum keempat uji lainnya.

---

## 2. Rancangan uji

Setiap uji berbentuk sama, dan bentuk itu yang penting:

1. sisipkan data yang **memenuhi** invarian, `COMMIT` — harus **berhasil**;
2. ubah satu nilai sehingga invariannya **dilanggar**, `COMMIT` — harus **GAGAL**;
3. periksa bahwa setelah kegagalan itu, data di tabel **tidak berubah**.

Langkah 1 bukan basa-basi: uji yang selalu gagal juga "menolak pelanggaran", dan tanpa langkah 1
tidak ada yang membedakan penjaga dari kerusakan.

### N-0 — bentuk refresh mana yang didapat

```sql
-- dijalankan setelah MV dibuat oleh sesi DDL
SELECT MVIEW_NAME, REFRESH_MODE, REFRESH_METHOD, FAST_REFRESHABLE, STALENESS
FROM   USER_MVIEWS
WHERE  MVIEW_NAME LIKE 'MV_REKONSILIASI%';
```

**Yang harus terbaca:** `REFRESH_MODE = COMMIT`. Bila `FAST_REFRESHABLE` bukan `DIRECT_LOAD_AND_DML`
atau setara, MV-nya `COMPLETE ON COMMIT` — sah, tetapi ongkos serialisasinya jauh lebih besar dan
itu harus tercatat di `CATATAN-DDL.md`, bukan ditemukan belakangan.

### N-1 — INV-47: jumlah baris turunan sama dengan besaran induknya

| Langkah | Tindakan | Harus |
|---|---|---|
| 1 | satu versi, satu besaran induk bernilai 100, dua baris turunan bernilai 60 dan 40; `COMMIT` | **berhasil** |
| 2 | ubah satu baris turunan dari 40 menjadi 41; `COMMIT` | **GAGAL** |
| 3 | `SELECT` baris itu | masih 40 |
| 4 | hapus satu baris turunan sehingga jumlahnya 60; `COMMIT` | **GAGAL** |

Langkah 4 penting tersendiri: banyak penjaga menangkap `UPDATE` dan melewatkan `DELETE`.

### N-2 — INV-50: penyebaran berjumlah 100 per sumbu pihak

| Langkah | Tindakan | Harus |
|---|---|---|
| 1 | penyebaran OR 60 + RI 40 pada satu induk; `COMMIT` | **berhasil** |
| 2 | ubah RI menjadi 39; `COMMIT` | **GAGAL** |
| 3 | tambah baris ketiga bernilai 1 sehingga jumlahnya 100 lagi; `COMMIT` | **berhasil** |

Langkah 3 membuktikan penjaganya memeriksa **jumlah**, bukan **cacah baris**.

### N-3 — INV-51: retensi + penyerahan sama dengan 100

| Langkah | Tindakan | Harus |
|---|---|---|
| 1 | retensi 30, penyerahan 70; `COMMIT` | **berhasil** |
| 2 | ubah retensi menjadi 30,000001; `COMMIT` | **GAGAL** |

Langkah 2 sekaligus menguji **presisi**: bila ia lolos, toleransi pembulatan tersembunyi sedang
bekerja di suatu tempat, dan itu temuan tersendiri — ADR-0003 menuntut presisi seragam tanpa
pembulatan di tengah rantai.

### N-4 — INV-31 bentuk lintas baris

| Langkah | Tindakan | Harus |
|---|---|---|
| 1 | baris `QUOTA_SHARE` dengan persen terisi dan lines kosong; `COMMIT` | **berhasil** |
| 2 | isi juga `JUMLAH_LINES_SURPLUS` pada baris yang sama; `COMMIT` | **GAGAL** |
| 3 | baris `SURPLUS` dengan keduanya kosong; `COMMIT` | **GAGAL** |

### N-5 — uji negatif untuk keenam TRIGGER

Bentuk yang sama berlaku pada trigger, dan alasannya lebih kuat: trigger punya sakelar mati.

| # | Invarian | Yang disisipkan | Harus |
|---|---|---|---|
| N-5a | INV-19 | `UPDATE` salah satu kolom kunci alami `KONTRAK` | GAGAL |
| N-5b | INV-22 | perpindahan `DRAFT` → `DISETUJUI` langsung | GAGAL |
| N-5c | INV-23 | perpindahan keluar dari `DISETUJUI` | GAGAL |
| N-5d | INV-24 | `UPDATE` nilai apa pun pada versi `DISETUJUI`, **dan** pada salah satu entitas anaknya | GAGAL keduanya |
| N-5e | INV-28 | perpindahan tanpa baris `CATATAN_PERSETUJUAN` | GAGAL |
| N-5f | INV-54 | addendum bertanggal di luar periode kontraknya | GAGAL |

N-5d menuntut **dua** penyisipan: satu pada baris versinya, satu pada entitas anak. Menegakkan
pembekuan hanya pada baris induk adalah kesalahan yang paling mudah dibuat dan paling mahal —
seluruh nilai kontrak justru ada di anak-anaknya.

---

### N-6 — INV-64…INV-68, lima kunci alami yang dinomori 24 September 2026

Kelimanya **CONSTRAINT**, bukan trigger dan bukan *materialized view*. Karena itu keduanya berlaku:
uji negatifnya **lebih murah** — constraint tidak punya sakelar mati dan tidak dapat basi — tetapi
uji **positifnya** justru yang penting, sebab kesalahan lingkup pada kunci alami **tidak menolak
data yang salah, melainkan menolak data yang BENAR.**

**Maka setiap baris di bawah punya dua uji, dan yang kedua itu yang pernah gagal tiga kali di modul
ini.**

| # | Invarian | Yang disisipkan | Harus |
|---|---|---|---|
| **N-6a** | INV-64 | dua baris `BAGIAN` dengan `ID_LAYER` sama | **GAGAL** |
| **N-6a+** | INV-64 | dua baris `BAGIAN` pada **dua layer berbeda** dalam satu versi | **BERHASIL** |
| **N-6b** | INV-65 | dua `RINCIAN_PENYEBARAN` ber-`ID_JENIS_REASURANSI` sama di bawah satu `ID_PENYEBARAN` | **GAGAL** |
| **N-6b+** | INV-65 | jenis reasuransi **yang sama** di bawah **dua `PENYEBARAN` berbeda** pada versi yang sama | **BERHASIL** — ini yang akan ditolak bila lingkupnya keliru dinaikkan ke versi |
| **N-6c** | INV-66 | dua `PORTOFOLIO` (masuk, premium) pada satu versi | **GAGAL** |
| **N-6c+** | INV-66 | (masuk, premium) dan (keluar, premium) pada satu versi | **BERHASIL** |
| **N-6d** | INV-67 | dokumen yang sama dilampirkan dua kali ke satu versi | **GAGAL** |
| **N-6d+** | INV-67 | dokumen yang sama dilampirkan ke **dua versi** dari satu kontrak | **BERHASIL** — addendum yang merujuk slip yang sama adalah hal biasa |
| **N-6e** | INV-68 | dua baris berkode sama, dijalankan pada **keenam** tabel acuan | **GAGAL** enam kali |

> **Uji berakhiran `+` bukan pelengkap kerapian.** Ketiga kekeliruan lingkup yang sudah terjadi di
> modul ini seluruhnya **lolos uji negatif** — constraint yang terlalu ketat memang menolak
> duplikat. Yang menangkapnya hanya uji yang menyisipkan **data sah** dan menuntutnya **diterima**.
>
> **N-6b+ dan N-6d+ adalah yang paling berharga di seluruh daftar ini**, karena keduanya persis
> bentuk yang pernah salah.

**N-6d dan N-6d+ menunggu satu hal.** INV-67 adalah satu-satunya dari kelima yang lingkupnya
**tidak dapat dibaca dari penulisnya** — lampiran ditangani mesin bawaan Pega yang tidak terekspor
(`L-10`). Ujinya tetap ditulis dan tetap dijalankan; yang berubah bila ekspor kedua datang adalah
**invariannya**, bukan ujinya.

---

---

### N-7 — INV-69, INV-70, INV-71 · diterapkan 24 September 2026 (diff `D-5`)

Dari **`GRL-20`** dan **`GRL-19`**, ronde D grilling Adjustment.

| # | Invarian | Yang disisipkan | Harus |
|---|---|---|---|
| **N-7a** | INV-69 | versi `TIDAK_MATERIAL` + satu baris `NILAI_SELISIH` premi | **GAGAL** |
| **N-7a+** | INV-69 | versi `TIDAK_MATERIAL` + satu baris selisih **tanggal pelaporan** | **BERHASIL** — bukan besaran uang maupun porsi |
| **N-7b** | INV-70 | versi `MATERIAL` + perubahan `PENGECUALIAN` | **GAGAL** |
| **N-7b+** | INV-70 | versi `MATERIAL` + perubahan limit, pengecualian **tidak** disentuh | **BERHASIL** |
| **N-7c+** | INV-69 + INV-70 | **versi WARISAN yang melanggar keduanya** | **BERHASIL DIMUAT** |
| **N-7d** | INV-71 | dua `DOKUMEN_ADDENDUM` bernomor sama | **GAGAL** |
| **N-7d+** | INV-71 | satu dokumen memayungi versi dari **dua kontrak berbeda** | **BERHASIL** |

> ### `N-7c+` adalah yang paling berharga di seluruh berkas ini
>
> Di sistem lama aturan materialitas **tidak pernah ditegakkan di sisi simpan** (`TDA-10`); ia hanya
> penguncian layar, dan penguncian layar **dapat dilewati**. Maka baris warisan **tidak boleh
> diandaikan patuh**.
>
> **Tanpa `N-7c+`, `INV-69` dan `INV-70` akan menolak data lama pada hari peralihan** — dan
> penolakannya baru terlihat ketika migrasi berhenti di tengah. `UA-3` yang mengukur berapa banyak;
> **perlakuan atas yang melanggar adalah keputusan pemilik proses**, bukan keputusan berkas ini.

> **`N-7d+` menjaga arah sebaliknya:** satu dokumen **memang boleh** memayungi kontrak dari cedant
> yang berbeda (`DB-3` dibantah). Uji yang hanya negatif akan lolos meski constraint-nya terlalu
> ketat — dan terlalu ketat di sini berarti **menolak dokumen yang sah**.

**Kadar `N-7a` … `N-7c+`:** keduanya menguji **MV**, dan MV **tidak dapat diuji tanpa instans
Oracle** (`L-3`). Ujinya **ditulis dan belum dijalankan**, sama dengan `N-1` … `N-5f`.
**`N-7d`/`N-7d+` dapat dijalankan** begitu ada instans — ia `UNIQUE` biasa.

---

## 3. Pemantau kebasian — penegakan yang bisa berhenti diam-diam menuntut pengawas

Trigger yang dimatikan dan materialized view yang basi punya sifat yang sama: **keduanya berhenti
menjaga tanpa satu pun kegagalan terlihat**. Constraint biasa tidak punya sifat itu — ia gagal
keras atau tidak ada sama sekali.

Maka keduanya, dan **hanya** keduanya, menuntut pemantau:

```sql
-- MV yang tidak lagi ON COMMIT, atau tidak lagi segar
SELECT MVIEW_NAME, REFRESH_MODE, STALENESS, LAST_REFRESH_DATE
FROM   USER_MVIEWS
WHERE  MVIEW_NAME LIKE 'MV_REKONSILIASI%'
  AND (REFRESH_MODE <> 'COMMIT' OR STALENESS <> 'FRESH');

-- trigger yang tidak menyala
SELECT TRIGGER_NAME, STATUS
FROM   USER_TRIGGERS
WHERE  STATUS <> 'ENABLED';

-- constraint yang tidak menyala atau tidak tervalidasi
SELECT CONSTRAINT_NAME, STATUS, VALIDATED
FROM   USER_CONSTRAINTS
WHERE  STATUS <> 'ENABLED' OR VALIDATED <> 'VALIDATED';
```

**Kedua kueri terakhir sengaja lebih luas dari MV-nya.** `NOVALIDATE` pada sebuah constraint
menghasilkan bentuk yang sama: penjaga yang terbaca di katalog dan tidak pernah memeriksa data yang
sudah ada.

Pemantau ini berjalan terjadwal, dan hasilnya **bukan nol adalah kejadian**, bukan laporan. Siapa
yang menerimanya ditetapkan bersama rencana operasi; yang ditetapkan sekarang adalah **bahwa ia
harus ada**, dan bahwa keempat invarian rekonsiliasi tidak boleh dinyatakan tegak tanpanya.

---

## 4. Ongkos serialisasi, dan apa yang terjadi bila MV dimatikan

Ini dinyatakan sekarang, di muka, karena ia **akan** muncul lagi sebagai keluhan performa dan orang
**akan** tergoda mematikan MV-nya.

**Sifat yang diketahui:** `REFRESH ON COMMIT` menyerialkan commit yang menyentuh kelompok baris yang
sama. Dua orang yang menyunting penyebaran pada versi kontrak yang sama tidak dapat commit
bersamaan; yang kedua menunggu. Pada kontrak yang berbeda tidak ada tunggu-menunggu.

**Bila ia dimatikan — dinyatakan supaya keputusannya sadar:**

| Yang hilang | Akibatnya |
|---|---|
| INV-47 | jumlah baris turunan dapat berbeda dari induknya tanpa ada yang menolak |
| INV-50 | penyebaran dapat berjumlah selain 100, dan kapasitas yang tersebar tidak lagi sama dengan yang dibagikan |
| INV-51 | retensi dan penyerahan dapat tidak berjumlah 100 |
| INV-31 | cabang dapat mengisi kedua kolom sekaligus |

Keempatnya **bukan aturan tampilan**; keempatnya menjaga agar uang yang dibagikan sama dengan uang
yang ada. Sistem lama tidak punya satu pun dari keempatnya, dan itulah sebabnya tabel acuan
kapasitas yang persentasenya tidak berjumlah seratus tidak pernah tertangkap siapa pun (butir 7
daftar eskalasi).

**Maka mematikannya adalah keputusan yang harus diambil dengan nama** — siapa yang memutuskan,
kapan, dan invarian mana yang dilepaskan — bukan tindakan operasional yang diambil diam-diam saat
laporan performa datang. Bila ia dilepaskan, keempatnya **turun ke aplikasi** dan
`PETA-INVARIAN-KE-DDL.md` diperbarui; ia tidak boleh menjadi invarian yang tidak ditegakkan di mana
pun.
