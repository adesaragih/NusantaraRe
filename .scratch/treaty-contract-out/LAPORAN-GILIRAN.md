# Laporan giliran — Treaty Contract Out (cabang `modul/treaty-contract-out`)

> Satu bab per tiket. Angka uji dijalankan ulang tiap commit. Bab **TELEMETRI EKSEKUSI** diisi dari yang
> terukur di dalam sesi; angka token/biaya sejati TIDAK terlihat dari dalam sesi dan tidak dikarang
> (sumber sah: `/cost` atau `claude --print --output-format json`).

## Tiket 01 — skema relasional + migrasi + tipe dirapikan (PREFACTOR)

**Commit:** lihat `git log --oneline` baris `treaty-contract-out: tiket 01`.

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| migrasi | `APP_RNM/internal/repository/migrations/300`–`306` (+ `_down`) | `T_TREATYYEAR`, `T_TREATYCONTRACT`, `T_TREATYREINSURER`, `T_MTREATYSECURITY`, `T_TREATYBUSINESS`, `T_PROPORTIONALARRG`, `T_TREATYCO_JEJAK`; 7 sequence; 7 index. Uang/persen `NUMBER(38,8)`, tanggal `DATE`, nama kolom VERBATIM warisan |
| dokumen | `.scratch/treaty-contract-out/STRUKTUR-TABEL-TREATY-CONTRACT-OUT.md` | acuan tunggal kolom; dijaga `TestKolomDDLCocokDenganStruktur` + `TestGolonganTipeDDLCocokDenganStruktur` |
| models | `internal/models/tco_arrangement.go` | enam entitas bertipe rapi (`*apd.Decimal`, `time.Time`); sentinel `ParentReinsTypeTanpaInduk = "00"` |
| repository | `tco_migrasidata.go` (murni), `tco_pindah.go` (Oracle), `tco_identitas.go`, `tco_kontrak_hilir.go` | konversi + laporan temuan; pemindahan satu transaksi + rekonsiliasi tepat + penyelarasan sequence; identitas `'1'+lpad` yang gagal terang bila tidak muat; 6 pembaca hilir read-only |
| skemauji | `skemauji/tco_tiruan.go` + 2 sisipan di `skemauji.go` | tiruan 6 tabel warisan bertipe deklarasi DBA; pengisi fixture |
| services / cmd | `services/tco_migrasi.go`; flag `-migrate-data-treaty-contract-out` di `cmd/api/main.go` | pintu pemindahan; menolak `IS_PEGA_PROD=true` |
| uji | `tco_migrasidata_test.go`, `tco_pindah_test.go`, `tco_identitas_test.go`, `tco_kontrak_hilir_test.go`, `tco_migrasi_test.go` (tanpa Oracle); `tco_pindah_db_test.go` (bertag `db`) | 32 uji murni/penjaga baru; 3 uji `db` (SKIP tanpa skema uji) |

### Kode bersama yang disentuh (aditif / dipersempit, dilaporkan)

| Berkas | Perubahan | Sebab |
| --- | --- | --- |
| `repository/migrasi_test.go` | `TestSeluruhFKPohonPolisBerkaskade` dan `TestSetiapFKPohonPolisBerindex` **dipersempit** ke rentang 050–079 lewat `milikPremiumList` (baru); cacah CREATE 50→71; cacah CREATE TABLE 20→27 | keduanya menuntut FK berkaskade pada SETIAP migrasi ≥050, padahal anak modul ini menggantung pada kunci gabungan tanpa FK (spec §2). Kebijakan FK modul ini dijaga sendiri di `tco_migrasi_test.go` |
| `repository/strukturkolom_test.go` | `letakStruktur` +1 dokumen | penjaga STRUKTUR wajib mengenal dokumen modul |
| `repository/batasanpemakaian_test.go` | cacah pemanggil `skemauji.Buka()` 9→10 | uji `db` baru membuka koneksi sendiri |
| `repository/skemauji/skemauji.go` | `Pasang` +1 gelung tiruan TCO; `Bongkar` daftar tiruan +6 | tiruan warisan modul ini |
| `cmd/api/main.go` | +1 flag, +1 fungsi, eksklusi flag | pintu tco2 |

### Ralat tiket bertanggal 28-09-2026 (XML membantah)

