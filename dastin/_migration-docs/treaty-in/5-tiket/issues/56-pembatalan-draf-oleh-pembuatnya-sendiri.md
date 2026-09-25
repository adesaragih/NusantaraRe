---
status: aktif
---

# 56: Pengisi kontrak membatalkan draf yang ia buat sendiri, dan kontraknya langsung terbuka untuk versi berikutnya

*Asal: `DAFTAR-PEKERJAAN.md` `P-57` · `ADR-0055` perubahan 24 Sep 2026 · `INV-04`.*

**What to build:** **PK** membatalkan versi `DRAFT` yang **ia buat sendiri**. Versinya menjadi `DIBATALKAN` —
terminal — **barisnya tetap tersimpan**, nomornya tidak dipakai ulang, dan kontraknya **langsung
terbuka** untuk versi berikutnya.

**Persyaratan:** `ADR-0055` (`DRAFT → DIBATALKAN`) · `INV-04` · `INV-25`

**Tidak termasuk:** **Pembatalan oleh orang lain** — tidak ada, dan itu keputusan: hanya pembuatnya.

**Jalur gagal:** Orang lain membatalkan draf -> ditolak · Membatalkan versi yang sudah diajukan -> ditolak;
jalurnya `KEMBALIKAN` lalu batal · Nomor versi yang dibatalkan dipakai ulang -> ditolak.

**Uji:** **Negatif:** pelaku bukan pembuat; versi bukan `DRAFT`; pakai ulang nomornya.
**Positif:** sesudah pembatalan, versi berikutnya **langsung diterima** — itu yang membuktikan
`INV-25` membebaskan slotnya.

**Menggantikan:** **Tidak ada padanannya — ini kemampuan BARU.** Sistem lama **tidak punya cara membuang
draf** yang tidak merusak apa pun: satu-satunya jalan adalah mengajukannya supaya ditolak, dan
penolakan **menghapus barisnya** (`TDA-02`).

**Blocked by:** `45` · `54`

**Dasar:**
```
EVIDENCED(ADR-0055 perubahan 24 Sep 2026 - L-7 ditutup)
        DECIDED(ADR-0055, INV-04, INV-25)
```

- [ ] hanya pembuatnya dapat membatalkan, dan hanya saat `DRAFT`
- [ ] baris tetap tersimpan, nomor tidak dipakai ulang
- [ ] uji positif: versi berikutnya langsung diterima sesudahnya
- [ ] **CARA MENYALAKANNYA** ditulis — ini golongan BARU
