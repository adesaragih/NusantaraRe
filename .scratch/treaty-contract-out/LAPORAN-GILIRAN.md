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

## Jeda — 28-09-2026 (atas permintaan work owner, sesudah tiket 03)

Giliran dijeda atas perintah work owner, sebelum tiket 04. Aturan berhenti GILIRAN-6 (laporan pertama paling cepat
sesudah tiket 04) dikalahkan oleh perintah jeda langsung. Konteks sesi tidak menipis; jeda ini bukan karena konteks.

| Butir | Keadaan |
| --- | --- |
| Tiket terakhir selesai | **03** — tahun treaty, periode, anti-dobel — commit `678fd25` (sebelumnya 01 `1872d26`, 02 `1f2aab5`) |
| Tiket yang dijeda | **12** — lampiran di tahun treaty (urutan ke-4 brief, sesudah 03) |
| Pekerjaan tiket 12 yang sudah ada di disk | **nol berkas.** Yang dilakukan baru pembacaan: korpus rantai lampiran dan jalur dokumen yang ada di `main` |
| Commit WIP | tidak dibuat — pohon kerja bersih saat jeda (`git status` kosong sesudah `678fd25`) |
| Bab ini | ditulis sesudah pemeriksaan itu dan **belum di-commit** |
| Sisa urutan | 12 → 04 → 05 → 07 → 08 → 06 → 11 → 09 → 10, lalu uji penuh dan `/code-review` |

### Temuan pembacaan tiket 12 (belum dituangkan ke kode)

Nomor baris = baris mentah berkas korpus (satu tag per baris), diverifikasi dengan `awk 'NR==n'`.

| Unsur | Bukti | Arah yang direncanakan |
| --- | --- | --- |
| panel lampiran | `Section/InputTreatyContract.xml` b11721 `Attachment for`, b13074 menyertakan `GridTreatyArrangementAttachment` | panel di form tahun treaty yang sudah punya ID |
| label panel | `Section/GridTreatyArrangementAttachment.xml`: b578 `Add attachment`, b1023 `Refresh` (→ `LoadAttachmentTreatyOut`), b1785 `For Treaty Contract Out`, b2391 `Download` (→ `DownloadAll_Act`), b2659 `Download All` (→ `TreatyOutDownloadAll_Act`), b3032 `File Name`, b3170 `Type` (sel `.pyCategory` b3705), b3470 tautan nama berkas (→ `TreatyOutDownloadOne`), b3897 `Delete` (→ `DeleteAttachmentTreaty`) | label VERBATIM, diuji terhadap korpus |
| kunci lampiran di Pega | `Activity/TreatyOutSaveAttachment.xml` b1402 dan `DeleteAttachmentTreaty.xml` b252: `TreatyYear + TreatyYearID`; `RDBList/GetAllAttachment2_Sql.xml` membaca `M_ATTACHMENTTREATY_2 where treatyid = {TreatyIn.ID}` | tidak dibawa; kunci baru = ID tahun treaty; uji statik menolak rujukan ID treaty inward |
| penulisan di Pega | `RDBList/InsertAtatchment_Sql.xml` → `POOLDATA.PEGA_M_ATTACHMENT` dengan CLOB JSON dan `COMMIT` | tidak dibawa (prosedur tidak dipanggil, nol COMMIT, nol JSON) |
| pesan tanpa berkas | `TreatyOutSaveAttachment.xml` b376 `"Tidak ada file yg diattach"` | dibawa sebagai penolakan berkas kosong |
| master kategori | `RDBList/CategoryAttach_SQL.xml` `select * from CATEGORY_ATTACH_REAS order by note`; `Activity/SetCategoryAttachTreatyin.xml` b500 memakai `.NOTE` | dibaca saja, kolom `NOTE` saja; kolom ID master tidak terbukti di korpus, jadi tidak dipakai |
| jalur dokumen `main` | `services/unggahan.go` (folder `UNGGAHAN_DIR`, batas 25 MiB, outbox `storage-unggah`/`storage-hapus`, pelaksana stub `PelaksanaBerkasLokal`), `repository/efekkeluar.go` (`AntreEfek`, `PungutEfek`, `TuntaskanEfek`), `services/efekkeluar.go` (`ResolverEndpoint`, `KunciUnggahBerkas`, `LayakDicobaUlang`, `ErrPenyimpananBelumDisetujui`), `services/tokenstorage.go` | dipakai ulang; tabel dan fungsi milik Claim Life tidak ditulis |

### Keputusan rancangan yang sudah diambil untuk tiket 12

1. **Migrasi 307** `T_TREATYYEAR_LAMPIRAN`: `ID`, `IDTREATYYEAR` (FK ke `T_TREATYYEAR` tanpa kaskade), `FILENAME`,
   `FILEMIMETYPE`, `CATEGORY` (teks `NOTE`), `IMAGEID` (kunci penyimpanan acak, lahir bersama barisnya), `T_STORAGE_ID`
   (terisi hanya setelah unggahan dipastikan), `UKURAN NUMBER(19)`, `USERID`, `TGLUPLOAD`. Sequence
   `SEQ_T_TREATYYEAR_LAMPIRAN`. Yang ikut diperbarui: bab STRUKTUR, `sequenceDikenalTCO`, kunci cacah CREATE, dan uji
   kebijakan FK (2 → 3 FK).
2. **Pengulangan tidak menggandakan berkas** karena kunci penyimpanan (`IMAGEID`) tetap per lampiran. Pelaksana efek
   bersifat idempoten: baris sudah dihapus berarti selesai tanpa kerja, dan `T_STORAGE_ID` yang sudah terisi juga
   berarti selesai. Bila berkas antrean sudah hilang, pelaksana menanyakan penyimpanan apakah kunci itu ada.
3. **Pekerja outbox milik modul ini sendiri.** `PekerjaEfek` bawaan tidak dipakai apa adanya. Jalur menyerahnya menulis
   jejak Claim Life (`PerekamJejakOracle` → `KlaimLife.SisipJejak`), dan itu penulisan lintas modul. Pekerja modul ini
   memakai ulang fungsi outbox repository, `Backoff`, dan `LayakDicobaUlang`, lalu menulis jejaknya ke `T_TREATYCO_JEJAK`.
4. **Tidak ada pekerja outbox yang berjalan di `cmd/api`** untuk modul mana pun (nol pemanggil `SatuPutaran` di kode
   produksi). Karena itu efek dijalankan sesudah commit unggahan. Endpoint `ulangi` memakai jalan yang sama.
5. **Hapus selalu mengantre `storage-hapus`.** Objek yang sudah tidak ada di penyimpanan tidak menggagalkan penghapusan.
6. **Klien penyimpanan di balik antarmuka.** Di DEV dipakai stub lokal di bawah `UNGGAHAN_DIR`. Rangkaian jarak jauh
   terdiri dari resolver runtime (`ResolverEndpoint` + `KunciUnggahBerkas`), cache token berbatas margin, dan transport.
   Transport tidak diimplementasikan dan gagal dengan `ErrPenyimpananBelumDisetujui`, yang permanen. Endpoint luar
   tidak dipanggil.
7. **Unduhan di layar lewat `fetch` + Blob**, bukan `<a href>`, supaya header identitas `X-Pelaku` ikut terkirim.

⚠️ Pengamatan di luar modul ini, **tidak diverifikasi dengan menjalankannya dan tidak disentuh**: tautan unduh dokumen
Claim Life (`PanelDokumenPeserta.tsx`, `<a href>`) tidak membawa header identitas. Dalam mode `AUTH_STUB` ia tampaknya
akan dijawab 401.

### Langkah berikutnya (sesi lanjutan)

1. Mulai dari cabang `modul/treaty-contract-out` @ `678fd25` di worktree `.worktrees/treaty-contract-out`.
2. Kerjakan tiket 12 menurut rancangan di atas. Bab "Pembacaan ulang XML" dapat diambil dari tabel temuan di atas.
3. Lanjutkan 04 → 05 → 07 → 08 → 06 → 11 → 09 → 10, satu commit per tiket.
4. Laporan pertama ke work owner paling cepat sesudah tiket 04, sesuai GILIRAN-6.
5. Di akhir: uji penuh (`go test ./...`, `go vet -tags=db`, `npm run typecheck`, `vitest run`), lalu `/code-review`.

### TELEMETRI EKSEKUSI — tiket 12 (dijeda, pembacaan saja)

| Ukuran | Nilai |
| --- | --- |
| Berkas korpus dibaca | ±16: `GridTreatyArrangementAttachment.xml`, `TreatyOutAttachContent.xml`, `TreatyOutSaveAttachment.xml`, `LoadAttachmentTreatyOut.xml`, `TreatyOutDownloadOne.xml`, `TreatyOutDownloadAll_Act.xml`, `DeleteAttachmentTreaty.xml`, `LoadAttachment.xml`, `SetCategory_act.xml`, `SetCategoryAttachTreatyin.xml`, 9 RDBList lampiran/penyimpanan, `InputTreatyContract.xml` (b11700–b13240), dan 4 `CategoryAttach`/`AttachmentLife` modul lain untuk kolom master |
| Berkas kode dibaca | ±14: `unggahan.go`, `dokumen.go`, `dokumenunggah.go`, `efekkeluar.go` (services + repository), `antrean.go`, `tokenstorage.go` (services + repository), `linkservice.go`, `handlers/dokumen.go`, `efekkeluar_statik_test.go`, `PanelDokumenPeserta.tsx`, `api.ts` (bagian dokumen), migrasi 007/015/019/020 |
| Perintah dijalankan | ±24, semuanya baca, ditambah satu penulisan bab ini |
| Berkas ditulis / disunting | 1: bab ini (+±75 baris) |
| Putaran instrumen gagal lalu diulang | 1: satu loop `for f in $(grep -rl …)` memecah jalur korpus yang berspasi, lalu diulang dengan jalur yang dikutip |
| Token / biaya | tidak terlihat dari dalam sesi, jadi tidak dikarang |

