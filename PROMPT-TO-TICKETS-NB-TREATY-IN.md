# PROMPT TO-TICKETS — NB Treaty In

> Salin seluruh isi berkas ini sebagai prompt ke sesi eksekutor.
>
> ⛔ **Prasyarat:** `.scratch\nb-treaty-in\spec.md` sudah ada — 60.167 B, 48 user story, 96
> acceptance criteria — dan `VERIFIKASI-KEADAAN.md` berputusan **layak**.
>
> ⚠️ **Skill `to-tickets` tidak dapat dipanggil sendiri oleh agen** (CLAUDE.md §8). Ketikkan
> `/to-tickets` sebagai manusia, lalu berkas ini menjadi briefnya. Bila skill tidak dipakai, tulis
> tiket langsung dalam bentuk rumah yang dijelaskan di Bab 3.

---

## 0. LINGKUP — DIKUNCI

Hanya modul **`D:\XML\RNM_BRD\NB Treaty In`**.

`D:\XML\RNM_BRD\` adalah korpus **READ-ONLY**. Menulis hanya ke `OUTPUT_HASIL_RNM\`.
`D:\XML\nusantara-re\` **terlarang**.

Keluaran: **`OUTPUT_HASIL_RNM\.scratch\nb-treaty-in\issues\`** — berkas bernomor `NN-judul.md`.

⛔ `spec.md`, `grilling-ronde-*.md`, dan lembar jawaban **tidak disunting** oleh ronde ini.

---

## 1. SUMBER

| # | Sumber | Sifat |
| ---: | --- | --- |
| 1 | `.scratch\nb-treaty-in\spec.md` | **sumber utama** — tiket lahir dari sini |
| 2 | `.scratch\nb-treaty-in\PERTANYAAN-untuk-*.md` | keputusan work owner, mengikat |
| 3 | `.scratch\nb-treaty-in\VERIFIKASI-KEADAAN.md` · `KEADAAN-NB-TREATY-IN.md` | keadaan terukur |
| 4 | `docs\adr\*.md` | 15 ADR; rujuk yang relevan, **jangan buat ADR baru** |
| 5 | `grilling-ronde-1..4.md` | latar; angkanya sebagian basi |

Contoh bentuk rumah: **144 tiket** di sebelas modul, `.scratch\*\issues\*.md`.

---

## 2. ⛔ ATURAN KHUSUS RONDE INI — JANGAN SEBUT NAMA TABEL

**Struktur tabel penyimpanan NB Treaty In belum dirancang**, dan tidak dapat dirancang sampai
contoh isi JSON polis diterima dari DBA (**P29**).

Karena itu tiket ronde ini **tidak boleh menyebut nama tabel atau kolom fisik**. Tulis dalam bahasa
domain: *"data realisasi treaty"*, *"riwayat akseptasi"*, *"putusan per tingkat"* — bukan nama
tabel.

**Ini bukan kekurangan.** Dua modul sudah membuktikan caranya berhasil:

| Modul | Tiket menyebut nama tabel |
| --- | --- |
| `claim-facin` | **0 dari 16** |
| `komite-claim-facin` | **0 dari 14** |
| `claim-prop` | 1 dari 16 |
| `komite-claim-prop` | 1 dari 15 |

Tiga puluh tiket lengkap tanpa satu nama tabel pun.

### Pengecualian tunggal — tiket `00`

`00-skema-penyimpanan-dan-migrasi.md` ditulis sebagai **tiket gantung**: pendek, berstatus
`needs-info`, menyatakan apa yang ditunggu dan apa yang sudah diketahui. **Jangan dikosongkan** —
nomornya dipesan supaya tidak diambil tiket lain.

Isinya cukup: apa yang ditunggu (**P29** contoh JSON, **P1** naskah procedure), apa yang sudah
diketahui (JSON dibuang, penyimpanan jadi tabel datar, satu procedure melayani NB **dan** EDM),
dan pernyataan tegas bahwa **tiket lain tidak bergantung padanya**.

⚠️ Catat juga di tiket itu: `PEGA_JSON_POLIS_TREATYIN` dipanggil NB Treaty In, NB FacIn, **dan**
EDM Treaty In dengan delapan parameter yang sama — EDM mengisi dua parameter yang NB kirim sebagai
`NULL` dan `'0'`. **Rancangan tabel kelak harus menampung keduanya.**

---

## 3. BENTUK TIKET

Kepala berkas:

```
# NN: Judul — <kalimat singkat yang membedakannya>

**Status:** ready-for-agent | needs-info | blocked
**Blocked by:** NN · NN          (kosongkan bila tidak ada)
**Menutup:** AC a · b · c *(n AC)* — US x–y
```

Bab, berurutan. Yang bertanda ★ hampir selalu ada di 144 tiket yang sudah ada:

| Bab | Dipakai |
| --- | ---: |
| ★ `## Hasil & nilai pengguna` | 142/144 |
| ★ `## Area codebase` | 115/144 |
| ★ `## ADR terkait` | 113/144 |
| ★ `## Acceptance criteria` | 113/144 |
| ★ `## Rule Pega sumber` | 109/144 |
| `## Perintah verifikasi` | 75/144 |
| `## Blocker` | 59/144 |
| `## Butir [terbuka] yang menyentuh tiket ini` | 42/144 |
| `## Catatan` | 48/144 |

