# PROMPT — GILIRAN 12 *(**sesi BARU** di `OUTPUT_HASIL_RNM`, cabang `main` @ `43c54de` atau lebih baru)*: **unduh dokumen 401 → bl dokumen lengkap → sensus remark lintas tiga modul → uji asap baca-saja ke DEV → status tiket**

> Hanya konteks **Claim Life, PremiumList Life, Komite Claim Life**. GILIRAN-3 *(A/B/C)*, 4–11 tetap rujukan; aturan berhenti GILIRAN-6
> dan sesi baru tiap giliran *(induk §0.5)* berlaku. Berkas `tco_*` dan worktree Treaty Contract Out **tidak** disentuh.

## 0. KEADAAN AWAL — DIVERIFIKASI ASISTEN 28-09-2026

| Klaim laporan GILIRAN-11 | Diperiksa ulang | Hasil |
| --- | --- | --- |
| 11 commit `ff1893c` → `43c54de`; berhenti karena selesai | 10 commit milik sesi ini di `main` *(di luar commit Treaty)*; pohon bersih kecuali `App.tsx` milik work owner | ✅ |
| delapan langkah `SaveOutStandingLife_Act` ter-remark | `pyStepsBlockName` bernilai `//` di **delapan** tempat: b3495, b4632, b5009, b6178, b7074, b7293, b7512, b10649 | ✅ temuan sah; gerbang STNC dan dokumen lengkap dibuang dengan benar |
| `PeriksaDokumenLengkap` tanpa pemanggil | `services/dokumen.go:127`, nol pemanggil produksi | ✅ → **bl** |
| panduan uji | `PANDUAN-UJI-LAYAR-TIGA-MODUL.md` *(±92 KB)* | ✅ |
| uji | dijalankan ulang asisten di `main`: **640 PASS · 0 FAIL**; dengan tag `db` **38 SKIP**; vitest **390**; vet, gofmt, tsc bersih | ✅ |
| ⛔ **tautan unduh dokumen 401** | `components/claimlife/PanelDokumenPeserta.tsx:157` masih `<a href={tautanDokumen(b.id)} download>`; backend `pelaku.go` membaca identitas **hanya** dari header `X-Pelaku`/`X-Peran` → setiap unduhan dijawab 401 | ⛔ **paket 0** |
| ⚠️ memori "cetak `pyStepsBlockName`" | pelajaran ini berlaku ke **seluruh** activity yang sudah ditiru, bukan hanya `SaveOutStandingLife_Act` | ⚠️ **paket 2** |

## 1. KEPUTUSAN GILIRAN INI

| Butir | Isi |
| --- | --- |
| **bl** — OQ-N6 `[DIPUTUSKAN; veto work owner]` | Kelengkapan dokumen per kategori **ter-remark** di XML → di sistem lama ia **tidak pernah berlaku**. Mengikuti XML dan aturan "kode mati dibuang": `PeriksaDokumenLengkap` **dibuang** beserta ujinya; tiket 03 dan `spec.md` mencatatnya bertanggal; OQ-N6 ditutup. Bila bisnis menghendaki gerbang itu, ia keputusan **baru**, bukan replikasi |
| OQ-N5 *(sumber `ProdDateTime` / keputusan Komite #4 untuk Claim Life)*, langkah 11.10 *(DOL kosong)*, OQ-N1…N4, OQ-M1…M7, OQ-PL-09/10/11, OQ-K-04a/05/05b | **tetap untuk work owner** — tidak diputuskan executor; perilaku sekarang *(Arasapas ditahan beralasan; DOL kosong ditolak)* dipertahankan |

## 2. URUTAN — satu commit per paket

| # | Paket | Isi |
| ---: | --- | --- |
| 0 | **Unduh dokumen membawa identitas** | `PanelDokumenPeserta.tsx:157`: ganti `<a href download>` dengan unduhan lewat klien `fetch` yang sama *(header `X-Pelaku`/`X-Peran` ikut)*, hasil disimpan sebagai `Blob` dengan nama berkas asli; tombol `View Office Online` b3502 memakai jalur yang sama; **uji kontrak dua sisi**: handler menolak tanpa header, klien selalu mengirimnya; penjaga statik: nol `href` ke `/api/` di seluruh `frontend/src`. Periksa pola yang sama di layar PremiumList dan Komite *(lampiran, unduhan)*. Commit `fix: unduh dokumen membawa identitas pelaku` |
| 1 | **bl** | §1. Commit `claim-life: bl — PeriksaDokumenLengkap dibuang, ter-remark di XML` |
| 2 | **Sensus remark lintas tiga modul** | untuk **setiap** activity yang sudah ditiru di tiga modul *(daftar dari PARITAS masing-masing)*: cetak langkah beserta `pyStepsBlockName`; tiap langkah ber-`//` yang ternyata **ditiru** di kode → dibuang *(atau dicatat bila penirunya disengaja sebagai keputusan tercatat)*; tiap langkah hidup yang **tidak** ditiru → dicatat sebagai celah. Hasil = tabel per modul di PARITAS bab *"Sensus remark 28-09-2026"* + ralat tiket bertanggal. Satu commit per modul |
| 3 | **Uji asap baca-saja ke DEV** | nyalakan backend `main` dengan `.env` work owner *(`. .\muat-env.ps1` atau setara)*, lalu **hanya `GET`** ke setiap rute baca tiga modul *(daftar dari `handlers.go`)* dengan header stub; catat kode jawaban dan cacah baris, **tanpa** menyalin isi baris; `POST`/`PUT`/`DELETE` **tidak** dijalankan *(menulis ke DEV = persetujuan work owner)*. Hasil ditambahkan ke `PANDUAN-UJI-LAYAR-TIGA-MODUL.md` bab *"Uji asap baca-saja"*; rute yang menjawab 5xx diperbaiki di giliran ini |
| 4 | **Status tiket** | tiket yang tersentuh paket 0–3 diperbarui *(Status, `[x]` berbukti, ralat)* |

## 3. ATURAN · LAPORAN · TELEMETRI

Sama dengan GILIRAN-4 §2–§3, GILIRAN-6, GILIRAN-9 §2; tambahan: **setiap pembacaan activity mencetak `pyStepsBlockName`**.
Laporan: baris pertama alasan berhenti; tabel **paket → commit → rule XML → rute/komponen**; ralat tiket; OQ; angka uji tiap commit
dengan dan tanpa tag `db`; hasil uji asap *(rute → kode jawaban → cacah)*; bab **TELEMETRI EKSEKUSI**.

---

*Disusun 28 September 2026 sesudah verifikasi `5ac6571..43c54de` (uji atas pohon ter-commit; delapan `pyStepsBlockName //` di
`SaveOutStandingLife_Act.xml`; `PanelDokumenPeserta.tsx:157`; `services/dokumen.go:127`).*
