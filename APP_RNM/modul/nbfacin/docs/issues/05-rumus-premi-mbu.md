# 05: Rumus premi MBU — bentuk bersarang

**What to build:** Premi MBU dihitung benar, dan bersamanya **aturan pembagi komposit terbukti dari
struktur ekspresinya sendiri** — bukan dari tafsir kita.

MBU memakai bentuk **bersarang**, dan di situlah nilainya. Pembagi **dalam** menempel pada rate;
pembagi **luar** menempel pada pro-rata. Keterikatan faktor→pembagi karena itu **tertulis eksplisit**,
tidak perlu disimpulkan:

```
@Math.divide( @Math.divide((.TSI*(.Rate+.Loading)),100,4) * @if(.ProRatePercent=="",100,.ProRatePercent), 100, 4)
                                          ^^^                                                            ^^^
                              pembagi DALAM → .Rate (%)                        pembagi LUAR → ProRatePercent
```

Asal: `FillPremiMBU_FacIn` **L1144** dan **L1626**.

Dua hal yang mudah terlewat: `.Loading` **dijumlahkan ke rate sebelum dikalikan** (bukan sesudah),
dan pro-rata kosong diperlakukan sebagai **100**, bukan nol — membalik keduanya mengubah angka tanpa
gejala.

**Blocked by:** 03

**Status:** ready-for-human — diport 01-10-2026, menunggu tinjauan work owner

- [ ] Kedua bentuk (`L1144`, `L1626`) diimplementasikan sebagai cabang internal di balik pintu masuk yang sama
- [ ] Struktur **bersarang** dipertahankan — bukan diratakan menjadi satu pembagian dengan pembagi 10.000
- [ ] `.Loading` dijumlahkan ke rate **sebelum** perkalian, sesuai sumbernya
- [ ] Pro-rata kosong → **100**, direproduksi apa adanya
- [ ] Skala MBU **persen** diambil dari resolver (tiket 03), bukan ditulis ulang di sini
- [ ] **Rekonsiliasi eksak** per bentuk, tanpa toleransi
- [ ] Kasus uji yang membuktikan pembagi dalam terikat ke rate dan pembagi luar ke pro-rata — bila keduanya tertukar, uji harus **gagal**

## Comments

### 2026-10-01 — implementasi (agent)

Kode: `premium/mbu.go`. ✅ **Rekonsiliasi: 93 dari 93 baris coverage dari lima kasus MBU nyata cocok
eksak** (Python decimal; empat baris jadi fixture `mbu_test.go`, daftar-izin angka saja).

⚠️ **Temuan yang mengubah tiket:** "L1626" BUKAN rumus premi. `[terverifikasi]` Baris itu ada di
`pyStepsPreCondParamsWhen` langkah 2.2.1.5 `FillPremiMBU_FacIn` — syarat `.MinPremium = 100000`. Satu-
satunya langkah yang menulis `.Premium` adalah 2.2.1.3 (L1144). Kriteria "kedua bentuk" karenanya
tidak berlaku; `.MinPremium` dan diskon MBU (2.2.1.7–2.2.1.8) tidak menulis `.Premium` dan tidak diport.

| Kriteria | Keadaan |
| --- | --- |
| Struktur bersarang dipertahankan | ✅ dibuktikan contoh hitung yang membedakan dari ÷10.000 (0,5638 vs 0,5637) |
| `.Loading` dijumlahkan sebelum perkalian | ✅ contoh hitung; Loading kosong = 0 `[terverifikasi]` 93 baris |
| Pro-rata kosong → 100 | ✅ contoh hitung (data nyata selalu terisi) |
| Skala dari resolver | ✅ |
| Rekonsiliasi eksak | ✅ jalur sederhana; ⚠️ tidak ada baris nyata yang memicu pembulatan atau pro-rata ≠ 100 (`../PERTANYAAN-LANJUTAN.md` butir 3) |

`.ProRatePercent` di L1144 adalah medan COVERAGE (`ProRatePercentCoverage`), bukan
`pyWorkPage.OfferFacIn.ProRatePercent` milik PA. Uji mutasi: 4/4.

### 2026-10-01 — koreksi status satu kriteria (agent, tindak lanjut review)

⚠️ Kriteria *"bila pembagi dalam dan luar tertukar, uji harus gagal"* **tidak dapat dibuktikan seperti
tertulis**: kedua pembagi sama-sama 100, jadi menukarnya tidak mengubah angka. Yang dibuktikan uji adalah
bentuk bersarang vs diratakan ÷10.000 (0,5638 vs 0,5637). Tanda ✅ pada baris "Struktur bersarang"
berlaku untuk itu saja. Keputusan agent A4 dan A12 menunggu konfirmasi work owner.
