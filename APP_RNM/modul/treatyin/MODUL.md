# Modul `treatyin` — Treaty In

⚠️ **Terdaftar 1 Oktober 2026 — LAPISAN SKEMA empat belas tiket, bukan tiketnya.**
`backend/modul.go` ada, `inti/backend/daftar/modul_treatyin_gen.go` bangkit, dan slot menu `972`
menyalakan `M_NAV_MENU.DIMIGRASI`.

Yang berdiri (migrasi `400`–`415`): **23 tabel, 192 kolom, 23 sequence** — tiket `14` `15` `20`
`22`–`31` `39`. Ditambah satu jalur baca atas tabel acuan. Cacah kolom tiap tabel **sama persis**
dengan `KAMUS-KOLOM.md`, kecuali `VERSI_KONTRAK` yang bertambah satu kolom dari tiket `01` papan
Adjustment.

⛔ **KEEMPAT BELAS TIKET ITU BELUM SELESAI**, dan ini bukan formalitas. Papan tiketnya menetapkan ukuran
selesai: *"tiap tiket membawa uji negatif DAN uji positif"*, sebab *"'dapat diperagakan' tidak
tersedia sebagai bukti"*. Uji di modul ini **membaca teks DDL dan gudang tiruan** — tidak satu pun
menjalankan Oracle (`L-3`). Yang BELUM terbukti: versi yatim ditolak, nomor urut ganda ditolak,
tanggal terbalik ditolak, dan uji positif *"satu kontrak dengan tiga versi diterima, lapisan bekunya
tidak disalin"*. Daftarnya di [`docs/issues/README.md`](docs/issues/README.md).

⛔ **Sembilan invarian tabelnya berdiri tetapi BELUM ditegakkan di mana pun** — `INV-29`, `INV-38`,
`INV-39`, `INV-40`, `INV-53`, `INV-55`, `INV-56`, `INV-57`, dan `INV-49` yang bahkan belum punya
tempat berdiri. Sebagian besar membandingkan baris di tabel BERBEDA, yang `CHECK` Oracle tidak dapat
nyatakan; `INV-29` dan `INV-38` **dapat**, dan yang menahannya hanya ADR-0056. Seluruhnya menunggu
jalur simpan yang sama. Daftar lengkapnya di
[`docs/STRUKTUR-TABEL-TREATY-IN.md`](docs/STRUKTUR-TABEL-TREATY-IN.md).

⛔ **`UQ_LAYER` tidak menegakkan `INV-05` sepenuhnya**: `BAGIAN_LAYER` boleh kosong, dan Oracle
memperlakukan NULL sebagai tidak sama dengan NULL di kunci unik komposit. Lubangnya ada di
`Z00_KUNCI_ALAMI.sql`, bukan dibuat di sini; tagihannya di dalam migrasi `413`.

⛔ **Sequence berdiri tetapi BELUM TERPAKAI.** Tiket `14` menuntut *"tidak ada jalur lain yang dapat
memberi pengenal"*; hari ini belum ada jalur tulis sama sekali, sehingga tuntutan itu benar secara
hampa. Ia ditegakkan bersama jalur simpan.

ℹ️ **Layar modul ini ada karena mendaftarkan modul menuntutnya** `[keputusan work owner 01-10-2026]`.
Penjaga `daftar.modulAktif.test.ts` mewajibkan setiap modul terdaftar punya `frontend/`, sementara
papan tiket melarang mengarang layar (`L-4`). Jalan keluarnya: satu halaman **baca-saja** atas data
yang tiket `15` buat — nol alur karangan.

⚠️ **Lima penyelarasan spec dengan repo** tercatat di
[`docs/KEPUTUSAN-PENYELARASAN-REPO.md`](docs/KEPUTUSAN-PENYELARASAN-REPO.md). **Jangan menyunting DDL
tanpa membacanya.**

⛔ **Tabel di bawah dibaca penjaga** (`inti/backend/penjaga`): rentang migrasi dan slot menu. Ubah
nilainya hanya lewat pull request yang disetujui tim inti — dua modul tidak boleh berbagi nomor.

| Kunci | Nilai |
| --- | --- |
| Nama modul | `treatyin` |
| Folder korpus | `Treaty In` |
| GROUPMENU | `TREATY` |
| Pemilik | `@PEMILIK-TREATYIN` |
| Status | dimigrasi |
| Rentang migrasi | `400-439` |
| Slot menu | `972-973` |
| Prefix rute API | `/api/treaty-in` |
| Kontrak disediakan | — |
| Kontrak dipakai | — |

`Pemilik` adalah penanda pemegang modul. Wilayah berkas yang boleh disentuh cabang
`module/<nama>` dijaga `.github/workflows/penjaga-wilayah-cabang.yml` - CODEOWNERS
dipensiunkan 1 Oktober 2026.

## Isi folder

| Folder | Isi |
| --- | --- |
| `docs/` | `STRUKTUR-TABEL-TREATY-IN.md` (delapan tabel), `KEPUTUSAN-PENYELARASAN-REPO.md` (empat penyelarasan), `issues/` (tiket `14` dan `15`) |
| `backend/` | `modul.go` (`Pendaftaran()`), `models/`, `repository/`, `services/`, `handlers/`, `migrations/` |
| `frontend/` | `menu.ts`, `rute.tsx`, `api.ts`, `labels.ts`, `pages/AcuanTreatyIn.tsx` |

## Migrasi

Rentang `400-439` (tabel R2, urut hulu ke hilir: migrasi modul hilir yang merujuk tabel modul hulu
selalu berjalan sesudahnya). Slot menu `972-973` hanya menyalakan `DIMIGRASI` baris modul ini (satu `UPDATE`,
nol `INSERT` — menu datar 30-09-2026) saat modul mendapat layar pertamanya, di folder
`backend/migrations/` modul ini sendiri — bentuk SQL-nya di `APP_RNM/PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md`
bab 6. Nomor selalu tiga digit.
