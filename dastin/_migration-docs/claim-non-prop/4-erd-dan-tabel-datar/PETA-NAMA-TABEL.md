# PETA NAMA TABEL — `T_CLAIMNP_` dihapus, `T_CLAIM_` berlaku

> **Dasar**: `Diagram-Skema-Tabel-NusantaraRe.xlsx`, sheet *Daftar Relasi*, catatan bawah — *"[keputusan work owner] 20 September 2026 — awalan `T_CLAIMP_` DIGANTI `T_CLAIM_` di seluruh berkas ini."* Berkas itu berubah 20-09-2026 pukul 12:54 di luar sesi ini; rancangan Claim Non Prop menyesuaikan diri terhadapnya, bukan sebaliknya.

Penanda lini dilepas dari tabel klaim. Tabel bersama (`T_WORK_CLAIM`, `T_GENERAL_CLAIM`, `T_GENERAL_KOMITE`, `T_KOMITE_KOMITELIST`, `T_VIEW_SUGGEST`, `DOCUMENT_CLAIM`) **tidak berubah sama sekali**. Tabel Claim Life (`T_CLAIMLF_`) **tidak disentuh** — begitu pula berkas rujukan.

## A · Nama lama → nama baru

| Nama lama | Nama berlaku | Padanan di `ddl-usulan/` |
|---|---|---|
| `T_CLAIMNP_ADJUSTMENT` | **`T_CLAIM_ADJUSTMENT`** | AKSEPTASI + ADJUSTMENT |
| `T_CLAIMNP_ADJ_LAYER` | **`T_CLAIM_ADJ_LAYER`** | ⚠ TIDAK ADA di KLAIMNP |
| `T_CLAIMNP_ADJ_LAYER_CURRENCY` | **`T_CLAIM_ADJ_LAYER_CURRENCY`** | ⚠ TIDAK ADA di KLAIMNP |
| `T_CLAIMNP_ADJ_LOSS_ALLOCATION` | **`T_CLAIM_ADJ_LOSS_ALLOCATION`** | ALOKASI_LAYER (baris beku) |
| `T_CLAIMNP_ADJ_QUOTA_SHARE` | **`T_CLAIM_ADJ_QUOTA_SHARE`** | PENYEBARAN (baris beku) |
| `T_CLAIMNP_ADJ_SPREADING` | **`T_CLAIM_ADJ_SPREADING`** | PENYEBARAN (baris beku) |
| `(baru — dipecah dari T_CLAIMNP_SPREADING)` | **`T_CLAIM_BREAK_QS`** | PENYEBARAN · bagian QS |
| `T_CLAIMNP_CLAIM_AMOUNT` | **`T_CLAIM_CLAIM_AMOUNT`** | NILAI_KLAIM_MATA_UANG |
| `T_CLAIMNP_CLOSING_DATE` | **`T_CLAIM_CLOSING_DATE`** | TUTUP_BUKU |
| `T_CLAIMNP_CORRECTION` | **`T_CLAIM_CORRECTION`** | KOREKSI_NILAI |
| `T_CLAIMNP_DATE_GUARD` | **`T_CLAIM_DATE_GUARD`** | KLAIM_PENJAGA_TANGGAL |
| `T_CLAIMNP_ESTIMATION` | **`T_CLAIM_ESTIMATION`** | ESTIMASI_AWAL |
| `(baru — dipecah dari T_CLAIMNP_OBJECT)` | **`T_CLAIM_INTEREST`** | OBJEK_PERTANGGUNGAN · bagian interest |
| `T_CLAIMNP_MIG_CORRELATION` | **`T_CLAIM_MIG_CORRELATION`** | MIGRASI_KORELASI |
| `T_CLAIMNP_MIG_LANDING` | **`T_CLAIM_MIG_LANDING`** | MIGRASI_PENDARATAN |
| `T_CLAIMNP_MIG_REJECTED` | **`T_CLAIM_MIG_REJECTED`** | MIGRASI_NILAI_DITOLAK |
| `T_CLAIMNP_OBJECT` | **`T_CLAIM_OBJECT`** | OBJEK_PERTANGGUNGAN |
| `(baru — dipecah dari T_CLAIMNP_OBJECT)` | **`T_CLAIM_OBJECT_ITEM`** | OBJEK_PERTANGGUNGAN · bagian item |
| `T_CLAIMNP_OUTBOUND_ARCHIVE` | **`T_CLAIM_OUTBOUND_ARCHIVE`** | ARSIP_MUATAN_KELUAR |
| `T_CLAIMNP_RATE` | **`T_CLAIM_RATE`** | TARIF_BERLAKU |
| `T_CLAIMNP_RECEIVER` | **`T_CLAIM_RECEIVER`** | REKENING_PENERIMA |
| `T_CLAIMNP_REINSTATEMENT` | **`T_CLAIM_REINSTATEMENT`** | PREMI_PEMULIHAN |
| `T_CLAIMNP_RETENTION_CEDANT` | **`T_CLAIM_RETENTION_CEDANT`** | RETENSI_CEDANT |
| `T_CLAIMNP_RETRO` | **`T_CLAIM_RETRO`** | ⚠ TIDAK ADA di KLAIMNP |
| `T_CLAIMNP_SPREADING` | **`T_CLAIM_SPREADING`** | PENYEBARAN |
| `T_CLAIMNP_SPREADING_RISK` | **`T_CLAIM_SPREADING_RISK`** | ALOKASI_LAYER |
| `T_CLAIMNP_SPREAD_LOSS` | **`T_CLAIM_SPREAD_LOSS`** | PEMBAGIAN_KERUGIAN |

