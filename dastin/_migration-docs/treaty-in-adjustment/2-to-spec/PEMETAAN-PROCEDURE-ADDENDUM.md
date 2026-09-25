# Pemetaan dua prosedur addendum — ke mana tiap potong logikanya pergi

**Tanggal:** 24 September 2026 · **Langkah 7 to-spec**
**Sumber:** `../../treaty-in/PROCEDURE/PEGA_M_TREATY_IN_EDM.txt` (150 baris) dan
`PEGA_M_TREATY_IN_DETAIL_EDM.txt` (152 baris) — dibaca seluruhnya
**`SISI`** — **IRISAN.** Kedua prosedur ada di folder `PROCEDURE/` modul **induk**; temuannya milik
Treaty In, diadili di sini untuk jalur addendum saja dan diusulkan ke induk.

> **Tiga golongan, tidak ada yang keempat:**
> **PINDAH KE GOLANG** — aturan bisnis · **TURUN JADI CONSTRAINT** — penjagaan integritas ·
> **TIDAK DIBAWA** — cacat atau perancah, **dan cacatnya disebut**.

---

## 1. `PEGA_M_TREATY_IN_EDM` — penyimpan addendum

| # | Potong logika | Golongan | Keterangan |
|---|---|---|---|
| 1 | **`IF STSINPUT = '0'`** memilih sisip atau perbarui | **PINDAH KE GOLANG** | keputusan *"addendum ini sudah ada atau belum"* adalah aturan bisnis. Dan penandanya **berbeda dari saudaranya**: kontrak memakai `IDPega = 'UnknownId'`, addendum memakai parameter tersendiri — **dua cara menyatakan hal yang sama**, di dua prosedur yang disalin dari satu cetakan |
| 2 | `INSERT` ke `M_TREATY_IN_EDM` **dan** `TREATY_IN_EDM` dalam satu transaksi | **TURUN JADI CONSTRAINT** | yang dijaga: baris datar tidak boleh ada tanpa dokumennya. Di skema baru **tidak ada dua bentuk**, jadi penjagaannya menjadi kunci asing biasa |
| 3 | **`EDMDATE = SYSDATE` pada sisip DAN pada perbarui** | **TIDAK DIBAWA — cacat** | kolom yang namanya berbunyi *tanggal addendum* sebenarnya **tanggal sentuh terakhir**. Setiap penyimpanan ulang menggesernya, sehingga **tanggal pembuatan addendum tidak dapat dipulihkan dari kolom ini**. Di model baru, dibuat dan diubah adalah **dua fakta** (ADR-0006). **`EVIDENCED` — dan sapuan penulis membenarkannya: `EDMDate` punya NOL penulis di Pega**; ia hanya diisi prosedur ini |
| 4 | **`OLDID` diperbarui setiap kali** | **TURUN JADI CONSTRAINT** | penunjuk versi pendahulu seharusnya **beku** sesudah versi lahir. Di skema baru ia `ID_VERSI_KONTRAK_DASAR`, dan keberubahannya **dilarang** — `GRL-10` menyimpannya **tepat di satu tempat** |
| 5 | tidak ada `M_SITE_DATABASE`, tidak ada urutan | **PINDAH KE GOLANG** | pengenal addendum **dirakit di Pega**, bukan di basis data — itulah sebabnya bentuk `…/Rnn` sampai ke sini sebagai teks |
| 6 | **`replace(DataPega,'UnknownId',id)`** | **TIDAK DIBAWA — cacat** | penyulihan **teks buta atas seluruh dokumen JSON**; mengganti setiap kemunculan kata itu di mana pun, termasuk di dalam nilai teks yang kebetulan memuatnya |
| 7 | **`ErrMsg := 'Insert Rate Error'`** pada penyimpan **addendum** | **TIDAK DIBAWA — cacat** | pesan galat menyebut modul yang salah. Orang yang mencarinya di log akan mencari di modul yang keliru |
| 8 | **`COMMIT` di dalam prosedur** | **TURUN JADI CONSTRAINT — sebagai LARANGAN** | batas transaksi dimiliki **pemanggil**. Ini yang membuat cacat §2 tidak dapat dibatalkan |
| 9 | **`EXCEPTION WHEN OTHERS` tanpa `RAISE`, tiga tingkat** | **TIDAK DIBAWA — cacat** | menelan **setiap** galat, termasuk yang tidak diantisipasi |

