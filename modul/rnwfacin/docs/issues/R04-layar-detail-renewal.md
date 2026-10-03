# R04: Layar detail renewal

**What to build:** Layar detail renewal menampilkan seluruh rincian penawaran — objek, coverage,
spreading, komisi, pembayaran — dengan **menyisipkan komponen New Business**, sama seperti R03, dan
menambahkan hanya field yang memang khas renewal.

Ini layar terbesar dalam lingkup renewal.

`[terverifikasi]` `Section\InputRenewalDtl` menyisipkan **14 section** dan mengikat **37 properti
sendiri**; varian `_IsUW` menyisipkan **17 section** dan mengikat 36 properti. Yang disisipkan antara
lain `CoverageList` · `CoverageSpreadingList` · `CoverageCommisionList` · `PaymentCurrencyList` ·
`OfferFacIn_NusaRe` · `SummaryCoverage_Section` · `SummarySpreading_Section` · `CallDualScoringRisk` ·
`InputDtlObject_FacIn`.

⚠️ **Tiket ini sengaja dibiarkan utuh, tidak dipecah.** Pengerja **boleh** memecah varian `_IsUW`
menjadi sub-tugas bila satu context window tidak cukup — tetapi pemecahan itu keputusan pelaksanaan,
bukan pemecahan tiket.

⚠️ **Dua komponen terpisah** — `InputRenewalDtl` dan `InputRenewalDtl_IsUW` masing-masing komponen
sendiri. Keputusan sadar work owner, sengaja berbeda dari spec model-data NB §6; alasannya reproduksi
perilaku terekam (`CLAUDE.md` §1).

### ⏸ `InputDtlObject` — usang, tetapi diport apa adanya

`[terverifikasi]` `Section\InputRenewalDtl` menyisipkan `InputDtlObject` **empat kali** (L2764 · L2944 ·
L3207 · L3390), sementara section itu **tidak ada di NB maupun RNW**. Penggantinya
`InputDtlObject_FacIn` **ada di keduanya** dan muncul **26 kali** di berkas yang sama.

Keputusan work owner (K-041): **`InputDtlObject` usang, digantikan `InputDtlObject_FacIn`**; keempat
sisipan adalah sisa kelewat.

⛔ **Jangan hapus keempat sisipan itu diam-diam** (`CLAUDE.md` §1). Karena section tujuannya tidak ada,
sisipan itu tidak menghasilkan apa pun dan menjadi tidak berdampak dengan sendirinya.

⚠️ **`InputDtlObject` ≠ `InputDtlObject_FacIn`** — dua section berbeda, dibedakan hanya oleh sufiks.
Jangan disatukan, jangan diperlakukan sebagai salah ketik.

**Blocked by:** R03

**Status:** needs-info — layar; lingkup RNW sebatas layanan (butir 55), menunggu keputusan membangun layar · ⏸ pengerjaan rnwfacin dihentikan work owner 01-10-2026 (`PROMPT-LANJUT-IMPLEMENT-NB-SAJA.md`)

- [ ] `InputRenewalDtl` dan `InputRenewalDtl_IsUW` diport sebagai **dua komponen terpisah**
- [ ] Keduanya **menyisipkan** komponen NB; tidak ada komponen NB yang disalin
- [ ] Ke-37 dan 36 properti khas layar ini terikat benar
- [ ] Keempat sisipan `InputDtlObject` **diport apa adanya**, dengan komentar menyatakan statusnya usang dan section penggantinya
- [ ] Tidak ada penyatuan `InputDtlObject` dengan `InputDtlObject_FacIn`
- [ ] Tidak ada perhitungan premi, spreading, atau komisi yang ditulis ulang di tiket ini — seluruhnya diwarisi NB
