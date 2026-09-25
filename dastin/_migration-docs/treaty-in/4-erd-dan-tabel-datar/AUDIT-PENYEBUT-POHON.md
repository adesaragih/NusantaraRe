# Audit penyebut — setiap artefak yang angkanya datang dari pohon

**Tanggal:** 24 September 2026
**Sebab:** L-8 — `TITIK-BUTA-POHON.md`

> **SETIAP KESIMPULAN YANG PERNAH DIHITUNG DARI POHON ITU DIHITUNG ATAS SEMESTA YANG KURANG.**
>
> Maka yang diperiksa bukan korban yang terlihat, melainkan **daftar penyebut**. Satu baris per
> artefak: apakah angkanya berubah bila dihitung dari sapuan mandiri.

**"Kemungkinan besar aman" bukan jawaban.** Setiap baris di bawah ini punya pemeriksaan yang
dijalankan, bukan penalaran.

---

## 0. Dua hasil yang menentukan seluruh tabel

### 0.1 Kepala kontrak TIDAK tersentuh — nol dari 138

| | |
|---|---:|
| Properti kelas akar `ASM-FW-GISFW-Int-TREATY_IN` | **138** |
| Yang tidak ada di pohon | **0** |

Kelas akar tertangkap **utuh**. Itu masuk akal dan sekaligus melegakan: titik butanya adalah aturan
ber-*applies-to* **kelas anak**, dan aturan semacam itu tidak menulis properti kepala.

**Akibatnya langsung:** `SPEC-MODEL-DATA.md` **§10.1 (`KONTRAK`, 8 atribut) dan §10.2
(`VERSI_KONTRAK`, 47 atribut) tidak tersentuh sama sekali.** Yang harus dikerjakan ulang hanya §10.3
dan §10.4 — keduanya entitas anak.

### 0.2 Tidak ada ENTITAS bisnis baru — tetapi ada sembilan WADAH yang tidak pernah terlihat

Pertanyaan yang belum pernah ditanyakan: *adakah **entitas** yang tidak pernah terlihat?*
`BreakDownSprdList` membuktikan jawabannya bisa "ada". Diperiksa dua tingkat:

**Tingkat properti** — dari 377 nama properti yang tidak ada di pohon, **44 berbentuk wadah**
(pernah muncul diikuti `(` atau `.`, jadi ia halaman atau daftar, bukan skalar). Yang jatuh di 14
kelas penyandang 27 entitas: **sembilan**.

| Wadah | Kelas | Putusan |
|---|---|---|
| `BreakDownSprdList` · `BreakDownSprdListXOL` | `LimitsSpreading` | **`RINCIAN_PENYEBARAN`** — entitasnya sudah ada namanya |
| `RNMSpreadedList` · `RNMSpreadedListRI` | `LimitsDetail` | **turunan** — kembaran proporsional dari 10 daftar `…XOL` yang §12.3 sudah nyatakan turunan |
| `MDPMinList` | `Limits` | **daftar uang** — minimum deposit premium; paket uang, bukan entitas |
| `ReserveList` | `LimitsDetail` | **daftar uang** — cadangan premi; berpasangan `PremiumReservePct` yang sudah ada di §10.4 |
| `Brokerage` | `LimitsDetail` | **daftar uang** |
| `Achievement` | `LimitsDetail` | **GEL-3** |

**Tingkat kelas** — 56 kelas punya properti; **17 muncul di pohon, 39 tidak**. Ke-39 disisir satu
per satu, dan **tidak satu pun entitas bisnis baru di dalam gelombang ini**:

| Golongan | Kelas | Kedudukan |
|---|---|---|
| sudah dikenal sebagai batas atau gelombang lain | `Int-TREATYINDETAIL` · `Int-TREATYINOFFER` · `Data-TreatyInAchievement` · `Data-TreatyInRetroShare` · `Int-treaty_in_edm` | proyeksi datar, batas yang digambar, GEL-3, GEL-2, tabel addendum |
| master di luar skema | `Int-BUSINESS` · `Int-TREATYREINSURER` · `Int-AGENT` · `Int-CURRENCY` · `Int-TREATYGROUP` · `Int-REINSURANCETYPE` · `Int-CLIENT` · dan 8 lainnya | ADR-0023, ADR-0041 |
| perkakas dan penawaran | `Data-Search` · `Data-Quotation` · `Data-OfferTreatyIn*` · `Work` · `Work-NB` · `Int-Param` · `Int-OFFERJSON` · `Int-M_LINK_SERVICE` · `Int-T_STORAGE_IMAGE` | bukan model |
| **bentuk konversi, bukan entitas** | **`Data-TreatyInBusinessLimit` (18)** · `Data-TreatyInLimitCash` · `Data-TreatyInCessionLimit` · `Data-TreatyInCommision` · `Data-PolicyTreatyIn` | dipakai **hanya** oleh `TreatyInMappingDataconvert.xml` dan `…Prop.xml` — bentuk sumber konversi, bukan bentuk simpan |
| sudah DITUNDA | `Data-TreatyInProfitCommision` (`ME`, `ProfitComm`, `YDCF`) | §6, *profit commission* |

