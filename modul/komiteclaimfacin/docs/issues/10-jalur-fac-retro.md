# 10: Jalur Fac Retro — jalur yang melompati penyiapan wewenang

**Status:** ready-for-agent
**Blocked by:** 00 · 03
**Menutup:** AC 71 · 72 · 73 · 74 · 75 · 76 *(6 AC)* — US 37 · 38

## Hasil & nilai pengguna

Hari ini Penyesuaian bertanda **Fac Retro** **melompati penyiapan wewenang penyetujuan** — ⛔ dan itu **tidak terlihat oleh siapa pun yang membaca aturan**, sebab ia terkubur di dalam detail.

Sesudah tiket ini, ⭐ Jalur Fac Retro **dikenali sebagai jalur tersendiri**, disebut terang-terangan, sehingga siapa pun yang membaca aturan **melihat jalur itu ada** — ⛔ bukan menemukannya sesudah sesuatu terjadi.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Hari ini Penyesuaian bertanda **Fac Retro**
> **melompati penyiapan wewenang penyetujuan** — ⛔ dan itu **tidak terlihat oleh siapa pun yang membaca aturan**, sebab
> ia terkubur di dalam detail."* → yang dilompati **hanya perluasan tangga**, bukan wewenang:
>
> - `ApprovalKomite_Act` L2.2 menyetel `Local.Retro := 1` bila ada adjustment KMT ber-`IsFacRetro == 1`, lalu L3
>   `Exit-Activity` melewati L4–L8 (perluasan, pita SPV B, hitung ulang `KomiteLoop`). **Tangga retro sudah dibentuk sisi
>   klaim** saat kelahiran (`RosterKomiteCalon`: batas = `ValueAdjustment`, tanpa saringan batas bila melampaui
>   `LIMIT_TOP` Technic Div. Head).
> - Di tingkat akhir, retro juga melewati **konversi** (`KomitePost_Adjustment` S17) dan **log "AKSEPATSI"** (S19),
>   keduanya bergerbang `Local.FacRetro == 1 [T=3]`.
> - **Wewenang tetap ditegakkan** untuk retro: `SetValueKomite` S15 `SetProteksiSubmiteKomite` tetap dipanggil sesudah
>   S14, dan di sistem baru `Pemegang` (tiket 03) berlaku untuk semua kasus.
> - `IsFacRetro` / `IsFacretro` objek, item, dan adjustment ditulis ulang di tingkat akhir oleh `SaveAcceptation_KMT`
>   L2.1–L4 (spreading item `TreatyType == "10015"` → 1, selain itu 0), lewat kontrak.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Penanda Fac Retro | jenis treaty tertentu menandai penyesuaian sebagai Fac Retro |
| ⛔ Akibatnya | ⛔ penyiapan wewenang penyetujuan **keluar seketika** bila penanda itu menyala |
| ⚠️ Tiga jalur pemicu | ⛔ `[terverifikasi]` penanda dipicu **setidaknya tiga jalur**, ⚠️ salah satunya **tanpa gerbang sama sekali** |
| ⚠️ Medan jenis treaty | ⛔ **mencampur kode master dengan teks harfiah** |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K7** — ⚠️ **DITIRU APA ADANYA** — ⛔ **BERBEDA dari rekomendasi asisten**, yang mengusulkan menjadikan penanda jenis treaty sebagai **data**. ⛔ **Disengaja**

