---
status: aktif
---

# 09: Nomor dokumen baris warisan diisi dari arsip kertas, dan yang belum terisi terlihat sebagai daftar

*Asal: `DAFTAR-PEKERJAAN.md` `P-66` · `GRL-19` §3.1 · `STRUKTUR-ADDENDUM.md` §3.1.*

**What to build:** **PK** dapat mengisi nomor dokumen pada baris addendum warisan, satu per satu, dari arsip
kertas. Baris yang **belum** terisi muncul sebagai **daftar yang dapat dibaca** — bukan sebagai diam.

Ini **pekerjaan orang**, bukan kemampuan sistem yang berjalan sendiri. Migrasi **tidak mengisinya**:
ia tidak punya sumber.

**Persyaratan:** `GRL-19` §3.1 · `ADR-0042` (baris warisan ditandai, diterima apa adanya)

**Tidak termasuk:** **Pengisian otomatis oleh migrasi** — mustahil, dan itu temuan bukan kekurangan: sapuan
ekspor mengembalikan **nol** properti penyimpan nomor dokumen.
**Menolak baris warisan tanpa nomor** — dilarang; `ADR-0042` menerima warisan apa adanya.

**Jalur gagal:** Baris warisan tanpa nomor dokumen ditolak oleh invarian mana pun -> **cacat**; ia harus diterima dan terdaftar · Daftar yang belum terisi tidak dapat dibaca siapa pun -> pekerjaan ini tidak akan pernah dijadwalkan.

**Uji:** **Negatif:** invarian keunikan nomor dokumen **tidak boleh** menolak beberapa baris warisan
yang sama-sama kosong.
**Positif:** daftar baris tanpa nomor dapat dibaca, dan berkurang satu setiap kali sebuah nomor
diisi.

**Menggantikan:** tidak ada — sistem lama tidak pernah merekam nomor dokumen, sehingga tidak ada perilaku yang digantikan. **Ini kemampuan BARU yang menutup lubang yang dibuat oleh ketiadaan sumber.**

**Blocked by:** 04

**Dasar:**
```
EVIDENCED(sapuan-properti-dokumen@ekspor-2026-09 - NOL hasil, dikalibrasi)
        DECIDED(GRL-19, ADR-0042)
        DIASUMSIKAN-CLEAR(DB-16a)
```

- [ ] baris warisan dapat menyimpan nomor dokumen yang diketik orang
- [ ] beberapa baris warisan **tanpa** nomor tidak saling bertabrakan pada invarian keunikan
- [ ] daftar baris yang belum terisi **dapat dibaca**, dan berkurang saat diisi
- [ ] ongkos pengisiannya **diukur** dan dilaporkan — berapa baris, dan apakah sepadan
