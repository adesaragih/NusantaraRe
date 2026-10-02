# E02: Tiga predikat yang perilakunya berbeda di endorsement

**What to build:** Tiga predikat bernama sama dengan New Business tetapi **berperilaku berbeda** di
endorsement, ditambah satu penegasan yang mudah salah diimplementasikan.

`[terverifikasi]` Di endorsement, operator ber-workbasket **marketing dihitung sebagai underwriter** —
menentukan **siapa yang boleh menyetujui**. Predikat klaim menguji **keberadaan halaman data**, bukan
prefiks ID kasus. Predikat travel menerima **dua sumber**, bukan satu.

⛔ **Predikat "bukan endorsement" BUKAN negasi predikat "endorsement".** Keduanya membaca **jalur
properti berbeda** — satu agregat tersimpan, satu halaman aktif. Keduanya dapat bernilai benar
bersamaan, atau salah bersamaan.

**Asal (Pega).** `When\IsUW` · `When\IsClaim` · `When\IsTravel` · `When\IsEDM` · `When\IsNotEDM`

**Keputusan.** K-046 · K-050 · `CLAUDE.md` §4.5

**Blocked by:** E01

**Status:** sebagian — 01-10-2026. Lewat registry EDM: IsUW 4 workbasket (+Marketing), IsTravel dua
sumber, IsNotEDM implementasi sendiri (dua arah tidak sinkron diuji), ketiga K-046
(`predikat_test.go`). IsClaim (`PropertyHasValue(pyWorkPage.ClaimData)`, panic di generator) diport
tangan sebagai `sikapKhusus`; arti "punya nilai" untuk halaman `[belum terverifikasi]` → dijawab
pemanggil (`PemeriksaNilai`). Siapa yang boleh menyetujui: tetap `[pertanyaan terbuka]` Underwriting.

- [ ] Predikat underwriter endorsement mengakui **satu workbasket lebih banyak** daripada New Business
- [ ] Predikat klaim endorsement menguji keberadaan halaman data; New Business menguji prefiks ID
- [ ] Predikat travel endorsement menerima dua sumber; New Business satu
- [ ] Predikat "bukan endorsement" punya **implementasi sendiri** — membaca halaman aktif
- [ ] Kasus uji: kedua halaman sengaja **tidak sinkron** → membuktikan keduanya dapat bernilai sama
- [ ] **K-046** `K046_IsUW_MarketingSebagaiUnderwriter_EDM`
- [ ] **K-046** `K046_IsClaim_UjiPrefiks_vs_UjiHalaman`
- [ ] **K-046** `K046_IsTravel_LabelTidakDieksekusi` — label tampilan menyebut kode bisnis yang tidak dieksekusi

⚠️ Siapa yang boleh menyetujui di endorsement masih **`[pertanyaan terbuka]`** milik Underwriting.
Sampai dijawab, perilakunya **diport apa adanya**.
