# 14: Lingkup — yang TIDAK dibangun, dan diuji supaya tetap tidak terbangun

**Status:** selesai *(implementasi 2026-10-03, cabang `modul/nbtreatyin/implementasi`; semula: ready-for-agent)*
**Blocked by:** —
**Menutup:** AC 61 · 62 · 63 · 64 · 65 · 88 *(6 AC)* — US 45

## Hasil & nilai pengguna

Hari ini folder modul memuat **lebih banyak aturan daripada pekerjaan modul ini**. `[terverifikasi]`
Ia **dependency closure** — memuat aturan milik modul lain karena pewarisan kelas. ⛔ Membangun dari
isi folder apa adanya akan mewarisi **beban mati** dan pekerjaan yang bukan milik modul ini.

Sesudah tiket ini, sistem baru **tidak memuat** apa yang sudah dinyatakan di luar lingkup, dan
⭐ ada **uji yang menjaganya tetap begitu** — supaya ia tidak masuk kembali diam-diam pada ronde
pembangunan berikutnya.

## Area codebase

- Uji lingkup: daftar yang tidak boleh terpanggil

## Rule Pega sumber

| Yang **tidak** dibangun | Sebab | Butir |
| --- | --- | --- |
| **28** aturan yatim, **561** langkah penetapan | ⛔ tidak terjangkau; milik modul Fac | keadaan Bab 2 |
| **6** aturan pembongkar dokumen | penyimpanan pindah ke sumber relasional | P29 · P15 |
| aturan bernama sama yang menguji `= 1` | ⛔ **bukan aturan yang hidup** | P6 |
| penanda persetujuan **kedua** | ⛔ sudah tidak dipakai; ⭐ **berkas rule-nya nol di seluruh korpus** | P36 |
| medan "nomor surat" sebagai penanda arah | routing memakai medan tersendiri | P7 |
| pemanggilan penyalin data yang tidak dipakai | ⛔ tidak digunakan | P30 |

## ADR terkait

- **ADR-0009** — migrasi penuh; tidak ada koeksistensi dua penulis

## Acceptance criteria

- [x] **AC 61** — **28** aturan yatim **tidak dimigrasi**; nol terpanggil
- [x] **AC 62** — **6** aturan pembongkar dokumen **tidak dimigrasi**
- [x] **AC 63** — aturan bernama sama yang menguji `= 1` **tidak dimigrasi**
- [x] **AC 64** — penanda persetujuan kedua **tidak dibangun**; penampungnya tidak dibuat
- [x] **AC 65** — medan "nomor surat" **tidak** dipakai menyimpan penanda arah
- [x] **AC 88** — nol dari **561** langkah yatim terpanggil

## Perintah verifikasi

1. Jalankan uji lingkup — ⭐ **nol** dari daftar di atas terpanggil.
2. Cari penanda persetujuan kedua di seluruh sistem baru — ⭐ **tidak ada**.

## Catatan

⚠️ `[terverifikasi]` **Penanda persetujuan kedua dirujuk juga oleh modul Fac** — tiga berkasnya
merujuk aturan yang sama, yang berkasnya **tidak ada di seluruh korpus**. ⭐ Catat untuk modul itu;
⛔ **jangan kerjakan di tiket ini.**

## ⛔ RALAT K8 — AC 62 (putaran 2, 03-10-2026)

Bunyi lama (Acceptance criteria): *"**AC 62** — **6** aturan pembongkar dokumen **tidak dimigrasi**"*;
Rule Pega sumber: *"**6** aturan pembongkar dokumen | penyimpanan pindah ke sumber relasional | P29 · P15"*.

Bunyi baru — `[keputusan work owner]` **K8** (jalur NonProp/XOL dibangun, master dibaca baca-saja):

| Aturan | Nasib | Bukti |
| --- | --- | --- |
| `InputPolicyTreatyInDetail_NonProp` | ✅ **dimigrasi** (langkah 2-6 lewat `repository.MasterXOLDariJSON`, baca-saja) | `RDBList\BrowseTreatyInJoinEDM`; `models.InputDetailNonProp` |
| `SetTreatyIn_Act` | ✅ **sebagian** — langkah 3-5 (`RDBList\BrowseTreatyIn`) dan 13 (`TreatySetReinstatement`), dalam rantai `TreatyRealizationCheckXOLList` | `services.siapkanNonProp` |
| `InputPolicyTreatyInDetail_preACT` | langkah 16 dan 18 ✅ (K8); 9-10 (`M_TREATY_IN_DETAIL_EDM`, di luar tabel K8), 13, 17 ⛔ | `services.pilihBisnisNonProp` |
| `FetchMasterTreatyIn` | ⛔ tetap (master proporsional, bukan jalur XOL) | P29 |
| `InputPolicyTreatyOutDetail_preACT`, `InputPolicyTreatyOutDetail_NonProp` | ⛔ tetap — treaty KELUAR (K8 butir 4): `select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={pyWorkPage.PolicyTreatyIn.NoOffer}` | `RDBList\BrowseTreatyOut` |

AC 62 dengan bunyi baru: **nol** pembacaan JSON di luar `repository.MasterXOLDariJSON`, **nol** penulisan
JSON di mana pun, dan keempat aturan yang tetap ⛔ tidak terpanggil. Status AC tetap ✅.
