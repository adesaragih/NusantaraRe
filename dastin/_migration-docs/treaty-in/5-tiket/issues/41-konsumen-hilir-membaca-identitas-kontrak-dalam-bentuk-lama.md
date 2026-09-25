---
status: aktif
golongan: baru
---

# 41: Konsumen hilir yang tidak dikenal tetap dapat membaca identitas kontrak dalam bentuk lama

*Asal: `DAFTAR-PEKERJAAN.md` `P-54` · `ADR-0051` · `SPEC-MODEL-DATA.md` §10.1, §10.2.*

**What to build:** Sebuah **bentuk baca** menyajikan identitas kontrak **dalam bentuk sistem lama** —
pengenal bertipe teks berpola lama, cedant, asal bisnis, periode — sehingga konsumen hilir yang
**belum diketahui siapa saja** tetap dapat membacanya pada hari peralihan.

**Kenapa begini:** `ADR-0051` memerintahkan **anggap ada konsumen hilir**, dan alasannya terbaca dari
ekspor: tabel datar `TREATYINDETAIL` dan `PROPORTIONALARRG` ditulis pada **setiap** penyimpanan,
dengan 52 kolom bertipe `VARCHAR2(1000)`, dan **tidak satu pun aturan di dalam ekspor membacanya
kembali**. Tabel yang ditulis tetapi tidak dibaca oleh sistemnya sendiri **dibaca oleh sesuatu yang di
luar ekspor** — dan kita tidak tahu apa.

> **Diam bukan bukti.** Ketiadaan pembaca di dalam ekspor tidak membuktikan ketiadaan pembaca; ia
> hanya membuktikan pembacanya **tidak ada di sini**.

**Persyaratan:** `ADR-0051` · `INV-59` (hilir diberi **rujukan**, bukan salinan) · `INV-58` (bentuk
baca adalah **turunan yang tidak disimpan**) · `ADR-0042`

**Tidak termasuk:** **Menulis ulang tabel datar lama.** Bentuk baca ini **diturunkan**, bukan
dituliskan aplikasi — `SaveTreatyInDetail_Act` yang menulis nol tetap ke tiga kolomnya **tidak
ditiru**, dan `TETAPAN-DI-KODE.md` §3 sudah menyatakan tabel datar di model baru **diturunkan**.
**Jalur penerbitan ke luar** — `G2` menyatakannya **tegas tidak masuk**: `TREATYINOFFER` tidak punya
satu pun penulis yang terjangkau dan penggantinya belum diketahui (`DAFTAR-ESKALASI-MANAJEMEN.md`
butir 3).
**Isi kontrak selengkapnya dalam bentuk lama** — hanya **identitas** yang dijanjikan di sini. Menjanji
lebih berarti menjanji bentuk yang belum ada yang memintanya.

**Jalur gagal:** Bentuk baca yang **menyimpan** hasilnya sebagai tabel -> melanggar `INV-58`;
ia diturunkan saat dibaca · Pengenal bentuk lama dibangkitkan dari cap waktu -> ditolak `INV-02`; ia
**dibaca dari `NOMOR_KONTRAK_WARISAN`**, dan kontrak yang lahir di sistem baru **tidak punya** — itu
keadaan yang dinyatakan, bukan ditambal.

**Uji:** **Negatif:** minta bentuk lama untuk kontrak yang lahir di sistem baru -> **pernyataan bahwa
nomor lamanya tidak ada**, bukan nomor karangan.
**Positif — dan ia yang membuktikan bentuknya benar-benar sama:** ambil **sepuluh kontrak warisan**,
bandingkan keluaran bentuk baca ini dengan baris `TREATYINDETAIL` lama untuk kontrak yang sama.
**Kesepuluhnya harus identik pada kolom identitas.** Perbandingan terhadap satu kontrak tidak cukup:
sebuah bentuk baca yang benar untuk satu baris dan salah untuk sisanya lulus uji tunggal mana pun.

**Menggantikan:** tidak ada aturan yang digantikan. **Golongannya BARU** — jaminan kompatibilitas
baca memang belum pernah dinyatakan — sehingga `CARA MENYALAKANNYA` wajib.

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-6, CO-7, CO-8** |
| 2 | **siapa membaca, seberapa sering** — **pemilik proses Treaty In, mingguan** selama delapan minggu pertama sesudah peralihan, atas **catatan akses** ke bentuk baca ini. Catatan itulah satu-satunya cara mengetahui **siapa** konsumen hilirnya |
| 3 | **ambang berangka** — bentuk baca ini **dapat dicabut** ketika catatan akses menunjukkan **nol akses selama 8 minggu berturut-turut**. Sampai itu terjadi, ia dipertahankan — dan bila ada akses, pemiliknya **dicari dan dicatat namanya** |
| 4 | **siapa boleh menyalakan atau mencabutnya** — pemilik proses Treaty In |

> Butir 3 sengaja berbentuk **syarat pencabutan**, bukan syarat penyalaan: fitur ini menyala sejak
> hari pertama, dan yang belum diketahui adalah **kapan ia boleh mati**.

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(TREATYINDETAIL@Table/ dan PROPORTIONALARRG@Table/ - ditulis tiap simpan, NOL pembaca di dalam ekspor)
        DECIDED(ADR-0051, ADR-0042, INV-58, INV-59)
```

- [ ] bentuk baca identitas dalam bentuk lama berdiri, dan ia **diturunkan** — bukan tabel yang ditulis aplikasi
- [ ] kontrak yang lahir di sistem baru menghasilkan **pernyataan ketiadaan nomor lama**, bukan nomor karangan
- [ ] uji positif lulus atas **sepuluh** kontrak warisan, bukan satu
- [ ] **catatan akses** berdiri, sebab tanpanya butir 3 `CARA MENYALAKANNYA` tidak dapat diukur
- [ ] keempat butir **CARA MENYALAKANNYA** terisi, termasuk **syarat pencabutan berangka**
- [ ] jalur penerbitan ke luar dinyatakan **di luar irisan ini**, dengan eskalasi butir 3 sebagai penagihnya
