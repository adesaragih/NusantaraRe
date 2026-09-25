---
status: aktif
---

# 53: Section head, dept head, atau direktur mengembalikan versi ke draft, dan pengembalian tanpa alasan ditolak

*Asal: `DAFTAR-PEKERJAAN.md` `P-34` · `ADR-0055` (`KEMBALIKAN`).*

**What to build:** **SH / DH / DR** mengembalikan versi yang menunggunya ke `DRAFT`. Pengembalian **tanpa
alasan ditolak**, dan alasannya tersimpan di catatan persetujuan.

Perpindahan ini yang membuat pembekuan jenis dan materialitas (tiket `07`) **tidak menjadi
penguncian permanen**.

**Persyaratan:** `ADR-0055` (`KEMBALIKAN` ×3) · `INV-27`

**Tidak termasuk:** **Penolakan** — tiket `55`. Mengembalikan dan menolak adalah dua perpindahan berbeda dengan akibat berbeda.

**Jalur gagal:** Pengembalian tanpa alasan -> ditolak · Pengembalian oleh orang yang bukan pemegang antriannya -> ditolak karena peran.

**Uji:** **Negatif:** alasan kosong; alasan hanya spasi; pelaku salah.
**Positif:** sesudah dikembalikan, versinya **dapat disunting lagi** termasuk jenis dan
materialitasnya — itu yang membuktikan tiket `07` tidak mengunci permanen.

**Menggantikan:** `TDA-02` dan bentuk umumnya: sistem lama **tidak punya perpindahan mengembalikan** untuk
addendum sama sekali. Satu-satunya jalan keluar dari pengajuan adalah **ditolak**, dan penolakan
**menghapus barisnya**.

**Blocked by:** `49`

**Dasar:**
```
EVIDENCED(TreatyInDeclineConfirmation_postactEDM@ekspor-2026-09 - langkah 1-4 MATI, dua RDB remove hidup)
        DECIDED(ADR-0055)
```

- [ ] ketiga tingkat dapat mengembalikan; alasan **wajib**
- [ ] alasannya tersimpan di catatan persetujuan
- [ ] uji positif: sesudah kembali ke `DRAFT`, jenis dan materialitas dapat diubah lagi
