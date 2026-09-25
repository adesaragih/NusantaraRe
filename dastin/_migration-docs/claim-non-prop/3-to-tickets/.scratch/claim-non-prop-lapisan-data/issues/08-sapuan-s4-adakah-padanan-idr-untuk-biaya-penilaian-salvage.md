---
status: selesai
---

# 08: Sapuan S4: adakah padanan IDR untuk biaya penilaian, salvage, dan biaya lain

> **SELESAI 2026-09-18.** Dijalankan 18 September 2026 atas empat lapisan; hasilnya `SAPUAN-S3-S9.md` bagian S4.
> Dipindahkan dari papan aktif, tidak dihapus: tiket yang pernah ada adalah bukti bahwa pekerjaannya tidak dilewatkan.

*Asal: `T-35` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Jawaban apakah ketiganya memang tanpa pasangan IDR, atau kolomnya yang kurang.

Sapuan seluruh properti bersufiks `IDR` pada class `SpreadingRisk`, `Data-Adjustment`, dan `ClaimData`, dibandingkan terhadap daftar properti bernilai uang. Empat lapisan.

**Blocked by:**

- None (can start immediately)


**Dasar:** DERIVED dari review; DECIDED(ADR-0007) mewajibkan pasangan untuk setiap nilai uang.

- [ ] Daftar properti uang yang punya padanan IDR dan yang tidak, beserta berkas dan baris.

**Ketidakpastian:** **Menahan T-05.** Bila ketiganya perlu pasangan IDR, kolomnya bertambah di `ALOKASI_LAYER` dan `CHECK`-nya ikut terpasang sendiri lewat pembangkit.

## Hasil sapuan

- **Ya, ada kelas nilai uang yang hidup tanpa padanan IDR sama sekali** — tujuh properti, lima di antaranya masuk hitungan instruksi bayar.
- Penamaan alternatif (`Rupiah`, `Idr`, `_IDR`, `Rp`, `InRp`) ikut disapu: **nol hasil**. Yang ada hanya sufiks `RNM`, yaitu porsi pihak, bukan mata uang.
- **Menyentuh usulan ADR-0029 dan tidak saya selesaikan sendiri.** Dua pembacaan berdiri di atas bukti yang sama; membedakannya butuh jawaban orang. Diajukan sebagai **A16** ke akuntansi.
