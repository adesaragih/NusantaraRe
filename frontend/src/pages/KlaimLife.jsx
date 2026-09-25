import { useState } from 'react'

import { ambilKlaimLife, tampilUang } from '../services/api.js'

// Tiket 01 AC-2: halaman menampilkan klaim dan daftar barisnya, masing-masing
// dengan statusnya.
//
// ⛔ Nilai uang tidak pernah disentuh aritmetika JS dan tidak pernah dibungkus
// Number(). Ia datang sebagai teks desimal dan ditampilkan apa adanya
// (ADR-U-0003, ADR-U-0016).
export default function KlaimLife() {
  const [id, setId] = useState('')
  const [klaim, setKlaim] = useState(null)
  const [galat, setGalat] = useState(null)
  const [memuat, setMemuat] = useState(false)

  async function cari(e) {
    e.preventDefault()
    setGalat(null)
    setKlaim(null)
    setMemuat(true)
    try {
      setKlaim(await ambilKlaimLife(id.trim()))
    } catch (err) {
      setGalat(err?.response?.status === 404 ? 'Klaim tidak ada.' : 'Gagal membaca klaim.')
    } finally {
      setMemuat(false)
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
                          {/* Tiket 01 AC-3: status sebagai kata — bukan nama field,
                              dan BUKAN angka. Kode mentahnya tetap ada di kontrak API
                              bagi yang menelusuri, tetapi tidak pernah dicetak ke layar. */}
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