> **Daftar entitas selamat.** `STRUKTUR-DATA.md` tidak bertambah dan tidak berkurang. Yang berubah
> **atribut** dan **paket uang**, bukan entitas.

---

## 1. Daftar penyebut — satu baris per artefak

| Artefak | Penyebutnya | Berubah? | Dasar putusan |
|---|---|:--:|---|
| `PETA-TELUSUR-JSON.md` | **667 jalur** — 485 / 166 / 12 / 4 | **YA** | semestanya sapuan berawalan `TreatyIn.`; kurang **sekurang-kurangnya 74** jalur di dalam lingkup. §5-nya mengakui batasnya tetapi tidak pernah mengukurnya |
| `SPEC-MODEL-DATA.md` §1 | 667 diadili · 272 cermin · 213 masukan · 187 daun · 113+54 turunan · 11 dibuang · 4 ditunda | **YA** | semesta yang sama |
| `SPEC-MODEL-DATA.md` §2.3 | daftar entitas | **tidak** | §0.2 — nol entitas bisnis baru |
| `SPEC-MODEL-DATA.md` §3 | 187 jalur masukan per entitas | **YA, sebagian** | `LAYER`, `DETAIL_PROPORSIONAL`, `BAGIAN`, `PENYEBARAN` kurang. **Kepala tidak** — §0.1 |
| `SPEC-MODEL-DATA.md` §4 | 119 jalur daun turunan (54+65) | **YA** | 19 `Total…`/`Sum…` dan `RNMSpreadedList`(RI) masuk turunan, belum terhitung |
| `SPEC-MODEL-DATA.md` §5 | 12 dibuang | **tidak, untuk sekarang** | belum ada yang baru diadili sebagai dibuang; bertambah bila §10 memutuskan begitu |
| `SPEC-MODEL-DATA.md` §6 | 4 ditunda | **tidak** | `Data-TreatyInProfitCommision` memang sudah di sana |
| `SPEC-MODEL-DATA.md` §7 | **43 paket uang dari 65 keluarga** | **YA** | `MDPMinList`, `ReserveList`, `Brokerage`, `GrossPremiumMinList` tidak pernah terhitung. Ditambah lubang yang sudah dicatat: 16/22/5 **tidak pernah dituangkan sebagai daftar** |
| `SPEC-MODEL-DATA.md` §10.1, §10.2 | 8 dan 47 atribut | **tidak** | §0.1 — nol dari 138 |
| `SPEC-MODEL-DATA.md` §10.3 | `LAYER` 13 atribut | **YA** | 16 hilang, 13 menuntut keputusan |
| `SPEC-MODEL-DATA.md` §10.4 | `DETAIL_PROPORSIONAL` 9 atribut | **YA** | 38 hilang, 17 menuntut keputusan |
| `SPEC-MODEL-DATA.md` §12.4 | lima entitas dilengkapi | **tidak** | sumbernya **sudah** inventaris kelas — sumber mandiri yang sama |
| `INVENTARIS-STRUKTUR-DATA.md` | atribut per benda | **YA** | §0-nya menyatakan sumbernya "setiap jalur properti yang dibaca atau ditulis aturan mana pun" — semesta yang sama |
| `STRUKTUR-DATA.md` | 19+1+7+4+2+1 entitas | **tidak jumlahnya** | tetapi **arti dan kunci alami `RINCIAN_PENYEBARAN` berubah** — lihat `TITIK-BUTA-POHON.md` §6 |
| `struktur-treatyin-lama.md` | **985 simpul · 17 kelas · 5 tak dideklarasikan** | **YA** | 985 adalah yang **terjangkau dari akar**. 56 kelas punya properti; **39 tidak pernah muncul di pohon** |
| `peta-nama-tabel-treatyin.tsv` · `PETA-NAMA-TABEL-TREATYIN.md` | **44 tabel datar** | **YA** | tidak ada tabel untuk `BreakDownSprdList`/`…XOL` |
| `ERD-STRUKTUR-TREATYIN.html` | 44 kotak · 39 relasi | **YA** | ikut peta nama tabel |
| `TABEL-DATAR.md` | kolom per tabel datar | **YA, tidak langsung** | ia turunan dari entitas; berubah sejauh atribut entitasnya berubah |
| `ERD.md` · `ERD-TREATY-MASUK.md` · `ERD-TREATY-MASUK.html` | entitas rancangan | **tidak** | disusun dari DDL Oracle, bukan dari pohon |
| `SPEC-INVARIAN.md` | **63 invarian** | **tidak jumlahnya** | disusun dari model, bukan dari pohon. **Tetapi INV-47 dan INV-50 menyebut sumbu yang tidak ada** — §2 di bawah |
| `SEAM-ADJUSTMENT.md` | ruas seam | **tidak** | disusun dari model |
| `KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md` G1c | **"kesembilan belas properti kelas addendum"** | **tidak** | **diverifikasi ulang:** `Int-treaty_in_edm` punya tepat **19** properti. Angkanya benar, dan sumbernya memang inventaris kelas |
| `5-tiket/DAFTAR-PEKERJAAN.md` | 58 kemampuan · **27 entitas** | **tidak** | §0.2 — entitasnya tidak bertambah |
| `5-tiket/KEPUTUSAN-PEMBAGIAN-TIKET.md` | gerbang | **tidak** | tidak memuat angka berpenyebut pohon |

