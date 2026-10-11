# 68: Nilai awal coverage (medan kosong tidak tersimpan NULL)

> ⚠️ **Disusun agent sesi c3 - bukan hasil `/to-tickets`.** Kontrak dari sesi 9d (keputusan work owner 08-10-2026).

**Status:** backend selesai sesi c3 08-10-2026; frontend (sesi 0f) `NILAI_AWAL_COVERAGE` di `FormCoverage.tsx`.

## Keputusan work owner 08-10-2026

- Kutipan (dengan tangkapan layar form coverage Pega basis Sum Insured): *"ini nilai default saat isi coverage, dan jika tidak di
  isi, dikasih nilainya, jangan null"*.
- Tambahan untuk medan khusus basis: *"perlu juga, ikutin seperti xml dari pega"*.

## Dasar - BUKAN dari korpus (kecuali medan basis)

`[terverifikasi]` telusur sesi 9d `D:\migrasi\RNM\NB FacIn\`:
- `CountPremi_ACT` L3482 / L3627 "Default LossLimit = 100" hanya mengisi `Local.LossLimit` (variabel hitung lokal), bukan properti
  `.LostLimit`;
- `AddCoverageAutoFire` hanya mengisi CoverageBasis 1 dan IndemnityPercentage 100;
- rule default milik `Data-Coverage` tidak ada di korpus.

Nilai INDEMNITY / UNIT / RATE / DISCOUNT_PERCENTAGE / NET_RATE / LIMITOF_LIABILITY / PCT_LO_L / LOST_LIMIT / DISCOUNT / PREMIUM =
**keputusan work owner**, bukan Pega. Medan basis (telusur 9d): SUBLIMIT 100 = `CountPremi_ACT` L3638 `Local.SubLimit = 100` bila
tidak > 0 (L3717); FIRST_LOSS / EML_PML tanpa default di XML, rumus `TSI*x/100` memperlakukan kosong = 0; FIRST_SCALE tidak
(dihitung dari tabel First Loss Scale, SetFirstLossScale_ACT L1073).

## Nilai (kosong / absen -> nilai; isian pengguna, termasuk 0, tidak diubah)

| Kolom | Nilai | Kapan |
| --- | --- | --- |
| INDEMNITY (teks) | "0" | sesudah hitung |
| UNIT (Day (s)) | "0" | sesudah hitung |
| RATE, DISCOUNT_PERCENTAGE, NET_RATE, LIMITOF_LIABILITY, PCT_LO_L, DISCOUNT | 0 | sesudah hitung |
| LOST_LIMIT | 100 | sesudah hitung |
| PREMIUM | 0 | sesudah hitung (CountPremi) |
| SUBLIMIT | 100 | SEBELUM hitung |
| FIRST_LOSS, EML_PML | 0 | SEBELUM hitung |

TIDAK diubah: FIRST_SCALE, RATE_OJK, medan server (TSI, TSI_LIABILITY, PRO_RATE_PERCENT).

## Backend sesi c3

- **K68-1 Satu titik** `services/coverage.go` `siapkanCoverageDengan` - dipakai SEMUA penulisan T_COVERAGELIST: PUT …/objek (Save tab
  Coverage; coverage otomatis dan Copy coverage adalah aksi layar yang disimpan lewat PUT ini - `GET …/coverage-otomatis` hanya daftar
  pilihan) dan Save Data CSV objek (tiket 66, `gantiObjek`). PUT …/akumulasi (tahap UW) tidak lewat sini - tetap hanya kolom akumulasi.
- **K68-2 `nilaiAwalCoverage` SESUDAH hitung** `[keputusan agent]`: `hitungSimpan` menurunkan mode dari kosongnya Rate / Premium (A158)
  - Rate 0 sebelum hitung akan mengubah premi ketikan pengguna (mode amount) menjadi 0. Sesudah hitung Rate / Premium kosong hanya bila
  KEDUANYA kosong (premi 0 = rate 0 × …, cocok "rate 0 memberi premi 0"). Simpan ulang dengan nilai awal memberi hasil SAMA:
  LossLimit kosong / 0 = 100 (L3482), NetRate 0 = kosong (`adaNet`), DiscountPercentage 0 -> Discount 0, LimitOfLiability / PctLoL /
  INDEMNITY teks / UNIT tidak dipakai hitung.
- **K68-3 `nilaiAwalBasis` SEBELUM hitung** (permintaan 9d; TSILiability ikut nilai awal, simpan ulang konsisten), pada SALINAN larik
  coverage. Premi tidak berubah. ⚠️ **Tepi**: basis 4 dengan PctAdjustment / NetRate memakai `.Sublimit` MENTAH (CountPremi langkah
  45-50) - Sublimit kosong dulu TSILiability 0, kini TSI × 100 / 100 (= nilai Pega bila field berisi 100).
- Data lama NULL terisi saat disimpan ulang; GET tetap mengembalikan "" untuk NULL; tanpa backfill / migrasi.

**Uji** `TestNilaiAwalCoverage`, `TestNilaiAwalCoverageSimpanUlang` (percent, amount, kosong, diskon persen / nominal, Loss Limit 50,
basis 2 / 3 / 4, adjustable: premi = hitungSimpan isian asli, simpan ulang premi / rate / diskon / net rate / loss limit /
TSILiability / total sama), `TestNilaiAwalCoverageNetRate`, `TestNilaiAwalBasis` (termasuk tepi basis 4 adjustable),
`TestNilaiAwalTidakMengubahPemanggil`; diralat: `TestGantiObjekCoverage` (coverage kosong premi "0"), handler `TestObjekDeductible`
(unit "0"). **Mutasi 16 / 16 mati.**
