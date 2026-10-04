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
| `docs/` | `STRUKTUR-TABEL-TREATY-IN-ADJUSTMENT.md`, `KEPUTUSAN-TIKET-02-03.md`, `issues/` (tiket `01`, `02`, `03`, `05`) |
| `backend/` | `modul.go` (`Pendaftaran()`), `models/`, `repository/`, `services/`, `handlers/`, `migrations/` |
| `frontend/` | `menu.ts`, `rute.tsx`, `api.ts`, `labels.ts`, `pages/RantaiVersi.tsx` |

## Migrasi

Rentang `440-479` (tabel R2, urut hulu ke hilir: migrasi modul hilir yang merujuk tabel modul hulu
selalu berjalan sesudahnya). Slot menu `974-975` hanya menyalakan `DIMIGRASI` baris modul ini (satu `UPDATE`,
nol `INSERT` — menu datar 30-09-2026) saat modul mendapat layar pertamanya, di folder
`backend/migrations/` modul ini sendiri — bentuk SQL-nya di `APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`
bab 6. Nomor selalu tiga digit.

## Pernyataan untuk penjaga

### Kaskade ON DELETE CASCADE

**Nol kaskade di modul ini**, dan itu keputusan yang dikutip, bukan bawaan yang dibiarkan.

Modul ini punya satu kunci asing sendiri: `VERSI_KONTRAK.ID_VERSI_KONTRAK_DASAR` →
`VERSI_KONTRAK` (migrasi `440`). `4-erd-dan-tabel-datar/ERD.md` §2.2 menyatakannya **[hapus: tolak]**
— *"menghapus versi dasar akan membuat seluruh baris selisih kehilangan artinya"* — sehingga ia
berdiri tanpa klausa `ON DELETE`.

| Awalan berkas | Relasi |
| --- | --- |

Tabel di atas **sengaja kosong**, dan itu bukan kelalaian: tabel kosong menyatakan **nol berkas
modul ini yang boleh memuat `ON DELETE CASCADE`**, dan `TestKaskadeHanyaPadaRelasiTerdaftar` akan
menolak yang pertama kali menambahkannya tanpa mendaftarkannya di sini.

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