## Lanjutan 1 — titik awal (29-09-2026)

Brief `PROMPT-LANJUTAN-TREATY-CONTRACT-OUT-1.md` (commit `918eeed`). Pohon kerja bersih di `918eeed`; `0022865` (bab jeda)
dan `678fd25` (tiket 03) ada di riwayat. Rancangan tiket 12 dari bab Jeda dipakai tanpa dirancang ulang.

| Perintah | Hasil titik awal |
| --- | --- |
| `go vet ./...` · `go vet -tags=db ./...` · `gofmt -l .` | bersih |
| `go test ./...` | 776 lulus, 0 gagal |
| `go test -tags=db ./...` | 776 lulus, 46 dilewati (tanpa skema uji: `ORACLE_DSN` tidak dikonfigurasi di sesi ini), 0 gagal |
| `npx tsc --noEmit` · `npx vite build` | bersih |
| `npx vitest run` | 431 lulus di 38 berkas |

## Tiket 12 — lampiran di tahun treaty (FITUR BARU)

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| skema | `migrations/307_t_treatyyear_lampiran.sql` (+`_down`), bab STRUKTUR `T_TREATYYEAR_LAMPIRAN` | tabel kedelapan tco1; FK ke tahun tanpa kaskade; `IMAGEID` unik; sequence lebar 9 |
| models | `tco_lampiran.go` (+uji) | status turunan, kategori VERBATIM master, nama berkas antrean dan entri zip yang aman |
| repository | `tco_lampiran.go`, `tco_jejak.go` (+3 aksi) | rekam lampiran; efek unggah terakhir dibaca dari outbox bersama (`ROW_NUMBER … DIBUAT DESC`, disaring `MODUL`) |
| services | `tco_lampiran.go`, `tco_penyimpanan.go` (+uji) | layanan lampiran dengan pekerja modul sendiri; stub lokal; rangkaian jarak jauh; cache token |
| handlers | `tco_lampiran.go` (+uji, +uji `db`), `rute_treaty_contract_out.go` (+1 panggilan, +pemetaan galat) | 8 rute |
| frontend | `PanelLampiranTahun.tsx` (+uji), `labels.treaty-contract-out.ts` `LAMPIRAN_TCO`, `api.ts` (+9 ekspor), `InboxTreatyContract.tsx` (+panel) | unduhan lewat `fetch` berheader identitas; nol `<a href>` |

### Keputusan yang mengikat

- Lampiran opsional (AC 55): kegagalan penyimpanan hanya mengubah status; rekam, jejak, dan berkas antrean tetap.
- Pengulangan tidak menggandakan berkas (AC 58): kunci penyimpanan adalah `IMAGEID` rekam itu sendiri; pelaksana
  idempoten di tiga jalan (rekam terhapus, sudah terkirim, berkas antrean hilang tetapi penyimpanan memilikinya).
  Dibuktikan dengan mutasi: kunci acak per percobaan membuat uji pengulangan gagal (2 berkas, mau 1).
- Hapus berkas yang sudah tidak ada tidak menggagalkan penghapusan rekam. Dibuktikan dengan mutasi: toleransi dicabut,
  dua uji gagal.
- Menyerah (permanen atau jatah percobaan habis) masuk jejak `T_TREATYCO_JEJAK` di transaksi yang sama.
- Alamat di-resolve saat jalan per operasi (AC 59); penjaga statik modul (`TestTCOLampiranTanpaAlamatLiteral`)
  dibuktikan menggigit dengan berkas mutan.

### Kode bersama yang disentuh (aditif / dipersempit, dilaporkan)

| Berkas | Perubahan |
| --- | --- |
| `repository/migrasi_test.go` | kunci cacah CREATE 71→74 dan CREATE TABLE 27→28, dengan komentar bertanggal |
| `handlers/penyuntikan_test.go` | +1 entri `tco_lampiran.go` (lima penyuntikan wajib) |
| `repository/skemauji/tco_tiruan.go` (milik modul) | +tiruan `CATEGORY_ATTACH_REAS`, `IsiKategoriLampiranTCO` |
| `handlers/tco_db_test.go` (milik modul) | server uji memakai `UNGGAHAN_DIR` sementara |
| `handlers/rute_treaty_contract_out_test.go` (milik modul) | larangan jalur hapus tahun dipersempit ke jalur tahun itu sendiri dan diperluas ke `tco_lampiran.go` |
| fungsi bersama yang DIPANGGIL, tidak diubah | `PohonKlaim.AntreEfek/PungutEfek/TuntaskanEfek`, `Backoff`, `LayakDicobaUlang`, `Unggahan.tulisBerkas`, `models.MimeDokumen`, `models.ImageIDBaru`, `ResolverEndpoint`, `KunciUnggahBerkas/KunciURLBerkas/KunciHapusBerkas` |

### Ralat / OQ

Empat ralat bertanggal di tiket 12: paket `internal/clients` tidak ada; cache token dibangun di atas antarmuka karena
jalur token nyata menuntut persetujuan; tidak ada pekerja outbox di `cmd/api`; kolom ID master kategori tidak terbukti.
**OQ-TCO-08 (baru):** sumber token Oracle untuk penyimpanan nyata (`GET_TOKEN_STORAGE` + garam) belum dipasang, menunggu
persetujuan penyambungan. **OQ-TCO-09 (baru):** pekerja latar yang menjalankan `SatuPutaran` berkala belum ada; lampiran
yang gagal sementara dicoba lagi saat aksi berikutnya atau lewat `Ulangi`.

### Angka uji

| Perintah | Hasil |
| --- | --- |
| `go vet ./...` · `go vet -tags=db ./...` · `gofmt -l .` | bersih |
| `go test ./...` | 808 lulus, 0 gagal (+32 dari titik awal) |
| `go test -tags=db ./...` | 808 lulus, 47 dilewati (+1: `TestLampiranTahunTreatyLingkaranPenuh`, tanpa skema uji), 0 gagal |
| `npx tsc --noEmit` · `npx vite build` | bersih |
| `npx vitest run` | 452 lulus di 39 berkas (+21) |

⚠️ SQL tiket 12 belum pernah dijalankan terhadap Oracle: uji `db` dilewati karena skema uji tidak dikonfigurasi di sesi ini.

### TELEMETRI EKSEKUSI — tiket 12

| Ukuran | Nilai |
| --- | --- |
| Berkas dibaca | ±30 (±16 korpus pada sesi jeda + 2 pemeriksaan cacah `ServiceGoogle`/`LinkService`; ±14 pola kode) |
| Berkas XML korpus disensus | `GridTreatyArrangementAttachment.xml`, `InputTreatyContract.xml` (b11700–b13240), `TreatyOutSaveAttachment.xml`, `SetCategoryAttachTreatyin.xml`, `ServiceGoogle.xml`, `LinkService.xml`; 10 baris label + 5 baris aksi diverifikasi satu per satu |
| Perintah dijalankan | ±45 |
| Berkas ditulis / disunting | 13 baru (+±2.300 baris), 12 disunting (+±200) |
| Putaran instrumen gagal lalu diulang | 5: heredoc panjang gagal diurai (repository ditulis ulang lewat Write); dua kali heredoc memakan garis miring terbalik (uji label dan uji halaman, diperbaiki lewat skrip Write dan Edit); tiga nama pemalsuan bertabrakan dengan uji lain (`efekUji`, `resolverUji`, dll., diberi akhiran); penjaga alamat bersama menangkap pola `://` di uji statik modul sendiri (pola dirakit dari potongan) |
| Token / biaya | tidak terlihat dari dalam sesi, jadi tidak dikarang |

## Tiket 04 — kontrak treaty di dalam tahun

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| models | `tco_kontrak.go` (+uji) | gerbang simpan; tanggal akhir bawaan ditiru apa adanya |
| repository | `tco_kontrak.go` (+uji) | `MasterKontrakTCO`; perbarui dan ambil berbatas tahun |
| services | `tco_kontrak.go` (+uji) | `KontrakTreatyTCO` |
| handlers | `tco_kontrak.go` (+uji, +uji `db`), `rute_treaty_contract_out.go` (+1 panggilan, +pemetaan galat) | 4 rute |
| frontend | `PanelKontrakTahun.tsx`, `InboxTreatyContractReinsType.tsx` (+uji masing-masing), `KONTRAK_TCO`, `api.ts` (+3) | editor dari baris tahun dan dari menu |

### Kode bersama yang disentuh (aditif, dilaporkan)

| Berkas | Perubahan |
| --- | --- |
| `frontend/src/lib/daftarMenu.ts` | `ModulTetap` +`'tco-kontrak'`; `ENTRI_MENU` +1 (label VERBATIM harness b151) |
| `frontend/src/App.tsx` | +1 rute halaman |
| `frontend/src/components/Shell.test.ts` | butir menu 5 → 6 |
| `handlers/penyuntikan_test.go` | +1 entri `tco_kontrak.go` (tiga penyuntikan wajib) |
| `repository/tco_tahun.go` (milik modul, tiket 03) | placeholder anti-dobel `:4`,`:4` → `:4`,`:5` (ralat 8) |

