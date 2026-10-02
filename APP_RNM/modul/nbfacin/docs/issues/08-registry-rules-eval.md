# 08: Registry predikat — satu seam untuk seluruh rule `When`

**What to build:** Seluruh gerbang keputusan sistem lama hidup di **satu registry** yang dapat
ditanyai dengan nama, dan setiap predikat menyebut rule Pega asalnya. Tidak ada salinan logika
gerbang yang tersebar dan bisa menyimpang.

⚠️ **Kondisi rule `When` ada di DUA tag, dan keduanya wajib dibaca.** Membaca hanya yang pertama
membuat rule tampak "tanpa kondisi" padahal kondisinya ada — **170 berkas (28,3 %)** menyembunyikan
kondisinya di tag kedua: 13 tagnya kosong, 157 berisi teks placeholder. Seluruh **601** berkas punya
tag kedua terisi, dan **tidak ada satu pun** rule yang kondisinya benar-benar tidak terbaca.

Klaim lama "lima rule `When` kondisinya kosong" berasal dari membaca satu tag saja, dan **keliru**.

Predikat individual **bukan** seam. Satu pintu masuk, diuji table-driven.

**Blocked by:** None (can start immediately)

**Status:** ready-for-human — diport 01-10-2026, menunggu tinjauan work owner

- [ ] Satu pintu masuk evaluasi bernama, melayani seluruh registry — bukan satu fungsi publik per predikat
- [ ] Kondisi dibaca dari **kedua** tag; rule yang kondisinya hanya ada di tag kedua tetap terbaca benar
- [ ] Teks placeholder pada tag pertama **tidak** diperlakukan sebagai kondisi
- [ ] **Nama predikat yang tidak dikenal → gagal keras**, bukan mengembalikan salah — salah ketik tidak boleh berubah menjadi gerbang yang selalu tertutup
- [ ] Setiap predikat menyebut rule Pega asalnya dalam komentar
- [ ] Uji table-driven mencakup contoh dari kedua kelompok: yang kondisinya terbaca di tag pertama, dan yang tersembunyi di tag kedua


## Comments

### 2026-10-01 — implementasi (agent, `/implement`)

Kode: `APP_RNM/modul/nbfacin/backend/services/rules/` — `rules.go` (seam `Eval`), `registry_gen.go`
(bangkitan, 209 predikat), `bangkit/` (pembangkit), `logika/` (pemecah token bersama), tiga berkas uji.
Belum di-commit. Keputusan: `../KEPUTUSAN-30-09-2026.md` butir 19–23.

| Kriteria | Keadaan |
| --- | --- |
| Satu pintu masuk bernama untuk seluruh registry | ✅ `rules.Eval(nama, kasus) (bool, error)`; nama tidak peka huruf seperti Pega |
| Kondisi dari kedua tag; yang hanya di tag kedua tetap terbaca | ✅ registry dibangun dari ekspresi tersimpan (`pyLogic` + `pyConditionValue1`) untuk SEMUA predikat; `IsBonding` (tampilan placeholder) dan `IsAsuransiKredit` (tampilan terbaca) teruji |
| Placeholder tag pertama bukan kondisi | ✅ pembangkit tidak membaca `pyConditionString` sama sekali |
| Nama tak dikenal → gagal keras | ✅ panic |
| Setiap predikat menyebut rule asal | ✅ komentar + medan `asal` per entri, diperiksa `TestRegistryUtuh` |
| Uji table-driven kedua kelompok | ✅ `TestEvalPredikatNyata`, ditambah uji seluruh registry |

⚠️ **Angka tiket dibetulkan:** "601" adalah jumlah berkas tiga folder; registry NB memuat **209**
predikat dari 210 berkas. **13 predikat panic** — IsPKSASM (sikap spec), 11 berbentuk kondisi lain
(platform/CRM, `IsDeducType`, `IsTypeDeductType`, `IsShowInput`, `IsPASSG`, `StepStatusFail`), 1 merujuk
`pyIsIPad` yang tidak ada di korpus (dihitung dari `registry_gen.go`). Predikat lain yang merujuk
salah satunya ikut panic saat dievaluasi.

Uji: pembangkit deterministik (dua kali, hash sama); registry dicocokkan dengan pembacaan independen
korpus (4 ukuran sama); 11 mutasi atas `rules.go` semuanya tertangkap.

Batas yang disadari: belum ada pemanggil (`acceptance`, alur, layar); relatif `.X` diselesaikan oleh
`Kasus` pemanggil; properti tidak ada diperlakukan kosong `[dugaan]`; semua kondisi dievaluasi tanpa
hubung-singkat.

