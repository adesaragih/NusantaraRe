# 11: Terima & tampilkan hasil keputusan Komite

> ⚠️ **Cakupan diubah 2026-09-15.** Judul lama: *"Kontrak Komite — jalur balik hasil keputusan"*.
> Tiket ini **tidak lagi menulis `STS_REJECT`**. `[terverifikasi]` Penulisannya milik **Komite Claim
> Life**, bukan Claim — Life: `Komite Claim Life/Activity/KomitePostAdjustment.xml`
> (`ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`) yang menulis ke
> dua tingkat baris. Implementasinya ada di **tiket Komite 05**.
> Tiket ini kini **membaca dan menampilkan** hasil itu, lalu melanjutkan siklus klaim.

**Status:** claimed

**Blocked by:** 10 (kontrak Komite — penyerahan) · **Komite 05** (jalur balik — penulisan
`STS_REJECT`)

## Hasil & nilai pengguna

Sebagai **ReasLifeSPV**, saya melihat hasil keputusan Komite pada baris adjustment saya — sehingga
saya tahu apakah baris itu diaksep atau ditolak — dan bila ditolak, saya dapat menambah baris
adjustment baru untuk mengajukan ulang dengan angka yang diperbaiki.
*(User story 20, 21, 22 di spec)*

Inilah yang membuat **klaim tidak terminal**: penolakan menghasilkan putaran berikutnya, bukan akhir.

## Area codebase

`internal/services` (pembacaan hasil; pembuatan baris lanjutan; perhitungan status klaim turunan),
`internal/handlers` (endpoint status klaim + tambah baris), `frontend/` (tampilan hasil dan riwayat
putaran).

**Tidak** menulis `STS_REJECT` — itu milik Komite 05.

## Rule Pega sumber

| Rule | Identitas | Peran di tiket ini |
| --- | --- | --- |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` **penulis** hasil — 6 `Property-Set` bernilai `1`, 2 bernilai `2`, digerbangi `KomiteCount == KomiteLoop`. Tiket ini **membacanya**, tidak menjalankannya |
| `Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` | `[terverifikasi]` **satu rule bersama** — hash ternormalisasi `c50bfd9a12` identik di kedua modul |
| `Claim Life/Activity/SetIndexAdjustmentList.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `SETINDEXADJUSTMENTLIST` / `RULE-OBJ-ACTIVITY` | pewarisan 8 kolom ke baris lanjutan |

## ADR terkait

**ADR-0001** (Kontrak 2 — jalur balik; `AcceptStatus` dipetakan ke `STS_REJECT` **di batas**, tidak
disimpan sebagai status kedua), **ADR-0011** (terminal per baris; klaim tidak terminal; revisi =
baris baru), **ADR-0007** (jejak audit).

## Acceptance criteria

- [x] Hasil keputusan Komite **terbaca** pada baris `AdjustmentList` yang diserahkan — Aksep atau
      Ditolak. *(AC 4 spec Claim Life)*
- [ ] Klaim **tetap** dapat menerima baris adjustment baru setelah penolakan Komite.
- [x] Baris baru yang ditambahkan setelah penolakan berstatus Outstanding dan **mewarisi delapan
      kolom** dari baris pertama **tanpa** mewarisi status. *(AC 5 spec Claim Life)*
- [ ] `PremiumListDetail` dan header klaim **mencerminkan** baris terakhir setelah hasil diterapkan.
      *(AC 7 spec Claim Life)*
- [x] Status klaim "selesai" dihitung sebagai keadaan **turunan** dari kumpulan baris — bukan kolom
      tersimpan.
- [x] `AcceptStatus` **tidak** disimpan sebagai status kedua di konteks ini. *(**ADR-0001**)*
- [x] Hasil keputusan **hanya diterapkan** ketika putaran Komite sudah mencapai **tingkat terakhir**;
      hasil dari tingkat antara **tidak** mengubah status baris mana pun di konteks ini.
      *(AC 6 spec)* ⚠️ **Penegakannya milik Komite Claim Life tiket 05** — tiket ini hanya wajib
      **tidak menerapkan lebih awal**. Dicatat agar AC 6 punya jejak pemilik, bukan tampak terlewat.
- [x] Tiket ini **tidak** menulis `STS_REJECT` — diverifikasi dengan tidak adanya jalur tulis status
      baris di konteks Claim — Life.

## Catatan — mengapa cakupan diubah

`[terverifikasi]` Penulisan `STS_REJECT` dilakukan rule di modul **`Komite Claim Life`** dengan class
`ASM-FW-GCNMFW-Work-KomiteLife`. Menempatkannya di tiket Claim — Life akan membuat **dua tiket
`ready-for-agent` sama-sama mengklaim penulisan yang sama** — dua agent dapat mengimplementasikannya
berdua.

