# 05: Peran menggantikan nama orang — tertunda sampai pemetaan diterima

**Status:** needs-info — mekanisme dibangun, pemetaan peran menunggu IAM *(implementasi 2026-10-03, cabang `modul/nbtreatyin/implementasi`; semula: needs-info)*
**Blocked by:** pemetaan **nama → peran** dari `[IAM]` dan `[work owner]`
**Menutup:** AC 12 · 13 · 81 · 82 · 91 *(5 AC)* — US 17 · 18

## Hasil & nilai pengguna

Hari ini wewenang di **dua belas tempat** ditentukan dengan memeriksa *"apakah pengguna ini orang
tertentu"*. `[terverifikasi]` Nilainya **tidak disalin** ke artefak mana pun. ⛔ Aturan seperti itu
berhenti bekerja ketika orangnya pindah jabatan atau keluar, **tanpa pemberitahuan**. Dan peran
sendiri disimpan di **kolom nomor telepon**.

Sesudah tiket ini, wewenang ditentukan **peran**, dan peran disimpan di medan peran — ⭐ kepergian
seseorang tidak lagi mematahkan alur.

## Blocker

⛔ **Keputusannya jelas, pelaksanaannya belum bisa.** `[keputusan work owner]` P12 dan P28
menetapkan penggantian dengan peran; ⛔ **pemetaan nama → peran tidak ada di korpus.**

| Yang dibutuhkan | Kenapa |
| --- | --- |
| ⛔ untuk tiap dari **12 tempat**: peran penggantinya | tanpa itu, penggantian adalah tebakan |
| ⛔⛔ untuk **dua layar**: **arahnya** | `[terverifikasi]` nama yang sama dipakai **dua arah** — satu bagian muncul **hanya** untuk orang itu, bagian lain untuk **semua kecuali** orang itu. ⛔ Satu peran untuk keduanya akan **membalik** salah satunya |
| ⛔ peran ketiga | ⭐ tangga punya **tiga** posisi, peran yang terbaca baru **dua** |

## Area codebase

- Lapisan service: pemeriksaan wewenang
- Lapisan repository: penyimpanan peran pengguna

## Rule Pega sumber

| Yang dibaca | Rule |
| --- | --- |
| Peran di kolom telepon | `When\IsTreaty1.xml` · `When\IsSPVTreaty1.xml` — dibandingkan terhadap `OperatorID.pyTelephone` |
| Guard identitas orang | **12 berkas**, lewat `OperatorID.pyUserIdentifier` |
| Guard di layar | **4 layar · 12 tempat**, salah satunya memakai medan identitas yang berbeda |

## ADR terkait

- **ADR-0002** — RBAC memakai peran yang sudah ada; rangkap peran ditolak

## Acceptance criteria

- [ ] 🟡 **AC 12** — wewenang ditentukan **peran**, bukan nama orang
- [x] **AC 13** — peran disimpan di medan peran; ⛔ kolom nomor telepon kembali berisi nomor telepon
- [x] **AC 81** — ke-12 tempat ditandai **tertunda**, ⛔ tidak dibangun dengan peran yang ditebak
- [x] **AC 82** — ⛔ **arah** pemeriksaan pada dua layar **tidak ditebak**
- [x] **AC 91** — ⛔ peran karangan **tidak dibuat** untuk menutup kekurangan posisi ketiga

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| ⛔ **1** | pemetaan nama → peran untuk 12 tempat, **berikut arahnya** | ⛔ **MENAHAN** |
| ⛔ **2** | peran yang tersedia belum cukup untuk tiga posisi | ⛔ **MENAHAN** |
| **3** | penetapan, pembuatan, dan pencabutan peran belum dibahas | tidak menahan |

## Catatan

⛔ **Menebak peran berarti memberi atau mencabut wewenang atas dasar tebakan.** Tiket ini ditulis
lengkap supaya siap dikerjakan begitu pemetaannya tiba, ⛔ **bukan supaya dikerjakan sekarang.**

## Hasil implementasi 2026-10-03 — mekanisme dibangun, pemetaan tetap menunggu

- Tabel `M_NBTRIN_PERAN_TEMPAT` (migrasi 330: KODE_TEMPAT, PERAN, ARAH `MUNCUL`/`KECUALI`), **nol
  baris** — diisi IAM bersama work owner. Tempat tanpa baris = **tertunda** (AC 81); dua arah di satu
  tempat = tertunda (AC 82). Uji memakai peran fiktif `UJI-` (AC 91).
