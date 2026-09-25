---
status: tertahan
---

# 04: Dokumen addendum berdiri bernomor sendiri, dan satu dokumen dapat menyentuh beberapa kontrak

*Asal: `DAFTAR-PEKERJAAN.md` `P-61` · `GRL-19` · `INV-71` · `STRUKTUR-ADDENDUM.md` §3.*

**What to build:** **PK** mencatat sebuah **dokumen addendum** dengan **nomornya sendiri**, dan dokumen itu
dapat menunjuk versi di **beberapa kontrak sekaligus**. Persetujuan **tetap per kontrak** — dokumen
tidak membawa keadaan persetujuan apa pun.

Artefak: entitas `DOKUMEN_ADDENDUM`, hubungannya ke `VERSI_KONTRAK`, dan keunikan nomornya.

**PEMBUAT PERTAMA** untuk `DOKUMEN_ADDENDUM`.

**Persyaratan:** `GRL-19` · `INV-71` (keunikan nomor dokumen) · `KTV-3` (tidak ada penanda "dikirim ke luar"; jenis tetap bernilai dua)

**Tidak termasuk:** **Tanggal berlaku dokumen** — irisan 08.
**Pengisian nomor untuk baris warisan** — irisan 09; migrasi **tidak punya sumber** untuk mengisinya.
**Keadaan persetujuan pada dokumen** — tidak ada, dan itu keputusan: persetujuan melekat pada versi
per kontrak (*"1 kontrak di aksep 1 per 1"*).

**Jalur gagal:** Dua dokumen bernomor sama -> **ditolak** `INV-71`, pesannya menyebut nomor yang bentrok ·
Dokumen tanpa satu pun versi yang ditunjuknya -> ditolak.

**Uji:** **Negatif:** simpan dua dokumen bernomor sama; simpan dokumen tanpa versi.
**Positif — dan ia pokok irisan ini:** satu dokumen menunjuk versi di **tiga kontrak berbeda**
**diterima**, dan ketiga versi tetap berjalan di rantai persetujuannya masing-masing.

**Menggantikan:** `TDA-11` — *picker menyatukan kontrak dan addendum tanpa pembeda dan tanpa saringan
keadaan.* Sistem lama **tidak punya benda bernama dokumen sama sekali**: sapuan dua inventaris
properti atas dua belas akar kata mengembalikan **nol** properti penyimpan nomor dokumen, dan kelas
integrasi addendum memuat tepat 19 properti tanpa satu pun di antaranya.

**Blocked by:** **14** — *ditambahkan 24 September 2026.* Tiket ini menambahkan **kolom** atau **tabel anak** pada `VERSI_KONTRAK`, dan tidak satu pun tiket di papan membuat tabelnya. Irisan `14` adalah **PEMBUAT PERTAMA** `KONTRAK` dan `VERSI_KONTRAK`.

**Dasar:**
```
EVIDENCED(sapuan-properti-dokumen@ekspor-2026-09 - NOL hasil, dikalibrasi atas EDMState/EDMMaterialType/OLDID)
        DECIDED(GRL-19, KTV-3, INV-71)
        DIASUMSIKAN-CLEAR(DB-16a)
```

- [ ] `DOKUMEN_ADDENDUM` berdiri dengan nomornya sebagai kunci alami
- [ ] `INV-71` menolak nomor ganda, pesannya menyebut nomornya
- [ ] satu dokumen menunjuk versi di beberapa kontrak; persetujuan tiap versi **tidak saling menyentuh**
- [ ] dokumen **tidak** membawa kolom keadaan persetujuan — diperiksa, bukan diandaikan
- [ ] `DB-16a` tercatat di `ASUMSI-CLEAR.md`; bila dibantah, jenis bernilai **tiga** dan penanda "dikirim ke luar" kembali
