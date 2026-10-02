---
status: accepted
tanggal: 2026-09-14
sumber: grilling Ronde 2 Q12 (`.scratch/claim-life/grilling-ronde-2.md`), keputusan work owner
---

# Efek keluar Claim — Life dijalankan asinkron, tidak memblokir alur, dengan antre-ulang

Keempat efek keluar Claim — Life bersifat **asinkron**: kegagalan **dicatat**, **tidak memblokir**
alur klaim, dan disediakan **antre-ulang (retry)**. Perilaku "tidak memblokir" dipertahankan dari
sistem lama; **keandalannya ditambah**.

## Efek keluar yang tercakup `[terverifikasi]`

| Efek | Rule |
| --- | --- |
| Unggah berkas ke Google Storage | `Claim Life/Activity/InsertGoogleStorage_Act.xml` |
| Email | `Claim Life/Activity/SendEmailKlaimLF.xml` |
| Arasapas | `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml` — **satu-satunya salinan di korpus** (OQ-035) |
| Konversi ke produksi | `Claim Life/ConnectREST/convertJsonNusareToProductionClaimLife.xml` (`pyServiceName` = `convertJsonNusareToProductionClaimLife`, `pyBaseURLSelectionType = SETTING`) |

## Mengapa "tidak memblokir" adalah paritas, bukan perubahan

`[terverifikasi]` Di `Claim Life` ditemukan **pencatatan**, bukan gerbang keberhasilan:
`Claim Life/RDBList/InsertLogServiceClaim.xml` → `INSERT INTO pooldata.monitoring_klaim_log`.

Bandingkan konteks facultative, yang **punya** gerbang keberhasilan eksplisit —
`When/IsSuccessHitService.xml` menguji `.StatusService.StsKonversiFacIn` / `StsKonversiFacOut`
dan mengarahkan alur ke `SetToInbox_ACT` bila gagal. **Pola itu tidak ada di `Claim Life`.**

Jadi keputusan "tidak memblokir" **mempertahankan perilaku terbaca**, bukan mengarangnya.
Yang ditambahkan adalah **antre-ulang**, yang di Pega tidak terbaca ada.

## Consequences

- Alur klaim (Register → Outstanding → Medical Check → Claim Analis) **tidak pernah tertahan** oleh
  kegagalan layanan luar.
- Kegagalan harus **terlihat** — jejaknya masuk ke jalur audit **ADR-0007**, tidak cukup hanya ke
  log layanan.
- Efek keluar digerbangi flag lingkungan **ADR-0005**; di non-production ketiganya tidak berjalan,
  sehingga antre-ulang pun tidak aktif di sana.
- Alamat layanan **di-lookup runtime** dari `M_LINK_SERVICE` (**ADR-0013**, menggantikan ADR-0004);
  kegagalan **kunci kategori tidak ditemukan** harus dibedakan dari kegagalan jaringan agar
  antre-ulang tidak berputar sia-sia.
- **Konsekuensi yang perlu disadari:** karena tidak memblokir, sebuah klaim dapat mencapai status
  akhir sementara efek keluarnya masih tertunda. Apakah keadaan itu boleh dianggap "selesai"
  adalah pertanyaan bisnis yang **belum dijawab**.

## Pengecualian: Komite Claim Life (2026-09-15)

⚠️ ADR ini berlaku untuk **Claim — Life**. Konteks **Komite Claim Life** memutuskan sebaliknya —
efek keluarnya **wajib berhasil** (transactional outbox, at-least-once) — karena ia memuat integrasi
**Kasir** (pembayaran) yang tidak ada di Claim Life. Lihat **ADR-0015**.

Keduanya berlaku bersamaan pada konteks masing-masing; ini **bukan** pembatalan.

## OQ yang masih terbuka dan menyentuh ADR ini

| OQ | Yang belum diketahui |
| --- | --- |
| **OQ-025 / OQ-035** | `serviceInsertArasapasClaimLife_act` adalah satu-satunya salinan di korpus dan dipanggil juga oleh `Komite Claim Life` — kepemilikan dan kontraknya belum ditetapkan |
| **OQ-047** | Isi tabel `M_LINK_SERVICE` — daftar endpoint yang sebenarnya tidak diketahui |
| **OQ-002** | `POOLDATA.GET_TOKEN_STORAGE` (dipakai jalur unggah berkas) tanpa body di korpus |
| **OQ-029** | Kondisi asli `IsPEGAPROD` tidak terbaca — gerbang lingkungan menghindarinya, bukan menjawabnya |
