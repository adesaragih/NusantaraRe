// ⛔ DIPINDAHKAN dari `pages/FormKontrakTreatyIn.tsx` 5 Oktober 2026 —
// pemindahan MURNI, nol perubahan perilaku.

import { useState } from 'react'

import { StripTab } from '../../../../inti/frontend/components/ui/dasar'
import type { BarisLayerWarisan } from '../api'
import {
  JENIS_RNM_SHARE,
  JENIS_SHARE,
  KOLOM_RNM_SHARE,
  KOLOM_SHARE,
  SUB_TAB_SHARE,
} from '../labels'
import type { ModeForm } from '../mode'
import { barisLayer } from './angka'
import TabGridWarisan from './TabGridWarisan'

export function PanelRnmShare({
  layer,
  petunjukKosong,
  mode = 'lihat',
}: {
  layer: readonly BarisLayerWarisan[]
  petunjukKosong: string
  mode?: ModeForm
}) {
  return (
    <TabGridWarisan
      judul="RNM Share"
      kolom={KOLOM_RNM_SHARE}
      baris={barisLayer(
        layer,
        (b) => [
          b.layer, b.rnmShare, b.liabilityRNM, b.mdpRNM100,
          b.epiRNMQS100, b.rnmRetainedPremi, b.rnmQSPremi,
        ],
        JENIS_RNM_SHARE,
      )}
      petunjukKosong={petunjukKosong}
      mode={mode}
    />
  )
}

/**
 * Tab `Share` beserta strip sub-tabnya — bentuk yang gambar 16/17/34
 * perlihatkan.
 *
 * ⚠️ Satu sub-tab hari ini (`RNM Share`). Strip bersatu butir terlihat
 * ganjil, dan ia tetap dipasang: gambar `16` memperlihatkannya begitu, dan
 * strip yang baru muncul saat butir kedua datang akan memindahkan isi di
 * bawah mata pemakai tanpa ada yang mengubah apa pun.
 */
export default function SubTabShare({
  layer,
  petunjukKosong,
  mode = 'lihat',
}: {
  layer: readonly BarisLayerWarisan[]
  petunjukKosong: string
  mode?: ModeForm
}) {
  const [sub, setSub] = useState<string>(SUB_TAB_SHARE[0])
  const daftar: readonly string[] = SUB_TAB_SHARE
  const tampil = daftar.includes(sub) ? sub : SUB_TAB_SHARE[0]
  return (
    <>
      <TabGridWarisan
        judul="Share"
        kolom={KOLOM_SHARE}
        baris={barisLayer(
          layer,
          (b) => [
            b.layer, b.persenCession, b.jenisPenyebaran, b.persenBrokerage,
            b.cessionKeRI, b.qsor, b.qsri, b.liabilityQSOR, b.liabilityQSRI,
            b.epi100, b.riogr,
          ],
          JENIS_SHARE,
        )}
        petunjukKosong={petunjukKosong}
        mode={mode}
      />
      <StripTab tab={daftar} aktif={tampil} onPilih={setSub} />
      {tampil === 'RNM Share' && (
        <PanelRnmShare layer={layer} petunjukKosong={petunjukKosong} mode={mode} />
      )}
    </>
  )
}
