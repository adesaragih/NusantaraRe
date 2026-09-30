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
