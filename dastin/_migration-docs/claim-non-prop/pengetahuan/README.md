# pengetahuan/

Bahan yang **masuk**, bukan kesimpulan yang **keluar**.

Aturan isinya satu kalimat: **kalau isinya tidak saya karang, ia di sini.** Sumber mentah, data hasil ekstraksi, dan alat untuk mengambil data. Kesimpulan — ADR, FINDING, BLUEPRINT, CONTEXT, RINGKASAN — tinggal di folder induk, karena isinya adalah tafsir atas bahan di sini, dan tafsir harus bisa dibantah oleh bahannya.

Dipisahkan 18 September 2026. Seluruh rujukan di folder induk sudah disetel ke jalur baru; 16 jalur unik, semuanya terverifikasi menunjuk ke berkas yang ada.

---

## Isi

### Sumber — diterima apa adanya, tidak diubah

| Berkas | Asal | Keadaan |
|---|---|---|
| `DDL_Script_ClaimNonProp.xls` | ditempel pengguna 2026-09-18 10:36 | **terpakai** — 48 objek `POOLDATA`, terurai ke `ddl/` |
| `DDL_Script_ClaimNonProp2.xls` | muncul di folder 2026-09-18 10:57 | **belum pernah dibaca** — tidak disebut dalam instruksi mana pun |

### Hasil uraian sumber

| Berkas | Isi |
|---|---|
| `ddl/` | 49 berkas `.sql`, satu objek satu berkas. Setiap berkas berkepala `-- UNTUK DIBACA, BUKAN UNTUK DIJALANKAN`. Klausa `DROP TABLE` dan `ALTER ... DROP PRIMARY KEY` sudah dibuang |
| `ddl/TABLE_PC_ASM_FW_GCNMFW_WORK.sql` | perkecualian: **bukan** dari berkas `.xls`, ditempel langsung oleh pengguna. Tabel work Pega — **tidak dimigrasi** |
| `SCHEMA-ACTUAL.csv` | 698 baris. 598 dari `DDL_Script_ClaimNonProp.xls`, 100 dari DDL tabel work. Dibedakan lewat kolom `sumber` dan `req_id`. Pemisah `;` |

### Data mentah hasil sapuan XML

| Berkas | Isi | Dipakai oleh |
|---|---|---|
| `arithmetic-inventory.tsv` | 673 ekspresi aritmetika beserta skala pembagiannya | `BLUEPRINT.md`, `_selesai/OPEN-QUESTIONS.md` E1 dan E2 |
| `rekonsiliasi-kolom-vs-properti.tsv` | 408 kolom DDL × 659 properti Pega; 110 berpasangan, 298 kolom yatim, 549 properti yatim | `BLUEPRINT.md` §19 |

Keduanya **angka batas atas**, bukan angka pasti — dasarnya pencocokan nama.

### Alat — belum satu pun pernah dijalankan

| Berkas | Guna | Keadaan |
|---|---|---|
| `PULL-LIST.csv` | daftar tarikan 74 objek: 67 `CONFIRMED` (nama terbaca literal di SQL produksi), 7 `DERIVED` (ditebak dari nama class Pega). 19 objek masih ber-`OWNER` tanda tanya | **kanonik.** 49 dari 74 sudah punya DDL di `ddl/`; 25 belum |
| `PREFLIGHT.sql` | 335 baris, tujuh blok T1–T7 — memastikan versi, hak akses, dan keberadaan objek sebelum query apa pun ditulis | T1–T3 terjawab dari DDL tanpa dijalankan; **T4, T5, T7 masih menunggu** |
| `RECON.sql` | 433 baris — memulihkan tipe dan presisi properti, memetakan objek Oracle. Hasilnya ditempel ke `SCHEMA-ACTUAL.csv` | **belum bisa dipakai** — kepalanya memuat `SET PAGESIZE`/`SET LONG`, perintah SQL\*Plus, bukan TOAD. Perlu ditulis ulang |

---

## Batas yang berlaku untuk seluruh folder ini

`RECON.sql` menetapkannya di kepalanya sendiri, sebelum ada satu baris data pun:

```
READ-ONLY. Tidak ada DML. Sampel maks 20 baris.
Kolom teks bebas di-masking. Tanpa data nasabah.
```

Dan yang berlaku untuk `ddl/`: **untuk dibaca, bukan untuk dijalankan.**
