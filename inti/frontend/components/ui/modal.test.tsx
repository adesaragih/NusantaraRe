import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { Modal } from './dasar'

// Modal: empat jalan keluar SELALU ada, kecuali modal `tanpaTutup` - popup yang hanya boleh dilanjutkan
// (perintah work owner 06-10-2026, popup nomor polis NB Treaty In: "tidak mau ada cancel, tidak ada close").

const tutup = () => undefined

describe('Modal', () => {
  it('bawaan: tombol X dan Cancel ada', () => {
    const html = renderToStaticMarkup(
      <Modal judul="UJI" onTutup={tutup} aksi={<button type="button">OK</button>}>
        isi
      </Modal>,
    )
    expect(html).toContain('modal__close')
    expect(html).toContain('>Cancel<')
    expect(html).toContain('>OK<')
  })

  it('tanpaTutup: tanpa X dan tanpa Cancel - hanya aksi', () => {
    const html = renderToStaticMarkup(
      <Modal judul="UJI" onTutup={tutup} tanpaTutup aksi={<button type="button">OK</button>}>
        isi
      </Modal>,
    )
    expect(html).not.toContain('modal__close')
    expect(html).not.toContain('>Cancel<')
    expect(html).toContain('>OK<')
    // satu aksi di tengah ("OK nya ditengah")
    expect(html).toContain('modal__actions modal__actions--tengah')
  })

  it('tanpaTutup: Escape dan klik backdrop tidak menutup', () => {
    const sumber = readFileSync(join(__dirname, 'dasar.tsx'), 'utf8')
    expect(sumber).toContain('if (e.key === "Escape" && !tanpaTutup) mulaiTutup();')
    expect(sumber).toContain('if (!tanpaTutup && e.target === e.currentTarget) mulaiTutup();')
  })
})
