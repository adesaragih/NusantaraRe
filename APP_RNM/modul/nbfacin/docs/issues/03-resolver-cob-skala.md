# 03: Resolver lini bisnis → skala rasio, dan aturan pembagi komposit

**What to build:** Setiap rasio yang masuk perhitungan mendapat skalanya dari **satu resolver
terpusat** yang memetakan lini bisnis (COB) → per mille atau persen. Tidak ada tempat kedua yang bisa
menyimpang.

Properti rate yang sama memakai **dua satuan berbeda menurut lini bisnis**. Menyeragamkannya
menggeser premi **satu ordo besaran 10** pada separuh portofolio — jadi bahayanya bukan salah memilih
satuan, melainkan menyamakannya.

**Aturan penguraian pembagi komposit — inti tiket ini.** Pembagi gabungan pada satu operasi pembagian
adalah **hasil kali** sumbangan tiap faktor bersatuan, bukan satu satuan tunggal. Satuan rate dibaca
dari **faktor yang menempel padanya saja**, tidak pernah dari pembagi total:

```
@Math.divide((.TSI * .Rate * ProRatePercent), 100000, 4)
        100000  =  1000 (rate ber-‰)  ×  100 (ProRatePercent ber-%)
```

Pembagi `100000` karena itu **menegaskan** rate = per mille — bukan membantahnya. Siapa pun yang
membacanya sebagai "satuan lebih kecil dari per mille" akan salah satu ordo besaran.

Resolver mengikuti **rumus**, bukan label layar: tiga label layar terbukti bertentangan dengan
rumusnya, dan rumus yang menang di ketiganya.

**Blocked by:** 01

**Status:** ready-for-human — diport 01-10-2026, menunggu tinjauan work owner

- [ ] Peta skala lengkap: **PA · Layering · FIRE = per mille**; **MBU · ANEKA · BONDING · GOLF · MARINE CARGO = persen**
- [ ] Aturan pembagi komposit diterapkan: sumbangan per faktor, bukan pembagi total
- [ ] Kasus uji membuktikan pembagi `100000` menghasilkan **per mille**, bukan satuan lain
- [ ] **Lini bisnis yang belum ada di peta → gagal keras**, bukan skala default
- [ ] Skala **tidak** dibaca dari tabel konfigurasi yang dapat diubah tanpa deployment — ia fakta struktural korpus, bukan parameter bisnis
- [ ] Resolver adalah **satu-satunya** pengisi skala; tidak ada pemanggil yang menyetel skala sendiri


## Comments

### 2026-10-01 — implementasi (agent)

Kode: `APP_RNM/modul/nbfacin/backend/services/premium/resolver.go` (+ `resolver_test.go`); `premium.go`
memakainya. Keputusan: `../KEPUTUSAN-30-09-2026.md` butir 26 dan catatan bukti per lini.

| Kriteria | Keadaan |
| --- | --- |
| Peta skala lengkap 8 lini | ✅ `SatuanRate`: PA · Layering · FIRE = ‰; MBU · ANEKA · BONDING · GOLF · MARINE CARGO = % |
| Aturan pembagi komposit | ✅ `Satuan.Pembagi()` per faktor; `Calculate` PA L713 memakai rate × ProRatePercent |
| Uji pembagi 100000 → per mille | ✅ `TestPembagiKompositPerMil` + rekonsiliasi kasus PA metode 1 |
| Lini tak dikenal → gagal keras | ✅ panic, juga lewat `Calculate` |
| Tidak dari tabel konfigurasi | ✅ peta di kode |
| Resolver satu-satunya pengisi skala | ✅ `ProRatePercent` sengaja di luar resolver — satuannya milik medan, bukan lini (K-018) |

Lini yang ada di peta tetapi belum punya rumus premi → `ErrBentukBelumDiport` (bukan panic). 7 mutasi
atas resolver dan pemakainya, seluruhnya tertangkap.

⚠️ Label diturunkan: hanya PA, MBU, Layering yang buktinya dikutip K-018; lima lini lain bersandar pada
keputusan K-018 dengan penguat `[dugaan]`.

### 2026-10-01 — tindak lanjut `/code-review` (agent)

- Kriteria "resolver satu-satunya pengisi skala" kini **harfiah**: `SatuanProRata()` di `resolver.go`,
  `rasio` hanya dibangun lewat `rasioRate`/`rasioProRata`, dijaga uji AST
  `TestRasioHanyaDibangunDiResolver` (butir 27). Penjaga itu pada percobaan pertama langsung menangkap
  `kompilasi_gagal_rasio.go`, berkas wajib-gagal-kompilasi yang sengaja memakai `rasio{}`; ia masuk
  daftar-izin eksplisit beserta alasannya.
- Bukti di komentar kini menyebut berkas dan `class / nama`; label `[keputusan work owner]` untuk lima
  lini yang tidak dikutip buktinya di K-018; klaim L1146 dipisah dari "cocok eksak" (L1146 tidak dipakai).
- `TestPembagiKompositPerMil` (tautologis) dibuang — pembuktian 100000 → per mille tetap lewat
  `Calculate` (`TestPremiPAContohHitung` "prorata pembagi komposit" + rekonsiliasi metode 1).
- Satu tabel `dataSatuan` untuk simbol + pembagi; helper uji `harusPanic`; komentar kepala berkas;
  spec Modul 2 diberi catatan sinkron.

Uji mutasi: 9 mutasi, seluruhnya tertangkap.

### 2026-10-01 — Layering dicabut dari peta (keputusan work owner, butir 30)

Layering tidak dipakai lagi di sistem baru: peta `SatuanRate` kini **7 lini** (PA · FIRE = ‰; MBU · ANEKA ·
BONDING · GOLF · MARINE CARGO = %). `"Layering"` kini panic seperti lini di luar peta
(`TestSatuanRateLiniTakDikenalPanic`). Teks K-018 diberi amandemen, tidak dihapus
(`../00-KEPUTUSAN-WORK-OWNER.md`).
