---
status: selesai
---

# 03: Penamaan yang tidak nyaris kembar

> **SELESAI 2026-09-18.** Penggantian nama sudah dikerjakan: KEPUTUSAN_KOMITE, dan sufiks _FK tersisa nol.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.
>
> **SELESAI PENUH 19 September 2026.** Pilahan *"aturannya selesai, namanya draf"* dicabut: namanya kini berlaku juga.
>
> REQ-032 turun jadi verifikasi, dan **daftar singkatan dibatalkan** — tiket `02` mati. Nama constraint tidak lagi menunggu daftar apa pun, karena ia tidak lagi memakai singkatan.
>
> **Aturan penamaan final — empat, berlaku untuk seluruh skema:**
>
> 1. **Nama tabel dan kolom memakai kata utuh bahasa Indonesia**, `UPPER_SNAKE_CASE`, mengikuti glosarium. Tidak ada singkatan buatan. Bila sebuah nama melewati 30 byte, yang diganti **katanya** dengan kata utuh yang lebih pendek — **bukan dipotong jadi singkatan**.
> 2. **Nama constraint dan index tidak mengeja kolom.** Bentuknya `<peran>_<tabel>`, ditambah pembeda numerik bila sebuah tabel punya lebih dari satu constraint sejenis. Isinya dibaca dari katalog.
> 3. **Peran sebagai awalan tetap**: `PK_` `UQ_` `FK_` `CK_` `IX_`, ditambah `V_` untuk view dan **`SQ_` untuk sequence**. Bukan singkatan bahasa — penanda jenis objek yang baku di Oracle. *(`V_` dan `SQ_` ditambahkan 19 September 2026; `SQ_` menyusul bersama 22 sequence yang dibuat hari itu, dan daftarnya diperbarui supaya aturan ini tidak lagi tertinggal dari berkasnya sendiri.)*
> 4. **Tiap nama diperiksa panjangnya saat ditulis**, dan pembangkit **berhenti** bila ada yang melewati 30 byte.
>
> Daftar `_Avoid_` tetap berlaku penuh di nama tabel maupun kolom.
>
> **Dikerjakan 19 September 2026**: 105 objek bernama ditulis ulang di seluruh `ddl-usulan/`; satu nama tabel diganti karena constraint-nya melewati batas — `DAFTAR_KLAIM_PENJAGA_TANGGAL` → **`KLAIM_PENJAGA_TANGGAL`**, dan kolom `ID_DAFTAR` → `ID_PENJAGA_TANGGAL` supaya tidak menyebut kata yang sudah tidak ada di nama tabelnya. Pengenal di atas 30 byte: **nol**. Rinciannya `SPEC-MODEL-DATA.md` bagian 16.

*Asal: `T-14` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Tidak ada pasangan nama berjarak satu kata di skema yang belum berjalan.

`ADJUSTMENT.STATUS_AKSEPTASI` → **`KEPUTUSAN_KOMITE`**, dinamai menurut isinya: ia hasil keputusan Komite atas Adjustment, bukan keadaan entitas akseptasi. `FK_DPT_KLM_FK` → `FK_KLAIM_PENJAGA_TANGGAL_1`; sufiks `_FK` tersisa: **nol**.

**Blocked by:**

- None (can start immediately)


**Dasar:** EVIDENCED: `MEMORI_PEMAHAMAN.MD` §10.4 butir 4 mendaftar nama nyaris kembar sebagai technical debt sistem lama — `CountTotalInterest_Act` vs `CountTotalInsterest_Act`, `ProtectNilaiClaim` vs `ProteksiNilaiClaim`.

- [x] Tidak ada dua pengenal di skema yang berbeda hanya pada satu kata.
- [x] Tidak ada constraint bersufiks `_FK`.
- [x] Tidak ada singkatan buatan di nama tabel maupun kolom.
- [x] Tidak ada nama constraint yang mengeja kolomnya.
- [x] Pengenal di atas 30 byte: **nol**, diperiksa pembangkit atas 105 objek bernama.

**Ketidakpastian:** Tidak ada. REQ-032 memverifikasi, dan jawabannya tidak mengubah satu nama pun — 30 byte sah di setiap versi.