`[keputusan work owner 2026-09-15]` Kepemilikan ditetapkan: **Komite menulis, Claim Life membaca.**

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Pembacaan ulang XML — 27 September 2026

Jendela `D:\XML\RNM_BRD\Komite Claim Life\**\*.xml`. Satu rule, dibaca sebagai **penulis** —
tiket ini membacanya, tidak menjalankannya.

| Rule (path) | Baris pecahan | Yang dipastikan |
| --- | --- | --- |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` | **2220 · 2280** · **3198 · 3261** · **3476 · 3539** | ENAM `Property-Set` menulis `.STS_REJECT = 1` — tiga pasang, masing-masing ke `AdjustmentList` **dan** `PremiumListDetail` |
| idem | **6303 · 6349** | DUA `Property-Set` menulis `.STS_REJECT = 2`, juga berpasangan |
| idem | **5695 · 8119 · 8648 · 8887** | keempat gerbang menuntut `pyWorkPage.KomiteCount == pyWorkPage.KomiteLoop` — hasil hanya berlaku di tingkat TERAKHIR |
| idem | 6172 | penolakan menyetel `KomiteCount = KomiteLoop` — tangga ditutup seketika |
| idem | 9028 | persetujuan tingkat antara menyetel `KomiteCount = KomiteCount + 1` |

**Sensus dua arah.** Cara 1 — `<PropertiesName>` memuat `STS_REJECT` → **8**, terbagi 6 bernilai
`1` dan 2 bernilai `2`. Cara 2 — baris mana pun menyebut `STS_REJECT` → **20** (termasuk definisi
rule). Klaim tiket *"6 `Property-Set` bernilai 1, 2 bernilai 2"* **sahih** `[terverifikasi]`.

⭐ **Yang sensus ungkap dan tiket belum sebut:** tiap penulisan menyentuh **dua** tingkat sekaligus
— baris adjustment **dan** pesertanya — tetapi **bukan header**. Pencerminan ke header datang dari
`serviceInsertArasapasClaimLife_act`, yang sudah dicatat tiket 04.

**Pertanyaan yang XML tidak jawab:** apakah peserta ikut berubah status ketika baris **Outstanding
baru** lahir. Korpus hanya menulis status peserta saat KEPUTUSAN, tidak saat baris baru dibuat.
`[terbuka — work owner]`

## Implementasi — 27 September 2026 (tiket 11)

**Status: `claimed`** — **6 dari 8 AC tertutup.** Titik tetap `c6f6691`.

| Angka | Nilai | Cara 1 | Cara 2 | Sepakat? | Label |
| --- | --- | --- | --- | --- | --- |
| AC tertutup | **6 dari 8** | `grep -c '^- \[x\]'` → 6 | `grep -cE '^- \[[x ]\]'` → 8 dikurangi `grep -c '^- \[ \]'` → 2 | ✅ | `[terverifikasi]` |
| test Go | **222 PASS · 0 FAIL · 34 SKIP** *(dari 214/32)* | cacah awalan `--- PASS`/`--- FAIL`/`--- SKIP` dari `go test -tags=db ./internal/... -v` → **256** | `grep -rhoE '^func Test[A-Za-z0-9_]+' internal/ --include=*_test.go \| wc -l` → **256** | ✅ | `[terverifikasi]` |
| test JS | **13 PASS** *(dari 11)* | `npm test` → `Tests 13 passed` | `grep -c '  it(' src/services/api.test.ts` → 13 | ✅ | `[terverifikasi]` |
| modul frontend | **88** | `npm run build` → `88 modules transformed` | — | **belum punya data** | `[dugaan]` |

| Berkas | Isi |
| --- | --- |
| `services/hasilkomite.go` *(baru)* | `HasilKomite`, gerbang tingkat terakhir, `BarisLanjutan`, `Putaran` |
| `services/hasilkomite_test.go` · `_db_test.go` · `_statik_test.go` *(baru)* | 7 + 2 + 3 kasus |
| `repository/pohonklaim.go` | `SisipkanBaris` + `sisipBarisAdjustment` — INSERT baris adjustment kini **satu** salinan |
| `handlers/putaran.go` *(baru)* + `handlers.go` | satu rute |
| `frontend/` `api.ts` · `KlaimLife.tsx` · `api.test.ts` | `tambahPutaran`, `bolehPutaranBaru`, tombol "Putaran berikutnya" |

⭐ **Pemanggil produksi pertama `TambahBaris`.** Ia lahir di tiket 03 sebagai fungsi murni tanpa
pemanggil; `BarisLanjutan` memakainya di sini, sehingga aturan pewarisan yang diuji sejak tiket 03
akhirnya benar-benar berjalan di jalur pengguna.

**Penjaga, tiap-tiap dibuktikan gagal lalu dipulihkan:**

| Mutasi | Hasil |
| --- | --- |
| kolom `ACCEPT_STATUS` disisipkan ke DDL | merah, kolomnya disebut |
| `WarisiKolom` mulai menyalin `KomiteID` | merah — *"KomiteID bocor ke putaran baru"* |
| status lanjutan diwarisi, bukan Outstanding | merah |
| menulis lewat array peserta yang DIBAGI | merah — *"baris asal berubah"* |
| baris lanjutan diberi keputusan Aksep | merah, dua kali: perilaku **dan** penugasan literal |

⛔ **Dua mutasi semula DIAM, dan diam bukan lulus.** Menghapus `baru.KomiteID = ""` tidak membuat
apa pun merah — sebab ketiga pembersih eksplisit itu **no-op hari ini**: `TambahBaris` berangkat
dari baris kosong. Dan `p.Baris = append(...)` pada parameter nilai tidak menyentuh pemanggil
ketika kapasitasnya pas. Keduanya diulang dengan mutasi yang benar — `WarisiKolom` yang diubah,
dan penulisan lewat array yang dibagi — dan barulah merah. Sifat no-op itu **ditulis di kodenya**,
supaya tidak ada yang mengira ia sudah terbukti.

⛔ **Penjaga tiket 07 menuntut pendaftaran, dan itu benar.** `hasilkomite.go` menulis status
Outstanding pada baris baru, jadi ia penulis status — dan harus terdaftar beserta gerbangnya.
Didaftarkan dengan `WajibPeran(pelaku, PeranSimpanOutstanding)`.

⛔ **AC terakhir tidak dipura-purakan.** Bunyinya *"tiket ini tidak menulis `STS_REJECT`"*, dan
harfiahnya **tidak benar**: `BarisLanjutan` menulis `STS_REJECT = 0`. Yang dijaga dibuat lebih
tajam dan lebih berguna — konteks ini tidak pernah menulis **KEPUTUSAN** (`1` maupun `2`).

⭐ **Akibat lintas tiket yang dikunci:** medan bank **tidak** termasuk kedelapan kolom warisan,
sehingga baris lanjutan **tidak lolos gerbang rekening tiket 10** sampai rekeningnya diisi.
Rekening putaran lama tidak otomatis menjadi rekening putaran baru — dan itu dikunci supaya tidak
"diperbaiki" diam-diam kelak.

### ⛔ AC yang belum tertutup — 2

⛔ **SATU AKAR, LIMA TIKET — dinyatakan di sini karena di sinilah ia tertangkap.**
Keempat jalur tulis modul ini berakhir di `jejak.Rekam`: `statusbaris.go:306`, `tahap.go:149`,
`komite.go:403`, `hasilkomite.go:277`. Pemasangan `Jejak` produksi **nol** — `DenganJejak` hanya
dipanggil test. Akibatnya **setiap jalur tulis Claim Life menjawab HTTP 501 hari ini**: menolak
baris (tiket 05), memindah tahap (08), menyerahkan ke Komite (10), dan membuka putaran (11).
Semuanya menunggu butir **am** disahkan work owner.

Ini bukan temuan tiket 11 saja. AC serupa di tiket 04, 05, 08, dan 10 dibaca di bawah standar yang
sama: yang terpasang **bentuknya**, dan isinya menunggu tabel jejak. Saya tidak menyunting AC tiket
yang sudah di-commit; yang dilakukan adalah **menyatakannya**, sekali, di tempat ia terbukti.

| AC | Sebab |
| --- | --- |
| *"Klaim tetap dapat menerima baris adjustment baru setelah penolakan"* | ⛔ **dicabut centangnya sesudah review.** Aturannya berjalan dan diuji (`BarisLanjutan`, `TestKlaimTidakTerminalSetelahPenolakan`), tetapi **jalur simpannya** berakhir di `jejak.Rekam` yang selalu gagal → transaksi batal, rute menjawab 501. Menunggu **am** |
| *"`PremiumListDetail` dan header klaim mencerminkan baris terakhir"* | header **tertutup** — `CerminkanHeader` dipanggil dan diuji di test db. `PremiumListDetail` **tidak**: korpus hanya menulis status peserta saat KEPUTUSAN, bukan saat baris Outstanding baru lahir, dan apakah peserta ikut berubah adalah `[terbuka — work owner]`. Tidak dikarang |

### Hasil /code-review — titik tetap `c6f6691`

Kedua review **membangun elakannya, mengompilasinya, dan menjalankannya**. Tujuh elakan terbukti
hijau terhadap penjaga yang baru saja saya tulis.

**Spec — 1 cacat nyata + 3 centang tidak sah + 1 lubang penjaga.**

| Temuan | Tindakan |
| --- | --- |
| ⛔ **CACAT: `IS_CHECK` tidak dipulihkan.** `Tolak` mencabutnya; `Tambah` tidak memasangnya kembali; `TambahBaris` menulisnya ke `salin` yang saya buang. Korpus menggerbangi **6** prasyarat pada `.IsCheck = true` → baris lanjutan lahir lalu **tidak akan pernah diambil Komite**. Tujuan tiket ini persis itu | **diperbaiki** — `PasangPenandaDipilih` di transaksi yang sama, plus assertion di test db |
| AC *"klaim tetap dapat menerima baris baru"* dicentang padahal `jejak.Rekam` selalu gagal → rute menjawab 501 | **dicabut centangnya**, dan akar modul-lebarnya dinyatakan |
| `CerminkanHeader` menulis `ACCEPTED_NO = ""` tanpa syarat, mengabaikan peserta lain | **diterima, dinyatakan** — mengikuti kontrak tiket 04; arti "baris terakhir" pada klaim berpeserta banyak sudah `[terbuka]` sejak tiket 04 |
| penjaga keputusan dapat dielakkan `models.BarisAdjustment{KodeStatus: "1"}` (titik dua) | **diganti pemeriksaan SAAT JALAN** di `repository.PeriksaBarisBaru` |

**Standards — 3 elakan penjaga terbukti + 4 pelanggaran sensus.**

| Temuan | Tindakan |
| --- | --- |
| ⛔ **`` Go menganggap garis bawah sebagai huruf** → `KOMITE_ACCEPT_STATUS` lolos dari `ACCEPT_?STATUS`; pemindaian DDL tidak rekursif dan `.SQL` huruf besar terlewat | pola tanpa ``, telusur rekursif, sufiks tak peka huruf besar |
| ⛔ **konteks ini dapat menimpa keputusan baris LAMA**: menyisipkan `PerbaruiStatusBaris(...)` ke transaksi `Tambah` menulis `STS_REJECT = 1` pada baris lama, seluruh suite **hijau** | `hasilkomite.go` kini dilarang memanggil pengubah status baris lama sama sekali |
| ⛔ **`barisTerakhirPeserta` membocorkan penunjuk** ke array pemanggil meski parameternya nilai; test mengaku "dilindungi parameter nilai" padahal hanya memeriksa 2 medan | kembalikan **salinan**; test membandingkan **seluruh** baris asal dengan `reflect.DeepEqual` |
| §4a ×2 + kontradiksi + perintah audit hilang | **diralat**, lihat bab *Ralat sensus* |
| ⭐ **Speculative Generality**: `SumberHasilKomite`, `BacaHasil`, `PeriksaHasilFinal`, `KeputusanKomite`, `HasilKomite` — **nol** pemanggil produksi | **seluruh jalur baca DIBUANG** |
| `sisipBarisAdjustment` diperiksa byte demi byte terhadap INSERT aslinya | **identik** — 18 kolom, urutan sama, pemetaan NULL sama |

⛔ **Jebakan tiket 07 saya ulangi, sesudah menulis bahwa saya tidak akan.** Di tiket 12 saya menulis
*"pemanggil produksi ada sejak lahir — kekeliruan tiket 07 tidak diulang di sini"*, lalu di tiket 11
melahirkan tujuh nama tanpa satu pun pemanggil. Yang AC minta sudah berjalan tanpa semuanya:
keputusan ditulis ke `STS_REJECT` baris, dibaca `AmbilBaris`, ditampilkan sebagai kata oleh layar —
jalur yang ada sejak tiket 01. Rantai join hanya menambah rincian yang tidak satu pun AC minta.

⭐ **Penegakan "tidak diterapkan lebih awal" kini STRUKTURAL.** Konteks ini tidak dapat menerapkan
keputusan pada tingkat mana pun: `PeriksaBarisBaru` menolak baris baru berkeputusan, dan berkas ini
dilarang memanggil pengubah status baris lama. Struktur yang tidak memungkinkan lebih kuat
daripada gerbang yang memeriksa.

**lanjut dari sini:** tiket 11 selesai. Berikutnya tiket 13 — migrasi data penuh, yang
**menuntut persetujuan manusia** sebelum dijalankan.

