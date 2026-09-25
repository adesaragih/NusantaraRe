# Usulan diff ke artefak induk — **DISETUJUI; ENAM diterapkan, SATU diparkir**

**Tanggal:** 24 September 2026
**Aturan yang berlaku:** **artefak induk disunting HANYA sesudah diff-nya disetujui.**

> ## KEADAAN — diperbarui 24 September 2026
>
> Pemilik proses **menyetujui ketujuh diff**, dengan empat syarat dan satu urutan.
> **Lima diterapkan 24 September 2026 pagi; `D-1` menyusul sore hari ketika `F-2` ditutup.**
> **Enam diterapkan; satu diparkir.**
>
> | Diff | Artefak | Keadaan |
> |---|---|---|
> | **`D-3`** | `STRUKTUR-DATA.md` | **DITERAPKAN 24 Sep 2026** — pertama, sebab ia **mengikat** |
> | **`D-2`** | `SPEC-MODEL-DATA.md` §2.3 · §10.19a · **§10.19b** | **DITERAPKAN 24 Sep 2026** — beserta arti *"kosong untuk seluruh baris warisan"* (syarat 2) |
> | **`D-4`** | `SPEC-INVARIAN.md` §9 | **DITERAPKAN 24 Sep 2026** — beserta **pemantau kebasian**; `INV-69`/`INV-70` bertanda **KLAIM** (syarat 3) |
> | **`D-5`** | `UJI-NEGATIF-INVARIAN.md` `N-7` | **DITERAPKAN 24 Sep 2026** — tujuh uji, **empat positif** |
> | **`D-6`** | `PETA-TELUSUR-JSON.md` §4a | **DITERAPKAN 24 Sep 2026** — beserta **peringatan semesta `L-8`** |
> | **`D-1`** | `SPEC-MODEL-DATA.md` §10.2a | **DITERAPKAN 24 Sep 2026** — parkirnya dicabut oleh **penutupan `F-2`** di sesi to-spec induk. Nomor urutnya **ke-44**, bukan ke-48; isi diffnya **tidak diubah** |
> | **`D-7`** | `7-2-ACTUALVALUE.md` | **DIPARKIR** — dan **dasarnya sudah lewat**: `F-1` ternyata **sudah diterapkan** hari ini. Parkirnya **tidak dicabut sendiri**; pencabutannya wewenang pemilik proses |
>
> **Berkas usulan yang tidak menyatakan mana yang sudah diterapkan akan diterapkan dua kali.**

`GRL-01` mengikat: modul ini tidak punya spesifikasi sendiri. Hasil to-spec Adjustment **mendarat di
artefak induk**; folder `2-to-spec/` Adjustment hanya menerima catatan kerja dan usulan ini.

---

## D-1 · `SPEC-MODEL-DATA.md` §10.2 — `SIFAT_MATERIAL_ADDENDUM` KEMBALI

| | |
|---|---|
| **Bagian** | §10.2 `VERSI_KONTRAK`, tabel atribut |
| **Kenapa** | `GRL-12` **BATAL**; `GRL-20` mengembalikan materialitas sebagai **atribut masukan** |
| **Label** | **PELESTARIAN** — sistem lama memang menyimpannya (`TreatyIn.EDMMaterialType`) |

Baris yang diusulkan:

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `SIFAT_MATERIAL_ADDENDUM` | `EDMMaterialType` | E | **ya** — kosong pada versi pertama | dua nilai **`MATERIAL`** / **`TIDAK_MATERIAL`**. **Masukan, bukan turunan** (`GRL-20`). Beku sejak `AJUKAN`, berjejak selama `DRAFT` (`KTV-1`) |

> **§10.2a sudah memuat baris ini** sebagai atribut ke-48, ditambahkan 24 Sep 2026. Yang berubah
> **bukan keberadaannya melainkan dasarnya**: ia sempat dihapus oleh `GRL-18` atas dasar `GRL-12`,
> dan `GRL-20` mengembalikannya. **Usulannya: pertahankan barisnya, ganti catatannya.**

## D-2 · `SPEC-MODEL-DATA.md` §2.3 dan §10 — entitas `DOKUMEN_ADDENDUM`

| | |
|---|---|
| **Kenapa** | `GRL-19` — `GRL-01` **dikoreksi, bukan dibongkar**; jumlah entitas bertambah **satu** |
| **Label** | **BARU** — sapuan nomor dokumen **nol**, terkalibrasi |