Ditulis di bab *Pembacaan ulang XML* tiket 01: (1) tiket 02 menyebut `InboxTreatyContractReinsType` "layar master jenis reasuransi" — XML: ia editor KONTRAK per tahun + grid business/reinsurer/security; RD non-Old dipakai pemilih di form kontrak, bukan layar master terpisah. (2) brief §2 "kelompok sudah ada di sidebar" — belum ada. (3) xlsx hanya memuat pohon `InboxTreatyContract`. (4) tiket 05 "existing tidak menolak total ≠ 100" — XML menolak total > 100 (`SaveTreatyReinsurerDetail1_Act` b1452/b1561); = 100 memang tidak diwajibkan. (5) "dua tabel lain" tco1 = `T_TREATYCO_JEJAK` (tiket 01) + tabel lampiran tahun treaty (tiket 12).

### OQ dibuka

`OQ-TCO-01` bentuk teks tanggal warisan (tujuh bentuk yang dikenal pengurai, sisanya dilaporkan);
`OQ-TCO-02` arti/bentuk `IUDATE`; `OQ-TCO-03` isi hidup `PROPORTIONALLIST`/`OBJECT` (dicacah laporan);
`OQ-TCO-04` arti `PROPORTION`; `OQ-TCO-05` label "Underwriting Year"/"Transaction Year" yang bersilang
di `InputDtlTreatyContact.xml` b8004/b8284. Tidak ada OQ ditutup.

### Angka uji

| Perintah | Hasil |
| --- | --- |
| `go build ./...` · `go vet ./...` · `go vet -tags=db ./internal/repository/` · `gofmt -l .` | bersih |
| `go test ./...` | seluruh paket `ok`; 744 PASS (termasuk sub-uji), 0 FAIL, 0 SKIP |
| `go test -tags=db ./internal/repository/ -run Pindahkan` | 3 SKIP (Oracle tidak dikonfigurasi di sesi ini) |
| `npm run typecheck` · `vitest run` | bersih; 372 lulus di 34 berkas (frontend tidak disentuh tiket 01) |

### TELEMETRI EKSEKUSI — tiket 01

| Ukuran | Nilai |
| --- | --- |
| Berkas dibaca (aturan, spec, tiket, kode pola, ADR, penjaga) | ±70 berkas; spec 50 KB, 12 tiket 79 KB, CLAUDE.md 16 KB, tiga brief induk ±35 KB, ±25 berkas Go/TS pola |
| Berkas XML korpus disensus | 3 harness, 12 section, 18 activity, 24 RDBList (14 modul ini + 10 hilir), 7 report definition, 2 data page, 1 xlsx (r1–r139 dibaca) |
| Perintah dijalankan | ±45 (grep/sed korpus, go test/vet/build ×6, vitest ×2, py ×4, git ×5) |
| Berkas ditulis / disunting | 26 baru (+2.9k baris: 14 SQL, 7 Go produksi/uji, 3 dokumen, 1 skemauji, 1 services), 5 disunting (+89/−9) |
| Putaran instrumen gagal lalu diulang | 2 — (a) heredoc Bash mengubah `\n` menjadi baris baru di `main.go` (jebakan backslash yang tercatat di memory) → diperbaiki lewat Edit; (b) STRUKTUR sempat ditulis ke `.scratch` checkout utama → dipindah ke worktree |
| Token / biaya | tidak terlihat dari dalam sesi — tidak dikarang |

## Tiket 02 — jenis reasuransi: master dibaca + saringan non-life

**Commit tiket 01:** `1872d26`. **Commit tiket 02:** lihat `git log --oneline` baris `treaty-contract-out: tiket 02`.

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| repository | `tco_jenisreasuransi.go` (+`_test`) | `DaftarNonLife`: `FLAG = :1 AND TYPE IN (:2..:4) AND ID NOT LIKE :5..:16`, `ORDER BY NOTE ASC`; tabel kebenaran `LolosSaringanNonLifeTCO`; uji sisi korpus atas `pyFilterValue`/`NotStartsWith` |
| services | `tco_jenisreasuransi.go` (+`_test`) | pembaca disuntik (bawaan gagal terang); `WajibIdentitas`; master kosong → `ErrMasterJenisReasuransiKosong` |
| handlers | `rute_treaty_contract_out.go` (+`_test`), `tco_db_test.go` (`db`) | `GET /api/treaty-contract-out/jenis-reasuransi` → `{daftar, total}`; 401/403/503/400/500; seam HTTP dengan Router ber-stub identitas terhadap master tiruan berisi 12 blacklist + awalan + flag + tipe 4 |
| skemauji | `tco_tiruan.go` | tiruan `REINSURANCETYPE` + `IsiJenisReasuransiTCO` |
| penjaga | `tco_nama_jujur_test.go`, `namaJujur.test.ts` | nol pengenal `Old`/`testing`/`JSON_KLAIM` di kode modul (Go dan React) |
| frontend | `assets/labels.treaty-contract-out.ts` (+`_test`), `components/treaty-contract-out/PilihJenisReasuransi.tsx` (+`_test`), `services/api.ts` (aditif) | label VERBATIM b151/b359/b2652/b7925 dijaga terhadap korpus; pemilih tunggal berbasis `Pilih` |

