# PROMPT INDUK — **tiga modul sekaligus**: Claim Life *(main)* · PremiumList Life · Komite Claim Life — tiga sesi executor, tiga worktree, satu repo

> Perintah work owner 27 September 2026: *"KALAU BISA MENGERJAKAN 3 MODUL SEKALIGUS!"* Bisa, dengan
> **tiga sesi executor berjalan berdampingan**, masing-masing di **worktree git** dan **cabang**
> sendiri, memakai **satu** Oracle DEV yang sama. Berkas ini adalah aturan bersamanya. Tiap sesi
> menempel **brief modulnya** *(bab 3)* sesudah membaca berkas ini.

---

## 0. UNTUK WORK OWNER DULU — kenapa layar kosong, dan tiga perintah yang membereskannya

Diperiksa 27 September 2026 pukul 14.50 di mesin ini:

| Temuan | Bukti |
| --- | --- |
| Backend **sudah jalan** di `:8080`, tetapi **tanpa database** | `GET /healthz` → `{"status":"sehat","database":"tidak dikonfigurasi"}`; proses `api` dimulai 14.45 **tanpa** `muat-env.ps1` *(ORACLE_DSN kosong → "berjalan tanpa database")* |
| Dengan `.env` dimuat, Go **menjangkau Oracle** | backend uji di `:8081` sesudah `. .\muat-env.ps1` → `{"database":"terjangkau"}` *(dimatikan lagi)* |
| Tabel aplikasi **belum ada** di DEV | katalog `POOLDATA`: `T_MIGRASI`, `T_WORK_CLAIM`, `T_GENERAL_CLAIM` **tidak ada** — migrasi belum pernah dijalankan; nol tabrakan nama *(diperiksa ulang)* |

Jalankan di **satu** jendela PowerShell, berurutan *(keputusan **as** brief lanjutan 6: migrasi hanya
membuat objek baru, tidak menyentuh tabel warisan; `-migrate-down` menolak `POOLDATA`)*:

```powershell
Stop-Process -Name api -ErrorAction SilentlyContinue      # backend lama yang tanpa .env
Set-Location 'D:\XML\RNM_BRD\OUTPUT_HASIL_RNM\APP_RNM'
. .\muat-env.ps1                                           # titik, spasi, nama skrip
go run .\cmd\api -migrate                                  # sekali; membentuk tabel lalu keluar
go run .\cmd\api                                           # biarkan jendela ini hidup
```

Lalu muat ulang `http://localhost:5173`. Inbox akan **kosong tetapi hidup** *(belum ada klaim)*;
daftarkan satu klaim lewat `Register` dengan nomor premium list yang ada di `M_LIFE_PREMIUM_DETAIL`.

---

## 0.1 KEPUTUSAN **ax** `[DIPUTUSKAN work owner — 27 September 2026]`: `T_EFEK_KELUAR` → **`T_LOG_SERVICE_RNM`**

Perintah work owner: *"JIKA ITU UNTUK LOG SERVICE MAKA NAMANYA YANG JELAS, T_LOG_SERVICE_RNM."*

| Butir | Ketentuan |
| --- | --- |
| Nama | tabel **`T_LOG_SERVICE_RNM`**; index `IX_LOG_SERVICE_RNM_JADWAL` *(STATUS, JADWAL_BERIKUT)* dan `IX_LOG_SERVICE_RNM_RUJUKAN`; sequence `SEQ_LOG_SERVICE_RNM`. Kolom dan maknanya **tetap** *(antrean + catatan hasil tiap panggilan layanan luar: STATUS, PERCOBAAN, JADWAL_BERIKUT, GALAT_TERAKHIR)* — yang berubah hanya namanya |
| Cara | migrasi `015` **belum pernah dijalankan di Oracle mana pun** *(`T_MIGRASI` tidak ada di DEV)* → berkasnya **disunting di tempat** dan diganti nama `015_t_log_service_rnm.sql` + `_down.sql`; bukan migrasi baru |
| Kode | `repository/efekkeluar.go`, `services/antrean.go`, `repository/migrasi_test.go` *(nama tabel di uji bentuk)*; nama tipe Go `EfekKeluar` boleh tetap — yang diminta work owner adalah nama **tabel** |
| Dokumen | tiket 12 blok `### Ralat menurut keputusan work owner — 27 September 2026 (butir ax)`; `STRUKTUR-TABEL-CLAIM-LIFE.md`; PARITAS bila menyebutnya |
| Siapa, kapan | **sesi Claim Life, di `main`, sebagai commit pertama giliran** `claim-life: ax — T_EFEK_KELUAR menjadi T_LOG_SERVICE_RNM` — **sebelum** `git worktree add` §1, supaya kedua cabang lahir sudah membawa nama baru |
| Padanan warisan | `pooldata.monitoring_klaim_log` *(ditulis `InsertLogServiceClaim`)* tetap **tidak** dipakai: ia log tanpa antrean; tabel baru ini log **dan** antrean |

