# PROMPT — lanjutan 6: **KERANGKA TAMPILAN DULU** (homepage + komponen dasar dari `REFERENSI_UI`), lalu A3 layar demi layar di dalamnya

> **Skill:** ketik `/mattpocock-skills:implement` sebagai manusia, lalu tempel berkas ini **utuh**.
> Brief modul, lanjutan 1–3, **lanjutan 4 bab 0–8**, dan **lanjutan 5** berlaku seluruhnya. Berkas
> ini **mengubah urutan** atas instruksi work owner 27 September 2026: *"Buatkan dahulu struktur
> awalnya secara lengkap (homepage) dan komponen awal di dalamnya; referensi UI dari
> `D:\XML\RNM_BRD\REFERENSI_UI\frontend`. Menyelesaikan satu modul tanpa tampilan terlalu lama."*
>
> **GILIRAN INI:** Langkah 0 → **F0 kerangka tampilan** *(§2)* → **A3 per layar di dalam kerangka**
> *(§3)* → A4 → Bagian B. Rantai tanpa mengakhiri giliran; XML menang; pohon XML untuk penyarangan;
> satu pesan ke manusia di akhir.

---

## 0. KEADAAN AWAL — 27 September 2026, sesudah `e4bb014`

| | Keadaan |
| --- | --- |
| `HEAD` | `e4bb014` *A3 — sensus paritas*; 11 commit kode sesudah `f8bbf7a` *(A2 tutup: roster+kasus Komite, resolver, token, kategori ar1, outbox aq; perapian 405; audit jalur kegagalan; awalan nomor di-lookup; nomor di-commit sebelum pendaftaran; 13 centang AC bergeser)*. Working tree bersih |
| Uji | vet · vet db · gofmt nol · build · **261 PASS · 0 FAIL · 34 SKIP** · 16 JS · 88 modul; 49 berkas +3.029/−159 sejak `f8bbf7a` |
| Verifikasi independen | seluruh angka tereproduksi; `CheckTotalAdjustmentClaim` memang **tidak ada** di korpus *(dirujuk 10 section, berkasnya tidak diekspor — catat sebagai `[tidak ada di korpus]`, bukan ditebak)*; `PARITAS-LAYAR-DAN-AKSI.md` 164 baris: **24 ada · 10 belum · 6 tidak ditiru** |
| Permukaan | rute HTTP 11 · halaman React **2** · `api.ts` 17 fungsi · 1 berkas test JS · nol CSS · `App.tsx` = dua halaman ditumpuk |
| Sisa kecil | tujuh **tipe** stub masih terdefinisi *(`Antrean…`, `Jejak…`, `KasusKomite…`, `KategoriWajib…`, `Penomor…`, `Resolver…`, `Roster…`)* walau implementasinya sudah diganti — buang yang tak terpakai atau pindah ke `_test.go` |

---

## 1. `.env` DAN KONEKSI ORACLE — kenapa terasa "belum ada", dan apa yang benar

**Fakta yang dicek asisten 27 September 2026:** `APP_RNM\.env` **ada dan terisi** — `HTTP_ADDR=:8080`,
`ORACLE_DSN=oracle://POOLDATA:<sandi>@<host>:1521/DEV_NUSARE2`, `ORACLE_SCHEMA=POOLDATA`,
`IS_PEGA_PROD=false`; `frontend\.env` terisi *(`VITE_API_BASE_URL=` kosong, proxy dev ke `:8080`)*.
`muat-env.ps1` dan `muat-env.cmd` ada.

| Yang terjadi | Sebab | Yang benar |
| --- | --- | --- |
| `go run`/`go test` berkata `ORACLE_DSN belum dikonfigurasi` | **Go tidak membaca `.env`** *(disengaja — `PANDUAN-MENJALANKAN.txt` bab 3)*; hanya Vite yang membaca `frontend\.env` sendiri | di PowerShell yang sama: `Set-Location APP_RNM` → `. .\muat-env.ps1` *(titik, spasi)* → `go run .\cmd\api` — variabel hidup di jendela itu |
| 34 test **SKIP** walau `.env` terisi | test `db` menuntut `ORACLE_SKEMA_UJI=true` **dan skema bukan `POOLDATA`** *(`config.PagarSkemaUji`)*. Skema uji **membuat tiruan bernama sama dengan tabel warisan** *(`M_LIFE_PREMIUM_DETAIL`, `RETROCESSIONLIFE`, `EMAILKOMITE`, `GCP_IMAGE`, …)* lalu **membuangnya** — dijalankan di `POOLDATA` ia menghapus tabel warisan sungguhan | user Oracle **kosong** dari DBA. `[data DBA]` akun `POOLDATA` hanya punya `CONNECT`, `RESOURCE`, `CREATE TABLE/SEQUENCE/VIEW`, `UNLIMITED TABLESPACE` — **tidak** punya `CREATE USER`, dan tidak ada user uji lain di instance. Sampai DBA membuatnya, 34 SKIP adalah keadaan yang **benar**; executor **tidak** mengakalinya |
| Aplikasi jalan tetapi layar klaim 500 | tabel baru **belum pernah dibuat** di Oracle mana pun | keputusan **as** di bawah |

