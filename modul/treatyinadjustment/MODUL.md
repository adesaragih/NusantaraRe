# Modul `treatyinadjustment` — Treaty In Adjustment

⚠️ **Terdaftar 1 Oktober 2026 — LAPISAN SKEMA tiket `01` dan `05`.** `backend/modul.go` ada,
`inti/backend/daftar/modul_treatyinadjustment_gen.go` bangkit, dan slot menu `974` menyalakan
`M_NAV_MENU.DIMIGRASI`.

⛔ **Modul ini TIDAK membuat satu tabel pun**, dan itu bukan kelalaian. Model datanya **satu** dengan
Treaty In — pemisahan 25-09-2026 memindahkan papan tiketnya, bukan `SPEC-MODEL-DATA.md`-nya. Migrasi
`440` dan `441` **mengubah `VERSI_KONTRAK` milik modul `treatyin`**: menambah
`ID_VERSI_KONTRAK_DASAR` (tiket `01`) dan melonggarkan `NOMOR_URUT_VERSI` menjadi boleh kosong
(tiket `05`, langkah PERLUAS).

⛔ **Tiket `01` BELUM selesai sebagai tiket**: keempat penolakannya — versi penyesuaian tanpa dasar,
versi pertama berdasar, dasar berkeadaan `DITOLAK`, dasar berkeadaan `DIBATALKAN` — menuntut jalur
simpan yang belum ada. Yang selesai artefak skemanya.

⛔ **Tiket `02` dan `03` masih punya sisa pekerjaan SKEMA**, bukan hanya penegakan: dua nilai
`SIFAT_MATERIAL_ADDENDUM` (`02`) dan **bawaan** `TANGGAL_BERLAKU_ADDENDUM` (`03`) belum dinyatakan di
mana pun. Uraiannya di [`docs/KEPUTUSAN-TIKET-02-03.md`](docs/KEPUTUSAN-TIKET-02-03.md).