**27 tabel** `T_CLAIM_`. Dihitung dari `peta-nama-tabel.tsv`, yang dibangkitkan `alat/buat-skema-claimnp.py` bersama workbook-nya.

## B · Yang kini DIPAKAI BERSAMA dengan Claim Prop / Claim Fac In

Melepas penanda lini berarti tabel-tabel ini **satu tabel untuk beberapa lini**, bukan tabel bernama mirip. Induk dan kunci tamunya sudah diadu dengan berkas rujukan: **nihil** yang berbeda.

| Tabel | Lini pemakai |
|---|---|
| `T_CLAIM_ADJUSTMENT` | PROP · FAC IN · NON PROP |
| `T_CLAIM_ADJ_LOSS_ALLOCATION` | PROP · NON PROP |
| `T_CLAIM_ADJ_QUOTA_SHARE` | PROP · FAC IN · NON PROP |
| `T_CLAIM_ADJ_SPREADING` | PROP · FAC IN · NON PROP |
| `T_CLAIM_BREAK_QS` | PROP · FAC IN · NON PROP |
| `T_CLAIM_CLAIM_AMOUNT` | PROP · NON PROP |
| `T_CLAIM_ESTIMATION` | PROP · FAC IN · NON PROP |
| `T_CLAIM_INTEREST` | PROP · NON PROP |
| `T_CLAIM_OBJECT` | FAC IN · NON PROP |
| `T_CLAIM_OBJECT_ITEM` | FAC IN · NON PROP |
| `T_CLAIM_SPREADING` | PROP · FAC IN · NON PROP |

**11 tabel** dipakai bersama · **16 tabel** khas Claim Non Prop (`T_CLAIM_ADJ_LAYER`, `T_CLAIM_ADJ_LAYER_CURRENCY`, `T_CLAIM_CLOSING_DATE`, `T_CLAIM_CORRECTION`, `T_CLAIM_DATE_GUARD`, `T_CLAIM_MIG_CORRELATION`, `T_CLAIM_MIG_LANDING`, `T_CLAIM_MIG_REJECTED`, `T_CLAIM_OUTBOUND_ARCHIVE`, `T_CLAIM_RATE`, `T_CLAIM_RECEIVER`, `T_CLAIM_REINSTATEMENT`, `T_CLAIM_RETENTION_CEDANT`, `T_CLAIM_RETRO`, `T_CLAIM_SPREADING_RISK`, `T_CLAIM_SPREAD_LOSS`).

## C · Tiga tabel yang DIPECAH

Pemecahan ini bukan kerapian — ia membuat rancangan cocok dengan `struktur-claimdata-lama.md` **dan** dengan berkas rujukan sekaligus.

| Dulu | Kini | Sebab |
|---|---|---|
| `T_CLAIMNP_SPREADING` memuat `SpreadingClaim[]` **dan** `SpreadingBreakQS[]` dengan pembeda `JENIS_REASURANSI` | `T_CLAIM_SPREADING` + `T_CLAIM_BREAK_QS` | di struktur lama keduanya halaman terpisah (`SpreadingBreakQS[]` n=28), dan berkas rujukan sudah punya `T_CLAIM_BREAK_QS` sebagai tabel sejajar |
| `T_CLAIMNP_OBJECT` memuat `InterestList[]` **dan** `ObjectList[].ObjectItemList[]` | `T_CLAIM_INTEREST` + `T_CLAIM_OBJECT` + `T_CLAIM_OBJECT_ITEM` | di struktur lama `InterestList[]` (n=39) berdiri sendiri, dan `ObjectItemList[]` bersarang di dalam `ObjectList[]`. Berkas rujukan memakai ketiga nama itu persis |

## D · Yang TIDAK ikut berubah

| Hal | Keadaan |
|---|---|
| Skema Oracle `KLAIMNP` | **tetap**. Ia nama skema, bukan awalan tabel |
| Nama tabel di `ddl-usulan/` | **tetap** berbahasa Indonesia (`KLAIM`, `AKSEPTASI`, `NILAI_KLAIM_MATA_UANG`, …). Penerjemahnya ada di `ERD-ORACLE.xlsx` sheet **PETA-NAMA-T** dan di `peta-nama-tabel.tsv` |
| `T_CLAIMLF_*` (Claim Life) | **tetap**, sesuai keputusan yang sama |
| Nama berkas | **tetap**. Modulnya masih Claim Non Prop; yang berubah penamaan tabelnya |

## E · Satu nama yang mudah tertukar

`T_CLAIM_RETRO` (khas Non Prop, induk `T_CLAIM_SPREAD_LOSS`, dari `CNPSpreadLoss[].RetroList[]`) **bukan** `T_CLAIM_FAC_RETRO` (Prop dan Fac In, dari `ClaimData.FacRetroList`). Nama mirip, halaman asal berbeda, induk berbeda. Keduanya berdiri sendiri.

Tabel yang ada di berkas rujukan tetapi **tidak** dipakai Non Prop: `T_CLAIM_FAC_RETRO`.

---

Dibangkitkan dari `peta-nama-tabel.tsv` dan `Diagram-Skema-Tabel-NusantaraRe.xlsx`. Bila berkas ini berbeda dari keduanya, **alatnya yang salah**.
