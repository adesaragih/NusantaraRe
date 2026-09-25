> Modul  : Komite Claim Non Prop · Tahap 3 · 2026-09-21
> Peran  : juru catat
> Untuk  : DBA
> Masukan: `INVENTARIS-BUKTI.md` §2.5 (enam baris) dan §2.2 (`TANGGAL_CLOSING`) · DDL yang dipegang di `claim-non-prop/pengetahuan/ddl/`
> Status : MENUNGGU HASIL
> Sifat  : TAMBAH-SAJA

# Paket kueri — tujuh permintaan baca

**Seluruhnya baca-saja: nol `INSERT`, nol `UPDATE`, nol `DELETE`, nol DDL, nol Tracer.**
Tidak ada satu pun permintaan di bawah yang mengubah keadaan basis data. Yang terbesar
membaca beberapa ratus baris; sisanya adalah cacah dan tabel acuan kecil.

Nama objek di bawah diambil dari berkas DDL yang sudah kami pegang, bukan dari ingatan.
`ROWNUM` dipakai alih-alih `FETCH FIRST` supaya kuerinya sah di setiap versi Oracle.

Dua dari tujuh **membuka aliran yang saat ini beku** — pekerjaan yang benar-benar tidak dapat
dirancang tanpanya. Lima sisanya hanya mengubah ukuran pekerjaan atau menaikkan derajat
keyakinan. Kolom terakhir menyebutnya tegas.

---

## 1 · Log kiriman pembayaran ke Kasir

**Objek:** `POOLDATA.DIRECTTOKASIR_LOG` — DDL dipegang, isinya **nol baris** pernah diambil.

```sql
-- 1a · berapa banyak, dan sejak kapan
SELECT COUNT(*) AS CACAH,
       MIN(TGL_INPUT) AS PALING_LAMA,
       MAX(TGL_INPUT) AS PALING_BARU
  FROM POOLDATA.DIRECTTOKASIR_LOG;

-- 1b · contoh muatan, 200 baris terbaru
SELECT * FROM (
  SELECT TGL_INPUT, IDPEGA, NOAKSEPTASI, KET, DATA_JSON
    FROM POOLDATA.DIRECTTOKASIR_LOG
   ORDER BY TGL_INPUT DESC
) WHERE ROWNUM <= 200;

-- 1c · berapa kiriman terjadi per satu nomor akseptasi
SELECT NOAKSEPTASI, COUNT(*) AS CACAH_KIRIMAN
  FROM POOLDATA.DIRECTTOKASIR_LOG
 GROUP BY NOAKSEPTASI
HAVING COUNT(*) > 1
 ORDER BY CACAH_KIRIMAN DESC;
```

**Membuka:** bentuk muatan kiriman pembayaran, dan berapa kali kiriman benar-benar terjadi
per akseptasi. **1c menjawab pertanyaan yang tidak dapat dijawab dari rule** — pengirimnya
berada di dalam loop dengan pra-syarat yang tidak menyaring nama treaty, sehingga jumlah
kiriman nyata hanya terbaca dari log.

**Memblokir:** **YA — aliran A-4 beku seluruhnya tanpa ini.**

---

## 2 · Baris akseptasi outstanding

**Objek:** `POOLDATA.OS_AKSEPTASI_KLAIM`, terutama kolom `DATA_JSON` (CLOB) — DDL dipegang,
isinya **nol baris** pernah diambil.

```sql
-- 2a · cacah, dan sebaran penanda pembalikan
SELECT COUNT(*) AS CACAH FROM POOLDATA.OS_AKSEPTASI_KLAIM;

SELECT STS_REJECT, COUNT(*) AS CACAH
  FROM POOLDATA.OS_AKSEPTASI_KLAIM
 GROUP BY STS_REJECT
 ORDER BY 1;

-- 2b · contoh isi, 50 baris terbaru
SELECT * FROM (
  SELECT CASEID, NOCLAIM, TANGGAL, TGL_PROD, STS_REJECT, STS_KONVERSI, ERR_NOTE, DATA_JSON
    FROM POOLDATA.OS_AKSEPTASI_KLAIM
   ORDER BY TANGGAL DESC
) WHERE ROWNUM <= 50;
```

**Membuka:** arti kesembilan argumen prosedur akseptasi, bentuk simpan barisnya, dan
**cakupan pembalikan** — sebaran `STS_REJECT` menunjukkan apakah pembalikan adalah baris
tersendiri atau penanda pada baris yang sama.