ℹ️ **Layar modul ini ada karena mendaftarkan modul menuntutnya** `[keputusan work owner 01-10-2026]` —
baca-saja, nol alur karangan (~~`L-4`~~ — **dicabut 3 Oktober 2026**; lihat bab di bawah. Layar modul ini belum dibangun dari ekspornya, dan itu kini pekerjaan yang tertunda, bukan larangan).

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `treatyinadjustment` |
| Folder korpus | `Treaty In Adjustment` |
| GROUPMENU | `TREATY` |
| Pemilik | `@PEMILIK-TREATYINADJUSTMENT` |
| Status | dimigrasi |
| Rentang migrasi | `440-479` |
| Slot menu | `974-975` |
| Prefix rute API | `/api/treaty-in-adjustment` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-TREATY-IN-ADJUSTMENT.md`, `KEPUTUSAN-TIKET-02-03.md`, `LAYAR-ADJUSTMENT.md`, `PERTANYAAN-TERBUKA-LAYAR-ADJUSTMENT.md`, `issues/` (tiket `01`, `02`, `03`, `05`) |
| `backend/` | `modul.go` (`Pendaftaran()`), `models/`, `repository/`, `services/`, `handlers/`, `migrations/` |
| `frontend/` | `menu.ts`, `rute.tsx`, `api.ts`, `labels.ts`, `labelsPenyesuaian.ts`, `pages/` (`PenyesuaianKontrak.tsx`, `RantaiVersi.tsx`, `LampiranKontrak.tsx`), `komponen/SisiPenyesuaian.tsx` |

## Migrasi

Rentang `440-479` (tabel R2, urut hulu ke hilir: migrasi modul hilir yang merujuk tabel modul hulu
selalu berjalan sesudahnya). Slot menu `974-975` hanya menyalakan `DIMIGRASI` baris modul ini (satu `UPDATE`,
nol `INSERT` — menu datar 30-09-2026) saat modul mendapat layar pertamanya, di folder
`backend/migrations/` modul ini sendiri — bentuk SQL-nya di `APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`
bab 6. Nomor selalu tiga digit.

## Pernyataan untuk penjaga

### Kaskade ON DELETE CASCADE

⚠️ **RALAT 4 Oktober 2026.** Bab ini pernah berbunyi *"nol kaskade di modul ini"*, dan tabelnya
sengaja kosong. Itu benar selama modul ini tidak punya tabel sendiri. **Keputusan pemilik proses
4 Oktober 2026** ([`treatyin/docs/KEPUTUSAN-PENYELARASAN-REPO.md` §19](../treatyin/docs/KEPUTUSAN-PENYELARASAN-REPO.md))
memindahkan `NILAI_SELISIH` dan `NILAI_SEBELUM_PRO_RATE` ke sini, dan keduanya membawa kaskadenya.

Modul ini punya **tiga** kunci asing. Yang pertama miliknya sejak awal:
`VERSI_KONTRAK.ID_VERSI_KONTRAK_DASAR` → `VERSI_KONTRAK` (migrasi `440`).
`4-erd-dan-tabel-datar/ERD.md` §2.2 menyatakannya **[hapus: tolak]** — *"menghapus versi dasar
akan membuat seluruh baris selisih kehilangan artinya"* — sehingga ia berdiri tanpa klausa
`ON DELETE`.

| Awalan berkas | Relasi |
| --- | --- |
| `442_` | `VERSI_KONTRAK` → `NILAI_SELISIH` · `VERSI_KONTRAK` → `NILAI_SEBELUM_PRO_RATE` *(dua)* — `ERD.md` §2.6, ERD HTML baris 36 dan 37 |

⛔ **Daftar-IZIN, bukan pencabutan.** Berkas di luar `442_` tetap ditolak oleh
`TestKaskadeHanyaPadaRelasiTerdaftar` (penjaga inti) dan `TestKaskadeHanyaPadaBerkasTerdaftar`
(penjaga modul ini). Yang berubah isinya, bukan aturannya.

⛔ **`TestNolTabelBaru` juga menjadi daftar-izin**, bukan dicabut: kedua tabel itu boleh, tabel
ketiga tetap ditolak. Riwayat penolakan pertamanya — migrasi ini pernah ditulis sebagai `442_`
lalu dipindah ke `treatyin/426_` — tercatat di kepala `treatyin/.../426_nilai_selisih.sql`.

⚠️ **DITAGIH ke pemilik spec:** `treaty-in/SPEC-MODEL-DATA.md` masih menyatakan model datanya
**SATU**, dan pesan penjaga lama mengutipnya. Pernyataan itu kini tidak lagi benar seluruhnya.

⚠️ **Relasi kedua modul ini akan punya belum dapat dibuat.** `ERD.md` §2.6 menyatakan dua relasi
yang melintasi sekat, dan **keduanya** menuju `NILAI_SELISIH`:

```
VERSI_KONTRAK             1--o<  NILAI_SELISIH   [hapus: ikut hapus]   SEKAT
BESARAN_DAPAT_DISESUAIKAN 1--<   NILAI_SELISIH   [hapus: tolak]        SEKAT
```

`NILAI_SELISIH` **tidak ada di `ddl-usulan/` maupun `KAMUS-KOLOM.md`** — kolomnya belum diputuskan
siapa pun. Tiket `06` dan `13` bersandar padanya. Begitu tabelnya lahir, yang pertama masuk tabel di
atas sebagai kaskade.

## ⛔ `L-4` DICABUT — 3 Oktober 2026

`L-4` berbunyi: ***"tidak ada spesifikasi layar di mana pun."*** Atas dasar itu papan melarang
menyentuh UI, dan dua ronde menolak membangun layar.

**Pernyataan itu salah, dan sudah salah sejak awal.** Spesifikasinya ada di ekspor Pega 2026-09,
dan cacahnya dapat diperiksa dalam satu perintah:

| Modul | `Section/` | `Harness/` |
|---|---:|---:|
| Treaty In | **50** | **3** |
| Treaty In Adjustment | **68** | **6** |

```bash
ls "D:/XML_NURE/Treaty In/Section"/*.xml | wc -l            # 50
ls "D:/XML_NURE/Treaty In/Harness"/*.xml | wc -l            # 3
ls "D:/XML_NURE/Treaty In Adjustment/Section"/*.xml | wc -l # 68
ls "D:/XML_NURE/Treaty In Adjustment/Harness"/*.xml | wc -l # 6
```

Yang tidak ada bukan spesifikasinya, melainkan **pembacaannya**. Di dalam berkas itu tata letaknya
tertulis lengkap: `<pyValue>` properti yang diikat sebuah kontrol, `<pyLabelFor>` medan milik
sebuah label, `<pyPropertyTarget>` properti yang diisi pemilih, `<pyGridProps>` bentuk tabelnya.

> ### Larangan yang TETAP berlaku
>
> **Jangan mengarang layar.** Yang berubah hanya alasannya: dulu tidak boleh karena tidak ada
> acuannya, kini tidak boleh karena **acuannya ada dan harus dibaca**. Tiap medan yang dibangun
> membawa jejaknya — nama rule dan posisi bitanya — di `frontend/labels.ts`.

**Tiga jebakan yang dicatat saat pencabutan ini, supaya ronde berikutnya tidak kejeblos:**

1. **`<pyIncludedRuleXML>` dibuang lebih dulu.** Ekspor Pega menyematkan salinan utuh rule anak;
   sapuan datar membaca medan milik anak sebagai milik induk. Pada
   `Section/InputTreatyInOffer.xml` salinan itu **899.307 bita** dari 8.194.691.
2. **Berkas Section membundel indeks rule LAIN.** Caption `Treaty Contract Name`,
   `Contract Ref No`, `Teritorial Scope`, `Bordereaux Note`, `Accounting Mode`, dan `Treaty Year`
   muncul di `InputTreatyInOffer.xml` tetapi `pzIndexOwnerKey`-nya **`TREATYINNONPROPORTIONAL`**.
   Menyimpulkan pemilik sebuah medan dari **nama berkasnya** dapat salah total.
3. **Blok bergaris mati tidak dibangun.** Pega tidak dapat mengomentari tata letak, jadi blok
   dimatikan dengan kondisi mustahil — `1=2`, `1==2`, `Never` pada `pyContainerVisibleWhen` /
   `pyVisibleWhen` / `pyRowVisibleCondition`. ⚠️ `pyDisabledWhen>1=2` artinya **kebalikannya**:
   tidak pernah dinonaktifkan. Jangan tertukar.

## Kunci alami `NILAI_SELISIH` dan `NILAI_SEBELUM_PRO_RATE` — tiket 76 dan 77

Migrasi `443`. Keduanya berkunci alami **(`ID_VERSI_KONTRAK`, `KODE_BESARAN`, `KODE_MATA_UANG`)** —
`UQ_NILAI_SELISIH` dan `UQ_NILAI_SEBELUM_PRO_RATE`.

⛔ **Mata uang ADA DI DALAM kunci**, dan itu `ADR-0048` butir 3. Uji positif
`TestTiket76BesaranSamaDuaMataUangDiterima` membuktikannya: besaran yang sama dalam IDR dan USD
adalah **dua fakta yang sah**, dan `UQ` tanpa mata uang akan menolak yang kedua — menolaknya
diam-diam, sebab setiap uji negatif tetap hijau.

### ⚠️ Nomor invariannya BELUM ADA, dan itu ditagih

`Z00_KUNCI_ALAMI.sql` berbunyi *"constraint tanpa invarian tidak punya tempat untuk gagal"*, dan
`SPEC-INVARIAN.md` berhenti di **`INV-71`**. Kedua tiket menyatakan apa adanya *"kunci alaminya
belum punya nomor invarian"* sebagai persyaratan yang diketahui, bukan sebagai penghalang — dan
daftar periksa keduanya tetap menuntut kuncinya.

⛔ **Nomornya TIDAK dikarang di repo.** Memberi `INV-72`/`INV-73` dari sini berarti mengambil
wewenang pemilik `SPEC-INVARIAN.md`, dan nomor yang ditetapkan di dua tempat akan bertabrakan.
**Tagihan:** dua nomor invarian untuk kedua kunci alami di atas.

### ⚠️ `INV-58` dan kenapa `NILAI_SEBELUM_PRO_RATE` BUKAN pelanggarannya

*(Daftar periksa tiket 77 menuntut alasan ini berdiri di `MODUL.md`, bukan hanya di tiketnya.)*

`INV-58` melarang menyimpan nilai turunan. Tabel ini **tampak turunan dan ia bukan**, dan sebabnya
`GRL-15`: **mesin pro rata sengaja TIDAK dibangun kembali.**

Diukur 5 Oktober 2026 di `Activity/TreatyEDMProRateCalculation.xml` — satu-satunya rule di modul
ini yang menyentuh `ValueBeforeProrate`, dan langkahnya berbunyi persis:

```
pyStepsDescription : Copy value from ValueDifference to ValueBeforeProrate
PropertiesName     : TreatyIn.ValueBeforeProrate
PropertiesValue    : TreatyIn.ValueDifference
```

Ia **SALINAN yang diambil sebelum mesin pro rata mengubah angkanya** — bukan hasil hitungan. Dan
karena mesin itu tidak akan ada di sistem baru, angkanya **tidak dapat dipulihkan dengan
menghitung mundur**. Inilah pengecualian yang `INV-58` sediakan: **fakta terbukukan yang membawa
penunjuk asalnya**.

⛔ `TestTiket77BarisTetapSahWalauProRataTidakDipakai` menjaga pernyataan itu: baris tetap sah pada
versi yang `MEMAKAI_PRORATA`-nya **tidak** disetel. Tiket 77 menulis sebabnya sendiri — *"tabel
yang menolak baris itu mengaku dirinya hasil hitungan."*

### ⛔ Satu tuntutan tiket 77 yang TIDAK dapat dipenuhi — `MATA_UANG`

Tiket 77 menuntut kunci asing kedua:

```
MATA_UANG  1--<  NILAI_SEBELUM_PRO_RATE   [hapus: tolak]
```

`MATA_UANG` **dicabut** migrasi `434` ([`KEPUTUSAN` §16](../treatyin/docs/KEPUTUSAN-PENYELARASAN-REPO.md)):
kurs dan daftar mata uang kini dari `TREATYEXCHANGEYEARLY`. Kunci asing ke tabel yang tidak ada
tidak dapat dipasang, dan membangunnya kembali dilarang §16.

Jadi `KODE_MATA_UANG` berdiri sebagai **teks di dalam kunci alami, tanpa kunci asing** — setengah
dari yang tiket 77 tuntut, dan setengahnya dinyatakan alih-alih disamarkan. Polanya sama dengan
`KODE_BESARAN` terhadap `BESARAN_DAPAT_DISESUAIKAN`, penghalang yang tiket 76 catat sendiri.

**SIAPA DAPAT MENJAWAB:** pemilik proses — apakah `INV-44` masih menuntut daftar mata uang
berupa tabel sesudah §16, dan bila ya, tabel mana yang menggantikan `MATA_UANG`.
**SYARAT PEMBALIKAN:** begitu daftar mata uang punya tabel lagi, kedua kolom itu dinaikkan
menjadi kunci asing `tolak` lewat migrasi korektif — dan **harus sebelum tiket 06 memuat data**,
sebab kode yang tidak punya padanan akan menolak.

## Layar Adjustment — 5 Oktober 2026

Halaman awal modul kini layar `InputTreatyInAdjustment` sistem lama: daftar
penyesuaian, lalu **Old Data ‖ New Data** berdampingan, Attachment, deret tombol,
History. Catatan ukurnya di [`docs/LAYAR-ADJUSTMENT.md`](docs/LAYAR-ADJUSTMENT.md).

⛔ Datanya dari `TREATY_IN_EDM` + `M_TREATY_IN_EDM` — tabel warisan, **baca saja**,
nol migrasi. Bukan `VERSI_KONTRAK`: diukur nol baris. Rute
`GET /api/treaty-in-adjustment/penyesuaian-warisan` dan `…/penyesuaian-warisan/satu?id=`
(pengenalnya bergaris miring, jadi lewat parameter kueri).

⛔ Panel Old **nol medan dapat disunting**. Panel New dapat disunting di mode Edit
menurut `pyReadOnlyCondition` ekspornya, tetapi tombol Save **mati** — nol yang
terkirim. `NILAI_SELISIH` dan `NILAI_SEBELUM_PRO_RATE` tetap nol baris: tab Value
Difference menampilkan nilai TERSIMPAN dokumen, bukan hitung ulang.

Pertanyaan terbukanya: [`docs/PERTANYAAN-TERBUKA-LAYAR-ADJUSTMENT.md`](docs/PERTANYAAN-TERBUKA-LAYAR-ADJUSTMENT.md).

### Isi tab dibangkitkan dari ekspor

`alat/ekstrak_kerangka.py` membaca Section ekspor, membuang `<pyIncludedRuleXML>` dengan
**menghitung kedalaman sarang**, lalu menulis dua berkas bangkitan — ⛔ jangan disunting tangan:
`frontend/ekspor/kerangka.gen.ts` (37 tab: blok, grid, medan, desimal per sel, syarat) dan
`backend/repository/kunci_kerangka_gen.go` (daftar-izin kunci). Desimal per kolom mengikuti
ekspor (keputusan pemilik proses 5 Oktober 2026). Isi Retro tidak dibangun (§17). Uraian:
[`docs/LAYAR-ADJUSTMENT.md`](docs/LAYAR-ADJUSTMENT.md) §6.