---

## 1. TATA LETAK — worktree, cabang, port

| Modul | Tempat kerja | Cabang | Backend | Vite |
| --- | --- | --- | --- | --- |
| **Claim Life** *(lanjutan)* | `D:\XML\RNM_BRD\OUTPUT_HASIL_RNM\` *(worktree utama)* | `main` | `:8080` | `5173` |
| **PremiumList Life** | `OUTPUT_HASIL_RNM\.worktrees\premiumlist-life\` | `modul/premiumlist-life` | `:8082` | `5174` |
| **Komite Claim Life** | `OUTPUT_HASIL_RNM\.worktrees\komite-claim-life\` | `modul/komite-claim-life` | `:8083` | `5175` |

`.worktrees/` sudah masuk `.gitignore`. Membuatnya *(work owner atau sesi Claim Life, sekali)*:

```powershell
Set-Location 'D:\XML\RNM_BRD\OUTPUT_HASIL_RNM'
git worktree add -b modul/premiumlist-life  .worktrees\premiumlist-life  main
git worktree add -b modul/komite-claim-life .worktrees\komite-claim-life main
Copy-Item APP_RNM\.env          .worktrees\premiumlist-life\APP_RNM\.env
Copy-Item APP_RNM\frontend\.env .worktrees\premiumlist-life\APP_RNM\frontend\.env
Copy-Item APP_RNM\.env          .worktrees\komite-claim-life\APP_RNM\.env
Copy-Item APP_RNM\frontend\.env .worktrees\komite-claim-life\APP_RNM\frontend\.env
```

Lalu di tiap worktree sunting **hanya** `.env`-nya *(tidak masuk git)*: `HTTP_ADDR=:8082` /
`:8083`; `frontend\.env` → `DEV_PROXY_TARGET=http://localhost:8082` / `8083`; Vite dijalankan
`npm run dev -- --port 5174` / `5175`. Tiap worktree: `npm install` sekali di `frontend\`.
Korpus `D:\XML\RNM_BRD\` tetap READ-ONLY dari semua worktree.

---

## 2. PEMBAGIAN — supaya tiga cabang bertemu tanpa saling menimpa

| Hal | Aturan |
| --- | --- |
| **Nomor migrasi** | Claim Life `017`–`029` · Komite `030`–`049` · PremiumList Life `050`–`079`. Nama berkas `NNN_<isi>.sql` + `_down.sql`; `T_MIGRASI` mencatat per nama, jadi tiga cabang tidak saling menunggu |
| **Objek Oracle** | himpunan tabel/sequence per modul **terpisah** *(Claim Life `T_WORK_CLAIM`, `T_GENERAL_CLAIM`, `T_CLAIMLF_*`, `T_LOG_SERVICE_RNM`; Komite `T_GENERAL_KOMITE`, `T_KOMITE_KOMITELIST` — sudah dibuat migrasi `013` Claim Life; PremiumList `T_WORK_POLIS`, `T_PREMIUM_LIST*`, `T_VIEW_SUGGEST`)*. Sebelum `-migrate` pertama di cabangnya, executor memeriksa **tabrakan nama** di katalog *(pola keputusan as)* dan mencatat hasilnya |
| **Berkas backend** | paket tetap `internal/{handlers,services,repository,models}`; berkas milik modul **berawalan**: `polis_*.go` *(PremiumList)*, `komite_*.go` *(Komite)*; rute modul di `handlers/rute_<modul>.go` berisi satu fungsi `daftarkanRute<Modul>(mux, svc, stub)` yang dipanggil **satu baris** di `Router` — satu-satunya sentuhan ke `handlers.go` |
| **Berfile milik modul lain** | **tidak disunting**. Perubahan pada kode bersama *(`services.Service`, `models` bersama, `repository/db`)* hanya **aditif** *(fungsi/berkas baru)*, dicatat di laporan akhir dengan sebabnya |
| **Frontend** | `src/pages/<modul>/`, `src/assets/labels.<modul>.ts` *(pola bukti baris XML)*, satu `KelompokMenu` per modul di `Shell.tsx` yang isinya dari `butirMenu<Modul>` di berkas modul; `PARITAS-LAYAR-DAN-AKSI.md` dan `LAPORAN-GILIRAN.md` di `.scratch/<modul>/` |
| **Menu** | hanya yang berbukti XML modul itu *(lanjutan 7 §1)*: PremiumList Life punya harness portal `PremiumLife_harness` *(root sheet struktur)*; Komite **tidak** punya harness → satu butir inbox worklist |
| **Oracle bersama** | ketiga cabang boleh `-migrate` langkahnya sendiri di DEV *(objek terpisah, ledger per nama)*; **tidak pernah** `-migrate-down` di `POOLDATA`; skema uji kosong dari DBA belum ada → `db` test SKIP di semua cabang |
| **Penomoran & prosedur** | keputusan **o** berlaku di ketiga modul: prosedur Oracle **tidak dipanggil**, logikanya ditiru di Go *(`PenomorCounter` o1–o3, `SUMBER-PENOMORAN-DBA.md`)*; ADR-U-0029 nol `COMMIT` di SQL |
| **Kontrak antar modul** | Claim Life → Komite: muatan `serahkanKomite` + `T_LOG_SERVICE_RNM` *(A2)*; Komite → Claim Life: `STS_REJECT` tingkat akhir *(tiket Claim Life 11)*; PremiumList → Claim Life: `T_PREMIUM_LIST` dibaca Claim Life untuk `PolicyDataLife` **sesudah merge** *(keputusan av)* — PremiumList menyediakan pembaca `repository.PolisRingkas(plNumber)` yang didokumentasikan di tiket 08-nya |

---

## 3. BRIEF PER SESI

| Sesi | Tempel utuh | Titik mulai |
| --- | --- | --- |
| Claim Life | `PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE-LANJUTAN-10.md` | `main` @ `cfc4824` |
| PremiumList Life | `PROMPT-IMPLEMENTASI-MODUL-PREMIUMLIST-LIFE.md` | `modul/premiumlist-life` dari `main` |
| Komite Claim Life | `PROMPT-IMPLEMENTASI-MODUL-KOMITE-CLAIM-LIFE.md` | `modul/komite-claim-life` dari `main` |

Aturan yang sama di ketiganya: XML menang dan tiket diralat dengan bukti; penyarangan dibaca sebagai
**pohon**; mekanisme giliran lanjutan 8 §1 *(nol pesan di antara paket, catatan ke berkas laporan
giliran)*; larangan keamanan brief modul Claim Life *(nama orang, nomor polis, kredensial, URL,
username, salt token; fixture `UJI-*`; approval manusia untuk sistem luar)*; bab TELEMETRI EKSEKUSI.

---

## 4. PENYATUAN — urutan merge dan siapa

1. Sesi modul selesai → laporan akhir → work owner *(atau asisten, bila diminta)* di worktree utama:
   `git merge --no-ff modul/<nama>` → uji penuh di `main` *(Go, vitest, build)* → `go run .\cmd\api
   -migrate` bila ada langkah baru → commit merge.
2. Urutan yang disarankan: **Komite** dulu *(ia hanya membaca kontrak Claim Life yang sudah ada)*,
   lalu **PremiumList Life**, lalu Claim Life menyambungkan `PolicyDataLife` ke `T_PREMIUM_LIST`
   *(paket kecil di gilirannya berikutnya)*.
3. Konflik yang diperkirakan dan boleh diselesaikan tanpa bertanya: `handlers.go` *(baris pendaftaran
   rute)*, `Shell.tsx` *(kelompok menu)*, `labels.ts` *(impor)*, `.gitignore`. Konflik pada berkas
   milik modul lain **tidak boleh** — itu tanda aturan bab 2 dilanggar.

---

*Disusun 27 September 2026 sesudah pemeriksaan backend (`:8080` tanpa database, `:8081` dengan
`.env` terjangkau), katalog DEV (`T_MIGRASI` belum ada), `git worktree` tersedia (git 2.55), struktur
`Router`/`Shell`, dan kesiapan dokumen kedua modul lain (spec + tiket `ready-for-agent`: PremiumList
Life 10 tiket, Komite Claim Life 11 tiket).*
