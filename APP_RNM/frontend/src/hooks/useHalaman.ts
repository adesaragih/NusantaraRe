/**
 * Keadaan halaman untuk daftar berhalaman.
 *
 * Pola yang sama muncul di lima layar: satu `useState` untuk nomor halaman,
 * dan setiap penyaring baru harus mengembalikannya ke satu — kalau tidak,
 * pengguna mengubah saringan sambil berada di halaman tujuh dan melihat
 * daftar kosong yang terbaca seperti "tidak ada data".
 *
 * Hook ini ada karena kelalaian itu sudah terjadi: kembali-ke-satu mudah
 * dilupakan ketika ia ditulis ulang di tiap layar, dan gagalnya tidak
 * menghasilkan satu pun galat.
 */

import { useCallback, useState } from "react";

export function useHalaman(awal = 1) {
  const [halaman, setHalaman] = useState(awal);

  /** Pasang penyaring baru: nilainya berubah DAN halaman kembali ke satu. */
  const saring = useCallback(
    <T,>(set: (v: T) => void) =>
      (v: T) => {
        set(v);
        setHalaman(1);
      },
    [],
  );

  /**
   * Kembali ke halaman pertama.
   *
   * ⚠️ DIMEMOISASI (ronde 271 C.20, TV2-410). Sebelumnya ia arrow BARU
   * setiap render, sementara `saring` di atas — saudaranya, dikembalikan
   * dari hook yang sama — sudah memakai useCallback. Ketimpangan itu
   * bukan kosmetik: `ReportRealisasi.tsx:74` mendaftarkan `keAwal` sebagai
   * dependency sebuah useEffect, sehingga efeknya berjalan ulang pada
   * SETIAP render, menyetel state baru, yang memicu render berikutnya.
   *
   * ⚠️ Ini BUKAN klaim bahwa ia sebab ke-14 penolakan `realization-report`
   * pada log 23-09 09:30 (TV2-342). Log itu tidak direproduksi. Predikat
   * `retry` sudah mengecualikan NOT_PROVISIONED (`services/api.ts`
   * `bolehUlang`), jadi pengulangannya jelas bukan dari sana — dan ini
   * cacat terukur yang berdiri sendiri, diperbaiki karena benar, bukan
   * karena terbukti menyebabkan log itu.
   */
  const keAwal = useCallback(() => setHalaman(1), []);

  return { halaman, setHalaman, saring, keAwal };
}
