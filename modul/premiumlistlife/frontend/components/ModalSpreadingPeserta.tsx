// Popup spreading retro SATU peserta — klik baris grid "Participant Details"
// (keputusan work owner 07-10-2026; menggantikan tab Spreading semua peserta).
//
// BACA-SAJA. Kepala: Name Of Insured, Certificate No, Currency, Retroceded Share
// peserta, treaty year. Isi: lapisan (Reins Type, Retroceded Share lapisan) →
// reinsurer (%Share, Amount, Rate, Commision %, Ovr Comm %, Gross, Net); Reins Type
// dan bagian lapisan hanya ditulis di baris pertama lapisannya.
//
// ⛔ SEMUA angka lewat SATU formatter (`pemisahRibuan`, sama dengan grid peserta):
// teks desimal dari server, tidak pernah ke float, pemisah desimal selalu koma.

import { useEffect, useState } from 'react'

import { Field, Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { pemisahRibuan } from '../angka'
import { ambilSpreadingPeserta, type SpreadingSatuPeserta } from '../api'
import { PANEL_SPREADING } from '../labels'
import { tanggalTampil } from '../tanggal'

/** Satu baris tabel popup: satu reinsurer; lapisan hanya di baris pertamanya. */
export interface BarisTabelSpreading {
  kunci: string
  /** Baris pertama sebuah lapisan (garis pemisah antar lapisan). */
  awalLapisan: boolean
  reinsType: string
  bagian: string
  reinsurer: string
  pctShare: string
  amount: string
  rate: string
  komisi: string
  ovrComm: string
  gross: string
  net: string
}

type KolomTabel = Exclude<keyof BarisTabelSpreading, 'kunci' | 'awalLapisan'>

/** Kolom tabel popup, urut tampil (sama dengan tab Spreading lama tanpa No/peserta). */
export const KOLOM_SPREADING: readonly { k: KolomTabel; angka: boolean }[] = [
  { k: 'reinsType', angka: false },
  { k: 'bagian', angka: true },
  { k: 'reinsurer', angka: false },
  { k: 'pctShare', angka: true },
  { k: 'amount', angka: true },
  { k: 'rate', angka: true },
  { k: 'komisi', angka: true },
  { k: 'ovrComm', angka: true },
  { k: 'gross', angka: true },
  { k: 'net', angka: true },
]

/** Meratakan spreading satu peserta menjadi baris tabel (satu baris per reinsurer). */
export function barisTabelSpreading(p: SpreadingSatuPeserta): BarisTabelSpreading[] {
  const keluar: BarisTabelSpreading[] = []
  ;(p.lapisan ?? []).forEach((l, li) => {
    l.reinsurer.forEach((r, ri) => {
      const awal = ri === 0
      keluar.push({
        kunci: `${String(li)}-${String(ri)}`,
        awalLapisan: awal,
        reinsType: awal ? l.reinsType : '',
        bagian: awal ? pemisahRibuan(l.bagian) : '',
        reinsurer: r.nama,
        pctShare: pemisahRibuan(r.pctShare),
        amount: pemisahRibuan(r.amount),
        rate: pemisahRibuan(r.rate),
        komisi: pemisahRibuan(r.komisi),
        ovrComm: pemisahRibuan(r.ovrComm),
        gross: pemisahRibuan(r.gross),
        net: pemisahRibuan(r.net),
      })
    })
  })
  return keluar
}

/** Judul popup: `Spreading - <Name Of Insured>`. */
export function judulModalSpreading(nama: string): string {
  return nama.trim() === '' ? PANEL_SPREADING.judul : `${PANEL_SPREADING.judul} - ${nama.trim()}`
}

/** Treaty year tampil: label lalu periode dd/mm/yyyy bila ada. */
export function treatyYearTampil(p: Pick<SpreadingSatuPeserta, 'treatyYear' | 'treatyMulai' | 'treatyAkhir'>): string {
  const mulai = tanggalTampil(p.treatyMulai)
  const akhir = tanggalTampil(p.treatyAkhir)
  if (mulai === '' && akhir === '') return p.treatyYear
  return `${p.treatyYear} (${mulai} - ${akhir})`.trim()
}

export default function ModalSpreadingPeserta({
  polisID,
  detailID,
  nama,
  onTutup,
}: {
  polisID: string
  detailID: string
  /** Name Of Insured dari baris grid - judul tampil sebelum data tiba. */
  nama: string
  onTutup: () => void
}) {
  const [isi, setIsi] = useState<SpreadingSatuPeserta | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let hidup = true
    setIsi(null)
    setGalat(null)
    ambilSpreadingPeserta(polisID, detailID).then(
      (p) => {
        if (hidup) setIsi(p)
      },
      (e: unknown) => {
        if (hidup) setGalat(e)
      },
    )
    return () => {
      hidup = false
    }
  }, [polisID, detailID])

  const baris = isi === null ? [] : barisTabelSpreading(isi)
  const k = PANEL_SPREADING.kepala

  return (
    <Modal judul={judulModalSpreading(isi?.nama ?? nama)} onTutup={onTutup} labelBatal={PANEL_SPREADING.tutup} lebar>
      {galat !== null && <Gagal galat={galat} />}
      {galat === null && isi === null && <Memuat />}
      {isi !== null && (
        <div className="pl-spreading-modal">
          <div className="form-grid">
            <Field label={k.nama} value={isi.nama} onChange={() => {}} readOnly />
            <Field label={k.sertifikat} value={isi.sertifikat} onChange={() => {}} readOnly />
            <Field label={k.currency} value={isi.currency} onChange={() => {}} readOnly />
            <Field label={k.retrocededShare} value={pemisahRibuan(isi.retrocededShare)} onChange={() => {}} readOnly />
            <Field label={k.treatyYear} value={treatyYearTampil(isi)} onChange={() => {}} readOnly />
          </div>
          {baris.length === 0 ? (
            <Kosong pesan={PANEL_SPREADING.kosong} />
          ) : (
            <div className="pl-spreading">
              <table className="inbox__tabel">
                <thead>
                  <tr>
                    {KOLOM_SPREADING.map((c) => (
                      <th key={c.k} className={c.angka ? 'pl-spreading__angka' : undefined}>
                        {PANEL_SPREADING.kolom[c.k]}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {baris.map((b) => (
                    <tr key={b.kunci} className={b.awalLapisan ? 'pl-spreading__awal' : undefined}>
                      {KOLOM_SPREADING.map((c) => (
                        <td key={c.k} className={c.angka ? 'pl-spreading__angka' : undefined}>
                          {b[c.k]}
                        </td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}
    </Modal>
  )
}
