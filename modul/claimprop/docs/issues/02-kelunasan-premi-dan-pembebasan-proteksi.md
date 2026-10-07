# 02: Pemeriksaan kelunasan premi dan pembebasannya lewat data proteksi

**Status:** sebagian 07-10-2026 — gerbang premi lunas dibangun; View Status Payment Premi nonaktif (OQ-CP-03); bacaan ARASAPAS tidak terlihat dari DEV (OQ-CP-11) — RALAT 07-10-2026 (semula `ready-for-agent`)
**Blocked by:** 00 (PREFACTOR) · 01 (registrasi klaim dan nomor polis)
**Menutup:** AC 113 · 114 · 115 · 116 · 117 *(5 AC)* — US 5–6

## Hasil & nilai pengguna

Klaim tidak diproses di atas polis yang preminya belum dibayar. Claim Admin melihat sebabnya secara
eksplisit, dan kasus yang sudah disetujui manajemen tetap dapat berjalan lewat data proteksi —
tanpa perlu mematikan pemeriksaannya.

## Area codebase

Validasi akseptasi · integrasi baca ke schema `arasapas` · aturan pembebasan proteksi.

## Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `RDBList/CekLunasPremi_Sql.xml` | — | menjumlahkan nilai bertanda atas `arasapas.invoice` dan `detail_invoice`; ⚠️ pasangan tabelnya **tanpa prefiks schema** |
| `Activity/CekPremiLunas_Act.xml` | 2 | mengisi parameter untuk prefix `CLMP-` |
| `Activity/CekPremiLunas_Act.xml` | 3 | mengisi parameter untuk prefix `CLM-` — ⚠️ **tidak ada cabang untuk `CLMNP-`** |
| `Activity/CekPremiLunas_Act.xml` | 4 | menjalankan pencarian **tanpa gerbang** |
| `Activity/CekPremiLunas_Act.xml` | 5 | ⚠️ menambal koma-ke-titik atas hasil jumlah premi |
| `Activity/CekPremiLunas_Act.xml` | 6 | menandai belum lunas dan menyiapkan pesan bila hasilnya lebih besar dari nol |
| `Activity/CekPremiLunas_Act.xml` | 7 | gerbang induk — pembebasan **hanya dievaluasi ketika premi belum lunas** |
| `Activity/CekPremiLunas_Act.xml` | 7.1 | memanggil `RDBList/CekProteksiKlaim.xml` |
| `Activity/CekPremiLunas_Act.xml` | 7.2 | mengembalikan status menjadi lunas bila proteksi ditemukan |
| `Activity/CekPremiLunas_Act.xml` | 8 | menampilkan pesan bila status tetap belum lunas |
| `RDBList/CekProteksiKlaim.xml` | — | menyaring `pooldata.openproteksi_edm`; ⚠️ dua kode di-hardcode di dalam SQL |
| `Activity/GetDtlPaymentPremi_act.xml` | — | jalur rincian pembayaran premi |

## ADR terkait

**ADR-0003** (uang non-float — nilai premi tidak boleh melewati floating point).

## Acceptance criteria

- [ ] `[terverifikasi]` Premi polis diperiksa lunas sebelum akseptasi; bila belum lunas, akseptasi ditahan dan pesan sebabnya ditampilkan. ⚠️ **"Lunas" berarti jumlah bertanda ≤ 0, bukan = 0** *(AC 113 spec)*
- [ ] ⚠️ Kedua tabel premi diprefiks schema eksplisit. **Alasan menyimpang:** di Pega `arasapas.invoice` diprefiks sementara pasangannya `detail_invoice` tidak — dalam satu perintah yang sama *(AC 114 spec)*
- [ ] ⚠️ Perilaku pemeriksaan premi untuk prefix `CLMNP-` dibawa apa adanya dan **ditandai**, bukan diam-diam dilengkapi. **Catatan `[terbuka]`:** belum dapat digolongkan — apakah Go menambah cabangnya belum dijawab work owner, jadi apakah ini menyimpang pun belum dapat diputuskan. Yang terbukti: Pega mengisi parameter hanya untuk `CLMP-` dan `CLM-`, sedangkan pencariannya berjalan tanpa gerbang — sehingga untuk `CLMNP-` pencarian berjalan dengan parameter yang tidak pernah diisi *(AC 115 spec)*
- [ ] `[terverifikasi]` Pemeriksaan premi dapat dibebaskan lewat data proteksi; pembebasan **hanya dievaluasi ketika premi belum lunas** *(AC 116 spec)*
- [ ] ⚠️ Dua kode proteksi menjadi **nilai bernama**, bukan literal di dalam perintah SQL. **Alasan menyimpang:** di Pega keduanya dipaku ke dalam teks SQL sehingga tidak dapat diubah tanpa menyentuh rule *(AC 117 spec)*

## Catatan `[terbuka]` ringan — TIDAK memblokir

⚠️ `[terbuka]` Apakah prefix `CLMNP-` memang dikecualikan dari pemeriksaan premi, atau cabangnya
tertinggal? Pemilik: **work owner**. Bawa apa adanya; perilaku yang ditiru adalah perilaku yang
berjalan.

## Perintah verifikasi

```
jalankan test "premi belum lunas -> akseptasi ditahan dengan pesan sebabnya"
jalankan test "jumlah bertanda = 0 -> dianggap LUNAS"
jalankan test "jumlah bertanda < 0 -> dianggap LUNAS"
jalankan test "premi belum lunas + proteksi terbuka -> akseptasi LOLOS"
jalankan test "premi sudah lunas -> proteksi TIDAK dibaca"
periksa kedua tabel premi berprefiks schema               -> tidak ada yang telanjang
cari literal kode proteksi di dalam SQL                   -> nihil
```