### Ralat / OQ

Delapan ralat bertanggal di tiket 04 (gerbang simpan Pega yang dikomentari; anti-dobel kontrak; gerbang tahun mulai;
tanggal akhir bawaan beranomali; RD grid tidak diekspor; harness menu tanpa konteks; Username; placeholder berulang).
**OQ-TCO-10 (baru):** anomali 366 hari `SetTanggalTreatyContract` langkah 5 — ditiru apa adanya; usul: `mulai + 1 tahun
kalender`. **OQ-TCO-11 (baru):** satu jenis reasuransi satu kontrak per tahun — konfirmasi work owner.

### Kontrak hilir

Tidak berubah. Pembaca hilir tco3 (`KontrakHilirTCO`) tidak membaca `T_TREATYCONTRACT`; kontrak hanya membuka kombinasi
yang dibaca hilir dari reinsurer, business, dan klausul.

### Angka uji

| Perintah | Hasil |
| --- | --- |
| `go vet ./...` · `go vet -tags=db ./...` · `gofmt -l .` | bersih |
| `go test ./...` | 823 lulus, 0 gagal (+15) |
| `go test -tags=db ./...` | 823 lulus, 48 dilewati (+1: `TestKontrakTreatyLingkaranPenuh`, tanpa skema uji), 0 gagal |
| `npx tsc --noEmit` · `npx vite build` | bersih |
| `npx vitest run` | 484 lulus di 41 berkas (+32) |

⚠️ SQL tiket 04 belum pernah dijalankan terhadap Oracle: uji `db` dilewati karena skema uji tidak dikonfigurasi di sesi ini.

### TELEMETRI EKSEKUSI — tiket 04

| Ukuran | Nilai |
| --- | --- |
| Berkas dibaca | ±18 (8 korpus: 2 section besar, harness, `PanggilReinsType`, 5 aktivitas; ±10 pola kode Go/TS) |
| Berkas XML korpus disensus | `InputTreatyContractReinsType.xml` (24.261 baris, b1–b12030 untuk tiket ini), `InputTreatyContract.xml` b20760–b21420, `SaveTreatyContract_Act.xml` (12 langkah), `SetTanggalTreatyContract.xml` (9 langkah), `SetUbahTreatyContract.xml`, `NewInputTreatyContract_Act.xml`, `BrowseReinsTypeYear.xml`, `D_TreatyContract_Act.xml`; 28 baris label/aksi diverifikasi satu per satu |
| Perintah dijalankan | ±35 |
| Berkas ditulis / disunting | 12 baru (+±1.500 baris), 11 disunting (+±200) |
| Instrumen baru | `langkah.py` di scratchpad: pembuang langkah aktivitas (prasyarat, penanda `//`, transisi) — dipakai ulang tiket berikut |
| Putaran instrumen gagal lalu diulang | 3: nama pembantu uji `tgl` bertabrakan di models (diganti `tglTeks`); pemalsuan pemindai `pemindaiUji` tidak ada (dipakai `barisPalsu` yang sudah ada); grep langkah aktivitas tak terbaca karena tag bersarang (diganti pembuang XML) |
| Token / biaya | tidak terlihat dari dalam sesi, jadi tidak dikarang |

## Tiket 05 — reinsurer + total share

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| models | `tco_reinsurer.go`, `tco_kombinasi.go` (+uji) | aturan uang dan gerbang reinsurer; kombinasi |
| repository | `tco_reinsurer.go` (+uji), `tco_kontrak.go` `Kunci` | reinsurer per kombinasi; master `AGENT` |
| services | `tco_reinsurer.go` (+uji) | `ReinsurerTCO` |
| handlers | `tco_reinsurer.go` (+uji, +uji `db`) | 4 rute |
| frontend | `PanelReinsurerKombinasi.tsx` (+uji), `REINSURER_TCO`, `api.ts` (+3) | panel dari tombol `Reinsurer List` |

### Bukti merah lebih dulu untuk aturan uang

Uji `UraiPersenMasukTCO`, `TotalShareTCO`, dan `PeriksaTotalShareTCO` dijalankan dulu terhadap implementasi sementara
berbasis `float64`: gagal di tiga uji (`0.1 + 0.2 = 0.30000000000000004`; tiga sepertiga kehilangan skala `100.00000000`;
rentang dan skala tidak ditegakkan). Implementasi `apd` kemudian meluluskannya.

### Kode bersama yang disentuh (aditif, dilaporkan)

| Berkas | Perubahan |
| --- | --- |
| `handlers/penyuntikan_test.go` | +1 entri `tco_reinsurer.go` (empat penyuntikan wajib) |
| `repository/tco_jenisreasuransi.go` (milik modul) | `masterDibacaSajaTCO` +`AGENT` |
| `repository/skemauji/tco_tiruan.go` (milik modul) | +tiruan `AGENT`, `IsiAgentTCO` |

### Ralat / OQ

Enam ralat bertanggal di tiket 05; yang terpenting: total share > 100 DITOLAK di Pega (catatan tiket menyebut total
bukan gerbang). **OQ-TCO-12 (baru):** nama kolom fisik `STATUSACTIVE` master `AGENT` diturunkan dari properti RD.

### Kontrak hilir

Kolom yang dibaca hilir (`KolomReinsurerHilir`: TREATYYEAR, TREATYGROUPID, REINSTYPEID, REINSURERID, CLIENTID, NAME,
RICOMM, PCTSHARE) ditulis VERBATIM; `TestKontrakHilirTCO*` tetap hijau. Penulis kini mengisi `NAME`/`CLIENTID` dari
master, sehingga pembaca hilir mendapat nama perusahaan yang sama dengan master.

### Angka uji

| Perintah | Hasil |
| --- | --- |
| `go vet ./...` · `go vet -tags=db ./...` · `gofmt -l .` | bersih |
| `go test ./...` | 840 lulus, 0 gagal (+17) |
| `go test -tags=db ./...` | 840 lulus, 49 dilewati (+1: `TestReinsurerKombinasiLingkaranPenuh`, tanpa skema uji), 0 gagal |
| `npx tsc --noEmit` · `npx vite build` | bersih |
| `npx vitest run` | 514 lulus di 42 berkas (+30) |

⚠️ SQL tiket 05 belum pernah dijalankan terhadap Oracle: uji `db` dilewati karena skema uji tidak dikonfigurasi di sesi ini.

### TELEMETRI EKSEKUSI — tiket 05

| Ukuran | Nilai |
| --- | --- |
| Berkas dibaca | ±20 (10 korpus: 7 aktivitas, section 14.763 baris, 2 RD, 2 RDBList modul lain untuk kolom `AGENT`; ±10 pola kode) |
| Berkas XML korpus disensus | `SaveTreatyReinsurerDetail1_Act.xml` (15 langkah), `SetErrorMessageReinsurer.xml`, `SetUbahTreatyReinsurerList_Act.xml`, `NewTreatyReinsurerDetail_Act.xml`, `BrowseTreatyReinsurerList_Act.xml`, `SetTreatyReinsurerList_Act.xml`, `ViewDetailTreatyReinsurerGrid1.xml`, `BrowseDetailTreatyReisurer_RD.xml`, `BrowseAgentReinsSOA_RD.xml`, `GetMasterReinsurerList.xml`; 26 baris label/aksi diverifikasi |
| Perintah dijalankan | ±30 |
| Berkas ditulis / disunting | 11 baru (+±1.700 baris), 11 disunting (+±250) |
| Putaran instrumen gagal lalu diulang | 2: harapan total `100` di uji `db` salah (Oracle memberi skala 8 → `100.00000000`, diperbaiki sebelum dijalankan); dua metode bawaan tak terpakai dibuang |
| Token / biaya | tidak terlihat dari dalam sesi, jadi tidak dikarang |

## Tiket 07 — business + nonaktif

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| models | `tco_business.go` (+uji) | gerbang simpan; status aktif sebagaimana hilir membacanya |
| repository | `tco_business.go` (+uji) | `MasterBusinessKombinasiTCO`, master `BUSINESS` |
| services | `tco_business.go` (+uji) | `BusinessTCO` |
| handlers | `tco_business.go` (+uji, +uji `db`) | 5 rute |
| frontend | `PanelBusinessKombinasi.tsx` (+uji), `BUSINESS_TCO`, `api.ts` (+4) | panel dari tombol `Business List` |

### Kode bersama yang disentuh (aditif, dilaporkan)

| Berkas | Perubahan |
| --- | --- |
| `handlers/penyuntikan_test.go` | +1 entri `tco_business.go` (empat penyuntikan wajib) |
| `repository/tco_jenisreasuransi.go`, `skemauji/tco_tiruan.go` (milik modul) | `masterDibacaSajaTCO` +`BUSINESS`; tiruan `BUSINESS` + `IsiBusinessTCO` |

### Ralat / OQ

Enam ralat bertanggal di tiket 07 (grid Pega menyembunyikan baris nonaktif; nilai `IsActive`; pemilih tidak tersaring
grup; anti-dobel kode; hapus Pega berkelas `_LIFE`; `BIZNAME` dari master). **OQ-TCO-13 (baru):** nilai nonaktif `0` dan
tabel `BUSINESSGROUP`.

### Kontrak hilir

`KolomBusinessHilir` (ID, TREATYYEAR, TREATYGROUPID, REINSTYPEID) dan `ISACTIVE` yang disaring hilir (`isactive='1'`)
ditulis VERBATIM. Baris nonaktif berisi `0`, sehingga hilir tetap tidak membacanya.

