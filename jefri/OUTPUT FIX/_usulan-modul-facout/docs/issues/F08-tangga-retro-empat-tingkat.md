# F08: Tangga persetujuan retrosesi — empat tingkat

**What to build:** Penawaran retrosesi berjalan melalui tangga peran tersendiri sampai disetujui atau
ditolak, lengkap dengan percabangan grup dan jejak persetujuan.

⛔ **Ini BUKAN Seam 2.** `[terverifikasi]` `Flow\OfferFacRetro` memuat **nol shape Approval** di
keempat salinannya — ia **tidak** memakai mesin akseptasi limit-berjenjang Fac In. Tangga retro adalah
mesin tersendiri (**K-055**).

## Tangga empat tingkat berlaku **MENYELURUH untuk ketiga siklus**

✅ **K-058** menetapkan: rule yang sudah **diekspor ulang ke `DDL\` MENANG** atas salinan korpus.
`DDL\OfferFacRetro.xml` memuat tangga **empat tingkat**, maka **empat tingkat berlaku untuk NB, RNW
dan EDM** — tanpa pembatas lingkup.

`[terverifikasi]` Empat salinan, **keempatnya ber-`pyRuleSetVersion` IDENTIK `01-01-95`**:

| Salinan | Ukuran | Assignment | Gateway | Workbasket | Status |
| --- | ---: | ---: | ---: | --- | --- |
| `DDL\OfferFacRetro.xml` (ekspor ulang, commit 2026-08-31) | 171.308 B | 5 | 6 | Admin/Head/GroupLeader/TechnicalDirector | ✅ **BERLAKU** |
| `Endorsment Fac In\Flow\OfferFacRetro.xml` | 170.966 B | 5 | 6 | idem — **4 tingkat** | sejalan |
| `NB FacIn\Flow\OfferFacRetro.xml` | 134.214 B | 3 | 4 | Admin + Head — **2 tingkat** | ⛔ **BASI** |
| `RNW Fac In\Flow\OfferFacRetro.xml` | 134.214 B | 3 | 4 | **identik byte** dengan NB | ⛔ **BASI** |

⚠️ **Salinan NB/RNW yang 2 tingkat dinyatakan BASI — dicatat, bukan dihapus** (`PANDUAN-KERJA` §7).
Ia tetap berguna sebagai pembanding untuk mengukur apa yang berubah.

⛔ **Nomor versi tidak menolong di sini** — keempatnya `01-01-95`, isinya berbeda. Ini salah satu
dasar pertanyaan terbuka K-058 tentang metode deteksi drift.

## ✅ Layar keempat peran SAMA

**K-058:** perbedaan antar salinan **hanya pada ALURNYA** (2 tingkat vs 4 tingkat). **Tampilan menu
dan isinya SAMA untuk keempat peran** — `ReasFacOutGroupLeader` dan `ReasFacOutTechnicalDirector`
memakai **layar yang sama** dengan `ReasFacOutAdmin` dan `ReasFacOutHead`.

📌 Konsekuensi implementasi: **satu set layar**, empat tahap alur. Jangan membangun layar terpisah
per peran.

## Isi tangga

`[terverifikasi]` Dari `<pyMOName>` pada `DDL\OfferFacRetro.xml`: empat peran
`ReasFacOutAdmin` · `ReasFacOutHead` · `ReasFacOutGroupLeader` · `ReasFacOutTechnicalDirector`,
masing-masing **satu kali**; `Is it group?` / `IsGroup`; `Accept?` **4×**; `reject` **4×**;
`confirm` **4×**; `IsRISlip` **2×**; `OfferFacOut` **4×**; routing `pyRouteTo` = `Custom`.

`[terverifikasi]` Siklus hidup layarnya: `Activity\OfferFacOut_PreAct` bergerbang
`.pyWorkBasketName=="ReasFacOutAdmin" || .pyWorkBasketName=="ReasFacOutHead"`, memanggil
`GetOldDataRetro_ACT`, `SumCurrencyListAllRetro_Act`, `GetHistoryAkseptasiPega_Act`,
`ProtekReinsurerList_Act`, `CheckDataFacOut_Act`. `OfferFacOut_PostAct` menulis jejak
`ViewSuggest(<APPEND>).Approval` / `.CommentSuggest` / `.DateSuggest` dan membaca
`EmailTypeRetro` / `EmailTypeRetroSlip` bernilai `1` / `2` / `7`.

⚠️ **Satu keputusan manusia = satu transisi** (`GLOSARIUM.md`) — bukan satu proses yang menghitung
seluruh rantai sekaligus.

**Asal (Pega).** `DDL\OfferFacRetro.xml` · `Flow\OfferFacRetro` (tiga salinan korpus) ·
`Activity\OfferFacOut_PreAct` · `OfferFacOut_PostAct`

**Keputusan.** **K-055** · **K-058** (salinan `DDL\` menang; layar keempat peran sama) · K-053 ·
`CLAUDE.md` §4.6

**Blocked by:** F05

**Status:** blocked

- [ ] Empat tingkat berurutan, **satu keputusan = satu transisi**
- [ ] Percabangan grup hidup sebagai gerbang tersendiri
- [ ] Tiap tingkat punya jalur **terima** dan **tolak**
- [ ] ⛔ **Tidak memakai mesin akseptasi Fac In (Seam 2)** — nol shape Approval
- [ ] Jejak persetujuan tercatat per langkah (siapa, kapan, komentar) — **identitas disimpan sebagai rujukan, nilainya tidak disalin ke artefak** (K-025)
- [ ] Enumerasi tipe email `1`/`2`/`7` **diport apa adanya**; artinya `[pertanyaan terbuka]` → `panic` bila nilai lain muncul
- [ ] Empat tingkat berlaku **menyeluruh** untuk NB, RNW dan EDM (K-058) — tanpa cabang per siklus
- [ ] **Satu set layar** dipakai keempat peran; GroupLeader dan TechnicalDirector memakai layar Admin/Head
- [ ] Komentar menyebut salinan `DDL\OfferFacRetro.xml` sebagai sumber, dan menandai salinan NB/RNW 2 tingkat sebagai **basi**
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
