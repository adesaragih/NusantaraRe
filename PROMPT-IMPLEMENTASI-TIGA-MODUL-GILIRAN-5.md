# PROMPT — GILIRAN 5 *(sesi tunggal di `OUTPUT_HASIL_RNM`, cabang `main` @ `2da4cc1`)*: **paket 0 BERANDA + STRUKTUR MENU LENGKAP dari REFERENSI_UI (bg) → PremiumList 01 bagian 2 → 09 (+ av) → Komite 01 → 09 → Claim Life §3.1**

> Brief GILIRAN-3 *(A/B/C)* dan GILIRAN-4 **tetap berlaku**; brief ini menambahkan **paket 0** yang diminta work owner
> 28-09-2026: *"buatkan dahulu struktur awalnya secara lengkap (homepage) dan komponen awal di dalamnya, referensi UI dari
> `D:\XML\RNM_BRD\REFERENSI_UI\frontend`"*. Mekanisme giliran = lanjutan 8 §1; berhenti sah hanya pada batas tiket.

---

## 0. KEADAAN AWAL — DIVERIFIKASI ASISTEN 28-09-2026

| Hal | Keadaan |
| --- | --- |
| `main` | `2da4cc1` *(ralat tiket 01 PremiumList: `POSITION` = `Offer`/`Premium` posisi layar — `ProtectAccept.xml` b1207/b2288, `InputPolicyHolder.xml` b712, `InputOfferLife_preAct` b271, `countCategoryAttachment_act` b271; `STATUS` = `pyWorkStatus`)*; `T_WORK_POLIS` sudah punya kedua kolom *(050)* → nol migrasi baru ✅ |
| Uji | **Go 407 PASS · 0 FAIL · 38 SKIP** · vitest **254** · vet, gofmt, tsc bersih *(dijalankan ulang asisten)* |
| DEV | 28 langkah migrasi lengkap; `T_VIEW_SUGGEST` berkolom `INITIAL_SUGGEST` |
| `.env` work owner | `ORACLE_DSN` ke `192.168.122.100:1521/DEV_NUSARE2` akun `POOLDATA` **sudah ada sejak awal**; `UNGGAHAN_DIR` **sudah ditambahkan asisten** *(folder `APP_RNM\unggahan`, di-gitignore)* di `.env` main dan kedua worktree |
| Backend | berjalan di `:8080`, `healthz` = `{"status":"sehat","database":"terjangkau"}` *(dihidupkan asisten, log `.scratch/backend-8080.log`)*; Vite `:5173` proxy → `:8080` |
| REFERENSI_UI yang **sudah** ada di frontend kita | `components/ui/dasar.tsx` *(25 ekspor sama persis)*, `KelompokMenu`, `Shell` *(topbar, lipat panel, profil)*, `hooks/useHalaman`, `lib/{desimal,format,keadaanGalat,lipatMenu,singkatan,tanggalInput}`, `assets/styles.css` |
| REFERENSI_UI yang **belum** ada | **Beranda** *(layar awal sesudah identitas)*, struktur menu **seluruh modul**, `PaletMenu` *(Ctrl+K, `lib/daftarMenu.ts`, `lib/entitasNavigasi.ts`)*, logo topbar *(`assets/logo-topbar.png`)*, pola grid + saringan *(`RecordForm`, `BilahSaringRegistry`, `CariSebaris`, `PilihCari`, `lib/exportXlsx.ts`, `lib/saringRegistry.ts`)* |

## 1. KEPUTUSAN **bg** — BERANDA DAN STRUKTUR MENU LENGKAP `[DIPUTUSKAN — permintaan work owner 28-09-2026; veto work owner]`

