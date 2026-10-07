# 13: Efek keluar — Kasir, Arasapas, konversi non-life, dan email komite

**Status:** sebagian 07-10-2026 — efek keluar diantre di outbox hanya bila `IS_PEGA_PROD`; panggilan nyata berhenti `…BelumDisetujui` (OQ-CP-03); email komite diantre (hanya bila `IS_PEGA_PROD`) — RALAT 07-10-2026 (semula `ready-for-agent`)
**Blocked by:** 00 (PREFACTOR) · 08 (baris adjustment)
**Menutup:** AC 89 · 90 · 91 · 92 · 93 · 94 · 99 · 100 · 101 *(9 AC)* — US 50–54

## Hasil & nilai pengguna

Data akseptasi sampai ke Kasir, klaim tersinkron ke sistem inti reinsurance non-life, dan email
pemberitahuan sampai ke anggota komite yang berwenang — **dan bila salah satunya gagal, pekerjaan itu
tidak hilang diam-diam**. Claim Admin melihat status pengiriman, dan Finance mendapat nomor kasus
Kasir untuk rekonsiliasi.

⚠️⚠️ `[terverifikasi]` Hari ini **tidak ada jaring pengaman sama sekali**: nol `Exit-Activity`, nol
`Page-Set-Messages`, nol rollback, dan **39 dari 39** penangan exception kosong — termasuk pada kedua
langkah pemanggilan REST. Tabel log yang ada adalah **jejak, bukan antrean**: nol kolom status/retry,
satu INSERT dan satu SELECT di seluruh korpus, nol UPDATE, nol job pemroses ulang.

## Area codebase

Outbox transaksional · klien Kasir · klien konversi non-life · pengirim email · resolusi endpoint.

## Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/HitServiceToKasir_Act.xml` | — | 14 langkah tanpa penanganan kegagalan; nol retry |
| `Activity/HitServiceToKasir_Act.xml` | 9.8 · 13.5 | mencatat ke tabel log **empat kolom**: muatan JSON, kunci internal Pega, nomor akseptasi, dan pesan respons — ⚠️ **nomor kasus Kasir dibuang** |
| `ConnectREST/SendAcceptationToKasir.xml` | — | memetakan seluruh JSON ke satu properti tanpa menyebut sub-properti |
| `ConnectREST/SendAcceptationToKasir.xml` | pengaturan Connect-REST | ⭐ **BARU 2026-09-19 (ronde 4)** — **satu-satunya** rule berautentikasi di 329 berkas: `pyUseAuthentication = true`, dan kredensialnya diambil dari **profil autentikasi bernama** (`pyAuthProfileSelectionType`), **bukan** header yang ditulis di badan rule. Proksi `NO_AUTH`. Catatan pengembangnya: *"tambah auth"* |
| kelima `ConnectREST/` lain di Claim Prop | pengaturan Connect-REST | ⭐ **BARU** — `GETDTLPAYMENTCLAIM` · `KONVERSIKLAIMNONLIFE` · `SERVICEGOOGLE` · `GETPAYATTACHMENT` · `GETPREMIUMPAIDONTREATYIN` semuanya `pyUseAuthentication = **false**` — ⛔ **jangan** diberi autentikasi |
| `ConnectREST/ServiceGoogle.xml` | pengaturan Connect-REST | ⭐ **BARU 2026-09-19 (ronde 6)** — **batas waktu respons `300 000` milidetik (5 menit)**, catatan pengembang *"buat jadi 300rb"*. Satu-satunya batas waktu eksplisit yang terbaca di keenam Connect-REST |
| `Activity/SetPayableTreaty_Act.xml` | 4 · 5 | ⭐ **BARU** — ⚠️ `Obj-Save` dan `Obj-Refresh-And-Lock` **ber-remark**, jadi rule ini **tidak pernah menyimpan**. Catatan pengembangnya *"BUKA PROTEKSI SALVAGE KASIR"*. Termasuk daftar **jalur simpan mati** `[data work owner 2026-09-19]` |
| `ConnectREST/KonversiKlaimNonLife.xml` | — | sinkronisasi ke sistem inti non-life |
| `Activity/KonversiKlaim_Act.xml` | — | pemanggil konversi |
| `Activity/SendEmailKlaim.xml` · `Activity/SendEmailKlaimRejectClose.xml` | — | pemberitahuan email |
| `RDBList/GetEmailCeding_SQL.xml` | — | alamat ceding lewat fungsi di schema `gl` |
| keenam `ConnectREST/` | — | endpoint **tidak pernah literal**; berasal dari tabel tautan layanan dikunci sepasang kategori |

