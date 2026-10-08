// ⛔ DIPINDAHKAN dari `pages/FormKontrakTreatyIn.tsx` 5 Oktober 2026 —
// pemindahan MURNI, nol perubahan perilaku.
//
// ---------------------------------------------------------------------
// ⛔ ADD/DELETE DICABUT DARI KEDUA GRID, 6 Oktober 2026 — dan buktinya
// ---------------------------------------------------------------------
// Kedua grid di berkas ini adalah proyeksi PER LAYER dari larik
// `TreatyIn.Share`. Di ekspor, grid itu punya tombolnya — dan keduanya MATI:
//
//   Section/TreatyInTabsNonProportional.xml, grid `TreatyIn.Share`
//   (tab `RNM Share < Title < Share`):
//     sel 382  `Add`     pyCondition = `1=2`   -> Refresh:TreatyInNonAddItem
//     sel 398  `Delete`  pyCondition = `1=2`   -> DeleteRow
//
// Barisnya di Pega LAHIR dari `Update Summary` (`TreatyInNonAddItem`,
// `param.Type=="share"`), yang membuang larik lama lalu menyusunnya ulang
// satu per layer — tidak pernah diketik satu per satu.
//
// ⚠️ YANG HIDUP DI TAB `Share` adalah grid LAIN, dan keduanya BELUM
// dibangun di aplikasi ini — jadi pencabutan ini tidak mengambil apa pun
// yang Pega izinkan:
//
//   `TreatyIn.ShareReins`                  Reinsurer Name · Layer · % Share
//     sel 280 `Add` / 285 `Delete`         `TreatyIn.ViewState !='1'`
//   `TreatyIn.ShareFacultativeReinsurers`  Facultative Reinsurers · Layer · % Share
//     non-prop sel 312/317 · prop sel 216/221   `TreatyIn.ViewState !='1'`
//
// Sisi Prop diperiksa terpisah, sebab grid ini dirender untuk keduanya:
// `TreatyInTabsProportional`, `TreatyInShareProp`, `Share`, `DetailShare`
// nol memuat Add/Delete pada grid Share per layer.
//
// ⛔ SEL TETAP DAPAT DIUBAH di mode Edit. 12 dari 13 sel properti grid
// `TreatyIn.Share` ber-`pyEditOptions` = `Auto` di ekspor; yang dicabut
// hanya kemampuan MENAMBAH dan MENGHAPUS baris.

import { useState } from 'react'

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
import { StripTabNavigasi } from './navigasi'
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
      // Sel 382/398 `1=2` — lihat kepala berkas.
      bisaTambah={false}
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
        // Proyeksi per layer larik yang sama — sel 382/398 `1=2`.
        bisaTambah={false}
      />
      <StripTabNavigasi tab={daftar} aktif={tampil} onPilih={setSub} />
      {tampil === 'RNM Share' && (
        <PanelRnmShare layer={layer} petunjukKosong={petunjukKosong} mode={mode} />
      )}
    </>
  )
}
