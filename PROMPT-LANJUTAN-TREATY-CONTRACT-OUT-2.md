# PROMPT — LANJUTAN 2 MODUL **Treaty Contract Out** *(sesi yang sama atau sesi baru, folder `OUTPUT_HASIL_RNM\.worktrees\treaty-contract-out`, cabang `modul/treaty-contract-out` @ `8c27b9e` atau lebih baru)*: **menerapkan 14 jawaban work owner atas OQ-TCO-08 … 21**

> Brief modul `PROMPT-IMPLEMENTASI-MODUL-TREATY-CONTRACT-OUT.md` dan lanjutan 1 **tetap berlaku**. Hanya konteks Treaty Contract Out;
> berkas modul lain tidak disentuh; **tidak** merge ke `main` *(asisten menyatukan sesudah verifikasi)*. Setiap pembacaan activity
> mencetak `pyStepsBlockName`.

## 0. KEADAAN AWAL — DIVERIFIKASI ASISTEN 29-09-2026

`8c27b9e`: 17 commit sesudah jeda, pohon bersih; Go **840 PASS · 0 FAIL** *(uji tingkat atas)*, dengan tag `db` **55 SKIP**; vitest
**574**; build **80** modul; vet, gofmt, tsc bersih. Migrasi `300`–`307`. Nol SQL modul ini pernah dijalankan ke Oracle.

## 1. JAWABAN WORK OWNER *(29-09-2026, kutipan)* → yang harus dikerjakan

| OQ | Jawaban | Kerjakan |
| --- | --- | --- |
| **08** token penyimpanan | *"sekarang"* | pasang pelaksana **nyata** efek penyimpanan lampiran di balik env `PELAKSANA_STORAGE=nyata` *(bawaan tetap `stub`)*: token dari `GET_TOKEN_STORAGE` ditiru di Go dengan garam dari env **`STORAGE_TOKEN_SALT`** *(diisi work owner lewat saluran rahasia DBA; **tidak pernah** ditulis ke repo, log, uji, atau dokumen)*; alamat layanan dari resolver `M_LINK_SERVICE` saat jalan *(tidak disalin)*; uji memakai server tiruan lokal. **Tidak** memanggil layanan sungguhan dari uji atau dari sesi ini — penyambungan sungguhan terjadi saat work owner mengisi env dan menyalakan backend |
| **09** pekerja latar | *"perlu"* | pekerja berkala milik modul yang menjalankan `SatuPutaran` *(interval dari env, bawaan mati di uji)*; retry dan anti-dobel memakai kunci idempoten yang sudah ada; dihidupkan dari `cmd/api` |
| **10** tanggal akhir kontrak | *"Ganti jadi mulai + 1 tahun kalender"* | `SetTanggalTreatyContract` langkah 5 diganti: akhir = mulai **+ 1 tahun kalender** *(jepit akhir bulan seperti `ADD_MONTHS(…,12)` untuk 29 Februari)*; **penyimpangan sadar** dicatat bertanggal di tiket 04 dengan bukti b3007 dan anomali 366 hari; uji: tahun biasa, tahun kabisat, 29 Februari |
| **11** satu jenis satu kontrak per tahun | *"benar"* | aturan yang ada **dikonfirmasi**; label `[keputusan kami]` diganti `[keputusan work owner 29-09-2026]` |
| **12** kolom `STATUSACTIVE` master `AGENT` | *"benar"* | konfirmasi; tandai `[keputusan work owner]` |
| **13** nonaktif business | *"0 berarti nonaktif"* | konfirmasi; konstanta di `models`, uji |
| **14** anti-dobel per jenis klausul, wajib Object/Periode | *"setuju"* | konfirmasi |
| **15** daftar `ReinsTypeID` form klausul | *"benar"* | konfirmasi |
| **16** kolom master `TREATYDESC`/`OCCUPATION`/`CLAUSE` | *"benar"* | konfirmasi |
| **17** security dobel ditolak, `%Share` 0..100 | *"setuju"* | konfirmasi; tandai penyimpangan sadar dari Pega |
| **18** `KURS` diisi, skala 8, dua kurs = master rusak | *"setuju"* | konfirmasi |
| **19** tombol simpan tunggal | *"tidak perlu"* | **buang** `POST /kontrak-utuh` dan `PUT /kontrak/{kid}/utuh` beserta handler dan klien yang hanya melayani keduanya *(rute tanpa pemanggil)*; penyimpanan tetap **per panel seperti Pega**; bagian tiket 09 yang masih dipakai panel *(transaksi per simpan)* dipertahankan; tiket 09 diralat bertanggal *(sebagian `wontfix` — keputusan work owner)* |
| **20** "klausul milik kontrak ini" | *"dari induknya"* | angka popup hapus dihitung dari **`PARENTREINSTYPEID`** jenis reasuransi kontrak; uji |
| **21** hapus kontrak | *"hapus saja, samain dengan pega"* | kaskade hapus **seperti Pega**: seluruh anak kombinasi *(reinsurer, security, business, klausul menurut aturan Pega)* ikut terhapus, **termasuk** yang juga dipakai kontrak lain; popup konfirmasi `KonfirmasiHapusTCO` **tetap** dan wajib menyebut **kontrak lain yang ikut terdampak** *(cacahnya)* supaya penghapusan tidak diam; uji mutasi atas perbaikan konkurensi disesuaikan; tiket 10 diralat bertanggal |

Tiap OQ yang dijawab: blok **"Keputusan work owner 29-09-2026"** di tiket terkait dan di dokumen OQ modul; OQ ditutup.

## 2. URUTAN — satu commit per kelompok

| # | Kelompok | Commit |
| ---: | --- | --- |
| 1 | konfirmasi 11–18 *(dokumen + konstanta)* | `docs: treaty-contract-out — OQ-TCO-11..18 dikonfirmasi work owner` |
| 2 | **21** kaskade seperti Pega · **20** induk | `treaty-contract-out: OQ-TCO-20/21 — hapus seperti Pega, klausul dari induk` |
| 3 | **19** buang rute simpan utuh | `treaty-contract-out: OQ-TCO-19 — rute simpan kontrak utuh dibuang` |
| 4 | **10** tanggal akhir + 1 tahun kalender | `treaty-contract-out: OQ-TCO-10 — tanggal akhir mulai + 1 tahun kalender` |
| 5 | **09** pekerja latar · **08** pelaksana penyimpanan nyata | `treaty-contract-out: OQ-TCO-08/09 — pekerja latar dan pelaksana penyimpanan nyata` |
| 6 | uji penuh + `/code-review` singkat atas `8c27b9e..HEAD` | `fix: temuan /code-review — …` bila ada |

## 3. LAPORAN

Satu pesan: tabel **OQ → commit → rute/komponen**; ralat tiket; angka uji tiap commit dengan dan tanpa tag `db`; kunci env baru
*(`PELAKSANA_STORAGE`, `STORAGE_TOKEN_SALT`, interval pekerja)* dan cara mengisinya — **tanpa nilai**; bab **TELEMETRI EKSEKUSI**.

---

*Disusun 29 September 2026 dari jawaban work owner atas OQ-TCO-08 … 21 dan verifikasi `8c27b9e`.*