**Memblokir:** **YA — aliran A-5 beku seluruhnya tanpa ini.**

---

## 3 · Isi daftar anggota komite

**Objek:** `POOLDATA.EMAILKOMITE` — DDL dipegang, isinya **nol baris** pernah diambil. Tabel
kecil; ambil seluruhnya.

```sql
SELECT ID, NAME, EMAIL, OPERATOR_ID, JABATAN,
       DEGREE, LIMIT_BOTTOM, LIMIT_TOP,
       TYPE_KOMITE, TYPE_BUSINESS, STS_AKTIF,
       STS_REJECT, STS_ADJ, STS_KLAIM
  FROM POOLDATA.EMAILKOMITE
 ORDER BY TYPE_KOMITE, DEGREE, ID;
```

**Membuka:** isi semai daftar anggota pada sistem baru, dan berapa tingkat kewenangan yang
sebenarnya terpakai. `LIMIT_TOP` diambil sistem lama tetapi tidak pernah menyaring apa pun —
kueri ini memperlihatkan apakah kolomnya memang terisi.

**Memblokir:** tidak. Mengubah **ukuran** penyemaian pada tiket `01`, bukan rancangannya.

---

## 4 · Isi badan prosedur penerbit nomor, pada basis data yang berjalan

**Objek:** `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` — yang kami pegang adalah hasil
*reverse-engineer* dari sebuah berkas Excel, **bukan** ekspor langsung. Yang kami minta
adalah isi yang benar-benar berjalan hari ini.

```sql
-- 4a · badan prosedur apa adanya
SELECT LINE, TEXT
  FROM ALL_SOURCE
 WHERE OWNER = 'POOLDATA'
   AND NAME  = 'PROC_GENERATE_SEQUENCE_NUMBER'
   AND TYPE  = 'PROCEDURE'
 ORDER BY LINE;

-- 4b · kapan terakhir diubah, dan apakah sah
SELECT OBJECT_NAME, OBJECT_TYPE, STATUS, CREATED, LAST_DDL_TIME
  FROM ALL_OBJECTS
 WHERE OWNER = 'POOLDATA'
   AND OBJECT_NAME = 'PROC_GENERATE_SEQUENCE_NUMBER';
```

**Membuka:** satu koreksi penilai yang menggantung — lima penutupan kami bersandar pada isi
prosedur ini dan sampai sekarang bertanda *tafsir*. Kueri ini menaikkannya menjadi berbukti,
atau membatalkannya.

**Memblokir:** tidak memblokir tiket mana pun; mekanika penomoran sudah ditetapkan. Ia
menutup keraguan, bukan lubang rancangan.

---

## 5 · Isi tabel penentu periode buku

**Objek:** `POOLDATA.TANGGAL_CLOSING` — dirujuk **di dalam** prosedur penerbit nomor; DDL-nya
tidak kami pegang (`INVENTARIS-BUKTI.md` §2.2, bukan §2.5).

```sql
-- 5a · bentuk tabelnya
SELECT COLUMN_NAME, DATA_TYPE, DATA_LENGTH, NULLABLE
  FROM ALL_TAB_COLUMNS
 WHERE OWNER = 'POOLDATA'
   AND TABLE_NAME = 'TANGGAL_CLOSING'
 ORDER BY COLUMN_ID;

-- 5b · isinya
SELECT * FROM POOLDATA.TANGGAL_CLOSING ORDER BY 1;
```

**Membuka:** **nilai ambang periode buku** — satu-satunya hal yang ditahan pagar pada aliran
penomoran. Mekanikanya sudah ditulis penuh; hanya angkanya yang menunggu.

**Memblokir:** memblokir **satu nilai** di dalam tiket `08`, bukan tiketnya. Pengerja
membangun parameter bernama tanpa isi, dan isinya dipasang saat kueri ini pulang.

---

## 6 · Sebaran jenis usulan yang pernah masuk komite

**Objek:** menurut inventaris, "`ADJUSTMENT` sisi Claim".

