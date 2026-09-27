# PROMPT — GILIRAN TIGA MODUL 2 *(sesi tunggal, `main`)*: **paket 0 dokumen bd/018 → §2 A bagian 2 diagnosa (bf) → §2 B dokumen (be) → §2 C A4 kode → §3 PremiumList 01–09 (+ av) → §4 Komite 01–09**

> Brief `PROMPT-IMPLEMENTASI-TIGA-MODUL-GILIRAN-1.md` **tetap berlaku seluruhnya**; yang sudah tuntas dan
> tidak diulang: §5 penyatuan, §1 butir 1–2, §2 A bagian 1 *(migrasi 018)*. Mekanisme giliran = lanjutan 8
> §1; batas tiket = titik lapor sah *(GILIRAN-1 §6)*. Nol pertanyaan "mana yang didahulukan": urutannya §2 di
> bawah, semuanya, berurutan.

---

## 0. KEADAAN SESUDAH `94ee1a6` — DIVERIFIKASI ULANG

| Klaim laporan | Diperiksa ulang | Hasil |
| --- | --- | --- |
| penyatuan; konflik di `migrasi_test.go` dan `strukturkolom_test.go`; tiga cabang sejajar `c1d3b8c` | `df1353d` *(Komite)*, `c1d3b8c` *(PremiumList)*; kedua worktree di `c1d3b8c`; `main` `94ee1a6`; tiga pohon bersih | ✅ |
| Go 334 · 0 · 34; JS 242 | dijalankan ulang: **334 PASS · 0 FAIL · 34 SKIP**; vitest **242**; **50** modul; vet, gofmt, tsc bersih | ✅ |
| migrasi 018 | `T_CLAIMLF_DIAGNOSE` *(`ID`, `PREMIUM_LIST_DETAIL_ID` FK **ON DELETE CASCADE**, `URUTAN`, `ICD_CODE` 100, `DISEASE` 1000, `GROUP_DIAGNOSE` 255, `STS_REJECT` 8 = tipe peserta migrasi 003)*, `IX_DIAGNOSE_PESERTA`, `SEQ_CLAIMLF_DIAGNOSE`, `DISEASE` peserta → 1000; `_down` urut benar | ✅ sesuai bd |
| OQ-K.1 ditutup *(kolom `ID`)*; OQ-K.2 ditutup; OQ-L dibuka | `OQ-untuk-tim.md` 467, 495, 513; `penyakit.go:55` `kolomNomorPenyakit = "ID"` | ✅ |
| `GROUPDIAGNOSE` hanya di satu berkas | seluruh korpus *(19 folder modul)*: hanya `ClaimLifeDetailGCNM.xml`; katalog DEV POOLDATA: **tidak ada** kolom padanannya di tabel mana pun | ✅ OQ-L sah |
| gerbang `.STS_REJECT=='1' \|\| .STS_REJECT=='2'` ×4 di grid | b4682 *(`Add`)*, b5059 *(`Find Disease`)*, b5870 *(`GROUPDIAGNOSE`)*, b6152 *(`Delete`)* | ✅ — dipakai §2 baris 1 |
| `T_MIGRASI` DEV | masih `001`–`016` — **work owner belum menjalankan `-migrate`** | ⚠️ §4 |
| tiket diralat | ⛔ **tiket 08 dan 14 belum punya blok bertanggal bd/018** *(tiket 08 baris 94 masih menyebut `T_CLAIMLF_DIAGNOSE` sebagai "bila work owner menyetujui al")*; STRUKTUR sudah *(bab `T_CLAIMLF_DIAGNOSE`)*; PARITAS 24/25/13c belum menyebut rutenya | ⛔ paket 0 |
| penjaga kaskade | `TestKaskadeHanyaPadaEmpatRelasi` kini memuat lima relasi Claim Life + relasi 9 Komite — **namanya berbohong** | ⚠️ paket 0, kecil |
| nol kebocoran | seluruh berkas yang berubah | ✅ |

## 1. KEPUTUSAN **bf** — `GROUPDIAGNOSE` tanpa daftar pilihan *(OQ-L)* `[DIPUTUSKAN; veto work owner]`

Daftar pilihan dropdown b5863 hidup pada rule properti `GROUPDIAGNOSE` *(kelas `Data-DiagnoseLife`, local
list)* yang **tidak diekspor** dan **tidak ada** di katalog DEV. Maka: kolomnya ada *(018)*; backend `PUT
…/diagnosa/{diagId}` menerima `groupDiagnose` teks ≤ 255 **tanpa** pemeriksaan daftar, satu konstanta
berpenanda `[terbuka — OQ-L]`; frontend menampilkan kontrolnya sebagai `BelumTersedia` bernama dengan label
VERBATIM dari section *(baca `pyLabelPreview` kolom itu di b5780–b5860, catat barisnya)* sampai OQ-L
dijawab. OQ-L ditambah permintaan konkret: *"ekspor `Rule-Obj-Property GROUPDIAGNOSE` kelas
`ASM-FW-GISFW-Data-DiagnoseLife` beserta daftar lokalnya"*. Nilai **tidak dikarang**.

## 2. URUTAN GILIRAN INI — tanpa pesan di antaranya kecuali batas tiket

