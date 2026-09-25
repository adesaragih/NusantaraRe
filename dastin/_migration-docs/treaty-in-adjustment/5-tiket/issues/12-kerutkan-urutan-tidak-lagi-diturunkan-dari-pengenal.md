---
status: aktif
---

# 12: Kerutkan — tidak ada lagi pembaca yang menurunkan urutan dari pengenal warisan

*Asal: `DAFTAR-PEKERJAAN.md` `P-65` bagian 3 · `GRL-17` · `GRL-11`.*

**What to build:** Turunan *"versi berlaku"* dan setiap pengurutan versi dibaca **hanya** dari
`NOMOR_URUT_VERSI`. Tidak ada lagi kode di mana pun yang mengurai angka dari pengenal `/Rnn`.

Pengenal `/Rnn` **tetap disimpan dan tetap ditampilkan** — yang dicabut adalah **pembacaannya
sebagai urutan**, bukan keberadaannya.

**Persyaratan:** `GRL-17` · `GRL-11` · `GRL-09` (pengenal warisan dilestarikan apa adanya)

**Tidak termasuk:** **Penghapusan pengenal `/Rnn`** — dilarang. `GRL-09` melestarikannya; nomor itu sudah
beredar di luar sistem.

**Jalur gagal:** Sebuah pembaca yang masih mengurai `/Rnn` tertinggal -> ia akan **benar untuk hampir semua kontrak** dan salah untuk yang nomornya pernah dipakai ulang — kegagalan yang paling sulit ditemukan.

**Uji:** **Negatif:** sapuan seluruh kode untuk pengurai `/Rnn` mengembalikan **nol**, dan sapuannya
**dikalibrasi** terhadap satu pemakaian yang diketahui ada sebelum irisan ini.
**Positif:** turunan "versi berlaku" atas kontrak yang nomor `/Rnn`-nya **pernah bentrok**
mengembalikan versi yang benar.

**Menggantikan:** `TDA-12` — *offset pengurai nomor revisi meleset satu; patah pada revisi kesepuluh.*
Deret yang dihasilkannya `R09 -> R010 -> R11 -> R02`, bergantung pada perilaku `@substring` di luar
panjang teks — yang **tidak terbaca dari ekspor**.

**Blocked by:** 10

**Dasar:**
```
EVIDENCED(TreatyInRevisi_post@ekspor-2026-09 - @substring(ID,10,12))
        DECIDED(GRL-17, GRL-11, GRL-09)
```

- [ ] sapuan pengurai `/Rnn` mengembalikan **nol**, dan sapuannya **dikalibrasi**
- [ ] pengenal `/Rnn` tetap tersimpan dan tetap ditampilkan
- [ ] uji positif lulus atas kontrak yang nomornya pernah bentrok
