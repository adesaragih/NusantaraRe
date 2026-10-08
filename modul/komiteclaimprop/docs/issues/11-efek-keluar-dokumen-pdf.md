# 11: Efek keluar — dokumen PDF akseptasi

**Status:** dibangun SEBAGIAN — PDF dibangun, unggah + `DOCUMENT_CLAIM` belum disambung *(RALAT 08-10-2026; status lama: `ready-for-agent`)*

> **RALAT 08-10-2026** — implementasi satu modul (prompt `_brief/PROMPT-IMPLEMENTASI-MODUL-KOMITE-CLAIM-PROP.md` §7). Kalimat lama tetap di bawah, dikutip di sini:
>
> - `IsPrintAccept := 1` ditulis; ~~berkas PDF tidak dikarang: stream `FILEAcceptanceNote` tidak diekspor (OQ, sama dengan Claim Prop).~~ RALAT kedua 08-10-2026 (stream diekspor work owner): efek outbox `dokumen-akseptasi` diantre di S21 (hanya produksi); markup "ACCEPTED CLAIM INSURANCE" dirakit VERBATIM saat efek dikirim (`SusunDokumenAkseptasi`), nama berkas "Persetujuan Klaim   AcceptNo <no>.pdf", kategori AcceptanceNote, folder Claim. ~~`HTMLToPDF` (mesin PDF platform Pega) tanpa padanan di go.mod → pelaksana berhenti `ErrPenyimpananBelumDisetujui` (OQ-KCP-07).~~ RALAT ketiga 08-10-2026 (OQ-KCP-07 "A"): PDF A4 digambar `github.com/go-pdf/fpdf` dari halaman TempAcceptedNo (`models.PDFAcceptanceNote`); unggah Google Storage + `DOCUMENT_CLAIM` (S12-S13) belum disambung - pelaksana berhenti `ErrPenyimpananBelumDisetujui`. AC 85-86 (pengenal berkas nanodetik + UUID, jenis berkas diterima) ditempel di sini, menunggu unggahan PDF.


**Blocked by:** **08 (nomor akseptasi)**

## Hasil & nilai pengguna

Sebagai **bagian akseptasi**, dokumen akseptasi **dibuat otomatis** sebagai PDF dan tersimpan, jadi
tidak perlu disusun manual.

*(User story 26 di spec)*

## Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Komite Claim Prop/Activity/KomitePostAdjustment.xml` | `[terverifikasi]` **langkah 21** memanggil pembuat dokumen akseptasi |
| `Komite Claim Prop/Activity/HTMLToPDF.xml` | ⚠️ `[terverifikasi]` **rule bawaan Pega**, berkelas dasar — **bukan buatan Nusantara Re**. ⛔ **Tidak perlu dipindahkan**; yang dipindahkan adalah **pemakaiannya** |
| `Komite Claim Prop/Activity/InsertDocument_Act.xml` | `[terverifikasi]` menyimpan baris dokumen lalu memanggil pengunggah berkas |
| `Komite Claim Prop/Activity/InsertGoogleStorage_Act.xml` | `[terverifikasi]` mengunggah berkas ke penyimpanan luar, lalu mencatatnya ke tabel penyimpanan berkas |
| `Komite Claim Prop/Activity/GetLinkService.xml` | `[terverifikasi]` mengambil **alamat layanan** dari tabel alamat, disaring **dua kunci** — kategori dan sub-kategori — sehingga hasilnya satu baris |

## ⚠️ Sikap terhadap kegagalan — **berhenti dengan galat**

`[terverifikasi]` Rantai PDF **melempar galat dan berhenti** bila gagal — **berbeda** dari sikap
pengiriman ke Kasir, yang melompat dan melanjutkan (tiket 12). **Empat titik penanganannya:**

| Di mana | Bila gagal |
| --- | --- |
| pembuat PDF, langkah 6 | **melempar galat** bila markup kosong |
| pembuat PDF, langkah 9 | **mencatat galat dan berhenti** bila pembuat tidak menghasilkan isi |
| tiga langkah pengubah berkas jadi teks | **melempar galat** bila lampiran gagal |

⭐ **Modul ini punya dua sikap berbeda terhadap kegagalan, dan keduanya disengaja.**

⚠️ `[terverifikasi]` **Pengunggahan berkas TIDAK punya penanganan gagal.** Catatan pengembangnya
sendiri berbunyi *"FIX ERROR HANDLING"*, tetapi **jejaknya tidak terbaca** di keempat keluarga wadah
yang sudah disisir habis. ⛔ Jangan ditebak.

## ⚠️ TITIK YANG SENGAJA DIUBAH — urutan terhadap penyimpanan

`[keputusan work owner]` 2026-09-18. **Ini perubahan sadar, bukan peniruan.**

`[terverifikasi]` Di Pega **kedelapan efek keluar berjalan SEBELUM penyimpanan**. Penyimpanan ada di
**langkah 41**, sedangkan pengiriman ke Kasir di **34** dan email di **35**. **Nol yang sesudah.**

⚠️ Akibatnya di Pega: bila penyimpanan gagal, **uang sudah dikirim, email sudah sampai, dokumen
sudah dibuat, dan tiga tabel log sudah terisi** — sementara kasusnya sendiri tidak tersimpan.

> `[keputusan work owner]` — **URUTAN INI TIDAK DITIRU.** Fakta di atas dicatat sebagai **fakta**,
> **tidak mengikat rancangan**. Urutannya **akan disesuaikan**.
>
> ⛔ **Jangan mengunci urutan Pega sebagai syarat.** ⛔ Urutan penggantinya **belum diputuskan**,
> dan **tidak ditetapkan di tiket ini**.

## Yang harus diuji

**Diverifikasi oleh:** spec.md AC 48 · 49 · 50 · 51 · 52 · 69 · 70 · 73

- [ ] Nomor akseptasi terbit → dokumen PDF **dibuat** dan **tersimpan**, dan baris dokumen tercatat.
- [ ] Kegagalan pembuatan PDF **menghentikan dengan galat yang terlihat**, bukan diam-diam dilewati.
- [ ] Alamat layanan penyimpanan diambil dengan **dua kunci penyaring**, bukan baris pertama tanpa
      saringan.
- [ ] ⚠️ Kegagalan **pengunggahan berkas** — perilakunya mengikuti apa yang ada, dan bila memang
      tidak tertangani, **itu didokumentasikan**, bukan ditambal diam-diam.

## Butir `[terbuka]` yang menyentuh tiket ini

- **Pengunggahan berkas dan email tanpa penanganan gagal** — belum dibawa ke work owner.
- **Catatan pengembang "FIX ERROR HANDLING" yang tidak berjejak.**
- **Urutan terhadap penyimpanan** belum ditetapkan.
- **Baris mana di tabel alamat layanan yang menunjuk lingkungan uji** belum ditetapkan.

## Seam & verifikasi

Efek keluar diuji dengan **layanan sungguhan di lingkungan uji terpisah**
`[keputusan work owner]` 2026-09-18.
