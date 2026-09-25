---
status: aktif
---

# 57: Angka rupiah pada versi yang sudah disetujui tidak bergeser, hari ini maupun tahun depan

*Asal: `DAFTAR-PEKERJAAN.md` `P-27` · `ADR-0036` · `INV-43`.*

**What to build:** Kurs **dibekukan pada versinya** saat disetujui, dan **tidak dibaca ulang** saat
ditampilkan. Angka rupiah sebuah versi `DISETUJUI` sama hari ini dan tahun depan.

**Persyaratan:** `ADR-0036` (angka dasar persetujuan beku; hasil beku membawa penunjuk ke masukan yang dipakai) · `INV-43`

**Tidak termasuk:** **Kurs pada versi `DRAFT`** — ia memang boleh mengikuti kurs berjalan sampai disetujui.
**Perhitungan ulang** — `ADR-0043`: itu peristiwa bisnis tersendiri, bukan bagian tiket ini.

**Jalur gagal:** Mengubah baris kurs tahunan lalu membuka versi `DISETUJUI` -> angkanya **tidak berubah** ·
Bila berubah, tiket ini belum selesai.

**Uji:** **Negatif:** ubah kurs master, buka versi disetujui, bandingkan — **harus identik**.
**Positif:** versi `DRAFT` **mengikuti** kurs terbaru; pembekuan tidak boleh bocor ke draf.

**Menggantikan:** **Angka yang sudah disetujui dapat bergeser sendiri** — salah satu dari empat sebabnya
**kurs tahun berjalan**. Sistem lama membaca kurs saat menampilkan, sehingga nilai rupiah sebuah
kontrak yang sudah disetujui **berubah setiap ganti tahun** tanpa ada yang menyentuhnya.

**Blocked by:** `46` · `20`

**Dasar:**
```
EVIDENCED(TREATYEXCHANGEYEARLY@ekspor-2026-09 - kurs per tahun dibaca saat tampil)
        DECIDED(ADR-0036, INV-43)
```

- [ ] kurs tersimpan pada versi saat disetujui, beserta tanggal dan sumbernya
- [ ] uji negatif kurs-master-berubah lulus: angka versi disetujui **identik**
- [ ] uji positif lulus: versi `DRAFT` tetap mengikuti kurs terbaru
