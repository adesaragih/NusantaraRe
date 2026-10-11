// Popup penampil View File NB FacIn (work owner 11-10-2026 "untuk view-nya buatin juga, jadi tidak perlu di download") - pola panel
// lampiran bersama `inti/frontend/lampiran/PanelLampiranReas.tsx` (seperti Product Name Life):
// - View Office Online (xls…pptx, `ReasViewAttachment` b2951): URL bertanda tangan `GET …/office` -> form GET tersembunyi ke
//   penampil (`inti/frontend/lampiran/penampil.ts`) di bingkai popup - tidak membuka jendela / tab baru;
// - View pdf / gambar: isi dari rute unduh beridentitas `…/isi` -> objek URL lokal; tipe dari ekstensi, bukan dari jawaban.
// Hanya untuk berkas yang sudah tersimpan di penyimpanan (status '').

import { useEffect, useRef, useState } from 'react'

import { Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import {
  BINGKAI_PENAMPIL,
  bisaViewOffice,
  jenisViewOnline,
  mimeViewOnline,
  PARAM_PENAMPIL,
  PENAMPIL_OFFICE,
} from '../../../../inti/frontend/lampiran/penampil'
import { ambilIsiLampiran, tautanOfficeLampiran, type BarisLampiran } from '../api'
import { TEKS_LAMPIRAN as T } from '../labels'

export type JenisPenampil = 'office' | 'pdf' | 'gambar'

/** Ekstensi dari nama berkas (tanpa titik); '' bila tidak ada. */
export function ekstensiBerkas(nama: string): string {
  const i = nama.lastIndexOf('.')
  return i < 0 ? '' : nama.slice(i + 1)
}

/** Jenis penampil untuk satu berkas; null = tombol tidak ditawarkan (berkas lain tetap lewat unduh). */
export function jenisPenampil(nama: string): JenisPenampil | null {
  const e = ekstensiBerkas(nama)
  return bisaViewOffice(e) ? 'office' : jenisViewOnline(e)
}

export default function PenampilLampiran({
  caseId,
  berkas,
  onTutup,
}: {
  caseId: string
  berkas: BarisLampiran
  /** Ditutup pemakai (tanpa argumen) atau gagal memuat (galatnya dikirim). */
  onTutup: (galat?: unknown) => void
}) {
  const jenis = jenisPenampil(berkas.namaFile)
  const ekstensi = ekstensiBerkas(berkas.namaFile)
  const [siap, setSiap] = useState(false)
  const formOffice = useRef<HTMLFormElement>(null)
  const urlOffice = useRef<HTMLInputElement>(null)
  const bingkai = useRef<HTMLIFrameElement>(null)
  const gambar = useRef<HTMLImageElement>(null)
  const tutup = useRef(onTutup)
  useEffect(() => {
    tutup.current = onTutup
  })

  useEffect(() => {
    if (jenis === null) return
    let batal = false
    let objekURL: string | null = null
    const kerja =
      jenis === 'office'
        ? tautanOfficeLampiran(caseId, berkas.id).then(({ url }) => {
            if (batal || formOffice.current === null || urlOffice.current === null) return
            urlOffice.current.value = url
            formOffice.current.submit()
          })
        : ambilIsiLampiran(caseId, berkas.id).then((isi) => {
            if (batal) return
            objekURL = URL.createObjectURL(new Blob([isi], { type: mimeViewOnline(ekstensi) }))
            if (jenis === 'gambar' && gambar.current !== null) gambar.current.src = objekURL
            if (jenis === 'pdf' && bingkai.current !== null) bingkai.current.src = objekURL
          })
    kerja.then(
      () => {
        if (!batal) setSiap(true)
      },
      (err: unknown) => {
        if (!batal) tutup.current(err)
      },
    )
    return () => {
      batal = true
      if (objekURL !== null) URL.revokeObjectURL(objekURL)
    }
  }, [caseId, berkas.id, jenis, ekstensi])

  return (
    <Modal judul={berkas.namaFile} onTutup={() => onTutup()} labelBatal={T.tutup} penuh>
      {!siap && <Memuat />}
      <form ref={formOffice} method="get" action={PENAMPIL_OFFICE} target={BINGKAI_PENAMPIL} hidden>
        <input ref={urlOffice} type="hidden" name={PARAM_PENAMPIL} />
      </form>
      {jenis === 'gambar' ? (
        <img ref={gambar} alt={berkas.namaFile} className="lampiran-reas__penampil-gambar" />
      ) : (
        <iframe ref={bingkai} name={BINGKAI_PENAMPIL} title={berkas.namaFile} className="lampiran-reas__penampil-bingkai" />
      )}
    </Modal>
  )
}