### Angka uji

| Perintah | Hasil |
| --- | --- |
| `go vet ./...` · `go vet -tags=db ./...` · `gofmt -l .` | bersih |
| `go test ./...` | 854 lulus, 0 gagal (+14) |
| `go test -tags=db ./...` | 854 lulus, 50 dilewati (+1: `TestBusinessKombinasiLingkaranPenuh`, tanpa skema uji), 0 gagal |
| `npx tsc --noEmit` · `npx vite build` | bersih |
| `npx vitest run` | 536 lulus di 43 berkas (+22) |

⚠️ SQL tiket 07 belum pernah dijalankan terhadap Oracle: uji `db` dilewati karena skema uji tidak dikonfigurasi di sesi ini.

### TELEMETRI EKSEKUSI — tiket 07

| Ukuran | Nilai |
| --- | --- |
| Berkas dibaca | ±14 (9 korpus: 6 aktivitas, section 12.263 baris, 2 RD, 2 RDBList; SQL `BUSINESS` modul lain; ±5 pola kode) |
| Berkas XML korpus disensus | `SaveTreatyBusinessDetail_Act`, `NewTreatyBusinessDetail_Act`, `SetUbahTreatyBusinessList_Act`, `DeleteRowBusiness`, `BrowseTreatyBusinessList_Act`, `SetTreatyBusinessList_Act`, `ViewDetailTreatyBusinessGrid`, `BrowseTreatyBusiness_RD`, `BrowseFilterBusiness_RD`, `GetMasterBusinessList`, `DeleteRowBusinessList`; 24 baris label/aksi diverifikasi |
| Perintah dijalankan | ±22 |
| Berkas ditulis / disunting | 11 baru (+±1.300 baris), 10 disunting (+±200) |
| Putaran instrumen gagal lalu diulang | 1: salinan cadangan mutasi sempat ditaruh di `/tmp` alih-alih scratchpad — dipulihkan dan dibuang |
| Token / biaya | tidak terlihat dari dalam sesi, jadi tidak dikarang |

## Tiket 08 — klausul satu tabel, 25 jenis, validasinya

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| models | `tco_klausul.go` (+uji) | tabel aturan 25 jenis di KODE (AC 34), gerbang wajib-isi yang menyebut medan (AC 35), jenis ditahan (AC 36), Rp/Usd anak, batas total anak, peringatan spreading |
| repository | `tco_klausul.go` (+uji) | satu tabel `T_PROPORTIONALARRG` untuk seluruh jenis (AC 24), anak menulis NULL di sembilan kolom induk (AC 25), urut `ID ASC` (AC 31), master jenis/occupation/clause baca-saja (AC 27) |
| services | `tco_klausul.go` (+uji) | `KlausulTCO` — simpan satu transaksi, anti-dobel 409 VERBATIM (AC 30), induk beranak tidak boleh ganti jenis reasuransi |
| handlers | `tco_klausul.go` (+uji, +uji `db`) | 5 rute; nol jalur hapus |
| frontend | `PanelKlausulTahun.tsx`, `PanelJenisKlausul.tsx`, `InboxTreatyContractDescription.tsx` (+uji), label + ±50 baris korpus diuji, `api.ts` (+4) | tombol `List Description` hidup; butir menu ketiga kelompok |

### Bukti merah lebih dulu untuk aturan uang

`RpUsdAnakTCO` dimutasi ke `float64` (`strconv.ParseFloat` → kali → bagi): uji lama tidak menangkapnya; kasus bernilai
besar (`12193263122359396.42211401` × `152415.78762536`) ditambahkan dan **merah** di bawah mutasi, hijau di `apd`.
Mutasi dibuang; salinan asli dipulihkan dari scratchpad.

### Kode bersama yang disentuh (aditif, dilaporkan)

| Berkas | Perubahan |
| --- | --- |
| `handlers/penyuntikan_test.go` | +1 entri `tco_klausul.go` |
| `handlers/tco_db_test.go` (milik modul) | `sqlDBMentah()` untuk mengisi tiruan master |
| `repository/tco_jenisreasuransi.go`, `skemauji/tco_tiruan.go` (milik modul) | `masterDibacaSajaTCO` +`TREATYDESC`, `OCCUPATION`, `CLAUSE`; tiruan ketiganya + `IsiJenisKlausulTCO`, `IsiPilihanKlausulTCO` |
| `frontend/src/lib/daftarMenu.ts`, `App.tsx`, `components/Shell.test.ts` | butir menu `tco-klausul` (6 → 7) |

### Ralat / OQ

Sembilan ralat bertanggal di tiket 08 (TerritorialLimit bukan gerbang; medan CoinsPanel; exclusion empat subjenis;
spreading hanya peringatan TreatyLimitChild; `Data sudah pernah di Input` mati di Pega; pemilih jenis anak; nama kolom
master; LimitMB/Portfolio ditahan; ruang lingkup `TREATYYEARID`). **OQ baru:**

- **OQ-TCO-14** — kunci anti-dobel per jenis dan wajib-isi Object/Periode `[keputusan kami]`; Product + UW.
- **OQ-TCO-15** — daftar pilihan `ReinsTypeID` form klausul (`D_EnumerationList` tidak diekspor) → daftar tiket 02.
- **OQ-TCO-16** — nama kolom master `TREATYDESC`/`OCCUPATION`/`CLAUSE` dari properti RD; DBA.

OQ spec yang tetap terbuka: aturan LimitMB & Portfolio (AC 36), arti bisnis istilah klausul.

### Kontrak hilir

`KolomKlausulHilir` (TREATYYEAR, TREATYGROUPID, TREATYDESCID, REINSTYPEID, REINSTYPENAME, PCT, RP, USD) dan saringan induk
`PARENTREINSTYPEID` (`GetQuotaShare`, Claim Fac In) seluruhnya ditulis penulis klausul — dikunci
`TestKlausulMenulisKolomHilir`. Rp/Usd anak yang dibaca `GetLimitPLATreatyin` (Claim Prop) adalah angka turunan yang
tersimpan, sama seperti Pega.

### Angka uji

| Perintah | Hasil |
| --- | --- |
| `go vet ./...` · `go vet -tags=db ./...` · `gofmt -l .` | bersih |
| `go test ./...` | 878 lulus, 0 gagal (+24) |
| `go test -tags=db ./...` | 878 lulus, 51 dilewati (+1: `TestKlausulLingkaranPenuh`, tanpa skema uji), 0 gagal |
| `npx tsc --noEmit` · `npx vite build` | bersih |
| `npx vitest run` | 556 lulus di 46 berkas (+20, +3 berkas) |

⚠️ SQL tiket 08 belum pernah dijalankan terhadap Oracle: uji `db` dilewati karena skema uji tidak dikonfigurasi di sesi ini.

### TELEMETRI EKSEKUSI — tiket 08

| Ukuran | Nilai |
| --- | --- |
| Berkas dibaca | ±45 (±38 korpus: 25 aktivitas `SaveTreatyArr*`, `SetKirimIDDesc`, `HitungRpUsd`, `TreatyTestChildTotal_Act`, 2 RDBList simpan, harness, ±18 section form; ±7 pola kode) |
| Berkas XML korpus disensus | 25 aktivitas simpan (langkah + prasyarat), 18 section form (medan + label), harness klausul; ±50 baris label/aksi diverifikasi oleh uji |
| Perintah dijalankan | ±45 |
| Berkas ditulis / disunting | 15 baru (3.222 baris), 14 disunting (+411 −17) |
| Putaran instrumen gagal lalu diulang | 3: uji statis sempat memakai literal skema-alamat (ditangkap penjaga global); peta kolom sempat memicu positif palsu penjaga tulis-warisan (dipindah); `noUncheckedIndexedAccess` pada peta label (diubah ke `as const satisfies`) |
| Token / biaya | tidak terlihat dari dalam sesi, jadi tidak dikarang |

## Tiket 06 — security reinsurer, struktur bersih

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| models | `tco_security.go` (+uji) | gerbang security + `%Share`; penjaga nol float di tiga lapis |
| repository | `tco_security.go` (+uji) | `T_MTREATYSECURITY` berkolom bernama, kunci `ID` + `REAS_ID`, nol `TRIM`; `TestTCONolInsertPosisional` menjaga SELURUH sumber modul |
| services | `tco_security.go` (+uji) | `SecurityTCO` — satu transaksi berkunci kontrak, dobel 409, jejak |
| handlers | `tco_security.go` (+uji, +uji `db`) | 4 rute |
| frontend | `PanelSecurityReinsurer.tsx` (+uji), `SECURITY_TCO`, `api.ts` (+3) | tombol baris reinsurer `Security Reinsurer` hidup |

### Bukti merah lebih dulu

- Aturan persen (ADR-0003): nilai `%Share` 0..100 berskala 8 muat di `float64` tanpa kehilangan, jadi mutasi nilai tidak
  akan pernah merah. Penjaganya statis (`TestSecurityTanpaFloat`): `float64` disisipkan sementara ke layanan → **merah**;
  dipulihkan → hijau.
- AC 18/20: kunci UPDATE dimutasi menjadi `... AND TRIM(REAS_SECURITY) = TRIM(:3)` (perilaku warisan) → `TestSQLSecurityTCO`
  **merah** dua kali (kunci berbasis nama; `TRIM`); dipulihkan → hijau.

### Kode bersama yang disentuh (aditif, dilaporkan)

