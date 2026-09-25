---
status: tertahan
---

# 08: Dokumen addendum membawa tanggal berlakunya sendiri, dan boleh kosong

*Asal: `DAFTAR-PEKERJAAN.md` `P-62` · `KTV-2` · `STRUKTUR-ADDENDUM.md` §3.*

**What to build:** Sebuah dokumen addendum dapat menyimpan **tanggal berlakunya sendiri**, terpisah dari
tanggal berlaku versi. Kosong **diterima** — tidak setiap dokumen membawa tanggal.

**Persyaratan:** `KTV-2` · `ADR-0035` (kegagalan tidak disamarkan menjadi nilai — kosong berarti kosong, bukan tanggal mulai kontrak)

**Tidak termasuk:** **Tanggal berlaku pada versi** — irisan 03. Keduanya ruas berbeda pada pembawa berbeda, dan
irisan ini **tidak boleh** menuliskan nilainya ke versi.
**Aturan bisnis yang memakai tanggal dokumen** — belum ada; `DB-16b` yang menentukannya.

**Jalur gagal:** Tanggal dokumen diisi lalu dibaca sebagai tanggal berlaku versi -> **cacat**; keduanya tidak boleh saling menimpa · Kosong ditolak -> salah; kosong adalah nilai sah.

**Uji:** **Negatif:** pastikan mengisi tanggal dokumen **tidak** mengubah tanggal berlaku versi mana pun.
**Positif:** dokumen tanpa tanggal **diterima** dan tetap dapat menunjuk versi.

**Menggantikan:** tidak ada — sistem lama tidak punya entitas dokumen, sehingga tidak punya tanggal untuknya.

**Blocked by:** 04

**Dasar:**
```
DECIDED(KTV-2, ADR-0035)
        DIASUMSIKAN-CLEAR(DB-16b)
```

- [ ] kolom tanggal berlaku berdiri pada `DOKUMEN_ADDENDUM`, **boleh kosong**
- [ ] uji negatif lulus: tanggal dokumen tidak menyentuh tanggal versi
- [ ] `DB-16b` tercatat di `ASUMSI-CLEAR.md` **beserta tenggatnya**: bila dibantah, kolomnya **dicabut sebelum data dimuat** — sesudah migrasi berjalan, mencabutnya berongkos
