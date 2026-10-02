# F04: Derivasi Fac Retro untuk tiga lini bisnis sisanya

**What to build:** Derivasi retrosesi yang sudah hidup untuk Fire (F02) diperluas ke **Aneka/Golf**,
**Marine Cargo/MBU**, dan **PA/Travel**, sehingga seluruh lini bisnis dapat menurunkan Fac Out.

`[terverifikasi]` `Activity\SetDataFacOut_Act` memanggil **empat** sub-aktivitas per lini —
`SetDataFacOutFire_Act` (F02) plus `SetDataFacOutAnekaGolf_Act`, `SetDataFacOutCargoMBU_Act`,
`SetDataFacOutPATravel_Act` — masing-masing bergerbang predikat lini bisnisnya sendiri.

⚠️ **Tiap lini menulis ke koleksi anak yang berbeda.** `[terverifikasi]` `FacRetroList` punya **enam**
koleksi anak, bukan satu: `LocationList` · `CargoList` · `PersonList` · `VehicleList` ·
`Property.RiskLocation.AnekaList` · `Property.RiskLocation.OccupationList(n).AnekaList`. Menyamakan
keenamnya ke satu daftar akan kehilangan data.

`[terverifikasi]` **Penjumlah PA berdiri sendiri.** `Activity\SumFacOutPA_Act` bergerbang
`.TreatyType=="10015"` dan `IsPA`, menyetel `.PremiumSpreaded` dan `.TSISpreaded`. Ia **bukan**
bagian dari `SetDataFacOutPATravel_Act` dan tidak boleh dilebur ke dalamnya.

⚠️ **Dua dari empat sub-aktivitas berbeda isi di Endorsement.** `[terverifikasi]` kontrak 23 tag:
`SetDataFacOutFire_Act` dan `SetDataFacOutAnekaGolf_Act` **berbeda** NB↔EDM; `SetDataFacOutCargoMBU_Act`
dan `SetDataFacOutPATravel_Act` **identik**. Perbedaannya ditangani **F11**, bukan di sini.

**Asal (Pega).** `Activity\SetDataFacOutAnekaGolf_Act` · `SetDataFacOutCargoMBU_Act` ·
`SetDataFacOutPATravel_Act` · `SumFacOutPA_Act` · induk `SetDataFacOut_Act`

**Keputusan.** K-053 · K-054 · K-018 (skala rasio per lini) · K-010/K-012 · `CLAUDE.md` §4.6

**Blocked by:** F02

**Status:** blocked

- [ ] Keempat lini menurunkan Fac Out lewat **satu jalur derivasi yang sama**, bercabang di gerbang lini
- [ ] **Enam koleksi anak** `FacRetroList` terisi sesuai lininya, tidak dileburkan
- [ ] Penjumlah PA tetap **langkah tersendiri**, bergerbang penanda ketat `"10015"` **dan** lini PA
- [ ] Skala rasio mengikuti **resolver K-018**, bukan angka hafalan per lini
- [ ] Lini yang tidak ada di peta resolver → **`panic`** (`CLAUDE.md` §4.5)
- [ ] Nilai uang bertipe `Money`; rasio bertipe `Ratio`
- [ ] Perbedaan Endorsement pada dua sub-aktivitas **tidak ditangani di sini** — komentar menunjuk F11
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