| Berkas | Perubahan |
| --- | --- |
| `handlers/penyuntikan_test.go` | +1 entri `tco_security.go` (lima penyuntikan wajib) |
| `handlers/rute_treaty_contract_out.go` (milik modul) | +1 pendaftaran rute, +3 kasus galat |
| `repository/tco_kontrak_test.go` (milik modul) | peta placeholder +4 SQL security |

### Ralat / OQ

Tujuh ralat bertanggal di tiket 06 (ganti nama memutus rujukan di Pega; DELETE berkunci nama; `Local.IsUpdate` mati;
`%Share` tidak diperiksa; `CLIENTNAME` bukan kolom; `THN_TREATYID` residu; tiga kolom `[terbuka]` NULL). **OQ baru:**

- **OQ-TCO-17** — security dobel per reinsurer ditolak dan `%Share` wajib 0..100 `[keputusan kami]` (Pega tidak menegakkan
  keduanya); Product + UW.

Blocker tiket (`TOP_ID`, `TP_TREATY`, `USER_ID`) tetap terbuka — kolom dibawa bernama, tidak dibuang.

### Kontrak hilir

Tidak ada: `MTREATYSECURITY` tidak dibaca `Claim Prop`, `Komite Claim Prop`, maupun `Claim Fac In` (pencarian penuh atas
ketiga folder korpus, nol temuan).

### Angka uji

| Perintah | Hasil |
| --- | --- |
| `go vet ./...` · `go vet -tags=db ./...` · `gofmt -l .` | bersih |
| `go test ./...` | 894 lulus, 0 gagal (+16) |
| `go test -tags=db ./...` | 894 lulus, 52 dilewati (+1: `TestSecurityLingkaranPenuh`, tanpa skema uji), 0 gagal |
| `npx tsc --noEmit` · `npx vite build` | bersih |
| `npx vitest run` | 563 lulus di 47 berkas (+7, +1 berkas) |

⚠️ SQL tiket 06 belum pernah dijalankan terhadap Oracle: uji `db` (termasuk kaskade FK) dilewati karena skema uji tidak
dikonfigurasi di sesi ini.

### TELEMETRI EKSEKUSI — tiket 06

| Ukuran | Nilai |
| --- | --- |
| Berkas dibaca | ±22 (±12 korpus: 4 aktivitas, 2 data transform, 4 RDBList, 2 RD, 2 section; ±10 pola kode) |
| Berkas XML korpus disensus | `SaveSecurityReinsurer_Act`, `ShowEditSecurityReinsurer`, `DeleteSecurityReinsurer`, `SetErrorMessageReinsurer`, `InputNewSecurityReinsurer`, `SetSecurityReinsurer`, `InsertToMTreatySecurity`, `UpdateMTreatySecurity`, `DeleteSecurityReinsurer` (SQL), `DeleteFromTreatyReinsurer_Act`, `SelectSecurityReinsurer`, `BrowseAgentReinsSOA_RD`; 20 baris label/aksi diverifikasi oleh uji |
| Perintah dijalankan | ±35 |
| Berkas ditulis / disunting | 11 baru (1.602 baris), 8 disunting (+183 −2) |
| Putaran instrumen gagal lalu diulang | 3: heredoc panjang gagal diurai (berkas ditulis lewat Write); pola sunting panel reinsurer tanpa `try {` (diperbaiki); heredoc memakan backslash di skrip turunan (disunting lewat Edit) |
| Token / biaya | tidak terlihat dari dalam sesi, jadi tidak dikarang |

## Tiket 11 — kurs USD ke IDR

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| models | `tco_kurs.go` (+uji) | tanggal/nilai master diurai sekali (AC 53), pemilihan kurs berlaku (AC 46), konversi apd (AC 51), pesan VERBATIM; aturan klausul `Berkurs`/`Konversi` |
| repository | `tco_kurs.go` (+uji) | `TREATYEXCHANGEYEARLY` dibaca saja, `QUARTER`/`IDCURRENCY` di-bind (AC 47/48) |
| services | `tco_kurs.go` (+uji) | `KursTCO`; `KlausulTCO` menurunkan `Usd` induk dan menyimpan `KURS` |
| handlers | `tco_kurs.go` (+uji, +uji `db`) | `GET /tahun/{id}/kurs`, `GET /tahun/{id}/kurs/konversi` |
| frontend | `PanelJenisKlausul.tsx`, `KURS_TCO`, `api.ts` (+2) | kurs berlaku / pesan server, `Add` nonaktif tanpa kurs, pratinjau konversi server |

### Bukti merah lebih dulu untuk aturan uang

`RpDariUsdTCO` dimutasi ke `float64` (`strconv.ParseFloat` × lalu `FormatFloat`): `TestKonversiKursTCO` **merah**
(`2004267622717146800000.00000000` alih-alih `2004267622717146774324.82837777`); dipulihkan → hijau. Nilai harapan
dihitung ulang dengan `decimal` Python presisi 80 setelah satu literal tulisan tangan terbukti salah.

### Kode bersama yang disentuh (aditif, dilaporkan)

| Berkas | Perubahan |
| --- | --- |
| `repository/matauangid.go` (Claim Life) | **dipakai, tidak diubah** — `MataUang.Pengenal("USD")` |
| `handlers/penyuntikan_test.go` | +1 entri `tco_kurs.go`; entri `tco_klausul.go` +`DenganKurs` |
| `repository/tco_jenisreasuransi.go`, `skemauji/tco_tiruan.go` (milik modul) | `masterDibacaSajaTCO` +`CURRENCY`; tiruan `TREATYEXCHANGEYEARLY` + `CURRENCY`, `IsiKursTCO`, `IsiMataUangTCO` |
| `models/tco_klausul.go`, `services/tco_klausul.go` + uji (milik modul, tiket 08) | tujuh induk: `Usd` turunan; aturan `Berkurs`/`Konversi`; uji tiket 08 diperbarui ke perilaku ini |

### Ralat / OQ

Tujuh ralat bertanggal di tiket 11 (grid XOL tak pernah menemukan kurs; `KURS` warisan selalu kosong; `Usd` induk
turunan — meralat tiket 08; skala pembagian; dua kurs berlaku; `QUARTER`; StartDate kosong). **OQ baru:**

- **OQ-TCO-18** — `KURS` diisi kurs yang dipakai, skala 8 untuk ketujuh pembagian, dan dua kurs berlaku = master rusak
  (503) `[keputusan kami]`; Finance + DBA.

OQ spec yang tetap terbuka: arti `QUARTER = '0'`.

### Kontrak hilir

`KolomKlausulHilir` tidak berubah (RP, USD tetap dua kolom — AC 49). `USD` induk yang dibaca `GetLimitPLATreatyin`
(Claim Prop) kini selalu terisi `Rp ÷ Kurs`, seperti yang Pega simpan dari form yang `Usd`-nya hanya dibaca.

### Angka uji

| Perintah | Hasil |
| --- | --- |
| `go vet ./...` · `go vet -tags=db ./...` · `gofmt -l .` | bersih |
| `go test ./...` | 908 lulus, 0 gagal (+14) |
| `go test -tags=db ./...` | 908 lulus, 53 dilewati (+1: `TestKursLingkaranPenuh`, tanpa skema uji), 0 gagal |
| `npx tsc --noEmit` · `npx vite build` | bersih |
| `npx vitest run` | 568 lulus di 47 berkas (+5) |

⚠️ SQL tiket 11 belum pernah dijalankan terhadap Oracle: uji `db` dilewati karena skema uji tidak dikonfigurasi di sesi ini.

### TELEMETRI EKSEKUSI — tiket 11

| Ukuran | Nilai |
| --- | --- |
| Berkas dibaca | ±30 (±20 korpus: `testingKurs`, `RefreshKurs`, `HitungRpUsd_depan`, `CalculateTSIExcludeTreaty`, `NewTreatyArrEpi`, `GetMasterKursList`, dua SQL simpan, `NitipKurs`, harness, tujuh form induk, SQL `CURRENCY` dua modul lain; ±10 pola kode) |
| Berkas XML korpus disensus | 14 `NewTreatyArr*` (gerbang kurs), 7 form induk (medan `Usd` hanya dibaca), 4 pemanggil `CalculateTSIExcludeTreaty`; 7 baris bukti diuji ke korpus |
| Perintah dijalankan | ±40 |
| Berkas ditulis / disunting | 9 baru (1.038 baris), 14 disunting (+368 −41) |
| Putaran instrumen gagal lalu diulang | 4: literal harapan besar salah hitung (dihitung ulang dengan `decimal`); penjaga tulis-warisan menangkap `.CURRENCY` di tiruan (diganti konstanta); pola sunting penjaga penyuntikan tidak unik (konteks diperluas); regex skrip uji mendahului pola penggantian |
| Token / biaya | tidak terlihat dari dalam sesi, jadi tidak dikarang |

## Tiket 09 — simpan atomik lintas enam tabel

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| models | `tco_status_simpan.go` (+uji) | status "1" satu-satunya sukses (AC 39); penjaga teks `JSON_KLAIM` (AC 40) |
| repository | `tco_transaksi_utuh.go` (+uji); `tco_identitas.go`; 11 berkas `tco_*.go` | transaksi utuh: pembaca sadar-transaksi, identitas sementara, penetapan identitas (+ jejak, + FK security) |
| services | `tco_simpan_utuh.go` (+uji) | `SimpanUtuhTCO` memakai ULANG kelima penulis per baris di satu transaksi (AC 37); galat menyebut baris (AC 38) |
| handlers | `tco_simpan_utuh.go` (+uji, +uji `db`) | `POST .../kontrak-utuh`, `PUT .../kontrak/{kid}/utuh` |
| frontend | `api.ts`, `simpanUtuh.test.ts` | klien + aturan status; tombol menunggu OQ-TCO-19 |

