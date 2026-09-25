---
status: aktif
---

# 54: Tiap perpindahan keadaan meninggalkan satu catatan persetujuan, dan tidak ada jalur yang melewatinya

*Asal: `DAFTAR-PEKERJAAN.md` `P-46` · `ADR-0045` · `SPEC-MODEL-DATA.md` §10.20.*

**What to build:** Setiap perpindahan keadaan menulis **satu** baris `CATATAN_PERSETUJUAN` berisi pelaku,
tingkat, waktu, dan alasannya bila ada. **Tidak ada jalur yang melewatinya** — ia fakta mesin.

**PEMBUAT PERTAMA** untuk `CATATAN_PERSETUJUAN`.

**Persyaratan:** `ADR-0045` (jejak sebagai fakta mesin, tidak dapat dilewati) · `SPEC-MODEL-DATA.md` §10.20 · `INV-26`

**Tidak termasuk:** **`JEJAK_PERUBAHAN`** — entitas berbeda, sudah dibangun tiket `39`. Yang satu mencatat
**perpindahan keadaan**, yang lain mencatat **perubahan fakta**.
**`PERISTIWA_KONTRAK`** — lihat `PENGHALANG` di bawah; ia **tidak punya kemampuan** dan tidak
dibangun di sini.

**Jalur gagal:** Perpindahan yang berhasil tanpa catatan -> **mustahil**; bila mungkin, tiket ini belum
selesai · Catatan tanpa pelaku -> ditolak.

**Uji:** **Negatif:** coba lakukan tiap perpindahan lewat jalur yang melewati penulisan catatan.
**Positif:** **ketiga belas** perpindahan diuji, dan **ketiga belasnya** menghasilkan tepat satu
catatan — bukan nol, bukan dua.

**Menggantikan:** **`CommentList` sistem lama** hanya punya `OperatorName`, `IsApproved`, `Suggest`, dan
`Date` — **tidak merekam tingkat penyetuju sama sekali**. Dan zona waktunya campuran: komentar
persetujuan memakai jam Pega (GMT), *"Create Revision"* memakai jam Oracle, dan langkah jam Oracle
di `AddCommentList_Act` langkah 1 **mati**.

**Blocked by:** `45`

**Dasar:**
```
EVIDENCED(AddCommentList_Act@ekspor-2026-09 - langkah 1 MATI, CommentList tanpa ruas tingkat)
        DECIDED(ADR-0045, INV-26)
```

- [ ] `CATATAN_PERSETUJUAN` berdiri sesuai `KAMUS-KOLOM.md`, **termasuk ruas tingkat** yang sistem lama tidak punya
- [ ] ketiga belas perpindahan menghasilkan tepat **satu** catatan masing-masing
- [ ] **tidak ada jalur** yang dapat berpindah tanpa menulis catatan — diuji, bukan diargumentasikan
- [ ] waktu disimpan dalam **satu** zona yang dinyatakan, bukan campuran
