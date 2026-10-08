// Tab **Event Limits** (Non-Proportional) — EMPAT BARIS BERLABEL, SATU SET
// PER KONTRAK, gambar `28`.
//
//	RSMD Limit                  [ Currency ]  [ nilai ]
//	Earthquake Limit            [ Currency ]  [ nilai ]
//	Flood Limit (Jabodetabek)   [ Currency ]  [ nilai ]
//	Flood Limit (Nationwide)    [ Currency ]  [ nilai ]
//
// ---------------------------------------------------------------------
// ⛔ DIBENTUK ULANG 6 Oktober 2026 — bentuk DAN sumbernya keliru
// ---------------------------------------------------------------------
// Bentuk sebelumnya mengulang keempat baris PER LAYER × TREATY GROUP, dari
// nilai tingkat `Limits[].Detail[]`. Ekspor mengatakan lain —
// `Section/TreatyInTabsNonProportional.xml`, tab `Event Limits`, delapan sel
// `ALWAYS`:
//
//	TreatyIn.CurrencyRSMD       · TreatyIn.RSMDLimit
//	TreatyIn.CurrencyEarthquake · TreatyIn.Earthquake
//	TreatyIn.CurrencyFloodJab   · TreatyIn.FloodJab
//	TreatyIn.CurrencyFloodNat   · TreatyIn.FloodNation
//
// Properti AKAR — satu set per kontrak. Dan pengukuran 6 Oktober 2026 atas
// SELURUH 1.855 dokumen membuktikan keduanya memang tempat yang berbeda:
//
//	                 dokumen  akar berisi  Detail berisi  keduanya
//	Proportional       1.079            0            347         0
//	NonProportional      772           49              1         0
//
// ⇒ Grid per layer yang lama KOSONG pada 771 dari 772 kontrak Non-Prop,
// sementara 49 kontrak punya nilai akar yang tidak pernah tampil. (Nilai
// tingkat Detail milik sub-tab Event Limits di dalam Limits PROPORSIONAL —
// itu tetap di `TabLimitsProp`.)
//
// ⚠️ NILAI AKARNYA BELUM DAPAT DIBACA, dan itu dikatakan di layar: tidak ada
// kolom pendaratan yang memuatnya, dan kolom baru menuntut migrasi —
// rentang migrasi `treatyin` habis (larangan 6). Medan kosong di sini
// bermakna "belum terjangkau", BUKAN "kontrak ini tidak punya".
//
// ⚠️ Label sama, nasib baca-saja beda: hanya sel RSMD ber-
// `pyReadOnlyCondition` = `TreatyIn.ViewState = 1`; tiga lainnya tanpa
// syarat di selnya sendiri. Di layar ini keempatnya mengikuti `mode`, sebab
// harness Pega pun baca-saja di mode lihat.

import { useEffect, useState } from 'react'

import { Panel, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { ambilOpsiLimits, type PilihanWarisan } from '../api'
import { EVENT_LIMITS } from '../labels'
import { useProperti } from '../halaman'
import type { ModeForm } from '../mode'

/** Keempat baris — label sel dan kunci akarnya, urut ekspor/gambar 28. */
export const BARIS_EVENT_LIMITS = [
  { label: EVENT_LIMITS.rsmd, mataUang: 'CurrencyRSMD', nilai: 'RSMDLimit' },
  { label: EVENT_LIMITS.gempa, mataUang: 'CurrencyEarthquake', nilai: 'Earthquake' },
  { label: EVENT_LIMITS.banjirJab, mataUang: 'CurrencyFloodJab', nilai: 'FloodJab' },
  { label: EVENT_LIMITS.banjirNas, mataUang: 'CurrencyFloodNat', nilai: 'FloodNation' },
] as const

export default function TabEventLimits({ mode = 'lihat' }: { mode?: ModeForm }) {
  const bisaUbah = mode === 'ubah'
  // ⭐ Properti AKAR `TreatyIn.*` di PENAMPUNG HALAMAN — bertahan saat pindah
  // tab (dahulu `useState` tab ini, yang hilang begitu tab dilepas).
  const [cRsmd, setCRsmd] = useProperti('CurrencyRSMD', '')
  const [rsmd, setRsmd] = useProperti('RSMDLimit', '')
  const [cGempa, setCGempa] = useProperti('CurrencyEarthquake', '')
  const [gempa, setGempa] = useProperti('Earthquake', '')
  const [cJab, setCJab] = useProperti('CurrencyFloodJab', '')
  const [jab, setJab] = useProperti('FloodJab', '')
  const [cNas, setCNas] = useProperti('CurrencyFloodNat', '')
  const [nas, setNas] = useProperti('FloodNation', '')
  const isi: Record<string, string> = {
    CurrencyRSMD: cRsmd,
    RSMDLimit: rsmd,
    CurrencyEarthquake: cGempa,
    Earthquake: gempa,
    CurrencyFloodJab: cJab,
    FloodJab: jab,
    CurrencyFloodNat: cNas,
    FloodNation: nas,
  }
  const pengubah: Record<string, (v: string) => void> = {
    CurrencyRSMD: setCRsmd,
    RSMDLimit: setRsmd,
    CurrencyEarthquake: setCGempa,
    Earthquake: setGempa,
    CurrencyFloodJab: setCJab,
    FloodJab: setJab,
    CurrencyFloodNat: setCNas,
    FloodNation: setNas,
  }
  const setIsi = (kunci: string, v: string) => {
    pengubah[kunci]?.(v)
  }
  // `.Currency…` pxAutoComplete — daftar mata uang yang sama dengan tab
  // Limits (`BrowseCurrencyTreatyIn_RD`). Diminta hanya bila dapat diubah.
  const [mataUang, setMataUang] = useState<readonly PilihanWarisan[]>([])
  useEffect(() => {
    if (!bisaUbah) return
    let dibuang = false
    ambilOpsiLimits()
      .then((o) => {
        if (!dibuang) setMataUang(o.mataUang)
      })
      .catch(() => undefined)
    return () => {
      dibuang = true
    }
  }, [bisaUbah])

  return (
    <Panel judul={EVENT_LIMITS.judul}>
      <dl className="trin__ev-medan">
        {BARIS_EVENT_LIMITS.map((b) => (
          <div key={b.nilai} className="trin__ev-baris">
            <dt>{b.label}</dt>
            <dd className="trin__ev-mu">
              {bisaUbah ? (
                <Pilih
                  label=""
                  value={isi[b.mataUang] ?? ''}
                  opsi={mataUang.map((o) => ({ value: o.nama, label: o.nama }))}
                  onChange={(v) => {
                    setIsi(b.mataUang, v)
                  }}
                />
              ) : (
                <span className="trin__redup">—</span>
              )}
            </dd>
            <dd className="trin__ev-nilai">
              {bisaUbah ? (
                // ⛔ Tidak diformat saat diketik — koma desimal harus dapat diketik.
                <input
                  className="field__input"
                  type="text"
                  inputMode="decimal"
                  aria-label={b.label}
                  value={isi[b.nilai] ?? ''}
                  onChange={(e) => {
                    setIsi(b.nilai, e.target.value)
                  }}
                />
              ) : (
                <span className="trin__redup">—</span>
              )}
            </dd>
          </div>
        ))}
      </dl>
      <span className="trin__redup" role="note">
        {EVENT_LIMITS.belumTerjangkau}
      </span>
    </Panel>
  )
}
