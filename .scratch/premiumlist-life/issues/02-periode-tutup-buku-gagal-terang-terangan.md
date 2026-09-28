# 02: Periode tutup buku dibaca dari `POOLDATA.TANGGAL_CLOSING` — dan gagal terang-terangan bila kosong

**Status:** sebagian — periode produksi belum ditampilkan di layar Summary dan Detail sebelum menyimpan

**Blocked by:** **00 (skema tujuh tabel — PREFACTOR)**

## Hasil & nilai pengguna

Sebagai **tim akuntansi**, saya ingin setiap transaksi premium list dibukukan ke **periode produksi
yang benar** menurut tanggal tutup buku yang berlaku, dan saya ingin sistem **berhenti dan berteriak**
bila tanggal tutup buku tidak dapat dibaca — supaya kesalahan periode ketahuan saat itu juga, bukan
saat tutup buku bulan berikutnya ketika angkanya sudah terlanjur salah. *(User story 39 di spec)*

## Area codebase

`internal/repository` (pembacaan `TANGGAL_CLOSING`), `internal/services` (kebijakan periode produksi;
jam yang dapat dikendalikan), `internal/handlers` (pesan kesalahan yang menyebut tabel sumber),
`frontend/` (menampilkan periode produksi yang dipilih pada layar premium list).

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path |
| --- | --- | --- |
| `SubmitPremiumList_Act` | `ASM-FW-GISFW-WORK-LIFE` / `SUBMITPREMIUMLIST_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/SubmitPremiumList_Act.xml` (214.154 byte, tersimpan `20260122T072213`, **17 langkah**) |
| `GETTanggalClosing_SQL` | `ASM-FW-GISFW-INT-POLICYJSON` / `RNM!GETTANGGALCLOSING_SQL` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/GETTanggalClosing_SQL.xml` — `SELECT * FROM POOLDATA.TANGGAL_CLOSING` |

`[terverifikasi]` Rantai di `SubmitPremiumList_Act`:

| Step | Baris | Isi |
| ---: | ---: | --- |
| 2 | 715 | `RDB-List` "Get Tanggal Closing" → `GETTanggalClosing_SQL` |
| 3 | 892 | "Set Tanggal Closing" — `Local.TglProd = TglProd.pxResults(1).TANGGAL` (918–919) |
| 3 | **966** | ⚠️ `Local.TglProd = @if(Local.TglProd=="",25,Local.TglProd)` — **fallback diam** |
| 4 | **1211** | `Local.NextMonth = @if(@toDecimal(Local.currentdate)>Local.TglProd, @toDecimal(Local.CurrentMonth)+1, @toDecimal(Local.CurrentMonth))` |
| ~~15~~ | 3387 | ~~"Kalau acc di atas tanggal 25 akan masuk produksi bulan berikutnya"~~ — **REMARK** (`<pyStepsBlockName>//`, baris 3398); memuat precondition `@toDecimal(Local.currentdate)>Local.TglProd` (3491) |

`[terverifikasi]` Pergeseran periode: `@CurrentDate("yyyy","Asia/Jakarta") + Local.NextMonth +
"01T050000.000 GMT"` — **tanggal 1 bulan berikutnya, `05:00 GMT` = 12:00 WIB**.

⚠️ `[terverifikasi]` **Dua versi hidup berdampingan.** `SubmitPremiumList_Act` (`20260122`) membaca
ambang dari tabel; `InsertJsonPolisLife_Act` (`ASM-FW-GISFW-WORK-LIFE` / `INSERTJSONPOLISLIFE_ACT`,
`20260728` — **lebih baru**) masih menanam konstanta: step 4 dengan precondition
`@toDecimal(Local.currentdate)>25` (baris 1170). `[keputusan work owner]` **ikuti yang dari DB.**

