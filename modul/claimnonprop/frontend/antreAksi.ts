// Disalin dari `modul/claimprop/frontend/antreAksi.ts` (pola, bukan impor): keputusan work owner yang disebut di bawah
// adalah keputusan layar Claim Prop yang ditiru Claim Non Prop (prompt Claim Non Prop tahap 1).
// Antrean aksi layar kasus (laporan work owner 09-10-2026: "klik send komite, ga muncul pop up, harus klik lagi").
// Klik tombol sesudah mengetik di medan beraksi (MakeLowercase, hitung ulang) memicu blur medan itu lebih dulu, dan aksi
// medannya berjalan ke server. Dulu tombol nonaktif selama sibuk sehingga kliknya hilang; kini aksi kedua diantre dan
// dijalankan sesudah aksi pertama selesai, di atas halaman hasilnya (urutan sama dengan Pega: refresh medan, lalu
// tombol). Hanya satu aksi menunggu (yang terakhir); klik ganda aksi yang sama diabaikan.

export interface AksiDiminta {
  aksi: string
  indeks: number
}

/** `jalan` = kirim sekarang; `antre` = jalankan sesudah aksi berjalan selesai; `abaikan` = klik ganda. */
export function putuskanAksi(berjalan: AksiDiminta | null, minta: AksiDiminta): 'jalan' | 'antre' | 'abaikan' {
  if (berjalan === null) return 'jalan'
  if (berjalan.aksi === minta.aksi && berjalan.indeks === minta.indeks) return 'abaikan'
  return 'antre'
}
