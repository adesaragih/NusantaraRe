# R02: Layar periode renewal

**What to build:** Underwriter melihat dan menyunting **periode polis baru** berdampingan dengan
rujukan ke **polis lama** — nomor polis lama, tanggal renewal, serta tanggal berlaku dan berakhir —
dalam satu layar.

`[terverifikasi]` `Section\PeriodeRenewal` mengikat **22 properti**, di antaranya
`.QuotationData.OldPolicyNo` · `.QuotationData.RNWDate` · `.QuotationData.StatusBusiness` ·
`.PolicyData.StartDateTime` · `.PolicyData.EndDateTime` · `.PolicyData.OfferingDate` ·
`.IsSpecialAcceptance`. Varian `_IsUW` mengikat 18 properti.

⚠️ **Dua komponen terpisah, bukan satu komponen dua mode.** `PeriodeRenewal` dan
`PeriodeRenewal_IsUW` diport masing-masing sebagai komponen sendiri.

Ini **keputusan sadar work owner untuk siklus renewal**, dan ia **sengaja berbeda** dari rekomendasi
spec model-data NB §6 (satu komponen dengan prop `mode`). Alasannya reproduksi perilaku terekam
(`CLAUDE.md` §1): sistem lama memang punya dua section. **Bukan inkonsistensi** — perbedaannya
disengaja dan tercatat di sini.

**Blocked by:** R01

**Status:** needs-info — layar; lingkup RNW sebatas layanan (butir 55), menunggu keputusan membangun layar · ⏸ pengerjaan rnwfacin dihentikan work owner 01-10-2026 (`PROMPT-LANJUT-IMPLEMENT-NB-SAJA.md`)

- [ ] `PeriodeRenewal` dan `PeriodeRenewal_IsUW` diport sebagai **dua komponen terpisah**
- [ ] Nomor polis lama, tanggal renewal, dan pembeda siklus terlihat dan terikat benar
- [ ] Tanggal berlaku, berakhir, dan penawaran dapat disunting sesuai perilaku lama
- [x] Nilai awal berasal dari salinan polis lama (R01), bukan diisi ulang manual
- [ ] Komentar menyebut section Pega asal tiap komponen
- [ ] Keputusan "dua komponen terpisah" dicatat di kode sebagai pilihan sadar, dengan rujukan ke alasannya

## Comments

### 2026-10-01 — ditahan (agent)

Layar tidak dibangun (butir 55). Satu kriteria tingkat layanan sudah dipenuhi R01: nilai awal = salinan polis lama
(`BuatKasusRenewal`). `RNWDate` tetap belum ditulis (A34).