**Keputusan `as` — work owner, `[USULAN — rekomendasi setujui]`: jalankan `-migrate` ke DEV `POOLDATA`
supaya aplikasi berjalan dengan data warisan sungguhan dan layar berisi.** Bukti keamanannya
`[data DBA — dibaca 27 September 2026]`: seluruh **33 objek** yang dibuat 15 langkah migrasi *(11
tabel `T_*`, 9 sequence `SEQ_*`, 13 index `IX_*`/`UX_*`)*, **21 constraint** bernama, dan `T_MIGRASI`
**tidak satu pun ada** di `POOLDATA`; migrasi hanya `CREATE` objek baru dan `ALTER` tabel **miliknya
sendiri** *(`010`, `011`, `014`)*; pra-terbang bentuk berjalan; `IS_PEGA_PROD=false`. Yang perlu
diketahui sebelum menyetujui: `-migrate-down` **menolak** `POOLDATA` *(pagar yang disengaja)*, jadi
membatalkannya kelak = DBA membuang 33 objek itu manual; instance DEV dipakai bersama, objek baru
terlihat pengembang lain. Bila disetujui, **manusia** yang menjalankannya *(bukan executor — modul
§9)*:

```powershell
Set-Location D:\XML\RNM_BRD\OUTPUT_HASIL_RNM\APP_RNM
. .\muat-env.ps1
go run .\cmd\api -migrate      # cetak ringkasan objek yang dibuat, lalu berhenti
go run .\cmd\api               # backend :8080
# jendela lain:
Set-Location D:\XML\RNM_BRD\OUTPUT_HASIL_RNM\APP_RNM\frontend ; npm run dev
```

Sandi yang tertulis di chat berkali-kali tetap sebaiknya **diganti DBA** sesudah pekerjaan ini;
`.env` tidak pernah masuk git *(dicek: `.gitignore` pola `*.env`)*.

---

## 2. F0 — KERANGKA TAMPILAN DARI `REFERENSI_UI\frontend`

**Sumber `[dibaca asisten]`:** `D:\XML\RNM_BRD\REFERENSI_UI\frontend` — **READ-ONLY** *(disalin ke
pohon kita, tidak pernah disunting di tempatnya)*; aplikasi `treaty-v2-frontend`: Vite 6 + React 18 +
TypeScript + `@tanstack/react-query`, uji `node --test`. Strukturnya *(README-nya, ADR-011 butir 5)*:
`App.tsx` *(kerangka: Login → Shell dengan sidebar terlipat/laci + topbar + palet menu Ctrl+K + menu
profil; 2.376 baris)* · `assets/styles.css` *(3.956 baris; token `--navy #03045e`, `--blue #0077b6`,
`--sidebar-w 272px`, `--sidebar-w-terlipat 72px`, `--topbar-h 64px`, radius, ring, transisi)* ·
`assets/labels.ts` *(pusat label: `MENU`, `labelKolom`, `labelForm`, `LABEL_TAHAP`, `LABEL_KEPUTUSAN`
— tiap teks membawa komentar bukti sumbernya)* · `components/ui/dasar.tsx` *(`Field`, `FieldSandi`,
`FieldTanggal`, `Area`, `Pilih`, `Gagal`, `BelumTersedia`, `Modal`, `ModalBaris`, `StripTab`, `Panel`,
ikon, `Memuat`, `Kosong`, `BarisSkeleton`, `Halaman`, `UKURAN_HALAMAN`)* · `components/`
*(`KelompokMenu`, `PaletMenu`, `CariSebaris`, `PilihCari`, `RecordForm`, `BilahSaringRegistry`)* ·
`hooks/useHalaman` · `lib/` *(`format`, `desimal`, `tanggalInput`, `singkatan`, `lipatMenu`,
`pohonLipat`, `keadaanGalat`, `exportXlsx`, …)* · `store/sesi.ts` *(token di `sessionStorage`)* ·
`services/api.ts` *(`ApiFailure`, `auth`, `token`, `bolehUlang`)* · `PagarGalat.tsx` · `pages/` per
menu. Kita: Vite 5 + React 18 + TS + axios + zustand + vitest, **nol** CSS, dua halaman ditumpuk.

**Aturan adopsi `[DIPUTUSKAN — work owner]`:**

