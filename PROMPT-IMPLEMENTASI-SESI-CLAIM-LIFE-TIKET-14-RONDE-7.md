# PROMPT — sesi implementasi berikutnya: Claim Life tiket 14 **ronde 7** — mata uang header, satu ralat nama test, dan penutupan bila Oracle ada

> **Skill:** ketik `/mattpocock-skills:implement` sebagai manusia, lalu tempel berkas ini **utuh**.
> Brief induk **`PROMPT-IMPLEMENTASI-GO-REACT.md`** berlaku seluruhnya. Brief ronde 6
> **`PROMPT-IMPLEMENTASI-SESI-CLAIM-LIFE-TIKET-14-RONDE-6.md`** dan ronde 5 tetap berlaku untuk hal
> yang tidak diubah di sini *(Langkah A dan D, §5–§8)*.
>
> **Urutan baca sebelum satu baris kode:** brief induk §4–§8 → `APP_RNM\README-BACA-DULU.md` →
> tiket 14 bab `## Implementasi — ronde 6, 26 September 2026 malam` → §2–§3 berkas ini.
>
> **SESI INI:** claim-life · tiket **14 ronde 7** + tiket **01 (penutupan, bila G1)** + **persiapan
> tiket 02 (bila G3)**. Tanpa G1 dan tanpa keputusan **z**, sesi ini hanya satu ralat kata (§3-1):
> **jangan mulai untuk itu saja.**

---

## 0. KEADAAN AWAL — 26 September 2026 malam

| | Keadaan |
| --- | --- |
| `HEAD` | `347ee57` — tiket 14 *claimed* **40/53**, daftar terbuka 13 = kotak `[ ]`; tiket 01 *claimed* 4/7; tiket 02 *ready-for-agent*, blocker o. Working tree bersih |
| Uji yang lulus | `go vet`, `go vet -tags=db`, `gofmt` nol, `go build`; **99** test Go PASS; **20** test bertag `db` SKIP dengan pesan; `tsc --noEmit`; **5** test JS; `vite build` 87 modul |
| Diterapkan ronde 6 | **s1** `PeriksaNilaiWarisan` di **kedua** jalur tulis warisan *(`Simpan`, `skemauji.IsiBarisLama`)*, `models` utuh; **w2** `CLAIM_RETRO NUMBER(38,8)` di `002`, `models.Klaim.ClaimRetro Money`, `AmbilHeader` lewat `TO_CHAR`, `MarshalJSON` membawanya sebagai teks, `BarisLamaDari` menulis ke tiap baris, `BongkarBarisLama` membacanya; mata uang campur dilaporkan; pola angka menerima keluaran TM9 *(`.5`)*; penjaga nama `DOCUMENT_CLAIM` |
| Migrasi | **belum pernah dijalankan di Oracle mana pun**, tujuh sesi |
| Ditahan executor, menunggu Anda | **mata uang header** *(`T_GENERAL_CLAIM` tanpa kolom mata uang; `ClaimRetro.Currency` kosong saat header dibaca sendirian)* |

**Verifikasi independen 26 September 2026 malam atas `347ee57`:** seluruh angka laporan ronde 6
tereproduksi *(99 · 20 · 5 · 40/53 · 13 · 18 berkas +634/−38)*; ketiga temuan tinjauan
*(`MarshalJSON`, jalur tulis kedua, TM9)* dibaca di diff dan benar diperbaiki. Dua hal dicek
sendiri karena keduanya menyangkut **seluruh** jalur uang, bukan hanya ronde ini: **TM9 memang
mengeluarkan `.5` dan `-.25`** *(dicoba di instance pengembangan)*, dan **`apd` memang menerima
`.5`, `-.25`, `5.`** *(dicoba dengan pustaka yang sama)* — jadi pembacaan uang tiket 01 sampai
sekarang aman, bukan hanya pagar warisan yang baru. Yang belum tepat ada di §3: satu ralat kata.

---

## 1. GERBANG

| Gerbang | Yang membukanya | Keadaan |
| --- | --- | --- |
| **G1 Oracle** | user **kosong** baru + `ORACLE_DSN` + `ORACLE_SCHEMA` + `ORACLE_SKEMA_UJI=true` | tertutup — DBA belum membuat user; **tujuh** sesi |
| **G2 keputusan tiket 14** | **z** *(mata uang header; bukti §2)*; y, j bila ada | tertutup |
| **G3 keputusan tiket 02** | **o1, o2, o3** | tertutup — hanya oleh keputusan |

---

## 2. KEPUTUSAN WORK OWNER — dengan bukti