### Bukti merah lebih dulu

- Urutan identitas: mutasi yang menetapkan identitas SEBELUM baris ditulis → `TestSimpanUtuhGagalPadaKlausulKeN` **merah**
  (`tetapkan 1` pada kegagalan); dipulihkan → hijau.
- Mode sementara: mutasi yang mematikan identitas sementara → `TestIdentitasSementaraTCO` **merah**; dipulihkan → hijau.
- Penjaga `JSON_KLAIM`: teks itu disisipkan sementara ke pesan layanan → `TestTCOTanpaJSONKLAIM` **merah**; dipulihkan.

### Kode bersama yang disentuh (aditif, dilaporkan)

| Berkas | Perubahan |
| --- | --- |
| `handlers/penyuntikan_test.go` | +1 entri `tco_simpan_utuh.go` (enam penyuntikan wajib) |
| `repository/repository.go` (`Tx`), `services.DalamTransaksi` | **tidak diubah** — transaksi utuh dibawa lewat `context` di berkas modul |
| penjaga global `TestNolAlamatLayananDiKode` | tidak diubah; penangkap kode HTTP ditulis dengan menyematkan `http.ResponseWriter` supaya tidak memicu pola `http.Head` |

### Ralat / OQ

Lima ralat bertanggal di tiket 09 (tidak ada simpan utuh di Pega; tombol tunggal belum dibangun; identitas sementara;
pembaca sadar-transaksi; urutan permintaan). **OQ baru:**

- **OQ-TCO-19** — letak dan bentuk tombol simpan tunggal (layar Pega hanya punya `Save` per panel); work owner + UW.

### Kontrak hilir

Tidak berubah: simpan utuh menulis lewat penulis per baris yang sama, sehingga kolom yang dibaca hilir terisi persis
seperti simpan per panel. Identitas yang terlihat hilir selalu identitas tetap (nol `S#########T` tersisa — diuji `db`).

### Angka uji

| Perintah | Hasil |
| --- | --- |
| `go vet ./...` · `go vet -tags=db ./...` · `gofmt -l .` | bersih |
| `go test ./...` | 920 lulus, 0 gagal (+12) |
| `go test -tags=db ./...` | 920 lulus, 54 dilewati (+1: `TestSimpanUtuhAtomikLintasEnamTabel`, tanpa skema uji), 0 gagal |
| `npx tsc --noEmit` · `npx vite build` | bersih |
| `npx vitest run` | 569 lulus di 48 berkas (+1, +1 berkas) |

⚠️ Atomisitas lintas enam tabel — satu-satunya hal yang diuji tiket ini — BELUM terbukti terhadap Oracle: uji `db`
(gagal di klausul ke-3 → nol baris dan nol nomor sequence berubah; sukses → identitas tetap, nol identitas sementara
tersisa) dilewati karena skema uji tidak dikonfigurasi di sesi ini. Uji tanpa Oracle hanya membuktikan orkestrasinya.

### TELEMETRI EKSEKUSI — tiket 09

| Ukuran | Nilai |
| --- | --- |
| Berkas dibaca | ±25 (±9 korpus: enam rule Connect-SQL, `SaveTreatyContract_Act`, dua aktivitas galat; ±16 pola kode: lima layanan, repository identitas/jejak/Tx, penjaga) |
| Berkas XML korpus disensus | 7 rule simpan (baris COMMIT), 1 aktivitas simpan kontrak |
| Perintah dijalankan | ±30 |
| Berkas ditulis / disunting | 10 baru (1.185 baris), 16 disunting (+95 −24) |
| Putaran instrumen gagal lalu diulang | 3: penjaga `JSON_KLAIM` pertama ikut membaca berkas uji dan direktori (dibatasi ke produksi); konstanta aksi jejak huruf kecil; penjaga alamat-layanan global menangkap `http.Header` (diganti penyematan) |
| Token / biaya | tidak terlihat dari dalam sesi, jadi tidak dikarang |

## Tiket 10 — kaskade hapus, popup, klausul yang tetap hidup

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| repository | `tco_kaskade.go` (+uji) | dampak + kaskade kontrak (urutan Pega) dan reinsurer; saringan hitung = saringan hapus |
| services | `tco_kaskade.go` (+uji) | pratinjau; hapus satu transaksi: kunci, hitung ulang, bandingkan dengan yang dikonfirmasi, hapus, bandingkan lagi, jejak |
| handlers | `tco_kaskade.go` (+uji, +uji `db`) | 4 rute |
| frontend | `KonfirmasiHapusTCO.tsx` (+uji), `HAPUS_TCO`, `api.ts` (+4) | popup Ya/Batal di panel kontrak dan panel reinsurer |

### Bukti merah lebih dulu

- AC 44: mutasi yang menambah langkah hapus klausul ke `langkahHapusKontrakTCO` → `TestLangkahHapusKontrakTanpaKlausul`
  **merah** (langkah menyasar klausul; urutan berubah); dipulihkan → hijau. Celah yang ditemukan sebelum mutasi: penjaga
  pertama memeriksa daftar SQL tulisan tangan, bukan langkah yang dijalankan — langkahnya dipindah ke satu fungsi yang diuji.
- Angka popup: `TestHapusKontrakAngkaHarusSamaDenganPopup` menolak angka yang berubah SEBELUM hapus (nol penghapusan) dan
  angka terhapus yang berbeda (transaksi dibatalkan, nol jejak).

### Kode bersama yang disentuh (aditif, dilaporkan)

| Berkas | Perubahan |
| --- | --- |
| `handlers/penyuntikan_test.go` | +1 entri `tco_kaskade.go` |
| `frontend/src/services/api.ts` | +4 fungsi, tipe `DampakHapusTCO` — nama `DampakHapus` sudah dipakai Claim Life dan TypeScript akan MENGGABUNGKAN deklarasinya; tipe Claim Life tidak tersentuh |
| `components/ui/dasar.tsx` (`Modal`) | **dipakai, tidak diubah** |
| penjaga `TestHandlersTidakMengimporRepository` | tidak diubah; galat repository dialiaskan lewat `services.ErrKaskadeTidakUtuh` |

### Ralat / OQ

Lima ralat bertanggal di tiket 10 (Pega tanpa popup; angka popup = angka terhapus di server; klausul yang dihitung;
hapus reinsurer eksplisit; `DetailTreatyExclustion` panel rinci). **OQ baru:**

- **OQ-TCO-20** — definisi "klausul milik kontrak ini" untuk angka popup (REINSTYPEID atau PARENTREINSTYPEID tahun itu)
  `[keputusan kami]`; Product + UW.

### Kontrak hilir

Tidak berubah untuk klausul (tidak disentuh). Setelah kontrak dihapus, reinsurer/business kombinasinya hilang — pembaca
hilir `ReinsurerUntukHilir` / `BusinessUntukHilir` membaca nol baris untuk kombinasi itu, sama seperti Pega.

### Angka uji

| Perintah | Hasil |
| --- | --- |
| `go vet ./...` · `go vet -tags=db ./...` · `gofmt -l .` | bersih |
| `go test ./...` | 931 lulus, 0 gagal (+11) |
| `go test -tags=db ./...` | 931 lulus, 55 dilewati (+1: `TestKaskadeHapusKontrakKlausulTetapHidup`, tanpa skema uji), 0 gagal |
| `npx tsc --noEmit` · `npx vite build` | bersih |
| `npx vitest run` | 573 lulus di 49 berkas (+4, +1 berkas) |

⚠️ Kaskade dua arah (tiga anak hilang, klausul masih ada, bisnis ber-`TREATYYEARID` NULL ikut terhapus) BELUM terbukti
terhadap Oracle: uji `db` dilewati karena skema uji tidak dikonfigurasi di sesi ini.

### TELEMETRI EKSEKUSI — tiket 10

| Ukuran | Nilai |
| --- | --- |
| Berkas dibaca | ±20 (±9 korpus: `BrowseDeleteRowTreatyInContract`, `DeleteFromTREATYCONTRACT_SQL`, `DeleteTreatyReins_Act`, `DeleteFromTreatyReinsurer_Act`, `DetailTreatyExclustion` + `_Sec`, dua section tombol, grid exclusion; ±11 pola kode) |
| Berkas XML korpus disensus | 2 aktivitas hapus, 2 SQL kaskade, 1 flow action + section, 3 section tombol |
| Perintah dijalankan | ±30 |
| Berkas ditulis / disunting | 9 baru (1.155 baris), 10 disunting (+218 −6) |
| Putaran instrumen gagal lalu diulang | 4: nama `PerekamJejakOracle` bertabrakan dengan kode Claim Life (diganti); `DampakHapus` bergabung dengan tipe Claim Life (diganti `DampakHapusTCO`); penjaga arsitektur menolak impor repository di handler (alias services); regex uji `tahunID?` salah |
| Token / biaya | tidak terlihat dari dalam sesi, jadi tidak dikarang |

## Uji penuh + /code-review dua sumbu atas `251cb3b..HEAD` — 29-09-2026

Tinjauan dua sumbu (kebenaran/bug; kesesuaian spec, tiket, konvensi repo) atas 14 commit tiket + 2 commit dokumen.
Lima belas temuan; seluruhnya diperbaiki dalam tujuh commit `fix: temuan /code-review — …`.

