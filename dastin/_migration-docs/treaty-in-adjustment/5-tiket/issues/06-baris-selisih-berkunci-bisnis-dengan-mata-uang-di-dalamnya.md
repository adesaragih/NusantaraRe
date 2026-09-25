---
status: aktif
---

# 06: Baris selisih berkunci bisnis, dan mata uang ada di dalam kuncinya

*Asal: `DAFTAR-PEKERJAAN.md` `P-63` · `ADR-0048` butir 3 · `GRL-14` · `GRL-16` · bahan to-spec `E-1a`.*

**What to build:** Menyimpan versi penyesuaian menghasilkan **baris selisih per besaran yang berubah**,
dipadankan terhadap versi dasarnya lewat **kunci bisnis** — bukan menurut posisi baris. Mata uang
**ada di dalam kunci**, sehingga mengubah mata uang tampil sebagai baris dihapus ditambah baris
baru, bukan sebagai selisih angka.

Artefak: entitas `NILAI_SELISIH`, kunci padanannya, dan jalur yang mengisinya saat simpan.

**PEMBUAT PERTAMA** untuk `NILAI_SELISIH`.

**Persyaratan:** `ADR-0048` butir 3 · `ADR-0053` (mata uang adalah daftar) · bahan to-spec `E-1a` · `GRL-14` · `GRL-16`

**Tidak termasuk:** **`INV-69` dan `INV-70`** — irisan 11.
**Pohon `ActualValue`** — `GRL-14` menghapusnya; tidak ada pohon ketiga yang dibangun di sini.
**Perhitungan ulang selisih versi yang sudah disetujui** — `ADR-0036`: angka dasar persetujuan beku.

**Jalur gagal:** Baris disisipkan di **tengah** daftar pada versi baru -> baris sesudahnya **tetap dipadankan
dengan pasangan bisnisnya**, bukan bergeser · Daftar baru lebih panjang daripada daftar lama ->
barisnya muncul sebagai **baris baru**, bukan dikurangi baris yang tidak ada · Mata uang berubah ->
**dua baris**, bukan satu selisih.

**Uji:** **Negatif:** sisipkan baris di tengah dan pastikan padanan **tidak** bergeser; ubah mata uang
dan pastikan hasilnya bukan pengurangan lintas mata uang.
**Positif — dan ia WAJIB, sebab `GRL-16` bersandar padanya:** mengubah **share fakultatif**
menghasilkan sedikitnya **satu** baris selisih. Uji yang tidak pernah gagal atas besaran itu berarti
besaran itu **tidak ada di dalam himpunan yang dipadankan** — dan itu persis cacat yang `GRL-16`
temukan di sistem lama.

**Menggantikan:** `TDA-04` — *selisih dipadankan menurut posisi baris.* Dan `TDA-05` — *selisih dikurangkan
tanpa memeriksa mata uang.* Sistem lama memadankan lewat `pxListSubscript` dan `<CURRENT>`; bila
sebuah baris disisipkan di tengah, seluruh baris sesudahnya dipadankan dengan baris lama yang salah
**tanpa satu galat pun**. Ditambah `TDA-06` — sebagian selisih ditulis ke pohon yang salah lalu
ditimpa.

**Blocked by:** 01

**Dasar:**
```
EVIDENCED(TreatyEDMCalculateDifference@ekspor-2026-09, TreatyEDMDifferenceShare@ekspor-2026-09, TreatyEDMDifferenceDeduction@ekspor-2026-09)
        DECIDED(ADR-0048, ADR-0053, GRL-14, GRL-16)
```

- [ ] `NILAI_SELISIH` berdiri dengan kunci padanannya, dan **mata uang ada di dalam kunci**
- [ ] penyisipan baris di tengah daftar **tidak** menggeser padanan
- [ ] perubahan mata uang menghasilkan dua baris, bukan satu selisih
- [ ] **uji positif share fakultatif lulus** — perubahannya menghasilkan baris selisih
- [ ] selisih versi yang sudah disetujui **tidak berubah** setelah rumusnya diperbaiki
