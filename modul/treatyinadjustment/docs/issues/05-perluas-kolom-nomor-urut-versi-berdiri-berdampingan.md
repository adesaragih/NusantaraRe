---
status: aktif
---

# 05: Perluas — kolom nomor urut versi berdiri berdampingan dengan pengenal warisan

*Asal: `DAFTAR-PEKERJAAN.md` `P-65` bagian 1 · `GRL-17` · `ADR-0042`, `ADR-0043`.*

**What to build:** Kolom `NOMOR_URUT_VERSI` berdiri pada `VERSI_KONTRAK` dan **kosong**; pengenal warisan
`‹kontrak›/Rnn` tetap apa adanya. **Tidak ada pembaca yang berubah.**

Ini bagian **perluas** dari perubahan lebar: bentuk baru berdiri di samping yang lama sehingga
tidak ada yang patah.

**Persyaratan:** `GRL-17` · `ADR-0042` (sejarah pindah apa adanya) · `GRL-09` (pengenal warisan dilestarikan)

**Tidak termasuk:** **Pengisian nilainya** — irisan 10.
**Pencabutan pembacaan dari pengenal** — irisan 12.
**Pengenal `/Rnn`** tidak disentuh sama sekali; ia dilestarikan apa adanya.

**Jalur gagal:** Kolom berdiri `NOT NULL` -> setiap baris warisan gagal dimuat. Ia **wajib boleh kosong** pada tahap ini, dan itu pokok bentuk perluas.

**Uji:** **Negatif:** memuat baris warisan tanpa nomor urut **tidak boleh gagal**.
**Positif:** seluruh pembaca yang ada — pencarian, pembaca hilir bentuk lama — **mengembalikan hasil
yang sama persis** sebelum dan sesudah kolomnya berdiri.

**Menggantikan:** tidak ada — bentuk lama tidak dicabut di irisan ini; ia berdiri berdampingan.

**Blocked by:** **14** — *ditambahkan 24 September 2026.* Tiket ini menambahkan **kolom** atau **tabel anak** pada `VERSI_KONTRAK`, dan tidak satu pun tiket di papan membuat tabelnya. Irisan `14` adalah **PEMBUAT PERTAMA** `KONTRAK` dan `VERSI_KONTRAK`.

**Dasar:**
```
EVIDENCED(TreatyInRevisi_post@ekspor-2026-09 - offset @substring(ID,10,12) meleset satu)
        DECIDED(GRL-17, ADR-0042)
```

- [ ] `NOMOR_URUT_VERSI` berdiri, **boleh kosong**
- [ ] pengenal `/Rnn` tidak berubah pada satu baris pun
- [ ] seluruh pembaca yang ada mengembalikan hasil identik sebelum dan sesudah
