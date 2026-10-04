# 07: Pembulatan di dalam loop akumulasi — kasus uji wajib

**What to build:** Bukti bahwa perhitungan premi tetap cocok **untuk daftar panjang**, bukan hanya
untuk satu nilai.

Ini tiket pengujian, dan ia berdiri sendiri karena alasan yang spesifik: di sistem lama, sebagian
pembulatan berada **di dalam loop akumulasi**, sehingga galatnya **menumpuk per iterasi**. Port yang
benar untuk satu nilai masih bisa meleset untuk daftar panjang — dan itu bentuk kesalahan yang
**paling sulit terlihat**, karena setiap nilai tunggalnya lulus.

Lingkupnya **lintas lini bisnis**: galat akumulasi tidak mengenal batas COB, jadi menguji satu COB
saja tidak membuktikan apa pun tentang yang lain.

Dua sifat sistem lama yang memperkuat kebutuhan ini: validasi sisa spreading menuntut kesamaan
**persis** (jumlah bagian = 100, jumlah premi tersebar = premi) tanpa distribusi galat pembulatan,
dan pembulatan dilakukan lebih dulu sebelum dibandingkan.

**Blocked by:** 04, 05, 06

**Status:** ready-for-human — sebagian: akumulasi MBU per mata uang (4/5 kasus NB cocok) dan total FIRE per item/lokasi diport (cabang ≥2 item lewat double: belum tercakup, A48); lini lain tanpa pembanding; validasi sisa spreading belum bertiket

- [ ] Kasus uji akumulasi untuk **tiap lini bisnis** yang punya rumus di tiket 04–06
- [ ] Setiap kasus memakai daftar yang cukup panjang untuk memunculkan penumpukan galat — bukan dua atau tiga baris
- [ ] Pembulatan terjadi **di dalam** loop, di posisi yang sama seperti sumbernya — bukan sekali di akhir
- [ ] **Rekonsiliasi eksak** pada nilai agregat, bukan hanya pada tiap elemen
- [ ] Kasus uji yang membuktikan pemindahan pembulatan ke luar loop **membuat uji gagal** — kalau tidak, ujinya tidak menguji apa pun
- [ ] Validasi sisa direproduksi tanpa toleransi: kesamaan **persis** setelah pembulatan, tanpa distribusi galat

## Comments

### 2026-10-01 — sebagian (agent)

✅ **Layering:** premi coverage = Σ premi lapisan yang masing-masing SUDAH dibulatkan di dalam loop.
`TestPremiLayering` membuktikan bedanya: 7 lapis → 0,3948; dibulatkan sekali di luar loop → 0,3946. Mutasi
"bulatkan di luar loop" tertangkap.

⏸ **PA dan MBU:** akumulasinya ada di activity total per mata uang (mis. `FillPremiMBU_FacIn` langkah
2.6.2.1.1 `Local.PremiPerCurrency += .Premium`, `AddCurrencyListPA_ACT`) yang belum diport.
⏸ **Validasi sisa spreading:** milik `services/spreading`, yang menurut indeks masih "menunggu, bukan tiket".
Rekonsiliasi agregat dengan data nyata: terbuka.

### 2026-10-01 — bagian Layering DICABUT (keputusan work owner, butir 30)

Butir ✅ Layering di atas tidak berlaku lagi: kodenya dihapus bersama tiket 06. Tiket ini kini **tanpa
satu pun kasus akumulasi yang diport** — sisa lingkupnya PA dan MBU (total per mata uang) dan validasi sisa
spreading, keduanya ⏸ seperti dicatat di atas.

### 2026-10-01 — akumulasi MBU per mata uang diport

Kode: `backend/services/premium/akumulasi.go` — `PremiPerMataUang(coverage, urutanMaster)`, langkah 2.6 jalur NB
`FillPremiMBU_FacIn` (premi coverage dibulatkan L1144 di dalam loop, lalu dijumlahkan per mata uang).

