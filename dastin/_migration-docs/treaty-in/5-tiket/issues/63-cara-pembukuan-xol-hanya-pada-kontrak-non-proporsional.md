---
status: tertahan
---

# 63: Cara pembukuan XOL dicatat, dan hanya pada kontrak non-proporsional

*Asal: `DAFTAR-PEKERJAAN.md` `P-20` · `INV-34` · `SPEC-MODEL-DATA.md` §14.4.*

**What to build:** **PK** mencatat cara pembukuan XOL pada versi, dan `INV-34` menolaknya pada kontrak
**proporsional**.

**Persyaratan:** `INV-34` · `SPEC-MODEL-DATA.md` §14.4

**Tidak termasuk:** **Perilaku pembukuannya sendiri** — di luar gelombang ini; yang dibangun hanya pencatatan dan penolakannya.

**Jalur gagal:** Mencatat cara pembukuan XOL pada kontrak proporsional -> **ditolak**, pesannya menyebut sifat proporsi kontraknya.

**Uji:** **Negatif:** kontrak proporsional. **Positif:** kontrak non-proporsional menerima seluruh nilai sahnya — dan daftar nilai sahnya **belum diketahui** sampai `Uji X-2` kembali.

**Menggantikan:** Tidak ada cacat yang digantikan — `AccountingMode` tersimpan di sistem lama dan **belum diketahui apakah dipakai**.

**Blocked by:** `14` · **`Uji X-2`**

**Dasar:**
```
DECIDED(INV-34)
        DIASUMSIKAN-CLEAR(Uji X-2)
```

- [ ] `INV-34` menolak cara pembukuan XOL pada kontrak proporsional
- [ ] **PENGHALANG:** hasil `Uji X-2` atas `AccountingMode` · **siapa menjawab:** DBA, sesudah izin kueri baca-saja turun · **yang berubah:** apakah kolomnya **dipakai sama sekali**, dan himpunan nilai sahnya
