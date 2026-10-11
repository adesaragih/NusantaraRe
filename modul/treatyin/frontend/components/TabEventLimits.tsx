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
// ⭐ RALAT 8 Oktober 2026 — NILAINYA SUDAH DAPAT DIBACA.
//
// Keterangan lama berbunyi *"belum ada kolom pendaratan yang memuatnya"*. Itu
// benar saat ditulis 6 Oktober, dan KEDALUWARSA sejak migrasi `446`
// mendirikan `T_TREATY_HAZARD_LIMIT` — tabel yang kedelapan medan ini
// tempati, dan yang Save memang tulis.
//
// ⛔ Yang hilang selama ini PEMBACANYA: nol jalur membaca tabel itu kembali,
// sehingga tab ini selalu kosong walau kontraknya punya nilai. Pembacanya
// lahir 8 Oktober 2026 (`repository.BacaBatasBahaya`), dan nilainya tiba di
// sini lewat penampung halaman seperti medan tersimpan lain.
//
// ⚠️ Keterangan kedaluwarsa yang terdengar pasti adalah sebab cacat ini
// bertahan dua hari: ia menjawab pertanyaan "kenapa kosong?" dengan jawaban
// yang salah, dan nol orang memeriksanya lagi.
//
// ⚠️ Label sama, nasib baca-saja beda: hanya sel RSMD ber-
// `pyReadOnlyCondition` = `TreatyIn.ViewState = 1`; tiga lainnya tanpa
// syarat di selnya sendiri. Di layar ini keempatnya mengikuti `mode`, sebab
// harness Pega pun baca-saja di mode lihat.

import { useEffect, useState } from 'react'

import { FieldAngka, Panel } from '../../../../inti/frontend/components/ui/dasar'
import { PilihCari as Pilih } from '../../../../inti/frontend/components/ui/pilihSaring'
import { formatNumber } from '../../../../inti/frontend/lib/format'
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
                  // ⭐ Pilihan kosong berbunyi `Currency`, seperti di Pega —
                  // bukan teks kosong bawaan aplikasi.
                  kosong={EVENT_LIMITS.mataUangKosong}
                  opsi={mataUang.map((o) => ({ value: o.nama, label: o.nama }))}
                  onChange={(v) => {
                    setIsi(b.mataUang, v)
                  }}
                />
              ) : (
                <span>{isi[b.mataUang] || <span className="trin__redup">—</span>}</span>
              )}
            </dd>
            <dd className="trin__ev-nilai">
              {bisaUbah ? (
                // ⭐ `FieldAngka` — pemisah ribuan hidup saat mengetik, sama
                // dengan seluruh isian uang modul ini. Bentuk KABEL-nya tetap
                // titik desimal, jadi yang tersimpan tidak berubah.
                <FieldAngka
                  label=""
                  value={isi[b.nilai] ?? ''}
                  desimal={2}
                  // ⛔ `0,00` sebagai BAYANGAN, bukan nilai. Pega memang
                  // memperlihatkan `0,00` pada medan kosong, tetapi menuliskan
                  // nol ke dalamnya akan mengubah data: "belum diisi" dan
                  // "diisi nol" bukan hal yang sama, dan Save membedakannya.
                  placeholder={EVENT_LIMITS.nilaiKosong}
                  onChange={(v) => {
                    setIsi(b.nilai, v)
                  }}
                />
              ) : (
                // ⛔ Mode lihat DULU selalu `—`, bahkan ketika nilainya ada:
                // tab ini tidak pernah memperlihatkan apa pun. Sekarang ia
                // menampilkan angkanya, berpemisah ribuan seperti grid uang
                // lain, dan `—` hanya untuk yang benar-benar kosong.
                <span>
                  {isi[b.nilai] ? formatNumber(isi[b.nilai], 2) : <span className="trin__redup">—</span>}
                </span>
              )}
            </dd>
          </div>
        ))}
      </dl>
    </Panel>
  )
}
