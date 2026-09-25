# 08: Nomor akseptasi — terbit sekali, di tingkat akhir

**Status:** ready-for-agent

**Blocked by:** **06 (tangga maju atau selesai)**

## Hasil & nilai pengguna

Sebagai **bagian akseptasi**, nomor akseptasi terbit **tepat sekali**, saat penyetuju terakhir
menyetujui **tanpa syarat**. Tidak ada nomor ganda, dan nomor **tidak terbit** ketika persetujuannya
masih bersyarat.

*(User story 23, 24 di spec)*

## Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Komite Claim Prop/Activity/KomitePostAdjustment.xml` | `[terverifikasi]` blok **langkah 16** dan anak **16.1**–**16.9**. **Langkah 16** adalah induk, bergerbang **bukan persetujuan bersyarat**, dan merupakan **langkah berulang** |
| idem | `[terverifikasi]` **16.1 · 16.2 · 16.3 · 16.4 di-remark** — tidak berjalan · **16.5** mengambil kode produksi · **16.7** membangkitkan bulan-tahun dan nomor urut · **16.8** merangkai · **16.9** menuliskan nomor dan menandai status akseptasi |
| idem | `[terverifikasi]` **16.9** bergerbang *"penyetuju terakhir dan disetujui"* |

⭐ **Gerbangnya berlapis dua**: di tingkat induk *"bukan persetujuan bersyarat"*, di tingkat anak
*"penyetuju terakhir dan disetujui"*. Itulah yang membuat nomor terbit **sekali saja**.

**Rule pengambil data yang dipanggil** `[terverifikasi]`:

| Langkah | Rule |
| --- | --- |
| 16.1 *(remark)* | `ASM-FW-GISFW-INT-POLICYJSON!RNM!GETTANGGALCLOSING_SQL` |
| 16.4 *(remark)* | `ASM-FW-GCNMFW-INT-V_POLIS!GCNM!GENERATENOACCEPTTREATY` |
| **16.5** | `ASM-FW-GISFW-INT-POLICYJSON!RNM!GETKODEPRODNONLIFE_SQL` |
| **16.7** | `ASM-FW-GISFW-INT-POLICYJSON!RNM!GETSEQUENCENUMBER_SQL` |

⚠️ Dua rule yang langkahnya di-remark **tidak ada di ekspor modul ini** — konsisten: yang dimatikan
tidak ikut diekspor. ⛔ **Jangan dipindahkan.**

## Yang harus diuji

**Diverifikasi oleh:** spec.md AC 31 · 32 · 33 · 34 · 35

- [ ] Nomor akseptasi terbit **tepat sekali** — pada penyetuju terakhir yang menyetujui tanpa syarat.
- [ ] Nomor **tidak terbit** ketika persetujuan **bersyarat**, meski penyetuju terakhir menyetujui.
- [ ] Nomor **tidak terbit** pada penyetuju selain yang terakhir.
- [ ] Menjalankan ulang alur pada kasus yang sudah bernomor **tidak menerbitkan nomor kedua**.
- [ ] Status akseptasi pada baris penyesuaian induk terisi bersama nomornya.
- [ ] ⛔ Kedua rule yang langkahnya di-remark **tidak dipanggil**.

## Butir `[terbuka]` yang menyentuh tiket ini

Tidak ada.

## Seam & verifikasi

Memakai ulang seam Claim Prop.
