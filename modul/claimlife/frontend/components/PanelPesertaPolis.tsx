// View Detail peserta POLIS — Section/DetailPolisLife.xml (InputOSClaimLife b13267), keputusan work owner 11-10-2026.
//
// Tombol View Detail b631 → LoadDataPeserta_Act; tombolnya hilang sesudah ditekan (`ViewDetail.CARI1 !=1` b792). Grid
// `TempDetail.pxResults` b4537, kolom b4700–b6082; paginasi NextPrev ⟪ ⟨ ⟩ ⟫ (b1216, b1489, b2877, b3106) diganti
// komponen `Halaman` inti seperti layar lain. 500 baris per halaman (b572). Peserta versi polis yang SAMA dengan header
// klaim (PL_NUMBER + IDPEGA) — bukan peserta klaim. Baca saja.

import { useState } from 'react'

import { OS_PEGA } from '../labels'
import { tanggalTampil } from './PanelDataPolis'
import { ambilPesertaPolis, type HalamanPesertaPolis } from '../api'
import { Gagal, Halaman, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'

/** Grid satu halaman — dipisah supaya dapat diuji tanpa server. */
export function GridPesertaPolis({ hal, onPindah }: { hal: HalamanPesertaPolis; onPindah: (halaman: number) => void }) {
  if (hal.baris.length === 0) return <Kosong pesan="No insured in this policy version." />
  return (
    <>
      <div className="inbox__gulir">
        <table className="inbox__tabel">
          <thead>
            <tr>
              {OS_PEGA.kolomPolis.map((k) => (
                <th key={k}>{k}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {hal.baris.map((b) => (
              <tr key={b.urut}>
                <td>{b.urut}</td>
                <td>{b.nomorPolis || '—'}</td>
                <td>{b.nomorSertifikat || '—'}</td>
                <td>{b.namaTertanggung || '—'}</td>
                <td>{b.jenisKelamin || '—'}</td>
                <td>{tanggalTampil(b.tanggalLahir)}</td>
                <td>{b.umur || '—'}</td>
                <td>{b.plan || '—'}</td>
                <td>{tanggalTampil(b.tanggalMulai)}</td>
                <td>{tanggalTampil(b.tanggalEfektif)}</td>
                <td>{tanggalTampil(b.tanggalExpired)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <Halaman halaman={hal.halaman} ukuran={hal.ukuran} total={hal.total} onPindah={onPindah} />
    </>
  )
}

export function PanelPesertaPolis({ klaimID }: { klaimID: string }) {
  const [hal, setHal] = useState<HalamanPesertaPolis | null>(null)
  const [memuat, setMemuat] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)

  async function muat(halaman: number): Promise<void> {
    setMemuat(true)
    setGalat(null)
    try {
      setHal(await ambilPesertaPolis(klaimID, halaman))
    } catch (e) {
      setGalat(e)
    } finally {
      setMemuat(false)
    }
  }

  return (
    <section className="register__kartu os__kartu" aria-label={OS_PEGA.viewDetail}>
      {/* b792: tombol hanya tampil sebelum ditekan pertama kali. */}
      {hal === null && (
        <p className="os__simpan">
          <button type="button" className="btn btn--ghost" disabled={memuat} onClick={() => void muat(1)}>
            {OS_PEGA.viewDetail}
          </button>
        </p>
      )}
      {memuat && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}
      {hal !== null && (
        <GridPesertaPolis
          hal={hal}
          onPindah={(n) => {
            void muat(n)
          }}
        />
      )}
    </section>
  )
}