### `Hasil & nilai pengguna` ditulis dua paragraf

Pertama **hari ini** — apa yang terjadi sekarang, dengan cacatnya bila ada.
Kedua **sesudah tiket ini** — apa yang berubah bagi pengguna.

Bukan daftar tugas. Bahasa pengguna, bukan bahasa Pega.

---

## 4. YANG WAJIB — CAKUPAN AC

**Setiap dari 96 acceptance criteria di `spec.md` harus disebut oleh sekurang-kurangnya satu
tiket.** Tidak boleh ada AC yatim.

Buat tabel silang di berkas `issues\00-PETA-AC.md`: nomor AC → tiket yang menutupnya.

**Cacah dua cara:**
*(a)* kumpulkan seluruh nomor AC yang disebut di baris `**Menutup:**` seluruh tiket;
*(b)* baca bab `Acceptance criteria` tiap tiket dan kumpulkan nomor yang disebut di sana.

Kedua himpunan harus sama, dan gabungannya harus **1–96 tanpa lompat**. Bila tidak, laporkan
nomor yang hilang atau ganda — **jangan ditambal diam-diam**.

---

## 5. URUTAN DAN KETERGANTUNGAN

Tiket yang menyentuh **rantai perhitungan uang** berstatus `blocked`, `Blocked by: P18`. Tulis
tetap — isinya apa yang sudah diketahui — tetapi jangan ditandai `ready-for-agent`.

Sembilan medan **wajib sekaligus terkunci** di layar Dept Head berada di sini: selama P18 kosong,
layar itu tidak dapat disimpan. Tiket layar Dept Head **wajib menyebut ketergantungan itu**.

Sisanya — alur, tahap, wewenang, kronologi, masa berlaku, layar, penggolongan bisnis, riwayat
akseptasi — dapat berstatus `ready-for-agent`.

Perkiraan wajar: **12–18 tiket**, mengikuti sebelas modul lain yang menghasilkan 9–16.

---

## 6. DISIPLIN

Setiap pernyataan membawa bukti: **jalur berkas + tipe rule + nama rule**.

Penanda: `[terverifikasi]` · `[dugaan]` · `[terbuka]` · `[data DBA]` · `[keputusan work owner]` ·
`[penyimpangan sadar]`.

**Jangan menutup pertanyaan terbuka.** **Jangan membuat keputusan baru** — tiket menerjemahkan
spec menjadi pekerjaan, tidak menambah keputusan. Bila sebuah tiket terasa memutuskan sesuatu yang
tidak ada di spec, ia salah tulis: laporkan, jangan dilaksanakan.

**Jangan membuat ADR baru.** Rujuk 15 yang ada. Bila sebuah tiket menuntut ADR baru, **katakan
begitu di bab Catatan** dan biarkan manusia yang memutuskan.

**Jangan menyalin nilai berupa nama orang.** Nomor polis tidak ditulis apa adanya.
Nol rahasia, nol data nasabah, nol cuplikan data produksi. **Nol DDL, nol `CREATE TABLE`.**

Uang tidak pernah `float` — **ADR-0003**. Arah ketergantungan `handlers → services → repository`.

Jalur korpus memuat spasi — pakai Python, bukan loop shell.

---

## 7. BAB WAJIB — TELEMETRI EKSEKUSI

Tulis di `issues\00-PETA-AC.md`, bab terakhir.

Angka token sejati tidak terlihat dari dalam sesi. Ukur dari luar:

```
claude --print --output-format json "<prompt>" > hasil.json
```

Catat token keluaran · cache-read · jumlah panggilan alat · durasi · biaya.
Bila pengukuran dari luar tidak dilakukan, **katakan begitu**, dan bila angka diperoleh dengan
mengurangkan baseline, **katakan itu juga** — jangan menyajikannya seolah terukur langsung.

---

## 8. YANG MENANDAKAN RONDE INI BERHASIL

- Seluruh **96 AC** tertutup, dicacah dua cara, nol yatim dan nol ganda.
- **Nol nama tabel** di seluruh tiket kecuali tiket `00`, dan di sana pun hanya sebagai
  keterangan apa yang ditunggu.
- Tiket yang bergantung **P18** bertanda `blocked`, bukan `ready-for-agent`.
- Tiket `00` ada, pendek, dan menyatakan bahwa tiket lain tidak bergantung padanya.
- `00-PETA-AC.md` memuat tabel silang AC → tiket.
- Nol keputusan baru, nol ADR baru.
- Bab telemetri jujur tentang cara pengukurannya.

---

*Disusun 22 September 2026, dari spec 60.167 B yang lahir sesudah empat ronde grilling, 47 dari 49
pertanyaan terjawab, dan satu perbaikan berkas korpus yang menyusutkan lingkup modul 60 %.*