| # | Berkas (saat ditinjau) | Temuan | Perbaikan |
| ---: | --- | --- | --- |
| 1 | `repository/tco_transaksi_utuh.go` | identitas sementara per transaksi berawal dari 1 → dua simpan utuh serentak menyisipkan PK sama dan saling menunggu | `081870f` — awalan 10 digit acak per transaksi |
| 2 | `services/tco_kontrak.go`, `tco_tahun.go` | anti-dobel kontrak/tahun periksa-lalu-sisip tanpa kunci | `0883587` — kunci baris tahun (kontrak); `LOCK TABLE … EXCLUSIVE` (tahun) |
| 3 | `services/tco_security.go` | teks jejak security menyimpan identitas reinsurer sementara | `081870f` — kode agen reinsurer; uji db memeriksa KETERANGAN |
| 4 | `services/tco_kaskade.go` | kombinasi dibaca sebelum kunci → kaskade memakai kombinasi basi | `0883587` — dibaca ulang sesudah kunci, lewat transaksi |
| 5 | `services/tco_kontrak.go`, `tco_tahun.go` | ganti jenis reasuransi / tahun / grup memutus anak kombinasi | `0883587` — 409 selama beranak |
| 6 | `services/tco_lampiran.go` | aksi seorang pemakai menjalankan antrean seluruh modul, jejak menyerah atas nama yang salah | `4104878` — hanya efek milik lampiran aksi itu (`PungutEfekRujukanTCO`) |
| 7 | `repository/tco_kurs.go` | satu TOIDR rusak di periode lain menggagalkan semua periode | `6a013c9` — TOIDR diurai hanya untuk baris berlaku |
| 8 | `services/tco_klausul.go` | Rp/Usd anak tidak dihitung ulang saat induk berubah | `6a013c9` — dihitung ulang di transaksi yang sama (merah dulu) |
| 9 | `models/tco_klausul.go`, `tco_kurs.go`, `tco_reinsurer.go` | konteks apd bersama ditulis (data race); konteks buatan sendiri | `6a013c9` — `utils.DecimalContext()` + penjaga statis |
| 10 | `repository/tco_kaskade.go` | kombinasi dipakai bersama kontrak tahun lain → kaskade menghapus milik kontrak lain | `0883587` + `9fe4e96` — milik bersama tidak ikut; popup menyebutnya |
| 11 | `handlers/tco_simpan_utuh.go` | 503 bermakna disamarkan jadi "(galat server)" | `34933d1` — hanya 500 yang disamarkan |
| 12 | `services/tco_tahun.go` | nama grup dari klien, ID tidak diperiksa master | `17fc31c` — diperiksa ke TREATYGROUP, nama dari master |
| 13 | `services/tco_klausul.go` | master klausul dibaca ulang per baris dalam simpan utuh | `34933d1` — kurs, jenis reasuransi, TREATYDESC sekali per permintaan |
| 14 | `PanelKontrakTahun.tsx`, `PanelReinsurerKombinasi.tsx` | sesudah 409 popup memegang angka basi | `9fe4e96` — angka dimuat ulang |
| 15 | `services/api.ts` | URL blob unduhan dicabut di tik yang sama dengan klik | `4104878` — dicabut sesudah 1 detik |

### Yang dibiarkan, dan sebabnya

- **#2 tahun**: kunci TABEL, bukan baris - tahun treaty tidak punya baris induk untuk dikunci dan indeks unik tidak dapat
  ditambahkan tanpa risiko gagal atas data migrasi yang mungkin dobel. Penulis tahun jarang; pembaca tidak terhalang.
- **#10** menyimpang dari Pega (Pega menghapus seluruh anak kombinasi) - sengaja, supaya kontrak lain tidak rusak
  (**OQ-TCO-21**, baru).
- **#13 sebagian**: pencarian master SATU baris (agen, bisnis, occupation/clause) tetap per baris - berindeks dan tidak
  mengurai riwayat; yang di-cache hanya pembacaan master UTUH.
- **#6**: `JalankanAntrean` (seluruh modul) tetap ada untuk pekerja kelak; nol pemanggil produksi sekarang.
- Tidak satu pun perbaikan terbukti terhadap Oracle: uji `db` dilewati (skema uji tidak dikonfigurasi di sesi ini).

### Angka uji per commit

`go test ./...` · `go test -tags=db ./...` (lulus + dilewati) · `npx vitest run` (lulus / berkas). Commit perbaikan dihitung
dari salinan `git archive` berisi `APP_RNM/{internal,pkg,cmd}` saja: di salinan itu 9 uji yang membaca berkas di luar
subpohon (dokumen STRUKTUR, katalog, `.env`) gagal karena lingkungan — angka di bawah adalah jumlah uji (lulus + gagal
lingkungan); di worktree nyata seluruhnya lulus.

| Commit | Isi | Go | Go `db` | vitest |
| --- | --- | ---: | --- | --- |
| `1872d26` · `1f2aab5` · `678fd25` | tiket 01 · 02 · 03 | lihat bab tiket 01–03 | | |
| `0022865` · `918eeed` | dokumen | 776 | 776 + 46 | 431 / 38 |
| `7b1db9b` | tiket 12 | 808 | 808 + 47 | 452 / 39 |
| `6737fb9` | tiket 04 | 823 | 823 + 48 | 484 / 41 |
| `5523b76` | tiket 05 | 840 | 840 + 49 | 514 / 42 |
| `19f9c79` | tiket 07 | 854 | 854 + 50 | 536 / 43 |
| `75db1f8` | tiket 08 | 878 | 878 + 51 | 556 / 46 |
| `eb314b4` | tiket 06 | 894 | 894 + 52 | 563 / 47 |
| `7dca378` | tiket 11 | 908 | 908 + 53 | 568 / 47 |
| `7989fa4` | tiket 09 | 920 | 920 + 54 | 569 / 48 |
| `57d2af6` | tiket 10 | 931 | 931 + 55 | 573 / 49 |
| `081870f` | fix #1 #3 | 931 | 931 + 55 | 573 / 49 |
| `0883587` | fix #2 #4 #5 #10 | 935 | 935 + 55 | 573 / 49 |
| `6a013c9` | fix #7 #8 #9 | 938 | 938 + 55 | 573 / 49 |
| `34933d1` | fix #11 #13 | 939 | 939 + 55 | 573 / 49 |
| `4104878` | fix #6 #15 | 940 | 940 + 55 | 573 / 49 |
| `17fc31c` | fix #12 | 941 | 941 + 55 | 573 / 49 |
| `9fe4e96` | fix #14 #10 (popup) | **941 lulus, 0 gagal** | **941 + 55 dilewati** | **574 / 49** |

Akhir (worktree nyata, `9fe4e96`): `gofmt -l` bersih · `go vet ./...` dan `-tags=db` bersih · `npx tsc --noEmit` bersih ·
`npx vite build` bersih.

### OQ

**Dibuka sesi lanjutan ini** (seluruhnya masih terbuka): OQ-TCO-08, 09 (tiket 12), 10, 11 (04), 12 (05), 13 (07), 14, 15,
16 (08), 17 (06), 18 (11), 19 (09), 20 (10), **21** (tinjauan: kombinasi bersama). **Ditutup**: tidak ada di sesi ini.
OQ spec yang tetap terbuka: aturan LimitMB & Portfolio (AC 36), arti `QUARTER = '0'`, arti bisnis istilah klausul.

### Kontrak hilir (Claim Prop, Komite Claim Prop, Claim Fac In)

- **Klausul** (`KolomKlausulHilir` + `PARENTREINSTYPEID`): ditulis seluruhnya (`TestKlausulMenulisKolomHilir`); Rp/Usd anak
  turunan tersimpan dan kini ikut dihitung ulang saat induk berubah; `USD` tujuh induk = `Rp ÷ Kurs`.
- **Reinsurer / business** (`ReinsurerUntukHilir`, `BusinessUntukHilir`): kolom VERBATIM (tiket 05/07); hapus kontrak
  tidak lagi menghapus milik kontrak lain yang berbagi kombinasi.
- **Security**: tidak dibaca modul hilir (pencarian penuh, nol temuan).
- **Tahun/grup**: `TREATYGROUPNAME` kini selalu dari master sebelum disalin ke kombinasi.

### Kode bersama yang disentuh (sesi lanjutan 1, seluruhnya aditif)

| Berkas | Perubahan |
| --- | --- |
| `handlers/penyuntikan_test.go` | entri modul: lampiran, kontrak, reinsurer, business, klausul (+kurs), security, kurs, simpan utuh, kaskade; tahun +`DenganGrup` |
| `frontend/src/lib/daftarMenu.ts`, `App.tsx`, `components/Shell.test.ts` | butir menu `tco-kontrak`, `tco-klausul` |
| `frontend/src/services/api.ts` | fungsi & tipe modul ditambahkan di ekor berkas; `DampakHapusTCO` sengaja tidak memakai nama `DampakHapus` milik Claim Life |
| `repository/matauangid.go`, `components/ui/dasar.tsx` (`Modal`), `repository/efekkeluar.go`, `repository/repository.go` (`Tx`) | **dipakai, tidak diubah** (`PungutEfekRujukanTCO` ditulis di berkas modul) |
| penjaga global (`TestNolAlamatLayananDiKode`, `TestHandlersTidakMengimporRepository`, `TestTCOWarisanHanyaDibaca`, nama jujur) | **tidak diubah**; kode modul disesuaikan agar lolos |

### TELEMETRI EKSEKUSI — uji penuh + /code-review + perbaikan

