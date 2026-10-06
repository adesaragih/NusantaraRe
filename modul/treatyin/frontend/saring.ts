// Saringan kedekatan untuk pemilih "Choose Ceding" / "Choose Source of
// Business".
//
// ⛔ Berdiri SENDIRI, di luar berkas layar. Ia logika murni — nol React, nol
// fetch — dan itu yang membuatnya dapat diuji tanpa merender apa pun.

import type { PilihanWarisan } from './api'

/**
 * Peringkat kedekatan satu baris terhadap ketikan — kecil berarti lebih dekat.
 *
 * ⛔ Peringkat, bukan sekadar cocok-atau-tidak. `saringTerdekat` hanya
 * menampilkan baris pada peringkat TERBAIK yang ditemukan, sehingga ketika ada
 * yang berawalan persis, yang sekadar mengandung tidak ikut muncul.
 *
 * Mengembalikan `null` bila tidak cocok sama sekali.
 */
export function peringkatKedekatan(nama: string, id: string, ketikan: string): number | null {
  // ⛔ Huruf kecil DIPAKSA di sini, bukan diandaikan sudah dilakukan pemanggil.
  // Fungsi yang benar hanya bila pemanggilnya ingat adalah fungsi yang akan
  // salah pada pemanggil kedua.
  const q = ketikan.trim().toLowerCase()
  if (q === '') return 0
  const n = nama.toLowerCase()
  const i = id.toLowerCase()

  if (n === q || i === q) return 0
  if (n.startsWith(q) || i.startsWith(q)) return 1
  // Awal KATA — "adira" menemukan "ASURANSI ADIRA DINAMIKA" tanpa menarik
  // setiap nama yang kebetulan memuat potongan itu di tengah kata.
  if (n.split(/\s+/).some((kata) => kata.startsWith(q))) return 2
  if (n.includes(q) || i.includes(q)) return 3

  // ⛔ Peringkat terakhir: SELURUH kata ketikan harus ada, bukan salah satu.
  // "asuransi adira" dengan syarat salah-satu akan menarik tiap nama
  // berawalan "ASURANSI" — persis keluhan yang membuat aturan ini ada.
  const kata = q.split(/\s+/).filter((k) => k !== '')
  if (kata.length > 1 && kata.every((k) => n.includes(k))) return 4
  return null
}

/**
 * Menyaring ke baris PALING DEKAT saja.
 *
 * ⚠️ Yang dikembalikan hanya satu peringkat — yang terbaik yang ada. Bila satu
 * nama berawalan persis dengan ketikan, baris yang sekadar mengandungnya TIDAK
 * ikut tampil. Itu seluruh pokoknya: daftar pendek yang tepat lebih berguna
 * daripada daftar panjang yang memuat jawabannya di suatu tempat.
 */
export function saringTerdekat(
  daftar: readonly PilihanWarisan[],
  ketikan: string,
): PilihanWarisan[] {
  const q = ketikan.trim().toLowerCase()
  const bernilai: { b: PilihanWarisan; p: number }[] = []
  let terbaik = Number.MAX_SAFE_INTEGER
  for (const b of daftar) {
    const p = peringkatKedekatan(b.nama, b.id, q)
    if (p === null) continue
    if (p < terbaik) terbaik = p
    bernilai.push({ b, p })
  }
  return bernilai.filter((x) => x.p === terbaik).map((x) => x.b)
}