| # | Paket | Isi | Commit |
| ---: | --- | --- | --- |
| 0 | **Dokumen bd/018 + nama penjaga** | tiket 08: bab bertanggal *"Diagnosa banyak per peserta — bukti XML, 27 September 2026 (butir bd)"* *(b3914/b3923, Add b4690/addRow b4700, Delete b6160/deleteRow b6170, SetDisease b260/b307/b389, SetSTS_Reject b241/b257, gerbang ×4)*; baris 94 diralat; tiket 14: bab *"Migrasi 018 — `T_CLAIMLF_DIAGNOSE` (butir bd)"*; PARITAS 24/25/13c → rute yang akan ada; `TestKaskadeHanyaPadaEmpatRelasi` → nama yang menyebut isinya | `docs: bd/018 di tiket 08 dan 14; PARITAS diagnosa` · `fix: nama penjaga kaskade sesuai isinya` |
| 1 | **§2 A bagian 2 — diagnosa: rute + layar** | `POST /api/klaim-life/{id}/peserta/{pesertaId}/diagnosa` *(Add b4690 → baris kosong, `URUTAN` berikut)*; `PUT …/diagnosa/{diagId}` *(Choose b2509 → `SetDisease`: `ICD_CODE`+`DISEASE` dari baris hasil `GET /api/penyakit-life`; `groupDiagnose` per **bf**)*; `DELETE …/diagnosa/{diagId}` *(Delete b6160; `URUTAN` dirapatkan)*; daftar ikut `GET /api/klaim-life/{id}` per peserta. **Gerbang** *(b4682/b5059/b5870/b6152)*: di dalam baris grid `.` = baris diagnosa *(`STS_REJECT`-nya disalin `SetSTS_Reject`)*, di luar grid `.` = peserta — keduanya berarti **peserta sudah diputus** → ketiga rute menjawab **409** berkata-kata bila `STS_REJECT` peserta ∈ {`1`,`2`}; uji tabel. `SetSTS_Reject` b241/b257 disambungkan ke rute tolak/akseptasi yang **ada**: `STS_REJECT` tiap diagnosa = peserta, satu transaksi. Gerbang tahap/pemegang = gerbang pemuat `ClaimLifeDetailGCNM` *(`EditDateClaimLife_Section`, `RejectOSClaimLife_Sec`, `ViewClaimDetailLifeGCNM`; bukan `MedicalCheckClaimLife`)* — dibaca, dicatat. Layar: grid di Detail dengan `Add`/`Delete`, kolom `DISEASE`/`ICDCODE` read-only *(b5374/b5566)*, `GROUPDIAGNOSE` per bf; `CariDiagnosa.tsx` tombol `Choose` aktif; penanda `KlaimLife.tsx:536` dicabut; label VERBATIM; uji JS | `claim-life: A3 — Diagnosa (bd) bagian 2, rute dan layar` |
| 2 | **§2 B — dokumen (be)** | GILIRAN-1 §2 B utuh: `UNGGAHAN_DIR`; unggah/unduh/hapus lewat outbox dengan pelaksana stub; layar `DocumentLife` tiga tombol aktif; email = efek outbox | `claim-life: A3 — Dokumen (2), unggah/unduh/hapus lewat outbox (be)` |
| 3 | **§2 C — A4 kode** | GILIRAN-1 §2 C; tanpa eksekusi | `claim-life: A4 — migrasi data (kode)` |
| 4 | **§3 PremiumList Life 01 → 09** | GILIRAN-1 §3 utuh; sesudah tiket 04: paket **av** Claim Life | `premiumlist-life: tiket NN — <judul>` · `claim-life: av — PolicyDataLife dari PremiumList Life` |
| 5 | **§4 Komite Claim Life 01 → 09** | GILIRAN-1 §4 utuh | `komite-claim-life: tiket NN — <judul>` |

Seluruhnya di `main`; kedua worktree di-fast-forward pada akhir giliran. Migrasi baru hanya dari keputusan
tercatat: Claim Life `019`+, Komite `031`+, PremiumList `057`+. Tiap paket: bab **Pembacaan ulang XML** dengan
path + baris di tiket terkait; ralat bertanggal bila XML membantah tiket; PARITAS; uji murni + handler + JS;
`LAPORAN-GILIRAN-F0.md` +1 bab; nol kebocoran; `-migrate` **tidak** dijalankan executor.

## 3. YANG TIDAK BERUBAH

Semua keputusan GILIRAN-1 §1 *(OQ-K.1 = `ID`; nol tabrakan nama; stub, bukan izin; A4 kode saja; urutan
modul)* dan keputusan bd/be. Induk §0.3: brief modul yang tertempel di sesi ini dikerjakan sesi ini.

## 4. WORK OWNER — sekarang, sebelum giliran ini mulai

`T_MIGRASI` DEV masih berhenti di `016`. Jalankan **sekali dari `main`** *(perintah induk §0)*:
`. .\muat-env.ps1` → `go run .\cmd\api -migrate` → mengisi `017`, `018`, `030`, `050`–`056`. Tanpa ini,
layar Close Claim dan diagnosa akan ditolak Oracle *(kolom dan tabelnya belum ada di DEV)*.

## 5. LAPORAN · TELEMETRI

Persis GILIRAN-1 §6: satu pesan di akhir atau pada batas tiket; tabel per modul **tiket → commit → menu/tombol
XML → rute/kontrol**; ralat tiket; OQ ditutup/dibuka; angka uji tiap commit; bab **TELEMETRI EKSEKUSI** per paket.

---

*Disusun 27 September 2026 sesudah verifikasi `94ee1a6` (uji, vet, gofmt, tsc, build dijalankan ulang;
graf merge dibaca), pembacaan ulang migrasi 018 + down, `ClaimLifeDetailGCNM.xml` b3900–b6250 (gerbang
×4, dropdown b5863), pencarian `GROUPDIAGNOSE` atas seluruh korpus, katalog DEV (`T_MIGRASI`, kolom
GROUP/DIAG — agregat saja), OQ-untuk-tim, tiket 08/14, STRUKTUR, PARITAS, `penyakit.go`, `migrasi_test.go`.*
