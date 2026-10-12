// Tab **Retro** (Non-Proportional) — `Section/TreatyInTabsNonProportional.xml`,
// tab TABBED `Retro` @3291797: satu kotak centang `Has Share To Retro:`
// berpengikat `TreatyIn.IsMultipleRetro` (`pyDisabledWhen TreatyIn.ViewState = 1`).
//
// ⭐ Menggantikan tab `RNM Share` di strip Non-Prop (permintaan pemakai
// 10 Oktober 2026: "di pega tab ini tidak ada tapi tab Retro yang ada").
// `IsMultipleRetro` tersimpan sebagai TEKS `"true"`/`"false"` dan punya kolom
// pendaratan sendiri (`ISMULTIPLERETRO`), jadi Save membawanya lewat
// penampung halaman seperti properti akar lain.

import { Panel } from '../../../../inti/frontend/components/ui/dasar'
import { useProperti } from '../halaman'
import type { ModeForm } from '../mode'

export default function TabRetro({ awal = '', mode = 'lihat' }: { awal?: string; mode?: ModeForm }) {
  const [nilai, setNilai] = useProperti('IsMultipleRetro', awal)
  return (
    <Panel judul="Retro">
      <label className="trin__limit-medan">
        <span>Has Share To Retro:</span>{' '}
        <input
          type="checkbox"
          checked={nilai === 'true'}
          disabled={mode !== 'ubah'}
          onChange={(e) => {
            setNilai(e.target.checked ? 'true' : 'false')
          }}
        />
      </label>
    </Panel>
  )
}