| | Keputusan | Bukti | Keadaan |
| ---: | --- | --- | --- |
| z | **Mata uang header klaim.** `T_GENERAL_CLAIM` tidak punya kolom mata uang, padahal `CLAIM_RETRO` kini uang *(w2)*. Usulan: **(z1)** tambah kolom `CURRENCY VARCHAR2(8)` ke `002` *(nama sama dengan kolom warisan dan kolom peserta/adjustment yang sudah ada, supaya tidak lahir dua nama)*; `BongkarBarisLama` mengisinya dari baris — bila campur, tetap dilaporkan seperti sekarang dan header memakai baris pertama; `AmbilHeader` membacanya dan mengisi `ClaimRetro.Currency`; `Simpan` menulisnya; STRUKTUR diberi ralat bertanggal; spec `PremiumListSummary` tidak berubah *(kolom ini turunan dari baris, bukan field Pega baru)*. **(z2)** tanpa kolom: header selalu meminjam mata uang dari peserta pertama saat dibaca — lebih sedikit skema, tetapi `AmbilHeader` sendirian tetap kosong | agregat instance pengembangan *(cacah saja)*: **2** mata uang berbeda di seluruh tabel, **2.856** klaim, **0** klaim bermata-uang campur, **0** baris `CURRENCY` NULL. Satu mata uang per klaim adalah **kenyataan data**, bukan asumsi | `[USULAN]` — rekomendasi **z1** |
| y, j, o1–o3, v2, s′ | tetap sebagaimana brief ronde 5 dan 6 | | `[USULAN]` / `[terbuka]` |
| k, l | tetap | | berlaku |

---

## 3. TEMUAN VERIFIKASI RONDE 6 — satu ralat kata

| # | Temuan | Letak | Yang dikerjakan |
| ---: | --- | --- | --- |
| 1 | ⚠️ Bab ronde 6 menyebut **`TestSimpanMemanggilPagarNilaiWarisan`**, test yang **tidak ada** di kode: sesudah tinjauan ia menjadi `TestSetiapJalurTulisWarisanDipagari` *(cakupannya dua jalur, bukan satu)*, tetapi paragraf di bab tiket tidak ikut diganti | tiket 14 bab ronde 6, paragraf "Satu cacat saya sendiri tertangkap" | satu baris ralat di bab ronde 7 |
| 2 | ℹ️ Penjaga `TestSetiapJalurTulisWarisanDipagari` membaca **teks**: `_ = PeriksaNilaiWarisan(b)` akan lolos. Sudah dinyatakan sendiri di komentarnya | — | catatan; tidak diubah |

Yang **sudah benar**: pagar s1 di dua jalur dan cacahnya dikunci; `ClaimRetro` pulang-pergi sebagai
uang dan ikut JSON sebagai teks; TM9 dan `apd` saling cocok *(dicek ulang)*; mata uang campur
dilaporkan, bukan dipilih diam-diam; nama `DOCUMENT_CLAIM` hanya untuk warisan dan kelas Pega, dengan
penjaga; ralat 89→91; STRUKTUR `CLAIM_RETRO` desimal + blok ralat; tolakan atas *"002 tanpa ALTER"*
benar *(migrasi belum pernah jalan; §2 l)*; pengakuan lingkup melebar ditulis terang.

---

## 4. URUTAN SESI

**Langkah 0** — commit brief ini *(`docs: brief sesi tiket 14 ronde 7`)*, `git status --porcelain`
kosong, uji tanpa Oracle hijau.

**Langkah B — G0.** §3-1 saja.

**Langkah C — G2 bila terbuka.** **z1** → `002` `CURRENCY VARCHAR2(8)` *(hanya selama `T_MIGRASI`
belum ada di instance mana pun, §2 l)*; `models.Klaim` tidak perlu medan baru — `ClaimRetro.Currency`
yang diisi; `AmbilHeader`, `Simpan`, `BongkarBarisLama`, `BarisLamaDari` mengikuti; test murni
pulang-pergi mata uang header; `TestKolomDDLCocokDenganStruktur` menuntut STRUKTUR ikut diralat
*(blok bertanggal)*; `TestAC07HeaderMemuatFieldSummary` **tidak** diubah *(kolom ini bukan field
`PremiumListSummary`)*.

**Langkah A — G1 bila terbuka.** Persis brief ronde 6 §4 Langkah A.

**Langkah D — G3 bila terbuka.** Persis brief ronde 5 §4 Langkah D.

**Verifikasi penuh, penutup, tanda `resolved`** — persis brief ronde 5 §4.

---

## 5–8. DBA · GAYA · PERSETUJUAN · TELEMETRI

Persis brief ronde 6 §5–§8. Satu-satunya yang membuka jalan ke `resolved` masih **user kosong dari
DBA**; tujuh sesi berakhir `claimed` karena itu.

---

*Disusun 26 September 2026 malam dari verifikasi independen `347ee57`: 99 test dijalankan ulang,
diff 18 berkas dibaca utuh, bentuk keluaran TM9 dan pembacaan `apd` dicoba sendiri, agregat mata
uang per klaim dibaca dari instance pengembangan (cacah saja, nol baris).*