| # | Aturan |
| ---: | --- |
| 1 | **Struktur, CSS, komponen dasar, dan pola shell diambil dari referensi** — disalin ke `frontend/src` dengan nama berkas yang sama, lalu diadaptasi. Yang **tidak** diambil: halaman treaty *(`pages/offer`, `realization`, `citrix`, `borderaux`, `master/*`)*, berkas `.test.ts` referensi *(uji kita vitest)*, logo dan judul **e-Treaty** *(aset produk lain — pakai teks "Nusantara Re" dan favicon netral sampai aset resmi diberikan)* |
| 2 | Dependensi: **boleh** menambah `@tanstack/react-query` *(sama dengan referensi)*; `zustand` dipakai untuk sesi pelaku; `axios` boleh diganti klien `fetch` referensi *(`ApiFailure`, `bolehUlang`)* bila `api.ts` kita dimigrasi utuh — satu pilihan, dinyatakan di tiket |
| 3 | `lib/desimal.ts` referensi **diperiksa** sebelum dipakai: uang di proyek ini **teks desimal eksak**; fungsi yang lewat `Number`/float untuk uang **tidak** dipakai *(ADR-U-0003)*; format tampilan boleh |
| 4 | `assets/labels.ts` kita mengikuti pola referensi: **setiap teks menu/layar membawa bukti XML** — judul tahap VERBATIM `Register_Flow` *(`Input Register`, `Outstanding Claim`, `Medical Check`, `Claim Analis`)*, judul flow action VERBATIM *(`Attach Document Life`, `Claim Life - Upload CSV`, `Adjustment_Detail`, …)*, label tombol VERBATIM section *(`Save Adjustment`, `Save to Outstanding`, `Reject Outstanding`, `Send Back to Medical`, `Send Back to Admin`, `Close Claim`, …)*; terjemahan Indonesia boleh **di samping**, bukan menggantikan |
| 5 | Login = **stub** selama `AUTH_STUB=true`: layar masuk memilih akun dan peran *(`ReasLifeAdmin`, `ReasLifeMedicalAdvisor`, `ReasLifeSPV`)*, tersimpan di `store/sesi` *(pola referensi, `sessionStorage`)*, dikirim sebagai `X-Pelaku`/`X-Peran` *(ab)*; layar memasang pita peringatan "mode stub"; nol sandi disimpan; IAM sungguhan = tiket 07/ADR-U-0030 |
| 6 | **Beranda** = **antrian kerja pemegang peran**, meniru worklist/workbasket Pega dari `Register_Flow` *(`[terverifikasi]` lanjutan 4 bab 7)*: Admin melihat klaim miliknya di *Input Register* dan *Outstanding Claim* *(router `pxCreateOperator`)*; Medical Advisor melihat workbasket *Medical Check*; SPV melihat workbasket *Claim Analis*; ringkasan cacah per tahap; pintasan ke Register. Backend baru: `GET /api/klaim-life?posisi=&status=&halaman=&ukuran=` atas `T_WORK_CLAIM` ⋈ `T_GENERAL_CLAIM` *(ber-index, berbatas, gerbang peran lewat `pelakuDari`)* |
| 7 | **Menu lengkap sejak hari pertama**: kelompok **Claim Life** *(Beranda, Register, Outstanding, Medical Check, Claim Analis, Dokumen, Komite, Detail & Tutup, Cari Polis)*, kelompok **Komite Claim Life** *(Inbox per posisi, Riwayat tangga)*, dan **11 modul lain** dari `.scratch/` *(Claim FacIn, Claim Prop, Komite FacIn, Komite Prop, EDM Treaty In, Endorsement Life, Master Contract Retro Life, Master Product Name Life, NB Treaty In, PremiumList Life, Treaty Contract Out)* sebagai menu ber-`BelumTersedia` — supaya kerangka rumahnya utuh dan tiap modul berikutnya tinggal mengisi |
| 8 | Sidebar terlipat/laci, topbar, palet Ctrl+K, menu profil, `PagarGalat`, `Memuat`/`Kosong`/`Gagal` — **perilaku persis referensi** *(termasuk penutupan menu profil lewat klik luar dan Esc, lebar tablet melipat panel)*; uji vitest untuk: menu per peran, palet, login stub, klien API gagal, halaman ada di dalam shell |

**Paket kerja F0 — satu commit tiap paket, `frontend: kerangka — <isi>`:**

