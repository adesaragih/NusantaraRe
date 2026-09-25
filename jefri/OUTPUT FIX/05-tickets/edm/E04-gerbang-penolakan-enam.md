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

**Status:** blocked

- [ ] Keenam gerbang menolak pada kondisi yang benar **dan hanya** pada kondisi itu
- [ ] Gerbang pembayaran **hanya** berlaku untuk jenis endorsement pembatalan
- [ ] **Empat klep** membatalkan penolakan yang benar, masing-masing dengan penanda jenisnya sendiri
- [ ] **Gerbang kedua dan keenam tetap menolak** meski tabel pengecualian terisi — keduanya tanpa klep
- [ ] Klep kosong atau tabel tak terjangkau → **penolakan berlaku penuh**
- [ ] Bypass "endorsement RI slip" melewati gerbang yang benar — diport apa adanya
- [ ] Rule pengecekan pembatal **diimplementasikan**; **isi tabelnya di luar lingkup** dan bukan blocker
- [ ] Rujukan laporan bersifat **runtime** — tidak dapat divalidasi saat kompilasi; ketiadaannya **bukan `panic`** karena tabelnya punya pemilik (K-006)
- [ ] Tiap gerbang menyebut rule Pega asalnya dalam komentar (§4.6)
