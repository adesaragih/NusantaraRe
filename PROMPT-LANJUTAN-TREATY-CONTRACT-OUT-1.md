# PROMPT — LANJUTAN 1 MODUL **Treaty Contract Out** *(sesi BARU, folder `OUTPUT_HASIL_RNM\.worktrees\treaty-contract-out`, cabang `modul/treaty-contract-out` @ `0022865`)*: **tiket 12 → 04 → 05 → 07 → 08 → 06 → 11 → 09 → 10**, lalu uji penuh dan `/code-review`

> Buka **sesi Claude Code baru** di folder worktree di atas, `/mattpocock-skills:implement`, lalu tempel brief ini utuh.
> Brief induk modul ini **tetap berlaku seluruhnya**: `PROMPT-IMPLEMENTASI-MODUL-TREATY-CONTRACT-OUT.md` *(keputusan tco1–tco3,
> o, lampiran, kurs; peta harness; menu; pembagian berkas)*. Brief ini hanya menetapkan titik lanjut dan urutan sisa.
> Sesi ini **tidak** menyentuh berkas modul lain dan **tidak** merge ke `main`.

## 0. LANGKAH 0 — sebelum menulis kode apa pun

1. `git status --porcelain` harus **kosong**; `git log --oneline -5` harus memuat `0022865` *(bab jeda)* dan `678fd25` *(tiket 03)*.
2. Baca bab **"Jeda — 28-09-2026"** di `.scratch/treaty-contract-out/LAPORAN-GILIRAN.md` **utuh**. Rancangan tiket 12 yang sudah
   tercatat di sana dipakai, tidak dirancang ulang: tabel lampiran baru migrasi **307**; pekerja antrean **milik modul sendiri**
   *(pekerja bawaan menulis jalur gagalnya ke jejak audit Claim Life)*; kunci penyimpanan tetap per lampiran agar pengulangan tidak
   menggandakan berkas; label panel yang sudah dicek ke nomor baris korpus.
3. Uji hijau sebagai titik awal: `go vet ./...`, `go test ./...` **dan** `go test -tags=db ./...`, `gofmt -l`, `npx vitest run`,
   `npx tsc --noEmit`, `npx vite build`. Catat angkanya di bab pertama laporan giliran ini.

## 1. KEADAAN YANG SUDAH PASTI

| Hal | Keadaan |
| --- | --- |
| Tiket selesai | **01** skema + migrasi `300`+ *(`1872d26`)* · **02** jenis reasuransi *(`1f2aab5`)* · **03** tahun treaty *(`678fd25`)* |
| Tiket dijeda | **12** lampiran — baru tahap membaca korpus dan jalur dokumen `main`; **nol** berkas kode ditulis |
| Port | backend `:8093`, Vite `5185` *(`.env` worktree)* |
| Migrasi | `300`–`306` terpakai sampai tiket 03; tiket 12 = `307`; sisa rentang hingga `319` |
| Di luar modul ini | tautan unduh dokumen **Claim Life** 401 *(ditemukan sesi ini, dibuktikan asisten)* — **milik sesi `main`**, **jangan** diperbaiki di sini |

## 2. URUTAN SISA — satu commit per tiket `treaty-contract-out: tiket NN — <judul>`, tanpa pesan di antaranya

| # | Tiket | Bergantung pada | Catatan |
| ---: | --- | --- | --- |
| 1 | **12** lampiran di tahun treaty | 03 ✅ | rancangan bab Jeda; flow action `TreatyOutAttachContent`; outbox + pelaksana stub; endpoint penyimpanan nyata **tidak** dipanggil; unduhan lewat klien `fetch` dengan header identitas *(jangan ulangi cacat `<a href>` Claim Life)* |
| 2 | **04** kontrak treaty di dalam tahun | 02 ✅, 03 ✅ | `InboxTreatyContract` → `GridTreatyContract`, `InputTreatyContract` |
| 3 | **05** reinsurer + total share | 04 | uang `apd.Decimal`; total share per aturan XML |
| 4 | **07** business + nonaktif | 04 | |
| 5 | **08** klausul satu tabel, 25 jenis, validasinya | 04 | `InboxTreatyContractDescription` → `ViewDetailDescription` *(`BrowseDescriptionLimit` b5049/b5457/b8150, `testingKurs` b5110, `SetKirimIDDesc` b5154, `PanggilID` b5228, `RefreshErrorProportionalarrg` b5348/b5730)*; `WHEN` dibaca dari urutan aksi |
| 6 | **06** security reinsurer | 05 | struktur bersih `MTREATYSECURITY` *(tco1: tabel `T_`)* |
| 7 | **11** kurs USD ke IDR | 08 | `TREATYEXCHANGEYEARLY` dibaca saja; `NitipKurs` |
| 8 | **09** simpan atomik lintas enam tabel | 05, 06, 07, 08 | satu transaksi; nol `COMMIT` di teks SQL |
| 9 | **10** kaskade hapus, popup, klausul yang tetap hidup | 06, 07, 09 | `DetailTreatyExclustion` → `DetailTreatyExclustion_Sec` |
| 10 | **Uji penuh + `/code-review` dua sumbu** atas `251cb3b..HEAD` | semua | temuan diperbaiki dalam commit `fix: temuan /code-review — …`; yang dibiarkan ditulis beserta sebabnya |

Tiap tiket: bab **Pembacaan ulang XML** *(path + baris)*; ralat bertanggal bila XML membantah tiket; penulis medan dicari, bukan
hanya pembacanya; uji murni + handler + JS, dibuktikan merah lebih dulu untuk aturan uang; `PARITAS-LAYAR-DAN-AKSI.md` dan
`LAPORAN-GILIRAN.md` +1 bab per tiket; penjaga kata cadangan Oracle dan penjaga nama orang hijau; kode bersama hanya ditambah.

## 3. BERHENTI · JEDA · LAPORAN

- **Selesai** = tiket 10 dan `/code-review` tuntas. Laporan satu pesan: tabel **tiket → commit → harness/tombol XML → rute/komponen**;
  ralat tiket; OQ dibuka/ditutup; kontrak hilir *(Claim Prop, Komite Claim Prop, Claim Fac In)*; kode bersama yang disentuh; angka uji
  tiap commit **dengan dan tanpa tag `db`**; bab **TELEMETRI EKSEKUSI**.
- **Konteks menipis** = berhenti hanya pada batas tiket, baris pertama *"konteks menipis: kira-kira N% tersisa"*, commit bab jeda baru
  di `LAPORAN-GILIRAN.md` *(tiket terakhir selesai, tiket berikutnya, rancangan yang sudah terbaca)*. Pohon **harus** bersih sebelum
  berhenti — bab jeda ikut di-commit.
- **Diminta jeda oleh work owner** = sama dengan konteks menipis.

---

*Disusun 28 September 2026 sesudah verifikasi laporan jeda (`678fd25`, bab jeda di-commit asisten sebagai `0022865`, pohon
worktree bersih) dan pembuktian cacat tautan unduh Claim Life dari kode (`PanelDokumenPeserta.tsx:157`, `pelaku.go`).*
