// Popup `View Old Policy` - `ViewOldPolicy_EDM` (QR), `_QP`, `_TP`, `_TR` (tiket 03).
//
// Tombolnya di layar yatim `ShowLifePremiumSummary_EDM` (b64965/b65522/b66083/b66640) dipindah ke kepala
// `InputEDMLife` (RALAT R02, OQ-EDM-004): satu tombol, popup sesuai `.Type` kasus. Varian QP/TP/TR di Pega
// kosong karena wadah luarnya bersyarat `.Type = 'QR'` (R03) - di sini grid dalamnya diisi (spec AC 9).

import { useEffect, useState } from 'react'

import { Gagal, Halaman, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilPolisLama, type PolisLamaEDM } from '../api'
import { POLIS_LAMA_PESERTA, POLIS_LAMA_REKAP, type TipeEDM } from '../kolomKorpus'
import { POLIS_LAMA_EDM, UMUM_EDM } from '../labels'
import { UKURAN_HALAMAN_EDM } from '../tampilan'
import { TabelKorpus } from './TabelKorpus'

/** Tipe yang dikenal popup; tipe lain jatuh ke varian induk `QR` (seperti tombol `.Type=='QR'`). */
export function tipePopup(t: string): TipeEDM {
  return t === 'QP' || t === 'TP' || t === 'TR' ? t : 'QR'
}

export default function PolisLama({ kasusId, tipe, onTutup }: { kasusId: string; tipe: string; onTutup: () => void }) {
  const [halaman, setHalaman] = useState(1)
  const [isi, setIsi] = useState<PolisLamaEDM | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const t = tipePopup(tipe)

  useEffect(() => {
    let batal = false
    ambilPolisLama(kasusId, halaman)
      .then((p) => {
        if (!batal) setIsi(p)
      })
      .catch((e: unknown) => {
        if (!batal) setGalat(e)
      })
    return () => {
      batal = true
    }
  }, [kasusId, halaman])

  return (
    <Modal judul={POLIS_LAMA_EDM.viewOldPolicy} onTutup={onTutup} labelBatal={UMUM_EDM.tutup} penuh>
      {galat !== null && <Gagal galat={galat} />}
      {isi === null && galat === null && <Memuat pesan={UMUM_EDM.memuat} />}
      {isi !== null && isi.peserta.length === 0 && <Kosong pesan={UMUM_EDM.kosong} />}
      {isi !== null && isi.peserta.length > 0 && (
        <>
          <TabelKorpus kolom={POLIS_LAMA_PESERTA[t]} baris={isi.peserta.map((p) => p.nilai)} kunci={(_, i) => isi.peserta[i]?.id ?? String(i)} />
          <Halaman halaman={halaman} ukuran={UKURAN_HALAMAN_EDM} total={isi.total} onPindah={setHalaman} />
        </>
      )}
      {isi !== null && isi.rekap.length > 0 && (
        <TabelKorpus kolom={POLIS_LAMA_REKAP[t]} baris={isi.rekap} kunci={(b, i) => `${b.PL_NUMBER_EDM ?? ''}/${b.CURRENCY ?? ''}/${i}`} />
      )}
    </Modal>
  )
}