| Kriteria | Keadaan |
| --- | --- |
| Kasus uji akumulasi tiap lini berumus | ✅ MBU. PA: `[terverifikasi]` tidak ada penulis `CurrencyList.Premium` untuk PA (`AddCurrencyListPA_ACT` hanya menjumlahkan TSI) — tidak ada akumulasi premi untuk diport |
| Daftar cukup panjang | ✅ 24, 24, 38, 1, 6 coverage (93) dari kasus nyata |
| Pembulatan di dalam loop, posisi sama | ✅ lewat `Calculate` per coverage; `TestPremiPerMataUangPembulatanDiDalamLoop` (0,0001 + 0,0001 = 0,0002; di luar loop 0,0001) |
| Rekonsiliasi eksak pada agregat | ✅ 4 dari 5 kasus NB, kesamaan nilai desimal eksak (A31); ⏸ kasus #5 selisih terbuka (butir di register) |
| Memindah pembulatan ke luar loop membuat uji gagal | ✅ mutasi presisi L1144 4 → 20 (dalam dan luar) tertangkap |
| Validasi sisa spreading | ⏸ milik `services/spreading`, belum bertiket |

Uji mutasi `akumulasi.go` + presisi `mbu.go`: 7/7.

### 2026-10-01 malam — total FIRE per item dan per lokasi (agent)

Kode: `backend/services/premium/total.go` — `TotalFireLokasi`, port `SumTotalTSIPremiGross_Act` langkah 1 (`IsFire`):
premi item = Σ premi coverage (`Local.PremiCov`, penjumlahan eksak), daftar per mata uang lokasi dibangun ulang
(`Property-Remove`, 1.2) dan dijumlahkan per item, `Rate = @if(.TSI=0,0,@divide(.Premium,.TSI,20)*1000)` — seri pada
`@divide` ditolak (`ErrSeriDivide`, modenya belum terverifikasi; A37 hanya untuk `@Math.divide`).

Rekonsiliasi (`rekonsiliasi.barisTotal`, teks persis kecuali TSI per mata uang — A47): **nb-fire-1** item 1/1 + lokasi
1/1 cocok, **rnw-fire-1** 1/1 + 1/1 cocok; **edm-fire-1** belum tercakup (total EDM dari `CountGrossPremiEDM_Act`
1.1.4–1.1.5, A44). Sensus Python sebelumnya: item `TotalGrossPremi` cocok 11/11 di ketiga kasus; lokasi EDM berbeda di
sekitar desimal ke-17 dari jumlah sederhana — jalur EDM, tidak diusut.

Belum: langkah 2–6 (GOLF, ANEKA, MARINE, PA, MBU) — tanpa nilai pembanding di fixture; validasi sisa spreading
(`services/spreading`, belum bertiket).

### 2026-10-02 — ralat: penjumlahan lokasi melewati double (agent, register "Ralat ketiga")

`Local.TSI`/`Local.Premi` bertipe `double` (L293–301): cabang 1.3.4 tidak eksak di sistem lama. Komentar di atas
("penjumlahan eksak", "lokasi EDM … jalur EDM, tidak diusut") **diralat**: total EDM dihasilkan rule ini dan terulang
model double sampai digit terakhir. Kode: entri `LewatDouble` → belum tercakup (A48). Rekonsiliasi kini: item 11/11 cocok
(EDM 9, NB 1, RNW 1); lokasi NB dan RNW cocok; lokasi EDM belum tercakup. ⏸ Cabang 1.3.4 tanpa data nyata NB.

### 2026-10-02 — A48 diputuskan: tetap eksak (butir 63)

Rumusan yang tepat: item ke-2 dst. **dibulatkan ke `double`** (`Local.Premi`/`Local.TSI`, L293–301) lalu diakumulasi
`Decimal` (`Local.PremiCov`, L305–307) — bukan "menjumlahkan lewat double". Bukti model: 2/2 lokasi dari satu kasus
(`edm-fire-1` akar + `OldData`). Keputusan: port **tetap eksak**, lokasi ≥2 item bermata uang sama **belum tercakup**;
paritas `double` hanya lewat ADR bila kelak diwajibkan.