⛔ *Ralat 28-09-2026 (sensus remark GILIRAN-12):* `SubmitPremiumList_Act` langkah 2–4 membaca
`TANGGAL_CLOSING`, tetapi hasilnya hanya mengalir ke `TempGenerate.CARI1/2` (tak dibaca rule mana
pun) dan ke langkah 15 yang **ter-remark** (`//` b3398; b3491 miliknya). Dua aturan yang **hidup**:
nomor PL digulir `PROC_GENERATE_SEQUENCE_NUMBER` (membaca tabel itu sendiri), dan `ProdDateTime`
digeser `InsertJsonPolisLife_Act` langkah 4 (`>25` tertanam, b1170). Keputusan "ikuti yang dari DB"
tetap; penerapannya pada `ProdDateTime` (yang di Pega memakai 25) dikonfirmasi ulang — **OQ-PL-13**.

`[data DBA]` `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` juga menggulir periode lewat
`POOLDATA.TANGGAL_CLOSING` — jadi tabel ini adalah sumber tunggal aturan periode di kedua sisi.

## ADR terkait

**ADR-0006** (penomoran lewat stored procedure — periode adalah masukannya), **ADR-0005** (flag
lingkungan), **ADR-0007** (jejak audit).

## Acceptance criteria

- [x] Tanggal tutup buku dibaca dari `POOLDATA.TANGGAL_CLOSING` **setiap kali dibutuhkan**, bukan
      dari konstanta, bukan dari cache yang tidak pernah kedaluwarsa. *(AC 7 spec)* — bukti: `repository/polis_periode.go:Tanggal`, `repository/penomor.go:HariClosing` (dibaca tiap pemanggilan, nol cache)
- [x] Bila `POOLDATA.TANGGAL_CLOSING` **kosong atau tidak terbaca**, transaksi **ditolak** dengan
      pesan yang **menyebut tabel sumbernya**. Sistem **tidak** memakai nilai pengganti apa pun.
      *(AC 8 spec; `[keputusan work owner]`)* — bukti: `models/polis_periode.go:PeriodeProduksi`; uji `TestTabelKosongDITOLAK_BukanDiamDiamPakai25`
- [x] Tidak ada konstanta `25` — maupun angka ambang lain — di kode produksi. Test yang mencari
      literal ambang di lapisan services/repository **gagal** bila ada. — bukti: uji `TestNolAmbangTutupBukuTertanam`, `TestPenjagaAmbangMasihMenggigit`
- [x] Transaksi pada tanggal **setelah** tanggal tutup buku memperoleh periode **tanggal 1 bulan
      berikutnya**, pukul `05:00 GMT` (12:00 WIB). *(AC 9 spec)* — bukti: uji `TestPeriodeSelaluTanggalSatuJamLimaGMT`
- [x] Transaksi pada tanggal **sama dengan** tanggal tutup buku **tetap** di periode berjalan —
      perbandingannya `>`, bukan `>=`. *(baris 1211 dan 3491 keduanya memakai `>`)* — bukti: uji `TestPerbandinganLebihBesarBukanLebihBesarSama`
- [x] Pergantian tahun tertangani: tutup buku Desember menggeser ke Januari tahun berikutnya, bukan
      ke "bulan 13". — bukti: uji `TestPergantianTahunBukanBulanTigaBelas`
- [x] Jam dapat dikendalikan dari test — aturan periode diuji tanpa menunggu tanggal nyata. — bukti: `services/polis_periode.go:DenganJam`; `models/polis_periode.go:PeriodeProduksi` menerima jam dari pemanggil
- [ ] Periode produksi yang terpilih **terlihat pengguna** sebelum ia menyimpan, bukan hanya
      tersimpan diam-diam. — belum: periode tampil hanya di layar keputusan (`InputOffer.tsx`); layar Summary (`Submit`) dan tombol terbitkan nomor di Detail tidak menampilkannya sebelum menyimpan

## Blocker

**Tidak ada.**

## Catatan — perbedaan sengaja dari Pega

| Perilaku Pega | Sistem baru |
| --- | --- |
| `@if(Local.TglProd=="",25,Local.TglProd)` — diam-diam memakai `25` | **gagal terang-terangan**, transaksi ditolak |

Alasannya: fallback diam membukukan transaksi ke periode yang salah **tanpa meninggalkan jejak**, dan
kekeliruannya baru terlihat saat tutup buku. Ini kegagalan tersembunyi, bukan ketahanan.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```
