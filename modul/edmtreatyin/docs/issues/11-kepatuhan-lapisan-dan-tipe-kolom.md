# 11: Kepatuhan lapisan dan tipe kolom

> ## ⭐ PENAHAN GUGUR — 23 September 2026 sore
>
> `[keputusan work owner]` *"Selesaikan, jangan jadi permasalahan."* ⭐ Presisi dinaikkan ke ~~**`NUMBER(38,8)`** — **30 digit di depan koma**, delapan di belakang, batas tertinggi Oracle~~ ⭐ **`NUMBER(38,10)`** *(RALAT 04-10-2026 — lihat blok KOREKSI 06-10 di bawah)*. Dasarnya sapuan korpus: ambang dagang nyata sudah **tepat di batas** 12 digit *(`181500000000.00` · `150000000000.00`)*, dan ada sentinel **17 digit** *(`99999999999999999.99`)*. Penjumlahan lintas mata uang dapat melewati keduanya. ⭐ `NUMBER` di Oracle berpanjang **berubah-ubah** — hanya digit bermakna yang tersimpan, sehingga pelebaran ini **tidak memakan ruang tambahan**. Butir ditutup oleh bukti, bukan oleh DBA.
>
> ⭐ **`blocked` → `ready-for-agent`.** ⛔ **Nol butir `[data DBA]` tersisa di tiket ini.**
>
> Rinciannya: `modul/nbtreatyin/docs/KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`.

> ## ⛔ KOREKSI 06-10-2026 — tabel NB nyata (log: `../KOREKSI-DOKUMEN-2026-10-06.md`)
>
> | Bunyi lama (dikutip) | Bunyi baru | Bukti |
> | --- | --- | --- |
> | *"Presisi dinaikkan ke **`NUMBER(38,8)`** — **30 digit di depan koma**, delapan di belakang"* | ⭐ **`NUMBER(38,10)`** — **28 digit** di depan koma, **sepuluh** di belakang (RALAT 04-10-2026, perintah work owner). Galat lama 9 desimal (`592.629.512,880000276`) **tersimpan utuh**; yang berdesimal > 10 (`ShareValue` 24, `PremiumSpreaded` 20) dibulatkan di desimal kesebelas | `modul/nbtreatyin/backend/migrations/320_t_general_polis_treaty.sql` baris 26–28; `modul/nbtreatyin/docs/spec-penyimpanan-relasional.md` blok RALAT 04-10-2026 (baris 5–17); EDM 360–363 `NUMBER(38,10)` (`../STRUKTUR-TABEL-EDM-TREATY-IN.md` baris 8–10) |
> | *"spec ini **tertinggal**"* (AC 49 sembilan lawan AC 58 delapan) | ✅ Pertentangan **sudah diselaraskan** di spec 23-09 sore (AC 49 dikutip-coret); dengan skala 10, tuntutan lama *"sekurangnya sembilan"* kini **terpenuhi** | `../spec-penyimpanan-relasional.md` AC 49 blok *DISELARASKAN 23-09-2026 sore* |
> | *"Menutup: AC 45–53 · AC 57–58 (11 AC)"* | **AC 45–53 · AC 58 (10 AC)**; AC 57 (`REMARK`, `ID-27b`) milik **tiket 12** | `00-PETA-AC.md` blok KOREKSI butir 2 |

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026 sore)*
~~**Blocked by:** ⛔ `[data DBA]` **presisi fisik kolom uang** — dua belas digit di depan koma belum diuji terhadap nilai terbesar~~ ⛔ **gugur 23-09-2026 sore**
**Bergantung pada tiket NB:** **18** *(tipe kolom dan presisi uang)* · **20** *(transaksi tunggal dan skema eksplisit)*
**Menutup:** AC **45–53** · AC ~~**57–58** *(11 AC)*~~ **58** *(10 AC — AC 57 ke tiket 12; koreksi 06-10)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-1 · ID-2 · ID-43..ID-47

## Hasil & nilai pengguna

Sembilan aturan kepatuhan yang berlaku sama seperti pada polis baru, ditegakkan juga pada jalur
endorsemen — supaya jalur kedua tidak menjadi celah.

⭐ **Yang paling mudah terlewat:** satu antarmuka penyimpanan, bukan dua. Endorsemen yang mendapat
antarmukanya sendiri akan melewati penjaga yang dibangun untuk polis baru.

## Yang dibangun