### Kode bersama yang disentuh (aditif, dilaporkan)

| Berkas | Perubahan |
| --- | --- |
| `handlers/handlers.go` | +1 baris `daftarkanRuteTreatyContractOut(mux, svc, stubPelaku)` |
| `handlers/penyuntikan_test.go` | +1 entri peta suntikan (`rute_treaty_contract_out.go`) |
| `repository/batasanpemakaian_test.go` | cacah pemanggil `skemauji.Buka()` 10→11 |
| `repository/tco_migrasi_test.go` (milik modul) | penjaga baca-saja diperluas ke `masterDibacaSajaTCO` |
| `frontend/src/services/api.ts` | +1 bagian (2 interface, 1 fungsi) |

### Ralat / OQ

Ralat bertanggal di tiket 02: (1) harness `InboxTreatyContractReinsType` = editor kontrak, bukan layar master
(RD non-Old = pemilih form kontrak); (2) operator saringan `NotStartsWith` b581, bukan `NOT IN` — SQL `NOT LIKE`;
(3) nama RD menyebut `OLD_LJT_ID` yang tidak ada di saringannya. OQ dibuka: OQ-TCO-06 (pemilih kontrak memakai
blacklist), OQ-TCO-07 (nama fisik `SOANote`/`Code`).

### Angka uji

| Perintah | Hasil |
| --- | --- |
| `go build` · `go vet ./...` · `go vet -tags=db ./internal/...` · `gofmt -l .` | bersih |
| `go test ./...` | seluruh paket `ok`; 758 PASS, 0 FAIL |
| `go test -tags=db` | 3 uji HTTP TCO + 3 uji pemindahan SKIP tanpa Oracle |
| `npm run typecheck` · `vitest run` | bersih; 392 lulus di 37 berkas |

### TELEMETRI EKSEKUSI — tiket 02

| Ukuran | Nilai |
| --- | --- |
| Berkas dibaca | ±12 (dua RD ±1.000 baris terbelah, 3 SQL korpus lintas modul, katalog DEV, 4 pola kode Go/TS) |
| Berkas XML korpus disensus | 2 RD + 3 RDBList lintas modul + 1 sensus `grep` seluruh korpus (`from reinsurancetype`: 18/4/4/3/2) |
| Perintah dijalankan | ±14 |
| Berkas ditulis / disunting | 12 baru (+±900 baris), 5 disunting (+±40) |
| Putaran instrumen gagal lalu diulang | 2 — regex nama jujur tidak menangkap `TypeOld()` (uji instrumen sendiri yang menemukannya, diperbaiki); heredoc Bash mengubah backslash regex → dipindah ke Edit |
| Token / biaya | tidak terlihat dari dalam sesi — tidak dikarang |

## Tiket 03 — tahun treaty: CRUD, gerbang periode, anti-dobel

**Commit tiket 02:** `1f2aab5`. **Commit tiket 03:** lihat `git log --oneline` baris `treaty-contract-out: tiket 03`.

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| models | `tco_tahun.go` (+`_test`) | `PeriksaTahunTreaty` (grup wajib b388, tahun wajib b411, tahun angka b434/CheckYear) + `PeriksaPeriodeTCO` (AC 9, pesan menyebut kedua medan) |
| repository | `tco_tahun.go`, `tco_jejak.go`, `tco_gruptreaty.go` (+`tco_tahun_test`, `tco_pindai_uji_test`) | daftar `ID DESC` berhalaman; ambil; sisip (ID `SEQ_T_TREATYYEAR`); perbarui SELURUH kolom; `CariDobel` (TRUNC, `(:4 IS NULL OR ID <> :4)`); jejak `T_TREATYCO_JEJAK`; master `TREATYGROUP` dibaca saja |
| services | `tco_tahun.go` (`TahunTreatyTCO`), `tco_gruptreaty.go` (+`_test`) | gudang + pelaksana transaksi + jam DISUNTIK; `Simpan`: validasi → tx {dobel → 409, sisip/perbarui, jejak}; `GalatTahunTreatyDobel` menyebut ID lain |
| handlers | `rute_treaty_contract_out.go` (+`_test`), `tco_db_test.go` | `GET /tahun`, `POST /tahun`, `GET/PUT /tahun/{id}`, `GET /grup-treaty`; 401/404/409/422/503; POST menolak `id` dari klien, PUT menolak `id` badan ≠ jalur |
| frontend | `pages/treaty-contract-out/InboxTreatyContract.tsx` (+`_test`), `labels.treaty-contract-out.ts` `TAHUN_TCO` (+29 baris korpus diuji), `api.ts` (+`ambilTahunTreaty`, `ambilSatuTahunTreaty`, `simpanTahunTreaty`, `ambilGrupTreaty`) | grid enam kolom VERBATIM, tombol `Add`/`Edit`/`ReinsType`(🔜04)/`List Description`(🔜08), form `Input New Data` lengkap, catatan OQ-TCO-05 |

