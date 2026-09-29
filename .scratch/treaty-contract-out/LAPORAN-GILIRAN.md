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
