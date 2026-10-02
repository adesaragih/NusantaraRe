# Fixture daftar-izin premi — kasus nyata PA dan MBU

`premi.json`: masukan dan premi sistem lama dari kasus nyata `D:\migrasi\RNM\DDL\CONTOH`, diambil dengan
**daftar-izin medan angka** — tanpa nomor kasus, nama, alamat, atau kunci lain (keputusan work owner
01-10-2026, `docs/KEPUTUSAN-30-09-2026.md` butir 16). Berkas mentah tidak masuk repositori.

| Kasus | Lini | Medan | Catatan |
| --- | --- | --- | --- |
| 4 | PA | TSI, Rate, ProRatePercent, CalculateMethod_FacIn, Premium | tidak satu pun memuat `Discount` — premi lama cocok bila diskon nol (butir 17). "metode 3 #2" membedakan pembulatan: 504231451.61275755 → sistem lama .6128, potong memberi .6127 |
| 4 | MBU (MOTOR VEHICLE) | TSI, Rate, Loading, ProRatePercent coverage, Premium | `[terverifikasi]` seluruh 93 baris CoverageList kelima kasus MBU cocok eksak dengan L1144 (Python decimal, 01-10-2026); empat ini mewakili ragam rate dan ejaan pro-rata ("100.000…" dan "100"). Tidak satu pun memicu pembulatan — bentuk bersarang dibuktikan `TestPremiMBUBersarang` |

Dibaca `TestRekonsiliasiDaftarIzin` (paket ini) dan kerangka `services/rekonsiliasi` (tiket 16). Dipindah dari
literal Go ke JSON 01-10-2026, tanpa mengubah satu angka pun (diekstrak dengan skrip dari literalnya).

## `mbu_mata_uang.json` — agregat per mata uang (tiket 07)

Lima kasus **NB** MOTOR VEHICLE yang sama (kasus EDM tidak: jalur EDM tidak menulis premi ke CurrencyList), 93
coverage. Daftar-izin per coverage: TSI, Rate, Loading, ProRatePercent, `Currency.Name`, `FlagDelete`; per kasus:
`CurrencyList` (Name, Premium) di **akar** halaman — bukan `CedingCedantList(*).CurrencyList` (premi × pangsa) dan
bukan `OldData`. Premi lama per coverage sengaja tidak disalin: itulah yang dihitung. Ekstraksi diverifikasi Python
`decimal` (01-10-2026): ke-93 premi coverage = L1144; jumlah per mata uang = CurrencyList.Premium pada 4 dari 5 kasus.
`urutanMasterMataUang` diambil dari urutan CurrencyList tersimpan `[dugaan]` (isi tabel master tidak ada di korpus;
semua kasus hanya IDR).
