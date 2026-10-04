# inti: contoh uji menu DIMIGRASI='0' memakai edmtreatyin

> Cabang lokal `inti/nbtreatyin-uji-menu` @ `0fa4a3e6` (dari `origin/dev`). Merge **sebelum** PR `nbtreatyin`.

## Summary

Uji gigit `frontend/daftar.menuTabel.test.ts` memakai `nbtreatyin` sebagai contoh modul berbaris `DIMIGRASI='0'`.
`M_NAV_MENU.nbtreatyin` kini `DIMIGRASI='1'` (slot 968), jadi contohnya dipindah ke modul yang memang belum menyala.

```diff
-    // Contoh modul berbaris DIMIGRASI='0' - nbtreatyin (nbfacin menyala 02-10-2026, tiket 21).
-    const nbtreatyin = { nama: 'nbtreatyin', kelompok: 'NB Treaty In', halaman: ['nb'], halamanAwal: 'nb' }
-    expect(selisihMenuModul(BERSIH.baris, [...MODUL_FRONTEND, nbtreatyin])).toEqual(["modul frontend nbtreatyin: barisnya DIMIGRASI='0'"])
+    // Contoh modul berbaris DIMIGRASI='0' - edmtreatyin (nbtreatyin menyala 03-10-2026, slot 968).
+    const edmtreatyin = { nama: 'edmtreatyin', kelompok: 'EDM Treaty In', halaman: ['edm'], halamanAwal: 'edm' }
+    expect(selisihMenuModul(BERSIH.baris, [...MODUL_FRONTEND, edmtreatyin])).toEqual(["modul frontend edmtreatyin: barisnya DIMIGRASI='0'"])
```

## Evidence

- **Before:** di atas `origin/dev` + modul nbtreatyin, uji gigit gagal (contoh `nbtreatyin` sudah `DIMIGRASI='1'`).
- **After:** `npx vitest run frontend/daftar.menuTabel.test.ts` di atas `origin/dev`: 14/14 lulus.

## Merge Danger

**Door:** two-way — hanya berkas uji, nol kode produksi.

**Blast Radius:** uji-saja.

Isinya identik dengan bagian `frontend/` commit `88c303de` di cabang `modul/nbtreatyin/implementasi`; merge PR ini
lebih dulu membuat bagian itu tergabung tanpa konflik. Perubahan angka sensus `skemauji.Buka()` 16 → 17
(`modul/claimlife/backend/repository/batasanpemakaian_test.go`) **tidak** ikut di sini — hanya benar bila modul
nbtreatyin ada, jadi tetap di PR modul untuk ditinjau CODEOWNERS.
