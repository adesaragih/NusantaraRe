# E04: Enam gerbang penolakan dan empat klep pembatalnya

**What to build:** Enam pemeriksaan yang dapat **menolak** pembuatan kasus endorsement, empat di
antaranya dapat **dibatalkan** oleh tabel pengecualian.

`[terverifikasi]` Polanya seragam per gerbang: ambil data pemeriksa → siapkan penanda klep → panggil
laporan pembatal → pasang galat **kecuali** klep terisi.

⛔ **Hanya empat dari enam gerbang punya klep.** Gerbang "sudah ada endorsement lain yang belum
selesai" dan gerbang "tidak ter-spreading fac out" **tidak dapat dibatalkan**.

⛔ **Klep kosong berarti penolakan berlaku penuh.** Itu perilaku sistem lama dan arahnya aman —
diport apa adanya.

**Asal (Pega).** `Activity\SetErrorBatalEndorsement_Act` langkah 5 (24 sub-langkah) ·
`RDBList\GetEDMStatus_SQL` · `GetDataClaim_SQL` · `SearcStatusBayarArasaps_SQL` ·
`GetListRNWbyNopolis_SQL` · `GetFacoutList_SQL` · ReportDefinition `GetListEdm` ·
`BrowseOpenProteksiEdm_RD`

**Keputusan.** K-049 · K-046 · K-006

**Blocked by:** E03

**Status:** sebagian — 01-10-2026, `backend/services/gerbangtolak.go` (`PeriksaGerbangPenolakan`,
`JenisBisnisUntukPredikat`) + `predikatkasus.go` (`PredikatDari`: langkah 10–16 dari registry EDM);
test `gerbangtolak_test.go` (`TestGerbangLiniLewatRegistry`). Hasil query masukan (E17 dilewati).
`[terverifikasi]` IsLife EDM membaca `pyWorkPage.Quotation.BusinessOldId` (L1…L16), bukan jenis bisnis
"Life" dari langkah 8. ⛔ `[terverifikasi]` Prakondisi yang dieksekusi **tidak memuat satu pun nomor
polis**: 557 dari 560 kemunculan ada di cache editor `pyExpressionGadget/pyExpressionMapNew` (spec §4.1
mencacah 502 — cacah cache). ⛔ Bypass RI slip melompat dua kali: 5.1 → `jmp` (lewati gerbang 1), 5.10 →
`jmp2` (lewati gerbang 3–6); gerbang 2 tetap.

- [x] Keenam gerbang menolak pada kondisi yang benar **dan hanya** pada kondisi itu
- [x] Gerbang pembayaran **hanya** berlaku untuk jenis endorsement pembatalan (`EdmType` 1/2)
- [x] **Empat klep** membatalkan penolakan yang benar, masing-masing dengan penanda jenisnya sendiri
- [x] **Gerbang kedua dan keenam tetap menolak** meski seluruh klep terisi
- [ ] Klep kosong atau tabel tak terjangkau → **penolakan berlaku penuh** — klep kosong ✅ (klep yang tidak ada di peta = kosong); "tabel tak terjangkau" adalah kegagalan pelaksana query (repository, E17) yang belum ada — perilakunya kelak diuji di sana
- [x] Bypass "endorsement RI slip" — `TestBypassRISlip`
- [x] Rule pengecekan pembatal diimplementasikan sebagai keputusan atas hasil klep; isi tabel di luar lingkup
- [x] Rujukan laporan bersifat runtime — pelaksananya repository; ketiadaannya bukan `panic`
- [x] Tiap gerbang menyebut rule Pega asalnya dalam komentar (§4.6)
- [x] Tambahan dari korpus: langkah 7–19 (jenis bisnis kosong → "Invalid policy no!"; tidak satu pun predikat lini → "Invalid business, please contact IT!"; polis berpola RNML dibaca Life sebelum predikat)
