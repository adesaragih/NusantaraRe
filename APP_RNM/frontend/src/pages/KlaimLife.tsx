import { useState, type FormEvent } from 'react'

import { ambilKlaimLife, kodeStatusGalat, tampilUang, type Klaim } from '../services/api'

// ============================================================================
// pages/KlaimLife.tsx — halaman "buka satu klaim Life" (tiket 01 AC-2).
//
// Cara membaca berkas ini bila baru mengenal React:
//   - Komponen = fungsi yang MENGEMBALIKAN tampilan. Bagian yang mirip HTML di
//     dalam `return (...)` disebut JSX; karena berkas ini TypeScript, namanya TSX.
//   - `useState` = kotak penyimpan nilai. Setiap kali nilainya diganti lewat
//     fungsi set-nya (misalnya `setKlaim`), React menggambar ulang tampilan.
//   - `<string>` di belakang useState memberi tahu TypeScript isi kotak itu apa.
//   - `{ ... }` di dalam JSX = "sisipkan nilai JavaScript di sini".
//
// ⛔ Nilai uang tidak pernah disentuh aritmetika dan tidak pernah dibungkus
// Number(). Ia datang sebagai teks desimal dan ditampilkan apa adanya
// (ADR-U-0003, ADR-U-0016). Penjelasannya di services/api.ts bab 3.
// ============================================================================
export default function KlaimLife() {
  // Empat kotak keadaan halaman ini.
  const [id, setId] = useState<string>('') // yang sedang diketik pengguna
  const [klaim, setKlaim] = useState<Klaim | null>(null) // klaim yang berhasil dibaca
  const [galat, setGalat] = useState<string | null>(null) // pesan bila gagal
  const [memuat, setMemuat] = useState<boolean>(false) // sedang menunggu server?

  // Dipanggil saat tombol "Buka" ditekan (form dikirim).
  async function cari(e: FormEvent<HTMLFormElement>) {
    e.preventDefault() // jangan biarkan browser memuat ulang halaman
    setGalat(null)
    setKlaim(null)
    setMemuat(true)
    try {
      setKlaim(await ambilKlaimLife(id.trim()))
    } catch (err: unknown) {
      setGalat(kodeStatusGalat(err) === 404 ? 'Klaim tidak ada.' : 'Gagal membaca klaim.')
    } finally {
      setMemuat(false) // apa pun hasilnya, berhenti "memuat"
    }
  }

  return (
    <section>
      <h2>Klaim Life</h2>

      <form onSubmit={cari}>
        <label htmlFor="idKlaim">Pengenal klaim</label>{' '}
        <input
          id="idKlaim"
          value={id}
          onChange={(e) => setId(e.target.value)}
          placeholder="mis. UJI-KLAIM-1"
        />{' '}
        <button type="submit" disabled={memuat || id.trim() === ''}>
          {memuat ? 'Memuat…' : 'Buka'}
        </button>
      </form>

      {/* `galat && (...)` = tampilkan hanya bila galat terisi. */}
      {galat && <p role="alert">{galat}</p>}

      {klaim && (
        <article>
          <h3>{klaim.nomorKlaim || klaim.id}</h3>
          <dl>
            <dt>Nomor polis</dt>
            <dd>{klaim.nomorPolis || '—'}</dd>
            <dt>Nama bisnis</dt>
            <dd>{klaim.namaBisnis || '—'}</dd>
            <dt>Jumlah baris</dt>
            <dd>{klaim.cacahBaris}</dd>
          </dl>

          {/* `.map` = ulangi blok di bawah untuk SETIAP peserta.
              `key` wajib ada supaya React tahu baris mana yang berubah. */}
          {klaim.peserta.map((p) => (
            <section key={p.id}>
              <h4>
                Peserta {p.nomorSertifikat || p.id}{' '}
                <small>({p.baris.length} baris)</small>
              </h4>

              {p.baris.length === 0 ? (
                <p>Belum ada baris adjustment.</p>
              ) : (
                <table>
                  <thead>
                    <tr>
                      <th>Baris</th>
                      <th>Status</th>
                      <th>Jumlah klaim</th>
                      <th>Nomor akseptasi</th>
                    </tr>
                  </thead>
                  <tbody>
                    {p.baris.map((b) => (
                      <tr key={b.id}>
                        <td>{b.id}</td>
                        <td>
                          {/* Tiket 01 AC-3: status sebagai KATA — bukan nama field,
                              bukan angka. Kode mentahnya (b.kodeStatus) tetap ada di
                              kontrak API bagi yang menelusuri, tetapi tidak dicetak. */}
                          {b.status}
                        </td>
                        <td>{tampilUang(b.jumlahKlaim) || '—'}</td>
                        <td>{b.nomorAkseptasi || '—'}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </section>
          ))}
        </article>
      )}
    </section>
  )
}
