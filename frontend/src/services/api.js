import axios from 'axios'

// Satu klien HTTP untuk seluruh frontend.
//
// Nol alamat host di kode: baseURL datang dari env var (ADR-U-0004).
// Kosong berarti sama-asal, yang dipakai bila frontend disajikan di belakang
// proxy yang sama dengan backend.
const baseURL = import.meta.env.VITE_API_BASE_URL ?? ''

export const api = axios.create({
  baseURL,
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' },
})

// ⛔ Uang tidak pernah menjadi Number (ADR-U-0003, ADR-U-0016).
//
// Backend mengirim jumlah uang sebagai TEKS desimal — { amount: "1234.56",
// currency: "IDR" } — persis supaya JSON.parse tidak pernah mengubahnya
// menjadi float64. Helper di bawah menjaga batas itu tetap terlihat; jangan
// pernah membungkusnya dengan Number(), parseFloat(), atau aritmetika JS.

/** Ambil jumlah uang sebagai teks, apa adanya. */
export function jumlahUang(uang) {
  if (uang === null || uang === undefined) return ''
  if (typeof uang.amount === 'number') {
    throw new TypeError(
      'jumlah uang datang sebagai angka JSON; ia harus berupa teks desimal',
    )
  }
  return uang.amount ?? ''
}

/** Tampilkan uang beserta mata uangnya tanpa menyentuh nilainya. */
export function tampilUang(uang) {
  const jumlah = jumlahUang(uang)
  if (jumlah === '') return ''
  return `${jumlah} ${uang.currency ?? ''}`.trim()
}

export async function cekKesehatan() {
  const { data } = await api.get('/healthz')
  return data
}
