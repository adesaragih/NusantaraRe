# 04: Layar komite — 92 medan, dua kotak centang terkunci

**Status:** ready-for-agent
**Blocked by:** 00 · 02
**Menutup:** AC 30 · 31 · 32 · 33 · 34 · 35 *(6 AC)* — US 6–11

## Hasil & nilai pengguna

Hari ini Anggota komite **belum punya layar** untuk menilai penyesuaian, dan aturan siapa boleh menyunting apa belum ditetapkan.

Sesudah tiket ini, Anggota komite melihat rekapitulasi penyesuaian **dalam mata uang asli dan dalam rupiah**, dapat menuliskan catatan, dan ⭐ **dua kotak centang usulan hanya dapat disunting penyetuju pertama**.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Layar komite | 92 medan; ⭐ **tepat dua** terkunci bagi jenjang selain yang pertama |
| Dua kotak centang | usul **menutup klaim** dan usul **mencadangkan** |
| ⚠️ Jebakan alih | ⛔ pembanding di Pega bertipe **teks**, padahal pencacahnya **bilangan** |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K4** — ⭐ **Ditiru apa adanya** — usulan tindak lanjut dibuat **penyetuju pertama**; jenjang berikutnya **menilai**, ⛔ tidak mengganti

## Yang harus diuji

- [ ] Penyetuju **pertama** dapat menyunting kedua kotak centang usulan
- [ ] Penyetuju **kedua ke atas** melihat keduanya **terkunci**
- [ ] ⭐ **90 medan lain TIDAK terkunci** oleh aturan itu — ⛔ jangan menguncinya karena salah membaca
- [ ] Layar menampilkan rekapitulasi **dalam mata uang asli** dan **total dalam rupiah**
- [ ] Anggota komite dapat menuliskan **catatan** pada keputusannya
- [ ] ⛔ Pembandingan memakai **bilangan dengan bilangan** — ⚠️ Pega memaafkan teks, sistem baru tidak

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **11** | ⚠️ 29 medan bergantung jenis objek | ⚠️ menahan bentuk tampilan objek |

## Seam & verifikasi

**Seam:** lapisan layanan komite untuk data; layar hanya menampilkan.
1. Buka layar sebagai penyetuju **pertama** ⇒ ⭐ kedua kotak centang **dapat disunting**.
2. Buka sebagai penyetuju **kedua** ⇒ ⭐ keduanya **terkunci**.
3. Periksa 90 medan lain pada penyetuju kedua ⇒ ⛔ **tidak ikut terkunci**.
4. Bandingkan total rupiah dengan hitungan manual ⇒ ⭐ cocok.
