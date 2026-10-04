# E21: Jalur produksi endorsement ⛔ BLOCKED

**What to build:** Menuliskan selisih yang dihitung E13–E15 ke tabel produksi, beserta penomoran
versi polis.

⛔ **BLOCKED — menunggu tabel flat (P-10).** Prosedur tersimpan yang menulis dokumen polis di sistem
lama **diganti insert-ke-tabel-flat**, dan bentuknya dirancang **bersama New Business dan Renewal**,
bukan sendiri untuk endorsement.

⚠️ **Tiket ini tidak dihapus.** Lingkupnya sudah diketahui; menghilangkannya akan menyembunyikan
pekerjaan yang pasti ada.

`[terverifikasi]` Rantainya: penulisan dokumen polis → penyimpanan produksi → percabangan tujuh lini.
Ada **dua gerbang idempotensi** di lapis berbeda, sehingga rantai aman diulang — penting untuk
rekonsiliasi paralel run.

**Asal (Pega).** `Activity\SaveEDMToJsonPolicy_Act` · `SaveTreatyProduction_Act` ·
`SaveFacinProdAllEDM_Act` (tujuh cabang lini) · `InsertFacoutProductionEDM` · `GetEdmProdKe_Act`

**Keputusan.** K-046 · K-027 · K-050 · K-006 · `CLAUDE.md` §4.3

**Blocked by:** ⛔ **tabel flat (P-10)** · E13 · E15

**Status:** blocked — dilewati atas perintah work owner 01-10-2026 (tabel flat P-10).

- [ ] ⛔ Menunggu rancangan tabel flat bersama New Business dan Renewal
- [ ] **Dua gerbang idempotensi** di lapis berbeda; rantai aman diulang
- [ ] Versi polis nol-berbasis, versi baru = terakhir + 1 — **cara hitung versi terakhir aman** karena tidak ada penghapusan baris (K-050)
- [ ] Skema Oracle **tidak berubah** (§4.3)
- [ ] Nilai berkoma desimal (K-027)
- [ ] **K-046** `K046_GuardFacOut_HanyaType7` — gerbang idempotensi fac out hanya dipakai keluar pada satu jenis penyesuaian
- [ ] **K-046** `K046_JsonPolisMonitoring_HanyaPrefiksEDMT` — tabel pemantauan hanya ditulis untuk satu prefiks kasus
- [ ] **K-046** `K046_DuaSemantikPRODKE` — dua cara menentukan versi terakhir hidup berdampingan
- [ ] ⚠️ Satu cabang lini bisnis rule-nya **tidak ada di korpus** tetapi ada di folder pelengkap work owner — kandidat amandemen **K-006**, bukan `panic`
- [ ] Endpoint dari konfigurasi, tidak pernah literal (§4.4)
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