---

## 2. Yang paling berbahaya di daftar ini, dan kenapa

**`PETA-TELUSUR-JSON.md` dipakai sebagai BUKTI KELENGKAPAN.** Kepalanya berbunyi *"Ini bukan
lampiran. Ini bukti kelengkapan."* dan kriteria terimanya *"tidak ada satu jalur pun di luar empat
kategori — kategori kelima adalah tempat hal-hal bersembunyi."*

Kriteria itu **terpenuhi atas 667 jalur**, dan 667 bukan seluruhnya.

> **Bukti yang lengkap atas himpunan yang bolong tetap bolong, dan ia lebih berbahaya daripada
> tidak punya bukti sama sekali — karena ia menghentikan orang mencari.**

Yang harus ditulis di kepala berkas itu bukan koreksi angka, melainkan **batas semestanya**: empat
golongan berlaku atas jalur yang terjangkau dari halaman akar, dan jalur yang ditulis aturan
ber-*applies-to* kelas anak **tidak ada di dalamnya**.

`PETA-TELUSUR-JSON.md` §5 sudah menyebut lubang metodenya dan bahkan menamai dua properti yang
diketahui hilang. Yang tidak pernah dilakukan: **mengukurnya**. Selisih antara "ada lubang, kami
tahu" dan "lubangnya 74 di dalam lingkup, ini daftarnya" adalah selisih antara catatan dan bukti.

---

## 3. Apa yang dikerjakan atas daftar ini

| Baris | Perlakuan |
|---|---|
| §10.3, §10.4 | **dikerjakan ulang sekarang** bersama §10 untuk 23 entitas — bukan utang |
| §1, §3, §4, §7 `SPEC-MODEL-DATA.md` | angkanya diperbarui **saat §10 selesai**, karena §10-lah yang menghasilkan angka barunya |
| `PETA-TELUSUR-JSON.md` | **batas semestanya ditulis di kepala berkas**, dan jalur yang baru terbaca ikut diadili saat §10 mengadili atributnya |
| `peta-nama-tabel-treatyin.tsv`, ERD struktur | dijalankan ulang **setelah** §10, karena keduanya keluaran perkakas, bukan tulisan tangan |
| `struktur-treatyin-lama.md` | §8 batas bukti diberi **angka**, dan 985 disebut sebagai "terjangkau dari akar" |
| INV-47, INV-50 | diperbaiki **sekarang** — §2 `TITIK-BUTA-POHON.md` §6 dan pemeriksaan sumbu pihak |

---

## 4. Cara menjalankan ulang audit ini

```
python alat/sapu-properti-per-kelas.py
```

Lalu tiga pemeriksaan yang menghasilkan §0 di atas, seluruhnya atas kedua CSV keluarannya:

1. **kepala aman?** — hitung properti kelas akar `Int-TREATY_IN` dan berapa yang hilang. Nol berarti
   §10.1 dan §10.2 tidak tersentuh.
2. **wadah yang tidak terlihat?** — dari properti yang hilang, ambil yang pernah muncul diikuti `(`
   atau `.` di ekspor. Itu calon entitas.
3. **kelas yang tidak terlihat?** — kelas di sapuan dikurangi kelas di pohon. Sisakan yang
   ber-awalan `Data-TreatyIn`, dan periksa **siapa pemakainya**: kelas yang hanya dipakai aktivitas
   konversi adalah bentuk sumber, bukan bentuk simpan.