| Unsur | Isi | Bukti / batas |
| --- | --- | --- |
| **Beranda** | layar awal sesudah identitas *(stub)*: topbar berlogo, salam identitas + peran, **kartu per modul** *(nama, keadaan `aktif` / `belum dimigrasi`, tautan ke butir pertamanya)*; untuk tiga modul aktif kartu memuat **cacah antrian** per tahap dari endpoint daftar yang **sudah ada** *(Claim Life empat tahap; PremiumList; Komite)* — bila endpoint tidak mengembalikan total, tambahkan **satu** rute baca `GET /api/beranda` yang menghitung dari repository yang ada, tanpa tabel baru | `[kerangka aplikasi, bukan menu Pega]` — Beranda adalah pengganti layar awal portal, seperti `PremiumLife_harness` *(portal, kelas `Data-Portal`)*; Claim Life tidak punya harness portal terekspor *(OQ ke pemilik ekspor tetap terbuka)* |
| **Sidebar lengkap** | **17 kelompok** = folder modul korpus `D:\XML\RNM_BRD\`: Claim Fac In, Claim Life, Claim Non Prop, Claim Prop, EDM Treaty In, Endorsement Life, Endorsment Fac In, Komite Claim FacIn, Komite Claim Life, Komite Claim Non Prop, Komite Claim Prop, Master Contract Retro Life, Master Product Name Life, NB FacIn, NB Treaty In, PremiumList Life, RNW Fac In. **Butir** hanya untuk yang berbukti XML: Claim Life *(Inbox Claim Life, Register)*, PremiumList Life *(`PremiumList`)*, Komite Claim Life *(`Inbox Komite`)*. Kelompok lain **terlipat, tanpa butir**, berketerangan *"belum dimigrasi"* — **nol** butir dikarang | nama kelompok = nama folder korpus `[nama folder korpus]`, dicatat di `labels.ts`; kelompok Master/Offer/Realization/Citrix/Borderaux milik REFERENSI_UI **tidak** dibawa *(itu aplikasi Treaty, bukan korpus ini)* |
| **PaletMenu** Ctrl+K | port `components/PaletMenu.tsx` + `lib/daftarMenu.ts` + `lib/entitasNavigasi.ts` + uji `daftarMenu.sinkron.test.ts` *(dua arah sidebar ↔ palet)*; hanya butir yang ada | pola REFERENSI_UI apa adanya |
| **Logo** | salin `assets/logo-topbar.png` **hanya bila** itu logo perusahaan *(bukan tulisan "Treaty")*; bila tidak, topbar tanpa gambar | REFERENSI_UI hanya disalin, tidak disunting |
| **Pola grid + saringan** | port `RecordForm`, `BilahSaringRegistry`, `CariSebaris`, `PilihCari`, `lib/saringRegistry.ts`, `lib/exportXlsx.ts` beserta ujinya ke `components/`; **dipakai** oleh Inbox Claim Life *(empat tab)*, lalu Inbox PremiumList *(tiket 01 bagian 2)* dan Inbox Komite *(tiket 01)* — kolom tetap dari XML *(`InboxPremiumList` RD)*, ekspor xlsx hanya kolom yang tampil | komponen yang tidak dipakai layar mana pun **tidak** di-port *(kode mati)* |
| **Identitas** | tetap stub *(`VITE_STUB_PELAKU`/`VITE_STUB_PERAN`)*; **tanpa** form Login *(keputusan work owner 27-09)*; menu profil topbar menampilkan identitas + peran, tanpa `Keluar` | — |

Cara kerja: baca `REFERENSI_UI/frontend/src/App.tsx` *(b619–b750 topbar, lipat panel, profil, cari menu)*, `README.md`,
`components/PaletMenu.tsx`, `lib/daftarMenu.ts`, `lib/entitasNavigasi.ts`, `components/RecordForm.tsx`,
`BilahSaringRegistry.tsx`, `assets/styles.css` *(kelas `topbar__*`, `login`, `panel`, `halaman`)* — port **apa adanya**
*(nama Indonesia, uji ikut)*, sesuaikan hanya nama modul. Uji: cacah kelompok = 17; butir = 4; palet menemukan keempatnya;
Beranda merender kartu aktif dengan angka dari respons tiruan. Commit
`shell: Beranda + 17 kelompok modul + PaletMenu + pola grid dari REFERENSI_UI (bg)`.

> ⛔ **Catatan bertanggal 30-09-2026 — butir navigasi DICABUT** `[keputusan work owner 30-09-2026]`
> (`PROMPT-MENU-DATAR-PER-GROUPMENU.md`): *"menu jangan ada model seperti child. Buat grouping menu antar
> GROUPMENU dari tabel M_NAV_MENU. Butir inbox, register, premiumlist, komite, tco-tahun harusnya tidak
> perlu, karena 1 modul 1 menu."* Unsur **Sidebar lengkap** di atas tidak berlaku lagi untuk butirnya:
> sidebar kini kepala `GROUPMENU` + satu tombol per modul (label = nama folder korpus), klik membuka halaman
> awal modul, dan modul yang belum dimigrasi tampil sebagai tombol nonaktif, bukan kelompok terlipat.
> Migrasi `901_m_nav_menu_datar.sql` membuang kelima butir dari `M_NAV_MENU`; kode di cabang `dev`
> (`10df292`, `2a10257`, `9b9eb80`, `46d3533`). Beranda, PaletMenu, logo, pola grid, dan identitas tetap.

## 2. URUTAN GILIRAN INI — semua di `main`, nol pesan di antara paket

| # | Bagian | Rujukan |
| ---: | --- | --- |
| 0 | **bg — Beranda, sidebar 17 kelompok, PaletMenu, logo, pola grid** *(§1)* | REFERENSI_UI |
| 1 | PremiumList **tiket 01 bagian 2** — layar `Input Offer` *(`InputOfferLife` + `ConfirmSection`)*, 20 aturan `ProtectAccept` dengan pesan VERBATIM *(`Please choose no offer !`, `COB can't null`, `Sum insured number <n> is 0`, …, tiap prasyaratnya)*, kotak masuk `PremiumList` memakai pola grid §1 | brief 3-PREMIUMLIST §1–§2 |
| 2 | PremiumList **tiket 02 → 09** | brief 3-PREMIUMLIST §2 |
| 3 | Claim Life **av** sesudah tiket 04 | brief 3-CLAIM-LIFE §4 |
| 4 | Komite **tiket 01 → 09** *(Inbox Komite memakai pola grid §1)* | brief 3-KOMITE §1–§2 |
| 5 | Claim Life **§3.1** — 41 baris sensus diputuskan | brief 3-CLAIM-LIFE §3.1 |

## 3. ATURAN · LAPORAN · TELEMETRI

Sama dengan GILIRAN-4 §2–§3: XML sebagai pohon; ralat bertanggal; menu hanya yang berbukti *(kecuali kerangka §1 yang
ditandai)*; uji per paket; `-migrate` tidak dijalankan executor; endpoint luar tidak dipanggil; nol kebocoran. Laporan satu
pesan pada batas tiket: tabel **paket → tombol/kelas XML atau REFERENSI_UI → rute/komponen**, angka uji tiap commit, bab
**TELEMETRI EKSEKUSI** per paket.

---

*Disusun 28 September 2026 sesudah verifikasi `2da4cc1` (uji dijalankan ulang; XML `ProtectAccept`, `InputPolicyHolder`,
dua activity dibaca ulang), pemeriksaan `.env` dan backend (`healthz`), dan pembacaan `REFERENSI_UI/frontend` (README, App.tsx,
components, lib, assets) dibandingkan dengan `APP_RNM/frontend/src`.*
