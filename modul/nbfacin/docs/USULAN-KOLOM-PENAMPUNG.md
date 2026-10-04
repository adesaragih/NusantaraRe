# Usulan kolom untuk isi penampung `loader.Flatten` — ✅ diputuskan work owner 02-10-2026 (butir 70)

> **Keputusan (butir 70, diteruskan sesi `nusantarare-0f`, "setuju"):** P1, P2–P3, P5, P6, P7 **setuju** sebagaimana
> diusulkan (P2–P3: tabel baru, **bukan** dilipat ke `T_SHIP`, karena di sumber ia daftar berulang). P4: pilihan **(b)** —
> kolom sendiri di `POOLDATA.HISTORYAKSEPTASIPRODUCTION` di samping `POSISI`, **tafsiran** atas K-069 (7b) yang dipilih
> work owner; **menunggu DBA** (nama kolom + tipe). Penerapannya: tiket 22 bab *Amandemen rancangan*.
> ⚠️ **P7 dicabut penerapannya oleh butir 71** (A67, A68): data membantah premis V-47 — semua penunjuk kembali ke
> penampung sampai work owner memutuskan V-47 dan R1 per kunci. ✅ **Butir 72:** V-47 diubah — semua penunjuk disimpan apa
> adanya di 48 kolom teks mentah (tiket 22).

> **Usulan, bukan keputusan.** Keputusan work owner 02-10-2026 (butir 69, diteruskan sesi `nusantarare-0f`, "mau"):
> penampung (ADR-0023) **tidak** disimpan di Oracle; setiap medan di dalamnya harus mendapat **kolom** atau keputusan
> **"dibuang" eksplisit**; untuk `CoverageInitial` dan `AdditionalShip` agent mengajukan usulan kolom — **tidak** menambah
> kolom sendiri. Selama penampung berisi, pemuatan produksi/fase 1 belum boleh dinyatakan selesai (ADR-0023 akibat 2).
> Tipe diusulkan dari **bentuk nilai** yang terukur; nilai kasus tidak disalin ke sini. Satu-satunya nilai yang
> dikutip adalah lima kode enumerasi `IsCedingConfirm` — kode status, bukan data kasus; V-31 sudah menyebutnya.

## Isi penampung yang terukur — 02-10-2026

Jendela: keluaran `loader.Flatten` atas **5 fixture** (JSON, bentuk produksi) dan **115 contoh XML korpus** (diubah ke JSON
di luar repositori, lalu dihapus).

| # | Jalur :: medan | Fixture | Korpus | Bentuk nilai |
| :-: | --- | ---: | ---: | --- |
| P1 | `CargoList/CoverageList :: CoverageInitial` | 104 | 0 | berhuruf, ≤ 21 karakter |
| P2 | `CargoList/(PolicyData/)Ship/AdditionalShip[n] :: DWT`, `GRT`, `NRT` | 104 masing-masing | 0 | bulat, 1 karakter |
| P3 | `CargoList/(PolicyData/)Ship/AdditionalShip[n] :: ID` | 104 | 0 | bulat, ≤ 10 karakter |
| P4 | `ViewSuggest[n] :: IsCedingConfirm` | 34 | 805 | lima nilai enumerasi (`Offer`, `Policy`, `Binding`, `Accepted`, `Retrocession`), ≤ 12 karakter |
| P5 | `CurrencyList :: ID` | 0 | 5 | bulat, 5 karakter |
| P6 | `FacRetroList/PrintRISlip/FacOfferList/CurrencyList/Policy :: TSI` | 0 | 15 | titik desimal, ≤ 46 karakter |
| P7 | penunjuk `Idx*` / `Index*` (V-47) | 600 (321 + 45 + 5 + 208 + 21) | tidak dihitung ulang putaran ini (rincian per jalur: tiket 22) | bulat |