> ⛔ **RALAT 10-10-2026.** Butir lamanya dikutip utuh, tidak dihapus: *"**K7** — ⚠️ **DITIRU APA ADANYA** — ⛔
> **BERBEDA dari rekomendasi asisten**, yang mengusulkan menjadikan penanda jenis treaty sebagai **data**. ⛔
> **Disengaja**"* → `keputusan-sebelum-to-spec.md` mencatat Q3 sebagai **belum diputuskan**, jadi asal K7 tidak
> tertelusur di berkas yang boleh disunting (lihat RALAT di bab Q3 berkas itu). Yang mengikat kini: XML (L3, S17, S19,
> `SaveAcceptation_KMT` L2.1) + **KCF-02** ("Kasus Fac Retro melewati perluasan (L3); tangganya sudah dibentuk sisi
> klaim"). Nomor K7 juga dipakai `STRUKTUR-TABEL-CLAIM-FACIN.md` §5d untuk keputusan lain (jejak komite satu tabel);
> keduanya tidak berhubungan.

## Yang harus diuji

- [ ] ⭐ Penyesuaian bertanda Fac Retro mengikuti **jalur tersendiri** yang **melompati penyiapan wewenang**
- [ ] ⚠️ Nilai penanda jenis treaty tetap **tetapan di dalam kode** ⇒ ⛔ **setiap perubahan daftar treaty menuntut RILIS perangkat lunak**
- [ ] ⛔⛔ **Spec dan tiket ini menyebut jalur itu SECARA TERBUKA** sebagai jalur yang melompati pemeriksaan wewenang — ⭐ supaya ia **terlihat**, bukan terkubur
- [ ] ⚠️ Ketiga jalur pemicu **ditelusuri** sebelum dibangun — ⛔ termasuk yang **tanpa gerbang**
- [ ] ⚠️ Pencampuran kode master dengan teks harfiah pada medan jenis treaty **wajib beres** sebelum medan itu boleh menjadi kolom bernilai kode

> ⛔ **RALAT 10-10-2026.** Baris dan butir lamanya dikutip utuh, tidak dihapus: *"⛔ Akibatnya | ⛔ penyiapan wewenang
> penyetujuan **keluar seketika** bila penanda itu menyala"*, *"⭐ Penyesuaian bertanda Fac Retro mengikuti **jalur
> tersendiri** yang **melompati penyiapan wewenang**"* dan *"⛔⛔ **Spec dan tiket ini menyebut jalur itu SECARA
> TERBUKA** sebagai jalur yang melompati pemeriksaan wewenang — ⭐ supaya ia **terlihat**, bukan terkubur"* → yang
> keluar seketika adalah **perluasan tangga** (`ApprovalKomite_Act` L3), bukan pemeriksaan wewenang; lihat RALAT di
> kepala tiket. Uji yang benar: retro ⇒ List of Committee = tangga sisi klaim tanpa perluasan; tingkat akhir ⇒ nol
> konversi dan nol log "AKSEPATSI"; bukan anggota workbasket ⇒ tetap 403. Nilai `"10015"` tetap konstanta bernama di
> kode (bagian K7 yang berdiri). Tiga jalur pemicu dan pencampuran kode master milik sisi klaim (butir 28).
>
> ⚠️ **OQ (maksud XML tidak pasti):** `KomitePost_Adjustment` S25 (jalur kasir) bergerbang `Local.FacRetro == 1 [T=3]`
> dengan catatan *"JIKA RETRO KONVESI WAKTU PRINT DLA"*, tetapi ber-pre=false, jadi gerbangnya tidak pernah dinilai dan
> kasir tetap berjalan untuk retro. Dibangun ikut XML (kasir jalan); apakah maksudnya retro dilewati = pertanyaan untuk
> work owner.

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **28** | ⚠️ **Arti nilai penanda jenis treaty**, dan apakah **ketiga jalur pemicu** memang dikehendaki | ⚠️ tidak menahan pembangunan — ⭐ menahan **kepastian bahwa seluruh jalurnya sudah dikenali** |

## Seam & verifikasi

**Seam:** lapisan layanan komite — penyiapan wewenang penyetujuan.
1. Jalankan penyesuaian **bertanda Fac Retro** ⇒ ⭐ penyiapan wewenang **dilewati**.
2. Jalankan penyesuaian **biasa** ⇒ ⭐ penyiapan wewenang **dijalankan**.
⚠️ Picu penanda lewat **ketiga jalur** ⇒ ⭐ ketiganya menghasilkan perilaku **yang sama**.
3. Periksa dokumentasi jalur ⇒ ⭐ ia **disebut sebagai jalur tersendiri**, bukan catatan kaki.