- Tempat identitas di rule TERJANGKAU dan nasibnya:
  | Tempat | Nasib |
  | --- | --- |
  | `ListSuggest .ProductionDate` (tampil + wajib) | `LISTSUGGEST_PRODUCTIONDATE` lewat tabel |
  | tiga tombol Submit `DetailDeptHeadTreatyIn_UW` (`<ID-operator-1>`) | ⚠️ diganti **posisi kasus** Dept Head — diturunkan dari tangga P13, bukan peran tebakan; mohon konfirmasi |
  | label NON EDM / EDM `DetailPoliciesNonProportional` (`<ID-operator-2>`) | bagian XOL non-proporsional tidak dibangun (P29) |
  | `When\IsSPVCreate`, `IsTreaty1`, `IsSPVTreaty1` | hanya memilih Assignment4/6 (posisi sama) — tidak dibangun |

## ⭐ Putaran 2 — paket penyimpanan (03-10-2026)

Dasar: PROMPT-NB-TREATY-IN-PUTARAN-2 bab 0 butir 11–12, bab 2 K4/K16/K17; rincian kolom `docs/PERBANDINGAN-KOLOM-DIAGRAM.md`.

⛔ **RALAT** atas *Hasil implementasi 2026-10-03* butir pertama. Bunyi lama, dikutip: *"Tabel
`M_NBTRIN_PERAN_TEMPAT` (migrasi 330: KODE_TEMPAT, PERAN, ARAH `MUNCUL`/`KECUALI`), **nol baris** — diisi IAM
bersama work owner."*

Bunyi baru (`[keputusan work owner]` **K16**): tabel itu tidak ada di diagram grilling — migrasi 330 dan
pembacanya (`repository.DaftarPeranTempat`, `services.Gudang.DaftarPeranTempat`, tiruan) **dihapus**. Pemetaan
tempat → peran → arah kini **konstanta kode** `backend/models/peran_tempat.go` (`PemetaanPeranTempat`),
**kosong sampai IAM menjawab** (K12). Peran pengguna tetap dari `inti.Pelaku.Peran`. Perilaku saat kosong sama
dengan tabel kosong: setiap tempat tertunda (AC 81), arah tidak ditebak (AC 82; dua arah atau arah asing di
satu tempat = tertunda), nol peran karangan (AC 91). Uji: `models/peran_tempat_test.go`
(`TestTempatTampilMenurutArah`, `TestPemetaanPeranTempatKosongSampaiIAMMenjawab`), `handlers/alur_test.go`
`TestTempatBerperanTidakDitebak`.

**Status tetap needs-info** — pemetaan 12 tempat dan arahnya menunggu IAM; mengisinya = menyunting satu
variabel, bukan migrasi.

## ⭐ Putaran 2 — paket P8: work group di gerbang portal (03-10-2026)

Butir terbuka **baru** (tidak menahan tiket 04): wadah grid `Section\SFAPortal_OpportunitiesList.xml` bersyarat
`OperatorID.pyWorkGroup!='ReasLife' && OperatorID.pyWorkBasketList(2).pyWorkBasketName=='ReasTreatyInAdmin'`, dan
wadah dalamnya `!IsOperatorLife` (`When\IsOperatorLife.xml`: `OperatorID.pyWorkGroup = "ReasLife"`). Klausa
workbasket dibangun menurut nama (tiket 04, P8). ⛔ Klausa **work group** tidak dibangun: `inti.Pelaku` hanya
membawa AkunID + workbasket, dan M_LOGIN_GO hanya CODE organisasi/divisi/unit — pemetaan work group Pega
`ReasLife` ke data itu **tidak ada** (K12 kosong). Menyimpulkannya dari peran `ReasLife*` = menebak (AC 91).
Akibatnya: operator ber-work-group `ReasLife` yang juga anggota `ReasTreatyInAdmin` melihat daftar di sistem baru
tetapi tidak di Pega. Yang dibutuhkan dari `[IAM]`: padanan work group untuk akun (atau konfirmasi bahwa syarat itu
boleh dibuang). **Status tetap needs-info.**