### Kode bersama yang disentuh (aditif / dipersempit, dilaporkan)

| Berkas | Perubahan |
| --- | --- |
| `frontend/src/assets/labels.ts` | `MODUL.treatyContractOut` +1 (komentar ralat: folder korpus 20) |
| `frontend/src/components/Shell.tsx` | kelompok ke-18 |
| `frontend/src/lib/daftarMenu.ts` | `ModulTetap` +`'tco-tahun'`; `ENTRI_MENU` +1 |
| `frontend/src/App.tsx` | +1 rute halaman |
| `frontend/src/components/Shell.test.ts` · `pages/Beranda.test.ts` | cacah 17→18 kelompok/kartu, 4→5 butir, aktif +1; 14 tanpa butir tetap |
| `frontend/src/lib/daftarMenu.sinkron.test.ts` | penjaga palet dipersempit ke maksudnya (hasil tidak boleh BERKELOMPOK modul belum dimigrasi) — pencocokan potongan kata membuat `NB Treaty In` menemukan `InboxTreatyContract` yang sah |
| `handlers/penyuntikan_test.go` | +2 entri (`svc.GrupTreaty()`, `svc.TahunTreatyTCO()`) |
| `repository/skemauji/tco_tiruan.go` (milik modul) | +tiruan `TREATYGROUP` |

### Ralat / OQ

Lima ralat bertanggal di tiket 03 (kelompok sidebar belum ada; pemilih grup treaty dari `BrowseTreatyGroup_RD`;
`PROPORTION` menyimpan pilihan "Reinsurance Type"; label tahun bersilang OQ-TCO-05; penjaga palet dipersempit).
Nama layanan `TahunTreatyTCO` (bukan `TahunTreaty`) sebab nama itu sudah dipakai Claim Life (`spreading.go`).

### Angka uji

| Perintah | Hasil |
| --- | --- |
| `go build` · `go vet ./...` · `go vet -tags=db ./internal/...` · `gofmt -l .` | bersih |
| `go test ./...` | seluruh paket `ok`; 776 PASS, 0 FAIL |
| `npm run typecheck` · `vitest run` | bersih; 431 lulus di 38 berkas |

### TELEMETRI EKSEKUSI — tiket 03

| Ukuran | Nilai |
| --- | --- |
| Berkas dibaca | ±14 (7 section/activity/RD korpus, 7 pola kode Go/TS) |
| Berkas XML korpus disensus | `InputTreatyContract.xml` (23k baris, dipindai b1–b22800), `InputDtlTreatyContact.xml`, `GridTreatyContract.xml`, `SaveTreatyYear_Act.xml`, `SetTreatyYear_Act.xml`, `NewInputTreatyYear_Act.xml`, `CheckYear.xml`, `BrowseTreatyYear_RD.xml`, `BrowseTreatyGroup_RD.xml`; 29 baris label diverifikasi satu per satu |
| Perintah dijalankan | ±16 |
| Berkas ditulis / disunting | 16 baru (+±1.9k baris), 11 disunting (+±120) |
| Putaran instrumen gagal lalu diulang | 3 — heredoc panjang gagal diurai (0 baris jalan) → skrip lewat Write + `py`; `angkaSaja` sudah ada di models (dipakai ulang); `services.TahunTreaty` bertabrakan dengan Claim Life (diberi akhiran TCO) |
| Token / biaya | tidak terlihat dari dalam sesi — tidak dikarang |
