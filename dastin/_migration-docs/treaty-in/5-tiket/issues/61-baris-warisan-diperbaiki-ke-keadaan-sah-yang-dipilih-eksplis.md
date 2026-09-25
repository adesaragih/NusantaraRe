---
status: aktif
---

# 61: Baris warisan tak terpetakan diperbaiki ke keadaan sah yang dipilih secara eksplisit

*Asal: `DAFTAR-PEKERJAAN.md` `P-52` · `ADR-0054` · `ADR-0045`.*

**What to build:** **PK** memindahkan baris ber-`WARISAN_TAK_TERPETAKAN` ke keadaan sah **yang dipilihnya
secara eksplisit** — bukan ditebak sistem — dan pilihannya tercatat beserta siapa dan kapan.
Sesudahnya baris itu berjalan seperti baris lain.

**Persyaratan:** `ADR-0054` (`PERBAIKAN_WARISAN`) · `ADR-0045` (tercatat di jejak) · `ADR-0042` (sentuh-perbaiki: penanda warisannya dicabut)

**Tidak termasuk:** **Penebakan otomatis** — dilarang `ADR-0054`. Tidak ada jalur yang memilihkan keadaannya.

**Jalur gagal:** Perbaikan tanpa memilih keadaan -> ditolak · Perbaikan yang tidak meninggalkan jejak ->
mustahil · Baris yang diperbaiki tetap membawa penanda warisan -> cacat; `ADR-0042` mencabutnya.

**Uji:** **Negatif:** perbaikan tanpa pilihan; perbaikan ke keadaan di luar delapan.
**Positif — dan ia yang membuktikan tiket ini punya kriteria sendiri:** sesudah diperbaiki, baris
itu **dapat menjalani perpindahan biasa** dan **memenuhi invarian yang berlaku sekarang**. Itu
perilaku yang tidak dimiliki tiket `60`, dan sebab keduanya tidak digabung.

**Menggantikan:** Tidak ada — sistem lama tidak punya keadaan tak-terpetakan, sebab ia tidak punya himpunan tertutup.

**Blocked by:** `60`

**Dasar:**
```
DECIDED(ADR-0054, ADR-0045, ADR-0042)
```

- [ ] keadaan dipilih **orang**, eksplisit, dari delapan yang sah
- [ ] pilihannya tercatat di jejak beserta pelaku dan waktu
- [ ] penanda warisan **dicabut** sesudah perbaikan
- [ ] uji positif: baris yang diperbaiki menjalani perpindahan biasa dan memenuhi invarian sekarang
