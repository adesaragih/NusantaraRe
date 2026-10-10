// Satu medan layar realisasi, menurut definisinya (`medan.ts`).
//
// Angka uang ditampilkan berpasangan dengan kode mata uangnya (AC 85), dalam
// format sel Section Pega (K14, `sajian.ts`): hanya-baca selalu berformat;
// tersunting berformat selama tidak difokus bila sel ber-
// `pyShowReadonlyFormatting=true`, dan mentah saat diketik. Nilai tersimpan
// tidak pernah diubah oleh format (AC 24).

import { useEffect, useId, useRef, useState } from 'react'

import { Area, Field, Pilih, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import { nilai, type Halaman } from '../api'
import { JEDA_HITUNG_MS, hitungSaatKetik } from '../hitungLangsung'
import { teksPilihan, type Medan } from '../medan'
import { nilaiNol, sajikan, type Sajian } from '../sajian'
import InputAngka from './InputAngka'
import InputTanggal from './InputTanggal'

export interface PropsKotakMedan {
  medan: Medan
  halaman: Halaman
  wajib: boolean
  hanyaBaca: boolean
  opsiMataUang: Opsi[]
  opsiMO: Opsi[]
  onUbah: (jalur: string, v: string) => void
  onSelesai: (medan: Medan, v: string) => void
}

/** Sajian sel: dari definisinya; medan tanggal tanpa sajian tetap tanggal (AC 33). */
function sajianMedan(m: Medan): Sajian | undefined {
  return m.sajian ?? (m.jenis === 'tanggal' ? 'tanggal' : undefined)
}

export default function KotakMedan({ medan, halaman, wajib, hanyaBaca, opsiMataUang, opsiMO, onUbah, onSelesai }: PropsKotakMedan) {
  const [fokus, setFokus] = useState(false)
  const v = nilai(halaman, medan.jalur)
  const kunci = hanyaBaca || medan.kunci === true || medan.jenis === 'tampil'
  const kode = medan.mataUang ? nilai(halaman, medan.mataUang) : ''
  const pesan = halaman.pesan?.[medan.jalur]?.join('; ')
  const sajian = sajianMedan(medan)
  const idIsian = useId()
  // Hitung langsung saat mengetik (`hitungLangsung.ts`): jeda sesudah ketikan terakhir; nilai yang sudah terkirim tidak
  // dikirim ulang saat blur. `onSelesai` terkini dipakai saat jeda habis - halaman / hasil server terbaru, bukan saat
  // jeda dipasang.
  const jeda = useRef<ReturnType<typeof setTimeout> | null>(null)
  const terkirim = useRef<string | null>(null)
  const selesaiKini = useRef(onSelesai)
  selesaiKini.current = onSelesai
  useEffect(
    () => () => {
      if (jeda.current !== null) clearTimeout(jeda.current)
    },
    [],
  )

  if (kunci && medan.jenis !== 'centang') {
    // angka rata kanan; nol di bagian uang tampil "0" redup (sama dengan placeholder isian)
    const angka = sajian !== undefined && sajian !== 'tanggal' && !medan.opsi
    const nol = angka && sajian.nolPolos === true && nilaiNol(v)
    return (
      <div className="field nbti__medan">
        <span className="field__label">
          {medan.label}
          {wajib && <span className="field__req">*</span>}
        </span>
        <span
          className={'nbti__nilai' + (angka ? ' nbti__nilai--angka' : '') + (nol ? ' nbti__nilai--nol' : '')}
          data-jalur={medan.jalur}
        >
          {kode && v !== '' && <span className="nbti__kode">{kode}</span>}
          {medan.opsi ? teksPilihan(medan.opsi, v) : sajikan(v, sajian)}
        </span>
        {pesan && <div className="field__error">{pesan}</div>}
      </div>
    )
  }
  switch (medan.jenis) {
    case 'centang':
      // `pyDisabled` / hanya-baca: kotak tampil, tidak dapat diubah.
      return (
        <label className="nbti__centang">
          <input
            type="checkbox"
            checked={v === 'true'}
            disabled={kunci}
            onChange={(e) => onSelesai(medan, e.target.checked ? 'true' : 'false')}
          />
          {medan.label}
        </label>
      )
    case 'area':
      return <Area label={medan.label} value={v} required={wajib} error={pesan} onChange={(x) => onUbah(medan.jalur, x)} />
    case 'radio': {
      // pxRadioButtons: nilai tersimpan di luar daftar ditambahkan sebagai pilihan sendiri (tidak dibuang)
      const opsi = medan.opsi ?? []
      const semua = v !== '' && !opsi.some((o) => o.value === v) ? [...opsi, { value: v, label: v }] : opsi
      return (
        <div className="field">
          <span className="field__label">
            {medan.label}
            {wajib && <span className="field__req">*</span>}
          </span>
          <div className="nbti__radio" role="radiogroup" aria-label={medan.label}>
            {semua.map((o) => (
              <label key={o.value}>
                <input type="radio" name={`nbti-${medan.jalur}`} checked={v === o.value} onChange={() => onSelesai(medan, o.value)} />
                {o.label}
              </label>
            ))}
          </div>
          {pesan && <div className="field__error">{pesan}</div>}
        </div>
      )
    }
    case 'pilihan':
      return (
        <Pilih
          label={medan.label}
          value={v}
          required={wajib}
          error={pesan}
          opsi={medan.opsi ?? []}
          onChange={(x) => onSelesai(medan, x)}
        />
      )
    case 'angka':
      // isian angka (permintaan work owner 06-10-2026): hanya angka, ribuan otomatis, rata kanan, placeholder "0"
      return (
        <div className="nbti__dengan-kode">
          {kode && <span className="nbti__kode">{kode}</span>}
          <div className="field">
            <label className="field__label" htmlFor={idIsian}>
              {medan.label}
              {wajib && <span className="field__req">*</span>}
            </label>
            <InputAngka
              id={idIsian}
              label={medan.label}
              value={v}
              sajian={sajian ?? {}}
              onChange={(x) => {
                onUbah(medan.jalur, x)
                if (jeda.current !== null) clearTimeout(jeda.current)
                jeda.current = null
                if (!hitungSaatKetik(!!medan.aksi?.length, x)) return
                jeda.current = setTimeout(() => {
                  jeda.current = null
                  terkirim.current = x
                  selesaiKini.current(medan, x)
                }, JEDA_HITUNG_MS)
              }}
              onBlur={() => {
                if (jeda.current !== null) clearTimeout(jeda.current)
                jeda.current = null
                const kini = nilai(halaman, medan.jalur)
                const sudah = terkirim.current === kini
                terkirim.current = null
                if (!sudah) onSelesai(medan, kini)
              }}
            />
            {pesan && <div className="field__error">{pesan}</div>}
          </div>
        </div>
      )
    case 'mataUang':
    case 'mo':
      return (
        <Pilih
          label={medan.label}
          value={v}
          required={wajib}
          error={pesan}
          opsi={medan.jenis === 'mo' ? opsiMO : opsiMataUang}
          onChange={(x) => onSelesai(medan, x)}
        />
      )
    case 'tanggal':
      // Kotak teks dd-mm-yyyy + tombol kalender (work owner 10-10-2026: "bisa di copy paste dan di ketik lancar"),
      // bukan isian tanggal bawaan browser. Aksi medan (SystemSetOneYear, ProtectDate) berjalan saat fokus
      // MENINGGALKAN kotak dan tombolnya - pindah dari kotak ke tombol kalender bukan selesai.
      return (
        <div
          className="field"
          onBlur={(e) => {
            if (!e.currentTarget.contains(e.relatedTarget as Node | null)) onSelesai(medan, nilai(halaman, medan.jalur))
          }}
        >
          <label className="field__label" htmlFor={idIsian}>
            {medan.label}
            {wajib && <span className="field__req">*</span>}
          </label>
          <InputTanggal id={idIsian} value={v} galat={!!pesan} onChange={(x) => onUbah(medan.jalur, x)} />
          {pesan && <div className="field__error">{pesan}</div>}
        </div>
      )
    default: {
      const berformat = !fokus && sajian !== undefined && sajian !== 'tanggal' && sajian.formatSaatSunting === true
      return (
        <div
          className="nbti__dengan-kode"
          onFocus={() => setFokus(true)}
          onBlur={() => {
            setFokus(false)
            onSelesai(medan, nilai(halaman, medan.jalur))
          }}
        >
          {kode && <span className="nbti__kode">{kode}</span>}
          <Field
            label={medan.label}
            value={berformat ? sajikan(v, sajian) : v}
            required={wajib}
            error={pesan}
            type="text"
            onChange={(x) => onUbah(medan.jalur, x)}
          />
        </div>
      )
    }
  }
}