| Aturan | Yang dicegahnya |
| --- | --- |
| **satu transaksi** per generasi | generasi separuh tersimpan |
| skema basis data **eksplisit** | tabel terbaca dari skema yang salah |
| kolom pelaku dari **identitas login** | kolom pelaku kosong seperti di sistem lama |
| ⛔ nol tipe mengambang pada uang | sistem baru **menambah** galat baru |
| skala kolom uang | pemotongan digit yang dianggap presisi |
| kode dan penanda **tetap teks** | nol di depan hilang ⇒ penggolong jenis usaha gagal |
| teks kosong → **tak-bernilai** | kosong menjadi nol, dan nol punya arti lain |
| arah ketergantungan tidak terbalik | lapisan penyimpanan memanggil lapisan layanan |
| **satu** antarmuka penyimpanan | jalur endorsemen melewati penjaga jalur polis baru |

Ditambah panjang minimum kolom keterangan, dan aturan bahwa nilai berdesimal lebih dari skala kolom
**dibulatkan, bukan ditolak**.

## Batas — yang TIDAK termasuk

⛔ Penetapan tipe kolom itu sendiri — tiket NB **18**; tiket ini **menegakkannya** pada jalur endorsemen.
⛔ `CREATE TABLE` — presisi fisik dicocokkan DBA **di dalam** tiket ini.

## Cara mengujinya

Lewat seam `repository`, dan ⚠️ **sebagian test memeriksa nilai kolom langsung** — pulang-pergi
saja tidak cukup, sebab tulis dan baca yang sama-sama salah simetris tetap hijau.

⭐ Uji antarmuka tunggal dijalankan sebagai pemeriksaan susunan kode, bukan lewat data.

## Acceptance criteria

- [ ] **AC 45** — seluruh penyimpanan satu generasi berada dalam **satu transaksi**
- [ ] **AC 46** — setiap query menulis skema **eksplisit**
- [ ] **AC 47** — kolom pelaku diisi dari identitas login
- [ ] **AC 48** — nol kolom uang bertipe mengambang
- [ ] **AC 49** — kolom uang menerima skala desimal yang ditetapkan *(koreksi 06-10: sepuluh, `NUMBER(38,10)`)*
- [ ] **AC 50** — kode dan penanda tersimpan sebagai **teks**, nol di depan utuh
- [ ] **AC 51** — teks kosong pada kolom angka atau tanggal tersimpan **tak-bernilai**
- [ ] **AC 52** — arah ketergantungan **tidak pernah dibalik**
- [ ] **AC 53** — hanya ada **satu** antarmuka penyimpanan polis treaty
- [ ] ~~**AC 57** — kolom keterangan berpanjang sekurangnya batas yang ditetapkan~~ ⛔ *(koreksi 06-10)* pindah ke **tiket 12**
- [ ] **AC 58** — nilai berdesimal lebih dari skala kolom **dibulatkan, bukan ditolak** *(skala 10 — koreksi 06-10)*

## ⛔ Kenapa tiket ini `blocked`

~~`[data DBA]` **Dua belas digit di depan koma belum diuji.**~~ ⛔ **BUTIR GUGUR 23-09-2026 sore.** Presisi dinaikkan ke ~~**`NUMBER(38,8)`** — tiga puluh digit di depan koma, delapan di belakang~~ ⭐ **`NUMBER(38,10)`** — dua puluh delapan digit di depan koma, sepuluh di belakang *(RALAT 04-10-2026; koreksi 06-10)*. Ditutup oleh bukti korpus, bukan oleh jawaban DBA.

⚠️ **Dan satu pertentangan di dalam spec, dicatat di sini karena spec tidak boleh disunting:**
✅ *(koreksi 06-10)* **sudah diselaraskan** di spec 23-09 sore — paragraf di bawah tinggal sejarah.
`[terverifikasi]` **AC 49** berbunyi *"kolom uang menerima **sekurangnya sembilan** angka di belakang
koma"*, sementara **AC 58** pada spec yang sama menetapkan skala **delapan** dan menyatakan
kelebihannya **dibulatkan**.

⛔ Nilai produksi berdesimal sembilan **berubah** pada skala delapan. Spec penyimpanan polis baru
sudah diselaraskan ke delapan dan **mengutip bunyi lamanya**; ~~spec ini **tertinggal**~~ spec ini
sudah menyusul 23-09 sore. ~~⛔ **Tidak diputuskan di tiket ini.** Pemiliknya `[work owner]` dan `[data DBA]`.~~
⭐ *(koreksi 06-10)* Diputuskan: skala **10** (`[perintah work owner]` 04-10-2026, RALAT NB) — nilai
berdesimal sembilan **tidak berubah** lagi.