### 2026-10-01 — tindak lanjut `/code-review` (agent)

⚠️ **Entri di atas basi pada jumlah panic dan operand.** Yang berlaku (`../KEPUTUSAN-30-09-2026.md`
butir 24–25):

- **Operand kanan tanpa kutip → panic** (temuan review spec): 43 baris di 6 predikat, termasuk `IsUW`.
  Sebelumnya dibaca sebagai properti, sehingga `IsUW` terbuka bagi operator tanpa workbasket.
  Pertanyaan ke work owner: `../PERTANYAAN-NB08.md`. **Predikat panic kini 19.**
- Semua baris yang dirujuk dinilai **urut label** (A, B, …, AA) — hasil dan galat deterministik.

Perbaikan lain dari review: label `[dugaan kuat]` → `[dugaan]`; perintah audit untuk angka 209;
komentar kembar ISFLAGOLDDATA dibetulkan (pembangkit memakai sidik isi; hash ternormalisasi diukur
terpisah); pengurai dan pemecah token kini satu paket (`logika`), dipakai pembangkit dan penilai;
logika diurai sekali saat paket dimuat; pembangkit menyimpan struktur kondisi, bukan teks yang diurai
ulang; uji unit pembangkit + penjaga `TestRegistrySamaDenganKeluaranBangkit` (butuh env
`KORPUS_NB_WHEN`; dilewati tanpanya); komentar logika boolean; loop uji digabung.

Uji: 14 mutasi, seluruhnya tertangkap (kode keluar proses dihitung); dua tes diperketat agar memeriksa
alasan panic. Pembacaan silang korpus: 209/209 nama, logika, baris, 19/19 panic.

Batas yang tersisa: penjaga registry tidak jalan di CI (korpus tidak di repo); kunci registry hanya
nama — aman untuk NB, perlu kelas bila dipakai lintas folder; `IsFacout` (komentar "Banding") dan
`IsSpreadingDepan` (folder lain) menunggu tiket 09 dan butuh perluasan pembangkit.

### 2026-10-01 — operand tanpa kutip dibaca sebagai teks (keputusan work owner, butir 28)

Butir 24 diganti butir 28. Pembangkit kini membaca satu kata tanpa kutip (`ReasFacInDirector`, `A1`,
`Offer`) sebagai `teksLiteral`; `IsUW`, `IsLimitSBondKBG`, `IsLimitCreditCL`, `IsCedingConfirmOffer`,
`IsCedingConfirmPolicy` tidak lagi panic. Operator tanpa workbasket **tidak** membuka `IsUW`
(`TestEvalKataTelanjangSebagaiTeks`). Dua pengecualian agent, menunggu konfirmasi:

- **A13** — `IsVisible` (`= True`) tetap ditolak: teks `"True"` ≠ `"true"` milik `booleanLiteral`.
- **A14** — baris pembanding login operator (`.pyUserIdentifier`, `.pxCreateOperator`,
  `.pxUpdateOperator`) ditolak dan **nilainya tidak ditulis** (CLAUDE.md §4 butir 10). Saat regenerasi
  ditemukan lima nilai login sudah tertulis di `registry_gen.go` sejak pengerjaan awal tiket ini; kini nol
  (`TestLoginTidakDisalin`; cek dua cara di `../KEPUTUSAN-30-09-2026.md` bawah butir 35).

Predikat panic statis: 19 → 17 (dua cara + rekonsiliasi delta, `../KEPUTUSAN-30-09-2026.md`).

### 2026-10-01 — A13 dan A14 dikonfirmasi (butir 37)

`IsVisible` (`= True`) tetap ditolak; pembanding login operator tetap ditolak tanpa menulis nilainya.

### 2026-10-01 malam — literal identitas tidak disalin (agent)

Generator `bangkit`: `medanIdentitas` (`OldPolicyNo`, `PolicyNo`, `MarketingName`, `pyTelephone`) — pembanding dengan
literal tidak-kosong menjadi panic tanpa nilai, setara `medanLogin`. Predikat yang kini panic: `ISERRORSPREADING`,
`ISTBONDING`, `ISTREATY1`, `ISSPVTREATY1`. Generator juga menerima `-folder`/`-paket` (registry EDM, butir 56–57).
Penjaga baru: `TestRegistryTanpaLiteralIdentitas`, `TestPilihProfil`, `TestBangkitFolderLain`. Rincian: register bab
"Temuan data identitas".