`[terverifikasi]` `AdditionalShip` di fixture selalu **satu** unsur per `Ship` (`[1]`, 104 kali). ⚠️ **Ralat agent:** daftar
"medan tak ada di rancangan" di korpus yang ditulis sebelumnya (`LocationList.TableOfLimit` 186, `…OfferFacIn` 22,
`LocationList/Property.Property` 3, `…SurveyAgent` 2) **bukan data** — `[terverifikasi]` semuanya elemen XML kosong yang
hanya berisi spasi (316 / 34 / 4 / 2 elemen), yang pengubah XML→JSON agent ubah menjadi medan bernilai spasi. Cacat
instrumen pengukuran, bukan isi korpus; di JSON produksi tidak ada.

## Usulan per medan

| # | Usulan agent | Alternatif | Dasar |
| :-: | --- | --- | --- |
| P1 | kolom `T_COVERAGELIST.COVERAGE_INITIAL VARCHAR2(500)` | dibuang | teks deskriptif; rancangan memberi teks nama/keterangan `VARCHAR2(500)` (mis. `CONDITIONS`, `ACCUMULATION_DESCRIPTION`) |
| P2–P3 | tabel baru `T_ADDITIONALSHIP` (anak `T_SHIP`, berulang: `SEQ_NO`, `ROW_UID`), kolom `DWT`, `GRT`, `NRT` `VARCHAR2(50)`, `ADDITIONAL_SHIP_REF_ID VARCHAR2(50)` | dilipat ke `T_SHIP` (bila selalu satu unsur), atau dibuang | tiga nama sama dengan kolom `T_SHIP` yang rancangan beri `VARCHAR2(50)`; `ID` dinamai ulang pola **V-49** (`SHIP_REF_ID`) |
| P4 | ⛔ **tidak diusulkan satu tabel** — tiga pilihan; **menunggu work owner** (butir 69.4) | (a) `T_WORK_POLIS.IS_CEDING_CONFIRM VARCHAR2(50)` = nilai baris terakhir, sepola V-48 `POSISI`; (b) kolom di `POOLDATA.HISTORYAKSEPTASIPRODUCTION` — tabel lama, butuh DBA; (c) tabel riwayat flat baru — V-31 sudah menolak `T_VIEWSUGGEST` | K-069 (7b) salinan repo `docs/00-KEPUTUSAN-WORK-OWNER.md` L3788–3798 (dulu dikutip L3770–3780 dari salinan lama `D:\migrasi\RNM\OUTPUT\`): "Keputusan work owner: **kolom sendiri**" — tabel dan tipe tidak disebut. Nilainya per **baris riwayat**; (a) kehilangan riwayat, dan "baris terakhir" ViewSuggest belum terdefinisi |
| P5 | kolom `T_CURRENCYLIST.CURRENCY_REF_ID VARCHAR2(50)` | dibuang | angka rujukan master sebangun `Currency/ID` (V-49) — bukan kode mata uang (A60) |
| P6 | kolom `T_FR_CURRENCYLIST.POLICY_TSI` — tipe mengikuti keputusan presisi | dibuang | Policy ditolak dilipat (A53) karena TSI-nya menimpa `CurrencyList.TSI`; apakah keduanya selalu sama **belum diukur** |
| P7 | *(dicabut butir 71 — lihat kepala dokumen)* klasifikasi per medan menurut V-47: penunjuk induk langsung → **dibuang** (V-47 sudah memutuskannya); penunjuk leluhur → FK (tiket 22 butir 4, menunggu aturan) | — | V-47: "Yang menunjuk INDUK LANGSUNG dibuang karena PARENT_ID sudah menyatakan hal yang sama" |

Bila disetujui, kolom baru masuk rancangan (DDL draf / workbook adalah sumber READ-ONLY — perubahannya keputusan work
owner) dan tiket 23; Flatten dan generatornya menyusul.

*Tanpa nama orang, tanpa nilai kasus — hanya nama medan, hitungan, bentuk, dan lima kode enumerasi P4.*
