# 10: Jalur Fac Retro — jalur yang melompati penyiapan wewenang

**Status:** ready-for-agent
**Blocked by:** 00 · 03
**Menutup:** AC 71 · 72 · 73 · 74 · 75 · 76 *(6 AC)* — US 37 · 38

## Hasil & nilai pengguna

Hari ini Penyesuaian bertanda **Fac Retro** **melompati penyiapan wewenang penyetujuan** — ⛔ dan itu **tidak terlihat oleh siapa pun yang membaca aturan**, sebab ia terkubur di dalam detail.

Sesudah tiket ini, ⭐ Jalur Fac Retro **dikenali sebagai jalur tersendiri**, disebut terang-terangan, sehingga siapa pun yang membaca aturan **melihat jalur itu ada** — ⛔ bukan menemukannya sesudah sesuatu terjadi.

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

## Yang harus diuji

- [ ] ⭐ Penyesuaian bertanda Fac Retro mengikuti **jalur tersendiri** yang **melompati penyiapan wewenang**
- [ ] ⚠️ Nilai penanda jenis treaty tetap **tetapan di dalam kode** ⇒ ⛔ **setiap perubahan daftar treaty menuntut RILIS perangkat lunak**
- [ ] ⛔⛔ **Spec dan tiket ini menyebut jalur itu SECARA TERBUKA** sebagai jalur yang melompati pemeriksaan wewenang — ⭐ supaya ia **terlihat**, bukan terkubur
- [ ] ⚠️ Ketiga jalur pemicu **ditelusuri** sebelum dibangun — ⛔ termasuk yang **tanpa gerbang**
- [ ] ⚠️ Pencampuran kode master dengan teks harfiah pada medan jenis treaty **wajib beres** sebelum medan itu boleh menjadi kolom bernilai kode

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
