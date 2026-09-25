---
status: aktif
---

# 13: Perubahan persentase pemulihan limit menghasilkan baris selisih, dan versi tidak-material yang mengubahnya ditolak

*Asal: `DAFTAR-PEKERJAAN.md` `P-64` · `TDA-18` · `INV-69` · `STRUKTUR-ADDENDUM.md` §2.2.*

**What to build:** Mengubah persentase pemulihan limit pada sebuah versi penyesuaian **menghasilkan baris
`NILAI_SELISIH`** — sesuatu yang sistem lama tidak pernah lakukan.

Akibat langsungnya, dan ia yang harus diberitahukan sebelum peralihan:

> Versi yang mengubah persentase pemulihan limit **dan dinyatakan `TIDAK_MATERIAL`** **ditolak saat
> simpan** oleh `INV-69` — sebab perubahan itu kini menghasilkan baris selisih bertipe **porsi**.
> Di sistem lama hal yang sama **lolos tanpa sepatah galat**.

**Persyaratan:** `TDA-18` · `INV-69` · bahan to-spec `E-1a`

**Tidak termasuk:** **Perhitungan ulang atas versi warisan** — `ADR-0043`: menghitung ulang adalah peristiwa
bisnis tersendiri, bukan bagian migrasi. Baris warisan **tidak** diklasifikasi ulang.

**Jalur gagal:** Versi `TIDAK_MATERIAL` yang mengubah persentase pemulihan -> **ditolak saat simpan**, dan
pesannya menyebut besaran yang melanggarnya · Perubahan yang sama pada versi `MATERIAL` -> diterima,
dan baris selisihnya muncul.

**Uji:** **Negatif:** simpan versi `TIDAK_MATERIAL` yang mengubah persentase pemulihan limit.
**Positif:** versi `MATERIAL` yang mengubahnya **diterima**, dan baris selisihnya **ada** — sebab
ketiadaan baris di sini adalah persis cacat lama, dan nol terbaca seperti "cocok".

**Menggantikan:** `TDA-18` — ***`ReinstatementPct` disalin, tidak dikurangi.*** Di
`TreatyEDMDifferenceLimits` bentuknya `ValueDifference…ReinstatementPct = .ReinstatementPct`, bukan
pengurangan terhadap `OLDDATA`. Sapuan seluruh korpus 708 berkas, dua ekspor, dua bentuk penulis:
**10 penugasan, NOL pengurangan**. Kalibrasi **lulus** — `Limit` ditemukan 6 pengurangan nyata.
Akibatnya di sistem lama: **perubahan persentase pemulihan tidak pernah menghasilkan baris
selisih**.

**Blocked by:** 11

**Dasar:**
```
EVIDENCED(TreatyEDMDifferenceLimits@ekspor-2026-09 - 10 penugasan, nol pengurangan, kalibrasi lulus atas Limit)
        DECIDED(TDA-18, INV-69)
        DIASUMSIKAN-CLEAR(DB-20)
```

- [ ] perubahan persentase pemulihan limit menghasilkan baris `NILAI_SELISIH` bertipe porsi
- [ ] versi `TIDAK_MATERIAL` yang mengubahnya **ditolak saat simpan**, pesannya menyebut besarannya
- [ ] uji positif lulus: versi `MATERIAL` diterima **dan barisnya ada**
- [ ] baris warisan **tidak** diklasifikasi ulang (`ADR-0043`)
- [ ] **pemberitahuan pra-peralihan tertulis dan terkirim**: pekerjaan yang hari ini lolos tanpa galat akan ditolak sesudah peralihan
