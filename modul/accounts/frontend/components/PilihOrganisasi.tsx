// Pilihan "Insured Name" - dropdown bersaring inti (`PilihSaring`): daftar terbuka hanya saat diklik / diketik dan
// tertutup sesudah memilih, Escape, atau fokus pergi. Organisasi `CLIENT` (FLAG Org, puluhan ribu baris) dicari di
// server, paling banyak 50 per jawaban; kata kosong = 50 nama pertama. Memilih satu mengisi Insured Name dan Org ID.

import { useCallback, useRef, useState } from 'react'

import { Gagal } from '../../../../inti/frontend/components/ui/dasar'
import { PilihSaring, type OpsiSaring } from '../../../../inti/frontend/components/ui/pilihSaring'
import { cariOrganisasi } from '../api'
import type { InsuredTerpilih } from '../aturan'
import { ACC } from '../labels'

export default function PilihOrganisasi({
  terpilih,
  error,
  onPilih,
}: {
  terpilih: InsuredTerpilih | null
  error?: string
  onPilih: (o: InsuredTerpilih) => void
}) {
  const [opsi, setOpsi] = useState<OpsiSaring[]>([])
  const [memuat, setMemuat] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  // Jawaban pencarian yang sudah disusul pencarian lebih baru dibuang.
  const giliran = useRef(0)

  const cari = useCallback((kata: string) => {
    const ke = ++giliran.current
    setMemuat(true)
    cariOrganisasi(kata).then(
      (h) => {
        if (ke !== giliran.current) return
        setOpsi(h.daftar.map((o) => ({ value: o.id, label: o.nama, keterangan: o.idView })))
        setGalat(null)
        setMemuat(false)
      },
      (g: unknown) => {
        if (ke !== giliran.current) return
        setOpsi([])
        setGalat(g)
        setMemuat(false)
      },
    )
  }, [])

  return (
    <div>
      <PilihSaring
        label={ACC.insuredName}
        required
        value={terpilih?.id ?? ''}
        teksTerpilih={terpilih?.nama ?? ''}
        opsi={opsi}
        memuat={memuat}
        onCari={cari}
        onPilih={(o) => {
          onPilih({ id: o.value, orgId: o.keterangan ?? '', nama: o.label })
        }}
      />
      {galat !== null && <Gagal galat={galat} />}
      {error && <div className="field__error">{error}</div>}
    </div>
  )
}
