# 02: Periode tutup buku dibaca dari `POOLDATA.TANGGAL_CLOSING` — dan gagal terang-terangan bila kosong

**Status:** ready-for-agent

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

`[data DBA]` `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` juga menggulir periode lewat
`POOLDATA.TANGGAL_CLOSING` — jadi tabel ini adalah sumber tunggal aturan periode di kedua sisi.

## ADR terkait

**ADR-0006** (penomoran lewat stored procedure — periode adalah masukannya), **ADR-0005** (flag
lingkungan), **ADR-0007** (jejak audit).

## Acceptance criteria

- [ ] Tanggal tutup buku dibaca dari `POOLDATA.TANGGAL_CLOSING` **setiap kali dibutuhkan**, bukan
      dari konstanta, bukan dari cache yang tidak pernah kedaluwarsa. *(AC 7 spec)*
- [ ] Bila `POOLDATA.TANGGAL_CLOSING` **kosong atau tidak terbaca**, transaksi **ditolak** dengan
      pesan yang **menyebut tabel sumbernya**. Sistem **tidak** memakai nilai pengganti apa pun.
      *(AC 8 spec; `[keputusan work owner]`)*
- [ ] Tidak ada konstanta `25` — maupun angka ambang lain — di kode produksi. Test yang mencari
      literal ambang di lapisan services/repository **gagal** bila ada.
- [ ] Transaksi pada tanggal **setelah** tanggal tutup buku memperoleh periode **tanggal 1 bulan
      berikutnya**, pukul `05:00 GMT` (12:00 WIB). *(AC 9 spec)*
- [ ] Transaksi pada tanggal **sama dengan** tanggal tutup buku **tetap** di periode berjalan —
      perbandingannya `>`, bukan `>=`. *(baris 1211 dan 3491 keduanya memakai `>`)*
- [ ] Pergantian tahun tertangani: tutup buku Desember menggeser ke Januari tahun berikutnya, bukan
      ke "bulan 13".
- [ ] Jam dapat dikendalikan dari test — aturan periode diuji tanpa menunggu tanggal nyata.
- [ ] Periode produksi yang terpilih **terlihat pengguna** sebelum ia menyimpan, bukan hanya
      tersimpan diam-diam.

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
