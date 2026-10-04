# E13: Selisih pembayaran tingkat mata uang

**What to build:** Untuk setiap mata uang pada kasus endorsement, hitung **selisih pembayaran** —
nilai sesudah dikurangi nilai sebelum.

`[terverifikasi]` Rumus kanoniknya tunggal dan berulang di belasan penugasan. Nilai "sebelum"
diambil dari **lapis A** (dokumen polis lama), bukan dari nilai lama per baris.

⛔ **Netto menambahkan PPh dan PPN, tidak mengurangkannya.** Terlihat berlawanan intuisi pajak;
diport apa adanya.

⚠️ **Urutan mengikat:** transformasi perhitungan premi **menimpa** porsi periode yang disetel modul
before-image sebelum jalur produksi membacanya (K-048).

**Asal (Pega).** `Activity\CountEndorsementData` · `CountDataEDMElse` · `CountPaymentEdm_Act` ·
`CountPaymentEdmTSIObj_Act` · `DataTransform\CountPremiEDM_DT`

**Keputusan.** K-048 · K-046 · K-010/K-012

**Blocked by:** E12 · E06

**Status:** sebagian — 01-10-2026. Porsi periode tahap pembayaran:
`backend/services/porsipembayaran.go` — `CountPaymentEdm_Act` langkah 1 (keluar bila bukan EDM), 3, 4
(normalisasi tanggal polis ke tanggal Jakarta 05:00 GMT), 9 (FIX RATE untuk MarineCargo ∨
AdjShareCedant), 11 (`CountPaymentEdmTSIObj_Act`), blok 13–15; gerbang dari registry
(`PredikatPembayaranDari`). Rekonsiliasi kasus nyata `ProrateEDMEnd = 1` (lemah: konstanta blok 13).
Belum: rumus pembayaran 13.2/14.4/15.5, `SumTotalPayment`, netto +PPh+PPN, `CountPremiEDM_DT`.

- [ ] Selisih pembayaran = nilai sesudah − nilai sebelum, per mata uang
- [ ] Netto = premi − komisi − brokerage − potongan **+ PPh + PPN**
- [ ] Nilai sebelum bersumber **lapis A**, bukan lapis B
- [ ] Lima komponen "lama" tingkat pembayaran ikut terisi: komisi · PPN · PPh · brokerage · potongan
- [ ] Mata uang yang **tidak ada** di polis lama → seluruh nilai dianggap **baru**, nilai sebelum dipaksa nol
- [ ] Seluruh nilai bertipe `Money` dengan mata uangnya
- [ ] ⚠️ Urutan eksekusi terhadap modul before-image **dijaga** (K-048)
- [ ] **K-046** `K046_NettoPremi_PPhPPNDitambahkan`
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
