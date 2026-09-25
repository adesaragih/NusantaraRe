# Penelusuran tiga pohon cermin — Treaty In Adjustment

**Tanggal:** 24 September 2026 · **Langkah 1 to-spec** · **Perkakas:** `tools/telusur-pohon-cermin.py`

> ## BUKAN SUMBER KEBENARAN
>
> **Struktur kanonik ada di `../../treaty-in/2-to-spec/`.** Berkas ini hanya memuat yang **khas
> addendum** — `GRL-01` mengikat: modul ini **tidak punya spesifikasi sendiri**. Dua salinan penuh
> akan berbeda dalam sebulan, dan pembacanya tidak akan tahu mana yang menang.

> ### SEMESTA — dibaca sebelum satu baris dipercaya
>
> **Jalur yang ditambahkan di sini TIDAK BOLEH dihitung sebagai kelengkapan.**
>
> Peta telusur induk berlabel *"semestanya kurang"* akibat titik buta `L-8`, dan **340 properti
> belum diperiksa siapa pun**. Yang dapat dinyatakan: *seluruh jalur **di dalam semesta yang
> diperiksa** punya nasib.* Bukan: *seluruh jalur punya nasib.*

## 1. Penyebut — dihitung ulang, bukan diwarisi

`struktur-treatyin-lama.md` §5 menyebut **526 dari 985 simpul** berada di dalam keempat pohon cermin
(217 + 152 + 142 + 15). **Itu hitungan SIMPUL POHON.** Penyebut berkas ini dihitung dari **daftar
jalur** `PETA-TELUSUR-JSON.md` §6, dan keduanya **cara menghitung yang berbeda — tidak harus sama**.

| | Jumlah |
|---|---:|
| jalur beradjudikasi §6 | **667** |
| **lingkup modul ini** | **279** |
|   — jalur pohon cermin | 272 |
|   — skalar addendum yang ditarik masuk | **7** |
| di luar lingkup — pohon utama, milik Prompt A | **388** |

**Menutup: 279 + 388 = 667.**

**Ditolak perkakas:** 2 berkas bukan `.xml`. Tidak ada sebab penolakan lain.

## 2. Hitungan penutup — empat nasib, tidak ada yang kelima

| Nasib | Jumlah |
|---|---:|
| **DISIMPAN** | **154** |
| **TURUNAN** | 12 |
| **DIBUANG** | **113** |
| **DITUNDA** | **0** |
| **JUMLAH** | **279** |

> **TIDAK ADA JALUR TANPA NASIB.** `DITUNDA` bernilai nol bukan karena dikosongkan — perkakasnya
> memberi nasib `DITUNDA` kepada apa pun yang tidak dikenalinya, dan tidak satu pun jatuh ke sana.

**Nasib tidak ditentukan perkakas.** Ia datang dari putusan grilling yang dirujuk per baris;
perkakas hanya memasangkan dan memastikan tidak ada yang kosong.

## 3. `SISI` — dihitung dari SELISIH HIMPUNAN BERKAS, bukan dari nama berkas

| | Jumlah |
|---|---:|
| berkas hanya di ekspor Adjustment → **56-KHAS** | **56** |
| berkas di kedua ekspor → **IRISAN** | **323** |

**Per nama daun, klasifikasi ini tidak berguna** — nama seperti `Amount` muncul di mana-mana, dan
258 dari 279 jatuh ke "keduanya". **Maka sisi dihitung per AKAR, dari berkas yang MENULISNYA:**

| Akar | Penulis | 56-KHAS | IRISAN | **SISI** |
|---|---:|---:|---:|---|
| **`ValueDifference`** | 16 | 0 | 16 | **IRISAN** |
| **`OLDDATA`** | 1 | 1 | 0 | **56-KHAS** |
| `ActualValue` | 19 | 1 | 18 | KEDUANYA |
| `ValueBeforeProrate` | 1 | 0 | 1 | IRISAN |
| **`EDMState`** | 2 | 2 | 0 | **56-KHAS** |
| **`EDMMaterialType`** | 1 | 1 | 0 | **56-KHAS** |
| **`AddendumPremi`** | 1 | 1 | 0 | **56-KHAS** |
| `EDMEffective` | 2 | 1 | 1 | KEDUANYA |
| **`EDMDate`** | **0** | 0 | 0 | **TIDAK ADA PENULIS** |
| `RevisionState` | 3 | 0 | 3 | IRISAN |
| `ViewState` | 6 | 1 | 5 | KEDUANYA |
| `IsEditData` | 4 | 1 | 3 | KEDUANYA |
| `OLDID` | 4 | 2 | 2 | KEDUANYA |

### Tiga hal yang terbaca dari tabel itu, dan ketiganya berakibat

**1. Mesin selisih adalah IRISAN — ia milik Treaty In, bukan Adjustment.**
Keenam belas penulis `ValueDifference` **seluruhnya** di berkas yang ada di kedua ekspor. Maka setiap
temuan tentang mesin selisih — `TDA-04` pemadanan posisi, `TDA-05` mata uang tidak dibandingkan,
`TDA-16` selisih EGNPI nol — **milik induk**, diadili di sini untuk jalur addendum saja, dan
**dicatat sebagai usulan untuk induk**.

**2. `OLDDATA` justru 56-KHAS** — satu penulis, dan ia hanya ada di ekspor Adjustment. Pohon yang
memegang nilai lama memang lahir dari jalur addendum.

**3. `EDMDate` NOL PENULIS.** Tidak ada satu pun aturan Pega yang menulisinya. Ia diisi
`POOLDATA.PEGA_M_TREATY_IN_EDM` dengan `SYSDATE` — **pada sisip DAN pada perbarui**. Itu yang
membuatnya **tanggal sentuh terakhir**, bukan tanggal addendum, dan itu sebabnya ia **DIBUANG**.

## 4. Tabel penelusuran

Bentuk barisnya: `JALUR_JSON` · `NASIB` · `TABEL` · `KOLOM` · `TIPE` · `BOLEH_KOSONG` · `ALASAN` ·
`BUKTI` · `SISI`. Isinya dibangkitkan ke `tools/telusur-cermin.json` dan dituangkan di §5.

`BUKTI` berbentuk `EVIDENCED(NamaAturan@ekspor-2026-09)`. Ekspor ini **mungkin** dari lingkungan QA
(`L-9`); penanda itu yang memungkinkan satu sapuan mengeluarkan daftar terdampak bila jawabannya
kelak datang.