## 2. `PEGA_M_TREATY_IN_DETAIL_EDM` — proyeksi datar addendum

### 2.1 Cacat pertama: **tidak ada cabang `ELSE`**

```sql
IF IDPega = 'UnknownId' THEN
   … sisip …
END IF;

StsSave := 1;
IDPegaOut := id_TreatyIn_out;
```

Bila `IDPega` **bukan** `'UnknownId'` — yaitu setiap kali baris rincinya **sudah ada dan hendak
diperbarui**:

* **tidak satu pun pernyataan dijalankan**;
* `id_TreatyIn_out` **tidak pernah diberi nilai**, sehingga `IDPegaOut` keluar **`NULL`**;
* dan **`StsSave := 1` tetap dijalankan** — **berhasil**.

> **Golongan: TIDAK DIBAWA — cacat.** Ia sekeluarga dengan `PEGA_TREATY_IN` yang dipanggil dari layar
> addendum dan menghasilkan `UPDATE … WHERE ID = 'XXXXXXX/Rnn'` — **nol baris berubah, tanpa galat**,
> lalu melapor *"Data Sudah Disimpan Dengan ID : …"*.
>
> **Pekerjaan yang tidak dilakukan, dilaporkan sebagai selesai.**

**Uji datanya:** apakah ada `TREATYID` dengan lebih dari satu baris rinci untuk kombinasi
layer/kelompok yang sama — **Uji AN**, sudah masuk lampiran pengukuran.

### 2.2 Cacat kedua: **tujuh kolom tertinggal** — `TDA-17`

`TREATYINDETAILEDM` adalah **himpunan bagian murni** dari `TREATYINDETAIL`. Yang hilang:

```
DEDUCTIBLE   DEDUCTIBLE2   ADJ_RATE   PREMIUM_EARNED   MDP_PCT   ROL_PCT   MDP
```

Ketujuhnya besaran non-proporsional, dan **ketujuhnya justru yang paling mungkin berubah lewat
addendum**. Adjudikasinya `GRILL-D/01-TEMUAN.md` `TD-05`; tiga di antaranya **lubang induk**, bukan
pembuangan.

### 2.3 Sisanya

| # | Potong logika | Golongan |
|---|---|---|
| 1 | `SELECT COMMENCEMENT, TERMINATION FROM TREATY_IN_EDM` lalu `to_date(…,'YYYYMMDD')` | **PINDAH KE GOLANG** — dan di skema baru **tidak perlu diambil**, karena rincinya **berinduk** pada versinya |
| 2 | **`WHEN NO_DATA_FOUND THEN v_COMMENCEMENT := NULL`** | **TIDAK DIBAWA — cacat** | induk yang tidak ditemukan menjadi **tanggal kosong**, bukan penolakan. Baris rinci yatim tetap tersimpan. Di skema baru: **kunci asing**, dan ia menolak |
| 3 | `INSERT` dua tabel dalam satu transaksi | **TURUN JADI CONSTRAINT** |
| 4 | **`ErrMsg` dideklarasikan `OUT` dan tidak pernah diberi nilai sama sekali** | **TIDAK DIBAWA — cacat** — berbeda dari §1 butir 7 yang setidaknya **salah**; yang ini **tidak ada** |

---

## 3. Rekap

| Prosedur | Pindah | Constraint | **Tidak dibawa** |
|---|---:|---:|---:|
| `PEGA_M_TREATY_IN_EDM` | 2 | 3 | **4** |
| `PEGA_M_TREATY_IN_DETAIL_EDM` | 1 | 1 | **4** |

> **Perintah pemilik proses: sistem baru tidak memanggil stored procedure.** ADR-0056 sudah
> menetapkannya untuk seluruh skema `TREATY_MASUK`. **Nol `CREATE PROCEDURE` di keluaran mana pun**,
> dan **nol kolom blob** — diperiksa di `../../treaty-in/2-to-spec/ddl-usulan/`.

**Kedua prosedur terpetakan habis.** Tidak ada potong logika yang tidak bergolongan.
