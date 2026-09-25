# Usulan aturan `TA-12` — penanda keadaan diperiksa terhadap **sumbu** yang dimaksudnya

**Tanggal:** 24 September 2026 · **Lahir di:** putaran penulisan tiket batch 1 Treaty In
**Status: USULAN.** Belum berlaku; menunggu pemilik proses.
**Ditulis di sini, bukan di `GRILL-*/07-AUDIT.md`,** sebab ronde ini **bukan ronde grilling** —
membuatkannya folder ronde berarti mengarang ronde yang tidak pernah dijalankan.

---

## Bunyinya

> **Setiap penanda keadaan diperiksa terhadap SUMBU yang dimaksudnya, bukan hanya terhadap
> nilainya.**
>
> Dua berkas yang menandai hal yang sama dengan kata yang berbeda **belum tentu bertentangan** — dan
> justru itu bahayanya. Bila keduanya memakai **sumbu yang berbeda**, keduanya dapat **benar**
> sekaligus **menyesatkan**, dan pembacanya tidak punya cara mengetahuinya.

---

## Apa yang melahirkannya

Tiga kemampuan — `P-20`, `P-41`, `P-49` — ditandai dua kali dengan kata yang berbeda:

| Berkas | Penandanya | Sumbu yang dimaksudnya |
|---|---|---|
| `DAFTAR-PEKERJAAN.md` §4 | **`TERTAHAN`** | *"akan dikerjakan **begitu penghalangnya jatuh**"* |
| `KEPUTUSAN-PEMBAGIAN-TIKET.md` §2 | **tegas tidak masuk** | *"**tidak dikerjakan** di penyerahan ini, **apa pun jawabannya**"* |

Keduanya **benar** pada sumbunya masing-masing. Dan keduanya bersama-sama menghasilkan kekeliruan
yang dapat diperkirakan:

> **Siapa pun yang mencabut penghalang `P-41` akan mengira ia jadi dapat dikerjakan** — lalu
> mengerjakan sesuatu yang `G2` sudah keluarkan dari ruang lingkup.

Pemeriksaan biasa **tidak menangkapnya**: mencocokkan *nilai* penanda akan melaporkan "keduanya
berisi sesuatu", dan mencocokkan *kata* akan melaporkan "keduanya berbeda" tanpa dapat menyatakan
mana yang benar.

---

## Kenapa ia layak menjadi aturan, bukan catatan

**Ia instans keempat dalam satu hari.** Keempatnya berbentuk sama — dua berkas menyatakan keadaan
yang sama dengan kata yang berbeda — dan keempatnya baru ketahuan ketika seseorang kebetulan membaca
keduanya berdampingan:

| # | Yang bertentangan | Sumbu yang berbeda |
|---|---|---|
| 1 | `2-to-spec/TDA-17-PUTUSAN.md` menyatakan `PREMIUM_EARNED` **DISIMPAN**; `PENELUSURAN-JSON-KE-KOLOM.md` dan §10.3a menyatakannya **DITURUNKAN** | **bentuk** versus **nasib**. Bentuknya `Page List`; nasibnya turunan. Keduanya benar |
| 2 | `issues/README.md` mencetak **"dapat dimulai: 4"**; kenyataannya nol | **penghalang di dalam papan** versus **seluruh penghalang**. Angkanya benar untuk sumbu pertama |
| 3 | `DAFTAR-PEKERJAAN.md` §2 menulis **60/57**; §4 menulis **58/55**; hitungan mekanis memberi **66/63** | **kode berbaris** versus **kemampuan**. `P-59` ada sebagai kemampuan dan tidak ada sebagai baris |
| 4 | `P-20`, `P-41`, `P-49` — **`TERTAHAN`** versus **tegas tidak masuk** | **menunggu jawaban** versus **di luar ruang lingkup** |

**Empat kali dalam satu hari bukan kebetulan; ia pola.** Dan ketiga instans pertama **tidak
ditemukan oleh pemeriksaan** — ketiganya ditemukan karena satu sesi kebetulan membuka kedua
berkasnya.

---

## Cara menerapkannya

**Satu pertanyaan, diajukan sebelum sebuah penanda ditulis maupun dibaca:**

> **Penanda ini menjawab pertanyaan apa?**

Lalu tulis **pertanyaannya**, bukan hanya nilainya. Sebuah penanda yang tidak dapat menyebutkan
pertanyaannya adalah penanda yang **akan dibaca sebagai jawaban atas pertanyaan yang lain**.

**Bentuk yang diterima:**

| Bentuk | Contoh |
|---|---|
| ✅ | `TERTAHAN` — *menunggu jawaban `DB-20`; dikerjakan begitu jawabannya turun* |
| ✅ | `TIDAK MASUK (G2)` — *di luar ruang lingkup penyerahan pertama, apa pun jawabannya* |
| ❌ | `TERTAHAN` sendirian, pada sesuatu yang sebenarnya di luar ruang lingkup |

**Dan pada pemeriksaan dua arah antar berkas**, yang dicocokkan **bukan nilainya** melainkan
**pasangan (sumbu, nilai)**. Dua penanda bersumbu berbeda **tidak dilaporkan sebagai pertentangan**;
ia dilaporkan sebagai **dua fakta yang harus dibaca bersama**, dan berkas yang lebih rendah
wewenangnya menyebut yang lebih tinggi.

---

## Yang aturan ini TIDAK katakan

- Ia **tidak** menuntut satu berkas menjadi satu-satunya tempat penanda. Urutan wewenang sudah
  mengatur siapa yang menang; aturan ini soal **dapat dibacanya**, bukan soal siapa pemiliknya.
- Ia **tidak** berlaku untuk penanda bernilai tunggal yang sumbunya hanya satu — `status: aktif` pada
  tiket punya satu sumbu dan tidak dapat disalahbaca.
- Ia **tidak** menggantikan aturan *"empat golongan, tidak ada yang kelima"*. Justru sebaliknya:
  keempat golongan itu **satu sumbu**, dan aturan ini melarang golongan dari **sumbu lain**
  diselundupkan ke dalamnya.

---

## Bila usulan ini diterima

| | |
|---|---|
| **Nomornya** | `TA-12`, melanjutkan `TA-11` di `treaty-in-adjustment/GRILL-D/07-AUDIT.md` §3 |
| **Berlaku untuk** | setiap modul di `D:\XML_NURE`, sama seperti `TA-09`…`TA-11` |
| **Yang pertama menerapkannya** | putaran ini sendiri — ketiga penanda `P-20`/`P-41`/`P-49` sudah diperbaiki di `DAFTAR-PEKERJAAN.md` §4.1, dan keempat instans di atas sudah tercatat masing-masing di tempatnya |
| **Yang menagihnya bila ditolak** | berkas ini, dan ia **tidak dihapus** — usulan aturan yang ditolak tetap merekam bahwa polanya pernah muncul empat kali dalam satu hari |