Atribut yang diusulkan:

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_DOKUMEN_ADDENDUM` | *baru* | C1 | tidak | |
| `NOMOR_DOKUMEN` | **tidak ada di sistem lama** | T | tidak | **kunci alami**, unik **global**; `UA-19` mengujinya |
| `TANGGAL_BERLAKU` | *baru* | D | **ya** | `KTV-2`; **dicabut bila `DB-16b` dibantah** |

Dan pada `VERSI_KONTRAK`:

| Nama | Asal | Tipe | Boleh kosong | Catatan |
|---|---|---|---|---|
| `ID_DOKUMEN_ADDENDUM` | *baru* | R | **ya** | satu dokumen memayungi banyak versi, lintas kontrak. **Persetujuan tetap per versi** |

> **Kosong untuk SELURUH baris warisan.** Migrasi tidak punya sumber. Masuk **rencana peralihan**.

## D-3 · `4-erd-dan-tabel-datar/STRUKTUR-DATA.md` — daftar entitas

**Berkas ini MENGIKAT** soal daftar entitas (urutan wewenang butir 4). **Harus diselaraskan sebelum
entitas baru ditambahkan di mana pun.**

| Entitas | Bita | Arti | Kunci alami | Induk |
|---|---:|---|---|---|
| **`DOKUMEN_ADDENDUM`** | 17 | satu dokumen addendum yang disepakati cedant, memayungi satu atau beberapa versi kontrak | nomor dokumen, unik **global** | **— (berdiri sendiri)** |

## D-4 · `SPEC-INVARIAN.md` — `INV-69` dan `INV-70`

| Inv | Bunyi | Golongan |
|---|---|---|
| **`INV-69`** | Pada `VERSI_KONTRAK` bersifat **`TIDAK_MATERIAL`**, tidak ada baris `NILAI_SELISIH` yang besarannya bertipe uang atau porsi | **CONSTRAINT** lewat *materialized view* — ia lintas baris |
| **`INV-70`** | Pada `VERSI_KONTRAK` bersifat **`MATERIAL`**, teks kesepakatan dan identitas kontrak tidak berubah terhadap versi dasarnya | **CONSTRAINT**, idem |
| **`INV-71`** | `NOMOR_DOKUMEN` unik di seluruh `DOKUMEN_ADDENDUM` | **CONSTRAINT** — `UNIQUE (NOMOR_DOKUMEN)` |

Ketiganya lolos dua uji: **dapat gagal**, dan **dapat dikompilasi menjadi nama kolom**.

## D-5 · `UJI-NEGATIF-INVARIAN.md` — `N-7`

| # | Yang disisipkan | Harus |
|---|---|---|
| `N-7a` | versi `TIDAK_MATERIAL` + satu baris selisih premi | **GAGAL** |
| `N-7a+` | versi `TIDAK_MATERIAL` + satu baris selisih **tanggal pelaporan** | **BERHASIL** |
| `N-7b` | versi `MATERIAL` + perubahan `PENGECUALIAN` | **GAGAL** |
| `N-7b+` | versi `MATERIAL` + perubahan limit, pengecualian tidak disentuh | **BERHASIL** |
| **`N-7c+`** | **versi warisan yang melanggar keduanya** | **BERHASIL DIMUAT** |
| `N-7d` | dua `DOKUMEN_ADDENDUM` bernomor sama | **GAGAL** |
| `N-7d+` | satu dokumen memayungi versi dari **dua kontrak berbeda** | **BERHASIL** |

> **Uji berakhiran `+` bukan kerapian.** Ketiga kekeliruan lingkup di korpus ini **seluruhnya lolos
> uji negatif**; hanya uji positif yang menangkapnya. **`N-7c+` yang paling berharga** — tanpa ia,
> invariannya menolak data lama pada hari peralihan.

## D-6 · `PETA-TELUSUR-JSON.md` — nasib 279 jalur cermin

Menambahkan nasib untuk lingkup yang selama ini dikecualikan: **154 DISIMPAN · 12 TURUNAN ·
113 DIBUANG · 0 DITUNDA**, menutup ke **279**.

> **Dengan peringatan yang wajib ikut:** jalur ini **tidak boleh dihitung sebagai kelengkapan**.
> Semestanya tetap kurang **340 properti** (`L-8`).

## D-7 · `4-erd-dan-tabel-datar/7-2-ACTUALVALUE.md` — kalimat untuk pembaca arsip

Menambahkan kalimat §6 `STRUKTUR-ADDENDUM.md` ke berkas itu **dan** ke tempat arsip `JSONDATA`
dirujuk, sebab pembaca arsip tidak membuka berkas temuan.

---

## Dua temuan yang diusulkan ke induk — keduanya `SISI` **IRISAN**

| # | Temuan | Kenapa milik induk |
|---|---|---|
| **`TDA-18`** *(calon)* | **`ReinstatementPct` disalin, tidak dikurangi** — perubahan persentase reinstatement **tidak pernah menghasilkan baris selisih**. Sapuan seluruh korpus, kalibrasi lulus (`Limit` 6 pengurangan nyata): **10 penugasan, nol pengurangan** | `TreatyEDMDifferenceLimits` ada di **kedua ekspor** |
| **`TDA-17`** | tiga kolom tanpa rumah — `DEDUCTIBLE2`, `PREMIUM_EARNED`, `ROL_PCT` | lubang pada §10 **induk**, sebabnya `L-8` |

> **Keduanya belum melewati ronde TDA mana pun.** `TDA-18` lahir hari ini; `TDA-17` lahir sesudah
> ronde TDA ditutup. **Tidak diperbaiki diam-diam.**

---

## Ringkasan — tujuh diff, dua temuan

| Urutan | Diff | Artefak induk | Sifat | **Keadaan** |
|---:|---|---|---|---|
| 1 | `D-3` | `STRUKTUR-DATA.md` | entitas baru — **MENGIKAT** | **DITERAPKAN** |
| 2 | `D-2` | `SPEC-MODEL-DATA.md` §2.3 · §10.19a · §10.19b | entitas baru + arti "kosong" | **DITERAPKAN** |
| 3 | `D-4` | `SPEC-INVARIAN.md` §9 | tiga invarian + **pemantau kebasian** | **DITERAPKAN** |
| 4 | `D-5` | `UJI-NEGATIF-INVARIAN.md` `N-7` | tujuh uji, **empat positif** | **DITERAPKAN** |
| 5 | `D-6` | `PETA-TELUSUR-JSON.md` §4a | 279 nasib + peringatan semesta | **DITERAPKAN** |
| 6 | `D-1` | `SPEC-MODEL-DATA.md` §10.2a | ganti catatan, pertahankan baris | **DITERAPKAN** — `F-2` ditutup |
| — | `D-7` | `7-2-ACTUALVALUE.md` | kalimat pembaca arsip | **DIPARKIR** — menunggu izin; dasarnya sudah lewat |

**ENAM diterapkan 24 September 2026. SATU diparkir** — `D-7` — dengan pernyataan keputusan di
artefak induknya sendiri, bukan `TODO`.

> ### Yang mencabut parkir `D-1`, dan bagaimana
>
> **`F-2` ditutup** di sesi to-spec induk, putaran penutup. Penutupannya **bukan** menemukan keempat
> atribut yang hilang: pencariannya diselesaikan atas sumber mandiri dan **nol calon tersisa**, lalu
> **angka judulnya** yang diturunkan — dari 47/48 menjadi **45** yang dihitung (43 §10.2 + 1 §10.2a
> + 1 §10.19a).
>
> Itu memenuhi syarat yang `D-1` sendiri pasang: *"menambah satu memberi 45, bukan 48"*. Nomor urut
> baris `SIFAT_MATERIAL_ADDENDUM` kini **ke-44**, dan itu angka yang dapat dipertanggungjawabkan.
>
> **Isi `D-1` tidak diubah satu huruf pun.**

### Yang berubah dari isi usulan semula: **tidak ada**

Ketujuh diff diterapkan **apa adanya**. Keempat syarat pemilik proses **menambah**, tidak mengganti:
syarat 1 memarkir `D-1`, syarat 2 menambah §10.19b, syarat 3 menambah pemantau kebasian dan tanda
KLAIM, syarat 4 memarkir `D-7`.

> **Satu hal dilaporkan, tidak dibetulkan sendiri.** Syarat 4 mengandaikan koreksi `F-1` **belum**
> diterapkan pada `7-2-ACTUALVALUE.md`. **Ia sudah** — §1.2, §1.3, dan §1.4 dikoreksi hari ini di
> bawah butir `P-1`. Penahan `D-7` karenanya **sudah lewat**, tetapi parkirnya **tidak dicabut
> sendiri**: perintah berbunyi *"lima sekarang, dua diparkir"*, dan pencabutannya wewenang pemilik
> proses. **Satu kalimat izin cukup.**