**Kami tidak dapat menuliskan kueri siap-jalan untuk ini, dan sebabnya perlu Anda tahu.**
Daftar usulan di sistem lama adalah daftar bersarang milik Pega; pada keempat puluh sembilan
berkas DDL yang kami pegang **tidak ada tabel bernama `ADJUSTMENT`**, dan tabel kerja Pega
`DATAPEGA.PC_ASM_FW_GCNMFW_WORK` **tidak memiliki kolom** untuk jenis usulan — nilainya
berada di dalam kolom gumpalan `PZPVSTREAM`, yang tidak dapat dibaca `SELECT` biasa.

Yang kami minta lebih dulu adalah pemeriksaan katalog, bukan tebakan:

```sql
-- 6a · adakah kolom yang membawanya
SELECT OWNER, TABLE_NAME, COLUMN_NAME, DATA_TYPE
  FROM ALL_TAB_COLUMNS
 WHERE OWNER IN ('DATAPEGA','POOLDATA')
   AND (COLUMN_NAME LIKE '%ADJUST%' OR COLUMN_NAME LIKE '%KOMITE%'
        OR COLUMN_NAME LIKE '%COMMITTE%')
 ORDER BY OWNER, TABLE_NAME, COLUMN_NAME;
```

**Membuka:** apakah pertanyaan ini dapat dijawab dengan SQL sama sekali. Bila jawabannya
tidak, ia bukan permintaan ke DBA melainkan permintaan ekspor dari admin Pega — dan kami
akan mengajukannya sebagai hal yang berbeda.

**Memblokir:** tidak.

---

## 7 · Cacah klaim yang membawa nomor sirkulasi bukan miliknya

**Objek:** menurut inventaris, "`KLAIM` × case Komite".

**Kendala yang sama dengan nomor 6, dan lebih tajam.** Tabel kerja Pega memang mengekspos
`KOMITECOUNT` dan `FLAGONGOINGCOMMITTE` sebagai kolom, tetapi **tidak** nomor sirkulasi yang
ditulis ke klaim. Kueri 6a di atas juga menjawab bagian ini.

Yang **dapat** dijawab SQL hari ini adalah pertanyaan bertetangga — bukan pengganti, dan kami
menyebutnya begitu:

```sql
-- 7a · klaim yang penanda komitenya menyala tetapi tidak punya berkas komite anak
SELECT COUNT(*) AS CACAH
  FROM DATAPEGA.PC_ASM_FW_GCNMFW_WORK K
 WHERE K.PXOBJCLASS NOT LIKE '%KomiteTreatyNonProp'
   AND (K.FLAGONGOINGCOMMITTE IS NOT NULL OR K.KOMITECOUNT > 0)
   AND NOT EXISTS (
       SELECT 1
         FROM DATAPEGA.PC_ASM_FW_GCNMFW_WORK A
        WHERE A.PXOBJCLASS LIKE '%KomiteTreatyNonProp'
          AND A.PXCOVERINSKEY = K.PXCOVERINSKEY);
```

**Membuka:** ukuran pekerjaan migrasi data. Pertanyaan aslinya — berapa klaim membawa nomor
sirkulasi **milik berkas lain** — menunggu jawaban 6a.

**Memblokir:** tidak. Bentuk migrasi data sudah diputuskan: tali klaim ke sirkulasi dibaca
dari sisi sirkulasi, dan nomor yang tersimpan di klaim lama tidak dipercaya. Cacah ini
mengubah **ukuran**, bukan rancangan.

---

## Ringkas

| # | Yang diambil | Objek | Memblokir |
|---|---|---|---|
| 1 | Log kiriman pembayaran | `POOLDATA.DIRECTTOKASIR_LOG` | **YA — A-4** |
| 2 | Baris akseptasi outstanding | `POOLDATA.OS_AKSEPTASI_KLAIM` | **YA — A-5** |
| 3 | Daftar anggota komite | `POOLDATA.EMAILKOMITE` | tidak — ukuran penyemaian |
| 4 | Badan prosedur penerbit nomor | `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` | tidak — menutup keraguan |
| 5 | Penentu periode buku | `POOLDATA.TANGGAL_CLOSING` | satu **nilai** di tiket `08` |
| 6 | Sebaran jenis usulan | katalog dulu — kolomnya mungkin tidak ada | tidak |
| 7 | Cacah nomor sirkulasi yatim | katalog dulu — kolomnya tidak ada | tidak — ukuran migrasi |

Nomor 1 dan 2 yang kami minta lebih dulu bila harus memilih: keduanya menahan seluruh aliran,
dan sisanya tidak menahan apa pun.
