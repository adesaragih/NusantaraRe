// Form Add account - tata letak tangkapan layar Pega SFAGIS Account dalam gaya form Kelola User (keputusan work owner
// 04-10-2026): Insured Name * dan Org ID, Group Business * dan Owner, Description lebar penuh. "Territory" dihapus.
// Akun hanya diinput sekali: form ini hanya untuk akun BARU - nol View, nol Edit, nol hapus.

import { useState } from 'react'

import { Area, Field, Gagal, Modal, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { tambah, type Account, type GroupBusiness, type Isian } from '../api'
import { adaGalat, isianKosong, periksa, type GalatIsian, type InsuredTerpilih } from '../aturan'
import { ACC } from '../labels'
import PilihOrganisasi from './PilihOrganisasi'

export default function FormAccount({
  groupBusiness,
  onTutup,
  onTersimpan,
}: {
  groupBusiness: GroupBusiness[]
  onTutup: () => void
  onTersimpan: (a: Account) => void
}) {
  const [isi, setIsi] = useState(isianKosong)
  const [insured, setInsured] = useState<InsuredTerpilih | null>(null)
  const [galatMedan, setGalatMedan] = useState<GalatIsian>({})
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)

  // Sesudah Save ditolak, galat medan mengikuti isian - medan yang sudah diisi tidak lagi merah.
  const isiUlang = (baru: Isian) => {
    setIsi(baru)
    if (adaGalat(galatMedan)) setGalatMedan(periksa(baru))
  }

  const simpan = () => {
    const g = periksa(isi)
    setGalatMedan(g)
    if (adaGalat(g) || sibuk) return
    setSibuk(true)
    setGalat(null)
    tambah(isi).then(
      (a) => {
        setSibuk(false)
        onTersimpan(a)
      },
      (g: unknown) => {
        setSibuk(false)
        setGalat(g)
      },
    )
  }

  return (
    <Modal
      judul={ACC.judulBaru}
      onTutup={onTutup}
      onKirim={simpan}
      labelBatal={ACC.batal}
      lebar
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk}>
          {sibuk ? ACC.menyimpan : ACC.simpan}
        </button>
      }
    >
      <div className="accounts__form">
        {galat !== null && <Gagal galat={galat} />}
        <div className="form-grid">
          <PilihOrganisasi
            terpilih={insured}
            error={galatMedan.insured}
            onPilih={(o) => {
              setInsured(o)
              isiUlang({ ...isi, insuredId: o.id })
            }}
          />
          <Field
            label={ACC.orgId}
            value={insured?.orgId ?? ''}
            placeholder={ACC.orgOtomatis}
            onChange={() => undefined}
            readOnly
          />
          <Pilih
            label={ACC.groupBusiness}
            required
            kosong={ACC.pilih}
            value={isi.groupBusinessId}
            error={galatMedan.groupBusiness}
            opsi={groupBusiness.map((b) => ({ value: b.id, label: b.note }))}
            onChange={(v) => {
              isiUlang({ ...isi, groupBusinessId: v })
            }}
          />
          <Field label={ACC.owner} value="" placeholder={ACC.ownerOtomatis} onChange={() => undefined} readOnly />
          <Area
            label={ACC.description}
            value={isi.description}
            baris={5}
            onChange={(v) => {
              isiUlang({ ...isi, description: v })
            }}
          />
        </div>
      </div>
    </Modal>
  )
}