| Paket | Isi | Bukti selesai |
| --- | --- | --- |
| F0.1 | `styles.css`, `components/ui/dasar.tsx`, `lib/{format,desimal,tanggalInput,singkatan,keadaanGalat}`, `PagarGalat`, `assets/labels.ts` *(pola)* — disalin, diadaptasi, `tsc` bersih | `npm run build` hijau; `tsc` bersih |
| F0.2 | `store/sesi` *(pelaku + peran)*, `services/api.ts` *(klien + `ApiFailure` + header stub)*, **Login stub** | test: masuk → header terkirim; keluar → hilang |
| F0.3 | **Shell**: sidebar + `KelompokMenu` + `PaletMenu` + topbar + profil + laci ponsel; **menu lengkap** butir 7 dengan `BelumTersedia` | test: menu per peran; Ctrl+K; lipat/laci |
| F0.4 | **Beranda antrian kerja** + rute list backend butir 6 + `Halaman`/`useHalaman` | test murni service list; test JS beranda per peran |
| F0.5 | Dua halaman yang ada *(`RegisterKlaim`, `KlaimLife`)* **dipindah ke dalam shell** sebagai menu Register dan Detail; `App.tsx` lama dibuang | `npm run dev`: masuk → shell → beranda → Register → Detail, nol galat konsol |

`PANDUAN-MENJALANKAN.txt` bab frontend diperbarui *(login stub, peran, cara membuka beranda)*.

---

## 3. A3 — LAYAR DEMI LAYAR DI DALAM KERANGKA

Sensus `PARITAS-LAYAR-DAN-AKSI.md` *(24 ada · 10 belum · 6 tidak ditiru)* dan aturan paritas lanjutan
4 §3 tetap; **urutan kelompok diubah supaya tampilan cepat terlihat**, tiap kelompok = rute backend +
halaman/komponen di dalam shell + test, satu commit:

| Urut | Kelompok | Layar yang lahir |
| ---: | --- | --- |
| 1 | Register | `InputRegisterClaimLife` + Cari Polis *(`SearchPolicy_Harness/Section`)*, pilih peserta, tanggal kejadian, validasi DOL/STNC, upload CSV |
| 2 | Outstanding | `OSClaimLife` + `InputOSClaimLife`, `Adjustment_Detail` *(enam field bank)*, `RetroClaimLife`, Save to Outstanding, Reject Outstanding |
| 3 | Detail & Tutup | `ViewClaimDetailLifeGCNM` *(panel detail yang sama dipakai tiga tahap, tombol Save Adjustment bergerbang activity — bab 7)*, `ShowEditClaimLife`, `CloseClaim` |
| 4 | Dokumen | `AttachDocumentLife`, `DocumentLife`, unduh, hapus lampiran; kategori wajib **ar1** |
| 5 | Medis | `MedicalCheck`, `Diagnose_Harness/Section`, cari diagnosa, `SetDisease`, kirim balik; **al** diputuskan dari bukti di sini |
| 6 | Akseptasi | `AkseptasiClaimLife` + `InputAkseptasiClaimLife` *(Save, Send Back to Medical/Admin, Close Claim; gerbang `Type` TP/TR)* |
| 7 | Komite | `ClaimComite`, `Committe_Life`, kirim ke Komite, daftar roster |

Setiap layar: himpunan field **persis** section *(dibaca sebagai pohon)*, `<pyCondition>` → gerbang +
kontrol tersembunyi, tombol → rute; `PARITAS-…md` diperbarui per commit. Yang tidak ditiru dinyatakan
dengan bukti *(`NextPrev`, `setVisibility_Act`, …)*.

---

## 4. A4 DAN BAGIAN B

Persis lanjutan 5 §4–§5. Halaman Komite *(inbox per posisi, riwayat tangga)* lahir **di dalam shell
yang sama**, kelompok menu Komite Claim Life.

---

## 5. LANGKAH 0

`git add PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE-LANJUTAN-6.md` → commit `docs: brief lanjutan 6 —
kerangka tampilan dulu dari REFERENSI_UI, .env dan keputusan as` → `git status --porcelain` kosong →
uji hijau *(261 · 34 SKIP · 16 JS · 88 modul)* → SHA = titik tetap **F0.1**.

---

## 6. LAPORAN AKHIR GILIRAN

Persis lanjutan 5 §7, ditambah **bab tampilan**: daftar halaman yang dapat dibuka *(menu → halaman →
rute yang dipanggil)*, cara membukanya *(`npm run dev`, akun stub per peran)*, cacah komponen dari
referensi yang dipakai / diadaptasi / tidak dipakai beserta sebabnya, dan cacah aksi layar dari 31 yang
kini punya rute **dan** kontrol.

---

*Disusun 27 September 2026 sesudah verifikasi `e4bb014` (uji dijalankan ulang, 11 commit dibaca
statistiknya, sensus paritas dibaca), pembacaan `REFERENSI_UI\frontend` (README, `App.tsx` kerangka
Shell/Login, `labels.ts`, `dasar.tsx`, token `styles.css`, `sesi.ts`, `package.json`), pemeriksaan
`.env` (bentuk, tanpa nilai), hak akun `POOLDATA` di katalog, dan pemeriksaan tabrakan 33 objek + 21
constraint migrasi terhadap `POOLDATA` (nol tabrakan).*