| Ukuran | Nilai |
| --- | --- |
| Tinjauan | satu putaran `/code-review` bercabang (forked) dua sumbu atas `251cb3b..HEAD`; 15 temuan, 15 diperbaiki |
| Berkas dibaca | ±25 pola kode (layanan, repository, fake uji, penjaga), 0 berkas korpus baru |
| Perintah dijalankan | ±45 (termasuk 7 penghitungan per commit lewat `git archive`) |
| Commit | 7 perbaikan; +10 uji Go, +1 uji vitest; 1 migrasi tidak ditambah |
| Putaran instrumen gagal lalu diulang | 3: `git worktree add` di scratchpad gagal (jalur berkas repositori lain terlalu panjang; tidak meninggalkan sisa, `git worktree prune` bersih) → diganti `git archive`; impor `time` salah kelompok; tipe `DampakHapus` Claim Life (tahap sebelumnya) |
| Token / biaya | tidak terlihat dari dalam sesi, jadi tidak dikarang |

## Lanjutan 2 — jawaban work owner atas OQ-TCO-08 … 21 (29-09-2026, langsung di `main`)

Titik awal `main` @ `043116f`: Go 900 uji tingkat atas lulus (1.001 termasuk sub-uji), tag `db` 57 dilewati; vitest 612;
build 81 modul. Suntingan work owner yang belum di-commit (`App.tsx`, `labels.ts`, `KelompokMenu.tsx`, `dasar.tsx`,
`.scratch/cadangan/`) dibiarkan dan tidak ikut satu commit pun.

### Kelompok 1 — OQ-TCO-11 … 18 dikonfirmasi

Delapan jawaban konfirmasi. Label kode `[keputusan kami]` / `[dugaan kuat]` → `[keputusan work owner 29-09-2026]`
(`services/tco_kontrak.go`, `repository/tco_reinsurer.go`, `models/tco_business.go`, `models/tco_klausul.go`,
`services/tco_klausul.go`, `repository/tco_klausul.go`, `models/tco_security.go`, `services/tco_security.go`,
`models/tco_kurs.go`, tiruan skema, label frontend). Konstanta `BusinessNonaktif = "0"` dikunci uji baru
`TestNilaiIsActiveKeputusanWorkOwner`. Blok "Keputusan work owner 29-09-2026" di tiket 04, 05, 06, 07, 08, 11. Register
OQ modul dibuat: `OQ-TREATY-CONTRACT-OUT.md`.

### Kelompok 2 — OQ-TCO-21 hapus seperti Pega · OQ-TCO-20 klausul dari induknya

- `repository/tco_kaskade.go`: `langkahHapusKontrakTCO` kembali empat langkah Pega tanpa percabangan "kombinasi bersama"
  (saringan bisnis tahun-saja dibuang); `Bersama` tetap dihitung untuk popup. Saringan klausul popup
  `TREATYYEARID = :1 AND ((PARENTREINSTYPEID = :2 AND REINSTYPEID = :3) OR PARENTREINSTYPEID = :4)`, sentinel di-bind.
- `services/tco_kaskade.go`: `KonfirmasiHapus.Bersama` ikut dibandingkan sebelum dan sesudah hapus; jejak menyebut kontrak
  lain terdampak. Handler: kueri `DELETE …&bersama=` wajib.
- Frontend: popup menampilkan cacah kontrak lain sebagai peringatan (`role="alert"`), `hapusKontrak` mengirim `bersama`.
- Uji perbaikan konkurensi lanjutan 1 yang mengecualikan kombinasi bersama DIGANTI `TestLangkahHapusKontrakSepertiPegaWalauBersama`;
  uji db kaskade kini membuat tahun lain berkombinasi sama dan memeriksa `"bersama":1`.

### Kelompok 3 — OQ-TCO-19 rute simpan kontrak utuh dibuang

Dibuang (tanpa pemanggil sejak jawaban *"tidak perlu"*): `handlers/tco_simpan_utuh.go` (+uji, +uji `db`),
`services/tco_simpan_utuh.go`, `services/tco_simpan_utuh_cache.go` (+uji), `models/tco_status_simpan.go`, identitas
sementara dan `TetapkanIdentitasTCO` (`repository/tco_transaksi_utuh.go` → `tco_baca_tx.go`, hanya `DenganBacaTxTCO` /
`bacaTCO` yang tersisa), kait identitas sementara di `IdentitasBerikutTCO`, entri penjaga penyuntikan, klien frontend
`simpanKontrakUtuh` / `statusSimpanSukses` (+`simpanUtuh.test.ts`). Penjaga `JSON_KLAIM` dipindah ke
`models/tco_teks_galat_test.go`. Uji baru `TestRuteSimpanUtuhDibuang` (kedua rute tidak lagi dilayani).

### Kelompok 4 — OQ-TCO-10 tanggal akhir mulai + 1 tahun kalender

`models.AkhirKontrakBawaanTCO(mulai)` (argumen tahun treaty dibuang) = mulai + 1 tahun kalender, semantik
`ADD_MONTHS(mulai, 12)`; `tahunKabisatTeks` dibuang. Pembacaan ulang `SetTanggalTreatyContract` mencetak
`pyStepsBlockName` (langkah 7–9 `//`). Uji tabel baru dibuktikan merah terhadap rumus lama (4 kasus berbeda), lalu hijau;
kasus 29 Februari juga di seam layanan (`TestKontrakAkhirBawaan`). Penyimpangan sadar bertanggal di tiket 04.

### Kelompok 5 — OQ-TCO-08/09 pekerja latar dan pelaksana penyimpanan nyata

- **Transport** `services/tco_pengirim_storage.go` — `NewPengirimBerkasHTTPTCO`: POST JSON halaman `UploadDoc` (nama medan
  verbatim), `geturl` → URL bertanda tangan → GET, `delete`; status 400/422 = permanen, 404 = berkas tidak ada, selebihnya
  dicoba ulang; galat jaringan diringkas tanpa alamat/inang/IP. Satu-satunya berkas yang memegang klien HTTP keluar.
- **Token** `services/tco_penyimpanan_nyata.go` — `NewSumberTokenStorageTCO`: `GET_TOKEN_STORAGE` ditiru (pakai ulang
  `GCP_IMAGE` beserta kedaluwarsanya, atau `RakitToken` + simpan umur 1 menit, satu transaksi). Repositori baru
  `repository/tco_storage.go`: `AppStorageTCO` (`T_FOLDER_IMAGE`, masuk `masterDibacaSajaTCO`) dan
  `TokenStorageBerlakuTCO` (token + `INPUTDATE`).
- **Pemilih** `PenyimpananLampiranTCO` / `Service.DenganPenyimpananLampiranTCO(pelaksana, garam)` — bawaan stub; `nyata`
  merakit resolver `M_LINK_SERVICE` + cache token sekali per proses + transport. Tanpa Oracle: resolver gagal terang
  (`ResolverLinkServiceOracle` bersama panik dengan db nil — tidak disentuh, dielakkan di pemilih modul).
- **Pekerja** `LampiranTahunTCO.JalankanPekerja` + `cmd/api` `jalankanPekerjaLampiranTCO`; `LayananLampiranTCO` diekspor
  dari handler supaya rute dan pekerja memakai rakitan yang sama (penjaga penyuntikan menuntut pemilihnya).
- **Config** `PELAKSANA_STORAGE` (`stub`|`nyata`, nilai lain ditolak; `nyata` tanpa `STORAGE_TOKEN_SALT` ditolak — pesan
  menyebut nama kunci, bukan nilainya) dan `TCO_PEKERJA_LAMPIRAN_INTERVAL` (durasi Go; kosong = mati).
- **Penjaga** `TestNolAlamatLayananDiKode`: peta baru `berkasKlienHTTPDisetujui` (satu baris, jumlah dikunci) mengecualikan
  transport dari pemeriksaan klien HTTP **saja**; `TestTCOLampiranTanpaAlamatLiteral` idem, dengan cek "pengecualian
  terpakai". Komentar transport sempat memuat skema-alamat literal dan tertangkap penjaga modul — diganti teks.
- **Uji baru (14 tingkat atas):** transport ×4 (bentuk `UploadDoc` + unggah ulang satu objek; ambil/periksa/buang; peta
  status + server mati + batas waktu + App kosong, galat bersih dari inang/token/garam; ujung-ke-ujung resolver → cache →
  sumber token → HTTP), sumber token ×2, pemilih ×1, galat permanen lewat antrean ×1, pekerja ×1 (mati di interval ≤ 0;
  menjalankan antrean lalu berhenti bersama ctx), config ×3, SQL repositori ×2. Semua alamat dari `httptest` saat jalan.
- **Temuan saat membangun:** handler "lambat" uji batas waktu yang tidak membaca badan membuat `httptest.Server.Close`
  menunggu selamanya — diberi kanal pelepas. `errors.Unwrap` atas galat tanpa `%w` menghasilkan nil — dibuang.
- **Pembacaan ulang XML** (`pyStepsBlockName` dicetak): `InsertGoogleStorage_Act` (Claim Life), `GetUrlGoogleStorage_Act`,
  `DeleteGoogleStorage_Act` (Treaty Contract Out) — seluruh langkah kosong kecuali `EXIT` terakhir; nol langkah `//`.
- **OQ-TCO-22 dibuka** (`Folder`/`Durasi`/`Namafile`, lihat register). Tidak ada SQL yang dijalankan ke Oracle; tidak ada
  layanan sungguhan yang dipanggil.
