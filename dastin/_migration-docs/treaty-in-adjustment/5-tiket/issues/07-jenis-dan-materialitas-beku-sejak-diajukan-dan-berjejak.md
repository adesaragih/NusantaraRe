---
status: tertahan
---

# 07: Jenis dan materialitas beku sejak diajukan, dan tiap perubahannya meninggalkan jejak

*Asal: `DAFTAR-PEKERJAAN.md` `P-04` (diubah 24 Sep 2026) · `GRL-18` · `KTV-1` · `ADR-0045`.*

**What to build:** Selama versinya `DRAFT`, **PK** dapat membetulkan jenis dan materialitas, dan **tiap
perubahan meninggalkan baris jejak** berisi nilai lama, nilai baru, siapa, dan kapan. Begitu versinya
**diajukan**, keduanya **beku** — mengubahnya menuntut versinya dikembalikan ke `DRAFT` lebih dulu.

Artefak: penegakan titik beku pada perpindahan `AJUKAN`, dan baris `JEJAK_PERUBAHAN` untuk kedua
sumbu.

**Persyaratan:** `GRL-18` · `KTV-1` · `ADR-0045` (jejak perubahan sebagai fakta mesin, tidak dapat dilewati) · `ADR-0055` (perpindahan `KEMBALIKAN`)

**Tidak termasuk:** **Pembekuan seluruh kepala versi** — `P-38` sudah memegangnya pada `DISETUJUI`; irisan ini
hanya memajukan titik beku **untuk dua sumbu**.
**Penguncian ruas oleh materialitas** — irisan 02.

**Jalur gagal:** Mengubah jenis pada versi yang sudah diajukan -> **ditolak**, dan pesannya menyebut bahwa
versinya harus dikembalikan ke `DRAFT` lebih dulu · Mengubah materialitas pada versi yang menunggu
persetujuan -> ditolak · Perubahan pada `DRAFT` tanpa baris jejak -> **tidak mungkin**; jejaknya
fakta mesin.

**Uji:** **Negatif:** ubah jenis sesudah `AJUKAN`; ubah materialitas sesudah `AJUKAN`; ubah dua kali
lalu kembali ke nilai semula dan pastikan jejaknya **tiga baris**, bukan nol.
**Positif:** mengubah jenis pada versi `DRAFT` **berhasil**, dan sesudah `KEMBALIKAN` ia **berhasil
lagi** — pembekuan tidak boleh menjadi penguncian permanen.

**Menggantikan:** `P-04` bunyi lama: *"**PK** dapat mengisi dan mengubah seluruh kepala versi **selama
versinya `DRAFT`**."*

Yang bergeser: pembekuan sesudah `DRAFT` **dinyatakan**, bukan disiratkan; dan perubahan jenis
maupun materialitas **meninggalkan jejak**. Bunyi lama diam tentang apa yang terjadi sesudah `DRAFT`
dan diam sama sekali tentang jejak. Di sistem lama, radio picker **menimpa jenis tanpa syarat apa
pun**, kapan saja, dengan **nol pemeriksaan di sisi simpan** dan **tanpa mencatat nilai
sebelumnya** — sehingga pertanyaan *"jenis apa versi ini ketika kepala seksi menyetujuinya"* **tidak
dapat dijawab untuk satu pun baris warisan**.

**Blocked by:** 02

**Dasar:**
```
EVIDENCED(PickerTreatyInMasterRevisi@ekspor-2026-09, TreatyInEDMSetValue@ekspor-2026-09)
        DECIDED(GRL-18, KTV-1, ADR-0045, ADR-0055)
        DIASUMSIKAN-CLEAR(DB-20)
        DIASUMSIKAN-CLEAR(REV-3)
```

- [ ] perubahan jenis dan materialitas pada `DRAFT` menghasilkan baris `JEJAK_PERUBAHAN` bernilai lama dan baru
- [ ] perubahan sesudah `AJUKAN` **ditolak**, pesannya menyebut jalan keluarnya
- [ ] uji positif lulus: sesudah `KEMBALIKAN` keduanya dapat diubah lagi
- [ ] `DB-20` tercatat di `ASUMSI-CLEAR.md`; bila ia menjawab **beku sejak lahir**, kriteria selesai tiket ini **berbalik**