⚠️ `[terverifikasi]` Kelima *production level* Pega bernilai identik — pemisahan lingkungan sepenuhnya
bergantung isi tabel, dan satu cabang digerbangi **nama node aplikasi**.

## ADR terkait

**ADR-0015** (efek keluar wajib berhasil — transactional outbox) · **ADR-0013** (resolusi endpoint
lewat tabel tautan layanan) · **ADR-0010** (penyimpanan berkas).

⚠️ **Catatan keamanan — dicatat, BUKAN untuk ditindaklanjuti oleh tiket ini.** Kredensial autentikasi
Kasir ikut beredar di dalam berkas ekspor korpus. Ia **sebaiknya diganti sesudah migrasi**. Itu
urusan **tim Kasir**; dicatat di sini semata supaya tidak terlewat. ⛔ Nilainya tidak pernah disalin
ke berkas mana pun di repositori ini.

## Acceptance criteria

- [ ] ⚠️ Efek keluar memakai **outbox transaksional** — dicatat dalam transaksi yang sama dengan perubahan data, dikirim, dan **diulang** bila gagal. **Alasan menyimpang:** Pega tidak punya retry, tidak punya rollback, dan 39 dari 39 penangan exception-nya kosong; kegagalan hilang tanpa jejak *(AC 89 spec)*
- [ ] Test: gagalkan klien luar → pekerjaan **tetap ada di outbox** dan terkirim pada percobaan berikutnya *(AC 90 spec)*
- [ ] `[terverifikasi]` Endpoint **di-resolve saat runtime** dari tabel tautan layanan; **nol URL literal** di kode *(AC 91 spec)*
- [ ] ⚠️ Pemisahan lingkungan **tidak** bergantung pada identitas node aplikasi. **Alasan menyimpang:** di Pega satu cabang digerbangi nama node, sehingga perilaku berubah bila node diganti nama atau ditambah *(AC 92 spec)*
- [ ] ✅ `[data DBA]` **Tidak ada bug pada pembacaan kode respons Kasir.** Kasir mengirim **kedua ejaan sekaligus dengan nilai identik**, sehingga teks status sukses memang muncul selama ini. Label cacat dicabut *(AC 93 spec)*
- [ ] Status pengiriman ke Kasir dapat dilihat pengguna *(AC 94 spec)*
- [ ] ⭐ **BARU 2026-09-19 (ronde 4)** — `[terverifikasi]` Klien Kasir **berautentikasi**, dan kredensialnya berasal dari **konfigurasi runtime bernama** — bukan dipaku di kode, bukan di berkas sumber. Bentuk di Pega: profil autentikasi bernama yang dirujuk rule Connect-REST *(⛔ belum ada nomor AC spec — `spec.md` belum ditambal)*
- [ ] ⭐ **BARU** — `[terverifikasi]` **Hanya klien Kasir** yang berautentikasi. Kelima klien luar lain **tidak**; memberi mereka autentikasi adalah penyimpangan yang tidak diminta
- [ ] Test: kredensial Kasir salah/kosong → pengiriman **gagal terang-terangan** dan pekerjaan **tetap di outbox** (bukan gagal diam-diam)
- [ ] ⚠️ `[data DBA]` **Nomor kasus Kasir disimpan** bersama jejak pengiriman. **Alasan menyimpang:** Pega membuangnya — sensus 329 berkas menemukannya di **nol berkas**, sehingga rantai rekonsiliasi dengan Kasir putus *(AC 99 spec)*
- [ ] ⚠️ `[data DBA]` Kode respons Kasir **dibandingkan sebagai string, atau dinormalisasi eksplisit**. **Alasan menyimpang:** Kasir mengirim nilai bertipe string sedangkan Pega membandingkannya sebagai angka dan bergantung pada konversi diam-diam; Go tidak melakukan konversi itu *(AC 100 spec)*
- [ ] `[data DBA]` Parser menerima **kedua ejaan** kode dan pesan respons, **dan mencatat di log** bila yang datang hanya ejaan yang tidak diharapkan *(AC 101 spec)*

## Perintah verifikasi

```
jalankan test "klien Kasir gagal -> pekerjaan tetap di outbox"
jalankan test "outbox diproses ulang -> pekerjaan terkirim"
jalankan test "nomor kasus Kasir tersimpan bersama jejak pengiriman"
jalankan test "kode respons string '1' -> sukses"
jalankan test "kode respons ' 1' dan '01' -> tidak lolos diam-diam sebagai sukses"
jalankan test "hanya ejaan tak diharapkan yang datang -> tercatat di log"
cari URL literal di kode                                  -> nihil
cari ketergantungan pada nama node aplikasi               -> nihil
```
