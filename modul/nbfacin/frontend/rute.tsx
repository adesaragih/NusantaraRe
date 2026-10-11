// Rute modul NB FacIn - tiket 21; halaman depan tiket 25/26; Inward Facultative tiket 30.

import { useEffect, useRef, useState } from 'react'

import { Modal } from '../../../inti/frontend/components/ui/dasar'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import { STATUS_PORTAL, TEKS_EMAIL_TYPE, type StatusPortal } from './labels'
import type { HalamanNbFacIn } from './menu'
import CoverageCargo from './pages/CoverageCargo'
import FormOpportunity from './pages/FormOpportunity'
import InwardFacultative, { type KasusBaru } from './pages/InwardFacultative'
import PortalOpportunity from './pages/PortalOpportunity'
import ViewPolicy from './pages/ViewPolicy'
import { pasangKlikBaris } from './components/klikBaris'
import './nbfacin.css'

export function RuteNbFacIn({ halaman, masuk, onPindah, bukaKasus, onBeranda, onBukaMenu, onLihatBerkas, onBukaModul }: PropsRute<HalamanNbFacIn>) {
  // Case yang baru dibuat Create opportunity - dibawa ke layar Inward Facultative (assignment pertama).
  // Komponen rute tetap terpasang selama modul aktif, jadi keadaan ini bertahan saat pindah halaman.
  const [kasus, setKasus] = useState<KasusBaru | null>(null)
  // Pesan konfirmasi sesudah Submit berhasil - tampil di portal (tiket 57).
  const [pesanPortal, setPesanPortal] = useState('')
  // View Policy tetap terpasang sesudah pertama dibuka: case Treaty yang dibuka di tempat (`onBukaModul`) lalu Back kembali ke
  // daftar dengan pencarian yang sama (work owner 11-10-2026).
  const [vpDibuka, setVpDibuka] = useState(false)
  useEffect(() => {
    if (halaman === 'nbfacin-view-policy') setVpDibuka(true)
  }, [halaman])
  // Posisi switch In Progress / Resolved portal - bertahan sesudah kembali dari layar case (pola NB Treaty In).
  const [statusPortal, setStatusPortal] = useState<StatusPortal>(STATUS_PORTAL[0])
  // Case dibuka dari kotak masuk Beranda: Cancel / Submit kembali ke Beranda (pola NB Treaty In).
  const [dariBeranda, setDariBeranda] = useState(false)
  // Pesan Submit berhasil untuk case dari Beranda: Beranda inti tidak menerima pesan modul, jadi tampil di jendela dulu
  // (work owner 08-10-2026: "saat di aksep uw, tadi kenapa tidak ada info nb statusnya ya"); OK -> Beranda.
  const [pesanBeranda, setPesanBeranda] = useState('')

  // Case yang dibuka langsung dari daftar kotak masuk Beranda (`PropsRute.bukaKasus`). Ketukan yang SUDAH ada saat rute
  // dipasang adalah permintaan lama - App menyimpannya melewati logout / login (work owner 08-10-2026: "baru login langsung
  // kebuka gini, harusnya kan ada di menu home terlebih dahulu") - jadi hanya ketukan BARU yang membuka case.
  // Klik di mana pun pada baris master-detail = buka / tutup (work owner 10-10-2026), semua tabel NB FacIn.
  useEffect(() => pasangKlikBaris(), [])

  const ketukBuka = bukaKasus?.ketuk
  const idBuka = bukaKasus?.id
  const ketukSaatPasang = useRef(ketukBuka)
  useEffect(() => {
    if (ketukBuka !== undefined && idBuka !== undefined && ketukBuka !== ketukSaatPasang.current) {
      setPesanPortal('')
      setKasus({ caseId: idBuka })
      setDariBeranda(true)
      onPindah('nbfacin-inward')
    }
    // onPindah sengaja tidak menjadi dependensi: pemanggil membuat fungsi baru tiap render.
  }, [ketukBuka, idBuka])

  // Halaman layar case tanpa case (mis. halaman App yang tersisa dari sebelum logout): kembali ke Beranda, bukan layar kosong.
  // Dikecualikan selama jendela pesan Submit-dari-Beranda tampil (case sudah ditutup, halaman belum pindah).
  useEffect(() => {
    if (halaman === 'nbfacin-inward' && kasus === null && pesanBeranda === '') {
      if (onBeranda) onBeranda()
      else onPindah('nbfacin-portal')
    }
    // onBeranda / onPindah sengaja tidak menjadi dependensi: pemanggil membuat fungsi baru tiap render.
  }, [halaman, kasus, pesanBeranda])

  /** Keluar dari layar case: ke Beranda bila dibuka dari sana, selain itu ke portal. */
  const keluarKasus = () => {
    if (dariBeranda && onBeranda) {
      setDariBeranda(false)
      onBeranda()
    } else onPindah('nbfacin-portal')
  }
  return (
    <>
      {halaman === 'nbfacin-portal' && (
        <PortalOpportunity
          onBuat={() => {
            setPesanPortal('')
            onPindah('nbfacin-opportunity')
          }}
          onBuka={(caseId) => {
            setPesanPortal('')
            setDariBeranda(false)
            setKasus({ caseId })
            onPindah('nbfacin-inward')
          }}
          pesan={pesanPortal}
          onTutupPesan={() => setPesanPortal('')}
          status={statusPortal}
          onStatus={setStatusPortal}
        />
      )}
      {vpDibuka && (
        <div hidden={halaman !== 'nbfacin-view-policy'}>
          <ViewPolicy onLihatBerkas={onLihatBerkas} onBukaModul={onBukaModul} />
        </div>
      )}
      {halaman === 'nbfacin-opportunity' && (
        <FormOpportunity
          pemilik={masuk.nama ?? masuk.akunID}
          onBukaMenu={onBukaMenu}
          onDibuat={(k) => {
            setKasus(k)
            onPindah('nbfacin-inward')
          }}
        />
      )}
      {halaman === 'nbfacin-inward' && kasus && (
        <InwardFacultative
          kasus={kasus}
          onBatal={keluarKasus}
          onSelesai={(pesan) => {
            // Submit berhasil: case ditutup (Pega `finishAssignment`) dan kembali ke tempat case dibuka - Beranda bila dari
            // kotak masuk (work owner 07-10-2026: "hbs di send Teamleader knpa lari kesini ? harusnya tetap di inbox"),
            // selain itu portal dengan pesan konfirmasi.
            setKasus(null)
            if (dariBeranda && onBeranda) {
              setPesanBeranda(pesan)
              return
            }
            setPesanPortal(pesan)
            onPindah('nbfacin-portal')
          }}
        />
      )}
      {halaman === 'nbfacin-coverage-cargo' && <CoverageCargo />}
      {pesanBeranda !== '' && (
        <Modal
          judul={TEKS_EMAIL_TYPE.judulSelesai}
          labelBatal={TEKS_EMAIL_TYPE.ok}
          onTutup={() => {
            setPesanBeranda('')
            setDariBeranda(false)
            onBeranda?.()
          }}
        >
          <p>{pesanBeranda}</p>
        </Modal>
      )}
    </>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanNbFacIn> = RuteNbFacIn
